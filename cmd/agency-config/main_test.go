package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
	"open-transit-rt/internal/devices"
	"open-transit-rt/internal/feed/schedule"
	"open-transit-rt/internal/prediction"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

func TestValidationRunDerivesRealtimeArtifacts(t *testing.T) {
	validatorPath := writeRealtimeValidator(t)
	t.Setenv("GTFS_RT_VALIDATOR_PATH", validatorPath)
	t.Setenv("GTFS_RT_VALIDATOR_VERSION", "test-validator")
	t.Setenv("GTFS_RT_VALIDATOR_ARGS", "")

	for _, feedType := range []string{"vehicle_positions", "trip_updates", "alerts"} {
		t.Run(feedType, func(t *testing.T) {
			store := &fakePublicationStore{}
			artifacts := &fakeRealtimeArtifacts{payloads: map[string][]byte{
				feedType: []byte("protobuf-" + feedType),
			}}
			handler := newHandlerWithRealtime(
				"demo-agency",
				fakeScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "demo-agency", FeedVersionID: "feed-demo", RevisionTime: time.Now().UTC(), Payload: []byte("schedule zip bytes")}},
				store,
				fakeDeviceStore{},
				fakePinger{},
				auth.TestAuthenticator{Principal: auth.Principal{Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer}},
				artifacts,
			)

			body := []byte(fmt.Sprintf(`{"validator_id":"realtime-mobilitydata","feed_type":%q}`, feedType))
			req := httptest.NewRequest(http.MethodPost, "/admin/validation/run", bytes.NewReader(body))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
			}
			var result compliance.ValidationResult
			if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
				t.Fatalf("decode result: %v", err)
			}
			if result.Status != "passed" || result.FeedType != feedType || store.result.Status != "passed" {
				t.Fatalf("result = %+v stored = %+v, want passed %s validation", result, store.result, feedType)
			}
			if artifacts.calls[feedType] != 1 {
				t.Fatalf("artifact calls = %+v, want one call for %s", artifacts.calls, feedType)
			}
		})
	}
}

func TestValidationRunDerivesScheduleArtifact(t *testing.T) {
	validatorPath := writeScheduleValidator(t)
	t.Setenv("GTFS_VALIDATOR_PATH", validatorPath)
	t.Setenv("GTFS_VALIDATOR_VERSION", "test-validator")

	store := &fakePublicationStore{}
	handler := newHandlerWithRealtime(
		"demo-agency",
		fakeScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "demo-agency", FeedVersionID: "feed-demo", RevisionTime: time.Now().UTC(), Payload: []byte("schedule zip bytes")}},
		store,
		fakeDeviceStore{},
		fakePinger{},
		auth.TestAuthenticator{Principal: auth.Principal{Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer}},
		&fakeRealtimeArtifacts{},
	)

	req := httptest.NewRequest(http.MethodPost, "/admin/validation/run", bytes.NewReader([]byte(`{"validator_id":"static-mobilitydata","feed_type":"schedule"}`)))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var result compliance.ValidationResult
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if result.Status != "passed" || result.FeedType != "schedule" || result.FeedVersionID != "feed-demo" {
		t.Fatalf("result = %+v, want passed schedule validation recorded for feed-demo", result)
	}
	if store.result.Status != "passed" || store.result.FeedType != "schedule" || store.result.FeedVersionID != "feed-demo" {
		t.Fatalf("stored result = %+v, want persisted passed schedule validation", store.result)
	}
}

func TestValidationRunRejectsClientSuppliedRealtimePath(t *testing.T) {
	handler := newHandlerWithRealtime(
		"demo-agency",
		fakeScheduleBuilder{},
		&fakePublicationStore{},
		fakeDeviceStore{},
		fakePinger{},
		auth.TestAuthenticator{Principal: auth.Principal{Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer}},
		&fakeRealtimeArtifacts{},
	)
	req := httptest.NewRequest(http.MethodPost, "/admin/validation/run", bytes.NewReader([]byte(`{"validator_id":"realtime-mobilitydata","feed_type":"alerts","realtime_pb_path":"/tmp/evil.pb"}`)))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}

func TestValidationRunRejectsUnknownFields(t *testing.T) {
	handler := newHandlerWithRealtime(
		"demo-agency",
		fakeScheduleBuilder{},
		&fakePublicationStore{},
		fakeDeviceStore{},
		fakePinger{},
		auth.TestAuthenticator{Principal: auth.Principal{Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer}},
		&fakeRealtimeArtifacts{},
	)
	for _, body := range []string{
		`{"validator_id":"static-mobilitydata","feed_type":"schedule","agency_id":"demo-agency"}`,
		`{"validator_id":"realtime-mobilitydata","feed_type":"alerts","argv":["bad"]}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/admin/validation/run", bytes.NewReader([]byte(body)))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("body %s status = %d, want 400", body, rr.Code)
		}
	}
}

func TestValidationRunRejectsUnknownValidatorAndFeedType(t *testing.T) {
	handler := newHandlerWithRealtime(
		"demo-agency",
		fakeScheduleBuilder{},
		&fakePublicationStore{},
		fakeDeviceStore{},
		fakePinger{},
		auth.TestAuthenticator{Principal: auth.Principal{Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer}},
		&fakeRealtimeArtifacts{},
	)
	for _, body := range []string{
		`{"validator_id":"missing","feed_type":"schedule"}`,
		`{"validator_id":"static-mobilitydata","feed_type":"alerts"}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/admin/validation/run", bytes.NewReader([]byte(body)))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("body %s status = %d, want 400", body, rr.Code)
		}
	}
}

func TestAgencyConfigAdminRejectsUnauthenticatedAccess(t *testing.T) {
	handler := newHandlerWithRealtime(
		"demo-agency",
		fakeScheduleBuilder{},
		&fakePublicationStore{},
		fakeDeviceStore{},
		fakePinger{},
		authRejectAll{},
		&fakeRealtimeArtifacts{},
	)
	req := httptest.NewRequest(http.MethodPost, "/admin/validation/run", bytes.NewReader([]byte(`{"validator_id":"static-mobilitydata","feed_type":"schedule"}`)))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rr.Code)
	}
}

func TestPublicScheduleRemainsAnonymous(t *testing.T) {
	handler := newHandlerWithRealtime(
		"demo-agency",
		fakeScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "demo-agency", FeedVersionID: "feed-demo", RevisionTime: time.Now().UTC(), Payload: []byte("schedule zip bytes")}},
		&fakePublicationStore{},
		fakeDeviceStore{},
		fakePinger{},
		authRejectAll{},
		&fakeRealtimeArtifacts{},
	)
	req := httptest.NewRequest(http.MethodGet, "/public/gtfs/schedule.zip", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want anonymous 200", rr.Code)
	}
}

