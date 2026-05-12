package main

import (
	"net/http"
	"time"

	"open-transit-rt/internal/auth"
)

type connectorTestsView struct {
	GeneratedAt time.Time                `json:"generated_at"`
	AgencyID    string                   `json:"agency_id"`
	Boundary    string                   `json:"boundary"`
	Commands    []connectorTestCommand   `json:"commands"`
	ClaimFlags  connectorTestsClaimFlags `json:"claim_flags"`
}

type connectorTestCommand struct {
	ID                string   `json:"id"`
	Label             string   `json:"label"`
	CommandLine       string   `json:"command_line"`
	Validates         string   `json:"validates"`
	Inputs            string   `json:"inputs"`
	FailureNextAction string   `json:"failure_next_action"`
	DoesNotProve      string   `json:"does_not_prove"`
	DocsLinks         []string `json:"docs_links"`
}

type connectorTestsClaimFlags struct {
	BackendCommandExecutionEnabled  bool `json:"backend_command_execution_enabled"`
	ManifestCommandExecutionEnabled bool `json:"manifest_command_execution_enabled"`
	ExternalNetworkContacted        bool `json:"external_network_contacted"`
	ExternalEvidenceCreated         bool `json:"external_evidence_created"`
	ConsumerStatusesChanged         bool `json:"consumer_statuses_changed"`
	ComplianceClaimed               bool `json:"compliance_claimed"`
	VendorCompatibilityClaimed      bool `json:"vendor_compatibility_claimed"`
	ProductionReadinessClaimed      bool `json:"production_readiness_claimed"`
	ProductionGradeETAClaimed       bool `json:"production_grade_eta_claimed"`
}

func (h *handler) renderConnectorTests(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "connector-tests")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "connector-tests", page)
}

func (h *handler) renderConnectorTestsJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "connector-tests")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.ConnectorTests)
}

