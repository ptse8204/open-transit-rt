package tripupdates

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/proto"

	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/prediction"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

func TestNoopTripUpdatesSnapshotIsValidEmptyFeedWithDiagnostics(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	diagnostics := &fakeDiagnosticsRepo{}
	builder := newTestBuilder(t,
		&fakeScheduleRepo{active: gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"}},
		&fakeTelemetryRepo{},
		&fakeStateRepo{},
		prediction.NewNoopAdapter(),
		diagnostics,
	)

	snapshot, err := builder.Snapshot(context.Background(), generatedAt)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if snapshot.Diagnostics.Status != prediction.StatusNoop || snapshot.Diagnostics.Reason != prediction.ReasonNoopAdapter {
		t.Fatalf("diagnostics = %+v, want noop/noop_adapter", snapshot.Diagnostics)
	}
	if snapshot.DiagnosticsPersistenceOutcome != DiagnosticsPersistenceStored || len(diagnostics.records) != 1 {
		t.Fatalf("diagnostics persistence = %q records=%d, want stored one row", snapshot.DiagnosticsPersistenceOutcome, len(diagnostics.records))
	}

	message := unmarshalFeed(t, mustMarshalProto(t, snapshot))
	if message.GetHeader().GetTimestamp() != uint64(generatedAt.Unix()) {
		t.Fatalf("header timestamp = %d, want generated_at", message.GetHeader().GetTimestamp())
	}
	if message.GetHeader().GetGtfsRealtimeVersion() != "2.0" || message.GetHeader().GetIncrementality() != gtfsrt.FeedHeader_FULL_DATASET {
		t.Fatalf("header = %+v, want GTFS-RT 2.0 FULL_DATASET", message.GetHeader())
	}
	if len(message.Entity) != 0 {
		t.Fatalf("entities = %d, want empty no-op feed", len(message.Entity))
	}

	debugJSON, err := snapshot.MarshalDebugJSON()
	if err != nil {
		t.Fatalf("debug json: %v", err)
	}
	var debug Debug
	if err := json.Unmarshal(debugJSON, &debug); err != nil {
		t.Fatalf("decode debug: %v", err)
	}
	if debug.VehiclePositionsURL != testVehiclePositionsURL || debug.DiagnosticsPersistenceOutcome != DiagnosticsPersistenceStored {
		t.Fatalf("debug = %+v, want URL and persistence outcome", debug)
	}
}

func TestTripUpdatesBuilderPassesInputsAndSortsOutput(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	events := []telemetry.StoredEvent{
		tripUpdateStoredEvent(1, "bus-z", generatedAt.Add(-20*time.Second)),
		tripUpdateStoredEvent(2, "bus-a", generatedAt.Add(-10*time.Second)),
	}
	assignments := map[string]state.Assignment{
		"bus-z": {AgencyID: "demo-agency", VehicleID: "bus-z", State: state.StateInService, TripID: "trip-z"},
	}
	adapter := &fakeAdapter{
		result: prediction.Result{
			Diagnostics: prediction.Diagnostics{Status: prediction.StatusOK, Reason: "test_adapter"},
			TripUpdates: []prediction.TripUpdate{
				{
					EntityID:  "z-entity",
					VehicleID: "bus-z",
					TripID:    "trip-z",
					RouteID:   "route-10",
					StartDate: "20260421",
					StartTime: "08:00:00",
					StopTimeUpdates: []prediction.StopTimeUpdate{
						{StopID: "stop-3", StopSequence: 3},
						{StopID: "stop-1", StopSequence: 1},
					},
				},
				{
					EntityID:  "a-entity",
					VehicleID: "bus-a",
					TripID:    "trip-a",
					RouteID:   "route-10",
					StartDate: "20260421",
					StartTime: "08:05:00",
					StopTimeUpdates: []prediction.StopTimeUpdate{
						{StopID: "stop-2", StopSequence: 2},
					},
				},
			},
		},
	}
	builder := newTestBuilder(t,
		&fakeScheduleRepo{
			active: gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"},
			tripsByDate: map[string][]gtfs.TripCandidate{
				"20260421": {
					tripUpdateCandidate("trip-a", "20260421", "08:05:00", "demo-agency", "feed-demo"),
					tripUpdateCandidate("trip-z", "20260421", "08:00:00", "demo-agency", "feed-demo"),
				},
			},
		},
		&fakeTelemetryRepo{events: events},
		&fakeStateRepo{assignments: assignments},
		adapter,
		&fakeDiagnosticsRepo{},
	)

	snapshot, err := builder.Snapshot(context.Background(), generatedAt)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if adapter.calls != 1 {
		t.Fatalf("adapter calls = %d, want 1", adapter.calls)
	}
	if adapter.request.ActiveFeedVersion.ID != "feed-demo" || adapter.request.VehiclePositionsURL != testVehiclePositionsURL {
		t.Fatalf("adapter request = %+v, want active feed and vehicle positions URL", adapter.request)
	}
	if len(adapter.request.Telemetry) != 2 || len(adapter.request.Assignments) != 1 {
		t.Fatalf("adapter input counts telemetry=%d assignments=%d, want 2/1", len(adapter.request.Telemetry), len(adapter.request.Assignments))
	}

	message := unmarshalFeed(t, mustMarshalProto(t, snapshot))
	if len(message.Entity) != 2 {
		t.Fatalf("entity count = %d, want 2", len(message.Entity))
	}
	if message.Entity[0].GetId() != "a-entity" || message.Entity[1].GetId() != "z-entity" {
		t.Fatalf("entity order = %s,%s; want deterministic entity id sort", message.Entity[0].GetId(), message.Entity[1].GetId())
	}
	stopUpdates := message.Entity[1].GetTripUpdate().GetStopTimeUpdate()
	if len(stopUpdates) != 2 || stopUpdates[0].GetStopSequence() != 1 || stopUpdates[1].GetStopSequence() != 3 {
		t.Fatalf("stop update order = %+v, want ascending stop_sequence", stopUpdates)
	}
}

func TestHappyPathMockExternalAdapterUsesNormalizationAndDiagnosticsPersistence(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	arrival := generatedAt.Add(2 * time.Minute)
	confidence := 0.82
	diagnostics := &fakeDiagnosticsRepo{}
	adapter := &fakeAdapter{
		name: "mock-external",
		result: prediction.Result{
			Diagnostics: prediction.Diagnostics{
				Status: prediction.StatusOK,
				Reason: prediction.ReasonPredictionsAvailable,
				Details: map[string]any{
					"require_prediction_confidence": true,
					"adapter_contract":              "phase_29a_test_only",
				},
			},
			TripUpdates: []prediction.TripUpdate{
				{
					AgencyID:      "demo-agency",
					FeedVersionID: "feed-demo",
					EntityID:      "entity-b",
					VehicleID:     "bus-b",
					TripID:        "trip-b",
					RouteID:       "route-10",
					StartDate:     "20260421",
					StartTime:     "08:10:00",
					Confidence:    &confidence,
					StopTimeUpdates: []prediction.StopTimeUpdate{
						{StopID: "stop-2", StopSequence: 2, ArrivalTime: &arrival},
						{StopID: "stop-1", StopSequence: 1},
					},
				},
				{
					AgencyID:      "demo-agency",
					FeedVersionID: "feed-demo",
					EntityID:      "entity-a",
					VehicleID:     "bus-a",
					TripID:        "trip-a",
					RouteID:       "route-10",
					StartDate:     "20260421",
					StartTime:     "08:05:00",
					Confidence:    &confidence,
					StopTimeUpdates: []prediction.StopTimeUpdate{
						{StopID: "stop-1", StopSequence: 1},
					},
				},
			},
		},
	}
	builder := newTestBuilder(t,
		&fakeScheduleRepo{
			active: gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"},
			tripsByDate: map[string][]gtfs.TripCandidate{
				"20260421": {
					tripUpdateCandidate("trip-a", "20260421", "08:05:00", "demo-agency", "feed-demo"),
					tripUpdateCandidate("trip-b", "20260421", "08:10:00", "demo-agency", "feed-demo"),
				},
			},
		},
		&fakeTelemetryRepo{},
		&fakeStateRepo{},
		adapter,
		diagnostics,
	)

	snapshot, err := builder.Snapshot(context.Background(), generatedAt)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if snapshot.DiagnosticsPersistenceOutcome != DiagnosticsPersistenceStored || len(diagnostics.records) != 1 {
		t.Fatalf("diagnostics persistence = %q records=%d, want stored one row", snapshot.DiagnosticsPersistenceOutcome, len(diagnostics.records))
	}
	if diagnostics.records[0].AdapterName != "mock-external" || diagnostics.records[0].InputCounts.TripUpdatesOutput != 2 {
		t.Fatalf("diagnostics record = %+v, want mock adapter output persisted", diagnostics.records[0])
	}

	message := unmarshalFeed(t, mustMarshalProto(t, snapshot))
	if message.GetHeader().GetGtfsRealtimeVersion() != "2.0" || message.GetHeader().GetIncrementality() != gtfsrt.FeedHeader_FULL_DATASET {
		t.Fatalf("header = %+v, want unchanged public GTFS-RT contract", message.GetHeader())
	}
	if len(message.Entity) != 2 || message.Entity[0].GetId() != "entity-a" || message.Entity[1].GetId() != "entity-b" {
		t.Fatalf("entity order = %+v, want normalized entity ordering", message.Entity)
	}
	stopUpdates := message.Entity[1].GetTripUpdate().GetStopTimeUpdate()
	if len(stopUpdates) != 2 || stopUpdates[0].GetStopSequence() != 1 || stopUpdates[1].GetStopSequence() != 2 {
		t.Fatalf("stop updates = %+v, want normalized stop ordering", stopUpdates)
	}
}

func TestMockExternalAdapterOutputValidationRejectsUnsafePredictions(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	fresh := generatedAt.Add(90 * time.Second)
	stale := generatedAt.Add(-time.Second)
	highConfidence := 0.9
	lowConfidence := 0.2

	tests := []struct {
		name       string
		candidate  gtfs.TripCandidate
		update     prediction.TripUpdate
		wantReason string
	}{
		{
			name: "trip not in active feed",
			update: prediction.TripUpdate{
				EntityID: "missing-trip", TripID: "trip-missing", StartDate: "20260421", StartTime: "08:00:00", Confidence: &highConfidence,
				StopTimeUpdates: []prediction.StopTimeUpdate{{StopID: "stop-1", StopSequence: 1, ArrivalTime: &fresh}},
			},
			wantReason: adapterOutputTripNotInActiveFeed,
		},
		{
			name:      "impossible stop sequence",
			candidate: tripUpdateCandidate("trip-10", "20260421", "08:00:00", "demo-agency", "feed-demo"),
			update: prediction.TripUpdate{
				EntityID: "bad-sequence", TripID: "trip-10", StartDate: "20260421", StartTime: "08:00:00", Confidence: &highConfidence,
				StopTimeUpdates: []prediction.StopTimeUpdate{{StopID: "stop-99", StopSequence: 99, ArrivalTime: &fresh}},
			},
			wantReason: adapterOutputInvalidStopSequence,
		},
		{
			name:      "stale timestamp",
			candidate: tripUpdateCandidate("trip-10", "20260421", "08:00:00", "demo-agency", "feed-demo"),
			update: prediction.TripUpdate{
				EntityID: "stale-prediction", TripID: "trip-10", StartDate: "20260421", StartTime: "08:00:00", Confidence: &highConfidence,
				StopTimeUpdates: []prediction.StopTimeUpdate{{StopID: "stop-1", StopSequence: 1, ArrivalTime: &stale}},
			},
			wantReason: adapterOutputStalePrediction,
		},
		{
			name:      "wrong agency",
			candidate: tripUpdateCandidate("trip-10", "20260421", "08:00:00", "demo-agency", "feed-demo"),
			update: prediction.TripUpdate{
				AgencyID: "other-agency", EntityID: "wrong-agency", TripID: "trip-10", StartDate: "20260421", StartTime: "08:00:00", Confidence: &highConfidence,
				StopTimeUpdates: []prediction.StopTimeUpdate{{StopID: "stop-1", StopSequence: 1, ArrivalTime: &fresh}},
			},
			wantReason: adapterOutputWrongAgency,
		},
		{
			name:      "wrong feed version",
			candidate: tripUpdateCandidate("trip-10", "20260421", "08:00:00", "demo-agency", "feed-demo"),
			update: prediction.TripUpdate{
				FeedVersionID: "other-feed", EntityID: "wrong-feed", TripID: "trip-10", StartDate: "20260421", StartTime: "08:00:00", Confidence: &highConfidence,
				StopTimeUpdates: []prediction.StopTimeUpdate{{StopID: "stop-1", StopSequence: 1, ArrivalTime: &fresh}},
			},
			wantReason: adapterOutputWrongFeedVersion,
		},
		{
			name:      "unsupported added trip prediction",
			candidate: tripUpdateCandidate("trip-added", "20260421", "08:00:00", "demo-agency", "feed-demo"),
			update: prediction.TripUpdate{
				EntityID: "added-trip", TripID: "trip-added", StartDate: "20260421", StartTime: "08:00:00", Confidence: &highConfidence,
				ScheduleRelationship: prediction.ScheduleRelationshipAdded,
				StopTimeUpdates:      []prediction.StopTimeUpdate{{StopID: "stop-1", StopSequence: 1, ArrivalTime: &fresh}},
			},
			wantReason: adapterOutputUnsupportedDisruption,
		},
		{
			name:      "low confidence",
			candidate: tripUpdateCandidate("trip-10", "20260421", "08:00:00", "demo-agency", "feed-demo"),
			update: prediction.TripUpdate{
				EntityID: "low-confidence", TripID: "trip-10", StartDate: "20260421", StartTime: "08:00:00", Confidence: &lowConfidence,
				StopTimeUpdates: []prediction.StopTimeUpdate{{StopID: "stop-1", StopSequence: 1, ArrivalTime: &fresh}},
			},
			wantReason: adapterOutputLowConfidence,
		},
		{
			name:      "missing confidence",
			candidate: tripUpdateCandidate("trip-10", "20260421", "08:00:00", "demo-agency", "feed-demo"),
			update: prediction.TripUpdate{
				EntityID: "missing-confidence", TripID: "trip-10", StartDate: "20260421", StartTime: "08:00:00",
				StopTimeUpdates: []prediction.StopTimeUpdate{{StopID: "stop-1", StopSequence: 1, ArrivalTime: &fresh}},
			},
			wantReason: adapterOutputMissingConfidence,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trips := []gtfs.TripCandidate{}
			if tt.candidate.TripID != "" {
				trips = append(trips, tt.candidate)
			}
			adapter := &fakeAdapter{
				name: "mock-external",
				result: prediction.Result{
					Diagnostics: prediction.Diagnostics{
						Status:  prediction.StatusOK,
						Reason:  prediction.ReasonPredictionsAvailable,
						Details: map[string]any{"require_prediction_confidence": true},
					},
					TripUpdates: []prediction.TripUpdate{tt.update},
				},
			}
			builder := newTestBuilder(t,
				&fakeScheduleRepo{
					active:      gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"},
					tripsByDate: map[string][]gtfs.TripCandidate{"20260421": trips},
				},
				&fakeTelemetryRepo{},
				&fakeStateRepo{},
				adapter,
				&fakeDiagnosticsRepo{},
			)

			snapshot, err := builder.Snapshot(context.Background(), generatedAt)
			if err != nil {
				t.Fatalf("snapshot: %v", err)
			}
			if len(snapshot.TripUpdates) != 0 {
				t.Fatalf("trip updates = %+v, want unsafe mock external output withheld", snapshot.TripUpdates)
			}
			if snapshot.Diagnostics.Status != prediction.StatusError || snapshot.Diagnostics.Reason != prediction.ReasonAdapterOutputRejected {
				t.Fatalf("diagnostics = %+v, want rejected adapter output", snapshot.Diagnostics)
			}
			if got := snapshot.Diagnostics.Metrics.WithheldByReason[tt.wantReason]; got != 1 {
				t.Fatalf("withheld = %+v, want %s=1", snapshot.Diagnostics.Metrics.WithheldByReason, tt.wantReason)
			}
			message := unmarshalFeed(t, mustMarshalProto(t, snapshot))
			if len(message.Entity) != 0 {
				t.Fatalf("protobuf entities = %d, want valid empty feed after rejection", len(message.Entity))
			}
		})
	}
}