func TestPublicFeedsJSONIsQueryRoutedAndPublicMetadataOnly(t *testing.T) {
	now := time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{discoveries: map[string]compliance.FeedDiscovery{
		"agency-a": {
			AgencyID: "agency-a", AgencyName: "Agency A", GeneratedAt: now, PublicationEnvironment: compliance.EnvironmentDev,
			PublicBaseURL: "https://agency-a.example",
			Feeds:         []compliance.FeedMetadata{{FeedType: "schedule", CanonicalPublicURL: "https://agency-a.example/public/gtfs/schedule.zip", ActivationStatus: "active", ActiveFeedVersionID: "feed-a"}},
		},
		"agency-b": {
			AgencyID: "agency-b", AgencyName: "Agency B", GeneratedAt: now, PublicationEnvironment: compliance.EnvironmentDev,
			PublicBaseURL: "https://agency-b.example",
			Feeds:         []compliance.FeedMetadata{{FeedType: "schedule", CanonicalPublicURL: "https://agency-b.example/public/gtfs/schedule.zip", ActivationStatus: "active", ActiveFeedVersionID: "feed-b"}},
		},
	}}
	handler := newHandlerWithRealtime(
		"agency-a",
		fakeScheduleBuilder{},
		store,
		fakeDeviceStore{},
		fakePinger{},
		authRejectAll{},
		&fakeRealtimeArtifacts{},
	)

	req := httptest.NewRequest(http.MethodGet, "/public/feeds.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("default status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.discoveryAgencyID != "agency-a" || strings.Contains(rr.Body.String(), "agency-b.example") {
		t.Fatalf("default feed discovery agency=%q body=%s, want configured agency only", store.discoveryAgencyID, rr.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/public/feeds.json?agency_id=agency-b", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("query status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if store.discoveryAgencyID != "agency-b" || strings.Contains(body, "agency-a.example") || !strings.Contains(body, "agency-b.example") {
		t.Fatalf("query feed discovery agency=%q body=%s, want requested agency only", store.discoveryAgencyID, body)
	}
	for _, forbidden := range []string{"token", "token_hash", "private_notes", "raw_report", "operator_artifact", "evidence/private", "payload_json"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("public feeds.json exposes private/admin field %q: %s", forbidden, body)
		}
	}
}

func TestPublicAgencyFeedDiscoveryAndScheduleArePathRouted(t *testing.T) {
	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{discoveries: map[string]compliance.FeedDiscovery{
		"agency-a": {AgencyID: "agency-a", AgencyName: "Agency A", GeneratedAt: now, PublicBaseURL: "https://feeds.example/public/agencies/agency-a"},
		"agency-b": {AgencyID: "agency-b", AgencyName: "Agency B", GeneratedAt: now, PublicBaseURL: "https://feeds.example/public/agencies/agency-b"},
	}}
	builder := fakeScheduleBuilder{snapshotsByAgency: map[string]schedule.Snapshot{
		"agency-a": {AgencyID: "agency-a", FeedVersionID: "feed-a", RevisionTime: now, Payload: []byte("agency-a schedule")},
		"agency-b": {AgencyID: "agency-b", FeedVersionID: "feed-b", RevisionTime: now, Payload: []byte("agency-b schedule")},
	}}
	handler := newHandlerWithRealtime(
		"agency-a",
		builder,
		store,
		fakeDeviceStore{},
		fakePinger{},
		authRejectAll{},
		&fakeRealtimeArtifacts{},
	)

	req := httptest.NewRequest(http.MethodGet, "/public/agencies/agency-b/feeds.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("feeds status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.discoveryAgencyID != "agency-b" || strings.Contains(rr.Body.String(), "agency-a") || !strings.Contains(rr.Body.String(), "agency-b") {
		t.Fatalf("path feed discovery agency=%q body=%s, want agency-b only", store.discoveryAgencyID, rr.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/public/agencies/agency-b/gtfs/schedule.zip", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("schedule status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Body.String(); got != "agency-b schedule" {
		t.Fatalf("schedule payload = %q, want agency-b schedule", got)
	}

	req = httptest.NewRequest(http.MethodGet, "/public/agencies/.hidden/feeds.json", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("invalid agency status = %d, want 400", rr.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/public/agencies/agency-b/gtfs/schedule.json", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public debug-like route status = %d, want 404", rr.Code)
	}
}

func TestOperationsConsoleRejectsUnauthenticatedAccess(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, authRejectAll{})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rr.Code)
	}
}

func TestAgencyConfigAdminAgencyBoundaries(t *testing.T) {
	store := &fakePublicationStore{
		scorecard: compliance.Scorecard{AgencyID: "agency-a", OverallStatus: compliance.StatusYellow},
		consumers: []compliance.ConsumerRecord{{ConsumerName: "Maps A", Status: "not_started", UpdatedAt: time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)}},
	}
	handler := newHandlerWithRealtime(
		"agency-a",
		fakeScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "agency-a", FeedVersionID: "feed-a", RevisionTime: time.Now().UTC(), Payload: []byte("schedule zip bytes")}},
		store,
		fakeDeviceStore{},
		fakePinger{},
		auth.TestAuthenticator{Principal: auth.Principal{Subject: "admin-a@example.com", AgencyID: "agency-a", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer}},
		&fakeRealtimeArtifacts{},
	)

	req := httptest.NewRequest(http.MethodPost, "/admin/publication/bootstrap", bytes.NewReader([]byte(`{"agency_id":"agency-b","public_base_url":"https://agency-b.example"}`)))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("bootstrap conflict status = %d, want 403", rr.Code)
	}
	if store.bootstrapInput.AgencyID != "" {
		t.Fatalf("bootstrap ran despite conflict: %+v", store.bootstrapInput)
	}

	req = httptest.NewRequest(http.MethodPost, "/admin/publication/bootstrap", bytes.NewReader([]byte(`{"public_base_url":"https://agency-a.example","feed_base_url":"https://agency-a.example/public"}`)))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("bootstrap status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.bootstrapInput.AgencyID != "agency-a" || store.bootstrapInput.ActorID != "admin-a@example.com" {
		t.Fatalf("bootstrap identity = %+v, want principal agency/actor", store.bootstrapInput)
	}

	for _, tc := range []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "scorecard get", method: http.MethodGet, path: "/admin/compliance/scorecard?agency_id=agency-b"},
		{name: "scorecard post", method: http.MethodPost, path: "/admin/compliance/scorecard", body: `{"agency_id":"agency-b"}`},
		{name: "consumer get", method: http.MethodGet, path: "/admin/consumer-ingestion?agency_id=agency-b"},
		{name: "consumer post", method: http.MethodPost, path: "/admin/consumer-ingestion", body: `{"agency_id":"agency-b","consumer_name":"Maps B"}`},
		{name: "device rebind", method: http.MethodPost, path: "/admin/devices/rebind", body: `{"agency_id":"agency-b","device_id":"device-b-1","vehicle_id":"bus-b-1"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusForbidden {
				t.Fatalf("status = %d, want 403: %s", rr.Code, rr.Body.String())
			}
		})
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/compliance/scorecard", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || store.latestScorecardAgencyID != "agency-a" {
		t.Fatalf("scorecard status=%d agency=%q, want agency-a", rr.Code, store.latestScorecardAgencyID)
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/consumer-ingestion", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || store.listConsumersAgencyID != "agency-a" || strings.Contains(rr.Body.String(), "Maps B") {
		t.Fatalf("consumer status=%d agency=%q body=%s, want agency-a only", rr.Code, store.listConsumersAgencyID, rr.Body.String())
	}
}

func TestOperationsConsoleRendersEmptyState(t *testing.T) {
	store := &fakePublicationStore{
		discoveryErr: errors.New("no feed config"),
		scorecardErr: errors.New("no scorecard"),
	}
	srv := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Operations Console", "publication metadata is not configured yet", "telemetry repository is not available", "no Trip Updates diagnostics recorded yet"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
}

func TestOperationsSetupRendersTruthfulMissingStates(t *testing.T) {
	store := &fakePublicationStore{
		discoveryErr:         errors.New("no feed config"),
		scorecardErr:         errors.New("no scorecard"),
		publicationConfigErr: errors.New("no publication config"),
	}
	srv := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}, telemetry: fakeTelemetryRepository{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/setup", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Guided Setup Checklist",
		"publication metadata",
		"validation records",
		"device bindings",
		"telemetry repository",
		"docs/evidence tracker",
		"not observed yet",
		"prepared is not submitted or accepted",
		"Browser ZIP upload is deferred",
		"Validation is supporting evidence only",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"accepted by", "CAL-ITP/Caltrans compliant", "consumer ingestion confirmed"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("body overclaims %q: %s", forbidden, body)
		}
	}
}

func TestOperationsReadinessWorkflowRendersEvidenceBoundedRows(t *testing.T) {
	now := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	coverage := 75.0
	store := &fakePublicationStore{
		discovery: compliance.FeedDiscovery{
			AgencyID: "demo-agency", AgencyName: "Demo Agency", GeneratedAt: now, PublicationEnvironment: "dev",
			PublicBaseURL:         "https://feeds.example.org",
			TechnicalContactEmail: "ops@example.org",
			License:               compliance.License{Name: "CC BY 4.0", URL: "https://example.org/license"},
			Feeds: []compliance.FeedMetadata{
				{FeedType: "schedule", CanonicalPublicURL: "https://feeds.example.org/public/gtfs/schedule.zip", ActivationStatus: "active", ActiveFeedVersionID: "feed-v1", LastValidationStatus: "passed", LastValidationAt: &now, LastHealthStatus: "ok", LastHealthAt: &now},
				{FeedType: "vehicle_positions", CanonicalPublicURL: "https://feeds.example.org/public/gtfsrt/vehicle_positions.pb", ActivationStatus: "active", ActiveFeedVersionID: "feed-v1", LastValidationStatus: "passed", LastValidationAt: &now, LastHealthStatus: "ok", LastHealthAt: &now},
				{FeedType: "trip_updates", CanonicalPublicURL: "https://feeds.example.org/public/gtfsrt/trip_updates.pb", ActivationStatus: "active", ActiveFeedVersionID: "feed-v1", LastValidationStatus: "warning", LastValidationAt: &now, LastHealthStatus: "ok", LastHealthAt: &now},
				{FeedType: "alerts", CanonicalPublicURL: "https://feeds.example.org/public/gtfsrt/alerts.pb", ActivationStatus: "active", ActiveFeedVersionID: "feed-v1", LastValidationStatus: "passed", LastValidationAt: &now, LastHealthStatus: "ok", LastHealthAt: &now},
			},
			Readiness: compliance.Readiness{
				Discoverable: true, HTTPSURLs: true, LicenseComplete: true, ContactComplete: true,
				AllRequiredFeedsListed: true, CanonicalValidationComplete: true,
			},
		},
		scorecard: compliance.Scorecard{AgencyID: "demo-agency", SnapshotAt: now, OverallStatus: compliance.StatusYellow},
		consumers: []compliance.ConsumerRecord{{ConsumerName: "Google Maps", Status: "not_started", UpdatedAt: now}},
		tripDiagnostics: compliance.TripUpdatesDiagnosticsSummary{
			Recorded: true, SnapshotAt: now, AdapterName: "deterministic", DiagnosticsStatus: prediction.StatusOK,
			DiagnosticsReason: prediction.ReasonPartialPredictions,
			Metrics: prediction.Metrics{
				TripUpdatesCoverageRate: prediction.RateMetric{Numerator: 3, Denominator: 4, Percent: &coverage, Status: "measured"},
				FutureStopCoverageRate:  prediction.RateMetric{Numerator: 3, Denominator: 4, Percent: &coverage, Status: "measured"},
			},
		},
	}
	srv := newOperationsTestHandler(&handler{
		store: store,
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{{
			AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now, CreatedAt: now,
		}}},
		telemetry: fakeTelemetryRepository{latest: []telemetry.StoredEvent{{
			Event:      telemetry.Event{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Timestamp: now.Add(-30 * time.Second), Lat: 1, Lon: 2},
			ReceivedAt: now.Add(-29 * time.Second), IngestStatus: telemetry.IngestStatusAccepted,
		}}},
		state: fakeStateRepository{assignments: map[string]state.Assignment{"bus-1": {
			VehicleID: "bus-1", State: state.StateInService, TripID: "trip-1", Confidence: 0.9, ActiveFrom: now,
		}}},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/readiness", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"CAL-ITP-Style Readiness Workflow",
		"supports CAL-ITP-style readiness workflows",
		"does not claim CAL-ITP/Caltrans compliance",
		"Status source",
		"Next action",
		"Claim boundary",
		"Stable public URLs",
		"Static GTFS feed",
		"Vehicle Positions",
		"Trip Updates",
		"Alerts",
		"License/contact metadata",
		"Validation status",
		"Telemetry freshness",
		"Operations status",
		"Consumer packet preparedness",
		"feed discovery and published_feed records",
		"validation records",
		"telemetry latest rows",
		"scorecard snapshots",
		"docs/evidence tracker paths",
		"target-originated evidence",
		"Prepared packets are not submitted, under review, accepted, listed, displayed, or ingested.",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{
		"CAL-ITP/Caltrans compliant",
		"accepted by",
		"consumer ingestion confirmed",
		"production ready",
		"proves final-root",
		"hosted SaaS",
	} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("body overclaims %q: %s", forbidden, body)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/readiness?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict status = %d, want 403", rr.Code)
	}
}

func TestOperationsChecklistRoutesArePrivateScopedAndDeterministic(t *testing.T) {
	now := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{
		discovery: compliance.FeedDiscovery{
			AgencyID: "demo-agency", AgencyName: "Demo Agency", GeneratedAt: now, PublicationEnvironment: "pilot",
			PublicBaseURL:         "https://pilot.example.org",
			TechnicalContactEmail: "ops@example.org",
			License:               compliance.License{Name: "CC BY 4.0", URL: "https://example.org/license"},
			Feeds: []compliance.FeedMetadata{
				{FeedType: "schedule", CanonicalPublicURL: "https://pilot.example.org/public/gtfs/schedule.zip", ActiveFeedVersionID: "feed-v1", LastValidationStatus: "passed", LastValidationAt: &now},
				{FeedType: "vehicle_positions", CanonicalPublicURL: "http://localhost:8080/public/gtfsrt/vehicle_positions.pb", ActiveFeedVersionID: "feed-v1"},
				{FeedType: "trip_updates", CanonicalPublicURL: "https://feeds.agency.example/public/gtfsrt/trip_updates.pb", ActiveFeedVersionID: "feed-v1"},
				{FeedType: "alerts", CanonicalPublicURL: "https://feeds.real-agency.org/public/gtfsrt/alerts.pb", ActiveFeedVersionID: "feed-v1"},
			},
			Readiness: compliance.Readiness{AllRequiredFeedsListed: true, LicenseComplete: true, ContactComplete: true, HTTPSURLs: false},
		},
		scorecard: compliance.Scorecard{AgencyID: "demo-agency", SnapshotAt: now, OverallStatus: compliance.StatusYellow},
	}
	srv := newOperationsTestHandler(&handler{
		store:   store,
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now}}},
		telemetry: fakeTelemetryRepository{latest: []telemetry.StoredEvent{{
			Event: telemetry.Event{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Timestamp: now, Lat: 1, Lon: 2}, ReceivedAt: now,
		}}},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/checklist.json?agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json prefix", got)
	}
	var checklist operatorChecklistView
	if err := json.Unmarshal(rr.Body.Bytes(), &checklist); err != nil {
		t.Fatalf("decode checklist: %v", err)
	}
	assertChecklistShape(t, checklist)
	if checklist.AgencyID != "demo-agency" {
		t.Fatalf("agency_id = %q, want principal agency", checklist.AgencyID)
	}
	assertChecklistFlagsFalse(t, checklist.Flags)
	assertChecklistSafeStrings(t, rr.Body.String())
	assertChecklistDocsLinksSafe(t, checklist)
	assertChecklistNoPositiveClaims(t, rr.Body.String())

	groupIDs := make([]string, 0, len(checklist.Groups))
	for _, group := range checklist.Groups {
		groupIDs = append(groupIDs, group.ID)
	}
	wantGroups := []string{"setup", "feeds", "validation", "telemetry", "operations", "consumer_workflow"}
	if strings.Join(groupIDs, ",") != strings.Join(wantGroups, ",") {
		t.Fatalf("groups = %v, want %v", groupIDs, wantGroups)
	}
	labels := strings.Join(allHeuristicLabels(checklist), ",")
	for _, want := range []string{"pilot_or_reference_root", "local_only", "final_root_candidate_unverified", "no_final_root_evidence", "operator_entered_unverified", "approval_unknown"} {
		if !strings.Contains(labels, want) {
			t.Fatalf("heuristics %q missing %q", labels, want)
		}
	}

	normalized := checklist
	normalized.GeneratedAt = time.Time{}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/checklist.json", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	var second operatorChecklistView
	if err := json.Unmarshal(rr.Body.Bytes(), &second); err != nil {
		t.Fatalf("decode second checklist: %v", err)
	}
	second.GeneratedAt = time.Time{}
	if fmt.Sprintf("%#v", normalized) != fmt.Sprintf("%#v", second) {
		t.Fatalf("normalized checklist changed:\n%#v\n%#v", normalized, second)
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/checklist?agency_id=demo-agency", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("html Cache-Control = %q, want no-store", got)
	}
	body := rr.Body.String()
	for _, want := range []string{"Private Operator Checklist", "This checklist is private operator diagnostics", "not evidence", "not an evidence packet", "not compliance proof", "not agency approval", "not consumer acceptance", "not production readiness", "Setup", "Feeds", "Validation", "Telemetry", "Operations", "Consumer Workflow", "Placeholder-like", "Pilot/reference root", "No final-root evidence"} {
		if !strings.Contains(body, want) {
			t.Fatalf("html body missing %q: %s", want, body)
		}
	}
	assertChecklistSafeStrings(t, body)
	assertChecklistNoPositiveClaims(t, body)

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/checklist?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("html agency conflict status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/checklist.json?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("json agency conflict status = %d, want 403", rr.Code)
	}
}

func TestOperationsLaunchpadRoutesPrivateScopedGETOnlyNoStore(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		t.Run(string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			for _, path := range []string{"/admin/operations/launchpad", "/admin/operations/launchpad.json"} {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
				if rr.Code != http.StatusOK {
					t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
				}
				if got := rr.Header().Get("Cache-Control"); got != "no-store" {
					t.Fatalf("%s Cache-Control = %q, want no-store", path, got)
				}
			}
		})
	}

	unauth := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, authRejectAll{})
	for _, path := range []string{"/admin/operations/launchpad", "/admin/operations/launchpad.json"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		unauth.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("unauth %s status = %d, want 401", path, rr.Code)
		}
	}

	authenticated := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator}, Method: auth.MethodBearer,
	}})
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		for _, path := range []string{"/admin/operations/launchpad", "/admin/operations/launchpad.json"} {
			req := httptest.NewRequest(method, path, nil)
			rr := httptest.NewRecorder()
			authenticated.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status = %d, want 405", method, path, rr.Code)
			}
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/launchpad?agency_id=other-agency", nil)
	rr := httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict html status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/launchpad.json?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict json status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/launchpad.json", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public launchpad route status = %d, want 404", rr.Code)
	}
}

func TestOperationsLaunchpadJSONShapeFlagsAndNoLeakage(t *testing.T) {
	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{
		discovery: compliance.FeedDiscovery{
			AgencyID: "demo-agency", AgencyName: "Demo Agency", GeneratedAt: now, PublicationEnvironment: "pilot",
			PublicBaseURL:         "https://pilot.example.org",
			TechnicalContactEmail: "ops@example.org",
			License:               compliance.License{Name: "CC BY 4.0", URL: "https://example.org/license"},
			Feeds: []compliance.FeedMetadata{
				{FeedType: "schedule", CanonicalPublicURL: "https://pilot.example.org/public/gtfs/schedule.zip", ActiveFeedVersionID: "feed-v1", LastValidationStatus: "passed", LastValidationAt: &now},
				{FeedType: "vehicle_positions", CanonicalPublicURL: "https://pilot.example.org/public/gtfsrt/vehicle_positions.pb", ActiveFeedVersionID: "feed-v1"},
				{FeedType: "trip_updates", CanonicalPublicURL: "https://pilot.example.org/public/gtfsrt/trip_updates.pb", ActiveFeedVersionID: "feed-v1"},
				{FeedType: "alerts", CanonicalPublicURL: "https://pilot.example.org/public/gtfsrt/alerts.pb", ActiveFeedVersionID: "feed-v1"},
			},
			Readiness: compliance.Readiness{AllRequiredFeedsListed: true, LicenseComplete: true, ContactComplete: true, HTTPSURLs: true},
		},
		scorecard: compliance.Scorecard{AgencyID: "demo-agency", SnapshotAt: now, OverallStatus: compliance.StatusYellow},
	}
	srv := newOperationsTestHandler(&handler{
		store:   store,
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now}}},
		telemetry: fakeTelemetryRepository{latest: []telemetry.StoredEvent{{
			Event: telemetry.Event{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Timestamp: now, Lat: 1, Lon: 2}, ReceivedAt: now,
		}}},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/launchpad.json?agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json prefix", got)
	}
	var launchpad agencyLaunchpadView
	if err := json.Unmarshal(rr.Body.Bytes(), &launchpad); err != nil {
		t.Fatalf("decode launchpad: %v", err)
	}
	assertLaunchpadShape(t, launchpad)
	assertLaunchpadFlagsFalse(t, launchpad.ClaimFlags)
	assertLaunchpadSafeStrings(t, rr.Body.String())
	if launchpad.AgencyID != "demo-agency" {
		t.Fatalf("agency_id = %q, want demo-agency", launchpad.AgencyID)
	}
	var ids []string
	for _, section := range launchpad.Sections {
		ids = append(ids, section.ID)
	}
	wantIDs := []string{"setup", "gtfs", "metadata", "five_feeds", "telemetry", "validators", "readiness", "connector_conformance", "support_bundle", "decision_gate"}
	if strings.Join(ids, ",") != strings.Join(wantIDs, ",") {
		t.Fatalf("section ids = %v, want %v", ids, wantIDs)
	}
	for _, section := range launchpad.Sections {
		for _, link := range section.AdminLinks {
			if !strings.HasPrefix(link, "/admin/") {
				t.Fatalf("section %s unsafe admin link %q", section.ID, link)
			}
		}
	}

	missingHandler := newOperationsTestHandler(&handler{store: &fakePublicationStore{discoveryErr: errors.New("missing discovery"), scorecardErr: errors.New("missing scorecard"), consumersErr: errors.New("missing consumers")}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/launchpad.json", nil)
	rr = httptest.NewRecorder()
	missingHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("missing status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var missing agencyLaunchpadView
	if err := json.Unmarshal(rr.Body.Bytes(), &missing); err != nil {
		t.Fatalf("decode missing launchpad: %v", err)
	}
	for _, id := range []string{"gtfs", "five_feeds", "telemetry"} {
		if status := launchpadSectionStatus(missing, id); status == checklistStatusOK {
			t.Fatalf("missing-data section %s status = ok, want missing/unknown/review/blocker", id)
		}
	}
	assertLaunchpadFlagsFalse(t, missing.ClaimFlags)
}

func TestOperationsLaunchpadHTMLBoundariesNoFormsAndEscapes(t *testing.T) {
	store := &fakePublicationStore{
		discovery: compliance.FeedDiscovery{
			AgencyID: "demo-agency", AgencyName: `<script>alert("x")</script>`,
			License: compliance.License{Name: "Demo License", URL: "https://example.org/license"},
		},
		scorecardErr: errors.New("no scorecard"),
	}
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/launchpad", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Private Agency Launchpad", "creates no evidence", "contacts no external party", "changes no consumer status", "Setup", "GTFS", "Metadata", "Five feeds", "Telemetry", "Validators", "Readiness", "Connector conformance", "Support bundle", "Decision gate"} {
		if !strings.Contains(body, want) {
			t.Fatalf("html body missing %q: %s", want, body)
		}
	}
	if strings.Contains(body, `<script>alert("x")</script>`) {
		t.Fatalf("html did not escape script-like metadata: %s", body)
	}
	for _, forbidden := range []string{`<form`, `method="post"`, "/public/operations/launchpad", "agency approved", "consumer accepted", "production ready", "launch complete", "compliance achieved"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("launchpad html contains forbidden %q: %s", forbidden, body)
		}
	}
	assertLaunchpadSafeStrings(t, body)
}

func TestOperationsChecklistAccessMatrixMethodsAndRoutes(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		t.Run(string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			for _, path := range []string{"/admin/operations/checklist", "/admin/operations/checklist.json"} {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
				if rr.Code != http.StatusOK {
					t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
				}
			}
		})
	}

	unauth := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, authRejectAll{})
	for _, path := range []string{"/admin/operations/checklist", "/admin/operations/checklist.json"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		unauth.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("unauth %s status = %d, want 401", path, rr.Code)
		}
	}

	authenticated := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator}, Method: auth.MethodBearer,
	}})
	for _, path := range []string{"/admin/operations/checklist", "/admin/operations/checklist.json"} {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		rr := httptest.NewRecorder()
		authenticated.ServeHTTP(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("POST %s status = %d, want 405", path, rr.Code)
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/not-a-real-section", nil)
	rr := httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("unknown operations route status = %d, want 404", rr.Code)
	}
}

func TestOperationsChecklistHandlesMissingDataAndEscapesHTML(t *testing.T) {
	store := &fakePublicationStore{
		discovery: compliance.FeedDiscovery{
			AgencyID: "demo-agency", AgencyName: `<script>alert("x")</script>`, PublicBaseURL: "",
			License: compliance.License{Name: "Demo License", URL: "https://example.org/license"},
		},
		scorecardErr:         errors.New("no scorecard"),
		consumersErr:         errors.New("no consumers"),
		publicationConfigErr: errors.New("no publication config"),
	}
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/checklist", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if strings.Contains(body, `<script>alert("x")</script>`) || !strings.Contains(body, "&lt;script&gt;") {
		t.Fatalf("html did not escape script-like metadata: %s", body)
	}
	for _, want := range []string{"missing", "unknown", "needs_review"} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing-data html lacks status %q: %s", want, body)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/checklist.json", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("json status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var checklist operatorChecklistView
	if err := json.Unmarshal(rr.Body.Bytes(), &checklist); err != nil {
		t.Fatalf("decode checklist: %v", err)
	}
	assertChecklistShape(t, checklist)
	if !checklistContainsSignal(checklist, `<script>alert("x")</script>`) {
		t.Fatalf("json should preserve script-like metadata as data: %+v", checklist)
	}
}

func TestOperationsChecklistNavigationLinks(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	for _, path := range []string{"/admin/operations", "/admin/operations/setup", "/admin/operations/readiness"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
		}
		body := rr.Body.String()
		for _, want := range []string{"/admin/operations/checklist", "/admin/operations/checklist.json"} {
			if !strings.Contains(body, want) {
				t.Fatalf("%s missing link %q: %s", path, want, body)
			}
		}
	}
}

func TestDeploymentDoctorAndCaddyLocalRouteGuards(t *testing.T) {
	doctor, err := os.ReadFile(filepath.Join("..", "..", "scripts", "deployment-doctor.sh"))
	if err != nil {
		t.Fatalf("read deployment doctor: %v", err)
	}
	text := string(doctor)
	if !strings.Contains(text, `"/admin/gtfs-studio"`) {
		t.Fatalf("deployment doctor missing /admin/gtfs-studio private route")
	}
	if strings.Contains(text, `"/admin/gtfs"`) {
		t.Fatalf("deployment doctor still checks exact /admin/gtfs")
	}
	if !strings.Contains(text, `"/admin/operations/validation-health.json"`) {
		t.Fatalf("deployment doctor missing private validation-health JSON check")
	}
	if strings.Contains(text, `validation-health"`) && strings.Contains(text, `-X POST`) {
		t.Fatalf("deployment doctor must not POST validation-health")
	}
	caddy, err := os.ReadFile(filepath.Join("..", "..", "deploy", "Caddyfile.local"))
	if err != nil {
		t.Fatalf("read caddyfile: %v", err)
	}
	caddyText := string(caddy)
	caddyLines := nonCommentCaddyLines(caddyText)
	if !containsString(caddyLines, "@local_root {") || !containsString(caddyLines, "path /") || !containsString(caddyLines, `respond @local_root "Open Transit RT local app is running. Public feeds are under /public/ and admin routes require auth." 200`) {
		t.Fatalf("local Caddyfile missing exact root 200 handler:\n%s", caddyText)
	}
	if !containsString(caddyLines, `respond "not found" 404`) {
		t.Fatalf("local Caddyfile missing unmatched 404 fallback:\n%s", caddyText)
	}
	var lastRespond string
	for _, line := range caddyLines {
		if strings.HasPrefix(line, "respond ") {
			lastRespond = line
		}
		if strings.HasPrefix(line, `respond "`) && strings.HasSuffix(line, `" 200`) {
			t.Fatalf("local Caddyfile has unconditional 200 catch-all:\n%s", caddyText)
		}
	}
	if lastRespond != `respond "not found" 404` {
		t.Fatalf("local Caddyfile final respond = %q, want unmatched 404 fallback:\n%s", lastRespond, caddyText)
	}
}

