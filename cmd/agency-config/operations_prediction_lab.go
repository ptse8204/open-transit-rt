package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/prediction"
)

type predictionLabView struct {
	GeneratedAt     time.Time                  `json:"generated_at"`
	AgencyID        string                     `json:"agency_id"`
	Boundary        string                     `json:"boundary"`
	Summary         predictionLabSummary       `json:"summary"`
	Deterministic   predictionLabDeterministic `json:"deterministic"`
	WithheldReasons []predictionLabReason      `json:"withheld_reasons"`
	ShadowReview    predictionLabShadowReview  `json:"shadow_review"`
	ReviewRows      []predictionLabReviewRow   `json:"review_rows"`
	Commands        []predictionLabCommand     `json:"commands"`
	ClaimFlags      predictionLabClaimFlags    `json:"claim_flags"`
}

type predictionLabSummary struct {
	Status                       string `json:"status"`
	CurrentSignal                string `json:"current_signal"`
	AdapterName                  string `json:"adapter_name"`
	DiagnosticsStatus            string `json:"diagnostics_status"`
	DiagnosticsReason            string `json:"diagnostics_reason"`
	ActiveFeedVersionID          string `json:"active_feed_version_id"`
	EligiblePredictionCandidates int    `json:"eligible_prediction_candidates"`
	TripUpdatesEmitted           int    `json:"trip_updates_emitted"`
	WithheldCount                int    `json:"withheld_count"`
	NextAction                   string `json:"next_action"`
	DoesNotProve                 string `json:"does_not_prove"`
}

type predictionLabDeterministic struct {
	Status       string                       `json:"status"`
	Boundary     string                       `json:"boundary"`
	ReviewSignal string                       `json:"review_signal"`
	NextAction   string                       `json:"next_action"`
	DoesNotProve string                       `json:"does_not_prove"`
	Rows         []predictionLabDiagnosticRow `json:"rows"`
}

type predictionLabDiagnosticRow struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	DoesNotProve  string `json:"does_not_prove"`
}

type predictionLabReason struct {
	Reason       string `json:"reason"`
	Label        string `json:"label"`
	Count        int    `json:"count"`
	WhatItMeans  string `json:"what_it_means"`
	NextAction   string `json:"next_action"`
	DoesNotProve string `json:"does_not_prove"`
}

type predictionLabShadowReview struct {
	Status       string                   `json:"status"`
	Boundary     string                   `json:"boundary"`
	NextAction   string                   `json:"next_action"`
	DoesNotProve string                   `json:"does_not_prove"`
	Rows         []predictionLabShadowRow `json:"rows"`
}

type predictionLabShadowRow struct {
	ID              string `json:"id"`
	Label           string `json:"label"`
	Status          string `json:"status"`
	Reason          string `json:"reason"`
	Latency         string `json:"latency"`
	CountComparison string `json:"count_comparison"`
	FailureBehavior string `json:"failure_behavior"`
	FirstSafeCheck  string `json:"first_safe_check"`
	DoesNotProve    string `json:"does_not_prove"`
}

type predictionLabReviewRow struct {
	Severity     string `json:"severity"`
	Area         string `json:"area"`
	Signal       string `json:"signal"`
	NextAction   string `json:"next_action"`
	DoesNotProve string `json:"does_not_prove"`
}

type predictionLabCommand struct {
	ID             string `json:"id"`
	Label          string `json:"label"`
	CommandLine    string `json:"command_line"`
	ExpectedResult string `json:"expected_result"`
	DoesNotProve   string `json:"does_not_prove"`
}

type predictionLabClaimFlags struct {
	BrowserPredictorRunEnabled     bool `json:"browser_predictor_run_enabled"`
	ExternalNetworkContacted       bool `json:"external_network_contacted"`
	BackendCommandExecutionEnabled bool `json:"backend_command_execution_enabled"`
	ExternalEvidenceCreated        bool `json:"external_evidence_created"`
	FinalRootEvidenceCreated       bool `json:"final_root_evidence_created"`
	ConsumerStatusesChanged        bool `json:"consumer_statuses_changed"`
	ComplianceClaimed              bool `json:"compliance_claimed"`
	ProductionReadinessClaimed     bool `json:"production_readiness_claimed"`
	ProductionGradeETAClaimed      bool `json:"production_grade_eta_claimed"`
	RealWorldETAAccuracyClaimed    bool `json:"real_world_eta_accuracy_claimed"`
	VendorCompatibilityClaimed     bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed   bool `json:"hardware_certification_claimed"`
	SLAClaimed                     bool `json:"sla_claimed"`
	HostedSaaSClaimed              bool `json:"hosted_saas_claimed"`
	PublicLaunchClaimed            bool `json:"public_launch_claimed"`
	ConsumerAcceptanceClaimed      bool `json:"consumer_acceptance_claimed"`
	RawObservedRowsPersisted       bool `json:"raw_observed_rows_persisted"`
}

