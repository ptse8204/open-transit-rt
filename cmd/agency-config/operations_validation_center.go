package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
)

type operationsValidationCenterView struct {
	GeneratedAt       time.Time                              `json:"generated_at"`
	AgencyID          string                                 `json:"agency_id"`
	Boundary          string                                 `json:"boundary"`
	FeedRows          []operationsValidationCenterFeedRow    `json:"feed_rows"`
	ValidationHistory []operationsValidationCenterValidation `json:"validation_history"`
	ValidatorHealth   []operationsValidationCenterValidation `json:"validator_health"`
	GTFSQuality       []operationsValidationCenterQuality    `json:"gtfs_quality"`
	IssueDrilldowns   []operationsValidationCenterIssue      `json:"issue_drilldowns"`
	ReadinessTimeline []operationsValidationCenterTimeline   `json:"readiness_timeline"`
	Blockers          []operationsValidationCenterBlocker    `json:"blockers"`
	ConsumerTracker   []operationsValidationCenterConsumer   `json:"consumer_tracker"`
	Counts            operationsValidationCenterCounts       `json:"counts"`
	ClaimFlags        operationsValidationCenterClaimFlags   `json:"claim_flags"`
}

type operationsValidationCenterFeedRow struct {
	ID             string `json:"id"`
	Label          string `json:"label"`
	PublicPath     string `json:"public_path"`
	ConfiguredURL  string `json:"configured_url"`
	Status         string `json:"status"`
	LastChecked    string `json:"last_checked"`
	HTTPStatus     string `json:"http_status"`
	ContentType    string `json:"content_type"`
	Freshness      string `json:"freshness"`
	ValidatorState string `json:"validator_state"`
	HealthState    string `json:"health_state"`
	CurrentSignal  string `json:"current_signal"`
	WhatThisMeans  string `json:"what_this_means"`
	NextAction     string `json:"next_action"`
	DoesNotProve   string `json:"does_not_prove"`
}

type operationsValidationCenterValidation struct {
	ID                        string `json:"id"`
	FeedType                  string `json:"feed_type"`
	Label                     string `json:"label"`
	ValidatorID               string `json:"validator_id"`
	ValidatorName             string `json:"validator_name"`
	Status                    string `json:"status"`
	ToolingStatus             string `json:"tooling_status"`
	ArtifactStatus            string `json:"artifact_status"`
	LatestResultStatus        string `json:"latest_result_status"`
	LatestResultAt            string `json:"latest_result_at"`
	ActiveFeedVersionID       string `json:"active_feed_version_id"`
	LatestResultFeedVersionID string `json:"latest_result_feed_version_id"`
	StaleStatus               string `json:"stale_status"`
	HealthStatus              string `json:"health_status"`
	CurrentSignal             string `json:"current_signal"`
	WhatThisMeans             string `json:"what_this_means"`
	NextAction                string `json:"next_action"`
	DoesNotProve              string `json:"does_not_prove"`
}

type operationsValidationCenterQuality struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	WhatThisMeans string `json:"what_this_means"`
	NextAction    string `json:"next_action"`
	DoesNotProve  string `json:"does_not_prove"`
	DetailsURL    string `json:"details_url"`
}

type operationsValidationCenterIssue struct {
	ID                string   `json:"id"`
	Source            string   `json:"source"`
	SourceLabel       string   `json:"source_label"`
	Status            string   `json:"status"`
	Severity          string   `json:"severity"`
	Family            string   `json:"family"`
	Codes             []string `json:"codes"`
	Count             int      `json:"count"`
	SampleCount       int      `json:"sample_count"`
	OverflowCount     int      `json:"overflow_count"`
	LikelyOwner       string   `json:"likely_owner"`
	AffectedFiles     string   `json:"affected_files"`
	OperatorSummary   string   `json:"operator_summary"`
	WhyItMatters      string   `json:"why_it_matters"`
	RecommendedAction string   `json:"recommended_action"`
	SafeFixPath       string   `json:"safe_fix_path"`
	VerifyWith        string   `json:"verify_with"`
	EscalateIf        string   `json:"escalate_if"`
	DetailsURL        string   `json:"details_url"`
	DoesNotProve      string   `json:"does_not_prove"`
}

