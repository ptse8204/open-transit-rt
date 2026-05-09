package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"open-transit-rt/internal/db"
	"open-transit-rt/internal/feed"
	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

const (
	defaultTarget      = "http://localhost:8080"
	defaultScenarioDir = "testdata/telemetry-simulator"
	defaultOutputRoot  = ".cache/telemetry-simulator"
	defaultDeviceToken = "dev-device-token"
)

type scenario struct {
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	SyntheticOnly bool            `json:"synthetic_only"`
	ReferenceTime string          `json:"reference_time"`
	Requires      []string        `json:"requires,omitempty"`
	Events        []scenarioEvent `json:"events"`
}

type scenarioEvent struct {
	Label                  string          `json:"label"`
	PreserveIdentity       bool            `json:"preserve_identity,omitempty"`
	Timestamp              string          `json:"timestamp,omitempty"`
	OffsetSeconds          *int            `json:"offset_seconds,omitempty"`
	Payload                eventPayload    `json:"payload"`
	ExpectedHTTPStatuses   []int           `json:"expected_http_statuses,omitempty"`
	ExpectedIngestStatuses []string        `json:"expected_ingest_statuses,omitempty"`
	ExpectedResponseBody   json.RawMessage `json:"expected_response_body,omitempty"`
}

type eventPayload struct {
	AgencyID  string   `json:"agency_id"`
	DeviceID  string   `json:"device_id"`
	VehicleID string   `json:"vehicle_id"`
	DriverID  string   `json:"driver_id,omitempty"`
	Timestamp string   `json:"timestamp"`
	Lat       float64  `json:"lat"`
	Lon       float64  `json:"lon"`
	Bearing   *float64 `json:"bearing,omitempty"`
	SpeedMPS  *float64 `json:"speed_mps,omitempty"`
	AccuracyM *float64 `json:"accuracy_m,omitempty"`
	TripHint  string   `json:"trip_hint,omitempty"`
}

type ingestResponse struct {
	Accepted     bool   `json:"accepted"`
	IngestStatus string `json:"ingest_status"`
	AgencyID     string `json:"agency_id"`
	VehicleID    string `json:"vehicle_id"`
	ObservedAt   string `json:"observed_at"`
	ReceivedAt   string `json:"received_at"`
	Error        string `json:"error,omitempty"`
}

type eventResult struct {
	Label             string            `json:"label"`
	RequestPayload    eventPayload      `json:"request_payload"`
	HTTPStatus        int               `json:"http_status"`
	HTTPDurationMS    int64             `json:"http_duration_ms,omitempty"`
	IngestStatus      string            `json:"ingest_status,omitempty"`
	Accepted          bool              `json:"accepted"`
	ResponseBody      json.RawMessage   `json:"response_body,omitempty"`
	ResponseSHA256    string            `json:"response_sha256,omitempty"`
	StoredTelemetryID int64             `json:"stored_telemetry_id,omitempty"`
	Matched           bool              `json:"matched,omitempty"`
	MatcherDurationMS int64             `json:"matcher_duration_ms,omitempty"`
	Assignment        *state.Assignment `json:"assignment,omitempty"`
	SkippedReason     string            `json:"skipped_reason,omitempty"`
}

type runSummary struct {
	Scenario                        string        `json:"scenario"`
	Description                     string        `json:"description"`
	SyntheticOnly                   bool          `json:"synthetic_only"`
	Target                          string        `json:"target"`
	ReferenceTime                   time.Time     `json:"reference_time"`
	DryRun                          bool          `json:"dry_run"`
	RunMatcher                      bool          `json:"run_matcher"`
	DurationMS                      int64         `json:"duration_ms"`
	MatcherTotalDurationMS          int64         `json:"matcher_total_duration_ms"`
	VehiclePositionsDebugDurationMS int64         `json:"vehicle_positions_debug_duration_ms"`
	EventsSent                      int           `json:"events_sent"`
	EventsAccepted                  int           `json:"events_accepted"`
	EventsDuplicate                 int           `json:"events_duplicate"`
	EventsOutOfOrder                int           `json:"events_out_of_order"`
	EventsRejected                  int           `json:"events_rejected"`
	OutputDirectory                 string        `json:"output_directory"`
	ExternalEvidenceCreated         bool          `json:"external_evidence_created"`
	ConsumerStatusesChanged         bool          `json:"consumer_statuses_changed"`
	ClaimsBoundary                  []string      `json:"claims_boundary"`
	ScenarioRequirements            []string      `json:"scenario_requirements,omitempty"`
	Events                          []eventResult `json:"events"`
}

