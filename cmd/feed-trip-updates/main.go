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
	"open-transit-rt/internal/tenant"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type snapshotBuilder interface {
	Ready(ctx context.Context) error
	Snapshot(ctx context.Context, generatedAt time.Time) (tripupdates.Snapshot, error)
}

type agencySnapshotBuilder interface {
	SnapshotForAgency(ctx context.Context, agencyID string, generatedAt time.Time) (tripupdates.Snapshot, error)
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
		mustPredictionAdapter(prediction.NewAdapterFromEnv(os.LookupEnv, gtfs.NewPostgresRepository(pool), prediction.NewPostgresOperationsRepository(pool))),
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
	return prediction.NewAdapterFromEnv(os.LookupEnv, scheduleRepo, operationsRepo)
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
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("vehicle positions URL must not include query or fragment")
	}
	if parsed.Path == "/public/gtfsrt/vehicle_positions.pb" {
		return parsed.String(), nil
	}
	if _, matched, err := tenant.PublicAgencyPath(parsed.EscapedPath(), "/gtfsrt/vehicle_positions.pb"); matched && err == nil {
		return parsed.String(), nil
	}
	return "", fmt.Errorf("vehicle positions URL must end with /public/gtfsrt/vehicle_positions.pb or /public/agencies/{agency_id}/gtfsrt/vehicle_positions.pb")
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
		writePublicProto(w, r, builder, "")
	})
	mux.HandleFunc("/public/agencies/", func(w http.ResponseWriter, r *http.Request) {
		agencyID, suffix, matched, err := tenant.PublicAgencyRoute(r.URL.EscapedPath())
		if !matched {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "invalid agency route", http.StatusBadRequest)
			return
		}
		if suffix != "/gtfsrt/trip_updates.pb" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writePublicProto(w, r, builder, agencyID)
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

func writePublicProto(w http.ResponseWriter, r *http.Request, builder snapshotBuilder, agencyID string) {
	snapshot, err := tripUpdatesSnapshot(r.Context(), builder, agencyID, time.Now().UTC())
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
}

func tripUpdatesSnapshot(ctx context.Context, builder snapshotBuilder, agencyID string, generatedAt time.Time) (tripupdates.Snapshot, error) {
	if agencyID == "" {
		return builder.Snapshot(ctx, generatedAt)
	}
	if agencyBuilder, ok := builder.(agencySnapshotBuilder); ok {
		return agencyBuilder.SnapshotForAgency(ctx, agencyID, generatedAt)
	}
	snapshot, err := builder.Snapshot(ctx, generatedAt)
	if err != nil {
		return tripupdates.Snapshot{}, err
	}
	if snapshot.AgencyID != agencyID {
		return tripupdates.Snapshot{}, fmt.Errorf("trip updates builder cannot build requested agency")
	}
	return snapshot, nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
