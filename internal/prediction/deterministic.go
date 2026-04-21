package prediction

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

const (
	ReasonNoLatestTelemetry           = "no_latest_telemetry"
	ReasonStaleTelemetry              = "stale_telemetry"
	ReasonAssignmentTelemetryMismatch = "assignment_telemetry_mismatch"
	ReasonAssignmentFeedMismatch      = "assignment_feed_mismatch"
	ReasonNotInService                = "not_in_service"
	ReasonDeadheadNoPrediction        = "deadhead_no_prediction"
	ReasonLayoverNoPrediction         = "layover_no_prediction"
	ReasonMissingTripIdentity         = "missing_trip_identity"
	ReasonDegradedAssignment          = "degraded_assignment"
	ReasonBelowConfidenceThreshold    = "below_confidence_threshold"
	ReasonMissingCurrentStop          = "missing_current_stop_sequence"
	ReasonScheduleUnavailable         = "schedule_unavailable"
	ReasonCurrentStopNotInSchedule    = "current_stop_not_in_schedule"
	ReasonScheduleDeviationTooLarge   = "schedule_deviation_too_large"
	ReasonNoFutureStops               = "no_future_stops"
	ReasonDuplicateTripInstance       = "duplicate_trip_instance"
	ReasonCanceledTripRequiresAlert   = "canceled_trip_requires_service_alert"
	ReasonAddedTripWithheld           = "added_trip_withheld"
	ReasonShortTurnWithheld           = "short_turn_withheld"
	ReasonDetourPredictionWithheld    = "detour_prediction_withheld"
)

type DeterministicConfig struct {
	StaleTelemetryTTL       time.Duration
	AssignmentConfidenceMin float64
	MaxScheduleDeviation    time.Duration
	DuplicateConfidenceGap  float64
}

func DefaultDeterministicConfig() DeterministicConfig {
	return DeterministicConfig{
		StaleTelemetryTTL:       90 * time.Second,
		AssignmentConfidenceMin: state.DefaultConfig().MinConfidence,
		MaxScheduleDeviation:    45 * time.Minute,
		DuplicateConfidenceGap:  0.05,
	}
}

func (c DeterministicConfig) validated() DeterministicConfig {
	defaults := DefaultDeterministicConfig()
	if c.StaleTelemetryTTL == 0 {
		c.StaleTelemetryTTL = defaults.StaleTelemetryTTL
	}
	if c.AssignmentConfidenceMin == 0 {
		c.AssignmentConfidenceMin = defaults.AssignmentConfidenceMin
	}
	if c.MaxScheduleDeviation == 0 {
		c.MaxScheduleDeviation = defaults.MaxScheduleDeviation
	}
	if c.DuplicateConfidenceGap == 0 {
		c.DuplicateConfidenceGap = defaults.DuplicateConfidenceGap
	}
	return c
}

type DeterministicAdapter struct {
	schedules  gtfs.Repository
	operations OperationsRepository
	config     DeterministicConfig
}

func NewDeterministicAdapter(schedules gtfs.Repository, operations OperationsRepository, config DeterministicConfig) (*DeterministicAdapter, error) {
	if schedules == nil {
		return nil, fmt.Errorf("gtfs repository is required")
	}
	return &DeterministicAdapter{schedules: schedules, operations: operations, config: config.validated()}, nil
}

func (a *DeterministicAdapter) Name() string {
	return "deterministic"
}

