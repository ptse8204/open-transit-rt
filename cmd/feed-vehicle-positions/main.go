package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"open-transit-rt/internal/auth"
	appdb "open-transit-rt/internal/db"
	"open-transit-rt/internal/feed"
	"open-transit-rt/internal/server"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
	"open-transit-rt/internal/tenant"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type snapshotBuilder interface {
	Snapshot(ctx context.Context, generatedAt time.Time) (feed.VehiclePositionsSnapshot, error)
}

type agencySnapshotBuilder interface {
	SnapshotForAgency(ctx context.Context, agencyID string, generatedAt time.Time) (feed.VehiclePositionsSnapshot, error)
}

type vehiclePositionsHealthRecorder interface {
	SaveVehiclePositionsHealth(ctx context.Context, record feed.VehiclePositionsHealthRecord) error
}

type adminAuth interface {
	Require(...auth.Role) func(http.Handler) http.Handler
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

	adminAuth, err := auth.MiddlewareFromEnv(pool)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run("feed-vehicle-positions", newHandlerWithAuthAndHealth(builder, pool, adminAuth, feed.NewVehiclePositionsHealthRepository(pool))); err != nil {
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
	return newHandlerWithAuth(builder, ready, auth.TestAuthenticator{Principal: auth.Principal{
		Subject:  "test-admin",
		AgencyID: "demo-agency",
		Roles:    []auth.Role{auth.RoleAdmin, auth.RoleEditor, auth.RoleOperator, auth.RoleReadOnly},
		Method:   auth.MethodBearer,
	}})
}

func newHandlerWithAuth(builder snapshotBuilder, ready pinger, admin adminAuth) http.Handler {
	return newHandlerWithAuthAndHealth(builder, ready, admin, nil)
}

func newHandlerWithAuthAndHealth(builder snapshotBuilder, ready pinger, admin adminAuth, health vehiclePositionsHealthRecorder) http.Handler {
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
		writePublicProto(w, r, builder, health, "")
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
		if suffix != "/gtfsrt/vehicle_positions.pb" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writePublicProto(w, r, builder, health, agencyID)
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
			http.Error(w, "build vehicle positions snapshot", http.StatusInternalServerError)
			return
		}
		if snapshot.AgencyID != "" && snapshot.AgencyID != principal.AgencyID {
			http.Error(w, "feed debug belongs to another agency", http.StatusForbidden)
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
	adminRead := admin.Require(auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	mux.Handle("/public/gtfsrt/vehicle_positions.json", adminRead(debugHandler))
	mux.Handle("/admin/debug/gtfsrt/vehicle_positions.json", adminRead(debugHandler))

	return mux
}

func writePublicProto(w http.ResponseWriter, r *http.Request, builder snapshotBuilder, health vehiclePositionsHealthRecorder, agencyID string) {
	started := time.Now().UTC()
	snapshot, err := vehiclePositionsSnapshot(r.Context(), builder, agencyID, started)
	if err != nil {
		http.Error(w, "build vehicle positions snapshot", http.StatusInternalServerError)
		return
	}
	payload, err := snapshot.MarshalProto()
	if err != nil {
		http.Error(w, "marshal vehicle positions protobuf", http.StatusInternalServerError)
		return
	}
	persistVehiclePositionsHealth(health, snapshot, time.Since(started))
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("Last-Modified", snapshot.GeneratedAt.Format(http.TimeFormat))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

func vehiclePositionsSnapshot(ctx context.Context, builder snapshotBuilder, agencyID string, generatedAt time.Time) (feed.VehiclePositionsSnapshot, error) {
	if agencyID == "" {
		return builder.Snapshot(ctx, generatedAt)
	}
	if agencyBuilder, ok := builder.(agencySnapshotBuilder); ok {
		return agencyBuilder.SnapshotForAgency(ctx, agencyID, generatedAt)
	}
	snapshot, err := builder.Snapshot(ctx, generatedAt)
	if err != nil {
		return feed.VehiclePositionsSnapshot{}, err
	}
	if snapshot.AgencyID != agencyID {
		return feed.VehiclePositionsSnapshot{}, errors.New("vehicle positions builder cannot build requested agency")
	}
	return snapshot, nil
}

func persistVehiclePositionsHealth(health vehiclePositionsHealthRecorder, snapshot feed.VehiclePositionsSnapshot, generationLatency time.Duration) {
	if health == nil {
		return
	}
	record := feed.HealthRecordFromVehiclePositionsSnapshot(snapshot, generationLatency)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	go func() {
		defer cancel()
		_ = health.SaveVehiclePositionsHealth(ctx, record)
	}()
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
