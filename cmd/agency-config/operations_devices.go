package main

import (
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/devices"
)

type operationsDeviceOnboardingUseCase struct {
	Name      string
	When      string
	NextStep  string
	AdminOnly bool
}

type operationsDeviceRow struct {
	DeviceID         string
	VehicleID        string
	Status           string
	ValidFrom        time.Time
	LastUsedAt       *time.Time
	RotatedAt        *time.Time
	RevokedAt        *time.Time
	LatestObservedAt *time.Time
	LatestReceivedAt *time.Time
	LatestAgeSeconds *int64
	Freshness        string
	Assignment       string
	AssignmentAt     *time.Time
	AssignmentSource string
	NextAction       string
}

type operationsDeviceFleetOnboardingView struct {
	Boundary             string
	Status               string
	InventoryRows        []operationsDeviceFleetOnboardingRow
	BulkImportRows       []operationsDeviceFleetOnboardingRow
	TokenLifecycleRows   []operationsDeviceFleetOnboardingRow
	FreshnessTriageRows  []operationsDeviceFleetOnboardingRow
	BindingReviewRows    []operationsDeviceFleetOnboardingRow
	TechnicalHandoffRows []operationsDeviceFleetOnboardingRow
}

type operationsDeviceFleetOnboardingRow struct {
	ID                  string
	Label               string
	Status              string
	CurrentSignal       string
	OperatorStep        string
	TechnicalHelperStep string
	DoesNotProve        string
}

func operationsDeviceOnboardingUseCases() []operationsDeviceOnboardingUseCase {
	return []operationsDeviceOnboardingUseCase{
		{
			Name:      "New vehicle",
			When:      "Use when a vehicle or phone tracker is entering service for the first time.",
			NextStep:  "Create or rotate a one-time token, install it on the device, then send a sample telemetry event.",
			AdminOnly: true,
		},
		{
			Name:      "Credential rotation",
			When:      "Use after staff turnover, suspected exposure, or scheduled key hygiene.",
			NextStep:  "Rotate the token and replace it on the installed device before the next pull-out.",
			AdminOnly: true,
		},
		{
			Name:      "Vehicle swap",
			When:      "Use when the physical tracker remains active but the assigned vehicle changes.",
			NextStep:  "Rebind the device to the correct vehicle, then confirm fresh telemetry appears for that vehicle.",
			AdminOnly: true,
		},
		{
			Name:      "Telemetry verification",
			When:      "Use after setup to confirm the device is producing accepted observations.",
			NextStep:  "Review freshness, assignment state, and the per-device next action below.",
			AdminOnly: false,
		},
	}
}

