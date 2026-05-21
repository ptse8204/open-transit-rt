package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	connectorpkg "open-transit-rt/internal/connectors"
	"open-transit-rt/internal/prediction"
)

const predictionConnectorPostMaxBytes = 32 << 10

type predictionConnectorView struct {
	GeneratedAt                  time.Time                    `json:"generated_at"`
	AgencyID                     string                       `json:"agency_id"`
	Boundary                     string                       `json:"boundary"`
	Notice                       string                       `json:"notice,omitempty"`
	Error                        string                       `json:"error,omitempty"`
	CurrentMode                  predictionConnectorModeState `json:"current_mode"`
	Modes                        []predictionConnectorMode    `json:"modes"`
	Instances                    []connectorInstanceRow       `json:"instances"`
	Configured                   []connectorInstanceRow       `json:"configured_instances"`
	WithheldReasons              []predictionLabReason        `json:"withheld_reasons"`
	VehiclePositionsIndependence string                       `json:"vehicle_positions_independence"`
	ExternalHTTPRules            []predictionConnectorRule    `json:"external_http_rules"`
	NextAction                   string                       `json:"next_action"`
	DoesNotProve                 string                       `json:"does_not_prove"`
}

type predictionConnectorModeState struct {
	Mode             string `json:"mode"`
	AdapterName      string `json:"adapter_name"`
	State            string `json:"state"`
	PublicOutput     string `json:"public_output"`
	VehiclePositions string `json:"vehicle_positions"`
	WithheldSignal   string `json:"withheld_signal"`
	NextAction       string `json:"next_action"`
}

type predictionConnectorMode struct {
	ID              string `json:"id"`
	Label           string `json:"label"`
	Summary         string `json:"summary"`
	PublicOutput    string `json:"public_output"`
	FailureBehavior string `json:"failure_behavior"`
	FirstCheck      string `json:"first_check"`
	DoesNotEnable   string `json:"does_not_enable"`
}

type predictionConnectorRule struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	Required     string `json:"required"`
	BrowserKeeps string `json:"browser_keeps"`
}

type predictionConnectorConfigMetadata struct {
	Mode                        string `json:"mode"`
	AdapterName                 string `json:"adapter_name"`
	EndpointURLEnvRef           string `json:"endpoint_url_env_ref,omitempty"`
	AllowedHostsEnvRef          string `json:"allowed_hosts_env_ref,omitempty"`
	Path                        string `json:"path,omitempty"`
	TimeoutSeconds              int    `json:"timeout_seconds,omitempty"`
	ShadowMode                  bool   `json:"shadow_mode"`
	FailClosed                  bool   `json:"fail_closed"`
	VehiclePositionsIndependent bool   `json:"vehicle_positions_independent"`
	PublicOutputSource          string `json:"public_output_source"`
}

func (h *handler) renderPredictionConnector(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "prediction-setup")
	page.PredictionConnector = buildPredictionConnector(page, "", "")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "prediction-setup", page)
}

func (h *handler) renderPredictionConnectorJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "prediction-setup")
	page.PredictionConnector = buildPredictionConnector(page, "", "")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.PredictionConnector)
}

func (h *handler) operationsPredictionConnectorPost(w http.ResponseWriter, r *http.Request) {
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
	page := h.buildOperationsPage(r, principal, "prediction-setup")
	switch strings.TrimSpace(r.FormValue("action")) {
	case "save_prediction_connector":
		h.operationsPredictionConnectorMetadataPost(w, r, principal, page)
	default:
		http.Error(w, "unknown connector action", http.StatusBadRequest)
	}
}

func (h *handler) operationsPredictionConnectorMetadataPost(w http.ResponseWriter, r *http.Request, principal auth.Principal, page operationsPage) {
	writer, ok := h.connectorInstances.(connectorInstanceWriter)
	if !ok || writer == nil {
		page.PredictionConnector = buildPredictionConnector(page, "", "connector instance writer is not available in this runtime")
		renderOperationsTemplate(w, "prediction-setup", page)
		return
	}
	input, err := predictionConnectorInstanceInput(r, principal)
	if err != nil {
		page.PredictionConnector = buildPredictionConnector(page, "", err.Error())
		renderOperationsTemplate(w, "prediction-setup", page)
		return
	}
	if _, err := writer.UpsertInstance(r.Context(), input); err != nil {
		page.PredictionConnector = buildPredictionConnector(page, "", "prediction connector metadata could not be stored")
		renderOperationsTemplate(w, "prediction-setup", page)
		return
	}
	page = h.buildOperationsPage(r, principal, "prediction-setup")
	page.PredictionConnector = buildPredictionConnector(page, "Prediction connector metadata was saved. Deployment-owned env configuration is still required before any external predictor can affect Trip Updates.", "")
	renderOperationsTemplate(w, "prediction-setup", page)
}

