package avladapter

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/telemetry"
)

const (
	SendManifestSchemaVersion = "avl-adapter-send.v1"

	EnvTelemetryURL      = "AVL_ADAPTER_TELEMETRY_URL"
	EnvOutputDir         = "AVL_ADAPTER_OUTPUT_DIR"
	EnvTimeout           = "AVL_ADAPTER_TIMEOUT"
	EnvMaxRetries        = "AVL_ADAPTER_MAX_RETRIES"
	EnvRetryInitialDelay = "AVL_ADAPTER_RETRY_INITIAL_DELAY"
	EnvRetryMaxDelay     = "AVL_ADAPTER_RETRY_MAX_DELAY"
	EnvFailOnWarnings    = "AVL_ADAPTER_FAIL_ON_WARNINGS"
	EnvReferenceTime     = "AVL_ADAPTER_REFERENCE_TIME"
	EnvStaleThreshold    = "AVL_ADAPTER_STALE_THRESHOLD"
	EnvFutureThreshold   = "AVL_ADAPTER_FUTURE_THRESHOLD"
)

const (
	DefaultSendTimeout           = 10 * time.Second
	DefaultSendMaxRetries        = 2
	DefaultSendRetryInitialDelay = 250 * time.Millisecond
	DefaultSendRetryMaxDelay     = 2 * time.Second
)

const (
	CodeInvalidSendManifest      = "invalid_send_manifest"
	CodeInvalidSendConfig        = "invalid_send_config"
	CodeInvalidTelemetryTarget   = "invalid_telemetry_target"
	CodeInvalidOutputPath        = "invalid_output_path"
	CodeMissingCredentialMapping = "missing_credential_mapping"
	CodeMissingCredentialToken   = "missing_credential_token"
	CodeSendBlockedWarning       = "send_blocked_warning"
	CodeRedactionScanFailure     = "redaction_scan_failure"
)

var safeEnvName = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

type SendManifest struct {
	SchemaVersion   string               `json:"schema_version"`
	TelemetryURLEnv string               `json:"telemetry_url_env"`
	Credentials     []ManifestCredential `json:"credentials"`
}

type ManifestCredential struct {
	AgencyID  string `json:"agency_id"`
	DeviceID  string `json:"device_id"`
	VehicleID string `json:"vehicle_id"`
	TokenEnv  string `json:"token_env"`
	Notes     string `json:"notes,omitempty"`
}

type SendConfig struct {
	TelemetryURL      string
	OutputDir         string
	Timeout           time.Duration
	MaxRetries        int
	RetryInitialDelay time.Duration
	RetryMaxDelay     time.Duration
	FailOnWarnings    bool
	ReferenceTime     time.Time
	StaleThreshold    time.Duration
	FutureThreshold   time.Duration
	GeneratedAt       time.Time
}

type Target struct {
	URL      string
	Path     string
	Loopback bool
	HostRef  string
}

type OutputTarget struct {
	Dir   string
	Label string
	Ref   string
}

type PreparedCredential struct {
	Ref      string
	TokenEnv string
	Token    string
}

type SendPrepared struct {
	Manifest      SendManifest
	Config        SendConfig
	Target        Target
	Output        OutputTarget
	Credentials   map[credentialKey]PreparedCredential
	EventRefs     []string
	Forbidden     []string
	RawPayloadIDs []string
}

type SendReport struct {
	Summary                 Summary
	Diagnostics             []SendDiagnostic
	Manifest                RedactedManifest
	ResponseForbiddenValues []string
	Files                   map[string][]byte
	Stdout                  string
}

type Summary struct {
	GeneratedAt                     string `json:"generated_at"`
	Mode                            string `json:"mode"`
	DryRun                          bool   `json:"dry_run"`
	TelemetryURLPath                string `json:"telemetry_url_path"`
	TelemetryTargetLoopback         bool   `json:"telemetry_target_loopback"`
	TelemetryHostRef                string `json:"telemetry_host_ref,omitempty"`
	OutputLabel                     string `json:"output_label"`
	OutputRef                       string `json:"output_ref,omitempty"`
	TransformedCount                int    `json:"transformed_count"`
	SentCount                       int    `json:"sent_count"`
	SucceededCount                  int    `json:"succeeded_count"`
	FailedCount                     int    `json:"failed_count"`
	SkippedCount                    int    `json:"skipped_count"`
	RetryTotal                      int    `json:"retry_total"`
	DurationMS                      int64  `json:"duration_ms"`
	ExternalEvidenceCreated         bool   `json:"external_evidence_created"`
	ConsumerStatusesChanged         bool   `json:"consumer_statuses_changed"`
	ComplianceClaimed               bool   `json:"compliance_claimed"`
	ProductionReadinessClaimed      bool   `json:"production_readiness_claimed"`
	HostedSaaSClaimed               bool   `json:"hosted_saas_claimed"`
	AgencyAdoptionClaimed           bool   `json:"agency_adoption_claimed"`
	ConsumerAcceptanceClaimed       bool   `json:"consumer_acceptance_claimed"`
	VendorCompatibilityClaimed      bool   `json:"vendor_compatibility_claimed"`
	ProductionAVLReliabilityClaimed bool   `json:"production_avl_reliability_claimed"`
	ProductionGradeETAClaimed       bool   `json:"production_grade_eta_claimed"`
}

