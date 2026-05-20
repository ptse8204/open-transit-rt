package main

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/gtfs"
)

const (
	gtfsBrowserImportMaxBytes       = 64 << 20
	gtfsBrowserImportMemoryBytes    = 1 << 20
	gtfsImportPreflightReadLimit    = 4 << 20
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

type gtfsImportSourceReview struct {
	SourceType              string
	SourceURL               string
	ChecksumSHA256          string
	ByteCount               int64
	ImportTimestamp         *time.Time
	ActiveFeedVersion       string
	ScheduleIdentitySummary string
	UpdateComparison        string
	RollbackVisibility      string
	Preflight               []operationsGTFSChangeRow
	NextActions             []string
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
		h.renderGTFSImportStatus(w, r, principal, http.StatusRequestEntityTooLarge, "", "GTFS import request was blocked because the form body is invalid or too large.", nil, gtfsImportSourceReview{})
		return
	}
	if principal.Method == auth.MethodCookie && strings.TrimSpace(h.csrfSecret) != "" && strings.TrimSpace(r.FormValue("csrf_token")) != csrfToken(h.csrfSecret, principal) {
		h.renderGTFSImportStatus(w, r, principal, http.StatusForbidden, "", "GTFS import request was blocked because the CSRF token is invalid.", nil, gtfsImportSourceReview{})
		return
	}
	if err := rejectGTFSImportUnexpectedFields(r); err != nil {
		h.renderGTFSImportStatus(w, r, principal, http.StatusBadRequest, "", err.Error(), nil, gtfsImportSourceReview{})
		return
	}
	if strings.TrimSpace(r.FormValue("action")) != "import_gtfs" {
		h.renderGTFSImportStatus(w, r, principal, http.StatusBadRequest, "", "GTFS import accepts only the import_gtfs action.", nil, gtfsImportSourceReview{})
		return
	}
	if h.gtfsImport == nil {
		h.renderGTFSImportStatus(w, r, principal, http.StatusServiceUnavailable, "", "GTFS import service is not available in this runtime. Use the CLI import path or start the console with a database-backed runtime.", nil, gtfsImportSourceReview{})
		return
	}
	notes, err := gtfsImportNotes(r.FormValue("notes"))
	if err != nil {
		h.renderGTFSImportStatus(w, r, principal, http.StatusBadRequest, "", err.Error(), nil, gtfsImportSourceReview{})
		return
	}

	sourceType := strings.TrimSpace(r.FormValue("source_type"))
	tmpPath, cleanup, err := h.gtfsImportTempPath(r, sourceType)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		h.renderGTFSImportStatus(w, r, principal, http.StatusBadRequest, "", err.Error(), nil, gtfsImportSourceReview{})
		return
	}
	sourceReview := gtfsImportSourceReviewFromPath(tmpPath, sourceType, r.FormValue("gtfs_url"), time.Now().UTC().Truncate(time.Second))

	ctx, cancel := context.WithTimeout(r.Context(), gtfsBrowserImportTimeout())
	defer cancel()
	result, err := h.gtfsImport.ImportZip(ctx, gtfs.ImportOptions{
		AgencyID: principal.AgencyID,
		ZipPath:  tmpPath,
		ActorID:  principal.Subject,
		Notes:    notes,
	})
	view := gtfsImportResultFromResult(result)
	sourceReview.ActiveFeedVersion = firstNonEmpty(result.FeedVersionID, sourceReview.ActiveFeedVersion)
	sourceReview.ScheduleIdentitySummary = gtfsImportScheduleIdentity(result)
	if err != nil {
		var importErr *gtfs.ImportError
		if errors.As(err, &importErr) {
			view = gtfsImportResultFromResult(importErr.Result)
			sourceReview.ActiveFeedVersion = firstNonEmpty(importErr.Result.FeedVersionID, sourceReview.ActiveFeedVersion)
			sourceReview.ScheduleIdentitySummary = gtfsImportScheduleIdentity(importErr.Result)
			h.renderGTFSImportStatus(w, r, principal, http.StatusOK, "", gtfsImportSafeFailure(importErr.Result), view, sourceReview)
			return
		}
		h.renderGTFSImportStatus(w, r, principal, http.StatusInternalServerError, "", "GTFS import could not finish. Check database/import service health and retry from the CLI if the browser path remains unavailable.", view, sourceReview)
		return
	}
	h.renderGTFSImportStatus(w, r, principal, http.StatusOK, "GTFS import finished through the existing import pipeline. Review GTFS quality, validator health, and feed health before treating the feed as ready.", "", view, sourceReview)
}

