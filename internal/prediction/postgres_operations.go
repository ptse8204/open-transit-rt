package prediction

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOperationsRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresOperationsRepository(pool *pgxpool.Pool) *PostgresOperationsRepository {
	return &PostgresOperationsRepository{pool: pool}
}

func (r *PostgresOperationsRepository) ListActivePredictionOverrides(ctx context.Context, agencyID string, at time.Time) ([]OverrideRecord, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, agency_id, vehicle_id, override_type, route_id, trip_id, start_date, start_time, state, expires_at, cleared_at, reason, created_by, created_at
		FROM manual_override
		WHERE agency_id = $1
		  AND cleared_at IS NULL
		  AND (expires_at IS NULL OR expires_at > $2)
		ORDER BY created_at DESC, id DESC
	`, agencyID, at)
	if err != nil {
		return nil, fmt.Errorf("query active prediction overrides: %w", err)
	}
	defer rows.Close()

	var overrides []OverrideRecord
	for rows.Next() {
		override, err := scanOverride(rows)
		if err != nil {
			return nil, fmt.Errorf("scan prediction override: %w", err)
		}
		overrides = append(overrides, override)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate prediction overrides: %w", err)
	}
	return overrides, nil
}

func (r *PostgresOperationsRepository) CreatePredictionOverride(ctx context.Context, input OverrideInput) (OverrideRecord, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return OverrideRecord{}, fmt.Errorf("begin create prediction override: %w", err)
	}
	defer tx.Rollback(ctx)

	override, err := insertOverride(ctx, tx, input)
	if err != nil {
		return OverrideRecord{}, err
	}
	if err := insertAuditLog(ctx, tx, input.AgencyID, input.ActorID, "prediction_override.create", "manual_override", strconv.FormatInt(override.ID, 10), nil, override, input.Reason); err != nil {
		return OverrideRecord{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return OverrideRecord{}, fmt.Errorf("commit create prediction override: %w", err)
	}
	return override, nil
}

func (r *PostgresOperationsRepository) ReplacePredictionOverride(ctx context.Context, input OverrideInput) (OverrideRecord, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return OverrideRecord{}, fmt.Errorf("begin replace prediction override: %w", err)
	}
	defer tx.Rollback(ctx)

	at := input.Now
	if at.IsZero() {
		at = time.Now().UTC()
	}
	rows, err := tx.Query(ctx, `
		SELECT id, agency_id, vehicle_id, override_type, route_id, trip_id, start_date, start_time, state, expires_at, cleared_at, reason, created_by, created_at
		FROM manual_override
		WHERE agency_id = $1
		  AND vehicle_id = $2
		  AND cleared_at IS NULL
		  AND (expires_at IS NULL OR expires_at > $3)
		ORDER BY created_at DESC, id DESC
	`, input.AgencyID, input.VehicleID, at)
	if err != nil {
		return OverrideRecord{}, fmt.Errorf("query overrides for replace: %w", err)
	}
	var old []OverrideRecord
	for rows.Next() {
		override, err := scanOverride(rows)
		if err != nil {
			rows.Close()
			return OverrideRecord{}, fmt.Errorf("scan override for replace: %w", err)
		}
		old = append(old, override)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return OverrideRecord{}, fmt.Errorf("iterate overrides for replace: %w", err)
	}
	rows.Close()

	if _, err := tx.Exec(ctx, `
		UPDATE manual_override
		SET cleared_at = $3
		WHERE agency_id = $1
		  AND vehicle_id = $2
		  AND cleared_at IS NULL
		  AND (expires_at IS NULL OR expires_at > $3)
	`, input.AgencyID, input.VehicleID, at); err != nil {
		return OverrideRecord{}, fmt.Errorf("clear replaced prediction overrides: %w", err)
	}

	override, err := insertOverride(ctx, tx, input)
	if err != nil {
		return OverrideRecord{}, err
	}
	if err := insertAuditLog(ctx, tx, input.AgencyID, input.ActorID, "prediction_override.replace", "manual_override", strconv.FormatInt(override.ID, 10), old, override, input.Reason); err != nil {
		return OverrideRecord{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return OverrideRecord{}, fmt.Errorf("commit replace prediction override: %w", err)
	}
	return override, nil
}

func (r *PostgresOperationsRepository) ClearPredictionOverride(ctx context.Context, agencyID string, overrideID int64, actorID string, reason string, at time.Time) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin clear prediction override: %w", err)
	}
	defer tx.Rollback(ctx)

	if at.IsZero() {
		at = time.Now().UTC()
	}
	row := tx.QueryRow(ctx, `
		SELECT id, agency_id, vehicle_id, override_type, route_id, trip_id, start_date, start_time, state, expires_at, cleared_at, reason, created_by, created_at
		FROM manual_override
		WHERE agency_id = $1
		  AND id = $2
	`, agencyID, overrideID)
	old, err := scanOverride(row)
	if err != nil {
		return fmt.Errorf("query prediction override for clear: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		UPDATE manual_override
		SET cleared_at = $3
		WHERE agency_id = $1
		  AND id = $2
	`, agencyID, overrideID, at); err != nil {
		return fmt.Errorf("clear prediction override: %w", err)
	}
	next := old
	next.ClearedAt = &at
	if err := insertAuditLog(ctx, tx, agencyID, actorID, "prediction_override.clear", "manual_override", strconv.FormatInt(overrideID, 10), old, next, reason); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit clear prediction override: %w", err)
	}
	return nil
}

