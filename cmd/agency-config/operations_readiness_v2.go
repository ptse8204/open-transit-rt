package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
)

type operationsReadinessV2View struct {
	GeneratedAt time.Time                    `json:"generated_at"`
	AgencyID    string                       `json:"agency_id"`
	Boundary    string                       `json:"boundary"`
	FocusAreas  []operationsReadinessV2Focus `json:"focus_areas"`
	Rows        []operationsReadinessV2Row   `json:"rows"`
	Counts      operationsReadinessV2Counts  `json:"counts"`
	ClaimFlags  operationsReadinessV2Claims  `json:"claim_flags"`
}

type operationsReadinessV2Focus struct {
	ID                 string   `json:"id"`
	Label              string   `json:"label"`
	Status             string   `json:"status"`
	WhatThisHelpsWith  string   `json:"what_this_helps_with"`
	PrimarySignal      string   `json:"primary_signal"`
	NextAction         string   `json:"next_action"`
	WhatItDoesNotProve string   `json:"what_it_does_not_prove"`
	RowIDs             []string `json:"row_ids"`
	AdminLinks         []string `json:"admin_links"`
	DocsLinks          []string `json:"docs_links"`
}

type operationsReadinessV2Row struct {
	ID                 string   `json:"id"`
	ReadinessItem      string   `json:"readiness_item"`
	Status             string   `json:"status"`
	CurrentSignal      string   `json:"current_signal"`
	WhatThisMeans      string   `json:"what_this_means"`
	WhyItMatters       string   `json:"why_it_matters"`
	WhatToDoNext       string   `json:"what_to_do_next"`
	WhatItDoesNotProve string   `json:"what_it_does_not_prove"`
	AdminLinks         []string `json:"admin_links"`
	DocsLinks          []string `json:"docs_links"`
}

type operationsReadinessV2Counts struct {
	Rows     int            `json:"rows"`
	Statuses map[string]int `json:"statuses"`
}

type operationsReadinessV2Claims struct {
	ExternalEvidenceCreated    bool `json:"external_evidence_created"`
	FinalRootEvidenceCreated   bool `json:"final_root_evidence_created"`
	ConsumerStatusesChanged    bool `json:"consumer_statuses_changed"`
	ComplianceClaimed          bool `json:"compliance_claimed"`
	ProductionReadinessClaimed bool `json:"production_readiness_claimed"`
	AgencyApprovalClaimed      bool `json:"agency_approval_claimed"`
	ConsumerAcceptanceClaimed  bool `json:"consumer_acceptance_claimed"`
	PublicLaunchClaimed        bool `json:"public_launch_claimed"`
	HostedSaaSClaimed          bool `json:"hosted_saas_claimed"`
	SLAClaimed                 bool `json:"sla_claimed"`
	UptimeGuaranteeClaimed     bool `json:"uptime_guarantee_claimed"`
	VendorCompatibilityClaimed bool `json:"vendor_compatibility_claimed"`
	ProductionGradeETAClaimed  bool `json:"production_grade_eta_claimed"`
}

func (h *handler) renderReadinessV2(w http.ResponseWriter, r *http.Request) {
	principal, ok := authRequireOperationsRead(w, r)
	if !ok {
		return
	}
	page := h.buildOperationsPage(r, principal, "readiness")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "readiness", page)
}

func (h *handler) renderReadinessV2JSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := authRequireOperationsRead(w, r)
	if !ok {
		return
	}
	page := h.buildOperationsPage(r, principal, "readiness")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.ReadinessV2)
}

func authRequireOperationsRead(w http.ResponseWriter, r *http.Request) (auth.Principal, bool) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return auth.Principal{}, false
	}
	return principal, true
}