func TestValidatorHealthScriptDryRunOutputSafety(t *testing.T) {
	tmp := t.TempDir()
	fakeBin := filepath.Join(tmp, "bin")
	if err := os.Mkdir(fakeBin, 0o700); err != nil {
		t.Fatal(err)
	}
	curlPath := filepath.Join(fakeBin, "curl")
	if err := os.WriteFile(curlPath, []byte("#!/bin/sh\necho curl must not be called in dry-run >&2\nexit 99\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	outputDir := filepath.Join(tmp, "validator-health")
	cmd := exec.Command("sh", filepath.Join("..", "..", "scripts", "validator-health.sh"), "--dry-run")
	cmd.Env = append(os.Environ(),
		"PATH="+fakeBin+":"+os.Getenv("PATH"),
		"OUTPUT_DIR="+outputDir,
		"ALLOW_UNIGNORED_OUTPUT_DIR=true",
		"ADMIN_TOKEN=not-used-in-dry-run",
		"PUBLIC_BASE_URL=http://127.0.0.1:65535",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("dry-run failed: %v\n%s", err, out)
	}
	for _, name := range []string{"summary.json", "summary.md", "manifest.json", "manifest.md"} {
		if _, err := os.Stat(filepath.Join(outputDir, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
	summaryBytes, err := os.ReadFile(filepath.Join(outputDir, "summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	var summary compliance.ValidationHealthSummary
	if err := json.Unmarshal(summaryBytes, &summary); err != nil {
		t.Fatalf("summary json invalid: %v", err)
	}
	assertValidationHealthSummaryShape(t, summary)
	manifestBytes, err := os.ReadFile(filepath.Join(outputDir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("manifest json invalid: %v", err)
	}
	for _, key := range []string{"included_files", "excluded_categories", "output_dir", "created_at"} {
		if _, ok := manifest[key]; !ok {
			t.Fatalf("manifest missing %s: %+v", key, manifest)
		}
	}
	for _, name := range []string{"summary.json", "manifest.json", "summary.md", "manifest.md"} {
		body, err := os.ReadFile(filepath.Join(outputDir, name))
		if err != nil {
			t.Fatal(err)
		}
		assertValidatorHealthGeneratedFileNoSecretValues(t, string(body))
		if len(body) > 16000 {
			t.Fatalf("%s size = %d, want bounded", name, len(body))
		}
	}
	manifestMD, err := os.ReadFile(filepath.Join(outputDir, "manifest.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"raw validator reports", "Authorization headers", "cookies", "tokens", "database URLs", "private paths", "evidence packets", "consumer submission artifacts"} {
		if !strings.Contains(string(manifestMD), want) {
			t.Fatalf("manifest.md missing %q", want)
		}
	}
}

func TestValidatorHealthScriptReadOnlyNoTokenNoNetwork(t *testing.T) {
	tmp := t.TempDir()
	fakeBin := filepath.Join(tmp, "bin")
	if err := os.Mkdir(fakeBin, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fakeBin, "curl"), []byte("#!/bin/sh\necho unexpected network >&2\nexit 99\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	outputDir := filepath.Join(tmp, "read-only")
	cmd := exec.Command("sh", filepath.Join("..", "..", "scripts", "validator-health.sh"))
	cmd.Env = append(os.Environ(),
		"PATH="+fakeBin+":"+os.Getenv("PATH"),
		"OUTPUT_DIR="+outputDir,
		"ALLOW_UNIGNORED_OUTPUT_DIR=true",
		"VALIDATOR_TOOLING_MODE=stub",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("read-only no-token run failed: %v\n%s", err, out)
	}
	assertValidatorHealthScriptOutputFilesSafe(t, outputDir)
}

func TestValidatorHealthScriptRunValidatorsPostsActionAndCSRF(t *testing.T) {
	tmp := t.TempDir()
	fakeBin := filepath.Join(tmp, "bin")
	if err := os.Mkdir(fakeBin, 0o700); err != nil {
		t.Fatal(err)
	}
	capture := filepath.Join(tmp, "curl.args")
	fakeCurl := `#!/bin/sh
printf '%s\n' "$@" >"$CURL_CAPTURE"
found_action=false
found_csrf=false
prev=""
for arg in "$@"; do
  if [ "$prev" = "--data-urlencode" ] && [ "$arg" = "action=run_all" ]; then found_action=true; fi
  if [ "$prev" = "--data-urlencode" ] && [ "$arg" = "csrf_token=valid-csrf" ]; then found_csrf=true; fi
  prev="$arg"
done
if [ "$found_action" != "true" ] || [ "$found_csrf" != "true" ]; then
  exit 64
fi
cat <<'JSON'
{"generated_at":"2026-05-09T12:00:00Z","agency_id":"demo-agency","overall_status":"recorded","tooling_status":"configured","feeds":[{"feed_type":"schedule","validator_id":"static-mobilitydata","validator_name":"mobilitydata-gtfs-validator","tooling_status":"configured","artifact_status":"available","latest_result_status":"recorded","latest_result_at":"2026-05-09T12:00:00Z","active_feed_version_id":"feed-v1","latest_result_feed_version_id":"feed-v1","stale_status":"current","health_status":"recorded","next_action":"Keep this as private diagnostics.","claim_boundary":"Private diagnostics only."},{"feed_type":"vehicle_positions","validator_id":"realtime-mobilitydata","validator_name":"mobilitydata-gtfs-realtime-validator","tooling_status":"configured","artifact_status":"available","latest_result_status":"recorded","latest_result_at":"2026-05-09T12:00:00Z","active_feed_version_id":"feed-v1","latest_result_feed_version_id":"feed-v1","stale_status":"current","health_status":"recorded","next_action":"Keep this as private diagnostics.","claim_boundary":"Private diagnostics only."},{"feed_type":"trip_updates","validator_id":"realtime-mobilitydata","validator_name":"mobilitydata-gtfs-realtime-validator","tooling_status":"configured","artifact_status":"available","latest_result_status":"recorded","latest_result_at":"2026-05-09T12:00:00Z","active_feed_version_id":"feed-v1","latest_result_feed_version_id":"feed-v1","stale_status":"current","health_status":"recorded","next_action":"Keep this as private diagnostics.","claim_boundary":"Private diagnostics only."},{"feed_type":"alerts","validator_id":"realtime-mobilitydata","validator_name":"mobilitydata-gtfs-realtime-validator","tooling_status":"configured","artifact_status":"available","latest_result_status":"recorded","latest_result_at":"2026-05-09T12:00:00Z","active_feed_version_id":"feed-v1","latest_result_feed_version_id":"feed-v1","stale_status":"current","health_status":"recorded","next_action":"Keep this as private diagnostics.","claim_boundary":"Private diagnostics only."}],"external_evidence_created":false,"consumer_statuses_changed":false,"compliance_claimed":false,"production_readiness_claimed":false}
JSON
`
	if err := os.WriteFile(filepath.Join(fakeBin, "curl"), []byte(fakeCurl), 0o700); err != nil {
		t.Fatal(err)
	}
	outputDir := filepath.Join(tmp, "run-all")
	cmd := exec.Command("sh", filepath.Join("..", "..", "scripts", "validator-health.sh"))
	cmd.Env = append(os.Environ(),
		"PATH="+fakeBin+":"+os.Getenv("PATH"),
		"OUTPUT_DIR="+outputDir,
		"ALLOW_UNIGNORED_OUTPUT_DIR=true",
		"RUN_VALIDATORS=true",
		"ADMIN_TOKEN=script-admin-token-value",
		"CSRF_TOKEN=valid-csrf",
		"ADMIN_BASE_URL=http://127.0.0.1:8080",
		"VALIDATOR_TOOLING_MODE=stub",
		"CURL_CAPTURE="+capture,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run validators failed: %v\n%s", err, out)
	}
	args, err := os.ReadFile(capture)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(args), "action=run_all") || !strings.Contains(string(args), "csrf_token=valid-csrf") {
		t.Fatalf("curl args missing action/csrf: %s", args)
	}
	assertValidatorHealthScriptOutputFilesSafe(t, outputDir)
}

func TestValidatorHealthScriptRunValidatorsBlockedCSRFOutput(t *testing.T) {
	tmp := t.TempDir()
	fakeBin := filepath.Join(tmp, "bin")
	if err := os.Mkdir(fakeBin, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fakeBin, "curl"), []byte("#!/bin/sh\nprintf 'forbidden\\n' >&2\nexit 22\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	outputDir := filepath.Join(tmp, "blocked")
	baseEnv := append(os.Environ(),
		"PATH="+fakeBin+":"+os.Getenv("PATH"),
		"OUTPUT_DIR="+outputDir,
		"ALLOW_UNIGNORED_OUTPUT_DIR=true",
		"RUN_VALIDATORS=true",
		"ADMIN_TOKEN=script-admin-token-value",
		"CSRF_TOKEN=invalid-csrf",
		"ADMIN_BASE_URL=http://127.0.0.1:8080",
		"VALIDATOR_TOOLING_MODE=stub",
	)
	cmd := exec.Command("sh", filepath.Join("..", "..", "scripts", "validator-health.sh"))
	cmd.Env = baseEnv
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("non-strict blocked run should exit 0: %v\n%s", err, out)
	}
	assertValidatorHealthScriptOutputFilesSafe(t, outputDir)
	summaryBytes, err := os.ReadFile(filepath.Join(outputDir, "summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(summaryBytes), `"overall_status": "blocked"`) {
		t.Fatalf("summary did not record blocked status: %s", summaryBytes)
	}
	strictOutput := filepath.Join(tmp, "blocked-strict")
	cmd = exec.Command("sh", filepath.Join("..", "..", "scripts", "validator-health.sh"))
	cmd.Env = append(baseEnv, "OUTPUT_DIR="+strictOutput, "STRICT_VALIDATOR_HEALTH=true")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("strict blocked run should exit nonzero: %s", out)
	}
	assertValidatorHealthScriptOutputFilesSafe(t, strictOutput)
}

func TestValidatorHealthScriptRejectsEvidenceOutput(t *testing.T) {
	cmd := exec.Command("sh", filepath.Join("..", "..", "scripts", "validator-health.sh"), "--dry-run")
	cmd.Env = append(os.Environ(), "OUTPUT_DIR=docs/evidence/validator-health-test", "ALLOW_UNIGNORED_OUTPUT_DIR=true")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected evidence output rejection, got success: %s", out)
	}
}

func TestOperationsSetupPublicationFormRequiresAdminAndDerivesAgencyID(t *testing.T) {
	store := &fakePublicationStore{}
	srv := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	form := "action=publication_bootstrap&public_base_url=https%3A%2F%2Fagency.example&feed_base_url=https%3A%2F%2Fagency.example%2Ffeeds"
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/setup", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rr.Code)
	}

	srv = newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/setup", strings.NewReader(form+"&technical_contact_email= ops%40agency.example &license_name= CC-BY &license_url=https%3A%2F%2Fagency.example%2Flicense&publication_environment= pilot "))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.bootstrapInput.AgencyID != "demo-agency" || store.bootstrapInput.ActorID != "admin@example.com" {
		t.Fatalf("bootstrap input identity = %+v, want authenticated principal", store.bootstrapInput)
	}
	if store.bootstrapInput.TechnicalContactEmail != "ops@agency.example" || store.bootstrapInput.PublicationEnvironment != "pilot" {
		t.Fatalf("bootstrap input not trimmed = %+v", store.bootstrapInput)
	}
}

func TestOperationsSetupPublicationFormRejectsConflictingAgencyID(t *testing.T) {
	store := &fakePublicationStore{}
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}, telemetry: fakeTelemetryRepository{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	form := "action=publication_bootstrap&agency_id=other-agency&public_base_url=https%3A%2F%2Fagency.example&feed_base_url=https%3A%2F%2Fagency.example%2Ffeeds"
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/setup", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rr.Code)
	}
	if store.bootstrapInput.AgencyID != "" {
		t.Fatalf("bootstrap should not run on agency conflict: %+v", store.bootstrapInput)
	}
}

func TestOperationsSetupValidationFormMapsFeedTypeServerSide(t *testing.T) {
	validatorPath := writeRealtimeValidator(t)
	t.Setenv("GTFS_RT_VALIDATOR_PATH", validatorPath)
	t.Setenv("GTFS_RT_VALIDATOR_VERSION", "test-validator")
	t.Setenv("GTFS_RT_VALIDATOR_ARGS", "")

	store := &fakePublicationStore{}
	artifacts := &fakeRealtimeArtifacts{payloads: map[string][]byte{"alerts": []byte("protobuf-alerts")}}
	handler := newOperationsTestHandler(&handler{
		store:    store,
		schedule: fakeScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "demo-agency", FeedVersionID: "feed-demo", RevisionTime: time.Now().UTC(), Payload: []byte("schedule zip bytes")}},
		realtime: artifacts,
		devices:  fakeDeviceStore{},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	form := "action=run_validation&feed_type=alerts"
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/setup", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if artifacts.calls["alerts"] != 1 {
		t.Fatalf("artifact calls = %+v, want alerts validation via server mapping", artifacts.calls)
	}
	if store.result.FeedType != "alerts" || store.result.ValidatorName != "mobilitydata-gtfs-realtime-validator" {
		t.Fatalf("stored validation result = %+v, want realtime validator selected by feed type", store.result)
	}
	body := rr.Body.String()
	for _, forbidden := range []string{"validator_id", "argv"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("setup validation response leaks browser-supplied or raw validator detail %q: %s", forbidden, body)
		}
	}
}

func TestOperationsSetupValidationFormRejectsUnsafeBrowserFields(t *testing.T) {
	handler := newOperationsTestHandler(&handler{
		store:    &fakePublicationStore{},
		schedule: fakeScheduleBuilder{},
		realtime: &fakeRealtimeArtifacts{},
		devices:  fakeDeviceStore{},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	for _, form := range []string{
		"action=run_validation&feed_type=alerts&validator_id=realtime-mobilitydata",
		"action=run_validation&feed_type=alerts&realtime_pb_path=%2Ftmp%2Fevil.pb",
		"action=run_validation&feed_type=schedule&output_path=%2Ftmp%2Freport",
	} {
		req := httptest.NewRequest(http.MethodPost, "/admin/operations/setup", strings.NewReader(form))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 with safe rendered form error: %s", rr.Code, rr.Body.String())
		}
		body := rr.Body.String()
		if !strings.Contains(body, "validation setup form only accepts feed type") {
			t.Fatalf("body does not contain safe form error: %s", body)
		}
		if strings.Contains(body, "/tmp/evil.pb") || strings.Contains(body, "/tmp/report") {
			t.Fatalf("body leaks browser-supplied path: %s", body)
		}
	}
}

func TestOperationsSetupCookiePostRequiresCSRF(t *testing.T) {
	cfg := auth.JWTConfig{Secrets: []string{"test-secret"}, Issuer: "test-issuer", Audience: "test-audience", TTL: time.Hour}
	signer, err := auth.NewSigner(cfg)
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	token, _, err := signer.Sign("admin@example.com", "demo-agency", time.Hour)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	verifier, err := auth.NewVerifier(cfg)
	if err != nil {
		t.Fatalf("verifier: %v", err)
	}
	middleware := auth.NewMiddleware(verifier, auth.StaticRoleStore{Roles: []auth.Role{auth.RoleAdmin}}, "csrf-secret")
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}, csrfSecret: "csrf-secret"}, middleware)
	form := "action=publication_bootstrap&public_base_url=https%3A%2F%2Fagency.example&feed_base_url=https%3A%2F%2Fagency.example%2Ffeeds"
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/setup", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "admin_session", Value: token})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 for missing CSRF", rr.Code)
	}
}

