package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/gtfs"
)

func TestDraftListHidesDiscardedByDefault(t *testing.T) {
	store := &fakeDraftStore{
		drafts: []gtfs.DraftSummary{
			{Draft: gtfs.Draft{ID: "draft-active", AgencyID: "demo-agency", Name: "Active", Status: gtfs.DraftStatusDraft}},
			{Draft: gtfs.Draft{ID: "draft-discarded", AgencyID: "demo-agency", Name: "Discarded", Status: gtfs.DraftStatusDiscarded}},
		},
	}
	handler := newHandler(store, fakePinger{})

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/gtfs-studio?agency_id=demo-agency", nil)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.Code)
	}
	if store.includeDiscarded {
		t.Fatalf("includeDiscarded = true, want false by default")
	}
	body := resp.Body.String()
	if !strings.Contains(body, "Active") {
		t.Fatalf("body missing active draft: %s", body)
	}
	if strings.Contains(body, "draft-discarded") || strings.Contains(body, ">Discarded<") {
		t.Fatalf("body includes discarded draft by default: %s", body)
	}
}

func TestDraftListCanIncludeDiscarded(t *testing.T) {
	store := &fakeDraftStore{
		drafts: []gtfs.DraftSummary{
			{Draft: gtfs.Draft{ID: "draft-discarded", AgencyID: "demo-agency", Name: "Discarded", Status: gtfs.DraftStatusDiscarded}},
		},
	}
	handler := newHandler(store, fakePinger{})

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/gtfs-studio?agency_id=demo-agency&include_discarded=1", nil)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.Code)
	}
	if !store.includeDiscarded {
		t.Fatalf("includeDiscarded = false, want true")
	}
	if body := resp.Body.String(); !strings.Contains(body, "Discarded") || !strings.Contains(body, "discarded") {
		t.Fatalf("body missing discarded draft/status: %s", body)
	}
}

func TestDraftSummaryShowsVersionVisibility(t *testing.T) {
	store := &fakeDraftStore{
		draft: gtfs.Draft{
			ID:                         "draft-1",
			AgencyID:                   "demo-agency",
			Name:                       "Draft 1",
			Status:                     gtfs.DraftStatusPublished,
			BaseFeedVersionID:          "gtfs-import-1",
			LastPublishAttemptID:       7,
			LastPublishedFeedVersionID: "gtfs-studio-7",
		},
	}
	handler := newHandler(store, fakePinger{})

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/gtfs-studio/drafts/draft-1", nil)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.Code)
	}
	body := resp.Body.String()
	for _, want := range []string{"published", "gtfs-import-1", "7", "gtfs-studio-7"} {
		if !strings.Contains(body, want) {
			t.Fatalf("summary missing %q: %s", want, body)
		}
	}
}

func TestGTFSStudioAdminRejectsUnauthenticatedAccess(t *testing.T) {
	handler := newHandlerWithAuth(&fakeDraftStore{}, fakePinger{}, authRejectAll{}, "test-csrf")
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/gtfs-studio?agency_id=demo-agency", nil)
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", resp.Code)
	}
}

type fakePinger struct{}

func (fakePinger) Ping(context.Context) error { return nil }

type authRejectAll struct{}

func (authRejectAll) Require(...auth.Role) func(http.Handler) http.Handler {
	return func(_ http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		})
	}
}

type fakeDraftStore struct {
	drafts           []gtfs.DraftSummary
	draft            gtfs.Draft
	includeDiscarded bool
}

