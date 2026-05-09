package compliance

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestReliabilitySummaryShapeOrderMissingAndClaimFlags(t *testing.T) {
	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	endpoint := true
	freshness := 30.0
	latency := 120.0
	matched := 80.0
	summary := BuildReliabilitySummary(ReliabilityInput{
		GeneratedAt: now,
		AgencyID:    "demo-agency",
		FeedHealthRecords: []ReliabilityFeedHealthRecord{
			{FeedType: "vehicle_positions", SnapshotAt: now.Add(-time.Minute), EndpointAvailable: &endpoint, FreshnessSeconds: &freshness, GenerationLatencyMS: &latency, MatchedVehiclePercent: &matched},
		},
		Incidents: NormalizeReliabilityIncidentRollup(now, 0, nil, nil, nil, nil, nil, 10),
	})
	if summary.AgencyID != "demo-agency" || summary.GeneratedAt != now {
		t.Fatalf("summary metadata = %+v", summary)
	}
	if len(summary.Feeds) != 4 {
		t.Fatalf("feeds len = %d, want 4", len(summary.Feeds))
	}
	wantOrder := []string{"schedule", "vehicle_positions", "trip_updates", "alerts"}
	for i, want := range wantOrder {
		if summary.Feeds[i].FeedType != want {
			t.Fatalf("feed order[%d] = %q, want %q", i, summary.Feeds[i].FeedType, want)
		}
		assertReliabilityStatusAllowed(t, summary.Feeds[i].Status)
	}
	if summary.Feeds[0].Status == ReliabilityStatusOK || summary.Feeds[2].Status == ReliabilityStatusOK || summary.Feeds[3].Status == ReliabilityStatusOK {
		t.Fatalf("missing data became ok: %+v", summary.Feeds)
	}
	if summary.Feeds[1].Status != ReliabilityStatusOK {
		t.Fatalf("vehicle_positions status = %s, want ok", summary.Feeds[1].Status)
	}
	assertReliabilitySummaryFlagsFalse(t, summary)
	for _, section := range []ReliabilitySection{summary.BackupRestore, summary.Alerting, summary.AvailabilitySampling, summary.LongRunningOperations} {
		if section.Status == ReliabilityStatusOK {
			t.Fatalf("missing safe source became ok: %+v", section)
		}
		assertReliabilityStatusAllowed(t, section.Status)
	}
}

func TestReliabilityIncidentRollupSanitizesCapsAndCounts(t *testing.T) {
	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	oldest := now.Add(-2 * time.Hour)
	var recent []ReliabilityIncidentItem
	for i := 0; i < 20; i++ {
		recent = append(recent, ReliabilityIncidentItem{
			ID:       int64(i + 1),
			Type:     "prediction_review",
			Severity: "warning",
			Status:   "open",
			OpenedAt: now.Add(time.Duration(-i) * time.Minute),
			Title:    "raw details_json token https://private.example/host",
			Category: "payload_json",
		})
	}
	rollup := NormalizeReliabilityIncidentRollup(now, 20, map[string]int{"open": 12, "resolved": 8}, map[string]int{"warning": 20}, map[string]int{"prediction_review": 20}, &oldest, recent, 10)
	if rollup.Total != 20 || rollup.CountsByStatus["open"] != 12 || rollup.CountsBySeverity["warning"] != 20 || rollup.CountsByType["prediction_review"] != 20 {
		t.Fatalf("counts = %+v", rollup)
	}
	if rollup.OldestOpenAgeSeconds == nil || *rollup.OldestOpenAgeSeconds != 7200 {
		t.Fatalf("oldest open age = %+v, want 7200", rollup.OldestOpenAgeSeconds)
	}
	if len(rollup.Recent) != 10 {
		t.Fatalf("recent len = %d, want cap 10", len(rollup.Recent))
	}
	body := fmt.Sprintf("%+v", rollup)
	for _, forbidden := range []string{"details_json", "token", "https://", "payload_json", "private.example"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("incident rollup leaked %q: %s", forbidden, body)
		}
	}
}

func assertReliabilitySummaryFlagsFalse(t *testing.T, summary ReliabilitySummary) {
	t.Helper()
	flags := summary.ClaimFlags
	if flags.ExternalEvidenceCreated || flags.FinalRootEvidenceCreated || flags.ConsumerStatusesChanged ||
		flags.ComplianceClaimed || flags.ProductionReadinessClaimed || flags.SLAClaimed ||
		flags.UptimeGuaranteeClaimed || flags.HostedSaaSClaimed || flags.AgencyAdoptionClaimed ||
		flags.ConsumerAcceptanceClaimed || flags.VendorCompatibilityClaimed || flags.ProductionGradeETAClaimed {
		t.Fatalf("claim flags must all be false: %+v", flags)
	}
}

func assertReliabilityStatusAllowed(t *testing.T, status string) {
	t.Helper()
	switch status {
	case ReliabilityStatusOK, ReliabilityStatusNeedsReview, ReliabilityStatusMissing, ReliabilityStatusUnknown, ReliabilityStatusUnhealthy:
	default:
		t.Fatalf("status %q is not allowed", status)
	}
}
