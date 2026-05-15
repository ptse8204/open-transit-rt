package realtimequality

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadBacktestBrowserMissingRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".cache", "realtime-quality-backtest")
	browser, err := LoadBacktestBrowser(root, 5)
	if err != nil {
		t.Fatalf("load missing root: %v", err)
	}
	if browser.Status != "missing" || browser.RootRef != DefaultBacktestBrowserRoot || len(browser.Rows) != 0 || browser.Boundary == "" {
		t.Fatalf("unexpected missing-root browser: %+v", browser)
	}
}

func TestLoadBacktestBrowserReadsAggregateOnlySummaries(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".cache", "realtime-quality-backtest")
	writeBacktestBrowserFixture(t, root, "20260514T120000Z", false)

	browser, err := LoadBacktestBrowser(root, 5)
	if err != nil {
		t.Fatalf("load browser: %v", err)
	}
	if browser.Status != "needs_review" || len(browser.Rows) != 1 {
		t.Fatalf("browser = %+v, want one needs_review row", browser)
	}
	row := browser.Rows[0]
	if row.OutputRef != ".cache/realtime-quality-backtest/20260514T120000Z" || row.MaturityGate != "diagnostic_watch" || row.PredictionCoverage != "62.5% (5/8)" || row.MAEAbsoluteErrorSeconds != "29 sec" || !strings.Contains(row.WithheldByReason, "manual_override_review=1") || row.ConformanceSignal != "synthetic_covered (5/5 synthetic cases)" {
		t.Fatalf("unexpected browser row: %+v", row)
	}
	raw, err := json.Marshal(browser)
	if err != nil {
		t.Fatalf("marshal browser: %v", err)
	}
	for _, forbidden := range []string{root, filepath.Dir(root), "/Users/", "raw observed", "raw prediction", "authorization:", "postgres://"} {
		if strings.Contains(strings.ToLower(string(raw)), strings.ToLower(forbidden)) {
			t.Fatalf("browser output leaked forbidden string %q: %s", forbidden, raw)
		}
	}
}

func TestLoadBacktestBrowserBlocksUnsafeRootAndClaimFlags(t *testing.T) {
	if _, err := LoadBacktestBrowser(filepath.Join(t.TempDir(), "docs", "evidence", "realtime-quality-backtest"), 5); err == nil {
		t.Fatalf("unsafe evidence-like root was accepted")
	}

	root := filepath.Join(t.TempDir(), ".cache", "realtime-quality-backtest")
	writeBacktestBrowserFixture(t, root, "20260514T130000Z", true)
	browser, err := LoadBacktestBrowser(root, 5)
	if err != nil {
		t.Fatalf("load browser with blocked fixture: %v", err)
	}
	if browser.Status != "blocked" || len(browser.Rows) != 1 || browser.Rows[0].Status != "blocked" || browser.Rows[0].GeneratedAt != "not loaded" {
		t.Fatalf("unsafe fixture was not blocked: %+v", browser)
	}
}

func writeBacktestBrowserFixture(t *testing.T, root, name string, forbiddenClaim bool) {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create fixture dir: %v", err)
	}
	generatedAt := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	coverage := 62.5
	futureCoverage := 50.0
	mae := 29.0
	p90 := 40.0
	overall := MetricGroup{
		GroupType:                   "overall",
		CoverageDenominator:         8,
		CoverageNumerator:           5,
		FutureStopCoverageNumerator: 4,
		MatchedPredictionCount:      5,
		MissingPredictionCount:      1,
		MissingObservationCount:     1,
		StalePredictionCount:        1,
		WithheldByReason:            map[string]int{"manual_override_review": 1},
		MAEAbsoluteErrorSeconds:     &mae,
		P90AbsoluteErrorSeconds:     &p90,
		PredictionCoverage:          Rate{Numerator: 5, Denominator: 8, Status: "measured", Percent: &coverage},
		FutureStopCoverage:          Rate{Numerator: 4, Denominator: 8, Status: "measured", Percent: &futureCoverage},
		MaturityGate:                "diagnostic_watch",
	}
	summary := SummaryDocument{
		SchemaVersion:             BacktestSchemaVersion,
		GeneratedAt:               generatedAt,
		InputRecordCounts:         RecordCounts{ObservedRecords: 8, PredictionRecords: 7},
		Overall:                   overall,
		Conformance:               coveredBrowserConformanceReview(),
		GroupCount:                1,
		ProductionGradeETAClaimed: forbiddenClaim,
	}
	metrics := MetricsDocument{SchemaVersion: BacktestSchemaVersion, GeneratedAt: generatedAt, Groups: []MetricGroup{overall}}
	manifest := ManifestDocument{
		SchemaVersion: BacktestSchemaVersion,
		GeneratedAt:   generatedAt,
		OutputKind:    "private_local_realtime_quality_backtest",
		OutputFiles:   ExpectedBacktestOutputFiles(),
		SafetyChecks: map[string]bool{
			"docs_evidence_output_rejected": true,
			"evidence_like_output_rejected": true,
			"symlink_ancestors_rejected":    true,
			"raw_inputs_not_copied":         true,
			"private_paths_omitted":         true,
			"raw_rows_omitted":              true,
		},
		Boundaries: map[string]bool{
			"db_persistence":                   false,
			"migration_added":                  false,
			"operations_console_change":        false,
			"public_api_added":                 false,
			"consumer_tracker_changed":         false,
			"external_predictor_runtime_added": false,
		},
		AggregateOnly:    true,
		RawRowsPersisted: false,
	}
	writeBacktestBrowserJSON(t, filepath.Join(dir, "summary.json"), summary)
	writeBacktestBrowserJSON(t, filepath.Join(dir, "metrics.json"), metrics)
	writeBacktestBrowserJSON(t, filepath.Join(dir, "manifest.json"), manifest)
	for _, name := range []string{"summary.md", "metrics.md"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("private aggregate diagnostic summary\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
}

func coveredBrowserConformanceReview() ConformanceReview {
	cases := []ConformanceCase{
		{ID: "after-midnight-service", Status: "synthetic_covered", Signal: "after-midnight covered", DoesNotProve: "synthetic only"},
		{ID: "frequency-headway-service", Status: "synthetic_covered", Signal: "frequency covered", DoesNotProve: "synthetic only"},
		{ID: "service-calendar-start-instance", Status: "synthetic_covered", Signal: "service calendar covered", DoesNotProve: "synthetic only"},
		{ID: "blocked-unknown-ambiguous", Status: "synthetic_covered", Signal: "withheld covered", DoesNotProve: "synthetic only"},
		{ID: "shadow-fail-closed", Status: "synthetic_covered", Signal: "shadow covered", DoesNotProve: "synthetic only"},
	}
	return ConformanceReview{
		Status:        "synthetic_covered",
		Boundary:      "Private synthetic conformance summary only.",
		SyntheticOnly: true,
		AggregateOnly: true,
		CaseCount:     len(cases),
		Cases:         cases,
		DoesNotProve:  "Synthetic conformance rows do not prove real-world ETA accuracy.",
	}
}

func writeBacktestBrowserJSON(t *testing.T, path string, value any) {
	t.Helper()
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatalf("write fixture json: %v", err)
	}
}
