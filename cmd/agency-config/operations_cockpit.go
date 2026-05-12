package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	operationsStatusReady          = "Ready"
	operationsStatusNeedsReview    = "Needs review"
	operationsStatusMissing        = "Missing"
	operationsStatusBlocked        = "Blocked"
	operationsStatusUnknown        = "Unknown"
	operationsStatusDiagnosticOnly = "Diagnostic only"
)

type operationsCockpitView struct {
	GeneratedAt   time.Time                  `json:"generated_at"`
	AgencyID      string                     `json:"agency_id"`
	Boundary      string                     `json:"boundary"`
	SetupProgress []operationsCockpitRow     `json:"setup_progress"`
	PrimaryCards  []operationsCockpitCard    `json:"primary_cards"`
	ClaimFlags    operationsCockpitClaimFlag `json:"claim_flags"`
}

type operationsCockpitRow struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	AdminLink     string `json:"admin_link"`
	DoesNotProve  string `json:"does_not_prove"`
}

type operationsCockpitCard struct {
	ID            string   `json:"id"`
	Label         string   `json:"label"`
	Status        string   `json:"status"`
	CurrentSignal string   `json:"current_signal"`
	NextAction    string   `json:"next_action"`
	AdminLink     string   `json:"admin_link"`
	DocsLinks     []string `json:"docs_links"`
	DoesNotProve  string   `json:"does_not_prove"`
}

type operationsCockpitClaimFlag struct {
	ExternalEvidenceCreated      bool `json:"external_evidence_created"`
	ConsumerStatusesChanged      bool `json:"consumer_statuses_changed"`
	ComplianceClaimed            bool `json:"compliance_claimed"`
	ProductionReadinessClaimed   bool `json:"production_readiness_claimed"`
	AgencyApprovalClaimed        bool `json:"agency_approval_claimed"`
	ConsumerAcceptanceClaimed    bool `json:"consumer_acceptance_claimed"`
	PublicLaunchClaimed          bool `json:"public_launch_claimed"`
	HostedSaaSClaimed            bool `json:"hosted_saas_claimed"`
	VendorCompatibilityClaimed   bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed bool `json:"hardware_certification_claimed"`
	SLAClaimed                   bool `json:"sla_claimed"`
	UptimeGuaranteeClaimed       bool `json:"uptime_guarantee_claimed"`
	ProductionGradeETAClaimed    bool `json:"production_grade_eta_claimed"`
}

