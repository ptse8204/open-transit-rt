package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"open-transit-rt/internal/auth"
	connectorpkg "open-transit-rt/internal/connectors"
)

func TestMonitoringConnectorRendersNoSendAdminForm(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/monitoring", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("monitoring connector status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Monitoring Setup", "Health Digest Preview", "no-send", "MONITORING_DIGEST_DESTINATION_REF", `name="action" value="save_monitoring_connector"`, "does not send notifications"} {
		if !strings.Contains(body, want) {
			t.Fatalf("monitoring setup missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"http://", "https://", "password=", "token=", "secret="} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("monitoring setup exposed unsafe copy %q: %s", forbidden, body)
		}
	}
}

func TestDiscoveryConnectorRendersReadinessAdminForm(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/discovery", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("discovery connector status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Discovery Setup", "/public/feeds.json metadata", "No portal automation", "PUBLIC_FEED_BASE_URL", `name="action" value="save_discovery_connector"`, "consumer status mutation"} {
		if !strings.Contains(body, want) {
			t.Fatalf("discovery setup missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"http://", "https://", "password=", "token=", "secret="} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("discovery setup exposed unsafe copy %q: %s", forbidden, body)
		}
	}
}

func TestMonitoringConnectorPostStoresNoSendMetadataOnly(t *testing.T) {
	store := &fakeVehicleAVLConnectorStore{}
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: store}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	form := monitoringConnectorValidForm()
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/monitoring", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("monitoring post status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.saved.ConnectorType != connectorpkg.TypeMonitoringExport || store.saved.ConnectorKind != "health_digest_no_send" {
		t.Fatalf("saved connector = %+v", store.saved)
	}
	var metadata monitoringConnectorConfigMetadata
	if err := json.Unmarshal(store.saved.ConfigJSON, &metadata); err != nil {
		t.Fatalf("decode monitoring metadata: %v", err)
	}
	if !metadata.NoSendDefault || !metadata.NotificationDeliveryOff || metadata.DestinationRef != "MONITORING_DIGEST_DESTINATION_REF" {
		t.Fatalf("metadata = %+v", metadata)
	}
	assertConnectorConfigNoRawRuntimeValues(t, store.saved.ConfigJSON)
	if !strings.Contains(rr.Body.String(), "No notification or export delivery was attempted") {
		t.Fatalf("response missing no-send notice: %s", rr.Body.String())
	}
}

func TestDiscoveryConnectorPostStoresNoAutomationMetadataOnly(t *testing.T) {
	store := &fakeVehicleAVLConnectorStore{}
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: store}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	form := discoveryConnectorValidForm()
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/discovery", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("discovery post status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.saved.ConnectorType != connectorpkg.TypeConsumerDiscovery || store.saved.ConnectorKind != "feeds_json_readiness" {
		t.Fatalf("saved connector = %+v", store.saved)
	}
	var metadata discoveryConnectorConfigMetadata
	if err := json.Unmarshal(store.saved.ConfigJSON, &metadata); err != nil {
		t.Fatalf("decode discovery metadata: %v", err)
	}
	if metadata.PortalAutomationEnabled || !metadata.ConsumerStatusMutationOff || metadata.PublicBaseURLEnvRef != "PUBLIC_FEED_BASE_URL" {
		t.Fatalf("metadata = %+v", metadata)
	}
	assertConnectorConfigNoRawRuntimeValues(t, store.saved.ConfigJSON)
	if !strings.Contains(rr.Body.String(), "No portal automation or consumer status mutation was attempted") {
		t.Fatalf("response missing no-automation notice: %s", rr.Body.String())
	}
}

func TestMonitoringAndDiscoveryRejectRawRuntimeRefs(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: &fakeVehicleAVLConnectorStore{}}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	cases := []struct {
		name string
		path string
		form url.Values
	}{
		{name: "monitoring url", path: "/admin/operations/connectors/monitoring", form: monitoringConnectorValidForm()},
		{name: "discovery url", path: "/admin/operations/connectors/discovery", form: discoveryConnectorValidForm()},
	}
	cases[0].form.Set("destination_ref", "https://hooks.example.test/private")
	cases[1].form.Set("public_base_url_env_ref", "https://feeds.example.test")
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.path, strings.NewReader(tc.form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want validation page: %s", rr.Code, rr.Body.String())
			}
			if strings.Contains(rr.Body.String(), "metadata was saved") {
				t.Fatalf("unsafe metadata was saved: %s", rr.Body.String())
			}
			if !strings.Contains(strings.ToLower(rr.Body.String()), "reference label") {
				t.Fatalf("response missing reference-label error: %s", rr.Body.String())
			}
		})
	}
}

func monitoringConnectorValidForm() url.Values {
	return url.Values{
		"action":          {"save_monitoring_connector"},
		"display_name":    {"Health digest preview"},
		"owner":           {"ops@example.org"},
		"mode":            {"health_digest_no_send"},
		"destination_ref": {"MONITORING_DIGEST_DESTINATION_REF"},
	}
}

func discoveryConnectorValidForm() url.Values {
	return url.Values{
		"action":                    {"save_discovery_connector"},
		"display_name":              {"Feed discovery readiness"},
		"owner":                     {"sharing@example.org"},
		"public_base_url_env_ref":   {"PUBLIC_FEED_BASE_URL"},
		"license_contact_owner_ref": {"PUBLIC_FEED_METADATA_OWNER"},
	}
}

func assertConnectorConfigNoRawRuntimeValues(t *testing.T, raw []byte) {
	t.Helper()
	lower := strings.ToLower(string(raw))
	for _, forbidden := range []string{"http://", "https://", "password", "token=", "secret=", "api_key", "webhook"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("config leaked %q: %s", forbidden, string(raw))
		}
	}
}
