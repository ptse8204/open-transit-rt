package main

import (
	"fmt"
	"strings"
	"time"
)

type operationsFirstRunView struct {
	GeneratedAt                         time.Time                    `json:"generated_at"`
	AgencyID                            string                       `json:"agency_id"`
	Boundary                            string                       `json:"boundary"`
	LocalDemoDeploymentEvidenceBoundary string                       `json:"local_demo_deployment_evidence_boundary"`
	Paths                               []operationsFirstRunPath     `json:"paths"`
	Tasks                               []operationsFirstRunTask     `json:"tasks"`
	FeedURLs                            []operationsFirstRunFeedURL  `json:"feed_urls"`
	Counts                              operationsFirstRunCounts     `json:"counts"`
	ClaimFlags                          operationsFirstRunClaimFlags `json:"claim_flags"`
}

type operationsFirstRunPath struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	CurrentSignal string `json:"current_signal"`
	Meaning       string `json:"meaning"`
	FirstAction   string `json:"first_action"`
	UILink        string `json:"ui_link"`
	DocsLink      string `json:"docs_link"`
	DoesNotProve  string `json:"does_not_prove"`
}

type operationsFirstRunTask struct {
	Order         int    `json:"order"`
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	Meaning       string `json:"meaning"`
	NextAction    string `json:"next_action"`
	UILink        string `json:"ui_link"`
	DocsLink      string `json:"docs_link"`
	DoesNotProve  string `json:"does_not_prove"`
}

type operationsFirstRunFeedURL struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	Status       string `json:"status"`
	URL          string `json:"url"`
	CopyValue    string `json:"copy_value"`
	Source       string `json:"source"`
	Meaning      string `json:"meaning"`
	NextAction   string `json:"next_action"`
	DocsLink     string `json:"docs_link"`
	DoesNotProve string `json:"does_not_prove"`
}

type operationsFirstRunCounts struct {
	Tasks    int            `json:"tasks"`
	FeedURLs int            `json:"feed_urls"`
	Statuses map[string]int `json:"statuses"`
}

type operationsFirstRunClaimFlags struct {
	BackendCommandExecutionEnabled      bool `json:"backend_command_execution_enabled"`
	CacheDiagnosticsRead                bool `json:"cache_diagnostics_read"`
	ExternalNetworkContacted            bool `json:"external_network_contacted"`
	ExternalEvidenceCreated             bool `json:"external_evidence_created"`
	FinalRootEvidenceCreated            bool `json:"final_root_evidence_created"`
	ConsumerStatusesChanged             bool `json:"consumer_statuses_changed"`
	SecretsCollected                    bool `json:"secrets_collected"`
	ComplianceClaimed                   bool `json:"compliance_claimed"`
	ProductionReadinessClaimed          bool `json:"production_readiness_claimed"`
	AgencyApprovalClaimed               bool `json:"agency_approval_claimed"`
	ConsumerAcceptanceClaimed           bool `json:"consumer_acceptance_claimed"`
	PublicLaunchClaimed                 bool `json:"public_launch_claimed"`
	HostedSaaSClaimed                   bool `json:"hosted_saas_claimed"`
	VendorCompatibilityClaimed          bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed        bool `json:"hardware_certification_claimed"`
	ProductionAVLReliabilityClaimed     bool `json:"production_avl_reliability_claimed"`
	ProductionGradeETAQualityClaimed    bool `json:"production_grade_eta_quality_claimed"`
	SLAClaimed                          bool `json:"sla_claimed"`
	UptimeGuaranteeClaimed              bool `json:"uptime_guarantee_claimed"`
	DynamicBackendPluginLoadingEnabled  bool `json:"dynamic_backend_plugin_loading_enabled"`
	ReleaseCandidateApprovalClaimed     bool `json:"release_candidate_approval_claimed"`
	ManagedSupportCommitmentClaimed     bool `json:"managed_support_commitment_claimed"`
	FinalDeploymentOwnershipClaimed     bool `json:"final_deployment_ownership_claimed"`
	ConsumerIngestionWorkflowCompleted  bool `json:"consumer_ingestion_workflow_completed"`
	ProductionMultiTenantHostingClaimed bool `json:"production_multi_tenant_hosting_claimed"`
}

