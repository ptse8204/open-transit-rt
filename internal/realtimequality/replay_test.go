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

func TestReplayPhase29AfterMidnightAndFrequencyScenarios(t *testing.T) {
	tests := []struct {
		fixture              string
		vehicleID            string
		tripID               string
		reason               string
		scheduleRelationship string
	}{
		{
			fixture:              "after-midnight-service.json",
			vehicleID:            "bus-night-1",
			tripID:               "trip-night-2500",
			scheduleRelationship: "scheduled",
		},
		{
			fixture:              "frequency-exact-window.json",
			vehicleID:            "bus-frequency-exact",
			tripID:               "trip-frequency-exact",
			reason:               "frequency_exact_instance",
			scheduleRelationship: "scheduled",
		},
		{
			fixture:              "frequency-non-exact-window.json",
			vehicleID:            "bus-frequency-1",
			tripID:               "trip-frequency-0800",
			reason:               "frequency_non_exact_conservative",
			scheduleRelationship: "unscheduled",
		},
	}
	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			report := runReplayFixture(t, tt.fixture)
			assignment := assignmentForVehicle(t, report, tt.vehicleID)
			if assignment.TripID != tt.tripID || assignment.State != "in_service" {
				t.Fatalf("assignment = %+v, want in-service trip %s", assignment, tt.tripID)
			}
			if tt.reason != "" && !hasReplayReason(assignment.ReasonCodes, tt.reason) {
				t.Fatalf("reason_codes = %+v, want %s", assignment.ReasonCodes, tt.reason)
			}
			if len(report.TripUpdates) != 1 || report.TripUpdates[0].ScheduleRelation != tt.scheduleRelationship {
				t.Fatalf("trip_updates = %+v, want one %s update", report.TripUpdates, tt.scheduleRelationship)
			}
		})
	}
}

func TestReplayPhase29ManualOverrideExpiryTransition(t *testing.T) {
	before := runReplayFixture(t, "manual-override-before-expiry.json")
	beforeAssignment := assignmentForVehicle(t, before, "bus-override-expiry")
	if beforeAssignment.Source != "manual_override" || !hasReplayReason(beforeAssignment.ReasonCodes, "manual_override_active") {
		t.Fatalf("before expiry assignment = %+v, want authoritative manual override", beforeAssignment)
	}
	if before.Metrics.ManualOverrideAssignments != 1 || before.Metrics.WithheldByReason["missing_current_stop_sequence"] != 1 {
		t.Fatalf("before expiry metrics = %+v, want manual override and withheld reason visible", before.Metrics)
	}

	after := runReplayFixture(t, "manual-override-after-expiry.json")
	afterAssignment := assignmentForVehicle(t, after, "bus-override-expired")
	if afterAssignment.Source != "automatic" || hasReplayReason(afterAssignment.ReasonCodes, "manual_override_active") {
		t.Fatalf("after expiry assignment = %+v, want return to automatic matching", afterAssignment)
	}
	if after.Metrics.ManualOverrideAssignments != 0 || after.Metrics.EligiblePredictionCandidates != 1 || len(after.TripUpdates) != 1 {
		t.Fatalf("after expiry metrics/trip updates = %+v / %+v, want automatic prediction path restored", after.Metrics, after.TripUpdates)
	}
}

func TestReplayPhase29CancellationAndZeroDenominatorDiagnostics(t *testing.T) {
	report := runReplayFixture(t, "cancellation-alert-linkage.json")
	metrics := report.Metrics
	if metrics.CanceledTripsEmitted != 1 || metrics.CancellationAlertLinksExpected != 1 || metrics.CancellationAlertLinksMissing != 1 {
		t.Fatalf("metrics = %+v, want cancellation alert linkage counts visible", metrics)
	}
	if metrics.TripUpdatesCoverageRate.Status != "not_applicable" || metrics.TripUpdatesCoverageRate.Denominator != 0 {
		t.Fatalf("coverage = %+v, want zero-denominator not_applicable", metrics.TripUpdatesCoverageRate)
	}
}

func TestReplayPhase29HardPatternsKeepUncertaintyVisible(t *testing.T) {
	report := runReplayFixture(t, "stale-ambiguous-hard-pattern.json")
	metrics := report.Metrics
	if metrics.UnknownAssignments != 2 || metrics.AmbiguousAssignments != 1 || metrics.StaleTelemetryRows != 1 {
		t.Fatalf("metrics = %+v, want unknown, ambiguous, and stale counts visible", metrics)
	}
	if metrics.WithheldByReason["stale_telemetry"] != 1 || metrics.WithheldByReason["degraded_assignment"] != 1 {
		t.Fatalf("withheld_by_reason = %+v, want stale and degraded visibility", metrics.WithheldByReason)
	}
}

func runReplayFixture(t *testing.T, fixture string) Report {
	t.Helper()
	scenario, err := LoadScenario(filepath.Join("..", "..", "testdata", "replay", fixture))
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
	return report
}

func assignmentForVehicle(t *testing.T, report Report, vehicleID string) AssignmentReport {
	t.Helper()
	for _, assignment := range report.Assignments {
		if assignment.VehicleID == vehicleID {
			return assignment
		}
	}
	t.Fatalf("assignment for vehicle %s not found in %+v", vehicleID, report.Assignments)
	return AssignmentReport{}
}

func hasReplayReason(reasons []string, want string) bool {
	for _, reason := range reasons {
		if reason == want {
			return true
		}
	}
	return false
}
