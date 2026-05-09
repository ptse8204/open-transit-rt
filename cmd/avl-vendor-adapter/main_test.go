package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"open-transit-rt/internal/avladapter"
	"open-transit-rt/internal/telemetry"
)

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := run([]string{"help"}, &stdout, &stderr); err != nil {
		t.Fatalf("run help: %v", err)
	}
	assertHelpBoundary(t, stdout.String())
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("help output missing usage: %s", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunRequiresDryRun(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"--mapping", "../../testdata/avl-vendor/mapping.json", "../../testdata/avl-vendor/valid.json"}, &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "choose exactly one mode") {
		t.Fatalf("error = %v, want explicit mode selection error", err)
	}
}

func TestRunDryRunWritesTelemetryToStdoutAndDiagnosticsToStderr(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{
		"--dry-run",
		"--reference-time", "2026-05-04T12:00:00Z",
		"--mapping", "../../testdata/avl-vendor/mapping.json",
		"../../testdata/avl-vendor/valid.json",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run dry-run: %v; stderr=%s", err, stderr.String())
	}
	var events []telemetry.Event
	if err := json.Unmarshal(stdout.Bytes(), &events); err != nil {
		t.Fatalf("stdout is not telemetry JSON array: %v; stdout=%s", err, stdout.String())
	}
	if len(events) != 1 || events[0].AgencyID != "demo-agency" || !events[0].Valid() {
		t.Fatalf("events = %+v", events)
	}
	var diagnostics []avladapter.Diagnostic
	if err := json.Unmarshal(stderr.Bytes(), &diagnostics); err != nil {
		t.Fatalf("stderr is not diagnostics JSON array: %v; stderr=%s", err, stderr.String())
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %+v, want none", diagnostics)
	}
	combined := strings.ToLower(stdout.String() + stderr.String())
	for _, forbidden := range []string{"token", "authorization", "password", "secret", "credential", "postgres://", "private key"} {
		if strings.Contains(combined, forbidden) {
			t.Fatalf("dry-run output contains forbidden secret-like word %q: stdout=%s stderr=%s", forbidden, stdout.String(), stderr.String())
		}
	}
}

func TestRunDryRunBoundaryHelpDoesNotDependOnPhase29B(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := run([]string{"--help"}, &stdout, &stderr); err != nil {
		t.Fatalf("run help: %v", err)
	}
	assertHelpBoundary(t, stdout.String())
	if strings.Contains(stdout.String(), "Phase 29B") {
		t.Fatalf("help should use adapter-kit wording, got: %s", stdout.String())
	}
}

