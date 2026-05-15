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

	"open-transit-rt/internal/admincontrol"
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
	ConsumerPreparation    operationsConsumerPreparationView
	ReadinessItems         []readinessItemView
	ReadinessV2            operationsReadinessV2View
	Checklist              operatorChecklistView
	Cockpit                operationsCockpitView
	FirstRun               operationsFirstRunView
	Launchpad              agencyLaunchpadView
	SetupWizard            operationsSetupWizardView
	ConnectorHub           connectorHubView
	ConnectorWorkbench     connectorWorkbenchView
	ConnectorTests         connectorTestsView
	Help                   operationsHelpView
	ContextHelp            operationsContextHelp
	FeedReadiness          operationsFeedReadinessView
	FeedHealth             operationsFeedHealthView
	ValidationCenter       operationsValidationCenterView
	Realtime               operationsRealtimeView
	PredictionLab          predictionLabView
	Maintenance            operationsMaintenanceView
	Access                 operationsAccessView
	Audit                  operationsAuditView
	TelemetrySimulator     operationsTelemetrySimulatorView
	GTFSImportResult       *gtfsImportResultView
	GTFSImportSource       gtfsImportSourceReview
	GTFSImportNotice       string
	GTFSImportError        string
	GTFSQuality            compliance.GTFSQualityTriage
	GTFSQualityGuidance    operationsGTFSQualityGuidanceView
	GTFSWorkbench          operationsGTFSWorkbenchView
	GTFSQualityNotice      string
	GTFSQualityError       string
	ValidationHealth       compliance.ValidationHealthSummary
	ValidationHealthNotice string
	ValidationHealthError  string
	Reliability            compliance.ReliabilitySummary
	ReliabilityError       string
	AgencyScope            operationsAgencyScopeView
	IsAdmin                bool
	PrincipalRoles         []string
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
	AdapterDetails                map[string]any
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
	case "connectors/workbench":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderConnectorWorkbench(w, r)
	case "connectors/workbench.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderConnectorWorkbenchJSON(w, r)
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
	case "validation-center":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderValidationCenter(w, r)
	case "validation-center.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderValidationCenterJSON(w, r)
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
	case "realtime":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderRealtime(w, r)
	case "realtime.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderRealtimeJSON(w, r)
	case "prediction-lab":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderPredictionLab(w, r)
	case "prediction-lab.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderPredictionLabJSON(w, r)
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
	case "gtfs-workbench":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderGTFSWorkbench(w, r)
	case "gtfs-workbench.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderGTFSWorkbenchJSON(w, r)
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
	case "access":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderAccess(w, r)
	case "access.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderAccessJSON(w, r)
	case "audit":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderAudit(w, r)
	case "audit.json":
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderAuditJSON(w, r)
	case "feeds", "telemetry", "devices", "consumers", "evidence", "setup":
		w.Header().Set("Cache-Control", "no-store")
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

func (h *handler) renderGTFSWorkbench(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "gtfs-workbench")
	renderOperationsTemplate(w, "gtfs-workbench", page)
}

func (h *handler) renderGTFSWorkbenchJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "gtfs-workbench")
	writeJSON(w, http.StatusOK, page.GTFSWorkbench)
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

func (h *handler) operationsValidationHealthRefreshCommandJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	started := time.Now().UTC().Truncate(time.Second)
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	if err := r.ParseForm(); err != nil {
		result := admincontrol.NewResult(
			"validation_health.refresh",
			admincontrol.StatusBlocked,
			started,
			"Validator health refresh was blocked because the form body is invalid or too large.",
			[]string{"Retry with the private console refresh action and no client-supplied command fields."},
			[]admincontrol.Error{{Code: "request_body_blocked", Message: "form body invalid or too large"}},
		)
		writeJSON(w, http.StatusRequestEntityTooLarge, result)
		return
	}
	if principal.Method == auth.MethodCookie && strings.TrimSpace(h.csrfSecret) != "" && strings.TrimSpace(r.FormValue("csrf_token")) != csrfToken(h.csrfSecret, principal) {
		result := admincontrol.NewResult(
			"validation_health.refresh",
			admincontrol.StatusBlocked,
			started,
			"Validator health refresh was blocked because the CSRF token is invalid.",
			[]string{"Reload the private console and retry the refresh action."},
			[]admincontrol.Error{{Code: "csrf_invalid", Message: "invalid CSRF token"}},
		)
		writeJSON(w, http.StatusForbidden, result)
		return
	}
	if err := rejectValidationHealthUnexpectedFields(r); err != nil {
		result := admincontrol.NewResult(
			"validation_health.refresh",
			admincontrol.StatusBlocked,
			started,
			err.Error(),
			[]string{"Use the server-owned validation health refresh action without client-supplied execution fields."},
			[]admincontrol.Error{{Code: "unsupported_field", Message: err.Error()}},
		)
		writeJSON(w, http.StatusBadRequest, result)
		return
	}
	action := strings.TrimSpace(r.FormValue("action"))
	if action != "" && action != "refresh" && action != "validation_health.refresh" {
		result := admincontrol.NewResult(
			"validation_health.refresh",
			admincontrol.StatusBlocked,
			started,
			"Validator health refresh accepts only the validation_health.refresh action.",
			[]string{"Use the private refresh control, or use the existing admin-only validator run action when an actual validator run is intended."},
			[]admincontrol.Error{{Code: "unsupported_action", Message: "unsupported validation health refresh action"}},
		)
		writeJSON(w, http.StatusBadRequest, result)
		return
	}
	page := h.buildOperationsPage(r, principal, "validation-health")
	result := validationHealthRefreshCommandResult(started, page.ValidationHealth)
	writeJSON(w, http.StatusOK, result)
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

func validationHealthRefreshCommandResult(started time.Time, summary compliance.ValidationHealthSummary) admincontrol.Result {
	status := validationHealthCommandStatus(summary.OverallStatus)
	return admincontrol.NewResult(
		"validation_health.refresh",
		status,
		started,
		fmt.Sprintf("Validation health summary refreshed from existing private records. Overall status: %s; feed rows: %d.", firstNonEmpty(summary.OverallStatus, "unknown"), len(summary.Feeds)),
		validationHealthCommandNextActions(summary),
		nil,
	)
}

func validationHealthCommandStatus(status string) admincontrol.Status {
	switch status {
	case compliance.ValidationHealthStatusBlocked,
		compliance.ValidationHealthStatusMissingTooling,
		compliance.ValidationHealthStatusMisconfiguredTooling,
		compliance.ValidationHealthStatusArtifactUnavailable:
		return admincontrol.StatusBlocked
	case compliance.ValidationHealthStatusRecorded,
		compliance.ValidationHealthStatusRunnable,
		compliance.ValidationHealthStatusConfigured,
		compliance.ValidationHealthStatusInstalled,
		compliance.ValidationHealthStatusStub,
		compliance.ValidationHealthStatusConfiguredForTests,
		compliance.ValidationHealthStatusSkipped:
		return admincontrol.StatusOK
	case compliance.ValidationHealthStatusFailed,
		compliance.ValidationHealthStatusStale,
		compliance.ValidationHealthStatusNeedsReview,
		compliance.ValidationHealthStatusNotRun,
		compliance.ValidationHealthStatusUnknown:
		return admincontrol.StatusNeedsReview
	default:
		return admincontrol.StatusNeedsReview
	}
}

