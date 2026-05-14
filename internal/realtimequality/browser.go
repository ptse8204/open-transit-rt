package realtimequality

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	DefaultBacktestBrowserRoot  = ".cache/realtime-quality-backtest"
	defaultBacktestBrowserLimit = 5
)

var safeBacktestRunName = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

type BacktestBrowser struct {
	Status   string               `json:"status"`
	Boundary string               `json:"boundary"`
	RootRef  string               `json:"root_ref"`
	Message  string               `json:"message"`
	Rows     []BacktestBrowserRun `json:"rows"`
}

type BacktestBrowserRun struct {
	OutputRef               string `json:"output_ref"`
	Status                  string `json:"status"`
	GeneratedAt             string `json:"generated_at"`
	ObservedRecords         int    `json:"observed_records"`
	PredictionRecords       int    `json:"prediction_records"`
	GroupCount              int    `json:"group_count"`
	MaturityGate            string `json:"maturity_gate"`
	PredictionCoverage      string `json:"prediction_coverage"`
	FutureStopCoverage      string `json:"future_stop_coverage"`
	MAEAbsoluteErrorSeconds string `json:"mae_absolute_error_seconds"`
	P90AbsoluteErrorSeconds string `json:"p90_absolute_error_seconds"`
	WithheldByReason        string `json:"withheld_by_reason"`
	DiagnosticSignal        string `json:"diagnostic_signal"`
	DoesNotProve            string `json:"does_not_prove"`
}

func LoadBacktestBrowser(root string, limit int) (BacktestBrowser, error) {
	if strings.TrimSpace(root) == "" {
		root = DefaultBacktestBrowserRoot
	}
	if limit <= 0 {
		limit = defaultBacktestBrowserLimit
	}
	if limit > defaultMaxGroups {
		limit = defaultBacktestBrowserLimit
	}
	rootRef, ok := backtestBrowserRootRef(root)
	if !ok {
		return BacktestBrowser{}, fmt.Errorf("backtest browser root must be %s", DefaultBacktestBrowserRoot)
	}
	browser := BacktestBrowser{
		Status:   "missing",
		Boundary: "Private local aggregate backtest summaries only. This reader verifies exact output files, schema versions, aggregate-only manifest flags, and forbidden-claim flags before returning bounded metrics. It does not execute backtests, read source-row inputs, contact predictors, create evidence, alter public feeds, or prove production-grade ETA quality.",
		RootRef:  rootRef,
		Message:  "No local realtime quality backtest summaries are available yet.",
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return BacktestBrowser{}, fmt.Errorf("backtest browser root cannot be resolved")
	}
	if evidenceLikePath(abs) {
		return BacktestBrowser{}, fmt.Errorf("backtest browser root must not be evidence-like")
	}
	info, err := os.Lstat(abs)
	if os.IsNotExist(err) {
		return browser, nil
	}
	if err != nil {
		browser.Status = "blocked"
		browser.Message = "Local backtest summaries could not be listed from the private cache root."
		return browser, nil
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		browser.Status = "blocked"
		browser.Message = "Local backtest summary root failed private path safety checks."
		return browser, nil
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		browser.Status = "blocked"
		browser.Message = "Local backtest summaries could not be listed from the private cache root."
		return browser, nil
	}
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].Name() > entries[j].Name() })
	for _, entry := range entries {
		if len(browser.Rows) >= limit {
			break
		}
		if !entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !safeBacktestRunName.MatchString(entry.Name()) {
			continue
		}
		runRef := rootRef + "/" + entry.Name()
		row := loadBacktestBrowserRun(filepath.Join(abs, entry.Name()), runRef)
		browser.Rows = append(browser.Rows, row)
	}
	if len(browser.Rows) == 0 {
		return browser, nil
	}
	browser.Status = "ok"
	browser.Message = "Local aggregate backtest summaries are available for private review."
	for _, row := range browser.Rows {
		if row.Status == "blocked" {
			browser.Status = "blocked"
			browser.Message = "One or more local backtest summaries failed private safety checks."
			return browser, nil
		}
		if row.Status != "ok" {
			browser.Status = "needs_review"
			browser.Message = "One or more local backtest summaries need operator review."
		}
	}
	return browser, nil
}

