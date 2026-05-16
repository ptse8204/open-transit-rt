package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

func TestRunValidVehiclePositionsJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vehicle_positions.pb")
	writeTestFeed(t, path, &gtfsrt.FeedMessage{
		Header: testCommandHeader(time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)),
		Entity: []*gtfsrt.FeedEntity{{
			Id: proto.String("bus-1"),
			Vehicle: &gtfsrt.VehiclePosition{
				Vehicle:   &gtfsrt.VehicleDescriptor{Id: proto.String("bus-1")},
				Position:  &gtfsrt.Position{Latitude: proto.Float32(37.1), Longitude: proto.Float32(-122.1)},
				Timestamp: proto.Uint64(uint64(time.Date(2026, 5, 16, 11, 59, 30, 0, time.UTC).Unix())),
			},
		}},
	})
	var stdout, stderr bytes.Buffer
	code := run([]string{"--vehicle-positions", path, "--generated-at", "2026-05-16T12:00:00Z"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	for _, want := range []string{`"overall_status": "ok"`, `"feed_type": "vehicle_positions"`, `"external_evidence_created": false`} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q: %s", want, stdout.String())
		}
	}
}

func TestRunFailedConformanceReturnsOne(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.pb")
	writeTestFeed(t, path, &gtfsrt.FeedMessage{
		Header: testCommandHeader(time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)),
		Entity: []*gtfsrt.FeedEntity{{Id: proto.String("alert-1"), Alert: &gtfsrt.Alert{}}},
	})
	var stdout, stderr bytes.Buffer
	code := run([]string{"--alerts", path, "--generated-at", "2026-05-16T12:00:00Z", "--output", "text"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), "alert_header_text") {
		t.Fatalf("stdout missing alert failure: %s", stdout.String())
	}
}

func TestRunRejectsEvidenceLikePaths(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--vehicle-positions", "docs/evidence/captured/example.pb"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code=%d stderr=%s stdout=%s, want read refusal", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), "refusing to read evidence-like path") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}

func writeTestFeed(t *testing.T, path string, message *gtfsrt.FeedMessage) {
	t.Helper()
	payload, err := proto.Marshal(message)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func testCommandHeader(generatedAt time.Time) *gtfsrt.FeedHeader {
	incrementality := gtfsrt.FeedHeader_FULL_DATASET
	return &gtfsrt.FeedHeader{
		GtfsRealtimeVersion: proto.String("2.0"),
		Incrementality:      &incrementality,
		Timestamp:           proto.Uint64(uint64(generatedAt.Unix())),
	}
}
