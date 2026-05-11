package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
)

const (
	checklistStatusOK          = "ok"
	checklistStatusNeedsReview = "needs_review"
	checklistStatusMissing     = "missing"
	checklistStatusBlocked     = "blocked"
	checklistStatusUnknown     = "unknown"

	heuristicMissing                      = "missing"
	heuristicPlaceholderLike              = "placeholder_like"
	heuristicOperatorEnteredUnverified    = "operator_entered_unverified"
	heuristicApprovalUnknown              = "approval_unknown"
	heuristicApprovalArtifactNotRetained  = "approval_artifact_not_retained"
	heuristicLocalOnly                    = "local_only"
	heuristicPilotOrReferenceRoot         = "pilot_or_reference_root"
	heuristicFinalRootCandidateUnverified = "final_root_candidate_unverified"
	heuristicNoFinalRootEvidence          = "no_final_root_evidence"
)

type operatorChecklistView struct {
	GeneratedAt time.Time                `json:"generated_at"`
	AgencyID    string                   `json:"agency_id"`
	Groups      []operatorChecklistGroup `json:"groups"`
	Counts      operatorChecklistCounts  `json:"counts"`
	Flags       operatorChecklistFlags   `json:"flags"`
}

type operatorChecklistGroup struct {
	ID    string                 `json:"id"`
	Label string                 `json:"label"`
	Rows  []operatorChecklistRow `json:"rows"`
}

type operatorChecklistRow struct {
	ID              string   `json:"id"`
	Label           string   `json:"label"`
	Status          string   `json:"status"`
	Source          string   `json:"source"`
	CurrentSignal   string   `json:"current_signal"`
	NextAction      string   `json:"next_action"`
	ClaimBoundary   string   `json:"claim_boundary"`
	DocsLinks       []string `json:"docs_links"`
	HeuristicLabels []string `json:"heuristic_labels"`
}

type operatorChecklistCounts struct {
	Groups   int            `json:"groups"`
	Rows     int            `json:"rows"`
	Statuses map[string]int `json:"statuses"`
}

type operatorChecklistFlags struct {
	ExternalEvidenceCreated    bool `json:"external_evidence_created"`
	FinalRootEvidenceCreated   bool `json:"final_root_evidence_created"`
	ConsumerStatusesChanged    bool `json:"consumer_statuses_changed"`
	ComplianceClaimed          bool `json:"compliance_claimed"`
	ProductionReadinessClaimed bool `json:"production_readiness_claimed"`
	AgencyApprovalClaimed      bool `json:"agency_approval_claimed"`
	ConsumerAcceptanceClaimed  bool `json:"consumer_acceptance_claimed"`
}

func (h *handler) renderOperationsChecklist(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "checklist")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "checklist", page)
}

func (h *handler) renderOperationsChecklistJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "checklist")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.Checklist)
}

