package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
	connectorpkg "open-transit-rt/internal/connectors"
)

const validatorConnectorPostMaxBytes = 32 << 10

type validatorConnectorView struct {
	GeneratedAt       time.Time                `json:"generated_at"`
	AgencyID          string                   `json:"agency_id"`
	Boundary          string                   `json:"boundary"`
	Notice            string                   `json:"notice,omitempty"`
	Error             string                   `json:"error,omitempty"`
	Instances         []connectorInstanceRow   `json:"instances"`
	Configured        []connectorInstanceRow   `json:"configured_instances"`
	Allowlisted       []validatorConnectorSpec `json:"allowlisted_validators"`
	HealthRows        []validatorConnectorRun  `json:"health_rows"`
	Rules             []validatorConnectorRule `json:"rules"`
	NextAction        string                   `json:"next_action"`
	DoesNotProve      string                   `json:"does_not_prove"`
	RawCommandBlocked string                   `json:"raw_command_blocked"`
}

type validatorConnectorSpec struct {
	ID              string   `json:"id"`
	Label           string   `json:"label"`
	FeedTypes       []string `json:"feed_types"`
	PathEnvRef      string   `json:"path_env_ref"`
	VersionEnvRef   string   `json:"version_env_ref"`
	ArgsEnvRef      string   `json:"args_env_ref,omitempty"`
	FirstCheck      string   `json:"first_check"`
	FailureBehavior string   `json:"failure_behavior"`
	DoesNotProve    string   `json:"does_not_prove"`
}

type validatorConnectorRun struct {
	FeedType           string `json:"feed_type"`
	ValidatorID        string `json:"validator_id"`
	ValidatorName      string `json:"validator_name"`
	ToolingStatus      string `json:"tooling_status"`
	ArtifactStatus     string `json:"artifact_status"`
	LatestResultStatus string `json:"latest_result_status"`
	LatestResultAt     string `json:"latest_result_at"`
	HealthStatus       string `json:"health_status"`
	NextAction         string `json:"next_action"`
	DoesNotProve       string `json:"does_not_prove"`
}

type validatorConnectorRule struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Rule    string `json:"rule"`
	Blocked string `json:"blocked"`
}

type validatorConnectorConfigMetadata struct {
	ValidatorID        string   `json:"validator_id"`
	FeedTypes          []string `json:"feed_types"`
	ToolingPathEnvRef  string   `json:"tooling_path_env_ref"`
	VersionEnvRef      string   `json:"version_env_ref,omitempty"`
	ArgsEnvRef         string   `json:"args_env_ref,omitempty"`
	TimeoutSeconds     int      `json:"timeout_seconds"`
	ServerOwnedRunOnly bool     `json:"server_owned_run_only"`
	RawCommandsBlocked bool     `json:"raw_commands_blocked"`
}

func (h *handler) renderValidatorConnector(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "validator-setup")
	page.ValidatorConnector = buildValidatorConnector(page, "", "")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "validator-setup", page)
}

func (h *handler) renderValidatorConnectorJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "validator-setup")
	page.ValidatorConnector = buildValidatorConnector(page, "", "")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.ValidatorConnector)
}

func (h *handler) operationsValidatorConnectorPost(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	if auth.RejectAgencyConflict(w, r.FormValue("agency_id"), principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "validator-setup")
	switch strings.TrimSpace(r.FormValue("action")) {
	case "save_validator_connector":
		h.operationsValidatorConnectorMetadataPost(w, r, principal, page)
	default:
		http.Error(w, "unknown connector action", http.StatusBadRequest)
	}
}

func (h *handler) operationsValidatorConnectorMetadataPost(w http.ResponseWriter, r *http.Request, principal auth.Principal, page operationsPage) {
	writer, ok := h.connectorInstances.(connectorInstanceWriter)
	if !ok || writer == nil {
		page.ValidatorConnector = buildValidatorConnector(page, "", "connector instance writer is not available in this runtime")
		renderOperationsTemplate(w, "validator-setup", page)
		return
	}
	input, err := validatorConnectorInstanceInput(r, principal)
	if err != nil {
		page.ValidatorConnector = buildValidatorConnector(page, "", err.Error())
		renderOperationsTemplate(w, "validator-setup", page)
		return
	}
	if _, err := writer.UpsertInstance(r.Context(), input); err != nil {
		page.ValidatorConnector = buildValidatorConnector(page, "", "validator connector metadata could not be stored")
		renderOperationsTemplate(w, "validator-setup", page)
		return
	}
	page = h.buildOperationsPage(r, principal, "validator-setup")
	page.ValidatorConnector = buildValidatorConnector(page, "Validator connector metadata was saved. Use Validation Health for server-owned allowlisted runs; the browser did not accept a command.", "")
	renderOperationsTemplate(w, "validator-setup", page)
}

