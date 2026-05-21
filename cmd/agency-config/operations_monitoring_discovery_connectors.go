package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	connectorpkg "open-transit-rt/internal/connectors"
)

const monitoringDiscoveryConnectorPostMaxBytes = 32 << 10

type monitoringConnectorView struct {
	GeneratedAt       time.Time              `json:"generated_at"`
	AgencyID          string                 `json:"agency_id"`
	Boundary          string                 `json:"boundary"`
	Notice            string                 `json:"notice,omitempty"`
	Error             string                 `json:"error,omitempty"`
	Instances         []connectorInstanceRow `json:"instances"`
	Configured        []connectorInstanceRow `json:"configured_instances"`
	DigestPreview     []monitoringDigestRow  `json:"digest_preview"`
	Destinations      []monitoringModeRow    `json:"destinations"`
	NextAction        string                 `json:"next_action"`
	NoSendDefault     string                 `json:"no_send_default"`
	RedactionBoundary string                 `json:"redaction_boundary"`
	DoesNotProve      string                 `json:"does_not_prove"`
}

type discoveryConnectorView struct {
	GeneratedAt  time.Time               `json:"generated_at"`
	AgencyID     string                  `json:"agency_id"`
	Boundary     string                  `json:"boundary"`
	Notice       string                  `json:"notice,omitempty"`
	Error        string                  `json:"error,omitempty"`
	Instances    []connectorInstanceRow  `json:"instances"`
	Configured   []connectorInstanceRow  `json:"configured_instances"`
	Readiness    []discoveryReadinessRow `json:"readiness"`
	NextAction   string                  `json:"next_action"`
	NoAutomation string                  `json:"no_automation"`
	DoesNotProve string                  `json:"does_not_prove"`
}

type monitoringDigestRow struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
}

type monitoringModeRow struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	Reference    string `json:"reference"`
	SendState    string `json:"send_state"`
	FirstCheck   string `json:"first_check"`
	DoesNotProve string `json:"does_not_prove"`
}

type discoveryReadinessRow struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	DoesNotProve  string `json:"does_not_prove"`
}

type monitoringConnectorConfigMetadata struct {
	Mode                    string   `json:"mode"`
	DestinationRef          string   `json:"destination_ref"`
	DigestSources           []string `json:"digest_sources"`
	NoSendDefault           bool     `json:"no_send_default"`
	RedactedTestResultOnly  bool     `json:"redacted_test_result_only"`
	NotificationDeliveryOff bool     `json:"notification_delivery_off"`
}

type discoveryConnectorConfigMetadata struct {
	Mode                      string `json:"mode"`
	PublicBaseURLEnvRef       string `json:"public_base_url_env_ref"`
	LicenseContactOwnerRef    string `json:"license_contact_owner_ref"`
	FeedsJSONReadiness        string `json:"feeds_json_readiness"`
	PortalAutomationEnabled   bool   `json:"portal_automation_enabled"`
	ConsumerStatusMutationOff bool   `json:"consumer_status_mutation_off"`
}

func (h *handler) renderMonitoringConnector(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "monitoring-setup")
	page.MonitoringConnector = buildMonitoringConnector(page, "", "")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "monitoring-setup", page)
}

func (h *handler) renderMonitoringConnectorJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "monitoring-setup")
	page.MonitoringConnector = buildMonitoringConnector(page, "", "")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.MonitoringConnector)
}

func (h *handler) renderDiscoveryConnector(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "discovery-setup")
	page.DiscoveryConnector = buildDiscoveryConnector(page, "", "")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "discovery-setup", page)
}

func (h *handler) renderDiscoveryConnectorJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "discovery-setup")
	page.DiscoveryConnector = buildDiscoveryConnector(page, "", "")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.DiscoveryConnector)
}

func (h *handler) operationsMonitoringConnectorPost(w http.ResponseWriter, r *http.Request) {
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
	if strings.TrimSpace(r.FormValue("action")) != "save_monitoring_connector" {
		http.Error(w, "unknown connector action", http.StatusBadRequest)
		return
	}
	page := h.buildOperationsPage(r, principal, "monitoring-setup")
	writer, ok := h.connectorInstances.(connectorInstanceWriter)
	if !ok || writer == nil {
		page.MonitoringConnector = buildMonitoringConnector(page, "", "connector instance writer is not available in this runtime")
		renderOperationsTemplate(w, "monitoring-setup", page)
		return
	}
	input, err := monitoringConnectorInstanceInput(r, principal)
	if err != nil {
		page.MonitoringConnector = buildMonitoringConnector(page, "", err.Error())
		renderOperationsTemplate(w, "monitoring-setup", page)
		return
	}
	if _, err := writer.UpsertInstance(r.Context(), input); err != nil {
		page.MonitoringConnector = buildMonitoringConnector(page, "", "monitoring connector metadata could not be stored")
		renderOperationsTemplate(w, "monitoring-setup", page)
		return
	}
	page = h.buildOperationsPage(r, principal, "monitoring-setup")
	page.MonitoringConnector = buildMonitoringConnector(page, "Monitoring connector metadata was saved with no-send defaults. No notification or export delivery was attempted.", "")
	renderOperationsTemplate(w, "monitoring-setup", page)
}