func buildOperationsDeviceFleetOnboarding(page operationsPage) operationsDeviceFleetOnboardingView {
	summary := summarizeOperationsDeviceFleet(page.DeviceRows, page.Telemetry)
	status := operationsStatusReady
	if summary.TotalBindings == 0 {
		status = operationsStatusMissing
	} else if summary.StaleBindings > 0 || summary.NotSeenBindings > 0 || summary.UnknownAssignments > 0 || summary.LowConfidenceAssignments > 0 || summary.UnboundTelemetryRows > 0 {
		status = operationsStatusNeedsReview
	}

	return operationsDeviceFleetOnboardingView{
		Boundary: "Private fleet onboarding guidance only. This page does not collect device token values, send telemetry, import bulk secrets, contact vendors, collect evidence, certify hardware, prove vendor compatibility, prove production AVL reliability, or change consumer status.",
		Status:   status,
		InventoryRows: []operationsDeviceFleetOnboardingRow{
			{
				ID:                  "device_inventory",
				Label:               "Vehicle / device inventory review",
				Status:              deviceInventoryStatus(summary),
				CurrentSignal:       deviceInventorySignal(summary),
				OperatorStep:        "Compare the listed device IDs and vehicle IDs with the agency's private fleet roster before the next pull-out.",
				TechnicalHelperStep: "If rows are missing, prepare public-safe device_id and vehicle_id pairs, then use the audited rotate/rebind flow one device at a time.",
				DoesNotProve:        "Does not prove the fleet roster is complete, hardware is certified, or a vendor mapping is correct.",
			},
			{
				ID:                  "telemetry_coverage",
				Label:               "Telemetry coverage by binding",
				Status:              deviceCoverageStatus(summary),
				CurrentSignal:       deviceCoverageSignal(summary),
				OperatorStep:        "Review fresh, stale, and not-seen counts before relying on Vehicle Positions for service review.",
				TechnicalHelperStep: "Use simulator dry runs or device logs outside the browser to confirm reporting cadence without exposing credentials.",
				DoesNotProve:        "Does not prove real-world AVL reliability, SLA, uptime, or consumer ingestion.",
			},
		},
		BulkImportRows: []operationsDeviceFleetOnboardingRow{
			{
				ID:                  "bulk_import_scope",
				Label:               "Bulk import plan",
				Status:              operationsStatusDiagnosticOnly,
				CurrentSignal:       "Bulk onboarding remains a private planning checklist; the console does not import token values or generate bulk secrets.",
				OperatorStep:        "Prepare a private CSV with agency_id, device_id, vehicle_id, intended_status, install_owner, and notes only.",
				TechnicalHelperStep: "Reject columns for tokens, auth headers, passwords, endpoint secrets, raw payloads, private serials, or vendor account IDs before any script uses the file.",
				DoesNotProve:        "Does not prove a bulk import ran, a device was installed, or vendor hardware works.",
			},
			{
				ID:                  "dry_run_preview",
				Label:               "Dry-run preview before install",
				Status:              operationsStatusDiagnosticOnly,
				CurrentSignal:       "Use local dry-run helpers before installing credentials on devices.",
				OperatorStep:        "Ask a technical helper for a redacted dry-run summary, then rotate/rebind real tokens only when ready to install.",
				TechnicalHelperStep: "Run `scripts/device-onboarding.sh sample --dry-run` or `scripts/device-onboarding.sh simulate --dry-run`; do not paste returned secrets into tickets or docs.",
				DoesNotProve:        "Does not prove live telemetry was sent or accepted.",
			},
		},
		TokenLifecycleRows: []operationsDeviceFleetOnboardingRow{
			{
				ID:                  "rotate_rebind_only",
				Label:               "Rotate/rebind is the supported credential action",
				Status:              operationsStatusReady,
				CurrentSignal:       "Existing device store supports verify, rotate/rebind, and metadata listing; token recovery is intentionally unavailable.",
				OperatorStep:        "Use rotate/rebind when a device is new, moved, replaced, exposed, or assigned to a corrected vehicle ID.",
				TechnicalHelperStep: "Store the one-time returned token in deployment-owned secret storage immediately; never rely on this console to recover it later.",
				DoesNotProve:        "Does not prove secret storage, physical installation, or device custody controls are complete.",
			},
			{
				ID:                  "secret_delivery",
				Label:               "Secret delivery checklist",
				Status:              operationsStatusNeedsReview,
				CurrentSignal:       "Secret delivery is deployment-owned and must happen outside public docs, screenshots, issue comments, and committed examples.",
				OperatorStep:        "Confirm who installs the token, where it is stored, and when old copied commands or screenshots are destroyed.",
				TechnicalHelperStep: "Use private secret channels and redact command history; never ask the browser page to collect device token values.",
				DoesNotProve:        "Does not prove security compliance, audit completeness, or absence of prior exposure.",
			},
		},
		FreshnessTriageRows: []operationsDeviceFleetOnboardingRow{
			{
				ID:                  "not_seen_devices",
				Label:               "Not-seen device triage",
				Status:              notSeenStatus(summary),
				CurrentSignal:       strconv.Itoa(summary.NotSeenBindings) + " configured bindings have no latest accepted telemetry.",
				OperatorStep:        "Confirm the device is installed, powered, assigned to the listed vehicle, and expected to report today.",
				TechnicalHelperStep: "Check local device logs and token binding inputs privately; do not store rejected raw payloads in public artifacts.",
				DoesNotProve:        "Does not prove the device is offline; it may not have attempted to report yet.",
			},
			{
				ID:                  "stale_devices",
				Label:               "Stale telemetry triage",
				Status:              staleStatus(summary),
				CurrentSignal:       strconv.Itoa(summary.StaleBindings) + " configured bindings have stale latest telemetry.",
				OperatorStep:        "Review whether stale devices are out of service, parked, in a coverage gap, or misconfigured.",
				TechnicalHelperStep: "Compare observed and received times, reporting cadence, and network path using private logs only.",
				DoesNotProve:        "Does not prove production AVL reliability or a consumer-visible outage.",
			},
			{
				ID:                  "unknown_device_triage",
				Label:               "Unknown-device and rejected-payload triage",
				Status:              operationsStatusDiagnosticOnly,
				CurrentSignal:       "Unauthorized device payloads are rejected before storage; use aggregate symptoms and private logs rather than retained raw payloads.",
				OperatorStep:        "If a device reports 401 errors, verify agency_id, device_id, vehicle_id, and current token binding without copying the token value.",
				TechnicalHelperStep: "Use adapter conformance fixtures and local logs to separate unknown-device, mismatched-binding, invalid-token, and stale-source cases.",
				DoesNotProve:        "Does not prove an unknown device exists in stored telemetry, because rejected payloads are not persisted.",
			},
		},
		BindingReviewRows: []operationsDeviceFleetOnboardingRow{
			{
				ID:                  "device_vehicle_binding",
				Label:               "Device-to-vehicle binding review",
				Status:              bindingReviewStatus(summary),
				CurrentSignal:       bindingReviewSignal(summary),
				OperatorStep:        "Confirm each device row is attached to the expected public-safe vehicle ID before using trip descriptors operationally.",
				TechnicalHelperStep: "Correct mapping errors with rotate/rebind, then confirm fresh accepted telemetry appears for the corrected vehicle.",
				DoesNotProve:        "Does not prove a GTFS vehicle mapping, block assignment, or trip match is correct.",
			},
			{
				ID:                  "identifier_hygiene",
				Label:               "Identifier hygiene",
				Status:              operationsStatusNeedsReview,
				CurrentSignal:       "Device IDs and vehicle IDs should be stable and public-safe; private serials and vendor account IDs belong in deployment-owned records.",
				OperatorStep:        "Review identifiers before screenshots, support summaries, or contributor docs are created.",
				TechnicalHelperStep: "Keep vendor-specific IDs inside adapters or private mapping files, not core matching logic or public examples.",
				DoesNotProve:        "Does not prove privacy review, vendor approval, or public evidence readiness.",
			},
		},
		TechnicalHandoffRows: []operationsDeviceFleetOnboardingRow{
			{
				ID:                  "helper_packet",
				Label:               "Safe technical-helper handoff",
				Status:              operationsStatusDiagnosticOnly,
				CurrentSignal:       "Safe handoff may include counts, public-safe IDs, freshness buckets, route links, command names, and environment variable names only.",
				OperatorStep:        "Share what is stale, missing, or mismatched without sending token values, request credential headers, private endpoints, raw telemetry, database URLs, or evidence packets.",
				TechnicalHelperStep: "Return a redacted summary of commands run, observed status codes, freshness counts, and next action; keep raw logs private.",
				DoesNotProve:        "Does not create retained evidence, consumer proof, or production support certification.",
			},
			{
				ID:                  "post_install_review",
				Label:               "Post-install review",
				Status:              postInstallReviewStatus(summary),
				CurrentSignal:       postInstallReviewSignal(summary),
				OperatorStep:        "After install, confirm accepted telemetry, freshness, and assignment confidence before depending on public trip descriptors.",
				TechnicalHelperStep: "Escalate low-confidence or unknown assignments to matcher review instead of forcing a trip descriptor.",
				DoesNotProve:        "Does not prove ETA quality, compliance, consumer acceptance, or production readiness.",
			},
		},
	}
}

