package alerts

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"

	domainalerts "open-transit-rt/internal/alerts"
	"open-transit-rt/internal/feed"
)

const (
	StatusOK    = "ok"
	StatusEmpty = "empty"
)

type AlertRepository interface {
	ListAlerts(ctx context.Context, filter domainalerts.ListFilter) ([]domainalerts.Alert, error)
}

type HealthRepository interface {
	SaveAlertsSnapshot(ctx context.Context, record HealthRecord) error
}

type HealthRecord struct {
	AgencyID            string
	SnapshotAt          time.Time
	ActiveFeedVersionID string
	AlertsOutput        int
	Status              string
	Reason              string
}

type Config struct {
	AgencyID string
}

type Builder struct {
	alerts AlertRepository
	health HealthRepository
	config Config
}

func NewBuilder(alertRepo AlertRepository, healthRepo HealthRepository, config Config) (*Builder, error) {
	if alertRepo == nil {
		return nil, fmt.Errorf("alerts repository is required")
	}
	if healthRepo == nil {
		return nil, fmt.Errorf("alerts health repository is required")
	}
	if config.AgencyID == "" {
		return nil, fmt.Errorf("AGENCY_ID is required")
	}
	return &Builder{alerts: alertRepo, health: healthRepo, config: config}, nil
}

func (b *Builder) Snapshot(ctx context.Context, generatedAt time.Time) (Snapshot, error) {
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}
	generatedAt = generatedAt.UTC()
	alerts, err := b.alerts.ListAlerts(ctx, domainalerts.ListFilter{
		AgencyID:      b.config.AgencyID,
		PublishedOnly: true,
		At:            generatedAt,
		Limit:         5000,
	})
	if err != nil {
		return Snapshot{}, fmt.Errorf("list public alerts: %w", err)
	}
	snapshot := Snapshot{
		AgencyID:    b.config.AgencyID,
		GeneratedAt: generatedAt,
		Status:      StatusOK,
		Alerts:      normalizeAlerts(alerts),
	}
	if len(snapshot.Alerts) == 0 {
		snapshot.Status = StatusEmpty
		snapshot.Reason = "no_published_active_alerts"
	}
	if len(snapshot.Alerts) > 0 {
		snapshot.ActiveFeedVersionID = snapshot.Alerts[0].FeedVersionID
	}
	if err := b.health.SaveAlertsSnapshot(ctx, HealthRecord{
		AgencyID:            snapshot.AgencyID,
		SnapshotAt:          snapshot.GeneratedAt,
		ActiveFeedVersionID: snapshot.ActiveFeedVersionID,
		AlertsOutput:        len(snapshot.Alerts),
		Status:              snapshot.Status,
		Reason:              snapshot.Reason,
	}); err != nil {
		snapshot.DiagnosticsPersistenceOutcome = "failed"
		snapshot.DiagnosticsPersistenceError = err.Error()
	} else {
		snapshot.DiagnosticsPersistenceOutcome = "stored"
	}
	return snapshot, nil
}

type Snapshot struct {
	AgencyID                      string
	GeneratedAt                   time.Time
	ActiveFeedVersionID           string
	Status                        string
	Reason                        string
	Alerts                        []domainalerts.Alert
	DiagnosticsPersistenceOutcome string
	DiagnosticsPersistenceError   string
}

func (s Snapshot) BuildProto() (*gtfsrt.FeedMessage, error) {
	timestamp := uint64(s.GeneratedAt.Unix())
	incrementality := gtfsrt.FeedHeader_FULL_DATASET
	message := &gtfsrt.FeedMessage{
		Header: &gtfsrt.FeedHeader{
			GtfsRealtimeVersion: proto.String(feed.GTFSRealtimeVersion),
			Incrementality:      &incrementality,
			Timestamp:           &timestamp,
		},
		Entity: []*gtfsrt.FeedEntity{},
	}
	for _, alert := range normalizeAlerts(s.Alerts) {
		entity, err := buildEntity(alert)
		if err != nil {
			return nil, err
		}
		message.Entity = append(message.Entity, entity)
	}
	return message, nil
}