type SendDiagnostic struct {
	RecordIndex     int               `json:"record_index"`
	CredentialRef   string            `json:"credential_ref"`
	Outcome         string            `json:"outcome"`
	Attempts        int               `json:"attempts"`
	DurationMS      int64             `json:"duration_ms"`
	HTTPStatus      int               `json:"http_status,omitempty"`
	ResponseSHA256  string            `json:"response_sha256,omitempty"`
	SafeSuccessData map[string]string `json:"safe_success_fields,omitempty"`
}

type RedactedManifest struct {
	SchemaVersion   string                       `json:"schema_version"`
	TelemetryURLEnv string                       `json:"telemetry_url_env"`
	Credentials     []RedactedManifestCredential `json:"credentials"`
}

type RedactedManifestCredential struct {
	CredentialRef string `json:"credential_ref"`
	TokenEnv      string `json:"token_env"`
	Notes         string `json:"notes,omitempty"`
}

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Sleeper func(context.Context, time.Duration) error

type credentialKey struct {
	agencyID  string
	deviceID  string
	vehicleID string
}

type Environment func(string) string

func LoadSendManifest(r io.Reader) (SendManifest, []Diagnostic) {
	var manifest SendManifest
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&manifest); err != nil {
		return SendManifest{}, []Diagnostic{errorDiagnostic(CodeInvalidSendManifest, "invalid send manifest JSON", nil)}
	}
	return manifest, ValidateSendManifest(manifest)
}

func ValidateSendManifest(manifest SendManifest) []Diagnostic {
	var diagnostics []Diagnostic
	if manifest.SchemaVersion != SendManifestSchemaVersion {
		diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidSendManifest, "send manifest schema_version must be avl-adapter-send.v1", nil))
	}
	if manifest.TelemetryURLEnv != EnvTelemetryURL || !safeEnvName.MatchString(manifest.TelemetryURLEnv) {
		diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidSendManifest, "telemetry_url_env must be AVL_ADAPTER_TELEMETRY_URL", nil))
	}
	if secretLike(manifest.SchemaVersion) || secretLike(manifest.TelemetryURLEnv) {
		diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidSendManifest, "send manifest contains a secret-like value", nil))
	}
	seen := map[credentialKey]bool{}
	for index, row := range manifest.Credentials {
		idx := index
		if row.AgencyID == "" || row.DeviceID == "" || row.VehicleID == "" || row.TokenEnv == "" {
			diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidSendManifest, "agency_id, device_id, vehicle_id, and token_env are required in credential rows", &idx))
			continue
		}
		if !safeEnvName.MatchString(row.TokenEnv) {
			diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidSendManifest, "token_env must be a safe environment variable name", &idx))
		}
		for _, value := range []string{row.AgencyID, row.DeviceID, row.VehicleID, row.TokenEnv, row.Notes} {
			if secretLike(value) {
				diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidSendManifest, "send manifest contains a secret-like value", &idx))
				break
			}
		}
		key := credentialKey{agencyID: row.AgencyID, deviceID: row.DeviceID, vehicleID: row.VehicleID}
		if seen[key] {
			diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidSendManifest, "duplicate credential row for agency_id, device_id, and vehicle_id", &idx))
			continue
		}
		seen[key] = true
	}
	return diagnostics
}