type operationsValidationCenterTimeline struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	WhatThisMeans string `json:"what_this_means"`
	NextAction    string `json:"next_action"`
	DoesNotProve  string `json:"does_not_prove"`
}

type operationsValidationCenterBlocker struct {
	ID           string `json:"id"`
	Severity     string `json:"severity"`
	Area         string `json:"area"`
	Signal       string `json:"signal"`
	NextAction   string `json:"next_action"`
	DoesNotProve string `json:"does_not_prove"`
	ReviewURL    string `json:"review_url,omitempty"`
}

type operationsValidationCenterConsumer struct {
	Target       string `json:"target"`
	Status       string `json:"status"`
	Source       string `json:"source"`
	UpdatedAt    string `json:"updated_at"`
	NextAction   string `json:"next_action"`
	DoesNotProve string `json:"does_not_prove"`
}

type operationsValidationCenterCounts struct {
	FeedRows       int            `json:"feed_rows"`
	ValidationRows int            `json:"validation_rows"`
	IssueRows      int            `json:"issue_rows"`
	ConsumerRows   int            `json:"consumer_rows"`
	Statuses       map[string]int `json:"statuses"`
}

type operationsValidationCenterClaimFlags struct {
	ExternalEvidenceCreated      bool `json:"external_evidence_created"`
	FinalRootEvidenceCreated     bool `json:"final_root_evidence_created"`
	ConsumerStatusesChanged      bool `json:"consumer_statuses_changed"`
	ComplianceClaimed            bool `json:"compliance_claimed"`
	ProductionReadinessClaimed   bool `json:"production_readiness_claimed"`
	AgencyApprovalClaimed        bool `json:"agency_approval_claimed"`
	ConsumerAcceptanceClaimed    bool `json:"consumer_acceptance_claimed"`
	PublicLaunchClaimed          bool `json:"public_launch_claimed"`
	HostedSaaSClaimed            bool `json:"hosted_saas_claimed"`
	SLAClaimed                   bool `json:"sla_claimed"`
	UptimeGuaranteeClaimed       bool `json:"uptime_guarantee_claimed"`
	VendorCompatibilityClaimed   bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed bool `json:"hardware_certification_claimed"`
	ProductionGradeETAClaimed    bool `json:"production_grade_eta_claimed"`
}

func (h *handler) renderValidationCenter(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "validation-center")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "validation-center", page)
}

func (h *handler) renderValidationCenterJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "validation-center")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.ValidationCenter)
}

func buildOperationsValidationCenter(page operationsPage) operationsValidationCenterView {
	feedRows := validationCenterFeedRows(page.FeedHealth.Rows)
	validationRows := validationCenterValidationRows(page.ValidationHealth.Feeds)
	qualityRows := validationCenterQualityRows(page)
	issueRows := validationCenterIssueRows(page)
	timelineRows := validationCenterTimelineRows(page.ReadinessV2.Rows)
	blockerRows := validationCenterBlockerRows(timelineRows, feedRows, validationRows, issueRows)
	consumerRows := validationCenterConsumerRows(page.Consumers)
	return operationsValidationCenterView{
		GeneratedAt:       page.GeneratedAt,
		AgencyID:          page.AgencyID,
		Boundary:          "Private validation diagnostics only. This center combines existing private feed, validator, GTFS quality, readiness, reliability, and consumer prepared-tracker signals; it creates no evidence, changes no consumer status, contacts no external party, and proves no compliance, acceptance, final-root readiness, hosted operation, SLA, production readiness, vendor compatibility, hardware certification, public launch, or ETA quality.",
		FeedRows:          feedRows,
		ValidationHistory: validationRows,
		ValidatorHealth:   validationRows,
		GTFSQuality:       qualityRows,
		IssueDrilldowns:   issueRows,
		ReadinessTimeline: timelineRows,
		Blockers:          blockerRows,
		ConsumerTracker:   consumerRows,
		Counts:            validationCenterCounts(feedRows, validationRows, issueRows, consumerRows),
		ClaimFlags:        operationsValidationCenterClaimFlags{},
	}
}

