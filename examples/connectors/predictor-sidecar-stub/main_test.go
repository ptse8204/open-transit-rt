package main

import (
	"os"
	"testing"

	predictionsdk "open-transit-rt/examples/connectors/sdk/prediction"
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
	if response.NetworkSend || response.PublicTripUpdatesMutation {
		t.Fatalf("response enables send or public mutation: %+v", response)
	}
	if len(response.TripUpdates) != 0 {
		t.Fatalf("trip updates = %d, want 0 for stub", len(response.TripUpdates))
	}
	if response.VehiclePositionsReference != request.VehiclePositionsURL || response.ActiveFeedVersion != request.ActiveFeedVersion {
		t.Fatalf("response references = %+v, want request references", response)
	}
	foundUnknown := false
	for _, diagnostic := range response.Diagnostics {
		if diagnostic.VehicleID == "vehicle-low" && diagnostic.State == "unknown" && diagnostic.Reason == predictionsdk.ReasonBelowConfidence {
			foundUnknown = true
		}
	}
	if !foundUnknown {
		t.Fatalf("diagnostics = %+v, want low-confidence unknown", response.Diagnostics)
	}
}

func TestBuildResponseWithholdsMissingTelemetryAndDescriptor(t *testing.T) {
	request, err := readRequest("fixtures/prediction-input.json")
	if err != nil {
		t.Fatal(err)
	}
	response := BuildResponse(request, 0.8)
	var missingTelemetry bool
	var missingDescriptor bool
	for _, diagnostic := range response.Diagnostics {
		if diagnostic.VehicleID == "vehicle-no-telemetry" && diagnostic.Reason == predictionsdk.ReasonMissingTelemetry {
			missingTelemetry = true
		}
		if diagnostic.VehicleID == "vehicle-no-trip" && diagnostic.Reason == predictionsdk.ReasonMissingTripDescriptor {
			missingDescriptor = true
		}
	}
	if !missingTelemetry || !missingDescriptor {
		t.Fatalf("diagnostics = %+v, want missing telemetry and descriptor withholding", response.Diagnostics)
	}
}
