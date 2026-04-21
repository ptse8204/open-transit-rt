package prediction

import (
	"context"
	"testing"
	"time"

	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

func TestDeterministicAdapterEmitsMatchedScheduledTripUpdate(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 5, 0, 0, time.UTC)
	ops := &fakeOperationsRepo{}
	adapter := newDeterministicTestAdapter(t, ops, map[string][]gtfs.TripCandidate{
		"20260420": {predictionTrip("trip-10", "route-10", "08:00:00", "08:10:00", nil)},
	})

	result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(generatedAt, []telemetry.StoredEvent{
		predictionTelemetry(7, "bus-10", generatedAt),
	}, map[string]state.Assignment{
		"bus-10": predictionAssignment("bus-10", "trip-10", state.StateInService, state.AssignmentSourceAutomatic, 7),
	}))
	if err != nil {
		t.Fatalf("predict: %v", err)
	}
	if len(result.TripUpdates) != 1 {
		t.Fatalf("trip updates = %d, want 1", len(result.TripUpdates))
	}
	update := result.TripUpdates[0]
	if update.EntityID != "trip_update:trip-10:20260420:08:00:00" || update.ScheduleRelationship != ScheduleRelationshipScheduled {
		t.Fatalf("update = %+v, want stable scheduled entity", update)
	}
	if len(update.StopTimeUpdates) != 1 || update.StopTimeUpdates[0].StopSequence != 2 {
		t.Fatalf("stop updates = %+v, want future stop 2 only", update.StopTimeUpdates)
	}
	if got, want := update.StopTimeUpdates[0].ArrivalTime.UTC(), time.Date(2026, 4, 20, 15, 15, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("arrival = %s, want %s", got, want)
	}
	if result.Diagnostics.Metrics.EligiblePredictionCandidates != 1 || *result.Diagnostics.Metrics.CoveragePercent != 100 {
		t.Fatalf("metrics = %+v, want one eligible candidate and 100 coverage", result.Diagnostics.Metrics)
	}
	if len(ops.savedReviews) != 0 {
		t.Fatalf("reviews = %+v, want none for healthy prediction", ops.savedReviews)
	}
}

func TestDeterministicAdapterHandlesAfterMidnightServiceTimes(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 7, 40, 0, 0, time.UTC)
	adapter := newDeterministicTestAdapter(t, nil, map[string][]gtfs.TripCandidate{
		"20260420": {predictionTrip("trip-night", "night", "24:30:00", "25:15:00", nil)},
	})

	assignment := predictionAssignment("bus-night", "trip-night", state.StateInService, state.AssignmentSourceAutomatic, 8)
	assignment.StartTime = "24:30:00"
	result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(generatedAt, []telemetry.StoredEvent{
		predictionTelemetry(8, "bus-night", generatedAt),
	}, map[string]state.Assignment{
		"bus-night": assignment,
	}))
	if err != nil {
		t.Fatalf("predict: %v", err)
	}
	if len(result.TripUpdates) != 1 || len(result.TripUpdates[0].StopTimeUpdates) != 1 {
		t.Fatalf("result = %+v, want one after-midnight future stop", result)
	}
	if got, want := result.TripUpdates[0].StopTimeUpdates[0].ArrivalTime.UTC(), time.Date(2026, 4, 21, 8, 25, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("arrival = %s, want %s", got, want)
	}
}

