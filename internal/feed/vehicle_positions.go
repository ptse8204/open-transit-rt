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
	TripDescriptorOmissionSuppressedStaleTelemetry    = "suppressed_stale_telemetry"
	TripDescriptorOmissionStaleTelemetry              = "stale_telemetry"
	TripDescriptorOmissionNoAssignment                = "no_assignment"
	TripDescriptorOmissionAssignmentTelemetryMismatch = "assignment_telemetry_mismatch"
	TripDescriptorOmissionNotInService                = "not_in_service"
	TripDescriptorOmissionMissingTripID               = "missing_trip_id"
	TripDescriptorOmissionManualStateWithoutTrip      = "manual_state_without_trip"
	TripDescriptorOmissionDegradedAssignment          = "degraded_assignment"
	TripDescriptorOmissionMissingScheduleContext      = "missing_schedule_context"
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

func (b *VehiclePositionsBuilder) SnapshotForAgency(ctx context.Context, agencyID string, generatedAt time.Time) (VehiclePositionsSnapshot, error) {
	if agencyID == "" {
		return VehiclePositionsSnapshot{}, fmt.Errorf("agency_id is required")
	}
	if agencyID == b.config.AgencyID {
		return b.Snapshot(ctx, generatedAt)
	}
	config := b.config
	config.AgencyID = agencyID
	return (&VehiclePositionsBuilder{telemetry: b.telemetry, assignments: b.assignments, config: config}).Snapshot(ctx, generatedAt)
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
	if assignment.DegradedState == state.DegradedMissingScheduleData {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionMissingScheduleContext
		return vehicle
	}
	if degradedAssignmentBlocksTripDescriptor(assignment.DegradedState) {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionDegradedAssignment
		return vehicle
	}
	if assignment.State != state.StateInService {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionNotInService
		return vehicle
	}
	if assignment.AssignmentSource == state.AssignmentSourceManualOverride && assignment.TripID == "" {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionManualStateWithoutTrip
		return vehicle
	}
	if assignment.TripID == "" {
		vehicle.TripDescriptorOmissionReason = TripDescriptorOmissionMissingTripID
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
	ReviewSummary            VehiclePositionsReview `json:"review_summary"`
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
	PublicFeedDiagnostic         string            `json:"public_feed_diagnostic"`
	Assignment                   *state.Assignment `json:"assignment,omitempty"`
}

type DebugPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Bearing   float64 `json:"bearing,omitempty"`
	SpeedMPS  float64 `json:"speed_mps,omitempty"`
}

type VehiclePositionsReview struct {
	VehiclesInSnapshot             int            `json:"vehicles_in_snapshot"`
	VehiclesInProtobuf             int            `json:"vehicles_in_protobuf"`
	TripDescriptorsPublished       int            `json:"trip_descriptors_published"`
	TripDescriptorOmissions        map[string]int `json:"trip_descriptor_omissions"`
	StaleVehicles                  int            `json:"stale_vehicles"`
	SuppressedVehicles             int            `json:"suppressed_vehicles"`
	UnmatchedVehicles              int            `json:"unmatched_vehicles"`
	MissingScheduleContextVehicles int            `json:"missing_schedule_context_vehicles"`
	LowConfidenceVehicles          int            `json:"low_confidence_vehicles"`
	AssignmentTelemetryMismatches  int            `json:"assignment_telemetry_mismatches"`
	Truncated                      bool           `json:"truncated"`
	TruncatedVehicleCountMin       int            `json:"truncated_vehicle_count_min"`
	LatestTelemetryRowsRead        int            `json:"latest_telemetry_rows_read"`
	NoTelemetry                    bool           `json:"no_telemetry"`
	UsefulnessStatus               string         `json:"usefulness_status"`
	UsefulnessReasons              []string       `json:"usefulness_reasons"`
}

