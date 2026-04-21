package alerts

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

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) UpsertAlert(ctx context.Context, input UpsertInput) (Alert, error) {
	if input.AgencyID == "" {
		return Alert{}, fmt.Errorf("agency_id is required")
	}
	if input.AlertKey == "" {
		return Alert{}, fmt.Errorf("alert_key is required")
	}
	if input.HeaderText == "" {
		return Alert{}, fmt.Errorf("header_text is required")
	}
	if input.Now.IsZero() {
		input.Now = time.Now().UTC()
	}
	if input.SourceType == "" {
		input.SourceType = SourceOperator
	}
	if input.Cause == "" {
		input.Cause = "unknown_cause"
	}
	if input.Effect == "" {
		input.Effect = "unknown_effect"
	}
	actor := input.ActorID
	if actor == "" {
		actor = "system"
	}
	status := StatusDraft
	var publishedAt any
	if input.Publish {
		status = StatusPublished
		publishedAt = input.Now
	}
	metadata, err := json.Marshal(input.Metadata)
	if err != nil {
		return Alert{}, fmt.Errorf("marshal alert metadata: %w", err)
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Alert{}, fmt.Errorf("begin alert upsert: %w", err)
	}
	defer tx.Rollback(ctx)

	var oldValue map[string]any
	existing, exists, err := r.alertByKeyTx(ctx, tx, input.AgencyID, input.AlertKey)
	if err != nil {
		return Alert{}, err
	}
	if exists {
		oldValue = alertAuditValue(existing)
	}

	alert, err := scanAlert(tx.QueryRow(ctx, `
		INSERT INTO service_alert (
			agency_id, alert_key, status, cause, effect, header_text, description_text, url,
			active_start, active_end, feed_version_id, source_type, source_id, metadata_json,
			created_by, updated_by, published_at, archived_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NULLIF($11, ''), $12, $13, $14::jsonb,
			$15, $15, $16, NULL, $17)
		ON CONFLICT (agency_id, alert_key) DO UPDATE
		SET status = EXCLUDED.status,
		    cause = EXCLUDED.cause,
		    effect = EXCLUDED.effect,
		    header_text = EXCLUDED.header_text,
		    description_text = EXCLUDED.description_text,
		    url = EXCLUDED.url,
		    active_start = EXCLUDED.active_start,
		    active_end = EXCLUDED.active_end,
		    feed_version_id = EXCLUDED.feed_version_id,
		    source_type = EXCLUDED.source_type,
		    source_id = EXCLUDED.source_id,
		    metadata_json = EXCLUDED.metadata_json,
		    updated_by = EXCLUDED.updated_by,
		    published_at = COALESCE(EXCLUDED.published_at, service_alert.published_at),
		    archived_at = NULL,
		    updated_at = EXCLUDED.updated_at
		RETURNING id, agency_id, alert_key, status, cause, effect, header_text,
		          description_text, url, active_start, active_end, feed_version_id,
		          source_type, source_id, metadata_json, created_by, updated_by,
		          published_at, archived_at, created_at, updated_at
	`, input.AgencyID, input.AlertKey, status, input.Cause, input.Effect, input.HeaderText,
		nullString(input.DescriptionText), nullString(input.URL), input.ActiveStart, input.ActiveEnd,
		input.FeedVersionID, input.SourceType, nullString(input.SourceID), string(metadata), actor,
		publishedAt, input.Now))
	if err != nil {
		return Alert{}, fmt.Errorf("upsert service alert: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM service_alert_informed_entity WHERE service_alert_id = $1`, alert.ID); err != nil {
		return Alert{}, fmt.Errorf("delete old alert informed entities: %w", err)
	}
	for _, entity := range input.Entities {
		entityMetadata, err := json.Marshal(entity.Metadata)
		if err != nil {
			return Alert{}, fmt.Errorf("marshal informed entity metadata: %w", err)
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO service_alert_informed_entity (
				service_alert_id, agency_id, route_id, stop_id, trip_id, start_date,
				start_time, metadata_json
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb)
		`, alert.ID, input.AgencyID, nullString(entity.RouteID), nullString(entity.StopID),
			nullString(entity.TripID), nullString(entity.StartDate), nullString(entity.StartTime),
			string(entityMetadata)); err != nil {
			return Alert{}, fmt.Errorf("insert alert informed entity: %w", err)
		}
	}

	alert.Entities, err = r.listEntitiesTx(ctx, tx, alert.ID)
	if err != nil {
		return Alert{}, err
	}
	action := "service_alert.create"
	if exists {
		action = "service_alert.update"
	}
	if input.Publish {
		action = "service_alert.publish"
	}
	if err := insertAuditLog(ctx, tx, input.AgencyID, actor, action, "service_alert", strconv.FormatInt(alert.ID, 10), oldValue, alertAuditValue(alert), input.AlertKey); err != nil {
		return Alert{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return Alert{}, fmt.Errorf("commit alert upsert: %w", err)
	}
	return alert, nil
}

func (r *PostgresRepository) PublishAlert(ctx context.Context, agencyID string, alertID int64, actorID string, at time.Time) (Alert, error) {
	if at.IsZero() {
		at = time.Now().UTC()
	}
	if actorID == "" {
		actorID = "system"
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Alert{}, fmt.Errorf("begin alert publish: %w", err)
	}
	defer tx.Rollback(ctx)
	old, err := r.alertByIDTx(ctx, tx, agencyID, alertID)
	if err != nil {
		return Alert{}, err
	}
	alert, err := scanAlert(tx.QueryRow(ctx, `
		UPDATE service_alert
		SET status = 'published',
		    published_at = COALESCE(published_at, $3),
		    archived_at = NULL,
		    updated_by = $4,
		    updated_at = $3
		WHERE agency_id = $1 AND id = $2
		RETURNING id, agency_id, alert_key, status, cause, effect, header_text,
		          description_text, url, active_start, active_end, feed_version_id,
		          source_type, source_id, metadata_json, created_by, updated_by,
		          published_at, archived_at, created_at, updated_at
	`, agencyID, alertID, at, actorID))
	if err != nil {
		return Alert{}, fmt.Errorf("publish service alert: %w", err)
	}
	alert.Entities, err = r.listEntitiesTx(ctx, tx, alert.ID)
	if err != nil {
		return Alert{}, err
	}
	if err := insertAuditLog(ctx, tx, agencyID, actorID, "service_alert.publish", "service_alert", strconv.FormatInt(alert.ID, 10), alertAuditValue(old), alertAuditValue(alert), alert.AlertKey); err != nil {
		return Alert{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return Alert{}, fmt.Errorf("commit alert publish: %w", err)
	}
	return alert, nil
}

func (r *PostgresRepository) ArchiveAlert(ctx context.Context, agencyID string, alertID int64, actorID string, reason string, at time.Time) error {
	if at.IsZero() {
		at = time.Now().UTC()
	}
	if actorID == "" {
		actorID = "system"
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin alert archive: %w", err)
	}
	defer tx.Rollback(ctx)
	old, err := r.alertByIDTx(ctx, tx, agencyID, alertID)
	if err != nil {
		return err
	}
	alert, err := scanAlert(tx.QueryRow(ctx, `
		UPDATE service_alert
		SET status = 'archived',
		    archived_at = $3,
		    updated_by = $4,
		    updated_at = $3
		WHERE agency_id = $1 AND id = $2
		RETURNING id, agency_id, alert_key, status, cause, effect, header_text,
		          description_text, url, active_start, active_end, feed_version_id,
		          source_type, source_id, metadata_json, created_by, updated_by,
		          published_at, archived_at, created_at, updated_at
	`, agencyID, alertID, at, actorID))
	if err != nil {
		return fmt.Errorf("archive service alert: %w", err)
	}
	if err := insertAuditLog(ctx, tx, agencyID, actorID, "service_alert.archive", "service_alert", strconv.FormatInt(alertID, 10), alertAuditValue(old), alertAuditValue(alert), reason); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit alert archive: %w", err)
	}
	return nil
}

func (r *PostgresRepository) ListAlerts(ctx context.Context, filter ListFilter) ([]Alert, error) {
	if filter.Limit <= 0 {
		filter.Limit = 100
	}
	if filter.At.IsZero() {
		filter.At = time.Now().UTC()
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, agency_id, alert_key, status, cause, effect, header_text,
		       description_text, url, active_start, active_end, feed_version_id,
		       source_type, source_id, metadata_json, created_by, updated_by,
		       published_at, archived_at, created_at, updated_at
		FROM service_alert
		WHERE agency_id = $1
		  AND ($2 = '' OR status = $2)
		  AND (NOT $3 OR status = 'published')
		  AND (NOT $3 OR active_start IS NULL OR active_start <= $4)
		  AND (NOT $3 OR active_end IS NULL OR active_end >= $4)
		ORDER BY alert_key, id
		LIMIT $5
	`, filter.AgencyID, filter.Status, filter.PublishedOnly, filter.At, filter.Limit)
	if err != nil {
		return nil, fmt.Errorf("query service alerts: %w", err)
	}
	defer rows.Close()
	var alerts []Alert
	for rows.Next() {
		alert, err := scanAlert(rows)
		if err != nil {
			return nil, fmt.Errorf("scan service alert: %w", err)
		}
		alerts = append(alerts, alert)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate service alerts: %w", err)
	}
	for i := range alerts {
		entities, err := r.listEntities(ctx, alerts[i].ID)
		if err != nil {
			return nil, err
		}
		alerts[i].Entities = entities
	}
	return alerts, nil
}

func (r *PostgresRepository) alertByIDTx(ctx context.Context, tx pgx.Tx, agencyID string, alertID int64) (Alert, error) {
	alert, err := scanAlert(tx.QueryRow(ctx, `
		SELECT id, agency_id, alert_key, status, cause, effect, header_text,
		       description_text, url, active_start, active_end, feed_version_id,
		       source_type, source_id, metadata_json, created_by, updated_by,
		       published_at, archived_at, created_at, updated_at
		FROM service_alert
		WHERE agency_id = $1 AND id = $2
	`, agencyID, alertID))
	if err != nil {
		return Alert{}, fmt.Errorf("query service alert: %w", err)
	}
	alert.Entities, err = r.listEntitiesTx(ctx, tx, alert.ID)
	if err != nil {
		return Alert{}, err
	}
	return alert, nil
}

func (r *PostgresRepository) alertByKeyTx(ctx context.Context, tx pgx.Tx, agencyID string, alertKey string) (Alert, bool, error) {
	alert, err := scanAlert(tx.QueryRow(ctx, `
		SELECT id, agency_id, alert_key, status, cause, effect, header_text,
		       description_text, url, active_start, active_end, feed_version_id,
		       source_type, source_id, metadata_json, created_by, updated_by,
		       published_at, archived_at, created_at, updated_at
		FROM service_alert
		WHERE agency_id = $1 AND alert_key = $2
	`, agencyID, alertKey))
	if err != nil {
		if err == pgx.ErrNoRows {
			return Alert{}, false, nil
		}
		return Alert{}, false, fmt.Errorf("query service alert by key: %w", err)
	}
	alert.Entities, err = r.listEntitiesTx(ctx, tx, alert.ID)
	if err != nil {
		return Alert{}, false, err
	}
	return alert, true, nil
}

func (r *PostgresRepository) listEntities(ctx context.Context, alertID int64) ([]InformedEntity, error) {
	rows, err := r.pool.Query(ctx, entityQuery(), alertID)
	if err != nil {
		return nil, fmt.Errorf("query alert informed entities: %w", err)
	}
	defer rows.Close()
	return scanEntities(rows)
}

func (r *PostgresRepository) listEntitiesTx(ctx context.Context, tx pgx.Tx, alertID int64) ([]InformedEntity, error) {
	rows, err := tx.Query(ctx, entityQuery(), alertID)
	if err != nil {
		return nil, fmt.Errorf("query alert informed entities: %w", err)
	}
	defer rows.Close()
	return scanEntities(rows)
}

func entityQuery() string {
	return `
		SELECT id, service_alert_id, agency_id, route_id, stop_id, trip_id,
		       start_date, start_time, metadata_json, created_at
		FROM service_alert_informed_entity
		WHERE service_alert_id = $1
		ORDER BY route_id NULLS LAST, stop_id NULLS LAST, trip_id NULLS LAST, start_date NULLS LAST, start_time NULLS LAST, id
	`
}

func scanEntities(rows pgx.Rows) ([]InformedEntity, error) {
	var entities []InformedEntity
	for rows.Next() {
		var entity InformedEntity
		var routeID, stopID, tripID, startDate, startTime sql.NullString
		var metadataBytes []byte
		if err := rows.Scan(&entity.ID, &entity.ServiceAlertID, &entity.AgencyID, &routeID, &stopID, &tripID, &startDate, &startTime, &metadataBytes, &entity.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan alert informed entity: %w", err)
		}
		entity.RouteID = routeID.String
		entity.StopID = stopID.String
		entity.TripID = tripID.String
		entity.StartDate = startDate.String
		entity.StartTime = startTime.String
		entity.Metadata = map[string]any{}
		if len(metadataBytes) > 0 {
			_ = json.Unmarshal(metadataBytes, &entity.Metadata)
		}
		entities = append(entities, entity)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate alert informed entities: %w", err)
	}
	return entities, nil
}

type alertScanner interface {
	Scan(dest ...any) error
}

func scanAlert(row alertScanner) (Alert, error) {
	var alert Alert
	var description, urlValue, feedVersionID, sourceID sql.NullString
	var activeStart, activeEnd, publishedAt, archivedAt sql.NullTime
	var metadataBytes []byte
	if err := row.Scan(
		&alert.ID,
		&alert.AgencyID,
		&alert.AlertKey,
		&alert.Status,
		&alert.Cause,
		&alert.Effect,
		&alert.HeaderText,
		&description,
		&urlValue,
		&activeStart,
		&activeEnd,
		&feedVersionID,
		&alert.SourceType,
		&sourceID,
		&metadataBytes,
		&alert.CreatedBy,
		&alert.UpdatedBy,
		&publishedAt,
		&archivedAt,
		&alert.CreatedAt,
		&alert.UpdatedAt,
	); err != nil {
		return Alert{}, err
	}
	alert.DescriptionText = description.String
	alert.URL = urlValue.String
	alert.FeedVersionID = feedVersionID.String
	alert.SourceID = sourceID.String
	if activeStart.Valid {
		t := activeStart.Time
		alert.ActiveStart = &t
	}
	if activeEnd.Valid {
		t := activeEnd.Time
		alert.ActiveEnd = &t
	}
	if publishedAt.Valid {
		t := publishedAt.Time
		alert.PublishedAt = &t
	}
	if archivedAt.Valid {
		t := archivedAt.Time
		alert.ArchivedAt = &t
	}
	alert.Metadata = map[string]any{}
	if len(metadataBytes) > 0 {
		_ = json.Unmarshal(metadataBytes, &alert.Metadata)
	}
	return alert, nil
}

func alertAuditValue(alert Alert) map[string]any {
	return map[string]any{
		"id":          alert.ID,
		"alert_key":   alert.AlertKey,
		"status":      alert.Status,
		"cause":       alert.Cause,
		"effect":      alert.Effect,
		"header_text": alert.HeaderText,
		"source_type": alert.SourceType,
		"source_id":   alert.SourceID,
	}
}

func insertAuditLog(ctx context.Context, tx pgx.Tx, agencyID string, actorID string, action string, entityType string, entityID string, oldValue any, newValue any, reason string) error {
	oldJSON, err := marshalNullable(oldValue)
	if err != nil {
		return err
	}
	newJSON, err := marshalNullable(newValue)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO audit_log (agency_id, actor_id, action, entity_type, entity_id, old_value_json, new_value_json, reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, agencyID, actorID, action, entityType, entityID, oldJSON, newJSON, nullString(reason))
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

func marshalNullable(value any) (any, error) {
	if value == nil {
		return nil, nil
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal audit value: %w", err)
	}
	return string(payload), nil
}

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}