func (r *PostgresOperationsRepository) SavePredictionReviewItems(ctx context.Context, items []ReviewItem) error {
	for _, item := range items {
		if item.AgencyID == "" || item.Reason == "" {
			continue
		}
		status := item.Status
		if status == "" {
			status = ReviewStatusOpen
		}
		severity := item.Severity
		if severity == "" {
			severity = "warning"
		}
		details := cloneDetails(item.Details)
		details["prediction_review"] = true
		details["reason"] = item.Reason
		details["status"] = status
		details["snapshot_at"] = item.SnapshotAt.UTC().Format(time.RFC3339)
		if item.StartDate != "" {
			details["start_date"] = item.StartDate
		}
		if item.StartTime != "" {
			details["start_time"] = item.StartTime
		}
		payload, err := json.Marshal(details)
		if err != nil {
			return fmt.Errorf("marshal prediction review details: %w", err)
		}
		_, err = r.pool.Exec(ctx, `
			INSERT INTO incident (
				agency_id,
				incident_type,
				severity,
				route_id,
				vehicle_id,
				trip_id,
				status,
				details_json
			)
			VALUES ($1, 'prediction_review', $2, $3, $4, $5, $6, $7::jsonb)
		`, item.AgencyID, severity, nilIfEmpty(item.RouteID), nilIfEmpty(item.VehicleID), nilIfEmpty(item.TripID), status, string(payload))
		if err != nil {
			return fmt.Errorf("insert prediction review item: %w", err)
		}
	}
	return nil
}

func (r *PostgresOperationsRepository) ListPredictionReviewItems(ctx context.Context, filter ReviewFilter) ([]ReviewItem, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}
	status := filter.Status
	rows, err := r.pool.Query(ctx, `
		SELECT id, agency_id, created_at, severity, route_id, vehicle_id, trip_id, status, details_json, resolved_at
		FROM incident
		WHERE agency_id = $1
		  AND incident_type = 'prediction_review'
		  AND ($2 = '' OR status = $2)
		ORDER BY created_at DESC, id DESC
		LIMIT $3
	`, filter.AgencyID, string(status), limit)
	if err != nil {
		return nil, fmt.Errorf("query prediction review items: %w", err)
	}
	defer rows.Close()

	var items []ReviewItem
	for rows.Next() {
		item, err := scanReviewItem(rows)
		if err != nil {
			return nil, fmt.Errorf("scan prediction review item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate prediction review items: %w", err)
	}
	return items, nil
}

func (r *PostgresOperationsRepository) UpdatePredictionReviewStatus(ctx context.Context, agencyID string, reviewID int64, status ReviewStatus, actorID string, reason string, at time.Time) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin update prediction review status: %w", err)
	}
	defer tx.Rollback(ctx)

	if at.IsZero() {
		at = time.Now().UTC()
	}
	var oldStatus string
	if err := tx.QueryRow(ctx, `
		SELECT status
		FROM incident
		WHERE agency_id = $1
		  AND id = $2
		  AND incident_type = 'prediction_review'
	`, agencyID, reviewID).Scan(&oldStatus); err != nil {
		return fmt.Errorf("query prediction review for status update: %w", err)
	}
	var resolvedAt any
	if status == ReviewStatusResolved {
		resolvedAt = at
	}
	if _, err := tx.Exec(ctx, `
		UPDATE incident
		SET status = $3,
		    resolved_at = $4
		WHERE agency_id = $1
		  AND id = $2
		  AND incident_type = 'prediction_review'
	`, agencyID, reviewID, status, resolvedAt); err != nil {
		return fmt.Errorf("update prediction review status: %w", err)
	}
	oldValue := map[string]any{"status": oldStatus}
	newValue := map[string]any{"status": status}
	if err := insertAuditLog(ctx, tx, agencyID, actorID, "prediction_review.status_update", "incident", strconv.FormatInt(reviewID, 10), oldValue, newValue, reason); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit update prediction review status: %w", err)
	}
	return nil
}

type overrideScanner interface {
	Scan(dest ...any) error
}

func scanOverride(row overrideScanner) (OverrideRecord, error) {
	var override OverrideRecord
	var routeID, tripID, startDate, startTime, reason sql.NullString
	var expiresAt, clearedAt sql.NullTime
	if err := row.Scan(
		&override.ID,
		&override.AgencyID,
		&override.VehicleID,
		&override.OverrideType,
		&routeID,
		&tripID,
		&startDate,
		&startTime,
		&override.State,
		&expiresAt,
		&clearedAt,
		&reason,
		&override.CreatedBy,
		&override.CreatedAt,
	); err != nil {
		return OverrideRecord{}, err
	}
	override.RouteID = routeID.String
	override.TripID = tripID.String
	override.StartDate = startDate.String
	override.StartTime = startTime.String
	override.Reason = reason.String
	if expiresAt.Valid {
		override.ExpiresAt = &expiresAt.Time
	}
	if clearedAt.Valid {
		override.ClearedAt = &clearedAt.Time
	}
	return override, nil
}

