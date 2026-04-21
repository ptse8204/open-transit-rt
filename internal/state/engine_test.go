package state

import (
	"context"
	"errors"
	"testing"
	"time"

	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/telemetry"
)

func TestEngineAfterMidnightUsesPreviousServiceDate(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()
	engine := NewEngine(fakeScheduleWithTrips("overnight-agency", "America/Vancouver", "feed-night", []gtfs.TripCandidate{
		nightTrip(),
	}), assignments, DefaultConfig())

	event := storedEvent("overnight-agency", "night-bus-1", "trip-night-2430", time.Date(2026, 4, 21, 7, 30, 0, 0, time.UTC), 49.2827, -123.1207)
	assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
	if err != nil {
		t.Fatalf("match event: %v", err)
	}
	if assignment.State != StateInService {
		t.Fatalf("state = %s, want in_service", assignment.State)
	}
	if assignment.ServiceDate != "20260420" || assignment.StartDate != "20260420" {
		t.Fatalf("service/start date = %s/%s, want 20260420", assignment.ServiceDate, assignment.StartDate)
	}
	if assignment.StartTime != "24:30:00" {
		t.Fatalf("start_time = %s, want 24:30:00", assignment.StartTime)
	}
}

func TestEngineExactFrequencyInstancesAreNotCollapsed(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()
	engine := NewEngine(fakeScheduleWithTrips("freq-agency", "America/Vancouver", "feed-freq", []gtfs.TripCandidate{
		frequencyTrip("trip-loop-exact", 1),
	}), assignments, DefaultConfig())

	event := storedEvent("freq-agency", "freq-bus-1", "trip-loop-exact", time.Date(2026, 4, 20, 16, 10, 0, 0, time.UTC), 49.2827, -123.1207)
	assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
	if err != nil {
		t.Fatalf("match event: %v", err)
	}
	if assignment.TripID != "trip-loop-exact" {
		t.Fatalf("assignment = %+v, want trip-loop-exact", assignment)
	}
	if assignment.StartTime != "09:10:00" {
		t.Fatalf("start_time = %s, want generated instance 09:10:00", assignment.StartTime)
	}
	if !hasReason(assignment, ReasonFrequencyExactInstance) {
		t.Fatalf("reason_codes = %+v, want %s", assignment.ReasonCodes, ReasonFrequencyExactInstance)
	}
}

func TestEngineNonExactFrequencyUsesConservativeIdentity(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()
	engine := NewEngine(fakeScheduleWithTrips("freq-agency", "America/Vancouver", "feed-freq", []gtfs.TripCandidate{
		frequencyTrip("trip-loop", 0),
	}), assignments, DefaultConfig())

	event := storedEvent("freq-agency", "freq-bus-1", "trip-loop", time.Date(2026, 4, 20, 15, 10, 0, 0, time.UTC), 49.276, -123.115)
	assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
	if err != nil {
		t.Fatalf("match event: %v", err)
	}
	if assignment.StartTime != "08:00:00" {
		t.Fatalf("start_time = %s, want conservative frequency window identity 08:00:00", assignment.StartTime)
	}
	if !hasReason(assignment, ReasonFrequencyNonExact) {
		t.Fatalf("reason_codes = %+v, want %s", assignment.ReasonCodes, ReasonFrequencyNonExact)
	}
	if assignment.ScoreDetails["frequency_identity_type"] != "non_exact_window" || assignment.ScoreDetails["exact_scheduled_instance"] != false {
		t.Fatalf("score_details = %+v, want conservative non-exact window identity", assignment.ScoreDetails)
	}
}

func TestEngineUnknownPersistsExplicitRowAndIncident(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()
	assignments.current = &Assignment{
		AgencyID:         "demo-agency",
		VehicleID:        "bus-10",
		State:            StateInService,
		TripID:           "trip-10-0800",
		StartDate:        "20260420",
		StartTime:        "08:00:00",
		AssignmentSource: AssignmentSourceAutomatic,
		DegradedState:    DegradedNone,
		ActiveFrom:       time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC),
	}
	engine := NewEngine(fakeScheduleWithTrips("demo-agency", "America/Vancouver", "feed-demo", []gtfs.TripCandidate{
		dayTrip("trip-10-0800", "block-10", "shape-10"),
	}), assignments, DefaultConfig())

	event := storedEvent("demo-agency", "bus-10", "trip-10-0800", time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC), 49.2827, -123.1207)
	assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(2*time.Minute))
	if err != nil {
		t.Fatalf("match stale event: %v", err)
	}
	if assignment.State != StateUnknown {
		t.Fatalf("state = %s, want unknown", assignment.State)
	}
	if assignment.ServiceDate != "20260420" {
		t.Fatalf("service_date = %s, want resolvable 20260420", assignment.ServiceDate)
	}
	if assignments.closedCount != 1 {
		t.Fatalf("closed current rows = %d, want 1", assignments.closedCount)
	}
	if len(assignments.incidents) != 1 || assignments.incidents[0].Type != IncidentStaleTelemetry {
		t.Fatalf("incidents = %+v, want stale incident", assignments.incidents)
	}
}