func SendConfigFromEnv(manifest SendManifest, getenv Environment, now time.Time) (SendConfig, []Diagnostic) {
	if getenv == nil {
		getenv = os.Getenv
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	cfg := SendConfig{
		TelemetryURL:      strings.TrimSpace(getenv(manifest.TelemetryURLEnv)),
		OutputDir:         strings.TrimSpace(getenv(EnvOutputDir)),
		Timeout:           DefaultSendTimeout,
		MaxRetries:        DefaultSendMaxRetries,
		RetryInitialDelay: DefaultSendRetryInitialDelay,
		RetryMaxDelay:     DefaultSendRetryMaxDelay,
		ReferenceTime:     now.UTC(),
		StaleThreshold:    DefaultStaleThreshold,
		FutureThreshold:   DefaultFutureThreshold,
		GeneratedAt:       now.UTC(),
	}
	var diagnostics []Diagnostic
	parseDurationEnv(getenv, EnvTimeout, &cfg.Timeout, &diagnostics)
	parseIntEnv(getenv, EnvMaxRetries, &cfg.MaxRetries, &diagnostics)
	parseDurationEnv(getenv, EnvRetryInitialDelay, &cfg.RetryInitialDelay, &diagnostics)
	parseDurationEnv(getenv, EnvRetryMaxDelay, &cfg.RetryMaxDelay, &diagnostics)
	parseBoolEnv(getenv, EnvFailOnWarnings, &cfg.FailOnWarnings, &diagnostics)
	parseTimeEnv(getenv, EnvReferenceTime, &cfg.ReferenceTime, &diagnostics)
	parseDurationEnv(getenv, EnvStaleThreshold, &cfg.StaleThreshold, &diagnostics)
	parseDurationEnv(getenv, EnvFutureThreshold, &cfg.FutureThreshold, &diagnostics)
	if cfg.TelemetryURL == "" {
		diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidSendConfig, "telemetry URL environment variable is required for send mode", nil))
	}
	if cfg.Timeout <= 0 || cfg.MaxRetries < 0 || cfg.RetryInitialDelay < 0 || cfg.RetryMaxDelay < 0 || cfg.StaleThreshold <= 0 || cfg.FutureThreshold <= 0 {
		diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidSendConfig, "send mode durations and retry counts must be non-negative, with positive timeout and thresholds", nil))
	}
	return cfg, diagnostics
}

func PrepareSend(manifest SendManifest, cfg SendConfig, result Result, payload Payload, getenv Environment, cwd string) (SendPrepared, []Diagnostic) {
	if getenv == nil {
		getenv = os.Getenv
	}
	var diagnostics []Diagnostic
	target, targetDiagnostics := ValidateTelemetryTarget(cfg.TelemetryURL)
	diagnostics = append(diagnostics, targetDiagnostics...)
	output, outputDiagnostics := ResolveOutputTarget(cfg.OutputDir, cfg.GeneratedAt, cwd)
	diagnostics = append(diagnostics, outputDiagnostics...)
	preparedCredentials, credentialDiagnostics, forbidden := PrepareCredentials(manifest, getenv)
	diagnostics = append(diagnostics, credentialDiagnostics...)
	eventRefs := make([]string, 0, len(result.Events))
	for _, event := range result.Events {
		key := credentialKey{agencyID: event.AgencyID, deviceID: event.DeviceID, vehicleID: event.VehicleID}
		credential, ok := preparedCredentials[key]
		if !ok {
			diagnostics = append(diagnostics, errorDiagnostic(CodeMissingCredentialMapping, "send manifest is missing a credential mapping for a transformed event", nil))
			eventRefs = append(eventRefs, "")
			continue
		}
		eventRefs = append(eventRefs, credential.Ref)
	}
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Severity != SeverityWarning {
			continue
		}
		if diagnostic.Code == CodeStaleTimestamp || diagnostic.Code == CodeFutureTimestamp || cfg.FailOnWarnings {
			diagnostics = append(diagnostics, errorDiagnostic(CodeSendBlockedWarning, "send mode blocked by transformed-record warning", diagnostic.Index))
		}
	}
	return SendPrepared{
		Manifest:      manifest,
		Config:        cfg,
		Target:        target,
		Output:        output,
		Credentials:   preparedCredentials,
		EventRefs:     eventRefs,
		Forbidden:     forbidden,
		RawPayloadIDs: RawPayloadValues(payload),
	}, diagnostics
}

func ValidateTelemetryTarget(rawURL string) (Target, []Diagnostic) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return Target{}, []Diagnostic{errorDiagnostic(CodeInvalidTelemetryTarget, "telemetry target URL is invalid", nil)}
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.Path != "/v1/telemetry" {
		return Target{}, []Diagnostic{errorDiagnostic(CodeInvalidTelemetryTarget, "telemetry target must be exactly /v1/telemetry without userinfo, query, or fragment", nil)}
	}
	host := parsed.Hostname()
	loopback := isLoopbackHost(host)
	if parsed.Scheme != "https" && !loopback {
		return Target{}, []Diagnostic{errorDiagnostic(CodeInvalidTelemetryTarget, "non-loopback credentialed telemetry sends require HTTPS", nil)}
	}
	return Target{
		URL:      parsed.String(),
		Path:     parsed.Path,
		Loopback: loopback,
		HostRef:  "host:" + shortHash(strings.ToLower(parsed.Host)),
	}, nil
}

