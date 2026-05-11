package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/gtfs"
)

const (
	gtfsBrowserImportMaxBytes       = 64 << 20
	gtfsBrowserImportMemoryBytes    = 1 << 20
	defaultGTFSBrowserImportTimeout = 15 * time.Minute
	defaultGTFSURLDownloadTimeout   = 30 * time.Second
)

type gtfsImportRunner interface {
	ImportZip(ctx context.Context, opts gtfs.ImportOptions) (gtfs.ImportResult, error)
}

type gtfsImportResultView struct {
	ImportID       int64
	AgencyID       string
	FeedVersionID  string
	Status         string
	ErrorCount     int
	WarningCount   int
	InfoCount      int
	Counts         []countView
	ReportStored   bool
	FailureMessage string
}

func (h *handler) renderGTFSImport(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "gtfs-import")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "gtfs-import", page)
}

func (h *handler) operationsGTFSImportPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, gtfsBrowserImportMaxBytes+gtfsBrowserImportMemoryBytes)
	principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	if err := parseGTFSImportForm(r); err != nil {
		h.renderGTFSImportStatus(w, r, principal, http.StatusRequestEntityTooLarge, "", "GTFS import request was blocked because the form body is invalid or too large.", nil)
		return
	}
	if principal.Method == auth.MethodCookie && strings.TrimSpace(h.csrfSecret) != "" && strings.TrimSpace(r.FormValue("csrf_token")) != csrfToken(h.csrfSecret, principal) {
		h.renderGTFSImportStatus(w, r, principal, http.StatusForbidden, "", "GTFS import request was blocked because the CSRF token is invalid.", nil)
		return
	}
	if err := rejectGTFSImportUnexpectedFields(r); err != nil {
		h.renderGTFSImportStatus(w, r, principal, http.StatusBadRequest, "", err.Error(), nil)
		return
	}
	if strings.TrimSpace(r.FormValue("action")) != "import_gtfs" {
		h.renderGTFSImportStatus(w, r, principal, http.StatusBadRequest, "", "GTFS import accepts only the import_gtfs action.", nil)
		return
	}
	if h.gtfsImport == nil {
		h.renderGTFSImportStatus(w, r, principal, http.StatusServiceUnavailable, "", "GTFS import service is not available in this runtime. Use the CLI import path or start the console with a database-backed runtime.", nil)
		return
	}
	notes, err := gtfsImportNotes(r.FormValue("notes"))
	if err != nil {
		h.renderGTFSImportStatus(w, r, principal, http.StatusBadRequest, "", err.Error(), nil)
		return
	}

	sourceType := strings.TrimSpace(r.FormValue("source_type"))
	tmpPath, cleanup, err := h.gtfsImportTempPath(r, sourceType)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		h.renderGTFSImportStatus(w, r, principal, http.StatusBadRequest, "", err.Error(), nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), gtfsBrowserImportTimeout())
	defer cancel()
	result, err := h.gtfsImport.ImportZip(ctx, gtfs.ImportOptions{
		AgencyID: principal.AgencyID,
		ZipPath:  tmpPath,
		ActorID:  principal.Subject,
		Notes:    notes,
	})
	view := gtfsImportResultFromResult(result)
	if err != nil {
		var importErr *gtfs.ImportError
		if errors.As(err, &importErr) {
			view = gtfsImportResultFromResult(importErr.Result)
			h.renderGTFSImportStatus(w, r, principal, http.StatusOK, "", gtfsImportSafeFailure(importErr.Result), view)
			return
		}
		h.renderGTFSImportStatus(w, r, principal, http.StatusInternalServerError, "", "GTFS import could not finish. Check database/import service health and retry from the CLI if the browser path remains unavailable.", view)
		return
	}
	h.renderGTFSImportStatus(w, r, principal, http.StatusOK, "GTFS import finished through the existing import pipeline. Review GTFS quality, validator health, and feed health before treating the feed as ready.", "", view)
}

func (h *handler) renderGTFSImportStatus(w http.ResponseWriter, r *http.Request, principal auth.Principal, status int, notice string, message string, result *gtfsImportResultView) {
	page := h.buildOperationsPage(r, principal, "gtfs-import")
	page.GTFSImportNotice = notice
	page.GTFSImportError = message
	page.GTFSImportResult = result
	w.Header().Set("Cache-Control", "no-store")
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	renderOperationsTemplate(w, "gtfs-import", page)
}

func (h *handler) gtfsImportTempPath(r *http.Request, sourceType string) (string, func(), error) {
	switch sourceType {
	case "upload":
		return writeGTFSUploadTemp(r)
	case "url":
		return downloadGTFSURLTemp(r.Context(), r.FormValue("gtfs_url"))
	default:
		if hasGTFSUpload(r) {
			return writeGTFSUploadTemp(r)
		}
		if strings.TrimSpace(r.FormValue("gtfs_url")) != "" {
			return downloadGTFSURLTemp(r.Context(), r.FormValue("gtfs_url"))
		}
		return "", nil, fmt.Errorf("choose ZIP upload or URL import")
	}
}

