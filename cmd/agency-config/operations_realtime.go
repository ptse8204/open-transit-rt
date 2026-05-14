package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/state"
)

type operationsRealtimeView struct {
	GeneratedAt    time.Time                      `json:"generated_at"`
	AgencyID       string                         `json:"agency_id"`
	Boundary       string                         `json:"boundary"`
	StaleThreshold string                         `json:"stale_threshold"`
	Summary        operationsRealtimeSummary      `json:"summary"`
	Feeds          []operationsRealtimeFeedStatus `json:"feeds"`
	Fleet          []operationsRealtimeFleetRow   `json:"fleet"`
	ClaimFlags     operationsRealtimeClaimFlags   `json:"claim_flags"`
}

type operationsRealtimeSummary struct {
	Status              string `json:"status"`
	CurrentSignal       string `json:"current_signal"`
	LatestTelemetryRows int    `json:"latest_telemetry_rows"`
	FreshTelemetryRows  int    `json:"fresh_telemetry_rows"`
	StaleTelemetryRows  int    `json:"stale_telemetry_rows"`
	DeviceBindings      int    `json:"device_bindings"`
	DevicesReporting    int    `json:"devices_reporting"`
	DevicesNotSeen      int    `json:"devices_not_seen"`
	MatchedAssignments  int    `json:"matched_assignments"`
	UnknownAssignments  int    `json:"unknown_assignments"`
	LowConfidenceRows   int    `json:"low_confidence_rows"`
	ManualOverrides     int    `json:"manual_overrides"`
	NextAction          string `json:"next_action"`
	DoesNotProve        string `json:"does_not_prove"`
}

type operationsRealtimeFeedStatus struct {
	ID              string      `json:"id"`
	Label           string      `json:"label"`
	State           string      `json:"state"`
	Count           string      `json:"count"`
	LatestSignal    string      `json:"latest_signal"`
	StaleOrWithheld string      `json:"stale_or_withheld"`
	Adapter         string      `json:"adapter,omitempty"`
	NextAction      string      `json:"next_action"`
	AdminLink       string      `json:"admin_link"`
	DoesNotProve    string      `json:"does_not_prove"`
	Details         []countView `json:"details,omitempty"`
}

type operationsRealtimeFleetRow struct {
	VehicleID        string   `json:"vehicle_id"`
	DeviceID         string   `json:"device_id"`
	Freshness        string   `json:"freshness"`
	ObservedAt       string   `json:"observed_at"`
	ReceivedAt       string   `json:"received_at"`
	AgeSeconds       int64    `json:"age_seconds"`
	AssignmentState  string   `json:"assignment_state"`
	DegradedState    string   `json:"degraded_state,omitempty"`
	AssignmentSource string   `json:"assignment_source,omitempty"`
	RouteID          string   `json:"route_id,omitempty"`
	TripID           string   `json:"trip_id,omitempty"`
	Confidence       string   `json:"confidence,omitempty"`
	ReasonCodes      []string `json:"reason_codes,omitempty"`
	CurrentSignal    string   `json:"current_signal"`
	NextAction       string   `json:"next_action"`
	DoesNotProve     string   `json:"does_not_prove"`
}

type operationsRealtimeClaimFlags struct {
	BrowserTelemetrySendEnabled     bool `json:"browser_telemetry_send_enabled"`
	BackendCommandExecutionEnabled  bool `json:"backend_command_execution_enabled"`
	DeviceTokenCollectedByBrowser   bool `json:"device_token_collected_by_browser"`
	ExternalEvidenceCreated         bool `json:"external_evidence_created"`
	ConsumerStatusesChanged         bool `json:"consumer_statuses_changed"`
	ComplianceClaimed               bool `json:"compliance_claimed"`
	ProductionReadinessClaimed      bool `json:"production_readiness_claimed"`
	VendorCompatibilityClaimed      bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed    bool `json:"hardware_certification_claimed"`
	ProductionAVLReliabilityClaimed bool `json:"production_avl_reliability_claimed"`
	ProductionGradeETAClaimed       bool `json:"production_grade_eta_claimed"`
	RealWorldETAAccuracyClaimed     bool `json:"real_world_eta_accuracy_claimed"`
	SLAClaimed                      bool `json:"sla_claimed"`
	PublicLaunchClaimed             bool `json:"public_launch_claimed"`
	ConsumerAcceptanceClaimed       bool `json:"consumer_acceptance_claimed"`
}