func buildOperationsReadinessV2(page operationsPage) operationsReadinessV2View {
	rows := []operationsReadinessV2Row{
		readinessV2Row(
			"discovery_metadata",
			"Feed discovery and metadata",
			readinessV2DiscoveryStatus(page),
			readinessV2DiscoverySignal(page),
			"Operators can see the configured agency name, public feed root, license, technical contact, and feeds.json readiness flags.",
			"Consumer and Caltrans-style workflows need stable, discoverable feed metadata before any external review can be meaningful.",
			readinessV2DiscoveryNext(page),
			"This private view does not show final-root ownership, agency approval, consumer discovery, public launch, or CAL-ITP/Caltrans compliance.",
			[]string{"/admin/operations/setup", "/admin/operations/feeds"},
			[]string{"docs/requirements-calitp-compliance.md", "docs/tutorials/calitp-readiness-checklist.md"},
		),
		readinessV2Row(
			"feed_health",
			"Plain-language feed health",
			readinessV2FeedHealthStatus(page),
			readinessV2FeedHealthSignal(page),
			"The private feed-health view summarizes feeds.json, schedule, Vehicle Positions, Trip Updates, and Alerts with validator and reliability context.",
			"Readiness work is easier to triage when operators can distinguish missing metadata, validator blockers, stale health, and weak realtime signals.",
			readinessV2FeedHealthNext(page),
			"This private view does not show service-level or uptime proof, consumer acceptance, public launch completion, production readiness, or compliance.",
			[]string{"/admin/operations/feed-health", "/admin/operations/validation-health", "/admin/operations/reliability"},
			[]string{"docs/tutorials/no-cli-agency-first-run.md", "docs/requirements-2a-2f.md"},
		),
		readinessV2Row(
			"static_gtfs_quality",
			"Static GTFS quality",
			readinessV2GTFSQualityStatus(page),
			readinessV2GTFSQualitySignal(page),
			"The active schedule feed and private GTFS quality triage show canonical validator and internal importer status separately.",
			"Static GTFS is the foundation for service-day matching, trip descriptors, feed discovery, and all GTFS-Realtime outputs.",
			readinessV2GTFSQualityNext(page),
			"This private view does not show validator-clean production data, no-warning compliance, final-root approval, or consumer ingestion.",
			[]string{"/admin/operations/gtfs-import", "/admin/operations/gtfs-quality", "/admin/gtfs-studio"},
			[]string{"docs/tutorials/real-agency-gtfs-onboarding.md", "docs/tutorials/gtfs-validation-triage.md"},
		),
		readinessV2Row(
			"vehicle_positions",
			"Vehicle Positions readiness",
			readinessV2FeedRowStatus(page, "vehicle_positions"),
			readinessV2FeedRowSignal(page, "vehicle_positions"),
			"Vehicle Positions readiness combines feed metadata, validator health, reliability health, telemetry freshness, and conservative assignment state.",
			"Vehicle Positions are the first production-grade public output target and the base signal for any later prediction adapter.",
			readinessV2RealtimeNext(page, "vehicle_positions"),
			"This private view does not show production AVL reliability, real vendor compatibility, hardware certification, consumer display, or compliance.",
			[]string{"/admin/operations/feed-health", "/admin/operations/telemetry-simulator", "/admin/operations/telemetry", "/admin/operations/devices"},
			[]string{"docs/requirements-2a-2f.md", "docs/tutorials/device-avl-integration.md"},
		),
		readinessV2Row(
			"trip_updates",
			"Trip Updates adapter boundary",
			readinessV2FeedRowStatus(page, "trip_updates"),
			readinessV2TripUpdatesSignal(page),
			"Trip Updates readiness shows whether the pluggable prediction boundary has diagnostics, coverage, and validation signals.",
			"Trip Updates are required for complete realtime publication, but they must stay replaceable and must not be confused with Vehicle Positions publishing.",
			readinessV2RealtimeNext(page, "trip_updates"),
			"This private view does not show production-grade ETA quality, real-world ETA accuracy, consumer display, or compliance.",
			[]string{"/admin/operations/feed-health", "/admin/operations/feeds", "/admin/operations/reliability"},
			[]string{"docs/requirements-trip-updates.md", "docs/architecture.md"},
		),
		readinessV2Row(
			"alerts",
			"Service Alerts readiness",
			readinessV2FeedRowStatus(page, "alerts"),
			readinessV2FeedRowSignal(page, "alerts"),
			"Alerts readiness shows whether the Alerts feed is listed, validator-visible, and represented in feed-health context.",
			"A complete GTFS-Realtime set needs Alerts alongside Vehicle Positions and Trip Updates so disruptions can be communicated consistently.",
			readinessV2RealtimeNext(page, "alerts"),
			"This private view does not show agency-reviewed alert content, public launch completion, consumer display, or compliance.",
			[]string{"/admin/operations/feed-health", "/admin/operations/feeds"},
			[]string{"docs/requirements-calitp-compliance.md"},
		),
		readinessV2Row(
			"validation_health",
			"Validator health",
			readinessV2NormalizeStatus(page.ValidationHealth.OverallStatus),
			readinessV2ValidationSignal(page),
			"Private validator health summarizes schedule and realtime validator tooling, artifacts, latest records, staleness, and next actions.",
			"Validator records help operators avoid publishing or relying on feeds with known schema, tooling, or artifact blockers.",
			readinessV2ValidationNext(page),
			"This private view does not show consumer acceptance, no-error production compliance, public fetch success, or agency approval.",
			[]string{"/admin/operations/validation-health", "/admin/operations/gtfs-quality"},
			[]string{"docs/tutorials/gtfs-validation-triage.md", "docs/compliance-evidence-checklist.md"},
		),
		readinessV2Row(
			"operations_reliability",
			"Operations reliability diagnostics",
			readinessV2NormalizeStatus(page.Reliability.OverallStatus),
			readinessV2ReliabilitySignal(page),
			"Private reliability diagnostics summarize recent feed-health records and incident rollups where available.",
			"Operators need freshness, invalid-response, availability, and incident context before escalating a feed from setup into monitored operations.",
			readinessV2ReliabilityNext(page),
			"This private view does not show service-level commitments, uptime proof, hosted-service availability, production readiness, or managed support.",
			[]string{"/admin/operations/reliability", "/admin/operations/feed-health"},
			[]string{"docs/runbooks/small-agency-pilot-operations.md", "docs/requirements-2a-2f.md"},
		),
		readinessV2Row(
			"telemetry_devices",
			"Telemetry and device setup",
			readinessV2TelemetryDeviceStatus(page),
			readinessV2TelemetryDeviceSignal(page),
			"Telemetry and device setup show whether device bindings exist and whether recent accepted observations are visible for the agency.",
			"Realtime feed quality depends on device identity, accepted telemetry, stale handling, and conservative assignment behavior.",
			readinessV2TelemetryDeviceNext(page),
			"This private view does not show real fleet reliability, vendor AVL compatibility, hardware certification, or production AVL coverage.",
			[]string{"/admin/operations/devices", "/admin/operations/telemetry-simulator", "/admin/operations/telemetry"},
			[]string{"docs/tutorials/telemetry-simulator-and-device-trial.md", "docs/tutorials/agency-first-run.md", "docs/tutorials/operator-smoke-and-support-bundle.md"},
		),
		readinessV2Row(
			"operations_scorecard",
			"Operations scorecard",
			readinessV2ScorecardStatus(page),
			readinessV2ScorecardSignal(page),
			"The operations scorecard is the existing private rollup for deployment-health review and operator workflow status.",
			"A scorecard gives operators a repeatable review point before preparing support bundles, launch reviews, or issue triage.",
			readinessV2ScorecardNext(page),
			"This private view does not show agency adoption, managed hosting, service-level commitments, production readiness, public launch, or compliance.",
			[]string{"/admin/operations/evidence", "/admin/operations/reliability"},
			[]string{"docs/runbooks/small-agency-pilot-operations.md", "docs/tutorials/operator-smoke-and-support-bundle.md"},
		),
		readinessV2Row(
			"consumer_prepared_tracker",
			"Consumer prepared tracker",
			readinessV2ConsumerStatus(page),
			readinessV2ConsumerSignal(page),
			"The consumer tracker exposes prepared packet records and runtime workflow records without changing target-specific statuses.",
			"Major planner ingestion is an external outcome, so operators must keep prepared materials separate from submitted, reviewed, accepted, listed, displayed, or ingested claims.",
			readinessV2ConsumerNext(page),
			"This private view does not show submission, review, acceptance, listing, display, ingestion, consumer approval, or Caltrans/CAL-ITP compliance.",
			[]string{"/admin/operations/consumers"},
			[]string{"docs/requirements-calitp-compliance.md", "docs/compliance-evidence-checklist.md"},
		),
	}
	return operationsReadinessV2View{
		GeneratedAt: page.GeneratedAt,
		AgencyID:    page.AgencyID,
		Boundary:    "Private authenticated readiness checklist only; viewing it creates no evidence, changes no consumer status, contacts no external party, opens no public route, and records no approval, compliance, public launch, hosted-service, vendor, SLA, uptime, production-readiness, consumer-acceptance, or production-grade ETA outcome.",
		FocusAreas:  readinessV2FocusAreas(page),
		Rows:        rows,
		Counts:      readinessV2Counts(rows),
		ClaimFlags:  operationsReadinessV2Claims{},
	}
}

