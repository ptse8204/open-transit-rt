package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
)

type operationsHelpView struct {
	GeneratedAt     time.Time                     `json:"generated_at"`
	AgencyID        string                        `json:"agency_id"`
	Boundary        string                        `json:"boundary"`
	RoleTours       []operationsHelpRoleTour      `json:"role_tours"`
	FirstWeek       []operationsHelpFirstWeekItem `json:"first_week_checklist"`
	Glossary        []operationsHelpGlossaryTerm  `json:"glossary"`
	Recovery        []operationsHelpRecoveryRow   `json:"recovery_guidance"`
	TrainingGuide   operationsHelpTrainingGuide   `json:"training_guide"`
	QuickTasks      []operationsHelpQuickTask     `json:"quick_tasks"`
	Handoff         []operationsHelpHandoffItem   `json:"handoff_checklist"`
	DemoScenarios   []operationsHelpDemoScenario  `json:"demo_scenarios"`
	TrainerScript   []operationsHelpTrainerStep   `json:"trainer_script"`
	HelperChecklist []operationsHelpHelperItem    `json:"administrator_checklist"`
	Topics          []operationsHelpTopic         `json:"topics"`
	ContextualHelp  operationsContextHelp         `json:"contextual_help"`
	ClaimFlags      operationsHelpClaimFlags      `json:"claim_flags"`
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

type operationsHelpGlossaryTerm struct {
	ID               string   `json:"id"`
	Term             string   `json:"term"`
	PlainMeaning     string   `json:"plain_meaning"`
	TechnicalMeaning string   `json:"technical_meaning"`
	WhereToReview    string   `json:"where_to_review"`
	DocsLinks        []string `json:"docs_links"`
	DoesNotProve     string   `json:"does_not_prove"`
}

type operationsHelpRecoveryRow struct {
	ID                string `json:"id"`
	WhatOperatorSees  string `json:"what_operator_sees"`
	LikelyCause       string `json:"likely_cause"`
	SafeNextStep      string `json:"safe_next_step"`
	EscalationTrigger string `json:"escalation_trigger"`
	ConsoleLink       string `json:"console_link"`
	DoesNotProve      string `json:"does_not_prove"`
}

type operationsHelpTrainingGuide struct {
	Label    string `json:"label"`
	DocsPath string `json:"docs_path"`
	Audience string `json:"audience"`
	HowToUse string `json:"how_to_use"`
	Boundary string `json:"boundary"`
}

type operationsHelpQuickTask struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	PrimaryRole  string `json:"primary_role"`
	ConsoleLink  string `json:"console_link"`
	ReviewSteps  string `json:"review_steps"`
	DoneWhen     string `json:"done_when"`
	Escalation   string `json:"escalation"`
	DoesNotProve string `json:"does_not_prove"`
}

type operationsHelpHandoffItem struct {
	ID           string `json:"id"`
	Area         string `json:"area"`
	FromRole     string `json:"from_role"`
	ToRole       string `json:"to_role"`
	Confirm      string `json:"confirm"`
	ConsoleLink  string `json:"console_link"`
	DoesNotProve string `json:"does_not_prove"`
}

type operationsHelpDemoScenario struct {
	ID             string   `json:"id"`
	Label          string   `json:"label"`
	Audience       string   `json:"audience"`
	FixturePaths   []string `json:"fixture_paths"`
	ConsoleLink    string   `json:"console_link"`
	Exercise       string   `json:"exercise"`
	DoneWhen       string   `json:"done_when"`
	RecoveryPrompt string   `json:"recovery_prompt"`
	DoesNotProve   string   `json:"does_not_prove"`
}

type operationsHelpTrainerStep struct {
	ID             string `json:"id"`
	Segment        string `json:"segment"`
	Minutes        string `json:"minutes"`
	TalkTrack      string `json:"talk_track"`
	AskParticipant string `json:"ask_participant"`
	ConsoleLink    string `json:"console_link"`
	Boundary       string `json:"boundary"`
}

type operationsHelpHelperItem struct {
	ID                     string `json:"id"`
	Area                   string `json:"area"`
	Collect                string `json:"collect"`
	DoNotCollect           string `json:"do_not_collect"`
	ConsoleLink            string `json:"console_link"`
	DocsLink               string `json:"docs_link"`
	NeedsAuthorizationWhen string `json:"needs_authorization_when"`
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
		GeneratedAt:     generatedAt,
		AgencyID:        agencyID,
		Boundary:        "Private authenticated Operations Console help only. Viewing it is read-only guidance: it creates no evidence, contacts no outside party, changes no consumer status, and records no approval or outside outcome.",
		RoleTours:       operationsHelpRoleTours(),
		FirstWeek:       operationsHelpFirstWeekChecklist(),
		Glossary:        operationsHelpGlossary(),
		Recovery:        operationsHelpRecoveryGuidance(),
		TrainingGuide:   operationsHelpTrainingGuideLink(),
		QuickTasks:      operationsHelpQuickTasks(),
		Handoff:         operationsHelpHandoffChecklist(),
		DemoScenarios:   operationsHelpDemoScenarios(),
		TrainerScript:   operationsHelpTrainerScript(),
		HelperChecklist: operationsHelpHelperChecklist(),
		Topics:          topics,
		ContextualHelp:  context,
		ClaimFlags:      operationsHelpClaimFlags{},
	}
}

