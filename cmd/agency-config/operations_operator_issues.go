package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const operatorIssueDisplayLimit = 8

type operationsIssueCenterView struct {
	GeneratedAt    time.Time                       `json:"generated_at"`
	AgencyID       string                          `json:"agency_id"`
	Boundary       string                          `json:"boundary"`
	Issues         []operationsOperatorIssue       `json:"issues"`
	VisibleIssues  []operationsOperatorIssue       `json:"visible_issues"`
	RealtimeFeeds  []operationsRealtimeFeedIssue   `json:"realtime_feeds"`
	Counts         operationsOperatorIssueCounts   `json:"counts"`
	Recommendation operationsOperatorIssueGuidance `json:"recommendation"`
}

type operationsOperatorIssue struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Severity      string `json:"severity"`
	Status        string `json:"status"`
	Owner         string `json:"owner"`
	WhyItMatters  string `json:"why_it_matters"`
	NextAction    string `json:"next_action"`
	AdminLink     string `json:"admin_link"`
	SourceSignal  string `json:"source_signal"`
	SourceSurface string `json:"source_surface"`
}

type operationsRealtimeFeedIssue struct {
	ID                   string `json:"id"`
	Label                string `json:"label"`
	PublishState         string `json:"publish_state"`
	Reason               string `json:"reason"`
	NextFix              string `json:"next_fix"`
	ValidatorConnection  string `json:"validator_connection"`
	FeedHealthConnection string `json:"feed_health_connection"`
	AdminLink            string `json:"admin_link"`
}

type operationsOperatorIssueCounts struct {
	Total     int            `json:"total"`
	Displayed int            `json:"displayed"`
	Hidden    int            `json:"hidden"`
	ByStatus  map[string]int `json:"by_status"`
}

type operationsOperatorIssueGuidance struct {
	Summary    string `json:"summary"`
	NextAction string `json:"next_action"`
	AdminLink  string `json:"admin_link"`
}

type operatorIssueBuilder struct {
	seen   map[string]bool
	issues []operationsOperatorIssue
}

func buildOperationsIssueCenter(page operationsPage) operationsIssueCenterView {
	builder := operatorIssueBuilder{seen: map[string]bool{}}
	builder.addCoreIssues(page)
	builder.addFeedIssues(page)
	builder.addRealtimeIssues(page)
	builder.addReadinessIssues(page)
	builder.addMaintenanceIssues(page)
	issues := sortAndLimitOperatorIssues(builder.issues)
	if len(issues) == 0 {
		issues = []operationsOperatorIssue{{
			ID:            "continue_private_review",
			Label:         "Continue private review",
			Severity:      checklistStatusOK,
			Status:        checklistStatusOK,
			Owner:         "operator",
			WhyItMatters:  "No blocked, missing, or needs-review issue is visible in the bounded private summary.",
			NextAction:    "Continue periodic feed, validator, realtime, connector, readiness, and maintenance review.",
			AdminLink:     "/admin/operations/maintenance",
			SourceSignal:  "issue center found no visible blockers",
			SourceSurface: "Operations Console",
		}}
	}
	visible := visibleOperatorIssues(issues)
	return operationsIssueCenterView{
		GeneratedAt:    page.GeneratedAt,
		AgencyID:       page.AgencyID,
		Boundary:       "Private prioritized operator issue list only. It reads existing diagnostics and does not mutate feeds, create evidence, contact external systems, change consumer statuses, or prove compliance, production readiness, public launch, SLA, vendor compatibility, production AVL reliability, or ETA quality.",
		Issues:         issues,
		VisibleIssues:  visible,
		RealtimeFeeds:  buildOperatorRealtimeFeedIssues(page),
		Counts:         operatorIssueCounts(issues, min(len(issues), operatorIssueDisplayLimit)),
		Recommendation: operatorIssueRecommendation(issues),
	}
}

func (b *operatorIssueBuilder) add(issue operationsOperatorIssue) {
	issue.ID = strings.TrimSpace(issue.ID)
	if issue.ID == "" || b.seen[issue.ID] {
		return
	}
	issue.Status = normalizeChecklistStatus(firstNonEmpty(issue.Status, issue.Severity))
	issue.Severity = normalizeChecklistStatus(firstNonEmpty(issue.Severity, issue.Status))
	issue.Owner = firstNonEmpty(issue.Owner, "operator")
	issue.AdminLink = firstSafeAdminLink(issue.AdminLink)
	issue.SourceSurface = firstNonEmpty(issue.SourceSurface, "Operations Console")
	b.seen[issue.ID] = true
	b.issues = append(b.issues, issue)
}

