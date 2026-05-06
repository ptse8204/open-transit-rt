package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"open-transit-rt/internal/gtfs"
)

type fakeImporter struct {
	result gtfs.ImportResult
	err    error
	opts   gtfs.ImportOptions
	ctx    context.Context
}

func (f *fakeImporter) ImportZip(ctx context.Context, opts gtfs.ImportOptions) (gtfs.ImportResult, error) {
	f.ctx = ctx
	f.opts = opts
	return f.result, f.err
}

func TestRunRequiresFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run(context.Background(), nil, &stdout, &stderr, &fakeImporter{})
	if exitCode != 2 {
		t.Fatalf("exit code = %d, want 2", exitCode)
	}
	if !strings.Contains(stderr.String(), "-agency-id is required") {
		t.Fatalf("stderr = %q, want missing agency message", stderr.String())
	}
}

func TestRunOutputsSuccessfulImportResult(t *testing.T) {
	fake := &fakeImporter{result: gtfs.ImportResult{
		ImportID:      7,
		AgencyID:      "demo-agency",
		FeedVersionID: "gtfs-import-7",
		Status:        gtfs.ImportStatusPublished,
		ReportStored:  true,
	}}
	var stdout, stderr bytes.Buffer
	exitCode := run(context.Background(), []string{"-agency-id", "demo-agency", "-zip", "/tmp/gtfs.zip", "-actor-id", "tester"}, &stdout, &stderr, fake)
	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %q", exitCode, stderr.String())
	}
	if fake.opts.AgencyID != "demo-agency" || fake.opts.ZipPath != "/tmp/gtfs.zip" || fake.opts.ActorID != "tester" {
		t.Fatalf("opts = %+v, want CLI flags forwarded", fake.opts)
	}
	var result gtfs.ImportResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode stdout: %v", err)
	}
	if result.Status != gtfs.ImportStatusPublished || result.FeedVersionID != "gtfs-import-7" {
		t.Fatalf("result = %+v, want published JSON", result)
	}
}

func TestRunAppliesConfiguredTimeout(t *testing.T) {
	fake := &fakeImporter{result: gtfs.ImportResult{
		AgencyID:     "demo-agency",
		Status:       gtfs.ImportStatusPublished,
		ReportStored: true,
	}}
	var stdout, stderr bytes.Buffer
	exitCode := run(context.Background(), []string{"-agency-id", "demo-agency", "-zip", "/tmp/gtfs.zip", "-timeout", "7m"}, &stdout, &stderr, fake)
	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %q", exitCode, stderr.String())
	}
	deadline, ok := fake.ctx.Deadline()
	if !ok {
		t.Fatalf("import context has no deadline")
	}
	remaining := time.Until(deadline)
	if remaining < 6*time.Minute || remaining > 7*time.Minute {
		t.Fatalf("deadline remaining = %s, want close to 7m", remaining)
	}
}

func TestLoadImportTimeoutFromEnv(t *testing.T) {
	t.Setenv("GTFS_IMPORT_TIMEOUT", "23m")
	if got := loadImportTimeout(); got != 23*time.Minute {
		t.Fatalf("loadImportTimeout = %s, want 23m", got)
	}
	t.Setenv("GTFS_IMPORT_TIMEOUT", "not-a-duration")
	if got := loadImportTimeout(); got != defaultImportTimeout {
		t.Fatalf("invalid loadImportTimeout = %s, want default %s", got, defaultImportTimeout)
	}
}

func TestRunOutputsFailedImportResult(t *testing.T) {
	result := gtfs.ImportResult{
		ImportID:       8,
		AgencyID:       "demo-agency",
		Status:         gtfs.ImportStatusFailed,
		ErrorCount:     1,
		ReportStored:   false,
		FailureMessage: "publish failed and failure report could not be stored",
	}
	fake := &fakeImporter{
		result: result,
		err:    &gtfs.ImportError{Result: result, Err: errors.New("gtfs import publish failed and failure report could not be stored")},
	}
	var stdout, stderr bytes.Buffer
	exitCode := run(context.Background(), []string{"-agency-id", "demo-agency", "-zip", "/tmp/gtfs.zip"}, &stdout, &stderr, fake)
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "failure report could not be stored") {
		t.Fatalf("stderr = %q, want clear report-storage failure", stderr.String())
	}
	var decoded gtfs.ImportResult
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("decode stdout: %v", err)
	}
	if decoded.ReportStored {
		t.Fatalf("decoded result = %+v, want report_stored=false", decoded)
	}
}