func (f *fakeDraftStore) CreateDraft(context.Context, gtfs.CreateDraftOptions) (gtfs.Draft, error) {
	return gtfs.Draft{}, nil
}
func (f *fakeDraftStore) ListDrafts(_ context.Context, _ string, includeDiscarded bool) ([]gtfs.DraftSummary, error) {
	f.includeDiscarded = includeDiscarded
	if includeDiscarded {
		return f.drafts, nil
	}
	var visible []gtfs.DraftSummary
	for _, draft := range f.drafts {
		if draft.Status != gtfs.DraftStatusDiscarded {
			visible = append(visible, draft)
		}
	}
	return visible, nil
}
func (f *fakeDraftStore) GetDraft(context.Context, string) (gtfs.Draft, error) {
	return f.draft, nil
}
func (f *fakeDraftStore) DiscardDraft(context.Context, gtfs.DiscardDraftOptions) error {
	return nil
}
func (f *fakeDraftStore) PublishDraft(context.Context, gtfs.PublishDraftOptions) (gtfs.PublishDraftResult, error) {
	return gtfs.PublishDraftResult{}, nil
}
func (f *fakeDraftStore) UpsertAgency(context.Context, gtfs.DraftAgency) error { return nil }
func (f *fakeDraftStore) GetAgency(context.Context, string) (gtfs.DraftAgency, error) {
	return gtfs.DraftAgency{}, nil
}
func (f *fakeDraftStore) UpsertRoute(context.Context, gtfs.DraftRoute) error { return nil }
func (f *fakeDraftStore) ListRoutes(context.Context, string) ([]gtfs.DraftRoute, error) {
	return nil, nil
}
func (f *fakeDraftStore) RemoveRoute(context.Context, string, string) error { return nil }
func (f *fakeDraftStore) UpsertStop(context.Context, gtfs.DraftStop) error  { return nil }
func (f *fakeDraftStore) ListStops(context.Context, string) ([]gtfs.DraftStop, error) {
	return nil, nil
}
func (f *fakeDraftStore) RemoveStop(context.Context, string, string) error { return nil }
func (f *fakeDraftStore) UpsertCalendar(context.Context, gtfs.DraftCalendar) error {
	return nil
}
func (f *fakeDraftStore) ListCalendars(context.Context, string) ([]gtfs.DraftCalendar, error) {
	return nil, nil
}
func (f *fakeDraftStore) RemoveCalendar(context.Context, string, string) error { return nil }
func (f *fakeDraftStore) UpsertCalendarDate(context.Context, gtfs.DraftCalendarDate) error {
	return nil
}
func (f *fakeDraftStore) ListCalendarDates(context.Context, string) ([]gtfs.DraftCalendarDate, error) {
	return nil, nil
}
func (f *fakeDraftStore) RemoveCalendarDate(context.Context, string, string, string) error {
	return nil
}
func (f *fakeDraftStore) UpsertTrip(context.Context, gtfs.DraftTrip) error { return nil }
func (f *fakeDraftStore) ListTrips(context.Context, string) ([]gtfs.DraftTrip, error) {
	return nil, nil
}
func (f *fakeDraftStore) RemoveTrip(context.Context, string, string) error { return nil }
func (f *fakeDraftStore) UpsertStopTime(context.Context, gtfs.DraftStopTime) error {
	return nil
}
func (f *fakeDraftStore) ListStopTimes(context.Context, string) ([]gtfs.DraftStopTime, error) {
	return nil, nil
}
func (f *fakeDraftStore) RemoveStopTime(context.Context, string, string, int) error {
	return nil
}
func (f *fakeDraftStore) UpsertShapePoint(context.Context, gtfs.DraftShapePoint) error {
	return nil
}
func (f *fakeDraftStore) ListShapePoints(context.Context, string) ([]gtfs.DraftShapePoint, error) {
	return nil, nil
}
func (f *fakeDraftStore) RemoveShapePoint(context.Context, string, string, int) error {
	return nil
}
func (f *fakeDraftStore) UpsertFrequency(context.Context, gtfs.DraftFrequency) error {
	return nil
}
func (f *fakeDraftStore) ListFrequencies(context.Context, string) ([]gtfs.DraftFrequency, error) {
	return nil, nil
}
func (f *fakeDraftStore) RemoveFrequency(context.Context, string, string, string) error {
	return nil
}
