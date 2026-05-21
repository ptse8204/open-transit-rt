package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"open-transit-rt/internal/auth"
	connectorpkg "open-transit-rt/internal/connectors"
)

func TestConnectorHubSeparatesExamplesFromConfiguredInstances(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02ReadOnlyPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("connectors status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Configured Connector Instances",
		string(connectorpkg.StateExampleAvailable),
		"examples are not configuration",
		"No configured telemetry source instance.",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("connector page missing %q: %s", want, body)
		}
	}
	if strings.Contains(body, "configured telemetry source instance(s)") {
		t.Fatalf("example manifests were presented as configured instances: %s", body)
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/connectors.json", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("connectors json status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var hub connectorHubView
	if err := json.Unmarshal(rr.Body.Bytes(), &hub); err != nil {
		t.Fatalf("decode connector hub json: %v: %s", err, rr.Body.String())
	}
	if hub.InstanceSummary.ConfiguredInstances != 0 {
		t.Fatalf("configured instance count = %d, want 0", hub.InstanceSummary.ConfiguredInstances)
	}
	if hub.InstanceSummary.ExampleManifests == 0 {
		t.Fatalf("expected example manifests in summary: %+v", hub.InstanceSummary)
	}
	for _, row := range hub.Instances {
		if row.State != string(connectorpkg.StateExampleAvailable) && row.State != string(connectorpkg.StateNotConfigured) {
			t.Fatalf("unexpected configured state from examples-only hub: %+v", row)
		}
	}
}

func TestConnectorHubShowsConfiguredInstanceStateWithoutSecrets(t *testing.T) {
	checked := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	store := fakeConnectorInstanceStore{instances: []connectorpkg.Instance{{
		ID:            42,
		AgencyID:      "demo-agency",
		ConnectorType: connectorpkg.TypeTelemetrySource,
		ConnectorKind: "http_polling",
		DisplayName:   "Agency AVL poller",
		State:         connectorpkg.StateConfiguredNotTested,
		Owner:         "ops@example.com",
		ConfigJSON:    json.RawMessage(`{"field_map":"standard_vehicle_observation","source_label":"depot"}`),
		SecretRefs:    []string{"AVL_HTTP_TOKEN_REF"},
		DryRunStatus:  "not_run",
		LastCheckedAt: &checked,
	}}}
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: store}, auth.TestAuthenticator{Principal: phase02ReadOnlyPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("connectors status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Agency AVL poller",
		string(connectorpkg.StateConfiguredNotTested),
		"metadata keys: field_map, source_label",
		"AVL_HTTP_TOKEN_REF",
		"configured telemetry source instance(s): Agency AVL poller",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("connector page missing configured instance marker %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"super-secret", "raw_payload", "password="} {
		if strings.Contains(strings.ToLower(body), forbidden) {
			t.Fatalf("connector page leaked forbidden value %q: %s", forbidden, body)
		}
	}
}

type fakeConnectorInstanceStore struct {
	instances []connectorpkg.Instance
	err       error
}

func (f fakeConnectorInstanceStore) ListInstances(context.Context, string) ([]connectorpkg.Instance, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.instances, nil
}
