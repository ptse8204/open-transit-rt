package alerts

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"

	domainalerts "open-transit-rt/internal/alerts"
)

func TestAlertsSnapshotBuildsValidPublishedAlertFeed(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	activeStart := generatedAt.Add(-time.Hour)
	repo := fakeAlertRepo{alerts: []domainalerts.Alert{{
		AgencyID:        "demo-agency",
		AlertKey:        "canceled:trip-10:20260421:08:00:00",
		Status:          domainalerts.StatusPublished,
		Cause:           "other_cause",
		Effect:          "no_service",
		HeaderText:      "Trip 10 canceled",
		DescriptionText: "Trip 10 is canceled today.",
		ActiveStart:     &activeStart,
		FeedVersionID:   "feed-demo",
		Entities: []domainalerts.InformedEntity{{
			AgencyID:  "demo-agency",
			RouteID:   "route-10",
			TripID:    "trip-10",
			StartDate: "20260421",
			StartTime: "08:00:00",
		}},
	}}}
	builder, err := NewBuilder(repo, &fakeHealthRepo{}, Config{AgencyID: "demo-agency"})
	if err != nil {
		t.Fatalf("new builder: %v", err)
	}
	snapshot, err := builder.Snapshot(context.Background(), generatedAt)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	payload, err := snapshot.MarshalProto()
	if err != nil {
		t.Fatalf("marshal proto: %v", err)
	}
	var message gtfsrt.FeedMessage
	if err := proto.Unmarshal(payload, &message); err != nil {
		t.Fatalf("unmarshal feed: %v", err)
	}
	if err := proto.CheckInitialized(&message); err != nil {
		t.Fatalf("protobuf not initialized: %v", err)
	}
	if message.GetHeader().GetTimestamp() != uint64(generatedAt.Unix()) {
		t.Fatalf("header timestamp = %d, want generated_at", message.GetHeader().GetTimestamp())
	}
	if len(message.Entity) != 1 || message.Entity[0].GetAlert() == nil {
		t.Fatalf("entities = %+v, want one alert entity", message.Entity)
	}
	alert := message.Entity[0].GetAlert()
	if alert.GetEffect() != gtfsrt.Alert_NO_SERVICE || alert.GetCause() != gtfsrt.Alert_OTHER_CAUSE {
		t.Fatalf("alert cause/effect = %s/%s, want OTHER_CAUSE/NO_SERVICE", alert.GetCause(), alert.GetEffect())
	}
	if len(alert.GetInformedEntity()) != 1 || alert.GetInformedEntity()[0].GetTrip().GetTripId() != "trip-10" {
		t.Fatalf("informed entities = %+v, want trip selector", alert.GetInformedEntity())
	}
}

func TestAlertsDebugIncludesOutputAndPersistenceStatus(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	snapshot := Snapshot{
		AgencyID:                      "demo-agency",
		GeneratedAt:                   generatedAt,
		Status:                        StatusEmpty,
		Reason:                        "no_published_active_alerts",
		DiagnosticsPersistenceOutcome: "stored",
	}
	payload, err := snapshot.MarshalDebugJSON()
	if err != nil {
		t.Fatalf("marshal debug: %v", err)
	}
	var debug Debug
	if err := json.Unmarshal(payload, &debug); err != nil {
		t.Fatalf("decode debug: %v", err)
	}
	if debug.Status != StatusEmpty || debug.AlertsOutput != 0 || debug.DiagnosticsPersistenceOutcome != "stored" {
		t.Fatalf("debug = %+v, want empty stored diagnostics", debug)
	}
}

type fakeAlertRepo struct {
	alerts []domainalerts.Alert
}

func (f fakeAlertRepo) ListAlerts(context.Context, domainalerts.ListFilter) ([]domainalerts.Alert, error) {
	return f.alerts, nil
}

type fakeHealthRepo struct{}

func (f *fakeHealthRepo) SaveAlertsSnapshot(context.Context, HealthRecord) error {
	return nil
}
