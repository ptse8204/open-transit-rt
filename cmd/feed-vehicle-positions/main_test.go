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
	snapshot feed.VehiclePositionsSnapshot
	err      error
}

func (f *fakeSnapshotBuilder) Snapshot(context.Context, time.Time) (feed.VehiclePositionsSnapshot, error) {
	if f.err != nil {
		return feed.VehiclePositionsSnapshot{}, f.err
	}
	return f.snapshot, nil
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
