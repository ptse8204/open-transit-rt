package compliance

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestValidationToolingStatusFromCheckValidators(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		exitCode int
		want     string
	}{
		{name: "pinned success", mode: "pinned", exitCode: 0, want: ValidationHealthStatusConfigured},
		{name: "stub success", mode: "stub", exitCode: 0, want: ValidationHealthStatusStub},
		{name: "missing tooling", mode: "pinned", exitCode: 11, want: ValidationHealthStatusMissingTooling},
		{name: "misconfigured tooling", mode: "pinned", exitCode: 12, want: ValidationHealthStatusMisconfiguredTooling},
		{name: "unknown fallback", mode: "pinned", exitCode: 99, want: ValidationHealthStatusBlocked},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidationToolingStatusFromCheckValidators(tt.mode, tt.exitCode); got != tt.want {
				t.Fatalf("status = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildValidationHealthSummaryShapeOrderAndFlags(t *testing.T) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	summary := BuildValidationHealthSummary(ValidationHealthInput{
		GeneratedAt: now,
		AgencyID:    "demo-agency",
		Discovery:   validationHealthDiscovery(now),
		ToolingStatusByValidator: map[string]string{
			ValidationHealthStaticValidatorID:   ValidationHealthStatusConfigured,
			ValidationHealthRealtimeValidatorID: ValidationHealthStatusConfigured,
		},
		Records: []ValidationReportRecord{
			recordForHealth(1, "vehicle_positions", "feed-v1", "passed", 0, 0, now),
			recordForHealth(2, "schedule", "feed-v1", "passed", 0, 0, now),
			recordForHealth(3, "alerts", "feed-v1", "warning", 0, 2, now),
			recordForHealth(4, "trip_updates", "feed-v1", "failed", 1, 0, now),
		},
	})
	if len(summary.Feeds) != 4 {
		t.Fatalf("feeds = %d, want 4", len(summary.Feeds))
	}
	wantOrder := []string{"schedule", "vehicle_positions", "trip_updates", "alerts"}
	for i, want := range wantOrder {
		if summary.Feeds[i].FeedType != want {
			t.Fatalf("feed order[%d] = %q, want %q", i, summary.Feeds[i].FeedType, want)
		}
	}
	if summary.ExternalEvidenceCreated || summary.ConsumerStatusesChanged || summary.ComplianceClaimed || summary.ProductionReadinessClaimed {
		t.Fatalf("false claim flags must be false: %+v", summary)
	}
	if summary.OverallStatus != ValidationHealthStatusFailed {
		t.Fatalf("overall = %q, want failed", summary.OverallStatus)
	}
	assertValidationHealthJSONAllowlist(t, summary)
	assertNoValidationHealthLeakage(t, summary)
}

func TestBuildValidationHealthSummaryMissingArtifactAndStale(t *testing.T) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	discovery := validationHealthDiscovery(now)
	discovery.Feeds = discovery.Feeds[:2]
	summary := BuildValidationHealthSummary(ValidationHealthInput{
		GeneratedAt: now,
		AgencyID:    "demo-agency",
		Discovery:   discovery,
		ToolingStatusByValidator: map[string]string{
			ValidationHealthStaticValidatorID:   ValidationHealthStatusConfigured,
			ValidationHealthRealtimeValidatorID: ValidationHealthStatusConfigured,
		},
		Records: []ValidationReportRecord{
			recordForHealth(1, "schedule", "old-feed", "passed", 0, 0, now),
			recordForHealth(2, "vehicle_positions", "feed-v1", "passed", 0, 0, now),
		},
	})
	if summary.Feeds[0].HealthStatus != ValidationHealthStatusStale {
		t.Fatalf("schedule health = %q, want stale", summary.Feeds[0].HealthStatus)
	}
	if summary.Feeds[2].ArtifactStatus != ValidationHealthArtifactUnavailable || summary.Feeds[2].HealthStatus != ValidationHealthStatusArtifactUnavailable {
		t.Fatalf("trip_updates row = %+v, want artifact_unavailable", summary.Feeds[2])
	}
	if summary.Feeds[3].ArtifactStatus != ValidationHealthArtifactUnavailable || summary.Feeds[3].HealthStatus != ValidationHealthStatusArtifactUnavailable {
		t.Fatalf("alerts row = %+v, want artifact_unavailable", summary.Feeds[3])
	}
	if len(summary.Feeds) != 4 {
		t.Fatalf("feeds = %d, want 4", len(summary.Feeds))
	}
}

