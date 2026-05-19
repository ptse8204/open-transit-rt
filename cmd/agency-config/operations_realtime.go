package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/state"
)

type operationsRealtimeView struct {
	GeneratedAt    time.Time                      `json:"generated_at"`
	AgencyID       string                         `json:"agency_id"`
	Boundary       string                         `json:"boundary"`
	StaleThreshold string                         `json:"stale_threshold"`
	Summary        operationsRealtimeSummary      `json:"summary"`
	Feeds          []operationsRealtimeFeedStatus `json:"feeds"`
	Usefulness     operationsRealtimeUsefulness   `json:"usefulness"`
	Publishing     []operationsRealtimeFeedReview `json:"publishing_review"`
	ReplayGuidance operationsRealtimeReplayGuide  `json:"replay_guidance"`
	Fleet          []operationsRealtimeFleetRow   `json:"fleet"`
	Issues         []operationsRealtimeIssue      `json:"issues"`
	Guidance       []operationsRealtimeGuidance   `json:"guidance"`
	ClaimFlags     operationsRealtimeClaimFlags   `json:"claim_flags"`
}

type operationsRealtimeSummary struct {
	Status              string `json:"status"`
	CurrentSignal       string `json:"current_signal"`
	LatestTelemetryRows int    `json:"latest_telemetry_rows"`
	FreshTelemetryRows  int    `json:"fresh_telemetry_rows"`
	StaleTelemetryRows  int    `json:"stale_telemetry_rows"`
	DeviceBindings      int    `json:"device_bindings"`
	DevicesReporting    int    `json:"devices_reporting"`
	DevicesNotSeen      int    `json:"devices_not_seen"`
	MatchedAssignments  int    `json:"matched_assignments"`
	UnknownAssignments  int    `json:"unknown_assignments"`
	LowConfidenceRows   int    `json:"low_confidence_rows"`
	ManualOverrides     int    `json:"manual_overrides"`
	NextAction          string `json:"next_action"`
	DoesNotProve        string `json:"does_not_prove"`
}

type operationsRealtimeUsefulness struct {
	Status        string                              `json:"status"`
	Summary       string                              `json:"summary"`
	Boundary      string                              `json:"boundary"`
	Rows          []operationsRealtimeUsefulnessScore `json:"rows"`
	Freshness     []operationsRealtimeFreshnessReview `json:"freshness"`
	OmissionRules []operationsRealtimeOmissionRule    `json:"omission_rules"`
}

type operationsRealtimeUsefulnessScore struct {
	ID                   string      `json:"id"`
	Label                string      `json:"label"`
	Score                int         `json:"score"`
	ScoreLabel           string      `json:"score_label"`
	CurrentSignal        string      `json:"current_signal"`
	HelpfulSignal        string      `json:"helpful_signal"`
	NeedsReviewSignal    string      `json:"needs_review_signal"`
	ConsumerSafeBehavior string      `json:"consumer_safe_behavior"`
	NextAction           string      `json:"next_action"`
	DoesNotProve         string      `json:"does_not_prove"`
	Details              []countView `json:"details,omitempty"`
}

type operationsRealtimeFeedReview struct {
	ID               string                               `json:"id"`
	Label            string                               `json:"label"`
	Status           string                               `json:"status"`
	WhatLooksHealthy string                               `json:"what_looks_healthy"`
	NeedsAttention   string                               `json:"needs_attention"`
	NotProven        string                               `json:"not_proven"`
	NextAction       string                               `json:"next_action"`
	Signals          []operationsRealtimePublishingSignal `json:"signals"`
}

type operationsRealtimePublishingSignal struct {
	Label   string `json:"label"`
	Value   string `json:"value"`
	Meaning string `json:"meaning"`
}

type operationsRealtimeReplayGuide struct {
	Status       string                         `json:"status"`
	Summary      string                         `json:"summary"`
	BrowserStart string                         `json:"browser_start"`
	LocalReplay  string                         `json:"local_replay"`
	ReviewAfter  string                         `json:"review_after"`
	Boundary     string                         `json:"boundary"`
	Steps        []operationsRealtimeReplayStep `json:"steps"`
}

type operationsRealtimeReplayStep struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	Action       string `json:"action"`
	SafeBoundary string `json:"safe_boundary"`
}

type operationsRealtimeFreshnessReview struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	DoesNotProve  string `json:"does_not_prove"`
}

type operationsRealtimeOmissionRule struct {
	ID           string `json:"id"`
	Condition    string `json:"condition"`
	SafeBehavior string `json:"safe_behavior"`
	ReviewStep   string `json:"review_step"`
	DoesNotProve string `json:"does_not_prove"`
}

type operationsRealtimeFeedStatus struct {
	ID              string      `json:"id"`
	Label           string      `json:"label"`
	State           string      `json:"state"`
	Count           string      `json:"count"`
	LatestSignal    string      `json:"latest_signal"`
	StaleOrWithheld string      `json:"stale_or_withheld"`
	Adapter         string      `json:"adapter,omitempty"`
	NextAction      string      `json:"next_action"`
	AdminLink       string      `json:"admin_link"`
	DoesNotProve    string      `json:"does_not_prove"`
	Details         []countView `json:"details,omitempty"`
}

type operationsRealtimeFleetRow struct {
	VehicleID        string   `json:"vehicle_id"`
	DeviceID         string   `json:"device_id"`
	Freshness        string   `json:"freshness"`
	ObservedAt       string   `json:"observed_at"`
	ReceivedAt       string   `json:"received_at"`
	AgeSeconds       int64    `json:"age_seconds"`
	AssignmentState  string   `json:"assignment_state"`
	DegradedState    string   `json:"degraded_state,omitempty"`
	AssignmentSource string   `json:"assignment_source,omitempty"`
	RouteID          string   `json:"route_id,omitempty"`
	TripID           string   `json:"trip_id,omitempty"`
	Confidence       string   `json:"confidence,omitempty"`
	ReasonCodes      []string `json:"reason_codes,omitempty"`
	CurrentSignal    string   `json:"current_signal"`
	NextAction       string   `json:"next_action"`
	DoesNotProve     string   `json:"does_not_prove"`
}

type operationsRealtimeIssue struct {
	Severity     string `json:"severity"`
	Area         string `json:"area"`
	Signal       string `json:"signal"`
	NextAction   string `json:"next_action"`
	AdminLink    string `json:"admin_link,omitempty"`
	DoesNotProve string `json:"does_not_prove"`
}

type operationsRealtimeGuidance struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	WhatItMeans  string `json:"what_it_means"`
	ReviewSignal string `json:"review_signal"`
	NextAction   string `json:"next_action"`
	DoesNotProve string `json:"does_not_prove"`
}

