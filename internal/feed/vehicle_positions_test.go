package feed

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"

	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

func TestVehiclePositionsProtoMatchedEntity(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 0, 30, 0, time.UTC)
	event := feedStoredEvent(101, "bus-10", generatedAt.Add(-30*time.Second), `{"bearing":0,"speed_mps":8.2}`)
	assignment := matchedAssignment(event, "trip-10-0800")
	snapshot := buildTestSnapshot(t, generatedAt, []telemetry.StoredEvent{event}, map[string]state.Assignment{"bus-10": assignment}, testConfig())

	payload, err := snapshot.MarshalProto()
	if err != nil {
		t.Fatalf("marshal proto: %v", err)
	}
	message := unmarshalFeed(t, payload)
	if err := proto.CheckInitialized(message); err != nil {
		t.Fatalf("protobuf not initialized: %v", err)
	}
	if message.GetHeader().GetGtfsRealtimeVersion() != GTFSRealtimeVersion {
		t.Fatalf("gtfs version = %s, want %s", message.GetHeader().GetGtfsRealtimeVersion(), GTFSRealtimeVersion)
	}
	if message.GetHeader().GetTimestamp() != uint64(generatedAt.Unix()) {
		t.Fatalf("header timestamp = %d, want %d", message.GetHeader().GetTimestamp(), generatedAt.Unix())
	}
	if message.GetHeader().GetIncrementality() != gtfsrt.FeedHeader_FULL_DATASET {
		t.Fatalf("incrementality = %s, want FULL_DATASET", message.GetHeader().GetIncrementality())
	}
	if len(message.Entity) != 1 {
		t.Fatalf("entity count = %d, want 1", len(message.Entity))
	}
	entity := message.Entity[0]
	if entity.GetId() != "bus-10" || entity.GetVehicle().GetVehicle().GetId() != "bus-10" {
		t.Fatalf("stable ids not set: %+v", entity)
	}
	position := entity.GetVehicle().GetPosition()
	if position.GetLatitude() != float32(event.Lat) || position.GetLongitude() != float32(event.Lon) {
		t.Fatalf("position = %+v, want event coordinates", position)
	}
	if position.Bearing == nil || position.GetBearing() != 0 {
		t.Fatalf("bearing pointer/value = %v/%f, want explicit true north", position.Bearing, position.GetBearing())
	}
	if position.Speed == nil || position.GetSpeed() != float32(8.2) {
		t.Fatalf("speed pointer/value = %v/%f, want explicit speed", position.Speed, position.GetSpeed())
	}
	trip := entity.GetVehicle().GetTrip()
	if trip.GetTripId() != "trip-10-0800" || trip.GetRouteId() != "route-10" || trip.GetStartDate() != "20260420" || trip.GetStartTime() != "08:00:00" {
		t.Fatalf("trip descriptor = %+v, want matched assignment identity", trip)
	}
}