func readinessV2FocusAreas(page operationsPage) []operationsReadinessV2Focus {
	return []operationsReadinessV2Focus{
		readinessV2Focus(
			"public_feed_urls",
			"Public feed URLs",
			readinessV2PublicFeedURLStatus(page),
			"Helps operators confirm feeds.json, static GTFS, Vehicle Positions, Trip Updates, and Alerts have configured URLs before sharing them.",
			readinessV2PublicFeedURLSignal(page),
			readinessV2PublicFeedURLNext(page),
			"This private view does not show final-root ownership, source-of-truth listing, consumer ingestion, public launch, compliance, or production readiness.",
			[]string{"discovery_metadata", "feed_health"},
			[]string{"/admin/operations/feeds", "/admin/operations/feed-health"},
			[]string{"docs/release-candidate-readiness.md", "docs/requirements-calitp-compliance.md"},
		),
		readinessV2Focus(
			"static_gtfs",
			"Static GTFS",
			readinessV2GTFSQualityStatus(page),
			"Helps operators review schedule import, active feed version, required files, service dates, and validation triage.",
			readinessV2GTFSQualitySignal(page),
			readinessV2GTFSQualityNext(page),
			"This private view does not show validator-clean production data, agency approval, consumer acceptance, compliance, or final-root readiness.",
			[]string{"static_gtfs_quality"},
			[]string{"/admin/operations/gtfs-workbench", "/admin/operations/gtfs-quality", "/admin/operations/validation-health"},
			[]string{"docs/tutorials/gtfs-validation-triage.md", "docs/requirements-2a-2f.md"},
		),
		readinessV2Focus(
			"vehicle_positions",
			"Vehicle Positions",
			readinessV2FeedRowStatus(page, "vehicle_positions"),
			"Helps operators review the first public GTFS-Realtime target with telemetry freshness, assignment confidence, and feed-health context.",
			readinessV2FeedRowSignal(page, "vehicle_positions"),
			readinessV2RealtimeNext(page, "vehicle_positions"),
			"This private view does not show production AVL reliability, vendor compatibility, hardware certification, consumer display, compliance, or production readiness.",
			[]string{"vehicle_positions", "telemetry_devices"},
			[]string{"/admin/operations/realtime", "/admin/operations/telemetry", "/admin/operations/devices"},
			[]string{"docs/requirements-2a-2f.md", "docs/tutorials/device-avl-integration.md"},
		),
		readinessV2Focus(
			"trip_updates",
			"Trip Updates",
			readinessV2FeedRowStatus(page, "trip_updates"),
			"Helps operators keep the prediction boundary visible while reviewing generated, withheld, stale, ambiguous, and low-confidence Trip Updates behavior.",
			readinessV2TripUpdatesSignal(page),
			readinessV2RealtimeNext(page, "trip_updates"),
			"This private view does not show production-grade ETA quality, real-world ETA accuracy, consumer display, compliance, or production readiness.",
			[]string{"trip_updates"},
			[]string{"/admin/operations/realtime", "/admin/operations/prediction-lab", "/admin/operations/feed-health"},
			[]string{"docs/requirements-trip-updates.md", "docs/tutorials/prediction-eta-lab.md"},
		),
		readinessV2Focus(
			"alerts",
			"Alerts",
			readinessV2FeedRowStatus(page, "alerts"),
			"Helps operators verify the Alerts feed URL, feed-health row, and service-disruption review path.",
			readinessV2FeedRowSignal(page, "alerts"),
			readinessV2RealtimeNext(page, "alerts"),
			"This private view does not show agency-reviewed alert content, consumer display, compliance, public launch, or production readiness.",
			[]string{"alerts"},
			[]string{"/admin/operations/realtime", "/admin/alerts/console", "/admin/operations/feed-health"},
			[]string{"docs/requirements-calitp-compliance.md"},
		),
		readinessV2Focus(
			"validation",
			"Validation",
			readinessV2NormalizeStatus(page.ValidationHealth.OverallStatus),
			"Helps operators see static and realtime validator tooling, latest records, stale checks, blockers, and next actions.",
			readinessV2ValidationSignal(page),
			readinessV2ValidationNext(page),
			"This private view does not show validator-clean public feeds, compliance, consumer acceptance, public fetch success, or agency approval.",
			[]string{"validation_health", "static_gtfs_quality"},
			[]string{"/admin/operations/validation-health", "/admin/operations/validation-center", "/admin/operations/gtfs-quality"},
			[]string{"docs/tutorials/gtfs-validation-triage.md", "docs/dependencies.md"},
		),
		readinessV2Focus(
			"license_contact",
			"License and contact metadata",
			readinessV2LicenseContactStatus(page),
			"Helps operators verify open-license and technical-contact fields before any future public feed review.",
			readinessV2LicenseContactSignal(page),
			readinessV2LicenseContactNext(page),
			"This private view does not show legal approval, managed support, consumer acceptance, compliance, or source-of-truth listing.",
			[]string{"discovery_metadata"},
			[]string{"/admin/operations/setup", "/admin/operations/feeds"},
			[]string{"docs/requirements-calitp-compliance.md", "docs/release-candidate-readiness.md"},
		),
		readinessV2Focus(
			"uptime_operations",
			"Uptime and operations signals",
			readinessV2NormalizeStatus(page.Reliability.OverallStatus),
			"Helps operators review local feed-health, reliability, maintenance, and incident signals before routine operations.",
			readinessV2ReliabilitySignal(page),
			readinessV2ReliabilityNext(page),
			"This private view does not show uptime, SLA coverage, hosted service availability, managed support, compliance, or production readiness.",
			[]string{"operations_reliability", "operations_scorecard"},
			[]string{"/admin/operations/reliability", "/admin/operations/maintenance", "/admin/operations/feed-health"},
			[]string{"docs/runbooks/small-agency-pilot-operations.md", "docs/runbooks/monitoring-and-alerting.md"},
		),
		readinessV2Focus(
			"telemetry_device_state",
			"Telemetry and device state",
			readinessV2TelemetryDeviceStatus(page),
			"Helps operators review device bindings, accepted observations, stale handling, and conservative assignment behavior.",
			readinessV2TelemetryDeviceSignal(page),
			readinessV2TelemetryDeviceNext(page),
			"This private view does not show real fleet reliability, vendor AVL compatibility, hardware certification, production AVL coverage, compliance, or consumer acceptance.",
			[]string{"telemetry_devices", "vehicle_positions"},
			[]string{"/admin/operations/devices", "/admin/operations/telemetry", "/admin/operations/telemetry-simulator"},
			[]string{"docs/tutorials/telemetry-simulator-and-device-trial.md", "docs/connectors/catalog.md"},
		),
		readinessV2Focus(
			"consumer_preparedness",
			"Consumer preparedness",
			readinessV2ConsumerStatus(page),
			"Helps operators keep prepared packet records, public feed metadata, and external-target boundaries visible without moving statuses.",
			readinessV2ConsumerSignal(page),
			readinessV2ConsumerNext(page),
			"This private view does not show submission, review, acceptance, ingestion, listing, display, compliance, consumer approval, or public launch.",
			[]string{"consumer_prepared_tracker", "discovery_metadata"},
			[]string{"/admin/operations/consumers", "/admin/operations/readiness"},
			[]string{"docs/consumer-submission-evidence.md", "docs/evidence/consumer-submissions/README.md"},
		),
	}
}

