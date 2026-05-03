package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	domainalerts "open-transit-rt/internal/alerts"
	"open-transit-rt/internal/auth"
	appdb "open-transit-rt/internal/db"
	feedalerts "open-transit-rt/internal/feed/alerts"
	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/server"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type snapshotBuilder interface {
	Snapshot(ctx context.Context, generatedAt time.Time) (feedalerts.Snapshot, error)
}

type activeFeedChecker interface {
	ActiveFeedVersion(ctx context.Context, agencyID string) (gtfs.FeedVersion, error)
}

type alertStore interface {
	domainalerts.Repository
	ReconcileCanceledTripAlerts(ctx context.Context, agencyID string, actorID string, at time.Time) (domainalerts.ReconcileResult, error)
}

type adminAuth interface {
	Require(...auth.Role) func(http.Handler) http.Handler
}

type handler struct {
	agencyID   string
	builder    snapshotBuilder
	alerts     alertStore
	ready      pinger
	activeFeed activeFeedChecker
	admin      adminAuth
	csrfSecret string
}

func main() {
	agencyID := os.Getenv("AGENCY_ID")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := appdb.Connect(ctx, appdb.LoadConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	alertRepo := domainalerts.NewPostgresRepository(pool)
	builder, err := feedalerts.NewBuilder(alertRepo, feedalerts.NewPostgresHealthRepository(pool), feedalerts.Config{AgencyID: agencyID})
	if err != nil {
		log.Fatal(err)
	}

	adminAuth, err := auth.MiddlewareFromEnv(pool)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run("feed-alerts", newHandlerWithReadiness(agencyID, builder, alertRepo, pool, gtfs.NewPostgresRepository(pool), adminAuth)); err != nil {
		log.Fatal(err)
	}
}

func newHandler(builder snapshotBuilder, alerts alertStore, ready pinger) http.Handler {
	return newHandlerWithAuth(builder, alerts, ready, auth.TestAuthenticator{Principal: auth.Principal{
		Subject:  "test-admin",
		AgencyID: "demo-agency",
		Roles:    []auth.Role{auth.RoleAdmin, auth.RoleEditor, auth.RoleOperator, auth.RoleReadOnly},
		Method:   auth.MethodBearer,
	}})
}

func newHandlerWithAuth(builder snapshotBuilder, alerts alertStore, ready pinger, admin adminAuth) http.Handler {
	return newHandlerWithReadiness("demo-agency", builder, alerts, ready, readyActiveFeed{}, admin)
}

func newHandlerWithReadiness(agencyID string, builder snapshotBuilder, alerts alertStore, ready pinger, activeFeed activeFeedChecker, admin adminAuth) http.Handler {
	h := &handler{agencyID: agencyID, builder: builder, alerts: alerts, ready: ready, activeFeed: activeFeed, admin: admin, csrfSecret: os.Getenv("CSRF_SECRET")}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/readyz", h.readyz)
	mux.HandleFunc("/public/gtfsrt/alerts.pb", h.publicProto)
	adminRead := admin.Require(auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	mux.Handle("/public/gtfsrt/alerts.json", adminRead(http.HandlerFunc(h.publicJSON)))
	mux.Handle("/admin/debug/gtfsrt/alerts.json", adminRead(http.HandlerFunc(h.publicJSON)))
	mux.Handle("/admin/alerts/console", adminRead(http.HandlerFunc(h.alertsConsole)))
	mux.Handle("/admin/alerts/console/", adminRead(http.HandlerFunc(h.alertsConsoleAction)))
	mux.Handle("/admin/alerts", adminRead(http.HandlerFunc(h.adminAlerts)))
	mux.Handle("/admin/alerts/", adminRead(http.HandlerFunc(h.adminAlertAction)))
	return mux
}

type readyActiveFeed struct{}

func (readyActiveFeed) ActiveFeedVersion(_ context.Context, agencyID string) (gtfs.FeedVersion, error) {
	return gtfs.FeedVersion{ID: "test-active-feed", AgencyID: agencyID}, nil
}

func (h *handler) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"service": "feed-alerts", "status": "ok"})
}

func (h *handler) readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := h.ready.Ping(ctx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"service": "feed-alerts", "status": "unavailable", "error": "database unavailable"})
		return
	}
	if _, err := h.activeFeed.ActiveFeedVersion(ctx, h.agencyID); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"service": "feed-alerts", "status": "unavailable", "error": "active feed unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"service": "feed-alerts", "status": "ready"})
}

