package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const defaultMetricsFixture = "examples/connectors/monitoring-export/fixtures/metrics.json"

type MetricsInput struct {
	SyntheticOnly bool            `json:"synthetic_only"`
	AgencyID      string          `json:"agency_id"`
	Metrics       map[string]int  `json:"metrics"`
	Incidents     []IncidentInput `json:"incidents"`
}

type IncidentInput struct {
	ID              string `json:"id"`
	Severity        string `json:"severity"`
	Summary         string `json:"summary"`
	OperatorContact string `json:"operator_contact,omitempty"`
	RawPayload      string `json:"raw_payload,omitempty"`
}

type ExportBatch struct {
	AgencyID       string           `json:"agency_id"`
	SendEnabled    bool             `json:"send_enabled"`
	NetworkSend    bool             `json:"network_send"`
	Metrics        map[string]int   `json:"metrics"`
	Incidents      []IncidentOutput `json:"incidents"`
	RedactedFields []string         `json:"redacted_fields"`
}

type IncidentOutput struct {
	ID       string `json:"id"`
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
}

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

func readMetrics(path string) (MetricsInput, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return MetricsInput{}, err
	}
	var input MetricsInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return MetricsInput{}, err
	}
	if !input.SyntheticOnly {
		return MetricsInput{}, errors.New("metrics fixture must be marked synthetic_only")
	}
	if input.AgencyID == "" {
		return MetricsInput{}, errors.New("metrics input missing agency_id")
	}
	return input, nil
}

func BuildExportBatch(input MetricsInput) ExportBatch {
	batch := ExportBatch{
		AgencyID:       input.AgencyID,
		SendEnabled:    false,
		NetworkSend:    false,
		Metrics:        input.Metrics,
		RedactedFields: []string{"operator_contact", "raw_payload"},
	}
	for _, incident := range input.Incidents {
		batch.Incidents = append(batch.Incidents, IncidentOutput{
			ID:       incident.ID,
			Severity: incident.Severity,
			Summary:  incident.Summary,
		})
	}
	return batch
}
