package state

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/telemetry"
)

type Config struct {
	StaleThreshold        time.Duration
	MinConfidence         float64
	AmbiguityGap          float64
	ScheduleFitWindow     time.Duration
	ShapeDistanceMeters   float64
	ContinuityWindow      time.Duration
	BlockTransitionWindow time.Duration
}

func DefaultConfig() Config {
	return Config{
		StaleThreshold:        90 * time.Second,
		MinConfidence:         0.65,
		AmbiguityGap:          0.15,
		ScheduleFitWindow:     45 * time.Minute,
		ShapeDistanceMeters:   500,
		ContinuityWindow:      20 * time.Minute,
		BlockTransitionWindow: 20 * time.Minute,
	}
}

type Engine struct {
	schedules   gtfs.Repository
	assignments Repository
	config      Config
}

func NewEngine(schedules gtfs.Repository, assignments Repository, config Config) *Engine {
	return &Engine{schedules: schedules, assignments: assignments, config: mergeConfig(config)}
}

func mergeConfig(config Config) Config {
	defaults := DefaultConfig()
	if config.StaleThreshold == 0 {
		config.StaleThreshold = defaults.StaleThreshold
	}
	if config.MinConfidence == 0 {
		config.MinConfidence = defaults.MinConfidence
	}
	if config.AmbiguityGap == 0 {
		config.AmbiguityGap = defaults.AmbiguityGap
	}
	if config.ScheduleFitWindow == 0 {
		config.ScheduleFitWindow = defaults.ScheduleFitWindow
	}
	if config.ShapeDistanceMeters == 0 {
		config.ShapeDistanceMeters = defaults.ShapeDistanceMeters
	}
	if config.ContinuityWindow == 0 {
		config.ContinuityWindow = defaults.ContinuityWindow
	}
	if config.BlockTransitionWindow == 0 {
		config.BlockTransitionWindow = defaults.BlockTransitionWindow
	}
	return config
}