func buildOperationsCockpit(page operationsPage) operationsCockpitView {
	progress := []operationsCockpitRow{
		cockpitProgressRow("agency_metadata", "Agency metadata", cockpitMetadataStatus(page), cockpitMetadataSignal(page), "Review agency name, public URL, license, and technical contact before sharing feed URLs.", "/admin/operations/setup", "Does not prove agency approval or final-root ownership."),
		cockpitProgressRow("gtfs_imported", "GTFS imported", cockpitActiveFeedStatus(page), cockpitActiveFeedSignal(page), "Import or review the active schedule, then open GTFS quality and feed health.", "/admin/operations/gtfs-import", "Does not prove schedule correctness or canonical validator success."),
		cockpitProgressRow("active_feed_version", "Active feed version", cockpitActiveFeedStatus(page), firstNonEmpty(page.ActiveFeedVersion, "no active schedule feed version recorded"), "Keep the active version visible before validator, rollback, or support decisions.", "/admin/operations/gtfs-import", "Does not prove staged rollback execution is available in the browser."),
		cockpitProgressRow("five_feed_paths", "Five public feed paths", cockpitFeedHealthStatus(page), cockpitFeedHealthSignal(page), "Open feed health and review each public path, validator state, health state, and next action.", "/admin/operations/feed-health", "Does not prove consumer acceptance, final-root ownership, or compliance."),
		cockpitProgressRow("validators", "Validators", cockpitValidationStatus(page), fmt.Sprintf("overall=%s; tooling=%s", page.ValidationHealth.OverallStatus, page.ValidationHealth.ToolingStatus), "Run or review allowlisted validator health from the private page.", "/admin/operations/validation-health", "Does not prove validator-clean production data or consumer acceptance."),
		cockpitProgressRow("telemetry", "Telemetry", cockpitTelemetryStatus(page), telemetryEvidence(page), "Create or review device bindings, then send a safe sample telemetry event if needed.", "/admin/operations/devices", "Does not prove real device reliability or vendor compatibility."),
		cockpitProgressRow("readiness", "Readiness", cockpitReadinessStatus(page), fmt.Sprintf("%d readiness rows; overall validation=%s", len(page.ReadinessV2.Rows), page.ValidationHealth.OverallStatus), "Open readiness for plain-language capability gaps and claim boundaries.", "/admin/operations/readiness", "Does not prove CAL-ITP/Caltrans compliance."),
		cockpitProgressRow("maintenance", "Maintenance", maintenanceStatusToCockpit(page.Maintenance.OverallStatus), page.Maintenance.OverallStatus, "Review maintenance tasks, backup/restore configuration state, telemetry freshness, and support summary instructions.", "/admin/operations/maintenance", "Does not prove SLA, uptime, hosted service availability, or production readiness."),
	}

	cards := []operationsCockpitCard{
		cockpitCard("import_update_gtfs", "Import / update GTFS", cockpitActiveFeedStatus(page), cockpitActiveFeedSignal(page), "Use browser upload or safe URL import. After import, review GTFS quality, validator health, and feed health.", "/admin/operations/gtfs-import", []string{"docs/tutorials/no-cli-agency-first-run.md"}, "Does not prove source schedule correctness, agency approval, or compliance."),
		cockpitCard("review_feed_health", "Review feed health", cockpitFeedHealthStatus(page), cockpitFeedHealthSignal(page), "Review exactly five public paths and the recorded validator/health state for each one.", "/admin/operations/feed-health", []string{"docs/tutorials/no-cli-agency-first-run.md", "docs/deployment/off-host-validation.md"}, "Does not prove consumer acceptance, ingestion, listing, display, final-root ownership, or compliance."),
		cockpitCard("review_gtfs_quality", "Review GTFS quality", gtfsQualityStatus(page), gtfsQualitySignal(page), "Fix source GTFS or draft data first, then re-import or rerun the allowlisted static validator.", "/admin/operations/gtfs-quality", []string{"docs/tutorials/gtfs-validation-triage.md"}, "Does not auto-edit GTFS or approve warnings."),
		cockpitCard("run_review_validator_health", "Run or review validator health", cockpitValidationStatus(page), fmt.Sprintf("overall=%s; tooling=%s", page.ValidationHealth.OverallStatus, page.ValidationHealth.ToolingStatus), "Admins can request the allowlisted server-side validation run; other roles can review the current state.", "/admin/operations/validation-health", []string{"docs/deployment/off-host-validation.md"}, "Does not run browser-supplied commands, paths, URLs, argument lists, binaries, or timeouts."),
		cockpitCard("manage_devices_vehicles", "Manage devices / vehicles", deviceStatus(page), deviceEvidence(page), "Rotate or rebind a one-time token, then confirm the latest accepted telemetry time and assignment state.", "/admin/operations/devices", []string{"docs/tutorials/device-avl-integration.md"}, "Does not display stored token values or prove hardware certification."),
		cockpitCard("synthetic_telemetry", "Send synthetic telemetry / review readiness", telemetrySimulatorStatus(page), telemetrySimulatorSignal(page), "Use the browser guide for safe simulator commands; keep tokens in the operator shell, not the browser.", "/admin/operations/telemetry-simulator", []string{"docs/tutorials/telemetry-simulator-and-device-trial.md"}, "Does not prove real service telemetry quality or real vendor compatibility."),
		cockpitCard("realtime_feed_state", "Review Vehicle Positions / Trip Updates", realtimeStatus(page), realtimeSignal(page), "Review telemetry freshness, Vehicle Positions usefulness, Trip Updates adapter diagnostics, and withheld reasons.", "/admin/operations/feed-health", []string{"docs/requirements-trip-updates.md"}, "Does not prove production-grade ETA quality."),
		cockpitCard("manage_alerts", "Manage alerts", alertStatus(page), feedEvidence(page, "alerts"), "Open the Alerts Console for lifecycle work, then review the public Alerts feed state.", "/admin/alerts/console", []string{"docs/requirements-calitp-compliance.md"}, "Does not prove consumer display or agency approval."),
		cockpitCard("connector_readiness", "Review connector readiness", operationsStatusDiagnosticOnly, fmt.Sprintf("%d connector categories; %d manifest entries", len(page.ConnectorHub.Categories), len(page.ConnectorHub.Registry.Entries)), "Use Connector Hub and connector tests to keep telemetry, prediction, validator, monitoring, and discovery integrations behind safe boundaries.", "/admin/operations/connectors", []string{"docs/integration-adapter-kit.md"}, "Does not load arbitrary backend plugins or prove named vendor compatibility."),
		cockpitCard("maintenance_tasks", "Review maintenance tasks", maintenanceStatusToCockpit(page.Maintenance.OverallStatus), page.Maintenance.OverallStatus, "Review weekly/monthly tasks, backup/restore configuration, last checks, and support summary guidance.", "/admin/operations/maintenance", []string{"docs/tutorials/small-agency-maintenance-guide.md"}, "Does not prove SLA, uptime, or production readiness."),
		cockpitCard("support_summary", "Export local support summary", operationsStatusDiagnosticOnly, "support bundle remains an operator-run local helper", "Use the maintenance page to see why a technical helper may still run the support-bundle command.", "/admin/operations/maintenance", []string{"docs/tutorials/operator-smoke-and-support-bundle.md"}, "Does not create retained evidence or upload private diagnostics."),
	}

	return operationsCockpitView{
		GeneratedAt:   page.GeneratedAt,
		AgencyID:      page.AgencyID,
		Boundary:      "Private Agency Operations Cockpit only. It summarizes setup, feeds, GTFS quality, validators, telemetry, connectors, readiness, and maintenance without creating evidence, contacting external parties, changing consumer statuses, or making compliance, adoption, hosted-service, SLA, vendor, production, or ETA-quality claims.",
		SetupProgress: progress,
		PrimaryCards:  cards,
		ClaimFlags:    operationsCockpitClaimFlag{},
	}
}

