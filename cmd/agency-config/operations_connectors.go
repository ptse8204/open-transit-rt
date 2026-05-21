package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
	connectorpkg "open-transit-rt/internal/connectors"
)

const safePluginDefinition = "In Open Transit RT, a plugin is an optional sidecar, command adapter, manifest, or connector process. It is not arbitrary dynamic code loaded into the backend."

type connectorHubView struct {
	GeneratedAt      time.Time                `json:"generated_at"`
	AgencyID         string                   `json:"agency_id"`
	Boundary         string                   `json:"boundary"`
	PluginDefinition string                   `json:"plugin_definition"`
	InstanceSummary  connectorInstanceSummary `json:"instance_summary"`
	Instances        []connectorInstanceRow   `json:"instances"`
	Health           []connectorHealthRow     `json:"health"`
	Catalog          []connectorCatalogRow    `json:"catalog"`
	Categories       []connectorCategory      `json:"categories"`
	Registry         connectorpkg.Registry    `json:"registry"`
	ClaimFlags       connectorHubClaimFlags   `json:"claim_flags"`
}

type connectorInstanceStore interface {
	ListInstances(ctx context.Context, agencyID string) ([]connectorpkg.Instance, error)
}

type connectorInstanceSummary struct {
	ConnectorTypes      int    `json:"connector_types"`
	ConfiguredInstances int    `json:"configured_instances"`
	ExampleManifests    int    `json:"example_manifests"`
	OverallState        string `json:"overall_state"`
	Boundary            string `json:"boundary"`
}

type connectorInstanceRow struct {
	ID                     string   `json:"id"`
	ConnectorType          string   `json:"connector_type"`
	ConnectorKind          string   `json:"connector_kind"`
	DisplayName            string   `json:"display_name"`
	State                  string   `json:"state"`
	Owner                  string   `json:"owner"`
	ExamplesAvailable      int      `json:"examples_available"`
	DeploymentConfigExists string   `json:"deployment_config_exists"`
	ConfigMetadata         string   `json:"config_metadata"`
	SecretRefs             []string `json:"secret_refs"`
	DryRunStatus           string   `json:"dry_run_status"`
	ActivationReadiness    string   `json:"activation_readiness"`
	LastSignal             string   `json:"last_signal"`
	NextAction             string   `json:"next_action"`
	SafeLinks              []string `json:"safe_links"`
	Limits                 string   `json:"limits"`
}

type connectorCatalogRow struct {
	ID             string   `json:"id"`
	Group          string   `json:"group"`
	Label          string   `json:"label"`
	Status         string   `json:"status"`
	StartWith      string   `json:"start_with"`
	BrowserReview  string   `json:"browser_review"`
	FirstSafeCheck string   `json:"first_safe_check"`
	DoesNotProve   string   `json:"does_not_prove"`
	DocsLinks      []string `json:"docs_links"`
}

type connectorCategory struct {
	ID                 string   `json:"id"`
	Label              string   `json:"label"`
	Status             string   `json:"status"`
	Summary            string   `json:"summary"`
	ConnectorShape     string   `json:"connector_shape"`
	Inputs             []string `json:"inputs"`
	Outputs            []string `json:"outputs"`
	FailureBehavior    string   `json:"failure_behavior"`
	CommandSuggestions []string `json:"command_suggestions"`
	AdminLinks         []string `json:"admin_links"`
	DocsLinks          []string `json:"docs_links"`
	ClaimBoundary      string   `json:"claim_boundary"`
}

type connectorHealthRow struct {
	ID                 string   `json:"id"`
	Label              string   `json:"label"`
	Owner              string   `json:"owner"`
	Status             string   `json:"status"`
	Configured         string   `json:"configured"`
	DryRunReady        string   `json:"dry_run_ready"`
	SendState          string   `json:"send_state"`
	LastSyntheticCheck string   `json:"last_synthetic_check"`
	RedactionStatus    string   `json:"redaction_status"`
	KnownBlockers      []string `json:"known_blockers"`
	IssueCategory      string   `json:"issue_category"`
	IssueLinks         []string `json:"issue_links"`
	SetupChecklist     []string `json:"setup_checklist"`
	ChecklistCopy      string   `json:"checklist_copy"`
	DoesNotProve       string   `json:"does_not_prove"`
}

type connectorHubClaimFlags struct {
	DynamicBackendPluginLoadingEnabled bool `json:"dynamic_backend_plugin_loading_enabled"`
	VendorCompatibilityClaimed         bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed       bool `json:"hardware_certification_claimed"`
	ConsumerStatusesChanged            bool `json:"consumer_statuses_changed"`
	ExternalEvidenceCreated            bool `json:"external_evidence_created"`
	ComplianceClaimed                  bool `json:"compliance_claimed"`
	ProductionReadinessClaimed         bool `json:"production_readiness_claimed"`
	HostedSaaSClaimed                  bool `json:"hosted_saas_claimed"`
	ProductionGradeETAClaimed          bool `json:"production_grade_eta_claimed"`
}

func (h *handler) renderConnectorHub(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "connectors")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "connectors", page)
}

func (h *handler) renderConnectorHubJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "connectors")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.ConnectorHub)
}