func (b *operatorIssueBuilder) addCoreIssues(page operationsPage) {
	if strings.TrimSpace(page.ActiveFeedVersion) == "" {
		b.add(operationsOperatorIssue{
			ID:            "schedule_missing",
			Label:         "Import a schedule",
			Severity:      checklistStatusMissing,
			Owner:         "operator",
			WhyItMatters:  "Static GTFS is the base for feed URLs, validators, matching, Vehicle Positions, Trip Updates, and Alerts.",
			NextAction:    "Import a GTFS ZIP or review the active schedule before relying on realtime output.",
			AdminLink:     "/admin/operations/gtfs-import",
			SourceSignal:  firstNonEmpty(page.DiscoveryError, "no active schedule feed version recorded"),
			SourceSurface: "Schedule",
		})
	}
	if page.DiscoveryError != "" || strings.TrimSpace(page.Discovery.PublicBaseURL) == "" || !page.Discovery.Readiness.LicenseComplete || !page.Discovery.Readiness.ContactComplete {
		b.add(operationsOperatorIssue{
			ID:            "publication_metadata",
			Label:         "Review feed metadata",
			Severity:      issueSeverityFromStatus(cockpitMetadataStatus(page)),
			Owner:         "administrator",
			WhyItMatters:  "Feed URLs need a public base URL, license, and contact before they are useful to operators or future authorized sharing.",
			NextAction:    "Open Agency Setup and confirm agency profile, public base URL, license, contact, and environment fields.",
			AdminLink:     "/admin/operations/setup-wizard",
			SourceSignal:  cockpitMetadataSignal(page),
			SourceSurface: "Agency Setup",
		})
	}
	if page.ValidationHealth.OverallStatus != "" && normalizeIssueSourceStatus(page.ValidationHealth.OverallStatus) != checklistStatusOK {
		b.add(operationsOperatorIssue{
			ID:            "validator_health",
			Label:         "Fix validator health",
			Severity:      issueSeverityFromStatus(page.ValidationHealth.OverallStatus),
			Owner:         "administrator",
			WhyItMatters:  "Validator blockers can hide feed problems before operators share or rely on feed URLs.",
			NextAction:    "Review Validator Health, install or configure allowlisted tooling when needed, and rerun safe checks.",
			AdminLink:     "/admin/operations/validation-health",
			SourceSignal:  fmt.Sprintf("overall=%s; tooling=%s", page.ValidationHealth.OverallStatus, page.ValidationHealth.ToolingStatus),
			SourceSurface: "Validation",
		})
	}
}

func (b *operatorIssueBuilder) addFeedIssues(page operationsPage) {
	for _, row := range page.FeedHealth.Rows {
		status := normalizeChecklistStatus(row.Status)
		if status == checklistStatusOK {
			continue
		}
		b.add(operationsOperatorIssue{
			ID:            "feed_" + row.ID,
			Label:         "Fix " + row.Label,
			Severity:      status,
			Owner:         feedIssueOwner(row.ID),
			WhyItMatters:  feedIssueWhy(row.ID),
			NextAction:    row.NextAction,
			AdminLink:     firstOperatorLink(row.AdminLinks, "/admin/operations/feed-health"),
			SourceSignal:  row.CurrentSignal,
			SourceSurface: "Feed Health",
		})
	}
}

func (b *operatorIssueBuilder) addRealtimeIssues(page operationsPage) {
	summary := page.Realtime.Summary
	if normalizeChecklistStatus(summary.Status) != checklistStatusOK {
		b.add(operationsOperatorIssue{
			ID:            "realtime_fleet",
			Label:         "Review vehicle freshness",
			Severity:      summary.Status,
			Owner:         "operator",
			WhyItMatters:  "Fresh accepted telemetry is the input for publishable Vehicle Positions and defensible Trip Updates.",
			NextAction:    summary.NextAction,
			AdminLink:     "/admin/operations/realtime",
			SourceSignal:  summary.CurrentSignal,
			SourceSurface: "Realtime",
		})
	}
	for _, issue := range page.Realtime.Issues {
		status := normalizeChecklistStatus(issue.Severity)
		if status == checklistStatusOK {
			continue
		}
		if b.hasFeedIssueForRealtimeArea(issue.Area) {
			continue
		}
		b.add(operationsOperatorIssue{
			ID:            "realtime_" + safeIssueID(issue.Area),
			Label:         issue.Area,
			Severity:      status,
			Owner:         realtimeIssueOwner(issue.Area),
			WhyItMatters:  "Realtime output should prefer missing or withheld data over false trip, location, or ETA certainty.",
			NextAction:    issue.NextAction,
			AdminLink:     firstSafeAdminLink(firstNonEmpty(issue.AdminLink, "/admin/operations/realtime")),
			SourceSignal:  issue.Signal,
			SourceSurface: "Realtime",
		})
	}
}