func readinessV2Focus(id, label, status, helpsWith, signal, next, boundary string, rowIDs, adminLinks, docsLinks []string) operationsReadinessV2Focus {
	return operationsReadinessV2Focus{
		ID:                 firstNonEmpty(id, "readiness_focus"),
		Label:              firstNonEmpty(label, "Readiness focus"),
		Status:             readinessV2NormalizeStatus(status),
		WhatThisHelpsWith:  firstNonEmpty(helpsWith, "Helps operators review local readiness signals."),
		PrimarySignal:      firstNonEmpty(signal, "unknown"),
		NextAction:         firstNonEmpty(next, "Review the linked private Operations Console page."),
		WhatItDoesNotProve: firstNonEmpty(boundary, privateBoundary()),
		RowIDs:             cleanLaunchpadList(rowIDs),
		AdminLinks:         safeAdminLinks(adminLinks),
		DocsLinks:          safeDocsLinks(docsLinks),
	}
}

func readinessV2Row(id, item, status, signal, means, matters, next, boundary string, adminLinks, docsLinks []string) operationsReadinessV2Row {
	return operationsReadinessV2Row{
		ID:                 id,
		ReadinessItem:      firstNonEmpty(item, id),
		Status:             readinessV2NormalizeStatus(status),
		CurrentSignal:      firstNonEmpty(signal, "unknown"),
		WhatThisMeans:      firstNonEmpty(means, "This row summarizes an existing private operations signal."),
		WhyItMatters:       firstNonEmpty(matters, "Operators need this signal before escalating readiness decisions."),
		WhatToDoNext:       firstNonEmpty(next, "Review the linked private Operations Console page."),
		WhatItDoesNotProve: firstNonEmpty(boundary, privateBoundary()),
		AdminLinks:         safeAdminLinks(adminLinks),
		DocsLinks:          safeDocsLinks(docsLinks),
	}
}