func (a *DeterministicAdapter) PredictTripUpdates(ctx context.Context, request Request) (Result, error) {
	agency, err := a.schedules.Agency(ctx, request.AgencyID)
	if err != nil {
		return Result{}, fmt.Errorf("load agency for prediction: %w", err)
	}
	loc, err := time.LoadLocation(agency.Timezone)
	if err != nil {
		return Result{}, fmt.Errorf("load agency timezone for prediction: %w", err)
	}

	metrics := Metrics{
		TelemetryRowsConsidered: len(request.Telemetry),
		AssignmentsConsidered:   len(request.Assignments),
		WithheldByReason:        map[string]int{},
		DegradedByReason:        map[string]int{},
	}
	details := map[string]any{
		"coverage_denominator": "eligible in-service ETA prediction candidates; canceled trips are excluded and tracked separately",
		"review_lifecycle":     "prediction review items use open/resolved/deferred status",
	}
	reviews := []ReviewItem{}

	eventsByVehicle := make(map[string]telemetry.StoredEvent, len(request.Telemetry))
	for _, event := range request.Telemetry {
		eventsByVehicle[event.VehicleID] = event
	}

	overridesByVehicle := map[string][]OverrideRecord{}
	var disruptionUpdates []TripUpdate
	if a.operations != nil {
		overrides, err := a.operations.ListActivePredictionOverrides(ctx, request.AgencyID, request.GeneratedAt)
		if err != nil {
			return Result{}, fmt.Errorf("list active prediction overrides: %w", err)
		}
		for _, override := range overrides {
			overridesByVehicle[override.VehicleID] = append(overridesByVehicle[override.VehicleID], override)
		}
		disruptionUpdates = a.applyDisruptionOverrides(request, overrides, &metrics, &reviews)
	}

	cache := scheduleCache{adapter: a, request: request}
	var candidates []predictedCandidate
	vehicleIDs := sortedAssignmentVehicles(request.Assignments)
	for _, vehicleID := range vehicleIDs {
		if hasBlockingDisruption(overridesByVehicle[vehicleID]) {
			continue
		}
		candidate, ok := a.evaluateAssignment(ctx, request, loc, &cache, eventsByVehicle, vehicleID, request.Assignments[vehicleID], &metrics, &reviews)
		if ok {
			candidates = append(candidates, candidate)
		}
	}

	updates := append(disruptionUpdates, a.resolveDuplicateTripInstances(request, candidates, &metrics, &reviews)...)
	sort.SliceStable(updates, func(i, j int) bool {
		return tripUpdateKey(updates[i]) < tripUpdateKey(updates[j])
	})
	metrics.TripUpdatesEmitted = len(updates)
	for _, update := range updates {
		metrics.StopUpdatesEmitted += len(update.StopTimeUpdates)
	}
	metrics.CoveragePercent = percent(countETAUpdates(updates), metrics.EligiblePredictionCandidates)
	metrics.FutureStopCoveragePercent = percent(countUpdatesWithFutureStops(updates), metrics.EligiblePredictionCandidates)

	if a.operations != nil && len(reviews) > 0 {
		if err := a.operations.SavePredictionReviewItems(ctx, reviews); err != nil {
			details["review_persistence_error"] = err.Error()
		}
	}

	status := StatusOK
	reason := ReasonNoEligiblePredictions
	switch {
	case len(updates) > 0 && totalReasons(metrics.WithheldByReason) > 0:
		reason = ReasonPartialPredictions
	case len(updates) > 0:
		reason = ReasonPredictionsAvailable
	}
	return Result{
		TripUpdates: updates,
		Diagnostics: Diagnostics{
			Status:  status,
			Reason:  reason,
			Metrics: metrics,
			Details: details,
		},
	}, nil
}

type scheduleCache struct {
	adapter   *DeterministicAdapter
	request   Request
	byDate    map[string][]gtfs.TripCandidate
	loadError error
}

func (c *scheduleCache) trips(ctx context.Context, serviceDate string) ([]gtfs.TripCandidate, error) {
	if c.byDate == nil {
		c.byDate = map[string][]gtfs.TripCandidate{}
	}
	if trips, ok := c.byDate[serviceDate]; ok {
		return trips, nil
	}
	trips, err := c.adapter.schedules.ListTripCandidates(ctx, c.request.AgencyID, c.request.ActiveFeedVersion.ID, serviceDate)
	if err != nil {
		c.loadError = err
		return nil, err
	}
	c.byDate[serviceDate] = trips
	return trips, nil
}

type predictedCandidate struct {
	update     TripUpdate
	assignment state.Assignment
}

type scheduleInstance struct {
	trip                 gtfs.TripCandidate
	stopTimes            []gtfs.StopTime
	scheduleRelationship ScheduleRelationship
}