func buildOperationsFirstRun(page operationsPage) operationsFirstRunView {
	feedURLs := firstRunFeedURLs(page)
	tasks := []operationsFirstRunTask{
		firstRunTask(
			1,
			"metadata",
			"Metadata",
			firstRunMetadataStatus(page),
			firstRunMetadataSignal(page),
			"Agency name, public root, open license, and technical contact are the first inputs consumers and operators need before feed publication can be reviewed.",
			"Enter or review publication metadata in the guided setup page; keep placeholders marked missing or review-needed.",
			"/admin/operations/setup",
			"docs/tutorials/agency-first-run.md",
			"Metadata in this console does not prove agency approval, final-root ownership, or public publication authorization.",
		),
		firstRunTask(
			2,
			"gtfs",
			"GTFS",
			gtfsLaunchpadStatus(page),
			gtfsLaunchpadSignal(page),
			"A published static schedule is the base for Vehicle Positions, Trip Updates, Alerts references, validation, and readiness checks.",
			"Use browser upload or safe URL import first when you have a GTFS ZIP; publish a GTFS Studio draft only when typed authoring is the right source.",
			"/admin/operations/gtfs-import",
			"docs/tutorials/real-agency-gtfs-onboarding.md",
			"A stored GTFS feed does not prove validator-clean status, agency approval, or Caltrans compliance.",
		),
		firstRunTask(
			3,
			"five_feed_urls",
			"Five feed URLs",
			fiveFeedsLaunchpadStatus(page),
			firstRunFiveFeedSignal(feedURLs),
			"The operator needs copyable locations for feeds.json, GTFS Schedule, Vehicle Positions, Trip Updates, and Alerts before checking publication posture.",
			"Copy or review each path below, then open the feed health command center for validation, freshness, and next actions.",
			"/admin/operations/feed-health",
			"docs/tutorials/operator-smoke-and-support-bundle.md",
			"URL listing does not prove public fetch success, stable final-root ownership, validator health, or consumer acceptance.",
		),
		firstRunTask(
			4,
			"validation_health",
			"Validation health",
			validatorLaunchpadStatus(page),
			validatorLaunchpadSignal(page),
			"Static GTFS and realtime validators are private quality gates and should be reviewed before stronger readiness language.",
			"Open validator health, install or configure pinned tooling if needed, then run only the allowlisted validator actions.",
			"/admin/operations/validation-health",
			"docs/tutorials/gtfs-validation-triage.md",
			"Validator output is supporting diagnostics only; it does not prove consumer ingestion, compliance, or feed correctness in all operating conditions.",
		),
		firstRunTask(
			5,
			"telemetry",
			"Telemetry",
			firstRunTelemetryStatus(page),
			firstRunTelemetrySignal(page),
			"Vehicle Positions depend on accepted telemetry and device bindings; stale or absent rows should remain visible instead of hidden.",
			"Bind a device token, send synthetic or deployment-owned telemetry through the authenticated ingest path, and review freshness.",
			"/admin/operations/devices",
			"docs/tutorials/telemetry-simulator-and-device-trial.md",
			"Telemetry shown here does not prove vendor compatibility, hardware certification, fleet reliability, or production AVL reliability.",
		),
		firstRunTask(
			6,
			"vp_tu_alerts",
			"VP/TU/Alerts",
			firstRunRealtimeStatus(page),
			firstRunRealtimeSignal(page),
			"Vehicle Positions are the first production-grade realtime output, while Trip Updates stay behind the prediction adapter and Alerts remain a separate lifecycle feed.",
			"Review feed health for Vehicle Positions, Trip Updates, and Alerts; fix missing URLs or validation gaps without coupling Trip Updates to ingest.",
			"/admin/operations/feed-health",
			"docs/requirements-calitp-compliance.md",
			"Realtime feed availability does not prove production-grade ETA quality, consumer display, or complete disruption handling.",
		),
		firstRunTask(
			7,
			"readiness",
			"Readiness",
			readinessLaunchpadStatus(page),
			readinessLaunchpadSignal(page),
			"Readiness rows combine publication, feed health, validation, telemetry, reliability, and consumer-workflow signals into a private review list.",
			"Open readiness and checklist views; leave missing source records missing until the underlying operator records exist.",
			"/admin/operations/readiness",
			"docs/tutorials/calitp-readiness-checklist.md",
			"Readiness review does not claim Caltrans compliance, public launch completion, production-readiness, or consumer acceptance.",
		),
		firstRunTask(
			8,
			"connectors",
			"Connectors",
			checklistStatusNeedsReview,
			"Connector Hub and connector test instructions are available as read-only local contract guidance.",
			"Optional external systems should stay behind sidecar, manifest, or adapter boundaries with synthetic conformance, redaction, and fail-closed checks first.",
			"Review Connector Hub, connector tests, and external-connection readiness before any deployment-owned external integration.",
			"/admin/operations/connectors",
			"docs/external-connection-readiness.md",
			"Connector guidance does not enable dynamic backend plugin loading, contact vendors, or prove vendor compatibility.",
		),
		firstRunTask(
			9,
			"support_rc_checks",
			"Support/RC checks",
			checklistStatusNeedsReview,
			"Operator smoke, deployment doctor, reliability, support bundle, and release-candidate checks are linked as private operator-run diagnostics.",
			"These checks help separate local demo behavior, target deployment review, and retained evidence intake before any outside claim is made.",
			"Run the documented checks only in the operator environment, review redaction, and keep RC decisions outside this GET-only page.",
			"/admin/operations/reliability",
			"docs/release-candidate-readiness.md",
			"Support and RC checks do not prove SLA coverage, managed support commitment, deployment ownership, public launch, or release approval.",
		),
	}
	return operationsFirstRunView{
		GeneratedAt:                         page.GeneratedAt,
		AgencyID:                            page.AgencyID,
		Boundary:                            "Start Here is a private authenticated first-run guide. It is GET-only, reads existing private records, runs no commands, creates no retained evidence, changes no consumer status, and makes no approval, compliance, public-launch, hosted-service, vendor, SLA, or ETA-quality claim.",
		LocalDemoDeploymentEvidenceBoundary: "Local demo checks show local wiring only. Deployment checks review a target environment only. Evidence requires a separate authorized intake with retention, redaction, and claim mapping; this page creates none of that.",
		Paths: []operationsFirstRunPath{
			firstRunPath(
				"no_developer",
				"No-developer path",
				"Use the private browser console; admin-only buttons remain on their existing pages and read-only users can still review status.",
				"An operator can complete the visible path with guided setup, browser GTFS import, device binding, validator health, feed health, readiness, and connector review pages.",
				"Open the setup wizard, then return here after each source record changes.",
				"/admin/operations/setup-wizard",
				"docs/tutorials/agency-first-run.md",
				"Browser guidance does not bypass role checks, CSRF checks, deployment setup, validation tooling, or evidence requirements.",
			),
			firstRunPath(
				"developer",
				"Developer path",
				"Use documented Make and Go commands from a terminal; the console only points to the checks and does not execute them.",
				"Developers can bootstrap local services, import GTFS, run validators, simulator checks, connector conformance, and release-candidate checks through existing CLI workflows.",
				"Follow the docs for local app startup and focused checks, then review the resulting private records here.",
				"/admin/operations/checklist",
				"docs/tutorials/self-hosted-operator-trial.md",
				"Copyable command guidance does not create retained evidence, contact external systems from this page, or prove deployment ownership.",
			),
		},
		Tasks:      tasks,
		FeedURLs:   feedURLs,
		Counts:     firstRunCounts(tasks, feedURLs),
		ClaimFlags: operationsFirstRunClaimFlags{},
	}
}

