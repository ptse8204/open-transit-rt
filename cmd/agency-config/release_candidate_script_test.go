package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseCandidateCheckDryRunExactFilesAndClaimFlags(t *testing.T) {
	root := releaseCandidateRepoRoot(t)
	outputRel := releaseCandidateTempRel(t, root, "dry-run")
	cmd := releaseCandidateCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(), "OUTPUT_DIR="+outputRel, "FORCE=true")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("release candidate dry-run failed: %v\n%s", err, out)
	}
	assertReleaseCandidateExactFiles(t, root, outputRel)
	assertReleaseCandidateNoLeaks(t, root, outputRel, string(out))
	summary := readReleaseCandidateSummary(t, root, outputRel)
	if summary["dry_run"] != true {
		t.Fatalf("dry_run = %v, want true", summary["dry_run"])
	}
	flags := summary["claim_flags"].(map[string]any)
	for _, key := range []string{
		"retained_evidence_created",
		"consumer_statuses_changed",
		"compliance_claimed",
		"production_readiness_claimed",
		"hosted_saas_claimed",
		"agency_adoption_claimed",
		"consumer_acceptance_claimed",
		"vendor_compatibility_claimed",
		"production_grade_eta_claimed",
		"release_published",
		"tag_created",
		"image_pushed",
	} {
		if flags[key] != false {
			t.Fatalf("claim flag %s = %v, want false", key, flags[key])
		}
	}
}

func TestReleaseCandidateCheckRejectsEvidenceOutput(t *testing.T) {
	root := releaseCandidateRepoRoot(t)
	cmd := releaseCandidateCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(), "OUTPUT_DIR=docs/evidence/release-candidate-check", "ALLOW_UNIGNORED_OUTPUT_DIR=true")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected evidence-like output rejection, got success: %s", out)
	}
	if !strings.Contains(string(out), "OUTPUT_DIR must not be evidence-like") {
		t.Fatalf("unexpected error: %s", out)
	}
}

func TestReleaseCandidateCheckRejectsSymlinkOutput(t *testing.T) {
	root := releaseCandidateRepoRoot(t)
	base := filepath.Join(root, ".cache", "release-candidate-check-test", "symlink")
	target := filepath.Join(base, "target")
	link := filepath.Join(base, "link")
	t.Cleanup(func() { _ = os.RemoveAll(base) })
	if err := os.MkdirAll(target, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	cmd := releaseCandidateCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(), "OUTPUT_DIR="+filepath.ToSlash(filepath.Join(".cache", "release-candidate-check-test", "symlink", "link")))
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected symlink output rejection, got success: %s", out)
	}
	if !strings.Contains(string(out), "OUTPUT_DIR must not contain symlink directories") {
		t.Fatalf("unexpected error: %s", out)
	}
}

func TestReleaseCandidateCheckRejectsNonEmptyOutputWithoutForce(t *testing.T) {
	root := releaseCandidateRepoRoot(t)
	outputRel := releaseCandidateTempRel(t, root, "non-empty")
	outputAbs := filepath.Join(root, outputRel)
	if err := os.MkdirAll(outputAbs, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputAbs, "existing.txt"), []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := releaseCandidateCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(), "OUTPUT_DIR="+outputRel)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-empty output rejection, got success: %s", out)
	}
	if !strings.Contains(string(out), "OUTPUT_DIR exists and is non-empty") {
		t.Fatalf("unexpected error: %s", out)
	}
}

func TestReleaseCandidateCheckHelpAndDocsBoundary(t *testing.T) {
	root := releaseCandidateRepoRoot(t)
	cmd := releaseCandidateCommand(root, "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help failed: %v\n%s", err, out)
	}
	text := string(out) + readReleaseCandidateDoc(t, root)
	for _, want := range []string{
		"private local release-candidate diagnostics",
		"never tags",
		"creates retained evidence",
		"changes consumer statuses",
		"CAL-ITP/Caltrans compliance",
		"consumer acceptance",
		"vendor compatibility",
		"production readiness",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing boundary phrase %q", want)
		}
	}
}

func releaseCandidateCommand(root string, args ...string) *exec.Cmd {
	all := append([]string{filepath.Join(root, "scripts", "release-candidate-check.sh")}, args...)
	cmd := exec.Command(all[0], all[1:]...)
	cmd.Dir = root
	return cmd
}

func releaseCandidateRepoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func releaseCandidateTempRel(t *testing.T, root, name string) string {
	t.Helper()
	rel := filepath.ToSlash(filepath.Join(".cache", "release-candidate-check-test", name, strings.ReplaceAll(t.Name(), "/", "-")))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(root, rel)) })
	return rel
}

func readReleaseCandidateSummary(t *testing.T, root, outputRel string) map[string]any {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(root, outputRel, "summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	var summary map[string]any
	if err := json.Unmarshal(body, &summary); err != nil {
		t.Fatal(err)
	}
	return summary
}

func assertReleaseCandidateExactFiles(t *testing.T, root, outputRel string) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(root, outputRel))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			names = append(names, entry.Name())
		}
	}
	got := strings.Join(names, ",")
	want := "check-log.txt,manifest.json,manifest.md,summary.json,summary.md"
	if got != want {
		t.Fatalf("output files = %s, want %s", got, want)
	}
}

func assertReleaseCandidateNoLeaks(t *testing.T, root, outputRel, terminal string) {
	t.Helper()
	combined := terminal
	for _, name := range []string{"summary.json", "summary.md", "manifest.json", "manifest.md", "check-log.txt"} {
		body, err := os.ReadFile(filepath.Join(root, outputRel, name))
		if err != nil {
			t.Fatal(err)
		}
		combined += "\n" + string(body)
	}
	for _, forbidden := range []string{"Authorization: Bearer", "DATABASE_URL=", "TOKEN=", "PASSWORD=", "SECRET=", "/Users/"} {
		if strings.Contains(combined, forbidden) {
			t.Fatalf("release candidate output leaked %q\n%s", forbidden, combined)
		}
	}
}

func readReleaseCandidateDoc(t *testing.T, root string) string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(root, "docs", "release-candidate-readiness.md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}
