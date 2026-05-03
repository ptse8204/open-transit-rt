package compliance

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

func TestReadinessRequiresHTTPSOnlyInProductionDiscoverability(t *testing.T) {
	cfg := feedConfig{
		PublicBaseURL:          "http://localhost:8080",
		FeedBaseURL:            "http://localhost:8080/public",
		TechnicalContactEmail:  "ops@example.org",
		LicenseName:            "CC BY 4.0",
		LicenseURL:             "https://creativecommons.org/licenses/by/4.0/",
		PublicationEnvironment: EnvironmentDev,
	}
	feeds := []FeedMetadata{
		{FeedType: "schedule", CanonicalPublicURL: "http://localhost/public/gtfs/schedule.zip", LicenseName: cfg.LicenseName, LicenseURL: cfg.LicenseURL, ContactEmail: cfg.TechnicalContactEmail, LastValidationStatus: "passed"},
		{FeedType: "vehicle_positions", CanonicalPublicURL: "http://localhost/public/gtfsrt/vehicle_positions.pb", LicenseName: cfg.LicenseName, LicenseURL: cfg.LicenseURL, ContactEmail: cfg.TechnicalContactEmail, LastValidationStatus: "passed"},
		{FeedType: "trip_updates", CanonicalPublicURL: "http://localhost/public/gtfsrt/trip_updates.pb", LicenseName: cfg.LicenseName, LicenseURL: cfg.LicenseURL, ContactEmail: cfg.TechnicalContactEmail, LastValidationStatus: "passed"},
		{FeedType: "alerts", CanonicalPublicURL: "http://localhost/public/gtfsrt/alerts.pb", LicenseName: cfg.LicenseName, LicenseURL: cfg.LicenseURL, ContactEmail: cfg.TechnicalContactEmail, LastValidationStatus: "passed"},
	}
	dev := evaluateReadiness(cfg, feeds)
	if !dev.Discoverable || dev.HTTPSURLs {
		t.Fatalf("dev readiness = %+v, want discoverable with non-HTTPS flagged", dev)
	}
	cfg.PublicationEnvironment = EnvironmentProduction
	production := evaluateReadiness(cfg, feeds)
	if production.Discoverable || production.HTTPSURLs {
		t.Fatalf("production readiness = %+v, want non-HTTPS to block discoverability", production)
	}
}

func TestValidationScoreTreatsMissingValidatorsAsProductionRed(t *testing.T) {
	feeds := []FeedMetadata{{FeedType: "schedule", LastValidationStatus: "not_run"}}
	if got := validationScore(EnvironmentDev, feeds); got != StatusYellow {
		t.Fatalf("dev missing validator score = %s, want yellow", got)
	}
	if got := validationScore(EnvironmentProduction, feeds); got != StatusRed {
		t.Fatalf("production missing validator score = %s, want red", got)
	}
}

func TestValidationScoreUsesPassedWarningAndFailedStatuses(t *testing.T) {
	if got := validationScore(EnvironmentProduction, []FeedMetadata{{FeedType: "schedule", LastValidationStatus: "passed"}}); got != StatusGreen {
		t.Fatalf("passed validator score = %s, want green", got)
	}
	if got := validationScore(EnvironmentProduction, []FeedMetadata{{FeedType: "schedule", LastValidationStatus: "warning"}}); got != StatusYellow {
		t.Fatalf("warning validator score = %s, want yellow", got)
	}
	if got := validationScore(EnvironmentProduction, []FeedMetadata{{FeedType: "schedule", LastValidationStatus: "failed"}}); got != StatusRed {
		t.Fatalf("failed validator score = %s, want red", got)
	}
}

func TestRunValidationStoresNotRunWhenBinaryMissing(t *testing.T) {
	store := &fakeValidationStore{}
	registry := ValidatorRegistry{"static-test": {ID: "static-test", Name: "test-validator", FeedTypes: []string{"schedule"}, RequiresSchedule: true}}
	result, err := RunValidation(context.Background(), store, registry, ValidationRunInput{AgencyID: "demo-agency", FeedType: "schedule", ValidatorID: "static-test", ScheduleZIPPayload: []byte("zip")})
	if err != nil {
		t.Fatalf("run validation: %v", err)
	}
	if result.Status != "not_run" || store.result.Status != "not_run" {
		t.Fatalf("result = %+v stored = %+v, want not_run", result, store.result)
	}
	if store.result.Report["reason"] != "validator_binary_missing" {
		t.Fatalf("report = %+v, want missing binary reason", store.result.Report)
	}
}

func TestRunValidationNormalizesPassedJSONReport(t *testing.T) {
	store := &fakeValidationStore{}
	registry := ValidatorRegistry{"echo-json": {ID: "echo-json", Name: "test-validator", FeedTypes: []string{"schedule"}, Binary: "/bin/echo", Args: []string{`{"status":"passed","error_count":0,"warning_count":0,"info_count":3}`, "{output_dir}"}}}
	result, err := RunValidation(context.Background(), store, registry, ValidationRunInput{AgencyID: "demo-agency", FeedType: "schedule", ValidatorID: "echo-json"})
	if err != nil {
		t.Fatalf("run validation: %v", err)
	}
	if result.Status != "passed" || result.ErrorCount != 0 || result.WarningCount != 0 || result.InfoCount != 3 {
		t.Fatalf("result = %+v, want passed with info count", result)
	}
	if store.result.Report["raw_report"] == nil {
		t.Fatalf("report = %+v, want raw_report", store.result.Report)
	}
}