func buildConnectorHub(page operationsPage) connectorHubView {
	registry := connectorRegistryForSection(page.Section)
	instanceRows := connectorInstanceRows(page, registry)
	categories := []connectorCategory{
		connectorCategoryView(
			"telemetry_source",
			"Vehicle / GPS / AVL connectors",
			"available",
			"Use CSV replay, HTTP polling, webhook sidecars, generic JSON transforms, vendor-shaped synthetic examples, and authenticated telemetry POST paths.",
			"sidecar or command adapter that normalizes observations before authenticated POST /v1/telemetry",
			[]string{"external observations", "agency-owned mapping", "device credential reference"},
			[]string{"Open Transit RT telemetry.Event records", "private redacted adapter diagnostics"},
			"Fail closed before ingest when mapping, credentials, timestamps, or payload quality are unsafe.",
			[]string{"make telemetry-simulator", "make adapter-conformance", "make test-connector-examples"},
			[]string{"/admin/operations/connectors/workbench", "/admin/operations/devices", "/admin/operations/telemetry"},
			[]string{"docs/connectors/catalog.md", "docs/connectors/vehicle-avl-starter-kits.md", "docs/integration-adapter-kit.md", "docs/tutorials/device-avl-integration.md"},
			"Telemetry connectors do not prove real vendor compatibility, hardware certification, production AVL reliability, or agency adoption.",
		),
		connectorCategoryView(
			"prediction",
			"Prediction connectors",
			"available",
			"Review the deterministic built-in predictor, external HTTP predictor adapter, shadow mode, fail-closed behavior, and TheTransitClock candidate boundary.",
			"optional internal adapter or external HTTP sidecar behind internal/prediction.Adapter",
			[]string{"active GTFS feed version", "latest telemetry", "current assignments", "Vehicle Positions feed data or URL"},
			[]string{"Trip Updates feed", "bounded diagnostics"},
			"Return deterministic output or valid empty Trip Updates with adapter diagnostics when an optional sidecar fails.",
			[]string{"make adapter-conformance", "go test ./internal/prediction"},
			[]string{"/admin/operations/feeds", "/admin/operations/reliability"},
			[]string{"docs/connectors/catalog.md", "docs/requirements-trip-updates.md", "docs/integration-adapter-kit.md", "docs/tutorials/external-adapter-conformance.md"},
			"Prediction connectors do not prove production-grade ETA quality, consumer acceptance, or named predictor compatibility.",
		),
		connectorCategoryView(
			"validator",
			"Validator connectors",
			"available",
			"Review MobilityData static GTFS validation, MobilityData GTFS Realtime validation, allowlisted validator IDs, and private validation health.",
			"allowlisted command adapter configured by deployment-owned environment",
			[]string{"server-derived schedule ZIP", "server-derived realtime protobuf", "validator ID"},
			[]string{"normalized validation_report rows", "private validator health diagnostics"},
			"Missing or failing tooling records not-run, missing, warning, failed, or blocked states instead of pretending success.",
			[]string{"make validators-check", "make validator-health"},
			[]string{"/admin/operations/validation-health", "/admin/operations/gtfs-quality"},
			[]string{"docs/connectors/catalog.md", "docs/dependencies.md", "docs/tutorials/gtfs-validation-triage.md", "docs/release-candidate-readiness.md"},
			"Validator records are supporting signals only; they are not CAL-ITP/Caltrans compliance or consumer acceptance.",
		),
		connectorCategoryView(
			"monitoring_export",
			"Monitoring / export connectors",
			"available",
			"Review local health summaries, operations notify drafts, monitoring/export helpers, and the deployment-owned monitoring boundary.",
			"deployment-owned no-send exporter or redacted summary writer",
			[]string{"validator health summary", "deployment doctor summary", "reliability summary"},
			[]string{"private redacted monitoring/export summaries"},
			"Default examples write local private summaries only and require deployment-owned configuration before any external send.",
			[]string{"make operations-reliability", "make operations-notify"},
			[]string{"/admin/operations/reliability", "/admin/operations/validation-health"},
			[]string{"docs/connectors/catalog.md", "docs/tutorials/self-hosted-operations-notifications.md", "docs/deployment/reference-deployment-doctor.md"},
			"Monitoring connectors do not prove SLA coverage, uptime guarantees, hosted service availability, or production readiness.",
		),
		connectorCategoryView(
			"consumer_discovery",
			"Consumer / discovery connectors",
			"available",
			"Review /public/feeds.json, static GTFS, Vehicle Positions, Trip Updates, Alerts, and consumer packet preparedness without submission or acceptance claims.",
			"docs workflow and packet generator, not portal automation",
			[]string{"stable feed URL metadata", "license/contact metadata", "validation status"},
			[]string{"prepared consumer packet records", "private readiness summaries"},
			"Do not contact targets or change statuses without retained target-originated or operator-retained authorization artifacts.",
			[]string{"make external-connection-check", "make audit-final-claim-review"},
			[]string{"/admin/operations/consumers", "/admin/operations/readiness"},
			[]string{"docs/connectors/catalog.md", "docs/consumer-submission-evidence.md", "docs/evidence/consumer-submissions/README.md", "docs/compliance-evidence-checklist.md"},
			"Prepared packets are not submission, review, acceptance, ingestion, listing, display, compliance, or public launch proof.",
		),
		connectorCategoryView(
			"future_extension_model",
			"Future connector extension model",
			"policy",
			"Use manifest-based sidecars, keep arbitrary dynamic backend plugin loading disabled, and require conformance tests for new connector shapes.",
			"manifest-based sidecar or command adapter with explicit contracts",
			[]string{"connector manifest", "synthetic fixture", "redaction plan", "conformance cases"},
			[]string{"reviewed connector example", "private diagnostics", "bounded docs"},
			"Reject secrets, raw commands, status mutation, external submission automation, unreviewed runtime loading, and unsupported claims.",
			[]string{"make external-connection-check", "make adapter-conformance"},
			[]string{"/admin/operations/connectors/workbench", "/admin/operations/connectors/tests"},
			[]string{"docs/connectors/catalog.md", "docs/connectors/plugin-contract.md", "docs/connectors/contributing-connectors.md", "docs/extension-governance.md"},
			"Extension policy does not show runtime distribution support, third-party compatibility, compliance, consumer acceptance, or production readiness.",
		),
	}
	return connectorHubView{
		GeneratedAt:      page.GeneratedAt,
		AgencyID:         page.AgencyID,
		Boundary:         "Private authenticated Connectors only; viewing it creates no evidence, contacts no external party, changes no consumer status, and records no approval, compatibility, compliance, hosted-service, SLA, production-readiness, or ETA-quality outcome. Treat it as the starting point for manifest, redaction, fail-closed, and adapter-conformance review.",
		PluginDefinition: safePluginDefinition,
		InstanceSummary:  connectorInstanceSummaryForRows(instanceRows, registry),
		Instances:        instanceRows,
		Health:           connectorHealthRows(page, registry),
		Catalog:          connectorCatalogRows(),
		Categories:       categories,
		Registry:         registry,
		ClaimFlags:       connectorHubClaimFlags{},
	}
}

func (h *handler) connectorInstancesForPage(r *http.Request, agencyID string) ([]connectorpkg.Instance, string) {
	if h == nil || h.connectorInstances == nil {
		return nil, ""
	}
	instances, err := h.connectorInstances.ListInstances(r.Context(), agencyID)
	if err != nil {
		return nil, "connector instance records are not available"
	}
	return instances, ""
}

func connectorInstanceRows(page operationsPage, registry connectorpkg.Registry) []connectorInstanceRow {
	var rows []connectorInstanceRow
	for _, connectorType := range connectorpkg.SupportedTypes {
		instances := connectorInstancesByType(page.ConnectorInstances, connectorType)
		exampleCount := len(connectorRegistryEntries(registry, connectorType))
		if len(instances) == 0 {
			state := string(connectorpkg.StateNotConfigured)
			lastSignal := "No committed example manifest loaded for this connector type."
			if exampleCount > 0 {
				state = string(connectorpkg.StateExampleAvailable)
				lastSignal = strconv.Itoa(exampleCount) + " committed example manifest(s) available; examples are not configured connectors."
			}
			rows = append(rows, connectorInstanceRow{
				ID:                     connectorType + ":example-library",
				ConnectorType:          connectorType,
				ConnectorKind:          "example library",
				DisplayName:            connectorTypeLabel(connectorType) + " examples",
				State:                  state,
				Owner:                  "deployment owner",
				ExamplesAvailable:      exampleCount,
				DeploymentConfigExists: "no",
				ConfigMetadata:         "not configured",
				SecretRefs:             nil,
				DryRunStatus:           "not_run",
				ActivationReadiness:    "not ready; create a per-agency connector instance first",
				LastSignal:             lastSignal,
				NextAction:             "Create or review a per-agency connector instance before recording dry-run or activation status.",
				SafeLinks:              connectorSafeLinks(connectorType),
				Limits:                 "Example manifests are not deployment-owned configuration and do not prove compatibility, reliability, compliance, acceptance, or readiness.",
			})
			continue
		}
		for _, instance := range instances {
			rows = append(rows, connectorInstanceRowFromInstance(instance, exampleCount))
		}
	}
	return rows
}