func (h *handler) renderPredictionLab(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "prediction-lab")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "prediction-lab", page)
}

func (h *handler) renderPredictionLabJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "prediction-lab")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.PredictionLab)
}

func buildPredictionLab(page operationsPage) predictionLabView {
	view := predictionLabView{
		GeneratedAt: page.GeneratedAt,
		AgencyID:    page.AgencyID,
		Boundary:    "Private prediction diagnostics only. This page explains why Trip Updates and ETA-like outputs were emitted, withheld, or failed closed. It does not prove production-grade ETA quality, real-world ETA accuracy, consumer acceptance, CAL-ITP/Caltrans compliance, vendor compatibility, SLA coverage, hosted service readiness, or public launch readiness.",
		Commands: []predictionLabCommand{
			{
				ID:             "deterministic-prediction-tests",
				Label:          "Deterministic predictor unit checks",
				CommandLine:    "go test ./internal/prediction -run Deterministic",
				ExpectedResult: "Deterministic prediction and withheld-output behavior stays covered by local tests.",
				DoesNotProve:   "Real-world ETA accuracy, public consumer display, production readiness, or external predictor support.",
			},
			{
				ID:             "realtime-quality-replay",
				Label:          "Realtime quality replay",
				CommandLine:    "make realtime-quality",
				ExpectedResult: "Synthetic replay fixtures keep unknown, stale, ambiguous, low-confidence, and withheld states visible.",
				DoesNotProve:   "Production-grade ETA quality, vendor compatibility, hardware certification, SLA coverage, or compliance.",
			},
		},
		ClaimFlags: predictionLabClaimFlags{},
	}
	view.Summary = buildPredictionLabSummary(page.TripUpdatesQuality)
	view.Deterministic = buildPredictionLabDeterministic(page.TripUpdatesQuality)
	view.WithheldReasons = buildPredictionLabWithheldReasons(page.TripUpdatesQuality)
	view.ShadowReview = buildPredictionLabShadowReview(page.TripUpdatesQuality)
	view.ReviewRows = buildPredictionLabReviewRows(view.Summary, view.WithheldReasons)
	return view
}