func TestEngineMissingShapeReducesConfidenceButDoesNotBlockStrongEvidence(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()
	assignments.current = &Assignment{
		AgencyID:         "demo-agency",
		VehicleID:        "bus-10",
		State:            StateInService,
		TripID:           "trip-10-0800",
		StartDate:        "20260420",
		StartTime:        "08:00:00",
		AssignmentSource: AssignmentSourceAutomatic,
		BlockID:          "block-10",
		DegradedState:    DegradedNone,
		ActiveFrom:       time.Date(2026, 4, 20, 14, 59, 0, 0, time.UTC),
	}
	trip := dayTrip("trip-10-0800", "block-10", "")
	trip.ShapeID = ""
	trip.ShapePoints = nil
	engine := NewEngine(fakeScheduleWithTrips("demo-agency", "America/Vancouver", "feed-demo", []gtfs.TripCandidate{trip}), assignments, DefaultConfig())

	event := storedEvent("demo-agency", "bus-10", "trip-10-0800", time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC), 49.2827, -123.1207)
	assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
	if err != nil {
		t.Fatalf("match event: %v", err)
	}
	if assignment.State != StateInService {
		t.Fatalf("state = %s, want in_service despite missing shape", assignment.State)
	}
	if !hasReason(assignment, ReasonMissingShape) {
		t.Fatalf("reason_codes = %+v, want missing_shape", assignment.ReasonCodes)
	}
	if assignment.DegradedState != DegradedMissingShape {
		t.Fatalf("degraded_state = %s, want missing_shape", assignment.DegradedState)
	}
}

func TestEngineManualOverridePrecedence(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()
	assignments.override = &ManualOverride{
		ID:        7,
		AgencyID:  "demo-agency",
		VehicleID: "bus-10",
		Type:      "trip_assignment",
		RouteID:   "route-10",
		TripID:    "trip-override",
		StartDate: "20260420",
		StartTime: "08:30:00",
		State:     StateInService,
		Reason:    "dispatcher correction",
	}
	engine := NewEngine(fakeScheduleWithTrips("demo-agency", "America/Vancouver", "feed-demo", []gtfs.TripCandidate{
		dayTrip("trip-10-0800", "block-10", "shape-10"),
	}), assignments, DefaultConfig())

	event := storedEvent("demo-agency", "bus-10", "trip-10-0800", time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC), 49.2827, -123.1207)
	assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
	if err != nil {
		t.Fatalf("match event: %v", err)
	}
	if assignment.AssignmentSource != AssignmentSourceManualOverride || assignment.TripID != "trip-override" {
		t.Fatalf("assignment = %+v, want manual override trip", assignment)
	}
}

func TestEngineBlockTransitionReason(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()
	assignments.current = &Assignment{
		AgencyID:         "demo-agency",
		VehicleID:        "bus-10",
		State:            StateInService,
		TripID:           "trip-10-0800",
		BlockID:          "block-10",
		StartDate:        "20260420",
		StartTime:        "08:00:00",
		AssignmentSource: AssignmentSourceAutomatic,
		DegradedState:    DegradedNone,
		ActiveFrom:       time.Date(2026, 4, 20, 15, 20, 0, 0, time.UTC),
	}
	next := dayTrip("trip-10-0830", "block-10", "shape-10")
	next.StopTimes[0].ArrivalSeconds = 8*3600 + 30*60
	next.StopTimes[0].DepartureSeconds = 8*3600 + 30*60
	next.StopTimes[1].ArrivalSeconds = 8*3600 + 40*60
	next.StopTimes[1].DepartureSeconds = 8*3600 + 40*60
	engine := NewEngine(fakeScheduleWithTrips("demo-agency", "America/Vancouver", "feed-demo", []gtfs.TripCandidate{next}), assignments, DefaultConfig())

	event := storedEvent("demo-agency", "bus-10", "trip-10-0830", time.Date(2026, 4, 20, 15, 30, 0, 0, time.UTC), 49.2827, -123.1207)
	assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
	if err != nil {
		t.Fatalf("match event: %v", err)
	}
	if assignment.TripID != "trip-10-0830" || !hasReason(assignment, ReasonBlockTransitionMatch) {
		t.Fatalf("assignment = %+v, want block transition match", assignment)
	}
}

