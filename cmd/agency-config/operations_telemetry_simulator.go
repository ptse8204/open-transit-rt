package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
)

const (
	telemetrySimulatorScenarioDir     = "testdata/telemetry-simulator"
	telemetrySimulatorDefaultScenario = "on-route"
)

type operationsTelemetrySimulatorView struct {
	GeneratedAt        time.Time                         `json:"generated_at"`
	AgencyID           string                            `json:"agency_id"`
	Boundary           string                            `json:"boundary"`
	ScenarioDir        string                            `json:"scenario_dir"`
	TargetRules        []string                          `json:"target_rules"`
	CredentialHandling []string                          `json:"credential_handling"`
	DiagnosticsPolicy  string                            `json:"diagnostics_policy"`
	LoadError          string                            `json:"load_error,omitempty"`
	Scenarios          []operationsTelemetryScenario     `json:"scenarios"`
	Commands           []operationsTelemetryCommand      `json:"commands"`
	ClaimFlags         operationsTelemetrySimulatorClaim `json:"claim_flags"`
}

type operationsTelemetryScenario struct {
	ID                  string                       `json:"id"`
	Name                string                       `json:"name"`
	Description         string                       `json:"description"`
	SyntheticOnly       bool                         `json:"synthetic_only"`
	ReferenceTime       string                       `json:"reference_time"`
	EventCount          int                          `json:"event_count"`
	EventLabels         []string                     `json:"event_labels"`
	Requires            []string                     `json:"requires"`
	ExpectedHTTPStatus  []int                        `json:"expected_http_statuses"`
	ExpectedIngestState []string                     `json:"expected_ingest_statuses"`
	DefaultLocal        bool                         `json:"default_local"`
	NextAction          string                       `json:"next_action"`
	Commands            []operationsTelemetryCommand `json:"commands"`
	DoesNotProve        string                       `json:"does_not_prove"`
}

type operationsTelemetryCommand struct {
	ID                string `json:"id"`
	Label             string `json:"label"`
	CommandLine       string `json:"command_line"`
	WhatItDoes        string `json:"what_it_does"`
	OperatorPrep      string `json:"operator_prep"`
	FailureNextAction string `json:"failure_next_action"`
	DoesNotProve      string `json:"does_not_prove"`
}

type operationsTelemetrySimulatorClaim struct {
	BackendCommandExecutionEnabled bool `json:"backend_command_execution_enabled"`
	TelemetrySentByWebRequest      bool `json:"telemetry_sent_by_web_request"`
	DeviceTokenCollectedByBrowser  bool `json:"device_token_collected_by_browser"`
	CacheDiagnosticsReadEnabled    bool `json:"cache_diagnostics_read_enabled"`
	ExternalEvidenceCreated        bool `json:"external_evidence_created"`
	ConsumerStatusesChanged        bool `json:"consumer_statuses_changed"`
	VendorCompatibilityClaimed     bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed   bool `json:"hardware_certification_claimed"`
	ProductionAVLClaimed           bool `json:"production_avl_claimed"`
	RealRealtimeClaimed            bool `json:"real_realtime_claimed"`
	ProductionGradeETAClaimed      bool `json:"production_grade_eta_claimed"`
	ComplianceClaimed              bool `json:"compliance_claimed"`
}

type telemetrySimulatorFixture struct {
	Name          string                           `json:"name"`
	Description   string                           `json:"description"`
	SyntheticOnly bool                             `json:"synthetic_only"`
	ReferenceTime string                           `json:"reference_time"`
	Requires      []string                         `json:"requires"`
	Events        []telemetrySimulatorFixtureEvent `json:"events"`
}

type telemetrySimulatorFixtureEvent struct {
	Label                  string   `json:"label"`
	ExpectedHTTPStatuses   []int    `json:"expected_http_statuses"`
	ExpectedIngestStatuses []string `json:"expected_ingest_statuses"`
}

func (h *handler) renderTelemetrySimulator(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "telemetry-simulator")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "telemetry-simulator", page)
}

func (h *handler) renderTelemetrySimulatorJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "telemetry-simulator")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.TelemetrySimulator)
}

