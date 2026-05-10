package main

import (
	"os"
	"strings"
	"testing"
	"time"

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

func TestTransformDryRunDropsLowQuality(t *testing.T) {
	observations, err := readObservations("fixtures/observations.json")
	if err != nil {
		t.Fatal(err)
	}
	summary := Transform(observations, time.Date(2026, 5, 10, 15, 1, 0, 0, time.UTC))
	if len(summary.Events) != 2 {
		t.Fatalf("events = %d, want 2", len(summary.Events))
	}
	if len(summary.Drops) != 1 || !strings.Contains(summary.Drops[0].Reason, "low quality") {
		t.Fatalf("drops = %+v, want low-quality drop", summary.Drops)
	}
	for _, event := range summary.Events {
		if !event.DryRun || event.NetworkSend {
			t.Fatalf("event send flags = dry_run %v network_send %v", event.DryRun, event.NetworkSend)
		}
	}
}
