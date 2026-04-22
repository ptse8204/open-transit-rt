package compliance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultValidatorTimeout   = 2 * time.Minute
	DefaultValidatorOutputCap = 512 * 1024
	DefaultValidatorReportCap = 2 * 1024 * 1024
)

type ValidationStore interface {
	StoreValidationResult(ctx context.Context, result ValidationResult) error
}

type ValidatorSpec struct {
	ID               string
	Name             string
	Version          string
	FeedTypes        []string
	Binary           string
	Args             []string
	Timeout          time.Duration
	MaxStdoutBytes   int64
	MaxStderrBytes   int64
	MaxReportBytes   int64
	RequiresSchedule bool
	RequiresRealtime bool
}

type ValidatorRegistry map[string]ValidatorSpec

func ValidatorRegistryFromEnv() ValidatorRegistry {
	registry := ValidatorRegistry{}
	if spec := staticValidatorFromEnv(); spec.ID != "" {
		registry[spec.ID] = spec
	}
	if spec := realtimeValidatorFromEnv(); spec.ID != "" {
		registry[spec.ID] = spec
	}
	return registry
}

func staticValidatorFromEnv() ValidatorSpec {
	path := strings.TrimSpace(os.Getenv("GTFS_VALIDATOR_PATH"))
	spec := ValidatorSpec{
		ID:               "static-mobilitydata",
		Name:             "mobilitydata-gtfs-validator",
		Version:          getenv("GTFS_VALIDATOR_VERSION", "v7.1.0"),
		FeedTypes:        []string{"schedule"},
		Timeout:          durationFromEnv("GTFS_VALIDATOR_TIMEOUT", DefaultValidatorTimeout),
		MaxStdoutBytes:   int64FromEnv("GTFS_VALIDATOR_STDOUT_MAX_BYTES", DefaultValidatorOutputCap),
		MaxStderrBytes:   int64FromEnv("GTFS_VALIDATOR_STDERR_MAX_BYTES", DefaultValidatorOutputCap),
		MaxReportBytes:   int64FromEnv("GTFS_VALIDATOR_REPORT_MAX_BYTES", DefaultValidatorReportCap),
		RequiresSchedule: true,
	}
	if path == "" {
		return spec
	}
	if strings.HasSuffix(strings.ToLower(path), ".jar") {
		spec.Binary = "java"
		spec.Args = []string{"-jar", path, "-i", "{schedule_zip}", "-o", "{output_dir}"}
		return spec
	}
	spec.Binary = path
	spec.Args = []string{"-i", "{schedule_zip}", "-o", "{output_dir}"}
	return spec
}

func realtimeValidatorFromEnv() ValidatorSpec {
	path := strings.TrimSpace(os.Getenv("GTFS_RT_VALIDATOR_PATH"))
	spec := ValidatorSpec{
		ID:               "realtime-mobilitydata",
		Name:             "mobilitydata-gtfs-realtime-validator",
		Version:          getenv("GTFS_RT_VALIDATOR_VERSION", "pinned-digest-required"),
		FeedTypes:        []string{"vehicle_positions", "trip_updates", "alerts"},
		Timeout:          durationFromEnv("GTFS_RT_VALIDATOR_TIMEOUT", DefaultValidatorTimeout),
		MaxStdoutBytes:   int64FromEnv("GTFS_RT_VALIDATOR_STDOUT_MAX_BYTES", DefaultValidatorOutputCap),
		MaxStderrBytes:   int64FromEnv("GTFS_RT_VALIDATOR_STDERR_MAX_BYTES", DefaultValidatorOutputCap),
		MaxReportBytes:   int64FromEnv("GTFS_RT_VALIDATOR_REPORT_MAX_BYTES", DefaultValidatorReportCap),
		RequiresSchedule: true,
		RequiresRealtime: true,
	}
	if path == "" {
		return spec
	}
	spec.Binary = path
	args := strings.Fields(os.Getenv("GTFS_RT_VALIDATOR_ARGS"))
	if len(args) == 0 {
		args = []string{"--schedule", "{schedule_zip}", "--realtime", "{realtime_pb}", "--feed_type", "{feed_type}", "--output_dir", "{output_dir}"}
	}
	spec.Args = args
	return spec
}

