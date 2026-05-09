package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOperationsNotifyScriptDryRunNoSendAndExactOutput(t *testing.T) {
	root := operationsNotifyRepoRoot(t)
	fakeBin := operationsNotifyFakeNoSendBin(t)
	outputRel := operationsNotifyTempRel(t, root, "dry-run")
	cmd := operationsNotifyCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(),
		"PATH="+fakeBin+":"+os.Getenv("PATH"),
		"OUTPUT_DIR="+outputRel,
		"FORCE=true",
		"NOTIFY_WEBHOOK_URL=https://hooks.example/services/T000/B000/secret-token",
		"NOTIFY_EMAIL_TO=ops@example.test",
		"VALIDATOR_HEALTH_SUMMARY=.cache/operations-notify-test/missing-validator/summary.json",
		"DEPLOYMENT_DOCTOR_SUMMARY=.cache/operations-notify-test/missing-doctor/summary.json",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("dry-run failed: %v\n%s", err, out)
	}
	assertOperationsNotifyOutputSafe(t, root, outputRel, string(out), []string{"secret-token", "ops@example.test", "/tmp/", "/var/", "/etc/", "/Users/"})
	summary := readOperationsNotifySummary(t, root, outputRel)
	for _, flag := range []string{
		"external_evidence_created",
		"consumer_statuses_changed",
		"compliance_claimed",
		"production_readiness_claimed",
		"hosted_saas_claimed",
		"agency_adoption_claimed",
		"consumer_acceptance_claimed",
		"vendor_compatibility_claimed",
		"production_grade_eta_claimed",
		"notification_sent",
	} {
		if !summaryBool(t, summary, flag, false) {
			t.Fatalf("%s not present as false: %+v", flag, summary[flag])
		}
	}
	destinations := summary["destinations"].(map[string]any)
	if destinations["webhook_present"] != true || destinations["email_present"] != true {
		t.Fatalf("destination booleans = %+v, want true/true", destinations)
	}
	notification := summary["notification"].(map[string]any)
	if notification["not_sent"] != true {
		t.Fatalf("notification.not_sent = %+v", notification["not_sent"])
	}
}

func TestOperationsNotifySeverityMapping(t *testing.T) {
	root := operationsNotifyRepoRoot(t)
	cases := []struct {
		name       string
		validator  any
		doctor     any
		want       string
		strict     bool
		wantErr    bool
		malformedV bool
		malformedD bool
	}{
		{
			name:      "validator blocked doctor passed",
			validator: operationsNotifyValidatorSummary("blocked"),
			doctor:    operationsNotifyDoctorSummary("passed", nil),
			want:      "blocked",
		},
		{
			name:      "doctor blocker",
			validator: operationsNotifyValidatorSummary("recorded"),
			doctor:    operationsNotifyDoctorSummary("blocker", []map[string]any{{"category": "env", "name": "APP_ENV", "status": "blocker", "detail": "missing required reference env"}}),
			want:      "blocked",
		},
		{
			name:       "malformed default",
			doctor:     operationsNotifyDoctorSummary("passed", nil),
			want:       "needs_review",
			malformedV: true,
		},
		{
			name:   "missing default",
			doctor: operationsNotifyDoctorSummary("passed", nil),
			want:   "needs_review",
		},
		{
			name:    "missing strict",
			doctor:  operationsNotifyDoctorSummary("passed", nil),
			want:    "blocked",
			strict:  true,
			wantErr: true,
		},
		{
			name:       "malformed strict",
			doctor:     operationsNotifyDoctorSummary("passed", nil),
			want:       "blocked",
			strict:     true,
			wantErr:    true,
			malformedV: true,
		},
		{
			name:      "all healthy",
			validator: operationsNotifyValidatorSummary("recorded"),
			doctor:    operationsNotifyDoctorSummary("passed", nil),
			want:      "info",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			base := operationsNotifyTempRel(t, root, "severity")
			validatorRel := filepath.ToSlash(filepath.Join(base, "validator-health", "20260509T120000Z", "summary.json"))
			doctorRel := filepath.ToSlash(filepath.Join(base, "deployment-doctor", "20260509T120000Z", "summary.json"))
			if tc.malformedV {
				writeOperationsNotifyRaw(t, root, validatorRel, []byte(`{"overall_status":`))
			} else if tc.validator != nil {
				writeOperationsNotifyJSON(t, root, validatorRel, tc.validator)
			}
			if tc.malformedD {
				writeOperationsNotifyRaw(t, root, doctorRel, []byte(`{"overall_status":`))
			} else if tc.doctor != nil {
				writeOperationsNotifyJSON(t, root, doctorRel, tc.doctor)
			}
			outputRel := filepath.ToSlash(filepath.Join(base, "out"))
			cmd := operationsNotifyCommand(root)
			env := []string{
				"OUTPUT_DIR=" + outputRel,
				"VALIDATOR_HEALTH_SUMMARY=" + validatorRel,
				"DEPLOYMENT_DOCTOR_SUMMARY=" + doctorRel,
				"FORCE=true",
			}
			if tc.strict {
				env = append(env, "STRICT_OPERATIONS_NOTIFY=true")
			}
			cmd.Env = append(os.Environ(), env...)
			out, err := cmd.CombinedOutput()
			if tc.wantErr && err == nil {
				t.Fatalf("expected strict error, got success: %s", out)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("run failed: %v\n%s", err, out)
			}
			assertOperationsNotifyOutputSafe(t, root, outputRel, string(out), nil)
			summary := readOperationsNotifySummary(t, root, outputRel)
			notification := summary["notification"].(map[string]any)
			if got := notification["severity"]; got != tc.want {
				t.Fatalf("severity = %v, want %s\nsummary=%+v", got, tc.want, summary)
			}
		})
	}
}

