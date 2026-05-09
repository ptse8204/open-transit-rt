package realtimequality

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	ObservedSchemaVersion     = "realtime-quality-observed.v1"
	PredictionSchemaVersion   = "realtime-quality-predictions.v1"
	BacktestSchemaVersion     = "realtime-quality-backtest.v1"
	DefaultStalePredictionTTL = 5 * time.Minute
	defaultMaxInputRecords    = 100000
	defaultMaxGroups          = 500
	defaultMaxOutputBytes     = 1 << 20
)

var allowedGateLabels = map[string]bool{
	"insufficient_data": true,
	"diagnostic_pass":   true,
	"diagnostic_watch":  true,
}

type ObservedDataset struct {
	SchemaVersion  string          `json:"schema_version"`
	AgencyTimezone string          `json:"agency_timezone"`
	Records        []ObservedEvent `json:"records"`
}

type ObservedEvent struct {
	AgencyID            string    `json:"agency_id"`
	FeedVersionID       string    `json:"feed_version_id"`
	RouteID             string    `json:"route_id"`
	TripID              string    `json:"trip_id"`
	StartDate           string    `json:"start_date"`
	StartTime           string    `json:"start_time"`
	StopID              string    `json:"stop_id"`
	StopSequence        int       `json:"stop_sequence"`
	EventType           string    `json:"event_type"`
	ObservedTime        time.Time `json:"observed_time"`
	ServicePatternLabel string    `json:"service_pattern_label,omitempty"`
}

type PredictionDataset struct {
	SchemaVersion string             `json:"schema_version"`
	Records       []PredictionSample `json:"records"`
}

type PredictionSample struct {
	GeneratedTime        time.Time `json:"generated_time"`
	AdapterName          string    `json:"adapter_name"`
	AgencyID             string    `json:"agency_id"`
	FeedVersionID        string    `json:"feed_version_id"`
	RouteID              string    `json:"route_id"`
	TripID               string    `json:"trip_id"`
	StartDate            string    `json:"start_date"`
	StartTime            string    `json:"start_time"`
	StopSequence         int       `json:"stop_sequence"`
	EventType            string    `json:"event_type"`
	PredictedTime        time.Time `json:"predicted_time,omitempty"`
	Confidence           *float64  `json:"confidence,omitempty"`
	ScheduleRelationship string    `json:"schedule_relationship,omitempty"`
	WithheldReason       string    `json:"withheld_reason,omitempty"`
}

type BacktestOptions struct {
	GeneratedAt        time.Time
	StalePredictionTTL time.Duration
	MinMatchedSamples  int
	PassMinCoverage    float64
	PassMaxMAESeconds  float64
	PassMaxP90Seconds  float64
	MaxInputRecords    int
	MaxGroups          int
	MaxOutputBytes     int
}

type BacktestReport struct {
	Summary  SummaryDocument  `json:"summary"`
	Metrics  MetricsDocument  `json:"metrics"`
	Manifest ManifestDocument `json:"manifest"`
}

type SummaryDocument struct {
	SchemaVersion                      string       `json:"schema_version"`
	GeneratedAt                        string       `json:"generated_at"`
	InputRecordCounts                  RecordCounts `json:"input_record_counts"`
	Overall                            MetricGroup  `json:"overall"`
	GroupCount                         int          `json:"group_count"`
	ExternalEvidenceCreated            bool         `json:"external_evidence_created"`
	ConsumerStatusesChanged            bool         `json:"consumer_statuses_changed"`
	ComplianceClaimed                  bool         `json:"compliance_claimed"`
	ProductionReadinessClaimed         bool         `json:"production_readiness_claimed"`
	HostedSaaSClaimed                  bool         `json:"hosted_saas_claimed"`
	AgencyAdoptionClaimed              bool         `json:"agency_adoption_claimed"`
	ConsumerAcceptanceClaimed          bool         `json:"consumer_acceptance_claimed"`
	VendorCompatibilityClaimed         bool         `json:"vendor_compatibility_claimed"`
	ProductionGradeETAClaimed          bool         `json:"production_grade_eta_claimed"`
	PublicAPIAdded                     bool         `json:"public_api_added"`
	DatabasePersistenceAdded           bool         `json:"database_persistence_added"`
	OperationsConsoleBacktestViewAdded bool         `json:"operations_console_backtest_view_added"`
}

type RecordCounts struct {
	ObservedRecords   int `json:"observed_records"`
	PredictionRecords int `json:"prediction_records"`
}

type MetricsDocument struct {
	SchemaVersion string        `json:"schema_version"`
	GeneratedAt   string        `json:"generated_at"`
	Groups        []MetricGroup `json:"groups"`
}

