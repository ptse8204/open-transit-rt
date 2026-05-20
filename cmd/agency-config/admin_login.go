package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"open-transit-rt/internal/appconfig"
	"open-transit-rt/internal/auth"
)

const (
	adminLoginPath          = "/admin/login"
	firstAdminSetupPath     = "/admin/setup/first-admin"
	adminLogoutPath         = "/admin/logout"
	adminLoginStateTTL      = 10 * time.Minute
	adminLoginPostMaxBytes  = 8 << 10
	adminSetupPostMaxBytes  = 12 << 10
	genericLoginFailureText = "Sign in failed. Check your email and password."
)

type adminLoginFlow struct {
	stateKey   string
	now        func() time.Time
	mu         sync.Mutex
	stateStore map[string]time.Time
}

func newAdminLoginFlow(csrfSecret string) *adminLoginFlow {
	return &adminLoginFlow{
		stateKey:   firstNonEmpty(csrfSecret, os.Getenv("CSRF_SECRET"), os.Getenv("ADMIN_JWT_SECRET")),
		now:        time.Now,
		stateStore: map[string]time.Time{},
	}
}

func (h *handler) adminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	if r.URL.Path != adminLoginPath {
		http.NotFound(w, r)
		return
	}
	if !h.passwordLoginAvailable() {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.renderAdminLogin(w, r, http.StatusOK, "")
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, adminLoginPostMaxBytes)
		h.adminLoginPost(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) firstAdminSetup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	if r.URL.Path != firstAdminSetupPath {
		http.NotFound(w, r)
		return
	}
	if !h.passwordLoginAvailable() {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		token := strings.TrimSpace(r.URL.Query().Get("token"))
		h.renderFirstAdminSetup(w, r, http.StatusOK, token, "")
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, adminSetupPostMaxBytes)
		h.firstAdminSetupPost(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) adminLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.URL.Path != adminLogoutPath {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	expireAdminSessionCookie(w, r)
	http.Redirect(w, r, adminLoginPath, http.StatusSeeOther)
}

func (h *handler) passwordLoginAvailable() bool {
	if h == nil || h.passwords == nil || h.loginFlow == nil || strings.TrimSpace(h.loginFlow.stateKey) == "" {
		return false
	}
	return !strings.EqualFold(strings.TrimSpace(os.Getenv("PASSWORD_LOGIN_ENABLED")), "false")
}

func (h *handler) adminLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderAdminLogin(w, r, http.StatusRequestEntityTooLarge, genericLoginFailureText)
		return
	}
	if !h.loginFlow.consumeState(strings.TrimSpace(r.FormValue("state"))) {
		h.renderAdminLogin(w, r, http.StatusForbidden, "This sign-in form expired. Reload and try again.")
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	principal, err := h.passwords.AuthenticatePassword(r.Context(), h.agencyID, email, password, time.Now().UTC())
	if err != nil {
		h.renderAdminLogin(w, r, http.StatusUnauthorized, genericLoginFailureText)
		return
	}
	if err := h.issueAdminSession(w, r, principal.Subject, principal.AgencyID); err != nil {
		http.Error(w, "sign-in is unavailable", http.StatusServiceUnavailable)
		return
	}
	http.Redirect(w, r, "/admin/operations", http.StatusSeeOther)
}

func (h *handler) firstAdminSetupPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderFirstAdminSetup(w, r, http.StatusRequestEntityTooLarge, "", "The setup form was too large. Reload and try again.")
		return
	}
	token := strings.TrimSpace(r.FormValue("token"))
	if !h.loginFlow.consumeState(strings.TrimSpace(r.FormValue("state"))) {
		h.renderFirstAdminSetup(w, r, http.StatusForbidden, token, "This setup form expired. Reload and try again.")
		return
	}
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")
	if password != confirm {
		h.renderFirstAdminSetup(w, r, http.StatusBadRequest, token, "The password confirmation did not match.")
		return
	}
	principal, err := h.passwords.CompleteBootstrapPassword(r.Context(), token, password, time.Now().UTC())
	if err != nil {
		h.renderFirstAdminSetup(w, r, http.StatusUnauthorized, "", "This setup link is invalid or expired. Generate a new first-admin link from the server console.")
		return
	}
	if err := h.issueAdminSession(w, r, principal.Subject, principal.AgencyID); err != nil {
		http.Error(w, "setup completed but sign-in is unavailable", http.StatusServiceUnavailable)
		return
	}
	http.Redirect(w, r, "/admin/operations/setup-wizard", http.StatusSeeOther)
}

func (h *handler) renderAdminLogin(w http.ResponseWriter, r *http.Request, status int, message string) {
	state, err := h.loginFlow.issueState()
	if err != nil {
		http.Error(w, "sign-in is unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = adminLoginTemplate.Execute(w, map[string]string{
		"AgencyID": h.agencyID,
		"State":    state,
		"Message":  message,
	})
}

func (h *handler) renderFirstAdminSetup(w http.ResponseWriter, r *http.Request, status int, token string, message string) {
	state, err := h.loginFlow.issueState()
	if err != nil {
		http.Error(w, "setup is unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = firstAdminSetupTemplate.Execute(w, map[string]string{
		"State":   state,
		"Token":   token,
		"Message": message,
	})
}

