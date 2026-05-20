package prediction

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

func TestAdapterFactoryDefaultsAndExternalModes(t *testing.T) {
	adapter, err := NewAdapterFromEnv(testEnv(nil), &fakePredictionSchedule{}, nil)
	if err != nil {
		t.Fatalf("default adapter: %v", err)
	}
	if adapter.Name() != AdapterNameDeterministic {
		t.Fatalf("adapter = %q, want deterministic", adapter.Name())
	}

	adapter, err = NewAdapterFromEnv(testEnv(map[string]string{"TRIP_UPDATES_ADAPTER": "noop"}), &fakePredictionSchedule{}, nil)
	if err != nil {
		t.Fatalf("noop adapter: %v", err)
	}
	if adapter.Name() != AdapterNameNoop {
		t.Fatalf("adapter = %q, want noop", adapter.Name())
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	endpoint := server.URL + externalHTTPPath
	host := mustURLHost(t, endpoint)
	env := testEnv(map[string]string{
		"TRIP_UPDATES_ADAPTER":                       "external_http",
		"TRIP_UPDATES_EXTERNAL_HTTP_URL":             endpoint,
		"TRIP_UPDATES_EXTERNAL_HTTP_ALLOWED_HOSTS":   host,
		"TRIP_UPDATES_EXTERNAL_HTTP_TIMEOUT_SECONDS": "1",
	})
	adapter, err = NewAdapterFromEnv(env, &fakePredictionSchedule{}, nil)
	if err != nil {
		t.Fatalf("external_http adapter: %v", err)
	}
	if adapter.Name() != AdapterNameExternalHTTP {
		t.Fatalf("adapter = %q, want external_http", adapter.Name())
	}

	env = testEnv(map[string]string{
		"TRIP_UPDATES_ADAPTER":                       "external_http_shadow",
		"TRIP_UPDATES_EXTERNAL_HTTP_URL":             endpoint,
		"TRIP_UPDATES_EXTERNAL_HTTP_ALLOWED_HOSTS":   host,
		"TRIP_UPDATES_EXTERNAL_HTTP_TIMEOUT_SECONDS": "1",
	})
	adapter, err = NewAdapterFromEnv(env, &fakePredictionSchedule{}, nil)
	if err != nil {
		t.Fatalf("external_http_shadow adapter: %v", err)
	}
	if adapter.Name() != AdapterNameExternalHTTPShadow {
		t.Fatalf("adapter = %q, want external_http_shadow", adapter.Name())
	}
}

func TestExternalHTTPURLValidationRejectsUnsafeEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	goodURL := server.URL + externalHTTPPath
	goodHost := mustURLHost(t, goodURL)

	tests := []struct {
		name string
		raw  string
		host string
	}{
		{name: "missing allowlist", raw: goodURL, host: ""},
		{name: "non allowlisted host", raw: goodURL, host: "other.example"},
		{name: "userinfo", raw: "https://user:pass@" + goodHost + externalHTTPPath, host: goodHost},
		{name: "query", raw: goodURL + "?token=secret", host: goodHost},
		{name: "fragment", raw: goodURL + "#secret", host: goodHost},
		{name: "wrong path", raw: server.URL + "/predict", host: goodHost},
		{name: "unsafe http non loopback", raw: "http://predictor.example" + externalHTTPPath, host: "predictor.example"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewExternalHTTPAdapter(ExternalHTTPConfig{
				URL:              tt.raw,
				AllowedHosts:     splitCSV(tt.host),
				Timeout:          time.Second,
				MaxRequestBytes:  1024,
				MaxResponseBytes: 1024,
			})
			if err == nil {
				t.Fatalf("adapter config succeeded for unsafe endpoint")
			}
		})
	}
}