func buildOperatorChecklist(page operationsPage) operatorChecklistView {
	groups := []operatorChecklistGroup{
		{ID: "setup", Label: "Setup", Rows: []operatorChecklistRow{
			row("setup_metadata", "Agency metadata", metadataStatus(page.Discovery.AgencyName), "publication metadata", firstNonEmpty(page.Discovery.AgencyName, page.DiscoveryError, "blank agency name"), "Enter real agency metadata, then retain any separate operator review notes outside this checklist.", privateBoundary(), docs("docs/tutorials/agency-first-run.md", "docs/tutorials/real-agency-gtfs-onboarding.md"), classifyRequiredMetadata(page.Discovery.AgencyName)),
			row("setup_license_contact", "License and contact metadata", metadataStatus(page.Discovery.License.Name, page.Discovery.License.URL, page.Discovery.TechnicalContactEmail), "publication metadata", licenseContactEvidence(page), "Replace blanks or placeholders with operator-entered license and monitored technical contact values.", "Metadata in this checklist is not agency approval.", docs("docs/compliance-evidence-checklist.md", "docs/tutorials/calitp-readiness-checklist.md"), classifyRequiredMetadata(page.Discovery.License.Name, page.Discovery.License.URL, page.Discovery.TechnicalContactEmail)),
			row("setup_public_url", "Public root URL", urlStatus(page.Discovery.PublicBaseURL), "feed discovery", firstNonEmpty(page.Discovery.PublicBaseURL, "missing public base URL"), "Review the root URL class and keep final-root conclusions outside this checklist until separate retained proof exists.", "URL heuristics are not final-root proof.", docs("docs/deployment/reference-deployment-doctor.md", "docs/phase-42-reference-deployment-doctor.md"), classifyPublicURL(page.Discovery.PublicBaseURL)),
		}},
		{ID: "feeds", Label: "Feeds", Rows: []operatorChecklistRow{
			row("feeds_schedule", "Static GTFS feed", feedChecklistStatus(page, "schedule"), "feed discovery", feedEvidence(page, "schedule"), "Import or publish the schedule feed and rerun discovery checks.", "A listed feed is not compliance proof.", docs("docs/tutorials/real-agency-gtfs-onboarding.md", "docs/tutorials/gtfs-validation-triage.md"), feedHeuristics(page, "schedule")),
			row("feeds_vehicle_positions", "Vehicle Positions feed", feedChecklistStatus(page, "vehicle_positions"), "feed discovery", feedEvidence(page, "vehicle_positions"), "Confirm telemetry, public protobuf output, and feed discovery records.", "Vehicle Positions signals are not consumer acceptance.", docs("docs/phase-42-reference-deployment-doctor.md", "docs/tutorials/operator-smoke-and-support-bundle.md"), feedHeuristics(page, "vehicle_positions")),
			row("feeds_trip_updates", "Trip Updates feed", feedChecklistStatus(page, "trip_updates"), "feed discovery and diagnostics", tripUpdatesReadinessEvidence(page), "Review prediction diagnostics and validation before relying on Trip Updates.", "Trip Updates diagnostics are not production-grade ETA proof.", docs("docs/requirements-trip-updates.md", "docs/tutorials/operator-smoke-and-support-bundle.md"), feedHeuristics(page, "trip_updates")),
			row("feeds_alerts", "Alerts feed", feedChecklistStatus(page, "alerts"), "feed discovery and Alerts Console", feedEvidence(page, "alerts"), "Use the Alerts Console for alert lifecycle checks and validate output.", "Alerts feed signals are not consumer display proof.", docs("docs/requirements-calitp-compliance.md", "docs/tutorials/operator-smoke-and-support-bundle.md"), feedHeuristics(page, "alerts")),
		}},
		{ID: "validation", Label: "Validation", Rows: []operatorChecklistRow{
			row("validation_records", "Validator records", validationChecklistStatus(page), "validation records", validationEvidence(page), "Run allowlisted validators for schedule and realtime feeds.", "Validator output is supporting diagnostics, not compliance proof.", docs("docs/tutorials/gtfs-validation-triage.md", "docs/compliance-evidence-checklist.md"), validationHeuristics(page)),
		}},
		{ID: "telemetry", Label: "Telemetry", Rows: []operatorChecklistRow{
			row("telemetry_freshness", "Telemetry freshness", telemetryChecklistStatus(page), "telemetry repository", telemetryEvidence(page), "Use the telemetry simulator guide for synthetic sample commands, then resolve stale vehicles.", "Fresh telemetry is not vendor compatibility proof.", docs("docs/tutorials/telemetry-simulator-and-device-trial.md", "docs/tutorials/agency-first-run.md", "docs/tutorials/operator-smoke-and-support-bundle.md"), telemetryHeuristics(page)),
			row("telemetry_device_bindings", "Device bindings", deviceChecklistStatus(page), "device bindings", deviceEvidence(page), "Rotate or bind a device token and store the token outside this repo.", "Binding records do not prove fleet reliability.", docs("docs/tutorials/agency-first-run.md", "docs/runbooks/small-agency-pilot-operations.md"), deviceHeuristics(page)),
		}},
		{ID: "operations", Label: "Operations", Rows: []operatorChecklistRow{
			row("operations_scorecard", "Operations scorecard", operationsChecklistStatus(page), "scorecard snapshots", operationsEvidence(page), "Run scorecard, monitor, backup, and restore-drill workflows in the operator environment.", "Operations snapshots are diagnostics, not managed hosting proof.", docs("docs/runbooks/small-agency-pilot-operations.md", "docs/tutorials/operator-smoke-and-support-bundle.md"), operationsHeuristics(page)),
			row("operations_doctor_bundle", "Doctor, smoke, and support docs", checklistStatusOK, "repo docs", "operator smoke, support bundle, and deployment doctor docs are linked", "Use the docs to run private checks and attach outputs only where policy allows.", "Viewing docs does not create evidence.", docs("docs/tutorials/operator-smoke-and-support-bundle.md", "docs/deployment/reference-deployment-doctor.md"), []string{heuristicApprovalArtifactNotRetained}),
		}},
		{ID: "consumer_workflow", Label: "Consumer Workflow", Rows: []operatorChecklistRow{
			row("consumer_prepared_packets", "Prepared packet tracker", consumerChecklistStatus(page), "docs/evidence tracker", consumerPreparedSignal(page), "Review prepared packet docs; do not change statuses without target-originated records.", "Prepared packet records are not submission, review, listing, display, ingestion, or consumer acceptance.", docs("docs/evidence/consumer-submissions/README.md", "docs/evidence/consumer-submissions/status.json"), consumerHeuristics(page)),
		}},
	}
	return operatorChecklistView{
		GeneratedAt: page.GeneratedAt,
		AgencyID:    page.AgencyID,
		Groups:      groups,
		Counts:      checklistCounts(groups),
		Flags:       operatorChecklistFlags{},
	}
}