func (b *operatorIssueBuilder) addReadinessIssues(page operationsPage) {
	for _, row := range page.ReadinessV2.Rows {
		status := normalizeChecklistStatus(row.Status)
		if status == checklistStatusOK {
			continue
		}
		b.add(operationsOperatorIssue{
			ID:            "readiness_" + row.ID,
			Label:         row.ReadinessItem,
			Severity:      status,
			Owner:         readinessIssueOwner(row),
			WhyItMatters:  row.WhyItMatters,
			NextAction:    row.WhatToDoNext,
			AdminLink:     firstOperatorLink(row.AdminLinks, "/admin/operations/readiness"),
			SourceSignal:  row.CurrentSignal,
			SourceSurface: "Readiness",
		})
	}
}

func (b *operatorIssueBuilder) addMaintenanceIssues(page operationsPage) {
	for _, row := range page.Maintenance.SummaryRows {
		status := issueSeverityFromStatus(row.Status)
		if status == checklistStatusOK || status == checklistStatusUnknown {
			continue
		}
		b.add(operationsOperatorIssue{
			ID:            "maintenance_" + row.ID,
			Label:         row.Label,
			Severity:      status,
			Owner:         "deployment owner",
			WhyItMatters:  "Maintenance gaps can block reliable self-hosted operation even when the browser workflow is usable.",
			NextAction:    row.NextAction,
			AdminLink:     "/admin/operations/maintenance",
			SourceSignal:  row.CurrentSignal,
			SourceSurface: "Maintenance",
		})
	}
}

func buildOperatorRealtimeFeedIssues(page operationsPage) []operationsRealtimeFeedIssue {
	rows := map[string]operationsFeedHealthRow{}
	for _, row := range page.FeedHealth.Rows {
		rows[row.ID] = row
	}
	return []operationsRealtimeFeedIssue{
		operatorRealtimeFeedIssue("vehicle_positions", "Vehicle Positions", rows["vehicle_positions"], vehiclePositionsPublishState(page, rows["vehicle_positions"])),
		operatorRealtimeFeedIssue("trip_updates", "Trip Updates", rows["trip_updates"], tripUpdatesPublishState(page, rows["trip_updates"])),
		operatorRealtimeFeedIssue("alerts", "Alerts", rows["alerts"], alertsPublishState(page, rows["alerts"])),
	}
}

func operatorRealtimeFeedIssue(id string, label string, row operationsFeedHealthRow, state string) operationsRealtimeFeedIssue {
	return operationsRealtimeFeedIssue{
		ID:                   id,
		Label:                label,
		PublishState:         state,
		Reason:               realtimeStateReason(id, state, row),
		NextFix:              realtimeStateNextFix(id, state, row),
		ValidatorConnection:  firstNonEmpty(row.ValidatorContext, "validator context not available"),
		FeedHealthConnection: firstNonEmpty(row.HealthContext, row.CurrentSignal, "feed-health context not available"),
		AdminLink:            firstOperatorLink(row.AdminLinks, "/admin/operations/feed-health"),
	}
}

func vehiclePositionsPublishState(page operationsPage, row operationsFeedHealthRow) string {
	switch normalizeChecklistStatus(row.Status) {
	case checklistStatusBlocked:
		return "blocked"
	case checklistStatusMissing:
		return "missing"
	}
	if page.TelemetryError != "" {
		return "blocked"
	}
	if len(page.Telemetry) == 0 {
		return "missing"
	}
	if page.StaleCount >= len(page.Telemetry) {
		return "stale"
	}
	return "publishable"
}

func tripUpdatesPublishState(page operationsPage, row operationsFeedHealthRow) string {
	switch normalizeChecklistStatus(row.Status) {
	case checklistStatusBlocked:
		return "blocked"
	case checklistStatusMissing:
		return "missing"
	}
	quality := page.TripUpdatesQuality
	if !quality.Recorded {
		return "missing"
	}
	if quality.StaleTelemetryRows > 0 {
		return "stale"
	}
	if quality.TripUpdatesEmitted == 0 || countViewTotal(quality.WithheldByReason) > 0 || quality.UnknownAssignments > 0 || quality.AmbiguousAssignments > 0 {
		return "withheld"
	}
	return "publishable"
}