func TestDeterministicAdapterHandlesFrequencyInstances(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 5, 0, 0, time.UTC)
	exact := predictionTrip("trip-exact", "loop", "08:00:00", "08:08:00", []gtfs.Frequency{{
		TripID: "trip-exact", StartSeconds: 8 * 3600, EndSeconds: 9 * 3600, HeadwaySecs: 600, ExactTimes: 1, StartTime: "08:00:00", EndTime: "09:00:00",
	}})
	nonExact := predictionTrip("trip-window", "loop", "08:00:00", "08:08:00", []gtfs.Frequency{{
		TripID: "trip-window", StartSeconds: 8 * 3600, EndSeconds: 9 * 3600, HeadwaySecs: 600, ExactTimes: 0, StartTime: "08:00:00", EndTime: "09:00:00",
	}})
	adapter := newDeterministicTestAdapter(t, nil, map[string][]gtfs.TripCandidate{"20260420": {exact, nonExact}})

	exactAssignment := predictionAssignment("bus-exact", "trip-exact", state.StateInService, state.AssignmentSourceAutomatic, 10)
	exactAssignment.StartTime = "08:00:00"
	windowAssignment := predictionAssignment("bus-window", "trip-window", state.StateInService, state.AssignmentSourceAutomatic, 11)
	windowAssignment.StartTime = "08:00:00"
	result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(generatedAt, []telemetry.StoredEvent{
		predictionTelemetry(10, "bus-exact", generatedAt),
		predictionTelemetry(11, "bus-window", time.Date(2026, 4, 20, 15, 5, 0, 0, time.UTC)),
	}, map[string]state.Assignment{
		"bus-exact":  exactAssignment,
		"bus-window": windowAssignment,
	}))
	if err != nil {
		t.Fatalf("predict: %v", err)
	}
	if len(result.TripUpdates) != 2 {
		t.Fatalf("trip updates = %+v, want exact and non-exact frequency updates", result.TripUpdates)
	}
	relationships := map[string]ScheduleRelationship{}
	for _, update := range result.TripUpdates {
		relationships[update.TripID] = update.ScheduleRelationship
	}
	if relationships["trip-exact"] != ScheduleRelationshipScheduled || relationships["trip-window"] != ScheduleRelationshipUnscheduled {
		t.Fatalf("relationships = %+v, want scheduled exact and unscheduled non-exact", relationships)
	}
}

func TestDeterministicAdapterSuppressesDeadheadAndLayover(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 5, 0, 0, time.UTC)
	ops := &fakeOperationsRepo{}
	adapter := newDeterministicTestAdapter(t, ops, map[string][]gtfs.TripCandidate{
		"20260420": {predictionTrip("trip-10", "route-10", "08:00:00", "08:10:00", nil)},
	})
	deadhead := predictionAssignment("bus-deadhead", "trip-10", state.StateDeadhead, state.AssignmentSourceManualOverride, 12)
	layover := predictionAssignment("bus-layover", "trip-10", state.StateLayover, state.AssignmentSourceManualOverride, 13)

	result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(generatedAt, []telemetry.StoredEvent{
		predictionTelemetry(12, "bus-deadhead", generatedAt),
		predictionTelemetry(13, "bus-layover", generatedAt),
	}, map[string]state.Assignment{
		"bus-deadhead": deadhead,
		"bus-layover":  layover,
	}))
	if err != nil {
		t.Fatalf("predict: %v", err)
	}
	if len(result.TripUpdates) != 0 {
		t.Fatalf("trip updates = %+v, want no predictions for deadhead/layover", result.TripUpdates)
	}
	if result.Diagnostics.Metrics.WithheldByReason[ReasonDeadheadNoPrediction] != 1 || result.Diagnostics.Metrics.WithheldByReason[ReasonLayoverNoPrediction] != 1 {
		t.Fatalf("withheld = %+v, want deadhead and layover reasons", result.Diagnostics.Metrics.WithheldByReason)
	}
	if len(ops.savedReviews) != 2 {
		t.Fatalf("reviews = %+v, want persisted deadhead/layover review items", ops.savedReviews)
	}
}

