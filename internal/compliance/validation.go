package compliance

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
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
	command := expandCommand(input.Command, input)
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	result.Report = map[string]any{
		"command": command,
		"stdout":  stdout.String(),
		"stderr":  stderr.String(),
	}
	if err != nil {
		result.Status = "failed"
		result.ErrorCount = 1
		result.Report["error"] = err.Error()
	} else {
		result.Status = "passed"
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

func expandCommand(command string, input ValidationRunInput) string {
	replacements := map[string]string{
		"{schedule_zip}": shellQuote(input.ScheduleZIPPath),
		"{realtime_pb}":  shellQuote(input.RealtimePBPath),
		"{feed_type}":    shellQuote(input.FeedType),
		"{output_dir}":   shellQuote("validation-output"),
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
