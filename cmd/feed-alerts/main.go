package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	"open-transit-rt/internal/tenant"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type snapshotBuilder interface {
	Snapshot(ctx context.Context, generatedAt time.Time) (feedalerts.Snapshot, error)
}

type agencySnapshotBuilder interface {
	SnapshotForAgency(ctx context.Context, agencyID string, generatedAt time.Time) (feedalerts.Snapshot, error)
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
	mux.HandleFunc("/public/agencies/", h.publicAgencyRoute)
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
	h.writeProto(w, r, "")
}

func (h *handler) publicAgencyRoute(w http.ResponseWriter, r *http.Request) {
	agencyID, suffix, matched, err := tenant.PublicAgencyRoute(r.URL.EscapedPath())
	if !matched {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "invalid agency route", http.StatusBadRequest)
		return
	}
	if suffix != "/gtfsrt/alerts.pb" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.writeProto(w, r, agencyID)
}

func (h *handler) writeProto(w http.ResponseWriter, r *http.Request, agencyID string) {
	snapshot, err := h.alertsSnapshot(r.Context(), agencyID, time.Now().UTC())
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

func (h *handler) alertsSnapshot(ctx context.Context, agencyID string, generatedAt time.Time) (feedalerts.Snapshot, error) {
	if agencyID == "" || agencyID == h.agencyID {
		return h.builder.Snapshot(ctx, generatedAt)
	}
	if agencyBuilder, ok := h.builder.(agencySnapshotBuilder); ok {
		return agencyBuilder.SnapshotForAgency(ctx, agencyID, generatedAt)
	}
	snapshot, err := h.builder.Snapshot(ctx, generatedAt)
	if err != nil {
		return feedalerts.Snapshot{}, err
	}
	if snapshot.AgencyID != agencyID {
		return feedalerts.Snapshot{}, errors.New("alerts builder cannot build requested agency")
	}
	return snapshot, nil
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
	AgencyID       string
	Status         string
	Alerts         []domainalerts.Alert
	CSRFToken      string
	Error          string
	LifecycleRows  []alertConsoleReviewRow
	ServiceRows    []alertConsoleReviewRow
	DisruptionRows []alertDisruptionTemplateRow
	ValidationRows []alertConsoleReviewRow
	UsefulnessRows []alertConsoleReviewRow
	Boundary       string
	DoesNotProve   string
	ClaimFlags     alertConsoleClaimFlags
	GeneratedAt    time.Time
}

type alertConsoleReviewRow struct {
	ID           string
	Label        string
	Status       string
	Signal       string
	NextAction   string
	DoesNotProve string
}

type alertDisruptionTemplateRow struct {
	ID             string
	Situation      string
	Cause          string
	Effect         string
	EntityGuidance string
	WindowGuidance string
	ReviewStep     string
	DoesNotProve   string
}

type alertConsoleClaimFlags struct {
	ExternalEvidenceCreated       bool
	ConsumerStatusesChanged       bool
	ComplianceClaimed             bool
	ProductionReadinessClaimed    bool
	PublicLaunchClaimed           bool
	ConsumerAcceptanceClaimed     bool
	VendorCompatibilityClaimed    bool
	HardwareCertificationClaimed  bool
	SLAClaimed                    bool
	HostedSaaSClaimed             bool
	BrowserExternalContactEnabled bool
	RawPrivatePayloadsShown       bool
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
	if trimmed == "reconcile-cancellations" {
		if _, err := h.alerts.ReconcileCanceledTripAlerts(r.Context(), principal.AgencyID, principal.Subject, time.Now().UTC()); err != nil {
			http.Error(w, "reconcile canceled trip alerts", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/alerts/console", http.StatusSeeOther)
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
	now := time.Now().UTC()
	page := alertConsolePage{
		AgencyID:       principal.AgencyID,
		Status:         status,
		Alerts:         alerts,
		CSRFToken:      alertCSRFToken(h.csrfSecret, principal),
		Error:          formError,
		LifecycleRows:  alertLifecycleRows(alerts, now),
		ServiceRows:    alertServiceDisruptionRows(alerts, now),
		DisruptionRows: alertDisruptionTemplates(),
		ValidationRows: alertValidationRows(alerts),
		UsefulnessRows: alertUsefulnessRows(alerts, now),
		Boundary:       "Private Alerts Console workflow only. This page helps operators review alert lifecycle, cancellation linkage, validation steps, and feed usefulness without contacting consumers or collecting evidence.",
		DoesNotProve:   "Alerts Console records do not prove consumer display, consumer acceptance, public launch, compliance, production readiness, vendor compatibility, hardware certification, hosted service readiness, SLA coverage, or agency approval.",
		GeneratedAt:    now,
	}
	if err != nil {
		page.Error = "list alerts"
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := alertConsoleTemplate.Execute(w, page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func alertLifecycleRows(alerts []domainalerts.Alert, now time.Time) []alertConsoleReviewRow {
	counts := alertOperationCounts(alerts, now)
	return []alertConsoleReviewRow{
		{
			ID:           "lifecycle_counts",
			Label:        "Lifecycle counts",
			Status:       alertReviewStatus(len(alerts) > 0),
			Signal:       fmt.Sprintf("draft=%d; published=%d; archived=%d; active_published=%d; upcoming_published=%d; expired_published=%d", counts.status[domainalerts.StatusDraft], counts.status[domainalerts.StatusPublished], counts.status[domainalerts.StatusArchived], counts.activePublished, counts.upcomingPublished, counts.expiredPublished),
			NextAction:   "Review draft alerts, publish only agency-approved messages, and archive expired or resolved alerts.",
			DoesNotProve: "Lifecycle counts do not prove consumer display, public launch, or operational completeness.",
		},
		{
			ID:           "cancellation_linkage",
			Label:        "Canceled-trip linkage",
			Status:       alertReviewStatus(counts.reconciledCancellation > 0),
			Signal:       fmt.Sprintf("%d alerts were created or maintained by the cancellation reconciler.", counts.reconciledCancellation),
			NextAction:   "Use the reconciliation action only after cancellation overrides are reviewed; then validate Trip Updates and Alerts together.",
			DoesNotProve: "A reconciled cancellation alert does not prove consumer ingestion, rider display, or real disruption handling quality.",
		},
		{
			ID:           "authoring_review",
			Label:        "Operator-authored alerts",
			Status:       alertReviewStatus(counts.operatorAuthored > 0),
			Signal:       fmt.Sprintf("%d alerts are operator-authored or source-unlabeled.", counts.operatorAuthored),
			NextAction:   "Confirm header text, affected entities, active window, cause/effect, and archive policy before publishing.",
			DoesNotProve: "Operator-authored alert records do not prove agency approval or compliance.",
		},
		{
			ID:           "authoring_preflight",
			Label:        "Authoring preflight",
			Status:       "review_required",
			Signal:       fmt.Sprintf("agency_wide_or_unscoped=%d; published_without_active_end=%d", counts.agencyWide, counts.indefinitePublished),
			NextAction:   "Avoid agency-wide or indefinite alerts unless the disruption truly applies agency-wide; invalid RFC3339 windows are ignored by the form parser, so review saved windows before publishing.",
			DoesNotProve: "Preflight guidance does not prove alert correctness, agency approval, consumer display, or compliance.",
		},
	}
}

func alertServiceDisruptionRows(alerts []domainalerts.Alert, now time.Time) []alertConsoleReviewRow {
	counts := alertOperationCounts(alerts, now)
	staleNeedsReview := counts.expiredPublished + counts.indefinitePublished
	scopeNeedsReview := counts.agencyWide + counts.alertsWithoutEntities
	return []alertConsoleReviewRow{
		{
			ID:           "active_disruption_review",
			Label:        "Active disruption review",
			Status:       alertReviewStatus(counts.activePublished > 0 || counts.status[domainalerts.StatusDraft] > 0),
			Signal:       fmt.Sprintf("active_published=%d; draft_waiting=%d; validation_after_change=required", counts.activePublished, counts.status[domainalerts.StatusDraft]),
			NextAction:   "For each active disruption, confirm the affected entity, active window, Trip Updates relationship, and post-change Alerts validation.",
			DoesNotProve: "Active disruption review does not prove consumer display, agency approval, public launch, or compliance.",
		},
		{
			ID:           "stale_or_indefinite_review",
			Label:        "Stale or indefinite alert review",
			Status:       alertReviewStatus(staleNeedsReview == 0 && len(alerts) > 0),
			Signal:       fmt.Sprintf("expired_published=%d; published_without_active_end=%d", counts.expiredPublished, counts.indefinitePublished),
			NextAction:   "Archive expired alerts and add bounded end times unless an agency-wide open-ended disruption is intentionally reviewed.",
			DoesNotProve: "A clear stale-alert review does not prove operational completeness, consumer display, or SLA coverage.",
		},
		{
			ID:           "entity_scope_review",
			Label:        "Affected entity scope review",
			Status:       alertReviewStatus(scopeNeedsReview == 0 && len(alerts) > 0),
			Signal:       fmt.Sprintf("agency_wide_or_unscoped=%d; alerts_without_entities=%d", counts.agencyWide, counts.alertsWithoutEntities),
			NextAction:   "Prefer route, stop, or trip selectors when the disruption is not truly agency-wide.",
			DoesNotProve: "Entity scope review does not prove alert correctness, consumer display, or compliance.",
		},
		{
			ID:           "cancellation_pairing_review",
			Label:        "Cancellation pairing review",
			Status:       alertReviewStatus(counts.reconciledCancellation > 0),
			Signal:       fmt.Sprintf("cancellation_reconciler_alerts=%d; missing_alert_hint_action=reconcile_then_validate", counts.reconciledCancellation),
			NextAction:   "When canceled Trip Updates are present, reconcile canceled-trip alerts and validate both Trip Updates and Alerts.",
			DoesNotProve: "Cancellation pairing review does not prove real-world disruption handling quality or consumer ingestion.",
		},
	}
}

type alertOperationSummary struct {
	status                 map[string]int
	activePublished        int
	upcomingPublished      int
	expiredPublished       int
	reconciledCancellation int
	operatorAuthored       int
	agencyWide             int
	alertsWithoutEntities  int
	indefinitePublished    int
}

func alertOperationCounts(alerts []domainalerts.Alert, now time.Time) alertOperationSummary {
	counts := alertOperationSummary{status: map[string]int{}}
	for _, alert := range alerts {
		counts.status[alert.Status]++
		if alert.SourceType == domainalerts.SourceCancellationReconciler {
			counts.reconciledCancellation++
		}
		if alert.SourceType == "" || alert.SourceType == domainalerts.SourceOperator {
			counts.operatorAuthored++
		}
		if len(alert.Entities) == 0 {
			counts.agencyWide++
			counts.alertsWithoutEntities++
		}
		if alert.Status == domainalerts.StatusPublished && alert.ActiveEnd == nil {
			counts.indefinitePublished++
		}
		if alert.Status != domainalerts.StatusPublished {
			continue
		}
		switch alertWindowState(alert, now) {
		case "active":
			counts.activePublished++
		case "upcoming":
			counts.upcomingPublished++
		case "expired":
			counts.expiredPublished++
		}
	}
	return counts
}

func alertWindowState(alert domainalerts.Alert, now time.Time) string {
	if alert.ActiveEnd != nil && alert.ActiveEnd.Before(now) {
		return "expired"
	}
	if alert.ActiveStart != nil && alert.ActiveStart.After(now) {
		return "upcoming"
	}
	return "active"
}

func alertReviewStatus(ok bool) string {
	if ok {
		return "review_available"
	}
	return "needs_review"
}

func alertDisruptionTemplates() []alertDisruptionTemplateRow {
	none := "Template guidance does not prove consumer display, public launch, compliance, production readiness, or real disruption handling quality."
	return []alertDisruptionTemplateRow{
		{
			ID:             "canceled_trip",
			Situation:      "Canceled trip",
			Cause:          "other_cause",
			Effect:         "no_service",
			EntityGuidance: "Include route, trip_id, start_date, and start_time when known.",
			WindowGuidance: "Use the cancellation window; archive when service is restored or the operating day ends.",
			ReviewStep:     "Reconcile cancellation overrides, then validate Trip Updates and Alerts together.",
			DoesNotProve:   none,
		},
		{
			ID:             "detour",
			Situation:      "Detour",
			Cause:          "construction, police_activity, accident, or other_cause",
			Effect:         "detour or modified_service",
			EntityGuidance: "Prefer route and stop selectors; include trip selectors only when the detour is trip-specific.",
			WindowGuidance: "Set a bounded active window and archive once the detour ends.",
			ReviewStep:     "Review affected stops and Trip Updates withholding before publishing.",
			DoesNotProve:   none,
		},
		{
			ID:             "significant_delay",
			Situation:      "Significant delay",
			Cause:          "accident, weather, technical_problem, or other_cause",
			Effect:         "significant_delays or reduced_service",
			EntityGuidance: "Use route selectors first; add stops or trips only when the delay is localized.",
			WindowGuidance: "Set a bounded review window and archive when headways recover.",
			ReviewStep:     "Check Realtime Center freshness and Trip Updates usefulness before relying on ETA-like output.",
			DoesNotProve:   none,
		},
		{
			ID:             "stop_closure",
			Situation:      "Stop closure or stop moved",
			Cause:          "construction, maintenance, police_activity, or other_cause",
			Effect:         "stop_moved, no_service, or modified_service",
			EntityGuidance: "Include stop_id and route_id when known; avoid agency-wide alerts for local stop issues.",
			WindowGuidance: "Use the known closure window and archive promptly.",
			ReviewStep:     "Confirm the stop is in the active GTFS and run Alerts validation after publishing.",
			DoesNotProve:   none,
		},
		{
			ID:             "modified_or_added_service",
			Situation:      "Modified or added service",
			Cause:          "holiday, maintenance, or other_cause",
			Effect:         "modified_service or additional_service",
			EntityGuidance: "Use route selectors and add trip selectors only when the modified service is trip-specific.",
			WindowGuidance: "Use the special-service window; review GTFS Workbench if the change belongs in schedule data.",
			ReviewStep:     "Confirm whether the change should be static GTFS, Trip Updates, Alerts, or a combination.",
			DoesNotProve:   none,
		},
	}
}

func alertValidationRows(alerts []domainalerts.Alert) []alertConsoleReviewRow {
	published := countAlertsWithStatus(alerts, domainalerts.StatusPublished)
	return []alertConsoleReviewRow{
		{
			ID:           "gtfs_rt_alerts_validation",
			Label:        "GTFS-RT Alerts validation",
			Status:       alertReviewStatus(published > 0),
			Signal:       fmt.Sprintf("%d published alert records are visible in this private list.", published),
			NextAction:   "After publishing or archiving, run the configured realtime validator for Alerts through Validation Center or `make validate`.",
			DoesNotProve: "A local validation pass does not prove consumer acceptance, compliance, public launch, or production readiness.",
		},
		{
			ID:           "feed_health_review",
			Label:        "Alerts feed health review",
			Status:       "review_required",
			Signal:       "Review `/public/gtfsrt/alerts.pb` through Feed Health after lifecycle changes.",
			NextAction:   "Check freshness, validator context, and whether the feed is intentionally empty or contains active alerts.",
			DoesNotProve: "Feed reachability does not prove consumer display or disruption workflow completeness.",
		},
		{
			ID:           "missing_alert_hints",
			Label:        "Missing-alert hints",
			Status:       "review_required",
			Signal:       "Prediction Lab and Realtime Center expose cancellation-alert missing counters when canceled Trip Updates need alert linkage.",
			NextAction:   "When cancellation_alert_links_missing is nonzero, open this console, reconcile cancellations, then validate Trip Updates and Alerts.",
			DoesNotProve: "Missing-alert hints do not prove a real cancellation was agency-approved or displayed to riders.",
		},
	}
}

func alertUsefulnessRows(alerts []domainalerts.Alert, now time.Time) []alertConsoleReviewRow {
	activePublished := 0
	for _, alert := range alerts {
		if alert.Status == domainalerts.StatusPublished && alertWindowState(alert, now) == "active" {
			activePublished++
		}
	}
	status := "valid_empty_or_needs_review"
	next := "If there is an active disruption, create or publish an agency-approved alert; otherwise keep the feed valid and intentionally empty."
	if activePublished > 0 {
		status = "review_available"
		next = "Validate the Alerts feed, confirm affected entities and windows, then archive stale alerts when the disruption ends."
	}
	return []alertConsoleReviewRow{{
		ID:           "public_feed_usefulness",
		Label:        "Public Alerts feed usefulness",
		Status:       status,
		Signal:       fmt.Sprintf("%d active published alerts are visible in this private list.", activePublished),
		NextAction:   next,
		DoesNotProve: "Useful local alert records do not prove consumer ingestion, consumer acceptance, public launch, compliance, or production readiness.",
	}}
}

func countAlertsWithStatus(alerts []domainalerts.Alert, status string) int {
	total := 0
	for _, alert := range alerts {
		if alert.Status == status {
			total++
		}
	}
	return total
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
<p class="warning">{{.Boundary}} {{.DoesNotProve}}</p>
{{if .Error}}<p class="warning">{{.Error}}</p>{{end}}
<h2>Lifecycle Dashboard</h2>
<table><thead><tr><th>Review</th><th>Status</th><th>Signal</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .LifecycleRows}}<tr><td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td><td>{{.Status}}</td><td>{{.Signal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h2>Service Disruption Review</h2>
<table><thead><tr><th>Review</th><th>Status</th><th>Signal</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .ServiceRows}}<tr><td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td><td>{{.Status}}</td><td>{{.Signal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h2>Cancellation Linkage</h2>
<p class="warning">Use this only after canceled-trip overrides are reviewed. It creates or updates cancellation alerts from existing private overrides and links matching missing-alert review incidents; it does not contact consumers or prove display.</p>
<form method="post" action="/admin/alerts/console/reconcile-cancellations">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="agency_id" value="{{.AgencyID}}">
<button>Reconcile canceled-trip alerts</button>
</form>
<h2>Disruption Templates</h2>
<table><thead><tr><th>Situation</th><th>Cause</th><th>Effect</th><th>Entities</th><th>Window</th><th>Review step</th><th>Does not prove</th></tr></thead><tbody>
{{range .DisruptionRows}}<tr><td><strong>{{.Situation}}</strong><br><code>{{.ID}}</code></td><td>{{.Cause}}</td><td>{{.Effect}}</td><td>{{.EntityGuidance}}</td><td>{{.WindowGuidance}}</td><td>{{.ReviewStep}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h2>Validation And Feed Usefulness</h2>
<table><thead><tr><th>Review</th><th>Status</th><th>Signal</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .ValidationRows}}<tr><td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td><td>{{.Status}}</td><td>{{.Signal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
{{range .UsefulnessRows}}<tr><td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td><td>{{.Status}}</td><td>{{.Signal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h2>Claim Flags</h2>
<table><tbody>
<tr><th>external_evidence_created</th><td>{{.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th>consumer_statuses_changed</th><td>{{.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th>compliance_claimed</th><td>{{.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th>production_readiness_claimed</th><td>{{.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th>public_launch_claimed</th><td>{{.ClaimFlags.PublicLaunchClaimed}}</td></tr>
<tr><th>consumer_acceptance_claimed</th><td>{{.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th>browser_external_contact_enabled</th><td>{{.ClaimFlags.BrowserExternalContactEnabled}}</td></tr>
<tr><th>raw_private_payloads_shown</th><td>{{.ClaimFlags.RawPrivatePayloadsShown}}</td></tr>
</tbody></table>
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