func TestDeterministicAdapterHandlesDisruptionOverridesConservatively(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 5, 0, 0, time.UTC)
	ops := &fakeOperationsRepo{overrides: []OverrideRecord{
		{ID: 1, AgencyID: "demo-agency", VehicleID: "bus-cancel", OverrideType: "canceled_trip", TripID: "trip-10", RouteID: "route-10", StartDate: "20260420", StartTime: "08:00:00", State: "canceled"},
		{ID: 2, AgencyID: "demo-agency", VehicleID: "bus-added", OverrideType: "added_trip", TripID: "trip-added", RouteID: "route-10", StartDate: "20260420", StartTime: "08:30:00", State: "added"},
		{ID: 3, AgencyID: "demo-agency", VehicleID: "bus-short", OverrideType: "short_turn", TripID: "trip-10", RouteID: "route-10", StartDate: "20260420", StartTime: "08:00:00", State: "short_turn"},
		{ID: 4, AgencyID: "demo-agency", VehicleID: "bus-detour", OverrideType: "detour", TripID: "trip-10", RouteID: "route-10", StartDate: "20260420", StartTime: "08:00:00", State: "detour"},
	}}
	adapter := newDeterministicTestAdapter(t, ops, map[string][]gtfs.TripCandidate{
		"20260420": {predictionTrip("trip-10", "route-10", "08:00:00", "08:10:00", nil)},
	})

	result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(generatedAt, nil, nil))
	if err != nil {
		t.Fatalf("predict: %v", err)
	}
	if len(result.TripUpdates) != 1 || result.TripUpdates[0].ScheduleRelationship != ScheduleRelationshipCanceled {
		t.Fatalf("trip updates = %+v, want one canceled TripUpdate", result.TripUpdates)
	}
	metrics := result.Diagnostics.Metrics
	if metrics.CanceledTripsEmitted != 1 || metrics.CancellationAlertLinksExpected != 1 || metrics.CancellationAlertLinksMissing != 1 {
		t.Fatalf("cancellation metrics = %+v, want emitted cancellation and missing alert linkage", metrics)
	}
	if metrics.AddedTripsWithheld != 1 || metrics.ShortTurnsWithheld != 1 || metrics.DetoursWithheld != 1 {
		t.Fatalf("disruption metrics = %+v, want conservative withholds", metrics)
	}
	if len(ops.savedReviews) != 4 {
		t.Fatalf("saved reviews = %+v, want cancellation linkage plus added/short/detour reviews", ops.savedReviews)
	}
	if got := ops.savedReviews[0].Details["expected_alert_missing"]; got != true {
		t.Fatalf("cancellation review details = %+v, want expected_alert_missing=true", ops.savedReviews[0].Details)
	}
}

func TestDeterministicAdapterWithholdsAmbiguousDuplicateTripInstance(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 5, 0, 0, time.UTC)
	ops := &fakeOperationsRepo{}
	adapter := newDeterministicTestAdapter(t, ops, map[string][]gtfs.TripCandidate{
		"20260420": {predictionTrip("trip-10", "route-10", "08:00:00", "08:10:00", nil)},
	})

	first := predictionAssignment("bus-a", "trip-10", state.StateInService, state.AssignmentSourceAutomatic, 21)
	second := predictionAssignment("bus-b", "trip-10", state.StateInService, state.AssignmentSourceAutomatic, 22)
	result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(generatedAt, []telemetry.StoredEvent{
		predictionTelemetry(21, "bus-a", generatedAt),
		predictionTelemetry(22, "bus-b", generatedAt),
	}, map[string]state.Assignment{"bus-a": first, "bus-b": second}))
	if err != nil {
		t.Fatalf("predict: %v", err)
	}
	if len(result.TripUpdates) != 0 {
		t.Fatalf("trip updates = %+v, want duplicate trip instance withheld", result.TripUpdates)
	}
	if result.Diagnostics.Metrics.EligiblePredictionCandidates != 2 || result.Diagnostics.Metrics.WithheldByReason[ReasonDuplicateTripInstance] != 2 {
		t.Fatalf("metrics = %+v, want two eligible duplicate withholds", result.Diagnostics.Metrics)
	}
}

func newDeterministicTestAdapter(t *testing.T, ops OperationsRepository, trips map[string][]gtfs.TripCandidate) *DeterministicAdapter {
	t.Helper()
	adapter, err := NewDeterministicAdapter(&fakePredictionSchedule{tripsByDate: trips}, ops, DeterministicConfig{})
	if err != nil {
		t.Fatalf("new deterministic adapter: %v", err)
	}
	return adapter
}

func predictionRequest(generatedAt time.Time, events []telemetry.StoredEvent, assignments map[string]state.Assignment) Request {
	if assignments == nil {
		assignments = map[string]state.Assignment{}
	}
	return Request{
		AgencyID:          "demo-agency",
		GeneratedAt:       generatedAt,
		ActiveFeedVersion: gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"},
		Telemetry:         events,
		Assignments:       assignments,
	}
}