func (h *handler) operationsDiscoveryConnectorPost(w http.ResponseWriter, r *http.Request) {
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
	if strings.TrimSpace(r.FormValue("action")) != "save_discovery_connector" {
		http.Error(w, "unknown connector action", http.StatusBadRequest)
		return
	}
	page := h.buildOperationsPage(r, principal, "discovery-setup")
	writer, ok := h.connectorInstances.(connectorInstanceWriter)
	if !ok || writer == nil {
		page.DiscoveryConnector = buildDiscoveryConnector(page, "", "connector instance writer is not available in this runtime")
		renderOperationsTemplate(w, "discovery-setup", page)
		return
	}
	input, err := discoveryConnectorInstanceInput(r, principal, page)
	if err != nil {
		page.DiscoveryConnector = buildDiscoveryConnector(page, "", err.Error())
		renderOperationsTemplate(w, "discovery-setup", page)
		return
	}
	if _, err := writer.UpsertInstance(r.Context(), input); err != nil {
		page.DiscoveryConnector = buildDiscoveryConnector(page, "", "discovery connector metadata could not be stored")
		renderOperationsTemplate(w, "discovery-setup", page)
		return
	}
	page = h.buildOperationsPage(r, principal, "discovery-setup")
	page.DiscoveryConnector = buildDiscoveryConnector(page, "Discovery connector metadata was saved for readiness review only. No portal automation or consumer status mutation was attempted.", "")
	renderOperationsTemplate(w, "discovery-setup", page)
}

func buildMonitoringConnector(page operationsPage, notice string, errText string) monitoringConnectorView {
	instances, configured := connectorRowsForType(page, connectorpkg.TypeMonitoringExport)
	return monitoringConnectorView{
		GeneratedAt:       page.GeneratedAt,
		AgencyID:          page.AgencyID,
		Boundary:          "Private Monitoring Setup stores no-send destination reference labels and redacted digest metadata only. It does not send notifications, create exports, contact webhooks, store endpoint URLs, or create evidence.",
		Notice:            notice,
		Error:             errText,
		Instances:         instances,
		Configured:        configured,
		DigestPreview:     monitoringDigestRows(page),
		Destinations:      monitoringModeRows(),
		NextAction:        "Save a no-send destination reference, review the health digest preview, then let a deployment owner configure any real delivery outside the browser.",
		NoSendDefault:     "All monitoring/export connectors are no-send by default; browser review never delivers notifications or writes external exports.",
		RedactionBoundary: "Digest rows use existing private status summaries and omit raw validation reports, private URLs, tokens, payloads, and local paths.",
		DoesNotProve:      "Monitoring connector configuration does not prove notification delivery, SLA coverage, uptime guarantees, hosted service availability, compliance, consumer acceptance, or production readiness.",
	}
}

func buildDiscoveryConnector(page operationsPage, notice string, errText string) discoveryConnectorView {
	instances, configured := connectorRowsForType(page, connectorpkg.TypeConsumerDiscovery)
	return discoveryConnectorView{
		GeneratedAt:  page.GeneratedAt,
		AgencyID:     page.AgencyID,
		Boundary:     "Private Discovery Setup reviews /public/feeds.json metadata readiness and stores deployment reference labels only. It does not automate portal submissions, contact consumers, change consumer statuses, or create evidence.",
		Notice:       notice,
		Error:        errText,
		Instances:    instances,
		Configured:   configured,
		Readiness:    discoveryReadinessRows(page),
		NextAction:   "Confirm public feed metadata, license, and contact fields, then use sharing prep without submission automation.",
		NoAutomation: "Portal automation and consumer status mutation stay disabled. Prepared metadata is not submission, review, listing, display, acceptance, or compliance proof.",
		DoesNotProve: "Discovery connector configuration does not prove consumer submission, ingestion, listing, display, acceptance, agency approval, public launch, compliance, hosted service, SLA, or production readiness.",
	}
}