func validationCenterFeedRows(rows []operationsFeedHealthRow) []operationsValidationCenterFeedRow {
	out := make([]operationsValidationCenterFeedRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, operationsValidationCenterFeedRow{
			ID:             row.ID,
			Label:          row.Label,
			PublicPath:     row.PublicPath,
			ConfiguredURL:  row.ConfiguredURL,
			Status:         readinessV2NormalizeStatus(row.Status),
			LastChecked:    row.LastChecked,
			HTTPStatus:     row.LastKnownHTTPStatus,
			ContentType:    row.ContentType,
			Freshness:      row.Freshness,
			ValidatorState: row.ValidatorState,
			HealthState:    row.HealthState,
			CurrentSignal:  row.CurrentSignal,
			WhatThisMeans:  row.WhatThisMeans,
			NextAction:     row.NextAction,
			DoesNotProve:   row.DoesNotProve,
		})
	}
	return out
}

func validationCenterValidationRows(rows []compliance.ValidationHealthRow) []operationsValidationCenterValidation {
	out := make([]operationsValidationCenterValidation, 0, len(rows))
	for _, row := range rows {
		status := readinessV2NormalizeStatus(row.HealthStatus)
		out = append(out, operationsValidationCenterValidation{
			ID:                        "validation_" + row.FeedType,
			FeedType:                  row.FeedType,
			Label:                     validationCenterFeedLabel(row.FeedType),
			ValidatorID:               row.ValidatorID,
			ValidatorName:             row.ValidatorName,
			Status:                    status,
			ToolingStatus:             row.ToolingStatus,
			ArtifactStatus:            row.ArtifactStatus,
			LatestResultStatus:        row.LatestResultStatus,
			LatestResultAt:            formatTimeForText(row.LatestResultAt),
			ActiveFeedVersionID:       row.ActiveFeedVersionID,
			LatestResultFeedVersionID: row.LatestResultFeedVersionID,
			StaleStatus:               row.StaleStatus,
			HealthStatus:              row.HealthStatus,
			CurrentSignal:             fmt.Sprintf("validator=%s; result=%s; stale=%s; tooling=%s; artifact=%s", row.ValidatorName, row.LatestResultStatus, row.StaleStatus, row.ToolingStatus, row.ArtifactStatus),
			WhatThisMeans:             validationCenterValidationMeaning(row),
			NextAction:                row.NextAction,
			DoesNotProve:              "Validator rows are supporting diagnostics only. They do not prove compliance, consumer acceptance, final-root readiness, hosted availability, SLA, production readiness, public launch, or agency approval.",
		})
	}
	return out
}

func validationCenterQualityRows(page operationsPage) []operationsValidationCenterQuality {
	rows := []operationsValidationCenterQuality{
		validationCenterQualityRow("canonical_static", "Canonical static validator", page.GTFSQuality.Canonical),
		validationCenterQualityRow("internal_importer", "Internal import validation", page.GTFSQuality.InternalImporter),
	}
	return rows
}

func validationCenterQualityRow(id string, label string, section compliance.GTFSQualitySection) operationsValidationCenterQuality {
	status := readinessV2NormalizeStatus(section.Status)
	signal := section.OperatorSummary
	if strings.TrimSpace(signal) == "" {
		signal = "no current quality summary is available"
	}
	return operationsValidationCenterQuality{
		ID:            id,
		Label:         label,
		Status:        status,
		CurrentSignal: signal,
		WhatThisMeans: "GTFS quality summarizes private validator or importer notices into operator review guidance without editing schedule data.",
		NextAction:    firstNonEmpty(section.RecommendedAction, "Open GTFS Quality to review the private fix planner and source-specific guidance."),
		DoesNotProve:  "GTFS quality summaries do not prove validator-clean production data, compliance, consumer acceptance, agency approval, or production readiness.",
		DetailsURL:    "/admin/operations/gtfs-quality",
	}
}

func validationCenterIssueRows(page operationsPage) []operationsValidationCenterIssue {
	rows := []operationsValidationCenterIssue{}
	rows = append(rows, validationCenterIssuesFromSection("canonical_static", page.GTFSQuality.Canonical)...)
	rows = append(rows, validationCenterIssuesFromSection("internal_importer", page.GTFSQuality.InternalImporter)...)
	return rows
}

