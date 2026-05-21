package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"open-transit-rt/internal/auth"
)

func TestAdminPasswordLoginIssuesProductionSessionCookie(t *testing.T) {
	handler, store := newAdminLoginTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "https://admin.example.org/admin/login", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET login status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	state := extractLocalLoginState(t, rr.Body.String())

	form := url.Values{}
	form.Set("state", state)
	form.Set("email", "Admin@Example.ORG")
	form.Set("password", "correct horse battery staple")
	req = httptest.NewRequest(http.MethodPost, "https://admin.example.org/admin/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Forwarded-Proto", "https")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("POST login status = %d, want 303: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Location"); got != "/admin/operations" {
		t.Fatalf("Location = %q, want /admin/operations", got)
	}
	if store.email != "Admin@Example.ORG" || store.password != "correct horse battery staple" {
		t.Fatalf("store credentials = %q %q", store.email, store.password)
	}
	var session *http.Cookie
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "admin_session" {
			session = cookie
		}
	}
	if session == nil {
		t.Fatalf("login did not set admin_session: %#v", rr.Result().Cookies())
	}
	if !session.HttpOnly || !session.Secure || session.Path != "/admin" || session.SameSite != http.SameSiteLaxMode || session.MaxAge <= 0 {
		t.Fatalf("session cookie missing production attributes: %#v", session)
	}
	if strings.Contains(rr.Body.String(), session.Value) || strings.Contains(rr.Body.String(), "Bearer ") {
		t.Fatalf("login response leaked token text: %s", rr.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "https://admin.example.org/admin/operations", nil)
	req.AddCookie(session)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("cookie-authenticated operations status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
}

func TestAdminPasswordLoginFailureIsGenericAndStateIsSingleUse(t *testing.T) {
	handler, store := newAdminLoginTestHandler(t)
	store.err = auth.ErrInvalidCredentials
	req := httptest.NewRequest(http.MethodGet, "https://admin.example.org/admin/login", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	state := extractLocalLoginState(t, rr.Body.String())

	form := url.Values{"state": {state}, "email": {"missing@example.org"}, "password": {"wrong-password"}}
	req = httptest.NewRequest(http.MethodPost, "https://admin.example.org/admin/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("failed login status = %d, want 401", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, genericLoginFailureText) {
		t.Fatalf("failed login did not show generic text: %s", body)
	}
	if strings.Contains(body, "missing@example.org") || strings.Contains(body, "wrong-password") || strings.Contains(body, "no such user") {
		t.Fatalf("failed login leaked credential detail: %s", body)
	}

	req = httptest.NewRequest(http.MethodPost, "https://admin.example.org/admin/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("reused state status = %d, want 403", rr.Code)
	}
}

func TestFirstAdminSetupConsumesTokenSetsPasswordAndIssuesSession(t *testing.T) {
	handler, store := newAdminLoginTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "https://admin.example.org/admin/setup/first-admin?token=setup-token", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET setup status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	state := extractLocalLoginState(t, rr.Body.String())

	form := url.Values{}
	form.Set("state", state)
	form.Set("token", "setup-token")
	form.Set("password", "correct horse battery staple")
	form.Set("confirm_password", "correct horse battery staple")
	req = httptest.NewRequest(http.MethodPost, "https://admin.example.org/admin/setup/first-admin", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("POST setup status = %d, want 303: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Location"); got != "/admin/operations/setup-wizard" {
		t.Fatalf("Location = %q, want setup wizard", got)
	}
	if store.completedToken != "setup-token" || store.completedPassword != "correct horse battery staple" {
		t.Fatalf("completed setup = token %q password %q", store.completedToken, store.completedPassword)
	}
	if strings.Contains(rr.Body.String(), "setup-token") || strings.Contains(rr.Body.String(), "correct horse") {
		t.Fatalf("setup response leaked token or password: %s", rr.Body.String())
	}
}

func TestAdminLogoutRequiresExistingCookieCSRFAndExpiresSession(t *testing.T) {
	handler, _ := newAdminLoginTestHandler(t)
	signer, err := auth.NewSigner(auth.JWTConfig{Secrets: []string{"test-admin-secret"}, Issuer: "test-issuer", Audience: "test-audience", TTL: time.Hour})
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	token, _, err := signer.Sign("admin@example.org", "demo-agency", time.Hour)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "https://admin.example.org/admin/logout", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "admin_session", Value: token})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("logout without csrf status = %d, want 403: %s", rr.Code, rr.Body.String())
	}

	principal := auth.Principal{Subject: "admin@example.org", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodCookie}
	form := url.Values{"csrf_token": {auth.CSRFToken("test-csrf", principal)}}
	req = httptest.NewRequest(http.MethodPost, "https://admin.example.org/admin/logout", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "admin_session", Value: token})
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("logout status = %d, want 303: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Location"); got != "/admin/login" {
		t.Fatalf("Location = %q, want /admin/login", got)
	}
	var expired bool
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "admin_session" && cookie.MaxAge < 0 {
			expired = true
		}
	}
	if !expired {
		t.Fatalf("logout did not expire admin_session: %#v", rr.Result().Cookies())
	}
}

func newAdminLoginTestHandler(t *testing.T) (http.Handler, *fakePasswordStore) {
	t.Helper()
	t.Setenv("APP_ENV", "production")
	t.Setenv("AGENCY_ID", "demo-agency")
	t.Setenv("ADMIN_JWT_SECRET", "test-admin-secret")
	t.Setenv("ADMIN_JWT_ISSUER", "test-issuer")
	t.Setenv("ADMIN_JWT_AUDIENCE", "test-audience")
	t.Setenv("ADMIN_JWT_TTL", "1h")
	t.Setenv("CSRF_SECRET", "test-csrf")
	t.Setenv("PASSWORD_LOGIN_ENABLED", "true")
	cfg := auth.JWTConfig{Secrets: []string{"test-admin-secret"}, Issuer: "test-issuer", Audience: "test-audience", TTL: time.Hour}
	verifier, err := auth.NewVerifier(cfg)
	if err != nil {
		t.Fatalf("verifier: %v", err)
	}
	admin := auth.NewMiddleware(verifier, auth.StaticRoleStore{Roles: []auth.Role{auth.RoleAdmin, auth.RoleEditor, auth.RoleOperator, auth.RoleReadOnly}}, "test-csrf")
	passwords := &fakePasswordStore{principal: auth.Principal{Subject: "admin@example.org", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodCookie}}
	handler := newHandlerWithRealtimeAndPasswordStore(
		"demo-agency",
		fakeScheduleBuilder{},
		&fakePublicationStore{},
		fakeDeviceStore{},
		fakePinger{},
		admin,
		&fakeRealtimeArtifacts{},
		passwords,
	)
	return handler, passwords
}

type fakePasswordStore struct {
	principal         auth.Principal
	err               error
	email             string
	password          string
	completedToken    string
	completedPassword string
}

func (f *fakePasswordStore) AuthenticatePassword(_ context.Context, _ string, email string, password string, _ time.Time) (auth.Principal, error) {
	f.email = email
	f.password = password
	if f.err != nil {
		return auth.Principal{}, f.err
	}
	return f.principal, nil
}

func (f *fakePasswordStore) CompleteBootstrapPassword(_ context.Context, token string, password string, _ time.Time) (auth.Principal, error) {
	f.completedToken = token
	f.completedPassword = password
	if f.err != nil {
		return auth.Principal{}, f.err
	}
	if strings.TrimSpace(token) == "" {
		return auth.Principal{}, errors.New("missing token")
	}
	return f.principal, nil
}