func buildOperationsDeviceRows(bindings []devices.Binding, telemetryRows []telemetryView) []operationsDeviceRow {
	byDeviceVehicle := make(map[string]telemetryView, len(telemetryRows))
	for _, row := range telemetryRows {
		if row.DeviceID != "" || row.VehicleID != "" {
			byDeviceVehicle[deviceVehicleKey(row.DeviceID, row.VehicleID)] = row
		}
	}

	rows := make([]operationsDeviceRow, 0, len(bindings))
	for _, binding := range bindings {
		row := operationsDeviceRow{
			DeviceID:   binding.DeviceID,
			VehicleID:  binding.VehicleID,
			Status:     binding.Status,
			ValidFrom:  binding.ValidFrom,
			LastUsedAt: binding.LastUsedAt,
			RotatedAt:  binding.RotatedAt,
			RevokedAt:  binding.RevokedAt,
			Freshness:  "not seen",
			Assignment: "not available",
		}
		telemetryRow, ok := byDeviceVehicle[deviceVehicleKey(binding.DeviceID, binding.VehicleID)]
		if ok {
			applyOperationsDeviceTelemetry(&row, telemetryRow)
		}
		row.NextAction = operationsDeviceNextAction(row)
		rows = append(rows, row)
	}
	return rows
}