type operationsRealtimeClaimFlags struct {
	BrowserTelemetrySendEnabled     bool `json:"browser_telemetry_send_enabled"`
	BackendCommandExecutionEnabled  bool `json:"backend_command_execution_enabled"`
	DeviceTokenCollectedByBrowser   bool `json:"device_token_collected_by_browser"`
	ExternalEvidenceCreated         bool `json:"external_evidence_created"`
	ConsumerStatusesChanged         bool `json:"consumer_statuses_changed"`
	ComplianceClaimed               bool `json:"compliance_claimed"`
	ProductionReadinessClaimed      bool `json:"production_readiness_claimed"`
	VendorCompatibilityClaimed      bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed    bool `json:"hardware_certification_claimed"`
	ProductionAVLReliabilityClaimed bool `json:"production_avl_reliability_claimed"`
	ProductionGradeETAClaimed       bool `json:"production_grade_eta_claimed"`
	RealWorldETAAccuracyClaimed     bool `json:"real_world_eta_accuracy_claimed"`
	SLAClaimed                      bool `json:"sla_claimed"`
	PublicLaunchClaimed             bool `json:"public_launch_claimed"`
	ConsumerAcceptanceClaimed       bool `json:"consumer_acceptance_claimed"`
}

func (h *handler) renderRealtime(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "realtime")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "realtime", page)
}

func (h *handler) renderRealtimeJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "realtime")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.Realtime)
}

func buildOperationsRealtime(page operationsPage) operationsRealtimeView {
	fleet := buildRealtimeFleetRows(page)
	summary := buildRealtimeSummary(page, fleet)
	feeds := realtimeFeedStatuses(page)
	return operationsRealtimeView{
		GeneratedAt:    page.GeneratedAt,
		AgencyID:       page.AgencyID,
		Boundary:       "Private authenticated realtime operations diagnostics only. Viewing this page sends no telemetry, runs no backend command, creates no retained evidence, changes no consumer status, contacts no external party, and records no compliance, consumer-acceptance, vendor, hardware, SLA, public-launch, production-readiness, production AVL reliability, or production-grade ETA claim.",
		StaleThreshold: page.StaleThreshold.String(),
		Summary:        summary,
		Feeds:          feeds,
		Usefulness:     buildRealtimeUsefulnessReview(page, summary, fleet, feeds),
		Publishing:     buildRealtimePublishingReview(page, summary, fleet, feeds),
		ReplayGuidance: buildRealtimeReplayGuide(page),
		Fleet:          fleet,
		Issues:         realtimeOperatorIssues(page, summary, fleet, feeds),
		Guidance:       realtimeQualityGuidance(page, summary, feeds),
		ClaimFlags:     operationsRealtimeClaimFlags{},
	}
}

func buildRealtimeSummary(page operationsPage, fleet []operationsRealtimeFleetRow) operationsRealtimeSummary {
	summary := operationsRealtimeSummary{
		Status:              checklistStatusUnknown,
		LatestTelemetryRows: len(page.Telemetry),
		StaleTelemetryRows:  page.StaleCount,
		DeviceBindings:      len(page.Devices),
		DoesNotProve:        "This private view does not show production AVL reliability, vendor compatibility, hardware certification, consumer display, compliance, hosted operation, SLA, public launch, production readiness, or production-grade ETA quality.",
	}
	if summary.LatestTelemetryRows >= summary.StaleTelemetryRows {
		summary.FreshTelemetryRows = summary.LatestTelemetryRows - summary.StaleTelemetryRows
	}
	for _, row := range fleet {
		switch row.Freshness {
		case "fresh":
			summary.DevicesReporting++
		case "not seen":
			summary.DevicesNotSeen++
		}
		if row.AssignmentSource == string(state.AssignmentSourceManualOverride) {
			summary.ManualOverrides++
		}
		if isRealtimeMatched(row) {
			summary.MatchedAssignments++
		} else {
			summary.UnknownAssignments++
		}
		if isRealtimeLowConfidence(row.Confidence) {
			summary.LowConfidenceRows++
		}
	}
	if page.TelemetryError != "" {
		summary.Status = checklistStatusBlocked
		summary.CurrentSignal = page.TelemetryError
		summary.NextAction = "Confirm the telemetry service and database are available, then return to this private center."
		return summary
	}
	switch {
	case summary.LatestTelemetryRows == 0 && summary.DeviceBindings == 0:
		summary.Status = checklistStatusMissing
		summary.CurrentSignal = "no latest telemetry rows and no device bindings are visible to this console"
		summary.NextAction = "Create or rotate a device credential, install it on a device or simulator, then send an authenticated sample telemetry event."
	case summary.LatestTelemetryRows == 0:
		summary.Status = checklistStatusMissing
		summary.CurrentSignal = fmt.Sprintf("%d device bindings exist, but no latest telemetry rows are visible", summary.DeviceBindings)
		summary.NextAction = "Check device installation or run the synthetic simulator from an operator shell."
	case summary.StaleTelemetryRows == summary.LatestTelemetryRows:
		summary.Status = checklistStatusNeedsReview
		summary.CurrentSignal = fmt.Sprintf("%d latest telemetry rows are stale at threshold %s", summary.StaleTelemetryRows, page.StaleThreshold)
		summary.NextAction = "Check device power, network, reporting cadence, and simulator or ingest logs before relying on Vehicle Positions."
	case summary.UnknownAssignments > 0 || summary.LowConfidenceRows > 0:
		summary.Status = checklistStatusNeedsReview
		summary.CurrentSignal = fmt.Sprintf("%d fresh rows, %d stale rows, %d unknown or unavailable assignments, %d low-confidence rows", summary.FreshTelemetryRows, summary.StaleTelemetryRows, summary.UnknownAssignments, summary.LowConfidenceRows)
		summary.NextAction = "Review assignment reasons and keep trip descriptors unknown when confidence is not sufficient."
	default:
		summary.Status = checklistStatusOK
		summary.CurrentSignal = fmt.Sprintf("%d latest telemetry rows are visible; %d are fresh; %d assignments are matched", summary.LatestTelemetryRows, summary.FreshTelemetryRows, summary.MatchedAssignments)
		summary.NextAction = "Continue monitoring freshness, assignment confidence, Vehicle Positions, Trip Updates diagnostics, and Alerts lifecycle."
	}
	return summary
}

func buildRealtimeUsefulnessReview(page operationsPage, summary operationsRealtimeSummary, fleet []operationsRealtimeFleetRow, feeds []operationsRealtimeFeedStatus) operationsRealtimeUsefulness {
	feedByID := realtimeFeedStatusByID(feeds)
	rows := []operationsRealtimeUsefulnessScore{
		vehiclePositionsUsefulnessScore(page, summary, feedByID["vehicle_positions"]),
		tripUpdatesUsefulnessScore(page, feedByID["trip_updates"]),
		alertsUsefulnessScore(feedByID["alerts"]),
	}
	return operationsRealtimeUsefulness{
		Status:        realtimeUsefulnessStatus(rows),
		Summary:       fmt.Sprintf("Private usefulness scoring reviews %d realtime feed types, %d fleet rows, and %d latest telemetry rows.", len(rows), len(fleet), summary.LatestTelemetryRows),
		Boundary:      "Usefulness scores are private operator diagnostics only. They are not SLA, uptime, production readiness, production-grade ETA, real-world accuracy, consumer display, public launch, compliance, vendor, or hardware proof.",
		Rows:          rows,
		Freshness:     realtimeFreshnessReviewRows(page, summary, feedByID),
		OmissionRules: realtimeConsumerSafeOmissionRules(),
	}
}

