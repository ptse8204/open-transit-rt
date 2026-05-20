package main

import (
	"strings"

	"open-transit-rt/internal/auth"
)

type operationsAgencyScopeView struct {
	AgencyID       string   `json:"agency_id"`
	Status         string   `json:"status"`
	Source         string   `json:"source"`
	SwitcherStatus string   `json:"switcher_status"`
	QueryRule      string   `json:"query_rule"`
	Roles          []string `json:"roles"`
	NextAction     string   `json:"next_action"`
	DoesNotProve   string   `json:"does_not_prove"`
}

func buildOperationsAgencyScope(principal auth.Principal) operationsAgencyScopeView {
	roles := safePrincipalRoles(principal.Roles)
	return operationsAgencyScopeView{
		AgencyID:       strings.TrimSpace(principal.AgencyID),
		Status:         operationsStatusReady,
		Source:         "authenticated principal agency",
		SwitcherStatus: "locked to authenticated agency",
		QueryRule:      "agency_id query values must be path-segment-safe and match this agency; conflicting, encoded-slash, backslash, dot, or hidden values are rejected before page data is loaded",
		Roles:          roles,
		NextAction:     "Use a separate signed-in session for a different agency. Do not treat URL edits as agency switching.",
		DoesNotProve:   "Agency scope visibility does not show deployment-wide tenancy readiness or cross-agency administration support.",
	}
}