type ManifestDocument struct {
	SchemaVersion             string            `json:"schema_version"`
	GeneratedAt               string            `json:"generated_at"`
	OutputKind                string            `json:"output_kind"`
	OutputFiles               []string          `json:"output_files"`
	Inputs                    []InputManifest   `json:"inputs"`
	SafetyChecks              map[string]bool   `json:"safety_checks"`
	Boundaries                map[string]bool   `json:"boundaries"`
	AggregateOnly             bool              `json:"aggregate_only"`
	RawRowsPersisted          bool              `json:"raw_rows_persisted"`
	MaturityGateAllowedLabels []string          `json:"maturity_gate_allowed_labels"`
	Options                   map[string]string `json:"options"`
}

type InputManifest struct {
	Kind          string `json:"kind"`
	SchemaVersion string `json:"schema_version"`
	RecordCount   int    `json:"record_count"`
	SHA256        string `json:"sha256"`
}

type MetricGroup struct {
	GroupType                   string         `json:"group_type"`
	RouteID                     string         `json:"route_id,omitempty"`
	TimePeriod                  string         `json:"time_period,omitempty"`
	CoverageDenominator         int            `json:"coverage_denominator"`
	CoverageNumerator           int            `json:"coverage_numerator"`
	FutureStopCoverageNumerator int            `json:"future_stop_coverage_numerator"`
	MatchedPredictionCount      int            `json:"matched_prediction_count"`
	MissingPredictionCount      int            `json:"missing_prediction_count"`
	MissingObservationCount     int            `json:"missing_observation_count"`
	StalePredictionCount        int            `json:"stale_prediction_count"`
	WithheldByReason            map[string]int `json:"withheld_by_reason,omitempty"`
	MAEAbsoluteErrorSeconds     *float64       `json:"mae_absolute_error_seconds,omitempty"`
	MedianAbsoluteErrorSeconds  *float64       `json:"median_absolute_error_seconds,omitempty"`
	P90AbsoluteErrorSeconds     *float64       `json:"p90_absolute_error_seconds,omitempty"`
	MeanLeadTimeSeconds         *float64       `json:"mean_lead_time_seconds,omitempty"`
	MedianLeadTimeSeconds       *float64       `json:"median_lead_time_seconds,omitempty"`
	P90LeadTimeSeconds          *float64       `json:"p90_lead_time_seconds,omitempty"`
	PredictionCoverage          Rate           `json:"prediction_coverage"`
	FutureStopCoverage          Rate           `json:"future_stop_coverage"`
	MaturityGate                string         `json:"maturity_gate"`
}

type Rate struct {
	Numerator   int      `json:"numerator"`
	Denominator int      `json:"denominator"`
	Status      string   `json:"status"`
	Percent     *float64 `json:"percent,omitempty"`
}

type BacktestFiles struct {
	Output OutputTarget
	Files  map[string][]byte
	Stdout string
}

type OutputTarget struct {
	Dir   string
	Label string
	Ref   string
}

func LoadObservedDataset(path string) (ObservedDataset, []byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ObservedDataset{}, nil, fmt.Errorf("read observed events: %w", err)
	}
	var dataset ObservedDataset
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&dataset); err != nil {
		return ObservedDataset{}, raw, fmt.Errorf("decode observed events: %w", err)
	}
	if err := ValidateObservedDataset(dataset, defaultMaxInputRecords); err != nil {
		return ObservedDataset{}, raw, err
	}
	return dataset, raw, nil
}

func LoadPredictionDataset(path string) (PredictionDataset, []byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return PredictionDataset{}, nil, fmt.Errorf("read prediction samples: %w", err)
	}
	var dataset PredictionDataset
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&dataset); err != nil {
		return PredictionDataset{}, raw, fmt.Errorf("decode prediction samples: %w", err)
	}
	if err := ValidatePredictionDataset(dataset, defaultMaxInputRecords); err != nil {
		return PredictionDataset{}, raw, err
	}
	return dataset, raw, nil
}

func ValidateObservedDataset(dataset ObservedDataset, maxRecords int) error {
	if dataset.SchemaVersion != ObservedSchemaVersion {
		return fmt.Errorf("observed schema_version must be %s", ObservedSchemaVersion)
	}
	if dataset.AgencyTimezone == "" {
		return fmt.Errorf("observed agency_timezone is required")
	}
	if _, err := time.LoadLocation(dataset.AgencyTimezone); err != nil {
		return fmt.Errorf("observed agency_timezone is invalid: %w", err)
	}
	if maxRecords <= 0 {
		maxRecords = defaultMaxInputRecords
	}
	if len(dataset.Records) > maxRecords {
		return fmt.Errorf("observed record count exceeds %d", maxRecords)
	}
	for i, row := range dataset.Records {
		if err := validateObservedEvent(row); err != nil {
			return fmt.Errorf("observed record %d: %w", i, err)
		}
	}
	return nil
}

