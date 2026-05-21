package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
	connectorpkg "open-transit-rt/internal/connectors"
)

func TestValidatorConnectorRendersAdminMetadataForm(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/validators", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("validator connector status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Validator Setup",
		`<form method="post" action="/admin/operations/connectors/validators#validator-connector-form">`,
		`name="action" value="save_validator_connector"`,
		compliance.ValidationHealthStaticValidatorID,
		compliance.ValidationHealthRealtimeValidatorID,
		"GTFS_VALIDATOR_PATH",
		"GTFS_RT_VALIDATOR_PATH",
		"Raw commands",
		"Validation Health",
		"Open Schedule Quality",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("validator setup page missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"java -jar", "/tmp/", "http://", "https://", "password=", "token=", "secret="} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("validator setup page exposed unsafe copy %q: %s", forbidden, body)
		}
	}
}

func TestValidatorConnectorReadOnlyHasNoMutationForm(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02ReadOnlyPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/validators", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("validator connector status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if strings.Contains(body, `<form`) {
		t.Fatalf("read-only validator setup rendered mutation form: %s", body)
	}
	if !strings.Contains(body, "changes require an admin role") {
		t.Fatalf("read-only validator setup missing role boundary: %s", body)
	}
}

func TestValidatorConnectorPostStoresAllowlistedMetadataOnly(t *testing.T) {
	store := &fakeVehicleAVLConnectorStore{}
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: store}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	form := validatorConnectorValidForm()
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/validators", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("validator connector post status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.saved.ConnectorType != connectorpkg.TypeValidator || store.saved.ConnectorKind != "static_mobilitydata" {
		t.Fatalf("saved connector = %+v", store.saved)
	}
	if store.saved.State != connectorpkg.StateConfiguredNotTested {
		t.Fatalf("state = %q", store.saved.State)
	}
	if len(store.saved.SecretRefs) != 0 {
		t.Fatalf("validator connector should not store secret refs: %+v", store.saved.SecretRefs)
	}
	var metadata validatorConnectorConfigMetadata
	if err := json.Unmarshal(store.saved.ConfigJSON, &metadata); err != nil {
		t.Fatalf("decode saved validator metadata: %v", err)
	}
	if metadata.ValidatorID != compliance.ValidationHealthStaticValidatorID || metadata.ToolingPathEnvRef != "GTFS_VALIDATOR_PATH" {
		t.Fatalf("metadata = %+v", metadata)
	}
	if !metadata.ServerOwnedRunOnly || !metadata.RawCommandsBlocked {
		t.Fatalf("server-owned/raw-command flags missing: %+v", metadata)
	}
	raw := strings.ToLower(string(store.saved.ConfigJSON))
	for _, forbidden := range []string{"http://", "https://", "java -jar", "/tmp/", "password", "token=", "secret=", "argv"} {
		if strings.Contains(raw, forbidden) {
			t.Fatalf("saved validator config leaked %q: %s", forbidden, raw)
		}
	}
	if !strings.Contains(rr.Body.String(), "metadata was saved") {
		t.Fatalf("response missing bounded saved notice: %s", rr.Body.String())
	}
}

func TestValidatorConnectorPostRejectsRawCommandsAndUnsupportedIDs(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: &fakeVehicleAVLConnectorStore{}}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	for name, mutate := range map[string]func(url.Values){
		"unsupported id": func(v url.Values) { v.Set("validator_id", "custom-validator") },
		"raw command":    func(v url.Values) { v.Set("tooling_path_env_ref", "java -jar /tmp/validator.jar") },
		"raw path":       func(v url.Values) { v.Set("version_env_ref", "/private/version.txt") },
	} {
		t.Run(name, func(t *testing.T) {
			form := validatorConnectorValidForm()
			mutate(form)
			req := httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/validators", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want validation page: %s", rr.Code, rr.Body.String())
			}
			if strings.Contains(rr.Body.String(), "metadata was saved") {
				t.Fatalf("unsafe validator connector metadata was saved: %s", rr.Body.String())
			}
			body := strings.ToLower(rr.Body.String())
			if !strings.Contains(body, "allowlisted") && !strings.Contains(body, "reference label") {
				t.Fatalf("response missing validation error: %s", rr.Body.String())
			}
		})
	}
}

func validatorConnectorValidForm() url.Values {
	return url.Values{
		"action":               {"save_validator_connector"},
		"display_name":         {"Static validator"},
		"owner":                {"validation@example.org"},
		"validator_id":         {compliance.ValidationHealthStaticValidatorID},
		"tooling_path_env_ref": {"GTFS_VALIDATOR_PATH"},
		"version_env_ref":      {"GTFS_VALIDATOR_VERSION"},
		"timeout_seconds":      {"120"},
	}
}