func (e *Engine) MatchEvent(ctx context.Context, event telemetry.StoredEvent, now time.Time) (Assignment, error) {
	agency, agencyErr := e.schedules.Agency(ctx, event.AgencyID)
	serviceDays, serviceErr := []gtfs.ServiceDay(nil), error(nil)
	if agencyErr == nil {
		serviceDays, serviceErr = gtfs.ResolveServiceDays(event.Timestamp, agency.Timezone)
	}
	defaultServiceDate := ""
	defaultObservedLocalSeconds := any(nil)
	if len(serviceDays) > 0 {
		defaultServiceDate = serviceDays[0].Date
		defaultObservedLocalSeconds = serviceDays[0].ObservedLocalSeconds
	}

	if agencyErr != nil {
		return e.saveUnknown(ctx, event, defaultServiceDate, DegradedUnknown, []string{ReasonAgencyLookupFailed}, IncidentMatcherSystemFailure, map[string]any{
			"agency_error":           errorString(agencyErr),
			"service_date_known":     defaultServiceDate != "",
			"observed_local_seconds": defaultObservedLocalSeconds,
		})
	}
	if serviceErr != nil {
		return e.saveUnknown(ctx, event, defaultServiceDate, DegradedUnknown, []string{ReasonServiceDayFailed}, IncidentMatcherSystemFailure, map[string]any{
			"service_day_error":      errorString(serviceErr),
			"service_date_known":     defaultServiceDate != "",
			"observed_local_seconds": defaultObservedLocalSeconds,
		})
	}

	if now.Sub(event.Timestamp) > e.config.StaleThreshold {
		return e.saveUnknown(ctx, event, defaultServiceDate, DegradedStale, []string{ReasonStaleTelemetry}, IncidentStaleTelemetry, map[string]any{
			"observed_at":             event.Timestamp,
			"observed_local_seconds":  defaultObservedLocalSeconds,
			"now":                     now,
			"stale_threshold_seconds": e.config.StaleThreshold.Seconds(),
			"score_schema":            "loose_debug_v1",
		})
	}

	if e.assignments != nil {
		override, err := e.assignments.ActiveManualOverride(ctx, event.AgencyID, event.VehicleID, event.Timestamp)
		if err != nil {
			return Assignment{}, err
		}
		if override != nil {
			return e.saveManualOverride(ctx, event, *override, defaultServiceDate)
		}
	}

	feed, err := e.schedules.ActiveFeedVersion(ctx, event.AgencyID)
	if err != nil {
		return e.saveUnknown(ctx, event, defaultServiceDate, DegradedMissingScheduleData, []string{ReasonActiveFeedUnavailable}, IncidentMatcherSystemFailure, map[string]any{
			"active_feed_error":      errorString(err),
			"observed_local_seconds": defaultObservedLocalSeconds,
		})
	}

	var previous *Assignment
	if e.assignments != nil {
		previous, err = e.assignments.CurrentAssignment(ctx, event.AgencyID, event.VehicleID)
		if err != nil {
			return Assignment{}, err
		}
	}

	var scored []scoredCandidate
	for _, day := range serviceDays {
		trips, err := e.schedules.ListTripCandidates(ctx, event.AgencyID, feed.ID, day.Date)
		if err != nil {
			return e.saveUnknown(ctx, event, day.Date, DegradedMissingScheduleData, []string{ReasonScheduleQueryFailed}, IncidentMatcherSystemFailure, map[string]any{
				"feed_version_id":        feed.ID,
				"schedule_query_error":   errorString(err),
				"observed_local_seconds": day.ObservedLocalSeconds,
			})
		}
		for _, trip := range trips {
			instances := expandInstances(trip, day.ObservedLocalSeconds, e.config.ScheduleFitWindow)
			for _, instance := range instances {
				scored = append(scored, e.scoreCandidate(event, feed.ID, day, instance, previous))
			}
		}
	}

	if len(scored) == 0 {
		return e.saveUnknown(ctx, event, defaultServiceDate, DegradedMissingScheduleData, []string{ReasonNoScheduleCandidates}, IncidentMissingScheduleData, map[string]any{
			"feed_version_id":        feed.ID,
			"observed_local_seconds": defaultObservedLocalSeconds,
		})
	}

	sort.Slice(scored, func(i int, j int) bool {
		return scored[i].assignment.Confidence > scored[j].assignment.Confidence
	})
	best := scored[0]
	if len(scored) > 1 && best.assignment.Confidence-scored[1].assignment.Confidence < e.config.AmbiguityGap {
		return e.saveUnknown(ctx, event, best.assignment.ServiceDate, DegradedAmbiguous, []string{ReasonAmbiguousCandidates}, IncidentAssignmentAmbiguous, map[string]any{
			"trip_id":                best.assignment.TripID,
			"start_time":             best.assignment.StartTime,
			"observed_local_seconds": best.assignment.ScoreDetails["observed_local_seconds"],
			"best_trip_id":           best.assignment.TripID,
			"best_start_time":        best.assignment.StartTime,
			"best_confidence":        best.assignment.Confidence,
			"second_trip_id":         scored[1].assignment.TripID,
			"second_start_time":      scored[1].assignment.StartTime,
			"second_confidence":      scored[1].assignment.Confidence,
			"ambiguity_gap":          e.config.AmbiguityGap,
		})
	}
	if best.assignment.Confidence < e.config.MinConfidence {
		return e.saveUnknown(ctx, event, best.assignment.ServiceDate, DegradedLowConfidence, append(best.assignment.ReasonCodes, ReasonLowConfidence), IncidentLowConfidenceAssignment, map[string]any{
			"trip_id":                best.assignment.TripID,
			"start_time":             best.assignment.StartTime,
			"observed_local_seconds": best.assignment.ScoreDetails["observed_local_seconds"],
			"best_trip_id":           best.assignment.TripID,
			"best_start_time":        best.assignment.StartTime,
			"best_confidence":        best.assignment.Confidence,
			"minimum_confidence":     e.config.MinConfidence,
		})
	}

	return e.saveAssignment(ctx, best.assignment, nil)
}

type tripInstance struct {
	trip              gtfs.TripCandidate
	startTime         string
	startSeconds      int
	stopTimes         []gtfs.StopTime
	exactFrequency    bool
	nonExactFrequency bool
}

