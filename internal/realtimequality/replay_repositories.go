package realtimequality

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/prediction"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

type replaySchedule struct {
	agencyID      string
	timezone      string
	feedVersionID string
	tripsByDate   map[string][]gtfs.TripCandidate
}

func newReplaySchedule(scenario Scenario) (*replaySchedule, error) {
	schedules := &replaySchedule{
		agencyID:      scenario.AgencyID,
		timezone:      scenario.Timezone,
		feedVersionID: scenario.FeedVersionID,
		tripsByDate:   map[string][]gtfs.TripCandidate{},
	}
	for _, fixture := range scenario.Trips {
		trip, err := tripCandidateFromFixture(scenario, fixture)
		if err != nil {
			return nil, err
		}
		schedules.tripsByDate[fixture.ServiceDate] = append(schedules.tripsByDate[fixture.ServiceDate], trip)
	}
	for date := range schedules.tripsByDate {
		sort.SliceStable(schedules.tripsByDate[date], func(i int, j int) bool {
			left := schedules.tripsByDate[date][i]
			right := schedules.tripsByDate[date][j]
			if left.TripID != right.TripID {
				return left.TripID < right.TripID
			}
			return left.RouteID < right.RouteID
		})
	}
	return schedules, nil
}

func (r *replaySchedule) Agency(context.Context, string) (gtfs.Agency, error) {
	return gtfs.Agency{ID: r.agencyID, Timezone: r.timezone}, nil
}

func (r *replaySchedule) ActiveFeedVersion(context.Context, string) (gtfs.FeedVersion, error) {
	return gtfs.FeedVersion{ID: r.feedVersionID, AgencyID: r.agencyID}, nil
}

func (r *replaySchedule) ListTripCandidates(_ context.Context, _ string, _ string, serviceDate string) ([]gtfs.TripCandidate, error) {
	return append([]gtfs.TripCandidate(nil), r.tripsByDate[serviceDate]...), nil
}

func tripCandidateFromFixture(scenario Scenario, fixture TripFixture) (gtfs.TripCandidate, error) {
	stopTimes := make([]gtfs.StopTime, 0, len(fixture.StopTimes))
	for _, stop := range fixture.StopTimes {
		arrival, err := gtfs.ParseGTFSTime(stop.ArrivalTime)
		if err != nil {
			return gtfs.TripCandidate{}, fmt.Errorf("parse arrival for %s: %w", fixture.TripID, err)
		}
		departure, err := gtfs.ParseGTFSTime(stop.DepartureTime)
		if err != nil {
			return gtfs.TripCandidate{}, fmt.Errorf("parse departure for %s: %w", fixture.TripID, err)
		}
		stopTimes = append(stopTimes, gtfs.StopTime{
			TripID:            fixture.TripID,
			StopID:            stop.StopID,
			ArrivalSeconds:    arrival,
			DepartureSeconds:  departure,
			StopSequence:      stop.StopSequence,
			ShapeDistTraveled: stop.ShapeDistTraveled,
		})
	}
	shapePoints := make([]gtfs.ShapePoint, 0, len(fixture.ShapePoints))
	for _, point := range fixture.ShapePoints {
		shapePoints = append(shapePoints, gtfs.ShapePoint{
			ShapeID:      "shape-" + fixture.TripID,
			Lat:          point.Lat,
			Lon:          point.Lon,
			Sequence:     point.Sequence,
			DistTraveled: point.DistTraveled,
			HasDistance:  true,
		})
	}
	frequencies := make([]gtfs.Frequency, 0, len(fixture.Frequencies))
	for _, frequency := range fixture.Frequencies {
		start, err := gtfs.ParseGTFSTime(frequency.StartTime)
		if err != nil {
			return gtfs.TripCandidate{}, fmt.Errorf("parse frequency start for %s: %w", fixture.TripID, err)
		}
		end, err := gtfs.ParseGTFSTime(frequency.EndTime)
		if err != nil {
			return gtfs.TripCandidate{}, fmt.Errorf("parse frequency end for %s: %w", fixture.TripID, err)
		}
		frequencies = append(frequencies, gtfs.Frequency{
			TripID:       fixture.TripID,
			StartSeconds: start,
			EndSeconds:   end,
			HeadwaySecs:  frequency.HeadwaySecs,
			ExactTimes:   frequency.ExactTimes,
			StartTime:    frequency.StartTime,
			EndTime:      frequency.EndTime,
		})
	}
	return gtfs.TripCandidate{
		AgencyID:      scenario.AgencyID,
		FeedVersionID: scenario.FeedVersionID,
		ServiceDate:   fixture.ServiceDate,
		RouteID:       fixture.RouteID,
		ServiceID:     "replay-service",
		TripID:        fixture.TripID,
		BlockID:       fixture.BlockID,
		ShapeID:       "shape-" + fixture.TripID,
		StopTimes:     stopTimes,
		ShapePoints:   shapePoints,
		Frequencies:   frequencies,
	}, nil
}

