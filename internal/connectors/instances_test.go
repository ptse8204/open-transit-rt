package connectors

import (
	"encoding/json"
	"testing"
)

func TestConnectorInstanceStatesAreExplicitAndStable(t *testing.T) {
	want := []InstanceState{
		StateExampleAvailable,
		StateNotConfigured,
		StateConfiguredNotTested,
		StateDryRunPassed,
		StateReadyForActivation,
		StateActive,
		StateBlocked,
	}
	if len(SupportedInstanceStates) != len(want) {
		t.Fatalf("supported states = %v, want %v", SupportedInstanceStates, want)
	}
	for i, state := range want {
		if SupportedInstanceStates[i] != state {
			t.Fatalf("state[%d] = %q, want %q", i, SupportedInstanceStates[i], state)
		}
		parsed, err := ParseInstanceState(string(state))
		if err != nil {
			t.Fatalf("parse state %s: %v", state, err)
		}
		if parsed != state {
			t.Fatalf("parsed state = %q, want %q", parsed, state)
		}
	}
	if _, err := ParseInstanceState("configured"); err == nil {
		t.Fatalf("legacy configured state must be rejected")
	}
}

func TestConnectorInstanceConfigSummaryIsMetadataOnly(t *testing.T) {
	inst := Instance{
		ConfigJSON: json.RawMessage(`{"field_map":"vehicle_timestamp","source_label":"depot poller"}`),
		SecretRefs: []string{
			"  AVL_HTTP_TOKEN_REF  ",
			"",
			"DEVICE_TOKEN_REF",
		},
	}
	if !inst.DeploymentConfigExists() {
		t.Fatalf("deployment config should exist when metadata and secret refs are present")
	}
	keys := inst.ConfigKeys()
	if got, want := len(keys), 2; got != want {
		t.Fatalf("config key count = %d, want %d: %v", got, want, keys)
	}
	if keys[0] != "field_map" || keys[1] != "source_label" {
		t.Fatalf("config keys = %v", keys)
	}
	clean := cleanStringList(inst.SecretRefs)
	if got, want := len(clean), 2; got != want {
		t.Fatalf("secret refs count = %d, want %d: %v", got, want, clean)
	}
	if clean[0] != "AVL_HTTP_TOKEN_REF" || clean[1] != "DEVICE_TOKEN_REF" {
		t.Fatalf("secret refs = %v", clean)
	}
}

func TestConnectorInstanceUpsertValidationRejectsSecretsAndEndpoints(t *testing.T) {
	base := UpsertInstanceInput{
		AgencyID:      "demo-agency",
		ConnectorType: TypeTelemetrySource,
		ConnectorKind: "http_polling",
		DisplayName:   "Agency AVL poller",
		State:         StateConfiguredNotTested,
		ConfigJSON:    json.RawMessage(`{"source_shape":"http_polling","field_map":{"vehicle_id":"vehicle.id"}}`),
		SecretRefs:    []string{"AVL_HTTP_TOKEN_REF"},
		DryRunStatus:  "not_run",
	}
	if _, err := normalizeUpsertInstanceInput(base); err != nil {
		t.Fatalf("valid input rejected: %v", err)
	}
	for name, input := range map[string]UpsertInstanceInput{
		"endpoint": func() UpsertInstanceInput {
			next := base
			next.ConfigJSON = json.RawMessage(`{"endpoint":"https://private.example.test/avl"}`)
			return next
		}(),
		"secret_key": func() UpsertInstanceInput {
			next := base
			next.ConfigJSON = json.RawMessage(`{"api_token":"inline"}`)
			return next
		}(),
		"bad_secret_ref": func() UpsertInstanceInput {
			next := base
			next.SecretRefs = []string{"not-a-secret"}
			return next
		}(),
		"bad_state": func() UpsertInstanceInput {
			next := base
			next.State = InstanceState("configured")
			return next
		}(),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := normalizeUpsertInstanceInput(input); err == nil {
				t.Fatalf("unsafe input was accepted: %+v", input)
			}
		})
	}
}