func TestExternalHTTPTokenEnvValidationNeverSurfacesTokenValue(t *testing.T) {
	_, err := AdapterFactoryConfigFromEnv(testEnv(map[string]string{
		"TRIP_UPDATES_ADAPTER":                     "external_http",
		"TRIP_UPDATES_EXTERNAL_HTTP_URL":           "https://predictor.example" + externalHTTPPath,
		"TRIP_UPDATES_EXTERNAL_HTTP_ALLOWED_HOSTS": "predictor.example",
		"TRIP_UPDATES_EXTERNAL_HTTP_TOKEN_ENV":     "bad-name",
	}))
	if err == nil {
		t.Fatalf("invalid token env name succeeded")
	}

	secret := "super-secret-token-value"
	_, err = AdapterFactoryConfigFromEnv(testEnv(map[string]string{
		"TRIP_UPDATES_ADAPTER":                     "external_http",
		"TRIP_UPDATES_EXTERNAL_HTTP_URL":           "https://predictor.example" + externalHTTPPath,
		"TRIP_UPDATES_EXTERNAL_HTTP_ALLOWED_HOSTS": "predictor.example",
		"TRIP_UPDATES_EXTERNAL_HTTP_TOKEN_ENV":     "PREDICTOR_TOKEN",
		"PREDICTOR_TOKEN":                          secret,
	}))
	if err != nil && strings.Contains(err.Error(), secret) {
		t.Fatalf("error exposes token value: %v", err)
	}

	_, err = AdapterFactoryConfigFromEnv(testEnv(map[string]string{
		"TRIP_UPDATES_ADAPTER":                     "external_http",
		"TRIP_UPDATES_EXTERNAL_HTTP_URL":           "https://predictor.example" + externalHTTPPath,
		"TRIP_UPDATES_EXTERNAL_HTTP_ALLOWED_HOSTS": "predictor.example",
		"TRIP_UPDATES_EXTERNAL_HTTP_TOKEN_ENV":     "PREDICTOR_TOKEN",
	}))
	if err == nil {
		t.Fatalf("missing referenced token succeeded")
	}
	if strings.Contains(err.Error(), "PREDICTOR_TOKEN") {
		t.Fatalf("missing token error exposes token env name: %v", err)
	}
}

