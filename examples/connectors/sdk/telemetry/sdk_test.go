package telemetrysdk

import (
	"testing"
	"time"
)

func TestNormalizeBatchEmitsDryRunEventsAndDropsUnsafeObservations(t *testing.T) {
	now := time.Date(2026, 5, 10, 15, 1, 0, 0, time.UTC)
	summary := NormalizeBatch([]Observation{
		{
			AgencyID:   "agency-demo",
			DeviceID:   "device-1",
			VehicleID:  "vehicle-1",
			ObservedAt: now.Add(-30 * time.Second),
			Latitude:   34.0522,
			Longitude:  -118.2437,
			Quality:    0.95,
		},
		{
			AgencyID:   "agency-demo",
			DeviceID:   "device-low",
			VehicleID:  "vehicle-low",
			ObservedAt: now.Add(-20 * time.Second),
			Latitude:   34.0523,
			Longitude:  -118.2438,
			Quality:    0.2,
		},
	}, Options{Now: now, Source: "test-source"})

	if len(summary.Events) != 1 {
		t.Fatalf("events = %d, want 1: %+v", len(summary.Events), summary.Events)
	}
	event := summary.Events[0]
	if event.Source != "test-source" || !event.DryRun || event.NetworkSend {
		t.Fatalf("event flags/source = %+v, want dry-run no-send test-source", event)
	}
	if len(summary.Drops) != 1 || summary.Drops[0].Reason != ReasonLowQuality {
		t.Fatalf("drops = %+v, want low-quality drop", summary.Drops)
	}
}

func TestNormalizeBatchRejectsTimingIdentityAndCoordinateProblems(t *testing.T) {
	now := time.Date(2026, 5, 10, 15, 1, 0, 0, time.UTC)
	summary := NormalizeBatch([]Observation{
		{DeviceID: "missing", ObservedAt: now, Latitude: 34, Longitude: -118, Quality: 1},
		{AgencyID: "agency-demo", DeviceID: "zero-time", VehicleID: "vehicle", Latitude: 34, Longitude: -118, Quality: 1},
		{AgencyID: "agency-demo", DeviceID: "future", VehicleID: "vehicle", ObservedAt: now.Add(time.Minute), Latitude: 34, Longitude: -118, Quality: 1},
		{AgencyID: "agency-demo", DeviceID: "stale", VehicleID: "vehicle", ObservedAt: now.Add(-3 * time.Minute), Latitude: 34, Longitude: -118, Quality: 1},
		{AgencyID: "agency-demo", DeviceID: "coordinate", VehicleID: "vehicle", ObservedAt: now, Latitude: 91, Longitude: -118, Quality: 1},
	}, Options{Now: now})

	var reasons []string
	for _, drop := range summary.Drops {
		reasons = append(reasons, drop.Reason)
	}
	want := []string{ReasonMissingIdentity, ReasonInvalidTimestamp, ReasonFutureTimestamp, ReasonStaleObservation, ReasonInvalidCoordinate}
	if len(summary.Events) != 0 || len(reasons) != len(want) {
		t.Fatalf("events = %+v drops = %+v, want only drops %v", summary.Events, summary.Drops, want)
	}
	for i, reason := range want {
		if reasons[i] != reason {
			t.Fatalf("reason[%d] = %q, want %q (all reasons %v)", i, reasons[i], reason, reasons)
		}
	}
}