func RunValidation(ctx context.Context, store ValidationStore, registry ValidatorRegistry, input ValidationRunInput) (ValidationResult, error) {
	spec, ok := registry[input.ValidatorID]
	if !ok || input.ValidatorID == "" {
		return ValidationResult{}, fmt.Errorf("unknown validator_id")
	}
	result := ValidationResult{
		AgencyID:         input.AgencyID,
		FeedType:         input.FeedType,
		FeedVersionID:    input.FeedVersionID,
		ValidatorName:    spec.Name,
		ValidatorVersion: spec.Version,
		Status:           "not_run",
		Report:           map[string]any{"validator_id": input.ValidatorID},
	}
	if !spec.supports(input.FeedType) {
		return ValidationResult{}, fmt.Errorf("validator %s does not support feed_type %s", input.ValidatorID, input.FeedType)
	}
	if spec.Binary == "" {
		result.Report["reason"] = "validator_binary_missing"
		if err := store.StoreValidationResult(ctx, result); err != nil {
			return ValidationResult{}, err
		}
		return result, nil
	}
	if err := spec.validatePlaceholders(); err != nil {
		result.Report["reason"] = "validator_args_misconfigured"
		result.Report["error"] = err.Error()
		if err := store.StoreValidationResult(ctx, result); err != nil {
			return ValidationResult{}, err
		}
		return result, nil
	}
	if spec.Timeout <= 0 {
		spec.Timeout = DefaultValidatorTimeout
	}
	if spec.MaxStdoutBytes <= 0 {
		spec.MaxStdoutBytes = DefaultValidatorOutputCap
	}
	if spec.MaxStderrBytes <= 0 {
		spec.MaxStderrBytes = DefaultValidatorOutputCap
	}
	if spec.MaxReportBytes <= 0 {
		spec.MaxReportBytes = DefaultValidatorReportCap
	}

	workDir, err := os.MkdirTemp("", "open-transit-rt-validator-*")
	if err != nil {
		return ValidationResult{}, fmt.Errorf("create validator work dir: %w", err)
	}
	defer os.RemoveAll(workDir)
	outputDir := filepath.Join(workDir, "output")
	if err := os.Mkdir(outputDir, 0o700); err != nil {
		return ValidationResult{}, fmt.Errorf("create validator output dir: %w", err)
	}
	artifacts := map[string]string{"{output_dir}": outputDir, "{feed_type}": input.FeedType}
	if spec.RequiresSchedule {
		if len(input.ScheduleZIPPayload) == 0 {
			return ValidationResult{}, fmt.Errorf("schedule validator artifact is required")
		}
		path := filepath.Join(workDir, "schedule.zip")
		if err := os.WriteFile(path, input.ScheduleZIPPayload, 0o600); err != nil {
			return ValidationResult{}, fmt.Errorf("write schedule validator artifact: %w", err)
		}
		artifacts["{schedule_zip}"] = path
		result.Report["schedule_artifact_source"] = "internal_builder"
	}
	if spec.RequiresRealtime {
		if len(input.RealtimePBPayload) == 0 {
			result.Report["reason"] = "realtime_artifact_unavailable"
			if err := store.StoreValidationResult(ctx, result); err != nil {
				return ValidationResult{}, err
			}
			return result, nil
		}
		path := filepath.Join(workDir, input.FeedType+".pb")
		if err := os.WriteFile(path, input.RealtimePBPayload, 0o600); err != nil {
			return ValidationResult{}, fmt.Errorf("write realtime validator artifact: %w", err)
		}
		artifacts["{realtime_pb}"] = path
		if input.RealtimeArtifactSource != "" {
			result.Report["realtime_artifact_source"] = input.RealtimeArtifactSource
		}
	}

	args := expandArgs(spec.Args, artifacts)
	runCtx, cancel := context.WithTimeout(ctx, spec.Timeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, spec.Binary, args...)
	cmd.Dir = workDir
	stdout := &limitedBuffer{max: spec.MaxStdoutBytes}
	stderr := &limitedBuffer{max: spec.MaxStderrBytes}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	runErr := cmd.Run()

	normalized := normalizeValidatorOutput(stdout.String(), stderr.String(), outputDir, spec.MaxReportBytes)
	result.ErrorCount = normalized.errorCount
	result.WarningCount = normalized.warningCount
	result.InfoCount = normalized.infoCount
	result.Report = normalized.report
	result.Report["validator_id"] = input.ValidatorID
	result.Report["argv"] = redactedArgv(spec.Binary, args, workDir)
	result.Report["stdout"] = stdout.String()
	result.Report["stderr"] = stderr.String()
	if stdout.truncated {
		result.Report["stdout_truncated"] = true
	}
	if stderr.truncated {
		result.Report["stderr_truncated"] = true
	}
	if normalized.source != "" {
		result.Report["parsed_report_source"] = normalized.source
	}
	if runCtx.Err() == context.DeadlineExceeded {
		result.Status = "failed"
		result.ErrorCount = max(1, result.ErrorCount)
		result.Report["error"] = "validator timeout"
	} else if runErr != nil {
		result.Status = "failed"
		result.ErrorCount = max(1, result.ErrorCount)
		result.Report["error"] = runErr.Error()
	} else {
		result.Status = normalized.status()
	}
	if err := store.StoreValidationResult(ctx, result); err != nil {
		return ValidationResult{}, err
	}
	return result, nil
}

