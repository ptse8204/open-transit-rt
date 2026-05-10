package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const defaultPredictionFixture = "examples/connectors/predictor-sidecar-stub/fixtures/prediction-input.json"

type Request struct {
	SyntheticOnly       bool         `json:"synthetic_only"`
	ActiveFeedVersion   string       `json:"active_feed_version"`
	VehiclePositionsURL string       `json:"vehicle_positions_url"`
	Telemetry           []Telemetry  `json:"telemetry"`
	Assignments         []Assignment `json:"assignments"`
}

type Telemetry struct {
	AgencyID   string  `json:"agency_id"`
	VehicleID  string  `json:"vehicle_id"`
	ObservedAt string  `json:"observed_at"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type Assignment struct {
	VehicleID  string  `json:"vehicle_id"`
	TripID     string  `json:"trip_id"`
	StartDate  string  `json:"start_date"`
	StartTime  string  `json:"start_time"`
	Confidence float64 `json:"confidence"`
}

type Response struct {
	Mode             string       `json:"mode"`
	NetworkSend      bool         `json:"network_send"`
	TripUpdates      []TripUpdate `json:"trip_updates"`
	Diagnostics      []Diagnostic `json:"diagnostics"`
	VehiclePositions string       `json:"vehicle_positions_reference"`
}

type TripUpdate struct {
	TripID    string `json:"trip_id"`
	StartDate string `json:"start_date"`
	StartTime string `json:"start_time"`
}

type Diagnostic struct {
	VehicleID string `json:"vehicle_id"`
	State     string `json:"state"`
	Reason    string `json:"reason"`
}

func main() {
	path := defaultPredictionFixture
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	request, err := readRequest(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	response := BuildResponse(request, 0.8)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(response); err != nil {
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
		return Request{}, errors.New("prediction request fixture must be marked synthetic_only")
	}
	if request.ActiveFeedVersion == "" || request.VehiclePositionsURL == "" {
		return Request{}, errors.New("prediction request missing feed version or vehicle positions reference")
	}
	return request, nil
}

func BuildResponse(request Request, threshold float64) Response {
	response := Response{
		Mode:             "dry_run_noop",
		NetworkSend:      false,
		VehiclePositions: request.VehiclePositionsURL,
	}
	for _, assignment := range request.Assignments {
		if assignment.Confidence < threshold {
			response.Diagnostics = append(response.Diagnostics, Diagnostic{
				VehicleID: assignment.VehicleID,
				State:     "unknown",
				Reason:    "below confidence threshold",
			})
			continue
		}
		response.Diagnostics = append(response.Diagnostics, Diagnostic{
			VehicleID: assignment.VehicleID,
			State:     "eligible",
			Reason:    "stub does not emit public Trip Updates",
		})
	}
	return response
}