func ResolveOutputTarget(raw string, generatedAt time.Time, cwd string) (OutputTarget, []Diagnostic) {
	label := "default"
	if strings.TrimSpace(raw) == "" {
		raw = filepath.Join(".cache", "avl-vendor-adapter", generatedAt.Format("20060102T150405Z"))
	} else {
		label = "override"
	}
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return OutputTarget{}, []Diagnostic{errorDiagnostic(CodeInvalidOutputPath, "output path cannot be resolved", nil)}
		}
	}
	abs := raw
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(cwd, raw)
	}
	abs = filepath.Clean(abs)
	rel, err := filepath.Rel(cwd, abs)
	if err != nil {
		return OutputTarget{}, []Diagnostic{errorDiagnostic(CodeInvalidOutputPath, "output path cannot be resolved", nil)}
	}
	if rel == "." {
		return OutputTarget{}, []Diagnostic{errorDiagnostic(CodeInvalidOutputPath, "output path must be a directory below the repository root", nil)}
	}
	if strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return OutputTarget{}, []Diagnostic{errorDiagnostic(CodeInvalidOutputPath, "output path must stay below the repository root", nil)}
	}
	if rel == "docs"+string(filepath.Separator)+"evidence" || strings.HasPrefix(rel, "docs"+string(filepath.Separator)+"evidence"+string(filepath.Separator)) {
		return OutputTarget{}, []Diagnostic{errorDiagnostic(CodeInvalidOutputPath, "send output must not be written under docs/evidence", nil)}
	}
	if hasSymlinkAncestor(abs) {
		return OutputTarget{}, []Diagnostic{errorDiagnostic(CodeInvalidOutputPath, "send output path must not include symlinks", nil)}
	}
	if entries, err := os.ReadDir(abs); err == nil && len(entries) > 0 {
		return OutputTarget{}, []Diagnostic{errorDiagnostic(CodeInvalidOutputPath, "send output directory must be empty when it already exists", nil)}
	}
	ref := redactedOutputRef(rel)
	return OutputTarget{Dir: abs, Label: label, Ref: ref}, nil
}

func PrepareCredentials(manifest SendManifest, getenv Environment) (map[credentialKey]PreparedCredential, []Diagnostic, []string) {
	if getenv == nil {
		getenv = os.Getenv
	}
	credentials := map[credentialKey]PreparedCredential{}
	var diagnostics []Diagnostic
	var forbidden []string
	for index, row := range manifest.Credentials {
		idx := index
		token := getenv(row.TokenEnv)
		if token == "" {
			diagnostics = append(diagnostics, errorDiagnostic(CodeMissingCredentialToken, "credential token environment variable is required for send mode", &idx))
		} else {
			forbidden = append(forbidden, token)
		}
		key := credentialKey{agencyID: row.AgencyID, deviceID: row.DeviceID, vehicleID: row.VehicleID}
		credentials[key] = PreparedCredential{
			Ref:      CredentialRef(row.AgencyID, row.DeviceID, row.VehicleID),
			TokenEnv: row.TokenEnv,
			Token:    token,
		}
	}
	return credentials, diagnostics, forbidden
}

func SendEvents(ctx context.Context, events []telemetry.Event, prepared SendPrepared, client HTTPDoer, sleeper Sleeper) SendReport {
	if client == nil {
		client = http.DefaultClient
	}
	if sleeper == nil {
		sleeper = SleepContext
	}
	started := time.Now()
	diagnostics := make([]SendDiagnostic, 0, len(events))
	succeeded := 0
	failed := 0
	skipped := 0
	retryTotal := 0
	stopped := false
	var responseForbidden []string
	for index, event := range events {
		credential := prepared.Credentials[credentialKey{agencyID: event.AgencyID, deviceID: event.DeviceID, vehicleID: event.VehicleID}]
		if stopped {
			skipped++
			diagnostics = append(diagnostics, SendDiagnostic{RecordIndex: index, CredentialRef: credential.Ref, Outcome: "skipped_after_failure"})
			continue
		}
		diag, retries, ok, forbiddenValues := sendOne(ctx, index, event, prepared, credential, client, sleeper)
		responseForbidden = append(responseForbidden, forbiddenValues...)
		retryTotal += retries
		diagnostics = append(diagnostics, diag)
		if ok {
			succeeded++
			continue
		}
		failed++
		stopped = true
	}
	sent := succeeded + failed
	duration := time.Since(started)
	summary := buildSummary(prepared, len(events), sent, succeeded, failed, skipped, retryTotal, duration)
	return SendReport{
		Summary:                 summary,
		Diagnostics:             diagnostics,
		Manifest:                RedactManifest(prepared.Manifest),
		ResponseForbiddenValues: responseForbidden,
	}
}

