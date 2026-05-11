package monitoringsdk

type Input struct {
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
	StatusMutation bool             `json:"status_mutation"`
	EvidenceWrite  bool             `json:"evidence_write"`
	Metrics        map[string]int   `json:"metrics"`
	Incidents      []IncidentOutput `json:"incidents"`
	RedactedFields []string         `json:"redacted_fields"`
}

type IncidentOutput struct {
	ID       string `json:"id"`
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
}

func BuildNoSendBatch(input Input) ExportBatch {
	batch := ExportBatch{
		AgencyID:       input.AgencyID,
		SendEnabled:    false,
		NetworkSend:    false,
		StatusMutation: false,
		EvidenceWrite:  false,
		Metrics:        copyMetrics(input.Metrics),
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

func copyMetrics(metrics map[string]int) map[string]int {
	if len(metrics) == 0 {
		return nil
	}
	out := make(map[string]int, len(metrics))
	for key, value := range metrics {
		out[key] = value
	}
	return out
}