func (h *handler) publicProto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	snapshot, err := h.builder.Snapshot(r.Context(), time.Now().UTC())
	if err != nil {
		http.Error(w, "build alerts snapshot", http.StatusInternalServerError)
		return
	}
	payload, err := snapshot.MarshalProto()
	if err != nil {
		http.Error(w, "marshal alerts protobuf", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("Last-Modified", snapshot.GeneratedAt.Format(http.TimeFormat))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

func (h *handler) publicJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if h.agencyID != "" && principal.AgencyID != h.agencyID {
		http.Error(w, "feed debug belongs to another agency", http.StatusForbidden)
		return
	}
	snapshot, err := h.builder.Snapshot(r.Context(), time.Now().UTC())
	if err != nil {
		http.Error(w, "build alerts snapshot", http.StatusInternalServerError)
		return
	}
	if snapshot.AgencyID != "" && snapshot.AgencyID != principal.AgencyID {
		http.Error(w, "feed debug belongs to another agency", http.StatusForbidden)
		return
	}
	payload, err := snapshot.MarshalDebugJSON()
	if err != nil {
		http.Error(w, "marshal alerts debug json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Last-Modified", snapshot.GeneratedAt.Format(http.TimeFormat))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

func (h *handler) adminAlerts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
		if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
			return
		}
		alerts, err := h.alerts.ListAlerts(r.Context(), domainalerts.ListFilter{AgencyID: principal.AgencyID, Status: r.URL.Query().Get("status"), Limit: 200})
		if err != nil {
			http.Error(w, "list alerts", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"alerts": alerts})
	case http.MethodPost:
		principal, ok := auth.RequireRole(w, r, auth.RoleOperator, auth.RoleAdmin)
		if !ok {
			return
		}
		var input alertRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if auth.RejectAgencyConflict(w, input.AgencyID, principal) {
			return
		}
		input.AgencyID = principal.AgencyID
		input.ActorID = principal.Subject
		alert, err := h.alerts.UpsertAlert(r.Context(), input.toUpsertInput())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, alert)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) adminAlertAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	trimmed := strings.Trim(strings.TrimPrefix(r.URL.Path, "/admin/alerts/"), "/")
	if trimmed == "reconcile-cancellations" {
		principal, ok := auth.RequireRole(w, r, auth.RoleOperator, auth.RoleAdmin)
		if !ok {
			return
		}
		var req reconcileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if auth.RejectAgencyConflict(w, req.AgencyID, principal) {
			return
		}
		result, err := h.alerts.ReconcileCanceledTripAlerts(r.Context(), principal.AgencyID, principal.Subject, time.Now().UTC())
		if err != nil {
			http.Error(w, "reconcile canceled trip alerts", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, result)
		return
	}
	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	alertID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "invalid alert id", http.StatusBadRequest)
		return
	}
	var req alertActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	principal, ok := auth.RequireRole(w, r, auth.RoleOperator, auth.RoleAdmin)
	if !ok {
		return
	}
	if auth.RejectAgencyConflict(w, req.AgencyID, principal) {
		return
	}
	switch parts[1] {
	case "publish":
		alert, err := h.alerts.PublishAlert(r.Context(), principal.AgencyID, alertID, principal.Subject, time.Now().UTC())
		if err != nil {
			http.Error(w, "publish alert", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, alert)
	case "archive":
		if err := h.alerts.ArchiveAlert(r.Context(), principal.AgencyID, alertID, principal.Subject, req.Reason, time.Now().UTC()); err != nil {
			http.Error(w, "archive alert", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"archived": true})
	default:
		http.NotFound(w, r)
	}
}

type alertRequest struct {
	AgencyID        string                        `json:"agency_id"`
	AlertKey        string                        `json:"alert_key"`
	Cause           string                        `json:"cause"`
	Effect          string                        `json:"effect"`
	HeaderText      string                        `json:"header_text"`
	DescriptionText string                        `json:"description_text"`
	URL             string                        `json:"url"`
	ActiveStart     *time.Time                    `json:"active_start"`
	ActiveEnd       *time.Time                    `json:"active_end"`
	FeedVersionID   string                        `json:"feed_version_id"`
	SourceType      string                        `json:"source_type"`
	SourceID        string                        `json:"source_id"`
	Metadata        map[string]any                `json:"metadata"`
	ActorID         string                        `json:"actor_id"`
	Entities        []domainalerts.InformedEntity `json:"entities"`
	Publish         bool                          `json:"publish"`
}

