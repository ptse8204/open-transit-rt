package state

import (
	"context"
	"time"

	"open-transit-rt/internal/telemetry"
)

type VehicleServiceState string

const (
	StateUnknown      VehicleServiceState = "unknown"
	StateInService    VehicleServiceState = "in_service"
	StateLayover      VehicleServiceState = "layover"
	StateDeadhead     VehicleServiceState = "deadhead"
	StateOutOfService VehicleServiceState = "out_of_service"
)

type AssignmentSource string

const (
	AssignmentSourceAutomatic      AssignmentSource = "automatic"
	AssignmentSourceManualOverride AssignmentSource = "manual_override"
	AssignmentSourceUnknown        AssignmentSource = "unknown"
)

type DegradedState string

const (
	DegradedNone                DegradedState = "none"
	DegradedUnknown             DegradedState = "unknown"
	DegradedStale               DegradedState = "stale"
	DegradedAmbiguous           DegradedState = "ambiguous"
	DegradedMissingScheduleData DegradedState = "missing_schedule_data"
	DegradedLowConfidence       DegradedState = "low_confidence"
)

const (
	ReasonManualOverrideActive   = "manual_override_active"
	ReasonTripHintMatch          = "trip_hint_match"
	ReasonRouteHintMatch         = "route_hint_match"
	ReasonShapeProximityMatch    = "shape_proximity_match"
	ReasonMovementDirectionMatch = "movement_direction_match"
	ReasonStopProgressMatch      = "stop_progress_match"
	ReasonScheduleFitMatch       = "schedule_fit_match"
	ReasonContinuityMatch        = "continuity_match"
	ReasonBlockTransitionMatch   = "block_transition_match"
	ReasonFrequencyExactInstance = "frequency_exact_instance"
	ReasonFrequencyNonExact      = "frequency_non_exact_conservative"
	ReasonLowConfidence          = "low_confidence"
	ReasonAmbiguousCandidates    = "ambiguous_candidates"
	ReasonStaleTelemetry         = "stale_telemetry"
	ReasonNoScheduleCandidates   = "no_schedule_candidates"
	ReasonOffShape               = "off_shape"
	ReasonMissingShape           = "missing_shape"
	ReasonMissingStopTimes       = "missing_stop_times"
)

const (
	IncidentAssignmentUnknown        = "assignment_unknown"
	IncidentAssignmentAmbiguous      = "assignment_ambiguous"
	IncidentStaleTelemetry           = "stale_telemetry"
	IncidentMissingScheduleData      = "missing_schedule_data"
	IncidentLowConfidenceAssignment  = "low_confidence_assignment"
	IncidentBlockTransitionAmbiguous = "block_transition_ambiguous"
)

type Assignment struct {
	ID                  int64               `json:"id,omitempty"`
	AgencyID            string              `json:"agency_id,omitempty"`
	VehicleID           string              `json:"vehicle_id"`
	FeedVersionID       string              `json:"feed_version_id,omitempty"`
	TelemetryEventID    int64               `json:"telemetry_event_id,omitempty"`
	State               VehicleServiceState `json:"state"`
	ServiceDate         string              `json:"service_date,omitempty"`
	RouteID             string              `json:"route_id,omitempty"`
	TripID              string              `json:"trip_id,omitempty"`
	BlockID             string              `json:"block_id,omitempty"`
	StartDate           string              `json:"start_date,omitempty"`
	StartTime           string              `json:"start_time,omitempty"`
	CurrentStopSequence int                 `json:"current_stop_sequence,omitempty"`
	ShapeDistTraveled   float64             `json:"shape_dist_traveled,omitempty"`
	Confidence          float64             `json:"confidence"`
	AssignmentSource    AssignmentSource    `json:"assignment_source,omitempty"`
	ReasonCodes         []string            `json:"reason_codes,omitempty"`
	DegradedState       DegradedState       `json:"degraded_state,omitempty"`
	ScoreDetails        map[string]any      `json:"score_details_json,omitempty"`
	ManualOverrideID    int64               `json:"manual_override_id,omitempty"`
	ActiveFrom          time.Time           `json:"active_from,omitempty"`
}

type Incident struct {
	Type      string
	Severity  string
	RouteID   string
	VehicleID string
	TripID    string
	Details   map[string]any
}

type ManualOverride struct {
	ID        int64
	AgencyID  string
	VehicleID string
	Type      string
	RouteID   string
	TripID    string
	StartDate string
	StartTime string
	State     VehicleServiceState
	Reason    string
	ExpiresAt *time.Time
	CreatedAt time.Time
}

type Repository interface {
	ActiveManualOverride(ctx context.Context, agencyID string, vehicleID string, at time.Time) (*ManualOverride, error)
	CurrentAssignment(ctx context.Context, agencyID string, vehicleID string) (*Assignment, error)
	SaveAssignment(ctx context.Context, assignment Assignment, incidents []Incident) (Assignment, error)
}

type Matcher interface {
	Assign(event telemetry.Event) Assignment
}

type RuleBasedMatcher struct{}

func NewRuleBasedMatcher() *RuleBasedMatcher { return &RuleBasedMatcher{} }

func (m *RuleBasedMatcher) Assign(event telemetry.Event) Assignment {
	return Assignment{
		AgencyID:         event.AgencyID,
		VehicleID:        event.VehicleID,
		State:            StateUnknown,
		Confidence:       0.10,
		AssignmentSource: AssignmentSourceUnknown,
		ReasonCodes:      []string{ReasonNoScheduleCandidates},
		DegradedState:    DegradedUnknown,
		ScoreDetails:     map[string]any{"note": "placeholder matcher path; persisted deterministic matching uses Engine"},
		TelemetryEventID: 0,
	}
}
