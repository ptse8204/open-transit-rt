package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
	"open-transit-rt/internal/devices"
	"open-transit-rt/internal/prediction"
	"open-transit-rt/internal/state"
)

const (
	defaultTelemetryLimit        = 200
	defaultStaleSeconds          = 90
	gtfsQualityPostMaxBytes      = 64 << 10
	validationHealthPostMaxBytes = 64 << 10
)

type operationsPage struct {
	Title                  string
	AgencyID               string
	GeneratedAt            time.Time
	EnvironmentLabel       string
	CSRFToken              string
	Discovery              compliance.FeedDiscovery
	DiscoveryError         string
	PublicationConfig      compliance.PublicationConfig
	PublicationError       string
	SetupNotice            string
	SetupError             string
	SetupSteps             []setupStepView
	ValidationResult       *compliance.ValidationResult
	ActiveFeedVersion      string
	FeedsUpdatedAt         *time.Time
	TelemetryUpdatedAt     *time.Time
	TripUpdatesQuality     tripUpdatesQualityView
	ScorecardUpdatedAt     *time.Time
	ConsumersUpdatedAt     *time.Time
	EvidenceUpdatedAt      string
	Scorecard              *compliance.Scorecard
	ScorecardError         string
	Consumers              []consumerStatusView
	RuntimeConsumers       []consumerStatusView
	ReadinessItems         []readinessItemView
	ReadinessV2            operationsReadinessV2View
	Checklist              operatorChecklistView
	Cockpit                operationsCockpitView
	FirstRun               operationsFirstRunView
	Launchpad              agencyLaunchpadView
	SetupWizard            operationsSetupWizardView
	ConnectorHub           connectorHubView
	ConnectorTests         connectorTestsView
	Help                   operationsHelpView
	ContextHelp            operationsContextHelp
	FeedHealth             operationsFeedHealthView
	Maintenance            operationsMaintenanceView
	TelemetrySimulator     operationsTelemetrySimulatorView
	GTFSImportResult       *gtfsImportResultView
	GTFSImportSource       gtfsImportSourceReview
	GTFSImportNotice       string
	GTFSImportError        string
	GTFSQuality            compliance.GTFSQualityTriage
	GTFSQualityGuidance    operationsGTFSQualityGuidanceView
	GTFSQualityNotice      string
	GTFSQualityError       string
	ValidationHealth       compliance.ValidationHealthSummary
	ValidationHealthNotice string
	ValidationHealthError  string
	Reliability            compliance.ReliabilitySummary
	ReliabilityError       string
	IsAdmin                bool
	ConsumerError          string
	Telemetry              []telemetryView
	TelemetryError         string
	StaleCount             int
	Devices                []devices.Binding
	DeviceRows             []operationsDeviceRow
	DeviceOnboarding       []operationsDeviceOnboardingUseCase
	DeviceError            string
	DeviceToken            string
	DeviceTokenMeta        devices.RebindResult
	Links                  []evidenceLink
	NavGroups              []operationsNavGroup
	Section                string
	StaleThreshold         time.Duration
}

type consumerStatusView struct {
	Name        string
	Status      string
	UpdatedAt   *time.Time
	Source      string
	Notes       string
	PacketPath  string
	CurrentPath string
}

type setupStepView struct {
	Name       string
	Status     string
	Source     string
	Evidence   string
	NextAction string
	ActionURL  string
}

type readinessItemView struct {
	Name          string
	Status        string
	Source        string
	Evidence      string
	NextAction    string
	ClaimBoundary string
}

type telemetryView struct {
	VehicleID        string
	DeviceID         string
	ObservedAt       time.Time
	ReceivedAt       time.Time
	AgeSeconds       int64
	Stale            bool
	AssignmentState  string
	RouteID          string
	TripID           string
	Confidence       string
	ReasonCodes      []string
	DegradedState    string
	AssignmentSource string
	AssignmentAt     *time.Time
}

type tripUpdatesQualityView struct {
	Recorded                      bool
	Message                       string
	SnapshotAt                    *time.Time
	AdapterName                   string
	DiagnosticsStatus             string
	DiagnosticsReason             string
	ActiveFeedVersionID           string
	DiagnosticsPersistenceOutcome string
	UnknownAssignmentRate         string
	AmbiguousAssignmentRate       string
	StaleTelemetryRate            string
	TripUpdatesCoverageRate       string
	FutureStopCoverageRate        string
	EligiblePredictionCandidates  int
	TripUpdatesEmitted            int
	UnknownAssignments            int
	AmbiguousAssignments          int
	StaleTelemetryRows            int
	ManualOverrideAssignments     int
	CanceledTripsEmitted          int
	CancellationAlertLinksMissing int
	WithheldByReason              []countView
}

type countView struct {
	Label string
	Count int
}

type evidenceLink struct {
	Label     string
	Path      string
	UpdatedAt string
}

type tripUpdatesDiagnosticsReader interface {
	LatestTripUpdatesDiagnostics(ctx context.Context, agencyID string) (compliance.TripUpdatesDiagnosticsSummary, error)
}

type publicationConfigReader interface {
	PublicationConfig(ctx context.Context, agencyID string) (compliance.PublicationConfig, error)
}

type validationReportReader interface {
	LatestValidationReport(ctx context.Context, agencyID string, feedType string, validatorName string) (*compliance.ValidationReportRecord, error)
}

type reliabilityReader interface {
	LatestReliabilityFeedHealth(ctx context.Context, agencyID string, limit int) ([]compliance.ReliabilityFeedHealthRecord, error)
	ReliabilityIncidentRollup(ctx context.Context, agencyID string, now time.Time, recentLimit int) (compliance.ReliabilityIncidentRollup, error)
}

func (h *handler) operationsRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/admin/operations" {
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderOperations(w, r, "dashboard")
		return
	}
	if r.URL.Path == "/admin/operations.json" {
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderOperationsCockpitJSON(w, r)
		return
	}
	trimmed := strings.Trim(strings.TrimPrefix(r.URL.Path, "/admin/operations/"), "/")
	switch trimmed {
	case "launchpad":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderLaunchpad(w, r)
	case "launchpad.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderLaunchpadJSON(w, r)
	case "connectors":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderConnectorHub(w, r)
	case "connectors.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderConnectorHubJSON(w, r)
	case "connectors/tests":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderConnectorTests(w, r)
	case "connectors/tests.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderConnectorTestsJSON(w, r)
	case "help":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderOperationsHelp(w, r)
	case "help.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderOperationsHelpJSON(w, r)
	case "setup-wizard":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderSetupWizard(w, r)
	case "setup-wizard.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderSetupWizardJSON(w, r)
	case "feed-health":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderFeedHealth(w, r)
	case "feed-health.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderFeedHealthJSON(w, r)
	case "readiness":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderReadinessV2(w, r)
	case "readiness.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderReadinessV2JSON(w, r)
	case "telemetry-simulator":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderTelemetrySimulator(w, r)
	case "telemetry-simulator.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderTelemetrySimulatorJSON(w, r)
	case "gtfs-import":
		w.Header().Set("Cache-Control", "no-store")
		switch r.Method {
		case http.MethodGet:
			h.renderGTFSImport(w, r)
		case http.MethodPost:
			h.operationsGTFSImportPost(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	case "checklist":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderOperationsChecklist(w, r)
	case "checklist.json":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderOperationsChecklistJSON(w, r)
	case "gtfs-quality":
		w.Header().Set("Cache-Control", "no-store")
		switch r.Method {
		case http.MethodGet:
			h.renderGTFSQuality(w, r)
		case http.MethodPost:
			h.operationsGTFSQualityPost(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	case "validation-health":
		w.Header().Set("Cache-Control", "no-store")
		switch r.Method {
		case http.MethodGet:
			h.renderValidationHealth(w, r)
		case http.MethodPost:
			h.operationsValidationHealthPost(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	case "validation-health.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderValidationHealthJSON(w, r)
	case "reliability":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderReliability(w, r)
	case "reliability.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderReliabilityJSON(w, r)
	case "maintenance":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderMaintenance(w, r)
	case "maintenance.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderMaintenanceJSON(w, r)
	case "feeds", "telemetry", "devices", "consumers", "evidence", "setup":
		if trimmed == "devices" && r.Method == http.MethodPost {
			h.operationsDeviceRebind(w, r)
			return
		}
		if trimmed == "setup" && r.Method == http.MethodPost {
			h.operationsSetupPost(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderOperations(w, r, trimmed)
	default:
		http.NotFound(w, r)
	}
}

func (h *handler) renderOperations(w http.ResponseWriter, r *http.Request, section string) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, section)
	renderOperationsTemplate(w, section, page)
}

func (h *handler) renderOperationsCockpitJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "dashboard")
	writeJSON(w, http.StatusOK, page.Cockpit)
}

func (h *handler) renderGTFSQuality(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "gtfs-quality")
	renderOperationsTemplate(w, "gtfs-quality", page)
}

func (h *handler) renderValidationHealth(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "validation-health")
	renderOperationsTemplate(w, "validation-health", page)
}

func (h *handler) renderValidationHealthJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "validation-health")
	writeJSON(w, http.StatusOK, page.ValidationHealth)
}

func (h *handler) renderReliability(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "reliability")
	renderOperationsTemplate(w, "reliability", page)
}

func (h *handler) renderReliabilityJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "reliability")
	writeJSON(w, http.StatusOK, page.Reliability)
}

func (h *handler) renderMaintenance(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "maintenance")
	renderOperationsTemplate(w, "maintenance", page)
}

func (h *handler) renderMaintenanceJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "maintenance")
	writeJSON(w, http.StatusOK, page.Maintenance)
}

func (h *handler) operationsValidationHealthPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, validationHealthPostMaxBytes)
	principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
	if !ok {
		return
	}
	if !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	if err := r.ParseForm(); err != nil {
		page := h.buildOperationsPage(r, principal, "validation-health")
		page.ValidationHealthError = "Validator health request was blocked because the form body is invalid or too large."
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		renderOperationsTemplate(w, "validation-health", page)
		return
	}
	if principal.Method == auth.MethodCookie && strings.TrimSpace(h.csrfSecret) != "" && strings.TrimSpace(r.FormValue("csrf_token")) != csrfToken(h.csrfSecret, principal) {
		page := h.buildOperationsPage(r, principal, "validation-health")
		page.ValidationHealthError = "Validator health request was blocked because the CSRF token is invalid."
		w.WriteHeader(http.StatusForbidden)
		renderOperationsTemplate(w, "validation-health", page)
		return
	}
	if err := rejectValidationHealthUnexpectedFields(r); err != nil {
		page := h.buildOperationsPage(r, principal, "validation-health")
		page.ValidationHealthError = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		renderOperationsTemplate(w, "validation-health", page)
		return
	}
	if strings.TrimSpace(r.FormValue("action")) != "run_all" {
		page := h.buildOperationsPage(r, principal, "validation-health")
		page.ValidationHealthError = "unknown validator health action"
		w.WriteHeader(http.StatusBadRequest)
		renderOperationsTemplate(w, "validation-health", page)
		return
	}
	page := h.buildOperationsPage(r, principal, "validation-health")
	records, artifactOverrides := h.runValidationHealthAll(r, principal, page.Discovery)
	page = h.buildOperationsPage(r, principal, "validation-health")
	if len(records) > 0 || len(artifactOverrides) > 0 {
		page.ValidationHealth = h.validationHealthSummary(r, principal.AgencyID, page.Discovery, records, artifactOverrides)
	}
	page.ValidationHealthNotice = "Validator health run finished. Results were stored only as normal validation_report rows where validators ran."
	renderOperationsTemplate(w, "validation-health", page)
}

func rejectValidationHealthUnexpectedFields(r *http.Request) error {
	allowed := map[string]bool{"action": true, "csrf_token": true}
	blocked := []string{"feed_type", "validator_id", "validator_command", "validator_path", "output_path", "artifact_path", "report_path", "schedule_zip_path", "realtime_pb_path", "url", "URL", "argv", "args", "timeout"}
	for name := range r.Form {
		if allowed[name] {
			continue
		}
		lower := strings.ToLower(name)
		for _, blockedName := range blocked {
			if lower == strings.ToLower(blockedName) || strings.Contains(lower, "url") || strings.Contains(lower, "path") || strings.Contains(lower, "timeout") || strings.Contains(lower, "argv") || strings.Contains(lower, "args") {
				return fmt.Errorf("validator health accepts only action=run_all and CSRF fields")
			}
		}
		return fmt.Errorf("validator health accepts only action=run_all and CSRF fields")
	}
	return nil
}

func (h *handler) operationsGTFSQualityPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, gtfsQualityPostMaxBytes)
	principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
	if !ok {
		return
	}
	if !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	if err := r.ParseForm(); err != nil {
		page := h.buildOperationsPage(r, principal, "gtfs-quality")
		page.GTFSQualityError = "GTFS quality rerun request was blocked because the form body is invalid or too large."
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		renderOperationsTemplate(w, "gtfs-quality", page)
		return
	}
	if strings.TrimSpace(h.csrfSecret) != "" && strings.TrimSpace(r.FormValue("csrf_token")) != csrfToken(h.csrfSecret, principal) {
		page := h.buildOperationsPage(r, principal, "gtfs-quality")
		page.GTFSQualityError = "GTFS quality rerun request was blocked because the CSRF token is invalid."
		w.WriteHeader(http.StatusForbidden)
		renderOperationsTemplate(w, "gtfs-quality", page)
		return
	}
	if err := rejectGTFSQualityUnexpectedFields(r); err != nil {
		page := h.buildOperationsPage(r, principal, "gtfs-quality")
		page.GTFSQualityError = err.Error()
		renderOperationsTemplate(w, "gtfs-quality", page)
		return
	}
	if strings.TrimSpace(r.FormValue("action")) != "rerun_static_validator" {
		page := h.buildOperationsPage(r, principal, "gtfs-quality")
		page.GTFSQualityError = "unknown GTFS quality action"
		renderOperationsTemplate(w, "gtfs-quality", page)
		return
	}
	page := h.buildOperationsPage(r, principal, "gtfs-quality")
	activeFeedVersionID, activeRevision := scheduleFeedVersion(page.Discovery.Feeds)
	if activeFeedVersionID == "" {
		page.GTFSQualityError = "Blocking: no active published schedule feed version is available. Import or publish a schedule before rerunning the static MobilityData validator."
		renderOperationsTemplate(w, "gtfs-quality", page)
		return
	}
	result, err := h.runValidationForFeed(r, principal, "schedule", activeFeedVersionID)
	if err != nil {
		page.GTFSQualityError = "Blocking: static MobilityData validator could not run. Check validator tooling configuration and the active schedule, then retry."
		renderOperationsTemplate(w, "gtfs-quality", page)
		return
	}
	record := &compliance.ValidationReportRecord{Result: result, CreatedAt: time.Now().UTC().Truncate(time.Second)}
	page = h.buildOperationsPage(r, principal, "gtfs-quality")
	page.GTFSQuality.Canonical = compliance.BuildGTFSQualityTriage(compliance.GTFSQualityTriageInput{Canonical: record, ActiveFeedVersionID: activeFeedVersionID, ActiveFeedRevisionTime: activeRevision}).Canonical
	page.GTFSQualityNotice = "Static MobilityData validator rerun finished and was stored only as the normal validation result row."
	renderOperationsTemplate(w, "gtfs-quality", page)
}

func rejectGTFSQualityUnexpectedFields(r *http.Request) error {
	allowed := map[string]bool{"action": true, "csrf_token": true}
	blocked := []string{"validator_id", "validator_command", "validator_path", "output_path", "artifact_path", "report_path", "schedule_zip_path", "url", "URL", "argv", "args", "timeout"}
	for name := range r.Form {
		if allowed[name] {
			continue
		}
		lower := strings.ToLower(name)
		for _, blockedName := range blocked {
			if lower == strings.ToLower(blockedName) || strings.Contains(lower, "url") || strings.Contains(lower, "path") || strings.Contains(lower, "timeout") {
				return fmt.Errorf("GTFS quality rerun accepts only action and CSRF fields")
			}
		}
		return fmt.Errorf("GTFS quality rerun accepts only action and CSRF fields")
	}
	return nil
}

func (h *handler) operationsDeviceRebind(w http.ResponseWriter, r *http.Request) {
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
	result, err := h.devices.Rebind(r.Context(), devices.RebindInput{
		AgencyID:  principal.AgencyID,
		DeviceID:  r.FormValue("device_id"),
		VehicleID: r.FormValue("vehicle_id"),
		ActorID:   principal.Subject,
		Reason:    r.FormValue("reason"),
	})
	page := h.buildOperationsPage(r, principal, "devices")
	if err != nil {
		page.DeviceError = err.Error()
		renderOperationsTemplate(w, "devices", page)
		return
	}
	page.DeviceToken = result.Token
	result.Token = ""
	page.DeviceTokenMeta = result
	renderOperationsTemplate(w, "devices", page)
}

func (h *handler) operationsSetupPost(w http.ResponseWriter, r *http.Request) {
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
	page := h.buildOperationsPage(r, principal, "setup")
	switch strings.TrimSpace(r.FormValue("action")) {
	case "publication_bootstrap":
		input, err := setupPublicationInput(r, principal)
		if err != nil {
			page.SetupError = err.Error()
			renderOperationsTemplate(w, "setup", page)
			return
		}
		if err := h.store.BootstrapPublication(r.Context(), input); err != nil {
			page.SetupError = "publication metadata could not be stored"
			renderOperationsTemplate(w, "setup", page)
			return
		}
		page = h.buildOperationsPage(r, principal, "setup")
		page.SetupNotice = "Publication metadata was stored using the existing bootstrap/update workflow."
		renderOperationsTemplate(w, "setup", page)
	case "run_validation":
		if err := rejectSetupValidationUnsafeFields(r); err != nil {
			page.SetupError = err.Error()
			renderOperationsTemplate(w, "setup", page)
			return
		}
		feedType := strings.TrimSpace(r.FormValue("feed_type"))
		result, err := h.runValidationForFeed(r, principal, feedType, "")
		if err != nil {
			page.SetupError = err.Error()
			renderOperationsTemplate(w, "setup", page)
			return
		}
		page = h.buildOperationsPage(r, principal, "setup")
		page.ValidationResult = &result
		page.SetupNotice = "Validation finished and was stored as supporting evidence only."
		renderOperationsTemplate(w, "setup", page)
	default:
		http.Error(w, "unknown setup action", http.StatusBadRequest)
	}
}

func rejectSetupValidationUnsafeFields(r *http.Request) error {
	for _, name := range []string{
		"validator_id",
		"validator_command",
		"validator_path",
		"output_path",
		"artifact_path",
		"realtime_pb_path",
		"schedule_zip_path",
		"report_path",
		"argv",
		"args",
	} {
		if _, ok := r.Form[name]; ok {
			return fmt.Errorf("validation setup form only accepts feed type")
		}
	}
	return nil
}