func (s ValidatorSpec) supports(feedType string) bool {
	for _, candidate := range s.FeedTypes {
		if candidate == feedType {
			return true
		}
	}
	return false
}

func (s ValidatorSpec) Supports(feedType string) bool {
	return s.supports(feedType)
}

func (s ValidatorSpec) validatePlaceholders() error {
	joined := strings.Join(s.Args, "\x00")
	if !strings.Contains(joined, "{output_dir}") {
		return fmt.Errorf("validator args must include {output_dir}")
	}
	if s.RequiresSchedule && !strings.Contains(joined, "{schedule_zip}") {
		return fmt.Errorf("validator args must include {schedule_zip}")
	}
	if s.RequiresRealtime && !strings.Contains(joined, "{realtime_pb}") {
		return fmt.Errorf("validator args must include {realtime_pb}")
	}
	return nil
}

func expandArgs(args []string, replacements map[string]string) []string {
	expanded := make([]string, len(args))
	for i, arg := range args {
		expanded[i] = arg
		for old, next := range replacements {
			expanded[i] = strings.ReplaceAll(expanded[i], old, next)
		}
	}
	return expanded
}

func redactedArgv(binary string, args []string, workDir string) []string {
	argv := append([]string{binary}, args...)
	for i := range argv {
		if strings.HasPrefix(argv[i], workDir) {
			argv[i] = filepath.Join("{validator_work_dir}", strings.TrimPrefix(argv[i], workDir))
		}
	}
	return argv
}

type limitedBuffer struct {
	buf       bytes.Buffer
	max       int64
	written   int64
	truncated bool
}

func (b *limitedBuffer) Write(payload []byte) (int, error) {
	b.written += int64(len(payload))
	remaining := b.max - int64(b.buf.Len())
	if remaining <= 0 {
		b.truncated = true
		return len(payload), nil
	}
	if int64(len(payload)) > remaining {
		_, _ = b.buf.Write(payload[:remaining])
		b.truncated = true
		return len(payload), nil
	}
	_, _ = b.buf.Write(payload)
	return len(payload), nil
}

func (b *limitedBuffer) String() string {
	return b.buf.String()
}

