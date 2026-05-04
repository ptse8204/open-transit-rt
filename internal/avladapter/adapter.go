package avladapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"

	"open-transit-rt/internal/telemetry"
)

const (
	SeverityError   = "error"
	SeverityWarning = "warning"

	CodeDuplicateMapping      = "duplicate_mapping"
	CodeEmptyMappedIdentifier = "empty_mapped_identifier"
	CodeInvalidMappingJSON    = "invalid_mapping_json"
	CodeInvalidPayloadJSON    = "invalid_payload_json"
	CodeSourceMismatch        = "vendor_source_mismatch"
	CodeUnknownMapping        = "unknown_vendor_mapping"
	CodeMissingRequiredField  = "missing_required_field"
	CodeInvalidCoordinate     = "invalid_coordinate"
	CodeInvalidTelemetryEvent = "invalid_telemetry_event"
	CodeStaleTimestamp        = "stale_timestamp"
	CodeFutureTimestamp       = "future_timestamp"
	CodeLowGPSAccuracy        = "low_gps_accuracy"
	CodeDuplicateObservation  = "duplicate_observation"
	CodeOutOfOrderObservation = "out_of_order_observation"
)

const (
	DefaultStaleThreshold    = 90 * time.Second
	DefaultFutureThreshold   = 30 * time.Second
	DefaultLowAccuracyMeters = 50.0
)

type Diagnostic struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Index    *int   `json:"index,omitempty"`
}

type MappingFile struct {
	Mappings []MappingRow `json:"mappings"`
}

type MappingRow struct {
	VendorSource    string `json:"vendor_source"`
	VendorDeviceID  string `json:"vendor_device_id"`
	VendorVehicleID string `json:"vendor_vehicle_id"`
	AgencyID        string `json:"agency_id"`
	DeviceID        string `json:"device_id"`
	VehicleID       string `json:"vehicle_id"`
	Notes           string `json:"notes,omitempty"`
}

type Payload struct {
	VendorSource string         `json:"vendor_source"`
	Records      []VendorRecord `json:"records"`
}

type VendorRecord struct {
	VendorDeviceID  string   `json:"vendor_device_id"`
	VendorVehicleID string   `json:"vendor_vehicle_id"`
	ObservedAt      string   `json:"observed_at"`
	Lat             *float64 `json:"lat"`
	Lon             *float64 `json:"lon"`
	Bearing         *float64 `json:"bearing,omitempty"`
	SpeedMPS        *float64 `json:"speed_mps,omitempty"`
	AccuracyM       *float64 `json:"accuracy_m,omitempty"`
	TripHint        string   `json:"trip_hint,omitempty"`
}

type Options struct {
	ReferenceTime     time.Time
	StaleThreshold    time.Duration
	FutureThreshold   time.Duration
	LowAccuracyMeters float64
}

type Result struct {
	Events      []telemetry.Event
	Diagnostics []Diagnostic
}

type mapper struct {
	rows map[mappingKey]MappingRow
}

type mappingKey struct {
	vendorSource    string
	vendorDeviceID  string
	vendorVehicleID string
}

type vehicleKey struct {
	agencyID  string
	vehicleID string
}

func LoadMapping(r io.Reader) (MappingFile, []Diagnostic) {
	var mapping MappingFile
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&mapping); err != nil {
		return MappingFile{}, []Diagnostic{errorDiagnostic(CodeInvalidMappingJSON, fmt.Sprintf("invalid mapping JSON: %v", err), nil)}
	}
	return mapping, validateMapping(mapping)
}

