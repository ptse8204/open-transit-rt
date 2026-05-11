package telemetrysdk

import "time"

const (
	ReasonMissingIdentity   = "missing identity"
	ReasonInvalidTimestamp  = "invalid timestamp"
	ReasonFutureTimestamp   = "future timestamp"
	ReasonStaleObservation  = "stale observation"
	ReasonLowQuality        = "low quality"
	ReasonInvalidCoordinate = "invalid coordinates"
)

type Options struct {
	Now           time.Time
	Source        string
	MaxFutureSkew time.Duration
	MaxAge        time.Duration
	MinQuality    float64
}

type Observation struct {
	AgencyID   string
	DeviceID   string
	VehicleID  string
	ObservedAt time.Time
	Latitude   float64
	Longitude  float64
	Quality    float64
}

type Event struct {
	AgencyID    string  `json:"agency_id"`
	DeviceID    string  `json:"device_id"`
	VehicleID   string  `json:"vehicle_id"`
	ObservedAt  string  `json:"observed_at"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Source      string  `json:"source"`
	DryRun      bool    `json:"dry_run"`
	NetworkSend bool    `json:"network_send"`
}

type Drop struct {
	DeviceID  string `json:"device_id"`
	VehicleID string `json:"vehicle_id,omitempty"`
	Reason    string `json:"reason"`
}

type Summary struct {
	Events []Event `json:"events"`
	Drops  []Drop  `json:"drops"`
}

func NormalizeBatch(observations []Observation, opts Options) Summary {
	opts = normalizeOptions(opts)
	var summary Summary
	for _, observation := range observations {
		if reason := rejectReason(observation, opts); reason != "" {
			summary.Drops = append(summary.Drops, Drop{DeviceID: observation.DeviceID, VehicleID: observation.VehicleID, Reason: reason})
			continue
		}
		summary.Events = append(summary.Events, Event{
			AgencyID:    observation.AgencyID,
			DeviceID:    observation.DeviceID,
			VehicleID:   observation.VehicleID,
			ObservedAt:  observation.ObservedAt.UTC().Format(time.RFC3339),
			Latitude:    observation.Latitude,
			Longitude:   observation.Longitude,
			Source:      opts.Source,
			DryRun:      true,
			NetworkSend: false,
		})
	}
	return summary
}

func normalizeOptions(opts Options) Options {
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	if opts.Source == "" {
		opts.Source = "synthetic-telemetry-connector"
	}
	if opts.MaxFutureSkew == 0 {
		opts.MaxFutureSkew = 30 * time.Second
	}
	if opts.MaxAge == 0 {
		opts.MaxAge = 2 * time.Minute
	}
	if opts.MinQuality == 0 {
		opts.MinQuality = 0.5
	}
	return opts
}

func rejectReason(observation Observation, opts Options) string {
	switch {
	case observation.AgencyID == "" || observation.DeviceID == "" || observation.VehicleID == "":
		return ReasonMissingIdentity
	case observation.ObservedAt.IsZero():
		return ReasonInvalidTimestamp
	case observation.ObservedAt.After(opts.Now.Add(opts.MaxFutureSkew)):
		return ReasonFutureTimestamp
	case opts.Now.Sub(observation.ObservedAt) > opts.MaxAge:
		return ReasonStaleObservation
	case observation.Quality < opts.MinQuality:
		return ReasonLowQuality
	case observation.Latitude < -90 || observation.Latitude > 90 || observation.Longitude < -180 || observation.Longitude > 180:
		return ReasonInvalidCoordinate
	default:
		return ""
	}
}
