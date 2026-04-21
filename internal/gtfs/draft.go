package gtfs

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	DraftStatusDraft     = "draft"
	DraftStatusPublished = "published"
	DraftStatusDiscarded = "discarded"
)

var (
	ErrDraftNotEditable = errors.New("gtfs draft is not editable")
	ErrDraftNotFound    = errors.New("gtfs draft not found")
)

type DraftService struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

type CreateDraftOptions struct {
	AgencyID string
	Name     string
	ActorID  string
	Blank    bool
}

type DiscardDraftOptions struct {
	DraftID string
	ActorID string
	Reason  string
}

type PublishDraftOptions struct {
	DraftID string
	ActorID string
	Notes   string
}

type PublishDraftResult struct {
	PublishID      int64          `json:"publish_id,omitempty"`
	DraftID        string         `json:"draft_id"`
	AgencyID       string         `json:"agency_id"`
	FeedVersionID  string         `json:"feed_version_id,omitempty"`
	Status         string         `json:"status"`
	ErrorCount     int            `json:"error_count"`
	WarningCount   int            `json:"warning_count"`
	InfoCount      int            `json:"info_count"`
	Counts         map[string]int `json:"counts,omitempty"`
	ReportStored   bool           `json:"report_stored"`
	FailureMessage string         `json:"failure_message,omitempty"`
}

type DraftPublishError struct {
	Result PublishDraftResult
	Err    error
}

func (e *DraftPublishError) Error() string {
	return e.Err.Error()
}

func (e *DraftPublishError) Unwrap() error {
	return e.Err
}

