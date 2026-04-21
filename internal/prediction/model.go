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
	ReasonPredictionsAvailable  = "predictions_available"
	ReasonNoEligiblePredictions = "no_eligible_predictions"
	ReasonPartialPredictions    = "partial_predictions"
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
	Status  string
	Reason  string
	Metrics Metrics
	Details map[string]any
}

type InputCounts struct {
	TelemetryRows     int `json:"telemetry_rows"`
	AssignmentRows    int `json:"assignment_rows"`
	TripUpdatesOutput int `json:"trip_updates_output"`
}

type Metrics struct {
	TelemetryRowsConsidered        int            `json:"telemetry_rows_considered"`
	AssignmentsConsidered          int            `json:"assignments_considered"`
	EligiblePredictionCandidates   int            `json:"eligible_prediction_candidates"`
	TripUpdatesEmitted             int            `json:"trip_updates_emitted"`
	StopUpdatesEmitted             int            `json:"stop_updates_emitted"`
	WithheldByReason               map[string]int `json:"withheld_by_reason,omitempty"`
	DegradedByReason               map[string]int `json:"degraded_by_reason,omitempty"`
	CanceledTripsEmitted           int            `json:"canceled_trips_emitted"`
	CancellationAlertLinksExpected int            `json:"cancellation_alert_links_expected"`
	CancellationAlertLinksMissing  int            `json:"cancellation_alert_links_missing"`
	AddedTripsWithheld             int            `json:"added_trips_withheld"`
	ShortTurnsWithheld             int            `json:"short_turns_withheld"`
	DetoursWithheld                int            `json:"detours_withheld"`
	CoveragePercent                *float64       `json:"coverage_percent,omitempty"`
	FutureStopCoveragePercent      *float64       `json:"future_stop_coverage_percent,omitempty"`
}

type DiagnosticsRecord struct {
	AgencyID                    string
	SnapshotAt                  time.Time
	AdapterName                 string
	Status                      string
	Reason                      string
	ActiveFeedVersionID         string
	InputCounts                 InputCounts
	Metrics                     Metrics
	AdapterDetails              map[string]any
	VehiclePositionsURL         string
	DiagnosticsPersistenceState string
}

type DiagnosticsPersistenceResult struct {
	Stored bool
}

type DiagnosticsRepository interface {
	SaveTripUpdatesDiagnostics(ctx context.Context, record DiagnosticsRecord) (DiagnosticsPersistenceResult, error)
}

type ReviewStatus string

const (
	ReviewStatusOpen     ReviewStatus = "open"
	ReviewStatusResolved ReviewStatus = "resolved"
	ReviewStatusDeferred ReviewStatus = "deferred"
)

type OverrideRecord struct {
	ID           int64
	AgencyID     string
	VehicleID    string
	OverrideType string
	RouteID      string
	TripID       string
	StartDate    string
	StartTime    string
	State        string
	ExpiresAt    *time.Time
	ClearedAt    *time.Time
	Reason       string
	CreatedBy    string
	CreatedAt    time.Time
}

type OverrideInput struct {
	AgencyID     string
	VehicleID    string
	OverrideType string
	RouteID      string
	TripID       string
	StartDate    string
	StartTime    string
	State        string
	ExpiresAt    *time.Time
	Reason       string
	ActorID      string
	Now          time.Time
}

type ReviewItem struct {
	ID         int64
	AgencyID   string
	SnapshotAt time.Time
	VehicleID  string
	RouteID    string
	TripID     string
	StartDate  string
	StartTime  string
	Severity   string
	Reason     string
	Status     ReviewStatus
	Details    map[string]any
	CreatedAt  time.Time
	ResolvedAt *time.Time
}

type ReviewFilter struct {
	AgencyID string
	Status   ReviewStatus
	Limit    int
}

type OperationsRepository interface {
	ListActivePredictionOverrides(ctx context.Context, agencyID string, at time.Time) ([]OverrideRecord, error)
	CreatePredictionOverride(ctx context.Context, input OverrideInput) (OverrideRecord, error)
	ReplacePredictionOverride(ctx context.Context, input OverrideInput) (OverrideRecord, error)
	ClearPredictionOverride(ctx context.Context, agencyID string, overrideID int64, actorID string, reason string, at time.Time) error
	SavePredictionReviewItems(ctx context.Context, items []ReviewItem) error
	ListPredictionReviewItems(ctx context.Context, filter ReviewFilter) ([]ReviewItem, error)
	UpdatePredictionReviewStatus(ctx context.Context, agencyID string, reviewID int64, status ReviewStatus, actorID string, reason string, at time.Time) error
}