func TestExternalHTTPAdapterSendsSanitizedDTOAndMapsResponse(t *testing.T) {
	generatedAt := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	arrival := generatedAt.Add(2 * time.Minute)
	confidence := 0.91
	var requestBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != externalHTTPPath {
			t.Fatalf("path = %q, want %s", r.URL.Path, externalHTTPPath)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer runtime-token" {
			t.Fatalf("authorization header = %q, want bearer token", got)
		}
		var err error
		requestBody, err = ioReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request: %v", err)
		}
		writeExternalHTTPResponse(t, w, externalHTTPResponseDTO{
			SchemaVersion: "external_trip_updates.v1",
			AgencyID:      "demo-agency",
			FeedVersionID: "feed-demo",
			GeneratedAt:   generatedAt,
			TripUpdates: []externalTripUpdateDTO{{
				AgencyID:      "demo-agency",
				FeedVersionID: "feed-demo",
				EntityID:      "external-1",
				VehicleID:     "bus-10",
				TripID:        "trip-10",
				RouteID:       "route-10",
				StartDate:     "20260509",
				StartTime:     "08:00:00",
				Confidence:    &confidence,
				StopTimeUpdates: []externalStopUpdateDTO{{
					StopID:       "stop-2",
					StopSequence: 2,
					ArrivalTime:  &arrival,
				}},
			}},
		})
	}))
	defer server.Close()

	adapter, err := NewExternalHTTPAdapter(ExternalHTTPConfig{
		URL:              server.URL + externalHTTPPath,
		AllowedHosts:     []string{mustURLHost(t, server.URL)},
		Timeout:          time.Second,
		MaxRequestBytes:  4096,
		MaxResponseBytes: 4096,
		TokenEnv:         "PREDICTOR_TOKEN",
		TokenValue:       "runtime-token",
	})
	if err != nil {
		t.Fatalf("new adapter: %v", err)
	}
	result, err := adapter.PredictTripUpdates(context.Background(), Request{
		AgencyID:            "demo-agency",
		GeneratedAt:         generatedAt,
		ActiveFeedVersion:   gtfs.FeedVersion{ID: "feed-demo", AgencyID: "demo-agency"},
		VehiclePositionsURL: "https://feeds.example/public/gtfsrt/vehicle_positions.pb",
		Telemetry: []telemetry.StoredEvent{{
			ID: 99,
			Event: telemetry.Event{
				AgencyID:  "demo-agency",
				DeviceID:  "device-secret",
				VehicleID: "bus-10",
				DriverID:  "driver-secret",
				Timestamp: generatedAt,
				Lat:       34.05,
				Lon:       -118.25,
				TripHint:  "trip-10",
			},
			PayloadJSON: json.RawMessage(`{"vendor_token":"secret","device_id":"nested-device"}`),
		}},
		Assignments: map[string]state.Assignment{"bus-10": {
			ID:                  7,
			AgencyID:            "demo-agency",
			VehicleID:           "bus-10",
			FeedVersionID:       "feed-demo",
			TelemetryEventID:    99,
			State:               state.StateInService,
			ServiceDate:         "20260509",
			RouteID:             "route-10",
			TripID:              "trip-10",
			BlockID:             "block-10",
			StartDate:           "20260509",
			StartTime:           "08:00:00",
			CurrentStopSequence: 1,
			ShapeDistTraveled:   12.5,
			Confidence:          0.88,
			AssignmentSource:    state.AssignmentSourceManualOverride,
			ReasonCodes:         []string{state.ReasonManualOverrideActive},
			DegradedState:       state.DegradedNone,
			ScoreDetails:        map[string]any{"raw_override_reason": "private"},
			ManualOverrideID:    123,
		}},
	})
	if err != nil {
		t.Fatalf("predict: %v", err)
	}
	if result.Diagnostics.Status != StatusOK || len(result.TripUpdates) != 1 {
		t.Fatalf("result = %+v, want one accepted external update", result)
	}
	body := string(requestBody)
	for _, forbidden := range []string{
		"device_id", "driver_id", "payload_json", "vendor_token", "token", "Authorization",
		"score_details", "manual_override_id", "raw_override_reason", "telemetry_event_id",
	} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("external request includes forbidden field/value %q: %s", forbidden, body)
		}
	}
	for _, allowed := range []string{"vehicle_id", "timestamp", "lat", "lon", "trip_hint", "manual_override_active"} {
		if !strings.Contains(body, allowed) {
			t.Fatalf("external request missing allowed field %q: %s", allowed, body)
		}
	}
}