func (s VehiclePositionsSnapshot) ReviewSummary() VehiclePositionsReview {
	review := VehiclePositionsReview{
		VehiclesInSnapshot:       s.VehiclesInSnapshot,
		TripDescriptorOmissions:  map[string]int{},
		Truncated:                s.Truncated,
		TruncatedVehicleCountMin: s.TruncatedVehicleCountMin,
		LatestTelemetryRowsRead:  s.LatestTelemetryRowsRead,
		NoTelemetry:              s.NoTelemetry,
	}
	for _, vehicle := range s.Vehicles {
		if vehicle.IncludedInProtobuf {
			review.VehiclesInProtobuf++
		}
		if vehicle.TripDescriptorPublished {
			review.TripDescriptorsPublished++
		}
		if vehicle.TripDescriptorOmissionReason != "" && vehicle.TripDescriptorOmissionReason != TripDescriptorOmissionNone {
			review.TripDescriptorOmissions[vehicle.TripDescriptorOmissionReason]++
		}
		switch vehicle.TripDescriptorOmissionReason {
		case TripDescriptorOmissionStaleTelemetry, TripDescriptorOmissionSuppressedStaleTelemetry:
			review.StaleVehicles++
		}
		if vehicle.TripDescriptorOmissionReason == TripDescriptorOmissionSuppressedStaleTelemetry {
			review.SuppressedVehicles++
		}
		if vehicle.TripDescriptorOmissionReason == TripDescriptorOmissionNoAssignment {
			review.UnmatchedVehicles++
		}
		if vehicle.TripDescriptorOmissionReason == TripDescriptorOmissionMissingScheduleContext {
			review.MissingScheduleContextVehicles++
		}
		if vehicle.TripDescriptorOmissionReason == TripDescriptorOmissionBelowPublicationConfidence {
			review.LowConfidenceVehicles++
		}
		if vehicle.AssignmentTelemetryMismatch {
			review.AssignmentTelemetryMismatches++
		}
	}
	review.UsefulnessStatus, review.UsefulnessReasons = vehiclePositionsUsefulness(review)
	return review
}

func (s VehiclePositionsSnapshot) Debug() VehiclePositionsDebug {
	debug := VehiclePositionsDebug{
		AgencyID:                 s.AgencyID,
		GeneratedAt:              s.GeneratedAt,
		ReviewSummary:            s.ReviewSummary(),
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
			PublicFeedDiagnostic:         vehiclePublicFeedDiagnostic(vehicle),
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

func vehiclePositionsUsefulness(review VehiclePositionsReview) (string, []string) {
	if review.NoTelemetry || review.VehiclesInSnapshot == 0 {
		return "empty_no_telemetry", []string{"no accepted latest telemetry rows are available"}
	}
	var reasons []string
	if review.SuppressedVehicles > 0 {
		reasons = append(reasons, "some stale vehicles were suppressed")
	}
	if review.StaleVehicles > 0 {
		reasons = append(reasons, "stale telemetry is present")
	}
	if review.UnmatchedVehicles > 0 {
		reasons = append(reasons, "some vehicles have no current assignment")
	}
	if review.MissingScheduleContextVehicles > 0 {
		reasons = append(reasons, "active schedule context is missing for some assignments")
	}
	if review.LowConfidenceVehicles > 0 {
		reasons = append(reasons, "some assignments are below the publication confidence threshold")
	}
	if review.AssignmentTelemetryMismatches > 0 {
		reasons = append(reasons, "some assignments do not match the latest telemetry event")
	}
	if review.VehiclesInProtobuf == 0 {
		return "empty_after_suppression", appendIfMissing(reasons, "no vehicles are emitted in the protobuf")
	}
	if review.TripDescriptorsPublished == 0 {
		return "positions_without_trip_context", appendIfMissing(reasons, "vehicle positions emit without trip descriptors")
	}
	if len(reasons) > 0 {
		return "positions_with_partial_trip_context", reasons
	}
	return "positions_with_trip_context", []string{"fresh vehicle positions include conservative trip descriptors"}
}

func appendIfMissing(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func vehiclePublicFeedDiagnostic(vehicle VehicleSnapshot) string {
	switch vehicle.TripDescriptorOmissionReason {
	case TripDescriptorOmissionNone:
		if vehicle.TripDescriptorPublished {
			return "trip descriptor included from conservative assignment evidence"
		}
		return "vehicle position emitted without an omission reason"
	case TripDescriptorOmissionSuppressedStaleTelemetry:
		return "vehicle omitted from protobuf because telemetry is beyond the stale suppression threshold"
	case TripDescriptorOmissionStaleTelemetry:
		return "vehicle emitted without trip descriptor because telemetry is stale"
	case TripDescriptorOmissionNoAssignment:
		return "vehicle emitted without trip descriptor because no current assignment is available"
	case TripDescriptorOmissionAssignmentTelemetryMismatch:
		return "vehicle emitted without trip descriptor because assignment evidence does not match the latest telemetry event"
	case TripDescriptorOmissionMissingScheduleContext:
		return "vehicle emitted without trip descriptor because active schedule context is missing"
	case TripDescriptorOmissionBelowPublicationConfidence:
		return "vehicle emitted without trip descriptor because assignment confidence is below the publication threshold"
	case TripDescriptorOmissionDegradedAssignment:
		return "vehicle emitted without trip descriptor because assignment state is degraded"
	case TripDescriptorOmissionNotInService:
		return "vehicle emitted without trip descriptor because assignment state is not in service"
	case TripDescriptorOmissionMissingTripID:
		return "vehicle emitted without trip descriptor because the assignment has no trip_id"
	case TripDescriptorOmissionManualStateWithoutTrip:
		return "vehicle emitted without trip descriptor because manual state has no trip_id"
	default:
		return "vehicle position emitted with conservative trip context review needed"
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
