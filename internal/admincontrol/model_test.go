package admincontrol

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestResultJSONContractAndDefaults(t *testing.T) {
	started := time.Date(2026, 5, 14, 1, 2, 3, 456, time.UTC)
	result := NewResult("validation_health.refresh", StatusOK, started, "Validation health summary refreshed.", []string{"Review stale rows."}, nil)

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatalf("decode fields: %v", err)
	}
	for _, want := range []string{"action", "status", "started_at", "completed_at", "summary", "next_actions", "claim_flags"} {
		if _, ok := fields[want]; !ok {
			t.Fatalf("result missing field %q: %s", want, encoded)
		}
	}
	for got := range fields {
		if got != "action" && got != "status" && got != "started_at" && got != "completed_at" && got != "summary" && got != "next_actions" && got != "claim_flags" {
			t.Fatalf("unexpected result field %q in %s", got, encoded)
		}
	}
	if result.ClaimFlags.ExternalEvidenceCreated || result.ClaimFlags.ConsumerStatusesChanged || result.ClaimFlags.ComplianceClaimed || result.ClaimFlags.ProductionReadinessClaimed || result.ClaimFlags.AgencyApprovalClaimed || result.ClaimFlags.ConsumerAcceptanceClaimed || result.ClaimFlags.PublicLaunchClaimed || result.ClaimFlags.HostedSaaSClaimed || result.ClaimFlags.VendorCompatibilityClaimed || result.ClaimFlags.HardwareCertificationClaimed || result.ClaimFlags.SLAClaimed || result.ClaimFlags.UptimeGuaranteeClaimed || result.ClaimFlags.ProductionGradeETAClaimed {
		t.Fatalf("claim flags must default false: %+v", result.ClaimFlags)
	}
}

func TestResultStatusesNormalizeToBoundedVocabulary(t *testing.T) {
	for _, status := range []Status{StatusOK, StatusNeedsReview, StatusBlocked, StatusFailed} {
		result := NewResult("x", status, time.Time{}, "summary", nil, nil)
		if result.Status != status {
			t.Fatalf("status %q normalized to %q", status, result.Status)
		}
	}
	result := NewResult("x", Status("accepted"), time.Time{}, "summary", nil, nil)
	if result.Status != StatusBlocked {
		t.Fatalf("unsupported status = %q, want blocked", result.Status)
	}
}

func TestResultBoundsAndRedactsDiagnostics(t *testing.T) {
	long := strings.Repeat("a", MaxSummaryLength+20)
	result := NewResult(
		"validation_health.refresh",
		StatusNeedsReview,
		time.Time{},
		long,
		[]string{"Review /Users/example/private/path", strings.Repeat("b", MaxNextActionLength+20)},
		[]Error{
			{Code: "raw_report", Message: "stdout contains /tmp/private"},
			{Code: "ok", Message: strings.Repeat("c", MaxErrorTextLength+20)},
			{Code: "extra1", Message: "one"},
			{Code: "extra2", Message: "two"},
			{Code: "extra3", Message: "three"},
		},
	)
	if len([]rune(result.Summary)) != MaxSummaryLength {
		t.Fatalf("summary length = %d, want %d", len([]rune(result.Summary)), MaxSummaryLength)
	}
	if result.NextActions[0] != "redacted private diagnostic" {
		t.Fatalf("next action was not redacted: %+v", result.NextActions)
	}
	if len([]rune(result.NextActions[1])) != MaxNextActionLength {
		t.Fatalf("next action length = %d, want %d", len([]rune(result.NextActions[1])), MaxNextActionLength)
	}
	if len(result.Errors) != MaxErrorCount {
		t.Fatalf("errors = %d, want %d", len(result.Errors), MaxErrorCount)
	}
	if result.Errors[0].Code != "redacted private diagnostic" || result.Errors[0].Message != "redacted private diagnostic" {
		t.Fatalf("private diagnostic error was not redacted: %+v", result.Errors[0])
	}
	if len([]rune(result.Errors[1].Message)) != MaxErrorTextLength {
		t.Fatalf("error message length = %d, want %d", len([]rune(result.Errors[1].Message)), MaxErrorTextLength)
	}
}

