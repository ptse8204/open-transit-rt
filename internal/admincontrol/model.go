package admincontrol

import (
	"strings"
	"time"
	"unicode/utf8"
)

const (
	MaxSummaryLength    = 600
	MaxNextActionLength = 240
	MaxErrorCount       = 4
	MaxErrorCodeLength  = 64
	MaxErrorTextLength  = 240
)

type Status string

const (
	StatusOK          Status = "ok"
	StatusNeedsReview Status = "needs_review"
	StatusBlocked     Status = "blocked"
	StatusFailed      Status = "failed"
)

type LadderLevel string

const (
	LevelReadOnlyRefresh       LadderLevel = "read_only_refresh"
	LevelDryRun                LadderLevel = "dry_run"
	LevelReversiblePrivate     LadderLevel = "reversible_private_change"
	LevelPublishActivate       LadderLevel = "publish_activate"
	LevelDestructiveHardRevert LadderLevel = "destructive_or_hard_to_reverse"
)

type Definition struct {
	Action              string      `json:"action"`
	Label               string      `json:"label"`
	LadderLevel         LadderLevel `json:"ladder_level"`
	RequiredRole        string      `json:"required_role"`
	Confirmation        string      `json:"confirmation"`
	PublicFeedImpact    string      `json:"public_feed_impact"`
	PrivateImpact       string      `json:"private_impact"`
	RollbackPath        string      `json:"rollback_path"`
	DisabledByDefault   bool        `json:"disabled_by_default"`
	TechnicalHandoff    string      `json:"technical_handoff"`
	DoesNotProve        string      `json:"does_not_prove"`
	ServerOwnedMappings []string    `json:"server_owned_mappings"`
}

type Result struct {
	Action      string     `json:"action"`
	Status      Status     `json:"status"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt time.Time  `json:"completed_at"`
	Summary     string     `json:"summary"`
	NextActions []string   `json:"next_actions"`
	ClaimFlags  ClaimFlags `json:"claim_flags"`
	Errors      []Error    `json:"errors,omitempty"`
}

type ClaimFlags struct {
	ExternalEvidenceCreated      bool `json:"external_evidence_created"`
	ConsumerStatusesChanged      bool `json:"consumer_statuses_changed"`
	ComplianceClaimed            bool `json:"compliance_claimed"`
	ProductionReadinessClaimed   bool `json:"production_readiness_claimed"`
	AgencyApprovalClaimed        bool `json:"agency_approval_claimed"`
	ConsumerAcceptanceClaimed    bool `json:"consumer_acceptance_claimed"`
	PublicLaunchClaimed          bool `json:"public_launch_claimed"`
	HostedSaaSClaimed            bool `json:"hosted_saas_claimed"`
	VendorCompatibilityClaimed   bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed bool `json:"hardware_certification_claimed"`
	SLAClaimed                   bool `json:"sla_claimed"`
	UptimeGuaranteeClaimed       bool `json:"uptime_guarantee_claimed"`
	ProductionGradeETAClaimed    bool `json:"production_grade_eta_claimed"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewResult(action string, status Status, startedAt time.Time, summary string, nextActions []string, errors []Error) Result {
	completedAt := time.Now().UTC().Truncate(time.Second)
	if startedAt.IsZero() {
		startedAt = completedAt
	}
	return Result{
		Action:      cleanText(action, MaxErrorCodeLength),
		Status:      normalizeStatus(status),
		StartedAt:   startedAt.UTC().Truncate(time.Second),
		CompletedAt: completedAt,
		Summary:     cleanText(summary, MaxSummaryLength),
		NextActions: cleanTextList(nextActions, MaxNextActionLength),
		ClaimFlags:  ClaimFlags{},
		Errors:      cleanErrors(errors),
	}
}

func normalizeStatus(status Status) Status {
	switch status {
	case StatusOK, StatusNeedsReview, StatusBlocked, StatusFailed:
		return status
	default:
		return StatusBlocked
	}
}

func cleanTextList(values []string, maxRunes int) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		cleaned := cleanText(value, maxRunes)
		if cleaned == "" {
			continue
		}
		out = append(out, cleaned)
	}
	if out == nil {
		return []string{}
	}
	return out
}

func cleanErrors(values []Error) []Error {
	out := make([]Error, 0, len(values))
	for _, value := range values {
		if len(out) >= MaxErrorCount {
			break
		}
		code := cleanText(value.Code, MaxErrorCodeLength)
		message := cleanText(value.Message, MaxErrorTextLength)
		if code == "" && message == "" {
			continue
		}
		out = append(out, Error{Code: code, Message: message})
	}
	return out
}

func cleanText(value string, maxRunes int) string {
	cleaned := strings.TrimSpace(strings.Join(strings.Fields(value), " "))
	if cleaned == "" {
		return ""
	}
	if containsPrivateDiagnostic(cleaned) {
		return "redacted private diagnostic"
	}
	if maxRunes <= 0 || utf8.RuneCountInString(cleaned) <= maxRunes {
		return cleaned
	}
	runes := []rune(cleaned)
	return string(runes[:maxRunes])
}

