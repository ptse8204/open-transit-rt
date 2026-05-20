package main

import (
	"net/http"
	"time"

	"open-transit-rt/internal/auth"
)

type operationsAccessView struct {
	GeneratedAt  time.Time                    `json:"generated_at"`
	AgencyID     string                       `json:"agency_id"`
	Boundary     string                       `json:"boundary"`
	CurrentRoles []string                     `json:"current_roles"`
	Roles        []operationsRolePermission   `json:"roles"`
	Denied       []operationsAccessDeniedPath `json:"denied_guidance"`
}

type operationsRolePermission struct {
	ID                string `json:"id"`
	Label             string `json:"label"`
	Current           bool   `json:"current"`
	ReviewAccess      string `json:"review_access"`
	MutationAccess    string `json:"mutation_access"`
	AdministratorNote string `json:"administrator_note"`
	DoesNotProve      string `json:"does_not_prove"`
}

type operationsAccessDeniedPath struct {
	ID           string `json:"id"`
	Scenario     string `json:"scenario"`
	WhatHappened string `json:"what_happened"`
	NextAction   string `json:"next_action"`
	DoesNotProve string `json:"does_not_prove"`
}

func (h *handler) renderAccess(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "access")
	page.Access = buildOperationsAccessView(page.GeneratedAt, principal)
	renderOperationsTemplate(w, "access", page)
}

func (h *handler) renderAccessJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "access")
	writeJSON(w, http.StatusOK, buildOperationsAccessView(page.GeneratedAt, principal))
}

func buildOperationsAccessView(now time.Time, principal auth.Principal) operationsAccessView {
	currentRoles := safePrincipalRoles(principal.Roles)
	has := func(role auth.Role) bool { return principal.HasAny(role) }
	return operationsAccessView{
		GeneratedAt:  now,
		AgencyID:     principal.AgencyID,
		Boundary:     "Private role and agency-scope guidance only. This page explains configured access behavior without granting roles, switching agencies, creating evidence, or proving production multi-agency readiness.",
		CurrentRoles: currentRoles,
		Roles: []operationsRolePermission{
			{
				ID:                "admin",
				Label:             "Admin",
				Current:           has(auth.RoleAdmin),
				ReviewAccess:      "Can review private Operations Console pages.",
				MutationAccess:    "Can use admin-only browser actions where the route explicitly exposes them with role checks and form safety controls.",
				AdministratorNote: "Use for setup changes, browser GTFS import, validator runs, device credential changes, and other explicitly admin-only actions.",
				DoesNotProve:      "Admin role does not show deployment ownership, agency approval, hosted operation, or release readiness.",
			},
			{
				ID:                "editor",
				Label:             "Editor",
				Current:           has(auth.RoleEditor),
				ReviewAccess:      "Can review private pages that allow standard operations review roles.",
				MutationAccess:    "Can mutate only when a specific route explicitly allows editor access; admin-only actions remain unavailable.",
				AdministratorNote: "Use when schedule or content review needs a narrower role than admin.",
				DoesNotProve:      "Editor role does not imply authority to publish, submit, or approve external outcomes.",
			},
			{
				ID:                "operator",
				Label:             "Operator",
				Current:           has(auth.RoleOperator),
				ReviewAccess:      "Can review private operations, realtime, feed health, connector, maintenance, and help pages.",
				MutationAccess:    "Does not run admin-only setup, import, validation, or credential actions from the browser.",
				AdministratorNote: "Use for daily operational review and escalation to an administrator or deployment owner.",
				DoesNotProve:      "Operator role does not show service uptime, SLA coverage, or public launch readiness.",
			},
			{
				ID:                "read_only",
				Label:             "Read only",
				Current:           has(auth.RoleReadOnly),
				ReviewAccess:      "Can review private status, guidance, and diagnostics pages.",
				MutationAccess:    "Cannot use browser mutation forms.",
				AdministratorNote: "Use for oversight, walkthroughs, and support review without changing backend state.",
				DoesNotProve:      "Read-only access does not show an agency, consumer, vendor, or regulator has approved anything.",
			},
		},
		Denied: []operationsAccessDeniedPath{
			{
				ID:           "missing_role",
				Scenario:     "Role is not allowed",
				WhatHappened: "The signed-in session reached a private route, but none of its assigned roles are allowed for that action.",
				NextAction:   "Open this page to confirm the role set, then ask an admin to use the action or assign the appropriate role through the deployment-owned process.",
				DoesNotProve: "Denied access does not show data exists or does not exist for another agency.",
			},
			{
				ID:           "agency_scope_conflict",
				Scenario:     "Agency scope conflict",
				WhatHappened: "The requested agency scope did not match the signed-in agency, so the request stopped before page data loaded.",
				NextAction:   "Remove conflicting agency query values or use a separate signed-in session for the intended agency.",
				DoesNotProve: "Agency-scope denial does not reveal whether the requested agency has data.",
			},
			{
				ID:           "unsafe_agency_scope",
				Scenario:     "Unsafe agency identifier",
				WhatHappened: "The requested agency scope was not a single safe path segment, so the request stopped before page data loaded.",
				NextAction:   "Use the agency identifier from the signed-in session. Do not use encoded slashes, backslashes, dot segments, hidden segments, spaces, or punctuation in agency query values.",
				DoesNotProve: "Unsafe-scope denial does not reveal whether another agency or path exists.",
			},
			{
				ID:           "form_safety_check",
				Scenario:     "Form safety check failed",
				WhatHappened: "A browser form was submitted without the expected server-generated safety value.",
				NextAction:   "Reload the page and submit from the current private form; do not replay stale forms.",
				DoesNotProve: "A failed form safety check does not show the requested action would otherwise be valid.",
			},
		},
	}
}
