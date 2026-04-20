package main

import (
	"encoding/json"
	"log"
	"net/http"

	"open-transit-rt/internal/server"
)

type response struct {
	Service string   `json:"service"`
	Status  string   `json:"status"`
	Modes   []string `json:"modes"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{
			Service: "agency-config",
			Status:  "ok",
			Modes:   []string{"gtfs-import", "gtfs-studio", "realtime"},
		})
	})

	if err := server.Run("agency-config", mux); err != nil {
		log.Fatal(err)
	}
}