type scoredCandidate struct {
	assignment Assignment
}

func expandInstances(trip gtfs.TripCandidate, observedSeconds int, window time.Duration) []tripInstance {
	if len(trip.StopTimes) == 0 {
		return []tripInstance{{
			trip:      trip,
			startTime: "",
			stopTimes: nil,
		}}
	}

	baseStart := trip.StopTimes[0].DepartureSeconds
	if len(trip.Frequencies) == 0 {
		return []tripInstance{{
			trip:         trip,
			startTime:    gtfs.FormatGTFSTime(baseStart),
			startSeconds: baseStart,
			stopTimes:    cloneShiftedStopTimes(trip.StopTimes, 0),
		}}
	}

	windowSeconds := int(window.Seconds())
	var instances []tripInstance
	for _, f := range trip.Frequencies {
		switch f.ExactTimes {
		case 1:
			for start := f.StartSeconds; start < f.EndSeconds; start += f.HeadwaySecs {
				shift := start - baseStart
				stops := cloneShiftedStopTimes(trip.StopTimes, shift)
				last := stops[len(stops)-1].DepartureSeconds
				if observedSeconds < start-60 || observedSeconds > last+windowSeconds {
					continue
				}
				instances = append(instances, tripInstance{
					trip:           trip,
					startTime:      gtfs.FormatGTFSTime(start),
					startSeconds:   start,
					stopTimes:      stops,
					exactFrequency: true,
				})
			}
		default:
			duration := trip.StopTimes[len(trip.StopTimes)-1].DepartureSeconds - baseStart
			if observedSeconds < f.StartSeconds-windowSeconds || observedSeconds > f.EndSeconds+duration+windowSeconds {
				continue
			}
			shift := f.StartSeconds - baseStart
			instances = append(instances, tripInstance{
				trip:              trip,
				startTime:         f.StartTime,
				startSeconds:      f.StartSeconds,
				stopTimes:         cloneShiftedStopTimes(trip.StopTimes, shift),
				nonExactFrequency: true,
			})
		}
	}
	return instances
}

func cloneShiftedStopTimes(stopTimes []gtfs.StopTime, shift int) []gtfs.StopTime {
	shifted := make([]gtfs.StopTime, len(stopTimes))
	for i, st := range stopTimes {
		shifted[i] = st
		shifted[i].ArrivalSeconds += shift
		shifted[i].DepartureSeconds += shift
	}
	return shifted
}

