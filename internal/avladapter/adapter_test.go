package avladapter

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"open-transit-rt/internal/telemetry"
)

var referenceTime = time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC)

func TestTransformValidPayloadUsesMappingAuthorityAndTelemetryContract(t *testing.T) {
	result := transformFixture(t, "mapping.json", "valid.json")
	if HasHardErrors(result.Diagnostics) {
		t.Fatalf("diagnostics contain hard errors: %+v", result.Diagnostics)
	}
	if len(result.Events) != 1 {
		t.Fatalf("events len = %d, want 1", len(result.Events))
	}
	event := result.Events[0]
	if event.AgencyID != "demo-agency" || event.DeviceID != "device-1" || event.VehicleID != "bus-1" {
		t.Fatalf("mapped ids = %s/%s/%s", event.AgencyID, event.DeviceID, event.VehicleID)
	}
	if !roundTripsAsTelemetryEvent(event) || !event.Valid() {
		t.Fatalf("event does not satisfy telemetry contract: %+v", event)
	}
	raw, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	var decoded telemetry.Event
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal event: %v", err)
	}
	if !decoded.Valid() {
		t.Fatalf("decoded event invalid: %+v", decoded)
	}
}

func TestTransformDoesNotAllowPayloadToOverrideMappedIdentifiers(t *testing.T) {
	payload := `{
	  "vendor_source": "vendor-demo",
	  "records": [{
	    "vendor_device_id": "vendor-device-1",
	    "vendor_vehicle_id": "vendor-vehicle-1",
	    "agency_id": "wrong-agency",
	    "device_id": "wrong-device",
	    "vehicle_id": "wrong-bus",
	    "observed_at": "2026-05-04T12:00:00Z",
	    "lat": 49.2827,
	    "lon": -123.1207
	  }]
	}`
	mapping := loadMappingFixture(t, "mapping.json")
	result := TransformPayload(strings.NewReader(payload), mapping, Options{ReferenceTime: referenceTime})
	assertDiagnostic(t, result.Diagnostics, CodeInvalidPayloadJSON, SeverityError)
	if len(result.Events) != 0 {
		t.Fatalf("events len = %d, want 0 for payload with forbidden identifier fields", len(result.Events))
	}
}

func TestMappingHardErrors(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mapping     string
		wantCode    string
		wantMessage string
	}{
		{name: "duplicate mapping rows", mapping: "duplicate-mapping.json", wantCode: CodeDuplicateMapping},
		{name: "empty mapped ids", mapping: "empty-mapped-ids.json", wantCode: CodeEmptyMappedIdentifier},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, diagnostics := LoadMapping(bytes.NewReader(readFixture(t, tc.mapping)))
			assertDiagnostic(t, diagnostics, tc.wantCode, SeverityError)
		})
	}
}

func TestMappingRejectsSecretLikeUnknownFields(t *testing.T) {
	raw := `{
	  "mappings": [{
	    "vendor_source": "vendor-demo",
	    "vendor_device_id": "vendor-device-1",
	    "vendor_vehicle_id": "vendor-vehicle-1",
	    "agency_id": "demo-agency",
	    "device_id": "device-1",
	    "vehicle_id": "bus-1",
	    "token": "do-not-accept"
	  }]
	}`
	_, diagnostics := LoadMapping(strings.NewReader(raw))
	assertDiagnostic(t, diagnostics, CodeInvalidMappingJSON, SeverityError)
}

func TestPayloadHardErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		payload string
		code    string
	}{
		{name: "source mismatch", payload: "source-mismatch.json", code: CodeSourceMismatch},
		{name: "missing coordinate", payload: "missing-coordinate.json", code: CodeMissingRequiredField},
		{name: "invalid coordinate", payload: "invalid-coordinate.json", code: CodeInvalidCoordinate},
		{name: "unknown vendor vehicle", payload: "unknown-vendor-vehicle.json", code: CodeUnknownMapping},
		{name: "malformed payload", payload: "malformed.json", code: CodeInvalidPayloadJSON},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := transformFixture(t, "mapping.json", tc.payload)
			assertDiagnostic(t, result.Diagnostics, tc.code, SeverityError)
		})
	}
}

