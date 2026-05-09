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
	outDir := filepath.Join(t.TempDir(), "diagnostics")
	var stdout, stderr strings.Builder

	err := run(t.Context(), []string{
		"--scenario", filepath.Join("..", "..", "testdata", "telemetry-simulator", "on-route.json"),
		"--target", "https://reference.example.test",
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
}

func TestRunPostsThroughAuthenticatedTelemetryEndpoint(t *testing.T) {
	t.Setenv("DEVICE_TOKEN", "")
	outDir := filepath.Join(t.TempDir(), "diagnostics")
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
}

func TestPrepareOutputDirRejectsEvidenceDirectory(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(root, "..", ".."))
	target := filepath.Join(repoRoot, "docs", "evidence", "telemetry-simulator")
	if _, err := prepareOutputDir(target, false); err == nil {
		t.Fatal("expected docs/evidence output directory to be rejected")
	}
}
