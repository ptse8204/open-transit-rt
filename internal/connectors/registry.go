package connectors

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

const (
	exampleRegistryRoot = "examples/connectors"

	maxRegistryEntries      = 64
	maxRegistryManifestSize = 64 * 1024
	maxRegistryTextBytes    = 512
	maxRegistryItems        = 16
)

type Registry struct {
	Entries     []RegistryEntry      `json:"entries"`
	Diagnostics []RegistryDiagnostic `json:"diagnostics"`
}

type RegistryEntry struct {
	SourcePath        string                    `json:"source_path"`
	SchemaVersion     string                    `json:"schema_version"`
	ConnectorID       string                    `json:"connector_id"`
	ConnectorType     string                    `json:"type"`
	DisplayName       string                    `json:"display_name"`
	Description       string                    `json:"description"`
	ModeName          string                    `json:"mode_name"`
	DisabledByDefault bool                      `json:"disabled_by_default"`
	DocsLink          string                    `json:"docs_link"`
	InputContracts    []RegistryContract        `json:"input_contracts"`
	OutputContracts   []RegistryContract        `json:"output_contracts"`
	FailureBehavior   RegistryFailureBehavior   `json:"failure_behavior"`
	RedactionPolicy   RegistryRedactionPolicy   `json:"redaction_policy"`
	ClaimBoundary     RegistryClaimBoundary     `json:"claim_boundary"`
	ConformanceCases  []RegistryConformanceCase `json:"conformance_cases"`
}

type RegistryContract struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Schema         string   `json:"schema"`
	MediaTypes     []string `json:"media_types,omitempty"`
	RequiredFields []string `json:"required_fields,omitempty"`
	Produces       []string `json:"produces,omitempty"`
}

type RegistryFailureBehavior struct {
	TimeoutSeconds int    `json:"timeout_seconds"`
	RetryPolicy    string `json:"retry_policy"`
	DegradedState  string `json:"degraded_state"`
	FailClosed     bool   `json:"fail_closed"`
}

type RegistryRedactionPolicy struct {
	SecretStorage string   `json:"secret_storage"`
	RedactFields  []string `json:"redact_fields"`
}

type RegistryClaimBoundary struct {
	PositiveClaims []string `json:"positive_claims"`
	NotClaimed     []string `json:"not_claimed"`
}

type RegistryConformanceCase struct {
	ID             string `json:"id"`
	Description    string `json:"description"`
	FixturePath    string `json:"fixture_path"`
	ExpectedResult string `json:"expected_result"`
}

type RegistryDiagnostic struct {
	Level   string `json:"level"`
	Code    string `json:"code"`
	Path    string `json:"path,omitempty"`
	Message string `json:"message"`
}

func LoadExampleRegistry() Registry {
	repoRoot, diag := sourceRepoRoot()
	if diag != nil {
		return Registry{Diagnostics: []RegistryDiagnostic{*diag}}
	}
	return loadExampleRegistry(repoRoot)
}

func sourceRepoRoot() (string, *RegistryDiagnostic) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", &RegistryDiagnostic{
			Level:   "error",
			Code:    "registry_source_unavailable",
			Message: "connector registry source path is unavailable",
		}
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..")), nil
}

func loadExampleRegistry(repoRoot string) Registry {
	var registry Registry
	root := filepath.Join(repoRoot, filepath.FromSlash(exampleRegistryRoot))
	dirs, err := os.ReadDir(root)
	if err != nil {
		registry.Diagnostics = append(registry.Diagnostics, RegistryDiagnostic{
			Level:   "error",
			Code:    "registry_examples_unreadable",
			Path:    exampleRegistryRoot,
			Message: "example connector registry root is not readable",
		})
		return registry
	}

	var connectorDirs []os.DirEntry
	for _, dir := range dirs {
		if dir.Type()&os.ModeSymlink != 0 {
			registry.Diagnostics = append(registry.Diagnostics, RegistryDiagnostic{
				Level:   "error",
				Code:    "registry_symlink_rejected",
				Path:    filepath.ToSlash(filepath.Join(exampleRegistryRoot, dir.Name())),
				Message: "example connector registry entries must not be symlinks",
			})
			continue
		}
		if dir.IsDir() {
			connectorDirs = append(connectorDirs, dir)
		}
	}
	slices.SortFunc(connectorDirs, func(a, b os.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})
	if len(connectorDirs) > maxRegistryEntries {
		registry.Diagnostics = append(registry.Diagnostics, RegistryDiagnostic{
			Level:   "warning",
			Code:    "registry_entry_limit",
			Path:    exampleRegistryRoot,
			Message: fmt.Sprintf("registry is limited to %d example connector manifests", maxRegistryEntries),
		})
		connectorDirs = connectorDirs[:maxRegistryEntries]
	}

	seen := make(map[string]string)
	for _, dir := range connectorDirs {
		sourcePath := filepath.ToSlash(filepath.Join(exampleRegistryRoot, dir.Name(), "connector.json"))
		manifest, ok := readExampleManifest(filepath.Join(root, dir.Name(), "connector.json"), sourcePath, &registry.Diagnostics)
		if !ok {
			continue
		}
		if previous := seen[manifest.ConnectorID]; previous != "" {
			registry.Diagnostics = append(registry.Diagnostics, RegistryDiagnostic{
				Level:   "error",
				Code:    "registry_duplicate_connector_id",
				Path:    sourcePath,
				Message: "connector_id duplicates " + previous,
			})
			continue
		}
		seen[manifest.ConnectorID] = sourcePath
		registry.Entries = append(registry.Entries, registryEntryFromManifest(sourcePath, manifest, &registry.Diagnostics))
	}
	slices.SortFunc(registry.Entries, func(a, b RegistryEntry) int {
		return strings.Compare(a.ConnectorID, b.ConnectorID)
	})
	return registry
}