func (a *DeterministicAdapter) evaluateAssignment(
	ctx context.Context,
	request Request,
	loc *time.Location,
	cache *scheduleCache,
	eventsByVehicle map[string]telemetry.StoredEvent,
	vehicleID string,
	assignment state.Assignment,
	metrics *Metrics,
	reviews *[]ReviewItem,
) (predictedCandidate, bool) {
	event, ok := eventsByVehicle[vehicleID]
	if !ok {
		a.withhold(request, assignment, nil, ReasonNoLatestTelemetry, metrics, reviews)
		return predictedCandidate{}, false
	}
	if request.GeneratedAt.Sub(event.Timestamp) > a.config.StaleTelemetryTTL {
		a.withhold(request, assignment, &event, ReasonStaleTelemetry, metrics, reviews)
		return predictedCandidate{}, false
	}
	if assignment.State == state.StateDeadhead {
		a.withhold(request, assignment, &event, ReasonDeadheadNoPrediction, metrics, reviews)
		return predictedCandidate{}, false
	}
	if assignment.State == state.StateLayover {
		a.withhold(request, assignment, &event, ReasonLayoverNoPrediction, metrics, reviews)
		return predictedCandidate{}, false
	}
	if assignment.State != state.StateInService {
		a.withhold(request, assignment, &event, ReasonNotInService, metrics, reviews)
		return predictedCandidate{}, false
	}
	if assignment.TripID == "" || assignment.StartDate == "" || assignment.StartTime == "" {
		a.withhold(request, assignment, &event, ReasonMissingTripIdentity, metrics, reviews)
		return predictedCandidate{}, false
	}
	if assignment.FeedVersionID != "" && assignment.FeedVersionID != request.ActiveFeedVersion.ID {
		a.withhold(request, assignment, &event, ReasonAssignmentFeedMismatch, metrics, reviews)
		return predictedCandidate{}, false
	}
	if assignment.AssignmentSource == state.AssignmentSourceAutomatic && assignment.TelemetryEventID != 0 && assignment.TelemetryEventID != event.ID {
		a.withhold(request, assignment, &event, ReasonAssignmentTelemetryMismatch, metrics, reviews)
		return predictedCandidate{}, false
	}
	if assignment.DegradedState != "" && assignment.DegradedState != state.DegradedNone {
		a.withhold(request, assignment, &event, ReasonDegradedAssignment, metrics, reviews)
		metrics.DegradedByReason[string(assignment.DegradedState)]++
		return predictedCandidate{}, false
	}
	if assignment.AssignmentSource != state.AssignmentSourceManualOverride && assignment.Confidence < a.config.AssignmentConfidenceMin {
		a.withhold(request, assignment, &event, ReasonBelowConfidenceThreshold, metrics, reviews)
		return predictedCandidate{}, false
	}
	if assignment.CurrentStopSequence <= 0 {
		a.withhold(request, assignment, &event, ReasonMissingCurrentStop, metrics, reviews)
		return predictedCandidate{}, false
	}

	instance, err := cache.instance(ctx, assignment)
	if err != nil {
		a.withhold(request, assignment, &event, ReasonScheduleUnavailable, metrics, reviews)
		return predictedCandidate{}, false
	}
	currentStop, ok := stopBySequence(instance.stopTimes, assignment.CurrentStopSequence)
	if !ok {
		metrics.EligiblePredictionCandidates++
		a.withhold(request, assignment, &event, ReasonCurrentStopNotInSchedule, metrics, reviews)
		return predictedCandidate{}, false
	}
	metrics.EligiblePredictionCandidates++

	scheduledCurrent := serviceTime(assignment.StartDate, currentStop.DepartureSeconds, loc)
	delay := event.Timestamp.Sub(scheduledCurrent)
	if math.Abs(delay.Seconds()) > a.config.MaxScheduleDeviation.Seconds() {
		a.withhold(request, assignment, &event, ReasonScheduleDeviationTooLarge, metrics, reviews)
		return predictedCandidate{}, false
	}

	stopUpdates := make([]StopTimeUpdate, 0, len(instance.stopTimes))
	delaySeconds := int32(math.Round(delay.Seconds()))
	for _, stop := range instance.stopTimes {
		if stop.StopSequence <= assignment.CurrentStopSequence {
			continue
		}
		arrival := serviceTime(assignment.StartDate, stop.ArrivalSeconds, loc).Add(delay)
		departure := serviceTime(assignment.StartDate, stop.DepartureSeconds, loc).Add(delay)
		if !arrival.After(request.GeneratedAt) && !departure.After(request.GeneratedAt) {
			continue
		}
		stopUpdates = append(stopUpdates, StopTimeUpdate{
			StopID:                stop.StopID,
			StopSequence:          stop.StopSequence,
			ArrivalTime:           &arrival,
			DepartureTime:         &departure,
			ArrivalDelaySeconds:   &delaySeconds,
			DepartureDelaySeconds: &delaySeconds,
			ScheduleRelationship:  ScheduleRelationshipScheduled,
		})
	}
	if len(stopUpdates) == 0 {
		a.withhold(request, assignment, &event, ReasonNoFutureStops, metrics, reviews)
		return predictedCandidate{}, false
	}

	update := TripUpdate{
		EntityID:             tripUpdateEntityID(assignment.TripID, assignment.StartDate, assignment.StartTime),
		VehicleID:            vehicleID,
		TripID:               assignment.TripID,
		RouteID:              firstNonEmpty(assignment.RouteID, instance.trip.RouteID),
		StartDate:            assignment.StartDate,
		StartTime:            assignment.StartTime,
		ScheduleRelationship: instance.scheduleRelationship,
		StopTimeUpdates:      stopUpdates,
	}
	return predictedCandidate{update: update, assignment: assignment}, true
}

