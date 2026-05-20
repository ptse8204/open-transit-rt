package compliance

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestBuildGTFSQualityTriageLargeCanonicalReport(t *testing.T) {
	for _, total := range []int{10000, 50000} {
		t.Run(fmt.Sprintf("%d", total), func(t *testing.T) {
			record := largeCanonicalRecord(total)
			triage := BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &record})
			if len(triage.Canonical.Groups) > GTFSQualityMaxGroups {
				t.Fatalf("groups = %d, want capped to %d", len(triage.Canonical.Groups), GTFSQualityMaxGroups)
			}
			counted := triage.Canonical.OverflowCount
			for _, group := range triage.Canonical.Groups {
				counted += group.Count
				if len(group.Samples) > GTFSQualityMaxSamples {
					t.Fatalf("samples = %d, want capped", len(group.Samples))
				}
				if group.OverflowCount != group.Count-len(group.Samples) {
					t.Fatalf("group overflow = %d, want %d", group.OverflowCount, group.Count-len(group.Samples))
				}
			}
			if counted != total {
				t.Fatalf("counted notices = %d, want %d", counted, total)
			}
			again := BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &record})
			if !reflect.DeepEqual(triage.Canonical.Groups, again.Canonical.Groups) {
				t.Fatalf("triage ordering is not deterministic")
			}
			rendered := fmt.Sprintf("%+v", triage)
			for _, forbidden := range []string{"raw_report", "stdout", "stderr", "argv", "/tmp/private"} {
				if strings.Contains(rendered, forbidden) {
					t.Fatalf("triage leaked %q", forbidden)
				}
			}
		})
	}
}

func TestGTFSQualityTriageMalformedReports(t *testing.T) {
	cases := []map[string]any{
		{},
		{"raw_report": map[string]any{}},
		{"raw_report": map[string]any{"notices": "not-array"}},
		{"raw_report": map[string]any{"notices": []any{map[string]any{"message": "missing code"}}}},
		{"raw_report": map[string]any{"notices": []any{map[string]any{"code": "x", "nested": map[string]any{"raw_report": "secret"}}}}},
	}
	for i, report := range cases {
		record := ValidationReportRecord{Result: ValidationResult{AgencyID: "demo", FeedType: "schedule", ValidatorName: CanonicalStaticValidatorName, Status: "warning", WarningCount: 1, Report: report}, CreatedAt: time.Now()}
		triage := BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &record})
		if len(triage.Canonical.Groups) == 0 {
			t.Fatalf("case %d produced no fallback group", i)
		}
		if strings.Contains(fmt.Sprintf("%+v", triage), "map[raw_report") {
			t.Fatalf("case %d retained raw nested object: %+v", i, triage)
		}
	}
}

func TestGTFSQualityTriageTaxonomyAndSeverityOrdering(t *testing.T) {
	record := ValidationReportRecord{Result: ValidationResult{AgencyID: "demo", FeedType: "schedule", ValidatorName: CanonicalStaticValidatorName, Status: "warning", Report: map[string]any{"raw_report": map[string]any{"notices": []any{
		notice("unused_shape", "INFO", "shapes.txt", "unused shape"),
		notice("calendar_dates_service_id_foreign_key", "ERROR", "calendar_dates.txt", "missing foreign-key references"),
		notice("route_short_name_too_long", "WARNING", "routes.txt", "long name"),
		notice("duplicate_route_id", "ERROR", "routes.txt", "duplicate IDs"),
		notice("invalid_stop_time", "ERROR", "stop_times.txt", "bad stop_times"),
		notice("expired_calendar", "WARNING", "calendar.txt", "calendar/service-date issues"),
		notice("calendar_no_service", "WARNING", "calendar.txt", "calendar/service-date issues"),
		notice("shape_dist_traveled_decreases", "WARNING", "shapes.txt", "shape ordering"),
		notice("frequency_headway_invalid", "WARNING", "frequencies.txt", "frequency issues"),
		notice("block_id_gap", "WARNING", "trips.txt", "block transition"),
		notice("mystery_notice", "WARNING", "agency.txt", "unknown notices"),
	}}}}, CreatedAt: time.Now()}
	triage := BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &record})
	families := map[string]bool{}
	lastRank := -1
	for _, group := range triage.Canonical.Groups {
		families[group.Family] = true
		rank := severityRank(group.Severity)
		if rank < lastRank {
			t.Fatalf("severity ordering regressed: %+v", triage.Canonical.Groups)
		}
		lastRank = rank
		if group.Source != GTFSQualitySourceCanonicalValidator {
			t.Fatalf("group source = %q", group.Source)
		}
	}
	for _, family := range []string{"expired_calendar", "route_short_name_too_long", "unused_shape", "missing_or_foreign_key_reference", "bad_stop_times", "duplicate_ids", "calendar_service_dates", "shape_ordering", "frequency_issues", "block_transition_issues", "unknown"} {
		if !families[family] {
			t.Fatalf("missing family %s in %+v", family, triage.Canonical.Groups)
		}
	}
}

