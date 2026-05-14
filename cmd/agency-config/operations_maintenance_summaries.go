package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const maintenanceSummarySizeLimit = 1 << 20

var (
	maintenanceDeploymentDoctorRoot      = filepath.Join(".cache", "deployment-doctor")
	maintenanceOperationsReliabilityRoot = filepath.Join(".cache", "operations-reliability")
	maintenanceOperationsNotifyRoot      = filepath.Join(".cache", "operations-notify")
	maintenanceSupportBundleRoot         = filepath.Join(".cache", "support-bundles")
	maintenanceSummaryRunName            = regexp.MustCompile(`^\d{8}T\d{6}Z$`)
)

type maintenanceSummarySource struct {
	ID              string
	Label           string
	Root            string
	Kind            string
	FileName        string
	MissingAction   string
	DoesNotProve    string
	BuildSignal     func(map[string]any) (string, string, string)
	ClaimFlagsValid func(map[string]any) bool
}

func buildOperationsMaintenanceDiagnostics() operationsMaintenanceDiagnostics {
	sources := []maintenanceSummarySource{
		{
			ID:            "deployment_doctor",
			Label:         "Deployment doctor",
			Root:          maintenanceDeploymentDoctorRoot,
			Kind:          "deployment-doctor",
			FileName:      "summary.json",
			MissingAction: "Run `make deployment-doctor` from an operator shell when deployment diagnostics are needed.",
			DoesNotProve:  "Deployment doctor summaries do not prove production readiness, hosted service availability, SLA coverage, uptime, compliance, agency adoption, consumer acceptance, vendor compatibility, or disaster recovery.",
			BuildSignal: func(data map[string]any) (string, string, string) {
				overall := maintenanceString(data["overall_status"], "unknown")
				counts := maintenanceMap(data["counts"])
				categories := maintenanceMap(data["categories"])
				signal := fmt.Sprintf("overall=%s; blockers=%s; warnings=%s; backup=%s; restore=%s",
					overall,
					maintenanceNumberText(counts["blocker"]),
					maintenanceNumberText(counts["warning"]),
					maintenanceString(categories["backup_readiness"], "unknown"),
					maintenanceString(categories["restore_readiness"], "unknown"),
				)
				return maintenanceStatusFromDiagnostic(overall), signal, "Review blocker/warning counts and keep missing backup or restore values visible until a technical helper configures them."
			},
			ClaimFlagsValid: func(data map[string]any) bool {
				return maintenanceFalseFlags(data, "external_evidence_created", "final_root_evidence_created", "consumer_statuses_changed", "compliance_claimed", "production_readiness_claimed")
			},
		},
		{
			ID:            "operations_reliability",
			Label:         "Operations reliability",
			Root:          maintenanceOperationsReliabilityRoot,
			Kind:          "operations-reliability",
			FileName:      "summary.json",
			MissingAction: "Run `make operations-reliability` from an operator shell after validator-health, deployment-doctor, and notification summaries exist.",
			DoesNotProve:  "Operations reliability summaries do not prove SLA coverage, uptime guarantees, hosted service readiness, production readiness, compliance, consumer acceptance, or public launch.",
			BuildSignal: func(data map[string]any) (string, string, string) {
				overall := maintenanceString(data["overall_status"], "unknown")
				backup := maintenanceMap(data["backup_restore"])
				alerting := maintenanceMap(data["alerting"])
				availability := maintenanceMap(data["availability_sampling"])
				signal := fmt.Sprintf("overall=%s; backup_restore=%s; alerting=%s; availability=%s",
					overall,
					maintenanceString(backup["status"], "unknown"),
					maintenanceString(alerting["status"], "unknown"),
					maintenanceString(availability["status"], "unknown"),
				)
				return maintenanceStatusFromDiagnostic(overall), signal, "Review missing or needs-review sections before treating operations as routine."
			},
			ClaimFlagsValid: func(data map[string]any) bool {
				return maintenanceFalseFlags(maintenanceMap(data["claim_flags"]), "external_evidence_created", "final_root_evidence_created", "consumer_statuses_changed", "compliance_claimed", "production_readiness_claimed", "hosted_saas_claimed", "consumer_acceptance_claimed", "vendor_compatibility_claimed", "sla_claimed", "uptime_guarantee_claimed", "production_grade_eta_claimed")
			},
		},
		{
			ID:            "operations_notify",
			Label:         "Operations notification draft",
			Root:          maintenanceOperationsNotifyRoot,
			Kind:          "operations-notify",
			FileName:      "summary.json",
			MissingAction: "Run `make operations-notify` from an operator shell when a private local notification draft is needed.",
			DoesNotProve:  "Notification drafts do not prove any notification was sent, received, acknowledged, or tied to SLA/uptime coverage.",
			BuildSignal: func(data map[string]any) (string, string, string) {
				notification := maintenanceMap(data["notification"])
				counts := maintenanceMap(data["counts"])
				severity := maintenanceString(notification["severity"], maintenanceString(data["severity"], "unknown"))
				sent := maintenanceString(notification["not_sent"], maintenanceString(data["notification_sent"], "unknown"))
				signal := fmt.Sprintf("severity=%s; not_sent=%s; next_actions=%s; blocked_actions=%s",
					severity,
					sent,
					maintenanceNumberText(counts["next_actions"]),
					maintenanceNumberText(counts["blocked_actions"]),
				)
				return maintenanceStatusFromSeverity(severity), signal, "Review the private draft and next actions; send nothing unless a separate operator process authorizes it."
			},
			ClaimFlagsValid: func(data map[string]any) bool {
				return maintenanceFalseFlags(data, "external_evidence_created", "consumer_statuses_changed", "compliance_claimed", "production_readiness_claimed", "hosted_saas_claimed", "agency_adoption_claimed", "consumer_acceptance_claimed", "vendor_compatibility_claimed", "production_grade_eta_claimed", "notification_sent")
			},
		},
		{
			ID:            "support_bundle_manifest",
			Label:         "Support bundle manifest",
			Root:          maintenanceSupportBundleRoot,
			Kind:          "support-bundles",
			FileName:      "manifest.json",
			MissingAction: "Run `make support-bundle` from an operator shell only when a technical helper needs redaction-safe diagnostics.",
			DoesNotProve:  "Support bundle manifests do not prove a support bundle is safe to share outside the operator environment or qualify as retained evidence.",
			BuildSignal: func(data map[string]any) (string, string, string) {
				included := maintenanceListLength(data["included"])
				excluded := maintenanceListLength(data["excluded"])
				signal := fmt.Sprintf("included_categories=%d; excluded_categories=%d; external_evidence_created=%s; consumer_statuses_changed=%s",
					included,
					excluded,
					maintenanceString(data["external_evidence_created"], "false"),
					maintenanceString(data["consumer_statuses_changed"], "false"),
				)
				return operationsStatusDiagnosticOnly, signal, "Review redaction warnings before sharing any generated support bundle outside the operator environment."
			},
			ClaimFlagsValid: func(data map[string]any) bool {
				return maintenanceFalseFlags(data, "external_evidence_created", "consumer_statuses_changed")
			},
		},
	}
	rows := make([]operationsMaintenanceDiagnostic, 0, len(sources))
	for _, source := range sources {
		rows = append(rows, loadMaintenanceDiagnostic(source))
	}
	return operationsMaintenanceDiagnostics{
		Status:   maintenanceDiagnosticsOverall(rows),
		Boundary: "Local diagnostic summaries are read-only pointers into ignored `.cache` outputs. This page reads only bounded summary or manifest files from known helper roots and never reads raw logs, command output, backup dumps, private payloads, or retained evidence.",
		Rows:     rows,
	}
}

