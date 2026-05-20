package compliance

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"open-transit-rt/internal/prediction"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) ListAuditLog(ctx context.Context, agencyID string, limit int) ([]AuditLogRecord, error) {
	if agencyID == "" {
		return nil, fmt.Errorf("agency_id is required")
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
		  id,
		  created_at,
		  action,
		  entity_type,
		  COALESCE(entity_id, ''),
		  actor_id <> '',
		  NULLIF(BTRIM(COALESCE(reason, '')), '') IS NOT NULL,
		  old_value_json IS NOT NULL,
		  new_value_json IS NOT NULL
		FROM audit_log
		WHERE agency_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2
	`, agencyID, limit)
	if err != nil {
		return nil, fmt.Errorf("query audit log: %w", err)
	}
	defer rows.Close()
	var records []AuditLogRecord
	for rows.Next() {
		var record AuditLogRecord
		if err := rows.Scan(&record.ID, &record.CreatedAt, &record.Action, &record.EntityType, &record.EntityID, &record.ActorRecorded, &record.ReasonRecorded, &record.OldValueRecorded, &record.NewValueRecorded); err != nil {
			return nil, fmt.Errorf("scan audit log: %w", err)
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate audit log: %w", err)
	}
	return records, nil
}

func (r *PostgresRepository) BootstrapPublication(ctx context.Context, input BootstrapInput) error {
	if input.AgencyID == "" {
		return fmt.Errorf("agency_id is required")
	}
	if input.PublicationEnvironment == "" {
		input.PublicationEnvironment = EnvironmentDev
	}
	if input.ActorID == "" {
		input.ActorID = "system"
	}
	now := time.Now().UTC()
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin publication bootstrap: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		INSERT INTO feed_config (
			agency_id, public_base_url, feed_base_url, technical_contact_email,
			license_name, license_url, publication_environment, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (agency_id) DO UPDATE
		SET public_base_url = EXCLUDED.public_base_url,
		    feed_base_url = EXCLUDED.feed_base_url,
		    technical_contact_email = EXCLUDED.technical_contact_email,
		    license_name = EXCLUDED.license_name,
		    license_url = EXCLUDED.license_url,
		    publication_environment = EXCLUDED.publication_environment,
		    updated_at = EXCLUDED.updated_at
	`, input.AgencyID, input.PublicBaseURL, input.FeedBaseURL, input.TechnicalContactEmail,
		input.LicenseName, nullString(input.LicenseURL), input.PublicationEnvironment, now); err != nil {
		return fmt.Errorf("upsert feed config: %w", err)
	}

	activeFeedVersionID := ""
	_ = tx.QueryRow(ctx, `
		SELECT id
		FROM feed_version
		WHERE agency_id = $1 AND is_active
		ORDER BY activated_at DESC NULLS LAST, created_at DESC
		LIMIT 1
	`, input.AgencyID).Scan(&activeFeedVersionID)

	urls := canonicalURLs(input.PublicBaseURL, input.FeedBaseURL)
	for _, feedType := range RequiredFeedTypes {
		if _, err := tx.Exec(ctx, `
			INSERT INTO published_feed (
				agency_id, feed_type, canonical_public_url, license_name, license_url,
				contact_email, revision_timestamp, activation_status, active_feed_version_id,
				metadata_json, updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 'active', NULLIF($8, ''), $9::jsonb, $7)
			ON CONFLICT (agency_id, feed_type) DO UPDATE
			SET canonical_public_url = EXCLUDED.canonical_public_url,
			    license_name = EXCLUDED.license_name,
			    license_url = EXCLUDED.license_url,
			    contact_email = EXCLUDED.contact_email,
			    activation_status = EXCLUDED.activation_status,
			    active_feed_version_id = EXCLUDED.active_feed_version_id,
			    metadata_json = EXCLUDED.metadata_json,
			    updated_at = EXCLUDED.updated_at,
			    revision_timestamp = EXCLUDED.revision_timestamp
		`, input.AgencyID, feedType, urls[feedType], input.LicenseName, nullString(input.LicenseURL),
			input.TechnicalContactEmail, now, activeFeedVersionID, metadataForFeed(feedType)); err != nil {
			return fmt.Errorf("upsert published feed %s: %w", feedType, err)
		}
	}
	for _, consumer := range DefaultConsumers {
		if _, err := tx.Exec(ctx, `
			INSERT INTO consumer_ingestion (agency_id, consumer_name, status, packet_json, updated_at)
			VALUES ($1, $2, 'not_started', '{}'::jsonb, $3)
			ON CONFLICT (agency_id, consumer_name) DO NOTHING
		`, input.AgencyID, consumer, now); err != nil {
			return fmt.Errorf("seed consumer ingestion %s: %w", consumer, err)
		}
	}
	for _, gap := range DefaultMarketplaceGaps {
		if _, err := tx.Exec(ctx, `
			INSERT INTO marketplace_gap (agency_id, gap_key, status, updated_at)
			VALUES ($1, $2, 'not_started', $3)
			ON CONFLICT (agency_id, gap_key) DO NOTHING
		`, input.AgencyID, gap, now); err != nil {
			return fmt.Errorf("seed marketplace gap %s: %w", gap, err)
		}
	}
	auditPayload, _ := json.Marshal(input)
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (agency_id, actor_id, action, entity_type, entity_id, new_value_json, reason)
		VALUES ($1, $2, 'publication.bootstrap', 'feed_config', $1, $3, 'phase_8_publication_metadata')
	`, input.AgencyID, input.ActorID, string(auditPayload)); err != nil {
		return fmt.Errorf("insert publication bootstrap audit log: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit publication bootstrap: %w", err)
	}
	return nil
}

func (r *PostgresRepository) FeedDiscovery(ctx context.Context, agencyID string, generatedAt time.Time) (FeedDiscovery, error) {
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}
	cfg, err := r.feedConfig(ctx, agencyID)
	if err != nil {
		return FeedDiscovery{}, err
	}
	var agencyName string
	if err := r.pool.QueryRow(ctx, `SELECT name FROM agency WHERE id = $1`, agencyID).Scan(&agencyName); err != nil {
		return FeedDiscovery{}, fmt.Errorf("query agency name: %w", err)
	}
	feeds, err := r.feedMetadata(ctx, agencyID, cfg)
	if err != nil {
		return FeedDiscovery{}, err
	}
	readiness := evaluateReadiness(cfg, feeds)
	return FeedDiscovery{
		AgencyID:               agencyID,
		AgencyName:             agencyName,
		GeneratedAt:            generatedAt.UTC(),
		PublicationEnvironment: cfg.PublicationEnvironment,
		PublicBaseURL:          cfg.PublicBaseURL,
		TechnicalContactEmail:  cfg.TechnicalContactEmail,
		License: License{
			Name: cfg.LicenseName,
			URL:  cfg.LicenseURL,
		},
		Feeds:     feeds,
		Readiness: readiness,
	}, nil
}

func (r *PostgresRepository) PublicationConfig(ctx context.Context, agencyID string) (PublicationConfig, error) {
	cfg, err := r.feedConfig(ctx, agencyID)
	if err != nil {
		return PublicationConfig{}, err
	}
	return PublicationConfig{
		AgencyID:               agencyID,
		PublicBaseURL:          cfg.PublicBaseURL,
		FeedBaseURL:            cfg.FeedBaseURL,
		TechnicalContactEmail:  cfg.TechnicalContactEmail,
		LicenseName:            cfg.LicenseName,
		LicenseURL:             cfg.LicenseURL,
		PublicationEnvironment: cfg.PublicationEnvironment,
	}, nil
}

func (r *PostgresRepository) UpsertConsumer(ctx context.Context, input ConsumerInput) (ConsumerRecord, error) {
	if input.Status == "" {
		input.Status = "not_started"
	}
	packet, err := json.Marshal(input.Packet)
	if err != nil {
		return ConsumerRecord{}, fmt.Errorf("marshal consumer packet: %w", err)
	}
	now := time.Now().UTC()
	var submittedAt any
	if input.Status == "submitted" || input.Status == "resubmitted" {
		submittedAt = now
	}
	var record ConsumerRecord
	var submitted sql.NullTime
	var notes sql.NullString
	var packetBytes []byte
	err = r.pool.QueryRow(ctx, `
		INSERT INTO consumer_ingestion (agency_id, consumer_name, status, submitted_at, updated_at, notes, packet_json)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb)
		ON CONFLICT (agency_id, consumer_name) DO UPDATE
		SET status = EXCLUDED.status,
		    submitted_at = COALESCE(EXCLUDED.submitted_at, consumer_ingestion.submitted_at),
		    updated_at = EXCLUDED.updated_at,
		    notes = EXCLUDED.notes,
		    packet_json = EXCLUDED.packet_json
		RETURNING consumer_name, status, submitted_at, updated_at, notes, packet_json
	`, input.AgencyID, input.ConsumerName, input.Status, submittedAt, now, nullString(input.Notes), string(packet)).
		Scan(&record.ConsumerName, &record.Status, &submitted, &record.UpdatedAt, &notes, &packetBytes)
	if err != nil {
		return ConsumerRecord{}, fmt.Errorf("upsert consumer ingestion: %w", err)
	}
	if submitted.Valid {
		t := submitted.Time
		record.SubmittedAt = &t
	}
	record.Notes = notes.String
	record.Packet = map[string]any{}
	_ = json.Unmarshal(packetBytes, &record.Packet)
	return record, nil
}

func (r *PostgresRepository) ListConsumers(ctx context.Context, agencyID string) ([]ConsumerRecord, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT consumer_name, status, submitted_at, updated_at, notes, packet_json
		FROM consumer_ingestion
		WHERE agency_id = $1
		ORDER BY consumer_name
	`, agencyID)
	if err != nil {
		return nil, fmt.Errorf("query consumer ingestion: %w", err)
	}
	defer rows.Close()
	var records []ConsumerRecord
	for rows.Next() {
		var record ConsumerRecord
		var submitted sql.NullTime
		var notes sql.NullString
		var packetBytes []byte
		if err := rows.Scan(&record.ConsumerName, &record.Status, &submitted, &record.UpdatedAt, &notes, &packetBytes); err != nil {
			return nil, fmt.Errorf("scan consumer ingestion: %w", err)
		}
		if submitted.Valid {
			t := submitted.Time
			record.SubmittedAt = &t
		}
		record.Notes = notes.String
		record.Packet = map[string]any{}
		_ = json.Unmarshal(packetBytes, &record.Packet)
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate consumer ingestion: %w", err)
	}
	return records, nil
}