func TestRunSendWritesPrivateSummaryAndPostsTelemetry(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != "/v1/telemetry" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer synthetic-device-token" {
			t.Fatalf("authorization header was not sent")
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			t.Fatalf("content type = %q", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"accepted":true,"ingest_status":"accepted"}`))
	}))
	defer server.Close()

	outputDir := filepath.Join(".cache", "avl-vendor-adapter-cli-test")
	_ = os.RemoveAll(outputDir)
	t.Cleanup(func() { _ = os.RemoveAll(outputDir) })
	t.Setenv(avladapter.EnvTelemetryURL, server.URL+"/v1/telemetry")
	t.Setenv(avladapter.EnvOutputDir, outputDir)
	t.Setenv("AVL_ADAPTER_DEVICE_TOKEN", "synthetic-device-token")
	t.Setenv(avladapter.EnvReferenceTime, "2026-05-04T12:00:00Z")

	var stdout, stderr bytes.Buffer
	err := run([]string{
		"--send",
		"--mapping", "../../testdata/avl-vendor/mapping.json",
		"--manifest", "../../testdata/avl-vendor/send-manifest.json",
		"../../testdata/avl-vendor/valid.json",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run send: %v; stderr=%s", err, stderr.String())
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %s, want empty", stderr.String())
	}
	if strings.Contains(stdout.String(), server.Listener.Addr().String()) || strings.Contains(stdout.String(), "synthetic-device-token") {
		t.Fatalf("stdout leaked private data: %s", stdout.String())
	}
	for _, name := range []string{"summary.json", "summary.md", "manifest.json", "manifest.md", "diagnostics.json"} {
		if _, err := os.Stat(filepath.Join(outputDir, name)); err != nil {
			t.Fatalf("missing output file %s: %v", name, err)
		}
	}
}

func TestRunReviewDiagnosticsForSeparateBatchFixtures(t *testing.T) {
	for _, tc := range []struct {
		name    string
		fixture string
		code    string
	}{
		{name: "stale", fixture: "stale-timestamp.json", code: avladapter.CodeStaleTimestamp},
		{name: "future", fixture: "future-timestamp.json", code: avladapter.CodeFutureTimestamp},
		{name: "low accuracy", fixture: "low-gps-accuracy.json", code: avladapter.CodeLowGPSAccuracy},
		{name: "duplicate", fixture: "duplicate-batch.json", code: avladapter.CodeDuplicateObservation},
		{name: "out of order", fixture: "out-of-order-batch.json", code: avladapter.CodeOutOfOrderObservation},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			err := run([]string{
				"--dry-run",
				"--reference-time", "2026-05-04T12:00:00Z",
				"--mapping", "../../testdata/avl-vendor/mapping.json",
				"../../testdata/avl-vendor/" + tc.fixture,
			}, &stdout, &stderr)
			if err != nil {
				t.Fatalf("run dry-run: %v; stderr=%s", err, stderr.String())
			}
			var diagnostics []avladapter.Diagnostic
			if err := json.Unmarshal(stderr.Bytes(), &diagnostics); err != nil {
				t.Fatalf("stderr is not diagnostics JSON: %v; stderr=%s", err, stderr.String())
			}
			assertCLIDiagnostic(t, diagnostics, tc.code)
		})
	}
}

func assertHelpBoundary(t *testing.T, help string) {
	t.Helper()
	lower := strings.Join(strings.Fields(strings.ToLower(help)), " ")
	for _, phrase := range []string{
		"--dry-run",
		"--send",
		"not telemetry ingest status",
		"not vendor compatibility proof",
	} {
		if !strings.Contains(lower, phrase) {
			t.Fatalf("help output missing boundary phrase %q: %s", phrase, help)
		}
	}
}

func assertCLIDiagnostic(t *testing.T, diagnostics []avladapter.Diagnostic, code string) {
	t.Helper()
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code {
			return
		}
	}
	t.Fatalf("diagnostic code %q not found in %+v", code, diagnostics)
}

func TestRunHardErrorsPrintPartialDryRunOutputAndExitNonzero(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{
		"--dry-run",
		"--reference-time", "2026-05-04T12:00:00Z",
		"--mapping", "../../testdata/avl-vendor/mapping.json",
		"../../testdata/avl-vendor/batch-mixed.json",
	}, &stdout, &stderr)
	var exit exitError
	if !errors.As(err, &exit) || !exit.silent {
		t.Fatalf("error = %T %v, want silent exitError", err, err)
	}
	var events []telemetry.Event
	if err := json.Unmarshal(stdout.Bytes(), &events); err != nil {
		t.Fatalf("stdout is not telemetry JSON array: %v; stdout=%s", err, stdout.String())
	}
	if len(events) != 1 {
		t.Fatalf("events len = %d, want 1 partial dry-run transform output", len(events))
	}
	var diagnostics []avladapter.Diagnostic
	if err := json.Unmarshal(stderr.Bytes(), &diagnostics); err != nil {
		t.Fatalf("stderr is not diagnostics JSON array: %v; stderr=%s", err, stderr.String())
	}
	if !avladapter.HasHardErrors(diagnostics) {
		t.Fatalf("diagnostics should include hard errors: %+v", diagnostics)
	}
}

func TestRunHardErrorsWithNoValidRecordsPrintsEmptyArray(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{
		"--dry-run",
		"--reference-time", "2026-05-04T12:00:00Z",
		"--mapping", "../../testdata/avl-vendor/mapping.json",
		"../../testdata/avl-vendor/unknown-vendor-vehicle.json",
	}, &stdout, &stderr)
	var exit exitError
	if !errors.As(err, &exit) {
		t.Fatalf("error = %T %v, want exitError", err, err)
	}
	if strings.TrimSpace(stdout.String()) != "[]" {
		t.Fatalf("stdout = %q, want []", stdout.String())
	}
	if !json.Valid(stderr.Bytes()) {
		t.Fatalf("stderr is not valid JSON: %s", stderr.String())
	}
}