func loadBacktestBrowserRun(dir, ref string) BacktestBrowserRun {
	row := BacktestBrowserRun{
		OutputRef:        ref,
		Status:           "blocked",
		GeneratedAt:      "not loaded",
		MaturityGate:     "blocked",
		DiagnosticSignal: "Backtest summary failed private browser safety checks.",
		DoesNotProve:     "Backtest summaries do not prove production-grade ETA quality, real-world ETA accuracy, compliance, consumer acceptance, vendor compatibility, SLA coverage, hosted service readiness, or release readiness.",
	}
	if !backtestOutputFilesAreExact(dir) {
		return row
	}
	var summary SummaryDocument
	if err := readBacktestBrowserJSON(filepath.Join(dir, "summary.json"), &summary); err != nil {
		return row
	}
	var metrics MetricsDocument
	if err := readBacktestBrowserJSON(filepath.Join(dir, "metrics.json"), &metrics); err != nil {
		return row
	}
	var manifest ManifestDocument
	if err := readBacktestBrowserJSON(filepath.Join(dir, "manifest.json"), &manifest); err != nil {
		return row
	}
	if !backtestBrowserDocumentsAreSafe(summary, metrics, manifest) {
		return row
	}
	status := statusForBacktestGate(summary.Overall.MaturityGate)
	row.Status = status
	row.GeneratedAt = summary.GeneratedAt
	row.ObservedRecords = summary.InputRecordCounts.ObservedRecords
	row.PredictionRecords = summary.InputRecordCounts.PredictionRecords
	row.GroupCount = summary.GroupCount
	row.MaturityGate = summary.Overall.MaturityGate
	row.PredictionCoverage = formatBacktestRate(summary.Overall.PredictionCoverage)
	row.FutureStopCoverage = formatBacktestRate(summary.Overall.FutureStopCoverage)
	row.MAEAbsoluteErrorSeconds = formatBacktestSeconds(summary.Overall.MAEAbsoluteErrorSeconds)
	row.P90AbsoluteErrorSeconds = formatBacktestSeconds(summary.Overall.P90AbsoluteErrorSeconds)
	row.WithheldByReason = formatBacktestWithheld(summary.Overall.WithheldByReason)
	row.DiagnosticSignal = backtestDiagnosticSignal(summary.Overall)
	return row
}

func backtestOutputFilesAreExact(dir string) bool {
	expected := ExpectedBacktestOutputFiles()
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != len(expected) {
		return false
	}
	expectedSet := map[string]bool{}
	for _, name := range expected {
		expectedSet[name] = true
		info, err := os.Lstat(filepath.Join(dir, name))
		if err != nil || info.Mode()&os.ModeSymlink != 0 || info.IsDir() {
			return false
		}
	}
	for _, entry := range entries {
		if !expectedSet[entry.Name()] {
			return false
		}
	}
	return true
}

func readBacktestBrowserJSON(path string, value any) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(raw) > defaultMaxOutputBytes {
		return fmt.Errorf("backtest browser file too large")
	}
	return json.Unmarshal(raw, value)
}

