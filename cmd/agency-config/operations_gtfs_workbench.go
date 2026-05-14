package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"open-transit-rt/internal/compliance"
)

const gtfsWorkbenchImportHistoryLimit = 8

type gtfsImportHistoryReader interface {
	RecentGTFSImports(ctx context.Context, agencyID string, limit int) ([]compliance.GTFSImportRecord, error)
}

type operationsGTFSWorkbenchView struct {
	GeneratedAt       time.Time                         `json:"generated_at"`
	AgencyID          string                            `json:"agency_id"`
	Boundary          string                            `json:"boundary"`
	ActiveFeedVersion operationsGTFSActiveFeedVersion   `json:"active_feed_version"`
	Import            operationsGTFSImportSummary       `json:"import"`
	Quality           operationsGTFSQualitySummary      `json:"quality"`
	ValidationHealth  operationsGTFSValidationSummary   `json:"validation_health"`
	Preview           operationsGTFSPreviewSummary      `json:"preview"`
	FeedOutput        operationsGTFSFeedOutputSummary   `json:"feed_output"`
	Actions           []operationsGTFSWorkbenchAction   `json:"actions"`
	ClaimFlags        operationsGTFSWorkbenchClaimFlags `json:"claim_flags"`
}

type operationsGTFSActiveFeedVersion struct {
	Status             string     `json:"status"`
	FeedVersionID      string     `json:"feed_version_id"`
	RevisionTimestamp  *time.Time `json:"revision_timestamp"`
	ActivationStatus   string     `json:"activation_status"`
	CanonicalPublicURL string     `json:"canonical_public_url"`
	CurrentSignal      string     `json:"current_signal"`
	NextAction         string     `json:"next_action"`
	ClaimBoundary      string     `json:"claim_boundary"`
}

type operationsGTFSImportSummary struct {
	Status        string                    `json:"status"`
	HistoryStatus string                    `json:"history_status"`
	HistoryError  string                    `json:"history_error,omitempty"`
	Latest        *operationsGTFSImportRow  `json:"latest,omitempty"`
	Previous      *operationsGTFSImportRow  `json:"previous,omitempty"`
	History       []operationsGTFSImportRow `json:"history"`
	Diff          []operationsGTFSChangeRow `json:"diff"`
	NextAction    string                    `json:"next_action"`
	ClaimBoundary string                    `json:"claim_boundary"`
}

type operationsGTFSImportRow struct {
	ID                int64      `json:"id"`
	Status            string     `json:"status"`
	FeedVersionID     string     `json:"feed_version_id"`
	SourceName        string     `json:"source_name"`
	SourceSHA256      string     `json:"source_sha256"`
	SourceSHA256Short string     `json:"source_sha256_short"`
	SourceByteSize    int64      `json:"source_byte_size"`
	SourceByteText    string     `json:"source_byte_text"`
	ErrorCount        int        `json:"error_count"`
	WarningCount      int        `json:"warning_count"`
	InfoCount         int        `json:"info_count"`
	ActorID           string     `json:"actor_id"`
	StartedAt         time.Time  `json:"started_at"`
	CompletedAt       *time.Time `json:"completed_at"`
	Signal            string     `json:"signal"`
}

type operationsGTFSChangeRow struct {
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	ClaimBoundary string `json:"claim_boundary"`
}

type operationsGTFSQualitySummary struct {
	Status                 string `json:"status"`
	CanonicalStatus        string `json:"canonical_status"`
	InternalImporterStatus string `json:"internal_importer_status"`
	CanonicalResultStatus  string `json:"canonical_result_status"`
	InternalResultStatus   string `json:"internal_result_status"`
	CanonicalAction        string `json:"canonical_action"`
	InternalImporterAction string `json:"internal_importer_action"`
	BlockingGroups         int    `json:"blocking_groups"`
	NeedsReviewGroups      int    `json:"needs_review_groups"`
	InformationalGroups    int    `json:"informational_groups"`
	ClaimBoundary          string `json:"claim_boundary"`
}

type operationsGTFSValidationSummary struct {
	Status                    string     `json:"status"`
	ToolingStatus             string     `json:"tooling_status"`
	ArtifactStatus            string     `json:"artifact_status"`
	LatestResultStatus        string     `json:"latest_result_status"`
	LatestResultAt            *time.Time `json:"latest_result_at"`
	ActiveFeedVersionID       string     `json:"active_feed_version_id"`
	LatestResultFeedVersionID string     `json:"latest_result_feed_version_id"`
	StaleStatus               string     `json:"stale_status"`
	NextAction                string     `json:"next_action"`
	ClaimBoundary             string     `json:"claim_boundary"`
}