func (r *PostgresRepository) BuildAndStoreScorecard(ctx context.Context, agencyID string, at time.Time) (Scorecard, error) {
	discovery, err := r.FeedDiscovery(ctx, agencyID, at)
	if err != nil {
		return Scorecard{}, err
	}
	consumers, err := r.ListConsumers(ctx, agencyID)
	if err != nil {
		return Scorecard{}, err
	}
	validationStatus := validationScore(discovery.PublicationEnvironment, discovery.Feeds)
	discoverabilityStatus := boolScore(discovery.Readiness.Discoverable)
	consumerStatus := consumerScore(consumers)
	feedStatus := map[string]string{}
	for _, feed := range discovery.Feeds {
		feedStatus[feed.FeedType] = feedScore(feed)
	}
	scorecard := Scorecard{
		AgencyID:                agencyID,
		SnapshotAt:              discovery.GeneratedAt,
		PublicationEnvironment:  discovery.PublicationEnvironment,
		ScheduleStatus:          defaultStatus(feedStatus["schedule"]),
		VehiclePositionsStatus:  defaultStatus(feedStatus["vehicle_positions"]),
		TripUpdatesStatus:       defaultStatus(feedStatus["trip_updates"]),
		AlertsStatus:            defaultStatus(feedStatus["alerts"]),
		ValidationStatus:        validationStatus,
		DiscoverabilityStatus:   discoverabilityStatus,
		ConsumerIngestionStatus: consumerStatus,
		Details: map[string]any{
			"feeds":     discovery.Feeds,
			"readiness": discovery.Readiness,
			"consumers": consumers,
		},
	}
	scorecard.OverallStatus = worstStatus(scorecard.ScheduleStatus, scorecard.VehiclePositionsStatus, scorecard.TripUpdatesStatus, scorecard.AlertsStatus, scorecard.ValidationStatus, scorecard.DiscoverabilityStatus, scorecard.ConsumerIngestionStatus)
	details, err := json.Marshal(scorecard.Details)
	if err != nil {
		return Scorecard{}, fmt.Errorf("marshal scorecard details: %w", err)
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO compliance_scorecard_snapshot (
			agency_id, snapshot_at, publication_environment, overall_status, schedule_status,
			vehicle_positions_status, trip_updates_status, alerts_status, validation_status,
			discoverability_status, consumer_ingestion_status, details_json
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12::jsonb)
	`, scorecard.AgencyID, scorecard.SnapshotAt, scorecard.PublicationEnvironment, scorecard.OverallStatus,
		scorecard.ScheduleStatus, scorecard.VehiclePositionsStatus, scorecard.TripUpdatesStatus,
		scorecard.AlertsStatus, scorecard.ValidationStatus, scorecard.DiscoverabilityStatus,
		scorecard.ConsumerIngestionStatus, string(details))
	if err != nil {
		return Scorecard{}, fmt.Errorf("insert scorecard snapshot: %w", err)
	}
	return scorecard, nil
}

func (r *PostgresRepository) LatestScorecard(ctx context.Context, agencyID string) (Scorecard, error) {
	var scorecard Scorecard
	var detailsBytes []byte
	err := r.pool.QueryRow(ctx, `
		SELECT agency_id, snapshot_at, publication_environment, overall_status, schedule_status,
		       vehicle_positions_status, trip_updates_status, alerts_status, validation_status,
		       discoverability_status, consumer_ingestion_status, details_json
		FROM compliance_scorecard_snapshot
		WHERE agency_id = $1
		ORDER BY snapshot_at DESC, id DESC
		LIMIT 1
	`, agencyID).Scan(
		&scorecard.AgencyID,
		&scorecard.SnapshotAt,
		&scorecard.PublicationEnvironment,
		&scorecard.OverallStatus,
		&scorecard.ScheduleStatus,
		&scorecard.VehiclePositionsStatus,
		&scorecard.TripUpdatesStatus,
		&scorecard.AlertsStatus,
		&scorecard.ValidationStatus,
		&scorecard.DiscoverabilityStatus,
		&scorecard.ConsumerIngestionStatus,
		&detailsBytes,
	)
	if err != nil {
		return Scorecard{}, fmt.Errorf("query latest scorecard: %w", err)
	}
	scorecard.Details = map[string]any{}
	_ = json.Unmarshal(detailsBytes, &scorecard.Details)
	return scorecard, nil
}

func (r *PostgresRepository) LatestTripUpdatesDiagnostics(ctx context.Context, agencyID string) (TripUpdatesDiagnosticsSummary, error) {
	var snapshotAt time.Time
	var detailsBytes []byte
	err := r.pool.QueryRow(ctx, `
		SELECT snapshot_at, details_json
		FROM feed_health_snapshot
		WHERE agency_id = $1 AND feed_type = 'trip_updates'
		ORDER BY snapshot_at DESC, id DESC
		LIMIT 1
	`, agencyID).Scan(&snapshotAt, &detailsBytes)
	if err != nil {
		if err == pgx.ErrNoRows {
			return TripUpdatesDiagnosticsSummary{Recorded: false}, nil
		}
		return TripUpdatesDiagnosticsSummary{}, fmt.Errorf("query latest trip updates diagnostics: %w", err)
	}
	var details struct {
		AdapterName                   string             `json:"adapter_name"`
		DiagnosticsStatus             string             `json:"diagnostics_status"`
		DiagnosticsReason             string             `json:"diagnostics_reason"`
		ActiveFeedVersionID           string             `json:"active_feed_version_id"`
		PredictionMetrics             prediction.Metrics `json:"prediction_metrics"`
		VehiclePositionsURL           string             `json:"vehicle_positions_url"`
		DiagnosticsPersistenceOutcome string             `json:"diagnostics_persistence_outcome"`
		AdapterDetails                map[string]any     `json:"adapter_details"`
	}
	if err := json.Unmarshal(detailsBytes, &details); err != nil {
		return TripUpdatesDiagnosticsSummary{}, fmt.Errorf("decode latest trip updates diagnostics: %w", err)
	}
	return TripUpdatesDiagnosticsSummary{
		Recorded:                      true,
		SnapshotAt:                    snapshotAt.UTC(),
		AdapterName:                   details.AdapterName,
		DiagnosticsStatus:             details.DiagnosticsStatus,
		DiagnosticsReason:             details.DiagnosticsReason,
		ActiveFeedVersionID:           details.ActiveFeedVersionID,
		VehiclePositionsURL:           details.VehiclePositionsURL,
		DiagnosticsPersistenceOutcome: details.DiagnosticsPersistenceOutcome,
		AdapterDetails:                details.AdapterDetails,
		Metrics:                       details.PredictionMetrics,
	}, nil
}

func (r *PostgresRepository) LatestReliabilityFeedHealth(ctx context.Context, agencyID string, limit int) ([]ReliabilityFeedHealthRecord, error) {
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	rows, err := r.pool.Query(ctx, `
		SELECT feed_type, snapshot_at, endpoint_available, freshness_seconds,
		       generation_latency_ms, invalid_response_percent,
		       matched_vehicle_percent, coverage_percent
		FROM feed_health_snapshot
		WHERE agency_id = $1
		ORDER BY snapshot_at DESC, id DESC
		LIMIT $2
	`, agencyID, limit)
	if err != nil {
		return nil, fmt.Errorf("query reliability feed health snapshots: %w", err)
	}
	defer rows.Close()
	var records []ReliabilityFeedHealthRecord
	for rows.Next() {
		var record ReliabilityFeedHealthRecord
		var endpoint sql.NullBool
		var freshness, latency, invalid, matched, coverage sql.NullFloat64
		if err := rows.Scan(&record.FeedType, &record.SnapshotAt, &endpoint, &freshness, &latency, &invalid, &matched, &coverage); err != nil {
			return nil, fmt.Errorf("scan reliability feed health snapshot: %w", err)
		}
		if endpoint.Valid {
			v := endpoint.Bool
			record.EndpointAvailable = &v
		}
		if freshness.Valid {
			v := freshness.Float64
			record.FreshnessSeconds = &v
		}
		if latency.Valid {
			v := latency.Float64
			record.GenerationLatencyMS = &v
		}
		if invalid.Valid {
			v := invalid.Float64
			record.InvalidResponsePercent = &v
		}
		if matched.Valid {
			v := matched.Float64
			record.MatchedVehiclePercent = &v
		}
		if coverage.Valid {
			v := coverage.Float64
			record.CoveragePercent = &v
		}
		record.SnapshotAt = record.SnapshotAt.UTC()
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reliability feed health snapshots: %w", err)
	}
	return records, nil
}

func (r *PostgresRepository) ReliabilityIncidentRollup(ctx context.Context, agencyID string, now time.Time, recentLimit int) (ReliabilityIncidentRollup, error) {
	if recentLimit <= 0 || recentLimit > 25 {
		recentLimit = 10
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	type incidentAggregate struct {
		kind  string
		key   string
		count int
	}
	rows, err := r.pool.Query(ctx, `
		SELECT 'status' AS kind, status AS key, COUNT(*)::int AS count
		FROM incident
		WHERE agency_id = $1
		GROUP BY status
		UNION ALL
		SELECT 'severity' AS kind, severity AS key, COUNT(*)::int AS count
		FROM incident
		WHERE agency_id = $1
		GROUP BY severity
		UNION ALL
		SELECT 'type' AS kind, incident_type AS key, COUNT(*)::int AS count
		FROM incident
		WHERE agency_id = $1
		GROUP BY incident_type
	`, agencyID)
	if err != nil {
		return ReliabilityIncidentRollup{}, fmt.Errorf("query reliability incident counts: %w", err)
	}
	countsByStatus := map[string]int{}
	countsBySeverity := map[string]int{}
	countsByType := map[string]int{}
	total := 0
	for rows.Next() {
		var row incidentAggregate
		if err := rows.Scan(&row.kind, &row.key, &row.count); err != nil {
			rows.Close()
			return ReliabilityIncidentRollup{}, fmt.Errorf("scan reliability incident count: %w", err)
		}
		switch row.kind {
		case "status":
			countsByStatus[row.key] = row.count
			total += row.count
		case "severity":
			countsBySeverity[row.key] = row.count
		case "type":
			countsByType[row.key] = row.count
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return ReliabilityIncidentRollup{}, fmt.Errorf("iterate reliability incident counts: %w", err)
	}
	rows.Close()

	var oldest sql.NullTime
	if err := r.pool.QueryRow(ctx, `
		SELECT MIN(created_at)
		FROM incident
		WHERE agency_id = $1 AND status <> 'resolved'
	`, agencyID).Scan(&oldest); err != nil {
		return ReliabilityIncidentRollup{}, fmt.Errorf("query oldest open incident: %w", err)
	}
	var oldestPtr *time.Time
	if oldest.Valid {
		t := oldest.Time.UTC()
		oldestPtr = &t
	}

	recentRows, err := r.pool.Query(ctx, `
		SELECT id, incident_type, severity, status, created_at,
		       COALESCE(last_seen_at, resolved_at, created_at) AS updated_at
		FROM incident
		WHERE agency_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2
	`, agencyID, recentLimit)
	if err != nil {
		return ReliabilityIncidentRollup{}, fmt.Errorf("query recent reliability incidents: %w", err)
	}
	defer recentRows.Close()
	recent := make([]ReliabilityIncidentItem, 0, recentLimit)
	for recentRows.Next() {
		var item ReliabilityIncidentItem
		var updated time.Time
		if err := recentRows.Scan(&item.ID, &item.Type, &item.Severity, &item.Status, &item.OpenedAt, &updated); err != nil {
			return ReliabilityIncidentRollup{}, fmt.Errorf("scan recent reliability incident: %w", err)
		}
		item.OpenedAt = item.OpenedAt.UTC()
		t := updated.UTC()
		item.UpdatedAt = &t
		item.Title = "Incident " + item.Type
		item.Category = item.Type
		recent = append(recent, item)
	}
	if err := recentRows.Err(); err != nil {
		return ReliabilityIncidentRollup{}, fmt.Errorf("iterate recent reliability incidents: %w", err)
	}
	return NormalizeReliabilityIncidentRollup(now, total, countsByStatus, countsBySeverity, countsByType, oldestPtr, recent, recentLimit), nil
}

func (r *PostgresRepository) StoreValidationResult(ctx context.Context, result ValidationResult) error {
	report, err := json.Marshal(result.Report)
	if err != nil {
		return fmt.Errorf("marshal validation report: %w", err)
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO validation_report (
			agency_id, feed_version_id, feed_type, validator_name, validator_version,
			status, error_count, warning_count, info_count, report_json
		)
		VALUES ($1, NULLIF($2, ''), $3, $4, $5, $6, $7, $8, $9, $10::jsonb)
	`, result.AgencyID, result.FeedVersionID, result.FeedType, result.ValidatorName, result.ValidatorVersion,
		result.Status, result.ErrorCount, result.WarningCount, result.InfoCount, string(report))
	if err != nil {
		return fmt.Errorf("insert validation report: %w", err)
	}
	return nil
}

