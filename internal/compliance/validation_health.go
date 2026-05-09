package compliance

import (
	"sort"
	"strings"
	"time"
)

const (
	ValidationHealthStatusBlocked                = "blocked"
	ValidationHealthStatusFailed                 = "failed"
	ValidationHealthStatusMissingTooling         = "missing_tooling"
	ValidationHealthStatusMisconfiguredTooling   = "misconfigured_tooling"
	ValidationHealthStatusArtifactUnavailable    = "artifact_unavailable"
	ValidationHealthStatusStale                  = "stale"
	ValidationHealthStatusNeedsReview            = "needs_review"
	ValidationHealthStatusNotRun                 = "not_run"
	ValidationHealthStatusRunnable               = "runnable"
	ValidationHealthStatusRecorded               = "recorded"
	ValidationHealthStatusConfigured             = "configured"
	ValidationHealthStatusInstalled              = "installed"
	ValidationHealthStatusStub                   = "stub"
	ValidationHealthStatusConfiguredForTests     = "configured_for_tests"
	ValidationHealthStatusUnknown                = "unknown"
	ValidationHealthStatusSkipped                = "skipped"
	ValidationHealthArtifactAvailable            = "available"
	ValidationHealthArtifactUnavailable          = "artifact_unavailable"
	ValidationHealthArtifactUnknown              = "unknown"
	ValidationHealthStaleCurrent                 = "current"
	ValidationHealthStaleStale                   = "stale"
	ValidationHealthStaleUnknown                 = "unknown"
	ValidationHealthStaticValidatorID            = "static-mobilitydata"
	ValidationHealthRealtimeValidatorID          = "realtime-mobilitydata"
	ValidationHealthRealtimeValidatorName        = "mobilitydata-gtfs-realtime-validator"
	ValidationHealthClaimBoundaryPrivateOnly     = "Private diagnostics only; no evidence packet, consumer status change, compliance claim, acceptance claim, or production readiness claim."
	ValidationHealthStubClaimBoundaryPrivateOnly = "Private diagnostics only; stub tooling is for tests and is not production readiness."
)

var validationHealthFeedOrder = []string{"schedule", "vehicle_positions", "trip_updates", "alerts"}

type ValidationHealthSummary struct {
	GeneratedAt                time.Time             `json:"generated_at"`
	AgencyID                   string                `json:"agency_id"`
	OverallStatus              string                `json:"overall_status"`
	ToolingStatus              string                `json:"tooling_status"`
	Feeds                      []ValidationHealthRow `json:"feeds"`
	ExternalEvidenceCreated    bool                  `json:"external_evidence_created"`
	ConsumerStatusesChanged    bool                  `json:"consumer_statuses_changed"`
	ComplianceClaimed          bool                  `json:"compliance_claimed"`
	ProductionReadinessClaimed bool                  `json:"production_readiness_claimed"`
}

type ValidationHealthRow struct {
	FeedType                  string     `json:"feed_type"`
	ValidatorID               string     `json:"validator_id"`
	ValidatorName             string     `json:"validator_name"`
	ToolingStatus             string     `json:"tooling_status"`
	ArtifactStatus            string     `json:"artifact_status"`
	LatestResultStatus        string     `json:"latest_result_status"`
	LatestResultAt            *time.Time `json:"latest_result_at"`
	ActiveFeedVersionID       string     `json:"active_feed_version_id"`
	LatestResultFeedVersionID string     `json:"latest_result_feed_version_id"`
	StaleStatus               string     `json:"stale_status"`
	HealthStatus              string     `json:"health_status"`
	NextAction                string     `json:"next_action"`
	ClaimBoundary             string     `json:"claim_boundary"`
}

type ValidationHealthInput struct {
	GeneratedAt              time.Time
	AgencyID                 string
	Discovery                FeedDiscovery
	Registry                 ValidatorRegistry
	Records                  []ValidationReportRecord
	ToolingStatusByValidator map[string]string
	ArtifactStatusByFeed     map[string]string
}

