package gtfs

import (
	"archive/zip"
	"bytes"
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

func TestParseGTFSZipValidSmall(t *testing.T) {
	payload := zipFixture(t, "../../testdata/gtfs/valid-small", nil)
	feed, report := parseGTFSZip(payload, "demo-agency")
	if report.hasErrors() {
		t.Fatalf("unexpected validation errors: %+v", report.Errors)
	}
	if len(feed.Routes) != 1 || len(feed.Stops) != 3 || len(feed.Trips) != 1 || len(feed.StopTimes) != 3 {
		t.Fatalf("parsed counts = routes %d stops %d trips %d stop_times %d", len(feed.Routes), len(feed.Stops), len(feed.Trips), len(feed.StopTimes))
	}
	if feed.Trips[0].BlockID != "block-10" {
		t.Fatalf("block_id = %q, want block-10", feed.Trips[0].BlockID)
	}
}

func TestParseGTFSZipPreservesAfterMidnightTimes(t *testing.T) {
	payload := zipFixture(t, "../../testdata/gtfs/after-midnight", nil)
	feed, report := parseGTFSZip(payload, "overnight-agency")
	if report.hasErrors() {
		t.Fatalf("unexpected validation errors: %+v", report.Errors)
	}
	if got := feed.StopTimes[2].ArrivalTime; got != "26:10:00" {
		t.Fatalf("arrival_time = %q, want canonical imported text 26:10:00", got)
	}
	seconds, err := ParseGTFSTime(feed.StopTimes[2].ArrivalTime)
	if err != nil {
		t.Fatalf("parse after-midnight time: %v", err)
	}
	if seconds != 26*3600+10*60 {
		t.Fatalf("seconds = %d, want parsed validation/query seconds", seconds)
	}
}

func TestParseGTFSZipOptionalShapesAndFrequencies(t *testing.T) {
	payload := zipFixture(t, "../../testdata/gtfs/valid-small", map[string]string{
		"shapes.txt":      "",
		"frequencies.txt": "",
	})
	_, report := parseGTFSZip(payload, "demo-agency")
	if report.hasErrors() {
		t.Fatalf("missing optional files should be accepted, got errors: %+v", report.Errors)
	}
}

func TestParseGTFSZipRequiresServiceSource(t *testing.T) {
	payload := zipFixture(t, "../../testdata/gtfs/valid-small", map[string]string{
		"calendar.txt":       "",
		"calendar_dates.txt": "",
	})
	_, report := parseGTFSZip(payload, "demo-agency")
	if !report.hasErrors() || !hasImportCode(report, "missing_service_source") {
		t.Fatalf("report errors = %+v, want missing_service_source", report.Errors)
	}
}

func TestParseGTFSZipRejectsUnusableCalendarService(t *testing.T) {
	payload := zipFixture(t, "../../testdata/gtfs/valid-small", map[string]string{
		"calendar.txt": "service_id,monday,tuesday,wednesday,thursday,friday,saturday,sunday,start_date,end_date\nweekday,0,0,0,0,0,0,0,20260401,20261231\n",
	})
	_, report := parseGTFSZip(payload, "demo-agency")
	if !hasImportCode(report, "unusable_calendar_service") || !hasImportCode(report, "missing_usable_service_source") {
		t.Fatalf("report errors = %+v, want unusable calendar and missing usable service source", report.Errors)
	}
}

func TestParseGTFSZipRejectsCalendarDatesOnlyRemovals(t *testing.T) {
	payload := zipFixture(t, "../../testdata/gtfs/valid-small", map[string]string{
		"calendar.txt":       "",
		"calendar_dates.txt": "service_id,date,exception_type\nweekday,20260420,2\n",
	})
	_, report := parseGTFSZip(payload, "demo-agency")
	if !hasImportCode(report, "missing_usable_service_source") {
		t.Fatalf("report errors = %+v, want missing usable service source for removal-only calendar_dates", report.Errors)
	}
}

func TestParseGTFSZipRejectsUnsupportedRouteType(t *testing.T) {
	payload := zipFixture(t, "../../testdata/gtfs/valid-small", map[string]string{
		"routes.txt": "route_id,agency_id,route_short_name,route_long_name,route_type\nroute-10,demo-agency,10,Downtown/Uptown,9999\n",
	})
	_, report := parseGTFSZip(payload, "demo-agency")
	if !hasImportCode(report, "invalid_route_type") {
		t.Fatalf("report errors = %+v, want invalid_route_type", report.Errors)
	}
}

func TestParseGTFSZipMalformedFixtureFails(t *testing.T) {
	payload := zipFixture(t, "../../testdata/gtfs/malformed", nil)
	_, report := parseGTFSZip(payload, "bad-agency")
	if !report.hasErrors() {
		t.Fatalf("malformed fixture unexpectedly passed")
	}
	if !hasImportCode(report, "missing_required_file") || !hasImportCode(report, "unknown_stop") {
		t.Fatalf("report errors = %+v, want missing trips and unknown stop errors", report.Errors)
	}
}

func TestParseGTFSZipRejectsMultiAgencyConflict(t *testing.T) {
	payload := zipFixture(t, "../../testdata/gtfs/valid-small", map[string]string{
		"routes.txt": "route_id,agency_id,route_short_name,route_long_name,route_type\nroute-10,other-agency,10,Downtown/Uptown,3\n",
	})
	_, report := parseGTFSZip(payload, "demo-agency")
	if !hasImportCode(report, "route_agency_mismatch") {
		t.Fatalf("report errors = %+v, want route_agency_mismatch", report.Errors)
	}
}

func TestImportServiceIntegration(t *testing.T) {
	ctx := context.Background()
	pool, cleanup := setupGTFSImportIntegrationDB(t)
	defer cleanup()

	service := NewImportService(pool)
	service.now = func() time.Time {
		return time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	}

	t.Run("valid import publishes active feed with block and shape line", func(t *testing.T) {
		resetGTFSImportData(t, ctx, pool)
		path := writeZipFixture(t, "../../testdata/gtfs/valid-small", nil)
		result, err := service.ImportZip(ctx, ImportOptions{AgencyID: "demo-agency", ZipPath: path, ActorID: "test"})
		if err != nil {
			t.Fatalf("import zip: %v", err)
		}
		if result.Status != ImportStatusPublished || result.FeedVersionID == "" || !result.ReportStored {
			t.Fatalf("result = %+v, want published stored report", result)
		}

		repo := NewPostgresRepository(pool)
		active, err := repo.ActiveFeedVersion(ctx, "demo-agency")
		if err != nil {
			t.Fatalf("active feed: %v", err)
		}
		if active.ID != result.FeedVersionID {
			t.Fatalf("active feed = %s, want %s", active.ID, result.FeedVersionID)
		}
		candidates, err := repo.ListTripCandidates(ctx, "demo-agency", result.FeedVersionID, "20260420")
		if err != nil {
			t.Fatalf("list candidates: %v", err)
		}
		if len(candidates) != 1 || candidates[0].BlockID != "block-10" {
			t.Fatalf("candidates = %+v, want published block_id visible downstream", candidates)
		}

		var shapeLines int
		if err := pool.QueryRow(ctx, `
			SELECT count(*)
			FROM gtfs_shape_line
			WHERE agency_id = 'demo-agency'
			  AND feed_version_id = $1
			  AND shape_id = 'shape-10'
			  AND ST_NPoints(geom) = 3
		`, result.FeedVersionID).Scan(&shapeLines); err != nil {
			t.Fatalf("query shape line: %v", err)
		}
		if shapeLines != 1 {
			t.Fatalf("shape lines = %d, want one ordered line with 3 points", shapeLines)
		}
	})

	t.Run("repository feed versions imports and draft listings are agency scoped", func(t *testing.T) {
		resetGTFSImportData(t, ctx, pool)
		demoResult, err := service.ImportZip(ctx, ImportOptions{AgencyID: "demo-agency", ZipPath: writeZipFixture(t, "../../testdata/gtfs/valid-small", nil), ActorID: "admin-a@example.com"})
		if err != nil {
			t.Fatalf("import demo-agency: %v", err)
		}
		agencyBZip := writeZipFixture(t, "../../testdata/gtfs/valid-small", map[string]string{
			"agency.txt": "agency_id,agency_name,agency_url,agency_timezone,agency_lang\nagency-b,Agency B,http://agency-b.example,America/Vancouver,en\n",
			"routes.txt": "route_id,agency_id,route_short_name,route_long_name,route_type\nroute-b-20,agency-b,20,Agency B Route,3\n",
			"trips.txt":  "route_id,service_id,trip_id,trip_headsign,direction_id,block_id,shape_id\nroute-b-20,weekday,trip-b-20,Agency B,0,block-b-20,shape-10\n",
			"stop_times.txt": "trip_id,arrival_time,departure_time,stop_id,stop_sequence,shape_dist_traveled\n" +
				"trip-b-20,08:00:00,08:00:00,stop-1,1,0\n" +
				"trip-b-20,08:10:00,08:10:00,stop-2,2,1200\n" +
				"trip-b-20,08:20:00,08:20:00,stop-3,3,2500\n",
		})
		agencyBResult, err := service.ImportZip(ctx, ImportOptions{AgencyID: "agency-b", ZipPath: agencyBZip, ActorID: "admin-b@example.com"})
		if err != nil {
			t.Fatalf("import agency-b: %v", err)
		}

		repo := NewPostgresRepository(pool)
		activeA, err := repo.ActiveFeedVersion(ctx, "demo-agency")
		if err != nil {
			t.Fatalf("active feed demo-agency: %v", err)
		}
		if activeA.ID != demoResult.FeedVersionID || activeA.AgencyID != "demo-agency" {
			t.Fatalf("active demo-agency feed = %+v, want %s", activeA, demoResult.FeedVersionID)
		}
		activeB, err := repo.ActiveFeedVersion(ctx, "agency-b")
		if err != nil {
			t.Fatalf("active feed agency-b: %v", err)
		}
		if activeB.ID != agencyBResult.FeedVersionID || activeB.AgencyID != "agency-b" {
			t.Fatalf("active agency-b feed = %+v, want %s", activeB, agencyBResult.FeedVersionID)
		}

		candidatesA, err := repo.ListTripCandidates(ctx, "demo-agency", demoResult.FeedVersionID, "20260420")
		if err != nil {
			t.Fatalf("list demo-agency candidates: %v", err)
		}
		if len(candidatesA) != 1 || candidatesA[0].AgencyID != "demo-agency" || candidatesA[0].TripID != "trip-10-0800" {
			t.Fatalf("demo-agency candidates = %+v, want only demo trip", candidatesA)
		}
		crossCandidates, err := repo.ListTripCandidates(ctx, "demo-agency", agencyBResult.FeedVersionID, "20260420")
		if err != nil {
			t.Fatalf("list cross-agency candidates: %v", err)
		}
		if len(crossCandidates) != 0 {
			t.Fatalf("cross-agency candidates = %+v, want none", crossCandidates)
		}

		var importAgency string
		if err := pool.QueryRow(ctx, `
			SELECT agency_id
			FROM gtfs_import
			WHERE feed_version_id = $1
		`, demoResult.FeedVersionID).Scan(&importAgency); err != nil {
			t.Fatalf("query demo import agency: %v", err)
		}
		if importAgency != "demo-agency" {
			t.Fatalf("demo import agency = %q, want demo-agency", importAgency)
		}

		drafts := NewDraftService(pool)
		draftA, err := drafts.CreateDraft(ctx, CreateDraftOptions{AgencyID: "demo-agency", Name: "Agency A draft", ActorID: "admin-a@example.com", Blank: true})
		if err != nil {
			t.Fatalf("create demo-agency draft: %v", err)
		}
		draftB, err := drafts.CreateDraft(ctx, CreateDraftOptions{AgencyID: "agency-b", Name: "Agency B draft", ActorID: "admin-b@example.com", Blank: true})
		if err != nil {
			t.Fatalf("create agency-b draft: %v", err)
		}
		draftsA, err := drafts.ListDrafts(ctx, "demo-agency", true)
		if err != nil {
			t.Fatalf("list demo-agency drafts: %v", err)
		}
		if len(draftsA) != 1 || draftsA[0].ID != draftA.ID || draftsA[0].ID == draftB.ID {
			t.Fatalf("demo-agency drafts = %+v, want only %s", draftsA, draftA.ID)
		}
		draftsB, err := drafts.ListDrafts(ctx, "agency-b", true)
		if err != nil {
			t.Fatalf("list agency-b drafts: %v", err)
		}
		if len(draftsB) != 1 || draftsB[0].ID != draftB.ID || draftsB[0].ID == draftA.ID {
			t.Fatalf("agency-b drafts = %+v, want only %s", draftsB, draftB.ID)
		}
	})

	t.Run("failed import stores report and leaves no staged feed version", func(t *testing.T) {
		resetGTFSImportData(t, ctx, pool)
		path := writeZipFixture(t, "../../testdata/gtfs/malformed", nil)
		result, err := service.ImportZip(ctx, ImportOptions{AgencyID: "bad-agency", ZipPath: path, ActorID: "test"})
		if err == nil {
			t.Fatalf("malformed import unexpectedly succeeded")
		}
		if result.Status != ImportStatusFailed || !result.ReportStored || result.FeedVersionID != "" {
			t.Fatalf("result = %+v, want failed stored report without feed version", result)
		}
		var feedVersions int
		if err := pool.QueryRow(ctx, `SELECT count(*) FROM feed_version WHERE agency_id = 'bad-agency'`).Scan(&feedVersions); err != nil {
			t.Fatalf("count feed versions: %v", err)
		}
		if feedVersions != 0 {
			t.Fatalf("feed versions = %d, want no staged or active rows after validation failure", feedVersions)
		}
		var failedImports int
		if err := pool.QueryRow(ctx, `
			SELECT count(*)
			FROM gtfs_import gi
			JOIN validation_report vr ON vr.gtfs_import_id = gi.id
			WHERE gi.agency_id = 'bad-agency'
			  AND gi.status = 'failed'
			  AND gi.feed_version_id IS NULL
			  AND vr.status = 'failed'
		`).Scan(&failedImports); err != nil {
			t.Fatalf("count failed import reports: %v", err)
		}
		if failedImports != 1 {
			t.Fatalf("failed import reports = %d, want 1", failedImports)
		}
	})

	t.Run("publish failure stores gtfs import and validation report outside rolled back transaction", func(t *testing.T) {
		resetGTFSImportData(t, ctx, pool)
		_, err := pool.Exec(ctx, `
			CREATE OR REPLACE FUNCTION fail_gtfs_import_feed_version_insert()
			RETURNS trigger AS $$
			BEGIN
				RAISE EXCEPTION 'forced publish failure';
			END;
			$$ LANGUAGE plpgsql;

			CREATE TRIGGER fail_gtfs_import_feed_version_insert_trigger
			BEFORE INSERT ON feed_version
			FOR EACH ROW EXECUTE FUNCTION fail_gtfs_import_feed_version_insert();
		`)
		if err != nil {
			t.Fatalf("install failure trigger: %v", err)
		}
		t.Cleanup(func() {
			_, _ = pool.Exec(ctx, `
				DROP TRIGGER IF EXISTS fail_gtfs_import_feed_version_insert_trigger ON feed_version;
				DROP FUNCTION IF EXISTS fail_gtfs_import_feed_version_insert();
			`)
		})

		path := writeZipFixture(t, "../../testdata/gtfs/valid-small", nil)
		result, err := service.ImportZip(ctx, ImportOptions{AgencyID: "demo-agency", ZipPath: path, ActorID: "test"})
		if err == nil {
			t.Fatalf("forced publish failure unexpectedly succeeded")
		}
		if result.Status != ImportStatusFailed || !result.ReportStored || result.FeedVersionID != "" {
			t.Fatalf("result = %+v, want failed stored publish-failure report without feed version", result)
		}

		var feedVersions int
		if err := pool.QueryRow(ctx, `SELECT count(*) FROM feed_version WHERE agency_id = 'demo-agency'`).Scan(&feedVersions); err != nil {
			t.Fatalf("count feed versions: %v", err)
		}
		if feedVersions != 0 {
			t.Fatalf("feed versions = %d, want rolled back publish to leave none", feedVersions)
		}

		var failedReports int
		if err := pool.QueryRow(ctx, `
			SELECT count(*)
			FROM gtfs_import gi
			JOIN validation_report vr ON vr.gtfs_import_id = gi.id
			WHERE gi.agency_id = 'demo-agency'
			  AND gi.status = 'failed'
			  AND gi.feed_version_id IS NULL
			  AND vr.feed_version_id IS NULL
			  AND vr.status = 'failed'
			  AND vr.report_json::text LIKE '%publish_failed%'
		`).Scan(&failedReports); err != nil {
			t.Fatalf("count failed publish reports: %v", err)
		}
		if failedReports != 1 {
			t.Fatalf("failed publish reports = %d, want 1 gtfs_import plus validation_report", failedReports)
		}
	})

	t.Run("second valid import switches active feed", func(t *testing.T) {
		resetGTFSImportData(t, ctx, pool)
		first, err := service.ImportZip(ctx, ImportOptions{AgencyID: "demo-agency", ZipPath: writeZipFixture(t, "../../testdata/gtfs/valid-small", nil), ActorID: "test"})
		if err != nil {
			t.Fatalf("first import: %v", err)
		}
		second, err := service.ImportZip(ctx, ImportOptions{AgencyID: "demo-agency", ZipPath: writeZipFixture(t, "../../testdata/gtfs/valid-small", map[string]string{
			"trips.txt":      "route_id,service_id,trip_id,trip_headsign,direction_id,block_id,shape_id\nroute-10,weekday,trip-10-0900,Uptown,0,block-20,shape-10\n",
			"stop_times.txt": "trip_id,arrival_time,departure_time,stop_id,stop_sequence,shape_dist_traveled\ntrip-10-0900,09:00:00,09:00:00,stop-1,1,0\ntrip-10-0900,09:10:00,09:10:00,stop-2,2,1200\ntrip-10-0900,09:20:00,09:20:00,stop-3,3,2500\n",
		}), ActorID: "test"})
		if err != nil {
			t.Fatalf("second import: %v", err)
		}

		var firstActive, firstState, secondActive, secondState string
		if err := pool.QueryRow(ctx, `
			SELECT lifecycle_state, is_active::text
			FROM feed_version
			WHERE id = $1
		`, first.FeedVersionID).Scan(&firstState, &firstActive); err != nil {
			t.Fatalf("query first feed: %v", err)
		}
		if err := pool.QueryRow(ctx, `
			SELECT lifecycle_state, is_active::text
			FROM feed_version
			WHERE id = $1
		`, second.FeedVersionID).Scan(&secondState, &secondActive); err != nil {
			t.Fatalf("query second feed: %v", err)
		}
		if firstState != "retired" || firstActive != "false" || secondState != "active" || secondActive != "true" {
			t.Fatalf("feed states first=%s/%s second=%s/%s, want retired/false and active/true", firstState, firstActive, secondState, secondActive)
		}
	})
}

func hasImportCode(report ImportReport, code string) bool {
	for _, msg := range report.Errors {
		if msg.Code == code {
			return true
		}
	}
	return false
}

func zipFixture(t *testing.T, dir string, overrides map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read fixture dir: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if value, ok := overrides[name]; ok && value == "" {
			continue
		}
		var content []byte
		if value, ok := overrides[name]; ok {
			content = []byte(value)
		} else {
			content, err = os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				t.Fatalf("read fixture file %s: %v", name, err)
			}
		}
		file, err := writer.Create(name)
		if err != nil {
			t.Fatalf("create zip entry %s: %v", name, err)
		}
		if _, err := file.Write(content); err != nil {
			t.Fatalf("write zip entry %s: %v", name, err)
		}
	}
	for name, value := range overrides {
		if value == "" || fixtureHasFile(t, dir, name) {
			continue
		}
		file, err := writer.Create(name)
		if err != nil {
			t.Fatalf("create override zip entry %s: %v", name, err)
		}
		if _, err := file.Write([]byte(value)); err != nil {
			t.Fatalf("write override zip entry %s: %v", name, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close fixture zip: %v", err)
	}
	return buf.Bytes()
}

func fixtureHasFile(t *testing.T, dir string, name string) bool {
	t.Helper()
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}

func writeZipFixture(t *testing.T, dir string, overrides map[string]string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "gtfs.zip")
	if err := os.WriteFile(path, zipFixture(t, dir, overrides), 0o600); err != nil {
		t.Fatalf("write fixture zip: %v", err)
	}
	return path
}

