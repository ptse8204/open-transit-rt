package main

import (
	"net/http"
	"os"
	"time"

	"open-transit-rt/internal/auth"
)

type operationsHelpView struct {
	GeneratedAt    time.Time                `json:"generated_at"`
	AgencyID       string                   `json:"agency_id"`
	Boundary       string                   `json:"boundary"`
	Topics         []operationsHelpTopic    `json:"topics"`
	ContextualHelp operationsContextHelp    `json:"contextual_help"`
	ClaimFlags     operationsHelpClaimFlags `json:"claim_flags"`
}

type operationsHelpTopic struct {
	ID               string   `json:"id"`
	Label            string   `json:"label"`
	Summary          string   `json:"summary"`
	WhatToReview     string   `json:"what_to_review"`
	NextAction       string   `json:"next_action"`
	DoesNotProve     string   `json:"does_not_prove"`
	AdminLinks       []string `json:"admin_links"`
	DocsLinks        []string `json:"docs_links"`
	ClaimBoundary    string   `json:"claim_boundary"`
	PluginDefinition string   `json:"plugin_definition,omitempty"`
}

type operationsContextHelp struct {
	Section      string                `json:"section"`
	Label        string                `json:"label"`
	Topics       []operationsHelpTopic `json:"topics"`
	AllTopicsURL string                `json:"all_topics_url"`
	JSONURL      string                `json:"json_url"`
}

type operationsHelpClaimFlags struct {
	BackendCommandExecutionEnabled     bool `json:"backend_command_execution_enabled"`
	CacheDiagnosticsRead               bool `json:"cache_diagnostics_read"`
	ExternalNetworkContacted           bool `json:"external_network_contacted"`
	ExternalEvidenceCreated            bool `json:"external_evidence_created"`
	FinalRootEvidenceCreated           bool `json:"final_root_evidence_created"`
	ConsumerStatusesChanged            bool `json:"consumer_statuses_changed"`
	SecretsCollected                   bool `json:"secrets_collected"`
	ComplianceClaimed                  bool `json:"compliance_claimed"`
	ProductionReadinessClaimed         bool `json:"production_readiness_claimed"`
	AgencyApprovalClaimed              bool `json:"agency_approval_claimed"`
	ConsumerAcceptanceClaimed          bool `json:"consumer_acceptance_claimed"`
	PublicLaunchClaimed                bool `json:"public_launch_claimed"`
	HostedSaaSClaimed                  bool `json:"hosted_saas_claimed"`
	VendorCompatibilityClaimed         bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed       bool `json:"hardware_certification_claimed"`
	ProductionAVLReliabilityClaimed    bool `json:"production_avl_reliability_claimed"`
	ProductionGradeETAQualityClaimed   bool `json:"production_grade_eta_quality_claimed"`
	SLAClaimed                         bool `json:"sla_claimed"`
	UptimeGuaranteeClaimed             bool `json:"uptime_guarantee_claimed"`
	DynamicBackendPluginLoadingEnabled bool `json:"dynamic_backend_plugin_loading_enabled"`
}

func (h *handler) renderOperationsHelp(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsHelpPage(principal, "help")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "help", page)
}

func (h *handler) renderOperationsHelpJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsHelpPage(principal, "help")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.Help)
}

func (h *handler) buildOperationsHelpPage(principal auth.Principal, section string) operationsPage {
	now := time.Now().UTC().Truncate(time.Second)
	page := operationsPage{
		Title:            "Operations Console",
		AgencyID:         principal.AgencyID,
		GeneratedAt:      now,
		EnvironmentLabel: firstNonEmpty(os.Getenv("PUBLICATION_ENVIRONMENT"), "unknown"),
		CSRFToken:        csrfToken(h.csrfSecret, principal),
		Section:          section,
		NavGroups:        operationsNavGroups(section),
		StaleThreshold:   staleThreshold(),
		IsAdmin:          principal.HasAny(auth.RoleAdmin),
	}
	page.Help = buildOperationsHelpView(now, principal.AgencyID, section)
	page.ContextHelp = page.Help.ContextualHelp
	return page
}

func buildOperationsHelpView(generatedAt time.Time, agencyID string, section string) operationsHelpView {
	topics := operationsHelpTopics()
	context := contextualOperationsHelp(section, topics)
	return operationsHelpView{
		GeneratedAt:    generatedAt,
		AgencyID:       agencyID,
		Boundary:       "Private authenticated Operations Console help only. Viewing it is read-only guidance: it creates no evidence, contacts no outside party, changes no consumer status, and records no approval or outside outcome.",
		Topics:         topics,
		ContextualHelp: context,
		ClaimFlags:     operationsHelpClaimFlags{},
	}
}