func loadMaintenanceDiagnostic(source maintenanceSummarySource) operationsMaintenanceDiagnostic {
	row := operationsMaintenanceDiagnostic{
		ID:            source.ID,
		Label:         source.Label,
		Status:        operationsStatusMissing,
		SourceRef:     maintenanceSourceRootRef(source.Kind),
		GeneratedAt:   "not available",
		CurrentSignal: "no safe local summary found",
		NextAction:    source.MissingAction,
		DoesNotProve:  source.DoesNotProve,
	}
	rootRef, ok := maintenanceSummaryRootRef(source.Root, source.Kind)
	if !ok {
		row.Status = operationsStatusBlocked
		row.CurrentSignal = "configured summary root failed private cache boundary checks"
		row.NextAction = "Reset the summary root to the expected private `.cache` helper location."
		return row
	}
	path, runName, ok, blocked := latestMaintenanceSummaryPath(source.Root, source.Kind, source.FileName)
	if blocked {
		row.Status = operationsStatusBlocked
		row.SourceRef = rootRef
		row.CurrentSignal = "summary root failed private path safety checks"
		row.NextAction = "Regenerate local diagnostics under the expected non-symlink `.cache` root."
		return row
	}
	if !ok {
		row.SourceRef = rootRef
		return row
	}
	row.SourceRef = rootRef + "/" + runName + "/" + source.FileName
	data, err := readMaintenanceSummaryJSON(path)
	if err != nil {
		row.Status = operationsStatusBlocked
		row.CurrentSignal = "summary failed bounded JSON safety checks"
		row.NextAction = "Regenerate the local diagnostic summary before browser review."
		return row
	}
	if !source.ClaimFlagsValid(data) {
		row.Status = operationsStatusBlocked
		row.CurrentSignal = "summary contains a claim flag that must remain false"
		row.NextAction = "Do not use this summary; regenerate it after fixing the claim boundary."
		return row
	}
	status, signal, next := source.BuildSignal(data)
	row.Status = status
	row.GeneratedAt = firstNonEmpty(maintenanceString(data["generated_at_utc"], ""), maintenanceString(data["generated_at"], ""), maintenanceString(data["created_at_utc"], ""), maintenanceString(data["created_at"], ""), runName)
	row.CurrentSignal = maintenanceBoundedText(signal, "summary loaded")
	row.NextAction = maintenanceBoundedText(next, source.MissingAction)
	return row
}