func validationCenterIssuesFromSection(prefix string, section compliance.GTFSQualitySection) []operationsValidationCenterIssue {
	rows := make([]operationsValidationCenterIssue, 0, len(section.Groups))
	for index, group := range section.Groups {
		rows = append(rows, validationCenterIssueFromGroup(prefix, section, group, index))
	}
	if len(rows) == 0 && validationCenterSectionNeedsIssue(section.Status) {
		rows = append(rows, operationsValidationCenterIssue{
			ID:                prefix + "_section_status",
			Source:            section.Source,
			SourceLabel:       section.SourceLabel,
			Status:            readinessV2NormalizeStatus(section.Status),
			Severity:          readinessV2NormalizeStatus(section.Status),
			Family:            "section_status",
			Codes:             []string{"none recorded"},
			Count:             0,
			SampleCount:       0,
			OverflowCount:     section.OverflowCount,
			LikelyOwner:       "GTFS source owner with technical maintainer review",
			AffectedFiles:     "not available from the current sanitized summary",
			OperatorSummary:   firstNonEmpty(section.OperatorSummary, "No grouped issue rows are available for this source."),
			WhyItMatters:      "A non-ok source summary without grouped issues still needs operator review before stronger feed-readiness language.",
			RecommendedAction: firstNonEmpty(section.RecommendedAction, "Open GTFS Quality to review the private fix planner and source-specific guidance."),
			SafeFixPath:       "Review the source result in GTFS Quality, then fix source GTFS or rerun validation through the existing safe workflow.",
			VerifyWith:        "Rerun or refresh the appropriate validator and return to the private Validation Center.",
			EscalateIf:        "Escalate if this status blocks schedule or realtime review and grouped issue context is unavailable.",
			DetailsURL:        "/admin/operations/gtfs-quality",
			DoesNotProve:      validationCenterIssueBoundary(),
		})
	}
	return rows
}

func validationCenterIssueFromGroup(prefix string, section compliance.GTFSQualitySection, group compliance.GTFSQualityGroup, index int) operationsValidationCenterIssue {
	return operationsValidationCenterIssue{
		ID:                fmt.Sprintf("%s_%03d_%s", prefix, index+1, strings.ReplaceAll(group.Family, "_", "-")),
		Source:            firstNonEmpty(group.Source, section.Source),
		SourceLabel:       section.SourceLabel,
		Status:            readinessV2NormalizeStatus(group.Severity),
		Severity:          group.Severity,
		Family:            group.Family,
		Codes:             validationCenterSafeCodes(group.Codes),
		Count:             group.Count,
		SampleCount:       len(group.Samples),
		OverflowCount:     group.OverflowCount,
		LikelyOwner:       gtfsQualityLikelyOwner(group),
		AffectedFiles:     gtfsQualityAffectedFiles(group),
		OperatorSummary:   group.OperatorSummary,
		WhyItMatters:      group.WhyItMatters,
		RecommendedAction: group.RecommendedAction,
		SafeFixPath:       gtfsQualitySafeFixPath(section.Source, group),
		VerifyWith:        gtfsQualityVerifyWith(section.Source, group),
		EscalateIf:        gtfsQualityEscalation(group),
		DetailsURL:        "/admin/operations/gtfs-quality",
		DoesNotProve:      validationCenterIssueBoundary(),
	}
}

func validationCenterSectionNeedsIssue(status string) bool {
	switch readinessV2NormalizeStatus(status) {
	case checklistStatusBlocked, checklistStatusNeedsReview, checklistStatusUnknown:
		return true
	default:
		return false
	}
}

func validationCenterSafeCodes(codes []string) []string {
	out := make([]string, 0, len(codes))
	for _, code := range codes {
		out = append(out, validationCenterSafeCode(code))
	}
	if len(out) == 0 {
		return []string{"unknown"}
	}
	return out
}