func TestOperationsNotifySourceDiscoveryTimestampRules(t *testing.T) {
	root := operationsNotifyRepoRoot(t)
	base := operationsNotifyTempRel(t, root, "discovery")
	for _, rel := range []string{
		filepath.ToSlash(filepath.Join(".cache", "validator-health", "99991231T235900Z")),
		filepath.ToSlash(filepath.Join(".cache", "validator-health", "99991231T235901Z")),
		filepath.ToSlash(filepath.Join(".cache", "validator-health", "not-a-timestamp-ops-test")),
		filepath.ToSlash(filepath.Join(".cache", "deployment-doctor", "99991231T235900Z")),
		filepath.ToSlash(filepath.Join(".cache", "deployment-doctor", "99991231T235901Z")),
	} {
		t.Cleanup(func(rel string) func() {
			return func() { _ = os.RemoveAll(filepath.Join(root, rel)) }
		}(rel))
	}
	writeOperationsNotifyJSON(t, root, ".cache/validator-health/not-a-timestamp-ops-test/summary.json", operationsNotifyValidatorSummary("blocked"))
	writeOperationsNotifyJSON(t, root, ".cache/validator-health/99991231T235900Z/summary.json", operationsNotifyValidatorSummary("blocked"))
	writeOperationsNotifyJSON(t, root, ".cache/validator-health/99991231T235901Z/summary.json", operationsNotifyValidatorSummary("recorded"))
	mkdirOperationsNotify(t, root, ".cache/deployment-doctor/99991231T235900Z")
	writeOperationsNotifyJSON(t, root, ".cache/deployment-doctor/99991231T235901Z/summary.json", operationsNotifyDoctorSummary("passed", nil))
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	cmd := operationsNotifyCommand(root)
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"FORCE=true",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("discovery run failed: %v\n%s", err, out)
	}
	summary := readOperationsNotifySummary(t, root, outputRel)
	notification := summary["notification"].(map[string]any)
	if notification["severity"] != "info" {
		t.Fatalf("expected latest valid summaries to be healthy, got %+v", notification)
	}
}

