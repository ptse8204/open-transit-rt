package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"open-transit-rt/internal/appconfig"
	"open-transit-rt/internal/auth"
)

type operationsSessionBannerView struct {
	Subject string   `json:"subject"`
	Agency  string   `json:"agency"`
	Roles   []string `json:"roles"`
	Method  string   `json:"method"`
}

type operationsAuthStatusView struct {
	AgencyID           string                      `json:"agency_id"`
	GeneratedAt        time.Time                   `json:"generated_at"`
	Boundary           string                      `json:"boundary"`
	Session            operationsSessionBannerView `json:"session"`
	ActiveAuthMode     string                      `json:"active_auth_mode"`
	LocalDemoLogin     string                      `json:"local_demo_login"`
	PasswordLogin      string                      `json:"password_login"`
	SSOStatus          string                      `json:"sso_status"`
	SessionTTL         string                      `json:"session_ttl"`
	CookiePolicy       string                      `json:"cookie_policy"`
	BearerSupport      string                      `json:"bearer_support"`
	CSRFPolicy         string                      `json:"csrf_policy"`
	PasswordReset      string                      `json:"password_reset"`
	FutureSSODirection string                      `json:"future_sso_direction"`
	NextAction         string                      `json:"next_action"`
	DoesNotProve       string                      `json:"does_not_prove"`
}

func (h *handler) renderAuthStatus(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	page := h.buildOperationsPage(r, principal, "admin-sessions")
	renderOperationsTemplate(w, "admin-sessions", page)
}

func (h *handler) renderAuthStatusJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	page := h.buildOperationsPage(r, principal, "admin-sessions")
	writeJSON(w, http.StatusOK, page.AuthStatus)
}

func (h *handler) buildOperationsSessionBanner(principal auth.Principal) operationsSessionBannerView {
	return operationsSessionBannerView{
		Subject: firstNonEmpty(strings.TrimSpace(principal.Subject), "unknown subject"),
		Agency:  firstNonEmpty(strings.TrimSpace(principal.AgencyID), "unknown agency"),
		Roles:   safePrincipalRoles(principal.Roles),
		Method:  firstNonEmpty(strings.TrimSpace(principal.Method), "unknown"),
	}
}

func (h *handler) buildOperationsAuthStatus(r *http.Request, principal auth.Principal, generatedAt time.Time) operationsAuthStatusView {
	return operationsAuthStatusView{
		AgencyID:           principal.AgencyID,
		GeneratedAt:        generatedAt,
		Boundary:           "Private current-session status only. This page does not expose tokens, cookies, password hashes, reset tokens, provider secrets, or raw credential data.",
		Session:            h.buildOperationsSessionBanner(principal),
		ActiveAuthMode:     activeAuthMode(principal),
		LocalDemoLogin:     h.localDemoLoginStatus(r),
		PasswordLogin:      h.passwordLoginStatus(),
		SSOStatus:          "not configured / future; no OIDC endpoints, redirects, discovery, callback, JWKS, claim mapping, or provider login exists in this roadmap",
		SessionTTL:         adminSessionTTLStatus(),
		CookiePolicy:       "admin_session uses HttpOnly, SameSite=Lax, path=/admin, and Secure in production or HTTPS contexts.",
		BearerSupport:      "Bearer admin JWT remains supported for operator and API checks; roles are still loaded from the database.",
		CSRFPolicy:         "Unsafe cookie-authenticated admin requests require the generated CSRF token. Missing or stale form tokens are rejected before handlers mutate state.",
		PasswordReset:      "Admins generate short-lived one-time setup links from Users & Roles. Tokens are stored only as hashes and shown once.",
		FutureSSODirection: "Future SSO/OIDC must verify the external identity, map it to an internal subject, agency, and roles, then issue the same internal signed admin_session cookie.",
		NextAction:         "Use Users & Roles for password reset links and role changes. Keep SSO marked future until a provider adapter, callback, validation, and mapping flow exists.",
		DoesNotProve:       "This status page does not prove SSO support, compliance, production readiness, consumer acceptance, hosted service availability, uptime, or credential absence.",
	}
}

func activeAuthMode(principal auth.Principal) string {
	switch principal.Method {
	case auth.MethodCookie:
		return "browser admin_session cookie"
	case auth.MethodBearer:
		return "Bearer admin JWT"
	default:
		return firstNonEmpty(strings.TrimSpace(principal.Method), "unknown")
	}
}

func (h *handler) localDemoLoginStatus(r *http.Request) string {
	if h == nil || h.localLogin == nil || !h.localLogin.enabled {
		return "disabled"
	}
	if h.localLogin.available(r) {
		return "available for this local/demo loopback request"
	}
	if appconfig.IsProduction() {
		return "configured but disabled by production mode"
	}
	return "configured but unavailable for this host or missing local demo settings"
}

func (h *handler) passwordLoginStatus() string {
	if h == nil || h.passwords == nil {
		return "unavailable: password credential store is not configured"
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv("PASSWORD_LOGIN_ENABLED")), "false") {
		return "disabled by PASSWORD_LOGIN_ENABLED=false"
	}
	if h.loginFlow == nil || strings.TrimSpace(h.loginFlow.stateKey) == "" {
		return "unavailable: login state signing key is missing"
	}
	return "enabled"
}

func adminSessionTTLStatus() string {
	raw := strings.TrimSpace(os.Getenv("ADMIN_JWT_TTL"))
	if raw == "" {
		return auth.DefaultAdminTokenTTL.String() + " default"
	}
	ttl, err := time.ParseDuration(raw)
	if err != nil || ttl <= 0 {
		return "invalid ADMIN_JWT_TTL; startup signer/verifier configuration should reject it"
	}
	return ttl.String()
}