func row(id string, label string, status string, source string, currentSignal string, nextAction string, claimBoundary string, docsLinks []string, heuristics []string) operatorChecklistRow {
	return operatorChecklistRow{
		ID:              id,
		Label:           label,
		Status:          normalizeChecklistStatus(status),
		Source:          firstNonEmpty(source, "unknown"),
		CurrentSignal:   firstNonEmpty(currentSignal, "unknown"),
		NextAction:      firstNonEmpty(nextAction, "Review this row in the Operations Console."),
		ClaimBoundary:   firstNonEmpty(claimBoundary, privateBoundary()),
		DocsLinks:       safeDocsLinks(docsLinks),
		HeuristicLabels: safeHeuristics(heuristics),
	}
}

func checklistCounts(groups []operatorChecklistGroup) operatorChecklistCounts {
	counts := operatorChecklistCounts{Groups: len(groups), Statuses: map[string]int{
		checklistStatusOK:          0,
		checklistStatusNeedsReview: 0,
		checklistStatusMissing:     0,
		checklistStatusBlocked:     0,
		checklistStatusUnknown:     0,
	}}
	for _, group := range groups {
		for _, row := range group.Rows {
			counts.Rows++
			counts.Statuses[row.Status]++
		}
	}
	return counts
}

func normalizeChecklistStatus(status string) string {
	switch status {
	case checklistStatusOK, checklistStatusNeedsReview, checklistStatusMissing, checklistStatusBlocked, checklistStatusUnknown:
		return status
	default:
		return checklistStatusUnknown
	}
}

func privateBoundary() string {
	return "Private diagnostics only; not evidence, not an evidence packet, not compliance proof, not agency approval, not consumer acceptance, and not production readiness."
}

func docs(values ...string) []string {
	return values
}

func safeDocsLinks(values []string) []string {
	var out []string
	for _, value := range values {
		v := strings.TrimSpace(value)
		if v == "" || !strings.HasPrefix(v, "docs/") || strings.Contains(v, "..") || unsafePrivateString(v) {
			continue
		}
		out = append(out, v)
	}
	return out
}

