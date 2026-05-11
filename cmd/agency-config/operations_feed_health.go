package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
)

type operationsFeedHealthView struct {
	GeneratedAt time.Time                  `json:"generated_at"`
	AgencyID    string                     `json:"agency_id"`
	Boundary    string                     `json:"boundary"`
	Rows        []operationsFeedHealthRow  `json:"rows"`
	Counts      operationsFeedHealthCounts `json:"counts"`
	ClaimFlags  operationsFeedHealthClaims `json:"claim_flags"`
}

type operationsFeedHealthRow struct {
	ID               string   `json:"id"`
	Label            string   `json:"label"`
	PublicPath       string   `json:"public_path"`
	Status           string   `json:"status"`
	StatusText       string   `json:"status_text"`
	CurrentSignal    string   `json:"current_signal"`
	WhatThisMeans    string   `json:"what_this_means"`
	Freshness        string   `json:"freshness"`
	ValidatorContext string   `json:"validator_context"`
	HealthContext    string   `json:"health_context"`
	NextAction       string   `json:"next_action"`
	DoesNotProve     string   `json:"does_not_prove"`
	AdminLinks       []string `json:"admin_links"`
	DocsLinks        []string `json:"docs_links"`
}

type operationsFeedHealthCounts struct {
	Rows     int            `json:"rows"`
	Statuses map[string]int `json:"statuses"`
}

type operationsFeedHealthClaims struct {
	ExternalEvidenceCreated    bool `json:"external_evidence_created"`
	ConsumerStatusesChanged    bool `json:"consumer_statuses_changed"`
	ComplianceClaimed          bool `json:"compliance_claimed"`
	ProductionReadinessClaimed bool `json:"production_readiness_claimed"`
	SLAClaimed                 bool `json:"sla_claimed"`
	UptimeGuaranteeClaimed     bool `json:"uptime_guarantee_claimed"`
	ConsumerAcceptanceClaimed  bool `json:"consumer_acceptance_claimed"`
	PublicLaunchClaimed        bool `json:"public_launch_claimed"`
}

func (h *handler) renderFeedHealth(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "feed-health")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "feed-health", page)
}

func (h *handler) renderFeedHealthJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "feed-health")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.FeedHealth)
}

func buildOperationsFeedHealth(page operationsPage) operationsFeedHealthView {
	rows := []operationsFeedHealthRow{
		buildFeedsJSONHealthRow(page),
		buildFeedHealthRow(page, "schedule", "Static GTFS Schedule"),
		buildFeedHealthRow(page, "vehicle_positions", "Vehicle Positions"),
		buildFeedHealthRow(page, "trip_updates", "Trip Updates"),
		buildFeedHealthRow(page, "alerts", "Alerts"),
	}
	return operationsFeedHealthView{
		GeneratedAt: page.GeneratedAt,
		AgencyID:    page.AgencyID,
		Boundary:    "Private authenticated feed-health dashboard only; viewing it creates no evidence, changes no consumer status, contacts no external party, and records no compliance, SLA, uptime proof, consumer-acceptance, public-launch, or production-readiness outcome.",
		Rows:        rows,
		Counts:      feedHealthCounts(rows),
		ClaimFlags:  operationsFeedHealthClaims{},
	}
}