func TestOperationsNotifyRejectsSymlinkSourceAndEvidencePaths(t *testing.T) {
	root := operationsNotifyRepoRoot(t)
	base := operationsNotifyTempRel(t, root, "symlink")
	targetRel := filepath.ToSlash(filepath.Join(base, "target", "summary.json"))
	writeOperationsNotifyJSON(t, root, targetRel, operationsNotifyValidatorSummary("recorded"))
	linkRel := filepath.ToSlash(filepath.Join(base, "link-summary.json"))
	if err := os.Symlink(filepath.Join(root, targetRel), filepath.Join(root, linkRel)); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	cmd := operationsNotifyCommand(root)
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"VALIDATOR_HEALTH_SUMMARY="+linkRel,
		"DEPLOYMENT_DOCTOR_SUMMARY="+targetRel,
		"FORCE=true",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("non-strict rejected source should still produce output: %v\n%s", err, out)
	}
	summary := readOperationsNotifySummary(t, root, outputRel)
	sources := summary["source_summaries"].([]any)
	if sources[0].(map[string]any)["status"] != "malformed_source" {
		t.Fatalf("symlink source status = %+v", sources[0])
	}
	evidenceOut := filepath.ToSlash(filepath.Join("docs", "evidence", "operations-notify-test"))
	cmd = operationsNotifyCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(), "OUTPUT_DIR="+evidenceOut, "ALLOW_UNIGNORED_OUTPUT_DIR=true")
	if out, err = cmd.CombinedOutput(); err == nil {
		t.Fatalf("expected evidence output rejection, got success: %s", out)
	}
}

func TestOperationsNotifyNoRecursiveSourceReads(t *testing.T) {
	root := operationsNotifyRepoRoot(t)
	base := operationsNotifyTempRel(t, root, "recursive")
	sourceDir := filepath.ToSlash(filepath.Join(base, "validator-health", "20260509T120000Z"))
	writeOperationsNotifyJSON(t, root, filepath.ToSlash(filepath.Join(sourceDir, "summary.json")), operationsNotifyValidatorSummary("recorded"))
	writeOperationsNotifyJSON(t, root, filepath.ToSlash(filepath.Join(sourceDir, "manifest.json")), map[string]any{"included": []string{"summary.json"}})
	writeOperationsNotifyRaw(t, root, filepath.ToSlash(filepath.Join(sourceDir, "raw-report.json")), []byte(`{"raw_report":"DO-NOT-COPY-RAW-REPORT"}`))
	writeOperationsNotifyRaw(t, root, filepath.ToSlash(filepath.Join(sourceDir, "secret.env")), []byte("TOKEN=DO-NOT-COPY-TOKEN\n"))
	writeOperationsNotifyRaw(t, root, filepath.ToSlash(filepath.Join(sourceDir, "nested", "private.log")), []byte("DO-NOT-COPY-NESTED-LOG\n"))
	doctorRel := filepath.ToSlash(filepath.Join(base, "deployment-doctor", "20260509T120000Z", "summary.json"))
	writeOperationsNotifyJSON(t, root, doctorRel, operationsNotifyDoctorSummary("passed", nil))
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	cmd := operationsNotifyCommand(root)
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"VALIDATOR_HEALTH_SUMMARY="+filepath.ToSlash(filepath.Join(sourceDir, "summary.json")),
		"DEPLOYMENT_DOCTOR_SUMMARY="+doctorRel,
		"FORCE=true",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run failed: %v\n%s", err, out)
	}
	assertOperationsNotifyOutputSafe(t, root, outputRel, string(out), []string{"DO-NOT-COPY-RAW-REPORT", "DO-NOT-COPY-TOKEN", "DO-NOT-COPY-NESTED-LOG", "raw-report.json", "secret.env", "private.log"})
}

func TestOperationsNotifyLargeDeploymentDoctorSummary_1k(t *testing.T) {
	testOperationsNotifyLargeDeploymentDoctorSummary(t, 1000, false)
}

func TestOperationsNotifyLargeDeploymentDoctorSummary_10k(t *testing.T) {
	testOperationsNotifyLargeDeploymentDoctorSummary(t, 10000, false)
}

