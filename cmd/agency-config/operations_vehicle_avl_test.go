package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"open-transit-rt/internal/auth"
	connectorpkg "open-transit-rt/internal/connectors"
	"open-transit-rt/internal/devices"
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

func TestVehicleAVLDryRunPostStoresRedactedResult(t *testing.T) {
	store := &fakeVehicleAVLConnectorStore{instances: []connectorpkg.Instance{{
		ID:            101,
		AgencyID:      "demo-agency",
		ConnectorType: connectorpkg.TypeTelemetrySource,
		ConnectorKind: "generic_json_transform",
		DisplayName:   "Agency AVL poller",
		State:         connectorpkg.StateConfiguredNotTested,
		DryRunStatus:  "not_run",
	}}}
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: store}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	form := url.Values{
		"action":                {"record_vehicle_avl_dry_run"},
		"connector_instance_id": {"101"},
		"dry_run_status":        {"passed"},
		"redaction_scan_status": {"passed"},
		"accepted_count":        {"2"},
		"rejected_count":        {"1"},
		"dropped_count":         {"0"},
		"redacted_summary":      {"Synthetic fixture accepted two rows; no raw payload retained"},
	}
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/vehicle-avl", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("dry-run post status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.savedDryRun.ConnectorInstanceID != 101 || store.savedDryRun.Status != "passed" || store.savedDryRun.AcceptedCount != 2 {
		t.Fatalf("saved dry-run = %+v", store.savedDryRun)
	}
	raw := strings.ToLower(string(store.savedDryRun.RedactedSummary))
	for _, forbidden := range []string{"http://", "https://", "password=", "token=", "secret=", "raw_payload\":true"} {
		if strings.Contains(raw, forbidden) {
			t.Fatalf("dry-run summary leaked forbidden value %q: %s", forbidden, raw)
		}
	}
	body := rr.Body.String()
	for _, want := range []string{"dry-run result was recorded", string(connectorpkg.StateDryRunPassed), "accepted=2 rejected=1 dropped=0"} {
		if !strings.Contains(body, want) {
			t.Fatalf("dry-run response missing %q: %s", want, body)
		}
	}
}

func TestVehicleAVLDryRunPostRejectsRawSummary(t *testing.T) {
	store := &fakeVehicleAVLConnectorStore{instances: []connectorpkg.Instance{{
		ID:            101,
		AgencyID:      "demo-agency",
		ConnectorType: connectorpkg.TypeTelemetrySource,
		ConnectorKind: "generic_json_transform",
		DisplayName:   "Agency AVL poller",
		State:         connectorpkg.StateConfiguredNotTested,
		DryRunStatus:  "not_run",
	}}}
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: store}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	form := url.Values{
		"action":                {"record_vehicle_avl_dry_run"},
		"connector_instance_id": {"101"},
		"dry_run_status":        {"passed"},
		"redaction_scan_status": {"passed"},
		"redacted_summary":      {"payload sent to https://private.example.test with token=value"},
	}
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/vehicle-avl", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("dry-run validation status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "must not contain endpoints or inline secrets") {
		t.Fatalf("response missing raw summary validation error: %s", rr.Body.String())
	}
	if store.savedDryRun.ConnectorInstanceID != 0 {
		t.Fatalf("unsafe dry-run was saved: %+v", store.savedDryRun)
	}
}

func TestVehicleAVLActivationGateMarksReadyOnlyAfterChecksPass(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	config := json.RawMessage(`{"source_shape":"generic_json_transform","field_map":{"agency_id":"agency_id","device_id":"device.id","vehicle_id":"vehicle.id","observed_timestamp":"observed_at","lat":"position.lat","lon":"position.lon","quality":"quality"},"safety":{"target_path":"/v1/telemetry"}}`)
	store := &fakeVehicleAVLConnectorStore{
		instances: []connectorpkg.Instance{{
			ID:            101,
			AgencyID:      "demo-agency",
			ConnectorType: connectorpkg.TypeTelemetrySource,
			ConnectorKind: "generic_json_transform",
			DisplayName:   "Agency AVL poller",
			State:         connectorpkg.StateDryRunPassed,
			ConfigJSON:    config,
			SecretRefs:    []string{"AVL_HTTP_TOKEN_REF"},
			DryRunStatus:  "passed",
		}},
		jobs: []connectorpkg.DryRunJob{{
			ID:                  501,
			AgencyID:            "demo-agency",
			ConnectorInstanceID: 101,
			Status:              "passed",
			FinishedAt:          now,
			CreatedAt:           now,
			RedactionScanStatus: "passed",
			RedactedSummary:     json.RawMessage(`{"summary":"redacted fixture passed"}`),
		}},
	}
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1"}}}, connectorInstances: store}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/vehicle-avl", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("vehicle avl status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Activation Gate", "Mark ready for deployment-owned activation", "Device bindings exist", "Redaction scan passed"} {
		if !strings.Contains(body, want) {
			t.Fatalf("activation gate missing %q: %s", want, body)
		}
	}

	form := url.Values{"action": {"mark_vehicle_avl_ready"}, "connector_instance_id": {"101"}}
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/vehicle-avl", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("ready post status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.savedState.State != connectorpkg.StateReadyForActivation {
		t.Fatalf("saved state = %+v", store.savedState)
	}
	if !strings.Contains(rr.Body.String(), "ready for deployment-owned activation") {
		t.Fatalf("ready response missing readiness copy: %s", rr.Body.String())
	}
}

