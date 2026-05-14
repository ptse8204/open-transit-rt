package main

import (
	"net/http"
	"os"
	"time"

	"open-transit-rt/internal/auth"
)

type operationsHelpView struct {
	GeneratedAt    time.Time                     `json:"generated_at"`
	AgencyID       string                        `json:"agency_id"`
	Boundary       string                        `json:"boundary"`
	RoleTours      []operationsHelpRoleTour      `json:"role_tours"`
	FirstWeek      []operationsHelpFirstWeekItem `json:"first_week_checklist"`
	Topics         []operationsHelpTopic         `json:"topics"`
	ContextualHelp operationsContextHelp         `json:"contextual_help"`
	ClaimFlags     operationsHelpClaimFlags      `json:"claim_flags"`
}

type operationsHelpRoleTour struct {
	ID           string   `json:"id"`
	Label        string   `json:"label"`
	Who          string   `json:"who"`
	StartHere    string   `json:"start_here"`
	ReviewFirst  string   `json:"review_first"`
	FirstActions string   `json:"first_actions"`
	EscalateWhen string   `json:"escalate_when"`
	DoesNotProve string   `json:"does_not_prove"`
	AdminLinks   []string `json:"admin_links"`
	DocsLinks    []string `json:"docs_links"`
}

type operationsHelpFirstWeekItem struct {
	ID           string `json:"id"`
	Day          string `json:"day"`
	Role         string `json:"role"`
	Task         string `json:"task"`
	Review       string `json:"review"`
	DoneWhen     string `json:"done_when"`
	NextAction   string `json:"next_action"`
	ConsoleLink  string `json:"console_link"`
	DoesNotProve string `json:"does_not_prove"`
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
		Title:            operationsPageTitle(section),
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
		RoleTours:      operationsHelpRoleTours(),
		FirstWeek:      operationsHelpFirstWeekChecklist(),
		Topics:         topics,
		ContextualHelp: context,
		ClaimFlags:     operationsHelpClaimFlags{},
	}
}

func operationsHelpRoleTours() []operationsHelpRoleTour {
	return []operationsHelpRoleTour{
		helpRoleTour(
			"no_developer_evaluator",
			"No-developer evaluator",
			"Someone deciding whether the private browser workflow is understandable before asking for technical help.",
			"/admin/operations",
			"Start Here, Agency Setup, Feed Links & Health, Validation Center, Help",
			"Follow status chips and next actions only; do not try shell commands or deployment changes.",
			"Ask a technical helper when startup, validator tooling, public root configuration, or support bundle creation is needed.",
			"Does not prove agency approval, outside review, public launch, or release readiness.",
			[]string{"/admin/operations", "/admin/operations/setup-wizard", "/admin/operations/validation-center", "/admin/operations/help"},
			[]string{"docs/wiki/Small-Agency-Quick-Start.md", "docs/tutorials/browser-first-setup.md"},
		),
		helpRoleTour(
			"director_manager",
			"Director or manager",
			"Someone checking whether the agency has a clear operational path and knows which work needs staff or technical support.",
			"/admin/operations/launchpad",
			"Private Launchpad, Readiness, Consumers, Maintenance, Evidence",
			"Confirm the team knows the next private task, current blockers, and which outside claims remain unavailable.",
			"Ask for technical help when a blocker needs data import, feed validation tooling, deployment checks, or redaction review.",
			"Does not prove adoption, approval, compliance, public launch, managed service, or support coverage.",
			[]string{"/admin/operations/launchpad", "/admin/operations/readiness", "/admin/operations/consumers", "/admin/operations/maintenance"},
			[]string{"docs/agency-training-outline.md", "docs/support-boundaries.md"},
		),
		helpRoleTour(
			"daily_operator",
			"Daily operator",
			"Someone responsible for checking schedule, realtime, alerts, feed health, and routine maintenance state.",
			"/admin/operations/realtime",
			"Realtime Center, Feed Links & Health, Validation Center, Alerts, Maintenance",
			"Review stale, missing, blocked, and needs-review rows; prefer unknown over false certainty.",
			"Ask a technical helper when telemetry stops, validators are missing, backups need setup, or feed paths are not configured.",
			"Does not prove field reliability, target display, ETA quality, uptime, or service-level guarantees.",
			[]string{"/admin/operations/realtime", "/admin/operations/feeds", "/admin/operations/validation-center", "/admin/operations/maintenance"},
			[]string{"docs/tutorials/operator-smoke-and-support-bundle.md", "docs/runbooks/small-agency-pilot-operations.md"},
		),
		helpRoleTour(
			"technical_helper",
			"Technical helper",
			"Someone who can run local setup, validators, support bundle review, device setup, connector dry-runs, and deployment diagnostics.",
			"/admin/operations/setup",
			"Advanced Setup Details, GTFS Workbench, Validation Health, Devices, Connectors, Maintenance",
			"Use existing private admin routes and documented commands only; keep secrets and raw outputs out of browser guidance.",
			"Escalate to a maintainer before evidence retention, release packaging, portal work, public-root proof, or schema-changing work.",
			"Does not prove production readiness, vendor fit, hardware certification, release readiness, or external acceptance.",
			[]string{"/admin/operations/setup", "/admin/operations/gtfs-workbench", "/admin/operations/validation-health", "/admin/operations/connectors/workbench"},
			[]string{"docs/dependencies.md", "docs/validator-tooling.md", "docs/integration-adapter-kit.md"},
		),
		helpRoleTour(
			"integrator",
			"Integrator",
			"Someone evaluating telemetry, connector, monitoring, validator, or prediction sidecar boundaries using synthetic or local inputs.",
			"/admin/operations/connectors/workbench",
			"Connector Workbench, Connector Tests, Telemetry Simulator, Prediction & ETA Lab, Realtime Center",
			"Use synthetic fixtures and fail-closed adapter patterns before any real external integration is authorized.",
			"Ask for separate written authorization before real vendor/device payloads, credentials, or network sends are introduced.",
			"Does not prove named vendor compatibility, equipment certification, realtime reliability, or ETA accuracy.",
			[]string{"/admin/operations/connectors/workbench", "/admin/operations/connectors/tests", "/admin/operations/telemetry-simulator", "/admin/operations/prediction-lab"},
			[]string{"docs/integration-adapter-kit.md", "docs/tutorials/external-adapter-conformance.md"},
		),
	}
}