func TestTripUpdatesAdapterFailuresReturnVisibleDiagnostics(t *testing.T) {
	tests := []error{
		errors.New("predictor timeout"),
		errors.New("predictor unavailable"),
		errors.New("malformed predictor response"),
	}
	for _, adapterErr := range tests {
		t.Run(adapterErr.Error(), func(t *testing.T) {
			diagnostics := &fakeDiagnosticsRepo{}
			builder := newTestBuilder(t,
				&fakeScheduleRepo{active: gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"}},
				&fakeTelemetryRepo{},
				&fakeStateRepo{},
				&fakeAdapter{name: "mock-external", err: adapterErr},
				diagnostics,
			)

			snapshot, err := builder.Snapshot(context.Background(), time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC))
			if err != nil {
				t.Fatalf("snapshot: %v", err)
			}
			if len(snapshot.TripUpdates) != 0 {
				t.Fatalf("trip updates = %+v, want withheld output on adapter failure", snapshot.TripUpdates)
			}
			if snapshot.Diagnostics.Status != prediction.StatusError || snapshot.Diagnostics.Reason != prediction.ReasonAdapterError {
				t.Fatalf("diagnostics = %+v, want visible adapter error", snapshot.Diagnostics)
			}
			if len(diagnostics.records) != 1 || diagnostics.records[0].Status != prediction.StatusError {
				t.Fatalf("diagnostics records = %+v, want persisted adapter failure", diagnostics.records)
			}
		})
	}
}

