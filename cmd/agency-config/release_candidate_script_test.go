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
	source, ok := summary["source"].(map[string]any)
	if !ok || source["describe"] == nil || source["commit_sha"] == nil || source["pre_tag_review"] != true {
		t.Fatalf("source metadata missing describe, commit_sha, or pre_tag_review: %#v", summary["source"])
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
	sequence, ok := summary["review_sequence"].([]any)
	if !ok || len(sequence) < 6 {
		t.Fatalf("review_sequence length = %d, want at least 6", len(sequence))
	}
	for _, id := range []string{
		"check_links",
		"product_ui_smoke",
		"auth_production_boundary",
		"auth_password_login",
		"auth_bootstrap_single_use",
		"auth_logout_expiry",
		"auth_cookie_post_csrf",
		"dashboard_issue_priority",
		"setup_wizard_skip_reminder",
		"connector_examples_vs_configured",
		"connector_dry_run_redaction",
		"product_acceptance_audit",
		"product_language_audit",
		"ui_layout_audit",
		"operations_route_inventory",
		"api_contract_check",
		"stable_filter_check",
		"external_connection_check",
		"adapter_conformance",
		"connector_examples",
		"gtfsrt_conformance",
	} {
		check := releaseCandidateCheckByID(t, summary, id)
		if check["status"] != "not_checked" {
			t.Fatalf("%s dry-run status = %v, want not_checked", id, check["status"])
		}
	}
	inputs, ok := summary["release_note_inputs"].(map[string]any)
	if !ok || inputs["validation"] == nil || inputs["claim_boundaries"] == nil {
		t.Fatalf("release_note_inputs missing validation or claim boundaries: %#v", summary["release_note_inputs"])
	}
	matrix, ok := summary["package_audit_matrix"].([]any)
	if !ok || len(matrix) < 5 {
		t.Fatalf("package_audit_matrix length = %d, want at least 5", len(matrix))
	}
	summaryText := readReleaseCandidateOutputFile(t, root, outputRel, "summary.md")
	for _, want := range []string{
		"First Release-Candidate Workflow",
		"Release Note Inputs",
		"Local Package Audit Matrix",
		"pre-tag local diagnostics",
		"API/feed/extension contract check",
		"Product UI smoke",
		"Production auth boundary",
		"Password login issues admin_session",
		"Connector examples remain separate",
		"Connector dry-run records redacted results",
		"GTFS-RT conformance harness",
	} {
		if !strings.Contains(summaryText, want) {
			t.Fatalf("summary.md missing %q", want)
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
	base := filepath.Join(root, releaseCandidateTempRel(t, root, "symlink"))
	target := filepath.Join(base, "target")
	link := filepath.Join(base, "link")
	if err := os.MkdirAll(target, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	cmd := releaseCandidateCommand(root, "--dry-run")
	linkRel, err := filepath.Rel(root, link)
	if err != nil {
		t.Fatal(err)
	}
	cmd.Env = append(os.Environ(), "OUTPUT_DIR="+filepath.ToSlash(linkRel))
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

func TestReleaseCandidateCheckPinsValidatorToolingDespiteAmbientStub(t *testing.T) {
	root := releaseCandidateRepoRoot(t)
	outputRel := releaseCandidateTempRel(t, root, "validator-pinned")
	cmd := releaseCandidateCommand(root)
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"FORCE=true",
		"VALIDATOR_TOOLING_MODE=stub",
		"GTFS_VALIDATOR_PATH="+filepath.Join(root, outputRel, "missing-validator.jar"),
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected validator tooling blocker despite ambient stub mode, got success:\n%s", out)
	}
	assertReleaseCandidateExactFiles(t, root, outputRel)

	summary := readReleaseCandidateSummary(t, root, outputRel)
	if summary["overall_status"] != "blocker" {
		t.Fatalf("overall_status = %v, want blocker", summary["overall_status"])
	}
	check := releaseCandidateCheckByID(t, summary, "validators_check")
	if check["status"] != "blocker" {
		t.Fatalf("validators_check status = %v, want blocker; summary=%#v", check["status"], check)
	}
	detail, _ := check["detail"].(string)
	if strings.Contains(detail, "VALIDATOR_TOOLING_MODE=stub") {
		t.Fatalf("validators_check detail shows ambient stub bypass instead of pinned blocker: %q", detail)
	}
	if !strings.Contains(detail, "pinned tooling") {
		t.Fatalf("validators_check detail = %q, want actionable pinned tooling wording", detail)
	}
	if !strings.Contains(detail, "missing pinned tooling") && !strings.Contains(detail, "misconfigured pinned tooling") {
		t.Fatalf("validators_check detail = %q, want missing/misconfigured pinned tooling wording", detail)
	}
	summaryText := readReleaseCandidateOutputFile(t, root, outputRel, "summary.md")
	if !strings.Contains(summaryText, detail) {
		t.Fatalf("summary.md missing validators_check blocker detail %q\n%s", detail, summaryText)
	}
	logText := readReleaseCandidateOutputFile(t, root, outputRel, "check-log.txt")
	if !strings.Contains(logText, detail) {
		t.Fatalf("check-log.txt missing full validator blocker detail %q\n%s", detail, logText)
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
		"pre-tag",
		"local diagnostics",
		"API/feed/extension contract",
		"Production auth boundary",
		"first-admin bootstrap",
		"dashboard top-three",
		"Connector examples",
		"product UI",
		"GTFS-RT conformance",
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
	base := filepath.Join(root, ".cache", "release-candidate-check-test")
	if err := os.MkdirAll(base, 0o700); err != nil {
		t.Fatal(err)
	}
	prefix := name + "-" + strings.ReplaceAll(t.Name(), "/", "-") + "-"
	dir, err := os.MkdirTemp(base, prefix)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	rel, err := filepath.Rel(root, dir)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(rel)
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

func releaseCandidateCheckByID(t *testing.T, summary map[string]any, id string) map[string]any {
	t.Helper()
	checks, ok := summary["checks"].([]any)
	if !ok {
		t.Fatalf("summary checks missing or wrong type: %#v", summary["checks"])
	}
	for _, item := range checks {
		check, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("summary check has wrong type: %#v", item)
		}
		if check["id"] == id {
			return check
		}
	}
	t.Fatalf("summary missing check id %q", id)
	return nil
}

func readReleaseCandidateOutputFile(t *testing.T, root, outputRel, name string) string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(root, outputRel, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
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
