package compliance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ValidationStore interface {
	StoreValidationResult(ctx context.Context, result ValidationResult) error
}

func RunValidation(ctx context.Context, store ValidationStore, input ValidationRunInput) (ValidationResult, error) {
	result := ValidationResult{
		AgencyID:         input.AgencyID,
		FeedType:         input.FeedType,
		FeedVersionID:    input.FeedVersionID,
		ValidatorName:    input.ValidatorName,
		ValidatorVersion: input.ValidatorVersion,
		Status:           "not_run",
		Report:           map[string]any{},
	}
	if result.ValidatorName == "" {
		result.ValidatorName = "canonical-validator"
	}
	if input.Command == "" {
		result.Report = map[string]any{"reason": "validator_command_missing"}
		if err := store.StoreValidationResult(ctx, result); err != nil {
			return ValidationResult{}, err
		}
		return result, nil
	}
	outputDir, err := os.MkdirTemp("", "open-transit-rt-validator-*")
	if err != nil {
		return ValidationResult{}, fmt.Errorf("create validator output dir: %w", err)
	}
	defer os.RemoveAll(outputDir)

	command := expandCommand(input.Command, input, outputDir)
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	normalized := normalizeValidatorOutput(stdout.String(), stderr.String(), outputDir)
	result.ErrorCount = normalized.errorCount
	result.WarningCount = normalized.warningCount
	result.InfoCount = normalized.infoCount
	result.Report = normalized.report
	result.Report["command"] = command
	result.Report["stdout"] = stdout.String()
	result.Report["stderr"] = stderr.String()
	if normalized.source != "" {
		result.Report["parsed_report_source"] = normalized.source
	}
	if runErr != nil {
		result.Status = "failed"
		if result.ErrorCount == 0 {
			result.ErrorCount = 1
		}
		result.Report["error"] = runErr.Error()
	} else {
		result.Status = normalized.status()
	}
	if err := store.StoreValidationResult(ctx, result); err != nil {
		return ValidationResult{}, err
	}
	return result, nil
}

func StaticValidatorCommand(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasSuffix(strings.ToLower(path), ".jar") {
		return fmt.Sprintf("java -jar %s -i {schedule_zip} -o {output_dir}", shellQuote(path))
	}
	return fmt.Sprintf("%s -i {schedule_zip} -o {output_dir}", shellQuote(path))
}

func expandCommand(command string, input ValidationRunInput, outputDir string) string {
	replacements := map[string]string{
		"{schedule_zip}": shellQuote(input.ScheduleZIPPath),
		"{realtime_pb}":  shellQuote(input.RealtimePBPath),
		"{feed_type}":    shellQuote(input.FeedType),
		"{output_dir}":   shellQuote(outputDir),
	}
	expanded := command
	for old, next := range replacements {
		expanded = strings.ReplaceAll(expanded, old, next)
	}
	return expanded
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
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

func normalizeValidatorOutput(stdout string, stderr string, outputDir string) normalizedValidationOutput {
	result := normalizedValidationOutput{report: map[string]any{}}
	if raw, source, ok := firstJSONReport(stdout, stderr, outputDir); ok {
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

func firstJSONReport(stdout string, stderr string, outputDir string) (any, string, bool) {
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
		payload, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		if raw, ok := parseJSON(string(payload)); ok {
			return raw, file, true
		}
	}
	return nil, "", false
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
