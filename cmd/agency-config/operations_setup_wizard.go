package main

import (
	"net/http"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
)

type operationsSetupWizardView struct {
	GeneratedAt    time.Time                             `json:"generated_at"`
	AgencyID       string                                `json:"agency_id"`
	Boundary       string                                `json:"boundary"`
	Summary        operationsSetupWizardSummary          `json:"summary"`
	Blockers       []operationsSetupWizardBlocker        `json:"blockers"`
	Diagnostics    []operationsSetupWizardDiagnostic     `json:"diagnostics"`
	RoleVisibility []operationsSetupWizardRoleVisibility `json:"role_visibility"`
	TechnicalHelp  []operationsSetupWizardTechnicalHelp  `json:"technical_help"`
	Stages         []operationsSetupWizardStage          `json:"stages"`
	Counts         operationsSetupWizardCounts           `json:"counts"`
	ClaimFlags     setupWizardClaimFlags                 `json:"claim_flags"`
}

type operationsSetupWizardSummary struct {
	Status            string `json:"status"`
	CompletedStages   int    `json:"completed_stages"`
	NeedsReviewStages int    `json:"needs_review_stages"`
	MissingStages     int    `json:"missing_stages"`
	BlockedStages     int    `json:"blocked_stages"`
	UnknownStages     int    `json:"unknown_stages"`
	NextStageID       string `json:"next_stage_id"`
	NextStageLabel    string `json:"next_stage_label"`
	NextAction        string `json:"next_action"`
	NextActionLink    string `json:"next_action_link"`
	Meaning           string `json:"meaning"`
}

type operationsSetupWizardStage struct {
	ID            string   `json:"id"`
	Label         string   `json:"label"`
	Status        string   `json:"status"`
	CurrentSignal string   `json:"current_signal"`
	PrimaryAction string   `json:"primary_action"`
	ActionLabel   string   `json:"action_label"`
	AdminLink     string   `json:"admin_link"`
	DocsLinks     []string `json:"docs_links"`
	ClaimBoundary string   `json:"claim_boundary"`
}