func TestOperationsNotifyLargeDeploymentDoctorBlockers_10k(t *testing.T) {
	testOperationsNotifyLargeDeploymentDoctorSummary(t, 10000, true)
}

func TestOperationsNotifyLargeValidatorHealthSummary(t *testing.T) {
	root := operationsNotifyRepoRoot(t)
	base := operationsNotifyTempRel(t, root, "large-validator")
	validator := operationsNotifyValidatorSummary("recorded")
	validator.(map[string]any)["unexpected"] = strings.Split(strings.Repeat("x,", 10000), ",")
	validator.(map[string]any)["feeds"] = []map[string]any{
		{"feed_type": "schedule", "health_status": "recorded", "next_action": strings.Repeat("huge-action", 1000)},
	}
	validatorRel := filepath.ToSlash(filepath.Join(base, "validator-health", "20260509T120000Z", "summary.json"))
	doctorRel := filepath.ToSlash(filepath.Join(base, "deployment-doctor", "20260509T120000Z", "summary.json"))
	writeOperationsNotifyJSON(t, root, validatorRel, validator)
	writeOperationsNotifyJSON(t, root, doctorRel, operationsNotifyDoctorSummary("passed", nil))
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	runOperationsNotifySuccess(t, root, outputRel, validatorRel, doctorRel)
	assertOperationsNotifyOutputSafe(t, root, outputRel, "", []string{strings.Repeat("huge-action", 20), "unexpected"})
}

func TestOperationsNotifyHostileSourceSummaries(t *testing.T) {
	root := operationsNotifyRepoRoot(t)
	base := operationsNotifyTempRel(t, root, "hostile")
	validator := operationsNotifyValidatorSummary("blocked")
	validator.(map[string]any)["raw_report"] = map[string]any{"token": "SECRET-RAW"}
	validator.(map[string]any)["stdout"] = "Authorization: Bearer raw-secret-token"
	validator.(map[string]any)["argv"] = []string{"/tmp/private/raw"}
	doctor := operationsNotifyDoctorSummary("blocker", []map[string]any{{
		"category": "admin",
		"name":     "validation",
		"status":   "blocker",
		"detail":   "Review local private diagnostics only.",
	}})
	doctor["stderr"] = "postgres://user:pass@example/db"
	validatorRel := filepath.ToSlash(filepath.Join(base, "validator-health", "20260509T120000Z", "summary.json"))
	doctorRel := filepath.ToSlash(filepath.Join(base, "deployment-doctor", "20260509T120000Z", "summary.json"))
	writeOperationsNotifyJSON(t, root, validatorRel, validator)
	writeOperationsNotifyJSON(t, root, doctorRel, doctor)
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	runOperationsNotifySuccess(t, root, outputRel, validatorRel, doctorRel)
	assertOperationsNotifyOutputSafe(t, root, outputRel, "", []string{"SECRET-RAW", "raw-secret-token", "postgres://user:pass@", "/tmp/private", `"raw_report"`, `"stdout"`, `"stderr"`, `"argv"`})
}

func TestOperationsNotifyOversizedSourcesAndStrictMode(t *testing.T) {
	root := operationsNotifyRepoRoot(t)
	base := operationsNotifyTempRel(t, root, "oversized")
	validatorRel := filepath.ToSlash(filepath.Join(base, "validator-health", "20260509T120000Z", "summary.json"))
	doctorRel := filepath.ToSlash(filepath.Join(base, "deployment-doctor", "20260509T120000Z", "summary.json"))
	writeOperationsNotifyRaw(t, root, validatorRel, []byte(`{"overall_status":"recorded","padding":"`+strings.Repeat("x", 2048)+`"}`))
	writeOperationsNotifyRaw(t, root, doctorRel, []byte(`{"overall_status":`+strings.Repeat("x", 2048)))
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	cmd := operationsNotifyCommand(root)
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"VALIDATOR_HEALTH_SUMMARY="+validatorRel,
		"DEPLOYMENT_DOCTOR_SUMMARY="+doctorRel,
		"MAX_SOURCE_BYTES=256",
		"FORCE=true",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("default oversized run should exit 0: %v\n%s", err, out)
	}
	summary := readOperationsNotifySummary(t, root, outputRel)
	if got := summary["notification"].(map[string]any)["severity"]; got != "needs_review" {
		t.Fatalf("oversized severity = %v, want needs_review", got)
	}
	strictOut := filepath.ToSlash(filepath.Join(base, "strict-out"))
	cmd = operationsNotifyCommand(root)
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+strictOut,
		"VALIDATOR_HEALTH_SUMMARY="+validatorRel,
		"DEPLOYMENT_DOCTOR_SUMMARY="+doctorRel,
		"MAX_SOURCE_BYTES=256",
		"STRICT_OPERATIONS_NOTIFY=true",
		"FORCE=true",
	)
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("strict oversized run should fail: %s", out)
	}
	assertOperationsNotifyOutputSafe(t, root, strictOut, string(out), []string{strings.Repeat("x", 64), "/tmp/", "/var/", "/etc/", "/Users/"})
}

