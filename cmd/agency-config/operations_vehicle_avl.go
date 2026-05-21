package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	connectorpkg "open-transit-rt/internal/connectors"
)

const vehicleAVLPostMaxBytes = 64 << 10

type vehicleAVLSetupView struct {
	GeneratedAt    time.Time                 `json:"generated_at"`
	AgencyID       string                    `json:"agency_id"`
	Boundary       string                    `json:"boundary"`
	Notice         string                    `json:"notice,omitempty"`
	Error          string                    `json:"error,omitempty"`
	SourceShapes   []vehicleAVLShape         `json:"source_shapes"`
	FieldMappings  []vehicleAVLField         `json:"field_mappings"`
	Instances      []connectorInstanceRow    `json:"instances"`
	Configured     []connectorInstanceRow    `json:"configured_instances"`
	DryRuns        []vehicleAVLDryRunRow     `json:"dry_runs"`
	DryRunError    string                    `json:"dry_run_error,omitempty"`
	Activation     []vehicleAVLActivationRow `json:"activation"`
	ReadyInstances []connectorInstanceRow    `json:"ready_instances"`
	NextAction     string                    `json:"next_action"`
	DryRunBoundary string                    `json:"dry_run_boundary"`
	ActivationGate string                    `json:"activation_gate"`
	DoesNotProve   string                    `json:"does_not_prove"`
}

type vehicleAVLDryRunRow struct {
	ID          int64  `json:"id"`
	InstanceID  int64  `json:"instance_id"`
	Connector   string `json:"connector"`
	Status      string `json:"status"`
	Counts      string `json:"counts"`
	Redaction   string `json:"redaction"`
	Summary     string `json:"summary"`
	FinishedAt  string `json:"finished_at"`
	DoesNotKeep string `json:"does_not_keep"`
}