func TestBuildValidationHealthManyReports_1k(t *testing.T) {
	assertBuildValidationHealthManyReports(t, 1000)
}

func TestBuildValidationHealthManyReports_10k(t *testing.T) {
	assertBuildValidationHealthManyReports(t, 10000)
}

func TestBuildValidationHealthHostileHistory(t *testing.T) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	records := make([]ValidationReportRecord, 0, 1000)
	huge := strings.Repeat("Authorization: Bearer TOKEN=SECRET PASSWORD=postgres://user:pass@localhost/db /tmp/private ", 500)
	for i := 0; i < 1000; i++ {
		feedType := validationHealthFeedOrder[i%len(validationHealthFeedOrder)]
		record := recordForHealth(int64(i+1), feedType, "feed-v1", "warning", 0, 1, now.Add(time.Duration(i)*time.Second))
		record.Result.Report = map[string]any{
			"raw_report": map[string]any{"huge": huge},
			"stdout":     huge,
			"stderr":     huge,
			"argv":       []any{"/tmp/private/bin", "--token", "SECRET"},
		}
		records = append(records, record)
	}
	summary := BuildValidationHealthSummary(ValidationHealthInput{
		GeneratedAt: now,
		AgencyID:    "demo-agency",
		Discovery:   validationHealthDiscovery(now),
		ToolingStatusByValidator: map[string]string{
			ValidationHealthStaticValidatorID:   ValidationHealthStatusConfigured,
			ValidationHealthRealtimeValidatorID: ValidationHealthStatusConfigured,
		},
		Records: records,
	})
	if len(summary.Feeds) != 4 {
		t.Fatalf("feeds = %d, want 4", len(summary.Feeds))
	}
	assertNoValidationHealthLeakage(t, summary)
	payload, err := json.Marshal(summary)
	if err != nil {
		t.Fatal(err)
	}
	if len(payload) > 16000 {
		t.Fatalf("summary JSON size = %d, want bounded", len(payload))
	}
}

func BenchmarkBuildValidationHealth_1k(b *testing.B) {
	benchmarkBuildValidationHealth(b, 1000)
}

func BenchmarkBuildValidationHealth_10k(b *testing.B) {
	benchmarkBuildValidationHealth(b, 10000)
}

func BenchmarkBuildValidationHealthHostileHistory(b *testing.B) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	records := validationHealthRecords(10000, now)
	huge := strings.Repeat("Authorization Bearer TOKEN=SECRET PASSWORD=postgres://user:pass@localhost/db /tmp/private ", 100)
	for i := range records {
		records[i].Result.Report = map[string]any{"raw_report": huge, "stdout": huge, "stderr": huge, "argv": []any{huge}}
	}
	input := benchmarkValidationHealthInput(now, records)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = BuildValidationHealthSummary(input)
	}
}

func assertBuildValidationHealthManyReports(t *testing.T, total int) {
	t.Helper()
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	summary := BuildValidationHealthSummary(benchmarkValidationHealthInput(now, validationHealthRecords(total, now)))
	if len(summary.Feeds) != 4 {
		t.Fatalf("feeds = %d, want 4", len(summary.Feeds))
	}
	if summary.Feeds[0].LatestResultFeedVersionID != "feed-v1" {
		t.Fatalf("schedule latest feed = %q, want feed-v1", summary.Feeds[0].LatestResultFeedVersionID)
	}
	assertNoValidationHealthLeakage(t, summary)
	payload, _ := json.Marshal(summary)
	if len(payload) > 16000 {
		t.Fatalf("summary JSON size = %d, want bounded", len(payload))
	}
}

func benchmarkBuildValidationHealth(b *testing.B, total int) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	input := benchmarkValidationHealthInput(now, validationHealthRecords(total, now))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = BuildValidationHealthSummary(input)
	}
}

func benchmarkValidationHealthInput(now time.Time, records []ValidationReportRecord) ValidationHealthInput {
	return ValidationHealthInput{
		GeneratedAt: now,
		AgencyID:    "demo-agency",
		Discovery:   validationHealthDiscovery(now),
		ToolingStatusByValidator: map[string]string{
			ValidationHealthStaticValidatorID:   ValidationHealthStatusConfigured,
			ValidationHealthRealtimeValidatorID: ValidationHealthStatusConfigured,
		},
		Records: records,
	}
}