func TestOperationsNotifyRedactionFailureAlwaysFails(t *testing.T) {
	root := operationsNotifyRepoRoot(t)
	base := operationsNotifyTempRel(t, root, "redaction-failure")
	validator := operationsNotifyValidatorSummary("blocked")
	validator.(map[string]any)["feeds"] = []map[string]any{
		{"feed_type": "schedule", "health_status": "blocked", "next_action": "Authorization: Bearer leaked-secret-token"},
	}
	validatorRel := filepath.ToSlash(filepath.Join(base, "validator-health", "20260509T120000Z", "summary.json"))
	doctorRel := filepath.ToSlash(filepath.Join(base, "deployment-doctor", "20260509T120000Z", "summary.json"))
	writeOperationsNotifyJSON(t, root, validatorRel, validator)
	writeOperationsNotifyJSON(t, root, doctorRel, operationsNotifyDoctorSummary("passed", nil))
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	for _, strict := range []bool{false, true} {
		cmd := operationsNotifyCommand(root)
		env := []string{
			"OUTPUT_DIR=" + outputRel + fmt.Sprint(strict),
			"VALIDATOR_HEALTH_SUMMARY=" + validatorRel,
			"DEPLOYMENT_DOCTOR_SUMMARY=" + doctorRel,
			"FORCE=true",
		}
		if strict {
			env = append(env, "STRICT_OPERATIONS_NOTIFY=true")
		}
		cmd.Env = append(os.Environ(), env...)
		if out, err := cmd.CombinedOutput(); err == nil {
			t.Fatalf("redaction failure strict=%v should fail: %s", strict, out)
		}
	}
}

func TestOperationsNotifyDocsWordingBoundaries(t *testing.T) {
	root := operationsNotifyRepoRoot(t)
	paths := []string{
		"docs/phase-47-self-hosted-operations-notifications.md",
		"docs/handoffs/phase-47.md",
		"docs/tutorials/self-hosted-operations-notifications.md",
		"docs/roadmap-to-calitp-compliance-and-gap-closure.md",
	}
	for _, rel := range paths {
		body, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		lower := strings.ToLower(string(body))
		for _, phrase := range []string{
			"production health proof",
			"cal-itp readiness proof",
			"consumer-readiness proof",
			"notification delivered",
			"webhook sent",
			"email sent",
			"evidence packet created",
		} {
			if strings.Contains(lower, phrase) && !strings.Contains(lower, "not "+phrase) {
				t.Fatalf("%s contains forbidden positive phrase %q", rel, phrase)
			}
		}
		if strings.Contains(lower, "compliance gate") && !strings.Contains(lower, "not a compliance gate") {
			t.Fatalf("%s contains forbidden positive compliance gate wording", rel)
		}
	}
}