type config struct {
	target                  string
	scenarioName            string
	scenarioDir             string
	agencyID                string
	deviceID                string
	vehicleID               string
	deviceToken             string
	referenceTime           string
	outputDir               string
	force                   bool
	dryRun                  bool
	runMatcher              bool
	databaseURL             string
	allowUnignoredOutputDir bool
	listScenarios           bool
	timeout                 time.Duration
}

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "telemetry simulator failed: %s\n", redactForConsole(err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) error {
	startedAt := time.Now()
	cfg := configFromEnv()
	deviceTokenFlag := ""
	databaseURLFlag := ""
	fs := flag.NewFlagSet("telemetry-simulator", flag.ContinueOnError)
	fs.SetOutput(stdout)
	fs.StringVar(&cfg.target, "target", cfg.target, "base URL for the local/reference deployment")
	fs.StringVar(&cfg.scenarioName, "scenario", cfg.scenarioName, "scenario name or JSON file path")
	fs.StringVar(&cfg.scenarioDir, "scenario-dir", cfg.scenarioDir, "directory containing synthetic scenario JSON files")
	fs.StringVar(&cfg.agencyID, "agency-id", cfg.agencyID, "override agency_id for non-preserved scenario events")
	fs.StringVar(&cfg.deviceID, "device-id", cfg.deviceID, "override device_id for non-preserved scenario events")
	fs.StringVar(&cfg.vehicleID, "vehicle-id", cfg.vehicleID, "override vehicle_id for non-preserved scenario events")
	fs.StringVar(&deviceTokenFlag, "device-token", "", "device bearer token; prefer DEVICE_TOKEN env; required for non-dry-run sends")
	fs.StringVar(&cfg.referenceTime, "reference-time", cfg.referenceTime, "override scenario reference time as RFC3339")
	fs.StringVar(&cfg.outputDir, "output-dir", cfg.outputDir, "private diagnostics output directory; default is .cache/telemetry-simulator/<timestamp>")
	fs.BoolVar(&cfg.force, "force", cfg.force, "allow writing into an existing output directory")
	fs.BoolVar(&cfg.dryRun, "dry-run", cfg.dryRun, "render diagnostics without sending HTTP requests")
	fs.BoolVar(&cfg.runMatcher, "run-matcher", cfg.runMatcher, "after accepted ingest, run the DB-backed matcher and Vehicle Positions debug snapshot")
	fs.StringVar(&databaseURLFlag, "database-url", "", "Postgres URL used only with --run-matcher; prefer DATABASE_URL env")
	fs.BoolVar(&cfg.allowUnignoredOutputDir, "allow-unignored-output-dir", cfg.allowUnignoredOutputDir, "allow OUTPUT_DIR outside repo .cache; docs/evidence is always rejected")
	fs.BoolVar(&cfg.listScenarios, "list-scenarios", cfg.listScenarios, "list available synthetic scenarios and exit")
	fs.DurationVar(&cfg.timeout, "timeout", cfg.timeout, "HTTP and DB operation timeout")
	fs.Usage = func() {
		fmt.Fprintln(stdout, "Usage: telemetry-simulator [options]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if deviceTokenFlag != "" {
		cfg.deviceToken = deviceTokenFlag
	}
	if databaseURLFlag != "" {
		cfg.databaseURL = databaseURLFlag
	}
	if cfg.deviceToken == "" && isLoopbackTarget(cfg.target) {
		cfg.deviceToken = defaultDeviceToken
	}

	if cfg.listScenarios {
		return listScenarios(stdout, cfg.scenarioDir)
	}

	sc, path, err := loadScenario(cfg.scenarioName, cfg.scenarioDir)
	if err != nil {
		return err
	}
	if !sc.SyntheticOnly {
		return fmt.Errorf("scenario %q must declare synthetic_only=true", sc.Name)
	}
	if len(sc.Events) == 0 {
		return fmt.Errorf("scenario %q has no events", sc.Name)
	}
	referenceTime, err := resolveReferenceTime(sc.ReferenceTime, cfg.referenceTime)
	if err != nil {
		return err
	}
	if err := validateTargetAndToken(cfg); err != nil {
		return err
	}
	outputDir, err := prepareOutputDir(cfg.outputDir, cfg.force, cfg.allowUnignoredOutputDir)
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Scenario: %s\n", sc.Name)
	fmt.Fprintf(stdout, "Scenario file: %s\n", path)
	fmt.Fprintf(stdout, "Target: %s/v1/telemetry\n", strings.TrimRight(cfg.target, "/"))
	fmt.Fprintf(stdout, "Output: %s\n", outputDir)
	if cfg.dryRun {
		fmt.Fprintln(stdout, "Dry run: no telemetry will be sent.")
	}

	httpClient := &http.Client{Timeout: cfg.timeout}
	var matcher *matcherRunner
	if cfg.runMatcher && !cfg.dryRun {
		matcherAgencyID := cfg.agencyID
		if firstPayload, _, expandErr := expandEvent(sc.Events[0], referenceTime, cfg); expandErr == nil {
			matcherAgencyID = firstPayload.AgencyID
		}
		matcher, err = newMatcherRunner(ctx, cfg.databaseURL, matcherAgencyID)
		if err != nil {
			return err
		}
		defer matcher.Close()
	}

	results := make([]eventResult, 0, len(sc.Events))
	for _, raw := range sc.Events {
		payload, eventTime, err := expandEvent(raw, referenceTime, cfg)
		if err != nil {
			return fmt.Errorf("expand event %q: %w", raw.Label, err)
		}
		result := eventResult{Label: raw.Label, RequestPayload: payload}
		if cfg.dryRun {
			result.SkippedReason = "dry_run"
			results = append(results, result)
			continue
		}

		httpStartedAt := time.Now()
		status, body, err := postTelemetry(ctx, httpClient, cfg.target, cfg.deviceToken, payload)
		if err != nil {
			return fmt.Errorf("post event %q: %w", raw.Label, err)
		}
		result.HTTPDurationMS = elapsedMilliseconds(httpStartedAt)
		result.HTTPStatus = status
		result.ResponseSHA256 = sha256Hex(body)
		result.ResponseBody = redactResponseBody(body)
		var parsed ingestResponse
		_ = json.Unmarshal(body, &parsed)
		result.IngestStatus = parsed.IngestStatus
		result.Accepted = parsed.Accepted
		if err := validateExpectation(raw, result); err != nil {
			return err
		}

		if matcher != nil && status == http.StatusCreated && parsed.IngestStatus == string(telemetry.IngestStatusAccepted) {
			stored, err := matcher.lookupStoredEvent(ctx, payload, eventTime)
			if err != nil {
				return fmt.Errorf("lookup stored event %q: %w", raw.Label, err)
			}
			result.StoredTelemetryID = stored.ID
			matcherStartedAt := time.Now()
			assignment, err := matcher.engine.MatchEvent(ctx, stored, referenceTime)
			if err != nil {
				return fmt.Errorf("match event %q: %w", raw.Label, err)
			}
			result.MatcherDurationMS = elapsedMilliseconds(matcherStartedAt)
			result.Matched = true
			result.Assignment = &assignment
		} else if matcher != nil {
			result.SkippedReason = "matcher_runs_only_for_http_201_accepted_events"
		}
		results = append(results, result)
	}

	summary := summarize(sc, cfg, referenceTime, outputDir, results)
	if err := writeJSON(filepath.Join(outputDir, "summary.json"), summary); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outputDir, "events.json"), results); err != nil {
		return err
	}
	if matcher != nil {
		timing, err := writeMatcherDiagnostics(ctx, matcher, referenceTime, outputDir)
		if err != nil {
			return err
		}
		summary.VehiclePositionsDebugDurationMS = timing.VehiclePositionsDebugDurationMS
	}
	summary.DurationMS = elapsedMilliseconds(startedAt)
	for _, result := range results {
		summary.MatcherTotalDurationMS += result.MatcherDurationMS
	}
	if err := writeJSON(filepath.Join(outputDir, "summary.json"), summary); err != nil {
		return err
	}
	if err := scanOutputDirRedaction(outputDir); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Events sent: %d accepted: %d duplicate: %d out_of_order: %d rejected: %d\n", summary.EventsSent, summary.EventsAccepted, summary.EventsDuplicate, summary.EventsOutOfOrder, summary.EventsRejected)
	fmt.Fprintln(stdout, "Diagnostics are private local artifacts; no evidence packet or consumer status was changed.")
	return nil
}

