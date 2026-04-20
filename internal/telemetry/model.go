package telemetry

import "time"

type Event struct {
	AgencyID  string    `json:"agency_id"`
	DeviceID  string    `json:"device_id"`
	VehicleID string    `json:"vehicle_id"`
	DriverID  string    `json:"driver_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	Bearing   float64   `json:"bearing,omitempty"`
	SpeedMPS  float64   `json:"speed_mps,omitempty"`
	AccuracyM float64   `json:"accuracy_m,omitempty"`
	TripHint  string    `json:"trip_hint,omitempty"`
}

func (e Event) Valid() bool {
	if e.AgencyID == "" || e.DeviceID == "" || e.VehicleID == "" {
		return false
	}
	if e.Timestamp.IsZero() {
		return false
	}
	if e.Lat < -90 || e.Lat > 90 || e.Lon < -180 || e.Lon > 180 {
		return false
	}
	return true
}
