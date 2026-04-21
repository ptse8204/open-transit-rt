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
		&fakeScheduleRepo{active: gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"}},
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
	active gtfs.FeedVersion
	err    error
}

func (f *fakeScheduleRepo) Agency(context.Context, string) (gtfs.Agency, error) {
	return gtfs.Agency{}, errors.New("not implemented")
}

func (f *fakeScheduleRepo) ActiveFeedVersion(context.Context, string) (gtfs.FeedVersion, error) {
	if f.err != nil {
		return gtfs.FeedVersion{}, f.err
	}
	return f.active, nil
}

func (f *fakeScheduleRepo) ListTripCandidates(context.Context, string, string, string) ([]gtfs.TripCandidate, error) {
	return nil, errors.New("not implemented")
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
	result  prediction.Result
	err     error
	calls   int
	request prediction.Request
}

func (f *fakeAdapter) Name() string {
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

func (f *fakeDiagnosticsRepo) SaveTripUpdatesDiagnostics(_ context.Context, record prediction.DiagnosticsRecord) (prediction.DiagnosticsPersistenceResult, error) {
	f.records = append(f.records, record)
	if f.err != nil {
		return prediction.DiagnosticsPersistenceResult{}, f.err
	}
	return prediction.DiagnosticsPersistenceResult{Stored: true}, nil
}
