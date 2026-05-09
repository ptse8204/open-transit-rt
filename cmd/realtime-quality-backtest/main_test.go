package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := run([]string{"help"}, &stdout, &stderr); err != nil {
		t.Fatalf("help: %v", err)
	}
	help := strings.Join(strings.Fields(stdout.String()), " ")
	if !strings.Contains(help, "Usage:") || !strings.Contains(help, "not production-grade ETA proof") {
		t.Fatalf("help missing boundary wording: %s", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunWritesExactPrivateFiles(t *testing.T) {
	outputDir := filepath.Join(".cache", "realtime-quality-backtest-cli-test")
	_ = os.RemoveAll(outputDir)
	t.Cleanup(func() { _ = os.RemoveAll(outputDir) })
	var stdout, stderr bytes.Buffer
	err := run([]string{
		"--observed", "../../testdata/realtime-quality-backtest/observed-events.json",
		"--predictions", "../../testdata/realtime-quality-backtest/prediction-samples.json",
		"--output-dir", outputDir,
		"--generated-at", "2026-05-09T20:00:00Z",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v; stderr=%s", err, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %s, want empty", stderr.String())
	}
	var terminal map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &terminal); err != nil {
		t.Fatalf("stdout is not JSON: %v; %s", err, stdout.String())
	}
	if terminal["maturity_gate"] != "diagnostic_watch" {
		t.Fatalf("stdout = %+v, want diagnostic_watch", terminal)
	}
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if len(entries) != 5 {
		t.Fatalf("output files = %d, want 5", len(entries))
	}
	for _, name := range []string{"summary.json", "summary.md", "metrics.json", "metrics.md", "manifest.json"} {
		if _, err := os.Stat(filepath.Join(outputDir, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
	combined := stdout.String()
	for _, name := range []string{"summary.json", "summary.md", "metrics.json", "metrics.md", "manifest.json"} {
		raw, err := os.ReadFile(filepath.Join(outputDir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		combined += string(raw)
	}
	for _, forbidden := range []string{"../../testdata/realtime-quality-backtest", "/Users/", "authorization:", "postgres://", `"ready"`, `"production_ready"`, `"compliant"`, `"accepted"`} {
		if strings.Contains(strings.ToLower(combined), strings.ToLower(forbidden)) {
			t.Fatalf("output leaked forbidden text %q", forbidden)
		}
	}
}

func TestRunRejectsUnsafeOutputPaths(t *testing.T) {
	for _, outputDir := range []string{
		filepath.Join("docs", "evidence", "phase-50-cli-test"),
		filepath.Join("tmp", "evidence-like-backtest"),
		filepath.Join("..", "outside"),
	} {
		t.Run(outputDir, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			err := run([]string{
				"--observed", "../../testdata/realtime-quality-backtest/observed-events.json",
				"--predictions", "../../testdata/realtime-quality-backtest/prediction-samples.json",
				"--output-dir", outputDir,
				"--generated-at", "2026-05-09T20:00:00Z",
			}, &stdout, &stderr)
			if err == nil || !strings.Contains(err.Error(), "invalid output directory") {
				t.Fatalf("error = %v, want invalid output directory", err)
			}
		})
	}
}

func TestRunRejectsSymlinkAncestor(t *testing.T) {
	base := filepath.Join(".cache", "realtime-quality-backtest-symlink-test")
	_ = os.RemoveAll(base)
	t.Cleanup(func() { _ = os.RemoveAll(base) })
	target := filepath.Join(base, "target")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	link := filepath.Join(base, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("symlink: %v", err)
	}
	var stdout, stderr bytes.Buffer
	err := run([]string{
		"--observed", "../../testdata/realtime-quality-backtest/observed-events.json",
		"--predictions", "../../testdata/realtime-quality-backtest/prediction-samples.json",
		"--output-dir", filepath.Join(link, "out"),
		"--generated-at", "2026-05-09T20:00:00Z",
	}, &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("error = %v, want symlink rejection", err)
	}
}
