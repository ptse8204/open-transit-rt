package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	monitoringsdk "open-transit-rt/examples/connectors/sdk/monitoring"
)

const defaultMetricsFixture = "examples/connectors/monitoring-export/fixtures/metrics.json"

func main() {
	path := defaultMetricsFixture
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	input, err := readMetrics(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	batch := BuildExportBatch(input)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(batch); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readMetrics(path string) (monitoringsdk.Input, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return monitoringsdk.Input{}, err
	}
	var input monitoringsdk.Input
	if err := json.Unmarshal(raw, &input); err != nil {
		return monitoringsdk.Input{}, err
	}
	if !input.SyntheticOnly {
		return monitoringsdk.Input{}, errors.New("metrics fixture must be marked synthetic_only")
	}
	if input.AgencyID == "" {
		return monitoringsdk.Input{}, errors.New("metrics input missing agency_id")
	}
	return input, nil
}

func BuildExportBatch(input monitoringsdk.Input) monitoringsdk.ExportBatch {
	return monitoringsdk.BuildNoSendBatch(input)
}
