package devices

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

func TestHashTokenUsesPepper(t *testing.T) {
	left := NewPostgresStore(nil, Config{TokenPepper: "pepper-a"}).HashToken("device-token")
	right := NewPostgresStore(nil, Config{TokenPepper: "pepper-b"}).HashToken("device-token")
	if left == right {
		t.Fatalf("hash did not change when pepper changed")
	}
	if len(left) <= len("hmac-sha256:") || left[:12] != "hmac-sha256:" {
		t.Fatalf("hash = %q, want hmac-sha256 prefix", left)
	}
}

func TestGenerateTokenReturnsOpaqueValue(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	if len(token) < 32 {
		t.Fatalf("token length = %d, want opaque random token", len(token))
	}
}

func TestPostgresDeviceRebindInvalidatesOldTokenImmediately(t *testing.T) {
	ctx := context.Background()
	pool, cleanup := setupDeviceIntegrationDB(t)
	defer cleanup()

	store := NewPostgresStore(pool, Config{TokenPepper: "test-pepper"})
	oldToken := "old-device-token"
	seedDevice(t, ctx, pool, store, oldToken, "bus-1")

	if _, err := store.Verify(ctx, VerifyInput{Token: oldToken, AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1"}); err != nil {
		t.Fatalf("verify old token before rebind: %v", err)
	}
	if _, err := store.Verify(ctx, VerifyInput{Token: oldToken, AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-spoof"}); err == nil {
		t.Fatalf("spoofed vehicle binding verified successfully")
	}

	result, err := store.Rebind(ctx, RebindInput{
		AgencyID:  "demo-agency",
		DeviceID:  "device-1",
		VehicleID: "bus-2",
		ActorID:   "admin@example.com",
		Reason:    "rotation test",
		Now:       time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("rebind: %v", err)
	}
	if result.Token == "" || result.VehicleID != "bus-2" {
		t.Fatalf("result = %+v, want new token bound to bus-2", result)
	}
	if _, err := store.Verify(ctx, VerifyInput{Token: oldToken, AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1"}); err == nil {
		t.Fatalf("old token verified after rebind")
	}
	if _, err := store.Verify(ctx, VerifyInput{Token: result.Token, AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1"}); err == nil {
		t.Fatalf("new token verified for old vehicle binding")
	}
	if _, err := store.Verify(ctx, VerifyInput{Token: result.Token, AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-2"}); err != nil {
		t.Fatalf("new token did not verify for new binding: %v", err)
	}
	var auditRows int
	if err := pool.QueryRow(ctx, `
		SELECT count(*)
		FROM audit_log
		WHERE agency_id = 'demo-agency'
		  AND action = 'device.rebind'
		  AND entity_id = 'device-1'
	`).Scan(&auditRows); err != nil {
		t.Fatalf("count audit rows: %v", err)
	}
	if auditRows != 1 {
		t.Fatalf("audit rows = %d, want 1", auditRows)
	}
}

func seedDevice(t *testing.T, ctx context.Context, pool *pgxpool.Pool, store *PostgresStore, token string, vehicleID string) {
	t.Helper()
	if _, err := pool.Exec(ctx, `
		INSERT INTO agency (id, name, timezone)
		VALUES ('demo-agency', 'Demo Agency', 'America/Vancouver')
		ON CONFLICT (id) DO NOTHING
	`); err != nil {
		t.Fatalf("seed agency: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO device_credential (agency_id, device_id, vehicle_id, token_hash, status)
		VALUES ('demo-agency', 'device-1', $1, $2, 'active')
		ON CONFLICT (agency_id, device_id) DO UPDATE
		SET vehicle_id = EXCLUDED.vehicle_id,
		    token_hash = EXCLUDED.token_hash,
		    status = 'active',
		    rotated_at = NULL,
		    revoked_at = NULL
	`, vehicleID, store.HashToken(token)); err != nil {
		t.Fatalf("seed device: %v", err)
	}
}

func setupDeviceIntegrationDB(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("set INTEGRATION_TESTS=1 to run DB-backed device tests")
	}
	targetURL := os.Getenv("TEST_DATABASE_URL")
	if targetURL == "" {
		targetURL = "postgres://postgres:postgres@localhost:55432/open_transit_rt_test?sslmode=disable"
	}
	parsed, err := url.Parse(targetURL)
	if err != nil {
		t.Fatalf("parse TEST_DATABASE_URL: %v", err)
	}
	adminURL := *parsed
	adminURL.Path = "/postgres"
	adminDB, err := sql.Open("pgx", adminURL.String())
	if err != nil {
		t.Fatalf("open admin database: %v", err)
	}
	defer adminDB.Close()
	if err := adminDB.Ping(); err != nil {
		t.Fatalf("ping admin database: %v", err)
	}
	dbName := fmt.Sprintf("otrt_devices_test_%d", time.Now().UnixNano())
	if _, err := adminDB.Exec(`CREATE DATABASE ` + quoteIdentifier(dbName)); err != nil {
		t.Fatalf("create temporary database: %v", err)
	}
	cleanupDB := func() {
		adminDB, err := sql.Open("pgx", adminURL.String())
		if err == nil {
			defer adminDB.Close()
			_, _ = adminDB.Exec(`DROP DATABASE IF EXISTS ` + quoteIdentifier(dbName) + ` WITH (FORCE)`)
		}
	}
	testURL := *parsed
	testURL.Path = "/" + dbName
	if err := applyDeviceMigrations(testURL.String()); err != nil {
		cleanupDB()
		t.Fatalf("apply migrations: %v", err)
	}
	pool, err := pgxpool.New(context.Background(), testURL.String())
	if err != nil {
		cleanupDB()
		t.Fatalf("create test pool: %v", err)
	}
	return pool, func() {
		pool.Close()
		cleanupDB()
	}
}

func applyDeviceMigrations(databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, filepath.Join("..", "..", "db", "migrations"))
}

func quoteIdentifier(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