func ValidatePredictionDataset(dataset PredictionDataset, maxRecords int) error {
	if dataset.SchemaVersion != PredictionSchemaVersion {
		return fmt.Errorf("prediction schema_version must be %s", PredictionSchemaVersion)
	}
	if maxRecords <= 0 {
		maxRecords = defaultMaxInputRecords
	}
	if len(dataset.Records) > maxRecords {
		return fmt.Errorf("prediction record count exceeds %d", maxRecords)
	}
	for i, row := range dataset.Records {
		if err := validatePredictionSample(row); err != nil {
			return fmt.Errorf("prediction record %d: %w", i, err)
		}
	}
	return nil
}

func validateObservedEvent(row ObservedEvent) error {
	if err := validateSharedKeys(row.AgencyID, row.FeedVersionID, row.TripID, row.StartDate, row.StartTime, row.StopSequence, row.EventType); err != nil {
		return err
	}
	if row.RouteID == "" || row.StopID == "" {
		return fmt.Errorf("route_id and stop_id are required")
	}
	if row.ObservedTime.IsZero() {
		return fmt.Errorf("observed_time is required")
	}
	return validateSafeStrings(row.AgencyID, row.FeedVersionID, row.RouteID, row.TripID, row.StartDate, row.StartTime, row.StopID, row.ServicePatternLabel)
}

func validatePredictionSample(row PredictionSample) error {
	if err := validateSharedKeys(row.AgencyID, row.FeedVersionID, row.TripID, row.StartDate, row.StartTime, row.StopSequence, row.EventType); err != nil {
		return err
	}
	if row.GeneratedTime.IsZero() {
		return fmt.Errorf("generated_time is required")
	}
	if row.AdapterName == "" || row.RouteID == "" {
		return fmt.Errorf("adapter_name and route_id are required")
	}
	if row.PredictedTime.IsZero() && row.WithheldReason == "" {
		return fmt.Errorf("predicted_time is required unless withheld_reason is present")
	}
	if row.Confidence != nil && (*row.Confidence < 0 || *row.Confidence > 1) {
		return fmt.Errorf("confidence must be between 0 and 1")
	}
	return validateSafeStrings(row.AdapterName, row.AgencyID, row.FeedVersionID, row.RouteID, row.TripID, row.StartDate, row.StartTime, row.ScheduleRelationship, row.WithheldReason)
}

func validateSharedKeys(agencyID, feedVersionID, tripID, startDate, startTime string, stopSequence int, eventType string) error {
	if agencyID == "" || feedVersionID == "" || tripID == "" || startDate == "" || startTime == "" {
		return fmt.Errorf("agency_id, feed_version_id, trip_id, start_date, and start_time are required")
	}
	if stopSequence <= 0 {
		return fmt.Errorf("stop_sequence must be positive")
	}
	if eventType != "arrival" && eventType != "departure" {
		return fmt.Errorf("event_type must be arrival or departure")
	}
	return nil
}

func validateSafeStrings(values ...string) error {
	for _, value := range values {
		lower := strings.ToLower(value)
		for _, token := range []string{"authorization:", "cookie:", "bearer ", "postgres://", "postgresql://", "-----begin "} {
			if strings.Contains(lower, token) {
				return fmt.Errorf("input contains forbidden private value pattern")
			}
		}
	}
	return nil
}