func TestVehiclePositionsSnapshotPublicationDecisions(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 0, 30, 0, time.UTC)
	tests := []struct {
		name       string
		event      telemetry.StoredEvent
		assignment *state.Assignment
		wantReason string
		wantEntity bool
		wantTrip   bool
	}{
		{
			name:       "no assignment",
			event:      feedStoredEvent(1, "bus-no-assignment", generatedAt.Add(-30*time.Second), `{}`),
			wantReason: TripDescriptorOmissionNoAssignment,
			wantEntity: true,
		},
		{
			name:       "stale overrides no assignment",
			event:      feedStoredEvent(2, "bus-stale", generatedAt.Add(-2*time.Minute), `{}`),
			wantReason: TripDescriptorOmissionStaleTelemetry,
			wantEntity: true,
		},
		{
			name:       "suppressed stale",
			event:      feedStoredEvent(3, "bus-suppressed", generatedAt.Add(-10*time.Minute), `{}`),
			wantReason: TripDescriptorOmissionSuppressedStaleTelemetry,
			wantEntity: false,
		},
		{
			name:  "assignment telemetry mismatch",
			event: feedStoredEvent(4, "bus-mismatch", generatedAt.Add(-30*time.Second), `{}`),
			assignment: ptrAssignment(state.Assignment{
				AgencyID:         "demo-agency",
				VehicleID:        "bus-mismatch",
				TelemetryEventID: 999,
				State:            state.StateInService,
				RouteID:          "route-10",
				TripID:           "trip-10-0800",
				StartDate:        "20260420",
				StartTime:        "08:00:00",
				Confidence:       0.9,
				AssignmentSource: state.AssignmentSourceAutomatic,
				DegradedState:    state.DegradedNone,
			}),
			wantReason: TripDescriptorOmissionAssignmentTelemetryMismatch,
			wantEntity: true,
		},
		{
			name:  "below confidence",
			event: feedStoredEvent(5, "bus-low-confidence", generatedAt.Add(-30*time.Second), `{}`),
			assignment: ptrAssignment(state.Assignment{
				AgencyID:         "demo-agency",
				VehicleID:        "bus-low-confidence",
				TelemetryEventID: 5,
				State:            state.StateInService,
				RouteID:          "route-10",
				TripID:           "trip-10-0800",
				StartDate:        "20260420",
				StartTime:        "08:00:00",
				Confidence:       0.5,
				AssignmentSource: state.AssignmentSourceAutomatic,
				DegradedState:    state.DegradedNone,
			}),
			wantReason: TripDescriptorOmissionBelowPublicationConfidence,
			wantEntity: true,
		},
		{
			name:  "explicit unknown assignment",
			event: feedStoredEvent(6, "bus-unknown", generatedAt.Add(-30*time.Second), `{}`),
			assignment: ptrAssignment(state.Assignment{
				AgencyID:         "demo-agency",
				VehicleID:        "bus-unknown",
				TelemetryEventID: 6,
				State:            state.StateUnknown,
				Confidence:       0,
				AssignmentSource: state.AssignmentSourceUnknown,
				DegradedState:    state.DegradedUnknown,
			}),
			wantReason: TripDescriptorOmissionNotInService,
			wantEntity: true,
		},
		{
			name:  "degraded assignment",
			event: feedStoredEvent(7, "bus-degraded", generatedAt.Add(-30*time.Second), `{}`),
			assignment: ptrAssignment(state.Assignment{
				AgencyID:         "demo-agency",
				VehicleID:        "bus-degraded",
				TelemetryEventID: 7,
				State:            state.StateInService,
				RouteID:          "route-10",
				TripID:           "trip-10-0800",
				StartDate:        "20260420",
				StartTime:        "08:00:00",
				Confidence:       0.9,
				AssignmentSource: state.AssignmentSourceAutomatic,
				DegradedState:    state.DegradedAmbiguous,
			}),
			wantReason: TripDescriptorOmissionDegradedAssignment,
			wantEntity: true,
		},
		{
			name:  "manual override publishes below automatic threshold",
			event: feedStoredEvent(8, "bus-manual", generatedAt.Add(-30*time.Second), `{}`),
			assignment: ptrAssignment(state.Assignment{
				AgencyID:         "demo-agency",
				VehicleID:        "bus-manual",
				TelemetryEventID: 8,
				State:            state.StateInService,
				RouteID:          "route-10",
				TripID:           "trip-10-0800",
				StartDate:        "20260420",
				StartTime:        "08:00:00",
				Confidence:       0,
				AssignmentSource: state.AssignmentSourceManualOverride,
				DegradedState:    state.DegradedNone,
			}),
			wantReason: TripDescriptorOmissionNone,
			wantEntity: true,
			wantTrip:   true,
		},
		{
			name:  "matched",
			event: feedStoredEvent(9, "bus-match", generatedAt.Add(-30*time.Second), `{}`),
			assignment: ptrAssignment(state.Assignment{
				AgencyID:         "demo-agency",
				VehicleID:        "bus-match",
				TelemetryEventID: 9,
				State:            state.StateInService,
				RouteID:          "route-10",
				TripID:           "trip-10-0800",
				StartDate:        "20260420",
				StartTime:        "08:00:00",
				Confidence:       0.9,
				AssignmentSource: state.AssignmentSourceAutomatic,
				DegradedState:    state.DegradedNone,
			}),
			wantReason: TripDescriptorOmissionNone,
			wantEntity: true,
			wantTrip:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignments := map[string]state.Assignment{}
			if tt.assignment != nil {
				assignments[tt.event.VehicleID] = *tt.assignment
			}
			snapshot := buildTestSnapshot(t, generatedAt, []telemetry.StoredEvent{tt.event}, assignments, testConfig())
			if len(snapshot.Vehicles) != 1 {
				t.Fatalf("vehicle count = %d, want 1", len(snapshot.Vehicles))
			}
			vehicle := snapshot.Vehicles[0]
			if vehicle.TripDescriptorOmissionReason != tt.wantReason {
				t.Fatalf("reason = %s, want %s", vehicle.TripDescriptorOmissionReason, tt.wantReason)
			}
			if vehicle.IncludedInProtobuf != tt.wantEntity {
				t.Fatalf("included = %v, want %v", vehicle.IncludedInProtobuf, tt.wantEntity)
			}
			if vehicle.TripDescriptorPublished != tt.wantTrip {
				t.Fatalf("trip published = %v, want %v", vehicle.TripDescriptorPublished, tt.wantTrip)
			}
			debug := snapshot.Debug()
			if debug.Vehicles[0].TelemetryAgeSeconds != vehicle.TelemetryAgeSeconds {
				t.Fatalf("debug telemetry age = %f, want %f", debug.Vehicles[0].TelemetryAgeSeconds, vehicle.TelemetryAgeSeconds)
			}
		})
	}
}