func TestGTFSQualityRouteAuthMatrixAndMethods(t *testing.T) {
	roles := []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin}
	for _, role := range roles {
		t.Run("get_"+string(role), func(t *testing.T) {
			handler := newGTFSQualityTestHandler(t, role, &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC())}, fakeScheduleBuilder{})
			req := httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-quality", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("GET status = %d, want 200", rr.Code)
			}
			if got := rr.Header().Get("Cache-Control"); got != "no-store" {
				t.Fatalf("Cache-Control = %q, want no-store", got)
			}
		})
	}
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor} {
		t.Run("post_"+string(role), func(t *testing.T) {
			handler := newGTFSQualityTestHandler(t, role, &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC())}, fakeScheduleBuilder{})
			req := httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-quality", strings.NewReader(gtfsQualityRerunForm()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusForbidden {
				t.Fatalf("POST status = %d, want 403", rr.Code)
			}
		})
	}
	for _, method := range []string{http.MethodPut, http.MethodPatch, http.MethodDelete} {
		req := httptest.NewRequest(method, "/admin/operations/gtfs-quality", nil)
		rr := httptest.NewRecorder()
		newGTFSQualityTestHandler(t, auth.RoleAdmin, &fakePublicationStore{}, fakeScheduleBuilder{}).ServeHTTP(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s status = %d, want 405", method, rr.Code)
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/not-real", nil)
	rr := httptest.NewRecorder()
	newGTFSQualityTestHandler(t, auth.RoleAdmin, &fakePublicationStore{}, fakeScheduleBuilder{}).ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("unknown operations route = %d, want 404", rr.Code)
	}
}

func TestGTFSQualityCookiePostRequiresCSRF(t *testing.T) {
	cfg := auth.JWTConfig{Secrets: []string{"test-secret"}, Issuer: "test-issuer", Audience: "test-audience", TTL: time.Hour}
	signer, err := auth.NewSigner(cfg)
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	token, _, err := signer.Sign("admin@example.com", "demo-agency", time.Hour)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	verifier, err := auth.NewVerifier(cfg)
	if err != nil {
		t.Fatalf("verifier: %v", err)
	}
	middleware := auth.NewMiddleware(verifier, auth.StaticRoleStore{Roles: []auth.Role{auth.RoleAdmin}}, "csrf-secret")
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC())}, devices: fakeDeviceStore{}, schedule: fakeScheduleBuilder{}}, middleware)
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-quality", strings.NewReader("action=rerun_static_validator"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "admin_session", Value: token})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 for missing CSRF", rr.Code)
	}
}

func TestGTFSQualityPostStrictnessAndBodyCap(t *testing.T) {
	validatorPath := writeScheduleValidatorWithReport(t, `{"status":"warning","notices":[{"code":"route_short_name_too_long","severity":"WARNING","message":"review"}]}`)
	t.Setenv("GTFS_VALIDATOR_PATH", validatorPath)
	store := &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC())}
	scheduleBuilder := &countingScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "demo-agency", FeedVersionID: "feed-v1", RevisionTime: time.Now().UTC(), Payload: []byte("schedule zip")}}
	handler := newGTFSQualityTestHandler(t, auth.RoleAdmin, store, scheduleBuilder)
	for _, body := range []string{
		gtfsQualityRerunForm() + "&validator_id=static-mobilitydata",
		gtfsQualityRerunForm() + "&schedule_zip_path=/tmp/private/schedule.zip",
		gtfsQualityRerunForm() + "&URL=https%3A%2F%2Fevil.example%2Fgtfs.zip",
		gtfsQualityRerunForm() + "&timeout=1",
	} {
		req := httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-quality", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("strict field status = %d, want rendered blocker", rr.Code)
		}
		if scheduleBuilder.snapshotCalls != 0 || store.storeValidationCalls != 0 {
			t.Fatalf("unsafe form invoked validator: schedule=%d stores=%d", scheduleBuilder.snapshotCalls, store.storeValidationCalls)
		}
	}
	large := gtfsQualityRerunForm() + "&" + strings.Repeat("x", gtfsQualityPostMaxBytes+1)
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-quality", strings.NewReader(large))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized status = %d, want 413", rr.Code)
	}
	if scheduleBuilder.snapshotCalls != 0 || store.storeValidationCalls != 0 {
		t.Fatalf("oversized form invoked validator: schedule=%d stores=%d", scheduleBuilder.snapshotCalls, store.storeValidationCalls)
	}
}

func TestGTFSQualityPostRerunAndFailureBoundaries(t *testing.T) {
	validatorPath := writeScheduleValidatorWithReport(t, `{"status":"warning","notices":[{"code":"route_short_name_too_long","severity":"WARNING","message":"review /tmp/private <script>"}]}`)
	t.Setenv("GTFS_VALIDATOR_PATH", validatorPath)
	store := &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC())}
	scheduleBuilder := &countingScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "demo-agency", FeedVersionID: "feed-v1", RevisionTime: time.Now().UTC(), Payload: []byte("schedule zip")}}
	handler := newGTFSQualityTestHandler(t, auth.RoleAdmin, store, scheduleBuilder)
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-quality", strings.NewReader(gtfsQualityRerunForm()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	body := rr.Body.String()
	if rr.Code != http.StatusOK || !strings.Contains(body, "needs_review") {
		t.Fatalf("rerun status/body = %d %s, want needs_review", rr.Code, body)
	}
	for _, forbidden := range []string{"raw_report", "stdout", "stderr", "argv", "/tmp/private", "validator_clean", "compliant", "consumer_ready", "production_ready"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("rerun body leaked/claimed %q", forbidden)
		}
	}
	if store.result.ValidatorName != compliance.CanonicalStaticValidatorName || store.result.FeedVersionID != "feed-v1" {
		t.Fatalf("stored result = %+v, want canonical static for active feed", store.result)
	}

	t.Setenv("GTFS_VALIDATOR_PATH", "")
	store = &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC())}
	handler = newGTFSQualityTestHandler(t, auth.RoleAdmin, store, fakeScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "demo-agency", FeedVersionID: "feed-v1", RevisionTime: time.Now().UTC(), Payload: []byte("schedule zip")}})
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-quality", strings.NewReader(gtfsQualityRerunForm()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "blocking") {
		t.Fatalf("missing validator body = %d %s, want blocker page", rr.Code, rr.Body.String())
	}
	for _, forbidden := range []string{"stack trace", "stdout", "stderr", "/tmp/"} {
		if strings.Contains(rr.Body.String(), forbidden) {
			t.Fatalf("failure body leaked %q", forbidden)
		}
	}
}

func TestGTFSQualityConcurrentAdminReruns(t *testing.T) {
	validatorPath := writeScheduleValidatorWithReport(t, `{"status":"failed","notices":[{"code":"duplicate_trip_id","severity":"ERROR","message":"duplicate"}]}`)
	t.Setenv("GTFS_VALIDATOR_PATH", validatorPath)
	store := &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC())}
	handler := newGTFSQualityTestHandler(t, auth.RoleAdmin, store, &countingScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "demo-agency", FeedVersionID: "feed-v1", RevisionTime: time.Now().UTC(), Payload: []byte("schedule zip")}})
	errs := make(chan string, 2)
	for i := 0; i < 2; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-quality", strings.NewReader(gtfsQualityRerunForm()+"&validator_id=browser-supplied"))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK || strings.Contains(rr.Body.String(), "raw_report") || strings.Contains(rr.Body.String(), "/tmp/private") {
				errs <- fmt.Sprintf("status/body = %d %s", rr.Code, rr.Body.String())
				return
			}
			errs <- ""
		}()
	}
	for i := 0; i < 2; i++ {
		if err := <-errs; err != "" {
			t.Fatal(err)
		}
	}
	if store.storeValidationCalls != 0 {
		t.Fatalf("unsafe repeated POST stored validation rows: %d", store.storeValidationCalls)
	}
}

func TestGTFSQualityWarningOnlyWording(t *testing.T) {
	store := &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC()), validationRecords: []compliance.ValidationReportRecord{{
		ID: 1, CreatedAt: time.Now().UTC(), Result: compliance.ValidationResult{AgencyID: "demo-agency", FeedType: "schedule", FeedVersionID: "feed-v1", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "warning", WarningCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{noticeMap("route_short_name_too_long", "WARNING")}}}},
	}}}
	handler := newGTFSQualityTestHandler(t, auth.RoleAdmin, store, fakeScheduleBuilder{})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-quality", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	body := rr.Body.String()
	if !strings.Contains(body, "needs_review") {
		t.Fatalf("warning-only body missing needs_review: %s", body)
	}
	for _, forbidden := range []string{"validator_clean", "compliant", "accepted", "consumer_ready", "production_ready"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("warning-only body emitted forbidden positive label %q", forbidden)
		}
	}
}

func TestGTFSQualityGETReadOnlyAndAgencyIsolation(t *testing.T) {
	store := &fakePublicationStore{discoveries: map[string]compliance.FeedDiscovery{
		"agency-a": gtfsQualityDiscovery(time.Now().UTC()),
		"agency-b": {AgencyID: "agency-b", Feeds: []compliance.FeedMetadata{{FeedType: "schedule", ActiveFeedVersionID: "feed-b"}}},
	}, validationRecords: []compliance.ValidationReportRecord{
		{ID: 1, CreatedAt: time.Now().UTC(), Result: compliance.ValidationResult{AgencyID: "agency-b", FeedType: "schedule", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "failed", ErrorCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{noticeMap("duplicate_route_id", "ERROR")}}}}},
		{ID: 2, CreatedAt: time.Now().UTC(), Result: compliance.ValidationResult{AgencyID: "agency-a", FeedType: "schedule", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "warning", WarningCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{noticeMap("unused_shape", "WARNING")}}}}},
	}}
	schedule := &countingScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "agency-a", FeedVersionID: "feed-a", RevisionTime: time.Now().UTC(), Payload: []byte("zip")}}
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}, schedule: schedule, csrfSecret: "test-csrf"}, auth.TestAuthenticator{Principal: auth.Principal{Subject: "viewer@example.com", AgencyID: "agency-a", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-quality?agency_id=agency-a", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("same agency status = %d", rr.Code)
	}
	if strings.Contains(rr.Body.String(), "duplicate_route_id") || !strings.Contains(rr.Body.String(), "unused_shape") {
		t.Fatalf("agency isolation body = %s", rr.Body.String())
	}
	if schedule.snapshotCalls != 0 || store.storeValidationCalls != 0 {
		t.Fatalf("GET caused side effects: schedule=%d stores=%d", schedule.snapshotCalls, store.storeValidationCalls)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-quality?agency_id=agency-b", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("different agency status = %d, want 403", rr.Code)
	}
}