func RunBacktest(observed ObservedDataset, predictions PredictionDataset, observedRaw []byte, predictionsRaw []byte, options BacktestOptions) (BacktestReport, error) {
	if err := ValidateObservedDataset(observed, nonzero(options.MaxInputRecords, defaultMaxInputRecords)); err != nil {
		return BacktestReport{}, err
	}
	if err := ValidatePredictionDataset(predictions, nonzero(options.MaxInputRecords, defaultMaxInputRecords)); err != nil {
		return BacktestReport{}, err
	}
	options = normalizeOptions(options)
	loc, err := time.LoadLocation(observed.AgencyTimezone)
	if err != nil {
		return BacktestReport{}, err
	}

	predByKey := map[joinKey]PredictionSample{}
	for _, sample := range predictions.Records {
		key := keyForPrediction(sample)
		if current, ok := predByKey[key]; !ok || sample.GeneratedTime.After(current.GeneratedTime) {
			predByKey[key] = sample
		}
	}

	accumulators := map[groupKey]*metricAccumulator{}
	ensure := func(key groupKey) *metricAccumulator {
		acc := accumulators[key]
		if acc == nil {
			acc = &metricAccumulator{key: key, withheld: map[string]int{}}
			accumulators[key] = acc
		}
		return acc
	}
	overall := groupKey{kind: "overall"}
	seenObserved := map[joinKey]bool{}

	for _, obs := range observed.Records {
		key := keyForObserved(obs)
		seenObserved[key] = true
		keys := groupsFor(obs.RouteID, timePeriod(obs.ObservedTime.In(loc)))
		keys = append([]groupKey{overall}, keys...)
		for _, group := range keys {
			ensure(group).coverageDenominator++
		}
		sample, ok := predByKey[key]
		if !ok {
			for _, group := range keys {
				ensure(group).missingPrediction++
			}
			continue
		}
		if sample.WithheldReason != "" {
			for _, group := range keys {
				ensure(group).withheld[sample.WithheldReason]++
			}
			continue
		}
		stale := sample.GeneratedTime.After(obs.ObservedTime) || obs.ObservedTime.Sub(sample.GeneratedTime) > options.StalePredictionTTL
		if stale {
			for _, group := range keys {
				ensure(group).stalePrediction++
			}
			continue
		}
		if sample.PredictedTime.IsZero() {
			for _, group := range keys {
				ensure(group).missingPrediction++
			}
			continue
		}
		absErr := math.Abs(sample.PredictedTime.Sub(obs.ObservedTime).Seconds())
		future := sample.PredictedTime.After(sample.GeneratedTime)
		leadTime := sample.PredictedTime.Sub(sample.GeneratedTime).Seconds()
		for _, group := range keys {
			acc := ensure(group)
			acc.coverageNumerator++
			acc.matched++
			acc.errors = append(acc.errors, absErr)
			if future {
				acc.futureStopCoverageNumerator++
				acc.leadTimes = append(acc.leadTimes, leadTime)
			}
		}
	}

	for _, sample := range predictions.Records {
		if seenObserved[keyForPrediction(sample)] {
			continue
		}
		when := sample.PredictedTime
		if when.IsZero() {
			when = sample.GeneratedTime
		}
		keys := groupsFor(sample.RouteID, timePeriod(when.In(loc)))
		keys = append([]groupKey{overall}, keys...)
		for _, group := range keys {
			acc := ensure(group)
			if sample.WithheldReason != "" {
				acc.withheld[sample.WithheldReason]++
			} else {
				acc.missingObservation++
			}
		}
	}

	groups := make([]MetricGroup, 0, len(accumulators))
	for _, acc := range accumulators {
		groups = append(groups, acc.group(options))
	}
	sort.SliceStable(groups, func(i int, j int) bool {
		return groupSortKey(groups[i]) < groupSortKey(groups[j])
	})
	if len(groups) > options.MaxGroups {
		return BacktestReport{}, fmt.Errorf("metric group count exceeds %d", options.MaxGroups)
	}
	overallMetric := MetricGroup{GroupType: "overall", MaturityGate: "insufficient_data", PredictionCoverage: rate(0, 0), FutureStopCoverage: rate(0, 0)}
	for _, group := range groups {
		if group.GroupType == "overall" {
			overallMetric = group
			break
		}
	}
	if !allowedGateLabels[overallMetric.MaturityGate] {
		return BacktestReport{}, fmt.Errorf("unexpected maturity gate label %q", overallMetric.MaturityGate)
	}

	generated := options.GeneratedAt.UTC().Format(time.RFC3339)
	report := BacktestReport{
		Summary: SummaryDocument{
			SchemaVersion:                      BacktestSchemaVersion,
			GeneratedAt:                        generated,
			InputRecordCounts:                  RecordCounts{ObservedRecords: len(observed.Records), PredictionRecords: len(predictions.Records)},
			Overall:                            overallMetric,
			GroupCount:                         len(groups),
			ExternalEvidenceCreated:            false,
			ConsumerStatusesChanged:            false,
			ComplianceClaimed:                  false,
			ProductionReadinessClaimed:         false,
			HostedSaaSClaimed:                  false,
			AgencyAdoptionClaimed:              false,
			ConsumerAcceptanceClaimed:          false,
			VendorCompatibilityClaimed:         false,
			ProductionGradeETAClaimed:          false,
			PublicAPIAdded:                     false,
			DatabasePersistenceAdded:           false,
			OperationsConsoleBacktestViewAdded: false,
		},
		Metrics: MetricsDocument{
			SchemaVersion: BacktestSchemaVersion,
			GeneratedAt:   generated,
			Groups:        groups,
		},
		Manifest: ManifestDocument{
			SchemaVersion: BacktestSchemaVersion,
			GeneratedAt:   generated,
			OutputKind:    "private_local_realtime_quality_backtest",
			OutputFiles:   ExpectedBacktestOutputFiles(),
			Inputs: []InputManifest{
				{Kind: "observed_stop_events", SchemaVersion: observed.SchemaVersion, RecordCount: len(observed.Records), SHA256: sha256Hex(observedRaw)},
				{Kind: "prediction_samples", SchemaVersion: predictions.SchemaVersion, RecordCount: len(predictions.Records), SHA256: sha256Hex(predictionsRaw)},
			},
			SafetyChecks: map[string]bool{
				"docs_evidence_output_rejected": true,
				"evidence_like_output_rejected": true,
				"symlink_ancestors_rejected":    true,
				"unsafe_traversal_rejected":     true,
				"raw_inputs_not_copied":         true,
				"private_paths_omitted":         true,
				"raw_rows_omitted":              true,
			},
			Boundaries: map[string]bool{
				"db_persistence":                   false,
				"migration_added":                  false,
				"operations_console_change":        false,
				"public_api_added":                 false,
				"consumer_tracker_changed":         false,
				"external_predictor_runtime_added": false,
			},
			AggregateOnly:             true,
			RawRowsPersisted:          false,
			MaturityGateAllowedLabels: []string{"insufficient_data", "diagnostic_pass", "diagnostic_watch"},
			Options: map[string]string{
				"stale_prediction_ttl": options.StalePredictionTTL.String(),
				"min_matched_samples":  fmt.Sprintf("%d", options.MinMatchedSamples),
				"pass_min_coverage":    fmt.Sprintf("%.2f", options.PassMinCoverage),
				"pass_max_mae_seconds": fmt.Sprintf("%.0f", options.PassMaxMAESeconds),
				"pass_max_p90_seconds": fmt.Sprintf("%.0f", options.PassMaxP90Seconds),
			},
		},
	}
	return report, nil
}