func validationHealthCommandNextActions(summary compliance.ValidationHealthSummary) []string {
	seen := map[string]bool{}
	var actions []string
	for _, row := range summary.Feeds {
		next := strings.TrimSpace(row.NextAction)
		if next == "" || seen[next] {
			continue
		}
		actions = append(actions, next)
		seen[next] = true
		if len(actions) >= 4 {
			break
		}
	}
	if len(actions) == 0 {
		actions = append(actions, "Review validator health rows before stronger readiness language.")
	}
	actions = append(actions, "This refresh writes nothing, changes no public feed output, creates no evidence, and moves no consumer status.")
	return actions
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
		Title:            operationsPageTitle(section),
		AgencyID:         principal.AgencyID,
		GeneratedAt:      now,
		EnvironmentLabel: firstNonEmpty(os.Getenv("PUBLICATION_ENVIRONMENT"), "unknown"),
		CSRFToken:        csrfToken(h.csrfSecret, principal),
		Section:          section,
		NavGroups:        operationsNavGroups(section),
		StaleThreshold:   staleThreshold(),
		AgencyScope:      buildOperationsAgencyScope(principal),
		IsAdmin:          principal.HasAny(auth.RoleAdmin),
		PrincipalRoles:   safePrincipalRoles(principal.Roles),
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
	page.ConsumerPreparation = buildOperationsConsumerPreparation(page)

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
	page.GTFSWorkbench = h.buildGTFSWorkbenchView(r, page)
	page.Reliability, page.ReliabilityError = h.reliabilitySummary(r, principal.AgencyID, now)
	page.FeedReadiness = buildOperationsFeedReadiness(page)
	page.FeedHealth = buildOperationsFeedHealth(page)
	page.Realtime = buildOperationsRealtime(page)
	page.PredictionLab = buildPredictionLab(page)
	page.Maintenance = buildOperationsMaintenance(page)
	page.SetupSteps = setupSteps(page)
	page.ReadinessItems = readinessItems(page)
	page.ReadinessV2 = buildOperationsReadinessV2(page)
	page.ValidationCenter = buildOperationsValidationCenter(page)
	page.Checklist = buildOperatorChecklist(page)
	page.FirstRun = buildOperationsFirstRun(page)
	page.TelemetrySimulator = buildOperationsTelemetrySimulator(page)
	page.Launchpad = buildAgencyLaunchpad(page)
	page.SetupWizard = buildOperationsSetupWizard(page)
	page.ConnectorHub = buildConnectorHub(page)
	page.ConnectorWorkbench = buildConnectorWorkbench(page)
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
		AdapterDetails:                summary.AdapterDetails,
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
			NextAction: "Enter an operator-provided open license URL and technical contact, or keep this marked missing.",
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
			NextAction: "Review configured feed URLs and health records. Verification is supporting evidence only.",
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
			Name:          "Configured public route URLs",
			Status:        readinessStatus(page.Discovery.Readiness.AllRequiredFeedsListed && page.Discovery.Readiness.HTTPSURLs, page.DiscoveryError),
			Source:        "feed discovery and published_feed records",
			Evidence:      stableURLEvidence(page),
			NextAction:    "Confirm every configured public route URL is stable, HTTPS in deployment, and served through the intended feed root.",
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
			NextAction:    "Replace placeholders with operator-provided open license and monitored technical contact values.",
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

func safePrincipalRoles(roles []auth.Role) []string {
	if len(roles) == 0 {
		return []string{"none"}
	}
	out := make([]string, 0, len(roles))
	for _, role := range roles {
		value := strings.TrimSpace(string(role))
		if value != "" {
			out = append(out, value)
		}
	}
	if len(out) == 0 {
		return []string{"none"}
	}
	sort.Strings(out)
	return out
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
	"statusClass": func(value string) string {
		normalized := strings.ToLower(strings.TrimSpace(value))
		normalized = strings.NewReplacer("_", "-", " ", "-", "/", "-").Replace(normalized)
		if normalized == "" {
			return "unknown"
		}
		return normalized
	},
	"operationsCSS": operationsCSS,
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
<style>{{operationsCSS}}</style></head><body>
<a class="skip-link" href="#operations-main">Skip to main content</a>
<a class="skip-link" href="#operations-nav">Skip to section navigation</a>
<header class="operations-header" role="banner">
<p class="app-kicker">Private operations control plane</p>
<p class="app-breadcrumb"><a href="/admin/operations">Operations Console</a> / {{.Title}}</p>
<h1 id="operations-page-title">{{.Title}}</h1>
<p class="app-meta"><span>Agency: <strong>{{.AgencyID}}</strong></span><span>environment: <span class="pill">{{.EnvironmentLabel}}</span></span><span>generated: {{formatTime .GeneratedAt}}</span></p>
<section class="scope-banner" aria-labelledby="agency-scope-heading">
<h2 id="agency-scope-heading">Agency scope</h2>
<dl class="scope-grid">
<div><dt>Current agency</dt><dd><code>{{.AgencyScope.AgencyID}}</code></dd></div>
<div><dt>Status</dt><dd><span class="status-chip status-{{statusClass .AgencyScope.Status}}">{{.AgencyScope.Status}}</span></dd></div>
<div><dt>Source</dt><dd>{{.AgencyScope.Source}}</dd></div>
<div><dt>Roles</dt><dd>{{join .AgencyScope.Roles ", "}}</dd></div>
<div><dt>Switcher</dt><dd>{{.AgencyScope.SwitcherStatus}}</dd></div>
<div><dt>Query rule</dt><dd>{{.AgencyScope.QueryRule}}</dd></div>
</dl>
<p><strong>Next:</strong> {{.AgencyScope.NextAction}}</p>
<p class="muted"><strong>Does not prove:</strong> {{.AgencyScope.DoesNotProve}}</p>
</section>
</header>
<nav id="operations-nav" class="operations-nav" aria-label="Operations Console sections">
{{range .NavGroups}}<section class="nav-group" aria-labelledby="nav-group-{{.ID}}">
<p id="nav-group-{{.ID}}" class="nav-group-label">{{.Label}}</p>
<div class="nav-links">{{range .Items}}<a class="nav-link{{if .Current}} current{{end}}" href="{{.Href}}"{{if .Current}} aria-current="page"{{end}}>{{.Label}}{{if .ExternalAdminSurface}} <span class="nav-surface">admin surface</span>{{end}}</a>{{end}}</div>
</section>{{end}}
</nav>
{{if ne .Section "dashboard"}}{{template "contextHelpPanel" .}}{{end}}
<main id="operations-main" tabindex="-1" aria-labelledby="operations-page-title">
{{end}}

{{define "contextHelpPanel"}}
{{if .ContextHelp.Topics}}<aside class="context-help" aria-labelledby="context-help-heading">
<h2 id="context-help-heading">Help for {{.ContextHelp.Label}}</h2>
<div class="context-help-grid">{{range .ContextHelp.Topics}}<section class="context-help-topic"><h3>{{.Label}}</h3><p>{{.Summary}}</p><p><strong>Next:</strong> {{.NextAction}}</p><p><a href="/admin/operations/help#help-{{.ID}}">Open topic</a></p></section>{{end}}</div>
<p class="muted"><a href="{{.ContextHelp.AllTopicsURL}}">Open all help topics</a> · <a href="{{.ContextHelp.JSONURL}}">Export private help JSON</a></p>
</aside>{{end}}
{{end}}
{{define "layoutEnd"}}</main><script src="/admin/operations/assets/operations.js" defer></script></body></html>{{end}}

{{define "access"}}
{{template "layoutStart" .}}
<h2>Access &amp; Roles</h2>
<p class="warning">{{.Access.Boundary}}</p>
<p><a href="/admin/operations/access.json">Export private access JSON</a> · <a href="/admin/operations">Back to Start Here</a></p>
<table><tbody>
<tr><th>Agency</th><td><code>{{.Access.AgencyID}}</code></td></tr>
<tr><th>Current roles</th><td>{{join .Access.CurrentRoles ", "}}</td></tr>
<tr><th>Generated at</th><td>{{formatTime .Access.GeneratedAt}}</td></tr>
</tbody></table>
<h3>Role Permissions</h3>
<table><thead><tr><th>Role</th><th>Current session</th><th>Review access</th><th>Mutation access</th><th>Technical helper note</th><th>Does not prove</th></tr></thead><tbody>
{{range .Access.Roles}}<tr id="access-role-{{.ID}}"><td>{{.Label}}</td><td>{{.Current}}</td><td>{{.ReviewAccess}}</td><td>{{.MutationAccess}}</td><td>{{.TechnicalHelperNote}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Access-Denied Guidance</h3>
<table><thead><tr><th>Scenario</th><th>What happened</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .Access.Denied}}<tr id="access-denied-{{.ID}}"><td>{{.Scenario}}</td><td>{{.WhatHappened}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
{{template "layoutEnd" .}}
{{end}}

{{define "audit"}}
{{template "layoutStart" .}}
<h2>Audit Log</h2>
<p class="warning">{{.Audit.Boundary}}</p>
<p><a href="/admin/operations/audit.json">Export private audit JSON</a> · <a href="/admin/operations/access">Open Access &amp; Roles</a></p>
<table><tbody>
<tr><th>Agency</th><td><code>{{.Audit.AgencyID}}</code></td></tr>
<tr><th>Status</th><td>{{.Audit.Status}}</td></tr>
<tr><th>Generated at</th><td>{{formatTime .Audit.GeneratedAt}}</td></tr>
<tr><th>Next action</th><td>{{.Audit.NextAction}}</td></tr>
</tbody></table>
{{if .Audit.Rows}}
<table><caption>Recent scoped audit metadata</caption><thead><tr><th>ID</th><th>Created</th><th>Action</th><th>Entity</th><th>Actor recorded</th><th>Reason recorded</th><th>Old value</th><th>New value</th><th>Does not show</th></tr></thead><tbody>
{{range .Audit.Rows}}<tr id="audit-row-{{.ID}}"><td>{{.ID}}</td><td>{{.CreatedAt}}</td><td>{{.Action}}</td><td>{{.EntityType}} / {{.EntityRef}}</td><td>{{.ActorRecorded}}</td><td>{{.ReasonRecorded}}</td><td>{{.OldValueRecorded}}</td><td>{{.NewValueRecorded}}</td><td>{{.DoesNotShow}}</td></tr>{{end}}
</tbody></table>
{{else}}<p class="muted">{{.Audit.EmptyState}}</p>{{end}}
{{template "layoutEnd" .}}
{{end}}

{{define "help"}}
{{template "layoutStart" .}}
<h2>Operations Console Help</h2>
<p class="warning">{{.Help.Boundary}}</p>
<div class="card-grid" aria-label="Help empty or blocked state guidance">
<section class="card empty-state">
<h3>When a page is empty or blocked</h3>
<p><strong>What am I seeing?</strong> The console is showing the latest private records it can read for this agency.</p>
<p><strong>Is this bad?</strong> Not always. Empty often means first-run setup has not produced that source record yet; blocked means an operator action or configuration is needed.</p>
<p><strong>What should I do next?</strong> Open Start Here, follow the linked page, and keep missing records missing until the underlying setup, import, validator, telemetry, connector, or maintenance signal exists.</p>
<p><strong>Can I do it in the browser?</strong> Review can happen in the browser; admin-only browser actions stay on their existing private pages.</p>
<p><strong>When do I need a technical helper?</strong> Use one for local startup, validator/tooling setup, operator-shell commands, backup/restore configuration, deployment diagnostics, or external integration prep.</p>
<p><strong>What this does not prove:</strong> Help text does not prove compliance, agency adoption, consumer acceptance, final-root readiness, hosted service availability, production readiness, vendor compatibility, SLA coverage, hardware certification, or ETA quality.</p>
</section>
</div>
<p><a href="/admin/operations/help.json">Export private help JSON</a> · <a href="/admin/operations">Back to Operations Console</a></p>
<h3>Role-Based Tours</h3>
<div class="card-grid" aria-label="Role-based tours">
{{range .Help.RoleTours}}<section class="card" id="help-role-{{.ID}}">
<h3>{{.Label}}</h3>
<p>{{.Who}}</p>
<p><strong>Start here:</strong> <a href="{{.StartHere}}">{{.StartHere}}</a></p>
<p><strong>Review first:</strong> {{.ReviewFirst}}</p>
<p><strong>First actions:</strong> {{.FirstActions}}</p>
<p><strong>Ask for help when:</strong> {{.EscalateWhen}}</p>
<p><strong>Does not prove:</strong> {{.DoesNotProve}}</p>
<p><strong>Console:</strong> {{range .AdminLinks}}<a href="{{.}}">{{.}}</a> {{end}}</p>
<p><strong>Docs:</strong> {{range .DocsLinks}}<code>{{.}}</code> {{end}}</p>
</section>{{end}}
</div>
<h3>First-Week Checklist</h3>
<table><thead><tr><th>When</th><th>Role</th><th>Task</th><th>Review</th><th>Done when</th><th>Next action</th><th>Console</th><th>Does not prove</th></tr></thead><tbody>
{{range .Help.FirstWeek}}<tr id="help-first-week-{{.ID}}"><td>{{.Day}}</td><td>{{.Role}}</td><td>{{.Task}}</td><td>{{.Review}}</td><td>{{.DoneWhen}}</td><td>{{.NextAction}}</td><td><a href="{{.ConsoleLink}}">{{.ConsoleLink}}</a></td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Plain-Language Glossary</h3>
<table><thead><tr><th>Term</th><th>Plain meaning</th><th>Technical meaning</th><th>Where to review</th><th>Docs</th><th>Does not prove</th></tr></thead><tbody>
{{range .Help.Glossary}}<tr id="help-glossary-{{.ID}}"><td>{{.Term}}</td><td>{{.PlainMeaning}}</td><td>{{.TechnicalMeaning}}</td><td>{{.WhereToReview}}</td><td>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Common Mistake Recovery</h3>
<table><thead><tr><th>What the operator sees</th><th>Likely cause</th><th>Safe next step</th><th>Escalate when</th><th>Console</th><th>Does not prove</th></tr></thead><tbody>
{{range .Help.Recovery}}<tr id="help-recovery-{{.ID}}"><td>{{.WhatOperatorSees}}</td><td>{{.LikelyCause}}</td><td>{{.SafeNextStep}}</td><td>{{.EscalationTrigger}}</td><td><a href="{{.ConsoleLink}}">{{.ConsoleLink}}</a></td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Printable Staff Training Guide</h3>
<table><tbody>
<tr><th>Guide</th><td>{{.Help.TrainingGuide.Label}}</td></tr>
<tr><th>Repo path</th><td><code>{{.Help.TrainingGuide.DocsPath}}</code></td></tr>
<tr><th>Audience</th><td>{{.Help.TrainingGuide.Audience}}</td></tr>
<tr><th>How to use</th><td>{{.Help.TrainingGuide.HowToUse}}</td></tr>
<tr><th>Boundary</th><td>{{.Help.TrainingGuide.Boundary}}</td></tr>
</tbody></table>
<h3>Quick Tasks</h3>
<table><thead><tr><th>Task</th><th>Role</th><th>Review steps</th><th>Done when</th><th>Escalate when</th><th>Console</th><th>Does not prove</th></tr></thead><tbody>
{{range .Help.QuickTasks}}<tr id="help-quick-task-{{.ID}}"><td>{{.Label}}</td><td>{{.PrimaryRole}}</td><td>{{.ReviewSteps}}</td><td>{{.DoneWhen}}</td><td>{{.Escalation}}</td><td><a href="{{.ConsoleLink}}">{{.ConsoleLink}}</a></td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Staff Handoff Checklist</h3>
<table><thead><tr><th>Area</th><th>From</th><th>To</th><th>Confirm</th><th>Console</th><th>Does not prove</th></tr></thead><tbody>
{{range .Help.Handoff}}<tr id="help-handoff-{{.ID}}"><td>{{.Area}}</td><td>{{.FromRole}}</td><td>{{.ToRole}}</td><td>{{.Confirm}}</td><td><a href="{{.ConsoleLink}}">{{.ConsoleLink}}</a></td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
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
<h2 id="first-run-heading">Start Here: First Actions</h2>
<p>{{.Boundary}}</p>
<p class="muted">{{.LocalDemoDeploymentEvidenceBoundary}}</p>
<p><strong>Task status:</strong> {{.Counts.Tasks}} tasks · ok {{index .Counts.Statuses "ok"}} · needs review {{index .Counts.Statuses "needs_review"}} · missing {{index .Counts.Statuses "missing"}} · blocked {{index .Counts.Statuses "blocked"}} · unknown {{index .Counts.Statuses "unknown"}}</p>
<div class="card-grid" aria-label="First-run evaluator paths">
{{range .Paths}}<section class="card path-card path-{{.ID}}" id="first-run-path-{{.ID}}">
<h3>{{.Label}}</h3>
<p><strong>Current signal:</strong> {{.CurrentSignal}}</p>
<p><strong>What it means:</strong> {{.Meaning}}</p>
<p><strong>First action:</strong> {{.FirstAction}}</p>
<p><strong>Console:</strong> <a href="{{.UILink}}">{{.UILink}}</a></p>
<p><strong>Docs:</strong> <code>{{.DocsLink}}</code></p>
<p><strong>Does not prove:</strong> {{.DoesNotProve}}</p>
</section>{{end}}
</div>
<h3>Copy These Five Configured Feed URLs</h3>
<p class="section-note">Use these configured local/reference paths for private review. Missing stays missing until publication metadata or feed records exist.</p>
<div class="feed-copy-grid" aria-label="Copyable configured feed URLs">
{{range .FeedURLs}}<section class="feed-url-card" id="first-run-feed-{{.ID}}" data-copy-card>
<h3>{{.Label}}</h3>
<p><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span> <code>{{.ID}}</code></p>
<code class="copy-value" data-copy-value="{{.CopyValue}}">{{.CopyValue}}</code>
<p><strong>Current link:</strong> {{if .URL}}<a href="{{.URL}}">{{.URL}}</a>{{else}}missing{{end}}</p>
<p><strong>Next:</strong> {{.NextAction}}</p>
<p class="muted"><strong>Does not prove:</strong> {{.DoesNotProve}}</p>
</section>{{end}}
</div>
<h3>First-Run Acceptance Tasks</h3>
<table><thead><tr><th>Order</th><th>Task</th><th>Status</th><th>Current signal</th><th>What it means</th><th>Next action</th><th>Console</th><th>Docs</th><th>Does not prove</th></tr></thead><tbody>
{{range .Tasks}}<tr><td>{{.Order}}</td><td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.Meaning}}</td><td>{{.NextAction}}</td><td><a href="{{.UILink}}">{{.UILink}}</a></td><td><code>{{.DocsLink}}</code></td><td>{{.DoesNotProve}}</td></tr>{{end}}
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
<div class="hero start-here">
<h2>Agency Operations Cockpit / Start Here</h2>
<p>{{.Cockpit.Boundary}}</p>
<p><a href="/admin/operations.json">Export private cockpit JSON</a> · <a href="/admin/operations/maintenance">Open maintenance center</a></p>
</div>
<section class="card-grid" aria-label="Role-based task entry points">
<section class="card" id="role-entry-agency-evaluator"><h3>I am evaluating an agency</h3><p><strong>First step:</strong> <a href="/admin/operations/setup-wizard">Review setup progress</a>.</p><p><strong>Done when:</strong> local/private next actions and missing evidence gates are clear.</p><p><strong>Boundary:</strong> local review does not prove agency approval, compliance, or public launch.</p></section>
<section class="card" id="role-entry-operations-staff"><h3>I run daily operations</h3><p><strong>First step:</strong> <a href="/admin/operations/realtime">Check today&apos;s realtime state</a>.</p><p><strong>Done when:</strong> stale devices, Vehicle Positions, Trip Updates, Alerts, and feed freshness have next actions.</p><p><strong>Boundary:</strong> local health does not prove SLA, uptime, or production readiness.</p></section>
<section class="card" id="role-entry-technical-helper"><h3>I am helping technically</h3><p><strong>First step:</strong> <a href="/admin/operations/gtfs-workbench">Review current schedule</a>.</p><p><strong>Done when:</strong> active-vs-draft schedule state, validation meaning, import review, and rollback limits are understood.</p><p><strong>Boundary:</strong> browser review does not silently edit GTFS or prove schedule correctness.</p></section>
<section class="card" id="role-entry-release-reviewer"><h3>I am reviewing release state</h3><p><strong>First step:</strong> <a href="/admin/operations/validation-center">Review validation blockers</a>.</p><p><strong>Done when:</strong> <code>needs_review</code>, package/tag blockers, and prepared-only consumer state are visible.</p><p><strong>Boundary:</strong> diagnostics do not create a release-ready claim.</p></section>
<section class="card" id="role-entry-connector-evaluator"><h3>I am evaluating connectors</h3><p><strong>First step:</strong> <a href="/admin/operations/connectors/workbench">Choose connector recipe</a>.</p><p><strong>Done when:</strong> the first local/synthetic safety check and redaction boundary are clear.</p><p><strong>Boundary:</strong> synthetic checks do not prove live vendor or device compatibility.</p></section>
</section>
{{template "firstRunPanel" .FirstRun}}
{{template "contextHelpPanel" .}}
<h2>Setup Progress</h2>
<table><thead><tr><th>ID</th><th>Area</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .Cockpit.SetupProgress}}<tr id="cockpit-progress-{{.ID}}"><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td><a href="{{.AdminLink}}">{{.NextAction}}</a></td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h2>Primary Actions</h2>
<div class="card-grid" aria-label="Primary agency operations actions">
{{range .Cockpit.PrimaryCards}}<section class="card" id="cockpit-card-{{.ID}}">
<h3>{{.Label}}</h3>
<p class="status"><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></p>
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
<tr><th>Configured public route URLs</th><td>{{if .Discovery.Readiness.AllRequiredFeedsListed}}listed{{else}}missing or incomplete{{end}}</td></tr>
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
<tr><td>GTFS Workbench</td><td>{{.GTFSWorkbench.ActiveFeedVersion.Status}} active schedule; {{.GTFSWorkbench.Import.Status}} import history</td><td>{{formatTime .GTFSWorkbench.GeneratedAt}}</td><td><a href="/admin/operations/gtfs-workbench">review schedule workbench</a> · <a href="/admin/operations/gtfs-workbench.json">export JSON</a></td></tr>
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

<h2>Configured Feed URLs</h2>
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
{{range .Launchpad.Sections}}<tr><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td><ul>{{range .NextActions}}<li>{{.}}</li>{{end}}</ul></td><td>{{range .CommandSuggestions}}<code>{{.}}</code><br>{{end}}</td><td>{{range .AdminLinks}}<a href="{{.}}">{{.}}</a><br>{{end}}</td><td>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>

<h3>Decision Gate</h3>
<table><thead><tr><th>Option</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .Launchpad.DecisionNotes}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.Boundary}}</td></tr>{{end}}
</tbody></table>
<p class="muted">No POST action exists for this page. Missing data remains missing or unknown until the underlying private source records change.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "setup-wizard"}}
{{template "layoutStart" .}}
<h2>Agency Setup</h2>
<p class="warning">Set up the agency profile, schedule data, feed links, validation, and first telemetry review. This private page helps operators prepare the system; it does not publish feeds, contact outside services, or prove compliance.</p>
<p><a href="/admin/operations/setup-wizard.json">Export private setup JSON</a> · <a href="/admin/operations/setup">Open advanced setup details</a> · <a href="/admin/operations/checklist">Open private checklist</a></p>
<div class="card-grid" aria-label="Agency setup progress">
<section class="card">
<h3>Setup Progress</h3>
<p><span class="status-chip status-{{statusClass .SetupWizard.Summary.Status}}">{{.SetupWizard.Summary.Status}}</span></p>
<p>{{.SetupWizard.Summary.CompletedStages}} of {{.SetupWizard.Counts.Stages}} stages complete. {{.SetupWizard.Summary.NeedsReviewStages}} need review, {{.SetupWizard.Summary.MissingStages}} are missing, {{.SetupWizard.Summary.BlockedStages}} are blocked, and {{.SetupWizard.Summary.UnknownStages}} are unknown.</p>
<p class="muted">{{.SetupWizard.Summary.Meaning}}</p>
</section>
<section class="card">
<h3>Next Best Step</h3>
<p><strong>{{.SetupWizard.Summary.NextStageLabel}}</strong></p>
<p>{{.SetupWizard.Summary.NextAction}}</p>
{{if .SetupWizard.Summary.NextActionLink}}<p><a href="{{.SetupWizard.Summary.NextActionLink}}">{{.SetupWizard.Summary.NextActionLink}}</a></p>{{end}}
</section>
<section class="card">
<h3>Advanced Details</h3>
<p>Use the advanced setup page to edit safe publication metadata as an admin, then return here to review the setup sequence.</p>
<p><a href="/admin/operations/setup#publication-metadata">Review publication metadata</a></p>
</section>
</div>
{{if .SetupWizard.Blockers}}
<h3>Review Blocks And Next Actions</h3>
<table><thead><tr><th>Stage</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Console</th></tr></thead><tbody>
{{range .SetupWizard.Blockers}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{if .AdminLink}}<a href="{{.AdminLink}}">{{.ActionLabel}}</a>{{end}}</td></tr>{{end}}
</tbody></table>
{{end}}
<h3>Setup Diagnostics</h3>
<table><thead><tr><th>Diagnostic</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .SetupWizard.Diagnostics}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
<h3>Role Visibility</h3>
<table><thead><tr><th>Capability</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .SetupWizard.RoleVisibility}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
<h3>Technical Helper Cards</h3>
<div class="card-grid" aria-label="Technical helper escalation cards">
{{range .SetupWizard.TechnicalHelp}}<section class="card">
<h4>{{.Label}}</h4>
<p><strong>When needed:</strong> {{.WhenNeeded}}</p>
<p><strong>Next action:</strong> {{.NextAction}}</p>
{{if .AdminLink}}<p><a href="{{.AdminLink}}">Open console area</a></p>{{end}}
{{if .DocsLink}}<p><code>{{.DocsLink}}</code></p>{{end}}
<p class="muted">{{.ClaimBoundary}}</p>
</section>{{end}}
</div>
<div class="card-grid" aria-label="Agency setup stages">
{{range .SetupWizard.Stages}}<section class="card">
<h3>{{.Label}}</h3>
<p><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></p>
<p><strong>What we see:</strong> {{.CurrentSignal}}</p>
<p><strong>Next step:</strong> {{.PrimaryAction}}</p>
{{if .AdminLink}}<p><a href="{{.AdminLink}}">{{.ActionLabel}}</a></p>{{end}}
<p class="muted">{{.ClaimBoundary}}</p>
</section>{{end}}
</div>
<details>
<summary>What this page does not do</summary>
<p>{{.SetupWizard.Boundary}}</p>
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
</details>
<h3>Detailed Setup Signals</h3>
<table><thead><tr><th>ID</th><th>Stage</th><th>Status</th><th>Current signal</th><th>Primary action</th><th>Console</th><th>Docs</th><th>Boundary</th></tr></thead><tbody>
{{range .SetupWizard.Stages}}<tr><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td>{{.PrimaryAction}}</td><td>{{if .AdminLink}}<a href="{{.AdminLink}}">{{.ActionLabel}}</a>{{end}}</td><td>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
<p class="muted">This wizard is GET-only. It does not upload GTFS, mutate setup state, run validators, contact external systems, or create public routes.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "connectors"}}
{{template "layoutStart" .}}
<h2>Connector Hub</h2>
<p class="warning">{{.ConnectorHub.Boundary}}</p>
<div class="card-grid" aria-label="Connector empty or blocked state guidance">
<section class="card empty-state">
<h3>If no connector setup exists</h3>
<p><strong>What am I seeing?</strong> Connector Hub shows safe local adapter shapes, committed example manifests, and any registry diagnostics.</p>
<p><strong>Is this bad?</strong> No for a first browser review. It is a blocker only when a deployment depends on an external telemetry, prediction, validator, monitoring, or discovery integration.</p>
<p><strong>What should I do next?</strong> Pick the connector category, read its boundary, then run the fixed offline connector tests before any authorized integration work.</p>
<p><strong>Can I do it in the browser?</strong> You can review categories, docs, and manifest diagnostics here; the browser does not load plugins or contact systems.</p>
<p><strong>When do I need a technical helper?</strong> Use one to run conformance commands, configure sidecars, map credentials, or prepare deployment-owned external connections.</p>
<p><strong>What this does not prove:</strong> Connector guidance does not prove vendor compatibility, hardware certification, consumer acceptance, compliance, hosted operation, SLA coverage, production readiness, or ETA quality.</p>
</section>
</div>
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

