package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"open-transit-rt/internal/avladapter"
	"open-transit-rt/internal/telemetry"
)

const usageText = `Usage:
  go run ./cmd/avl-vendor-adapter help
  go run ./cmd/avl-vendor-adapter --dry-run --mapping testdata/avl-vendor/mapping.json [--reference-time RFC3339] testdata/avl-vendor/minimal-gps.json
  go run ./cmd/avl-vendor-adapter --send --mapping testdata/avl-vendor/mapping.json --manifest testdata/avl-vendor/send-manifest.json testdata/avl-vendor/minimal-gps.json

This is a synthetic adapter kit example. Choose exactly one mode: --dry-run or
--send. Send mode posts transformed records to the configured /v1/telemetry
endpoint and writes private redacted diagnostics under .cache/ by default.

Output streams:
  dry-run stdout: transformed Open Transit RT telemetry JSON array only. If no
                  records transform successfully, stdout is [].
  dry-run stderr: diagnostics as one stable JSON array. Diagnostics are dry-run
                  review output only and are not telemetry ingest status.
  send stdout: redaction-safe send counts and repo-relative output reference.
  send stderr: redaction-safe fatal diagnostics only.

Synthetic fixtures, dry-run output, and private send diagnostics are not vendor
compatibility proof, production AVL reliability proof, or external evidence.
`

var newSendHTTPClient = func(timeout time.Duration) avladapter.HTTPDoer {
	return &http.Client{Timeout: timeout}
}

var sendSleeper avladapter.Sleeper = avladapter.SleepContext

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

	flags := flag.NewFlagSet("avl-vendor-adapter", flag.ContinueOnError)
	flags.SetOutput(stderr)
	dryRun := flags.Bool("dry-run", false, "transform only and do not send telemetry")
	send := flags.Bool("send", false, "send transformed telemetry to the configured /v1/telemetry endpoint")
	mappingPath := flags.String("mapping", "", "path to synthetic mapping JSON")
	manifestPath := flags.String("manifest", "", "path to send manifest JSON")
	referenceTimeRaw := flags.String("reference-time", "", "RFC3339 reference time for stale/future diagnostics")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *dryRun == *send {
		return fmt.Errorf("choose exactly one mode: --dry-run or --send")
	}
	if *mappingPath == "" {
		return fmt.Errorf("--mapping is required")
	}
	if *send && *manifestPath == "" {
		return fmt.Errorf("--manifest is required for --send")
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("exactly one vendor payload fixture path is required")
	}

	if *send {
		return runSend(*mappingPath, *manifestPath, flags.Arg(0), stdout, stderr)
	}

	referenceTime := time.Now().UTC()
	if *referenceTimeRaw != "" {
		parsed, err := time.Parse(time.RFC3339, *referenceTimeRaw)
		if err != nil {
			return fmt.Errorf("--reference-time must be RFC3339: %w", err)
		}
		referenceTime = parsed
	}

	mappingFile, err := os.Open(*mappingPath)
	if err != nil {
		return fmt.Errorf("open mapping: %w", err)
	}
	defer mappingFile.Close()
	mapping, mappingDiagnostics := avladapter.LoadMapping(mappingFile)
	if len(mappingDiagnostics) > 0 {
		writeDiagnostics(stderr, mappingDiagnostics)
		writeEvents(stdout, nil)
		if avladapter.HasHardErrors(mappingDiagnostics) {
			return exitError{message: "mapping contains hard errors", silent: true}
		}
	}

	payloadFile, err := os.Open(flags.Arg(0))
	if err != nil {
		return fmt.Errorf("open vendor payload: %w", err)
	}
	defer payloadFile.Close()
	result := avladapter.TransformPayload(payloadFile, mapping, avladapter.Options{ReferenceTime: referenceTime})
	writeEvents(stdout, result.Events)
	writeDiagnostics(stderr, result.Diagnostics)
	if avladapter.HasHardErrors(result.Diagnostics) {
		return exitError{message: "vendor payload contains hard errors", silent: true}
	}
	return nil
}

func runSend(mappingPath string, manifestPath string, payloadPath string, stdout io.Writer, stderr io.Writer) error {
	mappingFile, err := os.Open(mappingPath)
	if err != nil {
		return fmt.Errorf("open mapping")
	}
	defer mappingFile.Close()
	mapping, mappingDiagnostics := avladapter.LoadMapping(mappingFile)

	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("open send manifest")
	}
	defer manifestFile.Close()
	manifest, manifestDiagnostics := avladapter.LoadSendManifest(manifestFile)

	now := time.Now().UTC()
	cfg, configDiagnostics := avladapter.SendConfigFromEnv(manifest, os.Getenv, now)

	payloadRaw, err := os.ReadFile(payloadPath)
	if err != nil {
		return fmt.Errorf("open vendor payload")
	}
	payload, payloadDiagnostics := avladapter.DecodePayload(bytes.NewReader(payloadRaw))
	result := avladapter.Result{Diagnostics: payloadDiagnostics}
	if !avladapter.HasHardErrors(payloadDiagnostics) {
		result = avladapter.TransformDecodedPayload(payload, mapping, avladapter.Options{
			ReferenceTime:   cfg.ReferenceTime,
			StaleThreshold:  cfg.StaleThreshold,
			FutureThreshold: cfg.FutureThreshold,
		})
	}
	cwd, _ := os.Getwd()
	prepared, prepareDiagnostics := avladapter.PrepareSend(manifest, cfg, result, payload, os.Getenv, cwd)
	diagnostics := append([]avladapter.Diagnostic{}, mappingDiagnostics...)
	diagnostics = append(diagnostics, manifestDiagnostics...)
	diagnostics = append(diagnostics, configDiagnostics...)
	diagnostics = append(diagnostics, result.Diagnostics...)
	diagnostics = append(diagnostics, prepareDiagnostics...)
	if avladapter.HasHardErrors(diagnostics) {
		writeDiagnostics(stderr, diagnostics)
		return exitError{message: "send preflight failed", silent: true}
	}

	report := avladapter.SendEvents(context.Background(), result.Events, prepared, newSendHTTPClient(cfg.Timeout), sendSleeper)
	report, outputDiagnostics := avladapter.BuildSendFiles(report, prepared)
	if len(outputDiagnostics) > 0 {
		writeDiagnostics(stderr, outputDiagnostics)
		return exitError{message: "send output failed", silent: true}
	}
	if err := avladapter.WriteSendFiles(report, prepared.Output); err != nil {
		writeDiagnostics(stderr, []avladapter.Diagnostic{{Code: avladapter.CodeInvalidOutputPath, Severity: avladapter.SeverityError, Message: "send output write failed"}})
		return exitError{message: "send output write failed", silent: true}
	}
	_, _ = fmt.Fprint(stdout, report.Stdout)
	if report.Summary.FailedCount > 0 {
		return exitError{message: "send failed", silent: true}
	}
	return nil
}

func writeEvents(stdout io.Writer, events []telemetry.Event) {
	raw, err := avladapter.MarshalEvents(events)
	if err != nil {
		raw = []byte("[]")
	}
	_, _ = fmt.Fprintln(stdout, string(raw))
}

func writeDiagnostics(stderr io.Writer, diagnostics []avladapter.Diagnostic) {
	raw, err := avladapter.MarshalDiagnostics(diagnostics)
	if err != nil {
		raw = []byte(`[]`)
	}
	_, _ = fmt.Fprintln(stderr, string(raw))
}

type exitError struct {
	message string
	silent  bool
}

func (e exitError) Error() string {
	return e.message
}
