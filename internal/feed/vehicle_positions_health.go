package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VehiclePositionsHealthRecord struct {
	AgencyID                string
	SnapshotAt              time.Time
	EndpointAvailable       bool
	FreshnessSeconds        *float64
	GenerationLatencyMS     *float64
	MatchedVehiclePercent   *float64
	VehiclesInSnapshot      int
	VehiclesInProtobuf      int
	TripDescriptors         int
	StaleVehicles           int
	UnmatchedVehicles       int
	Truncated               bool
	VehicleLimit            int
	LatestTelemetryRowsRead int
}

type VehiclePositionsHealthRepository struct {
	pool *pgxpool.Pool
}

func NewVehiclePositionsHealthRepository(pool *pgxpool.Pool) *VehiclePositionsHealthRepository {
	return &VehiclePositionsHealthRepository{pool: pool}
}

func (r *VehiclePositionsHealthRepository) SaveVehiclePositionsHealth(ctx context.Context, record VehiclePositionsHealthRecord) error {
	if r == nil || r.pool == nil {
		return fmt.Errorf("vehicle positions health repository is not configured")
	}
	if record.AgencyID == "" {
		return fmt.Errorf("agency_id is required")
	}
	snapshotAt := record.SnapshotAt.UTC()
	if snapshotAt.IsZero() {
		snapshotAt = time.Now().UTC()
	}
	details := map[string]any{
		"diagnostics_status":         "observed",
		"diagnostics_reason":         "vehicle_positions_feed_generated",
		"vehicles_in_snapshot":       boundedNonNegative(record.VehiclesInSnapshot),
		"vehicles_in_protobuf":       boundedNonNegative(record.VehiclesInProtobuf),
		"trip_descriptors":           boundedNonNegative(record.TripDescriptors),
		"stale_vehicles":             boundedNonNegative(record.StaleVehicles),
		"unmatched_vehicles":         boundedNonNegative(record.UnmatchedVehicles),
		"truncated":                  record.Truncated,
		"vehicle_limit":              boundedNonNegative(record.VehicleLimit),
		"latest_telemetry_rows_read": boundedNonNegative(record.LatestTelemetryRowsRead),
	}
	payload, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("marshal vehicle positions health details: %w", err)
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO feed_health_snapshot (
			agency_id,
			feed_type,
			snapshot_at,
			endpoint_available,
			freshness_seconds,
			generation_latency_ms,
			matched_vehicle_percent,
			details_json
		)
		VALUES ($1, 'vehicle_positions', $2, $3, $4, $5, $6, $7::jsonb)
	`, record.AgencyID, snapshotAt, record.EndpointAvailable, record.FreshnessSeconds,
		record.GenerationLatencyMS, record.MatchedVehiclePercent, string(payload))
	if err != nil {
		return fmt.Errorf("insert vehicle positions health snapshot: %w", err)
	}
	return nil
}

func HealthRecordFromVehiclePositionsSnapshot(snapshot VehiclePositionsSnapshot, generationLatency time.Duration) VehiclePositionsHealthRecord {
	included := 0
	publishedTrips := 0
	stale := 0
	unmatched := 0
	var newest *time.Time
	for _, vehicle := range snapshot.Vehicles {
		if vehicle.IncludedInProtobuf {
			included++
		}
		if vehicle.TripDescriptorPublished {
			publishedTrips++
		}
		if vehicle.TripDescriptorOmissionReason == TripDescriptorOmissionStaleTelemetry || vehicle.TripDescriptorOmissionReason == TripDescriptorOmissionSuppressedStaleTelemetry {
			stale++
		}
		if vehicle.TripDescriptorOmissionReason == TripDescriptorOmissionNoAssignment {
			unmatched++
		}
		observed := vehicle.TelemetryEvent.Timestamp.UTC()
		if !observed.IsZero() && (newest == nil || observed.After(*newest)) {
			t := observed
			newest = &t
		}
	}
	var freshness *float64
	if newest != nil {
		seconds := snapshot.GeneratedAt.Sub(*newest).Seconds()
		if seconds < 0 {
			seconds = 0
		}
		freshness = &seconds
	}
	latency := float64(generationLatency.Milliseconds())
	matched := matchedVehiclePercent(snapshot.Vehicles)
	return VehiclePositionsHealthRecord{
		AgencyID:                snapshot.AgencyID,
		SnapshotAt:              snapshot.GeneratedAt,
		EndpointAvailable:       true,
		FreshnessSeconds:        freshness,
		GenerationLatencyMS:     &latency,
		MatchedVehiclePercent:   matched,
		VehiclesInSnapshot:      snapshot.VehiclesInSnapshot,
		VehiclesInProtobuf:      included,
		TripDescriptors:         publishedTrips,
		StaleVehicles:           stale,
		UnmatchedVehicles:       unmatched,
		Truncated:               snapshot.Truncated,
		VehicleLimit:            snapshot.VehicleLimit,
		LatestTelemetryRowsRead: snapshot.LatestTelemetryRowsRead,
	}
}

func matchedVehiclePercent(vehicles []VehicleSnapshot) *float64 {
	if len(vehicles) == 0 {
		return nil
	}
	matched := 0
	for _, vehicle := range vehicles {
		if vehicle.TripDescriptorPublished {
			matched++
		}
	}
	value := float64(matched) / float64(len(vehicles)) * 100
	return &value
}

func boundedNonNegative(value int) int {
	if value < 0 {
		return 0
	}
	if value > 1000000 {
		return 1000000
	}
	return value
}