func operationsHelpRoleTours() []operationsHelpRoleTour {
	return []operationsHelpRoleTour{
		helpRoleTour(
			"no_code_evaluator",
			"No-code evaluator",
			"Someone deciding whether the private browser workflow is understandable before asking for technical help.",
			"/admin/operations",
			"Start, Agency Setup, Feed Links & Health, Validation Center, Help",
			"Follow status chips and next actions only; do not try shell commands or deployment changes.",
			"Ask an administrator when startup, validator tooling, public root configuration, or support bundle creation is needed.",
			"This private view does not show agency approval, outside review, public launch, or release readiness.",
			[]string{"/admin/operations", "/admin/operations/setup-wizard", "/admin/operations/validation-center", "/admin/operations/help"},
			[]string{"docs/wiki/Small-Agency-Quick-Start.md", "docs/tutorials/browser-first-setup.md"},
		),
		helpRoleTour(
			"director_manager",
			"Director or manager",
			"Someone checking whether the agency has a clear operational path and knows which work needs staff or technical support.",
			"/admin/operations/launchpad",
			"Agency Launchpad, Readiness, Consumer Preparation, Maintenance, Evidence Guidance",
			"Confirm the team knows the next private task, current blockers, and which outside claims remain unavailable.",
			"Ask for technical help when a blocker needs data import, feed validation tooling, deployment checks, or redaction review.",
			"This private view does not show adoption, approval, compliance, public launch, managed service, or support coverage.",
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
			"Ask an administrator when telemetry stops, validators are missing, backups need setup, or feed paths are not configured.",
			"This private view does not show field reliability, target display, ETA quality, uptime, or service-level guarantees.",
			[]string{"/admin/operations/realtime", "/admin/operations/feeds", "/admin/operations/validation-center", "/admin/operations/maintenance"},
			[]string{"docs/tutorials/operator-smoke-and-support-bundle.md", "docs/runbooks/small-agency-pilot-operations.md"},
		),
		helpRoleTour(
			"administrator",
			"Administrator",
			"Someone who can run local setup, validators, support bundle review, device setup, connector dry-runs, and deployment diagnostics.",
			"/admin/operations/setup",
			"Setup Details, Schedule Review, Validation Health, Devices, Connectors, Maintenance",
			"Use existing private admin routes and documented commands only; keep secrets and raw outputs out of browser guidance.",
			"Escalate to a maintainer before evidence retention, release packaging, portal work, public-root proof, or schema-changing work.",
			"This private view does not show production readiness, vendor fit, hardware certification, release readiness, or external acceptance.",
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
			"This private view does not show named vendor compatibility, equipment certification, realtime reliability, or ETA accuracy.",
			[]string{"/admin/operations/connectors/workbench", "/admin/operations/connectors/tests", "/admin/operations/telemetry-simulator", "/admin/operations/prediction-lab"},
			[]string{"docs/integration-adapter-kit.md", "docs/tutorials/external-adapter-conformance.md"},
		),
	}
}