type operationsGTFSPreviewSummary struct {
	Status        string `json:"status"`
	RowLimit      int    `json:"row_limit"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	ClaimBoundary string `json:"claim_boundary"`
}

type operationsGTFSFeedOutputSummary struct {
	Status        string `json:"status"`
	ScheduleURL   string `json:"schedule_url"`
	FeedsJSONURL  string `json:"feeds_json_url"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	ClaimBoundary string `json:"claim_boundary"`
}

type operationsGTFSWorkbenchAction struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	AdminLink     string `json:"admin_link"`
	ClaimBoundary string `json:"claim_boundary"`
}

type operationsGTFSWorkbenchClaimFlags struct {
	AutomaticGTFSEditEnabled       bool `json:"automatic_gtfs_edit_enabled"`
	SchedulePublishedFromWorkbench bool `json:"schedule_published_from_workbench"`
	ValidatorRunFromWorkbench      bool `json:"validator_run_from_workbench"`
	ExternalEvidenceCreated        bool `json:"external_evidence_created"`
	ConsumerStatusesChanged        bool `json:"consumer_statuses_changed"`
	ComplianceClaimed              bool `json:"compliance_claimed"`
	AgencyApprovalClaimed          bool `json:"agency_approval_claimed"`
	ConsumerAcceptanceClaimed      bool `json:"consumer_acceptance_claimed"`
	FinalRootReadinessClaimed      bool `json:"final_root_readiness_claimed"`
	PublicLaunchClaimed            bool `json:"public_launch_claimed"`
	HostedSaaSClaimed              bool `json:"hosted_saas_claimed"`
	ProductionReadinessClaimed     bool `json:"production_readiness_claimed"`
	VendorCompatibilityClaimed     bool `json:"vendor_compatibility_claimed"`
	ProductionGradeETAClaimed      bool `json:"production_grade_eta_claimed"`
}

func (h *handler) buildGTFSWorkbenchView(r *http.Request, page operationsPage) operationsGTFSWorkbenchView {
	view := operationsGTFSWorkbenchView{
		GeneratedAt: page.GeneratedAt,
		AgencyID:    page.AgencyID,
		Boundary:    "Private GTFS Workbench reads existing local records only. It does not auto-edit GTFS, publish schedules, run validators, create retained evidence, contact external systems, change consumer status, or prove compliance, consumer acceptance, public launch, production readiness, or final-root readiness.",
		ClaimFlags:  operationsGTFSWorkbenchClaimFlags{},
	}
	view.ActiveFeedVersion = buildGTFSActiveFeedVersion(page.Discovery)
	view.Import = h.buildGTFSImportSummary(r.Context(), page.AgencyID, view.ActiveFeedVersion.FeedVersionID)
	view.Quality = buildGTFSWorkbenchQualitySummary(page.GTFSQuality)
	view.ValidationHealth = buildGTFSWorkbenchValidationSummary(page.ValidationHealth)
	view.Preview = operationsGTFSPreviewSummary{
		Status:        "not_loaded",
		RowLimit:      0,
		CurrentSignal: "Bounded schedule preview tables are added in the preview checkpoint after this route is established.",
		NextAction:    "Use the import, quality, validation, and feed rows first; then review bounded previews once available.",
		ClaimBoundary: "A preview is a private operator aid only and does not change or approve GTFS.",
	}
	view.FeedOutput = buildGTFSWorkbenchFeedOutput(page.Discovery)
	view.Actions = buildGTFSWorkbenchActions(view, page)
	return view
}

func buildGTFSActiveFeedVersion(discovery compliance.FeedDiscovery) operationsGTFSActiveFeedVersion {
	versionID, revision := scheduleFeedVersion(discovery.Feeds)
	row := operationsGTFSActiveFeedVersion{
		Status:        "missing",
		FeedVersionID: versionID,
		ClaimBoundary: "Active schedule metadata is local publication state only; it does not prove source correctness, validator success, compliance, or consumer use.",
	}
	for _, feed := range discovery.Feeds {
		if feed.FeedType != "schedule" {
			continue
		}
		row.ActivationStatus = firstNonEmpty(feed.ActivationStatus, "unknown")
		row.CanonicalPublicURL = feed.CanonicalPublicURL
		row.RevisionTimestamp = revision
		break
	}
	if versionID == "" {
		row.CurrentSignal = "No active schedule feed version is recorded."
		row.NextAction = "Import a GTFS ZIP or publish a GTFS Studio draft, then return here to review the private schedule state."
		return row
	}
	row.Status = "ok"
	row.CurrentSignal = "Active local schedule feed version is recorded."
	row.NextAction = "Review latest import, GTFS quality, validator health, and feed output before relying on the schedule."
	return row
}