func buildFeedsJSONHealthRow(page operationsPage) operationsFeedHealthRow {
	status := checklistStatusMissing
	statusText := "Missing metadata"
	signal := firstNonEmpty(page.DiscoveryError, "no feed discovery metadata is available")
	freshness := "not generated"
	next := "Publish or import a schedule, then store publication metadata so feeds.json can describe every required feed."
	if page.DiscoveryError == "" {
		freshness = "generated " + formatTimeForText(&page.Discovery.GeneratedAt)
		url := strings.TrimRight(page.Discovery.PublicBaseURL, "/") + "/public/feeds.json"
		if strings.TrimSpace(page.Discovery.PublicBaseURL) == "" {
			url = "missing public base URL"
		}
		signal = fmt.Sprintf("%s; %d feed records; all required listed=%t; all HTTPS=%t; discoverable=%t", url, len(page.Discovery.Feeds), page.Discovery.Readiness.AllRequiredFeedsListed, page.Discovery.Readiness.HTTPSURLs, page.Discovery.Readiness.Discoverable)
		status = checklistStatusNeedsReview
		statusText = "Review metadata"
		next = "Confirm feeds.json lists schedule, Vehicle Positions, Trip Updates, and Alerts with HTTPS public URLs, license metadata, contact metadata, and discoverability signals."
		if page.Discovery.Readiness.AllRequiredFeedsListed && page.Discovery.Readiness.HTTPSURLs && page.Discovery.Readiness.Discoverable && page.Discovery.Readiness.LicenseComplete && page.Discovery.Readiness.ContactComplete {
			status = checklistStatusOK
			statusText = "Discovery metadata present"
			next = "Review individual feed validator and freshness rows before treating the feed set as ready."
		}
	}
	return operationsFeedHealthRow{
		ID:               "feeds_json",
		Label:            "feeds.json",
		PublicPath:       feedHealthPublicPath("feeds_json"),
		Status:           status,
		StatusText:       statusText,
		CurrentSignal:    signal,
		WhatThisMeans:    "The discovery document can tell operators and consumers where the feed URLs and metadata are expected to be.",
		Freshness:        freshness,
		ValidatorContext: "feeds.json is not itself a GTFS validator result; review each feed validator row separately.",
		HealthContext:    "Metadata presence is a readiness signal, not a public fetch or consumer-ingestion proof.",
		NextAction:       next,
		DoesNotProve:     "Does not prove consumer acceptance, listing, display, ingestion, final-root ownership, or CAL-ITP/Caltrans compliance.",
		AdminLinks:       []string{"/admin/operations/feeds", "/admin/operations/setup"},
		DocsLinks:        []string{"docs/requirements-calitp-compliance.md", "docs/tutorials/calitp-readiness-checklist.md"},
	}
}

func buildFeedHealthRow(page operationsPage, feedType string, label string) operationsFeedHealthRow {
	feed, hasFeed := feedHealthMetadata(page, feedType)
	validation := feedHealthValidationRow(page, feedType)
	reliability := feedHealthReliabilityRow(page, feedType)
	status := feedHealthStatus(page, feedType, feed, hasFeed, validation, reliability)
	return operationsFeedHealthRow{
		ID:               feedType,
		Label:            label,
		PublicPath:       feedHealthPublicPath(feedType),
		Status:           status,
		StatusText:       feedHealthStatusText(status),
		CurrentSignal:    feedHealthCurrentSignal(page, feedType, feed, hasFeed),
		WhatThisMeans:    feedHealthMeaning(feedType),
		Freshness:        feedHealthFreshness(feed, hasFeed, reliability),
		ValidatorContext: feedHealthValidatorContext(validation),
		HealthContext:    feedHealthHealthContext(page, feedType, feed, hasFeed, reliability),
		NextAction:       feedHealthNextAction(page, feedType, feed, hasFeed, validation, reliability, status),
		DoesNotProve:     feedHealthDoesNotProve(feedType),
		AdminLinks:       feedHealthAdminLinks(feedType),
		DocsLinks:        feedHealthDocsLinks(feedType),
	}
}

func feedHealthPublicPath(feedType string) string {
	switch feedType {
	case "feeds_json":
		return "/public/feeds.json"
	case "schedule":
		return "/public/gtfs/schedule.zip"
	case "vehicle_positions":
		return "/public/gtfsrt/vehicle_positions.pb"
	case "trip_updates":
		return "/public/gtfsrt/trip_updates.pb"
	case "alerts":
		return "/public/gtfsrt/alerts.pb"
	default:
		return ""
	}
}

func feedHealthMetadata(page operationsPage, feedType string) (compliance.FeedMetadata, bool) {
	for _, feed := range page.Discovery.Feeds {
		if feed.FeedType == feedType {
			return feed, true
		}
	}
	return compliance.FeedMetadata{}, false
}

func feedHealthValidationRow(page operationsPage, feedType string) *compliance.ValidationHealthRow {
	for i := range page.ValidationHealth.Feeds {
		if page.ValidationHealth.Feeds[i].FeedType == feedType {
			return &page.ValidationHealth.Feeds[i]
		}
	}
	return nil
}

func feedHealthReliabilityRow(page operationsPage, feedType string) *compliance.ReliabilityFeedRow {
	for i := range page.Reliability.Feeds {
		if page.Reliability.Feeds[i].FeedType == feedType {
			return &page.Reliability.Feeds[i]
		}
	}
	return nil
}