func BuildValidationHealthSummary(input ValidationHealthInput) ValidationHealthSummary {
	generatedAt := input.GeneratedAt.UTC()
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC().Truncate(time.Second)
	}
	agencyID := firstNonEmpty(input.AgencyID, input.Discovery.AgencyID)
	latest := latestValidationHealthRecords(input.Records)
	tooling := input.ToolingStatusByValidator
	if tooling == nil {
		tooling = ValidationToolingStatusByValidator(input.Registry)
	}
	artifacts := input.ArtifactStatusByFeed
	if artifacts == nil {
		artifacts = artifactStatusByFeed(input.Discovery)
	}

	rows := make([]ValidationHealthRow, 0, len(validationHealthFeedOrder))
	for _, feedType := range validationHealthFeedOrder {
		validatorID := ValidatorIDForHealthFeed(feedType)
		validatorName := ValidatorNameForHealthID(validatorID)
		record := latest[feedType+"\x00"+validatorName]
		activeFeedVersionID := activeFeedVersionIDForHealth(input.Discovery, feedType)
		artifactStatus := firstNonEmpty(artifacts[feedType], ValidationHealthArtifactUnknown)
		toolingStatus := firstNonEmpty(tooling[validatorID], ValidationHealthStatusUnknown)
		row := buildValidationHealthRow(feedType, validatorID, validatorName, toolingStatus, artifactStatus, activeFeedVersionID, record)
		rows = append(rows, row)
	}
	overall := ValidationHealthStatusUnknown
	for _, row := range rows {
		overall = worseValidationHealthStatus(overall, row.HealthStatus)
	}
	return ValidationHealthSummary{
		GeneratedAt:                generatedAt,
		AgencyID:                   agencyID,
		OverallStatus:              overall,
		ToolingStatus:              overallToolingStatus(tooling),
		Feeds:                      rows,
		ExternalEvidenceCreated:    false,
		ConsumerStatusesChanged:    false,
		ComplianceClaimed:          false,
		ProductionReadinessClaimed: false,
	}
}

func ValidatorIDForHealthFeed(feedType string) string {
	if strings.TrimSpace(feedType) == "schedule" {
		return ValidationHealthStaticValidatorID
	}
	return ValidationHealthRealtimeValidatorID
}

func ValidatorNameForHealthID(validatorID string) string {
	switch validatorID {
	case ValidationHealthStaticValidatorID:
		return CanonicalStaticValidatorName
	case ValidationHealthRealtimeValidatorID:
		return ValidationHealthRealtimeValidatorName
	default:
		return ""
	}
}

func ValidationToolingStatusFromCheckValidators(mode string, exitCode int) string {
	mode = strings.TrimSpace(mode)
	if mode == "" {
		mode = "pinned"
	}
	switch exitCode {
	case 0:
		if mode == "stub" {
			return ValidationHealthStatusStub
		}
		return ValidationHealthStatusConfigured
	case 11:
		return ValidationHealthStatusMissingTooling
	case 12:
		return ValidationHealthStatusMisconfiguredTooling
	default:
		return ValidationHealthStatusBlocked
	}
}

func ValidationToolingStatusByValidator(registry ValidatorRegistry) map[string]string {
	statuses := map[string]string{
		ValidationHealthStaticValidatorID:   validatorToolingStatus(registry[ValidationHealthStaticValidatorID]),
		ValidationHealthRealtimeValidatorID: validatorToolingStatus(registry[ValidationHealthRealtimeValidatorID]),
	}
	return statuses
}

func validatorToolingStatus(spec ValidatorSpec) string {
	if strings.TrimSpace(spec.ID) == "" {
		return ValidationHealthStatusMissingTooling
	}
	if strings.TrimSpace(spec.Binary) == "" {
		return ValidationHealthStatusMissingTooling
	}
	if err := spec.validatePlaceholders(); err != nil {
		return ValidationHealthStatusMisconfiguredTooling
	}
	return ValidationHealthStatusConfigured
}

func buildValidationHealthRow(feedType, validatorID, validatorName, toolingStatus, artifactStatus, activeFeedVersionID string, record *ValidationReportRecord) ValidationHealthRow {
	row := ValidationHealthRow{
		FeedType:            feedType,
		ValidatorID:         validatorID,
		ValidatorName:       validatorName,
		ToolingStatus:       toolingStatus,
		ArtifactStatus:      artifactStatus,
		LatestResultStatus:  ValidationHealthStatusNotRun,
		ActiveFeedVersionID: activeFeedVersionID,
		StaleStatus:         ValidationHealthStaleUnknown,
		HealthStatus:        ValidationHealthStatusUnknown,
		ClaimBoundary:       ValidationHealthClaimBoundaryPrivateOnly,
	}
	if toolingStatus == ValidationHealthStatusStub || toolingStatus == ValidationHealthStatusConfiguredForTests {
		row.ClaimBoundary = ValidationHealthStubClaimBoundaryPrivateOnly
	}
	if record != nil {
		t := record.CreatedAt.UTC()
		row.LatestResultAt = &t
		row.LatestResultFeedVersionID = record.Result.FeedVersionID
		row.LatestResultStatus = validationHealthResultStatus(record.Result)
		row.StaleStatus = validationHealthStaleStatus(activeFeedVersionID, record.Result.FeedVersionID)
	}
	row.HealthStatus = validationHealthRowStatus(row)
	row.NextAction = validationHealthNextAction(row)
	return row
}