func (e *Engine) scoreCandidate(event telemetry.StoredEvent, feedVersionID string, day gtfs.ServiceDay, instance tripInstance, previous *Assignment) scoredCandidate {
	reasons := make([]string, 0, 8)
	score := 0.0
	scoreDetails := map[string]any{
		"score_schema":             "loose_debug_v1",
		"observed_local_seconds":   day.ObservedLocalSeconds,
		"trip_id":                  instance.trip.TripID,
		"start_time":               instance.startTime,
		"frequency_identity_type":  "scheduled_trip",
		"exact_scheduled_instance": true,
	}
	degraded := DegradedNone

	if event.TripHint != "" && event.TripHint == instance.trip.TripID {
		score += 0.35
		reasons = append(reasons, ReasonTripHintMatch)
	}

	currentStopSequence, scheduleScore, closestStopDiff := scheduleFit(instance.stopTimes, day.ObservedLocalSeconds, e.config.ScheduleFitWindow)
	scheduleStopSequence := currentStopSequence
	if scheduleScore > 0 {
		score += scheduleScore
		reasons = append(reasons, ReasonScheduleFitMatch)
		scoreDetails["closest_stop_seconds"] = closestStopDiff
	} else if len(instance.stopTimes) == 0 {
		reasons = append(reasons, ReasonMissingStopTimes)
		degraded = DegradedMissingScheduleData
	}

	shapeDist := 0.0
	proj := projectToShape(event.Lat, event.Lon, instance.trip.ShapePoints)
	if proj.Found {
		shapeDist = proj.ShapeDistance
		scoreDetails["shape_distance_meters"] = proj.DistanceMeters
		if proj.DistanceMeters <= e.config.ShapeDistanceMeters {
			shapeScore := 0.25 * (1 - proj.DistanceMeters/e.config.ShapeDistanceMeters)
			if shapeScore < 0.05 {
				shapeScore = 0.05
			}
			score += shapeScore
			reasons = append(reasons, ReasonShapeProximityMatch)
			if event.Bearing != 0 && bearingDelta(event.Bearing, proj.Bearing) <= 60 {
				score += 0.10
				reasons = append(reasons, ReasonMovementDirectionMatch)
			}
		} else {
			reasons = append(reasons, ReasonOffShape)
			degraded = DegradedLowConfidence
		}
		if seq := stopProgressFromShape(instance.stopTimes, shapeDist); seq > 0 {
			currentStopSequence = seq
			scoreDetails["shape_stop_sequence"] = seq
			scoreDetails["schedule_stop_sequence"] = scheduleStopSequence
			if scheduleStopSequence == 0 || scheduleStopSequence == seq {
				score += 0.12
				reasons = append(reasons, ReasonStopProgressMatch)
			} else {
				score -= 0.15
			}
		}
	} else {
		score -= 0.10
		reasons = append(reasons, ReasonMissingShape)
		degraded = DegradedMissingShape
	}

	if previous != nil {
		sameInstance := previous.TripID == instance.trip.TripID &&
			previous.StartDate == day.Date &&
			previous.StartTime == instance.startTime
		continuityAge := event.Timestamp.Sub(previous.ActiveFrom)
		temporalContinuity := !previous.ActiveFrom.IsZero() && continuityAge >= 0 && continuityAge <= e.config.ContinuityWindow
		temporalBlockTransition := !previous.ActiveFrom.IsZero() && continuityAge >= 0 && continuityAge <= e.config.BlockTransitionWindow
		scoreDetails["previous_assignment_age_seconds"] = continuityAge.Seconds()
		if sameInstance && temporalContinuity {
			score += 0.20
			reasons = append(reasons, ReasonContinuityMatch)
		} else if previous.BlockID != "" && previous.BlockID == instance.trip.BlockID && !sameInstance && temporalBlockTransition {
			score += 0.15
			reasons = append(reasons, ReasonBlockTransitionMatch)
		}
	}

	if instance.exactFrequency {
		reasons = append(reasons, ReasonFrequencyExactInstance)
		scoreDetails["frequency_identity_type"] = "exact_generated_instance"
		scoreDetails["exact_scheduled_instance"] = true
	}
	if instance.nonExactFrequency {
		reasons = append(reasons, ReasonFrequencyNonExact)
		score += 0.05
		scoreDetails["frequency_identity_type"] = "non_exact_window"
		scoreDetails["exact_scheduled_instance"] = false
		scoreDetails["non_exact_frequency_window_start_time"] = instance.startTime
	}

	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	assignment := Assignment{
		AgencyID:            event.AgencyID,
		VehicleID:           event.VehicleID,
		FeedVersionID:       feedVersionID,
		TelemetryEventID:    event.ID,
		State:               StateInService,
		ServiceDate:         day.Date,
		RouteID:             instance.trip.RouteID,
		TripID:              instance.trip.TripID,
		BlockID:             instance.trip.BlockID,
		StartDate:           day.Date,
		StartTime:           instance.startTime,
		CurrentStopSequence: currentStopSequence,
		ShapeDistTraveled:   shapeDist,
		Confidence:          math.Round(score*1000) / 1000,
		AssignmentSource:    AssignmentSourceAutomatic,
		ReasonCodes:         dedupeReasons(reasons),
		DegradedState:       degraded,
		ScoreDetails:        scoreDetails,
		ActiveFrom:          event.Timestamp,
	}
	return scoredCandidate{assignment: assignment}
}