func buildPredictionLabSummary(quality tripUpdatesQualityView) predictionLabSummary {
	summary := predictionLabSummary{
		Status:              checklistStatusMissing,
		AdapterName:         firstNonEmpty(quality.AdapterName, "not recorded"),
		DiagnosticsStatus:   firstNonEmpty(quality.DiagnosticsStatus, "not recorded"),
		DiagnosticsReason:   firstNonEmpty(quality.DiagnosticsReason, "not recorded"),
		ActiveFeedVersionID: firstNonEmpty(quality.ActiveFeedVersionID, "not recorded"),
		DoesNotProve:        "Private Trip Updates diagnostics do not prove production-grade ETA quality, real-world ETA accuracy, compliance, consumer acceptance, hosted service, SLA, vendor compatibility, hardware certification, or release readiness.",
	}
	if !quality.Recorded {
		summary.CurrentSignal = firstNonEmpty(quality.Message, "no prediction diagnostics are available yet")
		summary.NextAction = "Send fresh telemetry, confirm an active schedule, then review Realtime Center and Feed Health."
		return summary
	}
	summary.EligiblePredictionCandidates = quality.EligiblePredictionCandidates
	summary.TripUpdatesEmitted = quality.TripUpdatesEmitted
	summary.WithheldCount = totalCountViews(quality.WithheldByReason)
	switch {
	case strings.EqualFold(quality.DiagnosticsStatus, prediction.StatusError):
		summary.Status = checklistStatusBlocked
		summary.CurrentSignal = fmt.Sprintf("adapter=%s failed closed with reason=%s", summary.AdapterName, summary.DiagnosticsReason)
		summary.NextAction = "Review adapter configuration and withheld reasons; Vehicle Positions can remain independent while Trip Updates stay empty or partial."
	case summary.TripUpdatesEmitted == 0 && summary.EligiblePredictionCandidates > 0:
		summary.Status = checklistStatusNeedsReview
		summary.CurrentSignal = fmt.Sprintf("%d eligible candidates, 0 Trip Updates emitted, %d withheld signals", summary.EligiblePredictionCandidates, summary.WithheldCount)
		summary.NextAction = "Review stale telemetry, unknown assignments, confidence thresholds, active feed state, and withheld reason rows."
	case summary.WithheldCount > 0:
		summary.Status = checklistStatusNeedsReview
		summary.CurrentSignal = fmt.Sprintf("%d Trip Updates emitted; %d outputs were withheld or degraded for review", summary.TripUpdatesEmitted, summary.WithheldCount)
		summary.NextAction = "Review withheld reason rows before treating any ETA-like output as useful to operators."
	case summary.TripUpdatesEmitted > 0:
		summary.Status = checklistStatusOK
		summary.CurrentSignal = fmt.Sprintf("%d Trip Updates emitted from %d eligible candidates", summary.TripUpdatesEmitted, summary.EligiblePredictionCandidates)
		summary.NextAction = "Continue monitoring deterministic diagnostics, freshness, assignment confidence, and future-stop coverage."
	default:
		summary.Status = checklistStatusUnknown
		summary.CurrentSignal = fmt.Sprintf("diagnostics=%s/%s with no eligible candidates", summary.DiagnosticsStatus, summary.DiagnosticsReason)
		summary.NextAction = "Confirm active schedule, service day, fresh telemetry, and matched assignments before expecting Trip Updates."
	}
	return summary
}

func buildPredictionLabDeterministic(quality tripUpdatesQualityView) predictionLabDeterministic {
	status := checklistStatusOK
	reviewSignal := "deterministic fallback is the default safe review path"
	nextAction := "Keep deterministic prediction as the baseline and prefer withheld output over false certainty."
	if !quality.Recorded {
		status = checklistStatusMissing
		reviewSignal = "no recorded Trip Updates diagnostics are available yet"
		nextAction = "Generate Trip Updates diagnostics from the private service path before reviewing ETA-like output."
	} else if quality.AdapterName != "" && quality.AdapterName != "deterministic" {
		status = checklistStatusNeedsReview
		reviewSignal = "current diagnostics were not recorded by the deterministic adapter"
		nextAction = "Confirm deterministic fallback remains available before enabling or reviewing external predictor output."
	}
	rows := []predictionLabDiagnosticRow{
		{
			ID:            "vehicle-positions-independent",
			Label:         "Vehicle Positions stay independent",
			Status:        checklistStatusOK,
			CurrentSignal: "Trip Updates diagnostics are reviewed separately from Vehicle Positions feed publication.",
			NextAction:    "Do not block Vehicle Positions on predictor availability or ETA-like output quality.",
			DoesNotProve:  "Vehicle Positions availability does not prove ETA quality, consumer display, compliance, or production readiness.",
		},
		{
			ID:            "assignment-confidence",
			Label:         "Assignment confidence",
			Status:        statusFromRateText(quality.UnknownAssignmentRate, quality.AmbiguousAssignmentRate),
			CurrentSignal: fmt.Sprintf("unknown=%s; ambiguous=%s; manual_overrides=%d", firstNonEmpty(quality.UnknownAssignmentRate, "not recorded"), firstNonEmpty(quality.AmbiguousAssignmentRate, "not recorded"), quality.ManualOverrideAssignments),
			NextAction:    "Keep trip descriptors unknown when assignment confidence is weak, ambiguous, or manually overridden.",
			DoesNotProve:  "A matched assignment does not prove real-world ETA accuracy or consumer display.",
		},
		{
			ID:            "telemetry-freshness",
			Label:         "Telemetry freshness",
			Status:        statusFromRateText(quality.StaleTelemetryRate),
			CurrentSignal: fmt.Sprintf("stale=%s; rows=%d", firstNonEmpty(quality.StaleTelemetryRate, "not recorded"), quality.StaleTelemetryRows),
			NextAction:    "Check device power, network, clocks, and reporting cadence before trusting ETA-like output.",
			DoesNotProve:  "Fresh telemetry does not prove hardware certification, vendor compatibility, SLA, or production AVL reliability.",
		},
		{
			ID:            "coverage",
			Label:         "Trip Updates coverage",
			Status:        statusFromRateText(quality.TripUpdatesCoverageRate, quality.FutureStopCoverageRate),
			CurrentSignal: fmt.Sprintf("coverage=%s; future_stop_coverage=%s", firstNonEmpty(quality.TripUpdatesCoverageRate, "not recorded"), firstNonEmpty(quality.FutureStopCoverageRate, "not recorded")),
			NextAction:    "Use coverage as a local diagnostic only; missing future stops should stay withheld until schedule and telemetry evidence improve.",
			DoesNotProve:  "Coverage percentages do not prove production-grade ETA quality, real-world accuracy, compliance, or consumer acceptance.",
		},
	}
	return predictionLabDeterministic{
		Status:       status,
		Boundary:     "Deterministic prediction remains the default safe baseline. It emits Trip Updates only when schedule, telemetry, assignment, freshness, and confidence checks support it.",
		ReviewSignal: reviewSignal,
		NextAction:   nextAction,
		DoesNotProve: "Deterministic diagnostics do not prove production-grade ETA quality, real-world ETA accuracy, consumer acceptance, or production readiness.",
		Rows:         rows,
	}
}

