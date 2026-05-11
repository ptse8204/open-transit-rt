package main

import (
	"net/http"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
)

type agencyLaunchpadView struct {
	GeneratedAt   time.Time                 `json:"generated_at"`
	AgencyID      string                    `json:"agency_id"`
	Boundary      string                    `json:"boundary"`
	Sections      []agencyLaunchpadSection  `json:"sections"`
	Counts        agencyLaunchpadCounts     `json:"counts"`
	ClaimFlags    agencyLaunchpadClaimFlags `json:"claim_flags"`
	DecisionNotes []agencyLaunchpadDecision `json:"decision_notes"`
}

type agencyLaunchpadSection struct {
	ID                 string   `json:"id"`
	Label              string   `json:"label"`
	Status             string   `json:"status"`
	CurrentSignal      string   `json:"current_signal"`
	NextActions        []string `json:"next_actions"`
	DocsLinks          []string `json:"docs_links"`
	CommandSuggestions []string `json:"command_suggestions"`
	AdminLinks         []string `json:"admin_links"`
	ClaimBoundary      string   `json:"claim_boundary"`
}

type agencyLaunchpadCounts struct {
	Sections int            `json:"sections"`
	Statuses map[string]int `json:"statuses"`
}

type agencyLaunchpadClaimFlags struct {
	ExternalEvidenceCreated    bool `json:"external_evidence_created"`
	FinalRootEvidenceCreated   bool `json:"final_root_evidence_created"`
	ConsumerStatusesChanged    bool `json:"consumer_statuses_changed"`
	ComplianceClaimed          bool `json:"compliance_claimed"`
	ProductionReadinessClaimed bool `json:"production_readiness_claimed"`
	AgencyApprovalClaimed      bool `json:"agency_approval_claimed"`
	ConsumerAcceptanceClaimed  bool `json:"consumer_acceptance_claimed"`
	PublicLaunchClaimed        bool `json:"public_launch_claimed"`
	HostedSaaSClaimed          bool `json:"hosted_saas_claimed"`
	VendorCompatibilityClaimed bool `json:"vendor_compatibility_claimed"`
	ProductionGradeETAClaimed  bool `json:"production_grade_eta_claimed"`
}

type agencyLaunchpadDecision struct {
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	Boundary      string `json:"boundary"`
}

func (h *handler) renderLaunchpad(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "launchpad")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "launchpad", page)
}

func (h *handler) renderLaunchpadJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "launchpad")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.Launchpad)
}