func operationsHelpFirstWeekChecklist() []operationsHelpFirstWeekItem {
	return []operationsHelpFirstWeekItem{
		helpFirstWeekItem("day_1_start", "Day 1", "No-developer evaluator", "Open Start Here and identify the first blocked or missing setup item.", "Agency scope, setup progress, first-run tasks, help panel", "The next private task is clear to a nontechnical reviewer.", "/admin/operations", "Does not prove the deployment is ready outside local review."),
		helpFirstWeekItem("day_1_setup", "Day 1", "Director or manager", "Review Agency Setup and assign owners for profile, feed metadata, schedule, telemetry, validation, and maintenance.", "Setup Wizard stages and role visibility", "Each setup area has a named staff owner or technical-helper owner.", "/admin/operations/setup-wizard", "Does not prove agency approval or outside readiness."),
		helpFirstWeekItem("day_2_gtfs", "Day 2", "Technical helper", "Import or review GTFS using the existing private schedule workflow.", "GTFS import result, Workbench, quality triage, required file checklist", "Schedule state is visible with clear next actions for missing or invalid records.", "/admin/operations/gtfs-workbench", "Does not prove validator-clean status or public source-of-truth listing."),
		helpFirstWeekItem("day_3_feeds", "Day 3", "Daily operator", "Check feed URLs, validation center, and plain-language feed health.", "Configured five feed URLs, validator context, local fetch context, blockers", "Missing or stale feed rows have an owner and next action.", "/admin/operations/feeds", "Does not prove final-root readiness, consumer review, or public launch."),
		helpFirstWeekItem("day_4_realtime", "Day 4", "Daily operator", "Review telemetry, devices, Vehicle Positions, Trip Updates, and Alerts together.", "Realtime Center, devices, telemetry freshness, Alerts lifecycle", "Unknown, stale, unmatched, and degraded states are understood before sharing realtime status.", "/admin/operations/realtime", "Does not prove production-grade ETA quality, field reliability, or target display."),
		helpFirstWeekItem("day_5_connectors", "Day 5", "Integrator", "Review connector recipes, conformance guidance, simulator paths, and prediction lab boundaries.", "Connector Workbench, Connector Tests, Telemetry Simulator, Prediction Lab", "Synthetic/local connector next steps are clear and no real credentials are needed.", "/admin/operations/connectors/workbench", "Does not prove vendor compatibility, hardware certification, or real device proof."),
		helpFirstWeekItem("day_5_maintenance", "Day 5", "Technical helper", "Review Maintenance Center, support bundle guidance, redaction warnings, and handoff needs.", "Backup/restore review, upgrade/rollback review, support bundle guidance, cadence rows", "A support and maintenance checklist exists without running destructive browser actions.", "/admin/operations/maintenance", "Does not prove hosted-service availability, SLA, uptime, release readiness, or support coverage."),
	}
}

func helpRoleTour(id string, label string, who string, startHere string, reviewFirst string, firstActions string, escalate string, doesNotProve string, adminLinks []string, docsLinks []string) operationsHelpRoleTour {
	return operationsHelpRoleTour{
		ID:           id,
		Label:        label,
		Who:          who,
		StartHere:    helpAdminLink(startHere),
		ReviewFirst:  reviewFirst,
		FirstActions: firstActions,
		EscalateWhen: escalate,
		DoesNotProve: doesNotProve,
		AdminLinks:   safeAdminLinks(adminLinks),
		DocsLinks:    safeDocsLinks(docsLinks),
	}
}

func helpFirstWeekItem(id string, day string, role string, task string, review string, doneWhen string, consoleLink string, doesNotProve string) operationsHelpFirstWeekItem {
	return operationsHelpFirstWeekItem{
		ID:           id,
		Day:          day,
		Role:         role,
		Task:         task,
		Review:       review,
		DoneWhen:     doneWhen,
		NextAction:   "Open the linked private console page and record only local operator notes unless a separate evidence gate is authorized.",
		ConsoleLink:  helpAdminLink(consoleLink),
		DoesNotProve: doesNotProve,
	}
}

func helpAdminLink(value string) string {
	links := safeAdminLinks([]string{value})
	if len(links) == 0 {
		return "/admin/operations"
	}
	return links[0]
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
		return "Start Here"
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
		return "Telemetry Freshness"
	case "telemetry-simulator":
		return "Telemetry Simulator"
	case "devices":
		return "Device Credentials"
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
		return "Operations Console Help"
	default:
		return "this section"
	}
}