func buildPredictionLabWithheldReasons(quality tripUpdatesQualityView) []predictionLabReason {
	if !quality.Recorded {
		return []predictionLabReason{{
			Reason:       "no_diagnostics_recorded",
			Label:        "No diagnostics recorded",
			Count:        1,
			WhatItMeans:  firstNonEmpty(quality.Message, "No Trip Updates diagnostics are available yet."),
			NextAction:   "Send fresh telemetry, confirm an active schedule, then review Realtime Center and Feed Health.",
			DoesNotProve: "Missing diagnostics do not prove predictor failure, public feed readiness, consumer display, or ETA quality.",
		}}
	}
	reasons := make([]predictionLabReason, 0, len(quality.WithheldByReason))
	for _, row := range quality.WithheldByReason {
		reasons = append(reasons, predictionLabReasonFor(row.Label, row.Count))
	}
	if len(reasons) == 0 {
		reasons = append(reasons, predictionLabReason{
			Reason:       "none_recorded",
			Label:        "No withheld reasons recorded",
			Count:        0,
			WhatItMeans:  "The latest diagnostics did not include withheld-output reason counts.",
			NextAction:   "Continue monitoring Trip Updates diagnostics, stale telemetry, assignment confidence, and future-stop coverage.",
			DoesNotProve: "A quiet diagnostic row does not prove real-world ETA accuracy or production readiness.",
		})
	}
	return reasons
}

func buildPredictionLabShadowReview(quality tripUpdatesQualityView) predictionLabShadowReview {
	review := predictionLabShadowReview{
		Status:       checklistStatusOK,
		Boundary:     "External predictor review is shadow/fail-closed diagnostics only. The browser does not contact sidecars, start predictors, test URLs, store credentials, or change public Trip Updates output.",
		NextAction:   "Keep external predictors disabled or shadow-only until bounded diagnostics, fail-closed behavior, and rollback are reviewed.",
		DoesNotProve: "Shadow diagnostics do not prove named predictor compatibility, real-world ETA accuracy, production readiness, SLA coverage, consumer acceptance, or compliance.",
	}
	shadow, ok := predictionLabShadowDetails(quality.AdapterDetails)
	if !ok {
		review.Status = checklistStatusUnknown
		review.Rows = []predictionLabShadowRow{{
			ID:              "external-http-shadow-not-recorded",
			Label:           "External HTTP shadow not recorded",
			Status:          "not recorded",
			Reason:          "no bounded external_http_shadow diagnostics are present in the latest Trip Updates record",
			Latency:         "not recorded",
			CountComparison: "deterministic output remains the only reviewed signal",
			FailureBehavior: "Default-safe: keep deterministic fallback and Vehicle Positions independent.",
			FirstSafeCheck:  "go test ./internal/prediction -run ExternalHTTP",
			DoesNotProve:    review.DoesNotProve,
		}}
		return review
	}
	review.Status = statusFromShadowStatus(shadow.Status)
	review.Rows = []predictionLabShadowRow{{
		ID:              "external-http-shadow",
		Label:           "External HTTP shadow",
		Status:          firstNonEmpty(shadow.Status, "unknown"),
		Reason:          firstNonEmpty(shadow.Reason, "not recorded"),
		Latency:         latencyText(shadow.LatencyMS),
		CountComparison: fmt.Sprintf("deterministic=%d; external=%d; delta=%+d", shadow.DeterministicCount, shadow.ExternalCount, shadow.CountDelta),
		FailureBehavior: "Shadow mode keeps deterministic Trip Updates as public output and records bounded diagnostic deltas only.",
		FirstSafeCheck:  "go test ./internal/prediction -run ExternalHTTP",
		DoesNotProve:    review.DoesNotProve,
	}}
	return review
}