func TestTripUpdatesMissingActiveFeedDoesNotCallAdapter(t *testing.T) {
	adapter := &fakeAdapter{}
	diagnostics := &fakeDiagnosticsRepo{}
	builder := newTestBuilder(t,
		&fakeScheduleRepo{err: pgx.ErrNoRows},
		&fakeTelemetryRepo{events: []telemetry.StoredEvent{tripUpdateStoredEvent(1, "bus-10", time.Now())}},
		&fakeStateRepo{},
		adapter,
		diagnostics,
	)

	snapshot, err := builder.Snapshot(context.Background(), time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if adapter.calls != 0 {
		t.Fatalf("adapter calls = %d, want 0 without active feed", adapter.calls)
	}
	if snapshot.Diagnostics.Reason != prediction.ReasonActiveFeedUnavailable {
		t.Fatalf("diagnostics = %+v, want active feed unavailable", snapshot.Diagnostics)
	}
	if len(diagnostics.records) != 1 || diagnostics.records[0].ActiveFeedVersionID != "" {
		t.Fatalf("diagnostics records = %+v, want persisted missing-active-feed trace", diagnostics.records)
	}
	message := unmarshalFeed(t, mustMarshalProto(t, snapshot))
	if len(message.Entity) != 0 {
		t.Fatalf("entity count = %d, want empty degraded feed", len(message.Entity))
	}
}

func TestTripUpdatesAdapterErrorReturnsEmptyDiagnostics(t *testing.T) {
	builder := newTestBuilder(t,
		&fakeScheduleRepo{active: gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"}},
		&fakeTelemetryRepo{},
		&fakeStateRepo{},
		&fakeAdapter{err: errors.New("predictor down")},
		&fakeDiagnosticsRepo{},
	)
	snapshot, err := builder.Snapshot(context.Background(), time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if snapshot.Diagnostics.Status != prediction.StatusError || snapshot.Diagnostics.Reason != prediction.ReasonAdapterError {
		t.Fatalf("diagnostics = %+v, want adapter error", snapshot.Diagnostics)
	}
	if len(snapshot.TripUpdates) != 0 {
		t.Fatalf("trip updates = %d, want empty output on adapter error", len(snapshot.TripUpdates))
	}
}

func TestTripUpdatesBuilderWithDeterministicAdapterProducesNonEmptyFeed(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 5, 0, 0, time.UTC)
	adapter, err := prediction.NewDeterministicAdapter(
		&fakeScheduleRepo{
			active: gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"},
			agency: gtfs.Agency{ID: "demo-agency", Timezone: "America/Vancouver"},
			tripsByDate: map[string][]gtfs.TripCandidate{
				"20260420": {{
					AgencyID:      "demo-agency",
					FeedVersionID: "feed-demo",
					ServiceDate:   "20260420",
					RouteID:       "route-10",
					ServiceID:     "weekday",
					TripID:        "trip-10",
					StopTimes: []gtfs.StopTime{
						{TripID: "trip-10", StopID: "stop-1", StopSequence: 1, ArrivalSeconds: 8 * 3600, DepartureSeconds: 8 * 3600},
						{TripID: "trip-10", StopID: "stop-2", StopSequence: 2, ArrivalSeconds: 8*3600 + 600, DepartureSeconds: 8*3600 + 600},
					},
				}},
			},
		},
		&fakePredictionOperations{},
		prediction.DeterministicConfig{},
	)
	if err != nil {
		t.Fatalf("new deterministic adapter: %v", err)
	}
	builder := newTestBuilder(t,
		&fakeScheduleRepo{
			active: gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"},
			agency: gtfs.Agency{ID: "demo-agency", Timezone: "America/Vancouver"},
			tripsByDate: map[string][]gtfs.TripCandidate{
				"20260420": {{
					AgencyID:      "demo-agency",
					FeedVersionID: "feed-demo",
					ServiceDate:   "20260420",
					RouteID:       "route-10",
					ServiceID:     "weekday",
					TripID:        "trip-10",
					StopTimes: []gtfs.StopTime{
						{TripID: "trip-10", StopID: "stop-1", StopSequence: 1, ArrivalSeconds: 8 * 3600, DepartureSeconds: 8 * 3600},
						{TripID: "trip-10", StopID: "stop-2", StopSequence: 2, ArrivalSeconds: 8*3600 + 600, DepartureSeconds: 8*3600 + 600},
					},
				}},
			},
		},
		&fakeTelemetryRepo{events: []telemetry.StoredEvent{tripUpdateStoredEvent(10, "bus-10", generatedAt)}},
		&fakeStateRepo{assignments: map[string]state.Assignment{
			"bus-10": {
				AgencyID:            "demo-agency",
				VehicleID:           "bus-10",
				FeedVersionID:       "feed-demo",
				TelemetryEventID:    10,
				State:               state.StateInService,
				ServiceDate:         "20260420",
				RouteID:             "route-10",
				TripID:              "trip-10",
				StartDate:           "20260420",
				StartTime:           "08:00:00",
				CurrentStopSequence: 1,
				Confidence:          0.9,
				AssignmentSource:    state.AssignmentSourceAutomatic,
				DegradedState:       state.DegradedNone,
			},
		}},
		adapter,
		&fakeDiagnosticsRepo{},
	)
	snapshot, err := builder.Snapshot(context.Background(), generatedAt)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if len(snapshot.TripUpdates) != 1 {
		t.Fatalf("trip updates = %+v, want deterministic non-empty output", snapshot.TripUpdates)
	}
	if snapshot.Diagnostics.Metrics.EligiblePredictionCandidates != 1 {
		t.Fatalf("metrics = %+v, want first-class prediction metrics", snapshot.Diagnostics.Metrics)
	}
	message := unmarshalFeed(t, mustMarshalProto(t, snapshot))
	if len(message.Entity) != 1 || len(message.Entity[0].GetTripUpdate().GetStopTimeUpdate()) != 1 {
		t.Fatalf("protobuf message = %+v, want one entity with future stop update", message)
	}
}

func newTestBuilder(
	t *testing.T,
	scheduleRepo gtfs.Repository,
	telemetryRepo telemetry.Repository,
	assignmentRepo state.Repository,
	adapter prediction.Adapter,
	diagnosticsRepo prediction.DiagnosticsRepository,
) *Builder {
	t.Helper()
	builder, err := NewBuilder(scheduleRepo, telemetryRepo, assignmentRepo, adapter, diagnosticsRepo, Config{
		AgencyID:            "demo-agency",
		MaxVehicles:         10,
		VehiclePositionsURL: testVehiclePositionsURL,
	})
	if err != nil {
		t.Fatalf("new builder: %v", err)
	}
	return builder
}

const testVehiclePositionsURL = "http://localhost:8083/public/gtfsrt/vehicle_positions.pb"

func tripUpdateStoredEvent(id int64, vehicleID string, observedAt time.Time) telemetry.StoredEvent {
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
		ReceivedAt:   observedAt.Add(5 * time.Second),
		IngestStatus: telemetry.IngestStatusAccepted,
		PayloadJSON:  json.RawMessage(`{}`),
	}
}

func tripUpdateCandidate(tripID string, serviceDate string, _ string, agencyID string, feedVersionID string) gtfs.TripCandidate {
	return gtfs.TripCandidate{
		AgencyID:      agencyID,
		FeedVersionID: feedVersionID,
		ServiceDate:   serviceDate,
		RouteID:       "route-10",
		ServiceID:     "weekday",
		TripID:        tripID,
		StopTimes: []gtfs.StopTime{
			{TripID: tripID, StopID: "stop-1", StopSequence: 1},
			{TripID: tripID, StopID: "stop-2", StopSequence: 2},
			{TripID: tripID, StopID: "stop-3", StopSequence: 3},
		},
	}
}

func mustMarshalProto(t *testing.T, snapshot Snapshot) []byte {
	t.Helper()
	payload, err := snapshot.MarshalProto()
	if err != nil {
		t.Fatalf("marshal proto: %v", err)
	}
	return payload
}

func unmarshalFeed(t *testing.T, payload []byte) *gtfsrt.FeedMessage {
	t.Helper()
	var message gtfsrt.FeedMessage
	if err := proto.Unmarshal(payload, &message); err != nil {
		t.Fatalf("unmarshal feed: %v", err)
	}
	if err := proto.CheckInitialized(&message); err != nil {
		t.Fatalf("protobuf not initialized: %v", err)
	}
	return &message
}

type fakeScheduleRepo struct {
	active      gtfs.FeedVersion
	agency      gtfs.Agency
	tripsByDate map[string][]gtfs.TripCandidate
	err         error
}

func (f *fakeScheduleRepo) Agency(context.Context, string) (gtfs.Agency, error) {
	if f.agency.ID != "" {
		return f.agency, nil
	}
	return gtfs.Agency{}, errors.New("not implemented")
}

func (f *fakeScheduleRepo) ActiveFeedVersion(context.Context, string) (gtfs.FeedVersion, error) {
	if f.err != nil {
		return gtfs.FeedVersion{}, f.err
	}
	return f.active, nil
}

func (f *fakeScheduleRepo) ListTripCandidates(_ context.Context, _ string, _ string, serviceDate string) ([]gtfs.TripCandidate, error) {
	return append([]gtfs.TripCandidate(nil), f.tripsByDate[serviceDate]...), nil
}

type fakeTelemetryRepo struct {
	events []telemetry.StoredEvent
	err    error
}

func (f *fakeTelemetryRepo) Store(context.Context, telemetry.Event, json.RawMessage) (telemetry.StoreResult, error) {
	return telemetry.StoreResult{}, errors.New("not implemented")
}

func (f *fakeTelemetryRepo) LatestByVehicle(context.Context, string, string) (telemetry.StoredEvent, error) {
	return telemetry.StoredEvent{}, errors.New("not implemented")
}

func (f *fakeTelemetryRepo) ListLatestByAgency(context.Context, string, int) ([]telemetry.StoredEvent, error) {
	if f.err != nil {
		return nil, f.err
	}
	return append([]telemetry.StoredEvent(nil), f.events...), nil
}

func (f *fakeTelemetryRepo) ListEvents(context.Context, string, int) ([]telemetry.StoredEvent, error) {
	return nil, errors.New("not implemented")
}

type fakeStateRepo struct {
	assignments map[string]state.Assignment
	err         error
}

func (f *fakeStateRepo) ActiveManualOverride(context.Context, string, string, time.Time) (*state.ManualOverride, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeStateRepo) CurrentAssignment(context.Context, string, string) (*state.Assignment, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeStateRepo) ListCurrentAssignments(_ context.Context, _ string, vehicleIDs []string) (map[string]state.Assignment, error) {
	if f.err != nil {
		return nil, f.err
	}
	result := make(map[string]state.Assignment, len(vehicleIDs))
	for _, vehicleID := range vehicleIDs {
		if assignment, ok := f.assignments[vehicleID]; ok {
			result[vehicleID] = assignment
		}
	}
	return result, nil
}

func (f *fakeStateRepo) SaveAssignment(context.Context, state.Assignment, []state.Incident) (state.Assignment, error) {
	return state.Assignment{}, errors.New("SaveAssignment must not be called by trip updates builder")
}

type fakeAdapter struct {
	name    string
	result  prediction.Result
	err     error
	calls   int
	request prediction.Request
}

func (f *fakeAdapter) Name() string {
	if f.name != "" {
		return f.name
	}
	return "fake"
}

func (f *fakeAdapter) PredictTripUpdates(_ context.Context, request prediction.Request) (prediction.Result, error) {
	f.calls++
	f.request = request
	if f.err != nil {
		return prediction.Result{}, f.err
	}
	return f.result, nil
}

type fakeDiagnosticsRepo struct {
	records []prediction.DiagnosticsRecord
	err     error
}

type fakePredictionOperations struct{}

func (fakePredictionOperations) ListActivePredictionOverrides(context.Context, string, time.Time) ([]prediction.OverrideRecord, error) {
	return nil, nil
}

func (fakePredictionOperations) CreatePredictionOverride(context.Context, prediction.OverrideInput) (prediction.OverrideRecord, error) {
	return prediction.OverrideRecord{}, nil
}

func (fakePredictionOperations) ReplacePredictionOverride(context.Context, prediction.OverrideInput) (prediction.OverrideRecord, error) {
	return prediction.OverrideRecord{}, nil
}

func (fakePredictionOperations) ClearPredictionOverride(context.Context, string, int64, string, string, time.Time) error {
	return nil
}

func (fakePredictionOperations) SavePredictionReviewItems(context.Context, []prediction.ReviewItem) error {
	return nil
}

func (fakePredictionOperations) ListPredictionReviewItems(context.Context, prediction.ReviewFilter) ([]prediction.ReviewItem, error) {
	return nil, nil
}

func (fakePredictionOperations) UpdatePredictionReviewStatus(context.Context, string, int64, prediction.ReviewStatus, string, string, time.Time) error {
	return nil
}

func (f *fakeDiagnosticsRepo) SaveTripUpdatesDiagnostics(_ context.Context, record prediction.DiagnosticsRecord) (prediction.DiagnosticsPersistenceResult, error) {
	f.records = append(f.records, record)
	if f.err != nil {
		return prediction.DiagnosticsPersistenceResult{}, f.err
	}
	return prediction.DiagnosticsPersistenceResult{Stored: true}, nil
}