func (c *scheduleCache) instance(ctx context.Context, assignment state.Assignment) (scheduleInstance, error) {
	trips, err := c.trips(ctx, assignment.StartDate)
	if err != nil {
		return scheduleInstance{}, err
	}
	for _, trip := range trips {
		if trip.TripID != assignment.TripID {
			continue
		}
		instance, ok := matchTripInstance(trip, assignment.StartTime)
		if ok {
			return instance, nil
		}
	}
	return scheduleInstance{}, fmt.Errorf("trip instance not found")
}

func matchTripInstance(trip gtfs.TripCandidate, startTime string) (scheduleInstance, bool) {
	if len(trip.StopTimes) == 0 {
		return scheduleInstance{}, false
	}
	baseStart := trip.StopTimes[0].DepartureSeconds
	if len(trip.Frequencies) == 0 {
		if startTime != "" && startTime != gtfs.FormatGTFSTime(baseStart) {
			return scheduleInstance{}, false
		}
		return scheduleInstance{trip: trip, stopTimes: cloneShiftedStopTimes(trip.StopTimes, 0), scheduleRelationship: ScheduleRelationshipScheduled}, true
	}
	for _, frequency := range trip.Frequencies {
		if frequency.ExactTimes == 1 {
			for start := frequency.StartSeconds; start < frequency.EndSeconds; start += frequency.HeadwaySecs {
				if gtfs.FormatGTFSTime(start) == startTime {
					return scheduleInstance{trip: trip, stopTimes: cloneShiftedStopTimes(trip.StopTimes, start-baseStart), scheduleRelationship: ScheduleRelationshipScheduled}, true
				}
			}
			continue
		}
		if frequency.StartTime == startTime {
			return scheduleInstance{trip: trip, stopTimes: cloneShiftedStopTimes(trip.StopTimes, frequency.StartSeconds-baseStart), scheduleRelationship: ScheduleRelationshipUnscheduled}, true
		}
	}
	return scheduleInstance{}, false
}

func cloneShiftedStopTimes(stopTimes []gtfs.StopTime, shift int) []gtfs.StopTime {
	shifted := make([]gtfs.StopTime, len(stopTimes))
	for i, stop := range stopTimes {
		shifted[i] = stop
		shifted[i].ArrivalSeconds += shift
		shifted[i].DepartureSeconds += shift
	}
	return shifted
}

