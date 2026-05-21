package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"open-transit-rt/internal/auth"
	connectorpkg "open-transit-rt/internal/connectors"
)

func TestVehicleAVLSetupRendersAdminMetadataForm(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/vehicle-avl", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("vehicle avl status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Vehicle / GPS / AVL Setup",
		`<form method="post" action="/admin/operations/connectors/vehicle-avl#vehicle-avl-form">`,
		`name="action" value="save_vehicle_avl_connector"`,
		`id="vehicle_avl_source_shape" name="source_shape"`,
		`value="generic_json_transform"`,
		`id="vehicle_avl_secret_ref" name="secret_ref"`,
		`name="field_device_id"`,
		`name="field_vehicle_id"`,
		`name="field_observed_timestamp"`,
		`name="field_lat"`,
		`name="field_lon"`,
		`<button type="submit">Save vehicle connector metadata</button>`,
		"Dry-run is required and remains server-owned",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("vehicle avl page missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"validator_command", "connector command", "http://", "https://"} {
		if strings.Contains(strings.ToLower(body), forbidden) {
			t.Fatalf("vehicle avl page exposed unsafe copy %q: %s", forbidden, body)
		}
	}
}

func TestVehicleAVLSetupReadOnlyHasNoMutationForm(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02ReadOnlyPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/vehicle-avl", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("vehicle avl status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if strings.Contains(body, `<form`) {
		t.Fatalf("read-only vehicle avl page rendered mutation form: %s", body)
	}
	if !strings.Contains(body, "changes require an admin role") {
		t.Fatalf("read-only vehicle avl page missing role boundary: %s", body)
	}
}

func TestVehicleAVLSetupPostStoresMetadataOnly(t *testing.T) {
	store := &fakeVehicleAVLConnectorStore{}
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: store}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	form := vehicleAVLValidForm()
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/vehicle-avl", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("vehicle avl post status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.saved.ConnectorType != connectorpkg.TypeTelemetrySource {
		t.Fatalf("connector type = %q", store.saved.ConnectorType)
	}
	if store.saved.State != connectorpkg.StateConfiguredNotTested {
		t.Fatalf("state = %q", store.saved.State)
	}
	if got := strings.Join(store.saved.SecretRefs, ","); got != "AVL_HTTP_TOKEN_REF" {
		t.Fatalf("secret refs = %q", got)
	}
	var metadata vehicleAVLConfigMetadata
	if err := json.Unmarshal(store.saved.ConfigJSON, &metadata); err != nil {
		t.Fatalf("decode saved metadata: %v", err)
	}
	if metadata.SourceShape != "generic_json_transform" {
		t.Fatalf("source shape = %q", metadata.SourceShape)
	}
	if metadata.FieldMap["vehicle_id"] != "vehicle.id" || metadata.FieldMap["lat"] != "position.lat" {
		t.Fatalf("field map = %+v", metadata.FieldMap)
	}
	raw := strings.ToLower(string(store.saved.ConfigJSON))
	for _, forbidden := range []string{"password", "secret", "token=", "http://", "https://"} {
		if strings.Contains(raw, forbidden) {
			t.Fatalf("saved config metadata leaked %q: %s", forbidden, raw)
		}
	}
	if !strings.Contains(rr.Body.String(), "Dry-run is still required before activation") {
		t.Fatalf("response missing dry-run reminder: %s", rr.Body.String())
	}
}

func TestVehicleAVLSetupPostRejectsUnsafeMetadata(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: &fakeVehicleAVLConnectorStore{}}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	for name, mutate := range map[string]func(url.Values){
		"bad secret ref": func(v url.Values) { v.Set("secret_ref", "raw-token-value") },
		"bad field map":  func(v url.Values) { v.Set("field_vehicle_id", "payload.api_token") },
	} {
		t.Run(name, func(t *testing.T) {
			form := vehicleAVLValidForm()
			mutate(form)
			req := httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/vehicle-avl", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want validation page: %s", rr.Code, rr.Body.String())
			}
			if !strings.Contains(rr.Body.String(), "bad") && !strings.Contains(rr.Body.String(), "invalid") && !strings.Contains(rr.Body.String(), "secret") {
				t.Fatalf("response missing validation error: %s", rr.Body.String())
			}
		})
	}
}

func vehicleAVLValidForm() url.Values {
	return url.Values{
		"action":                   {"save_vehicle_avl_connector"},
		"display_name":             {"Agency AVL poller"},
		"owner":                    {"ops@example.org"},
		"source_shape":             {"generic_json_transform"},
		"secret_ref":               {"AVL_HTTP_TOKEN_REF"},
		"field_agency_id":          {"agency_id"},
		"field_device_id":          {"device.id"},
		"field_vehicle_id":         {"vehicle.id"},
		"field_observed_timestamp": {"observed_at"},
		"field_lat":                {"position.lat"},
		"field_lon":                {"position.lon"},
		"field_quality":            {"quality"},
		"field_route_hint":         {"route_id"},
		"field_trip_hint":          {"trip_id"},
		"field_speed":              {"speed"},
		"field_bearing":            {"bearing"},
		"field_accuracy":           {"accuracy_meters"},
	}
}

type fakeVehicleAVLConnectorStore struct {
	instances []connectorpkg.Instance
	saved     connectorpkg.UpsertInstanceInput
}

func (f *fakeVehicleAVLConnectorStore) ListInstances(context.Context, string) ([]connectorpkg.Instance, error) {
	return f.instances, nil
}

func (f *fakeVehicleAVLConnectorStore) UpsertInstance(_ context.Context, input connectorpkg.UpsertInstanceInput) (connectorpkg.Instance, error) {
	f.saved = input
	instance := connectorpkg.Instance{
		ID:            101,
		AgencyID:      input.AgencyID,
		ConnectorType: input.ConnectorType,
		ConnectorKind: input.ConnectorKind,
		DisplayName:   input.DisplayName,
		State:         input.State,
		Owner:         input.Owner,
		ConfigJSON:    input.ConfigJSON,
		SecretRefs:    input.SecretRefs,
		DryRunStatus:  input.DryRunStatus,
	}
	f.instances = append(f.instances, instance)
	return instance, nil
}
