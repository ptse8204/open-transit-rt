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
	"open-transit-rt/internal/feed/tripupdates"
	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/prediction"
)

func TestLoadTripUpdatesConfigDerivesVehiclePositionsURLFromPublicFeedBase(t *testing.T) {
	t.Setenv("AGENCY_ID", "demo-agency")
	t.Setenv("FEED_BASE_URL", "http://localhost:8083/public/")
	t.Setenv("VEHICLE_POSITIONS_FEED_URL", "")
	config, err := loadTripUpdatesConfigFromEnv()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if config.VehiclePositionsURL != "http://localhost:8083/public/gtfsrt/vehicle_positions.pb" {
		t.Fatalf("vehicle positions url = %q, want Phase 3 public protobuf path", config.VehiclePositionsURL)
	}
}

func TestLoadTripUpdatesConfigRejectsAmbiguousFeedBase(t *testing.T) {
	t.Setenv("AGENCY_ID", "demo-agency")
	t.Setenv("FEED_BASE_URL", "http://localhost:8083")
	t.Setenv("VEHICLE_POSITIONS_FEED_URL", "")
	if _, err := loadTripUpdatesConfigFromEnv(); err == nil {
		t.Fatalf("load config succeeded, want FEED_BASE_URL missing /public error")
	}
}

func TestVehiclePositionsFeedURLWinsAsExactFullURL(t *testing.T) {
	t.Setenv("AGENCY_ID", "demo-agency")
	t.Setenv("FEED_BASE_URL", "http://bad.example/public")
	t.Setenv("VEHICLE_POSITIONS_FEED_URL", "https://feeds.example.com/public/gtfsrt/vehicle_positions.pb")
	config, err := loadTripUpdatesConfigFromEnv()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if config.VehiclePositionsURL != "https://feeds.example.com/public/gtfsrt/vehicle_positions.pb" {
		t.Fatalf("vehicle positions url = %q, want exact env URL", config.VehiclePositionsURL)
	}
}

func TestPredictionAdapterFromEnvDefaultsToDeterministicAndKeepsNoopFallback(t *testing.T) {
	t.Setenv("TRIP_UPDATES_ADAPTER", "")
	adapter, err := predictionAdapterFromEnv(fakeAdapterSchedule{}, nil)
	if err != nil {
		t.Fatalf("default adapter: %v", err)
	}
	if adapter.Name() != "deterministic" {
		t.Fatalf("adapter = %q, want deterministic default", adapter.Name())
	}

	t.Setenv("TRIP_UPDATES_ADAPTER", "noop")
	adapter, err = predictionAdapterFromEnv(fakeAdapterSchedule{}, nil)
	if err != nil {
		t.Fatalf("noop adapter: %v", err)
	}
	if adapter.Name() != "noop" {
		t.Fatalf("adapter = %q, want noop fallback", adapter.Name())
	}

	t.Setenv("TRIP_UPDATES_ADAPTER", "bad")
	if _, err := predictionAdapterFromEnv(fakeAdapterSchedule{}, nil); err == nil {
		t.Fatalf("invalid adapter succeeded, want error")
	}
}

func TestTripUpdatesHandlersReturnValidEmptyFeedAndDebug(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	snapshot := tripupdates.Snapshot{
		AgencyID:                      "demo-agency",
		GeneratedAt:                   generatedAt,
		AdapterName:                   "noop",
		VehiclePositionsURL:           "http://localhost:8083/public/gtfsrt/vehicle_positions.pb",
		Diagnostics:                   prediction.Diagnostics{Status: prediction.StatusNoop, Reason: prediction.ReasonNoopAdapter},
		DiagnosticsPersistenceOutcome: tripupdates.DiagnosticsPersistenceStored,
	}
	handler := newHandler(&fakeTripUpdatesBuilder{snapshot: snapshot}, okPinger{})

	req := httptest.NewRequest(http.MethodGet, "/public/gtfsrt/trip_updates.pb", nil)
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
		t.Fatalf("unmarshal feed: %v", err)
	}
	if message.GetHeader().GetTimestamp() != uint64(generatedAt.Unix()) || len(message.Entity) != 0 {
		t.Fatalf("message = %+v, want empty feed with generated_at timestamp", &message)
	}

	req = httptest.NewRequest(http.MethodGet, "/public/gtfsrt/trip_updates.json", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("json status = %d, want 200", rr.Code)
	}
	if rr.Header().Get("Last-Modified") != generatedAt.Format(http.TimeFormat) {
		t.Fatalf("json last-modified = %q, want snapshot generated_at", rr.Header().Get("Last-Modified"))
	}
	var debug tripupdates.Debug
	if err := json.Unmarshal(rr.Body.Bytes(), &debug); err != nil {
		t.Fatalf("decode debug: %v", err)
	}
	if debug.DiagnosticsStatus != prediction.StatusNoop || debug.DiagnosticsPersistenceOutcome != tripupdates.DiagnosticsPersistenceStored {
		t.Fatalf("debug = %+v, want no-op diagnostics", debug)
	}
}

func TestTripUpdatesDebugRejectsUnauthenticatedAccess(t *testing.T) {
	handler := newHandlerWithAuth(&fakeTripUpdatesBuilder{}, okPinger{}, authRejectAll{})
	req := httptest.NewRequest(http.MethodGet, "/public/gtfsrt/trip_updates.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rr.Code)
	}
}

func TestTripUpdatesHandlersRejectWrongMethodAndSurfaceSnapshotErrors(t *testing.T) {
	handler := newHandler(&fakeTripUpdatesBuilder{err: errors.New("database down")}, okPinger{})
	req := httptest.NewRequest(http.MethodPost, "/public/gtfsrt/trip_updates.pb", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST status = %d, want 405", rr.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/public/gtfsrt/trip_updates.pb", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("error status = %d, want 500", rr.Code)
	}
}

func TestTripUpdatesReadyz(t *testing.T) {
	handler := newHandler(&fakeTripUpdatesBuilder{}, errPinger{err: errors.New("down")})
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rr.Code)
	}

	handler = newHandler(&fakeTripUpdatesBuilder{}, okPinger{})
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
}

type fakeTripUpdatesBuilder struct {
	snapshot tripupdates.Snapshot
	err      error
}

func (f *fakeTripUpdatesBuilder) Snapshot(context.Context, time.Time) (tripupdates.Snapshot, error) {
	if f.err != nil {
		return tripupdates.Snapshot{}, f.err
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

func (e errPinger) Ping(context.Context) error {
	return e.err
}

type authRejectAll struct{}

func (authRejectAll) Require(...auth.Role) func(http.Handler) http.Handler {
	return func(_ http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		})
	}
}

type fakeAdapterSchedule struct{}

func (fakeAdapterSchedule) Agency(context.Context, string) (gtfs.Agency, error) {
	return gtfs.Agency{ID: "demo-agency", Timezone: "America/Vancouver"}, nil
}

func (fakeAdapterSchedule) ActiveFeedVersion(context.Context, string) (gtfs.FeedVersion, error) {
	return gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"}, nil
}

func (fakeAdapterSchedule) ListTripCandidates(context.Context, string, string, string) ([]gtfs.TripCandidate, error) {
	return nil, nil
}