func TestRunValidationNormalizesWarningJSONReport(t *testing.T) {
	store := &fakeValidationStore{}
	registry := ValidatorRegistry{"echo-json": {ID: "echo-json", Name: "test-validator", FeedTypes: []string{"alerts"}, Binary: "/bin/echo", Args: []string{`{"notices":[{"severity":"WARNING"},{"severity":"INFO"}]}`, "{output_dir}"}}}
	result, err := RunValidation(context.Background(), store, registry, ValidationRunInput{AgencyID: "demo-agency", FeedType: "alerts", ValidatorID: "echo-json"})
	if err != nil {
		t.Fatalf("run validation: %v", err)
	}
	if result.Status != "warning" || result.ErrorCount != 0 || result.WarningCount != 1 || result.InfoCount != 1 {
		t.Fatalf("result = %+v, want warning with notice counts", result)
	}
}

func TestRunValidationNormalizesFailedJSONReport(t *testing.T) {
	store := &fakeValidationStore{}
	script := filepath.Join(t.TempDir(), "validator.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nprintf '%s' '{\"summary\":{\"errors\":2,\"warnings\":1,\"infos\":4}}'\nexit 1\n"), 0o700); err != nil {
		t.Fatalf("write validator script: %v", err)
	}
	registry := ValidatorRegistry{"script-json": {ID: "script-json", Name: "test-validator", FeedTypes: []string{"trip_updates"}, Binary: script, Args: []string{"{output_dir}"}}}
	result, err := RunValidation(context.Background(), store, registry, ValidationRunInput{AgencyID: "demo-agency", FeedType: "trip_updates", ValidatorID: "script-json"})
	if err != nil {
		t.Fatalf("run validation: %v", err)
	}
	if result.Status != "failed" || result.ErrorCount != 2 || result.WarningCount != 1 || result.InfoCount != 4 {
		t.Fatalf("result = %+v, want failed with parsed counts", result)
	}
	if store.result.Report["error"] == "" {
		t.Fatalf("report = %+v, want command error retained", store.result.Report)
	}
}

func TestRunValidationRejectsUnknownValidatorID(t *testing.T) {
	if _, err := RunValidation(context.Background(), &fakeValidationStore{}, ValidatorRegistry{}, ValidationRunInput{AgencyID: "demo-agency", FeedType: "schedule", ValidatorID: "command"}); err == nil {
		t.Fatalf("run validation succeeded with unknown validator_id")
	}
}

