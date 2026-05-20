package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	telemetrysdk "open-transit-rt/examples/connectors/sdk/telemetry"
	"open-transit-rt/internal/connectors"
)

func TestConnectorManifestValidates(t *testing.T) {
	raw, err := os.Open("connector.json")
	if err != nil {
		t.Fatal(err)
	}
	defer raw.Close()

	manifest, err := connectors.DecodeManifest(raw)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.ConnectorType != connectors.TypeTelemetrySource {
		t.Fatalf("type = %q", manifest.ConnectorType)
	}
}

func TestTransformDryRunNormalizesMappedRecords(t *testing.T) {
	fixture, err := readFixture("fixtures/observations.json")
	if err != nil {
		t.Fatal(err)
	}
	result := Transform(fixture, time.Date(2026, 5, 10, 15, 1, 0, 0, time.UTC))
	if !result.DryRun || result.NetworkSend {
		t.Fatalf("send flags = dry_run %v network_send %v", result.DryRun, result.NetworkSend)
	}
	if result.TimeoutSeconds != 5 {
		t.Fatalf("timeout = %d, want fixture timeout", result.TimeoutSeconds)
	}
	if len(result.Events) != 2 {
		t.Fatalf("events = %d, want 2", len(result.Events))
	}
	if len(result.Drops) != 1 || result.Drops[0].Reason != telemetrysdk.ReasonLowQuality {
		t.Fatalf("drops = %+v, want low-quality drop", result.Drops)
	}
	for _, event := range result.Events {
		if event.Source != "synthetic-generic-json-transform" || !event.DryRun || event.NetworkSend {
			t.Fatalf("event = %+v, want dry-run generic source", event)
		}
	}
	rendered, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(rendered), "operator_note") || strings.Contains(string(rendered), "raw_payload") {
		t.Fatalf("result contains unmapped/raw fields: %s", rendered)
	}
}

func TestReadFixtureRejectsUnsafeSendAndBadMapping(t *testing.T) {
	fixture, err := readFixture("fixtures/observations.json")
	if err != nil {
		t.Fatal(err)
	}
	fixture.SendEnabled = true
	if err := validateFixture(fixture); err == nil || !strings.Contains(err.Error(), "send_enabled") {
		t.Fatalf("send-enabled validation error = %v", err)
	}
	fixture.SendEnabled = false
	delete(fixture.FieldMap, "vehicle_id")
	if err := validateFixture(fixture); err == nil || !strings.Contains(err.Error(), "field_map.vehicle_id") {
		t.Fatalf("missing mapping validation error = %v", err)
	}
}

func TestTransformDropsMalformedMappedRecord(t *testing.T) {
	fixture := Fixture{
		SyntheticOnly:  true,
		SendEnabled:    false,
		TimeoutSeconds: 5,
		FieldMap: map[string]string{
			"agency_id":  "agency",
			"device_id":  "device",
			"vehicle_id": "vehicle",
			"timestamp":  "observed_at",
			"latitude":   "lat",
			"longitude":  "lon",
			"quality":    "quality",
		},
		Records: []map[string]any{{
			"agency":      "agency-demo",
			"device":      "device-bad",
			"vehicle":     "vehicle-bad",
			"observed_at": "not-a-time",
			"lat":         34.0,
			"lon":         -118.0,
			"quality":     0.95,
		}},
	}
	result := Transform(fixture, time.Date(2026, 5, 10, 15, 1, 0, 0, time.UTC))
	if len(result.Events) != 0 || len(result.Drops) != 1 || result.Drops[0].Reason != telemetrysdk.ReasonInvalidTimestamp {
		t.Fatalf("result = %+v, want invalid timestamp drop only", result)
	}
}
