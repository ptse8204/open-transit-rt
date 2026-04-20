package telemetry

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

func TestAdvisoryLockKey(t *testing.T) {
	first := advisoryLockKey("demo-agency", "bus-10")
	second := advisoryLockKey("demo-agency", "bus-10")
	other := advisoryLockKey("demo-agency", "bus-11")

	if first != second {
		t.Fatalf("advisory lock key is not deterministic: %d != %d", first, second)
	}
	if first == other {
		t.Fatalf("different vehicle streams unexpectedly produced the same test key: %d", first)
	}
}

func TestPostgresRepositoryIntegration(t *testing.T) {
	ctx := context.Background()
	pool, cleanup := setupIntegrationDB(t)
	defer cleanup()

	repo := NewPostgresRepository(pool)

	t.Run("valid insert and latest query", func(t *testing.T) {
		resetTelemetry(t, ctx, pool)
		events := loadFixture(t, "../../testdata/telemetry/matched-vehicle.json")

		first := storeFixtureEvent(t, ctx, repo, events[0])
		if first.IngestStatus != IngestStatusAccepted {
			t.Fatalf("status = %s, want accepted", first.IngestStatus)
		}
		if first.ID == 0 || first.ReceivedAt.IsZero() {
			t.Fatalf("insert did not return DB id and received_at: %+v", first)
		}
		assertPayloadVehicle(t, first.PayloadJSON, events[0].VehicleID)

		second := storeFixtureEvent(t, ctx, repo, events[1])
		if second.IngestStatus != IngestStatusAccepted {
			t.Fatalf("status = %s, want accepted", second.IngestStatus)
		}

		latest, err := repo.LatestByVehicle(ctx, events[0].AgencyID, events[0].VehicleID)
		if err != nil {
			t.Fatalf("latest by vehicle: %v", err)
		}
		if !latest.Timestamp.Equal(events[1].Timestamp) {
			t.Fatalf("latest timestamp = %s, want %s", latest.Timestamp, events[1].Timestamp)
		}
	})

	t.Run("duplicate same payload", func(t *testing.T) {
		resetTelemetry(t, ctx, pool)
		event := loadFixture(t, "../../testdata/telemetry/matched-vehicle.json")[0]

		storeFixtureEvent(t, ctx, repo, event)
		duplicate := storeFixtureEvent(t, ctx, repo, event)
		if duplicate.IngestStatus != IngestStatusDuplicate {
			t.Fatalf("status = %s, want duplicate", duplicate.IngestStatus)
		}
	})

	t.Run("equal timestamp with different payload", func(t *testing.T) {
		resetTelemetry(t, ctx, pool)
		event := loadFixture(t, "../../testdata/telemetry/matched-vehicle.json")[0]
		changed := event
		changed.Lat = event.Lat + 0.001
		changed.SpeedMPS = event.SpeedMPS + 1

		storeFixtureEvent(t, ctx, repo, event)
		duplicate := storeFixtureEvent(t, ctx, repo, changed)
		if duplicate.IngestStatus != IngestStatusDuplicate {
			t.Fatalf("status = %s, want duplicate", duplicate.IngestStatus)
		}
	})

	t.Run("equal vehicle timestamp from different device", func(t *testing.T) {
		resetTelemetry(t, ctx, pool)
		event := loadFixture(t, "../../testdata/telemetry/matched-vehicle.json")[0]
		otherDevice := event
		otherDevice.DeviceID = "android-tablet-other"

		storeFixtureEvent(t, ctx, repo, event)
		duplicate := storeFixtureEvent(t, ctx, repo, otherDevice)
		if duplicate.IngestStatus != IngestStatusDuplicate {
			t.Fatalf("status = %s, want duplicate", duplicate.IngestStatus)
		}
	})

	t.Run("out of order after newer accepted event", func(t *testing.T) {
		resetTelemetry(t, ctx, pool)
		events := loadFixture(t, "../../testdata/telemetry/matched-vehicle.json")

		storeFixtureEvent(t, ctx, repo, events[1])
		outOfOrder := storeFixtureEvent(t, ctx, repo, events[0])
		if outOfOrder.IngestStatus != IngestStatusOutOfOrder {
			t.Fatalf("status = %s, want out_of_order", outOfOrder.IngestStatus)
		}

		latest, err := repo.LatestByVehicle(ctx, events[0].AgencyID, events[0].VehicleID)
		if err != nil {
			t.Fatalf("latest by vehicle: %v", err)
		}
		if !latest.Timestamp.Equal(events[1].Timestamp) {
			t.Fatalf("latest timestamp = %s, want newer %s", latest.Timestamp, events[1].Timestamp)
		}
	})

	t.Run("agency scoped latest per vehicle listing", func(t *testing.T) {
		resetTelemetry(t, ctx, pool)
		for _, event := range loadFixture(t, "../../testdata/telemetry/matched-vehicle.json") {
			storeFixtureEvent(t, ctx, repo, event)
		}
		for _, event := range loadFixture(t, "../../testdata/telemetry/swapped-vehicle.json") {
			storeFixtureEvent(t, ctx, repo, event)
		}
		for _, event := range loadFixture(t, "../../testdata/telemetry/after-midnight.json") {
			storeFixtureEvent(t, ctx, repo, event)
		}

		latest, err := repo.ListLatestByAgency(ctx, "demo-agency", 10)
		if err != nil {
			t.Fatalf("list latest: %v", err)
		}
		if len(latest) != 3 {
			t.Fatalf("latest count = %d, want one row for each of 3 demo vehicles", len(latest))
		}
		for _, event := range latest {
			if event.AgencyID != "demo-agency" {
				t.Fatalf("unscoped agency row returned: %+v", event)
			}
			if event.IngestStatus != IngestStatusAccepted {
				t.Fatalf("latest included non-accepted status: %+v", event)
			}
		}
	})

	t.Run("agency scoped debug events include persisted statuses", func(t *testing.T) {
		resetTelemetry(t, ctx, pool)
		events := loadFixture(t, "../../testdata/telemetry/matched-vehicle.json")
		storeFixtureEvent(t, ctx, repo, events[1])
		storeFixtureEvent(t, ctx, repo, events[0])

		debugEvents, err := repo.ListEvents(ctx, "demo-agency", 10)
		if err != nil {
			t.Fatalf("list events: %v", err)
		}
		if len(debugEvents) != 2 {
			t.Fatalf("debug event count = %d, want 2", len(debugEvents))
		}
		if debugEvents[0].AgencyID != "demo-agency" || debugEvents[1].AgencyID != "demo-agency" {
			t.Fatalf("debug events not agency scoped: %+v", debugEvents)
		}
		if debugEvents[0].IngestStatus != IngestStatusOutOfOrder {
			t.Fatalf("newest debug row status = %s, want out_of_order", debugEvents[0].IngestStatus)
		}
	})

	t.Run("unknown agency rejection", func(t *testing.T) {
		resetTelemetry(t, ctx, pool)
		event := loadFixture(t, "../../testdata/telemetry/matched-vehicle.json")[0]
		event.AgencyID = "missing-agency"

		_, err := repo.Store(ctx, event, mustMarshalEvent(t, event))
		if !errors.Is(err, ErrUnknownAgency) {
			t.Fatalf("err = %v, want ErrUnknownAgency", err)
		}
		events, listErr := repo.ListEvents(ctx, "missing-agency", 10)
		if listErr != nil {
			t.Fatalf("list missing-agency events: %v", listErr)
		}
		if len(events) != 0 {
			t.Fatalf("unknown agency persisted rows: %+v", events)
		}
	})
}