func TestVehicleAVLActivationGateBlocksWhenDeviceBindingMissing(t *testing.T) {
	config := json.RawMessage(`{"source_shape":"generic_json_transform","field_map":{"agency_id":"agency_id","device_id":"device.id","vehicle_id":"vehicle.id","observed_timestamp":"observed_at","lat":"position.lat","lon":"position.lon","quality":"quality"},"safety":{"target_path":"/v1/telemetry"}}`)
	store := &fakeVehicleAVLConnectorStore{
		instances: []connectorpkg.Instance{{
			ID:            101,
			AgencyID:      "demo-agency",
			ConnectorType: connectorpkg.TypeTelemetrySource,
			ConnectorKind: "generic_json_transform",
			DisplayName:   "Agency AVL poller",
			State:         connectorpkg.StateDryRunPassed,
			ConfigJSON:    config,
			SecretRefs:    []string{"AVL_HTTP_TOKEN_REF"},
			DryRunStatus:  "passed",
		}},
		jobs: []connectorpkg.DryRunJob{{ID: 501, ConnectorInstanceID: 101, Status: "passed", RedactionScanStatus: "passed", RedactedSummary: json.RawMessage(`{"summary":"redacted fixture passed"}`)}},
	}
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: store}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	form := url.Values{"action": {"mark_vehicle_avl_ready"}, "connector_instance_id": {"101"}}
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/vehicle-avl", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("blocked ready post status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.savedState.ID != 0 {
		t.Fatalf("state was updated despite missing device binding: %+v", store.savedState)
	}
	if !strings.Contains(rr.Body.String(), "activation readiness checks must pass") {
		t.Fatalf("response missing blocked activation message: %s", rr.Body.String())
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
	instances   []connectorpkg.Instance
	saved       connectorpkg.UpsertInstanceInput
	jobs        []connectorpkg.DryRunJob
	savedDryRun connectorpkg.CreateDryRunJobInput
	savedState  connectorpkg.UpdateInstanceStateInput
}

func (f *fakeVehicleAVLConnectorStore) ListInstances(context.Context, string) ([]connectorpkg.Instance, error) {
	return f.instances, nil
}

func (f *fakeVehicleAVLConnectorStore) ListDryRunJobs(context.Context, string, int) ([]connectorpkg.DryRunJob, error) {
	return f.jobs, nil
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

func (f *fakeVehicleAVLConnectorStore) CreateDryRunJob(_ context.Context, input connectorpkg.CreateDryRunJobInput) (connectorpkg.DryRunJob, error) {
	f.savedDryRun = input
	job := connectorpkg.DryRunJob{
		ID:                  501,
		AgencyID:            input.AgencyID,
		ConnectorInstanceID: input.ConnectorInstanceID,
		Status:              input.Status,
		StartedAt:           input.Now,
		FinishedAt:          input.Now,
		RedactedSummary:     input.RedactedSummary,
		AcceptedCount:       input.AcceptedCount,
		RejectedCount:       input.RejectedCount,
		DroppedCount:        input.DroppedCount,
		RedactionScanStatus: input.RedactionScanStatus,
		CreatedBy:           input.ActorID,
		CreatedAt:           input.Now,
	}
	f.jobs = append([]connectorpkg.DryRunJob{job}, f.jobs...)
	for i := range f.instances {
		if f.instances[i].ID == input.ConnectorInstanceID {
			f.instances[i].DryRunStatus = input.Status
			if input.Status == "passed" && input.RedactionScanStatus == "passed" {
				f.instances[i].State = connectorpkg.StateDryRunPassed
			} else {
				f.instances[i].State = connectorpkg.StateBlocked
			}
		}
	}
	return job, nil
}

func (f *fakeVehicleAVLConnectorStore) UpdateInstanceState(_ context.Context, input connectorpkg.UpdateInstanceStateInput) (connectorpkg.Instance, error) {
	f.savedState = input
	for i := range f.instances {
		if f.instances[i].ID == input.ID {
			f.instances[i].State = input.State
			return f.instances[i], nil
		}
	}
	return connectorpkg.Instance{}, nil
}