func TestGTFSQualityTriagePracticalCategoriesRiskAndGroupedCounts(t *testing.T) {
	record := ValidationReportRecord{Result: ValidationResult{AgencyID: "demo", FeedType: "schedule", ValidatorName: CanonicalStaticValidatorName, Status: "warning", Report: map[string]any{"raw_report": map[string]any{"notices": []any{
		map[string]any{"code": "missing_required_file", "severity": "ERROR", "fileName": "stop_times.txt", "message": "missing required file", "totalNotices": 4},
		map[string]any{"code": "agency_timezone_missing", "severity": "WARNING", "fileName": "agency.txt", "message": "agency timezone missing"},
		map[string]any{"code": "feed_info_license_missing", "severity": "WARNING", "fileName": "feed_info.txt", "message": "license missing"},
		map[string]any{"code": "route_color_invalid", "severity": "WARNING", "fileName": "routes.txt", "message": "route color invalid"},
		map[string]any{"code": "stop_lat_invalid", "severity": "WARNING", "fileName": "stops.txt", "message": "stop_lat invalid"},
	}}}}, CreatedAt: time.Now()}
	triage := BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &record})
	families := map[string]GTFSQualityGroup{}
	for _, group := range triage.Canonical.Groups {
		families[group.Family] = group
		if group.RiskLevel == "" {
			t.Fatalf("group missing risk level: %+v", group)
		}
	}
	for _, family := range []string{"missing_required_file", "agency_metadata", "license_contact_metadata", "route_metadata", "stop_location"} {
		if _, ok := families[family]; !ok {
			t.Fatalf("missing family %s in %+v", family, triage.Canonical.Groups)
		}
	}
	if got := families["missing_required_file"].Count; got != 4 {
		t.Fatalf("grouped totalNotices count = %d, want 4", got)
	}
	if got := families["missing_required_file"].RiskLevel; got != "blocks import or reliable feed use" {
		t.Fatalf("missing required risk = %q", got)
	}
}

func TestGTFSQualityTriageStaleActiveFeed(t *testing.T) {
	revision := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	current := ValidationReportRecord{Result: ValidationResult{AgencyID: "demo", FeedType: "schedule", FeedVersionID: "feed-v2", ValidatorName: CanonicalStaticValidatorName, Status: "passed", Report: map[string]any{"raw_report": map[string]any{"notices": []any{}}}}, CreatedAt: revision.Add(time.Minute)}
	if got := BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &current, ActiveFeedVersionID: "feed-v2", ActiveFeedRevisionTime: &revision}).Canonical; got.IsStale || got.Status == GTFSQualityNeedsReview {
		t.Fatalf("current result marked stale: %+v", got)
	}
	old := current
	old.CreatedAt = revision.Add(-time.Minute)
	if got := BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &old, ActiveFeedVersionID: "feed-v2", ActiveFeedRevisionTime: &revision}).Canonical; !got.IsStale || got.Status != GTFSQualityNeedsReview {
		t.Fatalf("old result = %+v, want stale needs_review", got)
	}
	otherFeed := current
	otherFeed.Result.FeedVersionID = "feed-v1"
	if got := BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &otherFeed, ActiveFeedVersionID: "feed-v2", ActiveFeedRevisionTime: &revision}).Canonical; !got.IsStale || got.Status != GTFSQualityNeedsReview {
		t.Fatalf("wrong feed result = %+v, want stale needs_review", got)
	}
}

func TestGTFSQualityTriageHostileReport(t *testing.T) {
	record := hostileCanonicalRecord(2000)
	triage := BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &record})
	if len(triage.Canonical.Groups) > GTFSQualityMaxGroups {
		t.Fatalf("groups = %d, want capped", len(triage.Canonical.Groups))
	}
	rendered := fmt.Sprintf("%+v", triage)
	for _, forbidden := range []string{"raw_report", "stdout", "stderr", "argv", "/tmp/private", strings.Repeat("x", 1000)} {
		if strings.Contains(rendered, forbidden) {
			t.Fatalf("hostile triage leaked %q", forbidden)
		}
	}
	for _, group := range triage.Canonical.Groups {
		for _, sample := range group.Samples {
			if len(sample) > GTFSQualityMaxSampleLength+3 {
				t.Fatalf("sample length = %d, want capped", len(sample))
			}
		}
	}
	again := BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &record})
	if !reflect.DeepEqual(triage.Canonical.Groups, again.Canonical.Groups) {
		t.Fatalf("hostile report ordering is not deterministic")
	}
}