func operationsHelpTopics() []operationsHelpTopic {
	return []operationsHelpTopic{
		helpTopic(
			"gtfs",
			"GTFS schedule",
			"GTFS is the static schedule, route, stop, calendar, shape, and trip source that realtime feeds reference.",
			"Review the active published feed version, import result, GTFS Studio draft status, and static validator findings before relying on a schedule.",
			"Import a GTFS ZIP, import a safe URL, or publish a GTFS Studio draft, then review GTFS quality guidance.",
			"Schedule records are readiness inputs only; they do not create agency approval, outside review, or stronger public claims.",
			[]string{"/admin/operations/gtfs-import", "/admin/operations/gtfs-quality", "/admin/gtfs-studio"},
			[]string{"docs/tutorials/real-agency-gtfs-onboarding.md", "docs/tutorials/gtfs-validation-triage.md", "docs/requirements-2a-2f.md"},
			"GTFS help is private operator guidance and does not change published feed state.",
			"",
		),
		helpTopic(
			"gtfs_rt",
			"GTFS Realtime feeds",
			"GTFS-RT covers Vehicle Positions, Trip Updates, and Alerts against the active schedule feed.",
			"Review feed URLs, freshness, protobuf validation context, Trip Updates diagnostics, and Alerts lifecycle state together.",
			"Open feed health, then inspect Vehicle Positions first and Trip Updates or Alerts where configured.",
			"Realtime feed visibility is a supporting signal only; it does not prove outside ingestion, display, launch, or ETA quality.",
			[]string{"/admin/operations/feed-health", "/admin/operations/feeds", "/admin/operations/reliability"},
			[]string{"docs/requirements-trip-updates.md", "docs/requirements-calitp-compliance.md", "docs/tutorials/operator-smoke-and-support-bundle.md"},
			"GTFS-RT help preserves Vehicle Positions first and keeps Trip Updates behind the prediction adapter boundary.",
			"",
		),
		helpTopic(
			"connectors",
			"Connectors and plugins",
			"Connectors are optional sidecars, command adapters, manifests, or connector processes around stable Open Transit RT boundaries.",
			"Review telemetry, prediction, validator, monitoring/export, and consumer/discovery connector paths before connecting external systems.",
			"Use Connector Hub and connector tests to check local contract fit with synthetic inputs before any authorized integration.",
			"Connector readiness does not prove named vendor fit, equipment certification, outside acceptance, or deployment outcome.",
			[]string{"/admin/operations/connectors", "/admin/operations/connectors/tests"},
			[]string{"docs/connectors/plugin-contract.md", "docs/integration-adapter-kit.md", "docs/tutorials/external-adapter-conformance.md"},
			"Connector help does not enable arbitrary dynamic backend code loading.",
			safePluginDefinition,
		),
		helpTopic(
			"readiness",
			"Readiness review",
			"Readiness rows translate private feed, metadata, validator, telemetry, reliability, and consumer packet signals into next actions.",
			"Review missing data, stale data, blocked rows, and what each row explicitly does not prove.",
			"Use the readiness checklist to decide the next private operator action while keeping missing inputs missing.",
			"Readiness review supports CAL-ITP-style workflows, but it is not compliance, approval, or outside consumer status.",
			[]string{"/admin/operations/readiness", "/admin/operations/checklist", "/admin/operations/feed-health"},
			[]string{"docs/tutorials/calitp-readiness-checklist.md", "docs/california-readiness-summary.md", "docs/compliance-evidence-checklist.md"},
			"Readiness help is a private review aid and does not create retained proof.",
			"",
		),
		helpTopic(
			"validators",
			"Validators",
			"Validators check static GTFS and GTFS-RT artifacts through server-owned allowlisted validator IDs.",
			"Review tooling status, artifact availability, latest reports, and GTFS quality triage before rerunning validators.",
			"Run validators only from the existing admin validator health action or existing CLI checks; help itself does not run anything.",
			"Validator rows are supporting signals only; they do not prove outside review or full public-data outcome.",
			[]string{"/admin/operations/validation-health", "/admin/operations/gtfs-quality"},
			[]string{"docs/dependencies.md", "docs/tutorials/gtfs-validation-triage.md", "docs/release-candidate-readiness.md"},
			"Validator help preserves existing validator execution semantics and allowlisted inputs.",
			"",
		),
		helpTopic(
			"telemetry",
			"Telemetry and devices",
			"Telemetry is authenticated vehicle observation input used for conservative matching and Vehicle Positions output.",
			"Review device bindings, latest observations, stale thresholds, assignment confidence, and simulator scenario guidance.",
			"Bind or rotate a device credential, send a synthetic sample if appropriate, then inspect stale or unmatched vehicles.",
			"Fresh telemetry is not proof of field reliability, named equipment fit, agency use, or AVL operating quality.",
			[]string{"/admin/operations/devices", "/admin/operations/telemetry", "/admin/operations/telemetry-simulator"},
			[]string{"docs/tutorials/device-avl-integration.md", "docs/tutorials/telemetry-simulator-and-device-trial.md", "docs/requirements-2a-2f.md"},
			"Telemetry help does not change the authenticated ingest contract.",
			"",
		),
		helpTopic(
			"claims_evidence",
			"Claims and evidence boundaries",
			"Claim boundaries keep private diagnostics, prepared packets, readiness rows, and retained evidence separate.",
			"Review consumer packet status, evidence links, open blockers, and retained-proof requirements before using stronger wording.",
			"Keep consumer targets at prepared unless target-originated retained evidence exists for that exact target and feed scope.",
			"Prepared packets and private diagnostics do not prove submission, review, listing, display, approval, launch, hosted service, or managed operations.",
			[]string{"/admin/operations/consumers", "/admin/operations/evidence", "/admin/operations/readiness"},
			[]string{"docs/evidence/evidence-track-router.md", "docs/evidence/consumer-submissions/README.md", "docs/repo-gaps.md"},
			"Claims help is read-only and does not write evidence or move tracker statuses.",
			"",
		),
	}
}