func (h *handler) renderRealtime(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "realtime")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "realtime", page)
}

func (h *handler) renderRealtimeJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "realtime")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.Realtime)
}

func buildOperationsRealtime(page operationsPage) operationsRealtimeView {
	fleet := buildRealtimeFleetRows(page)
	summary := buildRealtimeSummary(page, fleet)
	return operationsRealtimeView{
		GeneratedAt:    page.GeneratedAt,
		AgencyID:       page.AgencyID,
		Boundary:       "Private authenticated realtime operations diagnostics only. Viewing this page sends no telemetry, runs no backend command, creates no retained evidence, changes no consumer status, contacts no external party, and records no compliance, consumer-acceptance, vendor, hardware, SLA, public-launch, production-readiness, production AVL reliability, or production-grade ETA claim.",
		StaleThreshold: page.StaleThreshold.String(),
		Summary:        summary,
		Feeds:          realtimeFeedStatuses(page),
		Fleet:          fleet,
		ClaimFlags:     operationsRealtimeClaimFlags{},
	}
}

func buildRealtimeSummary(page operationsPage, fleet []operationsRealtimeFleetRow) operationsRealtimeSummary {
	summary := operationsRealtimeSummary{
		Status:              checklistStatusUnknown,
		LatestTelemetryRows: len(page.Telemetry),
		StaleTelemetryRows:  page.StaleCount,
		DeviceBindings:      len(page.Devices),
		DoesNotProve:        "Does not prove production AVL reliability, vendor compatibility, hardware certification, consumer display, compliance, hosted operation, SLA, public launch, production readiness, or production-grade ETA quality.",
	}
	if summary.LatestTelemetryRows >= summary.StaleTelemetryRows {
		summary.FreshTelemetryRows = summary.LatestTelemetryRows - summary.StaleTelemetryRows
	}
	for _, row := range fleet {
		switch row.Freshness {
		case "fresh":
			summary.DevicesReporting++
		case "not seen":
			summary.DevicesNotSeen++
		}
		if row.AssignmentSource == string(state.AssignmentSourceManualOverride) {
			summary.ManualOverrides++
		}
		if isRealtimeMatched(row) {
			summary.MatchedAssignments++
		} else {
			summary.UnknownAssignments++
		}
		if isRealtimeLowConfidence(row.Confidence) {
			summary.LowConfidenceRows++
		}
	}
	if page.TelemetryError != "" {
		summary.Status = checklistStatusBlocked
		summary.CurrentSignal = page.TelemetryError
		summary.NextAction = "Confirm the telemetry service and database are available, then return to this private center."
		return summary
	}
	switch {
	case summary.LatestTelemetryRows == 0 && summary.DeviceBindings == 0:
		summary.Status = checklistStatusMissing
		summary.CurrentSignal = "no latest telemetry rows and no device bindings are visible to this console"
		summary.NextAction = "Create or rotate a device credential, install it on a device or simulator, then send an authenticated sample telemetry event."
	case summary.LatestTelemetryRows == 0:
		summary.Status = checklistStatusMissing
		summary.CurrentSignal = fmt.Sprintf("%d device bindings exist, but no latest telemetry rows are visible", summary.DeviceBindings)
		summary.NextAction = "Check device installation or run the synthetic simulator from an operator shell."
	case summary.StaleTelemetryRows == summary.LatestTelemetryRows:
		summary.Status = checklistStatusNeedsReview
		summary.CurrentSignal = fmt.Sprintf("%d latest telemetry rows are stale at threshold %s", summary.StaleTelemetryRows, page.StaleThreshold)
		summary.NextAction = "Check device power, network, reporting cadence, and simulator or ingest logs before relying on Vehicle Positions."
	case summary.UnknownAssignments > 0 || summary.LowConfidenceRows > 0:
		summary.Status = checklistStatusNeedsReview
		summary.CurrentSignal = fmt.Sprintf("%d fresh rows, %d stale rows, %d unknown or unavailable assignments, %d low-confidence rows", summary.FreshTelemetryRows, summary.StaleTelemetryRows, summary.UnknownAssignments, summary.LowConfidenceRows)
		summary.NextAction = "Review assignment reasons and keep trip descriptors unknown when confidence is not sufficient."
	default:
		summary.Status = checklistStatusOK
		summary.CurrentSignal = fmt.Sprintf("%d latest telemetry rows are visible; %d are fresh; %d assignments are matched", summary.LatestTelemetryRows, summary.FreshTelemetryRows, summary.MatchedAssignments)
		summary.NextAction = "Continue monitoring freshness, assignment confidence, Vehicle Positions, Trip Updates diagnostics, and Alerts lifecycle."
	}
	return summary
}