func buildPredictionConnector(page operationsPage, notice string, errText string) predictionConnectorView {
	registry := connectorRegistryForSection("connectors")
	var instances []connectorInstanceRow
	var configured []connectorInstanceRow
	for _, row := range connectorInstanceRows(page, registry) {
		if row.ConnectorType == connectorpkg.TypePrediction {
			instances = append(instances, row)
			if row.State != string(connectorpkg.StateExampleAvailable) && row.State != string(connectorpkg.StateNotConfigured) {
				configured = append(configured, row)
			}
		}
	}
	return predictionConnectorView{
		GeneratedAt:                  page.GeneratedAt,
		AgencyID:                     page.AgencyID,
		Boundary:                     "Private Prediction Setup stores redacted metadata and deployment-owned environment reference labels only. It does not store predictor URLs, tokens, raw sidecar payloads, browser-run commands, screenshots, or public evidence.",
		Notice:                       notice,
		Error:                        errText,
		CurrentMode:                  predictionConnectorCurrentMode(page, configured),
		Modes:                        predictionConnectorModes(),
		Instances:                    instances,
		Configured:                   configured,
		WithheldReasons:              page.PredictionLab.WithheldReasons,
		VehiclePositionsIndependence: "Vehicle Positions publishing stays independent of optional prediction sidecars. A missing, failed, shadowed, or blocked prediction connector must not stop Vehicle Positions.",
		ExternalHTTPRules:            predictionConnectorRules(),
		NextAction:                   "Keep deterministic prediction as the default, save only deployment-owned reference labels for any external HTTP sidecar, then use Prediction Lab and adapter conformance before changing runtime env.",
		DoesNotProve:                 "Prediction connector configuration does not prove production-grade ETA quality, real-world ETA accuracy, consumer acceptance, named predictor compatibility, compliance, uptime, or production readiness.",
	}
}

func predictionConnectorCurrentMode(page operationsPage, rows []connectorInstanceRow) predictionConnectorModeState {
	state := predictionConnectorModeState{
		Mode:             "deterministic_default",
		AdapterName:      firstNonEmpty(page.TripUpdatesQuality.AdapterName, prediction.AdapterNameDeterministic),
		State:            string(connectorpkg.StateNotConfigured),
		PublicOutput:     "Trip Updates use the internal deterministic predictor when configured by runtime env; otherwise valid empty or withheld output remains acceptable.",
		VehiclePositions: "independent",
		WithheldSignal:   predictionConnectorWithheldSignal(page.PredictionLab.WithheldReasons),
		NextAction:       "Review Prediction Lab withheld reasons before considering external HTTP sidecar evaluation.",
	}
	for _, row := range rows {
		if predictionConnectorRowRank(row.State) < predictionConnectorRowRank(state.State) {
			continue
		}
		state.Mode = firstNonEmpty(row.ConnectorKind, state.Mode)
		state.State = row.State
		state.AdapterName = predictionConnectorAdapterName(row.ConnectorKind)
		state.PublicOutput = predictionConnectorPublicOutput(row.ConnectorKind)
		state.NextAction = row.NextAction
	}
	return state
}

func predictionConnectorRowRank(state string) int {
	switch state {
	case string(connectorpkg.StateBlocked):
		return 70
	case string(connectorpkg.StateActive):
		return 60
	case string(connectorpkg.StateReadyForActivation):
		return 50
	case string(connectorpkg.StateDryRunPassed):
		return 40
	case string(connectorpkg.StateConfiguredNotTested):
		return 30
	default:
		return 10
	}
}

func predictionConnectorAdapterName(mode string) string {
	switch mode {
	case "external_http_shadow":
		return prediction.AdapterNameExternalHTTPShadow
	case "external_http_fail_closed":
		return prediction.AdapterNameExternalHTTP
	case "deterministic_default":
		return prediction.AdapterNameDeterministic
	default:
		return firstNonEmpty(mode, prediction.AdapterNameDeterministic)
	}
}