func operationsHelpFirstWeekChecklist() []operationsHelpFirstWeekItem {
	return []operationsHelpFirstWeekItem{
		helpFirstWeekItem("day_1_start", "Day 1", "No-code evaluator", "Open Start and identify the first blocked or missing setup item.", "Agency scope, setup progress, first-run tasks, help panel", "The next private task is clear to a nontechnical reviewer.", "/admin/operations", "This private view does not show the deployment is ready outside local review."),
		helpFirstWeekItem("day_1_setup", "Day 1", "Director or manager", "Review Agency Setup and assign owners for profile, feed metadata, schedule, telemetry, validation, and maintenance.", "Setup Wizard stages and role visibility", "Each setup area has a named staff owner or deployment owner.", "/admin/operations/setup-wizard", "This private view does not show agency approval or outside readiness."),
		helpFirstWeekItem("day_2_gtfs", "Day 2", "Administrator", "Import or review GTFS using the existing private schedule workflow.", "GTFS import result, Workbench, quality triage, required file checklist", "Schedule state is visible with clear next actions for missing or invalid records.", "/admin/operations/gtfs-workbench", "This private view does not show validator-clean status or public source-of-truth listing."),
		helpFirstWeekItem("day_3_feeds", "Day 3", "Daily operator", "Check feed URLs, validation center, and plain-language feed health.", "Configured five feed URLs, validator context, local fetch context, blockers", "Missing or stale feed rows have an owner and next action.", "/admin/operations/feeds", "This private view does not show final-root readiness, consumer review, or public launch."),
		helpFirstWeekItem("day_4_realtime", "Day 4", "Daily operator", "Review telemetry, devices, Vehicle Positions, Trip Updates, and Alerts together.", "Realtime Center, devices, telemetry freshness, Alerts lifecycle", "Unknown, stale, unmatched, and degraded states are understood before sharing realtime status.", "/admin/operations/realtime", "This private view does not show production-grade ETA quality, field reliability, or target display."),
		helpFirstWeekItem("day_5_connectors", "Day 5", "Integrator", "Review connector recipes, conformance guidance, simulator paths, and prediction lab boundaries.", "Connector Workbench, Connector Tests, Telemetry Simulator, Trip Updates", "Synthetic/local connector next steps are clear and no real credentials are needed.", "/admin/operations/connectors/workbench", "This private view does not show vendor compatibility, hardware certification, or real device proof."),
		helpFirstWeekItem("day_5_maintenance", "Day 5", "Administrator", "Review Maintenance Center, support bundle guidance, redaction warnings, and handoff needs.", "Backup/restore review, upgrade/rollback review, support bundle guidance, cadence rows", "A support and maintenance checklist exists without running destructive browser actions.", "/admin/operations/maintenance", "This private view does not show hosted-service availability, SLA, uptime, release readiness, or support coverage."),
	}
}

func operationsHelpGlossary() []operationsHelpGlossaryTerm {
	return []operationsHelpGlossaryTerm{
		helpGlossaryTerm("gtfs", "GTFS", "The schedule package: routes, stops, trips, calendars, and agency metadata.", "A static feed ZIP with files such as agency, stops, routes, trips, stop times, calendar, and shapes.", "Schedule Review, Import Schedule ZIP, Schedule Quality", []string{"docs/requirements-2a-2f.md", "docs/tutorials/real-agency-gtfs-onboarding.md"}, "A loaded schedule does not show agency approval, source-of-truth listing, or validator-clean status."),
		helpGlossaryTerm("gtfs_rt", "GTFS Realtime", "The live feed family that describes vehicles, predicted trip changes, and alerts.", "Protocol buffer feeds for Vehicle Positions, Trip Updates, and Alerts that reference the active GTFS schedule.", "Realtime Center, Feed Links & Health, Validation Center", []string{"docs/requirements-trip-updates.md", "docs/requirements-calitp-compliance.md"}, "A visible realtime feed does not show consumer display, public launch, or realtime reliability."),
		helpGlossaryTerm("vehicle_positions", "Vehicle Positions", "Where vehicles are observed right now or recently.", "GTFS-RT entities emitted from authenticated telemetry after conservative matching and freshness checks.", "Realtime Center, Telemetry, Devices", []string{"docs/requirements-2a-2f.md"}, "Fresh Vehicle Positions do not prove field reliability, equipment certification, or consumer use."),
		helpGlossaryTerm("trip_updates", "Trip Updates", "Predicted schedule changes or arrivals when the system has enough confidence.", "GTFS-RT Trip Updates remain behind a prediction adapter and may be withheld when data is stale, ambiguous, or low confidence.", "Prediction & ETA Lab, Realtime Center", []string{"docs/requirements-trip-updates.md"}, "Trip Updates diagnostics do not prove production-grade ETA quality or real-world ETA accuracy."),
		helpGlossaryTerm("alerts", "Alerts", "Service messages such as delays, detours, or disruptions.", "GTFS-RT Alerts are authored and published separately from telemetry ingest and prediction.", "Alerts Console, Realtime Center", []string{"docs/requirements-calitp-compliance.md"}, "Alerts availability does not show agency approval, consumer display, or disruption-workflow completeness."),
		helpGlossaryTerm("telemetry", "Telemetry", "Vehicle observations sent from a device, simulator, adapter, or integration.", "Authenticated events posted to the telemetry ingest API and used for freshness, matching, and Vehicle Positions output.", "Telemetry, Devices, Telemetry Simulator", []string{"docs/tutorials/device-avl-integration.md", "docs/tutorials/telemetry-simulator-and-device-trial.md"}, "Telemetry presence does not show real-device reliability, vendor support, or field operations quality."),
		helpGlossaryTerm("validators", "Validators", "Tools that check schedule and realtime feed format and quality.", "Server-owned allowlisted validator IDs and off-host guidance summarize static GTFS and GTFS-RT diagnostic results.", "Validation Health, Validation Center, Schedule Quality", []string{"docs/validator-tooling.md", "docs/tutorials/gtfs-validation-triage.md"}, "Validator rows are supporting signals only; they do not prove compliance or outside review."),
		helpGlossaryTerm("readiness", "Readiness", "A private checklist of what looks missing, blocked, stale, or ready for local review.", "Aggregated private signals from metadata, feeds, validation, telemetry, reliability, maintenance, and prepared packet records.", "Readiness, Checklist, Validation Center", []string{"docs/tutorials/calitp-readiness-checklist.md", "docs/california-readiness-summary.md"}, "Readiness supports review workflows but does not show compliance, final-root status, or public launch."),
		helpGlossaryTerm("evidence", "Evidence", "Retained proof used only when a separately authorized claim needs it.", "Public-safe, redacted, claim-specific artifacts controlled by evidence gates and protected path rules.", "Evidence, Consumers, Feed Links & Health", []string{"docs/evidence/redaction-policy.md", "docs/evidence/evidence-track-router.md"}, "Help text and diagnostics are not evidence packets and do not move consumer status."),
		helpGlossaryTerm("connectors", "Connectors", "Ways to connect telemetry, predictors, validators, monitoring, or external processes safely.", "Optional sidecars, manifests, examples, or adapters around stable Open Transit RT contracts, using synthetic/local checks first.", "Connectors, Connector Workbench, Connector Tests", []string{"docs/integration-adapter-kit.md", "docs/connectors/plugin-contract.md"}, "Connector checks do not prove named vendor compatibility or hardware certification."),
		helpGlossaryTerm("support_bundle", "Support bundle", "A redacted package of diagnostic context for troubleshooting.", "A guided support artifact should avoid secrets, raw private output, credentials, database URLs, and unredacted operator data.", "Maintenance, Help", []string{"docs/tutorials/operator-smoke-and-support-bundle.md", "docs/support-boundaries.md"}, "A support bundle does not show managed support, SLA, uptime, release readiness, or public launch."),
	}
}

