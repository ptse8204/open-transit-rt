package state

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	appgtfs "open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/telemetry"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func TestPostgresMatcherIntegration(t *testing.T) {
	ctx := context.Background()
	pool, cleanup := setupStateIntegrationDB(t)
	defer cleanup()

	schedules := appgtfs.NewPostgresRepository(pool)
	assignments := NewPostgresRepository(pool)
	telemetryRepo := telemetry.NewPostgresRepository(pool)
	engine := NewEngine(schedules, assignments, DefaultConfig())

	t.Run("matches persisted latest telemetry", func(t *testing.T) {
		resetMatcherData(t, ctx, pool)
		seedDemoSchedule(t, ctx, pool, true)

		event := telemetry.Event{
			AgencyID:  "demo-agency",
			DeviceID:  "device-bus-10",
			VehicleID: "bus-10",
			Timestamp: time.Date(2026, 4, 20, 15, 2, 0, 0, time.UTC),
			Lat:       49.2827,
			Lon:       -123.1207,
			Bearing:   120,
			SpeedMPS:  8,
			AccuracyM: 8,
			TripHint:  "trip-10-0800",
		}
		storeTelemetry(t, ctx, telemetryRepo, event)
		latest, err := telemetryRepo.LatestByVehicle(ctx, "demo-agency", "bus-10")
		if err != nil {
			t.Fatalf("latest telemetry: %v", err)
		}

		assignment, err := engine.MatchEvent(ctx, latest, latest.Timestamp.Add(30*time.Second))
		if err != nil {
			t.Fatalf("match event: %v", err)
		}
		if assignment.State != StateInService || assignment.TripID != "trip-10-0800" {
			t.Fatalf("assignment = %+v, want matched trip", assignment)
		}
		if assignment.TelemetryEventID != latest.ID {
			t.Fatalf("telemetry_event_id = %d, want %d", assignment.TelemetryEventID, latest.ID)
		}

		current, err := assignments.CurrentAssignment(ctx, "demo-agency", "bus-10")
		if err != nil {
			t.Fatalf("current assignment: %v", err)
		}
		if current == nil || current.TripID != "trip-10-0800" || current.DegradedState != DegradedNone {
			t.Fatalf("persisted current = %+v, want healthy matched assignment", current)
		}
	})

	t.Run("stale telemetry writes unknown assignment and incident", func(t *testing.T) {
		resetMatcherData(t, ctx, pool)
		seedDemoSchedule(t, ctx, pool, true)

		_, err := assignments.SaveAssignment(ctx, Assignment{
			AgencyID:         "demo-agency",
			VehicleID:        "bus-stale",
			State:            StateInService,
			ServiceDate:      "20260420",
			RouteID:          "route-10",
			TripID:           "trip-10-0800",
			BlockID:          "block-10",
			StartDate:        "20260420",
			StartTime:        "08:00:00",
			Confidence:       0.9,
			AssignmentSource: AssignmentSourceAutomatic,
			ReasonCodes:      []string{ReasonTripHintMatch},
			DegradedState:    DegradedNone,
			ScoreDetails:     map[string]any{"score_schema": "loose_debug_v1"},
			ActiveFrom:       time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC),
		}, nil)
		if err != nil {
			t.Fatalf("seed previous assignment: %v", err)
		}

		event := telemetry.Event{
			AgencyID:  "demo-agency",
			DeviceID:  "device-bus-stale",
			VehicleID: "bus-stale",
			Timestamp: time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC),
			Lat:       49.2827,
			Lon:       -123.1207,
			TripHint:  "trip-10-0800",
		}
		stored := storeTelemetry(t, ctx, telemetryRepo, event)
		assignment, err := engine.MatchEvent(ctx, stored, stored.Timestamp.Add(2*time.Minute))
		if err != nil {
			t.Fatalf("match stale event: %v", err)
		}
		if assignment.State != StateUnknown || assignment.DegradedState != DegradedStale {
			t.Fatalf("assignment = %+v, want stale unknown", assignment)
		}

		var staleIncidents int
		if err := pool.QueryRow(ctx, `
			SELECT count(*)
			FROM incident
			WHERE agency_id = 'demo-agency'
			  AND vehicle_id = 'bus-stale'
			  AND incident_type = 'stale_telemetry'
			  AND vehicle_trip_assignment_id = $1
		`, assignment.ID).Scan(&staleIncidents); err != nil {
			t.Fatalf("count stale incidents: %v", err)
		}
		if staleIncidents != 1 {
			t.Fatalf("stale incident count = %d, want 1", staleIncidents)
		}
	})

	t.Run("manual override wins over automatic matching", func(t *testing.T) {
		resetMatcherData(t, ctx, pool)
		seedDemoSchedule(t, ctx, pool, true)
		_, err := pool.Exec(ctx, `
			INSERT INTO manual_override (
				agency_id, vehicle_id, override_type, route_id, trip_id, start_date, start_time, state, reason, created_by, created_at
			)
			VALUES (
				'demo-agency', 'bus-override', 'trip_assignment', 'route-10', 'trip-manual', '20260420', '08:30:00', 'in_service', 'test override', 'test', now()
			)
		`)
		if err != nil {
			t.Fatalf("insert manual override: %v", err)
		}

		event := telemetry.Event{
			AgencyID:  "demo-agency",
			DeviceID:  "device-bus-override",
			VehicleID: "bus-override",
			Timestamp: time.Date(2026, 4, 20, 15, 2, 0, 0, time.UTC),
			Lat:       49.2827,
			Lon:       -123.1207,
			TripHint:  "trip-10-0800",
		}
		stored := storeTelemetry(t, ctx, telemetryRepo, event)
		assignment, err := engine.MatchEvent(ctx, stored, stored.Timestamp.Add(30*time.Second))
		if err != nil {
			t.Fatalf("match override event: %v", err)
		}
		if assignment.AssignmentSource != AssignmentSourceManualOverride || assignment.TripID != "trip-manual" {
			t.Fatalf("assignment = %+v, want manual override", assignment)
		}
	})

	t.Run("missing shape degrades but can still match", func(t *testing.T) {
		resetMatcherData(t, ctx, pool)
		seedDemoSchedule(t, ctx, pool, false)
		_, err := assignments.SaveAssignment(ctx, Assignment{
			AgencyID:         "demo-agency",
			VehicleID:        "bus-no-shape",
			State:            StateInService,
			ServiceDate:      "20260420",
			RouteID:          "route-10",
			TripID:           "trip-10-0800",
			BlockID:          "block-10",
			StartDate:        "20260420",
			StartTime:        "08:00:00",
			Confidence:       0.9,
			AssignmentSource: AssignmentSourceAutomatic,
			ReasonCodes:      []string{ReasonContinuityMatch},
			DegradedState:    DegradedNone,
			ScoreDetails:     map[string]any{"score_schema": "loose_debug_v1"},
			ActiveFrom:       time.Date(2026, 4, 20, 14, 59, 0, 0, time.UTC),
		}, nil)
		if err != nil {
			t.Fatalf("seed continuity assignment: %v", err)
		}

		event := telemetry.Event{
			AgencyID:  "demo-agency",
			DeviceID:  "device-bus-no-shape",
			VehicleID: "bus-no-shape",
			Timestamp: time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC),
			Lat:       49.2827,
			Lon:       -123.1207,
			TripHint:  "trip-10-0800",
		}
		stored := storeTelemetry(t, ctx, telemetryRepo, event)
		assignment, err := engine.MatchEvent(ctx, stored, stored.Timestamp.Add(30*time.Second))
		if err != nil {
			t.Fatalf("match missing shape event: %v", err)
		}
		if assignment.State != StateInService || !hasReason(assignment, ReasonMissingShape) {
			t.Fatalf("assignment = %+v, want degraded in-service match with missing_shape reason", assignment)
		}
	})
}

