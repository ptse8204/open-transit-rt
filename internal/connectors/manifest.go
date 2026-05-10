package connectors

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"path"
	"regexp"
	"slices"
	"strings"
)

const SchemaVersion = "open-transit-rt.connector.v1"

const (
	TypeTelemetrySource   = "telemetry_source"
	TypePrediction        = "prediction"
	TypeValidator         = "validator"
	TypeMonitoringExport  = "monitoring_export"
	TypeConsumerDiscovery = "consumer_discovery"
)

var SupportedTypes = []string{
	TypeTelemetrySource,
	TypePrediction,
	TypeValidator,
	TypeMonitoringExport,
	TypeConsumerDiscovery,
}

var supportedPositiveClaims = []string{
	"adapter contract only",
	"synthetic conformance tested",
	"disabled by default",
	"redacted diagnostics only",
	"no status mutation",
}

type Manifest struct {
	SchemaVersion    string            `json:"schema_version"`
	ConnectorID      string            `json:"connector_id"`
	ConnectorType    string            `json:"type"`
	DisplayName      string            `json:"display_name"`
	Description      string            `json:"description"`
	Mode             Mode              `json:"mode"`
	InputContracts   []Contract        `json:"input_contracts"`
	OutputContracts  []Contract        `json:"output_contracts"`
	FailureBehavior  FailureBehavior   `json:"failure_behavior"`
	RedactionPolicy  RedactionPolicy   `json:"redaction_policy"`
	ClaimBoundary    ClaimBoundary     `json:"claim_boundary"`
	DocsLink         string            `json:"docs_link"`
	ConformanceCases []ConformanceCase `json:"conformance_cases"`
}

type Mode struct {
	Name                         string `json:"name"`
	DisabledByDefault            bool   `json:"disabled_by_default"`
	SendsNotificationsByDefault  bool   `json:"sends_notifications_by_default,omitempty"`
	AutomatesConsumerSubmission  bool   `json:"automates_consumer_submission,omitempty"`
	MutatesStatus                bool   `json:"mutates_status,omitempty"`
	RequiresOperatorConfirmation bool   `json:"requires_operator_confirmation,omitempty"`
}

type Contract struct {
	Name                        string   `json:"name"`
	Description                 string   `json:"description"`
	Schema                      string   `json:"schema"`
	MediaTypes                  []string `json:"media_types,omitempty"`
	RequiredFields              []string `json:"required_fields,omitempty"`
	Produces                    []string `json:"produces,omitempty"`
	MutatesStatus               bool     `json:"mutates_status,omitempty"`
	SendsNotifications          bool     `json:"sends_notifications,omitempty"`
	AutomatesConsumerSubmission bool     `json:"automates_consumer_submission,omitempty"`
	RawValidatorCommand         []string `json:"raw_validator_command,omitempty"`
}

type FailureBehavior struct {
	TimeoutSeconds int    `json:"timeout_seconds"`
	RetryPolicy    string `json:"retry_policy"`
	DegradedState  string `json:"degraded_state"`
	FailClosed     bool   `json:"fail_closed"`
}

type RedactionPolicy struct {
	SecretStorage  string   `json:"secret_storage"`
	RedactFields   []string `json:"redact_fields"`
	LogsRawPayload bool     `json:"logs_raw_payload,omitempty"`
}

type ClaimBoundary struct {
	PositiveClaims []string `json:"positive_claims"`
	NotClaimed     []string `json:"not_claimed"`
}

type ConformanceCase struct {
	ID             string `json:"id"`
	Description    string `json:"description"`
	FixturePath    string `json:"fixture_path"`
	ExpectedResult string `json:"expected_result"`
}

type Violation struct {
	Field   string
	Message string
}

func (v Violation) Error() string {
	return v.Field + ": " + v.Message
}

func DecodeManifest(r io.Reader) (Manifest, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return Manifest{}, err
	}
	if violations := rejectUnsafeRawJSON(raw); len(violations) > 0 {
		return Manifest{}, violations[0]
	}
	var manifest Manifest
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&manifest); err != nil {
		return Manifest{}, fmt.Errorf("decode connector manifest: %w", err)
	}
	if violations := manifest.Validate(); len(violations) > 0 {
		return Manifest{}, violations[0]
	}
	return manifest, nil
}