func TestPostgresComplianceRecordsAreAgencyScoped(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("set INTEGRATION_TESTS=1 to run DB-backed compliance isolation tests")
	}
	ctx := context.Background()
	pool, cleanup := setupComplianceIntegrationDB(t)
	defer cleanup()

	_, err := pool.Exec(ctx, `
		INSERT INTO agency (id, name, timezone)
		VALUES ('agency-a', 'Agency A', 'America/Los_Angeles'),
		       ('agency-b', 'Agency B', 'America/Los_Angeles');
		INSERT INTO feed_version (id, agency_id, source_type, lifecycle_state, is_active, activated_at)
		VALUES ('feed-a', 'agency-a', 'seed', 'active', true, '2026-05-02T12:00:00Z'),
		       ('feed-b', 'agency-b', 'seed', 'active', true, '2026-05-02T12:00:00Z');
	`)
	if err != nil {
		t.Fatalf("seed agencies: %v", err)
	}
	repo := NewPostgresRepository(pool)
	for _, input := range []BootstrapInput{
		{AgencyID: "agency-a", PublicBaseURL: "https://agency-a.example", FeedBaseURL: "https://agency-a.example/public", TechnicalContactEmail: "ops-a@example.com", LicenseName: "CC BY", LicenseURL: "https://agency-a.example/license", PublicationEnvironment: EnvironmentDev, ActorID: "admin-a@example.com"},
		{AgencyID: "agency-b", PublicBaseURL: "https://agency-b.example", FeedBaseURL: "https://agency-b.example/public", TechnicalContactEmail: "ops-b@example.com", LicenseName: "CC BY", LicenseURL: "https://agency-b.example/license", PublicationEnvironment: EnvironmentDev, ActorID: "admin-b@example.com"},
	} {
		if err := repo.BootstrapPublication(ctx, input); err != nil {
			t.Fatalf("bootstrap %s: %v", input.AgencyID, err)
		}
	}
	if err := repo.StoreValidationResult(ctx, ValidationResult{AgencyID: "agency-a", FeedVersionID: "feed-a", FeedType: "schedule", ValidatorName: "test-validator", ValidatorVersion: "v1", Status: "passed", Report: map[string]any{"redacted": true}}); err != nil {
		t.Fatalf("store validation agency-a: %v", err)
	}
	if _, err := repo.UpsertConsumer(ctx, ConsumerInput{AgencyID: "agency-a", ConsumerName: "Consumer A", Status: "not_started", Packet: map[string]any{"packet": "a"}}); err != nil {
		t.Fatalf("upsert consumer agency-a: %v", err)
	}
	if _, err := repo.UpsertConsumer(ctx, ConsumerInput{AgencyID: "agency-b", ConsumerName: "Consumer B", Status: "not_started", Packet: map[string]any{"packet": "b"}}); err != nil {
		t.Fatalf("upsert consumer agency-b: %v", err)
	}

	discoveryA, err := repo.FeedDiscovery(ctx, "agency-a", time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("feed discovery agency-a: %v", err)
	}
	discoveryB, err := repo.FeedDiscovery(ctx, "agency-b", time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("feed discovery agency-b: %v", err)
	}
	if discoveryA.AgencyID != "agency-a" || strings.Contains(fmt.Sprintf("%+v", discoveryA), "agency-b.example") {
		t.Fatalf("discovery A leaked B data: %+v", discoveryA)
	}
	if discoveryB.AgencyID != "agency-b" || strings.Contains(fmt.Sprintf("%+v", discoveryB), "agency-a.example") {
		t.Fatalf("discovery B leaked A data: %+v", discoveryB)
	}
	if status := feedStatusFor(discoveryA, "schedule"); status != "passed" {
		t.Fatalf("agency-a schedule validation status = %q, want passed", status)
	}
	if status := feedStatusFor(discoveryB, "schedule"); status != "not_run" {
		t.Fatalf("agency-b schedule validation status = %q, want not_run", status)
	}

	consumersA, err := repo.ListConsumers(ctx, "agency-a")
	if err != nil {
		t.Fatalf("list consumers agency-a: %v", err)
	}
	if !containsConsumer(consumersA, "Consumer A") || containsConsumer(consumersA, "Consumer B") {
		t.Fatalf("agency-a consumers = %+v, want A only", consumersA)
	}

	scoreA, err := repo.BuildAndStoreScorecard(ctx, "agency-a", time.Date(2026, 5, 2, 12, 1, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build scorecard agency-a: %v", err)
	}
	scoreB, err := repo.BuildAndStoreScorecard(ctx, "agency-b", time.Date(2026, 5, 2, 12, 2, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build scorecard agency-b: %v", err)
	}
	if scoreA.AgencyID != "agency-a" || scoreB.AgencyID != "agency-b" {
		t.Fatalf("scorecard agencies = %s/%s, want scoped scorecards", scoreA.AgencyID, scoreB.AgencyID)
	}
	latestA, err := repo.LatestScorecard(ctx, "agency-a")
	if err != nil {
		t.Fatalf("latest scorecard agency-a: %v", err)
	}
	if latestA.AgencyID != "agency-a" {
		t.Fatalf("latest scorecard agency = %q, want agency-a", latestA.AgencyID)
	}

	var auditAgencies []string
	rows, err := pool.Query(ctx, `SELECT agency_id FROM audit_log WHERE action = 'publication.bootstrap' ORDER BY agency_id`)
	if err != nil {
		t.Fatalf("query audit agencies: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var agencyID string
		if err := rows.Scan(&agencyID); err != nil {
			t.Fatalf("scan audit agency: %v", err)
		}
		auditAgencies = append(auditAgencies, agencyID)
	}
	if strings.Join(auditAgencies, ",") != "agency-a,agency-b" {
		t.Fatalf("audit agencies = %+v, want written agency ids for both synthetic agencies", auditAgencies)
	}
}

type fakeValidationStore struct {
	result ValidationResult
}

func (f *fakeValidationStore) StoreValidationResult(_ context.Context, result ValidationResult) error {
	f.result = result
	return nil
}

func feedStatusFor(discovery FeedDiscovery, feedType string) string {
	for _, feed := range discovery.Feeds {
		if feed.FeedType == feedType {
			return feed.LastValidationStatus
		}
	}
	return ""
}

func containsConsumer(records []ConsumerRecord, name string) bool {
	for _, record := range records {
		if record.ConsumerName == name {
			return true
		}
	}
	return false
}

func setupComplianceIntegrationDB(t *testing.T) (*pgxpool.Pool, func()) {
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
	dbName := fmt.Sprintf("otrt_compliance_test_%d", time.Now().UnixNano())
	if _, err := adminDB.Exec(`CREATE DATABASE ` + quoteComplianceIdent(dbName)); err != nil {
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
	if err := goose.Up(db, complianceMigrationDir()); err != nil {
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
			_, _ = adminDB.Exec(`DROP DATABASE IF EXISTS ` + quoteComplianceIdent(dbName) + ` WITH (FORCE)`)
			_ = adminDB.Close()
		}
	}
}

func complianceMigrationDir() string {
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

func quoteComplianceIdent(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
