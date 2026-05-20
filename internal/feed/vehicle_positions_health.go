package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VehiclePositionsHealthRecord struct {
	AgencyID                       string
	SnapshotAt                     time.Time
	EndpointAvailable              bool
	FreshnessSeconds               *float64
	GenerationLatencyMS            *float64
	MatchedVehiclePercent          *float64
	VehiclesInSnapshot             int
	VehiclesInProtobuf             int
	TripDescriptors                int
	StaleVehicles                  int
	UnmatchedVehicles              int
	SuppressedVehicles             int
	MissingScheduleContextVehicles int
	LowConfidenceVehicles          int
	AssignmentMismatches           int
	TripDescriptorOmissions        map[string]int
	UsefulnessStatus               string
	UsefulnessReasons              []string
	Truncated                      bool
	VehicleLimit                   int
	LatestTelemetryRowsRead        int
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
		"diagnostics_status":                "observed",
		"diagnostics_reason":                "vehicle_positions_feed_generated",
		"vehicles_in_snapshot":              boundedNonNegative(record.VehiclesInSnapshot),
		"vehicles_in_protobuf":              boundedNonNegative(record.VehiclesInProtobuf),
		"trip_descriptors":                  boundedNonNegative(record.TripDescriptors),
		"stale_vehicles":                    boundedNonNegative(record.StaleVehicles),
		"unmatched_vehicles":                boundedNonNegative(record.UnmatchedVehicles),
		"suppressed_vehicles":               boundedNonNegative(record.SuppressedVehicles),
		"missing_schedule_context_vehicles": boundedNonNegative(record.MissingScheduleContextVehicles),
		"low_confidence_vehicles":           boundedNonNegative(record.LowConfidenceVehicles),
		"assignment_mismatches":             boundedNonNegative(record.AssignmentMismatches),
		"trip_descriptor_omissions":         boundedCountMap(record.TripDescriptorOmissions),
		"usefulness_status":                 boundedHealthText(record.UsefulnessStatus),
		"usefulness_reasons":                boundedStringList(record.UsefulnessReasons),
		"truncated":                         record.Truncated,
		"vehicle_limit":                     boundedNonNegative(record.VehicleLimit),
		"latest_telemetry_rows_read":        boundedNonNegative(record.LatestTelemetryRowsRead),
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
	review := snapshot.ReviewSummary()
	var newest *time.Time
	for _, vehicle := range snapshot.Vehicles {
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
		AgencyID:                       snapshot.AgencyID,
		SnapshotAt:                     snapshot.GeneratedAt,
		EndpointAvailable:              true,
		FreshnessSeconds:               freshness,
		GenerationLatencyMS:            &latency,
		MatchedVehiclePercent:          matched,
		VehiclesInSnapshot:             snapshot.VehiclesInSnapshot,
		VehiclesInProtobuf:             review.VehiclesInProtobuf,
		TripDescriptors:                review.TripDescriptorsPublished,
		StaleVehicles:                  review.StaleVehicles,
		UnmatchedVehicles:              review.UnmatchedVehicles,
		SuppressedVehicles:             review.SuppressedVehicles,
		MissingScheduleContextVehicles: review.MissingScheduleContextVehicles,
		LowConfidenceVehicles:          review.LowConfidenceVehicles,
		AssignmentMismatches:           review.AssignmentTelemetryMismatches,
		TripDescriptorOmissions:        review.TripDescriptorOmissions,
		UsefulnessStatus:               review.UsefulnessStatus,
		UsefulnessReasons:              review.UsefulnessReasons,
		Truncated:                      snapshot.Truncated,
		VehicleLimit:                   snapshot.VehicleLimit,
		LatestTelemetryRowsRead:        snapshot.LatestTelemetryRowsRead,
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

func boundedCountMap(input map[string]int) map[string]int {
	out := make(map[string]int, len(input))
	for key, value := range input {
		if key == "" {
			continue
		}
		out[key] = boundedNonNegative(value)
	}
	return out
}

func boundedHealthText(value string) string {
	if len(value) > 80 {
		return value[:80]
	}
	return value
}

func boundedStringList(values []string) []string {
	const maxItems = 8
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		out = append(out, boundedHealthText(value))
		if len(out) >= maxItems {
			break
		}
	}
	return out
}