func configFromEnv() config {
	cfg := config{
		target:       getenv("TARGET", getenv("PUBLIC_BASE_URL", defaultTarget)),
		scenarioName: getenv("SCENARIO", "on-route"),
		scenarioDir:  getenv("SCENARIO_DIR", defaultScenarioDir),
		agencyID:     os.Getenv("AGENCY_ID"),
		deviceID:     os.Getenv("DEVICE_ID"),
		vehicleID:    os.Getenv("VEHICLE_ID"),
		deviceToken:  os.Getenv("DEVICE_TOKEN"),
		outputDir:    os.Getenv("OUTPUT_DIR"),
		databaseURL:  os.Getenv("DATABASE_URL"),
		timeout:      15 * time.Second,
	}
	if cfg.agencyID == "" {
		cfg.agencyID = "demo-agency"
	}
	if cfg.deviceID == "" {
		cfg.deviceID = "device-1"
	}
	if cfg.vehicleID == "" {
		cfg.vehicleID = "bus-1"
	}
	if raw := os.Getenv("REFERENCE_TIME"); raw != "" {
		cfg.referenceTime = raw
	}
	if raw := os.Getenv("RUN_MATCHER"); raw != "" {
		cfg.runMatcher = parseBool(raw)
	}
	if raw := os.Getenv("DRY_RUN"); raw != "" {
		cfg.dryRun = parseBool(raw)
	}
	if raw := os.Getenv("FORCE"); raw != "" {
		cfg.force = parseBool(raw)
	}
	if raw := os.Getenv("ALLOW_UNIGNORED_OUTPUT_DIR"); raw != "" {
		cfg.allowUnignoredOutputDir = parseBool(raw)
	}
	if raw := os.Getenv("TIMEOUT"); raw != "" {
		if d, err := time.ParseDuration(raw); err == nil {
			cfg.timeout = d
		}
	}
	return cfg
}