func setupStateIntegrationDB(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("set INTEGRATION_TESTS=1 to run DB-backed matcher tests")
	}

	targetURL := os.Getenv("TEST_DATABASE_URL")
	if targetURL == "" {
		targetURL = "postgres://postgres:postgres@localhost:55432/open_transit_rt_test?sslmode=disable"
	}

	if pool, cleanup, err := setupStateTemporaryDatabase(t, targetURL); err == nil {
		return pool, cleanup
	} else {
		t.Logf("isolated database setup unavailable, falling back to schema mode: %v", err)
	}
	return setupStateTemporarySchema(t, targetURL)
}

func setupStateTemporaryDatabase(t *testing.T, targetURL string) (*pgxpool.Pool, func(), error) {
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

	dbName := fmt.Sprintf("otrt_state_test_%d", time.Now().UnixNano())
	if _, err := adminDB.Exec(`CREATE DATABASE ` + quoteStateIdentifier(dbName)); err != nil {
		return nil, nil, fmt.Errorf("create temporary database: %w", err)
	}
	cleanup := func() {
		adminDB, err := sql.Open("pgx", adminURL.String())
		if err == nil {
			defer adminDB.Close()
			_, _ = adminDB.Exec(`DROP DATABASE IF EXISTS ` + quoteStateIdentifier(dbName) + ` WITH (FORCE)`)
		}
	}

	testURL := *parsed
	testURL.Path = "/" + dbName
	if err := applyStateMigrations(testURL.String()); err != nil {
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
	return pool, func() {
		pool.Close()
		cleanup()
	}, nil
}

func setupStateTemporarySchema(t *testing.T, targetURL string) (*pgxpool.Pool, func()) {
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

	schemaName := fmt.Sprintf("otrt_state_test_%d", time.Now().UnixNano())
	if _, err := baseDB.Exec(`CREATE SCHEMA ` + quoteStateIdentifier(schemaName)); err != nil {
		t.Fatalf("create fallback schema: %v", err)
	}

	query := parsed.Query()
	query.Set("search_path", schemaName+",public")
	parsed.RawQuery = query.Encode()
	cleanup := func() {
		_, _ = baseDB.Exec(`DROP SCHEMA IF EXISTS ` + quoteStateIdentifier(schemaName) + ` CASCADE`)
		_ = baseDB.Close()
	}

	if err := applyStateMigrations(parsed.String()); err != nil {
		cleanup()
		t.Fatalf("apply migrations to fallback schema: %v", err)
	}
	pool, err := pgxpool.New(context.Background(), parsed.String())
	if err != nil {
		cleanup()
		t.Fatalf("create fallback pool: %v", err)
	}
	return pool, func() {
		pool.Close()
		cleanup()
	}
}

func applyStateMigrations(databaseURL string) error {
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
	if err := goose.Up(db, stateMigrationDir()); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

func stateMigrationDir() string {
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

func resetMatcherData(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, `
		TRUNCATE TABLE
			incident,
			vehicle_trip_assignment,
			manual_override,
			telemetry_event,
			gtfs_frequency,
			gtfs_shape_line,
			gtfs_shape_point,
			gtfs_stop_time,
			gtfs_trip,
			gtfs_stop,
			gtfs_route,
			feed_version
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("reset matcher data: %v", err)
	}
	_, err = pool.Exec(ctx, `
		INSERT INTO agency (id, name, timezone, contact_email, public_url)
		VALUES ('demo-agency', 'Demo Agency', 'America/Vancouver', 'dev@example.com', 'http://localhost')
		ON CONFLICT (id) DO UPDATE
		SET timezone = EXCLUDED.timezone
	`)
	if err != nil {
		t.Fatalf("seed agency: %v", err)
	}
}

func seedDemoSchedule(t *testing.T, ctx context.Context, pool *pgxpool.Pool, includeShape bool) {
	t.Helper()
	_, err := pool.Exec(ctx, `
		INSERT INTO feed_version (id, agency_id, source_type, lifecycle_state, is_active, activated_at)
		VALUES ('feed-demo', 'demo-agency', 'seed', 'active', true, now());

		INSERT INTO gtfs_route (id, feed_version_id, agency_id, short_name, long_name, route_type)
		VALUES ('route-10', 'feed-demo', 'demo-agency', '10', 'Route 10', 3);

		INSERT INTO gtfs_stop (id, feed_version_id, agency_id, name, lat, lon, geom)
		VALUES
			('stop-1', 'feed-demo', 'demo-agency', 'Stop 1', 49.2827, -123.1207, ST_SetSRID(ST_MakePoint(-123.1207, 49.2827), 4326)),
			('stop-2', 'feed-demo', 'demo-agency', 'Stop 2', 49.2760, -123.1150, ST_SetSRID(ST_MakePoint(-123.1150, 49.2760), 4326));

		INSERT INTO gtfs_calendar (service_id, feed_version_id, agency_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday, start_date, end_date)
		VALUES ('weekday', 'feed-demo', 'demo-agency', true, true, true, true, true, false, false, '20260401', '20261231');

		INSERT INTO gtfs_trip (id, feed_version_id, agency_id, route_id, service_id, block_id, shape_id, direction_id)
		VALUES ('trip-10-0800', 'feed-demo', 'demo-agency', 'route-10', 'weekday', 'block-10', 'shape-10', 0);

		INSERT INTO gtfs_stop_time (trip_id, feed_version_id, agency_id, arrival_time, departure_time, stop_id, stop_sequence, shape_dist_traveled)
		VALUES
			('trip-10-0800', 'feed-demo', 'demo-agency', '08:00:00', '08:00:00', 'stop-1', 1, 0),
			('trip-10-0800', 'feed-demo', 'demo-agency', '08:10:00', '08:10:00', 'stop-2', 2, 1200)
	`)
	if err != nil {
		t.Fatalf("seed schedule: %v", err)
	}
	if includeShape {
		_, err = pool.Exec(ctx, `
			INSERT INTO gtfs_shape_point (shape_id, feed_version_id, agency_id, lat, lon, sequence, dist_traveled, geom)
			VALUES
				('shape-10', 'feed-demo', 'demo-agency', 49.2827, -123.1207, 1, 0, ST_SetSRID(ST_MakePoint(-123.1207, 49.2827), 4326)),
				('shape-10', 'feed-demo', 'demo-agency', 49.2760, -123.1150, 2, 1200, ST_SetSRID(ST_MakePoint(-123.1150, 49.2760), 4326));
		`)
		if err != nil {
			t.Fatalf("seed shape: %v", err)
		}
	}
}

func storeTelemetry(t *testing.T, ctx context.Context, repo *telemetry.PostgresRepository, event telemetry.Event) telemetry.StoredEvent {
	t.Helper()
	raw, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal telemetry: %v", err)
	}
	result, err := repo.Store(ctx, event, raw)
	if err != nil {
		t.Fatalf("store telemetry: %v", err)
	}
	return result.StoredEvent
}

func quoteStateIdentifier(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