func validationHealthResultStatus(result ValidationResult) string {
	status := strings.ToLower(strings.TrimSpace(result.Status))
	switch status {
	case "failed", "error", "blocked":
		return ValidationHealthStatusFailed
	case "warning", "warnings":
		return ValidationHealthStatusNeedsReview
	case "not_run":
		return ValidationHealthStatusNotRun
	case "passed":
		if result.ErrorCount > 0 {
			return ValidationHealthStatusFailed
		}
		if result.WarningCount > 0 {
			return ValidationHealthStatusNeedsReview
		}
		return ValidationHealthStatusRecorded
	default:
		if result.ErrorCount > 0 {
			return ValidationHealthStatusFailed
		}
		if result.WarningCount > 0 {
			return ValidationHealthStatusNeedsReview
		}
		if status == "" {
			return ValidationHealthStatusNotRun
		}
		return ValidationHealthStatusUnknown
	}
}

func validationHealthStaleStatus(activeFeedVersionID, latestFeedVersionID string) string {
	if activeFeedVersionID == "" {
		return ValidationHealthStaleUnknown
	}
	if latestFeedVersionID == "" {
		return ValidationHealthStaleUnknown
	}
	if activeFeedVersionID != latestFeedVersionID {
		return ValidationHealthStaleStale
	}
	return ValidationHealthStaleCurrent
}

func validationHealthRowStatus(row ValidationHealthRow) string {
	switch row.ToolingStatus {
	case ValidationHealthStatusMissingTooling:
		return ValidationHealthStatusMissingTooling
	case ValidationHealthStatusMisconfiguredTooling:
		return ValidationHealthStatusMisconfiguredTooling
	case ValidationHealthStatusBlocked:
		return ValidationHealthStatusBlocked
	}
	if row.ArtifactStatus == ValidationHealthArtifactUnavailable {
		return ValidationHealthStatusArtifactUnavailable
	}
	if row.LatestResultStatus == ValidationHealthStatusFailed {
		return ValidationHealthStatusFailed
	}
	if row.StaleStatus == ValidationHealthStaleStale {
		return ValidationHealthStatusStale
	}
	if row.LatestResultStatus == ValidationHealthStatusNeedsReview {
		return ValidationHealthStatusNeedsReview
	}
	if row.LatestResultStatus == ValidationHealthStatusNotRun {
		return ValidationHealthStatusNotRun
	}
	if row.LatestResultStatus == ValidationHealthStatusRecorded {
		return ValidationHealthStatusRecorded
	}
	if row.ToolingStatus == ValidationHealthStatusStub || row.ToolingStatus == ValidationHealthStatusConfiguredForTests {
		return ValidationHealthStatusConfiguredForTests
	}
	if row.ToolingStatus == ValidationHealthStatusConfigured || row.ToolingStatus == ValidationHealthStatusInstalled {
		return ValidationHealthStatusRunnable
	}
	return ValidationHealthStatusUnknown
}

func validationHealthNextAction(row ValidationHealthRow) string {
	switch row.HealthStatus {
	case ValidationHealthStatusMissingTooling:
		return "Install the pinned validator tooling with make validators-install, then rerun the private health check."
	case ValidationHealthStatusMisconfiguredTooling:
		return "Fix the pinned validator environment or wrapper configuration, then rerun the private health check."
	case ValidationHealthStatusBlocked:
		return "Inspect validator tooling configuration on the server and rerun after the blocker is resolved."
	case ValidationHealthStatusArtifactUnavailable:
		if row.FeedType == "schedule" {
			return "Publish or activate a schedule ZIP for this agency before running static validator health."
		}
		return "Generate or expose the current server-owned protobuf artifact for this agency before running realtime validator health."
	case ValidationHealthStatusFailed:
		return "Review the bounded validation summary, fix source data or tooling, and rerun the allowlisted validator."
	case ValidationHealthStatusStale:
		return "Rerun the allowlisted validator against the active feed version."
	case ValidationHealthStatusNeedsReview:
		return "Review validator warnings as operator diagnostics, decide whether data changes are needed, and rerun after changes."
	case ValidationHealthStatusNotRun:
		return "Run the admin-only validator health action or existing allowlisted validation workflow."
	case ValidationHealthStatusConfiguredForTests:
		return "Use pinned validator tooling for operator diagnostics outside deterministic test stubs."
	case ValidationHealthStatusRecorded, ValidationHealthStatusRunnable:
		return "Keep this as private diagnostics and rerun after feed artifacts change."
	default:
		return "Check feed metadata, validator tooling, and latest validation records, then rerun the private health check."
	}
}