func (h *handler) buildGTFSImportSummary(ctx context.Context, agencyID string, activeFeedVersionID string) operationsGTFSImportSummary {
	summary := operationsGTFSImportSummary{
		Status:        "unknown",
		HistoryStatus: "unavailable",
		NextAction:    "Review the active schedule and GTFS quality rows; import history is unavailable in this runtime.",
		ClaimBoundary: "Import records describe local importer outcomes only. They do not prove canonical validator success, agency approval, compliance, consumer acceptance, or production readiness.",
	}
	reader, ok := h.store.(gtfsImportHistoryReader)
	if !ok {
		summary.HistoryError = "GTFS import history reader is not available in this runtime."
		summary.Diff = []operationsGTFSChangeRow{gtfsWorkbenchChangeRow("Import history", "unknown", "No GTFS import history reader is configured.", "Review active schedule metadata and GTFS quality rows.", summary.ClaimBoundary)}
		return summary
	}
	records, err := reader.RecentGTFSImports(ctx, agencyID, gtfsWorkbenchImportHistoryLimit)
	if err != nil {
		summary.HistoryError = "GTFS import history is not available in this runtime."
		summary.Diff = []operationsGTFSChangeRow{gtfsWorkbenchChangeRow("Import history", "unknown", "Recent GTFS imports could not be read.", "Review active schedule metadata and GTFS quality rows.", summary.ClaimBoundary)}
		return summary
	}
	summary.HistoryStatus = "recorded"
	for _, record := range records {
		summary.History = append(summary.History, gtfsImportRowView(record))
	}
	if len(summary.History) == 0 {
		summary.Status = "missing"
		summary.NextAction = "Import a GTFS ZIP or publish a GTFS Studio draft before schedule review."
		summary.Diff = []operationsGTFSChangeRow{gtfsWorkbenchChangeRow("Latest import", "missing", "No GTFS import record is stored for this agency.", "Use Browser GTFS Import for a ZIP upload or safe URL import.", summary.ClaimBoundary)}
		return summary
	}
	latest := summary.History[0]
	summary.Latest = &latest
	if len(summary.History) > 1 {
		previous := summary.History[1]
		summary.Previous = &previous
	}
	summary.Status = gtfsImportSummaryStatus(latest, activeFeedVersionID)
	summary.NextAction = gtfsImportNextAction(summary.Status, latest)
	summary.Diff = gtfsImportDiffRows(summary.Latest, summary.Previous, activeFeedVersionID, summary.ClaimBoundary)
	return summary
}

func gtfsImportRowView(record compliance.GTFSImportRecord) operationsGTFSImportRow {
	row := operationsGTFSImportRow{
		ID:                record.ID,
		Status:            firstNonEmpty(record.Status, "unknown"),
		FeedVersionID:     record.FeedVersionID,
		SourceName:        safeGTFSSourceName(record.SourceFilename),
		SourceSHA256:      strings.TrimSpace(record.SourceSHA256),
		SourceSHA256Short: shortSHA(record.SourceSHA256),
		SourceByteSize:    record.SourceByteSize,
		SourceByteText:    formatByteCount(record.SourceByteSize),
		ErrorCount:        record.ErrorCount,
		WarningCount:      record.WarningCount,
		InfoCount:         record.InfoCount,
		ActorID:           record.ActorID,
		StartedAt:         record.StartedAt.UTC(),
		CompletedAt:       record.CompletedAt,
	}
	row.Signal = fmt.Sprintf("%s import with %d errors, %d warnings, and %d info notices.", row.Status, row.ErrorCount, row.WarningCount, row.InfoCount)
	return row
}

