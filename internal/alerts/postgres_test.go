package alerts

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

func TestPostgresAlertsRepositoryIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("set INTEGRATION_TESTS=1 to run DB-backed alerts tests")
	}
	ctx := context.Background()
	pool, cleanup := setupAlertsIntegrationDB(t)
	defer cleanup()

	_, err := pool.Exec(ctx, `
		INSERT INTO agency (id, name, timezone)
		VALUES ('demo-agency', 'Demo Agency', 'America/Vancouver'),
		       ('other-agency', 'Other Agency', 'America/Vancouver');
		INSERT INTO feed_version (id, agency_id, source_type, lifecycle_state, is_active, activated_at)
		VALUES ('feed-demo', 'demo-agency', 'seed', 'active', true, '2026-04-21T12:00:00Z'),
		       ('feed-other', 'other-agency', 'seed', 'active', true, '2026-04-21T12:00:00Z');
	`)
	if err != nil {
		t.Fatalf("seed alerts database: %v", err)
	}
	repo := NewPostgresRepository(pool)
	now := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	alert, err := repo.UpsertAlert(ctx, UpsertInput{
		AgencyID:   "demo-agency",
		AlertKey:   "alert-1",
		Cause:      "other_cause",
		Effect:     "no_service",
		HeaderText: "Stop closed",
		ActorID:    "operator@example.com",
		Publish:    true,
		Now:        now,
		Entities:   []InformedEntity{{AgencyID: "demo-agency", RouteID: "route-10"}},
	})
	if err != nil {
		t.Fatalf("upsert alert: %v", err)
	}
	if alert.ID == 0 || alert.Status != StatusPublished || len(alert.Entities) != 1 {
		t.Fatalf("alert = %+v, want published alert with entity", alert)
	}
	alerts, err := repo.ListAlerts(ctx, ListFilter{AgencyID: "demo-agency", PublishedOnly: true, At: now})
	if err != nil {
		t.Fatalf("list alerts: %v", err)
	}
	if len(alerts) != 1 || alerts[0].AlertKey != "alert-1" {
		t.Fatalf("alerts = %+v, want alert-1", alerts)
	}
	if _, err := repo.UpsertAlert(ctx, UpsertInput{
		AgencyID:   "other-agency",
		AlertKey:   "other-alert",
		HeaderText: "Other agency alert",
		ActorID:    "operator@example.com",
		Publish:    true,
		Now:        now,
	}); err != nil {
		t.Fatalf("upsert other agency alert: %v", err)
	}
	alerts, err = repo.ListAlerts(ctx, ListFilter{AgencyID: "demo-agency", PublishedOnly: true, At: now})
	if err != nil {
		t.Fatalf("list scoped alerts after other agency insert: %v", err)
	}
	if len(alerts) != 1 || alerts[0].AgencyID != "demo-agency" || alerts[0].AlertKey != "alert-1" {
		t.Fatalf("scoped alerts = %+v, want only demo-agency alert", alerts)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO manual_override (
			agency_id, vehicle_id, override_type, route_id, trip_id, start_date,
			start_time, state, expires_at, reason, created_by
		)
		VALUES ('demo-agency', 'bus-10', 'canceled_trip', 'route-10', 'trip-10',
			'20260421', '08:00:00', 'canceled', '2026-04-21T14:00:00Z',
			'operator cancellation', 'operator@example.com');
		INSERT INTO incident (
			agency_id, incident_type, severity, route_id, vehicle_id, trip_id, status, details_json
		)
		VALUES ('demo-agency', 'prediction_review', 'warning', 'route-10', 'bus-10', 'trip-10',
			'open', '{"expected_alert_missing":true,"start_date":"20260421","start_time":"08:00:00"}'::jsonb);
		INSERT INTO manual_override (
			agency_id, vehicle_id, override_type, route_id, trip_id, start_date,
			start_time, state, expires_at, reason, created_by
		)
		VALUES ('other-agency', 'bus-other', 'canceled_trip', 'route-other', 'trip-other',
			'20260421', '08:00:00', 'canceled', '2026-04-21T14:00:00Z',
			'operator cancellation', 'operator@example.com');
		INSERT INTO incident (
			agency_id, incident_type, severity, route_id, vehicle_id, trip_id, status, details_json
		)
		VALUES ('other-agency', 'prediction_review', 'warning', 'route-other', 'bus-other', 'trip-other',
			'open', '{"expected_alert_missing":true,"start_date":"20260421","start_time":"08:00:00"}'::jsonb);
	`)
	if err != nil {
		t.Fatalf("seed cancellation linkage: %v", err)
	}
	result, err := repo.ReconcileCanceledTripAlerts(ctx, "demo-agency", "operator@example.com", now)
	if err != nil {
		t.Fatalf("reconcile canceled trip alerts: %v", err)
	}
	if result.CreatedOrUpdated != 1 || result.LinkedReviews != 1 {
		t.Fatalf("reconcile result = %+v, want one alert and one linked review", result)
	}
	var linkedStatus string
	var linkedAlertID float64
	if err := pool.QueryRow(ctx, `
		SELECT status, (details_json->>'service_alert_id')::float8
		FROM incident
		WHERE agency_id = 'demo-agency'
		  AND incident_type = 'prediction_review'
	`).Scan(&linkedStatus, &linkedAlertID); err != nil {
		t.Fatalf("query linked review: %v", err)
	}
	if linkedStatus != "resolved" || linkedAlertID == 0 {
		t.Fatalf("linked review status/id = %s/%v, want resolved linked alert", linkedStatus, linkedAlertID)
	}
	var otherOpen int
	if err := pool.QueryRow(ctx, `
		SELECT count(*)
		FROM incident
		WHERE agency_id = 'other-agency'
		  AND incident_type = 'prediction_review'
		  AND status = 'open'
	`).Scan(&otherOpen); err != nil {
		t.Fatalf("query other-agency review: %v", err)
	}
	if otherOpen != 1 {
		t.Fatalf("other-agency open reviews = %d, want untouched by demo-agency reconciliation", otherOpen)
	}
}

func setupAlertsIntegrationDB(t *testing.T) (*pgxpool.Pool, func()) {
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

	dbName := fmt.Sprintf("otrt_alerts_test_%d", time.Now().UnixNano())
	if _, err := adminDB.Exec(`CREATE DATABASE ` + quoteIdent(dbName)); err != nil {
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
	if err := goose.Up(db, alertsMigrationDir()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	_ = db.Close()

	pool, err := pgxpool.New(context.Background(), testURL.String())
	if err != nil {
		t.Fatalf("connect temporary database: %v", err)
	}
	cleanup := func() {
		pool.Close()
		adminDB, err := sql.Open("pgx", admin.String())
		if err == nil {
			_, _ = adminDB.Exec(`DROP DATABASE IF EXISTS ` + quoteIdent(dbName) + ` WITH (FORCE)`)
			_ = adminDB.Close()
		}
	}
	return pool, cleanup
}

func alertsMigrationDir() string {
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
