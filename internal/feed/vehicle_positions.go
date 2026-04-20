package feed

import (
	"encoding/json"
	"time"

	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

type VehiclePositionEntity struct {
	ID        string           `json:"id"`
	Timestamp time.Time        `json:"timestamp"`
	Trip      state.Assignment `json:"trip"`
	Position  Position         `json:"position"`
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Bearing   float64 `json:"bearing,omitempty"`
	SpeedMPS  float64 `json:"speed_mps,omitempty"`
}

type VehiclePositionsFeed struct {
	GeneratedAt time.Time               `json:"generated_at"`
	Entities    []VehiclePositionEntity `json:"entities"`
	Note        string                  `json:"note"`
}

func BuildVehiclePositionsJSON(events []telemetry.Event, assignments map[string]state.Assignment) ([]byte, error) {
	entities := make([]VehiclePositionEntity, 0, len(events))
	for _, e := range events {
		assignment := assignments[e.VehicleID]
		entities = append(entities, VehiclePositionEntity{
			ID:        e.VehicleID,
			Timestamp: e.Timestamp,
			Trip:      assignment,
			Position: Position{
				Latitude:  e.Lat,
				Longitude: e.Lon,
				Bearing:   e.Bearing,
				SpeedMPS:  e.SpeedMPS,
			},
		})
	}

	payload := VehiclePositionsFeed{
		GeneratedAt: time.Now().UTC(),
		Entities:    entities,
		Note:        "JSON placeholder only. Replace with canonical GTFS-RT protobuf output in the next implementation pass.",
	}
	return json.MarshalIndent(payload, "", "  ")
}