func buildOperationsTelemetrySimulator(page operationsPage) operationsTelemetrySimulatorView {
	scenarios, err := loadTelemetrySimulatorScenarios()
	view := operationsTelemetrySimulatorView{
		GeneratedAt: page.GeneratedAt,
		AgencyID:    page.AgencyID,
		Boundary:    "Private authenticated simulator guide only; viewing this page executes no command, sends no telemetry, reads no private diagnostics, collects no device token, creates no evidence, and changes no consumer status.",
		ScenarioDir: telemetrySimulatorScenarioDir,
		TargetRules: []string{
			"Loopback local sends may use the seeded demo credential from the operator shell.",
			"Reference deployments require HTTPS and an operator-owned private credential outside the browser.",
			"Dry runs preview the synthetic scenario shape without sending telemetry.",
		},
		CredentialHandling: []string{
			"Use the Devices page to rotate or bind credentials; store the one-time value outside this repo.",
			"Provide private credentials only in the operator shell environment when running the simulator.",
			"This browser page never asks for, stores, displays, or posts device credentials.",
		},
		DiagnosticsPolicy: "Simulator diagnostics are private local output under the operator environment. This page may name the default private output location but never reads it, copies it, or turns it into evidence.",
		Scenarios:         scenarios,
		Commands: []operationsTelemetryCommand{
			telemetrySimulatorCommand("list-scenarios", "List available scenarios", "scripts/telemetry-simulator.sh --list-scenarios", "Prints committed synthetic scenario names in the operator shell.", "Run from a clean checkout. No browser credential entry is needed.", "Confirm the checkout includes testdata/telemetry-simulator and rerun from the repo root."),
			telemetrySimulatorCommand("default-dry-run", "Preview default payload shape", "DRY_RUN=true make telemetry-simulator", "Renders the default synthetic on-route scenario without sending telemetry.", "Use before a first send to verify target and scenario shape.", "Review scenario metadata and local app setup before sending a sample event."),
			telemetrySimulatorCommand("default-send", "Send default local sample", "make telemetry-simulator", "Sends the default synthetic on-route event through authenticated telemetry ingest from the operator shell.", "Start the local app and use the seeded loopback demo credential or a private shell environment for reference deployments.", "Review device binding, target URL, and local app health; keep diagnostics private."),
			telemetrySimulatorCommand("default-matcher", "Run matcher diagnostics after ingest", "RUN_MATCHER=true make telemetry-simulator", "After accepted ingest, runs private DB-backed matching and Vehicle Positions debug diagnostics from the operator shell.", "Use only where the operator machine has local database access. Keep database URLs out of shell history where possible.", "Review DB access, active GTFS, and accepted telemetry before rerunning."),
		},
		ClaimFlags: operationsTelemetrySimulatorClaim{},
	}
	if err != nil {
		view.LoadError = err.Error()
	}
	return view
}

func loadTelemetrySimulatorScenarios() ([]operationsTelemetryScenario, error) {
	root, err := operationsSourceRepoRoot()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(root, filepath.FromSlash(telemetrySimulatorScenarioDir))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read simulator scenario directory: %w", err)
	}
	scenarios := make([]operationsTelemetryScenario, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read simulator scenario %s: %w", entry.Name(), err)
		}
		var fixture telemetrySimulatorFixture
		if err := json.Unmarshal(raw, &fixture); err != nil {
			return nil, fmt.Errorf("decode simulator scenario %s: %w", entry.Name(), err)
		}
		name := strings.TrimSpace(fixture.Name)
		if name == "" {
			name = strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		}
		if !fixture.SyntheticOnly {
			return nil, fmt.Errorf("simulator scenario %s must declare synthetic_only=true", name)
		}
		scenarios = append(scenarios, operationsTelemetryScenario{
			ID:                  name,
			Name:                name,
			Description:         strings.TrimSpace(fixture.Description),
			SyntheticOnly:       fixture.SyntheticOnly,
			ReferenceTime:       strings.TrimSpace(fixture.ReferenceTime),
			EventCount:          len(fixture.Events),
			EventLabels:         telemetrySimulatorEventLabels(fixture.Events),
			Requires:            cleanLaunchpadList(fixture.Requires),
			ExpectedHTTPStatus:  telemetrySimulatorExpectedHTTPStatuses(fixture.Events),
			ExpectedIngestState: telemetrySimulatorExpectedIngestStatuses(fixture.Events),
			DefaultLocal:        name == telemetrySimulatorDefaultScenario,
			NextAction:          telemetrySimulatorNextAction(name, fixture.Requires),
			Commands:            telemetrySimulatorScenarioCommands(name),
			DoesNotProve:        "Synthetic simulator scenarios do not prove real vendor compatibility, hardware certification, fleet reliability, real service telemetry quality, production-grade ETA quality, consumer acceptance, public launch, or CAL-ITP/Caltrans compliance.",
		})
	}
	sort.SliceStable(scenarios, func(i, j int) bool {
		return telemetrySimulatorScenarioOrder(scenarios[i].Name) < telemetrySimulatorScenarioOrder(scenarios[j].Name)
	})
	return scenarios, nil
}

func operationsSourceRepoRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolve operations source path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..")), nil
}

func telemetrySimulatorEventLabels(events []telemetrySimulatorFixtureEvent) []string {
	labels := make([]string, 0, len(events))
	for _, event := range events {
		label := strings.TrimSpace(event.Label)
		if label != "" {
			labels = append(labels, label)
		}
	}
	return labels
}

func telemetrySimulatorExpectedHTTPStatuses(events []telemetrySimulatorFixtureEvent) []int {
	seen := make(map[int]bool)
	var out []int
	for _, event := range events {
		for _, status := range event.ExpectedHTTPStatuses {
			if !seen[status] {
				seen[status] = true
				out = append(out, status)
			}
		}
	}
	sort.Ints(out)
	return out
}

func telemetrySimulatorExpectedIngestStatuses(events []telemetrySimulatorFixtureEvent) []string {
	seen := make(map[string]bool)
	var out []string
	for _, event := range events {
		for _, status := range event.ExpectedIngestStatuses {
			status = strings.TrimSpace(status)
			if status != "" && !seen[status] {
				seen[status] = true
				out = append(out, status)
			}
		}
	}
	sort.Strings(out)
	return out
}

func telemetrySimulatorScenarioCommands(name string) []operationsTelemetryCommand {
	scenario := strings.TrimSpace(name)
	if scenario == "" {
		scenario = telemetrySimulatorDefaultScenario
	}
	return []operationsTelemetryCommand{
		telemetrySimulatorCommand("dry-run-"+scenario, "Dry run "+scenario, "SCENARIO="+scenario+" DRY_RUN=true make telemetry-simulator", "Previews this synthetic scenario without sending telemetry.", "Run from the repo root after selecting the scenario.", "Review fixture requirements and local setup before trying a send."),
		telemetrySimulatorCommand("send-"+scenario, "Send "+scenario, "SCENARIO="+scenario+" make telemetry-simulator", "Sends this synthetic scenario through authenticated ingest from the operator shell.", "Start the app first and provide any private credential through the operator shell environment.", "Review target URL, device binding, and scenario requirements; keep diagnostics private."),
		telemetrySimulatorCommand("matcher-"+scenario, "Run matcher for "+scenario, "SCENARIO="+scenario+" RUN_MATCHER=true make telemetry-simulator", "Runs optional private matcher diagnostics after accepted ingest.", "Use only with local database access and an active matching GTFS feed.", "Confirm accepted telemetry, GTFS fixture requirements, and database access before rerunning."),
	}
}

func telemetrySimulatorCommand(id string, label string, command string, what string, prep string, failure string) operationsTelemetryCommand {
	return operationsTelemetryCommand{
		ID:                id,
		Label:             label,
		CommandLine:       command,
		WhatItDoes:        what,
		OperatorPrep:      prep,
		FailureNextAction: failure,
		DoesNotProve:      "This command is a private synthetic diagnostic. It does not prove real vendor compatibility, hardware certification, production fleet reliability, real service telemetry quality, production-grade ETA quality, consumer acceptance, public launch, or CAL-ITP/Caltrans compliance.",
	}
}

func telemetrySimulatorNextAction(name string, requires []string) string {
	if name == telemetrySimulatorDefaultScenario {
		return "Start the local app, confirm the demo device binding exists, then dry-run before sending the default sample."
	}
	if len(requires) == 0 {
		return "Dry-run this synthetic scenario first, then confirm local app setup before sending a sample."
	}
	return "Review the fixture requirements, prepare only synthetic or operator-owned local data, then dry-run before sending."
}

func telemetrySimulatorScenarioOrder(name string) int {
	order := map[string]int{
		"on-route":         10,
		"stale":            20,
		"out-of-order":     30,
		"unknown-device":   40,
		"low-quality-gps":  50,
		"after-midnight":   60,
		"block-transition": 70,
	}
	if value, ok := order[name]; ok {
		return value
	}
	return 1000
}