func operationsHelpRecoveryGuidance() []operationsHelpRecoveryRow {
	return []operationsHelpRecoveryRow{
		helpRecoveryRow("empty_start_here", "Start shows missing or blocked rows.", "The agency has not completed setup, import, telemetry, validation, or maintenance inputs yet.", "Open Agency Setup and follow the first missing private next action.", "An administrator is needed if startup, environment, validator tooling, or deployment settings are missing.", "/admin/operations/setup-wizard", "Clearing a missing row does not show outside readiness."),
		helpRecoveryRow("gtfs_import_confusing", "A schedule import produced errors or the Workbench still looks incomplete.", "The source GTFS may be missing required files, has validator issues, or has not been published as the active feed.", "Open Schedule Review and Schedule Quality, then fix source data before relying on realtime feeds.", "Escalate when the source owner must change service data or a validator tool is unavailable.", "/admin/operations/gtfs-workbench", "Import review does not show validator-clean status or source-of-truth listing."),
		helpRecoveryRow("feed_url_missing", "A feed URL is missing, local-only, or not ready to copy.", "Publication metadata, active schedule, realtime feed output, or public root configuration is incomplete.", "Open Feed Links & Health and review metadata, local fetch context, and source-of-truth guidance.", "Use technical help for deployment root configuration or off-host validation.", "/admin/operations/feeds", "A configured URL does not show final-root readiness or consumer action."),
		helpRecoveryRow("validator_blocked", "Validation Health says missing, stale, blocked, or tool unavailable.", "The host may lack pinned validator tooling, an artifact may be missing, or the active feed changed after the latest run.", "Use Validation Center and Validation Health to identify the specific feed and validator row before rerunning allowed checks.", "Escalate if validator tooling must be installed or off-host validation is needed.", "/admin/operations/validation-center", "Validator success remains a supporting signal only."),
		helpRecoveryRow("vehicles_stale", "Vehicles are stale, unmatched, degraded, or missing.", "Telemetry may have stopped, device binding may be inactive, GPS quality may be poor, or matching confidence may be intentionally unknown.", "Open Realtime Center, Devices, and Telemetry to find freshness, binding, and confidence signals.", "Escalate if devices, tokens, or external telemetry adapters need operational changes.", "/admin/operations/realtime", "Resolving staleness does not show field reliability or target display."),
		helpRecoveryRow("trip_updates_withheld", "Trip Updates are empty, fallback, or withheld.", "The deterministic predictor may be protecting riders from stale, ambiguous, low-confidence, or missing assignment data.", "Open Prediction & ETA Lab and read withheld reasons before changing prediction settings.", "Escalate before introducing external predictors or real-world ETA claims.", "/admin/operations/prediction-lab", "Withheld diagnostics do not prove production-grade ETA quality."),
		helpRecoveryRow("connector_unclear", "A connector recipe is unclear or a dry-run does not match expectations.", "The input shape, adapter boundary, fixture, or monitoring expectation may not match the local synthetic contract.", "Open Connector Workbench and Connector Tests, then choose the recipe closest to the actual integration shape.", "Escalate before real credentials, real vendor payloads, or network sends are introduced.", "/admin/operations/connectors/workbench", "Synthetic conformance does not show vendor compatibility or hardware certification."),
		helpRecoveryRow("consumer_status_confusion", "A prepared consumer packet is visible and someone asks whether it can be submitted.", "Prepared packet records are local readiness references only and remain separate from external contact or target workflow status.", "Open Consumers and Feed Links & Health to review prepared-only boundaries and future authorization gates.", "Escalate only when a separate written authorization starts a consumer or evidence track.", "/admin/operations/consumers", "Prepared packet visibility does not show submission, review, listing, display, ingestion, or acceptance."),
	}
}

