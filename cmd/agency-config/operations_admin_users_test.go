package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"open-transit-rt/internal/auth"
)

func TestAdminUsersPageRequiresAdminAndScopesAgency(t *testing.T) {
	store := newFakeAdminUserPasswordStore()
	readOnly := auth.Principal{Subject: "viewer@example.org", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer}
	handler := newAdminUsersTestHandler(t, readOnly, store)
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/admin/users", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("read-only admin users status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
	if store.listAgencyID != "" {
		t.Fatalf("read-only request loaded users for agency %q", store.listAgencyID)
	}

	admin := auth.Principal{Subject: "admin@example.org", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer}
	handler = newAdminUsersTestHandler(t, admin, store)
	req = httptest.NewRequest(http.MethodGet, "/admin/operations/admin/users", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("admin users status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Users &amp; Roles", "Private admin-only user management", "SSO/OIDC is not implemented", "Password setup and reset links are short lived"} {
		if !strings.Contains(body, want) {
			t.Fatalf("admin users body missing %q: %s", want, body)
		}
	}
	if store.listAgencyID != "demo-agency" {
		t.Fatalf("list agency = %q, want demo-agency", store.listAgencyID)
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/admin/users?agency_id=other-agency", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("cross-agency users status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
}

func TestAdminUsersCreateDisableAndResetAreScoped(t *testing.T) {
	store := newFakeAdminUserPasswordStore()
	admin := auth.Principal{Subject: "admin@example.org", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer}
	handler := newAdminUsersTestHandler(t, admin, store)

	form := url.Values{}
	form.Set("csrf_token", auth.CSRFToken("test-csrf", admin))
	form.Set("action", "create_user")
	form.Set("email", "Ops@Example.ORG")
	form.Set("display_name", "Operations Lead")
	form.Add("role", string(auth.RoleOperator))
	form.Add("role", string(auth.RoleReadOnly))
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/admin/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("create user status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.createInput.AgencyID != "demo-agency" || store.createInput.Email != "Ops@Example.ORG" || store.createInput.ActorID != "admin@example.org" {
		t.Fatalf("unexpected create input: %+v", store.createInput)
	}
	if got := authRoleStrings(store.createInput.Roles); strings.Join(got, ",") != "operator,read_only" {
		t.Fatalf("create roles = %v, want operator/read_only", got)
	}

	form = url.Values{}
	form.Set("csrf_token", auth.CSRFToken("test-csrf", admin))
	form.Set("agency_id", "other-agency")
	form.Set("action", "create_user")
	form.Set("email", "bad@example.org")
	form.Add("role", string(auth.RoleAdmin))
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/admin/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("cross-agency create status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
	if store.createInput.Email == "bad@example.org" {
		t.Fatalf("cross-agency create reached store: %+v", store.createInput)
	}

	form = url.Values{}
	form.Set("csrf_token", auth.CSRFToken("test-csrf", admin))
	form.Set("action", "disable_user")
	form.Set("user_id", "42")
	form.Set("reason", "staff left")
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/admin/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("disable user status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if store.disableInput.AgencyID != "demo-agency" || store.disableInput.UserID != 42 || store.disableInput.Reason != "staff left" {
		t.Fatalf("unexpected disable input: %+v", store.disableInput)
	}

	form = url.Values{}
	form.Set("csrf_token", auth.CSRFToken("test-csrf", admin))
	form.Set("action", "password_reset")
	form.Set("user_id", "42")
	req = httptest.NewRequest(http.MethodPost, "https://admin.example.org/admin/operations/admin/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("password reset status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if !strings.Contains(body, "/admin/setup/first-admin?token=") || !strings.Contains(body, "shown once") {
		t.Fatalf("password reset body missing one-time setup link: %s", body)
	}
	token := extractResetToken(t, body)
	if store.resetHash == "" || store.resetHash == token || strings.Contains(body, store.resetHash) {
		t.Fatalf("reset hash/token handling invalid: hash=%q token=%q body=%s", store.resetHash, token, body)
	}
	if store.resetInput.AgencyID != "demo-agency" || store.resetInput.UserID != 42 || store.resetInput.ActorID != "admin@example.org" {
		t.Fatalf("unexpected reset input: %+v", store.resetInput)
	}
}

func TestAdminUsersReadOnlyCannotMutate(t *testing.T) {
	store := newFakeAdminUserPasswordStore()
	readOnly := auth.Principal{Subject: "viewer@example.org", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleReadOnly}, Method: auth.MethodBearer}
	handler := newAdminUsersTestHandler(t, readOnly, store)
	form := url.Values{"action": {"create_user"}, "email": {"ops@example.org"}, "role": {string(auth.RoleAdmin)}}
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/admin/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("read-only create status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
	if store.createInput.Email != "" || store.disableInput.UserID != 0 || store.resetHash != "" {
		t.Fatalf("read-only mutation reached store: create=%+v disable=%+v resetHash=%q", store.createInput, store.disableInput, store.resetHash)
	}
}

func newAdminUsersTestHandler(t *testing.T, principal auth.Principal, store *fakeAdminUserPasswordStore) http.Handler {
	t.Helper()
	t.Setenv("AGENCY_ID", "demo-agency")
	t.Setenv("CSRF_SECRET", "test-csrf")
	return newHandlerWithRealtimeAndPasswordStore(
		"demo-agency",
		fakeScheduleBuilder{},
		&fakePublicationStore{},
		fakeDeviceStore{},
		fakePinger{},
		auth.TestAuthenticator{Principal: principal},
		&fakeRealtimeArtifacts{},
		store,
	)
}

type fakeAdminUserPasswordStore struct {
	fakePasswordStore
	listUsers    []auth.AdminUser
	listAgencyID string
	createInput  auth.AdminUserCreateInput
	disableInput auth.AdminUserDisableInput
	resetInput   auth.PasswordResetTokenInput
	resetHash    string
}

func newFakeAdminUserPasswordStore() *fakeAdminUserPasswordStore {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	return &fakeAdminUserPasswordStore{
		fakePasswordStore: fakePasswordStore{principal: auth.Principal{Subject: "admin@example.org", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodCookie}},
		listUsers: []auth.AdminUser{{
			ID:             42,
			AgencyID:       "demo-agency",
			Email:          "ops@example.org",
			DisplayName:    "Operations Lead",
			Subject:        "ops@example.org",
			Roles:          []auth.Role{auth.RoleOperator, auth.RoleReadOnly},
			Status:         "active",
			PasswordStatus: "not_configured",
			CreatedAt:      now,
		}},
	}
}

func (f *fakeAdminUserPasswordStore) ListAdminUsers(_ context.Context, agencyID string) ([]auth.AdminUser, error) {
	f.listAgencyID = agencyID
	return append([]auth.AdminUser(nil), f.listUsers...), nil
}

func (f *fakeAdminUserPasswordStore) CreateAdminUser(_ context.Context, input auth.AdminUserCreateInput) (auth.AdminUser, error) {
	f.createInput = input
	return auth.AdminUser{ID: 99, AgencyID: input.AgencyID, Email: input.Email, Subject: input.Email, Roles: input.Roles, Status: "active", PasswordStatus: "not_configured", CreatedAt: input.Now}, nil
}

func (f *fakeAdminUserPasswordStore) DisableAdminUser(_ context.Context, input auth.AdminUserDisableInput) error {
	f.disableInput = input
	return nil
}

func (f *fakeAdminUserPasswordStore) CreatePasswordResetToken(_ context.Context, input auth.PasswordResetTokenInput, tokenHash string) (auth.BootstrapResult, error) {
	f.resetInput = input
	f.resetHash = tokenHash
	expires := input.Now.Add(input.TTL)
	if expires.IsZero() {
		expires = time.Now().UTC().Add(30 * time.Minute)
	}
	return auth.BootstrapResult{AgencyID: input.AgencyID, Email: "ops@example.org", Subject: "ops@example.org", UserID: input.UserID, TokenID: 7, ExpiresAt: expires}, nil
}

func authRoleStrings(roles []auth.Role) []string {
	out := make([]string, 0, len(roles))
	for _, role := range roles {
		out = append(out, string(role))
	}
	return out
}

func extractResetToken(t *testing.T, body string) string {
	t.Helper()
	matches := regexp.MustCompile(`token=([A-Za-z0-9_-]+)`).FindStringSubmatch(body)
	if len(matches) != 2 {
		t.Fatalf("reset token not found in body: %s", body)
	}
	return matches[1]
}
