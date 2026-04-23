package schedule

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Snapshot struct {
	AgencyID      string
	FeedVersionID string
	GeneratedAt   time.Time
	RevisionTime  time.Time
	Payload       []byte
}

type Builder struct {
	pool     *pgxpool.Pool
	agencyID string
}

func NewBuilder(pool *pgxpool.Pool, agencyID string) (*Builder, error) {
	if pool == nil {
		return nil, fmt.Errorf("database pool is required")
	}
	if agencyID == "" {
		return nil, fmt.Errorf("AGENCY_ID is required")
	}
	return &Builder{pool: pool, agencyID: agencyID}, nil
}

func (b *Builder) Snapshot(ctx context.Context, generatedAt time.Time) (Snapshot, error) {
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}
	feed, err := b.activeFeed(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	payload, err := b.buildZIP(ctx, feed)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{
		AgencyID:      b.agencyID,
		FeedVersionID: feed.ID,
		GeneratedAt:   generatedAt.UTC(),
		RevisionTime:  feed.RevisionTime.UTC(),
		Payload:       payload,
	}, nil
}

func (b *Builder) SnapshotForFeedVersion(ctx context.Context, feedVersionID string, generatedAt time.Time) (Snapshot, error) {
	if feedVersionID == "" {
		return b.Snapshot(ctx, generatedAt)
	}
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}
	feed, err := b.feedByID(ctx, feedVersionID)
	if err != nil {
		return Snapshot{}, err
	}
	payload, err := b.buildZIP(ctx, feed)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{
		AgencyID:      b.agencyID,
		FeedVersionID: feed.ID,
		GeneratedAt:   generatedAt.UTC(),
		RevisionTime:  feed.RevisionTime.UTC(),
		Payload:       payload,
	}, nil
}

func (b *Builder) Ready(ctx context.Context) error {
	if _, err := b.activeFeed(ctx); err != nil {
		return fmt.Errorf("active schedule feed unavailable: %w", err)
	}
	return nil
}

type activeFeed struct {
	ID           string
	RevisionTime time.Time
}

func (b *Builder) activeFeed(ctx context.Context) (activeFeed, error) {
	var feed activeFeed
	err := b.pool.QueryRow(ctx, `
		SELECT id, COALESCE(activated_at, published_at, created_at)
		FROM feed_version
		WHERE agency_id = $1
		  AND is_active
		ORDER BY activated_at DESC NULLS LAST, created_at DESC
		LIMIT 1
	`, b.agencyID).Scan(&feed.ID, &feed.RevisionTime)
	if err != nil {
		return activeFeed{}, fmt.Errorf("query active schedule feed: %w", err)
	}
	return feed, nil
}

func (b *Builder) feedByID(ctx context.Context, feedVersionID string) (activeFeed, error) {
	var feed activeFeed
	err := b.pool.QueryRow(ctx, `
		SELECT id, COALESCE(activated_at, published_at, created_at)
		FROM feed_version
		WHERE agency_id = $1
		  AND id = $2
		  AND lifecycle_state IN ('staged', 'active', 'retired')
	`, b.agencyID, feedVersionID).Scan(&feed.ID, &feed.RevisionTime)
	if err != nil {
		return activeFeed{}, fmt.Errorf("query schedule feed version: %w", err)
	}
	return feed, nil
}

