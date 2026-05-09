package prediction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

const externalHTTPPath = "/v1/predict/trip-updates"

var externalHTTPTokenEnvPattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

type ExternalHTTPConfig struct {
	URL              string
	AllowedHosts     []string
	Timeout          time.Duration
	MaxRequestBytes  int64
	MaxResponseBytes int64
	TokenEnv         string
	TokenValue       string
	Client           *http.Client
}

type ExternalHTTPAdapter struct {
	endpoint         string
	client           *http.Client
	maxRequestBytes  int64
	maxResponseBytes int64
	tokenValue       string
}

func NewExternalHTTPAdapter(config ExternalHTTPConfig) (*ExternalHTTPAdapter, error) {
	validated, err := config.validated()
	if err != nil {
		return nil, err
	}
	client := validated.Client
	if client == nil {
		client = externalHTTPClient(validated.Timeout)
	}
	return &ExternalHTTPAdapter{
		endpoint:         validated.URL,
		client:           client,
		maxRequestBytes:  validated.MaxRequestBytes,
		maxResponseBytes: validated.MaxResponseBytes,
		tokenValue:       validated.TokenValue,
	}, nil
}

func (a *ExternalHTTPAdapter) Name() string {
	return AdapterNameExternalHTTP
}

func (a *ExternalHTTPAdapter) PredictTripUpdates(ctx context.Context, request Request) (Result, error) {
	start := time.Now()
	payload, err := json.Marshal(sanitizeExternalHTTPRequest(request))
	if err != nil {
		return externalHTTPFailure("encode_request", start, 0, 0, 0), nil
	}
	if int64(len(payload)) > a.maxRequestBytes {
		return externalHTTPFailure("request_too_large", start, len(payload), 0, 0), nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint, bytes.NewReader(payload))
	if err != nil {
		return externalHTTPFailure("build_request", start, len(payload), 0, 0), nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if a.tokenValue != "" {
		req.Header.Set("Authorization", "Bearer "+a.tokenValue)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return externalHTTPFailure("request_failed", start, len(payload), 0, 0), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
		return externalHTTPFailure("http_status", start, len(payload), 0, resp.StatusCode), nil
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, a.maxResponseBytes+1))
	if err != nil {
		return externalHTTPFailure("read_response", start, len(payload), 0, resp.StatusCode), nil
	}
	if int64(len(body)) > a.maxResponseBytes {
		return externalHTTPFailure("response_too_large", start, len(payload), len(body), resp.StatusCode), nil
	}

	var decoded externalHTTPResponseDTO
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&decoded); err != nil {
		return externalHTTPFailure("malformed_json", start, len(payload), len(body), resp.StatusCode), nil
	}
	result, failure := mapExternalHTTPResponse(request, decoded, start, len(payload), len(body), resp.StatusCode)
	if failure != "" {
		return externalHTTPFailure(failure, start, len(payload), len(body), resp.StatusCode), nil
	}
	return result, nil
}

func (c ExternalHTTPConfig) validated() (ExternalHTTPConfig, error) {
	c.URL = strings.TrimSpace(c.URL)
	if c.URL == "" {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_URL is required")
	}
	if len(c.AllowedHosts) == 0 {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_ALLOWED_HOSTS is required")
	}
	if c.Timeout <= 0 || c.Timeout > 30*time.Second {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_TIMEOUT_SECONDS must be between 1 and 30")
	}
	if c.MaxRequestBytes <= 0 || c.MaxRequestBytes > 5*1024*1024 {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_MAX_REQUEST_BYTES must be between 1 and 5242880")
	}
	if c.MaxResponseBytes <= 0 || c.MaxResponseBytes > 5*1024*1024 {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_MAX_RESPONSE_BYTES must be between 1 and 5242880")
	}
	if c.TokenEnv != "" {
		if err := validateExternalHTTPTokenEnvName(c.TokenEnv); err != nil {
			return ExternalHTTPConfig{}, err
		}
		if strings.TrimSpace(c.TokenValue) == "" {
			return ExternalHTTPConfig{}, fmt.Errorf("referenced TRIP_UPDATES_EXTERNAL_HTTP_TOKEN_ENV is not set")
		}
	}
	parsed, err := url.Parse(c.URL)
	if err != nil {
		return ExternalHTTPConfig{}, fmt.Errorf("parse TRIP_UPDATES_EXTERNAL_HTTP_URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_URL must be absolute")
	}
	if parsed.User != nil {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_URL must not include userinfo")
	}
	if parsed.RawQuery != "" {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_URL must not include query")
	}
	if parsed.Fragment != "" {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_URL must not include fragment")
	}
	if parsed.Path != externalHTTPPath {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_URL path must be %s", externalHTTPPath)
	}
	allowed := false
	for _, host := range c.AllowedHosts {
		if strings.TrimSpace(host) == parsed.Host {
			allowed = true
			break
		}
	}
	if !allowed {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_URL host is not allowlisted")
	}
	if parsed.Scheme != "https" && !(parsed.Scheme == "http" && isLoopbackHost(parsed.Hostname())) {
		return ExternalHTTPConfig{}, fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_URL must use https except loopback test stubs")
	}
	c.URL = parsed.String()
	return c, nil
}

