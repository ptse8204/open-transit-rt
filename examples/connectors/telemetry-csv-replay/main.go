package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

const defaultCSVFixture = "examples/connectors/telemetry-csv-replay/fixtures/replay.csv"

type Event struct {
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

func main() {
	path := defaultCSVFixture
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	events, err := ParseReplay(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(events); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ParseReplay(r io.Reader) ([]Event, error) {
	reader := csv.NewReader(r)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, errors.New("csv replay requires a header and at least one data row")
	}
	if len(rows[0]) != 7 {
		return nil, errors.New("csv replay header must have seven columns")
	}

	events := make([]Event, 0, len(rows)-1)
	for i, row := range rows[1:] {
		if len(row) != 7 {
			return nil, fmt.Errorf("row %d: expected seven columns", i+2)
		}
		if row[0] != "true" {
			return nil, fmt.Errorf("row %d: synthetic_only must be true", i+2)
		}
		observedAt, err := time.Parse(time.RFC3339, row[4])
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid observed_at: %w", i+2, err)
		}
		latitude, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid latitude: %w", i+2, err)
		}
		longitude, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid longitude: %w", i+2, err)
		}
		if row[1] == "" || row[2] == "" || row[3] == "" {
			return nil, fmt.Errorf("row %d: missing identity", i+2)
		}
		events = append(events, Event{
			AgencyID:    row[1],
			DeviceID:    row[2],
			VehicleID:   row[3],
			ObservedAt:  observedAt.Format(time.RFC3339),
			Latitude:    latitude,
			Longitude:   longitude,
			Source:      "synthetic-csv-replay",
			DryRun:      true,
			NetworkSend: false,
		})
	}
	return events, nil
}
