package state

import "open-transit-rt/internal/telemetry"

type VehicleServiceState string

const (
	StateUnknown   VehicleServiceState = "unknown"
	StateInService VehicleServiceState = "in_service"
	StateLayover   VehicleServiceState = "layover"
	StateDeadhead  VehicleServiceState = "deadhead"
)

type Assignment struct {
	VehicleID           string              `json:"vehicle_id"`
	State               VehicleServiceState `json:"state"`
	RouteID             string              `json:"route_id,omitempty"`
	TripID              string              `json:"trip_id,omitempty"`
	StartDate           string              `json:"start_date,omitempty"`
	StartTime           string              `json:"start_time,omitempty"`
	CurrentStopSequence int                 `json:"current_stop_sequence,omitempty"`
	Confidence          float64             `json:"confidence"`
}

type Matcher interface {
	Assign(event telemetry.Event) Assignment
}

type RuleBasedMatcher struct{}

func NewRuleBasedMatcher() *RuleBasedMatcher { return &RuleBasedMatcher{} }

func (m *RuleBasedMatcher) Assign(event telemetry.Event) Assignment {
	return Assignment{
		VehicleID:  event.VehicleID,
		State:      StateUnknown,
		Confidence: 0.10,
	}
}