func connectorInstanceSummaryForRows(rows []connectorInstanceRow, registry connectorpkg.Registry) connectorInstanceSummary {
	configured := 0
	state := string(connectorpkg.StateNotConfigured)
	for _, row := range rows {
		if row.State != string(connectorpkg.StateExampleAvailable) && row.State != string(connectorpkg.StateNotConfigured) {
			configured++
			state = row.State
		}
		if row.State == string(connectorpkg.StateActive) {
			state = row.State
		}
		if row.State == string(connectorpkg.StateBlocked) && state != string(connectorpkg.StateActive) {
			state = row.State
		}
	}
	if configured == 0 && connectorRegistryEntryCount(registry) > 0 {
		state = string(connectorpkg.StateExampleAvailable)
	}
	return connectorInstanceSummary{
		ConnectorTypes:      len(connectorpkg.SupportedTypes),
		ConfiguredInstances: configured,
		ExampleManifests:    connectorRegistryEntryCount(registry),
		OverallState:        state,
		Boundary:            "Configured connector instances are deployment-owned agency records. Committed example manifests are separated and never count as configured, dry-run-passed, ready, or active connectors.",
	}
}

func connectorInstanceRowFromInstance(instance connectorpkg.Instance, exampleCount int) connectorInstanceRow {
	return connectorInstanceRow{
		ID:                     strconv.FormatInt(instance.ID, 10),
		ConnectorType:          instance.ConnectorType,
		ConnectorKind:          firstNonEmpty(instance.ConnectorKind, "unspecified"),
		DisplayName:            firstNonEmpty(instance.DisplayName, instance.ConnectorKind, "Connector instance"),
		State:                  string(instance.State),
		Owner:                  firstNonEmpty(instance.Owner, "unassigned"),
		ExamplesAvailable:      exampleCount,
		DeploymentConfigExists: yesNo(instance.DeploymentConfigExists()),
		ConfigMetadata:         connectorConfigMetadataSummary(instance),
		SecretRefs:             instance.SecretRefs,
		DryRunStatus:           firstNonEmpty(instance.DryRunStatus, "not_run"),
		ActivationReadiness:    connectorActivationReadiness(instance),
		LastSignal:             connectorInstanceLastSignal(instance),
		NextAction:             connectorInstanceNextAction(instance),
		SafeLinks:              connectorSafeLinks(instance.ConnectorType),
		Limits:                 "Instance state is an internal configuration signal only. It does not prove vendor compatibility, external delivery, compliance, consumer acceptance, SLA coverage, AVL reliability, or ETA quality.",
	}
}

func connectorConfigMetadataSummary(instance connectorpkg.Instance) string {
	keys := instance.ConfigKeys()
	if len(keys) == 0 && len(instance.SecretRefs) == 0 {
		return "no non-secret config metadata or secret refs recorded"
	}
	var parts []string
	if len(keys) > 0 {
		parts = append(parts, "metadata keys: "+strings.Join(keys, ", "))
	}
	if len(instance.SecretRefs) > 0 {
		parts = append(parts, "secret refs: "+strings.Join(instance.SecretRefs, ", "))
	}
	return strings.Join(parts, "; ")
}

func connectorActivationReadiness(instance connectorpkg.Instance) string {
	switch instance.State {
	case connectorpkg.StateActive:
		return "active by deployment-owned signal; browser did not start a process"
	case connectorpkg.StateReadyForActivation:
		return "ready for deployment-owned activation"
	case connectorpkg.StateDryRunPassed:
		return "dry-run passed; activation still requires device/token and deployment checks"
	case connectorpkg.StateConfiguredNotTested:
		return "configured but dry-run is required before activation"
	case connectorpkg.StateBlocked:
		return "blocked; fix owner, mapping, dry-run, redaction, or deployment references"
	default:
		return "not ready"
	}
}

func connectorInstanceLastSignal(instance connectorpkg.Instance) string {
	if instance.LastCheckedAt != nil {
		return "last checked " + instance.LastCheckedAt.UTC().Format(time.RFC3339)
	}
	if instance.ActivatedAt != nil {
		return "activated " + instance.ActivatedAt.UTC().Format(time.RFC3339)
	}
	return "no health check or accepted/rejected telemetry signal recorded"
}

func connectorInstanceNextAction(instance connectorpkg.Instance) string {
	switch instance.State {
	case connectorpkg.StateActive:
		return "Monitor health checks and accepted/rejected telemetry without exposing secrets or raw payloads."
	case connectorpkg.StateReadyForActivation:
		return "Have a deployment owner activate the connector outside the browser and then record the signal."
	case connectorpkg.StateDryRunPassed:
		return "Review activation readiness, device bindings, token refs, stale rules, and redaction before activation."
	case connectorpkg.StateConfiguredNotTested:
		return "Run or record a server-owned dry-run before activation can be reviewed."
	case connectorpkg.StateBlocked:
		return "Resolve the blocker and rerun the dry-run or health check."
	default:
		return "Complete connector metadata before dry-run or activation review."
	}
}

func connectorSafeLinks(connectorType string) []string {
	switch connectorType {
	case connectorpkg.TypeTelemetrySource:
		return []string{"/admin/operations/connectors/workbench", "/admin/operations/devices", "/admin/operations/telemetry"}
	case connectorpkg.TypePrediction:
		return []string{"/admin/operations/prediction-lab", "/admin/operations/realtime", "/admin/operations/connectors/tests"}
	case connectorpkg.TypeValidator:
		return []string{"/admin/operations/validation-health", "/admin/operations/validation-center", "/admin/operations/connectors/tests"}
	case connectorpkg.TypeMonitoringExport:
		return []string{"/admin/operations/maintenance", "/admin/operations/reliability", "/admin/operations/connectors/tests"}
	case connectorpkg.TypeConsumerDiscovery:
		return []string{"/admin/operations/feeds", "/admin/operations/consumers", "/admin/operations/readiness"}
	default:
		return []string{"/admin/operations/connectors/tests"}
	}
}

func connectorTypeLabel(connectorType string) string {
	switch connectorType {
	case connectorpkg.TypeTelemetrySource:
		return "Vehicle / GPS / AVL"
	case connectorpkg.TypePrediction:
		return "Prediction"
	case connectorpkg.TypeValidator:
		return "Validator"
	case connectorpkg.TypeMonitoringExport:
		return "Monitoring / export"
	case connectorpkg.TypeConsumerDiscovery:
		return "Discovery"
	default:
		return firstNonEmpty(connectorType, "Connector")
	}
}