func realtimeFeedStatusByID(feeds []operationsRealtimeFeedStatus) map[string]operationsRealtimeFeedStatus {
	out := map[string]operationsRealtimeFeedStatus{}
	for _, feed := range feeds {
		out[feed.ID] = feed
	}
	return out
}

func buildRealtimePublishingReview(page operationsPage, summary operationsRealtimeSummary, fleet []operationsRealtimeFleetRow, feeds []operationsRealtimeFeedStatus) []operationsRealtimeFeedReview {
	feedByID := realtimeFeedStatusByID(feeds)
	return []operationsRealtimeFeedReview{
		vehiclePositionsPublishingReview(page, summary, fleet, feedByID["vehicle_positions"]),
		tripUpdatesPublishingReview(page, feedByID["trip_updates"]),
		alertsPublishingReview(page, feedByID["alerts"]),
	}
}

func vehiclePositionsPublishingReview(page operationsPage, summary operationsRealtimeSummary, fleet []operationsRealtimeFleetRow, feed operationsRealtimeFeedStatus) operationsRealtimeFeedReview {
	suppressAfter := suppressStaleVehicleAfter()
	seenVehicles := 0
	publishedVehicles := 0
	suppressedVehicles := 0
	unmatchedVehicles := 0
	lowConfidenceVehicles := 0
	for _, row := range fleet {
		if row.Freshness == "not seen" {
			continue
		}
		seenVehicles++
		if suppressAfter > 0 && time.Duration(row.AgeSeconds)*time.Second > suppressAfter {
			suppressedVehicles++
			continue
		}
		publishedVehicles++
		if !isRealtimeMatched(row) {
			unmatchedVehicles++
		}
		if isRealtimeLowConfidence(row.Confidence) {
			lowConfidenceVehicles++
		}
	}
	tripDescriptorsOmitted := max(0, publishedVehicles-summary.MatchedAssignments)
	status := checklistStatusOK
	switch {
	case seenVehicles == 0:
		status = checklistStatusMissing
	case suppressedVehicles > 0 || unmatchedVehicles > 0 || summary.StaleTelemetryRows > 0 || lowConfidenceVehicles > 0 || realtimeFeedNeedsReview(feed.State):
		status = checklistStatusNeedsReview
	}
	coverage := "not available"
	if publishedVehicles > 0 {
		coverage = fmt.Sprintf("%d of %d publishable vehicle rows have trip descriptors", summary.MatchedAssignments, publishedVehicles)
	}
	return operationsRealtimeFeedReview{
		ID:               "vehicle_positions",
		Label:            "Vehicle Positions publishing review",
		Status:           status,
		WhatLooksHealthy: fmt.Sprintf("%d latest vehicle row(s), %d estimated in Vehicle Positions, %s.", seenVehicles, publishedVehicles, coverage),
		NeedsAttention:   "Review any suppressed, stale, unmatched, low-confidence, or assignment-mismatch rows before relying on public vehicle movement.",
		NotProven:        "This private review does not show production AVL reliability, vendor compatibility, hardware certification, consumer display, compliance, SLA, uptime, or production readiness.",
		NextAction:       "Open Telemetry and Feed Health, then inspect why omitted vehicles or trip descriptors stayed out of the public feed.",
		Signals: []operationsRealtimePublishingSignal{
			{"Vehicle count", strconv.Itoa(seenVehicles), "Latest accepted telemetry vehicles visible to this private console."},
			{"Estimated Vehicle Positions rows", strconv.Itoa(publishedVehicles), "Vehicles not older than the stale suppression threshold."},
			{"Stale vehicles", strconv.Itoa(summary.StaleTelemetryRows), "Vehicles older than the matching stale threshold; trip descriptors should stay unknown or omitted."},
			{"Unmatched vehicles", strconv.Itoa(unmatchedVehicles), "Telemetry rows without a publishable trip assignment."},
			{"Suppressed vehicles", strconv.Itoa(suppressedVehicles), "Vehicles old enough to be omitted from Vehicle Positions entirely."},
			{"Trip descriptor coverage", coverage, "How many estimated Vehicle Positions rows include a trip descriptor."},
			{"Why not published", vehiclePositionsNotPublishedReason(suppressedVehicles, unmatchedVehicles, tripDescriptorsOmitted, lowConfidenceVehicles), "Reasons this page can infer without exposing raw payloads or private debug blobs."},
		},
	}
}

func vehiclePositionsNotPublishedReason(suppressedVehicles int, unmatchedVehicles int, tripDescriptorsOmitted int, lowConfidenceVehicles int) string {
	reasons := []string{}
	if suppressedVehicles > 0 {
		reasons = append(reasons, fmt.Sprintf("%d suppressed stale vehicle(s)", suppressedVehicles))
	}
	if unmatchedVehicles > 0 {
		reasons = append(reasons, fmt.Sprintf("%d unmatched vehicle(s)", unmatchedVehicles))
	}
	if tripDescriptorsOmitted > 0 {
		reasons = append(reasons, fmt.Sprintf("%d trip descriptor(s) omitted", tripDescriptorsOmitted))
	}
	if lowConfidenceVehicles > 0 {
		reasons = append(reasons, fmt.Sprintf("%d low-confidence assignment(s)", lowConfidenceVehicles))
	}
	if len(reasons) == 0 {
		return "No suppressed vehicles or omitted trip descriptors are visible in this bounded summary."
	}
	return strings.Join(reasons, "; ")
}

