package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"open-transit-rt/internal/realtimequality"
)

const usageText = `Usage:
  go run ./cmd/realtime-quality-backtest --observed testdata/realtime-quality-backtest/observed-events.json --predictions testdata/realtime-quality-backtest/prediction-samples.json

This private CLI compares observed stop events against prediction samples and
writes aggregate realtime quality diagnostics under .cache/ by default.

Output files are exactly:
  summary.json
  summary.md
  metrics.json
  metrics.md
  manifest.json

The output is local engineering diagnostics only. It is not evidence, not a
consumer submission artifact, not a public API, and not production-grade ETA
proof.
`

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		var exit exitError
		if errors.As(err, &exit) && exit.silent {
			os.Exit(1)
		}
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		_, _ = fmt.Fprint(stdout, usageText)
		return nil
	}
	flags := flag.NewFlagSet("realtime-quality-backtest", flag.ContinueOnError)
	flags.SetOutput(stderr)
	observedPath := flags.String("observed", "", "path to observed stop events JSON")
	predictionPath := flags.String("predictions", "", "path to prediction samples JSON")
	outputDir := flags.String("output-dir", os.Getenv("OUTPUT_DIR"), "private output directory; default .cache/realtime-quality-backtest/<timestamp>")
	generatedAtRaw := flags.String("generated-at", "", "optional RFC3339 output timestamp")
	staleTTLRaw := flags.String("stale-ttl", realtimequality.DefaultStalePredictionTTL.String(), "stale prediction threshold duration")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	if *observedPath == "" || *predictionPath == "" {
		return fmt.Errorf("--observed and --predictions are required")
	}

	generatedAt := time.Now().UTC()
	if *generatedAtRaw != "" {
		parsed, err := time.Parse(time.RFC3339, *generatedAtRaw)
		if err != nil {
			return fmt.Errorf("--generated-at must be RFC3339: %w", err)
		}
		generatedAt = parsed.UTC()
	}
	staleTTL, err := time.ParseDuration(*staleTTLRaw)
	if err != nil {
		return fmt.Errorf("--stale-ttl must be a Go duration: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve working directory: %w", err)
	}
	output, err := realtimequality.ResolveBacktestOutputTarget(*outputDir, generatedAt, cwd)
	if err != nil {
		return fmt.Errorf("invalid output directory: %w", err)
	}
	observed, observedRaw, err := realtimequality.LoadObservedDataset(*observedPath)
	if err != nil {
		return err
	}
	predictions, predictionsRaw, err := realtimequality.LoadPredictionDataset(*predictionPath)
	if err != nil {
		return err
	}
	report, err := realtimequality.RunBacktest(observed, predictions, observedRaw, predictionsRaw, realtimequality.BacktestOptions{
		GeneratedAt:        generatedAt,
		StalePredictionTTL: staleTTL,
	})
	if err != nil {
		return err
	}
	files, err := realtimequality.BuildBacktestFiles(report, output, realtimequality.ForbiddenBacktestOutputValues(*observedPath, *predictionPath), 0)
	if err != nil {
		return err
	}
	if err := realtimequality.WriteBacktestFiles(files); err != nil {
		return exitError{message: err.Error(), silent: false}
	}
	_, _ = fmt.Fprint(stdout, files.Stdout)
	return nil
}

type exitError struct {
	message string
	silent  bool
}

func (e exitError) Error() string {
	return e.message
}
