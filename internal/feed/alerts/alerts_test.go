package alerts

import (
	"encoding/json"
	"testing"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

func TestAlertsSnapshotIsValidEmptyFeedWithExplicitTimestamp(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	builder, err := NewBuilder("demo-agency")
	if err != nil {
		t.Fatalf("new builder: %v", err)
	}
	snapshot := builder.Snapshot(generatedAt)
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
	if message.GetHeader().GetGtfsRealtimeVersion() != "2.0" || message.GetHeader().GetIncrementality() != gtfsrt.FeedHeader_FULL_DATASET {
		t.Fatalf("header = %+v, want GTFS-RT 2.0 FULL_DATASET", message.GetHeader())
	}
	if len(message.Entity) != 0 {
		t.Fatalf("entities = %d, want empty deferred alerts feed", len(message.Entity))
	}
}

func TestAlertsDebugIsJSONOnlyDeferredStatus(t *testing.T) {
	generatedAt := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	snapshot := Snapshot{
		AgencyID:    "demo-agency",
		GeneratedAt: generatedAt,
		Status:      StatusDeferred,
		Reason:      ReasonAlertsAuthoringMissing,
	}
	payload, err := snapshot.MarshalDebugJSON()
	if err != nil {
		t.Fatalf("marshal debug: %v", err)
	}
	var debug Debug
	if err := json.Unmarshal(payload, &debug); err != nil {
		t.Fatalf("decode debug: %v", err)
	}
	if debug.Status != StatusDeferred || debug.Reason != ReasonAlertsAuthoringMissing || debug.Diagnostics != "json_only" {
		t.Fatalf("debug = %+v, want JSON-only deferred status", debug)
	}
}
