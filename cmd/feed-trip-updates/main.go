package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	appdb "open-transit-rt/internal/db"
	"open-transit-rt/internal/feed/tripupdates"
	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/prediction"
	"open-transit-rt/internal/server"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type snapshotBuilder interface {
	Ready(ctx context.Context) error
	Snapshot(ctx context.Context, generatedAt time.Time) (tripupdates.Snapshot, error)
}

type adminAuth interface {
	Require(...auth.Role) func(http.Handler) http.Handler
}

func main() {
	config, err := loadTripUpdatesConfigFromEnv()
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

	builder, err := tripupdates.NewBuilder(
		gtfs.NewPostgresRepository(pool),
		telemetry.NewPostgresRepository(pool),
		state.NewPostgresRepository(pool),
		mustPredictionAdapter(predictionAdapterFromEnv(gtfs.NewPostgresRepository(pool), prediction.NewPostgresOperationsRepository(pool))),
		prediction.NewPostgresDiagnosticsRepository(pool),
		config,
	)
	if err != nil {
		log.Fatal(err)
	}

	adminAuth, err := auth.MiddlewareFromEnv(pool)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run("feed-trip-updates", newHandlerWithAuth(builder, pool, adminAuth)); err != nil {
		log.Fatal(err)
	}
}

func loadTripUpdatesConfigFromEnv() (tripupdates.Config, error) {
	vehiclePositionsURL, err := vehiclePositionsURLFromEnv()
	if err != nil {
		return tripupdates.Config{}, err
	}
	config := tripupdates.Config{
		AgencyID:            os.Getenv("AGENCY_ID"),
		MaxVehicles:         getenvInt("TRIP_UPDATES_MAX_VEHICLES", 2000),
		VehiclePositionsURL: vehiclePositionsURL,
	}
	return config.Validated()
}

func predictionAdapterFromEnv(scheduleRepo gtfs.Repository, operationsRepo prediction.OperationsRepository) (prediction.Adapter, error) {
	switch strings.ToLower(getenvString("TRIP_UPDATES_ADAPTER", "deterministic")) {
	case "noop":
		return prediction.NewNoopAdapter(), nil
	case "deterministic":
		return prediction.NewDeterministicAdapter(scheduleRepo, operationsRepo, prediction.DeterministicConfig{
			StaleTelemetryTTL:       time.Duration(getenvInt("TRIP_UPDATES_STALE_TELEMETRY_TTL_SECONDS", 90)) * time.Second,
			AssignmentConfidenceMin: getenvFloat("TRIP_UPDATES_ASSIGNMENT_CONFIDENCE_THRESHOLD", state.DefaultConfig().MinConfidence),
			MaxScheduleDeviation:    time.Duration(getenvInt("TRIP_UPDATES_MAX_SCHEDULE_DEVIATION_SECONDS", 2700)) * time.Second,
			DuplicateConfidenceGap:  getenvFloat("TRIP_UPDATES_DUPLICATE_CONFIDENCE_GAP", 0.05),
		})
	default:
		return nil, fmt.Errorf("TRIP_UPDATES_ADAPTER must be noop or deterministic")
	}
}

func mustPredictionAdapter(adapter prediction.Adapter, err error) prediction.Adapter {
	if err != nil {
		log.Fatal(err)
	}
	return adapter
}

func vehiclePositionsURLFromEnv() (string, error) {
	if raw := os.Getenv("VEHICLE_POSITIONS_FEED_URL"); raw != "" {
		return validateVehiclePositionsURL(raw)
	}
	base := os.Getenv("FEED_BASE_URL")
	if base == "" {
		return "", fmt.Errorf("VEHICLE_POSITIONS_FEED_URL or FEED_BASE_URL is required")
	}
	if !strings.HasSuffix(strings.TrimRight(base, "/"), "/public") {
		return "", fmt.Errorf("FEED_BASE_URL must include /public and point to the public feed root")
	}
	return validateVehiclePositionsURL(strings.TrimRight(base, "/") + "/gtfsrt/vehicle_positions.pb")
}

func validateVehiclePositionsURL(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse vehicle positions URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("vehicle positions URL must be absolute")
	}
	if parsed.Path != "/public/gtfsrt/vehicle_positions.pb" {
		return "", fmt.Errorf("vehicle positions URL must end with /public/gtfsrt/vehicle_positions.pb")
	}
	return parsed.String(), nil
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

func getenvString(key string, fallback string) string {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	return raw
}

func newHandler(builder snapshotBuilder, ready pinger) http.Handler {
	return newHandlerWithAuth(builder, ready, auth.TestAuthenticator{Principal: auth.Principal{
		Subject:  "test-admin",
		AgencyID: "demo-agency",
		Roles:    []auth.Role{auth.RoleAdmin, auth.RoleEditor, auth.RoleOperator, auth.RoleReadOnly},
		Method:   auth.MethodBearer,
	}})
}

func newHandlerWithAuth(builder snapshotBuilder, ready pinger, admin adminAuth) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"service": "feed-trip-updates",
			"status":  "ok",
		})
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := ready.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{
				"service": "feed-trip-updates",
				"status":  "unavailable",
				"error":   "database unavailable",
			})
			return
		}
		if err := builder.Ready(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{
				"service": "feed-trip-updates",
				"status":  "unavailable",
				"error":   "active feed unavailable",
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"service": "feed-trip-updates",
			"status":  "ready",
		})
	})

	mux.HandleFunc("/public/gtfsrt/trip_updates.pb", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		snapshot, err := builder.Snapshot(r.Context(), time.Now().UTC())
		if err != nil {
			http.Error(w, "build trip updates snapshot", http.StatusInternalServerError)
			return
		}
		payload, err := snapshot.MarshalProto()
		if err != nil {
			http.Error(w, "marshal trip updates protobuf", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/x-protobuf")
		w.Header().Set("Last-Modified", snapshot.GeneratedAt.Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	})

	debugHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		principal, ok := auth.PrincipalFromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		snapshot, err := builder.Snapshot(r.Context(), time.Now().UTC())
		if err != nil {
			http.Error(w, "build trip updates snapshot", http.StatusInternalServerError)
			return
		}
		if snapshot.AgencyID != "" && snapshot.AgencyID != principal.AgencyID {
			http.Error(w, "feed debug belongs to another agency", http.StatusForbidden)
			return
		}
		payload, err := snapshot.MarshalDebugJSON()
		if err != nil {
			http.Error(w, "marshal trip updates debug json", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Last-Modified", snapshot.GeneratedAt.Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	})
	adminRead := admin.Require(auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	mux.Handle("/public/gtfsrt/trip_updates.json", adminRead(debugHandler))
	mux.Handle("/admin/debug/gtfsrt/trip_updates.json", adminRead(debugHandler))

	return mux
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