func validationCenterSafeCode(code string) string {
	value := strings.TrimSpace(code)
	if value == "" {
		return "unknown"
	}
	lower := strings.ToLower(value)
	if len(value) > 80 || strings.Contains(lower, "/") || strings.Contains(lower, "\\") || strings.Contains(lower, "token") || strings.Contains(lower, "secret") || strings.Contains(lower, "password") || strings.Contains(lower, "authorization") || strings.Contains(lower, "bearer") || strings.Contains(lower, "cookie") || strings.Contains(lower, "postgres") || strings.Contains(lower, "http://") || strings.Contains(lower, "https://") {
		return "redacted"
	}
	return value
}

func validationCenterIssueBoundary() string {
	return "Issue drilldowns are sanitized operator guidance only. They do not edit GTFS, change drafts, publish schedules, prove compliance, prove consumer acceptance, prove agency approval, or prove production readiness."
}

func validationCenterTimelineRows(rows []operationsReadinessV2Row) []operationsValidationCenterTimeline {
	out := make([]operationsValidationCenterTimeline, 0, len(rows))
	for _, row := range rows {
		out = append(out, operationsValidationCenterTimeline{
			ID:            row.ID,
			Label:         row.ReadinessItem,
			Status:        readinessV2NormalizeStatus(row.Status),
			CurrentSignal: row.CurrentSignal,
			WhatThisMeans: row.WhatThisMeans,
			NextAction:    row.WhatToDoNext,
			DoesNotProve:  row.WhatItDoesNotProve,
		})
	}
	return out
}

func validationCenterBlockerRows(timelineRows []operationsValidationCenterTimeline, feedRows []operationsValidationCenterFeedRow, validationRows []operationsValidationCenterValidation, issueRows []operationsValidationCenterIssue) []operationsValidationCenterBlocker {
	rows := []operationsValidationCenterBlocker{}
	for _, row := range timelineRows {
		if !validationCenterNeedsAction(row.Status) {
			continue
		}
		rows = append(rows, operationsValidationCenterBlocker{
			ID:           "readiness_" + row.ID,
			Severity:     readinessV2NormalizeStatus(row.Status),
			Area:         row.Label,
			Signal:       row.CurrentSignal,
			NextAction:   row.NextAction,
			DoesNotProve: row.DoesNotProve,
			ReviewURL:    validationCenterReviewURL(row.ID),
		})
	}
	for _, row := range feedRows {
		if !validationCenterNeedsAction(row.Status) {
			continue
		}
		rows = append(rows, operationsValidationCenterBlocker{
			ID:           "feed_" + row.ID,
			Severity:     readinessV2NormalizeStatus(row.Status),
			Area:         row.Label,
			Signal:       row.CurrentSignal,
			NextAction:   row.NextAction,
			DoesNotProve: row.DoesNotProve,
			ReviewURL:    validationCenterFeedReviewURL(row.ID),
		})
	}
	for _, row := range validationRows {
		if !validationCenterNeedsAction(row.Status) {
			continue
		}
		rows = append(rows, operationsValidationCenterBlocker{
			ID:           "validator_" + row.FeedType,
			Severity:     readinessV2NormalizeStatus(row.Status),
			Area:         row.Label + " validation",
			Signal:       row.CurrentSignal,
			NextAction:   row.NextAction,
			DoesNotProve: row.DoesNotProve,
			ReviewURL:    "/admin/operations/validation-health",
		})
	}
	for _, row := range issueRows {
		if !validationCenterNeedsAction(row.Status) {
			continue
		}
		rows = append(rows, operationsValidationCenterBlocker{
			ID:           "issue_" + row.ID,
			Severity:     readinessV2NormalizeStatus(row.Status),
			Area:         row.SourceLabel + " / " + row.Family,
			Signal:       row.OperatorSummary,
			NextAction:   row.RecommendedAction,
			DoesNotProve: row.DoesNotProve,
			ReviewURL:    row.DetailsURL,
		})
	}
	if len(rows) > 30 {
		overflow := len(rows) - 30
		rows = rows[:30]
		rows = append(rows, operationsValidationCenterBlocker{
			ID:           "blocker_overflow",
			Severity:     checklistStatusNeedsReview,
			Area:         "Additional private blocker rows",
			Signal:       fmt.Sprintf("%d additional rows are hidden by the Center cap", overflow),
			NextAction:   "Open Feed Health, Validator Health, GTFS Quality, and Readiness to review the complete private row set.",
			DoesNotProve: "Hidden row count does not prove release readiness or absence of other issues.",
			ReviewURL:    "/admin/operations/readiness",
		})
	}
	return rows
}