func rejectGTFSImportUnexpectedFields(r *http.Request) error {
	allowed := map[string]bool{
		"action":      true,
		"csrf_token":  true,
		"source_type": true,
		"gtfs_url":    true,
		"notes":       true,
	}
	for name := range r.Form {
		if !allowed[name] {
			return fmt.Errorf("GTFS import accepts only source, URL or ZIP, notes, action, and CSRF fields")
		}
	}
	if r.MultipartForm != nil {
		for name := range r.MultipartForm.File {
			if name != "gtfs_zip" {
				return fmt.Errorf("GTFS import accepts only the gtfs_zip file field")
			}
		}
	}
	return nil
}

func parseGTFSImportForm(r *http.Request) error {
	contentType := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
	if strings.HasPrefix(contentType, "multipart/form-data") {
		return r.ParseMultipartForm(gtfsBrowserImportMemoryBytes)
	}
	return r.ParseForm()
}

func hasGTFSUpload(r *http.Request) bool {
	if r.MultipartForm == nil {
		return false
	}
	files := r.MultipartForm.File["gtfs_zip"]
	return len(files) > 0
}

func writeGTFSUploadTemp(r *http.Request) (string, func(), error) {
	file, header, err := r.FormFile("gtfs_zip")
	if err != nil {
		return "", nil, fmt.Errorf("choose a GTFS ZIP file")
	}
	defer file.Close()
	if header != nil && header.Filename != "" && !strings.HasSuffix(strings.ToLower(header.Filename), ".zip") {
		return "", nil, fmt.Errorf("GTFS upload must be a .zip file")
	}
	return writeGTFSReaderTemp(file)
}

func downloadGTFSURLTemp(ctx context.Context, raw string) (string, func(), error) {
	parsed, err := validateGTFSImportURL(raw)
	if err != nil {
		return "", nil, err
	}
	client := gtfsImportHTTPClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return "", nil, fmt.Errorf("build GTFS URL request: %w", err)
	}
	req.Header.Set("Accept", "application/zip, application/octet-stream, */*")
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("download GTFS ZIP: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", nil, fmt.Errorf("download GTFS ZIP returned HTTP %d", resp.StatusCode)
	}
	return writeGTFSReaderTemp(resp.Body)
}

func writeGTFSReaderTemp(reader io.Reader) (string, func(), error) {
	tmp, err := os.CreateTemp("", "open-transit-rt-gtfs-*.zip")
	if err != nil {
		return "", nil, fmt.Errorf("create temporary GTFS ZIP: %w", err)
	}
	path := tmp.Name()
	cleanup := func() { _ = os.Remove(path) }
	written, copyErr := io.Copy(tmp, io.LimitReader(reader, gtfsBrowserImportMaxBytes+1))
	closeErr := tmp.Close()
	if copyErr != nil {
		cleanup()
		return "", nil, fmt.Errorf("store temporary GTFS ZIP: %w", copyErr)
	}
	if closeErr != nil {
		cleanup()
		return "", nil, fmt.Errorf("close temporary GTFS ZIP: %w", closeErr)
	}
	if written == 0 {
		cleanup()
		return "", nil, fmt.Errorf("GTFS ZIP is empty")
	}
	if written > gtfsBrowserImportMaxBytes {
		cleanup()
		return "", nil, fmt.Errorf("GTFS ZIP exceeds the browser import limit of 64 MiB")
	}
	return path, cleanup, nil
}

func validateGTFSImportURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("enter a GTFS ZIP URL")
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("GTFS ZIP URL is invalid")
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return nil, fmt.Errorf("GTFS ZIP URL must use http or https")
	}
	if parsed.Scheme == "http" && !allowInsecureBrowserGTFSImportHTTP() {
		return nil, fmt.Errorf("GTFS ZIP URL must use https unless insecure browser imports are explicitly allowed for local testing")
	}
	if parsed.User != nil {
		return nil, fmt.Errorf("GTFS ZIP URL must not include embedded credentials")
	}
	if parsed.Fragment != "" {
		return nil, fmt.Errorf("GTFS ZIP URL must not include a fragment")
	}
	if gtfsImportQueryLooksSecret(parsed.RawQuery) {
		return nil, fmt.Errorf("GTFS ZIP URL query looks secret-bearing; use a public URL or upload the ZIP instead")
	}
	host := parsed.Hostname()
	if host == "" {
		return nil, fmt.Errorf("GTFS ZIP URL must include a host")
	}
	if !allowPrivateBrowserGTFSImportURLs() && gtfsImportHostLooksPrivate(host) {
		return nil, fmt.Errorf("GTFS ZIP URL host is private, local, or internal; use upload for local files")
	}
	return parsed, nil
}

func gtfsImportHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	if !allowPrivateBrowserGTFSImportURLs() {
		dialer := &net.Dialer{Timeout: 10 * time.Second}
		transport.DialContext = func(ctx context.Context, network string, address string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(address)
			if err != nil {
				return nil, err
			}
			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, err
			}
			if len(ips) == 0 {
				return nil, fmt.Errorf("resolve GTFS ZIP URL host: no addresses")
			}
			for _, addr := range ips {
				if gtfsImportIPLooksPrivate(addr.IP) {
					return nil, fmt.Errorf("GTFS ZIP URL resolved to a private, local, or internal address")
				}
			}
			return dialer.DialContext(ctx, network, net.JoinHostPort(ips[0].IP.String(), port))
		}
	}
	client := &http.Client{
		Timeout:   defaultGTFSURLDownloadTimeout,
		Transport: transport,
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return fmt.Errorf("stopped after 10 GTFS ZIP URL redirects")
		}
		if _, err := validateGTFSImportURL(req.URL.String()); err != nil {
			return err
		}
		return nil
	}
	return client
}

func gtfsImportHostLooksPrivate(host string) bool {
	host = strings.TrimSpace(strings.TrimSuffix(strings.ToLower(host), "."))
	if host == "" || host == "localhost" || strings.HasSuffix(host, ".localhost") || strings.HasSuffix(host, ".local") || strings.HasSuffix(host, ".internal") || !strings.Contains(host, ".") {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return gtfsImportIPLooksPrivate(ip)
	}
	return false
}

func gtfsImportIPLooksPrivate(ip net.IP) bool {
	return ip == nil || ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() || ip.IsMulticast()
}

func gtfsImportQueryLooksSecret(rawQuery string) bool {
	query := strings.ToLower(rawQuery)
	if query == "" {
		return false
	}
	for _, marker := range []string{"token", "secret", "password", "passwd", "authorization", "cookie", "x-amz-signature", "sig=", "key="} {
		if strings.Contains(query, marker) {
			return true
		}
	}
	return false
}

func gtfsImportNotes(userNotes string) (string, error) {
	userNotes = strings.TrimSpace(userNotes)
	if len(userNotes) > 500 {
		return "", fmt.Errorf("GTFS import notes are too long")
	}
	if gtfsImportTextLooksSecret(userNotes) {
		return "", fmt.Errorf("GTFS import notes look secret-bearing; remove credentials, tokens, cookies, or private paths")
	}
	if userNotes == "" {
		return "browser_gtfs_import", nil
	}
	return "browser_gtfs_import: " + userNotes, nil
}

func gtfsImportTextLooksSecret(value string) bool {
	lower := strings.ToLower(value)
	for _, marker := range []string{"authorization:", "set-cookie", "cookie:", "token", "secret", "password", "database_url", "restore_database_url", "postgres://", "private key", "/users/", "/var/lib", "/etc/"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func gtfsImportResultFromResult(result gtfs.ImportResult) *gtfsImportResultView {
	if result.Status == "" && result.ImportID == 0 && result.FeedVersionID == "" {
		return nil
	}
	counts := make([]countView, 0, len(result.Counts))
	for label, count := range result.Counts {
		counts = append(counts, countView{Label: label, Count: count})
	}
	sort.Slice(counts, func(i, j int) bool { return counts[i].Label < counts[j].Label })
	return &gtfsImportResultView{
		ImportID:       result.ImportID,
		AgencyID:       result.AgencyID,
		FeedVersionID:  result.FeedVersionID,
		Status:         result.Status,
		ErrorCount:     result.ErrorCount,
		WarningCount:   result.WarningCount,
		InfoCount:      result.InfoCount,
		Counts:         counts,
		ReportStored:   result.ReportStored,
		FailureMessage: gtfsImportDisplayFailure(result.FailureMessage),
	}
}

func gtfsImportSafeFailure(result gtfs.ImportResult) string {
	if failure := gtfsImportDisplayFailure(result.FailureMessage); failure != "" {
		return "GTFS import finished with stored validation feedback: " + failure + ". Review GTFS quality before retrying or publishing."
	}
	return "GTFS import finished with validation feedback. Review GTFS quality before retrying or publishing."
}

func gtfsImportDisplayFailure(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if gtfsImportTextLooksSecret(value) {
		return "details withheld from the page; review the stored validation report"
	}
	return value
}

func gtfsBrowserImportTimeout() time.Duration {
	raw := strings.TrimSpace(os.Getenv("BROWSER_GTFS_IMPORT_TIMEOUT"))
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("GTFS_IMPORT_TIMEOUT"))
	}
	if raw == "" {
		return defaultGTFSBrowserImportTimeout
	}
	parsed, err := time.ParseDuration(raw)
	if err != nil || parsed <= 0 {
		return defaultGTFSBrowserImportTimeout
	}
	return parsed
}

func allowPrivateBrowserGTFSImportURLs() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("ALLOW_BROWSER_GTFS_IMPORT_PRIVATE_URLS")), "true")
}

func allowInsecureBrowserGTFSImportHTTP() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("ALLOW_BROWSER_GTFS_IMPORT_INSECURE_HTTP")), "true")
}