type operationsSetupWizardBlocker struct {
	StageID       string `json:"stage_id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	ActionLabel   string `json:"action_label"`
	AdminLink     string `json:"admin_link"`
}

type operationsSetupWizardDiagnostic struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	ClaimBoundary string `json:"claim_boundary"`
}

type operationsSetupWizardRoleVisibility struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	ClaimBoundary string `json:"claim_boundary"`
}

type operationsSetupWizardTechnicalHelp struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	WhenNeeded    string `json:"when_needed"`
	NextAction    string `json:"next_action"`
	AdminLink     string `json:"admin_link"`
	DocsLink      string `json:"docs_link"`
	ClaimBoundary string `json:"claim_boundary"`
}

type operationsSetupWizardCounts struct {
	Stages   int            `json:"stages"`
	Statuses map[string]int `json:"statuses"`
}

type setupWizardClaimFlags struct {
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

func (h *handler) renderSetupWizard(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "setup-wizard")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "setup-wizard", page)
}

func (h *handler) renderSetupWizardJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "setup-wizard")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.SetupWizard)
}

func buildOperationsSetupWizard(page operationsPage) operationsSetupWizardView {
	stages := []operationsSetupWizardStage{
		setupWizardStage(
			"agency_profile",
			"Agency profile",
			metadataStatus(page.Discovery.AgencyName),
			firstNonEmpty(page.Discovery.AgencyName, page.DiscoveryError, "agency profile metadata is missing"),
			"Review agency name and timezone-related GTFS inputs before stronger setup conclusions.",
			"/admin/operations/setup",
			[]string{"docs/tutorials/agency-first-run.md", "docs/tutorials/real-agency-gtfs-onboarding.md"},
			"Agency profile metadata is an operator-entered setup signal, not agency approval.",
		),
		setupWizardStage(
			"publication_metadata",
			"Public feed information",
			metadataStatus(page.Discovery.License.Name, page.Discovery.License.URL, page.Discovery.TechnicalContactEmail),
			licenseContactEvidence(page),
			"Store or review public base URL, feed base URL, open license, contact, and environment values.",
			"/admin/operations/setup#publication-metadata",
			[]string{"docs/compliance-evidence-checklist.md", "docs/tutorials/calitp-readiness-checklist.md"},
			"Publication metadata completeness is not final-root proof, compliance proof, or agency approval.",
		),
		setupWizardStage(
			"gtfs",
			"Schedule data",
			gtfsLaunchpadStatus(page),
			gtfsLaunchpadSignal(page),
			"Use browser import, the CLI import path, or GTFS Studio draft publish; then review validation feedback.",
			"/admin/operations/gtfs-import",
			[]string{"docs/tutorials/real-agency-gtfs-onboarding.md", "docs/tutorials/gtfs-validation-triage.md"},
			"GTFS setup status does not claim validator-clean data, public launch, or Caltrans/CAL-ITP compliance.",
		),
		setupWizardStage(
			"feeds",
			"Feed links",
			fiveFeedsLaunchpadStatus(page),
			fiveFeedsLaunchpadSignal(page),
			"Review plain-language health for feeds.json, schedule, Vehicle Positions, Trip Updates, and Alerts without changing public routes.",
			"/admin/operations/feed-health",
			[]string{"docs/requirements-calitp-compliance.md", "docs/tutorials/operator-smoke-and-support-bundle.md"},
			"Listed feed URLs are readiness signals, not consumer acceptance or correctness proof.",
		),
		setupWizardStage(
			"telemetry",
			"Vehicle telemetry",
			telemetryLaunchpadStatus(page),
			telemetryLaunchpadSignal(page),
			"Bind devices and send authenticated sample telemetry through existing ingest flows.",
			"/admin/operations/telemetry-simulator",
			[]string{"docs/tutorials/telemetry-simulator-and-device-trial.md", "docs/tutorials/device-avl-integration.md"},
			"Telemetry visibility does not show real vendor compatibility, hardware certification, or fleet reliability.",
		),
		setupWizardStage(
			"validators",
			"Validation",
			validatorLaunchpadStatus(page),
			validatorLaunchpadSignal(page),
			"Review private validator health and run only server-side allowlisted validators from existing admin paths.",
			"/admin/operations/validation-health",
			[]string{"docs/release-candidate-readiness.md", "docs/tutorials/gtfs-validation-triage.md"},
			"Validator records are supporting diagnostics only, not compliance proof or consumer acceptance.",
		),
		setupWizardStage(
			"connectors",
			"Optional connectors",
			checklistStatusNeedsReview,
			"connector hub describes sidecar, manifest, command-adapter, and conformance paths without dynamic backend plugin loading",
			"Review connector boundaries before connecting optional external systems.",
			"/admin/operations/connectors",
			[]string{"docs/connectors/plugin-contract.md", "docs/tutorials/external-adapter-conformance.md", "docs/integration-adapter-kit.md"},
			"Connector setup does not claim vendor compatibility, external acceptance, hosted service availability, or production ETA quality.",
		),
		setupWizardStage(
			"readiness",
			"Readiness review",
			readinessLaunchpadStatus(page),
			readinessLaunchpadSignal(page),
			"Review readiness rows and private checklist rows, keeping missing evidence marked missing.",
			"/admin/operations/readiness",
			[]string{"docs/tutorials/calitp-readiness-checklist.md", "docs/california-readiness-summary.md"},
			"Readiness review does not claim CAL-ITP/Caltrans compliance, public launch, or production readiness.",
		),
	}
	return operationsSetupWizardView{
		GeneratedAt:    page.GeneratedAt,
		AgencyID:       page.AgencyID,
		Boundary:       "Private authenticated setup wizard only; viewing it creates no evidence, changes no state, contacts no external party, opens no public route, and records no approval, compliance, public launch, hosted-service, vendor, SLA, or production-readiness outcome.",
		Summary:        setupWizardSummary(stages),
		Blockers:       setupWizardBlockers(stages),
		Diagnostics:    setupWizardDiagnostics(page),
		RoleVisibility: setupWizardRoleVisibility(page),
		TechnicalHelp:  setupWizardTechnicalHelp(),
		Stages:         stages,
		Counts:         setupWizardCounts(stages),
		ClaimFlags:     setupWizardClaimFlags{},
	}
}

func setupWizardStage(id string, label string, status string, signal string, action string, adminLink string, docsLinks []string, boundary string) operationsSetupWizardStage {
	links := safeAdminLinks([]string{adminLink})
	cleanAdminLink := ""
	if len(links) > 0 {
		cleanAdminLink = links[0]
	}
	return operationsSetupWizardStage{
		ID:            id,
		Label:         label,
		Status:        normalizeChecklistStatus(status),
		CurrentSignal: firstNonEmpty(signal, "unknown"),
		PrimaryAction: firstNonEmpty(action, "Review this stage in the private Operations Console."),
		ActionLabel:   setupWizardActionLabel(id),
		AdminLink:     cleanAdminLink,
		DocsLinks:     safeDocsLinks(docsLinks),
		ClaimBoundary: firstNonEmpty(boundary, privateBoundary()),
	}
}

func setupWizardCounts(stages []operationsSetupWizardStage) operationsSetupWizardCounts {
	counts := operationsSetupWizardCounts{Stages: len(stages), Statuses: map[string]int{
		checklistStatusOK:          0,
		checklistStatusNeedsReview: 0,
		checklistStatusMissing:     0,
		checklistStatusBlocked:     0,
		checklistStatusUnknown:     0,
	}}
	for _, stage := range stages {
		counts.Statuses[normalizeChecklistStatus(stage.Status)]++
	}
	return counts
}

func setupWizardSummary(stages []operationsSetupWizardStage) operationsSetupWizardSummary {
	counts := setupWizardCounts(stages)
	next := operationsSetupWizardStage{}
	for _, stage := range stages {
		if stage.Status != checklistStatusOK {
			next = stage
			break
		}
	}
	if next.ID == "" && len(stages) > 0 {
		next = stages[len(stages)-1]
	}
	status := checklistStatusOK
	if counts.Statuses[checklistStatusBlocked] > 0 {
		status = checklistStatusBlocked
	} else if counts.Statuses[checklistStatusMissing] > 0 || counts.Statuses[checklistStatusNeedsReview] > 0 {
		status = checklistStatusNeedsReview
	} else if counts.Statuses[checklistStatusUnknown] > 0 {
		status = checklistStatusUnknown
	}
	return operationsSetupWizardSummary{
		Status:            status,
		CompletedStages:   counts.Statuses[checklistStatusOK],
		NeedsReviewStages: counts.Statuses[checklistStatusNeedsReview],
		MissingStages:     counts.Statuses[checklistStatusMissing],
		BlockedStages:     counts.Statuses[checklistStatusBlocked],
		UnknownStages:     counts.Statuses[checklistStatusUnknown],
		NextStageID:       next.ID,
		NextStageLabel:    next.Label,
		NextAction:        next.PrimaryAction,
		NextActionLink:    next.AdminLink,
		Meaning:           "Setup progress is a private operator guide for local/reference readiness. It does not show approval, compliance, consumer acceptance, hosted operation, or production readiness.",
	}
}

func setupWizardActionLabel(id string) string {
	switch id {
	case "agency_profile":
		return "Review profile"
	case "publication_metadata":
		return "Review feed information"
	case "gtfs":
		return "Open schedule import"
	case "feeds":
		return "Check feed links"
	case "telemetry":
		return "Review telemetry"
	case "validators":
		return "Open validation"
	case "connectors":
		return "Review connectors"
	case "readiness":
		return "Review readiness"
	default:
		return "Open section"
	}
}

func setupWizardBlockers(stages []operationsSetupWizardStage) []operationsSetupWizardBlocker {
	blockers := make([]operationsSetupWizardBlocker, 0)
	for _, stage := range stages {
		if stage.Status == checklistStatusOK {
			continue
		}
		blockers = append(blockers, operationsSetupWizardBlocker{
			StageID:       stage.ID,
			Label:         stage.Label,
			Status:        stage.Status,
			CurrentSignal: stage.CurrentSignal,
			NextAction:    stage.PrimaryAction,
			ActionLabel:   stage.ActionLabel,
			AdminLink:     stage.AdminLink,
		})
	}
	return blockers
}

func setupWizardDiagnostics(page operationsPage) []operationsSetupWizardDiagnostic {
	return []operationsSetupWizardDiagnostic{
		setupWizardDiagnostic(
			"public_base_url",
			"Public base URL",
			presenceStatus(firstNonEmpty(page.PublicationConfig.PublicBaseURL, page.Discovery.PublicBaseURL)),
			presenceSignal(firstNonEmpty(page.PublicationConfig.PublicBaseURL, page.Discovery.PublicBaseURL), "public base URL"),
			"Store a public base URL only after the operator confirms the deployment-owned value.",
			"Configured URL metadata is a setup signal, not final-root proof.",
		),
		setupWizardDiagnostic(
			"feed_base_url",
			"Feed base URL",
			presenceStatus(page.PublicationConfig.FeedBaseURL),
			presenceSignal(page.PublicationConfig.FeedBaseURL, "feed base URL"),
			"Store the feed base URL on the advanced setup page before reviewing public feed links.",
			"Feed URL configuration does not show public launch, consumer acceptance, or production readiness.",
		),
		setupWizardDiagnostic(
			"license_contact",
			"License and contact",
			presenceStatus(firstNonEmpty(page.PublicationConfig.LicenseName, page.Discovery.License.Name), firstNonEmpty(page.PublicationConfig.TechnicalContactEmail, page.Discovery.TechnicalContactEmail)),
			licenseContactDiagnosticSignal(page),
			"Store operator-provided open license and monitored technical contact values, or keep the row marked missing.",
			"Metadata completeness is not retained approval evidence.",
		),
		setupWizardDiagnostic(
			"publication_environment",
			"Publication environment",
			knownPresenceStatus(firstNonEmpty(page.PublicationConfig.PublicationEnvironment, page.Discovery.PublicationEnvironment, page.EnvironmentLabel)),
			knownPresenceSignal(firstNonEmpty(page.PublicationConfig.PublicationEnvironment, page.Discovery.PublicationEnvironment, page.EnvironmentLabel), "publication environment"),
			"Keep local/reference environments distinct from any future release or final-root review.",
			"Environment labels do not prove hosted operation or release readiness.",
		),
		setupWizardDiagnostic(
			"active_schedule",
			"Active schedule",
			presenceStatus(page.ActiveFeedVersion),
			presenceSignal(page.ActiveFeedVersion, "active schedule feed version"),
			"Import a GTFS ZIP, import a safe URL, or publish a GTFS Studio draft before realtime review.",
			"An active schedule does not show validator-clean data or agency approval.",
		),
		setupWizardDiagnostic(
			"validator_tooling",
			"Validator tooling",
			setupWizardValidationToolingStatus(page.ValidationHealth.ToolingStatus),
			firstNonEmpty(page.ValidationHealth.ToolingStatus, "validator tooling status is unknown"),
			"Review Validator Health and run only allowlisted server-side validation actions.",
			"Validator tooling availability does not show compliance or consumer acceptance.",
		),
		setupWizardDiagnostic(
			"device_bindings",
			"Device bindings",
			presenceStatusFromCount(len(page.Devices)),
			deviceBindingDiagnosticSignal(page),
			"Bind devices only through the private token lifecycle; store one-time secrets outside the repo.",
			"Device binding visibility does not show hardware certification or vendor compatibility.",
		),
		setupWizardDiagnostic(
			"telemetry_freshness",
			"Telemetry freshness",
			telemetryStatusForSetup(page),
			telemetryFreshnessDiagnosticSignal(page),
			"Use the telemetry simulator or authenticated device ingest, then review stale and unmatched states.",
			"Telemetry visibility does not show fleet reliability.",
		),
	}
}

func setupWizardRoleVisibility(page operationsPage) []operationsSetupWizardRoleVisibility {
	mutationStatus := checklistStatusBlocked
	mutationSignal := "current role can review setup but cannot submit setup forms"
	mutationNext := "Ask an admin to store publication metadata, run validation, or import GTFS from the browser."
	if page.IsAdmin {
		mutationStatus = checklistStatusOK
		mutationSignal = "current role can review and submit admin-only setup forms"
		mutationNext = "Use the admin-only forms only after confirming values and claim boundaries."
	}
	return []operationsSetupWizardRoleVisibility{
		setupWizardRoleRow(
			"current_roles",
			"Current role visibility",
			checklistStatusOK,
			"roles: "+strings.Join(page.PrincipalRoles, ", "),
			"Review private setup pages for the authenticated agency only.",
			"Role display is not an audit log, support entitlement, or production multi-tenant proof.",
		),
		setupWizardRoleRow(
			"setup_mutations",
			"Setup mutations",
			mutationStatus,
			mutationSignal,
			mutationNext,
			"Browser setup mutations remain private, admin-only, CSRF-protected, and agency-scoped.",
		),
		setupWizardRoleRow(
			"gtfs_import_mutation",
			"GTFS import mutation",
			mutationStatus,
			mutationSignal,
			"Use GTFS Import only when an admin has reviewed source, validation, and rollback implications.",
			"Import permission does not show schedule correctness, public launch, or consumer acceptance.",
		),
	}
}

func setupWizardTechnicalHelp() []operationsSetupWizardTechnicalHelp {
	return []operationsSetupWizardTechnicalHelp{
		setupWizardHelpRow(
			"local_app_startup",
			"Local app startup",
			"The browser cannot reach the private console, validators are missing, or local runtime services are unavailable.",
			"Run the documented local bootstrap and validation commands from an operator terminal.",
			"/admin/operations/maintenance",
			"docs/tutorials/agency-first-run.md",
			"Startup checks are local diagnostics, not hosted SaaS availability.",
		),
		setupWizardHelpRow(
			"gtfs_source_review",
			"GTFS source review",
			"The source ZIP is large, the service period is unclear, rollback is needed, or staged comparison is required.",
			"Use the real-agency GTFS onboarding guide before importing or publishing.",
			"/admin/operations/gtfs-import",
			"docs/tutorials/real-agency-gtfs-onboarding.md",
			"Source review does not create retained evidence or agency approval.",
		),
		setupWizardHelpRow(
			"feed_root_review",
			"Feed root review",
			"Public/feed base URLs, HTTPS behavior, redirects, or final-root ownership need operator confirmation.",
			"Keep final-root evidence work separate until separately authorized.",
			"/admin/operations/feed-health",
			"docs/agency-owned-domain-readiness.md",
			"Feed-root review does not collect final-root evidence.",
		),
		setupWizardHelpRow(
			"support_bundle",
			"Support bundle",
			"An operator needs a redacted diagnostic bundle for maintainer review.",
			"Use maintenance and support-bundle guidance; do not attach secrets, raw private data, or retained evidence without authorization.",
			"/admin/operations/maintenance",
			"docs/tutorials/operator-smoke-and-support-bundle.md",
			"Support guidance does not create support entitlement, SLA, or uptime guarantees.",
		),
	}
}

func setupWizardDiagnostic(id string, label string, status string, signal string, nextAction string, boundary string) operationsSetupWizardDiagnostic {
	return operationsSetupWizardDiagnostic{
		ID:            id,
		Label:         label,
		Status:        normalizeChecklistStatus(status),
		CurrentSignal: firstNonEmpty(signal, "unknown"),
		NextAction:    firstNonEmpty(nextAction, "Review this diagnostic in the private Operations Console."),
		ClaimBoundary: firstNonEmpty(boundary, privateBoundary()),
	}
}

func setupWizardRoleRow(id string, label string, status string, signal string, nextAction string, boundary string) operationsSetupWizardRoleVisibility {
	return operationsSetupWizardRoleVisibility{
		ID:            id,
		Label:         label,
		Status:        normalizeChecklistStatus(status),
		CurrentSignal: firstNonEmpty(signal, "unknown"),
		NextAction:    firstNonEmpty(nextAction, "Review current role permissions."),
		ClaimBoundary: firstNonEmpty(boundary, privateBoundary()),
	}
}

func setupWizardHelpRow(id string, label string, whenNeeded string, nextAction string, adminLink string, docsLink string, boundary string) operationsSetupWizardTechnicalHelp {
	return operationsSetupWizardTechnicalHelp{
		ID:            id,
		Label:         label,
		WhenNeeded:    firstNonEmpty(whenNeeded, "administrator help is needed"),
		NextAction:    firstNonEmpty(nextAction, "Use the linked private console or repo docs."),
		AdminLink:     firstSafeAdminLink(adminLink),
		DocsLink:      firstSafeDocsLink(docsLink),
		ClaimBoundary: firstNonEmpty(boundary, privateBoundary()),
	}
}

func presenceStatus(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return checklistStatusMissing
		}
	}
	return checklistStatusOK
}

func presenceStatusFromCount(count int) string {
	if count == 0 {
		return checklistStatusMissing
	}
	return checklistStatusOK
}

func presenceSignal(value string, label string) string {
	if strings.TrimSpace(value) == "" {
		return label + " is missing"
	}
	return label + " is configured"
}

func knownPresenceStatus(value string) string {
	clean := strings.TrimSpace(value)
	if clean == "" || strings.EqualFold(clean, "unknown") {
		return checklistStatusMissing
	}
	return checklistStatusOK
}

func knownPresenceSignal(value string, label string) string {
	clean := strings.TrimSpace(value)
	if clean == "" || strings.EqualFold(clean, "unknown") {
		return label + " is missing"
	}
	return label + " is configured"
}

func licenseContactDiagnosticSignal(page operationsPage) string {
	licenseConfigured := strings.TrimSpace(firstNonEmpty(page.PublicationConfig.LicenseName, page.Discovery.License.Name)) != ""
	contactConfigured := strings.TrimSpace(firstNonEmpty(page.PublicationConfig.TechnicalContactEmail, page.Discovery.TechnicalContactEmail)) != ""
	switch {
	case licenseConfigured && contactConfigured:
		return "license and technical contact are configured"
	case licenseConfigured:
		return "license is configured; technical contact is missing"
	case contactConfigured:
		return "technical contact is configured; license is missing"
	default:
		return "license and technical contact are missing"
	}
}

func setupWizardValidationToolingStatus(status string) string {
	switch strings.TrimSpace(status) {
	case "configured", "installed", "runnable", "recorded":
		return checklistStatusOK
	case "missing_tooling", "not_run":
		return checklistStatusMissing
	case "blocked", "failed", "misconfigured_tooling":
		return checklistStatusBlocked
	case "stub", "configured_for_tests", "artifact_unavailable", "stale", "needs_review", "skipped":
		return checklistStatusNeedsReview
	default:
		return checklistStatusUnknown
	}
}

func deviceBindingDiagnosticSignal(page operationsPage) string {
	if page.DeviceError != "" {
		return page.DeviceError
	}
	if len(page.Devices) == 0 {
		return "device bindings are missing"
	}
	return "device bindings are configured"
}

func telemetryStatusForSetup(page operationsPage) string {
	if page.TelemetryError != "" {
		return checklistStatusUnknown
	}
	if page.TelemetryUpdatedAt == nil {
		return checklistStatusMissing
	}
	if page.StaleCount > 0 {
		return checklistStatusNeedsReview
	}
	return checklistStatusOK
}

func telemetryFreshnessDiagnosticSignal(page operationsPage) string {
	if page.TelemetryError != "" {
		return page.TelemetryError
	}
	if page.TelemetryUpdatedAt == nil {
		return "telemetry has not been observed yet"
	}
	if page.StaleCount > 0 {
		return "latest telemetry includes stale rows"
	}
	return "latest telemetry is available"
}

func firstSafeAdminLink(link string) string {
	links := safeAdminLinks([]string{link})
	if len(links) == 0 {
		return ""
	}
	return links[0]
}

func firstSafeDocsLink(link string) string {
	links := safeDocsLinks([]string{link})
	if len(links) == 0 {
		return ""
	}
	return links[0]
}