func activeFeedVersionIDForHealth(discovery FeedDiscovery, feedType string) string {
	for _, feed := range discovery.Feeds {
		if feed.FeedType == feedType {
			return feed.ActiveFeedVersionID
		}
	}
	if feedType != "schedule" {
		for _, feed := range discovery.Feeds {
			if feed.FeedType == "schedule" {
				return feed.ActiveFeedVersionID
			}
		}
	}
	return ""
}

func artifactStatusByFeed(discovery FeedDiscovery) map[string]string {
	statuses := map[string]string{}
	for _, feedType := range validationHealthFeedOrder {
		statuses[feedType] = ValidationHealthArtifactUnavailable
	}
	for _, feed := range discovery.Feeds {
		if !isValidationHealthFeed(feed.FeedType) {
			continue
		}
		if feed.FeedType == "schedule" {
			if strings.TrimSpace(feed.ActiveFeedVersionID) != "" && strings.TrimSpace(feed.CanonicalPublicURL) != "" {
				statuses[feed.FeedType] = ValidationHealthArtifactAvailable
			}
			continue
		}
		if strings.TrimSpace(feed.CanonicalPublicURL) != "" {
			statuses[feed.FeedType] = ValidationHealthArtifactAvailable
		}
	}
	return statuses
}

func latestValidationHealthRecords(records []ValidationReportRecord) map[string]*ValidationReportRecord {
	latest := map[string]*ValidationReportRecord{}
	for i := range records {
		record := records[i]
		if !isValidationHealthFeed(record.Result.FeedType) {
			continue
		}
		validatorName := ValidatorNameForHealthID(ValidatorIDForHealthFeed(record.Result.FeedType))
		if record.Result.ValidatorName != validatorName {
			continue
		}
		key := record.Result.FeedType + "\x00" + record.Result.ValidatorName
		current := latest[key]
		if current == nil || record.CreatedAt.After(current.CreatedAt) || record.CreatedAt.Equal(current.CreatedAt) && record.ID > current.ID {
			copyRecord := record
			copyRecord.Result.Report = nil
			latest[key] = &copyRecord
		}
	}
	return latest
}

func isValidationHealthFeed(feedType string) bool {
	for _, candidate := range validationHealthFeedOrder {
		if feedType == candidate {
			return true
		}
	}
	return false
}

func overallToolingStatus(statuses map[string]string) string {
	overall := ValidationHealthStatusUnknown
	for _, validatorID := range []string{ValidationHealthStaticValidatorID, ValidationHealthRealtimeValidatorID} {
		overall = worseValidationHealthStatus(overall, firstNonEmpty(statuses[validatorID], ValidationHealthStatusUnknown))
	}
	return overall
}

func worseValidationHealthStatus(a, b string) string {
	if validationHealthRank(b) > validationHealthRank(a) {
		return b
	}
	return a
}

func validationHealthRank(status string) int {
	switch status {
	case ValidationHealthStatusBlocked:
		return 100
	case ValidationHealthStatusFailed:
		return 90
	case ValidationHealthStatusMissingTooling:
		return 80
	case ValidationHealthStatusMisconfiguredTooling:
		return 70
	case ValidationHealthStatusArtifactUnavailable:
		return 60
	case ValidationHealthStatusStale:
		return 50
	case ValidationHealthStatusNeedsReview:
		return 40
	case ValidationHealthStatusNotRun:
		return 30
	case ValidationHealthStatusRunnable, ValidationHealthStatusConfigured, ValidationHealthStatusInstalled:
		return 20
	case ValidationHealthStatusRecorded, ValidationHealthStatusStub, ValidationHealthStatusConfiguredForTests, ValidationHealthStatusSkipped:
		return 10
	default:
		return 0
	}
}

func SortValidationHealthRows(rows []ValidationHealthRow) {
	sort.SliceStable(rows, func(i, j int) bool {
		return validationHealthFeedRank(rows[i].FeedType) < validationHealthFeedRank(rows[j].FeedType)
	})
}

func validationHealthFeedRank(feedType string) int {
	for i, candidate := range validationHealthFeedOrder {
		if feedType == candidate {
			return i
		}
	}
	return len(validationHealthFeedOrder)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
