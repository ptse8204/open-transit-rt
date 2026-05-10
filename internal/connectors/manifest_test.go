package connectors

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidSyntheticManifestsCoverSupportedConnectorTypes(t *testing.T) {
	fixtures := map[string]string{
		TypeTelemetrySource:   "valid/valid-telemetry-source.json",
		TypePrediction:        "valid/valid-prediction.json",
		TypeValidator:         "valid/valid-validator.json",
		TypeMonitoringExport:  "valid/valid-monitoring-export.json",
		TypeConsumerDiscovery: "valid/valid-consumer-discovery.json",
	}
	for connectorType, fixture := range fixtures {
		t.Run(connectorType, func(t *testing.T) {
			manifest := loadFixture(t, fixture)
			if manifest.ConnectorType != connectorType {
				t.Fatalf("type = %q, want %q", manifest.ConnectorType, connectorType)
			}
		})
	}
}

func TestSchemaVersionAndConnectorTypeAreExact(t *testing.T) {
	manifest := validManifest(t)
	manifest.SchemaVersion = "open-transit-rt.connector.v2"
	assertViolation(t, manifest.Validate(), "schema_version")

	manifest = validManifest(t)
	manifest.ConnectorType = "fare_payment"
	assertViolation(t, manifest.Validate(), "type")
}

func TestRequiredManifestSections(t *testing.T) {
	manifest := validManifest(t)
	manifest.ConnectorID = ""
	manifest.DisplayName = ""
	manifest.Description = ""
	manifest.Mode = Mode{}
	manifest.InputContracts = nil
	manifest.OutputContracts = nil
	manifest.FailureBehavior = FailureBehavior{}
	manifest.RedactionPolicy = RedactionPolicy{}
	manifest.ClaimBoundary = ClaimBoundary{}
	manifest.DocsLink = ""
	manifest.ConformanceCases = nil

	violations := manifest.Validate()
	for _, field := range []string{
		"connector_id",
		"display_name",
		"description",
		"mode.name",
		"input_contracts",
		"output_contracts",
		"failure_behavior.timeout_seconds",
		"redaction_policy.secret_storage",
		"claim_boundary.positive_claims",
		"docs_link",
		"conformance_cases",
	} {
		assertViolation(t, violations, field)
	}
}

func TestDecodeRejectsSecretsAndPrivatePaths(t *testing.T) {
	for _, tc := range []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "secret key",
			raw:  `{"schema_version":"open-transit-rt.connector.v1","api_key":"sk-1234567890123456"}`,
			want: "api_key",
		},
		{
			name: "secret value",
			raw:  strings.Replace(string(readFixture(t, "valid/valid-telemetry-source.json")), `"description": "Synthetic telemetry source manifest for dry-run GPS adapter conformance."`, `"description": "Bearer abcdefghijklmnop"`, 1),
			want: "description",
		},
		{
			name: "private path",
			raw:  strings.Replace(string(readFixture(t, "valid/valid-telemetry-source.json")), `"testdata/connectors/valid/synthetic-telemetry-source-input.json"`, `"/Users/example/private/input.json"`, 1),
			want: "conformance_cases[0].fixture_path",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecodeManifest(strings.NewReader(tc.raw))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("DecodeManifest error = %v, want field %s", err, tc.want)
			}
		})
	}
}

func TestRejectsUnsafeURLs(t *testing.T) {
	for _, docsLink := range []string{
		"http://example.org/connectors/demo",
		"https://user@example.org/connectors/demo",
		"https://localhost/connectors/demo",
		"https://10.0.0.5/connectors/demo",
		"../docs/connectors/demo.md",
	} {
		t.Run(docsLink, func(t *testing.T) {
			manifest := validManifest(t)
			manifest.DocsLink = docsLink
			assertViolation(t, manifest.Validate(), "docs_link")
		})
	}
}

func TestRejectsUnsupportedPositiveClaims(t *testing.T) {
	for _, claim := range []string{
		"Caltrans compliant",
		"production ready",
		"real vendor compatibility",
		"custom positive claim",
	} {
		t.Run(claim, func(t *testing.T) {
			manifest := validManifest(t)
			manifest.ClaimBoundary.PositiveClaims = []string{claim}
			assertViolation(t, manifest.Validate(), "claim_boundary.positive_claims[0]")
		})
	}
}

func TestRejectsRawValidatorCommands(t *testing.T) {
	manifest := loadFixture(t, "valid/valid-validator.json")
	manifest.OutputContracts[0].RawValidatorCommand = []string{"java", "-jar", "/tmp/validator.jar"}
	assertViolation(t, manifest.Validate(), "output_contracts[0].raw_validator_command")

	raw := strings.Replace(string(readFixture(t, "valid/valid-validator.json")), `"produces": [`, `"command": "java -jar validator.jar", "produces": [`, 1)
	_, err := DecodeManifest(strings.NewReader(raw))
	if err == nil || !strings.Contains(err.Error(), "command") {
		t.Fatalf("DecodeManifest raw command error = %v", err)
	}
}

func TestRejectsNotificationSubmissionAndStatusMutation(t *testing.T) {
	manifest := loadFixture(t, "valid/valid-monitoring-export.json")
	manifest.Mode.SendsNotificationsByDefault = true
	assertViolation(t, manifest.Validate(), "mode.sends_notifications_by_default")

	manifest = loadFixture(t, "valid/valid-consumer-discovery.json")
	manifest.Mode.AutomatesConsumerSubmission = true
	assertViolation(t, manifest.Validate(), "mode.automates_consumer_submission")

	manifest = validManifest(t)
	manifest.OutputContracts[0].MutatesStatus = true
	assertViolation(t, manifest.Validate(), "output_contracts[0].mutates_status")
}

func TestConformanceCasesMustUseSyntheticFixtureScope(t *testing.T) {
	for _, fixturePath := range []string{
		"../testdata/connectors/input.json",
		"testdata/../connectors/input.json",
		"testdata/avl-vendor/valid.json",
		"/tmp/input.json",
	} {
		t.Run(fixturePath, func(t *testing.T) {
			manifest := validManifest(t)
			manifest.ConformanceCases[0].FixturePath = fixturePath
			assertViolation(t, manifest.Validate(), "conformance_cases[0].fixture_path")
		})
	}
}

func TestInvalidFixturesAreRejected(t *testing.T) {
	fixtures, err := filepath.Glob(filepath.Join("..", "..", "testdata", "connectors", "invalid", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(fixtures) == 0 {
		t.Fatal("expected invalid connector fixtures")
	}
	for _, fixture := range fixtures {
		t.Run(filepath.Base(fixture), func(t *testing.T) {
			raw, err := os.ReadFile(fixture)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeManifest(strings.NewReader(string(raw))); err == nil {
				t.Fatalf("DecodeManifest(%s) succeeded, want rejection", fixture)
			}
		})
	}
}

func validManifest(t *testing.T) Manifest {
	t.Helper()
	return loadFixture(t, "valid/valid-telemetry-source.json")
}

func loadFixture(t *testing.T, name string) Manifest {
	t.Helper()
	manifest, err := DecodeManifest(strings.NewReader(string(readFixture(t, name))))
	if err != nil {
		t.Fatalf("load %s: %v", name, err)
	}
	return manifest
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata", "connectors", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return raw
}

func assertViolation(t *testing.T, violations []Violation, field string) {
	t.Helper()
	for _, violation := range violations {
		if violation.Field == field {
			return
		}
	}
	t.Fatalf("violation for %s not found in %+v", field, violations)
}