func (a *DeterministicAdapter) applyDisruptionOverrides(request Request, overrides []OverrideRecord, metrics *Metrics, reviews *[]ReviewItem) []TripUpdate {
	var updates []TripUpdate
	for _, override := range overrides {
		switch override.OverrideType {
		case "canceled_trip":
			if override.TripID == "" || override.StartDate == "" || override.StartTime == "" {
				a.withholdOverride(request, override, ReasonMissingTripIdentity, metrics, reviews)
				continue
			}
			metrics.CanceledTripsEmitted++
			metrics.CancellationAlertLinksExpected++
			metrics.CancellationAlertLinksMissing++
			updates = append(updates, TripUpdate{
				EntityID:             tripUpdateEntityID(override.TripID, override.StartDate, override.StartTime),
				VehicleID:            override.VehicleID,
				TripID:               override.TripID,
				RouteID:              override.RouteID,
				StartDate:            override.StartDate,
				StartTime:            override.StartTime,
				ScheduleRelationship: ScheduleRelationshipCanceled,
			})
			*reviews = append(*reviews, ReviewItem{
				AgencyID:   request.AgencyID,
				SnapshotAt: request.GeneratedAt,
				VehicleID:  override.VehicleID,
				RouteID:    override.RouteID,
				TripID:     override.TripID,
				StartDate:  override.StartDate,
				StartTime:  override.StartTime,
				Severity:   "warning",
				Reason:     ReasonCanceledTripRequiresAlert,
				Status:     ReviewStatusOpen,
				Details: map[string]any{
					"expected_alert_missing":                true,
					"cancellation_alert_linkage_status":     "missing_alert_authoring_deferred",
					"linked_review_reason":                  ReasonCanceledTripRequiresAlert,
					"prediction_override_id":                override.ID,
					"coverage_denominator_included":         false,
					"coverage_denominator_exclusion_reason": "canceled_trip_metric_tracked_separately",
				},
			})
		case "added_trip":
			metrics.AddedTripsWithheld++
			a.withholdOverride(request, override, ReasonAddedTripWithheld, metrics, reviews)
		case "short_turn":
			metrics.ShortTurnsWithheld++
			a.withholdOverride(request, override, ReasonShortTurnWithheld, metrics, reviews)
		case "detour":
			metrics.DetoursWithheld++
			a.withholdOverride(request, override, ReasonDetourPredictionWithheld, metrics, reviews)
		}
	}
	return updates
}

func (a *DeterministicAdapter) resolveDuplicateTripInstances(request Request, candidates []predictedCandidate, metrics *Metrics, reviews *[]ReviewItem) []TripUpdate {
	byKey := map[string][]predictedCandidate{}
	for _, candidate := range candidates {
		byKey[tripUpdateKey(candidate.update)] = append(byKey[tripUpdateKey(candidate.update)], candidate)
	}
	var keys []string
	for key := range byKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var updates []TripUpdate
	for _, key := range keys {
		group := byKey[key]
		if len(group) == 1 {
			updates = append(updates, group[0].update)
			continue
		}
		sort.SliceStable(group, func(i, j int) bool {
			leftManual := group[i].assignment.AssignmentSource == state.AssignmentSourceManualOverride
			rightManual := group[j].assignment.AssignmentSource == state.AssignmentSourceManualOverride
			if leftManual != rightManual {
				return leftManual
			}
			return group[i].assignment.Confidence > group[j].assignment.Confidence
		})
		if group[0].assignment.AssignmentSource == state.AssignmentSourceManualOverride && group[1].assignment.AssignmentSource != state.AssignmentSourceManualOverride {
			updates = append(updates, group[0].update)
			continue
		}
		if group[0].assignment.Confidence-group[1].assignment.Confidence >= a.config.DuplicateConfidenceGap {
			updates = append(updates, group[0].update)
			continue
		}
		for _, candidate := range group {
			a.withhold(request, candidate.assignment, nil, ReasonDuplicateTripInstance, metrics, reviews)
		}
	}
	return updates
}

func hasBlockingDisruption(overrides []OverrideRecord) bool {
	for _, override := range overrides {
		switch override.OverrideType {
		case "added_trip", "short_turn", "detour", "canceled_trip":
			return true
		}
	}
	return false
}