func BuildSendFiles(report SendReport, prepared SendPrepared) (SendReport, []Diagnostic) {
	files := map[string][]byte{}
	putJSON := func(name string, value any) []Diagnostic {
		raw, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			return []Diagnostic{errorDiagnostic(CodeInvalidSendConfig, "send output could not be encoded", nil)}
		}
		files[name] = append(raw, '\n')
		return nil
	}
	if diagnostics := putJSON("summary.json", report.Summary); len(diagnostics) > 0 {
		return report, diagnostics
	}
	if diagnostics := putJSON("manifest.json", report.Manifest); len(diagnostics) > 0 {
		return report, diagnostics
	}
	if diagnostics := putJSON("diagnostics.json", report.Diagnostics); len(diagnostics) > 0 {
		return report, diagnostics
	}
	files["summary.md"] = []byte(renderSummaryMarkdown(report.Summary))
	files["manifest.md"] = []byte(renderManifestMarkdown(report.Manifest))
	forbidden := append([]string{}, prepared.Forbidden...)
	forbidden = append(forbidden, prepared.RawPayloadIDs...)
	forbidden = append(forbidden, report.ResponseForbiddenValues...)
	for name, raw := range files {
		if hits := ScanRedaction(string(raw), forbidden); len(hits) > 0 {
			return report, []Diagnostic{errorDiagnostic(CodeRedactionScanFailure, "redaction scan failed for "+name, nil)}
		}
	}
	report.Files = files
	report.Stdout = fmt.Sprintf("{\"mode\":\"send\",\"dry_run\":false,\"transformed_count\":%d,\"sent_count\":%d,\"succeeded_count\":%d,\"failed_count\":%d,\"skipped_count\":%d,\"output_ref\":%q}\n",
		report.Summary.TransformedCount,
		report.Summary.SentCount,
		report.Summary.SucceededCount,
		report.Summary.FailedCount,
		report.Summary.SkippedCount,
		report.Summary.OutputRef,
	)
	if hits := ScanRedaction(report.Stdout, forbidden); len(hits) > 0 {
		return report, []Diagnostic{errorDiagnostic(CodeRedactionScanFailure, "redaction scan failed for terminal output", nil)}
	}
	return report, nil
}

func WriteSendFiles(report SendReport, output OutputTarget) error {
	expected := []string{"diagnostics.json", "manifest.json", "manifest.md", "summary.json", "summary.md"}
	if len(report.Files) != len(expected) {
		return fmt.Errorf("send output file set is incomplete")
	}
	for _, name := range expected {
		if _, ok := report.Files[name]; !ok {
			return fmt.Errorf("send output file set is missing %s", name)
		}
	}
	if err := os.MkdirAll(output.Dir, 0o755); err != nil {
		return fmt.Errorf("create send output directory: %w", err)
	}
	if hasSymlinkAncestor(output.Dir) {
		return fmt.Errorf("send output path must not include symlinks")
	}
	for _, name := range expected {
		path := filepath.Join(output.Dir, name)
		if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("send output file path must not be a symlink")
		}
		if err := os.WriteFile(path, report.Files[name], 0o644); err != nil {
			return fmt.Errorf("write send output file: %w", err)
		}
	}
	entries, err := os.ReadDir(output.Dir)
	if err != nil {
		return fmt.Errorf("read send output directory: %w", err)
	}
	if len(entries) != len(expected) {
		return fmt.Errorf("send output directory contains unexpected files")
	}
	return nil
}

func RawPayloadValues(payload Payload) []string {
	seen := map[string]bool{}
	add := func(value string) {
		value = strings.TrimSpace(value)
		if value != "" && len(value) >= 4 {
			seen[value] = true
		}
	}
	add(payload.VendorSource)
	for _, record := range payload.Records {
		add(record.VendorDeviceID)
		add(record.VendorVehicleID)
	}
	values := make([]string, 0, len(seen))
	for value := range seen {
		values = append(values, value)
	}
	sort.Strings(values)
	return values
}

