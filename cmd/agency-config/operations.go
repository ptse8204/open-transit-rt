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
	defaultTelemetryLimit = 200
	defaultStaleSeconds   = 90
)

type operationsPage struct {
	Title              string
	AgencyID           string
	GeneratedAt        time.Time
	EnvironmentLabel   string
	CSRFToken          string
	Discovery          compliance.FeedDiscovery
	DiscoveryError     string
	PublicationConfig  compliance.PublicationConfig
	PublicationError   string
	SetupNotice        string
	SetupError         string
	SetupSteps         []setupStepView
	ValidationResult   *compliance.ValidationResult
	ActiveFeedVersion  string
	FeedsUpdatedAt     *time.Time
	TelemetryUpdatedAt *time.Time
	TripUpdatesQuality tripUpdatesQualityView
	ScorecardUpdatedAt *time.Time
	ConsumersUpdatedAt *time.Time
	EvidenceUpdatedAt  string
	Scorecard          *compliance.Scorecard
	ScorecardError     string
	Consumers          []consumerStatusView
	RuntimeConsumers   []consumerStatusView
	ConsumerError      string
	Telemetry          []telemetryView
	TelemetryError     string
	StaleCount         int
	Devices            []devices.Binding
	DeviceError        string
	DeviceToken        string
	DeviceTokenMeta    devices.RebindResult
	Links              []evidenceLink
	Section            string
	StaleThreshold     time.Duration
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

func (h *handler) operationsRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/admin/operations" {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderOperations(w, r, "dashboard")
		return
	}
	trimmed := strings.Trim(strings.TrimPrefix(r.URL.Path, "/admin/operations/"), "/")
	switch trimmed {
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
		StaleThreshold:   staleThreshold(),
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
	page.SetupSteps = setupSteps(page)
	return page
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
			NextAction: "Use CLI GTFS ZIP import for agency ZIPs or GTFS Studio for typed draft authoring and publish.",
			ActionURL:  "/admin/gtfs-studio",
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
			NextAction: "Run one allowlisted validator from this page or the existing admin validation API.",
			ActionURL:  "/admin/operations/setup#validation",
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
}).Parse(`
{{define "layoutStart"}}
<!doctype html><html><head><meta charset="utf-8"><title>{{.Title}}</title>
<style>
body{font-family:system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;margin:2rem;line-height:1.4;color:#1f2933}
nav a{margin-right:1rem} table{border-collapse:collapse;width:100%;margin:1rem 0} th,td{border:1px solid #d8dee4;padding:.45rem;text-align:left;vertical-align:top}
th{background:#f6f8fa}.pill{display:inline-block;border:1px solid #c8d1dc;border-radius:3px;padding:.1rem .35rem;background:#f6f8fa}
.warning{background:#fff8c5}.ok{background:#dafbe1}.bad{background:#ffebe9}.muted{color:#59636e}.token{border:1px solid #f0c36d;background:#fff8c5;padding:1rem}
form{margin:1rem 0} label{display:block;margin:.35rem 0} input,select,textarea{min-width:22rem;max-width:100%;padding:.35rem} button{padding:.4rem .7rem}
</style></head><body>
<h1>{{.Title}}</h1>
<p>Agency: <strong>{{.AgencyID}}</strong> · environment: <span class="pill">{{.EnvironmentLabel}}</span> · generated: {{formatTime .GeneratedAt}}</p>
<nav>
<a href="/admin/operations">Dashboard</a>
<a href="/admin/operations/feeds">Feeds</a>
<a href="/admin/operations/telemetry">Telemetry</a>
<a href="/admin/operations/devices">Devices</a>
<a href="/admin/alerts/console">Alerts</a>
<a href="/admin/operations/consumers">Consumers</a>
<a href="/admin/operations/evidence">Evidence</a>
<a href="/admin/operations/setup">Setup</a>
<a href="/admin/gtfs-studio">GTFS Studio</a>
</nav>
{{end}}
{{define "layoutEnd"}}</body></html>{{end}}

{{define "dashboard"}}
{{template "layoutStart" .}}
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
<tr><td>Feeds / validation</td><td>{{if .DiscoveryError}}not configured{{else}}{{len .Discovery.Feeds}} feed records{{end}}</td><td>{{formatTimePtr .FeedsUpdatedAt}}</td><td><a href="/admin/operations/feeds">review feed URLs and validation</a></td></tr>
<tr><td>Telemetry freshness</td><td>{{if .TelemetryError}}{{.TelemetryError}}{{else}}{{len .Telemetry}} vehicles; {{.StaleCount}} stale{{end}}</td><td>{{formatTimePtr .TelemetryUpdatedAt}}</td><td><a href="/admin/operations/telemetry">inspect vehicle freshness</a></td></tr>
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
{{template "layoutEnd" .}}
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
{{if .TelemetryError}}<p class="warning">{{.TelemetryError}}. Next action: confirm the telemetry service and database are running.</p>{{else if not .Telemetry}}<p class="warning">No telemetry has been accepted yet. Next action: create or rotate a device token, configure the device, then send a sample telemetry event.</p>{{else}}
<table><thead><tr><th>Vehicle</th><th>Device</th><th>Observed</th><th>Age seconds</th><th>Freshness</th><th>Assignment</th><th>Trip</th><th>Route</th><th>Confidence</th><th>Reasons</th><th>Assignment time</th></tr></thead><tbody>
{{range .Telemetry}}<tr><td>{{.VehicleID}}</td><td>{{.DeviceID}}</td><td>{{formatTime .ObservedAt}}</td><td>{{.AgeSeconds}}</td><td>{{if .Stale}}stale{{else}}fresh{{end}}</td><td>{{if .AssignmentState}}{{.AssignmentState}}{{else}}not available{{end}}{{if .DegradedState}} / {{.DegradedState}}{{end}}</td><td>{{.TripID}}</td><td>{{.RouteID}}</td><td>{{.Confidence}}</td><td>{{join .ReasonCodes ", "}}</td><td>{{formatTimePtr .AssignmentAt}}</td></tr>{{end}}
</tbody></table>{{end}}
<p class="muted">Safe diagnostics omit raw telemetry payloads, full score details, token fields, and private debug blobs.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "devices"}}
{{template "layoutStart" .}}
<h2>Device Credentials</h2>
<p class="warning">Device tokens are secrets. Store a one-time token immediately; it will not be shown again by this console.</p>
<p>The current supported browser flow is rotate/rebind. If a device has no credential yet, this uses the existing rebind API path; Phase 18 does not add a separate first-time creation API.</p>
{{if .DeviceToken}}<div class="token"><h3>One-time token</h3><p>Device: {{.DeviceTokenMeta.DeviceID}} · Vehicle: {{.DeviceTokenMeta.VehicleID}} · Rotated: {{.DeviceTokenMeta.RotatedAt}}</p><p><code>{{.DeviceToken}}</code></p></div>{{end}}
{{if .DeviceError}}<p class="warning">{{.DeviceError}}</p>{{end}}
<form method="post" action="/admin/operations/devices">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="agency_id" value="{{.AgencyID}}">
<label>Device ID <input name="device_id" required></label>
<label>Vehicle ID <input name="vehicle_id" required></label>
<label>Reason <input name="reason" placeholder="rotation or rebind reason"></label>
<button>Rotate / rebind token</button>
</form>
{{if not .Devices}}<p class="warning">No device bindings are recorded. Next action: use the rotate/rebind form for the first device and store the returned token securely.</p>{{else}}
<table><thead><tr><th>Device</th><th>Vehicle</th><th>Status</th><th>Valid from</th><th>Last used</th><th>Rotated</th><th>Revoked</th></tr></thead><tbody>
{{range .Devices}}<tr><td>{{.DeviceID}}</td><td>{{.VehicleID}}</td><td>{{.Status}}</td><td>{{formatTime .ValidFrom}}</td><td>{{formatTimePtr .LastUsedAt}}</td><td>{{formatTimePtr .RotatedAt}}</td><td>{{formatTimePtr .RevokedAt}}</td></tr>{{end}}
</tbody></table>{{end}}
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
<p class="muted">These links help operators find repo/deployment evidence. They do not assert consumer acceptance, hosted SaaS availability, agency endorsement, or universal production readiness.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "setup"}}
{{template "layoutStart" .}}
<h2>Guided Setup Checklist</h2>
{{if .SetupNotice}}<p class="ok">{{.SetupNotice}}</p>{{end}}
{{if .SetupError}}<p class="bad">{{.SetupError}}</p>{{end}}
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
<label>Public base URL <input name="public_base_url" maxlength="2048" required value="{{if .PublicationConfig.PublicBaseURL}}{{.PublicationConfig.PublicBaseURL}}{{else}}{{.Discovery.PublicBaseURL}}{{end}}"></label>
<label>Feed base URL <input name="feed_base_url" maxlength="2048" required value="{{.PublicationConfig.FeedBaseURL}}"></label>
<label>Technical contact email <input name="technical_contact_email" maxlength="320" value="{{if .PublicationConfig.TechnicalContactEmail}}{{.PublicationConfig.TechnicalContactEmail}}{{else}}{{.Discovery.TechnicalContactEmail}}{{end}}"></label>
<label>License name <input name="license_name" maxlength="160" value="{{if .PublicationConfig.LicenseName}}{{.PublicationConfig.LicenseName}}{{else}}{{.Discovery.License.Name}}{{end}}"></label>
<label>License URL <input name="license_url" maxlength="2048" value="{{if .PublicationConfig.LicenseURL}}{{.PublicationConfig.LicenseURL}}{{else}}{{.Discovery.License.URL}}{{end}}"></label>
<label>Publication environment <input name="publication_environment" maxlength="64" value="{{publicationEnvValue .}}"></label>
<button>Store publication metadata</button>
</form>

<h2>GTFS Import And Authoring</h2>
<p>Source: feed discovery. Browser ZIP upload is deferred in Phase 26 because upload security, size limits, validation, and role checks need a dedicated design.</p>
<table><tbody>
<tr><th>CLI ZIP import</th><td>Use the existing GTFS import flow documented in <code>docs/tutorials/real-agency-gtfs-onboarding.md</code>.</td></tr>
<tr><th>Typed authoring</th><td><a href="/admin/gtfs-studio">Open GTFS Studio</a> for draft authoring and publish.</td></tr>
<tr><th>Validation triage</th><td>Use <code>docs/tutorials/gtfs-validation-triage.md</code> and the validation form below.</td></tr>
<tr><th>Active feed verification</th><td><a href="/admin/operations/feeds">Review feed discovery and validation records</a>.</td></tr>
</tbody></table>

<h2 id="validation">Validation</h2>
<p class="muted">Source: validation records. The browser chooses only feed type; the server maps it to an allowlisted validator. Validation is supporting evidence only, not consumer acceptance or compliance.</p>
{{if .ValidationResult}}<p class="ok">Last run from this page: {{.ValidationResult.FeedType}} validation {{.ValidationResult.Status}} with {{.ValidationResult.ErrorCount}} errors and {{.ValidationResult.WarningCount}} warnings.</p>{{end}}
<form method="post" action="/admin/operations/setup">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="action" value="run_validation">
<label>Feed type <select name="feed_type">
<option value="schedule">schedule</option>
<option value="vehicle_positions">vehicle_positions</option>
<option value="trip_updates">trip_updates</option>
<option value="alerts">alerts</option>
</select></label>
<button>Run allowlisted validation</button>
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
`))
