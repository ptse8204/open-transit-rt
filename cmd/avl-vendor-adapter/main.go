package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"open-transit-rt/internal/avladapter"
	"open-transit-rt/internal/telemetry"
)

const usageText = `Usage:
  go run ./cmd/avl-vendor-adapter help
  go run ./cmd/avl-vendor-adapter --dry-run --mapping testdata/avl-vendor/mapping.json [--reference-time RFC3339] testdata/avl-vendor/valid.json

Phase 29B supports dry-run transforms only. Running without --dry-run fails
because network send mode is not implemented.

Output streams:
  stdout: transformed Open Transit RT telemetry JSON array. If no records
          transform successfully, stdout is [].
  stderr: diagnostics as one stable JSON array. Diagnostics are dry-run review
          output only and are not telemetry ingest acceptance status.
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

	flags := flag.NewFlagSet("avl-vendor-adapter", flag.ContinueOnError)
	flags.SetOutput(stderr)
	dryRun := flags.Bool("dry-run", false, "required; transform only and do not send telemetry")
	mappingPath := flags.String("mapping", "", "path to synthetic mapping JSON")
	referenceTimeRaw := flags.String("reference-time", "", "RFC3339 reference time for stale/future diagnostics")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if !*dryRun {
		return fmt.Errorf("send mode is not implemented in Phase 29B; rerun with --dry-run to print transformed telemetry JSON")
	}
	if *mappingPath == "" {
		return fmt.Errorf("--mapping is required")
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("exactly one vendor payload fixture path is required")
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