func safeGTFSSourceName(value string) string {
	trimmed := strings.TrimSpace(strings.ReplaceAll(value, "\\", "/"))
	if trimmed == "" {
		return "not recorded"
	}
	if idx := strings.LastIndex(trimmed, "/"); idx >= 0 {
		trimmed = trimmed[idx+1:]
	}
	if trimmed == "" {
		return "not recorded"
	}
	if len(trimmed) > 120 {
		return trimmed[:117] + "..."
	}
	return trimmed
}

func shortSHA(value string) string {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) <= 12 {
		return trimmed
	}
	return trimmed[:12]
}

func formatByteCount(value int64) string {
	if value < 0 {
		return "not available"
	}
	return fmt.Sprintf("%d bytes", value)
}

func gtfsImportSummaryStatus(latest operationsGTFSImportRow, activeFeedVersionID string) string {
	if latest.Status == "failed" || latest.ErrorCount > 0 {
		return "blocked"
	}
	if activeFeedVersionID != "" && latest.FeedVersionID != "" && latest.FeedVersionID != activeFeedVersionID {
		return "needs_review"
	}
	if latest.WarningCount > 0 {
		return "needs_review"
	}
	if latest.Status == "published" {
		return "ok"
	}
	return "needs_review"
}

func gtfsImportNextAction(status string, latest operationsGTFSImportRow) string {
	switch status {
	case "blocked":
		return "Fix the source GTFS or import input, then rerun the existing import flow."
	case "needs_review":
		return "Compare this import with the active schedule, then review GTFS quality and validation rows."
	case "ok":
		return "Continue local operator review with quality, validation, feed output, and preview sections."
	default:
		if latest.ID == 0 {
			return "Import a GTFS ZIP or publish a GTFS Studio draft before schedule review."
		}
		return "Review import history and active schedule state before stronger wording."
	}
}

func gtfsImportDiffRows(latest *operationsGTFSImportRow, previous *operationsGTFSImportRow, activeFeedVersionID string, boundary string) []operationsGTFSChangeRow {
	if latest == nil {
		return []operationsGTFSChangeRow{gtfsWorkbenchChangeRow("Latest import", "missing", "No GTFS import record is stored for this agency.", "Import a GTFS ZIP or publish a draft through existing private flows.", boundary)}
	}
	rows := []operationsGTFSChangeRow{
		gtfsWorkbenchChangeRow("Active schedule comparison", gtfsActiveImportComparisonStatus(*latest, activeFeedVersionID), gtfsActiveImportComparisonSignal(*latest, activeFeedVersionID), "Keep draft, import, and active schedule state separate before publishing or relying on outputs.", boundary),
	}
	if previous == nil {
		rows = append(rows, gtfsWorkbenchChangeRow("Previous import comparison", "missing", "No prior import is available for checksum or byte-size comparison.", "Treat this as the first recorded importer signal and review source data directly.", boundary))
		return rows
	}
	if latest.SourceSHA256 != "" && previous.SourceSHA256 != "" && latest.SourceSHA256 == previous.SourceSHA256 {
		rows = append(rows, gtfsWorkbenchChangeRow("Source checksum", "ok", "Latest import source checksum matches the previous recorded import.", "Still review schedule previews and validation because matching checksums do not prove correctness.", boundary))
	} else {
		rows = append(rows, gtfsWorkbenchChangeRow("Source checksum", "needs_review", "Latest import source checksum differs from the previous recorded import.", "Review source-system export changes and validation triage before relying on the schedule.", boundary))
	}
	if latest.SourceByteSize == previous.SourceByteSize {
		rows = append(rows, gtfsWorkbenchChangeRow("Source byte count", "ok", "Latest import source byte count matches the previous recorded import.", "Use previews and validation to understand whether content semantics changed.", boundary))
	} else {
		rows = append(rows, gtfsWorkbenchChangeRow("Source byte count", "needs_review", fmt.Sprintf("Latest import source byte count changed from %s to %s.", previous.SourceByteText, latest.SourceByteText), "Review route, stop, trip, and calendar summaries before relying on the new schedule.", boundary))
	}
	return rows
}

func gtfsActiveImportComparisonStatus(latest operationsGTFSImportRow, activeFeedVersionID string) string {
	if activeFeedVersionID == "" {
		return "missing"
	}
	if latest.FeedVersionID == "" {
		return "needs_review"
	}
	if latest.FeedVersionID == activeFeedVersionID {
		return "ok"
	}
	return "needs_review"
}

