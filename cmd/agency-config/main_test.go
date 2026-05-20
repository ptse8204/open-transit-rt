package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"open-transit-rt/internal/admincontrol"
	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
	"open-transit-rt/internal/devices"
	"open-transit-rt/internal/feed/schedule"
	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/prediction"
	"open-transit-rt/internal/realtimequality"
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

func TestPublicScheduleUsesCachedSnapshot(t *testing.T) {
	builder := &countingScheduleBuilder{snapshot: schedule.Snapshot{
		AgencyID:      "demo-agency",
		FeedVersionID: "feed-demo",
		RevisionTime:  time.Now().UTC(),
		Payload:       []byte("schedule zip bytes"),
	}}
	handler := newHandlerWithRealtime(
		"demo-agency",
		builder,
		&fakePublicationStore{},
		fakeDeviceStore{},
		fakePinger{},
		authRejectAll{},
		&fakeRealtimeArtifacts{},
	)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/public/gtfs/schedule.zip", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want 200", i+1, rr.Code)
		}
	}
	if builder.snapshotCalls != 1 {
		t.Fatalf("schedule snapshots built %d times, want 1 cached build", builder.snapshotCalls)
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

func TestLocalAdminLoginCreatesBrowserSessionWithoutTokenExposure(t *testing.T) {
	handler := newLocalAdminLoginTestHandler(t, "dev", "true")
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/admin/local-login", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Set-Cookie"); got != "" {
		t.Fatalf("GET set cookie %q, want none", got)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Start setup") || !strings.Contains(body, "Local demo sign-in") {
		t.Fatalf("GET body missing local sign-in UI: %s", body)
	}
	if strings.Contains(body, "eyJ") || strings.Contains(body, "admin_session=") || strings.Contains(body, "Bearer ") {
		t.Fatalf("GET body exposed token-like text: %s", body)
	}
	state := extractLocalLoginState(t, body)

	req = httptest.NewRequest(http.MethodPost, "http://localhost:8080/admin/local-login", strings.NewReader("state="+url.QueryEscape(state)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("POST status = %d, want 303: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Location"); got != "/admin/operations" {
		t.Fatalf("Location = %q, want /admin/operations", got)
	}
	cookies := rr.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != "admin_session" {
		t.Fatalf("cookies = %#v, want one admin_session cookie", cookies)
	}
	session := cookies[0]
	if !session.HttpOnly || session.Path != "/admin" || session.SameSite != http.SameSiteLaxMode || session.MaxAge <= 0 {
		t.Fatalf("session cookie missing safe attributes: %#v", session)
	}
	if strings.Contains(rr.Body.String(), session.Value) || strings.Contains(rr.Body.String(), "Bearer ") {
		t.Fatalf("POST body exposed token-like text: %s", rr.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "http://localhost:8080/admin/operations", nil)
	req.AddCookie(session)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("authenticated operations status = %d, want 200: %s", rr.Code, rr.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "http://localhost:8080/admin/local-login", strings.NewReader("state="+url.QueryEscape(state)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("reused state status = %d, want 403", rr.Code)
	}
}

func TestLocalAdminLoginIsDisabledOutsideLocalDemo(t *testing.T) {
	for _, tc := range []struct {
		name    string
		appEnv  string
		enabled string
		target  string
	}{
		{name: "production", appEnv: "production", enabled: "true", target: "http://localhost:8080/admin/local-login"},
		{name: "flag disabled", appEnv: "dev", enabled: "false", target: "http://localhost:8080/admin/local-login"},
		{name: "non local host", appEnv: "dev", enabled: "true", target: "http://example.com/admin/local-login"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handler := newLocalAdminLoginTestHandler(t, tc.appEnv, tc.enabled)
			req := httptest.NewRequest(http.MethodGet, tc.target, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404: %s", rr.Code, rr.Body.String())
			}
			if got := rr.Header().Get("Set-Cookie"); got != "" {
				t.Fatalf("disabled route set cookie %q", got)
			}
		})
	}
}

func TestLocalAdminLoginAcceptsLoopbackHostsAndExactRouteOnly(t *testing.T) {
	for _, target := range []string{
		"http://localhost:8080/admin/local-login",
		"http://127.0.0.1:8080/admin/local-login",
		"http://[::1]:8080/admin/local-login",
	} {
		t.Run(target, func(t *testing.T) {
			handler := newLocalAdminLoginTestHandler(t, "dev", "true")
			req := httptest.NewRequest(http.MethodGet, target, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200 for %s: %s", rr.Code, target, rr.Body.String())
			}
		})
	}

	handler := newLocalAdminLoginTestHandler(t, "dev", "true")
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/admin/local-login/extra", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("subroute status = %d, want 404", rr.Code)
	}
}

func TestLocalAdminLoginRejectsMissingAndTamperedState(t *testing.T) {
	handler := newLocalAdminLoginTestHandler(t, "dev", "true")
	for _, body := range []string{
		"",
		"state=not-a-valid-state",
	} {
		req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/admin/local-login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("body %q status = %d, want 403", body, rr.Code)
		}
		if got := rr.Header().Get("Set-Cookie"); got != "" {
			t.Fatalf("rejected POST set cookie %q", got)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/admin/local-login", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	state := extractLocalLoginState(t, rr.Body.String())
	replacement := "0"
	if strings.HasSuffix(state, replacement) {
		replacement = "1"
	}
	tampered := state[:len(state)-1] + replacement
	req = httptest.NewRequest(http.MethodPost, "http://localhost:8080/admin/local-login", strings.NewReader("state="+url.QueryEscape(tampered)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("tampered state status = %d, want 403", rr.Code)
	}
	if got := rr.Header().Get("Set-Cookie"); got != "" {
		t.Fatalf("tampered POST set cookie %q", got)
	}
}

func TestLocalAdminLoginPreservesBearerAndCookieCSRFBehavior(t *testing.T) {
	handler := newLocalAdminLoginTestHandler(t, "dev", "true")
	cfg := auth.JWTConfig{Secrets: []string{"test-admin-secret"}, Issuer: "test-issuer", Audience: "test-audience", TTL: time.Hour}
	signer, err := auth.NewSigner(cfg)
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	token, _, err := signer.Sign("admin@example.com", "demo-agency", time.Hour)
	if err != nil {
		t.Fatalf("sign bearer token: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/admin/operations", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.AddCookie(&http.Cookie{Name: "admin_session", Value: "invalid-cookie-token"})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("bearer-auth operations status = %d, want 200: %s", rr.Code, rr.Body.String())
	}

	session := localAdminLoginSessionCookie(t, handler)
	req = httptest.NewRequest(http.MethodPost, "http://localhost:8080/admin/operations/validation-health/refresh.json", strings.NewReader("action=refresh"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(session)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("cookie POST without CSRF status = %d, want 403: %s", rr.Code, rr.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "http://localhost:8080/admin/publication/bootstrap", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(session)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("cookie bootstrap POST without CSRF status = %d, want 403: %s", rr.Code, rr.Body.String())
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
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	body := rr.Body.String()
	for _, want := range []string{"Operations Console", "<title>Start</title>", `<h1 id="operations-page-title">Start</h1>`, "Normal browser path", "Administrator path", "Copy Feed URLs", "publication metadata is not configured yet", "telemetry repository is not available", "no Trip Updates diagnostics recorded yet"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
	for _, want := range []string{"Work through this in order", "Operations workflow", "Start setup", "Import Schedule", "Check feeds", "Connect vehicles", "Fix issues", "Share public URLs", "Maintain system"} {
		if !strings.Contains(body, want) {
			t.Fatalf("action-first body does not contain %q: %s", want, body)
		}
	}
}

func TestOperationsConsoleShowsServerOwnedAgencyScope(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "operator@example.com", AgencyID: "agency-a", Roles: []auth.Role{auth.RoleOperator, auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations?agency_id=agency-a", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("same agency status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Agency scope", "<code>agency-a</code>", "authenticated principal agency", "locked to authenticated agency", "agency_id query values must match this agency", "operator, read_only", "URL edits"} {
		if !strings.Contains(body, want) {
			t.Fatalf("agency scope body missing %q: %s", want, body)
		}
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations?agency_id=agency-b", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("conflicting agency status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
	if strings.Contains(rr.Body.String(), "agency-b") {
		t.Fatalf("forbidden response leaked conflicting agency: %s", rr.Body.String())
	}
}

func TestOperationsAccessRolesAndDeniedUX(t *testing.T) {
	srv := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator, auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/access", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("access status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Access &amp; Roles", "Private role and agency-scope guidance", "Admin", "Editor", "Operator", "Read only", "Role is not allowed", "Agency scope conflict", "Form safety check failed", "operator, read_only"} {
		if !strings.Contains(body, want) {
			t.Fatalf("access body missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"agency approved", "consumer accepted", "production ready", "hosted SaaS", "vendor compatible", "certified hardware"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("access body overclaims %q: %s", forbidden, body)
		}
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/access.json", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.HasPrefix(rr.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("access JSON response = %d %q %s", rr.Code, rr.Header().Get("Content-Type"), rr.Body.String())
	}
	var view operationsAccessView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode access JSON: %v", err)
	}
	if view.AgencyID != "demo-agency" || len(view.Roles) != 4 || len(view.Denied) != 3 || strings.Join(view.CurrentRoles, ",") != "operator,read_only" {
		t.Fatalf("unexpected access view: %+v", view)
	}

	readOnly := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-import", strings.NewReader("action=import_gtfs"))
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	readOnly.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("denied status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
	denied := rr.Body.String()
	for _, want := range []string{"Access denied", "role not allowed", "Open Access &amp; Roles"} {
		if !strings.Contains(denied, want) {
			t.Fatalf("denied response missing %q: %s", want, denied)
		}
	}
	for _, forbidden := range []string{"reader@example.com", "demo-agency", "raw_report", "Authorization", "Bearer ", "token", "database_url"} {
		if strings.Contains(denied, forbidden) {
			t.Fatalf("denied response leaked %q: %s", forbidden, denied)
		}
	}
}

func TestOperationsFeedsPageShowsPublicFeedReadinessReview(t *testing.T) {
	srv := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/feeds", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("feeds status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Feed URLs And Validation",
		"Private public-feed readiness review only",
		"Configured feed URL review",
		"Source-of-truth metadata checklist",
		"Source-of-truth listing guidance",
		"Off-host validation guidance",
		"Public docs portal alignment",
		"Future final-root/evidence checklist",
		"Copy guidance:",
		`id="feed-readiness-feeds_json"`,
		`data-copy-value="https://feeds.example.org/public/feeds.json"`,
		"https://feeds.example.org/public/gtfs/schedule.zip",
		"https://feeds.example.org/public/gtfsrt/vehicle_positions.pb",
		"https://feeds.example.org/public/gtfsrt/trip_updates.pb",
		"https://feeds.example.org/public/gtfsrt/alerts.pb",
		"feeds.json is metadata, not a GTFS validator artifact",
		"endpoint_available=true",
		"public_base_url=true; license=true; contact=true; all_required_listed=true; https=true; discoverable=true",
		"Provider or regional source-of-truth listing",
		"Screenshot and diagram policy",
		"Static schedule validator",
		"Realtime validators",
		"Small-host validation offload",
		"Browser-first docs path",
		"Feed URL share/copy guidance",
		"Future operator checklist",
		"docs/index.md",
		"Requires separate written authorization",
		"Detailed safety booleans remain available",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("feeds body missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"consumer accepted", "submission complete", "ingestion confirmed", "final-root ready", "production ready", "hosted SaaS", "vendor compatible", "certified hardware", "database_url", "Bearer "} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("feeds body overclaims or leaks %q: %s", forbidden, body)
		}
	}

	urlOnlyDiscovery := validationHealthTestDiscovery(time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	urlOnlyDiscovery.PublicBaseURL = "https://feeds.example.org"
	urlOnlyDiscovery.Readiness = compliance.Readiness{AllRequiredFeedsListed: true, LicenseComplete: true, ContactComplete: true, HTTPSURLs: true, Discoverable: true}
	urlOnly := buildOperationsFeedReadiness(operationsPage{Discovery: urlOnlyDiscovery})
	for _, row := range urlOnly.Rows {
		if row.ID != "feeds_json" && row.Status == operationsStatusReady {
			t.Fatalf("feed URL row %s status = ready with URL-only metadata, want validation/feed-health review first: %+v", row.ID, row)
		}
	}
}

func TestOperationsAuditLogBrowserScopedMetadata(t *testing.T) {
	now := time.Date(2026, 5, 14, 10, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{
		auditRows: []compliance.AuditLogRecord{
			{ID: 10, CreatedAt: now, Action: "device.rebind", EntityType: "device_binding", EntityID: "device-1", ActorRecorded: true, ReasonRecorded: true, OldValueRecorded: true, NewValueRecorded: true},
			{ID: 9, CreatedAt: now.Add(-time.Minute), Action: "prediction_override.create", EntityType: "manual_override", EntityID: "bus-1", ActorRecorded: true, ReasonRecorded: false, OldValueRecorded: false, NewValueRecorded: true},
		},
	}
	srv := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator, auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/audit", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("audit status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Audit Log", "Recent scoped audit metadata", "Visible rows", "2 of the latest 50", "actor=2; reason=1; old=1; new=2", "device.rebind", "manual_override", "Raw actor identifiers", "credential values"} {
		if !strings.Contains(body, want) {
			t.Fatalf("audit body missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"operator@example.com", "because dispatch requested it", "old_value_json", "new_value_json", "payload_json", "Authorization", "Bearer ", "postgres://", "/Users/", "secret"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("audit body leaked %q: %s", forbidden, body)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/audit.json", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.HasPrefix(rr.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("audit JSON response = %d %q %s", rr.Code, rr.Header().Get("Content-Type"), rr.Body.String())
	}
	var view operationsAuditView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode audit JSON: %v", err)
	}
	if store.auditAgencyID != "demo-agency" || view.AgencyID != "demo-agency" || len(view.Rows) != 2 {
		t.Fatalf("unexpected audit view agency=%q storeAgency=%q rows=%d view=%+v", view.AgencyID, store.auditAgencyID, len(view.Rows), view)
	}
	if view.Counts.VisibleRows != 2 || view.Counts.QueryLimit != operationsAuditLimit || view.Counts.ActorRecordedRows != 2 || view.Counts.ReasonRecordedRows != 1 || view.Counts.OldValueRecordedRows != 1 || view.Counts.NewValueRecordedRows != 2 || view.Counts.LatestCreatedAt != now.Format(time.RFC3339) {
		t.Fatalf("unexpected audit metadata counts: %+v", view.Counts)
	}
	if view.Rows[0].Action != "device.rebind" || view.Rows[0].ReasonRecorded != true || view.Rows[0].OldValueRecorded != true || view.Rows[0].NewValueRecorded != true {
		t.Fatalf("unexpected first audit row: %+v", view.Rows[0])
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/audit?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("conflicting agency audit status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
	if strings.Contains(rr.Body.String(), "other-agency") {
		t.Fatalf("forbidden audit response leaked conflicting agency: %s", rr.Body.String())
	}
}

func TestOperationsAgencyScopeConflictStopsBeforeAuditDataLoad(t *testing.T) {
	store := &fakePublicationStore{
		auditRows: []compliance.AuditLogRecord{
			{ID: 10, CreatedAt: time.Date(2026, 5, 14, 10, 0, 0, 0, time.UTC), Action: "secret.cross_agency_action", EntityType: "device_binding", EntityID: "secret-device", ActorRecorded: true},
		},
	}
	srv := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator, auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/audit.json?agency_id=other-agency", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("conflicting agency audit JSON status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
	if store.auditAgencyID != "" {
		t.Fatalf("audit store loaded agency %q despite conflicting query agency", store.auditAgencyID)
	}
	for _, forbidden := range []string{"other-agency", "secret.cross_agency_action", "secret-device"} {
		if strings.Contains(rr.Body.String(), forbidden) {
			t.Fatalf("forbidden audit response leaked %q: %s", forbidden, rr.Body.String())
		}
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
}

func TestOperationsCockpitJSONShapeStableCardsAndFlags(t *testing.T) {
	t.Setenv("VALIDATOR_TOOLING_MODE", "stub")
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations.json?agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	var view operationsCockpitView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode cockpit JSON: %v", err)
	}
	assertOperationsCockpitShape(t, view)
	assertOperationsCockpitFlagsFalse(t, view.ClaimFlags)
	if view.AgencyID != "demo-agency" {
		t.Fatalf("agency_id = %q, want demo-agency", view.AgencyID)
	}
	wantActions := []string{"start_setup", "import_gtfs", "check_feeds", "connect_vehicles", "review_realtime", "fix_issues", "share_public_urls", "maintain_system"}
	var gotActions []string
	for _, action := range view.ActionQueue {
		gotActions = append(gotActions, action.ID)
	}
	if strings.Join(gotActions, ",") != strings.Join(wantActions, ",") {
		t.Fatalf("actions = %v, want %v", gotActions, wantActions)
	}
	wantCards := []string{"import_update_gtfs", "review_feed_health", "review_gtfs_quality", "run_review_validator_health", "manage_devices_vehicles", "synthetic_telemetry", "realtime_feed_state", "manage_alerts", "connector_readiness", "maintenance_tasks", "support_summary"}
	var gotCards []string
	for _, card := range view.PrimaryCards {
		gotCards = append(gotCards, card.ID)
	}
	if strings.Join(gotCards, ",") != strings.Join(wantCards, ",") {
		t.Fatalf("cards = %v, want %v", gotCards, wantCards)
	}
	for _, forbidden := range []string{"compliance achieved", "consumer accepted", "agency approved", "hosted SaaS", "production ready", "vendor compatible"} {
		if strings.Contains(strings.ToLower(rr.Body.String()), strings.ToLower(forbidden)) {
			t.Fatalf("cockpit JSON contains forbidden claim %q: %s", forbidden, rr.Body.String())
		}
	}

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/admin/operations.json", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s status = %d, want 405", method, rr.Code)
		}
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations.json?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict status = %d, want 403", rr.Code)
	}
}

func TestOperationsCockpitHTMLShowsNoCLIPrimaryFlow(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Start",
		"Operations workflow",
		"Current status details",
		"More tools",
		`id="cockpit-action-start_setup"`,
		`id="cockpit-action-check_feeds"`,
		`id="cockpit-action-fix_issues"`,
		`id="cockpit-action-share_public_urls"`,
		`id="cockpit-action-maintain_system"`,
		`id="cockpit-card-import_update_gtfs"`,
		`id="cockpit-card-review_feed_health"`,
		`id="cockpit-card-review_gtfs_quality"`,
		`id="cockpit-card-run_review_validator_health"`,
		`id="cockpit-card-manage_devices_vehicles"`,
		`id="cockpit-card-realtime_feed_state"`,
		`id="cockpit-card-maintenance_tasks"`,
		"Agency scope and permissions",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("cockpit HTML missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"compliance achieved", "consumer accepted", "agency approved", "hosted SaaS", "production ready", "vendor compatible"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("cockpit HTML contains forbidden claim %q: %s", forbidden, body)
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
		"Setup Details",
		"Return to Agency Setup",
		"Setup Diagnostics",
		"Role Visibility",
		"Administrator escalation cards",
		"Publication metadata changes require an admin role",
		"Validation runs from this setup page require an admin role",
		"publication metadata",
		"validation records",
		"device bindings",
		"telemetry repository",
		"docs/evidence tracker",
		"not observed yet",
		"prepared is not submitted or accepted",
		"Browser import is admin-only",
		"Validation is supporting evidence only",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{`<form`, `method="post"`, `name="csrf_token"`, "database_url", "restore_database_url"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("read-only setup page contains forbidden %q: %s", forbidden, body)
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
			AgencyID: "demo-agency", AgencyName: `<script>alert("x")</script>`, GeneratedAt: now, PublicationEnvironment: "dev",
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
	if strings.Contains(body, `<script>alert("x")</script>`) {
		t.Fatalf("readiness v2 html did not escape script-like metadata: %s", body)
	}
	for _, want := range []string{
		"Readiness",
		"Private authenticated readiness checklist only",
		"CAL-ITP-Style Readiness Workflow Map",
		"Public feed URLs",
		"Static GTFS",
		"License and contact metadata",
		"Uptime and operations signals",
		"Consumer preparedness",
		"Readiness item",
		"Current signal",
		"What this means",
		"What this helps with",
		"What to do next",
		"What it does not show",
		"Feed discovery and metadata",
		"Plain-language feed health",
		"Static GTFS quality",
		"Vehicle Positions readiness",
		"Trip Updates adapter boundary",
		"Service Alerts readiness",
		"Validator health",
		"Operations reliability diagnostics",
		"Telemetry and device setup",
		"Operations scorecard",
		"Consumer prepared tracker",
		"target-specific statuses",
		"target-originated evidence",
		"This private view does not show submission, review, acceptance, listing, display, ingestion, consumer approval, or Caltrans/CAL-ITP compliance.",
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
		"external_evidence_created",
		"status counts=map",
		"all_required=",
		"public_base_url=",
		"HTTPS=",
		"license=",
		"contact=",
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

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/readiness.json?agency_id=demo-agency", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("json status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("json Cache-Control = %q, want no-store", got)
	}
	var readiness operationsReadinessV2View
	if err := json.Unmarshal(rr.Body.Bytes(), &readiness); err != nil {
		t.Fatalf("decode readiness v2: %v", err)
	}
	assertReadinessV2Shape(t, readiness)
	assertReadinessV2FlagsFalse(t, readiness.ClaimFlags)
	assertReadinessV2SafeStrings(t, rr.Body.String())
	if readiness.AgencyID != "demo-agency" {
		t.Fatalf("agency_id = %q, want demo-agency", readiness.AgencyID)
	}
	metadataGapDiscovery := store.discovery
	metadataGapDiscovery.Readiness.LicenseComplete = false
	metadataGapDiscovery.Readiness.ContactComplete = false
	focusByID := map[string]operationsReadinessV2Focus{}
	for _, focus := range readinessV2FocusAreas(operationsPage{Discovery: metadataGapDiscovery}) {
		focusByID[focus.ID] = focus
	}
	if focusByID["public_feed_urls"].Status != checklistStatusOK {
		t.Fatalf("public URL focus status = %q, want ok when URL/feed listing signals are complete", focusByID["public_feed_urls"].Status)
	}
	if focusByID["license_contact"].Status != checklistStatusNeedsReview {
		t.Fatalf("license/contact focus status = %q, want needs_review when metadata is incomplete", focusByID["license_contact"].Status)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(rr.Body.Bytes(), &top); err != nil {
		t.Fatalf("decode readiness v2 top-level: %v", err)
	}
	wantTop := map[string]bool{"generated_at": true, "agency_id": true, "boundary": true, "focus_areas": true, "rows": true, "counts": true, "claim_flags": true}
	for key := range top {
		if !wantTop[key] {
			t.Fatalf("readiness JSON should return only the v2 model, unexpected top-level key %q in %s", key, rr.Body.String())
		}
	}
	for key := range wantTop {
		if _, ok := top[key]; !ok {
			t.Fatalf("readiness JSON missing top-level key %q in %s", key, rr.Body.String())
		}
	}

	missingHandler := newOperationsTestHandler(&handler{store: &fakePublicationStore{discoveryErr: errors.New("missing discovery"), scorecardErr: errors.New("missing scorecard"), consumersErr: errors.New("missing consumers")}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/readiness.json", nil)
	rr = httptest.NewRecorder()
	missingHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("missing readiness json status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var missing operationsReadinessV2View
	if err := json.Unmarshal(rr.Body.Bytes(), &missing); err != nil {
		t.Fatalf("decode missing readiness v2: %v", err)
	}
	assertReadinessV2Shape(t, missing)
	for _, row := range missing.Rows {
		if row.Status == checklistStatusOK {
			t.Fatalf("missing-data readiness row %s status = ok, want missing/review/blocker/unknown", row.ID)
		}
	}
	assertReadinessV2FlagsFalse(t, missing.ClaimFlags)
}

func TestOperationsReadinessV2RoutesPrivateScopedGETOnlyNoStore(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		t.Run(string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			for _, path := range []string{"/admin/operations/readiness", "/admin/operations/readiness.json"} {
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
	for _, path := range []string{"/admin/operations/readiness", "/admin/operations/readiness.json"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		unauth.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("unauth %s status = %d, want 401", path, rr.Code)
		}
	}

	authenticated := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator}, Method: auth.MethodBearer,
	}})
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		for _, path := range []string{"/admin/operations/readiness", "/admin/operations/readiness.json"} {
			req := httptest.NewRequest(method, path, nil)
			rr := httptest.NewRecorder()
			authenticated.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status = %d, want 405", method, path, rr.Code)
			}
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/readiness?agency_id=other-agency", nil)
	rr := httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict html status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/readiness.json?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict json status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/readiness.json", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public readiness route status = %d, want 404", rr.Code)
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
	for _, want := range []string{"Readiness Checklist", "This checklist is private operator diagnostics", "not evidence", "not an evidence packet", "not compliance proof", "not agency approval", "not consumer acceptance", "not production readiness", "Setup", "Feeds", "Validation", "Telemetry", "Operations", "Consumer Workflow", "Placeholder-like", "Pilot/reference root", "No final-root evidence"} {
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
	assertFirstRunShape(t, launchpad.FirstRun)
	assertFirstRunFlagsFalse(t, launchpad.FirstRun.ClaimFlags)
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
	if launchpad.FirstRun.FeedURLs[0].CopyValue != "https://pilot.example.org/public/feeds.json" {
		t.Fatalf("feeds.json copy value = %q", launchpad.FirstRun.FeedURLs[0].CopyValue)
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
	assertFirstRunShape(t, missing.FirstRun)
	if status := firstRunTaskStatus(missing.FirstRun, "metadata"); status == checklistStatusOK {
		t.Fatalf("missing-data first-run metadata status = ok, want missing/review/blocker/unknown")
	}
	for _, row := range missing.FirstRun.FeedURLs {
		if row.URL == "" && row.CopyValue != "" {
			t.Fatalf("missing-data first-run feed %s copy value = %q, want empty", row.ID, row.CopyValue)
		}
	}
	assertLaunchpadFlagsFalse(t, missing.ClaimFlags)
	assertFirstRunFlagsFalse(t, missing.FirstRun.ClaimFlags)
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
	for _, want := range []string{"Agency Launchpad", "First-run details", "First-Run Tasks", "Copy Feed URLs", "Normal browser path", "Administrator path", "Validation health", "Realtime feeds: Vehicle Positions, Trip Updates, Alerts", "Maintenance and support checks", "Safety details", "creates no evidence", "contacts no external party", "changes no consumer status", "Setup", "GTFS", "Metadata", "Five expected feeds", "Telemetry", "Validators", "Readiness", "Connector conformance", "Support bundle", "Decision gate"} {
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

func TestOperationsDashboardFirstRunAcceptanceWorkflow(t *testing.T) {
	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{
		discovery: compliance.FeedDiscovery{
			AgencyID: "demo-agency", AgencyName: `<script>alert("x")</script>`, GeneratedAt: now, PublicationEnvironment: "pilot",
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

	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	body := rr.Body.String()
	if strings.Contains(body, `<script>alert("x")</script>`) {
		t.Fatalf("dashboard did not escape script-like metadata: %s", body)
	}
	for _, want := range []string{
		"Start",
		"Work through this in order",
		"Operations workflow",
		"Review realtime",
		"Fix issues",
		"Share public URLs",
		"Maintain system",
		"Task status:",
		"Normal browser path",
		"Administrator path",
		`class="feed-copy-grid"`,
		`class="status-chip status-needs-review"`,
		"First-Run Tasks",
		"Metadata",
		"GTFS",
		"Five configured feed URLs",
		"Validation health",
		"Telemetry",
		"Realtime feeds: Vehicle Positions, Trip Updates, Alerts",
		"Readiness",
		"Connectors",
		"Maintenance and support checks",
		"Copy Feed URLs",
		"https://pilot.example.org/public/feeds.json",
		"https://pilot.example.org/public/gtfs/schedule.zip",
		"https://pilot.example.org/public/gtfsrt/vehicle_positions.pb",
		"https://pilot.example.org/public/gtfsrt/trip_updates.pb",
		"https://pilot.example.org/public/gtfsrt/alerts.pb",
		"local wiring only",
		"separate authorized intake",
		"Detailed safety booleans remain available",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("dashboard body missing %q: %s", want, body)
		}
	}
	startHereIndex := strings.Index(body, "Work through this in order")
	helpIndex := strings.Index(body, "Help for Start")
	if startHereIndex < 0 || helpIndex < 0 || startHereIndex > helpIndex {
		t.Fatalf("dashboard should show Start before contextual help: start=%d help=%d body=%s", startHereIndex, helpIndex, body)
	}
	for _, forbidden := range []string{`<form`, `method="post"`, "/admin/operations/first-run", "/public/operations", "agency approved", "consumer accepted", "production ready", "launch complete", "compliance achieved"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("dashboard contains forbidden %q: %s", forbidden, body)
		}
	}
}

func TestSetupWizardRoutesPrivateScopedGETOnlyNoStore(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		t.Run(string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			for _, path := range []string{"/admin/operations/setup-wizard", "/admin/operations/setup-wizard.json"} {
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
	for _, path := range []string{"/admin/operations/setup-wizard", "/admin/operations/setup-wizard.json"} {
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
		for _, path := range []string{"/admin/operations/setup-wizard", "/admin/operations/setup-wizard.json"} {
			req := httptest.NewRequest(method, path, nil)
			rr := httptest.NewRecorder()
			authenticated.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status = %d, want 405", method, path, rr.Code)
			}
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/setup-wizard?agency_id=other-agency", nil)
	rr := httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict html status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/setup-wizard.json?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict json status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/setup-wizard.json", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public setup wizard route status = %d, want 404", rr.Code)
	}
}

func TestSetupWizardJSONShapeFlagsAndStages(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/setup-wizard.json?agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json prefix", got)
	}
	var wizard operationsSetupWizardView
	if err := json.Unmarshal(rr.Body.Bytes(), &wizard); err != nil {
		t.Fatalf("decode setup wizard: %v", err)
	}
	assertSetupWizardShape(t, wizard)
	assertSetupWizardFlagsFalse(t, wizard.ClaimFlags)
	assertSetupWizardSafeStrings(t, rr.Body.String())
	if wizard.AgencyID != "demo-agency" {
		t.Fatalf("agency_id = %q, want demo-agency", wizard.AgencyID)
	}
	var ids []string
	for _, stage := range wizard.Stages {
		ids = append(ids, stage.ID)
	}
	wantIDs := []string{"agency_profile", "publication_metadata", "gtfs", "feeds", "telemetry", "validators", "connectors", "readiness"}
	if strings.Join(ids, ",") != strings.Join(wantIDs, ",") {
		t.Fatalf("stage ids = %v, want %v", ids, wantIDs)
	}

	missingHandler := newOperationsTestHandler(&handler{store: &fakePublicationStore{discoveryErr: errors.New("missing discovery"), scorecardErr: errors.New("missing scorecard"), consumersErr: errors.New("missing consumers")}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/setup-wizard.json", nil)
	rr = httptest.NewRecorder()
	missingHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("missing status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var missing operationsSetupWizardView
	if err := json.Unmarshal(rr.Body.Bytes(), &missing); err != nil {
		t.Fatalf("decode missing setup wizard: %v", err)
	}
	for _, id := range []string{"agency_profile", "publication_metadata", "gtfs", "feeds", "telemetry"} {
		if status := setupWizardStageStatus(missing, id); status == checklistStatusOK {
			t.Fatalf("missing-data stage %s status = ok, want missing/unknown/review/blocker", id)
		}
	}
	assertSetupWizardFlagsFalse(t, missing.ClaimFlags)
}

func TestSetupWizardHTMLBoundariesNoFormsAndEscapes(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/setup-wizard", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Agency Setup", "Setup Progress", "Next Best Step", "Review Blocks And Next Actions", "Setup Diagnostics", "Role Visibility", "Administrator Cards", "Private authenticated setup wizard", "creates no evidence", "changes no state", "Agency profile", "Public feed information", "Schedule data", "Feed links", "Vehicle telemetry", "Validation", "Optional connectors", "Readiness review"} {
		if !strings.Contains(body, want) {
			t.Fatalf("html body missing %q: %s", want, body)
		}
	}
	if strings.Contains(body, `<script>alert("x")</script>`) {
		t.Fatalf("html did not escape script-like metadata: %s", body)
	}
	for _, forbidden := range []string{`<form`, `method="post"`, "/public/operations/setup-wizard", "agency approved", "consumer accepted", "production ready", "launch complete", "compliance achieved", "gtfs upload"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("setup wizard html contains forbidden %q: %s", forbidden, body)
		}
	}
	assertSetupWizardSafeStrings(t, body)
}

func TestGTFSImportRouteAuthMatrixAndBoundaries(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		t.Run("get_"+string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			req := httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-import", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("GET status = %d, want 200: %s", rr.Code, rr.Body.String())
			}
			if got := rr.Header().Get("Cache-Control"); got != "no-store" {
				t.Fatalf("Cache-Control = %q, want no-store", got)
			}
			body := rr.Body.String()
			for _, want := range []string{"Source Review Before Import", "not a preview-only action", "admin role", "CSRF protection", "private URLs", "Administrator needed"} {
				if !strings.Contains(body, want) {
					t.Fatalf("GET body missing source review copy %q: %s", want, body)
				}
			}
			if role == auth.RoleAdmin && !strings.Contains(body, `method="post"`) {
				t.Fatalf("admin GTFS import page does not include import forms: %s", body)
			}
			if role != auth.RoleAdmin && strings.Contains(body, `<form`) {
				t.Fatalf("non-admin GTFS import page includes mutation form: %s", body)
			}
		})
	}
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor} {
		t.Run("post_forbidden_"+string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}, gtfsImport: &fakeGTFSImportRunner{}}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			req := httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-import", strings.NewReader("action=import_gtfs&source_type=url&gtfs_url=https%3A%2F%2Fexample.org%2Fgtfs.zip"))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusForbidden {
				t.Fatalf("POST status = %d, want 403: %s", rr.Code, rr.Body.String())
			}
		})
	}

	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}, gtfsImport: &fakeGTFSImportRunner{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodPut, "/admin/operations/gtfs-import", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PUT status = %d, want 405", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-import?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/gtfs-import", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public route status = %d, want 404", rr.Code)
	}
}

func TestGTFSWorkbenchRoutesPrivateReadOnlyAndJSONBounded(t *testing.T) {
	now := time.Date(2026, 5, 13, 12, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{
		discovery: compliance.FeedDiscovery{
			AgencyID:      "demo-agency",
			AgencyName:    "Demo Agency",
			GeneratedAt:   now,
			PublicBaseURL: "https://agency.example",
			Feeds: []compliance.FeedMetadata{{
				FeedType:             "schedule",
				CanonicalPublicURL:   "https://agency.example/gtfs.zip",
				ActivationStatus:     "active",
				ActiveFeedVersionID:  "feed-v2",
				RevisionTimestamp:    &now,
				LastValidationStatus: "warning",
				LastValidationAt:     &now,
			}},
		},
		gtfsImports: []compliance.GTFSImportRecord{
			{
				ID:             2,
				AgencyID:       "demo-agency",
				FeedVersionID:  "feed-v2",
				SourceFilename: "/Users/private/operator/current-feed.zip",
				SourceSHA256:   "abcdef0123456789",
				SourceByteSize: 2048,
				Status:         "published",
				WarningCount:   1,
				StartedAt:      now,
				CompletedAt:    &now,
			},
			{
				ID:             1,
				AgencyID:       "demo-agency",
				FeedVersionID:  "feed-v1",
				SourceFilename: "previous-feed.zip",
				SourceSHA256:   "1111111111112222",
				SourceByteSize: 1024,
				Status:         "published",
				StartedAt:      now.Add(-time.Hour),
				CompletedAt:    &now,
			},
		},
		gtfsPreviews: map[string]compliance.GTFSSchedulePreview{
			"feed-v2": {
				AgencyID:      "demo-agency",
				FeedVersionID: "feed-v2",
				RowLimit:      10,
				Counts: compliance.GTFSSchedulePreviewCounts{
					Routes:        12,
					Stops:         2,
					Trips:         1,
					StopTimes:     4,
					Calendar:      1,
					CalendarDates: 1,
					ShapePoints:   2,
					Frequencies:   1,
				},
				Agency: compliance.GTFSScheduleAgencyPreview{AgencyID: "demo-agency", Name: "Demo Agency", Timezone: "America/Los_Angeles"},
				Routes: []compliance.GTFSScheduleRoutePreview{
					{ID: "route-01", ShortName: "1", LongName: "Main / <script>", RouteType: "3"},
					{ID: "route-02", ShortName: "2", LongName: "Crosstown", RouteType: "3"},
					{ID: "route-03", ShortName: "3", LongName: "Hill", RouteType: "3"},
					{ID: "route-04", ShortName: "4", LongName: "Lake", RouteType: "3"},
					{ID: "route-05", ShortName: "5", LongName: "Park", RouteType: "3"},
					{ID: "route-06", ShortName: "6", LongName: "Airport", RouteType: "3"},
					{ID: "route-07", ShortName: "7", LongName: "Depot", RouteType: "3"},
					{ID: "route-08", ShortName: "8", LongName: "School", RouteType: "3"},
					{ID: "route-09", ShortName: "9", LongName: "Clinic", RouteType: "3"},
					{ID: "route-10", ShortName: "10", LongName: "Loop", RouteType: "3"},
					{ID: "route-11", ShortName: "11", LongName: "Overflow", RouteType: "3"},
				},
				Stops: []compliance.GTFSScheduleStopPreview{
					{ID: "stop-01", Name: "Main & 1st", Lat: 34.1, Lon: -118.1},
					{ID: "stop-02", Name: "Main & 2nd", Lat: 34.2, Lon: -118.2},
				},
				Trips: []compliance.GTFSScheduleTripPreview{
					{ID: "trip-01", RouteID: "route-01", ServiceID: "weekday", BlockID: "block-1", ShapeID: "shape-1", DirectionID: "0"},
				},
				Calendar: []compliance.GTFSScheduleCalendarPreview{
					{ServiceID: "weekday", Days: "Mon, Tue, Wed, Thu, Fri", StartDate: "20260501", EndDate: "20261231"},
				},
				Frequencies: []compliance.GTFSScheduleFrequencyPreview{
					{TripID: "trip-01", StartTime: "06:00:00", EndTime: "09:00:00", HeadwaySecs: 900, ExactTimes: 0},
				},
			},
			"feed-v1": {
				AgencyID:      "demo-agency",
				FeedVersionID: "feed-v1",
				RowLimit:      10,
				Counts: compliance.GTFSSchedulePreviewCounts{
					Routes:        9,
					Stops:         3,
					Trips:         2,
					StopTimes:     8,
					Calendar:      1,
					CalendarDates: 0,
					ShapePoints:   1,
					Frequencies:   0,
				},
				Agency: compliance.GTFSScheduleAgencyPreview{AgencyID: "demo-agency", Name: "Demo Agency", Timezone: "America/Los_Angeles"},
				Routes: []compliance.GTFSScheduleRoutePreview{
					{ID: "route-01", ShortName: "1", LongName: "Old Main", RouteType: "3"},
					{ID: "route-legacy", ShortName: "L", LongName: "Legacy", RouteType: "3"},
				},
				Stops: []compliance.GTFSScheduleStopPreview{
					{ID: "stop-01", Name: "Main & 1st", Lat: 34.1, Lon: -118.1},
					{ID: "stop-03", Name: "Retired Stop", Lat: 34.3, Lon: -118.3},
				},
				Trips: []compliance.GTFSScheduleTripPreview{
					{ID: "trip-01", RouteID: "route-01", ServiceID: "weekday", BlockID: "block-old", ShapeID: "shape-old", DirectionID: "0"},
					{ID: "trip-retired", RouteID: "route-legacy", ServiceID: "weekday", BlockID: "block-old", ShapeID: "shape-old", DirectionID: "1"},
				},
				Calendar: []compliance.GTFSScheduleCalendarPreview{
					{ServiceID: "weekday", Days: "Mon, Tue, Wed, Thu, Fri", StartDate: "20260101", EndDate: "20260430"},
				},
			},
		},
		gtfsDrafts: []compliance.GTFSDraftRecord{{
			ID:                         "draft-1",
			AgencyID:                   "demo-agency",
			Name:                       "Draft <script>",
			Status:                     "draft",
			BaseFeedVersionID:          "feed-v1",
			LastPublishedFeedVersionID: "feed-v2",
			LastPublishAttemptID:       7,
			CreatedAt:                  now.Add(-2 * time.Hour),
			UpdatedAt:                  now.Add(-time.Hour),
		}},
		gtfsDraftPublishes: []compliance.GTFSDraftPublishRecord{{
			ID:            7,
			DraftID:       "draft-1",
			FeedVersionID: "feed-v2",
			Status:        "published",
			WarningCount:  1,
			StartedAt:     now.Add(-time.Hour),
			CompletedAt:   &now,
		}},
		feedVersions: []compliance.FeedVersionRecord{
			{
				ID:               "feed-v2",
				AgencyID:         "demo-agency",
				SourceType:       "gtfs_studio",
				LifecycleState:   "active",
				IsActive:         true,
				ValidationStatus: "warning",
				PublishedAt:      &now,
				ActivatedAt:      &now,
				CreatedAt:        now.Add(-time.Hour),
			},
			{
				ID:               "feed-v1",
				AgencyID:         "demo-agency",
				SourceType:       "gtfs_import",
				LifecycleState:   "retired",
				IsActive:         false,
				ValidationStatus: "passed",
				PublishedAt:      &now,
				ActivatedAt:      &now,
				RetiredAt:        &now,
				CreatedAt:        now.Add(-2 * time.Hour),
			},
		},
		validationRecords: []compliance.ValidationReportRecord{{
			ID:        1,
			CreatedAt: now,
			Result: compliance.ValidationResult{
				AgencyID:      "demo-agency",
				FeedType:      "schedule",
				FeedVersionID: "feed-v2",
				ValidatorName: compliance.CanonicalStaticValidatorName,
				Status:        "warning",
				WarningCount:  1,
				Report:        map[string]any{"raw_report": map[string]any{"notices": []any{map[string]any{"code": "stop_time_missing_time", "severity": "WARNING", "file": "stop_times.txt", "message": "route review only"}}}},
			},
		}},
	}
	srv := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	for _, path := range []string{"/admin/operations/gtfs-workbench", "/admin/operations/gtfs-workbench.json"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
		}
		if got := rr.Header().Get("Cache-Control"); got != "no-store" {
			t.Fatalf("%s Cache-Control = %q, want no-store", path, got)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-workbench", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	body := rr.Body.String()
	for _, want := range []string{"Schedule Review", "Current Schedule", "Latest Import", "Source checksum", "Agency Review Summary", "Required files", "Row counts", "Service dates", "Routes, stops, and trips", "What changed", "Validation Issue Triage", "Likely owner", "Plain-English meaning", "Suggested fix path", "Safe next action", "Schedule planner with operations review", "Import Change Signals", "Active Vs Previous Schedule Comparison", "File-Level Row Count Diff", "Route / Stop / Trip / Service Change Summary", "Draft-only rollback command design", "Draft Publish Review", "Draft Publish Checklist", "Schedule History And Rollback Guidance", "Rollback Guidance", "Recent Feed Versions", "Preview filters", "Required File Checklist", "Service Calendar Review", "Stop-time service coverage", "Frequency-based service", "Service exceptions", "Routes Preview", "Stops Preview", "Calendar / Service Preview", "No POST action exists"} {
		if !strings.Contains(body, want) {
			t.Fatalf("workbench HTML missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"/Users/private", "<script>", "route-11", "consumer accepted", "agency approved", "production ready", "validator-clean", "CAL-ITP/Caltrans compliant"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("workbench HTML leaks or overclaims %q: %s", forbidden, body)
		}
	}
	if strings.Contains(body, "<form") {
		t.Fatalf("read-only workbench rendered a mutation form: %s", body)
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-workbench.json", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	var decoded map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("decode workbench JSON: %v\n%s", err, rr.Body.String())
	}
	for _, key := range []string{"generated_at", "agency_id", "boundary", "review_summary", "active_feed_version", "import", "version_comparison", "issue_triage", "quality", "validation_health", "preview", "draft_review", "schedule_history", "feed_output", "actions", "claim_flags"} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("workbench JSON missing %q: %#v", key, decoded)
		}
	}
	reviewSummary, ok := decoded["review_summary"].([]any)
	if !ok || len(reviewSummary) != 7 {
		t.Fatalf("review_summary = %#v, want seven staff-facing rows", decoded["review_summary"])
	}
	seenReviewIDs := map[string]bool{}
	for _, item := range reviewSummary {
		row, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("review summary row = %#v, want object", item)
		}
		id, _ := row["id"].(string)
		seenReviewIDs[id] = true
		for _, key := range []string{"label", "status", "plain_language", "suggested_review", "does_not_prove"} {
			if strings.TrimSpace(fmt.Sprint(row[key])) == "" {
				t.Fatalf("review summary %s missing %s: %#v", id, key, row)
			}
		}
	}
	for _, id := range []string{"required_files", "row_counts", "service_dates", "route_stop_trip_review", "import_history", "what_changed", "issue_triage"} {
		if !seenReviewIDs[id] {
			t.Fatalf("review_summary missing %s in %#v", id, reviewSummary)
		}
	}
	issueTriage, ok := decoded["issue_triage"].(map[string]any)
	if !ok {
		t.Fatalf("issue_triage = %#v, want object", decoded["issue_triage"])
	}
	if issueTriage["status"] != "needs_review" || issueTriage["displayed_rows"] != float64(1) || issueTriage["total_rows"] != float64(1) {
		t.Fatalf("issue_triage status/counts = %#v, want one needs_review row", issueTriage)
	}
	issueRows, ok := issueTriage["rows"].([]any)
	if !ok || len(issueRows) != 1 {
		t.Fatalf("issue_triage rows = %#v, want one row", issueTriage["rows"])
	}
	issueRow, ok := issueRows[0].(map[string]any)
	if !ok {
		t.Fatalf("issue_triage row = %#v, want object", issueRows[0])
	}
	for _, key := range []string{"severity", "source_label", "family", "codes", "count", "likely_owner", "plain_english_meaning", "suggested_fix_path", "safe_next_action", "verify_with", "does_not_prove"} {
		if strings.TrimSpace(fmt.Sprint(issueRow[key])) == "" {
			t.Fatalf("issue_triage row missing %s: %#v", key, issueRow)
		}
	}
	if _, ok := issueRow["samples"]; ok {
		t.Fatalf("issue_triage row exposed raw samples: %#v", issueRow)
	}
	if !strings.Contains(fmt.Sprint(issueRow["likely_owner"]), "Schedule planner") || !strings.Contains(fmt.Sprint(issueRow["safe_next_action"]), "Before validation") {
		t.Fatalf("issue_triage row missing owner or safe next action wording: %#v", issueRow)
	}
	flags, ok := decoded["claim_flags"].(map[string]any)
	if !ok {
		t.Fatalf("claim_flags = %#v, want object", decoded["claim_flags"])
	}
	for key, value := range flags {
		if value == true {
			t.Fatalf("claim flag %s unexpectedly true in %#v", key, flags)
		}
	}
	preview, ok := decoded["preview"].(map[string]any)
	if !ok {
		t.Fatalf("preview = %#v, want object", decoded["preview"])
	}
	counts, ok := preview["counts"].(map[string]any)
	if !ok || counts["routes"] != float64(12) {
		t.Fatalf("preview counts = %#v, want 12 routes", preview["counts"])
	}
	serviceWarnings, ok := preview["service_warnings"].([]any)
	if !ok || len(serviceWarnings) < 4 {
		t.Fatalf("preview service_warnings = %#v, want service review rows", preview["service_warnings"])
	}
	serviceWarningsBody := fmt.Sprint(serviceWarnings)
	for _, want := range []string{"Service calendar source", "Service date range", "Stop-time service coverage", "Frequency-based service", "Service exceptions"} {
		if !strings.Contains(serviceWarningsBody, want) {
			t.Fatalf("service_warnings missing %q: %#v", want, serviceWarnings)
		}
	}
	routes, ok := preview["routes"].([]any)
	if !ok || len(routes) != 10 {
		t.Fatalf("preview routes = %#v, want capped 10 rows", preview["routes"])
	}
	versionComparison, ok := decoded["version_comparison"].(map[string]any)
	if !ok || versionComparison["status"] != "needs_review" || versionComparison["active_feed_version_id"] != "feed-v2" || versionComparison["previous_feed_version_id"] != "feed-v1" {
		t.Fatalf("version_comparison = %#v, want needs_review feed-v2 vs feed-v1", decoded["version_comparison"])
	}
	fileDiffs, ok := versionComparison["file_diffs"].([]any)
	if !ok || len(fileDiffs) != 8 {
		t.Fatalf("file_diffs = %#v, want eight GTFS file diff rows", versionComparison["file_diffs"])
	}
	seenFileDiffs := map[string]map[string]any{}
	for _, item := range fileDiffs {
		row, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("file diff row = %#v, want object", item)
		}
		file, _ := row["file"].(string)
		seenFileDiffs[file] = row
	}
	for _, file := range []string{"routes.txt", "stops.txt", "trips.txt", "stop_times.txt", "calendar.txt", "calendar_dates.txt", "shapes.txt", "frequencies.txt"} {
		if _, ok := seenFileDiffs[file]; !ok {
			t.Fatalf("file diffs missing %s in %#v", file, seenFileDiffs)
		}
	}
	if got := seenFileDiffs["routes.txt"]["delta_rows"]; got != float64(3) {
		t.Fatalf("routes delta = %#v, want 3", got)
	}
	if got := seenFileDiffs["frequencies.txt"]["delta_rows"]; got != float64(1) {
		t.Fatalf("frequencies delta = %#v, want 1", got)
	}
	entityDiffs, ok := versionComparison["entity_diffs"].([]any)
	if !ok || len(entityDiffs) == 0 {
		t.Fatalf("entity_diffs = %#v, want bounded entity summaries", versionComparison["entity_diffs"])
	}
	entityBody := fmt.Sprint(entityDiffs)
	for _, want := range []string{"Routes", "Stops", "Trips", "Service calendars", "Frequencies", "route-02", "route-legacy", "route-01"} {
		if !strings.Contains(entityBody, want) {
			t.Fatalf("entity diffs missing %q: %#v", want, entityDiffs)
		}
	}
	draftReview, ok := decoded["draft_review"].(map[string]any)
	if !ok || draftReview["status"] != "needs_review" {
		t.Fatalf("draft_review = %#v, want needs_review object", decoded["draft_review"])
	}
	scheduleHistory, ok := decoded["schedule_history"].(map[string]any)
	if !ok || scheduleHistory["status"] != "ok" {
		t.Fatalf("schedule_history = %#v, want ok object", decoded["schedule_history"])
	}
	jsonBody := rr.Body.String()
	for _, forbidden := range []string{"/Users/private", "route-11", "consumer_statuses_changed\":true", "compliance_claimed\":true", "production_readiness_claimed\":true"} {
		if strings.Contains(jsonBody, forbidden) {
			t.Fatalf("workbench JSON leaks or overclaims %q: %s", forbidden, jsonBody)
		}
	}

	req = httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-workbench", strings.NewReader("action=publish&gtfs_path=/tmp/private.zip"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST status = %d, want 405: %s", rr.Code, rr.Body.String())
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-workbench?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict status = %d, want 403", rr.Code)
	}
	unauth := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, authRejectAll{})
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-workbench", nil)
	rr = httptest.NewRecorder()
	unauth.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("unauth status = %d, want 401", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/gtfs-workbench", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public route status = %d, want 404", rr.Code)
	}
}

func TestGTFSImportUploadUsesTempFileAndImporter(t *testing.T) {
	importer := &fakeGTFSImportRunner{result: gtfs.ImportResult{
		ImportID:      42,
		AgencyID:      "demo-agency",
		FeedVersionID: "gtfs-import-42",
		Status:        gtfs.ImportStatusPublished,
		WarningCount:  1,
		InfoCount:     2,
		Counts:        map[string]int{"stops.txt": 4, "trips.txt": 2},
		ReportStored:  true,
	}}
	srv := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}, gtfsImport: importer}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	body, contentType := gtfsImportMultipartBody(t, "agency.zip", []byte("zip payload"), map[string]string{
		"action":      "import_gtfs",
		"source_type": "upload",
		"notes":       "operator supplied zip",
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-import", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if importer.calls != 1 {
		t.Fatalf("import calls = %d, want 1", importer.calls)
	}
	if importer.opts.AgencyID != "demo-agency" || importer.opts.ActorID != "admin@example.com" {
		t.Fatalf("opts = %+v, want principal agency and actor", importer.opts)
	}
	if !strings.Contains(importer.opts.Notes, "browser_gtfs_import") || !strings.Contains(importer.opts.Notes, "operator supplied zip") {
		t.Fatalf("notes = %q, want browser import note", importer.opts.Notes)
	}
	if !importer.pathExistsDuringCall {
		t.Fatalf("temporary ZIP path did not exist during import")
	}
	if string(importer.payload) != "zip payload" {
		t.Fatalf("payload = %q, want uploaded bytes", importer.payload)
	}
	if _, err := os.Stat(importer.opts.ZipPath); !os.IsNotExist(err) {
		t.Fatalf("temporary ZIP still exists or stat error is not not-exist: %v", err)
	}
	response := rr.Body.String()
	for _, want := range []string{"GTFS import finished", "gtfs-import-42", "stops.txt", "Validation report stored"} {
		if !strings.Contains(response, want) {
			t.Fatalf("response missing %q: %s", want, response)
		}
	}
	for _, forbidden := range []string{importer.opts.ZipPath, "zip payload", "/tmp/", "/var/lib", "agency approved", "consumer accepted", "CAL-ITP/Caltrans compliant"} {
		if forbidden != "" && strings.Contains(strings.ToLower(response), strings.ToLower(forbidden)) {
			t.Fatalf("response leaks or overclaims %q: %s", forbidden, response)
		}
	}
}

func TestGTFSImportZipPreflightSummarizesRequiredFilesAndServiceRows(t *testing.T) {
	zipPath := writeGTFSPreflightZip(t, map[string]string{
		"agency.txt":         "agency_id,agency_name,agency_url,agency_timezone\nA,Demo,https://example.org,America/Los_Angeles\n",
		"routes.txt":         "route_id,agency_id,route_short_name,route_long_name,route_type\nR,A,1,Main,3\n",
		"stops.txt":          "stop_id,stop_name,stop_lat,stop_lon\nS1,Main,34,-118\nS2,Second,34.1,-118.1\n",
		"trips.txt":          "route_id,service_id,trip_id\nR,WKD,T1\n",
		"stop_times.txt":     "trip_id,arrival_time,departure_time,stop_id,stop_sequence\nT1,23:50:00,23:50:00,S1,1\nT1,24:10:00,24:10:00,S2,2\n",
		"calendar_dates.txt": "service_id,date,exception_type\nWKD,20260520,1\n",
		"frequencies.txt":    "trip_id,start_time,end_time,headway_secs,exact_times\nT1,06:00:00,09:00:00,900,0\n",
	})
	rows := gtfsImportZipPreflight(zipPath)
	if len(rows) != 8 {
		t.Fatalf("preflight rows = %d, want 8: %+v", len(rows), rows)
	}
	seen := map[string]operationsGTFSChangeRow{}
	for _, row := range rows {
		seen[row.Label] = row
		if row.CurrentSignal == "" || row.NextAction == "" || row.ClaimBoundary == "" {
			t.Fatalf("preflight row missing guidance: %+v", row)
		}
		if strings.Contains(strings.ToLower(row.CurrentSignal), strings.ToLower(zipPath)) {
			t.Fatalf("preflight leaked temp path: %+v", row)
		}
	}
	for _, required := range []string{"agency.txt", "routes.txt", "stops.txt", "trips.txt", "stop_times.txt", "calendar.txt / calendar_dates.txt"} {
		if seen[required].Status != "ok" {
			t.Fatalf("%s status = %+v, want ok", required, seen[required])
		}
	}
	if seen["shapes.txt"].Status != "optional" {
		t.Fatalf("shapes status = %+v, want optional", seen["shapes.txt"])
	}
	if !strings.Contains(seen["calendar.txt / calendar_dates.txt"].CurrentSignal, "calendar_dates.txt has 1") {
		t.Fatalf("calendar preflight signal = %q", seen["calendar.txt / calendar_dates.txt"].CurrentSignal)
	}
}

func writeGTFSPreflightZip(t *testing.T, files map[string]string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "preflight.zip")
	out, err := os.Create(path)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(out)
	for name, body := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("create zip member %s: %v", name, err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Fatalf("write zip member %s: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := out.Close(); err != nil {
		t.Fatalf("close zip file: %v", err)
	}
	return path
}

func TestGTFSImportValidationFailureRendersBoundedResult(t *testing.T) {
	result := gtfs.ImportResult{
		ImportID:       7,
		AgencyID:       "demo-agency",
		Status:         gtfs.ImportStatusFailed,
		ErrorCount:     2,
		WarningCount:   1,
		ReportStored:   true,
		FailureMessage: "validation failed",
	}
	importer := &fakeGTFSImportRunner{result: result, err: &gtfs.ImportError{Result: result, Err: errors.New("raw validator path /tmp/private/report.json")}}
	srv := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}, gtfsImport: importer}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	body, contentType := gtfsImportMultipartBody(t, "agency.zip", []byte("bad zip payload"), map[string]string{
		"action":      "import_gtfs",
		"source_type": "upload",
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-import", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 bounded validation result: %s", rr.Code, rr.Body.String())
	}
	response := rr.Body.String()
	for _, want := range []string{"validation failed", "not published", "GTFS import finished with stored validation feedback"} {
		if !strings.Contains(response, want) {
			t.Fatalf("response missing %q: %s", want, response)
		}
	}
	for _, forbidden := range []string{"/tmp/private", "bad zip payload", "agency approved", "consumer accepted"} {
		if strings.Contains(strings.ToLower(response), strings.ToLower(forbidden)) {
			t.Fatalf("response leaks or overclaims %q: %s", forbidden, response)
		}
	}
}