func buildConnectorTests(page operationsPage) connectorTestsView {
	return connectorTestsView{
		GeneratedAt: page.GeneratedAt,
		AgencyID:    page.AgencyID,
		Boundary:    "Private authenticated connector test instructions only; this page generates fixed local/offline commands for an operator to run outside the web request. Viewing it executes nothing, writes no evidence, contacts no external party, changes no consumer status, and records no compatibility, compliance, production-readiness, or ETA-quality outcome.",
		Commands: []connectorTestCommand{
			connectorTestCommandView(
				"manifest-and-examples",
				"Manifest and example safety check",
				"make external-connection-check",
				"Connector manifests decode, example manifests remain sidecar/manifest/conformance bounded, and synthetic connector examples still pass local tests.",
				"Committed manifests under testdata/connectors and examples/connectors plus synthetic fixtures.",
				"Review the failing manifest, fixture, or example test. Do not bypass schema validation or add real credentials.",
				"Compliance, consumer acceptance, vendor compatibility, production readiness, public launch, or ETA quality.",
				[]string{"docs/connectors/plugin-contract.md", "docs/integration-adapter-kit.md"},
			),
			connectorTestCommandView(
				"adapter-conformance-full",
				"Full offline adapter conformance suite",
				"make adapter-conformance",
				"Telemetry, prediction, validator, and monitoring conformance cases are present, synthetic, and fail closed.",
				"testdata/adapter-conformance/suite.json and fixtures under testdata/adapter-conformance/fixtures.",
				"Open the suite case named in the failure and fix the synthetic fixture, manifest reference, or adapter boundary expectation.",
				"External network behavior, real validator execution, consumer submission, real vendor compatibility, production AVL reliability, or production-grade ETA quality.",
				[]string{"docs/tutorials/external-adapter-conformance.md", "docs/connectors/plugin-contract.md"},
			),
			connectorTestCommandView(
				"adapter-conformance-manifest",
				"Manifest-only conformance check",
				"go run ./cmd/adapter-conformance manifest --suite testdata/adapter-conformance",
				"Suite-referenced connector manifests decode through the safe manifest contract.",
				"Committed synthetic suite manifest paths only.",
				"Fix manifest fields through open-transit-rt.connector.v1; do not add raw commands, secrets, or stronger positive claims.",
				"Runtime plugin loading, validator execution, external contact, compliance, or compatibility.",
				[]string{"docs/connectors/plugin-contract.md"},
			),
			connectorTestCommandView(
				"adapter-conformance-telemetry",
				"Telemetry connector cases",
				"go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance",
				"Malformed, stale, future, wrong-agency, unknown-device, low-quality, duplicate, and out-of-order telemetry cases remain fail-closed.",
				"Synthetic telemetry conformance fixtures.",
				"Fix the dry-run adapter shape or fixture expectation. Keep real device credentials and private AVL payloads out of the repo.",
				"Real fleet reliability, hardware certification, real vendor support, agency adoption, or production AVL quality.",
				[]string{"docs/tutorials/external-adapter-conformance.md", "docs/tutorials/device-avl-integration.md"},
			),
			connectorTestCommandView(
				"adapter-conformance-prediction",
				"Prediction connector cases",
				"go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance",
				"Timeout, malformed, stale, wrong-agency, and low-confidence prediction outputs remain withheld or fail-closed.",
				"Synthetic prediction conformance fixtures.",
				"Fix sanitized request/response handling or diagnostics while keeping Vehicle Positions independent of predictor availability.",
				"Named predictor compatibility, consumer acceptance, production-grade ETA quality, or public Trip Updates mutation.",
				[]string{"docs/requirements-trip-updates.md", "docs/integration-adapter-kit.md"},
			),
			connectorTestCommandView(
				"adapter-conformance-validator",
				"Validator connector cases",
				"go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance",
				"Validator connector fixtures use allowlisted validator IDs and avoid raw validator commands.",
				"Synthetic validator allowlist fixture.",
				"Map validator behavior to server-owned IDs and normalized reports. Do not add arbitrary argument lists or private paths.",
				"Validator-clean feeds, CAL-ITP/Caltrans compliance, consumer acceptance, or public launch readiness.",
				[]string{"docs/tutorials/gtfs-validation-triage.md", "docs/dependencies.md"},
			),
			connectorTestCommandView(
				"adapter-conformance-monitoring",
				"Monitoring/export connector cases",
				"go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance",
				"Monitoring/export fixtures keep redaction and no-send defaults.",
				"Synthetic monitoring redaction and no-send fixtures.",
				"Fix redaction or no-send defaults. Keep webhook destinations, tokens, and external notification sends out of the example.",
				"SLA coverage, uptime guarantees, hosted service availability, production readiness, or retained evidence.",
				[]string{"docs/tutorials/self-hosted-operations-notifications.md", "docs/runbooks/small-agency-pilot-operations.md"},
			),
			connectorTestCommandView(
				"connector-example-tests",
				"Synthetic connector example tests",
				"make test-connector-examples",
				"All committed synthetic connector examples still compile and pass their local dry-run tests.",
				"Go examples under examples/connectors and their committed synthetic fixtures.",
				"Fix the example or fixture without adding network sends, private payloads, real credentials, or vendor-specific claims.",
				"Real integration proof, vendor compatibility, production reliability, consumer acceptance, or evidence collection.",
				[]string{"examples/README.md", "docs/integration-adapter-kit.md"},
			),
		},
		ClaimFlags: connectorTestsClaimFlags{},
	}
}

func connectorTestCommandView(id string, label string, commandLine string, validates string, inputs string, failureNextAction string, doesNotProve string, docsLinks []string) connectorTestCommand {
	return connectorTestCommand{
		ID:                firstNonEmpty(id, "connector-check"),
		Label:             firstNonEmpty(label, "Connector check"),
		CommandLine:       firstNonEmpty(commandLine, "make external-connection-check"),
		Validates:         firstNonEmpty(validates, "Local connector safety checks."),
		Inputs:            firstNonEmpty(inputs, "Committed synthetic fixtures only."),
		FailureNextAction: firstNonEmpty(failureNextAction, "Review the failing local fixture or manifest."),
		DoesNotProve:      firstNonEmpty(doesNotProve, "Compliance, compatibility, production readiness, or consumer acceptance."),
		DocsLinks:         safeDocsLinks(docsLinks),
	}
}
