package main

import (
	"fmt"
	"strings"
	"time"
)

const dashboardTopIssueLimit = 3

type operationsDashboardView struct {
	GeneratedAt       time.Time                     `json:"generated_at"`
	AgencyID          string                        `json:"agency_id"`
	Boundary          string                        `json:"boundary"`
	TopIssues         []operationsOperatorIssue     `json:"top_issues"`
	HealthySummaries  []operationsDashboardCategory `json:"healthy_summaries"`
	Categories        []operationsDashboardCategory `json:"categories"`
	HiddenIssueCount  int                           `json:"hidden_issue_count"`
	PrimaryNextAction string                        `json:"primary_next_action"`
}

type operationsDashboardCategory struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	Summary       string `json:"summary"`
	NextAction    string `json:"next_action"`
	AdminLink     string `json:"admin_link"`
	DoesNotProve  string `json:"does_not_prove"`
	HealthySignal string `json:"healthy_signal"`
}

func buildOperationsDashboard(page operationsPage) operationsDashboardView {
	top := topDashboardIssues(page.IssueCenter.Issues)
	categories := dashboardCategories(page)
	healthy := dashboardHealthyFallback(categories, dashboardTopIssueLimit-len(top))
	next := "Open the first top issue and complete its one action before moving to secondary tools."
	if len(top) == 0 {
		next = "No urgent issue is visible in the bounded private dashboard; continue setup, feed, realtime, connector, and maintenance review."
	}
	return operationsDashboardView{
		GeneratedAt:       page.GeneratedAt,
		AgencyID:          page.AgencyID,
		Boundary:          "Private dashboard summary only. It prioritizes existing records and does not mutate feeds, create evidence, contact external systems, change consumer status, or prove compliance, production readiness, uptime, vendor compatibility, production AVL reliability, or ETA quality.",
		TopIssues:         top,
		HealthySummaries:  healthy,
		Categories:        categories,
		HiddenIssueCount:  maxInt(0, len(page.IssueCenter.Issues)-len(top)),
		PrimaryNextAction: next,
	}
}

func topDashboardIssues(issues []operationsOperatorIssue) []operationsOperatorIssue {
	top := make([]operationsOperatorIssue, 0, dashboardTopIssueLimit)
	for _, issue := range issues {
		if normalizeChecklistStatus(issue.Status) == checklistStatusOK {
			continue
		}
		top = append(top, issue)
		if len(top) == dashboardTopIssueLimit {
			break
		}
	}
	return top
}

func dashboardHealthyFallback(categories []operationsDashboardCategory, limit int) []operationsDashboardCategory {
	if limit <= 0 {
		return nil
	}
	var healthy []operationsDashboardCategory
	for _, category := range categories {
		status := normalizeChecklistStatus(category.Status)
		if status == checklistStatusOK || category.Status == operationsStatusReady {
			healthy = append(healthy, category)
			if len(healthy) == limit {
				return healthy
			}
		}
	}
	for _, category := range categories {
		healthy = append(healthy, category)
		if len(healthy) == limit {
			return healthy
		}
	}
	return healthy
}

func dashboardCategories(page operationsPage) []operationsDashboardCategory {
	return []operationsDashboardCategory{
		dashboardCategory("setup", "Setup", cockpitMetadataStatus(page), cockpitMetadataSignal(page), "Continue the setup wizard until agency profile, feed URLs, license, and contact are recorded.", "/admin/operations/setup-wizard", "Setup completion does not prove agency approval or final-root ownership.", "Agency profile, contact, and license are visible."),
		dashboardCategory("schedule", "Schedule", cockpitActiveFeedStatus(page), cockpitActiveFeedSignal(page), "Import or review GTFS before relying on realtime outputs.", "/admin/operations/gtfs-workbench", "Schedule visibility does not prove the source GTFS is correct or validator-clean.", "An active published feed version is visible."),
		dashboardCategory("feed_urls", "Feed URLs", cockpitFeedHealthStatus(page), cockpitFeedHealthSignal(page), "Review five expected public feed paths and fix missing or stale rows first.", "/admin/operations/feed-health", "Feed URL visibility does not prove consumer ingestion, final-root ownership, or compliance.", "All expected feed paths have private review rows."),
		dashboardCategory("vehicle_data", "Vehicle Data", cockpitTelemetryStatus(page), telemetryEvidence(page), "Bind devices and review accepted, stale, and unmatched vehicle observations.", "/admin/operations/devices", "Vehicle data review does not prove vendor compatibility, hardware certification, or production AVL reliability.", "Recent vehicle telemetry is visible or explicitly absent."),
		dashboardCategory("realtime_output", "Realtime Output", realtimeStatus(page), realtimeSignal(page), "Check Vehicle Positions first, then Trip Updates and Alerts diagnostics.", "/admin/operations/realtime", "Realtime review does not prove production-grade ETA quality or consumer display.", "Realtime feed diagnostics are visible."),
		dashboardCategory("connectors", "Connectors", operationsStatusNeedsReview, fmt.Sprintf("%d connector categories; %d manifest entries", len(page.ConnectorHub.Categories), len(page.ConnectorHub.Registry.Entries)), "Review connector examples separately from deployment-owned configuration.", "/admin/operations/connectors", "Connector review does not prove a vendor integration or active connector exists.", "Connector catalog and checks are available for private review."),
		dashboardCategory("validation", "Validation", cockpitValidationStatus(page), fmt.Sprintf("overall=%s; tooling=%s", page.ValidationHealth.OverallStatus, page.ValidationHealth.ToolingStatus), "Review validator tooling and latest static/realtime results where configured.", "/admin/operations/validation-center", "Validator review does not prove compliance or consumer acceptance.", "Validator status is represented without raw commands."),
		dashboardCategory("deployment", "Deployment", maintenanceStatusToCockpit(page.Maintenance.OverallStatus), page.Maintenance.OverallStatus, "Review maintenance, backup, restore, and support-bundle readiness.", "/admin/operations/maintenance", "Deployment review does not prove uptime, SLA, hosted service, or production readiness.", "Maintenance summaries are available for private review."),
		dashboardCategory("users_access", "Users & Access", userAccessDashboardStatus(page), strings.Join(page.PrincipalRoles, ", "), "Review Login & Sessions and Users & Roles before adding staff access.", "/admin/operations/admin/sessions", "Access visibility does not prove SSO support or production multi-tenant readiness.", "Current signed-in subject, agency, and roles are visible."),
	}
}

func dashboardCategory(id, label, status, summary, next, link, doesNotProve, healthySignal string) operationsDashboardCategory {
	return operationsDashboardCategory{ID: id, Label: label, Status: status, Summary: summary, NextAction: next, AdminLink: link, DoesNotProve: doesNotProve, HealthySignal: healthySignal}
}

func userAccessDashboardStatus(page operationsPage) string {
	if len(page.PrincipalRoles) == 0 {
		return operationsStatusNeedsReview
	}
	return operationsStatusReady
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
