package avladapter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestSendPostsTelemetryWithBearerAuthAndJSONContentType(t *testing.T) {
	result := transformFixture(t, "mapping.json", "valid.json")
	prepared, diagnostics, requests := prepareSendForTest(t, result, "valid.json", http.StatusCreated)
	if HasHardErrors(diagnostics) {
		t.Fatalf("preflight diagnostics: %+v", diagnostics)
	}
	report := SendEvents(context.Background(), result.Events, prepared, requests.server.Client(), noSleep)
	if report.Summary.SucceededCount != 1 || report.Summary.SentCount != 1 {
		t.Fatalf("summary = %+v", report.Summary)
	}
	if requests.count != 1 {
		t.Fatalf("request count = %d, want 1", requests.count)
	}
	if requests.path != "/v1/telemetry" {
		t.Fatalf("path = %q", requests.path)
	}
	if requests.auth != "Bearer synthetic-device-token" {
		t.Fatalf("authorization header was not set correctly")
	}
	if !strings.HasPrefix(requests.contentType, "application/json") {
		t.Fatalf("content type = %q", requests.contentType)
	}
	if !strings.Contains(requests.body, `"agency_id":"demo-agency"`) {
		t.Fatalf("request body is not telemetry JSON: %s", requests.body)
	}
}

func TestStaleAndFutureBlockSendWithZeroRequests(t *testing.T) {
	for _, tc := range []struct {
		name    string
		fixture string
	}{
		{name: "stale", fixture: "stale-timestamp.json"},
		{name: "future", fixture: "future-timestamp.json"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := transformFixture(t, "mapping.json", tc.fixture)
			prepared, diagnostics, requests := prepareSendForTest(t, result, tc.fixture, http.StatusCreated)
			_ = prepared
			if !HasHardErrors(diagnostics) {
				t.Fatalf("diagnostics = %+v, want send blocker", diagnostics)
			}
			if requests.count != 0 {
				t.Fatalf("request count = %d, want 0 before preflight passes", requests.count)
			}
		})
	}
}

func TestWarningsSendByDefaultAndFailOnWarningsBlocks(t *testing.T) {
	result := transformFixture(t, "mapping.json", "low-gps-accuracy.json")
	prepared, diagnostics, requests := prepareSendForTest(t, result, "low-gps-accuracy.json", http.StatusAccepted)
	if HasHardErrors(diagnostics) {
		t.Fatalf("preflight diagnostics: %+v", diagnostics)
	}
	report := SendEvents(context.Background(), result.Events, prepared, requests.server.Client(), noSleep)
	if report.Summary.SucceededCount != 1 || requests.count != 1 {
		t.Fatalf("summary=%+v requests=%d", report.Summary, requests.count)
	}

	prepared.Config.FailOnWarnings = true
	_, diagnostics = PrepareSend(prepared.Manifest, prepared.Config, result, mustPayload(t, "low-gps-accuracy.json"), envForTest(requests.server.URL, prepared.Output.Ref), mustCwd(t))
	if !HasHardErrors(diagnostics) {
		t.Fatalf("diagnostics = %+v, want fail-on-warnings blocker", diagnostics)
	}
}

