package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	domainalerts "open-transit-rt/internal/alerts"
	"open-transit-rt/internal/auth"
	feedalerts "open-transit-rt/internal/feed/alerts"
	"open-transit-rt/internal/gtfs"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

func TestAlertsHandlersReturnValidFeedAndDebug(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	handler := newHandler(fakeAlertsBuilder{snapshot: feedalerts.Snapshot{
		AgencyID:                      "demo-agency",
		GeneratedAt:                   generatedAt,
		Status:                        feedalerts.StatusEmpty,
		Reason:                        "no_published_active_alerts",
		DiagnosticsPersistenceOutcome: "stored",
	}}, &fakeAlertStore{}, okPinger{})

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
	var debug feedalerts.Debug
	if err := json.Unmarshal(rr.Body.Bytes(), &debug); err != nil {
		t.Fatalf("decode debug: %v", err)
	}
	if debug.Status != feedalerts.StatusEmpty || debug.DiagnosticsPersistenceOutcome != "stored" {
		t.Fatalf("debug = %+v, want empty stored diagnostics", debug)
	}
}

func TestAlertsAdminCreateAndReconcile(t *testing.T) {
	store := &fakeAlertStore{}
	handler := newHandler(fakeAlertsBuilder{}, store, okPinger{})
	body := []byte(`{"agency_id":"demo-agency","alert_key":"alert-1","header_text":"Stop closed","actor_id":"operator@example.com","publish":true}`)
	req := httptest.NewRequest(http.MethodPost, "/admin/alerts", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("create status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.upsert.AlertKey != "alert-1" || !store.upsert.Publish {
		t.Fatalf("upsert = %+v, want published alert input", store.upsert)
	}

	req = httptest.NewRequest(http.MethodPost, "/admin/alerts/reconcile-cancellations", bytes.NewReader([]byte(`{"agency_id":"demo-agency","actor_id":"operator@example.com"}`)))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("reconcile status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if !store.reconciled {
		t.Fatalf("reconciled = false, want true")
	}
}

func TestAlertsAdminRejectsUnauthenticatedAccess(t *testing.T) {
	handler := newHandlerWithAuth(fakeAlertsBuilder{}, &fakeAlertStore{}, okPinger{}, authRejectAll{})
	req := httptest.NewRequest(http.MethodPost, "/admin/alerts", bytes.NewReader([]byte(`{}`)))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rr.Code)
	}
}

func TestAlertsConsoleRendersEmptyStateAndRejectsUnauthenticated(t *testing.T) {
	handler := newHandlerWithAuth(fakeAlertsBuilder{}, &fakeAlertStore{}, okPinger{}, authRejectAll{})
	req := httptest.NewRequest(http.MethodGet, "/admin/alerts/console", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rr.Code)
	}

	handler = newHandler(fakeAlertsBuilder{}, &fakeAlertStore{}, okPinger{})
	req = httptest.NewRequest(http.MethodGet, "/admin/alerts/console", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Alerts Console", "No alerts are recorded", "Operations Console"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body does not contain %q: %s", want, body)
		}
	}
}