func readinessV2Counts(rows []operationsReadinessV2Row) operationsReadinessV2Counts {
	counts := operationsReadinessV2Counts{Rows: len(rows), Statuses: map[string]int{
		checklistStatusOK:          0,
		checklistStatusNeedsReview: 0,
		checklistStatusMissing:     0,
		checklistStatusBlocked:     0,
		checklistStatusUnknown:     0,
	}}
	for _, row := range rows {
		counts.Statuses[readinessV2NormalizeStatus(row.Status)]++
	}
	return counts
}

func readinessV2NormalizeStatus(status string) string {
	switch status {
	case checklistStatusOK, "recorded", "passed":
		return checklistStatusOK
	case compliance.GTFSQualityInformational:
		return checklistStatusOK
	case checklistStatusNeedsReview, "warning", "stale", "not_run":
		return checklistStatusNeedsReview
	case checklistStatusMissing, "artifact_unavailable":
		return checklistStatusMissing
	case checklistStatusBlocked, "failed", "unhealthy", "missing_tooling", "misconfigured_tooling":
		return checklistStatusBlocked
	case compliance.GTFSQualityBlocking:
		return checklistStatusBlocked
	case checklistStatusUnknown, "":
		return checklistStatusUnknown
	default:
		return normalizeChecklistStatus(status)
	}
}

func readinessV2DiscoveryStatus(page operationsPage) string {
	if page.DiscoveryError != "" {
		return checklistStatusMissing
	}
	if page.Discovery.Readiness.AllRequiredFeedsListed && page.Discovery.Readiness.HTTPSURLs && page.Discovery.Readiness.Discoverable && page.Discovery.Readiness.LicenseComplete && page.Discovery.Readiness.ContactComplete {
		return checklistStatusOK
	}
	return checklistStatusNeedsReview
}

func readinessV2DiscoverySignal(page operationsPage) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	parts := []string{
		fmt.Sprintf("%d feed records are listed", len(page.Discovery.Feeds)),
		readinessV2BoolPhrase(page.Discovery.Readiness.AllRequiredFeedsListed, "all required feeds are listed", "one or more required feeds are missing"),
		readinessV2BoolPhrase(page.Discovery.Readiness.HTTPSURLs, "public feed URLs use HTTPS", "one or more public feed URLs are not HTTPS"),
		readinessV2BoolPhrase(page.Discovery.Readiness.Discoverable, "feeds.json is discoverable", "feeds.json discoverability needs review"),
		readinessV2BoolPhrase(page.Discovery.Readiness.LicenseComplete, "license metadata is present", "license metadata is incomplete"),
		readinessV2BoolPhrase(page.Discovery.Readiness.ContactComplete, "technical contact is present", "technical contact is missing"),
	}
	if strings.TrimSpace(page.Discovery.AgencyName) != "" {
		parts = append(parts, "agency name: "+page.Discovery.AgencyName)
	}
	if strings.TrimSpace(page.Discovery.PublicBaseURL) != "" {
		parts = append(parts, "public feed root: "+page.Discovery.PublicBaseURL)
	}
	return strings.Join(parts, "; ")
}

func readinessV2DiscoveryNext(page operationsPage) string {
	if page.DiscoveryError != "" {
		return "Bootstrap publication metadata after importing or publishing a schedule feed."
	}
	if !page.Discovery.Readiness.AllRequiredFeedsListed {
		return "Add or repair feed metadata so schedule, Vehicle Positions, Trip Updates, and Alerts are all listed."
	}
	if !page.Discovery.Readiness.HTTPSURLs || !page.Discovery.Readiness.Discoverable {
		return "Review public base and feed URLs, then keep final-root conclusions separate until retained proof exists."
	}
	if !page.Discovery.Readiness.LicenseComplete || !page.Discovery.Readiness.ContactComplete {
		return "Enter operator-confirmed open license and monitored technical contact metadata."
	}
	return "Keep metadata current and review feed-specific health rows before stronger workflow decisions."
}