type Draft struct {
	ID                         string
	AgencyID                   string
	Name                       string
	Status                     string
	BaseFeedVersionID          string
	LastPublishedFeedVersionID string
	LastPublishAttemptID       int64
	DiscardedAt                *time.Time
	DiscardedBy                string
	DiscardReason              string
	CreatedBy                  string
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

type DraftSummary struct {
	Draft
	LatestPublishStatus string
	LatestPublishID     int64
}

type DraftAgency struct {
	DraftID      string
	AgencyID     string
	Name         string
	Timezone     string
	ContactEmail string
	PublicURL    string
}

type DraftRoute struct {
	DraftID   string
	AgencyID  string
	ID        string
	ShortName string
	LongName  string
	RouteType int
}

type DraftStop struct {
	DraftID  string
	AgencyID string
	ID       string
	Name     string
	Lat      float64
	Lon      float64
}

type DraftCalendar struct {
	DraftID   string
	AgencyID  string
	ServiceID string
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Friday    bool
	Saturday  bool
	Sunday    bool
	StartDate string
	EndDate   string
}

type DraftCalendarDate struct {
	DraftID       string
	AgencyID      string
	ServiceID     string
	Date          string
	ExceptionType int
}

type DraftTrip struct {
	DraftID     string
	AgencyID    string
	ID          string
	RouteID     string
	ServiceID   string
	BlockID     string
	ShapeID     string
	DirectionID *int
}

type DraftStopTime struct {
	DraftID           string
	AgencyID          string
	TripID            string
	ArrivalTime       string
	DepartureTime     string
	StopID            string
	StopSequence      int
	PickupType        *int
	DropOffType       *int
	ShapeDistTraveled *float64
}

type DraftShapePoint struct {
	DraftID      string
	AgencyID     string
	ShapeID      string
	Lat          float64
	Lon          float64
	Sequence     int
	DistTraveled *float64
}

type DraftFrequency struct {
	DraftID     string
	AgencyID    string
	TripID      string
	StartTime   string
	EndTime     string
	HeadwaySecs int
	ExactTimes  int
}

func NewDraftService(pool *pgxpool.Pool) *DraftService {
	return &DraftService{
		pool: pool,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *DraftService) CreateDraft(ctx context.Context, opts CreateDraftOptions) (Draft, error) {
	if s == nil || s.pool == nil {
		return Draft{}, fmt.Errorf("gtfs draft service requires a database pool")
	}
	opts.AgencyID = strings.TrimSpace(opts.AgencyID)
	opts.Name = strings.TrimSpace(opts.Name)
	if opts.AgencyID == "" {
		return Draft{}, fmt.Errorf("agency_id is required")
	}
	if opts.Name == "" {
		return Draft{}, fmt.Errorf("draft name is required")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Draft{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	baseFeedVersionID := ""
	if !opts.Blank {
		baseFeedVersionID, err = activeFeedVersionID(ctx, tx, opts.AgencyID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return Draft{}, fmt.Errorf("query active feed for draft clone: %w", err)
		}
	}

	draftID := "draft-" + randomHex(8)
	now := s.now()
	draft := Draft{
		ID:                draftID,
		AgencyID:          opts.AgencyID,
		Name:              opts.Name,
		Status:            DraftStatusDraft,
		BaseFeedVersionID: baseFeedVersionID,
		CreatedBy:         opts.ActorID,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO gtfs_draft (id, agency_id, name, status, base_feed_version_id, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, 'draft', $4, $5, $6, $6)
	`, draftID, opts.AgencyID, opts.Name, nullString(baseFeedVersionID), nullString(opts.ActorID), now); err != nil {
		return Draft{}, fmt.Errorf("insert draft metadata: %w", err)
	}

	if baseFeedVersionID != "" {
		if err := clonePublishedFeedToDraft(ctx, tx, opts.AgencyID, baseFeedVersionID, draftID); err != nil {
			return Draft{}, err
		}
	} else {
		if err := insertBlankDraftAgency(ctx, tx, opts.AgencyID, draftID); err != nil {
			return Draft{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Draft{}, fmt.Errorf("commit draft create: %w", err)
	}
	return draft, nil
}

func (s *DraftService) ListDrafts(ctx context.Context, agencyID string, includeDiscarded bool) ([]DraftSummary, error) {
	query := `
		SELECT d.id, d.agency_id, d.name, d.status, d.base_feed_version_id,
		       d.last_published_feed_version_id, d.last_publish_attempt_id,
		       d.discarded_at, d.discarded_by, d.discard_reason,
		       d.created_by, d.created_at, d.updated_at,
		       COALESCE(p.status, ''), COALESCE(p.id, 0)
		FROM gtfs_draft d
		LEFT JOIN LATERAL (
			SELECT id, status
			FROM gtfs_draft_publish
			WHERE draft_id = d.id
			ORDER BY started_at DESC, id DESC
			LIMIT 1
		) p ON true
		WHERE d.agency_id = $1
		  AND ($2 OR d.status <> 'discarded')
		ORDER BY d.updated_at DESC, d.created_at DESC
	`
	rows, err := s.pool.Query(ctx, query, agencyID, includeDiscarded)
	if err != nil {
		return nil, fmt.Errorf("list drafts: %w", err)
	}
	defer rows.Close()

	var drafts []DraftSummary
	for rows.Next() {
		var summary DraftSummary
		if err := scanDraft(rows, &summary.Draft, &summary.LatestPublishStatus, &summary.LatestPublishID); err != nil {
			return nil, err
		}
		drafts = append(drafts, summary)
	}
	return drafts, rows.Err()
}

func (s *DraftService) GetDraft(ctx context.Context, draftID string) (Draft, error) {
	var draft Draft
	err := s.pool.QueryRow(ctx, `
		SELECT id, agency_id, name, status, base_feed_version_id,
		       last_published_feed_version_id, last_publish_attempt_id,
		       discarded_at, discarded_by, discard_reason,
		       created_by, created_at, updated_at
		FROM gtfs_draft
		WHERE id = $1
	`, draftID).Scan(
		&draft.ID, &draft.AgencyID, &draft.Name, &draft.Status,
		nullStringScan(&draft.BaseFeedVersionID),
		nullStringScan(&draft.LastPublishedFeedVersionID),
		nullInt64Scan(&draft.LastPublishAttemptID),
		&draft.DiscardedAt,
		nullStringScan(&draft.DiscardedBy),
		nullStringScan(&draft.DiscardReason),
		nullStringScan(&draft.CreatedBy),
		&draft.CreatedAt,
		&draft.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Draft{}, ErrDraftNotFound
	}
	if err != nil {
		return Draft{}, fmt.Errorf("get draft: %w", err)
	}
	return draft, nil
}

func (s *DraftService) DiscardDraft(ctx context.Context, opts DiscardDraftOptions) error {
	draft, err := s.ensureEditable(ctx, opts.DraftID)
	if err != nil {
		return err
	}
	now := s.now()
	tag, err := s.pool.Exec(ctx, `
		UPDATE gtfs_draft
		SET status = 'discarded',
		    discarded_at = $2,
		    discarded_by = $3,
		    discard_reason = $4,
		    updated_at = $2
		WHERE id = $1
		  AND status = 'draft'
	`, draft.ID, now, nullString(opts.ActorID), nullString(opts.Reason))
	if err != nil {
		return fmt.Errorf("discard draft: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return ErrDraftNotEditable
	}
	return nil
}

func (s *DraftService) PublishDraft(ctx context.Context, opts PublishDraftOptions) (PublishDraftResult, error) {
	draft, err := s.ensureEditable(ctx, opts.DraftID)
	if err != nil {
		result := PublishDraftResult{DraftID: opts.DraftID, Status: ImportStatusFailed, FailureMessage: err.Error()}
		return result, &DraftPublishError{Result: result, Err: err}
	}

	feed, report, err := s.draftFeed(ctx, draft)
	if err != nil {
		result := PublishDraftResult{DraftID: draft.ID, AgencyID: draft.AgencyID, Status: ImportStatusFailed, FailureMessage: err.Error()}
		return result, &DraftPublishError{Result: result, Err: err}
	}
	report.Metadata = map[string]string{
		"source":   "gtfs_studio",
		"draft_id": draft.ID,
	}

	publishID, err := s.insertDraftPublishAttempt(ctx, draft, opts, report)
	if err != nil {
		result := publishResult(draft, 0, "", ImportStatusFailed, report, false, "failed to store draft publish attempt")
		return result, &DraftPublishError{Result: result, Err: fmt.Errorf("store draft publish attempt: %w", err)}
	}

	if report.hasErrors() {
		if markErr := s.markDraftPublishFailed(ctx, draft.ID, publishID, report); markErr != nil {
			result := publishResult(draft, publishID, "", ImportStatusFailed, report, false, "failed to store validation failure report")
			return result, &DraftPublishError{Result: result, Err: fmt.Errorf("draft validation failed and failure report could not be stored: %w", markErr)}
		}
		if reportErr := insertValidationReport(ctx, txOrPool{pool: s.pool}, draft.AgencyID, "", report, validationReportLinks{DraftPublishID: publishID}); reportErr != nil {
			result := publishResult(draft, publishID, "", ImportStatusFailed, report, false, "failed to store validation report")
			return result, &DraftPublishError{Result: result, Err: fmt.Errorf("draft validation failed and validation report could not be stored: %w", reportErr)}
		}
		result := publishResult(draft, publishID, "", ImportStatusFailed, report, true, "validation failed")
		return result, &DraftPublishError{Result: result, Err: fmt.Errorf("draft validation failed with %d error(s)", len(report.Errors))}
	}

	feedVersionID := fmt.Sprintf("gtfs-studio-%d", publishID)
	if err := publishFeed(ctx, s.pool, s.now, publishFeedOptions{
		AgencyID:      draft.AgencyID,
		FeedVersionID: feedVersionID,
		SourceType:    "gtfs_studio",
		ActorID:       opts.ActorID,
		Notes:         opts.Notes,
		AuditAction:   "gtfs_studio_publish",
		EntityType:    "gtfs_draft",
		EntityID:      draft.ID,
		Feed:          feed,
		Report:        report,
		Links:         validationReportLinks{DraftPublishID: publishID},
		AfterPublish: func(ctx context.Context, tx pgx.Tx, publishedAt time.Time, report ImportReport, reportJSON []byte) error {
			if _, err := tx.Exec(ctx, `
				UPDATE gtfs_draft_publish
				SET status = 'published',
				    feed_version_id = $2,
				    error_count = $3,
				    warning_count = $4,
				    info_count = $5,
				    report_json = $6,
				    completed_at = $7
				WHERE id = $1
			`, publishID, feedVersionID, len(report.Errors), len(report.Warnings), len(report.Info), reportJSON, publishedAt); err != nil {
				return fmt.Errorf("mark draft publish published: %w", err)
			}
			if _, err := tx.Exec(ctx, `
				UPDATE gtfs_draft
				SET status = 'published',
				    last_published_feed_version_id = $2,
				    last_publish_attempt_id = $3,
				    updated_at = $4
				WHERE id = $1
			`, draft.ID, feedVersionID, publishID, publishedAt); err != nil {
				return fmt.Errorf("mark draft published: %w", err)
			}
			return nil
		},
	}); err != nil {
		failedReport := reportWithError(report, ImportMessage{Code: "publish_failed", Message: err.Error()})
		if markErr := s.markDraftPublishFailed(ctx, draft.ID, publishID, failedReport); markErr != nil {
			result := publishResult(draft, publishID, "", ImportStatusFailed, failedReport, false, "publish failed and failure report could not be stored")
			return result, &DraftPublishError{Result: result, Err: fmt.Errorf("draft publish failed and failure report could not be stored: publish error: %v; report error: %w", err, markErr)}
		}
		if reportErr := insertValidationReport(ctx, txOrPool{pool: s.pool}, draft.AgencyID, "", failedReport, validationReportLinks{DraftPublishID: publishID}); reportErr != nil {
			result := publishResult(draft, publishID, "", ImportStatusFailed, failedReport, false, "publish failed and validation report could not be stored")
			return result, &DraftPublishError{Result: result, Err: fmt.Errorf("draft publish failed and validation report could not be stored: publish error: %v; report error: %w", err, reportErr)}
		}
		result := publishResult(draft, publishID, "", ImportStatusFailed, failedReport, true, "publish failed")
		return result, &DraftPublishError{Result: result, Err: fmt.Errorf("draft publish failed: %w", err)}
	}

	return publishResult(draft, publishID, feedVersionID, ImportStatusPublished, report, true, ""), nil
}

func (s *DraftService) UpsertAgency(ctx context.Context, agency DraftAgency) error {
	draft, err := s.ensureEditable(ctx, agency.DraftID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(agency.AgencyID) == "" {
		agency.AgencyID = draft.AgencyID
	}
	if agency.AgencyID != draft.AgencyID {
		return fmt.Errorf("draft agency metadata must match draft agency_id")
	}
	if strings.TrimSpace(agency.Name) == "" || strings.TrimSpace(agency.Timezone) == "" {
		return fmt.Errorf("draft agency name and timezone are required")
	}
	if _, err := time.LoadLocation(agency.Timezone); err != nil {
		return fmt.Errorf("invalid draft agency timezone: %w", err)
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO gtfs_draft_agency (draft_id, agency_id, name, timezone, contact_email, public_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (draft_id) DO UPDATE
		SET name = EXCLUDED.name,
		    timezone = EXCLUDED.timezone,
		    contact_email = EXCLUDED.contact_email,
		    public_url = EXCLUDED.public_url,
		    updated_at = now()
	`, draft.ID, draft.AgencyID, agency.Name, agency.Timezone, nullString(agency.ContactEmail), nullString(agency.PublicURL))
	return err
}

func (s *DraftService) GetAgency(ctx context.Context, draftID string) (DraftAgency, error) {
	var agency DraftAgency
	err := s.pool.QueryRow(ctx, `
		SELECT draft_id, agency_id, name, timezone, contact_email, public_url
		FROM gtfs_draft_agency
		WHERE draft_id = $1
	`, draftID).Scan(&agency.DraftID, &agency.AgencyID, &agency.Name, &agency.Timezone, nullStringScan(&agency.ContactEmail), nullStringScan(&agency.PublicURL))
	if err != nil {
		return DraftAgency{}, err
	}
	return agency, nil
}

func (s *DraftService) UpsertRoute(ctx context.Context, route DraftRoute) error {
	draft, err := s.ensureEditable(ctx, route.DraftID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(route.ID) == "" {
		return fmt.Errorf("route id is required")
	}
	if !supportedRouteType(route.RouteType) {
		return fmt.Errorf("route_type must be 0-7 or 100-1702")
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO gtfs_draft_route (draft_id, agency_id, id, short_name, long_name, route_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (draft_id, id) DO UPDATE
		SET short_name = EXCLUDED.short_name,
		    long_name = EXCLUDED.long_name,
		    route_type = EXCLUDED.route_type,
		    updated_at = now()
	`, draft.ID, draft.AgencyID, route.ID, nullString(route.ShortName), nullString(route.LongName), route.RouteType)
	return err
}

func (s *DraftService) ListRoutes(ctx context.Context, draftID string) ([]DraftRoute, error) {
	rows, err := s.pool.Query(ctx, `SELECT draft_id, agency_id, id, short_name, long_name, route_type FROM gtfs_draft_route WHERE draft_id = $1 ORDER BY id`, draftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var routes []DraftRoute
	for rows.Next() {
		var r DraftRoute
		if err := rows.Scan(&r.DraftID, &r.AgencyID, &r.ID, nullStringScan(&r.ShortName), nullStringScan(&r.LongName), &r.RouteType); err != nil {
			return nil, err
		}
		routes = append(routes, r)
	}
	return routes, rows.Err()
}

func (s *DraftService) RemoveRoute(ctx context.Context, draftID string, id string) error {
	return s.removeDraftRow(ctx, draftID, `DELETE FROM gtfs_draft_route WHERE draft_id = $1 AND id = $2`, id)
}

func (s *DraftService) UpsertStop(ctx context.Context, stop DraftStop) error {
	draft, err := s.ensureEditable(ctx, stop.DraftID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(stop.ID) == "" || strings.TrimSpace(stop.Name) == "" {
		return fmt.Errorf("stop id and name are required")
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO gtfs_draft_stop (draft_id, agency_id, id, name, lat, lon)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (draft_id, id) DO UPDATE
		SET name = EXCLUDED.name, lat = EXCLUDED.lat, lon = EXCLUDED.lon, updated_at = now()
	`, draft.ID, draft.AgencyID, stop.ID, stop.Name, stop.Lat, stop.Lon)
	return err
}

func (s *DraftService) ListStops(ctx context.Context, draftID string) ([]DraftStop, error) {
	rows, err := s.pool.Query(ctx, `SELECT draft_id, agency_id, id, name, lat, lon FROM gtfs_draft_stop WHERE draft_id = $1 ORDER BY id`, draftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stops []DraftStop
	for rows.Next() {
		var stop DraftStop
		if err := rows.Scan(&stop.DraftID, &stop.AgencyID, &stop.ID, &stop.Name, &stop.Lat, &stop.Lon); err != nil {
			return nil, err
		}
		stops = append(stops, stop)
	}
	return stops, rows.Err()
}

func (s *DraftService) RemoveStop(ctx context.Context, draftID string, id string) error {
	return s.removeDraftRow(ctx, draftID, `DELETE FROM gtfs_draft_stop WHERE draft_id = $1 AND id = $2`, id)
}

func (s *DraftService) UpsertCalendar(ctx context.Context, calendar DraftCalendar) error {
	draft, err := s.ensureEditable(ctx, calendar.DraftID)
	if err != nil {
		return err
	}
	if err := validateDateText(calendar.StartDate); err != nil {
		return fmt.Errorf("invalid start_date: %w", err)
	}
	if err := validateDateText(calendar.EndDate); err != nil {
		return fmt.Errorf("invalid end_date: %w", err)
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO gtfs_draft_calendar (
			draft_id, agency_id, service_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday, start_date, end_date
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (draft_id, service_id) DO UPDATE
		SET monday = EXCLUDED.monday, tuesday = EXCLUDED.tuesday, wednesday = EXCLUDED.wednesday,
		    thursday = EXCLUDED.thursday, friday = EXCLUDED.friday, saturday = EXCLUDED.saturday,
		    sunday = EXCLUDED.sunday, start_date = EXCLUDED.start_date, end_date = EXCLUDED.end_date,
		    updated_at = now()
	`, draft.ID, draft.AgencyID, calendar.ServiceID, calendar.Monday, calendar.Tuesday, calendar.Wednesday, calendar.Thursday, calendar.Friday, calendar.Saturday, calendar.Sunday, calendar.StartDate, calendar.EndDate)
	return err
}

func (s *DraftService) ListCalendars(ctx context.Context, draftID string) ([]DraftCalendar, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT draft_id, agency_id, service_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday, start_date, end_date
		FROM gtfs_draft_calendar WHERE draft_id = $1 ORDER BY service_id
	`, draftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var calendars []DraftCalendar
	for rows.Next() {
		var c DraftCalendar
		if err := rows.Scan(&c.DraftID, &c.AgencyID, &c.ServiceID, &c.Monday, &c.Tuesday, &c.Wednesday, &c.Thursday, &c.Friday, &c.Saturday, &c.Sunday, &c.StartDate, &c.EndDate); err != nil {
			return nil, err
		}
		calendars = append(calendars, c)
	}
	return calendars, rows.Err()
}

func (s *DraftService) RemoveCalendar(ctx context.Context, draftID string, serviceID string) error {
	return s.removeDraftRow(ctx, draftID, `DELETE FROM gtfs_draft_calendar WHERE draft_id = $1 AND service_id = $2`, serviceID)
}

func (s *DraftService) UpsertCalendarDate(ctx context.Context, date DraftCalendarDate) error {
	draft, err := s.ensureEditable(ctx, date.DraftID)
	if err != nil {
		return err
	}
	if err := validateDateText(date.Date); err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}
	if date.ExceptionType != 1 && date.ExceptionType != 2 {
		return fmt.Errorf("exception_type must be 1 or 2")
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO gtfs_draft_calendar_date (draft_id, agency_id, service_id, date, exception_type)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (draft_id, service_id, date) DO UPDATE
		SET exception_type = EXCLUDED.exception_type, updated_at = now()
	`, draft.ID, draft.AgencyID, date.ServiceID, date.Date, date.ExceptionType)
	return err
}

func (s *DraftService) ListCalendarDates(ctx context.Context, draftID string) ([]DraftCalendarDate, error) {
	rows, err := s.pool.Query(ctx, `SELECT draft_id, agency_id, service_id, date, exception_type FROM gtfs_draft_calendar_date WHERE draft_id = $1 ORDER BY service_id, date`, draftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dates []DraftCalendarDate
	for rows.Next() {
		var d DraftCalendarDate
		if err := rows.Scan(&d.DraftID, &d.AgencyID, &d.ServiceID, &d.Date, &d.ExceptionType); err != nil {
			return nil, err
		}
		dates = append(dates, d)
	}
	return dates, rows.Err()
}

func (s *DraftService) RemoveCalendarDate(ctx context.Context, draftID string, serviceID string, date string) error {
	return s.removeDraftRow(ctx, draftID, `DELETE FROM gtfs_draft_calendar_date WHERE draft_id = $1 AND service_id = $2 AND date = $3`, serviceID, date)
}

func (s *DraftService) UpsertTrip(ctx context.Context, trip DraftTrip) error {
	draft, err := s.ensureEditable(ctx, trip.DraftID)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO gtfs_draft_trip (draft_id, agency_id, id, route_id, service_id, block_id, shape_id, direction_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (draft_id, id) DO UPDATE
		SET route_id = EXCLUDED.route_id, service_id = EXCLUDED.service_id, block_id = EXCLUDED.block_id,
		    shape_id = EXCLUDED.shape_id, direction_id = EXCLUDED.direction_id, updated_at = now()
	`, draft.ID, draft.AgencyID, trip.ID, trip.RouteID, trip.ServiceID, nullString(trip.BlockID), nullString(trip.ShapeID), nullInt(trip.DirectionID))
	return err
}

func (s *DraftService) ListTrips(ctx context.Context, draftID string) ([]DraftTrip, error) {
	rows, err := s.pool.Query(ctx, `SELECT draft_id, agency_id, id, route_id, service_id, block_id, shape_id, direction_id FROM gtfs_draft_trip WHERE draft_id = $1 ORDER BY id`, draftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var trips []DraftTrip
	for rows.Next() {
		var trip DraftTrip
		var direction sql.NullInt64
		if err := rows.Scan(&trip.DraftID, &trip.AgencyID, &trip.ID, &trip.RouteID, &trip.ServiceID, nullStringScan(&trip.BlockID), nullStringScan(&trip.ShapeID), &direction); err != nil {
			return nil, err
		}
		if direction.Valid {
			value := int(direction.Int64)
			trip.DirectionID = &value
		}
		trips = append(trips, trip)
	}
	return trips, rows.Err()
}

func (s *DraftService) RemoveTrip(ctx context.Context, draftID string, id string) error {
	return s.removeDraftRow(ctx, draftID, `DELETE FROM gtfs_draft_trip WHERE draft_id = $1 AND id = $2`, id)
}

func (s *DraftService) UpsertStopTime(ctx context.Context, stopTime DraftStopTime) error {
	draft, err := s.ensureEditable(ctx, stopTime.DraftID)
	if err != nil {
		return err
	}
	if stopTime.ArrivalTime == "" && stopTime.DepartureTime == "" {
		return fmt.Errorf("arrival_time or departure_time is required")
	}
	if stopTime.ArrivalTime != "" {
		if _, err := ParseGTFSTime(stopTime.ArrivalTime); err != nil {
			return err
		}
	}
	if stopTime.DepartureTime != "" {
		if _, err := ParseGTFSTime(stopTime.DepartureTime); err != nil {
			return err
		}
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO gtfs_draft_stop_time (
			draft_id, agency_id, trip_id, arrival_time, departure_time, stop_id, stop_sequence, pickup_type, drop_off_type, shape_dist_traveled
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (draft_id, trip_id, stop_sequence) DO UPDATE
		SET arrival_time = EXCLUDED.arrival_time, departure_time = EXCLUDED.departure_time, stop_id = EXCLUDED.stop_id,
		    pickup_type = EXCLUDED.pickup_type, drop_off_type = EXCLUDED.drop_off_type,
		    shape_dist_traveled = EXCLUDED.shape_dist_traveled, updated_at = now()
	`, draft.ID, draft.AgencyID, stopTime.TripID, nullString(stopTime.ArrivalTime), nullString(stopTime.DepartureTime), stopTime.StopID, stopTime.StopSequence, nullInt(stopTime.PickupType), nullInt(stopTime.DropOffType), nullFloat(stopTime.ShapeDistTraveled))
	return err
}

func (s *DraftService) ListStopTimes(ctx context.Context, draftID string) ([]DraftStopTime, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT draft_id, agency_id, trip_id, arrival_time, departure_time, stop_id, stop_sequence, pickup_type, drop_off_type, shape_dist_traveled
		FROM gtfs_draft_stop_time WHERE draft_id = $1 ORDER BY trip_id, stop_sequence
	`, draftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stopTimes []DraftStopTime
	for rows.Next() {
		var st DraftStopTime
		if err := rows.Scan(&st.DraftID, &st.AgencyID, &st.TripID, nullStringScan(&st.ArrivalTime), nullStringScan(&st.DepartureTime), &st.StopID, &st.StopSequence, nullIntScan(&st.PickupType), nullIntScan(&st.DropOffType), nullFloatScan(&st.ShapeDistTraveled)); err != nil {
			return nil, err
		}
		stopTimes = append(stopTimes, st)
	}
	return stopTimes, rows.Err()
}

func (s *DraftService) RemoveStopTime(ctx context.Context, draftID string, tripID string, sequence int) error {
	return s.removeDraftRow(ctx, draftID, `DELETE FROM gtfs_draft_stop_time WHERE draft_id = $1 AND trip_id = $2 AND stop_sequence = $3`, tripID, sequence)
}

func (s *DraftService) UpsertShapePoint(ctx context.Context, point DraftShapePoint) error {
	draft, err := s.ensureEditable(ctx, point.DraftID)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO gtfs_draft_shape_point (draft_id, agency_id, shape_id, lat, lon, sequence, dist_traveled)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (draft_id, shape_id, sequence) DO UPDATE
		SET lat = EXCLUDED.lat, lon = EXCLUDED.lon, dist_traveled = EXCLUDED.dist_traveled, updated_at = now()
	`, draft.ID, draft.AgencyID, point.ShapeID, point.Lat, point.Lon, point.Sequence, nullFloat(point.DistTraveled))
	return err
}

func (s *DraftService) ListShapePoints(ctx context.Context, draftID string) ([]DraftShapePoint, error) {
	rows, err := s.pool.Query(ctx, `SELECT draft_id, agency_id, shape_id, lat, lon, sequence, dist_traveled FROM gtfs_draft_shape_point WHERE draft_id = $1 ORDER BY shape_id, sequence`, draftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var points []DraftShapePoint
	for rows.Next() {
		var point DraftShapePoint
		if err := rows.Scan(&point.DraftID, &point.AgencyID, &point.ShapeID, &point.Lat, &point.Lon, &point.Sequence, nullFloatScan(&point.DistTraveled)); err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	return points, rows.Err()
}

func (s *DraftService) RemoveShapePoint(ctx context.Context, draftID string, shapeID string, sequence int) error {
	return s.removeDraftRow(ctx, draftID, `DELETE FROM gtfs_draft_shape_point WHERE draft_id = $1 AND shape_id = $2 AND sequence = $3`, shapeID, sequence)
}

func (s *DraftService) UpsertFrequency(ctx context.Context, frequency DraftFrequency) error {
	draft, err := s.ensureEditable(ctx, frequency.DraftID)
	if err != nil {
		return err
	}
	if _, err := ParseGTFSTime(frequency.StartTime); err != nil {
		return err
	}
	if _, err := ParseGTFSTime(frequency.EndTime); err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO gtfs_draft_frequency (draft_id, agency_id, trip_id, start_time, end_time, headway_secs, exact_times)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (draft_id, trip_id, start_time) DO UPDATE
		SET end_time = EXCLUDED.end_time, headway_secs = EXCLUDED.headway_secs, exact_times = EXCLUDED.exact_times, updated_at = now()
	`, draft.ID, draft.AgencyID, frequency.TripID, frequency.StartTime, frequency.EndTime, frequency.HeadwaySecs, frequency.ExactTimes)
	return err
}

func (s *DraftService) ListFrequencies(ctx context.Context, draftID string) ([]DraftFrequency, error) {
	rows, err := s.pool.Query(ctx, `SELECT draft_id, agency_id, trip_id, start_time, end_time, headway_secs, exact_times FROM gtfs_draft_frequency WHERE draft_id = $1 ORDER BY trip_id, start_time`, draftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var frequencies []DraftFrequency
	for rows.Next() {
		var f DraftFrequency
		if err := rows.Scan(&f.DraftID, &f.AgencyID, &f.TripID, &f.StartTime, &f.EndTime, &f.HeadwaySecs, &f.ExactTimes); err != nil {
			return nil, err
		}
		frequencies = append(frequencies, f)
	}
	return frequencies, rows.Err()
}

func (s *DraftService) RemoveFrequency(ctx context.Context, draftID string, tripID string, startTime string) error {
	return s.removeDraftRow(ctx, draftID, `DELETE FROM gtfs_draft_frequency WHERE draft_id = $1 AND trip_id = $2 AND start_time = $3`, tripID, startTime)
}

func (s *DraftService) ensureEditable(ctx context.Context, draftID string) (Draft, error) {
	draft, err := s.GetDraft(ctx, draftID)
	if err != nil {
		return Draft{}, err
	}
	if draft.Status != DraftStatusDraft {
		return Draft{}, fmt.Errorf("%w: status %s", ErrDraftNotEditable, draft.Status)
	}
	return draft, nil
}

func (s *DraftService) removeDraftRow(ctx context.Context, draftID string, query string, args ...any) error {
	if _, err := s.ensureEditable(ctx, draftID); err != nil {
		return err
	}
	allArgs := append([]any{draftID}, args...)
	_, err := s.pool.Exec(ctx, query, allArgs...)
	return err
}

func (s *DraftService) insertDraftPublishAttempt(ctx context.Context, draft Draft, opts PublishDraftOptions, report ImportReport) (int64, error) {
	reportJSON, err := json.Marshal(report)
	if err != nil {
		return 0, fmt.Errorf("marshal draft publish report: %w", err)
	}
	var id int64
	err = s.pool.QueryRow(ctx, `
		INSERT INTO gtfs_draft_publish (
			draft_id, agency_id, status, error_count, warning_count, info_count, report_json, actor_id, notes, started_at
		)
		VALUES ($1, $2, 'started', $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, draft.ID, draft.AgencyID, len(report.Errors), len(report.Warnings), len(report.Info), reportJSON, nullString(opts.ActorID), nullString(opts.Notes), s.now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE gtfs_draft
		SET last_publish_attempt_id = $2,
		    updated_at = $3
		WHERE id = $1
	`, draft.ID, id, s.now())
	return id, err
}

func (s *DraftService) markDraftPublishFailed(ctx context.Context, draftID string, publishID int64, report ImportReport) error {
	report.Status = ImportStatusFailed
	reportJSON, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("marshal failed draft publish report: %w", err)
	}
	now := s.now()
	if _, err := s.pool.Exec(ctx, `
		UPDATE gtfs_draft_publish
		SET status = 'failed',
		    feed_version_id = NULL,
		    error_count = $2,
		    warning_count = $3,
		    info_count = $4,
		    report_json = $5,
		    completed_at = $6
		WHERE id = $1
	`, publishID, len(report.Errors), len(report.Warnings), len(report.Info), reportJSON, now); err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `UPDATE gtfs_draft SET last_publish_attempt_id = $2, updated_at = $3 WHERE id = $1`, draftID, publishID, now)
	return err
}

func (s *DraftService) draftFeed(ctx context.Context, draft Draft) (parsedFeed, ImportReport, error) {
	report := ImportReport{
		Status: ImportStatusStarted,
		Counts: map[string]int{},
	}
	agency, err := s.GetAgency(ctx, draft.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			report.addError("agency.txt", 0, "missing_required_file", "draft agency metadata is required")
		} else {
			return parsedFeed{}, report, err
		}
	}
	routes, err := s.ListRoutes(ctx, draft.ID)
	if err != nil {
		return parsedFeed{}, report, err
	}
	stops, err := s.ListStops(ctx, draft.ID)
	if err != nil {
		return parsedFeed{}, report, err
	}
	calendars, err := s.ListCalendars(ctx, draft.ID)
	if err != nil {
		return parsedFeed{}, report, err
	}
	calendarDates, err := s.ListCalendarDates(ctx, draft.ID)
	if err != nil {
		return parsedFeed{}, report, err
	}
	trips, err := s.ListTrips(ctx, draft.ID)
	if err != nil {
		return parsedFeed{}, report, err
	}
	stopTimes, err := s.ListStopTimes(ctx, draft.ID)
	if err != nil {
		return parsedFeed{}, report, err
	}
	shapePoints, err := s.ListShapePoints(ctx, draft.ID)
	if err != nil {
		return parsedFeed{}, report, err
	}
	frequencies, err := s.ListFrequencies(ctx, draft.ID)
	if err != nil {
		return parsedFeed{}, report, err
	}

	feed := parsedFeed{
		AgencyID:           draft.AgencyID,
		ShapePointsByShape: map[string][]importShapePoint{},
	}
	if agency.AgencyID != "" {
		feed.Agencies = []importAgency{{
			ID:       agency.AgencyID,
			Name:     agency.Name,
			URL:      agency.PublicURL,
			Timezone: agency.Timezone,
			Email:    agency.ContactEmail,
		}}
		if _, err := time.LoadLocation(agency.Timezone); err != nil {
			report.addError("agency.txt", 0, "invalid_timezone", fmt.Sprintf("agency_timezone %q is invalid", agency.Timezone))
		}
	}
	for _, r := range routes {
		if !supportedRouteType(r.RouteType) {
			report.addError("routes.txt", 0, "invalid_route_type", "route_type must be 0-7 or 100-1702")
		}
		feed.Routes = append(feed.Routes, importRoute{ID: r.ID, AgencyID: draft.AgencyID, ShortName: r.ShortName, LongName: r.LongName, RouteType: r.RouteType})
	}
	for _, s := range stops {
		feed.Stops = append(feed.Stops, importStop{ID: s.ID, Name: s.Name, Lat: s.Lat, Lon: s.Lon})
	}
	for _, c := range calendars {
		feed.Calendars = append(feed.Calendars, importCalendar{
			ServiceID: c.ServiceID, Monday: c.Monday, Tuesday: c.Tuesday, Wednesday: c.Wednesday,
			Thursday: c.Thursday, Friday: c.Friday, Saturday: c.Saturday, Sunday: c.Sunday,
			StartDate: c.StartDate, EndDate: c.EndDate,
		})
	}
	for _, d := range calendarDates {
		feed.CalendarDates = append(feed.CalendarDates, importCalendarDate{ServiceID: d.ServiceID, Date: d.Date, ExceptionType: d.ExceptionType})
	}
	for _, t := range trips {
		feed.Trips = append(feed.Trips, importTrip{
			ID: t.ID, RouteID: t.RouteID, ServiceID: t.ServiceID, BlockID: t.BlockID, ShapeID: t.ShapeID, DirectionID: nullInt(t.DirectionID),
		})
	}
	for _, st := range stopTimes {
		if st.ArrivalTime != "" {
			if _, err := ParseGTFSTime(st.ArrivalTime); err != nil {
				report.addError("stop_times.txt", 0, "invalid_gtfs_time", err.Error())
			}
		}
		if st.DepartureTime != "" {
			if _, err := ParseGTFSTime(st.DepartureTime); err != nil {
				report.addError("stop_times.txt", 0, "invalid_gtfs_time", err.Error())
			}
		}
		feed.StopTimes = append(feed.StopTimes, importStopTime{
			TripID: st.TripID, ArrivalTime: st.ArrivalTime, DepartureTime: st.DepartureTime,
			StopID: st.StopID, StopSequence: st.StopSequence, PickupType: nullInt(st.PickupType),
			DropOffType: nullInt(st.DropOffType), ShapeDistTraveled: nullFloat(st.ShapeDistTraveled),
		})
	}
	for _, p := range shapePoints {
		feed.ShapePoints = append(feed.ShapePoints, importShapePoint{
			ShapeID: p.ShapeID, Lat: p.Lat, Lon: p.Lon, Sequence: p.Sequence, DistTraveled: nullFloat(p.DistTraveled),
		})
	}
	for _, f := range frequencies {
		if _, err := ParseGTFSTime(f.StartTime); err != nil {
			report.addError("frequencies.txt", 0, "invalid_gtfs_time", err.Error())
		}
		if _, err := ParseGTFSTime(f.EndTime); err != nil {
			report.addError("frequencies.txt", 0, "invalid_gtfs_time", err.Error())
		}
		feed.Frequencies = append(feed.Frequencies, importFrequency{TripID: f.TripID, StartTime: f.StartTime, EndTime: f.EndTime, HeadwaySecs: f.HeadwaySecs, ExactTimes: f.ExactTimes})
	}

	report.Counts["agency"] = len(feed.Agencies)
	report.Counts["routes"] = len(feed.Routes)
	report.Counts["stops"] = len(feed.Stops)
	report.Counts["calendar"] = len(feed.Calendars)
	report.Counts["calendar_dates"] = len(feed.CalendarDates)
	report.Counts["trips"] = len(feed.Trips)
	report.Counts["stop_times"] = len(feed.StopTimes)
	report.Counts["shapes"] = len(feed.ShapePoints)
	report.Counts["frequencies"] = len(feed.Frequencies)

	if len(feed.Agencies) == 0 {
		report.addError("agency.txt", 0, "missing_required_file", "agency metadata is required")
	}
	if len(feed.Routes) == 0 {
		report.addError("routes.txt", 0, "missing_required_file", "routes are required")
	}
	if len(feed.Stops) == 0 {
		report.addError("stops.txt", 0, "missing_required_file", "stops are required")
	}
	if len(feed.Trips) == 0 {
		report.addError("trips.txt", 0, "missing_required_file", "trips are required")
	}
	if len(feed.StopTimes) == 0 {
		report.addError("stop_times.txt", 0, "missing_required_file", "stop_times are required")
	}
	if len(feed.Calendars) == 0 && len(feed.CalendarDates) == 0 {
		report.addError("", 0, "missing_service_source", "at least one usable service source from calendars or calendar_dates is required")
	}
	validateFeed(&feed, &report)
	report.Status = ImportStatusPublished
	if report.hasErrors() {
		report.Status = ImportStatusFailed
	}
	return feed, report, nil
}

func activeFeedVersionID(ctx context.Context, tx pgx.Tx, agencyID string) (string, error) {
	var id string
	err := tx.QueryRow(ctx, `
		SELECT id
		FROM feed_version
		WHERE agency_id = $1
		  AND is_active
		ORDER BY activated_at DESC NULLS LAST, created_at DESC
		LIMIT 1
	`, agencyID).Scan(&id)
	return id, err
}

func clonePublishedFeedToDraft(ctx context.Context, tx pgx.Tx, agencyID string, feedVersionID string, draftID string) error {
	if _, err := tx.Exec(ctx, `
		INSERT INTO gtfs_draft_agency (draft_id, agency_id, name, timezone, contact_email, public_url)
		SELECT $1, id, name, timezone, contact_email, public_url
		FROM agency
		WHERE id = $2
	`, draftID, agencyID); err != nil {
		return fmt.Errorf("clone draft agency: %w", err)
	}
	statements := []string{
		`INSERT INTO gtfs_draft_route (draft_id, agency_id, id, short_name, long_name, route_type)
		 SELECT $1, agency_id, id, short_name, long_name, route_type FROM gtfs_route WHERE agency_id = $2 AND feed_version_id = $3`,
		`INSERT INTO gtfs_draft_stop (draft_id, agency_id, id, name, lat, lon)
		 SELECT $1, agency_id, id, name, lat, lon FROM gtfs_stop WHERE agency_id = $2 AND feed_version_id = $3`,
		`INSERT INTO gtfs_draft_calendar (draft_id, agency_id, service_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday, start_date, end_date)
		 SELECT $1, agency_id, service_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday, start_date, end_date FROM gtfs_calendar WHERE agency_id = $2 AND feed_version_id = $3`,
		`INSERT INTO gtfs_draft_calendar_date (draft_id, agency_id, service_id, date, exception_type)
		 SELECT $1, agency_id, service_id, date, exception_type FROM gtfs_calendar_date WHERE agency_id = $2 AND feed_version_id = $3`,
		`INSERT INTO gtfs_draft_trip (draft_id, agency_id, id, route_id, service_id, block_id, shape_id, direction_id)
		 SELECT $1, agency_id, id, route_id, service_id, block_id, shape_id, direction_id FROM gtfs_trip WHERE agency_id = $2 AND feed_version_id = $3`,
		`INSERT INTO gtfs_draft_stop_time (draft_id, agency_id, trip_id, arrival_time, departure_time, stop_id, stop_sequence, pickup_type, drop_off_type, shape_dist_traveled)
		 SELECT $1, agency_id, trip_id, arrival_time, departure_time, stop_id, stop_sequence, pickup_type, drop_off_type, shape_dist_traveled FROM gtfs_stop_time WHERE agency_id = $2 AND feed_version_id = $3`,
		`INSERT INTO gtfs_draft_shape_point (draft_id, agency_id, shape_id, lat, lon, sequence, dist_traveled)
		 SELECT $1, agency_id, shape_id, lat, lon, sequence, dist_traveled FROM gtfs_shape_point WHERE agency_id = $2 AND feed_version_id = $3`,
		`INSERT INTO gtfs_draft_frequency (draft_id, agency_id, trip_id, start_time, end_time, headway_secs, exact_times)
		 SELECT $1, agency_id, trip_id, start_time, end_time, headway_secs, exact_times FROM gtfs_frequency WHERE agency_id = $2 AND feed_version_id = $3`,
	}
	for _, stmt := range statements {
		if _, err := tx.Exec(ctx, stmt, draftID, agencyID, feedVersionID); err != nil {
			return fmt.Errorf("clone published feed to draft: %w", err)
		}
	}
	return nil
}

func insertBlankDraftAgency(ctx context.Context, tx pgx.Tx, agencyID string, draftID string) error {
	tag, err := tx.Exec(ctx, `
		INSERT INTO gtfs_draft_agency (draft_id, agency_id, name, timezone, contact_email, public_url)
		SELECT $1, id, name, timezone, contact_email, public_url
		FROM agency
		WHERE id = $2
	`, draftID, agencyID)
	if err != nil {
		return fmt.Errorf("insert blank draft agency: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("agency %q does not exist", agencyID)
	}
	return nil
}

func publishResult(draft Draft, publishID int64, feedVersionID string, status string, report ImportReport, stored bool, failure string) PublishDraftResult {
	return PublishDraftResult{
		PublishID:      publishID,
		DraftID:        draft.ID,
		AgencyID:       draft.AgencyID,
		FeedVersionID:  feedVersionID,
		Status:         status,
		ErrorCount:     len(report.Errors),
		WarningCount:   len(report.Warnings),
		InfoCount:      len(report.Info),
		Counts:         report.Counts,
		ReportStored:   stored,
		FailureMessage: failure,
	}
}

func scanDraft(rows pgx.Rows, draft *Draft, latestStatus *string, latestID *int64) error {
	return rows.Scan(
		&draft.ID, &draft.AgencyID, &draft.Name, &draft.Status,
		nullStringScan(&draft.BaseFeedVersionID),
		nullStringScan(&draft.LastPublishedFeedVersionID),
		nullInt64Scan(&draft.LastPublishAttemptID),
		&draft.DiscardedAt,
		nullStringScan(&draft.DiscardedBy),
		nullStringScan(&draft.DiscardReason),
		nullStringScan(&draft.CreatedBy),
		&draft.CreatedAt,
		&draft.UpdatedAt,
		latestStatus,
		latestID,
	)
}

func randomHex(bytesLen int) string {
	buf := make([]byte, bytesLen)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

func validateDateText(value string) error {
	if len(value) != 8 {
		return fmt.Errorf("date must use YYYYMMDD")
	}
	if _, err := time.Parse("20060102", value); err != nil {
		return err
	}
	return nil
}

type nullableString struct {
	dest *string
}

func nullStringScan(dest *string) *nullableString {
	return &nullableString{dest: dest}
}

func (n *nullableString) Scan(src any) error {
	var value sql.NullString
	if err := value.Scan(src); err != nil {
		return err
	}
	if value.Valid {
		*n.dest = value.String
	} else {
		*n.dest = ""
	}
	return nil
}

type nullableInt64 struct {
	dest *int64
}

func nullInt64Scan(dest *int64) *nullableInt64 {
	return &nullableInt64{dest: dest}
}

func (n *nullableInt64) Scan(src any) error {
	var value sql.NullInt64
	if err := value.Scan(src); err != nil {
		return err
	}
	if value.Valid {
		*n.dest = value.Int64
	} else {
		*n.dest = 0
	}
	return nil
}

type nullableInt struct {
	dest **int
}

func nullIntScan(dest **int) *nullableInt {
	return &nullableInt{dest: dest}
}

func (n *nullableInt) Scan(src any) error {
	var value sql.NullInt64
	if err := value.Scan(src); err != nil {
		return err
	}
	if value.Valid {
		converted := int(value.Int64)
		*n.dest = &converted
	} else {
		*n.dest = nil
	}
	return nil
}

type nullableFloat struct {
	dest **float64
}

func nullFloatScan(dest **float64) *nullableFloat {
	return &nullableFloat{dest: dest}
}

func (n *nullableFloat) Scan(src any) error {
	var value sql.NullFloat64
	if err := value.Scan(src); err != nil {
		return err
	}
	if value.Valid {
		converted := value.Float64
		*n.dest = &converted
	} else {
		*n.dest = nil
	}
	return nil
}
