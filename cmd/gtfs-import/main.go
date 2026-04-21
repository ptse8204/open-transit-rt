package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	appdb "open-transit-rt/internal/db"
	"open-transit-rt/internal/gtfs"
)

type importRunner interface {
	ImportZip(ctx context.Context, opts gtfs.ImportOptions) (gtfs.ImportResult, error)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pool, err := appdb.Connect(ctx, appdb.LoadConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	exitCode := run(ctx, os.Args[1:], os.Stdout, os.Stderr, gtfs.NewImportService(pool))
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func run(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer, runner importRunner) int {
	fs := flag.NewFlagSet("gtfs-import", flag.ContinueOnError)
	fs.SetOutput(stderr)
	agencyID := fs.String("agency-id", "", "agency id to import for")
	zipPath := fs.String("zip", "", "path to gtfs.zip")
	actorID := fs.String("actor-id", "system", "actor id for audit logging")
	notes := fs.String("notes", "", "optional import notes")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *agencyID == "" {
		fmt.Fprintln(stderr, "-agency-id is required")
		return 2
	}
	if *zipPath == "" {
		fmt.Fprintln(stderr, "-zip is required")
		return 2
	}
	if runner == nil {
		fmt.Fprintln(stderr, "import runner is required")
		return 2
	}

	result, err := runner.ImportZip(ctx, gtfs.ImportOptions{
		AgencyID: *agencyID,
		ZipPath:  *zipPath,
		ActorID:  *actorID,
		Notes:    *notes,
	})
	if encodeErr := json.NewEncoder(stdout).Encode(result); encodeErr != nil {
		fmt.Fprintf(stderr, "write import result: %v\n", encodeErr)
		return 1
	}
	if err != nil {
		var importErr *gtfs.ImportError
		if errors.As(err, &importErr) {
			fmt.Fprintf(stderr, "%v\n", importErr.Err)
		} else {
			fmt.Fprintf(stderr, "%v\n", err)
		}
		return 1
	}
	return 0
}