func containsPrivateDiagnostic(value string) bool {
	lower := strings.ToLower(value)
	for _, marker := range []string{
		"/users/",
		"/tmp/",
		"postgres://",
		"database_url",
		"authorization:",
		"bearer ",
		"cookie:",
		"set-cookie:",
		"raw_report",
		"raw validator",
		"validator_command",
		"stdout",
		"stderr",
		"argv",
		"token_hash",
	} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func ValidationHealthRefreshDefinition() Definition {
	return Definition{
		Action:           "validation_health.refresh",
		Label:            "Refresh validator health summary",
		LadderLevel:      LevelReadOnlyRefresh,
		RequiredRole:     "read_only",
		Confirmation:     "Refreshes the derived validator health summary for the authenticated agency.",
		PublicFeedImpact: "No public feed output changes.",
		PrivateImpact:    "Writes nothing; the response is computed from existing private records and server-owned artifact checks.",
		RollbackPath:     "No rollback is required because the refresh writes nothing.",
		TechnicalHandoff: "Use validator tooling setup or off-host validation docs if tooling or artifacts are blocked.",
		DoesNotProve:     "This private view does not show compliance, consumer acceptance, final-root readiness, hosted availability, production readiness, vendor compatibility, SLA, or ETA quality.",
		ServerOwnedMappings: []string{
			"validator identifiers",
			"artifact lookup",
			"report normalization",
		},
	}
}

func ValidationHealthRunAllDefinition() Definition {
	return Definition{
		Action:           "validation_health.run_all",
		Label:            "Run allowlisted validators",
		LadderLevel:      LevelReversiblePrivate,
		RequiredRole:     "admin",
		Confirmation:     "Runs server-owned validator health checks for the authenticated agency.",
		PublicFeedImpact: "No public feed output changes.",
		PrivateImpact:    "Successful validator runs may store normal validation_report rows only.",
		RollbackPath:     "Review or supersede validation_report rows with a later validator run; no browser rollback is required for public feeds.",
		TechnicalHandoff: "Use validator tooling setup or off-host validation docs if tooling or artifacts are blocked.",
		DoesNotProve:     "This private view does not show compliance, consumer acceptance, final-root readiness, hosted availability, production readiness, vendor compatibility, SLA, or ETA quality.",
		ServerOwnedMappings: []string{
			"validator identifiers",
			"validator binary paths",
			"artifact lookup",
			"timeouts",
			"report normalization",
		},
	}
}

func AlertsCancellationPreviewDefinition() Definition {
	return Definition{
		Action:            "alerts.cancellation_reconcile.preview",
		Label:             "Preview canceled-trip alert reconciliation",
		LadderLevel:       LevelDryRun,
		RequiredRole:      "operator",
		Confirmation:      "Computes a private preview from existing canceled-trip overrides and missing-alert review incidents.",
		PublicFeedImpact:  "No public feed output changes.",
		PrivateImpact:     "Writes nothing; reports aggregate preview counts only.",
		RollbackPath:      "No rollback is required because the preview writes nothing.",
		DisabledByDefault: true,
		TechnicalHandoff:  "Use the Alerts Console reconciliation action only after agency review of canceled-trip overrides.",
		DoesNotProve:      "This private view does not show compliance, consumer acceptance, consumer display, public launch, hosted availability, production readiness, vendor compatibility, hardware certification, SLA, or ETA quality.",
		ServerOwnedMappings: []string{
			"canceled-trip override lookup",
			"missing-alert review lookup",
			"alert entity mapping",
		},
	}
}

func RealtimeQualityBacktestDefinition() Definition {
	return Definition{
		Action:            "realtime_quality.backtest.dry_run",
		Label:             "Run private realtime quality backtest",
		LadderLevel:       LevelDryRun,
		RequiredRole:      "operator",
		Confirmation:      "Runs a bounded private/local aggregate backtest from server-owned or committed synthetic inputs.",
		PublicFeedImpact:  "No public feed output changes.",
		PrivateImpact:     "Writes only private aggregate diagnostics under allowed local diagnostic output when explicitly run by an operator.",
		RollbackPath:      "Remove or supersede private diagnostic output; no public feed rollback is required.",
		DisabledByDefault: true,
		TechnicalHandoff:  "Use realtime-quality backtest docs when observed/prediction inputs are missing or blocked.",
		DoesNotProve:      "This private view does not show compliance, consumer acceptance, real-world ETA accuracy, production-grade ETA quality, public launch, hosted availability, production readiness, vendor compatibility, hardware certification, or SLA.",
		ServerOwnedMappings: []string{
			"observed event input",
			"prediction sample input",
			"diagnostic output guard",
			"redaction scan",
		},
	}
}

func ConnectorConformanceReviewDefinition() Definition {
	return Definition{
		Action:            "connectors.conformance.review",
		Label:             "Review connector conformance results",
		LadderLevel:       LevelDryRun,
		RequiredRole:      "operator",
		Confirmation:      "Reviews committed synthetic connector manifests and conformance fixture coverage.",
		PublicFeedImpact:  "No public feed output changes.",
		PrivateImpact:     "Writes nothing; reads committed synthetic manifests and fixtures only.",
		RollbackPath:      "No rollback is required because the review writes nothing.",
		DisabledByDefault: true,
		TechnicalHandoff:  "Use connector conformance docs if manifests or synthetic fixtures are missing or blocked.",
		DoesNotProve:      "This private view does not show compliance, consumer acceptance, vendor compatibility, hardware certification, public launch, hosted availability, production readiness, SLA, production AVL reliability, or ETA quality.",
		ServerOwnedMappings: []string{
			"connector manifest registry",
			"synthetic fixture paths",
			"claim boundary lint",
		},
	}
}

func OperatorAssistantDefinitions() []Definition {
	return []Definition{
		ValidationHealthRefreshDefinition(),
		AlertsCancellationPreviewDefinition(),
		RealtimeQualityBacktestDefinition(),
		ConnectorConformanceReviewDefinition(),
		ValidationHealthRunAllDefinition(),
	}
}

func FindOperatorAssistantDefinition(action string) (Definition, bool) {
	for _, definition := range OperatorAssistantDefinitions() {
		if definition.Action == action {
			return definition, true
		}
	}
	return Definition{}, false
}