func tripUpdatesPublishingReview(page operationsPage, feed operationsRealtimeFeedStatus) operationsRealtimeFeedReview {
	quality := page.TripUpdatesQuality
	status := checklistStatusMissing
	emitted := 0
	eligible := 0
	withheld := 0
	source := "not recorded"
	latest := quality.Message
	if quality.Recorded {
		status = checklistStatusOK
		emitted = quality.TripUpdatesEmitted
		eligible = quality.EligiblePredictionCandidates
		withheld = countViewTotal(quality.WithheldByReason)
		source = firstNonEmpty(quality.AdapterName, "not available")
		latest = fmt.Sprintf("%s/%s at %s", quality.DiagnosticsStatus, quality.DiagnosticsReason, formatTimeForText(quality.SnapshotAt))
		if emitted == 0 || withheld > 0 || quality.UnknownAssignments > 0 || quality.AmbiguousAssignments > 0 || quality.StaleTelemetryRows > 0 || quality.CancellationAlertLinksMissing > 0 || realtimeFeedNeedsReview(feed.State) {
			status = checklistStatusNeedsReview
		}
	}
	signals := []operationsRealtimePublishingSignal{
		{"Prediction source", source, "Trip Updates remain behind the prediction adapter boundary."},
		{"Generated", strconv.Itoa(emitted), "Trip Updates emitted by the latest recorded diagnostics."},
		{"Eligible candidates", strconv.Itoa(eligible), "In-service prediction candidates considered safe enough to evaluate."},
		{"Withheld total", strconv.Itoa(withheld), "Rows withheld by explicit conservative reason."},
		{"Fallback reason", firstNonEmpty(quality.DiagnosticsReason, quality.Message, "not recorded"), "Why the latest diagnostic result generated, withheld, or fell back."},
		{"Stale telemetry", strconv.Itoa(quality.StaleTelemetryRows), "Stale telemetry rows should not produce invented ETA certainty."},
		{"Ambiguous assignments", strconv.Itoa(quality.AmbiguousAssignments), "Ambiguous assignment rows are withheld instead of guessed."},
		{"Low-confidence handling", tripUpdatesLowConfidenceHandling(quality.WithheldByReason), "Low-confidence predictor output remains withheld or falls back closed."},
	}
	for _, reason := range quality.WithheldByReason {
		signals = append(signals, operationsRealtimePublishingSignal{
			Label:   "Withheld reason: " + reason.Label,
			Value:   strconv.Itoa(reason.Count),
			Meaning: "Specific conservative Trip Updates withholding reason from the latest diagnostics.",
		})
	}
	return operationsRealtimeFeedReview{
		ID:               "trip_updates",
		Label:            "Trip Updates publishing review",
		Status:           status,
		WhatLooksHealthy: fmt.Sprintf("%d generated from %d eligible candidate(s); latest diagnostics: %s.", emitted, eligible, latest),
		NeedsAttention:   "Review generated versus withheld counts, fallback reason, stale telemetry, ambiguous assignments, low-confidence handling, and prediction source before relying on ETA-like output.",
		NotProven:        "This review does not show production-grade ETA quality, real-world ETA accuracy, consumer display, public launch, compliance, SLA, uptime, or production readiness.",
		NextAction:       "Open Prediction & ETA Lab for withheld reasons, adapter diagnostics, future-stop coverage, and backtest guidance.",
		Signals:          signals,
	}
}

func tripUpdatesLowConfidenceHandling(reasons []countView) string {
	for _, reason := range reasons {
		lower := strings.ToLower(reason.Label)
		if strings.Contains(lower, "low") && strings.Contains(lower, "confidence") {
			return fmt.Sprintf("%d low-confidence output(s) withheld", reason.Count)
		}
	}
	return "No explicit low-confidence withheld row is recorded; keep reviewing confidence thresholds and unknown assignments."
}

func alertsPublishingReview(page operationsPage, feed operationsRealtimeFeedStatus) operationsRealtimeFeedReview {
	status := checklistStatusNeedsReview
	if !realtimeFeedNeedsReview(feed.State) && page.TripUpdatesQuality.CancellationAlertLinksMissing == 0 {
		status = checklistStatusOK
	}
	if strings.TrimSpace(feed.State) == "" || feed.State == "not available yet" {
		status = checklistStatusMissing
	}
	return operationsRealtimeFeedReview{
		ID:               "alerts",
		Label:            "Alerts lifecycle review",
		Status:           status,
		WhatLooksHealthy: "Alerts feed metadata and validator/feed-health rows are visible, and cancellation alert linkage has no recorded missing rows.",
		NeedsAttention:   "Active alert counts and stale alert details remain in the Alerts Console; review missing cancellation/disruption links before relying on cancellation messaging.",
		NotProven:        "This review does not show agency approval, consumer display, public launch, compliance, SLA, uptime, or production readiness.",
		NextAction:       "Open the Alerts Console for active, stale, planned, archived, cancellation, and disruption review, then check Alerts feed health and validation.",
		Signals: []operationsRealtimePublishingSignal{
			{"Active alerts", "review in Alerts Console", "The Realtime Center links to the lifecycle surface instead of duplicating alert authoring state."},
			{"Stale alerts", "review in Alerts Console", "Archive or update stale alert records in the dedicated Alerts workflow."},
			{"Missing cancellation links", strconv.Itoa(page.TripUpdatesQuality.CancellationAlertLinksMissing), "Canceled Trip Updates that still need an Alerts lifecycle review."},
			{"Missing disruption links", "not automatically inferred", "Service disruptions require operator-authored alert review; telemetry alone does not publish alerts."},
			{"Service disruption review", firstNonEmpty(feed.LatestSignal, "no Alerts feed signal"), "Use feed health plus Alerts Console before sharing alert status."},
		},
	}
}

func countViewTotal(rows []countView) int {
	total := 0
	for _, row := range rows {
		total += row.Count
	}
	return total
}

func buildRealtimeReplayGuide(page operationsPage) operationsRealtimeReplayGuide {
	status := checklistStatusNeedsReview
	if len(page.TelemetrySimulator.Scenarios) > 0 {
		status = checklistStatusOK
	}
	return operationsRealtimeReplayGuide{
		Status:       status,
		Summary:      fmt.Sprintf("%d committed synthetic simulator scenario(s) are available for local browser review.", len(page.TelemetrySimulator.Scenarios)),
		BrowserStart: "/admin/operations/telemetry-simulator",
		LocalReplay:  "Use the browser dry-run preview first; an administrator may run fixed local simulator or realtime-quality commands when credentials and private DB access are needed.",
		ReviewAfter:  "After replay, return to Realtime Center, Feed Health, Prediction & ETA Lab, Validation Health, and Alerts Console.",
		Boundary:     "Replay guidance is local/synthetic review only. The browser does not execute shell commands, collect tokens, send telemetry, create evidence, contact external systems, prove vendor compatibility, prove production AVL reliability, or prove ETA quality.",
		Steps: []operationsRealtimeReplayStep{
			{"browser_preview", "Preview a scenario in the browser", "Open Telemetry Simulator and choose on-route, stale, low-quality GPS, after-midnight, or block-transition fixture metadata.", "Preview only; no telemetry send, no token collection, and no command execution."},
			{"local_send", "Run local synthetic send only when appropriate", "An administrator can use the fixed simulator command with deployment-owned tokens when intentional local ingest is required.", "Keep tokens outside browser notes and do not commit generated `.cache` diagnostics."},
			{"realtime_review", "Review realtime feed usefulness", "Check Vehicle Positions publishing review, Trip Updates withheld reasons, Alerts lifecycle review, and GTFS-Realtime validator health.", "Healthy local signals do not prove external approval, consumer display, production readiness, or real-world accuracy."},
		},
	}
}

