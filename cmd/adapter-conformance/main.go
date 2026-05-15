package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"open-transit-rt/internal/connectors"
)

const suiteSchemaVersion = "open-transit-rt.adapter_conformance.v1"

var requiredScenarios = map[string][]string{
	"telemetry": {
		"malformed",
		"stale",
		"future",
		"wrong-agency",
		"unknown-device",
		"low-quality",
		"duplicate",
		"out-of-order",
		"missing-required-field",
		"invalid-coordinate",
	},
	"prediction": {
		"timeout",
		"malformed",
		"stale",
		"wrong-agency",
		"low-confidence",
		"missing-vehicle-positions-ref",
		"public-mutation-attempt",
	},
	"validator": {
		"allowlist",
		"raw-command",
	},
	"monitoring": {
		"redaction",
		"no-send",
		"unredacted-destination",
	},
}

type suite struct {
	SchemaVersion      string     `json:"schema_version"`
	SyntheticOnly      bool       `json:"synthetic_only"`
	ConnectorManifests []string   `json:"connector_manifests"`
	Cases              []testCase `json:"cases"`
}

type testCase struct {
	ID              string   `json:"id"`
	Type            string   `json:"type"`
	Scenario        string   `json:"scenario"`
	Fixture         string   `json:"fixture"`
	ExpectedOutcome string   `json:"expected_outcome"`
	Assertions      []string `json:"assertions"`
	SyntheticOnly   bool     `json:"synthetic_only"`
}

type fixtureHeader struct {
	SyntheticOnly bool `json:"synthetic_only"`
}

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(code)
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		usage(stdout)
		return 0
	}

	command := args[0]
	switch command {
	case "run", "manifest", "telemetry", "prediction", "validator", "monitoring":
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", command)
		usage(stderr)
		return 2
	}

	cfg, err := parseCommand(command, args[1:], stdout)
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return 2
	}
	if cfg.help {
		commandUsage(command, stdout)
		return 0
	}

	if err := execute(command, cfg.suiteDir, stdout); err != nil {
		fmt.Fprintf(stderr, "adapter conformance failed: %s\n", err)
		return 1
	}
	return 0
}

type commandConfig struct {
	suiteDir string
	help     bool
}

func parseCommand(command string, args []string, out io.Writer) (commandConfig, error) {
	cfg := commandConfig{suiteDir: "testdata/adapter-conformance"}
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(out)
	fs.StringVar(&cfg.suiteDir, "suite", cfg.suiteDir, "synthetic adapter conformance suite directory")
	fs.BoolVar(&cfg.help, "help", false, "show command help")
	if err := fs.Parse(args); err != nil {
		return cfg, err
	}
	if fs.NArg() != 0 {
		return cfg, fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}
	return cfg, nil
}

func usage(w io.Writer) {
	fmt.Fprintln(w, `Usage:
  adapter-conformance help
  adapter-conformance manifest --suite testdata/adapter-conformance
  adapter-conformance telemetry --suite testdata/adapter-conformance
  adapter-conformance prediction --suite testdata/adapter-conformance
  adapter-conformance validator --suite testdata/adapter-conformance
  adapter-conformance monitoring --suite testdata/adapter-conformance
  adapter-conformance run --suite testdata/adapter-conformance

Offline only: validates synthetic fixtures and connector manifests without
sending network traffic, running validators, contacting consumers, writing
evidence, or mutating repository state.`)
}

func commandUsage(command string, w io.Writer) {
	fmt.Fprintf(w, "Usage: adapter-conformance %s --suite testdata/adapter-conformance\n", command)
}

