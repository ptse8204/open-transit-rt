package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	appdb "open-transit-rt/internal/db"
	"open-transit-rt/internal/feed/alerts"
	"open-transit-rt/internal/server"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type snapshotBuilder interface {
	Snapshot(generatedAt time.Time) alerts.Snapshot
}

func main() {
	builder, err := alerts.NewBuilder(os.Getenv("AGENCY_ID"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := appdb.Connect(ctx, appdb.LoadConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := server.Run("feed-alerts", newHandler(builder, pool)); err != nil {
		log.Fatal(err)
	}
}

func newHandler(builder snapshotBuilder, ready pinger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"service": "feed-alerts",
			"status":  "ok",
		})
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := ready.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{
				"service": "feed-alerts",
				"status":  "unavailable",
				"error":   "database unavailable",
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"service": "feed-alerts",
			"status":  "ready",
		})
	})

	mux.HandleFunc("/public/gtfsrt/alerts.pb", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		snapshot := builder.Snapshot(time.Now().UTC())
		payload, err := snapshot.MarshalProto()
		if err != nil {
			http.Error(w, "marshal alerts protobuf", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/x-protobuf")
		w.Header().Set("Last-Modified", snapshot.GeneratedAt.Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	})

	mux.HandleFunc("/public/gtfsrt/alerts.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		snapshot := builder.Snapshot(time.Now().UTC())
		payload, err := snapshot.MarshalDebugJSON()
		if err != nil {
			http.Error(w, "marshal alerts debug json", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Last-Modified", snapshot.GeneratedAt.Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	})

	return mux
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