func (r *PostgresRepository) LatestValidationReport(ctx context.Context, agencyID string, feedType string, validatorName string) (*ValidationReportRecord, error) {
	var record ValidationReportRecord
	var reportBytes []byte
	err := r.pool.QueryRow(ctx, `
		SELECT id, agency_id, COALESCE(feed_version_id, ''), feed_type, validator_name,
		       validator_version, status, error_count, warning_count, info_count, report_json, created_at
		FROM validation_report
		WHERE agency_id = $1
		  AND feed_type = $2
		  AND validator_name = $3
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, agencyID, feedType, validatorName).Scan(
		&record.ID,
		&record.Result.AgencyID,
		&record.Result.FeedVersionID,
		&record.Result.FeedType,
		&record.Result.ValidatorName,
		&record.Result.ValidatorVersion,
		&record.Result.Status,
		&record.Result.ErrorCount,
		&record.Result.WarningCount,
		&record.Result.InfoCount,
		&reportBytes,
		&record.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("query latest validation report: %w", err)
	}
	if err := json.Unmarshal(reportBytes, &record.Result.Report); err != nil {
		record.Result.Report = map[string]any{}
	}
	record.CreatedAt = record.CreatedAt.UTC()
	return &record, nil
}

func (r *PostgresRepository) RecentGTFSImports(ctx context.Context, agencyID string, limit int) ([]GTFSImportRecord, error) {
	if limit <= 0 || limit > 25 {
		limit = 10
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, agency_id, COALESCE(feed_version_id, ''), source_filename,
		       source_sha256, source_byte_size, status, error_count, warning_count,
		       info_count, COALESCE(actor_id, ''), started_at, completed_at
		FROM gtfs_import
		WHERE agency_id = $1
		ORDER BY started_at DESC, id DESC
		LIMIT $2
	`, agencyID, limit)
	if err != nil {
		return nil, fmt.Errorf("query recent GTFS imports: %w", err)
	}
	defer rows.Close()
	records := make([]GTFSImportRecord, 0, limit)
	for rows.Next() {
		var record GTFSImportRecord
		var completed sql.NullTime
		if err := rows.Scan(
			&record.ID,
			&record.AgencyID,
			&record.FeedVersionID,
			&record.SourceFilename,
			&record.SourceSHA256,
			&record.SourceByteSize,
			&record.Status,
			&record.ErrorCount,
			&record.WarningCount,
			&record.InfoCount,
			&record.ActorID,
			&record.StartedAt,
			&completed,
		); err != nil {
			return nil, fmt.Errorf("scan recent GTFS import: %w", err)
		}
		record.StartedAt = record.StartedAt.UTC()
		if completed.Valid {
			t := completed.Time.UTC()
			record.CompletedAt = &t
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recent GTFS imports: %w", err)
	}
	return records, nil
}

func (r *PostgresRepository) GTFSSchedulePreview(ctx context.Context, agencyID string, feedVersionID string, limit int) (GTFSSchedulePreview, error) {
	if limit <= 0 || limit > 25 {
		limit = 10
	}
	preview := GTFSSchedulePreview{
		AgencyID:      agencyID,
		FeedVersionID: feedVersionID,
		RowLimit:      limit,
	}
	var timezone sql.NullString
	if err := r.pool.QueryRow(ctx, `
		SELECT id, name, timezone
		FROM agency
		WHERE id = $1
	`, agencyID).Scan(&preview.Agency.AgencyID, &preview.Agency.Name, &timezone); err != nil {
		return GTFSSchedulePreview{}, fmt.Errorf("query GTFS preview agency: %w", err)
	}
	preview.Agency.Timezone = timezone.String

	countQueries := []struct {
		target *int
		query  string
	}{
		{&preview.Counts.Routes, `SELECT COUNT(*)::int FROM gtfs_route WHERE agency_id = $1 AND feed_version_id = $2`},
		{&preview.Counts.Stops, `SELECT COUNT(*)::int FROM gtfs_stop WHERE agency_id = $1 AND feed_version_id = $2`},
		{&preview.Counts.Trips, `SELECT COUNT(*)::int FROM gtfs_trip WHERE agency_id = $1 AND feed_version_id = $2`},
		{&preview.Counts.StopTimes, `SELECT COUNT(*)::int FROM gtfs_stop_time WHERE agency_id = $1 AND feed_version_id = $2`},
		{&preview.Counts.Calendar, `SELECT COUNT(*)::int FROM gtfs_calendar WHERE agency_id = $1 AND feed_version_id = $2`},
		{&preview.Counts.CalendarDates, `SELECT COUNT(*)::int FROM gtfs_calendar_date WHERE agency_id = $1 AND feed_version_id = $2`},
		{&preview.Counts.ShapePoints, `SELECT COUNT(*)::int FROM gtfs_shape_point WHERE agency_id = $1 AND feed_version_id = $2`},
		{&preview.Counts.Frequencies, `SELECT COUNT(*)::int FROM gtfs_frequency WHERE agency_id = $1 AND feed_version_id = $2`},
	}
	for _, item := range countQueries {
		if err := r.pool.QueryRow(ctx, item.query, agencyID, feedVersionID).Scan(item.target); err != nil {
			return GTFSSchedulePreview{}, fmt.Errorf("query GTFS preview counts: %w", err)
		}
	}

	routes, err := r.previewRoutes(ctx, agencyID, feedVersionID, limit)
	if err != nil {
		return GTFSSchedulePreview{}, err
	}
	preview.Routes = routes
	stops, err := r.previewStops(ctx, agencyID, feedVersionID, limit)
	if err != nil {
		return GTFSSchedulePreview{}, err
	}
	preview.Stops = stops
	trips, err := r.previewTrips(ctx, agencyID, feedVersionID, limit)
	if err != nil {
		return GTFSSchedulePreview{}, err
	}
	preview.Trips = trips
	calendar, err := r.previewCalendar(ctx, agencyID, feedVersionID, limit)
	if err != nil {
		return GTFSSchedulePreview{}, err
	}
	preview.Calendar = calendar
	frequencies, err := r.previewFrequencies(ctx, agencyID, feedVersionID, limit)
	if err != nil {
		return GTFSSchedulePreview{}, err
	}
	preview.Frequencies = frequencies
	return preview, nil
}

func (r *PostgresRepository) previewRoutes(ctx context.Context, agencyID string, feedVersionID string, limit int) ([]GTFSScheduleRoutePreview, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, COALESCE(short_name, ''), COALESCE(long_name, ''), route_type
		FROM gtfs_route
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY id
		LIMIT $3
	`, agencyID, feedVersionID, limit)
	if err != nil {
		return nil, fmt.Errorf("query GTFS preview routes: %w", err)
	}
	defer rows.Close()
	var routes []GTFSScheduleRoutePreview
	for rows.Next() {
		var route GTFSScheduleRoutePreview
		var routeType sql.NullInt64
		if err := rows.Scan(&route.ID, &route.ShortName, &route.LongName, &routeType); err != nil {
			return nil, fmt.Errorf("scan GTFS preview route: %w", err)
		}
		if routeType.Valid {
			route.RouteType = fmt.Sprintf("%d", routeType.Int64)
		}
		routes = append(routes, route)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate GTFS preview routes: %w", err)
	}
	return routes, nil
}

func (r *PostgresRepository) previewStops(ctx context.Context, agencyID string, feedVersionID string, limit int) ([]GTFSScheduleStopPreview, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, lat, lon
		FROM gtfs_stop
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY id
		LIMIT $3
	`, agencyID, feedVersionID, limit)
	if err != nil {
		return nil, fmt.Errorf("query GTFS preview stops: %w", err)
	}
	defer rows.Close()
	var stops []GTFSScheduleStopPreview
	for rows.Next() {
		var stop GTFSScheduleStopPreview
		if err := rows.Scan(&stop.ID, &stop.Name, &stop.Lat, &stop.Lon); err != nil {
			return nil, fmt.Errorf("scan GTFS preview stop: %w", err)
		}
		stops = append(stops, stop)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate GTFS preview stops: %w", err)
	}
	return stops, nil
}

