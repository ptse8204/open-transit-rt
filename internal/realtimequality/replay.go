package realtimequality

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"time"

	"open-transit-rt/internal/feed"
	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/prediction"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

const defaultVehiclePositionsURL = "http://localhost:8080/public/gtfsrt/vehicle_positions.pb"

type Scenario struct {
	Name                string                      `json:"name"`
	InputFixture        string                      `json:"input_fixture"`
	AgencyID            string                      `json:"agency_id"`
	Timezone            string                      `json:"timezone"`
	FeedVersionID       string                      `json:"feed_version_id"`
	GeneratedAt         time.Time                   `json:"generated_at"`
	Trips               []TripFixture               `json:"trips"`
	Telemetry           []TelemetryFixture          `json:"telemetry_events"`
	InputAssignments    []AssignmentFixture         `json:"input_assignments,omitempty"`
	ManualOverrides     []ManualOverrideFixture     `json:"manual_overrides,omitempty"`
	PredictionOverrides []PredictionOverrideFixture `json:"prediction_overrides,omitempty"`
	Expected            Expectations                `json:"expected"`
}

type TripFixture struct {
	ServiceDate string        `json:"service_date"`
	RouteID     string        `json:"route_id"`
	TripID      string        `json:"trip_id"`
	BlockID     string        `json:"block_id,omitempty"`
	StartTime   string        `json:"start_time"`
	StopTimes   []StopFixture `json:"stop_times"`
	ShapePoints []ShapePoint  `json:"shape_points,omitempty"`
}

type StopFixture struct {
	StopID            string  `json:"stop_id"`
	StopSequence      int     `json:"stop_sequence"`
	ArrivalTime       string  `json:"arrival_time"`
	DepartureTime     string  `json:"departure_time"`
	ShapeDistTraveled float64 `json:"shape_dist_traveled,omitempty"`
}

type ShapePoint struct {
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Sequence     int     `json:"sequence"`
	DistTraveled float64 `json:"dist_traveled"`
}

type TelemetryFixture struct {
	ID        int64     `json:"id"`
	DeviceID  string    `json:"device_id"`
	VehicleID string    `json:"vehicle_id"`
	Timestamp time.Time `json:"timestamp"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	Bearing   float64   `json:"bearing,omitempty"`
	TripHint  string    `json:"trip_hint,omitempty"`
}

type AssignmentFixture struct {
	VehicleID           string   `json:"vehicle_id"`
	State               string   `json:"state"`
	RouteID             string   `json:"route_id,omitempty"`
	TripID              string   `json:"trip_id,omitempty"`
	StartDate           string   `json:"start_date,omitempty"`
	StartTime           string   `json:"start_time,omitempty"`
	CurrentStopSequence int      `json:"current_stop_sequence,omitempty"`
	Confidence          float64  `json:"confidence"`
	Source              string   `json:"source,omitempty"`
	DegradedState       string   `json:"degraded_state,omitempty"`
	ReasonCodes         []string `json:"reason_codes,omitempty"`
	TelemetryEventID    int64    `json:"telemetry_event_id,omitempty"`
}

type ManualOverrideFixture struct {
	ID        int64  `json:"id"`
	VehicleID string `json:"vehicle_id"`
	Type      string `json:"type"`
	RouteID   string `json:"route_id,omitempty"`
	TripID    string `json:"trip_id,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	State     string `json:"state"`
	Reason    string `json:"reason,omitempty"`
}