func TestEngineAmbiguousCandidatesPersistUnknownIncident(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()
	first := dayTrip("trip-a", "block-a", "shape-10")
	second := dayTrip("trip-b", "block-b", "shape-10")
	engine := NewEngine(fakeScheduleWithTrips("demo-agency", "America/Vancouver", "feed-demo", []gtfs.TripCandidate{first, second}), assignments, DefaultConfig())

	event := storedEvent("demo-agency", "bus-ambiguous", "", time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC), 49.2827, -123.1207)
	assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
	if err != nil {
		t.Fatalf("match ambiguous event: %v", err)
	}
	if assignment.State != StateUnknown || assignment.DegradedState != DegradedAmbiguous {
		t.Fatalf("assignment = %+v, want ambiguous unknown", assignment)
	}
	if len(assignments.incidents) != 1 || assignments.incidents[0].Type != IncidentAssignmentAmbiguous {
		t.Fatalf("incidents = %+v, want ambiguous incident", assignments.incidents)
	}
}

func TestEngineNoScheduleCandidatesPersistsUnknown(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()
	engine := NewEngine(fakeScheduleWithTrips("demo-agency", "America/Vancouver", "feed-demo", nil), assignments, DefaultConfig())

	event := storedEvent("demo-agency", "bus-no-schedule", "", time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC), 49.4000, -123.3000)
	assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
	if err != nil {
		t.Fatalf("match no schedule event: %v", err)
	}
	if assignment.State != StateUnknown || !hasReason(assignment, ReasonNoScheduleCandidates) {
		t.Fatalf("assignment = %+v, want no-schedule unknown", assignment)
	}
	if assignment.ServiceDate != "20260420" {
		t.Fatalf("service_date = %s, want resolvable 20260420", assignment.ServiceDate)
	}
}

func TestEngineContinuityRequiresTemporalPlausibility(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()
	assignments.current = &Assignment{
		AgencyID:         "demo-agency",
		VehicleID:        "bus-10",
		State:            StateInService,
		TripID:           "trip-10-0800",
		StartDate:        "20260420",
		StartTime:        "08:00:00",
		AssignmentSource: AssignmentSourceAutomatic,
		DegradedState:    DegradedNone,
		ActiveFrom:       time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC),
	}
	engine := NewEngine(fakeScheduleWithTrips("demo-agency", "America/Vancouver", "feed-demo", []gtfs.TripCandidate{
		dayTrip("trip-10-0800", "block-10", "shape-10"),
	}), assignments, DefaultConfig())

	event := storedEvent("demo-agency", "bus-10", "trip-10-0800", time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC), 49.2827, -123.1207)
	assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
	if err != nil {
		t.Fatalf("match event: %v", err)
	}
	if hasReason(assignment, ReasonContinuityMatch) {
		t.Fatalf("reason_codes = %+v, did not expect continuity outside window", assignment.ReasonCodes)
	}
}

func TestEngineBlockTransitionRequiresTemporalPlausibility(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()
	assignments.current = &Assignment{
		AgencyID:         "demo-agency",
		VehicleID:        "bus-10",
		State:            StateInService,
		TripID:           "trip-10-0800",
		BlockID:          "block-10",
		StartDate:        "20260420",
		StartTime:        "08:00:00",
		AssignmentSource: AssignmentSourceAutomatic,
		DegradedState:    DegradedNone,
		ActiveFrom:       time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC),
	}
	next := dayTrip("trip-10-0830", "block-10", "shape-10")
	next.StopTimes[0].ArrivalSeconds = 8*3600 + 30*60
	next.StopTimes[0].DepartureSeconds = 8*3600 + 30*60
	next.StopTimes[1].ArrivalSeconds = 8*3600 + 40*60
	next.StopTimes[1].DepartureSeconds = 8*3600 + 40*60
	engine := NewEngine(fakeScheduleWithTrips("demo-agency", "America/Vancouver", "feed-demo", []gtfs.TripCandidate{next}), assignments, DefaultConfig())

	event := storedEvent("demo-agency", "bus-10", "trip-10-0830", time.Date(2026, 4, 20, 15, 30, 0, 0, time.UTC), 49.2827, -123.1207)
	assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
	if err != nil {
		t.Fatalf("match event: %v", err)
	}
	if hasReason(assignment, ReasonBlockTransitionMatch) {
		t.Fatalf("reason_codes = %+v, did not expect block transition outside window", assignment.ReasonCodes)
	}
}