func (r *PostgresRepository) previewTrips(ctx context.Context, agencyID string, feedVersionID string, limit int) ([]GTFSScheduleTripPreview, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, route_id, service_id, COALESCE(block_id, ''), COALESCE(shape_id, ''), direction_id
		FROM gtfs_trip
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY route_id, id
		LIMIT $3
	`, agencyID, feedVersionID, limit)
	if err != nil {
		return nil, fmt.Errorf("query GTFS preview trips: %w", err)
	}
	defer rows.Close()
	var trips []GTFSScheduleTripPreview
	for rows.Next() {
		var trip GTFSScheduleTripPreview
		var direction sql.NullInt64
		if err := rows.Scan(&trip.ID, &trip.RouteID, &trip.ServiceID, &trip.BlockID, &trip.ShapeID, &direction); err != nil {
			return nil, fmt.Errorf("scan GTFS preview trip: %w", err)
		}
		if direction.Valid {
			trip.DirectionID = fmt.Sprintf("%d", direction.Int64)
		}
		trips = append(trips, trip)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate GTFS preview trips: %w", err)
	}
	return trips, nil
}

func (r *PostgresRepository) previewCalendar(ctx context.Context, agencyID string, feedVersionID string, limit int) ([]GTFSScheduleCalendarPreview, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT service_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday, start_date, end_date
		FROM gtfs_calendar
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY service_id
		LIMIT $3
	`, agencyID, feedVersionID, limit)
	if err != nil {
		return nil, fmt.Errorf("query GTFS preview calendar: %w", err)
	}
	defer rows.Close()
	var calendar []GTFSScheduleCalendarPreview
	for rows.Next() {
		var row GTFSScheduleCalendarPreview
		var days [7]bool
		if err := rows.Scan(&row.ServiceID, &days[0], &days[1], &days[2], &days[3], &days[4], &days[5], &days[6], &row.StartDate, &row.EndDate); err != nil {
			return nil, fmt.Errorf("scan GTFS preview calendar: %w", err)
		}
		row.Days = calendarDaysText(days)
		calendar = append(calendar, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate GTFS preview calendar: %w", err)
	}
	return calendar, nil
}