func TestExternalHTTPFailuresReturnAdapterErrorDiagnostics(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		maxResp int64
		want    string
	}{
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "private failure body", http.StatusInternalServerError)
			},
			maxResp: 1024,
			want:    "http_status",
		},
		{
			name: "malformed json",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"schema_version":`))
			},
			maxResp: 1024,
			want:    "malformed_json",
		},
		{
			name: "oversized response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(strings.Repeat("x", 64)))
			},
			maxResp: 8,
			want:    "response_too_large",
		},
		{
			name: "missing start time",
			handler: func(w http.ResponseWriter, r *http.Request) {
				writeExternalHTTPResponse(t, w, externalHTTPResponseDTO{
					SchemaVersion: "external_trip_updates.v1",
					AgencyID:      "demo-agency",
					FeedVersionID: "feed-demo",
					GeneratedAt:   time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
					TripUpdates: []externalTripUpdateDTO{{
						AgencyID:      "demo-agency",
						FeedVersionID: "feed-demo",
						TripID:        "trip-10",
						StartDate:     "20260509",
						Confidence:    floatPtr(0.9),
					}},
				})
			},
			maxResp: 1024,
			want:    "missing_start_time",
		},
		{
			name: "wrong scope",
			handler: func(w http.ResponseWriter, r *http.Request) {
				writeExternalHTTPResponse(t, w, externalHTTPResponseDTO{
					SchemaVersion: "external_trip_updates.v1",
					AgencyID:      "other-agency",
					FeedVersionID: "feed-demo",
					GeneratedAt:   time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
				})
			},
			maxResp: 1024,
			want:    "wrong_scope",
		},
		{
			name: "stale response timestamp",
			handler: func(w http.ResponseWriter, r *http.Request) {
				writeExternalHTTPResponse(t, w, externalHTTPResponseDTO{
					SchemaVersion: "external_trip_updates.v1",
					AgencyID:      "demo-agency",
					FeedVersionID: "feed-demo",
					GeneratedAt:   time.Date(2026, 5, 9, 11, 59, 59, 0, time.UTC),
				})
			},
			maxResp: 1024,
			want:    "stale_response_timestamp",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			adapter, err := NewExternalHTTPAdapter(ExternalHTTPConfig{
				URL:              server.URL + externalHTTPPath,
				AllowedHosts:     []string{mustURLHost(t, server.URL)},
				Timeout:          time.Second,
				MaxRequestBytes:  1024,
				MaxResponseBytes: tt.maxResp,
			})
			if err != nil {
				t.Fatalf("new adapter: %v", err)
			}
			result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC), nil, nil))
			if err != nil {
				t.Fatalf("predict returned hard error: %v", err)
			}
			if result.Diagnostics.Status != StatusError || result.Diagnostics.Reason != ReasonAdapterError {
				t.Fatalf("diagnostics = %+v, want adapter_error", result.Diagnostics)
			}
			if result.Diagnostics.Details["failure_type"] != tt.want {
				t.Fatalf("details = %+v, want failure_type %s", result.Diagnostics.Details, tt.want)
			}
			if strings.Contains(fmtAny(result.Diagnostics.Details), "private failure body") {
				t.Fatalf("diagnostics include raw response body: %+v", result.Diagnostics.Details)
			}
		})
	}
}

func TestExternalHTTPRequestTooLargeReturnsAdapterErrorWithoutHTTPCall(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		http.Error(w, "should not be called", http.StatusInternalServerError)
	}))
	defer server.Close()
	adapter, err := NewExternalHTTPAdapter(ExternalHTTPConfig{
		URL:              server.URL + externalHTTPPath,
		AllowedHosts:     []string{mustURLHost(t, server.URL)},
		Timeout:          time.Second,
		MaxRequestBytes:  16,
		MaxResponseBytes: 1024,
	})
	if err != nil {
		t.Fatalf("new adapter: %v", err)
	}

	result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC), []telemetry.StoredEvent{
		predictionTelemetry(1, "bus-10", time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)),
	}, nil))
	if err != nil {
		t.Fatalf("predict returned hard error: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 0 {
		t.Fatalf("http calls = %d, want no call when request exceeds cap", got)
	}
	assertAdapterFailure(t, result, "request_too_large")
}

func TestExternalHTTPTimeoutReturnsRedactedAdapterError(t *testing.T) {
	secret := "super-secret-token-value"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		_, _ = w.Write([]byte(`{"private":"raw body should not appear"}`))
	}))
	defer server.Close()
	adapter, err := NewExternalHTTPAdapter(ExternalHTTPConfig{
		URL:              server.URL + externalHTTPPath,
		AllowedHosts:     []string{mustURLHost(t, server.URL)},
		Timeout:          time.Nanosecond,
		MaxRequestBytes:  1024,
		MaxResponseBytes: 1024,
		TokenEnv:         "PREDICTOR_TOKEN",
		TokenValue:       secret,
	})
	if err != nil {
		t.Fatalf("new adapter: %v", err)
	}

	result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC), nil, nil))
	if err != nil {
		t.Fatalf("predict returned hard error: %v", err)
	}
	assertAdapterFailure(t, result, "timeout")
	rendered := fmtAny(result.Diagnostics.Details)
	for _, forbidden := range []string{secret, "PREDICTOR_TOKEN", "raw body", "private", mustURLHost(t, server.URL), externalHTTPPath} {
		if strings.Contains(rendered, forbidden) {
			t.Fatalf("timeout diagnostics include forbidden detail %q: %s", forbidden, rendered)
		}
	}
}

func TestExternalHTTPRedirectIsNotFollowedAndDegradesSafely(t *testing.T) {
	var targetCalls int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&targetCalls, 1)
		writeExternalHTTPResponse(t, w, externalHTTPResponseDTO{
			SchemaVersion: "external_trip_updates.v1",
			AgencyID:      "demo-agency",
			FeedVersionID: "feed-demo",
			GeneratedAt:   time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
		})
	}))
	defer target.Close()
	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL+externalHTTPPath, http.StatusTemporaryRedirect)
	}))
	defer redirector.Close()
	adapter, err := NewExternalHTTPAdapter(ExternalHTTPConfig{
		URL:              redirector.URL + externalHTTPPath,
		AllowedHosts:     []string{mustURLHost(t, redirector.URL)},
		Timeout:          time.Second,
		MaxRequestBytes:  1024,
		MaxResponseBytes: 1024,
	})
	if err != nil {
		t.Fatalf("new adapter: %v", err)
	}

	result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC), nil, nil))
	if err != nil {
		t.Fatalf("predict returned hard error: %v", err)
	}
	if got := atomic.LoadInt32(&targetCalls); got != 0 {
		t.Fatalf("redirect target calls = %d, want redirect not followed", got)
	}
	assertAdapterFailure(t, result, "http_status")
	if result.Diagnostics.Details["http_status_code"] != http.StatusTemporaryRedirect {
		t.Fatalf("details = %+v, want redirect status code", result.Diagnostics.Details)
	}
}

func floatPtr(value float64) *float64 {
	return &value
}

func TestExternalHTTPShadowReturnsDeterministicOutputWithBoundedDiagnostics(t *testing.T) {
	deterministic := &stubPredictionAdapter{name: "deterministic", result: Result{
		Diagnostics: Diagnostics{Status: StatusOK, Reason: ReasonPredictionsAvailable},
		TripUpdates: []TripUpdate{{EntityID: "deterministic", TripID: "trip-10", StartDate: "20260509"}},
	}}
	external := &stubPredictionAdapter{name: "external_http", err: errors.New("token=secret host=predictor.example")}
	adapter, err := NewExternalHTTPShadowAdapter(deterministic, external)
	if err != nil {
		t.Fatalf("new shadow adapter: %v", err)
	}
	result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC), nil, nil))
	if err != nil {
		t.Fatalf("shadow predict: %v", err)
	}
	if len(result.TripUpdates) != 1 || result.TripUpdates[0].EntityID != "deterministic" {
		t.Fatalf("trip updates = %+v, want deterministic public output", result.TripUpdates)
	}
	rendered := fmtAny(result.Diagnostics.Details)
	if !strings.Contains(rendered, "external_http_shadow") || strings.Contains(rendered, "secret") || strings.Contains(rendered, "predictor.example") {
		t.Fatalf("shadow diagnostics not bounded/redacted: %+v", result.Diagnostics.Details)
	}
	shadow := shadowDetailsFromResult(t, result)
	if shadow["fallback_used"] != true || shadow["shadow_only"] != true || shadow["external_output_public"] != false {
		t.Fatalf("shadow fallback/public flags = %+v, want deterministic fallback and non-public external output", shadow)
	}
	if shadow["divergence_status"] != "external_error" {
		t.Fatalf("shadow divergence = %+v, want external_error", shadow)
	}
}

func TestExternalHTTPShadowSummarizesDivergenceAndConfidenceWithoutRawDetails(t *testing.T) {
	lowConfidence := 0.25
	highConfidence := 0.92
	generatedAt := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	deterministic := &stubPredictionAdapter{name: "deterministic", result: Result{
		Diagnostics: Diagnostics{Status: StatusOK, Reason: ReasonPredictionsAvailable},
		TripUpdates: []TripUpdate{{
			AgencyID:      "demo-agency",
			FeedVersionID: "feed-demo",
			VehicleID:     "bus-10",
			TripID:        "trip-10",
			StartDate:     "20260509",
			StartTime:     "08:00:00",
			Confidence:    &highConfidence,
		}},
	}}
	external := &stubPredictionAdapter{name: "external_http", result: Result{
		Diagnostics: Diagnostics{
			Status:  StatusOK,
			Reason:  ReasonNoEligiblePredictions,
			Details: map[string]any{"raw_response": "token=secret host=predictor.example"},
		},
		TripUpdates: []TripUpdate{
			{
				AgencyID:      "demo-agency",
				FeedVersionID: "feed-demo",
				VehicleID:     "bus-10",
				TripID:        "trip-10",
				StartDate:     "20260509",
				StartTime:     "08:00:00",
				Confidence:    &lowConfidence,
			},
			{
				AgencyID:      "demo-agency",
				FeedVersionID: "feed-demo",
				VehicleID:     "bus-11",
				TripID:        "trip-11",
				StartDate:     "20260509",
				StartTime:     "09:00:00",
				Confidence:    &highConfidence,
			},
		},
	}}
	adapter, err := NewExternalHTTPShadowAdapter(deterministic, external)
	if err != nil {
		t.Fatalf("new shadow adapter: %v", err)
	}
	result, err := adapter.PredictTripUpdates(context.Background(), predictionRequest(generatedAt, nil, nil))
	if err != nil {
		t.Fatalf("shadow predict: %v", err)
	}
	if len(result.TripUpdates) != 1 || result.TripUpdates[0].TripID != "trip-10" {
		t.Fatalf("trip updates = %+v, want deterministic output only", result.TripUpdates)
	}
	shadow := shadowDetailsFromResult(t, result)
	if shadow["divergence_status"] != "trip_identity_delta" || shadow["matching_identity_count"] != 1 || shadow["external_only_count"] != 1 || shadow["count_delta"] != 1 {
		t.Fatalf("shadow identity summary = %+v, want bounded divergence counts", shadow)
	}
	if shadow["external_low_confidence_count"] != 1 || shadow["external_missing_confidence_count"] != 0 {
		t.Fatalf("shadow confidence summary = %+v, want low-confidence external count", shadow)
	}
	reasons, ok := shadow["divergence_by_reason"].(map[string]int)
	if !ok {
		t.Fatalf("divergence_by_reason = %T %+v, want map[string]int", shadow["divergence_by_reason"], shadow["divergence_by_reason"])
	}
	if reasons["trip_identity_delta"] != 1 || reasons["external_low_confidence"] != 1 || reasons["reason_mismatch"] != 1 {
		t.Fatalf("divergence_by_reason = %+v, want identity/confidence/reason review", reasons)
	}
	rendered := fmtAny(result.Diagnostics.Details)
	if strings.Contains(rendered, "secret") || strings.Contains(rendered, "predictor.example") || strings.Contains(rendered, "raw_response") {
		t.Fatalf("shadow diagnostics include raw external details: %s", rendered)
	}
}

func testEnv(values map[string]string) EnvLookup {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}

func mustURLHost(t *testing.T, raw string) string {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	return parsed.Host
}

func writeExternalHTTPResponse(t *testing.T, w http.ResponseWriter, response externalHTTPResponseDTO) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}

func ioReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func fmtAny(value any) string {
	payload, _ := json.Marshal(value)
	return string(payload)
}

func assertAdapterFailure(t *testing.T, result Result, failureType string) {
	t.Helper()
	if result.Diagnostics.Status != StatusError || result.Diagnostics.Reason != ReasonAdapterError {
		t.Fatalf("diagnostics = %+v, want adapter_error", result.Diagnostics)
	}
	if len(result.TripUpdates) != 0 {
		t.Fatalf("trip updates = %+v, want no output on adapter failure", result.TripUpdates)
	}
	if result.Diagnostics.Details["failure_type"] != failureType {
		t.Fatalf("details = %+v, want failure_type %s", result.Diagnostics.Details, failureType)
	}
}

func shadowDetailsFromResult(t *testing.T, result Result) map[string]any {
	t.Helper()
	shadow, ok := result.Diagnostics.Details["external_http_shadow"].(map[string]any)
	if !ok {
		t.Fatalf("external_http_shadow = %T %+v, want map", result.Diagnostics.Details["external_http_shadow"], result.Diagnostics.Details)
	}
	return shadow
}

type stubPredictionAdapter struct {
	name   string
	result Result
	err    error
}

func (s *stubPredictionAdapter) Name() string {
	return s.name
}

func (s *stubPredictionAdapter) PredictTripUpdates(context.Context, Request) (Result, error) {
	return s.result, s.err
}
