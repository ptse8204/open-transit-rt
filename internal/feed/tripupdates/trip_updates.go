package tripupdates

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"

	"open-transit-rt/internal/feed"
	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/prediction"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

const DiagnosticsPersistenceStored = "stored"
const DiagnosticsPersistenceFailed = "failed"

type Config struct {
	AgencyID            string
	MaxVehicles         int
	VehiclePositionsURL string
}

func (c Config) Validated() (Config, error) {
	if c.AgencyID == "" {
		return Config{}, fmt.Errorf("AGENCY_ID is required")
	}
	if c.MaxVehicles < 1 {
		return Config{}, fmt.Errorf("TRIP_UPDATES_MAX_VEHICLES must be at least 1")
	}
	if c.VehiclePositionsURL == "" {
		return Config{}, fmt.Errorf("vehicle positions feed URL is required")
	}
	return c, nil
}

type Builder struct {
	schedules   gtfs.Repository
	telemetry   telemetry.Repository
	assignments state.Repository
	adapter     prediction.Adapter
	diagnostics prediction.DiagnosticsRepository
	config      Config
}

func NewBuilder(
	scheduleRepo gtfs.Repository,
	telemetryRepo telemetry.Repository,
	assignmentRepo state.Repository,
	adapter prediction.Adapter,
	diagnosticsRepo prediction.DiagnosticsRepository,
	config Config,
) (*Builder, error) {
	if scheduleRepo == nil {
		return nil, fmt.Errorf("gtfs repository is required")
	}
	if telemetryRepo == nil {
		return nil, fmt.Errorf("telemetry repository is required")
	}
	if assignmentRepo == nil {
		return nil, fmt.Errorf("assignment repository is required")
	}
	if adapter == nil {
		return nil, fmt.Errorf("prediction adapter is required")
	}
	if diagnosticsRepo == nil {
		return nil, fmt.Errorf("diagnostics repository is required")
	}
	validated, err := config.Validated()
	if err != nil {
		return nil, err
	}
	return &Builder{
		schedules:   scheduleRepo,
		telemetry:   telemetryRepo,
		assignments: assignmentRepo,
		adapter:     adapter,
		diagnostics: diagnosticsRepo,
		config:      validated,
	}, nil
}

func (b *Builder) Snapshot(ctx context.Context, generatedAt time.Time) (Snapshot, error) {
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}
	generatedAt = generatedAt.UTC()

	activeFeed, err := b.schedules.ActiveFeedVersion(ctx, b.config.AgencyID)
	if err != nil {
		if gtfs.IsNoRows(err) {
			snapshot := Snapshot{
				AgencyID:            b.config.AgencyID,
				GeneratedAt:         generatedAt,
				AdapterName:         b.adapter.Name(),
				VehiclePositionsURL: b.config.VehiclePositionsURL,
				Diagnostics: prediction.Diagnostics{
					Status: prediction.StatusError,
					Reason: prediction.ReasonActiveFeedUnavailable,
				},
			}
			b.persistDiagnostics(ctx, &snapshot)
			return snapshot, nil
		}
		return Snapshot{}, fmt.Errorf("active feed version: %w", err)
	}

	latest, err := b.telemetry.ListLatestByAgency(ctx, b.config.AgencyID, b.config.MaxVehicles+1)
	if err != nil {
		return Snapshot{}, fmt.Errorf("list latest telemetry: %w", err)
	}

	snapshot := Snapshot{
		AgencyID:            b.config.AgencyID,
		GeneratedAt:         generatedAt,
		ActiveFeedVersionID: activeFeed.ID,
		AdapterName:         b.adapter.Name(),
		VehicleLimit:        b.config.MaxVehicles,
		VehiclePositionsURL: b.config.VehiclePositionsURL,
	}
	if len(latest) > b.config.MaxVehicles {
		snapshot.Truncated = true
		snapshot.TruncatedVehicleCountMin = len(latest) - b.config.MaxVehicles
		latest = latest[:b.config.MaxVehicles]
	}
	snapshot.LatestTelemetryRows = len(latest)

	vehicleIDs := make([]string, 0, len(latest))
	for _, event := range latest {
		vehicleIDs = append(vehicleIDs, event.VehicleID)
	}
	assignments, err := b.assignments.ListCurrentAssignments(ctx, b.config.AgencyID, vehicleIDs)
	if err != nil {
		return Snapshot{}, fmt.Errorf("list current assignments: %w", err)
	}
	snapshot.AssignmentRows = len(assignments)

	result, err := b.adapter.PredictTripUpdates(ctx, prediction.Request{
		AgencyID:            b.config.AgencyID,
		GeneratedAt:         generatedAt,
		ActiveFeedVersion:   activeFeed,
		Telemetry:           append([]telemetry.StoredEvent(nil), latest...),
		Assignments:         copyAssignments(assignments),
		VehiclePositionsURL: b.config.VehiclePositionsURL,
	})
	if err != nil {
		snapshot.Diagnostics = prediction.Diagnostics{
			Status: prediction.StatusError,
			Reason: prediction.ReasonAdapterError,
		}
	} else {
		snapshot.TripUpdates = normalizeTripUpdates(result.TripUpdates)
		snapshot.Diagnostics = result.Diagnostics
		if snapshot.Diagnostics.Status == "" {
			snapshot.Diagnostics.Status = prediction.StatusOK
		}
	}

	b.persistDiagnostics(ctx, &snapshot)
	return snapshot, nil
}