func (b *Builder) buildZIP(ctx context.Context, feed activeFeed) ([]byte, error) {
	files := []zipFile{
		{name: "agency.txt", rows: b.agencyRows},
		{name: "feed_info.txt", rows: b.feedInfoRows},
		{name: "routes.txt", rows: b.routeRows},
		{name: "stops.txt", rows: b.stopRows},
		{name: "trips.txt", rows: b.tripRows},
		{name: "stop_times.txt", rows: b.stopTimeRows},
		{name: "calendar.txt", rows: b.calendarRows, optional: true},
		{name: "calendar_dates.txt", rows: b.calendarDateRows, optional: true},
		{name: "shapes.txt", rows: b.shapeRows, optional: true},
		{name: "frequencies.txt", rows: b.frequencyRows, optional: true},
	}
	var buffer bytes.Buffer
	zw := zip.NewWriter(&buffer)
	for _, file := range files {
		rows, err := file.rows(ctx, feed.ID)
		if err != nil {
			_ = zw.Close()
			return nil, err
		}
		if file.optional && len(rows) == 1 {
			continue
		}
		if err := writeZipCSV(zw, file.name, feed.RevisionTime, rows); err != nil {
			_ = zw.Close()
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("close schedule zip: %w", err)
	}
	return buffer.Bytes(), nil
}

type zipFile struct {
	name     string
	rows     func(context.Context, string) ([][]string, error)
	optional bool
}

func writeZipCSV(zw *zip.Writer, name string, modified time.Time, rows [][]string) error {
	header := &zip.FileHeader{Name: name, Method: zip.Deflate}
	header.SetModTime(modified.UTC())
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("create zip entry %s: %w", name, err)
	}
	cw := csv.NewWriter(writer)
	if err := cw.WriteAll(rows); err != nil {
		return fmt.Errorf("write csv %s: %w", name, err)
	}
	return nil
}

func (b *Builder) agencyRows(ctx context.Context, _ string) ([][]string, error) {
	rows := [][]string{{"agency_id", "agency_name", "agency_url", "agency_timezone", "agency_email"}}
	var name, timezone string
	var publicURL, email sql.NullString
	err := b.pool.QueryRow(ctx, `
		SELECT name, public_url, timezone, contact_email
		FROM agency
		WHERE id = $1
	`, b.agencyID).Scan(&name, &publicURL, &timezone, &email)
	if err != nil {
		return nil, fmt.Errorf("query agency for schedule zip: %w", err)
	}
	rows = append(rows, []string{b.agencyID, name, publicURL.String, timezone, email.String})
	return rows, nil
}

func (b *Builder) feedInfoRows(ctx context.Context, feedVersionID string) ([][]string, error) {
	rows := [][]string{{"feed_publisher_name", "feed_publisher_url", "feed_lang", "feed_start_date", "feed_end_date", "feed_version", "feed_contact_email"}}
	var name string
	var publicURL, email sql.NullString
	err := b.pool.QueryRow(ctx, `
		SELECT name, public_url, contact_email
		FROM agency
		WHERE id = $1
	`, b.agencyID).Scan(&name, &publicURL, &email)
	if err != nil {
		return nil, fmt.Errorf("query agency for feed_info: %w", err)
	}
	var startDate, endDate sql.NullString
	err = b.pool.QueryRow(ctx, `
		SELECT MIN(service_date), MAX(service_date)
		FROM (
			SELECT start_date AS service_date
			FROM gtfs_calendar
			WHERE agency_id = $1 AND feed_version_id = $2
			UNION ALL
			SELECT end_date AS service_date
			FROM gtfs_calendar
			WHERE agency_id = $1 AND feed_version_id = $2
			UNION ALL
			SELECT date AS service_date
			FROM gtfs_calendar_date
			WHERE agency_id = $1 AND feed_version_id = $2
		) service_dates
	`, b.agencyID, feedVersionID).Scan(&startDate, &endDate)
	if err != nil {
		return nil, fmt.Errorf("query service dates for feed_info: %w", err)
	}
	rows = append(rows, []string{name, publicURL.String, "en", startDate.String, endDate.String, feedVersionID, email.String})
	return rows, nil
}

func (b *Builder) routeRows(ctx context.Context, feedVersionID string) ([][]string, error) {
	rows := [][]string{{"route_id", "agency_id", "route_short_name", "route_long_name", "route_type"}}
	queryRows, err := b.pool.Query(ctx, `
		SELECT id, short_name, long_name, route_type
		FROM gtfs_route
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY id
	`, b.agencyID, feedVersionID)
	if err != nil {
		return nil, fmt.Errorf("query route rows: %w", err)
	}
	defer queryRows.Close()
	for queryRows.Next() {
		var id string
		var shortName, longName sql.NullString
		var routeType sql.NullInt64
		if err := queryRows.Scan(&id, &shortName, &longName, &routeType); err != nil {
			return nil, fmt.Errorf("scan route row: %w", err)
		}
		rows = append(rows, []string{id, b.agencyID, shortName.String, longName.String, intString(routeType)})
	}
	return rows, rowErr(queryRows, "route rows")
}