func readinessV2PublicFeedURLStatus(page operationsPage) string {
	if page.DiscoveryError != "" || strings.TrimSpace(page.Discovery.PublicBaseURL) == "" {
		return checklistStatusMissing
	}
	if page.Discovery.Readiness.AllRequiredFeedsListed && page.Discovery.Readiness.HTTPSURLs && page.Discovery.Readiness.Discoverable {
		return checklistStatusOK
	}
	return checklistStatusNeedsReview
}

func readinessV2PublicFeedURLSignal(page operationsPage) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	return strings.Join([]string{
		fmt.Sprintf("public feed root: %s", firstNonEmpty(page.Discovery.PublicBaseURL, "missing")),
		fmt.Sprintf("%d feed records are listed", len(page.Discovery.Feeds)),
		readinessV2BoolPhrase(page.Discovery.Readiness.AllRequiredFeedsListed, "all required feeds are listed", "one or more required feeds are missing"),
		readinessV2BoolPhrase(page.Discovery.Readiness.HTTPSURLs, "public feed URLs use HTTPS", "one or more public feed URLs are not HTTPS"),
		readinessV2BoolPhrase(page.Discovery.Readiness.Discoverable, "feeds.json is discoverable", "feeds.json discoverability needs review"),
	}, "; ")
}

func readinessV2PublicFeedURLNext(page operationsPage) string {
	if page.DiscoveryError != "" || strings.TrimSpace(page.Discovery.PublicBaseURL) == "" {
		return "Set publication metadata so feeds.json and the four feed URLs can be reviewed in the browser."
	}
	if !page.Discovery.Readiness.AllRequiredFeedsListed {
		return "List schedule, Vehicle Positions, Trip Updates, and Alerts in feeds.json before external review."
	}
	if !page.Discovery.Readiness.HTTPSURLs || !page.Discovery.Readiness.Discoverable {
		return "Review public base and feed URLs, then keep final-root conclusions separate until retained proof exists."
	}
	return "Open Feed URLs and Feed Health to review validation and freshness context before sharing configured URLs."
}

func readinessV2LicenseContactStatus(page operationsPage) string {
	if page.DiscoveryError != "" {
		return checklistStatusMissing
	}
	if page.Discovery.Readiness.LicenseComplete && page.Discovery.Readiness.ContactComplete {
		return checklistStatusOK
	}
	return checklistStatusNeedsReview
}

func readinessV2LicenseContactSignal(page operationsPage) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	parts := []string{
		"agency name: " + firstNonEmpty(page.Discovery.AgencyName, "missing"),
		readinessV2BoolPhrase(page.Discovery.Readiness.LicenseComplete, "license metadata is present", "license metadata is incomplete"),
		readinessV2BoolPhrase(page.Discovery.Readiness.ContactComplete, "technical contact is present", "technical contact is missing"),
	}
	if strings.TrimSpace(page.Discovery.License.URL) != "" {
		parts = append(parts, "license URL is configured")
	}
	return strings.Join(parts, "; ")
}

func readinessV2LicenseContactNext(page operationsPage) string {
	if page.DiscoveryError != "" {
		return "Bootstrap publication metadata after importing or publishing a schedule feed."
	}
	if !page.Discovery.Readiness.LicenseComplete {
		return "Add an operator-reviewed open-license name and URL before external feed review."
	}
	if !page.Discovery.Readiness.ContactComplete {
		return "Add a monitored technical contact before external feed review."
	}
	return "Keep license and technical contact metadata current with operator review."
}

func readinessV2FeedHealthStatus(page operationsPage) string {
	if len(page.FeedHealth.Rows) == 0 {
		return checklistStatusMissing
	}
	status := checklistStatusOK
	for _, row := range page.FeedHealth.Rows {
		switch row.Status {
		case checklistStatusBlocked:
			return checklistStatusBlocked
		case checklistStatusMissing:
			if status != checklistStatusBlocked {
				status = checklistStatusMissing
			}
		case checklistStatusNeedsReview:
			if status == checklistStatusOK {
				status = checklistStatusNeedsReview
			}
		case checklistStatusUnknown:
			if status == checklistStatusOK {
				status = checklistStatusUnknown
			}
		}
	}
	return status
}

func readinessV2FeedHealthSignal(page operationsPage) string {
	if len(page.FeedHealth.Rows) == 0 {
		return "feed-health model is not available"
	}
	return fmt.Sprintf("%d feed-health rows reviewed: %s", len(page.FeedHealth.Rows), readinessV2StatusSentence(page.FeedHealth.Counts.Statuses))
}

func readinessV2FeedHealthNext(page operationsPage) string {
	for _, row := range page.FeedHealth.Rows {
		if row.Status != checklistStatusOK {
			return row.NextAction
		}
	}
	return "Continue periodic private feed-health, validator, and reliability review."
}

