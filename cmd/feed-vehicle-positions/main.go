package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	appdb "open-transit-rt/internal/db"
	"open-transit-rt/internal/feed"
	"open-transit-rt/internal/server"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type snapshotBuilder interface {
	Snapshot(ctx context.Context, generatedAt time.Time) (feed.VehiclePositionsSnapshot, error)
}

func main() {
	config, err := loadFeedConfigFromEnv()
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

	telemetryRepo := telemetry.NewPostgresRepository(pool)
	assignmentRepo := state.NewPostgresRepository(pool)
	builder, err := feed.NewVehiclePositionsBuilder(telemetryRepo, assignmentRepo, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Run("feed-vehicle-positions", newHandler(builder, pool)); err != nil {
		log.Fatal(err)
	}
}

func loadFeedConfigFromEnv() (feed.VehiclePositionsConfig, error) {
	config := feed.VehiclePositionsConfig{
		AgencyID:                  os.Getenv("AGENCY_ID"),
		MaxVehicles:               getenvInt("VEHICLE_POSITIONS_MAX_VEHICLES", 2000),
		StaleTelemetryTTL:         time.Duration(getenvInt("STALE_TELEMETRY_TTL_SECONDS", 90)) * time.Second,
		SuppressStaleVehicleAfter: time.Duration(getenvInt("SUPPRESS_STALE_VEHICLE_AFTER_SECONDS", 300)) * time.Second,
		TripConfidenceThreshold:   getenvFloat("VEHICLE_POSITIONS_TRIP_CONFIDENCE_THRESHOLD", state.DefaultConfig().MinConfidence),
	}
	return config.Validated()
}

func getenvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return value
}

func getenvFloat(key string, fallback float64) float64 {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return -1
	}
	return value
}

func newHandler(builder snapshotBuilder, ready pinger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"service": "feed-vehicle-positions",
			"status":  "ok",
		})
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := ready.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{
				"service": "feed-vehicle-positions",
				"status":  "unavailable",
				"error":   "database unavailable",
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"service": "feed-vehicle-positions",
			"status":  "ready",
		})
	})

	mux.HandleFunc("/public/gtfsrt/vehicle_positions.pb", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		snapshot, err := builder.Snapshot(r.Context(), time.Now().UTC())
		if err != nil {
			http.Error(w, "build vehicle positions snapshot", http.StatusInternalServerError)
			return
		}
		payload, err := snapshot.MarshalProto()
		if err != nil {
			http.Error(w, "marshal vehicle positions protobuf", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/x-protobuf")
		w.Header().Set("Last-Modified", snapshot.GeneratedAt.Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	})

	mux.HandleFunc("/public/gtfsrt/vehicle_positions.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		snapshot, err := builder.Snapshot(r.Context(), time.Now().UTC())
		if err != nil {
			http.Error(w, "build vehicle positions snapshot", http.StatusInternalServerError)
			return
		}
		payload, err := snapshot.MarshalDebugJSON()
		if err != nil {
			http.Error(w, "marshal vehicle positions debug json", http.StatusInternalServerError)
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