func vehiclePositionsUsefulnessScore(page operationsPage, summary operationsRealtimeSummary, feed operationsRealtimeFeedStatus) operationsRealtimeUsefulnessScore {
	score := 0
	label := "missing_signal"
	helpful := "No accepted latest telemetry rows are visible."
	needs := "Device credentials, telemetry ingest, and Vehicle Positions feed-health rows need review."
	if summary.LatestTelemetryRows > 0 {
		score = 1
		label = "telemetry_visible"
		helpful = fmt.Sprintf("%d latest telemetry rows are visible.", summary.LatestTelemetryRows)
		needs = fmt.Sprintf("%d stale rows; %d unknown assignments; %d low-confidence rows.", summary.StaleTelemetryRows, summary.UnknownAssignments, summary.LowConfidenceRows)
		if summary.FreshTelemetryRows > 0 {
			score = 2
			label = "useful_with_review"
		}
		if summary.FreshTelemetryRows > 0 && summary.MatchedAssignments > 0 && summary.UnknownAssignments == 0 && summary.LowConfidenceRows == 0 && !realtimeFeedNeedsReview(feed.State) {
			score = 3
			label = "useful_for_local_review"
			needs = "Continue periodic stale telemetry, assignment, validator, and feed-health review."
		}
		if summary.StaleTelemetryRows == summary.LatestTelemetryRows {
			label = "stale_only"
		}
	}
	return operationsRealtimeUsefulnessScore{
		ID:                   "vehicle_positions",
		Label:                "Vehicle Positions usefulness",
		Score:                score,
		ScoreLabel:           label,
		CurrentSignal:        firstNonEmpty(feed.LatestSignal, summary.CurrentSignal, "Vehicle Positions signal is not available."),
		HelpfulSignal:        helpful,
		NeedsReviewSignal:    needs,
		ConsumerSafeBehavior: "Emit only defensible vehicle position fields; suppress stale rows when configured and omit trip descriptors when assignment confidence is weak.",
		NextAction:           "Review telemetry freshness, stale suppression, assignment confidence, Vehicle Positions feed health, and realtime validation together.",
		DoesNotProve:         "This private view does not show field reliability, vendor compatibility, hardware certification, consumer display, compliance, SLA, uptime, or production readiness.",
		Details: []countView{
			{Label: "fresh_telemetry", Count: summary.FreshTelemetryRows},
			{Label: "stale_telemetry", Count: summary.StaleTelemetryRows},
			{Label: "matched_assignments", Count: summary.MatchedAssignments},
			{Label: "unknown_assignments", Count: summary.UnknownAssignments},
		},
	}
}

func tripUpdatesUsefulnessScore(page operationsPage, feed operationsRealtimeFeedStatus) operationsRealtimeUsefulnessScore {
	score := 0
	label := "diagnostics_missing"
	helpful := "Trip Updates diagnostics are not recorded."
	needs := "Review telemetry and assignment confidence before relying on prediction output."
	details := []countView{}
	if page.TripUpdatesQuality.Recorded {
		score = 1
		label = "withheld_or_empty"
		helpful = fmt.Sprintf("%d Trip Updates emitted from %d eligible candidates.", page.TripUpdatesQuality.TripUpdatesEmitted, page.TripUpdatesQuality.EligiblePredictionCandidates)
		needs = fmt.Sprintf("%d unknown assignments; %d ambiguous assignments; %d stale telemetry rows.", page.TripUpdatesQuality.UnknownAssignments, page.TripUpdatesQuality.AmbiguousAssignments, page.TripUpdatesQuality.StaleTelemetryRows)
		details = append(details,
			countView{Label: "emitted", Count: page.TripUpdatesQuality.TripUpdatesEmitted},
			countView{Label: "eligible_candidates", Count: page.TripUpdatesQuality.EligiblePredictionCandidates},
			countView{Label: "unknown_assignments", Count: page.TripUpdatesQuality.UnknownAssignments},
			countView{Label: "ambiguous_assignments", Count: page.TripUpdatesQuality.AmbiguousAssignments},
			countView{Label: "stale_telemetry", Count: page.TripUpdatesQuality.StaleTelemetryRows},
		)
		details = append(details, page.TripUpdatesQuality.WithheldByReason...)
		if page.TripUpdatesQuality.TripUpdatesEmitted > 0 {
			score = 2
			label = "generated_with_review"
		}
		if page.TripUpdatesQuality.TripUpdatesEmitted > 0 && page.TripUpdatesQuality.UnknownAssignments == 0 && page.TripUpdatesQuality.AmbiguousAssignments == 0 && page.TripUpdatesQuality.StaleTelemetryRows == 0 && !realtimeFeedNeedsReview(feed.State) {
			score = 3
			label = "useful_for_local_review"
			needs = "Continue reviewing withheld reasons, future-stop coverage, adapter status, and validation before stronger language."
		}
	}
	return operationsRealtimeUsefulnessScore{
		ID:                   "trip_updates",
		Label:                "Trip Updates usefulness",
		Score:                score,
		ScoreLabel:           label,
		CurrentSignal:        firstNonEmpty(feed.LatestSignal, page.TripUpdatesQuality.Message, "Trip Updates signal is not available."),
		HelpfulSignal:        helpful,
		NeedsReviewSignal:    needs,
		ConsumerSafeBehavior: "Withhold Trip Updates or emit valid empty/fallback output when prediction evidence is stale, ambiguous, low confidence, or missing future-stop support.",
		NextAction:           "Review emitted counts, withheld reasons, stale telemetry, adapter fallback state, and realtime validation before relying on ETAs.",
		DoesNotProve:         "This private view does not show production-grade ETA quality, real-world ETA accuracy, consumer display, public launch, compliance, SLA, or production readiness.",
		Details:              details,
	}
}

func alertsUsefulnessScore(feed operationsRealtimeFeedStatus) operationsRealtimeUsefulnessScore {
	score := 0
	label := "missing_signal"
	helpful := "Alerts feed state is not available."
	needs := "Open Alerts Console and feed health before relying on alert output."
	if strings.TrimSpace(feed.State) != "" && feed.State != "not available yet" {
		score = 1
		label = "lifecycle_needs_review"
		helpful = "Alerts feed state is visible in private feed health."
		needs = "Active alert count is not exposed in this bounded summary; lifecycle review still belongs in the Alerts Console."
		if !realtimeFeedNeedsReview(feed.State) {
			score = 2
			label = "configured_lifecycle_review"
		}
	}
	return operationsRealtimeUsefulnessScore{
		ID:                   "alerts",
		Label:                "Alerts usefulness",
		Score:                score,
		ScoreLabel:           label,
		CurrentSignal:        firstNonEmpty(feed.LatestSignal, "Alerts signal is not available."),
		HelpfulSignal:        helpful,
		NeedsReviewSignal:    needs,
		ConsumerSafeBehavior: "Publish only authored and reviewed Alerts records; do not infer a disruption or an all-clear state from missing private rows.",
		NextAction:           "Review active, planned, archived, and canceled-service alert workflows in the Alerts Console, then check feed health and validation.",
		DoesNotProve:         "This private view does not show agency approval, consumer display, public launch, compliance, SLA, uptime, or production readiness.",
	}
}

func realtimeUsefulnessStatus(rows []operationsRealtimeUsefulnessScore) string {
	if len(rows) == 0 {
		return checklistStatusUnknown
	}
	hasZero := false
	hasReview := false
	for _, row := range rows {
		if row.Score <= 0 {
			hasZero = true
			continue
		}
		if row.Score < 3 {
			hasReview = true
		}
	}
	switch {
	case hasZero:
		return checklistStatusMissing
	case hasReview:
		return checklistStatusNeedsReview
	default:
		return checklistStatusOK
	}
}

