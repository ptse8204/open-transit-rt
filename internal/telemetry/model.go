package telemetry

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type IngestStatus string

const (
	IngestStatusAccepted   IngestStatus = "accepted"
	IngestStatusDuplicate  IngestStatus = "duplicate"
	IngestStatusOutOfOrder IngestStatus = "out_of_order"
	IngestStatusRejected   IngestStatus = "rejected"
)

const (
	DefaultStaleSourceThreshold = 90 * time.Second
	DefaultFutureThreshold      = 30 * time.Second
	DefaultLowAccuracyMeters    = 50.0
	QualityFlagStaleTimestamp   = "stale_timestamp"
	QualityFlagFutureTimestamp  = "future_timestamp"
	QualityFlagLowGPSAccuracy   = "low_gps_accuracy"
)

var ErrUnknownAgency = errors.New("unknown agency")

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

type StoredEvent struct {
	ID int64 `json:"id"`
	Event
	ReceivedAt   time.Time       `json:"received_at"`
	IngestStatus IngestStatus    `json:"ingest_status"`
	PayloadJSON  json.RawMessage `json:"payload_json,omitempty"`
}

type StoreResult struct {
	StoredEvent
}

type Repository interface {
	Store(ctx context.Context, event Event, payload json.RawMessage) (StoreResult, error)
	LatestByVehicle(ctx context.Context, agencyID string, vehicleID string) (StoredEvent, error)
	ListLatestByAgency(ctx context.Context, agencyID string, limit int) ([]StoredEvent, error)
	ListEvents(ctx context.Context, agencyID string, limit int) ([]StoredEvent, error)
}

func (e Event) Valid() bool {
	if strings.TrimSpace(e.AgencyID) != e.AgencyID || strings.TrimSpace(e.DeviceID) != e.DeviceID || strings.TrimSpace(e.VehicleID) != e.VehicleID {
		return false
	}
	if e.AgencyID == "" || e.DeviceID == "" || e.VehicleID == "" {
		return false
	}
	if e.Timestamp.IsZero() {
		return false
	}
	if e.Lat < -90 || e.Lat > 90 || e.Lon < -180 || e.Lon > 180 {
		return false
	}
	if e.Bearing < 0 || e.Bearing > 360 || e.SpeedMPS < 0 || e.AccuracyM < 0 {
		return false
	}
	return true
}

func (e Event) QualityFlags(referenceTime time.Time) []string {
	if referenceTime.IsZero() {
		referenceTime = time.Now().UTC()
	}
	referenceTime = referenceTime.UTC()
	timestamp := e.Timestamp.UTC()
	var flags []string
	if timestamp.Before(referenceTime.Add(-DefaultStaleSourceThreshold)) {
		flags = append(flags, QualityFlagStaleTimestamp)
	}
	if timestamp.After(referenceTime.Add(DefaultFutureThreshold)) {
		flags = append(flags, QualityFlagFutureTimestamp)
	}
	if e.AccuracyM > DefaultLowAccuracyMeters {
		flags = append(flags, QualityFlagLowGPSAccuracy)
	}
	return flags
}

func HasQualityFlag(flags []string, target string) bool {
	for _, flag := range flags {
		if flag == target {
			return true
		}
	}
	return false
}