func TestEngineConfigMergesPartialCustomValues(t *testing.T) {
	engine := NewEngine(fakeScheduleWithTrips("demo-agency", "America/Vancouver", "feed-demo", nil), newFakeAssignments(), Config{MinConfidence: 0.8})
	if engine.config.MinConfidence != 0.8 {
		t.Fatalf("min confidence = %f, want custom 0.8", engine.config.MinConfidence)
	}
	if engine.config.StaleThreshold != DefaultConfig().StaleThreshold {
		t.Fatalf("stale threshold = %s, want default %s", engine.config.StaleThreshold, DefaultConfig().StaleThreshold)
	}
	if engine.config.ContinuityWindow != DefaultConfig().ContinuityWindow {
		t.Fatalf("continuity window = %s, want default %s", engine.config.ContinuityWindow, DefaultConfig().ContinuityWindow)
	}
}

func TestNewEnginePanicsOnInvalidConstruction(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("NewEngine did not panic for missing repositories")
		}
	}()
	_ = NewEngine(nil, nil, DefaultConfig())
}

func TestEngineSystemFailuresUseDistinctReasons(t *testing.T) {
	ctx := context.Background()
	assignments := newFakeAssignments()

	t.Run("agency lookup failure", func(t *testing.T) {
		engine := NewEngine(&failingSchedule{agencyErr: errors.New("database unavailable")}, assignments, DefaultConfig())
		event := storedEvent("demo-agency", "bus-10", "", time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC), 49.2827, -123.1207)
		assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
		if err != nil {
			t.Fatalf("match event: %v", err)
		}
		if !hasReason(assignment, ReasonAgencyLookupFailed) || hasReason(assignment, ReasonNoScheduleCandidates) {
			t.Fatalf("reason_codes = %+v, want agency failure without no_schedule", assignment.ReasonCodes)
		}
	})

	t.Run("active feed failure", func(t *testing.T) {
		engine := NewEngine(&failingSchedule{agency: gtfs.Agency{ID: "demo-agency", Timezone: "America/Vancouver"}, feedErr: errors.New("no active feed")}, assignments, DefaultConfig())
		event := storedEvent("demo-agency", "bus-10", "", time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC), 49.2827, -123.1207)
		assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
		if err != nil {
			t.Fatalf("match event: %v", err)
		}
		if !hasReason(assignment, ReasonActiveFeedUnavailable) || hasReason(assignment, ReasonNoScheduleCandidates) {
			t.Fatalf("reason_codes = %+v, want active feed failure without no_schedule", assignment.ReasonCodes)
		}
	})

	t.Run("schedule query failure", func(t *testing.T) {
		engine := NewEngine(&failingSchedule{
			agency:   gtfs.Agency{ID: "demo-agency", Timezone: "America/Vancouver"},
			feed:     gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"},
			queryErr: errors.New("query timeout"),
		}, assignments, DefaultConfig())
		event := storedEvent("demo-agency", "bus-10", "", time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC), 49.2827, -123.1207)
		assignment, err := engine.MatchEvent(ctx, event, event.Timestamp.Add(30*time.Second))
		if err != nil {
			t.Fatalf("match event: %v", err)
		}
		if !hasReason(assignment, ReasonScheduleQueryFailed) || hasReason(assignment, ReasonNoScheduleCandidates) {
			t.Fatalf("reason_codes = %+v, want query failure without no_schedule", assignment.ReasonCodes)
		}
	})
}

type fakeSchedule struct {
	agency gtfs.Agency
	feed   gtfs.FeedVersion
	trips  []gtfs.TripCandidate
}

type failingSchedule struct {
	agency    gtfs.Agency
	feed      gtfs.FeedVersion
	agencyErr error
	feedErr   error
	queryErr  error
}