func buildValidatorConnector(page operationsPage, notice string, errText string) validatorConnectorView {
	registry := connectorRegistryForSection("connectors")
	var instances []connectorInstanceRow
	var configured []connectorInstanceRow
	for _, row := range connectorInstanceRows(page, registry) {
		if row.ConnectorType == connectorpkg.TypeValidator {
			instances = append(instances, row)
			if row.State != string(connectorpkg.StateExampleAvailable) && row.State != string(connectorpkg.StateNotConfigured) {
				configured = append(configured, row)
			}
		}
	}
	return validatorConnectorView{
		GeneratedAt:       page.GeneratedAt,
		AgencyID:          page.AgencyID,
		Boundary:          "Private Validator Setup stores allowlisted validator IDs and deployment-owned environment reference labels only. It never stores raw commands, argument lists, private artifact paths, validator output blobs, screenshots, or evidence.",
		Notice:            notice,
		Error:             errText,
		Instances:         instances,
		Configured:        configured,
		Allowlisted:       validatorConnectorSpecs(),
		HealthRows:        validatorConnectorHealthRows(page.ValidationHealth.Feeds),
		Rules:             validatorConnectorRules(),
		NextAction:        "Choose the static or realtime validator, save only env reference labels, then review Validation Health before running the existing allowlisted server-owned action.",
		DoesNotProve:      "Validator connector configuration and validator results are supporting diagnostics only. They do not prove CAL-ITP/Caltrans compliance, consumer acceptance, final-root readiness, public launch, production readiness, hosted service, SLA, or agency approval.",
		RawCommandBlocked: "Browser forms accept validator_id and env ref labels only. Raw commands, binary paths, argv, private file paths, endpoint URLs, and output artifacts are blocked.",
	}
}

func validatorConnectorSpecs() []validatorConnectorSpec {
	return []validatorConnectorSpec{
		{
			ID:              compliance.ValidationHealthStaticValidatorID,
			Label:           "Static GTFS validator",
			FeedTypes:       []string{"schedule"},
			PathEnvRef:      "GTFS_VALIDATOR_PATH",
			VersionEnvRef:   "GTFS_VALIDATOR_VERSION",
			FirstCheck:      "make validators-check",
			FailureBehavior: "Missing or misconfigured tooling records private missing_tooling or misconfigured_tooling status.",
			DoesNotProve:    "Validator pass does not prove source-of-truth listing, compliance, or consumer acceptance.",
		},
		{
			ID:              compliance.ValidationHealthRealtimeValidatorID,
			Label:           "Realtime GTFS-RT validator",
			FeedTypes:       []string{"vehicle_positions", "trip_updates", "alerts"},
			PathEnvRef:      "GTFS_RT_VALIDATOR_PATH",
			VersionEnvRef:   "GTFS_RT_VALIDATOR_VERSION",
			ArgsEnvRef:      "GTFS_RT_VALIDATOR_ARGS",
			FirstCheck:      "make gtfsrt-conformance",
			FailureBehavior: "Missing artifacts, malformed protobuf, or tooling failures record private diagnostics without creating compliance proof.",
			DoesNotProve:    "Validator pass does not prove production AVL reliability, ETA quality, public launch, or consumer acceptance.",
		},
	}
}

func validatorConnectorRules() []validatorConnectorRule {
	return []validatorConnectorRule{
		{ID: "validator_id", Label: "Allowlisted validator ID", Rule: "Must be static-mobilitydata or realtime-mobilitydata.", Blocked: "Unsupported validator IDs and uploaded validator manifests."},
		{ID: "env_refs", Label: "Tooling env refs", Rule: "Use uppercase deployment reference labels such as GTFS_VALIDATOR_PATH.", Blocked: "Raw binary paths, private directories, endpoint URLs, and credentials."},
		{ID: "execution", Label: "Execution", Rule: "Runs stay in the existing server-owned Validation Health flow.", Blocked: "Browser-supplied commands, argv, binaries, working directories, and timeouts outside bounds."},
		{ID: "status", Label: "Result meaning", Rule: "Results are private diagnostics and last-run summaries.", Blocked: "Compliance, acceptance, final-root, production, SLA, uptime, or outside-review claims."},
	}
}