func validationCenterNeedsAction(status string) bool {
	switch readinessV2NormalizeStatus(status) {
	case checklistStatusBlocked, checklistStatusMissing, checklistStatusNeedsReview:
		return true
	default:
		return false
	}
}

func validationCenterReviewURL(id string) string {
	switch id {
	case "discovery_metadata":
		return "/admin/operations/setup"
	case "feed_health", "vehicle_positions", "trip_updates", "alerts":
		return "/admin/operations/feed-health"
	case "static_gtfs_quality":
		return "/admin/operations/gtfs-quality"
	case "validation_health":
		return "/admin/operations/validation-health"
	case "operations_reliability":
		return "/admin/operations/reliability"
	case "telemetry_devices":
		return "/admin/operations/telemetry"
	case "operations_scorecard", "consumer_prepared_tracker":
		return "/admin/operations/readiness"
	default:
		return "/admin/operations/readiness"
	}
}

func validationCenterFeedReviewURL(id string) string {
	switch id {
	case "feeds_json":
		return "/admin/operations/feeds"
	case "schedule":
		return "/admin/operations/gtfs-quality"
	case "vehicle_positions", "trip_updates", "alerts":
		return "/admin/operations/realtime"
	default:
		return "/admin/operations/feed-health"
	}
}

func validationCenterConsumerRows(rows []consumerStatusView) []operationsValidationCenterConsumer {
	out := make([]operationsValidationCenterConsumer, 0, len(rows))
	for _, row := range rows {
		out = append(out, operationsValidationCenterConsumer{
			Target:       row.Name,
			Status:       row.Status,
			Source:       row.Source,
			UpdatedAt:    formatTimeForText(row.UpdatedAt),
			NextAction:   "Keep the target-specific packet record prepared-only until separate written authorization and retained target-originated evidence support a status change.",
			DoesNotProve: "Prepared tracker status does not prove submission, review, acceptance, ingestion, listing, display, consumer approval, or compliance.",
		})
	}
	return out
}

func validationCenterCounts(feedRows []operationsValidationCenterFeedRow, validationRows []operationsValidationCenterValidation, issueRows []operationsValidationCenterIssue, consumerRows []operationsValidationCenterConsumer) operationsValidationCenterCounts {
	statuses := map[string]int{
		checklistStatusOK:          0,
		checklistStatusNeedsReview: 0,
		checklistStatusMissing:     0,
		checklistStatusBlocked:     0,
		checklistStatusUnknown:     0,
	}
	for _, row := range feedRows {
		statuses[readinessV2NormalizeStatus(row.Status)]++
	}
	for _, row := range validationRows {
		statuses[readinessV2NormalizeStatus(row.Status)]++
	}
	for _, row := range issueRows {
		statuses[readinessV2NormalizeStatus(row.Status)]++
	}
	return operationsValidationCenterCounts{
		FeedRows:       len(feedRows),
		ValidationRows: len(validationRows),
		IssueRows:      len(issueRows),
		ConsumerRows:   len(consumerRows),
		Statuses:       statuses,
	}
}

func validationCenterFeedLabel(feedType string) string {
	switch feedType {
	case "schedule":
		return "Static GTFS Schedule"
	case "vehicle_positions":
		return "Vehicle Positions"
	case "trip_updates":
		return "Trip Updates"
	case "alerts":
		return "Alerts"
	default:
		return firstNonEmpty(feedType, "unknown feed")
	}
}

func validationCenterValidationMeaning(row compliance.ValidationHealthRow) string {
	switch row.FeedType {
	case "schedule":
		return "Static validation summarizes the active schedule artifact using server-owned validator tooling when available."
	case "vehicle_positions", "trip_updates", "alerts":
		return "Realtime validation summarizes server-owned GTFS-Realtime protobuf artifacts without exposing raw validator output."
	default:
		return "Validation health summarizes server-owned tooling, artifact availability, latest result, and staleness."
	}
}
