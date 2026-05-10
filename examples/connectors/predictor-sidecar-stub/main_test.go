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
	if manifest.ConnectorType != connectors.TypePrediction {
		t.Fatalf("type = %q", manifest.ConnectorType)
	}
}

func TestBuildResponseIsNoopAndKeepsLowConfidenceUnknown(t *testing.T) {
	request, err := readRequest("fixtures/prediction-input.json")
	if err != nil {
		t.Fatal(err)
	}
	response := BuildResponse(request, 0.8)
	if response.NetworkSend {
		t.Fatal("response enables network send")
	}
	if len(response.TripUpdates) != 0 {
		t.Fatalf("trip updates = %d, want 0 for stub", len(response.TripUpdates))
	}
	foundUnknown := false
	for _, diagnostic := range response.Diagnostics {
		if diagnostic.VehicleID == "vehicle-low" && diagnostic.State == "unknown" {
			foundUnknown = true
		}
	}
	if !foundUnknown {
		t.Fatalf("diagnostics = %+v, want low-confidence unknown", response.Diagnostics)
	}
}