func realtimeFreshnessReviewRows(page operationsPage, summary operationsRealtimeSummary, feeds map[string]operationsRealtimeFeedStatus) []operationsRealtimeFreshnessReview {
	tripSignal := page.TripUpdatesQuality.Message
	tripStatus := checklistStatusMissing
	if page.TripUpdatesQuality.Recorded {
		tripStatus = checklistStatusNeedsReview
		if page.TripUpdatesQuality.TripUpdatesEmitted > 0 && page.TripUpdatesQuality.StaleTelemetryRows == 0 {
			tripStatus = checklistStatusOK
		}
		tripSignal = fmt.Sprintf("%s/%s at %s", page.TripUpdatesQuality.DiagnosticsStatus, page.TripUpdatesQuality.DiagnosticsReason, formatTimeForText(page.TripUpdatesQuality.SnapshotAt))
	}
	return []operationsRealtimeFreshnessReview{
		{
			ID:            "telemetry_freshness",
			Label:         "Telemetry freshness",
			Status:        summary.Status,
			CurrentSignal: fmt.Sprintf("%d latest rows; %d fresh; %d stale at %s", summary.LatestTelemetryRows, summary.FreshTelemetryRows, summary.StaleTelemetryRows, page.StaleThreshold),
			NextAction:    "Check stale devices before changing matching thresholds or prediction settings.",
			DoesNotProve:  "Fresh telemetry does not show field reliability, SLA, uptime, or vendor compatibility.",
		},
		{
			ID:            "device_state",
			Label:         "Device state triage",
			Status:        realtimeDeviceFreshnessStatus(summary),
			CurrentSignal: fmt.Sprintf("%d bindings; %d reporting; %d not seen", summary.DeviceBindings, summary.DevicesReporting, summary.DevicesNotSeen),
			NextAction:    "Review device bindings and one-time token lifecycle when devices are not seen.",
			DoesNotProve:  "A visible binding does not show hardware certification or production AVL reliability.",
		},
		realtimeFeedFreshnessRow("vehicle_positions_feed", "Vehicle Positions feed freshness", feeds["vehicle_positions"]),
		{
			ID:            "trip_updates_diagnostics",
			Label:         "Trip Updates diagnostics freshness",
			Status:        tripStatus,
			CurrentSignal: tripSignal,
			NextAction:    "Review diagnostics recency, withheld reasons, and adapter fallback before relying on prediction output.",
			DoesNotProve:  "Recent diagnostics do not prove ETA accuracy or consumer display.",
		},
		realtimeFeedFreshnessRow("alerts_feed", "Alerts feed freshness", feeds["alerts"]),
	}
}

func realtimeFeedFreshnessRow(id string, label string, feed operationsRealtimeFeedStatus) operationsRealtimeFreshnessReview {
	status := checklistStatusMissing
	if strings.TrimSpace(feed.State) != "" {
		status = checklistStatusNeedsReview
		if !realtimeFeedNeedsReview(feed.State) {
			status = checklistStatusOK
		}
	}
	return operationsRealtimeFreshnessReview{
		ID:            id,
		Label:         label,
		Status:        status,
		CurrentSignal: fmt.Sprintf("%s; %s; %s", firstNonEmpty(feed.State, "not available"), firstNonEmpty(feed.LatestSignal, "no latest signal"), firstNonEmpty(feed.StaleOrWithheld, "no stale or withheld summary")),
		NextAction:    firstNonEmpty(feed.NextAction, "Open feed health and validation center before relying on this feed."),
		DoesNotProve:  feed.DoesNotProve,
	}
}

func realtimeDeviceFreshnessStatus(summary operationsRealtimeSummary) string {
	switch {
	case summary.DeviceBindings == 0:
		return checklistStatusMissing
	case summary.DevicesNotSeen > 0:
		return checklistStatusNeedsReview
	default:
		return checklistStatusOK
	}
}

func realtimeConsumerSafeOmissionRules() []operationsRealtimeOmissionRule {
	return []operationsRealtimeOmissionRule{
		{
			ID:           "stale_vehicle_position",
			Condition:    "Telemetry is stale or older than the configured suppression threshold.",
			SafeBehavior: "Prefer suppressing stale vehicles or omitting stale-sensitive fields over presenting old movement as current.",
			ReviewStep:   "Check device power, network, timestamps, and ingest cadence before changing thresholds.",
			DoesNotProve: "Suppressing stale rows does not show fleet reliability or uptime.",
		},
		{
			ID:           "unknown_assignment",
			Condition:    "Assignment is unknown, ambiguous, manually withheld, or below confidence threshold.",
			SafeBehavior: "Emit Vehicle Positions without a trip descriptor or keep the public trip context unknown.",
			ReviewStep:   "Review service day, after-midnight trips, frequency service, block continuity, and active overrides.",
			DoesNotProve: "A later match does not show consumer display or ETA quality.",
		},
		{
			ID:           "trip_updates_withheld",
			Condition:    "Trip Updates lack safe prediction evidence, future-stop support, or adapter output.",
			SafeBehavior: "Withhold Trip Updates or emit valid empty/fallback output instead of inventing ETAs.",
			ReviewStep:   "Review withheld reasons, adapter diagnostics, stale telemetry, and validation before changing prediction behavior.",
			DoesNotProve: "Withheld or emitted Trip Updates do not prove production-grade ETA quality or real-world accuracy.",
		},
		{
			ID:           "alerts_not_authored",
			Condition:    "No active alert is authored for a disruption or cancellation signal.",
			SafeBehavior: "Do not infer or publish a service alert automatically from telemetry or prediction data.",
			ReviewStep:   "Use the Alerts Console lifecycle review and cancellation-link guidance before publishing alert records.",
			DoesNotProve: "An absent alert does not show there is no disruption or that an agency approved messaging.",
		},
	}
}

func buildRealtimeFleetRows(page operationsPage) []operationsRealtimeFleetRow {
	rows := make([]operationsRealtimeFleetRow, 0, len(page.Telemetry)+len(page.DeviceRows))
	seenDeviceVehicle := map[string]bool{}
	for _, telemetryRow := range page.Telemetry {
		row := realtimeFleetRowFromTelemetry(telemetryRow)
		rows = append(rows, row)
		seenDeviceVehicle[deviceVehicleKey(row.DeviceID, row.VehicleID)] = true
	}
	for _, deviceRow := range page.DeviceRows {
		key := deviceVehicleKey(deviceRow.DeviceID, deviceRow.VehicleID)
		if seenDeviceVehicle[key] {
			continue
		}
		rows = append(rows, realtimeFleetRowFromDevice(deviceRow))
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].VehicleID == rows[j].VehicleID {
			return rows[i].DeviceID < rows[j].DeviceID
		}
		return rows[i].VehicleID < rows[j].VehicleID
	})
	return rows
}

