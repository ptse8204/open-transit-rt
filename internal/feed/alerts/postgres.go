package alerts

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresHealthRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresHealthRepository(pool *pgxpool.Pool) *PostgresHealthRepository {
	return &PostgresHealthRepository{pool: pool}
}

func (r *PostgresHealthRepository) SaveAlertsSnapshot(ctx context.Context, record HealthRecord) error {
	details := map[string]any{
		"diagnostics_status":     record.Status,
		"diagnostics_reason":     record.Reason,
		"active_feed_version_id": record.ActiveFeedVersionID,
		"alerts_output":          record.AlertsOutput,
	}
	payload, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("marshal alerts health details: %w", err)
	}
	var coverage *float64
	if record.AlertsOutput > 0 {
		value := 100.0
		coverage = &value
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
		VALUES ($1, 'alerts', $2, true, $3, $4::jsonb)
	`, record.AgencyID, record.SnapshotAt, coverage, string(payload))
	if err != nil {
		return fmt.Errorf("insert alerts health snapshot: %w", err)
	}
	return nil
}