func calendarDaysText(days [7]bool) string {
	labels := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	var out []string
	for i, enabled := range days {
		if enabled {
			out = append(out, labels[i])
		}
	}
	if len(out) == 0 {
		return "no weekly days"
	}
	return strings.Join(out, ", ")
}

func (r *PostgresRepository) previewFrequencies(ctx context.Context, agencyID string, feedVersionID string, limit int) ([]GTFSScheduleFrequencyPreview, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT trip_id, start_time, end_time, headway_secs, exact_times
		FROM gtfs_frequency
		WHERE agency_id = $1 AND feed_version_id = $2
		ORDER BY trip_id, start_time
		LIMIT $3
	`, agencyID, feedVersionID, limit)
	if err != nil {
		return nil, fmt.Errorf("query GTFS preview frequencies: %w", err)
	}
	defer rows.Close()
	var frequencies []GTFSScheduleFrequencyPreview
	for rows.Next() {
		var frequency GTFSScheduleFrequencyPreview
		if err := rows.Scan(&frequency.TripID, &frequency.StartTime, &frequency.EndTime, &frequency.HeadwaySecs, &frequency.ExactTimes); err != nil {
			return nil, fmt.Errorf("scan GTFS preview frequency: %w", err)
		}
		frequencies = append(frequencies, frequency)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate GTFS preview frequencies: %w", err)
	}
	return frequencies, nil
}

func (r *PostgresRepository) RecentGTFSDrafts(ctx context.Context, agencyID string, limit int) ([]GTFSDraftRecord, error) {
	if limit <= 0 || limit > 25 {
		limit = 10
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, agency_id, name, status, COALESCE(base_feed_version_id, ''),
		       COALESCE(last_published_feed_version_id, ''), COALESCE(last_publish_attempt_id, 0),
		       created_at, updated_at
		FROM gtfs_draft
		WHERE agency_id = $1
		ORDER BY updated_at DESC, created_at DESC, id
		LIMIT $2
	`, agencyID, limit)
	if err != nil {
		return nil, fmt.Errorf("query recent GTFS drafts: %w", err)
	}
	defer rows.Close()
	records := make([]GTFSDraftRecord, 0, limit)
	for rows.Next() {
		var record GTFSDraftRecord
		if err := rows.Scan(&record.ID, &record.AgencyID, &record.Name, &record.Status, &record.BaseFeedVersionID, &record.LastPublishedFeedVersionID, &record.LastPublishAttemptID, &record.CreatedAt, &record.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan recent GTFS draft: %w", err)
		}
		record.CreatedAt = record.CreatedAt.UTC()
		record.UpdatedAt = record.UpdatedAt.UTC()
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recent GTFS drafts: %w", err)
	}
	return records, nil
}