func normalizeOptions(options BacktestOptions) BacktestOptions {
	if options.GeneratedAt.IsZero() {
		options.GeneratedAt = time.Now().UTC()
	}
	if options.StalePredictionTTL <= 0 {
		options.StalePredictionTTL = DefaultStalePredictionTTL
	}
	if options.MinMatchedSamples <= 0 {
		options.MinMatchedSamples = 1
	}
	if options.PassMinCoverage <= 0 {
		options.PassMinCoverage = 70
	}
	if options.PassMaxMAESeconds <= 0 {
		options.PassMaxMAESeconds = 180
	}
	if options.PassMaxP90Seconds <= 0 {
		options.PassMaxP90Seconds = 300
	}
	options.MaxInputRecords = nonzero(options.MaxInputRecords, defaultMaxInputRecords)
	options.MaxGroups = nonzero(options.MaxGroups, defaultMaxGroups)
	options.MaxOutputBytes = nonzero(options.MaxOutputBytes, defaultMaxOutputBytes)
	return options
}

type joinKey struct {
	agencyID, feedVersionID, tripID, startDate, startTime, eventType string
	stopSequence                                                     int
}

func keyForObserved(row ObservedEvent) joinKey {
	return joinKey{agencyID: row.AgencyID, feedVersionID: row.FeedVersionID, tripID: row.TripID, startDate: row.StartDate, startTime: row.StartTime, stopSequence: row.StopSequence, eventType: row.EventType}
}

func keyForPrediction(row PredictionSample) joinKey {
	return joinKey{agencyID: row.AgencyID, feedVersionID: row.FeedVersionID, tripID: row.TripID, startDate: row.StartDate, startTime: row.StartTime, stopSequence: row.StopSequence, eventType: row.EventType}
}

type groupKey struct {
	kind       string
	routeID    string
	timePeriod string
}

func groupsFor(routeID, period string) []groupKey {
	return []groupKey{
		{kind: "route", routeID: routeID},
		{kind: "time_period", timePeriod: period},
		{kind: "route_time_period", routeID: routeID, timePeriod: period},
	}
}

func timePeriod(local time.Time) string {
	hour := local.Hour()
	switch {
	case hour < 5:
		return "overnight"
	case hour < 9:
		return "morning_peak"
	case hour < 15:
		return "midday"
	case hour < 19:
		return "afternoon_peak"
	default:
		return "evening"
	}
}

type metricAccumulator struct {
	key                         groupKey
	coverageDenominator         int
	coverageNumerator           int
	futureStopCoverageNumerator int
	matched                     int
	missingPrediction           int
	missingObservation          int
	stalePrediction             int
	withheld                    map[string]int
	errors                      []float64
	leadTimes                   []float64
}