func buildOperationsMaintenanceInfrastructureChecks() operationsMaintenancePanel {
	rows := []operationsMaintenancePanelRow{
		maintenancePanelRow(
			"database_connectivity",
			"Database connectivity",
			operationsStatusMissing,
			"no safe deployment-doctor summary found",
			"Keep database status marked missing until a technical helper runs private deployment diagnostics.",
			"Run `make deployment-doctor` from an operator shell; values and raw migrator output stay outside the browser.",
			"Database diagnostics do not prove production readiness or data-loss safety.",
		),
		maintenancePanelRow(
			"migration_status",
			"Migration status",
			operationsStatusMissing,
			"no safe deployment-doctor summary found",
			"Do not treat missing migration status as safe for upgrade or rollback.",
			"Use the deployment doctor or migrator status command from an operator shell before schema work.",
			"Migration diagnostics do not prove backward compatibility by themselves.",
		),
		maintenancePanelRow(
			"postgis_extension",
			"PostGIS extension",
			operationsStatusMissing,
			"no safe deployment-doctor summary found",
			"Keep spatial database capability marked missing until the private diagnostic summary records it.",
			"Run the deployment doctor against the intended private database target when configured.",
			"PostGIS status does not prove all geospatial queries are production safe.",
		),
		maintenancePanelRow(
			"validator_tooling",
			"Validator tooling",
			operationsStatusMissing,
			"no safe deployment-doctor summary found",
			"Treat missing validator tooling as a reason to use off-host validation guidance.",
			"Run `make deployment-doctor` or `./scripts/check-validators.sh` from an operator shell when a fresh tooling check is needed.",
			"Validator tooling presence does not prove feeds are validator clean or compliant.",
		),
		maintenancePanelRow(
			"backup_storage_access",
			"Backup storage access",
			operationsStatusMissing,
			"no safe deployment-doctor summary found",
			"Keep backup storage status marked missing until a private backup target is configured and checked.",
			"Run deployment diagnostics from an operator shell; do not expose backup paths, dumps, or raw filesystem output in the browser.",
			"Backup storage access does not prove a backup exists or restore will succeed.",
		),
	}
	rootRef, ok := maintenanceSummaryRootRef(maintenanceDeploymentDoctorRoot, "deployment-doctor")
	if !ok {
		return maintenanceInfrastructurePanel(rows, operationsStatusBlocked, "configured deployment-doctor root failed private cache boundary checks")
	}
	path, runName, ok, blocked := latestMaintenanceSummaryPath(maintenanceDeploymentDoctorRoot, "deployment-doctor", "summary.json")
	if blocked {
		return maintenanceInfrastructurePanel(rows, operationsStatusBlocked, "deployment-doctor summary root failed private path safety checks")
	}
	if !ok {
		return maintenanceInfrastructurePanel(rows, "", "")
	}
	source := rootRef + "/" + runName + "/summary.json"
	data, err := readMaintenanceSummaryJSON(path)
	if err != nil {
		return maintenanceInfrastructurePanel(rows, operationsStatusBlocked, "deployment-doctor summary failed bounded JSON safety checks")
	}
	if !maintenanceFalseFlags(data, "external_evidence_created", "final_root_evidence_created", "consumer_statuses_changed", "compliance_claimed", "production_readiness_claimed") {
		return maintenanceInfrastructurePanel(rows, operationsStatusBlocked, "deployment-doctor summary contains a claim flag that must remain false")
	}
	categories := maintenanceMap(data["categories"])
	updates := []struct {
		id       string
		category string
	}{
		{id: "database_connectivity", category: "database"},
		{id: "migration_status", category: "migrations"},
		{id: "postgis_extension", category: "postgis"},
		{id: "validator_tooling", category: "validators"},
		{id: "backup_storage_access", category: "backup_readiness"},
	}
	for i := range rows {
		for _, update := range updates {
			if rows[i].ID != update.id {
				continue
			}
			categoryStatus := maintenanceString(categories[update.category], "unknown")
			rows[i].Status = maintenanceStatusFromDiagnostic(categoryStatus)
			rows[i].CurrentSignal = maintenanceBoundedText(fmt.Sprintf("%s=%s; source=%s", update.category, categoryStatus, source), "deployment-doctor category loaded")
		}
	}
	return maintenanceInfrastructurePanel(rows, "", "")
}