func connectorRowsForType(page operationsPage, connectorType string) ([]connectorInstanceRow, []connectorInstanceRow) {
	registry := connectorRegistryForSection("connectors")
	var instances []connectorInstanceRow
	var configured []connectorInstanceRow
	for _, row := range connectorInstanceRows(page, registry) {
		if row.ConnectorType != connectorType {
			continue
		}
		instances = append(instances, row)
		if row.State != string(connectorpkg.StateExampleAvailable) && row.State != string(connectorpkg.StateNotConfigured) {
			configured = append(configured, row)
		}
	}
	return instances, configured
}

func monitoringModeRows() []monitoringModeRow {
	return []monitoringModeRow{
		{ID: "health_digest_no_send", Label: "Health digest preview", Reference: "MONITORING_DIGEST_DESTINATION_REF", SendState: "no_send", FirstCheck: "make operations-reliability", DoesNotProve: "Notification delivery, SLA, uptime, or production readiness."},
		{ID: "redacted_export_no_send", Label: "Redacted export preview", Reference: "MONITORING_EXPORT_DESTINATION_REF", SendState: "no_send", FirstCheck: "make operations-notify", DoesNotProve: "External delivery, retained evidence, compliance, or hosted-service monitoring."},
	}
}

func monitoringDigestRows(page operationsPage) []monitoringDigestRow {
	return []monitoringDigestRow{
		{ID: "validator_health", Label: "Validator health", Status: firstNonEmpty(page.ValidationHealth.OverallStatus, checklistStatusUnknown), CurrentSignal: "tooling=" + firstNonEmpty(page.ValidationHealth.ToolingStatus, checklistStatusUnknown), NextAction: "Open Validation Health before including validator state in any digest."},
		{ID: "feed_health", Label: "Feed health", Status: feedHealthDigestStatus(page), CurrentSignal: fmt.Sprintf("%d feed row(s)", len(page.FeedHealth.Rows)), NextAction: "Review feed health rows before sharing a digest preview."},
		{ID: "reliability", Label: "Reliability", Status: firstNonEmpty(page.Reliability.OverallStatus, checklistStatusUnknown), CurrentSignal: fmt.Sprintf("%d reliability feed row(s)", len(page.Reliability.Feeds)), NextAction: "Review reliability diagnostics before configuring external monitoring."},
	}
}

func feedHealthDigestStatus(page operationsPage) string {
	if len(page.FeedHealth.Rows) == 0 {
		return checklistStatusUnknown
	}
	status := checklistStatusOK
	for _, row := range page.FeedHealth.Rows {
		normalized := readinessV2NormalizeStatus(row.Status)
		if normalizeChecklistStatus(normalized) != checklistStatusOK {
			status = normalized
			if normalized == checklistStatusBlocked || normalized == checklistStatusMissing {
				return normalized
			}
		}
	}
	return status
}

func discoveryReadinessRows(page operationsPage) []discoveryReadinessRow {
	feedsCount := len(page.Discovery.Feeds)
	return []discoveryReadinessRow{
		discoveryReadiness("feeds_json", "/public/feeds.json metadata", feedsCount > 0, fmt.Sprintf("%d configured feed metadata row(s)", feedsCount), "Configure all public feed metadata before sharing discovery links."),
		discoveryReadiness("public_base_url", "Public base URL", strings.TrimSpace(page.Discovery.PublicBaseURL) != "", yesNo(strings.TrimSpace(page.Discovery.PublicBaseURL) != ""), "Configure a deployment-owned public root before discovery review."),
		discoveryReadiness("license", "License metadata", strings.TrimSpace(page.Discovery.License.Name) != "", firstNonEmpty(page.Discovery.License.Name, "missing"), "Add open license metadata before sharing preparation."),
		discoveryReadiness("technical_contact", "Technical contact", strings.TrimSpace(page.Discovery.TechnicalContactEmail) != "", yesNo(strings.TrimSpace(page.Discovery.TechnicalContactEmail) != ""), "Add a technical contact before sharing preparation."),
		discoveryReadiness("consumer_status", "Consumer status mutation", true, "disabled", "Keep portal automation and consumer status mutation disabled unless separately authorized."),
	}
}

func discoveryReadiness(id string, label string, ok bool, signal string, next string) discoveryReadinessRow {
	status := checklistStatusBlocked
	if ok {
		status = checklistStatusOK
	}
	return discoveryReadinessRow{
		ID:            id,
		Label:         label,
		Status:        status,
		CurrentSignal: firstNonEmpty(signal, "not configured"),
		NextAction:    next,
		DoesNotProve:  "Discovery readiness does not prove submission, listing, display, consumer acceptance, public launch, compliance, or production readiness.",
	}
}

