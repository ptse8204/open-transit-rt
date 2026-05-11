package monitoringsdk

import "testing"

func TestBuildNoSendBatchRedactsAndDoesNotMutate(t *testing.T) {
	input := Input{
		AgencyID: "agency-demo",
		Metrics:  map[string]int{"fresh": 12},
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
	for _, field := range batch.RedactedFields {
		if field == "operator_contact" {
			return
		}
	}
	t.Fatalf("redacted fields = %v, want operator_contact", batch.RedactedFields)
}
