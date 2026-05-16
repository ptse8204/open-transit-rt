package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidatePublicFeedsScriptDryRunSummaryFlagsAndNoLeaks(t *testing.T) {
	root := referenceCheckRepoRoot(t)
	outputRel := filepath.ToSlash(filepath.Join(referenceCheckTempRel(t, root, "validate-public-feeds"), "out"))
	cmd := exec.Command("sh", filepath.Join(root, "scripts", "validate-public-feeds.sh"), "--public-base-url", "https://feeds.example.org", "--dry-run", "--output-dir", outputRel)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "FORCE=true", "ADMIN_TOKEN=SECRET-RAW", "DATABASE_URL=postgres://user:pass@example/db")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("validate-public-feeds dry-run failed: %v\n%s", err, out)
	}
	summary := readReferenceCheckJSON(t, root, outputRel, "summary.json")
	rows := summary["rows"].([]any)
	if len(rows) != 5 {
		t.Fatalf("rows = %d, want 5: %+v", len(rows), rows)
	}
	wantPaths := []string{"/public/feeds.json", "/public/gtfs/schedule.zip", "/public/gtfsrt/vehicle_positions.pb", "/public/gtfsrt/trip_updates.pb", "/public/gtfsrt/alerts.pb"}
	for i, row := range rows {
		got := row.(map[string]any)["public_path"]
		if got != wantPaths[i] {
			t.Fatalf("row %d public_path = %v, want %s", i, got, wantPaths[i])
		}
	}
	assertReferenceCheckFlagsFalse(t, summary["claim_flags"].(map[string]any))
	assertReferenceCheckNoLeakage(t, root, outputRel, string(out))
}

func TestOCIReferenceCheckScriptDryRunSummaryFlagsAndNoLeaks(t *testing.T) {
	root := referenceCheckRepoRoot(t)
	outputRel := filepath.ToSlash(filepath.Join(referenceCheckTempRel(t, root, "oci-reference-check"), "out"))
	cmd := exec.Command("sh", filepath.Join(root, "scripts", "oci-reference-check.sh"), "--public-base-url", "https://feeds.example.org", "--dry-run", "--output-dir", outputRel)
	cmd.Dir = root
	cmd.Env = append(os.Environ(),
		"FORCE=true",
		"BACKUP_DIR=/private/secret-backup",
		"RESTORE_DRILL_DATABASE_URL=postgres://user:pass@example/db",
		"DEVICE_TOKEN=Bearer raw-secret-token",
		"OCI_HOST=192.0.2.10",
		"OCI_KEY=/Users/example/.ssh/private-key",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("oci-reference-check dry-run failed: %v\n%s", err, out)
	}
	summary := readReferenceCheckJSON(t, root, outputRel, "summary.json")
	assertReferenceCheckFlagsFalse(t, summary["claim_flags"].(map[string]any))
	if len(summary["loopback_health"].([]any)) != 5 {
		t.Fatalf("loopback rows = %+v, want 5 rows", summary["loopback_health"])
	}
	publicSummary := summary["public_feed_summary"].(map[string]any)
	if len(publicSummary["rows"].([]any)) != 5 {
		t.Fatalf("public feed rows = %+v, want 5 rows", publicSummary["rows"])
	}
	assertReferenceCheckNoLeakage(t, root, outputRel, string(out))
}

func TestReferenceCheckScriptsRejectEvidenceOutputAndSecretURLs(t *testing.T) {
	root := referenceCheckRepoRoot(t)
	for _, tc := range []struct {
		name string
		args []string
	}{
		{name: "evidence_output", args: []string{"--output-dir", "docs/evidence/reference-check-test", "--dry-run"}},
		{name: "userinfo_url", args: []string{"--public-base-url", "https://user:secret@feeds.example.org", "--dry-run"}},
	} {
		t.Run("validate_"+tc.name, func(t *testing.T) {
			cmd := exec.Command("sh", append([]string{filepath.Join(root, "scripts", "validate-public-feeds.sh")}, tc.args...)...)
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "FORCE=true", "ALLOW_UNIGNORED_OUTPUT_DIR=true")
			if out, err := cmd.CombinedOutput(); err == nil {
				t.Fatalf("expected validate-public-feeds rejection, got success: %s", out)
			}
		})
		t.Run("oci_"+tc.name, func(t *testing.T) {
			cmd := exec.Command("sh", append([]string{filepath.Join(root, "scripts", "oci-reference-check.sh")}, tc.args...)...)
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "FORCE=true", "ALLOW_UNIGNORED_OUTPUT_DIR=true")
			if out, err := cmd.CombinedOutput(); err == nil {
				t.Fatalf("expected oci-reference-check rejection, got success: %s", out)
			}
		})
	}
}

func referenceCheckRepoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func referenceCheckTempRel(t *testing.T, root, name string) string {
	t.Helper()
	base := filepath.Join(root, ".cache", "reference-check-test")
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

func readReferenceCheckJSON(t *testing.T, root, outputRel, name string) map[string]any {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(root, outputRel, name))
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func assertReferenceCheckFlagsFalse(t *testing.T, flags map[string]any) {
	t.Helper()
	for key, value := range flags {
		if value != false {
			t.Fatalf("claim flag %s = %+v, want false", key, value)
		}
	}
}

func assertReferenceCheckNoLeakage(t *testing.T, root, outputRel, stdout string) {
	t.Helper()
	combined := stdout
	filepath.WalkDir(filepath.Join(root, outputRel), func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		body, readErr := os.ReadFile(path)
		if readErr == nil {
			combined += string(body)
		}
		return nil
	})
	for _, forbidden := range []string{
		"SECRET-RAW",
		"raw-secret-token",
		"postgres://user:pass",
		"/private/secret",
		"/Users/example",
		"Authorization:",
		"Bearer raw",
		"BEGIN PRIVATE KEY",
		"consumer accepted",
		"production ready",
		"compliance achieved",
	} {
		if strings.Contains(strings.ToLower(combined), strings.ToLower(forbidden)) {
			t.Fatalf("reference check output leaked or overclaimed %q\n%s", forbidden, combined)
		}
	}
}