func alertsPublishState(page operationsPage, row operationsFeedHealthRow) string {
	switch normalizeChecklistStatus(row.Status) {
	case checklistStatusBlocked:
		return "blocked"
	case checklistStatusMissing:
		return "missing"
	}
	if page.TripUpdatesQuality.CancellationAlertLinksMissing > 0 {
		return "blocked"
	}
	return "publishable"
}

func realtimeStateReason(id string, state string, row operationsFeedHealthRow) string {
	switch state {
	case "publishable":
		return "The private source signals are present enough for local operator review; keep validation and freshness checks running."
	case "missing":
		return "Required metadata, feed output, telemetry, diagnostics, or lifecycle records are not visible yet."
	case "stale":
		return "The latest telemetry or diagnostics include stale rows, so realtime output should stay conservative."
	case "withheld":
		return "The prediction boundary is intentionally withholding output because the current evidence is ambiguous or incomplete."
	case "blocked":
		return "A validator, feed-health, lifecycle, or repository blocker needs attention before relying on this feed."
	default:
		return firstNonEmpty(row.CurrentSignal, id+" state needs review")
	}
}

func realtimeStateNextFix(id string, state string, row operationsFeedHealthRow) string {
	if state == "publishable" {
		return "Continue private feed health, validator, and freshness review before any authorized external sharing."
	}
	if strings.TrimSpace(row.NextAction) != "" {
		return row.NextAction
	}
	switch id {
	case "vehicle_positions":
		return "Review Devices, Telemetry, and Vehicle Positions validation together."
	case "trip_updates":
		return "Review Prediction Lab, withheld reasons, stale telemetry, and Trip Updates validation."
	case "alerts":
		return "Review Alerts Console lifecycle rows and Alerts feed health."
	default:
		return "Open Feed Health and fix the first blocked row."
	}
}

func sortAndLimitOperatorIssues(issues []operationsOperatorIssue) []operationsOperatorIssue {
	out := append([]operationsOperatorIssue(nil), issues...)
	sort.SliceStable(out, func(i, j int) bool {
		left := operatorIssueRank(out[i].Severity)
		right := operatorIssueRank(out[j].Severity)
		if left != right {
			return left < right
		}
		leftSurface := operatorIssueSurfaceRank(out[i].SourceSurface)
		rightSurface := operatorIssueSurfaceRank(out[j].SourceSurface)
		if leftSurface != rightSurface {
			return leftSurface < rightSurface
		}
		return out[i].Label < out[j].Label
	})
	return out
}

func visibleOperatorIssues(issues []operationsOperatorIssue) []operationsOperatorIssue {
	if len(issues) <= operatorIssueDisplayLimit {
		return append([]operationsOperatorIssue(nil), issues...)
	}
	visible := append([]operationsOperatorIssue(nil), issues[:operatorIssueDisplayLimit]...)
	hidden := len(issues) - operatorIssueDisplayLimit
	visible = append(visible, operationsOperatorIssue{
		ID:            "additional_issue_rows",
		Label:         "More issue rows available",
		Severity:      checklistStatusNeedsReview,
		Status:        checklistStatusNeedsReview,
		Owner:         "operator",
		WhyItMatters:  "The dashboard keeps the first view focused, but the full private issue list remains available below.",
		NextAction:    fmt.Sprintf("Open the full issue table for %d additional issue row(s).", hidden),
		AdminLink:     "/admin/operations#all-operator-issues",
		SourceSignal:  "dashboard display limit reached",
		SourceSurface: "Operations Console",
	})
	return visible
}

func operatorIssueCounts(issues []operationsOperatorIssue, displayed int) operationsOperatorIssueCounts {
	counts := operationsOperatorIssueCounts{Total: len(issues), Displayed: displayed, ByStatus: map[string]int{}}
	if counts.Total > counts.Displayed {
		counts.Hidden = counts.Total - counts.Displayed
	}
	for _, issue := range issues {
		counts.ByStatus[normalizeChecklistStatus(issue.Severity)]++
	}
	return counts
}

func operatorIssueRecommendation(issues []operationsOperatorIssue) operationsOperatorIssueGuidance {
	if len(issues) == 0 {
		return operationsOperatorIssueGuidance{Summary: "No visible issue rows need immediate action.", NextAction: "Continue maintenance review.", AdminLink: "/admin/operations/maintenance"}
	}
	first := issues[0]
	return operationsOperatorIssueGuidance{
		Summary:    fmt.Sprintf("Start with %s because it is %s and owned by %s.", first.Label, first.Severity, first.Owner),
		NextAction: first.NextAction,
		AdminLink:  first.AdminLink,
	}
}

