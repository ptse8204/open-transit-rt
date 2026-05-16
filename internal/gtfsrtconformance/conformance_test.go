package gtfsrtconformance

import (
	"testing"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

func TestVehiclePositionsConformanceAcceptsLocalProfile(t *testing.T) {
	payload := marshalTestFeed(t, &gtfsrt.FeedMessage{
		Header: testHeader(t, time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)),
		Entity: []*gtfsrt.FeedEntity{{
			Id: proto.String("bus-1"),
			Vehicle: &gtfsrt.VehiclePosition{
				Vehicle:   &gtfsrt.VehicleDescriptor{Id: proto.String("bus-1")},
				Position:  &gtfsrt.Position{Latitude: proto.Float32(37.1), Longitude: proto.Float32(-122.1)},
				Timestamp: proto.Uint64(uint64(time.Date(2026, 5, 16, 11, 59, 30, 0, time.UTC).Unix())),
				Trip: &gtfsrt.TripDescriptor{
					TripId:    proto.String("trip-1"),
					RouteId:   proto.String("route-1"),
					StartDate: proto.String("20260516"),
				},
			},
		}},
	})
	report := CheckPayload(FeedVehiclePositions, "synthetic.pb", payload, Options{GeneratedAt: time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)})
	if report.Status != StatusOK || report.EntityCount != 1 {
		t.Fatalf("report = %+v, want ok with one entity", report)
	}
}

func TestConformanceFlagsDuplicateEntityIDs(t *testing.T) {
	payload := marshalTestFeed(t, &gtfsrt.FeedMessage{
		Header: testHeader(t, time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)),
		Entity: []*gtfsrt.FeedEntity{
			{Id: proto.String("dup"), Vehicle: validVehicle("bus-1")},
			{Id: proto.String("dup"), Vehicle: validVehicle("bus-2")},
		},
	})
	report := CheckPayload(FeedVehiclePositions, "synthetic.pb", payload, Options{GeneratedAt: time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)})
	if report.Status != StatusFailed {
		t.Fatalf("status = %s, want failed for duplicate entity id; checks=%+v", report.Status, report.Checks)
	}
}

func TestTripUpdatesMissingStartDateNeedsReview(t *testing.T) {
	payload := marshalTestFeed(t, &gtfsrt.FeedMessage{
		Header: testHeader(t, time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)),
		Entity: []*gtfsrt.FeedEntity{{
			Id: proto.String("tu-1"),
			TripUpdate: &gtfsrt.TripUpdate{
				Trip: &gtfsrt.TripDescriptor{TripId: proto.String("trip-1")},
				StopTimeUpdate: []*gtfsrt.TripUpdate_StopTimeUpdate{{
					StopSequence: proto.Uint32(1),
					Arrival:      &gtfsrt.TripUpdate_StopTimeEvent{Time: proto.Int64(time.Date(2026, 5, 16, 12, 10, 0, 0, time.UTC).Unix())},
				}},
			},
		}},
	})
	report := CheckPayload(FeedTripUpdates, "synthetic.pb", payload, Options{GeneratedAt: time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)})
	if report.Status != StatusNeedsReview {
		t.Fatalf("status = %s, want needs_review for missing start_date; checks=%+v", report.Status, report.Checks)
	}
}

func TestAlertsMissingHeaderTextFails(t *testing.T) {
	payload := marshalTestFeed(t, &gtfsrt.FeedMessage{
		Header: testHeader(t, time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)),
		Entity: []*gtfsrt.FeedEntity{{
			Id: proto.String("alert-1"),
			Alert: &gtfsrt.Alert{
				InformedEntity: []*gtfsrt.EntitySelector{{AgencyId: proto.String("demo-agency")}},
				ActivePeriod:   []*gtfsrt.TimeRange{{Start: proto.Uint64(uint64(time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC).Unix()))}},
			},
		}},
	})
	report := CheckPayload(FeedAlerts, "synthetic.pb", payload, Options{GeneratedAt: time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)})
	if report.Status != StatusFailed {
		t.Fatalf("status = %s, want failed for missing alert header text; checks=%+v", report.Status, report.Checks)
	}
}

func TestBuildReportOrdersFeedsAndKeepsClaimFlagsFalse(t *testing.T) {
	generatedAt := time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)
	empty := marshalTestFeed(t, &gtfsrt.FeedMessage{Header: testHeader(t, generatedAt)})
	report := BuildReport([]FeedInput{
		{FeedType: FeedAlerts, Source: "alerts.pb", Payload: empty},
		{FeedType: FeedVehiclePositions, Source: "vehicle_positions.pb", Payload: empty},
		{FeedType: FeedTripUpdates, Source: "trip_updates.pb", Payload: empty},
	}, Options{GeneratedAt: generatedAt})
	if report.OverallStatus != StatusOK || len(report.Feeds) != 3 {
		t.Fatalf("report = %+v, want ok with three empty local feeds", report)
	}
	if report.Feeds[0].FeedType != FeedVehiclePositions || report.Feeds[1].FeedType != FeedTripUpdates || report.Feeds[2].FeedType != FeedAlerts {
		t.Fatalf("feed order = %+v", report.Feeds)
	}
	if report.ClaimFlags != (ClaimFlags{}) {
		t.Fatalf("claim flags = %+v, want all false", report.ClaimFlags)
	}
}

func validVehicle(id string) *gtfsrt.VehiclePosition {
	return &gtfsrt.VehiclePosition{
		Vehicle:   &gtfsrt.VehicleDescriptor{Id: proto.String(id)},
		Position:  &gtfsrt.Position{Latitude: proto.Float32(37.1), Longitude: proto.Float32(-122.1)},
		Timestamp: proto.Uint64(uint64(time.Date(2026, 5, 16, 11, 59, 30, 0, time.UTC).Unix())),
	}
}

func testHeader(t *testing.T, generatedAt time.Time) *gtfsrt.FeedHeader {
	t.Helper()
	incrementality := gtfsrt.FeedHeader_FULL_DATASET
	return &gtfsrt.FeedHeader{
		GtfsRealtimeVersion: proto.String("2.0"),
		Incrementality:      &incrementality,
		Timestamp:           proto.Uint64(uint64(generatedAt.Unix())),
	}
}

func marshalTestFeed(t *testing.T, message *gtfsrt.FeedMessage) []byte {
	t.Helper()
	payload, err := proto.Marshal(message)
	if err != nil {
		t.Fatalf("marshal test feed: %v", err)
	}
	return payload
}
