package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestOperationsReliabilityScriptExactFilesDefaultFlagsAndNoSend(t *testing.T) {
	root := operationsReliabilityRepoRoot(t)
	base := operationsReliabilityTempRel(t, root, "exact")
	validatorRel := filepath.ToSlash(filepath.Join(base, "validator-health", "20260509T120000Z", "summary.json"))
	doctorRel := filepath.ToSlash(filepath.Join(base, "deployment-doctor", "20260509T120000Z", "summary.json"))
	notifyRel := filepath.ToSlash(filepath.Join(base, "operations-notify", "20260509T120000Z", "summary.json"))
	writeOperationsReliabilityJSON(t, root, validatorRel, map[string]any{"overall_status": "recorded"})
	writeOperationsReliabilityJSON(t, root, doctorRel, map[string]any{"overall_status": "passed"})
	writeOperationsReliabilityJSON(t, root, notifyRel, map[string]any{"overall_status": "info", "destinations": map[string]any{"webhook_present": true}})
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	cmd := operationsReliabilityCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"FORCE=true",
		"VALIDATOR_HEALTH_SUMMARY="+validatorRel,
		"DEPLOYMENT_DOCTOR_SUMMARY="+doctorRel,
		"OPERATIONS_NOTIFY_SUMMARY="+notifyRel,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run failed: %v\n%s", err, out)
	}
	assertOperationsReliabilityExactFiles(t, root, outputRel)
	summary := readOperationsReliabilitySummary(t, root, outputRel)
	monitoring := summary["monitoring_export"].(map[string]any)
	if monitoring["not_sent"] != true || monitoring["webhook_present"] != true || monitoring["webhook_send_enabled"] != false || monitoring["email_send_enabled"] != false || monitoring["destination_values_recorded"] != false {
		t.Fatalf("monitoring export no-send summary invalid: %+v", monitoring)
	}
	privateOps := summary["private_ops_summary"].(map[string]any)
	if privateOps["notification_not_sent"] != true || privateOps["monitoring_export_status"] == "" || privateOps["summary_json"] != "summary.json" || privateOps["manifest_json"] != "manifest.json" {
		t.Fatalf("private ops summary invalid: %+v", privateOps)
	}
	exportFormats := privateOps["export_formats"].(map[string]any)
	for _, key := range []string{"feed_health", "connector_health", "validator_posture", "telemetry_freshness", "maintenance_tasks"} {
		if exportFormats[key] == "" {
			t.Fatalf("export format %s missing: %+v", key, exportFormats)
		}
	}
	if privateOps["redaction_boundary"] == "" {
		t.Fatalf("redaction boundary missing: %+v", privateOps)
	}
	flags := summary["claim_flags"].(map[string]any)
	for _, flag := range []string{
		"external_evidence_created",
		"final_root_evidence_created",
		"consumer_statuses_changed",
		"compliance_claimed",
		"production_readiness_claimed",
		"sla_claimed",
		"uptime_guarantee_claimed",
		"hosted_saas_claimed",
		"agency_adoption_claimed",
		"consumer_acceptance_claimed",
		"vendor_compatibility_claimed",
		"production_grade_eta_claimed",
	} {
		if flags[flag] != false {
			t.Fatalf("flag %s = %+v, want false", flag, flags[flag])
		}
	}
	assertOperationsReliabilityNoLeakage(t, root, outputRel, string(out))
}

func TestOperationsReliabilityScriptMissingSourcesAreNotOK(t *testing.T) {
	root := operationsReliabilityRepoRoot(t)
	base := operationsReliabilityTempRel(t, root, "missing")
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	cmd := operationsReliabilityCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"FORCE=true",
		"VALIDATOR_HEALTH_SUMMARY="+filepath.ToSlash(filepath.Join(base, "missing-validator", "summary.json")),
		"DEPLOYMENT_DOCTOR_SUMMARY="+filepath.ToSlash(filepath.Join(base, "missing-doctor", "summary.json")),
		"OPERATIONS_NOTIFY_SUMMARY="+filepath.ToSlash(filepath.Join(base, "missing-notify", "summary.json")),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("missing-source run failed: %v\n%s", err, out)
	}
	summary := readOperationsReliabilitySummary(t, root, outputRel)
	for _, section := range []string{"backup_restore", "alerting", "availability_sampling", "long_running_operations"} {
		row := summary[section].(map[string]any)
		if row["status"] == "ok" {
			t.Fatalf("%s missing source became ok: %+v", section, row)
		}
	}
}

func TestOperationsReliabilityScriptRejectsUnsafePathsAndSources(t *testing.T) {
	root := operationsReliabilityRepoRoot(t)
	cmd := operationsReliabilityCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(), "OUTPUT_DIR=docs/evidence/operations-reliability-test", "ALLOW_UNIGNORED_OUTPUT_DIR=true")
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("expected evidence output rejection, got success: %s", out)
	}

	base := operationsReliabilityTempRel(t, root, "unsafe")
	sourceRel := filepath.ToSlash(filepath.Join(base, "validator-health", "20260509T120000Z", "summary.json"))
	writeOperationsReliabilityRaw(t, root, sourceRel, []byte(`{"overall_status":"ok","database_url":"postgres://user:pass@localhost/db"}`))
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	cmd = operationsReliabilityCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"FORCE=true",
		"VALIDATOR_HEALTH_SUMMARY="+sourceRel,
	)
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("expected DB URL source rejection, got success: %s", out)
	}

	targetRel := filepath.ToSlash(filepath.Join(base, "target", "summary.json"))
	writeOperationsReliabilityJSON(t, root, targetRel, map[string]any{"overall_status": "ok"})
	linkRel := filepath.ToSlash(filepath.Join(base, "link-summary.json"))
	if err := os.Symlink(filepath.Join(root, targetRel), filepath.Join(root, linkRel)); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	cmd = operationsReliabilityCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"FORCE=true",
		"VALIDATOR_HEALTH_SUMMARY="+linkRel,
	)
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("expected symlink source rejection, got success: %s", out)
	}
}

