package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
)

const operationsAuditLimit = 50

type auditLogReader interface {
	ListAuditLog(ctx context.Context, agencyID string, limit int) ([]compliance.AuditLogRecord, error)
}

type operationsAuditView struct {
	GeneratedAt time.Time            `json:"generated_at"`
	AgencyID    string               `json:"agency_id"`
	Boundary    string               `json:"boundary"`
	Status      string               `json:"status"`
	Rows        []operationsAuditRow `json:"rows"`
	EmptyState  string               `json:"empty_state"`
	NextAction  string               `json:"next_action"`
}

type operationsAuditRow struct {
	ID               int64  `json:"id"`
	CreatedAt        string `json:"created_at"`
	Action           string `json:"action"`
	EntityType       string `json:"entity_type"`
	EntityRef        string `json:"entity_ref"`
	ActorRecorded    bool   `json:"actor_recorded"`
	ReasonRecorded   bool   `json:"reason_recorded"`
	OldValueRecorded bool   `json:"old_value_recorded"`
	NewValueRecorded bool   `json:"new_value_recorded"`
	DoesNotShow      string `json:"does_not_show"`
}

func (h *handler) renderAudit(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "audit")
	page.Audit = h.buildOperationsAuditView(r.Context(), page.GeneratedAt, principal.AgencyID)
	renderOperationsTemplate(w, "audit", page)
}

func (h *handler) renderAuditJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "audit")
	writeJSON(w, http.StatusOK, h.buildOperationsAuditView(r.Context(), page.GeneratedAt, principal.AgencyID))
}

func (h *handler) buildOperationsAuditView(ctx context.Context, now time.Time, agencyID string) operationsAuditView {
	view := operationsAuditView{
		GeneratedAt: now,
		AgencyID:    agencyID,
		Boundary:    "Private read-only audit browser. It shows scoped audit metadata only and never renders raw old/new JSON, raw reasons, credential values, payloads, private file paths, or cross-agency rows.",
		Status:      operationsStatusMissing,
		EmptyState:  "audit log reader is not available in this runtime",
		NextAction:  "Use this page after a database-backed runtime is configured; continue to rely on route-level role and agency checks meanwhile.",
	}
	reader, ok := h.store.(auditLogReader)
	if !ok {
		return view
	}
	records, err := reader.ListAuditLog(ctx, agencyID, operationsAuditLimit)
	if err != nil {
		view.Status = operationsStatusBlocked
		view.EmptyState = "audit log metadata could not be loaded"
		view.NextAction = "Ask a technical helper to check database connectivity and audit-log query permissions from the operator environment."
		return view
	}
	if len(records) == 0 {
		view.Status = operationsStatusReady
		view.EmptyState = "no audit rows recorded for this agency"
		view.NextAction = "Mutation-capable workflows will add audit rows when their existing audit writes are used."
		return view
	}
	view.Status = operationsStatusReady
	view.EmptyState = ""
	view.NextAction = "Review recent scoped mutation metadata and investigate unexpected actions from the source workflow; raw values remain server-owned."
	for _, record := range records {
		view.Rows = append(view.Rows, operationsAuditRow{
			ID:               record.ID,
			CreatedAt:        record.CreatedAt.UTC().Format(time.RFC3339),
			Action:           auditBoundedText(record.Action, "action recorded"),
			EntityType:       auditBoundedText(record.EntityType, "entity recorded"),
			EntityRef:        auditBoundedText(record.EntityID, "entity id not shown"),
			ActorRecorded:    record.ActorRecorded,
			ReasonRecorded:   record.ReasonRecorded,
			OldValueRecorded: record.OldValueRecorded,
			NewValueRecorded: record.NewValueRecorded,
			DoesNotShow:      "Raw actor identifiers, reasons, JSON diffs, payloads, credential values, and private paths are not rendered.",
		})
	}
	return view
}

func auditBoundedText(value, fallback string) string {
	cleaned := strings.TrimSpace(strings.Join(strings.Fields(value), " "))
	if cleaned == "" || auditPrivateText(cleaned) {
		return fallback
	}
	if len(cleaned) > 80 {
		return cleaned[:67] + " [truncated]"
	}
	return cleaned
}

func auditPrivateText(value string) bool {
	lower := strings.ToLower(value)
	for _, forbidden := range []string{
		"authorization:",
		"bearer ",
		"set-cookie",
		"admin_session",
		"database_url",
		"restore_database_url",
		"postgres://",
		"postgresql://",
		"payload_json",
		"raw_report",
		"stdout",
		"stderr",
		"argv",
		"file://",
		"/users/",
		"/var/lib",
		"/etc/",
	} {
		if strings.Contains(lower, forbidden) {
			return true
		}
	}
	return false
}