func (a *metricAccumulator) group(options BacktestOptions) MetricGroup {
	errorsSorted := append([]float64(nil), a.errors...)
	sort.Float64s(errorsSorted)
	leadTimesSorted := append([]float64(nil), a.leadTimes...)
	sort.Float64s(leadTimesSorted)
	group := MetricGroup{
		GroupType:                   a.key.kind,
		RouteID:                     a.key.routeID,
		TimePeriod:                  a.key.timePeriod,
		CoverageDenominator:         a.coverageDenominator,
		CoverageNumerator:           a.coverageNumerator,
		FutureStopCoverageNumerator: a.futureStopCoverageNumerator,
		MatchedPredictionCount:      a.matched,
		MissingPredictionCount:      a.missingPrediction,
		MissingObservationCount:     a.missingObservation,
		StalePredictionCount:        a.stalePrediction,
		WithheldByReason:            sortedReasonMap(a.withheld),
		PredictionCoverage:          rate(a.coverageNumerator, a.coverageDenominator),
		FutureStopCoverage:          rate(a.futureStopCoverageNumerator, a.coverageDenominator),
		MaturityGate:                "insufficient_data",
	}
	if len(errorsSorted) > 0 {
		group.MAEAbsoluteErrorSeconds = floatPtr(round1(mean(errorsSorted)))
		group.MedianAbsoluteErrorSeconds = floatPtr(round1(percentileNearest(errorsSorted, 0.50)))
		group.P90AbsoluteErrorSeconds = floatPtr(round1(percentileNearest(errorsSorted, 0.90)))
	}
	if len(leadTimesSorted) > 0 {
		group.MeanLeadTimeSeconds = floatPtr(round1(mean(leadTimesSorted)))
		group.MedianLeadTimeSeconds = floatPtr(round1(percentileNearest(leadTimesSorted, 0.50)))
		group.P90LeadTimeSeconds = floatPtr(round1(percentileNearest(leadTimesSorted, 0.90)))
	}
	if a.coverageDenominator > 0 && a.matched >= options.MinMatchedSamples {
		coverage := 0.0
		if group.PredictionCoverage.Percent != nil {
			coverage = *group.PredictionCoverage.Percent
		}
		maePass := group.MAEAbsoluteErrorSeconds != nil && *group.MAEAbsoluteErrorSeconds <= options.PassMaxMAESeconds
		p90Pass := group.P90AbsoluteErrorSeconds != nil && *group.P90AbsoluteErrorSeconds <= options.PassMaxP90Seconds
		if coverage >= options.PassMinCoverage && maePass && p90Pass {
			group.MaturityGate = "diagnostic_pass"
		} else {
			group.MaturityGate = "diagnostic_watch"
		}
	}
	return group
}

func rate(numerator, denominator int) Rate {
	out := Rate{Numerator: numerator, Denominator: denominator, Status: "not_applicable"}
	if denominator <= 0 {
		return out
	}
	out.Status = "measured"
	percent := round1(float64(numerator) / float64(denominator) * 100)
	out.Percent = &percent
	return out
}

