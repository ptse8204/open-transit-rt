package schedule

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"io"
	"testing"
	"time"
)

func TestWriteZipCSVUsesDeterministicEntryTimestampAndRows(t *testing.T) {
	revision := time.Date(2026, 4, 21, 20, 15, 0, 0, time.UTC)
	rows := [][]string{{"route_id", "route_type"}, {"route-10", "3"}}
	first := buildTestZip(t, revision, rows)
	second := buildTestZip(t, revision, rows)
	if !bytes.Equal(first, second) {
		t.Fatalf("zip bytes differ for identical rows and revision timestamp")
	}
	reader, err := zip.NewReader(bytes.NewReader(first), int64(len(first)))
	if err != nil {
		t.Fatalf("read zip: %v", err)
	}
	if len(reader.File) != 1 || reader.File[0].Name != "routes.txt" {
		t.Fatalf("zip files = %+v, want routes.txt", reader.File)
	}
	if reader.File[0].Method != zip.Deflate {
		t.Fatalf("zip method = %d, want deflate", reader.File[0].Method)
	}
	if !reader.File[0].Modified.Equal(revision) {
		t.Fatalf("modified = %s, want revision %s", reader.File[0].Modified, revision)
	}
	rc, err := reader.File[0].Open()
	if err != nil {
		t.Fatalf("open routes.txt: %v", err)
	}
	defer rc.Close()
	content, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read routes.txt: %v", err)
	}
	parsed, err := csv.NewReader(bytes.NewReader(content)).ReadAll()
	if err != nil {
		t.Fatalf("parse routes csv: %v", err)
	}
	if parsed[1][0] != "route-10" || parsed[1][1] != "3" {
		t.Fatalf("rows = %+v, want route-10 row", parsed)
	}
}

func buildTestZip(t *testing.T, revision time.Time, rows [][]string) []byte {
	t.Helper()
	var buffer bytes.Buffer
	zw := zip.NewWriter(&buffer)
	if err := writeZipCSV(zw, "routes.txt", revision, rows); err != nil {
		t.Fatalf("write zip csv: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buffer.Bytes()
}
