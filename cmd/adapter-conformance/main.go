package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
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
	"consumer_discovery": {
		"feed-url-metadata",
		"status-mutation-blocked",
		"submission-automation-blocked",
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
	case "run", "manifest", "telemetry", "prediction", "validator", "monitoring", "consumer_discovery":
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
  adapter-conformance consumer_discovery --suite testdata/adapter-conformance
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
	case "consumer_discovery":
		switch tc.Scenario {
		case "feed-url-metadata":
			required = []string{"offline", "prepared_only"}
		case "status-mutation-blocked":
			required = []string{"offline", "status_mutation_blocked"}
		case "submission-automation-blocked":
			required = []string{"offline", "submission_blocked"}
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
	if err := scanUnsafeJSON(body, label); err != nil {
		return err
	}
	return nil
}

func scanUnsafeJSON(body []byte, label string) error {
	var decoded any
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil
	}
	return walkUnsafeJSON(label, "", decoded)
}

func walkUnsafeJSON(label string, field string, value any) error {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			childField := key
			if field != "" {
				childField = field + "." + key
			}
			lowerKey := strings.ToLower(key)
			if unsafeFixtureSecretKeyPattern.MatchString(lowerKey) {
				return fmt.Errorf("%s contains forbidden text in %s: secret-bearing fields are not allowed", label, childField)
			}
			if lowerKey == "command" || lowerKey == "raw_command" || lowerKey == "shell" || lowerKey == "argv" || lowerKey == "args" {
				return fmt.Errorf("%s contains forbidden text in %s: raw command fields are not allowed", label, childField)
			}
			if err := walkUnsafeJSON(label, childField, child); err != nil {
				return err
			}
		}
	case []any:
		for i, child := range typed {
			childField := fmt.Sprintf("%s[%d]", field, i)
			if field == "" {
				childField = fmt.Sprintf("[%d]", i)
			}
			if err := walkUnsafeJSON(label, childField, child); err != nil {
				return err
			}
		}
	case string:
		if unsafeFixtureSecretValuePattern.MatchString(typed) {
			return fmt.Errorf("%s contains forbidden text in %s: secret-like values are not allowed", label, field)
		}
		if isPrivateFixturePath(typed) {
			return fmt.Errorf("%s contains forbidden text in %s: private filesystem paths are not allowed", label, field)
		}
		if containsUnsafeFixtureEndpoint(typed) {
			return fmt.Errorf("%s contains forbidden text in %s: private endpoint strings are not allowed", label, field)
		}
	}
	return nil
}

func isPrivateFixturePath(value string) bool {
	trimmed := strings.TrimSpace(value)
	return strings.HasPrefix(trimmed, "file://") ||
		strings.HasPrefix(trimmed, "../") ||
		strings.Contains(trimmed, "/../") ||
		strings.HasPrefix(trimmed, "/Users/") ||
		strings.HasPrefix(trimmed, "/home/") ||
		strings.HasPrefix(trimmed, "/var/") ||
		strings.HasPrefix(trimmed, "/tmp/") ||
		strings.HasPrefix(trimmed, "/etc/")
}

func containsUnsafeFixtureEndpoint(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return false
	}
	if strings.Contains(normalized, "localhost") || strings.Contains(normalized, "[::1]") || strings.Contains(normalized, "::1") || strings.Contains(normalized, ".local") {
		return true
	}
	return unsafeFixturePrivateEndpointPattern.MatchString(normalized)
}

var (
	unsafeFixturePrivateEndpointPattern = regexp.MustCompile(`(^|[^0-9a-z])((10|127|0)\.\d{1,3}\.\d{1,3}\.\d{1,3}|192\.168\.\d{1,3}\.\d{1,3}|172\.(1[6-9]|2[0-9]|3[0-1])\.\d{1,3}\.\d{1,3})(:[0-9]{1,5})?($|[^0-9a-z])`)
	unsafeFixtureSecretKeyPattern       = regexp.MustCompile(`(?i)(secret|password|passwd|token|api[_-]?key|private[_-]?key|credential)`)
	unsafeFixtureSecretValuePattern     = regexp.MustCompile(`(?i)\b(bearer\s+[a-z0-9._~+/=-]{12,}|sk-[a-z0-9]{12,}|ghp_[a-z0-9_]{12,}|-----begin [a-z ]*private key-----)\b`)
)

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