func TransformPayload(payloadReader io.Reader, mapping MappingFile, options Options) Result {
	options = options.withDefaults()
	m, diagnostics := newMapper(mapping)
	if hasErrors(diagnostics) {
		return Result{Diagnostics: diagnostics}
	}

	var payload Payload
	decoder := json.NewDecoder(payloadReader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidPayloadJSON, fmt.Sprintf("invalid vendor payload JSON: %v", err), nil))
		return Result{Diagnostics: diagnostics}
	}

	events := make([]telemetry.Event, 0, len(payload.Records))
	latestByVehicle := map[vehicleKey]time.Time{}
	for index, record := range payload.Records {
		idx := index
		if payload.VendorSource == "" {
			diagnostics = append(diagnostics, errorDiagnostic(CodeMissingRequiredField, "vendor_source is required", &idx))
			continue
		}
		if record.VendorDeviceID == "" || record.VendorVehicleID == "" || record.ObservedAt == "" {
			diagnostics = append(diagnostics, errorDiagnostic(CodeMissingRequiredField, "vendor_device_id, vendor_vehicle_id, and observed_at are required", &idx))
			continue
		}
		if record.Lat == nil || record.Lon == nil {
			diagnostics = append(diagnostics, errorDiagnostic(CodeMissingRequiredField, "lat and lon are required", &idx))
			continue
		}

		row, ok := m.lookup(payload.VendorSource, record.VendorDeviceID, record.VendorVehicleID)
		if !ok {
			if m.hasVendorIDs(record.VendorDeviceID, record.VendorVehicleID) {
				diagnostics = append(diagnostics, errorDiagnostic(CodeSourceMismatch, "vendor_source does not match the mapping row for this vendor device and vehicle", &idx))
			} else {
				diagnostics = append(diagnostics, errorDiagnostic(CodeUnknownMapping, "no mapping exists for vendor_source, vendor_device_id, and vendor_vehicle_id", &idx))
			}
			continue
		}
		if *record.Lat < -90 || *record.Lat > 90 || *record.Lon < -180 || *record.Lon > 180 {
			diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidCoordinate, "lat must be between -90 and 90 and lon must be between -180 and 180", &idx))
			continue
		}
		observedAt, err := time.Parse(time.RFC3339, record.ObservedAt)
		if err != nil {
			diagnostics = append(diagnostics, errorDiagnostic(CodeMissingRequiredField, "observed_at must be RFC3339 with timezone", &idx))
			continue
		}

		event := telemetry.Event{
			AgencyID:  row.AgencyID,
			DeviceID:  row.DeviceID,
			VehicleID: row.VehicleID,
			Timestamp: observedAt,
			Lat:       *record.Lat,
			Lon:       *record.Lon,
			TripHint:  record.TripHint,
		}
		if record.Bearing != nil {
			event.Bearing = *record.Bearing
		}
		if record.SpeedMPS != nil {
			event.SpeedMPS = *record.SpeedMPS
		}
		if record.AccuracyM != nil {
			event.AccuracyM = *record.AccuracyM
		}
		if !roundTripsAsTelemetryEvent(event) || !event.Valid() {
			diagnostics = append(diagnostics, errorDiagnostic(CodeInvalidTelemetryEvent, "transformed record does not satisfy the Open Transit RT telemetry contract", &idx))
			continue
		}

		if observedAt.Before(options.ReferenceTime.Add(-options.StaleThreshold)) {
			diagnostics = append(diagnostics, warningDiagnostic(CodeStaleTimestamp, "observed_at is stale relative to reference time", &idx))
		}
		if observedAt.After(options.ReferenceTime.Add(options.FutureThreshold)) {
			diagnostics = append(diagnostics, warningDiagnostic(CodeFutureTimestamp, "observed_at is in the future relative to reference time", &idx))
		}
		if record.AccuracyM != nil && *record.AccuracyM > options.LowAccuracyMeters {
			diagnostics = append(diagnostics, warningDiagnostic(CodeLowGPSAccuracy, "accuracy_m is above the dry-run review threshold", &idx))
		}
		key := vehicleKey{agencyID: event.AgencyID, vehicleID: event.VehicleID}
		if previous, ok := latestByVehicle[key]; ok {
			switch {
			case observedAt.Equal(previous):
				diagnostics = append(diagnostics, warningDiagnostic(CodeDuplicateObservation, "record has the same observed_at as an earlier dry-run record for this mapped vehicle", &idx))
			case observedAt.Before(previous):
				diagnostics = append(diagnostics, warningDiagnostic(CodeOutOfOrderObservation, "record is older than an earlier dry-run record for this mapped vehicle", &idx))
			}
		}
		if previous, ok := latestByVehicle[key]; !ok || observedAt.After(previous) {
			latestByVehicle[key] = observedAt
		}
		events = append(events, event)
	}
	return Result{Events: events, Diagnostics: diagnostics}
}

func MarshalEvents(events []telemetry.Event) ([]byte, error) {
	if events == nil {
		events = []telemetry.Event{}
	}
	return json.MarshalIndent(events, "", "  ")
}