func setupIntegrationDB(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("set INTEGRATION_TESTS=1 to run DB-backed telemetry tests")
	}

	targetURL := os.Getenv("TEST_DATABASE_URL")
	if targetURL == "" {
		targetURL = "postgres://postgres:postgres@localhost:55432/open_transit_rt_test?sslmode=disable"
	}

	if pool, cleanup, err := setupTemporaryDatabase(t, targetURL); err == nil {
		return pool, cleanup
	} else {
		t.Logf("isolated database setup unavailable, falling back to schema mode: %v", err)
	}

	return setupTemporarySchema(t, targetURL)
}

func setupTemporaryDatabase(t *testing.T, targetURL string) (*pgxpool.Pool, func(), error) {
	t.Helper()

	parsed, err := url.Parse(targetURL)
	if err != nil {
		return nil, nil, fmt.Errorf("parse TEST_DATABASE_URL: %w", err)
	}

	adminURL := *parsed
	adminURL.Path = "/postgres"
	adminDB, err := sql.Open("pgx", adminURL.String())
	if err != nil {
		return nil, nil, fmt.Errorf("open admin database: %w", err)
	}
	defer adminDB.Close()
	if err := adminDB.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping admin database: %w", err)
	}

	dbName := fmt.Sprintf("otrt_test_%d", time.Now().UnixNano())
	if _, err := adminDB.Exec(`CREATE DATABASE ` + quoteIdentifier(dbName)); err != nil {
		return nil, nil, fmt.Errorf("create temporary database: %w", err)
	}

	testURL := *parsed
	testURL.Path = "/" + dbName
	cleanup := func() {
		adminDB, err := sql.Open("pgx", adminURL.String())
		if err == nil {
			defer adminDB.Close()
			_, _ = adminDB.Exec(`DROP DATABASE IF EXISTS ` + quoteIdentifier(dbName) + ` WITH (FORCE)`)
		}
	}

	if err := applyMigrations(testURL.String()); err != nil {
		cleanup()
		return nil, nil, err
	}

	pool, err := pgxpool.New(context.Background(), testURL.String())
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("create test pool: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		cleanup()
		return nil, nil, fmt.Errorf("ping test pool: %w", err)
	}
	seedAgencies(t, context.Background(), pool)

	return pool, func() {
		pool.Close()
		cleanup()
	}, nil
}

