package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"open-transit-rt/internal/appconfig"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEditor   Role = "editor"
	RoleOperator Role = "operator"
	RoleReadOnly Role = "read_only"

	MethodBearer = "bearer"
	MethodCookie = "cookie"
)

type Principal struct {
	Subject  string
	AgencyID string
	Roles    []Role
	Method   string
}

type contextKey struct{}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(contextKey{}).(Principal)
	return principal, ok
}

func ContextWithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, contextKey{}, principal)
}

func (p Principal) HasAny(roles ...Role) bool {
	for _, have := range p.Roles {
		for _, want := range roles {
			if have == want {
				return true
			}
		}
	}
	return false
}

type RoleStore interface {
	RolesForSubject(ctx context.Context, agencyID string, subject string) ([]Role, error)
}

type PostgresRoleStore struct {
	pool *pgxpool.Pool
}

func NewPostgresRoleStore(pool *pgxpool.Pool) *PostgresRoleStore {
	return &PostgresRoleStore{pool: pool}
}

func (s *PostgresRoleStore) RolesForSubject(ctx context.Context, agencyID string, subject string) ([]Role, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT rb.role
		FROM agency_user au
		JOIN role_binding rb
		  ON rb.agency_id = au.agency_id
		 AND rb.agency_user_id = au.id
		WHERE au.agency_id = $1
		  AND (au.auth_subject = $2 OR au.email = $2)
		ORDER BY rb.role
	`, agencyID, subject)
	if err != nil {
		return nil, fmt.Errorf("query role bindings: %w", err)
	}
	defer rows.Close()
	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role); err != nil {
			return nil, fmt.Errorf("scan role binding: %w", err)
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate role bindings: %w", err)
	}
	return roles, nil
}

type Middleware struct {
	Verifier   *Verifier
	RoleStore  RoleStore
	CSRFSecret string
}

func NewMiddleware(verifier *Verifier, roles RoleStore, csrfSecret string) *Middleware {
	return &Middleware{Verifier: verifier, RoleStore: roles, CSRFSecret: csrfSecret}
}

func MiddlewareFromEnv(pool *pgxpool.Pool) (*Middleware, error) {
	cfg, err := JWTConfigFromEnv()
	if err != nil {
		return nil, err
	}
	verifier, err := NewVerifier(cfg)
	if err != nil {
		return nil, err
	}
	csrfSecret := strings.TrimSpace(os.Getenv("CSRF_SECRET"))
	if csrfSecret == "" {
		return nil, fmt.Errorf("CSRF_SECRET is required")
	}
	if appconfig.IsProduction() && appconfig.IsDevPlaceholder(csrfSecret) {
		return nil, fmt.Errorf("CSRF_SECRET must not use a development placeholder in production")
	}
	return NewMiddleware(verifier, NewPostgresRoleStore(pool), csrfSecret), nil
}

func JWTConfigFromEnv() (JWTConfig, error) {
	ttl := DefaultAdminTokenTTL
	if raw := strings.TrimSpace(os.Getenv("ADMIN_JWT_TTL")); raw != "" {
		parsed, err := time.ParseDuration(raw)
		if err != nil {
			return JWTConfig{}, fmt.Errorf("parse ADMIN_JWT_TTL: %w", err)
		}
		ttl = parsed
	}
	skew := DefaultClockSkew
	if raw := strings.TrimSpace(os.Getenv("ADMIN_JWT_CLOCK_SKEW")); raw != "" {
		parsed, err := time.ParseDuration(raw)
		if err != nil {
			return JWTConfig{}, fmt.Errorf("parse ADMIN_JWT_CLOCK_SKEW: %w", err)
		}
		skew = parsed
	}
	secrets := []string{strings.TrimSpace(os.Getenv("ADMIN_JWT_SECRET"))}
	for _, old := range strings.Split(os.Getenv("ADMIN_JWT_OLD_SECRETS"), ",") {
		if old = strings.TrimSpace(old); old != "" {
			secrets = append(secrets, old)
		}
	}
	if appconfig.IsProduction() && appconfig.IsDevPlaceholder(secrets[0]) {
		return JWTConfig{}, fmt.Errorf("ADMIN_JWT_SECRET must not use a development placeholder in production")
	}
	return JWTConfig{
		Secrets:   secrets,
		Issuer:    strings.TrimSpace(os.Getenv("ADMIN_JWT_ISSUER")),
		Audience:  strings.TrimSpace(os.Getenv("ADMIN_JWT_AUDIENCE")),
		ClockSkew: skew,
		TTL:       ttl,
	}, nil
}

func (m *Middleware) Require(roles ...Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, err := m.authenticate(r)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if !principal.HasAny(roles...) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			if principal.Method == MethodCookie && unsafeMethod(r.Method) {
				if !m.validCSRF(r, principal) {
					http.Error(w, "invalid csrf token", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r.WithContext(ContextWithPrincipal(r.Context(), principal)))
		})
	}
}

func (m *Middleware) authenticate(r *http.Request) (Principal, error) {
	token, method := bearerToken(r)
	if token == "" {
		token, method = cookieToken(r)
	}
	if token == "" {
		return Principal{}, ErrInvalidToken
	}
	claims, err := m.Verifier.Verify(token)
	if err != nil {
		return Principal{}, err
	}
	roles, err := m.RoleStore.RolesForSubject(r.Context(), claims.AgencyID, claims.Subject)
	if err != nil {
		return Principal{}, err
	}
	if len(roles) == 0 {
		return Principal{}, errors.New("no roles")
	}
	return Principal{Subject: claims.Subject, AgencyID: claims.AgencyID, Roles: roles, Method: method}, nil
}

func bearerToken(r *http.Request) (string, string) {
	raw := strings.TrimSpace(r.Header.Get("Authorization"))
	if raw == "" {
		return "", ""
	}
	fields := strings.Fields(raw)
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		return "", ""
	}
	return fields[1], MethodBearer
}

func cookieToken(r *http.Request) (string, string) {
	cookie, err := r.Cookie("admin_session")
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return "", ""
	}
	return cookie.Value, MethodCookie
}

func unsafeMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func (m *Middleware) validCSRF(r *http.Request, principal Principal) bool {
	token := strings.TrimSpace(r.Header.Get("X-CSRF-Token"))
	if token == "" {
		_ = r.ParseForm()
		token = strings.TrimSpace(r.FormValue("csrf_token"))
	}
	return hmac.Equal([]byte(token), []byte(CSRFToken(m.CSRFSecret, principal)))
}

func CSRFToken(secret string, principal Principal) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte("csrf-v1\x00" + principal.Subject + "\x00" + principal.AgencyID))
	return hex.EncodeToString(mac.Sum(nil))
}

func RequireRole(w http.ResponseWriter, r *http.Request, roles ...Role) (Principal, bool) {
	principal, ok := PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return Principal{}, false
	}
	if !principal.HasAny(roles...) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return Principal{}, false
	}
	return principal, true
}

func RequireAgencyQueryMatch(w http.ResponseWriter, r *http.Request, principal Principal) bool {
	requestAgency := strings.TrimSpace(r.URL.Query().Get("agency_id"))
	if requestAgency != "" && requestAgency != principal.AgencyID {
		http.Error(w, "agency_id conflicts with authenticated agency", http.StatusForbidden)
		return false
	}
	return true
}

func RejectAgencyConflict(w http.ResponseWriter, agencyID string, principal Principal) bool {
	if strings.TrimSpace(agencyID) != "" && agencyID != principal.AgencyID {
		http.Error(w, "agency_id conflicts with authenticated agency", http.StatusForbidden)
		return true
	}
	return false
}

type StaticRoleStore struct {
	Roles []Role
}

func (s StaticRoleStore) RolesForSubject(context.Context, string, string) ([]Role, error) {
	return append([]Role(nil), s.Roles...), nil
}

type TestAuthenticator struct {
	Principal Principal
}

func (t TestAuthenticator) Require(roles ...Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !t.Principal.HasAny(roles...) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r.WithContext(ContextWithPrincipal(r.Context(), t.Principal)))
		})
	}
}

func RowsToRoles(rows pgx.Rows) ([]Role, error) {
	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}
