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
	"open-transit-rt/internal/prediction"
)

func TestPredictionConnectorRendersAdminMetadataForm(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/prediction", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("prediction connector status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Prediction Setup",
		`<form method="post" action="/admin/operations/connectors/prediction#prediction-connector-form">`,
		`name="action" value="save_prediction_connector"`,
		`value="deterministic_default"`,
		`value="external_http_shadow"`,
		`value="external_http_fail_closed"`,
		`name="endpoint_url_env_ref"`,
		`name="allowed_hosts_env_ref"`,
		`value="` + prediction.ExternalHTTPTripUpdatesPath + `"`,
		"Vehicle Positions publishing stays independent",
		"External HTTP prediction is not enabled by this page",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("prediction setup page missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"http://", "https://", "password=", "token=", "secret=", "production-grade ETA quality."} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("prediction setup page exposed unsafe copy %q: %s", forbidden, body)
		}
	}
}

func TestPredictionConnectorReadOnlyHasNoMutationForm(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02ReadOnlyPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/connectors/prediction", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("prediction connector status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if strings.Contains(body, `<form`) {
		t.Fatalf("read-only prediction setup rendered mutation form: %s", body)
	}
	if !strings.Contains(body, "changes require an admin role") {
		t.Fatalf("read-only prediction setup missing role boundary: %s", body)
	}
}

func TestPredictionConnectorPostStoresReferenceMetadataOnly(t *testing.T) {
	store := &fakeVehicleAVLConnectorStore{}
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: store}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	form := predictionConnectorValidForm()
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/prediction", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("prediction connector post status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.saved.ConnectorType != connectorpkg.TypePrediction || store.saved.ConnectorKind != "external_http_shadow" {
		t.Fatalf("saved connector = %+v", store.saved)
	}
	if store.saved.State != connectorpkg.StateConfiguredNotTested {
		t.Fatalf("state = %q", store.saved.State)
	}
	if got := strings.Join(store.saved.SecretRefs, ","); got != "PREDICTOR_TOKEN" {
		t.Fatalf("secret refs = %q", got)
	}
	var metadata predictionConnectorConfigMetadata
	if err := json.Unmarshal(store.saved.ConfigJSON, &metadata); err != nil {
		t.Fatalf("decode saved prediction metadata: %v", err)
	}
	if metadata.Mode != "external_http_shadow" || !metadata.ShadowMode || metadata.Path != prediction.ExternalHTTPTripUpdatesPath {
		t.Fatalf("metadata = %+v", metadata)
	}
	if !metadata.VehiclePositionsIndependent {
		t.Fatalf("vehicle positions independence not recorded: %+v", metadata)
	}
	raw := strings.ToLower(string(store.saved.ConfigJSON))
	for _, forbidden := range []string{"http://", "https://", "password", "token", "secret", "api_key"} {
		if strings.Contains(raw, forbidden) {
			t.Fatalf("saved config metadata leaked %q: %s", forbidden, raw)
		}
	}
	if !strings.Contains(rr.Body.String(), "metadata was saved") {
		t.Fatalf("response missing saved notice: %s", rr.Body.String())
	}
}

func TestPredictionConnectorPostRejectsUnsafeEndpointRefsAndPaths(t *testing.T) {
	handler := newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}, connectorInstances: &fakeVehicleAVLConnectorStore{}}, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	for name, mutate := range map[string]func(url.Values){
		"raw endpoint": func(v url.Values) { v.Set("endpoint_url_env_ref", "https://predictor.example/v1/predict/trip-updates") },
		"raw token":    func(v url.Values) { v.Set("token_ref", "token=secret") },
		"wrong path":   func(v url.Values) { v.Set("path", "/predict") },
	} {
		t.Run(name, func(t *testing.T) {
			form := predictionConnectorValidForm()
			mutate(form)
			req := httptest.NewRequest(http.MethodPost, "/admin/operations/connectors/prediction", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want validation page: %s", rr.Code, rr.Body.String())
			}
			if strings.Contains(rr.Body.String(), "metadata was saved") {
				t.Fatalf("unsafe prediction connector metadata was saved: %s", rr.Body.String())
			}
			body := strings.ToLower(rr.Body.String())
			if !strings.Contains(body, "must") && !strings.Contains(body, "path") && !strings.Contains(body, "ref") {
				t.Fatalf("response missing validation error: %s", rr.Body.String())
			}
		})
	}
}

func predictionConnectorValidForm() url.Values {
	return url.Values{
		"action":                {"save_prediction_connector"},
		"display_name":          {"Agency prediction shadow"},
		"owner":                 {"integrations@example.org"},
		"mode":                  {"external_http_shadow"},
		"endpoint_url_env_ref":  {"TRIP_UPDATES_EXTERNAL_HTTP_URL"},
		"allowed_hosts_env_ref": {"TRIP_UPDATES_EXTERNAL_HTTP_ALLOWED_HOSTS"},
		"path":                  {prediction.ExternalHTTPTripUpdatesPath},
		"token_ref":             {"PREDICTOR_TOKEN"},
		"timeout_seconds":       {"2"},
	}
}