type vehicleAVLActivationRow struct {
	InstanceID    int64  `json:"instance_id"`
	Connector     string `json:"connector"`
	CheckID       string `json:"check_id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
}

type vehicleAVLShape struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Summary    string `json:"summary"`
	FirstCheck string `json:"first_check"`
	DoesNotRun string `json:"does_not_run"`
}

type vehicleAVLField struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
	Example  string `json:"example"`
}

type vehicleAVLConfigMetadata struct {
	SourceShape string            `json:"source_shape"`
	FieldMap    map[string]string `json:"field_map"`
	Safety      map[string]any    `json:"safety"`
}

type connectorInstanceWriter interface {
	UpsertInstance(ctx context.Context, input connectorpkg.UpsertInstanceInput) (connectorpkg.Instance, error)
}

type connectorDryRunWriter interface {
	CreateDryRunJob(ctx context.Context, input connectorpkg.CreateDryRunJobInput) (connectorpkg.DryRunJob, error)
}

type connectorInstanceStateWriter interface {
	UpdateInstanceState(ctx context.Context, input connectorpkg.UpdateInstanceStateInput) (connectorpkg.Instance, error)
}

func (h *handler) renderVehicleAVLSetup(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "vehicle-avl-setup")
	page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", "")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "vehicle-avl-setup", page)
}

func (h *handler) renderVehicleAVLSetupJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "vehicle-avl-setup")
	page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", "")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.VehicleAVLSetup)
}

func (h *handler) operationsVehicleAVLPost(w http.ResponseWriter, r *http.Request) {
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
	page := h.buildOperationsPage(r, principal, "vehicle-avl-setup")
	switch strings.TrimSpace(r.FormValue("action")) {
	case "save_vehicle_avl_connector":
		h.operationsVehicleAVLMetadataPost(w, r, principal, page)
	case "record_vehicle_avl_dry_run":
		h.operationsVehicleAVLDryRunPost(w, r, principal, page)
	case "mark_vehicle_avl_ready":
		h.operationsVehicleAVLReadyPost(w, r, principal, page)
	default:
		http.Error(w, "unknown connector action", http.StatusBadRequest)
	}
}

func (h *handler) operationsVehicleAVLReadyPost(w http.ResponseWriter, r *http.Request, principal auth.Principal, page operationsPage) {
	writer, ok := h.connectorInstances.(connectorInstanceStateWriter)
	if !ok || writer == nil {
		page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", "connector state writer is not available in this runtime")
		renderOperationsTemplate(w, "vehicle-avl-setup", page)
		return
	}
	instanceID, err := strconv.ParseInt(strings.TrimSpace(r.FormValue("connector_instance_id")), 10, 64)
	if err != nil || instanceID <= 0 {
		page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", "connector instance is required")
		renderOperationsTemplate(w, "vehicle-avl-setup", page)
		return
	}
	if !vehicleAVLReadyInstanceIDs(page)[instanceID] {
		page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", "activation readiness checks must pass before marking a connector ready")
		renderOperationsTemplate(w, "vehicle-avl-setup", page)
		return
	}
	if _, err := writer.UpdateInstanceState(r.Context(), connectorpkg.UpdateInstanceStateInput{
		AgencyID: principal.AgencyID,
		ID:       instanceID,
		State:    connectorpkg.StateReadyForActivation,
		ActorID:  principal.Subject,
		Now:      time.Now().UTC(),
	}); err != nil {
		page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", "connector could not be marked ready for deployment-owned activation")
		renderOperationsTemplate(w, "vehicle-avl-setup", page)
		return
	}
	page = h.buildOperationsPage(r, principal, "vehicle-avl-setup")
	page.VehicleAVLSetup = buildVehicleAVLSetup(page, "Connector was marked ready for deployment-owned activation. The browser did not start a sidecar or contact an external source.", "")
	renderOperationsTemplate(w, "vehicle-avl-setup", page)
}

func (h *handler) operationsVehicleAVLMetadataPost(w http.ResponseWriter, r *http.Request, principal auth.Principal, page operationsPage) {
	writer, ok := h.connectorInstances.(connectorInstanceWriter)
	if !ok || writer == nil {
		page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", "connector instance writer is not available in this runtime")
		renderOperationsTemplate(w, "vehicle-avl-setup", page)
		return
	}
	input, err := vehicleAVLInstanceInput(r, principal)
	if err != nil {
		page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", err.Error())
		renderOperationsTemplate(w, "vehicle-avl-setup", page)
		return
	}
	if _, err := writer.UpsertInstance(r.Context(), input); err != nil {
		page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", "vehicle connector metadata could not be stored")
		renderOperationsTemplate(w, "vehicle-avl-setup", page)
		return
	}
	page = h.buildOperationsPage(r, principal, "vehicle-avl-setup")
	page.VehicleAVLSetup = buildVehicleAVLSetup(page, "Vehicle / GPS / AVL connector metadata was saved. Dry-run is still required before activation.", "")
	renderOperationsTemplate(w, "vehicle-avl-setup", page)
}

func (h *handler) operationsVehicleAVLDryRunPost(w http.ResponseWriter, r *http.Request, principal auth.Principal, page operationsPage) {
	writer, ok := h.connectorInstances.(connectorDryRunWriter)
	if !ok || writer == nil {
		page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", "connector dry-run writer is not available in this runtime")
		renderOperationsTemplate(w, "vehicle-avl-setup", page)
		return
	}
	input, err := vehicleAVLDryRunInput(r, principal)
	if err != nil {
		page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", err.Error())
		renderOperationsTemplate(w, "vehicle-avl-setup", page)
		return
	}
	if _, err := writer.CreateDryRunJob(r.Context(), input); err != nil {
		page.VehicleAVLSetup = buildVehicleAVLSetup(page, "", "vehicle connector dry-run result could not be recorded")
		renderOperationsTemplate(w, "vehicle-avl-setup", page)
		return
	}
	page = h.buildOperationsPage(r, principal, "vehicle-avl-setup")
	page.VehicleAVLSetup = buildVehicleAVLSetup(page, "Vehicle / GPS / AVL dry-run result was recorded with a redacted summary only.", "")
	renderOperationsTemplate(w, "vehicle-avl-setup", page)
}

func buildVehicleAVLSetup(page operationsPage, notice string, errText string) vehicleAVLSetupView {
	registry := connectorRegistryForSection("connectors")
	var instances []connectorInstanceRow
	var configured []connectorInstanceRow
	for _, row := range connectorInstanceRows(page, registry) {
		if row.ConnectorType == connectorpkg.TypeTelemetrySource {
			instances = append(instances, row)
			if row.State != string(connectorpkg.StateExampleAvailable) && row.State != string(connectorpkg.StateNotConfigured) {
				configured = append(configured, row)
			}
		}
	}
	return vehicleAVLSetupView{
		GeneratedAt:    page.GeneratedAt,
		AgencyID:       page.AgencyID,
		Boundary:       "Private Vehicle / GPS / AVL setup stores redacted metadata and deployment-owned secret reference labels only. It does not store secret values, raw payloads, private endpoints, screenshots, or evidence.",
		Notice:         notice,
		Error:          errText,
		SourceShapes:   vehicleAVLShapes(),
		FieldMappings:  vehicleAVLFields(),
		Instances:      instances,
		Configured:     configured,
		DryRuns:        vehicleAVLDryRunRows(page, configured),
		DryRunError:    page.ConnectorDryRunError,
		Activation:     vehicleAVLActivationRows(page, configured),
		ReadyInstances: vehicleAVLReadyInstances(page, configured),
		NextAction:     "Choose the closest source shape, map the required fields, save metadata, then run a server-owned dry-run before activation review.",
		DryRunBoundary: "Dry-run is required and remains server-owned. This browser page records configuration metadata only.",
		ActivationGate: "Activation stays blocked until mapping, dry-run, device bindings, token refs, stale/future/quality rules, and redaction checks pass.",
		DoesNotProve:   "This setup page does not prove vendor compatibility, hardware certification, production AVL reliability, consumer acceptance, compliance, uptime, or ETA quality.",
	}
}

func vehicleAVLReadyInstanceIDs(page operationsPage) map[int64]bool {
	ready := make(map[int64]bool)
	for _, row := range vehicleAVLReadyInstances(page, configuredVehicleAVLRows(page)) {
		id, err := strconv.ParseInt(row.ID, 10, 64)
		if err == nil {
			ready[id] = true
		}
	}
	return ready
}

func configuredVehicleAVLRows(page operationsPage) []connectorInstanceRow {
	registry := connectorRegistryForSection("connectors")
	var configured []connectorInstanceRow
	for _, row := range connectorInstanceRows(page, registry) {
		if row.ConnectorType == connectorpkg.TypeTelemetrySource && row.State != string(connectorpkg.StateExampleAvailable) && row.State != string(connectorpkg.StateNotConfigured) {
			configured = append(configured, row)
		}
	}
	return configured
}

func vehicleAVLReadyInstances(page operationsPage, instances []connectorInstanceRow) []connectorInstanceRow {
	readyIDs := make(map[int64]bool)
	for _, check := range vehicleAVLActivationRows(page, instances) {
		if check.Status != "passed" {
			readyIDs[check.InstanceID] = false
			continue
		}
		if _, ok := readyIDs[check.InstanceID]; !ok {
			readyIDs[check.InstanceID] = true
		}
	}
	var ready []connectorInstanceRow
	for _, row := range instances {
		id, err := strconv.ParseInt(row.ID, 10, 64)
		if err == nil && readyIDs[id] && row.State != string(connectorpkg.StateReadyForActivation) && row.State != string(connectorpkg.StateActive) {
			ready = append(ready, row)
		}
	}
	return ready
}

func vehicleAVLActivationRows(page operationsPage, instances []connectorInstanceRow) []vehicleAVLActivationRow {
	instanceByID := make(map[int64]connectorpkg.Instance)
	for _, instance := range page.ConnectorInstances {
		if instance.ConnectorType == connectorpkg.TypeTelemetrySource {
			instanceByID[instance.ID] = instance
		}
	}
	latestJob := latestDryRunJobsByInstance(page.ConnectorDryRunJobs)
	var rows []vehicleAVLActivationRow
	for _, row := range instances {
		id, err := strconv.ParseInt(row.ID, 10, 64)
		if err != nil {
			continue
		}
		instance := instanceByID[id]
		job := latestJob[id]
		rows = append(rows,
			activationRow(id, row.DisplayName, "configured", "Connector configured", row.State != string(connectorpkg.StateExampleAvailable) && row.State != string(connectorpkg.StateNotConfigured), row.State, "Save connector metadata before activation review."),
			activationRow(id, row.DisplayName, "field_map", "Required field map", vehicleAVLRequiredFieldMapValid(instance), "required telemetry fields mapped", "Map agency, device, vehicle, timestamp, latitude, longitude, and quality fields."),
			activationRow(id, row.DisplayName, "dry_run", "Dry-run passed", row.State == string(connectorpkg.StateDryRunPassed) || row.State == string(connectorpkg.StateReadyForActivation) || row.State == string(connectorpkg.StateActive), row.DryRunStatus, "Record a passed server-owned dry-run result."),
			activationRow(id, row.DisplayName, "device_bindings", "Device bindings exist", len(page.Devices) > 0, fmt.Sprintf("%d device binding(s)", len(page.Devices)), "Create or review device bindings before activation."),
			activationRow(id, row.DisplayName, "token_refs", "Token reference labels", len(instance.SecretRefs) > 0, fmt.Sprintf("%d secret reference label(s)", len(instance.SecretRefs)), "Store only deployment-owned secret reference labels, not token values."),
			activationRow(id, row.DisplayName, "target_path", "Telemetry target path", vehicleAVLTargetPathSafe(instance), "/v1/telemetry", "Use authenticated POST /v1/telemetry through a deployment-owned sidecar."),
			activationRow(id, row.DisplayName, "freshness_rules", "Stale/future/quality rules", page.StaleThreshold > 0, page.StaleThreshold.String(), "Keep stale telemetry and quality rules configured before activation."),
			activationRow(id, row.DisplayName, "redaction_scan", "Redaction scan passed", job.ID > 0 && job.RedactionScanStatus == "passed", firstNonEmpty(job.RedactionScanStatus, "not recorded"), "Record a passed redaction scan with the dry-run result."),
		)
	}
	return rows
}

func activationRow(instanceID int64, connector string, checkID string, label string, passed bool, signal string, next string) vehicleAVLActivationRow {
	status := "blocked"
	if passed {
		status = "passed"
	}
	return vehicleAVLActivationRow{
		InstanceID:    instanceID,
		Connector:     connector,
		CheckID:       checkID,
		Label:         label,
		Status:        status,
		CurrentSignal: firstNonEmpty(signal, "not configured"),
		NextAction:    next,
	}
}

func latestDryRunJobsByInstance(jobs []connectorpkg.DryRunJob) map[int64]connectorpkg.DryRunJob {
	out := make(map[int64]connectorpkg.DryRunJob)
	for _, job := range jobs {
		if existing, ok := out[job.ConnectorInstanceID]; !ok || job.CreatedAt.After(existing.CreatedAt) || job.ID > existing.ID {
			out[job.ConnectorInstanceID] = job
		}
	}
	return out
}

func vehicleAVLRequiredFieldMapValid(instance connectorpkg.Instance) bool {
	var metadata vehicleAVLConfigMetadata
	if err := json.Unmarshal(instance.ConfigJSON, &metadata); err != nil {
		return false
	}
	for _, field := range vehicleAVLFields() {
		if field.Required && strings.TrimSpace(metadata.FieldMap[field.ID]) == "" {
			return false
		}
	}
	return true
}

func vehicleAVLTargetPathSafe(instance connectorpkg.Instance) bool {
	var metadata vehicleAVLConfigMetadata
	if err := json.Unmarshal(instance.ConfigJSON, &metadata); err != nil {
		return false
	}
	value, _ := metadata.Safety["target_path"].(string)
	return value == "/v1/telemetry"
}

func vehicleAVLDryRunRows(page operationsPage, instances []connectorInstanceRow) []vehicleAVLDryRunRow {
	names := make(map[int64]string)
	for _, row := range instances {
		id, err := strconv.ParseInt(row.ID, 10, 64)
		if err == nil {
			names[id] = row.DisplayName
		}
	}
	var rows []vehicleAVLDryRunRow
	for _, job := range page.ConnectorDryRunJobs {
		if names[job.ConnectorInstanceID] == "" {
			continue
		}
		rows = append(rows, vehicleAVLDryRunRow{
			ID:          job.ID,
			InstanceID:  job.ConnectorInstanceID,
			Connector:   names[job.ConnectorInstanceID],
			Status:      job.Status,
			Counts:      fmt.Sprintf("accepted=%d rejected=%d dropped=%d", job.AcceptedCount, job.RejectedCount, job.DroppedCount),
			Redaction:   job.RedactionScanStatus,
			Summary:     dryRunSummaryText(job.RedactedSummary),
			FinishedAt:  job.FinishedAt.UTC().Format(time.RFC3339),
			DoesNotKeep: "Raw payloads are not retained by default.",
		})
	}
	return rows
}

func dryRunSummaryText(raw json.RawMessage) string {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return "summary unavailable"
	}
	if value, ok := obj["summary"].(string); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return "redacted summary recorded"
}

func vehicleAVLShapes() []vehicleAVLShape {
	return []vehicleAVLShape{
		{"generic_json_transform", "Generic JSON transform", "Map a redacted JSON payload shape into authenticated telemetry observations.", "Save metadata, then dry-run with a bounded synthetic or redacted fixture.", "Does not execute arbitrary transforms from the browser."},
		{"http_polling", "HTTP polling", "Review a deployment-owned poller shape that normalizes observations before posting telemetry.", "Save metadata with endpoint refs only, then dry-run server-side.", "Does not store URLs, credentials, or start polling."},
		{"webhook_sidecar", "Webhook receiver / sidecar", "Review a sidecar that accepts vendor payloads and posts normalized telemetry.", "Save field map and token ref labels, then dry-run sidecar fixtures.", "Does not expose public receivers or accept raw payloads here."},
		{"csv_replay", "CSV replay", "Review bounded replay files for synthetic or operator-redacted observations.", "Run connector example tests before any real replay.", "Does not retain raw private CSV payloads."},
		{"vendor_payload_adapter", "Vendor-shaped adapter", "Review vendor-shaped payload transform through cmd/avl-vendor-adapter or an equivalent sidecar.", "Dry-run the adapter with redacted fixtures before activation.", "Does not claim named vendor compatibility."},
	}
}

func vehicleAVLFields() []vehicleAVLField {
	return []vehicleAVLField{
		{"agency_id", "Agency ID", true, "agency_id"},
		{"device_id", "Device ID", true, "device.id"},
		{"vehicle_id", "Vehicle ID", true, "vehicle.id"},
		{"observed_timestamp", "Observed timestamp", true, "observed_at"},
		{"lat", "Latitude", true, "position.lat"},
		{"lon", "Longitude", true, "position.lon"},
		{"quality", "Quality", true, "quality"},
		{"route_hint", "Route hint", false, "route_id"},
		{"trip_hint", "Trip hint", false, "trip_id"},
		{"speed", "Speed", false, "speed"},
		{"bearing", "Bearing", false, "bearing"},
		{"accuracy", "Accuracy", false, "accuracy_meters"},
	}
}

func vehicleAVLInstanceInput(r *http.Request, principal auth.Principal) (connectorpkg.UpsertInstanceInput, error) {
	shape := strings.TrimSpace(r.FormValue("source_shape"))
	if !vehicleAVLShapeAllowed(shape) {
		return connectorpkg.UpsertInstanceInput{}, fmt.Errorf("unsupported vehicle data source shape")
	}
	displayName := strings.TrimSpace(r.FormValue("display_name"))
	owner := strings.TrimSpace(r.FormValue("owner"))
	secretRef := strings.TrimSpace(r.FormValue("secret_ref"))
	if err := connectorpkg.ValidateSecretRef(secretRef); err != nil {
		return connectorpkg.UpsertInstanceInput{}, err
	}
	fieldMap := make(map[string]string)
	for _, field := range vehicleAVLFields() {
		value := strings.TrimSpace(r.FormValue("field_" + field.ID))
		if field.Required && value == "" {
			return connectorpkg.UpsertInstanceInput{}, fmt.Errorf("%s mapping is required", field.Label)
		}
		if value == "" {
			continue
		}
		if err := validateFieldMappingValue(value); err != nil {
			return connectorpkg.UpsertInstanceInput{}, fmt.Errorf("%s mapping is invalid: %w", field.Label, err)
		}
		fieldMap[field.ID] = value
	}
	metadata := vehicleAVLConfigMetadata{
		SourceShape: shape,
		FieldMap:    fieldMap,
		Safety: map[string]any{
			"browser_exec_enabled": false,
			"dry_run_required":     true,
			"network_send_enabled": false,
			"target_path":          "/v1/telemetry",
		},
	}
	configJSON, err := json.Marshal(metadata)
	if err != nil {
		return connectorpkg.UpsertInstanceInput{}, err
	}
	var secretRefs []string
	if secretRef != "" {
		secretRefs = append(secretRefs, secretRef)
	}
	return connectorpkg.UpsertInstanceInput{
		AgencyID:      principal.AgencyID,
		ConnectorType: connectorpkg.TypeTelemetrySource,
		ConnectorKind: shape,
		DisplayName:   displayName,
		State:         connectorpkg.StateConfiguredNotTested,
		Owner:         owner,
		ConfigJSON:    configJSON,
		SecretRefs:    secretRefs,
		DryRunStatus:  "not_run",
		ActorID:       principal.Subject,
		Now:           time.Now().UTC(),
	}, nil
}

func vehicleAVLDryRunInput(r *http.Request, principal auth.Principal) (connectorpkg.CreateDryRunJobInput, error) {
	instanceID, err := strconv.ParseInt(strings.TrimSpace(r.FormValue("connector_instance_id")), 10, 64)
	if err != nil || instanceID <= 0 {
		return connectorpkg.CreateDryRunJobInput{}, fmt.Errorf("connector instance is required")
	}
	summary := strings.TrimSpace(r.FormValue("redacted_summary"))
	if summary == "" {
		return connectorpkg.CreateDryRunJobInput{}, fmt.Errorf("redacted summary is required")
	}
	if err := validateRedactedDryRunSummary(summary); err != nil {
		return connectorpkg.CreateDryRunJobInput{}, err
	}
	summaryJSON, err := json.Marshal(map[string]any{
		"summary":              summary,
		"raw_payload_retained": false,
		"browser_executed":     false,
	})
	if err != nil {
		return connectorpkg.CreateDryRunJobInput{}, err
	}
	accepted, err := dryRunCount(r, "accepted_count")
	if err != nil {
		return connectorpkg.CreateDryRunJobInput{}, err
	}
	rejected, err := dryRunCount(r, "rejected_count")
	if err != nil {
		return connectorpkg.CreateDryRunJobInput{}, err
	}
	dropped, err := dryRunCount(r, "dropped_count")
	if err != nil {
		return connectorpkg.CreateDryRunJobInput{}, err
	}
	return connectorpkg.CreateDryRunJobInput{
		AgencyID:            principal.AgencyID,
		ConnectorInstanceID: instanceID,
		Status:              strings.TrimSpace(r.FormValue("dry_run_status")),
		RedactedSummary:     summaryJSON,
		AcceptedCount:       accepted,
		RejectedCount:       rejected,
		DroppedCount:        dropped,
		RedactionScanStatus: strings.TrimSpace(r.FormValue("redaction_scan_status")),
		ActorID:             principal.Subject,
		Now:                 time.Now().UTC(),
	}, nil
}

func dryRunCount(r *http.Request, name string) (int, error) {
	raw := strings.TrimSpace(r.FormValue(name))
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 || value > 100000 {
		return 0, fmt.Errorf("%s must be a non-negative bounded count", strings.ReplaceAll(name, "_", " "))
	}
	return value, nil
}

func validateRedactedDryRunSummary(value string) error {
	if len(value) > 300 {
		return fmt.Errorf("redacted summary is too long")
	}
	lower := strings.ToLower(value)
	if strings.Contains(value, "://") || strings.Contains(lower, "password=") || strings.Contains(lower, "token=") || strings.Contains(lower, "secret=") {
		return fmt.Errorf("redacted summary must not contain endpoints or inline secrets")
	}
	return nil
}

func vehicleAVLShapeAllowed(shape string) bool {
	for _, candidate := range vehicleAVLShapes() {
		if shape == candidate.ID {
			return true
		}
	}
	return false
}

var fieldMappingPattern = regexp.MustCompile(`^[A-Za-z0-9_.\[\]-]{1,120}$`)

func validateFieldMappingValue(value string) error {
	if !fieldMappingPattern.MatchString(value) {
		return fmt.Errorf("use a bounded field path such as vehicle.id or position.lat")
	}
	lower := strings.ToLower(value)
	if strings.Contains(lower, "password") || strings.Contains(lower, "secret") || strings.Contains(lower, "token") || strings.Contains(value, "://") {
		return fmt.Errorf("field mappings must not contain secrets or endpoints")
	}
	return nil
}