func (h *handler) issueAdminSession(w http.ResponseWriter, r *http.Request, subject string, agencyID string) error {
	cfg, err := auth.JWTConfigFromEnv()
	if err != nil {
		return err
	}
	signer, err := auth.NewSigner(cfg)
	if err != nil {
		return err
	}
	token, claims, err := signer.Sign(subject, agencyID, cfg.TTL)
	if err != nil {
		return err
	}
	expires := time.Unix(claims.Expires, 0).UTC()
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    token,
		Path:     "/admin",
		Expires:  expires,
		MaxAge:   int(time.Until(expires).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   adminCookieSecure(r),
	})
	return nil
}

func expireAdminSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/admin",
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   adminCookieSecure(r),
	})
}

func adminCookieSecure(r *http.Request) bool {
	return appconfig.IsProduction() || r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

func (f *adminLoginFlow) issueState() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate admin login state: %w", err)
	}
	expires := f.now().UTC().Add(adminLoginStateTTL)
	nonce := base64.RawURLEncoding.EncodeToString(raw[:])
	payload := nonce + "." + strconv.FormatInt(expires.Unix(), 10)
	signature := f.signState(payload)
	state := payload + "." + signature
	f.mu.Lock()
	f.cleanupLocked(f.now().UTC())
	f.stateStore[state] = expires
	f.mu.Unlock()
	return state, nil
}

func (f *adminLoginFlow) consumeState(state string) bool {
	if state == "" || !f.validStateSignature(state) {
		return false
	}
	now := f.now().UTC()
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cleanupLocked(now)
	expires, ok := f.stateStore[state]
	if !ok || now.After(expires) {
		delete(f.stateStore, state)
		return false
	}
	delete(f.stateStore, state)
	return true
}

func (f *adminLoginFlow) validStateSignature(state string) bool {
	parts := strings.Split(state, ".")
	if len(parts) != 3 {
		return false
	}
	payload := parts[0] + "." + parts[1]
	want := f.signState(payload)
	return hmac.Equal([]byte(want), []byte(parts[2]))
}

func (f *adminLoginFlow) signState(payload string) string {
	mac := hmac.New(sha256.New, []byte(f.stateKey))
	_, _ = mac.Write([]byte("admin-login-state-v1\x00" + payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (f *adminLoginFlow) cleanupLocked(now time.Time) {
	for state, expires := range f.stateStore {
		if now.After(expires) {
			delete(f.stateStore, state)
		}
	}
}

var adminLoginTemplate = template.Must(template.New("admin-login").Funcs(template.FuncMap{
	"operationsCSS": operationsCSS,
}).Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="robots" content="noindex,nofollow">
<link rel="icon" href="data:,">
<title>Admin Login</title>
<style>{{operationsCSS}}</style>
</head>
<body class="operations-body">
<main id="operations-main" tabindex="-1" aria-labelledby="admin-login-title">
<section class="login-panel" aria-label="Production admin login">
<p class="eyebrow">Open Transit RT</p>
<h1 id="admin-login-title">Admin Login</h1>
<p class="warning">Use deployment-owned credentials. Local demo sign-in is separate and remains disabled in production.</p>
{{if .Message}}<p class="warning" role="alert">{{.Message}}</p>{{end}}
<form method="post" action="/admin/login">
<input type="hidden" name="state" value="{{.State}}">
<input type="hidden" name="agency_id" value="{{.AgencyID}}">
<label for="admin-login-email">Email</label>
<input id="admin-login-email" name="email" type="email" autocomplete="username" required>
<label for="admin-login-password">Password</label>
<input id="admin-login-password" name="password" type="password" autocomplete="current-password" required>
<button type="submit">Sign in</button>
</form>
<p class="muted">SSO/OIDC is not configured in this roadmap. Future SSO must map identity to internal roles before issuing this app's session cookie.</p>
</section>
</main>
</body>
</html>`))

var firstAdminSetupTemplate = template.Must(template.New("first-admin-setup").Funcs(template.FuncMap{
	"operationsCSS": operationsCSS,
}).Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="robots" content="noindex,nofollow">
<link rel="icon" href="data:,">
<title>First Admin Setup</title>
<style>{{operationsCSS}}</style>
</head>
<body class="operations-body">
<main id="operations-main" tabindex="-1" aria-labelledby="first-admin-setup-title">
<section class="login-panel" aria-label="First admin setup">
<p class="eyebrow">Open Transit RT</p>
<h1 id="first-admin-setup-title">First Admin Setup</h1>
<p class="warning">Use the one-time server-generated link before it expires. The token is stored only as a hash.</p>
{{if .Message}}<p class="warning" role="alert">{{.Message}}</p>{{end}}
<form method="post" action="/admin/setup/first-admin">
<input type="hidden" name="state" value="{{.State}}">
<input type="hidden" name="token" value="{{.Token}}">
<label for="first-admin-password">Password</label>
<input id="first-admin-password" name="password" type="password" autocomplete="new-password" minlength="12" required>
<label for="first-admin-confirm">Confirm password</label>
<input id="first-admin-confirm" name="confirm_password" type="password" autocomplete="new-password" minlength="12" required>
<button type="submit">Create admin password</button>
</form>
<p class="muted">Generate a new setup link from the server console if this page reports an invalid or expired link.</p>
</section>
</main>
</body>
</html>`))

var _ passwordAuthStore = (*auth.PostgresAdminStore)(nil)