func realtimeFleetRowFromTelemetry(row telemetryView) operationsRealtimeFleetRow {
	freshness := "fresh"
	if row.Stale {
		freshness = "stale"
	}
	out := operationsRealtimeFleetRow{
		VehicleID:        row.VehicleID,
		DeviceID:         row.DeviceID,
		Freshness:        freshness,
		ObservedAt:       formatTimeForText(&row.ObservedAt),
		ReceivedAt:       formatTimeForText(&row.ReceivedAt),
		AgeSeconds:       row.AgeSeconds,
		AssignmentState:  firstNonEmpty(row.AssignmentState, string(state.StateUnknown)),
		DegradedState:    row.DegradedState,
		AssignmentSource: row.AssignmentSource,
		RouteID:          row.RouteID,
		TripID:           row.TripID,
		Confidence:       row.Confidence,
		ReasonCodes:      append([]string(nil), row.ReasonCodes...),
		DoesNotProve:     "This private view does not show real fleet reliability, vendor compatibility, hardware certification, consumer display, compliance, or production-grade ETA quality.",
	}
	out.CurrentSignal = realtimeFleetSignal(out)
	out.NextAction = realtimeFleetNextAction(out)
	return out
}

func realtimeFleetRowFromDevice(row operationsDeviceRow) operationsRealtimeFleetRow {
	out := operationsRealtimeFleetRow{
		VehicleID:        row.VehicleID,
		DeviceID:         row.DeviceID,
		Freshness:        row.Freshness,
		ObservedAt:       formatTimeForText(row.LatestObservedAt),
		ReceivedAt:       formatTimeForText(row.LatestReceivedAt),
		AssignmentState:  "not available",
		AssignmentSource: row.AssignmentSource,
		DoesNotProve:     "A device binding does not show hardware certification, vendor compatibility, production AVL reliability, consumer display, compliance, or production readiness.",
	}
	if row.LatestAgeSeconds != nil {
		out.AgeSeconds = *row.LatestAgeSeconds
	}
	out.CurrentSignal = realtimeFleetSignal(out)
	out.NextAction = realtimeFleetNextAction(out)
	return out
}

func realtimeFeedStatuses(page operationsPage) []operationsRealtimeFeedStatus {
	source := []operationsRealtimeUsefulnessRow{
		page.FeedHealth.RealtimeUsefulness.VehiclePositions,
		page.FeedHealth.RealtimeUsefulness.TripUpdates,
		page.FeedHealth.RealtimeUsefulness.Alerts,
	}
	rows := make([]operationsRealtimeFeedStatus, 0, len(source))
	for _, row := range source {
		rows = append(rows, operationsRealtimeFeedStatus{
			ID:              row.ID,
			Label:           row.Label,
			State:           row.State,
			Count:           row.Count,
			LatestSignal:    row.LatestSignal,
			StaleOrWithheld: row.StaleOrHeld,
			Adapter:         row.Adapter,
			NextAction:      row.NextAction,
			AdminLink:       row.AdminLink,
			DoesNotProve:    row.DoesNotProve,
			Details:         row.Details,
		})
	}
	return rows
}

func realtimeOperatorIssues(page operationsPage, summary operationsRealtimeSummary, fleet []operationsRealtimeFleetRow, feeds []operationsRealtimeFeedStatus) []operationsRealtimeIssue {
	issues := []operationsRealtimeIssue{}
	if page.TelemetryError != "" {
		issues = append(issues, operationsRealtimeIssue{
			Severity:     checklistStatusBlocked,
			Area:         "Telemetry repository",
			Signal:       page.TelemetryError,
			NextAction:   "Confirm the private telemetry service and database are available.",
			AdminLink:    "/admin/operations/telemetry",
			DoesNotProve: "A restored repository connection does not show production AVL reliability or uptime.",
		})
	}
	if summary.LatestTelemetryRows == 0 {
		issues = append(issues, operationsRealtimeIssue{
			Severity:     checklistStatusMissing,
			Area:         "Fleet freshness",
			Signal:       summary.CurrentSignal,
			NextAction:   summary.NextAction,
			AdminLink:    "/admin/operations/devices",
			DoesNotProve: "A sample row does not show vendor compatibility, hardware certification, or production readiness.",
		})
	}
	for _, row := range fleet {
		if row.Freshness == "stale" {
			issues = append(issues, operationsRealtimeIssue{
				Severity:     checklistStatusNeedsReview,
				Area:         "Stale telemetry",
				Signal:       fmt.Sprintf("%s / %s was last observed %s; age=%d seconds", firstNonEmpty(row.VehicleID, "unknown vehicle"), firstNonEmpty(row.DeviceID, "unknown device"), row.ObservedAt, row.AgeSeconds),
				NextAction:   row.NextAction,
				AdminLink:    "/admin/operations/telemetry",
				DoesNotProve: row.DoesNotProve,
			})
		}
		if row.Freshness == "not seen" {
			issues = append(issues, operationsRealtimeIssue{
				Severity:     checklistStatusMissing,
				Area:         "Device not seen",
				Signal:       fmt.Sprintf("%s / %s has a binding but no latest accepted telemetry row", firstNonEmpty(row.VehicleID, "unknown vehicle"), firstNonEmpty(row.DeviceID, "unknown device")),
				NextAction:   row.NextAction,
				AdminLink:    "/admin/operations/devices",
				DoesNotProve: row.DoesNotProve,
			})
		}
		if !isRealtimeMatched(row) && row.Freshness != "not seen" {
			issues = append(issues, operationsRealtimeIssue{
				Severity:     checklistStatusNeedsReview,
				Area:         "Assignment confidence",
				Signal:       fmt.Sprintf("%s assignment=%s degraded=%s confidence=%s reasons=%s", firstNonEmpty(row.VehicleID, "unknown vehicle"), firstNonEmpty(row.AssignmentState, "not available"), firstNonEmpty(row.DegradedState, "none"), firstNonEmpty(row.Confidence, "not available"), strings.Join(row.ReasonCodes, ", ")),
				NextAction:   realtimeFleetNextAction(row),
				AdminLink:    "/admin/operations/telemetry",
				DoesNotProve: "Low-confidence or unknown assignments should not become public trip certainty.",
			})
		}
	}
	for _, feed := range feeds {
		if !realtimeFeedNeedsReview(feed.State) {
			continue
		}
		issues = append(issues, operationsRealtimeIssue{
			Severity:     checklistStatusNeedsReview,
			Area:         feed.Label,
			Signal:       fmt.Sprintf("%s; %s; %s", feed.State, feed.Count, feed.StaleOrWithheld),
			NextAction:   feed.NextAction,
			AdminLink:    feed.AdminLink,
			DoesNotProve: feed.DoesNotProve,
		})
	}
	if len(issues) == 0 {
		return []operationsRealtimeIssue{{
			Severity:     checklistStatusOK,
			Area:         "Current realtime review",
			Signal:       "No stale, not-seen, low-confidence, or withheld realtime rows are visible in this bounded summary.",
			NextAction:   "Continue periodic private monitoring of telemetry, assignments, Vehicle Positions, Trip Updates diagnostics, and Alerts lifecycle.",
			AdminLink:    "/admin/operations/feed-health",
			DoesNotProve: "A quiet private dashboard does not show compliance, consumer display, production uptime, or ETA quality.",
		}}
	}
	const limit = 12
	if len(issues) > limit {
		issues = append(issues[:limit], operationsRealtimeIssue{
			Severity:     checklistStatusNeedsReview,
			Area:         "Bounded issue list",
			Signal:       fmt.Sprintf("%d additional realtime issue rows are hidden from this browser summary", len(issues)-limit),
			NextAction:   "Use the private JSON export and underlying diagnostics pages for deeper local review.",
			AdminLink:    "/admin/operations/realtime.json",
			DoesNotProve: "A bounded issue count is not an exhaustive operational proof.",
		})
	}
	return issues
}

