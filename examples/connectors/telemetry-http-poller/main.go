package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const defaultFixture = "examples/connectors/telemetry-http-poller/fixtures/observations.json"

type Observation struct {
	AgencyID  string  `json:"agency_id"`
	DeviceID  string  `json:"device_id"`
	VehicleID string  `json:"vehicle_id"`
	Timestamp string  `json:"timestamp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Quality   float64 `json:"quality"`
}

type ObservationFixture struct {
	SyntheticOnly bool          `json:"synthetic_only"`
	Observations  []Observation `json:"observations"`
}

type TelemetryEvent struct {
	AgencyID    string  `json:"agency_id"`
	DeviceID    string  `json:"device_id"`
	VehicleID   string  `json:"vehicle_id"`
	ObservedAt  string  `json:"observed_at"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Source      string  `json:"source"`
	DryRun      bool    `json:"dry_run"`
	NetworkSend bool    `json:"network_send"`
}

type Drop struct {
	DeviceID string `json:"device_id"`
	Reason   string `json:"reason"`
}

type Summary struct {
	Events []TelemetryEvent `json:"events"`
	Drops  []Drop           `json:"drops"`
}

func main() {
	path := defaultFixture
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	observations, err := readObservations(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	summary := Transform(observations, time.Date(2026, 5, 10, 15, 1, 0, 0, time.UTC))
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readObservations(path string) ([]Observation, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var fixture ObservationFixture
	if err := json.Unmarshal(raw, &fixture); err != nil {
		return nil, err
	}
	if !fixture.SyntheticOnly {
		return nil, errors.New("observation fixture must be marked synthetic_only")
	}
	observations := fixture.Observations
	if len(observations) == 0 {
		return nil, errors.New("observation fixture is empty")
	}
	return observations, nil
}

func Transform(observations []Observation, now time.Time) Summary {
	var summary Summary
	for _, observation := range observations {
		observedAt, err := time.Parse(time.RFC3339, observation.Timestamp)
		switch {
		case observation.AgencyID == "" || observation.DeviceID == "" || observation.VehicleID == "":
			summary.Drops = append(summary.Drops, Drop{DeviceID: observation.DeviceID, Reason: "missing identity"})
		case err != nil:
			summary.Drops = append(summary.Drops, Drop{DeviceID: observation.DeviceID, Reason: "invalid timestamp"})
		case observedAt.After(now.Add(30 * time.Second)):
			summary.Drops = append(summary.Drops, Drop{DeviceID: observation.DeviceID, Reason: "future timestamp"})
		case now.Sub(observedAt) > 2*time.Minute:
			summary.Drops = append(summary.Drops, Drop{DeviceID: observation.DeviceID, Reason: "stale observation"})
		case observation.Quality < 0.5:
			summary.Drops = append(summary.Drops, Drop{DeviceID: observation.DeviceID, Reason: "low quality"})
		case observation.Latitude < -90 || observation.Latitude > 90 || observation.Longitude < -180 || observation.Longitude > 180:
			summary.Drops = append(summary.Drops, Drop{DeviceID: observation.DeviceID, Reason: "invalid coordinates"})
		default:
			summary.Events = append(summary.Events, TelemetryEvent{
				AgencyID:    observation.AgencyID,
				DeviceID:    observation.DeviceID,
				VehicleID:   observation.VehicleID,
				ObservedAt:  observedAt.Format(time.RFC3339),
				Latitude:    observation.Latitude,
				Longitude:   observation.Longitude,
				Source:      "synthetic-http-poller",
				DryRun:      true,
				NetworkSend: false,
			})
		}
	}
	return summary
}
