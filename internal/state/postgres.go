package state

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) ActiveManualOverride(ctx context.Context, agencyID string, vehicleID string, at time.Time) (*ManualOverride, error) {
	var override ManualOverride
	var routeID, tripID, startDate, startTime, reason sql.NullString
	var expiresAt sql.NullTime
	err := r.pool.QueryRow(ctx, `
		SELECT id, agency_id, vehicle_id, override_type, route_id, trip_id, start_date, start_time, state, expires_at, reason, created_at
		FROM manual_override
		WHERE agency_id = $1
		  AND vehicle_id = $2
		  AND cleared_at IS NULL
		  AND (expires_at IS NULL OR expires_at > $3)
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, agencyID, vehicleID, at).Scan(
		&override.ID,
		&override.AgencyID,
		&override.VehicleID,
		&override.Type,
		&routeID,
		&tripID,
		&startDate,
		&startTime,
		&override.State,
		&expiresAt,
		&reason,
		&override.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query active manual override: %w", err)
	}
	override.RouteID = routeID.String
	override.TripID = tripID.String
	override.StartDate = startDate.String
	override.StartTime = startTime.String
	override.Reason = reason.String
	if expiresAt.Valid {
		override.ExpiresAt = &expiresAt.Time
	}
	return &override, nil
}

func (r *PostgresRepository) CurrentAssignment(ctx context.Context, agencyID string, vehicleID string) (*Assignment, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, agency_id, vehicle_id, feed_version_id, telemetry_event_id, service_date, route_id, trip_id, block_id, start_date, start_time,
		       current_stop_sequence, shape_dist_traveled, state, confidence, assignment_source, reason_codes, degraded_state, score_details_json,
		       manual_override_id, active_from
		FROM vehicle_trip_assignment
		WHERE agency_id = $1
		  AND vehicle_id = $2
		  AND active_to IS NULL
		ORDER BY active_from DESC, id DESC
		LIMIT 1
	`, agencyID, vehicleID)
	assignment, err := scanAssignment(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query current assignment: %w", err)
	}
	return &assignment, nil
}