func gtfsActiveImportComparisonSignal(latest operationsGTFSImportRow, activeFeedVersionID string) string {
	if activeFeedVersionID == "" {
		return "No active schedule feed version is recorded."
	}
	if latest.FeedVersionID == "" {
		return "Latest import has no linked feed version."
	}
	if latest.FeedVersionID == activeFeedVersionID {
		return "Latest import feed version matches the active schedule feed version."
	}
	return "Latest import feed version differs from the active schedule feed version."
}

func gtfsWorkbenchChangeRow(label string, status string, signal string, nextAction string, boundary string) operationsGTFSChangeRow {
	return operationsGTFSChangeRow{Label: label, Status: status, CurrentSignal: signal, NextAction: nextAction, ClaimBoundary: boundary}
}

func buildGTFSWorkbenchQualitySummary(triage compliance.GTFSQualityTriage) operationsGTFSQualitySummary {
	summary := operationsGTFSQualitySummary{
		CanonicalStatus:        firstNonEmpty(triage.Canonical.Status, "unknown"),
		InternalImporterStatus: firstNonEmpty(triage.InternalImporter.Status, "unknown"),
		CanonicalResultStatus:  firstNonEmpty(triage.Canonical.ResultStatus, "not_run"),
		InternalResultStatus:   firstNonEmpty(triage.InternalImporter.ResultStatus, "not_run"),
		CanonicalAction:        triage.Canonical.RecommendedAction,
		InternalImporterAction: triage.InternalImporter.RecommendedAction,
		ClaimBoundary:          "GTFS quality rows separate internal importer checks from canonical MobilityData static validator diagnostics. They are private review signals only.",
	}
	for _, section := range []compliance.GTFSQualitySection{triage.Canonical, triage.InternalImporter} {
		for _, group := range section.Groups {
			switch group.Severity {
			case compliance.GTFSQualityBlocking:
				summary.BlockingGroups++
			case compliance.GTFSQualityNeedsReview:
				summary.NeedsReviewGroups++
			default:
				summary.InformationalGroups++
			}
		}
	}
	summary.Status = worseGTFSWorkbenchStatus(summary.CanonicalStatus, summary.InternalImporterStatus)
	return summary
}

func worseGTFSWorkbenchStatus(statuses ...string) string {
	worst := "unknown"
	for _, status := range statuses {
		switch status {
		case "blocked", "failed", "missing_tooling", "artifact_unavailable", "misconfigured_tooling":
			return "blocked"
		case "needs_review", "stale", "not_run":
			if worst != "blocked" {
				worst = "needs_review"
			}
		case "ok", "informational", "recorded", "runnable", "configured_for_tests":
			if worst == "unknown" {
				worst = "ok"
			}
		}
	}
	return worst
}

func buildGTFSWorkbenchValidationSummary(summary compliance.ValidationHealthSummary) operationsGTFSValidationSummary {
	out := operationsGTFSValidationSummary{
		Status:        "unknown",
		ToolingStatus: firstNonEmpty(summary.ToolingStatus, "unknown"),
		NextAction:    "Run or review private validator health after an active schedule exists.",
		ClaimBoundary: "Validator health is a supporting private signal only; it does not prove canonical validator success, compliance, consumer acceptance, or production readiness.",
	}
	for _, row := range summary.Feeds {
		if row.FeedType != "schedule" {
			continue
		}
		out.Status = firstNonEmpty(row.HealthStatus, "unknown")
		out.ToolingStatus = firstNonEmpty(row.ToolingStatus, out.ToolingStatus)
		out.ArtifactStatus = firstNonEmpty(row.ArtifactStatus, "unknown")
		out.LatestResultStatus = firstNonEmpty(row.LatestResultStatus, "not_run")
		out.LatestResultAt = row.LatestResultAt
		out.ActiveFeedVersionID = row.ActiveFeedVersionID
		out.LatestResultFeedVersionID = row.LatestResultFeedVersionID
		out.StaleStatus = firstNonEmpty(row.StaleStatus, "unknown")
		out.NextAction = firstNonEmpty(row.NextAction, out.NextAction)
		out.ClaimBoundary = firstNonEmpty(row.ClaimBoundary, out.ClaimBoundary)
		return out
	}
	return out
}

