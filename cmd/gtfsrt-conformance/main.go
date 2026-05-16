package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"open-transit-rt/internal/gtfsrtconformance"
)

const maxPayloadBytes = 10 * 1024 * 1024

type config struct {
	vehiclePositions string
	tripUpdates      string
	alerts           string
	output           string
	generatedAt      string
	help             bool
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	cfg, err := parse(args, stdout)
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return 2
	}
	if cfg.help {
		usage(stdout)
		return 0
	}
	inputs, err := readInputs(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "gtfsrt conformance failed: %s\n", err)
		return 1
	}
	generatedAt, err := parseGeneratedAt(cfg.generatedAt)
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return 2
	}
	report := gtfsrtconformance.BuildReport(inputs, gtfsrtconformance.Options{GeneratedAt: generatedAt})
	if err := writeReport(stdout, cfg.output, report); err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return 2
	}
	if report.OverallStatus == gtfsrtconformance.StatusFailed {
		return 1
	}
	return 0
}

func parse(args []string, out io.Writer) (config, error) {
	cfg := config{output: "json"}
	fs := flag.NewFlagSet("gtfsrt-conformance", flag.ContinueOnError)
	fs.SetOutput(out)
	fs.StringVar(&cfg.vehiclePositions, "vehicle-positions", "", "local Vehicle Positions protobuf path")
	fs.StringVar(&cfg.tripUpdates, "trip-updates", "", "local Trip Updates protobuf path")
	fs.StringVar(&cfg.alerts, "alerts", "", "local Alerts protobuf path")
	fs.StringVar(&cfg.output, "output", cfg.output, "json or text")
	fs.StringVar(&cfg.generatedAt, "generated-at", "", "RFC3339 timestamp used for diagnostics")
	fs.BoolVar(&cfg.help, "help", false, "show help")
	if err := fs.Parse(args); err != nil {
		return cfg, err
	}
	if fs.NArg() != 0 {
		return cfg, fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}
	if cfg.help {
		return cfg, nil
	}
	if cfg.vehiclePositions == "" && cfg.tripUpdates == "" && cfg.alerts == "" {
		return cfg, errors.New("at least one --vehicle-positions, --trip-updates, or --alerts path is required")
	}
	switch cfg.output {
	case "json", "text":
	default:
		return cfg, errors.New("--output must be json or text")
	}
	return cfg, nil
}

func usage(w io.Writer) {
	fmt.Fprintln(w, `Usage:
  gtfsrt-conformance --vehicle-positions path/to/vehicle_positions.pb
  gtfsrt-conformance --trip-updates path/to/trip_updates.pb --alerts path/to/alerts.pb --output text

Offline only: parses local GTFS-RT protobuf artifacts and checks a bounded
interoperability profile. It does not fetch URLs, contact consumers, write
evidence, move consumer statuses, or prove compliance, consumer acceptance,
production readiness, SLA/uptime, vendor compatibility, hardware
certification, production-grade ETA quality, or real-world ETA accuracy.`)
}

func readInputs(cfg config) ([]gtfsrtconformance.FeedInput, error) {
	candidates := []struct {
		feedType string
		path     string
	}{
		{gtfsrtconformance.FeedVehiclePositions, cfg.vehiclePositions},
		{gtfsrtconformance.FeedTripUpdates, cfg.tripUpdates},
		{gtfsrtconformance.FeedAlerts, cfg.alerts},
	}
	inputs := []gtfsrtconformance.FeedInput{}
	for _, candidate := range candidates {
		if strings.TrimSpace(candidate.path) == "" {
			continue
		}
		payload, source, err := readPayload(candidate.path)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, gtfsrtconformance.FeedInput{
			FeedType: candidate.feedType,
			Source:   source,
			Payload:  payload,
		})
	}
	return inputs, nil
}

func readPayload(rawPath string) ([]byte, string, error) {
	clean := filepath.Clean(rawPath)
	if isEvidenceLike(clean) {
		return nil, "", fmt.Errorf("refusing to read evidence-like path %q", rawPath)
	}
	info, err := os.Stat(clean)
	if err != nil {
		return nil, "", err
	}
	if info.IsDir() {
		return nil, "", fmt.Errorf("%q is a directory", rawPath)
	}
	if info.Size() > maxPayloadBytes {
		return nil, "", fmt.Errorf("%q exceeds %d bytes", rawPath, maxPayloadBytes)
	}
	payload, err := os.ReadFile(clean)
	if err != nil {
		return nil, "", err
	}
	return payload, filepath.ToSlash(clean), nil
}

func parseGeneratedAt(raw string) (time.Time, error) {
	if strings.TrimSpace(raw) == "" {
		return time.Now().UTC().Truncate(time.Second), nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("--generated-at must be RFC3339: %w", err)
	}
	return t.UTC(), nil
}

func writeReport(w io.Writer, output string, report gtfsrtconformance.Report) error {
	switch output {
	case "json":
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(report)
	case "text":
		fmt.Fprintf(w, "GTFS-RT conformance: %s\n", report.OverallStatus)
		for _, feed := range report.Feeds {
			fmt.Fprintf(w, "- %s: %s (%d entities)\n", feed.FeedType, feed.Status, feed.EntityCount)
			for _, check := range feed.Checks {
				fmt.Fprintf(w, "  %s %s: %s\n", check.Severity, check.ID, check.Message)
			}
		}
		return nil
	default:
		return errors.New("--output must be json or text")
	}
}

func isEvidenceLike(path string) bool {
	normalized := strings.ToLower(filepath.ToSlash(path))
	for _, marker := range []string{"docs/evidence", "/evidence/", "evidence/", "/submission", "submission/", "/proof", "proof/"} {
		if strings.Contains(normalized, marker) {
			return true
		}
	}
	return false
}