func buildAgencyLaunchpad(page operationsPage) agencyLaunchpadView {
	sections := []agencyLaunchpadSection{
		launchpadSection(
			"setup",
			"Setup",
			rollupSetupStatus(page),
			setupLaunchpadSignal(page),
			[]string{"Run the local app or reference deployment setup.", "Review private checklist rows before changing deployment settings."},
			[]string{"docs/tutorials/agency-first-run.md", "docs/tutorials/self-hosted-operator-trial.md", "docs/tutorials/agency-launchpad.md"},
			[]string{"make agency-app-up", "make deployment-doctor"},
			[]string{"/admin/operations/setup", "/admin/operations/checklist"},
			"Setup status is private operator diagnostics, not public launch or deployment approval.",
		),
		launchpadSection(
			"gtfs",
			"GTFS",
			gtfsLaunchpadStatus(page),
			gtfsLaunchpadSignal(page),
			[]string{"Import a GTFS ZIP through the browser or CLI, or publish a GTFS Studio draft.", "Use GTFS quality triage for validator and importer findings."},
			[]string{"docs/tutorials/real-agency-gtfs-onboarding.md", "docs/tutorials/gtfs-validation-triage.md"},
			[]string{"make agency-pilot-up", "go run ./cmd/gtfs-import --help"},
			[]string{"/admin/operations/gtfs-import", "/admin/gtfs-studio", "/admin/operations/gtfs-quality"},
			"GTFS publication signals are not agency approval or CAL-ITP/Caltrans compliance.",
		),
		launchpadSection(
			"metadata",
			"Metadata",
			metadataLaunchpadStatus(page),
			metadataLaunchpadSignal(page),
			[]string{"Replace placeholder metadata with operator-entered license, contact, and URL values.", "Keep separate authorization records outside this page if stronger claims are ever needed."},
			[]string{"docs/compliance-evidence-checklist.md", "docs/tutorials/calitp-readiness-checklist.md"},
			[]string{"make audit-final-claim-review"},
			[]string{"/admin/operations/setup", "/admin/operations/readiness"},
			"Metadata completeness is a readiness signal only and is not final-root proof or approval.",
		),
		launchpadSection(
			"five_feeds",
			"Five feeds",
			fiveFeedsLaunchpadStatus(page),
			fiveFeedsLaunchpadSignal(page),
			[]string{"Confirm feeds.json plus schedule, Vehicle Positions, Trip Updates, and Alerts are listed.", "Use the plain-language feed health dashboard before relying on any public feed path."},
			[]string{"docs/tutorials/operator-smoke-and-support-bundle.md", "docs/requirements-calitp-compliance.md"},
			[]string{"make operator-smoke", "make validator-health"},
			[]string{"/admin/operations/feed-health", "/admin/operations/feeds", "/admin/operations/validation-health"},
			"Feed availability is not consumer acceptance, correctness proof, or compliance proof.",
		),
		launchpadSection(
			"telemetry",
			"Telemetry",
			telemetryLaunchpadStatus(page),
			telemetryLaunchpadSignal(page),
			[]string{"Bind or rotate device credentials.", "Send synthetic telemetry through the authenticated ingest path before reviewing Vehicle Positions."},
			[]string{"docs/tutorials/telemetry-simulator-and-device-trial.md", "docs/tutorials/device-avl-integration.md"},
			[]string{"make telemetry-simulator"},
			[]string{"/admin/operations/devices", "/admin/operations/telemetry-simulator", "/admin/operations/telemetry"},
			"Telemetry diagnostics do not prove vendor compatibility or production AVL reliability.",
		),
		launchpadSection(
			"validators",
			"Validators",
			validatorLaunchpadStatus(page),
			validatorLaunchpadSignal(page),
			[]string{"Install or check pinned validators.", "Run allowlisted validation from private validator health or CLI checks."},
			[]string{"docs/release-candidate-readiness.md", "docs/tutorials/gtfs-validation-triage.md"},
			[]string{"make validators-check", "make validator-health"},
			[]string{"/admin/operations/validation-health", "/admin/operations/gtfs-quality"},
			"Validator output is a supporting signal only, not compliance, consumer acceptance, or correctness proof.",
		),
		launchpadSection(
			"readiness",
			"Readiness",
			readinessLaunchpadStatus(page),
			readinessLaunchpadSignal(page),
			[]string{"Review readiness rows and private checklist rows.", "Keep missing data missing until source records exist."},
			[]string{"docs/tutorials/calitp-readiness-checklist.md", "docs/california-readiness-summary.md"},
			[]string{"make release-candidate-check", "make audit-final-claim-review"},
			[]string{"/admin/operations/readiness", "/admin/operations/checklist"},
			"Readiness review does not claim CAL-ITP/Caltrans compliance, public launch, or production readiness.",
		),
		launchpadSection(
			"connector_conformance",
			"Connector conformance",
			checklistStatusNeedsReview,
			"sidecar manifests, synthetic conformance suite, and generic examples are available for local checks",
			[]string{"Validate manifests for any optional sidecar connector.", "Run synthetic conformance before connecting external systems."},
			[]string{"docs/connectors/plugin-contract.md", "docs/tutorials/external-adapter-conformance.md", "docs/integration-adapter-kit.md"},
			[]string{"make external-connection-check", "make adapter-conformance"},
			[]string{"/admin/operations/connectors", "/admin/operations/readiness"},
			"Connector conformance is local contract quality only, not vendor compatibility or external acceptance.",
		),
		launchpadSection(
			"support_bundle",
			"Support bundle",
			checklistStatusNeedsReview,
			"operator smoke, deployment doctor, reliability, and support-bundle workflows are available as private diagnostics",
			[]string{"Run smoke and support-bundle tooling only in the operator environment.", "Review redaction before sharing any artifact outside the deployment team."},
			[]string{"docs/tutorials/operator-smoke-and-support-bundle.md", "docs/evidence/redaction-policy.md", "docs/deployment/reference-deployment-doctor.md"},
			[]string{"make operator-smoke", "make support-bundle", "make deployment-doctor"},
			[]string{"/admin/operations/reliability", "/admin/operations/evidence"},
			"Support bundles are private diagnostics unless a separate evidence intake approves retention and disclosure.",
		),
		launchpadSection(
			"decision_gate",
			"Decision gate",
			checklistStatusNeedsReview,
			"claim flags remain false and no approval or compliance outcome is recorded by this launchpad",
			[]string{"Decide whether to continue local hardening, pause for missing inputs, or request a separate evidence intake.", "Do not move consumer statuses or publish stronger wording from this page."},
			[]string{"docs/post-60-product-roadmap.md", "docs/open-questions.md", "docs/phase-60-final-claim-review-and-public-closeout.md"},
			[]string{"make audit-final-claim-review"},
			[]string{"/admin/operations/checklist.json", "/admin/operations/launchpad.json"},
			"The decision gate records no approval outcome, compliance outcome, public launch outcome, or production readiness outcome.",
		),
	}
	return agencyLaunchpadView{
		GeneratedAt: page.GeneratedAt,
		AgencyID:    page.AgencyID,
		Boundary:    "Private authenticated operator workflow only; viewing it creates no evidence, contacts no external party, changes no consumer status, and records no approval or compliance outcome.",
		Sections:    sections,
		Counts:      launchpadCounts(sections),
		ClaimFlags:  agencyLaunchpadClaimFlags{},
		DecisionNotes: []agencyLaunchpadDecision{
			{Label: "Continue", Status: checklistStatusNeedsReview, CurrentSignal: "product quality and external connection maturity work can continue inside normal repo checks", NextAction: "Use normal development checkpoints and baseline checks.", Boundary: "Continue is an internal maintainer choice, not public launch or production readiness."},
			{Label: "Pause", Status: checklistStatusNeedsReview, CurrentSignal: "missing authorization, public-safe retention, or real deployment inputs should stop evidence-like work", NextAction: "Record blockers in docs or open questions without collecting retained artifacts.", Boundary: "Pause does not create evidence or imply failure of a public obligation."},
			{Label: "Evidence intake", Status: checklistStatusBlocked, CurrentSignal: "no intake exists in this launchpad", NextAction: "Define authorization, retention, redaction, claim target, allowed tools, and stop conditions before any evidence track.", Boundary: "No intake means no evidence phase."},
		},
	}
}