type replayStateRepository struct {
	agencyID        string
	assignments     map[string]state.Assignment
	manualOverrides map[string]state.ManualOverride
}

func newReplayStateRepository(scenario Scenario) *replayStateRepository {
	repo := &replayStateRepository{
		agencyID:        scenario.AgencyID,
		assignments:     map[string]state.Assignment{},
		manualOverrides: map[string]state.ManualOverride{},
	}
	for _, fixture := range scenario.ManualOverrides {
		repo.manualOverrides[fixture.VehicleID] = state.ManualOverride{
			ID:        fixture.ID,
			AgencyID:  scenario.AgencyID,
			VehicleID: fixture.VehicleID,
			Type:      fixture.Type,
			RouteID:   fixture.RouteID,
			TripID:    fixture.TripID,
			StartDate: fixture.StartDate,
			StartTime: fixture.StartTime,
			State:     state.VehicleServiceState(fixture.State),
			Reason:    fixture.Reason,
			ExpiresAt: fixture.ExpiresAt,
			CreatedAt: scenario.GeneratedAt,
		}
	}
	return repo
}

func (r *replayStateRepository) ActiveManualOverride(_ context.Context, _ string, vehicleID string, at time.Time) (*state.ManualOverride, error) {
	override, ok := r.manualOverrides[vehicleID]
	if !ok {
		return nil, nil
	}
	if override.ExpiresAt != nil && !at.Before(*override.ExpiresAt) {
		return nil, nil
	}
	return &override, nil
}

func (r *replayStateRepository) CurrentAssignment(_ context.Context, _ string, vehicleID string) (*state.Assignment, error) {
	assignment, ok := r.assignments[vehicleID]
	if !ok {
		return nil, nil
	}
	copied := assignment
	return &copied, nil
}

func (r *replayStateRepository) ListCurrentAssignments(_ context.Context, _ string, vehicleIDs []string) (map[string]state.Assignment, error) {
	out := map[string]state.Assignment{}
	if len(vehicleIDs) == 0 {
		for vehicleID, assignment := range r.assignments {
			out[vehicleID] = assignment
		}
		return out, nil
	}
	for _, vehicleID := range vehicleIDs {
		if assignment, ok := r.assignments[vehicleID]; ok {
			out[vehicleID] = assignment
		}
	}
	return out, nil
}

func (r *replayStateRepository) SaveAssignment(_ context.Context, assignment state.Assignment, _ []state.Incident) (state.Assignment, error) {
	r.assignments[assignment.VehicleID] = assignment
	return assignment, nil
}

func (r *replayStateRepository) currentAssignments() map[string]state.Assignment {
	out := make(map[string]state.Assignment, len(r.assignments))
	for vehicleID, assignment := range r.assignments {
		out[vehicleID] = assignment
	}
	return out
}

type replayTelemetryRepository struct {
	events []telemetry.StoredEvent
}

func (r replayTelemetryRepository) Store(context.Context, telemetry.Event, json.RawMessage) (telemetry.StoreResult, error) {
	return telemetry.StoreResult{}, fmt.Errorf("replay telemetry repository is read-only")
}

func (r replayTelemetryRepository) LatestByVehicle(_ context.Context, _ string, vehicleID string) (telemetry.StoredEvent, error) {
	var latest telemetry.StoredEvent
	found := false
	for _, event := range r.events {
		if event.VehicleID != vehicleID {
			continue
		}
		if !found || isNewerTelemetry(event, latest) {
			latest = event
			found = true
		}
	}
	if found {
		return latest, nil
	}
	return telemetry.StoredEvent{}, fmt.Errorf("vehicle %s not found", vehicleID)
}