func loadScenario(name string, dir string) (scenario, string, error) {
	path := name
	if !strings.ContainsRune(name, os.PathSeparator) && filepath.Ext(name) == "" {
		path = filepath.Join(dir, name+".json")
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		return scenario{}, "", fmt.Errorf("read scenario: %w", err)
	}
	var sc scenario
	if err := json.Unmarshal(payload, &sc); err != nil {
		return scenario{}, "", fmt.Errorf("parse scenario %s: %w", path, err)
	}
	if sc.Name == "" {
		return scenario{}, "", fmt.Errorf("scenario %s has empty name", path)
	}
	return sc, path, nil
}

func listScenarios(stdout io.Writer, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("list scenarios: %w", err)
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		sc, _, err := loadScenario(strings.TrimSuffix(entry.Name(), ".json"), dir)
		if err != nil {
			return err
		}
		names = append(names, fmt.Sprintf("%s\t%s", sc.Name, sc.Description))
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintln(stdout, name)
	}
	return nil
}

func resolveReferenceTime(scenarioReference string, override string) (time.Time, error) {
	raw := firstNonEmpty(override, scenarioReference)
	if raw == "" {
		return time.Now().UTC(), nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse reference time %q: %w", raw, err)
	}
	return t.UTC(), nil
}

func expandEvent(raw scenarioEvent, referenceTime time.Time, cfg config) (eventPayload, time.Time, error) {
	payload := raw.Payload
	if !raw.PreserveIdentity {
		payload.AgencyID = firstNonEmpty(cfg.agencyID, payload.AgencyID)
		payload.DeviceID = firstNonEmpty(cfg.deviceID, payload.DeviceID)
		payload.VehicleID = firstNonEmpty(cfg.vehicleID, payload.VehicleID)
	}
	eventTime := referenceTime
	if raw.Timestamp != "" {
		parsed, err := time.Parse(time.RFC3339, raw.Timestamp)
		if err != nil {
			return eventPayload{}, time.Time{}, err
		}
		eventTime = parsed.UTC()
	}
	if raw.OffsetSeconds != nil {
		eventTime = referenceTime.Add(time.Duration(*raw.OffsetSeconds) * time.Second)
	}
	payload.Timestamp = eventTime.Format(time.RFC3339)
	if payload.AgencyID == "" || payload.DeviceID == "" || payload.VehicleID == "" {
		return eventPayload{}, time.Time{}, errors.New("agency_id, device_id, and vehicle_id are required")
	}
	if payload.Lat < -90 || payload.Lat > 90 || payload.Lon < -180 || payload.Lon > 180 {
		return eventPayload{}, time.Time{}, errors.New("lat/lon out of range")
	}
	return payload, eventTime, nil
}