func TestGTFSQualityLatestSelectionAndEmptyStates(t *testing.T) {
	now := time.Now().UTC()
	store := &fakePublicationStore{discovery: gtfsQualityDiscovery(now), validationRecords: []compliance.ValidationReportRecord{
		{ID: 1, CreatedAt: now.Add(-time.Hour), Result: compliance.ValidationResult{AgencyID: "demo-agency", FeedType: "schedule", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "failed", ErrorCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{noticeMap("newer_lower_id", "ERROR")}}}}},
		{ID: 10, CreatedAt: now, Result: compliance.ValidationResult{AgencyID: "demo-agency", FeedType: "schedule", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "warning", WarningCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{noticeMap("older", "WARNING")}}}}},
		{ID: 11, CreatedAt: now, Result: compliance.ValidationResult{AgencyID: "demo-agency", FeedType: "schedule", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "failed", ErrorCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{noticeMap("higher_id", "ERROR")}}}}},
		{ID: 2, CreatedAt: now.Add(time.Hour), Result: compliance.ValidationResult{AgencyID: "demo-agency", FeedType: "schedule", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "warning", WarningCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{noticeMap("newest_lower_id", "WARNING")}}}}},
		{ID: 99, CreatedAt: now.Add(time.Hour), Result: compliance.ValidationResult{AgencyID: "demo-agency", FeedType: "vehicle_positions", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "failed", ErrorCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{noticeMap("wrong_feed", "ERROR")}}}}},
		{ID: 100, CreatedAt: now.Add(time.Hour), Result: compliance.ValidationResult{AgencyID: "demo-agency", FeedType: "schedule", ValidatorName: "other-validator", Status: "failed", ErrorCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{noticeMap("wrong_validator", "ERROR")}}}}},
		{ID: 101, CreatedAt: now.Add(time.Hour), Result: compliance.ValidationResult{AgencyID: "other", FeedType: "schedule", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "failed", ErrorCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{noticeMap("wrong_agency", "ERROR")}}}}},
	}}
	handler := newGTFSQualityTestHandler(t, auth.RoleAdmin, store, fakeScheduleBuilder{})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-quality", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	body := rr.Body.String()
	if !strings.Contains(body, "newest_lower_id") || strings.Contains(body, "higher_id") || strings.Contains(body, "wrong_feed") || strings.Contains(body, "wrong_validator") || strings.Contains(body, "wrong_agency") {
		t.Fatalf("latest selection body = %s", body)
	}
	empty := &fakePublicationStore{discovery: gtfsQualityDiscovery(now)}
	handler = newGTFSQualityTestHandler(t, auth.RoleAdmin, empty, fakeScheduleBuilder{})
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-quality", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "Run the appropriate validation") {
		t.Fatalf("empty canonical state = %d %s", rr.Code, rr.Body.String())
	}
	missingSchedule := &fakePublicationStore{discovery: compliance.FeedDiscovery{AgencyID: "demo-agency"}}
	handler = newGTFSQualityTestHandler(t, auth.RoleAdmin, missingSchedule, fakeScheduleBuilder{})
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-quality", strings.NewReader(gtfsQualityRerunForm()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "no active published schedule") {
		t.Fatalf("missing schedule = %d %s", rr.Code, rr.Body.String())
	}
}

func TestGTFSQualityHTMLEscapingAndUsabilitySmoke(t *testing.T) {
	store := &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC()), validationRecords: []compliance.ValidationReportRecord{{
		ID: 1, CreatedAt: time.Now().UTC(), Result: compliance.ValidationResult{AgencyID: "demo-agency", FeedType: "schedule", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "warning", WarningCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{map[string]any{"code": "route_short_name_too_long", "severity": "WARNING", "message": "<script>alert(1)</script>"}}}}},
	}, {
		ID: 2, CreatedAt: time.Now().UTC(), Result: compliance.ValidationResult{AgencyID: "demo-agency", FeedType: "schedule", ValidatorName: compliance.InternalGTFSImportValidatorName, Status: "failed", ErrorCount: 1, Report: map[string]any{"errors": []any{map[string]any{"code": "missing_trip_reference", "message": "missing reference"}}}},
	}}}
	handler := newGTFSQualityTestHandler(t, auth.RoleAdmin, store, fakeScheduleBuilder{})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-quality", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	body := rr.Body.String()
	for _, want := range []string{"<h2>GTFS Quality Triage</h2>", "Canonical MobilityData static validator", "Open Transit RT internal import validation", "Validator output is diagnostics", "Rerun static MobilityData validator"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q", want)
		}
	}
	if strings.Contains(body, "<script>alert(1)</script>") || !strings.Contains(body, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Fatalf("script was not escaped: %s", body)
	}
	readOnly := newGTFSQualityTestHandler(t, auth.RoleReadOnly, store, fakeScheduleBuilder{})
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-quality", nil)
	rr = httptest.NewRecorder()
	readOnly.ServeHTTP(rr, req)
	if strings.Contains(rr.Body.String(), "Rerun static MobilityData validator</button>") {
		t.Fatalf("read-only page showed admin rerun form")
	}
}

func TestGTFSQualityPageLargeReportRender(t *testing.T) {
	store := &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC()), validationRecords: []compliance.ValidationReportRecord{largeOperationsCanonicalRecord(50000)}}
	handler := newGTFSQualityTestHandler(t, auth.RoleAdmin, store, fakeScheduleBuilder{})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-quality", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if rr.Body.Len() > 650000 {
		t.Fatalf("rendered body = %d bytes, want bounded", rr.Body.Len())
	}
	for _, forbidden := range []string{"raw_report", "stdout", "stderr", "argv", "/tmp/private", strings.Repeat("x", 1000)} {
		if strings.Contains(rr.Body.String(), forbidden) {
			t.Fatalf("large render leaked %q", forbidden)
		}
	}
}

func BenchmarkRenderGTFSQualityPage(b *testing.B) {
	store := &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC()), validationRecords: []compliance.ValidationReportRecord{largeOperationsCanonicalRecord(50000)}}
	handler := newGTFSQualityTestHandler(b, auth.RoleAdmin, store, fakeScheduleBuilder{})
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-quality", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			b.Fatalf("status = %d", rr.Code)
		}
	}
}

func TestValidationHealthRouteAuthMatrixMethodsAndHeaders(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		handler := newValidationHealthTestHandler(t, role, &fakePublicationStore{discovery: validationHealthTestDiscovery(time.Now().UTC())}, fakeScheduleBuilder{})
		for _, path := range []string{"/admin/operations/validation-health", "/admin/operations/validation-health.json"} {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("%s role %s status = %d, want 200: %s", path, role, rr.Code, rr.Body.String())
			}
			if got := rr.Header().Get("Cache-Control"); got != "no-store" {
				t.Fatalf("%s Cache-Control = %q, want no-store", path, got)
			}
			if strings.HasSuffix(path, ".json") && !strings.HasPrefix(rr.Header().Get("Content-Type"), "application/json") {
				t.Fatalf("%s Content-Type = %q, want application/json", path, rr.Header().Get("Content-Type"))
			}
		}
	}
	unauth := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, authRejectAll{})
	for _, path := range []string{"/admin/operations/validation-health", "/admin/operations/validation-health.json"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		unauth.ServeHTTP(rr, req)
		if rr.Code == http.StatusOK {
			t.Fatalf("%s unauthenticated status = 200, want rejection", path)
		}
	}
	for _, method := range []string{http.MethodPut, http.MethodPatch, http.MethodDelete} {
		req := httptest.NewRequest(method, "/admin/operations/validation-health", nil)
		rr := httptest.NewRecorder()
		newValidationHealthTestHandler(t, auth.RoleAdmin, &fakePublicationStore{}, fakeScheduleBuilder{}).ServeHTTP(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s status = %d, want 405", method, rr.Code)
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/not-real", nil)
	rr := httptest.NewRecorder()
	newValidationHealthTestHandler(t, auth.RoleAdmin, &fakePublicationStore{}, fakeScheduleBuilder{}).ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("unknown operations route = %d, want 404", rr.Code)
	}
}

func TestOperationsReliabilityRoutesPrivateScopedGETOnlyNoStore(t *testing.T) {
	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	endpoint := true
	freshness := 25.0
	latency := 50.0
	store := &fakePublicationStore{
		discovery: validationHealthTestDiscovery(now),
		reliabilityHealth: []compliance.ReliabilityFeedHealthRecord{
			{FeedType: "vehicle_positions", SnapshotAt: now, EndpointAvailable: &endpoint, FreshnessSeconds: &freshness, GenerationLatencyMS: &latency},
		},
		reliabilityIncidents: compliance.NormalizeReliabilityIncidentRollup(now, 0, nil, nil, nil, nil, nil, 10),
	}
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
			Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
		}})
		for _, path := range []string{"/admin/operations/reliability", "/admin/operations/reliability.json"} {
			req := httptest.NewRequest(http.MethodGet, path+"?agency_id=demo-agency", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("%s role %s status = %d, want 200: %s", path, role, rr.Code, rr.Body.String())
			}
			if got := rr.Header().Get("Cache-Control"); got != "no-store" {
				t.Fatalf("%s Cache-Control = %q, want no-store", path, got)
			}
		}
	}
	srv := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator}, Method: auth.MethodBearer,
	}})
	for _, path := range []string{"/admin/operations/reliability", "/admin/operations/reliability.json"} {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		rr := httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("POST %s status = %d, want 405", path, rr.Code)
		}
		req = httptest.NewRequest(http.MethodGet, path+"?agency_id=other-agency", nil)
		rr = httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("conflict %s status = %d, want 403", path, rr.Code)
		}
	}
	unauth := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, authRejectAll{})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/reliability.json", nil)
	rr := httptest.NewRecorder()
	unauth.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK {
		t.Fatalf("unauthenticated reliability JSON returned 200")
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/reliability.json", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public reliability route status = %d, want 404", rr.Code)
	}
}

func TestOperationsReliabilityJSONShapeOrderMissingFlagsAndNoLeakage(t *testing.T) {
	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	endpoint := true
	freshness := 25.0
	latency := 50.0
	oldest := now.Add(-time.Hour)
	store := &fakePublicationStore{
		discovery: validationHealthTestDiscovery(now),
		reliabilityHealth: []compliance.ReliabilityFeedHealthRecord{
			{FeedType: "vehicle_positions", SnapshotAt: now, EndpointAvailable: &endpoint, FreshnessSeconds: &freshness, GenerationLatencyMS: &latency},
		},
		reliabilityIncidents: compliance.NormalizeReliabilityIncidentRollup(now, 2,
			map[string]int{"open": 1, "resolved": 1},
			map[string]int{"warning": 2},
			map[string]int{"prediction_review": 2},
			&oldest,
			[]compliance.ReliabilityIncidentItem{{
				ID: 1, Type: "prediction_review", Severity: "warning", Status: "open", OpenedAt: now.Add(-time.Hour),
				Title: "raw details_json token https://private.example", Category: "payload_json",
			}},
			10),
	}
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/reliability.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var summary compliance.ReliabilitySummary
	if err := json.Unmarshal(rr.Body.Bytes(), &summary); err != nil {
		t.Fatalf("decode summary: %v", err)
	}
	if len(summary.Feeds) != 4 || summary.Feeds[0].FeedType != "schedule" || summary.Feeds[1].FeedType != "vehicle_positions" || summary.Feeds[2].FeedType != "trip_updates" || summary.Feeds[3].FeedType != "alerts" {
		t.Fatalf("feed order/shape = %+v", summary.Feeds)
	}
	if summary.Feeds[0].Status == compliance.ReliabilityStatusOK || summary.Feeds[2].Status == compliance.ReliabilityStatusOK || summary.Feeds[3].Status == compliance.ReliabilityStatusOK {
		t.Fatalf("missing feed data became ok: %+v", summary.Feeds)
	}
	if summary.ClaimFlags.ExternalEvidenceCreated || summary.ClaimFlags.FinalRootEvidenceCreated || summary.ClaimFlags.ConsumerStatusesChanged || summary.ClaimFlags.ComplianceClaimed || summary.ClaimFlags.ProductionReadinessClaimed || summary.ClaimFlags.SLAClaimed || summary.ClaimFlags.UptimeGuaranteeClaimed || summary.ClaimFlags.HostedSaaSClaimed || summary.ClaimFlags.AgencyAdoptionClaimed || summary.ClaimFlags.ConsumerAcceptanceClaimed || summary.ClaimFlags.VendorCompatibilityClaimed || summary.ClaimFlags.ProductionGradeETAClaimed {
		t.Fatalf("claim flags not all false: %+v", summary.ClaimFlags)
	}
	body := rr.Body.String()
	for _, forbidden := range []string{"details_json", "token", "https://private.example", "payload_json", "raw payload", "postgres://", "Authorization", "Cookie"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("reliability JSON leaked %q: %s", forbidden, body)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/reliability", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	html := rr.Body.String()
	if !strings.Contains(html, ">3600<") {
		t.Fatalf("html oldest-open age did not render numeric seconds: %s", html)
	}
	if strings.Contains(html, "0x") {
		t.Fatalf("html appears to contain pointer formatting: %s", html)
	}
	for _, forbidden := range []string{"TOKEN=secret", "https://private.example", "payload_json", "postgres://", "Authorization", "Cookie"} {
		if strings.Contains(html, forbidden) {
			t.Fatalf("reliability HTML leaked %q: %s", forbidden, html)
		}
	}
}

func TestValidationHealthJSONContractOrderAndNoLeakage(t *testing.T) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{discovery: validationHealthTestDiscovery(now), validationRecords: []compliance.ValidationReportRecord{
		validationHealthRecord(1, "trip_updates", "feed-v1", "failed", now.Add(time.Minute)),
		validationHealthRecord(2, "schedule", "feed-v1", "passed", now.Add(2*time.Minute)),
		validationHealthRecord(3, "vehicle_positions", "feed-v1", "warning", now.Add(3*time.Minute)),
		validationHealthRecord(4, "alerts", "feed-v1", "passed", now.Add(4*time.Minute)),
	}}
	handler := newValidationHealthTestHandler(t, auth.RoleReadOnly, store, fakeScheduleBuilder{})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/validation-health.json?agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d: %s", rr.Code, rr.Body.String())
	}
	var summary compliance.ValidationHealthSummary
	if err := json.Unmarshal(rr.Body.Bytes(), &summary); err != nil {
		t.Fatalf("decode summary: %v", err)
	}
	assertValidationHealthSummaryShape(t, summary)
	if summary.Feeds[0].FeedType != "schedule" || summary.Feeds[1].FeedType != "vehicle_positions" || summary.Feeds[2].FeedType != "trip_updates" || summary.Feeds[3].FeedType != "alerts" {
		t.Fatalf("unexpected feed order: %+v", summary.Feeds)
	}
	assertValidationHealthHTTPNoLeakage(t, rr.Body.String())
	assertValidationHealthJSONAllowlist(t, rr.Body.Bytes())
}

func TestValidationHealthHTMLMatchesJSONRows(t *testing.T) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{discovery: validationHealthTestDiscovery(now), validationRecords: []compliance.ValidationReportRecord{
		validationHealthRecord(1, "schedule", "feed-v1", "passed", now),
		validationHealthRecord(2, "vehicle_positions", "feed-v1", "passed", now),
		validationHealthRecord(3, "trip_updates", "feed-v1", "passed", now),
		validationHealthRecord(4, "alerts", "feed-v1", "passed", now),
	}}
	handler := newValidationHealthTestHandler(t, auth.RoleAdmin, store, fakeScheduleBuilder{})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/validation-health.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	var summary compliance.ValidationHealthSummary
	if err := json.Unmarshal(rr.Body.Bytes(), &summary); err != nil {
		t.Fatal(err)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/validation-health", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	body := rr.Body.String()
	if !strings.Contains(body, "<td>"+summary.OverallStatus+"</td>") {
		t.Fatalf("html overall status does not match json %q", summary.OverallStatus)
	}
	for _, row := range summary.Feeds {
		if !strings.Contains(body, "<td>"+row.FeedType+"</td>") {
			t.Fatalf("html missing feed row %q", row.FeedType)
		}
	}
	for _, want := range []string{"external_evidence_created", "consumer_statuses_changed", "compliance_claimed", "production_readiness_claimed", "private diagnostics"} {
		if !strings.Contains(body, want) {
			t.Fatalf("html missing %q", want)
		}
	}
	for _, forbidden := range []string{"compliance gate", "production gate"} {
		if strings.Contains(strings.ToLower(body), forbidden) {
			t.Fatalf("html contains forbidden wording %q", forbidden)
		}
	}
}

func TestValidationHealthPostStrictnessCSRFAndBodyCap(t *testing.T) {
	store := &fakePublicationStore{discovery: validationHealthTestDiscovery(time.Now().UTC())}
	srv := newValidationHealthTestHandler(t, auth.RoleAdmin, store, fakeScheduleBuilder{})
	for _, field := range []string{"feed_type", "validator_id", "validator_path", "validator_command", "output_path", "artifact_path", "report_path", "url", "argv", "args", "timeout_seconds"} {
		req := httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health", strings.NewReader(validationHealthRunAllForm()+"&"+field+"=browser"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("field %s status = %d, want 400", field, rr.Code)
		}
		if got := rr.Header().Get("Cache-Control"); got != "no-store" {
			t.Fatalf("field %s Cache-Control = %q, want no-store", field, got)
		}
	}
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health", strings.NewReader(strings.Repeat("x", validationHealthPostMaxBytes+1)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("large form status = %d, want 413", rr.Code)
	}
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health", strings.NewReader("action=run_all"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("bearer missing csrf status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	cookieAuth := auth.TestAuthenticator{Principal: auth.Principal{Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodCookie}}
	srv = newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}, schedule: fakeScheduleBuilder{}, csrfSecret: "test-csrf"}, cookieAuth)
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health", strings.NewReader("action=run_all"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("missing csrf status = %d, want 403", rr.Code)
	}
}