func DecodePayload(r io.Reader) (Payload, []Diagnostic) {
	var payload Payload
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		return Payload{}, []Diagnostic{errorDiagnostic(CodeInvalidPayloadJSON, "invalid vendor payload JSON", nil)}
	}
	return payload, nil
}

func CredentialRef(agencyID string, deviceID string, vehicleID string) string {
	return "cred:" + shortHash(agencyID+"|"+deviceID+"|"+vehicleID)
}

func SleepContext(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func ScanRedaction(content string, forbidden []string) []string {
	var hits []string
	lower := strings.ToLower(content)
	for _, value := range forbidden {
		if value == "" {
			continue
		}
		if strings.Contains(content, value) {
			hits = append(hits, "forbidden_value")
		}
	}
	for _, pattern := range []string{
		`authorization\s*:`,
		`\bbearer\s+[a-z0-9._~+/=-]+`,
		`\bcookie\s*[:=]`,
		`postgres(?:ql)?://`,
		`-----begin [^-]*private key-----`,
	} {
		if regexp.MustCompile(pattern).MatchString(lower) {
			hits = append(hits, "forbidden_pattern")
		}
	}
	return hits
}

func RedactManifest(manifest SendManifest) RedactedManifest {
	credentials := make([]RedactedManifestCredential, 0, len(manifest.Credentials))
	for _, row := range manifest.Credentials {
		credentials = append(credentials, RedactedManifestCredential{
			CredentialRef: CredentialRef(row.AgencyID, row.DeviceID, row.VehicleID),
			TokenEnv:      row.TokenEnv,
			Notes:         row.Notes,
		})
	}
	return RedactedManifest{
		SchemaVersion:   manifest.SchemaVersion,
		TelemetryURLEnv: manifest.TelemetryURLEnv,
		Credentials:     credentials,
	}
}

func ParsePayloadForSend(raw []byte) (Payload, Result) {
	payload, diagnostics := DecodePayload(bytes.NewReader(raw))
	if len(diagnostics) > 0 {
		return Payload{}, Result{Diagnostics: diagnostics}
	}
	return payload, Result{}
}

func TransformDecodedPayload(payload Payload, mapping MappingFile, options Options) Result {
	raw, err := json.Marshal(payload)
	if err != nil {
		return Result{Diagnostics: []Diagnostic{errorDiagnostic(CodeInvalidPayloadJSON, "invalid vendor payload JSON", nil)}}
	}
	return TransformPayload(bytes.NewReader(raw), mapping, options)
}

func sendOne(ctx context.Context, index int, event telemetry.Event, prepared SendPrepared, credential PreparedCredential, client HTTPDoer, sleeper Sleeper) (SendDiagnostic, int, bool, []string) {
	started := time.Now()
	attempts := 0
	retries := 0
	var lastStatus int
	var lastHash string
	var safe map[string]string
	var outcome string
	for {
		attempts++
		raw, err := json.Marshal(event)
		if err != nil {
			return SendDiagnostic{RecordIndex: index, CredentialRef: credential.Ref, Outcome: "encode_failed", Attempts: attempts, DurationMS: durationMS(started)}, retries, false, nil
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, prepared.Target.URL, bytes.NewReader(raw))
		if err != nil {
			return SendDiagnostic{RecordIndex: index, CredentialRef: credential.Ref, Outcome: "request_failed", Attempts: attempts, DurationMS: durationMS(started)}, retries, false, nil
		}
		req.Header.Set("Authorization", "Bearer "+credential.Token)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			outcome = "retryable_error"
			if attempts <= prepared.Config.MaxRetries {
				if sleepErr := sleeper(ctx, retryDelay(prepared.Config, retries)); sleepErr != nil {
					return SendDiagnostic{RecordIndex: index, CredentialRef: credential.Ref, Outcome: "retry_sleep_failed", Attempts: attempts, DurationMS: durationMS(started)}, retries, false, nil
				}
				retries++
				continue
			}
			return SendDiagnostic{RecordIndex: index, CredentialRef: credential.Ref, Outcome: outcome, Attempts: attempts, DurationMS: durationMS(started)}, retries, false, nil
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		_ = resp.Body.Close()
		forbiddenValues := responseForbiddenValues(body)
		lastStatus = resp.StatusCode
		if len(body) > 0 {
			sum := sha256.Sum256(body)
			lastHash = hex.EncodeToString(sum[:])
			safe = parseSafeSuccessFields(body)
		}
		switch {
		case resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusAccepted:
			return SendDiagnostic{RecordIndex: index, CredentialRef: credential.Ref, Outcome: "succeeded", Attempts: attempts, DurationMS: durationMS(started), HTTPStatus: lastStatus, ResponseSHA256: lastHash, SafeSuccessData: safe}, retries, true, forbiddenValues
		case isRetryableStatus(resp.StatusCode):
			outcome = "retryable_status"
			if attempts <= prepared.Config.MaxRetries {
				if sleepErr := sleeper(ctx, retryDelay(prepared.Config, retries)); sleepErr != nil {
					return SendDiagnostic{RecordIndex: index, CredentialRef: credential.Ref, Outcome: "retry_sleep_failed", Attempts: attempts, DurationMS: durationMS(started), HTTPStatus: lastStatus, ResponseSHA256: lastHash}, retries, false, forbiddenValues
				}
				retries++
				continue
			}
			return SendDiagnostic{RecordIndex: index, CredentialRef: credential.Ref, Outcome: outcome, Attempts: attempts, DurationMS: durationMS(started), HTTPStatus: lastStatus, ResponseSHA256: lastHash}, retries, false, forbiddenValues
		default:
			return SendDiagnostic{RecordIndex: index, CredentialRef: credential.Ref, Outcome: "terminal_status", Attempts: attempts, DurationMS: durationMS(started), HTTPStatus: lastStatus, ResponseSHA256: lastHash}, retries, false, forbiddenValues
		}
	}
}

func buildSummary(prepared SendPrepared, transformed int, sent int, succeeded int, failed int, skipped int, retryTotal int, duration time.Duration) Summary {
	return Summary{
		GeneratedAt:                     prepared.Config.GeneratedAt.Format(time.RFC3339),
		Mode:                            "send",
		DryRun:                          false,
		TelemetryURLPath:                prepared.Target.Path,
		TelemetryTargetLoopback:         prepared.Target.Loopback,
		TelemetryHostRef:                prepared.Target.HostRef,
		OutputLabel:                     prepared.Output.Label,
		OutputRef:                       prepared.Output.Ref,
		TransformedCount:                transformed,
		SentCount:                       sent,
		SucceededCount:                  succeeded,
		FailedCount:                     failed,
		SkippedCount:                    skipped,
		RetryTotal:                      retryTotal,
		DurationMS:                      duration.Milliseconds(),
		ExternalEvidenceCreated:         false,
		ConsumerStatusesChanged:         false,
		ComplianceClaimed:               false,
		ProductionReadinessClaimed:      false,
		HostedSaaSClaimed:               false,
		AgencyAdoptionClaimed:           false,
		ConsumerAcceptanceClaimed:       false,
		VendorCompatibilityClaimed:      false,
		ProductionAVLReliabilityClaimed: false,
		ProductionGradeETAClaimed:       false,
	}
}

func isRetryableStatus(status int) bool {
	return status == http.StatusRequestTimeout || status == http.StatusTooManyRequests || status >= 500
}

func retryDelay(cfg SendConfig, retry int) time.Duration {
	delay := cfg.RetryInitialDelay
	for i := 0; i < retry; i++ {
		delay *= 2
	}
	if cfg.RetryMaxDelay > 0 && delay > cfg.RetryMaxDelay {
		return cfg.RetryMaxDelay
	}
	return delay
}

func parseSafeSuccessFields(body []byte) map[string]string {
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}
	allowed := map[string]bool{"accepted": true, "ingest_status": true}
	safe := map[string]string{}
	for key, value := range raw {
		if !allowed[key] {
			continue
		}
		switch typed := value.(type) {
		case string:
			safe[key] = typed
		case bool:
			safe[key] = strconv.FormatBool(typed)
		}
	}
	if len(safe) == 0 {
		return nil
	}
	return safe
}