func TestGTFSImportURLDownloadAndUnsafeRejection(t *testing.T) {
	importer := &fakeGTFSImportRunner{result: gtfs.ImportResult{ImportID: 8, AgencyID: "demo-agency", FeedVersionID: "gtfs-import-8", Status: gtfs.ImportStatusPublished, ReportStored: true}}
	srv := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}, gtfsImport: importer}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})

	req := httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-import", strings.NewReader(url.Values{
		"action":      {"import_gtfs"},
		"source_type": {"url"},
		"gtfs_url":    {"http://localhost:12345/gtfs.zip?token=secret"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("unsafe status = %d, want 400: %s", rr.Code, rr.Body.String())
	}
	if importer.calls != 0 {
		t.Fatalf("unsafe URL invoked importer %d time(s), want 0", importer.calls)
	}
	if strings.Contains(rr.Body.String(), "localhost:12345") || strings.Contains(rr.Body.String(), "secret") {
		t.Fatalf("unsafe URL response leaks raw URL: %s", rr.Body.String())
	}

	t.Setenv("ALLOW_BROWSER_GTFS_IMPORT_PRIVATE_URLS", "true")
	t.Setenv("ALLOW_BROWSER_GTFS_IMPORT_INSECURE_HTTP", "true")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("downloaded zip"))
	}))
	defer server.Close()
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-import", strings.NewReader(url.Values{
		"action":      {"import_gtfs"},
		"source_type": {"url"},
		"gtfs_url":    {server.URL + "/gtfs.zip"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("download status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if importer.calls != 1 {
		t.Fatalf("download import calls = %d, want 1", importer.calls)
	}
	if string(importer.payload) != "downloaded zip" {
		t.Fatalf("download payload = %q, want downloaded bytes", importer.payload)
	}
	if strings.Contains(rr.Body.String(), server.URL) || strings.Contains(rr.Body.String(), "127.0.0.1") || strings.Contains(strings.ToLower(rr.Body.String()), "localhost") {
		t.Fatalf("download response leaked local URL: %s", rr.Body.String())
	}
}

func TestGTFSImportRejectsUnexpectedFieldsAndCookieCSRF(t *testing.T) {
	importer := &fakeGTFSImportRunner{result: gtfs.ImportResult{ImportID: 9, AgencyID: "demo-agency", Status: gtfs.ImportStatusPublished, ReportStored: true}}
	srv := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}, gtfsImport: importer}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-import", strings.NewReader(url.Values{
		"action":          {"import_gtfs"},
		"source_type":     {"url"},
		"gtfs_url":        {"https://example.org/gtfs.zip"},
		"agency_id":       {"other"},
		"feed_version_id": {"feed-v2"},
		"zip_path":        {"/tmp/private.zip"},
		"argv":            {"--unsafe"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("unexpected field status = %d, want 400: %s", rr.Code, rr.Body.String())
	}
	if importer.calls != 0 {
		t.Fatalf("unexpected fields invoked importer %d time(s), want 0", importer.calls)
	}
	for _, forbidden := range []string{"/tmp/private.zip", "--unsafe", "feed-v2"} {
		if strings.Contains(rr.Body.String(), forbidden) {
			t.Fatalf("unexpected field response leaks %q: %s", forbidden, rr.Body.String())
		}
	}

	cookieAuth := auth.TestAuthenticator{Principal: auth.Principal{Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodCookie}}
	srv = newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}, gtfsImport: importer, csrfSecret: "test-csrf"}, cookieAuth)
	body, contentType := gtfsImportMultipartBody(t, "agency.zip", []byte("zip payload"), map[string]string{
		"action":      "import_gtfs",
		"source_type": "upload",
	})
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-import", body)
	req.Header.Set("Content-Type", contentType)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("cookie missing csrf status = %d, want 403: %s", rr.Code, rr.Body.String())
	}

	token := csrfToken("test-csrf", auth.Principal{Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodCookie})
	body, contentType = gtfsImportMultipartBody(t, "agency.zip", []byte("zip payload"), map[string]string{
		"action":      "import_gtfs",
		"source_type": "upload",
	})
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/gtfs-import?csrf_token="+token, body)
	req.Header.Set("Content-Type", contentType)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("cookie query csrf status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
}

func TestFeedHealthRoutesPrivateScopedGETOnlyNoStore(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		t.Run(string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			for _, path := range []string{"/admin/operations/feed-health", "/admin/operations/feed-health.json"} {
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
	for _, path := range []string{"/admin/operations/feed-health", "/admin/operations/feed-health.json"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		unauth.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("unauth %s status = %d, want 401", path, rr.Code)
		}
	}

	authenticated := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator}, Method: auth.MethodBearer,
	}})
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		for _, path := range []string{"/admin/operations/feed-health", "/admin/operations/feed-health.json"} {
			req := httptest.NewRequest(method, path, nil)
			rr := httptest.NewRecorder()
			authenticated.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status = %d, want 405", method, path, rr.Code)
			}
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/feed-health?agency_id=other-agency", nil)
	rr := httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict html status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/feed-health.json?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict json status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/feed-health.json", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public feed health route status = %d, want 404", rr.Code)
	}
}

func TestFeedHealthJSONShapeFlagsRowsAndMissingData(t *testing.T) {
	t.Setenv("VALIDATOR_TOOLING_MODE", "stub")
	srv := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/feed-health.json?agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json prefix", got)
	}
	var health operationsFeedHealthView
	if err := json.Unmarshal(rr.Body.Bytes(), &health); err != nil {
		t.Fatalf("decode feed health: %v", err)
	}
	assertFeedHealthShape(t, health)
	assertFeedHealthFlagsFalse(t, health.ClaimFlags)
	assertFeedHealthSafeStrings(t, rr.Body.String())
	if health.AgencyID != "demo-agency" {
		t.Fatalf("agency_id = %q, want demo-agency", health.AgencyID)
	}
	var ids []string
	rowsByID := map[string]operationsFeedHealthRow{}
	for _, row := range health.Rows {
		ids = append(ids, row.ID)
		rowsByID[row.ID] = row
	}
	wantIDs := []string{"feeds_json", "schedule", "vehicle_positions", "trip_updates", "alerts"}
	if strings.Join(ids, ",") != strings.Join(wantIDs, ",") {
		t.Fatalf("row ids = %v, want %v", ids, wantIDs)
	}
	wantPaths := map[string]string{
		"feeds_json":        "/public/feeds.json",
		"schedule":          "/public/gtfs/schedule.zip",
		"vehicle_positions": "/public/gtfsrt/vehicle_positions.pb",
		"trip_updates":      "/public/gtfsrt/trip_updates.pb",
		"alerts":            "/public/gtfsrt/alerts.pb",
	}
	for id, wantPath := range wantPaths {
		if rowsByID[id].PublicPath != wantPath {
			t.Fatalf("row %s public_path = %q, want %q", id, rowsByID[id].PublicPath, wantPath)
		}
	}
	if rowsByID["feeds_json"].Status != checklistStatusOK || !strings.Contains(rowsByID["feeds_json"].CurrentSignal, "all HTTPS=true") || !strings.Contains(rowsByID["feeds_json"].CurrentSignal, "discoverable=true") {
		t.Fatalf("feeds_json row did not include HTTPS/discoverability readiness: %+v", rowsByID["feeds_json"])
	}
	if !strings.Contains(rowsByID["vehicle_positions"].HealthContext, "active schedule context=feed-v1") {
		t.Fatalf("vehicle_positions row missing active schedule context: %+v", rowsByID["vehicle_positions"])
	}
	for _, id := range []string{"vehicle_positions", "trip_updates", "alerts"} {
		if strings.Contains(strings.ToLower(rowsByID[id].Freshness), "generated") {
			t.Fatalf("realtime row %s should not label revision metadata as generated freshness: %+v", id, rowsByID[id])
		}
	}
	for _, link := range rowsByID["alerts"].AdminLinks {
		if link == "/admin/alerts/console" {
			t.Fatalf("alerts feed health row links to unregistered alerts console route: %+v", rowsByID["alerts"].AdminLinks)
		}
	}

	notDiscoverableStore := feedHealthTestStore(t)
	notDiscoverableStore.discovery.Readiness.HTTPSURLs = false
	notDiscoverableStore.discovery.Readiness.Discoverable = false
	notDiscoverableHandler := newOperationsTestHandler(&handler{store: notDiscoverableStore, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/feed-health.json", nil)
	rr = httptest.NewRecorder()
	notDiscoverableHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("not-discoverable status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var notDiscoverable operationsFeedHealthView
	if err := json.Unmarshal(rr.Body.Bytes(), &notDiscoverable); err != nil {
		t.Fatalf("decode not-discoverable feed health: %v", err)
	}
	if notDiscoverable.Rows[0].ID != "feeds_json" || notDiscoverable.Rows[0].Status == checklistStatusOK {
		t.Fatalf("feeds_json status with non-HTTPS/non-discoverable metadata = %+v, want needs review", notDiscoverable.Rows[0])
	}

	missingHandler := newOperationsTestHandler(&handler{store: &fakePublicationStore{discoveryErr: errors.New("missing discovery"), scorecardErr: errors.New("missing scorecard"), consumersErr: errors.New("missing consumers")}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/feed-health.json", nil)
	rr = httptest.NewRecorder()
	missingHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("missing status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var missing operationsFeedHealthView
	if err := json.Unmarshal(rr.Body.Bytes(), &missing); err != nil {
		t.Fatalf("decode missing feed health: %v", err)
	}
	assertFeedHealthShape(t, missing)
	for _, row := range missing.Rows {
		if row.Status == checklistStatusOK {
			t.Fatalf("missing-data row %s status = ok, want missing/review/blocker/unknown", row.ID)
		}
	}
	assertFeedHealthFlagsFalse(t, missing.ClaimFlags)
}

func TestFeedHealthHTMLPlainLanguageBoundariesAndEscapes(t *testing.T) {
	t.Setenv("VALIDATOR_TOOLING_MODE", "stub")
	store := feedHealthTestStore(t)
	store.discovery.AgencyName = `<script>alert("x")</script>`
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/feed-health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Feed Health Dashboard", "command center tracks exactly five configured public route paths", "/public/feeds.json", "/public/gtfs/schedule.zip", "/public/gtfsrt/vehicle_positions.pb", "/public/gtfsrt/trip_updates.pb", "/public/gtfsrt/alerts.pb", "feeds.json", "Static GTFS Schedule", "Vehicle Positions", "Trip Updates", "Alerts", "Public path", "What this means", "Freshness", "Validator context", "Health context", "active schedule context", "Next action", "Limits"} {
		if !strings.Contains(body, want) {
			t.Fatalf("html body missing %q: %s", want, body)
		}
	}
	if strings.Contains(body, `<script>alert("x")</script>`) {
		t.Fatalf("html did not escape script-like metadata: %s", body)
	}
	for _, forbidden := range []string{`<form`, `method="post"`, "/public/operations/feed-health", "agency approved", "consumer accepted", "production ready", "launch complete", "compliance achieved", "SLA coverage", "uptime guarantee"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("feed health html contains forbidden %q: %s", forbidden, body)
		}
	}
	assertFeedHealthSafeStrings(t, body)
}

func TestConnectorHubRoutesPrivateScopedGETOnlyNoStore(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		t.Run(string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			for _, path := range []string{"/admin/operations/connectors", "/admin/operations/connectors.json"} {
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
	for _, path := range []string{"/admin/operations/connectors", "/admin/operations/connectors.json"} {
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
		for _, path := range []string{"/admin/operations/connectors", "/admin/operations/connectors.json"} {
			req := httptest.NewRequest(method, path, nil)
			rr := httptest.NewRecorder()
			authenticated.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status = %d, want 405", method, path, rr.Code)
			}
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors?agency_id=other-agency", nil)
	rr := httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict html status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/connectors.json?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict json status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/connectors.json", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public connector route status = %d, want 404", rr.Code)
	}
}

func TestConnectorHubJSONShapeFlagsAndCategories(t *testing.T) {
	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{
		discovery: compliance.FeedDiscovery{
			AgencyID: "demo-agency", AgencyName: "Demo Agency", GeneratedAt: now, PublicationEnvironment: "pilot",
			PublicBaseURL:         "https://pilot.example.org",
			TechnicalContactEmail: "ops@example.org",
			License:               compliance.License{Name: "CC BY 4.0", URL: "https://example.org/license"},
		},
	}
	srv := newOperationsTestHandler(&handler{
		store:   store,
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now}}},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors.json?agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json prefix", got)
	}
	var hub connectorHubView
	if err := json.Unmarshal(rr.Body.Bytes(), &hub); err != nil {
		t.Fatalf("decode connector hub: %v", err)
	}
	assertConnectorHubShape(t, hub)
	assertConnectorHubFlagsFalse(t, hub.ClaimFlags)
	assertConnectorHubSafeStrings(t, rr.Body.String())
	if hub.AgencyID != "demo-agency" {
		t.Fatalf("agency_id = %q, want demo-agency", hub.AgencyID)
	}
	if hub.PluginDefinition != safePluginDefinition {
		t.Fatalf("plugin definition = %q, want safe definition", hub.PluginDefinition)
	}
	wantHealthIDs := []string{"telemetry_source", "prediction", "validator", "monitoring_export", "consumer_discovery", "future_extension_model"}
	var gotHealthIDs []string
	for _, row := range hub.Health {
		gotHealthIDs = append(gotHealthIDs, row.ID)
		if row.ChecklistCopy == "" {
			t.Fatalf("connector health checklist is not safe/copyable: %+v", row)
		}
	}
	if strings.Join(gotHealthIDs, ",") != strings.Join(wantHealthIDs, ",") {
		t.Fatalf("health ids = %v, want %v", gotHealthIDs, wantHealthIDs)
	}
	var ids []string
	for _, category := range hub.Categories {
		ids = append(ids, category.ID)
	}
	if len(hub.Catalog) != 28 {
		t.Fatalf("catalog rows = %d, want 28", len(hub.Catalog))
	}
	wantCatalogIDs := []string{"csv_replay_adapter", "http_polling_adapter", "webhook_sidecar_adapter", "generic_json_transform_adapter", "vendor_shaped_synthetic_examples", "authenticated_telemetry_post", "deterministic_builtin_predictor", "external_http_predictor_adapter", "shadow_mode_predictor", "fail_closed_predictor_behavior", "thetransitclock_candidate_notes", "mobilitydata_static_gtfs_validator", "mobilitydata_gtfs_realtime_validator", "allowlisted_validator_ids", "private_validation_health", "local_health_summaries", "operations_notify_draft", "monitoring_export_helper", "deployment_owned_monitoring_boundary", "public_feeds_json", "static_gtfs_url", "vehicle_positions_url", "trip_updates_url", "alerts_url", "consumer_packet_preparedness", "manifest_based_sidecars", "no_dynamic_backend_plugin_loading", "conformance_tests_required"}
	var gotCatalogIDs []string
	for _, row := range hub.Catalog {
		gotCatalogIDs = append(gotCatalogIDs, row.ID)
		if row.Group == "" || row.Label == "" || row.Status == "" || row.StartWith == "" || row.BrowserReview == "" || row.FirstSafeCheck == "" || row.DoesNotProve == "" || len(row.DocsLinks) == 0 {
			t.Fatalf("invalid catalog row: %+v", row)
		}
	}
	if strings.Join(gotCatalogIDs, ",") != strings.Join(wantCatalogIDs, ",") {
		t.Fatalf("catalog ids = %v, want %v", gotCatalogIDs, wantCatalogIDs)
	}
	wantIDs := []string{"telemetry_source", "prediction", "validator", "monitoring_export", "consumer_discovery", "future_extension_model"}
	if strings.Join(ids, ",") != strings.Join(wantIDs, ",") {
		t.Fatalf("category ids = %v, want %v", ids, wantIDs)
	}
	var registryIDs []string
	for _, entry := range hub.Registry.Entries {
		registryIDs = append(registryIDs, entry.ConnectorID)
		if entry.SourcePath == "" || !strings.HasPrefix(entry.SourcePath, "examples/connectors/") {
			t.Fatalf("registry entry has unsafe source path: %+v", entry)
		}
		if entry.SchemaVersion != "open-transit-rt.connector.v1" || !strings.HasPrefix(entry.ConnectorID, "example.") || entry.DisplayName == "" || entry.ConnectorType == "" || entry.ModeName == "" || entry.DocsLink == "" {
			t.Fatalf("registry entry missing bounded manifest summary: %+v", entry)
		}
		if !entry.DisabledByDefault || !entry.FailureBehavior.FailClosed || len(entry.InputContracts) == 0 || len(entry.OutputContracts) == 0 || len(entry.ConformanceCases) == 0 {
			t.Fatalf("registry entry must remain disabled, fail-closed, and conformance-backed: %+v", entry)
		}
	}
	wantRegistryIDs := []string{"example.consumer-discovery-metadata", "example.generic-json-transform", "example.monitoring-export", "example.predictor-sidecar-stub", "example.telemetry-csv-replay", "example.telemetry-http-poller", "example.telemetry-webhook-sidecar", "example.validator-allowlist"}
	if strings.Join(registryIDs, ",") != strings.Join(wantRegistryIDs, ",") {
		t.Fatalf("registry ids = %v, want %v", registryIDs, wantRegistryIDs)
	}
	if len(hub.Registry.Diagnostics) != 0 {
		t.Fatalf("registry diagnostics = %+v, want none", hub.Registry.Diagnostics)
	}
}

func TestConnectorHubHTMLBoundariesNoFormsAndEscapes(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Connectors", "Connector Health Review", "Vehicle data setup", "Prediction setup", "Validator setup", "Monitoring export setup", "Feed discovery setup", "Future extension setup", "Setup checklist", "keep_send_enabled=false", "keep_network_send=false", "Connector Catalog", "CSV replay adapter", "HTTP polling adapter", "Webhook sidecar adapter", "Generic JSON transform adapter", "TheTransitClock candidate notes", "Consumer packet preparedness", "No arbitrary dynamic backend plugin loading", "Safe plugin definition", "optional sidecar, command adapter, manifest, or connector process", "not arbitrary dynamic code loaded into the backend", "Vehicle / GPS / AVL connectors", "Prediction connectors", "Validator connectors", "Monitoring / export connectors", "Consumer / discovery connectors", "Future connector extension model", "Manifest Registry", "Synthetic telemetry HTTP poller", "Synthetic telemetry webhook sidecar", "Synthetic predictor sidecar stub", "Synthetic monitoring export", "Synthetic consumer discovery metadata", "disabled by default", "fail closed", "synthetic cases"} {
		if !strings.Contains(body, want) {
			t.Fatalf("html body missing %q: %s", want, body)
		}
	}
	if strings.Contains(body, `<script>alert("x")</script>`) {
		t.Fatalf("html did not escape script-like metadata: %s", body)
	}
	for _, forbidden := range []string{`<form`, `method="post"`, "/public/operations/connectors", "agency approved", "consumer accepted", "production ready", "launch complete", "compliance achieved", "vendor compatible", "certified hardware", "marketplace"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("connector hub html contains forbidden %q: %s", forbidden, body)
		}
	}
	assertConnectorHubSafeStrings(t, body)
}