func TestSendOutputFileSetAndRedactedFields(t *testing.T) {
	result := transformFixture(t, "mapping.json", "valid.json")
	prepared, diagnostics, requests := prepareSendForTest(t, result, "valid.json", http.StatusCreated)
	if HasHardErrors(diagnostics) {
		t.Fatalf("preflight diagnostics: %+v", diagnostics)
	}
	report := SendEvents(context.Background(), result.Events, prepared, requests.server.Client(), noSleep)
	report, diagnostics = BuildSendFiles(report, prepared)
	if HasHardErrors(diagnostics) {
		t.Fatalf("build diagnostics: %+v", diagnostics)
	}
	if err := WriteSendFiles(report, prepared.Output); err != nil {
		t.Fatalf("write files: %v", err)
	}
	entries, err := os.ReadDir(prepared.Output.Dir)
	if err != nil {
		t.Fatalf("read output dir: %v", err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	want := []string{"diagnostics.json", "manifest.json", "manifest.md", "summary.json", "summary.md"}
	if strings.Join(names, ",") != strings.Join(want, ",") {
		t.Fatalf("files = %v, want %v", names, want)
	}
	summaryRaw := string(report.Files["summary.json"])
	if !strings.Contains(summaryRaw, `"telemetry_url_path": "/v1/telemetry"`) ||
		!strings.Contains(summaryRaw, `"external_evidence_created": false`) ||
		!strings.Contains(summaryRaw, `"consumer_statuses_changed": false`) {
		t.Fatalf("summary missing required fields: %s", summaryRaw)
	}
	if strings.Contains(summaryRaw, requests.server.Listener.Addr().String()) || strings.Contains(summaryRaw, prepared.Output.Dir) {
		t.Fatalf("summary leaked raw host or absolute output path: %s", summaryRaw)
	}
	diagnosticsRaw := string(report.Files["diagnostics.json"])
	for _, forbidden := range []string{"demo-agency", "device-1", "bus-1", "vendor-device-1", "vendor-vehicle-1"} {
		if strings.Contains(diagnosticsRaw, forbidden) {
			t.Fatalf("diagnostics leaked raw id %q: %s", forbidden, diagnosticsRaw)
		}
	}
	if !strings.Contains(diagnosticsRaw, `"record_index"`) || !strings.Contains(diagnosticsRaw, `"credential_ref"`) {
		t.Fatalf("diagnostics missing redacted references: %s", diagnosticsRaw)
	}
}

func TestOutputRefKeepsCachePathAndRedactsNonCacheOverride(t *testing.T) {
	cwd := mustCwd(t)
	cacheOutput, diagnostics := ResolveOutputTarget(filepath.Join(".cache", "avl-vendor-adapter", "phase-48-test"), referenceTime, cwd)
	if HasHardErrors(diagnostics) {
		t.Fatalf("cache output diagnostics: %+v", diagnostics)
	}
	if cacheOutput.Ref != ".cache/avl-vendor-adapter/phase-48-test" {
		t.Fatalf("cache output ref = %q", cacheOutput.Ref)
	}

	privateRel := filepath.Join("private", "operator-acme", "out")
	privateOutput, diagnostics := ResolveOutputTarget(privateRel, referenceTime, cwd)
	if HasHardErrors(diagnostics) {
		t.Fatalf("private output diagnostics: %+v", diagnostics)
	}
	if strings.Contains(privateOutput.Ref, "private") || strings.Contains(privateOutput.Ref, "operator-acme") {
		t.Fatalf("private output ref leaked raw path: %q", privateOutput.Ref)
	}
	if !strings.HasPrefix(privateOutput.Ref, "output:") {
		t.Fatalf("private output ref = %q, want redacted output ref", privateOutput.Ref)
	}

	result := transformFixture(t, "mapping.json", "valid.json")
	requests := newRequestRecorder(http.StatusCreated)
	defer requests.server.Close()
	prepared, diagnostics := preparedForServer(t, result, "valid.json", requests.server.URL, privateRel, false)
	if HasHardErrors(diagnostics) {
		t.Fatalf("preflight diagnostics: %+v", diagnostics)
	}
	report := SendEvents(context.Background(), result.Events, prepared, requests.server.Client(), noSleep)
	report, diagnostics = BuildSendFiles(report, prepared)
	if HasHardErrors(diagnostics) {
		t.Fatalf("build diagnostics: %+v", diagnostics)
	}
	summary := string(report.Files["summary.json"])
	if strings.Contains(summary, filepath.ToSlash(privateRel)) || strings.Contains(report.Stdout, filepath.ToSlash(privateRel)) {
		t.Fatalf("non-cache output leaked raw relative path: summary=%s stdout=%s", summary, report.Stdout)
	}
	if !strings.Contains(summary, `"output_ref": "output:`) || !strings.Contains(report.Stdout, `"output_ref":"output:`) {
		t.Fatalf("redacted output ref missing: summary=%s stdout=%s", summary, report.Stdout)
	}
}

func TestSendManifestValidationFailures(t *testing.T) {
	cases := []string{
		`{"schema_version":"wrong","telemetry_url_env":"AVL_ADAPTER_TELEMETRY_URL","credentials":[]}`,
		`{"schema_version":"avl-adapter-send.v1","telemetry_url_env":"bad-env","credentials":[]}`,
		`{"schema_version":"avl-adapter-send.v1","telemetry_url_env":"AVL_ADAPTER_TELEMETRY_URL","credentials":[{"agency_id":"demo-agency","device_id":"device-1","vehicle_id":"bus-1","token_env":"AVL_ADAPTER_DEVICE_TOKEN"},{"agency_id":"demo-agency","device_id":"device-1","vehicle_id":"bus-1","token_env":"AVL_ADAPTER_DEVICE_TOKEN_2"}]}`,
		`{"schema_version":"avl-adapter-send.v1","telemetry_url_env":"AVL_ADAPTER_TELEMETRY_URL","credentials":[{"agency_id":"demo-agency","device_id":"device-1","vehicle_id":"bus-1","token_env":"Bearer raw-token"}]}`,
		`{"schema_version":"avl-adapter-send.v1","telemetry_url_env":"AVL_ADAPTER_TELEMETRY_URL","credentials":[{"agency_id":"demo-agency","device_id":"device-1","vehicle_id":"bus-1","token_env":"AVL_ADAPTER_DEVICE_TOKEN","notes":"token=raw"}]}`,
		`{"schema_version":"avl-adapter-send.v1","telemetry_url_env":"AVL_ADAPTER_TELEMETRY_URL","credentials":[{"agency_id":"demo-agency","device_id":"device-1","vehicle_id":"bus-1","token_env":"AVL_ADAPTER_DEVICE_TOKEN","notes":"-----BEGIN PRIVATE KEY-----"}]}`,
		`{"schema_version":"avl-adapter-send.v1","telemetry_url_env":"AVL_ADAPTER_TELEMETRY_URL","extra":true,"credentials":[]}`,
		`{"schema_version":"avl-adapter-send.v1","telemetry_url_env":"AVL_ADAPTER_TELEMETRY_URL","credentials":[{"agency_id":"demo-agency","device_id":"device-1","vehicle_id":"bus-1","token_env":"AVL_ADAPTER_DEVICE_TOKEN","extra":true}]}`,
	}
	for _, raw := range cases {
		_, diagnostics := LoadSendManifest(strings.NewReader(raw))
		if !HasHardErrors(diagnostics) {
			t.Fatalf("manifest %s diagnostics = %+v, want hard error", raw, diagnostics)
		}
	}

	result := transformFixture(t, "mapping.json", "valid.json")
	manifest := validSendManifest()
	manifest.Credentials = nil
	cfg := validSendConfig("http://127.0.0.1:1/v1/telemetry", safeOutputPath(t))
	_, diagnostics := PrepareSend(manifest, cfg, result, mustPayload(t, "valid.json"), envForTest("http://127.0.0.1:1", cfg.OutputDir), mustCwd(t))
	assertDiagnostic(t, diagnostics, CodeMissingCredentialMapping, SeverityError)
}

func TestTelemetryTargetValidationRejectsUnsafeTargets(t *testing.T) {
	valid, diagnostics := ValidateTelemetryTarget("http://127.0.0.1:8080/v1/telemetry")
	if HasHardErrors(diagnostics) || !valid.Loopback || valid.Path != "/v1/telemetry" {
		t.Fatalf("loopback target diagnostics=%+v target=%+v", diagnostics, valid)
	}
	for _, raw := range []string{
		"http://user@example.test/v1/telemetry",
		"https://example.test/v1/telemetry?x=1",
		"https://example.test/v1/telemetry#fragment",
		"https://example.test/v1/events",
		"http://example.test/v1/telemetry",
	} {
		_, diagnostics := ValidateTelemetryTarget(raw)
		if !HasHardErrors(diagnostics) {
			t.Fatalf("target %q diagnostics=%+v, want hard error", raw, diagnostics)
		}
	}
}

func TestStatusSuccessAndRetryBehavior(t *testing.T) {
	for _, status := range []int{http.StatusCreated, http.StatusAccepted} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			result := transformFixture(t, "mapping.json", "valid.json")
			prepared, diagnostics, requests := prepareSendForTest(t, result, "valid.json", status)
			if HasHardErrors(diagnostics) {
				t.Fatalf("preflight diagnostics: %+v", diagnostics)
			}
			report := SendEvents(context.Background(), result.Events, prepared, requests.server.Client(), noSleep)
			if report.Summary.SucceededCount != 1 {
				t.Fatalf("summary = %+v", report.Summary)
			}
		})
	}

	for _, status := range []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound, http.StatusMethodNotAllowed} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			result := transformFixture(t, "mapping.json", "valid.json")
			prepared, diagnostics, requests := prepareSendForTest(t, result, "valid.json", status)
			if HasHardErrors(diagnostics) {
				t.Fatalf("preflight diagnostics: %+v", diagnostics)
			}
			report := SendEvents(context.Background(), result.Events, prepared, requests.server.Client(), noSleep)
			if report.Diagnostics[0].Attempts != 1 || requests.count != 1 {
				t.Fatalf("diagnostics=%+v requests=%d, want no retry", report.Diagnostics, requests.count)
			}
		})
	}
}

