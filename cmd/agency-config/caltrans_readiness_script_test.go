package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCaltransReadinessCheckDryRunExactFilesFlagsAndStatuses(t *testing.T) {
	root := caltransReadinessRepoRoot(t)
	outputRel := caltransReadinessTempRel(t, root, "dry-run")
	cmd := caltransReadinessCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(), "OUTPUT_DIR="+outputRel, "FORCE=true")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("caltrans readiness dry-run failed: %v\n%s", err, out)
	}
	assertCaltransReadinessExactFiles(t, root, outputRel)
	assertCaltransReadinessNoLeaks(t, root, outputRel, string(out))
	summary := readCaltransReadinessSummary(t, root, outputRel)
	if summary["dry_run"] != true {
		t.Fatalf("dry_run = %v, want true", summary["dry_run"])
	}
	assertCaltransReadinessFlagsFalse(t, summary)
	assertCaltransReadinessNoClaimStatuses(t, summary)
	assertCaltransReadinessConsumerPreparedOnly(t, summary)
	if status := caltransReadinessRowStatus(summary, "public_fetchability"); status != "not_checked" {
		t.Fatalf("public_fetchability status = %q, want not_checked", status)
	}
	if status := caltransReadinessRowStatus(summary, "stable_urls"); status != "missing" {
		t.Fatalf("stable_urls status = %q, want missing without inputs", status)
	}
}

func TestCaltransReadinessCheckWithSafeFeedsSummary(t *testing.T) {
	root := caltransReadinessRepoRoot(t)
	base := caltransReadinessTempRel(t, root, "feeds")
	feedsRel := filepath.ToSlash(filepath.Join(base, "inputs", "feeds.json"))
	writeCaltransReadinessJSON(t, root, feedsRel, map[string]any{
		"public_base_url":         "https://feeds.example.org",
		"technical_contact_email": "ops@example.org",
		"license":                 map[string]any{"name": "CC BY 4.0", "url": "https://example.org/license"},
		"publication_environment": "pilot",
		"agency_id":               "demo-agency",
		"agency_name":             "Demo Agency",
		"generated_at":            "2026-05-10T12:00:00Z",
		"readiness":               map[string]any{"all_required_feeds_listed": true, "https_urls": true, "license_complete": true, "contact_complete": true},
		"feeds": []map[string]any{
			{"feed_type": "schedule", "canonical_public_url": "https://feeds.example.org/public/gtfs/schedule.zip"},
			{"feed_type": "vehicle_positions", "canonical_public_url": "https://feeds.example.org/public/gtfsrt/vehicle_positions.pb"},
			{"feed_type": "trip_updates", "canonical_public_url": "https://feeds.example.org/public/gtfsrt/trip_updates.pb"},
			{"feed_type": "alerts", "canonical_public_url": "https://feeds.example.org/public/gtfsrt/alerts.pb"},
		},
	})
	validatorRel := filepath.ToSlash(filepath.Join(base, "inputs", "validator-health.json"))
	reliabilityRel := filepath.ToSlash(filepath.Join(base, "inputs", "operations-reliability.json"))
	tripRel := filepath.ToSlash(filepath.Join(base, "inputs", "trip-id-consistency.json"))
	writeCaltransReadinessJSON(t, root, validatorRel, map[string]any{"overall_status": "recorded"})
	writeCaltransReadinessJSON(t, root, reliabilityRel, map[string]any{"overall_status": "recorded"})
	writeCaltransReadinessJSON(t, root, tripRel, map[string]any{"status": "recorded"})
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	cmd := caltransReadinessCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+outputRel,
		"FORCE=true",
		"FEEDS_JSON_PATH="+feedsRel,
		"VALIDATOR_HEALTH_SUMMARY="+validatorRel,
		"OPERATIONS_RELIABILITY_SUMMARY="+reliabilityRel,
		"TRIP_ID_CONSISTENCY_SUMMARY="+tripRel,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("safe summary run failed: %v\n%s", err, out)
	}
	summary := readCaltransReadinessSummary(t, root, outputRel)
	for _, area := range []string{"stable_urls", "static_gtfs", "vehicle_positions", "trip_updates", "alerts", "https", "open_license", "contact", "validators", "freshness", "trip_id_consistency", "consumer_packet_preparedness"} {
		if status := caltransReadinessRowStatus(summary, area); status != "present" {
			t.Fatalf("%s status = %q, want present", area, status)
		}
	}
	if status := caltransReadinessRowStatus(summary, "public_fetchability"); status != "not_checked" {
		t.Fatalf("public_fetchability status = %q, want not_checked when fetch disabled", status)
	}
	assertCaltransReadinessFlagsFalse(t, summary)
	assertCaltransReadinessNoClaimStatuses(t, summary)
}

func TestCaltransReadinessCheckNonHTTPSRootIsNotPresent(t *testing.T) {
	root := caltransReadinessRepoRoot(t)
	outputRel := caltransReadinessTempRel(t, root, "http-root")
	cmd := caltransReadinessCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(), "OUTPUT_DIR="+outputRel, "FORCE=true", "PUBLIC_BASE_URL=http://localhost:8080")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("http root run failed: %v\n%s", err, out)
	}
	summary := readCaltransReadinessSummary(t, root, outputRel)
	if status := caltransReadinessRowStatus(summary, "stable_urls"); status != "present" {
		t.Fatalf("stable_urls status = %q, want present because URLs are configured", status)
	}
	if status := caltransReadinessRowStatus(summary, "https"); status == "present" {
		t.Fatalf("https status = present for http root; summary=%+v", summary)
	}
	if status := caltransReadinessRowStatus(summary, "public_fetchability"); status != "not_checked" {
		t.Fatalf("public_fetchability status = %q, want not_checked by default", status)
	}
	assertCaltransReadinessNoClaimStatuses(t, summary)
}

