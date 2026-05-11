package main

import (
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

func TestParseReplay(t *testing.T) {
	file, err := os.Open("fixtures/replay.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	summary, err := ParseReplay(file, time.Date(2026, 5, 10, 16, 1, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if len(summary.Events) != 2 {
		t.Fatalf("events = %d, want 2", len(summary.Events))
	}
	if len(summary.Drops) != 1 || summary.Drops[0].Reason != telemetrysdk.ReasonLowQuality {
		t.Fatalf("drops = %+v, want low-quality drop", summary.Drops)
	}
	for _, event := range summary.Events {
		if !event.DryRun || event.NetworkSend {
			t.Fatalf("event send flags = dry_run %v network_send %v", event.DryRun, event.NetworkSend)
		}
	}
}

func TestParseReplayFailsClosedOnBadTimestamp(t *testing.T) {
	summary, err := ParseReplay(strings.NewReader("synthetic_only,agency_id,device_id,vehicle_id,observed_at,latitude,longitude,quality\ntrue,agency-demo,device-1,vehicle-1,not-time,1,2,1\n"), time.Date(2026, 5, 10, 16, 1, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if len(summary.Events) != 0 || len(summary.Drops) != 1 || summary.Drops[0].Reason != telemetrysdk.ReasonInvalidTimestamp {
		t.Fatalf("summary = %+v, want invalid timestamp drop", summary)
	}
}
