package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"open-transit-rt/internal/auth"
)

func TestOperationsFocusedConfigPagesArePrivateAndNarrow(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02ReadOnlyPrincipal()})
	for _, path := range []string{
		"/admin/operations/config",
		"/admin/operations/config/agency",
		"/admin/operations/config/feeds",
		"/admin/operations/config/auth",
		"/admin/operations/config/deployment",
		"/admin/operations/config/advanced",
	} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
			}
			if got := rr.Header().Get("Cache-Control"); got != "no-store" {
				t.Fatalf("%s Cache-Control = %q, want no-store", path, got)
			}
			body := rr.Body.String()
			for _, want := range []string{"Focused", "Current Settings", "Back to Dashboard"} {
				if !strings.Contains(body, want) {
					t.Fatalf("%s missing focused config marker %q: %s", path, want, body)
				}
			}
			if strings.Contains(body, `<form`) {
				t.Fatalf("%s rendered a mutation form for read-only user: %s", path, body)
			}
			req = httptest.NewRequest(http.MethodGet, "/public"+strings.TrimPrefix(path, "/admin"), nil)
			rr = httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusNotFound {
				t.Fatalf("public mirror %s status = %d, want 404", path, rr.Code)
			}
		})
	}
}

func TestOperationsFocusedFeedConfigOwnsPublicationMetadataForm(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/config/feeds", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("config feeds status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		`<form method="post" action="/admin/operations/setup#publication-metadata">`,
		`name="action" value="publication_bootstrap"`,
		`id="config_public_base_url" type="url" name="public_base_url"`,
		`id="config_feed_base_url" type="url" name="feed_base_url"`,
		`id="config_publication_environment" name="publication_environment"`,
		`<button type="submit">Store publication metadata</button>`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("focused feed config missing %q: %s", want, body)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/setup", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("setup status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body = rr.Body.String()
	for _, stale := range []string{"setup_public_base_url", "setup_feed_base_url", `name="action" value="publication_bootstrap"`} {
		if strings.Contains(body, stale) {
			t.Fatalf("setup diagnostics still owns publication form marker %q: %s", stale, body)
		}
	}
	if !strings.Contains(body, "/admin/operations/config/feeds") {
		t.Fatalf("setup diagnostics missing focused config link: %s", body)
	}
}