func feedHealthStatus(page operationsPage, feedType string, feed compliance.FeedMetadata, hasFeed bool, validation *compliance.ValidationHealthRow, reliability *compliance.ReliabilityFeedRow) string {
	if page.DiscoveryError != "" || !hasFeed || strings.TrimSpace(feed.CanonicalPublicURL) == "" {
		return checklistStatusMissing
	}
	if reliability != nil && reliability.Status == compliance.ReliabilityStatusUnhealthy {
		return checklistStatusBlocked
	}
	if validation != nil {
		switch validation.HealthStatus {
		case compliance.ValidationHealthStatusFailed, compliance.ValidationHealthStatusBlocked, compliance.ValidationHealthStatusMissingTooling, compliance.ValidationHealthStatusMisconfiguredTooling:
			return checklistStatusBlocked
		case compliance.ValidationHealthStatusArtifactUnavailable, compliance.ValidationHealthStatusNotRun, compliance.ValidationHealthStatusStale, compliance.ValidationHealthStatusNeedsReview:
			return checklistStatusNeedsReview
		}
	}
	if reliability != nil {
		switch reliability.Status {
		case compliance.ReliabilityStatusNeedsReview:
			return checklistStatusNeedsReview
		case compliance.ReliabilityStatusMissing:
			if feedType != "schedule" {
				return checklistStatusNeedsReview
			}
		case compliance.ReliabilityStatusUnknown:
			if feedType != "schedule" {
				return checklistStatusNeedsReview
			}
		}
	}
	if feed.LastValidationStatus == "" || feed.LastValidationStatus == "not_run" || feed.LastValidationStatus == "warning" {
		return checklistStatusNeedsReview
	}
	if feed.LastValidationStatus == "failed" || feed.LastValidationStatus == "blocked" {
		return checklistStatusBlocked
	}
	return checklistStatusOK
}

func feedHealthStatusText(status string) string {
	switch normalizeChecklistStatus(status) {
	case checklistStatusOK:
		return "Looks ready to review"
	case checklistStatusNeedsReview:
		return "Needs review"
	case checklistStatusMissing:
		return "Missing"
	case checklistStatusBlocked:
		return "Blocked"
	default:
		return "Unknown"
	}
}

func feedHealthCurrentSignal(page operationsPage, feedType string, feed compliance.FeedMetadata, hasFeed bool) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	if !hasFeed {
		return "no " + feedType + " metadata record is listed in feed discovery"
	}
	return fmt.Sprintf("URL=%q, active feed version=%q, validation=%s", feed.CanonicalPublicURL, feed.ActiveFeedVersionID, firstNonEmpty(feed.LastValidationStatus, "not_run"))
}

func feedHealthMeaning(feedType string) string {
	switch feedType {
	case "schedule":
		return "Static GTFS is the schedule foundation used by route, trip, stop, and service-day matching."
	case "vehicle_positions":
		return "Vehicle Positions show the latest publishable vehicle locations from accepted telemetry and conservative assignment state."
	case "trip_updates":
		return "Trip Updates show prediction output from the configured adapter boundary when enough defensible assignment data exists."
	case "alerts":
		return "Alerts show current service-alert records from the private Alerts Console lifecycle."
	default:
		return "This row summarizes a public transit feed output."
	}
}

func feedHealthFreshness(feed compliance.FeedMetadata, hasFeed bool, reliability *compliance.ReliabilityFeedRow) string {
	parts := []string{}
	if hasFeed {
		if feed.RevisionTimestamp != nil {
			parts = append(parts, "metadata revised "+formatTimeForText(feed.RevisionTimestamp))
		}
		if feed.LastValidationAt != nil {
			parts = append(parts, "validated "+formatTimeForText(feed.LastValidationAt))
		}
		if feed.LastHealthAt != nil {
			parts = append(parts, "health checked "+formatTimeForText(feed.LastHealthAt))
		}
	}
	if reliability != nil && reliability.SnapshotAt != nil {
		parts = append(parts, "snapshot "+formatTimeForText(reliability.SnapshotAt))
	}
	if len(parts) == 0 {
		return "no metadata revision, validation, or health timestamp recorded"
	}
	return strings.Join(parts, "; ")
}

func feedHealthValidatorContext(row *compliance.ValidationHealthRow) string {
	if row == nil {
		return "no validator health row is available"
	}
	return fmt.Sprintf("%s via %s: latest=%s, stale=%s, tooling=%s, artifact=%s", row.ValidatorID, row.ValidatorName, row.LatestResultStatus, row.StaleStatus, row.ToolingStatus, row.ArtifactStatus)
}

