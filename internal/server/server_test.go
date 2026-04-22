package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWrapAddsRequestIDAndRedactsSecretQueryValues(t *testing.T) {
	t.Setenv("METRICS_ENABLED", "")
	handler := wrap("test-service", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/admin/debug?token=secret&agency_id=demo-agency", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rr.Code)
	}
	if rr.Header().Get(requestIDHeader) == "" {
		t.Fatalf("missing request id header")
	}
	path := redactedPath(req)
	if strings.Contains(path, "secret") || !strings.Contains(path, "token=redacted") || !strings.Contains(path, "agency_id=present") {
		t.Fatalf("redacted path = %q", path)
	}
}

func TestMetricsEndpointOnlyWhenEnabled(t *testing.T) {
	t.Setenv("METRICS_ENABLED", "")
	handler := wrap("test-service", http.NotFoundHandler())
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("metrics disabled status = %d, want 404", rr.Code)
	}

	t.Setenv("METRICS_ENABLED", "true")
	handler = wrap("test-service", http.NotFoundHandler())
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("metrics enabled status = %d, want 200", rr.Code)
	}
}
