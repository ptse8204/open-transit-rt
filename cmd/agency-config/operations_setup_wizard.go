package main

import (
	"net/http"
	"time"

	"open-transit-rt/internal/auth"
)

type operationsSetupWizardView struct {
	GeneratedAt time.Time                    `json:"generated_at"`
	AgencyID    string                       `json:"agency_id"`
	Boundary    string                       `json:"boundary"`
	Stages      []operationsSetupWizardStage `json:"stages"`
	Counts      operationsSetupWizardCounts  `json:"counts"`
	ClaimFlags  setupWizardClaimFlags        `json:"claim_flags"`
}

type operationsSetupWizardStage struct {
	ID            string   `json:"id"`
	Label         string   `json:"label"`
	Status        string   `json:"status"`
	CurrentSignal string   `json:"current_signal"`
	PrimaryAction string   `json:"primary_action"`
	AdminLink     string   `json:"admin_link"`
	DocsLinks     []string `json:"docs_links"`
	ClaimBoundary string   `json:"claim_boundary"`
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
			"Publication metadata",
			metadataStatus(page.Discovery.License.Name, page.Discovery.License.URL, page.Discovery.TechnicalContactEmail),
			licenseContactEvidence(page),
			"Store or review public base URL, feed base URL, open license, contact, and environment values.",
			"/admin/operations/setup#publication-metadata",
			[]string{"docs/compliance-evidence-checklist.md", "docs/tutorials/calitp-readiness-checklist.md"},
			"Publication metadata completeness is not final-root proof, compliance proof, or agency approval.",
		),
		setupWizardStage(
			"gtfs",
			"GTFS",
			gtfsLaunchpadStatus(page),
			gtfsLaunchpadSignal(page),
			"Use browser import, the CLI import path, or GTFS Studio draft publish; then review validation feedback.",
			"/admin/operations/gtfs-import",
			[]string{"docs/tutorials/real-agency-gtfs-onboarding.md", "docs/tutorials/gtfs-validation-triage.md"},
			"GTFS setup status does not claim validator-clean data, public launch, or Caltrans/CAL-ITP compliance.",
		),
		setupWizardStage(
			"feeds",
			"Feeds",
			fiveFeedsLaunchpadStatus(page),
			fiveFeedsLaunchpadSignal(page),
			"Review plain-language health for feeds.json, schedule, Vehicle Positions, Trip Updates, and Alerts without changing public routes.",
			"/admin/operations/feed-health",
			[]string{"docs/requirements-calitp-compliance.md", "docs/tutorials/operator-smoke-and-support-bundle.md"},
			"Listed feed URLs are readiness signals, not consumer acceptance or correctness proof.",
		),
		setupWizardStage(
			"telemetry",
			"Telemetry",
			telemetryLaunchpadStatus(page),
			telemetryLaunchpadSignal(page),
			"Bind devices and send authenticated sample telemetry through existing ingest flows.",
			"/admin/operations/telemetry",
			[]string{"docs/tutorials/telemetry-simulator-and-device-trial.md", "docs/tutorials/device-avl-integration.md"},
			"Telemetry visibility does not prove real vendor compatibility, hardware certification, or fleet reliability.",
		),
		setupWizardStage(
			"validators",
			"Validators",
			validatorLaunchpadStatus(page),
			validatorLaunchpadSignal(page),
			"Review private validator health and run only server-side allowlisted validators from existing admin paths.",
			"/admin/operations/validation-health",
			[]string{"docs/release-candidate-readiness.md", "docs/tutorials/gtfs-validation-triage.md"},
			"Validator records are supporting diagnostics only, not compliance proof or consumer acceptance.",
		),
		setupWizardStage(
			"connectors",
			"Connectors",
			checklistStatusNeedsReview,
			"connector hub describes sidecar, manifest, command-adapter, and conformance paths without dynamic backend plugin loading",
			"Review connector boundaries before connecting optional external systems.",
			"/admin/operations/connectors",
			[]string{"docs/connectors/plugin-contract.md", "docs/tutorials/external-adapter-conformance.md", "docs/integration-adapter-kit.md"},
			"Connector setup does not claim vendor compatibility, external acceptance, hosted SaaS, or production ETA quality.",
		),
		setupWizardStage(
			"readiness",
			"Readiness",
			readinessLaunchpadStatus(page),
			readinessLaunchpadSignal(page),
			"Review readiness rows and private checklist rows, keeping missing evidence marked missing.",
			"/admin/operations/readiness",
			[]string{"docs/tutorials/calitp-readiness-checklist.md", "docs/california-readiness-summary.md"},
			"Readiness review does not claim CAL-ITP/Caltrans compliance, public launch, or production readiness.",
		),
	}
	return operationsSetupWizardView{
		GeneratedAt: page.GeneratedAt,
		AgencyID:    page.AgencyID,
		Boundary:    "Private authenticated setup wizard only; viewing it creates no evidence, changes no state, contacts no external party, opens no public route, and records no approval, compliance, public launch, hosted-service, vendor, SLA, or production-readiness outcome.",
		Stages:      stages,
		Counts:      setupWizardCounts(stages),
		ClaimFlags:  setupWizardClaimFlags{},
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