func operationsHelpTrainingGuideLink() operationsHelpTrainingGuide {
	return operationsHelpTrainingGuide{
		Label:    "Printable Staff Training Guide",
		DocsPath: "docs/operator-training-guide.md",
		Audience: "Directors, daily operators, administrators, integrators, and no-code evaluators.",
		HowToUse: "Use the guide for staff onboarding, weekly review, troubleshooting practice, and handoff notes before opening stronger evidence or release gates.",
		Boundary: "The guide is training material only. It does not create evidence, move consumer status, authorize external contact, or prove public outcomes.",
	}
}

func operationsHelpQuickTasks() []operationsHelpQuickTask {
	return []operationsHelpQuickTask{
		helpQuickTask("import_gtfs", "Import or review GTFS", "Administrator", "/admin/operations/gtfs-workbench", "Open Schedule Review, inspect active feed state, import history, required files, quality triage, and validation context.", "Schedule state and next action are visible to the operator.", "Escalate when source GTFS data must be corrected or validator tooling is unavailable.", "This private view does not show validator-clean status, source-of-truth listing, agency approval, or compliance."),
		helpQuickTask("check_feed_health", "Check feed health", "Daily operator", "/admin/operations/feeds", "Review five configured feed URLs, metadata, local fetch context, validator context, and off-host guidance.", "Every missing, stale, blocked, or needs-review feed row has an owner or next action.", "Escalate when public root configuration or off-host validation is needed.", "This private view does not show final-root readiness, consumer action, public launch, SLA, or uptime."),
		helpQuickTask("run_validation", "Review validation", "Administrator", "/admin/operations/validation-center", "Use Validation Center and Validation Health to identify the feed, validator, artifact, latest status, and blocker.", "Validation status is understood and the next safe rerun path is clear.", "Escalate before changing validator tooling, server paths, or deployment configuration.", "This private view does not show compliance, outside review, or consumer status."),
		helpQuickTask("simulate_telemetry", "Add or simulate telemetry", "Daily operator", "/admin/operations/telemetry-simulator", "Review device binding, simulator scenario, credential boundary, telemetry freshness, and Realtime Center results.", "Synthetic or local telemetry review has a clear result without exposing credentials.", "Escalate if live devices, tokens, or adapter payloads need operational changes.", "This private view does not show real-device reliability, vendor compatibility, or field operations quality."),
		helpQuickTask("review_connectors", "Review connectors", "Integrator", "/admin/operations/connectors/workbench", "Choose a recipe, inspect synthetic normalization preview, dry-run guidance, conformance results, and fail-closed boundaries.", "The next connector experiment uses local or synthetic inputs only.", "Escalate before real credentials, real vendor payloads, or network sends are introduced.", "This private view does not show named vendor compatibility or hardware certification."),
		helpQuickTask("prepare_administrator", "Prepare for an administrator", "Director or manager", "/admin/operations/help", "Collect the current private page, blocker, owner, and intended next action without copying secrets or raw private output.", "The administrator knows which console section and docs page to inspect first.", "Escalate to maintainer review before evidence retention, release packaging, or schema-changing work.", "This private view does not show managed support, release readiness, or outside approval."),
		helpQuickTask("support_bundle", "Prepare support bundle guidance", "Administrator", "/admin/operations/maintenance", "Open Maintenance, review support-bundle/redaction guidance, backup/restore review, cadence rows, and infrastructure category checks.", "Support context is ready for redacted local review without browser-executed destructive actions.", "Escalate before sharing private artifacts outside the operator-approved channel.", "This private view does not show paid support, SLA, uptime, hosted-service availability, or release readiness."),
	}
}