func testOperationsNotifyLargeDeploymentDoctorSummary(t *testing.T, n int, blockers bool) {
	root := operationsNotifyRepoRoot(t)
	base := operationsNotifyTempRel(t, root, fmt.Sprintf("large-doctor-%d-%v", n, blockers))
	status := "warning"
	rowStatus := "warning"
	if blockers {
		status = "blocker"
		rowStatus = "blocker"
	}
	checks := make([]map[string]any, 0, n)
	for i := 0; i < n; i++ {
		checks = append(checks, map[string]any{
			"category": "env",
			"name":     fmt.Sprintf("CHECK_%05d", i),
			"status":   rowStatus,
			"detail":   "Resolve this local/reference diagnostic item before relying on notification summaries.",
		})
	}
	validatorRel := filepath.ToSlash(filepath.Join(base, "validator-health", "20260509T120000Z", "summary.json"))
	doctorRel := filepath.ToSlash(filepath.Join(base, "deployment-doctor", "20260509T120000Z", "summary.json"))
	writeOperationsNotifyJSON(t, root, validatorRel, operationsNotifyValidatorSummary("recorded"))
	writeOperationsNotifyJSON(t, root, doctorRel, operationsNotifyDoctorSummary(status, checks))
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	start := time.Now()
	runOperationsNotifySuccess(t, root, outputRel, validatorRel, doctorRel)
	t.Logf("operations notify large summary n=%d blockers=%v elapsed=%s", n, blockers, time.Since(start))
	summary := readOperationsNotifySummary(t, root, outputRel)
	actions := summary["next_actions"].([]any)
	if len(actions) > 20 {
		t.Fatalf("next_actions len = %d, want <=20", len(actions))
	}
	if got := int(summary["overflow_count"].(float64)); got < n-20 {
		t.Fatalf("overflow_count = %d, want at least %d", got, n-20)
	}
	assertOperationsNotifyOutputSafe(t, root, outputRel, "", []string{"raw_report", "SECRET=", "PASSWORD=", "Authorization: Bearer"})
}

func operationsNotifyRepoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func operationsNotifyCommand(root string, args ...string) *exec.Cmd {
	all := append([]string{filepath.Join(root, "scripts", "operations-notify.sh")}, args...)
	cmd := exec.Command("sh", all...)
	cmd.Dir = root
	return cmd
}

func operationsNotifyTempRel(t *testing.T, root, prefix string) string {
	t.Helper()
	base := filepath.Join(root, ".cache")
	if err := os.MkdirAll(base, 0o700); err != nil {
		t.Fatal(err)
	}
	dir, err := os.MkdirTemp(base, prefix+"-")
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

func operationsNotifyFakeNoSendBin(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range []string{"curl", "mail", "sendmail", "nc"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("#!/bin/sh\necho unexpected send/network command >&2\nexit 99\n"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func mkdirOperationsNotify(t *testing.T, root, rel string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, rel), 0o700); err != nil {
		t.Fatal(err)
	}
}

func writeOperationsNotifyJSON(t *testing.T, root, rel string, value any) {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	writeOperationsNotifyRaw(t, root, rel, append(body, '\n'))
}

func writeOperationsNotifyRaw(t *testing.T, root, rel string, body []byte) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
}

func runOperationsNotifySuccess(t *testing.T, root, outputRel, validatorRel, doctorRel string) {
	t.Helper()
	cmd := operationsNotifyCommand(root)
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"VALIDATOR_HEALTH_SUMMARY="+validatorRel,
		"DEPLOYMENT_DOCTOR_SUMMARY="+doctorRel,
		"FORCE=true",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("operations notify failed: %v\n%s", err, out)
	}
	assertOperationsNotifyOutputSafe(t, root, outputRel, string(out), nil)
}

func readOperationsNotifySummary(t *testing.T, root, outputRel string) map[string]any {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(root, outputRel, "summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("summary json invalid: %v\n%s", err, body)
	}
	return decoded
}