func (r *PostgresRepository) RecentGTFSDraftPublishes(ctx context.Context, agencyID string, limit int) ([]GTFSDraftPublishRecord, error) {
	if limit <= 0 || limit > 25 {
		limit = 10
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, draft_id, COALESCE(feed_version_id, ''), status, error_count,
		       warning_count, info_count, COALESCE(actor_id, ''), started_at, completed_at
		FROM gtfs_draft_publish
		WHERE agency_id = $1
		ORDER BY started_at DESC, id DESC
		LIMIT $2
	`, agencyID, limit)
	if err != nil {
		return nil, fmt.Errorf("query recent GTFS draft publishes: %w", err)
	}
	defer rows.Close()
	records := make([]GTFSDraftPublishRecord, 0, limit)
	for rows.Next() {
		var record GTFSDraftPublishRecord
		var completed sql.NullTime
		if err := rows.Scan(&record.ID, &record.DraftID, &record.FeedVersionID, &record.Status, &record.ErrorCount, &record.WarningCount, &record.InfoCount, &record.ActorID, &record.StartedAt, &completed); err != nil {
			return nil, fmt.Errorf("scan recent GTFS draft publish: %w", err)
		}
		record.StartedAt = record.StartedAt.UTC()
		if completed.Valid {
			t := completed.Time.UTC()
			record.CompletedAt = &t
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recent GTFS draft publishes: %w", err)
	}
	return records, nil
}

func (r *PostgresRepository) RecentFeedVersions(ctx context.Context, agencyID string, limit int) ([]FeedVersionRecord, error) {
	if limit <= 0 || limit > 25 {
		limit = 10
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, agency_id, source_type, lifecycle_state, is_active, validation_status,
		       published_at, activated_at, retired_at, created_at
		FROM feed_version
		WHERE agency_id = $1
		ORDER BY is_active DESC, activated_at DESC NULLS LAST, created_at DESC, id
		LIMIT $2
	`, agencyID, limit)
	if err != nil {
		return nil, fmt.Errorf("query recent feed versions: %w", err)
	}
	defer rows.Close()
	records := make([]FeedVersionRecord, 0, limit)
	for rows.Next() {
		var record FeedVersionRecord
		var published, activated, retired sql.NullTime
		if err := rows.Scan(&record.ID, &record.AgencyID, &record.SourceType, &record.LifecycleState, &record.IsActive, &record.ValidationStatus, &published, &activated, &retired, &record.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan recent feed version: %w", err)
		}
		record.CreatedAt = record.CreatedAt.UTC()
		if published.Valid {
			t := published.Time.UTC()
			record.PublishedAt = &t
		}
		if activated.Valid {
			t := activated.Time.UTC()
			record.ActivatedAt = &t
		}
		if retired.Valid {
			t := retired.Time.UTC()
			record.RetiredAt = &t
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recent feed versions: %w", err)
	}
	return records, nil
}

type feedConfig struct {
	PublicBaseURL          string
	FeedBaseURL            string
	TechnicalContactEmail  string
	LicenseName            string
	LicenseURL             string
	PublicationEnvironment string
}

func (r *PostgresRepository) feedConfig(ctx context.Context, agencyID string) (feedConfig, error) {
	var cfg feedConfig
	var licenseURL sql.NullString
	err := r.pool.QueryRow(ctx, `
		SELECT public_base_url, feed_base_url, technical_contact_email, license_name,
		       license_url, publication_environment
		FROM feed_config
		WHERE agency_id = $1
	`, agencyID).Scan(&cfg.PublicBaseURL, &cfg.FeedBaseURL, &cfg.TechnicalContactEmail, &cfg.LicenseName, &licenseURL, &cfg.PublicationEnvironment)
	if err != nil {
		return feedConfig{}, fmt.Errorf("query feed config: %w", err)
	}
	cfg.LicenseURL = licenseURL.String
	return cfg, nil
}

func (r *PostgresRepository) feedMetadata(ctx context.Context, agencyID string, cfg feedConfig) ([]FeedMetadata, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT feed_type, canonical_public_url, activation_status, active_feed_version_id,
		       revision_timestamp, license_name, license_url, contact_email
		FROM published_feed
		WHERE agency_id = $1
		ORDER BY feed_type
	`, agencyID)
	if err != nil {
		return nil, fmt.Errorf("query published feed metadata: %w", err)
	}
	defer rows.Close()
	var feeds []FeedMetadata
	for rows.Next() {
		var feed FeedMetadata
		var feedVersionID, licenseURL sql.NullString
		var revision sql.NullTime
		if err := rows.Scan(&feed.FeedType, &feed.CanonicalPublicURL, &feed.ActivationStatus, &feedVersionID, &revision, &feed.LicenseName, &licenseURL, &feed.ContactEmail); err != nil {
			return nil, fmt.Errorf("scan published feed metadata: %w", err)
		}
		feed.ActiveFeedVersionID = feedVersionID.String
		if revision.Valid {
			t := revision.Time.UTC()
			feed.RevisionTimestamp = &t
		}
		feed.LicenseURL = licenseURL.String
		if feed.LicenseName == "" {
			feed.LicenseName = cfg.LicenseName
		}
		if feed.LicenseURL == "" {
			feed.LicenseURL = cfg.LicenseURL
		}
		if feed.ContactEmail == "" {
			feed.ContactEmail = cfg.TechnicalContactEmail
		}
		validationStatus, validationAt := r.latestValidation(ctx, agencyID, feed.FeedType)
		feed.LastValidationStatus = validationStatus
		feed.LastValidationAt = validationAt
		healthStatus, healthAt := r.latestHealth(ctx, agencyID, feed.FeedType)
		feed.LastHealthStatus = healthStatus
		feed.LastHealthAt = healthAt
		feeds = append(feeds, feed)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate published feed metadata: %w", err)
	}
	sort.SliceStable(feeds, func(i int, j int) bool {
		return feedOrder(feeds[i].FeedType) < feedOrder(feeds[j].FeedType)
	})
	return feeds, nil
}

func (r *PostgresRepository) latestValidation(ctx context.Context, agencyID string, feedType string) (string, *time.Time) {
	var status string
	var createdAt time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT status, created_at
		FROM validation_report
		WHERE agency_id = $1 AND feed_type = $2
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, agencyID, feedType).Scan(&status, &createdAt)
	if err != nil {
		return "not_run", nil
	}
	t := createdAt.UTC()
	return status, &t
}

func (r *PostgresRepository) latestHealth(ctx context.Context, agencyID string, feedType string) (string, *time.Time) {
	var endpoint sql.NullBool
	var snapshotAt time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT endpoint_available, snapshot_at
		FROM feed_health_snapshot
		WHERE agency_id = $1 AND feed_type = $2
		ORDER BY snapshot_at DESC, id DESC
		LIMIT 1
	`, agencyID, feedType).Scan(&endpoint, &snapshotAt)
	if err != nil {
		return "not_run", nil
	}
	status := "degraded"
	if endpoint.Valid && endpoint.Bool {
		status = "ok"
	}
	if endpoint.Valid && !endpoint.Bool {
		status = "unhealthy"
	}
	t := snapshotAt.UTC()
	return status, &t
}

