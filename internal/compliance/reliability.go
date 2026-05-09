package compliance

import (
	"sort"
	"strings"
	"time"
)

const (
	ReliabilityStatusOK          = "ok"
	ReliabilityStatusNeedsReview = "needs_review"
	ReliabilityStatusMissing     = "missing"
	ReliabilityStatusUnknown     = "unknown"
	ReliabilityStatusUnhealthy   = "unhealthy"
)

var reliabilityFeedOrder = []string{"schedule", "vehicle_positions", "trip_updates", "alerts"}

type ReliabilityClaimFlags struct {
	ExternalEvidenceCreated    bool `json:"external_evidence_created"`
	FinalRootEvidenceCreated   bool `json:"final_root_evidence_created"`
	ConsumerStatusesChanged    bool `json:"consumer_statuses_changed"`
	ComplianceClaimed          bool `json:"compliance_claimed"`
	ProductionReadinessClaimed bool `json:"production_readiness_claimed"`
	SLAClaimed                 bool `json:"sla_claimed"`
	UptimeGuaranteeClaimed     bool `json:"uptime_guarantee_claimed"`
	HostedSaaSClaimed          bool `json:"hosted_saas_claimed"`
	AgencyAdoptionClaimed      bool `json:"agency_adoption_claimed"`
	ConsumerAcceptanceClaimed  bool `json:"consumer_acceptance_claimed"`
	VendorCompatibilityClaimed bool `json:"vendor_compatibility_claimed"`
	ProductionGradeETAClaimed  bool `json:"production_grade_eta_claimed"`
}

type ReliabilitySummary struct {
	GeneratedAt           time.Time                 `json:"generated_at"`
	AgencyID              string                    `json:"agency_id"`
	OverallStatus         string                    `json:"overall_status"`
	Feeds                 []ReliabilityFeedRow      `json:"feeds"`
	Incidents             ReliabilityIncidentRollup `json:"incidents"`
	BackupRestore         ReliabilitySection        `json:"backup_restore"`
	Alerting              ReliabilitySection        `json:"alerting"`
	AvailabilitySampling  ReliabilitySection        `json:"availability_sampling"`
	LongRunningOperations ReliabilitySection        `json:"long_running_operations"`
	ClaimFlags            ReliabilityClaimFlags     `json:"claim_flags"`
}

type ReliabilityFeedRow struct {
	FeedType               string     `json:"feed_type"`
	Status                 string     `json:"status"`
	Source                 string     `json:"source"`
	SnapshotAt             *time.Time `json:"snapshot_at"`
	EndpointAvailable      *bool      `json:"endpoint_available"`
	FreshnessSeconds       *float64   `json:"freshness_seconds"`
	GenerationLatencyMS    *float64   `json:"generation_latency_ms"`
	InvalidResponsePercent *float64   `json:"invalid_response_percent"`
	MatchedVehiclePercent  *float64   `json:"matched_vehicle_percent"`
	CoveragePercent        *float64   `json:"coverage_percent"`
	DiagnosticThreshold    string     `json:"diagnostic_threshold"`
	NextAction             string     `json:"next_action"`
}

type ReliabilityFeedHealthRecord struct {
	FeedType               string
	SnapshotAt             time.Time
	EndpointAvailable      *bool
	FreshnessSeconds       *float64
	GenerationLatencyMS    *float64
	InvalidResponsePercent *float64
	MatchedVehiclePercent  *float64
	CoveragePercent        *float64
}

type ReliabilityIncidentRollup struct {
	Status               string                    `json:"status"`
	Source               string                    `json:"source"`
	Total                int                       `json:"total"`
	CountsByStatus       map[string]int            `json:"counts_by_status"`
	CountsBySeverity     map[string]int            `json:"counts_by_severity"`
	CountsByType         map[string]int            `json:"counts_by_type"`
	OldestOpenAgeSeconds *int64                    `json:"oldest_open_age_seconds,omitempty"`
	Recent               []ReliabilityIncidentItem `json:"recent"`
	RecentLimit          int                       `json:"recent_limit"`
	NextAction           string                    `json:"next_action"`
}

