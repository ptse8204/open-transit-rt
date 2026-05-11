package predictionsdk

const (
	ReasonEligibleNoop          = "eligible but stub does not emit public Trip Updates"
	ReasonBelowConfidence       = "below confidence threshold"
	ReasonMissingTelemetry      = "no telemetry for assignment"
	ReasonMissingTripDescriptor = "missing trip descriptor"
)

type Request struct {
	SyntheticOnly       bool         `json:"synthetic_only"`
	ActiveFeedVersion   string       `json:"active_feed_version"`
	VehiclePositionsURL string       `json:"vehicle_positions_url"`
	Telemetry           []Telemetry  `json:"telemetry"`
	Assignments         []Assignment `json:"assignments"`
}

type Telemetry struct {
	AgencyID   string  `json:"agency_id"`
	VehicleID  string  `json:"vehicle_id"`
	ObservedAt string  `json:"observed_at"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type Assignment struct {
	VehicleID  string  `json:"vehicle_id"`
	TripID     string  `json:"trip_id"`
	StartDate  string  `json:"start_date"`
	StartTime  string  `json:"start_time"`
	Confidence float64 `json:"confidence"`
}

type Options struct {
	ConfidenceThreshold float64
	Mode                string
}

type Response struct {
	Mode                      string       `json:"mode"`
	NetworkSend               bool         `json:"network_send"`
	PublicTripUpdatesMutation bool         `json:"public_trip_updates_mutation"`
	TripUpdates               []TripUpdate `json:"trip_updates"`
	Diagnostics               []Diagnostic `json:"diagnostics"`
	VehiclePositionsReference string       `json:"vehicle_positions_reference"`
	ActiveFeedVersion         string       `json:"active_feed_version"`
	Counts                    Counts       `json:"counts"`
}

type TripUpdate struct {
	TripID    string `json:"trip_id"`
	StartDate string `json:"start_date"`
	StartTime string `json:"start_time"`
}

type Diagnostic struct {
	VehicleID string `json:"vehicle_id"`
	State     string `json:"state"`
	Reason    string `json:"reason"`
}

type Counts struct {
	TelemetryRows     int `json:"telemetry_rows"`
	Assignments       int `json:"assignments"`
	EligibleNoop      int `json:"eligible_noop"`
	WithheldUnknown   int `json:"withheld_unknown"`
	MissingTelemetry  int `json:"missing_telemetry"`
	MissingDescriptor int `json:"missing_descriptor"`
}

func BuildDryRunResponse(request Request, opts Options) Response {
	opts = normalizeOptions(opts)
	response := Response{
		Mode:                      opts.Mode,
		NetworkSend:               false,
		PublicTripUpdatesMutation: false,
		VehiclePositionsReference: request.VehiclePositionsURL,
		ActiveFeedVersion:         request.ActiveFeedVersion,
		Counts: Counts{
			TelemetryRows: len(request.Telemetry),
			Assignments:   len(request.Assignments),
		},
	}
	telemetryByVehicle := make(map[string]bool, len(request.Telemetry))
	for _, row := range request.Telemetry {
		if row.VehicleID != "" {
			telemetryByVehicle[row.VehicleID] = true
		}
	}
	for _, assignment := range request.Assignments {
		switch {
		case assignment.TripID == "" || assignment.StartDate == "" || assignment.StartTime == "":
			response.Counts.MissingDescriptor++
			response.Diagnostics = append(response.Diagnostics, Diagnostic{VehicleID: assignment.VehicleID, State: "withheld", Reason: ReasonMissingTripDescriptor})
		case assignment.Confidence < opts.ConfidenceThreshold:
			response.Counts.WithheldUnknown++
			response.Diagnostics = append(response.Diagnostics, Diagnostic{VehicleID: assignment.VehicleID, State: "unknown", Reason: ReasonBelowConfidence})
		case !telemetryByVehicle[assignment.VehicleID]:
			response.Counts.MissingTelemetry++
			response.Diagnostics = append(response.Diagnostics, Diagnostic{VehicleID: assignment.VehicleID, State: "withheld", Reason: ReasonMissingTelemetry})
		default:
			response.Counts.EligibleNoop++
			response.Diagnostics = append(response.Diagnostics, Diagnostic{VehicleID: assignment.VehicleID, State: "eligible", Reason: ReasonEligibleNoop})
		}
	}
	return response
}

func normalizeOptions(opts Options) Options {
	if opts.ConfidenceThreshold == 0 {
		opts.ConfidenceThreshold = 0.8
	}
	if opts.Mode == "" {
		opts.Mode = "dry_run_noop"
	}
	return opts
}
