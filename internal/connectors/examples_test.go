package connectors

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestExampleConnectorManifestsAndFixtures(t *testing.T) {
	root := filepath.Join("..", "..")
	expected := []string{
		"monitoring-export",
		"predictor-sidecar-stub",
		"telemetry-csv-replay",
		"telemetry-http-poller",
	}
	manifests, err := filepath.Glob(filepath.Join(root, "examples", "connectors", "*", "connector.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(manifests) != len(expected) {
		t.Fatalf("example manifests = %d, want %d: %v", len(manifests), len(expected), manifests)
	}
	var names []string
	for _, manifestPath := range manifests {
		names = append(names, filepath.Base(filepath.Dir(manifestPath)))
	}
	sort.Strings(names)
	if strings.Join(names, ",") != strings.Join(expected, ",") {
		t.Fatalf("example connector names = %v, want %v", names, expected)
	}

	for _, manifestPath := range manifests {
		t.Run(filepath.Base(filepath.Dir(manifestPath)), func(t *testing.T) {
			raw, err := os.Open(manifestPath)
			if err != nil {
				t.Fatal(err)
			}
			manifest, err := DecodeManifest(raw)
			closeErr := raw.Close()
			if err != nil {
				t.Fatal(err)
			}
			if closeErr != nil {
				t.Fatal(closeErr)
			}
			if !strings.HasPrefix(manifest.ConnectorID, "example.") {
				t.Fatalf("connector_id = %q, want example prefix", manifest.ConnectorID)
			}
			for _, tc := range manifest.ConformanceCases {
				fixturePath := filepath.Join(root, filepath.FromSlash(tc.FixturePath))
				if _, err := os.Stat(fixturePath); err != nil {
					t.Fatalf("fixture %s: %v", tc.FixturePath, err)
				}
				assertSyntheticFixture(t, fixturePath)
			}
		})
	}
}

func assertSyntheticFixture(t *testing.T, path string) {
	t.Helper()
	switch filepath.Ext(path) {
	case ".json":
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var header struct {
			SyntheticOnly bool `json:"synthetic_only"`
		}
		if err := json.Unmarshal(raw, &header); err != nil {
			t.Fatalf("decode %s: %v", path, err)
		}
		if !header.SyntheticOnly {
			t.Fatalf("%s must set synthetic_only true", path)
		}
	case ".csv":
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		rows, err := csv.NewReader(f).ReadAll()
		if err != nil {
			t.Fatalf("decode %s: %v", path, err)
		}
		if len(rows) < 2 || len(rows[0]) == 0 || rows[0][0] != "synthetic_only" {
			t.Fatalf("%s must start with synthetic_only column", path)
		}
		for i, row := range rows[1:] {
			if len(row) == 0 || row[0] != "true" {
				t.Fatalf("%s row %d must set synthetic_only true", path, i+2)
			}
		}
	default:
		t.Fatalf("unsupported example fixture type: %s", path)
	}
}
