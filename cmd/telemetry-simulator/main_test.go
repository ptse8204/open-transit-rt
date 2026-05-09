package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDryRunWritesPrivateDiagnosticsWithoutToken(t *testing.T) {
	t.Setenv("DEVICE_TOKEN", "")
	outDir := cacheTestDir(t)
	var stdout, stderr strings.Builder

	err := run(t.Context(), []string{
		"--scenario", filepath.Join("..", "..", "testdata", "telemetry-simulator", "on-route.json"),
		"--target", "http://reference.example.test",
		"--device-token", "",
		"--dry-run",
		"--output-dir", outDir,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run dry-run: %v\nstderr=%s", err, stderr.String())
	}

	payload, err := os.ReadFile(filepath.Join(outDir, "summary.json"))
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	if strings.Contains(string(payload), defaultDeviceToken) || strings.Contains(string(payload), "Authorization") {
		t.Fatalf("diagnostics leaked token material or auth header: %s", payload)
	}
	var summary runSummary
	if err := json.Unmarshal(payload, &summary); err != nil {
		t.Fatalf("parse summary: %v", err)
	}
	if !summary.SyntheticOnly || summary.ExternalEvidenceCreated || summary.ConsumerStatusesChanged {
		t.Fatalf("unexpected evidence boundary fields: %+v", summary)
	}
	if summary.EventsSent != 0 || !summary.DryRun {
		t.Fatalf("dry run summary = %+v", summary)
	}
	if summary.DurationMS <= 0 {
		t.Fatalf("dry run duration_ms = %d, want positive", summary.DurationMS)
	}
	if len(summary.Events) != 1 || summary.Events[0].HTTPDurationMS != 0 || summary.Events[0].MatcherDurationMS != 0 {
		t.Fatalf("dry run event timings = %+v", summary.Events)
	}
}

func TestRunPostsThroughAuthenticatedTelemetryEndpoint(t *testing.T) {
	t.Setenv("DEVICE_TOKEN", "")
	outDir := cacheTestDir(t)
	const token = "synthetic-test-token"
	seenAuth := false
	seenPath := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path == "/v1/telemetry"
		seenAuth = r.Header.Get("Authorization") == "Bearer "+token
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"accepted":true,"ingest_status":"accepted","agency_id":"demo-agency","vehicle_id":"bus-1","observed_at":"2026-04-20T15:05:00Z","received_at":"2026-04-20T15:05:01Z"}`))
	}))
	defer server.Close()

	var stdout, stderr strings.Builder
	err := run(t.Context(), []string{
		"--scenario", filepath.Join("..", "..", "testdata", "telemetry-simulator", "on-route.json"),
		"--target", server.URL,
		"--device-token", token,
		"--output-dir", outDir,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v\nstderr=%s", err, stderr.String())
	}
	if !seenPath || !seenAuth {
		t.Fatalf("request did not use authenticated /v1/telemetry path: path=%v auth=%v", seenPath, seenAuth)
	}
	payload, err := os.ReadFile(filepath.Join(outDir, "events.json"))
	if err != nil {
		t.Fatalf("read events: %v", err)
	}
	if strings.Contains(string(payload), token) || strings.Contains(string(payload), "Authorization") {
		t.Fatalf("events diagnostics leaked token material or auth header: %s", payload)
	}
	var events []eventResult
	if err := json.Unmarshal(payload, &events); err != nil {
		t.Fatalf("parse events: %v", err)
	}
	if len(events) != 1 || events[0].HTTPDurationMS <= 0 {
		t.Fatalf("http timing not recorded: %+v", events)
	}
	summaryPayload, err := os.ReadFile(filepath.Join(outDir, "summary.json"))
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	var summary runSummary
	if err := json.Unmarshal(summaryPayload, &summary); err != nil {
		t.Fatalf("parse summary: %v", err)
	}
	if summary.DurationMS <= 0 {
		t.Fatalf("summary duration_ms = %d, want positive", summary.DurationMS)
	}
}

func TestRunBatchScenarioCountsAndTimings(t *testing.T) {
	t.Setenv("DEVICE_TOKEN", "")
	baseDir := cacheTestDir(t)
	scenarioPath := filepath.Join(baseDir, "batch.json")
	if err := os.MkdirAll(filepath.Dir(scenarioPath), 0o700); err != nil {
		t.Fatalf("mkdir scenario dir: %v", err)
	}
	if err := os.WriteFile(scenarioPath, []byte(`{
  "name": "batch",
  "description": "synthetic batch test",
  "synthetic_only": true,
  "reference_time": "2026-04-20T15:05:00Z",
  "events": [
    {
      "label": "accepted",
      "offset_seconds": 0,
      "payload": {"agency_id":"demo-agency","device_id":"device-1","vehicle_id":"bus-1","timestamp":"2026-04-20T15:05:00Z","lat":49.27935,"lon":-123.11785},
      "expected_http_statuses": [201],
      "expected_ingest_statuses": ["accepted"]
    },
    {
      "label": "duplicate",
      "offset_seconds": 0,
      "payload": {"agency_id":"demo-agency","device_id":"device-1","vehicle_id":"bus-1","timestamp":"2026-04-20T15:05:00Z","lat":49.27935,"lon":-123.11785},
      "expected_http_statuses": [202],
      "expected_ingest_statuses": ["duplicate"]
    },
    {
      "label": "rejected",
      "offset_seconds": 30,
      "payload": {"agency_id":"demo-agency","device_id":"device-1","vehicle_id":"bus-1","timestamp":"2026-04-20T15:05:30Z","lat":49.27935,"lon":-123.11785},
      "expected_http_statuses": [401]
    }
  ]
}`), 0o600); err != nil {
		t.Fatalf("write scenario: %v", err)
	}
	responseByCount := []struct {
		status int
		body   string
	}{
		{http.StatusCreated, `{"accepted":true,"ingest_status":"accepted","agency_id":"demo-agency","vehicle_id":"bus-1","observed_at":"2026-04-20T15:05:00Z","received_at":"2026-04-20T15:05:01Z"}`},
		{http.StatusAccepted, `{"accepted":false,"ingest_status":"duplicate","agency_id":"demo-agency","vehicle_id":"bus-1","observed_at":"2026-04-20T15:05:00Z","received_at":"2026-04-20T15:05:02Z"}`},
		{http.StatusUnauthorized, `{"error":"invalid device credential"}`},
	}
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if count >= len(responseByCount) {
			t.Fatalf("unexpected request %d", count+1)
		}
		response := responseByCount[count]
		count++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.status)
		_, _ = w.Write([]byte(response.body))
	}))
	defer server.Close()

	outDir := filepath.Join(baseDir, "diagnostics")
	var stdout, stderr strings.Builder
	err := run(t.Context(), []string{
		"--scenario", scenarioPath,
		"--target", server.URL,
		"--device-token", "synthetic-test-token",
		"--output-dir", outDir,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run batch: %v\nstderr=%s", err, stderr.String())
	}
	if count != 3 {
		t.Fatalf("request count = %d, want 3", count)
	}
	payload, err := os.ReadFile(filepath.Join(outDir, "summary.json"))
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	var summary runSummary
	if err := json.Unmarshal(payload, &summary); err != nil {
		t.Fatalf("parse summary: %v", err)
	}
	if summary.EventsSent != 3 || summary.EventsAccepted != 1 || summary.EventsDuplicate != 1 || summary.EventsRejected != 1 {
		t.Fatalf("unexpected counts: %+v", summary)
	}
	if len(summary.Events) != 3 {
		t.Fatalf("events length = %d, want 3", len(summary.Events))
	}
	for _, event := range summary.Events {
		if event.HTTPDurationMS <= 0 {
			t.Fatalf("event %q http_duration_ms = %d, want positive", event.Label, event.HTTPDurationMS)
		}
	}
}

func TestValidateTargetAndTokenURLSafety(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config
		wantErr bool
	}{
		{
			name: "loopback http allowed",
			cfg:  config{target: "http://127.0.0.1:8080", deviceToken: "token"},
		},
		{
			name: "localhost http allowed",
			cfg:  config{target: "http://localhost:8080", deviceToken: "token"},
		},
		{
			name:    "non-loopback http rejected",
			cfg:     config{target: "http://example.com", deviceToken: "token"},
			wantErr: true,
		},
		{
			name: "non-loopback http dry-run allowed",
			cfg:  config{target: "http://example.com", dryRun: true},
		},
		{
			name: "https allowed",
			cfg:  config{target: "https://example.com", deviceToken: "token"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTargetAndToken(tt.cfg)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestPrepareOutputDirSafety(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(root, "..", ".."))
	cacheDir := filepath.Join(repoRoot, ".cache", "telemetry-simulator-tests", strings.ReplaceAll(t.Name(), "/", "-"), "allowed")
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(repoRoot, ".cache", "telemetry-simulator-tests", strings.ReplaceAll(t.Name(), "/", "-")))
	})
	allowed, err := prepareOutputDir(cacheDir, false, false)
	if err != nil {
		t.Fatalf(".cache output dir rejected: %v", err)
	}
	if info, err := os.Stat(allowed); err != nil {
		t.Fatalf("stat allowed dir: %v", err)
	} else if info.Mode().Perm() != 0o700 {
		t.Fatalf("allowed dir mode = %o, want 700", info.Mode().Perm())
	}

	if _, err := prepareOutputDir(filepath.Join(repoRoot, "docs", "tmp", "telemetry-simulator"), false, false); err == nil {
		t.Fatal("expected docs/tmp output directory to be rejected without override")
	}
	if _, err := prepareOutputDir(filepath.Join(repoRoot, "docs", "evidence", "telemetry-simulator"), false, true); err == nil {
		t.Fatal("expected docs/evidence output directory to be rejected even with override")
	}

	link := filepath.Join(repoRoot, ".cache", "telemetry-simulator-tests", strings.ReplaceAll(t.Name(), "/", "-"), "link")
	target := filepath.Join(repoRoot, ".cache", "telemetry-simulator-tests", strings.ReplaceAll(t.Name(), "/", "-"), "target")
	if err := os.MkdirAll(target, 0o700); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := prepareOutputDir(link, false, false); err == nil {
		t.Fatal("expected symlink output directory to be rejected")
	}
}

func cacheTestDir(t *testing.T) string {
	t.Helper()
	root, err := repoRoot()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, ".cache", "telemetry-simulator-tests", strings.ReplaceAll(t.Name(), "/", "-"))
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("remove stale cache test dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return dir
}

func TestRedactionScanRejectsSecrets(t *testing.T) {
	cases := map[string]string{
		"authorization": "Authorization: Bearer secret",
		"bearer":        "Bearer secret",
		"device_token":  "DEVICE_TOKEN=secret",
		"database_url":  "DATABASE_URL=postgres://user:pass@example/db",
		"postgres_url":  "postgres://user:pass@example/db",
		"cookie":        "Cookie: session=secret",
		"private_key":   "-----BEGIN PRIVATE KEY-----",
	}
	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join(cacheTestDir(t), name)
			if err := os.MkdirAll(dir, 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "bad.txt"), []byte(value), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := scanOutputDirRedaction(dir); err == nil {
				t.Fatal("expected redaction scan to reject secret-shaped output")
			}
		})
	}
}
