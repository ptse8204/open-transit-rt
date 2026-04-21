package gtfs

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Agency(ctx context.Context, agencyID string) (Agency, error) {
	var agency Agency
	err := r.pool.QueryRow(ctx, `
		SELECT id, timezone
		FROM agency
		WHERE id = $1
	`, agencyID).Scan(&agency.ID, &agency.Timezone)
	if err != nil {
		return Agency{}, fmt.Errorf("query agency timezone: %w", err)
	}
	return agency, nil
}

func (r *PostgresRepository) ActiveFeedVersion(ctx context.Context, agencyID string) (FeedVersion, error) {
	var feed FeedVersion
	err := r.pool.QueryRow(ctx, `
		SELECT id, agency_id
		FROM feed_version
		WHERE agency_id = $1
		  AND is_active
		ORDER BY activated_at DESC NULLS LAST, created_at DESC
		LIMIT 1
	`, agencyID).Scan(&feed.ID, &feed.AgencyID)
	if err != nil {
		return FeedVersion{}, fmt.Errorf("query active feed version: %w", err)
	}
	return feed, nil
}

func (r *PostgresRepository) ListTripCandidates(ctx context.Context, agencyID string, feedVersionID string, serviceDate string) ([]TripCandidate, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT t.id, t.route_id, t.service_id, t.block_id, t.shape_id, t.direction_id
		FROM gtfs_trip t
		WHERE t.agency_id = $1
		  AND t.feed_version_id = $2
		  AND (
		    EXISTS (
		      SELECT 1
		      FROM gtfs_calendar c
		      WHERE c.agency_id = t.agency_id
		        AND c.feed_version_id = t.feed_version_id
		        AND c.service_id = t.service_id
		        AND c.start_date <= $3
		        AND c.end_date >= $3
		        AND CASE EXTRACT(ISODOW FROM to_date($3, 'YYYYMMDD'))::int
		          WHEN 1 THEN c.monday
		          WHEN 2 THEN c.tuesday
		          WHEN 3 THEN c.wednesday
		          WHEN 4 THEN c.thursday
		          WHEN 5 THEN c.friday
		          WHEN 6 THEN c.saturday
		          WHEN 7 THEN c.sunday
		        END
		    )
		    OR EXISTS (
		      SELECT 1
		      FROM gtfs_calendar_date cd
		      WHERE cd.agency_id = t.agency_id
		        AND cd.feed_version_id = t.feed_version_id
		        AND cd.service_id = t.service_id
		        AND cd.date = $3
		        AND cd.exception_type = 1
		    )
		  )
		  AND NOT EXISTS (
		    SELECT 1
		    FROM gtfs_calendar_date cd
		    WHERE cd.agency_id = t.agency_id
		      AND cd.feed_version_id = t.feed_version_id
		      AND cd.service_id = t.service_id
		      AND cd.date = $3
		      AND cd.exception_type = 2
		  )
		ORDER BY t.route_id, t.id
	`, agencyID, feedVersionID, serviceDate)
	if err != nil {
		return nil, fmt.Errorf("query trip candidates: %w", err)
	}
	defer rows.Close()

	var trips []TripCandidate
	for rows.Next() {
		var trip TripCandidate
		var direction *int
		var blockID, shapeID sql.NullString
		if err := rows.Scan(&trip.TripID, &trip.RouteID, &trip.ServiceID, &blockID, &shapeID, &direction); err != nil {
			return nil, fmt.Errorf("scan trip candidate: %w", err)
		}
		trip.AgencyID = agencyID
		trip.FeedVersionID = feedVersionID
		trip.ServiceDate = serviceDate
		trip.BlockID = blockID.String
		trip.ShapeID = shapeID.String
		trip.DirectionID = direction
		trips = append(trips, trip)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate trip candidates: %w", err)
	}

	for i := range trips {
		stopTimes, err := r.stopTimes(ctx, agencyID, feedVersionID, trips[i].TripID)
		if err != nil {
			return nil, err
		}
		shapePoints, err := r.shapePoints(ctx, agencyID, feedVersionID, trips[i].ShapeID)
		if err != nil {
			return nil, err
		}
		frequencies, err := r.frequencies(ctx, agencyID, feedVersionID, trips[i].TripID)
		if err != nil {
			return nil, err
		}
		trips[i].StopTimes = stopTimes
		trips[i].ShapePoints = shapePoints
		trips[i].Frequencies = frequencies
	}

	return trips, nil
}

func (r *PostgresRepository) stopTimes(ctx context.Context, agencyID string, feedVersionID string, tripID string) ([]StopTime, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT trip_id, stop_id, COALESCE(arrival_time, departure_time), COALESCE(departure_time, arrival_time), stop_sequence, COALESCE(shape_dist_traveled, 0)
		FROM gtfs_stop_time
		WHERE agency_id = $1
		  AND feed_version_id = $2
		  AND trip_id = $3
		ORDER BY stop_sequence
	`, agencyID, feedVersionID, tripID)
	if err != nil {
		return nil, fmt.Errorf("query stop times: %w", err)
	}
	defer rows.Close()

	var stopTimes []StopTime
	for rows.Next() {
		var st StopTime
		var arrival, departure string
		if err := rows.Scan(&st.TripID, &st.StopID, &arrival, &departure, &st.StopSequence, &st.ShapeDistTraveled); err != nil {
			return nil, fmt.Errorf("scan stop time: %w", err)
		}
		arrivalSeconds, err := ParseGTFSTime(arrival)
		if err != nil {
			return nil, err
		}
		departureSeconds, err := ParseGTFSTime(departure)
		if err != nil {
			return nil, err
		}
		st.ArrivalSeconds = arrivalSeconds
		st.DepartureSeconds = departureSeconds
		stopTimes = append(stopTimes, st)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stop times: %w", err)
	}
	return stopTimes, nil
}