func evaluateReadiness(cfg feedConfig, feeds []FeedMetadata) Readiness {
	feedMap := map[string]FeedMetadata{}
	for _, feed := range feeds {
		feedMap[feed.FeedType] = feed
	}
	allRequired := true
	httpsURLs := true
	licenseComplete := cfg.LicenseName != "" && cfg.LicenseURL != ""
	contactComplete := cfg.TechnicalContactEmail != ""
	canonicalValidationComplete := true
	activeScheduleListed := false
	realtimeFeedsListed := true
	for _, feedType := range RequiredFeedTypes {
		feed, ok := feedMap[feedType]
		if !ok || feed.CanonicalPublicURL == "" {
			allRequired = false
			httpsURLs = false
			canonicalValidationComplete = false
			if feedType != "schedule" {
				realtimeFeedsListed = false
			}
			continue
		}
		if feedType == "schedule" && strings.EqualFold(feed.ActivationStatus, "active") && strings.TrimSpace(feed.ActiveFeedVersionID) != "" {
			activeScheduleListed = true
		}
		parsed, err := url.Parse(feed.CanonicalPublicURL)
		if err != nil || parsed.Scheme != "https" {
			httpsURLs = false
		}
		if feed.LicenseName == "" || feed.LicenseURL == "" {
			licenseComplete = false
		}
		if feed.ContactEmail == "" {
			contactComplete = false
		}
		if feed.LastValidationStatus != "passed" && feed.LastValidationStatus != "warning" {
			canonicalValidationComplete = false
		}
	}
	discoverable := allRequired && licenseComplete && contactComplete
	if cfg.PublicationEnvironment == EnvironmentProduction {
		discoverable = discoverable && httpsURLs && canonicalValidationComplete
	}
	return Readiness{
		Discoverable:                     discoverable,
		HTTPSURLs:                        httpsURLs,
		LicenseComplete:                  licenseComplete,
		ContactComplete:                  contactComplete,
		AllRequiredFeedsListed:           allRequired,
		CanonicalValidationComplete:      canonicalValidationComplete,
		StablePublicBaseURL:              stablePublicBaseURL(cfg.PublicBaseURL),
		PublicationEnvironmentConfigured: strings.TrimSpace(cfg.PublicationEnvironment) != "",
		ActiveScheduleListed:             activeScheduleListed,
		RealtimeFeedsListed:              realtimeFeedsListed,
	}
}

