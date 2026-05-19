package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
)

func TestOperationsIssueCenterPrioritizesDeduplicatesAndSummarizesRealtime(t *testing.T) {
	enableIssueCenterValidatorTools(t)
	page := buildIssueCenterTestPage(t)
	center := page.IssueCenter
	if center.Counts.Total == 0 || center.Counts.Displayed == 0 {
		t.Fatalf("issue center did not expose issues: %+v", center.Counts)
	}
	seen := map[string]bool{}
	for _, issue := range center.Issues {
		if seen[issue.ID] {
			t.Fatalf("duplicate issue id %q in %+v", issue.ID, center.Issues)
		}
		seen[issue.ID] = true
		if issue.Owner == "" || issue.WhyItMatters == "" || issue.NextAction == "" || issue.SourceSignal == "" {
			t.Fatalf("issue missing required operator fields: %+v", issue)
		}
	}
	for _, want := range []string{"feed_trip_updates", "realtime_fleet"} {
		if !seen[want] {
			t.Fatalf("issue center missing %q in %+v", want, center.Issues)
		}
	}
	if center.Issues[0].Severity != checklistStatusMissing {
		t.Fatalf("first issue severity = %q, want missing: %+v", center.Issues[0].Severity, center.Issues)
	}
	if center.Counts.Displayed > operatorIssueDisplayLimit {
		t.Fatalf("displayed issue count = %d, want at most %d", center.Counts.Displayed, operatorIssueDisplayLimit)
	}
	if center.Counts.Hidden > 0 {
		visibleSeen := map[string]bool{}
		for _, issue := range center.VisibleIssues {
			visibleSeen[issue.ID] = true
		}
		if !visibleSeen["additional_issue_rows"] {
			t.Fatalf("overflow issue row missing from visible issues: %+v", center.VisibleIssues)
		}
		if seen["additional_issue_rows"] {
			t.Fatalf("synthetic overflow row leaked into full issue list: %+v", center.Issues)
		}
	}
	states := map[string]string{}
	for _, row := range center.RealtimeFeeds {
		states[row.ID] = row.PublishState
		if !allowedRealtimePublishState(row.PublishState) {
			t.Fatalf("unexpected realtime publish state: %+v", row)
		}
		if row.Reason == "" || row.NextFix == "" || row.ValidatorConnection == "" || row.FeedHealthConnection == "" {
			t.Fatalf("realtime usefulness row missing detail: %+v", row)
		}
	}
	if states["vehicle_positions"] != "missing" || states["trip_updates"] != "missing" || states["alerts"] != "publishable" {
		t.Fatalf("unexpected realtime states: %+v", states)
	}
}

func TestOperationsDashboardAndRealtimeShowIssueCenter(t *testing.T) {
	enableIssueCenterValidatorTools(t)
	store := issueCenterTestStore(t)
	srv := issueCenterTestServer(store)
	for _, path := range []string{"/admin/operations", "/admin/operations/feed-health", "/admin/operations/realtime"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
		}
		body := rr.Body.String()
		if !strings.Contains(body, "Fix These First") && path == "/admin/operations" {
			t.Fatalf("dashboard missing issue center: %s", body)
		}
		if path == "/admin/operations" {
			for _, want := range []string{"Owner", "Why it matters", "Fix Trip Updates", "Review vehicle freshness", "All issue rows", "#all-operator-issues"} {
				if !strings.Contains(body, want) {
					t.Fatalf("dashboard missing %q: %s", want, body)
				}
			}
		}
		if path == "/admin/operations/feed-health" || path == "/admin/operations/realtime" {
			for _, want := range []string{"GTFS-RT Usefulness", "Vehicle Positions", "Trip Updates", "Alerts", "missing", "publishable"} {
				if !strings.Contains(body, want) {
					t.Fatalf("%s missing %q: %s", path, want, body)
				}
			}
		}
	}
}

func enableIssueCenterValidatorTools(t *testing.T) {
	t.Helper()
	t.Setenv("GTFS_VALIDATOR_PATH", writeScheduleValidator(t))
	t.Setenv("GTFS_RT_VALIDATOR_PATH", writeRealtimeValidator(t))
}

func allowedRealtimePublishState(state string) bool {
	switch state {
	case "publishable", "missing", "stale", "withheld", "blocked":
		return true
	default:
		return false
	}
}

func buildIssueCenterTestPage(t *testing.T) operationsPage {
	t.Helper()
	store := issueCenterTestStore(t)
	h := &handler{store: store, devices: fakeDeviceStore{}, telemetry: fakeTelemetryRepository{}, csrfSecret: "test-csrf"}
	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	return h.buildOperationsPage(req, auth.Principal{Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer}, "dashboard")
}

func issueCenterTestServer(store *fakePublicationStore) http.Handler {
	return newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}, telemetry: fakeTelemetryRepository{}, csrfSecret: "test-csrf"}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})
}

func issueCenterTestStore(t *testing.T) *fakePublicationStore {
	t.Helper()
	now := time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC)
	store := feedHealthTestStore(t)
	store.discovery.GeneratedAt = now
	store.discovery.PublicationEnvironment = "reference"
	store.tripDiagnostics = compliance.TripUpdatesDiagnosticsSummary{}
	feeds := store.discovery.Feeds[:0]
	for _, feed := range store.discovery.Feeds {
		if feed.FeedType == "trip_updates" {
			continue
		}
		feeds = append(feeds, feed)
	}
	store.discovery.Feeds = feeds
	return store
}