func MarshalDiagnostics(diagnostics []Diagnostic) ([]byte, error) {
	if diagnostics == nil {
		diagnostics = []Diagnostic{}
	}
	sort.SliceStable(diagnostics, func(i, j int) bool {
		if diagnostics[i].Index == nil && diagnostics[j].Index != nil {
			return false
		}
		if diagnostics[i].Index != nil && diagnostics[j].Index == nil {
			return true
		}
		if diagnostics[i].Index != nil && diagnostics[j].Index != nil && *diagnostics[i].Index != *diagnostics[j].Index {
			return *diagnostics[i].Index < *diagnostics[j].Index
		}
		if diagnostics[i].Severity != diagnostics[j].Severity {
			return diagnostics[i].Severity < diagnostics[j].Severity
		}
		return diagnostics[i].Code < diagnostics[j].Code
	})
	return json.MarshalIndent(diagnostics, "", "  ")
}

func HasHardErrors(diagnostics []Diagnostic) bool {
	return hasErrors(diagnostics)
}

func validateMapping(mapping MappingFile) []Diagnostic {
	_, diagnostics := newMapper(mapping)
	return diagnostics
}

func newMapper(mapping MappingFile) (mapper, []Diagnostic) {
	diagnostics := []Diagnostic{}
	rows := map[mappingKey]MappingRow{}
	for index, row := range mapping.Mappings {
		idx := index
		if row.VendorSource == "" || row.VendorDeviceID == "" || row.VendorVehicleID == "" {
			diagnostics = append(diagnostics, errorDiagnostic(CodeMissingRequiredField, "vendor_source, vendor_device_id, and vendor_vehicle_id are required in mapping rows", &idx))
			continue
		}
		if row.AgencyID == "" || row.DeviceID == "" || row.VehicleID == "" {
			diagnostics = append(diagnostics, errorDiagnostic(CodeEmptyMappedIdentifier, "agency_id, device_id, and vehicle_id are required in mapping rows", &idx))
			continue
		}
		key := mappingKey{vendorSource: row.VendorSource, vendorDeviceID: row.VendorDeviceID, vendorVehicleID: row.VendorVehicleID}
		if _, exists := rows[key]; exists {
			diagnostics = append(diagnostics, errorDiagnostic(CodeDuplicateMapping, "duplicate mapping row for vendor_source, vendor_device_id, and vendor_vehicle_id", &idx))
			continue
		}
		rows[key] = row
	}
	return mapper{rows: rows}, diagnostics
}

func (m mapper) lookup(source string, deviceID string, vehicleID string) (MappingRow, bool) {
	row, ok := m.rows[mappingKey{vendorSource: source, vendorDeviceID: deviceID, vendorVehicleID: vehicleID}]
	return row, ok
}

func (m mapper) hasVendorIDs(deviceID string, vehicleID string) bool {
	for key := range m.rows {
		if key.vendorDeviceID == deviceID && key.vendorVehicleID == vehicleID {
			return true
		}
	}
	return false
}

func (o Options) withDefaults() Options {
	if o.ReferenceTime.IsZero() {
		o.ReferenceTime = time.Now().UTC()
	}
	if o.StaleThreshold == 0 {
		o.StaleThreshold = DefaultStaleThreshold
	}
	if o.FutureThreshold == 0 {
		o.FutureThreshold = DefaultFutureThreshold
	}
	if o.LowAccuracyMeters == 0 {
		o.LowAccuracyMeters = DefaultLowAccuracyMeters
	}
	return o
}

func roundTripsAsTelemetryEvent(event telemetry.Event) bool {
	raw, err := json.Marshal(event)
	if err != nil {
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var decoded telemetry.Event
	if err := decoder.Decode(&decoded); err != nil {
		return false
	}
	return decoded.Valid()
}

func errorDiagnostic(code string, message string, index *int) Diagnostic {
	return Diagnostic{Code: code, Severity: SeverityError, Message: message, Index: index}
}

func warningDiagnostic(code string, message string, index *int) Diagnostic {
	return Diagnostic{Code: code, Severity: SeverityWarning, Message: message, Index: index}
}

func hasErrors(diagnostics []Diagnostic) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == SeverityError {
			return true
		}
	}
	return false
}