func TestVehiclePositionsEmptyAndAllSuppressedFeedsAreSuccessful(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 0, 30, 0, time.UTC)

	empty := buildTestSnapshot(t, generatedAt, nil, nil, testConfig())
	emptyMessage := unmarshalFeed(t, mustMarshalProto(t, empty))
	if len(emptyMessage.Entity) != 0 || emptyMessage.GetHeader().GetGtfsRealtimeVersion() != GTFSRealtimeVersion {
		t.Fatalf("empty message = %+v, want populated header and no entities", emptyMessage)
	}
	if !empty.Debug().NoTelemetry {
		t.Fatalf("empty debug no_telemetry = false, want true")
	}

	suppressedEvent := feedStoredEvent(1, "bus-suppressed", generatedAt.Add(-10*time.Minute), `{}`)
	allSuppressed := buildTestSnapshot(t, generatedAt, []telemetry.StoredEvent{suppressedEvent}, nil, testConfig())
	suppressedMessage := unmarshalFeed(t, mustMarshalProto(t, allSuppressed))
	if len(suppressedMessage.Entity) != 0 || suppressedMessage.GetHeader().GetGtfsRealtimeVersion() != GTFSRealtimeVersion {
		t.Fatalf("suppressed message = %+v, want populated header and no entities", suppressedMessage)
	}
	if allSuppressed.Debug().NoTelemetry {
		t.Fatalf("all-suppressed debug no_telemetry = true, want false")
	}
}

func TestVehiclePositionsTruncatesBeforePublicationAndSortsOutput(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 0, 30, 0, time.UTC)
	config := testConfig()
	config.MaxVehicles = 2
	events := []telemetry.StoredEvent{
		feedStoredEvent(1, "z-bus", generatedAt.Add(-10*time.Second), `{}`),
		feedStoredEvent(2, "a-bus", generatedAt.Add(-20*time.Second), `{}`),
		feedStoredEvent(3, "m-bus", generatedAt.Add(-30*time.Second), `{}`),
	}
	snapshot := buildTestSnapshot(t, generatedAt, events, nil, config)
	if !snapshot.Truncated || snapshot.VehiclesInSnapshot != 2 || snapshot.LatestTelemetryRowsRead != 3 || snapshot.TruncatedVehicleCountMin != 1 {
		t.Fatalf("truncation fields = %+v, want truncated 3->2", snapshot)
	}
	if snapshot.Vehicles[0].VehicleID != "a-bus" || snapshot.Vehicles[1].VehicleID != "z-bus" {
		t.Fatalf("vehicle order = %s, %s; want deterministic vehicle id sort after cap", snapshot.Vehicles[0].VehicleID, snapshot.Vehicles[1].VehicleID)
	}
	message := unmarshalFeed(t, mustMarshalProto(t, snapshot))
	if got := []string{message.Entity[0].GetId(), message.Entity[1].GetId()}; got[0] != "a-bus" || got[1] != "z-bus" {
		t.Fatalf("protobuf entity order = %+v, want sorted capped vehicles", got)
	}
}

func TestVehiclePositionsFrequencyAndPayloadOptionalFields(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 0, 30, 0, time.UTC)
	frequencyEvent := feedStoredEvent(1, "freq-bus", generatedAt.Add(-30*time.Second), `{"bearing":0}`)
	frequencyAssignment := matchedAssignment(frequencyEvent, "trip-loop")
	frequencyAssignment.ReasonCodes = []string{state.ReasonFrequencyNonExact}
	snapshot := buildTestSnapshot(t, generatedAt, []telemetry.StoredEvent{frequencyEvent}, map[string]state.Assignment{"freq-bus": frequencyAssignment}, testConfig())
	message := unmarshalFeed(t, mustMarshalProto(t, snapshot))
	vehicle := message.Entity[0].GetVehicle()
	if vehicle.GetTrip().GetScheduleRelationship() != gtfsrt.TripDescriptor_UNSCHEDULED {
		t.Fatalf("relationship = %s, want UNSCHEDULED", vehicle.GetTrip().GetScheduleRelationship())
	}
	if vehicle.GetPosition().Bearing == nil || vehicle.GetPosition().GetBearing() != 0 {
		t.Fatalf("bearing = %v/%f, want explicit zero", vehicle.GetPosition().Bearing, vehicle.GetPosition().GetBearing())
	}
	if vehicle.GetPosition().Speed != nil {
		t.Fatalf("speed was set from missing payload field: %v", vehicle.GetPosition().Speed)
	}

	malformedEvent := feedStoredEvent(2, "bad-bearing", generatedAt.Add(-30*time.Second), `{"bearing":"north","speed_mps":null}`)
	malformedAssignment := matchedAssignment(malformedEvent, "trip-10-0800")
	snapshot = buildTestSnapshot(t, generatedAt, []telemetry.StoredEvent{malformedEvent}, map[string]state.Assignment{"bad-bearing": malformedAssignment}, testConfig())
	message = unmarshalFeed(t, mustMarshalProto(t, snapshot))
	if message.Entity[0].GetVehicle().GetPosition().Bearing != nil {
		t.Fatalf("bearing was set from malformed payload: %v", message.Entity[0].GetVehicle().GetPosition().Bearing)
	}
	if message.Entity[0].GetVehicle().GetPosition().Speed != nil {
		t.Fatalf("speed was set from null payload: %v", message.Entity[0].GetVehicle().GetPosition().Speed)
	}
}

