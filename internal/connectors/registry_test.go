package connectors

import (
	"encoding/json"
	"path"
	"sort"
	"strings"
	"testing"
)

func TestLoadExampleRegistryReturnsCommittedExampleManifests(t *testing.T) {
	registry := LoadExampleRegistry()
	if len(registry.Diagnostics) != 0 {
		t.Fatalf("LoadExampleRegistry diagnostics = %+v", registry.Diagnostics)
	}
	if len(registry.Entries) != 8 {
		t.Fatalf("entries = %d, want 8: %+v", len(registry.Entries), registry.Entries)
	}

	var ids []string
	for _, entry := range registry.Entries {
		ids = append(ids, entry.ConnectorID)
		if entry.SchemaVersion != SchemaVersion {
			t.Fatalf("%s schema_version = %q", entry.SourcePath, entry.SchemaVersion)
		}
		if !strings.HasPrefix(entry.ConnectorID, "example.") {
			t.Fatalf("%s connector_id = %q, want example prefix", entry.SourcePath, entry.ConnectorID)
		}
		if !strings.HasPrefix(entry.SourcePath, exampleRegistryRoot+"/") || path.Base(entry.SourcePath) != "connector.json" {
			t.Fatalf("source_path = %q, want committed example connector manifest path", entry.SourcePath)
		}
		if path.IsAbs(entry.SourcePath) || strings.Contains(entry.SourcePath, "..") {
			t.Fatalf("source_path = %q, want clean relative path", entry.SourcePath)
		}
		if entry.DisplayName == "" || entry.Description == "" || entry.ModeName == "" || entry.DocsLink == "" {
			t.Fatalf("%s has empty registry summary fields: %+v", entry.SourcePath, entry)
		}
		if !entry.DisabledByDefault || !entry.FailureBehavior.FailClosed {
			t.Fatalf("%s must remain disabled-by-default and fail-closed: %+v", entry.SourcePath, entry)
		}
		if len(entry.InputContracts) == 0 || len(entry.OutputContracts) == 0 || len(entry.ConformanceCases) == 0 {
			t.Fatalf("%s must expose bounded contracts and conformance cases: %+v", entry.SourcePath, entry)
		}
		assertRegistryEntryBounded(t, entry)
	}
	if !sort.StringsAreSorted(ids) {
		t.Fatalf("registry entries are not sorted by connector_id: %v", ids)
	}
}

func TestExampleRegistryJSONDoesNotExposeUnsafeIntegrationState(t *testing.T) {
	registry := LoadExampleRegistry()
	raw, err := json.Marshal(registry)
	if err != nil {
		t.Fatal(err)
	}
	body := strings.ToLower(string(raw))
	for _, forbidden := range []string{
		"raw_validator_command",
		"raw_command",
		"shell",
		"argv",
		"api_key",
		"password",
		"/users/",
		"/tmp/",
		"file://",
		"production ready",
		"vendor compatible",
		"consumer accepted",
	} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("registry JSON contains forbidden %q: %s", forbidden, body)
		}
	}
}

func TestRegistryEntryFromManifestBoundsLargeFieldsAndLists(t *testing.T) {
	manifest := validManifest(t)
	manifest.DisplayName = strings.Repeat("a", maxRegistryTextBytes+10)
	manifest.RedactionPolicy.RedactFields = repeatedStrings("field", maxRegistryItems+3)
	manifest.InputContracts = repeatedContracts(manifest.InputContracts[0], maxRegistryItems+2)
	manifest.ConformanceCases = repeatedCases(manifest.ConformanceCases[0], maxRegistryItems+2)

	var diagnostics []RegistryDiagnostic
	entry := registryEntryFromManifest("examples/connectors/demo/connector.json", manifest, &diagnostics)
	if len(entry.DisplayName) != maxRegistryTextBytes {
		t.Fatalf("display_name length = %d, want %d", len(entry.DisplayName), maxRegistryTextBytes)
	}
	if len(entry.RedactionPolicy.RedactFields) != maxRegistryItems {
		t.Fatalf("redact_fields length = %d, want %d", len(entry.RedactionPolicy.RedactFields), maxRegistryItems)
	}
	if len(entry.InputContracts) != maxRegistryItems {
		t.Fatalf("input_contracts length = %d, want %d", len(entry.InputContracts), maxRegistryItems)
	}
	if len(entry.ConformanceCases) != maxRegistryItems {
		t.Fatalf("conformance_cases length = %d, want %d", len(entry.ConformanceCases), maxRegistryItems)
	}
	if len(diagnostics) == 0 {
		t.Fatal("expected bounding diagnostics")
	}
	for _, diagnostic := range diagnostics {
		if diagnostic.Level != "warning" || diagnostic.Code != "registry_entry_bounded" {
			t.Fatalf("diagnostic = %+v, want registry_entry_bounded warning", diagnostic)
		}
	}
}

func assertRegistryEntryBounded(t *testing.T, entry RegistryEntry) {
	t.Helper()
	for field, value := range map[string]string{
		"source_path":  entry.SourcePath,
		"connector_id": entry.ConnectorID,
		"type":         entry.ConnectorType,
		"display_name": entry.DisplayName,
		"description":  entry.Description,
		"mode_name":    entry.ModeName,
		"docs_link":    entry.DocsLink,
	} {
		if len(value) > maxRegistryTextBytes {
			t.Fatalf("%s length = %d, exceeds bound %d", field, len(value), maxRegistryTextBytes)
		}
	}
	if len(entry.InputContracts) > maxRegistryItems || len(entry.OutputContracts) > maxRegistryItems || len(entry.ConformanceCases) > maxRegistryItems {
		t.Fatalf("entry list lengths exceed registry bounds: %+v", entry)
	}
}

func repeatedStrings(prefix string, count int) []string {
	values := make([]string, 0, count)
	for i := 0; i < count; i++ {
		values = append(values, prefix)
	}
	return values
}

func repeatedContracts(contract Contract, count int) []Contract {
	contracts := make([]Contract, 0, count)
	for i := 0; i < count; i++ {
		contracts = append(contracts, contract)
	}
	return contracts
}

func repeatedCases(tc ConformanceCase, count int) []ConformanceCase {
	cases := make([]ConformanceCase, 0, count)
	for i := 0; i < count; i++ {
		cases = append(cases, tc)
	}
	return cases
}