func execute(command string, suiteDir string, stdout io.Writer) error {
	root, err := repoRoot()
	if err != nil {
		return err
	}
	suitePath, err := safeSuiteDir(root, suiteDir)
	if err != nil {
		return err
	}
	s, err := loadSuite(suitePath)
	if err != nil {
		return err
	}
	if command == "manifest" || command == "run" {
		if err := validateManifests(root, s); err != nil {
			return err
		}
	}
	if command != "manifest" {
		filter := ""
		if command != "run" {
			filter = command
		}
		if err := validateCases(suitePath, s, filter); err != nil {
			return err
		}
	}
	if command == "run" {
		fmt.Fprintf(stdout, "adapter conformance suite passed: %s\n", rel(root, suitePath))
	} else {
		fmt.Fprintf(stdout, "adapter conformance %s passed: %s\n", command, rel(root, suitePath))
	}
	return nil
}

func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		next := filepath.Dir(dir)
		if next == dir {
			return "", errors.New("could not find repo root")
		}
		dir = next
	}
}

func safeSuiteDir(root string, raw string) (string, error) {
	if raw == "" {
		return "", errors.New("--suite is required")
	}
	clean := filepath.Clean(raw)
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") || strings.Contains(clean, string(filepath.Separator)+".."+string(filepath.Separator)) {
		return "", errors.New("--suite must be a clean relative path")
	}
	if strings.Contains(filepath.ToSlash(clean), "docs/evidence") || strings.Contains(filepath.ToSlash(clean), "evidence") {
		return "", errors.New("--suite must not point at evidence paths")
	}
	path := filepath.Join(root, clean)
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("--suite is not a directory: %s", raw)
	}
	return path, nil
}

func loadSuite(suiteDir string) (suite, error) {
	body, err := os.ReadFile(filepath.Join(suiteDir, "suite.json"))
	if err != nil {
		return suite{}, err
	}
	if err := scanUnsafe(body, "suite.json"); err != nil {
		return suite{}, err
	}
	var s suite
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&s); err != nil {
		return suite{}, fmt.Errorf("decode suite.json: %w", err)
	}
	if s.SchemaVersion != suiteSchemaVersion {
		return suite{}, fmt.Errorf("suite schema_version must be %s", suiteSchemaVersion)
	}
	if !s.SyntheticOnly {
		return suite{}, errors.New("suite must be marked synthetic_only")
	}
	if len(s.ConnectorManifests) == 0 {
		return suite{}, errors.New("suite must reference connector manifests")
	}
	if len(s.Cases) == 0 {
		return suite{}, errors.New("suite must include conformance cases")
	}
	return s, nil
}