func TestWarningsUseFixedReferenceTime(t *testing.T) {
	for _, tc := range []struct {
		name    string
		payload string
		code    string
	}{
		{name: "stale", payload: "stale-timestamp.json", code: CodeStaleTimestamp},
		{name: "future", payload: "future-timestamp.json", code: CodeFutureTimestamp},
		{name: "low accuracy", payload: "low-gps-accuracy.json", code: CodeLowGPSAccuracy},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := transformFixture(t, "mapping.json", tc.payload)
			assertDiagnostic(t, result.Diagnostics, tc.code, SeverityWarning)
			if HasHardErrors(result.Diagnostics) {
				t.Fatalf("warning-only payload has hard errors: %+v", result.Diagnostics)
			}
			if len(result.Events) != 1 {
				t.Fatalf("events len = %d, want 1", len(result.Events))
			}
		})
	}
}

func TestBatchMixedValidityKeepsValidOutputAndHardErrors(t *testing.T) {
	result := transformFixture(t, "mapping.json", "batch-mixed.json")
	assertDiagnostic(t, result.Diagnostics, CodeUnknownMapping, SeverityError)
	if len(result.Events) != 1 {
		t.Fatalf("events len = %d, want only valid record", len(result.Events))
	}
	if result.Events[0].VehicleID != "bus-1" {
		t.Fatalf("valid event vehicle = %q", result.Events[0].VehicleID)
	}
}

func TestDuplicateAndOutOfOrderAreDryRunWarnings(t *testing.T) {
	result := transformFixture(t, "mapping.json", "duplicate-out-of-order.json")
	assertDiagnostic(t, result.Diagnostics, CodeDuplicateObservation, SeverityWarning)
	assertDiagnostic(t, result.Diagnostics, CodeOutOfOrderObservation, SeverityWarning)
	if HasHardErrors(result.Diagnostics) {
		t.Fatalf("duplicate/out-of-order dry-run observations should not be hard errors: %+v", result.Diagnostics)
	}
}

func TestOptionalTripHintRemainsTelemetryHint(t *testing.T) {
	result := transformFixture(t, "mapping.json", "optional-trip-hint.json")
	if HasHardErrors(result.Diagnostics) {
		t.Fatalf("diagnostics contain hard errors: %+v", result.Diagnostics)
	}
	if len(result.Events) != 1 || result.Events[0].TripHint != "trip-10-0800" {
		t.Fatalf("trip hint not preserved as telemetry field: %+v", result.Events)
	}
}

func TestMarshalOutputsStableJSONArrays(t *testing.T) {
	events, err := MarshalEvents(nil)
	if err != nil {
		t.Fatalf("marshal nil events: %v", err)
	}
	if string(events) != "[]" {
		t.Fatalf("nil events JSON = %q, want []", string(events))
	}
	diagnostics, err := MarshalDiagnostics([]Diagnostic{
		warningDiagnostic(CodeLowGPSAccuracy, "warning", ptr(1)),
		errorDiagnostic(CodeUnknownMapping, "error", ptr(0)),
	})
	if err != nil {
		t.Fatalf("marshal diagnostics: %v", err)
	}
	if !json.Valid(diagnostics) || !strings.Contains(string(diagnostics), `"code"`) || !strings.Contains(string(diagnostics), `"severity"`) || !strings.Contains(string(diagnostics), `"index"`) {
		t.Fatalf("diagnostics JSON missing required fields: %s", diagnostics)
	}
}

func loadMappingFixture(t *testing.T, name string) MappingFile {
	t.Helper()
	mapping, diagnostics := LoadMapping(bytes.NewReader(readFixture(t, name)))
	if HasHardErrors(diagnostics) {
		t.Fatalf("mapping diagnostics: %+v", diagnostics)
	}
	return mapping
}

func transformFixture(t *testing.T, mappingName string, payloadName string) Result {
	t.Helper()
	mapping := loadMappingFixture(t, mappingName)
	return TransformPayload(bytes.NewReader(readFixture(t, payloadName)), mapping, Options{ReferenceTime: referenceTime})
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata", "avl-vendor", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return raw
}

func assertDiagnostic(t *testing.T, diagnostics []Diagnostic, code string, severity string) {
	t.Helper()
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code && diagnostic.Severity == severity {
			return
		}
	}
	t.Fatalf("diagnostic %s/%s not found in %+v", code, severity, diagnostics)
}

func ptr(value int) *int {
	return &value
}