func firstRunTask(order int, id string, label string, status string, signal string, meaning string, nextAction string, uiLink string, docsLink string, doesNotProve string) operationsFirstRunTask {
	return operationsFirstRunTask{
		Order:         order,
		ID:            strings.TrimSpace(id),
		Label:         firstNonEmpty(label, id),
		Status:        normalizeChecklistStatus(status),
		CurrentSignal: firstNonEmpty(signal, "unknown"),
		Meaning:       firstNonEmpty(meaning, "This step summarizes an existing private operations signal."),
		NextAction:    firstNonEmpty(nextAction, "Review the linked private Operations Console section."),
		UILink:        firstRunAdminLink(uiLink),
		DocsLink:      firstRunDocLink(docsLink),
		DoesNotProve:  firstNonEmpty(doesNotProve, privateBoundary()),
	}
}

func firstRunPath(id string, label string, signal string, meaning string, firstAction string, uiLink string, docsLink string, doesNotProve string) operationsFirstRunPath {
	return operationsFirstRunPath{
		ID:            strings.TrimSpace(id),
		Label:         firstNonEmpty(label, id),
		CurrentSignal: firstNonEmpty(signal, "unknown"),
		Meaning:       firstNonEmpty(meaning, "This path summarizes a first-run operator workflow."),
		FirstAction:   firstNonEmpty(firstAction, "Open the linked private Operations Console page."),
		UILink:        firstRunAdminLink(uiLink),
		DocsLink:      firstRunDocLink(docsLink),
		DoesNotProve:  firstNonEmpty(doesNotProve, privateBoundary()),
	}
}