func TestConnectorWorkbenchRoutesPrivateScopedGETOnlyNoStore(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		t.Run(string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			for _, path := range []string{"/admin/operations/connectors/workbench", "/admin/operations/connectors/workbench.json"} {
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
	for _, path := range []string{"/admin/operations/connectors/workbench", "/admin/operations/connectors/workbench.json"} {
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
		for _, path := range []string{"/admin/operations/connectors/workbench", "/admin/operations/connectors/workbench.json"} {
			req := httptest.NewRequest(method, path, nil)
			rr := httptest.NewRecorder()
			authenticated.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status = %d, want 405", method, path, rr.Code)
			}
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/workbench?agency_id=other-agency", nil)
	rr := httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict html status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/workbench.json?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict json status = %d, want 403", rr.Code)
	}
	for _, path := range []string{"/public/operations/connectors/workbench", "/public/operations/connectors/workbench.json"} {
		req = httptest.NewRequest(http.MethodGet, path, nil)
		rr = httptest.NewRecorder()
		authenticated.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("public connector workbench route %s status = %d, want 404", path, rr.Code)
		}
	}
}

func TestConnectorWorkbenchJSONShapeFlagsRecipesAndManifestReview(t *testing.T) {
	srv := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/workbench.json?agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json prefix", got)
	}
	var view connectorWorkbenchView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode connector workbench: %v", err)
	}
	assertConnectorWorkbenchShape(t, view)
	assertConnectorWorkbenchFlagsFalse(t, view.ClaimFlags)
	assertConnectorWorkbenchSafeStrings(t, rr.Body.String())
	if view.AgencyID != "demo-agency" {
		t.Fatalf("agency_id = %q, want demo-agency", view.AgencyID)
	}
	wantRecipes := []string{"csv_telemetry_sandbox", "api_polling_recipe", "webhook_transform_boundary", "synthetic_only", "predictor_sidecar", "monitoring_export", "public_feed_url_verification", "consumer_discovery_metadata"}
	var gotRecipes []string
	for _, recipe := range view.Recipes {
		gotRecipes = append(gotRecipes, recipe.ID)
	}
	if strings.Join(gotRecipes, ",") != strings.Join(wantRecipes, ",") {
		t.Fatalf("recipe ids = %v, want %v", gotRecipes, wantRecipes)
	}
	wantDecisionRows := []string{"csv_vehicle_locations", "gps_polling_api", "avl_can_post", "synthetic_only", "prediction_sidecar", "monitoring_export", "off_host_validation", "consumer_discovery_metadata"}
	var gotDecisionRows []string
	for _, row := range view.DecisionTree {
		gotDecisionRows = append(gotDecisionRows, row.ID)
	}
	if strings.Join(gotDecisionRows, ",") != strings.Join(wantDecisionRows, ",") {
		t.Fatalf("decision row ids = %v, want %v", gotDecisionRows, wantDecisionRows)
	}
	wantTemplates := []string{"telemetry_source", "prediction_sidecar", "validator_off_host", "monitoring_export", "consumer_discovery"}
	var gotTemplates []string
	for _, row := range view.RedactionTemplates {
		gotTemplates = append(gotTemplates, row.ID)
	}
	if strings.Join(gotTemplates, ",") != strings.Join(wantTemplates, ",") {
		t.Fatalf("redaction template ids = %v, want %v", gotTemplates, wantTemplates)
	}
	wantLintChecks := []string{"secret_and_endpoint_scan", "command_and_plugin_boundary", "status_submission_no_send", "claim_boundary", "fixture_scope"}
	var gotLintChecks []string
	for _, row := range view.ManifestReview.LintChecks {
		gotLintChecks = append(gotLintChecks, row.ID)
	}
	if strings.Join(gotLintChecks, ",") != strings.Join(wantLintChecks, ",") {
		t.Fatalf("manifest lint ids = %v, want %v", gotLintChecks, wantLintChecks)
	}
	var registryIDs []string
	for _, row := range view.ManifestReview.Rows {
		registryIDs = append(registryIDs, row.ConnectorID)
		if row.SourcePath == "" || !strings.HasPrefix(row.SourcePath, "examples/connectors/") || row.FirstCheck != "make external-connection-check" {
			t.Fatalf("manifest row has unsafe source/check: %+v", row)
		}
		if !row.DisabledByDefault || !row.FailClosed || len(row.InputContracts) == 0 || len(row.OutputContracts) == 0 || row.ConformanceCaseCount == 0 {
			t.Fatalf("manifest row must remain disabled, fail-closed, and conformance-backed: %+v", row)
		}
	}
	wantRegistryIDs := []string{"example.consumer-discovery-metadata", "example.generic-json-transform", "example.monitoring-export", "example.predictor-sidecar-stub", "example.telemetry-csv-replay", "example.telemetry-http-poller", "example.telemetry-webhook-sidecar", "example.validator-allowlist"}
	if strings.Join(registryIDs, ",") != strings.Join(wantRegistryIDs, ",") {
		t.Fatalf("manifest ids = %v, want %v", registryIDs, wantRegistryIDs)
	}
}

func TestConnectorWorkbenchHTMLBoundariesNoFormsAndEscapes(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/workbench", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Connector Workbench", "Connection Decision Tree", "csv_vehicle_locations", "consumer_discovery_metadata", "Redaction-First Templates", "Telemetry Source Template", "Consumer Discovery Template", "send_enabled=false", "public_mutation=false", "submit_enabled=false", "Recipe Chooser", "Dry-Run Command Cards", "Synthetic Telemetry Normalization Preview", "Webhook And AVL Transform Boundaries", "Prediction Sidecar Guide", "external_http_shadow", "Vehicle Positions stay independent", "Monitoring Export Guide", "no_send_export_batch", "Consumer / Discovery Guide", "prepared_packet_review", "no_submit_no_status_mutation", "network_send=false", "Synthetic Conformance Viewer", "adapter-conformance-full", "telemetry-malformed", "telemetry-missing-required-field", "prediction-timeout", "prediction-public-mutation-attempt", "validator-allowlist", "validator-raw-command", "monitoring-no-send", "monitoring-unredacted-destination", "consumer-discovery-feed-url-metadata", "consumer-discovery-status-mutation-blocked", "Receiver is deployment-owned", "Transform before telemetry ingest", "Credentials stay server-owned", "Review before any intentional send", "I have a CSV of vehicle locations", "I have a GPS API", "I have an AVL source that can POST", "I want synthetic telemetry only", "I want an external predictor", "I want monitoring summaries", "I want off-host validation", "I want feed discovery metadata", "Example Manifest Registry Review", "Manifest Lint Summary", "Positive claim allowlist", "Safe plugin definition", "Synthetic telemetry CSV replay", "Synthetic telemetry webhook sidecar", "Synthetic consumer discovery metadata", "device-low", "low quality", "network send enabled: false", "disabled by default", "fail closed", "does not upload manifests"} {
		if !strings.Contains(body, want) {
			t.Fatalf("html body missing %q: %s", want, body)
		}
	}
	if strings.Contains(body, `<script>alert("x")</script>`) {
		t.Fatalf("html did not escape script-like metadata: %s", body)
	}
	for _, forbidden := range []string{`<form`, `method="post"`, "/public/operations/connectors", "agency approved", "consumer accepted", "production ready", "launch complete", "compliance achieved", "vendor compatible", "certified hardware", "marketplace", "test connection", "start sidecar"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("connector workbench html contains forbidden %q: %s", forbidden, body)
		}
	}
	assertConnectorWorkbenchSafeStrings(t, body)
}

func TestConnectorTestsRoutesPrivateScopedGETOnlyNoStore(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		t.Run(string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			for _, path := range []string{"/admin/operations/connectors/tests", "/admin/operations/connectors/tests.json"} {
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
	for _, path := range []string{"/admin/operations/connectors/tests", "/admin/operations/connectors/tests.json"} {
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
		for _, path := range []string{"/admin/operations/connectors/tests", "/admin/operations/connectors/tests.json"} {
			req := httptest.NewRequest(method, path, nil)
			rr := httptest.NewRecorder()
			authenticated.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status = %d, want 405", method, path, rr.Code)
			}
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/tests?agency_id=other-agency", nil)
	rr := httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict html status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/tests.json?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict json status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/connectors/tests.json", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public connector test route status = %d, want 404", rr.Code)
	}
}

func TestConnectorTestsJSONShapeFlagsAndCommands(t *testing.T) {
	srv := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/tests.json?agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json prefix", got)
	}
	var view connectorTestsView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode connector tests: %v", err)
	}
	assertConnectorTestsShape(t, view)
	assertConnectorTestsFlagsFalse(t, view.ClaimFlags)
	assertConnectorTestsSafeStrings(t, rr.Body.String())
	wantCommands := []string{
		"make external-connection-check",
		"make adapter-conformance",
		"go run ./cmd/adapter-conformance manifest --suite testdata/adapter-conformance",
		"go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance",
		"go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance",
		"go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance",
		"go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance",
		"go run ./cmd/adapter-conformance consumer_discovery --suite testdata/adapter-conformance",
		"make test-connector-examples",
	}
	var gotCommands []string
	for _, command := range view.Commands {
		gotCommands = append(gotCommands, command.CommandLine)
	}
	if strings.Join(gotCommands, "\n") != strings.Join(wantCommands, "\n") {
		t.Fatalf("commands = %v, want %v", gotCommands, wantCommands)
	}
}

func TestConnectorTestsHTMLInstructionsOnly(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/tests", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Connector Test Instructions", "make external-connection-check", "make adapter-conformance", "make test-connector-examples", "Telemetry connector cases", "Prediction connector cases", "Validator connector cases", "Monitoring/export connector cases", "Consumer/discovery connector cases", "consumer_discovery", "does not execute commands", "read manifest-provided commands"} {
		if !strings.Contains(body, want) {
			t.Fatalf("html body missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{`<form`, `method="post"`, "/public/operations/connectors", "agency approved", "consumer accepted", "production ready", "launch complete", "compliance achieved", "vendor compatible", "certified hardware", "marketplace"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("connector tests html contains forbidden %q: %s", forbidden, body)
		}
	}
	assertConnectorTestsSafeStrings(t, body)
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

func TestOperationsConsoleNavigationIsGroupedAndRouteStable(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		`aria-label="Operations Console sections"`,
		"Start",
		`href="/admin/operations" aria-current="page">Start</a>`,
		`href="/admin/operations/devices">Devices</a>`,
		`href="/admin/operations/telemetry-simulator">Simulator</a>`,
		"Schedule Review",
		"Realtime",
		`href="/admin/operations/prediction-lab">Trip Updates</a>`,
		"Connectors",
		"Feeds",
		"Maintain",
		"Help",
		`href="/admin/operations" aria-current="page"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("navigation missing %q: %s", want, body)
		}
	}
	for _, href := range []string{
		"/admin/operations",
		"/admin/operations/launchpad",
		"/admin/operations/setup-wizard",
		"/admin/operations/connectors",
		"/admin/operations/connectors/workbench",
		"/admin/operations/connectors/tests",
		"/admin/operations/gtfs-workbench",
		"/admin/operations/gtfs-import",
		"/admin/operations/feed-health",
		"/admin/operations/prediction-lab",
		"/admin/operations/readiness",
		"/admin/operations/feeds",
		"/admin/operations/gtfs-quality",
		"/admin/operations/validation-health",
		"/admin/operations/reliability",
		"/admin/operations/maintenance",
		"/admin/operations/access",
		"/admin/operations/audit",
		"/admin/operations/telemetry",
		"/admin/operations/telemetry-simulator",
		"/admin/operations/devices",
		"/admin/alerts/console",
		"/admin/operations/help",
		"/admin/operations/consumers",
		"/admin/operations/evidence",
		"/admin/operations/setup",
		"/admin/operations/checklist",
		"/admin/gtfs-studio",
	} {
		if !strings.Contains(body, `href="`+href+`"`) {
			t.Fatalf("navigation missing stable href %q: %s", href, body)
		}
	}
	if got := strings.Count(body, `aria-current="page"`); got != 1 {
		t.Fatalf("aria-current count = %d, want 1: %s", got, body)
	}
	for _, oldLabel := range []string{">Dashboard</a>", ">Start Here</a>", ">Devices &amp; Tokens</a>", ">Telemetry Simulator</a>"} {
		if strings.Contains(body, oldLabel) {
			t.Fatalf("navigation still contains old label %q: %s", oldLabel, body)
		}
	}
	for _, forbidden := range []string{"agency approved", "consumer accepted", "production ready", "launch complete", "compliance achieved", "vendor compatible", "certified hardware"} {
		if strings.Contains(strings.ToLower(body), forbidden) {
			t.Fatalf("navigation body contains forbidden claim %q: %s", forbidden, body)
		}
	}
}

func TestOperationsRouteRegistryCentralizesCanonicalInventory(t *testing.T) {
	groups := make(map[string]bool)
	for _, group := range operationsRouteGroups() {
		if group.ID == "" || group.Label == "" {
			t.Fatalf("route group has empty field: %+v", group)
		}
		if groups[group.ID] {
			t.Fatalf("duplicate route group: %s", group.ID)
		}
		groups[group.ID] = true
	}

	sections := make(map[string]bool)
	paths := make(map[string]bool)
	for _, route := range operationsRoutes() {
		if route.Section == "" || route.Path == "" || route.NavLabel == "" || route.GroupID == "" {
			t.Fatalf("route registry row has empty required field: %+v", route)
		}
		if sections[route.Section] {
			t.Fatalf("duplicate route section: %s", route.Section)
		}
		if paths[route.Path] {
			t.Fatalf("duplicate route path: %s", route.Path)
		}
		if !groups[route.GroupID] {
			t.Fatalf("route %s references unknown group %q", route.Path, route.GroupID)
		}
		if strings.HasPrefix(route.Path, "/public/") {
			t.Fatalf("route registry must not include public admin route: %+v", route)
		}
		if len(route.Methods) == 0 {
			t.Fatalf("route %s has no method posture", route.Path)
		}
		if route.ExternalAdminSurface {
			if route.JSONPath != "" {
				t.Fatalf("external admin surface should not define Operations JSON pair: %+v", route)
			}
		} else {
			if !strings.HasPrefix(route.Path, "/admin/operations") {
				t.Fatalf("canonical Operations route has unexpected path: %+v", route)
			}
			if route.PageTitle == "" {
				t.Fatalf("canonical Operations route missing page title: %+v", route)
			}
			if !route.NoStore {
				t.Fatalf("canonical Operations route must be no-store: %+v", route)
			}
			if operationsPageTitle(route.Section) != route.PageTitle {
				t.Fatalf("route %s page title mismatch: got %q want %q", route.Section, operationsPageTitle(route.Section), route.PageTitle)
			}
		}
		if _, ok := operationsRouteBySection(route.Section); !ok {
			t.Fatalf("route section lookup failed for %s", route.Section)
		}
		if _, ok := operationsRouteByPath(route.Path); !ok {
			t.Fatalf("route path lookup failed for %s", route.Path)
		}
		if route.JSONPath != "" {
			if !strings.HasSuffix(route.JSONPath, ".json") {
				t.Fatalf("JSON route lacks .json suffix: %+v", route)
			}
			if _, ok := operationsRouteByPath(route.JSONPath); !ok {
				t.Fatalf("JSON route lookup failed for %s", route.JSONPath)
			}
		}
		sections[route.Section] = true
		paths[route.Path] = true
	}

	if got := len(operationsCanonicalHTMLRoutes()); got != 28 {
		t.Fatalf("canonical HTML route count = %d, want 28", got)
	}
	jsonRoutes := operationsCanonicalJSONRoutes()
	if got := len(jsonRoutes); got != 20 {
		t.Fatalf("canonical JSON route count = %d, want 20: %v", got, jsonRoutes)
	}
	if !containsString(jsonRoutes, "/admin/operations/checklist.json") {
		t.Fatalf("registry must include checklist JSON route: %v", jsonRoutes)
	}
	commands := operationsCommandRoutes()
	if len(commands) != 1 || commands[0].Path != "/admin/operations/validation-health/refresh.json" || commands[0].Method != http.MethodPost || !commands[0].NoStore {
		t.Fatalf("command route registry mismatch: %+v", commands)
	}
	externals := operationsExternalAdminSurfaceRoutes()
	externalPaths := make([]string, 0, len(externals))
	for _, route := range externals {
		externalPaths = append(externalPaths, route.Path)
	}
	for _, want := range []string{"/admin/gtfs-studio", "/admin/alerts/console"} {
		if !containsString(externalPaths, want) {
			t.Fatalf("external admin surface missing %q: %+v", want, externals)
		}
	}
	if title := operationsPageTitle("not-a-real-section"); title != "Operations Console" {
		t.Fatalf("unknown section title = %q, want Operations Console", title)
	}
}

func TestOperationsRouteTitlesAndFirstClickLabelOrder(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	for _, tc := range []struct {
		path  string
		title string
	}{
		{path: "/admin/operations", title: "Start"},
		{path: "/admin/operations/gtfs-workbench", title: "Schedule Review"},
		{path: "/admin/operations/gtfs-import", title: "Import Schedule"},
		{path: "/admin/operations/prediction-lab", title: "Trip Updates"},
		{path: "/admin/operations/telemetry", title: "Telemetry"},
		{path: "/admin/operations/devices", title: "Devices"},
		{path: "/admin/operations/connectors/workbench", title: "Connector Workbench"},
		{path: "/admin/operations/access", title: "Access &amp; Roles"},
		{path: "/admin/operations/audit", title: "Audit Log"},
		{path: "/admin/operations/help", title: "Help"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
			}
			body := rr.Body.String()
			for _, want := range []string{"<title>" + tc.title + "</title>", `<h1 id="operations-page-title">` + tc.title + "</h1>"} {
				if !strings.Contains(body, want) {
					t.Fatalf("%s missing page-specific title %q: %s", tc.path, want, body)
				}
			}
			if tc.path == "/admin/operations" {
				firstClick := strings.Index(body, `<h1 id="operations-page-title">Start</h1>`)
				helpPanel := strings.Index(body, `class="context-help"`)
				if firstClick < 0 || helpPanel < 0 || firstClick > helpPanel {
					t.Fatalf("Start label should appear before contextual help: firstClick=%d helpPanel=%d body=%s", firstClick, helpPanel, body)
				}
			}
		})
	}
}

func TestOperationsConsoleNavigationActiveStateForRepresentativeSections(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	for _, tc := range []struct {
		path string
		href string
	}{
		{path: "/admin/operations/gtfs-workbench", href: "/admin/operations/gtfs-workbench"},
		{path: "/admin/operations/feed-health", href: "/admin/operations/feed-health"},
		{path: "/admin/operations/gtfs-quality", href: "/admin/operations/gtfs-quality"},
		{path: "/admin/operations/prediction-lab", href: "/admin/operations/prediction-lab"},
		{path: "/admin/operations/telemetry", href: "/admin/operations/telemetry"},
		{path: "/admin/operations/connectors/workbench", href: "/admin/operations/connectors/workbench"},
		{path: "/admin/operations/connectors/tests", href: "/admin/operations/connectors/tests"},
		{path: "/admin/operations/help", href: "/admin/operations/help"},
		{path: "/admin/operations/consumers", href: "/admin/operations/consumers"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
			}
			body := rr.Body.String()
			if !strings.Contains(body, `href="`+tc.href+`" aria-current="page"`) {
				t.Fatalf("%s did not mark current nav item %q: %s", tc.path, tc.href, body)
			}
			if got := strings.Count(body, `aria-current="page"`); got != 1 {
				t.Fatalf("%s aria-current count = %d, want 1: %s", tc.path, got, body)
			}
		})
	}
}

func TestOperationsLegacyPrivatePagesUseNoStore(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	for _, path := range []string{
		"/admin/operations/feeds",
		"/admin/operations/telemetry",
		"/admin/operations/devices",
		"/admin/operations/consumers",
		"/admin/operations/evidence",
		"/admin/operations/setup",
	} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
			}
			if got := rr.Header().Get("Cache-Control"); got != "no-store" {
				t.Fatalf("%s Cache-Control = %q, want no-store", path, got)
			}
		})
	}
}

func TestOperationsHelpRoutesPrivateScopedGETOnlyNoStore(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		t.Run(string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			for _, path := range []string{"/admin/operations/help", "/admin/operations/help.json", "/admin/operations/help/", "/admin/operations/help.json/"} {
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

	unauth := newOperationsTestHandler(&handler{}, authRejectAll{})
	for _, path := range []string{"/admin/operations/help", "/admin/operations/help.json", "/admin/operations/help/", "/admin/operations/help.json/"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		unauth.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("unauth %s status = %d, want 401", path, rr.Code)
		}
		if got := rr.Header().Get("Cache-Control"); got != "no-store" {
			t.Fatalf("unauth %s Cache-Control = %q, want no-store", path, got)
		}
	}

	authenticated := newOperationsTestHandler(&handler{}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator}, Method: auth.MethodBearer,
	}})
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		for _, path := range []string{"/admin/operations/help", "/admin/operations/help.json", "/admin/operations/help/", "/admin/operations/help.json/"} {
			req := httptest.NewRequest(method, path, nil)
			rr := httptest.NewRecorder()
			authenticated.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status = %d, want 405", method, path, rr.Code)
			}
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/help?agency_id=other-agency", nil)
	rr := httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict html status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/help.json?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict json status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/help.json", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public help route status = %d, want 404", rr.Code)
	}
}

func TestOperationsHelpJSONShapeFlagsAndNoLeakage(t *testing.T) {
	handler := newOperationsTestHandler(&handler{}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/help.json?agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json prefix", got)
	}
	var view operationsHelpView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode help JSON: %v", err)
	}
	assertOperationsHelpShape(t, view)
	assertOperationsHelpFlagsFalse(t, view.ClaimFlags)
	assertOperationsHelpSafeStrings(t, rr.Body.String())
	if view.AgencyID != "demo-agency" {
		t.Fatalf("agency_id = %q, want demo-agency", view.AgencyID)
	}
	if len(view.ContextualHelp.Topics) == 0 || view.ContextualHelp.AllTopicsURL != "/admin/operations/help" || view.ContextualHelp.JSONURL != "/admin/operations/help.json" {
		t.Fatalf("invalid contextual help: %+v", view.ContextualHelp)
	}
}

func TestOperationsHelpHTMLRendersTopicsBoundariesAndNoForms(t *testing.T) {
	handler := newOperationsTestHandler(&handler{}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/help", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Help &amp; Tutorials",
		`id="help-gtfs"`,
		`id="help-gtfs_rt"`,
		`id="help-connectors"`,
		`id="help-readiness"`,
		`id="help-validators"`,
		`id="help-telemetry"`,
		`id="help-claims_evidence"`,
		`id="help-role-no_code_evaluator"`,
		`id="help-role-director_manager"`,
		`id="help-role-daily_operator"`,
		`id="help-role-administrator"`,
		`id="help-role-integrator"`,
		"Role-Based Tours",
		"First-Week Checklist",
		"Plain-Language Glossary",
		"Common Mistake Recovery",
		"Printable Staff Training Guide",
		"Quick Tasks",
		"Staff Handoff Checklist",
		"Demo Scenario Catalog",
		"Trainer Script",
		"Administrator Checklist",
		"docs/operator-training-guide.md",
		"No-code evaluator",
		"Director or manager",
		"Daily operator",
		"Administrator",
		"Integrator",
		`id="help-first-week-day_1_start"`,
		`id="help-first-week-day_5_maintenance"`,
		`id="help-glossary-gtfs"`,
		`id="help-glossary-trip_updates"`,
		`id="help-glossary-support_bundle"`,
		`id="help-recovery-validator_blocked"`,
		`id="help-recovery-consumer_status_confusion"`,
		`id="help-quick-task-import_gtfs"`,
		`id="help-quick-task-support_bundle"`,
		`id="help-handoff-setup_owner"`,
		`id="help-handoff-claim_owner"`,
		`id="help-demo-scenario-baseline_start"`,
		`id="help-demo-scenario-after_midnight"`,
		`id="help-demo-scenario-connector_conformance"`,
		`id="help-trainer-step-open_boundary"`,
		`id="help-trainer-step-closeout"`,
		`id="help-helper-check-startup_context"`,
		`id="help-helper-check-evidence_context"`,
		"testdata/gtfs/valid-small",
		"testdata/adapter-conformance/suite.json",
		"Trip Updates are empty, fallback, or withheld",
		"Prepared packet visibility does not show submission",
		safePluginDefinition,
		"Detailed safety booleans remain available",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("help HTML missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{`<form`, `method="post"`, "/public/operations/help", "agency approved", "consumer accepted", "production ready", "launch complete", "compliance achieved", "vendor compatible", "certified hardware"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("help HTML contains forbidden string %q: %s", forbidden, body)
		}
	}
	assertOperationsHelpSafeStrings(t, body)
}

func TestOperationsSharedLayoutRendersContextualHelpForMajorSections(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	for _, tc := range []struct {
		path string
		want []string
	}{
		{path: "/admin/operations/gtfs-import", want: []string{`Help for GTFS import`, `help-gtfs`, `help-validators`}},
		{path: "/admin/operations/feed-health", want: []string{`Help for feed health`, `help-gtfs_rt`, `help-readiness`}},
		{path: "/admin/operations/connectors", want: []string{`Help for Connectors`, `help-connectors`, `help-claims_evidence`}},
		{path: "/admin/operations/telemetry-simulator", want: []string{`Help for Telemetry Simulator`, `help-telemetry`, `help-gtfs_rt`}},
		{path: "/admin/operations/readiness", want: []string{`Help for readiness`, `help-readiness`, `help-claims_evidence`}},
		{path: "/admin/operations/validation-health", want: []string{`Help for validator health`, `help-validators`, `help-claims_evidence`}},
		{path: "/admin/operations/evidence", want: []string{`Help for evidence`, `help-claims_evidence`, `help-readiness`}},
	} {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
			}
			body := rr.Body.String()
			for _, want := range append([]string{`class="context-help"`, `<strong>Next:</strong>`, `/admin/operations/help`, `/admin/operations/help.json`}, tc.want...) {
				if !strings.Contains(body, want) {
					t.Fatalf("%s contextual help missing %q: %s", tc.path, want, body)
				}
			}
			if strings.Contains(body, `<form method="post" action="/admin/operations/help"`) {
				t.Fatalf("%s contextual help must not add help POST form: %s", tc.path, body)
			}
		})
	}
}

func TestOperationsConsoleEmptyStateGuidanceAnswersFirstRunQuestions(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	paths := []string{
		"/admin/operations/gtfs-import",
		"/admin/operations/feed-health",
		"/admin/operations/gtfs-quality",
		"/admin/operations/validation-health",
		"/admin/operations/devices",
		"/admin/operations/telemetry",
		"/admin/operations/telemetry-simulator",
		"/admin/operations/connectors",
		"/admin/operations/connectors/tests",
		"/admin/operations/maintenance",
		"/admin/operations/help",
	}
	want := []string{
		"What am I seeing?",
		"Is this bad?",
		"What should I do next?",
		"Can I do it in the browser?",
		"When is an administrator needed?",
		"Limits",
	}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
			}
			body := rr.Body.String()
			for _, text := range want {
				if !strings.Contains(body, text) {
					t.Fatalf("%s empty-state guidance missing %q: %s", path, text, body)
				}
			}
			if !strings.Contains(body, `class="card empty-state"`) && !strings.Contains(body, `class="section-note"`) {
				t.Fatalf("%s empty-state guidance missing a stable container class: %s", path, body)
			}
		})
	}
}

func TestOperationsConsoleSharedLayoutHasAccessibilityAndMobileLandmarks(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		`<html lang="en">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<link rel="icon" href="data:,">`,
		`class="skip-link" href="#operations-main"`,
		`class="skip-link" href="#operations-nav"`,
		`<header class="operations-header" role="banner">`,
		`<h1 id="operations-page-title">Start</h1>`,
		`<div class="operations-frame">`,
		`<nav id="operations-nav" class="operations-nav" aria-label="Operations Console sections">`,
		`<section class="nav-group" aria-labelledby="nav-group-maintenance">`,
		`<p id="nav-group-maintenance" class="nav-group-label">Maintain</p>`,
		`<main id="operations-main" tabindex="-1" aria-labelledby="operations-page-title">`,
		`<script src="/admin/operations/assets/operations.js" defer></script>`,
		`</main>
</div>
<script src="/admin/operations/assets/operations.js" defer></script></body></html>`,
		`:focus-visible`,
		`summary:focus-visible`,
		`@media (prefers-contrast:more)`,
		`tbody tr:focus-within`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("shared layout missing %q: %s", want, body)
		}
	}
}

func TestOperationsProgressiveAssetPrivateAllowlistedAndNoBuild(t *testing.T) {
	srv := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/assets/operations.js", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("asset status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	if got := rr.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options = %q, want nosniff", got)
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/javascript") {
		t.Fatalf("Content-Type = %q, want application/javascript", got)
	}
	body := rr.Body.String()
	for _, want := range []string{
		"OpenTransitOperations",
		"safeAdminPath",
		"/admin/operations/validation-health/refresh.json",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("asset missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{
		"http://", "https://", "cdn.", "unpkg", "maps.googleapis", "/public/gtfsrt", "/v1/events",
		"csrf_token.", "bearer_token", "device_token", "Authorization", "Cookie", "postgres://", "raw_report", "stdout", "stderr", "/Users/", "docs/evidence/captured",
	} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("asset contains forbidden private or external pattern %q: %s", forbidden, body)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/public/operations/assets/operations.js", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public asset status = %d, want 404", rr.Code)
	}
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/assets/operations.js", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST asset status = %d, want 405", rr.Code)
	}
	unauth := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, authRejectAll{})
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/assets/operations.js", nil)
	rr = httptest.NewRecorder()
	unauth.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK {
		t.Fatalf("unauthenticated asset returned 200")
	}
}

func TestOperationsProgressiveReviewToolsRenderWithNoJSFallback(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	for _, tc := range []struct {
		path string
		want []string
	}{
		{
			path: "/admin/operations",
			want: []string{
				`data-copy-card`,
				`class="copy-value" data-copy-value=`,
				`Copy Feed URLs`,
			},
		},
		{
			path: "/admin/operations/feed-health",
			want: []string{
				`class="review-tools" data-review-tools data-review-target="feed-health-review-rows"`,
				`data-review-filter`,
				`data-review-search`,
				`data-review-sort`,
				`data-review-remember`,
				`data-review-reset`,
				`id="feed-health-review-rows"`,
				`data-review-row`,
				`Filters only change this browser view`,
			},
		},
		{
			path: "/admin/operations/validation-health",
			want: []string{
				`data-admin-refresh="/admin/operations/validation-health/refresh.json"`,
				`id="validation-refresh-status" class="review-status" aria-live="polite"`,
				`Reloads existing private records only`,
				`does not run validators, change feeds, create evidence, contact consumers, or prove readiness`,
			},
		},
	} {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("%s status = %d, want 200: %s", tc.path, rr.Code, rr.Body.String())
			}
			body := rr.Body.String()
			for _, want := range tc.want {
				if !strings.Contains(body, want) {
					t.Fatalf("%s body missing %q: %s", tc.path, want, body)
				}
			}
			for _, forbidden := range []string{"Copy approved", "Copy production", "Copy compliant", "consumer-ready", "certified"} {
				if strings.Contains(body, forbidden) {
					t.Fatalf("%s body contains unsupported enhancement copy %q: %s", tc.path, forbidden, body)
				}
			}
		})
	}
}

func TestOperationsCoreRoutesUseSharedAppShellAndDesignTokens(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	for _, path := range []string{
		"/admin/operations",
		"/admin/operations/setup-wizard",
		"/admin/operations/gtfs-import",
		"/admin/operations/feed-health",
		"/admin/operations/readiness",
		"/admin/operations/validation-health",
		"/admin/operations/telemetry",
		"/admin/operations/connectors",
		"/admin/operations/maintenance",
		"/admin/operations/help",
	} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
			}
			body := rr.Body.String()
			for _, want := range []string{
				`<header class="operations-header" role="banner">`,
				`Private operations`,
				`class="app-breadcrumb"`,
				`class="app-meta"`,
				`<div class="operations-frame">`,
				`<nav id="operations-nav" class="operations-nav" aria-label="Operations Console sections">`,
				`<main id="operations-main" tabindex="-1" aria-labelledby="operations-page-title">`,
				`:root{`,
				`--color-surface`,
				`--space-4`,
				`--radius-3`,
				`--color-focus`,
				`.status-ready-for-local-review`,
				`@media (prefers-reduced-motion:reduce)`,
				`@media (prefers-contrast:more)`,
			} {
				if !strings.Contains(body, want) {
					t.Fatalf("%s shared shell missing %q: %s", path, want, body)
				}
			}
		})
	}
}

func TestOperationsConsoleDesignSystemAvoidsFragileDecoration(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		`Draft Schedule Editor <span class="nav-surface">separate tool</span>`,
		`Alerts <span class="nav-surface">separate tool</span>`,
		`Ready for local review`,
		`Copy Feed URLs`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("design-system route marker missing %q: %s", want, body)
		}
	}
	consumerReq := httptest.NewRequest(http.MethodGet, "/admin/operations/consumers", nil)
	consumerRR := httptest.NewRecorder()
	handler.ServeHTTP(consumerRR, consumerReq)
	if consumerRR.Code != http.StatusOK {
		t.Fatalf("consumer tracker status = %d, want 200: %s", consumerRR.Code, consumerRR.Body.String())
	}
	if !strings.Contains(consumerRR.Body.String(), "Prepared Consumer Packet Tracker") {
		t.Fatalf("consumer page missing prepared packet heading: %s", consumerRR.Body.String())
	}
	lower := strings.ToLower(body)
	for _, forbidden := range []string{
		"border-left:",
		"background-clip:text",
		"letter-spacing:-",
		"consumer submission evidence",
		"copy these five feed urls",
	} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("design system contains fragile or overclaiming marker %q: %s", forbidden, body)
		}
	}
}

func TestOperationsConsoleMobileCSSKeepsTablesFormsAndLongValuesUsable(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		`*{box-sizing:border-box}`,
		`overflow-wrap:anywhere`,
		`@media (max-width:700px)`,
		`table{display:block;max-width:100%;overflow-x:auto;-webkit-overflow-scrolling:touch}`,
		`input,select,textarea{min-width:0;width:100%}`,
		`.nav-link,button{width:100%}`,
		`a:focus-visible,button:focus-visible,input:focus-visible,select:focus-visible,textarea:focus-visible,summary:focus-visible,main:focus-visible`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("mobile/focus CSS missing %q: %s", want, body)
		}
	}
}

func TestOperationsConsoleFormsUseLabelsAndSubmitButtonsWithoutChangingContracts(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	for _, tc := range []struct {
		path string
		want []string
	}{
		{
			path: "/admin/operations/gtfs-import",
			want: []string{
				`<form method="post" enctype="multipart/form-data" action="/admin/operations/gtfs-import?csrf_token=`,
				`name="action" value="import_gtfs"`,
				`name="source_type" value="upload"`,
				`for="gtfs_upload_zip"`,
				`id="gtfs_upload_zip" type="file" name="gtfs_zip"`,
				`for="gtfs_upload_notes"`,
				`id="gtfs_upload_notes" name="notes"`,
				`<button type="submit">Import ZIP</button>`,
				`<form method="post" action="/admin/operations/gtfs-import">`,
				`name="source_type" value="url"`,
				`for="gtfs_import_url"`,
				`id="gtfs_import_url" type="url" name="gtfs_url"`,
				`for="gtfs_url_notes"`,
				`id="gtfs_url_notes" name="notes"`,
				`<button type="submit">Import URL</button>`,
			},
		},
		{
			path: "/admin/operations/gtfs-quality",
			want: []string{
				`<form method="post" action="/admin/operations/gtfs-quality">`,
				`name="action" value="rerun_static_validator"`,
				`<button type="submit">Rerun static MobilityData validator</button>`,
			},
		},
		{
			path: "/admin/operations/validation-health",
			want: []string{
				`<form method="post" action="/admin/operations/validation-health">`,
				`name="action" value="run_all"`,
				`<button type="submit">Run allowlisted validators</button>`,
			},
		},
		{
			path: "/admin/operations/devices",
			want: []string{
				`<form method="post" action="/admin/operations/devices">`,
				`name="agency_id" value="demo-agency"`,
				`for="device_rebind_device_id"`,
				`id="device_rebind_device_id" name="device_id"`,
				`for="device_rebind_vehicle_id"`,
				`id="device_rebind_vehicle_id" name="vehicle_id"`,
				`for="device_rebind_reason"`,
				`id="device_rebind_reason" name="reason"`,
				`<button type="submit">Rotate / rebind token</button>`,
			},
		},
		{
			path: "/admin/operations/setup",
			want: []string{
				`<form method="post" action="/admin/operations/setup">`,
				`name="action" value="publication_bootstrap"`,
				`for="setup_public_base_url"`,
				`id="setup_public_base_url" type="url" name="public_base_url"`,
				`for="setup_feed_base_url"`,
				`id="setup_feed_base_url" type="url" name="feed_base_url"`,
				`for="setup_technical_contact_email"`,
				`id="setup_technical_contact_email" type="email" name="technical_contact_email"`,
				`for="setup_license_name"`,
				`id="setup_license_name" name="license_name"`,
				`for="setup_license_url"`,
				`id="setup_license_url" type="url" name="license_url"`,
				`for="setup_publication_environment"`,
				`id="setup_publication_environment" name="publication_environment"`,
				`<button type="submit">Store publication metadata</button>`,
				`name="action" value="run_validation"`,
				`for="setup_validation_feed_type"`,
				`id="setup_validation_feed_type" name="feed_type"`,
				`<button type="submit">Run allowlisted validation</button>`,
			},
		},
	} {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
			}
			body := rr.Body.String()
			for _, want := range tc.want {
				if !strings.Contains(body, want) {
					t.Fatalf("%s missing form contract %q: %s", tc.path, want, body)
				}
			}
		})
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
	for _, want := range []string{"record_small_host_resources", "record_service_dependency_review", "record_postgres_capacity_review", "record_upgrade_rollback_review", "RESTORE_DRILL_DATABASE_URL", "hosted_saas_claimed"} {
		if !strings.Contains(text, want) {
			t.Fatalf("deployment doctor missing Phase 104 guard %q", want)
		}
	}
	ociCheck, err := os.ReadFile(filepath.Join("..", "..", "scripts", "oci-reference-check.sh"))
	if err != nil {
		t.Fatalf("read oci reference check: %v", err)
	}
	if !strings.Contains(string(ociCheck), "/healthz") || strings.Contains(string(ociCheck), "127.0.0.1:${port}/health\"") {
		t.Fatalf("oci reference check must probe /healthz loopback health")
	}
	caddy, err := os.ReadFile(filepath.Join("..", "..", "deploy", "Caddyfile.local"))
	if err != nil {
		t.Fatalf("read caddyfile: %v", err)
	}
	caddyText := string(caddy)
	caddyLines := nonCommentCaddyLines(caddyText)
	if !containsString(caddyLines, "@local_root {") || !containsString(caddyLines, "path /") || !containsString(caddyLines, `respond @local_root "Open Transit RT local app is running. Start local browser setup at /admin/local-login. Public feeds are under /public/ and admin routes require auth." 200`) {
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

func TestGTFSQualityGuidanceShowsActionableFixPathsSafely(t *testing.T) {
	now := time.Now().UTC()
	store := &fakePublicationStore{discovery: gtfsQualityDiscovery(now), validationRecords: []compliance.ValidationReportRecord{
		{
			ID: 1, CreatedAt: now, Result: compliance.ValidationResult{
				AgencyID: "demo-agency", FeedType: "schedule", FeedVersionID: "feed-v1", ValidatorName: compliance.CanonicalStaticValidatorName, Status: "failed", ErrorCount: 1, WarningCount: 2,
				Report: map[string]any{"raw_report": map[string]any{"notices": []any{
					map[string]any{"code": "expired_calendar", "severity": "ERROR", "message": "expired", "path": "/tmp/private/calendar.txt"},
					map[string]any{"code": "stop_times_arrival_time_missing", "severity": "WARNING", "message": "bad stop time"},
					map[string]any{"code": "frequency_headway_invalid", "severity": "WARNING", "message": "bad frequency"},
				}}},
			},
		},
		{
			ID: 2, CreatedAt: now, Result: compliance.ValidationResult{
				AgencyID: "demo-agency", FeedType: "schedule", FeedVersionID: "feed-v1", ValidatorName: compliance.InternalGTFSImportValidatorName, Status: "failed", ErrorCount: 1,
				Report: map[string]any{"errors": []any{map[string]any{"code": "missing_trip_reference", "message": "missing trip reference", "path": "/tmp/private/trips.txt"}}},
			},
		},
	}}
	handler := newGTFSQualityTestHandler(t, auth.RoleAdmin, store, fakeScheduleBuilder{})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/gtfs-quality", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Fix Workflow",
		"Fix Planner",
		"Private Fix Checklist",
		"manual_review_only",
		"Likely owner",
		"Affected files",
		"Safe fix path",
		"Risk level",
		"can break service availability or realtime usefulness",
		"Safe draft suggestion",
		"Draft suggestion record",
		"Before validation plan",
		"After validation plan",
		"Verify with",
		"Escalate if",
		"Schedule planner or GTFS source owner",
		"calendar.txt / calendar_dates.txt",
		"stop_times.txt / trips.txt",
		"frequencies.txt / trips.txt",
		"GTFS export owner or source-system admin",
		"Re-import through browser or CLI",
		"does not edit GTFS",
		"No automatic production edit",
		"Advisory only; no persisted draft suggestion record",
		"Detailed safety booleans remain available",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("GTFS quality guidance missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{
		"/tmp/private",
		"raw_report",
		"stdout",
		"stderr",
		"argv",
		"validator_clean",
		"consumer_ready",
		"production_ready",
		"automatic_gtfs_edit_enabled</code></th><td>true",
		"draft_mutation_enabled</code></th><td>true",
		"draft_suggestion_records_created</code></th><td>true",
		"schedule_publish_enabled</code></th><td>true",
		"validator_semantics_changed</code></th><td>true",
		"compliance_claimed</code></th><td>true",
		"consumer_statuses_changed</code></th><td>true",
	} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("GTFS quality guidance leaked or overclaimed %q: %s", forbidden, body)
		}
	}
}