func TestCaltransReadinessCheckRejectsEvidenceOutputAndUnsafeSource(t *testing.T) {
	root := caltransReadinessRepoRoot(t)
	cmd := caltransReadinessCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(), "OUTPUT_DIR=docs/evidence/caltrans-readiness-check")
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("expected evidence output rejection, got success: %s", out)
	}

	base := caltransReadinessTempRel(t, root, "unsafe-source")
	feedsRel := filepath.ToSlash(filepath.Join(base, "inputs", "feeds.json"))
	writeCaltransReadinessRaw(t, root, feedsRel, []byte(`{"public_base_url":"https://feeds.example.org","database_url":"postgres://user:pass@localhost/db"}`))
	outputRel := filepath.ToSlash(filepath.Join(base, "out"))
	cmd = caltransReadinessCommand(root, "--dry-run")
	cmd.Env = append(os.Environ(), "OUTPUT_DIR="+outputRel, "FORCE=true", "FEEDS_JSON_PATH="+feedsRel)
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("expected unsafe source rejection, got success: %s", out)
	}
}

func TestCaltransReadinessCheckHelpAndDocBoundary(t *testing.T) {
	root := caltransReadinessRepoRoot(t)
	cmd := caltransReadinessCommand(root, "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help failed: %v\n%s", err, out)
	}
	doc, err := os.ReadFile(filepath.Join(root, "docs", "caltrans-readiness-gap-report.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(out) + string(doc)
	for _, want := range []string{
		"gap diagnostics",
		"does not refresh official requirements",
		"does not create retained evidence",
		"does not claim CAL-ITP/Caltrans compliance",
		"consumer acceptance",
		"production readiness",
		"vendor compatibility",
		"production-grade ETA",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing boundary phrase %q", want)
		}
	}
}

func caltransReadinessCommand(root string, args ...string) *exec.Cmd {
	all := append([]string{filepath.Join(root, "scripts", "caltrans-readiness-check.sh")}, args...)
	cmd := exec.Command("sh", all...)
	cmd.Dir = root
	return cmd
}

func caltransReadinessRepoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func caltransReadinessTempRel(t *testing.T, root, name string) string {
	t.Helper()
	rel := filepath.ToSlash(filepath.Join(".cache", "caltrans-readiness-check-test", name, strings.ReplaceAll(t.Name(), "/", "-")))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(root, rel)) })
	return rel
}

func writeCaltransReadinessJSON(t *testing.T, root, rel string, payload any) {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	writeCaltransReadinessRaw(t, root, rel, body)
}

func writeCaltransReadinessRaw(t *testing.T, root, rel string, body []byte) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
}

func readCaltransReadinessSummary(t *testing.T, root, outputRel string) map[string]any {
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

func assertCaltransReadinessExactFiles(t *testing.T, root, outputRel string) {
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
	want := "gap-review.txt,manifest.json,manifest.md,summary.json,summary.md"
	if got != want {
		t.Fatalf("files = %s, want %s", got, want)
	}
}

func assertCaltransReadinessNoLeaks(t *testing.T, root, outputRel, stdout string) {
	t.Helper()
	combined := stdout
	for _, name := range []string{"summary.json", "summary.md", "manifest.json", "manifest.md", "gap-review.txt"} {
		body, err := os.ReadFile(filepath.Join(root, outputRel, name))
		if err != nil {
			t.Fatal(err)
		}
		combined += "\n" + string(body)
	}
	for _, forbidden := range []string{"Authorization:", "Bearer ", "DATABASE_URL", "postgres://", "BEGIN PRIVATE KEY", "secret-token", "https://hooks.", "/Users/"} {
		if strings.Contains(combined, forbidden) {
			t.Fatalf("caltrans readiness output leaked %q\n%s", forbidden, combined)
		}
	}
}

func assertCaltransReadinessFlagsFalse(t *testing.T, summary map[string]any) {
	t.Helper()
	flags := summary["claim_flags"].(map[string]any)
	for key, value := range flags {
		if value != false {
			t.Fatalf("claim flag %s = %v, want false", key, value)
		}
	}
}

func assertCaltransReadinessNoClaimStatuses(t *testing.T, summary map[string]any) {
	t.Helper()
	for _, raw := range summary["rows"].([]any) {
		row := raw.(map[string]any)
		status := row["status"].(string)
		switch status {
		case "ok", "passed", "compliant", "certified", "accepted", "ingested", "listed", "displayed":
			t.Fatalf("claim-upgrading status emitted: %+v", row)
		case "present", "partial", "missing", "not_checked", "needs_review", "blocked":
		default:
			t.Fatalf("unexpected row status %q in %+v", status, row)
		}
	}
}

func assertCaltransReadinessConsumerPreparedOnly(t *testing.T, summary map[string]any) {
	t.Helper()
	tracker := summary["consumer_tracker"].(map[string]any)
	if tracker["prepared_only"] != true {
		t.Fatalf("consumer tracker prepared_only = %v, want true", tracker["prepared_only"])
	}
	targets := tracker["expected_targets"].([]any)
	got := make([]string, 0, len(targets))
	for _, target := range targets {
		got = append(got, target.(string))
	}
	want := []string{"Google Maps", "Apple Maps", "Transit App", "Bing Maps", "Moovit", "Mobility Database", "transit.land"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("expected targets = %v, want %v", got, want)
	}
}

func caltransReadinessRowStatus(summary map[string]any, area string) string {
	for _, raw := range summary["rows"].([]any) {
		row := raw.(map[string]any)
		if row["area"] == area {
			return row["status"].(string)
		}
	}
	return ""
}