func setupTemporarySchema(t *testing.T, targetURL string) (*pgxpool.Pool, func()) {
	t.Helper()
	parsed, err := url.Parse(targetURL)
	if err != nil {
		t.Fatalf("parse TEST_DATABASE_URL for schema fallback: %v", err)
	}

	baseDB, err := sql.Open("pgx", targetURL)
	if err != nil {
		t.Fatalf("open fallback test database: %v", err)
	}
	if err := baseDB.Ping(); err != nil {
		_ = baseDB.Close()
		t.Fatalf("fallback test database is unavailable: %v", err)
	}

	schemaName := fmt.Sprintf("otrt_test_%d", time.Now().UnixNano())
	if _, err := baseDB.Exec(`CREATE SCHEMA ` + quoteIdentifier(schemaName)); err != nil {
		t.Fatalf("create fallback schema: %v", err)
	}

	query := parsed.Query()
	query.Set("search_path", schemaName+",public")
	parsed.RawQuery = query.Encode()

	cleanup := func() {
		_, _ = baseDB.Exec(`DROP SCHEMA IF EXISTS ` + quoteIdentifier(schemaName) + ` CASCADE`)
		_ = baseDB.Close()
	}

	if err := applyMigrations(parsed.String()); err != nil {
		cleanup()
		t.Fatalf("apply migrations to fallback schema: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), parsed.String())
	if err != nil {
		cleanup()
		t.Fatalf("create fallback test pool: %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		cleanup()
		t.Fatalf("ping fallback test pool: %v", err)
	}
	seedAgencies(t, context.Background(), pool)

	return pool, func() {
		pool.Close()
		cleanup()
	}
}

func applyMigrations(databaseURL string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set migration dialect: %w", err)
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open migration database: %w", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping migration database: %w", err)
	}
	if err := goose.Up(db, migrationDir()); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

func migrationDir() string {
	if raw := os.Getenv("MIGRATIONS_DIR"); raw != "" {
		if _, err := os.Stat(raw); err == nil {
			return raw
		}
		if _, err := os.Stat(filepath.Join("..", "..", raw)); err == nil {
			return filepath.Join("..", "..", raw)
		}
	}
	return filepath.Join("..", "..", "db", "migrations")
}

func seedAgencies(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	for _, agency := range []struct {
		id       string
		name     string
		timezone string
	}{
		{id: "demo-agency", name: "Demo Agency", timezone: "America/Vancouver"},
		{id: "overnight-agency", name: "Overnight Agency", timezone: "America/Vancouver"},
		{id: "freq-agency", name: "Frequency Agency", timezone: "America/Vancouver"},
	} {
		_, err := pool.Exec(ctx, `
			INSERT INTO agency (id, name, timezone, contact_email, public_url)
			VALUES ($1, $2, $3, 'dev@example.com', 'http://localhost')
			ON CONFLICT (id) DO UPDATE
			SET name = EXCLUDED.name,
			    timezone = EXCLUDED.timezone
		`, agency.id, agency.name, agency.timezone)
		if err != nil {
			t.Fatalf("seed agency %s: %v", agency.id, err)
		}
	}
}

func resetTelemetry(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	if _, err := pool.Exec(ctx, `DELETE FROM telemetry_event`); err != nil {
		t.Fatalf("reset telemetry: %v", err)
	}
}

func loadFixture(t *testing.T, path string) []Event {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	var events []Event
	if err := json.Unmarshal(raw, &events); err != nil {
		t.Fatalf("decode fixture %s: %v", path, err)
	}
	return events
}

func storeFixtureEvent(t *testing.T, ctx context.Context, repo *PostgresRepository, event Event) StoredEvent {
	t.Helper()
	result, err := repo.Store(ctx, event, mustMarshalEvent(t, event))
	if err != nil {
		t.Fatalf("store fixture event: %v", err)
	}
	return result.StoredEvent
}

func mustMarshalEvent(t *testing.T, event Event) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	return raw
}

func assertPayloadVehicle(t *testing.T, payload json.RawMessage, vehicleID string) {
	t.Helper()
	var parsed struct {
		VehicleID string `json:"vehicle_id"`
	}
	if err := json.Unmarshal(payload, &parsed); err != nil {
		t.Fatalf("payload is not JSON: %v", err)
	}
	if parsed.VehicleID != vehicleID {
		t.Fatalf("payload vehicle_id = %s, want %s", parsed.VehicleID, vehicleID)
	}
}

func quoteIdentifier(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
