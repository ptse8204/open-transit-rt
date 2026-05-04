package main

import (
	"bytes"
	"encoding/json"
	"errors"
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
	if !strings.Contains(stdout.String(), "Usage:") || !strings.Contains(stdout.String(), "--dry-run") {
		t.Fatalf("help output missing usage: %s", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunRequiresDryRun(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"--mapping", "../../testdata/avl-vendor/mapping.json", "../../testdata/avl-vendor/valid.json"}, &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "send mode is not implemented") {
		t.Fatalf("error = %v, want send-mode-not-implemented", err)
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