func setupPublicationInput(r *http.Request, principal auth.Principal) (compliance.BootstrapInput, error) {
	publicBaseURL, err := setupFormValue(r, "public_base_url", 2048)
	if err != nil {
		return compliance.BootstrapInput{}, err
	}
	feedBaseURL, err := setupFormValue(r, "feed_base_url", 2048)
	if err != nil {
		return compliance.BootstrapInput{}, err
	}
	technicalContactEmail, err := setupFormValue(r, "technical_contact_email", 320)
	if err != nil {
		return compliance.BootstrapInput{}, err
	}
	licenseName, err := setupFormValue(r, "license_name", 160)
	if err != nil {
		return compliance.BootstrapInput{}, err
	}
	licenseURL, err := setupFormValue(r, "license_url", 2048)
	if err != nil {
		return compliance.BootstrapInput{}, err
	}
	publicationEnvironment, err := setupFormValue(r, "publication_environment", 64)
	if err != nil {
		return compliance.BootstrapInput{}, err
	}
	if publicBaseURL == "" {
		return compliance.BootstrapInput{}, fmt.Errorf("public base URL is required")
	}
	if feedBaseURL == "" {
		return compliance.BootstrapInput{}, fmt.Errorf("feed base URL is required")
	}
	if publicationEnvironment == "" || publicationEnvironment == "unknown" {
		publicationEnvironment = compliance.EnvironmentDev
	}
	return compliance.BootstrapInput{
		AgencyID:               principal.AgencyID,
		PublicBaseURL:          publicBaseURL,
		FeedBaseURL:            feedBaseURL,
		TechnicalContactEmail:  technicalContactEmail,
		LicenseName:            licenseName,
		LicenseURL:             licenseURL,
		PublicationEnvironment: publicationEnvironment,
		ActorID:                principal.Subject,
	}, nil
}

func setupFormValue(r *http.Request, name string, maxLen int) (string, error) {
	value := strings.TrimSpace(r.FormValue(name))
	if len(value) > maxLen {
		return "", fmt.Errorf("%s is too long", strings.ReplaceAll(name, "_", " "))
	}
	return value, nil
}

func (h *handler) buildOperationsPage(r *http.Request, principal auth.Principal, section string) operationsPage {
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
		Links: []evidenceLink{
			{Label: "OCI hosted evidence packet", Path: "docs/evidence/captured/oci-pilot/2026-04-24/README.md", UpdatedAt: "2026-04-24"},
			{Label: "Phase 23 agency-owned domain blocker", Path: "docs/agency-owned-domain-readiness.md", UpdatedAt: "Phase 23"},
			{Label: "Real-agency GTFS evidence scaffold", Path: "docs/evidence/real-agency-gtfs/README.md", UpdatedAt: "Phase 24"},
			{Label: "Device/AVL evidence scaffold", Path: "docs/evidence/device-avl/README.md", UpdatedAt: "Phase 25"},
			{Label: "Consumer submission tracker", Path: "docs/evidence/consumer-submissions/README.md", UpdatedAt: "2026-04-26"},
			{Label: "Consumer packet status JSON", Path: "docs/evidence/consumer-submissions/status.json", UpdatedAt: "2026-04-27"},
			{Label: "California readiness summary", Path: "docs/california-readiness-summary.md", UpdatedAt: "Phase 20"},
			{Label: "Compliance evidence checklist", Path: "docs/compliance-evidence-checklist.md", UpdatedAt: "repo docs"},
			{Label: "Small-agency pilot operations runbook", Path: "docs/runbooks/small-agency-pilot-operations.md", UpdatedAt: "Phase 17"},
			{Label: "Evidence redaction policy", Path: "docs/evidence/redaction-policy.md", UpdatedAt: "Phase 15"},
		},
		EvidenceUpdatedAt: "2026-04-26",
	}

	if reader, ok := h.store.(publicationConfigReader); ok {
		cfg, err := reader.PublicationConfig(r.Context(), principal.AgencyID)
		if err != nil {
			page.PublicationError = "publication metadata is not configured yet"
		} else {
			page.PublicationConfig = cfg
		}
	} else {
		page.PublicationError = "publication metadata config reader is not available in this runtime"
	}

	discovery, err := h.store.FeedDiscovery(r.Context(), principal.AgencyID, now)
	if err != nil {
		page.DiscoveryError = "publication metadata is not configured yet"
	} else {
		page.Discovery = discovery
		page.EnvironmentLabel = firstNonEmpty(discovery.PublicationEnvironment, page.EnvironmentLabel)
		page.ActiveFeedVersion = activeFeedVersion(discovery.Feeds)
		page.FeedsUpdatedAt = latestFeedTime(discovery)
	}

	scorecard, err := h.store.LatestScorecard(r.Context(), principal.AgencyID)
	if err != nil {
		page.ScorecardError = "no scorecard has been stored yet"
	} else {
		page.Scorecard = &scorecard
		t := scorecard.SnapshotAt.UTC()
		page.ScorecardUpdatedAt = &t
	}

	consumers, err := h.store.ListConsumers(r.Context(), principal.AgencyID)
	if err != nil {
		page.ConsumerError = "consumer status records are not available"
		page.Consumers = consumerTrackerStatuses(nil)
	} else {
		page.Consumers = consumerTrackerStatuses(consumers)
		page.RuntimeConsumers = runtimeConsumerStatuses(consumers)
		page.ConsumersUpdatedAt = latestConsumerTime(consumers)
	}

	page.Telemetry, page.TelemetryUpdatedAt, page.StaleCount, page.TelemetryError = h.telemetryViews(r, principal.AgencyID, now)
	page.TripUpdatesQuality = h.tripUpdatesQualityView(r, principal.AgencyID)

	bindings, err := h.devices.ListBindings(r.Context(), principal.AgencyID)
	if err != nil {
		page.DeviceError = "device bindings are not available"
	} else {
		page.Devices = bindings
	}
	page.DeviceRows = buildOperationsDeviceRows(page.Devices, page.Telemetry)
	page.DeviceOnboarding = operationsDeviceOnboardingUseCases()
	page.GTFSQuality = h.gtfsQualityTriage(r, principal.AgencyID, page.Discovery)
	page.GTFSQualityGuidance = buildOperationsGTFSQualityGuidance(page)
	page.ValidationHealth = h.validationHealthSummary(r, principal.AgencyID, page.Discovery, nil, nil)
	page.Reliability, page.ReliabilityError = h.reliabilitySummary(r, principal.AgencyID, now)
	page.FeedHealth = buildOperationsFeedHealth(page)
	page.Maintenance = buildOperationsMaintenance(page)
	page.SetupSteps = setupSteps(page)
	page.ReadinessItems = readinessItems(page)
	page.ReadinessV2 = buildOperationsReadinessV2(page)
	page.Checklist = buildOperatorChecklist(page)
	page.FirstRun = buildOperationsFirstRun(page)
	page.TelemetrySimulator = buildOperationsTelemetrySimulator(page)
	page.Launchpad = buildAgencyLaunchpad(page)
	page.SetupWizard = buildOperationsSetupWizard(page)
	page.ConnectorHub = buildConnectorHub(page)
	page.ConnectorTests = buildConnectorTests(page)
	page.Cockpit = buildOperationsCockpit(page)
	page.Help = buildOperationsHelpView(page.GeneratedAt, page.AgencyID, page.Section)
	page.ContextHelp = page.Help.ContextualHelp
	return page
}

func (h *handler) reliabilitySummary(r *http.Request, agencyID string, now time.Time) (compliance.ReliabilitySummary, string) {
	reader, ok := h.store.(reliabilityReader)
	if !ok {
		return compliance.BuildReliabilitySummary(compliance.ReliabilityInput{GeneratedAt: now, AgencyID: agencyID}), "reliability database reader is not available in this runtime"
	}
	records, err := reader.LatestReliabilityFeedHealth(r.Context(), agencyID, 200)
	if err != nil {
		return compliance.BuildReliabilitySummary(compliance.ReliabilityInput{GeneratedAt: now, AgencyID: agencyID}), "feed health snapshots are not available"
	}
	incidents, err := reader.ReliabilityIncidentRollup(r.Context(), agencyID, now, 10)
	if err != nil {
		return compliance.BuildReliabilitySummary(compliance.ReliabilityInput{GeneratedAt: now, AgencyID: agencyID, FeedHealthRecords: records}), "incident rollup is not available"
	}
	return compliance.BuildReliabilitySummary(compliance.ReliabilityInput{GeneratedAt: now, AgencyID: agencyID, FeedHealthRecords: records, Incidents: incidents}), ""
}

func (h *handler) validationHealthSummary(r *http.Request, agencyID string, discovery compliance.FeedDiscovery, extraRecords []compliance.ValidationReportRecord, artifactOverrides map[string]string) compliance.ValidationHealthSummary {
	records := append([]compliance.ValidationReportRecord(nil), extraRecords...)
	if reader, ok := h.store.(validationReportReader); ok {
		for _, feedType := range []string{"schedule", "vehicle_positions", "trip_updates", "alerts"} {
			validatorName := compliance.ValidatorNameForHealthID(compliance.ValidatorIDForHealthFeed(feedType))
			record, err := reader.LatestValidationReport(r.Context(), agencyID, feedType, validatorName)
			if err == nil && record != nil {
				records = append(records, *record)
			}
		}
	}
	artifactStatus := map[string]string{}
	for _, feed := range discovery.Feeds {
		if feed.FeedType == "schedule" {
			if strings.TrimSpace(feed.ActiveFeedVersionID) != "" && strings.TrimSpace(feed.CanonicalPublicURL) != "" {
				artifactStatus[feed.FeedType] = compliance.ValidationHealthArtifactAvailable
			} else {
				artifactStatus[feed.FeedType] = compliance.ValidationHealthArtifactUnavailable
			}
			continue
		}
		if feed.FeedType == "vehicle_positions" || feed.FeedType == "trip_updates" || feed.FeedType == "alerts" {
			if strings.TrimSpace(feed.CanonicalPublicURL) != "" {
				artifactStatus[feed.FeedType] = compliance.ValidationHealthArtifactAvailable
			} else {
				artifactStatus[feed.FeedType] = compliance.ValidationHealthArtifactUnavailable
			}
		}
	}
	for _, feedType := range []string{"schedule", "vehicle_positions", "trip_updates", "alerts"} {
		if artifactStatus[feedType] == "" {
			artifactStatus[feedType] = compliance.ValidationHealthArtifactUnavailable
		}
	}
	for feedType, status := range artifactOverrides {
		if status != "" {
			artifactStatus[feedType] = status
		}
	}
	return compliance.BuildValidationHealthSummary(compliance.ValidationHealthInput{
		GeneratedAt:              time.Now().UTC().Truncate(time.Second),
		AgencyID:                 agencyID,
		Discovery:                discovery,
		Registry:                 compliance.ValidatorRegistryFromEnv(),
		Records:                  records,
		ToolingStatusByValidator: validationHealthToolingStatusByValidator(),
		ArtifactStatusByFeed:     artifactStatus,
	})
}

func validationHealthToolingStatusByValidator() map[string]string {
	if strings.TrimSpace(os.Getenv("VALIDATOR_TOOLING_MODE")) == "stub" {
		return map[string]string{
			compliance.ValidationHealthStaticValidatorID:   compliance.ValidationHealthStatusStub,
			compliance.ValidationHealthRealtimeValidatorID: compliance.ValidationHealthStatusStub,
		}
	}
	return compliance.ValidationToolingStatusByValidator(compliance.ValidatorRegistryFromEnv())
}

func (h *handler) runValidationHealthAll(r *http.Request, principal auth.Principal, discovery compliance.FeedDiscovery) ([]compliance.ValidationReportRecord, map[string]string) {
	var records []compliance.ValidationReportRecord
	artifactStatus := map[string]string{}
	activeScheduleVersion, _ := scheduleFeedVersion(discovery.Feeds)
	for _, feedType := range []string{"schedule", "vehicle_positions", "trip_updates", "alerts"} {
		if feedType == "schedule" && activeScheduleVersion == "" {
			artifactStatus[feedType] = compliance.ValidationHealthArtifactUnavailable
			continue
		}
		if feedType != "schedule" && !feedDiscoveryHasPublicURL(discovery, feedType) {
			artifactStatus[feedType] = compliance.ValidationHealthArtifactUnavailable
			continue
		}
		result, err := h.runValidationForFeed(r, principal, feedType, activeScheduleVersion)
		if err != nil {
			artifactStatus[feedType] = compliance.ValidationHealthArtifactUnavailable
			continue
		}
		artifactStatus[feedType] = compliance.ValidationHealthArtifactAvailable
		records = append(records, compliance.ValidationReportRecord{Result: result, CreatedAt: time.Now().UTC().Truncate(time.Second)})
	}
	return records, artifactStatus
}

func feedDiscoveryHasPublicURL(discovery compliance.FeedDiscovery, feedType string) bool {
	for _, feed := range discovery.Feeds {
		if feed.FeedType == feedType && strings.TrimSpace(feed.CanonicalPublicURL) != "" {
			return true
		}
	}
	return false
}

func (h *handler) gtfsQualityTriage(r *http.Request, agencyID string, discovery compliance.FeedDiscovery) compliance.GTFSQualityTriage {
	activeFeedVersionID, activeRevision := scheduleFeedVersion(discovery.Feeds)
	reader, ok := h.store.(validationReportReader)
	if !ok {
		return compliance.BuildGTFSQualityTriage(compliance.GTFSQualityTriageInput{ActiveFeedVersionID: activeFeedVersionID, ActiveFeedRevisionTime: activeRevision})
	}
	canonical, err := reader.LatestValidationReport(r.Context(), agencyID, "schedule", compliance.CanonicalStaticValidatorName)
	if err != nil {
		canonical = nil
	}
	internalImporter, err := reader.LatestValidationReport(r.Context(), agencyID, "schedule", compliance.InternalGTFSImportValidatorName)
	if err != nil {
		internalImporter = nil
	}
	return compliance.BuildGTFSQualityTriage(compliance.GTFSQualityTriageInput{
		Canonical:              canonical,
		InternalImporter:       internalImporter,
		ActiveFeedVersionID:    activeFeedVersionID,
		ActiveFeedRevisionTime: activeRevision,
	})
}

func (h *handler) tripUpdatesQualityView(r *http.Request, agencyID string) tripUpdatesQualityView {
	reader, ok := h.store.(tripUpdatesDiagnosticsReader)
	if !ok {
		return tripUpdatesQualityView{Message: "no Trip Updates diagnostics recorded yet"}
	}
	summary, err := reader.LatestTripUpdatesDiagnostics(r.Context(), agencyID)
	if err != nil {
		return tripUpdatesQualityView{Message: "Trip Updates diagnostics are not available"}
	}
	if !summary.Recorded {
		return tripUpdatesQualityView{Message: "no Trip Updates diagnostics recorded yet"}
	}
	snapshotAt := summary.SnapshotAt.UTC()
	metrics := summary.Metrics
	return tripUpdatesQualityView{
		Recorded:                      true,
		SnapshotAt:                    &snapshotAt,
		AdapterName:                   summary.AdapterName,
		DiagnosticsStatus:             summary.DiagnosticsStatus,
		DiagnosticsReason:             summary.DiagnosticsReason,
		ActiveFeedVersionID:           summary.ActiveFeedVersionID,
		DiagnosticsPersistenceOutcome: summary.DiagnosticsPersistenceOutcome,
		UnknownAssignmentRate:         rateText(metrics.UnknownAssignmentRate),
		AmbiguousAssignmentRate:       rateText(metrics.AmbiguousAssignmentRate),
		StaleTelemetryRate:            rateText(metrics.StaleTelemetryRate),
		TripUpdatesCoverageRate:       rateText(metrics.TripUpdatesCoverageRate),
		FutureStopCoverageRate:        rateText(metrics.FutureStopCoverageRate),
		EligiblePredictionCandidates:  metrics.EligiblePredictionCandidates,
		TripUpdatesEmitted:            metrics.TripUpdatesEmitted,
		UnknownAssignments:            metrics.UnknownAssignments,
		AmbiguousAssignments:          metrics.AmbiguousAssignments,
		StaleTelemetryRows:            metrics.StaleTelemetryRows,
		ManualOverrideAssignments:     metrics.ManualOverrideAssignments,
		CanceledTripsEmitted:          metrics.CanceledTripsEmitted,
		CancellationAlertLinksMissing: metrics.CancellationAlertLinksMissing,
		WithheldByReason:              countViews(metrics.WithheldByReason),
	}
}

func rateText(rate prediction.RateMetric) string {
	if rate.Status == "not_applicable" {
		if rate.NotApplicableReason != "" {
			return "not applicable: " + rate.NotApplicableReason
		}
		return "not applicable"
	}
	if rate.Percent == nil {
		return "not available"
	}
	return fmt.Sprintf("%.1f%% (%d/%d)", *rate.Percent, rate.Numerator, rate.Denominator)
}

