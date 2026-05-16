package main

import (
	"os"
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
	if manifest.SchemaVersion != connectors.SchemaVersion {
		t.Fatalf("schema = %q", manifest.SchemaVersion)
	}
}

func TestTransformWebhookDryRunDropsLowQuality(t *testing.T) {
	envelope, err := readWebhook("fixtures/webhook.json")
	if err != nil {
		t.Fatal(err)
	}
	summary := Transform(envelope, time.Date(2026, 5, 10, 17, 1, 0, 0, time.UTC))
	if len(summary.Events) != 1 {
		t.Fatalf("events = %d, want 1", len(summary.Events))
	}
	if len(summary.Drops) != 1 || summary.Drops[0].Reason != telemetrysdk.ReasonLowQuality {
		t.Fatalf("drops = %+v, want low-quality drop", summary.Drops)
	}
	if !summary.Events[0].DryRun || summary.Events[0].NetworkSend {
		t.Fatalf("event send flags = dry_run %v network_send %v", summary.Events[0].DryRun, summary.Events[0].NetworkSend)
	}
}

func TestTransformWebhookDropsInvalidTimestamp(t *testing.T) {
	summary := Transform(WebhookEnvelope{
		SyntheticOnly: true,
		Source:        "synthetic-webhook-sidecar",
		Observations: []WebhookObservation{{
			AgencyID:  "agency-demo",
			DeviceID:  "device-bad-time",
			VehicleID: "vehicle-bad-time",
			Timestamp: "not-a-time",
			Latitude:  34,
			Longitude: -118,
			Quality:   1,
		}},
	}, time.Date(2026, 5, 10, 17, 1, 0, 0, time.UTC))
	if len(summary.Events) != 0 || len(summary.Drops) != 1 || summary.Drops[0].Reason != telemetrysdk.ReasonInvalidTimestamp {
		t.Fatalf("summary = %+v, want invalid timestamp drop", summary)
	}
}