func TestRetryableFailuresUseSleeper(t *testing.T) {
	result := transformFixture(t, "mapping.json", "valid.json")
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 1 {
			http.Error(w, "retry later", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"accepted":true,"ingest_status":"accepted"}`))
	}))
	defer server.Close()
	prepared, diagnostics := preparedForServer(t, result, "valid.json", server.URL, safeOutputPath(t), false)
	if HasHardErrors(diagnostics) {
		t.Fatalf("preflight diagnostics: %+v", diagnostics)
	}
	sleeps := 0
	report := SendEvents(context.Background(), result.Events, prepared, server.Client(), func(context.Context, time.Duration) error {
		sleeps++
		return nil
	})
	if report.Summary.SucceededCount != 1 || report.Summary.RetryTotal != 1 || sleeps != 1 || requests != 2 {
		t.Fatalf("summary=%+v sleeps=%d requests=%d", report.Summary, sleeps, requests)
	}
}

func TestTerminalFailureStopsLaterRecords(t *testing.T) {
	result := transformFixture(t, "multi-vehicle-mapping.json", "multi-vehicle-gps.json")
	manifest := validSendManifest()
	manifest.Credentials = append(manifest.Credentials, ManifestCredential{
		AgencyID:  "demo-agency",
		DeviceID:  "device-2",
		VehicleID: "bus-2",
		TokenEnv:  "AVL_ADAPTER_DEVICE_TOKEN_2",
	})
	requests := newRequestRecorder(http.StatusBadRequest)
	defer requests.server.Close()
	prepared, diagnostics := preparedForManifestAndServer(t, result, "multi-vehicle-gps.json", manifest, requests.server.URL, safeOutputPath(t), false)
	if HasHardErrors(diagnostics) {
		t.Fatalf("preflight diagnostics: %+v", diagnostics)
	}
	report := SendEvents(context.Background(), result.Events, prepared, requests.server.Client(), noSleep)
	if report.Summary.FailedCount != 1 || report.Summary.SkippedCount != 1 || requests.count != 1 {
		t.Fatalf("summary=%+v requests=%d", report.Summary, requests.count)
	}
	if report.Diagnostics[1].Outcome != "skipped_after_failure" {
		t.Fatalf("diagnostics = %+v", report.Diagnostics)
	}
}

func TestNoRawResponseBodyLeakage(t *testing.T) {
	result := transformFixture(t, "mapping.json", "valid.json")
	requests := newRequestRecorderWithBody(http.StatusCreated, `{"accepted":true,"ingest_status":"accepted","message":"PRIVATE_BODY_MARKER"}`)
	defer requests.server.Close()
	prepared, diagnostics := preparedForServer(t, result, "valid.json", requests.server.URL, safeOutputPath(t), false)
	if HasHardErrors(diagnostics) {
		t.Fatalf("preflight diagnostics: %+v", diagnostics)
	}
	report := SendEvents(context.Background(), result.Events, prepared, requests.server.Client(), noSleep)
	report, diagnostics = BuildSendFiles(report, prepared)
	if HasHardErrors(diagnostics) {
		t.Fatalf("build diagnostics: %+v", diagnostics)
	}
	combined := report.Stdout
	for _, raw := range report.Files {
		combined += string(raw)
	}
	if strings.Contains(combined, "PRIVATE_BODY_MARKER") {
		t.Fatalf("raw response body marker leaked: %s", combined)
	}
}

func TestRedactionScanRejectsForbiddenValues(t *testing.T) {
	for _, content := range []string{
		"Authorization: Bearer value",
		"Bearer raw-token",
		"Cookie: session=value",
		"postgres://user:pass@localhost/db",
		"-----BEGIN PRIVATE KEY-----",
		"contains synthetic-device-token",
		"contains vendor-device-1",
	} {
		if hits := ScanRedaction(content, []string{"synthetic-device-token", "vendor-device-1"}); len(hits) == 0 {
			t.Fatalf("content %q passed redaction scan", content)
		}
	}
}

func TestOutputRejectsSymlinksAndDocsEvidence(t *testing.T) {
	cwd := mustCwd(t)
	_, diagnostics := ResolveOutputTarget(filepath.Join("docs", "evidence", "phase-48"), referenceTime, cwd)
	assertDiagnostic(t, diagnostics, CodeInvalidOutputPath, SeverityError)

	link := filepath.Join(cwd, ".cache", "avl-vendor-symlink-test")
	target := filepath.Join(cwd, ".cache", "avl-vendor-symlink-target")
	_ = os.Remove(link)
	_ = os.RemoveAll(target)
	t.Cleanup(func() {
		_ = os.Remove(link)
		_ = os.RemoveAll(target)
	})
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("symlink: %v", err)
	}
	_, diagnostics = ResolveOutputTarget(filepath.Join(".cache", "avl-vendor-symlink-test", "out"), referenceTime, cwd)
	assertDiagnostic(t, diagnostics, CodeInvalidOutputPath, SeverityError)
}

type requestRecorder struct {
	server      *httptest.Server
	count       int
	path        string
	auth        string
	contentType string
	body        string
}

func newRequestRecorder(status int) *requestRecorder {
	return newRequestRecorderWithBody(status, `{"accepted":true,"ingest_status":"accepted"}`)
}

func newRequestRecorderWithBody(status int, body string) *requestRecorder {
	recorder := &requestRecorder{}
	recorder.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder.count++
		recorder.path = r.URL.Path
		recorder.auth = r.Header.Get("Authorization")
		recorder.contentType = r.Header.Get("Content-Type")
		raw, _ := io.ReadAll(r.Body)
		recorder.body = string(raw)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	return recorder
}

func prepareSendForTest(t *testing.T, result Result, payloadName string, status int) (SendPrepared, []Diagnostic, *requestRecorder) {
	t.Helper()
	requests := newRequestRecorder(status)
	t.Cleanup(requests.server.Close)
	prepared, diagnostics := preparedForServer(t, result, payloadName, requests.server.URL, safeOutputPath(t), false)
	return prepared, diagnostics, requests
}

func preparedForServer(t *testing.T, result Result, payloadName string, serverURL string, outputPath string, failOnWarnings bool) (SendPrepared, []Diagnostic) {
	t.Helper()
	return preparedForManifestAndServer(t, result, payloadName, validSendManifest(), serverURL, outputPath, failOnWarnings)
}

func preparedForManifestAndServer(t *testing.T, result Result, payloadName string, manifest SendManifest, serverURL string, outputPath string, failOnWarnings bool) (SendPrepared, []Diagnostic) {
	t.Helper()
	cfg := validSendConfig(serverURL+"/v1/telemetry", outputPath)
	cfg.FailOnWarnings = failOnWarnings
	return PrepareSend(manifest, cfg, result, mustPayload(t, payloadName), envForTest(serverURL, outputPath), mustCwd(t))
}

func validSendManifest() SendManifest {
	return SendManifest{
		SchemaVersion:   SendManifestSchemaVersion,
		TelemetryURLEnv: EnvTelemetryURL,
		Credentials: []ManifestCredential{{
			AgencyID:  "demo-agency",
			DeviceID:  "device-1",
			VehicleID: "bus-1",
			TokenEnv:  "AVL_ADAPTER_DEVICE_TOKEN",
			Notes:     "Synthetic local reference binding.",
		}},
	}
}

func validSendConfig(targetURL string, outputPath string) SendConfig {
	return SendConfig{
		TelemetryURL:      targetURL,
		OutputDir:         outputPath,
		Timeout:           DefaultSendTimeout,
		MaxRetries:        DefaultSendMaxRetries,
		RetryInitialDelay: time.Millisecond,
		RetryMaxDelay:     time.Millisecond,
		ReferenceTime:     referenceTime,
		StaleThreshold:    DefaultStaleThreshold,
		FutureThreshold:   DefaultFutureThreshold,
		GeneratedAt:       referenceTime,
	}
}

func envForTest(serverURL string, outputPath string) Environment {
	return func(name string) string {
		switch name {
		case EnvTelemetryURL:
			return serverURL + "/v1/telemetry"
		case EnvOutputDir:
			return outputPath
		case "AVL_ADAPTER_DEVICE_TOKEN", "AVL_ADAPTER_DEVICE_TOKEN_2":
			return "synthetic-device-token"
		default:
			return ""
		}
	}
}

func mustPayload(t *testing.T, name string) Payload {
	t.Helper()
	payload, diagnostics := DecodePayload(strings.NewReader(string(readFixture(t, name))))
	if HasHardErrors(diagnostics) {
		t.Fatalf("payload diagnostics: %+v", diagnostics)
	}
	return payload
}

func safeOutputPath(t *testing.T) string {
	t.Helper()
	path := filepath.Join(".cache", "avl-vendor-adapter-test-"+strings.ReplaceAll(t.Name(), "/", "-")+"-"+time.Now().Format("150405.000000000"))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(mustCwd(t), path)) })
	return path
}

func mustCwd(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(cwd, "..", ".."))
}

func noSleep(context.Context, time.Duration) error { return nil }

func TestSendReportJSONDoesNotContainRawIDs(t *testing.T) {
	result := transformFixture(t, "mapping.json", "valid.json")
	prepared, diagnostics, requests := prepareSendForTest(t, result, "valid.json", http.StatusCreated)
	if HasHardErrors(diagnostics) {
		t.Fatalf("preflight diagnostics: %+v", diagnostics)
	}
	report := SendEvents(context.Background(), result.Events, prepared, requests.server.Client(), noSleep)
	raw, err := json.Marshal(report.Diagnostics)
	if err != nil {
		t.Fatalf("marshal diagnostics: %v", err)
	}
	for _, forbidden := range []string{"demo-agency", "device-1", "bus-1"} {
		if strings.Contains(string(raw), forbidden) {
			t.Fatalf("diagnostics leaked %q: %s", forbidden, raw)
		}
	}
}