type normalizedValidationOutput struct {
	errorCount   int
	warningCount int
	infoCount    int
	explicit     string
	source       string
	report       map[string]any
}

func (n normalizedValidationOutput) status() string {
	if n.errorCount > 0 || n.explicit == "failed" {
		return "failed"
	}
	if n.warningCount > 0 || n.explicit == "warning" {
		return "warning"
	}
	if n.explicit == "not_run" {
		return "not_run"
	}
	return "passed"
}

func normalizeValidatorOutput(stdout string, stderr string, outputDir string, maxReportBytes int64) normalizedValidationOutput {
	result := normalizedValidationOutput{report: map[string]any{}}
	if raw, source, ok := firstJSONReport(stdout, stderr, outputDir, maxReportBytes); ok {
		result.source = source
		result.report["raw_report"] = raw
		result.errorCount, result.warningCount, result.infoCount = countsFromJSON(raw)
		result.explicit = statusFromJSON(raw)
		return result
	}
	combined := stdout + "\n" + stderr
	result.errorCount = textCount(combined, `(?i)\b(errors?|fatal)\b\s*[:=]\s*(\d+)`)
	result.warningCount = textCount(combined, `(?i)\b(warnings?)\b\s*[:=]\s*(\d+)`)
	result.infoCount = textCount(combined, `(?i)\b(info|infos|information|notices?)\b\s*[:=]\s*(\d+)`)
	result.report["raw_text"] = strings.TrimSpace(combined)
	if result.errorCount == 0 && result.warningCount == 0 && result.infoCount == 0 {
		result.source = "exit_status"
	} else {
		result.source = "text_counts"
	}
	return result
}

func firstJSONReport(stdout string, stderr string, outputDir string, maxReportBytes int64) (any, string, bool) {
	for _, candidate := range []struct {
		source string
		text   string
	}{
		{source: "stdout", text: stdout},
		{source: "stderr", text: stderr},
	} {
		if raw, ok := parseJSON(candidate.text); ok {
			return raw, candidate.source, true
		}
	}
	files, _ := filepath.Glob(filepath.Join(outputDir, "*.json"))
	sortReportFiles(files)
	for _, file := range files {
		payload, err := readBoundedFile(file, maxReportBytes)
		if err != nil {
			continue
		}
		if raw, ok := parseJSON(string(payload)); ok {
			return raw, filepath.Base(file), true
		}
	}
	return nil, "", false
}

func readBoundedFile(path string, maxBytes int64) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(io.LimitReader(file, maxBytes+1))
}

func parseJSON(text string) (any, bool) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return nil, false
	}
	var raw any
	if err := json.Unmarshal([]byte(trimmed), &raw); err == nil {
		return raw, true
	}
	start := strings.IndexAny(trimmed, "{[")
	if start < 0 {
		return nil, false
	}
	for end := len(trimmed); end > start; end-- {
		var candidate any
		if err := json.Unmarshal([]byte(trimmed[start:end]), &candidate); err == nil {
			return candidate, true
		}
	}
	return nil, false
}

func countsFromJSON(raw any) (int, int, int) {
	errors, warnings, infos := countFromValue(raw)
	if errors > 0 || warnings > 0 || infos > 0 {
		return errors, warnings, infos
	}
	if object, ok := raw.(map[string]any); ok {
		errors = numberField(object, "error_count", "errors_count", "errorCount", "errorsCount", "num_errors", "numErrors", "errors")
		warnings = numberField(object, "warning_count", "warnings_count", "warningCount", "warningsCount", "num_warnings", "numWarnings", "warnings")
		infos = numberField(object, "info_count", "infos_count", "infoCount", "infosCount", "notice_count", "notices_count", "num_infos", "numInfos", "infos", "info", "notices")
	}
	return errors, warnings, infos
}