func feedHealthHealthContext(page operationsPage, feedType string, feed compliance.FeedMetadata, hasFeed bool, reliability *compliance.ReliabilityFeedRow) string {
	parts := []string{}
	if hasFeed {
		if feed.LastHealthStatus != "" {
			parts = append(parts, "feed discovery health="+feed.LastHealthStatus)
		}
	}
	if reliability != nil {
		parts = append(parts, fmt.Sprintf("private reliability status=%s; threshold=%s", reliability.Status, feedHealthThresholdText(reliability.DiagnosticThreshold)))
	}
	if feedType == "trip_updates" {
		if page.TripUpdatesQuality.Recorded {
			parts = append(parts, "Trip Updates diagnostics="+page.TripUpdatesQuality.DiagnosticsStatus+"/"+page.TripUpdatesQuality.DiagnosticsReason)
		} else {
			parts = append(parts, page.TripUpdatesQuality.Message)
		}
	}
	if len(parts) == 0 {
		return "no private feed-health snapshot is recorded"
	}
	return strings.Join(parts, "; ")
}

func feedHealthThresholdText(threshold string) string {
	clean := strings.TrimSpace(threshold)
	if clean == "" {
		return "not recorded"
	}
	lower := strings.ToLower(clean)
	if strings.Contains(lower, "sla") && strings.Contains(lower, "uptime") {
		return "private diagnostic threshold; not a service-level or uptime proof"
	}
	return clean
}

func feedHealthNextAction(page operationsPage, feedType string, feed compliance.FeedMetadata, hasFeed bool, validation *compliance.ValidationHealthRow, reliability *compliance.ReliabilityFeedRow, status string) string {
	if page.DiscoveryError != "" || !hasFeed || strings.TrimSpace(feed.CanonicalPublicURL) == "" {
		if feedType == "schedule" {
			return "Import GTFS by browser or CLI, then store publication metadata and rerun validation."
		}
		return "Confirm the feed is configured and listed in feeds.json, then run validator health."
	}
	if validation != nil && validation.HealthStatus != compliance.ValidationHealthStatusRecorded {
		return validation.NextAction
	}
	if reliability != nil && reliability.Status != compliance.ReliabilityStatusOK {
		return reliability.NextAction
	}
	if status == checklistStatusOK {
		return "Continue periodic private validation and freshness review; keep stronger public claims out until evidence exists."
	}
	return "Review validator and feed-health context before relying on this feed."
}

func feedHealthDoesNotProve(feedType string) string {
	switch feedType {
	case "schedule":
		return "Does not prove validator-clean production data, agency approval, final-root proof, consumer acceptance, or CAL-ITP/Caltrans compliance."
	case "vehicle_positions":
		return "Does not prove production AVL reliability, vendor compatibility, hardware certification, consumer display, or compliance."
	case "trip_updates":
		return "Does not prove production-grade ETA quality, real-world ETA accuracy, consumer display, or compliance."
	case "alerts":
		return "Does not prove consumer display, agency approval, public launch completion, or compliance."
	default:
		return "Does not prove consumer acceptance, public launch, SLA, uptime proof, production readiness, or compliance."
	}
}

func feedHealthAdminLinks(feedType string) []string {
	switch feedType {
	case "schedule":
		return []string{"/admin/operations/gtfs-import", "/admin/operations/gtfs-quality", "/admin/operations/validation-health"}
	case "vehicle_positions":
		return []string{"/admin/operations/telemetry", "/admin/operations/validation-health", "/admin/operations/reliability"}
	case "trip_updates":
		return []string{"/admin/operations/feeds", "/admin/operations/validation-health", "/admin/operations/reliability"}
	case "alerts":
		return []string{"/admin/operations/feeds", "/admin/operations/validation-health", "/admin/operations/reliability"}
	default:
		return []string{"/admin/operations/feeds"}
	}
}

func feedHealthDocsLinks(feedType string) []string {
	switch feedType {
	case "schedule":
		return []string{"docs/tutorials/real-agency-gtfs-onboarding.md", "docs/tutorials/gtfs-validation-triage.md"}
	case "vehicle_positions":
		return []string{"docs/requirements-2a-2f.md", "docs/tutorials/device-avl-integration.md"}
	case "trip_updates":
		return []string{"docs/requirements-trip-updates.md", "docs/tutorials/operator-smoke-and-support-bundle.md"}
	case "alerts":
		return []string{"docs/requirements-calitp-compliance.md"}
	default:
		return []string{"docs/requirements-calitp-compliance.md"}
	}
}

func feedHealthCounts(rows []operationsFeedHealthRow) operationsFeedHealthCounts {
	counts := operationsFeedHealthCounts{Rows: len(rows), Statuses: map[string]int{
		checklistStatusOK:          0,
		checklistStatusNeedsReview: 0,
		checklistStatusMissing:     0,
		checklistStatusBlocked:     0,
		checklistStatusUnknown:     0,
	}}
	for _, row := range rows {
		counts.Statuses[normalizeChecklistStatus(row.Status)]++
	}
	return counts
}