func operationsHelpHandoffChecklist() []operationsHelpHandoffItem {
	return []operationsHelpHandoffItem{
		helpHandoffItem("setup_owner", "Agency setup", "Director or manager", "Administrator", "Agency profile, public feed metadata, contact, license, and role visibility have an owner.", "/admin/operations/setup-wizard", "This private view does not show agency approval, legal approval, or public readiness."),
		helpHandoffItem("schedule_owner", "Schedule data", "Administrator", "Daily operator", "Active schedule state, import history, quality triage, and validation context are understandable.", "/admin/operations/gtfs-workbench", "This private view does not show validator-clean status or source-of-truth listing."),
		helpHandoffItem("realtime_owner", "Realtime operations", "Daily operator", "Administrator", "Telemetry freshness, device binding, Vehicle Positions, Trip Updates, and Alerts states have known next actions.", "/admin/operations/realtime", "This private view does not show field reliability, target display, or ETA quality."),
		helpHandoffItem("connector_owner", "Connector experiments", "Integrator", "Administrator", "Recipes, synthetic fixtures, conformance expectations, and no-send boundaries are documented for the next local experiment.", "/admin/operations/connectors/workbench", "This private view does not show vendor compatibility, hardware certification, or live integration fitness."),
		helpHandoffItem("maintenance_owner", "Maintenance", "Administrator", "Director or manager", "Backup/restore review, upgrade/rollback review, support-bundle guidance, and cadence rows have owners.", "/admin/operations/maintenance", "This private view does not show hosted-service availability, SLA, uptime, managed support, or release readiness."),
		helpHandoffItem("claim_owner", "Claims and evidence", "Director or manager", "Maintainer or authorized evidence owner", "Prepared packet, final-root, consumer, agency pilot, vendor/device, ETA-quality, and compliance gates stay separate.", "/admin/operations/consumers", "Does not create evidence, authorize external contact, move consumer status, or prove public outcomes."),
	}
}

func operationsHelpDemoScenarios() []operationsHelpDemoScenario {
	return []operationsHelpDemoScenario{
		helpDemoScenario("baseline_start", "Baseline local agency walkthrough", "No-code evaluator and trainer", []string{"testdata/gtfs/valid-small", "testdata/telemetry-simulator/on-route.json"}, "/admin/operations", "Use Start to explain schedule import, feed URLs, realtime status, and help boundaries from a clean local demo.", "Participant can name the next private action and what the demo does not show.", "If the page is empty, use Agency Setup and keep missing rows missing until local inputs exist.", "This private view does not show agency adoption, public launch, final-root readiness, or consumer display."),
		helpDemoScenario("after_midnight", "After-midnight service drill", "Daily operator and administrator", []string{"testdata/gtfs/after-midnight", "testdata/telemetry-simulator/after-midnight.json", "testdata/replay/after-midnight-service.json"}, "/admin/operations/realtime", "Review how service-day language, stale thresholds, and conservative matching behave for late-night trips.", "Operator can explain why agency-local service day matters before relying on realtime output.", "If a vehicle looks unmatched, inspect trip candidates and avoid inventing a trip descriptor.", "This private view does not show real-world overnight reliability or ETA accuracy."),
		helpDemoScenario("frequency_service", "Frequency-based service drill", "Administrator and integrator", []string{"testdata/gtfs/frequency-based", "testdata/telemetry/frequency-based.json", "testdata/replay/frequency-exact-window.json"}, "/admin/operations/prediction-lab", "Use frequency fixtures to discuss headways, repeated trip instances, and withheld Trip Updates when confidence is insufficient.", "Participant can identify why unknown is safer than false certainty.", "If output is withheld, read the Trip Updates diagnostics before changing thresholds.", "This private view does not show production-grade ETA quality or real-world headway prediction accuracy."),
		helpDemoScenario("stale_unknown_device", "Stale and unknown-device recovery drill", "Daily operator", []string{"testdata/telemetry-simulator/stale.json", "testdata/telemetry-simulator/unknown-device.json", "testdata/replay/stale-telemetry-withheld.json"}, "/admin/operations/devices", "Practice tracing stale telemetry and unknown devices through Devices, Telemetry, and Realtime Center.", "Operator can decide when to escalate device binding or token work to an administrator.", "If a token or binding change is needed, keep token values out of notes and issue trackers.", "This private view does not show device fleet reliability, vendor compatibility, or hardware certification."),
		helpDemoScenario("alerts_disruption", "Disruption and alert lifecycle drill", "Daily operator and director", []string{"testdata/replay/cancellation-alert-linkage.json", "testdata/replay/disruption-diagnostics-baseline.json"}, "/admin/alerts/console", "Review how canceled-trip hints, alert lifecycle, and public Alerts usefulness guidance fit together.", "Participant can separate an operator alert draft from public launch or consumer acceptance.", "If a real disruption requires publication, use existing alert workflow and agency policy outside this training drill.", "This private view does not show public alert display, consumer ingestion, or agency approval."),
		helpDemoScenario("connector_conformance", "Connector conformance tabletop", "Integrator and administrator", []string{"testdata/adapter-conformance/suite.json", "testdata/connectors/valid/synthetic-telemetry-source-input.json", "testdata/connectors/valid/synthetic-prediction-input.json"}, "/admin/operations/connectors/workbench", "Map a proposed CSV/API/POST/prediction/monitoring integration to synthetic adapter boundaries before any real credentials exist.", "Participant can choose the safest local recipe and name what authorization would be needed next.", "If a vendor payload is proposed, stop and request separate written authorization and redaction review.", "This private view does not show named vendor compatibility, real device proof, or hardware certification."),
	}
}