func TestGTFSQualityTriageRedactsSecretLikeSamples(t *testing.T) {
	record := ValidationReportRecord{Result: ValidationResult{AgencyID: "demo", FeedType: "schedule", ValidatorName: CanonicalStaticValidatorName, Status: "warning", WarningCount: 1, Report: map[string]any{"raw_report": map[string]any{"notices": []any{
		map[string]any{"code": "mystery_notice", "severity": "WARNING", "file": "agency.txt", "message": "Authorization: Bearer TOKEN=SECRET database_url=postgres://user:pass@localhost/db Cookie admin_session webhook https://private.example/hook /Users/private/report.json"},
	}}}}, CreatedAt: time.Now()}
	triage := BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &record})
	rendered := fmt.Sprintf("%+v", triage)
	for _, forbidden := range []string{"Authorization", "Bearer", "TOKEN=SECRET", "database_url", "postgres://", "Cookie", "admin_session", "webhook", "private.example", "/Users/private"} {
		if strings.Contains(rendered, forbidden) {
			t.Fatalf("secret-like sample leaked %q: %+v", forbidden, triage)
		}
	}
	for _, want := range []string{"{redacted_secret}", "{redacted_database}", "{redacted_url}", "{private_path}"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("redacted sample missing %q: %+v", want, triage)
		}
	}
}

func BenchmarkBuildGTFSQualityTriage(b *testing.B) {
	for _, total := range []int{10000, 50000} {
		record := largeCanonicalRecord(total)
		b.Run(fmt.Sprintf("%d", total), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &record})
			}
		})
	}
}

func BenchmarkBuildGTFSQualityTriageHostileReport(b *testing.B) {
	record := hostileCanonicalRecord(50000)
	for i := 0; i < b.N; i++ {
		_ = BuildGTFSQualityTriage(GTFSQualityTriageInput{Canonical: &record})
	}
}

func largeCanonicalRecord(total int) ValidationReportRecord {
	notices := make([]any, 0, total)
	codes := []string{"expired_calendar", "route_short_name_too_long", "unused_shape", "stop_times_arrival_time_missing", "duplicate_trip_id", "shape_dist_traveled_decreases", "frequency_headway_invalid", "block_id_gap"}
	severities := []string{"ERROR", "WARNING", "INFO"}
	for i := 0; i < total; i++ {
		notices = append(notices, notice(codes[i%len(codes)], severities[i%len(severities)], "file.txt", fmt.Sprintf("notice %d /tmp/private <script>", i)))
	}
	return ValidationReportRecord{ID: 1, Result: ValidationResult{AgencyID: "demo", FeedType: "schedule", FeedVersionID: "feed-v1", ValidatorName: CanonicalStaticValidatorName, Status: "warning", ErrorCount: total / 3, WarningCount: total / 3, InfoCount: total / 3, Report: map[string]any{"raw_report": map[string]any{"notices": notices}, "stdout": "secret stdout", "stderr": "secret stderr", "argv": []any{"/tmp/private/validator"}}}, CreatedAt: time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)}
}

func hostileCanonicalRecord(total int) ValidationReportRecord {
	notices := make([]any, 0, total+2)
	huge := strings.Repeat("x", 5000)
	for i := 0; i < total; i++ {
		notices = append(notices, map[string]any{
			"code":        fmt.Sprintf("unknown_notice_%03d", i%250),
			"severity":    "WARNING",
			"message":     "<script>alert(1)</script> " + huge,
			"entity":      map[string]any{"nested": map[string]any{"path": "/tmp/private/raw.zip", "raw_report": huge}},
			"output_path": "/tmp/private/output",
		})
	}
	notices = append(notices, "mixed", []any{"array"})
	return ValidationReportRecord{ID: 1, Result: ValidationResult{AgencyID: "demo", FeedType: "schedule", ValidatorName: CanonicalStaticValidatorName, Status: "warning", WarningCount: total, Report: map[string]any{"raw_report": map[string]any{"notices": notices}, "stdout": huge, "stderr": huge, "argv": []any{"/tmp/private/bin"}}}, CreatedAt: time.Now()}
}

func notice(code string, severity string, file string, message string) map[string]any {
	return map[string]any{"code": code, "severity": severity, "filename": file, "message": message}
}