func readExampleManifest(path string, sourcePath string, diagnostics *[]RegistryDiagnostic) (Manifest, bool) {
	info, err := os.Lstat(path)
	if err != nil {
		*diagnostics = append(*diagnostics, RegistryDiagnostic{
			Level:   "error",
			Code:    "registry_manifest_missing",
			Path:    sourcePath,
			Message: "example connector manifest is not readable",
		})
		return Manifest{}, false
	}
	if info.Mode()&os.ModeSymlink != 0 {
		*diagnostics = append(*diagnostics, RegistryDiagnostic{
			Level:   "error",
			Code:    "registry_symlink_rejected",
			Path:    sourcePath,
			Message: "example connector manifest must not be a symlink",
		})
		return Manifest{}, false
	}
	if !info.Mode().IsRegular() {
		*diagnostics = append(*diagnostics, RegistryDiagnostic{
			Level:   "error",
			Code:    "registry_manifest_not_regular",
			Path:    sourcePath,
			Message: "example connector manifest must be a regular file",
		})
		return Manifest{}, false
	}
	if info.Size() > maxRegistryManifestSize {
		*diagnostics = append(*diagnostics, RegistryDiagnostic{
			Level:   "error",
			Code:    "registry_manifest_too_large",
			Path:    sourcePath,
			Message: fmt.Sprintf("example connector manifest exceeds %d bytes", maxRegistryManifestSize),
		})
		return Manifest{}, false
	}

	f, err := os.Open(path)
	if err != nil {
		*diagnostics = append(*diagnostics, RegistryDiagnostic{
			Level:   "error",
			Code:    "registry_manifest_open_failed",
			Path:    sourcePath,
			Message: "example connector manifest could not be opened",
		})
		return Manifest{}, false
	}
	defer f.Close()

	manifest, err := DecodeManifest(io.LimitReader(f, maxRegistryManifestSize+1))
	if err != nil {
		*diagnostics = append(*diagnostics, RegistryDiagnostic{
			Level:   "error",
			Code:    "registry_manifest_invalid",
			Path:    sourcePath,
			Message: err.Error(),
		})
		return Manifest{}, false
	}
	return manifest, true
}