func countViews(counts map[string]int) []countView {
	keys := make([]string, 0, len(counts))
	for key, count := range counts {
		if count <= 0 {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	views := make([]countView, 0, len(keys))
	for _, key := range keys {
		views = append(views, countView{Label: key, Count: counts[key]})
	}
	return views
}

func (h *handler) telemetryViews(r *http.Request, agencyID string, now time.Time) ([]telemetryView, *time.Time, int, string) {
	if h.telemetry == nil {
		return nil, nil, 0, "telemetry repository is not available in this runtime"
	}
	latest, err := h.telemetry.ListLatestByAgency(r.Context(), agencyID, defaultTelemetryLimit)
	if err != nil {
		return nil, nil, 0, "latest telemetry is not available"
	}
	vehicleIDs := make([]string, 0, len(latest))
	for _, event := range latest {
		vehicleIDs = append(vehicleIDs, event.VehicleID)
	}
	assignments := map[string]state.Assignment{}
	if h.state != nil {
		if rows, err := h.state.ListCurrentAssignments(r.Context(), agencyID, vehicleIDs); err == nil {
			assignments = rows
		}
	}
	threshold := staleThreshold()
	var newest *time.Time
	var staleCount int
	views := make([]telemetryView, 0, len(latest))
	for _, event := range latest {
		age := now.Sub(event.Timestamp)
		stale := age > threshold
		if stale {
			staleCount++
		}
		observed := event.Timestamp.UTC()
		if newest == nil || observed.After(*newest) {
			t := observed
			newest = &t
		}
		view := telemetryView{
			VehicleID:  event.VehicleID,
			DeviceID:   event.DeviceID,
			ObservedAt: observed,
			ReceivedAt: event.ReceivedAt.UTC(),
			AgeSeconds: int64(age.Seconds()),
			Stale:      stale,
		}
		if assignment, ok := assignments[event.VehicleID]; ok {
			view.AssignmentState = string(assignment.State)
			view.RouteID = assignment.RouteID
			view.TripID = assignment.TripID
			view.Confidence = fmt.Sprintf("%.2f", assignment.Confidence)
			view.ReasonCodes = append([]string(nil), assignment.ReasonCodes...)
			view.DegradedState = string(assignment.DegradedState)
			view.AssignmentSource = string(assignment.AssignmentSource)
			if !assignment.ActiveFrom.IsZero() {
				t := assignment.ActiveFrom.UTC()
				view.AssignmentAt = &t
			}
		}
		views = append(views, view)
	}
	return views, newest, staleCount, ""
}

func activeFeedVersion(feeds []compliance.FeedMetadata) string {
	for _, feed := range feeds {
		if feed.FeedType == "schedule" && feed.ActiveFeedVersionID != "" {
			return feed.ActiveFeedVersionID
		}
	}
	for _, feed := range feeds {
		if feed.ActiveFeedVersionID != "" {
			return feed.ActiveFeedVersionID
		}
	}
	return ""
}

func scheduleFeedVersion(feeds []compliance.FeedMetadata) (string, *time.Time) {
	for _, feed := range feeds {
		if feed.FeedType == "schedule" {
			return feed.ActiveFeedVersionID, feed.RevisionTimestamp
		}
	}
	return "", nil
}

func latestFeedTime(discovery compliance.FeedDiscovery) *time.Time {
	var latest *time.Time
	candidates := []time.Time{discovery.GeneratedAt}
	for _, feed := range discovery.Feeds {
		if feed.RevisionTimestamp != nil {
			candidates = append(candidates, *feed.RevisionTimestamp)
		}
		if feed.LastValidationAt != nil {
			candidates = append(candidates, *feed.LastValidationAt)
		}
		if feed.LastHealthAt != nil {
			candidates = append(candidates, *feed.LastHealthAt)
		}
	}
	for _, candidate := range candidates {
		t := candidate.UTC()
		if latest == nil || t.After(*latest) {
			latest = &t
		}
	}
	return latest
}

func latestConsumerTime(consumers []compliance.ConsumerRecord) *time.Time {
	var latest *time.Time
	for _, consumer := range consumers {
		t := consumer.UpdatedAt.UTC()
		if latest == nil || t.After(*latest) {
			latest = &t
		}
	}
	return latest
}

func setupSteps(page operationsPage) []setupStepView {
	steps := []setupStepView{
		{
			Name:       "Agency metadata",
			Status:     missingOrValue(page.Discovery.AgencyName, "missing"),
			Source:     "publication metadata",
			Evidence:   firstNonEmpty(page.Discovery.AgencyName, page.PublicationError, page.DiscoveryError),
			NextAction: "Review agency name in GTFS Studio and publication metadata before making deployment claims.",
			ActionURL:  "/admin/operations/setup#publication-metadata",
		},
		{
			Name:       "License and contact metadata",
			Status:     readinessStatus(page.Discovery.Readiness.LicenseComplete && page.Discovery.Readiness.ContactComplete, page.DiscoveryError),
			Source:     "publication metadata",
			Evidence:   licenseContactEvidence(page),
			NextAction: "Enter an agency-approved open license URL and technical contact, or keep this marked missing.",
			ActionURL:  "/admin/operations/setup#publication-metadata",
		},
		{
			Name:       "GTFS import or GTFS Studio path",
			Status:     missingOrValue(page.ActiveFeedVersion, "missing active feed"),
			Source:     "feed discovery",
			Evidence:   feedEvidence(page, "schedule"),
			NextAction: "Use browser GTFS ZIP upload or safe URL import first; use GTFS Studio when typed draft authoring is needed.",
			ActionURL:  "/admin/operations/gtfs-import",
		},
		{
			Name:       "Publication bootstrap",
			Status:     readinessStatus(page.Discovery.Readiness.AllRequiredFeedsListed, page.DiscoveryError),
			Source:     "feed discovery",
			Evidence:   allFeedURLsEvidence(page),
			NextAction: "Use the publication form to store public/feed base URLs and metadata after confirming the values.",
			ActionURL:  "/admin/operations/setup#publication-metadata",
		},
		{
			Name:       "Device token setup",
			Status:     countStatus(len(page.Devices), "binding records"),
			Source:     "device bindings",
			Evidence:   deviceEvidence(page),
			NextAction: "Rotate/rebind a one-time device token and store it outside this repo.",
			ActionURL:  "/admin/operations/devices",
		},
		{
			Name:       "First telemetry event",
			Status:     telemetryStatus(page),
			Source:     "telemetry repository",
			Evidence:   telemetryEvidence(page),
			NextAction: "Send a sample telemetry event using the device onboarding helper, then review freshness.",
			ActionURL:  "/admin/operations/telemetry",
		},
		{
			Name:       "First validation run",
			Status:     validationStatus(page),
			Source:     "validation records",
			Evidence:   validationEvidence(page),
			NextAction: "Review private validator health, then run allowlisted validators from the health page or existing admin validation workflow.",
			ActionURL:  "/admin/operations/validation-health",
		},
		{
			Name:       "Public feed verification",
			Status:     readinessStatus(page.Discovery.Readiness.AllRequiredFeedsListed, page.DiscoveryError),
			Source:     "feed discovery",
			Evidence:   allFeedURLsEvidence(page),
			NextAction: "Review public feed URLs and health records. Verification is supporting evidence only.",
			ActionURL:  "/admin/operations/feeds",
		},
		{
			Name:       "Alerts setup",
			Status:     feedStatus(page, "alerts"),
			Source:     "feed discovery and Alerts Console",
			Evidence:   feedEvidence(page, "alerts"),
			NextAction: "Use the Alerts Console to create, publish, or archive service alerts as needed.",
			ActionURL:  "/admin/alerts/console",
		},
		{
			Name:       "Consumer packet/status review",
			Status:     countStatus(len(page.Consumers), "prepared docs tracker targets"),
			Source:     "docs/evidence tracker",
			Evidence:   "Phase 20 docs tracker records prepared packets only.",
			NextAction: "Review packet paths and submission workflow; do not change statuses without target-originated evidence.",
			ActionURL:  "/admin/operations/consumers",
		},
		{
			Name:       "Evidence/readiness review",
			Status:     countStatus(len(page.Links), "evidence links"),
			Source:     "evidence links",
			Evidence:   "Repo evidence links are navigation aids and do not prove consumer acceptance or compliance.",
			NextAction: "Review OCI pilot evidence, agency-owned-domain blocker, GTFS/device scaffolds, and readiness docs.",
			ActionURL:  "/admin/operations/evidence",
		},
	}
	return steps
}

func readinessItems(page operationsPage) []readinessItemView {
	consumerEvidence := "docs tracker targets are prepared packet records only"
	if len(page.Consumers) > 0 {
		consumerEvidence = fmt.Sprintf("%d docs tracker targets shown as prepared packet records only", len(page.Consumers))
	}
	if len(page.RuntimeConsumers) > 0 {
		consumerEvidence += fmt.Sprintf("; %d runtime deployment workflow records are visible separately", len(page.RuntimeConsumers))
	}
	return []readinessItemView{
		{
			Name:          "Stable public URLs",
			Status:        readinessStatus(page.Discovery.Readiness.AllRequiredFeedsListed && page.Discovery.Readiness.HTTPSURLs, page.DiscoveryError),
			Source:        "feed discovery and published_feed records",
			Evidence:      stableURLEvidence(page),
			NextAction:    "Confirm every public URL is stable, HTTPS in deployment, and served through the intended public feed root.",
			ClaimBoundary: "URL records are readiness signals, not agency-owned-domain or consumer acceptance evidence.",
		},
		{
			Name:          "Static GTFS feed",
			Status:        feedStatus(page, "schedule"),
			Source:        "published_feed schedule record and validation records",
			Evidence:      feedEvidence(page, "schedule"),
			NextAction:    "Import or publish the agency schedule, then run the allowlisted static GTFS validator.",
			ClaimBoundary: "A listed schedule feed is not a CAL-ITP/Caltrans compliance claim by itself.",
		},
		{
			Name:          "Vehicle Positions",
			Status:        feedStatus(page, "vehicle_positions"),
			Source:        "published_feed Vehicle Positions record, validation records, and feed health snapshots",
			Evidence:      feedEvidence(page, "vehicle_positions"),
			NextAction:    "Verify fresh telemetry, public protobuf output, feed health, and GTFS-Realtime validation.",
			ClaimBoundary: "Vehicle Positions readiness does not prove real device reliability or consumer acceptance.",
		},
		{
			Name:          "Trip Updates",
			Status:        feedStatus(page, "trip_updates"),
			Source:        "published_feed Trip Updates record, validation records, and Trip Updates diagnostics",
			Evidence:      tripUpdatesReadinessEvidence(page),
			NextAction:    "Review prediction diagnostics, coverage, withheld cases, and GTFS-Realtime validation.",
			ClaimBoundary: "Trip Updates diagnostics are not production-grade ETA quality evidence.",
		},
		{
			Name:          "Alerts",
			Status:        feedStatus(page, "alerts"),
			Source:        "published_feed Alerts record, validation records, and Alerts Console workflow",
			Evidence:      feedEvidence(page, "alerts"),
			NextAction:    "Use the Alerts Console for lifecycle checks, then validate the public Alerts feed.",
			ClaimBoundary: "Alert feed availability does not prove consumer display or agency approval.",
		},
		{
			Name:          "License/contact metadata",
			Status:        readinessStatus(page.Discovery.Readiness.LicenseComplete && page.Discovery.Readiness.ContactComplete, page.DiscoveryError),
			Source:        "feed_config, published_feed, and /public/feeds.json metadata",
			Evidence:      licenseContactEvidence(page),
			NextAction:    "Replace placeholders with agency-approved open license and monitored technical contact values.",
			ClaimBoundary: "Metadata completeness is not agency approval unless separate retained approval exists.",
		},
		{
			Name:          "Validation status",
			Status:        page.ValidationHealth.OverallStatus,
			Source:        "validation records and private validator health",
			Evidence:      validationEvidence(page),
			NextAction:    "Open private validator health diagnostics and run the admin-only allowlisted validator action if needed.",
			ClaimBoundary: "Validator results support readiness review, but do not prove consumer acceptance.",
		},
		{
			Name:          "Telemetry freshness",
			Status:        telemetryStatus(page),
			Source:        "telemetry latest rows and current assignment summaries",
			Evidence:      telemetryEvidence(page),
			NextAction:    "Send validated telemetry through device credentials and resolve stale or unmatched vehicles.",
			ClaimBoundary: "Fresh local telemetry is not real vendor compatibility or fleet reliability proof.",
		},
		{
			Name:          "Operations status",
			Status:        operationsStatus(page),
			Source:        "scorecard snapshots and Operations Console summaries",
			Evidence:      operationsEvidence(page),
			NextAction:    "Run scorecard export, feed monitor, backup, and restore-drill workflows in the operator environment.",
			ClaimBoundary: "Operations records are deployment health signals, not managed hosting or universal deployment claims.",
		},
		{
			Name:          "Consumer packet preparedness",
			Status:        countStatus(len(page.Consumers), "prepared docs tracker targets"),
			Source:        "docs/evidence tracker paths and runtime consumer workflow records",
			Evidence:      consumerEvidence,
			NextAction:    "Review packet paths and official submission workflow; change statuses only with target-originated evidence.",
			ClaimBoundary: "Prepared packets are not submitted, under review, accepted, listed, displayed, or ingested.",
		},
	}
}

func missingOrValue(value string, missing string) string {
	if strings.TrimSpace(value) == "" {
		return missing
	}
	return value
}

func readinessStatus(ok bool, missingReason string) string {
	if missingReason != "" {
		return "missing"
	}
	if ok {
		return "recorded"
	}
	return "missing"
}

func countStatus(count int, label string) string {
	if count == 0 {
		return "missing"
	}
	return fmt.Sprintf("%d %s", count, label)
}

func licenseContactEvidence(page operationsPage) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	return fmt.Sprintf("license name=%q, license URL=%q, technical contact=%q", page.Discovery.License.Name, page.Discovery.License.URL, page.Discovery.TechnicalContactEmail)
}

func feedStatus(page operationsPage, feedType string) string {
	if page.DiscoveryError != "" {
		return "missing"
	}
	for _, feed := range page.Discovery.Feeds {
		if feed.FeedType == feedType {
			if feed.CanonicalPublicURL == "" {
				return "missing URL"
			}
			if feed.LastValidationStatus == "" || feed.LastValidationStatus == "not_run" {
				return "URL listed; validation not run"
			}
			return "URL listed; validation " + feed.LastValidationStatus
		}
	}
	return "missing"
}

func feedEvidence(page operationsPage, feedType string) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	for _, feed := range page.Discovery.Feeds {
		if feed.FeedType == feedType {
			return fmt.Sprintf("URL=%q, active feed version=%q, validation=%s at %s", feed.CanonicalPublicURL, feed.ActiveFeedVersionID, firstNonEmpty(feed.LastValidationStatus, "not_run"), formatTimeForText(feed.LastValidationAt))
		}
	}
	return "no " + feedType + " feed metadata record"
}

func allFeedURLsEvidence(page operationsPage) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	return fmt.Sprintf("%d feed discovery records; all required listed=%t", len(page.Discovery.Feeds), page.Discovery.Readiness.AllRequiredFeedsListed)
}

func stableURLEvidence(page operationsPage) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	return fmt.Sprintf("%d feed discovery records; all required listed=%t; all HTTPS=%t", len(page.Discovery.Feeds), page.Discovery.Readiness.AllRequiredFeedsListed, page.Discovery.Readiness.HTTPSURLs)
}

func deviceEvidence(page operationsPage) string {
	if page.DeviceError != "" {
		return page.DeviceError
	}
	if len(page.Devices) == 0 {
		return "no device binding records"
	}
	return fmt.Sprintf("%d device binding records; tokens are not rendered", len(page.Devices))
}

func telemetryStatus(page operationsPage) string {
	if page.TelemetryError != "" {
		return "not available"
	}
	if len(page.Telemetry) == 0 {
		return "not observed yet"
	}
	if page.StaleCount > 0 {
		return fmt.Sprintf("%d vehicles observed; %d stale", len(page.Telemetry), page.StaleCount)
	}
	return fmt.Sprintf("%d vehicles observed", len(page.Telemetry))
}

func telemetryEvidence(page operationsPage) string {
	if page.TelemetryError != "" {
		return page.TelemetryError
	}
	if page.TelemetryUpdatedAt == nil {
		return "no latest telemetry rows"
	}
	return "latest observed telemetry at " + formatTimeForText(page.TelemetryUpdatedAt)
}

func validationStatus(page operationsPage) string {
	if page.DiscoveryError != "" {
		return "not run yet"
	}
	for _, feed := range page.Discovery.Feeds {
		if feed.LastValidationAt != nil {
			return "supporting records exist"
		}
	}
	return "not run yet"
}

func validationEvidence(page operationsPage) string {
	if page.DiscoveryError != "" {
		return page.DiscoveryError
	}
	var parts []string
	for _, feed := range page.Discovery.Feeds {
		parts = append(parts, feed.FeedType+"="+firstNonEmpty(feed.LastValidationStatus, "not_run")+" at "+formatTimeForText(feed.LastValidationAt))
	}
	if len(parts) == 0 {
		return "no validation records"
	}
	return strings.Join(parts, "; ")
}

func tripUpdatesReadinessEvidence(page operationsPage) string {
	feed := feedEvidence(page, "trip_updates")
	if !page.TripUpdatesQuality.Recorded {
		return feed + "; diagnostics=" + page.TripUpdatesQuality.Message
	}
	return fmt.Sprintf("%s; diagnostics=%s/%s at %s; coverage=%s; future stops=%s",
		feed,
		page.TripUpdatesQuality.DiagnosticsStatus,
		page.TripUpdatesQuality.DiagnosticsReason,
		formatTimeForText(page.TripUpdatesQuality.SnapshotAt),
		page.TripUpdatesQuality.TripUpdatesCoverageRate,
		page.TripUpdatesQuality.FutureStopCoverageRate,
	)
}

func operationsStatus(page operationsPage) string {
	if page.Scorecard != nil {
		return page.Scorecard.OverallStatus
	}
	if page.ScorecardError != "" {
		return "missing scorecard"
	}
	return "not available"
}

func operationsEvidence(page operationsPage) string {
	if page.Scorecard != nil {
		return fmt.Sprintf("latest scorecard overall=%s at %s", page.Scorecard.OverallStatus, formatTimeForText(page.ScorecardUpdatedAt))
	}
	return firstNonEmpty(page.ScorecardError, "no scorecard snapshot available")
}

func formatTimeForText(t *time.Time) string {
	if t == nil || t.IsZero() {
		return "not available"
	}
	return t.UTC().Format(time.RFC3339)
}

func runtimeConsumerStatuses(records []compliance.ConsumerRecord) []consumerStatusView {
	statuses := make([]consumerStatusView, 0, len(records))
	for _, record := range records {
		updated := record.UpdatedAt.UTC()
		statuses = append(statuses, consumerStatusView{
			Name:      record.ConsumerName,
			Status:    record.Status,
			UpdatedAt: &updated,
			Source:    "runtime DB deployment workflow record",
			Notes:     record.Notes,
		})
	}
	return statuses
}

func consumerTrackerStatuses(records []compliance.ConsumerRecord) []consumerStatusView {
	byName := map[string]compliance.ConsumerRecord{}
	for _, record := range records {
		byName[record.ConsumerName] = record
	}
	targets := []consumerStatusView{
		{Name: "Google Maps", Status: "prepared", Source: "docs/evidence tracker", CurrentPath: "docs/evidence/consumer-submissions/current/google-maps.md", PacketPath: "docs/evidence/consumer-submissions/packets/google-maps/README.md"},
		{Name: "Apple Maps", Status: "prepared", Source: "docs/evidence tracker", CurrentPath: "docs/evidence/consumer-submissions/current/apple-maps.md", PacketPath: "docs/evidence/consumer-submissions/packets/apple-maps/README.md"},
		{Name: "Transit App", Status: "prepared", Source: "docs/evidence tracker", CurrentPath: "docs/evidence/consumer-submissions/current/transit-app.md", PacketPath: "docs/evidence/consumer-submissions/packets/transit-app/README.md"},
		{Name: "Bing Maps", Status: "prepared", Source: "docs/evidence tracker", CurrentPath: "docs/evidence/consumer-submissions/current/bing-maps.md", PacketPath: "docs/evidence/consumer-submissions/packets/bing-maps/README.md"},
		{Name: "Moovit", Status: "prepared", Source: "docs/evidence tracker", CurrentPath: "docs/evidence/consumer-submissions/current/moovit.md", PacketPath: "docs/evidence/consumer-submissions/packets/moovit/README.md"},
		{Name: "Mobility Database", Status: "prepared", Source: "docs/evidence tracker", CurrentPath: "docs/evidence/consumer-submissions/current/mobility-database.md", PacketPath: "docs/evidence/consumer-submissions/packets/mobility-database/README.md"},
		{Name: "transit.land", Status: "prepared", Source: "docs/evidence tracker", CurrentPath: "docs/evidence/consumer-submissions/current/transit-land.md", PacketPath: "docs/evidence/consumer-submissions/packets/transit-land/README.md"},
	}
	for i := range targets {
		targets[i].Notes = "Prepared packet only; not submitted, under review, accepted, or ingested."
		if record, ok := byName[targets[i].Name]; ok {
			updated := record.UpdatedAt.UTC()
			targets[i].UpdatedAt = &updated
			targets[i].Notes += " Runtime DB workflow record currently says " + record.Status + "."
		}
	}
	return targets
}

func staleThreshold() time.Duration {
	raw := strings.TrimSpace(os.Getenv("STALE_TELEMETRY_TTL_SECONDS"))
	if raw == "" {
		return defaultStaleSeconds * time.Second
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return defaultStaleSeconds * time.Second
	}
	return time.Duration(seconds) * time.Second
}