func (a *DeterministicAdapter) withhold(request Request, assignment state.Assignment, event *telemetry.StoredEvent, reason string, metrics *Metrics, reviews *[]ReviewItem) {
	metrics.WithheldByReason[reason]++
	details := map[string]any{
		"reason":                    reason,
		"assignment_id":             assignment.ID,
		"assignment_source":         assignment.AssignmentSource,
		"degraded_state":            assignment.DegradedState,
		"coverage_denominator_note": "only ETA-eligible in-service prediction candidates count in coverage percent",
	}
	if event != nil {
		details["telemetry_event_id"] = event.ID
		details["observed_at"] = event.Timestamp.UTC().Format(time.RFC3339)
	}
	*reviews = append(*reviews, ReviewItem{
		AgencyID:   request.AgencyID,
		SnapshotAt: request.GeneratedAt,
		VehicleID:  assignment.VehicleID,
		RouteID:    assignment.RouteID,
		TripID:     assignment.TripID,
		StartDate:  assignment.StartDate,
		StartTime:  assignment.StartTime,
		Severity:   reviewSeverity(reason),
		Reason:     reason,
		Status:     ReviewStatusOpen,
		Details:    details,
	})
}

func (a *DeterministicAdapter) withholdOverride(request Request, override OverrideRecord, reason string, metrics *Metrics, reviews *[]ReviewItem) {
	metrics.WithheldByReason[reason]++
	*reviews = append(*reviews, ReviewItem{
		AgencyID:   request.AgencyID,
		SnapshotAt: request.GeneratedAt,
		VehicleID:  override.VehicleID,
		RouteID:    override.RouteID,
		TripID:     override.TripID,
		StartDate:  override.StartDate,
		StartTime:  override.StartTime,
		Severity:   reviewSeverity(reason),
		Reason:     reason,
		Status:     ReviewStatusOpen,
		Details: map[string]any{
			"prediction_override_id": override.ID,
			"override_type":          override.OverrideType,
			"override_state":         override.State,
		},
	})
}

func reviewSeverity(reason string) string {
	switch reason {
	case ReasonDeadheadNoPrediction, ReasonLayoverNoPrediction, ReasonNotInService:
		return "info"
	default:
		return "warning"
	}
}

func serviceTime(serviceDate string, seconds int, loc *time.Location) time.Time {
	day, err := time.ParseInLocation("20060102", serviceDate, loc)
	if err != nil {
		return time.Unix(0, 0).UTC()
	}
	return day.Add(time.Duration(seconds) * time.Second).UTC()
}

func stopBySequence(stopTimes []gtfs.StopTime, sequence int) (gtfs.StopTime, bool) {
	for _, stop := range stopTimes {
		if stop.StopSequence == sequence {
			return stop, true
		}
	}
	return gtfs.StopTime{}, false
}

func sortedAssignmentVehicles(assignments map[string]state.Assignment) []string {
	vehicleIDs := make([]string, 0, len(assignments))
	for vehicleID := range assignments {
		vehicleIDs = append(vehicleIDs, vehicleID)
	}
	sort.Strings(vehicleIDs)
	return vehicleIDs
}

func tripUpdateEntityID(tripID string, startDate string, startTime string) string {
	return "trip_update:" + tripID + ":" + startDate + ":" + startTime
}

func tripUpdateKey(update TripUpdate) string {
	if update.EntityID != "" {
		return update.EntityID
	}
	return tripUpdateEntityID(update.TripID, update.StartDate, update.StartTime)
}

func percent(numerator int, denominator int) *float64 {
	if denominator <= 0 {
		return nil
	}
	value := float64(numerator) / float64(denominator) * 100
	return &value
}

func countUpdatesWithFutureStops(updates []TripUpdate) int {
	count := 0
	for _, update := range updates {
		if update.ScheduleRelationship != ScheduleRelationshipCanceled && len(update.StopTimeUpdates) > 0 {
			count++
		}
	}
	return count
}

func countETAUpdates(updates []TripUpdate) int {
	count := 0
	for _, update := range updates {
		if update.ScheduleRelationship != ScheduleRelationshipCanceled {
			count++
		}
	}
	return count
}

func totalReasons(reasons map[string]int) int {
	total := 0
	for _, count := range reasons {
		total += count
	}
	return total
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
