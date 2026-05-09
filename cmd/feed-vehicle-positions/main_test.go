package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/feed"
	"open-transit-rt/internal/telemetry"
)

func TestLoadFeedConfigRequiresAgencyID(t *testing.T) {
	t.Setenv("AGENCY_ID", "")
	_, err := loadFeedConfigFromEnv()
	if err == nil {
		t.Fatalf("loadFeedConfigFromEnv succeeded, want missing AGENCY_ID error")
	}
}

func TestVehiclePositionsProtobufHandlerHeadersAndEmptyFeed(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 0, 30, 0, time.UTC)
	handler := newHandler(&fakeSnapshotBuilder{
		snapshot: feed.VehiclePositionsSnapshot{
			AgencyID:     "demo-agency",
			GeneratedAt:  generatedAt,
			VehicleLimit: 2000,
			NoTelemetry:  true,
		},
	}, okPinger{})

	req := httptest.NewRequest(http.MethodGet, "/public/gtfsrt/vehicle_positions.pb", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if rr.Header().Get("Content-Type") != "application/x-protobuf" {
		t.Fatalf("content-type = %q, want application/x-protobuf", rr.Header().Get("Content-Type"))
	}
	if rr.Header().Get("Last-Modified") != generatedAt.Format(http.TimeFormat) {
		t.Fatalf("last-modified = %q, want snapshot generated_at", rr.Header().Get("Last-Modified"))
	}
	var message gtfsrt.FeedMessage
	if err := proto.Unmarshal(rr.Body.Bytes(), &message); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if message.GetHeader().GetGtfsRealtimeVersion() != feed.GTFSRealtimeVersion || len(message.Entity) != 0 {
		t.Fatalf("message = %+v, want valid empty feed", &message)
	}
}

