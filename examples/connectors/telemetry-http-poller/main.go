package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	telemetrysdk "open-transit-rt/examples/connectors/sdk/telemetry"
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

func Transform(observations []Observation, now time.Time) telemetrysdk.Summary {
	var parsed []telemetrysdk.Observation
	var summary telemetrysdk.Summary
	for _, observation := range observations {
		observedAt, err := time.Parse(time.RFC3339, observation.Timestamp)
		if err != nil {
			summary.Drops = append(summary.Drops, telemetrysdk.Drop{DeviceID: observation.DeviceID, VehicleID: observation.VehicleID, Reason: telemetrysdk.ReasonInvalidTimestamp})
			continue
		}
		parsed = append(parsed, telemetrysdk.Observation{
			AgencyID:   observation.AgencyID,
			DeviceID:   observation.DeviceID,
			VehicleID:  observation.VehicleID,
			ObservedAt: observedAt,
			Latitude:   observation.Latitude,
			Longitude:  observation.Longitude,
			Quality:    observation.Quality,
		})
	}
	normalized := telemetrysdk.NormalizeBatch(parsed, telemetrysdk.Options{Now: now, Source: "synthetic-http-poller"})
	summary.Events = append(summary.Events, normalized.Events...)
	summary.Drops = append(summary.Drops, normalized.Drops...)
	return summary
}
