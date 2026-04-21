package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"

	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

const GTFSRealtimeVersion = "2.0"

const (
	TripDescriptorOmissionNone                        = "none"
	TripDescriptorOmissionVehicleNotIncluded          = "vehicle_not_included"
	TripDescriptorOmissionSuppressedStaleTelemetry    = "suppressed_stale_telemetry"
	TripDescriptorOmissionStaleTelemetry              = "stale_telemetry"
	TripDescriptorOmissionNoAssignment                = "no_assignment"
	TripDescriptorOmissionAssignmentTelemetryMismatch = "assignment_telemetry_mismatch"
	TripDescriptorOmissionNotInService                = "not_in_service"
	TripDescriptorOmissionMissingTripID               = "missing_trip_id"
	TripDescriptorOmissionManualStateWithoutTrip      = "manual_state_without_trip"
	TripDescriptorOmissionDegradedAssignment          = "degraded_assignment"
	TripDescriptorOmissionBelowPublicationConfidence  = "below_publication_confidence"
)

type VehiclePositionsConfig struct {
	AgencyID                  string
	MaxVehicles               int
	StaleTelemetryTTL         time.Duration
	SuppressStaleVehicleAfter time.Duration
	TripConfidenceThreshold   float64
}

func (c VehiclePositionsConfig) Validated() (VehiclePositionsConfig, error) {
	if c.AgencyID == "" {
		return VehiclePositionsConfig{}, fmt.Errorf("AGENCY_ID is required")
	}
	if c.MaxVehicles < 1 {
		return VehiclePositionsConfig{}, fmt.Errorf("VEHICLE_POSITIONS_MAX_VEHICLES must be at least 1")
	}
	if c.StaleTelemetryTTL <= 0 {
		return VehiclePositionsConfig{}, fmt.Errorf("STALE_TELEMETRY_TTL_SECONDS must be greater than 0")
	}
	if c.SuppressStaleVehicleAfter <= 0 {
		return VehiclePositionsConfig{}, fmt.Errorf("SUPPRESS_STALE_VEHICLE_AFTER_SECONDS must be greater than 0")
	}
	if c.SuppressStaleVehicleAfter < c.StaleTelemetryTTL {
		return VehiclePositionsConfig{}, fmt.Errorf("SUPPRESS_STALE_VEHICLE_AFTER_SECONDS must be greater than or equal to STALE_TELEMETRY_TTL_SECONDS")
	}
	if c.TripConfidenceThreshold < 0 || c.TripConfidenceThreshold > 1 {
		return VehiclePositionsConfig{}, fmt.Errorf("VEHICLE_POSITIONS_TRIP_CONFIDENCE_THRESHOLD must be between 0 and 1")
	}
	return c, nil
}

type VehiclePositionsBuilder struct {
	telemetry   telemetry.Repository
	assignments state.Repository
	config      VehiclePositionsConfig
}

func NewVehiclePositionsBuilder(telemetryRepo telemetry.Repository, assignmentRepo state.Repository, config VehiclePositionsConfig) (*VehiclePositionsBuilder, error) {
	if telemetryRepo == nil {
		return nil, fmt.Errorf("telemetry repository is required")
	}
	if assignmentRepo == nil {
		return nil, fmt.Errorf("assignment repository is required")
	}
	validated, err := config.Validated()
	if err != nil {
		return nil, err
	}
	return &VehiclePositionsBuilder{telemetry: telemetryRepo, assignments: assignmentRepo, config: validated}, nil
}