type ReliabilityIncidentItem struct {
	ID        int64      `json:"id"`
	Type      string     `json:"type"`
	Severity  string     `json:"severity"`
	Status    string     `json:"status"`
	OpenedAt  time.Time  `json:"opened_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Title     string     `json:"title"`
	Category  string     `json:"category"`
}

type ReliabilitySection struct {
	Status     string            `json:"status"`
	Source     string            `json:"source"`
	Summary    string            `json:"summary"`
	Signals    map[string]string `json:"signals,omitempty"`
	NextAction string            `json:"next_action"`
}

type ReliabilityInput struct {
	GeneratedAt       time.Time
	AgencyID          string
	FeedHealthRecords []ReliabilityFeedHealthRecord
	Incidents         ReliabilityIncidentRollup
	BackupRestore     *ReliabilitySection
	Alerting          *ReliabilitySection
	Availability      *ReliabilitySection
	LongRunning       *ReliabilitySection
}

func BuildReliabilitySummary(input ReliabilityInput) ReliabilitySummary {
	generatedAt := input.GeneratedAt.UTC()
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC().Truncate(time.Second)
	}
	feeds := buildReliabilityFeedRows(input.FeedHealthRecords)
	incidents := input.Incidents
	if incidents.Source == "" {
		incidents = ReliabilityIncidentRollup{
			Status:           ReliabilityStatusUnknown,
			Source:           "incident table",
			CountsByStatus:   map[string]int{},
			CountsBySeverity: map[string]int{},
			CountsByType:     map[string]int{},
			RecentLimit:      0,
			NextAction:       "Incident rollup is unavailable; review database access before using this diagnostic.",
		}
	}
	incidents.Status = normalizeReliabilityStatus(incidents.Status)
	sections := []ReliabilitySection{
		reliabilitySectionOrMissing(input.BackupRestore, "safe private backup/restore summary", "Run a private backup/restore diagnostic summary before reviewing this section."),
		reliabilitySectionOrMissing(input.Alerting, "safe private alerting summary", "Run a private alerting diagnostic summary before reviewing this section."),
		reliabilitySectionOrMissing(input.Availability, "safe private availability sampling summary", "Run local diagnostic sampling before reviewing this section."),
		reliabilitySectionOrMissing(input.LongRunning, "safe private long-running operations summary", "Run private long-running operations diagnostics before reviewing this section."),
	}
	overall := ReliabilityStatusUnknown
	for _, row := range feeds {
		overall = worseReliabilityStatus(overall, row.Status)
	}
	overall = worseReliabilityStatus(overall, incidents.Status)
	for _, section := range sections {
		overall = worseReliabilityStatus(overall, section.Status)
	}
	return ReliabilitySummary{
		GeneratedAt:           generatedAt,
		AgencyID:              strings.TrimSpace(input.AgencyID),
		OverallStatus:         overall,
		Feeds:                 feeds,
		Incidents:             incidents,
		BackupRestore:         sections[0],
		Alerting:              sections[1],
		AvailabilitySampling:  sections[2],
		LongRunningOperations: sections[3],
		ClaimFlags:            ReliabilityClaimFlags{},
	}
}

func buildReliabilityFeedRows(records []ReliabilityFeedHealthRecord) []ReliabilityFeedRow {
	latest := map[string]ReliabilityFeedHealthRecord{}
	for _, record := range records {
		feedType := strings.TrimSpace(record.FeedType)
		if !requiredReliabilityFeed(feedType) {
			continue
		}
		current, ok := latest[feedType]
		if !ok || record.SnapshotAt.After(current.SnapshotAt) {
			latest[feedType] = record
		}
	}
	rows := make([]ReliabilityFeedRow, 0, len(reliabilityFeedOrder))
	for _, feedType := range reliabilityFeedOrder {
		record, ok := latest[feedType]
		if !ok {
			status := ReliabilityStatusUnknown
			next := "No feed health snapshot exists for this feed. Missing data is not treated as ok."
			if feedType != "schedule" {
				status = ReliabilityStatusMissing
				next = "Generate this feed and persist a private feed health snapshot before treating it as observed."
			}
			rows = append(rows, ReliabilityFeedRow{
				FeedType:            feedType,
				Status:              status,
				Source:              "feed_health_snapshot",
				DiagnosticThreshold: "missing data remains unknown or missing",
				NextAction:          next,
			})
			continue
		}
		t := record.SnapshotAt.UTC()
		row := ReliabilityFeedRow{
			FeedType:               feedType,
			Status:                 reliabilityFeedStatus(record),
			Source:                 "feed_health_snapshot",
			SnapshotAt:             &t,
			EndpointAvailable:      record.EndpointAvailable,
			FreshnessSeconds:       record.FreshnessSeconds,
			GenerationLatencyMS:    record.GenerationLatencyMS,
			InvalidResponsePercent: record.InvalidResponsePercent,
			MatchedVehiclePercent:  record.MatchedVehiclePercent,
			CoveragePercent:        record.CoveragePercent,
			DiagnosticThreshold:    "private diagnostic threshold; not an SLA or uptime guarantee",
		}
		row.NextAction = reliabilityFeedNextAction(row)
		rows = append(rows, row)
	}
	return rows
}

func reliabilityFeedStatus(record ReliabilityFeedHealthRecord) string {
	if record.EndpointAvailable != nil && !*record.EndpointAvailable {
		return ReliabilityStatusUnhealthy
	}
	if record.InvalidResponsePercent != nil && *record.InvalidResponsePercent > 0 {
		return ReliabilityStatusUnhealthy
	}
	if record.EndpointAvailable == nil && record.FreshnessSeconds == nil && record.GenerationLatencyMS == nil && record.MatchedVehiclePercent == nil && record.CoveragePercent == nil {
		return ReliabilityStatusUnknown
	}
	if record.FreshnessSeconds != nil && *record.FreshnessSeconds > 90 {
		return ReliabilityStatusNeedsReview
	}
	if record.GenerationLatencyMS != nil && *record.GenerationLatencyMS > 5000 {
		return ReliabilityStatusNeedsReview
	}
	if record.MatchedVehiclePercent != nil && *record.MatchedVehiclePercent < 50 {
		return ReliabilityStatusNeedsReview
	}
	if record.CoveragePercent != nil && *record.CoveragePercent < 50 {
		return ReliabilityStatusNeedsReview
	}
	if record.EndpointAvailable != nil && *record.EndpointAvailable {
		return ReliabilityStatusOK
	}
	return ReliabilityStatusNeedsReview
}

func reliabilityFeedNextAction(row ReliabilityFeedRow) string {
	switch row.Status {
	case ReliabilityStatusOK:
		return "Continue local diagnostic sampling and review trends; this is not an SLA or production-readiness claim."
	case ReliabilityStatusUnhealthy:
		return "Investigate the latest observed feed failure before relying on this feed."
	case ReliabilityStatusNeedsReview:
		return "Review freshness, latency, validity, and coverage signals against private diagnostic thresholds."
	case ReliabilityStatusMissing:
		return "Persist a private feed health snapshot before treating this feed as observed."
	default:
		return "Instrumentation is absent or no database row exists; leave status unknown until observed."
	}
}

func reliabilitySectionOrMissing(section *ReliabilitySection, source, next string) ReliabilitySection {
	if section == nil {
		return ReliabilitySection{
			Status:     ReliabilityStatusMissing,
			Source:     source,
			Summary:    "No safe private diagnostic summary was found.",
			NextAction: next,
		}
	}
	out := *section
	out.Status = normalizeReliabilityStatus(out.Status)
	if out.Status == "" {
		out.Status = ReliabilityStatusUnknown
	}
	if out.Source == "" {
		out.Source = source
	}
	if out.Summary == "" {
		out.Summary = "Safe private diagnostic summary is present but did not provide a bounded summary."
	}
	if out.NextAction == "" {
		out.NextAction = "Review this private diagnostic source before making any operational decision."
	}
	return out
}

func NormalizeReliabilityIncidentRollup(now time.Time, total int, byStatus, bySeverity, byType map[string]int, oldestOpen *time.Time, recent []ReliabilityIncidentItem, limit int) ReliabilityIncidentRollup {
	if limit <= 0 {
		limit = 10
	}
	if len(recent) > limit {
		recent = recent[:limit]
	}
	var age *int64
	if oldestOpen != nil && !oldestOpen.IsZero() {
		seconds := int64(now.UTC().Sub(oldestOpen.UTC()).Seconds())
		if seconds < 0 {
			seconds = 0
		}
		age = &seconds
	}
	status := ReliabilityStatusOK
	if total > 0 {
		status = ReliabilityStatusNeedsReview
	}
	return ReliabilityIncidentRollup{
		Status:               status,
		Source:               "incident table",
		Total:                total,
		CountsByStatus:       sortedCountMap(byStatus),
		CountsBySeverity:     sortedCountMap(bySeverity),
		CountsByType:         sortedCountMap(byType),
		OldestOpenAgeSeconds: age,
		Recent:               sanitizeReliabilityIncidentItems(recent),
		RecentLimit:          limit,
		NextAction:           reliabilityIncidentNextAction(total),
	}
}

func reliabilityIncidentNextAction(total int) string {
	if total == 0 {
		return "No incidents were found in the incident table for this private diagnostic sample."
	}
	return "Review open or acknowledged incidents in the operations workflow; this rollup intentionally omits raw details."
}

func sanitizeReliabilityIncidentItems(items []ReliabilityIncidentItem) []ReliabilityIncidentItem {
	out := make([]ReliabilityIncidentItem, 0, len(items))
	for _, item := range items {
		item.Type = safeReliabilityToken(item.Type, "unknown")
		item.Severity = safeReliabilityToken(item.Severity, "unknown")
		item.Status = safeReliabilityToken(item.Status, "unknown")
		item.Title = safeReliabilityTitle(item.Title, item.Type, item.ID)
		item.Category = safeReliabilityToken(item.Category, item.Type)
		out = append(out, item)
	}
	return out
}

func safeReliabilityTitle(value, incidentType string, id int64) string {
	value = strings.TrimSpace(value)
	if value == "" || unsafeReliabilityText(value) {
		if incidentType == "" {
			incidentType = "unknown"
		}
		return "Incident " + safeReliabilityToken(incidentType, "unknown")
	}
	return truncateReliabilityText(value, 80)
}

func safeReliabilityToken(value, fallback string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" || unsafeReliabilityText(value) {
		return fallback
	}
	var b strings.Builder
	for _, r := range value {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return fallback
	}
	return truncateReliabilityText(b.String(), 64)
}

func unsafeReliabilityText(value string) bool {
	lower := strings.ToLower(value)
	for _, marker := range []string{"token", "secret", "password", "authorization", "cookie", "postgres://", "database_url", "webhook", "payload", "details_json", "/users/", "/tmp/", "/var/", "/etc/", "http://", "https://"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func truncateReliabilityText(value string, max int) string {
	if len(value) <= max {
		return value
	}
	if max <= 12 {
		return value[:max]
	}
	return strings.TrimSpace(value[:max-12]) + " [truncated]"
}

func requiredReliabilityFeed(feedType string) bool {
	for _, expected := range reliabilityFeedOrder {
		if feedType == expected {
			return true
		}
	}
	return false
}

func normalizeReliabilityStatus(status string) string {
	switch strings.TrimSpace(status) {
	case ReliabilityStatusOK, ReliabilityStatusNeedsReview, ReliabilityStatusMissing, ReliabilityStatusUnknown, ReliabilityStatusUnhealthy:
		return strings.TrimSpace(status)
	default:
		return ReliabilityStatusUnknown
	}
}

func worseReliabilityStatus(a, b string) string {
	rank := map[string]int{
		ReliabilityStatusOK:          0,
		ReliabilityStatusUnknown:     1,
		ReliabilityStatusMissing:     2,
		ReliabilityStatusNeedsReview: 3,
		ReliabilityStatusUnhealthy:   4,
	}
	a = normalizeReliabilityStatus(a)
	b = normalizeReliabilityStatus(b)
	if rank[b] > rank[a] {
		return b
	}
	return a
}

func sortedCountMap(in map[string]int) map[string]int {
	tmp := map[string]int{}
	for key, value := range in {
		if value <= 0 {
			continue
		}
		tmp[safeReliabilityToken(key, "unknown")] += value
	}
	out := map[string]int{}
	keys := make([]string, 0, len(tmp))
	for key := range tmp {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		out[key] = tmp[key]
	}
	return out
}
