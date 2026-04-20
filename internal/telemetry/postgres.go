package telemetry

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
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

func (r *PostgresRepository) Store(ctx context.Context, event Event, payload json.RawMessage) (StoreResult, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return StoreResult{}, fmt.Errorf("begin telemetry transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var agencyExists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM agency WHERE id = $1)`, event.AgencyID).Scan(&agencyExists); err != nil {
		return StoreResult{}, fmt.Errorf("check agency: %w", err)
	}
	if !agencyExists {
		return StoreResult{}, ErrUnknownAgency
	}

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock($1::bigint)`, advisoryLockKey(event.AgencyID, event.VehicleID)); err != nil {
		return StoreResult{}, fmt.Errorf("lock telemetry vehicle stream: %w", err)
	}

	status := IngestStatusAccepted
	var latestObservedAt time.Time
	err = tx.QueryRow(ctx, `
		SELECT observed_at
		FROM telemetry_event
		WHERE agency_id = $1
		  AND vehicle_id = $2
		  AND ingest_status = 'accepted'
		ORDER BY observed_at DESC, id DESC
		LIMIT 1
	`, event.AgencyID, event.VehicleID).Scan(&latestObservedAt)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return StoreResult{}, fmt.Errorf("query latest accepted telemetry: %w", err)
	}
	if err == nil {
		switch {
		case event.Timestamp.Equal(latestObservedAt):
			status = IngestStatusDuplicate
		case event.Timestamp.Before(latestObservedAt):
			status = IngestStatusOutOfOrder
		}
	}

	stored := StoredEvent{
		Event:        event,
		IngestStatus: status,
		PayloadJSON:  append(json.RawMessage(nil), payload...),
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO telemetry_event (
			agency_id,
			device_id,
			vehicle_id,
			observed_at,
			lat,
			lon,
			geom,
			bearing,
			speed_mps,
			accuracy_m,
			trip_hint,
			payload_json,
			ingest_status
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			ST_SetSRID(ST_MakePoint($6, $5), 4326),
			$7, $8, $9, $10, $11::jsonb, $12
		)
		RETURNING id, received_at
	`,
		event.AgencyID,
		event.DeviceID,
		event.VehicleID,
		event.Timestamp,
		event.Lat,
		event.Lon,
		event.Bearing,
		event.SpeedMPS,
		event.AccuracyM,
		event.TripHint,
		string(payload),
		status,
	).Scan(&stored.ID, &stored.ReceivedAt)
	if err != nil {
		return StoreResult{}, fmt.Errorf("insert telemetry event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return StoreResult{}, fmt.Errorf("commit telemetry transaction: %w", err)
	}

	return StoreResult{StoredEvent: stored}, nil
}

func (r *PostgresRepository) LatestByVehicle(ctx context.Context, agencyID string, vehicleID string) (StoredEvent, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, agency_id, device_id, vehicle_id, observed_at, received_at, lat, lon, bearing, speed_mps, accuracy_m, trip_hint, payload_json, ingest_status
		FROM telemetry_event
		WHERE agency_id = $1
		  AND vehicle_id = $2
		  AND ingest_status = 'accepted'
		ORDER BY observed_at DESC, id DESC
		LIMIT 1
	`, agencyID, vehicleID)
	return scanStoredEvent(row)
}

func (r *PostgresRepository) ListLatestByAgency(ctx context.Context, agencyID string, limit int) ([]StoredEvent, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, agency_id, device_id, vehicle_id, observed_at, received_at, lat, lon, bearing, speed_mps, accuracy_m, trip_hint, payload_json, ingest_status
		FROM (
			SELECT DISTINCT ON (vehicle_id)
				id, agency_id, device_id, vehicle_id, observed_at, received_at, lat, lon, bearing, speed_mps, accuracy_m, trip_hint, payload_json, ingest_status
			FROM telemetry_event
			WHERE agency_id = $1
			  AND ingest_status = 'accepted'
			ORDER BY vehicle_id, observed_at DESC, id DESC
		) latest
		ORDER BY observed_at DESC, id DESC
		LIMIT $2
	`, agencyID, limit)
	if err != nil {
		return nil, fmt.Errorf("query latest agency telemetry: %w", err)
	}
	defer rows.Close()
	return scanStoredEvents(rows)
}

func (r *PostgresRepository) ListEvents(ctx context.Context, agencyID string, limit int) ([]StoredEvent, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, agency_id, device_id, vehicle_id, observed_at, received_at, lat, lon, bearing, speed_mps, accuracy_m, trip_hint, payload_json, ingest_status
		FROM telemetry_event
		WHERE agency_id = $1
		ORDER BY received_at DESC, id DESC
		LIMIT $2
	`, agencyID, limit)
	if err != nil {
		return nil, fmt.Errorf("query agency telemetry events: %w", err)
	}
	defer rows.Close()
	return scanStoredEvents(rows)
}

func advisoryLockKey(agencyID string, vehicleID string) int64 {
	sum := sha256.Sum256([]byte("telemetry_event\x00" + agencyID + "\x00" + vehicleID))
	return int64(binary.BigEndian.Uint64(sum[:8]))
}

type rowScanner interface {
	Scan(dest ...any) error
}

type rowsScanner interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}

func scanStoredEvent(row rowScanner) (StoredEvent, error) {
	var stored StoredEvent
	var payload []byte
	if err := row.Scan(
		&stored.ID,
		&stored.AgencyID,
		&stored.DeviceID,
		&stored.VehicleID,
		&stored.Timestamp,
		&stored.ReceivedAt,
		&stored.Lat,
		&stored.Lon,
		&stored.Bearing,
		&stored.SpeedMPS,
		&stored.AccuracyM,
		&stored.TripHint,
		&payload,
		&stored.IngestStatus,
	); err != nil {
		return StoredEvent{}, err
	}
	stored.PayloadJSON = append(json.RawMessage(nil), payload...)
	return stored, nil
}

func scanStoredEvents(rows rowsScanner) ([]StoredEvent, error) {
	var events []StoredEvent
	for rows.Next() {
		stored, err := scanStoredEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("scan telemetry event: %w", err)
		}
		events = append(events, stored)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate telemetry events: %w", err)
	}
	return events, nil
}