func (b *VehiclePositionsBuilder) Snapshot(ctx context.Context, generatedAt time.Time) (VehiclePositionsSnapshot, error) {
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}
	generatedAt = generatedAt.UTC()

	latest, err := b.telemetry.ListLatestByAgency(ctx, b.config.AgencyID, b.config.MaxVehicles+1)
	if err != nil {
		return VehiclePositionsSnapshot{}, fmt.Errorf("list latest telemetry: %w", err)
	}

	snapshot := VehiclePositionsSnapshot{
		AgencyID:                b.config.AgencyID,
		GeneratedAt:             generatedAt,
		VehicleLimit:            b.config.MaxVehicles,
		LatestTelemetryRowsRead: len(latest),
	}
	if len(latest) > b.config.MaxVehicles {
		snapshot.Truncated = true
		snapshot.TruncatedVehicleCountMin = len(latest) - b.config.MaxVehicles
		latest = latest[:b.config.MaxVehicles]
	}

	vehicleIDs := make([]string, 0, len(latest))
	for _, event := range latest {
		vehicleIDs = append(vehicleIDs, event.VehicleID)
	}
	assignments, err := b.assignments.ListCurrentAssignments(ctx, b.config.AgencyID, vehicleIDs)
	if err != nil {
		return VehiclePositionsSnapshot{}, fmt.Errorf("list current assignments: %w", err)
	}

	snapshot.Vehicles = make([]VehicleSnapshot, 0, len(latest))
	for _, event := range latest {
		assignment, hasAssignment := assignments[event.VehicleID]
		vehicle := b.vehicleSnapshot(event, assignment, hasAssignment, generatedAt)
		snapshot.Vehicles = append(snapshot.Vehicles, vehicle)
	}
	sort.Slice(snapshot.Vehicles, func(i int, j int) bool {
		return snapshot.Vehicles[i].VehicleID < snapshot.Vehicles[j].VehicleID
	})
	snapshot.VehiclesInSnapshot = len(snapshot.Vehicles)
	snapshot.NoTelemetry = len(snapshot.Vehicles) == 0
	return snapshot, nil
}

func (b *VehiclePositionsBuilder) vehicleSnapshot(event telemetry.StoredEvent, assignment state.Assignment, hasAssignment bool, generatedAt time.Time) VehicleSnapshot {
	age := generatedAt.Sub(event.Timestamp)
	if age < 0 {
		age = 0
	}
	vehicle := VehicleSnapshot{
		VehicleID:           event.VehicleID,
		TelemetryEvent:      event,
		TelemetryAgeSeconds: age.Seconds(),
		HasAssignment:       hasAssignment,
		IncludedInProtobuf:  true,
	}
	if hasAssignment {
		vehicle.Assignment = assignment
	}

	if hasAssignment && assignment.AssignmentSource == state.AssignmentSourceAutomatic && assignment.TelemetryEventID != 0 && assignment.TelemetryEventID != event.ID {
		vehicle.AssignmentTelemetryMismatch = true
	}

	if age > b.config.SuppressStaleVehicleAfter {
		vehicle.IncludedInProtobuf = false
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionSuppressedStaleTelemetry
		return vehicle
	}
	if age > b.config.StaleTelemetryTTL {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionStaleTelemetry
		return vehicle
	}
	if !hasAssignment {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionNoAssignment
		return vehicle
	}
	if vehicle.AssignmentTelemetryMismatch {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionAssignmentTelemetryMismatch
		return vehicle
	}
	if assignment.State != state.StateInService {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionNotInService
		return vehicle
	}
	if assignment.TripID == "" {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionMissingTripID
		return vehicle
	}
	if assignment.AssignmentSource == state.AssignmentSourceManualOverride && assignment.TripID == "" {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionManualStateWithoutTrip
		return vehicle
	}
	if degradedAssignmentBlocksTripDescriptor(assignment.DegradedState) {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionDegradedAssignment
		return vehicle
	}
	if assignment.AssignmentSource != state.AssignmentSourceManualOverride && assignment.Confidence < b.config.TripConfidenceThreshold {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionBelowPublicationConfidence
		return vehicle
	}

	vehicle.AssignmentPublishable = true
	vehicle.TripDescriptorPublished = true
	vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionNone
	return vehicle
}

type VehiclePositionsSnapshot struct {
	AgencyID                 string
	GeneratedAt              time.Time
	VehicleLimit             int
	LatestTelemetryRowsRead  int
	VehiclesInSnapshot       int
	Truncated                bool
	TruncatedVehicleCountMin int
	NoTelemetry              bool
	Vehicles                 []VehicleSnapshot
}

type VehicleSnapshot struct {
	VehicleID                    string
	TelemetryEvent               telemetry.StoredEvent
	TelemetryAgeSeconds          float64
	HasAssignment                bool
	Assignment                   state.Assignment
	IncludedInProtobuf           bool
	AssignmentPublishable        bool
	AssignmentTelemetryMismatch  bool
	TripDescriptorPublished      bool
	TripDescriptorOmissionReason string
}

