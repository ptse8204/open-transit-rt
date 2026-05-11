package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBootstrapDevHelpDocumentsPreflightAndBoundaries(t *testing.T) {
	root := bootstrapDevRepoRoot(t)
	cmd := exec.Command(filepath.Join(root, "scripts", "bootstrap-dev.sh"), "--help")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bootstrap help failed: %v\n%s", err, out)
	}
	text := string(out)
	for _, want := range []string{
		"scripts/bootstrap-dev.sh --check",
		"Docker CLI missing",
		"Docker daemon stopped",
		"Port 55432 occupied",
		"local development setup only",
		"not hosted service",
		"not production-readiness",
		"not compliance proof",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("bootstrap help missing %q\n%s", want, text)
		}
	}
}

func TestBootstrapDevRejectsUnknownArgumentBeforeMutation(t *testing.T) {
	root := bootstrapDevRepoRoot(t)
	cmd := exec.Command(filepath.Join(root, "scripts", "bootstrap-dev.sh"), "--unknown")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected unknown argument to fail, got success:\n%s", out)
	}
	text := string(out)
	for _, want := range []string{
		"Unknown argument: --unknown",
		"scripts/bootstrap-dev.sh --check",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("unknown-argument output missing %q\n%s", want, text)
		}
	}
}

func TestAgencyLocalAppHelpDocumentsFirstRunBlockers(t *testing.T) {
	root := bootstrapDevRepoRoot(t)
	cmd := exec.Command(filepath.Join(root, "scripts", "agency-local-app.sh"), "--help")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("agency local app help failed: %v\n%s", err, out)
	}
	text := string(out)
	for _, want := range []string{
		"Docker is missing or stopped",
		"Docker Compose plugin",
		"host ports 8080 or 55432",
		"first-run image/module pulls",
		"local evaluation only",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("agency-local-app help missing %q\n%s", want, text)
		}
	}
}

func bootstrapDevRepoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}
