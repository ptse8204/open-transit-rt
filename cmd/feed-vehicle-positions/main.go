package main

import (
	"log"
	"net/http"
	"time"

	"open-transit-rt/internal/feed"
	"open-transit-rt/internal/server"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

func main() {
	matcher := state.NewRuleBasedMatcher()
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/public/gtfsrt/vehicle_positions.json", func(w http.ResponseWriter, _ *http.Request) {
		sample := telemetry.Event{
			AgencyID:  "demo-agency",
			DeviceID:  "device-001",
			VehicleID: "vehicle-001",
			Timestamp: time.Now().UTC(),
			Lat:       49.2827,
			Lon:       -123.1207,
			Bearing:   90,
			SpeedMPS:  8.2,
		}
		assignment := matcher.Assign(sample)
		payload, err := feed.BuildVehiclePositionsJSON([]telemetry.Event{sample}, map[string]state.Assignment{
			sample.VehicleID: assignment,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	})

	if err := server.Run("feed-vehicle-positions", mux); err != nil {
		log.Fatal(err)
	}
}