func (s VehiclePositionsSnapshot) BuildProto() (*gtfsrt.FeedMessage, error) {
	timestamp := uint64(s.GeneratedAt.Unix())
	incrementality := gtfsrt.FeedHeader_FULL_DATASET
	message := &gtfsrt.FeedMessage{
		Header: &gtfsrt.FeedHeader{
			GtfsRealtimeVersion: proto.String(GTFSRealtimeVersion),
			Incrementality:      &incrementality,
			Timestamp:           &timestamp,
		},
		Entity: []*gtfsrt.FeedEntity{},
	}

	for _, vehicle := range s.Vehicles {
		if !vehicle.IncludedInProtobuf {
			continue
		}
		entity, err := vehicle.buildProtoEntity()
		if err != nil {
			return nil, err
		}
		message.Entity = append(message.Entity, entity)
	}
	return message, nil
}

func (s VehiclePositionsSnapshot) MarshalProto() ([]byte, error) {
	message, err := s.BuildProto()
	if err != nil {
		return nil, err
	}
	payload, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("marshal vehicle positions protobuf: %w", err)
	}
	return payload, nil
}

func (v VehicleSnapshot) buildProtoEntity() (*gtfsrt.FeedEntity, error) {
	id := v.VehicleID
	observedAt := uint64(v.TelemetryEvent.Timestamp.Unix())
	lat := float32(v.TelemetryEvent.Lat)
	lon := float32(v.TelemetryEvent.Lon)
	position := &gtfsrt.Position{
		Latitude:  &lat,
		Longitude: &lon,
	}
	if hasNumericPayloadField(v.TelemetryEvent.PayloadJSON, "bearing") {
		bearing := float32(v.TelemetryEvent.Bearing)
		position.Bearing = &bearing
	}
	if hasNumericPayloadField(v.TelemetryEvent.PayloadJSON, "speed_mps") {
		speed := float32(v.TelemetryEvent.SpeedMPS)
		position.Speed = &speed
	}

	vehicle := &gtfsrt.VehiclePosition{
		Vehicle: &gtfsrt.VehicleDescriptor{
			Id: proto.String(v.VehicleID),
		},
		Position:  position,
		Timestamp: &observedAt,
	}
	if v.TripDescriptorPublished {
		vehicle.Trip = buildTripDescriptor(v.Assignment)
	}

	return &gtfsrt.FeedEntity{
		Id:      &id,
		Vehicle: vehicle,
	}, nil
}

func buildTripDescriptor(assignment state.Assignment) *gtfsrt.TripDescriptor {
	descriptor := &gtfsrt.TripDescriptor{}
	if assignment.TripID != "" {
		descriptor.TripId = proto.String(assignment.TripID)
	}
	if assignment.RouteID != "" {
		descriptor.RouteId = proto.String(assignment.RouteID)
	}
	if assignment.StartDate != "" {
		descriptor.StartDate = proto.String(assignment.StartDate)
	}
	if assignment.StartTime != "" {
		descriptor.StartTime = proto.String(assignment.StartTime)
	}
	relationship := gtfsrt.TripDescriptor_SCHEDULED
	if hasReason(assignment, state.ReasonFrequencyNonExact) {
		relationship = gtfsrt.TripDescriptor_UNSCHEDULED
	}
	descriptor.ScheduleRelationship = &relationship
	return descriptor
}

type VehiclePositionsDebug struct {
	AgencyID                 string                 `json:"agency_id"`
	GeneratedAt              time.Time              `json:"generated_at"`
	Truncated                bool                   `json:"truncated"`
	VehicleLimit             int                    `json:"vehicle_limit"`
	LatestTelemetryRowsRead  int                    `json:"latest_telemetry_rows_read"`
	VehiclesInSnapshot       int                    `json:"vehicles_in_snapshot"`
	TruncatedVehicleCountMin int                    `json:"truncated_vehicle_count_min"`
	NoTelemetry              bool                   `json:"no_telemetry"`
	Vehicles                 []VehicleDebugSnapshot `json:"vehicles"`
}

