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
	"testing"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
	"open-transit-rt/internal/devices"
	"open-transit-rt/internal/feed/schedule"
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
	result       compliance.ValidationResult
	discovery    compliance.FeedDiscovery
	discoveryErr error
}

func (f *fakePublicationStore) BootstrapPublication(context.Context, compliance.BootstrapInput) error {
	return nil
}

func (f *fakePublicationStore) FeedDiscovery(context.Context, string, time.Time) (compliance.FeedDiscovery, error) {
	if f.discoveryErr != nil {
		return compliance.FeedDiscovery{}, f.discoveryErr
	}
	return f.discovery, nil
}

func (f *fakePublicationStore) UpsertConsumer(context.Context, compliance.ConsumerInput) (compliance.ConsumerRecord, error) {
	return compliance.ConsumerRecord{}, nil
}

func (f *fakePublicationStore) ListConsumers(context.Context, string) ([]compliance.ConsumerRecord, error) {
	return nil, nil
}

func (f *fakePublicationStore) LatestScorecard(context.Context, string) (compliance.Scorecard, error) {
	return compliance.Scorecard{}, nil
}

func (f *fakePublicationStore) BuildAndStoreScorecard(context.Context, string, time.Time) (compliance.Scorecard, error) {
	return compliance.Scorecard{}, nil
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