func mean(values []float64) float64 {
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func percentileNearest(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	index := int(math.Ceil(p*float64(len(sorted)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

func sortedReasonMap(in map[string]int) map[string]int {
	if len(in) == 0 {
		return nil
	}
	out := map[string]int{}
	keys := make([]string, 0, len(in))
	for key := range in {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		out[key] = in[key]
	}
	return out
}

func groupSortKey(group MetricGroup) string {
	return group.GroupType + "\x00" + group.RouteID + "\x00" + group.TimePeriod
}

func ExpectedBacktestOutputFiles() []string {
	return []string{"summary.json", "summary.md", "metrics.json", "metrics.md", "manifest.json"}
}

func ResolveBacktestOutputTarget(raw string, generatedAt time.Time, cwd string) (OutputTarget, error) {
	label := "default"
	if strings.TrimSpace(raw) == "" {
		raw = filepath.Join(".cache", "realtime-quality-backtest", generatedAt.UTC().Format("20060102T150405Z"))
	} else {
		label = "override"
	}
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return OutputTarget{}, fmt.Errorf("output path cannot be resolved")
		}
	}
	abs := raw
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(cwd, raw)
	}
	abs = filepath.Clean(abs)
	rel, err := filepath.Rel(cwd, abs)
	if err != nil {
		return OutputTarget{}, fmt.Errorf("output path cannot be resolved")
	}
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return OutputTarget{}, fmt.Errorf("output path must stay below the repository root")
	}
	relSlash := filepath.ToSlash(rel)
	lower := strings.ToLower(relSlash)
	if lower == "docs/evidence" || strings.HasPrefix(lower, "docs/evidence/") || strings.Contains(lower, "evidence") {
		return OutputTarget{}, fmt.Errorf("output path must not be evidence-like or under docs/evidence")
	}
	if lower != ".cache" && !strings.HasPrefix(lower, ".cache/") {
		return OutputTarget{}, fmt.Errorf("output path must stay under .cache")
	}
	if hasSymlinkAncestorBelow(cwd, abs) {
		return OutputTarget{}, fmt.Errorf("output path must not include symlink ancestors")
	}
	if entries, err := os.ReadDir(abs); err == nil && len(entries) > 0 {
		return OutputTarget{}, fmt.Errorf("output directory must be empty when it already exists")
	}
	return OutputTarget{Dir: abs, Label: label, Ref: redactedOutputRef(relSlash)}, nil
}

func BuildBacktestFiles(report BacktestReport, output OutputTarget, forbidden []string, maxOutputBytes int) (BacktestFiles, error) {
	maxOutputBytes = nonzero(maxOutputBytes, defaultMaxOutputBytes)
	files := map[string][]byte{}
	putJSON := func(name string, value any) error {
		raw, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			return fmt.Errorf("encode %s: %w", name, err)
		}
		raw = append(raw, '\n')
		if len(raw) > maxOutputBytes {
			return fmt.Errorf("%s exceeds output size cap", name)
		}
		files[name] = raw
		return nil
	}
	if err := putJSON("summary.json", report.Summary); err != nil {
		return BacktestFiles{}, err
	}
	if err := putJSON("metrics.json", report.Metrics); err != nil {
		return BacktestFiles{}, err
	}
	if err := putJSON("manifest.json", report.Manifest); err != nil {
		return BacktestFiles{}, err
	}
	files["summary.md"] = []byte(renderBacktestSummaryMarkdown(report.Summary))
	files["metrics.md"] = []byte(renderBacktestMetricsMarkdown(report.Metrics))
	for name, raw := range files {
		if len(raw) > maxOutputBytes {
			return BacktestFiles{}, fmt.Errorf("%s exceeds output size cap", name)
		}
		if hits := ScanBacktestRedaction(string(raw), forbidden); len(hits) > 0 {
			return BacktestFiles{}, fmt.Errorf("redaction scan failed for %s", name)
		}
	}
	stdout := fmt.Sprintf("{\"mode\":\"realtime_quality_backtest\",\"output_ref\":%q,\"observed_records\":%d,\"prediction_records\":%d,\"maturity_gate\":%q}\n",
		output.Ref,
		report.Summary.InputRecordCounts.ObservedRecords,
		report.Summary.InputRecordCounts.PredictionRecords,
		report.Summary.Overall.MaturityGate,
	)
	if hits := ScanBacktestRedaction(stdout, forbidden); len(hits) > 0 {
		return BacktestFiles{}, fmt.Errorf("redaction scan failed for terminal output")
	}
	return BacktestFiles{Output: output, Files: files, Stdout: stdout}, nil
}

func WriteBacktestFiles(files BacktestFiles) error {
	expected := ExpectedBacktestOutputFiles()
	if len(files.Files) != len(expected) {
		return fmt.Errorf("backtest output file set is incomplete")
	}
	for _, name := range expected {
		if _, ok := files.Files[name]; !ok {
			return fmt.Errorf("backtest output file set is missing %s", name)
		}
	}
	if err := os.MkdirAll(files.Output.Dir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	for _, name := range expected {
		path := filepath.Join(files.Output.Dir, name)
		if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("output file path must not be a symlink")
		}
		if err := os.WriteFile(path, files.Files[name], 0o644); err != nil {
			return fmt.Errorf("write output file: %w", err)
		}
	}
	entries, err := os.ReadDir(files.Output.Dir)
	if err != nil {
		return fmt.Errorf("read output directory: %w", err)
	}
	if len(entries) != len(expected) {
		return fmt.Errorf("output directory contains unexpected files")
	}
	return nil
}

