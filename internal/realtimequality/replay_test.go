package realtimequality

import (
	"context"
	"path/filepath"
	"sort"
	"testing"
)

func TestReplayFixturesMatchCurrentBehavior(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("..", "..", "testdata", "replay", "*.json"))
	if err != nil {
		t.Fatalf("glob replay fixtures: %v", err)
	}
	if len(paths) == 0 {
		t.Fatal("no replay fixtures found")
	}
	sort.Strings(paths)
	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			scenario, err := LoadScenario(path)
			if err != nil {
				t.Fatalf("load scenario: %v", err)
			}
			report, err := Run(context.Background(), scenario)
			if err != nil {
				t.Fatalf("run scenario: %v", err)
			}
			if mismatches := Compare(report, scenario.Expected); len(mismatches) > 0 {
				t.Fatalf("replay mismatches:\n%v", mismatches)
			}
		})
	}
}

func TestReplayFixturesUseFixedClockAndStableOrdering(t *testing.T) {
	scenario, err := LoadScenario(filepath.Join("..", "..", "testdata", "replay", "matched-current-behavior.json"))
	if err != nil {
		t.Fatalf("load scenario: %v", err)
	}
	first, err := Run(context.Background(), scenario)
	if err != nil {
		t.Fatalf("first replay: %v", err)
	}
	second, err := Run(context.Background(), scenario)
	if err != nil {
		t.Fatalf("second replay: %v", err)
	}
	if !first.GeneratedAt.Equal(scenario.GeneratedAt) || !second.GeneratedAt.Equal(scenario.GeneratedAt) {
		t.Fatalf("replay did not use fixture clock: first=%s second=%s fixture=%s", first.GeneratedAt, second.GeneratedAt, scenario.GeneratedAt)
	}
	if mismatches := Compare(second, scenario.Expected); len(mismatches) > 0 {
		t.Fatalf("second replay mismatches:\n%v", mismatches)
	}
	if len(first.Assignments) != len(second.Assignments) || len(first.VehiclePositions) != len(second.VehiclePositions) || len(first.TripUpdates) != len(second.TripUpdates) {
		t.Fatalf("replay output sizes changed between runs")
	}
}

func TestReplayMetricsKeepUncertaintyVisible(t *testing.T) {
	scenario, err := LoadScenario(filepath.Join("..", "..", "testdata", "replay", "ambiguous-assignment-visible.json"))
	if err != nil {
		t.Fatalf("load scenario: %v", err)
	}
	report, err := Run(context.Background(), scenario)
	if err != nil {
		t.Fatalf("run scenario: %v", err)
	}
	metrics := report.Metrics
	if metrics.UnknownAssignments != 1 || metrics.AmbiguousAssignments != 1 || metrics.DegradedAssignments != 1 {
		t.Fatalf("metrics = %+v, want unknown/ambiguous/degraded counts preserved", metrics)
	}
	if metrics.WithheldByReason["degraded_assignment"] != 1 {
		t.Fatalf("withheld_by_reason = %+v, want degraded assignment visible", metrics.WithheldByReason)
	}
	if metrics.TripUpdatesCoverageRate.Status != "not_applicable" || metrics.TripUpdatesCoverageRate.Percent != nil {
		t.Fatalf("coverage rate = %+v, want honest zero-denominator not_applicable", metrics.TripUpdatesCoverageRate)
	}
}