func TestValidationHealthPostRunAllPartialSuccessAndBounded(t *testing.T) {
	t.Setenv("GTFS_VALIDATOR_PATH", writeScheduleValidator(t))
	t.Setenv("GTFS_RT_VALIDATOR_PATH", writeRealtimeValidator(t))
	t.Setenv("GTFS_RT_VALIDATOR_VERSION", "test-validator")
	t.Setenv("GTFS_RT_VALIDATOR_ARGS", "")
	now := time.Now().UTC()
	store := &fakePublicationStore{discovery: validationHealthTestDiscovery(now)}
	artifacts := &fakeRealtimeArtifacts{payloads: map[string][]byte{
		"vehicle_positions": []byte("vp"),
		"alerts":            []byte("alerts"),
	}, errors: map[string]error{"trip_updates": errors.New("private /tmp/path TOKEN=secret")}}
	handler := newOperationsTestHandler(&handler{
		store:    store,
		devices:  fakeDeviceStore{},
		schedule: fakeScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "demo-agency", FeedVersionID: "feed-v1", RevisionTime: now, Payload: []byte("schedule zip")}},
		realtime: artifacts,
	}, auth.TestAuthenticator{Principal: auth.Principal{Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer}})
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health", strings.NewReader(validationHealthRunAllForm()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.storeValidationCalls != 3 {
		t.Fatalf("stored validation calls = %d, want 3", store.storeValidationCalls)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "trip_updates") || !strings.Contains(body, compliance.ValidationHealthArtifactUnavailable) {
		t.Fatalf("partial failure not represented safely: %s", body)
	}
	assertValidationHealthHTTPNoLeakage(t, body)
	if len(body) > 50000 {
		t.Fatalf("html size = %d, want bounded", len(body))
	}
}

func TestValidationHealthConcurrentRunAllDoesNotPanicOrLeak(t *testing.T) {
	store := &fakePublicationStore{discovery: validationHealthTestDiscovery(time.Now().UTC())}
	handler := newValidationHealthTestHandler(t, auth.RoleAdmin, store, fakeScheduleBuilder{snapshot: schedule.Snapshot{AgencyID: "demo-agency", FeedVersionID: "feed-v1", RevisionTime: time.Now().UTC(), Payload: []byte("schedule zip")}})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health", strings.NewReader(validationHealthRunAllForm()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Errorf("status = %d, want 200: %s", rr.Code, rr.Body.String())
			}
			assertValidationHealthHTTPNoLeakage(t, rr.Body.String())
		}()
	}
	wg.Wait()
}

func TestValidationHealthLargeHistoryPageAndJSONBounded(t *testing.T) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	records := make([]compliance.ValidationReportRecord, 0, 10000)
	for i := 0; i < 10000; i++ {
		feedType := []string{"schedule", "vehicle_positions", "trip_updates", "alerts"}[i%4]
		record := validationHealthRecord(int64(i+1), feedType, "feed-v1", "warning", now.Add(time.Duration(i)*time.Second))
		record.Result.Report = map[string]any{"raw_report": strings.Repeat("stdout stderr argv /tmp/private TOKEN=SECRET postgres://user:pass@localhost/db", 100)}
		records = append(records, record)
	}
	store := &fakePublicationStore{discovery: validationHealthTestDiscovery(now), validationRecords: records}
	handler := newValidationHealthTestHandler(t, auth.RoleAdmin, store, fakeScheduleBuilder{})
	for _, path := range []string{"/admin/operations/validation-health.json", "/admin/operations/validation-health"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d: %s", path, rr.Code, rr.Body.String())
		}
		if len(rr.Body.Bytes()) > 50000 {
			t.Fatalf("%s output size = %d, want bounded", path, len(rr.Body.Bytes()))
		}
		assertValidationHealthHTTPNoLeakage(t, rr.Body.String())
	}
}

func BenchmarkRenderValidationHealthPage(b *testing.B) {
	store := &fakePublicationStore{discovery: validationHealthTestDiscovery(time.Now().UTC()), validationRecords: []compliance.ValidationReportRecord{
		validationHealthRecord(1, "schedule", "feed-v1", "passed", time.Now().UTC()),
		validationHealthRecord(2, "vehicle_positions", "feed-v1", "passed", time.Now().UTC()),
		validationHealthRecord(3, "trip_updates", "feed-v1", "warning", time.Now().UTC()),
		validationHealthRecord(4, "alerts", "feed-v1", "passed", time.Now().UTC()),
	}}
	handler := newValidationHealthTestHandler(b, auth.RoleAdmin, store, fakeScheduleBuilder{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/admin/operations/validation-health", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			b.Fatalf("status = %d", rr.Code)
		}
	}
}

func BenchmarkRenderValidationHealthJSON(b *testing.B) {
	store := &fakePublicationStore{discovery: validationHealthTestDiscovery(time.Now().UTC()), validationRecords: []compliance.ValidationReportRecord{
		validationHealthRecord(1, "schedule", "feed-v1", "passed", time.Now().UTC()),
		validationHealthRecord(2, "vehicle_positions", "feed-v1", "passed", time.Now().UTC()),
		validationHealthRecord(3, "trip_updates", "feed-v1", "warning", time.Now().UTC()),
		validationHealthRecord(4, "alerts", "feed-v1", "passed", time.Now().UTC()),
	}}
	handler := newValidationHealthTestHandler(b, auth.RoleAdmin, store, fakeScheduleBuilder{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/admin/operations/validation-health.json", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			b.Fatalf("status = %d", rr.Code)
		}
	}
}

func TestOperationsConsoleRendersSafeTripUpdatesQualitySummary(t *testing.T) {
	now := time.Date(2026, 4, 26, 12, 0, 0, 0, time.UTC)
	coverage := 50.0
	store := &fakePublicationStore{
		tripDiagnostics: compliance.TripUpdatesDiagnosticsSummary{
			Recorded:                      true,
			SnapshotAt:                    now,
			AdapterName:                   "deterministic",
			DiagnosticsStatus:             prediction.StatusOK,
			DiagnosticsReason:             prediction.ReasonPartialPredictions,
			ActiveFeedVersionID:           "feed-demo",
			DiagnosticsPersistenceOutcome: "stored",
			Metrics: prediction.Metrics{
				TelemetryRowsConsidered:      2,
				AssignmentsConsidered:        2,
				EligiblePredictionCandidates: 2,
				TripUpdatesEmitted:           1,
				UnknownAssignments:           1,
				AmbiguousAssignments:         1,
				StaleTelemetryRows:           1,
				ManualOverrideAssignments:    1,
				WithheldByReason:             map[string]int{prediction.ReasonDegradedAssignment: 1},
				UnknownAssignmentRate:        prediction.RateMetric{Numerator: 1, Denominator: 2, Percent: &coverage, Status: "measured", DenominatorDefinition: "current unknown assignments / current assignments considered"},
				AmbiguousAssignmentRate:      prediction.RateMetric{Numerator: 1, Denominator: 2, Percent: &coverage, Status: "measured", DenominatorDefinition: "current ambiguous assignments / current assignments considered"},
				StaleTelemetryRate:           prediction.RateMetric{Numerator: 1, Denominator: 2, Percent: &coverage, Status: "measured", DenominatorDefinition: "stale latest telemetry rows / telemetry rows considered"},
				TripUpdatesCoverageRate:      prediction.RateMetric{Numerator: 1, Denominator: 2, Percent: &coverage, Status: "measured", DenominatorDefinition: "emitted non-canceled Trip Updates / eligible in-service ETA candidates"},
				FutureStopCoverageRate:       prediction.RateMetric{Numerator: 1, Denominator: 2, Percent: &coverage, Status: "measured", DenominatorDefinition: "non-canceled Trip Updates with at least one future stop update / eligible in-service ETA candidates"},
			},
		},
	}
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/feeds", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Trip Updates Quality Diagnostics", "deterministic", "partial_predictions", "50.0% (1/2)", "degraded_assignment: 1"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"payload_json", "score_details", "private_debug", "token_hash"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("body leaks %q: %s", forbidden, body)
		}
	}
}

func TestOperationsConsoleRendersDemoStateWithSafeTelemetryDiagnostics(t *testing.T) {
	t.Setenv("PUBLICATION_ENVIRONMENT", "pilot")
	now := time.Date(2026, 4, 26, 12, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{
		discovery: compliance.FeedDiscovery{
			AgencyID: "demo-agency", AgencyName: "Demo Agency", GeneratedAt: now, PublicationEnvironment: "dev",
			PublicBaseURL: "http://localhost:8080",
			Feeds: []compliance.FeedMetadata{{
				FeedType: "schedule", CanonicalPublicURL: "http://localhost:8080/public/gtfs/schedule.zip",
				ActivationStatus: "active", ActiveFeedVersionID: "gtfs-import-3", LastValidationStatus: "passed", LastValidationAt: &now,
			}},
			Readiness: compliance.Readiness{AllRequiredFeedsListed: true, LicenseComplete: true, ContactComplete: true, CanonicalValidationComplete: true},
		},
		scorecard: compliance.Scorecard{AgencyID: "demo-agency", SnapshotAt: now, OverallStatus: compliance.StatusYellow},
		consumers: []compliance.ConsumerRecord{{ConsumerName: "Google Maps", Status: "not_started", UpdatedAt: now}},
	}
	handler := newOperationsTestHandler(&handler{
		store: store,
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{{
			AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now, CreatedAt: now,
		}}},
		telemetry: fakeTelemetryRepository{latest: []telemetry.StoredEvent{{
			ID: 42,
			Event: telemetry.Event{
				AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Timestamp: now.Add(-30 * time.Second), Lat: 1, Lon: 2,
			},
			ReceivedAt: now.Add(-29 * time.Second), IngestStatus: telemetry.IngestStatusAccepted, PayloadJSON: []byte(`{"secret":"hidden"}`),
		}}},
		state: fakeStateRepository{assignments: map[string]state.Assignment{"bus-1": {
			VehicleID: "bus-1", State: state.StateInService, RouteID: "route-1", TripID: "trip-1", Confidence: 0.91,
			ReasonCodes: []string{state.ReasonTripHintMatch}, DegradedState: state.DegradedNone, AssignmentSource: state.AssignmentSourceAutomatic,
			ScoreDetails: map[string]any{"private_debug": true}, ActiveFrom: now.Add(-25 * time.Second),
		}}},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/telemetry", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"bus-1", "trip-1", "route-1", "0.91", state.ReasonTripHintMatch} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"payload_json", "secret", "score_details", "private_debug"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("body leaks %q: %s", forbidden, body)
		}
	}
}

func TestOperationsConsoleViewsAreAgencyScoped(t *testing.T) {
	now := time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)
	handler := newOperationsTestHandler(&handler{
		store: &fakePublicationStore{
			discovery: compliance.FeedDiscovery{
				AgencyID: "agency-a", AgencyName: "Agency A", GeneratedAt: now,
				PublicBaseURL: "https://agency-a.example",
				Feeds:         []compliance.FeedMetadata{{FeedType: "schedule", CanonicalPublicURL: "https://agency-a.example/public/gtfs/schedule.zip"}},
			},
			scorecard: compliance.Scorecard{AgencyID: "agency-a", OverallStatus: compliance.StatusYellow},
			consumers: []compliance.ConsumerRecord{{ConsumerName: "Consumer A", Status: "not_started", UpdatedAt: now}},
		},
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{
			{AgencyID: "agency-a", DeviceID: "device-a-1", VehicleID: "bus-a-1", Status: "active", ValidFrom: now, CreatedAt: now},
			{AgencyID: "agency-b", DeviceID: "device-b-1", VehicleID: "bus-b-1", Status: "active", ValidFrom: now, CreatedAt: now},
		}},
		telemetry: fakeTelemetryRepository{latest: []telemetry.StoredEvent{
			{Event: telemetry.Event{AgencyID: "agency-a", DeviceID: "device-a-1", VehicleID: "bus-a-1", Timestamp: now, Lat: 1, Lon: 2}, ReceivedAt: now, IngestStatus: telemetry.IngestStatusAccepted},
			{Event: telemetry.Event{AgencyID: "agency-b", DeviceID: "device-b-1", VehicleID: "bus-b-1", Timestamp: now, Lat: 3, Lon: 4}, ReceivedAt: now, IngestStatus: telemetry.IngestStatusAccepted},
		}},
		state: fakeStateRepository{assignments: map[string]state.Assignment{
			"bus-a-1": {VehicleID: "bus-a-1", State: state.StateInService, RouteID: "route-a-10", TripID: "trip-a-10", Confidence: 0.9, ActiveFrom: now},
		}},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader-a@example.com", AgencyID: "agency-a", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})

	for _, section := range []string{"", "/readiness", "/feeds", "/telemetry", "/devices", "/consumers", "/evidence", "/setup"} {
		req := httptest.NewRequest(http.MethodGet, "/admin/operations"+section+"?agency_id=agency-b", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("section %q status = %d, want 403", section, rr.Code)
		}
	}

	for _, path := range []string{"/admin/operations/telemetry", "/admin/operations/devices", "/admin/operations/consumers"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
		}
		body := rr.Body.String()
		if strings.Contains(body, "bus-b-1") || strings.Contains(body, "device-b-1") || strings.Contains(body, "Consumer B") {
			t.Fatalf("%s leaked agency-b data: %s", path, body)
		}
	}
}

func TestOperationsDeviceRebindShowsTokenOnlyOnPost(t *testing.T) {
	deviceStore := &fakeDeviceStoreWithToken{token: "one-time-token"}
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: deviceStore}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	form := strings.NewReader("agency_id=demo-agency&device_id=device-1&vehicle_id=bus-1&reason=rotate")
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/devices", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if !strings.Contains(body, "one-time-token") {
		t.Fatalf("POST body does not show one-time token: %s", body)
	}
	if strings.Contains(body, "token_hash") {
		t.Fatalf("POST body leaks token hash: %s", body)
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/devices", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if strings.Contains(rr.Body.String(), "one-time-token") {
		t.Fatalf("GET body repeats one-time token: %s", rr.Body.String())
	}
}

func TestOperationsDeviceRebindRequiresAdminRole(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/devices", strings.NewReader("agency_id=demo-agency&device_id=device-1&vehicle_id=bus-1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rr.Code)
	}
}

func TestOperationsCookiePostRequiresCSRF(t *testing.T) {
	cfg := auth.JWTConfig{Secrets: []string{"test-secret"}, Issuer: "test-issuer", Audience: "test-audience", TTL: time.Hour}
	signer, err := auth.NewSigner(cfg)
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	token, _, err := signer.Sign("admin@example.com", "demo-agency", time.Hour)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	verifier, err := auth.NewVerifier(cfg)
	if err != nil {
		t.Fatalf("verifier: %v", err)
	}
	middleware := auth.NewMiddleware(verifier, auth.StaticRoleStore{Roles: []auth.Role{auth.RoleAdmin}}, "csrf-secret")
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}, csrfSecret: "csrf-secret"}, middleware)
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/devices", strings.NewReader("agency_id=demo-agency&device_id=device-1&vehicle_id=bus-1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "admin_session", Value: token})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 for missing CSRF", rr.Code)
	}
}