type VehicleDebugSnapshot struct {
	VehicleID                    string            `json:"vehicle_id"`
	TelemetryEventID             int64             `json:"telemetry_event_id"`
	ObservedAt                   time.Time         `json:"observed_at"`
	TelemetryAgeSeconds          float64           `json:"telemetry_age_seconds"`
	Position                     DebugPosition     `json:"position"`
	IncludedInProtobuf           bool              `json:"included_in_protobuf"`
	HasAssignment                bool              `json:"has_assignment"`
	AssignmentPublishable        bool              `json:"assignment_publishable"`
	AssignmentTelemetryMismatch  bool              `json:"assignment_telemetry_mismatch"`
	TripDescriptorPublished      bool              `json:"trip_descriptor_published"`
	TripDescriptorOmissionReason string            `json:"trip_descriptor_omission_reason"`
	Assignment                   *state.Assignment `json:"assignment,omitempty"`
}

type DebugPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Bearing   float64 `json:"bearing,omitempty"`
	SpeedMPS  float64 `json:"speed_mps,omitempty"`
}

func (s VehiclePositionsSnapshot) Debug() VehiclePositionsDebug {
	debug := VehiclePositionsDebug{
		AgencyID:                 s.AgencyID,
		GeneratedAt:              s.GeneratedAt,
		Truncated:                s.Truncated,
		VehicleLimit:             s.VehicleLimit,
		LatestTelemetryRowsRead:  s.LatestTelemetryRowsRead,
		VehiclesInSnapshot:       s.VehiclesInSnapshot,
		TruncatedVehicleCountMin: s.TruncatedVehicleCountMin,
		NoTelemetry:              s.NoTelemetry,
		Vehicles:                 make([]VehicleDebugSnapshot, 0, len(s.Vehicles)),
	}
	for _, vehicle := range s.Vehicles {
		debugVehicle := VehicleDebugSnapshot{
			VehicleID:                    vehicle.VehicleID,
			TelemetryEventID:             vehicle.TelemetryEvent.ID,
			ObservedAt:                   vehicle.TelemetryEvent.Timestamp,
			TelemetryAgeSeconds:          vehicle.TelemetryAgeSeconds,
			Position:                     DebugPosition{Latitude: vehicle.TelemetryEvent.Lat, Longitude: vehicle.TelemetryEvent.Lon},
			IncludedInProtobuf:           vehicle.IncludedInProtobuf,
			HasAssignment:                vehicle.HasAssignment,
			AssignmentPublishable:        vehicle.AssignmentPublishable,
			AssignmentTelemetryMismatch:  vehicle.AssignmentTelemetryMismatch,
			TripDescriptorPublished:      vehicle.TripDescriptorPublished,
			TripDescriptorOmissionReason: vehicle.TripDescriptorOmissionReason,
		}
		if hasNumericPayloadField(vehicle.TelemetryEvent.PayloadJSON, "bearing") {
			debugVehicle.Position.Bearing = vehicle.TelemetryEvent.Bearing
		}
		if hasNumericPayloadField(vehicle.TelemetryEvent.PayloadJSON, "speed_mps") {
			debugVehicle.Position.SpeedMPS = vehicle.TelemetryEvent.SpeedMPS
		}
		if vehicle.HasAssignment {
			assignment := vehicle.Assignment
			debugVehicle.Assignment = &assignment
		}
		debug.Vehicles = append(debug.Vehicles, debugVehicle)
	}
	return debug
}

func (s VehiclePositionsSnapshot) MarshalDebugJSON() ([]byte, error) {
	payload, err := json.MarshalIndent(s.Debug(), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal vehicle positions debug json: %w", err)
	}
	return payload, nil
}

func degradedAssignmentBlocksTripDescriptor(degraded state.DegradedState) bool {
	switch degraded {
	case state.DegradedUnknown, state.DegradedStale, state.DegradedAmbiguous, state.DegradedMissingScheduleData, state.DegradedLowConfidence:
		return true
	default:
		return false
	}
}

func hasReason(assignment state.Assignment, reason string) bool {
	for _, candidate := range assignment.ReasonCodes {
		if candidate == reason {
			return true
		}
	}
	return false
}

func hasNumericPayloadField(payload json.RawMessage, field string) bool {
	if len(payload) == 0 {
		return false
	}
	var parsed map[string]any
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return false
	}
	value, ok := parsed[field]
	if !ok {
		return false
	}
	switch value.(type) {
	case float64:
		return true
	default:
		return false
	}
}
