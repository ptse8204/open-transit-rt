package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	appdb "open-transit-rt/internal/db"
	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/server"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type draftStore interface {
	CreateDraft(context.Context, gtfs.CreateDraftOptions) (gtfs.Draft, error)
	ListDrafts(context.Context, string, bool) ([]gtfs.DraftSummary, error)
	GetDraft(context.Context, string) (gtfs.Draft, error)
	DiscardDraft(context.Context, gtfs.DiscardDraftOptions) error
	PublishDraft(context.Context, gtfs.PublishDraftOptions) (gtfs.PublishDraftResult, error)
	UpsertAgency(context.Context, gtfs.DraftAgency) error
	GetAgency(context.Context, string) (gtfs.DraftAgency, error)
	UpsertRoute(context.Context, gtfs.DraftRoute) error
	ListRoutes(context.Context, string) ([]gtfs.DraftRoute, error)
	RemoveRoute(context.Context, string, string) error
	UpsertStop(context.Context, gtfs.DraftStop) error
	ListStops(context.Context, string) ([]gtfs.DraftStop, error)
	RemoveStop(context.Context, string, string) error
	UpsertCalendar(context.Context, gtfs.DraftCalendar) error
	ListCalendars(context.Context, string) ([]gtfs.DraftCalendar, error)
	RemoveCalendar(context.Context, string, string) error
	UpsertCalendarDate(context.Context, gtfs.DraftCalendarDate) error
	ListCalendarDates(context.Context, string) ([]gtfs.DraftCalendarDate, error)
	RemoveCalendarDate(context.Context, string, string, string) error
	UpsertTrip(context.Context, gtfs.DraftTrip) error
	ListTrips(context.Context, string) ([]gtfs.DraftTrip, error)
	RemoveTrip(context.Context, string, string) error
	UpsertStopTime(context.Context, gtfs.DraftStopTime) error
	ListStopTimes(context.Context, string) ([]gtfs.DraftStopTime, error)
	RemoveStopTime(context.Context, string, string, int) error
	UpsertShapePoint(context.Context, gtfs.DraftShapePoint) error
	ListShapePoints(context.Context, string) ([]gtfs.DraftShapePoint, error)
	RemoveShapePoint(context.Context, string, string, int) error
	UpsertFrequency(context.Context, gtfs.DraftFrequency) error
	ListFrequencies(context.Context, string) ([]gtfs.DraftFrequency, error)
	RemoveFrequency(context.Context, string, string, string) error
}

type adminAuth interface {
	Require(...auth.Role) func(http.Handler) http.Handler
}

type studioHandler struct {
	drafts     draftStore
	ready      pinger
	admin      adminAuth
	csrfSecret string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := appdb.Connect(ctx, appdb.LoadConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	adminAuth, err := auth.MiddlewareFromEnv(pool)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run("gtfs-studio", newHandlerWithAuth(gtfs.NewDraftService(pool), pool, adminAuth, os.Getenv("CSRF_SECRET"))); err != nil {
		log.Fatal(err)
	}
}

func newHandler(drafts draftStore, ready pinger) http.Handler {
	return newHandlerWithAuth(drafts, ready, auth.TestAuthenticator{Principal: auth.Principal{
		Subject:  "test-admin",
		AgencyID: "demo-agency",
		Roles:    []auth.Role{auth.RoleAdmin, auth.RoleEditor, auth.RoleOperator, auth.RoleReadOnly},
		Method:   auth.MethodBearer,
	}}, "test-csrf")
}

func newHandlerWithAuth(drafts draftStore, ready pinger, admin adminAuth, csrfSecret string) http.Handler {
	h := &studioHandler{drafts: drafts, ready: ready, admin: admin, csrfSecret: csrfSecret}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/readyz", h.readyz)
	adminRead := admin.Require(auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	mux.Handle("/admin/gtfs-studio", adminRead(http.HandlerFunc(h.listDrafts)))
	mux.Handle("/admin/gtfs-studio/drafts", adminRead(http.HandlerFunc(h.createDraft)))
	mux.Handle("/admin/gtfs-studio/drafts/", adminRead(http.HandlerFunc(h.draftRoutes)))
	return mux
}

func (h *studioHandler) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"service": "gtfs-studio", "status": "ok"})
}

