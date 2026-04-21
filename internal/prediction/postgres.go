package prediction

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDiagnosticsRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresDiagnosticsRepository(pool *pgxpool.Pool) *PostgresDiagnosticsRepository {
	return &PostgresDiagnosticsRepository{pool: pool}
}

func (r *PostgresDiagnosticsRepository) SaveTripUpdatesDiagnostics(ctx context.Context, record DiagnosticsRecord) (DiagnosticsPersistenceResult, error) {
	details := map[string]any{
		"adapter_name":                    record.AdapterName,
		"diagnostics_status":              record.Status,
		"diagnostics_reason":              record.Reason,
		"active_feed_version_id":          record.ActiveFeedVersionID,
		"input_counts":                    record.InputCounts,
		"vehicle_positions_url":           record.VehiclePositionsURL,
		"diagnostics_persistence_outcome": record.DiagnosticsPersistenceState,
	}
	payload, err := json.Marshal(details)
	if err != nil {
		return DiagnosticsPersistenceResult{}, fmt.Errorf("marshal trip updates diagnostics: %w", err)
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO feed_health_snapshot (
			agency_id,
			feed_type,
			snapshot_at,
			endpoint_available,
			coverage_percent,
			details_json
		)
		VALUES ($1, 'trip_updates', $2, true, $3, $4::jsonb)
	`, record.AgencyID, record.SnapshotAt, coveragePercent(record.InputCounts), string(payload))
	if err != nil {
		return DiagnosticsPersistenceResult{}, fmt.Errorf("insert trip updates diagnostics: %w", err)
	}
	return DiagnosticsPersistenceResult{Stored: true}, nil
}

func coveragePercent(counts InputCounts) *float64 {
	if counts.AssignmentRows <= 0 {
		return nil
	}
	value := float64(counts.TripUpdatesOutput) / float64(counts.AssignmentRows) * 100
	return &value
}