func launchpadSection(id string, label string, status string, signal string, nextActions []string, docsLinks []string, commandSuggestions []string, adminLinks []string, boundary string) agencyLaunchpadSection {
	return agencyLaunchpadSection{
		ID:                 id,
		Label:              label,
		Status:             normalizeChecklistStatus(status),
		CurrentSignal:      firstNonEmpty(signal, "unknown"),
		NextActions:        cleanLaunchpadList(nextActions),
		DocsLinks:          safeDocsLinks(docsLinks),
		CommandSuggestions: safeCommandSuggestions(commandSuggestions),
		AdminLinks:         safeAdminLinks(adminLinks),
		ClaimBoundary:      firstNonEmpty(boundary, privateBoundary()),
	}
}

func launchpadCounts(sections []agencyLaunchpadSection) agencyLaunchpadCounts {
	counts := agencyLaunchpadCounts{Sections: len(sections), Statuses: map[string]int{
		checklistStatusOK:          0,
		checklistStatusNeedsReview: 0,
		checklistStatusMissing:     0,
		checklistStatusBlocked:     0,
		checklistStatusUnknown:     0,
	}}
	for _, section := range sections {
		counts.Statuses[normalizeChecklistStatus(section.Status)]++
	}
	return counts
}

func cleanLaunchpadList(values []string) []string {
	var out []string
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" || unsafePrivateString(trimmed) {
			continue
		}
		out = append(out, trimmed)
	}
	if len(out) == 0 {
		return []string{"Review the linked private Operations Console section."}
	}
	return out
}

func safeCommandSuggestions(values []string) []string {
	var out []string
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" || unsafePrivateString(trimmed) {
			continue
		}
		if strings.Contains(trimmed, "curl ") || strings.Contains(trimmed, "http://") || strings.Contains(trimmed, "https://") || strings.Contains(trimmed, ">") {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}