func buildRealtimeFleetRows(page operationsPage) []operationsRealtimeFleetRow {
	rows := make([]operationsRealtimeFleetRow, 0, len(page.Telemetry)+len(page.DeviceRows))
	seenDeviceVehicle := map[string]bool{}
	for _, telemetryRow := range page.Telemetry {
		row := realtimeFleetRowFromTelemetry(telemetryRow)
		rows = append(rows, row)
		seenDeviceVehicle[deviceVehicleKey(row.DeviceID, row.VehicleID)] = true
	}
	for _, deviceRow := range page.DeviceRows {
		key := deviceVehicleKey(deviceRow.DeviceID, deviceRow.VehicleID)
		if seenDeviceVehicle[key] {
			continue
		}
		rows = append(rows, realtimeFleetRowFromDevice(deviceRow))
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].VehicleID == rows[j].VehicleID {
			return rows[i].DeviceID < rows[j].DeviceID
		}
		return rows[i].VehicleID < rows[j].VehicleID
	})
	return rows
}

func realtimeFleetRowFromTelemetry(row telemetryView) operationsRealtimeFleetRow {
	freshness := "fresh"
	if row.Stale {
		freshness = "stale"
	}
	out := operationsRealtimeFleetRow{
		VehicleID:        row.VehicleID,
		DeviceID:         row.DeviceID,
		Freshness:        freshness,
		ObservedAt:       formatTimeForText(&row.ObservedAt),
		ReceivedAt:       formatTimeForText(&row.ReceivedAt),
		AgeSeconds:       row.AgeSeconds,
		AssignmentState:  firstNonEmpty(row.AssignmentState, string(state.StateUnknown)),
		DegradedState:    row.DegradedState,
		AssignmentSource: row.AssignmentSource,
		RouteID:          row.RouteID,
		TripID:           row.TripID,
		Confidence:       row.Confidence,
		ReasonCodes:      append([]string(nil), row.ReasonCodes...),
		DoesNotProve:     "Does not prove real fleet reliability, vendor compatibility, hardware certification, consumer display, compliance, or production-grade ETA quality.",
	}
	out.CurrentSignal = realtimeFleetSignal(out)
	out.NextAction = realtimeFleetNextAction(out)
	return out
}