func readinessV2GTFSQualityStatus(page operationsPage) string {
	if strings.TrimSpace(page.ActiveFeedVersion) == "" {
		return checklistStatusMissing
	}
	return readinessV2WorstStatus(
		readinessV2NormalizeGTFSQualityStatus(page.GTFSQuality.Canonical.Status),
		readinessV2NormalizeGTFSQualityStatus(page.GTFSQuality.InternalImporter.Status),
	)
}

func readinessV2NormalizeGTFSQualityStatus(status string) string {
	switch status {
	case compliance.GTFSQualityInformational:
		return checklistStatusOK
	case compliance.GTFSQualityNeedsReview:
		return checklistStatusNeedsReview
	case compliance.GTFSQualityBlocking:
		return checklistStatusBlocked
	case compliance.GTFSQualityUnknown, "":
		return checklistStatusUnknown
	default:
		return readinessV2NormalizeStatus(status)
	}
}

func readinessV2GTFSQualitySignal(page operationsPage) string {
	if strings.TrimSpace(page.ActiveFeedVersion) == "" {
		return "no active schedule feed version is available"
	}
	return fmt.Sprintf("active schedule version %q; static validator is %s; internal importer is %s", page.ActiveFeedVersion, readinessV2GTFSQualityStatusText(page.GTFSQuality.Canonical.Status), readinessV2GTFSQualityStatusText(page.GTFSQuality.InternalImporter.Status))
}

func readinessV2GTFSQualityNext(page operationsPage) string {
	if strings.TrimSpace(page.ActiveFeedVersion) == "" {
		return "Import or publish a schedule before rerunning static validation."
	}
	if page.GTFSQuality.Canonical.RecommendedAction != "" {
		return page.GTFSQuality.Canonical.RecommendedAction
	}
	return "Review GTFS quality triage and rerun the allowlisted static validator when needed."
}

func readinessV2FeedRowStatus(page operationsPage, feedType string) string {
	for _, row := range page.FeedHealth.Rows {
		if row.ID == feedType {
			return row.Status
		}
	}
	return checklistStatusMissing
}

func readinessV2FeedRowSignal(page operationsPage, feedType string) string {
	for _, row := range page.FeedHealth.Rows {
		if row.ID == feedType {
			return row.CurrentSignal + "; " + row.HealthContext
		}
	}
	return "no " + feedType + " feed-health row is available"
}

func readinessV2TripUpdatesSignal(page operationsPage) string {
	base := readinessV2FeedRowSignal(page, "trip_updates")
	if page.TripUpdatesQuality.Recorded {
		return base + fmt.Sprintf("; diagnostics=%s/%s; coverage=%s; future stop coverage=%s", page.TripUpdatesQuality.DiagnosticsStatus, page.TripUpdatesQuality.DiagnosticsReason, page.TripUpdatesQuality.TripUpdatesCoverageRate, page.TripUpdatesQuality.FutureStopCoverageRate)
	}
	return base + "; " + page.TripUpdatesQuality.Message
}

func readinessV2RealtimeNext(page operationsPage, feedType string) string {
	for _, row := range page.FeedHealth.Rows {
		if row.ID == feedType {
			return row.NextAction
		}
	}
	return "Confirm the feed is configured and listed, then run private validator and feed-health review."
}

func readinessV2ValidationSignal(page operationsPage) string {
	if len(page.ValidationHealth.Feeds) == 0 {
		return "no validator health rows are available"
	}
	return fmt.Sprintf("validator health is %s; tooling is %s; %d feed validator rows are available", firstNonEmpty(page.ValidationHealth.OverallStatus, checklistStatusUnknown), firstNonEmpty(page.ValidationHealth.ToolingStatus, checklistStatusUnknown), len(page.ValidationHealth.Feeds))
}

func readinessV2ValidationNext(page operationsPage) string {
	for _, row := range page.ValidationHealth.Feeds {
		if row.HealthStatus != "recorded" && row.NextAction != "" {
			return row.NextAction
		}
	}
	return "Keep validator health current and review feed-specific validator rows before external readiness decisions."
}

func readinessV2ReliabilitySignal(page operationsPage) string {
	if page.ReliabilityError != "" {
		return page.ReliabilityError
	}
	return fmt.Sprintf("reliability status is %s; %d feed rows and %d recent incidents are visible", firstNonEmpty(page.Reliability.OverallStatus, checklistStatusUnknown), len(page.Reliability.Feeds), len(page.Reliability.Incidents.Recent))
}

func readinessV2ReliabilityNext(page operationsPage) string {
	if page.ReliabilityError != "" {
		return "Run or configure private reliability feed-health snapshots before relying on reliability status."
	}
	for _, row := range page.Reliability.Feeds {
		if row.Status != checklistStatusOK && row.NextAction != "" {
			return row.NextAction
		}
	}
	return "Continue periodic private reliability and incident review."
}

func readinessV2TelemetryDeviceStatus(page operationsPage) string {
	if page.DeviceError != "" || page.TelemetryError != "" {
		return checklistStatusUnknown
	}
	if len(page.Devices) == 0 || len(page.Telemetry) == 0 {
		return checklistStatusMissing
	}
	if page.StaleCount > 0 {
		return checklistStatusNeedsReview
	}
	return checklistStatusOK
}