func validateTargetAndToken(cfg config) error {
	parsed, err := url.ParseRequestURI(strings.TrimRight(cfg.target, "/"))
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("target URL must use http or https")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return errors.New("target URL must not include credentials, query string, or fragment")
	}
	if cfg.dryRun {
		return nil
	}
	if strings.TrimSpace(cfg.deviceToken) == "" {
		return errors.New("DEVICE_TOKEN or --device-token is required for non-dry-run simulation")
	}
	if parsed.Scheme == "http" && !isLoopbackHost(parsed.Hostname()) {
		return errors.New("non-loopback credentialed telemetry sends require https")
	}
	return nil
}

func prepareOutputDir(raw string, force bool, allowUnignored bool) (string, error) {
	if raw == "" {
		raw = filepath.Join(defaultOutputRoot, time.Now().UTC().Format("20060102T150405Z"))
	}
	clean := filepath.Clean(raw)
	abs, err := filepath.Abs(clean)
	if err != nil {
		return "", err
	}
	root, err := repoRoot()
	if err != nil {
		return "", err
	}
	evidence, err := filepath.Abs(filepath.Join(root, "docs", "evidence"))
	if err != nil {
		return "", err
	}
	if abs == evidence || strings.HasPrefix(abs, evidence+string(os.PathSeparator)) {
		return "", errors.New("simulator diagnostics must not be written under docs/evidence")
	}
	cache, err := filepath.Abs(filepath.Join(root, ".cache"))
	if err != nil {
		return "", err
	}
	if !allowUnignored && abs != cache && !strings.HasPrefix(abs, cache+string(os.PathSeparator)) {
		return "", fmt.Errorf("OUTPUT_DIR must resolve under %s unless ALLOW_UNIGNORED_OUTPUT_DIR=true", cache)
	}
	info, err := os.Lstat(abs)
	if err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return "", errors.New("refusing to write diagnostics into a symlink")
		}
		entries, err := os.ReadDir(abs)
		if err != nil {
			return "", err
		}
		if len(entries) > 0 && !force {
			return "", fmt.Errorf("output directory %s exists and is not empty; pass --force to reuse it", abs)
		}
		return abs, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if err := os.MkdirAll(abs, 0o700); err != nil {
		return "", err
	}
	return abs, nil
}

func repoRoot() (string, error) {
	dir, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return filepath.Abs(".")
		}
		dir = parent
	}
}

func postTelemetry(ctx context.Context, client *http.Client, target string, token string, payload eventPayload) (int, []byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, err
	}
	endpoint := strings.TrimRight(target, "/") + "/v1/telemetry"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, respBody, nil
}

func validateExpectation(raw scenarioEvent, result eventResult) error {
	statuses := raw.ExpectedHTTPStatuses
	if len(statuses) == 0 {
		statuses = []int{http.StatusCreated, http.StatusAccepted}
	}
	if !containsInt(statuses, result.HTTPStatus) {
		return fmt.Errorf("event %q returned HTTP %d, expected one of %v", raw.Label, result.HTTPStatus, statuses)
	}
	if len(raw.ExpectedIngestStatuses) == 0 || result.IngestStatus == "" {
		return nil
	}
	if !containsString(raw.ExpectedIngestStatuses, result.IngestStatus) {
		return fmt.Errorf("event %q returned ingest_status %q, expected one of %v", raw.Label, result.IngestStatus, raw.ExpectedIngestStatuses)
	}
	return nil
}