func TestGTFSQualityFixPlannerBoundsRowsAndNoMutationFlags(t *testing.T) {
	groups := make([]compliance.GTFSQualityGroup, 0, gtfsQualityFixPlannerMaxRows+5)
	for i := 0; i < gtfsQualityFixPlannerMaxRows+5; i++ {
		groups = append(groups, compliance.GTFSQualityGroup{
			Source:            compliance.GTFSQualitySourceCanonicalValidator,
			Family:            fmt.Sprintf("unknown_%02d", i),
			Codes:             []string{fmt.Sprintf("code_%02d", i)},
			Severity:          compliance.GTFSQualityNeedsReview,
			Count:             1,
			OperatorSummary:   "planner row",
			WhyItMatters:      "bounded planner rows keep the page usable",
			RecommendedAction: "review source GTFS",
			Samples:           []string{"code=sample"},
		})
	}
	page := operationsPage{
		AgencyID:          "demo-agency",
		ActiveFeedVersion: "feed-v1",
		GTFSQuality: compliance.GTFSQualityTriage{
			Canonical: compliance.GTFSQualitySection{
				Source:      compliance.GTFSQualitySourceCanonicalValidator,
				SourceLabel: "Canonical MobilityData static validator",
				Status:      compliance.GTFSQualityNeedsReview,
				Groups:      groups,
			},
			InternalImporter: compliance.GTFSQualitySection{Status: compliance.GTFSQualityInformational},
		},
	}
	guidance := buildOperationsGTFSQualityGuidance(page)
	if guidance.FixPlanner.TotalRows != gtfsQualityFixPlannerMaxRows+5 || guidance.FixPlanner.DisplayedRows != gtfsQualityFixPlannerMaxRows || guidance.FixPlanner.HiddenRows != 5 {
		t.Fatalf("planner bounds = total %d displayed %d hidden %d", guidance.FixPlanner.TotalRows, guidance.FixPlanner.DisplayedRows, guidance.FixPlanner.HiddenRows)
	}
	if !strings.Contains(guidance.FixPlanner.Checklist, "5 additional grouped issue row(s) are hidden") {
		t.Fatalf("checklist missing hidden-row notice: %s", guidance.FixPlanner.Checklist)
	}
	if guidance.ClaimFlags.AutomaticGTFSEditEnabled || guidance.ClaimFlags.DraftMutationEnabled || guidance.ClaimFlags.DraftSuggestionRecordsCreated || guidance.ClaimFlags.SchedulePublishEnabled || guidance.ClaimFlags.ConsumerStatusesChanged || guidance.ClaimFlags.ComplianceClaimed {
		t.Fatalf("unexpected mutation or claim flag: %+v", guidance.ClaimFlags)
	}
	for _, row := range guidance.FixPlanner.Rows {
		if row.RiskLevel == "" || !strings.Contains(row.NoAutoApplyBoundary, "No automatic production edit") || !strings.Contains(row.DraftSuggestionRecord, "no persisted draft suggestion record") {
			t.Fatalf("planner row missing safe boundary: %+v", row)
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

func TestValidationHealthRefreshCommandJSONPrivateScopedAndSafe(t *testing.T) {
	now := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC)
	store := &fakePublicationStore{discovery: validationHealthTestDiscovery(now), validationRecords: []compliance.ValidationReportRecord{
		validationHealthRecord(1, "schedule", "feed-v1", "passed", now),
		validationHealthRecord(2, "vehicle_positions", "feed-v1", "warning", now),
	}}
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		handler := newValidationHealthTestHandler(t, role, store, fakeScheduleBuilder{})
		req := httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health/refresh.json", strings.NewReader("action=validation_health.refresh"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("role %s status = %d, want 200: %s", role, rr.Code, rr.Body.String())
		}
		if got := rr.Header().Get("Cache-Control"); got != "no-store" {
			t.Fatalf("role %s Cache-Control = %q, want no-store", role, got)
		}
		if !strings.HasPrefix(rr.Header().Get("Content-Type"), "application/json") {
			t.Fatalf("role %s Content-Type = %q, want application/json", role, rr.Header().Get("Content-Type"))
		}
		var result admincontrol.Result
		if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
			t.Fatalf("decode command result: %v", err)
		}
		if result.Action != "validation_health.refresh" || result.Status == "" || result.Summary == "" || len(result.NextActions) == 0 {
			t.Fatalf("invalid command result: %+v", result)
		}
		assertAdminCommandResultJSONAllowlist(t, rr.Body.Bytes())
		assertAdminCommandFlagsFalse(t, result.ClaimFlags)
		assertValidationHealthHTTPNoLeakage(t, rr.Body.String())
	}

	srv := newValidationHealthTestHandler(t, auth.RoleReadOnly, store, fakeScheduleBuilder{})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/validation-health/refresh.json", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET refresh status = %d, want 405", rr.Code)
	}
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health/refresh.json?agency_id=other-agency", strings.NewReader("action=refresh"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict refresh status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodPost, "/public/operations/validation-health/refresh.json", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public refresh status = %d, want 404", rr.Code)
	}
	unauth := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, authRejectAll{})
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health/refresh.json", strings.NewReader("action=refresh"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	unauth.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK {
		t.Fatalf("unauthenticated refresh returned 200")
	}
}

func TestValidationHealthRefreshCommandJSONStrictnessCSRFAndBodyCap(t *testing.T) {
	store := &fakePublicationStore{discovery: validationHealthTestDiscovery(time.Now().UTC())}
	srv := newValidationHealthTestHandler(t, auth.RoleReadOnly, store, fakeScheduleBuilder{})
	for _, field := range []string{"feed_type", "validator_id", "validator_path", "validator_command", "output_path", "artifact_path", "report_path", "schedule_zip_path", "realtime_pb_path", "path", "url", "URL", "argv", "args", "timeout", "timeout_seconds", "raw_report", "stdout", "stderr"} {
		req := httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health/refresh.json", strings.NewReader("action=refresh&"+field+"=browser"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("field %s status = %d, want 400", field, rr.Code)
		}
		assertValidationHealthHTTPNoLeakage(t, rr.Body.String())
	}
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health/refresh.json", strings.NewReader(strings.Repeat("x", validationHealthPostMaxBytes+1)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("large refresh form status = %d, want 413", rr.Code)
	}

	cookieAuth := auth.TestAuthenticator{Principal: auth.Principal{Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodCookie}}
	srv = newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}, csrfSecret: "test-csrf"}, cookieAuth)
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health/refresh.json", strings.NewReader("action=refresh"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("missing CSRF refresh status = %d, want 403", rr.Code)
	}
	principal := auth.Principal{Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodCookie}
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health/refresh.json", strings.NewReader("action=refresh&csrf_token="+url.QueryEscape(csrfToken("test-csrf", principal))))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("valid CSRF refresh status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
}

func TestValidationCenterRoutesPrivateScopedGETOnlyNoStore(t *testing.T) {
	store := feedHealthTestStore(t)
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
			Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
		}})
		for _, path := range []string{"/admin/operations/validation-center", "/admin/operations/validation-center.json"} {
			req := httptest.NewRequest(http.MethodGet, path+"?agency_id=demo-agency", nil)
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

	srv := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator}, Method: auth.MethodBearer,
	}})
	for _, path := range []string{"/admin/operations/validation-center", "/admin/operations/validation-center.json"} {
		for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
			req := httptest.NewRequest(method, path, nil)
			rr := httptest.NewRecorder()
			srv.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status = %d, want 405", method, path, rr.Code)
			}
		}
		req := httptest.NewRequest(http.MethodGet, path+"?agency_id=other-agency", nil)
		rr := httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("agency conflict %s status = %d, want 403", path, rr.Code)
		}
	}

	unauth := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, authRejectAll{})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/validation-center.json", nil)
	rr := httptest.NewRecorder()
	unauth.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK {
		t.Fatalf("unauthenticated validation center JSON returned 200")
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/validation-center.json", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public validation center route status = %d, want 404", rr.Code)
	}
}

func TestValidationCenterJSONShapeFeedRowsFlagsAndNoLeakage(t *testing.T) {
	now := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC)
	store := feedHealthTestStore(t)
	private := validationHealthRecord(99, "schedule", "feed-v1", "failed", now)
	private.Result.Report = map[string]any{
		"raw_report": map[string]any{"notices": []any{map[string]any{"code": "private_path", "severity": "ERROR", "message": "/Users/private TOKEN=SECRET", "path": "/Users/private/gtfs.zip"}}},
		"stdout":     "TOKEN=SECRET",
		"stderr":     "PASSWORD=SECRET",
		"argv":       []any{"/Users/private/validator"},
	}
	store.validationRecords = append(store.validationRecords, private)
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/validation-center.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var center operationsValidationCenterView
	if err := json.Unmarshal(rr.Body.Bytes(), &center); err != nil {
		t.Fatalf("decode validation center JSON: %v", err)
	}
	assertValidationCenterShape(t, center)
	assertValidationCenterFlagsFalse(t, center.ClaimFlags)
	assertValidationCenterJSONAllowlist(t, rr.Body.Bytes())
	assertValidationCenterSafeStrings(t, rr.Body.String())

	wantPaths := []string{
		"/public/feeds.json",
		"/public/gtfs/schedule.zip",
		"/public/gtfsrt/vehicle_positions.pb",
		"/public/gtfsrt/trip_updates.pb",
		"/public/gtfsrt/alerts.pb",
	}
	var gotPaths []string
	for _, row := range center.FeedRows {
		gotPaths = append(gotPaths, row.PublicPath)
	}
	if strings.Join(gotPaths, ",") != strings.Join(wantPaths, ",") {
		t.Fatalf("feed paths = %v, want %v", gotPaths, wantPaths)
	}
	gotConsumers := make([]string, 0, len(center.ConsumerTracker))
	for _, row := range center.ConsumerTracker {
		gotConsumers = append(gotConsumers, row.Target+":"+row.Status)
	}
	wantConsumers := []string{"Google Maps:prepared", "Apple Maps:prepared", "Transit App:prepared", "Bing Maps:prepared", "Moovit:prepared", "Mobility Database:prepared", "transit.land:prepared"}
	if strings.Join(gotConsumers, ",") != strings.Join(wantConsumers, ",") {
		t.Fatalf("consumer rows = %v, want %v", gotConsumers, wantConsumers)
	}
}

func TestValidationCenterHTMLPlainLanguageReadOnlyAndNoLeakage(t *testing.T) {
	store := feedHealthTestStore(t)
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/validation-center", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Feed Health And Validation Center",
		"Feed Health vs Validation",
		"Five Feed URL Panel",
		"Validation History",
		"Validator Health",
		"GTFS Quality Summary",
		"Prepared Consumer Tracker",
		"Safety details",
		"read-only",
		"does not run validators",
		"prepared only",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("validation center HTML missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"<form", "method=\"post\"", "/admin/validation/run", "raw_report", "stdout", "stderr", "argv", "/Users/private", "TOKEN=SECRET", "PASSWORD=SECRET", "postgres://", "production_ready", "compliance_achieved", "consumer_ready"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("validation center HTML leaked or overclaimed %q: %s", forbidden, body)
		}
	}
}

