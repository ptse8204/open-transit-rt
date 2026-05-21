package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"open-transit-rt/internal/auth"
)

func TestOperationsSessionBannerAndLoginSessionsPage(t *testing.T) {
	principal := auth.Principal{Subject: "reader@example.org", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodCookie}
	handler := newAdminUsersTestHandler(t, principal, newFakeAdminUserPasswordStore())
	t.Setenv("PASSWORD_LOGIN_ENABLED", "true")
	t.Setenv("ADMIN_JWT_TTL", "2h")

	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("dashboard status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Signed in as", "reader@example.org", "Roles read_only", `href="/admin/operations/admin/sessions#logout"`, "Auth cookie"} {
		if !strings.Contains(body, want) {
			t.Fatalf("dashboard body missing session marker %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"Bearer ", "admin_session=", "PASSWORD=", "SECRET="} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("dashboard session banner leaked forbidden text %q: %s", forbidden, body)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/admin/sessions", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("login sessions status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body = rr.Body.String()
	for _, want := range []string{"Login &amp; Sessions", "Password login", "enabled", "SSO/OIDC", "not configured / future", "2h0m0s", "Future SSO direction", `action="/admin/logout"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("login sessions body missing %q: %s", want, body)
		}
	}
	if strings.Contains(body, "callback") && strings.Contains(body, "implemented") && !strings.Contains(body, "no OIDC endpoints") {
		t.Fatalf("login sessions body may imply SSO endpoint support: %s", body)
	}
}

func TestOperationsLoginSessionsJSON(t *testing.T) {
	principal := auth.Principal{Subject: "operator@example.org", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleOperator}, Method: auth.MethodBearer}
	handler := newAdminUsersTestHandler(t, principal, newFakeAdminUserPasswordStore())
	t.Setenv("PASSWORD_LOGIN_ENABLED", "false")

	req := httptest.NewRequest(http.MethodGet, "/admin/operations/admin/sessions.json", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("login sessions JSON status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var view operationsAuthStatusView
	if err := json.Unmarshal(rr.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode login sessions JSON: %v: %s", err, rr.Body.String())
	}
	if view.Session.Subject != "operator@example.org" || view.Session.Method != auth.MethodBearer || view.ActiveAuthMode != "Bearer admin JWT" {
		t.Fatalf("unexpected session view: %+v", view.Session)
	}
	if !strings.Contains(view.PasswordLogin, "disabled") || !strings.Contains(view.SSOStatus, "not configured / future") {
		t.Fatalf("unexpected auth statuses: password=%q sso=%q", view.PasswordLogin, view.SSOStatus)
	}
	if strings.Contains(rr.Body.String(), "admin_session=") || strings.Contains(rr.Body.String(), "token=") {
		t.Fatalf("login sessions JSON leaked cookie/token text: %s", rr.Body.String())
	}
}
