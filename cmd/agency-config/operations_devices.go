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
