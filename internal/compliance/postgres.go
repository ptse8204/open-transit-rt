package compliance

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
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
		Scan(&record.ConsumerName, &record.Status, &submitted, &record.UpdatedAt, &record.Notes, &packetBytes)
	if err != nil {
		return ConsumerRecord{}, fmt.Errorf("upsert consumer ingestion: %w", err)
	}
	if submitted.Valid {
		t := submitted.Time
		record.SubmittedAt = &t
	}
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
		var packetBytes []byte
		if err := rows.Scan(&record.ConsumerName, &record.Status, &submitted, &record.UpdatedAt, &record.Notes, &packetBytes); err != nil {
			return nil, fmt.Errorf("scan consumer ingestion: %w", err)
		}
		if submitted.Valid {
			t := submitted.Time
			record.SubmittedAt = &t
		}
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
	for _, feedType := range RequiredFeedTypes {
		feed, ok := feedMap[feedType]
		if !ok || feed.CanonicalPublicURL == "" {
			allRequired = false
			httpsURLs = false
			canonicalValidationComplete = false
			continue
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
		Discoverable:                discoverable,
		HTTPSURLs:                   httpsURLs,
		LicenseComplete:             licenseComplete,
		ContactComplete:             contactComplete,
		AllRequiredFeedsListed:      allRequired,
		CanonicalValidationComplete: canonicalValidationComplete,
	}
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