func (b *Builder) stopRows(ctx context.Context, feedVersionID string) ([][]string, error) {
	rows := [][]string{{"stop_id", "stop_name", "stop_lat", "stop_lon"}}
	queryRows, err := b.pool.Query(ctx, `
		SELECT id, name, lat, lon
		FROM gtfs_stop
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY id
	`, b.agencyID, feedVersionID)
	if err != nil {
		return nil, fmt.Errorf("query stop rows: %w", err)
	}
	defer queryRows.Close()
	for queryRows.Next() {
		var id, name string
		var lat, lon float64
		if err := queryRows.Scan(&id, &name, &lat, &lon); err != nil {
			return nil, fmt.Errorf("scan stop row: %w", err)
		}
		rows = append(rows, []string{id, name, floatString(lat), floatString(lon)})
	}
	return rows, rowErr(queryRows, "stop rows")
}

func (b *Builder) calendarRows(ctx context.Context, feedVersionID string) ([][]string, error) {
	rows := [][]string{{"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date"}}
	queryRows, err := b.pool.Query(ctx, `
		SELECT service_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday, start_date, end_date
		FROM gtfs_calendar
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY service_id
	`, b.agencyID, feedVersionID)
	if err != nil {
		return nil, fmt.Errorf("query calendar rows: %w", err)
	}
	defer queryRows.Close()
	for queryRows.Next() {
		var serviceID, startDate, endDate string
		var monday, tuesday, wednesday, thursday, friday, saturday, sunday bool
		if err := queryRows.Scan(&serviceID, &monday, &tuesday, &wednesday, &thursday, &friday, &saturday, &sunday, &startDate, &endDate); err != nil {
			return nil, fmt.Errorf("scan calendar row: %w", err)
		}
		rows = append(rows, []string{serviceID, boolString(monday), boolString(tuesday), boolString(wednesday), boolString(thursday), boolString(friday), boolString(saturday), boolString(sunday), startDate, endDate})
	}
	return rows, rowErr(queryRows, "calendar rows")
}

func (b *Builder) calendarDateRows(ctx context.Context, feedVersionID string) ([][]string, error) {
	rows := [][]string{{"service_id", "date", "exception_type"}}
	queryRows, err := b.pool.Query(ctx, `
		SELECT service_id, date, exception_type
		FROM gtfs_calendar_date
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY service_id, date
	`, b.agencyID, feedVersionID)
	if err != nil {
		return nil, fmt.Errorf("query calendar date rows: %w", err)
	}
	defer queryRows.Close()
	for queryRows.Next() {
		var serviceID, date string
		var exceptionType int
		if err := queryRows.Scan(&serviceID, &date, &exceptionType); err != nil {
			return nil, fmt.Errorf("scan calendar date row: %w", err)
		}
		rows = append(rows, []string{serviceID, date, strconv.Itoa(exceptionType)})
	}
	return rows, rowErr(queryRows, "calendar date rows")
}

func (b *Builder) tripRows(ctx context.Context, feedVersionID string) ([][]string, error) {
	rows := [][]string{{"route_id", "service_id", "trip_id", "direction_id", "block_id", "shape_id"}}
	queryRows, err := b.pool.Query(ctx, `
		SELECT route_id, service_id, id, direction_id, block_id, shape_id
		FROM gtfs_trip
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY route_id, id
	`, b.agencyID, feedVersionID)
	if err != nil {
		return nil, fmt.Errorf("query trip rows: %w", err)
	}
	defer queryRows.Close()
	for queryRows.Next() {
		var routeID, serviceID, tripID string
		var directionID sql.NullInt64
		var blockID, shapeID sql.NullString
		if err := queryRows.Scan(&routeID, &serviceID, &tripID, &directionID, &blockID, &shapeID); err != nil {
			return nil, fmt.Errorf("scan trip row: %w", err)
		}
		rows = append(rows, []string{routeID, serviceID, tripID, intString(directionID), blockID.String, shapeID.String})
	}
	return rows, rowErr(queryRows, "trip rows")
}

