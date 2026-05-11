package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	predictionsdk "open-transit-rt/examples/connectors/sdk/prediction"
)

const defaultPredictionFixture = "examples/connectors/predictor-sidecar-stub/fixtures/prediction-input.json"

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

func readRequest(path string) (predictionsdk.Request, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return predictionsdk.Request{}, err
	}
	var request predictionsdk.Request
	if err := json.Unmarshal(raw, &request); err != nil {
		return predictionsdk.Request{}, err
	}
	if !request.SyntheticOnly {
		return predictionsdk.Request{}, errors.New("prediction request fixture must be marked synthetic_only")
	}
	if request.ActiveFeedVersion == "" || request.VehiclePositionsURL == "" {
		return predictionsdk.Request{}, errors.New("prediction request missing feed version or vehicle positions reference")
	}
	return request, nil
}

func BuildResponse(request predictionsdk.Request, threshold float64) predictionsdk.Response {
	return predictionsdk.BuildDryRunResponse(request, predictionsdk.Options{ConfidenceThreshold: threshold})
}
