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

	telemetrysdk "open-transit-rt/examples/connectors/sdk/telemetry"
)

const defaultCSVFixture = "examples/connectors/telemetry-csv-replay/fixtures/replay.csv"

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

	summary, err := ParseReplay(file, time.Date(2026, 5, 10, 16, 1, 0, 0, time.UTC))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ParseReplay(r io.Reader, now time.Time) (telemetrysdk.Summary, error) {
	reader := csv.NewReader(r)
	rows, err := reader.ReadAll()
	if err != nil {
		return telemetrysdk.Summary{}, err
	}
	if len(rows) < 2 {
		return telemetrysdk.Summary{}, errors.New("csv replay requires a header and at least one data row")
	}
	if len(rows[0]) != 8 {
		return telemetrysdk.Summary{}, errors.New("csv replay header must have eight columns")
	}

	var parsed []telemetrysdk.Observation
	var summary telemetrysdk.Summary
	for i, row := range rows[1:] {
		if len(row) != 8 {
			return telemetrysdk.Summary{}, fmt.Errorf("row %d: expected eight columns", i+2)
		}
		if row[0] != "true" {
			return telemetrysdk.Summary{}, fmt.Errorf("row %d: synthetic_only must be true", i+2)
		}
		observedAt, err := time.Parse(time.RFC3339, row[4])
		if err != nil {
			summary.Drops = append(summary.Drops, telemetrysdk.Drop{DeviceID: row[2], VehicleID: row[3], Reason: telemetrysdk.ReasonInvalidTimestamp})
			continue
		}
		latitude, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			summary.Drops = append(summary.Drops, telemetrysdk.Drop{DeviceID: row[2], VehicleID: row[3], Reason: telemetrysdk.ReasonInvalidCoordinate})
			continue
		}
		longitude, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			summary.Drops = append(summary.Drops, telemetrysdk.Drop{DeviceID: row[2], VehicleID: row[3], Reason: telemetrysdk.ReasonInvalidCoordinate})
			continue
		}
		quality, err := strconv.ParseFloat(row[7], 64)
		if err != nil {
			summary.Drops = append(summary.Drops, telemetrysdk.Drop{DeviceID: row[2], VehicleID: row[3], Reason: telemetrysdk.ReasonLowQuality})
			continue
		}
		parsed = append(parsed, telemetrysdk.Observation{
			AgencyID:   row[1],
			DeviceID:   row[2],
			VehicleID:  row[3],
			ObservedAt: observedAt,
			Latitude:   latitude,
			Longitude:  longitude,
			Quality:    quality,
		})
	}
	normalized := telemetrysdk.NormalizeBatch(parsed, telemetrysdk.Options{Now: now, Source: "synthetic-csv-replay"})
	summary.Events = append(summary.Events, normalized.Events...)
	summary.Drops = append(summary.Drops, normalized.Drops...)
	return summary, nil
}