func monitoringConnectorInstanceInput(r *http.Request, principal auth.Principal) (connectorpkg.UpsertInstanceInput, error) {
	mode := strings.TrimSpace(r.FormValue("mode"))
	if mode != "health_digest_no_send" && mode != "redacted_export_no_send" {
		return connectorpkg.UpsertInstanceInput{}, fmt.Errorf("unsupported monitoring connector mode")
	}
	destinationRef, err := deploymentRefLabel(r.FormValue("destination_ref"), "destination ref")
	if err != nil {
		return connectorpkg.UpsertInstanceInput{}, err
	}
	displayName := firstNonEmpty(strings.TrimSpace(r.FormValue("display_name")), "Monitoring export preview")
	metadata := monitoringConnectorConfigMetadata{
		Mode:                    mode,
		DestinationRef:          destinationRef,
		DigestSources:           []string{"validator_health", "feed_health", "reliability"},
		NoSendDefault:           true,
		RedactedTestResultOnly:  true,
		NotificationDeliveryOff: true,
	}
	configRaw, err := json.Marshal(metadata)
	if err != nil {
		return connectorpkg.UpsertInstanceInput{}, fmt.Errorf("encode monitoring connector metadata: %w", err)
	}
	return connectorpkg.UpsertInstanceInput{
		AgencyID:      principal.AgencyID,
		ConnectorType: connectorpkg.TypeMonitoringExport,
		ConnectorKind: mode,
		DisplayName:   displayName,
		State:         connectorpkg.StateConfiguredNotTested,
		Owner:         strings.TrimSpace(r.FormValue("owner")),
		ConfigJSON:    json.RawMessage(configRaw),
		DryRunStatus:  "not_run",
		ActorID:       principal.Subject,
		Now:           time.Now().UTC(),
	}, nil
}

func discoveryConnectorInstanceInput(r *http.Request, principal auth.Principal, page operationsPage) (connectorpkg.UpsertInstanceInput, error) {
	publicBaseRef, err := deploymentRefLabel(r.FormValue("public_base_url_env_ref"), "public base URL env ref")
	if err != nil {
		return connectorpkg.UpsertInstanceInput{}, err
	}
	ownerRef, err := deploymentRefLabel(r.FormValue("license_contact_owner_ref"), "license/contact owner ref")
	if err != nil {
		return connectorpkg.UpsertInstanceInput{}, err
	}
	displayName := firstNonEmpty(strings.TrimSpace(r.FormValue("display_name")), "Feed discovery readiness")
	metadata := discoveryConnectorConfigMetadata{
		Mode:                      "feeds_json_readiness",
		PublicBaseURLEnvRef:       publicBaseRef,
		LicenseContactOwnerRef:    ownerRef,
		FeedsJSONReadiness:        discoveryFeedsJSONReadiness(page),
		PortalAutomationEnabled:   false,
		ConsumerStatusMutationOff: true,
	}
	configRaw, err := json.Marshal(metadata)
	if err != nil {
		return connectorpkg.UpsertInstanceInput{}, fmt.Errorf("encode discovery connector metadata: %w", err)
	}
	return connectorpkg.UpsertInstanceInput{
		AgencyID:      principal.AgencyID,
		ConnectorType: connectorpkg.TypeConsumerDiscovery,
		ConnectorKind: "feeds_json_readiness",
		DisplayName:   displayName,
		State:         connectorpkg.StateConfiguredNotTested,
		Owner:         strings.TrimSpace(r.FormValue("owner")),
		ConfigJSON:    json.RawMessage(configRaw),
		DryRunStatus:  "not_run",
		ActorID:       principal.Subject,
		Now:           time.Now().UTC(),
	}, nil
}

func discoveryFeedsJSONReadiness(page operationsPage) string {
	if len(page.Discovery.Feeds) == 0 || strings.TrimSpace(page.Discovery.PublicBaseURL) == "" || strings.TrimSpace(page.Discovery.License.Name) == "" || strings.TrimSpace(page.Discovery.TechnicalContactEmail) == "" {
		return checklistStatusNeedsReview
	}
	return checklistStatusOK
}

func deploymentRefLabel(value string, label string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("%s is required", label)
	}
	lower := strings.ToLower(value)
	if strings.Contains(lower, "://") ||
		strings.ContainsAny(value, "\x00\r\n\t /\\") ||
		strings.Contains(lower, "password=") ||
		strings.Contains(lower, "token=") ||
		strings.Contains(lower, "secret=") {
		return "", fmt.Errorf("%s must be an uppercase deployment reference label, not a URL, path, or inline secret", label)
	}
	if err := connectorpkg.ValidateSecretRef(value); err != nil {
		return "", fmt.Errorf("%s must be an uppercase deployment reference label", label)
	}
	return value, nil
}