func TestOperationsConsumersDoNotInventAcceptanceClaims(t *testing.T) {
	now := time.Date(2026, 4, 26, 12, 0, 0, 0, time.UTC)
	handler := newOperationsTestHandler(&handler{
		store:   &fakePublicationStore{consumers: []compliance.ConsumerRecord{{ConsumerName: "Google Maps", Status: "not_started", UpdatedAt: now}}},
		devices: fakeDeviceStore{},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/consumers", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Google Maps", "not_started", "Mobility Database", "transit.land", "docs/evidence tracker"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
	if strings.Contains(strings.ToLower(body), "accepted by") {
		t.Fatalf("body invents acceptance claim: %s", body)
	}
}

func assertChecklistShape(t *testing.T, checklist operatorChecklistView) {
	t.Helper()
	if checklist.AgencyID == "" || checklist.Groups == nil || checklist.Counts.Groups != len(checklist.Groups) {
		t.Fatalf("invalid checklist top-level shape: %+v", checklist)
	}
	allowedStatuses := map[string]bool{"ok": true, "needs_review": true, "missing": true, "blocked": true, "unknown": true}
	seenIDs := map[string]bool{}
	rowCount := 0
	for _, group := range checklist.Groups {
		if group.ID == "" || group.Label == "" || len(group.Rows) == 0 {
			t.Fatalf("invalid group shape: %+v", group)
		}
		for _, row := range group.Rows {
			rowCount++
			if row.ID == "" || row.Label == "" || row.Source == "" || row.CurrentSignal == "" || row.NextAction == "" || row.ClaimBoundary == "" || row.DocsLinks == nil || row.HeuristicLabels == nil {
				t.Fatalf("invalid row shape: %+v", row)
			}
			if seenIDs[row.ID] {
				t.Fatalf("duplicate row id %q", row.ID)
			}
			seenIDs[row.ID] = true
			if !allowedStatuses[row.Status] {
				t.Fatalf("row %q status = %q, want neutral status", row.ID, row.Status)
			}
		}
	}
	if checklist.Counts.Rows != rowCount {
		t.Fatalf("counts rows = %d, want %d", checklist.Counts.Rows, rowCount)
	}
}

func nonCommentCaddyLines(text string) []string {
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		lines = append(lines, trimmed)
	}
	return lines
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func assertChecklistFlagsFalse(t *testing.T, flags operatorChecklistFlags) {
	t.Helper()
	if flags.ExternalEvidenceCreated || flags.FinalRootEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.AgencyApprovalClaimed || flags.ConsumerAcceptanceClaimed {
		t.Fatalf("checklist flags must all be false: %+v", flags)
	}
}

func assertLaunchpadShape(t *testing.T, launchpad agencyLaunchpadView) {
	t.Helper()
	if launchpad.AgencyID == "" || launchpad.Boundary == "" || len(launchpad.Sections) == 0 || launchpad.Counts.Sections != len(launchpad.Sections) {
		t.Fatalf("invalid launchpad top-level shape: %+v", launchpad)
	}
	allowedStatuses := map[string]bool{"ok": true, "needs_review": true, "missing": true, "blocked": true, "unknown": true}
	seenIDs := map[string]bool{}
	for _, section := range launchpad.Sections {
		if section.ID == "" || section.Label == "" || section.CurrentSignal == "" || section.ClaimBoundary == "" || len(section.NextActions) == 0 || section.DocsLinks == nil || section.CommandSuggestions == nil || section.AdminLinks == nil {
			t.Fatalf("invalid launchpad section shape: %+v", section)
		}
		if seenIDs[section.ID] {
			t.Fatalf("duplicate launchpad section id %q", section.ID)
		}
		seenIDs[section.ID] = true
		if !allowedStatuses[section.Status] {
			t.Fatalf("section %q status = %q, want neutral status", section.ID, section.Status)
		}
		for _, link := range section.DocsLinks {
			if !strings.HasPrefix(link, "docs/") {
				t.Fatalf("section %s has unsafe docs link %q", section.ID, link)
			}
		}
	}
}

func assertLaunchpadFlagsFalse(t *testing.T, flags agencyLaunchpadClaimFlags) {
	t.Helper()
	if flags.ExternalEvidenceCreated || flags.FinalRootEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.AgencyApprovalClaimed || flags.ConsumerAcceptanceClaimed || flags.PublicLaunchClaimed || flags.HostedSaaSClaimed || flags.VendorCompatibilityClaimed || flags.ProductionGradeETAClaimed {
		t.Fatalf("launchpad flags must all be false: %+v", flags)
	}
}

func assertLaunchpadSafeStrings(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"raw-token-value", "authorization:", "set-cookie", ".cache", "database_url", "restore_database_url", "payload_json", "raw telemetry", "token_hash", "file://", "/users/", "/opt/open-transit-rt", "/var/lib", "/etc/", "postgres://"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("launchpad leaks forbidden private string %q: %s", forbidden, body)
		}
	}
	for _, forbidden := range []string{"agency_approved", "final_root_approved", "consumer_ready", "production_ready", "public_launch_complete"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("launchpad emits forbidden label %q: %s", forbidden, body)
		}
	}
}

func launchpadSectionStatus(launchpad agencyLaunchpadView, id string) string {
	for _, section := range launchpad.Sections {
		if section.ID == id {
			return section.Status
		}
	}
	return ""
}

func assertChecklistSafeStrings(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"raw-token-value", "authorization:", "set-cookie", ".cache", "database_url", "restore_database_url", "payload_json", "raw telemetry", "token_hash", "file://", "/users/", "/opt/open-transit-rt", "/var/lib", "/etc/"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("checklist leaks forbidden private string %q: %s", forbidden, body)
		}
	}
	for _, forbidden := range []string{"agency_approved", "final_root_approved", "compliant", "accepted", "consumer_ready", "production_ready"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("checklist emits forbidden label %q: %s", forbidden, body)
		}
	}
}

func assertChecklistNoPositiveClaims(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"evidence packet created", "compliance evidence", "final-root evidence created", "agency approved", "consumer accepted", "production ready"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("checklist emits positive claim %q: %s", forbidden, body)
		}
	}
}

func assertChecklistDocsLinksSafe(t *testing.T, checklist operatorChecklistView) {
	t.Helper()
	for _, group := range checklist.Groups {
		for _, row := range group.Rows {
			for _, link := range row.DocsLinks {
				if !strings.HasPrefix(link, "docs/") {
					t.Fatalf("row %s has non-repo-relative docs link %q", row.ID, link)
				}
				lower := strings.ToLower(link)
				for _, forbidden := range []string{".cache", "file://", "/users", "localhost", "/opt/open-transit-rt", "/var/lib", "/etc"} {
					if strings.Contains(lower, forbidden) {
						t.Fatalf("row %s docs link %q contains private path marker %q", row.ID, link, forbidden)
					}
				}
			}
		}
	}
}

func allHeuristicLabels(checklist operatorChecklistView) []string {
	var labels []string
	for _, group := range checklist.Groups {
		for _, row := range group.Rows {
			labels = append(labels, row.HeuristicLabels...)
		}
	}
	return labels
}

func checklistContainsSignal(checklist operatorChecklistView, signal string) bool {
	for _, group := range checklist.Groups {
		for _, row := range group.Rows {
			if strings.Contains(row.CurrentSignal, signal) {
				return true
			}
		}
	}
	return false
}

func TestAgencyConfigReadyzRequiresDBActiveFeedAndPublicationMetadata(t *testing.T) {
	cases := []struct {
		name     string
		pinger   fakePinger
		schedule fakeScheduleBuilder
		store    *fakePublicationStore
		want     int
	}{
		{
			name:   "database unavailable",
			pinger: fakePinger{err: errors.New("down")},
			store:  &fakePublicationStore{discovery: readyDiscovery()},
			want:   http.StatusServiceUnavailable,
		},
		{
			name:     "active schedule feed missing",
			schedule: fakeScheduleBuilder{readyErr: errors.New("no active feed")},
			store:    &fakePublicationStore{discovery: readyDiscovery()},
			want:     http.StatusServiceUnavailable,
		},
		{
			name:  "publication config missing",
			store: &fakePublicationStore{discoveryErr: errors.New("no feed_config")},
			want:  http.StatusServiceUnavailable,
		},
		{
			name:  "published feed metadata incomplete",
			store: &fakePublicationStore{discovery: compliance.FeedDiscovery{Readiness: compliance.Readiness{AllRequiredFeedsListed: false}}},
			want:  http.StatusServiceUnavailable,
		},
		{
			name:  "ready",
			store: &fakePublicationStore{discovery: readyDiscovery()},
			want:  http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.store == nil {
				tc.store = &fakePublicationStore{discovery: readyDiscovery()}
			}
			handler := newHandlerWithRealtime(
				"demo-agency",
				tc.schedule,
				tc.store,
				fakeDeviceStore{},
				tc.pinger,
				auth.TestAuthenticator{Principal: auth.Principal{Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer}},
				&fakeRealtimeArtifacts{},
			)
			req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != tc.want {
				t.Fatalf("status = %d, want %d: %s", rr.Code, tc.want, rr.Body.String())
			}
		})
	}
}

func writeRealtimeValidator(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "gtfs-rt-validator.sh")
	script := `#!/bin/sh
realtime=""
schedule=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --realtime) shift; realtime="$1" ;;
    --schedule) shift; schedule="$1" ;;
  esac
  shift
done
test -s "$schedule" || exit 3
test -s "$realtime" || exit 4
printf '%s' '{"status":"passed","error_count":0,"warning_count":0,"info_count":1}'
`
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
		t.Fatalf("write validator: %v", err)
	}
	return path
}

func writeScheduleValidator(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "gtfs-validator.sh")
	script := `#!/bin/sh
schedule=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    -i) shift; schedule="$1" ;;
  esac
  shift
done
test -s "$schedule" || exit 3
printf '%s' '{"status":"passed","error_count":0,"warning_count":0,"info_count":1}'
`
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
		t.Fatalf("write validator: %v", err)
	}
	return path
}

func writeScheduleValidatorWithReport(t testing.TB, report string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "gtfs-validator.sh")
	script := `#!/bin/sh
schedule=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    -i) shift; schedule="$1" ;;
  esac
  shift
done
test -s "$schedule" || exit 3
printf '%s' '` + report + `'
`
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
		t.Fatalf("write validator: %v", err)
	}
	return path
}

func newGTFSQualityTestHandler(t testing.TB, role auth.Role, store *fakePublicationStore, scheduleBuilder scheduleBuilder) http.Handler {
	t.Helper()
	if store == nil {
		store = &fakePublicationStore{discovery: gtfsQualityDiscovery(time.Now().UTC())}
	}
	return newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}, schedule: scheduleBuilder, csrfSecret: "test-csrf"}, auth.TestAuthenticator{Principal: auth.Principal{Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer}})
}

func gtfsQualityDiscovery(now time.Time) compliance.FeedDiscovery {
	return compliance.FeedDiscovery{
		AgencyID: "demo-agency",
		Feeds: []compliance.FeedMetadata{{
			FeedType:            "schedule",
			CanonicalPublicURL:  "https://feeds.example.org/public/gtfs/schedule.zip",
			ActivationStatus:    "active",
			ActiveFeedVersionID: "feed-v1",
			RevisionTimestamp:   &now,
		}},
	}
}

func noticeMap(code string, severity string) map[string]any {
	return map[string]any{"code": code, "severity": severity, "message": code}
}

func gtfsQualityRerunForm() string {
	principal := auth.Principal{Subject: "user@example.com", AgencyID: "demo-agency"}
	return "action=rerun_static_validator&csrf_token=" + auth.CSRFToken("test-csrf", principal)
}

func newValidationHealthTestHandler(t testing.TB, role auth.Role, store *fakePublicationStore, scheduleBuilder scheduleBuilder) http.Handler {
	t.Helper()
	if store == nil {
		store = &fakePublicationStore{discovery: validationHealthTestDiscovery(time.Now().UTC())}
	}
	return newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}, schedule: scheduleBuilder, realtime: &fakeRealtimeArtifacts{}, csrfSecret: "test-csrf"}, auth.TestAuthenticator{Principal: auth.Principal{Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer}})
}

func validationHealthRunAllForm() string {
	principal := auth.Principal{Subject: "user@example.com", AgencyID: "demo-agency"}
	return "action=run_all&csrf_token=" + auth.CSRFToken("test-csrf", principal)
}

func validationHealthTestDiscovery(now time.Time) compliance.FeedDiscovery {
	return compliance.FeedDiscovery{AgencyID: "demo-agency", Feeds: []compliance.FeedMetadata{
		{FeedType: "schedule", CanonicalPublicURL: "https://feeds.example.org/public/gtfs/schedule.zip", ActivationStatus: "active", ActiveFeedVersionID: "feed-v1", RevisionTimestamp: &now},
		{FeedType: "vehicle_positions", CanonicalPublicURL: "https://feeds.example.org/public/gtfsrt/vehicle_positions.pb", ActivationStatus: "active", ActiveFeedVersionID: "feed-v1", RevisionTimestamp: &now},
		{FeedType: "trip_updates", CanonicalPublicURL: "https://feeds.example.org/public/gtfsrt/trip_updates.pb", ActivationStatus: "active", ActiveFeedVersionID: "feed-v1", RevisionTimestamp: &now},
		{FeedType: "alerts", CanonicalPublicURL: "https://feeds.example.org/public/gtfsrt/alerts.pb", ActivationStatus: "active", ActiveFeedVersionID: "feed-v1", RevisionTimestamp: &now},
	}}
}

func validationHealthRecord(id int64, feedType, feedVersionID, status string, createdAt time.Time) compliance.ValidationReportRecord {
	resultStatus := status
	warnings := 0
	errors := 0
	if status == "warning" {
		warnings = 1
	}
	if status == "failed" {
		errors = 1
	}
	return compliance.ValidationReportRecord{
		ID:        id,
		CreatedAt: createdAt,
		Result: compliance.ValidationResult{
			AgencyID:      "demo-agency",
			FeedType:      feedType,
			FeedVersionID: feedVersionID,
			ValidatorName: compliance.ValidatorNameForHealthID(compliance.ValidatorIDForHealthFeed(feedType)),
			Status:        resultStatus,
			ErrorCount:    errors,
			WarningCount:  warnings,
			Report:        map[string]any{"raw_report": "secret", "stdout": "secret", "stderr": "secret", "argv": []any{"/tmp/private"}},
		},
	}
}

func assertValidationHealthSummaryShape(t *testing.T, summary compliance.ValidationHealthSummary) {
	t.Helper()
	if len(summary.Feeds) != 4 || summary.GeneratedAt.IsZero() {
		t.Fatalf("invalid summary shape: %+v", summary)
	}
	if summary.ExternalEvidenceCreated || summary.ConsumerStatusesChanged || summary.ComplianceClaimed || summary.ProductionReadinessClaimed {
		t.Fatalf("false flags must stay false: %+v", summary)
	}
	for _, row := range summary.Feeds {
		if row.FeedType == "" || row.ValidatorID == "" || row.ValidatorName == "" || row.HealthStatus == "" || row.NextAction == "" || row.ClaimBoundary == "" {
			t.Fatalf("invalid row shape: %+v", row)
		}
	}
}

func assertValidationHealthJSONAllowlist(t *testing.T, payload []byte) {
	t.Helper()
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatal(err)
	}
	wantTop := map[string]bool{"generated_at": true, "agency_id": true, "overall_status": true, "tooling_status": true, "feeds": true, "external_evidence_created": true, "consumer_statuses_changed": true, "compliance_claimed": true, "production_readiness_claimed": true}
	for key := range decoded {
		if !wantTop[key] {
			t.Fatalf("unexpected top-level field %q in %s", key, payload)
		}
	}
	wantRow := map[string]bool{"feed_type": true, "validator_id": true, "validator_name": true, "tooling_status": true, "artifact_status": true, "latest_result_status": true, "latest_result_at": true, "active_feed_version_id": true, "latest_result_feed_version_id": true, "stale_status": true, "health_status": true, "next_action": true, "claim_boundary": true}
	for _, item := range decoded["feeds"].([]any) {
		row := item.(map[string]any)
		for key := range row {
			if !wantRow[key] {
				t.Fatalf("unexpected row field %q in %s", key, payload)
			}
		}
	}
}

