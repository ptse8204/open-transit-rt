package gtfs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestDraftServiceIntegration(t *testing.T) {
	ctx := context.Background()
	pool, cleanup := setupGTFSImportIntegrationDB(t)
	defer cleanup()

	imports := NewImportService(pool)
	imports.now = func() time.Time {
		return time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	}
	drafts := NewDraftService(pool)
	drafts.now = imports.now

	t.Run("blank draft without active feed and all core entity CRUD", func(t *testing.T) {
		resetGTFSImportData(t, ctx, pool)
		seedDraftAgency(t, ctx, pool, "demo-agency")

		draft, err := drafts.CreateDraft(ctx, CreateDraftOptions{AgencyID: "demo-agency", Name: "Blank draft", ActorID: "tester"})
		if err != nil {
			t.Fatalf("create blank draft without active feed: %v", err)
		}
		if draft.BaseFeedVersionID != "" {
			t.Fatalf("base feed = %q, want blank draft with no base feed", draft.BaseFeedVersionID)
		}

		if err := drafts.UpsertAgency(ctx, DraftAgency{DraftID: draft.ID, AgencyID: "demo-agency", Name: "Draft Demo", Timezone: "America/Vancouver", ContactEmail: "draft@example.com", PublicURL: "https://example.test"}); err != nil {
			t.Fatalf("upsert draft agency: %v", err)
		}
		if err := drafts.UpsertRoute(ctx, DraftRoute{DraftID: draft.ID, ID: "route-1", ShortName: "1", LongName: "Route One", RouteType: 3}); err != nil {
			t.Fatalf("upsert route: %v", err)
		}
		if err := drafts.UpsertStop(ctx, DraftStop{DraftID: draft.ID, ID: "stop-1", Name: "First", Lat: 49.2, Lon: -123.1}); err != nil {
			t.Fatalf("upsert stop 1: %v", err)
		}
		if err := drafts.UpsertStop(ctx, DraftStop{DraftID: draft.ID, ID: "stop-2", Name: "Second", Lat: 49.21, Lon: -123.11}); err != nil {
			t.Fatalf("upsert stop 2: %v", err)
		}
		if err := drafts.UpsertCalendar(ctx, DraftCalendar{DraftID: draft.ID, ServiceID: "weekday", Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, StartDate: "20260401", EndDate: "20261231"}); err != nil {
			t.Fatalf("upsert calendar: %v", err)
		}
		if err := drafts.UpsertCalendarDate(ctx, DraftCalendarDate{DraftID: draft.ID, ServiceID: "weekday", Date: "20260420", ExceptionType: 1}); err != nil {
			t.Fatalf("upsert calendar date: %v", err)
		}
		direction := 0
		if err := drafts.UpsertTrip(ctx, DraftTrip{DraftID: draft.ID, ID: "trip-1", RouteID: "route-1", ServiceID: "weekday", BlockID: "block-1", ShapeID: "shape-1", DirectionID: &direction}); err != nil {
			t.Fatalf("upsert trip: %v", err)
		}
		zero := 0.0
		if err := drafts.UpsertStopTime(ctx, DraftStopTime{DraftID: draft.ID, TripID: "trip-1", ArrivalTime: "08:00:00", DepartureTime: "08:00:00", StopID: "stop-1", StopSequence: 1, ShapeDistTraveled: &zero}); err != nil {
			t.Fatalf("upsert stop time 1: %v", err)
		}
		if err := drafts.UpsertStopTime(ctx, DraftStopTime{DraftID: draft.ID, TripID: "trip-1", ArrivalTime: "08:10:00", DepartureTime: "08:10:00", StopID: "stop-2", StopSequence: 2}); err != nil {
			t.Fatalf("upsert stop time 2: %v", err)
		}
		if err := drafts.UpsertShapePoint(ctx, DraftShapePoint{DraftID: draft.ID, ShapeID: "shape-1", Lat: 49.2, Lon: -123.1, Sequence: 1, DistTraveled: &zero}); err != nil {
			t.Fatalf("upsert shape point 1: %v", err)
		}
		dist := 1000.0
		if err := drafts.UpsertShapePoint(ctx, DraftShapePoint{DraftID: draft.ID, ShapeID: "shape-1", Lat: 49.21, Lon: -123.11, Sequence: 2, DistTraveled: &dist}); err != nil {
			t.Fatalf("upsert shape point 2: %v", err)
		}
		if err := drafts.UpsertFrequency(ctx, DraftFrequency{DraftID: draft.ID, TripID: "trip-1", StartTime: "08:00:00", EndTime: "09:00:00", HeadwaySecs: 600, ExactTimes: 1}); err != nil {
			t.Fatalf("upsert frequency: %v", err)
		}

		assertLen(t, "routes", mustList(drafts.ListRoutes(ctx, draft.ID)), 1)
		assertLen(t, "stops", mustList(drafts.ListStops(ctx, draft.ID)), 2)
		assertLen(t, "trips", mustList(drafts.ListTrips(ctx, draft.ID)), 1)
		assertLen(t, "stop_times", mustList(drafts.ListStopTimes(ctx, draft.ID)), 2)
		assertLen(t, "calendars", mustList(drafts.ListCalendars(ctx, draft.ID)), 1)
		assertLen(t, "calendar_dates", mustList(drafts.ListCalendarDates(ctx, draft.ID)), 1)
		assertLen(t, "shape_points", mustList(drafts.ListShapePoints(ctx, draft.ID)), 2)
		assertLen(t, "frequencies", mustList(drafts.ListFrequencies(ctx, draft.ID)), 1)

		if err := drafts.RemoveFrequency(ctx, draft.ID, "trip-1", "08:00:00"); err != nil {
			t.Fatalf("remove frequency: %v", err)
		}
		if err := drafts.RemoveShapePoint(ctx, draft.ID, "shape-1", 2); err != nil {
			t.Fatalf("remove shape point: %v", err)
		}
		if err := drafts.RemoveStopTime(ctx, draft.ID, "trip-1", 2); err != nil {
			t.Fatalf("remove stop time: %v", err)
		}
		if err := drafts.RemoveTrip(ctx, draft.ID, "trip-1"); err != nil {
			t.Fatalf("remove trip: %v", err)
		}
		if err := drafts.RemoveCalendarDate(ctx, draft.ID, "weekday", "20260420"); err != nil {
			t.Fatalf("remove calendar date: %v", err)
		}
		if err := drafts.RemoveCalendar(ctx, draft.ID, "weekday"); err != nil {
			t.Fatalf("remove calendar: %v", err)
		}
		if err := drafts.RemoveStop(ctx, draft.ID, "stop-2"); err != nil {
			t.Fatalf("remove stop: %v", err)
		}
		if err := drafts.RemoveRoute(ctx, draft.ID, "route-1"); err != nil {
			t.Fatalf("remove route: %v", err)
		}
	})

	t.Run("clone edit stays separate until publish and then becomes read only", func(t *testing.T) {
		resetGTFSImportData(t, ctx, pool)
		importResult, err := imports.ImportZip(ctx, ImportOptions{AgencyID: "demo-agency", ZipPath: writeZipFixture(t, "../../testdata/gtfs/valid-small", nil), ActorID: "tester"})
		if err != nil {
			t.Fatalf("import active feed: %v", err)
		}

		explicitBlank, err := drafts.CreateDraft(ctx, CreateDraftOptions{AgencyID: "demo-agency", Name: "Explicit blank", ActorID: "tester", Blank: true})
		if err != nil {
			t.Fatalf("create explicit blank draft: %v", err)
		}
		if explicitBlank.BaseFeedVersionID != "" {
			t.Fatalf("explicit blank base = %q, want none", explicitBlank.BaseFeedVersionID)
		}

		draft, err := drafts.CreateDraft(ctx, CreateDraftOptions{AgencyID: "demo-agency", Name: "Cloned draft", ActorID: "tester"})
		if err != nil {
			t.Fatalf("create cloned draft: %v", err)
		}
		if draft.BaseFeedVersionID != importResult.FeedVersionID {
			t.Fatalf("base feed = %q, want %q", draft.BaseFeedVersionID, importResult.FeedVersionID)
		}
		if err := drafts.UpsertRoute(ctx, DraftRoute{DraftID: draft.ID, ID: "route-10", ShortName: "10", LongName: "Edited Route", RouteType: 3}); err != nil {
			t.Fatalf("edit draft route: %v", err)
		}
		var publishedName string
		if err := pool.QueryRow(ctx, `
			SELECT long_name FROM gtfs_route WHERE feed_version_id = $1 AND id = 'route-10'
		`, importResult.FeedVersionID).Scan(&publishedName); err != nil {
			t.Fatalf("query published route: %v", err)
		}
		if publishedName == "Edited Route" {
			t.Fatalf("draft edit mutated published route before publish")
		}

		repo := NewPostgresRepository(pool)
		candidates, err := repo.ListTripCandidates(ctx, "demo-agency", importResult.FeedVersionID, "20260420")
		if err != nil {
			t.Fatalf("list active candidates before publish: %v", err)
		}
		if len(candidates) != 1 || candidates[0].FeedVersionID != importResult.FeedVersionID {
			t.Fatalf("active candidates before publish = %+v", candidates)
		}

		published, err := drafts.PublishDraft(ctx, PublishDraftOptions{DraftID: draft.ID, ActorID: "tester", Notes: "publish edited route"})
		if err != nil {
			t.Fatalf("publish draft: %v", err)
		}
		if published.FeedVersionID == "" || published.FeedVersionID == importResult.FeedVersionID {
			t.Fatalf("published feed version = %q, want new studio feed", published.FeedVersionID)
		}
		var sourceType, lifecycle string
		var active bool
		if err := pool.QueryRow(ctx, `SELECT source_type, lifecycle_state, is_active FROM feed_version WHERE id = $1`, published.FeedVersionID).Scan(&sourceType, &lifecycle, &active); err != nil {
			t.Fatalf("query studio feed version: %v", err)
		}
		if sourceType != "gtfs_studio" || lifecycle != "active" || !active {
			t.Fatalf("studio feed = %s/%s/%v, want gtfs_studio/active/true", sourceType, lifecycle, active)
		}
		if err := pool.QueryRow(ctx, `SELECT long_name FROM gtfs_route WHERE feed_version_id = $1 AND id = 'route-10'`, published.FeedVersionID).Scan(&publishedName); err != nil {
			t.Fatalf("query edited published route: %v", err)
		}
		if publishedName != "Edited Route" {
			t.Fatalf("published route name = %q, want Edited Route", publishedName)
		}
		storedDraft, err := drafts.GetDraft(ctx, draft.ID)
		if err != nil {
			t.Fatalf("get published draft: %v", err)
		}
		if storedDraft.Status != DraftStatusPublished || storedDraft.LastPublishedFeedVersionID != published.FeedVersionID || storedDraft.LastPublishAttemptID == 0 {
			t.Fatalf("draft traceability = %+v, want published status/feed/attempt", storedDraft)
		}
		if err := drafts.UpsertRoute(ctx, DraftRoute{DraftID: draft.ID, ID: "route-11", RouteType: 3}); !errors.Is(err, ErrDraftNotEditable) {
			t.Fatalf("edit published draft error = %v, want ErrDraftNotEditable", err)
		}
	})

	t.Run("discarded draft is hidden by default and publish is rejected before attempt", func(t *testing.T) {
		resetGTFSImportData(t, ctx, pool)
		seedDraftAgency(t, ctx, pool, "demo-agency")
		draft, err := drafts.CreateDraft(ctx, CreateDraftOptions{AgencyID: "demo-agency", Name: "Discard me", ActorID: "tester", Blank: true})
		if err != nil {
			t.Fatalf("create draft: %v", err)
		}
		if err := drafts.DiscardDraft(ctx, DiscardDraftOptions{DraftID: draft.ID, ActorID: "tester", Reason: "not needed"}); err != nil {
			t.Fatalf("discard draft: %v", err)
		}
		visible, err := drafts.ListDrafts(ctx, "demo-agency", false)
		if err != nil {
			t.Fatalf("list visible drafts: %v", err)
		}
		if len(visible) != 0 {
			t.Fatalf("visible drafts = %+v, want discarded hidden by default", visible)
		}
		withDiscarded, err := drafts.ListDrafts(ctx, "demo-agency", true)
		if err != nil {
			t.Fatalf("list with discarded: %v", err)
		}
		if len(withDiscarded) != 1 || withDiscarded[0].Status != DraftStatusDiscarded {
			t.Fatalf("with discarded = %+v, want discarded draft", withDiscarded)
		}
		result, err := drafts.PublishDraft(ctx, PublishDraftOptions{DraftID: draft.ID, ActorID: "tester"})
		if err == nil || !errors.Is(err, ErrDraftNotEditable) {
			t.Fatalf("publish discarded err = %v, result = %+v, want ErrDraftNotEditable", err, result)
		}
		var attempts int
		if err := pool.QueryRow(ctx, `SELECT count(*) FROM gtfs_draft_publish WHERE draft_id = $1`, draft.ID).Scan(&attempts); err != nil {
			t.Fatalf("count publish attempts: %v", err)
		}
		if attempts != 0 {
			t.Fatalf("publish attempts = %d, want none because discarded draft is rejected before publisher", attempts)
		}
		var feedVersions int
		if err := pool.QueryRow(ctx, `SELECT count(*) FROM feed_version WHERE agency_id = 'demo-agency'`).Scan(&feedVersions); err != nil {
			t.Fatalf("count feed versions: %v", err)
		}
		if feedVersions != 0 {
			t.Fatalf("feed versions = %d, want none", feedVersions)
		}
	})
}

func seedDraftAgency(t *testing.T, ctx context.Context, pool *pgxpool.Pool, agencyID string) {
	t.Helper()
	_, err := pool.Exec(ctx, `
		INSERT INTO agency (id, name, timezone, contact_email, public_url)
		VALUES ($1, 'Demo Agency', 'America/Vancouver', 'dev@example.com', 'http://localhost')
	`, agencyID)
	if err != nil {
		t.Fatalf("seed agency: %v", err)
	}
}

func mustList[T any](values []T, err error) []T {
	if err != nil {
		panic(err)
	}
	return values
}

func assertLen[T any](t *testing.T, name string, values []T, want int) {
	t.Helper()
	if len(values) != want {
		t.Fatalf("%s len = %d, want %d", name, len(values), want)
	}
}