func TestAlertsConsoleCreatePublishArchiveAndRoleBoundary(t *testing.T) {
	readOnly := newHandlerWithAuth(fakeAlertsBuilder{}, &fakeAlertStore{}, okPinger{}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "reader@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer,
	}})
	req := httptest.NewRequest(http.MethodPost, "/admin/alerts/console", strings.NewReader("agency_id=demo-agency&alert_key=a1&header_text=Alert"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	readOnly.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("read-only POST status = %d, want 403", rr.Code)
	}

	store := &fakeAlertStore{}
	handler := newHandler(fakeAlertsBuilder{}, store, okPinger{})
	req = httptest.NewRequest(http.MethodPost, "/admin/alerts/console", strings.NewReader("agency_id=demo-agency&alert_key=a1&header_text=Alert&route_id=r1&trip_id=t1&publish=on"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("create status = %d, want 303: %s", rr.Code, rr.Body.String())
	}
	if store.upsert.AlertKey != "a1" || !store.upsert.Publish || len(store.upsert.Entities) != 1 || store.upsert.Entities[0].RouteID != "r1" {
		t.Fatalf("upsert = %+v, want console create with affected entity", store.upsert)
	}

	req = httptest.NewRequest(http.MethodPost, "/admin/alerts/console/7/publish", strings.NewReader("agency_id=demo-agency"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther || store.publishedID != 7 {
		t.Fatalf("publish status = %d id=%d, want 303 id 7", rr.Code, store.publishedID)
	}

	req = httptest.NewRequest(http.MethodPost, "/admin/alerts/console/7/archive", strings.NewReader("agency_id=demo-agency&reason=done"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther || store.archivedID != 7 {
		t.Fatalf("archive status = %d id=%d, want 303 id 7", rr.Code, store.archivedID)
	}
}

func TestAlertsHandlersRejectWrongMethodAndReadyz(t *testing.T) {
	handler := newHandler(fakeAlertsBuilder{snapshot: feedalerts.Snapshot{GeneratedAt: time.Now().UTC()}}, &fakeAlertStore{}, okPinger{})
	req := httptest.NewRequest(http.MethodPost, "/public/gtfsrt/alerts.pb", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST status = %d, want 405", rr.Code)
	}

	handler = newHandler(fakeAlertsBuilder{}, &fakeAlertStore{}, errPinger{err: errors.New("down")})
	req = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rr.Code)
	}

	handler = newHandlerWithReadiness("demo-agency", fakeAlertsBuilder{}, &fakeAlertStore{}, okPinger{}, fakeActiveFeed{err: errors.New("missing active feed")}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject:  "test-admin",
		AgencyID: "demo-agency",
		Roles:    []auth.Role{auth.RoleAdmin, auth.RoleEditor, auth.RoleOperator, auth.RoleReadOnly},
		Method:   auth.MethodBearer,
	}})
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rr.Code)
	}

	handler = newHandlerWithReadiness("demo-agency", fakeAlertsBuilder{}, &fakeAlertStore{}, okPinger{}, fakeActiveFeed{}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject:  "test-admin",
		AgencyID: "demo-agency",
		Roles:    []auth.Role{auth.RoleAdmin, auth.RoleEditor, auth.RoleOperator, auth.RoleReadOnly},
		Method:   auth.MethodBearer,
	}})
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
}

type fakeAlertsBuilder struct {
	snapshot feedalerts.Snapshot
}

func (f fakeAlertsBuilder) Snapshot(context.Context, time.Time) (feedalerts.Snapshot, error) {
	return f.snapshot, nil
}

type fakeAlertStore struct {
	upsert      domainalerts.UpsertInput
	reconciled  bool
	publishedID int64
	archivedID  int64
	alerts      []domainalerts.Alert
}

func (f *fakeAlertStore) UpsertAlert(_ context.Context, input domainalerts.UpsertInput) (domainalerts.Alert, error) {
	f.upsert = input
	return domainalerts.Alert{ID: 1, AgencyID: input.AgencyID, AlertKey: input.AlertKey}, nil
}

func (f *fakeAlertStore) PublishAlert(_ context.Context, _ string, alertID int64, _ string, _ time.Time) (domainalerts.Alert, error) {
	f.publishedID = alertID
	return domainalerts.Alert{ID: alertID, Status: domainalerts.StatusPublished}, nil
}

func (f *fakeAlertStore) ArchiveAlert(_ context.Context, _ string, alertID int64, _ string, _ string, _ time.Time) error {
	f.archivedID = alertID
	return nil
}

func (f *fakeAlertStore) ListAlerts(context.Context, domainalerts.ListFilter) ([]domainalerts.Alert, error) {
	return f.alerts, nil
}

func (f *fakeAlertStore) ReconcileCanceledTripAlerts(context.Context, string, string, time.Time) (domainalerts.ReconcileResult, error) {
	f.reconciled = true
	return domainalerts.ReconcileResult{CreatedOrUpdated: 1}, nil
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

type fakeActiveFeed struct {
	err error
}

func (f fakeActiveFeed) ActiveFeedVersion(_ context.Context, agencyID string) (gtfs.FeedVersion, error) {
	if f.err != nil {
		return gtfs.FeedVersion{}, f.err
	}
	return gtfs.FeedVersion{ID: "feed-demo", AgencyID: agencyID}, nil
}

type authRejectAll struct{}

func (authRejectAll) Require(...auth.Role) func(http.Handler) http.Handler {
	return func(_ http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		})
	}
}
