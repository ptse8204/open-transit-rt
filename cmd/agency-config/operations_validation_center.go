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
	consumerRows := validationCenterConsumerRows(page.Consumers)
	return operationsValidationCenterView{
		GeneratedAt:       page.GeneratedAt,
		AgencyID:          page.AgencyID,
		Boundary:          "Private validation diagnostics only. This center combines existing private feed, validator, GTFS quality, readiness, reliability, and consumer prepared-tracker signals; it creates no evidence, changes no consumer status, contacts no external party, and proves no compliance, acceptance, final-root readiness, hosted operation, SLA, production readiness, vendor compatibility, hardware certification, public launch, or ETA quality.",
		FeedRows:          feedRows,
		ValidationHistory: validationRows,
		ValidatorHealth:   validationRows,
		GTFSQuality:       qualityRows,
		ReadinessTimeline: []operationsValidationCenterTimeline{},
		Blockers:          []operationsValidationCenterBlocker{},
		ConsumerTracker:   consumerRows,
		Counts:            validationCenterCounts(feedRows, validationRows, consumerRows),
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
		NextAction:    firstNonEmpty(section.RecommendedAction, "Open GTFS Quality to review source-specific guidance."),
		DoesNotProve:  "GTFS quality summaries do not prove validator-clean production data, compliance, consumer acceptance, agency approval, or production readiness.",
		DetailsURL:    "/admin/operations/gtfs-quality",
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

func validationCenterCounts(feedRows []operationsValidationCenterFeedRow, validationRows []operationsValidationCenterValidation, consumerRows []operationsValidationCenterConsumer) operationsValidationCenterCounts {
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
	return operationsValidationCenterCounts{
		FeedRows:       len(feedRows),
		ValidationRows: len(validationRows),
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