func connectorRegistryEntryCount(registry connectorpkg.Registry) int {
	return len(registry.Entries)
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func connectorHealthRows(page operationsPage, registry connectorpkg.Registry) []connectorHealthRow {
	return []connectorHealthRow{
		connectorHealthRowView(
			"telemetry_source",
			"Vehicle data setup",
			"administrator",
			connectorTypeStatus(registry, page.ConnectorInstances, connectorpkg.TypeTelemetrySource, connectorTelemetryBlockers(page, registry)),
			connectorConfiguredSignal(registry, page.ConnectorInstances, connectorpkg.TypeTelemetrySource, "telemetry source"),
			connectorDryRunSignal(registry, page.ConnectorInstances, connectorpkg.TypeTelemetrySource, "make test-connector-examples"),
			"Example sends are disabled; authenticated /v1/telemetry ingest needs a deployment-owned device token outside this page.",
			"Browser does not record runs. First check: make test-connector-examples.",
			connectorRedactionSignal(registry, connectorpkg.TypeTelemetrySource),
			connectorTelemetryBlockers(page, registry),
			"telemetry",
			[]string{"/admin/operations/devices", "/admin/operations/telemetry", "/admin/operations/realtime", "/admin/operations/connectors/tests"},
			[]string{
				"choose_manifest=example.telemetry-csv-replay|example.telemetry-http-poller|example.telemetry-webhook-sidecar|example.generic-json-transform",
				"keep_send_enabled=false",
				"keep_network_send=false",
				"map_agency_device_vehicle_timestamp_lat_lon_quality",
				"run_make_test_connector_examples_before_real_source",
			},
			"Real device proof, vendor compatibility, hardware certification, production AVL reliability, compliance, consumer acceptance, or production readiness.",
		),
		connectorHealthRowView(
			"prediction",
			"Prediction setup",
			"developer/integrator",
			connectorTypeStatus(registry, page.ConnectorInstances, connectorpkg.TypePrediction, connectorPredictionBlockers(page, registry)),
			connectorConfiguredSignal(registry, page.ConnectorInstances, connectorpkg.TypePrediction, "prediction"),
			connectorDryRunSignal(registry, page.ConnectorInstances, connectorpkg.TypePrediction, "go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance"),
			"Optional sidecars are disabled by default; Vehicle Positions publishing stays independent.",
			"Browser does not record runs. First check: adapter-conformance prediction.",
			connectorRedactionSignal(registry, connectorpkg.TypePrediction),
			connectorPredictionBlockers(page, registry),
			"prediction",
			[]string{"/admin/operations/prediction-lab", "/admin/operations/realtime", "/admin/operations/feed-health", "/admin/operations/connectors/tests"},
			[]string{
				"choose_manifest=example.predictor-sidecar-stub",
				"keep_public_mutation=false",
				"require_vehicle_positions_reference",
				"withhold_on_timeout_malformed_stale_wrong_agency_low_confidence",
				"run_adapter_conformance_prediction",
			},
			"Production-grade ETA quality, real-world accuracy, named predictor support, consumer acceptance, or production readiness.",
		),
		connectorHealthRowView(
			"validator",
			"Validator setup",
			"administrator",
			connectorTypeStatus(registry, page.ConnectorInstances, connectorpkg.TypeValidator, connectorValidatorBlockers(page, registry)),
			connectorConfiguredSignal(registry, page.ConnectorInstances, connectorpkg.TypeValidator, "validator"),
			connectorDryRunSignal(registry, page.ConnectorInstances, connectorpkg.TypeValidator, "go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance"),
			"Validator IDs are allowlisted; browser pages never accept raw validator commands.",
			"Browser does not record runs. First check: adapter-conformance validator.",
			connectorRedactionSignal(registry, connectorpkg.TypeValidator),
			connectorValidatorBlockers(page, registry),
			"validation",
			[]string{"/admin/operations/validation-health", "/admin/operations/validation-center", "/admin/operations/gtfs-quality", "/admin/operations/connectors/tests"},
			[]string{
				"choose_manifest=example.validator-allowlist",
				"keep_validator_id_allowlisted",
				"keep_artifact_paths_private_and_redacted",
				"record_not_run_or_blocked_instead_of_success",
				"run_adapter_conformance_validator",
			},
			"Validator-clean proof, CAL-ITP/Caltrans compliance, consumer acceptance, public launch, or production readiness.",
		),
		connectorHealthRowView(
			"monitoring_export",
			"Monitoring export setup",
			"deployment owner",
			connectorTypeStatus(registry, page.ConnectorInstances, connectorpkg.TypeMonitoringExport, connectorMonitoringBlockers(page, registry)),
			connectorConfiguredSignal(registry, page.ConnectorInstances, connectorpkg.TypeMonitoringExport, "monitoring export"),
			connectorDryRunSignal(registry, page.ConnectorInstances, connectorpkg.TypeMonitoringExport, "go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance"),
			"Example notification and export sends are disabled; any destination is deployment-owned outside this page.",
			"Browser does not record runs. First check: adapter-conformance monitoring.",
			connectorRedactionSignal(registry, connectorpkg.TypeMonitoringExport),
			connectorMonitoringBlockers(page, registry),
			"monitoring",
			[]string{"/admin/operations/maintenance", "/admin/operations/reliability", "/admin/operations/validation-health", "/admin/operations/connectors/tests"},
			[]string{
				"choose_manifest=example.monitoring-export",
				"keep_send_by_default=false",
				"redact_destination_and_contact_fields",
				"review_health_digest_before_delivery",
				"run_adapter_conformance_monitoring",
			},
			"SLA coverage, uptime guarantee, hosted service availability, notification delivery, compliance, retained evidence, or production readiness.",
		),
		connectorHealthRowView(
			"consumer_discovery",
			"Feed discovery setup",
			"deployment owner",
			connectorTypeStatus(registry, page.ConnectorInstances, connectorpkg.TypeConsumerDiscovery, connectorConsumerDiscoveryBlockers(page, registry)),
			connectorConfiguredSignal(registry, page.ConnectorInstances, connectorpkg.TypeConsumerDiscovery, "consumer discovery"),
			connectorDryRunSignal(registry, page.ConnectorInstances, connectorpkg.TypeConsumerDiscovery, "go run ./cmd/adapter-conformance consumer_discovery --suite testdata/adapter-conformance"),
			"Submission automation and consumer status mutation stay disabled.",
			"Browser does not record runs. First check: adapter-conformance consumer_discovery.",
			connectorRedactionSignal(registry, connectorpkg.TypeConsumerDiscovery),
			connectorConsumerDiscoveryBlockers(page, registry),
			"consumer discovery",
			[]string{"/admin/operations/feeds", "/admin/operations/consumers", "/admin/operations/readiness", "/admin/operations/connectors/tests"},
			[]string{
				"choose_manifest=example.consumer-discovery-metadata",
				"keep_submit_enabled=false",
				"keep_status_mutation=false",
				"review_public_base_url_license_contact_feed_urls",
				"run_adapter_conformance_consumer_discovery",
			},
			"Consumer submission, review, acceptance, ingestion, listing, display, compliance, public launch, or production readiness.",
		),
		connectorHealthRowView(
			"future_extension_model",
			"Future extension setup",
			"developer/integrator",
			connectorFutureExtensionStatus(registry),
			"Manifest-based examples are the only reviewed extension shape.",
			"Ready when make external-connection-check and adapter-conformance pass.",
			"Arbitrary backend plugin loading and browser command execution are disabled.",
			"Browser does not record runs. First check: make external-connection-check.",
			connectorRegistryRedactionStatus(registry),
			connectorRegistryBlockers(registry),
			"extension governance",
			[]string{"/admin/operations/connectors/workbench", "/admin/operations/connectors/tests"},
			[]string{
				"keep_plugin_model=sidecar_or_manifest",
				"reject_dynamic_backend_loading",
				"reject_manifest_commands_and_private_endpoints",
				"add_synthetic_fixture_and_conformance_case",
				"run_external_connection_check",
			},
			"Runtime plugin distribution, unreviewed extension execution, vendor compatibility, compliance, or production readiness.",
		),
	}
}

func connectorHealthRowView(id string, label string, owner string, status string, configured string, dryRunReady string, sendState string, lastSyntheticCheck string, redactionStatus string, blockers []string, issueCategory string, issueLinks []string, checklist []string, doesNotProve string) connectorHealthRow {
	cleanChecklist := cleanLaunchpadList(checklist)
	return connectorHealthRow{
		ID:                 id,
		Label:              firstNonEmpty(label, "Connector setup"),
		Owner:              firstNonEmpty(owner, "administrator"),
		Status:             firstNonEmpty(status, checklistStatusNeedsReview),
		Configured:         firstNonEmpty(configured, "Review connector manifest configuration."),
		DryRunReady:        firstNonEmpty(dryRunReady, "Run fixed local checks before integration work."),
		SendState:          firstNonEmpty(sendState, "Sends are disabled by default."),
		LastSyntheticCheck: firstNonEmpty(lastSyntheticCheck, "Browser does not record synthetic check runs."),
		RedactionStatus:    firstNonEmpty(redactionStatus, "Review redaction policy before use."),
		KnownBlockers:      cleanConnectorOptionalList(blockers),
		IssueCategory:      firstNonEmpty(issueCategory, "connectors"),
		IssueLinks:         safeAdminLinks(issueLinks),
		SetupChecklist:     cleanChecklist,
		ChecklistCopy:      strings.Join(cleanChecklist, "\n"),
		DoesNotProve:       firstNonEmpty(doesNotProve, "Compatibility, compliance, production readiness, or consumer acceptance."),
	}
}

func cleanConnectorOptionalList(values []string) []string {
	var out []string
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" || unsafePrivateString(trimmed) {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}

func connectorTypeStatus(registry connectorpkg.Registry, instances []connectorpkg.Instance, connectorType string, blockers []string) string {
	if len(registry.Diagnostics) > 0 {
		return checklistStatusBlocked
	}
	if len(connectorRegistryEntries(registry, connectorType)) == 0 {
		return checklistStatusBlocked
	}
	typed := connectorInstancesByType(instances, connectorType)
	if len(typed) == 0 {
		return string(connectorpkg.StateExampleAvailable)
	}
	state := strongestConnectorState(typed)
	if state == connectorpkg.StateBlocked {
		return string(connectorpkg.StateBlocked)
	}
	if len(blockers) > 0 {
		return checklistStatusNeedsReview
	}
	return string(state)
}

func connectorFutureExtensionStatus(registry connectorpkg.Registry) string {
	if len(registry.Diagnostics) > 0 {
		return checklistStatusBlocked
	}
	return "policy"
}

func connectorConfiguredSignal(registry connectorpkg.Registry, instances []connectorpkg.Instance, connectorType string, label string) string {
	entries := connectorRegistryEntries(registry, connectorType)
	typed := connectorInstancesByType(instances, connectorType)
	if len(typed) > 0 {
		return strconv.Itoa(len(typed)) + " configured " + label + " instance(s): " + strings.Join(connectorInstanceNames(typed), ", ")
	}
	if len(entries) == 0 {
		return "No configured " + label + " instance and no " + label + " example manifest loaded."
	}
	return "No configured " + label + " instance. " + strconv.Itoa(len(entries)) + " committed example manifest(s) are available but are not agency configuration: " + strings.Join(connectorRegistryEntryIDs(entries), ", ")
}

func connectorDryRunSignal(registry connectorpkg.Registry, instances []connectorpkg.Instance, connectorType string, command string) string {
	entries := connectorRegistryEntries(registry, connectorType)
	if len(registry.Diagnostics) > 0 {
		return "Blocked until registry diagnostics are fixed."
	}
	if len(entries) == 0 {
		return "Blocked until a committed synthetic manifest is loaded."
	}
	typed := connectorInstancesByType(instances, connectorType)
	if len(typed) == 0 {
		return "No agency dry-run recorded. Configure an instance first; example-only check is " + command
	}
	return "Configured instance dry-run status: " + strings.Join(connectorInstanceDryRunStatuses(typed), ", ")
}

func connectorInstancesByType(instances []connectorpkg.Instance, connectorType string) []connectorpkg.Instance {
	var typed []connectorpkg.Instance
	for _, instance := range instances {
		if instance.ConnectorType == connectorType {
			typed = append(typed, instance)
		}
	}
	return typed
}

func strongestConnectorState(instances []connectorpkg.Instance) connectorpkg.InstanceState {
	best := connectorpkg.StateConfiguredNotTested
	rank := map[connectorpkg.InstanceState]int{
		connectorpkg.StateBlocked:             70,
		connectorpkg.StateActive:              60,
		connectorpkg.StateReadyForActivation:  50,
		connectorpkg.StateDryRunPassed:        40,
		connectorpkg.StateConfiguredNotTested: 30,
		connectorpkg.StateNotConfigured:       20,
		connectorpkg.StateExampleAvailable:    10,
	}
	for _, instance := range instances {
		if rank[instance.State] > rank[best] {
			best = instance.State
		}
	}
	return best
}

func connectorInstanceNames(instances []connectorpkg.Instance) []string {
	names := make([]string, 0, len(instances))
	for _, instance := range instances {
		names = append(names, firstNonEmpty(instance.DisplayName, instance.ConnectorKind, "connector"))
	}
	return names
}

func connectorInstanceDryRunStatuses(instances []connectorpkg.Instance) []string {
	statuses := make([]string, 0, len(instances))
	for _, instance := range instances {
		statuses = append(statuses, firstNonEmpty(instance.DisplayName, instance.ConnectorKind, "connector")+"="+firstNonEmpty(instance.DryRunStatus, "not_run"))
	}
	return statuses
}

func connectorRedactionSignal(registry connectorpkg.Registry, connectorType string) string {
	entries := connectorRegistryEntries(registry, connectorType)
	if len(entries) == 0 {
		return "No redaction policy loaded."
	}
	return "Loaded policies use secret storage " + strings.Join(connectorRegistrySecretStorage(entries), "/") + " with bounded diagnostics."
}

func connectorRegistryRedactionStatus(registry connectorpkg.Registry) string {
	if len(registry.Diagnostics) > 0 {
		return "Registry diagnostics must be fixed before redaction review is trusted."
	}
	return "Manifest validation rejects secrets, private endpoints, private paths, raw commands, status mutation, and unsupported claims."
}

func connectorTelemetryBlockers(page operationsPage, registry connectorpkg.Registry) []string {
	blockers := connectorRegistryTypeBlockers(registry, connectorpkg.TypeTelemetrySource, "telemetry source")
	if page.DeviceError != "" {
		blockers = append(blockers, "Device binding store is unavailable; open Devices before real telemetry setup.")
	}
	if len(page.Devices) == 0 {
		blockers = append(blockers, "No active device binding is visible; bind devices before real telemetry ingest.")
	}
	if page.TelemetryError != "" {
		blockers = append(blockers, "Telemetry diagnostics are unavailable; review telemetry before connector setup.")
	}
	return blockers
}

func connectorPredictionBlockers(page operationsPage, registry connectorpkg.Registry) []string {
	blockers := connectorRegistryTypeBlockers(registry, connectorpkg.TypePrediction, "prediction")
	if page.ActiveFeedVersion == "" {
		blockers = append(blockers, "No active GTFS feed version is visible; prediction input review is incomplete.")
	}
	return blockers
}

func connectorValidatorBlockers(page operationsPage, registry connectorpkg.Registry) []string {
	blockers := connectorRegistryTypeBlockers(registry, connectorpkg.TypeValidator, "validator")
	status := strings.TrimSpace(page.ValidationHealth.OverallStatus)
	switch status {
	case "", compliance.ValidationHealthStatusConfigured, compliance.ValidationHealthStatusInstalled, compliance.ValidationHealthStatusRecorded, compliance.ValidationHealthStatusRunnable, compliance.ValidationHealthStatusConfiguredForTests, compliance.ValidationHealthStatusStub, compliance.ValidationHealthStatusSkipped:
	default:
		blockers = append(blockers, "Validation health is "+status+"; open Validation Health before relying on validator connectors.")
	}
	return blockers
}

func connectorMonitoringBlockers(page operationsPage, registry connectorpkg.Registry) []string {
	blockers := connectorRegistryTypeBlockers(registry, connectorpkg.TypeMonitoringExport, "monitoring export")
	if page.ReliabilityError != "" {
		blockers = append(blockers, "Reliability summaries are unavailable; open Maintenance before enabling monitoring export.")
	}
	return blockers
}

func connectorConsumerDiscoveryBlockers(page operationsPage, registry connectorpkg.Registry) []string {
	blockers := connectorRegistryTypeBlockers(registry, connectorpkg.TypeConsumerDiscovery, "consumer discovery")
	if page.DiscoveryError != "" {
		blockers = append(blockers, "Feed discovery metadata is unavailable; open Feeds before sharing preparation.")
	}
	if strings.TrimSpace(page.Discovery.PublicBaseURL) == "" {
		blockers = append(blockers, "Public base URL is not configured for discovery review.")
	}
	if strings.TrimSpace(page.Discovery.TechnicalContactEmail) == "" {
		blockers = append(blockers, "Technical contact metadata is missing for discovery review.")
	}
	if strings.TrimSpace(page.Discovery.License.Name) == "" {
		blockers = append(blockers, "License metadata is missing for discovery review.")
	}
	return blockers
}

func connectorRegistryTypeBlockers(registry connectorpkg.Registry, connectorType string, label string) []string {
	blockers := connectorRegistryBlockers(registry)
	if len(connectorRegistryEntries(registry, connectorType)) == 0 {
		blockers = append(blockers, "No "+label+" example manifest loaded.")
	}
	return blockers
}

func connectorRegistryBlockers(registry connectorpkg.Registry) []string {
	if len(registry.Diagnostics) == 0 {
		return nil
	}
	blockers := make([]string, 0, len(registry.Diagnostics))
	for _, diagnostic := range registry.Diagnostics {
		blockers = append(blockers, "Registry diagnostic "+firstNonEmpty(diagnostic.Code, diagnostic.Level)+" needs review.")
	}
	return blockers
}

func connectorRegistryEntries(registry connectorpkg.Registry, connectorType string) []connectorpkg.RegistryEntry {
	var entries []connectorpkg.RegistryEntry
	for _, entry := range registry.Entries {
		if entry.ConnectorType == connectorType {
			entries = append(entries, entry)
		}
	}
	return entries
}

func connectorRegistryEntryIDs(entries []connectorpkg.RegistryEntry) []string {
	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		ids = append(ids, firstNonEmpty(entry.ConnectorID, "connector"))
	}
	return ids
}

func connectorRegistrySecretStorage(entries []connectorpkg.RegistryEntry) []string {
	seen := map[string]bool{}
	values := make([]string, 0, len(entries))
	for _, entry := range entries {
		value := firstNonEmpty(entry.RedactionPolicy.SecretStorage, "none")
		if seen[value] {
			continue
		}
		seen[value] = true
		values = append(values, value)
	}
	return values
}

func connectorCatalogRows() []connectorCatalogRow {
	return []connectorCatalogRow{
		connectorCatalogRowView(
			"csv_replay_adapter",
			"Vehicle / GPS / AVL",
			"CSV replay adapter",
			"covered",
			"`examples/connectors/telemetry-csv-replay` with synthetic or operator-redacted rows.",
			"/admin/operations/connectors/workbench and /admin/operations/telemetry-simulator",
			"make test-connector-examples",
			"Production data quality, real device proof, vendor compatibility, hardware certification, compliance, consumer acceptance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/vehicle-avl-starter-kits.md", "docs/tutorials/telemetry-simulator-and-device-trial.md"},
		),
		connectorCatalogRowView(
			"http_polling_adapter",
			"Vehicle / GPS / AVL",
			"HTTP polling adapter",
			"covered",
			"`examples/connectors/telemetry-http-poller` with committed synthetic observations.",
			"/admin/operations/connectors/workbench",
			"make external-connection-check",
			"Live endpoint behavior, named API support, vendor compatibility, production AVL reliability, compliance, or public launch.",
			[]string{"docs/connectors/catalog.md", "docs/integration-adapter-kit.md", "docs/tutorials/device-avl-integration.md"},
		),
		connectorCatalogRowView(
			"webhook_sidecar_adapter",
			"Vehicle / GPS / AVL",
			"Webhook sidecar adapter",
			"covered",
			"`examples/connectors/telemetry-webhook-sidecar` as a dry-run transform shape.",
			"/admin/operations/connectors/workbench and /admin/operations/devices",
			"go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance",
			"Vendor support, hardware certification, production AVL reliability, SLA, compliance, consumer acceptance, or agency approval.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/vehicle-avl-starter-kits.md", "docs/evidence/redaction-policy.md"},
		),
		connectorCatalogRowView(
			"generic_json_transform_adapter",
			"Vehicle / GPS / AVL",
			"Generic JSON transform adapter",
			"covered",
			"`examples/connectors/generic-json-transform` with explicit field mapping and send disabled.",
			"/admin/operations/connectors/workbench and /admin/operations/realtime",
			"go test ./examples/connectors/generic-json-transform",
			"Vendor compatibility, real AVL reliability, production readiness, consumer acceptance, compliance, or ETA quality.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/vehicle-avl-starter-kits.md", "docs/connectors/redaction-first-recipes.md"},
		),
		connectorCatalogRowView(
			"vendor_shaped_synthetic_examples",
			"Vehicle / GPS / AVL",
			"Vendor-shaped synthetic examples",
			"covered",
			"`testdata/avl-vendor` fixture families for mapping, stale, future, duplicate, and low-quality cases.",
			"/admin/operations/connectors/workbench",
			"go test ./cmd/avl-vendor-adapter",
			"Named vendor support, hardware certification, real source behavior, production AVL reliability, or agency approval.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/vehicle-avl-starter-kits.md", "docs/evidence/redaction-policy.md"},
		),
		connectorCatalogRowView(
			"authenticated_telemetry_post",
			"Vehicle / GPS / AVL",
			"Authenticated telemetry POST",
			"available",
			"`POST /v1/telemetry` with a deployment-owned device bearer token.",
			"/admin/operations/devices and /admin/operations/telemetry",
			"make telemetry-simulator",
			"Real device proof, vendor compatibility, production AVL reliability, compliance, consumer acceptance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/device-avl-integration.md", "docs/tutorials/device-token-lifecycle.md"},
		),
		connectorCatalogRowView(
			"deterministic_builtin_predictor",
			"Prediction",
			"Deterministic built-in predictor",
			"available",
			"Existing deterministic prediction adapter boundary.",
			"/admin/operations/prediction-lab and /admin/operations/realtime",
			"go test ./internal/prediction",
			"Production-grade ETA quality, real-world ETA accuracy, consumer acceptance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/requirements-trip-updates.md", "docs/tutorials/prediction-eta-lab.md"},
		),
		connectorCatalogRowView(
			"external_http_predictor_adapter",
			"Prediction",
			"External HTTP predictor adapter",
			"review_only",
			"`examples/connectors/predictor-sidecar-stub` with sanitized synthetic request/response fixtures.",
			"/admin/operations/connectors/workbench and /admin/operations/prediction-lab",
			"go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance",
			"Named predictor compatibility, live service behavior, production readiness, consumer acceptance, or ETA quality.",
			[]string{"docs/connectors/catalog.md", "docs/requirements-trip-updates.md", "docs/tutorials/external-adapter-conformance.md"},
		),
		connectorCatalogRowView(
			"shadow_mode_predictor",
			"Prediction",
			"Shadow-mode predictor",
			"review_only",
			"Compare external predictor output while public Trip Updates stay on the reviewed path.",
			"/admin/operations/prediction-lab",
			"go run ./examples/connectors/predictor-sidecar-stub",
			"External predictor acceptance, production-grade ETA quality, real-world accuracy, or public-feed improvement.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/prediction-eta-lab.md", "docs/requirements-trip-updates.md"},
		),
		connectorCatalogRowView(
			"fail_closed_predictor_behavior",
			"Prediction",
			"Fail-closed predictor behavior",
			"covered",
			"Timeout, malformed, stale, wrong-agency, low-confidence, missing Vehicle Positions, and public-mutation cases.",
			"/admin/operations/connectors/tests",
			"make adapter-conformance",
			"ETA quality, real-world accuracy, named predictor compatibility, production readiness, or consumer acceptance.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/external-adapter-conformance.md", "docs/requirements-trip-updates.md"},
		),
		connectorCatalogRowView(
			"thetransitclock_candidate_notes",
			"Prediction",
			"TheTransitClock candidate notes",
			"candidate_only",
			"Treat TheTransitClock as a future candidate behind the external predictor adapter boundary only.",
			"/admin/operations/connectors/workbench",
			"make adapter-conformance",
			"TheTransitClock integration, compatibility, certification, production ETA quality, or real-world accuracy.",
			[]string{"docs/connectors/catalog.md", "docs/dependencies.md", "docs/requirements-trip-updates.md"},
		),
		connectorCatalogRowView(
			"mobilitydata_static_gtfs_validator",
			"Validator",
			"MobilityData static GTFS validator",
			"available",
			"Server-owned static GTFS validator tooling when installed by an administrator.",
			"/admin/operations/validation-health and /admin/operations/validation-center",
			"make validators-check",
			"Validator-clean proof, CAL-ITP/Caltrans compliance, consumer acceptance, public launch, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/dependencies.md", "docs/tutorials/gtfs-validation-triage.md"},
		),
		connectorCatalogRowView(
			"mobilitydata_gtfs_realtime_validator",
			"Validator",
			"MobilityData GTFS Realtime validator",
			"available",
			"Server-owned GTFS Realtime validator wrapper when installed by an administrator.",
			"/admin/operations/validation-health and /admin/operations/realtime",
			"make gtfsrt-conformance",
			"Validator-clean proof, consumer acceptance, production readiness, public launch, or compliance.",
			[]string{"docs/connectors/catalog.md", "docs/dependencies.md", "docs/tutorials/external-adapter-conformance.md"},
		),
		connectorCatalogRowView(
			"allowlisted_validator_ids",
			"Validator",
			"Allowlisted validator IDs",
			"covered",
			"`examples/connectors/validator-allowlist` with safe validator ID/feed-type mappings only.",
			"/admin/operations/connectors/tests",
			"go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance",
			"Raw validator command safety beyond the allowlist, validator-clean output, compliance, or consumer acceptance.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/plugin-contract.md", "docs/tutorials/gtfs-validation-triage.md"},
		),
		connectorCatalogRowView(
			"private_validation_health",
			"Validator",
			"Private validation health",
			"available",
			"Private health rows for validator installation, last run, blockers, and next action.",
			"/admin/operations/validation-health",
			"make validate",
			"External approval, consumer acceptance, compliance, production readiness, or public launch.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/gtfs-validation-triage.md", "docs/release-candidate-readiness.md"},
		),
		connectorCatalogRowView(
			"local_health_summaries",
			"Monitoring / Export",
			"Local health summaries",
			"available",
			"Private feed, validator, telemetry, maintenance, and reliability summaries.",
			"/admin/operations/feed-health and /admin/operations/maintenance",
			"make operations-reliability",
			"SLA coverage, uptime guarantee, hosted service availability, production readiness, or retained evidence.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/small-agency-maintenance-guide.md", "docs/runbooks/monitoring-and-alerting.md"},
		),
		connectorCatalogRowView(
			"operations_notify_draft",
			"Monitoring / Export",
			"Operations notify draft",
			"review_only",
			"Deployment-owned notification draft workflow with no send by default.",
			"/admin/operations/maintenance",
			"make operations-notify",
			"Notification delivery, incident response maturity, SLA, uptime, hosted service availability, or evidence creation.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/self-hosted-operations-notifications.md", "docs/evidence/redaction-policy.md"},
		),
		connectorCatalogRowView(
			"monitoring_export_helper",
			"Monitoring / Export",
			"Monitoring/export helper",
			"covered",
			"`examples/connectors/monitoring-export` with redacted no-send batches.",
			"/admin/operations/connectors/workbench",
			"go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance",
			"Notification delivery, SLA coverage, uptime guarantee, hosted service availability, production readiness, or retained evidence.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/self-hosted-operations-notifications.md", "docs/connectors/redaction-first-recipes.md"},
		),
		connectorCatalogRowView(
			"deployment_owned_monitoring_boundary",
			"Monitoring / Export",
			"Deployment-owned monitoring boundary",
			"needs_review",
			"Any real monitoring destination, credential, retention, and alert routing stays outside this browser page.",
			"/admin/operations/maintenance",
			"make external-connection-check",
			"Hosted service availability, paid support, SLA coverage, uptime guarantee, production readiness, or agency approval.",
			[]string{"docs/connectors/catalog.md", "docs/runbooks/monitoring-and-alerting.md", "docs/support-boundaries.md"},
		),
		connectorCatalogRowView(
			"public_feeds_json",
			"Consumer / Discovery",
			"`/public/feeds.json`",
			"available",
			"Public feed metadata endpoint for a running local or self-hosted instance.",
			"/admin/operations/feeds and /admin/operations/readiness",
			"make smoke",
			"Consumer submission, review, acceptance, ingestion, listing, display, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/release-candidate-readiness.md", "docs/consumer-submission-evidence.md"},
		),
		connectorCatalogRowView(
			"static_gtfs_url",
			"Consumer / Discovery",
			"Static GTFS URL",
			"available",
			"`/public/gtfs/schedule.zip` when an active feed version exists.",
			"/admin/operations/feeds and /admin/operations/feed-health",
			"make validate-public-feeds",
			"Validator-clean proof, consumer acceptance, public launch, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/gtfs-validation-triage.md", "docs/release-candidate-readiness.md"},
		),
		connectorCatalogRowView(
			"vehicle_positions_url",
			"Consumer / Discovery",
			"Vehicle Positions URL",
			"available",
			"`/public/gtfsrt/vehicle_positions.pb` when realtime publishing is configured.",
			"/admin/operations/realtime and /admin/operations/feed-health",
			"make gtfsrt-conformance",
			"Consumer ingestion, public display, production AVL reliability, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/telemetry-simulator-and-device-trial.md", "docs/requirements-calitp-compliance.md"},
		),
		connectorCatalogRowView(
			"trip_updates_url",
			"Consumer / Discovery",
			"Trip Updates URL",
			"available",
			"`/public/gtfsrt/trip_updates.pb` behind the pluggable prediction boundary.",
			"/admin/operations/realtime and /admin/operations/prediction-lab",
			"make gtfsrt-conformance",
			"Production-grade ETA quality, real-world accuracy, consumer acceptance, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/requirements-trip-updates.md", "docs/tutorials/prediction-eta-lab.md"},
		),
		connectorCatalogRowView(
			"alerts_url",
			"Consumer / Discovery",
			"Alerts URL",
			"available",
			"`/public/gtfsrt/alerts.pb` plus the separate private Alerts Console.",
			"/admin/operations/realtime and /admin/alerts/console",
			"make gtfsrt-conformance",
			"Consumer acceptance, complete disruption operations, public display, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/requirements-calitp-compliance.md", "docs/release-candidate-readiness.md"},
		),
		connectorCatalogRowView(
			"consumer_packet_preparedness",
			"Consumer / Discovery",
			"Consumer packet preparedness",
			"prepared_only",
			"Prepared URL and metadata review without external submission automation or status movement.",
			"/admin/operations/consumers and /admin/operations/readiness",
			"scripts/check-consumer-tracker.sh",
			"Consumer submission, review, acceptance, ingestion, listing, display, compliance, or public launch.",
			[]string{"docs/connectors/catalog.md", "docs/consumer-submission-evidence.md", "docs/evidence/consumer-submissions/README.md"},
		),
		connectorCatalogRowView(
			"manifest_based_sidecars",
			"Future Extension Model",
			"Manifest-based sidecars",
			"policy",
			"Connector manifests describe optional sidecars and command adapters; backend contracts stay explicit.",
			"/admin/operations/connectors/workbench",
			"make external-connection-check",
			"Runtime installation, automatic backend extension, vendor compatibility, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/plugin-contract.md", "docs/extension-governance.md"},
		),
		connectorCatalogRowView(
			"no_dynamic_backend_plugin_loading",
			"Future Extension Model",
			"No arbitrary dynamic backend plugin loading",
			"policy",
			safePluginDefinition,
			"/admin/operations/connectors",
			"make external-connection-check",
			"Third-party runtime distribution support, remote plugin installation, or unreviewed extension execution.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/plugin-contract.md", "docs/extension-governance.md"},
		),
		connectorCatalogRowView(
			"conformance_tests_required",
			"Future Extension Model",
			"Conformance tests required",
			"policy",
			"New connector shapes need manifest validation, synthetic fixtures, fail-closed behavior, and docs.",
			"/admin/operations/connectors/tests",
			"make adapter-conformance",
			"Real integration proof, production readiness, vendor compatibility, compliance, or consumer acceptance.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/contributing-connectors.md", "docs/tutorials/external-adapter-conformance.md"},
		),
	}
}