func TestValidationHealthDefinitionsBounded(t *testing.T) {
	refresh := ValidationHealthRefreshDefinition()
	if refresh.Action != "validation_health.refresh" || refresh.LadderLevel != LevelReadOnlyRefresh || refresh.RequiredRole != "read_only" {
		t.Fatalf("unexpected refresh definition: %+v", refresh)
	}
	if !strings.Contains(refresh.PublicFeedImpact, "No public feed output changes") || !strings.Contains(strings.ToLower(refresh.PrivateImpact), "writes nothing") {
		t.Fatalf("refresh definition does not state bounded impact: %+v", refresh)
	}

	runAll := ValidationHealthRunAllDefinition()
	if runAll.Action != "validation_health.run_all" || runAll.LadderLevel != LevelReversiblePrivate || runAll.RequiredRole != "admin" {
		t.Fatalf("unexpected run-all definition: %+v", runAll)
	}
	if !strings.Contains(runAll.PrivateImpact, "validation_report") || !strings.Contains(runAll.PublicFeedImpact, "No public feed output changes") {
		t.Fatalf("run-all definition does not state private diagnostic write boundary: %+v", runAll)
	}
}

func TestOperatorAssistantDefinitionsAreServerOwnedAndBounded(t *testing.T) {
	definitions := OperatorAssistantDefinitions()
	wantActions := []string{
		"validation_health.refresh",
		"alerts.cancellation_reconcile.preview",
		"realtime_quality.backtest.dry_run",
		"connectors.conformance.review",
		"validation_health.run_all",
	}
	if len(definitions) != len(wantActions) {
		t.Fatalf("definitions = %d, want %d: %+v", len(definitions), len(wantActions), definitions)
	}
	seen := map[string]bool{}
	for i, definition := range definitions {
		if definition.Action != wantActions[i] {
			t.Fatalf("definition %d action = %q, want %q", i, definition.Action, wantActions[i])
		}
		if seen[definition.Action] {
			t.Fatalf("duplicate action %q", definition.Action)
		}
		seen[definition.Action] = true
		if definition.Label == "" || definition.RequiredRole == "" || definition.Confirmation == "" || definition.PublicFeedImpact == "" || definition.PrivateImpact == "" || definition.RollbackPath == "" || definition.TechnicalHandoff == "" || definition.DoesNotProve == "" {
			t.Fatalf("definition has empty required fields: %+v", definition)
		}
		if !strings.Contains(definition.PublicFeedImpact, "No public feed output changes") {
			t.Fatalf("definition %s does not bound public feed impact: %+v", definition.Action, definition)
		}
		if !strings.Contains(definition.DoesNotProve, "This private view does not show") {
			t.Fatalf("definition %s does not state non-claims: %+v", definition.Action, definition)
		}
		if definition.Action != "validation_health.refresh" && definition.Action != "validation_health.run_all" && !definition.DisabledByDefault {
			t.Fatalf("future operator assistant action %s must remain disabled by default", definition.Action)
		}
		encoded, err := json.Marshal(definition)
		if err != nil {
			t.Fatalf("marshal definition %s: %v", definition.Action, err)
		}
		lower := strings.ToLower(string(encoded))
		for _, forbidden := range []string{
			"raw_command",
			"validator_command",
			"argv",
			"authorization:",
			"bearer ",
			"postgres://",
			"/users/",
			"consumer accepted",
			"compliance achieved",
			"production ready",
			"vendor compatible",
			"certified hardware",
		} {
			if strings.Contains(lower, forbidden) {
				t.Fatalf("definition %s contains forbidden %q: %s", definition.Action, forbidden, encoded)
			}
		}
	}
	if _, ok := FindOperatorAssistantDefinition("alerts.cancellation_reconcile.preview"); !ok {
		t.Fatal("catalog lookup did not find alerts preview definition")
	}
	if _, ok := FindOperatorAssistantDefinition("unknown.action"); ok {
		t.Fatal("catalog lookup found unknown action")
	}
}
