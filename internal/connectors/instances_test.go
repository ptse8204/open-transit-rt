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