func safeAdminLinks(values []string) []string {
	var out []string
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" || !strings.HasPrefix(trimmed, "/admin/") || strings.Contains(trimmed, "..") || unsafePrivateString(trimmed) {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}

func rollupSetupStatus(page operationsPage) string {
	var statuses []string
	for _, step := range page.SetupSteps {
		statuses = append(statuses, step.Status)
	}
	return rollupStatuses(statuses)
}

func setupLaunchpadSignal(page operationsPage) string {
	if page.PublicationError != "" || page.DiscoveryError != "" {
		return firstNonEmpty(page.PublicationError, page.DiscoveryError)
	}
	return "publication setup, private checklist, and guided setup pages are available"
}

func gtfsLaunchpadStatus(page operationsPage) string {
	if strings.TrimSpace(page.ActiveFeedVersion) == "" {
		return checklistStatusMissing
	}
	return checklistStatusNeedsReview
}

func gtfsLaunchpadSignal(page operationsPage) string {
	if strings.TrimSpace(page.ActiveFeedVersion) == "" {
		return "missing active GTFS feed version"
	}
	return "active GTFS feed version is present"
}

func metadataLaunchpadStatus(page operationsPage) string {
	status := metadataStatus(page.Discovery.AgencyName, page.Discovery.License.Name, page.Discovery.License.URL, page.Discovery.TechnicalContactEmail)
	if status == checklistStatusOK {
		return checklistStatusNeedsReview
	}
	return status
}

func metadataLaunchpadSignal(page operationsPage) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	return "agency name, license, contact, and public root metadata are reviewed as private readiness signals"
}

func fiveFeedsLaunchpadStatus(page operationsPage) string {
	if page.DiscoveryError != "" {
		return checklistStatusMissing
	}
	requiredFeeds := map[string]bool{"schedule": false, "vehicle_positions": false, "trip_updates": false, "alerts": false}
	for _, feed := range page.Discovery.Feeds {
		if _, ok := requiredFeeds[feed.FeedType]; ok && strings.TrimSpace(feed.CanonicalPublicURL) != "" {
			requiredFeeds[feed.FeedType] = true
		}
	}
	if strings.TrimSpace(page.Discovery.PublicBaseURL) == "" {
		return checklistStatusMissing
	}
	for _, ok := range requiredFeeds {
		if !ok {
			return checklistStatusMissing
		}
	}
	if !page.Discovery.Readiness.HTTPSURLs {
		return checklistStatusNeedsReview
	}
	return checklistStatusOK
}

func fiveFeedsLaunchpadSignal(page operationsPage) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	return "feeds.json plus schedule, Vehicle Positions, Trip Updates, and Alerts metadata are checked from feed discovery"
}

func telemetryLaunchpadStatus(page operationsPage) string {
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

func telemetryLaunchpadSignal(page operationsPage) string {
	if page.TelemetryError != "" {
		return page.TelemetryError
	}
	if len(page.Telemetry) == 0 {
		return "no accepted telemetry observed"
	}
	return "latest telemetry rows are visible in the private Operations Console"
}

func validatorLaunchpadStatus(page operationsPage) string {
	switch page.ValidationHealth.OverallStatus {
	case "blocked", "failed", "misconfigured_tooling":
		return checklistStatusBlocked
	case "missing_tooling", "artifact_unavailable", "not_run":
		return checklistStatusMissing
	case "stale", "needs_review", "recorded", "runnable", "configured", "installed", "stub", "configured_for_tests", "skipped":
		return checklistStatusNeedsReview
	default:
		return checklistStatusUnknown
	}
}

func validatorLaunchpadSignal(page operationsPage) string {
	if page.ValidationHealthError != "" {
		return page.ValidationHealthError
	}
	return "validator health overall " + firstNonEmpty(page.ValidationHealth.OverallStatus, checklistStatusUnknown) + " with tooling " + firstNonEmpty(page.ValidationHealth.ToolingStatus, checklistStatusUnknown)
}

func readinessLaunchpadStatus(page operationsPage) string {
	var statuses []string
	for _, item := range page.ReadinessItems {
		statuses = append(statuses, item.Status)
	}
	return rollupStatuses(statuses)
}

func readinessLaunchpadSignal(page operationsPage) string {
	if len(page.ReadinessItems) == 0 {
		return "readiness rows are not available"
	}
	return "private readiness workflow has current rows and next actions"
}

func rollupStatuses(statuses []string) string {
	if len(statuses) == 0 {
		return checklistStatusUnknown
	}
	seen := map[string]bool{}
	for _, status := range statuses {
		seen[normalizeChecklistStatus(status)] = true
	}
	switch {
	case seen[checklistStatusBlocked]:
		return checklistStatusBlocked
	case seen[checklistStatusMissing]:
		return checklistStatusMissing
	case seen[checklistStatusNeedsReview]:
		return checklistStatusNeedsReview
	case seen[checklistStatusUnknown]:
		return checklistStatusUnknown
	default:
		return checklistStatusOK
	}
}