type operationsDeviceFleetSummary struct {
	TotalBindings            int
	ActiveBindings           int
	InactiveBindings         int
	FreshBindings            int
	StaleBindings            int
	NotSeenBindings          int
	UnknownAssignments       int
	LowConfidenceAssignments int
	UnboundTelemetryRows     int
}

func summarizeOperationsDeviceFleet(rows []operationsDeviceRow, telemetryRows []telemetryView) operationsDeviceFleetSummary {
	summary := operationsDeviceFleetSummary{TotalBindings: len(rows)}
	bound := make(map[string]bool, len(rows))
	for _, row := range rows {
		bound[deviceVehicleKey(row.DeviceID, row.VehicleID)] = true
		if row.Status == "active" || row.Status == "" {
			summary.ActiveBindings++
		} else {
			summary.InactiveBindings++
		}
		switch row.Freshness {
		case "fresh":
			summary.FreshBindings++
		case "stale":
			summary.StaleBindings++
		default:
			summary.NotSeenBindings++
		}
		if strings.Contains(strings.ToLower(row.Assignment), "unknown") {
			summary.UnknownAssignments++
		}
		if confidence, ok := parseOperationsDeviceConfidence(row.Assignment); ok && confidence < 0.70 {
			summary.LowConfidenceAssignments++
		}
	}
	for _, row := range telemetryRows {
		if !bound[deviceVehicleKey(row.DeviceID, row.VehicleID)] {
			summary.UnboundTelemetryRows++
		}
	}
	return summary
}

func deviceInventoryStatus(summary operationsDeviceFleetSummary) string {
	if summary.TotalBindings == 0 {
		return operationsStatusMissing
	}
	if summary.InactiveBindings > 0 {
		return operationsStatusNeedsReview
	}
	return operationsStatusReady
}

func deviceInventorySignal(summary operationsDeviceFleetSummary) string {
	return strconv.Itoa(summary.TotalBindings) + " bindings; active=" + strconv.Itoa(summary.ActiveBindings) + "; inactive_or_revoked=" + strconv.Itoa(summary.InactiveBindings)
}

func deviceCoverageStatus(summary operationsDeviceFleetSummary) string {
	if summary.TotalBindings == 0 {
		return operationsStatusMissing
	}
	if summary.NotSeenBindings > 0 || summary.StaleBindings > 0 || summary.UnboundTelemetryRows > 0 {
		return operationsStatusNeedsReview
	}
	return operationsStatusReady
}

func deviceCoverageSignal(summary operationsDeviceFleetSummary) string {
	return "fresh=" + strconv.Itoa(summary.FreshBindings) + "; stale=" + strconv.Itoa(summary.StaleBindings) + "; not_seen=" + strconv.Itoa(summary.NotSeenBindings) + "; unlisted_accepted_rows=" + strconv.Itoa(summary.UnboundTelemetryRows)
}

func notSeenStatus(summary operationsDeviceFleetSummary) string {
	if summary.TotalBindings == 0 {
		return operationsStatusMissing
	}
	if summary.NotSeenBindings > 0 {
		return operationsStatusNeedsReview
	}
	return operationsStatusReady
}

func staleStatus(summary operationsDeviceFleetSummary) string {
	if summary.StaleBindings > 0 {
		return operationsStatusNeedsReview
	}
	if summary.TotalBindings == 0 {
		return operationsStatusMissing
	}
	return operationsStatusReady
}