func operationsHelpTrainerScript() []operationsHelpTrainerStep {
	return []operationsHelpTrainerStep{
		helpTrainerStep("open_boundary", "Opening boundary", "5", "State that this is private training over local or synthetic data only.", "Ask each participant to name one thing the session must not prove.", "/admin/operations/help", "Creates no evidence, no outside contact, and no consumer status change."),
		helpTrainerStep("role_path", "Role path assignment", "10", "Assign each participant one role path: evaluator, director, operator, administrator, or integrator.", "Ask the participant to open the role's first console page and read the first next action.", "/admin/operations/help", "Does not create support entitlement, adoption proof, or approval."),
		helpTrainerStep("demo_scenario", "Scenario walkthrough", "20", "Choose one demo scenario and follow only the listed console and fixture references.", "Ask what a blocked or missing row means and when to escalate.", "/admin/operations/help", "Uses synthetic/local fixtures only and does not show field outcomes."),
		helpTrainerStep("recovery_drill", "Recovery drill", "15", "Pick one common mistake row and rehearse the safe first step before any technical action.", "Ask what not to copy into notes or tickets.", "/admin/operations/help", "Does not authorize credential sharing, private data retention, or external contact."),
		helpTrainerStep("handoff", "Administrator handoff", "10", "Complete a handoff using page, blocker, owner, intended next action, and authorization need.", "Ask whether the handoff needs a maintainer, administrator, or separate evidence gate.", "/admin/operations/help", "Does not start an evidence, release, vendor, or consumer workflow."),
		helpTrainerStep("closeout", "Closeout", "5", "Review claim boundaries and leave statuses unchanged unless a real authorized workflow exists.", "Ask participants to repeat that all consumer targets remain prepared.", "/admin/operations/consumers", "This private view does not show adoption, consumer acceptance, compliance, public launch, hosted service, SLA, or production readiness."),
	}
}

func operationsHelpHelperChecklist() []operationsHelpHelperItem {
	return []operationsHelpHelperItem{
		helpHelperItem("startup_context", "Startup context", "App URL, command attempted, visible status, and local-only environment label.", "Secrets, tokens, database URLs, private logs, or screenshots with credentials.", "/admin/operations/setup", "docs/tutorials/agency-first-run.md", "Real deployment, public-root proof, or hosted-service claims are requested."),
		helpHelperItem("schedule_context", "Schedule context", "GTFS source type, import status, active feed version, validator row, and fixture path when demo-only.", "Real agency private schedule files unless explicitly authorized and redacted.", "/admin/operations/gtfs-workbench", "docs/tutorials/gtfs-validation-triage.md", "Source-of-truth listing, agency approval, or compliance proof is requested."),
		helpHelperItem("realtime_context", "Realtime context", "Vehicle ID, stale/unmatched state, simulator scenario, device binding status, and withheld reason.", "Device tokens, raw private AVL payloads, or unredacted vendor logs.", "/admin/operations/realtime", "docs/tutorials/telemetry-simulator-and-device-trial.md", "Real device proof, vendor proof, or ETA-quality claims are requested."),
		helpHelperItem("connector_context", "Connector context", "Chosen recipe, synthetic fixture, conformance case, expected input shape, and no-send posture.", "Real credentials, portal tokens, private endpoint URLs, or vendor payload samples.", "/admin/operations/connectors/workbench", "docs/tutorials/external-adapter-conformance.md", "Live vendor/device integration or network sends are proposed."),
		helpHelperItem("monitoring_context", "Monitoring context", "Health digest row, monitoring export row, no-send channel state, and local summary path if already produced.", "Webhook URLs, email credentials, destination values, or private alert recipient lists.", "/admin/operations/maintenance", "docs/tutorials/self-hosted-operations-notifications.md", "Live notification delivery, hosted monitoring, SLA, or uptime claims are requested."),
		helpHelperItem("evidence_context", "Evidence and claims", "Which claim is requested, current prepared-only status, required gate, and missing retained proof.", "Protected evidence paths, private portal records, or target-originated artifacts unless the evidence gate authorizes them.", "/admin/operations/evidence", "docs/evidence/evidence-track-router.md", "Any evidence retention, consumer submission, or status movement is requested."),
		helpHelperItem("release_context", "Release or package context", "Current local diagnostic status, release note draft path, blocker, and exact command that was run.", "Signing keys, registry tokens, private release credentials, or publication secrets.", "/admin/operations/readiness", "docs/release-candidate-readiness.md", "Tagging, package publication, image push, GitHub Release creation, or release-readiness claims are requested."),
	}
}

