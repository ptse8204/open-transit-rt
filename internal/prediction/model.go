package prediction

import (
	"context"
	"time"

	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

const (
	StatusNoop  = "noop"
	StatusOK    = "ok"
	StatusError = "error"

	ReasonNoopAdapter           = "noop_adapter"
	ReasonActiveFeedUnavailable = "active_feed_unavailable"
	ReasonAdapterError          = "adapter_error"
)

type Adapter interface {
	Name() string
	PredictTripUpdates(ctx context.Context, request Request) (Result, error)
}

type Request struct {
	AgencyID            string
	GeneratedAt         time.Time
	ActiveFeedVersion   gtfs.FeedVersion
	Telemetry           []telemetry.StoredEvent
	Assignments         map[string]state.Assignment
	VehiclePositionsURL string
}

type Result struct {
	TripUpdates []TripUpdate
	Diagnostics Diagnostics
}

type TripUpdate struct {
	EntityID             string
	VehicleID            string
	TripID               string
	RouteID              string
	StartDate            string
	StartTime            string
	ScheduleRelationship ScheduleRelationship
	StopTimeUpdates      []StopTimeUpdate
}

type ScheduleRelationship string

const (
	ScheduleRelationshipScheduled   ScheduleRelationship = "scheduled"
	ScheduleRelationshipUnscheduled ScheduleRelationship = "unscheduled"
	ScheduleRelationshipCanceled    ScheduleRelationship = "canceled"
	ScheduleRelationshipAdded       ScheduleRelationship = "added"
)

type StopTimeUpdate struct {
	StopID                string
	StopSequence          int
	ArrivalTime           *time.Time
	DepartureTime         *time.Time
	ArrivalDelaySeconds   *int32
	DepartureDelaySeconds *int32
	ScheduleRelationship  ScheduleRelationship
}

type Diagnostics struct {
	Status string
	Reason string
}

type InputCounts struct {
	TelemetryRows     int `json:"telemetry_rows"`
	AssignmentRows    int `json:"assignment_rows"`
	TripUpdatesOutput int `json:"trip_updates_output"`
}

type DiagnosticsRecord struct {
	AgencyID                    string
	SnapshotAt                  time.Time
	AdapterName                 string
	Status                      string
	Reason                      string
	ActiveFeedVersionID         string
	InputCounts                 InputCounts
	VehiclePositionsURL         string
	DiagnosticsPersistenceState string
}

type DiagnosticsPersistenceResult struct {
	Stored bool
}

type DiagnosticsRepository interface {
	SaveTripUpdatesDiagnostics(ctx context.Context, record DiagnosticsRecord) (DiagnosticsPersistenceResult, error)
}