func TestValidationCenterIssueDrilldownsFixOwnersAndNoRawSamples(t *testing.T) {
	now := time.Date(2026, 5, 14, 13, 0, 0, 0, time.UTC)
	store := feedHealthTestStore(t)
	store.validationRecords = append(store.validationRecords,
		compliance.ValidationReportRecord{
			ID:        200,
			CreatedAt: now,
			Result: compliance.ValidationResult{
				AgencyID:      "demo-agency",
				FeedType:      "schedule",
				FeedVersionID: "feed-v1",
				ValidatorName: compliance.CanonicalStaticValidatorName,
				Status:        "failed",
				ErrorCount:    1,
				WarningCount:  1,
				Report: map[string]any{
					"raw_report": map[string]any{"notices": []any{
						map[string]any{"code": "expired_calendar", "severity": "ERROR", "message": "/Users/private TOKEN=SECRET", "path": "/Users/private/calendar.txt"},
						map[string]any{"code": "frequency_headway_invalid", "severity": "WARNING", "message": "frequency needs review", "filename": "frequencies.txt"},
					}},
					"stdout": "TOKEN=SECRET",
					"stderr": "PASSWORD=SECRET",
					"argv":   []any{"/Users/private/validator"},
				},
			},
		},
		compliance.ValidationReportRecord{
			ID:        201,
			CreatedAt: now,
			Result: compliance.ValidationResult{
				AgencyID:      "demo-agency",
				FeedType:      "schedule",
				FeedVersionID: "feed-v1",
				ValidatorName: compliance.InternalGTFSImportValidatorName,
				Status:        "failed",
				ErrorCount:    1,
				Report: map[string]any{"errors": []any{
					map[string]any{"code": "missing_trip_reference", "message": "missing trip reference", "path": "/Users/private/trips.txt"},
				}},
			},
		},
	)
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/validation-center.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("json status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var center operationsValidationCenterView
	if err := json.Unmarshal(rr.Body.Bytes(), &center); err != nil {
		t.Fatalf("decode validation center JSON: %v", err)
	}
	if len(center.IssueDrilldowns) < 3 {
		t.Fatalf("issue drilldowns = %+v, want at least 3", center.IssueDrilldowns)
	}
	body := rr.Body.String()
	for _, want := range []string{"Schedule planner or GTFS source owner", "calendar.txt / calendar_dates.txt", "frequencies.txt / trips.txt", "GTFS export owner or source-system admin", "reported referring file and referenced GTFS file"} {
		if !strings.Contains(body, want) {
			t.Fatalf("validation center issues missing %q: %s", want, body)
		}
	}
	assertValidationCenterSafeStrings(t, body)
	for _, issue := range center.IssueDrilldowns {
		if issue.SampleCount > 0 && (strings.Contains(strings.Join(issue.Codes, ","), "/Users") || strings.Contains(strings.Join(issue.Codes, ","), "TOKEN")) {
			t.Fatalf("issue codes leaked private text: %+v", issue)
		}
		if issue.LikelyOwner == "" || issue.RiskLevel == "" || issue.AffectedFiles == "" || issue.SafeFixPath == "" || issue.VerifyWith == "" || issue.EscalateIf == "" || issue.DoesNotProve == "" {
			t.Fatalf("issue missing guidance fields: %+v", issue)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/validation-center", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	html := rr.Body.String()
	for _, want := range []string{"Issue Drilldowns", "Likely owner", "Risk level", "Affected files", "Safe fix path", "Verify with", "Sample count"} {
		if !strings.Contains(html, want) {
			t.Fatalf("validation center issue HTML missing %q: %s", want, html)
		}
	}
	assertValidationCenterSafeStrings(t, html)
	for _, forbidden := range []string{"/Users/private", "TOKEN=SECRET", "PASSWORD=SECRET", "raw_report", "stdout", "stderr", "argv"} {
		if strings.Contains(html, forbidden) {
			t.Fatalf("validation center issue HTML leaked %q: %s", forbidden, html)
		}
	}
}

func TestValidationCenterReadinessTimelineAndBlockerQueue(t *testing.T) {
	now := time.Date(2026, 5, 14, 14, 0, 0, 0, time.UTC)
	discovery := validationHealthTestDiscovery(now)
	discovery.Readiness = compliance.Readiness{AllRequiredFeedsListed: false, LicenseComplete: false, ContactComplete: false, HTTPSURLs: true}
	discovery.Feeds = discovery.Feeds[:1]
	discovery.Feeds[0].CanonicalPublicURL = ""
	store := &fakePublicationStore{discovery: discovery, validationRecords: []compliance.ValidationReportRecord{
		validationHealthRecord(1, "schedule", "feed-v1", "failed", now),
	}}
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/validation-center.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("json status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var center operationsValidationCenterView
	if err := json.Unmarshal(rr.Body.Bytes(), &center); err != nil {
		t.Fatalf("decode validation center JSON: %v", err)
	}
	if len(center.ReadinessTimeline) != 11 {
		t.Fatalf("readiness timeline rows = %d, want 11: %+v", len(center.ReadinessTimeline), center.ReadinessTimeline)
	}
	if len(center.Blockers) == 0 {
		t.Fatalf("expected blocker rows for incomplete discovery")
	}
	body := rr.Body.String()
	for _, want := range []string{"Feed discovery and metadata", "Validator health", "Static GTFS Schedule", "blocked", "missing"} {
		if !strings.Contains(body, want) {
			t.Fatalf("validation center timeline/blockers missing %q: %s", want, body)
		}
	}
	for _, blocker := range center.Blockers {
		if blocker.ID == "" || blocker.Severity == "" || blocker.Area == "" || blocker.Signal == "" || blocker.NextAction == "" || blocker.DoesNotProve == "" || blocker.ReviewURL == "" {
			t.Fatalf("invalid blocker row: %+v", blocker)
		}
		if !strings.HasPrefix(blocker.ReviewURL, "/admin/") {
			t.Fatalf("unsafe blocker review URL: %+v", blocker)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/validation-center", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	html := rr.Body.String()
	for _, want := range []string{"Readiness Timeline", "Current Blockers", "Review"} {
		if !strings.Contains(html, want) {
			t.Fatalf("validation center timeline HTML missing %q: %s", want, html)
		}
	}
	assertValidationCenterSafeStrings(t, html)
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

func TestOperationsMaintenanceRoutesJSONShapeFlagsAndPrivateBoundaries(t *testing.T) {
	t.Setenv("VALIDATOR_TOOLING_MODE", "stub")
	t.Setenv("BACKUP_DIR", "/private/withheld")
	roots := writeMaintenanceSummaryFixtures(t)
	withMaintenanceSummaryRoots(t, roots)
	now := time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC)
	store := feedHealthTestStore(t)
	store.discovery.GeneratedAt = now
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	for _, path := range []string{"/admin/operations/maintenance", "/admin/operations/maintenance.json"} {
		req := httptest.NewRequest(http.MethodGet, path+"?agency_id=demo-agency", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
		}
		if got := rr.Header().Get("Cache-Control"); got != "no-store" {
			t.Fatalf("%s Cache-Control = %q, want no-store", path, got)
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/maintenance.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	var view operationsMaintenanceView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode maintenance JSON: %v", err)
	}
	assertMaintenanceShape(t, view)
	assertMaintenanceFlagsFalse(t, view.ClaimFlags)
	wantRows := []string{"deployed_version", "active_feed_version", "last_gtfs_import", "last_five_feed_check", "validator_state", "backup_configuration", "restore_drill_configuration", "telemetry_freshness", "service_health"}
	var gotRows []string
	for _, row := range view.SummaryRows {
		gotRows = append(gotRows, row.ID)
	}
	if strings.Join(gotRows, ",") != strings.Join(wantRows, ",") {
		t.Fatalf("maintenance rows = %v, want %v", gotRows, wantRows)
	}
	if view.SupportSummary.OutputPath != ".cache/support-bundles/<timestamp>" {
		t.Fatalf("support output path = %q, want timestamp placeholder", view.SupportSummary.OutputPath)
	}
	if view.Diagnostics.Status != operationsStatusBlocked || len(view.Diagnostics.Rows) != 4 {
		t.Fatalf("unexpected diagnostics summary: %+v", view.Diagnostics)
	}
	if view.SmallHostReadiness.Status != operationsStatusNeedsReview || len(view.SmallHostReadiness.Rows) != 5 {
		t.Fatalf("unexpected small-host readiness panel: %+v", view.SmallHostReadiness)
	}
	if view.BackupRestore.Status != operationsStatusMissing || len(view.BackupRestore.Rows) != 4 || view.UpgradeRollback.Status != operationsStatusNeedsReview || len(view.UpgradeRollback.Rows) != 4 {
		t.Fatalf("unexpected maintenance panels: backup=%+v upgrade=%+v", view.BackupRestore, view.UpgradeRollback)
	}
	if view.SupportReview.Status != operationsStatusNeedsReview || len(view.SupportReview.Rows) != 4 || view.CadencePlan.Status != operationsStatusNeedsReview || len(view.CadencePlan.Rows) != 4 {
		t.Fatalf("unexpected support/cadence panels: support=%+v cadence=%+v", view.SupportReview, view.CadencePlan)
	}
	if view.MonitoringExport.Status != operationsStatusDiagnosticOnly || len(view.MonitoringExport.Rows) != 5 {
		t.Fatalf("unexpected monitoring export panel: %+v", view.MonitoringExport)
	}
	if view.Infrastructure.Status != operationsStatusBlocked || len(view.Infrastructure.Rows) != 10 {
		t.Fatalf("unexpected infrastructure panel: %+v", view.Infrastructure)
	}
	body := rr.Body.String()
	for _, want := range []string{"values withheld", "not configured", "make support-bundle", ".cache/support-bundles/", "deployment_doctor", "operations_reliability", "operations_notify", "support_bundle_manifest", "backup=blocker", "not_sent=true", "small_host_readiness", "small_host_preflight_sequence", "off_host_validation_choice", "resource_budget_review", "backup_restore_recovery_path", "upgrade_recovery_stop_points", "backup_configuration_presence", "restore_drill_configuration_presence", "upgrade_precheck", "rollback_precheck", "browser_destructive_actions", "release_artifact_boundary", "support_bundle_output_scope", "redaction_review", "evidence_boundary", "private_output_warning", "daily_operating_check", "weekly_maintenance_check", "monthly_recovery_check", "as_needed_support_check", "operations_notify_health_digest", "redacted_channel_guidance", "monitoring_export_summary_json", "redacted_operations_export_formats", "feed_health", "connector_health", "validator_posture", "telemetry_freshness", "maintenance_tasks", "no_send_default", "browser_send_enabled=false", "webhook_send_enabled=false", "email_send_enabled=false", "database_connectivity", "migration_status", "postgis_extension", "validator_tooling", "backup_storage_access", "small_host_resources", "service_dependencies", "proxy_exposure", "postgres_capacity", "upgrade_rollback_checklist"} {
		if !strings.Contains(body, want) {
			t.Fatalf("maintenance JSON missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"/private/withheld", "postgres://", "Authorization", "Bearer ", "consumer accepted", "production ready", "compliance achieved", "sla_covered"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("maintenance JSON leaked or overclaimed %q: %s", forbidden, body)
		}
	}

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/admin/operations/maintenance", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s status = %d, want 405", method, rr.Code)
		}
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/maintenance?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/maintenance.json", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public maintenance route status = %d, want 404", rr.Code)
	}
}

func TestOperationsMaintenanceSummaryReadersRejectUnsafeRootsAndClaimFlags(t *testing.T) {
	previous := maintenanceDeploymentDoctorRoot
	maintenanceDeploymentDoctorRoot = filepath.Join(t.TempDir(), "deployment-doctor")
	t.Cleanup(func() { maintenanceDeploymentDoctorRoot = previous })
	diagnostics := buildOperationsMaintenanceDiagnostics()
	if diagnostics.Rows[0].Status != operationsStatusBlocked || !strings.Contains(diagnostics.Rows[0].CurrentSignal, "failed private cache") {
		t.Fatalf("unsafe root was not blocked: %+v", diagnostics.Rows[0])
	}

	roots := writeMaintenanceSummaryFixtures(t)
	writeMaintenanceJSON(t, filepath.Join(roots.deploymentDoctor, "20260512T010313Z", "summary.json"), map[string]any{
		"generated_at_utc":             "20260512T010313Z",
		"overall_status":               "passed",
		"counts":                       map[string]any{"blocker": 0, "warning": 0},
		"categories":                   map[string]any{"backup_readiness": "passed", "restore_readiness": "passed"},
		"external_evidence_created":    false,
		"final_root_evidence_created":  false,
		"consumer_statuses_changed":    false,
		"compliance_claimed":           false,
		"production_readiness_claimed": true,
	})
	withMaintenanceSummaryRoots(t, roots)
	diagnostics = buildOperationsMaintenanceDiagnostics()
	if diagnostics.Rows[0].Status != operationsStatusBlocked || !strings.Contains(diagnostics.Rows[0].CurrentSignal, "claim flag") {
		t.Fatalf("claim flag summary was not blocked: %+v", diagnostics.Rows[0])
	}
}

type maintenanceSummaryTestRoots struct {
	deploymentDoctor      string
	operationsReliability string
	operationsNotify      string
	supportBundles        string
}

func withMaintenanceSummaryRoots(t *testing.T, roots maintenanceSummaryTestRoots) {
	t.Helper()
	previousDeployment := maintenanceDeploymentDoctorRoot
	previousReliability := maintenanceOperationsReliabilityRoot
	previousNotify := maintenanceOperationsNotifyRoot
	previousSupport := maintenanceSupportBundleRoot
	maintenanceDeploymentDoctorRoot = roots.deploymentDoctor
	maintenanceOperationsReliabilityRoot = roots.operationsReliability
	maintenanceOperationsNotifyRoot = roots.operationsNotify
	maintenanceSupportBundleRoot = roots.supportBundles
	t.Cleanup(func() {
		maintenanceDeploymentDoctorRoot = previousDeployment
		maintenanceOperationsReliabilityRoot = previousReliability
		maintenanceOperationsNotifyRoot = previousNotify
		maintenanceSupportBundleRoot = previousSupport
	})
}

func writeMaintenanceSummaryFixtures(t *testing.T) maintenanceSummaryTestRoots {
	t.Helper()
	base := filepath.Join(t.TempDir(), ".cache")
	roots := maintenanceSummaryTestRoots{
		deploymentDoctor:      filepath.Join(base, "deployment-doctor"),
		operationsReliability: filepath.Join(base, "operations-reliability"),
		operationsNotify:      filepath.Join(base, "operations-notify"),
		supportBundles:        filepath.Join(base, "support-bundles"),
	}
	writeMaintenanceJSON(t, filepath.Join(roots.deploymentDoctor, "20260512T010313Z", "summary.json"), map[string]any{
		"generated_at_utc":               "20260512T010313Z",
		"overall_status":                 "blocker",
		"counts":                         map[string]any{"blocker": 2, "warning": 1, "unavailable": 0},
		"categories":                     map[string]any{"backup_readiness": "blocker", "database": "skipped", "migrations": "skipped", "postgis": "skipped", "restore_readiness": "blocker", "validators": "passed", "small_host_resources": "warning", "service_dependencies": "passed", "proxy_exposure": "passed", "postgres_capacity": "warning", "upgrade_rollback": "passed"},
		"external_evidence_created":      false,
		"final_root_evidence_created":    false,
		"consumer_statuses_changed":      false,
		"compliance_claimed":             false,
		"production_readiness_claimed":   false,
		"hosted_saas_claimed":            false,
		"sla_claimed":                    false,
		"uptime_guarantee_claimed":       false,
		"vendor_compatibility_claimed":   false,
		"hardware_certification_claimed": false,
		"production_grade_eta_claimed":   false,
	})
	writeMaintenanceJSON(t, filepath.Join(roots.operationsReliability, "20260512T010400Z", "summary.json"), map[string]any{
		"generated_at":          "20260512T010400Z",
		"overall_status":        "missing",
		"backup_restore":        map[string]any{"status": "missing"},
		"alerting":              map[string]any{"status": "missing"},
		"availability_sampling": map[string]any{"status": "unknown"},
		"claim_flags":           map[string]any{"external_evidence_created": false, "final_root_evidence_created": false, "consumer_statuses_changed": false, "compliance_claimed": false, "production_readiness_claimed": false, "hosted_saas_claimed": false, "consumer_acceptance_claimed": false, "vendor_compatibility_claimed": false, "sla_claimed": false, "uptime_guarantee_claimed": false, "production_grade_eta_claimed": false},
	})
	writeMaintenanceJSON(t, filepath.Join(roots.operationsNotify, "20260512T010500Z", "summary.json"), map[string]any{
		"generated_at":                 "20260512T010500Z",
		"notification":                 map[string]any{"severity": "needs_review", "not_sent": true},
		"counts":                       map[string]any{"next_actions": 2, "blocked_actions": 1},
		"external_evidence_created":    false,
		"consumer_statuses_changed":    false,
		"compliance_claimed":           false,
		"production_readiness_claimed": false,
		"hosted_saas_claimed":          false,
		"agency_adoption_claimed":      false,
		"consumer_acceptance_claimed":  false,
		"vendor_compatibility_claimed": false,
		"production_grade_eta_claimed": false,
		"notification_sent":            false,
	})
	writeMaintenanceJSON(t, filepath.Join(roots.supportBundles, "20260512T010600Z", "manifest.json"), map[string]any{
		"created_at_utc":            "20260512T010600Z",
		"external_evidence_created": false,
		"consumer_statuses_changed": false,
		"included":                  []string{"system versions", "git summary"},
		"excluded":                  []string{"credential values", "private payloads"},
	})
	return roots
}

func writeMaintenanceJSON(t *testing.T, path string, value any) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create maintenance fixture dir: %v", err)
	}
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal maintenance fixture: %v", err)
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatalf("write maintenance fixture: %v", err)
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

func TestValidationHealthCommandDefinitionsBoundedForBrowserControl(t *testing.T) {
	refresh := admincontrol.ValidationHealthRefreshDefinition()
	if refresh.Action != "validation_health.refresh" || refresh.LadderLevel != admincontrol.LevelReadOnlyRefresh || refresh.RequiredRole != "read_only" {
		t.Fatalf("unexpected refresh command definition: %+v", refresh)
	}
	for _, want := range []string{"No public feed output changes", "Writes nothing", "This private view does not show compliance"} {
		if !strings.Contains(refresh.PublicFeedImpact+refresh.PrivateImpact+refresh.DoesNotProve, want) {
			t.Fatalf("refresh command definition missing %q: %+v", want, refresh)
		}
	}
	runAll := admincontrol.ValidationHealthRunAllDefinition()
	if runAll.Action != "validation_health.run_all" || runAll.LadderLevel != admincontrol.LevelReversiblePrivate || runAll.RequiredRole != "admin" {
		t.Fatalf("unexpected run-all command definition: %+v", runAll)
	}
	for _, want := range []string{"No public feed output changes", "validation_report", "This private view does not show compliance"} {
		if !strings.Contains(runAll.PublicFeedImpact+runAll.PrivateImpact+runAll.DoesNotProve, want) {
			t.Fatalf("run-all command definition missing %q: %+v", want, runAll)
		}
	}
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
	for _, want := range []string{"Detailed safety booleans remain available", "private diagnostics"} {
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
	for _, field := range []string{"feed_type", "validator_id", "validator_path", "validator_command", "output_path", "artifact_path", "report_path", "schedule_zip_path", "realtime_pb_path", "path", "url", "URL", "argv", "args", "timeout", "timeout_seconds", "raw_report", "stdout", "stderr"} {
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

func TestValidationHealthRunAllRequiresAdminRole(t *testing.T) {
	store := &fakePublicationStore{discovery: validationHealthTestDiscovery(time.Now().UTC())}
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor} {
		srv := newValidationHealthTestHandler(t, role, store, fakeScheduleBuilder{})
		req := httptest.NewRequest(http.MethodPost, "/admin/operations/validation-health", strings.NewReader(validationHealthRunAllForm()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("role %s status = %d, want 403", role, rr.Code)
		}
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
			AdapterDetails: map[string]any{"external_http_shadow": map[string]any{
				"status":                           prediction.StatusError,
				"reason":                           prediction.ReasonAdapterError,
				"latency_ms":                       42.0,
				"deterministic_trip_updates_count": 1.0,
				"external_trip_updates_count":      0.0,
				"count_delta":                      -1.0,
				"raw_response":                     "token=secret host=predictor.example",
			}},
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

func TestRealtimeOperationsCenterPrivateReadOnlyFleetOverview(t *testing.T) {
	t.Setenv("SUPPRESS_STALE_VEHICLE_AFTER_SECONDS", "120")
	now := time.Now().UTC().Truncate(time.Second)
	store := feedHealthTestStore(t)
	store.tripDiagnostics = compliance.TripUpdatesDiagnosticsSummary{
		Recorded:            true,
		SnapshotAt:          now,
		AdapterName:         "deterministic",
		DiagnosticsStatus:   "recorded",
		DiagnosticsReason:   "partial_predictions",
		ActiveFeedVersionID: "feed-v1",
		Metrics: prediction.Metrics{
			EligiblePredictionCandidates: 2,
			TripUpdatesEmitted:           1,
			UnknownAssignments:           1,
			AmbiguousAssignments:         1,
			StaleTelemetryRows:           1,
			WithheldByReason: map[string]int{
				prediction.ReasonBelowConfidenceThreshold: 1,
				prediction.ReasonStaleTelemetry:           1,
			},
			CancellationAlertLinksMissing: 1,
		},
	}
	srv := newOperationsTestHandler(&handler{
		store: store,
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{
			{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now.Add(-time.Hour), CreatedAt: now.Add(-time.Hour)},
			{AgencyID: "demo-agency", DeviceID: "device-2", VehicleID: "bus-2", Status: "active", ValidFrom: now.Add(-time.Hour), CreatedAt: now.Add(-time.Hour)},
			{AgencyID: "demo-agency", DeviceID: "device-3", VehicleID: "bus-3", Status: "active", ValidFrom: now.Add(-time.Hour), CreatedAt: now.Add(-time.Hour)},
		}},
		telemetry: fakeTelemetryRepository{latest: []telemetry.StoredEvent{
			{
				ID: 11,
				Event: telemetry.Event{
					AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Timestamp: now.Add(-30 * time.Second), Lat: 1, Lon: 2,
				},
				ReceivedAt: now.Add(-29 * time.Second), IngestStatus: telemetry.IngestStatusAccepted, PayloadJSON: []byte(`{"secret":"hidden"}`),
			},
			{
				ID: 22,
				Event: telemetry.Event{
					AgencyID: "demo-agency", DeviceID: "device-2", VehicleID: "bus-2", Timestamp: now.Add(-5 * time.Minute), Lat: 3, Lon: 4,
				},
				ReceivedAt: now.Add(-5*time.Minute + time.Second), IngestStatus: telemetry.IngestStatusAccepted, PayloadJSON: []byte(`{"raw_payload":"hidden"}`),
			},
		}},
		state: fakeStateRepository{assignments: map[string]state.Assignment{
			"bus-1": {
				VehicleID: "bus-1", State: state.StateInService, RouteID: "route-1", TripID: "trip-1", Confidence: 0.91,
				ReasonCodes: []string{state.ReasonTripHintMatch}, DegradedState: state.DegradedNone, AssignmentSource: state.AssignmentSourceAutomatic,
				ScoreDetails: map[string]any{"private_debug": true}, ActiveFrom: now.Add(-25 * time.Second),
			},
			"bus-2": {
				VehicleID: "bus-2", State: state.StateUnknown, Confidence: 0.20,
				ReasonCodes: []string{state.ReasonLowConfidence}, DegradedState: state.DegradedLowConfidence, AssignmentSource: state.AssignmentSourceAutomatic,
				ScoreDetails: map[string]any{"private_debug": true}, ActiveFrom: now.Add(-4 * time.Minute),
			},
		}},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})

	for _, path := range []string{"/admin/operations/realtime", "/admin/operations/realtime.json"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
		}
		if got := rr.Header().Get("Cache-Control"); got != "no-store" {
			t.Fatalf("%s Cache-Control = %q, want no-store", path, got)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/realtime", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	body := rr.Body.String()
	for _, want := range []string{"Realtime", "Fleet Freshness", "Realtime Feed Usefulness Review", "Feed Usefulness Details", "Healthy when", "Needs attention", "Not proven", "Vehicle Positions publishing review", "Vehicle count", "Estimated Vehicle Positions rows", "Suppressed vehicles", "Trip descriptor coverage", "Why not published", "Trip Updates publishing review", "Prediction source", "Fallback reason", "Low-confidence handling", "Withheld reason: below_confidence_threshold", "Alerts lifecycle review", "Active alerts", "Stale alerts", "Missing cancellation links", "Service disruption review", "Synthetic / Local Replay Guide", "Preview a scenario in the browser", "Freshness And Lifecycle Review", "Consumer-Safe Omission Rules", "Vehicle Positions usefulness", "Trip Updates usefulness", "Alerts usefulness", "Consumer-safe behavior", "emitted from", "valid empty/fallback output", "Needs Operator Review", "Realtime Quality Guidance", "Out-of-order or low-quality GPS", "Trip Updates withheld or fallback", "bus-1", "bus-2", "bus-3", "fresh", "stale", "not seen", "keep trip descriptors unknown", "Vehicle Positions", "Trip Updates", "Alerts"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"<form", "payload_json", "raw_payload", "secret", "score_details", "private_debug", "token_hash", "Bearer", "production-grade ETA quality</strong>"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("realtime page leaks or overstates %q: %s", forbidden, body)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/realtime.json", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	var view operationsRealtimeView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode realtime json: %v", err)
	}
	if view.AgencyID != "demo-agency" || view.Boundary == "" || len(view.Fleet) != 3 || len(view.Feeds) != 3 || len(view.Issues) == 0 || len(view.Guidance) < 6 {
		t.Fatalf("unexpected realtime JSON shape: %+v", view)
	}
	if view.Usefulness.Status == "" || len(view.Usefulness.Rows) != 3 || len(view.Usefulness.Freshness) != 5 || len(view.Usefulness.OmissionRules) != 4 || view.Usefulness.Boundary == "" {
		t.Fatalf("unexpected realtime usefulness shape: %+v", view.Usefulness)
	}
	if len(view.Publishing) != 3 || view.ReplayGuidance.Status == "" || view.ReplayGuidance.BrowserStart != "/admin/operations/telemetry-simulator" || len(view.ReplayGuidance.Steps) != 3 || view.ReplayGuidance.Boundary == "" {
		t.Fatalf("unexpected realtime publishing/replay shape: publishing=%+v replay=%+v", view.Publishing, view.ReplayGuidance)
	}
	if view.Issues[0].Severity == "" || view.Issues[0].NextAction == "" || view.Guidance[0].DoesNotProve == "" {
		t.Fatalf("realtime review guidance is not actionable: issues=%+v guidance=%+v", view.Issues, view.Guidance)
	}
	seenPublishing := map[string]operationsRealtimeFeedReview{}
	for _, row := range view.Publishing {
		seenPublishing[row.ID] = row
		if row.Label == "" || row.Status == "" || row.WhatLooksHealthy == "" || row.NeedsAttention == "" || row.NotProven == "" || row.NextAction == "" || len(row.Signals) == 0 {
			t.Fatalf("publishing review row is not actionable: %+v", row)
		}
		for _, signal := range row.Signals {
			if signal.Label == "" || signal.Value == "" || signal.Meaning == "" {
				t.Fatalf("publishing signal missing operator wording: row=%+v signal=%+v", row, signal)
			}
		}
	}
	if seenPublishing["vehicle_positions"].Signals[4].Label != "Suppressed vehicles" || seenPublishing["trip_updates"].Signals[0].Label != "Prediction source" || seenPublishing["alerts"].Signals[2].Label != "Missing cancellation links" {
		t.Fatalf("publishing review rows missing expected feed-specific signals: %+v", seenPublishing)
	}
	for _, row := range view.Usefulness.Rows {
		if row.ID == "" || row.Label == "" || row.ScoreLabel == "" || row.CurrentSignal == "" || row.HelpfulSignal == "" || row.NeedsReviewSignal == "" || row.ConsumerSafeBehavior == "" || row.NextAction == "" || row.DoesNotProve == "" {
			t.Fatalf("realtime usefulness row is not actionable: %+v", row)
		}
		if row.Score < 0 || row.Score > 3 {
			t.Fatalf("realtime usefulness score out of range: %+v", row)
		}
	}
	for _, rule := range view.Usefulness.OmissionRules {
		if rule.Condition == "" || rule.SafeBehavior == "" || rule.ReviewStep == "" || rule.DoesNotProve == "" {
			t.Fatalf("omission rule missing operator guidance: %+v", rule)
		}
	}
	if view.Summary.LatestTelemetryRows != 2 || view.Summary.StaleTelemetryRows != 1 || view.Summary.DeviceBindings != 3 || view.Summary.DevicesNotSeen != 1 || view.Summary.MatchedAssignments != 1 || view.Summary.UnknownAssignments != 2 || view.Summary.LowConfidenceRows != 1 {
		t.Fatalf("unexpected realtime summary: %+v", view.Summary)
	}
	assertRealtimeFlagsFalse(t, view.ClaimFlags)
	realtimePayload := rr.Body.String()
	for _, forbidden := range []string{"payload_json", "raw_payload", "secret", "score_details", "private_debug", "token_hash", "Bearer", `"lat":`, `"lon":`} {
		if strings.Contains(realtimePayload, forbidden) {
			t.Fatalf("realtime JSON leaks %q: %s", forbidden, realtimePayload)
		}
	}

	unauth := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, authRejectAll{})
	for _, path := range []string{"/admin/operations/realtime", "/admin/operations/realtime.json"} {
		req = httptest.NewRequest(http.MethodGet, path, nil)
		rr = httptest.NewRecorder()
		unauth.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("unauth %s status = %d, want 401", path, rr.Code)
		}
		req = httptest.NewRequest(http.MethodPost, path, nil)
		rr = httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("POST %s status = %d, want 405", path, rr.Code)
		}
		req = httptest.NewRequest(http.MethodGet, path+"?agency_id=other-agency", nil)
		rr = httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("agency conflict %s status = %d, want 403", path, rr.Code)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/public/operations/realtime", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public realtime route status = %d, want 404", rr.Code)
	}
}

func TestPredictionLabRoutesPrivateScopedGETOnlyNoStore(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin} {
		t.Run(string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "user@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			for _, path := range []string{"/admin/operations/prediction-lab", "/admin/operations/prediction-lab.json"} {
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
	for _, path := range []string{"/admin/operations/prediction-lab", "/admin/operations/prediction-lab.json"} {
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
		for _, path := range []string{"/admin/operations/prediction-lab", "/admin/operations/prediction-lab.json"} {
			req := httptest.NewRequest(method, path, nil)
			rr := httptest.NewRecorder()
			authenticated.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status = %d, want 405", method, path, rr.Code)
			}
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/prediction-lab?agency_id=other-agency", nil)
	rr := httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict html status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/prediction-lab.json?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	authenticated.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("agency conflict json status = %d, want 403", rr.Code)
	}
	for _, path := range []string{"/public/operations/prediction-lab", "/public/operations/prediction-lab.json"} {
		req = httptest.NewRequest(http.MethodGet, path, nil)
		rr = httptest.NewRecorder()
		authenticated.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("public prediction lab route %s status = %d, want 404", path, rr.Code)
		}
	}
}

func TestPredictionLabJSONShapeFlagsAndDeterministicDiagnostics(t *testing.T) {
	now := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC)
	coverage := 50.0
	backtestRoot := filepath.Join(t.TempDir(), ".cache", "realtime-quality-backtest")
	writePredictionLabBacktestFixture(t, backtestRoot)
	withPredictionLabBacktestRoot(t, backtestRoot)
	store := &fakePublicationStore{
		tripDiagnostics: compliance.TripUpdatesDiagnosticsSummary{
			Recorded:                      true,
			SnapshotAt:                    now,
			AdapterName:                   "deterministic",
			DiagnosticsStatus:             prediction.StatusOK,
			DiagnosticsReason:             prediction.ReasonPartialPredictions,
			ActiveFeedVersionID:           "feed-demo",
			DiagnosticsPersistenceOutcome: "stored",
			AdapterDetails: map[string]any{"external_http_shadow": map[string]any{
				"status":                           prediction.StatusError,
				"reason":                           prediction.ReasonAdapterError,
				"latency_ms":                       42.0,
				"deterministic_trip_updates_count": 1.0,
				"external_trip_updates_count":      0.0,
				"count_delta":                      -1.0,
				"raw_response":                     "token=secret host=predictor.example",
			}},
			Metrics: prediction.Metrics{
				TelemetryRowsConsidered:      2,
				AssignmentsConsidered:        2,
				EligiblePredictionCandidates: 2,
				TripUpdatesEmitted:           1,
				UnknownAssignments:           1,
				AmbiguousAssignments:         1,
				StaleTelemetryRows:           1,
				ManualOverrideAssignments:    1,
				WithheldByReason: map[string]int{
					prediction.ReasonDegradedAssignment:       2,
					prediction.ReasonBelowConfidenceThreshold: 1,
				},
				UnknownAssignmentRate:   prediction.RateMetric{Numerator: 1, Denominator: 2, Percent: &coverage, Status: "measured", DenominatorDefinition: "current unknown assignments / current assignments considered"},
				AmbiguousAssignmentRate: prediction.RateMetric{Numerator: 1, Denominator: 2, Percent: &coverage, Status: "measured", DenominatorDefinition: "current ambiguous assignments / current assignments considered"},
				StaleTelemetryRate:      prediction.RateMetric{Numerator: 1, Denominator: 2, Percent: &coverage, Status: "measured", DenominatorDefinition: "stale latest telemetry rows / telemetry rows considered"},
				TripUpdatesCoverageRate: prediction.RateMetric{Numerator: 1, Denominator: 2, Percent: &coverage, Status: "measured", DenominatorDefinition: "emitted non-canceled Trip Updates / eligible in-service ETA candidates"},
				FutureStopCoverageRate:  prediction.RateMetric{Numerator: 1, Denominator: 2, Percent: &coverage, Status: "measured", DenominatorDefinition: "non-canceled Trip Updates with at least one future stop update / eligible in-service ETA candidates"},
			},
		},
	}
	srv := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/prediction-lab.json?agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json prefix", got)
	}
	var view predictionLabView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode prediction lab: %v", err)
	}
	assertPredictionLabShape(t, view)
	assertPredictionLabFlagsFalse(t, view.ClaimFlags)
	assertPredictionLabSafeStrings(t, rr.Body.String())
	if view.AgencyID != "demo-agency" || view.Summary.AdapterName != "deterministic" || view.Summary.WithheldCount != 3 || view.Summary.TripUpdatesEmitted != 1 || len(view.WithheldReasons) != 2 || len(view.ReviewRows) == 0 {
		t.Fatalf("unexpected prediction lab summary: %+v reasons=%+v reviews=%+v", view.Summary, view.WithheldReasons, view.ReviewRows)
	}
	if view.Deterministic.Status != checklistStatusOK || len(view.Deterministic.Rows) != 4 {
		t.Fatalf("unexpected deterministic diagnostics: %+v", view.Deterministic)
	}
	if view.ShadowReview.Status != checklistStatusNeedsReview || len(view.ShadowReview.Rows) != 1 || view.ShadowReview.Rows[0].Status != prediction.StatusError || !strings.Contains(view.ShadowReview.Rows[0].CountComparison, "delta=-1") {
		t.Fatalf("unexpected shadow review: %+v", view.ShadowReview)
	}
	if view.Backtests.Status != "needs_review" || len(view.Backtests.Rows) != 1 || view.Backtests.Rows[0].MaturityGate != "diagnostic_watch" || view.Backtests.Rows[0].PredictionCoverage != "62.5% (5/8)" || view.Backtests.Rows[0].ConformanceSignal != "synthetic_covered (5/5 synthetic cases)" {
		t.Fatalf("unexpected backtest browser: %+v", view.Backtests)
	}
	seen := map[string]bool{}
	for _, reason := range view.WithheldReasons {
		seen[reason.Reason] = true
		if reason.WhatItMeans == "" || reason.NextAction == "" || reason.DoesNotProve == "" {
			t.Fatalf("withheld reason lacks guidance: %+v", reason)
		}
	}
	for _, want := range []string{prediction.ReasonDegradedAssignment, prediction.ReasonBelowConfidenceThreshold} {
		if !seen[want] {
			t.Fatalf("missing withheld reason %q in %+v", want, view.WithheldReasons)
		}
	}
}

func TestPredictionLabHTMLBoundariesNoFormsAndEscapes(t *testing.T) {
	now := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC)
	backtestRoot := filepath.Join(t.TempDir(), ".cache", "realtime-quality-backtest")
	writePredictionLabBacktestFixture(t, backtestRoot)
	withPredictionLabBacktestRoot(t, backtestRoot)
	store := &fakePublicationStore{
		discovery: compliance.FeedDiscovery{AgencyID: "demo-agency", AgencyName: `<script>alert("x")</script>`, GeneratedAt: now},
		tripDiagnostics: compliance.TripUpdatesDiagnosticsSummary{
			Recorded:            true,
			SnapshotAt:          now,
			AdapterName:         "deterministic",
			DiagnosticsStatus:   prediction.StatusOK,
			DiagnosticsReason:   prediction.ReasonPartialPredictions,
			ActiveFeedVersionID: "feed-demo",
			AdapterDetails: map[string]any{"external_http_shadow": map[string]any{
				"status":                           prediction.StatusOK,
				"reason":                           prediction.ReasonPredictionsAvailable,
				"latency_ms":                       17.0,
				"deterministic_trip_updates_count": 1.0,
				"external_trip_updates_count":      1.0,
				"count_delta":                      0.0,
				"raw_response":                     "token=secret host=predictor.example",
			}},
			Metrics: prediction.Metrics{
				EligiblePredictionCandidates: 1,
				WithheldByReason:             map[string]int{prediction.ReasonStaleTelemetry: 1},
			},
		},
	}
	handler := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/prediction-lab", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("html status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Prediction &amp; ETA Lab", "Trip Updates Decision", "Safe Fallback", "Deterministic Predictor Diagnostics", "Why ETAs Are Missing", "External Predictor Shadow Review", "External HTTP shadow", "deterministic=1; external=1; delta=&#43;0", "Backtest Summary", ".cache/realtime-quality-backtest/20260514T120000Z", "manual_override_review=1", "synthetic_covered (5/5 synthetic cases)", "Conservative Handling Guide", "Telemetry is stale", "Assignment is ambiguous", "Future ETA Proof Gates", "Real observed arrival/departure comparison", "Required before collecting", "Stale Telemetry", "Vehicle Positions stay independent", "Needs Operator Review", "Fixed Local Checks", "Detailed safety booleans remain available", "make realtime-quality", "make realtime-quality-backtest"} {
		if !strings.Contains(body, want) {
			t.Fatalf("prediction lab html missing %q: %s", want, body)
		}
	}
	if strings.Contains(body, `<script>alert("x")</script>`) {
		t.Fatalf("html did not escape script-like metadata: %s", body)
	}
	for _, forbidden := range []string{`<form`, `method="post"`, "/public/operations/prediction-lab", "test predictor", "test connection", "start sidecar", "validated ETA performance", "consumer-ready", "agency approved", "consumer accepted", "production ready", "launch complete", "compliance achieved", "vendor compatible", "certified hardware"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("prediction lab html contains forbidden %q: %s", forbidden, body)
		}
	}
	assertPredictionLabSafeStrings(t, body)
}

func withPredictionLabBacktestRoot(t *testing.T, root string) {
	t.Helper()
	previous := predictionLabBacktestRoot
	predictionLabBacktestRoot = root
	t.Cleanup(func() {
		predictionLabBacktestRoot = previous
	})
}

func writePredictionLabBacktestFixture(t *testing.T, root string) {
	t.Helper()
	dir := filepath.Join(root, "20260514T120000Z")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create prediction lab backtest fixture: %v", err)
	}
	generatedAt := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	coverage := 62.5
	futureCoverage := 50.0
	mae := 29.0
	p90 := 40.0
	overall := realtimequality.MetricGroup{
		GroupType:                   "overall",
		CoverageDenominator:         8,
		CoverageNumerator:           5,
		FutureStopCoverageNumerator: 4,
		MatchedPredictionCount:      5,
		MissingPredictionCount:      1,
		MissingObservationCount:     1,
		StalePredictionCount:        1,
		WithheldByReason:            map[string]int{"manual_override_review": 1},
		MAEAbsoluteErrorSeconds:     &mae,
		P90AbsoluteErrorSeconds:     &p90,
		PredictionCoverage:          realtimequality.Rate{Numerator: 5, Denominator: 8, Status: "measured", Percent: &coverage},
		FutureStopCoverage:          realtimequality.Rate{Numerator: 4, Denominator: 8, Status: "measured", Percent: &futureCoverage},
		MaturityGate:                "diagnostic_watch",
	}
	summary := realtimequality.SummaryDocument{
		SchemaVersion:     realtimequality.BacktestSchemaVersion,
		GeneratedAt:       generatedAt,
		InputRecordCounts: realtimequality.RecordCounts{ObservedRecords: 8, PredictionRecords: 7},
		Overall:           overall,
		Conformance:       coveredPredictionLabConformanceReview(),
		GroupCount:        1,
	}
	metrics := realtimequality.MetricsDocument{
		SchemaVersion: realtimequality.BacktestSchemaVersion,
		GeneratedAt:   generatedAt,
		Groups:        []realtimequality.MetricGroup{overall},
	}
	manifest := realtimequality.ManifestDocument{
		SchemaVersion: realtimequality.BacktestSchemaVersion,
		GeneratedAt:   generatedAt,
		OutputKind:    "private_local_realtime_quality_backtest",
		OutputFiles:   realtimequality.ExpectedBacktestOutputFiles(),
		SafetyChecks: map[string]bool{
			"docs_evidence_output_rejected": true,
			"evidence_like_output_rejected": true,
			"symlink_ancestors_rejected":    true,
			"raw_inputs_not_copied":         true,
			"private_paths_omitted":         true,
			"raw_rows_omitted":              true,
		},
		Boundaries: map[string]bool{
			"db_persistence":                   false,
			"migration_added":                  false,
			"operations_console_change":        false,
			"public_api_added":                 false,
			"consumer_tracker_changed":         false,
			"external_predictor_runtime_added": false,
		},
		AggregateOnly:    true,
		RawRowsPersisted: false,
	}
	writePredictionLabBacktestJSON(t, filepath.Join(dir, "summary.json"), summary)
	writePredictionLabBacktestJSON(t, filepath.Join(dir, "metrics.json"), metrics)
	writePredictionLabBacktestJSON(t, filepath.Join(dir, "manifest.json"), manifest)
	for _, name := range []string{"summary.md", "metrics.md"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("private aggregate diagnostic summary\n"), 0o644); err != nil {
			t.Fatalf("write prediction lab backtest %s: %v", name, err)
		}
	}
}

func coveredPredictionLabConformanceReview() realtimequality.ConformanceReview {
	cases := []realtimequality.ConformanceCase{
		{ID: "after-midnight-service", Status: "synthetic_covered", Signal: "after-midnight covered", DoesNotProve: "synthetic only"},
		{ID: "frequency-headway-service", Status: "synthetic_covered", Signal: "frequency covered", DoesNotProve: "synthetic only"},
		{ID: "service-calendar-start-instance", Status: "synthetic_covered", Signal: "service calendar covered", DoesNotProve: "synthetic only"},
		{ID: "blocked-unknown-ambiguous", Status: "synthetic_covered", Signal: "withheld covered", DoesNotProve: "synthetic only"},
		{ID: "shadow-fail-closed", Status: "synthetic_covered", Signal: "shadow covered", DoesNotProve: "synthetic only"},
	}
	return realtimequality.ConformanceReview{
		Status:        "synthetic_covered",
		Boundary:      "Private synthetic conformance summary only.",
		SyntheticOnly: true,
		AggregateOnly: true,
		CaseCount:     len(cases),
		Cases:         cases,
		DoesNotProve:  "Synthetic conformance rows do not prove real-world ETA accuracy.",
	}
}

func writePredictionLabBacktestJSON(t *testing.T, path string, value any) {
	t.Helper()
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal prediction lab backtest fixture: %v", err)
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatalf("write prediction lab backtest json: %v", err)
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
	if count := strings.Count(body, "one-time-token"); count != 1 {
		t.Fatalf("POST body shows one-time token %d times, want exactly once: %s", count, body)
	}
	for _, forbidden := range []string{"token_hash", "raw_token", "authorization", "Bearer "} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("POST body leaks %q: %s", forbidden, body)
		}
	}
	if deviceStore.rebindCalls != 1 {
		t.Fatalf("Rebind called %d times, want 1", deviceStore.rebindCalls)
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/devices", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if strings.Contains(rr.Body.String(), "one-time-token") {
		t.Fatalf("GET body repeats one-time token: %s", rr.Body.String())
	}
}

func TestOperationsDevicesNonAdminsShowGuidanceWithoutRebindForm(t *testing.T) {
	now := time.Now().UTC()
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor} {
		t.Run(string(role), func(t *testing.T) {
			handler := newOperationsTestHandler(&handler{
				store: &fakePublicationStore{},
				devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{{
					AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now.Add(-time.Hour), CreatedAt: now.Add(-time.Hour),
				}}},
				telemetry: fakeTelemetryRepository{latest: []telemetry.StoredEvent{{
					Event: telemetry.Event{
						AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Timestamp: now.Add(-30 * time.Second), Lat: 1, Lon: 2,
					},
					ReceivedAt: now.Add(-29 * time.Second), IngestStatus: telemetry.IngestStatusAccepted,
				}}},
				state: fakeStateRepository{assignments: map[string]state.Assignment{"bus-1": {
					VehicleID: "bus-1", State: state.StateInService, RouteID: "route-1", TripID: "trip-1", Confidence: 0.91, ActiveFrom: now.Add(-25 * time.Second),
				}}},
			}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			req := httptest.NewRequest(http.MethodGet, "/admin/operations/devices", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
			}
			body := rr.Body.String()
			for _, want := range []string{"Guided Onboarding Use Cases", "Telemetry verification", "bus-1", "fresh", "in_service / route route-1 / trip trip-1 / confidence 0.91", "No immediate action", "Admins can rotate or rebind device tokens"} {
				if !strings.Contains(body, want) {
					t.Fatalf("body does not contain %q: %s", want, body)
				}
			}
			for _, forbidden := range []string{"<form", "Rotate / rebind token", "name=\"csrf_token\"", "name=\"agency_id\"", "name=\"device_id\"", "name=\"vehicle_id\"", "one-time-token"} {
				if strings.Contains(body, forbidden) {
					t.Fatalf("non-admin devices page exposes mutation control or token %q: %s", forbidden, body)
				}
			}
		})
	}

	admin := newOperationsTestHandler(&handler{
		store: &fakePublicationStore{},
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{{
			AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now.Add(-time.Hour), CreatedAt: now.Add(-time.Hour),
		}}},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/devices", nil)
	rr := httptest.NewRecorder()
	admin.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("admin status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	for _, want := range []string{"<form", "method=\"post\"", "name=\"csrf_token\"", "name=\"agency_id\"", "name=\"device_id\"", "name=\"vehicle_id\"", "Rotate / rebind token"} {
		if !strings.Contains(rr.Body.String(), want) {
			t.Fatalf("admin body missing %q: %s", want, rr.Body.String())
		}
	}
}

func TestOperationsDevicesSummarizesTelemetryAndAssignmentsSafely(t *testing.T) {
	now := time.Now().UTC()
	handler := newOperationsTestHandler(&handler{
		store: &fakePublicationStore{},
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{
			{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now.Add(-2 * time.Hour), CreatedAt: now.Add(-2 * time.Hour)},
			{AgencyID: "demo-agency", DeviceID: "device-2", VehicleID: "bus-2", Status: "active", ValidFrom: now.Add(-2 * time.Hour), CreatedAt: now.Add(-2 * time.Hour)},
		}},
		telemetry: fakeTelemetryRepository{latest: []telemetry.StoredEvent{
			{
				Event: telemetry.Event{
					AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Timestamp: now.Add(-30 * time.Second), Lat: 1, Lon: 2,
				},
				ReceivedAt: now.Add(-29 * time.Second), IngestStatus: telemetry.IngestStatusAccepted,
				PayloadJSON: []byte(`{"token_hash":"payload-secret-token","private_debug":true,"vendor_id":"payload-secret-vendor"}`),
			},
			{
				Event: telemetry.Event{
					AgencyID: "demo-agency", DeviceID: "device-2", VehicleID: "bus-2", Timestamp: now.Add(-5 * time.Minute), Lat: 3, Lon: 4,
				},
				ReceivedAt: now.Add(-5 * time.Minute), IngestStatus: telemetry.IngestStatusAccepted,
				PayloadJSON: []byte(`{"raw_payload":"payload-secret-raw"}`),
			},
		}},
		state: fakeStateRepository{assignments: map[string]state.Assignment{"bus-1": {
			VehicleID: "bus-1", State: state.StateInService, RouteID: "route-1", TripID: "trip-1", Confidence: 0.91,
			ScoreDetails: map[string]any{"private_debug": true}, ActiveFrom: now.Add(-25 * time.Second),
		}}},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/devices", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"bus-1", "fresh", "in_service / route route-1 / trip trip-1 / confidence 0.91", "No immediate action",
		"bus-2", "stale", "not available", "Check device power, network, and reporting cadence",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"PayloadJSON", "payload_json", "raw_payload", "token_hash", "private_debug", "vendor_id", "payload-secret-token", "payload-secret-vendor", "payload-secret-raw"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("devices page leaks private field %q: %s", forbidden, body)
		}
	}
}