func countFromValue(value any) (int, int, int) {
	switch typed := value.(type) {
	case map[string]any:
		errors := numberField(typed, "error_count", "errors_count", "errorCount", "errorsCount", "num_errors", "numErrors", "errors")
		warnings := numberField(typed, "warning_count", "warnings_count", "warningCount", "warningsCount", "num_warnings", "numWarnings", "warnings")
		infos := numberField(typed, "info_count", "infos_count", "infoCount", "infosCount", "notice_count", "notices_count", "num_infos", "numInfos", "infos", "info", "notices")
		errors += arrayLen(typed, "errors", "error")
		warnings += arrayLen(typed, "warnings", "warning")
		infos += arrayLen(typed, "infos", "info", "notices")
		if notices, ok := typed["notices"].([]any); ok {
			noticeErrors, noticeWarnings, noticeInfos := countNotices(notices)
			if noticeErrors+noticeWarnings+noticeInfos > 0 {
				infos -= len(notices)
				errors += noticeErrors
				warnings += noticeWarnings
				infos += noticeInfos
			}
		}
		for _, child := range []string{"summary", "validation_summary", "validationSummary", "report", "results"} {
			if nested, ok := typed[child]; ok {
				nestedErrors, nestedWarnings, nestedInfos := countFromValue(nested)
				errors += nestedErrors
				warnings += nestedWarnings
				infos += nestedInfos
			}
		}
		return errors, warnings, infos
	case []any:
		return countNotices(typed)
	default:
		return 0, 0, 0
	}
}

func countNotices(notices []any) (int, int, int) {
	var errors, warnings, infos int
	for _, notice := range notices {
		object, ok := notice.(map[string]any)
		if !ok {
			continue
		}
		severity := strings.ToLower(stringField(object, "severity", "level", "type"))
		switch severity {
		case "error", "fatal", "critical", "failure", "failed":
			errors++
		case "warning", "warn":
			warnings++
		case "info", "informational", "notice":
			infos++
		}
	}
	return errors, warnings, infos
}

func statusFromJSON(raw any) string {
	if object, ok := raw.(map[string]any); ok {
		status := strings.ToLower(stringField(object, "status", "validation_status", "validationStatus"))
		switch status {
		case "not_run", "passed", "warning", "failed":
			return status
		case "pass", "success", "ok":
			return "passed"
		case "warn":
			return "warning"
		case "error", "failure":
			return "failed"
		}
		for _, child := range []string{"summary", "validation_summary", "validationSummary", "report", "results"} {
			if status := statusFromJSON(object[child]); status != "" {
				return status
			}
		}
	}
	return ""
}

func numberField(object map[string]any, keys ...string) int {
	for _, key := range keys {
		value, ok := object[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case float64:
			return int(typed)
		case int:
			return typed
		case string:
			parsed, err := strconv.Atoi(strings.TrimSpace(typed))
			if err == nil {
				return parsed
			}
		}
	}
	return 0
}

func arrayLen(object map[string]any, keys ...string) int {
	for _, key := range keys {
		if values, ok := object[key].([]any); ok {
			return len(values)
		}
	}
	return 0
}

func stringField(object map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := object[key].(string); ok {
			return value
		}
	}
	return ""
}

func textCount(text string, pattern string) int {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(text, -1)
	total := 0
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		value, err := strconv.Atoi(match[2])
		if err == nil {
			total += value
		}
	}
	return total
}

func sortReportFiles(files []string) {
	for i := 0; i < len(files)-1; i++ {
		for j := i + 1; j < len(files); j++ {
			if filepath.Base(files[j]) == "report.json" || files[j] < files[i] {
				files[i], files[j] = files[j], files[i]
			}
		}
	}
}

func getenv(key string, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	if raw := strings.TrimSpace(os.Getenv(key)); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}

func int64FromEnv(key string, fallback int64) int64 {
	if raw := strings.TrimSpace(os.Getenv(key)); raw != "" {
		if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}

func max(left int, right int) int {
	if left > right {
		return left
	}
	return right
}
