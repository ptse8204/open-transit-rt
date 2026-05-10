package main

import (
	"os"
	"strings"
	"testing"

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

	events, err := ParseReplay(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("events = %d, want 2", len(events))
	}
	for _, event := range events {
		if !event.DryRun || event.NetworkSend {
			t.Fatalf("event send flags = dry_run %v network_send %v", event.DryRun, event.NetworkSend)
		}
	}
}

func TestParseReplayFailsClosedOnBadTimestamp(t *testing.T) {
	_, err := ParseReplay(strings.NewReader("synthetic_only,agency_id,device_id,vehicle_id,observed_at,latitude,longitude\ntrue,agency-demo,device-1,vehicle-1,not-time,1,2\n"))
	if err == nil {
		t.Fatal("ParseReplay succeeded, want failure")
	}
}