func assertOperationsNotifyOutputSafe(t *testing.T, root, outputRel, terminal string, extraForbidden []string) {
	t.Helper()
	outDir := filepath.Join(root, outputRel)
	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatal(err)
	}
	got := make([]string, 0, len(entries))
	for _, entry := range entries {
		got = append(got, entry.Name())
	}
	want := []string{"manifest.json", "manifest.md", "notification.txt", "summary.json", "summary.md"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("output files = %v, want %v", got, want)
	}
	var combined bytes.Buffer
	combined.WriteString(terminal)
	for _, name := range want {
		body, err := os.ReadFile(filepath.Join(outDir, name))
		if err != nil {
			t.Fatal(err)
		}
		combined.Write(body)
		if name == "summary.json" || name == "manifest.json" {
			var decoded any
			if err := json.Unmarshal(body, &decoded); err != nil {
				t.Fatalf("%s invalid JSON: %v", name, err)
			}
			if name == "manifest.json" {
				manifest := decoded.(map[string]any)
				for _, item := range manifest["source_files"].([]any) {
					source := item.(string)
					if source != "<redacted-source>" && !strings.HasPrefix(source, ".cache/") {
						t.Fatalf("manifest source file is not repo-relative .cache or redacted: %q", source)
					}
				}
			}
		}
		switch name {
		case "notification.txt":
			if !bytes.HasPrefix(body, []byte("DRAFT — NOT SENT")) {
				t.Fatalf("notification.txt missing required prefix: %s", body[:min(len(body), 80)])
			}
			if len(body) > 16*1024 {
				t.Fatalf("notification.txt size = %d", len(body))
			}
		case "summary.md":
			if len(body) > 24*1024 {
				t.Fatalf("summary.md size = %d", len(body))
			}
		case "summary.json":
			if len(body) > 64*1024 {
				t.Fatalf("summary.json size = %d", len(body))
			}
		}
	}
	for _, forbidden := range append([]string{
		"Authorization: Bearer",
		"admin_session=",
		"Cookie:",
		"DATABASE_URL=",
		"postgres://user:pass@",
		"TOKEN=",
		"SECRET=",
		"PASSWORD=",
		"BEGIN PRIVATE KEY",
		"/Users/",
		"/tmp/private",
		"webhook secret",
	}, extraForbidden...) {
		if forbidden != "" && strings.Contains(combined.String(), forbidden) {
			t.Fatalf("operations-notify output leaked %q\n%s", forbidden, combined.String())
		}
	}
}

func operationsNotifyValidatorSummary(status string) any {
	return map[string]any{
		"generated_at":   "2026-05-09T12:00:00Z",
		"overall_status": status,
		"feeds": []map[string]any{
			{"feed_type": "schedule", "health_status": status, "next_action": "Review validator-health private diagnostics."},
			{"feed_type": "vehicle_positions", "health_status": status, "next_action": "Review validator-health private diagnostics."},
			{"feed_type": "trip_updates", "health_status": status, "next_action": "Review validator-health private diagnostics."},
			{"feed_type": "alerts", "health_status": status, "next_action": "Review validator-health private diagnostics."},
		},
		"external_evidence_created":    false,
		"consumer_statuses_changed":    false,
		"compliance_claimed":           false,
		"production_readiness_claimed": false,
	}
}

func operationsNotifyDoctorSummary(status string, checks []map[string]any) map[string]any {
	if checks == nil {
		checks = []map[string]any{}
	}
	counts := map[string]any{"passed": 0, "blocker": 0, "warning": 0, "skipped": 0, "unavailable": 0}
	if len(checks) == 0 {
		counts[status] = 1
	} else {
		for _, check := range checks {
			key := fmt.Sprint(check["status"])
			counts[key] = int(counts[key].(int)) + 1
		}
	}
	return map[string]any{
		"generated_at_utc":             "20260509T120000Z",
		"overall_status":               status,
		"counts":                       counts,
		"checks":                       checks,
		"external_evidence_created":    false,
		"final_root_evidence_created":  false,
		"consumer_statuses_changed":    false,
		"compliance_claimed":           false,
		"production_readiness_claimed": false,
	}
}

func summaryBool(t *testing.T, summary map[string]any, key string, want bool) bool {
	t.Helper()
	got, ok := summary[key].(bool)
	if !ok {
		t.Fatalf("summary[%s] missing bool: %+v", key, summary[key])
	}
	return got == want
}
