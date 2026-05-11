package predictionsdk

import "testing"

func TestBuildDryRunResponseWithholdsAndNeverMutatesPublicTripUpdates(t *testing.T) {
	request := Request{
		ActiveFeedVersion:   "feed-version-synthetic",
		VehiclePositionsURL: "https://example.org/gtfsrt/vehicle_positions.pb",
		Telemetry:           []Telemetry{{VehicleID: "vehicle-good"}},
		Assignments: []Assignment{
			{VehicleID: "vehicle-good", TripID: "trip-good", StartDate: "20260510", StartTime: "10:00:00", Confidence: 0.92},
			{VehicleID: "vehicle-low", TripID: "trip-low", StartDate: "20260510", StartTime: "10:05:00", Confidence: 0.2},
			{VehicleID: "vehicle-missing", TripID: "trip-missing", StartDate: "20260510", StartTime: "10:10:00", Confidence: 0.95},
			{VehicleID: "vehicle-no-trip", Confidence: 0.95},
		},
	}

	response := BuildDryRunResponse(request, Options{ConfidenceThreshold: 0.8})
	if response.NetworkSend || response.PublicTripUpdatesMutation || len(response.TripUpdates) != 0 {
		t.Fatalf("response mutates or emits output: %+v", response)
	}
	if response.VehiclePositionsReference != request.VehiclePositionsURL || response.ActiveFeedVersion != request.ActiveFeedVersion {
		t.Fatalf("response references = %+v, want request feed and Vehicle Positions reference", response)
	}
	if response.Counts.EligibleNoop != 1 || response.Counts.WithheldUnknown != 1 || response.Counts.MissingTelemetry != 1 || response.Counts.MissingDescriptor != 1 {
		t.Fatalf("counts = %+v, want one of each state", response.Counts)
	}
	wantReasons := map[string]string{
		"vehicle-good":    ReasonEligibleNoop,
		"vehicle-low":     ReasonBelowConfidence,
		"vehicle-missing": ReasonMissingTelemetry,
		"vehicle-no-trip": ReasonMissingTripDescriptor,
	}
	for _, diagnostic := range response.Diagnostics {
		if wantReasons[diagnostic.VehicleID] != diagnostic.Reason {
			t.Fatalf("diagnostic = %+v, want reason %q", diagnostic, wantReasons[diagnostic.VehicleID])
		}
	}
}