func ForbiddenBacktestOutputValues(paths ...string) []string {
	seen := map[string]bool{}
	for _, path := range paths {
		addForbiddenValue(seen, path)
		if abs, err := filepath.Abs(path); err == nil {
			addForbiddenValue(seen, abs)
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		addForbiddenValue(seen, cwd)
	}
	values := make([]string, 0, len(seen))
	for value := range seen {
		values = append(values, value)
	}
	sort.Strings(values)
	return values
}

func addForbiddenValue(seen map[string]bool, value string) {
	value = strings.TrimSpace(value)
	if len(value) >= 4 {
		seen[value] = true
	}
}

func ScanBacktestRedaction(content string, forbidden []string) []string {
	var hits []string
	lower := strings.ToLower(content)
	for _, value := range forbidden {
		if value != "" && strings.Contains(content, value) {
			hits = append(hits, "forbidden_value")
		}
	}
	for _, pattern := range []string{
		`authorization\s*:`,
		`\bbearer\s+[a-z0-9._~+/=-]+`,
		`\bcookie\s*[:=]`,
		`postgres(?:ql)?://`,
		`-----begin [^-]*private key-----`,
		`/users/[^ \n\t"]+`,
	} {
		if regexp.MustCompile(pattern).MatchString(lower) {
			hits = append(hits, "forbidden_pattern")
		}
	}
	for _, forbiddenGate := range []string{`"ready"`, `"production_ready"`, `"compliant"`, `"accepted"`} {
		if strings.Contains(lower, forbiddenGate) {
			hits = append(hits, "forbidden_gate_label")
		}
	}
	return hits
}

func renderBacktestSummaryMarkdown(summary SummaryDocument) string {
	var b strings.Builder
	b.WriteString("# Realtime Quality Backtest Summary\n\n")
	b.WriteString("Private local diagnostics only. This is not an evidence packet, not a consumer submission artifact, not a public API, and not production-grade ETA proof.\n\n")
	b.WriteString(fmt.Sprintf("- generated_at: `%s`\n", summary.GeneratedAt))
	b.WriteString(fmt.Sprintf("- observed_records: `%d`\n", summary.InputRecordCounts.ObservedRecords))
	b.WriteString(fmt.Sprintf("- prediction_records: `%d`\n", summary.InputRecordCounts.PredictionRecords))
	b.WriteString(fmt.Sprintf("- maturity_gate: `%s`\n", summary.Overall.MaturityGate))
	b.WriteString(fmt.Sprintf("- prediction_coverage: `%s`\n", formatRate(summary.Overall.PredictionCoverage)))
	if summary.Overall.MAEAbsoluteErrorSeconds != nil {
		b.WriteString(fmt.Sprintf("- mae_absolute_error_seconds: `%.1f`\n", *summary.Overall.MAEAbsoluteErrorSeconds))
	}
	if summary.Overall.MeanLeadTimeSeconds != nil {
		b.WriteString(fmt.Sprintf("- mean_lead_time_seconds: `%.1f`\n", *summary.Overall.MeanLeadTimeSeconds))
	}
	b.WriteString("- external_evidence_created: `false`\n")
	b.WriteString("- consumer_statuses_changed: `false`\n")
	return b.String()
}

func renderBacktestMetricsMarkdown(metrics MetricsDocument) string {
	var b strings.Builder
	b.WriteString("# Realtime Quality Backtest Metrics\n\n")
	b.WriteString("| group | coverage | future_stop_coverage | mae_s | median_s | p90_s | mean_lead_s | median_lead_s | p90_lead_s | stale | missing_prediction | missing_observation | gate |\n")
	b.WriteString("| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |\n")
	for _, group := range metrics.Groups {
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %d | %d | %d | %s |\n",
			markdownGroupLabel(group),
			formatRate(group.PredictionCoverage),
			formatRate(group.FutureStopCoverage),
			formatOptionalFloat(group.MAEAbsoluteErrorSeconds),
			formatOptionalFloat(group.MedianAbsoluteErrorSeconds),
			formatOptionalFloat(group.P90AbsoluteErrorSeconds),
			formatOptionalFloat(group.MeanLeadTimeSeconds),
			formatOptionalFloat(group.MedianLeadTimeSeconds),
			formatOptionalFloat(group.P90LeadTimeSeconds),
			group.StalePredictionCount,
			group.MissingPredictionCount,
			group.MissingObservationCount,
			group.MaturityGate,
		))
	}
	return b.String()
}

func markdownGroupLabel(group MetricGroup) string {
	switch group.GroupType {
	case "route":
		return "route:" + group.RouteID
	case "time_period":
		return "time_period:" + group.TimePeriod
	case "route_time_period":
		return "route:" + group.RouteID + " / " + group.TimePeriod
	default:
		return group.GroupType
	}
}

func formatRate(rate Rate) string {
	if rate.Percent == nil {
		return "n/a"
	}
	return fmt.Sprintf("%.1f%%", *rate.Percent)
}

func formatOptionalFloat(value *float64) string {
	if value == nil {
		return "n/a"
	}
	return fmt.Sprintf("%.1f", *value)
}

func sha256Hex(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func hasSymlinkAncestorBelow(root string, path string) bool {
	root = filepath.Clean(root)
	path = filepath.Clean(path)
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false
	}
	current := root
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return true
		}
	}
	return false
}

func redactedOutputRef(rel string) string {
	rel = filepath.ToSlash(filepath.Clean(rel))
	if strings.HasPrefix(rel, ".cache/") {
		return rel
	}
	return "path:" + shortHash(rel)
}

func shortHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])[:12]
}

func nonzero(value int, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

func round1(value float64) float64 {
	return math.Round(value*10) / 10
}

func floatPtr(value float64) *float64 {
	return &value
}

func IsOutputPathError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "output path") || errors.Is(err, os.ErrPermission)
}