type predictionLabShadowDetailsView struct {
	Status             string
	Reason             string
	LatencyMS          int
	DeterministicCount int
	ExternalCount      int
	CountDelta         int
}

func predictionLabShadowDetails(details map[string]any) (predictionLabShadowDetailsView, bool) {
	if len(details) == 0 {
		return predictionLabShadowDetailsView{}, false
	}
	raw, ok := details["external_http_shadow"]
	if !ok {
		return predictionLabShadowDetailsView{}, false
	}
	shadowMap, ok := raw.(map[string]any)
	if !ok {
		return predictionLabShadowDetailsView{}, false
	}
	return predictionLabShadowDetailsView{
		Status:             safeDiagnosticToken(shadowMap["status"]),
		Reason:             safeDiagnosticToken(shadowMap["reason"]),
		LatencyMS:          intFromAny(shadowMap["latency_ms"]),
		DeterministicCount: intFromAny(shadowMap["deterministic_trip_updates_count"]),
		ExternalCount:      intFromAny(shadowMap["external_trip_updates_count"]),
		CountDelta:         intFromAny(shadowMap["count_delta"]),
	}, true
}

func statusFromShadowStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case prediction.StatusOK:
		return checklistStatusOK
	case prediction.StatusError:
		return checklistStatusNeedsReview
	default:
		return checklistStatusUnknown
	}
}

func latencyText(ms int) string {
	if ms <= 0 {
		return "not recorded"
	}
	return fmt.Sprintf("%d ms", ms)
}

func safeDiagnosticToken(value any) string {
	switch typed := value.(type) {
	case string:
		cleaned := strings.TrimSpace(typed)
		if cleaned == "" || strings.ContainsAny(cleaned, "/:\\") || strings.Contains(strings.ToLower(cleaned), "token") || strings.Contains(strings.ToLower(cleaned), "secret") || strings.Contains(strings.ToLower(cleaned), "host") {
			return ""
		}
		return cleaned
	default:
		return ""
	}
}

func intFromAny(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case float32:
		return int(typed)
	default:
		return 0
	}
}