func firstRunFeedURLs(page operationsPage) []operationsFirstRunFeedURL {
	return []operationsFirstRunFeedURL{
		firstRunFeedURL(
			"feeds_json",
			"feeds.json",
			firstRunFeedsJSONURL(page),
			"publication metadata public base URL",
			"Discovery metadata lists feed URLs, contact, license, and update context for operators and downstream review.",
			"Set publication metadata before copying this outside the private console.",
			"docs/requirements-calitp-compliance.md",
			"A feeds.json URL does not prove source-of-truth website listing, final-root approval, or consumer ingestion.",
		),
		firstRunFeedURL(
			"schedule",
			"GTFS Schedule",
			firstRunFeedURLByType(page, "schedule"),
			"published_feed schedule metadata",
			"The schedule ZIP is the static GTFS base used by realtime references and validation.",
			"Import or publish GTFS, then validate the schedule before external sharing.",
			"docs/tutorials/real-agency-gtfs-onboarding.md",
			"A schedule URL does not prove validator-clean status, agency approval, or open-license publication.",
		),
		firstRunFeedURL(
			"vehicle_positions",
			"Vehicle Positions",
			firstRunFeedURLByType(page, "vehicle_positions"),
			"published_feed Vehicle Positions metadata",
			"Vehicle Positions are the first high-quality realtime feed target for Open Transit RT.",
			"Confirm fresh telemetry, public protobuf output, feed health, and realtime validation.",
			"docs/requirements-calitp-compliance.md",
			"A Vehicle Positions URL does not prove real-device reliability or consumer acceptance.",
		),
		firstRunFeedURL(
			"trip_updates",
			"Trip Updates",
			firstRunFeedURLByType(page, "trip_updates"),
			"published_feed Trip Updates metadata",
			"Trip Updates should remain replaceable behind the prediction adapter boundary.",
			"Review prediction diagnostics, withheld cases, and realtime validation before relying on Trip Updates.",
			"docs/requirements-trip-updates.md",
			"A Trip Updates URL does not prove production-grade ETA quality or stop-level prediction accuracy.",
		),
		firstRunFeedURL(
			"alerts",
			"Alerts",
			firstRunFeedURLByType(page, "alerts"),
			"published_feed Alerts metadata",
			"Alerts complete the expected realtime feed set and stay separate from telemetry ingest.",
			"Use the Alerts Console for lifecycle checks, then validate the Alerts feed.",
			"docs/requirements-calitp-compliance.md",
			"An Alerts URL does not prove agency approval, consumer display, or disruption-workflow completeness.",
		),
	}
}

func firstRunFeedURL(id string, label string, url string, source string, meaning string, nextAction string, docsLink string, doesNotProve string) operationsFirstRunFeedURL {
	status := checklistStatusMissing
	copyValue := "missing"
	trimmed := strings.TrimSpace(url)
	if trimmed != "" {
		status = checklistStatusNeedsReview
		copyValue = trimmed
	}
	return operationsFirstRunFeedURL{
		ID:           strings.TrimSpace(id),
		Label:        firstNonEmpty(label, id),
		Status:       status,
		URL:          trimmed,
		CopyValue:    copyValue,
		Source:       firstNonEmpty(source, "feed discovery"),
		Meaning:      firstNonEmpty(meaning, "This URL is copied from existing private publication metadata."),
		NextAction:   firstNonEmpty(nextAction, "Review the feed URL before using it outside this console."),
		DocsLink:     firstRunDocLink(docsLink),
		DoesNotProve: firstNonEmpty(doesNotProve, privateBoundary()),
	}
}