func (r *PostgresRepository) SaveAssignment(ctx context.Context, assignment Assignment, incidents []Incident) (Assignment, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Assignment{}, fmt.Errorf("begin assignment transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if assignment.ActiveFrom.IsZero() {
		assignment.ActiveFrom = time.Now().UTC()
	}
	if assignment.AssignmentSource == "" {
		assignment.AssignmentSource = AssignmentSourceAutomatic
	}
	if assignment.DegradedState == "" {
		assignment.DegradedState = DegradedNone
	}
	if assignment.ScoreDetails == nil {
		assignment.ScoreDetails = map[string]any{"score_schema": "loose_debug_v1"}
	}

	current, err := currentAssignmentInTx(ctx, tx, assignment.AgencyID, assignment.VehicleID)
	if err != nil {
		return Assignment{}, err
	}
	if current != nil && repeatedDegradedAssignment(*current, assignment) {
		return *current, nil
	}

	if _, err := tx.Exec(ctx, `
		UPDATE vehicle_trip_assignment
		SET active_to = $3
		WHERE agency_id = $1
		  AND vehicle_id = $2
		  AND active_to IS NULL
	`, assignment.AgencyID, assignment.VehicleID, assignment.ActiveFrom); err != nil {
		return Assignment{}, fmt.Errorf("close current assignment: %w", err)
	}

	scoreDetails, err := json.Marshal(assignment.ScoreDetails)
	if err != nil {
		return Assignment{}, fmt.Errorf("marshal score details: %w", err)
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO vehicle_trip_assignment (
			agency_id,
			vehicle_id,
			feed_version_id,
			telemetry_event_id,
			service_date,
			route_id,
			trip_id,
			block_id,
			start_date,
			start_time,
			current_stop_sequence,
			shape_dist_traveled,
			state,
			confidence,
			assignment_source,
			reason_codes,
			degraded_state,
			score_details_json,
			active_from,
			manual_override_id
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18::jsonb, $19, $20
		)
		RETURNING id
	`,
		assignment.AgencyID,
		assignment.VehicleID,
		nilIfEmpty(assignment.FeedVersionID),
		nilIfZero(assignment.TelemetryEventID),
		nilIfEmpty(assignment.ServiceDate),
		nilIfEmpty(assignment.RouteID),
		nilIfEmpty(assignment.TripID),
		nilIfEmpty(assignment.BlockID),
		nilIfEmpty(assignment.StartDate),
		nilIfEmpty(assignment.StartTime),
		nilIfZeroInt(assignment.CurrentStopSequence),
		assignment.ShapeDistTraveled,
		assignment.State,
		assignment.Confidence,
		assignment.AssignmentSource,
		assignment.ReasonCodes,
		assignment.DegradedState,
		string(scoreDetails),
		assignment.ActiveFrom,
		nilIfZero(assignment.ManualOverrideID),
	).Scan(&assignment.ID)
	if err != nil {
		return Assignment{}, fmt.Errorf("insert assignment: %w", err)
	}

	for _, incident := range incidents {
		if incident.Type == "" {
			continue
		}
		details, err := json.Marshal(incident.Details)
		if err != nil {
			return Assignment{}, fmt.Errorf("marshal incident details: %w", err)
		}
		severity := incident.Severity
		if severity == "" {
			severity = "warning"
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO incident (
				agency_id,
				incident_type,
				severity,
				route_id,
				vehicle_id,
				trip_id,
				vehicle_trip_assignment_id,
				details_json
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb)
		`,
			assignment.AgencyID,
			incident.Type,
			severity,
			nilIfEmpty(firstNonEmpty(incident.RouteID, assignment.RouteID)),
			firstNonEmpty(incident.VehicleID, assignment.VehicleID),
			nilIfEmpty(firstNonEmpty(incident.TripID, assignment.TripID)),
			assignment.ID,
			string(details),
		); err != nil {
			return Assignment{}, fmt.Errorf("insert incident: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Assignment{}, fmt.Errorf("commit assignment transaction: %w", err)
	}
	return assignment, nil
}

func currentAssignmentInTx(ctx context.Context, tx pgx.Tx, agencyID string, vehicleID string) (*Assignment, error) {
	row := tx.QueryRow(ctx, `
		SELECT id, agency_id, vehicle_id, feed_version_id, telemetry_event_id, service_date, route_id, trip_id, block_id, start_date, start_time,
		       current_stop_sequence, shape_dist_traveled, state, confidence, assignment_source, reason_codes, degraded_state, score_details_json,
		       manual_override_id, active_from
		FROM vehicle_trip_assignment
		WHERE agency_id = $1
		  AND vehicle_id = $2
		  AND active_to IS NULL
		ORDER BY active_from DESC, id DESC
		LIMIT 1
	`, agencyID, vehicleID)
	assignment, err := scanAssignment(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query current assignment for save: %w", err)
	}
	return &assignment, nil
}

func repeatedDegradedAssignment(current Assignment, next Assignment) bool {
	if current.State != StateUnknown || next.State != StateUnknown {
		return false
	}
	if current.DegradedState != next.DegradedState {
		return false
	}
	return equalStringSlices(current.ReasonCodes, next.ReasonCodes)
}

func equalStringSlices(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

type assignmentScanner interface {
	Scan(dest ...any) error
}

func scanAssignment(row assignmentScanner) (Assignment, error) {
	var assignment Assignment
	var feedVersionID, serviceDate, routeID, tripID, blockID, startDate, startTime sql.NullString
	var telemetryEventID, currentStopSequence, manualOverrideID sql.NullInt64
	var shapeDist sql.NullFloat64
	var scoreDetails []byte
	if err := row.Scan(
		&assignment.ID,
		&assignment.AgencyID,
		&assignment.VehicleID,
		&feedVersionID,
		&telemetryEventID,
		&serviceDate,
		&routeID,
		&tripID,
		&blockID,
		&startDate,
		&startTime,
		&currentStopSequence,
		&shapeDist,
		&assignment.State,
		&assignment.Confidence,
		&assignment.AssignmentSource,
		&assignment.ReasonCodes,
		&assignment.DegradedState,
		&scoreDetails,
		&manualOverrideID,
		&assignment.ActiveFrom,
	); err != nil {
		return Assignment{}, err
	}
	assignment.FeedVersionID = feedVersionID.String
	assignment.ServiceDate = serviceDate.String
	assignment.RouteID = routeID.String
	assignment.TripID = tripID.String
	assignment.BlockID = blockID.String
	assignment.StartDate = startDate.String
	assignment.StartTime = startTime.String
	if telemetryEventID.Valid {
		assignment.TelemetryEventID = telemetryEventID.Int64
	}
	if currentStopSequence.Valid {
		assignment.CurrentStopSequence = int(currentStopSequence.Int64)
	}
	if shapeDist.Valid {
		assignment.ShapeDistTraveled = shapeDist.Float64
	}
	if manualOverrideID.Valid {
		assignment.ManualOverrideID = manualOverrideID.Int64
	}
	if len(scoreDetails) > 0 {
		_ = json.Unmarshal(scoreDetails, &assignment.ScoreDetails)
	}
	return assignment, nil
}

func nilIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func nilIfZero(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

func nilIfZeroInt(value int) any {
	if value == 0 {
		return nil
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