func readinessV2TelemetryDeviceSignal(page operationsPage) string {
	parts := []string{
		fmt.Sprintf("%d device bindings are configured", len(page.Devices)),
	}
	if page.DeviceError != "" {
		parts = append(parts, "device records are unavailable: "+page.DeviceError)
	}
	if page.TelemetryError != "" {
		parts = append(parts, "telemetry records are unavailable: "+page.TelemetryError)
	} else {
		staleText := "no stale vehicles"
		if page.StaleCount > 0 {
			staleText = fmt.Sprintf("%d stale vehicles above %s", page.StaleCount, page.StaleThreshold)
		}
		parts = append(parts, fmt.Sprintf("%d latest telemetry rows; %s", len(page.Telemetry), staleText))
	}
	return strings.Join(parts, "; ")
}

func readinessV2TelemetryDeviceNext(page operationsPage) string {
	if len(page.Devices) == 0 {
		return "Bind or rotate a device token and store the token outside this repo."
	}
	if len(page.Telemetry) == 0 {
		return "Send an accepted telemetry event through device credentials and review the latest telemetry page."
	}
	if page.StaleCount > 0 {
		return "Resolve stale telemetry or adjust device reporting before relying on realtime feed freshness."
	}
	return "Continue monitoring stale and unmatched vehicle behavior."
}

func readinessV2ScorecardStatus(page operationsPage) string {
	if page.ScorecardError != "" {
		return checklistStatusMissing
	}
	if page.Scorecard == nil {
		return checklistStatusUnknown
	}
	return checklistStatusNeedsReview
}

func readinessV2ScorecardSignal(page operationsPage) string {
	if page.ScorecardError != "" {
		return page.ScorecardError
	}
	if page.Scorecard == nil {
		return "no scorecard snapshot is available"
	}
	return fmt.Sprintf("scorecard status is %s; snapshot recorded at %s", page.Scorecard.OverallStatus, page.Scorecard.SnapshotAt.UTC().Format(time.RFC3339))
}

func readinessV2ScorecardNext(page operationsPage) string {
	if page.Scorecard == nil {
		return "Run the scorecard and supporting operator workflows in the private operator environment."
	}
	return "Review scorecard details alongside reliability, validation health, and support-bundle workflows."
}

func readinessV2ConsumerStatus(page operationsPage) string {
	if len(page.Consumers) == 0 {
		return checklistStatusMissing
	}
	return checklistStatusNeedsReview
}

func readinessV2ConsumerSignal(page operationsPage) string {
	if len(page.Consumers) == 0 {
		return "no prepared consumer tracker targets are visible"
	}
	signal := fmt.Sprintf("%d prepared tracker targets are visible", len(page.Consumers))
	if len(page.RuntimeConsumers) > 0 {
		signal += fmt.Sprintf("; %d runtime consumer workflow records are visible separately", len(page.RuntimeConsumers))
	}
	return signal
}

func readinessV2ConsumerNext(page operationsPage) string {
	if len(page.Consumers) == 0 {
		return "Review prepared packet docs and keep statuses unchanged until target-originated evidence exists."
	}
	return "Use prepared packet materials for operator review only; change target statuses only with retained target-originated records."
}

func readinessV2BoolPhrase(ok bool, yes string, no string) string {
	if ok {
		return yes
	}
	return no
}

func readinessV2GTFSQualityStatusText(status string) string {
	switch readinessV2NormalizeGTFSQualityStatus(status) {
	case checklistStatusOK:
		return "recorded with informational or no blocking issues"
	case checklistStatusNeedsReview:
		return "recorded with warnings to review"
	case checklistStatusBlocked:
		return "blocked by errors or missing run requirements"
	case checklistStatusMissing:
		return "missing"
	default:
		return "not observed yet"
	}
}

func readinessV2StatusSentence(counts map[string]int) string {
	parts := []string{}
	if counts[checklistStatusOK] > 0 {
		parts = append(parts, fmt.Sprintf("%d ok", counts[checklistStatusOK]))
	}
	if counts[checklistStatusNeedsReview] > 0 {
		parts = append(parts, fmt.Sprintf("%d need review", counts[checklistStatusNeedsReview]))
	}
	if counts[checklistStatusMissing] > 0 {
		parts = append(parts, fmt.Sprintf("%d missing", counts[checklistStatusMissing]))
	}
	if counts[checklistStatusBlocked] > 0 {
		parts = append(parts, fmt.Sprintf("%d blocked", counts[checklistStatusBlocked]))
	}
	if counts[checklistStatusUnknown] > 0 {
		parts = append(parts, fmt.Sprintf("%d unknown", counts[checklistStatusUnknown]))
	}
	if len(parts) == 0 {
		return "no rows are available"
	}
	return strings.Join(parts, ", ")
}

func readinessV2WorstStatus(statuses ...string) string {
	seen := map[string]bool{}
	for _, status := range statuses {
		seen[readinessV2NormalizeStatus(status)] = true
	}
	switch {
	case seen[checklistStatusBlocked]:
		return checklistStatusBlocked
	case seen[checklistStatusMissing]:
		return checklistStatusMissing
	case seen[checklistStatusNeedsReview]:
		return checklistStatusNeedsReview
	case seen[checklistStatusUnknown]:
		return checklistStatusUnknown
	default:
		return checklistStatusOK
	}
}
