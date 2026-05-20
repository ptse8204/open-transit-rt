package monitoringsdk

import "testing"

func TestBuildNoSendBatchRedactsAndDoesNotMutate(t *testing.T) {
	input := Input{
		AgencyID: "agency-demo",
		Metrics:  map[string]int{"fresh": 12},
		FeedHealth: []HealthInput{{
			ID:         "vehicle_positions",
			Status:     "needs_review",
			Summary:    "Vehicle Positions stale bucket needs review",
			NextAction: "Open Feed Health and review freshness",
		}},
		ConnectorHealth: []HealthInput{{
			ID:         "telemetry_source",
			Status:     "needs_review",
			Summary:    "Do not copy https://webhook.example/private",
			NextAction: "Review private connector setup",
		}},
		Incidents: []IncidentInput{{
			ID:              "incident-1",
			Severity:        "warning",
			Summary:         "Synthetic incident",
			OperatorContact: "redacted@example.invalid",
			RawPayload:      "{\"synthetic\":true}",
		}},
	}
	batch := BuildNoSendBatch(input)
	if batch.SendEnabled || batch.NetworkSend || batch.StatusMutation || batch.EvidenceWrite {
		t.Fatalf("batch mutates or sends: %+v", batch)
	}
	if len(batch.Incidents) != 1 {
		t.Fatalf("incidents = %d, want 1", len(batch.Incidents))
	}
	if batch.Incidents[0].ID != "incident-1" || batch.Incidents[0].Summary == "" {
		t.Fatalf("incident public fields not preserved: %+v", batch.Incidents[0])
	}
	if len(batch.FeedHealth) != 1 || batch.FeedHealth[0].ID != "vehicle_positions" {
		t.Fatalf("feed health export = %+v", batch.FeedHealth)
	}
	if len(batch.ConnectorHealth) != 1 || batch.ConnectorHealth[0].Summary != "<redacted>" {
		t.Fatalf("connector health redaction = %+v", batch.ConnectorHealth)
	}
	for _, field := range batch.RedactedFields {
		if field == "operator_contact" {
			return
		}
	}
	t.Fatalf("redacted fields = %v, want operator_contact", batch.RedactedFields)
}
