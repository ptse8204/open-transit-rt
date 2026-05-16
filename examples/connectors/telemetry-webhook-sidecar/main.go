package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	telemetrysdk "open-transit-rt/examples/connectors/sdk/telemetry"
)

const defaultWebhookFixture = "examples/connectors/telemetry-webhook-sidecar/fixtures/webhook.json"

type WebhookEnvelope struct {
	SyntheticOnly bool                 `json:"synthetic_only"`
	Source        string               `json:"source"`
	ReceivedAt    string               `json:"received_at"`
	Observations  []WebhookObservation `json:"observations"`
}

type WebhookObservation struct {
	AgencyID  string  `json:"agency_id"`
	DeviceID  string  `json:"device_id"`
	VehicleID string  `json:"vehicle_id"`
	Timestamp string  `json:"timestamp"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Quality   float64 `json:"quality"`
}

func main() {
	path := defaultWebhookFixture
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	envelope, err := readWebhook(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	summary := Transform(envelope, time.Date(2026, 5, 10, 17, 1, 0, 0, time.UTC))
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readWebhook(path string) (WebhookEnvelope, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return WebhookEnvelope{}, err
	}
	var envelope WebhookEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return WebhookEnvelope{}, err
	}
	if !envelope.SyntheticOnly {
		return WebhookEnvelope{}, errors.New("webhook fixture must be marked synthetic_only")
	}
	if len(envelope.Observations) == 0 {
		return WebhookEnvelope{}, errors.New("webhook fixture is empty")
	}
	return envelope, nil
}

func Transform(envelope WebhookEnvelope, now time.Time) telemetrysdk.Summary {
	source := envelope.Source
	if source == "" {
		source = "synthetic-webhook-sidecar"
	}
	var parsed []telemetrysdk.Observation
	var summary telemetrysdk.Summary
	for _, observation := range envelope.Observations {
		observedAt, err := time.Parse(time.RFC3339, observation.Timestamp)
		if err != nil {
			summary.Drops = append(summary.Drops, telemetrysdk.Drop{DeviceID: observation.DeviceID, VehicleID: observation.VehicleID, Reason: telemetrysdk.ReasonInvalidTimestamp})
			continue
		}
		parsed = append(parsed, telemetrysdk.Observation{
			AgencyID:   observation.AgencyID,
			DeviceID:   observation.DeviceID,
			VehicleID:  observation.VehicleID,
			ObservedAt: observedAt,
			Latitude:   observation.Latitude,
			Longitude:  observation.Longitude,
			Quality:    observation.Quality,
		})
	}
	normalized := telemetrysdk.NormalizeBatch(parsed, telemetrysdk.Options{Now: now, Source: source})
	summary.Events = append(summary.Events, normalized.Events...)
	summary.Drops = append(summary.Drops, normalized.Drops...)
	return summary
}