func safeHeuristics(values []string) []string {
	allowed := map[string]bool{
		heuristicMissing:                      true,
		heuristicPlaceholderLike:              true,
		heuristicOperatorEnteredUnverified:    true,
		heuristicApprovalUnknown:              true,
		heuristicApprovalArtifactNotRetained:  true,
		heuristicLocalOnly:                    true,
		heuristicPilotOrReferenceRoot:         true,
		heuristicFinalRootCandidateUnverified: true,
		heuristicNoFinalRootEvidence:          true,
	}
	var out []string
	seen := map[string]bool{}
	for _, value := range values {
		if allowed[value] && !seen[value] {
			out = append(out, value)
			seen[value] = true
		}
	}
	if len(out) == 0 {
		return []string{heuristicApprovalUnknown}
	}
	return out
}

func unsafePrivateString(value string) bool {
	lower := strings.ToLower(value)
	for _, forbidden := range []string{".cache", "file://", "/users", "localhost", "/opt/open-transit-rt", "/var/lib", "/etc", "database_url", "restore_database_url", "authorization", "cookie", "token"} {
		if strings.Contains(lower, forbidden) {
			return true
		}
	}
	return false
}

func metadataStatus(values ...string) string {
	labels := classifyRequiredMetadata(values...)
	for _, label := range labels {
		if label == heuristicMissing {
			return checklistStatusMissing
		}
		if label == heuristicPlaceholderLike {
			return checklistStatusNeedsReview
		}
	}
	return checklistStatusNeedsReview
}

func classifyRequiredMetadata(values ...string) []string {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return []string{heuristicMissing}
		}
	}
	for _, value := range values {
		if looksPlaceholderLike(value) {
			return []string{heuristicPlaceholderLike, heuristicApprovalUnknown}
		}
	}
	return []string{heuristicOperatorEnteredUnverified, heuristicApprovalUnknown}
}

func looksPlaceholderLike(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return false
	}
	for _, marker := range []string{"placeholder", "example", "demo", "sample", "test", "tbd", "todo", "changeme", "your agency", "agency name", "noreply"} {
		if strings.Contains(normalized, marker) {
			return true
		}
	}
	return false
}

func urlStatus(value string) string {
	labels := classifyPublicURL(value)
	for _, label := range labels {
		switch label {
		case heuristicMissing:
			return checklistStatusMissing
		case heuristicLocalOnly, heuristicPilotOrReferenceRoot, heuristicFinalRootCandidateUnverified, heuristicNoFinalRootEvidence:
			return checklistStatusNeedsReview
		}
	}
	return checklistStatusUnknown
}

func classifyPublicURL(value string) []string {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return []string{heuristicMissing, heuristicNoFinalRootEvidence}
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Hostname() == "" {
		return []string{heuristicPilotOrReferenceRoot, heuristicNoFinalRootEvidence}
	}
	host := strings.ToLower(parsed.Hostname())
	if host == "localhost" || host == "127.0.0.1" || host == "::1" || strings.HasPrefix(host, "10.") || strings.HasPrefix(host, "192.168.") || strings.HasPrefix(host, "172.16.") {
		return []string{heuristicLocalOnly, heuristicNoFinalRootEvidence}
	}
	if strings.Contains(host, "example") || strings.Contains(host, "reference") || strings.Contains(host, "pilot") || strings.HasSuffix(host, ".local") || strings.HasSuffix(host, ".test") {
		return []string{heuristicPilotOrReferenceRoot, heuristicNoFinalRootEvidence}
	}
	if parsed.Scheme == "https" {
		return []string{heuristicFinalRootCandidateUnverified, heuristicNoFinalRootEvidence}
	}
	return []string{heuristicPilotOrReferenceRoot, heuristicNoFinalRootEvidence}
}

func feedChecklistStatus(page operationsPage, feedType string) string {
	if page.DiscoveryError != "" {
		return checklistStatusMissing
	}
	for _, feed := range page.Discovery.Feeds {
		if feed.FeedType == feedType {
			if strings.TrimSpace(feed.CanonicalPublicURL) == "" {
				return checklistStatusMissing
			}
			return checklistStatusNeedsReview
		}
	}
	return checklistStatusMissing
}