func (r *PostgresRepository) shapePoints(ctx context.Context, agencyID string, feedVersionID string, shapeID string) ([]ShapePoint, error) {
	if shapeID == "" {
		return nil, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT shape_id, lat, lon, sequence, dist_traveled
		FROM gtfs_shape_point
		WHERE agency_id = $1
		  AND feed_version_id = $2
		  AND shape_id = $3
		ORDER BY sequence
	`, agencyID, feedVersionID, shapeID)
	if err != nil {
		return nil, fmt.Errorf("query shape points: %w", err)
	}
	defer rows.Close()

	var points []ShapePoint
	for rows.Next() {
		var point ShapePoint
		var dist *float64
		if err := rows.Scan(&point.ShapeID, &point.Lat, &point.Lon, &point.Sequence, &dist); err != nil {
			return nil, fmt.Errorf("scan shape point: %w", err)
		}
		if dist != nil {
			point.DistTraveled = *dist
			point.HasDistance = true
		}
		points = append(points, point)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shape points: %w", err)
	}
	return points, nil
}

func (r *PostgresRepository) frequencies(ctx context.Context, agencyID string, feedVersionID string, tripID string) ([]Frequency, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT trip_id, start_time, end_time, headway_secs, exact_times
		FROM gtfs_frequency
		WHERE agency_id = $1
		  AND feed_version_id = $2
		  AND trip_id = $3
		ORDER BY start_time
	`, agencyID, feedVersionID, tripID)
	if err != nil {
		return nil, fmt.Errorf("query frequencies: %w", err)
	}
	defer rows.Close()

	var frequencies []Frequency
	for rows.Next() {
		var f Frequency
		if err := rows.Scan(&f.TripID, &f.StartTime, &f.EndTime, &f.HeadwaySecs, &f.ExactTimes); err != nil {
			return nil, fmt.Errorf("scan frequency: %w", err)
		}
		startSeconds, err := ParseGTFSTime(f.StartTime)
		if err != nil {
			return nil, err
		}
		endSeconds, err := ParseGTFSTime(f.EndTime)
		if err != nil {
			return nil, err
		}
		f.StartSeconds = startSeconds
		f.EndSeconds = endSeconds
		frequencies = append(frequencies, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate frequencies: %w", err)
	}
	return frequencies, nil
}

func IsNoRows(err error) bool {
	return err == pgx.ErrNoRows
}