func (r alertRequest) toUpsertInput() domainalerts.UpsertInput {
	return domainalerts.UpsertInput{
		AgencyID:        r.AgencyID,
		AlertKey:        r.AlertKey,
		Cause:           r.Cause,
		Effect:          r.Effect,
		HeaderText:      r.HeaderText,
		DescriptionText: r.DescriptionText,
		URL:             r.URL,
		ActiveStart:     r.ActiveStart,
		ActiveEnd:       r.ActiveEnd,
		FeedVersionID:   r.FeedVersionID,
		SourceType:      r.SourceType,
		SourceID:        r.SourceID,
		Metadata:        r.Metadata,
		ActorID:         r.ActorID,
		Entities:        r.Entities,
		Publish:         r.Publish,
		Now:             time.Now().UTC(),
	}
}

type alertActionRequest struct {
	AgencyID string `json:"agency_id"`
	ActorID  string `json:"actor_id"`
	Reason   string `json:"reason"`
}

type reconcileRequest struct {
	AgencyID string `json:"agency_id"`
	ActorID  string `json:"actor_id"`
}

type alertConsolePage struct {
	AgencyID  string
	Status    string
	Alerts    []domainalerts.Alert
	CSRFToken string
	Error     string
}

func (h *handler) alertsConsole(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
		if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
			return
		}
		h.renderAlertsConsole(w, r, principal, "")
	case http.MethodPost:
		principal, ok := auth.RequireRole(w, r, auth.RoleOperator, auth.RoleAdmin)
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
		input := domainalerts.UpsertInput{
			AgencyID:        principal.AgencyID,
			AlertKey:        r.FormValue("alert_key"),
			Cause:           r.FormValue("cause"),
			Effect:          r.FormValue("effect"),
			HeaderText:      r.FormValue("header_text"),
			DescriptionText: r.FormValue("description_text"),
			URL:             r.FormValue("url"),
			SourceType:      domainalerts.SourceOperator,
			ActorID:         principal.Subject,
			Entities:        alertEntitiesFromForm(r),
			Publish:         checkbox(r, "publish"),
			Now:             time.Now().UTC(),
		}
		if start := parseOptionalTime(r.FormValue("active_start")); start != nil {
			input.ActiveStart = start
		}
		if end := parseOptionalTime(r.FormValue("active_end")); end != nil {
			input.ActiveEnd = end
		}
		if _, err := h.alerts.UpsertAlert(r.Context(), input); err != nil {
			h.renderAlertsConsole(w, r, principal, err.Error())
			return
		}
		http.Redirect(w, r, "/admin/alerts/console", http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) alertsConsoleAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	principal, ok := auth.RequireRole(w, r, auth.RoleOperator, auth.RoleAdmin)
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
	trimmed := strings.Trim(strings.TrimPrefix(r.URL.Path, "/admin/alerts/console/"), "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	alertID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "invalid alert id", http.StatusBadRequest)
		return
	}
	switch parts[1] {
	case "publish":
		if _, err := h.alerts.PublishAlert(r.Context(), principal.AgencyID, alertID, principal.Subject, time.Now().UTC()); err != nil {
			http.Error(w, "publish alert", http.StatusBadRequest)
			return
		}
	case "archive":
		if err := h.alerts.ArchiveAlert(r.Context(), principal.AgencyID, alertID, principal.Subject, r.FormValue("reason"), time.Now().UTC()); err != nil {
			http.Error(w, "archive alert", http.StatusBadRequest)
			return
		}
	default:
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/admin/alerts/console", http.StatusSeeOther)
}