func predictionConnectorPublicOutput(mode string) string {
	switch mode {
	case "external_http_shadow":
		return "deterministic Trip Updates remain public; external HTTP output is diagnostics-only shadow data"
	case "external_http_fail_closed":
		return "external HTTP can be selected by deployment env later; failures must produce valid empty Trip Updates with diagnostics"
	case "deterministic_default":
		return "internal deterministic prediction remains the default public Trip Updates source"
	default:
		return "public output stays behind the prediction adapter boundary"
	}
}

func predictionConnectorWithheldSignal(reasons []predictionLabReason) string {
	if len(reasons) == 0 {
		return "no withheld reasons recorded yet"
	}
	parts := make([]string, 0, len(reasons))
	for _, reason := range reasons {
		parts = append(parts, firstNonEmpty(reason.Label, reason.Reason)+"="+strconv.Itoa(reason.Count))
	}
	return strings.Join(parts, ", ")
}

func predictionConnectorModes() []predictionConnectorMode {
	return []predictionConnectorMode{
		{
			ID:              "deterministic_default",
			Label:           "Deterministic default",
			Summary:         "Use the built-in deterministic adapter and conservative withholding rules.",
			PublicOutput:    "Internal deterministic output or valid empty/withheld Trip Updates.",
			FailureBehavior: "No external sidecar dependency.",
			FirstCheck:      "go test ./internal/prediction",
			DoesNotEnable:   "Does not enable a live external predictor or prove ETA quality.",
		},
		{
			ID:              "external_http_shadow",
			Label:           "External HTTP shadow",
			Summary:         "Evaluate an external HTTP sidecar without replacing public Trip Updates.",
			PublicOutput:    "Deterministic Trip Updates remain public; sidecar result is bounded diagnostics.",
			FailureBehavior: "Shadow failure is recorded and deterministic fallback remains public.",
			FirstCheck:      "go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance",
			DoesNotEnable:   "Does not contact a sidecar from the browser or prove named predictor compatibility.",
		},
		{
			ID:              "external_http_fail_closed",
			Label:           "External HTTP fail-closed",
			Summary:         "Prepare deployment-owned external_http runtime configuration with strict endpoint and byte/time caps.",
			PublicOutput:    "Only deployment env can select external_http; failures emit valid empty Trip Updates with diagnostics.",
			FailureBehavior: "Timeout, malformed response, wrong agency, low confidence, or unsafe output fails closed.",
			FirstCheck:      "go test ./internal/prediction -run ExternalHTTP",
			DoesNotEnable:   "Does not make external prediction the default and does not prove public ETA quality.",
		},
	}
}

func predictionConnectorRules() []predictionConnectorRule {
	return []predictionConnectorRule{
		{ID: "endpoint_url_env_ref", Label: "Endpoint URL env ref", Required: "Uppercase deployment label such as TRIP_UPDATES_EXTERNAL_HTTP_URL; value stays outside the browser.", BrowserKeeps: "reference label only"},
		{ID: "allowed_hosts_env_ref", Label: "Allowed hosts env ref", Required: "Uppercase deployment label such as TRIP_UPDATES_EXTERNAL_HTTP_ALLOWED_HOSTS.", BrowserKeeps: "reference label only"},
		{ID: "path", Label: "Exact HTTP path", Required: prediction.ExternalHTTPTripUpdatesPath, BrowserKeeps: "fixed path string"},
		{ID: "token_ref", Label: "Bearer token ref", Required: "Optional uppercase deployment secret label; no token value.", BrowserKeeps: "secret_refs entry only"},
		{ID: "timeout", Label: "Timeout", Required: "1-30 seconds; default 2 seconds.", BrowserKeeps: "bounded integer"},
	}
}