{{define "connector-workbench"}}
{{template "layoutStart" .}}
<h2>Connector Workbench</h2>
<p class="warning">{{.ConnectorWorkbench.Boundary}}</p>
<div class="card-grid" aria-label="Connector workbench boundary guidance">
<section class="card empty-state">
<h3>Start with a recipe</h3>
<p><strong>What am I seeing?</strong> The Workbench turns common connector situations into local/synthetic review paths.</p>
<p><strong>Is this bad?</strong> No. It is intentionally a planning and review surface, not a browser runner.</p>
<p><strong>What should I do next?</strong> Choose the recipe closest to your source, read its first safe check, then use an operator shell for the fixed command.</p>
<p><strong>Can I do it in the browser?</strong> You can review recipes and manifests here; the browser does not run checks, contact external systems, or send telemetry.</p>
<p><strong>When do I need a technical helper?</strong> Use one when credentials, sidecars, deployment-owned webhook receivers, or off-host validators are involved.</p>
<p><strong>What this does not prove:</strong> Workbench review does not prove real integration, compatibility, hardware certification, compliance, consumer acceptance, hosted service, SLA, production readiness, or ETA quality.</p>
</section>
</div>
<table><tbody>
<tr><th><code>backend_command_execution_enabled</code></th><td>{{.ConnectorWorkbench.ClaimFlags.BackendCommandExecutionEnabled}}</td></tr>
<tr><th><code>browser_network_send_enabled</code></th><td>{{.ConnectorWorkbench.ClaimFlags.BrowserNetworkSendEnabled}}</td></tr>
<tr><th><code>manifest_command_execution_enabled</code></th><td>{{.ConnectorWorkbench.ClaimFlags.ManifestCommandExecutionEnabled}}</td></tr>
<tr><th><code>dynamic_backend_plugin_loading_enabled</code></th><td>{{.ConnectorWorkbench.ClaimFlags.DynamicBackendPluginLoadingEnabled}}</td></tr>
<tr><th><code>external_network_contacted</code></th><td>{{.ConnectorWorkbench.ClaimFlags.ExternalNetworkContacted}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.ConnectorWorkbench.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.ConnectorWorkbench.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.ConnectorWorkbench.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>hardware_certification_claimed</code></th><td>{{.ConnectorWorkbench.ClaimFlags.HardwareCertificationClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.ConnectorWorkbench.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>production_grade_eta_claimed</code></th><td>{{.ConnectorWorkbench.ClaimFlags.ProductionGradeETAClaimed}}</td></tr>
</tbody></table>
<h3>Recipe Chooser</h3>
<div class="card-grid" aria-label="Connector recipe chooser">
{{range .ConnectorWorkbench.Recipes}}
<section class="card">
<h4>{{.Label}}</h4>
<p class="status">Status: {{.Status}}</p>
<p>{{.OperatorStory}}</p>
<p><strong>What this is:</strong> {{.WhatThisIs}}</p>
<p><strong>What you need:</strong> {{join .WhatYouNeed ", "}}</p>
<p><strong>Runs where:</strong> {{.RunsWhere}}</p>
<p><strong>First safe check:</strong> <code>{{.FirstSafeCheck}}</code></p>
<p><strong>Good result:</strong> {{.GoodResult}}</p>
<p><strong>If it fails:</strong> {{.IfItFails}}</p>
<p><strong>Does not prove:</strong> {{.DoesNotProve}}</p>
{{if .ManifestIDs}}<p><strong>Example manifests:</strong> {{range .ManifestIDs}}<code>{{.}}</code> {{end}}</p>{{end}}
{{if .AdminLinks}}<p><strong>Console:</strong> {{range .AdminLinks}}<a href="{{.}}">{{.}}</a> {{end}}</p>{{end}}
{{if .DocsLinks}}<p><strong>Docs:</strong> {{range .DocsLinks}}<code>{{.}}</code> {{end}}</p>{{end}}
</section>
{{end}}
</div>
<h3>Dry-Run Command Cards</h3>
<p>These are fixed operator-shell instructions. The browser does not execute them, read command output, send telemetry, or contact external systems.</p>
<table><thead><tr><th>Dry run</th><th>Instruction</th><th>Runs where</th><th>Inputs</th><th>Expected result</th><th>If it fails</th><th>Does not prove</th><th>Docs</th></tr></thead><tbody>
{{range .ConnectorWorkbench.DryRunCommands}}
<tr>
<td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td>
<td><code>{{.CommandLine}}</code></td>
<td>{{.RunsWhere}}</td>
<td>{{.Inputs}}</td>
<td>{{.ExpectedResult}}</td>
<td>{{.FailureNextAction}}</td>
<td>{{.DoesNotProve}}</td>
<td>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</td>
</tr>
{{end}}
</tbody></table>
<h3>Synthetic Telemetry Normalization Preview</h3>
<p class="warning">{{.ConnectorWorkbench.TelemetryPreview.Boundary}}</p>
<table><thead><tr><th>Source</th><th>Fixture</th><th>Status</th><th>Rows</th><th>Expected events</th><th>Expected drops</th><th>Instruction</th><th>Does not prove</th></tr></thead><tbody>
{{range .ConnectorWorkbench.TelemetryPreview.Sources}}
<tr>
<td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td>
<td><code>{{.FixturePath}}</code><br>synthetic only: {{.SyntheticOnly}}</td>
<td>{{.Status}}</td>
<td>{{.ObservedRows}}</td>
<td>{{.ExpectedEvents}}</td>
<td>{{.ExpectedDrops}}</td>
<td><code>{{.CommandLine}}</code></td>
<td>{{.DoesNotProve}}</td>
</tr>
{{end}}
</tbody></table>
<table><thead><tr><th>Source</th><th>Device / vehicle</th><th>Observed</th><th>Quality</th><th>Outcome</th><th>Reason</th><th>Dry run</th><th>Network send</th></tr></thead><tbody>
{{range .ConnectorWorkbench.TelemetryPreview.Rows}}
<tr>
<td><code>{{.SourceID}}</code></td>
<td><code>{{.DeviceID}}</code><br><code>{{.VehicleID}}</code></td>
<td>{{.ObservedAt}}</td>
<td>{{.Quality}}</td>
<td>{{.Outcome}}</td>
<td>{{.Reason}}</td>
<td>{{.DryRun}}</td>
<td>{{.NetworkSend}}</td>
</tr>
{{else}}
<tr><td colspan="8">No synthetic telemetry preview rows were loaded. Review fixture status and run the dry-run commands from an operator shell.</td></tr>
{{end}}
</tbody></table>
<p class="muted">Preview counts: {{.ConnectorWorkbench.TelemetryPreview.Counts.Events}} events, {{.ConnectorWorkbench.TelemetryPreview.Counts.Drops}} drops, network send enabled: {{.ConnectorWorkbench.TelemetryPreview.Counts.NetworkSendEnabled}}.</p>
<h3>{{.ConnectorWorkbench.WebhookBoundary.Title}}</h3>
<p class="warning">{{.ConnectorWorkbench.WebhookBoundary.Boundary}}</p>
<table><thead><tr><th>Boundary</th><th>What it means</th><th>Allowed inputs</th><th>Blocked inputs</th><th>First safe check</th><th>Fail-closed rule</th><th>Redaction rule</th><th>Does not prove</th><th>Review</th></tr></thead><tbody>
{{range .ConnectorWorkbench.WebhookBoundary.Rows}}
<tr>
<td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td>
<td>{{.WhatThisMeans}}</td>
<td>{{range .AllowedInputs}}{{.}}<br>{{end}}</td>
<td>{{range .BlockedInputs}}{{.}}<br>{{end}}</td>
<td><code>{{.FirstSafeCheck}}</code></td>
<td>{{.FailClosedRule}}</td>
<td>{{.RedactionRule}}</td>
<td>{{.DoesNotProve}}</td>
<td>{{range .ReviewLinks}}<a href="{{.}}">{{.}}</a><br>{{end}}</td>
</tr>
{{end}}
</tbody></table>
<p><strong>Docs:</strong> {{range .ConnectorWorkbench.WebhookBoundary.DocsLinks}}<code>{{.}}</code> {{end}}</p>
<h3>{{.ConnectorWorkbench.PredictionGuide.Title}}</h3>
<p class="warning">{{.ConnectorWorkbench.PredictionGuide.Boundary}}</p>
<table><thead><tr><th>Mode</th><th>Status</th><th>What this is</th><th>Inputs</th><th>Outputs</th><th>Failure behavior</th><th>First safe check</th><th>Does not prove</th><th>Review</th></tr></thead><tbody>
{{range .ConnectorWorkbench.PredictionGuide.Rows}}
<tr>
<td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td>
<td>{{.Status}}</td>
<td>{{.WhatThisIs}}</td>
<td>{{range .Inputs}}{{.}}<br>{{end}}</td>
<td>{{range .Outputs}}{{.}}<br>{{end}}</td>
<td>{{.FailureBehavior}}</td>
<td><code>{{.FirstSafeCheck}}</code></td>
<td>{{.DoesNotProve}}</td>
<td>{{range .ReviewLinks}}<a href="{{.}}">{{.}}</a><br>{{end}}{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</td>
</tr>
{{end}}
</tbody></table>
<p><strong>Docs:</strong> {{range .ConnectorWorkbench.PredictionGuide.DocsLinks}}<code>{{.}}</code> {{end}}</p>
<h3>{{.ConnectorWorkbench.MonitoringGuide.Title}}</h3>
<p class="warning">{{.ConnectorWorkbench.MonitoringGuide.Boundary}}</p>
<table><thead><tr><th>Recipe</th><th>Status</th><th>What this is</th><th>Inputs</th><th>Outputs</th><th>Failure behavior</th><th>First safe check</th><th>Does not prove</th><th>Review</th></tr></thead><tbody>
{{range .ConnectorWorkbench.MonitoringGuide.Rows}}
<tr>
<td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td>
<td>{{.Status}}</td>
<td>{{.WhatThisIs}}</td>
<td>{{range .Inputs}}{{.}}<br>{{end}}</td>
<td>{{range .Outputs}}{{.}}<br>{{end}}</td>
<td>{{.FailureBehavior}}</td>
<td><code>{{.FirstSafeCheck}}</code></td>
<td>{{.DoesNotProve}}</td>
<td>{{range .ReviewLinks}}<a href="{{.}}">{{.}}</a><br>{{end}}{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</td>
</tr>
{{end}}
</tbody></table>
<p><strong>Docs:</strong> {{range .ConnectorWorkbench.MonitoringGuide.DocsLinks}}<code>{{.}}</code> {{end}}</p>
<h3>Synthetic Conformance Viewer</h3>
<p class="warning">{{.ConnectorWorkbench.Conformance.Boundary}}</p>
<table><tbody>
<tr><th>Suite</th><td><code>{{.ConnectorWorkbench.Conformance.SuitePath}}</code></td></tr>
<tr><th>Status</th><td>{{.ConnectorWorkbench.Conformance.Status}}</td></tr>
<tr><th>Synthetic only</th><td>{{.ConnectorWorkbench.Conformance.SyntheticOnly}}</td></tr>
<tr><th>Manifest count</th><td>{{.ConnectorWorkbench.Conformance.ManifestCount}}</td></tr>
<tr><th>Case count</th><td>{{.ConnectorWorkbench.Conformance.CaseCount}}</td></tr>
</tbody></table>
<h4>Runner Guidance</h4>
<table><thead><tr><th>Command</th><th>Instruction</th><th>Inputs</th><th>Expected result</th><th>If it fails</th><th>Does not prove</th><th>Docs</th></tr></thead><tbody>
{{range .ConnectorWorkbench.Conformance.RunnerCommands}}
<tr>
<td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td>
<td><code>{{.CommandLine}}</code></td>
<td>{{.Inputs}}</td>
<td>{{.ExpectedResult}}</td>
<td>{{.FailureNextAction}}</td>
<td>{{.DoesNotProve}}</td>
<td>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</td>
</tr>
{{end}}
</tbody></table>
{{range .ConnectorWorkbench.Conformance.Groups}}
<h4>{{.Label}}</h4>
<p>Status: {{.Status}} · cases: {{.CaseCount}}</p>
{{$group := .}}
<table><thead><tr><th>Group</th><th>Required scenarios</th><th>Command</th><th>Case</th><th>Scenario</th><th>Fixture</th><th>Expected</th><th>Assertions</th><th>Status</th><th>Does not prove</th></tr></thead><tbody>
{{range .Cases}}
<tr>
<td><code>{{$group.ID}}</code></td>
<td>{{join $group.RequiredScenarios ", "}}</td>
<td><code>{{$group.CommandLine}}</code></td>
<td><code>{{.ID}}</code></td>
<td>{{.Scenario}}</td>
<td><code>{{.FixturePath}}</code><br>synthetic only: {{.SyntheticOnly}}</td>
<td>{{.ExpectedOutcome}}</td>
<td>{{join .Assertions ", "}}</td>
<td>{{.Status}}</td>
<td>{{$group.DoesNotProve}}</td>
</tr>
{{end}}
</tbody></table>
{{end}}
<h3>{{.ConnectorWorkbench.ManifestReview.Title}}</h3>
<p>{{.ConnectorWorkbench.ManifestReview.Summary}}</p>
<p><strong>Safe plugin definition:</strong> {{.ConnectorWorkbench.ManifestReview.PluginDefinition}}</p>
{{if .ConnectorWorkbench.ManifestReview.Diagnostics}}
<table><thead><tr><th>Level</th><th>Code</th><th>Path</th><th>Message</th></tr></thead><tbody>
{{range .ConnectorWorkbench.ManifestReview.Diagnostics}}<tr><td>{{.Level}}</td><td><code>{{.Code}}</code></td><td><code>{{.Path}}</code></td><td>{{.Message}}</td></tr>{{end}}
</tbody></table>
{{end}}
<table><thead><tr><th>Manifest</th><th>Type</th><th>Mode</th><th>Safety</th><th>Contracts</th><th>Conformance</th><th>Boundary</th><th>Docs</th></tr></thead><tbody>
{{range .ConnectorWorkbench.ManifestReview.Rows}}
<tr>
<td><strong>{{.DisplayName}}</strong><br><code>{{.ConnectorID}}</code><br><code>{{.SourcePath}}</code></td>
<td>{{.ConnectorType}}</td>
<td>{{.Mode}}<br>{{if .DisabledByDefault}}disabled by default{{else}}review before use{{end}}</td>
<td>{{if .FailClosed}}fail closed{{else}}review failure behavior{{end}}<br>secret storage: {{.SecretStorage}}</td>
<td>inputs: {{join .InputContracts ", "}}<br>outputs: {{join .OutputContracts ", "}}</td>
<td>{{.ConformanceCaseCount}} synthetic cases<br><code>{{.FirstCheck}}</code></td>
<td>{{.Boundary}}<br>{{.DoesNotProve}}</td>
<td><code>{{.DocsLink}}</code></td>
</tr>
{{else}}
<tr><td colspan="8">No committed connector manifests were loaded. Review diagnostics and run <code>make external-connection-check</code>.</td></tr>
{{end}}
</tbody></table>
<p><a href="/admin/operations/connectors">Open Connector Hub</a> for category overview. <a href="/admin/operations/connectors/tests">Open Connector Tests</a> for fixed offline commands.</p>
<p class="muted">Connector Workbench is GET-only generated guidance. It does not upload manifests, load backend plugins, execute commands, launch sidecars, contact external parties, write evidence, or change consumer status.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "connector-tests"}}
{{template "layoutStart" .}}
<h2>Connector Test Instructions</h2>
<p class="warning">{{.ConnectorTests.Boundary}}</p>
<div class="card-grid" aria-label="Connector test empty or blocked state guidance">
<section class="card empty-state">
<h3>If connector tests are not run yet</h3>
<p><strong>What am I seeing?</strong> This page lists fixed local/offline commands and the synthetic inputs each command checks.</p>
<p><strong>Is this bad?</strong> Not for review. It becomes a blocker when connector work depends on unverified manifests, examples, or adapter cases.</p>
<p><strong>What should I do next?</strong> Copy the relevant command into an operator terminal, review any failure, and keep the fix inside synthetic fixtures or adapter boundaries.</p>
<p><strong>Can I do it in the browser?</strong> No. The browser only shows commands and boundaries; it does not execute commands or read command output.</p>
<p><strong>When do I need a technical helper?</strong> Use one to run Go/Make checks, diagnose adapter failures, or set up local toolchains.</p>
<p><strong>What this does not prove:</strong> Passing local synthetic checks does not prove real vendor compatibility, consumer acceptance, compliance, production readiness, external network behavior, or ETA quality.</p>
</section>
</div>
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

{{define "gtfs-workbench"}}
{{template "layoutStart" .}}
<h2>GTFS Workbench</h2>
<p class="warning">{{.GTFSWorkbench.Boundary}}</p>
<p><a href="/admin/operations/gtfs-workbench.json">Export private GTFS Workbench JSON</a> · <a href="/admin/operations/gtfs-import">Import Schedule ZIP</a> · <a href="/admin/gtfs-studio">Open Draft Schedule Editor</a> · <a href="/admin/operations/gtfs-quality">Open Schedule Quality</a> · <a href="/admin/operations/validation-health">Open Schedule Validation</a></p>
<div class="card-grid" aria-label="GTFS Workbench summary cards">
<section class="card">
<h3>Current Schedule</h3>
<p><span class="status-chip status-{{statusClass .GTFSWorkbench.ActiveFeedVersion.Status}}">{{.GTFSWorkbench.ActiveFeedVersion.Status}}</span></p>
<p><strong>Active feed version:</strong> {{if .GTFSWorkbench.ActiveFeedVersion.FeedVersionID}}<code>{{.GTFSWorkbench.ActiveFeedVersion.FeedVersionID}}</code>{{else}}not recorded{{end}}</p>
<p><strong>Revision time:</strong> {{formatTimePtr .GTFSWorkbench.ActiveFeedVersion.RevisionTimestamp}}</p>
<p><strong>Public schedule URL:</strong> {{if .GTFSWorkbench.ActiveFeedVersion.CanonicalPublicURL}}<code>{{.GTFSWorkbench.ActiveFeedVersion.CanonicalPublicURL}}</code>{{else}}not configured{{end}}</p>
<p><strong>Current signal:</strong> {{.GTFSWorkbench.ActiveFeedVersion.CurrentSignal}}</p>
<p><strong>Next action:</strong> {{.GTFSWorkbench.ActiveFeedVersion.NextAction}}</p>
<p class="muted">{{.GTFSWorkbench.ActiveFeedVersion.ClaimBoundary}}</p>
</section>
<section class="card">
<h3>Latest Import</h3>
<p><span class="status-chip status-{{statusClass .GTFSWorkbench.Import.Status}}">{{.GTFSWorkbench.Import.Status}}</span></p>
{{with .GTFSWorkbench.Import.Latest}}
<p><strong>Source file:</strong> <code>{{.SourceName}}</code></p>
<p><strong>Source checksum:</strong> <code>{{.SourceSHA256}}</code></p>
<p><strong>Source size:</strong> {{.SourceByteText}}</p>
<p><strong>Feed version:</strong> {{if .FeedVersionID}}<code>{{.FeedVersionID}}</code>{{else}}not linked{{end}}</p>
<p><strong>Import signal:</strong> {{.Signal}}</p>
{{else}}
<p>No GTFS import record is available.</p>
{{end}}
<p><strong>Next action:</strong> {{.GTFSWorkbench.Import.NextAction}}</p>
<p class="muted">{{.GTFSWorkbench.Import.ClaimBoundary}}</p>
</section>
<section class="card">
<h3>Quality And Validation</h3>
<p><strong>Schedule Quality:</strong> <span class="status-chip status-{{statusClass .GTFSWorkbench.Quality.Status}}">{{.GTFSWorkbench.Quality.Status}}</span></p>
<p>Canonical MobilityData static validator: {{.GTFSWorkbench.Quality.CanonicalStatus}}. Internal importer: {{.GTFSWorkbench.Quality.InternalImporterStatus}}.</p>
<p><strong>Schedule Validation:</strong> <span class="status-chip status-{{statusClass .GTFSWorkbench.ValidationHealth.Status}}">{{.GTFSWorkbench.ValidationHealth.Status}}</span></p>
<p>{{.GTFSWorkbench.ValidationHealth.NextAction}}</p>
<p class="muted">Internal import checks and MobilityData static validator diagnostics remain separate signals.</p>
</section>
<section class="card">
<h3>Feed Output</h3>
<p><span class="status-chip status-{{statusClass .GTFSWorkbench.FeedOutput.Status}}">{{.GTFSWorkbench.FeedOutput.Status}}</span></p>
<p><strong>Schedule ZIP URL:</strong> {{if .GTFSWorkbench.FeedOutput.ScheduleURL}}<code>{{.GTFSWorkbench.FeedOutput.ScheduleURL}}</code>{{else}}not configured{{end}}</p>
<p><strong>feeds.json URL:</strong> {{if .GTFSWorkbench.FeedOutput.FeedsJSONURL}}<code>{{.GTFSWorkbench.FeedOutput.FeedsJSONURL}}</code>{{else}}not configured{{end}}</p>
<p><strong>Next action:</strong> {{.GTFSWorkbench.FeedOutput.NextAction}}</p>
<p class="muted">{{.GTFSWorkbench.FeedOutput.ClaimBoundary}}</p>
</section>
</div>
<h3>Next Operator Actions</h3>
<table><thead><tr><th>Step</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Console</th><th>Boundary</th></tr></thead><tbody>
{{range .GTFSWorkbench.Actions}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{if .AdminLink}}<a href="{{.AdminLink}}">{{.AdminLink}}</a>{{end}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
<h3>Import Change Signals</h3>
<table><thead><tr><th>Review item</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .GTFSWorkbench.Import.Diff}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
<h3>Active Vs Previous Schedule Comparison</h3>
<p class="warning">{{.GTFSWorkbench.VersionComparison.ClaimBoundary}}</p>
<table><tbody>
<tr><th>Status</th><td><span class="status-chip status-{{statusClass .GTFSWorkbench.VersionComparison.Status}}">{{.GTFSWorkbench.VersionComparison.Status}}</span></td></tr>
<tr><th>Active feed version</th><td>{{if .GTFSWorkbench.VersionComparison.ActiveFeedVersionID}}<code>{{.GTFSWorkbench.VersionComparison.ActiveFeedVersionID}}</code>{{else}}not recorded{{end}}</td></tr>
<tr><th>Previous feed version</th><td>{{if .GTFSWorkbench.VersionComparison.PreviousFeedVersionID}}<code>{{.GTFSWorkbench.VersionComparison.PreviousFeedVersionID}}</code>{{else}}not visible{{end}}</td></tr>
<tr><th>Previous lifecycle</th><td>{{.GTFSWorkbench.VersionComparison.PreviousLifecycleState}}</td></tr>
<tr><th>Previous validation</th><td>{{.GTFSWorkbench.VersionComparison.PreviousValidationStatus}}</td></tr>
<tr><th>Current signal</th><td>{{.GTFSWorkbench.VersionComparison.CurrentSignal}}</td></tr>
<tr><th>Next action</th><td>{{.GTFSWorkbench.VersionComparison.NextAction}}</td></tr>
</tbody></table>
{{if .GTFSWorkbench.VersionComparison.FileDiffs}}
<h4>File-Level Row Count Diff</h4>
<table><thead><tr><th>File</th><th>Status</th><th>Previous rows</th><th>Active rows</th><th>Delta</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .GTFSWorkbench.VersionComparison.FileDiffs}}<tr><td><code>{{.File}}</code></td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.PreviousRows}}</td><td>{{.ActiveRows}}</td><td>{{.DeltaRows}}</td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
{{end}}
{{if .GTFSWorkbench.VersionComparison.EntityDiffs}}
<h4>Route / Stop / Trip / Service Change Summary</h4>
<table><thead><tr><th>Entity</th><th>Status</th><th>Rows</th><th>Added sample</th><th>Removed sample</th><th>Changed sample</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .GTFSWorkbench.VersionComparison.EntityDiffs}}<tr><td>{{.Entity}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.PreviousRows}} -> {{.ActiveRows}}</td><td>{{if .AddedSample}}{{range .AddedSample}}<code>{{.}}</code><br>{{end}}{{else}}none in bounded sample{{end}}</td><td>{{if .RemovedSample}}{{range .RemovedSample}}<code>{{.}}</code><br>{{end}}{{else}}none in bounded sample{{end}}</td><td>{{if .ChangedSample}}{{range .ChangedSample}}<code>{{.}}</code><br>{{end}}{{else}}none in bounded sample{{end}}</td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
{{end}}
{{if .GTFSWorkbench.VersionComparison.ReviewRows}}
<h4>Version Comparison Review Checklist</h4>
<table><thead><tr><th>Review item</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .GTFSWorkbench.VersionComparison.ReviewRows}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
{{end}}
{{if .GTFSWorkbench.Import.History}}
<h3>Recent Import History</h3>
<table><thead><tr><th>ID</th><th>Status</th><th>Feed version</th><th>Source</th><th>Checksum</th><th>Size</th><th>Counts</th><th>Started</th><th>Completed</th></tr></thead><tbody>
{{range .GTFSWorkbench.Import.History}}<tr><td>{{.ID}}</td><td>{{.Status}}</td><td>{{if .FeedVersionID}}<code>{{.FeedVersionID}}</code>{{else}}not linked{{end}}</td><td><code>{{.SourceName}}</code></td><td><code>{{.SourceSHA256Short}}</code></td><td>{{.SourceByteText}}</td><td>{{.ErrorCount}} errors, {{.WarningCount}} warnings, {{.InfoCount}} info</td><td>{{formatTime .StartedAt}}</td><td>{{formatTimePtr .CompletedAt}}</td></tr>{{end}}
</tbody></table>
{{else}}<p class="muted">Recent import history is not available from this runtime.</p>{{end}}
<h3>Draft Publish Review</h3>
<p class="warning">{{.GTFSWorkbench.DraftReview.ClaimBoundary}}</p>
<table><tbody>
<tr><th>Status</th><td><span class="status-chip status-{{statusClass .GTFSWorkbench.DraftReview.Status}}">{{.GTFSWorkbench.DraftReview.Status}}</span></td></tr>
<tr><th>History</th><td>{{.GTFSWorkbench.DraftReview.HistoryStatus}}</td></tr>
<tr><th>Current signal</th><td>{{.GTFSWorkbench.DraftReview.CurrentSignal}}</td></tr>
<tr><th>Next action</th><td>{{.GTFSWorkbench.DraftReview.NextAction}}</td></tr>
</tbody></table>
<h4>Draft Publish Checklist</h4>
<table><thead><tr><th>Review item</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .GTFSWorkbench.DraftReview.Checklist}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
{{if .GTFSWorkbench.DraftReview.Drafts}}<h4>Recent Drafts</h4>
<table><thead><tr><th>Draft</th><th>Status</th><th>Base feed</th><th>Published feed</th><th>Last publish attempt</th><th>Updated</th></tr></thead><tbody>
{{range .GTFSWorkbench.DraftReview.Drafts}}<tr><td><code>{{.ID}}</code><br>{{.Name}}</td><td>{{.Status}}</td><td>{{if .BaseFeedVersionID}}<code>{{.BaseFeedVersionID}}</code>{{else}}not set{{end}}</td><td>{{if .LastPublishedFeedVersionID}}<code>{{.LastPublishedFeedVersionID}}</code>{{else}}not published{{end}}</td><td>{{.LastPublishAttemptID}}</td><td>{{formatTime .UpdatedAt}}</td></tr>{{end}}
</tbody></table>{{else}}<p class="muted">No recent GTFS Studio draft records are available.</p>{{end}}
{{if .GTFSWorkbench.DraftReview.PublishAttempts}}<h4>Recent Draft Publish Attempts</h4>
<table><thead><tr><th>ID</th><th>Draft</th><th>Status</th><th>Feed version</th><th>Counts</th><th>Started</th><th>Completed</th></tr></thead><tbody>
{{range .GTFSWorkbench.DraftReview.PublishAttempts}}<tr><td>{{.ID}}</td><td><code>{{.DraftID}}</code></td><td>{{.Status}}</td><td>{{if .FeedVersionID}}<code>{{.FeedVersionID}}</code>{{else}}not linked{{end}}</td><td>{{.ErrorCount}} errors, {{.WarningCount}} warnings, {{.InfoCount}} info</td><td>{{formatTime .StartedAt}}</td><td>{{formatTimePtr .CompletedAt}}</td></tr>{{end}}
</tbody></table>{{else}}<p class="muted">No recent GTFS Studio publish attempts are available.</p>{{end}}
<p class="muted">The Workbench has no draft publish POST route. Publishing stays in GTFS Studio and requires the existing admin confirmation path.</p>
<h3>Schedule History And Rollback Guidance</h3>
<p class="warning">{{.GTFSWorkbench.ScheduleHistory.ClaimBoundary}}</p>
<table><tbody>
<tr><th>Status</th><td><span class="status-chip status-{{statusClass .GTFSWorkbench.ScheduleHistory.Status}}">{{.GTFSWorkbench.ScheduleHistory.Status}}</span></td></tr>
<tr><th>History</th><td>{{.GTFSWorkbench.ScheduleHistory.HistoryStatus}}</td></tr>
<tr><th>Current signal</th><td>{{.GTFSWorkbench.ScheduleHistory.CurrentSignal}}</td></tr>
<tr><th>Next action</th><td>{{.GTFSWorkbench.ScheduleHistory.NextAction}}</td></tr>
</tbody></table>
<h4>Rollback Guidance</h4>
<table><thead><tr><th>Review item</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .GTFSWorkbench.ScheduleHistory.RollbackGuidance}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
{{if .GTFSWorkbench.ScheduleHistory.FeedVersions}}<h4>Recent Feed Versions</h4>
<table><thead><tr><th>Feed version</th><th>Source</th><th>Lifecycle</th><th>Active</th><th>Validation</th><th>Published</th><th>Activated</th><th>Retired</th><th>Created</th></tr></thead><tbody>
{{range .GTFSWorkbench.ScheduleHistory.FeedVersions}}<tr><td><code>{{.ID}}</code></td><td>{{.SourceType}}</td><td>{{.LifecycleState}}</td><td>{{.IsActive}}</td><td>{{.ValidationStatus}}</td><td>{{formatTimePtr .PublishedAt}}</td><td>{{formatTimePtr .ActivatedAt}}</td><td>{{formatTimePtr .RetiredAt}}</td><td>{{formatTime .CreatedAt}}</td></tr>{{end}}
</tbody></table>{{else}}<p class="muted">No recent feed-version rows are available for rollback review.</p>{{end}}
<p class="muted">This Workbench is read-only. It does not execute rollback, alter active feeds, or write rollback evidence.</p>
<h3>Preview Tables</h3>
<table><tbody>
<tr><th>Status</th><td>{{.GTFSWorkbench.Preview.Status}}</td></tr>
<tr><th>Row cap</th><td>{{.GTFSWorkbench.Preview.RowLimit}} rows per table</td></tr>
<tr><th>Current signal</th><td>{{.GTFSWorkbench.Preview.CurrentSignal}}</td></tr>
<tr><th>Next action</th><td>{{.GTFSWorkbench.Preview.NextAction}}</td></tr>
<tr><th>Boundary</th><td>{{.GTFSWorkbench.Preview.ClaimBoundary}}</td></tr>
</tbody></table>
<section class="review-tools" data-review-tools data-review-target="gtfs-preview-sections" aria-label="GTFS preview review tools">
<h3>Preview filters</h3>
<div class="review-controls">
<label for="gtfs-preview-filter">Show <select id="gtfs-preview-filter" data-review-filter><option value="all">All</option><option value="needs_action">Needs action</option><option value="blocked">Blocked</option><option value="optional">Optional</option><option value="ok">OK</option></select></label>
<label for="gtfs-preview-search">Find <input id="gtfs-preview-search" data-review-search placeholder="File, table, route, stop, service"></label>
<label for="gtfs-preview-sort">Sort by <select id="gtfs-preview-sort" data-review-sort><option value="needs_action">Needs action first</option><option value="name">Name A-Z</option><option value="status">Status A-Z</option></select></label>
<label><input type="checkbox" data-review-remember> Remember these review settings on this device</label>
<button type="button" data-review-reset>Reset review settings</button>
</div>
<p id="gtfs-preview-status" class="review-status" aria-live="polite" data-review-status>Showing all private preview rows.</p>
<p class="muted">Filters only change this browser view. They do not import, edit, publish, run validators, create evidence, contact external systems, or change consumer status.</p>
</section>
<h4>Required File Checklist</h4>
<table id="gtfs-preview-sections"><thead><tr><th>Item</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .GTFSWorkbench.Preview.RequiredFiles}}<tr data-review-row data-review-status="{{.Status}}" data-review-name="{{.File}}"><td><code>{{.File}}</code></td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
{{range .GTFSWorkbench.Preview.Sections}}<tr data-review-row data-review-status="{{.Status}}" data-review-name="{{.Label}}"><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}} Showing {{.RowsShown}} of {{.TotalRows}} rows; {{.OverflowCount}} omitted by cap.</td><td>Review the bounded rows below, then fix the source GTFS or draft data if the signal needs action.</td><td>{{.ClaimBoundary}}</td></tr>{{end}}
</tbody></table>
{{if .GTFSWorkbench.Preview.Agency}}<h4>Agency Preview</h4>
<table><thead><tr><th>Agency ID</th><th>Name</th><th>Timezone</th></tr></thead><tbody>
{{range .GTFSWorkbench.Preview.Agency}}<tr><td><code>{{.AgencyID}}</code></td><td>{{.Name}}</td><td>{{.Timezone}}</td></tr>{{end}}
</tbody></table>{{end}}
{{if .GTFSWorkbench.Preview.Routes}}<h4>Routes Preview</h4>
<table><thead><tr><th>Route ID</th><th>Short name</th><th>Long name</th><th>Type</th></tr></thead><tbody>
{{range .GTFSWorkbench.Preview.Routes}}<tr><td><code>{{.ID}}</code></td><td>{{.ShortName}}</td><td>{{.LongName}}</td><td>{{.RouteType}}</td></tr>{{end}}
</tbody></table>{{end}}
{{if .GTFSWorkbench.Preview.Stops}}<h4>Stops Preview</h4>
<table><thead><tr><th>Stop ID</th><th>Name</th><th>Lat</th><th>Lon</th></tr></thead><tbody>
{{range .GTFSWorkbench.Preview.Stops}}<tr><td><code>{{.ID}}</code></td><td>{{.Name}}</td><td>{{.Lat}}</td><td>{{.Lon}}</td></tr>{{end}}
</tbody></table>{{end}}
{{if .GTFSWorkbench.Preview.Trips}}<h4>Trips Preview</h4>
<table><thead><tr><th>Trip ID</th><th>Route</th><th>Service</th><th>Block</th><th>Shape</th><th>Direction</th></tr></thead><tbody>
{{range .GTFSWorkbench.Preview.Trips}}<tr><td><code>{{.ID}}</code></td><td><code>{{.RouteID}}</code></td><td><code>{{.ServiceID}}</code></td><td>{{.BlockID}}</td><td>{{.ShapeID}}</td><td>{{.DirectionID}}</td></tr>{{end}}
</tbody></table>{{end}}
{{if .GTFSWorkbench.Preview.Calendar}}<h4>Calendar / Service Preview</h4>
<table><thead><tr><th>Service ID</th><th>Days</th><th>Start</th><th>End</th></tr></thead><tbody>
{{range .GTFSWorkbench.Preview.Calendar}}<tr><td><code>{{.ServiceID}}</code></td><td>{{.Days}}</td><td>{{.StartDate}}</td><td>{{.EndDate}}</td></tr>{{end}}
</tbody></table>{{end}}
{{if .GTFSWorkbench.Preview.Frequencies}}<h4>Frequencies Preview</h4>
<table><thead><tr><th>Trip ID</th><th>Start</th><th>End</th><th>Headway seconds</th><th>Exact times</th></tr></thead><tbody>
{{range .GTFSWorkbench.Preview.Frequencies}}<tr><td><code>{{.TripID}}</code></td><td>{{.StartTime}}</td><td>{{.EndTime}}</td><td>{{.HeadwaySecs}}</td><td>{{.ExactTimes}}</td></tr>{{end}}
</tbody></table>{{end}}
<h3>Claim Flags</h3>
<table><tbody>
<tr><th><code>automatic_gtfs_edit_enabled</code></th><td>{{.GTFSWorkbench.ClaimFlags.AutomaticGTFSEditEnabled}}</td></tr>
<tr><th><code>schedule_published_from_workbench</code></th><td>{{.GTFSWorkbench.ClaimFlags.SchedulePublishedFromWorkbench}}</td></tr>
<tr><th><code>validator_run_from_workbench</code></th><td>{{.GTFSWorkbench.ClaimFlags.ValidatorRunFromWorkbench}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.GTFSWorkbench.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.GTFSWorkbench.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.GTFSWorkbench.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>agency_approval_claimed</code></th><td>{{.GTFSWorkbench.ClaimFlags.AgencyApprovalClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.GTFSWorkbench.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th><code>final_root_readiness_claimed</code></th><td>{{.GTFSWorkbench.ClaimFlags.FinalRootReadinessClaimed}}</td></tr>
<tr><th><code>public_launch_claimed</code></th><td>{{.GTFSWorkbench.ClaimFlags.PublicLaunchClaimed}}</td></tr>
<tr><th><code>hosted_saas_claimed</code></th><td>{{.GTFSWorkbench.ClaimFlags.HostedSaaSClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.GTFSWorkbench.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.GTFSWorkbench.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>production_grade_eta_claimed</code></th><td>{{.GTFSWorkbench.ClaimFlags.ProductionGradeETAClaimed}}</td></tr>
</tbody></table>
<p class="muted">No POST action exists for this Workbench page. It does not import, publish, run validators, edit drafts, create evidence, contact external systems, or change consumer statuses.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "gtfs-import"}}
{{template "layoutStart" .}}
<h2>Browser GTFS Import</h2>
<p class="warning">Private admin-only import path. Raw GTFS ZIP bytes are written to temporary runtime storage for the import attempt and then removed. This page creates no retained evidence, contacts no consumers, records no agency approval, and makes no CAL-ITP/Caltrans compliance, public launch, hosted-service, vendor compatibility, production-readiness, or production-grade ETA claim.</p>
{{if .GTFSImportNotice}}<p class="ok">{{.GTFSImportNotice}}</p>{{end}}
{{if .GTFSImportError}}<p class="bad">{{.GTFSImportError}}</p>{{end}}
<div class="card-grid" aria-label="GTFS import empty or blocked state guidance">
<section class="card empty-state">
<h3>If no GTFS is imported yet</h3>
<p><strong>What am I seeing?</strong> The active schedule row is missing until an imported or published GTFS feed version exists.</p>
<p><strong>Is this bad?</strong> It is normal on first run, but it blocks useful feed health, validation, telemetry matching, and realtime review.</p>
<p><strong>What should I do next?</strong> Import a GTFS ZIP or safe GTFS URL, then review GTFS Quality, Validation Health, and Feed Health.</p>
<p><strong>Can I do it in the browser?</strong> Yes, admins can use ZIP upload or safe URL import on this page.</p>
<p><strong>When do I need a technical helper?</strong> Use one for local app startup, very large scripted imports, unavailable import service, staged comparisons, or rollback work.</p>
<p><strong>What this does not prove:</strong> A successful import does not prove canonical validator success, schedule correctness, agency approval, consumer acceptance, compliance, final-root readiness, hosted operation, or production readiness.</p>
</section>
</div>
<h3>Source Review Before Import</h3>
<table><tbody>
<tr><th>Allowed source</th><td>Use a GTFS ZIP upload from the operator workstation or a safe HTTP(S) URL. The browser form does not allow local/private URLs unless the runtime explicitly enables local testing overrides.</td></tr>
<tr><th>What changes</th><td>A successful published import updates the active schedule feed version through the existing importer. This is not a preview-only action.</td></tr>
<tr><th>Review before submit</th><td>Confirm the source file, expected agency identity, service period, route/stops/trips coverage, license/contact metadata, and rollback plan before an admin starts the import.</td></tr>
<tr><th>Safety controls</th><td>Import requires an admin role, CSRF protection, form size limits, temporary runtime storage, server-owned import paths, and bounded result rendering.</td></tr>
<tr><th>Technical helper needed</th><td>Use a technical helper for large files, scripted imports, staged diffing, rollback execution, source-permission uncertainty, or validator tooling failures.</td></tr>
</tbody></table>
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
<p class="muted">Browser import accepts a ZIP upload or a safe HTTP(S) URL, runs the existing importer, and stores only normal import and validation records. Private/local URLs are blocked unless the runtime explicitly enables local testing overrides. After import, review GTFS quality, validator health, and all five configured feed paths before moving the dataset to wider operator review.</p>
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

{{define "validation-center"}}
{{template "layoutStart" .}}
<h2>Feed Health And Validation Center</h2>
<p class="warning">{{.ValidationCenter.Boundary}}</p>
<p><a href="/admin/operations/validation-center.json">Export private center JSON</a> · <a href="/admin/operations/feed-health">Open feed health</a> · <a href="/admin/operations/validation-health">Open validator health</a> · <a href="/admin/operations/gtfs-quality">Open GTFS quality</a> · <a href="/admin/operations/readiness">Open readiness checklist</a></p>
<div class="card-grid" aria-label="Validation center summary">
<section class="card">
<h3>Feed Rows</h3>
<p class="status"><span class="status-chip status-unknown">{{.ValidationCenter.Counts.FeedRows}}</span></p>
<p>Five private feed rows are expected: feed discovery, static schedule, Vehicle Positions, Trip Updates, and Alerts.</p>
</section>
<section class="card">
<h3>Validator Rows</h3>
<p class="status"><span class="status-chip status-unknown">{{.ValidationCenter.Counts.ValidationRows}}</span></p>
<p>Static and realtime validator rows stay separate from internal importer quality signals.</p>
</section>
<section class="card">
<h3>Issue Drilldowns</h3>
<p class="status"><span class="status-chip status-unknown">{{.ValidationCenter.Counts.IssueRows}}</span></p>
<p>Grouped issues include likely owner, affected files, safe fix path, and verification guidance without raw validator samples.</p>
</section>
<section class="card">
<h3>Prepared Tracker</h3>
<p class="status"><span class="status-chip status-unknown">{{.ValidationCenter.Counts.ConsumerRows}}</span></p>
<p>Prepared packet records remain prepared only. This page does not submit, review, or contact any target.</p>
</section>
</div>
<div class="card-grid" aria-label="Feed health and validation explanation">
<section class="card empty-state">
<h3>Feed Health vs Validation</h3>
<p><strong>Feed health</strong> is the private route and freshness signal available from configured feed records, reliability snapshots, and public path metadata.</p>
<p><strong>Validation</strong> is the server-owned validator tooling and latest result signal for static GTFS or GTFS-Realtime artifacts.</p>
<p><strong>GTFS quality</strong> summarizes static validator and importer notices into operator guidance without editing schedule data.</p>
<p><strong>What this does not prove:</strong> These diagnostics do not prove compliance, consumer acceptance, final-root readiness, hosted availability, SLA coverage, uptime, public launch, production readiness, vendor compatibility, hardware certification, or production-grade ETA quality.</p>
</section>
</div>
<h3>Five Feed URL Panel</h3>
<table><thead><tr><th>Feed</th><th>Public path</th><th>Configured URL</th><th>Status</th><th>HTTP</th><th>Content type</th><th>Freshness</th><th>Validator</th><th>Health</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .ValidationCenter.FeedRows}}<tr><td>{{.Label}}</td><td><code>{{.PublicPath}}</code></td><td><code>{{.ConfiguredURL}}</code></td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.HTTPStatus}}</td><td>{{.ContentType}}</td><td>{{.Freshness}}</td><td>{{.ValidatorState}}</td><td>{{.HealthState}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Validation History</h3>
<table><thead><tr><th>Feed</th><th>Validator</th><th>Status</th><th>Tooling</th><th>Artifact</th><th>Latest result</th><th>Latest at</th><th>Active feed</th><th>Result feed</th><th>Stale</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .ValidationCenter.ValidationHistory}}<tr><td>{{.Label}}</td><td>{{.ValidatorID}}<br><span class="muted">{{.ValidatorName}}</span></td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.ToolingStatus}}</td><td>{{.ArtifactStatus}}</td><td>{{.LatestResultStatus}}</td><td>{{.LatestResultAt}}</td><td>{{.ActiveFeedVersionID}}</td><td>{{.LatestResultFeedVersionID}}</td><td>{{.StaleStatus}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Validator Health</h3>
<table><thead><tr><th>Feed</th><th>Health status</th><th>Current signal</th><th>What this means</th></tr></thead><tbody>
{{range .ValidationCenter.ValidatorHealth}}<tr><td>{{.Label}}</td><td>{{.HealthStatus}}</td><td>{{.CurrentSignal}}</td><td>{{.WhatThisMeans}}</td></tr>{{end}}
</tbody></table>
<h3>GTFS Quality Summary</h3>
<table><thead><tr><th>Source</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .ValidationCenter.GTFSQuality}}<tr><td><a href="{{.DetailsURL}}">{{.Label}}</a></td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Issue Drilldowns</h3>
{{if .ValidationCenter.IssueDrilldowns}}
<table><thead><tr><th>Source</th><th>Severity</th><th>Family</th><th>Codes</th><th>Count</th><th>Sample count</th><th>Likely owner</th><th>Affected files</th><th>Operator summary</th><th>Why it matters</th><th>Recommended action</th><th>Safe fix path</th><th>Verify with</th><th>Escalate if</th><th>Boundary</th></tr></thead><tbody>
{{range .ValidationCenter.IssueDrilldowns}}<tr><td><a href="{{.DetailsURL}}">{{.SourceLabel}}</a></td><td><span class="status-chip status-{{statusClass .Status}}">{{.Severity}}</span></td><td>{{.Family}}</td><td>{{join .Codes ", "}}</td><td>{{.Count}}</td><td>{{.SampleCount}}{{if .OverflowCount}}; {{.OverflowCount}} omitted{{end}}</td><td>{{.LikelyOwner}}</td><td>{{.AffectedFiles}}</td><td>{{.OperatorSummary}}</td><td>{{.WhyItMatters}}</td><td>{{.RecommendedAction}}</td><td>{{.SafeFixPath}}</td><td>{{.VerifyWith}}</td><td>{{.EscalateIf}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
{{else}}
<p class="muted">No grouped GTFS quality issue drilldowns are available. Open GTFS Quality when a source-specific validator or importer result needs review.</p>
{{end}}
<h3>Readiness Timeline</h3>
<table><thead><tr><th>Step</th><th>Status</th><th>Current signal</th><th>What this means</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .ValidationCenter.ReadinessTimeline}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.WhatThisMeans}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Current Blockers</h3>
{{if .ValidationCenter.Blockers}}
<table><thead><tr><th>Severity</th><th>Area</th><th>Signal</th><th>Next action</th><th>Review</th><th>Boundary</th></tr></thead><tbody>
{{range .ValidationCenter.Blockers}}<tr><td><span class="status-chip status-{{statusClass .Severity}}">{{.Severity}}</span></td><td>{{.Area}}</td><td>{{.Signal}}</td><td>{{.NextAction}}</td><td>{{if .ReviewURL}}<a href="{{.ReviewURL}}">Review</a>{{else}}not linked{{end}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
{{else}}
<p class="muted">No current blocker rows are present in the private Center summary. This is not a release-readiness, compliance, consumer acceptance, SLA, uptime, or public-launch claim.</p>
{{end}}
<h3>Prepared Consumer Tracker</h3>
<table><thead><tr><th>Target</th><th>Status</th><th>Source</th><th>Updated</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .ValidationCenter.ConsumerTracker}}<tr><td>{{.Target}}</td><td>{{.Status}}</td><td>{{.Source}}</td><td>{{.UpdatedAt}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Claim Flags</h3>
<table><tbody>
<tr><th><code>external_evidence_created</code></th><td>{{.ValidationCenter.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>final_root_evidence_created</code></th><td>{{.ValidationCenter.ClaimFlags.FinalRootEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.ValidationCenter.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.ValidationCenter.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.ValidationCenter.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>agency_approval_claimed</code></th><td>{{.ValidationCenter.ClaimFlags.AgencyApprovalClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.ValidationCenter.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th><code>public_launch_claimed</code></th><td>{{.ValidationCenter.ClaimFlags.PublicLaunchClaimed}}</td></tr>
<tr><th><code>hosted_saas_claimed</code></th><td>{{.ValidationCenter.ClaimFlags.HostedSaaSClaimed}}</td></tr>
<tr><th><code>sla_claimed</code></th><td>{{.ValidationCenter.ClaimFlags.SLAClaimed}}</td></tr>
<tr><th><code>uptime_guarantee_claimed</code></th><td>{{.ValidationCenter.ClaimFlags.UptimeGuaranteeClaimed}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.ValidationCenter.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>hardware_certification_claimed</code></th><td>{{.ValidationCenter.ClaimFlags.HardwareCertificationClaimed}}</td></tr>
<tr><th><code>production_grade_eta_claimed</code></th><td>{{.ValidationCenter.ClaimFlags.ProductionGradeETAClaimed}}</td></tr>
</tbody></table>
<p class="muted">Validation Center is read-only. It reuses existing private summaries and does not run validators, execute commands, mutate feeds, move consumer statuses, write evidence, create releases, or contact external services.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "feed-health"}}
{{template "layoutStart" .}}
<h2>Feed Health Dashboard</h2>
<p class="warning">{{.FeedHealth.Boundary}}</p>
<p>This command center tracks exactly five configured public route paths: <code>/public/feeds.json</code>, <code>/public/gtfs/schedule.zip</code>, <code>/public/gtfsrt/vehicle_positions.pb</code>, <code>/public/gtfsrt/trip_updates.pb</code>, and <code>/public/gtfsrt/alerts.pb</code>.</p>
<div class="card-grid" aria-label="Feed health empty or blocked state guidance">
<section class="card empty-state">
<h3>If feed rows are missing or blocked</h3>
<p><strong>What am I seeing?</strong> Five expected feed paths are shown with the private metadata, validator, reliability, and freshness signals available right now.</p>
<p><strong>Is this bad?</strong> Missing is expected before setup or import; blocked means a feed, validator artifact, or reliability signal needs operator attention before stronger review.</p>
<p><strong>What should I do next?</strong> Use each row's next action, then return here after GTFS import, validator health, telemetry, or reliability records change.</p>
<p><strong>Can I do it in the browser?</strong> You can review the rows and follow private console links; some fixes happen in GTFS Import, Setup, Validation Health, Devices, or Alerts.</p>
<p><strong>When do I need a technical helper?</strong> Use one for deployment public-root configuration, off-host validation, loopback feed checks, reliability snapshots, or service diagnostics.</p>
<p><strong>What this does not prove:</strong> Local feed rows do not prove consumer acceptance, final-root ownership, compliance, hosted availability, service-level guarantees, uptime, production readiness, or public launch.</p>
</section>
</div>
<p><a href="/admin/operations/feed-health.json">Export private feed health JSON</a> · <a href="/admin/operations/validation-health">Open validator health</a> · <a href="/admin/operations/reliability">Open reliability diagnostics</a></p>
<section class="review-tools" data-review-tools data-review-target="feed-health-review-rows" aria-label="Review tools">
<h3>Review tools</h3>
<div class="review-controls">
<label for="feed-health-review-filter">Show <select id="feed-health-review-filter" data-review-filter><option value="all">All</option><option value="needs_action">Needs action</option><option value="blocked">Blocked</option><option value="missing">Missing</option><option value="stale">Stale</option><option value="unknown">Unknown</option><option value="fresh">Fresh</option></select></label>
<label for="feed-health-review-search">Find <input id="feed-health-review-search" data-review-search placeholder="Feed, status, or next action"></label>
<label for="feed-health-review-sort">Sort by <select id="feed-health-review-sort" data-review-sort><option value="needs_action">Needs action first</option><option value="name">Name A-Z</option><option value="status">Status A-Z</option></select></label>
<label><input type="checkbox" data-review-remember> Remember these review settings on this device</label>
<button type="button" data-review-reset>Reset review settings</button>
</div>
<p id="feed-health-review-status" class="review-status" aria-live="polite" data-review-status>Showing all private diagnostic rows.</p>
<p class="muted">Filters only change this browser view. They do not run validators, change feeds, create evidence, contact consumers, or prove readiness.</p>
</section>
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
<div class="card-grid" id="feed-health-review-rows" aria-label="Plain-language feed health rows">
{{range .FeedHealth.Rows}}
<section class="card" data-review-row data-review-status="{{.Status}}" data-review-name="{{.Label}}" data-review-updated="{{.LastChecked}}">
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
<p class="warning">{{.FeedReadiness.Boundary}}</p>
{{if .DiscoveryError}}<p class="warning">No feed metadata is available. Next action: publish or import a GTFS feed, then bootstrap publication metadata.</p>{{else}}
<h3>Configured feed URL review</h3>
<div class="feed-copy-grid" aria-label="Configured feed URL review">
{{range .FeedReadiness.Rows}}<section class="feed-url-card" id="feed-readiness-{{.ID}}" data-copy-card>
<h3>{{.Label}}</h3>
<p><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></p>
<p><strong>Public path:</strong> <code>{{.PublicPath}}</code></p>
<p><strong>Configured URL:</strong></p>
<code class="copy-value" data-copy-value="{{.CopyValue}}">{{.ConfiguredURL}}</code>
<p><strong>Metadata source:</strong> {{.MetadataSource}}</p>
<p><strong>Metadata status:</strong> {{.MetadataStatus}}</p>
<p><strong>Validation context:</strong> {{.ValidationContext}}</p>
<p><strong>Local public fetch context:</strong> {{.PublicFetchContext}}</p>
<p><strong>Meaning:</strong> {{.Meaning}}</p>
<p><strong>Copy guidance:</strong> {{.CopyGuidance}}</p>
<p><strong>Review before sharing:</strong></p>
<ul>{{range .ReviewChecklist}}<li>{{.}}</li>{{end}}</ul>
<p><strong>Does not prove:</strong> {{.DoesNotProve}}</p>
{{if .DocsLink}}<p><strong>Docs:</strong> <code>{{.DocsLink}}</code></p>{{end}}
</section>{{end}}
</div>
<h3>Source-of-truth metadata checklist</h3>
<table><thead><tr><th>Item</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .FeedReadiness.Metadata}}<tr id="feed-readiness-metadata-{{.ID}}"><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Source-of-truth listing guidance</h3>
<table><thead><tr><th>Item</th><th>Status</th><th>Current signal</th><th>Operator step</th><th>Technical helper step</th><th>Docs</th><th>Does not prove</th></tr></thead><tbody>
{{range .FeedReadiness.SourceOfTruth}}<tr id="feed-readiness-source-{{.ID}}"><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.OperatorStep}}</td><td>{{.TechnicalHelperStep}}</td><td>{{if .DocsLink}}<code>{{.DocsLink}}</code>{{end}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Off-host validation guidance</h3>
<table><thead><tr><th>Item</th><th>Status</th><th>Current signal</th><th>Operator step</th><th>Technical helper step</th><th>Docs</th><th>Does not prove</th></tr></thead><tbody>
{{range .FeedReadiness.OffHost}}<tr id="feed-readiness-off-host-{{.ID}}"><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.OperatorStep}}</td><td>{{.TechnicalHelperStep}}</td><td>{{if .DocsLink}}<code>{{.DocsLink}}</code>{{end}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Public docs portal alignment</h3>
<table><thead><tr><th>Item</th><th>Status</th><th>Current signal</th><th>Operator step</th><th>Technical helper step</th><th>Docs</th><th>Does not prove</th></tr></thead><tbody>
{{range .FeedReadiness.DocsPortal}}<tr id="feed-readiness-docs-{{.ID}}"><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.OperatorStep}}</td><td>{{.TechnicalHelperStep}}</td><td>{{if .DocsLink}}<code>{{.DocsLink}}</code>{{end}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Future final-root/evidence checklist</h3>
<table><thead><tr><th>Gate</th><th>Current status</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .FeedReadiness.FutureChecklist}}<tr id="feed-readiness-future-{{.ID}}"><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .CurrentStatus}}">{{.CurrentStatus}}</span></td><td>{{.NextAction}}</td><td>{{.Boundary}}</td></tr>{{end}}
</tbody></table>
<details><summary>Claim flags for this feed readiness review</summary>
<table><tbody>
<tr><th><code>external_evidence_created</code></th><td>{{.FeedReadiness.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>final_root_evidence_created</code></th><td>{{.FeedReadiness.ClaimFlags.FinalRootEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.FeedReadiness.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th><code>consumer_submission_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.ConsumerSubmissionClaimed}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>public_launch_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.PublicLaunchClaimed}}</td></tr>
<tr><th><code>hosted_saas_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.HostedSaaSClaimed}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>hardware_certification_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.HardwareCertificationClaimed}}</td></tr>
<tr><th><code>production_grade_eta_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.ProductionGradeETAClaimed}}</td></tr>
<tr><th><code>real_world_eta_accuracy_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.RealWorldETAAccuracyClaimed}}</td></tr>
<tr><th><code>sla_coverage_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.SLACoverageClaimed}}</td></tr>
<tr><th><code>uptime_guarantee_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.UptimeGuaranteeClaimed}}</td></tr>
<tr><th><code>final_root_readiness_claimed</code></th><td>{{.FeedReadiness.ClaimFlags.FinalRootReadinessClaimed}}</td></tr>
<tr><th><code>external_browser_fetch_performed</code></th><td>{{.FeedReadiness.ClaimFlags.ExternalBrowserFetchPerformed}}</td></tr>
</tbody></table>
</details>
{{template "feedTable" .}}
{{template "tripUpdatesQuality" .}}
<h3>Feed discovery document</h3>
<table><thead><tr><th>Item</th><th>URL</th><th>Validation</th><th>Last checked</th></tr></thead><tbody>
<tr><td>feeds.json</td><td>{{.Discovery.PublicBaseURL}}/public/feeds.json</td><td>not a validator result</td><td>{{formatTime .Discovery.GeneratedAt}}</td></tr>
</tbody></table>
{{end}}
<p class="muted">This view shows private configured metadata and local diagnostic summaries only. Any future retained evidence, final-root review, or consumer-status movement requires separate written authorization.</p>
<p><a href="/admin/operations/feed-health">Open plain-language feed health</a> · <a href="/admin/operations/gtfs-quality">Review GTFS quality triage actions</a> · <a href="/admin/operations/validation-health">Review private validator health diagnostics</a></p>
{{template "layoutEnd" .}}
{{end}}

