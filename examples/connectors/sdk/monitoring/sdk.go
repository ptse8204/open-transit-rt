package monitoringsdk

import "strings"

type Input struct {
	SyntheticOnly      bool            `json:"synthetic_only"`
	AgencyID           string          `json:"agency_id"`
	Metrics            map[string]int  `json:"metrics"`
	Incidents          []IncidentInput `json:"incidents"`
	FeedHealth         []HealthInput   `json:"feed_health,omitempty"`
	ConnectorHealth    []HealthInput   `json:"connector_health,omitempty"`
	ValidatorPosture   []HealthInput   `json:"validator_posture,omitempty"`
	TelemetryFreshness []HealthInput   `json:"telemetry_freshness,omitempty"`
	MaintenanceTasks   []HealthInput   `json:"maintenance_tasks,omitempty"`
}

type IncidentInput struct {
	ID              string `json:"id"`
	Severity        string `json:"severity"`
	Summary         string `json:"summary"`
	OperatorContact string `json:"operator_contact,omitempty"`
	RawPayload      string `json:"raw_payload,omitempty"`
}

type HealthInput struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Summary    string `json:"summary"`
	NextAction string `json:"next_action"`
	PrivateRef string `json:"private_ref,omitempty"`
	RawPayload string `json:"raw_payload,omitempty"`
}

type ExportBatch struct {
	AgencyID           string           `json:"agency_id"`
	SendEnabled        bool             `json:"send_enabled"`
	NetworkSend        bool             `json:"network_send"`
	StatusMutation     bool             `json:"status_mutation"`
	EvidenceWrite      bool             `json:"evidence_write"`
	Metrics            map[string]int   `json:"metrics"`
	Incidents          []IncidentOutput `json:"incidents"`
	FeedHealth         []HealthOutput   `json:"feed_health,omitempty"`
	ConnectorHealth    []HealthOutput   `json:"connector_health,omitempty"`
	ValidatorPosture   []HealthOutput   `json:"validator_posture,omitempty"`
	TelemetryFreshness []HealthOutput   `json:"telemetry_freshness,omitempty"`
	MaintenanceTasks   []HealthOutput   `json:"maintenance_tasks,omitempty"`
	RedactedFields     []string         `json:"redacted_fields"`
}

type IncidentOutput struct {
	ID       string `json:"id"`
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
}

type HealthOutput struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Summary    string `json:"summary"`
	NextAction string `json:"next_action"`
}

func BuildNoSendBatch(input Input) ExportBatch {
	batch := ExportBatch{
		AgencyID:           input.AgencyID,
		SendEnabled:        false,
		NetworkSend:        false,
		StatusMutation:     false,
		EvidenceWrite:      false,
		Metrics:            copyMetrics(input.Metrics),
		FeedHealth:         copyHealthRows(input.FeedHealth),
		ConnectorHealth:    copyHealthRows(input.ConnectorHealth),
		ValidatorPosture:   copyHealthRows(input.ValidatorPosture),
		TelemetryFreshness: copyHealthRows(input.TelemetryFreshness),
		MaintenanceTasks:   copyHealthRows(input.MaintenanceTasks),
		RedactedFields:     []string{"operator_contact", "raw_payload", "private_ref", "private_endpoint", "destination_value", "token"},
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

func copyHealthRows(rows []HealthInput) []HealthOutput {
	if len(rows) == 0 {
		return nil
	}
	out := make([]HealthOutput, 0, len(rows))
	for _, row := range rows {
		out = append(out, HealthOutput{
			ID:         safeMonitoringText(row.ID, "row"),
			Status:     safeMonitoringText(row.Status, "unknown"),
			Summary:    safeMonitoringText(row.Summary, "review private diagnostics"),
			NextAction: safeMonitoringText(row.NextAction, "review private diagnostics"),
		})
	}
	return out
}

func safeMonitoringText(value string, fallback string) string {
	text := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(value, "\r", " "), "\n", " "))
	if text == "" {
		return fallback
	}
	lower := strings.ToLower(text)
	for _, forbidden := range []string{
		"authorization:",
		"bearer ",
		"set-cookie",
		"database_url",
		"postgres://",
		"token_hash",
		"payload_json",
		"raw_payload",
		"webhook",
		"private_endpoint",
		"/users/",
		"/home/",
		"/var/",
		"/tmp/",
		"/etc/",
	} {
		if strings.Contains(lower, forbidden) {
			return "<redacted>"
		}
	}
	if len(text) > 220 {
		return strings.TrimSpace(text[:205]) + " [truncated]"
	}
	return text
}