func (h *handler) renderAlertsConsole(w http.ResponseWriter, r *http.Request, principal auth.Principal, formError string) {
	status := r.URL.Query().Get("status")
	alerts, err := h.alerts.ListAlerts(r.Context(), domainalerts.ListFilter{AgencyID: principal.AgencyID, Status: status, Limit: 200})
	page := alertConsolePage{AgencyID: principal.AgencyID, Status: status, Alerts: alerts, CSRFToken: alertCSRFToken(h.csrfSecret, principal), Error: formError}
	if err != nil {
		page.Error = "list alerts"
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := alertConsoleTemplate.Execute(w, page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func alertEntitiesFromForm(r *http.Request) []domainalerts.InformedEntity {
	entity := domainalerts.InformedEntity{
		RouteID:   strings.TrimSpace(r.FormValue("route_id")),
		StopID:    strings.TrimSpace(r.FormValue("stop_id")),
		TripID:    strings.TrimSpace(r.FormValue("trip_id")),
		StartDate: strings.TrimSpace(r.FormValue("start_date")),
		StartTime: strings.TrimSpace(r.FormValue("start_time")),
	}
	if entity.RouteID == "" && entity.StopID == "" && entity.TripID == "" {
		return nil
	}
	return []domainalerts.InformedEntity{entity}
}

func parseOptionalTime(raw string) *time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil
	}
	t := parsed.UTC()
	return &t
}

func alertCSRFToken(secret string, principal auth.Principal) string {
	if strings.TrimSpace(secret) == "" {
		return ""
	}
	return auth.CSRFToken(secret, principal)
}

func checkbox(r *http.Request, key string) bool {
	value := r.FormValue(key)
	return value == "1" || value == "true" || value == "on"
}

var alertConsoleTemplate = template.Must(template.New("alerts-console").Funcs(template.FuncMap{
	"formatTimePtr": func(t *time.Time) string {
		if t == nil || t.IsZero() {
			return "not available"
		}
		return t.UTC().Format(time.RFC3339)
	},
	"formatTime": func(t time.Time) string {
		if t.IsZero() {
			return "not available"
		}
		return t.UTC().Format(time.RFC3339)
	},
}).Parse(`<!doctype html><html><head><meta charset="utf-8"><title>Alerts Console</title>
<style>
body{font-family:system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;margin:2rem;line-height:1.4;color:#1f2933}
nav a{margin-right:1rem} table{border-collapse:collapse;width:100%;margin:1rem 0} th,td{border:1px solid #d8dee4;padding:.45rem;text-align:left;vertical-align:top}
th{background:#f6f8fa}.warning{background:#fff8c5;padding:.5rem} label{display:block;margin:.35rem 0} input,textarea,select{min-width:22rem;max-width:100%;padding:.35rem}
</style></head><body>
<h1>Alerts Console</h1>
<p>Agency: <strong>{{.AgencyID}}</strong></p>
<nav><a href="/admin/operations">Operations Console</a><a href="/admin/operations/feeds">Feeds</a><a href="/admin/gtfs-studio">GTFS Studio</a></nav>
<p class="warning">Alerts shown here are operator records. They are not evidence that any consumer has accepted or displayed the alert.</p>
{{if .Error}}<p class="warning">{{.Error}}</p>{{end}}
<h2>Alerts</h2>
<p>Filter: <a href="/admin/alerts/console">all</a> <a href="/admin/alerts/console?status=draft">draft</a> <a href="/admin/alerts/console?status=published">published</a> <a href="/admin/alerts/console?status=archived">archived</a></p>
{{if not .Alerts}}<p class="warning">No alerts are recorded for this filter. Next action: create a draft alert below when an agency-approved service message exists.</p>{{else}}
<table><thead><tr><th>ID</th><th>Status</th><th>Key</th><th>Header</th><th>Cause/effect</th><th>Active window</th><th>Affected entities</th><th>Actions</th></tr></thead><tbody>
{{range .Alerts}}<tr><td>{{.ID}}</td><td>{{.Status}}</td><td>{{.AlertKey}}</td><td>{{.HeaderText}}</td><td>{{.Cause}} / {{.Effect}}</td><td>{{formatTimePtr .ActiveStart}} to {{formatTimePtr .ActiveEnd}}</td><td>{{range .Entities}}route={{.RouteID}} trip={{.TripID}} stop={{.StopID}} {{end}}</td><td>
<form method="post" action="/admin/alerts/console/{{.ID}}/publish"><input type="hidden" name="csrf_token" value="{{$.CSRFToken}}"><input type="hidden" name="agency_id" value="{{$.AgencyID}}"><button>Publish</button></form>
<form method="post" action="/admin/alerts/console/{{.ID}}/archive"><input type="hidden" name="csrf_token" value="{{$.CSRFToken}}"><input type="hidden" name="agency_id" value="{{$.AgencyID}}"><input name="reason" placeholder="archive reason"><button>Archive</button></form>
</td></tr>{{end}}
</tbody></table>{{end}}

<h2>Create Or Update Alert</h2>
<form method="post" action="/admin/alerts/console">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="agency_id" value="{{.AgencyID}}">
<label>Alert key <input name="alert_key" required></label>
<label>Header <input name="header_text" required></label>
<label>Description <textarea name="description_text"></textarea></label>
<label>Cause <input name="cause" value="unknown_cause"></label>
<label>Effect <input name="effect" value="unknown_effect"></label>
<label>URL <input name="url"></label>
<label>Active start RFC3339 <input name="active_start" placeholder="2026-04-26T12:00:00Z"></label>
<label>Active end RFC3339 <input name="active_end" placeholder="2026-04-26T14:00:00Z"></label>
<fieldset><legend>Affected entity, optional</legend>
<label>Route ID <input name="route_id"></label>
<label>Trip ID <input name="trip_id"></label>
<label>Stop ID <input name="stop_id"></label>
<label>Start date <input name="start_date" placeholder="YYYYMMDD"></label>
<label>Start time <input name="start_time" placeholder="HH:MM:SS"></label>
</fieldset>
<label><input type="checkbox" name="publish"> Publish immediately</label>
<button>Save alert</button>
</form>
</body></html>`))

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
