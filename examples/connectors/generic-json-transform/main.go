package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	telemetrysdk "open-transit-rt/examples/connectors/sdk/telemetry"
)

const defaultFixture = "examples/connectors/generic-json-transform/fixtures/observations.json"

type Fixture struct {
	SyntheticOnly  bool              `json:"synthetic_only"`
	SendEnabled    bool              `json:"send_enabled"`
	TimeoutSeconds int               `json:"timeout_seconds"`
	FieldMap       map[string]string `json:"field_map"`
	Records        []map[string]any  `json:"records"`
}

type Result struct {
	DryRun         bool                 `json:"dry_run"`
	NetworkSend    bool                 `json:"network_send"`
	TimeoutSeconds int                  `json:"timeout_seconds"`
	Diagnostics    []string             `json:"diagnostics,omitempty"`
	Events         []telemetrysdk.Event `json:"events"`
	Drops          []telemetrysdk.Drop  `json:"drops"`
}

func main() {
	path := defaultFixture
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	fixture, err := readFixture(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	result := Transform(fixture, time.Date(2026, 5, 10, 15, 1, 0, 0, time.UTC))
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readFixture(path string) (Fixture, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Fixture{}, err
	}
	var fixture Fixture
	if err := json.Unmarshal(raw, &fixture); err != nil {
		return Fixture{}, err
	}
	if err := validateFixture(fixture); err != nil {
		return Fixture{}, err
	}
	return fixture, nil
}

func validateFixture(fixture Fixture) error {
	if !fixture.SyntheticOnly {
		return errors.New("fixture must be marked synthetic_only")
	}
	if fixture.SendEnabled {
		return errors.New("send_enabled must stay false for this dry-run connector")
	}
	if fixture.TimeoutSeconds <= 0 || fixture.TimeoutSeconds > 30 {
		return errors.New("timeout_seconds must be between 1 and 30")
	}
	for _, key := range []string{"agency_id", "device_id", "vehicle_id", "timestamp", "latitude", "longitude", "quality"} {
		if fixture.FieldMap[key] == "" {
			return fmt.Errorf("field_map.%s is required", key)
		}
	}
	if len(fixture.Records) == 0 {
		return errors.New("records must not be empty")
	}
	return nil
}

func Transform(fixture Fixture, now time.Time) Result {
	result := Result{DryRun: true, NetworkSend: false, TimeoutSeconds: fixture.TimeoutSeconds}
	var observations []telemetrysdk.Observation
	for _, record := range fixture.Records {
		observation, drop, ok := mapRecord(fixture.FieldMap, record)
		if !ok {
			result.Drops = append(result.Drops, drop)
			continue
		}
		observations = append(observations, observation)
	}
	normalized := telemetrysdk.NormalizeBatch(observations, telemetrysdk.Options{Now: now, Source: "synthetic-generic-json-transform"})
	result.Events = append(result.Events, normalized.Events...)
	result.Drops = append(result.Drops, normalized.Drops...)
	if len(result.Events) == 0 {
		result.Diagnostics = append(result.Diagnostics, "no telemetry emitted")
	}
	return result
}

func mapRecord(fieldMap map[string]string, record map[string]any) (telemetrysdk.Observation, telemetrysdk.Drop, bool) {
	deviceID := stringField(record, fieldMap["device_id"])
	vehicleID := stringField(record, fieldMap["vehicle_id"])
	observedAt, err := time.Parse(time.RFC3339, stringField(record, fieldMap["timestamp"]))
	if err != nil {
		return telemetrysdk.Observation{}, telemetrysdk.Drop{DeviceID: deviceID, VehicleID: vehicleID, Reason: telemetrysdk.ReasonInvalidTimestamp}, false
	}
	lat, ok := numberField(record, fieldMap["latitude"])
	if !ok {
		return telemetrysdk.Observation{}, telemetrysdk.Drop{DeviceID: deviceID, VehicleID: vehicleID, Reason: telemetrysdk.ReasonInvalidCoordinate}, false
	}
	lon, ok := numberField(record, fieldMap["longitude"])
	if !ok {
		return telemetrysdk.Observation{}, telemetrysdk.Drop{DeviceID: deviceID, VehicleID: vehicleID, Reason: telemetrysdk.ReasonInvalidCoordinate}, false
	}
	quality, ok := numberField(record, fieldMap["quality"])
	if !ok {
		quality = 0
	}
	return telemetrysdk.Observation{
		AgencyID:   stringField(record, fieldMap["agency_id"]),
		DeviceID:   deviceID,
		VehicleID:  vehicleID,
		ObservedAt: observedAt,
		Latitude:   lat,
		Longitude:  lon,
		Quality:    quality,
	}, telemetrysdk.Drop{}, true
}

func stringField(record map[string]any, key string) string {
	value, _ := record[key].(string)
	return value
}

func numberField(record map[string]any, key string) (float64, bool) {
	switch value := record[key].(type) {
	case float64:
		return value, true
	case int:
		return float64(value), true
	case string:
		parsed, err := strconv.ParseFloat(value, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}