func TestOperationsReliabilityScriptAllowsGeneratedSecretStatusCategory(t *testing.T) {
	root := operationsReliabilityRepoRoot(t)
	base := operationsReliabilityTempRel(t, root, "generated-secret-category")
	validatorRel := filepath.ToSlash(filepath.Join(base, "validator-health", "20260509T120000Z", "summary.json"))
	doctorRel := filepath.ToSlash(filepath.Join(base, "deployment-doctor", "20260509T120000Z", "summary.json"))
	notifyRel := filepath.ToSlash(filepath.Join(base, "operations-notify", "20260509T120000Z", "summary.json"))
	writeOperationsReliabilityJSON(t, root, validatorRel, map[string]any{"overall_status": "recorded"})
	writeOperationsReliabilityJSON(t, root, doctorRel, map[string]any{
		"overall_status": "warning",
		"checks": []map[string]string{{
			"category": "install_recovery",
			"name":     "environment_preflight",
			"status":   "blocker",
			"detail":   "environment_preflight=blocker; categories=env:blocker,generated_secret:blocker,url:blocker",
		}},
	})
	writeOperationsReliabilityJSON(t, root, notifyRel, map[string]any{"overall_status": "info"})
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	cmd := operationsReliabilityCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"FORCE=true",
		"VALIDATOR_HEALTH_SUMMARY="+validatorRel,
		"DEPLOYMENT_DOCTOR_SUMMARY="+doctorRel,
		"OPERATIONS_NOTIFY_SUMMARY="+notifyRel,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated_secret category should be accepted: %v\n%s", err, out)
	}

	unsafeRel := filepath.ToSlash(filepath.Join(base, "deployment-doctor", "20260509T120001Z", "summary.json"))
	writeOperationsReliabilityJSON(t, root, unsafeRel, map[string]any{"overall_status": "warning", "summary": "secret: inline-value"})
	cmd = operationsReliabilityCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"FORCE=true",
		"VALIDATOR_HEALTH_SUMMARY="+validatorRel,
		"DEPLOYMENT_DOCTOR_SUMMARY="+unsafeRel,
		"OPERATIONS_NOTIFY_SUMMARY="+notifyRel,
	)
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("expected inline secret marker rejection, got success: %s", out)
	}
}

func TestOperationsReliabilityScriptRejectsOversizedSource(t *testing.T) {
	root := operationsReliabilityRepoRoot(t)
	base := operationsReliabilityTempRel(t, root, "oversized")
	sourceRel := filepath.ToSlash(filepath.Join(base, "validator-health", "20260509T120000Z", "summary.json"))
	writeOperationsReliabilityRaw(t, root, sourceRel, []byte(`{"overall_status":"`+strings.Repeat("x", 2048)+`"}`))
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	cmd := operationsReliabilityCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"FORCE=true",
		"MAX_SOURCE_BYTES=128",
		"VALIDATOR_HEALTH_SUMMARY="+sourceRel,
	)
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("expected oversized source rejection, got success: %s", out)
	}
}

func operationsReliabilityCommand(root string, args ...string) *exec.Cmd {
	all := append([]string{filepath.Join(root, "scripts", "operations-reliability.sh")}, args...)
	cmd := exec.Command("sh", all...)
	cmd.Dir = root
	return cmd
}

func operationsReliabilityRepoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func operationsReliabilityTempRel(t *testing.T, root, name string) string {
	t.Helper()
	base := filepath.Join(root, ".cache", "operations-reliability-test")
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

func writeOperationsReliabilityJSON(t *testing.T, root, rel string, payload any) {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	writeOperationsReliabilityRaw(t, root, rel, body)
}

func writeOperationsReliabilityRaw(t *testing.T, root, rel string, body []byte) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
}

func readOperationsReliabilitySummary(t *testing.T, root, outputRel string) map[string]any {
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

func assertOperationsReliabilityExactFiles(t *testing.T, root, outputRel string) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(root, outputRel))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	got := strings.Join(names, ",")
	want := "manifest.json,manifest.md,reliability-review.txt,summary.json,summary.md"
	if got != want {
		t.Fatalf("files = %s, want %s", got, want)
	}
}

func assertOperationsReliabilityNoLeakage(t *testing.T, root, outputRel, stdout string) {
	t.Helper()
	combined := stdout
	for _, name := range []string{"summary.json", "summary.md", "manifest.json", "manifest.md", "reliability-review.txt"} {
		body, err := os.ReadFile(filepath.Join(root, outputRel, name))
		if err != nil {
			t.Fatal(err)
		}
		combined += string(body)
		if len(body) > 65536 {
			t.Fatalf("%s size = %d, want bounded", name, len(body))
		}
	}
	for _, forbidden := range []string{"postgres://", "DATABASE_URL", "Authorization:", "Cookie:", "BEGIN PRIVATE KEY", "secret-token", "https://hooks.", "/Users/"} {
		if strings.Contains(combined, forbidden) {
			t.Fatalf("operations-reliability output leaked %q\n%s", forbidden, combined)
		}
	}
}
