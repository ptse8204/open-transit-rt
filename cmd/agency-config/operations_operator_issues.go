package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"open-transit-rt/internal/compliance"
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
	ID               string `json:"id"`
	Label            string `json:"label"`
	Severity         string `json:"severity"`
	Status           string `json:"status"`
	Owner            string `json:"owner"`
	CurrentSignal    string `json:"current_signal"`
	WhyItMatters     string `json:"why_it_matters"`
	NextAction       string `json:"next_action"`
	RouteLink        string `json:"route_link"`
	Source           string `json:"source"`
	Freshness        string `json:"freshness"`
	DeduplicationKey string `json:"deduplication_key"`
	AdminLink        string `json:"admin_link"`
	SourceSignal     string `json:"source_signal"`
	SourceSurface    string `json:"source_surface"`
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
	builder.addValidationIssues(page)
	builder.addGTFSQualityIssues(page)
	builder.addRealtimeIssues(page)
	builder.addDeviceIssues(page)
	builder.addConnectorIssues(page)
	builder.addReadinessIssues(page)
	builder.addReliabilityIssues(page)
	builder.addMaintenanceIssues(page)
	issues := sortAndLimitOperatorIssues(builder.issues)
	if len(issues) == 0 {
		issues = []operationsOperatorIssue{{
			ID:               "continue_private_review",
			Label:            "Continue private review",
			Severity:         checklistStatusOK,
			Status:           checklistStatusOK,
			Owner:            "operator",
			CurrentSignal:    "issue center found no visible blockers",
			WhyItMatters:     "No blocked, missing, or needs-review issue is visible in the bounded private summary.",
			NextAction:       "Continue periodic feed, validator, realtime, connector, readiness, and maintenance review.",
			RouteLink:        "/admin/operations/maintenance",
			Source:           "Operations Console",
			Freshness:        issueFreshness(page.GeneratedAt),
			DeduplicationKey: "continue_private_review",
			AdminLink:        "/admin/operations/maintenance",
			SourceSignal:     "issue center found no visible blockers",
			SourceSurface:    "Operations Console",
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
	issue.ID = operatorIssueSafeID(firstNonEmpty(issue.ID, issue.Label, issue.Source))
	key := strings.TrimSpace(firstNonEmpty(issue.DeduplicationKey, issue.ID))
	if auditPrivateText(key) || unsafePrivateString(key) {
		key = issue.ID
	}
	key = operatorIssueSafeID(key)
	if issue.ID == "" || b.seen[issue.ID] || b.seen["key:"+key] {
		return
	}
	issue.Status = normalizeChecklistStatus(firstNonEmpty(issue.Status, issue.Severity))
	issue.Severity = normalizeChecklistStatus(firstNonEmpty(issue.Severity, issue.Status))
	issue.Owner = operatorIssueOwnerCategory(issue.Owner)
	issue.RouteLink = firstSafeAdminLink(firstNonEmpty(issue.RouteLink, issue.AdminLink))
	issue.AdminLink = firstSafeAdminLink(firstNonEmpty(issue.AdminLink, issue.RouteLink))
	issue.Label = operatorIssueDisplayText(issue.Label, "Review issue")
	issue.Source = operatorIssueDisplayText(firstNonEmpty(issue.Source, issue.SourceSurface, "Operations Console"), "Operations Console")
	issue.SourceSurface = firstNonEmpty(issue.SourceSurface, issue.Source)
	issue.SourceSurface = operatorIssueDisplayText(issue.SourceSurface, "Operations Console")
	issue.CurrentSignal = operatorIssueDisplayText(firstNonEmpty(issue.CurrentSignal, issue.SourceSignal, "not recorded"), "Private diagnostic omitted from browser table")
	issue.SourceSignal = firstNonEmpty(issue.SourceSignal, issue.CurrentSignal)
	issue.SourceSignal = operatorIssueDisplayText(issue.SourceSignal, "Private diagnostic omitted from browser table")
	issue.WhyItMatters = operatorIssueDisplayText(issue.WhyItMatters, "The source page has a private diagnostic that needs review.")
	issue.NextAction = operatorIssueDisplayText(issue.NextAction, "Open the linked private source page and review the bounded diagnostic there.")
	issue.Freshness = operatorIssueDisplayText(firstNonEmpty(issue.Freshness, "freshness not recorded separately"), "freshness not recorded separately")
	issue.DeduplicationKey = key
	b.seen[issue.ID] = true
	b.seen["key:"+issue.DeduplicationKey] = true
	b.issues = append(b.issues, issue)
}

func (b *operatorIssueBuilder) addCoreIssues(page operationsPage) {
	if strings.TrimSpace(page.ActiveFeedVersion) == "" {
		b.add(operationsOperatorIssue{
			ID:            "schedule_missing",
			Label:         "Import a schedule",
			Severity:      checklistStatusMissing,
			Owner:         "operator",
			CurrentSignal: firstNonEmpty(page.DiscoveryError, "no active schedule feed version recorded"),
			WhyItMatters:  "Static GTFS is the base for feed URLs, validators, matching, Vehicle Positions, Trip Updates, and Alerts.",
			NextAction:    "Import a GTFS ZIP or review the active schedule before relying on realtime output.",
			RouteLink:     "/admin/operations/gtfs-import",
			Source:        "Schedule",
			Freshness:     issueFreshness(page.GeneratedAt),
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
			CurrentSignal: cockpitMetadataSignal(page),
			WhyItMatters:  "Feed URLs need a public base URL, license, and contact before they are useful to operators or future authorized sharing.",
			NextAction:    "Open Agency Setup and confirm agency profile, public base URL, license, contact, and environment fields.",
			RouteLink:     "/admin/operations/setup-wizard",
			Source:        "Agency Setup",
			Freshness:     issueFreshness(page.GeneratedAt),
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
			CurrentSignal: fmt.Sprintf("overall=%s; tooling=%s", page.ValidationHealth.OverallStatus, page.ValidationHealth.ToolingStatus),
			WhyItMatters:  "Validator blockers can hide feed problems before operators share or rely on feed URLs.",
			NextAction:    "Review Validator Health, install or configure allowlisted tooling when needed, and rerun safe checks.",
			RouteLink:     "/admin/operations/validation-health",
			Source:        "Validation",
			Freshness:     issueFreshness(page.ValidationHealth.GeneratedAt),
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
			CurrentSignal: row.CurrentSignal,
			WhyItMatters:  feedIssueWhy(row.ID),
			NextAction:    row.NextAction,
			RouteLink:     firstOperatorLink(row.AdminLinks, "/admin/operations/feed-health"),
			Source:        "Feed Health",
			Freshness:     firstNonEmpty(row.Freshness, row.LastChecked, row.LastGenerated, "feed freshness not recorded"),
			AdminLink:     firstOperatorLink(row.AdminLinks, "/admin/operations/feed-health"),
			SourceSignal:  row.CurrentSignal,
			SourceSurface: "Feed Health",
		})
	}
}

func (b *operatorIssueBuilder) addValidationIssues(page operationsPage) {
	for _, blocker := range page.ValidationCenter.Blockers {
		status := normalizeChecklistStatus(blocker.Severity)
		if status == checklistStatusOK || b.seen[blocker.ID] {
			continue
		}
		b.add(operationsOperatorIssue{
			ID:               blocker.ID,
			Label:            blocker.Area,
			Severity:         status,
			Owner:            validationIssueOwner(blocker.Area, blocker.ReviewURL),
			CurrentSignal:    blocker.Signal,
			WhyItMatters:     "The validation center groups feed health, validator, GTFS quality, readiness, and consumer-prep blockers so operators can fix the first shared cause.",
			NextAction:       blocker.NextAction,
			RouteLink:        firstSafeAdminLink(blocker.ReviewURL),
			Source:           "Validation Center",
			Freshness:        issueFreshness(page.ValidationCenter.GeneratedAt),
			DeduplicationKey: blocker.ID,
			AdminLink:        firstSafeAdminLink(blocker.ReviewURL),
			SourceSignal:     blocker.Signal,
			SourceSurface:    "Validation Center",
		})
	}
}

func (b *operatorIssueBuilder) addGTFSQualityIssues(page operationsPage) {
	for _, row := range page.GTFSQualityGuidance.FixPlanner.Rows {
		status := issueSeverityFromStatus(row.Severity)
		if status == checklistStatusOK {
			continue
		}
		b.add(operationsOperatorIssue{
			ID:               "gtfs_quality_" + row.ID,
			Label:            "Fix " + firstNonEmpty(row.Family, "GTFS quality issue"),
			Severity:         status,
			Owner:            firstNonEmpty(row.LikelyOwner, "operator"),
			CurrentSignal:    fmt.Sprintf("%s; count=%d; affected=%s", row.IssueSummary, row.Count, row.AffectedFiles),
			WhyItMatters:     firstNonEmpty(row.WhyItMatters, "GTFS quality issues can affect imports, matching, and GTFS-Realtime usefulness."),
			NextAction:       firstNonEmpty(row.SafeFixSuggestion, row.BeforeValidationPlan, "Review the GTFS Quality fix planner."),
			RouteLink:        "/admin/operations/gtfs-quality",
			Source:           "GTFS Quality",
			Freshness:        issueFreshness(page.GTFSQualityGuidance.GeneratedAt),
			DeduplicationKey: "gtfs_quality_" + row.ID,
			AdminLink:        "/admin/operations/gtfs-quality",
			SourceSignal:     fmt.Sprintf("%s; count=%d; affected=%s", row.IssueSummary, row.Count, row.AffectedFiles),
			SourceSurface:    "GTFS Quality",
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
			CurrentSignal: summary.CurrentSignal,
			WhyItMatters:  "Fresh accepted telemetry is the input for publishable Vehicle Positions and defensible Trip Updates.",
			NextAction:    summary.NextAction,
			RouteLink:     "/admin/operations/realtime",
			Source:        "Realtime",
			Freshness:     issueFreshness(page.Realtime.GeneratedAt),
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
			CurrentSignal: issue.Signal,
			WhyItMatters:  "Realtime output should prefer missing or withheld data over false trip, location, or ETA certainty.",
			NextAction:    issue.NextAction,
			RouteLink:     firstSafeAdminLink(firstNonEmpty(issue.AdminLink, "/admin/operations/realtime")),
			Source:        "Realtime",
			Freshness:     issueFreshness(page.Realtime.GeneratedAt),
			AdminLink:     firstSafeAdminLink(firstNonEmpty(issue.AdminLink, "/admin/operations/realtime")),
			SourceSignal:  issue.Signal,
			SourceSurface: "Realtime",
		})
	}
}

func (b *operatorIssueBuilder) addDeviceIssues(page operationsPage) {
	for _, row := range append(append([]operationsDeviceFleetOnboardingRow{}, page.DeviceFleetOnboarding.InventoryRows...), append(page.DeviceFleetOnboarding.FreshnessTriageRows, page.DeviceFleetOnboarding.BindingReviewRows...)...) {
		status := issueSeverityFromStatus(row.Status)
		if status == checklistStatusOK || status == checklistStatusUnknown {
			continue
		}
		b.add(operationsOperatorIssue{
			ID:               "device_" + row.ID,
			Label:            row.Label,
			Severity:         status,
			Owner:            deviceIssueOwner(row),
			CurrentSignal:    row.CurrentSignal,
			WhyItMatters:     "Device onboarding and telemetry freshness determine whether Vehicle Positions can be useful without exposing credentials or false assignment certainty.",
			NextAction:       firstNonEmpty(row.OperatorStep, row.AdministratorStep),
			RouteLink:        "/admin/operations/devices",
			Source:           "Devices",
			Freshness:        issueFreshness(page.GeneratedAt),
			DeduplicationKey: "device_" + row.ID,
			AdminLink:        "/admin/operations/devices",
			SourceSignal:     row.CurrentSignal,
			SourceSurface:    "Devices",
		})
	}
}

func (b *operatorIssueBuilder) addConnectorIssues(page operationsPage) {
	for _, diagnostic := range page.ConnectorHub.Registry.Diagnostics {
		status := connectorDiagnosticSeverity(diagnostic.Level)
		if status == checklistStatusOK {
			continue
		}
		signal := strings.TrimSpace(strings.Join([]string{diagnostic.Code, diagnostic.Path, diagnostic.Message}, "; "))
		b.add(operationsOperatorIssue{
			ID:               "connector_registry_" + safeIssueID(firstNonEmpty(diagnostic.Code, diagnostic.Message)),
			Label:            "Fix connector registry diagnostic",
			Severity:         status,
			Owner:            "developer/integrator",
			CurrentSignal:    signal,
			WhyItMatters:     "Connector examples and manifests are the starting point for safe sidecar setup and must stay bounded, synthetic, and redacted.",
			NextAction:       "Open Connectors, fix the manifest or example diagnostic, then rerun offline connector checks.",
			RouteLink:        "/admin/operations/connectors",
			Source:           "Connectors",
			Freshness:        issueFreshness(page.ConnectorHub.GeneratedAt),
			DeduplicationKey: "connector_registry_" + safeIssueID(firstNonEmpty(diagnostic.Code, diagnostic.Message)),
			AdminLink:        "/admin/operations/connectors",
			SourceSignal:     signal,
			SourceSurface:    "Connectors",
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
			CurrentSignal: row.CurrentSignal,
			WhyItMatters:  row.WhyItMatters,
			NextAction:    row.WhatToDoNext,
			RouteLink:     firstOperatorLink(row.AdminLinks, "/admin/operations/readiness"),
			Source:        "Readiness",
			Freshness:     issueFreshness(page.ReadinessV2.GeneratedAt),
			AdminLink:     firstOperatorLink(row.AdminLinks, "/admin/operations/readiness"),
			SourceSignal:  row.CurrentSignal,
			SourceSurface: "Readiness",
		})
	}
}

func (b *operatorIssueBuilder) addReliabilityIssues(page operationsPage) {
	for _, row := range page.Reliability.Feeds {
		status := issueSeverityFromStatus(row.Status)
		if status == checklistStatusOK || b.seen["feed_"+row.FeedType] {
			continue
		}
		b.add(operationsOperatorIssue{
			ID:               "reliability_feed_" + row.FeedType,
			Label:            "Review " + row.FeedType + " reliability",
			Severity:         status,
			Owner:            "deployment owner",
			CurrentSignal:    reliabilityFeedSignal(row),
			WhyItMatters:     "Reliability snapshots help deployment owners spot freshness, endpoint, and validity problems without making SLA or uptime claims.",
			NextAction:       row.NextAction,
			RouteLink:        "/admin/operations/reliability",
			Source:           "Reliability",
			Freshness:        reliabilityFeedFreshness(row),
			DeduplicationKey: "reliability_feed_" + row.FeedType,
			AdminLink:        "/admin/operations/reliability",
			SourceSignal:     reliabilityFeedSignal(row),
			SourceSurface:    "Reliability",
		})
	}
	if status := issueSeverityFromStatus(page.Reliability.Incidents.Status); status != checklistStatusOK && status != checklistStatusUnknown {
		b.add(operationsOperatorIssue{
			ID:               "reliability_incidents",
			Label:            "Review reliability incidents",
			Severity:         status,
			Owner:            "deployment owner",
			CurrentSignal:    fmt.Sprintf("open incidents=%d; source=%s", page.Reliability.Incidents.Total, page.Reliability.Incidents.Source),
			WhyItMatters:     "Open private incidents can explain repeated feed, telemetry, or maintenance failures.",
			NextAction:       page.Reliability.Incidents.NextAction,
			RouteLink:        "/admin/operations/reliability",
			Source:           "Reliability",
			Freshness:        issueFreshness(page.Reliability.GeneratedAt),
			DeduplicationKey: "reliability_incidents",
			AdminLink:        "/admin/operations/reliability",
			SourceSignal:     fmt.Sprintf("open incidents=%d; source=%s", page.Reliability.Incidents.Total, page.Reliability.Incidents.Source),
			SourceSurface:    "Reliability",
		})
	}
	for _, section := range reliabilitySections(page) {
		status := issueSeverityFromStatus(section.status)
		if status == checklistStatusOK || status == checklistStatusUnknown {
			continue
		}
		b.add(operationsOperatorIssue{
			ID:               "reliability_" + section.id,
			Label:            section.label,
			Severity:         status,
			Owner:            "deployment owner",
			CurrentSignal:    section.summary,
			WhyItMatters:     "Self-hosted deployments need local operational checks before relying on feed uptime, recovery, or long-running jobs.",
			NextAction:       section.nextAction,
			RouteLink:        "/admin/operations/reliability",
			Source:           "Reliability",
			Freshness:        issueFreshness(page.Reliability.GeneratedAt),
			DeduplicationKey: "reliability_" + section.id,
			AdminLink:        "/admin/operations/reliability",
			SourceSignal:     section.summary,
			SourceSurface:    "Reliability",
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
			CurrentSignal: row.CurrentSignal,
			WhyItMatters:  "Maintenance gaps can block reliable self-hosted operation even when the browser workflow is usable.",
			NextAction:    row.NextAction,
			RouteLink:     "/admin/operations/maintenance",
			Source:        "Maintenance",
			Freshness:     issueFreshness(page.Maintenance.GeneratedAt),
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
		ID:               "additional_issue_rows",
		Label:            "More issue rows available",
		Severity:         checklistStatusNeedsReview,
		Status:           checklistStatusNeedsReview,
		Owner:            "operator",
		CurrentSignal:    "dashboard display limit reached",
		WhyItMatters:     "The dashboard keeps the first view focused, but the full private issue list remains available below.",
		NextAction:       fmt.Sprintf("Open the full issue table for %d additional issue row(s).", hidden),
		RouteLink:        "/admin/operations#all-operator-issues",
		Source:           "Operations Console",
		Freshness:        "same request",
		DeduplicationKey: "additional_issue_rows",
		AdminLink:        "/admin/operations#all-operator-issues",
		SourceSignal:     "dashboard display limit reached",
		SourceSurface:    "Operations Console",
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

func operatorIssueSafeID(value string) string {
	cleaned := strings.ToLower(strings.TrimSpace(value))
	if cleaned == "" {
		return "issue"
	}
	cleaned = strings.NewReplacer(
		"database_url", "private_value",
		"restore_database_url", "private_value",
		"authorization", "private_header",
		"bearer", "private_header",
		"set-cookie", "private_header",
		"cookie", "private_header",
		"token", "credential",
		"postgresql", "database",
		"postgres", "database",
		"payload_json", "private_payload",
		"raw_report", "private_report",
		"stdout", "private_output",
		"stderr", "private_output",
		"argv", "private_args",
		"file://", "private_file",
		"/users/", "/private/",
		"/var/lib", "/private",
		"/etc/", "/private/",
	).Replace(cleaned)
	return safeIssueID(cleaned)
}

func operatorIssueDisplayText(value, fallback string) string {
	cleaned := strings.TrimSpace(strings.Join(strings.Fields(value), " "))
	if cleaned == "" || auditPrivateText(cleaned) || unsafePrivateString(cleaned) {
		return fallback
	}
	const limit = 180
	if len(cleaned) > limit {
		return cleaned[:limit-15] + " [truncated]"
	}
	return cleaned
}

func operatorIssueOwnerCategory(owner string) string {
	lower := strings.ToLower(strings.TrimSpace(owner))
	switch lower {
	case "operator", "administrator", "deployment owner", "developer/integrator":
		return lower
	}
	switch {
	case strings.Contains(lower, "connector"), strings.Contains(lower, "integrator"), strings.Contains(lower, "developer"), strings.Contains(lower, "prediction"), strings.Contains(lower, "adapter"):
		return "developer/integrator"
	case strings.Contains(lower, "deployment"), strings.Contains(lower, "maintenance"), strings.Contains(lower, "reliability"), strings.Contains(lower, "backup"), strings.Contains(lower, "restore"), strings.Contains(lower, "database"), strings.Contains(lower, "proxy"), strings.Contains(lower, "systemd"):
		return "deployment owner"
	case strings.Contains(lower, "admin"), strings.Contains(lower, "token"), strings.Contains(lower, "credential"), strings.Contains(lower, "license"), strings.Contains(lower, "contact"), strings.Contains(lower, "public base"):
		return "administrator"
	case strings.Contains(lower, "technical maintainer") && !strings.Contains(lower, "gtfs source owner"):
		return "developer/integrator"
	default:
		return "operator"
	}
}

func validationIssueOwner(area string, link string) string {
	lower := strings.ToLower(area + " " + link)
	switch {
	case strings.Contains(lower, "connector"), strings.Contains(lower, "prediction"):
		return "developer/integrator"
	case strings.Contains(lower, "maintenance"), strings.Contains(lower, "reliability"), strings.Contains(lower, "setup"):
		return "deployment owner"
	case strings.Contains(lower, "validator"):
		return "administrator"
	default:
		return "operator"
	}
}

func deviceIssueOwner(row operationsDeviceFleetOnboardingRow) string {
	lower := strings.ToLower(row.Label + " " + row.AdministratorStep)
	if strings.Contains(lower, "token") || strings.Contains(lower, "secret") || strings.Contains(lower, "credential") || strings.Contains(lower, "administrator") {
		return "administrator"
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

func connectorDiagnosticSeverity(level string) string {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "error", "fatal", "blocked":
		return checklistStatusBlocked
	case "warning", "warn":
		return checklistStatusNeedsReview
	case "info", "ok":
		return checklistStatusOK
	default:
		return checklistStatusNeedsReview
	}
}

type operatorReliabilitySection struct {
	id         string
	label      string
	status     string
	summary    string
	nextAction string
}

func reliabilitySections(page operationsPage) []operatorReliabilitySection {
	return []operatorReliabilitySection{
		{id: "backup_restore", label: "Review backup and restore diagnostics", status: page.Reliability.BackupRestore.Status, summary: page.Reliability.BackupRestore.Summary, nextAction: page.Reliability.BackupRestore.NextAction},
		{id: "alerting", label: "Review alerting diagnostics", status: page.Reliability.Alerting.Status, summary: page.Reliability.Alerting.Summary, nextAction: page.Reliability.Alerting.NextAction},
		{id: "availability_sampling", label: "Review availability sampling", status: page.Reliability.AvailabilitySampling.Status, summary: page.Reliability.AvailabilitySampling.Summary, nextAction: page.Reliability.AvailabilitySampling.NextAction},
		{id: "long_running_operations", label: "Review long-running operations", status: page.Reliability.LongRunningOperations.Status, summary: page.Reliability.LongRunningOperations.Summary, nextAction: page.Reliability.LongRunningOperations.NextAction},
	}
}

func reliabilityFeedSignal(row compliance.ReliabilityFeedRow) string {
	parts := []string{
		"status=" + firstNonEmpty(row.Status, checklistStatusUnknown),
		"source=" + firstNonEmpty(row.Source, "not recorded"),
	}
	if row.EndpointAvailable != nil {
		parts = append(parts, fmt.Sprintf("endpoint_available=%t", *row.EndpointAvailable))
	}
	if row.FreshnessSeconds != nil {
		parts = append(parts, fmt.Sprintf("freshness=%.0fs", *row.FreshnessSeconds))
	}
	if row.InvalidResponsePercent != nil {
		parts = append(parts, fmt.Sprintf("invalid_response=%.1f%%", *row.InvalidResponsePercent))
	}
	if row.MatchedVehiclePercent != nil {
		parts = append(parts, fmt.Sprintf("matched_vehicles=%.1f%%", *row.MatchedVehiclePercent))
	}
	if row.CoveragePercent != nil {
		parts = append(parts, fmt.Sprintf("coverage=%.1f%%", *row.CoveragePercent))
	}
	return strings.Join(parts, "; ")
}

func reliabilityFeedFreshness(row compliance.ReliabilityFeedRow) string {
	if row.SnapshotAt == nil || row.SnapshotAt.IsZero() {
		return "snapshot not recorded"
	}
	return "snapshot " + formatTimeForText(row.SnapshotAt)
}

func issueFreshness(t time.Time) string {
	if t.IsZero() {
		return "freshness not recorded separately"
	}
	return "generated " + formatTimeForText(&t)
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