func validatorConnectorHealthRows(rows []compliance.ValidationHealthRow) []validatorConnectorRun {
	out := make([]validatorConnectorRun, 0, len(rows))
	for _, row := range rows {
		out = append(out, validatorConnectorRun{
			FeedType:           row.FeedType,
			ValidatorID:        row.ValidatorID,
			ValidatorName:      row.ValidatorName,
			ToolingStatus:      row.ToolingStatus,
			ArtifactStatus:     row.ArtifactStatus,
			LatestResultStatus: row.LatestResultStatus,
			LatestResultAt:     formatTimeForText(row.LatestResultAt),
			HealthStatus:       row.HealthStatus,
			NextAction:         row.NextAction,
			DoesNotProve:       row.ClaimBoundary,
		})
	}
	return out
}

func validatorConnectorInstanceInput(r *http.Request, principal auth.Principal) (connectorpkg.UpsertInstanceInput, error) {
	spec, err := validatorConnectorSpecByID(r.FormValue("validator_id"))
	if err != nil {
		return connectorpkg.UpsertInstanceInput{}, err
	}
	displayName := strings.TrimSpace(r.FormValue("display_name"))
	if displayName == "" {
		displayName = spec.Label
	}
	owner := strings.TrimSpace(r.FormValue("owner"))
	pathRef, err := validatorConnectorRef(r.FormValue("tooling_path_env_ref"), "tooling path env ref")
	if err != nil {
		return connectorpkg.UpsertInstanceInput{}, err
	}
	versionRef := strings.TrimSpace(r.FormValue("version_env_ref"))
	if versionRef != "" {
		versionRef, err = validatorConnectorRef(versionRef, "version env ref")
		if err != nil {
			return connectorpkg.UpsertInstanceInput{}, err
		}
	}
	argsRef := strings.TrimSpace(r.FormValue("args_env_ref"))
	if argsRef != "" {
		argsRef, err = validatorConnectorRef(argsRef, "args env ref")
		if err != nil {
			return connectorpkg.UpsertInstanceInput{}, err
		}
	}
	timeoutSeconds, err := validatorConnectorTimeout(r.FormValue("timeout_seconds"))
	if err != nil {
		return connectorpkg.UpsertInstanceInput{}, err
	}
	metadata := validatorConnectorConfigMetadata{
		ValidatorID:        spec.ID,
		FeedTypes:          spec.FeedTypes,
		ToolingPathEnvRef:  pathRef,
		VersionEnvRef:      versionRef,
		ArgsEnvRef:         argsRef,
		TimeoutSeconds:     timeoutSeconds,
		ServerOwnedRunOnly: true,
		RawCommandsBlocked: true,
	}
	configRaw, err := json.Marshal(metadata)
	if err != nil {
		return connectorpkg.UpsertInstanceInput{}, fmt.Errorf("encode validator connector metadata: %w", err)
	}
	return connectorpkg.UpsertInstanceInput{
		AgencyID:      principal.AgencyID,
		ConnectorType: connectorpkg.TypeValidator,
		ConnectorKind: strings.ReplaceAll(spec.ID, "-", "_"),
		DisplayName:   displayName,
		State:         connectorpkg.StateConfiguredNotTested,
		Owner:         owner,
		ConfigJSON:    json.RawMessage(configRaw),
		DryRunStatus:  "not_run",
		ActorID:       principal.Subject,
		Now:           time.Now().UTC(),
	}, nil
}

func validatorConnectorSpecByID(value string) (validatorConnectorSpec, error) {
	value = strings.TrimSpace(value)
	for _, spec := range validatorConnectorSpecs() {
		if spec.ID == value {
			return spec, nil
		}
	}
	return validatorConnectorSpec{}, fmt.Errorf("validator_id is not allowlisted")
}

func validatorConnectorRef(value string, label string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("%s is required", label)
	}
	lower := strings.ToLower(value)
	if strings.Contains(lower, "://") ||
		strings.ContainsAny(value, "\x00\r\n\t /\\") ||
		strings.Contains(lower, " -") ||
		strings.Contains(lower, "java ") ||
		strings.Contains(lower, ".jar") ||
		strings.Contains(lower, "password=") ||
		strings.Contains(lower, "token=") ||
		strings.Contains(lower, "secret=") {
		return "", fmt.Errorf("%s must be an uppercase deployment reference label, not a command, path, URL, or inline secret", label)
	}
	if err := connectorpkg.ValidateSecretRef(value); err != nil {
		return "", fmt.Errorf("%s must be an uppercase deployment reference label", label)
	}
	return value, nil
}

func validatorConnectorTimeout(value string) (int, error) {
	if strings.TrimSpace(value) == "" {
		return int((2 * time.Minute).Seconds()), nil
	}
	timeoutSeconds, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || timeoutSeconds < 1 || timeoutSeconds > 600 {
		return 0, fmt.Errorf("timeout seconds must be between 1 and 600")
	}
	return timeoutSeconds, nil
}
