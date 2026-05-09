package feed

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

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func TestVehiclePositionsHealthRepositoryIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("set INTEGRATION_TESTS=1 to run DB-backed Vehicle Positions health tests")
	}
	ctx := context.Background()
	pool, cleanup := setupFeedIntegrationDB(t)
	defer cleanup()

	_, err := pool.Exec(ctx, `
		INSERT INTO agency (id, name, timezone)
		VALUES ('demo-agency', 'Demo Agency', 'America/Los_Angeles'),
		       ('other-agency', 'Other Agency', 'America/Los_Angeles');
	`)
	if err != nil {
		t.Fatalf("seed agencies: %v", err)
	}
	freshness := 25.5
	latency := 42.0
	matched := 66.7
	repo := NewVehiclePositionsHealthRepository(pool)
	if err := repo.SaveVehiclePositionsHealth(ctx, VehiclePositionsHealthRecord{
		AgencyID:                "demo-agency",
		SnapshotAt:              time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
		EndpointAvailable:       true,
		FreshnessSeconds:        &freshness,
		GenerationLatencyMS:     &latency,
		MatchedVehiclePercent:   &matched,
		VehiclesInSnapshot:      3,
		VehiclesInProtobuf:      2,
		TripDescriptors:         1,
		StaleVehicles:           1,
		UnmatchedVehicles:       1,
		Truncated:               true,
		VehicleLimit:            2000,
		LatestTelemetryRowsRead: 4,
	}); err != nil {
		t.Fatalf("save vehicle positions health: %v", err)
	}

	var agencyID, feedType string
	var snapshotAt time.Time
	var endpoint bool
	var gotFreshness, gotLatency, gotMatched float64
	var detailsBytes []byte
	if err := pool.QueryRow(ctx, `
		SELECT agency_id, feed_type, snapshot_at, endpoint_available,
		       freshness_seconds, generation_latency_ms, matched_vehicle_percent,
		       details_json
		FROM feed_health_snapshot
		WHERE agency_id = 'demo-agency' AND feed_type = 'vehicle_positions'
	`).Scan(&agencyID, &feedType, &snapshotAt, &endpoint, &gotFreshness, &gotLatency, &gotMatched, &detailsBytes); err != nil {
		t.Fatalf("query vehicle positions health snapshot: %v", err)
	}
	if agencyID != "demo-agency" || feedType != "vehicle_positions" || !endpoint || !snapshotAt.Equal(time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)) {
		t.Fatalf("row identity = %s/%s/%v/%v, want demo vehicle_positions snapshot", agencyID, feedType, endpoint, snapshotAt)
	}
	if gotFreshness != freshness || gotLatency != latency || gotMatched != matched {
		t.Fatalf("metrics = freshness %.1f latency %.1f matched %.1f, want %.1f %.1f %.1f", gotFreshness, gotLatency, gotMatched, freshness, latency, matched)
	}
	var details map[string]any
	if err := json.Unmarshal(detailsBytes, &details); err != nil {
		t.Fatalf("unmarshal details: %v", err)
	}
	wantNumbers := map[string]float64{
		"vehicles_in_snapshot":       3,
		"vehicles_in_protobuf":       2,
		"trip_descriptors":           1,
		"stale_vehicles":             1,
		"unmatched_vehicles":         1,
		"vehicle_limit":              2000,
		"latest_telemetry_rows_read": 4,
	}
	for key, want := range wantNumbers {
		if details[key] != want {
			t.Fatalf("details[%s] = %+v, want %.0f in %+v", key, details[key], want, details)
		}
	}
	if details["truncated"] != true || details["diagnostics_status"] != "observed" {
		t.Fatalf("details = %+v, want truncated true and observed status", details)
	}
	body := string(detailsBytes)
	for _, forbidden := range []string{"bus-1", "vehicle_id", "payload", "token", "Authorization", "Cookie", "postgres://"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("details leaked %q: %s", forbidden, body)
		}
	}
}

func setupFeedIntegrationDB(t *testing.T) (*pgxpool.Pool, func()) {
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
	dbName := fmt.Sprintf("otrt_feed_test_%d", time.Now().UnixNano())
	if _, err := adminDB.Exec(`CREATE DATABASE ` + quoteFeedIdent(dbName)); err != nil {
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
	if err := goose.Up(db, feedMigrationDir()); err != nil {
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
			_, _ = adminDB.Exec(`DROP DATABASE IF EXISTS ` + quoteFeedIdent(dbName) + ` WITH (FORCE)`)
			_ = adminDB.Close()
		}
	}
}

func feedMigrationDir() string {
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

func quoteFeedIdent(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
