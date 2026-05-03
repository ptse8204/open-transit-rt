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
	"path/filepath"
	"strings"
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

	for _, section := range []string{"", "/feeds", "/telemetry", "/devices", "/consumers", "/evidence", "/setup"} {
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

type fakeScheduleBuilder struct {
	snapshot schedule.Snapshot
	err      error
	readyErr error
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

type fakePublicationStore struct {
	result                  compliance.ValidationResult
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
	f.result = result
	return nil
}

type fakeRealtimeArtifacts struct {
	payloads map[string][]byte
	calls    map[string]int
}

func (f *fakeRealtimeArtifacts) RealtimePB(_ context.Context, feedType string) ([]byte, string, error) {
	if f.calls == nil {
		f.calls = map[string]int{}
	}
	f.calls[feedType]++
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