func backtestBrowserDocumentsAreSafe(summary SummaryDocument, metrics MetricsDocument, manifest ManifestDocument) bool {
	if summary.SchemaVersion != BacktestSchemaVersion || metrics.SchemaVersion != BacktestSchemaVersion || manifest.SchemaVersion != BacktestSchemaVersion {
		return false
	}
	if summary.GeneratedAt == "" || metrics.GeneratedAt == "" || manifest.GeneratedAt == "" {
		return false
	}
	if summary.ExternalEvidenceCreated || summary.ConsumerStatusesChanged || summary.ComplianceClaimed || summary.ProductionReadinessClaimed || summary.HostedSaaSClaimed || summary.AgencyAdoptionClaimed || summary.ConsumerAcceptanceClaimed || summary.VendorCompatibilityClaimed || summary.ProductionGradeETAClaimed || summary.PublicAPIAdded || summary.DatabasePersistenceAdded {
		return false
	}
	if manifest.OutputKind != "private_local_realtime_quality_backtest" || !manifest.AggregateOnly || manifest.RawRowsPersisted {
		return false
	}
	if !sameStringSet(manifest.OutputFiles, ExpectedBacktestOutputFiles()) {
		return false
	}
	for _, required := range []string{"docs_evidence_output_rejected", "evidence_like_output_rejected", "symlink_ancestors_rejected", "raw_inputs_not_copied", "private_paths_omitted", "raw_rows_omitted"} {
		if !manifest.SafetyChecks[required] {
			return false
		}
	}
	for _, forbidden := range []string{"db_persistence", "migration_added", "operations_console_change", "public_api_added", "consumer_tracker_changed", "external_predictor_runtime_added"} {
		if manifest.Boundaries[forbidden] {
			return false
		}
	}
	if !allowedGateLabels[summary.Overall.MaturityGate] || len(metrics.Groups) > defaultMaxGroups || summary.GroupCount > defaultMaxGroups {
		return false
	}
	return true
}

func sameStringSet(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	seen := map[string]int{}
	for _, value := range left {
		seen[value]++
	}
	for _, value := range right {
		seen[value]--
	}
	for _, count := range seen {
		if count != 0 {
			return false
		}
	}
	return true
}

func statusForBacktestGate(gate string) string {
	switch gate {
	case "diagnostic_pass":
		return "ok"
	case "diagnostic_watch", "insufficient_data":
		return "needs_review"
	default:
		return "blocked"
	}
}

func backtestDiagnosticSignal(overall MetricGroup) string {
	switch overall.MaturityGate {
	case "diagnostic_pass":
		return "Aggregate fixture backtest met local diagnostic thresholds."
	case "diagnostic_watch":
		return "Aggregate fixture backtest stayed within diagnostic review thresholds but needs operator review."
	case "insufficient_data":
		return "Aggregate fixture backtest does not have enough matched samples for a diagnostic signal."
	default:
		return "Backtest summary failed private browser safety checks."
	}
}

func formatBacktestRate(rate Rate) string {
	if rate.Percent == nil {
		return fmt.Sprintf("%s (%d/%d)", firstNonBlank(rate.Status, "not recorded"), rate.Numerator, rate.Denominator)
	}
	return fmt.Sprintf("%.1f%% (%d/%d)", *rate.Percent, rate.Numerator, rate.Denominator)
}

func formatBacktestSeconds(value *float64) string {
	if value == nil {
		return "not recorded"
	}
	return fmt.Sprintf("%.0f sec", *value)
}

func formatBacktestWithheld(counts map[string]int) string {
	if len(counts) == 0 {
		return "none recorded"
	}
	keys := make([]string, 0, len(counts))
	for key, count := range counts {
		if count > 0 {
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return "none recorded"
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", key, counts[key]))
	}
	return strings.Join(parts, ", ")
}

func backtestBrowserRootRef(root string) (string, bool) {
	clean := filepath.ToSlash(filepath.Clean(root))
	if clean == DefaultBacktestBrowserRoot || strings.HasSuffix(clean, "/"+DefaultBacktestBrowserRoot) {
		return DefaultBacktestBrowserRoot, true
	}
	return "", false
}

func evidenceLikePath(path string) bool {
	lower := strings.ToLower(filepath.ToSlash(path))
	return strings.Contains(lower, "docs/evidence") || strings.Contains(lower, "evidence")
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
