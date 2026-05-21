package main

import "testing"

func TestTopDashboardIssuesSkipsHealthyAndCapsAtThree(t *testing.T) {
	issues := []operationsOperatorIssue{
		{ID: "healthy", Status: checklistStatusOK},
		{ID: "missing-feed", Status: checklistStatusMissing},
		{ID: "blocked-telemetry", Status: checklistStatusBlocked},
		{ID: "review-validation", Status: checklistStatusNeedsReview},
		{ID: "unknown-extra", Status: checklistStatusUnknown},
	}
	got := topDashboardIssues(issues)
	if len(got) != dashboardTopIssueLimit {
		t.Fatalf("top issues = %d, want %d: %+v", len(got), dashboardTopIssueLimit, got)
	}
	want := []string{"missing-feed", "blocked-telemetry", "review-validation"}
	for i, id := range want {
		if got[i].ID != id {
			t.Fatalf("top issue %d = %q, want %q: %+v", i, got[i].ID, id, got)
		}
	}
}

func TestDashboardHealthyFallbackFillsWhenFewerThanThreeIssues(t *testing.T) {
	categories := []operationsDashboardCategory{
		{ID: "setup", Label: "Setup", Status: checklistStatusOK, HealthySignal: "setup current"},
		{ID: "feeds", Label: "Feeds", Status: operationsStatusReady, HealthySignal: "feeds current"},
		{ID: "validation", Label: "Validation", Status: checklistStatusNeedsReview, HealthySignal: "needs review"},
	}
	got := dashboardHealthyFallback(categories, 2)
	if len(got) != 2 {
		t.Fatalf("healthy fallback = %d, want 2: %+v", len(got), got)
	}
	if got[0].ID != "setup" || got[1].ID != "feeds" {
		t.Fatalf("healthy fallback = %+v, want healthy/current rows first", got)
	}
	got = dashboardHealthyFallback(categories, 3)
	if len(got) != 3 || got[2].ID != "validation" {
		t.Fatalf("healthy fallback should fill remaining slots with compact summaries: %+v", got)
	}
}