func TestVehiclePositionsDeterministicProtoBytes(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 0, 30, 0, time.UTC)
	event := feedStoredEvent(1, "bus-10", generatedAt.Add(-30*time.Second), `{}`)
	assignment := matchedAssignment(event, "trip-10-0800")
	snapshot := buildTestSnapshot(t, generatedAt, []telemetry.StoredEvent{event}, map[string]state.Assignment{"bus-10": assignment}, testConfig())

	first := mustMarshalProto(t, snapshot)
	second := mustMarshalProto(t, snapshot)
	if !bytes.Equal(first, second) {
		t.Fatalf("protobuf bytes are not deterministic for identical snapshot")
	}
}

func buildTestSnapshot(t *testing.T, generatedAt time.Time, events []telemetry.StoredEvent, assignments map[string]state.Assignment, config VehiclePositionsConfig) VehiclePositionsSnapshot {
	t.Helper()
	builder, err := NewVehiclePositionsBuilder(&fakeTelemetryRepo{events: events}, &fakeStateRepo{assignments: assignments}, config)
	if err != nil {
		t.Fatalf("new builder: %v", err)
	}
	snapshot, err := builder.Snapshot(context.Background(), generatedAt)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	return snapshot
}

func testConfig() VehiclePositionsConfig {
	return VehiclePositionsConfig{
		AgencyID:                  "demo-agency",
		MaxVehicles:               10,
		StaleTelemetryTTL:         90 * time.Second,
		SuppressStaleVehicleAfter: 300 * time.Second,
		TripConfidenceThreshold:   0.65,
	}
}

func feedStoredEvent(id int64, vehicleID string, observedAt time.Time, payload string) telemetry.StoredEvent {
	return telemetry.StoredEvent{
		ID: id,
		Event: telemetry.Event{
			AgencyID:  "demo-agency",
			DeviceID:  "device-" + vehicleID,
			VehicleID: vehicleID,
			Timestamp: observedAt,
			Lat:       49.2827,
			Lon:       -123.1207,
			Bearing:   0,
			SpeedMPS:  8.2,
		},
		ReceivedAt:   observedAt.Add(5 * time.Second),
		IngestStatus: telemetry.IngestStatusAccepted,
		PayloadJSON:  json.RawMessage(payload),
	}
}

func matchedAssignment(event telemetry.StoredEvent, tripID string) state.Assignment {
	return state.Assignment{
		AgencyID:         event.AgencyID,
		VehicleID:        event.VehicleID,
		TelemetryEventID: event.ID,
		State:            state.StateInService,
		RouteID:          "route-10",
		TripID:           tripID,
		StartDate:        "20260420",
		StartTime:        "08:00:00",
		Confidence:       0.9,
		AssignmentSource: state.AssignmentSourceAutomatic,
		ReasonCodes:      []string{state.ReasonTripHintMatch},
		DegradedState:    state.DegradedNone,
	}
}

func ptrAssignment(assignment state.Assignment) *state.Assignment {
	return &assignment
}

func mustMarshalProto(t *testing.T, snapshot VehiclePositionsSnapshot) []byte {
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
	return &message
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

func (f *fakeTelemetryRepo) ListLatestByAgency(_ context.Context, _ string, limit int) ([]telemetry.StoredEvent, error) {
	if f.err != nil {
		return nil, f.err
	}
	if limit > len(f.events) {
		limit = len(f.events)
	}
	return append([]telemetry.StoredEvent(nil), f.events[:limit]...), nil
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
	return state.Assignment{}, errors.New("not implemented")
}