func feedOrder(feedType string) int {
	switch feedType {
	case "schedule":
		return 0
	case "vehicle_positions":
		return 1
	case "trip_updates":
		return 2
	case "alerts":
		return 3
	default:
		return 99
	}
}

func csrfToken(secret string, principal auth.Principal) string {
	if strings.TrimSpace(secret) == "" {
		return ""
	}
	return auth.CSRFToken(secret, principal)
}

func renderOperationsTemplate(w http.ResponseWriter, name string, data operationsPage) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := operationsTemplates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var operationsTemplates = template.Must(template.New("operations").Funcs(template.FuncMap{
	"formatTime": func(t time.Time) string {
		if t.IsZero() {
			return "not available"
		}
		return t.UTC().Format(time.RFC3339)
	},
	"formatTimePtr": func(t *time.Time) string {
		if t == nil || t.IsZero() {
			return "not available"
		}
		return t.UTC().Format(time.RFC3339)
	},
	"formatBoolPtr": func(v *bool) string {
		if v == nil {
			return "not observed"
		}
		if *v {
			return "true"
		}
		return "false"
	},
	"formatFloatPtr": func(v *float64) string {
		if v == nil {
			return "not observed"
		}
		return fmt.Sprintf("%.2f", *v)
	},
	"formatInt64Ptr": func(v *int64) string {
		if v == nil {
			return "not available"
		}
		return strconv.FormatInt(*v, 10)
	},
	"join": strings.Join,
	"lower": func(value string) string {
		return strings.ToLower(strings.ReplaceAll(value, "_", "-"))
	},
	"feedURL": func(discovery compliance.FeedDiscovery, feedType string) string {
		for _, feed := range discovery.Feeds {
			if feed.FeedType == feedType {
				return feed.CanonicalPublicURL
			}
		}
		return ""
	},
	"sortedFeeds": func(feeds []compliance.FeedMetadata) []compliance.FeedMetadata {
		out := append([]compliance.FeedMetadata(nil), feeds...)
		sort.SliceStable(out, func(i, j int) bool { return feedOrder(out[i].FeedType) < feedOrder(out[j].FeedType) })
		return out
	},
	"publicationEnvValue": func(page operationsPage) string {
		if page.PublicationConfig.PublicationEnvironment != "" {
			return page.PublicationConfig.PublicationEnvironment
		}
		if page.Discovery.PublicationEnvironment != "" {
			return page.Discovery.PublicationEnvironment
		}
		if page.EnvironmentLabel != "unknown" {
			return page.EnvironmentLabel
		}
		return ""
	},
	"humanHeuristic":           humanHeuristicLabel,
	"gtfsQualityLikelyOwner":   gtfsQualityLikelyOwner,
	"gtfsQualityAffectedFiles": gtfsQualityAffectedFiles,
	"gtfsQualitySafeFixPath":   gtfsQualitySafeFixPath,
	"gtfsQualityVerifyWith":    gtfsQualityVerifyWith,
	"gtfsQualityEscalation":    gtfsQualityEscalation,
	"gtfsQualityGuidanceClass": gtfsQualityGuidanceClass,
}).Parse(`
{{define "layoutStart"}}
<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>{{.Title}}</title>
<style>
*{box-sizing:border-box}body{font-family:system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;margin:2rem;line-height:1.4;color:#1f2933;overflow-wrap:anywhere}.skip-link{position:absolute;left:-999px;top:.5rem;background:#1f2933;color:#fff;padding:.5rem .75rem;border-radius:4px;z-index:10}.skip-link:focus,.skip-link:focus-visible{left:.5rem}.operations-header{margin-bottom:1rem}a:focus-visible,button:focus-visible,input:focus-visible,select:focus-visible,textarea:focus-visible,main:focus-visible{outline:3px solid #2563eb;outline-offset:2px}
.operations-nav{display:grid;grid-template-columns:repeat(auto-fit,minmax(13rem,1fr));gap:.75rem;margin:1rem 0 1.25rem}.nav-group{border:1px solid #d8dee4;border-radius:6px;padding:.55rem;background:#fff}.nav-group-label{font-weight:700;margin:0 0 .4rem}.nav-links{display:flex;flex-wrap:wrap;gap:.35rem}.nav-link{border:1px solid #d8dee4;border-radius:4px;padding:.45rem .6rem;min-height:2.25rem;text-decoration:none;color:#1f2933;background:#fff}.nav-link:focus,.nav-link:hover{border-color:#6b7280;background:#f6f8fa}.nav-link.current{border-color:#1f2933;background:#1f2933;color:#fff}
table{border-collapse:collapse;width:100%;margin:1rem 0} th,td{border:1px solid #d8dee4;padding:.45rem;text-align:left;vertical-align:top}
th{background:#f6f8fa}.pill{display:inline-block;border:1px solid #c8d1dc;border-radius:3px;padding:.1rem .35rem;background:#f6f8fa}
.hero{border:1px solid #c8d1dc;background:#f8fafc;padding:1rem;border-radius:6px;margin:1rem 0}.card-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(16rem,1fr));gap:1rem;margin:1rem 0}.card{border:1px solid #d8dee4;border-radius:6px;padding:1rem;background:#fff}.card h3{margin-top:0}.card p{margin:.4rem 0}.status{font-weight:600}.copy-value{display:block;border:1px solid #d8dee4;background:#f6f8fa;border-radius:4px;padding:.45rem;white-space:pre-wrap;overflow-wrap:anywhere}.context-help{border:1px solid #c8d1dc;background:#f8fafc;border-radius:6px;padding:1rem;margin:1rem 0}.context-help h2{font-size:1.05rem;margin:0 0 .6rem}.context-help-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(14rem,1fr));gap:.75rem}.context-help-topic{border-left:3px solid #2563eb;padding-left:.65rem}.context-help-topic h3{font-size:1rem;margin:.1rem 0}.context-help-topic p{margin:.3rem 0}
.warning{background:#fff8c5}.ok{background:#dafbe1}.bad{background:#ffebe9}.muted{color:#59636e}.token{border:1px solid #f0c36d;background:#fff8c5;padding:1rem}
form{margin:1rem 0} label{display:block;margin:.35rem 0} input,select,textarea{min-width:22rem;max-width:100%;padding:.45rem} button{padding:.5rem .8rem;min-height:2.25rem}
@media (max-width:700px){body{margin:0;padding:1rem}.operations-nav,.card-grid,.context-help-grid{grid-template-columns:1fr}.nav-links{display:grid;grid-template-columns:1fr}.nav-link,button{width:100%}table{display:block;max-width:100%;overflow-x:auto;-webkit-overflow-scrolling:touch}input,select,textarea{min-width:0;width:100%}}
</style></head><body>
<a class="skip-link" href="#operations-main">Skip to main content</a>
<header class="operations-header">
<h1>{{.Title}}</h1>
<p>Agency: <strong>{{.AgencyID}}</strong> · environment: <span class="pill">{{.EnvironmentLabel}}</span> · generated: {{formatTime .GeneratedAt}}</p>
</header>
<nav class="operations-nav" aria-label="Operations Console sections">
{{range .NavGroups}}<section class="nav-group" aria-label="{{.Label}}">
<p class="nav-group-label">{{.Label}}</p>
<div class="nav-links">{{range .Items}}<a class="nav-link{{if .Current}} current{{end}}" href="{{.Href}}"{{if .Current}} aria-current="page"{{end}}>{{.Label}}</a>{{end}}</div>
</section>{{end}}
</nav>
{{if .ContextHelp.Topics}}<aside class="context-help" aria-labelledby="context-help-heading">
<h2 id="context-help-heading">Help for {{.ContextHelp.Label}}</h2>
<div class="context-help-grid">{{range .ContextHelp.Topics}}<section class="context-help-topic"><h3>{{.Label}}</h3><p>{{.Summary}}</p><p><strong>Next:</strong> {{.NextAction}}</p><p><a href="/admin/operations/help#help-{{.ID}}">Open topic</a></p></section>{{end}}</div>
<p class="muted"><a href="{{.ContextHelp.AllTopicsURL}}">Open all help topics</a> · <a href="{{.ContextHelp.JSONURL}}">Export private help JSON</a></p>
</aside>{{end}}
<main id="operations-main" tabindex="-1">
{{end}}
{{define "layoutEnd"}}</main></body></html>{{end}}

{{define "help"}}
{{template "layoutStart" .}}
<h2>Operations Console Help</h2>
<p class="warning">{{.Help.Boundary}}</p>
<p><a href="/admin/operations/help.json">Export private help JSON</a> · <a href="/admin/operations">Back to Operations Console</a></p>
<div class="card-grid" aria-label="Help topics">
{{range .Help.Topics}}<section class="card" id="help-{{.ID}}">
<h3>{{.Label}}</h3>
<p>{{.Summary}}</p>
{{if .PluginDefinition}}<p><strong>Plugin definition:</strong> {{.PluginDefinition}}</p>{{end}}
<p><strong>Review:</strong> {{.WhatToReview}}</p>
<p><strong>Next action:</strong> {{.NextAction}}</p>
<p><strong>Boundary:</strong> {{.DoesNotProve}}</p>
<p><strong>Claim boundary:</strong> {{.ClaimBoundary}}</p>
<p><strong>Console:</strong> {{range .AdminLinks}}<a href="{{.}}">{{.}}</a> {{end}}</p>
<p><strong>Docs:</strong> {{range .DocsLinks}}<code>{{.}}</code> {{end}}</p>
</section>{{end}}
</div>
<h3>Claim Flags</h3>
<table><tbody>
<tr><th><code>backend_command_execution_enabled</code></th><td>{{.Help.ClaimFlags.BackendCommandExecutionEnabled}}</td></tr>
<tr><th><code>cache_diagnostics_read</code></th><td>{{.Help.ClaimFlags.CacheDiagnosticsRead}}</td></tr>
<tr><th><code>external_network_contacted</code></th><td>{{.Help.ClaimFlags.ExternalNetworkContacted}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.Help.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>final_root_evidence_created</code></th><td>{{.Help.ClaimFlags.FinalRootEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.Help.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>secrets_collected</code></th><td>{{.Help.ClaimFlags.SecretsCollected}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.Help.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.Help.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>agency_approval_claimed</code></th><td>{{.Help.ClaimFlags.AgencyApprovalClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.Help.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th><code>public_launch_claimed</code></th><td>{{.Help.ClaimFlags.PublicLaunchClaimed}}</td></tr>
<tr><th><code>hosted_saas_claimed</code></th><td>{{.Help.ClaimFlags.HostedSaaSClaimed}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.Help.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>hardware_certification_claimed</code></th><td>{{.Help.ClaimFlags.HardwareCertificationClaimed}}</td></tr>
<tr><th><code>production_avl_reliability_claimed</code></th><td>{{.Help.ClaimFlags.ProductionAVLReliabilityClaimed}}</td></tr>
<tr><th><code>production_grade_eta_quality_claimed</code></th><td>{{.Help.ClaimFlags.ProductionGradeETAQualityClaimed}}</td></tr>
<tr><th><code>sla_claimed</code></th><td>{{.Help.ClaimFlags.SLAClaimed}}</td></tr>
<tr><th><code>uptime_guarantee_claimed</code></th><td>{{.Help.ClaimFlags.UptimeGuaranteeClaimed}}</td></tr>
<tr><th><code>dynamic_backend_plugin_loading_enabled</code></th><td>{{.Help.ClaimFlags.DynamicBackendPluginLoadingEnabled}}</td></tr>
</tbody></table>
<p class="muted">Help is private guidance. Stronger outside statements require separate retained authorization and source-specific proof.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "firstRunPanel"}}
<section class="hero" aria-labelledby="first-run-heading">
<h2 id="first-run-heading">Start Here</h2>
<p>{{.Boundary}}</p>
<p class="muted">{{.LocalDemoDeploymentEvidenceBoundary}}</p>
<p><strong>Task status:</strong> {{.Counts.Tasks}} tasks · ok {{index .Counts.Statuses "ok"}} · needs review {{index .Counts.Statuses "needs_review"}} · missing {{index .Counts.Statuses "missing"}} · blocked {{index .Counts.Statuses "blocked"}} · unknown {{index .Counts.Statuses "unknown"}}</p>
<div class="card-grid" aria-label="First-run evaluator paths">
{{range .Paths}}<section class="card" id="first-run-path-{{.ID}}">
<h3>{{.Label}}</h3>
<p><strong>Current signal:</strong> {{.CurrentSignal}}</p>
<p><strong>What it means:</strong> {{.Meaning}}</p>
<p><strong>First action:</strong> {{.FirstAction}}</p>
<p><strong>Console:</strong> <a href="{{.UILink}}">{{.UILink}}</a></p>
<p><strong>Docs:</strong> <code>{{.DocsLink}}</code></p>
<p><strong>Does not prove:</strong> {{.DoesNotProve}}</p>
</section>{{end}}
</div>
<h3>First-Run Acceptance Tasks</h3>
<table><thead><tr><th>Order</th><th>Task</th><th>Status</th><th>Current signal</th><th>What it means</th><th>Next action</th><th>Console</th><th>Docs</th><th>Does not prove</th></tr></thead><tbody>
{{range .Tasks}}<tr><td>{{.Order}}</td><td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td>{{.Meaning}}</td><td>{{.NextAction}}</td><td><a href="{{.UILink}}">{{.UILink}}</a></td><td><code>{{.DocsLink}}</code></td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Copy These Five Feed URLs</h3>
<table><thead><tr><th>Feed</th><th>Status</th><th>Copy value</th><th>Current link</th><th>Meaning</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .FeedURLs}}<tr><td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td><td>{{.Status}}</td><td><code class="copy-value">{{.CopyValue}}</code></td><td>{{if .URL}}<a href="{{.URL}}">{{.URL}}</a>{{else}}missing{{end}}</td><td>{{.Meaning}}</td><td>{{.NextAction}}<br><code>{{.DocsLink}}</code></td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<details>
<summary>Claim flags for this first-run guide</summary>
<table><tbody>
<tr><th><code>backend_command_execution_enabled</code></th><td>{{.ClaimFlags.BackendCommandExecutionEnabled}}</td></tr>
<tr><th><code>cache_diagnostics_read</code></th><td>{{.ClaimFlags.CacheDiagnosticsRead}}</td></tr>
<tr><th><code>external_network_contacted</code></th><td>{{.ClaimFlags.ExternalNetworkContacted}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>final_root_evidence_created</code></th><td>{{.ClaimFlags.FinalRootEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>secrets_collected</code></th><td>{{.ClaimFlags.SecretsCollected}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>agency_approval_claimed</code></th><td>{{.ClaimFlags.AgencyApprovalClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th><code>public_launch_claimed</code></th><td>{{.ClaimFlags.PublicLaunchClaimed}}</td></tr>
<tr><th><code>hosted_saas_claimed</code></th><td>{{.ClaimFlags.HostedSaaSClaimed}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>hardware_certification_claimed</code></th><td>{{.ClaimFlags.HardwareCertificationClaimed}}</td></tr>
<tr><th><code>production_avl_reliability_claimed</code></th><td>{{.ClaimFlags.ProductionAVLReliabilityClaimed}}</td></tr>
<tr><th><code>production_grade_eta_quality_claimed</code></th><td>{{.ClaimFlags.ProductionGradeETAQualityClaimed}}</td></tr>
<tr><th><code>sla_claimed</code></th><td>{{.ClaimFlags.SLAClaimed}}</td></tr>
<tr><th><code>uptime_guarantee_claimed</code></th><td>{{.ClaimFlags.UptimeGuaranteeClaimed}}</td></tr>
<tr><th><code>dynamic_backend_plugin_loading_enabled</code></th><td>{{.ClaimFlags.DynamicBackendPluginLoadingEnabled}}</td></tr>
<tr><th><code>release_candidate_approval_claimed</code></th><td>{{.ClaimFlags.ReleaseCandidateApprovalClaimed}}</td></tr>
<tr><th><code>managed_support_commitment_claimed</code></th><td>{{.ClaimFlags.ManagedSupportCommitmentClaimed}}</td></tr>
<tr><th><code>final_deployment_ownership_claimed</code></th><td>{{.ClaimFlags.FinalDeploymentOwnershipClaimed}}</td></tr>
<tr><th><code>consumer_ingestion_workflow_completed</code></th><td>{{.ClaimFlags.ConsumerIngestionWorkflowCompleted}}</td></tr>
<tr><th><code>production_multi_tenant_hosting_claimed</code></th><td>{{.ClaimFlags.ProductionMultiTenantHostingClaimed}}</td></tr>
</tbody></table>
</details>
</section>
{{end}}

{{define "dashboard"}}
{{template "layoutStart" .}}
<div class="hero">
<h2>Agency Operations Cockpit</h2>
<p>{{.Cockpit.Boundary}}</p>
<p><a href="/admin/operations.json">Export private cockpit JSON</a> · <a href="/admin/operations/maintenance">Open maintenance center</a></p>
</div>
{{template "firstRunPanel" .FirstRun}}
<h2>Setup Progress</h2>
<table><thead><tr><th>ID</th><th>Area</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .Cockpit.SetupProgress}}<tr id="cockpit-progress-{{.ID}}"><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td><span class="pill">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td><a href="{{.AdminLink}}">{{.NextAction}}</a></td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h2>Primary Actions</h2>
<div class="card-grid" aria-label="Primary agency operations actions">
{{range .Cockpit.PrimaryCards}}<section class="card" id="cockpit-card-{{.ID}}">
<h3>{{.Label}}</h3>
<p class="status">{{.Status}}</p>
<p><strong>Current signal:</strong> {{.CurrentSignal}}</p>
<p><strong>What should I do next?</strong> <a href="{{.AdminLink}}">{{.NextAction}}</a></p>
<p><strong>Does not prove:</strong> {{.DoesNotProve}}</p>
{{if .DocsLinks}}<p>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</p>{{end}}
</section>{{end}}
</div>
<h2>Readiness</h2>
{{if .DiscoveryError}}<p class="warning">{{.DiscoveryError}}. Next action: bootstrap publication metadata after a feed is available.</p>{{else}}
<p>Active GTFS feed version: {{if .ActiveFeedVersion}}<strong>{{.ActiveFeedVersion}}</strong>{{else}}not available{{end}}</p>
<table><tbody>
<tr><th>Public URLs</th><td>{{if .Discovery.Readiness.AllRequiredFeedsListed}}listed{{else}}missing or incomplete{{end}}</td></tr>
<tr><th>License</th><td>{{if .Discovery.Readiness.LicenseComplete}}complete{{else}}missing{{end}}</td></tr>
<tr><th>Contact</th><td>{{if .Discovery.Readiness.ContactComplete}}complete{{else}}missing{{end}}</td></tr>
<tr><th>HTTPS URLs</th><td>{{if .Discovery.Readiness.HTTPSURLs}}yes{{else}}not all HTTPS; local/dev URLs may be HTTP{{end}}</td></tr>
<tr><th>Canonical validation</th><td>{{if .Discovery.Readiness.CanonicalValidationComplete}}current passed/warning records exist{{else}}not complete{{end}}</td></tr>
</tbody></table>{{end}}

<h2>Dashboard Sections</h2>
<table><thead><tr><th>Section</th><th>Status</th><th>Last updated</th><th>Next action</th></tr></thead><tbody>
<tr><td>Private agency launchpad</td><td>{{len .Launchpad.Sections}} workflow sections</td><td>{{formatTime .Launchpad.GeneratedAt}}</td><td><a href="/admin/operations/launchpad">open launchpad</a> · <a href="/admin/operations/launchpad.json">export JSON</a></td></tr>
<tr><td>Setup wizard</td><td>{{len .SetupWizard.Stages}} staged setup rows</td><td>{{formatTime .SetupWizard.GeneratedAt}}</td><td><a href="/admin/operations/setup-wizard">open wizard</a> · <a href="/admin/operations/setup-wizard.json">export JSON</a></td></tr>
<tr><td>Connector Hub</td><td>{{len .ConnectorHub.Categories}} connector categories</td><td>{{formatTime .ConnectorHub.GeneratedAt}}</td><td><a href="/admin/operations/connectors">review connector paths</a> · <a href="/admin/operations/connectors.json">export JSON</a></td></tr>
<tr><td>Browser GTFS import</td><td>admin-only ZIP upload or URL import</td><td>{{formatTime .GeneratedAt}}</td><td><a href="/admin/operations/gtfs-import">import GTFS with validation feedback</a></td></tr>
<tr><td>Feed health dashboard</td><td>{{len .FeedHealth.Rows}} plain-language rows</td><td>{{formatTime .FeedHealth.GeneratedAt}}</td><td><a href="/admin/operations/feed-health">open feed health</a> · <a href="/admin/operations/feed-health.json">export JSON</a></td></tr>
<tr><td>Private operator checklist</td><td>{{len .Checklist.Groups}} grouped diagnostics</td><td>{{formatTime .GeneratedAt}}</td><td><a href="/admin/operations/checklist">open checklist</a> · <a href="/admin/operations/checklist.json">export JSON</a></td></tr>
<tr><td>Feeds / validation</td><td>{{if .DiscoveryError}}not configured{{else}}{{len .Discovery.Feeds}} feed records{{end}}</td><td>{{formatTimePtr .FeedsUpdatedAt}}</td><td><a href="/admin/operations/feeds">review feed URLs and validation</a></td></tr>
<tr><td>GTFS quality triage</td><td>{{.GTFSQuality.Canonical.Status}} static validator; {{.GTFSQuality.InternalImporter.Status}} internal importer</td><td>{{formatTimePtr .GTFSQuality.Canonical.ValidationTimestamp}}</td><td><a href="/admin/operations/gtfs-quality">review GTFS validator notices and operator actions</a></td></tr>
<tr><td>Validator health</td><td>{{.ValidationHealth.OverallStatus}} overall; tooling {{.ValidationHealth.ToolingStatus}}</td><td>{{formatTime .ValidationHealth.GeneratedAt}}</td><td><a href="/admin/operations/validation-health">review private validator diagnostics</a> · <a href="/admin/operations/validation-health.json">JSON</a></td></tr>
<tr><td>Operations reliability</td><td>{{.Reliability.OverallStatus}} overall</td><td>{{formatTime .Reliability.GeneratedAt}}</td><td><a href="/admin/operations/reliability">review private reliability diagnostics</a> · <a href="/admin/operations/reliability.json">JSON</a></td></tr>
<tr><td>Maintenance center</td><td>{{.Maintenance.OverallStatus}} overall</td><td>{{formatTime .Maintenance.GeneratedAt}}</td><td><a href="/admin/operations/maintenance">review maintenance tasks</a> · <a href="/admin/operations/maintenance.json">JSON</a></td></tr>
<tr><td>CAL-ITP-style readiness workflow</td><td>{{len .ReadinessV2.Rows}} checklist v2 rows</td><td>{{formatTime .ReadinessV2.GeneratedAt}}</td><td><a href="/admin/operations/readiness">review readiness gaps and next actions</a> · <a href="/admin/operations/readiness.json">export JSON</a></td></tr>
<tr><td>Telemetry freshness</td><td>{{if .TelemetryError}}{{.TelemetryError}}{{else}}{{len .Telemetry}} vehicles; {{.StaleCount}} stale{{end}}</td><td>{{formatTimePtr .TelemetryUpdatedAt}}</td><td><a href="/admin/operations/telemetry">inspect vehicle freshness</a></td></tr>
<tr><td>Telemetry simulator guide</td><td>{{if .TelemetrySimulator.LoadError}}{{.TelemetrySimulator.LoadError}}{{else}}{{len .TelemetrySimulator.Scenarios}} synthetic scenarios{{end}}</td><td>{{formatTime .TelemetrySimulator.GeneratedAt}}</td><td><a href="/admin/operations/telemetry-simulator">review simulator commands</a> · <a href="/admin/operations/telemetry-simulator.json">export JSON</a></td></tr>
<tr><td>Trip Updates quality</td><td>{{if .TripUpdatesQuality.Recorded}}{{.TripUpdatesQuality.DiagnosticsStatus}} / {{.TripUpdatesQuality.DiagnosticsReason}}{{else}}{{.TripUpdatesQuality.Message}}{{end}}</td><td>{{formatTimePtr .TripUpdatesQuality.SnapshotAt}}</td><td><a href="/admin/operations/feeds">review realtime quality summary</a></td></tr>
<tr><td>Scorecard</td><td>{{if .Scorecard}}{{.Scorecard.OverallStatus}}{{else}}{{.ScorecardError}}{{end}}</td><td>{{formatTimePtr .ScorecardUpdatedAt}}</td><td><a href="/admin/operations/evidence">find scorecard evidence</a></td></tr>
<tr><td>Consumer status</td><td>{{if .ConsumerError}}{{.ConsumerError}}{{else}}{{len .Consumers}} targets shown{{end}}</td><td>{{formatTimePtr .ConsumersUpdatedAt}}</td><td><a href="/admin/operations/consumers">review evidence-only statuses</a></td></tr>
<tr><td>Evidence links</td><td>repo documentation links</td><td>{{.EvidenceUpdatedAt}}</td><td><a href="/admin/operations/evidence">open evidence index</a></td></tr>
</tbody></table>

<h2>Public Feed URLs</h2>
{{if .DiscoveryError}}<p>No public feed metadata is available yet.</p>{{else}}{{template "feedTable" .}}{{end}}
{{template "tripUpdatesQuality" .}}
<p class="muted">Validation and public fetch records are supporting evidence only. They are not consumer acceptance or CAL-ITP/Caltrans compliance by themselves.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "launchpad"}}
{{template "layoutStart" .}}
<h2>Private Agency Launchpad</h2>
<p class="warning">This launchpad is private operator diagnostics. It creates no evidence, contacts no external party, changes no consumer status, and records no approval, compliance, public launch, hosted service, vendor, SLA, or production-grade ETA claim.</p>
<p><a href="/admin/operations/launchpad.json">Export private launchpad JSON</a> · <a href="/admin/operations/checklist">Open private checklist</a> · <a href="/admin/operations/readiness">Open readiness review</a></p>
{{template "firstRunPanel" .FirstRun}}
<table><tbody>
<tr><th>Boundary</th><td>{{.Launchpad.Boundary}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.Launchpad.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>final_root_evidence_created</code></th><td>{{.Launchpad.ClaimFlags.FinalRootEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.Launchpad.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.Launchpad.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.Launchpad.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>agency_approval_claimed</code></th><td>{{.Launchpad.ClaimFlags.AgencyApprovalClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.Launchpad.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th><code>public_launch_claimed</code></th><td>{{.Launchpad.ClaimFlags.PublicLaunchClaimed}}</td></tr>
<tr><th><code>hosted_saas_claimed</code></th><td>{{.Launchpad.ClaimFlags.HostedSaaSClaimed}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.Launchpad.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>production_grade_eta_claimed</code></th><td>{{.Launchpad.ClaimFlags.ProductionGradeETAClaimed}}</td></tr>
</tbody></table>

<h3>Workflow Sections</h3>
<table><thead><tr><th>ID</th><th>Section</th><th>Status</th><th>Current signal</th><th>Next actions</th><th>Commands</th><th>Links</th><th>Docs</th><th>Boundary</th></tr></thead><tbody>
{{range .Launchpad.Sections}}<tr><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td><ul>{{range .NextActions}}<li>{{.}}</li>{{end}}</ul></td><td>{{range .CommandSuggestions}}<code>{{.}}</code><br>{{end}}</td><td>{{range .AdminLinks}}<a href="{{.}}">{{.}}</a><br>{{end}}</td><td>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>

<h3>Decision Gate</h3>
<table><thead><tr><th>Option</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .Launchpad.DecisionNotes}}<tr><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.Boundary}}</td></tr>{{end}}
</tbody></table>
<p class="muted">No POST action exists for this page. Missing data remains missing or unknown until the underlying private source records change.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "setup-wizard"}}
{{template "layoutStart" .}}
<h2>Setup Wizard</h2>
<p class="warning">{{.SetupWizard.Boundary}}</p>
<p><a href="/admin/operations/setup-wizard.json">Export private setup wizard JSON</a> · <a href="/admin/operations/setup">Open guided setup</a> · <a href="/admin/operations/checklist">Open private checklist</a></p>
<table><tbody>
<tr><th><code>external_evidence_created</code></th><td>{{.SetupWizard.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>final_root_evidence_created</code></th><td>{{.SetupWizard.ClaimFlags.FinalRootEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.SetupWizard.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.SetupWizard.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.SetupWizard.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>agency_approval_claimed</code></th><td>{{.SetupWizard.ClaimFlags.AgencyApprovalClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.SetupWizard.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th><code>public_launch_claimed</code></th><td>{{.SetupWizard.ClaimFlags.PublicLaunchClaimed}}</td></tr>
<tr><th><code>hosted_saas_claimed</code></th><td>{{.SetupWizard.ClaimFlags.HostedSaaSClaimed}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.SetupWizard.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>production_grade_eta_claimed</code></th><td>{{.SetupWizard.ClaimFlags.ProductionGradeETAClaimed}}</td></tr>
</tbody></table>
<table><thead><tr><th>ID</th><th>Stage</th><th>Status</th><th>Current signal</th><th>Primary action</th><th>Console</th><th>Docs</th><th>Boundary</th></tr></thead><tbody>
{{range .SetupWizard.Stages}}<tr><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td>{{.PrimaryAction}}</td><td>{{if .AdminLink}}<a href="{{.AdminLink}}">{{.AdminLink}}</a>{{end}}</td><td>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
<p class="muted">This wizard is GET-only. It does not upload GTFS, mutate setup state, run validators, contact external systems, or create public routes.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "connectors"}}
{{template "layoutStart" .}}
<h2>Connector Hub</h2>
<p class="warning">{{.ConnectorHub.Boundary}}</p>
<table><tbody>
<tr><th>Safe plugin definition</th><td>{{.ConnectorHub.PluginDefinition}}</td></tr>
<tr><th><code>dynamic_backend_plugin_loading_enabled</code></th><td>{{.ConnectorHub.ClaimFlags.DynamicBackendPluginLoadingEnabled}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.ConnectorHub.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.ConnectorHub.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.ConnectorHub.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.ConnectorHub.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.ConnectorHub.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
</tbody></table>
<div class="card-grid" aria-label="Connector categories">
{{range .ConnectorHub.Categories}}
<section class="card">
<h3>{{.Label}}</h3>
<p class="status">Status: {{.Status}}</p>
<p>{{.Summary}}</p>
<p><strong>Connector shape:</strong> {{.ConnectorShape}}</p>
<p><strong>Inputs:</strong> {{join .Inputs ", "}}</p>
<p><strong>Outputs:</strong> {{join .Outputs ", "}}</p>
<p><strong>Failure behavior:</strong> {{.FailureBehavior}}</p>
<p><strong>Boundary:</strong> {{.ClaimBoundary}}</p>
{{if .AdminLinks}}<p><strong>Console:</strong> {{range .AdminLinks}}<a href="{{.}}">{{.}}</a> {{end}}</p>{{end}}
{{if .DocsLinks}}<p><strong>Docs:</strong> {{range .DocsLinks}}<code>{{.}}</code> {{end}}</p>{{end}}
{{if .CommandSuggestions}}<p><strong>Checks:</strong> {{range .CommandSuggestions}}<code>{{.}}</code> {{end}}</p>{{end}}
</section>
{{end}}
</div>
<h3>Manifest Registry</h3>
<p>Read-only registry of committed synthetic connector example manifests. It does not accept uploads, load backend plugins, execute manifest commands, contact external systems, create retained evidence, or change consumer status.</p>
{{if .ConnectorHub.Registry.Diagnostics}}
<table><thead><tr><th>Level</th><th>Code</th><th>Path</th><th>Message</th></tr></thead><tbody>
{{range .ConnectorHub.Registry.Diagnostics}}<tr><td>{{.Level}}</td><td><code>{{.Code}}</code></td><td><code>{{.Path}}</code></td><td>{{.Message}}</td></tr>{{end}}
</tbody></table>
{{end}}
<table><thead><tr><th>Connector</th><th>Type</th><th>Mode</th><th>Contracts</th><th>Failure / redaction</th><th>Conformance</th><th>Boundary</th><th>Docs</th></tr></thead><tbody>
{{range .ConnectorHub.Registry.Entries}}
<tr>
<td><strong>{{.DisplayName}}</strong><br><code>{{.ConnectorID}}</code><br><code>{{.SourcePath}}</code></td>
<td>{{.ConnectorType}}</td>
<td>{{.ModeName}}<br>{{if .DisabledByDefault}}disabled by default{{else}}review before use{{end}}</td>
<td>inputs: {{range .InputContracts}}{{.Name}} (<code>{{.Schema}}</code>)<br>{{end}}outputs: {{range .OutputContracts}}{{.Name}} (<code>{{.Schema}}</code>)<br>{{end}}</td>
<td>{{if .FailureBehavior.FailClosed}}fail closed{{else}}review failure behavior{{end}}; {{.FailureBehavior.DegradedState}}<br>secret storage: {{.RedactionPolicy.SecretStorage}}</td>
<td>{{len .ConformanceCases}} synthetic cases</td>
<td>claims: {{join .ClaimBoundary.PositiveClaims ", "}}<br>not claimed: {{join .ClaimBoundary.NotClaimed ", "}}</td>
<td><code>{{.DocsLink}}</code></td>
</tr>
{{else}}
<tr><td colspan="8">No committed connector manifests were loaded. Review diagnostics and run <code>make external-connection-check</code>.</td></tr>
{{end}}
</tbody></table>
<p><a href="/admin/operations/connectors/tests">Open connector test instructions</a> for fixed offline checks.</p>
<p class="muted">Connector Hub is read-only. It exposes safe integration paths and local checks; it does not run external systems, collect retained evidence, contact vendors or consumers, or change consumer status.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "connector-tests"}}
{{template "layoutStart" .}}
<h2>Connector Test Instructions</h2>
<p class="warning">{{.ConnectorTests.Boundary}}</p>
<table><tbody>
<tr><th><code>backend_command_execution_enabled</code></th><td>{{.ConnectorTests.ClaimFlags.BackendCommandExecutionEnabled}}</td></tr>
<tr><th><code>manifest_command_execution_enabled</code></th><td>{{.ConnectorTests.ClaimFlags.ManifestCommandExecutionEnabled}}</td></tr>
<tr><th><code>external_network_contacted</code></th><td>{{.ConnectorTests.ClaimFlags.ExternalNetworkContacted}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.ConnectorTests.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.ConnectorTests.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.ConnectorTests.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
</tbody></table>
<table><thead><tr><th>Check</th><th>Copyable instruction</th><th>What it validates</th><th>Inputs</th><th>Failure next action</th><th>Does not prove</th><th>Docs</th></tr></thead><tbody>
{{range .ConnectorTests.Commands}}
<tr>
<td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td>
<td><code>{{.CommandLine}}</code></td>
<td>{{.Validates}}</td>
<td>{{.Inputs}}</td>
<td>{{.FailureNextAction}}</td>
<td>{{.DoesNotProve}}</td>
<td>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</td>
</tr>
{{end}}
</tbody></table>
<p class="muted">This page is GET-only generated guidance. It does not execute commands, read manifest-provided commands, run validators, start sidecars, write files, contact external parties, or change consumer status.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "gtfs-import"}}
{{template "layoutStart" .}}
<h2>Browser GTFS Import</h2>
<p class="warning">Private admin-only import path. Raw GTFS ZIP bytes are written to temporary runtime storage for the import attempt and then removed. This page creates no retained evidence, contacts no consumers, records no agency approval, and makes no CAL-ITP/Caltrans compliance, public launch, hosted-service, vendor compatibility, production-readiness, or production-grade ETA claim.</p>
{{if .GTFSImportNotice}}<p class="ok">{{.GTFSImportNotice}}</p>{{end}}
{{if .GTFSImportError}}<p class="bad">{{.GTFSImportError}}</p>{{end}}
<h3>Current Active Schedule</h3>
<table><tbody>
<tr><th>Active feed version</th><td>{{if .ActiveFeedVersion}}<code>{{.ActiveFeedVersion}}</code>{{else}}missing active schedule{{end}}</td></tr>
<tr><th>Schedule source review</th><td>Use this page to compare a new import's row counts and validation blockers with the current active schedule. Full staged diff and browser rollback execution are not available yet.</td></tr>
<tr><th>Rollback visibility</th><td>The active feed version is visible here. Prior feed-version listing and rollback execution remain technical-helper workflows until a safe browser rollback view is implemented.</td></tr>
</tbody></table>
{{if .GTFSImportSource.SourceType}}
<h3>GTFS Source Review</h3>
<table><tbody>
<tr><th>Source type</th><td>{{.GTFSImportSource.SourceType}}</td></tr>
<tr><th>Source URL</th><td>{{.GTFSImportSource.SourceURL}}</td></tr>
<tr><th>Checksum SHA-256</th><td>{{.GTFSImportSource.ChecksumSHA256}}</td></tr>
<tr><th>Byte count</th><td>{{.GTFSImportSource.ByteCount}}</td></tr>
<tr><th>Import timestamp</th><td>{{formatTimePtr .GTFSImportSource.ImportTimestamp}}</td></tr>
<tr><th>Active feed version after import</th><td>{{if .GTFSImportSource.ActiveFeedVersion}}<code>{{.GTFSImportSource.ActiveFeedVersion}}</code>{{else}}not available{{end}}</td></tr>
<tr><th>Schedule identity summary</th><td>{{.GTFSImportSource.ScheduleIdentitySummary}}</td></tr>
<tr><th>Update comparison</th><td>{{.GTFSImportSource.UpdateComparison}}</td></tr>
<tr><th>Rollback visibility</th><td>{{.GTFSImportSource.RollbackVisibility}}</td></tr>
</tbody></table>
{{end}}
{{if .GTFSImportResult}}
<h3>Last Import Attempt From This Page</h3>
<table><tbody>
<tr><th>Status</th><td>{{.GTFSImportResult.Status}}</td></tr>
<tr><th>Import ID</th><td>{{.GTFSImportResult.ImportID}}</td></tr>
<tr><th>Feed version</th><td>{{if .GTFSImportResult.FeedVersionID}}{{.GTFSImportResult.FeedVersionID}}{{else}}not published{{end}}</td></tr>
<tr><th>Errors</th><td>{{.GTFSImportResult.ErrorCount}}</td></tr>
<tr><th>Warnings</th><td>{{.GTFSImportResult.WarningCount}}</td></tr>
<tr><th>Info</th><td>{{.GTFSImportResult.InfoCount}}</td></tr>
<tr><th>Validation report stored</th><td>{{.GTFSImportResult.ReportStored}}</td></tr>
{{if .GTFSImportResult.FailureMessage}}<tr><th>Failure message</th><td>{{.GTFSImportResult.FailureMessage}}</td></tr>{{end}}
</tbody></table>
{{if .GTFSImportResult.Counts}}<h3>Import Counts</h3><table><thead><tr><th>GTFS file</th><th>Rows</th></tr></thead><tbody>{{range .GTFSImportResult.Counts}}<tr><td>{{.Label}}</td><td>{{.Count}}</td></tr>{{end}}</tbody></table>{{end}}
{{end}}
{{if .IsAdmin}}
<div class="card-grid" aria-label="GTFS import options">
<section class="card">
<h3>Upload ZIP</h3>
<form method="post" enctype="multipart/form-data" action="/admin/operations/gtfs-import?csrf_token={{.CSRFToken}}">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="action" value="import_gtfs">
<input type="hidden" name="source_type" value="upload">
<label for="gtfs_upload_zip">GTFS ZIP</label><input id="gtfs_upload_zip" type="file" name="gtfs_zip" accept=".zip,application/zip,application/octet-stream" required>
<label for="gtfs_upload_notes">Notes</label><textarea id="gtfs_upload_notes" name="notes" maxlength="500" rows="3" placeholder="Optional operator note without credentials or private paths"></textarea>
<button type="submit">Import ZIP</button>
</form>
</section>
<section class="card">
<h3>Import From URL</h3>
<form method="post" action="/admin/operations/gtfs-import">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="action" value="import_gtfs">
<input type="hidden" name="source_type" value="url">
<label for="gtfs_import_url">GTFS ZIP URL</label><input id="gtfs_import_url" type="url" name="gtfs_url" maxlength="2048" placeholder="https://agency.example/gtfs.zip" required>
<label for="gtfs_url_notes">Notes</label><textarea id="gtfs_url_notes" name="notes" maxlength="500" rows="3" placeholder="Optional operator note without credentials or private paths"></textarea>
<button type="submit">Import URL</button>
</form>
</section>
</div>
{{else}}
<p class="warning">Import actions require an admin role. Read-only, operator, and editor roles can review this page and use the linked quality/health views, but cannot upload or import GTFS from the browser.</p>
{{end}}
<h3>Next Actions</h3>
<table><tbody>
<tr><th>Validation feedback</th><td><a href="/admin/operations/gtfs-quality">Review GTFS quality triage</a> and stored import validation messages before relying on the feed.</td></tr>
<tr><th>Feed health</th><td><a href="/admin/operations/feed-health">Review the five-path feed health command center</a> after a successful publish.</td></tr>
<tr><th>Validator health</th><td><a href="/admin/operations/validation-health">Run or review allowlisted validator health</a>; browser requests cannot supply validator commands, paths, URLs, argument lists, binaries, artifacts, or timeouts.</td></tr>
<tr><th>Update decision</th><td>Compare current active schedule, new import counts, warnings, and blockers. If staged comparison is required, use a technical helper until browser staging comparison is implemented.</td></tr>
<tr><th>Rollback</th><td>Use the active feed version shown above and the operator rollback documentation. This page does not fake a rollback button.</td></tr>
<tr><th>Typed edits</th><td><a href="/admin/gtfs-studio">Open GTFS Studio</a> when an agency needs draft authoring instead of ZIP import.</td></tr>
<tr><th>CLI fallback</th><td>Keep using the documented CLI import path for large files, scripted imports, or runtimes where browser import is unavailable.</td></tr>
</tbody></table>
<p class="muted">Browser import accepts a ZIP upload or a safe HTTP(S) URL, runs the existing importer, and stores only normal import and validation records. Private/local URLs are blocked unless the runtime explicitly enables local testing overrides. After import, review GTFS quality, validator health, and all five public feed paths before treating the dataset as ready for wider operator review.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "readiness"}}
{{template "layoutStart" .}}
<h2>Readiness Checklist V2</h2>
<p class="warning">{{.ReadinessV2.Boundary}}</p>
<p><a href="/admin/operations/checklist">Open private operator checklist</a> · <a href="/admin/operations/checklist.json">Export private checklist JSON</a></p>
<p><a href="/admin/operations/feed-health">Open plain-language feed health</a> · <a href="/admin/operations/feed-health.json">Export private feed health JSON</a></p>
<p><a href="/admin/operations/readiness.json">Export private readiness v2 JSON</a> · <a href="/admin/operations/gtfs-quality">Open authenticated GTFS quality triage</a> · <a href="/admin/operations/validation-health">Open private validator health diagnostics</a></p>
<p class="muted">Each Readiness item card explains the current private signal, why it matters, the next operator action, and the boundary around stronger claims. Consumer tracker states remain prepared unless retained target-originated evidence supports a target-specific change. Claim flags are available in the private JSON export and remain false.</p>
<div class="card-grid" aria-label="Readiness checklist v2 rows">
{{range .ReadinessV2.Rows}}
<section class="card">
<h3>{{.ReadinessItem}}</h3>
<p class="status">Status <span class="pill">{{.Status}}</span></p>
<p><strong>Current signal:</strong> {{.CurrentSignal}}</p>
<p><strong>What this means:</strong> {{.WhatThisMeans}}</p>
<p><strong>Why it matters:</strong> {{.WhyItMatters}}</p>
<p><strong>What to do next:</strong> {{.WhatToDoNext}}</p>
<p><strong>What it does not prove:</strong> {{.WhatItDoesNotProve}}</p>
<p><strong>Console:</strong> {{range .AdminLinks}}<a href="{{.}}">{{.}}</a> {{end}}</p>
<p><strong>Docs:</strong> {{range .DocsLinks}}<code>{{.}}</code> {{end}}</p>
</section>
{{end}}
</div>
<p class="muted">No external evidence is created by viewing this page, and this workflow does not contact consumers, validators, agency systems, or external portals.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "feed-health"}}
{{template "layoutStart" .}}
<h2>Feed Health Dashboard</h2>
<p class="warning">{{.FeedHealth.Boundary}}</p>
<p>This command center tracks exactly five public paths: <code>/public/feeds.json</code>, <code>/public/gtfs/schedule.zip</code>, <code>/public/gtfsrt/vehicle_positions.pb</code>, <code>/public/gtfsrt/trip_updates.pb</code>, and <code>/public/gtfsrt/alerts.pb</code>.</p>
<p><a href="/admin/operations/feed-health.json">Export private feed health JSON</a> · <a href="/admin/operations/validation-health">Open validator health</a> · <a href="/admin/operations/reliability">Open reliability diagnostics</a></p>
<table><tbody>
<tr><th><code>external_evidence_created</code></th><td>{{.FeedHealth.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.FeedHealth.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.FeedHealth.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.FeedHealth.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>sla_claimed</code></th><td>{{.FeedHealth.ClaimFlags.SLAClaimed}}</td></tr>
<tr><th><code>uptime_guarantee_claimed</code></th><td>{{.FeedHealth.ClaimFlags.UptimeGuaranteeClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.FeedHealth.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th><code>public_launch_claimed</code></th><td>{{.FeedHealth.ClaimFlags.PublicLaunchClaimed}}</td></tr>
</tbody></table>
<div class="card-grid" aria-label="Plain-language feed health rows">
{{range .FeedHealth.Rows}}
<section class="card">
<h3>{{.Label}}</h3>
{{if .PublicPath}}<p><strong>Public path:</strong> <code>{{.PublicPath}}</code></p>{{end}}
<p><strong>Configured URL:</strong> <code>{{.ConfiguredURL}}</code></p>
<p><strong>Last known HTTP status:</strong> {{.LastKnownHTTPStatus}}</p>
<p><strong>Byte count:</strong> {{.ByteCount}}</p>
<p><strong>Content type:</strong> {{.ContentType}}</p>
<p><strong>Checksum:</strong> {{.Checksum}}</p>
<p><strong>Last generated:</strong> {{.LastGenerated}}</p>
<p><strong>Last checked:</strong> {{.LastChecked}}</p>
<p><strong>Validator state:</strong> {{.ValidatorState}}</p>
<p><strong>Health state:</strong> {{.HealthState}}</p>
<p class="status">{{.StatusText}} <span class="pill">{{.Status}}</span></p>
<p><strong>Current signal:</strong> {{.CurrentSignal}}</p>
<p><strong>What this means:</strong> {{.WhatThisMeans}}</p>
<p><strong>Freshness:</strong> {{.Freshness}}</p>
<p><strong>Validator context:</strong> {{.ValidatorContext}}</p>
<p><strong>Health context:</strong> {{.HealthContext}}</p>
<p><strong>Next action:</strong> {{.NextAction}}</p>
<p><strong>Does not prove:</strong> {{.DoesNotProve}}</p>
{{if .AdminLinks}}<p><strong>Console:</strong> {{range .AdminLinks}}<a href="{{.}}">{{.}}</a> {{end}}</p>{{end}}
{{if .DocsLinks}}<p><strong>Docs:</strong> {{range .DocsLinks}}<code>{{.}}</code> {{end}}</p>{{end}}
</section>
{{end}}
</div>
<h3>Realtime Usefulness</h3>
<div class="card-grid" aria-label="Realtime usefulness rows">
<section class="card" id="realtime-usefulness-vehicle-positions"><h3>{{.FeedHealth.RealtimeUsefulness.VehiclePositions.Label}}</h3><p class="status">{{.FeedHealth.RealtimeUsefulness.VehiclePositions.State}}</p><p><strong>Count:</strong> {{.FeedHealth.RealtimeUsefulness.VehiclePositions.Count}}</p><p><strong>Latest signal:</strong> {{.FeedHealth.RealtimeUsefulness.VehiclePositions.LatestSignal}}</p><p><strong>Stale or suppressed:</strong> {{.FeedHealth.RealtimeUsefulness.VehiclePositions.StaleOrHeld}}</p><p><strong>Next action:</strong> <a href="{{.FeedHealth.RealtimeUsefulness.VehiclePositions.AdminLink}}">{{.FeedHealth.RealtimeUsefulness.VehiclePositions.NextAction}}</a></p><p><strong>Does not prove:</strong> {{.FeedHealth.RealtimeUsefulness.VehiclePositions.DoesNotProve}}</p></section>
<section class="card" id="realtime-usefulness-trip-updates"><h3>{{.FeedHealth.RealtimeUsefulness.TripUpdates.Label}}</h3><p class="status">{{.FeedHealth.RealtimeUsefulness.TripUpdates.State}}</p><p><strong>Adapter:</strong> {{.FeedHealth.RealtimeUsefulness.TripUpdates.Adapter}}</p><p><strong>Count:</strong> {{.FeedHealth.RealtimeUsefulness.TripUpdates.Count}}</p><p><strong>Latest signal:</strong> {{.FeedHealth.RealtimeUsefulness.TripUpdates.LatestSignal}}</p><p><strong>Withheld / stale:</strong> {{.FeedHealth.RealtimeUsefulness.TripUpdates.StaleOrHeld}}</p>{{if .FeedHealth.RealtimeUsefulness.TripUpdates.Details}}<p><strong>Withheld reasons:</strong> {{range .FeedHealth.RealtimeUsefulness.TripUpdates.Details}}<span class="pill">{{.Label}}: {{.Count}}</span> {{end}}</p>{{end}}<p><strong>Next action:</strong> <a href="{{.FeedHealth.RealtimeUsefulness.TripUpdates.AdminLink}}">{{.FeedHealth.RealtimeUsefulness.TripUpdates.NextAction}}</a></p><p><strong>Does not prove:</strong> {{.FeedHealth.RealtimeUsefulness.TripUpdates.DoesNotProve}}</p></section>
<section class="card" id="realtime-usefulness-alerts"><h3>{{.FeedHealth.RealtimeUsefulness.Alerts.Label}}</h3><p class="status">{{.FeedHealth.RealtimeUsefulness.Alerts.State}}</p><p><strong>Count:</strong> {{.FeedHealth.RealtimeUsefulness.Alerts.Count}}</p><p><strong>Latest signal:</strong> {{.FeedHealth.RealtimeUsefulness.Alerts.LatestSignal}}</p><p><strong>Stale or held:</strong> {{.FeedHealth.RealtimeUsefulness.Alerts.StaleOrHeld}}</p><p><strong>Next action:</strong> <a href="{{.FeedHealth.RealtimeUsefulness.Alerts.AdminLink}}">{{.FeedHealth.RealtimeUsefulness.Alerts.NextAction}}</a></p><p><strong>Does not prove:</strong> {{.FeedHealth.RealtimeUsefulness.Alerts.DoesNotProve}}</p></section>
</div>
<p class="muted">This dashboard summarizes existing private records only. Missing records stay missing or unknown until the underlying source records change.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "feeds"}}
{{template "layoutStart" .}}
<h2>Feed URLs And Validation</h2>
{{if .DiscoveryError}}<p class="warning">No feed metadata is available. Next action: publish or import a GTFS feed, then bootstrap publication metadata.</p>{{else}}
{{template "feedTable" .}}
{{template "tripUpdatesQuality" .}}
<h3>Feed discovery document</h3>
<table><thead><tr><th>Item</th><th>URL</th><th>Validation</th><th>Last checked</th></tr></thead><tbody>
<tr><td>feeds.json</td><td>{{.Discovery.PublicBaseURL}}/public/feeds.json</td><td>not a validator result</td><td>{{formatTime .Discovery.GeneratedAt}}</td></tr>
</tbody></table>
{{end}}
<p class="muted">This view shows repo/deployment evidence only. Third-party consumer acceptance requires retained confirmation from the named consumer.</p>
<p><a href="/admin/operations/feed-health">Open plain-language feed health</a> · <a href="/admin/operations/gtfs-quality">Review GTFS quality triage actions</a> · <a href="/admin/operations/validation-health">Review private validator health diagnostics</a></p>
{{template "layoutEnd" .}}
{{end}}

{{define "gtfs-quality"}}
{{template "layoutStart" .}}
<h2>GTFS Quality Triage</h2>
{{if .GTFSQualityNotice}}<p class="ok">{{.GTFSQualityNotice}}</p>{{end}}
{{if .GTFSQualityError}}<p class="bad">{{.GTFSQualityError}}</p>{{end}}
<p class="warning">Validator output is diagnostics and supporting signal only. It is not consumer acceptance, not CAL-ITP/Caltrans compliance, not an evidence packet, and not production-readiness proof.</p>
<table><tbody>
<tr><th>Active schedule feed version</th><td>{{if .ActiveFeedVersion}}<code>{{.ActiveFeedVersion}}</code>{{else}}missing active schedule; next action: import or publish a schedule before rerunning validation{{end}}</td></tr>
<tr><th>Rerun boundary</th><td>Rerun uses only the authenticated agency active published schedule ZIP and the server-side static MobilityData validator mapping.</td></tr>
<tr><th>Guidance boundary</th><td>{{.GTFSQualityGuidance.Boundary}}</td></tr>
</tbody></table>
<h3>Fix Workflow</h3>
<div class="card-grid">
{{range .GTFSQualityGuidance.Workflow}}<section class="card"><h3>{{.Label}}</h3><p>{{.Summary}}</p><p><strong>Next outcome:</strong> {{.NextOutcome}}</p><p class="muted">{{.DoesNotDo}}</p><p>{{range .AdminLinks}}<a href="{{.}}">{{.}}</a><br>{{end}}</p><p>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</p></section>{{end}}
</div>
<h3>Claim Flags</h3>
<table><tbody>
<tr><th><code>automatic_gtfs_edit_enabled</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.AutomaticGTFSEditEnabled}}</td></tr>
<tr><th><code>draft_mutation_enabled</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.DraftMutationEnabled}}</td></tr>
<tr><th><code>schedule_publish_enabled</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.SchedulePublishEnabled}}</td></tr>
<tr><th><code>validator_semantics_changed</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.ValidatorSemanticsChanged}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>agency_approval_claimed</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.AgencyApprovalClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th><code>public_launch_claimed</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.PublicLaunchClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>production_grade_eta_claimed</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.ProductionGradeETAClaimed}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>hardware_certification_claimed</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.HardwareCertificationClaimed}}</td></tr>
<tr><th><code>production_avl_reliability_claimed</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.ProductionAVLReliabilityClaimed}}</td></tr>
</tbody></table>
{{if .IsAdmin}}
<h3>Rerun Static Validator</h3>
<form method="post" action="/admin/operations/gtfs-quality">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="action" value="rerun_static_validator">
<button type="submit">Rerun static MobilityData validator</button>
</form>
{{else}}<p class="muted">Rerun action is available only to admins.</p>{{end}}
{{template "gtfsQualitySection" .GTFSQuality.Canonical}}
{{template "gtfsQualitySection" .GTFSQuality.InternalImporter}}
{{template "layoutEnd" .}}
{{end}}

{{define "validation-health"}}
{{template "layoutStart" .}}
<h2>Validator Health</h2>
{{if .ValidationHealthNotice}}<p class="ok">{{.ValidationHealthNotice}}</p>{{end}}
{{if .ValidationHealthError}}<p class="bad">{{.ValidationHealthError}}</p>{{end}}
<p class="warning">This page is private diagnostics only. It does not create evidence packets, contact consumers, change consumer statuses, claim CAL-ITP/Caltrans compliance, or claim production readiness.</p>
<div class="card-grid">
<section class="card"><h3>Internal import validation</h3><p>Open Transit RT importer checks required GTFS structure and blocks unsafe activation paths. It helps explain import failures, but it is not the canonical MobilityData validator.</p></section>
<section class="card"><h3>Canonical static validation</h3><p>MobilityData static GTFS validation reviews the active schedule artifact when pinned tooling is installed and the schedule artifact is available.</p></section>
<section class="card"><h3>GTFS-Realtime validation</h3><p>Realtime validation reviews server-owned Vehicle Positions, Trip Updates, and Alerts protobuf artifacts. Browser requests cannot supply commands, paths, argument lists, artifacts, validator binaries, URLs, or timeouts.</p></section>
</div>
<table><tbody>
<tr><th>Overall status</th><td>{{.ValidationHealth.OverallStatus}}</td></tr>
<tr><th>Tooling status</th><td>{{.ValidationHealth.ToolingStatus}}</td></tr>
<tr><th>Generated at</th><td>{{formatTime .ValidationHealth.GeneratedAt}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.ValidationHealth.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.ValidationHealth.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.ValidationHealth.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.ValidationHealth.ProductionReadinessClaimed}}</td></tr>
</tbody></table>
{{if .IsAdmin}}
<h3>Run All Validators</h3>
<form method="post" action="/admin/operations/validation-health">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="action" value="run_all">
<button type="submit">Run allowlisted validators</button>
</form>
{{else}}<p class="muted">Run action is available only to admins.</p>{{end}}
{{template "validationHealthRows" .ValidationHealth}}
<p class="muted">Static schedule health uses the canonical MobilityData static validator. Realtime health uses the MobilityData GTFS-Realtime validator for Vehicle Positions, Trip Updates, and Alerts. Open Transit RT internal GTFS import validation remains context in GTFS quality triage and is not canonical validator health.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "reliability"}}
{{template "layoutStart" .}}
<h2>Operations Reliability</h2>
{{if .ReliabilityError}}<p class="warning">{{.ReliabilityError}}.</p>{{end}}
<p class="warning">This page is private operations diagnostics only. It does not create evidence, change consumer statuses, claim compliance, claim production readiness, claim SLA coverage, guarantee uptime, claim hosted service availability, claim agency adoption, claim consumer acceptance, claim vendor compatibility, or claim production-grade ETA quality.</p>
<table><tbody>
<tr><th>Overall status</th><td>{{.Reliability.OverallStatus}}</td></tr>
<tr><th>Generated at</th><td>{{formatTime .Reliability.GeneratedAt}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.Reliability.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>final_root_evidence_created</code></th><td>{{.Reliability.ClaimFlags.FinalRootEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.Reliability.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.Reliability.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.Reliability.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>sla_claimed</code></th><td>{{.Reliability.ClaimFlags.SLAClaimed}}</td></tr>
<tr><th><code>uptime_guarantee_claimed</code></th><td>{{.Reliability.ClaimFlags.UptimeGuaranteeClaimed}}</td></tr>
<tr><th><code>hosted_saas_claimed</code></th><td>{{.Reliability.ClaimFlags.HostedSaaSClaimed}}</td></tr>
<tr><th><code>agency_adoption_claimed</code></th><td>{{.Reliability.ClaimFlags.AgencyAdoptionClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.Reliability.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.Reliability.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>production_grade_eta_claimed</code></th><td>{{.Reliability.ClaimFlags.ProductionGradeETAClaimed}}</td></tr>
</tbody></table>

<h3>Feeds</h3>
<table><thead><tr><th>Feed</th><th>Status</th><th>Source</th><th>Snapshot</th><th>Endpoint</th><th>Freshness seconds</th><th>Generation latency ms</th><th>Invalid response percent</th><th>Matched vehicle percent</th><th>Coverage percent</th><th>Threshold</th><th>Next action</th></tr></thead><tbody>
{{range .Reliability.Feeds}}<tr><td>{{.FeedType}}</td><td>{{.Status}}</td><td>{{.Source}}</td><td>{{formatTimePtr .SnapshotAt}}</td><td>{{formatBoolPtr .EndpointAvailable}}</td><td>{{formatFloatPtr .FreshnessSeconds}}</td><td>{{formatFloatPtr .GenerationLatencyMS}}</td><td>{{formatFloatPtr .InvalidResponsePercent}}</td><td>{{formatFloatPtr .MatchedVehiclePercent}}</td><td>{{formatFloatPtr .CoveragePercent}}</td><td>{{.DiagnosticThreshold}}</td><td>{{.NextAction}}</td></tr>{{end}}
</tbody></table>

<h3>Incidents</h3>
<table><tbody>
<tr><th>Status</th><td>{{.Reliability.Incidents.Status}}</td></tr>
<tr><th>Source</th><td>{{.Reliability.Incidents.Source}}</td></tr>
<tr><th>Total</th><td>{{.Reliability.Incidents.Total}}</td></tr>
<tr><th>Oldest open age seconds</th><td>{{formatInt64Ptr .Reliability.Incidents.OldestOpenAgeSeconds}}</td></tr>
<tr><th>Next action</th><td>{{.Reliability.Incidents.NextAction}}</td></tr>
</tbody></table>
<table><thead><tr><th>ID</th><th>Type</th><th>Severity</th><th>Status</th><th>Opened</th><th>Updated</th><th>Title</th><th>Category</th></tr></thead><tbody>
{{range .Reliability.Incidents.Recent}}<tr><td>{{.ID}}</td><td>{{.Type}}</td><td>{{.Severity}}</td><td>{{.Status}}</td><td>{{formatTime .OpenedAt}}</td><td>{{formatTimePtr .UpdatedAt}}</td><td>{{.Title}}</td><td>{{.Category}}</td></tr>{{end}}
</tbody></table>

<h3>Safe Source Sections</h3>
<table><thead><tr><th>Section</th><th>Status</th><th>Source</th><th>Summary</th><th>Next action</th></tr></thead><tbody>
<tr><td>backup_restore</td><td>{{.Reliability.BackupRestore.Status}}</td><td>{{.Reliability.BackupRestore.Source}}</td><td>{{.Reliability.BackupRestore.Summary}}</td><td>{{.Reliability.BackupRestore.NextAction}}</td></tr>
<tr><td>alerting</td><td>{{.Reliability.Alerting.Status}}</td><td>{{.Reliability.Alerting.Source}}</td><td>{{.Reliability.Alerting.Summary}}</td><td>{{.Reliability.Alerting.NextAction}}</td></tr>
<tr><td>availability_sampling</td><td>{{.Reliability.AvailabilitySampling.Status}}</td><td>{{.Reliability.AvailabilitySampling.Source}}</td><td>{{.Reliability.AvailabilitySampling.Summary}}</td><td>{{.Reliability.AvailabilitySampling.NextAction}}</td></tr>
<tr><td>long_running_operations</td><td>{{.Reliability.LongRunningOperations.Status}}</td><td>{{.Reliability.LongRunningOperations.Source}}</td><td>{{.Reliability.LongRunningOperations.Summary}}</td><td>{{.Reliability.LongRunningOperations.NextAction}}</td></tr>
</tbody></table>
<p class="muted">Incident rows are sanitized and capped. Raw details_json, raw payloads, private text, logs, tokens, hostnames, and webhook values are not included.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "validationHealthRows"}}
<table><thead><tr><th>Feed</th><th>Validator</th><th>Tooling</th><th>Artifact</th><th>Latest result</th><th>Latest at</th><th>Active feed version</th><th>Result feed version</th><th>Stale</th><th>Health</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .Feeds}}<tr><td>{{.FeedType}}</td><td>{{.ValidatorID}}<br><span class="muted">{{.ValidatorName}}</span></td><td>{{.ToolingStatus}}</td><td>{{.ArtifactStatus}}</td><td>{{.LatestResultStatus}}</td><td>{{formatTimePtr .LatestResultAt}}</td><td>{{.ActiveFeedVersionID}}</td><td>{{.LatestResultFeedVersionID}}</td><td>{{.StaleStatus}}</td><td>{{.HealthStatus}}</td><td>{{.NextAction}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
{{end}}

{{define "gtfsQualitySection"}}
<h3>{{.SourceLabel}}</h3>
<table><tbody>
<tr><th>Source</th><td>{{.Source}}</td></tr>
<tr><th>Status</th><td>{{.Status}}</td></tr>
<tr><th>Feed version</th><td>{{if .FeedVersionID}}<code>{{.FeedVersionID}}</code>{{else}}not recorded{{end}}</td></tr>
<tr><th>Latest result timestamp</th><td>{{formatTimePtr .ValidationTimestamp}}</td></tr>
{{if .IsStale}}<tr><th>Needs review</th><td>{{.StaleReason}}</td></tr>{{end}}
<tr><th>Operator summary</th><td>{{.OperatorSummary}}</td></tr>
<tr><th>Recommended action</th><td>{{.RecommendedAction}}</td></tr>
{{if .OverflowCount}}<tr><th>Hidden issue overflow</th><td>{{.OverflowCount}} notices omitted by group cap</td></tr>{{end}}
</tbody></table>
{{if .Groups}}
<table><thead><tr><th>Severity</th><th>Family</th><th>Codes</th><th>Count</th><th>Operator summary</th><th>Why it matters</th><th>Recommended action</th><th>Likely owner</th><th>Affected files</th><th>Safe fix path</th><th>Verify with</th><th>Escalate if</th><th>Samples</th><th>Overflow</th></tr></thead><tbody>
{{range .Groups}}<tr class="gtfs-quality-{{gtfsQualityGuidanceClass .}}"><td>{{.Severity}}</td><td>{{.Family}}</td><td>{{join .Codes ", "}}</td><td>{{.Count}}</td><td>{{.OperatorSummary}}</td><td>{{.WhyItMatters}}</td><td>{{.RecommendedAction}}</td><td>{{gtfsQualityLikelyOwner .}}</td><td>{{gtfsQualityAffectedFiles .}}</td><td>{{gtfsQualitySafeFixPath .Source .}}</td><td>{{gtfsQualityVerifyWith .Source .}}</td><td>{{gtfsQualityEscalation .}}</td><td>{{range .Samples}}<code>{{.}}</code><br>{{end}}</td><td>{{.OverflowCount}}</td></tr>{{end}}
</tbody></table>
{{else}}<p class="warning">No issue groups are available for this source. Next action: {{.RecommendedAction}}</p>{{end}}
{{end}}

{{define "feedTable"}}
<table><thead><tr><th>Feed</th><th>URL</th><th>Validation</th><th>Validation time</th><th>Health</th><th>Health time</th><th>Active feed version</th></tr></thead><tbody>
{{range sortedFeeds .Discovery.Feeds}}<tr><td>{{.FeedType}}</td><td>{{if .CanonicalPublicURL}}<code>{{.CanonicalPublicURL}}</code>{{else}}missing{{end}}</td><td>{{.LastValidationStatus}}</td><td>{{formatTimePtr .LastValidationAt}}</td><td>{{.LastHealthStatus}}</td><td>{{formatTimePtr .LastHealthAt}}</td><td>{{.ActiveFeedVersionID}}</td></tr>{{end}}
</tbody></table>
{{end}}

{{define "tripUpdatesQuality"}}
<h3>Trip Updates Quality Diagnostics</h3>
{{if not .TripUpdatesQuality.Recorded}}<p class="warning">{{.TripUpdatesQuality.Message}}.</p>{{else}}
<table><tbody>
<tr><th>Recorded</th><td>{{formatTimePtr .TripUpdatesQuality.SnapshotAt}}</td></tr>
<tr><th>Adapter</th><td>{{.TripUpdatesQuality.AdapterName}}</td></tr>
<tr><th>Status</th><td>{{.TripUpdatesQuality.DiagnosticsStatus}} / {{.TripUpdatesQuality.DiagnosticsReason}}</td></tr>
<tr><th>Active feed version</th><td>{{.TripUpdatesQuality.ActiveFeedVersionID}}</td></tr>
<tr><th>Unknown assignment rate</th><td>{{.TripUpdatesQuality.UnknownAssignmentRate}}</td></tr>
<tr><th>Ambiguous assignment rate</th><td>{{.TripUpdatesQuality.AmbiguousAssignmentRate}}</td></tr>
<tr><th>Stale telemetry rate</th><td>{{.TripUpdatesQuality.StaleTelemetryRate}}</td></tr>
<tr><th>Trip Updates coverage</th><td>{{.TripUpdatesQuality.TripUpdatesCoverageRate}}</td></tr>
<tr><th>Future-stop coverage</th><td>{{.TripUpdatesQuality.FutureStopCoverageRate}}</td></tr>
<tr><th>Counts</th><td>{{.TripUpdatesQuality.TripUpdatesEmitted}} emitted; {{.TripUpdatesQuality.EligiblePredictionCandidates}} eligible ETA candidates; {{.TripUpdatesQuality.UnknownAssignments}} unknown; {{.TripUpdatesQuality.AmbiguousAssignments}} ambiguous; {{.TripUpdatesQuality.StaleTelemetryRows}} stale telemetry; {{.TripUpdatesQuality.ManualOverrideAssignments}} manual overrides; {{.TripUpdatesQuality.CanceledTripsEmitted}} canceled emitted; {{.TripUpdatesQuality.CancellationAlertLinksMissing}} cancellation alerts missing</td></tr>
<tr><th>Withheld by reason</th><td>{{if .TripUpdatesQuality.WithheldByReason}}{{range .TripUpdatesQuality.WithheldByReason}}<span class="pill">{{.Label}}: {{.Count}}</span> {{end}}{{else}}none recorded{{end}}</td></tr>
<tr><th>Diagnostics persistence</th><td>{{.TripUpdatesQuality.DiagnosticsPersistenceOutcome}}</td></tr>
</tbody></table>
{{end}}
<p class="muted">This summary is based only on recorded Trip Updates diagnostics. It omits raw telemetry payloads, full score details, token fields, and private debug blobs.</p>
{{end}}

{{define "telemetry"}}
{{template "layoutStart" .}}
<h2>Telemetry Freshness</h2>
<p>Stale threshold: {{.StaleThreshold}}</p>
<div class="card-grid">
<section class="card"><h3>Make Vehicle Positions non-empty</h3><p>Create or rotate a device token, configure a device or synthetic simulator from an operator shell, send accepted telemetry to <code>/v1/telemetry</code>, then review this page and Feed Health.</p></section>
<section class="card"><h3>Why Trip Updates may be empty</h3><p>Trip Updates can be empty when telemetry is missing or stale, assignment confidence is too low, a vehicle is unknown, or the prediction adapter withholds output. Prefer empty or unknown over false certainty.</p></section>
</div>
{{if .TelemetryError}}<p class="warning">{{.TelemetryError}}. Next action: confirm the telemetry service and database are running.</p>{{else if not .Telemetry}}<p class="warning">No telemetry has been accepted yet. Next action: create or rotate a device token, configure the device, then send a sample telemetry event.</p>{{else}}
<table><thead><tr><th>Vehicle</th><th>Device</th><th>Observed</th><th>Age seconds</th><th>Freshness</th><th>Assignment</th><th>Trip</th><th>Route</th><th>Confidence</th><th>Reasons</th><th>Assignment time</th></tr></thead><tbody>
{{range .Telemetry}}<tr><td>{{.VehicleID}}</td><td>{{.DeviceID}}</td><td>{{formatTime .ObservedAt}}</td><td>{{.AgeSeconds}}</td><td>{{if .Stale}}stale{{else}}fresh{{end}}</td><td>{{if .AssignmentState}}{{.AssignmentState}}{{else}}not available{{end}}{{if .DegradedState}} / {{.DegradedState}}{{end}}</td><td>{{.TripID}}</td><td>{{.RouteID}}</td><td>{{.Confidence}}</td><td>{{join .ReasonCodes ", "}}</td><td>{{formatTimePtr .AssignmentAt}}</td></tr>{{end}}
</tbody></table>{{end}}
<p class="muted">Safe diagnostics omit raw telemetry payloads, full score details, token fields, and private debug blobs.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "telemetry-simulator"}}
{{template "layoutStart" .}}
<h2>Telemetry Simulator Guide</h2>
<p class="warning">{{.TelemetrySimulator.Boundary}}</p>
<p>The browser guide reads committed synthetic scenario metadata from <code>{{.TelemetrySimulator.ScenarioDir}}</code>. It does not read private simulator output, run commands, collect device tokens, or send telemetry from the web request.</p>
<div class="card-grid">
<section class="card"><h3>Target rules</h3>{{range .TelemetrySimulator.TargetRules}}<p>{{.}}</p>{{end}}</section>
<section class="card"><h3>Credential handling</h3>{{range .TelemetrySimulator.CredentialHandling}}<p>{{.}}</p>{{end}}</section>
<section class="card"><h3>Diagnostics policy</h3><p>{{.TelemetrySimulator.DiagnosticsPolicy}}</p></section>
</div>
{{if .TelemetrySimulator.LoadError}}<p class="bad">{{.TelemetrySimulator.LoadError}}. Next action: confirm the committed scenario fixtures are present before running simulator commands.</p>{{end}}
<h3>Operator Commands</h3>
<table><thead><tr><th>Command</th><th>What it does</th><th>Operator prep</th><th>Failure next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .TelemetrySimulator.Commands}}<tr><td><code>{{.CommandLine}}</code></td><td>{{.WhatItDoes}}</td><td>{{.OperatorPrep}}</td><td>{{.FailureNextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Synthetic Scenarios</h3>
{{if not .TelemetrySimulator.Scenarios}}<p class="warning">No simulator scenarios are available. Next action: restore the committed synthetic fixtures under <code>{{.TelemetrySimulator.ScenarioDir}}</code>.</p>{{else}}
<table><thead><tr><th>Scenario</th><th>Purpose</th><th>Events</th><th>Requirements</th><th>Expected statuses</th><th>Commands</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .TelemetrySimulator.Scenarios}}<tr><td>{{.Name}}{{if .DefaultLocal}}<br><span class="pill">default local</span>{{end}}<br><span class="muted">reference: {{.ReferenceTime}}</span></td><td>{{.Description}}</td><td>{{.EventCount}}{{if .EventLabels}}<br><span class="muted">{{join .EventLabels ", "}}</span>{{end}}</td><td>{{if .Requires}}{{range .Requires}}<span class="pill">{{.}}</span> {{end}}{{else}}none recorded{{end}}</td><td>HTTP {{range .ExpectedHTTPStatus}}<span class="pill">{{.}}</span> {{end}}<br>ingest {{if .ExpectedIngestState}}{{range .ExpectedIngestState}}<span class="pill">{{.}}</span> {{end}}{{else}}not recorded{{end}}</td><td>{{range .Commands}}<code>{{.CommandLine}}</code><br>{{end}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>{{end}}
<h3>Claim Flags</h3>
<table><tbody>
<tr><th><code>backend_command_execution_enabled</code></th><td>{{.TelemetrySimulator.ClaimFlags.BackendCommandExecutionEnabled}}</td></tr>
<tr><th><code>telemetry_sent_by_web_request</code></th><td>{{.TelemetrySimulator.ClaimFlags.TelemetrySentByWebRequest}}</td></tr>
<tr><th><code>device_token_collected_by_browser</code></th><td>{{.TelemetrySimulator.ClaimFlags.DeviceTokenCollectedByBrowser}}</td></tr>
<tr><th><code>cache_diagnostics_read_enabled</code></th><td>{{.TelemetrySimulator.ClaimFlags.CacheDiagnosticsReadEnabled}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.TelemetrySimulator.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.TelemetrySimulator.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.TelemetrySimulator.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>hardware_certification_claimed</code></th><td>{{.TelemetrySimulator.ClaimFlags.HardwareCertificationClaimed}}</td></tr>
<tr><th><code>production_avl_claimed</code></th><td>{{.TelemetrySimulator.ClaimFlags.ProductionAVLClaimed}}</td></tr>
<tr><th><code>real_realtime_claimed</code></th><td>{{.TelemetrySimulator.ClaimFlags.RealRealtimeClaimed}}</td></tr>
<tr><th><code>production_grade_eta_claimed</code></th><td>{{.TelemetrySimulator.ClaimFlags.ProductionGradeETAClaimed}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.TelemetrySimulator.ClaimFlags.ComplianceClaimed}}</td></tr>
</tbody></table>
<p class="muted">Use the device page for credential rotation, the simulator guide for synthetic <code>/v1/telemetry</code> sends, and the telemetry page for accepted-event freshness. Keep simulator diagnostics local/private unless a future evidence-specific authorization and redaction process exists.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "devices"}}
{{template "layoutStart" .}}
<h2>Device Credentials</h2>
<p class="warning">Device tokens are secrets. Store a one-time token immediately; it will not be shown again by this console.</p>
<p>The supported browser flow is rotate/rebind. If a device has no credential yet, this uses the existing rebind API path.</p>
<div class="card-grid">
<section class="card"><h3>Token status</h3><p>The table shows credential status and dates, never stored token values. New tokens are shown only once after an admin rotate/rebind action.</p></section>
<section class="card"><h3>Vehicle binding</h3><p>Each device row links a device to a vehicle, latest accepted telemetry time, freshness, assignment state, match confidence where available, and a next action.</p></section>
<section class="card"><h3>Realtime setup</h3><p>Vehicle Positions need accepted fresh telemetry. Trip Updates may still be empty until matching confidence and prediction diagnostics justify output.</p></section>
</div>
<h3>Guided Onboarding Use Cases</h3>
<div class="card-grid">
{{range .DeviceOnboarding}}<section class="card"><h3>{{.Name}}</h3><p>{{.When}}</p><p><strong>Next:</strong> {{.NextStep}}</p>{{if .AdminOnly}}<p class="muted">Admin required.</p>{{end}}</section>{{end}}
</div>
{{if .DeviceToken}}<div class="token"><h3>One-time token</h3><p>Device: {{.DeviceTokenMeta.DeviceID}} · Vehicle: {{.DeviceTokenMeta.VehicleID}} · Rotated: {{.DeviceTokenMeta.RotatedAt}}</p><p><code>{{.DeviceToken}}</code></p></div>{{end}}
{{if .DeviceError}}<p class="warning">{{.DeviceError}}</p>{{end}}
{{if .IsAdmin}}
<form method="post" action="/admin/operations/devices">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="agency_id" value="{{.AgencyID}}">
<label for="device_rebind_device_id">Device ID</label><input id="device_rebind_device_id" name="device_id" required>
<label for="device_rebind_vehicle_id">Vehicle ID</label><input id="device_rebind_vehicle_id" name="vehicle_id" required>
<label for="device_rebind_reason">Reason</label><input id="device_rebind_reason" name="reason" placeholder="rotation or rebind reason">
<button type="submit">Rotate / rebind token</button>
</form>
{{else}}
<p class="muted">Admins can rotate or rebind device tokens. This view shows operational status only.</p>
{{end}}
{{if not .DeviceRows}}<p class="warning">No device bindings are recorded. Next action: ask an admin to create or rotate the first device token and store the returned token securely.</p>{{else}}
<table><thead><tr><th>Device</th><th>Vehicle</th><th>Status</th><th>Credential dates</th><th>Latest telemetry</th><th>Freshness</th><th>Assignment</th><th>Next action</th></tr></thead><tbody>
{{range .DeviceRows}}<tr><td>{{.DeviceID}}</td><td>{{.VehicleID}}</td><td>{{.Status}}</td><td>valid from {{formatTime .ValidFrom}}<br>last used {{formatTimePtr .LastUsedAt}}<br>rotated {{formatTimePtr .RotatedAt}}<br>revoked {{formatTimePtr .RevokedAt}}</td><td>observed {{formatTimePtr .LatestObservedAt}}<br>received {{formatTimePtr .LatestReceivedAt}}<br>age seconds {{formatInt64Ptr .LatestAgeSeconds}}</td><td>{{.Freshness}}</td><td>{{.Assignment}}{{if .AssignmentSource}}<br><span class="muted">source: {{.AssignmentSource}}</span>{{end}}{{if .AssignmentAt}}<br><span class="muted">at {{formatTimePtr .AssignmentAt}}</span>{{end}}</td><td>{{.NextAction}}</td></tr>{{end}}
</tbody></table>{{end}}
<p class="muted">Safe diagnostics omit raw telemetry payloads, token hashes, private debug fields, and hardware-specific identifiers.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "maintenance"}}
{{template "layoutStart" .}}
<h2>Maintenance Center</h2>
<p class="warning">{{.Maintenance.Boundary}}</p>
<p><a href="/admin/operations/maintenance.json">Export private maintenance JSON</a> · <a href="/admin/operations/feed-health">Open feed health</a> · <a href="/admin/operations/validation-health">Open validator health</a></p>
<table><tbody>
<tr><th>Overall status</th><td>{{.Maintenance.OverallStatus}}</td></tr>
<tr><th>Generated at</th><td>{{formatTime .Maintenance.GeneratedAt}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.Maintenance.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.Maintenance.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.Maintenance.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.Maintenance.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>sla_claimed</code></th><td>{{.Maintenance.ClaimFlags.SLAClaimed}}</td></tr>
<tr><th><code>uptime_guarantee_claimed</code></th><td>{{.Maintenance.ClaimFlags.UptimeGuaranteeClaimed}}</td></tr>
</tbody></table>
<h3>Summary</h3>
<table><thead><tr><th>ID</th><th>Item</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .Maintenance.SummaryRows}}<tr id="maintenance-row-{{.ID}}"><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Next Maintenance Tasks</h3>
<table><thead><tr><th>ID</th><th>Cadence</th><th>Task</th><th>Status</th><th>Owner</th><th>Next step</th></tr></thead><tbody>
{{range .Maintenance.Tasks}}<tr id="maintenance-task-{{.ID}}"><td><code>{{.ID}}</code></td><td>{{.Cadence}}</td><td>{{.Task}}</td><td>{{.Status}}</td><td>{{.Owner}}</td><td>{{.NextStep}}</td></tr>{{end}}
</tbody></table>
<h3>Support Summary</h3>
<table><tbody>
<tr><th>Status</th><td>{{.Maintenance.SupportSummary.Status}}</td></tr>
<tr><th>Command</th><td><code>{{.Maintenance.SupportSummary.Command}}</code></td></tr>
<tr><th>Private output path</th><td><code>{{.Maintenance.SupportSummary.OutputPath}}</code></td></tr>
<tr><th>Instructions</th><td>{{range .Maintenance.SupportSummary.Instructions}}<p>{{.}}</p>{{end}}</td></tr>
</tbody></table>
<p class="muted">Use this page for routine private maintenance decisions. It intentionally says not configured or not available instead of treating missing records as ok.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "consumers"}}
{{template "layoutStart" .}}
<h2>Consumer Submission Evidence</h2>
<p class="muted">The Phase 20 docs/evidence tracker is the source for prepared packet state. These statuses are not submission, review, acceptance, or ingestion evidence.</p>
{{if .ConsumerError}}<p class="warning">{{.ConsumerError}}. The docs/evidence tracker guidance remains visible below.</p>{{end}}
<table><thead><tr><th>Target</th><th>Docs tracker status</th><th>Source</th><th>Current record</th><th>Packet path</th><th>Notes</th></tr></thead><tbody>
{{range .Consumers}}<tr><td>{{.Name}}</td><td>{{.Status}}</td><td>{{.Source}}</td><td><code>{{.CurrentPath}}</code></td><td><code>{{.PacketPath}}</code></td><td>{{.Notes}}</td></tr>{{end}}
</tbody></table>
<h3>Runtime Deployment Workflow Records</h3>
{{if .RuntimeConsumers}}<table><thead><tr><th>Target</th><th>Runtime status</th><th>Source</th><th>Updated</th><th>Notes</th></tr></thead><tbody>
{{range .RuntimeConsumers}}<tr><td>{{.Name}}</td><td>{{.Status}}</td><td>{{.Source}}</td><td>{{formatTimePtr .UpdatedAt}}</td><td>{{.Notes}}</td></tr>{{end}}
</tbody></table>{{else}}<p class="warning">No runtime consumer workflow records are available. This does not change the docs tracker prepared packet state.</p>{{end}}
<p>Docs tracker repo file path: <code>docs/evidence/consumer-submissions/README.md</code></p>
{{template "layoutEnd" .}}
{{end}}

{{define "evidence"}}
{{template "layoutStart" .}}
<h2>Evidence And Runbook Links</h2>
<p class="muted">These markdown files are repository file paths, not web routes served by this app.</p>
<table><thead><tr><th>Record</th><th>Repo file path</th><th>Last updated</th></tr></thead><tbody>
{{range .Links}}<tr><td>{{.Label}}</td><td><code>{{.Path}}</code></td><td>{{.UpdatedAt}}</td></tr>{{end}}
</tbody></table>
<p class="muted">These links help operators find repo/deployment evidence. They do not assert consumer acceptance, hosted service availability, agency endorsement, or universal production readiness.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "setup"}}
{{template "layoutStart" .}}
<h2>Guided Setup Checklist</h2>
{{if .SetupNotice}}<p class="ok">{{.SetupNotice}}</p>{{end}}
{{if .SetupError}}<p class="bad">{{.SetupError}}</p>{{end}}
<p><a href="/admin/operations/checklist">Open private operator checklist</a> · <a href="/admin/operations/checklist.json">Export private checklist JSON</a></p>
<p class="muted">Each status is tied to a named source. Missing records stay missing until publication metadata, feed discovery, validation records, device bindings, telemetry, docs tracker records, or evidence links support a stronger statement.</p>
<table><thead><tr><th>Step</th><th>Status</th><th>Status source</th><th>Evidence signal</th><th>Next action</th></tr></thead><tbody>
{{range .SetupSteps}}<tr><td>{{.Name}}</td><td>{{.Status}}</td><td>{{.Source}}</td><td>{{.Evidence}}</td><td>{{if .ActionURL}}<a href="{{.ActionURL}}">{{.NextAction}}</a>{{else}}{{.NextAction}}{{end}}</td></tr>{{end}}
</tbody></table>

<h2 id="publication-metadata">Publication Metadata</h2>
<p class="muted">Source: publication metadata and feed discovery. This form uses the existing publication bootstrap/update repository behavior and derives agency ID from the authenticated admin principal.</p>
{{if .PublicationError}}<p class="warning">{{.PublicationError}}. Existing JSON admin API path: <code>/admin/publication/bootstrap</code>.</p>{{end}}
<table><tbody>
<tr><th>Agency ID</th><td><code>{{.AgencyID}}</code> (read-only authenticated principal)</td></tr>
<tr><th>Agency name</th><td>{{if .Discovery.AgencyName}}{{.Discovery.AgencyName}}{{else}}missing{{end}}</td></tr>
<tr><th>Public base URL</th><td>{{if .PublicationConfig.PublicBaseURL}}{{.PublicationConfig.PublicBaseURL}}{{else if .Discovery.PublicBaseURL}}{{.Discovery.PublicBaseURL}}{{else}}missing{{end}}</td></tr>
<tr><th>Feed base URL</th><td>{{if .PublicationConfig.FeedBaseURL}}{{.PublicationConfig.FeedBaseURL}}{{else}}missing{{end}}</td></tr>
<tr><th>License name</th><td>{{if .PublicationConfig.LicenseName}}{{.PublicationConfig.LicenseName}}{{else if .Discovery.License.Name}}{{.Discovery.License.Name}}{{else}}missing{{end}}</td></tr>
<tr><th>License URL</th><td>{{if .PublicationConfig.LicenseURL}}{{.PublicationConfig.LicenseURL}}{{else if .Discovery.License.URL}}{{.Discovery.License.URL}}{{else}}missing{{end}}</td></tr>
<tr><th>Technical contact</th><td>{{if .PublicationConfig.TechnicalContactEmail}}{{.PublicationConfig.TechnicalContactEmail}}{{else if .Discovery.TechnicalContactEmail}}{{.Discovery.TechnicalContactEmail}}{{else}}missing{{end}}</td></tr>
<tr><th>Publication environment</th><td>{{if .PublicationConfig.PublicationEnvironment}}{{.PublicationConfig.PublicationEnvironment}}{{else if .Discovery.PublicationEnvironment}}{{.Discovery.PublicationEnvironment}}{{else}}missing{{end}}</td></tr>
</tbody></table>
<form method="post" action="/admin/operations/setup">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="action" value="publication_bootstrap">
<label for="setup_public_base_url">Public base URL</label><input id="setup_public_base_url" type="url" name="public_base_url" maxlength="2048" required value="{{if .PublicationConfig.PublicBaseURL}}{{.PublicationConfig.PublicBaseURL}}{{else}}{{.Discovery.PublicBaseURL}}{{end}}">
<label for="setup_feed_base_url">Feed base URL</label><input id="setup_feed_base_url" type="url" name="feed_base_url" maxlength="2048" required value="{{.PublicationConfig.FeedBaseURL}}">
<label for="setup_technical_contact_email">Technical contact email</label><input id="setup_technical_contact_email" type="email" name="technical_contact_email" maxlength="320" value="{{if .PublicationConfig.TechnicalContactEmail}}{{.PublicationConfig.TechnicalContactEmail}}{{else}}{{.Discovery.TechnicalContactEmail}}{{end}}">
<label for="setup_license_name">License name</label><input id="setup_license_name" name="license_name" maxlength="160" value="{{if .PublicationConfig.LicenseName}}{{.PublicationConfig.LicenseName}}{{else}}{{.Discovery.License.Name}}{{end}}">
<label for="setup_license_url">License URL</label><input id="setup_license_url" type="url" name="license_url" maxlength="2048" value="{{if .PublicationConfig.LicenseURL}}{{.PublicationConfig.LicenseURL}}{{else}}{{.Discovery.License.URL}}{{end}}">
<label for="setup_publication_environment">Publication environment</label><input id="setup_publication_environment" name="publication_environment" maxlength="64" value="{{publicationEnvValue .}}">
<button type="submit">Store publication metadata</button>
</form>

<h2>GTFS Import And Authoring</h2>
<p>Source: feed discovery and the existing GTFS importer. Browser import is admin-only, size-limited, temporary-file based, and uses the same validation and publish pipeline as the CLI import path.</p>
<table><tbody>
<tr><th>Browser import</th><td><a href="/admin/operations/gtfs-import">Import a GTFS ZIP by upload or safe URL</a>.</td></tr>
<tr><th>CLI ZIP import</th><td>Use the existing GTFS import flow documented in <code>docs/tutorials/real-agency-gtfs-onboarding.md</code>.</td></tr>
<tr><th>Typed authoring</th><td><a href="/admin/gtfs-studio">Open GTFS Studio</a> for draft authoring and publish.</td></tr>
<tr><th>Validation triage</th><td>Use <code>docs/tutorials/gtfs-validation-triage.md</code> and the validation form below.</td></tr>
<tr><th>GTFS quality triage</th><td><a href="/admin/operations/gtfs-quality">Review canonical validator and internal importer actions</a>.</td></tr>
<tr><th>Validator health</th><td><a href="/admin/operations/validation-health">Review private validator installation, artifact, stale-result, and next-action diagnostics</a>.</td></tr>
<tr><th>Active feed verification</th><td><a href="/admin/operations/feeds">Review feed discovery and validation records</a>.</td></tr>
</tbody></table>

<h2 id="validation">Validation</h2>
<p class="muted">Source: validation records. The browser chooses only feed type; the server maps it to an allowlisted validator. Validation is supporting evidence only, not consumer acceptance or compliance.</p>
<p><a href="/admin/operations/validation-health">Open private validator health diagnostics</a></p>
{{if .ValidationResult}}<p class="ok">Last run from this page: {{.ValidationResult.FeedType}} validation {{.ValidationResult.Status}} with {{.ValidationResult.ErrorCount}} errors and {{.ValidationResult.WarningCount}} warnings.</p>{{end}}
<form method="post" action="/admin/operations/setup">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="action" value="run_validation">
<label for="setup_validation_feed_type">Feed type</label><select id="setup_validation_feed_type" name="feed_type">
<option value="schedule">schedule</option>
<option value="vehicle_positions">vehicle_positions</option>
<option value="trip_updates">trip_updates</option>
<option value="alerts">alerts</option>
</select>
<button type="submit">Run allowlisted validation</button>
</form>
{{if .DiscoveryError}}<p class="warning">No validation records are available because publication metadata is missing.</p>{{else}}{{template "feedTable" .}}{{end}}

<h2>Device And Telemetry Setup</h2>
<p>Source: device bindings and telemetry repository. Device tokens are one-time secrets and are only shown by the existing rotate/rebind flow.</p>
<table><tbody>
<tr><th>Device bindings</th><td>{{if .DeviceError}}{{.DeviceError}}{{else}}{{len .Devices}} binding records{{end}}</td></tr>
<tr><th>Latest telemetry</th><td>{{if .TelemetryError}}{{.TelemetryError}}{{else if .TelemetryUpdatedAt}}{{formatTimePtr .TelemetryUpdatedAt}}{{else}}not observed yet{{end}}</td></tr>
<tr><th>Stale telemetry</th><td>{{if .TelemetryError}}not available{{else}}{{.StaleCount}} stale latest rows using threshold {{.StaleThreshold}}{{end}}</td></tr>
<tr><th>Next action</th><td><a href="/admin/operations/devices">Manage device bindings</a>; use <code>scripts/device-onboarding.sh sample --dry-run</code> or <code>simulate --dry-run</code> to preview helper calls.</td></tr>
</tbody></table>

<h2>Alerts, Overrides, Consumers, Evidence</h2>
<table><tbody>
<tr><th>Alerts</th><td>Source: feed discovery and Alerts Console. <a href="/admin/alerts/console">Create, publish, or archive alerts</a>. Alerts feed availability does not prove consumer acceptance.</td></tr>
<tr><th>Manual overrides/review</th><td>Deferred in Phase 26 because a safe browser view would need carefully bounded summaries and must avoid raw diagnostics or new mutation semantics.</td></tr>
<tr><th>Consumer packets</th><td>Source: docs/evidence tracker. <a href="/admin/operations/consumers">Review all seven prepared packet records</a>; prepared is not submitted or accepted.</td></tr>
<tr><th>Evidence/readiness</th><td>Source: evidence links. <a href="/admin/operations/evidence">Open evidence link index</a>.</td></tr>
</tbody></table>
{{template "layoutEnd" .}}
{{end}}

{{define "checklist"}}
{{template "layoutStart" .}}
<h2>Private Operator Checklist</h2>
<p class="warning">This checklist is private operator diagnostics. It is not evidence, not an evidence packet, not compliance proof, not agency approval, not consumer acceptance, and not production readiness.</p>
<p><a href="/admin/operations/checklist.json">Export private checklist JSON</a> · <a href="/admin/operations/gtfs-quality">Open GTFS quality triage</a> · <a href="/admin/operations/validation-health">Open validator health diagnostics</a></p>
<table><tbody>
<tr><th><code>external_evidence_created</code></th><td>{{.Checklist.Flags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>final_root_evidence_created</code></th><td>{{.Checklist.Flags.FinalRootEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.Checklist.Flags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.Checklist.Flags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.Checklist.Flags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>agency_approval_claimed</code></th><td>{{.Checklist.Flags.AgencyApprovalClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.Checklist.Flags.ConsumerAcceptanceClaimed}}</td></tr>
</tbody></table>
{{range .Checklist.Groups}}
<h3>{{.Label}}</h3>
<table><thead><tr><th>ID</th><th>Row</th><th>Status</th><th>Source</th><th>Current signal</th><th>Next action</th><th>Boundary</th><th>Heuristics</th><th>Docs</th></tr></thead><tbody>
{{range .Rows}}<tr><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.Source}}</td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.ClaimBoundary}}</td><td>{{range .HeuristicLabels}}<span class="pill">{{humanHeuristic .}}</span> {{end}}</td><td>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</td></tr>{{end}}
</tbody></table>
{{end}}
{{template "layoutEnd" .}}
{{end}}
`))
