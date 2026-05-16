package gtfsrtconformance

import (
	"fmt"
	"sort"
	"strings"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

const (
	SchemaVersion = "open-transit-rt.gtfsrt_conformance.v1"

	FeedVehiclePositions = "vehicle_positions"
	FeedTripUpdates      = "trip_updates"
	FeedAlerts           = "alerts"

	StatusOK          = "ok"
	StatusNeedsReview = "needs_review"
	StatusFailed      = "failed"

	SeverityInfo        = "info"
	SeverityNeedsReview = "needs_review"
	SeverityFailed      = "failed"
)

var feedOrder = map[string]int{
	FeedVehiclePositions: 0,
	FeedTripUpdates:      1,
	FeedAlerts:           2,
}

type Options struct {
	GeneratedAt   time.Time
	MaxFutureSkew time.Duration
}

type Report struct {
	SchemaVersion      string       `json:"schema_version"`
	GeneratedAt        time.Time    `json:"generated_at"`
	Boundary           string       `json:"boundary"`
	SyntheticLocalOnly bool         `json:"synthetic_local_only"`
	OverallStatus      string       `json:"overall_status"`
	Feeds              []FeedReport `json:"feeds"`
	ClaimFlags         ClaimFlags   `json:"claim_flags"`
}

type ClaimFlags struct {
	ExternalEvidenceCreated       bool `json:"external_evidence_created"`
	ConsumerStatusesChanged       bool `json:"consumer_statuses_changed"`
	ComplianceClaimed             bool `json:"compliance_claimed"`
	ConsumerAcceptanceClaimed     bool `json:"consumer_acceptance_claimed"`
	ProductionReadinessClaimed    bool `json:"production_readiness_claimed"`
	SLAClaimed                    bool `json:"sla_claimed"`
	UptimeGuaranteeClaimed        bool `json:"uptime_guarantee_claimed"`
	VendorCompatibilityClaimed    bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed  bool `json:"hardware_certification_claimed"`
	ProductionGradeETAClaimed     bool `json:"production_grade_eta_claimed"`
	RealWorldETAAccuracyClaimed   bool `json:"real_world_eta_accuracy_claimed"`
	StableReleaseReadinessClaimed bool `json:"stable_release_readiness_claimed"`
}

type FeedInput struct {
	FeedType string
	Source   string
	Payload  []byte
}

type FeedReport struct {
	FeedType     string  `json:"feed_type"`
	Source       string  `json:"source"`
	Status       string  `json:"status"`
	EntityCount  int     `json:"entity_count"`
	HeaderTime   *string `json:"header_time,omitempty"`
	DoesNotProve string  `json:"does_not_prove"`
	Checks       []Check `json:"checks"`
}

type Check struct {
	ID       string `json:"id"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func BuildReport(inputs []FeedInput, opts Options) Report {
	generatedAt := opts.GeneratedAt.UTC()
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC().Truncate(time.Second)
	}
	reports := make([]FeedReport, 0, len(inputs))
	for _, input := range inputs {
		opts.GeneratedAt = generatedAt
		reports = append(reports, CheckPayload(input.FeedType, input.Source, input.Payload, opts))
	}
	sort.SliceStable(reports, func(i, j int) bool {
		left, lok := feedOrder[reports[i].FeedType]
		right, rok := feedOrder[reports[j].FeedType]
		if lok && rok && left != right {
			return left < right
		}
		return reports[i].FeedType < reports[j].FeedType
	})
	return Report{
		SchemaVersion:      SchemaVersion,
		GeneratedAt:        generatedAt,
		Boundary:           "Offline local GTFS-RT interoperability diagnostics only. This report does not contact consumers, write evidence, change consumer status, or prove compliance, consumer acceptance, production readiness, SLA/uptime, vendor compatibility, hardware certification, production-grade ETA quality, or real-world ETA accuracy.",
		SyntheticLocalOnly: true,
		OverallStatus:      overallStatus(reports),
		Feeds:              reports,
		ClaimFlags:         ClaimFlags{},
	}
}

func CheckPayload(feedType string, source string, payload []byte, opts Options) FeedReport {
	report := FeedReport{
		FeedType:     strings.TrimSpace(feedType),
		Source:       strings.TrimSpace(source),
		Status:       StatusOK,
		DoesNotProve: doesNotProve(feedType),
	}
	if !supportedFeed(feedType) {
		report.add(SeverityFailed, "feed_type_supported", fmt.Sprintf("unsupported GTFS-RT feed type %q", feedType))
		return report
	}
	if len(payload) == 0 {
		report.add(SeverityFailed, "payload_non_empty", "protobuf payload is empty")
		return report
	}
	var message gtfsrt.FeedMessage
	if err := proto.Unmarshal(payload, &message); err != nil {
		report.add(SeverityFailed, "payload_parseable", "protobuf payload could not be parsed: "+err.Error())
		return report
	}
	if err := proto.CheckInitialized(&message); err != nil {
		report.add(SeverityFailed, "payload_initialized", "protobuf payload has missing required fields: "+err.Error())
		return report
	}
	report.add(SeverityInfo, "payload_parseable", "protobuf payload parses as a GTFS-RT FeedMessage")
	checkHeader(&message, &report, opts)
	report.EntityCount = len(message.Entity)
	checkEntityIDs(message.Entity, &report)
	switch feedType {
	case FeedVehiclePositions:
		checkVehiclePositions(message.Entity, &report)
	case FeedTripUpdates:
		checkTripUpdates(message.Entity, &report)
	case FeedAlerts:
		checkAlerts(message.Entity, &report)
	}
	if report.EntityCount == 0 {
		report.add(SeverityInfo, "empty_full_dataset", "empty full-dataset GTFS-RT feed is parseable; operators should still review freshness and feed health context")
	}
	return report
}

func checkHeader(message *gtfsrt.FeedMessage, report *FeedReport, opts Options) {
	header := message.GetHeader()
	if header == nil {
		report.add(SeverityFailed, "header_present", "feed header is missing")
		return
	}
	if header.GetGtfsRealtimeVersion() == "2.0" {
		report.add(SeverityInfo, "header_version", "gtfs_realtime_version is 2.0")
	} else if strings.TrimSpace(header.GetGtfsRealtimeVersion()) == "" {
		report.add(SeverityFailed, "header_version", "gtfs_realtime_version is missing")
	} else {
		report.add(SeverityNeedsReview, "header_version", fmt.Sprintf("gtfs_realtime_version is %q; Open Transit RT generated feeds use 2.0 for local interoperability review", header.GetGtfsRealtimeVersion()))
	}
	if header.Incrementality == nil {
		report.add(SeverityNeedsReview, "header_incrementality", "incrementality is omitted; full-dataset publication is the expected local review profile")
	} else if header.GetIncrementality() == gtfsrt.FeedHeader_FULL_DATASET {
		report.add(SeverityInfo, "header_incrementality", "incrementality is FULL_DATASET")
	} else {
		report.add(SeverityNeedsReview, "header_incrementality", "incrementality is not FULL_DATASET; review consumer expectations before publishing differential feeds")
	}
	if header.Timestamp == nil {
		report.add(SeverityNeedsReview, "header_timestamp", "header timestamp is missing")
		return
	}
	headerTime := time.Unix(int64(header.GetTimestamp()), 0).UTC().Format(time.RFC3339)
	report.HeaderTime = &headerTime
	maxFutureSkew := opts.MaxFutureSkew
	if maxFutureSkew <= 0 {
		maxFutureSkew = 2 * time.Minute
	}
	now := opts.GeneratedAt.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if time.Unix(int64(header.GetTimestamp()), 0).UTC().After(now.Add(maxFutureSkew)) {
		report.add(SeverityNeedsReview, "header_timestamp", fmt.Sprintf("header timestamp %s is more than %s in the future of generated_at %s", headerTime, maxFutureSkew, now.Format(time.RFC3339)))
		return
	}
	report.add(SeverityInfo, "header_timestamp", "header timestamp is present")
}

func checkEntityIDs(entities []*gtfsrt.FeedEntity, report *FeedReport) {
	seen := map[string]bool{}
	for idx, entity := range entities {
		id := strings.TrimSpace(entity.GetId())
		if id == "" {
			report.add(SeverityFailed, "entity_id_present", fmt.Sprintf("entity at index %d has an empty id", idx))
			continue
		}
		if seen[id] {
			report.add(SeverityFailed, "entity_id_unique", fmt.Sprintf("entity id %q appears more than once", id))
			continue
		}
		seen[id] = true
	}
	if len(entities) > 0 && len(seen) == len(entities) {
		report.add(SeverityInfo, "entity_id_unique", "all entity ids are non-empty and unique")
	}
}

func checkVehiclePositions(entities []*gtfsrt.FeedEntity, report *FeedReport) {
	for _, entity := range entities {
		vehicle := entity.GetVehicle()
		if vehicle == nil {
			report.add(SeverityFailed, "vehicle_position_payload", fmt.Sprintf("entity %q does not contain a VehiclePosition payload", entity.GetId()))
			continue
		}
		if strings.TrimSpace(vehicle.GetVehicle().GetId()) == "" {
			report.add(SeverityNeedsReview, "vehicle_id_present", fmt.Sprintf("entity %q has no vehicle.id", entity.GetId()))
		}
		position := vehicle.GetPosition()
		if position == nil {
			report.add(SeverityFailed, "vehicle_position_coordinates", fmt.Sprintf("entity %q has no position", entity.GetId()))
		} else if !validLatLon(float64(position.GetLatitude()), float64(position.GetLongitude())) {
			report.add(SeverityFailed, "vehicle_position_coordinates", fmt.Sprintf("entity %q has invalid latitude/longitude", entity.GetId()))
		}
		if vehicle.Timestamp == nil {
			report.add(SeverityNeedsReview, "vehicle_position_timestamp", fmt.Sprintf("entity %q has no vehicle timestamp", entity.GetId()))
		}
		trip := vehicle.GetTrip()
		if trip != nil && strings.TrimSpace(trip.GetTripId()) != "" && strings.TrimSpace(trip.GetStartDate()) == "" && strings.TrimSpace(trip.GetStartTime()) == "" {
			report.add(SeverityNeedsReview, "vehicle_trip_identity", fmt.Sprintf("entity %q has a trip_id without start_date or start_time; repeated trips may be ambiguous", entity.GetId()))
		}
	}
	if len(entities) > 0 {
		report.add(SeverityInfo, "vehicle_position_payload", "Vehicle Positions payloads were inspected for vehicle identity, coordinates, timestamps, and conservative trip identity")
	}
}

func checkTripUpdates(entities []*gtfsrt.FeedEntity, report *FeedReport) {
	for _, entity := range entities {
		update := entity.GetTripUpdate()
		if update == nil {
			report.add(SeverityFailed, "trip_update_payload", fmt.Sprintf("entity %q does not contain a TripUpdate payload", entity.GetId()))
			continue
		}
		trip := update.GetTrip()
		if trip == nil {
			report.add(SeverityFailed, "trip_update_trip_identity", fmt.Sprintf("entity %q has no trip descriptor", entity.GetId()))
			continue
		}
		if strings.TrimSpace(trip.GetTripId()) == "" {
			report.add(SeverityFailed, "trip_update_trip_identity", fmt.Sprintf("entity %q has no trip_id", entity.GetId()))
		}
		if strings.TrimSpace(trip.GetStartDate()) == "" {
			report.add(SeverityNeedsReview, "trip_update_start_date", fmt.Sprintf("entity %q has no start_date; repeated or after-midnight trips may be ambiguous", entity.GetId()))
		}
		if len(update.GetStopTimeUpdate()) == 0 && trip.GetScheduleRelationship() != gtfsrt.TripDescriptor_CANCELED {
			report.add(SeverityNeedsReview, "trip_update_stop_times", fmt.Sprintf("entity %q has no stop_time_update rows; review consumer behavior for empty non-canceled updates", entity.GetId()))
		}
	}
	if len(entities) > 0 {
		report.add(SeverityInfo, "trip_update_payload", "Trip Updates payloads were inspected for trip identity, start_date, and stop-time coverage")
	}
}

func checkAlerts(entities []*gtfsrt.FeedEntity, report *FeedReport) {
	for _, entity := range entities {
		alert := entity.GetAlert()
		if alert == nil {
			report.add(SeverityFailed, "alert_payload", fmt.Sprintf("entity %q does not contain an Alert payload", entity.GetId()))
			continue
		}
		if !translatedTextPresent(alert.GetHeaderText()) {
			report.add(SeverityFailed, "alert_header_text", fmt.Sprintf("entity %q has no header_text translation", entity.GetId()))
		}
		if len(alert.GetInformedEntity()) == 0 {
			report.add(SeverityNeedsReview, "alert_informed_entity", fmt.Sprintf("entity %q has no informed_entity rows", entity.GetId()))
		}
		if len(alert.GetActivePeriod()) == 0 {
			report.add(SeverityNeedsReview, "alert_active_period", fmt.Sprintf("entity %q has no active_period; review consumer interpretation", entity.GetId()))
		}
	}
	if len(entities) > 0 {
		report.add(SeverityInfo, "alert_payload", "Alerts payloads were inspected for header text, informed entities, and active periods")
	}
}

func (r *FeedReport) add(severity string, id string, message string) {
	r.Checks = append(r.Checks, Check{ID: id, Severity: severity, Message: message})
	switch severity {
	case SeverityFailed:
		r.Status = StatusFailed
	case SeverityNeedsReview:
		if r.Status != StatusFailed {
			r.Status = StatusNeedsReview
		}
	}
}

func overallStatus(reports []FeedReport) string {
	status := StatusOK
	for _, report := range reports {
		if report.Status == StatusFailed {
			return StatusFailed
		}
		if report.Status == StatusNeedsReview {
			status = StatusNeedsReview
		}
	}
	return status
}

func supportedFeed(feedType string) bool {
	_, ok := feedOrder[strings.TrimSpace(feedType)]
	return ok
}

func validLatLon(lat float64, lon float64) bool {
	return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

func translatedTextPresent(text *gtfsrt.TranslatedString) bool {
	if text == nil {
		return false
	}
	for _, translation := range text.GetTranslation() {
		if strings.TrimSpace(translation.GetText()) != "" {
			return true
		}
	}
	return false
}

func doesNotProve(feedType string) string {
	switch feedType {
	case FeedVehiclePositions:
		return "Does not prove production AVL reliability, vendor compatibility, hardware certification, consumer display, compliance, SLA, uptime, or production readiness."
	case FeedTripUpdates:
		return "Does not prove production-grade ETA quality, real-world ETA accuracy, consumer display, public launch, compliance, SLA, or production readiness."
	case FeedAlerts:
		return "Does not prove agency approval, consumer display, public launch, compliance, SLA, uptime, or production readiness."
	default:
		return "Does not prove compliance, consumer acceptance, production readiness, SLA, uptime, vendor compatibility, hardware certification, or ETA quality."
	}
}