func cockpitProgressRow(id, label, status, signal, next, link, doesNotProve string) operationsCockpitRow {
	return operationsCockpitRow{ID: id, Label: label, Status: status, CurrentSignal: signal, NextAction: next, AdminLink: link, DoesNotProve: doesNotProve}
}

func cockpitCard(id, label, status, signal, next, link string, docs []string, doesNotProve string) operationsCockpitCard {
	return operationsCockpitCard{ID: id, Label: label, Status: status, CurrentSignal: signal, NextAction: next, AdminLink: link, DocsLinks: docs, DoesNotProve: doesNotProve}
}

func cockpitMetadataStatus(page operationsPage) string {
	if page.DiscoveryError != "" || page.PublicationError != "" && strings.TrimSpace(page.Discovery.PublicBaseURL) == "" {
		return operationsStatusMissing
	}
	if page.Discovery.Readiness.LicenseComplete && page.Discovery.Readiness.ContactComplete && strings.TrimSpace(page.Discovery.AgencyName) != "" {
		return operationsStatusReady
	}
	return operationsStatusNeedsReview
}

func cockpitMetadataSignal(page operationsPage) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	return fmt.Sprintf("agency=%q; public_base_url=%q; license_complete=%t; contact_complete=%t", page.Discovery.AgencyName, page.Discovery.PublicBaseURL, page.Discovery.Readiness.LicenseComplete, page.Discovery.Readiness.ContactComplete)
}

func cockpitActiveFeedStatus(page operationsPage) string {
	if strings.TrimSpace(page.ActiveFeedVersion) == "" {
		return operationsStatusMissing
	}
	return operationsStatusReady
}

func cockpitActiveFeedSignal(page operationsPage) string {
	if strings.TrimSpace(page.ActiveFeedVersion) == "" {
		return "no active schedule feed version is available"
	}
	return "active schedule feed version " + page.ActiveFeedVersion
}

func cockpitFeedHealthStatus(page operationsPage) string {
	counts := page.FeedHealth.Counts.Statuses
	if counts[checklistStatusBlocked] > 0 {
		return operationsStatusBlocked
	}
	if counts[checklistStatusMissing] > 0 {
		return operationsStatusMissing
	}
	if counts[checklistStatusNeedsReview] > 0 || counts[checklistStatusUnknown] > 0 {
		return operationsStatusNeedsReview
	}
	if page.FeedHealth.Counts.Rows == 5 {
		return operationsStatusReady
	}
	return operationsStatusUnknown
}

func cockpitFeedHealthSignal(page operationsPage) string {
	counts := page.FeedHealth.Counts.Statuses
	return fmt.Sprintf("%d rows; ready=%d needs_review=%d missing=%d blocked=%d unknown=%d", page.FeedHealth.Counts.Rows, counts[checklistStatusOK], counts[checklistStatusNeedsReview], counts[checklistStatusMissing], counts[checklistStatusBlocked], counts[checklistStatusUnknown])
}

func cockpitValidationStatus(page operationsPage) string {
	switch normalizeChecklistStatus(page.ValidationHealth.OverallStatus) {
	case checklistStatusBlocked:
		return operationsStatusBlocked
	case checklistStatusMissing:
		return operationsStatusMissing
	case checklistStatusNeedsReview:
		return operationsStatusNeedsReview
	}
	switch page.ValidationHealth.OverallStatus {
	case "recorded":
		return operationsStatusReady
	case "missing_tooling", "misconfigured_tooling", "blocked", "failed":
		return operationsStatusBlocked
	case "artifact_unavailable", "stale", "needs_review", "not_run", "runnable", "configured_for_tests":
		return operationsStatusNeedsReview
	case "":
		return operationsStatusUnknown
	default:
		return operationsStatusNeedsReview
	}
}

func cockpitTelemetryStatus(page operationsPage) string {
	if page.TelemetryError != "" {
		return operationsStatusUnknown
	}
	if len(page.Telemetry) == 0 {
		return operationsStatusMissing
	}
	if page.StaleCount > 0 {
		return operationsStatusNeedsReview
	}
	return operationsStatusReady
}