func (h *handler) renderGTFSImportStatus(w http.ResponseWriter, r *http.Request, principal auth.Principal, status int, notice string, message string, result *gtfsImportResultView, source gtfsImportSourceReview) {
	page := h.buildOperationsPage(r, principal, "gtfs-import")
	page.GTFSImportNotice = notice
	page.GTFSImportError = message
	page.GTFSImportResult = result
	page.GTFSImportSource = source
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

func gtfsImportSourceReviewFromPath(path string, sourceType string, rawURL string, at time.Time) gtfsImportSourceReview {
	review := gtfsImportSourceReview{
		SourceType:              gtfsImportSourceTypeLabel(sourceType, rawURL),
		SourceURL:               gtfsImportSourceURL(sourceType, rawURL),
		ImportTimestamp:         &at,
		ScheduleIdentitySummary: "not available until import completes",
		UpdateComparison:        "Current active schedule versus new import comparison is not exposed by the current browser model. Review counts and validation blockers before relying on an activation decision.",
		RollbackVisibility:      "The browser shows the active feed version after import. Prior feed-version listing or rollback execution is not implemented here; use the documented operator rollback path when needed.",
		NextActions: []string{
			"Review GTFS quality.",
			"Review feed health.",
			"Run or review validator health.",
			"Use rollback documentation or operator support if the active feed must be reverted.",
		},
	}
	file, err := os.Open(path)
	if err != nil {
		review.ChecksumSHA256 = "not available"
		review.Preflight = gtfsImportZipPreflight(path)
		return review
	}
	defer file.Close()
	hash := sha256.New()
	n, err := io.Copy(hash, file)
	if err != nil {
		review.ChecksumSHA256 = "not available"
		review.Preflight = gtfsImportZipPreflight(path)
		return review
	}
	review.ByteCount = n
	review.ChecksumSHA256 = hex.EncodeToString(hash.Sum(nil))
	review.Preflight = gtfsImportZipPreflight(path)
	return review
}

func gtfsImportZipPreflight(zipPath string) []operationsGTFSChangeRow {
	const boundary = "ZIP preflight reads only the temporary import source and renders bounded file/row signals. It does not validate GTFS semantics, publish a feed, retain evidence, or prove schedule correctness."
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return []operationsGTFSChangeRow{gtfsWorkbenchChangeRow("ZIP structure", "blocked", "Temporary source could not be opened as a GTFS ZIP before import.", "Use a valid GTFS ZIP export and retry the existing import path.", boundary)}
	}
	defer reader.Close()
	counts := map[string]int{}
	truncated := map[string]bool{}
	for _, file := range reader.File {
		name := strings.ToLower(path.Base(strings.TrimSpace(file.Name)))
		if name == "." || name == "/" || !strings.HasSuffix(name, ".txt") {
			continue
		}
		rows, capped := gtfsImportZipTextRows(file)
		counts[name] += rows
		truncated[name] = truncated[name] || capped
	}
	return []operationsGTFSChangeRow{
		gtfsImportZipRequiredRow("agency.txt", counts, truncated, true, boundary),
		gtfsImportZipRequiredRow("routes.txt", counts, truncated, true, boundary),
		gtfsImportZipRequiredRow("stops.txt", counts, truncated, true, boundary),
		gtfsImportZipRequiredRow("trips.txt", counts, truncated, true, boundary),
		gtfsImportZipRequiredRow("stop_times.txt", counts, truncated, true, boundary),
		gtfsImportZipCalendarRow(counts, truncated, boundary),
		gtfsImportZipRequiredRow("frequencies.txt", counts, truncated, false, boundary),
		gtfsImportZipRequiredRow("shapes.txt", counts, truncated, false, boundary),
	}
}

func gtfsImportZipTextRows(file *zip.File) (int, bool) {
	rc, err := file.Open()
	if err != nil {
		return 0, false
	}
	defer rc.Close()
	limit := int64(gtfsImportPreflightReadLimit)
	buf := make([]byte, 32*1024)
	newlines := 0
	var read int64
	var last byte
	sawBytes := false
	for {
		if read >= limit {
			return dataRowsFromLineScan(newlines, sawBytes, last), true
		}
		max := int64(len(buf))
		if remaining := limit - read; remaining < max {
			max = remaining
		}
		n, err := rc.Read(buf[:max])
		if n > 0 {
			read += int64(n)
			newlines += strings.Count(string(buf[:n]), "\n")
			last = buf[n-1]
			sawBytes = true
		}
		if err == io.EOF {
			return dataRowsFromLineScan(newlines, sawBytes, last), false
		}
		if err != nil {
			return dataRowsFromLineScan(newlines, sawBytes, last), false
		}
	}
}

