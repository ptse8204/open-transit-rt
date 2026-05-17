package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"open-transit-rt/internal/compliance"
)

const (
	gtfsWorkbenchImportHistoryLimit = 8
	gtfsWorkbenchPreviewRowLimit    = 10
	gtfsWorkbenchDiffSampleLimit    = 5
)

type gtfsImportHistoryReader interface {
	RecentGTFSImports(ctx context.Context, agencyID string, limit int) ([]compliance.GTFSImportRecord, error)
}

type gtfsSchedulePreviewReader interface {
	GTFSSchedulePreview(ctx context.Context, agencyID string, feedVersionID string, limit int) (compliance.GTFSSchedulePreview, error)
}

type gtfsDraftReviewReader interface {
	RecentGTFSDrafts(ctx context.Context, agencyID string, limit int) ([]compliance.GTFSDraftRecord, error)
	RecentGTFSDraftPublishes(ctx context.Context, agencyID string, limit int) ([]compliance.GTFSDraftPublishRecord, error)
}

type gtfsScheduleHistoryReader interface {
	RecentFeedVersions(ctx context.Context, agencyID string, limit int) ([]compliance.FeedVersionRecord, error)
}

type operationsGTFSWorkbenchView struct {
	GeneratedAt       time.Time                          `json:"generated_at"`
	AgencyID          string                             `json:"agency_id"`
	Boundary          string                             `json:"boundary"`
	ReviewSummary     []operationsGTFSReviewSummaryRow   `json:"review_summary"`
	ActiveFeedVersion operationsGTFSActiveFeedVersion    `json:"active_feed_version"`
	Import            operationsGTFSImportSummary        `json:"import"`
	VersionComparison operationsGTFSVersionComparison    `json:"version_comparison"`
	IssueTriage       operationsGTFSWorkbenchIssueTriage `json:"issue_triage"`
	Quality           operationsGTFSQualitySummary       `json:"quality"`
	ValidationHealth  operationsGTFSValidationSummary    `json:"validation_health"`
	Preview           operationsGTFSPreviewSummary       `json:"preview"`
	DraftReview       operationsGTFSDraftReviewSummary   `json:"draft_review"`
	ScheduleHistory   operationsGTFSScheduleHistory      `json:"schedule_history"`
	FeedOutput        operationsGTFSFeedOutputSummary    `json:"feed_output"`
	Actions           []operationsGTFSWorkbenchAction    `json:"actions"`
	ClaimFlags        operationsGTFSWorkbenchClaimFlags  `json:"claim_flags"`
}

type operationsGTFSReviewSummaryRow struct {
	ID              string `json:"id"`
	Label           string `json:"label"`
	Status          string `json:"status"`
	PlainLanguage   string `json:"plain_language"`
	SuggestedReview string `json:"suggested_review"`
	DoesNotProve    string `json:"does_not_prove"`
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

type operationsGTFSWorkbenchIssueTriage struct {
	Status        string                         `json:"status"`
	CurrentSignal string                         `json:"current_signal"`
	DisplayedRows int                            `json:"displayed_rows"`
	TotalRows     int                            `json:"total_rows"`
	HiddenRows    int                            `json:"hidden_rows"`
	Boundary      string                         `json:"boundary"`
	Rows          []operationsGTFSWorkbenchIssue `json:"rows"`
}

type operationsGTFSWorkbenchIssue struct {
	Severity            string   `json:"severity"`
	SourceLabel         string   `json:"source_label"`
	Family              string   `json:"family"`
	Codes               []string `json:"codes"`
	Count               int      `json:"count"`
	LikelyOwner         string   `json:"likely_owner"`
	PlainEnglishMeaning string   `json:"plain_english_meaning"`
	SuggestedFixPath    string   `json:"suggested_fix_path"`
	SafeNextAction      string   `json:"safe_next_action"`
	VerifyWith          string   `json:"verify_with"`
	DoesNotProve        string   `json:"does_not_prove"`
}

type operationsGTFSVersionComparison struct {
	Status                   string                        `json:"status"`
	HistoryStatus            string                        `json:"history_status"`
	ActiveFeedVersionID      string                        `json:"active_feed_version_id"`
	PreviousFeedVersionID    string                        `json:"previous_feed_version_id"`
	PreviousLifecycleState   string                        `json:"previous_lifecycle_state"`
	PreviousValidationStatus string                        `json:"previous_validation_status"`
	RowLimit                 int                           `json:"row_limit"`
	SampleLimit              int                           `json:"sample_limit"`
	CurrentSignal            string                        `json:"current_signal"`
	NextAction               string                        `json:"next_action"`
	ClaimBoundary            string                        `json:"claim_boundary"`
	FileDiffs                []operationsGTFSFileDiffRow   `json:"file_diffs"`
	EntityDiffs              []operationsGTFSEntityDiffRow `json:"entity_diffs"`
	ReviewRows               []operationsGTFSChangeRow     `json:"review_rows"`
}

type operationsGTFSFileDiffRow struct {
	File          string `json:"file"`
	Status        string `json:"status"`
	PreviousRows  int    `json:"previous_rows"`
	ActiveRows    int    `json:"active_rows"`
	DeltaRows     int    `json:"delta_rows"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	ClaimBoundary string `json:"claim_boundary"`
}

type operationsGTFSEntityDiffRow struct {
	Entity        string   `json:"entity"`
	Status        string   `json:"status"`
	PreviousRows  int      `json:"previous_rows"`
	ActiveRows    int      `json:"active_rows"`
	AddedSample   []string `json:"added_sample"`
	RemovedSample []string `json:"removed_sample"`
	ChangedSample []string `json:"changed_sample"`
	CurrentSignal string   `json:"current_signal"`
	NextAction    string   `json:"next_action"`
	ClaimBoundary string   `json:"claim_boundary"`
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
	Status        string                                    `json:"status"`
	RowLimit      int                                       `json:"row_limit"`
	CurrentSignal string                                    `json:"current_signal"`
	NextAction    string                                    `json:"next_action"`
	ClaimBoundary string                                    `json:"claim_boundary"`
	Counts        compliance.GTFSSchedulePreviewCounts      `json:"counts"`
	RequiredFiles []operationsGTFSRequiredFileStatus        `json:"required_files"`
	Sections      []operationsGTFSPreviewSection            `json:"sections"`
	Agency        []operationsGTFSScheduleAgencyRow         `json:"agency"`
	Routes        []compliance.GTFSScheduleRoutePreview     `json:"routes"`
	Stops         []compliance.GTFSScheduleStopPreview      `json:"stops"`
	Trips         []compliance.GTFSScheduleTripPreview      `json:"trips"`
	Calendar      []compliance.GTFSScheduleCalendarPreview  `json:"calendar"`
	Frequencies   []compliance.GTFSScheduleFrequencyPreview `json:"frequencies"`
}

type operationsGTFSRequiredFileStatus struct {
	File          string `json:"file"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	ClaimBoundary string `json:"claim_boundary"`
}

type operationsGTFSPreviewSection struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	RowsShown     int    `json:"rows_shown"`
	TotalRows     int    `json:"total_rows"`
	OverflowCount int    `json:"overflow_count"`
	CurrentSignal string `json:"current_signal"`
	ClaimBoundary string `json:"claim_boundary"`
}