func TestVehiclePositionsPathRoutedPublicFeedBuildsRequestedAgencyOnly(t *testing.T) {
	generatedAt := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	builder := &fakeSnapshotBuilder{snapshotsByAgency: map[string]feed.VehiclePositionsSnapshot{
		"agency-a": {AgencyID: "agency-a", GeneratedAt: generatedAt, NoTelemetry: true},
		"agency-b": {AgencyID: "agency-b", GeneratedAt: generatedAt, NoTelemetry: true},
	}}
	handler := newHandler(builder, okPinger{})

	req := httptest.NewRequest(http.MethodGet, "/public/agencies/agency-a/gtfsrt/vehicle_positions.pb", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("agency-a status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if builder.agencyCalls["agency-a"] != 1 || builder.agencyCalls["agency-b"] != 0 || builder.defaultCalls != 0 {
		t.Fatalf("calls = default %d agency %+v, want only agency-a", builder.defaultCalls, builder.agencyCalls)
	}

	req = httptest.NewRequest(http.MethodGet, "/public/agencies/agency-b/gtfsrt/vehicle_positions.pb", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("agency-b status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if builder.agencyCalls["agency-a"] != 1 || builder.agencyCalls["agency-b"] != 1 || builder.defaultCalls != 0 {
		t.Fatalf("calls = default %d agency %+v, want one call per requested agency", builder.defaultCalls, builder.agencyCalls)
	}

	req = httptest.NewRequest(http.MethodGet, "/public/agencies/agency-a/gtfsrt/vehicle_positions.json", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("per-agency debug status = %d, want 404", rr.Code)
	}
}

func TestVehiclePositionsJSONHandlerUsesSnapshotDebug(t *testing.T) {
	generatedAt := time.Date(2026, 4, 20, 15, 0, 30, 0, time.UTC)
	snapshot := feed.VehiclePositionsSnapshot{
		AgencyID:           "demo-agency",
		GeneratedAt:        generatedAt,
		VehicleLimit:       2000,
		VehiclesInSnapshot: 1,
		Vehicles: []feed.VehicleSnapshot{{
			VehicleID:                    "bus-10",
			TelemetryEvent:               telemetry.StoredEvent{ID: 7, Event: telemetry.Event{VehicleID: "bus-10", Timestamp: generatedAt.Add(-30 * time.Second), Lat: 49.2, Lon: -123.1}},
			TelemetryAgeSeconds:          30,
			IncludedInProtobuf:           true,
			TripDescriptorOmissionReason: feed.TripDescriptorOmissionNoAssignment,
		}},
	}
	handler := newHandler(&fakeSnapshotBuilder{snapshot: snapshot}, okPinger{})

	req := httptest.NewRequest(http.MethodGet, "/public/gtfsrt/vehicle_positions.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if rr.Header().Get("Last-Modified") != generatedAt.Format(http.TimeFormat) {
		t.Fatalf("last-modified = %q, want snapshot generated_at", rr.Header().Get("Last-Modified"))
	}
	var payload struct {
		NoTelemetry bool `json:"no_telemetry"`
		Vehicles    []struct {
			VehicleID                    string  `json:"vehicle_id"`
			TelemetryAgeSeconds          float64 `json:"telemetry_age_seconds"`
			IncludedInProtobuf           bool    `json:"included_in_protobuf"`
			TripDescriptorOmissionReason string  `json:"trip_descriptor_omission_reason"`
		} `json:"vehicles"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	if payload.NoTelemetry || len(payload.Vehicles) != 1 || payload.Vehicles[0].TelemetryAgeSeconds != 30 || payload.Vehicles[0].TripDescriptorOmissionReason != feed.TripDescriptorOmissionNoAssignment {
		t.Fatalf("payload = %+v, want snapshot debug fields", payload)
	}
}

func TestVehiclePositionsDebugRejectsUnauthenticatedAccess(t *testing.T) {
	handler := newHandlerWithAuth(&fakeSnapshotBuilder{}, okPinger{}, authRejectAll{})
	req := httptest.NewRequest(http.MethodGet, "/public/gtfsrt/vehicle_positions.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rr.Code)
	}
}

func TestVehiclePositionsDebugRejectsCrossAgencyPrincipal(t *testing.T) {
	handler := newHandlerWithAuth(&fakeSnapshotBuilder{snapshot: feed.VehiclePositionsSnapshot{
		AgencyID:    "agency-a",
		GeneratedAt: time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC),
	}}, okPinger{}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject:  "admin-b@example.com",
		AgencyID: "agency-b",
		Roles:    []auth.Role{auth.RoleReadOnly},
		Method:   auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/debug/gtfsrt/vehicle_positions.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rr.Code)
	}
}

func TestVehiclePositionsHandlersRejectWrongMethodAndSurfaceSnapshotErrors(t *testing.T) {
	handler := newHandler(&fakeSnapshotBuilder{err: errors.New("database down")}, okPinger{})

	req := httptest.NewRequest(http.MethodPost, "/public/gtfsrt/vehicle_positions.pb", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST status = %d, want 405", rr.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/public/gtfsrt/vehicle_positions.pb", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("error status = %d, want 500", rr.Code)
	}
}

func TestVehiclePositionsHealthPersistenceSuccess(t *testing.T) {
	generatedAt := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	health := &fakeHealthRecorder{done: make(chan feed.VehiclePositionsHealthRecord, 1)}
	handler := newHandlerWithAuthAndHealth(&fakeSnapshotBuilder{snapshot: feed.VehiclePositionsSnapshot{
		AgencyID:           "demo-agency",
		GeneratedAt:        generatedAt,
		VehicleLimit:       2000,
		VehiclesInSnapshot: 1,
		Vehicles: []feed.VehicleSnapshot{{
			VehicleID:                    "bus-1",
			TelemetryEvent:               telemetry.StoredEvent{Event: telemetry.Event{VehicleID: "bus-1", Timestamp: generatedAt.Add(-10 * time.Second), Lat: 1, Lon: 2}},
			IncludedInProtobuf:           true,
			TripDescriptorPublished:      true,
			TripDescriptorOmissionReason: feed.TripDescriptorOmissionNone,
		}},
	}}, okPinger{}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}}, health)
	req := httptest.NewRequest(http.MethodGet, "/public/gtfsrt/vehicle_positions.pb", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	select {
	case record := <-health.done:
		if record.AgencyID != "demo-agency" || !record.EndpointAvailable || record.VehiclesInSnapshot != 1 || record.TripDescriptors != 1 {
			t.Fatalf("record = %+v, want bounded vehicle positions health", record)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for health persistence")
	}
}

func TestVehiclePositionsHealthPersistenceFailureDoesNotChangePublicStatus(t *testing.T) {
	generatedAt := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	handler := newHandlerWithAuthAndHealth(&fakeSnapshotBuilder{snapshot: feed.VehiclePositionsSnapshot{
		AgencyID:    "demo-agency",
		GeneratedAt: generatedAt,
		NoTelemetry: true,
	}}, okPinger{}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}}, &fakeHealthRecorder{err: errors.New("insert failed"), done: make(chan feed.VehiclePositionsHealthRecord, 1)})
	req := httptest.NewRequest(http.MethodGet, "/public/gtfsrt/vehicle_positions.pb", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 despite health persistence failure: %s", rr.Code, rr.Body.String())
	}
	var message gtfsrt.FeedMessage
	if err := proto.Unmarshal(rr.Body.Bytes(), &message); err != nil {
		t.Fatalf("response was not valid protobuf: %v", err)
	}
}

func TestVehiclePositionsReadyz(t *testing.T) {
	handler := newHandler(&fakeSnapshotBuilder{}, errPinger{err: errors.New("down")})
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rr.Code)
	}

	handler = newHandler(&fakeSnapshotBuilder{}, okPinger{})
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
}

type fakeSnapshotBuilder struct {
	snapshot          feed.VehiclePositionsSnapshot
	snapshotsByAgency map[string]feed.VehiclePositionsSnapshot
	err               error
	defaultCalls      int
	agencyCalls       map[string]int
}

type fakeHealthRecorder struct {
	err  error
	done chan feed.VehiclePositionsHealthRecord
}

func (f *fakeHealthRecorder) SaveVehiclePositionsHealth(_ context.Context, record feed.VehiclePositionsHealthRecord) error {
	if f.done != nil {
		f.done <- record
	}
	return f.err
}

func (f *fakeSnapshotBuilder) Snapshot(context.Context, time.Time) (feed.VehiclePositionsSnapshot, error) {
	f.defaultCalls++
	if f.err != nil {
		return feed.VehiclePositionsSnapshot{}, f.err
	}
	return f.snapshot, nil
}

func (f *fakeSnapshotBuilder) SnapshotForAgency(_ context.Context, agencyID string, _ time.Time) (feed.VehiclePositionsSnapshot, error) {
	if f.agencyCalls == nil {
		f.agencyCalls = make(map[string]int)
	}
	f.agencyCalls[agencyID]++
	if f.err != nil {
		return feed.VehiclePositionsSnapshot{}, f.err
	}
	if f.snapshotsByAgency != nil {
		return f.snapshotsByAgency[agencyID], nil
	}
	snapshot := f.snapshot
	snapshot.AgencyID = agencyID
	return snapshot, nil
}

type okPinger struct{}

func (okPinger) Ping(context.Context) error {
	return nil
}

type errPinger struct {
	err error
}

type authRejectAll struct{}

func (authRejectAll) Require(...auth.Role) func(http.Handler) http.Handler {
	return func(_ http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		})
	}
}

func (e errPinger) Ping(context.Context) error {
	return e.err
}