func (f *failingSchedule) Agency(_ context.Context, _ string) (gtfs.Agency, error) {
	if f.agencyErr != nil {
		return gtfs.Agency{}, f.agencyErr
	}
	return f.agency, nil
}

func (f *failingSchedule) ActiveFeedVersion(_ context.Context, _ string) (gtfs.FeedVersion, error) {
	if f.feedErr != nil {
		return gtfs.FeedVersion{}, f.feedErr
	}
	return f.feed, nil
}

func (f *failingSchedule) ListTripCandidates(_ context.Context, _ string, _ string, _ string) ([]gtfs.TripCandidate, error) {
	if f.queryErr != nil {
		return nil, f.queryErr
	}
	return nil, nil
}

func fakeScheduleWithTrips(agencyID string, timezone string, feedVersionID string, trips []gtfs.TripCandidate) *fakeSchedule {
	return &fakeSchedule{
		agency: gtfs.Agency{ID: agencyID, Timezone: timezone},
		feed:   gtfs.FeedVersion{ID: feedVersionID, AgencyID: agencyID},
		trips:  trips,
	}
}

func (f *fakeSchedule) Agency(_ context.Context, agencyID string) (gtfs.Agency, error) {
	if agencyID != f.agency.ID {
		return gtfs.Agency{}, errors.New("missing agency")
	}
	return f.agency, nil
}

func (f *fakeSchedule) ActiveFeedVersion(_ context.Context, agencyID string) (gtfs.FeedVersion, error) {
	if agencyID != f.feed.AgencyID {
		return gtfs.FeedVersion{}, errors.New("missing active feed")
	}
	return f.feed, nil
}

func (f *fakeSchedule) ListTripCandidates(_ context.Context, agencyID string, feedVersionID string, serviceDate string) ([]gtfs.TripCandidate, error) {
	if agencyID != f.feed.AgencyID || feedVersionID != f.feed.ID {
		return nil, nil
	}
	trips := make([]gtfs.TripCandidate, len(f.trips))
	for i, trip := range f.trips {
		trip.AgencyID = agencyID
		trip.FeedVersionID = feedVersionID
		trip.ServiceDate = serviceDate
		trips[i] = trip
	}
	return trips, nil
}

type fakeAssignments struct {
	current     *Assignment
	override    *ManualOverride
	saved       []Assignment
	incidents   []Incident
	closedCount int
}

func newFakeAssignments() *fakeAssignments {
	return &fakeAssignments{}
}

func (f *fakeAssignments) ActiveManualOverride(_ context.Context, _ string, _ string, _ time.Time) (*ManualOverride, error) {
	return f.override, nil
}

func (f *fakeAssignments) CurrentAssignment(_ context.Context, _ string, _ string) (*Assignment, error) {
	return f.current, nil
}

func (f *fakeAssignments) SaveAssignment(_ context.Context, assignment Assignment, incidents []Incident) (Assignment, error) {
	if f.current != nil {
		f.closedCount++
	}
	assignment.ID = int64(len(f.saved) + 1)
	f.current = &assignment
	f.saved = append(f.saved, assignment)
	f.incidents = append(f.incidents, incidents...)
	return assignment, nil
}

func storedEvent(agencyID string, vehicleID string, tripHint string, observedAt time.Time, lat float64, lon float64) telemetry.StoredEvent {
	return telemetry.StoredEvent{
		ID: 101,
		Event: telemetry.Event{
			AgencyID:  agencyID,
			DeviceID:  "device-" + vehicleID,
			VehicleID: vehicleID,
			Timestamp: observedAt,
			Lat:       lat,
			Lon:       lon,
			Bearing:   120,
			SpeedMPS:  8,
			AccuracyM: 8,
			TripHint:  tripHint,
		},
		ReceivedAt:   observedAt.Add(5 * time.Second),
		IngestStatus: telemetry.IngestStatusAccepted,
	}
}

