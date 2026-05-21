package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	connectorpkg "open-transit-rt/internal/connectors"
)

const vehicleAVLPostMaxBytes = 64 << 10

type vehicleAVLSetupView struct {
	GeneratedAt    time.Time              `json:"generated_at"`
	AgencyID       string                 `json:"agency_id"`
	Boundary       string                 `json:"boundary"`
	Notice         string                 `json:"notice,omitempty"`
	Error          string                 `json:"error,omitempty"`
	SourceShapes   []vehicleAVLShape      `json:"source_shapes"`
	FieldMappings  []vehicleAVLField      `json:"field_mappings"`
	Instances      []connectorInstanceRow `json:"instances"`
	NextAction     string                 `json:"next_action"`
	DryRunBoundary string                 `json:"dry_run_boundary"`
	ActivationGate string                 `json:"activation_gate"`
	DoesNotProve   string                 `json:"does_not_prove"`
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
	if strings.TrimSpace(r.FormValue("action")) != "save_vehicle_avl_connector" {
		http.Error(w, "unknown connector action", http.StatusBadRequest)
		return
	}
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

func buildVehicleAVLSetup(page operationsPage, notice string, errText string) vehicleAVLSetupView {
	registry := connectorRegistryForSection("connectors")
	var instances []connectorInstanceRow
	for _, row := range connectorInstanceRows(page, registry) {
		if row.ConnectorType == connectorpkg.TypeTelemetrySource {
			instances = append(instances, row)
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
		NextAction:     "Choose the closest source shape, map the required fields, save metadata, then run a server-owned dry-run before activation review.",
		DryRunBoundary: "Dry-run is required and remains server-owned. This browser page records configuration metadata only.",
		ActivationGate: "Activation stays blocked until mapping, dry-run, device bindings, token refs, stale/future/quality rules, and redaction checks pass.",
		DoesNotProve:   "This setup page does not prove vendor compatibility, hardware certification, production AVL reliability, consumer acceptance, compliance, uptime, or ETA quality.",
	}
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