{{define "gtfs-quality"}}
{{template "layoutStart" .}}
<h2>GTFS Quality Triage</h2>
{{if .GTFSQualityNotice}}<p class="ok">{{.GTFSQualityNotice}}</p>{{end}}
{{if .GTFSQualityError}}<p class="bad">{{.GTFSQualityError}}</p>{{end}}
<p class="warning">Validator output is diagnostics and supporting signal only. It is not consumer acceptance, not CAL-ITP/Caltrans compliance, not an evidence packet, and not production-readiness proof.</p>
<div class="card-grid" aria-label="GTFS quality empty or blocked state guidance">
<section class="card empty-state">
<h3>If quality data is empty or stale</h3>
<p><strong>What am I seeing?</strong> GTFS Quality combines the latest internal importer result, canonical static validator result, active feed version, and operator guidance.</p>
<p><strong>Is this bad?</strong> Empty is expected before GTFS import or validator setup. Blocking errors should be fixed before relying on schedule or realtime outputs.</p>
<p><strong>What should I do next?</strong> Import or publish GTFS first, identify the source-data owner, fix the source, and rerun the allowlisted static validator when available.</p>
<p><strong>Can I do it in the browser?</strong> Admins can rerun the configured static validator from this page after an active schedule exists.</p>
<p><strong>When do I need a technical helper?</strong> Use one for validator installation, source-system export issues, complex calendar/shape/block fixes, or CLI-only imports.</p>
<p><strong>What this does not prove:</strong> Quality guidance and validator output do not prove compliance, consumer acceptance, agency approval, public launch, hosted operation, or production readiness.</p>
</section>
</div>
<table><tbody>
<tr><th>Active schedule feed version</th><td>{{if .ActiveFeedVersion}}<code>{{.ActiveFeedVersion}}</code>{{else}}missing active schedule; next action: import or publish a schedule before rerunning validation{{end}}</td></tr>
<tr><th>Rerun boundary</th><td>Rerun uses only the authenticated agency active published schedule ZIP and the server-side static MobilityData validator mapping.</td></tr>
<tr><th>Guidance boundary</th><td>{{.GTFSQualityGuidance.Boundary}}</td></tr>
</tbody></table>
<h3>Fix Workflow</h3>
<div class="card-grid">
{{range .GTFSQualityGuidance.Workflow}}<section class="card"><h3>{{.Label}}</h3><p>{{.Summary}}</p><p><strong>Next outcome:</strong> {{.NextOutcome}}</p><p class="muted">{{.DoesNotDo}}</p><p>{{range .AdminLinks}}<a href="{{.}}">{{.}}</a><br>{{end}}</p><p>{{range .DocsLinks}}<code>{{.}}</code><br>{{end}}</p></section>{{end}}
</div>
<h3>Fix Planner</h3>
{{with .GTFSQualityGuidance.FixPlanner}}
<section class="card" id="gtfs-quality-fix-planner">
<p><strong>Status:</strong> <span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></p>
<p>{{.Summary}}</p>
<p><strong>Draft suggestion mode:</strong> <code>{{.DraftSuggestionMode}}</code></p>
<p><strong>Displayed rows:</strong> {{.DisplayedRows}} of {{.TotalRows}}{{if .HiddenRows}}; {{.HiddenRows}} hidden by display cap{{end}}</p>
<p class="muted">{{.Boundary}}</p>
</section>
<table><thead><tr><th>Before validation plan</th><th>After validation plan</th></tr></thead><tbody>
<tr><td>{{range .BeforeValidation}}{{.}}<br>{{end}}</td><td>{{range .AfterValidation}}{{.}}<br>{{end}}</td></tr>
</tbody></table>
{{if .Rows}}
<table><thead><tr><th>Severity</th><th>Source</th><th>Family</th><th>Codes</th><th>Count</th><th>Likely owner</th><th>Affected files</th><th>Issue</th><th>Why it matters</th><th>Safe fix suggestion</th><th>Safe draft suggestion</th><th>Draft suggestion record</th><th>Before validation plan</th><th>After validation plan</th><th>Verify with</th><th>Escalate if</th><th>Samples</th><th>Boundary</th></tr></thead><tbody>
{{range .Rows}}<tr><td>{{.Severity}}</td><td>{{.SourceLabel}}</td><td>{{.Family}}</td><td>{{join .Codes ", "}}</td><td>{{.Count}}</td><td>{{.LikelyOwner}}</td><td>{{.AffectedFiles}}</td><td>{{.IssueSummary}}</td><td>{{.WhyItMatters}}</td><td>{{.SafeFixSuggestion}}</td><td>{{.DraftSuggestion}}</td><td>{{.DraftSuggestionRecord}}</td><td>{{.BeforeValidationPlan}}</td><td>{{.AfterValidationPlan}}</td><td>{{.VerifyWith}}</td><td>{{.EscalateIf}}</td><td>{{range .Samples}}<code>{{.}}</code><br>{{end}}</td><td>{{.NoAutoApplyBoundary}}</td></tr>{{end}}
</tbody></table>
{{else}}<p class="warning">No fix planner rows are available yet. Import or publish GTFS and run validation before exporting a checklist.</p>{{end}}
<h3>Private Fix Checklist</h3>
<pre>{{.Checklist}}</pre>
{{end}}
<h3>Claim Flags</h3>
<table><tbody>
<tr><th><code>automatic_gtfs_edit_enabled</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.AutomaticGTFSEditEnabled}}</td></tr>
<tr><th><code>draft_mutation_enabled</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.DraftMutationEnabled}}</td></tr>
<tr><th><code>draft_suggestion_records_created</code></th><td>{{.GTFSQualityGuidance.ClaimFlags.DraftSuggestionRecordsCreated}}</td></tr>
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
<div class="card-grid" aria-label="Validation health empty or blocked state guidance">
<section class="card empty-state">
<h3>If validators are missing or not run</h3>
<p><strong>What am I seeing?</strong> Validation Health shows server-owned validator tooling state, artifact availability, latest results, staleness, and next actions.</p>
<p><strong>Is this bad?</strong> Not on first run. It becomes a blocker when feed review depends on unavailable artifacts, missing tooling, failed validation, or stale reports.</p>
<p><strong>What should I do next?</strong> Import GTFS, confirm feed artifacts exist, then run or review the allowlisted validator health action.</p>
<p><strong>Can I do it in the browser?</strong> Admins can run the configured all-feed validator health action; other roles can review the current state.</p>
<p><strong>When do I need a technical helper?</strong> Use one to install pinned validators, configure off-host validators, inspect server logs, or fix missing artifacts.</p>
<p><strong>What this does not prove:</strong> Validator rows are supporting diagnostics only and do not prove compliance, consumer acceptance, final-root readiness, hosted availability, or production readiness.</p>
</section>
</div>
<div class="card-grid">
<section class="card"><h3>Internal import validation</h3><p>Open Transit RT importer checks required GTFS structure and blocks unsafe activation paths. It helps explain import failures, but it is not the canonical MobilityData validator.</p></section>
<section class="card"><h3>Canonical static validation</h3><p>MobilityData static GTFS validation reviews the active schedule artifact when pinned tooling is installed and the schedule artifact is available.</p></section>
<section class="card"><h3>GTFS-Realtime validation</h3><p>Realtime validation reviews server-owned Vehicle Positions, Trip Updates, and Alerts protobuf artifacts. Browser requests cannot supply commands, paths, argument lists, artifacts, validator binaries, URLs, or timeouts.</p></section>
</div>
<div class="card-grid" aria-label="Private command model">
<section class="card">
<h3>Read-only refresh command</h3>
<p><code>validation_health.refresh</code> recomputes this private summary from existing records and server-owned artifact checks. It writes nothing, changes no public feed output, creates no evidence, and moves no consumer status.</p>
<p><strong>Result statuses:</strong> <code>ok</code>, <code>needs_review</code>, <code>blocked</code>, or <code>failed</code>. These are private workflow outcomes only.</p>
<p><strong>Private JSON route:</strong> <code>POST /admin/operations/validation-health/refresh.json</code></p>
</section>
<section class="card">
<h3>Validator run command boundary</h3>
<p><code>validation_health.run_all</code> remains admin-only. It can store normal <code>validation_report</code> rows where validators run, but it does not change public feeds, create retained evidence, move consumer statuses, or prove compliance.</p>
<p>Browser requests still cannot supply validator IDs, commands, paths, URLs, argument arrays, artifacts, output paths, reports, or timeouts.</p>
</section>
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
<section class="review-tools" aria-label="Review tools">
<h3>Review tools</h3>
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<button type="button" data-admin-refresh="/admin/operations/validation-health/refresh.json" aria-describedby="validation-refresh-status">Refresh validator summary</button>
<p id="validation-refresh-status" class="review-status" aria-live="polite">Reloads existing private records only. It does not run validators, change feeds, create evidence, contact consumers, or prove readiness.</p>
</section>
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

