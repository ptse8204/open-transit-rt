package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"net"
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
	localAdminLoginPath         = "/admin/local-login"
	localAdminSessionDefaultTTL = 2 * time.Hour
	localAdminStateTTL          = 10 * time.Minute
	localAdminLoginPostMaxBytes = 4 << 10
)

type localAdminLogin struct {
	enabled    bool
	subject    string
	agencyID   string
	ttl        time.Duration
	stateKey   string
	now        func() time.Time
	mu         sync.Mutex
	stateStore map[string]time.Time
}

func newLocalAdminLoginFromEnv(defaultAgencyID string, csrfSecret string) *localAdminLogin {
	return &localAdminLogin{
		enabled: strings.EqualFold(strings.TrimSpace(os.Getenv("LOCAL_ADMIN_LOGIN_ENABLED")), "true"),
		subject: firstNonEmpty(
			os.Getenv("LOCAL_ADMIN_SUBJECT"),
			os.Getenv("ADMIN_SUBJECT"),
			"admin@example.com",
		),
		agencyID: firstNonEmpty(
			os.Getenv("LOCAL_ADMIN_AGENCY_ID"),
			os.Getenv("AGENCY_ID"),
			defaultAgencyID,
		),
		ttl:        localAdminSessionTTLFromEnv(),
		stateKey:   firstNonEmpty(csrfSecret, os.Getenv("CSRF_SECRET"), os.Getenv("ADMIN_JWT_SECRET")),
		now:        time.Now,
		stateStore: map[string]time.Time{},
	}
}

func localAdminSessionTTLFromEnv() time.Duration {
	raw := strings.TrimSpace(os.Getenv("LOCAL_ADMIN_SESSION_TTL"))
	if raw == "" {
		return localAdminSessionDefaultTTL
	}
	ttl, err := time.ParseDuration(raw)
	if err != nil || ttl <= 0 {
		return localAdminSessionDefaultTTL
	}
	if ttl > localAdminSessionDefaultTTL {
		return localAdminSessionDefaultTTL
	}
	return ttl
}

func (h *handler) localAdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	if r.URL.Path != localAdminLoginPath {
		http.NotFound(w, r)
		return
	}
	if h.localLogin == nil || !h.localLogin.available(r) {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.renderLocalAdminLogin(w, r, "")
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, localAdminLoginPostMaxBytes)
		h.localAdminLoginPost(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (l *localAdminLogin) available(r *http.Request) bool {
	return l.enabled &&
		!appconfig.IsProduction() &&
		strings.TrimSpace(l.subject) != "" &&
		strings.TrimSpace(l.agencyID) != "" &&
		strings.TrimSpace(l.stateKey) != "" &&
		localAdminHost(r.Host)
}

func localAdminHost(hostport string) bool {
	host := strings.TrimSpace(hostport)
	if host == "" {
		return false
	}
	if parsedHost, _, err := net.SplitHostPort(host); err == nil {
		host = parsedHost
	}
	host = strings.Trim(strings.ToLower(host), "[]")
	switch host {
	case "localhost", "127.0.0.1", "::1":
		return true
	default:
		return false
	}
}

func (h *handler) renderLocalAdminLogin(w http.ResponseWriter, r *http.Request, message string) {
	state, err := h.localLogin.issueState()
	if err != nil {
		http.Error(w, "local sign-in is unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = localAdminLoginTemplate.Execute(w, map[string]string{
		"AgencyID": h.localLogin.agencyID,
		"Subject":  h.localLogin.subject,
		"State":    state,
		"TTL":      friendlyDuration(h.localLogin.ttl),
		"Message":  message,
	})
}

func (h *handler) localAdminLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderLocalAdminLoginWithStatus(w, r, http.StatusRequestEntityTooLarge, "The local sign-in form was too large. Reload this page and try again.")
		return
	}
	if !h.localLogin.consumeState(strings.TrimSpace(r.FormValue("state"))) {
		h.renderLocalAdminLoginWithStatus(w, r, http.StatusForbidden, "This local sign-in expired. Reload this page and try again.")
		return
	}
	cfg, err := auth.JWTConfigFromEnv()
	if err != nil {
		http.Error(w, "local sign-in is unavailable", http.StatusServiceUnavailable)
		return
	}
	signer, err := auth.NewSigner(cfg)
	if err != nil {
		http.Error(w, "local sign-in is unavailable", http.StatusServiceUnavailable)
		return
	}
	token, _, err := signer.Sign(h.localLogin.subject, h.localLogin.agencyID, h.localLogin.ttl)
	if err != nil {
		http.Error(w, "local sign-in is unavailable", http.StatusServiceUnavailable)
		return
	}
	expires := h.localLogin.now().UTC().Add(h.localLogin.ttl)
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    token,
		Path:     "/admin",
		Expires:  expires,
		MaxAge:   int(h.localLogin.ttl.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
	http.Redirect(w, r, "/admin/operations", http.StatusSeeOther)
}