func TestOperationsDevicesShowsFleetOnboardingV2GuidanceSafely(t *testing.T) {
	now := time.Now().UTC()
	handler := newOperationsTestHandler(&handler{
		store: &fakePublicationStore{},
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{
			{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now.Add(-2 * time.Hour), CreatedAt: now.Add(-2 * time.Hour)},
			{AgencyID: "demo-agency", DeviceID: "device-2", VehicleID: "bus-2", Status: "active", ValidFrom: now.Add(-2 * time.Hour), CreatedAt: now.Add(-2 * time.Hour)},
		}},
		telemetry: fakeTelemetryRepository{latest: []telemetry.StoredEvent{
			{
				Event: telemetry.Event{
					AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Timestamp: now.Add(-30 * time.Second), Lat: 1, Lon: 2,
				},
				ReceivedAt: now.Add(-29 * time.Second), IngestStatus: telemetry.IngestStatusAccepted,
				PayloadJSON: []byte(`{"token_hash":"payload-secret-token","private_debug":true,"vendor_id":"payload-secret-vendor"}`),
			},
			{
				Event: telemetry.Event{
					AgencyID: "demo-agency", DeviceID: "device-3", VehicleID: "bus-3", Timestamp: now.Add(-45 * time.Second), Lat: 3, Lon: 4,
				},
				ReceivedAt: now.Add(-44 * time.Second), IngestStatus: telemetry.IngestStatusAccepted,
				PayloadJSON: []byte(`{"raw_payload":"payload-secret-raw"}`),
			},
		}},
		state: fakeStateRepository{assignments: map[string]state.Assignment{"bus-1": {
			VehicleID: "bus-1", State: state.StateUnknown, Confidence: 0.40, ActiveFrom: now.Add(-25 * time.Second),
		}}},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/devices", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Fleet Onboarding V2 Review",
		"Vehicle / device inventory review",
		"fresh=1; stale=0; not_seen=1; unlisted_accepted_rows=1",
		"Bulk import plan",
		"the console does not import token values or generate bulk secrets",
		"Rotate/rebind is the supported credential action",
		"token recovery is intentionally unavailable",
		"Ingest Diagnostics",
		"Authenticated ingest contract",
		"safe quality flags such as stale_timestamp or low_gps_accuracy",
		"future timestamps are rejected before storage",
		"invalid motion fields",
		"Duplicate and out_of_order events are stored as non-accepted rows",
		"Not-seen device triage",
		"1 configured bindings have no latest accepted telemetry",
		"Unknown-device and rejected-payload triage",
		"Unauthorized device payloads are rejected before storage",
		"unknown_assignments=1; low_confidence_assignments=1; unlisted_accepted_rows=1",
		"Safe administrator handoff",
		"token values, request credential headers, private endpoints, raw telemetry, database URLs, or evidence packets",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"PayloadJSON", "payload_json", "raw_payload", "token_hash", "private_debug", "vendor_id", "payload-secret-token", "payload-secret-vendor", "payload-secret-raw", "Authorization:", "Bearer "} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("fleet onboarding leaks private field %q: %s", forbidden, body)
		}
	}
}

func TestOperationsDevicesAndTelemetryAvoidUnsupportedClaims(t *testing.T) {
	now := time.Now().UTC()
	handler := newOperationsTestHandler(&handler{
		store: &fakePublicationStore{},
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{{
			AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now.Add(-time.Hour), CreatedAt: now.Add(-time.Hour),
		}}},
		telemetry: fakeTelemetryRepository{latest: []telemetry.StoredEvent{{
			Event:      telemetry.Event{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Timestamp: now.Add(-30 * time.Second), Lat: 1, Lon: 2},
			ReceivedAt: now.Add(-29 * time.Second), IngestStatus: telemetry.IngestStatusAccepted,
		}}},
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	for _, path := range []string{"/admin/operations/devices", "/admin/operations/telemetry"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
		}
		rawBody := rr.Body.String()
		body := strings.ToLower(rawBody)
		for _, want := range []string{
			"creates no retained evidence",
			"contacts no vendors or consumers",
			"changes no consumer status",
			"does not show hardware certification, vendor compatibility, production AVL reliability, consumer acceptance, compliance, hosted service, or production readiness",
		} {
			if !strings.Contains(rawBody, want) {
				t.Fatalf("%s body missing boundary copy %q: %s", path, want, rawBody)
			}
		}
		for _, forbidden := range []string{
			"cal-itp/caltrans compliant",
			"consumer accepted",
			"accepted by",
			"agency approved",
			"agency adopted",
			"production ready",
			"public launch complete",
			"hosted saas",
			"uptime guarantee",
			"vendor compatible",
			"hardware certified",
			"production avl reliable",
			"real realtime",
			"production-grade eta",
		} {
			if strings.Contains(body, forbidden) {
				t.Fatalf("%s body contains unsupported claim %q: %s", path, forbidden, rr.Body.String())
			}
		}
	}
}

func TestOperationsTelemetrySimulatorGuideListsSyntheticScenariosSafely(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/telemetry-simulator", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Telemetry Simulator Guide",
		"can preview committed synthetic fixture metadata but executes no command",
		"reads no private diagnostics",
		"collects no device token",
		"Browser Dry-Run Preview",
		"Preview synthetic dry run",
		"ready_for_browser_preview",
		"Redacted synthetic event preview",
		"synthetic point",
		"This browser page never asks for, stores, displays, or posts device credentials.",
		"on-route",
		"stale",
		"out-of-order",
		"unknown-device",
		"low-quality-gps",
		"after-midnight",
		"block-transition",
		"SCENARIO=on-route DRY_RUN=true make telemetry-simulator",
		"First local/synthetic dry-run safety check",
		"does not test a live vendor",
		"SCENARIO=on-route make telemetry-simulator",
		"SCENARIO=on-route RUN_MATCHER=true make telemetry-simulator",
		"Detailed safety booleans remain available",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("simulator guide missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{
		`method="post"`,
		"name=\"device_token\"",
		"DEVICE_TOKEN=",
		"Authorization:",
		"Bearer ",
		"\"payload\"",
		"\"lat\"",
		"\"lon\"",
		"payload_json",
		"token_hash",
		"private_debug",
		"vendor_id",
	} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("simulator guide exposes forbidden token/payload/control text %q: %s", forbidden, body)
		}
	}
}

func TestOperationsTelemetrySimulatorJSONIsPrivateBoundedAndNoExecution(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/telemetry-simulator.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var view operationsTelemetrySimulatorView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode simulator JSON: %v: %s", err, rr.Body.String())
	}
	if view.AgencyID != "demo-agency" {
		t.Fatalf("agency = %q, want demo-agency", view.AgencyID)
	}
	if view.ScenarioDir != telemetrySimulatorScenarioDir {
		t.Fatalf("scenario dir = %q, want %q", view.ScenarioDir, telemetrySimulatorScenarioDir)
	}
	if view.LoadError != "" {
		t.Fatalf("load error = %q", view.LoadError)
	}
	if len(view.Scenarios) < 7 {
		t.Fatalf("scenario count = %d, want at least 7", len(view.Scenarios))
	}
	if view.Scenarios[0].Name != "on-route" || !view.Scenarios[0].DefaultLocal {
		t.Fatalf("first scenario = %+v, want default on-route first", view.Scenarios[0])
	}
	if view.SelectedScenario != "on-route" {
		t.Fatalf("selected scenario = %q, want on-route", view.SelectedScenario)
	}
	if view.DryRunPreview.Status != "ready_for_browser_preview" || view.DryRunPreview.ScenarioID != "on-route" || len(view.DryRunPreview.Events) == 0 {
		t.Fatalf("invalid dry-run preview: %+v", view.DryRunPreview)
	}
	if !strings.Contains(view.DryRunPreview.Boundary, "does not execute shell commands") {
		t.Fatalf("dry-run boundary should explicitly avoid shell execution: %+v", view.DryRunPreview)
	}
	flags := view.ClaimFlags
	if flags.BackendCommandExecutionEnabled || flags.TelemetrySentByWebRequest || flags.DeviceTokenCollectedByBrowser || flags.CacheDiagnosticsReadEnabled ||
		flags.ExternalEvidenceCreated || flags.ConsumerStatusesChanged || flags.VendorCompatibilityClaimed || flags.HardwareCertificationClaimed ||
		flags.ProductionAVLClaimed || flags.RealRealtimeClaimed || flags.ProductionGradeETAClaimed || flags.ComplianceClaimed {
		t.Fatalf("claim flags should all be false: %+v", flags)
	}
	body := rr.Body.String()
	for _, forbidden := range []string{"\"payload\"", "\"lat\"", "\"lon\"", "DEVICE_TOKEN=", "Authorization:", "Bearer ", "token_hash", "private_debug", "vendor_id"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("simulator JSON exposes forbidden field or secret-like text %q: %s", forbidden, body)
		}
	}

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/admin/operations/telemetry-simulator", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s status = %d, want 405", method, rr.Code)
		}
	}
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/telemetry-simulator?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("cross-agency status = %d, want 403", rr.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations/telemetry-simulator", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("public simulator route status = %d, want 404", rr.Code)
	}
}

func TestOperationsTelemetrySimulatorBrowserDryRunSelectionIsPreviewOnly(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/telemetry-simulator.json?scenario=stale", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var view operationsTelemetrySimulatorView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode simulator JSON: %v: %s", err, rr.Body.String())
	}
	if view.SelectedScenario != "stale" || view.DryRunPreview.ScenarioID != "stale" {
		t.Fatalf("scenario selection did not drive preview: selected=%q preview=%+v", view.SelectedScenario, view.DryRunPreview)
	}
	if view.DryRunPreview.Status != "ready_for_browser_preview" || len(view.DryRunPreview.Events) == 0 {
		t.Fatalf("dry-run preview should be ready with redacted events: %+v", view.DryRunPreview)
	}
	payload := rr.Body.String()
	for _, forbidden := range []string{"\"payload\"", "\"lat\"", "\"lon\"", "DEVICE_TOKEN=", "Authorization:", "Bearer ", "token_hash", "private_debug", "vendor_id", "backend_command_execution_enabled\":true", "telemetry_sent_by_web_request\":true"} {
		if strings.Contains(payload, forbidden) {
			t.Fatalf("browser preview exposed forbidden field, secret-like text, or command/send flag %q: %s", forbidden, payload)
		}
	}
}

func TestOperationsTelemetrySimulatorAvoidsUnsupportedClaims(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
	for _, path := range []string{"/admin/operations/telemetry-simulator", "/admin/operations/telemetry-simulator.json"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
		}
		body := strings.ToLower(rr.Body.String())
		for _, forbidden := range []string{
			"cal-itp/caltrans compliant",
			"consumer accepted",
			"accepted by",
			"agency approved",
			"agency adopted",
			"production ready",
			"public launch complete",
			"hosted saas",
			"uptime guarantee",
			"vendor compatible",
			"hardware certified",
			"production avl reliable",
			"real realtime",
			"production-grade eta quality is proven",
		} {
			if strings.Contains(body, forbidden) {
				t.Fatalf("%s body contains unsupported claim %q: %s", path, forbidden, rr.Body.String())
			}
		}
	}
}

func TestOperationsDeviceRebindRequiresAdminRole(t *testing.T) {
	for _, role := range []auth.Role{auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor} {
		t.Run(string(role), func(t *testing.T) {
			deviceStore := &fakeDeviceStoreWithToken{token: "forbidden-token"}
			handler := newOperationsTestHandler(&handler{store: &fakePublicationStore{}, devices: deviceStore}, auth.TestAuthenticator{Principal: auth.Principal{
				Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{role}, Method: auth.MethodBearer,
			}})
			req := httptest.NewRequest(http.MethodPost, "/admin/operations/devices", strings.NewReader("agency_id=demo-agency&device_id=device-1&vehicle_id=bus-1"))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusForbidden {
				t.Fatalf("status = %d, want 403", rr.Code)
			}
			if deviceStore.rebindCalls != 0 {
				t.Fatalf("Rebind called %d times for non-admin role %s", deviceStore.rebindCalls, role)
			}
			if strings.Contains(rr.Body.String(), "forbidden-token") {
				t.Fatalf("forbidden response leaked token: %s", rr.Body.String())
			}
		})
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
	for _, want := range []string{
		"Prepared-Only Consumer Packet Explanation",
		"Private prepared-only consumer packet review",
		"7 of 7 docs tracker targets are visible as prepared-only records.",
		"Target Boundary Review",
		"Future Authorization Gates",
		"Workflow Separation",
		"Google Maps",
		"Apple Maps",
		"Transit App",
		"Bing Maps",
		"Moovit",
		"Mobility Database",
		"transit.land",
		"not_started",
		"docs/evidence tracker",
		"Runtime deployment note is",
		"docs tracker status remains",
		"Requires separate written authorization",
		"Detailed safety booleans remain available",
		"every docs/evidence consumer target at prepared",
		"Database workflow notes are shown separately",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"accepted by", "submission complete", "ingestion confirmed", "listed by consumer", "displayed by consumer", "production ready", "hosted SaaS", "vendor compatible", "certified hardware", "database_url", "Bearer "} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("body invents forbidden claim %q: %s", forbidden, body)
		}
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

func assertSetupWizardShape(t *testing.T, wizard operationsSetupWizardView) {
	t.Helper()
	if wizard.AgencyID == "" || wizard.Boundary == "" || wizard.Summary.Status == "" || wizard.Summary.NextStageID == "" || wizard.Summary.NextStageLabel == "" || wizard.Summary.NextAction == "" || wizard.Summary.NextActionLink == "" || wizard.Summary.Meaning == "" || len(wizard.Blockers) == 0 || len(wizard.Diagnostics) != 8 || len(wizard.RoleVisibility) != 3 || len(wizard.TechnicalHelp) != 4 || len(wizard.Stages) != 8 || wizard.Counts.Stages != len(wizard.Stages) {
		t.Fatalf("invalid setup wizard top-level shape: %+v", wizard)
	}
	allowedStatuses := map[string]bool{"ok": true, "needs_review": true, "missing": true, "blocked": true, "unknown": true}
	seenIDs := map[string]bool{}
	for _, stage := range wizard.Stages {
		if stage.ID == "" || stage.Label == "" || stage.Status == "" || stage.CurrentSignal == "" || stage.PrimaryAction == "" || stage.ActionLabel == "" || stage.AdminLink == "" || len(stage.DocsLinks) == 0 || stage.ClaimBoundary == "" {
			t.Fatalf("invalid setup wizard stage shape: %+v", stage)
		}
		if seenIDs[stage.ID] {
			t.Fatalf("duplicate setup wizard stage id %q", stage.ID)
		}
		seenIDs[stage.ID] = true
		if !allowedStatuses[stage.Status] {
			t.Fatalf("stage %q status = %q, want neutral status", stage.ID, stage.Status)
		}
		if !strings.HasPrefix(stage.AdminLink, "/admin/") {
			t.Fatalf("stage %s has unsafe admin link %q", stage.ID, stage.AdminLink)
		}
		for _, link := range stage.DocsLinks {
			if !strings.HasPrefix(link, "docs/") {
				t.Fatalf("stage %s has unsafe docs link %q", stage.ID, link)
			}
		}
	}
	for _, blocker := range wizard.Blockers {
		if blocker.StageID == "" || blocker.Label == "" || blocker.Status == "" || blocker.CurrentSignal == "" || blocker.NextAction == "" || blocker.ActionLabel == "" || blocker.AdminLink == "" {
			t.Fatalf("invalid setup wizard blocker shape: %+v", blocker)
		}
		if !strings.HasPrefix(blocker.AdminLink, "/admin/") {
			t.Fatalf("blocker %s has unsafe admin link %q", blocker.StageID, blocker.AdminLink)
		}
	}
	for _, diagnostic := range wizard.Diagnostics {
		if diagnostic.ID == "" || diagnostic.Label == "" || diagnostic.Status == "" || diagnostic.CurrentSignal == "" || diagnostic.NextAction == "" || diagnostic.ClaimBoundary == "" {
			t.Fatalf("invalid setup wizard diagnostic shape: %+v", diagnostic)
		}
		if !allowedStatuses[diagnostic.Status] {
			t.Fatalf("diagnostic %q status = %q, want neutral status", diagnostic.ID, diagnostic.Status)
		}
	}
	for _, role := range wizard.RoleVisibility {
		if role.ID == "" || role.Label == "" || role.Status == "" || role.CurrentSignal == "" || role.NextAction == "" || role.ClaimBoundary == "" {
			t.Fatalf("invalid setup wizard role shape: %+v", role)
		}
		if !allowedStatuses[role.Status] {
			t.Fatalf("role visibility %q status = %q, want neutral status", role.ID, role.Status)
		}
	}
	for _, help := range wizard.TechnicalHelp {
		if help.ID == "" || help.Label == "" || help.WhenNeeded == "" || help.NextAction == "" || help.AdminLink == "" || help.DocsLink == "" || help.ClaimBoundary == "" {
			t.Fatalf("invalid setup wizard administrator shape: %+v", help)
		}
		if !strings.HasPrefix(help.AdminLink, "/admin/") {
			t.Fatalf("administrator %s has unsafe admin link %q", help.ID, help.AdminLink)
		}
		if !strings.HasPrefix(help.DocsLink, "docs/") {
			t.Fatalf("administrator %s has unsafe docs link %q", help.ID, help.DocsLink)
		}
	}
}

func assertSetupWizardFlagsFalse(t *testing.T, flags setupWizardClaimFlags) {
	t.Helper()
	if flags.ExternalEvidenceCreated || flags.FinalRootEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.AgencyApprovalClaimed || flags.ConsumerAcceptanceClaimed || flags.PublicLaunchClaimed || flags.HostedSaaSClaimed || flags.VendorCompatibilityClaimed || flags.ProductionGradeETAClaimed {
		t.Fatalf("setup wizard flags must all be false: %+v", flags)
	}
}

func assertFirstRunShape(t *testing.T, firstRun operationsFirstRunView) {
	t.Helper()
	if firstRun.AgencyID == "" || firstRun.Boundary == "" || firstRun.LocalDemoDeploymentEvidenceBoundary == "" || len(firstRun.Paths) != 2 || len(firstRun.Tasks) != 9 || len(firstRun.FeedURLs) != 5 || firstRun.Counts.Tasks != len(firstRun.Tasks) || firstRun.Counts.FeedURLs != len(firstRun.FeedURLs) {
		t.Fatalf("invalid first-run top-level shape: %+v", firstRun)
	}
	allowedStatuses := map[string]bool{"ok": true, "needs_review": true, "missing": true, "blocked": true, "unknown": true}
	wantTaskIDs := []string{"metadata", "gtfs", "five_feed_urls", "validation_health", "telemetry", "vp_tu_alerts", "readiness", "connectors", "support_checks"}
	var gotTaskIDs []string
	for index, task := range firstRun.Tasks {
		gotTaskIDs = append(gotTaskIDs, task.ID)
		if task.Order != index+1 || task.ID == "" || task.Label == "" || task.Status == "" || task.CurrentSignal == "" || task.Meaning == "" || task.NextAction == "" || task.UILink == "" || task.DocsLink == "" || task.DoesNotProve == "" {
			t.Fatalf("invalid first-run task shape: %+v", task)
		}
		if !allowedStatuses[task.Status] {
			t.Fatalf("first-run task %q status = %q, want neutral status", task.ID, task.Status)
		}
		if !strings.HasPrefix(task.UILink, "/admin/") {
			t.Fatalf("first-run task %s unsafe UI link %q", task.ID, task.UILink)
		}
		if !strings.HasPrefix(task.DocsLink, "docs/") {
			t.Fatalf("first-run task %s unsafe docs link %q", task.ID, task.DocsLink)
		}
	}
	if strings.Join(gotTaskIDs, ",") != strings.Join(wantTaskIDs, ",") {
		t.Fatalf("first-run task ids = %v, want %v", gotTaskIDs, wantTaskIDs)
	}
	wantPathIDs := []string{"no_code", "developer"}
	var gotPathIDs []string
	for _, path := range firstRun.Paths {
		gotPathIDs = append(gotPathIDs, path.ID)
		if path.ID == "" || path.Label == "" || path.CurrentSignal == "" || path.Meaning == "" || path.FirstAction == "" || path.UILink == "" || path.DocsLink == "" || path.DoesNotProve == "" {
			t.Fatalf("invalid first-run path shape: %+v", path)
		}
		if !strings.HasPrefix(path.UILink, "/admin/") {
			t.Fatalf("first-run path %s unsafe UI link %q", path.ID, path.UILink)
		}
		if !strings.HasPrefix(path.DocsLink, "docs/") {
			t.Fatalf("first-run path %s unsafe docs link %q", path.ID, path.DocsLink)
		}
	}
	if strings.Join(gotPathIDs, ",") != strings.Join(wantPathIDs, ",") {
		t.Fatalf("first-run path ids = %v, want %v", gotPathIDs, wantPathIDs)
	}
	wantFeedIDs := []string{"feeds_json", "schedule", "vehicle_positions", "trip_updates", "alerts"}
	var gotFeedIDs []string
	for _, row := range firstRun.FeedURLs {
		gotFeedIDs = append(gotFeedIDs, row.ID)
		if row.ID == "" || row.Label == "" || row.Status == "" || row.Source == "" || row.Meaning == "" || row.NextAction == "" || row.DocsLink == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid first-run feed URL shape: %+v", row)
		}
		if row.URL == "" && row.CopyValue != "" {
			t.Fatalf("first-run feed %s has missing URL but copy value %q", row.ID, row.CopyValue)
		}
		if row.URL != "" && row.CopyValue == "" {
			t.Fatalf("first-run feed %s has URL but no copy value", row.ID)
		}
		if !allowedStatuses[row.Status] {
			t.Fatalf("first-run feed %q status = %q, want neutral status", row.ID, row.Status)
		}
		if !strings.HasPrefix(row.DocsLink, "docs/") {
			t.Fatalf("first-run feed %s unsafe docs link %q", row.ID, row.DocsLink)
		}
	}
	if strings.Join(gotFeedIDs, ",") != strings.Join(wantFeedIDs, ",") {
		t.Fatalf("first-run feed ids = %v, want %v", gotFeedIDs, wantFeedIDs)
	}
}

func assertFirstRunFlagsFalse(t *testing.T, flags operationsFirstRunClaimFlags) {
	t.Helper()
	if flags.BackendCommandExecutionEnabled || flags.CacheDiagnosticsRead || flags.ExternalNetworkContacted || flags.ExternalEvidenceCreated || flags.FinalRootEvidenceCreated || flags.ConsumerStatusesChanged || flags.SecretsCollected || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.AgencyApprovalClaimed || flags.ConsumerAcceptanceClaimed || flags.PublicLaunchClaimed || flags.HostedSaaSClaimed || flags.VendorCompatibilityClaimed || flags.HardwareCertificationClaimed || flags.ProductionAVLReliabilityClaimed || flags.ProductionGradeETAQualityClaimed || flags.SLAClaimed || flags.UptimeGuaranteeClaimed || flags.DynamicBackendPluginLoadingEnabled || flags.ReleaseCandidateApprovalClaimed || flags.ManagedSupportCommitmentClaimed || flags.FinalDeploymentOwnershipClaimed || flags.ConsumerIngestionWorkflowCompleted || flags.ProductionMultiTenantHostingClaimed {
		t.Fatalf("first-run flags must all be false: %+v", flags)
	}
}

func assertOperationsHelpShape(t *testing.T, view operationsHelpView) {
	t.Helper()
	if view.GeneratedAt.IsZero() || view.AgencyID == "" || view.Boundary == "" || view.TrainingGuide.DocsPath == "" || len(view.Topics) != 7 || len(view.RoleTours) != 5 || len(view.FirstWeek) != 7 || len(view.Glossary) != 11 || len(view.Recovery) != 8 || len(view.QuickTasks) != 7 || len(view.Handoff) != 6 || len(view.DemoScenarios) != 6 || len(view.TrainerScript) != 6 || len(view.HelperChecklist) != 7 {
		t.Fatalf("invalid help top-level shape: %+v", view)
	}
	if view.TrainingGuide.DocsPath != "docs/operator-training-guide.md" || view.TrainingGuide.Label == "" || view.TrainingGuide.Audience == "" || view.TrainingGuide.HowToUse == "" || view.TrainingGuide.Boundary == "" {
		t.Fatalf("invalid training guide link: %+v", view.TrainingGuide)
	}
	wantRoles := []string{"no_code_evaluator", "director_manager", "daily_operator", "administrator", "integrator"}
	var gotRoles []string
	for _, tour := range view.RoleTours {
		gotRoles = append(gotRoles, tour.ID)
		if tour.ID == "" || tour.Label == "" || tour.Who == "" || tour.StartHere == "" || tour.ReviewFirst == "" || tour.FirstActions == "" || tour.EscalateWhen == "" || tour.DoesNotProve == "" || len(tour.AdminLinks) == 0 || len(tour.DocsLinks) == 0 {
			t.Fatalf("invalid role tour shape: %+v", tour)
		}
		if !strings.HasPrefix(tour.StartHere, "/admin/") {
			t.Fatalf("role tour %s has unsafe start link %q", tour.ID, tour.StartHere)
		}
		for _, link := range tour.AdminLinks {
			if !strings.HasPrefix(link, "/admin/") {
				t.Fatalf("role tour %s has unsafe admin link %q", tour.ID, link)
			}
		}
		for _, link := range tour.DocsLinks {
			if !strings.HasPrefix(link, "docs/") {
				t.Fatalf("role tour %s has unsafe docs link %q", tour.ID, link)
			}
		}
	}
	if strings.Join(gotRoles, ",") != strings.Join(wantRoles, ",") {
		t.Fatalf("role ids = %v, want %v", gotRoles, wantRoles)
	}
	for _, item := range view.FirstWeek {
		if item.ID == "" || item.Day == "" || item.Role == "" || item.Task == "" || item.Review == "" || item.DoneWhen == "" || item.NextAction == "" || item.ConsoleLink == "" || item.DoesNotProve == "" {
			t.Fatalf("invalid first-week item shape: %+v", item)
		}
		if !strings.HasPrefix(item.ConsoleLink, "/admin/") {
			t.Fatalf("first-week item %s has unsafe console link %q", item.ID, item.ConsoleLink)
		}
	}
	for _, term := range view.Glossary {
		if term.ID == "" || term.Term == "" || term.PlainMeaning == "" || term.TechnicalMeaning == "" || term.WhereToReview == "" || term.DoesNotProve == "" || len(term.DocsLinks) == 0 {
			t.Fatalf("invalid glossary term shape: %+v", term)
		}
		for _, link := range term.DocsLinks {
			if !strings.HasPrefix(link, "docs/") {
				t.Fatalf("glossary term %s has unsafe docs link %q", term.ID, link)
			}
		}
	}
	for _, row := range view.Recovery {
		if row.ID == "" || row.WhatOperatorSees == "" || row.LikelyCause == "" || row.SafeNextStep == "" || row.EscalationTrigger == "" || row.ConsoleLink == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid recovery row shape: %+v", row)
		}
		if !strings.HasPrefix(row.ConsoleLink, "/admin/") {
			t.Fatalf("recovery row %s has unsafe console link %q", row.ID, row.ConsoleLink)
		}
	}
	for _, task := range view.QuickTasks {
		if task.ID == "" || task.Label == "" || task.PrimaryRole == "" || task.ConsoleLink == "" || task.ReviewSteps == "" || task.DoneWhen == "" || task.Escalation == "" || task.DoesNotProve == "" {
			t.Fatalf("invalid quick task shape: %+v", task)
		}
		if !strings.HasPrefix(task.ConsoleLink, "/admin/") {
			t.Fatalf("quick task %s has unsafe console link %q", task.ID, task.ConsoleLink)
		}
	}
	for _, item := range view.Handoff {
		if item.ID == "" || item.Area == "" || item.FromRole == "" || item.ToRole == "" || item.Confirm == "" || item.ConsoleLink == "" || item.DoesNotProve == "" {
			t.Fatalf("invalid handoff item shape: %+v", item)
		}
		if !strings.HasPrefix(item.ConsoleLink, "/admin/") {
			t.Fatalf("handoff item %s has unsafe console link %q", item.ID, item.ConsoleLink)
		}
	}
	for _, scenario := range view.DemoScenarios {
		if scenario.ID == "" || scenario.Label == "" || scenario.Audience == "" || len(scenario.FixturePaths) == 0 || scenario.ConsoleLink == "" || scenario.Exercise == "" || scenario.DoneWhen == "" || scenario.RecoveryPrompt == "" || scenario.DoesNotProve == "" {
			t.Fatalf("invalid demo scenario shape: %+v", scenario)
		}
		if !strings.HasPrefix(scenario.ConsoleLink, "/admin/") {
			t.Fatalf("demo scenario %s has unsafe console link %q", scenario.ID, scenario.ConsoleLink)
		}
		for _, link := range scenario.FixturePaths {
			if !strings.HasPrefix(link, "testdata/") || strings.Contains(link, "..") {
				t.Fatalf("demo scenario %s has unsafe fixture link %q", scenario.ID, link)
			}
		}
	}
	for _, step := range view.TrainerScript {
		if step.ID == "" || step.Segment == "" || step.Minutes == "" || step.TalkTrack == "" || step.AskParticipant == "" || step.ConsoleLink == "" || step.Boundary == "" {
			t.Fatalf("invalid trainer step shape: %+v", step)
		}
		if !strings.HasPrefix(step.ConsoleLink, "/admin/") {
			t.Fatalf("trainer step %s has unsafe console link %q", step.ID, step.ConsoleLink)
		}
	}
	for _, item := range view.HelperChecklist {
		if item.ID == "" || item.Area == "" || item.Collect == "" || item.DoNotCollect == "" || item.ConsoleLink == "" || item.DocsLink == "" || item.NeedsAuthorizationWhen == "" {
			t.Fatalf("invalid helper checklist item shape: %+v", item)
		}
		if !strings.HasPrefix(item.ConsoleLink, "/admin/") {
			t.Fatalf("helper checklist item %s has unsafe console link %q", item.ID, item.ConsoleLink)
		}
		if !strings.HasPrefix(item.DocsLink, "docs/") {
			t.Fatalf("helper checklist item %s has unsafe docs link %q", item.ID, item.DocsLink)
		}
	}
	wantIDs := []string{"gtfs", "gtfs_rt", "connectors", "readiness", "validators", "telemetry", "claims_evidence"}
	var gotIDs []string
	for _, topic := range view.Topics {
		gotIDs = append(gotIDs, topic.ID)
		if topic.ID == "" || topic.Label == "" || topic.Summary == "" || topic.WhatToReview == "" || topic.NextAction == "" || topic.DoesNotProve == "" || topic.ClaimBoundary == "" || len(topic.AdminLinks) == 0 || len(topic.DocsLinks) == 0 {
			t.Fatalf("invalid help topic shape: %+v", topic)
		}
		for _, link := range topic.AdminLinks {
			if !strings.HasPrefix(link, "/admin/") {
				t.Fatalf("topic %s has unsafe admin link %q", topic.ID, link)
			}
		}
		for _, link := range topic.DocsLinks {
			if !strings.HasPrefix(link, "docs/") {
				t.Fatalf("topic %s has unsafe docs link %q", topic.ID, link)
			}
		}
		if topic.ID == "connectors" && topic.PluginDefinition != safePluginDefinition {
			t.Fatalf("connectors topic plugin definition = %q, want safe definition", topic.PluginDefinition)
		}
	}
	if strings.Join(gotIDs, ",") != strings.Join(wantIDs, ",") {
		t.Fatalf("topic ids = %v, want %v", gotIDs, wantIDs)
	}
}

func assertOperationsHelpFlagsFalse(t *testing.T, flags operationsHelpClaimFlags) {
	t.Helper()
	if flags.BackendCommandExecutionEnabled || flags.CacheDiagnosticsRead || flags.ExternalNetworkContacted || flags.ExternalEvidenceCreated || flags.FinalRootEvidenceCreated || flags.ConsumerStatusesChanged || flags.SecretsCollected || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.AgencyApprovalClaimed || flags.ConsumerAcceptanceClaimed || flags.PublicLaunchClaimed || flags.HostedSaaSClaimed || flags.VendorCompatibilityClaimed || flags.HardwareCertificationClaimed || flags.ProductionAVLReliabilityClaimed || flags.ProductionGradeETAQualityClaimed || flags.SLAClaimed || flags.UptimeGuaranteeClaimed || flags.DynamicBackendPluginLoadingEnabled {
		t.Fatalf("help flags must all be false: %+v", flags)
	}
}

func assertOperationsHelpSafeStrings(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"raw-token-value", "authorization:", "set-cookie", ".cache", "database_url", "restore_database_url", "payload_json", "raw telemetry", "token_hash", "file://", "/users/", "/opt/open-transit-rt", "/var/lib", "/etc/", "postgres://", "raw_report", "stdout", "stderr", "argv", "portal automation"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("help leaks forbidden private string %q: %s", forbidden, body)
		}
	}
	for _, forbidden := range []string{"agency_approved", "final_root_approved", "consumer_ready", "production_ready", "public_launch_complete", "compliance_achieved", "vendor_compatible", "hardware_certified", "dynamic_plugin_loading_enabled"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("help emits forbidden label %q: %s", forbidden, body)
		}
	}
}

func assertConnectorHubShape(t *testing.T, hub connectorHubView) {
	t.Helper()
	if hub.AgencyID == "" || hub.Boundary == "" || hub.PluginDefinition == "" || len(hub.Health) != 6 || len(hub.Catalog) != 28 || len(hub.Categories) != 6 || len(hub.Registry.Entries) != 8 {
		t.Fatalf("invalid connector hub top-level shape: %+v", hub)
	}
	seenHealthIDs := map[string]bool{}
	for _, row := range hub.Health {
		if row.ID == "" || row.Label == "" || row.Owner == "" || row.Status == "" || row.Configured == "" || row.DryRunReady == "" || row.SendState == "" || row.LastSyntheticCheck == "" || row.RedactionStatus == "" || row.IssueCategory == "" || len(row.IssueLinks) == 0 || len(row.SetupChecklist) == 0 || row.ChecklistCopy == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid connector health row shape: %+v", row)
		}
		if seenHealthIDs[row.ID] {
			t.Fatalf("duplicate connector health id %q", row.ID)
		}
		seenHealthIDs[row.ID] = true
		for _, link := range row.IssueLinks {
			if !strings.HasPrefix(link, "/admin/") {
				t.Fatalf("connector health row %s has unsafe issue link %q", row.ID, link)
			}
		}
		assertConnectorHealthChecklistSafe(t, row)
	}
	seenCatalogIDs := map[string]bool{}
	for _, row := range hub.Catalog {
		if row.ID == "" || row.Group == "" || row.Label == "" || row.Status == "" || row.StartWith == "" || row.BrowserReview == "" || row.FirstSafeCheck == "" || row.DoesNotProve == "" || len(row.DocsLinks) == 0 {
			t.Fatalf("invalid connector catalog row shape: %+v", row)
		}
		if seenCatalogIDs[row.ID] {
			t.Fatalf("duplicate connector catalog id %q", row.ID)
		}
		seenCatalogIDs[row.ID] = true
		for _, link := range row.DocsLinks {
			if !strings.HasPrefix(link, "docs/") {
				t.Fatalf("catalog row %s has unsafe docs link %q", row.ID, link)
			}
		}
	}
	seenIDs := map[string]bool{}
	for _, category := range hub.Categories {
		if category.ID == "" || category.Label == "" || category.Status == "" || category.Summary == "" || category.ConnectorShape == "" || len(category.Inputs) == 0 || len(category.Outputs) == 0 || category.FailureBehavior == "" || category.ClaimBoundary == "" || category.DocsLinks == nil || category.CommandSuggestions == nil || category.AdminLinks == nil {
			t.Fatalf("invalid connector category shape: %+v", category)
		}
		if seenIDs[category.ID] {
			t.Fatalf("duplicate connector category id %q", category.ID)
		}
		seenIDs[category.ID] = true
		for _, link := range category.DocsLinks {
			if !strings.HasPrefix(link, "docs/") {
				t.Fatalf("category %s has unsafe docs link %q", category.ID, link)
			}
		}
		for _, link := range category.AdminLinks {
			if !strings.HasPrefix(link, "/admin/") {
				t.Fatalf("category %s has unsafe admin link %q", category.ID, link)
			}
		}
	}
	seenRegistryIDs := map[string]bool{}
	for _, entry := range hub.Registry.Entries {
		if entry.SourcePath == "" || entry.SchemaVersion == "" || entry.ConnectorID == "" || entry.ConnectorType == "" || entry.DisplayName == "" || entry.Description == "" || entry.ModeName == "" || entry.DocsLink == "" || len(entry.InputContracts) == 0 || len(entry.OutputContracts) == 0 || len(entry.ConformanceCases) == 0 {
			t.Fatalf("invalid connector registry entry shape: %+v", entry)
		}
		if seenRegistryIDs[entry.ConnectorID] {
			t.Fatalf("duplicate connector registry id %q", entry.ConnectorID)
		}
		seenRegistryIDs[entry.ConnectorID] = true
		if !strings.HasPrefix(entry.SourcePath, "examples/connectors/") || strings.Contains(entry.SourcePath, "..") || strings.HasPrefix(entry.SourcePath, "/") {
			t.Fatalf("registry entry %s has unsafe source path %q", entry.ConnectorID, entry.SourcePath)
		}
		if !strings.HasPrefix(entry.DocsLink, "examples/connectors/") && !strings.HasPrefix(entry.DocsLink, "docs/") {
			t.Fatalf("registry entry %s has unsafe docs link %q", entry.ConnectorID, entry.DocsLink)
		}
	}
}