func issueSeverityFromStatus(status string) string {
	switch normalizeIssueSourceStatus(status) {
	case checklistStatusBlocked:
		return checklistStatusBlocked
	case checklistStatusMissing:
		return checklistStatusMissing
	case checklistStatusNeedsReview:
		return checklistStatusNeedsReview
	case checklistStatusOK:
		return checklistStatusOK
	default:
		return checklistStatusUnknown
	}
}

func normalizeIssueSourceStatus(status string) string {
	normalized := strings.ToLower(strings.TrimSpace(status))
	normalized = strings.NewReplacer(" ", "_", "-", "_").Replace(normalized)
	switch normalized {
	case checklistStatusOK, "ready", "ready_for_local_review", "recorded", "passed", "available", "publishable", "generated":
		return checklistStatusOK
	case checklistStatusBlocked, "failed", "unhealthy", "missing_tooling", "misconfigured_tooling":
		return checklistStatusBlocked
	case checklistStatusMissing, "not_run", "not_recorded", "not_available", "diagnostics_missing", "artifact_unavailable":
		return checklistStatusMissing
	case checklistStatusNeedsReview, "needs_action", "stale", "warning", "withheld", "empty_or_withheld", "degraded", "runnable", "configured_for_tests":
		return checklistStatusNeedsReview
	default:
		return normalizeChecklistStatus(normalized)
	}
}

func operatorIssueRank(status string) int {
	switch normalizeChecklistStatus(status) {
	case checklistStatusBlocked:
		return 0
	case checklistStatusMissing:
		return 1
	case checklistStatusNeedsReview:
		return 2
	case checklistStatusUnknown:
		return 3
	case checklistStatusOK:
		return 9
	default:
		return 4
	}
}

func operatorIssueSurfaceRank(surface string) int {
	switch strings.ToLower(strings.TrimSpace(surface)) {
	case "feed health":
		return 0
	case "realtime":
		return 1
	case "schedule":
		return 2
	case "agency setup":
		return 3
	case "readiness":
		return 4
	case "validation":
		return 5
	case "maintenance":
		return 6
	default:
		return 7
	}
}

func feedIssueOwner(id string) string {
	switch id {
	case "feeds_json", "schedule":
		return "administrator"
	case "vehicle_positions":
		return "operator"
	case "trip_updates":
		return "developer/integrator"
	case "alerts":
		return "operator"
	default:
		return "operator"
	}
}

func feedIssueWhy(id string) string {
	switch id {
	case "feeds_json":
		return "The discovery document is how operators confirm the expected feed set and metadata before any future sharing."
	case "schedule":
		return "Static GTFS is the foundation for route, stop, trip, service-day, and realtime matching work."
	case "vehicle_positions":
		return "Vehicle Positions are the first public realtime output target and need fresh telemetry plus validation context."
	case "trip_updates":
		return "Trip Updates must stay behind the adapter boundary and should be withheld when assignment or prediction evidence is weak."
	case "alerts":
		return "Alerts complete the realtime feed set and need lifecycle review before disruption messaging is relied on."
	default:
		return "Feed-health rows point to the first private action needed before relying on configured URLs."
	}
}

func readinessIssueOwner(row operationsReadinessV2Row) string {
	for _, link := range row.AdminLinks {
		if strings.Contains(link, "connectors") {
			return "developer/integrator"
		}
		if strings.Contains(link, "maintenance") || strings.Contains(link, "setup") {
			return "deployment owner"
		}
	}
	return "operator"
}

func realtimeIssueOwner(area string) string {
	lower := strings.ToLower(area)
	switch {
	case strings.Contains(lower, "trip updates") || strings.Contains(lower, "prediction"):
		return "developer/integrator"
	case strings.Contains(lower, "device") || strings.Contains(lower, "telemetry"):
		return "operator"
	default:
		return "operator"
	}
}

func (b *operatorIssueBuilder) hasFeedIssueForRealtimeArea(area string) bool {
	switch safeIssueID(area) {
	case "vehicle_positions":
		return b.seen["feed_vehicle_positions"]
	case "trip_updates":
		return b.seen["feed_trip_updates"]
	case "alerts":
		return b.seen["feed_alerts"]
	default:
		return false
	}
}

func firstOperatorLink(values []string, fallback string) string {
	for _, value := range values {
		if safe := firstSafeAdminLink(value); safe != "" {
			return safe
		}
	}
	return firstSafeAdminLink(fallback)
}

func safeIssueID(value string) string {
	out := strings.ToLower(strings.TrimSpace(value))
	out = strings.NewReplacer(" ", "_", "-", "_", "/", "_", "&", "and").Replace(out)
	out = strings.Trim(out, "_")
	if out == "" {
		return "issue"
	}
	return out
}