func registryEntryFromManifest(sourcePath string, manifest Manifest, diagnostics *[]RegistryDiagnostic) RegistryEntry {
	entry := RegistryEntry{
		SourcePath:        boundedRegistryText(sourcePath, sourcePath, "source_path", diagnostics),
		SchemaVersion:     boundedRegistryText(sourcePath, manifest.SchemaVersion, "schema_version", diagnostics),
		ConnectorID:       boundedRegistryText(sourcePath, manifest.ConnectorID, "connector_id", diagnostics),
		ConnectorType:     boundedRegistryText(sourcePath, manifest.ConnectorType, "type", diagnostics),
		DisplayName:       boundedRegistryText(sourcePath, manifest.DisplayName, "display_name", diagnostics),
		Description:       boundedRegistryText(sourcePath, manifest.Description, "description", diagnostics),
		ModeName:          boundedRegistryText(sourcePath, manifest.Mode.Name, "mode.name", diagnostics),
		DisabledByDefault: manifest.Mode.DisabledByDefault,
		DocsLink:          boundedRegistryText(sourcePath, manifest.DocsLink, "docs_link", diagnostics),
		FailureBehavior: RegistryFailureBehavior{
			TimeoutSeconds: manifest.FailureBehavior.TimeoutSeconds,
			RetryPolicy:    boundedRegistryText(sourcePath, manifest.FailureBehavior.RetryPolicy, "failure_behavior.retry_policy", diagnostics),
			DegradedState:  boundedRegistryText(sourcePath, manifest.FailureBehavior.DegradedState, "failure_behavior.degraded_state", diagnostics),
			FailClosed:     manifest.FailureBehavior.FailClosed,
		},
		RedactionPolicy: RegistryRedactionPolicy{
			SecretStorage: boundedRegistryText(sourcePath, manifest.RedactionPolicy.SecretStorage, "redaction_policy.secret_storage", diagnostics),
			RedactFields:  boundedRegistryTextList(sourcePath, "redaction_policy.redact_fields", manifest.RedactionPolicy.RedactFields, diagnostics),
		},
		ClaimBoundary: RegistryClaimBoundary{
			PositiveClaims: boundedRegistryTextList(sourcePath, "claim_boundary.positive_claims", manifest.ClaimBoundary.PositiveClaims, diagnostics),
			NotClaimed:     boundedRegistryTextList(sourcePath, "claim_boundary.not_claimed", manifest.ClaimBoundary.NotClaimed, diagnostics),
		},
	}

	for i, contract := range boundedRegistryContracts(sourcePath, "input_contracts", manifest.InputContracts, diagnostics) {
		entry.InputContracts = append(entry.InputContracts, registryContract(sourcePath, fmt.Sprintf("input_contracts[%d]", i), contract, diagnostics))
	}
	for i, contract := range boundedRegistryContracts(sourcePath, "output_contracts", manifest.OutputContracts, diagnostics) {
		entry.OutputContracts = append(entry.OutputContracts, registryContract(sourcePath, fmt.Sprintf("output_contracts[%d]", i), contract, diagnostics))
	}
	for i, tc := range boundedRegistryCases(sourcePath, manifest.ConformanceCases, diagnostics) {
		prefix := fmt.Sprintf("conformance_cases[%d]", i)
		entry.ConformanceCases = append(entry.ConformanceCases, RegistryConformanceCase{
			ID:             boundedRegistryText(sourcePath, tc.ID, prefix+".id", diagnostics),
			Description:    boundedRegistryText(sourcePath, tc.Description, prefix+".description", diagnostics),
			FixturePath:    boundedRegistryText(sourcePath, tc.FixturePath, prefix+".fixture_path", diagnostics),
			ExpectedResult: boundedRegistryText(sourcePath, tc.ExpectedResult, prefix+".expected_result", diagnostics),
		})
	}
	return entry
}

func registryContract(sourcePath string, prefix string, contract Contract, diagnostics *[]RegistryDiagnostic) RegistryContract {
	return RegistryContract{
		Name:           boundedRegistryText(sourcePath, contract.Name, prefix+".name", diagnostics),
		Description:    boundedRegistryText(sourcePath, contract.Description, prefix+".description", diagnostics),
		Schema:         boundedRegistryText(sourcePath, contract.Schema, prefix+".schema", diagnostics),
		MediaTypes:     boundedRegistryTextList(sourcePath, prefix+".media_types", contract.MediaTypes, diagnostics),
		RequiredFields: boundedRegistryTextList(sourcePath, prefix+".required_fields", contract.RequiredFields, diagnostics),
		Produces:       boundedRegistryTextList(sourcePath, prefix+".produces", contract.Produces, diagnostics),
	}
}

func boundedRegistryContracts(sourcePath string, field string, contracts []Contract, diagnostics *[]RegistryDiagnostic) []Contract {
	if len(contracts) <= maxRegistryItems {
		return contracts
	}
	*diagnostics = append(*diagnostics, registryLimitDiagnostic(sourcePath, field))
	return contracts[:maxRegistryItems]
}

func boundedRegistryCases(sourcePath string, cases []ConformanceCase, diagnostics *[]RegistryDiagnostic) []ConformanceCase {
	if len(cases) <= maxRegistryItems {
		return cases
	}
	*diagnostics = append(*diagnostics, registryLimitDiagnostic(sourcePath, "conformance_cases"))
	return cases[:maxRegistryItems]
}

func boundedRegistryTextList(sourcePath string, field string, values []string, diagnostics *[]RegistryDiagnostic) []string {
	if len(values) > maxRegistryItems {
		*diagnostics = append(*diagnostics, registryLimitDiagnostic(sourcePath, field))
		values = values[:maxRegistryItems]
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, boundedRegistryText(sourcePath, value, field, diagnostics))
	}
	return out
}

func boundedRegistryText(sourcePath string, value string, field string, diagnostics *[]RegistryDiagnostic) string {
	if len(value) <= maxRegistryTextBytes {
		return value
	}
	*diagnostics = append(*diagnostics, registryLimitDiagnostic(sourcePath, field))
	return value[:maxRegistryTextBytes]
}

func registryLimitDiagnostic(sourcePath string, field string) RegistryDiagnostic {
	return RegistryDiagnostic{
		Level:   "warning",
		Code:    "registry_entry_bounded",
		Path:    sourcePath,
		Message: field + " was bounded in the connector registry entry",
	}
}
