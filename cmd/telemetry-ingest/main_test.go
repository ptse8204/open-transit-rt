package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"open-transit-rt/internal/telemetry"
)

type fakeRepo struct {
	storeFn func(context.Context, telemetry.Event, json.RawMessage) (telemetry.StoreResult, error)
	listFn  func(context.Context, string, int) ([]telemetry.StoredEvent, error)
}

func (f fakeRepo) Store(ctx context.Context, event telemetry.Event, payload json.RawMessage) (telemetry.StoreResult, error) {
	return f.storeFn(ctx, event, payload)
}

func (f fakeRepo) LatestByVehicle(context.Context, string, string) (telemetry.StoredEvent, error) {
	return telemetry.StoredEvent{}, nil
}

func (f fakeRepo) ListLatestByAgency(context.Context, string, int) ([]telemetry.StoredEvent, error) {
	return nil, nil
}

func (f fakeRepo) ListEvents(ctx context.Context, agencyID string, limit int) ([]telemetry.StoredEvent, error) {
	return f.listFn(ctx, agencyID, limit)
}

type fakePinger struct {
	err error
}

func (f fakePinger) Ping(context.Context) error {
	return f.err
}

func TestTelemetryPostStatuses(t *testing.T) {
	for _, tc := range []struct {
		name         string
		status       telemetry.IngestStatus
		wantHTTPCode int
		wantAccepted bool
	}{
		{name: "accepted", status: telemetry.IngestStatusAccepted, wantHTTPCode: http.StatusCreated, wantAccepted: true},
		{name: "duplicate", status: telemetry.IngestStatusDuplicate, wantHTTPCode: http.StatusAccepted, wantAccepted: false},
		{name: "out_of_order", status: telemetry.IngestStatusOutOfOrder, wantHTTPCode: http.StatusAccepted, wantAccepted: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			receivedAt := time.Date(2026, 4, 20, 16, 0, 0, 0, time.UTC)
			handler := newHandler(fakeRepo{
				storeFn: func(_ context.Context, event telemetry.Event, _ json.RawMessage) (telemetry.StoreResult, error) {
					return telemetry.StoreResult{StoredEvent: telemetry.StoredEvent{
						Event:        event,
						ReceivedAt:   receivedAt,
						IngestStatus: tc.status,
					}}, nil
				},
				listFn: func(context.Context, string, int) ([]telemetry.StoredEvent, error) { return nil, nil },
			}, fakePinger{})

			req := httptest.NewRequest(http.MethodPost, "/v1/telemetry", strings.NewReader(validTelemetryPayload()))
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantHTTPCode {
				t.Fatalf("status code = %d, want %d; body=%s", rec.Code, tc.wantHTTPCode, rec.Body.String())
			}
			var response ingestResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if response.Accepted != tc.wantAccepted {
				t.Fatalf("accepted = %v, want %v", response.Accepted, tc.wantAccepted)
			}
			if !response.ReceivedAt.Equal(receivedAt) {
				t.Fatalf("received_at = %s, want DB timestamp %s", response.ReceivedAt, receivedAt)
			}
		})
	}
}

func TestTelemetryPostRejectsInvalidRequests(t *testing.T) {
	storeCalled := false
	handler := newHandler(fakeRepo{
		storeFn: func(context.Context, telemetry.Event, json.RawMessage) (telemetry.StoreResult, error) {
			storeCalled = true
			return telemetry.StoreResult{}, nil
		},
		listFn: func(context.Context, string, int) ([]telemetry.StoredEvent, error) { return nil, nil },
	}, fakePinger{})

	for _, tc := range []struct {
		name string
		body string
	}{
		{name: "invalid json", body: `{"agency_id":`},
		{name: "invalid payload", body: `{"agency_id":"demo-agency","device_id":"device-1","vehicle_id":"bus-1","timestamp":"2026-04-20T15:02:00Z","lat":99,"lon":-123.1}`},
		{name: "timezone-less timestamp", body: `{"agency_id":"demo-agency","device_id":"device-1","vehicle_id":"bus-1","timestamp":"2026-04-20T15:02:00","lat":49.2,"lon":-123.1}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			storeCalled = false
			req := httptest.NewRequest(http.MethodPost, "/v1/telemetry", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status code = %d, want 400; body=%s", rec.Code, rec.Body.String())
			}
			if storeCalled {
				t.Fatalf("store called for invalid request")
			}
		})
	}
}

func TestTelemetryPostUnknownAgency(t *testing.T) {
	handler := newHandler(fakeRepo{
		storeFn: func(context.Context, telemetry.Event, json.RawMessage) (telemetry.StoreResult, error) {
			return telemetry.StoreResult{}, telemetry.ErrUnknownAgency
		},
		listFn: func(context.Context, string, int) ([]telemetry.StoredEvent, error) { return nil, nil },
	}, fakePinger{})

	req := httptest.NewRequest(http.MethodPost, "/v1/telemetry", strings.NewReader(validTelemetryPayload()))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status code = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
}

func TestEventsEndpointRequiresAgencyAndBoundedLimit(t *testing.T) {
	var gotAgency string
	var gotLimit int
	handler := newHandler(fakeRepo{
		storeFn: func(context.Context, telemetry.Event, json.RawMessage) (telemetry.StoreResult, error) {
			return telemetry.StoreResult{}, nil
		},
		listFn: func(_ context.Context, agencyID string, limit int) ([]telemetry.StoredEvent, error) {
			gotAgency = agencyID
			gotLimit = limit
			return []telemetry.StoredEvent{{
				ID: 1,
				Event: telemetry.Event{
					AgencyID:  agencyID,
					DeviceID:  "device-1",
					VehicleID: "bus-1",
					Timestamp: time.Date(2026, 4, 20, 15, 2, 0, 0, time.UTC),
					Lat:       49.2,
					Lon:       -123.1,
				},
				ReceivedAt:   time.Date(2026, 4, 20, 15, 2, 1, 0, time.UTC),
				IngestStatus: telemetry.IngestStatusOutOfOrder,
			}}, nil
		},
	}, fakePinger{})

	for _, path := range []string{"/v1/events", "/v1/events?agency_id=demo-agency", "/v1/events?agency_id=demo-agency&limit=501"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("%s status code = %d, want 400", path, rec.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/events?agency_id=demo-agency&limit=25", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if gotAgency != "demo-agency" || gotLimit != 25 {
		t.Fatalf("repo called with agency=%q limit=%d", gotAgency, gotLimit)
	}
}

func TestReadyzReportsDBUnavailable(t *testing.T) {
	handler := newHandler(fakeRepo{}, fakePinger{err: errors.New("db down")})

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status code = %d, want 503", rec.Code)
	}
}

func validTelemetryPayload() string {
	return `{
		"agency_id":"demo-agency",
		"device_id":"device-1",
		"vehicle_id":"bus-1",
		"timestamp":"2026-04-20T15:02:00Z",
		"lat":49.2,
		"lon":-123.1
	}`
}