func realtimeFleetRowFromDevice(row operationsDeviceRow) operationsRealtimeFleetRow {
	out := operationsRealtimeFleetRow{
		VehicleID:        row.VehicleID,
		DeviceID:         row.DeviceID,
		Freshness:        row.Freshness,
		ObservedAt:       formatTimeForText(row.LatestObservedAt),
		ReceivedAt:       formatTimeForText(row.LatestReceivedAt),
		AssignmentState:  "not available",
		AssignmentSource: row.AssignmentSource,
		DoesNotProve:     "A device binding does not prove hardware certification, vendor compatibility, production AVL reliability, consumer display, compliance, or production readiness.",
	}
	if row.LatestAgeSeconds != nil {
		out.AgeSeconds = *row.LatestAgeSeconds
	}
	out.CurrentSignal = realtimeFleetSignal(out)
	out.NextAction = realtimeFleetNextAction(out)
	return out
}

func realtimeFeedStatuses(page operationsPage) []operationsRealtimeFeedStatus {
	source := []operationsRealtimeUsefulnessRow{
		page.FeedHealth.RealtimeUsefulness.VehiclePositions,
		page.FeedHealth.RealtimeUsefulness.TripUpdates,
		page.FeedHealth.RealtimeUsefulness.Alerts,
	}
	rows := make([]operationsRealtimeFeedStatus, 0, len(source))
	for _, row := range source {
		rows = append(rows, operationsRealtimeFeedStatus{
			ID:              row.ID,
			Label:           row.Label,
			State:           row.State,
			Count:           row.Count,
			LatestSignal:    row.LatestSignal,
			StaleOrWithheld: row.StaleOrHeld,
			Adapter:         row.Adapter,
			NextAction:      row.NextAction,
			AdminLink:       row.AdminLink,
			DoesNotProve:    row.DoesNotProve,
			Details:         row.Details,
		})
	}
	return rows
}

func realtimeFleetSignal(row operationsRealtimeFleetRow) string {
	if row.Freshness == "not seen" {
		return "device binding is visible, but no latest accepted telemetry row is visible"
	}
	parts := []string{row.Freshness + " telemetry"}
	if row.ObservedAt != "" && row.ObservedAt != "not available" {
		parts = append(parts, "observed "+row.ObservedAt)
	}
	if row.AssignmentState != "" && row.AssignmentState != "not available" {
		parts = append(parts, "assignment "+row.AssignmentState)
	}
	if row.DegradedState != "" && row.DegradedState != string(state.DegradedNone) {
		parts = append(parts, "degraded "+row.DegradedState)
	}
	if row.Confidence != "" {
		parts = append(parts, "confidence "+row.Confidence)
	}
	return strings.Join(parts, "; ")
}

func realtimeFleetNextAction(row operationsRealtimeFleetRow) string {
	if row.Freshness == "not seen" {
		return "Install or check the device credential, then send an authenticated sample telemetry event from the device or simulator."
	}
	if row.Freshness == "stale" {
		return "Check device power, network, and reporting cadence before changing assignments."
	}
	if row.AssignmentSource == string(state.AssignmentSourceManualOverride) {
		return "Review the active operator override; automatic matching should not replace it until it expires or is cleared."
	}
	if !isRealtimeMatched(row) {
		return "Keep the trip descriptor unknown until matching confidence improves or an operator override is applied."
	}
	if isRealtimeLowConfidence(row.Confidence) {
		return "Review low-confidence matching evidence before relying on route or trip context."
	}
	return "Continue monitoring freshness and assignment confidence."
}

func isRealtimeMatched(row operationsRealtimeFleetRow) bool {
	if row.AssignmentState == "" || row.AssignmentState == "not available" || row.AssignmentState == string(state.StateUnknown) {
		return false
	}
	if row.DegradedState != "" && row.DegradedState != string(state.DegradedNone) {
		return false
	}
	if isRealtimeLowConfidence(row.Confidence) {
		return false
	}
	return row.TripID != "" || row.AssignmentSource == string(state.AssignmentSourceManualOverride)
}

func isRealtimeLowConfidence(value string) bool {
	if strings.TrimSpace(value) == "" {
		return false
	}
	confidence, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}
	return confidence < 0.70
}
