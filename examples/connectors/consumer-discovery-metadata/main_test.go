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
	if manifest.ConnectorType != connectors.TypeConsumerDiscovery {
		t.Fatalf("type = %q", manifest.ConnectorType)
	}
}

func TestBuildDecisionKeepsPreparedOnlyAndNoSubmit(t *testing.T) {
	metadata, err := readMetadata("fixtures/feeds.json")
	if err != nil {
		t.Fatal(err)
	}
	decision := BuildDecision(metadata)
	if !decision.ReadyForLocalReview || !decision.PreparedOnly {
		t.Fatalf("decision = %+v, want ready prepared-only review", decision)
	}
	if decision.SubmitEnabled || decision.StatusMutation || decision.NetworkSend || decision.EvidenceWrite {
		t.Fatalf("decision enables unsafe behavior: %+v", decision)
	}
}

func TestBuildDecisionBlocksUnsafeFlagsAndMissingURLs(t *testing.T) {
	metadata := FeedMetadata{
		SyntheticOnly:        true,
		AgencyID:             "demo-agency",
		FeedBaseURL:          "https://example.org/open-transit-rt",
		StaticGTFSURL:        "https://example.org/open-transit-rt/public/gtfs/schedule.zip",
		VehiclePositionsURL:  "https://example.org/open-transit-rt/public/gtfsrt/vehicle_positions.pb",
		TripUpdatesURL:       "https://example.org/open-transit-rt/public/gtfsrt/trip_updates.pb",
		AlertsURL:            "https://example.org/open-transit-rt/public/gtfsrt/alerts.pb",
		LicenseURL:           "https://example.org/open-transit-rt/license",
		TechnicalContactRole: "deployment owner",
		PreparedOnly:         true,
		SubmitEnabled:        true,
	}
	decision := BuildDecision(metadata)
	if decision.ReadyForLocalReview {
		t.Fatalf("decision = %+v, want blocked when submit flag is present", decision)
	}
	if decision.SubmitEnabled || decision.StatusMutation || decision.NetworkSend || decision.EvidenceWrite {
		t.Fatalf("decision should keep unsafe output flags false: %+v", decision)
	}

	metadata.SubmitEnabled = false
	metadata.TripUpdatesURL = ""
	decision = BuildDecision(metadata)
	if decision.ReadyForLocalReview || len(decision.MissingFields) == 0 {
		t.Fatalf("decision = %+v, want missing field blocker", decision)
	}
}

func TestReadMetadataRejectsNonSyntheticFixture(t *testing.T) {
	tmp := t.TempDir() + "/feeds.json"
	if err := os.WriteFile(tmp, []byte(`{"synthetic_only":false}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readMetadata(tmp); err == nil {
		t.Fatal("readMetadata succeeded with non-synthetic fixture")
	}
}