func bindingReviewStatus(summary operationsDeviceFleetSummary) string {
	if summary.TotalBindings == 0 {
		return operationsStatusMissing
	}
	if summary.UnknownAssignments > 0 || summary.LowConfidenceAssignments > 0 || summary.UnboundTelemetryRows > 0 {
		return operationsStatusNeedsReview
	}
	return operationsStatusReady
}

func bindingReviewSignal(summary operationsDeviceFleetSummary) string {
	return "unknown_assignments=" + strconv.Itoa(summary.UnknownAssignments) + "; low_confidence_assignments=" + strconv.Itoa(summary.LowConfidenceAssignments) + "; unlisted_accepted_rows=" + strconv.Itoa(summary.UnboundTelemetryRows)
}

func postInstallReviewStatus(summary operationsDeviceFleetSummary) string {
	if summary.TotalBindings == 0 || summary.NotSeenBindings == summary.TotalBindings {
		return operationsStatusMissing
	}
	if summary.StaleBindings > 0 || summary.UnknownAssignments > 0 || summary.LowConfidenceAssignments > 0 {
		return operationsStatusNeedsReview
	}
	return operationsStatusReady
}

func postInstallReviewSignal(summary operationsDeviceFleetSummary) string {
	return "accepted_or_seen=" + strconv.Itoa(summary.FreshBindings+summary.StaleBindings) + "; not_seen=" + strconv.Itoa(summary.NotSeenBindings) + "; unknown_or_low_confidence=" + strconv.Itoa(summary.UnknownAssignments+summary.LowConfidenceAssignments)
}

func applyOperationsDeviceTelemetry(row *operationsDeviceRow, telemetryRow telemetryView) {
	observed := telemetryRow.ObservedAt.UTC()
	received := telemetryRow.ReceivedAt.UTC()
	ageSeconds := telemetryRow.AgeSeconds
	row.LatestObservedAt = &observed
	row.LatestReceivedAt = &received
	row.LatestAgeSeconds = &ageSeconds
	if telemetryRow.Stale {
		row.Freshness = "stale"
	} else {
		row.Freshness = "fresh"
	}
	row.Assignment = operationsDeviceAssignmentSummary(telemetryRow)
	row.AssignmentAt = telemetryRow.AssignmentAt
	row.AssignmentSource = telemetryRow.AssignmentSource
}

func operationsDeviceAssignmentSummary(row telemetryView) string {
	state := strings.TrimSpace(row.AssignmentState)
	if state == "" {
		return "not available"
	}
	parts := []string{state}
	if strings.TrimSpace(row.RouteID) != "" {
		parts = append(parts, "route "+row.RouteID)
	}
	if strings.TrimSpace(row.TripID) != "" {
		parts = append(parts, "trip "+row.TripID)
	}
	if strings.TrimSpace(row.Confidence) != "" {
		parts = append(parts, "confidence "+row.Confidence)
	}
	return strings.Join(parts, " / ")
}

func operationsDeviceNextAction(row operationsDeviceRow) string {
	if row.Status != "" && row.Status != "active" {
		return "Review credential status before using this device for service."
	}
	if row.RevokedAt != nil {
		return "Rotate or replace the revoked credential before sending telemetry."
	}
	switch row.Freshness {
	case "not seen":
		return "Install the one-time token on the device and send an authenticated sample telemetry event."
	case "stale":
		return "Check device power, network, and reporting cadence, then confirm fresh telemetry."
	}
	if row.Assignment == "not available" {
		return "Confirm vehicle and trip hints or review matching before relying on public trip descriptors."
	}
	if strings.Contains(row.Assignment, "unknown") {
		return "Leave the trip descriptor unknown until matching confidence improves or an operator override is applied."
	}
	if confidence, ok := parseOperationsDeviceConfidence(row.Assignment); ok && confidence < 0.70 {
		return "Review low-confidence matching evidence before acting on the assignment."
	}
	return "No immediate action. Continue monitoring freshness and assignment confidence."
}

func parseOperationsDeviceConfidence(summary string) (float64, bool) {
	const marker = "confidence "
	idx := strings.LastIndex(summary, marker)
	if idx < 0 {
		return 0, false
	}
	value := strings.TrimSpace(summary[idx+len(marker):])
	confidence, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return confidence, true
}

func deviceVehicleKey(deviceID, vehicleID string) string {
	return deviceID + "\x00" + vehicleID
}
