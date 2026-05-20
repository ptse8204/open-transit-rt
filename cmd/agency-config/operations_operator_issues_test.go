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
	seenKeys := map[string]bool{}
	seenSources := map[string]bool{}
	for _, issue := range center.Issues {
		if seen[issue.ID] {
			t.Fatalf("duplicate issue id %q in %+v", issue.ID, center.Issues)
		}
		seen[issue.ID] = true
		if seenKeys[issue.DeduplicationKey] {
			t.Fatalf("duplicate issue dedupe key %q in %+v", issue.DeduplicationKey, center.Issues)
		}
		seenKeys[issue.DeduplicationKey] = true
		seenSources[issue.Source] = true
		if issue.Owner == "" || issue.CurrentSignal == "" || issue.WhyItMatters == "" || issue.NextAction == "" || issue.RouteLink == "" || issue.Source == "" || issue.Freshness == "" || issue.DeduplicationKey == "" {
			t.Fatalf("issue missing required V2 operator fields: %+v", issue)
		}
		if !strings.HasPrefix(issue.RouteLink, "/admin/") {
			t.Fatalf("issue route link is not private admin route: %+v", issue)
		}
		switch issue.Owner {
		case "operator", "administrator", "deployment owner", "developer/integrator":
		default:
			t.Fatalf("issue owner = %q, want plain owner category: %+v", issue.Owner, issue)
		}
	}
	for _, want := range []string{"feed_trip_updates", "realtime_fleet"} {
		if !seen[want] {
			t.Fatalf("issue center missing %q in %+v", want, center.Issues)
		}
	}
	for _, wantSource := range []string{"Feed Health", "Validation Center", "Realtime", "Devices", "Readiness", "Maintenance"} {
		if !seenSources[wantSource] {
			t.Fatalf("issue center missing source %q in %+v", wantSource, center.Issues)
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
	for _, path := range []string{"/admin/operations", "/admin/operations/launchpad", "/admin/operations/feed-health", "/admin/operations/validation-center", "/admin/operations/realtime"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
		}
		body := rr.Body.String()
		if !strings.Contains(body, "Top Issues") && path == "/admin/operations" {
			t.Fatalf("dashboard missing issue center: %s", body)
		}
		if path == "/admin/operations" {
			for _, want := range []string{"Category Summary", "Owner", "Current signal", "Why it matters", "Source", "Freshness", "Fix Trip Updates", "Review vehicle freshness", "All issue rows", "dashboard-top-issue-"} {
				if !strings.Contains(body, want) {
					t.Fatalf("dashboard missing %q: %s", want, body)
				}
			}
			if got := strings.Count(body, `id="dashboard-top-issue-`); got > 3 {
				t.Fatalf("dashboard top issue count = %d, want at most 3: %s", got, body)
			}
		}
		if path == "/admin/operations/launchpad" {
			for _, want := range []string{"Operator Issue Center", "Dedupe key", "Fix Trip Updates", "Review vehicle freshness"} {
				if !strings.Contains(body, want) {
					t.Fatalf("launchpad missing %q: %s", want, body)
				}
			}
		}
		if path == "/admin/operations/feed-health" || path == "/admin/operations/realtime" {
			for _, want := range []string{"Open issue center", "GTFS-RT Usefulness", "Vehicle Positions", "Trip Updates", "Alerts", "missing", "publishable"} {
				if !strings.Contains(body, want) {
					t.Fatalf("%s missing %q: %s", path, want, body)
				}
			}
		}
		if path == "/admin/operations/validation-center" && !strings.Contains(body, "Open issue center") {
			t.Fatalf("validation center missing issue center link: %s", body)
		}
	}
}

func TestOperationsIssueCenterDeduplicatesByStableKey(t *testing.T) {
	builder := operatorIssueBuilder{seen: map[string]bool{}}
	builder.add(operationsOperatorIssue{
		ID:               "feed_trip_updates",
		Label:            "Fix Trip Updates",
		Severity:         checklistStatusMissing,
		Owner:            "developer/integrator",
		CurrentSignal:    "missing feed row",
		WhyItMatters:     "Trip Updates should not be relied on when missing.",
		NextAction:       "Open feed health.",
		RouteLink:        "/admin/operations/feed-health",
		Source:           "Feed Health",
		Freshness:        "generated 2026-05-19T12:00:00Z",
		DeduplicationKey: "trip_updates_missing",
	})
	builder.add(operationsOperatorIssue{
		ID:               "validation_trip_updates",
		Label:            "Trip Updates validation",
		Severity:         checklistStatusMissing,
		Owner:            "developer/integrator",
		CurrentSignal:    "same missing Trip Updates blocker",
		WhyItMatters:     "Duplicate source row should collapse.",
		NextAction:       "Open validation center.",
		RouteLink:        "/admin/operations/validation-center",
		Source:           "Validation Center",
		Freshness:        "generated 2026-05-19T12:00:00Z",
		DeduplicationKey: "trip_updates_missing",
	})
	if len(builder.issues) != 1 {
		t.Fatalf("dedupe by stable key kept %d rows, want 1: %+v", len(builder.issues), builder.issues)
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