func feedHeuristics(page operationsPage, feedType string) []string {
	if page.DiscoveryError != "" {
		return []string{heuristicMissing}
	}
	for _, feed := range page.Discovery.Feeds {
		if feed.FeedType == feedType {
			return classifyPublicURL(feed.CanonicalPublicURL)
		}
	}
	return []string{heuristicMissing, heuristicNoFinalRootEvidence}
}

func validationChecklistStatus(page operationsPage) string {
	if page.DiscoveryError != "" {
		return checklistStatusUnknown
	}
	for _, feed := range page.Discovery.Feeds {
		if feed.LastValidationAt != nil {
			return checklistStatusNeedsReview
		}
	}
	return checklistStatusMissing
}

func validationHeuristics(page operationsPage) []string {
	if page.DiscoveryError != "" {
		return []string{heuristicApprovalUnknown}
	}
	for _, feed := range page.Discovery.Feeds {
		if feed.LastValidationAt != nil {
			return []string{heuristicOperatorEnteredUnverified, heuristicApprovalUnknown}
		}
	}
	return []string{heuristicMissing}
}

func telemetryChecklistStatus(page operationsPage) string {
	if page.TelemetryError != "" {
		return checklistStatusUnknown
	}
	if len(page.Telemetry) == 0 {
		return checklistStatusMissing
	}
	if page.StaleCount > 0 {
		return checklistStatusNeedsReview
	}
	return checklistStatusOK
}

func telemetryHeuristics(page operationsPage) []string {
	if page.TelemetryError != "" {
		return []string{heuristicApprovalUnknown}
	}
	if len(page.Telemetry) == 0 {
		return []string{heuristicMissing}
	}
	return []string{heuristicOperatorEnteredUnverified, heuristicApprovalUnknown}
}

func deviceChecklistStatus(page operationsPage) string {
	if page.DeviceError != "" {
		return checklistStatusUnknown
	}
	if len(page.Devices) == 0 {
		return checklistStatusMissing
	}
	return checklistStatusNeedsReview
}

func deviceHeuristics(page operationsPage) []string {
	if page.DeviceError != "" {
		return []string{heuristicApprovalUnknown}
	}
	if len(page.Devices) == 0 {
		return []string{heuristicMissing}
	}
	return []string{heuristicOperatorEnteredUnverified, heuristicApprovalUnknown}
}

func operationsChecklistStatus(page operationsPage) string {
	if page.ScorecardError != "" {
		return checklistStatusMissing
	}
	if page.Scorecard == nil {
		return checklistStatusUnknown
	}
	return checklistStatusNeedsReview
}

func operationsHeuristics(page operationsPage) []string {
	if page.Scorecard == nil {
		return []string{heuristicMissing}
	}
	return []string{heuristicOperatorEnteredUnverified, heuristicApprovalUnknown}
}

func consumerChecklistStatus(page operationsPage) string {
	if len(page.Consumers) == 0 {
		return checklistStatusMissing
	}
	return checklistStatusNeedsReview
}

func consumerPreparedSignal(page operationsPage) string {
	if len(page.Consumers) == 0 {
		return "no prepared packet tracker rows rendered"
	}
	return "prepared packet tracker rows: Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, transit.land"
}

func consumerHeuristics(page operationsPage) []string {
	if len(page.Consumers) == 0 {
		return []string{heuristicMissing}
	}
	return []string{heuristicApprovalArtifactNotRetained, heuristicApprovalUnknown}
}

func humanHeuristicLabel(value string) string {
	switch value {
	case heuristicMissing:
		return "Missing"
	case heuristicPlaceholderLike:
		return "Placeholder-like"
	case heuristicOperatorEnteredUnverified:
		return "Operator-entered; approval unknown"
	case heuristicApprovalUnknown:
		return "Approval unknown"
	case heuristicApprovalArtifactNotRetained:
		return "Approval artifact not retained"
	case heuristicLocalOnly:
		return "Local only"
	case heuristicPilotOrReferenceRoot:
		return "Pilot/reference root"
	case heuristicFinalRootCandidateUnverified:
		return "Final-root candidate; unverified"
	case heuristicNoFinalRootEvidence:
		return "No final-root evidence"
	default:
		return "Unknown"
	}
}