{{define "prediction-lab"}}
{{template "layoutStart" .}}
<h2>Prediction &amp; ETA Lab</h2>
<p class="warning">{{.PredictionLab.Boundary}}</p>
<p><a href="/admin/operations/prediction-lab.json">Export private prediction lab JSON</a> · <a href="/admin/operations/realtime">Open Realtime Center</a> · <a href="/admin/operations/feed-health">Open Feed Health</a></p>
<div class="card-grid" aria-label="Prediction lab summary">
<section class="card">
<h3>Trip Updates Decision</h3>
<p class="status"><span class="status-chip status-{{statusClass .PredictionLab.Summary.Status}}">{{.PredictionLab.Summary.Status}}</span></p>
<p><strong>Current signal:</strong> {{.PredictionLab.Summary.CurrentSignal}}</p>
<p><strong>Adapter:</strong> {{.PredictionLab.Summary.AdapterName}}</p>
<p><strong>Diagnostics:</strong> {{.PredictionLab.Summary.DiagnosticsStatus}} / {{.PredictionLab.Summary.DiagnosticsReason}}</p>
<p><strong>Counts:</strong> {{.PredictionLab.Summary.TripUpdatesEmitted}} emitted; {{.PredictionLab.Summary.EligiblePredictionCandidates}} eligible; {{.PredictionLab.Summary.WithheldCount}} withheld.</p>
<p><strong>Next action:</strong> {{.PredictionLab.Summary.NextAction}}</p>
<p><strong>Does not prove:</strong> {{.PredictionLab.Summary.DoesNotProve}}</p>
</section>
<section class="card">
<h3>Safe Fallback</h3>
<p>{{.PredictionLab.Deterministic.Boundary}}</p>
<p><strong>Status:</strong> {{.PredictionLab.Deterministic.Status}}</p>
<p><strong>Review signal:</strong> {{.PredictionLab.Deterministic.ReviewSignal}}</p>
<p><strong>Next action:</strong> {{.PredictionLab.Deterministic.NextAction}}</p>
<p><strong>Does not prove:</strong> {{.PredictionLab.Deterministic.DoesNotProve}}</p>
</section>
</div>
<h3>Deterministic Predictor Diagnostics</h3>
<table><thead><tr><th>Diagnostic</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .PredictionLab.Deterministic.Rows}}
<tr><td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>
{{end}}
</tbody></table>
<h3>Why ETAs Are Missing</h3>
<table><thead><tr><th>Reason</th><th>Count</th><th>What it means</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .PredictionLab.WithheldReasons}}
<tr><td><strong>{{.Label}}</strong><br><code>{{.Reason}}</code></td><td>{{.Count}}</td><td>{{.WhatItMeans}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>
{{else}}
<tr><td colspan="5">No withheld reason rows are available yet. Send fresh telemetry, confirm an active schedule, then review Realtime Center and Feed Health.</td></tr>
{{end}}
</tbody></table>
<h3>External Predictor Shadow Review</h3>
<p class="warning">{{.PredictionLab.ShadowReview.Boundary}}</p>
<p><strong>Status:</strong> {{.PredictionLab.ShadowReview.Status}} · <strong>Next action:</strong> {{.PredictionLab.ShadowReview.NextAction}}</p>
<table><thead><tr><th>Mode</th><th>Status</th><th>Reason</th><th>Latency</th><th>Count comparison</th><th>Failure behavior</th><th>First safe check</th><th>Does not prove</th></tr></thead><tbody>
{{range .PredictionLab.ShadowReview.Rows}}
<tr><td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td><td>{{.Status}}</td><td>{{.Reason}}</td><td>{{.Latency}}</td><td>{{.CountComparison}}</td><td>{{.FailureBehavior}}</td><td><code>{{.FirstSafeCheck}}</code></td><td>{{.DoesNotProve}}</td></tr>
{{end}}
</tbody></table>
<h3>Backtest Summary</h3>
<p class="warning">{{.PredictionLab.Backtests.Boundary}}</p>
<p><strong>Status:</strong> {{.PredictionLab.Backtests.Status}} · <strong>Cache root:</strong> <code>{{.PredictionLab.Backtests.RootRef}}</code> · {{.PredictionLab.Backtests.Message}}</p>
<table><thead><tr><th>Output</th><th>Status</th><th>Generated</th><th>Inputs</th><th>Coverage</th><th>Error</th><th>Withheld</th><th>Conformance</th><th>Signal</th><th>Does not prove</th></tr></thead><tbody>
{{range .PredictionLab.Backtests.Rows}}
<tr><td><code>{{.OutputRef}}</code></td><td>{{.Status}}<br><code>{{.MaturityGate}}</code></td><td>{{.GeneratedAt}}</td><td>{{.ObservedRecords}} observed<br>{{.PredictionRecords}} predictions<br>{{.GroupCount}} groups</td><td>prediction: {{.PredictionCoverage}}<br>future stop: {{.FutureStopCoverage}}</td><td>MAE: {{.MAEAbsoluteErrorSeconds}}<br>P90: {{.P90AbsoluteErrorSeconds}}</td><td>{{.WithheldByReason}}</td><td>{{.ConformanceSignal}}</td><td>{{.DiagnosticSignal}}</td><td>{{.DoesNotProve}}</td></tr>
{{else}}
<tr><td colspan="10">No aggregate backtest summaries are available yet. Run the fixed local command from an operator shell when synthetic/local backtest review is needed.</td></tr>
{{end}}
</tbody></table>
<h3>Conservative Handling Guide</h3>
<p class="warning">{{.PredictionLab.HandlingGuide.Boundary}}</p>
<p><strong>Status:</strong> {{.PredictionLab.HandlingGuide.Status}} · <strong>Next action:</strong> {{.PredictionLab.HandlingGuide.NextAction}}</p>
<table><thead><tr><th>Situation</th><th>Review signal</th><th>Safe behavior</th><th>Operator step</th><th>Does not prove</th></tr></thead><tbody>
{{range .PredictionLab.HandlingGuide.Rows}}
<tr><td><strong>{{.Situation}}</strong><br><code>{{.ID}}</code></td><td>{{.ReviewSignal}}</td><td>{{.SafeBehavior}}</td><td>{{.OperatorStep}}</td><td>{{.DoesNotProve}}</td></tr>
{{end}}
</tbody></table>
<h3>Future ETA Proof Gates</h3>
<p class="warning">{{.PredictionLab.ProofChecklist.Boundary}}</p>
<p><strong>Status:</strong> {{.PredictionLab.ProofChecklist.Status}} · <strong>Next action:</strong> {{.PredictionLab.ProofChecklist.NextAction}}</p>
<table><thead><tr><th>Gate</th><th>Required review</th><th>Authorization</th><th>Does not prove</th></tr></thead><tbody>
{{range .PredictionLab.ProofChecklist.Rows}}
<tr><td><strong>{{.FutureGate}}</strong><br><code>{{.ID}}</code></td><td>{{.RequiredReview}}</td><td>{{.SeparateAuthorization}}</td><td>{{.DoesNotProve}}</td></tr>
{{end}}
</tbody></table>
<h3>Needs Operator Review</h3>
<table><thead><tr><th>Severity</th><th>Area</th><th>Signal</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .PredictionLab.ReviewRows}}
<tr><td>{{.Severity}}</td><td>{{.Area}}</td><td>{{.Signal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>
{{end}}
</tbody></table>
<h3>Fixed Local Checks</h3>
<p>These commands are operator-shell guidance only. The browser does not run them, capture output, contact sidecars, or create evidence.</p>
<table><thead><tr><th>Check</th><th>Instruction</th><th>Expected result</th><th>Does not prove</th></tr></thead><tbody>
{{range .PredictionLab.Commands}}
<tr><td><strong>{{.Label}}</strong><br><code>{{.ID}}</code></td><td><code>{{.CommandLine}}</code></td><td>{{.ExpectedResult}}</td><td>{{.DoesNotProve}}</td></tr>
{{end}}
</tbody></table>
<h3>Claim Flags</h3>
<table><tbody>
<tr><th><code>browser_predictor_run_enabled</code></th><td>{{.PredictionLab.ClaimFlags.BrowserPredictorRunEnabled}}</td></tr>
<tr><th><code>external_network_contacted</code></th><td>{{.PredictionLab.ClaimFlags.ExternalNetworkContacted}}</td></tr>
<tr><th><code>backend_command_execution_enabled</code></th><td>{{.PredictionLab.ClaimFlags.BackendCommandExecutionEnabled}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.PredictionLab.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.PredictionLab.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>production_grade_eta_claimed</code></th><td>{{.PredictionLab.ClaimFlags.ProductionGradeETAClaimed}}</td></tr>
<tr><th><code>real_world_eta_accuracy_claimed</code></th><td>{{.PredictionLab.ClaimFlags.RealWorldETAAccuracyClaimed}}</td></tr>
</tbody></table>
{{template "layoutEnd" .}}
{{end}}

{{define "realtime"}}
{{template "layoutStart" .}}
<h2>Realtime Operations Center</h2>
<p class="warning">{{.Realtime.Boundary}}</p>
<p><a href="/admin/operations/realtime.json">Export private realtime JSON</a> · <a href="/admin/operations/telemetry">Open telemetry freshness</a> · <a href="/admin/operations/devices">Open device credentials</a> · <a href="/admin/operations/telemetry-simulator">Open simulator guide</a></p>
<div class="card-grid" aria-label="Realtime status summary">
<section class="card">
<h3>Fleet Freshness</h3>
<p class="status"><span class="status-chip status-{{statusClass .Realtime.Summary.Status}}">{{.Realtime.Summary.Status}}</span></p>
<p><strong>Current signal:</strong> {{.Realtime.Summary.CurrentSignal}}</p>
<p><strong>Latest telemetry:</strong> {{.Realtime.Summary.LatestTelemetryRows}} rows; {{.Realtime.Summary.FreshTelemetryRows}} fresh; {{.Realtime.Summary.StaleTelemetryRows}} stale.</p>
<p><strong>Devices:</strong> {{.Realtime.Summary.DeviceBindings}} bindings; {{.Realtime.Summary.DevicesReporting}} reporting; {{.Realtime.Summary.DevicesNotSeen}} not seen.</p>
<p><strong>Next action:</strong> {{.Realtime.Summary.NextAction}}</p>
</section>
<section class="card">
<h3>Assignments</h3>
<p><strong>Matched:</strong> {{.Realtime.Summary.MatchedAssignments}}</p>
<p><strong>Unknown or unavailable:</strong> {{.Realtime.Summary.UnknownAssignments}}</p>
<p><strong>Low confidence:</strong> {{.Realtime.Summary.LowConfidenceRows}}</p>
<p><strong>Manual overrides:</strong> {{.Realtime.Summary.ManualOverrides}}</p>
<p class="muted">Unknown or withheld is a safe conservative state when the system does not have enough evidence.</p>
</section>
{{range .Realtime.Feeds}}
<section class="card">
<h3>{{.Label}}</h3>
<p class="status"><span class="status-chip status-{{statusClass .State}}">{{.State}}</span></p>
<p><strong>Count:</strong> {{.Count}}</p>
<p><strong>Latest signal:</strong> {{.LatestSignal}}</p>
<p><strong>Stale or withheld:</strong> {{.StaleOrWithheld}}</p>
{{if .Adapter}}<p><strong>Adapter:</strong> {{.Adapter}}</p>{{end}}
{{if .Details}}<p><strong>Details:</strong> {{range $index, $detail := .Details}}{{if $index}}; {{end}}{{$detail.Label}}={{$detail.Count}}{{end}}</p>{{end}}
<p><strong>Next action:</strong> <a href="{{.AdminLink}}">{{.NextAction}}</a></p>
<p><strong>Does not prove:</strong> {{.DoesNotProve}}</p>
</section>
{{end}}
</div>
<h3>Realtime Feed Usefulness Review</h3>
<section class="card" id="realtime-usefulness-review">
<p class="status"><span class="status-chip status-{{statusClass .Realtime.Usefulness.Status}}">{{.Realtime.Usefulness.Status}}</span></p>
<p>{{.Realtime.Usefulness.Summary}}</p>
<p class="muted">{{.Realtime.Usefulness.Boundary}}</p>
</section>
<table><thead><tr><th>Feed</th><th>Score</th><th>Current signal</th><th>Helpful signal</th><th>Needs review</th><th>Consumer-safe behavior</th><th>Next action</th><th>Details</th><th>Boundary</th></tr></thead><tbody>
{{range .Realtime.Usefulness.Rows}}<tr><td>{{.Label}}</td><td>{{.Score}} / 3<br><span class="status-chip status-{{statusClass .ScoreLabel}}">{{.ScoreLabel}}</span></td><td>{{.CurrentSignal}}</td><td>{{.HelpfulSignal}}</td><td>{{.NeedsReviewSignal}}</td><td>{{.ConsumerSafeBehavior}}</td><td>{{.NextAction}}</td><td>{{range .Details}}<span class="pill">{{.Label}}: {{.Count}}</span> {{end}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Freshness And Lifecycle Review</h3>
<table><thead><tr><th>Area</th><th>Status</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .Realtime.Usefulness.Freshness}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Consumer-Safe Omission Rules</h3>
<table><thead><tr><th>Condition</th><th>Safe behavior</th><th>Review step</th><th>Boundary</th></tr></thead><tbody>
{{range .Realtime.Usefulness.OmissionRules}}<tr><td>{{.Condition}}</td><td>{{.SafeBehavior}}</td><td>{{.ReviewStep}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Needs Operator Review</h3>
<table><thead><tr><th>Severity</th><th>Area</th><th>Signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .Realtime.Issues}}<tr><td><span class="status-chip status-{{statusClass .Severity}}">{{.Severity}}</span></td><td>{{if .AdminLink}}<a href="{{.AdminLink}}">{{.Area}}</a>{{else}}{{.Area}}{{end}}</td><td>{{.Signal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Realtime Quality Guidance</h3>
<table><thead><tr><th>Topic</th><th>What it means</th><th>Review signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .Realtime.Guidance}}<tr><td>{{.Label}}</td><td>{{.WhatItMeans}}</td><td>{{.ReviewSignal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
{{if .TelemetryError}}<p class="warning">{{.TelemetryError}}. Next action: confirm the telemetry service and database are running.</p>{{end}}
{{if not .Realtime.Fleet}}<p class="warning">No fleet telemetry or device binding rows are visible. Next action: create or rotate a device credential, install it on a device or simulator, then send an authenticated sample telemetry event.</p>{{else}}
<h3>Fleet Freshness And Assignment Overview</h3>
<table><thead><tr><th>Vehicle</th><th>Device</th><th>Freshness</th><th>Observed</th><th>Age seconds</th><th>Assignment</th><th>Route</th><th>Trip</th><th>Confidence</th><th>Reasons</th><th>Current signal</th><th>Next action</th><th>Boundary</th></tr></thead><tbody>
{{range .Realtime.Fleet}}<tr><td>{{.VehicleID}}</td><td>{{.DeviceID}}</td><td><span class="status-chip status-{{statusClass .Freshness}}">{{.Freshness}}</span></td><td>{{.ObservedAt}}</td><td>{{.AgeSeconds}}</td><td>{{if .AssignmentState}}{{.AssignmentState}}{{else}}not available{{end}}{{if .DegradedState}} / {{.DegradedState}}{{end}}{{if .AssignmentSource}}<br><span class="muted">source: {{.AssignmentSource}}</span>{{end}}</td><td>{{.RouteID}}</td><td>{{.TripID}}</td><td>{{.Confidence}}</td><td>{{join .ReasonCodes ", "}}</td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
{{end}}
<h3>Claim Flags</h3>
<table><tbody>
<tr><th><code>browser_telemetry_send_enabled</code></th><td>{{.Realtime.ClaimFlags.BrowserTelemetrySendEnabled}}</td></tr>
<tr><th><code>backend_command_execution_enabled</code></th><td>{{.Realtime.ClaimFlags.BackendCommandExecutionEnabled}}</td></tr>
<tr><th><code>device_token_collected_by_browser</code></th><td>{{.Realtime.ClaimFlags.DeviceTokenCollectedByBrowser}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.Realtime.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.Realtime.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.Realtime.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.Realtime.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>vendor_compatibility_claimed</code></th><td>{{.Realtime.ClaimFlags.VendorCompatibilityClaimed}}</td></tr>
<tr><th><code>hardware_certification_claimed</code></th><td>{{.Realtime.ClaimFlags.HardwareCertificationClaimed}}</td></tr>
<tr><th><code>production_avl_reliability_claimed</code></th><td>{{.Realtime.ClaimFlags.ProductionAVLReliabilityClaimed}}</td></tr>
<tr><th><code>production_grade_eta_claimed</code></th><td>{{.Realtime.ClaimFlags.ProductionGradeETAClaimed}}</td></tr>
<tr><th><code>real_world_eta_accuracy_claimed</code></th><td>{{.Realtime.ClaimFlags.RealWorldETAAccuracyClaimed}}</td></tr>
<tr><th><code>sla_claimed</code></th><td>{{.Realtime.ClaimFlags.SLAClaimed}}</td></tr>
<tr><th><code>public_launch_claimed</code></th><td>{{.Realtime.ClaimFlags.PublicLaunchClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.Realtime.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
</tbody></table>
<p class="muted">Realtime Center is read-only. It does not change <code>/v1/telemetry</code>, public GTFS-Realtime feeds, device credentials, assignments, Alerts, evidence records, consumer tracker states, releases, or external services.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "telemetry"}}
{{template "layoutStart" .}}
<h2>Telemetry Freshness</h2>
<p class="warning">Private telemetry freshness diagnostics only. Viewing this page creates no retained evidence, contacts no vendors or consumers, changes no consumer status, and does not prove hardware certification, vendor compatibility, production AVL reliability, consumer acceptance, compliance, hosted service, or production readiness.</p>
<p>Stale threshold: {{.StaleThreshold}}</p>
<div class="card-grid" aria-label="Telemetry empty or blocked state guidance">
<section class="card empty-state">
<h3>If no telemetry appears</h3>
<p><strong>What am I seeing?</strong> Telemetry Freshness shows accepted latest observations, stale state, and conservative assignment results for this agency.</p>
<p><strong>Is this bad?</strong> It is normal before devices or simulator sends exist, but Vehicle Positions and Trip Updates will remain empty or limited until fresh telemetry arrives.</p>
<p><strong>What should I do next?</strong> Create or rotate a device credential, send a sample through authenticated ingest, then review freshness and Feed Health.</p>
<p><strong>Can I do it in the browser?</strong> You can review telemetry and rotate credentials from the browser; sending telemetry happens through a device or operator-shell simulator.</p>
<p><strong>When do I need a technical helper?</strong> Use one for simulator commands, device networking, ingest target configuration, database-backed matcher diagnostics, or deployment troubleshooting.</p>
<p><strong>What this does not prove:</strong> Fresh telemetry does not prove hardware certification, vendor compatibility, production AVL reliability, consumer acceptance, compliance, hosted operation, or production readiness.</p>
</section>
</div>
<div class="card-grid">
<section class="card"><h3>Make Vehicle Positions non-empty</h3><p>Create or rotate a device token, configure a device or synthetic simulator from an operator shell, send accepted telemetry to <code>/v1/telemetry</code>, then review this page and Feed Health.</p></section>
<section class="card"><h3>Why Trip Updates may be empty</h3><p>Trip Updates can be empty when telemetry is missing or stale, assignment confidence is too low, a vehicle is unknown, or the prediction adapter withholds output. Prefer empty or unknown over false certainty.</p></section>
<section class="card"><h3>Does not prove</h3><p>Fresh or visible telemetry does not prove hardware certification, vendor compatibility, production AVL reliability, consumer acceptance, compliance, hosted service, or production readiness.</p></section>
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
<div class="card-grid" aria-label="Telemetry simulator empty or blocked state guidance">
<section class="card empty-state">
<h3>If simulator scenarios are unavailable</h3>
<p><strong>What am I seeing?</strong> The page is a read-only guide to committed synthetic scenarios and fixed operator-shell commands.</p>
<p><strong>Is this bad?</strong> It is not bad for browser review, but missing fixtures block safe synthetic telemetry practice.</p>
<p><strong>What should I do next?</strong> Restore committed simulator fixtures or copy a fixed command into an operator shell after the local app and credentials are ready.</p>
<p><strong>Can I do it in the browser?</strong> No. The browser shows commands and boundaries only; it does not send telemetry or collect tokens.</p>
<p><strong>When do I need a technical helper?</strong> Use one for local app startup, shell environment setup, seeded credentials, matcher diagnostics, or failed simulator commands.</p>
<p><strong>What this does not prove:</strong> Synthetic telemetry does not prove real fleet reliability, vendor compatibility, hardware certification, consumer acceptance, compliance, hosted operation, or ETA quality.</p>
</section>
</div>
<div class="card-grid">
<section class="card"><h3>Target rules</h3>{{range .TelemetrySimulator.TargetRules}}<p>{{.}}</p>{{end}}</section>
<section class="card"><h3>Credential handling</h3>{{range .TelemetrySimulator.CredentialHandling}}<p>{{.}}</p>{{end}}</section>
<section class="card"><h3>Diagnostics policy</h3><p>{{.TelemetrySimulator.DiagnosticsPolicy}}</p></section>
</div>
<p><strong>First local/synthetic dry-run safety check:</strong> start with a command containing <code>DRY_RUN=true</code>. This previews committed synthetic payload shape only and does not test a live vendor, live AVL API, real device, or public feed consumer.</p>
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
<p class="warning">Private device credential diagnostics only. Viewing or rotating credentials creates no retained evidence, contacts no vendors or consumers, changes no consumer status, and does not prove hardware certification, vendor compatibility, production AVL reliability, consumer acceptance, compliance, hosted service, or production readiness.</p>
<p>The supported browser flow is rotate/rebind. If a device has no credential yet, this uses the existing rebind API path.</p>
<div class="card-grid" aria-label="Device credential empty or blocked state guidance">
<section class="card empty-state">
<h3>If no devices are listed</h3>
<p><strong>What am I seeing?</strong> Device Credentials shows configured device-to-vehicle bindings, token status dates, latest accepted telemetry, and assignment context.</p>
<p><strong>Is this bad?</strong> It is expected before first setup, but it blocks live telemetry and useful Vehicle Positions until at least one credential is installed and reporting.</p>
<p><strong>What should I do next?</strong> Ask an admin to rotate or create the first device token, install it on the device or simulator, then check Telemetry Freshness.</p>
<p><strong>Can I do it in the browser?</strong> Admins can rotate or rebind one-time credentials in the browser; read-only users can review status only.</p>
<p><strong>When do I need a technical helper?</strong> Use one for installing tokens on hardware, configuring device network targets, simulator sends, or diagnosing stale/no telemetry.</p>
<p><strong>What this does not prove:</strong> A device binding does not prove hardware certification, vendor compatibility, production AVL reliability, consumer acceptance, compliance, hosted operation, or production readiness.</p>
</section>
</div>
<div class="card-grid">
<section class="card"><h3>Token status</h3><p>The table shows credential status and dates, never stored token values. New tokens are shown only once after an admin rotate/rebind action.</p></section>
<section class="card"><h3>Vehicle binding</h3><p>Each device row links a device to a vehicle, latest accepted telemetry time, freshness, assignment state, match confidence where available, and a next action.</p></section>
<section class="card"><h3>Realtime setup</h3><p>Vehicle Positions need accepted fresh telemetry. Trip Updates may still be empty until matching confidence and prediction diagnostics justify output.</p></section>
<section class="card"><h3>Does not prove</h3><p>Device bindings and credential rotation do not prove hardware certification, vendor compatibility, production AVL reliability, consumer acceptance, compliance, hosted service, or production readiness.</p></section>
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
<div class="card-grid" aria-label="Maintenance empty or blocked state guidance">
<section class="card empty-state">
<h3>If backup, restore, or service checks are missing</h3>
<p><strong>What am I seeing?</strong> Maintenance summarizes configured and not-configured signals for version, active feed, five-feed checks, validators, backups, restore drills, telemetry, and service diagnostics.</p>
<p><strong>Is this bad?</strong> Missing is expected on a local first run; it becomes a blocker before depending on routine operations, recovery, support, or stronger deployment claims.</p>
<p><strong>What should I do next?</strong> Review each summary row, configure missing private backup/restore values where appropriate, and run the linked diagnostics from the operator environment.</p>
<p><strong>Can I do it in the browser?</strong> You can review status and next steps here; backup/restore configuration and support bundles are operator-shell work.</p>
<p><strong>When do I need a technical helper?</strong> Use one for backup paths, restore-drill targets, service health checks, support bundles, deployment logs, or redaction review.</p>
<p><strong>What this does not prove:</strong> Maintenance rows do not prove SLA coverage, uptime, hosted availability, production readiness, compliance, agency adoption, consumer acceptance, or disaster-recovery success.</p>
</section>
</div>
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
<h3>Local Diagnostic Summaries</h3>
<p class="warning">{{.Maintenance.Diagnostics.Boundary}}</p>
<p><strong>Status:</strong> {{.Maintenance.Diagnostics.Status}}</p>
<table><thead><tr><th>ID</th><th>Source</th><th>Status</th><th>Generated</th><th>Current signal</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .Maintenance.Diagnostics.Rows}}<tr id="maintenance-diagnostic-{{.ID}}"><td><code>{{.ID}}</code></td><td>{{.Label}}<br><code>{{.SourceRef}}</code></td><td>{{.Status}}</td><td>{{.GeneratedAt}}</td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Infrastructure Checks</h3>
<p class="warning">{{.Maintenance.Infrastructure.Boundary}}</p>
<p><strong>Status:</strong> {{.Maintenance.Infrastructure.Status}} · <strong>Next action:</strong> {{.Maintenance.Infrastructure.NextAction}}</p>
<table><thead><tr><th>ID</th><th>Item</th><th>Status</th><th>Current signal</th><th>Operator step</th><th>Technical helper step</th><th>Does not prove</th></tr></thead><tbody>
{{range .Maintenance.Infrastructure.Rows}}<tr id="maintenance-infrastructure-{{.ID}}"><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td>{{.OperatorStep}}</td><td>{{.TechnicalHelperStep}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Backup And Restore Review</h3>
<p class="warning">{{.Maintenance.BackupRestore.Boundary}}</p>
<p><strong>Status:</strong> {{.Maintenance.BackupRestore.Status}} · <strong>Next action:</strong> {{.Maintenance.BackupRestore.NextAction}}</p>
<table><thead><tr><th>ID</th><th>Item</th><th>Status</th><th>Current signal</th><th>Operator step</th><th>Technical helper step</th><th>Does not prove</th></tr></thead><tbody>
{{range .Maintenance.BackupRestore.Rows}}<tr id="maintenance-backup-restore-{{.ID}}"><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td>{{.OperatorStep}}</td><td>{{.TechnicalHelperStep}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Upgrade And Rollback Review</h3>
<p class="warning">{{.Maintenance.UpgradeRollback.Boundary}}</p>
<p><strong>Status:</strong> {{.Maintenance.UpgradeRollback.Status}} · <strong>Next action:</strong> {{.Maintenance.UpgradeRollback.NextAction}}</p>
<table><thead><tr><th>ID</th><th>Item</th><th>Status</th><th>Current signal</th><th>Operator step</th><th>Technical helper step</th><th>Does not prove</th></tr></thead><tbody>
{{range .Maintenance.UpgradeRollback.Rows}}<tr id="maintenance-upgrade-rollback-{{.ID}}"><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td>{{.OperatorStep}}</td><td>{{.TechnicalHelperStep}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Support Bundle And Redaction Review</h3>
<p class="warning">{{.Maintenance.SupportReview.Boundary}}</p>
<p><strong>Status:</strong> {{.Maintenance.SupportReview.Status}} · <strong>Next action:</strong> {{.Maintenance.SupportReview.NextAction}}</p>
<table><thead><tr><th>ID</th><th>Item</th><th>Status</th><th>Current signal</th><th>Operator step</th><th>Technical helper step</th><th>Does not prove</th></tr></thead><tbody>
{{range .Maintenance.SupportReview.Rows}}<tr id="maintenance-support-review-{{.ID}}"><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td>{{.OperatorStep}}</td><td>{{.TechnicalHelperStep}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Maintenance Cadence Plan</h3>
<p class="warning">{{.Maintenance.CadencePlan.Boundary}}</p>
<p><strong>Status:</strong> {{.Maintenance.CadencePlan.Status}} · <strong>Next action:</strong> {{.Maintenance.CadencePlan.NextAction}}</p>
<table><thead><tr><th>ID</th><th>Item</th><th>Status</th><th>Current signal</th><th>Operator step</th><th>Technical helper step</th><th>Does not prove</th></tr></thead><tbody>
{{range .Maintenance.CadencePlan.Rows}}<tr id="maintenance-cadence-{{.ID}}"><td><code>{{.ID}}</code></td><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.CurrentSignal}}</td><td>{{.OperatorStep}}</td><td>{{.TechnicalHelperStep}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
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
<h2>Prepared Consumer Packet Tracker</h2>
<p class="muted">The Phase 20 docs/evidence tracker is the source for prepared packet state. These statuses are not submission, review, acceptance, or ingestion evidence.</p>
{{if .ConsumerError}}<p class="warning">{{.ConsumerError}}. The docs/evidence tracker guidance remains visible below.</p>{{end}}
<section class="panel warning-panel" aria-labelledby="consumer-prepared-only-heading">
<h3 id="consumer-prepared-only-heading">Prepared-Only Consumer Packet Explanation</h3>
<p>{{.ConsumerPreparation.Boundary}}</p>
<p><strong>Status:</strong> {{.ConsumerPreparation.Status}} · <strong>Summary:</strong> {{.ConsumerPreparation.Summary}} · <strong>Runtime records:</strong> {{.ConsumerPreparation.RuntimeRecordCount}}</p>
<p><strong>Operator rule:</strong> {{.ConsumerPreparation.OperatorRule}}</p>
</section>
<table><thead><tr><th>Target</th><th>Docs tracker status</th><th>Source</th><th>Current record</th><th>Packet path</th><th>Notes</th></tr></thead><tbody>
{{range .Consumers}}<tr><td>{{.Name}}</td><td>{{.Status}}</td><td>{{.Source}}</td><td><code>{{.CurrentPath}}</code></td><td><code>{{.PacketPath}}</code></td><td>{{.Notes}}</td></tr>{{end}}
</tbody></table>
<h3>Target Boundary Review</h3>
<table><thead><tr><th>Target</th><th>Status</th><th>Current record</th><th>Packet path</th><th>Meaning</th><th>Next action</th><th>Does not prove</th></tr></thead><tbody>
{{range .ConsumerPreparation.Targets}}<tr id="consumer-prepared-{{.ID}}"><td>{{.Name}}</td><td>{{.Status}}</td><td><code>{{.CurrentPath}}</code></td><td><code>{{.PacketPath}}</code></td><td>{{.Meaning}}</td><td>{{.NextAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Future Authorization Gates</h3>
<table><thead><tr><th>Gate</th><th>Status</th><th>Required before action</th><th>Blocked in this track</th><th>Does not prove now</th></tr></thead><tbody>
{{range .ConsumerPreparation.FutureGates}}<tr id="consumer-gate-{{.ID}}"><td>{{.Label}}</td><td>{{.Status}}</td><td>{{.RequiredAuthorization}}</td><td>{{.BlockedAction}}</td><td>{{.DoesNotProve}}</td></tr>{{end}}
</tbody></table>
<h3>Workflow Separation</h3>
<table><thead><tr><th>Workflow</th><th>Boundary</th><th>Operator handling</th></tr></thead><tbody>
{{range .ConsumerPreparation.Separations}}<tr id="consumer-separation-{{.ID}}"><td>{{.Label}}</td><td>{{.Boundary}}</td><td>{{.OperatorHandling}}</td></tr>{{end}}
</tbody></table>
<details>
<summary>Claim flags for this consumer packet review</summary>
<table><tbody>
<tr><th><code>consumer_statuses_changed</code></th><td>{{.ConsumerPreparation.ClaimFlags.ConsumerStatusesChanged}}</td></tr>
<tr><th><code>consumer_submission_claimed</code></th><td>{{.ConsumerPreparation.ClaimFlags.ConsumerSubmissionClaimed}}</td></tr>
<tr><th><code>consumer_review_claimed</code></th><td>{{.ConsumerPreparation.ClaimFlags.ConsumerReviewClaimed}}</td></tr>
<tr><th><code>consumer_acceptance_claimed</code></th><td>{{.ConsumerPreparation.ClaimFlags.ConsumerAcceptanceClaimed}}</td></tr>
<tr><th><code>consumer_ingestion_claimed</code></th><td>{{.ConsumerPreparation.ClaimFlags.ConsumerIngestionClaimed}}</td></tr>
<tr><th><code>consumer_listing_claimed</code></th><td>{{.ConsumerPreparation.ClaimFlags.ConsumerListingClaimed}}</td></tr>
<tr><th><code>consumer_display_claimed</code></th><td>{{.ConsumerPreparation.ClaimFlags.ConsumerDisplayClaimed}}</td></tr>
<tr><th><code>external_contact_performed</code></th><td>{{.ConsumerPreparation.ClaimFlags.ExternalContactPerformed}}</td></tr>
<tr><th><code>external_evidence_created</code></th><td>{{.ConsumerPreparation.ClaimFlags.ExternalEvidenceCreated}}</td></tr>
<tr><th><code>final_root_evidence_created</code></th><td>{{.ConsumerPreparation.ClaimFlags.FinalRootEvidenceCreated}}</td></tr>
<tr><th><code>compliance_claimed</code></th><td>{{.ConsumerPreparation.ClaimFlags.ComplianceClaimed}}</td></tr>
<tr><th><code>production_readiness_claimed</code></th><td>{{.ConsumerPreparation.ClaimFlags.ProductionReadinessClaimed}}</td></tr>
<tr><th><code>hosted_saas_claimed</code></th><td>{{.ConsumerPreparation.ClaimFlags.HostedSaaSClaimed}}</td></tr>
<tr><th><code>public_launch_claimed</code></th><td>{{.ConsumerPreparation.ClaimFlags.PublicLaunchClaimed}}</td></tr>
</tbody></table>
</details>
<h3>Runtime Deployment Workflow Records</h3>
{{if .RuntimeConsumers}}<table><thead><tr><th>Target</th><th>Runtime status</th><th>Source</th><th>Updated</th><th>Notes</th></tr></thead><tbody>
{{range .RuntimeConsumers}}<tr><td>{{.Name}}</td><td>{{.Status}}</td><td>{{.Source}}</td><td>{{formatTimePtr .UpdatedAt}}</td><td>{{.Notes}}</td></tr>{{end}}
</tbody></table>{{else}}<p class="warning">No runtime consumer workflow records are available. This does not change the docs tracker prepared packet state.</p>{{end}}
<p>Docs tracker repo file path: <code>docs/evidence/consumer-submissions/README.md</code></p>
{{template "layoutEnd" .}}
{{end}}

{{define "evidence"}}
{{template "layoutStart" .}}
<h2>Evidence Links And Runbooks</h2>
<p class="muted">These markdown files are repository file paths, not web routes served by this app. Links are navigation aids and do not prove retained evidence exists.</p>
<table><thead><tr><th>Record</th><th>Repo file path</th><th>Last updated</th></tr></thead><tbody>
{{range .Links}}<tr><td>{{.Label}}</td><td><code>{{.Path}}</code></td><td>{{.UpdatedAt}}</td></tr>{{end}}
</tbody></table>
<p class="muted">These links help operators find repo/deployment evidence. They do not assert consumer acceptance, hosted service availability, agency endorsement, or universal production readiness.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "setup"}}
{{template "layoutStart" .}}
<h2>Advanced Setup Details</h2>
{{if .SetupNotice}}<p class="ok">{{.SetupNotice}}</p>{{end}}
{{if .SetupError}}<p class="bad">{{.SetupError}}</p>{{end}}
<p><a href="/admin/operations/setup-wizard">Return to Agency Setup</a> · <a href="/admin/operations/checklist">Open private operator checklist</a> · <a href="/admin/operations/checklist.json">Export private checklist JSON</a></p>
<p class="muted">Each status is tied to a named source. Missing records stay missing until publication metadata, feed discovery, validation records, device bindings, telemetry, docs tracker records, or evidence links support a stronger statement.</p>
<h3>Setup Diagnostics</h3>
<table><thead><tr><th>Diagnostic</th><th>Status</th><th>Current signal</th><th>Next action</th></tr></thead><tbody>
{{range .SetupWizard.Diagnostics}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td></tr>{{end}}
</tbody></table>
<h3>Role Visibility</h3>
<table><thead><tr><th>Capability</th><th>Status</th><th>Current signal</th><th>Next action</th></tr></thead><tbody>
{{range .SetupWizard.RoleVisibility}}<tr><td>{{.Label}}</td><td><span class="status-chip status-{{statusClass .Status}}">{{.Status}}</span></td><td>{{.CurrentSignal}}</td><td>{{.NextAction}}</td></tr>{{end}}
</tbody></table>
<div class="card-grid" aria-label="Technical helper escalation cards">
{{range .SetupWizard.TechnicalHelp}}<section class="card"><h3>{{.Label}}</h3><p><strong>When needed:</strong> {{.WhenNeeded}}</p><p><strong>Next action:</strong> {{.NextAction}}</p>{{if .AdminLink}}<p><a href="{{.AdminLink}}">Open console area</a></p>{{end}}</section>{{end}}
</div>
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
{{if .IsAdmin}}
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
{{else}}
<p class="warning">Publication metadata changes require an admin role. This account can review setup status but cannot submit setup forms.</p>
{{end}}

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
{{if .IsAdmin}}
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
{{else}}
<p class="warning">Validation runs from this setup page require an admin role. Review-only users can open Validator Health but cannot start allowlisted validators from the browser.</p>
{{end}}
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
