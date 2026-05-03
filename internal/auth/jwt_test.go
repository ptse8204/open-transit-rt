package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func TestJWTRequiresCoreClaimsAndValidatesAudienceIssuer(t *testing.T) {
	cfg := JWTConfig{Secrets: []string{"test-secret"}, Issuer: "open-transit-rt", Audience: "admin-api", ClockSkew: time.Minute, TTL: time.Hour}
	signer, err := NewSigner(cfg)
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}
	verifier, err := NewVerifier(cfg)
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}
	token, claims, err := signer.Sign("operator@example.com", "demo-agency", time.Hour)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if claims.JTI == "" {
		t.Fatalf("jti is empty")
	}
	verified, err := verifier.Verify(token)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if verified.Subject != "operator@example.com" || verified.AgencyID != "demo-agency" || verified.Issuer != cfg.Issuer || verified.Audience != cfg.Audience {
		t.Fatalf("claims = %+v, want scoped admin claims", verified)
	}

	wrongAudience, _ := NewVerifier(JWTConfig{Secrets: []string{"test-secret"}, Issuer: "open-transit-rt", Audience: "other", ClockSkew: time.Minute})
	if _, err := wrongAudience.Verify(token); err == nil {
		t.Fatalf("verify succeeded with wrong audience")
	}
}

func TestJWTAcceptsOldSecretForRotation(t *testing.T) {
	oldCfg := JWTConfig{Secrets: []string{"old-secret"}, Issuer: "open-transit-rt", Audience: "admin-api"}
	signer, err := NewSigner(oldCfg)
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}
	token, _, err := signer.Sign("operator@example.com", "demo-agency", time.Hour)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	verifier, err := NewVerifier(JWTConfig{Secrets: []string{"new-secret", "old-secret"}, Issuer: "open-transit-rt", Audience: "admin-api"})
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}
	if _, err := verifier.Verify(token); err != nil {
		t.Fatalf("verify with old secret during rotation: %v", err)
	}
}

func TestCSRFTokenBindsToPrincipal(t *testing.T) {
	principal := Principal{Subject: "operator@example.com", AgencyID: "demo-agency"}
	token := CSRFToken("csrf-secret", principal)
	if token == "" {
		t.Fatalf("csrf token is empty")
	}
	other := CSRFToken("csrf-secret", Principal{Subject: principal.Subject, AgencyID: "other-agency"})
	if token == other {
		t.Fatalf("csrf token did not bind to agency")
	}
}

func TestPostgresRoleStoreScopesRolesByClaimAgency(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("set INTEGRATION_TESTS=1 to run DB-backed auth role tests")
	}
	ctx := context.Background()
	pool, cleanup := setupAuthIntegrationDB(t)
	defer cleanup()

	_, err := pool.Exec(ctx, `
		INSERT INTO agency (id, name, timezone)
		VALUES ('agency-a', 'Agency A', 'America/Los_Angeles'),
		       ('agency-b', 'Agency B', 'America/Los_Angeles');
		WITH a_user AS (
			INSERT INTO agency_user (agency_id, email, auth_subject)
			VALUES ('agency-a', 'operator@example.com', 'shared-subject')
			RETURNING id
		)
		INSERT INTO role_binding (agency_id, agency_user_id, role)
		SELECT 'agency-a', id, 'read_only' FROM a_user;
		WITH b_user AS (
			INSERT INTO agency_user (agency_id, email, auth_subject)
			VALUES ('agency-b', 'operator@example.com', 'shared-subject')
			RETURNING id
		)
		INSERT INTO role_binding (agency_id, agency_user_id, role)
		SELECT 'agency-b', id, 'admin' FROM b_user;
	`)
	if err != nil {
		t.Fatalf("seed role bindings: %v", err)
	}

	store := NewPostgresRoleStore(pool)
	rolesA, err := store.RolesForSubject(ctx, "agency-a", "shared-subject")
	if err != nil {
		t.Fatalf("roles agency-a: %v", err)
	}
	if len(rolesA) != 1 || rolesA[0] != RoleReadOnly {
		t.Fatalf("agency-a roles = %+v, want only read_only", rolesA)
	}
	rolesB, err := store.RolesForSubject(ctx, "agency-b", "shared-subject")
	if err != nil {
		t.Fatalf("roles agency-b: %v", err)
	}
	if len(rolesB) != 1 || rolesB[0] != RoleAdmin {
		t.Fatalf("agency-b roles = %+v, want only admin", rolesB)
	}
}

func setupAuthIntegrationDB(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	targetURL := os.Getenv("TEST_DATABASE_URL")
	if targetURL == "" {
		targetURL = "postgres://postgres:postgres@localhost:55432/open_transit_rt_test?sslmode=disable"
	}
	parsed, err := url.Parse(targetURL)
	if err != nil {
		t.Fatalf("parse TEST_DATABASE_URL: %v", err)
	}
	admin := *parsed
	admin.Path = "/postgres"
	adminDB, err := sql.Open("pgx", admin.String())
	if err != nil {
		t.Fatalf("open admin database: %v", err)
	}
	defer adminDB.Close()
	dbName := fmt.Sprintf("otrt_auth_test_%d", time.Now().UnixNano())
	if _, err := adminDB.Exec(`CREATE DATABASE ` + quoteAuthIdent(dbName)); err != nil {
		t.Fatalf("create temporary database: %v", err)
	}
	testURL := *parsed
	testURL.Path = "/" + dbName
	db, err := sql.Open("pgx", testURL.String())
	if err != nil {
		t.Fatalf("open temporary database: %v", err)
	}
	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("set goose dialect: %v", err)
	}
	if err := goose.Up(db, authMigrationDir()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	_ = db.Close()
	pool, err := pgxpool.New(context.Background(), testURL.String())
	if err != nil {
		t.Fatalf("connect temporary database: %v", err)
	}
	return pool, func() {
		pool.Close()
		adminDB, err := sql.Open("pgx", admin.String())
		if err == nil {
			_, _ = adminDB.Exec(`DROP DATABASE IF EXISTS ` + quoteAuthIdent(dbName) + ` WITH (FORCE)`)
			_ = adminDB.Close()
		}
	}
}

func authMigrationDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "db/migrations"
	}
	for {
		candidate := filepath.Join(dir, "db", "migrations")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "db/migrations"
		}
		dir = parent
	}
}

func quoteAuthIdent(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