func summarize(sc scenario, cfg config, referenceTime time.Time, outputDir string, results []eventResult) runSummary {
	summary := runSummary{
		Scenario:                sc.Name,
		Description:             sc.Description,
		SyntheticOnly:           true,
		Target:                  strings.TrimRight(cfg.target, "/"),
		ReferenceTime:           referenceTime,
		DryRun:                  cfg.dryRun,
		RunMatcher:              cfg.runMatcher,
		OutputDirectory:         outputDir,
		ExternalEvidenceCreated: false,
		ConsumerStatusesChanged: false,
		ClaimsBoundary: []string{
			"synthetic-only diagnostics",
			"no vendor compatibility claim",
			"no production AVL reliability claim",
			"no production-grade ETA claim",
			"no CAL-ITP or Caltrans compliance claim",
		},
		ScenarioRequirements: sc.Requires,
		Events:               results,
	}
	for _, result := range results {
		if !cfg.dryRun {
			summary.EventsSent++
		}
		switch result.IngestStatus {
		case string(telemetry.IngestStatusAccepted):
			summary.EventsAccepted++
		case string(telemetry.IngestStatusDuplicate):
			summary.EventsDuplicate++
		case string(telemetry.IngestStatusOutOfOrder):
			summary.EventsOutOfOrder++
		default:
			if !cfg.dryRun {
				summary.EventsRejected++
			}
		}
	}
	return summary
}

type matcherRunner struct {
	pool        interface{ Close() }
	telemetry   *telemetry.PostgresRepository
	assignments *state.PostgresRepository
	engine      *state.Engine
	agencyID    string
}

func newMatcherRunner(ctx context.Context, databaseURL string, agencyID string) (*matcherRunner, error) {
	dbConfig := db.LoadConfigFromEnv()
	if databaseURL != "" {
		dbConfig.DatabaseURL = databaseURL
	}
	dbConfig.MaxConns = 4
	pool, err := db.Connect(ctx, dbConfig)
	if err != nil {
		return nil, err
	}
	telemetryRepo := telemetry.NewPostgresRepository(pool)
	assignmentRepo := state.NewPostgresRepository(pool)
	engine, err := state.NewEngine(gtfs.NewPostgresRepository(pool), assignmentRepo, state.DefaultConfig())
	if err != nil {
		pool.Close()
		return nil, err
	}
	return &matcherRunner{pool: pool, telemetry: telemetryRepo, assignments: assignmentRepo, engine: engine, agencyID: agencyID}, nil
}

func (m *matcherRunner) Close() {
	m.pool.Close()
}

func (m *matcherRunner) lookupStoredEvent(ctx context.Context, payload eventPayload, eventTime time.Time) (telemetry.StoredEvent, error) {
	events, err := m.telemetry.ListEvents(ctx, payload.AgencyID, 100)
	if err != nil {
		return telemetry.StoredEvent{}, err
	}
	for _, event := range events {
		if event.DeviceID == payload.DeviceID &&
			event.VehicleID == payload.VehicleID &&
			event.IngestStatus == telemetry.IngestStatusAccepted &&
			event.Timestamp.Equal(eventTime) {
			return event, nil
		}
	}
	return telemetry.StoredEvent{}, errors.New("accepted stored telemetry event not found")
}

type matcherDiagnosticsTiming struct {
	VehiclePositionsDebugDurationMS int64
}

func writeMatcherDiagnostics(ctx context.Context, matcher *matcherRunner, referenceTime time.Time, outputDir string) (matcherDiagnosticsTiming, error) {
	latest, err := matcher.telemetry.ListLatestByAgency(ctx, matcher.agencyID, 2000)
	if err != nil {
		return matcherDiagnosticsTiming{}, err
	}
	vehicleIDs := make([]string, 0, len(latest))
	for _, event := range latest {
		vehicleIDs = append(vehicleIDs, event.VehicleID)
	}
	assignments, err := matcher.assignments.ListCurrentAssignments(ctx, matcher.agencyID, vehicleIDs)
	if err != nil {
		return matcherDiagnosticsTiming{}, err
	}
	if err := writeJSON(filepath.Join(outputDir, "assignments.json"), assignments); err != nil {
		return matcherDiagnosticsTiming{}, err
	}
	builder, err := feed.NewVehiclePositionsBuilder(matcher.telemetry, matcher.assignments, feed.VehiclePositionsConfig{
		AgencyID:                  matcher.agencyID,
		MaxVehicles:               2000,
		StaleTelemetryTTL:         state.DefaultConfig().StaleThreshold,
		SuppressStaleVehicleAfter: 5 * time.Minute,
		TripConfidenceThreshold:   state.DefaultConfig().MinConfidence,
	})
	if err != nil {
		return matcherDiagnosticsTiming{}, err
	}
	vpStartedAt := time.Now()
	snapshot, err := builder.Snapshot(ctx, referenceTime)
	if err != nil {
		return matcherDiagnosticsTiming{}, err
	}
	debug, err := snapshot.MarshalDebugJSON()
	if err != nil {
		return matcherDiagnosticsTiming{}, err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "vehicle_positions_debug.json"), debug, 0o600); err != nil {
		return matcherDiagnosticsTiming{}, err
	}
	return matcherDiagnosticsTiming{VehiclePositionsDebugDurationMS: elapsedMilliseconds(vpStartedAt)}, nil
}