func connectorCatalogRowView(id string, group string, label string, status string, startWith string, browserReview string, firstCheck string, doesNotProve string, docsLinks []string) connectorCatalogRow {
	return connectorCatalogRow{
		ID:             firstNonEmpty(id, "connector_catalog_row"),
		Group:          firstNonEmpty(group, "Connector"),
		Label:          firstNonEmpty(label, "Connector row"),
		Status:         firstNonEmpty(status, checklistStatusNeedsReview),
		StartWith:      firstNonEmpty(startWith, "Review the Connector Workbench before adapting this connector shape."),
		BrowserReview:  firstNonEmpty(browserReview, "/admin/operations/connectors/workbench"),
		FirstSafeCheck: firstNonEmpty(firstCheck, "make external-connection-check"),
		DoesNotProve:   firstNonEmpty(doesNotProve, "Compatibility, compliance, production readiness, or consumer acceptance."),
		DocsLinks:      safeDocsLinks(docsLinks),
	}
}

func connectorRegistryForSection(section string) connectorpkg.Registry {
	return connectorpkg.LoadExampleRegistry()
}

func connectorCategoryView(id string, label string, status string, summary string, shape string, inputs []string, outputs []string, failure string, commands []string, adminLinks []string, docsLinks []string, boundary string) connectorCategory {
	return connectorCategory{
		ID:                 id,
		Label:              label,
		Status:             firstNonEmpty(status, checklistStatusNeedsReview),
		Summary:            firstNonEmpty(summary, "Review connector documentation before use."),
		ConnectorShape:     firstNonEmpty(shape, "sidecar, command adapter, manifest, or connector process"),
		Inputs:             cleanLaunchpadList(inputs),
		Outputs:            cleanLaunchpadList(outputs),
		FailureBehavior:    firstNonEmpty(failure, "Fail closed and keep diagnostics private."),
		CommandSuggestions: safeCommandSuggestions(commands),
		AdminLinks:         safeAdminLinks(adminLinks),
		DocsLinks:          safeDocsLinks(docsLinks),
		ClaimBoundary:      firstNonEmpty(boundary, privateBoundary()),
	}
}