func stablePublicBaseURL(raw string) bool {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	if host == "localhost" || strings.HasSuffix(host, ".local") {
		return false
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() {
			return false
		}
		return true
	}
	return strings.Contains(host, ".")
}

func canonicalURLs(publicBaseURL string, feedBaseURL string) map[string]string {
	publicBase := strings.TrimRight(publicBaseURL, "/")
	feedBase := strings.TrimRight(feedBaseURL, "/")
	if feedBase == "" {
		feedBase = publicBase + "/public"
	}
	return map[string]string{
		"schedule":          publicBase + "/public/gtfs/schedule.zip",
		"vehicle_positions": feedBase + "/gtfsrt/vehicle_positions.pb",
		"trip_updates":      feedBase + "/gtfsrt/trip_updates.pb",
		"alerts":            feedBase + "/gtfsrt/alerts.pb",
	}
}

func metadataForFeed(feedType string) string {
	payload, _ := json.Marshal(map[string]any{
		"revision_timestamp_semantics": "realtime feeds change revision_timestamp only when publication/bootstrap metadata changes; schedule changes on active feed publication",
		"feed_type":                    feedType,
	})
	return string(payload)
}

func validationScore(environment string, feeds []FeedMetadata) string {
	result := StatusGreen
	for _, feed := range feeds {
		switch feed.LastValidationStatus {
		case "passed":
		case "warning":
			result = worse(result, StatusYellow)
		case "not_run", "":
			if environment == EnvironmentProduction {
				result = worse(result, StatusRed)
			} else {
				result = worse(result, StatusYellow)
			}
		default:
			result = worse(result, StatusRed)
		}
	}
	if len(feeds) == 0 {
		if environment == EnvironmentProduction {
			return StatusRed
		}
		return StatusYellow
	}
	return result
}

func feedScore(feed FeedMetadata) string {
	if feed.ActivationStatus != "active" || feed.CanonicalPublicURL == "" {
		return StatusRed
	}
	if feed.LastHealthStatus == "unhealthy" {
		return StatusRed
	}
	if feed.LastHealthStatus == "not_run" || feed.LastValidationStatus == "not_run" {
		return StatusYellow
	}
	return StatusGreen
}

func consumerScore(consumers []ConsumerRecord) string {
	if len(consumers) == 0 {
		return StatusRed
	}
	hasAccepted := false
	hasStarted := false
	for _, consumer := range consumers {
		if consumer.Status == "accepted" {
			hasAccepted = true
		}
		if consumer.Status != "not_started" {
			hasStarted = true
		}
	}
	if hasAccepted {
		return StatusGreen
	}
	if hasStarted {
		return StatusYellow
	}
	return StatusRed
}

func boolScore(ok bool) string {
	if ok {
		return StatusGreen
	}
	return StatusRed
}

func defaultStatus(status string) string {
	if status == "" {
		return StatusRed
	}
	return status
}

func worstStatus(statuses ...string) string {
	result := StatusGreen
	for _, status := range statuses {
		result = worse(result, status)
	}
	return result
}

func worse(left string, right string) string {
	if left == StatusRed || right == StatusRed {
		return StatusRed
	}
	if left == StatusYellow || right == StatusYellow {
		return StatusYellow
	}
	return StatusGreen
}

func feedOrder(feedType string) int {
	for i, candidate := range RequiredFeedTypes {
		if candidate == feedType {
			return i
		}
	}
	return len(RequiredFeedTypes)
}

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}