func (b *Builder) Ready(ctx context.Context) error {
	if _, err := b.schedules.ActiveFeedVersion(ctx, b.config.AgencyID); err != nil {
		return fmt.Errorf("active feed version unavailable: %w", err)
	}
	return nil
}

func (b *Builder) persistDiagnostics(ctx context.Context, snapshot *Snapshot) {
	record := prediction.DiagnosticsRecord{
		AgencyID:            snapshot.AgencyID,
		SnapshotAt:          snapshot.GeneratedAt,
		AdapterName:         snapshot.AdapterName,
		Status:              snapshot.Diagnostics.Status,
		Reason:              snapshot.Diagnostics.Reason,
		ActiveFeedVersionID: snapshot.ActiveFeedVersionID,
		InputCounts: prediction.InputCounts{
			TelemetryRows:     snapshot.LatestTelemetryRows,
			AssignmentRows:    snapshot.AssignmentRows,
			TripUpdatesOutput: len(snapshot.TripUpdates),
		},
		Metrics:                     snapshot.Diagnostics.Metrics,
		AdapterDetails:              snapshot.Diagnostics.Details,
		VehiclePositionsURL:         snapshot.VehiclePositionsURL,
		DiagnosticsPersistenceState: DiagnosticsPersistenceStored,
	}
	if _, err := b.diagnostics.SaveTripUpdatesDiagnostics(ctx, record); err != nil {
		snapshot.DiagnosticsPersistenceOutcome = DiagnosticsPersistenceFailed
		snapshot.DiagnosticsPersistenceError = err.Error()
		return
	}
	snapshot.DiagnosticsPersistenceOutcome = DiagnosticsPersistenceStored
}

type Snapshot struct {
	AgencyID                      string
	GeneratedAt                   time.Time
	ActiveFeedVersionID           string
	AdapterName                   string
	VehicleLimit                  int
	LatestTelemetryRows           int
	AssignmentRows                int
	Truncated                     bool
	TruncatedVehicleCountMin      int
	VehiclePositionsURL           string
	Diagnostics                   prediction.Diagnostics
	DiagnosticsPersistenceOutcome string
	DiagnosticsPersistenceError   string
	TripUpdates                   []prediction.TripUpdate
}

func (s Snapshot) BuildProto() (*gtfsrt.FeedMessage, error) {
	timestamp := uint64(s.GeneratedAt.Unix())
	incrementality := gtfsrt.FeedHeader_FULL_DATASET
	message := &gtfsrt.FeedMessage{
		Header: &gtfsrt.FeedHeader{
			GtfsRealtimeVersion: proto.String(feed.GTFSRealtimeVersion),
			Incrementality:      &incrementality,
			Timestamp:           &timestamp,
		},
		Entity: []*gtfsrt.FeedEntity{},
	}

	for _, update := range normalizeTripUpdates(s.TripUpdates) {
		entity, err := buildEntity(update)
		if err != nil {
			return nil, err
		}
		message.Entity = append(message.Entity, entity)
	}
	return message, nil
}