func (r replayTelemetryRepository) ListLatestByAgency(_ context.Context, _ string, limit int) ([]telemetry.StoredEvent, error) {
	byVehicle := map[string]telemetry.StoredEvent{}
	for _, event := range r.events {
		current, ok := byVehicle[event.VehicleID]
		if !ok || isNewerTelemetry(event, current) {
			byVehicle[event.VehicleID] = event
		}
	}
	latest := make([]telemetry.StoredEvent, 0, len(byVehicle))
	for _, event := range byVehicle {
		latest = append(latest, event)
	}
	sort.SliceStable(latest, func(i int, j int) bool {
		if latest[i].Timestamp.Equal(latest[j].Timestamp) {
			return latest[i].ID > latest[j].ID
		}
		return latest[i].Timestamp.After(latest[j].Timestamp)
	})
	if limit > 0 && len(latest) > limit {
		latest = latest[:limit]
	}
	return latest, nil
}

func isNewerTelemetry(candidate telemetry.StoredEvent, current telemetry.StoredEvent) bool {
	if candidate.Timestamp.Equal(current.Timestamp) {
		return candidate.ID > current.ID
	}
	return candidate.Timestamp.After(current.Timestamp)
}

func (r replayTelemetryRepository) ListEvents(_ context.Context, _ string, limit int) ([]telemetry.StoredEvent, error) {
	events := append([]telemetry.StoredEvent(nil), r.events...)
	if limit > 0 && len(events) > limit {
		events = events[:limit]
	}
	return events, nil
}

type replayPredictionOperations struct {
	overrides []prediction.OverrideRecord
	reviews   []prediction.ReviewItem
}

func newReplayPredictionOperations(scenario Scenario) *replayPredictionOperations {
	ops := &replayPredictionOperations{}
	for _, fixture := range scenario.PredictionOverrides {
		ops.overrides = append(ops.overrides, prediction.OverrideRecord{
			ID:           fixture.ID,
			AgencyID:     scenario.AgencyID,
			VehicleID:    fixture.VehicleID,
			OverrideType: fixture.Type,
			RouteID:      fixture.RouteID,
			TripID:       fixture.TripID,
			StartDate:    fixture.StartDate,
			StartTime:    fixture.StartTime,
			State:        fixture.State,
			Reason:       fixture.Reason,
			CreatedAt:    scenario.GeneratedAt,
		})
	}
	return ops
}

func (r *replayPredictionOperations) ListActivePredictionOverrides(context.Context, string, time.Time) ([]prediction.OverrideRecord, error) {
	return append([]prediction.OverrideRecord(nil), r.overrides...), nil
}

func (r *replayPredictionOperations) CreatePredictionOverride(context.Context, prediction.OverrideInput) (prediction.OverrideRecord, error) {
	return prediction.OverrideRecord{}, fmt.Errorf("replay prediction operations repository is read-only")
}

func (r *replayPredictionOperations) ReplacePredictionOverride(context.Context, prediction.OverrideInput) (prediction.OverrideRecord, error) {
	return prediction.OverrideRecord{}, fmt.Errorf("replay prediction operations repository is read-only")
}

func (r *replayPredictionOperations) ClearPredictionOverride(context.Context, string, int64, string, string, time.Time) error {
	return fmt.Errorf("replay prediction operations repository is read-only")
}

func (r *replayPredictionOperations) SavePredictionReviewItems(_ context.Context, items []prediction.ReviewItem) error {
	r.reviews = append(r.reviews, items...)
	return nil
}

func (r *replayPredictionOperations) ListPredictionReviewItems(context.Context, prediction.ReviewFilter) ([]prediction.ReviewItem, error) {
	return append([]prediction.ReviewItem(nil), r.reviews...), nil
}

func (r *replayPredictionOperations) UpdatePredictionReviewStatus(context.Context, string, int64, prediction.ReviewStatus, string, string, time.Time) error {
	return fmt.Errorf("replay prediction operations repository is read-only")
}