func predictionTelemetry(id int64, vehicleID string, observedAt time.Time) telemetry.StoredEvent {
	return telemetry.StoredEvent{
		ID: id,
		Event: telemetry.Event{
			AgencyID:  "demo-agency",
			DeviceID:  "device-" + vehicleID,
			VehicleID: vehicleID,
			Timestamp: observedAt,
			Lat:       49.2827,
			Lon:       -123.1207,
		},
		IngestStatus: telemetry.IngestStatusAccepted,
	}
}

func predictionAssignment(vehicleID string, tripID string, serviceState state.VehicleServiceState, source state.AssignmentSource, telemetryEventID int64) state.Assignment {
	return state.Assignment{
		ID:                  telemetryEventID,
		AgencyID:            "demo-agency",
		VehicleID:           vehicleID,
		FeedVersionID:       "feed-demo",
		TelemetryEventID:    telemetryEventID,
		State:               serviceState,
		ServiceDate:         "20260420",
		RouteID:             "route-10",
		TripID:              tripID,
		StartDate:           "20260420",
		StartTime:           "08:00:00",
		CurrentStopSequence: 1,
		Confidence:          0.9,
		AssignmentSource:    source,
		DegradedState:       state.DegradedNone,
	}
}

func predictionTrip(tripID string, routeID string, firstTime string, secondTime string, frequencies []gtfs.Frequency) gtfs.TripCandidate {
	firstSeconds, _ := gtfs.ParseGTFSTime(firstTime)
	secondSeconds, _ := gtfs.ParseGTFSTime(secondTime)
	return gtfs.TripCandidate{
		AgencyID:      "demo-agency",
		FeedVersionID: "feed-demo",
		ServiceDate:   "20260420",
		RouteID:       routeID,
		ServiceID:     "weekday",
		TripID:        tripID,
		StopTimes: []gtfs.StopTime{
			{TripID: tripID, StopID: "stop-1", StopSequence: 1, ArrivalSeconds: firstSeconds, DepartureSeconds: firstSeconds},
			{TripID: tripID, StopID: "stop-2", StopSequence: 2, ArrivalSeconds: secondSeconds, DepartureSeconds: secondSeconds},
		},
		Frequencies: frequencies,
	}
}

type fakePredictionSchedule struct {
	tripsByDate map[string][]gtfs.TripCandidate
}

func (f *fakePredictionSchedule) Agency(context.Context, string) (gtfs.Agency, error) {
	return gtfs.Agency{ID: "demo-agency", Timezone: "America/Vancouver"}, nil
}

func (f *fakePredictionSchedule) ActiveFeedVersion(context.Context, string) (gtfs.FeedVersion, error) {
	return gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"}, nil
}

func (f *fakePredictionSchedule) ListTripCandidates(_ context.Context, _ string, _ string, serviceDate string) ([]gtfs.TripCandidate, error) {
	return append([]gtfs.TripCandidate(nil), f.tripsByDate[serviceDate]...), nil
}

type fakeOperationsRepo struct {
	overrides    []OverrideRecord
	savedReviews []ReviewItem
}

func (f *fakeOperationsRepo) ListActivePredictionOverrides(context.Context, string, time.Time) ([]OverrideRecord, error) {
	return append([]OverrideRecord(nil), f.overrides...), nil
}

func (f *fakeOperationsRepo) CreatePredictionOverride(context.Context, OverrideInput) (OverrideRecord, error) {
	return OverrideRecord{}, nil
}

func (f *fakeOperationsRepo) ReplacePredictionOverride(context.Context, OverrideInput) (OverrideRecord, error) {
	return OverrideRecord{}, nil
}

func (f *fakeOperationsRepo) ClearPredictionOverride(context.Context, string, int64, string, string, time.Time) error {
	return nil
}

func (f *fakeOperationsRepo) SavePredictionReviewItems(_ context.Context, items []ReviewItem) error {
	f.savedReviews = append(f.savedReviews, items...)
	return nil
}

func (f *fakeOperationsRepo) ListPredictionReviewItems(context.Context, ReviewFilter) ([]ReviewItem, error) {
	return nil, nil
}

func (f *fakeOperationsRepo) UpdatePredictionReviewStatus(context.Context, string, int64, ReviewStatus, string, string, time.Time) error {
	return nil
}