func validateManifests(root string, s suite) error {
	for _, manifestPath := range s.ConnectorManifests {
		path, err := cleanRepoPath(root, manifestPath)
		if err != nil {
			return fmt.Errorf("manifest %s: %w", manifestPath, err)
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		_, decodeErr := connectors.DecodeManifest(f)
		closeErr := f.Close()
		if decodeErr != nil {
			return fmt.Errorf("manifest %s: %w", manifestPath, decodeErr)
		}
		if closeErr != nil {
			return closeErr
		}
	}
	return nil
}

func validateCases(suiteDir string, s suite, filter string) error {
	seen := map[string]map[string]bool{}
	for caseType := range requiredScenarios {
		seen[caseType] = map[string]bool{}
	}
	ids := map[string]bool{}
	for _, tc := range s.Cases {
		if filter != "" && tc.Type != filter {
			continue
		}
		if err := validateCase(suiteDir, tc); err != nil {
			return err
		}
		if ids[tc.ID] {
			return fmt.Errorf("duplicate case id %s", tc.ID)
		}
		ids[tc.ID] = true
		seen[tc.Type][tc.Scenario] = true
	}
	requiredTypes := []string{filter}
	if filter == "" {
		requiredTypes = keys(requiredScenarios)
	}
	for _, caseType := range requiredTypes {
		for _, scenario := range requiredScenarios[caseType] {
			if !seen[caseType][scenario] {
				return fmt.Errorf("missing required %s scenario %s", caseType, scenario)
			}
		}
	}
	return nil
}

func validateCase(suiteDir string, tc testCase) error {
	if strings.TrimSpace(tc.ID) == "" || strings.TrimSpace(tc.Type) == "" || strings.TrimSpace(tc.Scenario) == "" {
		return fmt.Errorf("case has missing id, type, or scenario: %+v", tc)
	}
	if !tc.SyntheticOnly {
		return fmt.Errorf("case %s must be synthetic_only", tc.ID)
	}
	if _, ok := requiredScenarios[tc.Type]; !ok {
		return fmt.Errorf("case %s has unsupported type %s", tc.ID, tc.Type)
	}
	if !contains(requiredScenarios[tc.Type], tc.Scenario) {
		return fmt.Errorf("case %s has unsupported scenario %s", tc.ID, tc.Scenario)
	}
	switch tc.ExpectedOutcome {
	case "reject", "withhold", "allow", "redact", "no_send":
	default:
		return fmt.Errorf("case %s has unsupported expected_outcome %s", tc.ID, tc.ExpectedOutcome)
	}
	if err := requiredAssertions(tc); err != nil {
		return err
	}
	path, err := cleanSuitePath(suiteDir, tc.Fixture)
	if err != nil {
		return fmt.Errorf("case %s fixture: %w", tc.ID, err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := scanUnsafe(body, tc.Fixture); err != nil {
		return err
	}
	var header fixtureHeader
	if err := json.Unmarshal(body, &header); err != nil {
		return fmt.Errorf("case %s fixture is invalid JSON: %w", tc.ID, err)
	}
	if !header.SyntheticOnly {
		return fmt.Errorf("case %s fixture must be synthetic_only", tc.ID)
	}
	return nil
}

func requiredAssertions(tc testCase) error {
	required := []string{"offline", "fail_closed"}
	switch tc.Type {
	case "validator":
		switch tc.Scenario {
		case "allowlist":
			required = []string{"offline", "allowlisted_validator_id"}
		case "raw-command":
			required = []string{"offline", "command_blocked"}
		default:
			required = []string{"offline"}
		}
	case "monitoring":
		switch tc.Scenario {
		case "redaction":
			required = []string{"offline", "redacted"}
		case "no-send":
			required = []string{"offline", "no_send"}
		case "unredacted-destination":
			required = []string{"offline", "redacted", "no_send"}
		default:
			required = []string{"offline"}
		}
	}
	for _, assertion := range required {
		if !contains(tc.Assertions, assertion) {
			return fmt.Errorf("case %s missing assertion %s", tc.ID, assertion)
		}
	}
	return nil
}

func cleanRepoPath(root string, raw string) (string, error) {
	clean := filepath.Clean(raw)
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") || strings.Contains(clean, string(filepath.Separator)+".."+string(filepath.Separator)) {
		return "", errors.New("path must be clean and repo-relative")
	}
	if strings.Contains(filepath.ToSlash(clean), "docs/evidence") {
		return "", errors.New("path must not be under docs/evidence")
	}
	return filepath.Join(root, clean), nil
}

func cleanSuitePath(suiteDir string, raw string) (string, error) {
	clean := filepath.Clean(raw)
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") || strings.Contains(clean, string(filepath.Separator)+".."+string(filepath.Separator)) {
		return "", errors.New("path must be clean and suite-relative")
	}
	return filepath.Join(suiteDir, clean), nil
}

func scanUnsafe(body []byte, label string) error {
	text := strings.ToLower(string(body))
	for _, forbidden := range []string{
		"authorization: bearer",
		"database_url=",
		"-----begin private key-----",
		"webhook.example",
		"consumer accepted",
		"caltrans compliant",
		"production ready",
		"vendor compatible",
	} {
		if strings.Contains(text, forbidden) {
			return fmt.Errorf("%s contains forbidden text %q", label, forbidden)
		}
	}
	return nil
}

func rel(root string, path string) string {
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(relative)
}

func keys(values map[string][]string) []string {
	out := make([]string, 0, len(values))
	for key := range values {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func contains(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}