func (m Manifest) Validate() []Violation {
	var violations []Violation
	requiredString(&violations, "schema_version", m.SchemaVersion)
	requiredString(&violations, "connector_id", m.ConnectorID)
	requiredString(&violations, "type", m.ConnectorType)
	requiredString(&violations, "display_name", m.DisplayName)
	requiredString(&violations, "description", m.Description)
	requiredString(&violations, "docs_link", m.DocsLink)

	if m.SchemaVersion != "" && m.SchemaVersion != SchemaVersion {
		violations = append(violations, Violation{"schema_version", "must equal " + SchemaVersion})
	}
	if m.ConnectorType != "" && !slices.Contains(SupportedTypes, m.ConnectorType) {
		violations = append(violations, Violation{"type", "unsupported connector type"})
	}
	if !connectorIDPattern.MatchString(m.ConnectorID) {
		violations = append(violations, Violation{"connector_id", "must use lowercase letters, numbers, dots, underscores, or hyphens"})
	}
	if err := validateDocsLink(m.DocsLink, "docs_link"); err != nil {
		violations = append(violations, *err)
	}

	violations = append(violations, m.Mode.validate()...)
	violations = append(violations, validateContracts("input_contracts", m.InputContracts)...)
	violations = append(violations, validateContracts("output_contracts", m.OutputContracts)...)
	violations = append(violations, m.FailureBehavior.validate()...)
	violations = append(violations, m.RedactionPolicy.validate()...)
	violations = append(violations, m.ClaimBoundary.validate()...)
	violations = append(violations, validateConformanceCases(m.ConformanceCases)...)

	if len(m.InputContracts) == 0 {
		violations = append(violations, Violation{"input_contracts", "at least one input contract is required"})
	}
	if len(m.OutputContracts) == 0 {
		violations = append(violations, Violation{"output_contracts", "at least one output contract is required"})
	}
	if len(m.ConformanceCases) == 0 {
		violations = append(violations, Violation{"conformance_cases", "at least one conformance case is required"})
	}
	if m.ConnectorType == TypeConsumerDiscovery && m.Mode.AutomatesConsumerSubmission {
		violations = append(violations, Violation{"mode.automates_consumer_submission", "consumer submission automation is not allowed"})
	}
	return violations
}

func (m Mode) validate() []Violation {
	var violations []Violation
	requiredString(&violations, "mode.name", m.Name)
	if m.SendsNotificationsByDefault {
		violations = append(violations, Violation{"mode.sends_notifications_by_default", "notification sending must not be enabled by default"})
	}
	if m.AutomatesConsumerSubmission {
		violations = append(violations, Violation{"mode.automates_consumer_submission", "consumer submission automation is not allowed"})
	}
	if m.MutatesStatus {
		violations = append(violations, Violation{"mode.mutates_status", "connectors must not mutate internal status"})
	}
	if !m.DisabledByDefault {
		violations = append(violations, Violation{"mode.disabled_by_default", "connectors must be disabled until an operator configures them"})
	}
	return violations
}

func validateContracts(field string, contracts []Contract) []Violation {
	var violations []Violation
	for i, contract := range contracts {
		prefix := fmt.Sprintf("%s[%d]", field, i)
		requiredString(&violations, prefix+".name", contract.Name)
		requiredString(&violations, prefix+".description", contract.Description)
		requiredString(&violations, prefix+".schema", contract.Schema)
		if contract.MutatesStatus {
			violations = append(violations, Violation{prefix + ".mutates_status", "connectors must not mutate internal status"})
		}
		if contract.SendsNotifications {
			violations = append(violations, Violation{prefix + ".sends_notifications", "notification output must require an explicit operator-controlled integration"})
		}
		if contract.AutomatesConsumerSubmission {
			violations = append(violations, Violation{prefix + ".automates_consumer_submission", "consumer submission automation is not allowed"})
		}
		if len(contract.RawValidatorCommand) > 0 {
			violations = append(violations, Violation{prefix + ".raw_validator_command", "raw validator commands are not allowed; use server-side validator IDs"})
		}
		violations = append(violations, validateSafeStrings(prefix, []string{contract.Name, contract.Description, contract.Schema})...)
		violations = append(violations, validateSafeStrings(prefix+".media_types", contract.MediaTypes)...)
		violations = append(violations, validateSafeStrings(prefix+".required_fields", contract.RequiredFields)...)
		violations = append(violations, validateSafeStrings(prefix+".produces", contract.Produces)...)
	}
	return violations
}

