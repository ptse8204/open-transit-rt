package gtfs

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ImportStatusStarted   = "started"
	ImportStatusFailed    = "failed"
	ImportStatusPublished = "published"
)

type ImportService struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

type validationReportLinks struct {
	ImportID       int64
	DraftPublishID int64
}

type publishFeedOptions struct {
	AgencyID      string
	FeedVersionID string
	SourceType    string
	ActorID       string
	Notes         string
	AuditAction   string
	EntityType    string
	EntityID      string
	Feed          parsedFeed
	Report        ImportReport
	Links         validationReportLinks
	AfterPublish  func(ctx context.Context, tx pgx.Tx, publishedAt time.Time, report ImportReport, reportJSON []byte) error
}

type ImportOptions struct {
	AgencyID string
	ZipPath  string
	ActorID  string
	Notes    string
}

type ImportResult struct {
	ImportID       int64          `json:"import_id,omitempty"`
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

type ImportReport struct {
	Status   string            `json:"status"`
	Errors   []ImportMessage   `json:"errors,omitempty"`
	Warnings []ImportMessage   `json:"warnings,omitempty"`
	Info     []ImportMessage   `json:"info,omitempty"`
	Counts   map[string]int    `json:"counts,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type ImportMessage struct {
	File    string `json:"file,omitempty"`
	Row     int    `json:"row,omitempty"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ImportError struct {
	Result ImportResult
	Err    error
}

func (e *ImportError) Error() string {
	return e.Err.Error()
}

func (e *ImportError) Unwrap() error {
	return e.Err
}

func NewImportService(pool *pgxpool.Pool) *ImportService {
	return &ImportService{
		pool: pool,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *ImportService) ImportZip(ctx context.Context, opts ImportOptions) (ImportResult, error) {
	if s == nil || s.pool == nil {
		return ImportResult{}, fmt.Errorf("gtfs import service requires a database pool")
	}
	opts.AgencyID = strings.TrimSpace(opts.AgencyID)
	opts.ZipPath = strings.TrimSpace(opts.ZipPath)
	if opts.AgencyID == "" {
		return ImportResult{}, fmt.Errorf("agency_id is required")
	}
	if opts.ZipPath == "" {
		return ImportResult{}, fmt.Errorf("zip path is required")
	}

	payload, err := os.ReadFile(opts.ZipPath)
	if err != nil {
		return ImportResult{}, fmt.Errorf("read gtfs zip: %w", err)
	}
	source := importSource{
		Filename: filepath.Base(opts.ZipPath),
		SHA256:   sha256Hex(payload),
		ByteSize: int64(len(payload)),
	}

	feed, report := parseGTFSZip(payload, opts.AgencyID)
	report.Metadata = map[string]string{
		"source_filename": source.Filename,
		"source_sha256":   source.SHA256,
	}

	agency, agencyOK := feed.selectedAgency()
	if agencyOK {
		if err := s.upsertAgency(ctx, agency); err != nil {
			result := reportResult(opts.AgencyID, 0, "", ImportStatusFailed, report, false, "")
			return result, &ImportError{Result: result, Err: fmt.Errorf("upsert agency before import report: %w", err)}
		}
	} else if !s.agencyExists(ctx, opts.AgencyID) {
		result := reportResult(opts.AgencyID, 0, "", ImportStatusFailed, report, false, "agency metadata could not be stored before report")
		return result, &ImportError{Result: result, Err: fmt.Errorf("gtfs import failed and failure report could not be stored: agency %q does not exist and agency.txt did not contain usable matching metadata", opts.AgencyID)}
	}

	importID, err := s.insertImportAttempt(ctx, opts, source, report)
	if err != nil {
		result := reportResult(opts.AgencyID, 0, "", ImportStatusFailed, report, false, "failed to store import attempt")
		return result, &ImportError{Result: result, Err: fmt.Errorf("store gtfs import attempt: %w", err)}
	}

	if report.hasErrors() {
		if err := s.markImportFailed(ctx, importID, report); err != nil {
			result := reportResult(opts.AgencyID, importID, "", ImportStatusFailed, report, false, "failed to store validation failure report")
			return result, &ImportError{Result: result, Err: fmt.Errorf("gtfs import validation failed and failure report could not be stored: %w", err)}
		}
		if err := s.insertValidationReport(ctx, nil, importID, opts.AgencyID, "", report); err != nil {
			result := reportResult(opts.AgencyID, importID, "", ImportStatusFailed, report, false, "failed to store validation report")
			return result, &ImportError{Result: result, Err: fmt.Errorf("gtfs import validation failed and validation report could not be stored: %w", err)}
		}
		result := reportResult(opts.AgencyID, importID, "", ImportStatusFailed, report, true, "validation failed")
		return result, &ImportError{Result: result, Err: fmt.Errorf("gtfs import validation failed with %d error(s)", len(report.Errors))}
	}

	feedVersionID := fmt.Sprintf("gtfs-import-%d", importID)
	if err := s.publish(ctx, opts, importID, feedVersionID, feed, report); err != nil {
		failedReport := reportWithError(report, ImportMessage{
			Code:    "publish_failed",
			Message: err.Error(),
		})
		if markErr := s.markImportFailed(ctx, importID, failedReport); markErr != nil {
			result := reportResult(opts.AgencyID, importID, "", ImportStatusFailed, failedReport, false, "publish failed and failure report could not be stored")
			return result, &ImportError{Result: result, Err: fmt.Errorf("gtfs import publish failed and failure report could not be stored: publish error: %v; report error: %w", err, markErr)}
		}
		if reportErr := s.insertValidationReport(ctx, nil, importID, opts.AgencyID, "", failedReport); reportErr != nil {
			result := reportResult(opts.AgencyID, importID, "", ImportStatusFailed, failedReport, false, "publish failed and validation report could not be stored")
			return result, &ImportError{Result: result, Err: fmt.Errorf("gtfs import publish failed and validation report could not be stored: publish error: %v; report error: %w", err, reportErr)}
		}
		result := reportResult(opts.AgencyID, importID, "", ImportStatusFailed, failedReport, true, "publish failed")
		return result, &ImportError{Result: result, Err: fmt.Errorf("gtfs import publish failed: %w", err)}
	}

	return reportResult(opts.AgencyID, importID, feedVersionID, ImportStatusPublished, report, true, ""), nil
}

type importSource struct {
	Filename string
	SHA256   string
	ByteSize int64
}

func sha256Hex(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func reportResult(agencyID string, importID int64, feedVersionID string, status string, report ImportReport, stored bool, failure string) ImportResult {
	return ImportResult{
		ImportID:       importID,
		AgencyID:       agencyID,
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

func reportWithError(report ImportReport, msg ImportMessage) ImportReport {
	report.Errors = append(report.Errors, msg)
	report.Status = ImportStatusFailed
	return report
}

func (r ImportReport) hasErrors() bool {
	return len(r.Errors) > 0
}

func (r ImportReport) validationStatus() string {
	if len(r.Errors) > 0 {
		return "failed"
	}
	if len(r.Warnings) > 0 {
		return "warning"
	}
	return "passed"
}

func (s *ImportService) upsertAgency(ctx context.Context, agency importAgency) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO agency (id, name, timezone, contact_email, public_url)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    timezone = EXCLUDED.timezone,
		    contact_email = COALESCE(EXCLUDED.contact_email, agency.contact_email),
		    public_url = EXCLUDED.public_url,
		    updated_at = now()
	`, agency.ID, agency.Name, agency.Timezone, nullString(agency.Email), agency.URL)
	return err
}

func upsertAgencyTx(ctx context.Context, tx pgx.Tx, agency importAgency) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO agency (id, name, timezone, contact_email, public_url)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    timezone = EXCLUDED.timezone,
		    contact_email = COALESCE(EXCLUDED.contact_email, agency.contact_email),
		    public_url = EXCLUDED.public_url,
		    updated_at = now()
	`, agency.ID, agency.Name, agency.Timezone, nullString(agency.Email), agency.URL)
	if err != nil {
		return fmt.Errorf("upsert agency in publish transaction: %w", err)
	}
	return nil
}

func (s *ImportService) agencyExists(ctx context.Context, agencyID string) bool {
	var exists bool
	err := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM agency WHERE id = $1)`, agencyID).Scan(&exists)
	return err == nil && exists
}

func (s *ImportService) insertImportAttempt(ctx context.Context, opts ImportOptions, source importSource, report ImportReport) (int64, error) {
	reportJSON, err := json.Marshal(report)
	if err != nil {
		return 0, fmt.Errorf("marshal import report: %w", err)
	}
	var id int64
	err = s.pool.QueryRow(ctx, `
		INSERT INTO gtfs_import (
			agency_id, source_filename, source_sha256, source_byte_size, status,
			error_count, warning_count, info_count, report_json, actor_id, notes, started_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`, opts.AgencyID, source.Filename, source.SHA256, source.ByteSize, ImportStatusStarted,
		len(report.Errors), len(report.Warnings), len(report.Info), reportJSON, nullString(opts.ActorID), nullString(opts.Notes), s.now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ImportService) markImportFailed(ctx context.Context, importID int64, report ImportReport) error {
	report.Status = ImportStatusFailed
	reportJSON, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("marshal failed import report: %w", err)
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE gtfs_import
		SET status = $2,
		    feed_version_id = NULL,
		    error_count = $3,
		    warning_count = $4,
		    info_count = $5,
		    report_json = $6,
		    completed_at = $7
		WHERE id = $1
	`, importID, ImportStatusFailed, len(report.Errors), len(report.Warnings), len(report.Info), reportJSON, s.now())
	return err
}

func (s *ImportService) publish(ctx context.Context, opts ImportOptions, importID int64, feedVersionID string, feed parsedFeed, report ImportReport) error {
	return publishFeed(ctx, s.pool, s.now, publishFeedOptions{
		AgencyID:      opts.AgencyID,
		FeedVersionID: feedVersionID,
		SourceType:    "gtfs_import",
		ActorID:       opts.ActorID,
		Notes:         opts.Notes,
		AuditAction:   "gtfs_import_publish",
		EntityType:    "feed_version",
		EntityID:      feedVersionID,
		Feed:          feed,
		Report:        report,
		Links:         validationReportLinks{ImportID: importID},
		AfterPublish: func(ctx context.Context, tx pgx.Tx, publishedAt time.Time, report ImportReport, reportJSON []byte) error {
			if _, err := tx.Exec(ctx, `
				UPDATE gtfs_import
				SET status = $2,
				    feed_version_id = $3,
				    error_count = $4,
				    warning_count = $5,
				    info_count = $6,
				    report_json = $7,
				    completed_at = $8
				WHERE id = $1
			`, importID, ImportStatusPublished, feedVersionID, len(report.Errors), len(report.Warnings), len(report.Info), reportJSON, publishedAt); err != nil {
				return fmt.Errorf("mark import published: %w", err)
			}
			return nil
		},
	})
}

func publishFeed(ctx context.Context, pool *pgxpool.Pool, nowFunc func() time.Time, opts publishFeedOptions) error {
	if pool == nil {
		return fmt.Errorf("publish feed requires a database pool")
	}
	if nowFunc == nil {
		nowFunc = func() time.Time { return time.Now().UTC() }
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if agency, ok := opts.Feed.selectedAgency(); ok {
		if err := upsertAgencyTx(ctx, tx, agency); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO feed_version (id, agency_id, source_type, lifecycle_state, is_active, validation_status, notes, created_at)
		VALUES ($1, $2, $3, 'staged', false, $4, $5, $6)
	`, opts.FeedVersionID, opts.AgencyID, opts.SourceType, opts.Report.validationStatus(), nullString(opts.Notes), nowFunc()); err != nil {
		return fmt.Errorf("insert staged feed version: %w", err)
	}

	if err := insertFeedRows(ctx, tx, opts.AgencyID, opts.FeedVersionID, opts.Feed); err != nil {
		return err
	}

	if err := insertValidationReport(ctx, tx, opts.AgencyID, opts.FeedVersionID, opts.Report, opts.Links); err != nil {
		return err
	}

	now := nowFunc()
	if _, err := tx.Exec(ctx, `
		UPDATE feed_version
		SET is_active = false,
		    lifecycle_state = 'retired',
		    retired_at = $3
		WHERE agency_id = $1
		  AND is_active
		  AND id <> $2
	`, opts.AgencyID, opts.FeedVersionID, now); err != nil {
		return fmt.Errorf("retire prior active feed version: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE feed_version
		SET is_active = true,
		    lifecycle_state = 'active',
		    published_at = $3,
		    activated_at = $3
		WHERE agency_id = $1
		  AND id = $2
	`, opts.AgencyID, opts.FeedVersionID, now); err != nil {
		return fmt.Errorf("activate feed version: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE published_feed
		SET active_feed_version_id = $2,
		    activation_status = 'active',
		    revision_timestamp = $3,
		    updated_at = $3
		WHERE agency_id = $1
		  AND feed_type = 'schedule'
	`, opts.AgencyID, opts.FeedVersionID, now); err != nil {
		return fmt.Errorf("update published feed metadata: %w", err)
	}

	opts.Report.Status = ImportStatusPublished
	reportJSON, err := json.Marshal(opts.Report)
	if err != nil {
		return fmt.Errorf("marshal published report: %w", err)
	}
	if opts.AfterPublish != nil {
		if err := opts.AfterPublish(ctx, tx, now, opts.Report, reportJSON); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (agency_id, actor_id, action, entity_type, entity_id, new_value_json, reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, opts.AgencyID, defaultString(opts.ActorID, "system"), opts.AuditAction, opts.EntityType, opts.EntityID, reportJSON, nullString(opts.Notes)); err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit feed publish: %w", err)
	}
	return nil
}

func (s *ImportService) insertValidationReport(ctx context.Context, tx pgx.Tx, importID int64, agencyID string, feedVersionID string, report ImportReport) error {
	return insertValidationReport(ctx, txOrPool{s.pool, tx}, agencyID, feedVersionID, report, validationReportLinks{ImportID: importID})
}

type txOrPool struct {
	pool *pgxpool.Pool
	tx   pgx.Tx
}

func (e txOrPool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	if e.tx != nil {
		return e.tx.Exec(ctx, sql, arguments...)
	}
	return e.pool.Exec(ctx, sql, arguments...)
}

func insertValidationReport(ctx context.Context, exec interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}, agencyID string, feedVersionID string, report ImportReport, links validationReportLinks) error {
	reportJSON, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("marshal validation report: %w", err)
	}
	args := []any{
		agencyID,
		nullString(feedVersionID),
		"schedule",
		"open-transit-rt-internal-gtfs-import",
		"phase-4",
		report.validationStatus(),
		len(report.Errors),
		len(report.Warnings),
		len(report.Info),
		reportJSON,
		nullInt64(links.ImportID),
		nullInt64(links.DraftPublishID),
	}
	query := `
		INSERT INTO validation_report (
			agency_id, feed_version_id, feed_type, validator_name, validator_version, status,
			error_count, warning_count, info_count, report_json, gtfs_import_id, gtfs_draft_publish_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err = exec.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert validation report: %w", err)
	}
	return nil
}

func insertFeedRows(ctx context.Context, tx pgx.Tx, agencyID string, feedVersionID string, feed parsedFeed) error {
	for _, route := range feed.Routes {
		if _, err := tx.Exec(ctx, `
			INSERT INTO gtfs_route (id, feed_version_id, agency_id, short_name, long_name, route_type)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, route.ID, feedVersionID, agencyID, nullString(route.ShortName), nullString(route.LongName), route.RouteType); err != nil {
			return fmt.Errorf("insert route %s: %w", route.ID, err)
		}
	}
	for _, stop := range feed.Stops {
		if _, err := tx.Exec(ctx, `
			INSERT INTO gtfs_stop (id, feed_version_id, agency_id, name, lat, lon, geom)
			VALUES ($1, $2, $3, $4, $5, $6, ST_SetSRID(ST_MakePoint($6, $5), 4326))
		`, stop.ID, feedVersionID, agencyID, stop.Name, stop.Lat, stop.Lon); err != nil {
			return fmt.Errorf("insert stop %s: %w", stop.ID, err)
		}
	}
	for _, calendar := range feed.Calendars {
		if _, err := tx.Exec(ctx, `
			INSERT INTO gtfs_calendar (
				service_id, feed_version_id, agency_id, monday, tuesday, wednesday, thursday,
				friday, saturday, sunday, start_date, end_date
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, calendar.ServiceID, feedVersionID, agencyID, calendar.Monday, calendar.Tuesday, calendar.Wednesday,
			calendar.Thursday, calendar.Friday, calendar.Saturday, calendar.Sunday, calendar.StartDate, calendar.EndDate); err != nil {
			return fmt.Errorf("insert calendar %s: %w", calendar.ServiceID, err)
		}
	}
	for _, date := range feed.CalendarDates {
		if _, err := tx.Exec(ctx, `
			INSERT INTO gtfs_calendar_date (service_id, feed_version_id, agency_id, date, exception_type)
			VALUES ($1, $2, $3, $4, $5)
		`, date.ServiceID, feedVersionID, agencyID, date.Date, date.ExceptionType); err != nil {
			return fmt.Errorf("insert calendar date %s/%s: %w", date.ServiceID, date.Date, err)
		}
	}
	for _, trip := range feed.Trips {
		if _, err := tx.Exec(ctx, `
			INSERT INTO gtfs_trip (id, feed_version_id, agency_id, route_id, service_id, block_id, shape_id, direction_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, trip.ID, feedVersionID, agencyID, trip.RouteID, trip.ServiceID, nullString(trip.BlockID), nullString(trip.ShapeID), trip.DirectionID); err != nil {
			return fmt.Errorf("insert trip %s: %w", trip.ID, err)
		}
	}
	for _, stopTime := range feed.StopTimes {
		if _, err := tx.Exec(ctx, `
			INSERT INTO gtfs_stop_time (
				trip_id, feed_version_id, agency_id, arrival_time, departure_time, stop_id,
				stop_sequence, pickup_type, drop_off_type, shape_dist_traveled
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, stopTime.TripID, feedVersionID, agencyID, nullString(stopTime.ArrivalTime), nullString(stopTime.DepartureTime),
			stopTime.StopID, stopTime.StopSequence, stopTime.PickupType, stopTime.DropOffType, stopTime.ShapeDistTraveled); err != nil {
			return fmt.Errorf("insert stop time %s/%d: %w", stopTime.TripID, stopTime.StopSequence, err)
		}
	}
	for _, point := range feed.ShapePoints {
		if _, err := tx.Exec(ctx, `
			INSERT INTO gtfs_shape_point (shape_id, feed_version_id, agency_id, lat, lon, sequence, dist_traveled, geom)
			VALUES ($1, $2, $3, $4, $5, $6, $7, ST_SetSRID(ST_MakePoint($5, $4), 4326))
		`, point.ShapeID, feedVersionID, agencyID, point.Lat, point.Lon, point.Sequence, point.DistTraveled); err != nil {
			return fmt.Errorf("insert shape point %s/%d: %w", point.ShapeID, point.Sequence, err)
		}
	}
	for shapeID, points := range feed.ShapePointsByShape {
		if len(points) < 2 {
			continue
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO gtfs_shape_line (shape_id, feed_version_id, agency_id, geom)
			SELECT $1, $2, $3, ST_MakeLine(geom ORDER BY sequence)
			FROM gtfs_shape_point
			WHERE agency_id = $3
			  AND feed_version_id = $2
			  AND shape_id = $1
			GROUP BY shape_id
		`, shapeID, feedVersionID, agencyID); err != nil {
			return fmt.Errorf("insert shape line %s: %w", shapeID, err)
		}
	}
	for _, frequency := range feed.Frequencies {
		if _, err := tx.Exec(ctx, `
			INSERT INTO gtfs_frequency (trip_id, feed_version_id, agency_id, start_time, end_time, headway_secs, exact_times)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, frequency.TripID, feedVersionID, agencyID, frequency.StartTime, frequency.EndTime, frequency.HeadwaySecs, frequency.ExactTimes); err != nil {
			return fmt.Errorf("insert frequency %s/%s: %w", frequency.TripID, frequency.StartTime, err)
		}
	}
	return nil
}

func nullString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func nullInt(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}

func nullFloat(value *float64) any {
	if value == nil {
		return nil
	}
	return *value
}

func nullInt64(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

type parsedFeed struct {
	AgencyID           string
	Agencies           []importAgency
	Routes             []importRoute
	Stops              []importStop
	Calendars          []importCalendar
	CalendarDates      []importCalendarDate
	Trips              []importTrip
	StopTimes          []importStopTime
	ShapePoints        []importShapePoint
	ShapePointsByShape map[string][]importShapePoint
	Frequencies        []importFrequency
}

func (f parsedFeed) selectedAgency() (importAgency, bool) {
	for _, agency := range f.Agencies {
		if agency.ID == f.AgencyID {
			return agency, true
		}
	}
	return importAgency{}, false
}

type importAgency struct {
	ID       string
	Name     string
	URL      string
	Timezone string
	Email    string
}

type importRoute struct {
	ID        string
	AgencyID  string
	ShortName string
	LongName  string
	RouteType int
}

type importStop struct {
	ID   string
	Name string
	Lat  float64
	Lon  float64
}

type importCalendar struct {
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

type importCalendarDate struct {
	ServiceID     string
	Date          string
	ExceptionType int
}

type importTrip struct {
	ID          string
	RouteID     string
	ServiceID   string
	BlockID     string
	ShapeID     string
	DirectionID any
}

type importStopTime struct {
	TripID            string
	ArrivalTime       string
	DepartureTime     string
	StopID            string
	StopSequence      int
	PickupType        any
	DropOffType       any
	ShapeDistTraveled any
}

type importShapePoint struct {
	ShapeID      string
	Lat          float64
	Lon          float64
	Sequence     int
	DistTraveled any
}

type importFrequency struct {
	TripID      string
	StartTime   string
	EndTime     string
	HeadwaySecs int
	ExactTimes  int
}

func parseGTFSZip(payload []byte, agencyID string) (parsedFeed, ImportReport) {
	report := ImportReport{
		Status: ImportStatusStarted,
		Counts: map[string]int{},
	}
	feed := parsedFeed{
		AgencyID:           agencyID,
		ShapePointsByShape: map[string][]importShapePoint{},
	}
	reader, err := zip.NewReader(bytes.NewReader(payload), int64(len(payload)))
	if err != nil {
		report.addError("", 0, "invalid_zip", fmt.Sprintf("read gtfs zip: %v", err))
		report.Status = ImportStatusFailed
		return feed, report
	}
	files := map[string]*zip.File{}
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		name := strings.TrimPrefix(filepath.ToSlash(file.Name), "./")
		if strings.Contains(name, "/") {
			continue
		}
		files[name] = file
	}

	required := []string{"agency.txt", "routes.txt", "stops.txt", "trips.txt", "stop_times.txt"}
	for _, name := range required {
		if files[name] == nil {
			report.addError(name, 0, "missing_required_file", fmt.Sprintf("%s is required", name))
		}
	}
	if files["calendar.txt"] == nil && files["calendar_dates.txt"] == nil {
		report.addError("", 0, "missing_service_source", "at least one usable service source from calendar.txt or calendar_dates.txt is required")
	}

	feed.Agencies = parseAgencies(openCSV(files["agency.txt"], &report), agencyID, &report)
	feed.Routes = parseRoutes(openCSV(files["routes.txt"], &report), &report)
	feed.Stops = parseStops(openCSV(files["stops.txt"], &report), &report)
	feed.Calendars = parseCalendars(openCSV(files["calendar.txt"], &report), &report)
	feed.CalendarDates = parseCalendarDates(openCSV(files["calendar_dates.txt"], &report), &report)
	feed.Trips = parseTrips(openCSV(files["trips.txt"], &report), &report)
	feed.StopTimes = parseStopTimes(openCSV(files["stop_times.txt"], &report), &report)
	feed.ShapePoints = parseShapes(openCSV(files["shapes.txt"], &report), &report)
	feed.Frequencies = parseFrequencies(openCSV(files["frequencies.txt"], &report), &report)

	report.Counts["agency"] = len(feed.Agencies)
	report.Counts["routes"] = len(feed.Routes)
	report.Counts["stops"] = len(feed.Stops)
	report.Counts["calendar"] = len(feed.Calendars)
	report.Counts["calendar_dates"] = len(feed.CalendarDates)
	report.Counts["trips"] = len(feed.Trips)
	report.Counts["stop_times"] = len(feed.StopTimes)
	report.Counts["shapes"] = len(feed.ShapePoints)
	report.Counts["frequencies"] = len(feed.Frequencies)

	validateFeed(&feed, &report)
	report.Status = ImportStatusPublished
	if report.hasErrors() {
		report.Status = ImportStatusFailed
	}
	return feed, report
}

type csvTable struct {
	name string
	rows []csvRow
}

type csvRow struct {
	file   string
	number int
	values map[string]string
}

func openCSV(file *zip.File, report *ImportReport) csvTable {
	if file == nil {
		return csvTable{}
	}
	rc, err := file.Open()
	if err != nil {
		report.addError(file.Name, 0, "read_file_failed", err.Error())
		return csvTable{}
	}
	defer rc.Close()
	payload, err := io.ReadAll(rc)
	if err != nil {
		report.addError(file.Name, 0, "read_file_failed", err.Error())
		return csvTable{}
	}
	reader := csv.NewReader(bytes.NewReader(payload))
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		report.addError(file.Name, 0, "parse_csv_failed", err.Error())
		return csvTable{}
	}
	if len(records) == 0 {
		report.addError(file.Name, 0, "missing_header", "file is empty")
		return csvTable{}
	}
	header := make([]string, len(records[0]))
	for i, raw := range records[0] {
		header[i] = strings.TrimPrefix(strings.TrimSpace(raw), "\ufeff")
	}
	rows := make([]csvRow, 0, len(records)-1)
	for i, record := range records[1:] {
		values := make(map[string]string, len(header))
		for j, key := range header {
			if j < len(record) {
				values[key] = strings.TrimSpace(record[j])
			} else {
				values[key] = ""
			}
		}
		rows = append(rows, csvRow{file: file.Name, number: i + 2, values: values})
	}
	return csvTable{name: file.Name, rows: rows}
}

func parseAgencies(table csvTable, requestedAgencyID string, report *ImportReport) []importAgency {
	var agencies []importAgency
	for _, row := range table.rows {
		id := row.required("agency_id", report)
		agency := importAgency{
			ID:       id,
			Name:     row.required("agency_name", report),
			URL:      row.required("agency_url", report),
			Timezone: row.required("agency_timezone", report),
			Email:    row.get("agency_email"),
		}
		if agency.Timezone != "" {
			if _, err := time.LoadLocation(agency.Timezone); err != nil {
				report.addError(row.file, row.number, "invalid_timezone", fmt.Sprintf("agency_timezone %q is invalid", agency.Timezone))
			}
		}
		if agency.ID == requestedAgencyID {
			agencies = append(agencies, agency)
		}
	}
	if len(agencies) == 0 && table.name != "" {
		report.addError(table.name, 0, "agency_not_found", fmt.Sprintf("agency.txt must contain requested agency_id %q", requestedAgencyID))
	}
	if len(agencies) > 1 {
		report.addError(table.name, 0, "duplicate_agency", fmt.Sprintf("agency.txt contains multiple rows for agency_id %q", requestedAgencyID))
	}
	return agencies
}

func parseRoutes(table csvTable, report *ImportReport) []importRoute {
	var routes []importRoute
	seen := map[string]bool{}
	for _, row := range table.rows {
		id := row.required("route_id", report)
		if id != "" && seen[id] {
			report.addError(row.file, row.number, "duplicate_route", fmt.Sprintf("route_id %q is duplicated", id))
		}
		seen[id] = true
		routeType := row.requiredInt("route_type", report)
		if !supportedRouteType(routeType) {
			report.addError(row.file, row.number, "invalid_route_type", "route_type must be 0-7 or an extended GTFS route type from 100 through 1702")
		}
		routes = append(routes, importRoute{
			ID:        id,
			AgencyID:  row.get("agency_id"),
			ShortName: row.get("route_short_name"),
			LongName:  row.get("route_long_name"),
			RouteType: routeType,
		})
	}
	return routes
}

func parseStops(table csvTable, report *ImportReport) []importStop {
	var stops []importStop
	seen := map[string]bool{}
	for _, row := range table.rows {
		id := row.required("stop_id", report)
		if id != "" && seen[id] {
			report.addError(row.file, row.number, "duplicate_stop", fmt.Sprintf("stop_id %q is duplicated", id))
		}
		seen[id] = true
		lat := row.requiredFloat("stop_lat", report)
		lon := row.requiredFloat("stop_lon", report)
		if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
			report.addError(row.file, row.number, "invalid_stop_coordinate", "stop coordinates are outside valid latitude/longitude range")
		}
		stops = append(stops, importStop{
			ID:   id,
			Name: row.required("stop_name", report),
			Lat:  lat,
			Lon:  lon,
		})
	}
	return stops
}

func parseCalendars(table csvTable, report *ImportReport) []importCalendar {
	var calendars []importCalendar
	seen := map[string]bool{}
	for _, row := range table.rows {
		serviceID := row.required("service_id", report)
		if serviceID != "" && seen[serviceID] {
			report.addError(row.file, row.number, "duplicate_calendar_service", fmt.Sprintf("calendar service_id %q is duplicated", serviceID))
		}
		seen[serviceID] = true
		calendar := importCalendar{
			ServiceID: serviceID,
			Monday:    row.requiredBool("monday", report),
			Tuesday:   row.requiredBool("tuesday", report),
			Wednesday: row.requiredBool("wednesday", report),
			Thursday:  row.requiredBool("thursday", report),
			Friday:    row.requiredBool("friday", report),
			Saturday:  row.requiredBool("saturday", report),
			Sunday:    row.requiredBool("sunday", report),
			StartDate: row.requiredDate("start_date", report),
			EndDate:   row.requiredDate("end_date", report),
		}
		if calendar.StartDate != "" && calendar.EndDate != "" && calendar.StartDate > calendar.EndDate {
			report.addError(row.file, row.number, "invalid_calendar_range", "start_date must be before or equal to end_date")
		}
		calendars = append(calendars, calendar)
	}
	return calendars
}

func parseCalendarDates(table csvTable, report *ImportReport) []importCalendarDate {
	var dates []importCalendarDate
	seen := map[string]bool{}
	for _, row := range table.rows {
		serviceID := row.required("service_id", report)
		date := row.requiredDate("date", report)
		key := serviceID + "\x00" + date
		if serviceID != "" && date != "" && seen[key] {
			report.addError(row.file, row.number, "duplicate_calendar_date", fmt.Sprintf("calendar date %s/%s is duplicated", serviceID, date))
		}
		seen[key] = true
		exceptionType := row.requiredInt("exception_type", report)
		if exceptionType != 1 && exceptionType != 2 {
			report.addError(row.file, row.number, "invalid_exception_type", "exception_type must be 1 or 2")
		}
		dates = append(dates, importCalendarDate{ServiceID: serviceID, Date: date, ExceptionType: exceptionType})
	}
	return dates
}

func parseTrips(table csvTable, report *ImportReport) []importTrip {
	var trips []importTrip
	seen := map[string]bool{}
	for _, row := range table.rows {
		id := row.required("trip_id", report)
		if id != "" && seen[id] {
			report.addError(row.file, row.number, "duplicate_trip", fmt.Sprintf("trip_id %q is duplicated", id))
		}
		seen[id] = true
		directionID := row.optionalInt("direction_id", report)
		if directionID != nil && *directionID != 0 && *directionID != 1 {
			report.addError(row.file, row.number, "invalid_direction_id", "direction_id must be 0 or 1 when present")
		}
		trips = append(trips, importTrip{
			ID:          id,
			RouteID:     row.required("route_id", report),
			ServiceID:   row.required("service_id", report),
			BlockID:     row.get("block_id"),
			ShapeID:     row.get("shape_id"),
			DirectionID: nullInt(directionID),
		})
	}
	return trips
}

func parseStopTimes(table csvTable, report *ImportReport) []importStopTime {
	var stopTimes []importStopTime
	seen := map[string]bool{}
	for _, row := range table.rows {
		tripID := row.required("trip_id", report)
		arrival := row.get("arrival_time")
		departure := row.get("departure_time")
		if arrival == "" && departure == "" {
			report.addError(row.file, row.number, "missing_stop_time", "arrival_time or departure_time is required")
		}
		if arrival != "" {
			row.validateGTFSTime("arrival_time", report)
		}
		if departure != "" {
			row.validateGTFSTime("departure_time", report)
		}
		sequence := row.requiredInt("stop_sequence", report)
		key := tripID + "\x00" + strconv.Itoa(sequence)
		if tripID != "" && sequence != 0 && seen[key] {
			report.addError(row.file, row.number, "duplicate_stop_sequence", fmt.Sprintf("trip_id %q stop_sequence %d is duplicated", tripID, sequence))
		}
		seen[key] = true
		stopTimes = append(stopTimes, importStopTime{
			TripID:            tripID,
			ArrivalTime:       arrival,
			DepartureTime:     departure,
			StopID:            row.required("stop_id", report),
			StopSequence:      sequence,
			PickupType:        nullInt(row.optionalInt("pickup_type", report)),
			DropOffType:       nullInt(row.optionalInt("drop_off_type", report)),
			ShapeDistTraveled: nullFloat(row.optionalFloat("shape_dist_traveled", report)),
		})
	}
	return stopTimes
}

func parseShapes(table csvTable, report *ImportReport) []importShapePoint {
	var points []importShapePoint
	seen := map[string]bool{}
	for _, row := range table.rows {
		shapeID := row.required("shape_id", report)
		sequence := row.requiredInt("shape_pt_sequence", report)
		key := shapeID + "\x00" + strconv.Itoa(sequence)
		if shapeID != "" && sequence != 0 && seen[key] {
			report.addError(row.file, row.number, "duplicate_shape_sequence", fmt.Sprintf("shape_id %q shape_pt_sequence %d is duplicated", shapeID, sequence))
		}
		seen[key] = true
		lat := row.requiredFloat("shape_pt_lat", report)
		lon := row.requiredFloat("shape_pt_lon", report)
		if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
			report.addError(row.file, row.number, "invalid_shape_coordinate", "shape point coordinates are outside valid latitude/longitude range")
		}
		points = append(points, importShapePoint{
			ShapeID:      shapeID,
			Lat:          lat,
			Lon:          lon,
			Sequence:     sequence,
			DistTraveled: nullFloat(row.optionalFloat("shape_dist_traveled", report)),
		})
	}
	return points
}

func parseFrequencies(table csvTable, report *ImportReport) []importFrequency {
	var frequencies []importFrequency
	seen := map[string]bool{}
	for _, row := range table.rows {
		tripID := row.required("trip_id", report)
		startTime := row.required("start_time", report)
		endTime := row.required("end_time", report)
		if startTime != "" {
			row.validateGTFSTime("start_time", report)
		}
		if endTime != "" {
			row.validateGTFSTime("end_time", report)
		}
		key := tripID + "\x00" + startTime
		if tripID != "" && startTime != "" && seen[key] {
			report.addError(row.file, row.number, "duplicate_frequency", fmt.Sprintf("frequency %s/%s is duplicated", tripID, startTime))
		}
		seen[key] = true
		headway := row.requiredInt("headway_secs", report)
		if headway <= 0 {
			report.addError(row.file, row.number, "invalid_headway", "headway_secs must be greater than 0")
		}
		exactTimes := 0
		if raw := row.get("exact_times"); raw != "" {
			exactTimes = row.requiredInt("exact_times", report)
			if exactTimes != 0 && exactTimes != 1 {
				report.addError(row.file, row.number, "invalid_exact_times", "exact_times must be 0 or 1")
			}
		}
		frequencies = append(frequencies, importFrequency{
			TripID:      tripID,
			StartTime:   startTime,
			EndTime:     endTime,
			HeadwaySecs: headway,
			ExactTimes:  exactTimes,
		})
	}
	return frequencies
}

func validateFeed(feed *parsedFeed, report *ImportReport) {
	routes := map[string]bool{}
	for _, route := range feed.Routes {
		routes[route.ID] = true
		if route.AgencyID != "" && route.AgencyID != feed.AgencyID {
			report.addError("routes.txt", 0, "route_agency_mismatch", fmt.Sprintf("route_id %q references agency_id %q, want %q", route.ID, route.AgencyID, feed.AgencyID))
		}
	}
	stops := map[string]bool{}
	for _, stop := range feed.Stops {
		stops[stop.ID] = true
	}
	usableServices := map[string]bool{}
	for _, calendar := range feed.Calendars {
		if !calendar.hasActiveWeekday() {
			report.addError("calendar.txt", 0, "unusable_calendar_service", fmt.Sprintf("service_id %q has no active weekdays", calendar.ServiceID))
			continue
		}
		if calendar.ServiceID != "" && calendar.StartDate != "" && calendar.EndDate != "" && calendar.StartDate <= calendar.EndDate {
			usableServices[calendar.ServiceID] = true
		}
	}
	for _, date := range feed.CalendarDates {
		if date.ExceptionType == 1 && date.ServiceID != "" && date.Date != "" {
			usableServices[date.ServiceID] = true
		}
	}
	if len(usableServices) == 0 {
		report.addError("", 0, "missing_usable_service_source", "calendar.txt or calendar_dates.txt must contain at least one usable service: an active calendar weekday or a calendar_dates exception_type=1 addition")
	}

	shapes := map[string]bool{}
	byShape := map[string][]importShapePoint{}
	for _, point := range feed.ShapePoints {
		shapes[point.ShapeID] = true
		byShape[point.ShapeID] = append(byShape[point.ShapeID], point)
	}
	for shapeID, points := range byShape {
		sort.Slice(points, func(i, j int) bool { return points[i].Sequence < points[j].Sequence })
		var previousSequence int
		var previousDist float64
		var previousHasDist bool
		for i, point := range points {
			if i > 0 && point.Sequence <= previousSequence {
				report.addError("shapes.txt", 0, "invalid_shape_order", fmt.Sprintf("shape_id %q sequence values must be strictly increasing", shapeID))
			}
			previousSequence = point.Sequence
			if dist, ok := point.DistTraveled.(float64); ok {
				if previousHasDist && dist < previousDist {
					report.addError("shapes.txt", 0, "invalid_shape_distance_order", fmt.Sprintf("shape_id %q shape_dist_traveled must be nondecreasing", shapeID))
				}
				previousDist = dist
				previousHasDist = true
			}
		}
		feed.ShapePointsByShape[shapeID] = points
	}

	trips := map[string]importTrip{}
	for _, trip := range feed.Trips {
		trips[trip.ID] = trip
		if !routes[trip.RouteID] {
			report.addError("trips.txt", 0, "unknown_route", fmt.Sprintf("trip_id %q references unknown route_id %q", trip.ID, trip.RouteID))
		}
		if !usableServices[trip.ServiceID] {
			report.addError("trips.txt", 0, "unknown_service", fmt.Sprintf("trip_id %q references unknown service_id %q", trip.ID, trip.ServiceID))
		}
		if trip.ShapeID != "" && len(feed.ShapePoints) > 0 && !shapes[trip.ShapeID] {
			report.addError("trips.txt", 0, "unknown_shape", fmt.Sprintf("trip_id %q references unknown shape_id %q", trip.ID, trip.ShapeID))
		}
	}
	for _, stopTime := range feed.StopTimes {
		if _, ok := trips[stopTime.TripID]; !ok {
			report.addError("stop_times.txt", 0, "unknown_trip", fmt.Sprintf("stop_times references unknown trip_id %q", stopTime.TripID))
		}
		if !stops[stopTime.StopID] {
			report.addError("stop_times.txt", 0, "unknown_stop", fmt.Sprintf("stop_times references unknown stop_id %q", stopTime.StopID))
		}
	}
	for _, frequency := range feed.Frequencies {
		if _, ok := trips[frequency.TripID]; !ok {
			report.addError("frequencies.txt", 0, "unknown_frequency_trip", fmt.Sprintf("frequencies references unknown trip_id %q", frequency.TripID))
		}
		start, startErr := ParseGTFSTime(frequency.StartTime)
		end, endErr := ParseGTFSTime(frequency.EndTime)
		if startErr == nil && endErr == nil && end <= start {
			report.addError("frequencies.txt", 0, "invalid_frequency_range", fmt.Sprintf("frequency for trip_id %q must have end_time after start_time", frequency.TripID))
		}
	}
}

func supportedRouteType(routeType int) bool {
	return (routeType >= 0 && routeType <= 7) || (routeType >= 100 && routeType <= 1702)
}

func (c importCalendar) hasActiveWeekday() bool {
	return c.Monday || c.Tuesday || c.Wednesday || c.Thursday || c.Friday || c.Saturday || c.Sunday
}

func (r csvRow) get(field string) string {
	return strings.TrimSpace(r.values[field])
}

func (r csvRow) required(field string, report *ImportReport) string {
	value := r.get(field)
	if value == "" {
		report.addError(r.file, r.number, "missing_required_field", fmt.Sprintf("%s is required", field))
	}
	return value
}

func (r csvRow) requiredInt(field string, report *ImportReport) int {
	raw := r.required(field, report)
	if raw == "" {
		return 0
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		report.addError(r.file, r.number, "invalid_integer", fmt.Sprintf("%s must be an integer", field))
		return 0
	}
	return value
}

func (r csvRow) optionalInt(field string, report *ImportReport) *int {
	raw := r.get(field)
	if raw == "" {
		return nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		report.addError(r.file, r.number, "invalid_integer", fmt.Sprintf("%s must be an integer", field))
		return nil
	}
	return &value
}

func (r csvRow) requiredFloat(field string, report *ImportReport) float64 {
	raw := r.required(field, report)
	if raw == "" {
		return 0
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		report.addError(r.file, r.number, "invalid_float", fmt.Sprintf("%s must be a finite number", field))
		return 0
	}
	return value
}

func (r csvRow) optionalFloat(field string, report *ImportReport) *float64 {
	raw := r.get(field)
	if raw == "" {
		return nil
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		report.addError(r.file, r.number, "invalid_float", fmt.Sprintf("%s must be a finite number", field))
		return nil
	}
	return &value
}

func (r csvRow) requiredBool(field string, report *ImportReport) bool {
	raw := r.required(field, report)
	switch raw {
	case "1":
		return true
	case "0":
		return false
	default:
		report.addError(r.file, r.number, "invalid_boolean", fmt.Sprintf("%s must be 0 or 1", field))
		return false
	}
}

func (r csvRow) requiredDate(field string, report *ImportReport) string {
	raw := r.required(field, report)
	if raw == "" {
		return ""
	}
	if len(raw) != 8 {
		report.addError(r.file, r.number, "invalid_date", fmt.Sprintf("%s must use YYYYMMDD", field))
		return raw
	}
	if _, err := time.Parse("20060102", raw); err != nil {
		report.addError(r.file, r.number, "invalid_date", fmt.Sprintf("%s must use YYYYMMDD", field))
	}
	return raw
}

func (r csvRow) validateGTFSTime(field string, report *ImportReport) {
	if _, err := ParseGTFSTime(r.get(field)); err != nil {
		report.addError(r.file, r.number, "invalid_gtfs_time", fmt.Sprintf("%s is invalid: %v", field, err))
	}
}

func (r *ImportReport) addError(file string, row int, code string, message string) {
	r.Errors = append(r.Errors, ImportMessage{File: file, Row: row, Code: code, Message: message})
}

func (r *ImportReport) addWarning(file string, row int, code string, message string) {
	r.Warnings = append(r.Warnings, ImportMessage{File: file, Row: row, Code: code, Message: message})
}

func IsImportError(err error) bool {
	var importErr *ImportError
	return errors.As(err, &importErr)
}
