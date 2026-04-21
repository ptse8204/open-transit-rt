package alerts

import (
	"context"
	"time"
)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusArchived  = "archived"

	SourceOperator               = "operator"
	SourceCancellationReconciler = "cancellation_reconciler"
)

type Alert struct {
	ID              int64
	AgencyID        string
	AlertKey        string
	Status          string
	Cause           string
	Effect          string
	HeaderText      string
	DescriptionText string
	URL             string
	ActiveStart     *time.Time
	ActiveEnd       *time.Time
	FeedVersionID   string
	SourceType      string
	SourceID        string
	Metadata        map[string]any
	CreatedBy       string
	UpdatedBy       string
	PublishedAt     *time.Time
	ArchivedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Entities        []InformedEntity
}

type InformedEntity struct {
	ID             int64
	ServiceAlertID int64
	AgencyID       string
	RouteID        string
	StopID         string
	TripID         string
	StartDate      string
	StartTime      string
	Metadata       map[string]any
	CreatedAt      time.Time
}

type UpsertInput struct {
	AgencyID        string
	AlertKey        string
	Cause           string
	Effect          string
	HeaderText      string
	DescriptionText string
	URL             string
	ActiveStart     *time.Time
	ActiveEnd       *time.Time
	FeedVersionID   string
	SourceType      string
	SourceID        string
	Metadata        map[string]any
	ActorID         string
	Entities        []InformedEntity
	Publish         bool
	Now             time.Time
}

type ListFilter struct {
	AgencyID      string
	Status        string
	PublishedOnly bool
	At            time.Time
	Limit         int
}

type Repository interface {
	UpsertAlert(ctx context.Context, input UpsertInput) (Alert, error)
	PublishAlert(ctx context.Context, agencyID string, alertID int64, actorID string, at time.Time) (Alert, error)
	ArchiveAlert(ctx context.Context, agencyID string, alertID int64, actorID string, reason string, at time.Time) error
	ListAlerts(ctx context.Context, filter ListFilter) ([]Alert, error)
}