func (h *studioHandler) readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := h.ready.Ping(ctx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"service": "gtfs-studio", "status": "unavailable", "error": "database unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"service": "gtfs-studio", "status": "ready"})
}

func (h *studioHandler) listDrafts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin/gtfs-studio" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	includeDiscarded := r.URL.Query().Get("include_discarded") == "1"
	drafts, err := h.drafts.ListDrafts(r.Context(), principal.AgencyID, includeDiscarded)
	if err != nil {
		http.Error(w, "list drafts", http.StatusInternalServerError)
		return
	}
	render(w, "drafts", draftsPage{AgencyID: principal.AgencyID, IncludeDiscarded: includeDiscarded, Drafts: drafts, CSRFToken: h.csrfToken(r)})
}

func (h *studioHandler) createDraft(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	principal, ok := auth.RequireRole(w, r, auth.RoleEditor, auth.RoleAdmin)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	if auth.RejectAgencyConflict(w, r.FormValue("agency_id"), principal) {
		return
	}
	draft, err := h.drafts.CreateDraft(r.Context(), gtfs.CreateDraftOptions{
		AgencyID: principal.AgencyID,
		Name:     r.FormValue("name"),
		ActorID:  principal.Subject,
		Blank:    r.FormValue("mode") == "blank",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/admin/gtfs-studio/drafts/"+draft.ID, http.StatusSeeOther)
}

func (h *studioHandler) draftRoutes(w http.ResponseWriter, r *http.Request) {
	trimmed := strings.TrimPrefix(r.URL.Path, "/admin/gtfs-studio/drafts/")
	parts := strings.Split(strings.Trim(trimmed, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	draftID := parts[0]
	if len(parts) == 1 {
		h.draftSummary(w, r, draftID)
		return
	}
	switch parts[1] {
	case "discard":
		h.discardDraft(w, r, draftID)
	case "publish":
		h.publishDraft(w, r, draftID)
	default:
		remove := len(parts) == 3 && parts[2] == "remove"
		h.entity(w, r, draftID, parts[1], remove)
	}
}

func (h *studioHandler) draftSummary(w http.ResponseWriter, r *http.Request, draftID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	draft, err := h.drafts.GetDraft(r.Context(), draftID)
	if err != nil {
		http.Error(w, "draft not found", http.StatusNotFound)
		return
	}
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !h.requireDraftAgency(w, draft, principal) {
		return
	}
	render(w, "summary", summaryPage{Draft: draft, CSRFToken: h.csrfToken(r)})
}

func (h *studioHandler) discardDraft(w http.ResponseWriter, r *http.Request, draftID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	principal, ok := auth.RequireRole(w, r, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !h.requireDraftIDAgency(w, r, draftID, principal) {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	if err := h.drafts.DiscardDraft(r.Context(), gtfs.DiscardDraftOptions{DraftID: draftID, ActorID: principal.Subject, Reason: r.FormValue("reason")}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/admin/gtfs-studio/drafts/"+draftID, http.StatusSeeOther)
}

func (h *studioHandler) publishDraft(w http.ResponseWriter, r *http.Request, draftID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
	if !ok || !h.requireDraftIDAgency(w, r, draftID, principal) {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	result, err := h.drafts.PublishDraft(r.Context(), gtfs.PublishDraftOptions{DraftID: draftID, ActorID: principal.Subject, Notes: r.FormValue("notes")})
	status := http.StatusOK
	if err != nil {
		status = http.StatusBadRequest
	}
	writeJSON(w, status, result)
}

func (h *studioHandler) entity(w http.ResponseWriter, r *http.Request, draftID string, entity string, remove bool) {
	if remove {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.removeEntity(w, r, draftID, entity)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.listEntity(w, r, draftID, entity)
	case http.MethodPost:
		h.upsertEntity(w, r, draftID, entity)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *studioHandler) listEntity(w http.ResponseWriter, r *http.Request, draftID string, entity string) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !h.requireDraftIDAgency(w, r, draftID, principal) {
		return
	}
	var payload any
	var err error
	switch entity {
	case "agency":
		payload, err = h.drafts.GetAgency(r.Context(), draftID)
	case "routes":
		payload, err = h.drafts.ListRoutes(r.Context(), draftID)
	case "stops":
		payload, err = h.drafts.ListStops(r.Context(), draftID)
	case "trips":
		payload, err = h.drafts.ListTrips(r.Context(), draftID)
	case "stop_times":
		payload, err = h.drafts.ListStopTimes(r.Context(), draftID)
	case "calendars":
		payload, err = h.drafts.ListCalendars(r.Context(), draftID)
	case "calendar_dates":
		payload, err = h.drafts.ListCalendarDates(r.Context(), draftID)
	case "shape_points":
		payload, err = h.drafts.ListShapePoints(r.Context(), draftID)
	case "frequencies":
		payload, err = h.drafts.ListFrequencies(r.Context(), draftID)
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "list entity", http.StatusInternalServerError)
		return
	}
	render(w, "entity", entityPage{DraftID: draftID, Entity: entity, Fields: fieldsFor(entity), Payload: payload, CSRFToken: h.csrfToken(r)})
}

func (h *studioHandler) upsertEntity(w http.ResponseWriter, r *http.Request, draftID string, entity string) {
	principal, ok := auth.RequireRole(w, r, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !h.requireDraftIDAgency(w, r, draftID, principal) {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	if auth.RejectAgencyConflict(w, r.FormValue("agency_id"), principal) {
		return
	}
	err := h.upsertEntityForm(r, draftID, entity, principal.AgencyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/admin/gtfs-studio/drafts/"+draftID+"/"+entity, http.StatusSeeOther)
}

func (h *studioHandler) removeEntity(w http.ResponseWriter, r *http.Request, draftID string, entity string) {
	principal, ok := auth.RequireRole(w, r, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !h.requireDraftIDAgency(w, r, draftID, principal) {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	var err error
	switch entity {
	case "routes":
		err = h.drafts.RemoveRoute(r.Context(), draftID, r.FormValue("id"))
	case "stops":
		err = h.drafts.RemoveStop(r.Context(), draftID, r.FormValue("id"))
	case "trips":
		err = h.drafts.RemoveTrip(r.Context(), draftID, r.FormValue("id"))
	case "stop_times":
		err = h.drafts.RemoveStopTime(r.Context(), draftID, r.FormValue("trip_id"), atoi(r.FormValue("stop_sequence")))
	case "calendars":
		err = h.drafts.RemoveCalendar(r.Context(), draftID, r.FormValue("service_id"))
	case "calendar_dates":
		err = h.drafts.RemoveCalendarDate(r.Context(), draftID, r.FormValue("service_id"), r.FormValue("date"))
	case "shape_points":
		err = h.drafts.RemoveShapePoint(r.Context(), draftID, r.FormValue("shape_id"), atoi(r.FormValue("sequence")))
	case "frequencies":
		err = h.drafts.RemoveFrequency(r.Context(), draftID, r.FormValue("trip_id"), r.FormValue("start_time"))
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/admin/gtfs-studio/drafts/"+draftID+"/"+entity, http.StatusSeeOther)
}

func (h *studioHandler) upsertEntityForm(r *http.Request, draftID string, entity string, agencyID string) error {
	switch entity {
	case "agency":
		return h.drafts.UpsertAgency(r.Context(), gtfs.DraftAgency{DraftID: draftID, AgencyID: agencyID, Name: r.FormValue("name"), Timezone: r.FormValue("timezone"), ContactEmail: r.FormValue("contact_email"), PublicURL: r.FormValue("public_url")})
	case "routes":
		return h.drafts.UpsertRoute(r.Context(), gtfs.DraftRoute{DraftID: draftID, ID: r.FormValue("id"), ShortName: r.FormValue("short_name"), LongName: r.FormValue("long_name"), RouteType: atoi(r.FormValue("route_type"))})
	case "stops":
		return h.drafts.UpsertStop(r.Context(), gtfs.DraftStop{DraftID: draftID, ID: r.FormValue("id"), Name: r.FormValue("name"), Lat: atof(r.FormValue("lat")), Lon: atof(r.FormValue("lon"))})
	case "calendars":
		return h.drafts.UpsertCalendar(r.Context(), gtfs.DraftCalendar{DraftID: draftID, ServiceID: r.FormValue("service_id"), Monday: checkbox(r, "monday"), Tuesday: checkbox(r, "tuesday"), Wednesday: checkbox(r, "wednesday"), Thursday: checkbox(r, "thursday"), Friday: checkbox(r, "friday"), Saturday: checkbox(r, "saturday"), Sunday: checkbox(r, "sunday"), StartDate: r.FormValue("start_date"), EndDate: r.FormValue("end_date")})
	case "calendar_dates":
		return h.drafts.UpsertCalendarDate(r.Context(), gtfs.DraftCalendarDate{DraftID: draftID, ServiceID: r.FormValue("service_id"), Date: r.FormValue("date"), ExceptionType: atoi(r.FormValue("exception_type"))})
	case "trips":
		var direction *int
		if raw := strings.TrimSpace(r.FormValue("direction_id")); raw != "" {
			value := atoi(raw)
			direction = &value
		}
		return h.drafts.UpsertTrip(r.Context(), gtfs.DraftTrip{DraftID: draftID, ID: r.FormValue("id"), RouteID: r.FormValue("route_id"), ServiceID: r.FormValue("service_id"), BlockID: r.FormValue("block_id"), ShapeID: r.FormValue("shape_id"), DirectionID: direction})
	case "stop_times":
		return h.drafts.UpsertStopTime(r.Context(), gtfs.DraftStopTime{DraftID: draftID, TripID: r.FormValue("trip_id"), ArrivalTime: r.FormValue("arrival_time"), DepartureTime: r.FormValue("departure_time"), StopID: r.FormValue("stop_id"), StopSequence: atoi(r.FormValue("stop_sequence")), PickupType: optionalInt(r.FormValue("pickup_type")), DropOffType: optionalInt(r.FormValue("drop_off_type")), ShapeDistTraveled: optionalFloat(r.FormValue("shape_dist_traveled"))})
	case "shape_points":
		return h.drafts.UpsertShapePoint(r.Context(), gtfs.DraftShapePoint{DraftID: draftID, ShapeID: r.FormValue("shape_id"), Lat: atof(r.FormValue("lat")), Lon: atof(r.FormValue("lon")), Sequence: atoi(r.FormValue("sequence")), DistTraveled: optionalFloat(r.FormValue("dist_traveled"))})
	case "frequencies":
		return h.drafts.UpsertFrequency(r.Context(), gtfs.DraftFrequency{DraftID: draftID, TripID: r.FormValue("trip_id"), StartTime: r.FormValue("start_time"), EndTime: r.FormValue("end_time"), HeadwaySecs: atoi(r.FormValue("headway_secs")), ExactTimes: atoi(r.FormValue("exact_times"))})
	default:
		return fmt.Errorf("unknown entity")
	}
}

type draftsPage struct {
	AgencyID         string
	IncludeDiscarded bool
	Drafts           []gtfs.DraftSummary
	CSRFToken        string
}

type summaryPage struct {
	gtfs.Draft
	CSRFToken string
}

type entityPage struct {
	DraftID   string
	Entity    string
	Fields    []string
	Payload   any
	CSRFToken string
}

var pages = template.Must(template.New("pages").Parse(`
{{define "drafts"}}
<!doctype html><title>GTFS Studio</title><h1>GTFS Studio</h1>
<p>Agency: {{.AgencyID}}</p>
<form method="post" action="/admin/gtfs-studio/drafts">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input name="agency_id" value="{{.AgencyID}}">
<input name="name" placeholder="Draft name">
<select name="mode"><option value="clone">Clone active feed</option><option value="blank">Blank draft</option></select>
<button>Create draft</button>
</form>
<table><thead><tr><th>Name</th><th>Status</th><th>Base feed</th><th>Published feed</th><th>Latest publish</th></tr></thead><tbody>
{{range .Drafts}}<tr><td><a href="/admin/gtfs-studio/drafts/{{.ID}}">{{.Name}}</a></td><td>{{.Status}}</td><td>{{.BaseFeedVersionID}}</td><td>{{.LastPublishedFeedVersionID}}</td><td>{{.LatestPublishStatus}} {{.LatestPublishID}}</td></tr>{{end}}
</tbody></table>
{{if .IncludeDiscarded}}<p>Showing discarded drafts.</p>{{else}}<p>Discarded drafts are hidden. Add <code>include_discarded=1</code> to include them.</p>{{end}}
{{end}}

{{define "summary"}}
<!doctype html><title>Draft {{.ID}}</title><h1>{{.Name}}</h1>
<dl>
<dt>Status</dt><dd>{{.Status}}</dd>
<dt>Base feed version</dt><dd>{{.BaseFeedVersionID}}</dd>
<dt>Latest publish attempt</dt><dd>{{.LastPublishAttemptID}}</dd>
<dt>Published feed version</dt><dd>{{.LastPublishedFeedVersionID}}</dd>
</dl>
<nav>
<a href="{{.ID}}/agency">agency metadata</a>
<a href="{{.ID}}/routes">routes</a>
<a href="{{.ID}}/stops">stops</a>
<a href="{{.ID}}/trips">trips</a>
<a href="{{.ID}}/stop_times">stop_times</a>
<a href="{{.ID}}/calendars">calendars</a>
<a href="{{.ID}}/calendar_dates">calendar_dates</a>
<a href="{{.ID}}/shape_points">shape points</a>
<a href="{{.ID}}/frequencies">frequencies</a>
</nav>
<form method="post" action="{{.ID}}/publish"><input type="hidden" name="csrf_token" value="{{.CSRFToken}}"><input name="notes" placeholder="notes"><button>Publish</button></form>
<form method="post" action="{{.ID}}/discard"><input type="hidden" name="csrf_token" value="{{.CSRFToken}}"><input name="reason" placeholder="reason"><button>Discard</button></form>
{{end}}

{{define "entity"}}
<!doctype html><title>{{.Entity}}</title><h1>{{.Entity}}</h1>
<p>Minimal operational row editor for {{.Entity}}.</p>
<pre>{{printf "%+v" .Payload}}</pre>
<form method="post">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
{{range .Fields}}<label>{{.}} <input name="{{.}}"></label><br>{{end}}
<button>Save row</button>
</form>
{{end}}
`))

func render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pages.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func checkbox(r *http.Request, key string) bool {
	value := r.FormValue(key)
	return value == "1" || value == "true" || value == "on"
}

func (h *studioHandler) csrfToken(r *http.Request) string {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok || h.csrfSecret == "" {
		return ""
	}
	return auth.CSRFToken(h.csrfSecret, principal)
}

func (h *studioHandler) requireDraftIDAgency(w http.ResponseWriter, r *http.Request, draftID string, principal auth.Principal) bool {
	draft, err := h.drafts.GetDraft(r.Context(), draftID)
	if err != nil {
		http.Error(w, "draft not found", http.StatusNotFound)
		return false
	}
	return h.requireDraftAgency(w, draft, principal)
}

func (h *studioHandler) requireDraftAgency(w http.ResponseWriter, draft gtfs.Draft, principal auth.Principal) bool {
	if draft.AgencyID != principal.AgencyID {
		http.Error(w, "draft belongs to another agency", http.StatusForbidden)
		return false
	}
	return true
}

func atoi(raw string) int {
	value, _ := strconv.Atoi(strings.TrimSpace(raw))
	return value
}

func atof(raw string) float64 {
	value, _ := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	return value
}

func optionalInt(raw string) *int {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	value := atoi(raw)
	return &value
}

func optionalFloat(raw string) *float64 {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	value := atof(raw)
	return &value
}

func fieldsFor(entity string) []string {
	switch entity {
	case "agency":
		return []string{"agency_id", "name", "timezone", "contact_email", "public_url"}
	case "routes":
		return []string{"id", "short_name", "long_name", "route_type"}
	case "stops":
		return []string{"id", "name", "lat", "lon"}
	case "trips":
		return []string{"id", "route_id", "service_id", "block_id", "shape_id", "direction_id"}
	case "stop_times":
		return []string{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence", "pickup_type", "drop_off_type", "shape_dist_traveled"}
	case "calendars":
		return []string{"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date"}
	case "calendar_dates":
		return []string{"service_id", "date", "exception_type"}
	case "shape_points":
		return []string{"shape_id", "lat", "lon", "sequence", "dist_traveled"}
	case "frequencies":
		return []string{"trip_id", "start_time", "end_time", "headway_secs", "exact_times"}
	default:
		return nil
	}
}
