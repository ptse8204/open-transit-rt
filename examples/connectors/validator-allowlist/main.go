package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

const defaultValidatorFixture = "examples/connectors/validator-allowlist/fixtures/request.json"

var allowlistedValidators = map[string]map[string]bool{
	"static-mobilitydata": {
		"schedule": true,
	},
	"realtime-mobilitydata": {
		"vehicle_positions": true,
		"trip_updates":      true,
		"alerts":            true,
	},
}

type Request struct {
	SyntheticOnly bool   `json:"synthetic_only"`
	ValidatorID   string `json:"validator_id"`
	FeedType      string `json:"feed_type"`
	ArtifactRef   string `json:"artifact_ref"`
}

type Decision struct {
	ValidatorID        string `json:"validator_id"`
	FeedType           string `json:"feed_type"`
	ArtifactRef        string `json:"artifact_ref"`
	Allowed            bool   `json:"allowed"`
	Result             string `json:"result"`
	RawCommandAllowed  bool   `json:"raw_command_allowed"`
	NetworkSend        bool   `json:"network_send"`
	StatusMutation     bool   `json:"status_mutation"`
	EvidenceWrite      bool   `json:"evidence_write"`
	FailureBehavior    string `json:"failure_behavior"`
	ClaimBoundary      string `json:"claim_boundary"`
	OperatorNextAction string `json:"operator_next_action"`
}

func main() {
	path := defaultValidatorFixture
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	request, err := readRequest(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	decision := BuildDecision(request)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(decision); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readRequest(path string) (Request, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Request{}, err
	}
	var request Request
	if err := json.Unmarshal(raw, &request); err != nil {
		return Request{}, err
	}
	if !request.SyntheticOnly {
		return Request{}, errors.New("validator request fixture must be marked synthetic_only")
	}
	if request.ValidatorID == "" || request.FeedType == "" || request.ArtifactRef == "" {
		return Request{}, errors.New("validator request missing validator_id, feed_type, or artifact_ref")
	}
	if !isSafeFixtureRef(request.ArtifactRef) {
		return Request{}, errors.New("validator artifact_ref must use a safe fixture:// reference")
	}
	return request, nil
}

func BuildDecision(request Request) Decision {
	safeArtifactRef := "redacted"
	artifactRefOK := isSafeFixtureRef(request.ArtifactRef)
	if artifactRefOK {
		safeArtifactRef = request.ArtifactRef
	}
	allowed := validatorSupportsFeedType(request.ValidatorID, request.FeedType) && artifactRefOK
	result := "blocked"
	nextAction := "Use a server-owned allowlisted validator ID."
	if allowed {
		result = "allowlisted"
		nextAction = "Run validation through Open Transit RT server-owned validator mappings."
	}
	return Decision{
		ValidatorID:        request.ValidatorID,
		FeedType:           request.FeedType,
		ArtifactRef:        safeArtifactRef,
		Allowed:            allowed,
		Result:             result,
		RawCommandAllowed:  false,
		NetworkSend:        false,
		StatusMutation:     false,
		EvidenceWrite:      false,
		FailureBehavior:    "fail closed before validator execution",
		ClaimBoundary:      "allowlist decision only; not compliance, consumer acceptance, production readiness, or evidence",
		OperatorNextAction: nextAction,
	}
}

func validatorSupportsFeedType(validatorID, feedType string) bool {
	supported, ok := allowlistedValidators[validatorID]
	return ok && supported[feedType]
}

func isSafeFixtureRef(ref string) bool {
	const prefix = "fixture://"
	if !strings.HasPrefix(ref, prefix) {
		return false
	}
	name := strings.TrimPrefix(ref, prefix)
	if name == "" || strings.Contains(name, "..") {
		return false
	}
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
			continue
		}
		return false
	}
	return true
}
