package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"open-transit-rt/internal/server"
	"open-transit-rt/internal/telemetry"
)

var (
	mu     sync.Mutex
	events []telemetry.Event
)

type ingestResponse struct {
	Accepted   bool      `json:"accepted"`
	VehicleID  string    `json:"vehicle_id"`
	ReceivedAt time.Time `json:"received_at"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mu.Lock()
		count := len(events)
		mu.Unlock()
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service": "telemetry-ingest",
			"status":  "ok",
			"count":   count,
		})
	})

	mux.HandleFunc("/v1/telemetry", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var evt telemetry.Event
		if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if !evt.Valid() {
			http.Error(w, "invalid telemetry payload", http.StatusBadRequest)
			return
		}

		mu.Lock()
		events = append(events, evt)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ingestResponse{
			Accepted:   true,
			VehicleID:  evt.VehicleID,
			ReceivedAt: time.Now().UTC(),
		})
	})

	mux.HandleFunc("/v1/events", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mu.Lock()
		copyEvents := append([]telemetry.Event(nil), events...)
		mu.Unlock()
		_ = json.NewEncoder(w).Encode(copyEvents)
	})

	if err := server.Run("telemetry-ingest", mux); err != nil {
		log.Fatal(err)
	}
}
