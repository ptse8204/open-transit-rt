package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
)

const defaultDiscoveryFixture = "examples/connectors/consumer-discovery-metadata/fixtures/feeds.json"

type FeedMetadata struct {
	SyntheticOnly        bool   `json:"synthetic_only"`
	AgencyID             string `json:"agency_id"`
	FeedBaseURL          string `json:"feed_base_url"`
	StaticGTFSURL        string `json:"static_gtfs_url"`
	VehiclePositionsURL  string `json:"vehicle_positions_url"`
	TripUpdatesURL       string `json:"trip_updates_url"`
	AlertsURL            string `json:"alerts_url"`
	LicenseURL           string `json:"license_url"`
	TechnicalContactRole string `json:"technical_contact_role"`
	PreparedOnly         bool   `json:"prepared_only"`
	SubmitEnabled        bool   `json:"submit_enabled"`
	StatusMutation       bool   `json:"status_mutation"`
	NetworkSend          bool   `json:"network_send"`
	EvidenceWrite        bool   `json:"evidence_write"`
}

type DiscoveryDecision struct {
	AgencyID            string   `json:"agency_id"`
	ReadyForLocalReview bool     `json:"ready_for_local_review"`
	MissingFields       []string `json:"missing_fields"`
	PreparedOnly        bool     `json:"prepared_only"`
	SubmitEnabled       bool     `json:"submit_enabled"`
	StatusMutation      bool     `json:"status_mutation"`
	NetworkSend         bool     `json:"network_send"`
	EvidenceWrite       bool     `json:"evidence_write"`
	FailureBehavior     string   `json:"failure_behavior"`
	ClaimBoundary       string   `json:"claim_boundary"`
	OperatorNextAction  string   `json:"operator_next_action"`
}

func main() {
	path := defaultDiscoveryFixture
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	metadata, err := readMetadata(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	decision := BuildDecision(metadata)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(decision); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readMetadata(path string) (FeedMetadata, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return FeedMetadata{}, err
	}
	var metadata FeedMetadata
	if err := json.Unmarshal(raw, &metadata); err != nil {
		return FeedMetadata{}, err
	}
	if !metadata.SyntheticOnly {
		return FeedMetadata{}, errors.New("consumer discovery fixture must be marked synthetic_only")
	}
	return metadata, nil
}

func BuildDecision(metadata FeedMetadata) DiscoveryDecision {
	missing := missingMetadataFields(metadata)
	unsafeFlags := metadata.SubmitEnabled || metadata.StatusMutation || metadata.NetworkSend || metadata.EvidenceWrite || !metadata.PreparedOnly
	ready := len(missing) == 0 && !unsafeFlags
	nextAction := "Fix missing metadata and keep submission/status movement disabled."
	if ready {
		nextAction = "Review prepared-only metadata in the private console before any separately authorized external workflow."
	}
	return DiscoveryDecision{
		AgencyID:            metadata.AgencyID,
		ReadyForLocalReview: ready,
		MissingFields:       missing,
		PreparedOnly:        metadata.PreparedOnly,
		SubmitEnabled:       false,
		StatusMutation:      false,
		NetworkSend:         false,
		EvidenceWrite:       false,
		FailureBehavior:     "fail closed before submission, status mutation, network send, or evidence write",
		ClaimBoundary:       "prepared metadata only; not submission, acceptance, compliance, public launch, or production readiness",
		OperatorNextAction:  nextAction,
	}
}

func missingMetadataFields(metadata FeedMetadata) []string {
	checks := []struct {
		name  string
		value string
	}{
		{"agency_id", metadata.AgencyID},
		{"feed_base_url", metadata.FeedBaseURL},
		{"static_gtfs_url", metadata.StaticGTFSURL},
		{"vehicle_positions_url", metadata.VehiclePositionsURL},
		{"trip_updates_url", metadata.TripUpdatesURL},
		{"alerts_url", metadata.AlertsURL},
		{"license_url", metadata.LicenseURL},
		{"technical_contact_role", metadata.TechnicalContactRole},
	}
	var missing []string
	for _, check := range checks {
		if check.value == "" || !safePublicURLOrRole(check.name, check.value) {
			missing = append(missing, check.name)
		}
	}
	return missing
}

func safePublicURLOrRole(name string, value string) bool {
	if name == "agency_id" || name == "technical_contact_role" {
		return value != ""
	}
	parsed, err := url.Parse(value)
	return err == nil && parsed.Scheme == "https" && parsed.Hostname() != ""
}