func assertConnectorHubFlagsFalse(t *testing.T, flags connectorHubClaimFlags) {
	t.Helper()
	if flags.DynamicBackendPluginLoadingEnabled || flags.VendorCompatibilityClaimed || flags.HardwareCertificationClaimed || flags.ConsumerStatusesChanged || flags.ExternalEvidenceCreated || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.HostedSaaSClaimed || flags.ProductionGradeETAClaimed {
		t.Fatalf("connector hub flags must all be false: %+v", flags)
	}
}

func assertConnectorHealthChecklistSafe(t *testing.T, row connectorHealthRow) {
	t.Helper()
	joined := strings.ToLower(strings.Join(append(append([]string{}, row.SetupChecklist...), row.ChecklistCopy), "\n"))
	for _, forbidden := range []string{"authorization:", "bearer ", "set-cookie", "database_url", "token_hash", "payload_json", "file://", "/users/", "/home/", "/var/", "/tmp/", "/etc/", "postgres://", "localhost", "127.0.0.1", "192.168.", "10.0.0.", ".local", "http://"} {
		if strings.Contains(joined, forbidden) {
			t.Fatalf("connector health checklist %s leaks forbidden %q: %s", row.ID, forbidden, joined)
		}
	}
	for _, forbidden := range []string{"api_key", "private_key", "password", "credential"} {
		if strings.Contains(joined, forbidden) {
			t.Fatalf("connector health checklist %s includes secret-like key %q: %s", row.ID, forbidden, joined)
		}
	}
}

func assertConnectorTestsShape(t *testing.T, view connectorTestsView) {
	t.Helper()
	if view.GeneratedAt.IsZero() || view.AgencyID == "" || view.Boundary == "" || len(view.Commands) != 9 {
		t.Fatalf("invalid connector tests top-level shape: %+v", view)
	}
	seenIDs := map[string]bool{}
	for _, command := range view.Commands {
		if command.ID == "" || command.Label == "" || command.CommandLine == "" || command.Validates == "" || command.Inputs == "" || command.FailureNextAction == "" || command.DoesNotProve == "" || len(command.DocsLinks) == 0 {
			t.Fatalf("invalid connector test command shape: %+v", command)
		}
		if seenIDs[command.ID] {
			t.Fatalf("duplicate connector test command id %q", command.ID)
		}
		seenIDs[command.ID] = true
		for _, link := range command.DocsLinks {
			if !strings.HasPrefix(link, "docs/") && !strings.HasPrefix(link, "examples/") {
				t.Fatalf("command %s has unsafe docs link %q", command.ID, link)
			}
		}
	}
}

func assertConnectorTestsFlagsFalse(t *testing.T, flags connectorTestsClaimFlags) {
	t.Helper()
	if flags.BackendCommandExecutionEnabled || flags.ManifestCommandExecutionEnabled || flags.ExternalNetworkContacted || flags.ExternalEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.VendorCompatibilityClaimed || flags.ProductionReadinessClaimed || flags.ProductionGradeETAClaimed {
		t.Fatalf("connector test flags must all be false: %+v", flags)
	}
}

func assertConnectorWorkbenchShape(t *testing.T, view connectorWorkbenchView) {
	t.Helper()
	if view.GeneratedAt.IsZero() || view.AgencyID == "" || view.Boundary == "" || len(view.DecisionTree) != 8 || len(view.Recipes) != 8 || len(view.RedactionTemplates) != 5 || len(view.DryRunCommands) != 5 || view.TelemetryPreview.Boundary == "" || len(view.TelemetryPreview.Sources) != 2 || len(view.TelemetryPreview.Rows) != 6 || view.WebhookBoundary.Title == "" || len(view.WebhookBoundary.Rows) != 4 || len(view.WebhookBoundary.DocsLinks) != 3 || view.PredictionGuide.Title == "" || len(view.PredictionGuide.Rows) != 3 || len(view.PredictionGuide.DocsLinks) != 3 || view.MonitoringGuide.Title == "" || len(view.MonitoringGuide.Rows) != 3 || len(view.MonitoringGuide.DocsLinks) != 3 || view.ConsumerGuide.Title == "" || len(view.ConsumerGuide.Rows) != 3 || len(view.ConsumerGuide.DocsLinks) != 3 || view.Conformance.Boundary == "" || view.Conformance.SuitePath != "testdata/adapter-conformance/suite.json" || view.Conformance.Status == "" || !view.Conformance.SyntheticOnly || view.Conformance.ManifestCount != 13 || view.Conformance.CaseCount != 25 || len(view.Conformance.Groups) != 5 || len(view.Conformance.RunnerCommands) != 4 || view.ManifestReview.Title == "" || view.ManifestReview.PluginDefinition != safePluginDefinition || len(view.ManifestReview.Rows) != 8 || len(view.ManifestReview.LintChecks) != 5 {
		t.Fatalf("invalid connector workbench top-level shape: %+v", view)
	}
	seenDecisions := map[string]bool{}
	for _, row := range view.DecisionTree {
		if row.ID == "" || row.SourceSignal == "" || row.UseWhen == "" || row.Boundary == "" || row.FirstSafeCheck == "" || row.StopIf == "" || row.NextAdminLink == "" || len(row.DocsLinks) == 0 || row.DoesNotProve == "" {
			t.Fatalf("invalid connector workbench decision row: %+v", row)
		}
		if seenDecisions[row.ID] {
			t.Fatalf("duplicate connector workbench decision row %q", row.ID)
		}
		seenDecisions[row.ID] = true
		if !strings.HasPrefix(row.NextAdminLink, "/admin/") {
			t.Fatalf("decision row %s has unsafe admin link %q", row.ID, row.NextAdminLink)
		}
		for _, link := range row.DocsLinks {
			if !strings.HasPrefix(link, "docs/") && !strings.HasPrefix(link, "examples/") {
				t.Fatalf("decision row %s has unsafe docs link %q", row.ID, link)
			}
		}
	}
	seenRecipes := map[string]bool{}
	for _, recipe := range view.Recipes {
		if recipe.ID == "" || recipe.Label == "" || recipe.OperatorStory == "" || recipe.Status == "" || recipe.WhatThisIs == "" || len(recipe.WhatYouNeed) == 0 || recipe.RunsWhere == "" || recipe.FirstSafeCheck == "" || recipe.GoodResult == "" || recipe.IfItFails == "" || recipe.DoesNotProve == "" || len(recipe.AdminLinks) == 0 || len(recipe.DocsLinks) == 0 || len(recipe.ManifestIDs) == 0 {
			t.Fatalf("invalid connector workbench recipe: %+v", recipe)
		}
		if seenRecipes[recipe.ID] {
			t.Fatalf("duplicate connector workbench recipe id %q", recipe.ID)
		}
		seenRecipes[recipe.ID] = true
		for _, link := range recipe.AdminLinks {
			if !strings.HasPrefix(link, "/admin/") {
				t.Fatalf("recipe %s has unsafe admin link %q", recipe.ID, link)
			}
		}
		for _, link := range recipe.DocsLinks {
			if !strings.HasPrefix(link, "docs/") && !strings.HasPrefix(link, "examples/") {
				t.Fatalf("recipe %s has unsafe docs link %q", recipe.ID, link)
			}
		}
	}
	seenTemplates := map[string]bool{}
	for _, template := range view.RedactionTemplates {
		if template.ID == "" || template.Label == "" || template.AppliesTo == "" || template.DataClassification == "" || len(template.AllowedFields) == 0 || len(template.RedactFields) == 0 || len(template.BlockedFields) == 0 || template.NoSendDefault == "" || template.FailClosedRule == "" || template.FirstSafeCheck == "" || template.DoesNotProve == "" || len(template.DocsLinks) == 0 {
			t.Fatalf("invalid connector workbench redaction template: %+v", template)
		}
		if seenTemplates[template.ID] {
			t.Fatalf("duplicate connector workbench redaction template %q", template.ID)
		}
		seenTemplates[template.ID] = true
		for _, link := range template.DocsLinks {
			if !strings.HasPrefix(link, "docs/") && !strings.HasPrefix(link, "examples/") {
				t.Fatalf("redaction template %s has unsafe docs link %q", template.ID, link)
			}
		}
	}
	seenRows := map[string]bool{}
	for _, row := range view.ManifestReview.Rows {
		if row.SourcePath == "" || row.ConnectorID == "" || row.DisplayName == "" || row.ConnectorType == "" || row.Mode == "" || row.SecretStorage == "" || len(row.InputContracts) == 0 || len(row.OutputContracts) == 0 || row.ConformanceCaseCount == 0 || row.DocsLink == "" || row.Boundary == "" || row.FirstCheck == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid connector workbench manifest row: %+v", row)
		}
		if seenRows[row.ConnectorID] {
			t.Fatalf("duplicate connector workbench manifest id %q", row.ConnectorID)
		}
		seenRows[row.ConnectorID] = true
		if !strings.HasPrefix(row.SourcePath, "examples/connectors/") || strings.Contains(row.SourcePath, "..") || strings.HasPrefix(row.SourcePath, "/") {
			t.Fatalf("manifest row %s has unsafe source path %q", row.ConnectorID, row.SourcePath)
		}
		if !strings.HasPrefix(row.DocsLink, "examples/connectors/") && !strings.HasPrefix(row.DocsLink, "docs/") {
			t.Fatalf("manifest row %s has unsafe docs link %q", row.ConnectorID, row.DocsLink)
		}
	}
	seenLintChecks := map[string]bool{}
	for _, lint := range view.ManifestReview.LintChecks {
		if lint.ID == "" || lint.Label == "" || lint.Status == "" || lint.EnforcedBy == "" || lint.Blocks == "" || lint.OperatorAction == "" || lint.DoesNotProve == "" {
			t.Fatalf("invalid connector workbench manifest lint: %+v", lint)
		}
		if seenLintChecks[lint.ID] {
			t.Fatalf("duplicate connector workbench manifest lint %q", lint.ID)
		}
		seenLintChecks[lint.ID] = true
	}
	if view.TelemetryPreview.Counts.Sources != 2 || view.TelemetryPreview.Counts.Rows != 6 || view.TelemetryPreview.Counts.Events != 4 || view.TelemetryPreview.Counts.Drops != 2 || view.TelemetryPreview.Counts.NetworkSendEnabled {
		t.Fatalf("invalid connector workbench telemetry preview counts: %+v", view.TelemetryPreview.Counts)
	}
	seenDryRuns := map[string]bool{}
	for _, command := range view.DryRunCommands {
		if command.ID == "" || command.Label == "" || command.CommandLine == "" || command.RunsWhere == "" || command.Inputs == "" || command.ExpectedResult == "" || command.FailureNextAction == "" || command.DoesNotProve == "" || len(command.DocsLinks) == 0 {
			t.Fatalf("invalid connector workbench dry-run command: %+v", command)
		}
		if seenDryRuns[command.ID] {
			t.Fatalf("duplicate connector workbench dry-run id %q", command.ID)
		}
		seenDryRuns[command.ID] = true
		for _, link := range command.DocsLinks {
			if !strings.HasPrefix(link, "docs/") && !strings.HasPrefix(link, "examples/") {
				t.Fatalf("dry-run %s has unsafe docs link %q", command.ID, link)
			}
		}
	}
	for _, source := range view.TelemetryPreview.Sources {
		if source.ID == "" || source.Label == "" || source.FixturePath == "" || source.Status == "" || !source.SyntheticOnly || source.ObservedRows != 3 || source.ExpectedEvents != 2 || source.ExpectedDrops != 1 || source.CommandLine == "" || source.DoesNotProve == "" {
			t.Fatalf("invalid connector workbench telemetry source: %+v", source)
		}
		if !strings.HasPrefix(source.FixturePath, "examples/connectors/") || strings.Contains(source.FixturePath, "..") || strings.HasPrefix(source.FixturePath, "/") {
			t.Fatalf("unsafe connector workbench fixture path: %+v", source)
		}
	}
	for _, row := range view.TelemetryPreview.Rows {
		if row.SourceID == "" || row.DeviceID == "" || row.VehicleID == "" || row.ObservedAt == "" || row.Quality == "" || row.Outcome == "" || !row.DryRun || row.NetworkSend {
			t.Fatalf("invalid connector workbench telemetry row: %+v", row)
		}
		if row.Outcome == "drop" && row.Reason == "" {
			t.Fatalf("drop row missing reason: %+v", row)
		}
	}
	seenWebhookRows := map[string]bool{}
	for _, row := range view.WebhookBoundary.Rows {
		if row.ID == "" || row.Label == "" || row.WhatThisMeans == "" || len(row.AllowedInputs) == 0 || len(row.BlockedInputs) == 0 || row.FirstSafeCheck == "" || row.FailClosedRule == "" || row.RedactionRule == "" || row.DoesNotProve == "" || len(row.ReviewLinks) == 0 {
			t.Fatalf("invalid connector workbench webhook row: %+v", row)
		}
		if seenWebhookRows[row.ID] {
			t.Fatalf("duplicate connector workbench webhook row %q", row.ID)
		}
		seenWebhookRows[row.ID] = true
		for _, link := range row.ReviewLinks {
			if !strings.HasPrefix(link, "/admin/") {
				t.Fatalf("webhook row %s has unsafe review link %q", row.ID, link)
			}
		}
	}
	assertConnectorWorkbenchGuideShape(t, "prediction", view.PredictionGuide)
	assertConnectorWorkbenchGuideShape(t, "monitoring", view.MonitoringGuide)
	assertConnectorWorkbenchGuideShape(t, "consumer", view.ConsumerGuide)
	assertConnectorWorkbenchConformanceShape(t, view.Conformance)
}

func assertConnectorWorkbenchGuideShape(t *testing.T, label string, guide connectorWorkbenchGuide) {
	t.Helper()
	seen := map[string]bool{}
	if guide.Title == "" || guide.Boundary == "" || len(guide.Rows) == 0 || len(guide.DocsLinks) == 0 {
		t.Fatalf("invalid %s guide: %+v", label, guide)
	}
	for _, row := range guide.Rows {
		if row.ID == "" || row.Label == "" || row.Status == "" || row.WhatThisIs == "" || len(row.Inputs) == 0 || len(row.Outputs) == 0 || row.FailureBehavior == "" || row.FirstSafeCheck == "" || row.DoesNotProve == "" || len(row.ReviewLinks) == 0 || len(row.DocsLinks) == 0 {
			t.Fatalf("invalid %s guide row: %+v", label, row)
		}
		if seen[row.ID] {
			t.Fatalf("duplicate %s guide row id %q", label, row.ID)
		}
		seen[row.ID] = true
		for _, link := range row.ReviewLinks {
			if !strings.HasPrefix(link, "/admin/") {
				t.Fatalf("%s guide row %s has unsafe review link %q", label, row.ID, link)
			}
		}
		for _, link := range row.DocsLinks {
			if !strings.HasPrefix(link, "docs/") && !strings.HasPrefix(link, "examples/") {
				t.Fatalf("%s guide row %s has unsafe docs link %q", label, row.ID, link)
			}
		}
	}
}

func assertConnectorWorkbenchConformanceShape(t *testing.T, view connectorWorkbenchConformanceView) {
	t.Helper()
	wantCases := map[string]int{"telemetry": 10, "prediction": 7, "validator": 2, "monitoring": 3, "consumer_discovery": 3}
	seen := map[string]bool{}
	if view.Boundary == "" || view.SuitePath == "" || view.Status == "" || !view.SyntheticOnly || view.ManifestCount != 13 || view.CaseCount != 25 || len(view.Groups) != 5 || len(view.RunnerCommands) != 4 {
		t.Fatalf("invalid connector workbench conformance view: %+v", view)
	}
	for _, command := range view.RunnerCommands {
		if command.ID == "" || command.Label == "" || command.CommandLine == "" || command.Inputs == "" || command.ExpectedResult == "" || command.FailureNextAction == "" || command.DoesNotProve == "" || len(command.DocsLinks) == 0 {
			t.Fatalf("invalid connector conformance runner command: %+v", command)
		}
	}
	for _, group := range view.Groups {
		if group.ID == "" || group.Label == "" || group.Status != "covered" || group.CaseCount != wantCases[group.ID] || len(group.RequiredScenarios) == 0 || group.CommandLine == "" || group.DoesNotProve == "" || len(group.Cases) != wantCases[group.ID] {
			t.Fatalf("invalid connector conformance group: %+v", group)
		}
		if seen[group.ID] {
			t.Fatalf("duplicate connector conformance group %q", group.ID)
		}
		seen[group.ID] = true
		for _, tc := range group.Cases {
			if tc.ID == "" || tc.Scenario == "" || tc.FixturePath == "" || tc.ExpectedOutcome == "" || len(tc.Assertions) == 0 || tc.Status != "covered" || !tc.SyntheticOnly {
				t.Fatalf("invalid connector conformance case: %+v", tc)
			}
			if !strings.HasPrefix(tc.FixturePath, "testdata/adapter-conformance/fixtures/") || strings.Contains(tc.FixturePath, "..") || strings.HasPrefix(tc.FixturePath, "/") {
				t.Fatalf("unsafe connector conformance fixture path: %+v", tc)
			}
		}
	}
}

func assertConnectorWorkbenchFlagsFalse(t *testing.T, flags connectorWorkbenchClaimFlags) {
	t.Helper()
	if flags.BackendCommandExecutionEnabled || flags.BrowserNetworkSendEnabled || flags.ManifestCommandExecutionEnabled || flags.DynamicBackendPluginLoadingEnabled || flags.ExternalNetworkContacted || flags.ExternalEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.VendorCompatibilityClaimed || flags.HardwareCertificationClaimed || flags.ProductionReadinessClaimed || flags.HostedSaaSClaimed || flags.SLAClaimed || flags.ProductionGradeETAClaimed {
		t.Fatalf("connector workbench flags must all be false: %+v", flags)
	}
}

func assertOperationsCockpitShape(t *testing.T, view operationsCockpitView) {
	t.Helper()
	if view.GeneratedAt.IsZero() || view.AgencyID == "" || view.Boundary == "" || len(view.ActionQueue) != 8 || len(view.SetupProgress) != 8 || len(view.PrimaryCards) != 11 {
		t.Fatalf("invalid cockpit shape: %+v", view)
	}
	seen := map[string]bool{}
	for _, action := range view.ActionQueue {
		if action.ID == "" || action.Label == "" || action.Status == "" || action.Signal == "" || action.ActionLabel == "" || action.AdminLink == "" || action.HelpNeeded == "" || action.DoesNotProve == "" {
			t.Fatalf("invalid cockpit action: %+v", action)
		}
		if seen["action:"+action.ID] {
			t.Fatalf("duplicate cockpit action id %q", action.ID)
		}
		seen["action:"+action.ID] = true
		if !strings.HasPrefix(action.AdminLink, "/admin/") {
			t.Fatalf("unsafe action admin link %q", action.AdminLink)
		}
	}
	for _, row := range view.SetupProgress {
		if row.ID == "" || row.Label == "" || row.Status == "" || row.CurrentSignal == "" || row.NextAction == "" || row.AdminLink == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid cockpit progress row: %+v", row)
		}
		if seen["progress:"+row.ID] {
			t.Fatalf("duplicate cockpit progress id %q", row.ID)
		}
		seen["progress:"+row.ID] = true
		if !strings.HasPrefix(row.AdminLink, "/admin/") {
			t.Fatalf("unsafe progress admin link %q", row.AdminLink)
		}
	}
	for _, card := range view.PrimaryCards {
		if card.ID == "" || card.Label == "" || card.Status == "" || card.CurrentSignal == "" || card.NextAction == "" || card.AdminLink == "" || card.DoesNotProve == "" {
			t.Fatalf("invalid cockpit card: %+v", card)
		}
		if seen["card:"+card.ID] {
			t.Fatalf("duplicate cockpit card id %q", card.ID)
		}
		seen["card:"+card.ID] = true
		if !strings.HasPrefix(card.AdminLink, "/admin/") {
			t.Fatalf("unsafe card admin link %q", card.AdminLink)
		}
		for _, link := range card.DocsLinks {
			if !strings.HasPrefix(link, "docs/") {
				t.Fatalf("unsafe card docs link %q", link)
			}
		}
	}
}

func assertOperationsCockpitFlagsFalse(t *testing.T, flags operationsCockpitClaimFlag) {
	t.Helper()
	if flags.ExternalEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.AgencyApprovalClaimed || flags.ConsumerAcceptanceClaimed || flags.PublicLaunchClaimed || flags.HostedSaaSClaimed || flags.VendorCompatibilityClaimed || flags.HardwareCertificationClaimed || flags.SLAClaimed || flags.UptimeGuaranteeClaimed || flags.ProductionGradeETAClaimed {
		t.Fatalf("cockpit flags must all be false: %+v", flags)
	}
}

func assertMaintenanceShape(t *testing.T, view operationsMaintenanceView) {
	t.Helper()
	if view.GeneratedAt.IsZero() || view.AgencyID == "" || view.Boundary == "" || view.OverallStatus == "" || len(view.SummaryRows) != 9 || view.Diagnostics.Boundary == "" || view.Diagnostics.Status == "" || len(view.Diagnostics.Rows) != 4 || view.SmallHostReadiness.Boundary == "" || view.SmallHostReadiness.Status == "" || view.SmallHostReadiness.NextAction == "" || len(view.SmallHostReadiness.Rows) != 5 || view.BackupRestore.Boundary == "" || view.BackupRestore.Status == "" || view.BackupRestore.NextAction == "" || len(view.BackupRestore.Rows) != 4 || view.UpgradeRollback.Boundary == "" || view.UpgradeRollback.Status == "" || view.UpgradeRollback.NextAction == "" || len(view.UpgradeRollback.Rows) != 4 || view.SupportReview.Boundary == "" || view.SupportReview.Status == "" || view.SupportReview.NextAction == "" || len(view.SupportReview.Rows) != 4 || view.CadencePlan.Boundary == "" || view.CadencePlan.Status == "" || view.CadencePlan.NextAction == "" || len(view.CadencePlan.Rows) != 4 || view.MonitoringExport.Boundary == "" || view.MonitoringExport.Status == "" || view.MonitoringExport.NextAction == "" || len(view.MonitoringExport.Rows) != 5 || view.Infrastructure.Boundary == "" || view.Infrastructure.Status == "" || view.Infrastructure.NextAction == "" || len(view.Infrastructure.Rows) != 10 || len(view.Tasks) != 7 {
		t.Fatalf("invalid maintenance shape: %+v", view)
	}
	seen := map[string]bool{}
	for _, row := range view.SummaryRows {
		if row.ID == "" || row.Label == "" || row.Status == "" || row.CurrentSignal == "" || row.NextAction == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid maintenance row: %+v", row)
		}
		if seen["row:"+row.ID] {
			t.Fatalf("duplicate maintenance row id %q", row.ID)
		}
		seen["row:"+row.ID] = true
	}
	for _, row := range view.Diagnostics.Rows {
		if row.ID == "" || row.Label == "" || row.Status == "" || row.SourceRef == "" || row.GeneratedAt == "" || row.CurrentSignal == "" || row.NextAction == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid maintenance diagnostic row: %+v", row)
		}
		if seen["diagnostic:"+row.ID] {
			t.Fatalf("duplicate maintenance diagnostic id %q", row.ID)
		}
		seen["diagnostic:"+row.ID] = true
	}
	for _, panel := range []struct {
		name string
		rows []operationsMaintenancePanelRow
	}{
		{name: "small-host", rows: view.SmallHostReadiness.Rows},
		{name: "backup", rows: view.BackupRestore.Rows},
		{name: "upgrade", rows: view.UpgradeRollback.Rows},
		{name: "support", rows: view.SupportReview.Rows},
		{name: "cadence", rows: view.CadencePlan.Rows},
		{name: "monitoring", rows: view.MonitoringExport.Rows},
		{name: "infrastructure", rows: view.Infrastructure.Rows},
	} {
		for _, row := range panel.rows {
			if row.ID == "" || row.Label == "" || row.Status == "" || row.CurrentSignal == "" || row.OperatorStep == "" || row.AdministratorStep == "" || row.DoesNotProve == "" {
				t.Fatalf("invalid maintenance panel row: %+v", row)
			}
			if seen[panel.name+":"+row.ID] {
				t.Fatalf("duplicate maintenance panel row id %q", row.ID)
			}
			seen[panel.name+":"+row.ID] = true
		}
	}
	for _, task := range view.Tasks {
		if task.ID == "" || task.Cadence == "" || task.Task == "" || task.Status == "" || task.Owner == "" || task.NextStep == "" {
			t.Fatalf("invalid maintenance task: %+v", task)
		}
		if seen["task:"+task.ID] {
			t.Fatalf("duplicate maintenance task id %q", task.ID)
		}
		seen["task:"+task.ID] = true
	}
	if view.SupportSummary.Status == "" || view.SupportSummary.Command == "" || view.SupportSummary.OutputPath == "" || len(view.SupportSummary.Instructions) == 0 {
		t.Fatalf("invalid maintenance support summary: %+v", view.SupportSummary)
	}
}

func assertMaintenanceFlagsFalse(t *testing.T, flags operationsMaintenanceClaimFlags) {
	t.Helper()
	if flags.ExternalEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.SLAClaimed || flags.UptimeGuaranteeClaimed || flags.HostedSaaSClaimed || flags.AgencyAdoptionClaimed || flags.ConsumerAcceptanceClaimed || flags.VendorCompatibilityClaimed || flags.ProductionGradeETAClaimed {
		t.Fatalf("maintenance flags must all be false: %+v", flags)
	}
}

func assertFeedHealthShape(t *testing.T, health operationsFeedHealthView) {
	t.Helper()
	if health.GeneratedAt.IsZero() || health.AgencyID == "" || health.Boundary == "" || len(health.Rows) != 5 || health.Counts.Rows != 5 {
		t.Fatalf("invalid feed health shape: %+v", health)
	}
	for _, row := range health.Rows {
		if row.ID == "" || row.Label == "" || row.PublicPath == "" || row.ConfiguredURL == "" || row.LastKnownHTTPStatus == "" || row.ByteCount == "" || row.ContentType == "" || row.Checksum == "" || row.LastGenerated == "" || row.LastChecked == "" || row.ValidatorState == "" || row.HealthState == "" || row.Status == "" || row.StatusText == "" || row.CurrentSignal == "" || row.WhatThisMeans == "" || row.Freshness == "" || row.ValidatorContext == "" || row.HealthContext == "" || row.NextAction == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid feed health row: %+v", row)
		}
		for _, link := range row.AdminLinks {
			if !strings.HasPrefix(link, "/admin/") {
				t.Fatalf("row %s has unsafe admin link %q", row.ID, link)
			}
		}
		for _, link := range row.DocsLinks {
			if !strings.HasPrefix(link, "docs/") {
				t.Fatalf("row %s has unsafe docs link %q", row.ID, link)
			}
		}
	}
	if health.RealtimeUsefulness.VehiclePositions.ID != "vehicle_positions" || health.RealtimeUsefulness.TripUpdates.ID != "trip_updates" || health.RealtimeUsefulness.Alerts.ID != "alerts" {
		t.Fatalf("invalid realtime usefulness shape: %+v", health.RealtimeUsefulness)
	}
}

func assertFeedHealthFlagsFalse(t *testing.T, flags operationsFeedHealthClaims) {
	t.Helper()
	if flags.ExternalEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.SLAClaimed || flags.UptimeGuaranteeClaimed || flags.ConsumerAcceptanceClaimed || flags.PublicLaunchClaimed {
		t.Fatalf("feed health flags must all be false: %+v", flags)
	}
}

func assertRealtimeFlagsFalse(t *testing.T, flags operationsRealtimeClaimFlags) {
	t.Helper()
	if flags.BrowserTelemetrySendEnabled || flags.BackendCommandExecutionEnabled || flags.DeviceTokenCollectedByBrowser || flags.ExternalEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.VendorCompatibilityClaimed || flags.HardwareCertificationClaimed || flags.ProductionAVLReliabilityClaimed || flags.ProductionGradeETAClaimed || flags.RealWorldETAAccuracyClaimed || flags.SLAClaimed || flags.PublicLaunchClaimed || flags.ConsumerAcceptanceClaimed {
		t.Fatalf("realtime flags must all be false: %+v", flags)
	}
}

func assertPredictionLabShape(t *testing.T, view predictionLabView) {
	t.Helper()
	if view.GeneratedAt.IsZero() || view.AgencyID == "" || view.Boundary == "" || view.Summary.CurrentSignal == "" || view.Summary.NextAction == "" || view.Summary.DoesNotProve == "" || view.Deterministic.Boundary == "" || view.Deterministic.Status == "" || view.Deterministic.ReviewSignal == "" || len(view.Deterministic.Rows) != 4 || len(view.WithheldReasons) == 0 || view.ShadowReview.Boundary == "" || view.ShadowReview.Status == "" || view.ShadowReview.NextAction == "" || view.ShadowReview.DoesNotProve == "" || len(view.ShadowReview.Rows) == 0 || view.Backtests.Boundary == "" || view.Backtests.Status == "" || view.Backtests.RootRef == "" || view.Backtests.Message == "" || view.HandlingGuide.Boundary == "" || view.HandlingGuide.Status == "" || view.HandlingGuide.NextAction == "" || view.HandlingGuide.DoesNotProve == "" || len(view.HandlingGuide.Rows) != 4 || view.ProofChecklist.Boundary == "" || view.ProofChecklist.Status == "" || view.ProofChecklist.NextAction == "" || view.ProofChecklist.DoesNotProve == "" || len(view.ProofChecklist.Rows) != 4 || len(view.ReviewRows) == 0 || len(view.Commands) != 3 {
		t.Fatalf("invalid prediction lab shape: %+v", view)
	}
	for _, row := range view.Deterministic.Rows {
		if row.ID == "" || row.Label == "" || row.Status == "" || row.CurrentSignal == "" || row.NextAction == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid prediction lab diagnostic row: %+v", row)
		}
	}
	for _, row := range view.WithheldReasons {
		if row.Reason == "" || row.Label == "" || row.WhatItMeans == "" || row.NextAction == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid prediction lab reason row: %+v", row)
		}
	}
	for _, row := range view.ShadowReview.Rows {
		if row.ID == "" || row.Label == "" || row.Status == "" || row.Reason == "" || row.Latency == "" || row.CountComparison == "" || row.FailureBehavior == "" || row.FirstSafeCheck == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid prediction lab shadow row: %+v", row)
		}
	}
	for _, row := range view.Backtests.Rows {
		if row.OutputRef == "" || row.Status == "" || row.GeneratedAt == "" || row.MaturityGate == "" || row.PredictionCoverage == "" || row.FutureStopCoverage == "" || row.MAEAbsoluteErrorSeconds == "" || row.P90AbsoluteErrorSeconds == "" || row.WithheldByReason == "" || row.ConformanceSignal == "" || row.DiagnosticSignal == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid prediction lab backtest row: %+v", row)
		}
	}
	for _, row := range view.HandlingGuide.Rows {
		if row.ID == "" || row.Situation == "" || row.ReviewSignal == "" || row.SafeBehavior == "" || row.OperatorStep == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid prediction lab handling row: %+v", row)
		}
	}
	for _, row := range view.ProofChecklist.Rows {
		if row.ID == "" || row.FutureGate == "" || row.RequiredReview == "" || row.SeparateAuthorization == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid prediction lab proof checklist row: %+v", row)
		}
	}
	for _, row := range view.ReviewRows {
		if row.Severity == "" || row.Area == "" || row.Signal == "" || row.NextAction == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid prediction lab review row: %+v", row)
		}
	}
	for _, command := range view.Commands {
		if command.ID == "" || command.Label == "" || command.CommandLine == "" || command.ExpectedResult == "" || command.DoesNotProve == "" {
			t.Fatalf("invalid prediction lab command: %+v", command)
		}
		if strings.Contains(command.CommandLine, "curl ") || strings.Contains(command.CommandLine, "http://") || strings.Contains(command.CommandLine, "https://") {
			t.Fatalf("prediction lab command must remain local/offline guidance: %+v", command)
		}
	}
}

func assertPredictionLabFlagsFalse(t *testing.T, flags predictionLabClaimFlags) {
	t.Helper()
	if flags.BrowserPredictorRunEnabled || flags.ExternalNetworkContacted || flags.BackendCommandExecutionEnabled || flags.ExternalEvidenceCreated || flags.FinalRootEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.ProductionGradeETAClaimed || flags.RealWorldETAAccuracyClaimed || flags.VendorCompatibilityClaimed || flags.HardwareCertificationClaimed || flags.SLAClaimed || flags.HostedSaaSClaimed || flags.PublicLaunchClaimed || flags.ConsumerAcceptanceClaimed || flags.RawObservedRowsPersisted {
		t.Fatalf("prediction lab flags must all be false: %+v", flags)
	}
}

func assertPredictionLabSafeStrings(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"raw-token-value", "authorization:", "set-cookie", "database_url", "restore_database_url", "payload_json", "raw telemetry payload", "raw observed", "raw prediction", "raw gtfs-rt", "token_hash", "file://", "/users/", "/opt/open-transit-rt", "/var/lib", "/etc/", "postgres://", "raw_report", "stdout", "stderr", "argv", "private_debug", "score_details", "bearer ", "predictor.example", "raw_response", "token=secret"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("prediction lab leaks forbidden private string %q: %s", forbidden, body)
		}
	}
	for _, forbidden := range []string{"agency_approved", "final_root_approved", "consumer_ready", "production_ready", "public_launch_complete", "compliance_achieved", "sla_covered", "uptime_guaranteed", "eta_accuracy_proven", "production_grade_eta_proven"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("prediction lab emits forbidden label %q: %s", forbidden, body)
		}
	}
}

func assertFeedHealthSafeStrings(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"raw-token-value", "authorization:", "set-cookie", ".cache", "database_url", "restore_database_url", "payload_json", "raw telemetry", "token_hash", "file://", "/users/", "/opt/open-transit-rt", "/var/lib", "/etc/", "postgres://", "raw_report", "stdout", "stderr", "argv"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("feed health leaks forbidden private string %q: %s", forbidden, body)
		}
	}
	for _, forbidden := range []string{"agency_approved", "final_root_approved", "consumer_ready", "production_ready", "public_launch_complete", "compliance_achieved", "sla_covered", "uptime_guaranteed"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("feed health emits forbidden label %q: %s", forbidden, body)
		}
	}
}

