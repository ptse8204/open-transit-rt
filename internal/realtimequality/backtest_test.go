package realtimequality

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestBacktestFixtureMetricsAndGateLabels(t *testing.T) {
	observed, observedRaw, err := LoadObservedDataset(filepath.Join("..", "..", "testdata", "realtime-quality-backtest", "observed-events.json"))
	if err != nil {
		t.Fatalf("load observed: %v", err)
	}
	predictions, predictionsRaw, err := LoadPredictionDataset(filepath.Join("..", "..", "testdata", "realtime-quality-backtest", "prediction-samples.json"))
	if err != nil {
		t.Fatalf("load predictions: %v", err)
	}
	report, err := RunBacktest(observed, predictions, observedRaw, predictionsRaw, BacktestOptions{
		GeneratedAt: time.Date(2026, 5, 9, 20, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("run backtest: %v", err)
	}
	overall := report.Summary.Overall
	if overall.CoverageDenominator != 11 || overall.CoverageNumerator != 5 {
		t.Fatalf("overall coverage = %d/%d, want 5/11", overall.CoverageNumerator, overall.CoverageDenominator)
	}
	if overall.MatchedPredictionCount != 5 || overall.MissingPredictionCount != 1 || overall.StalePredictionCount != 1 || overall.MissingObservationCount != 1 {
		t.Fatalf("overall counts = %+v, want matched/missing/stale/missing-observation visibility", overall)
	}
	for _, reason := range []string{"manual_override_review", "unknown_assignment", "ambiguous_assignment", "external_predictor_fail_closed"} {
		if overall.WithheldByReason[reason] != 1 {
			t.Fatalf("withheld_by_reason = %+v, want %s", overall.WithheldByReason, reason)
		}
	}
	if overall.MAEAbsoluteErrorSeconds == nil || *overall.MAEAbsoluteErrorSeconds != 29 {
		t.Fatalf("mae = %v, want 29.0 seconds", overall.MAEAbsoluteErrorSeconds)
	}
	if overall.MedianAbsoluteErrorSeconds == nil || *overall.MedianAbsoluteErrorSeconds != 30 {
		t.Fatalf("median = %v, want 30.0 seconds", overall.MedianAbsoluteErrorSeconds)
	}
	if overall.P90AbsoluteErrorSeconds == nil || *overall.P90AbsoluteErrorSeconds != 40 {
		t.Fatalf("p90 = %v, want 40.0 seconds", overall.P90AbsoluteErrorSeconds)
	}
	if overall.MeanLeadTimeSeconds == nil || *overall.MeanLeadTimeSeconds != 120 {
		t.Fatalf("mean lead time = %v, want 120.0 seconds", overall.MeanLeadTimeSeconds)
	}
	if overall.MedianLeadTimeSeconds == nil || *overall.MedianLeadTimeSeconds != 120 {
		t.Fatalf("median lead time = %v, want 120.0 seconds", overall.MedianLeadTimeSeconds)
	}
	if overall.P90LeadTimeSeconds == nil || *overall.P90LeadTimeSeconds != 150 {
		t.Fatalf("p90 lead time = %v, want 150.0 seconds", overall.P90LeadTimeSeconds)
	}
	if overall.PredictionCoverage.Percent == nil || *overall.PredictionCoverage.Percent != 45.5 {
		t.Fatalf("coverage percent = %+v, want 45.5", overall.PredictionCoverage)
	}
	if report.Summary.Conformance.Status != "synthetic_covered" || !report.Summary.Conformance.SyntheticOnly || !report.Summary.Conformance.AggregateOnly || len(report.Summary.Conformance.Cases) != 5 {
		t.Fatalf("conformance summary = %+v, want five covered synthetic cases", report.Summary.Conformance)
	}
	for _, key := range []string{"after_midnight", "frequency", "unknown_assignment", "ambiguous_assignment", "external_predictor_fail_closed", "shadow_adapter_samples"} {
		if report.Summary.Conformance.CaseCounts[key] <= 0 {
			t.Fatalf("conformance case counts = %+v, want positive %s", report.Summary.Conformance.CaseCounts, key)
		}
	}
	for _, group := range report.Metrics.Groups {
		if !allowedGateLabels[group.MaturityGate] {
			t.Fatalf("forbidden maturity gate label in group %+v", group)
		}
	}
	assertGroupPresent(t, report.Metrics.Groups, "time_period", "", "overnight")
	assertGroupPresent(t, report.Metrics.Groups, "route", "route-frequency", "")
	assertGroupPresent(t, report.Metrics.Groups, "route_time_period", "route-block", "afternoon_peak")
}

func TestBacktestZeroDenominatorIsInsufficientData(t *testing.T) {
	observed := ObservedDataset{SchemaVersion: ObservedSchemaVersion, AgencyTimezone: "America/Los_Angeles"}
	predictions := PredictionDataset{SchemaVersion: PredictionSchemaVersion}
	report, err := RunBacktest(observed, predictions, []byte(`{"records":[]}`), []byte(`{"records":[]}`), BacktestOptions{
		GeneratedAt: time.Date(2026, 5, 9, 20, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("run zero denominator: %v", err)
	}
	if report.Summary.Overall.MaturityGate != "insufficient_data" {
		t.Fatalf("gate = %q, want insufficient_data", report.Summary.Overall.MaturityGate)
	}
	if report.Summary.Overall.PredictionCoverage.Status != "not_applicable" {
		t.Fatalf("coverage = %+v, want not_applicable", report.Summary.Overall.PredictionCoverage)
	}
	if report.Summary.Overall.MeanLeadTimeSeconds != nil || report.Summary.Overall.P90LeadTimeSeconds != nil {
		t.Fatalf("lead time metrics = %+v, want omitted for zero denominator", report.Summary.Overall)
	}
}

func TestBacktestSyntheticScaleUsesIndexedJoin(t *testing.T) {
	const count = 10000
	observed, predictions := syntheticScaleDatasets(count)
	report, err := RunBacktest(observed, predictions, []byte(`{"synthetic":"observed"}`), []byte(`{"synthetic":"predictions"}`), BacktestOptions{
		GeneratedAt:        time.Date(2026, 5, 9, 20, 0, 0, 0, time.UTC),
		StalePredictionTTL: 30 * time.Minute,
		MaxInputRecords:    count + 10,
		MaxGroups:          200,
	})
	if err != nil {
		t.Fatalf("run synthetic scale backtest: %v", err)
	}
	overall := report.Summary.Overall
	if overall.CoverageDenominator != count || overall.CoverageNumerator != count || overall.MatchedPredictionCount != count {
		t.Fatalf("overall scale metrics = %+v, want all %d records matched", overall, count)
	}
	if overall.MAEAbsoluteErrorSeconds == nil || *overall.MAEAbsoluteErrorSeconds != 20 {
		t.Fatalf("scale mae = %v, want 20.0", overall.MAEAbsoluteErrorSeconds)
	}
	if overall.MeanLeadTimeSeconds == nil || *overall.MeanLeadTimeSeconds != 120 {
		t.Fatalf("scale mean lead time = %v, want 120.0", overall.MeanLeadTimeSeconds)
	}
	if len(report.Metrics.Groups) > 200 {
		t.Fatalf("group count = %d, want bounded under 200", len(report.Metrics.Groups))
	}
}

func TestBacktestSchemaValidationRejectsMalformedInputs(t *testing.T) {
	if err := ValidateObservedDataset(ObservedDataset{SchemaVersion: "bad"}, 0); err == nil || !strings.Contains(err.Error(), "schema_version") {
		t.Fatalf("observed validation error = %v, want schema_version", err)
	}
	confidence := 1.2
	err := ValidatePredictionDataset(PredictionDataset{
		SchemaVersion: PredictionSchemaVersion,
		Records: []PredictionSample{{
			GeneratedTime: time.Date(2026, 5, 9, 20, 0, 0, 0, time.UTC),
			AdapterName:   "deterministic",
			AgencyID:      "demo",
			FeedVersionID: "feed",
			RouteID:       "route",
			TripID:        "trip",
			StartDate:     "20260509",
			StartTime:     "08:00:00",
			StopSequence:  1,
			EventType:     "arrival",
			PredictedTime: time.Date(2026, 5, 9, 20, 5, 0, 0, time.UTC),
			Confidence:    &confidence,
		}},
	}, 0)
	if err == nil || !strings.Contains(err.Error(), "confidence") {
		t.Fatalf("prediction validation error = %v, want confidence", err)
	}
}

func TestBacktestOutputPathSafety(t *testing.T) {
	root := t.TempDir()
	when := time.Date(2026, 5, 9, 20, 0, 0, 0, time.UTC)
	if _, err := ResolveBacktestOutputTarget(filepath.Join("docs", "evidence", "phase-50"), when, root); err == nil {
		t.Fatal("docs/evidence output path was accepted")
	}
	if _, err := ResolveBacktestOutputTarget(filepath.Join("reports", "evidence-like"), when, root); err == nil {
		t.Fatal("evidence-like output path was accepted")
	}
	if _, err := ResolveBacktestOutputTarget(filepath.Join("..", "outside"), when, root); err == nil {
		t.Fatal("unsafe traversal output path was accepted")
	}
	if _, err := ResolveBacktestOutputTarget(filepath.Join("tmp", "backtest"), when, root); err == nil {
		t.Fatal("non-.cache output path was accepted")
	}
	target := filepath.Join(root, ".cache", "target")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	link := filepath.Join(root, ".cache", "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("symlink: %v", err)
	}
	if _, err := ResolveBacktestOutputTarget(filepath.Join(".cache", "link", "out"), when, root); err == nil {
		t.Fatal("symlink ancestor output path was accepted")
	}
	output, err := ResolveBacktestOutputTarget("", when, root)
	if err != nil {
		t.Fatalf("default output: %v", err)
	}
	if output.Ref != ".cache/realtime-quality-backtest/20260509T200000Z" {
		t.Fatalf("default ref = %q", output.Ref)
	}
}

func TestBacktestFilesAreExactAndRedacted(t *testing.T) {
	observed, observedRaw, err := LoadObservedDataset(filepath.Join("..", "..", "testdata", "realtime-quality-backtest", "observed-events.json"))
	if err != nil {
		t.Fatalf("load observed: %v", err)
	}
	predictions, predictionsRaw, err := LoadPredictionDataset(filepath.Join("..", "..", "testdata", "realtime-quality-backtest", "prediction-samples.json"))
	if err != nil {
		t.Fatalf("load predictions: %v", err)
	}
	when := time.Date(2026, 5, 9, 20, 0, 0, 0, time.UTC)
	report, err := RunBacktest(observed, predictions, observedRaw, predictionsRaw, BacktestOptions{GeneratedAt: when})
	if err != nil {
		t.Fatalf("run backtest: %v", err)
	}
	root := t.TempDir()
	output, err := ResolveBacktestOutputTarget(filepath.Join(".cache", "phase-50-test"), when, root)
	if err != nil {
		t.Fatalf("resolve output: %v", err)
	}
	privateObservedPath := filepath.Join(root, "private", "observed.json")
	privatePredictionPath := filepath.Join(root, "private", "predictions.json")
	files, err := BuildBacktestFiles(report, output, ForbiddenBacktestOutputValues(privateObservedPath, privatePredictionPath), 0)
	if err != nil {
		t.Fatalf("build files: %v", err)
	}
	if err := WriteBacktestFiles(files); err != nil {
		t.Fatalf("write files: %v", err)
	}
	entries, err := os.ReadDir(output.Dir)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	expected := ExpectedBacktestOutputFiles()
	sortStrings(names)
	sortStrings(expected)
	if strings.Join(names, ",") != strings.Join(expected, ",") {
		t.Fatalf("files = %v, want %v", names, expected)
	}
	combined := ""
	for _, name := range names {
		raw, err := os.ReadFile(filepath.Join(output.Dir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		combined += string(raw)
	}
	for _, forbidden := range []string{privateObservedPath, privatePredictionPath, "/Users/", `"ready"`, `"production_ready"`, `"compliant"`, `"accepted"`} {
		if strings.Contains(strings.ToLower(combined), strings.ToLower(forbidden)) {
			t.Fatalf("output leaked forbidden value %q", forbidden)
		}
	}
	var manifest ManifestDocument
	if err := json.Unmarshal(files.Files["manifest.json"], &manifest); err != nil {
		t.Fatalf("manifest JSON invalid: %v", err)
	}
	if !manifest.AggregateOnly || manifest.RawRowsPersisted {
		t.Fatalf("manifest aggregate/raw flags = %+v", manifest)
	}
}

func assertGroupPresent(t *testing.T, groups []MetricGroup, kind string, route string, period string) {
	t.Helper()
	for _, group := range groups {
		if group.GroupType == kind && group.RouteID == route && group.TimePeriod == period {
			return
		}
	}
	t.Fatalf("group %s/%s/%s not found in %+v", kind, route, period, groups)
}

func sortStrings(values []string) {
	sort.Strings(values)
}

func syntheticScaleDatasets(count int) (ObservedDataset, PredictionDataset) {
	baseObserved := time.Date(2026, 5, 9, 15, 0, 0, 0, time.UTC)
	observed := ObservedDataset{
		SchemaVersion:  ObservedSchemaVersion,
		AgencyTimezone: "America/Los_Angeles",
		Records:        make([]ObservedEvent, 0, count),
	}
	predictions := PredictionDataset{
		SchemaVersion: PredictionSchemaVersion,
		Records:       make([]PredictionSample, 0, count),
	}
	confidence := 0.9
	for i := 0; i < count; i++ {
		routeID := fmt.Sprintf("route-%02d", i%20)
		tripID := fmt.Sprintf("trip-%05d", i)
		stopSequence := i%50 + 1
		observedTime := baseObserved.Add(time.Duration(i%720) * time.Minute)
		generatedTime := observedTime.Add(-100 * time.Second)
		predictedTime := observedTime.Add(20 * time.Second)
		observed.Records = append(observed.Records, ObservedEvent{
			AgencyID:      "scale-agency",
			FeedVersionID: "scale-feed",
			RouteID:       routeID,
			TripID:        tripID,
			StartDate:     "20260509",
			StartTime:     fmt.Sprintf("%02d:%02d:00", 5+(i%16), i%60),
			StopID:        fmt.Sprintf("stop-%02d", stopSequence),
			StopSequence:  stopSequence,
			EventType:     "arrival",
			ObservedTime:  observedTime,
		})
		predictions.Records = append(predictions.Records, PredictionSample{
			GeneratedTime:        generatedTime,
			AdapterName:          "deterministic",
			AgencyID:             "scale-agency",
			FeedVersionID:        "scale-feed",
			RouteID:              routeID,
			TripID:               tripID,
			StartDate:            "20260509",
			StartTime:            fmt.Sprintf("%02d:%02d:00", 5+(i%16), i%60),
			StopSequence:         stopSequence,
			EventType:            "arrival",
			PredictedTime:        predictedTime,
			Confidence:           &confidence,
			ScheduleRelationship: "scheduled",
		})
	}
	for i, j := 0, len(predictions.Records)-1; i < j; i, j = i+1, j-1 {
		predictions.Records[i], predictions.Records[j] = predictions.Records[j], predictions.Records[i]
	}
	return observed, predictions
}