func validationHealthRecords(total int, now time.Time) []ValidationReportRecord {
	records := make([]ValidationReportRecord, 0, total)
	for i := 0; i < total; i++ {
		feedType := validationHealthFeedOrder[i%len(validationHealthFeedOrder)]
		createdAt := now.Add(time.Duration(i) * time.Second)
		records = append(records, recordForHealth(int64(i+1), feedType, "old-feed", "warning", 0, 1, createdAt))
	}
	for i, feedType := range validationHealthFeedOrder {
		records = append(records, recordForHealth(int64(total+i+1), feedType, "feed-v1", "passed", 0, 0, now.Add(24*time.Hour+time.Duration(i)*time.Second)))
	}
	return records
}

func validationHealthDiscovery(now time.Time) FeedDiscovery {
	return FeedDiscovery{AgencyID: "demo-agency", Feeds: []FeedMetadata{
		{FeedType: "schedule", CanonicalPublicURL: "https://feeds.example.org/public/gtfs/schedule.zip", ActiveFeedVersionID: "feed-v1", RevisionTimestamp: &now},
		{FeedType: "vehicle_positions", CanonicalPublicURL: "https://feeds.example.org/public/gtfsrt/vehicle_positions.pb", ActiveFeedVersionID: "feed-v1", RevisionTimestamp: &now},
		{FeedType: "trip_updates", CanonicalPublicURL: "https://feeds.example.org/public/gtfsrt/trip_updates.pb", ActiveFeedVersionID: "feed-v1", RevisionTimestamp: &now},
		{FeedType: "alerts", CanonicalPublicURL: "https://feeds.example.org/public/gtfsrt/alerts.pb", ActiveFeedVersionID: "feed-v1", RevisionTimestamp: &now},
	}}
}

func recordForHealth(id int64, feedType, feedVersionID, status string, errors, warnings int, createdAt time.Time) ValidationReportRecord {
	return ValidationReportRecord{
		ID:        id,
		CreatedAt: createdAt,
		Result: ValidationResult{
			AgencyID:      "demo-agency",
			FeedType:      feedType,
			FeedVersionID: feedVersionID,
			ValidatorName: ValidatorNameForHealthID(ValidatorIDForHealthFeed(feedType)),
			Status:        status,
			ErrorCount:    errors,
			WarningCount:  warnings,
			Report:        map[string]any{"raw_report": "must not be retained", "stdout": "must not leak", "argv": []any{"/tmp/private"}},
		},
	}
}

func assertValidationHealthJSONAllowlist(t *testing.T, summary ValidationHealthSummary) {
	t.Helper()
	payload, err := json.Marshal(summary)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatal(err)
	}
	wantTop := map[string]bool{
		"generated_at": true, "agency_id": true, "overall_status": true, "tooling_status": true, "feeds": true,
		"external_evidence_created": true, "consumer_statuses_changed": true, "compliance_claimed": true, "production_readiness_claimed": true,
	}
	for key := range decoded {
		if !wantTop[key] {
			t.Fatalf("unexpected top-level field %q in %s", key, payload)
		}
	}
	feeds := decoded["feeds"].([]any)
	wantRow := map[string]bool{
		"feed_type": true, "validator_id": true, "validator_name": true, "tooling_status": true, "artifact_status": true,
		"latest_result_status": true, "latest_result_at": true, "active_feed_version_id": true, "latest_result_feed_version_id": true,
		"stale_status": true, "health_status": true, "next_action": true, "claim_boundary": true,
	}
	for _, item := range feeds {
		row := item.(map[string]any)
		for key := range row {
			if !wantRow[key] {
				t.Fatalf("unexpected row field %q in %s", key, payload)
			}
		}
	}
}

func assertNoValidationHealthLeakage(t *testing.T, summary ValidationHealthSummary) {
	t.Helper()
	payload, err := json.Marshal(summary)
	if err != nil {
		t.Fatal(err)
	}
	body := string(payload)
	for _, forbidden := range []string{"raw_report", "stdout", "stderr", "argv", "Authorization", "Bearer", "TOKEN=", "SECRET", "PASSWORD=", "postgres://", "/tmp/private", "admin_session", "Cookie"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("validation health leaked %q in %s", forbidden, body)
		}
	}
}