func assertValidationHealthHTTPNoLeakage(t testing.TB, body string) {
	t.Helper()
	for _, forbidden := range []string{"raw_report", "stdout", "stderr", "argv", "/tmp/private", "TOKEN=", "SECRET", "PASSWORD=", "postgres://", "Authorization", "Bearer", "Cookie", "admin_session", "DATABASE_URL"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("validation health leaked %q in %s", forbidden, body)
		}
	}
}

func assertValidatorHealthScriptOutputFilesSafe(t testing.TB, outputDir string) {
	t.Helper()
	for _, name := range []string{"summary.json", "summary.md", "manifest.json", "manifest.md"} {
		body, err := os.ReadFile(filepath.Join(outputDir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		assertValidatorHealthGeneratedFileNoSecretValues(t, string(body))
		if len(body) > 16000 {
			t.Fatalf("%s size = %d, want bounded", name, len(body))
		}
	}
}

func assertValidatorHealthGeneratedFileNoSecretValues(t testing.TB, body string) {
	t.Helper()
	for _, forbidden := range []string{
		"Authorization: Bearer",
		"Cookie:",
		"admin_session=",
		"csrf_token=",
		"DATABASE_URL=",
		"postgres://user:pass@",
		"/tmp/private",
		"/Users/",
		"script-admin-token-value",
		"TOKEN=",
		"SECRET=",
		"PASSWORD=",
		"BEGIN PRIVATE KEY",
	} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("validator-health output leaked %q in %s", forbidden, body)
		}
	}
}

func largeOperationsCanonicalRecord(total int) compliance.ValidationReportRecord {
	notices := make([]any, 0, total)
	codes := []string{"expired_calendar", "route_short_name_too_long", "unused_shape", "stop_times_arrival_time_missing", "duplicate_trip_id", "shape_dist_traveled_decreases", "frequency_headway_invalid", "block_id_gap"}
	severities := []string{"ERROR", "WARNING", "INFO"}
	for i := 0; i < total; i++ {
		notices = append(notices, map[string]any{"code": codes[i%len(codes)], "severity": severities[i%len(severities)], "message": strings.Repeat("x", 5000), "path": "/tmp/private/raw.zip"})
	}
	return compliance.ValidationReportRecord{ID: 1, CreatedAt: time.Now().UTC(), Result: compliance.ValidationResult{AgencyID: "demo-agency", FeedType: "schedule", FeedVersionID: "feed-v1", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "warning", WarningCount: total, Report: map[string]any{"raw_report": map[string]any{"notices": notices}, "stdout": "secret", "stderr": "secret", "argv": []any{"/tmp/private/bin"}}}}
}

type countingScheduleBuilder struct {
	snapshot      schedule.Snapshot
	err           error
	readyErr      error
	snapshotCalls int
}

func (f *countingScheduleBuilder) Ready(context.Context) error {
	return f.readyErr
}

func (f *countingScheduleBuilder) Snapshot(context.Context, time.Time) (schedule.Snapshot, error) {
	f.snapshotCalls++
	if f.err != nil {
		return schedule.Snapshot{}, f.err
	}
	return f.snapshot, nil
}

func (f *countingScheduleBuilder) SnapshotForFeedVersion(_ context.Context, feedVersionID string, _ time.Time) (schedule.Snapshot, error) {
	f.snapshotCalls++
	if f.err != nil {
		return schedule.Snapshot{}, f.err
	}
	snapshot := f.snapshot
	if feedVersionID != "" {
		snapshot.FeedVersionID = feedVersionID
	}
	return snapshot, nil
}

type fakeScheduleBuilder struct {
	snapshot          schedule.Snapshot
	snapshotsByAgency map[string]schedule.Snapshot
	err               error
	readyErr          error
}

func (f fakeScheduleBuilder) Ready(context.Context) error {
	return f.readyErr
}

func (f fakeScheduleBuilder) Snapshot(context.Context, time.Time) (schedule.Snapshot, error) {
	if f.err != nil {
		return schedule.Snapshot{}, f.err
	}
	return f.snapshot, nil
}

func (f fakeScheduleBuilder) SnapshotForFeedVersion(_ context.Context, feedVersionID string, _ time.Time) (schedule.Snapshot, error) {
	if f.err != nil {
		return schedule.Snapshot{}, f.err
	}
	snapshot := f.snapshot
	if feedVersionID != "" {
		snapshot.FeedVersionID = feedVersionID
	}
	return snapshot, nil
}

func (f fakeScheduleBuilder) SnapshotForAgency(_ context.Context, agencyID string, _ time.Time) (schedule.Snapshot, error) {
	if f.err != nil {
		return schedule.Snapshot{}, f.err
	}
	if f.snapshotsByAgency != nil {
		return f.snapshotsByAgency[agencyID], nil
	}
	snapshot := f.snapshot
	snapshot.AgencyID = agencyID
	return snapshot, nil
}

type fakePublicationStore struct {
	mu                      sync.Mutex
	result                  compliance.ValidationResult
	validationRecords       []compliance.ValidationReportRecord
	storeValidationCalls    int
	bootstrapInput          compliance.BootstrapInput
	bootstrapErr            error
	publicationConfig       compliance.PublicationConfig
	publicationConfigErr    error
	discovery               compliance.FeedDiscovery
	discoveries             map[string]compliance.FeedDiscovery
	discoveryErr            error
	discoveryAgencyID       string
	scorecard               compliance.Scorecard
	scorecardErr            error
	latestScorecardAgencyID string
	buildScorecardAgencyID  string
	consumers               []compliance.ConsumerRecord
	consumersErr            error
	listConsumersAgencyID   string
	tripDiagnostics         compliance.TripUpdatesDiagnosticsSummary
	tripDiagnosticsErr      error
	reliabilityHealth       []compliance.ReliabilityFeedHealthRecord
	reliabilityIncidents    compliance.ReliabilityIncidentRollup
	reliabilityHealthErr    error
	reliabilityIncidentsErr error
}

func (f *fakePublicationStore) BootstrapPublication(_ context.Context, input compliance.BootstrapInput) error {
	f.bootstrapInput = input
	if f.bootstrapErr != nil {
		return f.bootstrapErr
	}
	return nil
}

func (f *fakePublicationStore) PublicationConfig(context.Context, string) (compliance.PublicationConfig, error) {
	if f.publicationConfigErr != nil {
		return compliance.PublicationConfig{}, f.publicationConfigErr
	}
	return f.publicationConfig, nil
}

func (f *fakePublicationStore) FeedDiscovery(_ context.Context, agencyID string, _ time.Time) (compliance.FeedDiscovery, error) {
	f.discoveryAgencyID = agencyID
	if f.discoveryErr != nil {
		return compliance.FeedDiscovery{}, f.discoveryErr
	}
	if f.discoveries != nil {
		return f.discoveries[agencyID], nil
	}
	return f.discovery, nil
}

func (f *fakePublicationStore) UpsertConsumer(context.Context, compliance.ConsumerInput) (compliance.ConsumerRecord, error) {
	return compliance.ConsumerRecord{}, nil
}

func (f *fakePublicationStore) ListConsumers(_ context.Context, agencyID string) ([]compliance.ConsumerRecord, error) {
	f.listConsumersAgencyID = agencyID
	if f.consumersErr != nil {
		return nil, f.consumersErr
	}
	return f.consumers, nil
}

func (f *fakePublicationStore) LatestScorecard(_ context.Context, agencyID string) (compliance.Scorecard, error) {
	f.latestScorecardAgencyID = agencyID
	if f.scorecardErr != nil {
		return compliance.Scorecard{}, f.scorecardErr
	}
	return f.scorecard, nil
}

func (f *fakePublicationStore) LatestTripUpdatesDiagnostics(context.Context, string) (compliance.TripUpdatesDiagnosticsSummary, error) {
	if f.tripDiagnosticsErr != nil {
		return compliance.TripUpdatesDiagnosticsSummary{}, f.tripDiagnosticsErr
	}
	return f.tripDiagnostics, nil
}

func (f *fakePublicationStore) BuildAndStoreScorecard(_ context.Context, agencyID string, _ time.Time) (compliance.Scorecard, error) {
	f.buildScorecardAgencyID = agencyID
	scorecard := f.scorecard
	if scorecard.AgencyID == "" {
		scorecard.AgencyID = agencyID
	}
	return scorecard, nil
}

func (f *fakePublicationStore) StoreValidationResult(_ context.Context, result compliance.ValidationResult) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.result = result
	f.storeValidationCalls++
	f.validationRecords = append(f.validationRecords, compliance.ValidationReportRecord{ID: int64(len(f.validationRecords) + 1), Result: result, CreatedAt: time.Now().UTC()})
	return nil
}

func (f *fakePublicationStore) LatestValidationReport(_ context.Context, agencyID string, feedType string, validatorName string) (*compliance.ValidationReportRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var latest *compliance.ValidationReportRecord
	for i := range f.validationRecords {
		record := f.validationRecords[i]
		if record.Result.AgencyID != agencyID || record.Result.FeedType != feedType || record.Result.ValidatorName != validatorName {
			continue
		}
		if latest == nil || record.CreatedAt.After(latest.CreatedAt) || record.CreatedAt.Equal(latest.CreatedAt) && record.ID > latest.ID {
			copyRecord := record
			latest = &copyRecord
		}
	}
	if latest == nil {
		return nil, errors.New("not found")
	}
	return latest, nil
}

func (f *fakePublicationStore) LatestReliabilityFeedHealth(context.Context, string, int) ([]compliance.ReliabilityFeedHealthRecord, error) {
	if f.reliabilityHealthErr != nil {
		return nil, f.reliabilityHealthErr
	}
	return f.reliabilityHealth, nil
}

func (f *fakePublicationStore) ReliabilityIncidentRollup(context.Context, string, time.Time, int) (compliance.ReliabilityIncidentRollup, error) {
	if f.reliabilityIncidentsErr != nil {
		return compliance.ReliabilityIncidentRollup{}, f.reliabilityIncidentsErr
	}
	return f.reliabilityIncidents, nil
}

type fakeRealtimeArtifacts struct {
	payloads map[string][]byte
	errors   map[string]error
	calls    map[string]int
}

func (f *fakeRealtimeArtifacts) RealtimePB(_ context.Context, feedType string) ([]byte, string, error) {
	if f.calls == nil {
		f.calls = map[string]int{}
	}
	f.calls[feedType]++
	if err := f.errors[feedType]; err != nil {
		return nil, "", err
	}
	if payload := f.payloads[feedType]; len(payload) > 0 {
		return payload, "internal_builder", nil
	}
	return []byte("protobuf-" + feedType), "internal_builder", nil
}

type fakeDeviceStore struct{}

func (fakeDeviceStore) Verify(context.Context, devices.VerifyInput) (devices.Credential, error) {
	return devices.Credential{}, nil
}

func (fakeDeviceStore) Rebind(context.Context, devices.RebindInput) (devices.RebindResult, error) {
	return devices.RebindResult{}, nil
}

func (fakeDeviceStore) ListBindings(context.Context, string) ([]devices.Binding, error) {
	return nil, nil
}

type fakeDeviceStoreWithBindings struct {
	bindings []devices.Binding
}

func (f fakeDeviceStoreWithBindings) Verify(context.Context, devices.VerifyInput) (devices.Credential, error) {
	return devices.Credential{}, nil
}

func (f fakeDeviceStoreWithBindings) Rebind(context.Context, devices.RebindInput) (devices.RebindResult, error) {
	return devices.RebindResult{}, nil
}

func (f fakeDeviceStoreWithBindings) ListBindings(_ context.Context, agencyID string) ([]devices.Binding, error) {
	var bindings []devices.Binding
	for _, binding := range f.bindings {
		if binding.AgencyID == "" || binding.AgencyID == agencyID {
			bindings = append(bindings, binding)
		}
	}
	return bindings, nil
}

type fakeDeviceStoreWithToken struct {
	token string
}

func (f *fakeDeviceStoreWithToken) Verify(context.Context, devices.VerifyInput) (devices.Credential, error) {
	return devices.Credential{}, nil
}

func (f *fakeDeviceStoreWithToken) Rebind(_ context.Context, input devices.RebindInput) (devices.RebindResult, error) {
	return devices.RebindResult{AgencyID: input.AgencyID, DeviceID: input.DeviceID, VehicleID: input.VehicleID, Token: f.token, RotatedAt: "2026-04-26T12:00:00Z"}, nil
}

func (f *fakeDeviceStoreWithToken) ListBindings(context.Context, string) ([]devices.Binding, error) {
	return nil, nil
}

type fakeTelemetryRepository struct {
	latest []telemetry.StoredEvent
}

func (f fakeTelemetryRepository) Store(context.Context, telemetry.Event, json.RawMessage) (telemetry.StoreResult, error) {
	return telemetry.StoreResult{}, nil
}

func (f fakeTelemetryRepository) LatestByVehicle(context.Context, string, string) (telemetry.StoredEvent, error) {
	return telemetry.StoredEvent{}, nil
}

func (f fakeTelemetryRepository) ListLatestByAgency(_ context.Context, agencyID string, _ int) ([]telemetry.StoredEvent, error) {
	var latest []telemetry.StoredEvent
	for _, event := range f.latest {
		if event.Event.AgencyID == "" || event.Event.AgencyID == agencyID {
			latest = append(latest, event)
		}
	}
	return latest, nil
}

func (f fakeTelemetryRepository) ListEvents(context.Context, string, int) ([]telemetry.StoredEvent, error) {
	return nil, nil
}

type fakeStateRepository struct {
	assignments map[string]state.Assignment
}

func (f fakeStateRepository) ActiveManualOverride(context.Context, string, string, time.Time) (*state.ManualOverride, error) {
	return nil, nil
}

func (f fakeStateRepository) CurrentAssignment(context.Context, string, string) (*state.Assignment, error) {
	return nil, nil
}

func (f fakeStateRepository) ListCurrentAssignments(context.Context, string, []string) (map[string]state.Assignment, error) {
	return f.assignments, nil
}

func (f fakeStateRepository) SaveAssignment(context.Context, state.Assignment, []state.Incident) (state.Assignment, error) {
	return state.Assignment{}, nil
}

func newOperationsTestHandler(h *handler, admin adminAuth) http.Handler {
	if h.store == nil {
		h.store = &fakePublicationStore{}
	}
	if h.devices == nil {
		h.devices = fakeDeviceStore{}
	}
	if h.csrfSecret == "" {
		h.csrfSecret = "test-csrf"
	}
	mux := http.NewServeMux()
	adminRead := admin.Require(auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	mux.Handle("/admin/operations", adminRead(http.HandlerFunc(h.operationsRoot)))
	mux.Handle("/admin/operations/launchpad", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/launchpad.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/checklist", adminRead(http.HandlerFunc(h.operationsRoot)))
	mux.Handle("/admin/operations/checklist.json", adminRead(http.HandlerFunc(h.operationsRoot)))
	mux.Handle("/admin/operations/gtfs-quality", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.Body = http.MaxBytesReader(w, r.Body, gtfsQualityPostMaxBytes)
		}
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/validation-health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if r.Method == http.MethodPost {
			r.Body = http.MaxBytesReader(w, r.Body, validationHealthPostMaxBytes)
		}
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/validation-health.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/reliability", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/reliability.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/", adminRead(http.HandlerFunc(h.operationsRoot)))
	return mux
}

func readyDiscovery() compliance.FeedDiscovery {
	return compliance.FeedDiscovery{Readiness: compliance.Readiness{AllRequiredFeedsListed: true}}
}

type fakePinger struct {
	err error
}

func (f fakePinger) Ping(context.Context) error {
	return f.err
}

type authRejectAll struct{}

func (authRejectAll) Require(...auth.Role) func(http.Handler) http.Handler {
	return func(_ http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		})
	}
}