func validateExternalHTTPTokenEnvName(name string) error {
	if !externalHTTPTokenEnvPattern.MatchString(name) {
		return fmt.Errorf("TRIP_UPDATES_EXTERNAL_HTTP_TOKEN_ENV must be an uppercase environment variable name")
	}
	return nil
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

type externalHTTPRequestDTO struct {
	SchemaVersion       string                  `json:"schema_version"`
	AgencyID            string                  `json:"agency_id"`
	FeedVersionID       string                  `json:"feed_version_id"`
	GeneratedAt         time.Time               `json:"generated_at"`
	VehiclePositionsURL string                  `json:"vehicle_positions_url"`
	Telemetry           []externalTelemetryDTO  `json:"telemetry"`
	Assignments         []externalAssignmentDTO `json:"assignments"`
}

type externalTelemetryDTO struct {
	VehicleID string    `json:"vehicle_id"`
	Timestamp time.Time `json:"timestamp"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	Bearing   *float64  `json:"bearing,omitempty"`
	SpeedMPS  *float64  `json:"speed_mps,omitempty"`
	AccuracyM *float64  `json:"accuracy_m,omitempty"`
	TripHint  string    `json:"trip_hint,omitempty"`
}

type externalAssignmentDTO struct {
	VehicleID            string   `json:"vehicle_id"`
	FeedVersionID        string   `json:"feed_version_id,omitempty"`
	State                string   `json:"state"`
	ServiceDate          string   `json:"service_date,omitempty"`
	RouteID              string   `json:"route_id,omitempty"`
	TripID               string   `json:"trip_id,omitempty"`
	BlockID              string   `json:"block_id,omitempty"`
	StartDate            string   `json:"start_date,omitempty"`
	StartTime            string   `json:"start_time,omitempty"`
	CurrentStopSequence  int      `json:"current_stop_sequence,omitempty"`
	ShapeDistTraveled    float64  `json:"shape_dist_traveled,omitempty"`
	Confidence           float64  `json:"confidence"`
	AssignmentSource     string   `json:"assignment_source,omitempty"`
	ReasonCodes          []string `json:"reason_codes,omitempty"`
	DegradedState        string   `json:"degraded_state,omitempty"`
	ManualOverrideActive bool     `json:"manual_override_active"`
}

type externalHTTPResponseDTO struct {
	SchemaVersion string                  `json:"schema_version"`
	AgencyID      string                  `json:"agency_id"`
	FeedVersionID string                  `json:"feed_version_id"`
	GeneratedAt   time.Time               `json:"generated_at"`
	TripUpdates   []externalTripUpdateDTO `json:"trip_updates"`
}

type externalTripUpdateDTO struct {
	AgencyID             string                  `json:"agency_id"`
	FeedVersionID        string                  `json:"feed_version_id"`
	EntityID             string                  `json:"entity_id,omitempty"`
	VehicleID            string                  `json:"vehicle_id,omitempty"`
	TripID               string                  `json:"trip_id"`
	RouteID              string                  `json:"route_id,omitempty"`
	StartDate            string                  `json:"start_date"`
	StartTime            string                  `json:"start_time"`
	ScheduleRelationship ScheduleRelationship    `json:"schedule_relationship,omitempty"`
	StopTimeUpdates      []externalStopUpdateDTO `json:"stop_time_updates"`
	Confidence           *float64                `json:"confidence"`
}

type externalStopUpdateDTO struct {
	StopID                string               `json:"stop_id,omitempty"`
	StopSequence          int                  `json:"stop_sequence"`
	ArrivalTime           *time.Time           `json:"arrival_time,omitempty"`
	DepartureTime         *time.Time           `json:"departure_time,omitempty"`
	ArrivalDelaySeconds   *int32               `json:"arrival_delay_seconds,omitempty"`
	DepartureDelaySeconds *int32               `json:"departure_delay_seconds,omitempty"`
	ScheduleRelationship  ScheduleRelationship `json:"schedule_relationship,omitempty"`
}

func sanitizeExternalHTTPRequest(request Request) externalHTTPRequestDTO {
	dto := externalHTTPRequestDTO{
		SchemaVersion:       "external_trip_updates.v1",
		AgencyID:            request.AgencyID,
		FeedVersionID:       request.ActiveFeedVersion.ID,
		GeneratedAt:         request.GeneratedAt,
		VehiclePositionsURL: request.VehiclePositionsURL,
		Telemetry:           make([]externalTelemetryDTO, 0, len(request.Telemetry)),
		Assignments:         make([]externalAssignmentDTO, 0, len(request.Assignments)),
	}
	for _, event := range request.Telemetry {
		dto.Telemetry = append(dto.Telemetry, sanitizedTelemetry(event))
	}
	sort.SliceStable(dto.Telemetry, func(i, j int) bool {
		return dto.Telemetry[i].VehicleID < dto.Telemetry[j].VehicleID
	})
	for _, vehicleID := range sortedAssignmentKeys(request.Assignments) {
		dto.Assignments = append(dto.Assignments, sanitizedAssignment(request.Assignments[vehicleID]))
	}
	return dto
}

func sanitizedTelemetry(event telemetry.StoredEvent) externalTelemetryDTO {
	dto := externalTelemetryDTO{
		VehicleID: event.VehicleID,
		Timestamp: event.Timestamp,
		Lat:       event.Lat,
		Lon:       event.Lon,
		TripHint:  event.TripHint,
	}
	if event.Bearing != 0 {
		dto.Bearing = &event.Bearing
	}
	if event.SpeedMPS != 0 {
		dto.SpeedMPS = &event.SpeedMPS
	}
	if event.AccuracyM != 0 {
		dto.AccuracyM = &event.AccuracyM
	}
	return dto
}

func sanitizedAssignment(assignment state.Assignment) externalAssignmentDTO {
	reasons := append([]string(nil), assignment.ReasonCodes...)
	sort.Strings(reasons)
	return externalAssignmentDTO{
		VehicleID:            assignment.VehicleID,
		FeedVersionID:        assignment.FeedVersionID,
		State:                string(assignment.State),
		ServiceDate:          assignment.ServiceDate,
		RouteID:              assignment.RouteID,
		TripID:               assignment.TripID,
		BlockID:              assignment.BlockID,
		StartDate:            assignment.StartDate,
		StartTime:            assignment.StartTime,
		CurrentStopSequence:  assignment.CurrentStopSequence,
		ShapeDistTraveled:    assignment.ShapeDistTraveled,
		Confidence:           assignment.Confidence,
		AssignmentSource:     string(assignment.AssignmentSource),
		ReasonCodes:          reasons,
		DegradedState:        string(assignment.DegradedState),
		ManualOverrideActive: assignment.AssignmentSource == state.AssignmentSourceManualOverride || hasAssignmentReason(assignment, state.ReasonManualOverrideActive),
	}
}

func sortedAssignmentKeys(assignments map[string]state.Assignment) []string {
	keys := make([]string, 0, len(assignments))
	for key := range assignments {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func mapExternalHTTPResponse(request Request, decoded externalHTTPResponseDTO, start time.Time, requestBytes int, responseBytes int, statusCode int) (Result, string) {
	if decoded.SchemaVersion != "external_trip_updates.v1" {
		return Result{}, "unsupported_schema"
	}
	if decoded.AgencyID == "" || decoded.FeedVersionID == "" {
		return Result{}, "missing_scope"
	}
	if decoded.AgencyID != request.AgencyID || decoded.FeedVersionID != request.ActiveFeedVersion.ID {
		return Result{}, "wrong_scope"
	}
	if decoded.GeneratedAt.IsZero() || decoded.GeneratedAt.Before(request.GeneratedAt) {
		return Result{}, "stale_response_timestamp"
	}
	updates := make([]TripUpdate, 0, len(decoded.TripUpdates))
	for _, update := range decoded.TripUpdates {
		if update.AgencyID == "" || update.FeedVersionID == "" {
			return Result{}, "missing_trip_update_scope"
		}
		if update.AgencyID != request.AgencyID || update.FeedVersionID != request.ActiveFeedVersion.ID {
			return Result{}, "wrong_trip_update_scope"
		}
		if update.Confidence == nil {
			return Result{}, "missing_confidence"
		}
		if update.StartTime == "" {
			return Result{}, "missing_start_time"
		}
		updates = append(updates, mapExternalTripUpdate(update))
	}
	reason := ReasonNoEligiblePredictions
	if len(updates) > 0 {
		reason = ReasonPredictionsAvailable
	}
	return Result{
		TripUpdates: updates,
		Diagnostics: Diagnostics{
			Status: StatusOK,
			Reason: reason,
			Metrics: Metrics{
				TelemetryRowsConsidered: len(request.Telemetry),
				AssignmentsConsidered:   len(request.Assignments),
				TripUpdatesEmitted:      len(updates),
			},
			Details: externalHTTPSuccessDetails(start, requestBytes, responseBytes, statusCode, len(updates)),
		},
	}, ""
}

func mapExternalTripUpdate(update externalTripUpdateDTO) TripUpdate {
	mapped := TripUpdate{
		AgencyID:             update.AgencyID,
		FeedVersionID:        update.FeedVersionID,
		EntityID:             update.EntityID,
		VehicleID:            update.VehicleID,
		TripID:               update.TripID,
		RouteID:              update.RouteID,
		StartDate:            update.StartDate,
		StartTime:            update.StartTime,
		ScheduleRelationship: update.ScheduleRelationship,
		Confidence:           update.Confidence,
		StopTimeUpdates:      make([]StopTimeUpdate, 0, len(update.StopTimeUpdates)),
	}
	for _, stop := range update.StopTimeUpdates {
		mapped.StopTimeUpdates = append(mapped.StopTimeUpdates, StopTimeUpdate{
			StopID:                stop.StopID,
			StopSequence:          stop.StopSequence,
			ArrivalTime:           stop.ArrivalTime,
			DepartureTime:         stop.DepartureTime,
			ArrivalDelaySeconds:   stop.ArrivalDelaySeconds,
			DepartureDelaySeconds: stop.DepartureDelaySeconds,
			ScheduleRelationship:  stop.ScheduleRelationship,
		})
	}
	return mapped
}

func externalHTTPFailure(failure string, start time.Time, requestBytes int, responseBytes int, statusCode int) Result {
	return Result{
		Diagnostics: Diagnostics{
			Status: StatusError,
			Reason: ReasonAdapterError,
			Details: map[string]any{
				"adapter_contract": "external_trip_updates.v1",
				"failure_type":     failure,
				"latency_ms":       time.Since(start).Milliseconds(),
				"request_bytes":    requestBytes,
				"response_bytes":   responseBytes,
				"http_status_code": statusCode,
			},
		},
	}
}

func externalHTTPSuccessDetails(start time.Time, requestBytes int, responseBytes int, statusCode int, updateCount int) map[string]any {
	return map[string]any{
		"adapter_contract":      "external_trip_updates.v1",
		"latency_ms":            time.Since(start).Milliseconds(),
		"request_bytes":         requestBytes,
		"response_bytes":        responseBytes,
		"http_status_code":      statusCode,
		"external_update_count": updateCount,
	}
}

type ExternalHTTPShadowAdapter struct {
	deterministic Adapter
	external      Adapter
}

func NewExternalHTTPShadowAdapter(deterministic Adapter, external Adapter) (*ExternalHTTPShadowAdapter, error) {
	if deterministic == nil {
		return nil, fmt.Errorf("deterministic adapter is required")
	}
	if external == nil {
		return nil, fmt.Errorf("external adapter is required")
	}
	return &ExternalHTTPShadowAdapter{deterministic: deterministic, external: external}, nil
}

func (a *ExternalHTTPShadowAdapter) Name() string {
	return AdapterNameExternalHTTPShadow
}

func (a *ExternalHTTPShadowAdapter) PredictTripUpdates(ctx context.Context, request Request) (Result, error) {
	result, err := a.deterministic.PredictTripUpdates(ctx, request)
	if err != nil {
		return result, err
	}
	start := time.Now()
	shadow, shadowErr := a.external.PredictTripUpdates(ctx, request)
	details := copyDiagnosticsDetails(result.Diagnostics.Details)
	details["external_http_shadow"] = boundedShadowDiagnostics(result, shadow, shadowErr, start)
	result.Diagnostics.Details = details
	return result, nil
}

func boundedShadowDiagnostics(deterministic Result, shadow Result, shadowErr error, start time.Time) map[string]any {
	status := shadow.Diagnostics.Status
	reason := shadow.Diagnostics.Reason
	if shadowErr != nil {
		status = StatusError
		reason = ReasonAdapterError
	}
	return map[string]any{
		"status":                           status,
		"reason":                           reason,
		"latency_ms":                       time.Since(start).Milliseconds(),
		"deterministic_trip_updates_count": len(deterministic.TripUpdates),
		"external_trip_updates_count":      len(shadow.TripUpdates),
		"count_delta":                      len(shadow.TripUpdates) - len(deterministic.TripUpdates),
	}
}

func copyDiagnosticsDetails(details map[string]any) map[string]any {
	copied := make(map[string]any, len(details))
	for key, value := range details {
		copied[key] = value
	}
	return copied
}