func cockpitReadinessStatus(page operationsPage) string {
	if len(page.ReadinessV2.Rows) == 0 {
		return operationsStatusUnknown
	}
	if page.ReadinessV2.Counts.Statuses[checklistStatusBlocked] > 0 {
		return operationsStatusBlocked
	}
	if page.ReadinessV2.Counts.Statuses[checklistStatusMissing] > 0 {
		return operationsStatusMissing
	}
	if page.ReadinessV2.Counts.Statuses[checklistStatusNeedsReview] > 0 || page.ReadinessV2.Counts.Statuses[checklistStatusUnknown] > 0 {
		return operationsStatusNeedsReview
	}
	return operationsStatusReady
}

func maintenanceStatusToCockpit(status string) string {
	switch strings.TrimSpace(status) {
	case "ready", "ok":
		return operationsStatusReady
	case "blocked", "failed", "unhealthy":
		return operationsStatusBlocked
	case "missing":
		return operationsStatusMissing
	case "needs_review":
		return operationsStatusNeedsReview
	case "":
		return operationsStatusUnknown
	default:
		return operationsStatusDiagnosticOnly
	}
}

func gtfsQualityStatus(page operationsPage) string {
	statuses := []string{page.GTFSQuality.Canonical.Status, page.GTFSQuality.InternalImporter.Status}
	for _, status := range statuses {
		switch strings.TrimSpace(status) {
		case "blocked", "failed", "error":
			return operationsStatusBlocked
		}
	}
	for _, status := range statuses {
		switch strings.TrimSpace(status) {
		case "missing", "not_run", "stale", "needs_review", "warning":
			return operationsStatusNeedsReview
		}
	}
	if strings.TrimSpace(page.ActiveFeedVersion) == "" {
		return operationsStatusMissing
	}
	return operationsStatusReady
}

func gtfsQualitySignal(page operationsPage) string {
	return fmt.Sprintf("canonical=%s; internal_importer=%s; active_feed=%s", firstNonEmpty(page.GTFSQuality.Canonical.Status, "not available"), firstNonEmpty(page.GTFSQuality.InternalImporter.Status, "not available"), firstNonEmpty(page.ActiveFeedVersion, "missing"))
}

func deviceStatus(page operationsPage) string {
	if page.DeviceError != "" {
		return operationsStatusUnknown
	}
	if len(page.Devices) == 0 {
		return operationsStatusMissing
	}
	for _, row := range page.DeviceRows {
		if row.Freshness == "stale" || strings.Contains(strings.ToLower(row.Assignment), "unknown") {
			return operationsStatusNeedsReview
		}
	}
	return operationsStatusReady
}

func telemetrySimulatorStatus(page operationsPage) string {
	if page.TelemetrySimulator.LoadError != "" {
		return operationsStatusBlocked
	}
	return operationsStatusDiagnosticOnly
}

func telemetrySimulatorSignal(page operationsPage) string {
	if page.TelemetrySimulator.LoadError != "" {
		return page.TelemetrySimulator.LoadError
	}
	return fmt.Sprintf("%d committed synthetic scenarios; browser execution disabled", len(page.TelemetrySimulator.Scenarios))
}

func realtimeStatus(page operationsPage) string {
	if len(page.Telemetry) == 0 {
		return operationsStatusMissing
	}
	if page.TripUpdatesQuality.Recorded {
		if page.TripUpdatesQuality.TripUpdatesEmitted == 0 || page.StaleCount > 0 {
			return operationsStatusNeedsReview
		}
		return operationsStatusReady
	}
	return operationsStatusNeedsReview
}

func realtimeSignal(page operationsPage) string {
	vp := fmt.Sprintf("Vehicle Positions source telemetry rows=%d, stale=%d", len(page.Telemetry), page.StaleCount)
	if !page.TripUpdatesQuality.Recorded {
		return vp + "; Trip Updates diagnostics=" + page.TripUpdatesQuality.Message
	}
	return fmt.Sprintf("%s; Trip Updates adapter=%s emitted=%d withheld_reasons=%d", vp, page.TripUpdatesQuality.AdapterName, page.TripUpdatesQuality.TripUpdatesEmitted, len(page.TripUpdatesQuality.WithheldByReason))
}

func alertStatus(page operationsPage) string {
	for _, row := range page.FeedHealth.Rows {
		if row.ID == "alerts" {
			switch row.Status {
			case checklistStatusOK:
				return operationsStatusReady
			case checklistStatusBlocked:
				return operationsStatusBlocked
			case checklistStatusMissing:
				return operationsStatusMissing
			default:
				return operationsStatusNeedsReview
			}
		}
	}
	return operationsStatusUnknown
}