func (f FailureBehavior) validate() []Violation {
	var violations []Violation
	if f.TimeoutSeconds <= 0 || f.TimeoutSeconds > 60 {
		violations = append(violations, Violation{"failure_behavior.timeout_seconds", "must be between 1 and 60"})
	}
	requiredString(&violations, "failure_behavior.retry_policy", f.RetryPolicy)
	requiredString(&violations, "failure_behavior.degraded_state", f.DegradedState)
	if !f.FailClosed {
		violations = append(violations, Violation{"failure_behavior.fail_closed", "connectors must fail closed"})
	}
	return violations
}

func (r RedactionPolicy) validate() []Violation {
	var violations []Violation
	requiredString(&violations, "redaction_policy.secret_storage", r.SecretStorage)
	if r.SecretStorage != "" && r.SecretStorage != "none" && r.SecretStorage != "env_reference_only" {
		violations = append(violations, Violation{"redaction_policy.secret_storage", "must be none or env_reference_only"})
	}
	if len(r.RedactFields) == 0 {
		violations = append(violations, Violation{"redaction_policy.redact_fields", "at least one redacted field name is required"})
	}
	if r.LogsRawPayload {
		violations = append(violations, Violation{"redaction_policy.logs_raw_payload", "raw payload logging is not allowed"})
	}
	return violations
}

func (c ClaimBoundary) validate() []Violation {
	var violations []Violation
	if len(c.PositiveClaims) == 0 {
		violations = append(violations, Violation{"claim_boundary.positive_claims", "at least one positive claim is required"})
	}
	if len(c.NotClaimed) == 0 {
		violations = append(violations, Violation{"claim_boundary.not_claimed", "at least one claim exclusion is required"})
	}
	for i, claim := range c.PositiveClaims {
		normalized := strings.ToLower(strings.TrimSpace(claim))
		if !slices.Contains(supportedPositiveClaims, normalized) {
			violations = append(violations, Violation{fmt.Sprintf("claim_boundary.positive_claims[%d]", i), "unsupported positive claim"})
		}
		for _, banned := range bannedClaimFragments {
			if strings.Contains(normalized, banned) {
				violations = append(violations, Violation{fmt.Sprintf("claim_boundary.positive_claims[%d]", i), "claim requires retained evidence outside connector manifests"})
			}
		}
	}
	return violations
}

func validateConformanceCases(cases []ConformanceCase) []Violation {
	var violations []Violation
	for i, tc := range cases {
		prefix := fmt.Sprintf("conformance_cases[%d]", i)
		requiredString(&violations, prefix+".id", tc.ID)
		requiredString(&violations, prefix+".description", tc.Description)
		requiredString(&violations, prefix+".fixture_path", tc.FixturePath)
		requiredString(&violations, prefix+".expected_result", tc.ExpectedResult)
		if tc.ExpectedResult != "" && tc.ExpectedResult != "pass" && tc.ExpectedResult != "fail" {
			violations = append(violations, Violation{prefix + ".expected_result", "must be pass or fail"})
		}
		if tc.FixturePath != "" {
			clean := path.Clean(tc.FixturePath)
			if clean != tc.FixturePath || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "/") || !allowedFixturePath(clean) {
				violations = append(violations, Violation{prefix + ".fixture_path", "must be a clean relative path under testdata/connectors, testdata/adapter-conformance, or examples/connectors"})
			}
		}
	}
	return violations
}

func allowedFixturePath(clean string) bool {
	return strings.HasPrefix(clean, "testdata/connectors/") ||
		strings.HasPrefix(clean, "testdata/adapter-conformance/") ||
		strings.HasPrefix(clean, "examples/connectors/")
}

func validateSafeURL(value string, field string) *Violation {
	if value == "" {
		return nil
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return &Violation{field, "must be an absolute URL"}
	}
	if parsed.Scheme != "https" {
		return &Violation{field, "must use https"}
	}
	if parsed.User != nil {
		return &Violation{field, "must not include userinfo"}
	}
	host := parsed.Hostname()
	if ip := net.ParseIP(host); ip != nil && (ip.IsPrivate() || ip.IsLoopback()) {
		return &Violation{field, "must not point to private or loopback hosts"}
	}
	if strings.EqualFold(host, "localhost") || strings.HasSuffix(strings.ToLower(host), ".local") {
		return &Violation{field, "must not point to private local hosts"}
	}
	return nil
}