func predictionConnectorInstanceInput(r *http.Request, principal auth.Principal) (connectorpkg.UpsertInstanceInput, error) {
	mode := strings.TrimSpace(r.FormValue("mode"))
	if !predictionConnectorModeAllowed(mode) {
		return connectorpkg.UpsertInstanceInput{}, fmt.Errorf("unsupported prediction connector mode")
	}
	displayName := strings.TrimSpace(r.FormValue("display_name"))
	if displayName == "" {
		displayName = predictionConnectorDefaultName(mode)
	}
	owner := strings.TrimSpace(r.FormValue("owner"))
	metadata := predictionConnectorConfigMetadata{
		Mode:                        mode,
		AdapterName:                 predictionConnectorAdapterName(mode),
		ShadowMode:                  mode == "external_http_shadow",
		FailClosed:                  mode == "external_http_fail_closed",
		VehiclePositionsIndependent: true,
		PublicOutputSource:          predictionConnectorPublicOutput(mode),
	}
	var secretRefs []string
	if mode != "deterministic_default" {
		endpointRef, err := predictionConnectorRef(r.FormValue("endpoint_url_env_ref"), "endpoint URL env ref")
		if err != nil {
			return connectorpkg.UpsertInstanceInput{}, err
		}
		allowedHostsRef, err := predictionConnectorRef(r.FormValue("allowed_hosts_env_ref"), "allowed hosts env ref")
		if err != nil {
			return connectorpkg.UpsertInstanceInput{}, err
		}
		path := strings.TrimSpace(r.FormValue("path"))
		if path != prediction.ExternalHTTPTripUpdatesPath {
			return connectorpkg.UpsertInstanceInput{}, fmt.Errorf("external HTTP predictor path must be exactly %s", prediction.ExternalHTTPTripUpdatesPath)
		}
		timeoutSeconds, err := predictionConnectorTimeout(r.FormValue("timeout_seconds"))
		if err != nil {
			return connectorpkg.UpsertInstanceInput{}, err
		}
		metadata.EndpointURLEnvRef = endpointRef
		metadata.AllowedHostsEnvRef = allowedHostsRef
		metadata.Path = path
		metadata.TimeoutSeconds = timeoutSeconds
		tokenRef := strings.TrimSpace(r.FormValue("token_ref"))
		if tokenRef != "" {
			ref, err := predictionConnectorRef(tokenRef, "token ref")
			if err != nil {
				return connectorpkg.UpsertInstanceInput{}, err
			}
			secretRefs = append(secretRefs, ref)
		}
	}
	configRaw, err := json.Marshal(metadata)
	if err != nil {
		return connectorpkg.UpsertInstanceInput{}, fmt.Errorf("encode prediction connector metadata: %w", err)
	}
	return connectorpkg.UpsertInstanceInput{
		AgencyID:      principal.AgencyID,
		ConnectorType: connectorpkg.TypePrediction,
		ConnectorKind: mode,
		DisplayName:   displayName,
		State:         connectorpkg.StateConfiguredNotTested,
		Owner:         owner,
		ConfigJSON:    json.RawMessage(configRaw),
		SecretRefs:    secretRefs,
		DryRunStatus:  "not_run",
		ActorID:       principal.Subject,
		Now:           time.Now().UTC(),
	}, nil
}

func predictionConnectorModeAllowed(mode string) bool {
	switch mode {
	case "deterministic_default", "external_http_shadow", "external_http_fail_closed":
		return true
	default:
		return false
	}
}

func predictionConnectorDefaultName(mode string) string {
	switch mode {
	case "external_http_shadow":
		return "External HTTP prediction shadow"
	case "external_http_fail_closed":
		return "External HTTP prediction fail-closed"
	default:
		return "Deterministic prediction"
	}
}

func predictionConnectorRef(value string, label string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("%s is required", label)
	}
	if predictionConnectorUnsafeValue(value) {
		return "", fmt.Errorf("%s must be a deployment reference label, not a URL, token, or inline secret", label)
	}
	if err := connectorpkg.ValidateSecretRef(value); err != nil {
		return "", fmt.Errorf("%s must be an uppercase deployment reference label", label)
	}
	return value, nil
}

func predictionConnectorTimeout(value string) (int, error) {
	if strings.TrimSpace(value) == "" {
		return 2, nil
	}
	timeoutSeconds, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || timeoutSeconds < 1 || timeoutSeconds > 30 {
		return 0, fmt.Errorf("timeout seconds must be between 1 and 30")
	}
	return timeoutSeconds, nil
}

func predictionConnectorUnsafeValue(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	return strings.Contains(lower, "://") ||
		strings.ContainsAny(value, "\x00\r\n\t") ||
		strings.Contains(lower, "password=") ||
		strings.Contains(lower, "token=") ||
		strings.Contains(lower, "secret=")
}