func (h *handler) renderLocalAdminLoginWithStatus(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)
	h.renderLocalAdminLogin(w, r, message)
}

func (l *localAdminLogin) issueState() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate local login state: %w", err)
	}
	expires := l.now().UTC().Add(localAdminStateTTL)
	nonce := base64.RawURLEncoding.EncodeToString(raw[:])
	payload := nonce + "." + strconv.FormatInt(expires.Unix(), 10)
	signature := l.signState(payload)
	state := payload + "." + signature
	l.mu.Lock()
	l.cleanupLocked(l.now().UTC())
	l.stateStore[state] = expires
	l.mu.Unlock()
	return state, nil
}

func (l *localAdminLogin) consumeState(state string) bool {
	if state == "" || !l.validStateSignature(state) {
		return false
	}
	now := l.now().UTC()
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cleanupLocked(now)
	expires, ok := l.stateStore[state]
	if !ok || now.After(expires) {
		delete(l.stateStore, state)
		return false
	}
	delete(l.stateStore, state)
	return true
}

func (l *localAdminLogin) validStateSignature(state string) bool {
	parts := strings.Split(state, ".")
	if len(parts) != 3 {
		return false
	}
	payload := parts[0] + "." + parts[1]
	want := l.signState(payload)
	return hmac.Equal([]byte(want), []byte(parts[2]))
}

func (l *localAdminLogin) signState(payload string) string {
	mac := hmac.New(sha256.New, []byte(l.stateKey))
	_, _ = mac.Write([]byte("local-admin-login-v1\x00" + payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (l *localAdminLogin) cleanupLocked(now time.Time) {
	for state, expires := range l.stateStore {
		if now.After(expires) {
			delete(l.stateStore, state)
		}
	}
}

func friendlyDuration(d time.Duration) string {
	if d%time.Hour == 0 {
		hours := int(d / time.Hour)
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	if d%time.Minute == 0 {
		minutes := int(d / time.Minute)
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	}
	return d.String()
}

var localAdminLoginTemplate = template.Must(template.New("local-admin-login").Funcs(template.FuncMap{
	"operationsCSS": operationsCSS,
}).Parse(`<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:,"><title>Local Admin Sign-In</title><style>{{operationsCSS}}</style></head><body>
<main id="operations-main" tabindex="-1" aria-labelledby="local-login-title">
<section class="hero start-here">
<p class="app-kicker">Local demo sign-in</p>
<h1 id="local-login-title">Open Transit RT</h1>
<p class="lead">Start a private browser session for this local evaluator. This page works only on localhost when local demo sign-in is enabled.</p>
{{if .Message}}<p class="warning">{{.Message}}</p>{{end}}
<form method="post" action="/admin/local-login">
<input type="hidden" name="state" value="{{.State}}">
<table><tbody>
<tr><th>Agency</th><td><code>{{.AgencyID}}</code></td></tr>
<tr><th>Local user</th><td>{{.Subject}}</td></tr>
<tr><th>Session length</th><td>{{.TTL}}</td></tr>
<tr><th>Scope</th><td>Local/self-hosted evaluation only. This does not change public feeds, create evidence, or contact outside systems.</td></tr>
</tbody></table>
<button type="submit">Start setup</button>
</form>
<p><a href="/public/feeds.json">View local feed discovery</a></p>
<p class="muted">Production deployments must use deployment-owned admin access. Unsafe browser actions still require the private form safety check after sign-in.</p>
</section>
</main>
</body></html>`))