func responseForbiddenValues(body []byte) []string {
	if len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		value := strings.TrimSpace(string(body))
		if len(value) >= 4 {
			return []string{value}
		}
		return nil
	}
	var values []string
	allowed := map[string]bool{"accepted": true, "ingest_status": true}
	for key, value := range raw {
		if allowed[key] {
			continue
		}
		if text, ok := value.(string); ok && len(strings.TrimSpace(text)) >= 4 {
			values = append(values, text)
		}
	}
	return values
}

func durationMS(started time.Time) int64 {
	return time.Since(started).Milliseconds()
}

func renderSummaryMarkdown(summary Summary) string {
	return fmt.Sprintf(`# AVL adapter send summary

- mode: send
- dry_run: false
- telemetry_url_path: %s
- telemetry_target_loopback: %t
- telemetry_host_ref: %s
- output_label: %s
- output_ref: %s
- transformed_count: %d
- sent_count: %d
- succeeded_count: %d
- failed_count: %d
- skipped_count: %d
- retry_total: %d
- external_evidence_created: false
- consumer_statuses_changed: false
- compliance_claimed: false
- production_readiness_claimed: false
- hosted_saas_claimed: false
- agency_adoption_claimed: false
- consumer_acceptance_claimed: false
- vendor_compatibility_claimed: false
- production_avl_reliability_claimed: false
- production_grade_eta_claimed: false
`, summary.TelemetryURLPath, summary.TelemetryTargetLoopback, summary.TelemetryHostRef, summary.OutputLabel, summary.OutputRef, summary.TransformedCount, summary.SentCount, summary.SucceededCount, summary.FailedCount, summary.SkippedCount, summary.RetryTotal)
}