func maintenanceInfrastructurePanel(rows []operationsMaintenancePanelRow, forcedStatus, forcedSignal string) operationsMaintenancePanel {
	if forcedStatus != "" || forcedSignal != "" {
		for i := range rows {
			if forcedStatus != "" {
				rows[i].Status = forcedStatus
			}
			if forcedSignal != "" {
				rows[i].CurrentSignal = forcedSignal
			}
		}
	}
	return operationsMaintenancePanel{
		Status:     maintenancePanelOverall(rows),
		Boundary:   "Infrastructure checks are read-only deployment-doctor category summaries. The browser does not connect to databases, inspect private paths, run validators, run migrations, or execute disk checks.",
		NextAction: "Run deployment diagnostics from an operator shell when fresh database, migration, PostGIS, validator, or backup-storage status is needed.",
		Rows:       rows,
	}
}

func latestMaintenanceSummaryPath(root, kind, fileName string) (string, string, bool, bool) {
	abs, err := filepath.Abs(root)
	if err != nil || maintenanceEvidenceLikePath(abs) {
		return "", "", false, true
	}
	info, err := os.Lstat(abs)
	if os.IsNotExist(err) {
		return "", "", false, false
	}
	if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return "", "", false, true
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return "", "", false, true
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !maintenanceSummaryRunName.MatchString(entry.Name()) {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	for _, name := range names {
		path := filepath.Join(abs, name, fileName)
		info, err := os.Lstat(path)
		if err != nil || info.Mode()&os.ModeSymlink != 0 || info.IsDir() {
			continue
		}
		return path, name, true, false
	}
	return "", "", false, false
}

func readMaintenanceSummaryJSON(path string) (map[string]any, error) {
	info, err := os.Lstat(path)
	if err != nil || info.Mode()&os.ModeSymlink != 0 || info.IsDir() || info.Size() <= 0 || info.Size() > maintenanceSummarySizeLimit {
		return nil, fmt.Errorf("unsafe maintenance summary")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if maintenancePrivateText(string(raw)) {
		return nil, fmt.Errorf("private maintenance summary text")
	}
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty maintenance summary")
	}
	return data, nil
}

func maintenanceDiagnosticsOverall(rows []operationsMaintenanceDiagnostic) string {
	status := operationsStatusReady
	for _, row := range rows {
		status = worseMaintenanceStatus(status, row.Status)
	}
	return status
}

func maintenanceSummaryRootRef(root, kind string) (string, bool) {
	clean := filepath.ToSlash(filepath.Clean(root))
	want := maintenanceSourceRootRef(kind)
	if clean == want || strings.HasSuffix(clean, "/"+want) {
		return want, true
	}
	return "", false
}

func maintenanceSourceRootRef(kind string) string {
	return ".cache/" + kind
}

func maintenanceStatusFromDiagnostic(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "passed", "pass", "ok", "ready", "recorded":
		return operationsStatusReady
	case "blocker", "blocked", "unhealthy", "failed", "error":
		return operationsStatusBlocked
	case "missing", "not_found":
		return operationsStatusMissing
	case "warning", "warnings", "unavailable", "skipped", "needs_review", "stale", "degraded":
		return operationsStatusNeedsReview
	case "diagnostic_only":
		return operationsStatusDiagnosticOnly
	default:
		return operationsStatusUnknown
	}
}