func buildGTFSWorkbenchFeedOutput(discovery compliance.FeedDiscovery) operationsGTFSFeedOutputSummary {
	out := operationsGTFSFeedOutputSummary{
		Status:        "missing",
		NextAction:    "Configure publication metadata and active schedule feed output before sharing URLs for review.",
		ClaimBoundary: "Feed URLs are local/deployment records only. Copying them does not prove final-root ownership, consumer ingestion, compliance, or public launch.",
	}
	for _, feed := range discovery.Feeds {
		if feed.FeedType != "schedule" {
			continue
		}
		out.ScheduleURL = feed.CanonicalPublicURL
		break
	}
	if strings.TrimSpace(discovery.PublicBaseURL) != "" {
		out.FeedsJSONURL = strings.TrimRight(discovery.PublicBaseURL, "/") + "/public/feeds.json"
	}
	if out.ScheduleURL == "" && out.FeedsJSONURL == "" {
		out.CurrentSignal = "No schedule URL or feed discovery URL is configured."
		return out
	}
	out.Status = "needs_review"
	out.CurrentSignal = "Schedule and feed discovery URLs are configured for private review."
	out.NextAction = "Open Feed Health to review fetch and validation signals before sharing any URL outside authorized review."
	return out
}

func buildGTFSWorkbenchActions(view operationsGTFSWorkbenchView, page operationsPage) []operationsGTFSWorkbenchAction {
	actions := []operationsGTFSWorkbenchAction{
		{
			ID:            "publication_metadata",
			Label:         "Agency and feed metadata",
			Status:        statusFromMissing(page.PublicationError == ""),
			CurrentSignal: firstNonEmpty(page.PublicationError, "publication metadata is configured for private review"),
			NextAction:    "Review agency profile, feed base URL, license, and contact metadata before schedule review.",
			AdminLink:     "/admin/operations/setup#publication-metadata",
			ClaimBoundary: "Metadata completeness does not prove final-root readiness, compliance, or agency approval.",
		},
		{
			ID:            "active_schedule",
			Label:         "Active local schedule",
			Status:        view.ActiveFeedVersion.Status,
			CurrentSignal: view.ActiveFeedVersion.CurrentSignal,
			NextAction:    view.ActiveFeedVersion.NextAction,
			AdminLink:     "/admin/operations/gtfs-import",
			ClaimBoundary: view.ActiveFeedVersion.ClaimBoundary,
		},
		{
			ID:            "latest_import",
			Label:         "Latest import",
			Status:        view.Import.Status,
			CurrentSignal: gtfsImportActionSignal(view.Import),
			NextAction:    view.Import.NextAction,
			AdminLink:     "/admin/operations/gtfs-import",
			ClaimBoundary: view.Import.ClaimBoundary,
		},
		{
			ID:            "quality_triage",
			Label:         "GTFS quality triage",
			Status:        view.Quality.Status,
			CurrentSignal: fmt.Sprintf("Canonical=%s; internal importer=%s.", view.Quality.CanonicalStatus, view.Quality.InternalImporterStatus),
			NextAction:    "Review likely owner, affected files, safe fix path, and verification rows.",
			AdminLink:     "/admin/operations/gtfs-quality",
			ClaimBoundary: view.Quality.ClaimBoundary,
		},
		{
			ID:            "schedule_validation",
			Label:         "Schedule validator health",
			Status:        view.ValidationHealth.Status,
			CurrentSignal: fmt.Sprintf("Artifact=%s; latest result=%s; stale=%s.", firstNonEmpty(view.ValidationHealth.ArtifactStatus, "unknown"), firstNonEmpty(view.ValidationHealth.LatestResultStatus, "not_run"), firstNonEmpty(view.ValidationHealth.StaleStatus, "unknown")),
			NextAction:    view.ValidationHealth.NextAction,
			AdminLink:     "/admin/operations/validation-health",
			ClaimBoundary: view.ValidationHealth.ClaimBoundary,
		},
		{
			ID:            "feed_output",
			Label:         "Feed output review",
			Status:        view.FeedOutput.Status,
			CurrentSignal: view.FeedOutput.CurrentSignal,
			NextAction:    view.FeedOutput.NextAction,
			AdminLink:     "/admin/operations/feed-health",
			ClaimBoundary: view.FeedOutput.ClaimBoundary,
		},
	}
	return actions
}

func statusFromMissing(ok bool) string {
	if ok {
		return "ok"
	}
	return "missing"
}

func gtfsImportActionSignal(summary operationsGTFSImportSummary) string {
	if summary.Latest == nil {
		return "No GTFS import is recorded."
	}
	return summary.Latest.Signal
}