type PredictionOverrideFixture struct {
	ID        int64  `json:"id"`
	VehicleID string `json:"vehicle_id"`
	Type      string `json:"type"`
	RouteID   string `json:"route_id,omitempty"`
	TripID    string `json:"trip_id,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	State     string `json:"state"`
	Reason    string `json:"reason,omitempty"`
}

type Expectations struct {
	Assignments      []AssignmentExpectation      `json:"assignments"`
	VehiclePositions []VehiclePositionExpectation `json:"vehicle_positions"`
	TripUpdates      []TripUpdateExpectation      `json:"trip_updates"`
	WithheldReasons  map[string]int               `json:"withheld_reasons,omitempty"`
	Metrics          MetricsExpectation           `json:"metrics"`
}

type AssignmentExpectation struct {
	VehicleID     string   `json:"vehicle_id"`
	State         string   `json:"state"`
	Source        string   `json:"source,omitempty"`
	DegradedState string   `json:"degraded_state,omitempty"`
	TripID        string   `json:"trip_id,omitempty"`
	ReasonCodes   []string `json:"reason_codes,omitempty"`
}

type VehiclePositionExpectation struct {
	VehicleID          string `json:"vehicle_id"`
	Included           bool   `json:"included"`
	TripDescriptor     bool   `json:"trip_descriptor"`
	TripOmissionReason string `json:"trip_omission_reason"`
}

type TripUpdateExpectation struct {
	EntityID          string `json:"entity_id"`
	TripID            string `json:"trip_id"`
	ScheduleRelation  string `json:"schedule_relationship"`
	FutureStopUpdates int    `json:"future_stop_updates"`
}

type MetricsExpectation struct {
	TelemetryRowsConsidered      int             `json:"telemetry_rows_considered"`
	AssignmentsConsidered        int             `json:"assignments_considered"`
	UnknownAssignments           int             `json:"unknown_assignments"`
	AmbiguousAssignments         int             `json:"ambiguous_assignments"`
	DegradedAssignments          int             `json:"degraded_assignments"`
	StaleTelemetryRows           int             `json:"stale_telemetry_rows"`
	ManualOverrideAssignments    int             `json:"manual_override_assignments"`
	EligiblePredictionCandidates int             `json:"eligible_prediction_candidates"`
	TripUpdatesEmitted           int             `json:"trip_updates_emitted"`
	WithheldByReason             map[string]int  `json:"withheld_by_reason,omitempty"`
	UnknownAssignmentRate        RateExpectation `json:"unknown_assignment_rate"`
	AmbiguousAssignmentRate      RateExpectation `json:"ambiguous_assignment_rate"`
	StaleTelemetryRate           RateExpectation `json:"stale_telemetry_rate"`
	TripUpdatesCoverageRate      RateExpectation `json:"trip_updates_coverage_rate"`
	FutureStopCoverageRate       RateExpectation `json:"future_stop_coverage_rate"`
}

type RateExpectation struct {
	Numerator   int      `json:"numerator"`
	Denominator int      `json:"denominator"`
	Status      string   `json:"status"`
	Percent     *float64 `json:"percent,omitempty"`
}

type Report struct {
	ScenarioName      string                  `json:"scenario_name"`
	GeneratedAt       time.Time               `json:"generated_at"`
	Assignments       []AssignmentReport      `json:"assignments"`
	VehiclePositions  []VehiclePositionReport `json:"vehicle_positions"`
	TripUpdates       []TripUpdateReport      `json:"trip_updates"`
	Metrics           prediction.Metrics      `json:"metrics"`
	DiagnosticsStatus string                  `json:"diagnostics_status"`
	DiagnosticsReason string                  `json:"diagnostics_reason"`
}

type AssignmentReport struct {
	VehicleID     string   `json:"vehicle_id"`
	State         string   `json:"state"`
	Source        string   `json:"source"`
	DegradedState string   `json:"degraded_state"`
	TripID        string   `json:"trip_id,omitempty"`
	ReasonCodes   []string `json:"reason_codes,omitempty"`
}

type VehiclePositionReport struct {
	VehicleID          string `json:"vehicle_id"`
	Included           bool   `json:"included"`
	TripDescriptor     bool   `json:"trip_descriptor"`
	TripOmissionReason string `json:"trip_omission_reason"`
}

type TripUpdateReport struct {
	EntityID          string `json:"entity_id"`
	TripID            string `json:"trip_id"`
	ScheduleRelation  string `json:"schedule_relationship"`
	FutureStopUpdates int    `json:"future_stop_updates"`
}

func LoadScenario(path string) (Scenario, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return Scenario{}, fmt.Errorf("read replay scenario: %w", err)
	}
	var scenario Scenario
	if err := json.Unmarshal(payload, &scenario); err != nil {
		return Scenario{}, fmt.Errorf("decode replay scenario %s: %w", path, err)
	}
	if err := scenario.Validate(); err != nil {
		return Scenario{}, fmt.Errorf("validate replay scenario %s: %w", path, err)
	}
	return scenario, nil
}

func (s Scenario) Validate() error {
	switch {
	case s.Name == "":
		return fmt.Errorf("name is required")
	case s.InputFixture == "":
		return fmt.Errorf("input_fixture is required")
	case s.AgencyID == "":
		return fmt.Errorf("agency_id is required")
	case s.Timezone == "":
		return fmt.Errorf("timezone is required")
	case s.FeedVersionID == "":
		return fmt.Errorf("feed_version_id is required")
	case s.GeneratedAt.IsZero():
		return fmt.Errorf("generated_at is required")
	case s.Expected.Metrics.UnknownAssignmentRate.Status == "":
		return fmt.Errorf("expected.metrics.unknown_assignment_rate.status is required")
	case s.Expected.Metrics.TripUpdatesCoverageRate.Status == "":
		return fmt.Errorf("expected.metrics.trip_updates_coverage_rate.status is required")
	}
	return nil
}

func Run(ctx context.Context, scenario Scenario) (Report, error) {
	schedules, err := newReplaySchedule(scenario)
	if err != nil {
		return Report{}, err
	}
	stateRepo := newReplayStateRepository(scenario)
	events := replayEvents(scenario)
	if len(scenario.InputAssignments) > 0 {
		for _, fixture := range scenario.InputAssignments {
			stateRepo.assignments[fixture.VehicleID] = assignmentFromFixture(scenario, fixture)
		}
	} else {
		engine, err := state.NewEngine(schedules, stateRepo, state.Config{})
		if err != nil {
			return Report{}, err
		}
		for _, event := range events {
			if _, err := engine.MatchEvent(ctx, event, scenario.GeneratedAt); err != nil {
				return Report{}, fmt.Errorf("match event %s: %w", event.VehicleID, err)
			}
		}
	}
	assignments := stateRepo.currentAssignments()

	vp, err := feed.NewVehiclePositionsBuilder(replayTelemetryRepository{events: events}, stateRepo, feed.VehiclePositionsConfig{
		AgencyID:                  scenario.AgencyID,
		MaxVehicles:               2000,
		StaleTelemetryTTL:         state.DefaultConfig().StaleThreshold,
		SuppressStaleVehicleAfter: 5 * time.Minute,
		TripConfidenceThreshold:   state.DefaultConfig().MinConfidence,
	})
	if err != nil {
		return Report{}, err
	}
	vpSnapshot, err := vp.Snapshot(ctx, scenario.GeneratedAt)
	if err != nil {
		return Report{}, fmt.Errorf("build vehicle positions snapshot: %w", err)
	}

	ops := newReplayPredictionOperations(scenario)
	adapter, err := prediction.NewDeterministicAdapter(schedules, ops, prediction.DeterministicConfig{})
	if err != nil {
		return Report{}, err
	}
	result, err := adapter.PredictTripUpdates(ctx, prediction.Request{
		AgencyID:            scenario.AgencyID,
		GeneratedAt:         scenario.GeneratedAt,
		ActiveFeedVersion:   gtfs.FeedVersion{ID: scenario.FeedVersionID, AgencyID: scenario.AgencyID},
		Telemetry:           events,
		Assignments:         assignments,
		VehiclePositionsURL: defaultVehiclePositionsURL,
	})
	if err != nil {
		return Report{}, fmt.Errorf("predict trip updates: %w", err)
	}

	return Report{
		ScenarioName:      scenario.Name,
		GeneratedAt:       scenario.GeneratedAt,
		Assignments:       assignmentReports(assignments),
		VehiclePositions:  vehiclePositionReports(vpSnapshot),
		TripUpdates:       tripUpdateReports(result.TripUpdates),
		Metrics:           result.Diagnostics.Metrics,
		DiagnosticsStatus: result.Diagnostics.Status,
		DiagnosticsReason: result.Diagnostics.Reason,
	}, nil
}

func Compare(report Report, expected Expectations) []string {
	var mismatches []string
	if !reflect.DeepEqual(report.Assignments, expectedAssignmentReports(expected.Assignments)) {
		mismatches = append(mismatches, fmt.Sprintf("assignments got %+v want %+v", report.Assignments, expected.Assignments))
	}
	if !reflect.DeepEqual(report.VehiclePositions, expectedVehiclePositionReports(expected.VehiclePositions)) {
		mismatches = append(mismatches, fmt.Sprintf("vehicle_positions got %+v want %+v", report.VehiclePositions, expected.VehiclePositions))
	}
	if !reflect.DeepEqual(report.TripUpdates, expectedTripUpdateReports(expected.TripUpdates)) {
		mismatches = append(mismatches, fmt.Sprintf("trip_updates got %+v want %+v", report.TripUpdates, expected.TripUpdates))
	}
	metrics := report.Metrics
	want := expected.Metrics
	if metrics.TelemetryRowsConsidered != want.TelemetryRowsConsidered ||
		metrics.AssignmentsConsidered != want.AssignmentsConsidered ||
		metrics.UnknownAssignments != want.UnknownAssignments ||
		metrics.AmbiguousAssignments != want.AmbiguousAssignments ||
		metrics.DegradedAssignments != want.DegradedAssignments ||
		metrics.StaleTelemetryRows != want.StaleTelemetryRows ||
		metrics.ManualOverrideAssignments != want.ManualOverrideAssignments ||
		metrics.EligiblePredictionCandidates != want.EligiblePredictionCandidates ||
		metrics.TripUpdatesEmitted != want.TripUpdatesEmitted {
		mismatches = append(mismatches, fmt.Sprintf("metrics counts got %+v want %+v", metrics, want))
	}
	if !reflect.DeepEqual(metrics.WithheldByReason, want.WithheldByReason) {
		mismatches = append(mismatches, fmt.Sprintf("withheld_by_reason got %+v want %+v", metrics.WithheldByReason, want.WithheldByReason))
	}
	checkRate := func(name string, got prediction.RateMetric, want RateExpectation) {
		if got.Numerator != want.Numerator || got.Denominator != want.Denominator || got.Status != want.Status {
			mismatches = append(mismatches, fmt.Sprintf("%s got %+v want %+v", name, got, want))
			return
		}
		if want.Percent == nil {
			if got.Percent != nil {
				mismatches = append(mismatches, fmt.Sprintf("%s percent got %v want nil", name, *got.Percent))
			}
			return
		}
		if got.Percent == nil || *got.Percent != *want.Percent {
			mismatches = append(mismatches, fmt.Sprintf("%s percent got %v want %v", name, got.Percent, *want.Percent))
		}
	}
	checkRate("unknown_assignment_rate", metrics.UnknownAssignmentRate, want.UnknownAssignmentRate)
	checkRate("ambiguous_assignment_rate", metrics.AmbiguousAssignmentRate, want.AmbiguousAssignmentRate)
	checkRate("stale_telemetry_rate", metrics.StaleTelemetryRate, want.StaleTelemetryRate)
	checkRate("trip_updates_coverage_rate", metrics.TripUpdatesCoverageRate, want.TripUpdatesCoverageRate)
	checkRate("future_stop_coverage_rate", metrics.FutureStopCoverageRate, want.FutureStopCoverageRate)
	return mismatches
}

func replayEvents(scenario Scenario) []telemetry.StoredEvent {
	events := make([]telemetry.StoredEvent, 0, len(scenario.Telemetry))
	for _, fixture := range scenario.Telemetry {
		events = append(events, telemetry.StoredEvent{
			ID: fixture.ID,
			Event: telemetry.Event{
				AgencyID:  scenario.AgencyID,
				DeviceID:  fixture.DeviceID,
				VehicleID: fixture.VehicleID,
				Timestamp: fixture.Timestamp,
				Lat:       fixture.Lat,
				Lon:       fixture.Lon,
				Bearing:   fixture.Bearing,
				TripHint:  fixture.TripHint,
			},
			ReceivedAt:   fixture.Timestamp,
			IngestStatus: telemetry.IngestStatusAccepted,
		})
	}
	return events
}

func assignmentFromFixture(scenario Scenario, fixture AssignmentFixture) state.Assignment {
	source := state.AssignmentSource(fixture.Source)
	if source == "" {
		source = state.AssignmentSourceAutomatic
	}
	degraded := state.DegradedState(fixture.DegradedState)
	if degraded == "" {
		degraded = state.DegradedNone
	}
	return state.Assignment{
		AgencyID:            scenario.AgencyID,
		VehicleID:           fixture.VehicleID,
		FeedVersionID:       scenario.FeedVersionID,
		TelemetryEventID:    fixture.TelemetryEventID,
		State:               state.VehicleServiceState(fixture.State),
		ServiceDate:         fixture.StartDate,
		RouteID:             fixture.RouteID,
		TripID:              fixture.TripID,
		StartDate:           fixture.StartDate,
		StartTime:           fixture.StartTime,
		CurrentStopSequence: fixture.CurrentStopSequence,
		Confidence:          fixture.Confidence,
		AssignmentSource:    source,
		ReasonCodes:         append([]string(nil), fixture.ReasonCodes...),
		DegradedState:       degraded,
		ActiveFrom:          scenario.GeneratedAt,
	}
}

func assignmentReports(assignments map[string]state.Assignment) []AssignmentReport {
	reports := make([]AssignmentReport, 0, len(assignments))
	for _, assignment := range assignments {
		reasons := append([]string(nil), assignment.ReasonCodes...)
		sort.Strings(reasons)
		reports = append(reports, AssignmentReport{
			VehicleID:     assignment.VehicleID,
			State:         string(assignment.State),
			Source:        string(assignment.AssignmentSource),
			DegradedState: string(assignment.DegradedState),
			TripID:        assignment.TripID,
			ReasonCodes:   reasons,
		})
	}
	sort.SliceStable(reports, func(i int, j int) bool { return reports[i].VehicleID < reports[j].VehicleID })
	return reports
}

func vehiclePositionReports(snapshot feed.VehiclePositionsSnapshot) []VehiclePositionReport {
	reports := make([]VehiclePositionReport, 0, len(snapshot.Vehicles))
	for _, vehicle := range snapshot.Vehicles {
		reports = append(reports, VehiclePositionReport{
			VehicleID:          vehicle.VehicleID,
			Included:           vehicle.IncludedInProtobuf,
			TripDescriptor:     vehicle.TripDescriptorPublished,
			TripOmissionReason: vehicle.TripDescriptorOmissionReason,
		})
	}
	sort.SliceStable(reports, func(i int, j int) bool { return reports[i].VehicleID < reports[j].VehicleID })
	return reports
}

func tripUpdateReports(updates []prediction.TripUpdate) []TripUpdateReport {
	reports := make([]TripUpdateReport, 0, len(updates))
	for _, update := range updates {
		reports = append(reports, TripUpdateReport{
			EntityID:          update.EntityID,
			TripID:            update.TripID,
			ScheduleRelation:  string(update.ScheduleRelationship),
			FutureStopUpdates: len(update.StopTimeUpdates),
		})
	}
	sort.SliceStable(reports, func(i int, j int) bool { return reports[i].EntityID < reports[j].EntityID })
	return reports
}

func expectedAssignmentReports(expectations []AssignmentExpectation) []AssignmentReport {
	reports := make([]AssignmentReport, 0, len(expectations))
	for _, expectation := range expectations {
		reasons := append([]string(nil), expectation.ReasonCodes...)
		sort.Strings(reasons)
		reports = append(reports, AssignmentReport{
			VehicleID:     expectation.VehicleID,
			State:         expectation.State,
			Source:        expectation.Source,
			DegradedState: expectation.DegradedState,
			TripID:        expectation.TripID,
			ReasonCodes:   reasons,
		})
	}
	sort.SliceStable(reports, func(i int, j int) bool { return reports[i].VehicleID < reports[j].VehicleID })
	return reports
}

func expectedVehiclePositionReports(expectations []VehiclePositionExpectation) []VehiclePositionReport {
	reports := make([]VehiclePositionReport, 0, len(expectations))
	for _, expectation := range expectations {
		reports = append(reports, VehiclePositionReport{
			VehicleID:          expectation.VehicleID,
			Included:           expectation.Included,
			TripDescriptor:     expectation.TripDescriptor,
			TripOmissionReason: expectation.TripOmissionReason,
		})
	}
	sort.SliceStable(reports, func(i int, j int) bool { return reports[i].VehicleID < reports[j].VehicleID })
	return reports
}

func expectedTripUpdateReports(expectations []TripUpdateExpectation) []TripUpdateReport {
	reports := make([]TripUpdateReport, 0, len(expectations))
	for _, expectation := range expectations {
		reports = append(reports, TripUpdateReport{
			EntityID:          expectation.EntityID,
			TripID:            expectation.TripID,
			ScheduleRelation:  expectation.ScheduleRelation,
			FutureStopUpdates: expectation.FutureStopUpdates,
		})
	}
	sort.SliceStable(reports, func(i int, j int) bool { return reports[i].EntityID < reports[j].EntityID })
	return reports
}