func (s Snapshot) MarshalProto() ([]byte, error) {
	message, err := s.BuildProto()
	if err != nil {
		return nil, err
	}
	payload, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("marshal trip updates protobuf: %w", err)
	}
	return payload, nil
}

func (s Snapshot) MarshalDebugJSON() ([]byte, error) {
	payload, err := json.MarshalIndent(Debug{
		AgencyID:                      s.AgencyID,
		GeneratedAt:                   s.GeneratedAt,
		ActiveFeedVersionID:           s.ActiveFeedVersionID,
		AdapterName:                   s.AdapterName,
		VehicleLimit:                  s.VehicleLimit,
		LatestTelemetryRows:           s.LatestTelemetryRows,
		AssignmentRows:                s.AssignmentRows,
		TripUpdatesOutput:             len(s.TripUpdates),
		Truncated:                     s.Truncated,
		TruncatedVehicleCountMin:      s.TruncatedVehicleCountMin,
		VehiclePositionsURL:           s.VehiclePositionsURL,
		DiagnosticsStatus:             s.Diagnostics.Status,
		DiagnosticsReason:             s.Diagnostics.Reason,
		PredictionMetrics:             s.Diagnostics.Metrics,
		DiagnosticsDetails:            s.Diagnostics.Details,
		DiagnosticsPersistenceOutcome: s.DiagnosticsPersistenceOutcome,
		DiagnosticsPersistenceError:   s.DiagnosticsPersistenceError,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal trip updates debug json: %w", err)
	}
	return payload, nil
}

type Debug struct {
	AgencyID                      string             `json:"agency_id"`
	GeneratedAt                   time.Time          `json:"generated_at"`
	ActiveFeedVersionID           string             `json:"active_feed_version_id"`
	AdapterName                   string             `json:"adapter_name"`
	VehicleLimit                  int                `json:"vehicle_limit"`
	LatestTelemetryRows           int                `json:"latest_telemetry_rows"`
	AssignmentRows                int                `json:"assignment_rows"`
	TripUpdatesOutput             int                `json:"trip_updates_output"`
	Truncated                     bool               `json:"truncated"`
	TruncatedVehicleCountMin      int                `json:"truncated_vehicle_count_min"`
	VehiclePositionsURL           string             `json:"vehicle_positions_url"`
	DiagnosticsStatus             string             `json:"diagnostics_status"`
	DiagnosticsReason             string             `json:"diagnostics_reason"`
	PredictionMetrics             prediction.Metrics `json:"prediction_metrics"`
	DiagnosticsDetails            map[string]any     `json:"diagnostics_details,omitempty"`
	DiagnosticsPersistenceOutcome string             `json:"diagnostics_persistence_outcome"`
	DiagnosticsPersistenceError   string             `json:"diagnostics_persistence_error,omitempty"`
}

func buildEntity(update prediction.TripUpdate) (*gtfsrt.FeedEntity, error) {
	entityID := update.EntityID
	if entityID == "" {
		entityID = update.TripID
	}
	if entityID == "" {
		return nil, fmt.Errorf("trip update entity id or trip id is required")
	}
	tripUpdate := &gtfsrt.TripUpdate{
		Trip:           buildTripDescriptor(update),
		StopTimeUpdate: buildStopTimeUpdates(update.StopTimeUpdates),
	}
	if update.VehicleID != "" {
		tripUpdate.Vehicle = &gtfsrt.VehicleDescriptor{Id: proto.String(update.VehicleID)}
	}
	return &gtfsrt.FeedEntity{Id: proto.String(entityID), TripUpdate: tripUpdate}, nil
}

func buildTripDescriptor(update prediction.TripUpdate) *gtfsrt.TripDescriptor {
	descriptor := &gtfsrt.TripDescriptor{}
	if update.TripID != "" {
		descriptor.TripId = proto.String(update.TripID)
	}
	if update.RouteID != "" {
		descriptor.RouteId = proto.String(update.RouteID)
	}
	if update.StartDate != "" {
		descriptor.StartDate = proto.String(update.StartDate)
	}
	if update.StartTime != "" {
		descriptor.StartTime = proto.String(update.StartTime)
	}
	relationship := tripScheduleRelationship(update.ScheduleRelationship)
	descriptor.ScheduleRelationship = &relationship
	return descriptor
}

func buildStopTimeUpdates(updates []prediction.StopTimeUpdate) []*gtfsrt.TripUpdate_StopTimeUpdate {
	normalized := append([]prediction.StopTimeUpdate(nil), updates...)
	sort.SliceStable(normalized, func(i int, j int) bool {
		return normalized[i].StopSequence < normalized[j].StopSequence
	})
	result := make([]*gtfsrt.TripUpdate_StopTimeUpdate, 0, len(normalized))
	for _, update := range normalized {
		stopUpdate := &gtfsrt.TripUpdate_StopTimeUpdate{}
		if update.StopID != "" {
			stopUpdate.StopId = proto.String(update.StopID)
		}
		if update.StopSequence > 0 {
			stopUpdate.StopSequence = proto.Uint32(uint32(update.StopSequence))
		}
		if update.ArrivalTime != nil || update.ArrivalDelaySeconds != nil {
			stopUpdate.Arrival = stopTimeEvent(update.ArrivalTime, update.ArrivalDelaySeconds)
		}
		if update.DepartureTime != nil || update.DepartureDelaySeconds != nil {
			stopUpdate.Departure = stopTimeEvent(update.DepartureTime, update.DepartureDelaySeconds)
		}
		relationship := stopScheduleRelationship(update.ScheduleRelationship)
		stopUpdate.ScheduleRelationship = &relationship
		result = append(result, stopUpdate)
	}
	return result
}

func stopTimeEvent(t *time.Time, delay *int32) *gtfsrt.TripUpdate_StopTimeEvent {
	event := &gtfsrt.TripUpdate_StopTimeEvent{}
	if t != nil {
		event.Time = proto.Int64(t.Unix())
	}
	if delay != nil {
		event.Delay = delay
	}
	return event
}

func tripScheduleRelationship(relationship prediction.ScheduleRelationship) gtfsrt.TripDescriptor_ScheduleRelationship {
	switch relationship {
	case prediction.ScheduleRelationshipUnscheduled:
		return gtfsrt.TripDescriptor_UNSCHEDULED
	case prediction.ScheduleRelationshipCanceled:
		return gtfsrt.TripDescriptor_CANCELED
	case prediction.ScheduleRelationshipAdded:
		return gtfsrt.TripDescriptor_ADDED
	default:
		return gtfsrt.TripDescriptor_SCHEDULED
	}
}

func stopScheduleRelationship(relationship prediction.ScheduleRelationship) gtfsrt.TripUpdate_StopTimeUpdate_ScheduleRelationship {
	switch relationship {
	case prediction.ScheduleRelationshipUnscheduled:
		return gtfsrt.TripUpdate_StopTimeUpdate_UNSCHEDULED
	default:
		return gtfsrt.TripUpdate_StopTimeUpdate_SCHEDULED
	}
}

func normalizeTripUpdates(updates []prediction.TripUpdate) []prediction.TripUpdate {
	normalized := append([]prediction.TripUpdate(nil), updates...)
	for i := range normalized {
		normalized[i].StopTimeUpdates = append([]prediction.StopTimeUpdate(nil), normalized[i].StopTimeUpdates...)
		sort.SliceStable(normalized[i].StopTimeUpdates, func(left int, right int) bool {
			return normalized[i].StopTimeUpdates[left].StopSequence < normalized[i].StopTimeUpdates[right].StopSequence
		})
	}
	sort.SliceStable(normalized, func(i int, j int) bool {
		return tripUpdateSortKey(normalized[i]) < tripUpdateSortKey(normalized[j])
	})
	return normalized
}

func tripUpdateSortKey(update prediction.TripUpdate) string {
	if update.EntityID != "" {
		return update.EntityID
	}
	return update.TripID
}

func copyAssignments(assignments map[string]state.Assignment) map[string]state.Assignment {
	copied := make(map[string]state.Assignment, len(assignments))
	for vehicleID, assignment := range assignments {
		copied[vehicleID] = assignment
	}
	return copied
}
