package main

import (
	"context"
	"encoding/json"
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
	h := &handler{agencyID: agencyID, builder: builder, alerts: alerts, ready: ready, activeFeed: activeFeed, admin: admin}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/readyz", h.readyz)
	mux.HandleFunc("/public/gtfsrt/alerts.pb", h.publicProto)
	adminRead := admin.Require(auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	mux.Handle("/public/gtfsrt/alerts.json", adminRead(http.HandlerFunc(h.publicJSON)))
	mux.Handle("/admin/debug/gtfsrt/alerts.json", adminRead(http.HandlerFunc(h.publicJSON)))
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
	snapshot, err := h.builder.Snapshot(r.Context(), time.Now().UTC())
	if err != nil {
		http.Error(w, "build alerts snapshot", http.StatusInternalServerError)
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

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