func scheduleFit(stopTimes []gtfs.StopTime, observedSeconds int, window time.Duration) (int, float64, int) {
	if len(stopTimes) == 0 {
		return 0, 0, 0
	}
	bestSeq := stopTimes[0].StopSequence
	bestDiff := math.MaxInt
	for _, st := range stopTimes {
		diff := absInt(observedSeconds - st.DepartureSeconds)
		if diff < bestDiff {
			bestDiff = diff
			bestSeq = st.StopSequence
		}
	}
	windowSeconds := int(window.Seconds())
	if bestDiff > windowSeconds {
		return bestSeq, 0, bestDiff
	}
	score := 0.25 * (1 - float64(bestDiff)/float64(windowSeconds))
	if score < 0.05 {
		score = 0.05
	}
	return bestSeq, score, bestDiff
}

func stopProgressFromShape(stopTimes []gtfs.StopTime, shapeDist float64) int {
	if len(stopTimes) == 0 {
		return 0
	}
	bestSeq := 0
	bestDist := math.MaxFloat64
	for _, st := range stopTimes {
		diff := math.Abs(st.ShapeDistTraveled - shapeDist)
		if diff < bestDist {
			bestDist = diff
			bestSeq = st.StopSequence
		}
	}
	return bestSeq
}

func (e *Engine) saveManualOverride(ctx context.Context, event telemetry.StoredEvent, override ManualOverride, serviceDate string) (Assignment, error) {
	if override.StartDate != "" {
		serviceDate = override.StartDate
	}
	assignment := Assignment{
		AgencyID:         event.AgencyID,
		VehicleID:        event.VehicleID,
		TelemetryEventID: event.ID,
		State:            override.State,
		ServiceDate:      serviceDate,
		RouteID:          override.RouteID,
		TripID:           override.TripID,
		StartDate:        override.StartDate,
		StartTime:        override.StartTime,
		Confidence:       1,
		AssignmentSource: AssignmentSourceManualOverride,
		ReasonCodes:      []string{ReasonManualOverrideActive},
		DegradedState:    DegradedNone,
		ScoreDetails:     map[string]any{"score_schema": "loose_debug_v1", "override_type": override.Type, "reason": override.Reason},
		ManualOverrideID: override.ID,
		ActiveFrom:       event.Timestamp,
	}
	if assignment.State == "" {
		assignment.State = StateInService
	}
	return e.saveAssignment(ctx, assignment, nil)
}

func (e *Engine) saveUnknown(ctx context.Context, event telemetry.StoredEvent, serviceDate string, degraded DegradedState, reasons []string, incidentType string, details map[string]any) (Assignment, error) {
	if details == nil {
		details = map[string]any{}
	}
	details["score_schema"] = "loose_debug_v1"
	assignment := Assignment{
		AgencyID:         event.AgencyID,
		VehicleID:        event.VehicleID,
		TelemetryEventID: event.ID,
		State:            StateUnknown,
		ServiceDate:      serviceDate,
		Confidence:       0,
		AssignmentSource: AssignmentSourceAutomatic,
		ReasonCodes:      dedupeReasons(reasons),
		DegradedState:    degraded,
		ScoreDetails:     details,
		ActiveFrom:       event.Timestamp,
	}
	if assignment.DegradedState == "" {
		assignment.DegradedState = DegradedUnknown
	}
	incident := Incident{
		Type:      incidentType,
		Severity:  "warning",
		VehicleID: event.VehicleID,
		Details:   details,
	}
	return e.saveAssignment(ctx, assignment, []Incident{incident})
}

func (e *Engine) saveAssignment(ctx context.Context, assignment Assignment, incidents []Incident) (Assignment, error) {
	if e.assignments == nil {
		return assignment, nil
	}
	return e.assignments.SaveAssignment(ctx, assignment, incidents)
}

func dedupeReasons(reasons []string) []string {
	seen := make(map[string]bool, len(reasons))
	var deduped []string
	for _, reason := range reasons {
		if reason == "" || seen[reason] {
			continue
		}
		seen[reason] = true
		deduped = append(deduped, reason)
	}
	return deduped
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

var ErrNoAssignmentRepository = errors.New("assignment repository is required")

func (e *Engine) Validate() error {
	if e.schedules == nil {
		return fmt.Errorf("schedule repository is required")
	}
	if e.assignments == nil {
		return ErrNoAssignmentRepository
	}
	return nil
}