func predictionLabReasonFor(reason string, count int) predictionLabReason {
	base := predictionLabReason{
		Reason:       firstNonEmpty(reason, "unknown"),
		Label:        humanizeReason(reason),
		Count:        count,
		WhatItMeans:  "The predictor withheld output because the available evidence was not safe enough for ETA-like Trip Updates.",
		NextAction:   "Review schedule state, fresh telemetry, assignment confidence, and active feed diagnostics before changing prediction settings.",
		DoesNotProve: "A withheld reason does not prove production-grade ETA quality, real-world accuracy, consumer display, or compliance.",
	}
	switch reason {
	case prediction.ReasonNoLatestTelemetry:
		base.WhatItMeans = "No latest telemetry row was available for the vehicle."
		base.NextAction = "Check device credentials, simulator/device reporting, and telemetry ingest before expecting Trip Updates."
	case prediction.ReasonStaleTelemetry:
		base.WhatItMeans = "The last accepted telemetry was too old for safe ETA-like output."
		base.NextAction = "Check vehicle device power, network, clock sync, and reporting cadence."
	case prediction.ReasonDegradedAssignment:
		base.WhatItMeans = "The current assignment was degraded, ambiguous, or not strong enough to publish confident Trip Updates."
		base.NextAction = "Review route/trip hints, service day, after-midnight service, block continuity, and operator overrides."
	case prediction.ReasonBelowConfidenceThreshold:
		base.WhatItMeans = "Assignment confidence was below the configured prediction threshold."
		base.NextAction = "Improve telemetry quality or apply an explicit operator override when staff have reliable trip evidence."
	case prediction.ReasonScheduleUnavailable:
		base.WhatItMeans = "The active schedule data needed for future stop times was not available."
		base.NextAction = "Review GTFS Workbench, active feed version, and schedule validation before relying on Trip Updates."
	case prediction.ReasonNoFutureStops:
		base.WhatItMeans = "No future stop updates were safe to emit for the current trip state."
		base.NextAction = "Review current stop sequence, service day, and schedule timing before expecting ETA-like output."
	case prediction.ReasonDuplicateTripInstance:
		base.WhatItMeans = "Multiple candidates pointed at the same trip instance and the safer behavior was to withhold uncertain output."
		base.NextAction = "Review block continuity, repeated trip instances, and assignment confidence before publishing trip certainty."
	case prediction.ReasonCanceledTripRequiresAlert:
		base.WhatItMeans = "A canceled trip requires an Alerts lifecycle review before a cancellation-style Trip Update is safe."
		base.NextAction = "Review the Alerts Console and feed health before relying on cancellation Trip Updates."
	}
	return base
}

func buildPredictionLabReviewRows(summary predictionLabSummary, reasons []predictionLabReason) []predictionLabReviewRow {
	rows := []predictionLabReviewRow{}
	if summary.Status != checklistStatusOK {
		rows = append(rows, predictionLabReviewRow{
			Severity:     summary.Status,
			Area:         "Trip Updates decision",
			Signal:       summary.CurrentSignal,
			NextAction:   summary.NextAction,
			DoesNotProve: summary.DoesNotProve,
		})
	}
	for _, reason := range reasons {
		if reason.Count <= 0 || reason.Reason == "none_recorded" {
			continue
		}
		rows = append(rows, predictionLabReviewRow{
			Severity:     checklistStatusNeedsReview,
			Area:         reason.Label,
			Signal:       fmt.Sprintf("%d rows withheld for %s", reason.Count, reason.Reason),
			NextAction:   reason.NextAction,
			DoesNotProve: reason.DoesNotProve,
		})
	}
	if len(rows) == 0 {
		rows = append(rows, predictionLabReviewRow{
			Severity:     checklistStatusOK,
			Area:         "Current prediction review",
			Signal:       "No active withheld-output review rows are visible in this bounded summary.",
			NextAction:   "Continue monitoring deterministic diagnostics and local aggregate backtests when available.",
			DoesNotProve: "No visible review rows does not prove production-grade ETA quality, real-world accuracy, or consumer display.",
		})
	}
	const limit = 10
	if len(rows) > limit {
		rows = append(rows[:limit], predictionLabReviewRow{
			Severity:     checklistStatusNeedsReview,
			Area:         "Bounded review list",
			Signal:       fmt.Sprintf("%d additional review rows are hidden from this browser summary", len(rows)-limit),
			NextAction:   "Use the private JSON export and underlying diagnostics for deeper local review.",
			DoesNotProve: "A bounded issue count is not an exhaustive operational proof.",
		})
	}
	return rows
}

func statusFromRateText(values ...string) string {
	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		if normalized == "" || normalized == "not recorded" || strings.Contains(normalized, "not applicable") {
			continue
		}
		if strings.HasPrefix(normalized, "0.0%") || strings.HasPrefix(normalized, "0%") {
			continue
		}
		return checklistStatusNeedsReview
	}
	return checklistStatusOK
}

func totalCountViews(rows []countView) int {
	total := 0
	for _, row := range rows {
		total += row.Count
	}
	return total
}

func humanizeReason(reason string) string {
	trimmed := strings.TrimSpace(reason)
	if trimmed == "" {
		return "Unknown reason"
	}
	parts := strings.Split(trimmed, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}
