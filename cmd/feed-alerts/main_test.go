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

	"open-transit-rt/internal/feed/alerts"
)

func TestAlertsHandlersReturnValidEmptyFeedAndDebug(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	handler := newHandler(fakeAlertsBuilder{snapshot: alerts.Snapshot{
		AgencyID:    "demo-agency",
		GeneratedAt: generatedAt,
		Status:      alerts.StatusDeferred,
		Reason:      alerts.ReasonAlertsAuthoringMissing,
	}}, okPinger{})

	req := httptest.NewRequest(http.MethodGet, "/public/gtfsrt/alerts.pb", nil)
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

	req = httptest.NewRequest(http.MethodGet, "/public/gtfsrt/alerts.json", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("json status = %d, want 200", rr.Code)
	}
	if rr.Header().Get("Last-Modified") != generatedAt.Format(http.TimeFormat) {
		t.Fatalf("json last-modified = %q, want snapshot generated_at", rr.Header().Get("Last-Modified"))
	}
	var debug alerts.Debug
	if err := json.Unmarshal(rr.Body.Bytes(), &debug); err != nil {
		t.Fatalf("decode debug: %v", err)
	}
	if debug.Status != alerts.StatusDeferred || debug.Diagnostics != "json_only" {
		t.Fatalf("debug = %+v, want deferred JSON-only diagnostics", debug)
	}
}

func TestAlertsHandlersRejectWrongMethod(t *testing.T) {
	handler := newHandler(fakeAlertsBuilder{snapshot: alerts.Snapshot{GeneratedAt: time.Now().UTC()}}, okPinger{})
	req := httptest.NewRequest(http.MethodPost, "/public/gtfsrt/alerts.pb", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST status = %d, want 405", rr.Code)
	}
}

func TestAlertsReadyz(t *testing.T) {
	handler := newHandler(fakeAlertsBuilder{}, errPinger{err: errors.New("down")})
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rr.Code)
	}

	handler = newHandler(fakeAlertsBuilder{}, okPinger{})
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
}

type fakeAlertsBuilder struct {
	snapshot alerts.Snapshot
}

func (f fakeAlertsBuilder) Snapshot(time.Time) alerts.Snapshot {
	return f.snapshot
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
