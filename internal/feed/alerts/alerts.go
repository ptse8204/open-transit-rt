package alerts

import (
	"encoding/json"
	"fmt"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"

	"open-transit-rt/internal/feed"
)

const (
	StatusDeferred               = "deferred"
	ReasonAlertsAuthoringMissing = "alerts_authoring_not_implemented"
)

type Snapshot struct {
	AgencyID    string
	GeneratedAt time.Time
	Status      string
	Reason      string
}

type Builder struct {
	agencyID string
}

func NewBuilder(agencyID string) (*Builder, error) {
	if agencyID == "" {
		return nil, fmt.Errorf("AGENCY_ID is required")
	}
	return &Builder{agencyID: agencyID}, nil
}

func (b *Builder) Snapshot(generatedAt time.Time) Snapshot {
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}
	return Snapshot{
		AgencyID:    b.agencyID,
		GeneratedAt: generatedAt.UTC(),
		Status:      StatusDeferred,
		Reason:      ReasonAlertsAuthoringMissing,
	}
}

func (s Snapshot) BuildProto() *gtfsrt.FeedMessage {
	timestamp := uint64(s.GeneratedAt.Unix())
	incrementality := gtfsrt.FeedHeader_FULL_DATASET
	return &gtfsrt.FeedMessage{
		Header: &gtfsrt.FeedHeader{
			GtfsRealtimeVersion: proto.String(feed.GTFSRealtimeVersion),
			Incrementality:      &incrementality,
			Timestamp:           &timestamp,
		},
		Entity: []*gtfsrt.FeedEntity{},
	}
}

func (s Snapshot) MarshalProto() ([]byte, error) {
	payload, err := proto.Marshal(s.BuildProto())
	if err != nil {
		return nil, fmt.Errorf("marshal alerts protobuf: %w", err)
	}
	return payload, nil
}

func (s Snapshot) MarshalDebugJSON() ([]byte, error) {
	payload, err := json.MarshalIndent(Debug{
		AgencyID:    s.AgencyID,
		GeneratedAt: s.GeneratedAt,
		Status:      s.Status,
		Reason:      s.Reason,
		Diagnostics: "json_only",
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal alerts debug json: %w", err)
	}
	return payload, nil
}

type Debug struct {
	AgencyID    string    `json:"agency_id"`
	GeneratedAt time.Time `json:"generated_at"`
	Status      string    `json:"status"`
	Reason      string    `json:"reason"`
	Diagnostics string    `json:"diagnostics"`
}