type reviewScanner interface {
	Scan(dest ...any) error
}

func scanReviewItem(row reviewScanner) (ReviewItem, error) {
	var item ReviewItem
	var routeID, vehicleID, tripID sql.NullString
	var detailsBytes []byte
	var status string
	var resolvedAt sql.NullTime
	if err := row.Scan(
		&item.ID,
		&item.AgencyID,
		&item.CreatedAt,
		&item.Severity,
		&routeID,
		&vehicleID,
		&tripID,
		&status,
		&detailsBytes,
		&resolvedAt,
	); err != nil {
		return ReviewItem{}, err
	}
	item.RouteID = routeID.String
	item.VehicleID = vehicleID.String
	item.TripID = tripID.String
	item.Status = ReviewStatus(status)
	if resolvedAt.Valid {
		item.ResolvedAt = &resolvedAt.Time
	}
	if len(detailsBytes) > 0 {
		_ = json.Unmarshal(detailsBytes, &item.Details)
	}
	item.Reason, _ = item.Details["reason"].(string)
	item.StartDate, _ = item.Details["start_date"].(string)
	item.StartTime, _ = item.Details["start_time"].(string)
	if raw, _ := item.Details["snapshot_at"].(string); raw != "" {
		if parsed, err := time.Parse(time.RFC3339, raw); err == nil {
			item.SnapshotAt = parsed
		}
	}
	if item.SnapshotAt.IsZero() {
		item.SnapshotAt = item.CreatedAt
	}
	return item, nil
}

func insertOverride(ctx context.Context, tx pgx.Tx, input OverrideInput) (OverrideRecord, error) {
	if input.ActorID == "" {
		return OverrideRecord{}, fmt.Errorf("actor id is required for prediction override")
	}
	if input.Now.IsZero() {
		input.Now = time.Now().UTC()
	}
	var override OverrideRecord
	var routeID, tripID, startDate, startTime, reason sql.NullString
	var expiresAt, clearedAt sql.NullTime
	row := tx.QueryRow(ctx, `
		INSERT INTO manual_override (
			agency_id,
			vehicle_id,
			override_type,
			route_id,
			trip_id,
			start_date,
			start_time,
			state,
			expires_at,
			reason,
			created_by,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, agency_id, vehicle_id, override_type, route_id, trip_id, start_date, start_time, state, expires_at, cleared_at, reason, created_by, created_at
	`,
		input.AgencyID,
		input.VehicleID,
		input.OverrideType,
		nilIfEmpty(input.RouteID),
		nilIfEmpty(input.TripID),
		nilIfEmpty(input.StartDate),
		nilIfEmpty(input.StartTime),
		input.State,
		input.ExpiresAt,
		nilIfEmpty(input.Reason),
		input.ActorID,
		input.Now,
	)
	if err := row.Scan(
		&override.ID,
		&override.AgencyID,
		&override.VehicleID,
		&override.OverrideType,
		&routeID,
		&tripID,
		&startDate,
		&startTime,
		&override.State,
		&expiresAt,
		&clearedAt,
		&reason,
		&override.CreatedBy,
		&override.CreatedAt,
	); err != nil {
		return OverrideRecord{}, fmt.Errorf("insert prediction override: %w", err)
	}
	override.RouteID = routeID.String
	override.TripID = tripID.String
	override.StartDate = startDate.String
	override.StartTime = startTime.String
	override.Reason = reason.String
	if expiresAt.Valid {
		override.ExpiresAt = &expiresAt.Time
	}
	if clearedAt.Valid {
		override.ClearedAt = &clearedAt.Time
	}
	return override, nil
}

func insertAuditLog(ctx context.Context, tx pgx.Tx, agencyID string, actorID string, action string, entityType string, entityID string, oldValue any, newValue any, reason string) error {
	if actorID == "" {
		return fmt.Errorf("actor id is required for audit log")
	}
	oldPayload, err := json.Marshal(oldValue)
	if err != nil {
		return fmt.Errorf("marshal audit old value: %w", err)
	}
	newPayload, err := json.Marshal(newValue)
	if err != nil {
		return fmt.Errorf("marshal audit new value: %w", err)
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO audit_log (
			agency_id,
			actor_id,
			action,
			entity_type,
			entity_id,
			old_value_json,
			new_value_json,
			reason
		)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7::jsonb, $8)
	`, agencyID, actorID, action, entityType, nilIfEmpty(entityID), string(oldPayload), string(newPayload), nilIfEmpty(reason))
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

func cloneDetails(details map[string]any) map[string]any {
	cloned := make(map[string]any, len(details)+4)
	for key, value := range details {
		cloned[key] = value
	}
	return cloned
}

func nilIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}