func helpDemoScenario(id string, label string, audience string, fixturePaths []string, consoleLink string, exercise string, doneWhen string, recoveryPrompt string, doesNotProve string) operationsHelpDemoScenario {
	return operationsHelpDemoScenario{
		ID:             id,
		Label:          label,
		Audience:       audience,
		FixturePaths:   safeFixtureLinks(fixturePaths),
		ConsoleLink:    helpAdminLink(consoleLink),
		Exercise:       exercise,
		DoneWhen:       doneWhen,
		RecoveryPrompt: recoveryPrompt,
		DoesNotProve:   doesNotProve,
	}
}

func helpTrainerStep(id string, segment string, minutes string, talkTrack string, askParticipant string, consoleLink string, boundary string) operationsHelpTrainerStep {
	return operationsHelpTrainerStep{
		ID:             id,
		Segment:        segment,
		Minutes:        minutes,
		TalkTrack:      talkTrack,
		AskParticipant: askParticipant,
		ConsoleLink:    helpAdminLink(consoleLink),
		Boundary:       boundary,
	}
}

func helpHelperItem(id string, area string, collect string, doNotCollect string, consoleLink string, docsLink string, needsAuthorizationWhen string) operationsHelpHelperItem {
	links := safeDocsLinks([]string{docsLink})
	if len(links) == 0 {
		links = []string{"docs/operator-training-guide.md"}
	}
	return operationsHelpHelperItem{
		ID:                     id,
		Area:                   area,
		Collect:                collect,
		DoNotCollect:           doNotCollect,
		ConsoleLink:            helpAdminLink(consoleLink),
		DocsLink:               links[0],
		NeedsAuthorizationWhen: needsAuthorizationWhen,
	}
}

func safeFixtureLinks(values []string) []string {
	var out []string
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" || !strings.HasPrefix(trimmed, "testdata/") || strings.Contains(trimmed, "..") || unsafePrivateString(trimmed) {
			continue
		}
		out = append(out, trimmed)
	}
	return out
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

func helpGlossaryTerm(id string, term string, plain string, technical string, where string, docsLinks []string, doesNotProve string) operationsHelpGlossaryTerm {
	return operationsHelpGlossaryTerm{
		ID:               id,
		Term:             term,
		PlainMeaning:     plain,
		TechnicalMeaning: technical,
		WhereToReview:    where,
		DocsLinks:        safeDocsLinks(docsLinks),
		DoesNotProve:     doesNotProve,
	}
}

func helpRecoveryRow(id string, sees string, cause string, next string, escalation string, consoleLink string, doesNotProve string) operationsHelpRecoveryRow {
	return operationsHelpRecoveryRow{
		ID:                id,
		WhatOperatorSees:  sees,
		LikelyCause:       cause,
		SafeNextStep:      next,
		EscalationTrigger: escalation,
		ConsoleLink:       helpAdminLink(consoleLink),
		DoesNotProve:      doesNotProve,
	}
}

func helpQuickTask(id string, label string, role string, consoleLink string, steps string, doneWhen string, escalation string, doesNotProve string) operationsHelpQuickTask {
	return operationsHelpQuickTask{
		ID:           id,
		Label:        label,
		PrimaryRole:  role,
		ConsoleLink:  helpAdminLink(consoleLink),
		ReviewSteps:  steps,
		DoneWhen:     doneWhen,
		Escalation:   escalation,
		DoesNotProve: doesNotProve,
	}
}

func helpHandoffItem(id string, area string, fromRole string, toRole string, confirm string, consoleLink string, doesNotProve string) operationsHelpHandoffItem {
	return operationsHelpHandoffItem{
		ID:           id,
		Area:         area,
		FromRole:     fromRole,
		ToRole:       toRole,
		Confirm:      confirm,
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
			"Realtime feed visibility is a supporting signal only; it does not show outside ingestion, display, launch, or ETA quality.",
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
			"Use Connectors and connector tests to check local contract fit with synthetic inputs before any authorized integration.",
			"Connector readiness does not show named vendor fit, equipment certification, outside acceptance, or deployment outcome.",
			[]string{"/admin/operations/connectors", "/admin/operations/connectors/tests"},
			[]string{"docs/connectors/plugin-contract.md", "docs/integration-adapter-kit.md", "docs/tutorials/external-adapter-conformance.md"},
			"Connector help does not enable arbitrary dynamic backend code loading.",
			safePluginDefinition,
		),
		helpTopic(
			"readiness",
			"Readiness review",
			"Readiness rows translate private feed, metadata, validator, telemetry, reliability, and consumer packet signals into next actions.",
			"Review missing data, stale data, blocked rows, and what each row explicitly does not show.",
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
		return "Start"
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
		return "Connectors"
	case "connector-tests":
		return "connector tests"
	case "telemetry":
		return "Telemetry"
	case "telemetry-simulator":
		return "Telemetry Simulator"
	case "devices":
		return "Devices"
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
		return "Help & Tutorials"
	default:
		return "this section"
	}
}