func (b *Builder) stopTimeRows(ctx context.Context, feedVersionID string) ([][]string, error) {
	rows := [][]string{{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence", "pickup_type", "drop_off_type", "shape_dist_traveled"}}
	queryRows, err := b.pool.Query(ctx, `
		SELECT trip_id, arrival_time, departure_time, stop_id, stop_sequence, pickup_type, drop_off_type, shape_dist_traveled
		FROM gtfs_stop_time
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY trip_id, stop_sequence
	`, b.agencyID, feedVersionID)
	if err != nil {
		return nil, fmt.Errorf("query stop time rows: %w", err)
	}
	defer queryRows.Close()
	for queryRows.Next() {
		var tripID, stopID string
		var arrival, departure sql.NullString
		var sequence int
		var pickup, dropOff sql.NullInt64
		var dist sql.NullFloat64
		if err := queryRows.Scan(&tripID, &arrival, &departure, &stopID, &sequence, &pickup, &dropOff, &dist); err != nil {
			return nil, fmt.Errorf("scan stop time row: %w", err)
		}
		rows = append(rows, []string{tripID, arrival.String, departure.String, stopID, strconv.Itoa(sequence), intString(pickup), intString(dropOff), nullFloatString(dist)})
	}
	return rows, rowErr(queryRows, "stop time rows")
}

func (b *Builder) shapeRows(ctx context.Context, feedVersionID string) ([][]string, error) {
	rows := [][]string{{"shape_id", "shape_pt_lat", "shape_pt_lon", "shape_pt_sequence", "shape_dist_traveled"}}
	queryRows, err := b.pool.Query(ctx, `
		SELECT shape_id, lat, lon, sequence, dist_traveled
		FROM gtfs_shape_point
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY shape_id, sequence
	`, b.agencyID, feedVersionID)
	if err != nil {
		return nil, fmt.Errorf("query shape rows: %w", err)
	}
	defer queryRows.Close()
	for queryRows.Next() {
		var shapeID string
		var lat, lon float64
		var sequence int
		var dist sql.NullFloat64
		if err := queryRows.Scan(&shapeID, &lat, &lon, &sequence, &dist); err != nil {
			return nil, fmt.Errorf("scan shape row: %w", err)
		}
		rows = append(rows, []string{shapeID, floatString(lat), floatString(lon), strconv.Itoa(sequence), nullFloatString(dist)})
	}
	return rows, rowErr(queryRows, "shape rows")
}

func (b *Builder) frequencyRows(ctx context.Context, feedVersionID string) ([][]string, error) {
	rows := [][]string{{"trip_id", "start_time", "end_time", "headway_secs", "exact_times"}}
	queryRows, err := b.pool.Query(ctx, `
		SELECT trip_id, start_time, end_time, headway_secs, exact_times
		FROM gtfs_frequency
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY trip_id, start_time
	`, b.agencyID, feedVersionID)
	if err != nil {
		return nil, fmt.Errorf("query frequency rows: %w", err)
	}
	defer queryRows.Close()
	for queryRows.Next() {
		var tripID, startTime, endTime string
		var headway, exact int
		if err := queryRows.Scan(&tripID, &startTime, &endTime, &headway, &exact); err != nil {
			return nil, fmt.Errorf("scan frequency row: %w", err)
		}
		rows = append(rows, []string{tripID, startTime, endTime, strconv.Itoa(headway), strconv.Itoa(exact)})
	}
	return rows, rowErr(queryRows, "frequency rows")
}

func rowErr(rows pgx.Rows, label string) error {
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate %s: %w", label, err)
	}
	return nil
}

func boolString(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func intString(value sql.NullInt64) string {
	if !value.Valid {
		return ""
	}
	return strconv.FormatInt(value.Int64, 10)
}

func floatString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func nullFloatString(value sql.NullFloat64) string {
	if !value.Valid {
		return ""
	}
	return floatString(value.Float64)
}