func maintenanceStatusFromSeverity(severity string) string {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "info", "ok", "ready":
		return operationsStatusReady
	case "blocked", "blocker":
		return operationsStatusBlocked
	case "needs_review", "warning":
		return operationsStatusNeedsReview
	case "missing", "missing_source":
		return operationsStatusMissing
	default:
		return operationsStatusUnknown
	}
}

func maintenanceFalseFlags(data map[string]any, keys ...string) bool {
	for _, key := range keys {
		value, ok := data[key]
		if !ok {
			continue
		}
		if truthyMaintenanceValue(value) {
			return false
		}
	}
	return true
}

func truthyMaintenanceValue(value any) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		switch strings.ToLower(strings.TrimSpace(typed)) {
		case "", "false", "0", "no", "not_sent":
			return false
		default:
			return true
		}
	case float64:
		return typed != 0
	default:
		return false
	}
}

func maintenanceMap(value any) map[string]any {
	if typed, ok := value.(map[string]any); ok {
		return typed
	}
	return map[string]any{}
}

func maintenanceString(value any, fallback string) string {
	switch typed := value.(type) {
	case string:
		return maintenanceBoundedText(typed, fallback)
	case bool:
		if typed {
			return "true"
		}
		return "false"
	case float64:
		return fmt.Sprintf("%.0f", typed)
	default:
		return fallback
	}
}

func maintenanceNumberText(value any) string {
	switch typed := value.(type) {
	case float64:
		return fmt.Sprintf("%.0f", typed)
	case int:
		return fmt.Sprintf("%d", typed)
	default:
		return "0"
	}
}

func maintenanceListLength(value any) int {
	if typed, ok := value.([]any); ok {
		return len(typed)
	}
	return 0
}

func maintenanceBoundedText(value, fallback string) string {
	cleaned := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(value, "\n", " "), "\r", " "))
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	if cleaned == "" || maintenancePrivateText(cleaned) {
		return fallback
	}
	const limit = 240
	if len(cleaned) > limit {
		return cleaned[:limit-15] + " [truncated]"
	}
	return cleaned
}

func maintenanceEvidenceLikePath(path string) bool {
	lower := strings.ToLower(filepath.ToSlash(path))
	return strings.Contains(lower, "docs/evidence") || strings.Contains(lower, "/evidence/") || strings.HasSuffix(lower, "/evidence")
}

func maintenancePrivateText(text string) bool {
	lower := strings.ToLower(text)
	for _, forbidden := range []string{
		"authorization:",
		"bearer ",
		"set-cookie",
		"database_url=",
		"restore_database_url=",
		"postgres://",
		"postgresql://",
		"payload_json",
		"file://",
		"/users/",
		"/var/lib",
		"/etc/",
	} {
		if strings.Contains(lower, forbidden) {
			return true
		}
	}
	return false
}
