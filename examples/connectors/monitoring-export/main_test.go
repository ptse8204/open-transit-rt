package main

import (
	"os"
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
	if manifest.ConnectorType != connectors.TypeMonitoringExport {
		t.Fatalf("type = %q", manifest.ConnectorType)
	}
}

func TestBuildExportBatchRedactsAndDoesNotSend(t *testing.T) {
	input, err := readMetrics("fixtures/metrics.json")
	if err != nil {
		t.Fatal(err)
	}
	batch := BuildExportBatch(input)
	if batch.SendEnabled || batch.NetworkSend {
		t.Fatalf("send flags = send_enabled %v network_send %v", batch.SendEnabled, batch.NetworkSend)
	}
	if batch.StatusMutation || batch.EvidenceWrite {
		t.Fatalf("batch mutates status or writes evidence: %+v", batch)
	}
	if len(batch.Incidents) != 1 {
		t.Fatalf("incidents = %d, want 1", len(batch.Incidents))
	}
	if batch.Incidents[0].ID == "" || batch.Incidents[0].Summary == "" {
		t.Fatalf("redacted incident lost public fields: %+v", batch.Incidents[0])
	}
}