func dataRowsFromLineScan(newlines int, sawBytes bool, last byte) int {
	if !sawBytes {
		return 0
	}
	lines := newlines
	if last != '\n' {
		lines++
	}
	dataRows := lines - 1
	if dataRows < 0 {
		return 0
	}
	return dataRows
}

func gtfsImportZipRequiredRow(file string, counts map[string]int, truncated map[string]bool, required bool, boundary string) operationsGTFSChangeRow {
	rows, ok := counts[file]
	if !ok {
		if required {
			return gtfsWorkbenchChangeRow(file, "blocked", file+" is missing from the temporary ZIP preflight.", "Export a complete GTFS ZIP with required files before retrying import.", boundary)
		}
		return gtfsWorkbenchChangeRow(file, "optional", file+" is not present in the temporary ZIP preflight.", "If this file is expected, confirm the source export settings before import.", boundary)
	}
	status := "ok"
	if required && rows == 0 {
		status = "blocked"
	}
	if !required && rows == 0 {
		status = "needs_review"
	}
	return gtfsWorkbenchChangeRow(file, status, gtfsImportZipRowSignal(file, rows, truncated[file]), "Compare this bounded preflight count with importer counts, GTFS quality, and validator results after import.", boundary)
}

func gtfsImportZipCalendarRow(counts map[string]int, truncated map[string]bool, boundary string) operationsGTFSChangeRow {
	calendarRows, hasCalendar := counts["calendar.txt"]
	exceptionRows, hasExceptions := counts["calendar_dates.txt"]
	if !hasCalendar && !hasExceptions {
		return gtfsWorkbenchChangeRow("calendar.txt / calendar_dates.txt", "blocked", "Neither calendar.txt nor calendar_dates.txt is present in the temporary ZIP preflight.", "Add service calendars or exception-only service dates before retrying import.", boundary)
	}
	status := "ok"
	if calendarRows+exceptionRows == 0 {
		status = "blocked"
	}
	signal := fmt.Sprintf("calendar.txt has %d data row(s); calendar_dates.txt has %d data row(s).", calendarRows, exceptionRows)
	if truncated["calendar.txt"] || truncated["calendar_dates.txt"] {
		signal += " One or both scans reached the bounded preflight read cap."
	}
	return gtfsWorkbenchChangeRow("calendar.txt / calendar_dates.txt", status, signal, "Confirm service date coverage, holiday exceptions, and after-midnight trips after import.", boundary)
}

func gtfsImportZipRowSignal(file string, rows int, capped bool) string {
	signal := fmt.Sprintf("%s has %d bounded data row(s) in the temporary ZIP preflight.", file, rows)
	if capped {
		signal += " The scan reached the bounded preflight read cap."
	}
	return signal
}

func gtfsImportSourceTypeLabel(sourceType string, rawURL string) string {
	if sourceType == "url" || strings.TrimSpace(rawURL) != "" {
		return "URL import"
	}
	return "uploaded ZIP"
}

func gtfsImportSourceURL(sourceType string, rawURL string) string {
	if sourceType != "url" && strings.TrimSpace(rawURL) == "" {
		return "not applicable for upload"
	}
	if gtfsImportTextLooksSecret(rawURL) {
		return "withheld because URL looked secret-bearing"
	}
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Hostname() == "" {
		return "URL accepted for import; display omitted"
	}
	if gtfsImportHostLooksPrivate(parsed.Hostname()) {
		return "withheld because URL host is local or private"
	}
	display := parsed.Scheme + "://" + parsed.Host
	if base := path.Base(parsed.EscapedPath()); base != "." && base != "/" && base != "" {
		display += "/" + base
	}
	if parsed.RawQuery != "" {
		display += "?query=omitted"
	}
	return display
}

func gtfsImportScheduleIdentity(result gtfs.ImportResult) string {
	if result.Status == "" {
		return "not available until import completes"
	}
	parts := []string{}
	for _, key := range []string{"routes", "stops", "trips", "stop_times", "shapes"} {
		if value, ok := result.Counts[key]; ok {
			parts = append(parts, fmt.Sprintf("%s=%d", key, value))
		}
	}
	if len(parts) == 0 {
		return "import completed without row-count summary"
	}
	return strings.Join(parts, "; ")
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