func helpTopic(id string, label string, summary string, review string, next string, doesNotProve string, adminLinks []string, docsLinks []string, boundary string, pluginDefinition string) operationsHelpTopic {
	return operationsHelpTopic{
		ID:               id,
		Label:            label,
		Summary:          summary,
		WhatToReview:     review,
		NextAction:       next,
		DoesNotProve:     doesNotProve,
		AdminLinks:       safeAdminLinks(adminLinks),
		DocsLinks:        safeDocsLinks(docsLinks),
		ClaimBoundary:    firstNonEmpty(boundary, "Private help only; it does not create evidence or change outside status."),
		PluginDefinition: pluginDefinition,
	}
}

func contextualOperationsHelp(section string, topics []operationsHelpTopic) operationsContextHelp {
	ids := []string{"gtfs", "gtfs_rt", "readiness"}
	switch section {
	case "dashboard", "launchpad":
		ids = []string{"gtfs", "telemetry", "readiness"}
	case "setup", "setup-wizard":
		ids = []string{"gtfs", "readiness", "claims_evidence"}
	case "gtfs-import", "gtfs-quality":
		ids = []string{"gtfs", "validators", "readiness"}
	case "feeds", "feed-health":
		ids = []string{"gtfs_rt", "validators", "readiness"}
	case "connectors", "connector-tests":
		ids = []string{"connectors", "telemetry", "claims_evidence"}
	case "telemetry", "telemetry-simulator", "devices":
		ids = []string{"telemetry", "gtfs_rt", "connectors"}
	case "readiness", "checklist", "reliability":
		ids = []string{"readiness", "validators", "claims_evidence"}
	case "validation-health":
		ids = []string{"validators", "readiness", "claims_evidence"}
	case "consumers", "evidence", "help":
		ids = []string{"claims_evidence", "readiness", "validators"}
	}
	return operationsContextHelp{
		Section:      section,
		Label:        operationsHelpSectionLabel(section),
		Topics:       helpTopicsByID(topics, ids),
		AllTopicsURL: "/admin/operations/help",
		JSONURL:      "/admin/operations/help.json",
	}
}

func helpTopicsByID(topics []operationsHelpTopic, ids []string) []operationsHelpTopic {
	byID := map[string]operationsHelpTopic{}
	for _, topic := range topics {
		byID[topic.ID] = topic
	}
	out := make([]operationsHelpTopic, 0, len(ids))
	for _, id := range ids {
		if topic, ok := byID[id]; ok {
			out = append(out, topic)
		}
	}
	return out
}

func operationsHelpSectionLabel(section string) string {
	switch section {
	case "dashboard":
		return "dashboard"
	case "launchpad":
		return "launchpad"
	case "setup-wizard":
		return "setup wizard"
	case "setup":
		return "setup"
	case "gtfs-import":
		return "GTFS import"
	case "gtfs-quality":
		return "GTFS quality"
	case "feeds":
		return "feeds"
	case "feed-health":
		return "feed health"
	case "connectors":
		return "Connector Hub"
	case "connector-tests":
		return "connector tests"
	case "telemetry":
		return "telemetry"
	case "telemetry-simulator":
		return "telemetry simulator"
	case "devices":
		return "devices"
	case "readiness":
		return "readiness"
	case "checklist":
		return "checklist"
	case "validation-health":
		return "validator health"
	case "reliability":
		return "reliability"
	case "consumers":
		return "consumers"
	case "evidence":
		return "evidence"
	case "help":
		return "help"
	default:
		return "this section"
	}
}
