package prediction

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

func TestPostgresDiagnosticsRepositoryIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("set INTEGRATION_TESTS=1 to run DB-backed diagnostics tests")
	}
	ctx := context.Background()
	pool, cleanup := setupPredictionIntegrationDB(t)
	defer cleanup()

	_, err := pool.Exec(ctx, `
		INSERT INTO agency (id, name, timezone)
		VALUES ('demo-agency', 'Demo Agency', 'America/Vancouver')
	`)
	if err != nil {
		t.Fatalf("seed agency: %v", err)
	}
	repo := NewPostgresDiagnosticsRepository(pool)
	result, err := repo.SaveTripUpdatesDiagnostics(ctx, DiagnosticsRecord{
		AgencyID:            "demo-agency",
		SnapshotAt:          time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC),
		AdapterName:         "noop",
		Status:              StatusNoop,
		Reason:              ReasonNoopAdapter,
		ActiveFeedVersionID: "feed-demo",
		InputCounts: InputCounts{
			TelemetryRows:     2,
			AssignmentRows:    1,
			TripUpdatesOutput: 0,
		},
		VehiclePositionsURL:         "http://localhost:8083/public/gtfsrt/vehicle_positions.pb",
		DiagnosticsPersistenceState: "stored",
	})
	if err != nil {
		t.Fatalf("save diagnostics: %v", err)
	}
	if !result.Stored {
		t.Fatalf("stored = false, want true")
	}

	var detailsBytes []byte
	var feedType string
	var coverage *float64
	if err := pool.QueryRow(ctx, `
		SELECT feed_type, coverage_percent, details_json
		FROM feed_health_snapshot
		WHERE agency_id = 'demo-agency'
	`).Scan(&feedType, &coverage, &detailsBytes); err != nil {
		t.Fatalf("query diagnostics row: %v", err)
	}
	if feedType != "trip_updates" {
		t.Fatalf("feed_type = %q, want trip_updates", feedType)
	}
	if coverage != nil {
		t.Fatalf("coverage = %v, want nil when no eligible coverage denominator is recorded", *coverage)
	}
	var details map[string]any
	if err := json.Unmarshal(detailsBytes, &details); err != nil {
		t.Fatalf("unmarshal details: %v", err)
	}
	assertDetail(t, details, "adapter_name", "noop")
	assertDetail(t, details, "diagnostics_status", StatusNoop)
	assertDetail(t, details, "diagnostics_reason", ReasonNoopAdapter)
	assertDetail(t, details, "active_feed_version_id", "feed-demo")
	assertDetail(t, details, "vehicle_positions_url", "http://localhost:8083/public/gtfsrt/vehicle_positions.pb")
	assertDetail(t, details, "diagnostics_persistence_outcome", "stored")
	counts, ok := details["input_counts"].(map[string]any)
	if !ok {
		t.Fatalf("input_counts = %#v, want object", details["input_counts"])
	}
	if counts["telemetry_rows"].(float64) != 2 || counts["assignment_rows"].(float64) != 1 || counts["trip_updates_output"].(float64) != 0 {
		t.Fatalf("input_counts = %+v, want persisted counts", counts)
	}
}

func TestPostgresOperationsRepositoryIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("set INTEGRATION_TESTS=1 to run DB-backed prediction operations tests")
	}
	ctx := context.Background()
	pool, cleanup := setupPredictionIntegrationDB(t)
	defer cleanup()

	_, err := pool.Exec(ctx, `
		INSERT INTO agency (id, name, timezone)
		VALUES ('demo-agency', 'Demo Agency', 'America/Vancouver')
	`)
	if err != nil {
		t.Fatalf("seed agency: %v", err)
	}

	repo := NewPostgresOperationsRepository(pool)
	now := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	expires := now.Add(time.Hour)
	created, err := repo.CreatePredictionOverride(ctx, OverrideInput{
		AgencyID:     "demo-agency",
		VehicleID:    "bus-10",
		OverrideType: "trip_assignment",
		TripID:       "trip-10",
		StartDate:    "20260421",
		StartTime:    "08:00:00",
		State:        "in_service",
		ExpiresAt:    &expires,
		Reason:       "fix bad match",
		ActorID:      "operator@example.com",
		Now:          now,
	})
	if err != nil {
		t.Fatalf("create override: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("created override id = 0")
	}

	active, err := repo.ListActivePredictionOverrides(ctx, "demo-agency", now)
	if err != nil {
		t.Fatalf("list active overrides: %v", err)
	}
	if len(active) != 1 || active[0].TripID != "trip-10" {
		t.Fatalf("active overrides = %+v, want created override", active)
	}

	replaced, err := repo.ReplacePredictionOverride(ctx, OverrideInput{
		AgencyID:     "demo-agency",
		VehicleID:    "bus-10",
		OverrideType: "service_state",
		State:        "deadhead",
		Reason:       "vehicle deadheading",
		ActorID:      "operator@example.com",
		Now:          now.Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("replace override: %v", err)
	}
	active, err = repo.ListActivePredictionOverrides(ctx, "demo-agency", now.Add(2*time.Minute))
	if err != nil {
		t.Fatalf("list replaced overrides: %v", err)
	}
	if len(active) != 1 || active[0].ID != replaced.ID || active[0].State != "deadhead" {
		t.Fatalf("active overrides after replace = %+v, want replacement only", active)
	}

	if err := repo.ClearPredictionOverride(ctx, "demo-agency", replaced.ID, "operator@example.com", "resume automatic", now.Add(3*time.Minute)); err != nil {
		t.Fatalf("clear override: %v", err)
	}
	active, err = repo.ListActivePredictionOverrides(ctx, "demo-agency", now.Add(4*time.Minute))
	if err != nil {
		t.Fatalf("list cleared overrides: %v", err)
	}
	if len(active) != 0 {
		t.Fatalf("active overrides after clear = %+v, want none", active)
	}

	past := now.Add(-time.Minute)
	if _, err := repo.CreatePredictionOverride(ctx, OverrideInput{
		AgencyID:     "demo-agency",
		VehicleID:    "bus-expired",
		OverrideType: "service_state",
		State:        "layover",
		ExpiresAt:    &past,
		Reason:       "expired",
		ActorID:      "operator@example.com",
		Now:          now,
	}); err != nil {
		t.Fatalf("create expired override: %v", err)
	}
	active, err = repo.ListActivePredictionOverrides(ctx, "demo-agency", now)
	if err != nil {
		t.Fatalf("list after expired override: %v", err)
	}
	if len(active) != 0 {
		t.Fatalf("active overrides with expired row = %+v, want none", active)
	}

	var auditRows int
	if err := pool.QueryRow(ctx, `
		SELECT count(*)
		FROM audit_log
		WHERE agency_id = 'demo-agency'
		  AND action IN ('prediction_override.create', 'prediction_override.replace', 'prediction_override.clear')
	`).Scan(&auditRows); err != nil {
		t.Fatalf("count audit rows: %v", err)
	}
	if auditRows < 4 {
		t.Fatalf("audit rows = %d, want create/replace/clear coverage", auditRows)
	}

	if err := repo.SavePredictionReviewItems(ctx, []ReviewItem{{
		AgencyID:   "demo-agency",
		SnapshotAt: now,
		VehicleID:  "bus-10",
		RouteID:    "route-10",
		TripID:     "trip-10",
		StartDate:  "20260421",
		StartTime:  "08:00:00",
		Severity:   "warning",
		Reason:     ReasonDeadheadNoPrediction,
		Status:     ReviewStatusOpen,
		Details:    map[string]any{"source": "test"},
	}}); err != nil {
		t.Fatalf("save review item: %v", err)
	}
	reviews, err := repo.ListPredictionReviewItems(ctx, ReviewFilter{AgencyID: "demo-agency", Status: ReviewStatusOpen})
	if err != nil {
		t.Fatalf("list review items: %v", err)
	}
	if len(reviews) != 1 || reviews[0].Reason != ReasonDeadheadNoPrediction || reviews[0].Status != ReviewStatusOpen {
		t.Fatalf("reviews = %+v, want open prediction review", reviews)
	}
	if err := repo.UpdatePredictionReviewStatus(ctx, "demo-agency", reviews[0].ID, ReviewStatusDeferred, "operator@example.com", "check later", now.Add(5*time.Minute)); err != nil {
		t.Fatalf("defer review item: %v", err)
	}
	deferred, err := repo.ListPredictionReviewItems(ctx, ReviewFilter{AgencyID: "demo-agency", Status: ReviewStatusDeferred})
	if err != nil {
		t.Fatalf("list deferred review items: %v", err)
	}
	if len(deferred) != 1 || deferred[0].Status != ReviewStatusDeferred {
		t.Fatalf("deferred reviews = %+v, want deferred state", deferred)
	}
}

func assertDetail(t *testing.T, details map[string]any, key string, want string) {
	t.Helper()
	if got, _ := details[key].(string); got != want {
		t.Fatalf("%s = %q, want %q", key, got, want)
	}
}

func setupPredictionIntegrationDB(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	targetURL := os.Getenv("TEST_DATABASE_URL")
	if targetURL == "" {
		targetURL = "postgres://postgres:postgres@localhost:55432/open_transit_rt_test?sslmode=disable"
	}
	pool, cleanup, err := setupPredictionTemporaryDatabase(t, targetURL)
	if err != nil {
		t.Fatalf("setup temporary database: %v", err)
	}
	return pool, cleanup
}

func setupPredictionTemporaryDatabase(t *testing.T, targetURL string) (*pgxpool.Pool, func(), error) {
	t.Helper()
	parsed, err := url.Parse(targetURL)
	if err != nil {
		return nil, nil, err
	}
	admin := *parsed
	admin.Path = "/postgres"
	adminDB, err := sql.Open("pgx", admin.String())
	if err != nil {
		return nil, nil, err
	}
	defer adminDB.Close()

	dbName := fmt.Sprintf("otrt_prediction_test_%d", time.Now().UnixNano())
	if _, err := adminDB.Exec(`CREATE DATABASE ` + quoteIdent(dbName)); err != nil {
		return nil, nil, err
	}
	testURL := *parsed
	testURL.Path = "/" + dbName

	db, err := sql.Open("pgx", testURL.String())
	if err != nil {
		return nil, nil, err
	}
	if err := goose.SetDialect("postgres"); err != nil {
		db.Close()
		return nil, nil, err
	}
	if err := goose.Up(db, predictionMigrationDir()); err != nil {
		db.Close()
		return nil, nil, err
	}
	db.Close()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, testURL.String())
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		pool.Close()
		adminDB, err := sql.Open("pgx", admin.String())
		if err == nil {
			_, _ = adminDB.Exec(`DROP DATABASE IF EXISTS ` + quoteIdent(dbName) + ` WITH (FORCE)`)
			_ = adminDB.Close()
		}
	}
	return pool, cleanup, nil
}

func predictionMigrationDir() string {
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

func quoteIdent(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