type operationsGTFSScheduleAgencyRow struct {
	AgencyID string `json:"agency_id"`
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

type operationsGTFSFeedOutputSummary struct {
	Status        string `json:"status"`
	ScheduleURL   string `json:"schedule_url"`
	FeedsJSONURL  string `json:"feeds_json_url"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	ClaimBoundary string `json:"claim_boundary"`
}

type operationsGTFSDraftReviewSummary struct {
	Status          string                          `json:"status"`
	HistoryStatus   string                          `json:"history_status"`
	CurrentSignal   string                          `json:"current_signal"`
	NextAction      string                          `json:"next_action"`
	ClaimBoundary   string                          `json:"claim_boundary"`
	Drafts          []operationsGTFSDraftRow        `json:"drafts"`
	PublishAttempts []operationsGTFSDraftPublishRow `json:"publish_attempts"`
	Checklist       []operationsGTFSChangeRow       `json:"checklist"`
}

type operationsGTFSDraftRow struct {
	ID                         string    `json:"id"`
	Name                       string    `json:"name"`
	Status                     string    `json:"status"`
	BaseFeedVersionID          string    `json:"base_feed_version_id"`
	LastPublishedFeedVersionID string    `json:"last_published_feed_version_id"`
	LastPublishAttemptID       int64     `json:"last_publish_attempt_id"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

type operationsGTFSDraftPublishRow struct {
	ID            int64      `json:"id"`
	DraftID       string     `json:"draft_id"`
	FeedVersionID string     `json:"feed_version_id"`
	Status        string     `json:"status"`
	ErrorCount    int        `json:"error_count"`
	WarningCount  int        `json:"warning_count"`
	InfoCount     int        `json:"info_count"`
	StartedAt     time.Time  `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	Signal        string     `json:"signal"`
}

type operationsGTFSScheduleHistory struct {
	Status           string                         `json:"status"`
	HistoryStatus    string                         `json:"history_status"`
	CurrentSignal    string                         `json:"current_signal"`
	NextAction       string                         `json:"next_action"`
	ClaimBoundary    string                         `json:"claim_boundary"`
	FeedVersions     []operationsGTFSFeedVersionRow `json:"feed_versions"`
	RollbackGuidance []operationsGTFSChangeRow      `json:"rollback_guidance"`
}

type operationsGTFSFeedVersionRow struct {
	ID               string     `json:"id"`
	SourceType       string     `json:"source_type"`
	LifecycleState   string     `json:"lifecycle_state"`
	IsActive         bool       `json:"is_active"`
	ValidationStatus string     `json:"validation_status"`
	PublishedAt      *time.Time `json:"published_at"`
	ActivatedAt      *time.Time `json:"activated_at"`
	RetiredAt        *time.Time `json:"retired_at"`
	CreatedAt        time.Time  `json:"created_at"`
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
	view.Preview = h.buildGTFSPreviewSummary(r.Context(), page.AgencyID, view.ActiveFeedVersion.FeedVersionID)
	view.DraftReview = h.buildGTFSDraftReviewSummary(r.Context(), page.AgencyID, view.ActiveFeedVersion.FeedVersionID)
	view.ScheduleHistory = h.buildGTFSScheduleHistory(r.Context(), page.AgencyID, view.ActiveFeedVersion.FeedVersionID)
	view.VersionComparison = h.buildGTFSVersionComparison(r.Context(), page.AgencyID, view.ActiveFeedVersion.FeedVersionID, view.ScheduleHistory)
	view.IssueTriage = buildGTFSWorkbenchIssueTriage(page.GTFSQualityGuidance.FixPlanner)
	view.FeedOutput = buildGTFSWorkbenchFeedOutput(page.Discovery)
	view.Actions = buildGTFSWorkbenchActions(view, page)
	view.ReviewSummary = buildGTFSWorkbenchReviewSummary(view)
	return view
}

func buildGTFSWorkbenchReviewSummary(view operationsGTFSWorkbenchView) []operationsGTFSReviewSummaryRow {
	boundary := "This summary is private operator guidance only. It does not approve the schedule, prove the validator has no remaining notices, prove agency approval, prove consumer acceptance, or prove compliance."
	return []operationsGTFSReviewSummaryRow{
		{
			ID:              "required_files",
			Label:           "Required files",
			Status:          gtfsWorkbenchRequiredFilesStatus(view.Preview.RequiredFiles),
			PlainLanguage:   gtfsWorkbenchRequiredFilesSummary(view.Preview.RequiredFiles),
			SuggestedReview: "Start with blocked required files before reviewing route, stop, trip, service, or realtime behavior.",
			DoesNotProve:    boundary,
		},
		{
			ID:              "row_counts",
			Label:           "Row counts",
			Status:          view.Preview.Status,
			PlainLanguage:   gtfsWorkbenchRowCountSummary(view.Preview.Counts),
			SuggestedReview: "Compare row counts against the agency's expected service size, then review capped table previews below.",
			DoesNotProve:    boundary,
		},
		{
			ID:              "service_dates",
			Label:           "Service dates",
			Status:          gtfsWorkbenchServiceDateStatus(view.Preview),
			PlainLanguage:   gtfsWorkbenchServiceDateSummary(view.Preview),
			SuggestedReview: "Confirm weekday service, exceptions, holidays, and after-midnight service with the schedule owner.",
			DoesNotProve:    boundary,
		},
		{
			ID:              "route_stop_trip_review",
			Label:           "Routes, stops, and trips",
			Status:          gtfsWorkbenchCoreTableStatus(view.Preview.Counts),
			PlainLanguage:   gtfsWorkbenchCoreTableSummary(view.Preview.Counts),
			SuggestedReview: "Use the bounded route, stop, and trip previews to decide whether the source GTFS owner, GIS owner, or operations staff should review the data.",
			DoesNotProve:    boundary,
		},
		{
			ID:              "import_history",
			Label:           "Import history",
			Status:          view.Import.Status,
			PlainLanguage:   gtfsWorkbenchImportHistorySummary(view.Import),
			SuggestedReview: "Review the latest and previous imports before treating the active schedule as the right local version.",
			DoesNotProve:    boundary,
		},
		{
			ID:              "what_changed",
			Label:           "What changed",
			Status:          view.VersionComparison.Status,
			PlainLanguage:   view.VersionComparison.CurrentSignal,
			SuggestedReview: "Review file-level and entity-level diffs before relying on route, stop, trip, service, or frequency changes.",
			DoesNotProve:    boundary,
		},
		{
			ID:              "issue_triage",
			Label:           "Issue triage",
			Status:          view.IssueTriage.Status,
			PlainLanguage:   view.IssueTriage.CurrentSignal,
			SuggestedReview: "Assign each issue to the likely data owner, fix the source GTFS or reviewed draft, then rerun validation through the allowlisted path.",
			DoesNotProve:    boundary,
		},
	}
}

func gtfsWorkbenchRequiredFilesStatus(files []operationsGTFSRequiredFileStatus) string {
	if len(files) == 0 {
		return "missing"
	}
	status := "ok"
	for _, file := range files {
		switch file.Status {
		case "blocked":
			return "blocked"
		case "optional":
			if status == "ok" {
				status = "needs_review"
			}
		}
	}
	return status
}

func gtfsWorkbenchRequiredFilesSummary(files []operationsGTFSRequiredFileStatus) string {
	if len(files) == 0 {
		return "Required file checks are not available until an active schedule preview can load."
	}
	var blocked, optional, ok int
	for _, file := range files {
		switch file.Status {
		case "blocked":
			blocked++
		case "optional":
			optional++
		default:
			ok++
		}
	}
	return fmt.Sprintf("%d required or expected file checks are ready, %d are blocked, and %d optional file checks need context.", ok, blocked, optional)
}

func gtfsWorkbenchRowCountSummary(counts compliance.GTFSSchedulePreviewCounts) string {
	if counts.Routes+counts.Stops+counts.Trips+counts.StopTimes+counts.Calendar+counts.CalendarDates+counts.ShapePoints+counts.Frequencies == 0 {
		return "No schedule row counts are available yet."
	}
	return fmt.Sprintf("%d routes, %d stops, %d trips, %d stop times, %d calendar rows, %d exceptions, %d shape points, and %d frequency rows are stored for the active schedule preview.", counts.Routes, counts.Stops, counts.Trips, counts.StopTimes, counts.Calendar, counts.CalendarDates, counts.ShapePoints, counts.Frequencies)
}

func gtfsWorkbenchServiceDateStatus(preview operationsGTFSPreviewSummary) string {
	if preview.Counts.Calendar+preview.Counts.CalendarDates == 0 {
		return "blocked"
	}
	return "needs_review"
}

func gtfsWorkbenchServiceDateSummary(preview operationsGTFSPreviewSummary) string {
	if len(preview.Calendar) == 0 {
		if preview.Counts.CalendarDates > 0 {
			return fmt.Sprintf("No bounded calendar rows are visible, but %d calendar exception rows are stored.", preview.Counts.CalendarDates)
		}
		return "No service calendar rows are visible in the active preview."
	}
	start := preview.Calendar[0].StartDate
	end := preview.Calendar[0].EndDate
	for _, row := range preview.Calendar[1:] {
		if row.StartDate != "" && (start == "" || row.StartDate < start) {
			start = row.StartDate
		}
		if row.EndDate != "" && (end == "" || row.EndDate > end) {
			end = row.EndDate
		}
	}
	return fmt.Sprintf("%d service calendar rows are stored with %d exception rows; bounded preview service dates run from %s to %s.", preview.Counts.Calendar, preview.Counts.CalendarDates, firstNonEmpty(start, "not visible"), firstNonEmpty(end, "not visible"))
}

func gtfsWorkbenchCoreTableStatus(counts compliance.GTFSSchedulePreviewCounts) string {
	if counts.Routes == 0 || counts.Stops == 0 || counts.Trips == 0 || counts.StopTimes == 0 {
		return "blocked"
	}
	return "needs_review"
}

func gtfsWorkbenchCoreTableSummary(counts compliance.GTFSSchedulePreviewCounts) string {
	if counts.Routes+counts.Stops+counts.Trips+counts.StopTimes == 0 {
		return "Route, stop, trip, and stop-time counts are not available yet."
	}
	return fmt.Sprintf("The active preview has %d route rows, %d stop rows, %d trip rows, and %d stop-time rows.", counts.Routes, counts.Stops, counts.Trips, counts.StopTimes)
}

func gtfsWorkbenchImportHistorySummary(summary operationsGTFSImportSummary) string {
	if len(summary.History) == 0 {
		return "No import history rows are available in this runtime."
	}
	if summary.Latest == nil {
		return fmt.Sprintf("%d import history rows are available, but no latest import row was selected.", len(summary.History))
	}
	return fmt.Sprintf("Latest import %d is %s for feed version %s, with %d errors, %d warnings, and %d earlier import rows visible.", summary.Latest.ID, summary.Latest.Status, firstNonEmpty(summary.Latest.FeedVersionID, "not linked"), summary.Latest.ErrorCount, summary.Latest.WarningCount, len(summary.History)-1)
}

func buildGTFSWorkbenchIssueTriage(planner operationsGTFSQualityFixPlanner) operationsGTFSWorkbenchIssueTriage {
	const limit = 8
	triage := operationsGTFSWorkbenchIssueTriage{
		Status:        planner.Status,
		CurrentSignal: planner.Summary,
		TotalRows:     planner.TotalRows,
		Boundary:      "Issue triage rows are sanitized operator guidance only. They do not edit GTFS, mutate drafts, publish schedules, run validators, prove compliance, prove consumer acceptance, prove agency approval, or prove production readiness.",
	}
	if triage.Status == "" {
		triage.Status = "unknown"
	}
	rows := planner.Rows
	if len(rows) > limit {
		triage.HiddenRows = len(rows) - limit
		rows = rows[:limit]
	}
	for _, row := range rows {
		triage.Rows = append(triage.Rows, operationsGTFSWorkbenchIssue{
			Severity:            row.Severity,
			SourceLabel:         row.SourceLabel,
			Family:              row.Family,
			Codes:               append([]string(nil), row.Codes...),
			Count:               row.Count,
			LikelyOwner:         row.LikelyOwner,
			PlainEnglishMeaning: row.IssueSummary,
			SuggestedFixPath:    row.SafeFixSuggestion,
			SafeNextAction:      row.BeforeValidationPlan,
			VerifyWith:          row.VerifyWith,
			DoesNotProve:        triage.Boundary,
		})
	}
	triage.DisplayedRows = len(triage.Rows)
	if planner.HiddenRows > 0 {
		triage.HiddenRows += planner.HiddenRows
	}
	if triage.TotalRows == 0 {
		triage.CurrentSignal = "No grouped GTFS quality issue rows are available from the current validator/importer records."
	}
	return triage
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

func (h *handler) buildGTFSPreviewSummary(ctx context.Context, agencyID string, feedVersionID string) operationsGTFSPreviewSummary {
	summary := operationsGTFSPreviewSummary{
		Status:        "missing",
		RowLimit:      gtfsWorkbenchPreviewRowLimit,
		CurrentSignal: "No active schedule feed version is recorded, so preview tables cannot load.",
		NextAction:    "Import a GTFS ZIP or publish a GTFS Studio draft before reviewing schedule previews.",
		ClaimBoundary: "Preview tables are private operator aids only. They do not edit, publish, approve, validate, or certify GTFS.",
	}
	if strings.TrimSpace(feedVersionID) == "" {
		return summary
	}
	reader, ok := h.store.(gtfsSchedulePreviewReader)
	if !ok {
		summary.Status = "unknown"
		summary.CurrentSignal = "This runtime does not expose a GTFS schedule preview reader."
		summary.NextAction = "Use import, quality, validation, and feed output rows; ask a technical helper if table-level review is required."
		return summary
	}
	preview, err := reader.GTFSSchedulePreview(ctx, agencyID, feedVersionID, gtfsWorkbenchPreviewRowLimit)
	if err != nil {
		summary.Status = "unknown"
		summary.CurrentSignal = "GTFS schedule preview rows could not be read from existing private schedule tables."
		summary.NextAction = "Review importer and database state, then retry this private page."
		return summary
	}
	summary.RowLimit = preview.RowLimit
	summary.Counts = preview.Counts
	summary.RequiredFiles = gtfsRequiredFileStatuses(preview)
	summary.Sections = gtfsPreviewSections(preview)
	if preview.Agency.AgencyID != "" || preview.Agency.Name != "" {
		summary.Agency = []operationsGTFSScheduleAgencyRow{{
			AgencyID: preview.Agency.AgencyID,
			Name:     preview.Agency.Name,
			Timezone: preview.Agency.Timezone,
		}}
	}
	summary.Routes = preview.Routes
	summary.Stops = preview.Stops
	summary.Trips = preview.Trips
	summary.Calendar = preview.Calendar
	summary.Frequencies = preview.Frequencies
	summary.Status = "ok"
	for _, file := range summary.RequiredFiles {
		if file.Status == "blocked" {
			summary.Status = "blocked"
			break
		}
	}
	summary.CurrentSignal = fmt.Sprintf("Previewing up to %d rows per table for active schedule feed version %s.", summary.RowLimit, feedVersionID)
	summary.NextAction = "Use these bounded previews to decide which source GTFS owner should review routes, stops, trips, service calendars, or frequency-based service."
	return summary
}

func gtfsRequiredFileStatuses(preview compliance.GTFSSchedulePreview) []operationsGTFSRequiredFileStatus {
	boundary := "File checklist rows summarize local GTFS tables only and do not prove external validator success or approval."
	return []operationsGTFSRequiredFileStatus{
		requiredFileStatus("agency.txt", gtfsPreviewBoolStatus(preview.Agency.Name != ""), fmt.Sprintf("Agency record: %s.", firstNonEmpty(preview.Agency.Name, "not available")), "Confirm agency name and timezone in source GTFS before relying on public metadata.", boundary),
		requiredFileStatus("routes.txt", gtfsPreviewCountStatus(preview.Counts.Routes), fmt.Sprintf("%d route rows stored.", preview.Counts.Routes), "Review route names, modes, and source-system ownership.", boundary),
		requiredFileStatus("stops.txt", gtfsPreviewCountStatus(preview.Counts.Stops), fmt.Sprintf("%d stop rows stored.", preview.Counts.Stops), "Review stop names and coordinates with the schedule/source-data owner.", boundary),
		requiredFileStatus("trips.txt", gtfsPreviewCountStatus(preview.Counts.Trips), fmt.Sprintf("%d trip rows stored.", preview.Counts.Trips), "Review representative trips, route links, blocks, shapes, and direction IDs.", boundary),
		requiredFileStatus("stop_times.txt", gtfsPreviewCountStatus(preview.Counts.StopTimes), fmt.Sprintf("%d stop time rows stored.", preview.Counts.StopTimes), "Review stop-time ordering, after-midnight service, and repeated trip patterns in source GTFS.", boundary),
		requiredFileStatus("calendar.txt / calendar_dates.txt", gtfsPreviewCountStatus(preview.Counts.Calendar+preview.Counts.CalendarDates), fmt.Sprintf("%d calendar rows and %d exception rows stored.", preview.Counts.Calendar, preview.Counts.CalendarDates), "Review service dates, exceptions, and holiday service before relying on active trips.", boundary),
		requiredFileStatus("frequencies.txt", gtfsPreviewOptionalCountStatus(preview.Counts.Frequencies), fmt.Sprintf("%d frequency rows stored.", preview.Counts.Frequencies), "If frequency-based service is expected, confirm headways and exact-times values.", boundary),
		requiredFileStatus("shapes.txt", gtfsPreviewOptionalCountStatus(preview.Counts.ShapePoints), fmt.Sprintf("%d shape point rows stored.", preview.Counts.ShapePoints), "If route geometry is expected, confirm shape coverage and distance data with source GTFS.", boundary),
	}
}

func requiredFileStatus(file string, status string, signal string, nextAction string, boundary string) operationsGTFSRequiredFileStatus {
	return operationsGTFSRequiredFileStatus{File: file, Status: status, CurrentSignal: signal, NextAction: nextAction, ClaimBoundary: boundary}
}

func gtfsPreviewBoolStatus(ok bool) string {
	if ok {
		return "ok"
	}
	return "blocked"
}

func gtfsPreviewCountStatus(count int) string {
	return gtfsPreviewBoolStatus(count > 0)
}

func gtfsPreviewOptionalCountStatus(count int) string {
	if count > 0 {
		return "ok"
	}
	return "optional"
}

func gtfsPreviewSections(preview compliance.GTFSSchedulePreview) []operationsGTFSPreviewSection {
	boundary := "Rows are capped private previews and may omit additional records."
	return []operationsGTFSPreviewSection{
		previewSection("agency", "Agency", gtfsPreviewBoolStatus(preview.Agency.Name != ""), lenIf(preview.Agency.Name != ""), lenIf(preview.Agency.Name != ""), "Agency identity available from local agency metadata.", boundary),
		previewSection("routes", "Routes", gtfsPreviewCountStatus(preview.Counts.Routes), len(preview.Routes), preview.Counts.Routes, "Route rows show route IDs, names, and route types.", boundary),
		previewSection("stops", "Stops", gtfsPreviewCountStatus(preview.Counts.Stops), len(preview.Stops), preview.Counts.Stops, "Stop rows show stop IDs, names, and coordinates.", boundary),
		previewSection("trips", "Trips", gtfsPreviewCountStatus(preview.Counts.Trips), len(preview.Trips), preview.Counts.Trips, "Trip rows show route, service, block, shape, and direction links.", boundary),
		previewSection("calendar", "Calendar / Service", gtfsPreviewCountStatus(preview.Counts.Calendar+preview.Counts.CalendarDates), len(preview.Calendar), preview.Counts.Calendar, "Calendar rows show weekly service windows; exception count is summarized separately.", boundary),
		previewSection("frequencies", "Frequencies", gtfsPreviewOptionalCountStatus(preview.Counts.Frequencies), len(preview.Frequencies), preview.Counts.Frequencies, "Frequency rows show headway-based service where present.", boundary),
	}
}

func previewSection(id string, label string, status string, rowsShown int, totalRows int, signal string, boundary string) operationsGTFSPreviewSection {
	overflow := totalRows - rowsShown
	if overflow < 0 {
		overflow = 0
	}
	return operationsGTFSPreviewSection{
		ID:            id,
		Label:         label,
		Status:        status,
		RowsShown:     rowsShown,
		TotalRows:     totalRows,
		OverflowCount: overflow,
		CurrentSignal: signal,
		ClaimBoundary: boundary,
	}
}

func lenIf(ok bool) int {
	if ok {
		return 1
	}
	return 0
}

func (h *handler) buildGTFSDraftReviewSummary(ctx context.Context, agencyID string, activeFeedVersionID string) operationsGTFSDraftReviewSummary {
	boundary := "Draft review summarizes GTFS Studio state only. It does not publish from the Workbench, auto-fix GTFS, approve service changes, create evidence, or prove compliance, consumer acceptance, agency approval, or production readiness."
	summary := operationsGTFSDraftReviewSummary{
		Status:        "unknown",
		HistoryStatus: "unavailable",
		CurrentSignal: "GTFS Studio draft state is not available in this runtime.",
		NextAction:    "Use GTFS Studio for draft authoring and return here after publish attempts are recorded.",
		ClaimBoundary: boundary,
		Checklist: []operationsGTFSChangeRow{
			gtfsWorkbenchChangeRow("Workbench publish action", "ok", "No publish action exists on the GTFS Workbench.", "Use GTFS Studio admin publish after reviewing draft data and validation feedback.", boundary),
		},
	}
	reader, ok := h.store.(gtfsDraftReviewReader)
	if !ok {
		return summary
	}
	drafts, draftErr := reader.RecentGTFSDrafts(ctx, agencyID, gtfsWorkbenchImportHistoryLimit)
	publishes, publishErr := reader.RecentGTFSDraftPublishes(ctx, agencyID, gtfsWorkbenchImportHistoryLimit)
	if draftErr != nil && publishErr != nil {
		return summary
	}
	summary.HistoryStatus = "recorded"
	for _, draft := range drafts {
		summary.Drafts = append(summary.Drafts, gtfsDraftRowView(draft))
	}
	for _, publish := range publishes {
		summary.PublishAttempts = append(summary.PublishAttempts, gtfsDraftPublishRowView(publish))
	}
	summary.Status = gtfsDraftReviewStatus(summary.Drafts, summary.PublishAttempts, activeFeedVersionID)
	summary.CurrentSignal = gtfsDraftReviewSignal(summary)
	summary.NextAction = gtfsDraftReviewNextAction(summary.Status)
	summary.Checklist = append(summary.Checklist,
		gtfsWorkbenchChangeRow("Draft and active schedule separation", "ok", "Draft records and active published feed versions remain separate.", "Compare draft base and last published feed version before an admin publish.", boundary),
		gtfsWorkbenchChangeRow("Publish confirmation path", "needs_review", "Publishing remains in GTFS Studio and requires admin confirmation.", "Open GTFS Studio, review draft rows, then publish only after source owner review.", boundary),
		gtfsWorkbenchChangeRow("Post-publish validation", gtfsDraftPublishValidationStatus(summary.PublishAttempts), gtfsDraftPublishValidationSignal(summary.PublishAttempts), "After publish, review GTFS quality, validation health, feed health, and preview tables.", boundary),
	)
	return summary
}

func gtfsDraftRowView(record compliance.GTFSDraftRecord) operationsGTFSDraftRow {
	return operationsGTFSDraftRow{
		ID:                         record.ID,
		Name:                       boundedDisplayText(record.Name, 120),
		Status:                     firstNonEmpty(record.Status, "unknown"),
		BaseFeedVersionID:          record.BaseFeedVersionID,
		LastPublishedFeedVersionID: record.LastPublishedFeedVersionID,
		LastPublishAttemptID:       record.LastPublishAttemptID,
		CreatedAt:                  record.CreatedAt.UTC(),
		UpdatedAt:                  record.UpdatedAt.UTC(),
	}
}

func gtfsDraftPublishRowView(record compliance.GTFSDraftPublishRecord) operationsGTFSDraftPublishRow {
	row := operationsGTFSDraftPublishRow{
		ID:            record.ID,
		DraftID:       record.DraftID,
		FeedVersionID: record.FeedVersionID,
		Status:        firstNonEmpty(record.Status, "unknown"),
		ErrorCount:    record.ErrorCount,
		WarningCount:  record.WarningCount,
		InfoCount:     record.InfoCount,
		StartedAt:     record.StartedAt.UTC(),
		CompletedAt:   record.CompletedAt,
	}
	row.Signal = fmt.Sprintf("%s publish attempt with %d errors, %d warnings, and %d info notices.", row.Status, row.ErrorCount, row.WarningCount, row.InfoCount)
	return row
}

func boundedDisplayText(value string, max int) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "not recorded"
	}
	if max <= 0 || len(trimmed) <= max {
		return trimmed
	}
	if max <= 3 {
		return trimmed[:max]
	}
	return trimmed[:max-3] + "..."
}

func gtfsDraftReviewStatus(drafts []operationsGTFSDraftRow, publishes []operationsGTFSDraftPublishRow, activeFeedVersionID string) string {
	if len(drafts) == 0 && len(publishes) == 0 {
		return "missing"
	}
	if len(publishes) > 0 {
		latest := publishes[0]
		if latest.Status == "failed" || latest.ErrorCount > 0 {
			return "blocked"
		}
		if latest.WarningCount > 0 {
			return "needs_review"
		}
		if latest.FeedVersionID != "" && activeFeedVersionID != "" && latest.FeedVersionID != activeFeedVersionID {
			return "needs_review"
		}
	}
	for _, draft := range drafts {
		if draft.Status == "draft" {
			return "needs_review"
		}
	}
	return "ok"
}

func gtfsDraftReviewSignal(summary operationsGTFSDraftReviewSummary) string {
	if len(summary.Drafts) == 0 && len(summary.PublishAttempts) == 0 {
		return "No GTFS Studio drafts or publish attempts are recorded."
	}
	return fmt.Sprintf("%d recent drafts and %d recent publish attempts are recorded.", len(summary.Drafts), len(summary.PublishAttempts))
}

func gtfsDraftReviewNextAction(status string) string {
	switch status {
	case "blocked":
		return "Fix draft validation or publish errors in GTFS Studio before relying on the active schedule."
	case "needs_review":
		return "Review draft rows, latest publish feedback, and source-owner approval before any admin publish."
	case "ok":
		return "Continue local operator review with quality, validation, feed health, and preview rows."
	case "missing":
		return "Use GTFS Studio only if typed draft authoring is needed; otherwise continue with ZIP import review."
	default:
		return "Review GTFS Studio state with a technical helper if draft authoring is part of this workflow."
	}
}

func gtfsDraftPublishValidationStatus(publishes []operationsGTFSDraftPublishRow) string {
	if len(publishes) == 0 {
		return "missing"
	}
	latest := publishes[0]
	if latest.Status == "failed" || latest.ErrorCount > 0 {
		return "blocked"
	}
	if latest.WarningCount > 0 {
		return "needs_review"
	}
	return "ok"
}

func gtfsDraftPublishValidationSignal(publishes []operationsGTFSDraftPublishRow) string {
	if len(publishes) == 0 {
		return "No draft publish attempt is recorded."
	}
	return publishes[0].Signal
}

func (h *handler) buildGTFSScheduleHistory(ctx context.Context, agencyID string, activeFeedVersionID string) operationsGTFSScheduleHistory {
	boundary := "Schedule history is a private review of local feed_version records. It does not execute rollback, publish a release, create evidence, prove final-root readiness, or prove consumer ingestion."
	summary := operationsGTFSScheduleHistory{
		Status:        "unknown",
		HistoryStatus: "unavailable",
		CurrentSignal: "Schedule feed-version history is not available in this runtime.",
		NextAction:    "Use import and draft publish records, then ask a technical helper to inspect feed_version state if rollback review is needed.",
		ClaimBoundary: boundary,
		RollbackGuidance: []operationsGTFSChangeRow{
			gtfsWorkbenchChangeRow("Browser rollback execution", "ok", "No rollback POST route exists in the GTFS Workbench.", "Use documented operator procedures and technical-helper review before any rollback outside this page.", boundary),
		},
	}
	reader, ok := h.store.(gtfsScheduleHistoryReader)
	if !ok {
		return summary
	}
	records, err := reader.RecentFeedVersions(ctx, agencyID, gtfsWorkbenchImportHistoryLimit)
	if err != nil {
		return summary
	}
	summary.HistoryStatus = "recorded"
	for _, record := range records {
		summary.FeedVersions = append(summary.FeedVersions, feedVersionRowView(record))
	}
	summary.Status = scheduleHistoryStatus(summary.FeedVersions, activeFeedVersionID)
	summary.CurrentSignal = scheduleHistorySignal(summary.FeedVersions, activeFeedVersionID)
	summary.NextAction = scheduleHistoryNextAction(summary.Status)
	summary.RollbackGuidance = append(summary.RollbackGuidance,
		gtfsWorkbenchChangeRow("Active feed visibility", scheduleHistoryActiveStatus(summary.FeedVersions, activeFeedVersionID), scheduleHistorySignal(summary.FeedVersions, activeFeedVersionID), "Confirm the active feed version before rollback or support-bundle decisions.", boundary),
		gtfsWorkbenchChangeRow("Rollback candidate review", rollbackCandidateStatus(summary.FeedVersions), rollbackCandidateSignal(summary.FeedVersions), "If rollback is required, confirm candidate source type, validation state, feed-health impact, and operator approval outside this Workbench.", boundary),
		gtfsWorkbenchChangeRow("After rollback review", "needs_review", "Any rollback would require follow-up GTFS quality, validation health, feed health, realtime assignment, and public URL review.", "Record the outcome as normal private operations state only; do not create evidence or consumer claims here.", boundary),
	)
	return summary
}

func feedVersionRowView(record compliance.FeedVersionRecord) operationsGTFSFeedVersionRow {
	return operationsGTFSFeedVersionRow{
		ID:               record.ID,
		SourceType:       record.SourceType,
		LifecycleState:   record.LifecycleState,
		IsActive:         record.IsActive,
		ValidationStatus: record.ValidationStatus,
		PublishedAt:      record.PublishedAt,
		ActivatedAt:      record.ActivatedAt,
		RetiredAt:        record.RetiredAt,
		CreatedAt:        record.CreatedAt.UTC(),
	}
}

func scheduleHistoryStatus(rows []operationsGTFSFeedVersionRow, activeFeedVersionID string) string {
	if len(rows) == 0 {
		return "missing"
	}
	activeCount := 0
	activeFound := false
	for _, row := range rows {
		if row.IsActive {
			activeCount++
			if activeFeedVersionID == "" || row.ID == activeFeedVersionID {
				activeFound = true
			}
		}
	}
	if activeCount != 1 || !activeFound {
		return "needs_review"
	}
	return "ok"
}

func scheduleHistorySignal(rows []operationsGTFSFeedVersionRow, activeFeedVersionID string) string {
	if len(rows) == 0 {
		return "No local feed_version rows are available for schedule history review."
	}
	active := "not found"
	for _, row := range rows {
		if row.IsActive {
			active = row.ID
			break
		}
	}
	if activeFeedVersionID != "" && active != activeFeedVersionID {
		return fmt.Sprintf("Local history active feed is %s; feed discovery active schedule is %s.", active, activeFeedVersionID)
	}
	return fmt.Sprintf("%d recent feed versions are available; active feed version is %s.", len(rows), active)
}

func scheduleHistoryNextAction(status string) string {
	switch status {
	case "ok":
		return "Keep reviewing import, draft publish, preview, validation, and feed health rows before operational decisions."
	case "needs_review":
		return "Resolve active feed-version mismatch or duplicate-active records before rollback or publish decisions."
	case "missing":
		return "Import or publish a schedule before rollback review."
	default:
		return "Ask a technical helper to inspect schedule history if rollback review is required."
	}
}

func scheduleHistoryActiveStatus(rows []operationsGTFSFeedVersionRow, activeFeedVersionID string) string {
	return scheduleHistoryStatus(rows, activeFeedVersionID)
}

func rollbackCandidateStatus(rows []operationsGTFSFeedVersionRow) string {
	if len(rows) <= 1 {
		return "missing"
	}
	for _, row := range rows {
		if !row.IsActive && row.LifecycleState == "retired" {
			return "needs_review"
		}
	}
	return "unknown"
}

func rollbackCandidateSignal(rows []operationsGTFSFeedVersionRow) string {
	if len(rows) <= 1 {
		return "No prior local feed version is visible as a rollback candidate in this bounded history."
	}
	for _, row := range rows {
		if !row.IsActive && row.LifecycleState == "retired" {
			return "At least one retired local feed version exists; review outside this Workbench before any rollback procedure."
		}
	}
	return "Recent history has multiple feed versions, but no retired rollback candidate is visible in the bounded list."
}

func (h *handler) buildGTFSVersionComparison(ctx context.Context, agencyID string, activeFeedVersionID string, history operationsGTFSScheduleHistory) operationsGTFSVersionComparison {
	boundary := "Version comparison is a private read-only review of local published GTFS tables. It does not edit GTFS, execute rollback, prove canonical validator success, prove compliance, create evidence, or prove consumer ingestion."
	summary := operationsGTFSVersionComparison{
		Status:              "missing",
		HistoryStatus:       history.HistoryStatus,
		ActiveFeedVersionID: activeFeedVersionID,
		RowLimit:            gtfsWorkbenchPreviewRowLimit,
		SampleLimit:         gtfsWorkbenchDiffSampleLimit,
		CurrentSignal:       "No active schedule feed version is recorded for version comparison.",
		NextAction:          "Import or publish a schedule before reviewing version-to-version changes.",
		ClaimBoundary:       boundary,
		ReviewRows: []operationsGTFSChangeRow{
			gtfsWorkbenchChangeRow("Published feed versions only", "ok", "This comparison uses published feed_version rows and does not read GTFS Studio draft rows as active data.", "Keep draft editing and published feed review separate.", boundary),
		},
	}
	if strings.TrimSpace(activeFeedVersionID) == "" {
		return summary
	}
	candidate := previousFeedVersionForComparison(history.FeedVersions, activeFeedVersionID)
	if candidate == nil {
		summary.CurrentSignal = "No previous published feed version is visible in the bounded schedule history."
		summary.NextAction = "Treat this as a first recorded active schedule or ask a technical helper to inspect older feed versions if rollback review is needed."
		summary.ReviewRows = append(summary.ReviewRows,
			gtfsWorkbenchChangeRow("Previous version candidate", "missing", "No non-active previous feed version is visible in this bounded Workbench history.", "Import or publish a later schedule before active-vs-previous comparison is available here.", boundary),
		)
		return summary
	}
	summary.PreviousFeedVersionID = candidate.ID
	summary.PreviousLifecycleState = firstNonEmpty(candidate.LifecycleState, "unknown")
	summary.PreviousValidationStatus = firstNonEmpty(candidate.ValidationStatus, "not_run")

	reader, ok := h.store.(gtfsSchedulePreviewReader)
	if !ok {
		summary.Status = "unknown"
		summary.CurrentSignal = "This runtime does not expose the GTFS schedule preview reader needed for version comparison."
		summary.NextAction = "Use feed-version history and import records, then ask a technical helper for table-level diff review."
		summary.ReviewRows = append(summary.ReviewRows,
			gtfsWorkbenchChangeRow("Preview reader", "unknown", "GTFS table previews are unavailable in this runtime.", "Review database/repository wiring before relying on version comparison.", boundary),
		)
		return summary
	}
	previousPreview, previousErr := reader.GTFSSchedulePreview(ctx, agencyID, candidate.ID, gtfsWorkbenchPreviewRowLimit)
	activePreview, activeErr := reader.GTFSSchedulePreview(ctx, agencyID, activeFeedVersionID, gtfsWorkbenchPreviewRowLimit)
	if previousErr != nil || activeErr != nil {
		summary.Status = "unknown"
		summary.CurrentSignal = "One or both GTFS schedule previews could not be read for version comparison."
		summary.NextAction = "Review feed-version table state and schedule rows with a technical helper before rollback decisions."
		summary.ReviewRows = append(summary.ReviewRows,
			gtfsWorkbenchChangeRow("Previous preview", statusFromErr(previousErr), previewReadSignal("previous", candidate.ID, previousErr), "Resolve missing or unreadable schedule rows before rollback review.", boundary),
			gtfsWorkbenchChangeRow("Active preview", statusFromErr(activeErr), previewReadSignal("active", activeFeedVersionID, activeErr), "Resolve missing or unreadable schedule rows before relying on the active schedule.", boundary),
		)
		return summary
	}
	summary.FileDiffs = gtfsVersionFileDiffs(previousPreview.Counts, activePreview.Counts, boundary)
	summary.EntityDiffs = gtfsVersionEntityDiffs(previousPreview, activePreview, boundary)
	summary.Status = gtfsVersionComparisonStatus(summary.FileDiffs, summary.EntityDiffs)
	summary.CurrentSignal = gtfsVersionComparisonSignal(summary)
	summary.NextAction = gtfsVersionComparisonNextAction(summary.Status)
	summary.ReviewRows = append(summary.ReviewRows,
		gtfsWorkbenchChangeRow("Rollback candidate", rollbackCandidateStatusForComparison(*candidate), rollbackCandidateSignalForComparison(*candidate), "Before any rollback outside this page, confirm validation state, feed health, realtime assignment impact, operator approval, and audit expectations.", boundary),
		gtfsWorkbenchChangeRow("Realtime assignment review", "needs_review", "Current vehicle assignments may reference the active feed version that would become retired after rollback.", "After any external rollback procedure, review realtime matching, Vehicle Positions trip descriptors, Trip Updates withholding, and Alerts links before relying on outputs.", boundary),
		gtfsWorkbenchChangeRow("Draft-only rollback command design", "blocked", "No executable rollback command or browser POST route exists in this Workbench.", "A future rollback command must be admin-only, CSRF-protected for cookie auth, confirmation-based, agency-scoped, transactional, audited, and followed by validation/feed-health review.", boundary),
	)
	return summary
}

func previousFeedVersionForComparison(rows []operationsGTFSFeedVersionRow, activeFeedVersionID string) *operationsGTFSFeedVersionRow {
	for _, row := range rows {
		if row.ID == "" || row.ID == activeFeedVersionID || row.IsActive {
			continue
		}
		if row.LifecycleState == "retired" {
			copyRow := row
			return &copyRow
		}
	}
	for _, row := range rows {
		if row.ID == "" || row.ID == activeFeedVersionID || row.IsActive {
			continue
		}
		copyRow := row
		return &copyRow
	}
	return nil
}

func statusFromErr(err error) string {
	if err != nil {
		return "unknown"
	}
	return "ok"
}

func previewReadSignal(label string, feedVersionID string, err error) string {
	if err != nil {
		return fmt.Sprintf("%s preview for feed version %s could not be read.", label, feedVersionID)
	}
	return fmt.Sprintf("%s preview for feed version %s is readable.", label, feedVersionID)
}

func gtfsVersionFileDiffs(previous compliance.GTFSSchedulePreviewCounts, active compliance.GTFSSchedulePreviewCounts, boundary string) []operationsGTFSFileDiffRow {
	return []operationsGTFSFileDiffRow{
		gtfsVersionFileDiff("routes.txt", previous.Routes, active.Routes, "Review route additions, removals, route type changes, and route naming with the schedule owner.", boundary),
		gtfsVersionFileDiff("stops.txt", previous.Stops, active.Stops, "Review stop additions, removals, name changes, and coordinate changes with the schedule owner.", boundary),
		gtfsVersionFileDiff("trips.txt", previous.Trips, active.Trips, "Review trip additions, removals, route/service/block/shape links, and repeated trip patterns.", boundary),
		gtfsVersionFileDiff("stop_times.txt", previous.StopTimes, active.StopTimes, "Review stop-time volume changes, ordering, after-midnight times, and stop references in source GTFS.", boundary),
		gtfsVersionFileDiff("calendar.txt", previous.Calendar, active.Calendar, "Review service IDs, weekday patterns, and service date ranges before relying on trip coverage.", boundary),
		gtfsVersionFileDiff("calendar_dates.txt", previous.CalendarDates, active.CalendarDates, "Review added/removed service exceptions, holidays, and canceled-service dates.", boundary),
		gtfsVersionFileDiff("shapes.txt", previous.ShapePoints, active.ShapePoints, "Review shape coverage and distance data where route geometry affects matching or previews.", boundary),
		gtfsVersionFileDiff("frequencies.txt", previous.Frequencies, active.Frequencies, "Review frequency rows as service-affecting changes for trip instance and matching behavior.", boundary),
	}
}

func gtfsVersionFileDiff(file string, previousRows int, activeRows int, nextAction string, boundary string) operationsGTFSFileDiffRow {
	delta := activeRows - previousRows
	status := "ok"
	signal := fmt.Sprintf("%s row count is unchanged at %d.", file, activeRows)
	if delta != 0 {
		status = "needs_review"
		signal = fmt.Sprintf("%s row count changed from %d to %d (%+d).", file, previousRows, activeRows, delta)
	}
	return operationsGTFSFileDiffRow{
		File:          file,
		Status:        status,
		PreviousRows:  previousRows,
		ActiveRows:    activeRows,
		DeltaRows:     delta,
		CurrentSignal: signal,
		NextAction:    nextAction,
		ClaimBoundary: boundary,
	}
}

func gtfsVersionEntityDiffs(previous compliance.GTFSSchedulePreview, active compliance.GTFSSchedulePreview, boundary string) []operationsGTFSEntityDiffRow {
	return []operationsGTFSEntityDiffRow{
		gtfsVersionEntityDiff("Routes", previous.Counts.Routes, active.Counts.Routes, routeSignatureMap(previous.Routes), routeSignatureMap(active.Routes), "Review sampled route ID additions, removals, and changed names or route types.", boundary),
		gtfsVersionEntityDiff("Stops", previous.Counts.Stops, active.Counts.Stops, stopSignatureMap(previous.Stops), stopSignatureMap(active.Stops), "Review sampled stop ID additions, removals, renamed stops, and coordinate changes.", boundary),
		gtfsVersionEntityDiff("Trips", previous.Counts.Trips, active.Counts.Trips, tripSignatureMap(previous.Trips), tripSignatureMap(active.Trips), "Review sampled trip ID additions, removals, route/service/block changes, and shape changes.", boundary),
		gtfsVersionEntityDiff("Service calendars", previous.Counts.Calendar+previous.Counts.CalendarDates, active.Counts.Calendar+active.Counts.CalendarDates, calendarSignatureMap(previous.Calendar), calendarSignatureMap(active.Calendar), "Review sampled service ID changes plus the calendar_dates row-count signal before relying on service coverage.", boundary),
		gtfsVersionEntityDiff("Frequencies", previous.Counts.Frequencies, active.Counts.Frequencies, frequencySignatureMap(previous.Frequencies), frequencySignatureMap(active.Frequencies), "Review sampled frequency changes because headways and exact-times values affect trip instance interpretation.", boundary),
	}
}

func gtfsVersionEntityDiff(entity string, previousRows int, activeRows int, previous map[string]string, active map[string]string, nextAction string, boundary string) operationsGTFSEntityDiffRow {
	row := operationsGTFSEntityDiffRow{
		Entity:        entity,
		Status:        "ok",
		PreviousRows:  previousRows,
		ActiveRows:    activeRows,
		AddedSample:   boundedKeyDiff(active, previous, gtfsWorkbenchDiffSampleLimit),
		RemovedSample: boundedKeyDiff(previous, active, gtfsWorkbenchDiffSampleLimit),
		ChangedSample: boundedChangedKeys(previous, active, gtfsWorkbenchDiffSampleLimit),
		NextAction:    nextAction,
		ClaimBoundary: boundary,
	}
	if previousRows != activeRows || len(row.AddedSample) > 0 || len(row.RemovedSample) > 0 || len(row.ChangedSample) > 0 {
		row.Status = "needs_review"
	}
	row.CurrentSignal = entityDiffSignal(row)
	return row
}

func routeSignatureMap(rows []compliance.GTFSScheduleRoutePreview) map[string]string {
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		out[boundedDisplayText(row.ID, 160)] = strings.Join([]string{row.ShortName, row.LongName, row.RouteType}, "\x00")
	}
	return out
}

func stopSignatureMap(rows []compliance.GTFSScheduleStopPreview) map[string]string {
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		out[boundedDisplayText(row.ID, 160)] = fmt.Sprintf("%s\x00%.6f\x00%.6f", row.Name, row.Lat, row.Lon)
	}
	return out
}

func tripSignatureMap(rows []compliance.GTFSScheduleTripPreview) map[string]string {
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		out[boundedDisplayText(row.ID, 160)] = strings.Join([]string{row.RouteID, row.ServiceID, row.BlockID, row.ShapeID, row.DirectionID}, "\x00")
	}
	return out
}

func calendarSignatureMap(rows []compliance.GTFSScheduleCalendarPreview) map[string]string {
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		out[boundedDisplayText(row.ServiceID, 160)] = strings.Join([]string{row.Days, row.StartDate, row.EndDate}, "\x00")
	}
	return out
}

func frequencySignatureMap(rows []compliance.GTFSScheduleFrequencyPreview) map[string]string {
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		key := boundedDisplayText(row.TripID+" @ "+row.StartTime, 160)
		out[key] = fmt.Sprintf("%s\x00%d\x00%d", row.EndTime, row.HeadwaySecs, row.ExactTimes)
	}
	return out
}

func boundedKeyDiff(left map[string]string, right map[string]string, limit int) []string {
	keys := make([]string, 0)
	for key := range left {
		if _, ok := right[key]; !ok {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return boundedStringSlice(keys, limit)
}

func boundedChangedKeys(previous map[string]string, active map[string]string, limit int) []string {
	keys := make([]string, 0)
	for key, previousValue := range previous {
		activeValue, ok := active[key]
		if !ok {
			continue
		}
		if activeValue != previousValue {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return boundedStringSlice(keys, limit)
}

func boundedStringSlice(values []string, limit int) []string {
	if limit <= 0 || len(values) <= limit {
		return values
	}
	return append([]string(nil), values[:limit]...)
}

func entityDiffSignal(row operationsGTFSEntityDiffRow) string {
	parts := []string{
		fmt.Sprintf("rows %d -> %d", row.PreviousRows, row.ActiveRows),
		fmt.Sprintf("%d added samples", len(row.AddedSample)),
		fmt.Sprintf("%d removed samples", len(row.RemovedSample)),
		fmt.Sprintf("%d changed samples", len(row.ChangedSample)),
	}
	if row.Status == "ok" {
		return "Bounded sample comparison shows no row-count or sampled identity changes (" + strings.Join(parts, "; ") + ")."
	}
	return "Review bounded sample comparison (" + strings.Join(parts, "; ") + ")."
}

func gtfsVersionComparisonStatus(files []operationsGTFSFileDiffRow, entities []operationsGTFSEntityDiffRow) string {
	if len(files) == 0 && len(entities) == 0 {
		return "missing"
	}
	for _, row := range files {
		if row.Status == "needs_review" || row.Status == "blocked" {
			return "needs_review"
		}
	}
	for _, row := range entities {
		if row.Status == "needs_review" || row.Status == "blocked" {
			return "needs_review"
		}
	}
	return "ok"
}

func gtfsVersionComparisonSignal(summary operationsGTFSVersionComparison) string {
	if summary.PreviousFeedVersionID == "" || summary.ActiveFeedVersionID == "" {
		return "Version comparison needs both a previous and active feed version."
	}
	if summary.Status == "ok" {
		return fmt.Sprintf("Active feed version %s and previous feed version %s have matching row counts and bounded samples.", summary.ActiveFeedVersionID, summary.PreviousFeedVersionID)
	}
	return fmt.Sprintf("Active feed version %s differs from previous feed version %s in row counts or bounded samples.", summary.ActiveFeedVersionID, summary.PreviousFeedVersionID)
}

func gtfsVersionComparisonNextAction(status string) string {
	switch status {
	case "ok":
		return "Continue reviewing validation, feed health, and realtime assignment context before operational decisions."
	case "needs_review":
		return "Review source schedule changes with the schedule owner, then verify validation, feed health, and realtime implications."
	case "missing":
		return "Import or publish another schedule before active-vs-previous comparison is available."
	default:
		return "Ask a technical helper to inspect feed-version and schedule table state before rollback or publish decisions."
	}
}

func rollbackCandidateStatusForComparison(row operationsGTFSFeedVersionRow) string {
	if row.ID == "" {
		return "missing"
	}
	if row.LifecycleState != "retired" {
		return "needs_review"
	}
	if row.ValidationStatus == "failed" {
		return "blocked"
	}
	return "needs_review"
}

func rollbackCandidateSignalForComparison(row operationsGTFSFeedVersionRow) string {
	if row.ID == "" {
		return "No rollback candidate is selected."
	}
	return fmt.Sprintf("Candidate %s is %s with validation status %s.", row.ID, firstNonEmpty(row.LifecycleState, "unknown"), firstNonEmpty(row.ValidationStatus, "not_run"))
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
		summary.Diff = []operationsGTFSChangeRow{gtfsWorkbenchChangeRow("Latest import", "missing", "No GTFS import record is stored for this agency.", "Use Import GTFS for a ZIP upload or safe URL import.", summary.ClaimBoundary)}
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
		ClaimBoundary:          "GTFS quality rows and the private fix planner separate internal importer checks from canonical MobilityData static validator diagnostics. They are advisory review signals only.",
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
			ID:            "draft_publish_review",
			Label:         "Draft publish review",
			Status:        view.DraftReview.Status,
			CurrentSignal: view.DraftReview.CurrentSignal,
			NextAction:    view.DraftReview.NextAction,
			AdminLink:     "/admin/gtfs-studio",
			ClaimBoundary: view.DraftReview.ClaimBoundary,
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
