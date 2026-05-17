package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAdapterConformanceCommandsPassSyntheticSuite(t *testing.T) {
	for _, args := range [][]string{
		{"help"},
		{"manifest", "--suite", "testdata/adapter-conformance"},
		{"telemetry", "--suite", "testdata/adapter-conformance"},
		{"prediction", "--suite", "testdata/adapter-conformance"},
		{"validator", "--suite", "testdata/adapter-conformance"},
		{"monitoring", "--suite", "testdata/adapter-conformance"},
		{"consumer_discovery", "--suite", "testdata/adapter-conformance"},
		{"run", "--suite", "testdata/adapter-conformance"},
	} {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(args, &stdout, &stderr); code != 0 {
				t.Fatalf("run(%v) code=%d stderr=%s stdout=%s", args, code, stderr.String(), stdout.String())
			}
		})
	}
}

func TestAdapterConformanceV2ScenarioCoverage(t *testing.T) {
	root := repoRootForTest(t)
	s, err := loadSuite(filepath.Join(root, "testdata", "adapter-conformance"))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]string{
		"telemetry":          {"malformed", "stale", "future", "wrong-agency", "unknown-device", "low-quality", "duplicate", "out-of-order", "missing-required-field", "invalid-coordinate"},
		"prediction":         {"timeout", "malformed", "stale", "wrong-agency", "low-confidence", "missing-vehicle-positions-ref", "public-mutation-attempt"},
		"validator":          {"allowlist", "raw-command"},
		"monitoring":         {"redaction", "no-send", "unredacted-destination"},
		"consumer_discovery": {"feed-url-metadata", "status-mutation-blocked", "submission-automation-blocked"},
	}
	for caseType, scenarios := range want {
		if !sameStrings(requiredScenarios[caseType], scenarios) {
			t.Fatalf("requiredScenarios[%s] = %v, want %v", caseType, requiredScenarios[caseType], scenarios)
		}
		seen := map[string]testCase{}
		for _, tc := range s.Cases {
			if tc.Type == caseType {
				seen[tc.Scenario] = tc
			}
		}
		for _, scenario := range scenarios {
			tc, ok := seen[scenario]
			if !ok {
				t.Fatalf("missing %s scenario %s", caseType, scenario)
			}
			if !tc.SyntheticOnly {
				t.Fatalf("case %s must be synthetic_only", tc.ID)
			}
			if err := requiredAssertions(tc); err != nil {
				t.Fatalf("case %s required assertions: %v", tc.ID, err)
			}
		}
	}
}

func TestAdapterConformanceRejectsMissingScenario(t *testing.T) {
	root := repoRootForTest(t)
	tmp, tmpRel := tempSuiteDir(t, root, "missing-scenario")
	copyFile(t, filepath.Join(root, "testdata", "adapter-conformance", "suite.json"), filepath.Join(tmp, "suite.json"))
	copyDir(t, filepath.Join(root, "testdata", "adapter-conformance", "fixtures"), filepath.Join(tmp, "fixtures"))
	body, err := os.ReadFile(filepath.Join(tmp, "suite.json"))
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.Replace(string(body), `"scenario": "stale"`, `"scenario": "malformed"`, 1)
	if err := os.WriteFile(filepath.Join(tmp, "suite.json"), []byte(trimmed), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"telemetry", "--suite", tmpRel}, &stdout, &stderr)
	if code == 0 || !strings.Contains(stderr.String(), "missing required telemetry scenario stale") {
		t.Fatalf("expected missing stale scenario failure, code=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
}

func sameStrings(got []string, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func TestAdapterConformanceRejectsEvidenceSuitePath(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"run", "--suite", "docs/evidence/adapter-conformance"}, &stdout, &stderr)
	if code == 0 || !strings.Contains(stderr.String(), "evidence") {
		t.Fatalf("expected evidence path rejection, code=%d stderr=%s", code, stderr.String())
	}
}

func TestAdapterConformanceRejectsUnsafeFixtureText(t *testing.T) {
	root := repoRootForTest(t)
	tmp, tmpRel := tempSuiteDir(t, root, "unsafe-fixture")
	fixtures := filepath.Join(tmp, "fixtures")
	copyFile(t, filepath.Join(root, "testdata", "adapter-conformance", "suite.json"), filepath.Join(tmp, "suite.json"))
	copyDir(t, filepath.Join(root, "testdata", "adapter-conformance", "fixtures"), fixtures)
	if err := os.WriteFile(filepath.Join(fixtures, "telemetry-stale.json"), []byte(`{"synthetic_only":true,"note":"Authorization: Bearer token"}`+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"telemetry", "--suite", tmpRel}, &stdout, &stderr)
	if code == 0 || !strings.Contains(stderr.String(), "forbidden text") {
		t.Fatalf("expected unsafe fixture rejection, code=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
}

func repoRootForTest(t *testing.T) string {
	t.Helper()
	root, err := repoRoot()
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func tempSuiteDir(t *testing.T, root string, name string) (string, string) {
	t.Helper()
	dir := filepath.Join(root, ".cache", "adapter-conformance-test", name, strings.ReplaceAll(t.Name(), "/", "-"))
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return dir, relForTest(t, root, dir)
}

func relForTest(t *testing.T, root string, path string) string {
	t.Helper()
	rel, err := filepath.Rel(root, path)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(rel)
}

func copyDir(t *testing.T, src string, dst string) {
	t.Helper()
	if err := os.MkdirAll(dst, 0o700); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		copyFile(t, filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name()))
	}
}

func copyFile(t *testing.T, src string, dst string) {
	t.Helper()
	body, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, body, 0o600); err != nil {
		t.Fatal(err)
	}
}