func validateDocsLink(value string, field string) *Violation {
	if value == "" {
		return nil
	}
	if strings.Contains(value, "://") {
		return validateSafeURL(value, field)
	}
	clean := path.Clean(value)
	if clean != value || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "/") {
		return &Violation{field, "must be a clean relative docs path or safe https URL"}
	}
	if !strings.HasPrefix(clean, "docs/") && !strings.HasPrefix(clean, "examples/") {
		return &Violation{field, "must point under docs/ or examples/"}
	}
	return nil
}

func validateSafeStrings(field string, values []string) []Violation {
	var violations []Violation
	for _, value := range values {
		if isPrivatePath(value) {
			violations = append(violations, Violation{field, "private filesystem paths are not allowed"})
		}
		if strings.Contains(value, "://") {
			if err := validateSafeURL(value, field); err != nil {
				violations = append(violations, *err)
			}
		}
	}
	return violations
}

func rejectUnsafeRawJSON(raw []byte) []Violation {
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil
	}
	var violations []Violation
	walkJSON("", decoded, &violations)
	return violations
}

func walkJSON(field string, value any, violations *[]Violation) {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			childField := key
			if field != "" {
				childField = field + "." + key
			}
			lowerKey := strings.ToLower(key)
			if secretKeyPattern.MatchString(lowerKey) && !allowedSecretField(lowerKey) {
				*violations = append(*violations, Violation{childField, "secret-bearing fields are not allowed in connector manifests"})
			}
			if lowerKey == "command" || lowerKey == "raw_command" || lowerKey == "shell" || lowerKey == "argv" || lowerKey == "args" {
				*violations = append(*violations, Violation{childField, "raw commands are not allowed in connector manifests"})
			}
			walkJSON(childField, child, violations)
		}
	case []any:
		for i, child := range typed {
			walkJSON(fmt.Sprintf("%s[%d]", field, i), child, violations)
		}
	case string:
		if secretValuePattern.MatchString(typed) {
			*violations = append(*violations, Violation{field, "secret-like values are not allowed in connector manifests"})
		}
		if isPrivatePath(typed) {
			*violations = append(*violations, Violation{field, "private filesystem paths are not allowed"})
		}
	}
}

func allowedSecretField(key string) bool {
	switch key {
	case "secret_storage":
		return true
	default:
		return false
	}
}

func requiredString(violations *[]Violation, field string, value string) {
	if strings.TrimSpace(value) == "" {
		*violations = append(*violations, Violation{field, "is required"})
	}
}

func isPrivatePath(value string) bool {
	trimmed := strings.TrimSpace(value)
	if strings.HasPrefix(trimmed, "file://") || strings.HasPrefix(trimmed, "../") || strings.Contains(trimmed, "/../") {
		return true
	}
	return strings.HasPrefix(trimmed, "/Users/") ||
		strings.HasPrefix(trimmed, "/home/") ||
		strings.HasPrefix(trimmed, "/var/") ||
		strings.HasPrefix(trimmed, "/tmp/") ||
		strings.HasPrefix(trimmed, "/etc/")
}

var (
	connectorIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)
	secretKeyPattern   = regexp.MustCompile(`(?i)(secret|password|passwd|token|api[_-]?key|private[_-]?key|credential)`)
	secretValuePattern = regexp.MustCompile(
		`(?i)\b(bearer\s+[a-z0-9._~+/=-]{12,}|sk-[a-z0-9]{12,}|ghp_[a-z0-9_]{12,}|-----begin [a-z ]*private key-----)\b`,
	)
	bannedClaimFragments = []string{
		"caltrans compliant",
		"caltrans compliance",
		"cal-itp compliant",
		"cal-itp compliance",
		"consumer accepted",
		"consumer acceptance",
		"google maps accepted",
		"apple maps accepted",
		"agency adoption",
		"production ready",
		"production readiness",
		"hosted saas",
		"paid support",
		"sla",
		"vendor compatible",
		"vendor compatibility",
		"hardware certified",
		"public launch",
		"production avl reliability",
		"production-grade eta",
		"production grade eta",
		"real vendor compatibility",
		"agency approved",
	}
)