func assertValidationCenterShape(t *testing.T, center operationsValidationCenterView) {
	t.Helper()
	if center.GeneratedAt.IsZero() || center.AgencyID == "" || center.Boundary == "" || len(center.FeedRows) != 5 || len(center.ValidationHistory) != 4 || len(center.ValidatorHealth) != 4 || len(center.GTFSQuality) != 2 || len(center.ReadinessTimeline) != 11 || len(center.ConsumerTracker) != 7 {
		t.Fatalf("invalid validation center shape: %+v", center)
	}
	if center.Counts.FeedRows != len(center.FeedRows) || center.Counts.ValidationRows != len(center.ValidationHistory) || center.Counts.IssueRows != len(center.IssueDrilldowns) || center.Counts.ConsumerRows != len(center.ConsumerTracker) || len(center.Counts.Statuses) == 0 {
		t.Fatalf("invalid validation center counts: %+v", center.Counts)
	}
	for _, row := range center.FeedRows {
		if row.ID == "" || row.Label == "" || row.PublicPath == "" || row.ConfiguredURL == "" || row.Status == "" || row.LastChecked == "" || row.HTTPStatus == "" || row.ContentType == "" || row.Freshness == "" || row.ValidatorState == "" || row.HealthState == "" || row.CurrentSignal == "" || row.WhatThisMeans == "" || row.NextAction == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid validation center feed row: %+v", row)
		}
	}
	for _, row := range center.ValidationHistory {
		if row.ID == "" || row.FeedType == "" || row.Label == "" || row.ValidatorID == "" || row.ValidatorName == "" || row.Status == "" || row.ToolingStatus == "" || row.ArtifactStatus == "" || row.LatestResultStatus == "" || row.StaleStatus == "" || row.HealthStatus == "" || row.CurrentSignal == "" || row.WhatThisMeans == "" || row.NextAction == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid validation center validator row: %+v", row)
		}
	}
	for _, row := range center.GTFSQuality {
		if row.ID == "" || row.Label == "" || row.Status == "" || row.CurrentSignal == "" || row.WhatThisMeans == "" || row.NextAction == "" || row.DoesNotProve == "" || row.DetailsURL != "/admin/operations/gtfs-quality" {
			t.Fatalf("invalid validation center GTFS quality row: %+v", row)
		}
	}
	for _, row := range center.IssueDrilldowns {
		if row.ID == "" || row.Source == "" || row.SourceLabel == "" || row.Status == "" || row.Severity == "" || row.Family == "" || len(row.Codes) == 0 || row.LikelyOwner == "" || row.AffectedFiles == "" || row.OperatorSummary == "" || row.WhyItMatters == "" || row.RecommendedAction == "" || row.SafeFixPath == "" || row.VerifyWith == "" || row.EscalateIf == "" || row.DetailsURL != "/admin/operations/gtfs-quality" || row.DoesNotProve == "" {
			t.Fatalf("invalid validation center issue row: %+v", row)
		}
	}
	for _, row := range center.ReadinessTimeline {
		if row.ID == "" || row.Label == "" || row.Status == "" || row.CurrentSignal == "" || row.WhatThisMeans == "" || row.NextAction == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid validation center timeline row: %+v", row)
		}
	}
	for _, row := range center.Blockers {
		if row.ID == "" || row.Severity == "" || row.Area == "" || row.Signal == "" || row.NextAction == "" || row.DoesNotProve == "" || row.ReviewURL == "" {
			t.Fatalf("invalid validation center blocker row: %+v", row)
		}
		if !strings.HasPrefix(row.ReviewURL, "/admin/") {
			t.Fatalf("unsafe validation center blocker URL: %+v", row)
		}
	}
	for _, row := range center.ConsumerTracker {
		if row.Target == "" || row.Status != "prepared" || row.Source == "" || row.NextAction == "" || row.DoesNotProve == "" {
			t.Fatalf("invalid validation center consumer row: %+v", row)
		}
	}
}

func assertValidationCenterFlagsFalse(t *testing.T, flags operationsValidationCenterClaimFlags) {
	t.Helper()
	if flags.ExternalEvidenceCreated || flags.FinalRootEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.AgencyApprovalClaimed || flags.ConsumerAcceptanceClaimed || flags.PublicLaunchClaimed || flags.HostedSaaSClaimed || flags.SLAClaimed || flags.UptimeGuaranteeClaimed || flags.VendorCompatibilityClaimed || flags.HardwareCertificationClaimed || flags.ProductionGradeETAClaimed {
		t.Fatalf("validation center claim flags must all be false: %+v", flags)
	}
}

func assertValidationCenterJSONAllowlist(t *testing.T, payload []byte) {
	t.Helper()
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatal(err)
	}
	wantTop := map[string]bool{"generated_at": true, "agency_id": true, "boundary": true, "feed_rows": true, "validation_history": true, "validator_health": true, "gtfs_quality": true, "issue_drilldowns": true, "readiness_timeline": true, "blockers": true, "consumer_tracker": true, "counts": true, "claim_flags": true}
	for key := range decoded {
		if !wantTop[key] {
			t.Fatalf("unexpected validation center top-level field %q in %s", key, payload)
		}
	}
	wantFeed := map[string]bool{"id": true, "label": true, "public_path": true, "configured_url": true, "status": true, "last_checked": true, "http_status": true, "content_type": true, "freshness": true, "validator_state": true, "health_state": true, "current_signal": true, "what_this_means": true, "next_action": true, "does_not_prove": true}
	for _, item := range decoded["feed_rows"].([]any) {
		for key := range item.(map[string]any) {
			if !wantFeed[key] {
				t.Fatalf("unexpected validation center feed row field %q in %s", key, payload)
			}
		}
	}
	wantValidation := map[string]bool{"id": true, "feed_type": true, "label": true, "validator_id": true, "validator_name": true, "status": true, "tooling_status": true, "artifact_status": true, "latest_result_status": true, "latest_result_at": true, "active_feed_version_id": true, "latest_result_feed_version_id": true, "stale_status": true, "health_status": true, "current_signal": true, "what_this_means": true, "next_action": true, "does_not_prove": true}
	for _, item := range decoded["validation_history"].([]any) {
		for key := range item.(map[string]any) {
			if !wantValidation[key] {
				t.Fatalf("unexpected validation center validation row field %q in %s", key, payload)
			}
		}
	}
	wantIssue := map[string]bool{"id": true, "source": true, "source_label": true, "status": true, "severity": true, "family": true, "codes": true, "count": true, "sample_count": true, "overflow_count": true, "likely_owner": true, "risk_level": true, "affected_files": true, "operator_summary": true, "why_it_matters": true, "recommended_action": true, "safe_fix_path": true, "verify_with": true, "escalate_if": true, "details_url": true, "does_not_prove": true}
	for _, item := range decoded["issue_drilldowns"].([]any) {
		for key := range item.(map[string]any) {
			if !wantIssue[key] {
				t.Fatalf("unexpected validation center issue field %q in %s", key, payload)
			}
		}
	}
	wantTimeline := map[string]bool{"id": true, "label": true, "status": true, "current_signal": true, "what_this_means": true, "next_action": true, "does_not_prove": true}
	for _, item := range decoded["readiness_timeline"].([]any) {
		for key := range item.(map[string]any) {
			if !wantTimeline[key] {
				t.Fatalf("unexpected validation center timeline field %q in %s", key, payload)
			}
		}
	}
	wantBlocker := map[string]bool{"id": true, "severity": true, "area": true, "signal": true, "next_action": true, "does_not_prove": true, "review_url": true}
	for _, item := range decoded["blockers"].([]any) {
		for key := range item.(map[string]any) {
			if !wantBlocker[key] {
				t.Fatalf("unexpected validation center blocker field %q in %s", key, payload)
			}
		}
	}
	wantFlags := map[string]bool{"external_evidence_created": true, "final_root_evidence_created": true, "consumer_statuses_changed": true, "compliance_claimed": true, "production_readiness_claimed": true, "agency_approval_claimed": true, "consumer_acceptance_claimed": true, "public_launch_claimed": true, "hosted_saas_claimed": true, "sla_claimed": true, "uptime_guarantee_claimed": true, "vendor_compatibility_claimed": true, "hardware_certification_claimed": true, "production_grade_eta_claimed": true}
	flags := decoded["claim_flags"].(map[string]any)
	for key := range flags {
		if !wantFlags[key] {
			t.Fatalf("unexpected validation center claim flag %q in %s", key, payload)
		}
	}
}

func assertValidationCenterSafeStrings(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"raw_report", "stdout", "stderr", "argv", "/users/private", "/tmp/private", "token=secret", "password=secret", "postgres://", "authorization:", "bearer ", "cookie", "admin_session", "database_url"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("validation center leaked forbidden private string %q: %s", forbidden, body)
		}
	}
	for _, forbidden := range []string{"agency_approved", "final_root_approved", "consumer_ready", "production_ready", "public_launch_complete", "compliance_achieved", "sla_covered", "uptime_guaranteed", "vendor_certified", "hardware_certified"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("validation center emits forbidden label %q: %s", forbidden, body)
		}
	}
}

func assertReadinessV2Shape(t *testing.T, readiness operationsReadinessV2View) {
	t.Helper()
	if readiness.GeneratedAt.IsZero() || readiness.AgencyID == "" || readiness.Boundary == "" || len(readiness.FocusAreas) != 10 || len(readiness.Rows) != 11 || readiness.Counts.Rows != 11 {
		t.Fatalf("invalid readiness v2 shape: %+v", readiness)
	}
	wantFocusIDs := []string{"public_feed_urls", "static_gtfs", "vehicle_positions", "trip_updates", "alerts", "validation", "license_contact", "uptime_operations", "telemetry_device_state", "consumer_preparedness"}
	var gotFocusIDs []string
	for _, focus := range readiness.FocusAreas {
		gotFocusIDs = append(gotFocusIDs, focus.ID)
		if focus.ID == "" || focus.Label == "" || focus.Status == "" || focus.WhatThisHelpsWith == "" || focus.PrimarySignal == "" || focus.NextAction == "" || focus.WhatItDoesNotProve == "" || len(focus.RowIDs) == 0 || len(focus.AdminLinks) == 0 || len(focus.DocsLinks) == 0 {
			t.Fatalf("invalid readiness focus area: %+v", focus)
		}
		for _, link := range focus.AdminLinks {
			if !strings.HasPrefix(link, "/admin/") {
				t.Fatalf("focus %s has unsafe admin link %q", focus.ID, link)
			}
		}
		for _, link := range focus.DocsLinks {
			if !strings.HasPrefix(link, "docs/") {
				t.Fatalf("focus %s has unsafe docs link %q", focus.ID, link)
			}
			if _, err := os.Stat(filepath.Join("..", "..", link)); err != nil {
				t.Fatalf("focus %s docs link %q should exist: %v", focus.ID, link, err)
			}
		}
	}
	if strings.Join(gotFocusIDs, ",") != strings.Join(wantFocusIDs, ",") {
		t.Fatalf("focus ids = %v, want %v", gotFocusIDs, wantFocusIDs)
	}
	wantIDs := []string{"discovery_metadata", "feed_health", "static_gtfs_quality", "vehicle_positions", "trip_updates", "alerts", "validation_health", "operations_reliability", "telemetry_devices", "operations_scorecard", "consumer_prepared_tracker"}
	var gotIDs []string
	for _, row := range readiness.Rows {
		gotIDs = append(gotIDs, row.ID)
		if row.ID == "" || row.ReadinessItem == "" || row.Status == "" || row.CurrentSignal == "" || row.WhatThisMeans == "" || row.WhyItMatters == "" || row.WhatToDoNext == "" || row.WhatItDoesNotProve == "" {
			t.Fatalf("invalid readiness v2 row: %+v", row)
		}
		for _, link := range row.AdminLinks {
			if !strings.HasPrefix(link, "/admin/") {
				t.Fatalf("row %s has unsafe admin link %q", row.ID, link)
			}
		}
		for _, link := range row.DocsLinks {
			if !strings.HasPrefix(link, "docs/") {
				t.Fatalf("row %s has unsafe docs link %q", row.ID, link)
			}
		}
	}
	if strings.Join(gotIDs, ",") != strings.Join(wantIDs, ",") {
		t.Fatalf("row ids = %v, want %v", gotIDs, wantIDs)
	}
}

func assertReadinessV2FlagsFalse(t *testing.T, flags operationsReadinessV2Claims) {
	t.Helper()
	if flags.ExternalEvidenceCreated || flags.FinalRootEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.AgencyApprovalClaimed || flags.ConsumerAcceptanceClaimed || flags.PublicLaunchClaimed || flags.HostedSaaSClaimed || flags.SLAClaimed || flags.UptimeGuaranteeClaimed || flags.VendorCompatibilityClaimed || flags.ProductionGradeETAClaimed {
		t.Fatalf("readiness v2 flags must all be false: %+v", flags)
	}
}

func assertReadinessV2SafeStrings(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"raw-token-value", "authorization:", "set-cookie", ".cache", "database_url", "restore_database_url", "payload_json", "raw telemetry", "token_hash", "file://", "/users/", "/opt/open-transit-rt", "/var/lib", "/etc/", "postgres://", "raw_report", "stdout", "stderr", "argv"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("readiness v2 leaks forbidden private string %q: %s", forbidden, body)
		}
	}
	for _, forbidden := range []string{"agency_approved", "final_root_approved", "consumer_ready", "production_ready", "public_launch_complete", "compliance_achieved", "sla_covered", "uptime_guaranteed"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("readiness v2 emits forbidden label %q: %s", forbidden, body)
		}
	}
}

func assertConnectorHubSafeStrings(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"raw-token-value", "authorization:", "set-cookie", ".cache", "database_url", "restore_database_url", "payload_json", "raw telemetry", "token_hash", "file://", "/users/", "/opt/open-transit-rt", "/var/lib", "/etc/", "postgres://"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("connector hub leaks forbidden private string %q: %s", forbidden, body)
		}
	}
	for _, forbidden := range []string{"agency_approved", "final_root_approved", "consumer_ready", "production_ready", "public_launch_complete", "vendor_compatible", "hardware_certified", "dynamic_plugin_loading_enabled"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("connector hub emits forbidden label %q: %s", forbidden, body)
		}
	}
}

func assertConnectorWorkbenchSafeStrings(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"raw-token-value", "authorization:", "set-cookie", ".cache", "database_url", "restore_database_url", "payload_json", "token_hash", "file://", "/users/", "/opt/open-transit-rt", "/var/lib", "/etc/", "postgres://", "raw_validator_command", "raw_command", "client-supplied shell", "http://localhost", "127.0.0.1", "192.168.", "10.0.0.", ".local"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("connector workbench leaks forbidden private string %q: %s", forbidden, body)
		}
	}
	for _, forbidden := range []string{"agency_approved", "final_root_approved", "consumer_ready", "production_ready", "public_launch_complete", "vendor_compatible", "hardware_certified", "dynamic_plugin_loading_enabled"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("connector workbench emits forbidden label %q: %s", forbidden, body)
		}
	}
}

func assertConnectorTestsSafeStrings(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"raw-token-value", "authorization:", "set-cookie", ".cache", "database_url", "restore_database_url", "payload_json", "raw telemetry", "token_hash", "file://", "/users/", "/opt/open-transit-rt", "/var/lib", "/etc/", "postgres://", "raw_validator_command", "raw_command", "shell"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("connector tests leak forbidden private string %q: %s", forbidden, body)
		}
	}
	for _, forbidden := range []string{"agency_approved", "final_root_approved", "consumer_ready", "production_ready", "public_launch_complete", "vendor_compatible", "hardware_certified", "dynamic_plugin_loading_enabled"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("connector tests emit forbidden label %q: %s", forbidden, body)
		}
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

func assertSetupWizardSafeStrings(t *testing.T, body string) {
	t.Helper()
	lower := strings.ToLower(body)
	for _, forbidden := range []string{"raw-token-value", "authorization:", "set-cookie", ".cache", "database_url", "restore_database_url", "payload_json", "raw telemetry", "token_hash", "file://", "/users/", "/opt/open-transit-rt", "/var/lib", "/etc/", "postgres://"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("setup wizard leaks forbidden private string %q: %s", forbidden, body)
		}
	}
	for _, forbidden := range []string{"agency_approved", "final_root_approved", "consumer_ready", "production_ready", "public_launch_complete", "compliance_achieved"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("setup wizard emits forbidden label %q: %s", forbidden, body)
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

func firstRunTaskStatus(firstRun operationsFirstRunView, id string) string {
	for _, task := range firstRun.Tasks {
		if task.ID == id {
			return task.Status
		}
	}
	return ""
}

func setupWizardStageStatus(wizard operationsSetupWizardView, id string) string {
	for _, stage := range wizard.Stages {
		if stage.ID == id {
			return stage.Status
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

func feedHealthTestStore(t testing.TB) *fakePublicationStore {
	t.Helper()
	now := time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC)
	discovery := validationHealthTestDiscovery(now)
	discovery.GeneratedAt = now
	discovery.PublicBaseURL = "https://feeds.example.org"
	discovery.TechnicalContactEmail = "ops@example.org"
	discovery.License = compliance.License{Name: "CC BY 4.0", URL: "https://example.org/license"}
	discovery.Readiness = compliance.Readiness{AllRequiredFeedsListed: true, LicenseComplete: true, ContactComplete: true, HTTPSURLs: true, Discoverable: true}
	for i := range discovery.Feeds {
		discovery.Feeds[i].LastValidationStatus = "passed"
		discovery.Feeds[i].LastValidationAt = &now
		discovery.Feeds[i].LastHealthStatus = "ok"
		discovery.Feeds[i].LastHealthAt = &now
	}
	endpointAvailable := true
	freshness := 20.0
	latency := 120.0
	invalid := 0.0
	matched := 96.0
	coverage := 98.0
	records := []compliance.ReliabilityFeedHealthRecord{}
	for _, feedType := range []string{"schedule", "vehicle_positions", "trip_updates", "alerts"} {
		records = append(records, compliance.ReliabilityFeedHealthRecord{
			FeedType:               feedType,
			SnapshotAt:             now,
			EndpointAvailable:      &endpointAvailable,
			FreshnessSeconds:       &freshness,
			GenerationLatencyMS:    &latency,
			InvalidResponsePercent: &invalid,
			MatchedVehiclePercent:  &matched,
			CoveragePercent:        &coverage,
		})
	}
	return &fakePublicationStore{
		discovery: discovery,
		validationRecords: []compliance.ValidationReportRecord{
			validationHealthRecord(1, "schedule", "feed-v1", "passed", now),
			validationHealthRecord(2, "vehicle_positions", "feed-v1", "passed", now),
			validationHealthRecord(3, "trip_updates", "feed-v1", "passed", now),
			validationHealthRecord(4, "alerts", "feed-v1", "passed", now),
		},
		reliabilityHealth:    records,
		reliabilityIncidents: compliance.NormalizeReliabilityIncidentRollup(now, 0, nil, nil, nil, nil, nil, 10),
		tripDiagnostics:      compliance.TripUpdatesDiagnosticsSummary{Recorded: true, SnapshotAt: now, AdapterName: "deterministic", DiagnosticsStatus: "recorded", DiagnosticsReason: "test diagnostics", ActiveFeedVersionID: "feed-v1"},
	}
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

func assertAdminCommandResultJSONAllowlist(t *testing.T, payload []byte) {
	t.Helper()
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatal(err)
	}
	wantTop := map[string]bool{"action": true, "status": true, "started_at": true, "completed_at": true, "summary": true, "next_actions": true, "claim_flags": true, "errors": true}
	for key := range decoded {
		if !wantTop[key] {
			t.Fatalf("unexpected command result field %q in %s", key, payload)
		}
	}
	flags, ok := decoded["claim_flags"].(map[string]any)
	if !ok {
		t.Fatalf("claim_flags missing or not object in %s", payload)
	}
	wantFlags := map[string]bool{"external_evidence_created": true, "consumer_statuses_changed": true, "compliance_claimed": true, "production_readiness_claimed": true, "agency_approval_claimed": true, "consumer_acceptance_claimed": true, "public_launch_claimed": true, "hosted_saas_claimed": true, "vendor_compatibility_claimed": true, "hardware_certification_claimed": true, "sla_claimed": true, "uptime_guarantee_claimed": true, "production_grade_eta_claimed": true}
	for key := range flags {
		if !wantFlags[key] {
			t.Fatalf("unexpected command claim flag %q in %s", key, payload)
		}
	}
}

func assertAdminCommandFlagsFalse(t *testing.T, flags admincontrol.ClaimFlags) {
	t.Helper()
	if flags.ExternalEvidenceCreated || flags.ConsumerStatusesChanged || flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.AgencyApprovalClaimed || flags.ConsumerAcceptanceClaimed || flags.PublicLaunchClaimed || flags.HostedSaaSClaimed || flags.VendorCompatibilityClaimed || flags.HardwareCertificationClaimed || flags.SLAClaimed || flags.UptimeGuaranteeClaimed || flags.ProductionGradeETAClaimed {
		t.Fatalf("command claim flags must all be false: %+v", flags)
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

func (f *countingScheduleBuilder) SnapshotKey(context.Context) (schedule.SnapshotKey, error) {
	if f.err != nil {
		return schedule.SnapshotKey{}, f.err
	}
	return schedule.SnapshotKey{
		AgencyID:      f.snapshot.AgencyID,
		FeedVersionID: f.snapshot.FeedVersionID,
		RevisionTime:  f.snapshot.RevisionTime,
	}, nil
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

func (f fakeScheduleBuilder) SnapshotKey(context.Context) (schedule.SnapshotKey, error) {
	if f.err != nil {
		return schedule.SnapshotKey{}, f.err
	}
	return schedule.SnapshotKey{
		AgencyID:      f.snapshot.AgencyID,
		FeedVersionID: f.snapshot.FeedVersionID,
		RevisionTime:  f.snapshot.RevisionTime,
	}, nil
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

func (f fakeScheduleBuilder) SnapshotKeyForAgency(_ context.Context, agencyID string) (schedule.SnapshotKey, error) {
	if f.err != nil {
		return schedule.SnapshotKey{}, f.err
	}
	snapshot := f.snapshot
	if f.snapshotsByAgency != nil {
		snapshot = f.snapshotsByAgency[agencyID]
	} else {
		snapshot.AgencyID = agencyID
	}
	return schedule.SnapshotKey{
		AgencyID:      snapshot.AgencyID,
		FeedVersionID: snapshot.FeedVersionID,
		RevisionTime:  snapshot.RevisionTime,
	}, nil
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
	auditRows               []compliance.AuditLogRecord
	auditErr                error
	auditAgencyID           string
	gtfsImports             []compliance.GTFSImportRecord
	gtfsImportsErr          error
	gtfsPreview             compliance.GTFSSchedulePreview
	gtfsPreviews            map[string]compliance.GTFSSchedulePreview
	gtfsPreviewErr          error
	gtfsDrafts              []compliance.GTFSDraftRecord
	gtfsDraftsErr           error
	gtfsDraftPublishes      []compliance.GTFSDraftPublishRecord
	gtfsDraftPublishesErr   error
	feedVersions            []compliance.FeedVersionRecord
	feedVersionsErr         error
	reliabilityHealth       []compliance.ReliabilityFeedHealthRecord
	reliabilityIncidents    compliance.ReliabilityIncidentRollup
	reliabilityHealthErr    error
	reliabilityIncidentsErr error
}

func (f *fakePublicationStore) BootstrapPublication(_ context.Context, input compliance.BootstrapInput) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.bootstrapInput = input
	if f.bootstrapErr != nil {
		return f.bootstrapErr
	}
	return nil
}

func (f *fakePublicationStore) PublicationConfig(context.Context, string) (compliance.PublicationConfig, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.publicationConfigErr != nil {
		return compliance.PublicationConfig{}, f.publicationConfigErr
	}
	return f.publicationConfig, nil
}

func (f *fakePublicationStore) FeedDiscovery(_ context.Context, agencyID string, _ time.Time) (compliance.FeedDiscovery, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
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
	f.mu.Lock()
	defer f.mu.Unlock()
	f.listConsumersAgencyID = agencyID
	if f.consumersErr != nil {
		return nil, f.consumersErr
	}
	return append([]compliance.ConsumerRecord(nil), f.consumers...), nil
}

func (f *fakePublicationStore) LatestScorecard(_ context.Context, agencyID string) (compliance.Scorecard, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.latestScorecardAgencyID = agencyID
	if f.scorecardErr != nil {
		return compliance.Scorecard{}, f.scorecardErr
	}
	return f.scorecard, nil
}

func (f *fakePublicationStore) LatestTripUpdatesDiagnostics(context.Context, string) (compliance.TripUpdatesDiagnosticsSummary, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.tripDiagnosticsErr != nil {
		return compliance.TripUpdatesDiagnosticsSummary{}, f.tripDiagnosticsErr
	}
	return f.tripDiagnostics, nil
}

func (f *fakePublicationStore) ListAuditLog(_ context.Context, agencyID string, limit int) ([]compliance.AuditLogRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.auditAgencyID = agencyID
	if f.auditErr != nil {
		return nil, f.auditErr
	}
	if limit <= 0 || limit > len(f.auditRows) {
		limit = len(f.auditRows)
	}
	return append([]compliance.AuditLogRecord(nil), f.auditRows[:limit]...), nil
}

func (f *fakePublicationStore) BuildAndStoreScorecard(_ context.Context, agencyID string, _ time.Time) (compliance.Scorecard, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
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

func (f *fakePublicationStore) RecentGTFSImports(_ context.Context, agencyID string, limit int) ([]compliance.GTFSImportRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.gtfsImportsErr != nil {
		return nil, f.gtfsImportsErr
	}
	if limit <= 0 || limit > len(f.gtfsImports) {
		limit = len(f.gtfsImports)
	}
	out := make([]compliance.GTFSImportRecord, 0, limit)
	for _, record := range f.gtfsImports {
		if record.AgencyID != "" && record.AgencyID != agencyID {
			continue
		}
		out = append(out, record)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (f *fakePublicationStore) GTFSSchedulePreview(_ context.Context, agencyID string, feedVersionID string, limit int) (compliance.GTFSSchedulePreview, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.gtfsPreviewErr != nil {
		return compliance.GTFSSchedulePreview{}, f.gtfsPreviewErr
	}
	preview := f.gtfsPreview
	if f.gtfsPreviews != nil {
		if mapped, ok := f.gtfsPreviews[feedVersionID]; ok {
			preview = mapped
		} else if len(f.gtfsPreviews) > 0 && preview.FeedVersionID == "" {
			return compliance.GTFSSchedulePreview{}, errors.New("GTFS preview not found")
		}
	}
	if preview.AgencyID == "" {
		preview.AgencyID = agencyID
	}
	if preview.FeedVersionID == "" {
		preview.FeedVersionID = feedVersionID
	}
	if preview.RowLimit == 0 {
		preview.RowLimit = limit
	}
	if len(preview.Routes) > limit {
		preview.Routes = append([]compliance.GTFSScheduleRoutePreview(nil), preview.Routes[:limit]...)
	}
	if len(preview.Stops) > limit {
		preview.Stops = append([]compliance.GTFSScheduleStopPreview(nil), preview.Stops[:limit]...)
	}
	if len(preview.Trips) > limit {
		preview.Trips = append([]compliance.GTFSScheduleTripPreview(nil), preview.Trips[:limit]...)
	}
	if len(preview.Calendar) > limit {
		preview.Calendar = append([]compliance.GTFSScheduleCalendarPreview(nil), preview.Calendar[:limit]...)
	}
	if len(preview.Frequencies) > limit {
		preview.Frequencies = append([]compliance.GTFSScheduleFrequencyPreview(nil), preview.Frequencies[:limit]...)
	}
	return preview, nil
}

func (f *fakePublicationStore) RecentGTFSDrafts(_ context.Context, agencyID string, limit int) ([]compliance.GTFSDraftRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.gtfsDraftsErr != nil {
		return nil, f.gtfsDraftsErr
	}
	if limit <= 0 || limit > len(f.gtfsDrafts) {
		limit = len(f.gtfsDrafts)
	}
	out := make([]compliance.GTFSDraftRecord, 0, limit)
	for _, record := range f.gtfsDrafts {
		if record.AgencyID != "" && record.AgencyID != agencyID {
			continue
		}
		out = append(out, record)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (f *fakePublicationStore) RecentGTFSDraftPublishes(_ context.Context, agencyID string, limit int) ([]compliance.GTFSDraftPublishRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.gtfsDraftPublishesErr != nil {
		return nil, f.gtfsDraftPublishesErr
	}
	if limit <= 0 || limit > len(f.gtfsDraftPublishes) {
		limit = len(f.gtfsDraftPublishes)
	}
	out := make([]compliance.GTFSDraftPublishRecord, 0, limit)
	for _, record := range f.gtfsDraftPublishes {
		if record.ID == 0 && record.DraftID == "" {
			continue
		}
		out = append(out, record)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (f *fakePublicationStore) RecentFeedVersions(_ context.Context, agencyID string, limit int) ([]compliance.FeedVersionRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.feedVersionsErr != nil {
		return nil, f.feedVersionsErr
	}
	if limit <= 0 || limit > len(f.feedVersions) {
		limit = len(f.feedVersions)
	}
	out := make([]compliance.FeedVersionRecord, 0, limit)
	for _, record := range f.feedVersions {
		if record.AgencyID != "" && record.AgencyID != agencyID {
			continue
		}
		out = append(out, record)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (f *fakePublicationStore) LatestReliabilityFeedHealth(context.Context, string, int) ([]compliance.ReliabilityFeedHealthRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.reliabilityHealthErr != nil {
		return nil, f.reliabilityHealthErr
	}
	return append([]compliance.ReliabilityFeedHealthRecord(nil), f.reliabilityHealth...), nil
}

func (f *fakePublicationStore) ReliabilityIncidentRollup(context.Context, string, time.Time, int) (compliance.ReliabilityIncidentRollup, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.reliabilityIncidentsErr != nil {
		return compliance.ReliabilityIncidentRollup{}, f.reliabilityIncidentsErr
	}
	return f.reliabilityIncidents, nil
}

type fakeRealtimeArtifacts struct {
	mu       sync.Mutex
	payloads map[string][]byte
	errors   map[string]error
	calls    map[string]int
}

func (f *fakeRealtimeArtifacts) RealtimePB(_ context.Context, feedType string) ([]byte, string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.calls == nil {
		f.calls = map[string]int{}
	}
	f.calls[feedType]++
	if err := f.errors[feedType]; err != nil {
		return nil, "", err
	}
	if payload := f.payloads[feedType]; len(payload) > 0 {
		return append([]byte(nil), payload...), "internal_builder", nil
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
	token       string
	rebindCalls int
}

func (f *fakeDeviceStoreWithToken) Verify(context.Context, devices.VerifyInput) (devices.Credential, error) {
	return devices.Credential{}, nil
}

func (f *fakeDeviceStoreWithToken) Rebind(_ context.Context, input devices.RebindInput) (devices.RebindResult, error) {
	f.rebindCalls++
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
	mux.Handle("/admin/operations/assets/operations.js", adminRead(http.HandlerFunc(h.operationsAsset)))
	mux.Handle("/admin/operations", adminRead(http.HandlerFunc(h.operationsRoot)))
	mux.Handle("/admin/operations.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
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
	mux.Handle("/admin/operations/validation-center", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/validation-center.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/gtfs-workbench", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/gtfs-workbench.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/gtfs-import", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if r.Method == http.MethodPost {
			r.Body = http.MaxBytesReader(w, r.Body, gtfsBrowserImportMaxBytes+gtfsBrowserImportMemoryBytes)
		}
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
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
	mux.Handle("/admin/operations/validation-health/refresh.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if r.Method == http.MethodPost {
			r.Body = http.MaxBytesReader(w, r.Body, validationHealthPostMaxBytes)
		}
		adminRead(http.HandlerFunc(h.operationsValidationHealthRefreshCommandJSON)).ServeHTTP(w, r)
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
	mux.Handle("/admin/operations/maintenance", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/maintenance.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/telemetry-simulator", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/telemetry-simulator.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/help", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/help/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/help.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/help.json/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		adminRead(http.HandlerFunc(h.operationsRoot)).ServeHTTP(w, r)
	}))
	mux.Handle("/admin/operations/", adminRead(http.HandlerFunc(h.operationsRoot)))
	return mux
}

func newLocalAdminLoginTestHandler(t *testing.T, appEnv string, enabled string) http.Handler {
	t.Helper()
	t.Setenv("APP_ENV", appEnv)
	t.Setenv("AGENCY_ID", "demo-agency")
	t.Setenv("ADMIN_JWT_SECRET", "test-admin-secret")
	t.Setenv("ADMIN_JWT_ISSUER", "test-issuer")
	t.Setenv("ADMIN_JWT_AUDIENCE", "test-audience")
	t.Setenv("CSRF_SECRET", "test-csrf")
	t.Setenv("LOCAL_ADMIN_LOGIN_ENABLED", enabled)
	t.Setenv("LOCAL_ADMIN_SESSION_TTL", "2h")
	t.Setenv("LOCAL_ADMIN_SUBJECT", "admin@example.com")
	cfg := auth.JWTConfig{Secrets: []string{"test-admin-secret"}, Issuer: "test-issuer", Audience: "test-audience", TTL: time.Hour}
	verifier, err := auth.NewVerifier(cfg)
	if err != nil {
		t.Fatalf("verifier: %v", err)
	}
	admin := auth.NewMiddleware(verifier, auth.StaticRoleStore{Roles: []auth.Role{auth.RoleAdmin, auth.RoleEditor, auth.RoleOperator, auth.RoleReadOnly}}, "test-csrf")
	return newHandlerWithRealtime(
		"demo-agency",
		fakeScheduleBuilder{},
		&fakePublicationStore{},
		fakeDeviceStore{},
		fakePinger{},
		admin,
		&fakeRealtimeArtifacts{},
	)
}

func extractLocalLoginState(t *testing.T, body string) string {
	t.Helper()
	prefix := `name="state" value="`
	start := strings.Index(body, prefix)
	if start < 0 {
		t.Fatalf("local login body missing state input: %s", body)
	}
	start += len(prefix)
	end := strings.Index(body[start:], `"`)
	if end < 0 {
		t.Fatalf("local login body state input is not closed: %s", body)
	}
	state := body[start : start+end]
	if state == "" {
		t.Fatalf("local login state is empty")
	}
	return state
}

func localAdminLoginSessionCookie(t *testing.T, handler http.Handler) *http.Cookie {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/admin/local-login", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET local login status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	state := extractLocalLoginState(t, rr.Body.String())
	req = httptest.NewRequest(http.MethodPost, "http://localhost:8080/admin/local-login", strings.NewReader("state="+url.QueryEscape(state)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("POST local login status = %d, want 303: %s", rr.Code, rr.Body.String())
	}
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "admin_session" {
			return cookie
		}
	}
	t.Fatalf("local login did not set admin_session cookie")
	return nil
}

type fakeGTFSImportRunner struct {
	result               gtfs.ImportResult
	err                  error
	calls                int
	opts                 gtfs.ImportOptions
	payload              []byte
	pathExistsDuringCall bool
}

func (f *fakeGTFSImportRunner) ImportZip(ctx context.Context, opts gtfs.ImportOptions) (gtfs.ImportResult, error) {
	f.calls++
	f.opts = opts
	if _, err := os.Stat(opts.ZipPath); err == nil {
		f.pathExistsDuringCall = true
	}
	payload, err := os.ReadFile(opts.ZipPath)
	if err == nil {
		f.payload = append([]byte(nil), payload...)
	}
	if f.err != nil {
		return f.result, f.err
	}
	return f.result, ctx.Err()
}

func gtfsImportMultipartBody(t *testing.T, filename string, payload []byte, fields map[string]string) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatalf("write field %s: %v", key, err)
		}
	}
	part, err := writer.CreateFormFile("gtfs_zip", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(payload); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	return body, writer.FormDataContentType()
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