func setupGTFSImportIntegrationDB(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("set INTEGRATION_TESTS=1 to run DB-backed GTFS import tests")
	}

	targetURL := os.Getenv("TEST_DATABASE_URL")
	if targetURL == "" {
		targetURL = "postgres://postgres:postgres@localhost:55432/open_transit_rt_test?sslmode=disable"
	}

	if pool, cleanup, err := setupGTFSTemporaryDatabase(t, targetURL); err == nil {
		return pool, cleanup
	} else {
		t.Logf("isolated database setup unavailable, falling back to schema mode: %v", err)
	}
	return setupGTFSTemporarySchema(t, targetURL)
}

func setupGTFSTemporaryDatabase(t *testing.T, targetURL string) (*pgxpool.Pool, func(), error) {
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

	dbName := fmt.Sprintf("otrt_gtfs_import_test_%d", time.Now().UnixNano())
	if _, err := adminDB.Exec(`CREATE DATABASE ` + quoteGTFSIdentifier(dbName)); err != nil {
		return nil, nil, fmt.Errorf("create temporary database: %w", err)
	}
	cleanup := func() {
		adminDB, err := sql.Open("pgx", adminURL.String())
		if err == nil {
			defer adminDB.Close()
			_, _ = adminDB.Exec(`DROP DATABASE IF EXISTS ` + quoteGTFSIdentifier(dbName) + ` WITH (FORCE)`)
		}
	}

	testURL := *parsed
	testURL.Path = "/" + dbName
	if err := applyGTFSMigrations(testURL.String()); err != nil {
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

func setupGTFSTemporarySchema(t *testing.T, targetURL string) (*pgxpool.Pool, func()) {
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

	schemaName := fmt.Sprintf("otrt_gtfs_import_test_%d", time.Now().UnixNano())
	if _, err := baseDB.Exec(`CREATE SCHEMA ` + quoteGTFSIdentifier(schemaName)); err != nil {
		t.Fatalf("create fallback schema: %v", err)
	}

	query := parsed.Query()
	query.Set("search_path", schemaName+",public")
	parsed.RawQuery = query.Encode()
	cleanup := func() {
		_, _ = baseDB.Exec(`DROP SCHEMA IF EXISTS ` + quoteGTFSIdentifier(schemaName) + ` CASCADE`)
		_ = baseDB.Close()
	}

	if err := applyGTFSMigrations(parsed.String()); err != nil {
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

func applyGTFSMigrations(databaseURL string) error {
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
	if err := goose.Up(db, gtfsMigrationDir()); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

func gtfsMigrationDir() string {
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

func resetGTFSImportData(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, `
		TRUNCATE TABLE
			audit_log,
			validation_report,
			gtfs_draft_publish,
			gtfs_draft_frequency,
			gtfs_draft_shape_point,
			gtfs_draft_stop_time,
			gtfs_draft_trip,
			gtfs_draft_calendar_date,
			gtfs_draft_calendar,
			gtfs_draft_stop,
			gtfs_draft_route,
			gtfs_draft_agency,
			gtfs_draft_record,
			gtfs_draft,
			gtfs_import,
			gtfs_frequency,
			gtfs_shape_line,
			gtfs_shape_point,
			gtfs_stop_time,
			gtfs_trip,
			gtfs_stop,
			gtfs_route,
			published_feed,
			feed_version,
			agency
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("reset gtfs import data: %v", err)
	}
}

func quoteGTFSIdentifier(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