func realtimeQualityGuidance(page operationsPage, summary operationsRealtimeSummary, feeds []operationsRealtimeFeedStatus) []operationsRealtimeGuidance {
	feedByID := map[string]operationsRealtimeFeedStatus{}
	for _, feed := range feeds {
		feedByID[feed.ID] = feed
	}
	tripUpdates := feedByID["trip_updates"]
	vehiclePositions := feedByID["vehicle_positions"]
	alerts := feedByID["alerts"]
	return []operationsRealtimeGuidance{
		{
			ID:           "stale_telemetry",
			Label:        "Stale telemetry",
			WhatItMeans:  "A vehicle can have a last known position that is too old for confident realtime service.",
			ReviewSignal: fmt.Sprintf("%d stale latest telemetry rows at threshold %s", summary.StaleTelemetryRows, page.StaleThreshold),
			NextAction:   "Check device power, network, clock, and reporting cadence before changing service assignments.",
			DoesNotProve: "Fresh rows do not prove hardware certification, vendor compatibility, SLA, or production AVL reliability.",
		},
		{
			ID:           "low_confidence_assignments",
			Label:        "Unknown, ambiguous, or low-confidence assignment",
			WhatItMeans:  "The safer realtime output is to omit a trip descriptor when matching evidence is weak.",
			ReviewSignal: fmt.Sprintf("%d unknown or unavailable assignments; %d low-confidence rows", summary.UnknownAssignments, summary.LowConfidenceRows),
			NextAction:   "Review route/trip hints, service day, after-midnight service, block continuity, and any active operator override.",
			DoesNotProve: "A matched assignment does not show public consumer display or production-grade ETA quality.",
		},
		{
			ID:           "out_of_order_low_quality_gps",
			Label:        "Out-of-order or low-quality GPS",
			WhatItMeans:  "The browser summary only shows latest accepted rows; rejected or out-of-order samples belong in ingest logs and connector conformance outputs.",
			ReviewSignal: "If positions jump, stop updating, or timestamps move backward, treat route/trip certainty as suspect.",
			NextAction:   "Review connector conformance fixtures, telemetry simulator cases, and private ingest diagnostics without exposing raw payloads in HTML.",
			DoesNotProve: "Synthetic conformance does not show real vendor compatibility, hardware certification, or real-world fleet accuracy.",
		},
		{
			ID:           "vehicle_positions_debug",
			Label:        "Vehicle Positions debug",
			WhatItMeans:  "Vehicle Positions are useful only when authenticated telemetry and conservative assignment state stay fresh enough.",
			ReviewSignal: fmt.Sprintf("%s; %s", firstNonEmpty(vehiclePositions.State, "unknown"), firstNonEmpty(vehiclePositions.StaleOrWithheld, "no stale summary")),
			NextAction:   "Open telemetry freshness and feed health before treating Vehicle Positions output as ready for local review.",
			DoesNotProve: vehiclePositions.DoesNotProve,
		},
		{
			ID:           "trip_updates_withheld",
			Label:        "Trip Updates withheld or fallback",
			WhatItMeans:  "Trip Updates may be intentionally empty or withheld when the predictor lacks safe evidence.",
			ReviewSignal: fmt.Sprintf("%s; adapter=%s; %s", firstNonEmpty(tripUpdates.State, "unknown"), firstNonEmpty(tripUpdates.Adapter, "not configured"), firstNonEmpty(tripUpdates.StaleOrWithheld, "no withheld summary")),
			NextAction:   "Review withheld reasons, stale telemetry, prediction adapter status, and fallback state before relying on ETAs.",
			DoesNotProve: "Trip Updates diagnostics do not prove production-grade ETA quality or real-world accuracy.",
		},
		{
			ID:           "alerts_lifecycle",
			Label:        "Alerts lifecycle",
			WhatItMeans:  "Alerts need active/planned/archive review and feed validation separate from telemetry health.",
			ReviewSignal: fmt.Sprintf("%s; %s", firstNonEmpty(alerts.State, "unknown"), firstNonEmpty(alerts.LatestSignal, "no alert signal")),
			NextAction:   "Open the Alerts Console for safe edit links, then verify the Alerts feed health row.",
			DoesNotProve: "An alert visible in the private console does not show consumer display, public launch, agency approval, or compliance.",
		},
	}
}

func realtimeFeedNeedsReview(state string) bool {
	normalized := strings.ToLower(strings.TrimSpace(state))
	if normalized == "" {
		return true
	}
	switch normalized {
	case checklistStatusOK, "available", "ready", "recorded", "potentially non-empty":
		return false
	default:
		return true
	}
}

func realtimeFleetSignal(row operationsRealtimeFleetRow) string {
	if row.Freshness == "not seen" {
		return "device binding is visible, but no latest accepted telemetry row is visible"
	}
	parts := []string{row.Freshness + " telemetry"}
	if row.ObservedAt != "" && row.ObservedAt != "not available" {
		parts = append(parts, "observed "+row.ObservedAt)
	}
	if row.AssignmentState != "" && row.AssignmentState != "not available" {
		parts = append(parts, "assignment "+row.AssignmentState)
	}
	if row.DegradedState != "" && row.DegradedState != string(state.DegradedNone) {
		parts = append(parts, "degraded "+row.DegradedState)
	}
	if row.Confidence != "" {
		parts = append(parts, "confidence "+row.Confidence)
	}
	return strings.Join(parts, "; ")
}

func realtimeFleetNextAction(row operationsRealtimeFleetRow) string {
	if row.Freshness == "not seen" {
		return "Install or check the device credential, then send an authenticated sample telemetry event from the device or simulator."
	}
	if row.Freshness == "stale" {
		return "Check device power, network, and reporting cadence before changing assignments."
	}
	if row.AssignmentSource == string(state.AssignmentSourceManualOverride) {
		return "Review the active operator override; automatic matching should not replace it until it expires or is cleared."
	}
	if !isRealtimeMatched(row) {
		return "Keep the trip descriptor unknown until matching confidence improves or an operator override is applied."
	}
	if isRealtimeLowConfidence(row.Confidence) {
		return "Review low-confidence matching evidence before relying on route or trip context."
	}
	return "Continue monitoring freshness and assignment confidence."
}

func isRealtimeMatched(row operationsRealtimeFleetRow) bool {
	if row.AssignmentState == "" || row.AssignmentState == "not available" || row.AssignmentState == string(state.StateUnknown) {
		return false
	}
	if row.DegradedState != "" && row.DegradedState != string(state.DegradedNone) {
		return false
	}
	if isRealtimeLowConfidence(row.Confidence) {
		return false
	}
	return row.TripID != "" || row.AssignmentSource == string(state.AssignmentSourceManualOverride)
}

func isRealtimeLowConfidence(value string) bool {
	if strings.TrimSpace(value) == "" {
		return false
	}
	confidence, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}
	return confidence < 0.70
}
