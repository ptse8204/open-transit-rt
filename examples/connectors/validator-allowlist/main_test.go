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
	if manifest.ConnectorType != connectors.TypeValidator {
		t.Fatalf("type = %q", manifest.ConnectorType)
	}
}

func TestBuildDecisionAllowsOnlyServerOwnedIDsAndNeverRunsRawCommands(t *testing.T) {
	request, err := readRequest("fixtures/request.json")
	if err != nil {
		t.Fatal(err)
	}
	decision := BuildDecision(request)
	if !decision.Allowed || decision.Result != "allowlisted" {
		t.Fatalf("decision = %+v, want allowlisted", decision)
	}
	if decision.RawCommandAllowed || decision.NetworkSend || decision.StatusMutation || decision.EvidenceWrite {
		t.Fatalf("decision enables unsafe behavior: %+v", decision)
	}

	blocked := BuildDecision(Request{SyntheticOnly: true, ValidatorID: "unlisted-validator", FeedType: "schedule", ArtifactRef: "fixture://schedule.zip"})
	if blocked.Allowed || blocked.Result != "blocked" {
		t.Fatalf("blocked decision = %+v, want blocked", blocked)
	}
}

func TestBuildDecisionFailsClosedForWrongFeedTypeAndUnsafeArtifactRef(t *testing.T) {
	wrongFeed := BuildDecision(Request{SyntheticOnly: true, ValidatorID: "static-mobilitydata", FeedType: "alerts", ArtifactRef: "fixture://alerts.pb"})
	if wrongFeed.Allowed || wrongFeed.Result != "blocked" {
		t.Fatalf("wrong feed decision = %+v, want blocked", wrongFeed)
	}

	unsafeArtifact := BuildDecision(Request{SyntheticOnly: true, ValidatorID: "static-mobilitydata", FeedType: "schedule", ArtifactRef: "/Users/example/private/schedule.zip"})
	if unsafeArtifact.Allowed || unsafeArtifact.ArtifactRef != "redacted" {
		t.Fatalf("unsafe artifact decision = %+v, want blocked with redacted artifact", unsafeArtifact)
	}
}

func TestReadRequestRejectsUnsafeArtifactRef(t *testing.T) {
	tmp := t.TempDir() + "/request.json"
	if err := os.WriteFile(tmp, []byte(`{"synthetic_only":true,"validator_id":"static-mobilitydata","feed_type":"schedule","artifact_ref":"file:///tmp/private.zip"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readRequest(tmp); err == nil {
		t.Fatal("readRequest succeeded with unsafe artifact_ref")
	}
}