func firstRunCounts(tasks []operationsFirstRunTask, feedURLs []operationsFirstRunFeedURL) operationsFirstRunCounts {
	counts := operationsFirstRunCounts{Tasks: len(tasks), FeedURLs: len(feedURLs), Statuses: map[string]int{
		checklistStatusOK:          0,
		checklistStatusNeedsReview: 0,
		checklistStatusMissing:     0,
		checklistStatusBlocked:     0,
		checklistStatusUnknown:     0,
	}}
	for _, task := range tasks {
		counts.Statuses[normalizeChecklistStatus(task.Status)]++
	}
	return counts
}

func firstRunMetadataStatus(page operationsPage) string {
	if page.DiscoveryError != "" || page.PublicationError != "" {
		return checklistStatusMissing
	}
	return metadataStatus(page.Discovery.AgencyName, page.Discovery.License.Name, page.Discovery.License.URL, page.Discovery.TechnicalContactEmail)
}

func firstRunMetadataSignal(page operationsPage) string {
	if page.DiscoveryError != "" || page.PublicationError != "" {
		return firstNonEmpty(page.DiscoveryError, page.PublicationError)
	}
	return fmt.Sprintf("agency=%q; license=%q; contact=%q; public root=%q",
		firstNonEmpty(page.Discovery.AgencyName, "missing"),
		firstNonEmpty(page.Discovery.License.Name, "missing"),
		firstNonEmpty(page.Discovery.TechnicalContactEmail, "missing"),
		firstNonEmpty(page.Discovery.PublicBaseURL, "missing"),
	)
}

func firstRunTelemetryStatus(page operationsPage) string {
	return rollupStatuses([]string{deviceChecklistStatus(page), telemetryLaunchpadStatus(page)})
}

func firstRunTelemetrySignal(page operationsPage) string {
	return fmt.Sprintf("devices: %s; telemetry: %s", deviceEvidence(page), telemetryLaunchpadSignal(page))
}

func firstRunRealtimeStatus(page operationsPage) string {
	return rollupStatuses([]string{
		feedChecklistStatus(page, "vehicle_positions"),
		feedChecklistStatus(page, "trip_updates"),
		feedChecklistStatus(page, "alerts"),
	})
}

func firstRunRealtimeSignal(page operationsPage) string {
	return fmt.Sprintf("Vehicle Positions: %s; Trip Updates: %s; Alerts: %s",
		feedStatus(page, "vehicle_positions"),
		feedStatus(page, "trip_updates"),
		feedStatus(page, "alerts"),
	)
}

func firstRunFiveFeedSignal(feedURLs []operationsFirstRunFeedURL) string {
	present := 0
	for _, row := range feedURLs {
		if row.URL != "" {
			present++
		}
	}
	return fmt.Sprintf("%d of %d feed URL copy values are present", present, len(feedURLs))
}

func firstRunFeedsJSONURL(page operationsPage) string {
	base := strings.TrimRight(strings.TrimSpace(page.Discovery.PublicBaseURL), "/")
	if base == "" {
		return ""
	}
	return base + "/public/feeds.json"
}

func firstRunFeedURLByType(page operationsPage, feedType string) string {
	for _, feed := range page.Discovery.Feeds {
		if feed.FeedType == feedType {
			return strings.TrimSpace(feed.CanonicalPublicURL)
		}
	}
	return ""
}

func firstRunAdminLink(value string) string {
	links := safeAdminLinks([]string{value})
	if len(links) == 0 {
		return "/admin/operations"
	}
	return links[0]
}

func firstRunDocLink(value string) string {
	links := safeDocsLinks([]string{value})
	if len(links) == 0 {
		return "docs/tutorials/agency-first-run.md"
	}
	return links[0]
}