func (s Snapshot) MarshalProto() ([]byte, error) {
	message, err := s.BuildProto()
	if err != nil {
		return nil, err
	}
	payload, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("marshal alerts protobuf: %w", err)
	}
	return payload, nil
}

func (s Snapshot) MarshalDebugJSON() ([]byte, error) {
	payload, err := json.MarshalIndent(Debug{
		AgencyID:                      s.AgencyID,
		GeneratedAt:                   s.GeneratedAt,
		ActiveFeedVersionID:           s.ActiveFeedVersionID,
		Status:                        s.Status,
		Reason:                        s.Reason,
		AlertsOutput:                  len(s.Alerts),
		Alerts:                        s.Alerts,
		DiagnosticsPersistenceOutcome: s.DiagnosticsPersistenceOutcome,
		DiagnosticsPersistenceError:   s.DiagnosticsPersistenceError,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal alerts debug json: %w", err)
	}
	return payload, nil
}

type Debug struct {
	AgencyID                      string                 `json:"agency_id"`
	GeneratedAt                   time.Time              `json:"generated_at"`
	ActiveFeedVersionID           string                 `json:"active_feed_version_id"`
	Status                        string                 `json:"status"`
	Reason                        string                 `json:"reason,omitempty"`
	AlertsOutput                  int                    `json:"alerts_output"`
	Alerts                        []domainalerts.Alert   `json:"alerts"`
	DiagnosticsPersistenceOutcome string                 `json:"diagnostics_persistence_outcome"`
	DiagnosticsPersistenceError   string                 `json:"diagnostics_persistence_error,omitempty"`
	_                             map[string]interface{} `json:"-"`
}

func buildEntity(alert domainalerts.Alert) (*gtfsrt.FeedEntity, error) {
	if alert.AlertKey == "" {
		return nil, fmt.Errorf("alert key is required")
	}
	if alert.HeaderText == "" {
		return nil, fmt.Errorf("alert header text is required")
	}
	cause := alertCause(alert.Cause)
	effect := alertEffect(alert.Effect)
	rtAlert := &gtfsrt.Alert{
		Cause:          &cause,
		Effect:         &effect,
		ActivePeriod:   buildActivePeriods(alert),
		InformedEntity: buildInformedEntities(alert),
		HeaderText:     translated(alert.HeaderText),
	}
	if alert.DescriptionText != "" {
		rtAlert.DescriptionText = translated(alert.DescriptionText)
	}
	if alert.URL != "" {
		rtAlert.Url = translated(alert.URL)
	}
	return &gtfsrt.FeedEntity{Id: proto.String(alert.AlertKey), Alert: rtAlert}, nil
}

func buildActivePeriods(alert domainalerts.Alert) []*gtfsrt.TimeRange {
	if alert.ActiveStart == nil && alert.ActiveEnd == nil {
		return nil
	}
	rng := &gtfsrt.TimeRange{}
	if alert.ActiveStart != nil {
		rng.Start = proto.Uint64(uint64(alert.ActiveStart.UTC().Unix()))
	}
	if alert.ActiveEnd != nil {
		rng.End = proto.Uint64(uint64(alert.ActiveEnd.UTC().Unix()))
	}
	return []*gtfsrt.TimeRange{rng}
}

func buildInformedEntities(alert domainalerts.Alert) []*gtfsrt.EntitySelector {
	if len(alert.Entities) == 0 {
		return []*gtfsrt.EntitySelector{{AgencyId: proto.String(alert.AgencyID)}}
	}
	entities := append([]domainalerts.InformedEntity(nil), alert.Entities...)
	sort.SliceStable(entities, func(i int, j int) bool {
		return entitySortKey(entities[i]) < entitySortKey(entities[j])
	})
	result := make([]*gtfsrt.EntitySelector, 0, len(entities))
	for _, entity := range entities {
		selector := &gtfsrt.EntitySelector{AgencyId: proto.String(alert.AgencyID)}
		if entity.RouteID != "" {
			selector.RouteId = proto.String(entity.RouteID)
		}
		if entity.StopID != "" {
			selector.StopId = proto.String(entity.StopID)
		}
		if entity.TripID != "" {
			trip := &gtfsrt.TripDescriptor{TripId: proto.String(entity.TripID)}
			if entity.StartDate != "" {
				trip.StartDate = proto.String(entity.StartDate)
			}
			if entity.StartTime != "" {
				trip.StartTime = proto.String(entity.StartTime)
			}
			if entity.RouteID != "" {
				trip.RouteId = proto.String(entity.RouteID)
			}
			selector.Trip = trip
		}
		result = append(result, selector)
	}
	return result
}

func translated(text string) *gtfsrt.TranslatedString {
	return &gtfsrt.TranslatedString{
		Translation: []*gtfsrt.TranslatedString_Translation{{
			Text: proto.String(text),
		}},
	}
}

func alertCause(raw string) gtfsrt.Alert_Cause {
	switch strings.ToLower(raw) {
	case "other_cause":
		return gtfsrt.Alert_OTHER_CAUSE
	case "technical_problem":
		return gtfsrt.Alert_TECHNICAL_PROBLEM
	case "strike":
		return gtfsrt.Alert_STRIKE
	case "demonstration":
		return gtfsrt.Alert_DEMONSTRATION
	case "accident":
		return gtfsrt.Alert_ACCIDENT
	case "holiday":
		return gtfsrt.Alert_HOLIDAY
	case "weather":
		return gtfsrt.Alert_WEATHER
	case "maintenance":
		return gtfsrt.Alert_MAINTENANCE
	case "construction":
		return gtfsrt.Alert_CONSTRUCTION
	case "police_activity":
		return gtfsrt.Alert_POLICE_ACTIVITY
	case "medical_emergency":
		return gtfsrt.Alert_MEDICAL_EMERGENCY
	default:
		return gtfsrt.Alert_UNKNOWN_CAUSE
	}
}

func alertEffect(raw string) gtfsrt.Alert_Effect {
	switch strings.ToLower(raw) {
	case "no_service":
		return gtfsrt.Alert_NO_SERVICE
	case "reduced_service":
		return gtfsrt.Alert_REDUCED_SERVICE
	case "significant_delays":
		return gtfsrt.Alert_SIGNIFICANT_DELAYS
	case "detour":
		return gtfsrt.Alert_DETOUR
	case "additional_service":
		return gtfsrt.Alert_ADDITIONAL_SERVICE
	case "modified_service":
		return gtfsrt.Alert_MODIFIED_SERVICE
	case "other_effect":
		return gtfsrt.Alert_OTHER_EFFECT
	case "stop_moved":
		return gtfsrt.Alert_STOP_MOVED
	case "no_effect":
		return gtfsrt.Alert_NO_EFFECT
	case "accessibility_issue":
		return gtfsrt.Alert_ACCESSIBILITY_ISSUE
	default:
		return gtfsrt.Alert_UNKNOWN_EFFECT
	}
}

func normalizeAlerts(alerts []domainalerts.Alert) []domainalerts.Alert {
	normalized := append([]domainalerts.Alert(nil), alerts...)
	sort.SliceStable(normalized, func(i int, j int) bool {
		return normalized[i].AlertKey < normalized[j].AlertKey
	})
	return normalized
}

func entitySortKey(entity domainalerts.InformedEntity) string {
	return entity.RouteID + "|" + entity.StopID + "|" + entity.TripID + "|" + entity.StartDate + "|" + entity.StartTime
}