func writeJSON(path string, value any) error {
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	payload = append(payload, '\n')
	return os.WriteFile(path, payload, 0o600)
}

func redactResponseBody(body []byte) json.RawMessage {
	if len(body) == 0 {
		return nil
	}
	var parsed any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil
	}
	payload, err := json.Marshal(parsed)
	if err != nil {
		return nil
	}
	return payload
}

func sha256Hex(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func scanOutputDirRedaction(outputDir string) error {
	patterns := []struct {
		name string
		re   *regexp.Regexp
	}{
		{"Authorization header", regexp.MustCompile(`(?i)Authorization\s*:`)},
		{"Bearer token", regexp.MustCompile(`(?i)\bBearer\s+[A-Za-z0-9._~+/=-]+`)},
		{"DEVICE_TOKEN", regexp.MustCompile(`\bDEVICE_TOKEN\b`)},
		{"DATABASE_URL", regexp.MustCompile(`\bDATABASE_URL\b`)},
		{"Postgres password URL", regexp.MustCompile(`postgres(?:ql)?://[^/\s:@]+:[^@\s/]+@`)},
		{"Cookie header", regexp.MustCompile(`(?i)\bCookie\s*:`)},
		{"private key", regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`)},
	}
	err := filepath.WalkDir(outputDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeType != 0 {
			return nil
		}
		payload, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, pattern := range patterns {
			if pattern.re.Match(payload) {
				return fmt.Errorf("redaction scan failed for %s in %s", pattern.name, path)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("redaction scan: %w", err)
	}
	return nil
}

func redactForConsole(text string) string {
	replacements := []struct {
		re   *regexp.Regexp
		with string
	}{
		{regexp.MustCompile(`(?i)(Authorization\s*:\s*)Bearer\s+[A-Za-z0-9._~+/=-]+`), `${1}<redacted>`},
		{regexp.MustCompile(`(?i)\bBearer\s+[A-Za-z0-9._~+/=-]+`), `Bearer <redacted>`},
		{regexp.MustCompile(`(?i)(DEVICE_TOKEN\s*=\s*)[^\s]+`), `${1}<redacted>`},
		{regexp.MustCompile(`(?i)(DATABASE_URL\s*=\s*)[^\s]+`), `${1}<redacted>`},
		{regexp.MustCompile(`postgres(?:ql)?://([^:\s/@]+):([^@\s/]+)@`), `postgres://$1:<redacted>@`},
		{regexp.MustCompile(`(?i)(Cookie\s*:\s*)[^\r\n]+`), `${1}<redacted>`},
		{regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----[\s\S]*?-----END [A-Z ]*PRIVATE KEY-----`), `<redacted-private-key>`},
	}
	for _, replacement := range replacements {
		text = replacement.re.ReplaceAllString(text, replacement.with)
	}
	return text
}

func elapsedMilliseconds(startedAt time.Time) int64 {
	if startedAt.IsZero() {
		return 0
	}
	elapsed := time.Since(startedAt).Milliseconds()
	if elapsed == 0 {
		return 1
	}
	return elapsed
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func parseBool(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

func isLoopbackTarget(target string) bool {
	parsed, err := url.Parse(target)
	if err != nil {
		return false
	}
	return isLoopbackHost(parsed.Hostname())
}

func isLoopbackHost(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

func containsInt(values []int, candidate int) bool {
	for _, value := range values {
		if value == candidate {
			return true
		}
	}
	return false
}

func containsString(values []string, candidate string) bool {
	for _, value := range values {
		if value == candidate {
			return true
		}
	}
	return false
}