func dayTrip(tripID string, blockID string, shapeID string) gtfs.TripCandidate {
	trip := gtfs.TripCandidate{
		RouteID:   "route-10",
		ServiceID: "weekday",
		TripID:    tripID,
		BlockID:   blockID,
		ShapeID:   shapeID,
		StopTimes: []gtfs.StopTime{
			{TripID: tripID, StopID: "stop-1", ArrivalSeconds: 8 * 3600, DepartureSeconds: 8 * 3600, StopSequence: 1, ShapeDistTraveled: 0},
			{TripID: tripID, StopID: "stop-2", ArrivalSeconds: 8*3600 + 10*60, DepartureSeconds: 8*3600 + 10*60, StopSequence: 2, ShapeDistTraveled: 1200},
		},
	}
	if shapeID != "" {
		trip.ShapePoints = []gtfs.ShapePoint{
			{ShapeID: shapeID, Lat: 49.2827, Lon: -123.1207, Sequence: 1, DistTraveled: 0, HasDistance: true},
			{ShapeID: shapeID, Lat: 49.2760, Lon: -123.1150, Sequence: 2, DistTraveled: 1200, HasDistance: true},
		}
	}
	return trip
}

func nightTrip() gtfs.TripCandidate {
	trip := dayTrip("trip-night-2430", "block-night", "shape-night")
	trip.RouteID = "night-owl"
	trip.ServiceID = "night-service"
	trip.StopTimes = []gtfs.StopTime{
		{TripID: trip.TripID, StopID: "night-1", ArrivalSeconds: 24*3600 + 30*60, DepartureSeconds: 24*3600 + 30*60, StopSequence: 1, ShapeDistTraveled: 0},
		{TripID: trip.TripID, StopID: "night-2", ArrivalSeconds: 25*3600 + 15*60, DepartureSeconds: 25*3600 + 15*60, StopSequence: 2, ShapeDistTraveled: 1500},
	}
	trip.ShapePoints = []gtfs.ShapePoint{
		{ShapeID: "shape-night", Lat: 49.2827, Lon: -123.1207, Sequence: 1, DistTraveled: 0, HasDistance: true},
		{ShapeID: "shape-night", Lat: 49.2900, Lon: -123.1100, Sequence: 2, DistTraveled: 1500, HasDistance: true},
	}
	return trip
}

func frequencyTrip(tripID string, exactTimes int) gtfs.TripCandidate {
	trip := dayTrip(tripID, "block-loop", "shape-loop")
	trip.RouteID = "loop"
	trip.StopTimes = []gtfs.StopTime{
		{TripID: tripID, StopID: "freq-1", ArrivalSeconds: 9 * 3600, DepartureSeconds: 9 * 3600, StopSequence: 1, ShapeDistTraveled: 0},
		{TripID: tripID, StopID: "freq-2", ArrivalSeconds: 9*3600 + 8*60, DepartureSeconds: 9*3600 + 8*60, StopSequence: 2, ShapeDistTraveled: 900},
		{TripID: tripID, StopID: "freq-3", ArrivalSeconds: 9*3600 + 16*60, DepartureSeconds: 9*3600 + 16*60, StopSequence: 3, ShapeDistTraveled: 1800},
	}
	trip.ShapePoints = []gtfs.ShapePoint{
		{ShapeID: "shape-loop", Lat: 49.2827, Lon: -123.1207, Sequence: 1, DistTraveled: 0, HasDistance: true},
		{ShapeID: "shape-loop", Lat: 49.2760, Lon: -123.1150, Sequence: 2, DistTraveled: 900, HasDistance: true},
		{ShapeID: "shape-loop", Lat: 49.2700, Lon: -123.1000, Sequence: 3, DistTraveled: 1800, HasDistance: true},
	}
	start := "09:00:00"
	end := "10:00:00"
	if exactTimes == 0 {
		trip.StopTimes[0].ArrivalSeconds = 8 * 3600
		trip.StopTimes[0].DepartureSeconds = 8 * 3600
		trip.StopTimes[1].ArrivalSeconds = 8*3600 + 8*60
		trip.StopTimes[1].DepartureSeconds = 8*3600 + 8*60
		trip.StopTimes[2].ArrivalSeconds = 8*3600 + 16*60
		trip.StopTimes[2].DepartureSeconds = 8*3600 + 16*60
		start = "08:00:00"
		end = "09:00:00"
	}
	startSeconds, _ := gtfs.ParseGTFSTime(start)
	endSeconds, _ := gtfs.ParseGTFSTime(end)
	trip.Frequencies = []gtfs.Frequency{{
		TripID:       tripID,
		StartSeconds: startSeconds,
		EndSeconds:   endSeconds,
		HeadwaySecs:  600,
		ExactTimes:   exactTimes,
		StartTime:    start,
		EndTime:      end,
	}}
	return trip
}

func hasReason(assignment Assignment, reason string) bool {
	for _, got := range assignment.ReasonCodes {
		if got == reason {
			return true
		}
	}
	return false
}