func renderManifestMarkdown(manifest RedactedManifest) string {
	var b strings.Builder
	b.WriteString("# AVL adapter send manifest\n\n")
	b.WriteString("- schema_version: " + manifest.SchemaVersion + "\n")
	b.WriteString("- telemetry_url_env: " + manifest.TelemetryURLEnv + "\n")
	b.WriteString("- credentials:\n")
	for _, row := range manifest.Credentials {
		b.WriteString("  - credential_ref: " + row.CredentialRef + "\n")
		b.WriteString("    token_env: " + row.TokenEnv + "\n")
	}
	return b.String()
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func shortHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])[:12]
}

func redactedOutputRef(rel string) string {
	slashRel := filepath.ToSlash(rel)
	if slashRel == ".cache" || strings.HasPrefix(slashRel, ".cache/") {
		return slashRel
	}
	return "output:" + shortHash(slashRel)
}

func hasSymlinkAncestor(abs string) bool {
	abs = filepath.Clean(abs)
	volume := filepath.VolumeName(abs)
	rest := strings.TrimPrefix(abs, volume)
	parts := strings.Split(strings.Trim(rest, string(filepath.Separator)), string(filepath.Separator))
	current := volume + string(filepath.Separator)
	for _, part := range parts {
		if part == "" {
			continue
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		if err != nil {
			return true
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return true
		}
	}
	return false
}

func secretLike(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "authorization:") ||
		strings.Contains(lower, "bearer ") ||
		strings.Contains(lower, "postgres://") ||
		strings.Contains(lower, "postgresql://") ||
		strings.Contains(lower, "-----begin") && strings.Contains(lower, "private key-----") ||
		strings.Contains(lower, "cookie:") ||
		strings.Contains(lower, "token=") ||
		strings.Contains(lower, "token:") ||
		strings.Contains(lower, "password=") ||
		strings.Contains(lower, "password:") ||
		strings.Contains(lower, "secret=") ||
		strings.Contains(lower, "secret:")
}

func parseDurationEnv(getenv Environment, name string, target *time.Duration, diagnostics *[]Diagnostic) {
	raw := strings.TrimSpace(getenv(name))
	if raw == "" {
		return
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		*diagnostics = append(*diagnostics, errorDiagnostic(CodeInvalidSendConfig, name+" must be a Go duration", nil))
		return
	}
	*target = value
}

func parseIntEnv(getenv Environment, name string, target *int, diagnostics *[]Diagnostic) {
	raw := strings.TrimSpace(getenv(name))
	if raw == "" {
		return
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		*diagnostics = append(*diagnostics, errorDiagnostic(CodeInvalidSendConfig, name+" must be an integer", nil))
		return
	}
	*target = value
}

func parseBoolEnv(getenv Environment, name string, target *bool, diagnostics *[]Diagnostic) {
	raw := strings.TrimSpace(getenv(name))
	if raw == "" {
		return
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		*diagnostics = append(*diagnostics, errorDiagnostic(CodeInvalidSendConfig, name+" must be boolean", nil))
		return
	}
	*target = value
}

func parseTimeEnv(getenv Environment, name string, target *time.Time, diagnostics *[]Diagnostic) {
	raw := strings.TrimSpace(getenv(name))
	if raw == "" {
		return
	}
	value, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		*diagnostics = append(*diagnostics, errorDiagnostic(CodeInvalidSendConfig, name+" must be RFC3339", nil))
		return
	}
	*target = value
}
