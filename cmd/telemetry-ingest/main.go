package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	appdb "open-transit-rt/internal/db"
	"open-transit-rt/internal/devices"
	"open-transit-rt/internal/server"
	"open-transit-rt/internal/telemetry"
)

const maxTelemetryPayloadBytes = 1 << 20

type pinger interface {
	Ping(ctx context.Context) error
}

type ingestResponse struct {
	Accepted     bool                   `json:"accepted"`
	IngestStatus telemetry.IngestStatus `json:"ingest_status"`
	AgencyID     string                 `json:"agency_id"`
	VehicleID    string                 `json:"vehicle_id"`
	ObservedAt   time.Time              `json:"observed_at"`
	ReceivedAt   time.Time              `json:"received_at"`
}

type adminAuth interface {
	Require(...auth.Role) func(http.Handler) http.Handler
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := appdb.Connect(ctx, appdb.LoadConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := telemetry.NewPostgresRepository(pool)
	adminAuth, err := auth.MiddlewareFromEnv(pool)
	if err != nil {
		log.Fatal(err)
	}
	deviceConfig, err := devices.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run("telemetry-ingest", newHandlerWithSecurity(repo, devices.NewPostgresStore(pool, deviceConfig), pool, adminAuth, true)); err != nil {
		log.Fatal(err)
	}
}

func newHandler(repo telemetry.Repository, ready pinger) http.Handler {
	return newHandlerWithSecurity(repo, acceptingDeviceStore{}, ready, auth.TestAuthenticator{Principal: auth.Principal{
		Subject:  "test-admin",
		AgencyID: "demo-agency",
		Roles:    []auth.Role{auth.RoleAdmin, auth.RoleEditor, auth.RoleOperator, auth.RoleReadOnly},
		Method:   auth.MethodBearer,
	}}, false)
}

func newHandlerWithSecurity(repo telemetry.Repository, deviceStore devices.Store, ready pinger, admin adminAuth, requireDeviceToken bool) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"service": "telemetry-ingest",
			"status":  "ok",
		})
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := ready.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{
				"service": "telemetry-ingest",
				"status":  "unavailable",
				"error":   "database unavailable",
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"service": "telemetry-ingest",
			"status":  "ready",
		})
	})

	mux.HandleFunc("/v1/telemetry", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		evt, payload, err := decodeTelemetryPayload(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if requireDeviceToken {
			token := bearerToken(r)
			if token == "" {
				http.Error(w, "missing device token", http.StatusUnauthorized)
				return
			}
			if _, err := deviceStore.Verify(r.Context(), devices.VerifyInput{
				Token:     token,
				AgencyID:  evt.AgencyID,
				DeviceID:  evt.DeviceID,
				VehicleID: evt.VehicleID,
			}); err != nil {
				http.Error(w, "invalid device token or binding", http.StatusUnauthorized)
				return
			}
		}

		result, err := repo.Store(r.Context(), evt, payload)
		if errors.Is(err, telemetry.ErrUnknownAgency) {
			http.Error(w, "unknown agency", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "store telemetry event", http.StatusInternalServerError)
			return
		}

		statusCode := http.StatusAccepted
		if result.IngestStatus == telemetry.IngestStatusAccepted {
			statusCode = http.StatusCreated
		}
		writeJSON(w, statusCode, ingestResponse{
			Accepted:     result.IngestStatus == telemetry.IngestStatusAccepted,
			IngestStatus: result.IngestStatus,
			AgencyID:     result.AgencyID,
			VehicleID:    result.VehicleID,
			ObservedAt:   result.Timestamp,
			ReceivedAt:   result.ReceivedAt,
		})
	})

	eventsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
		if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
			return
		}
		limit, err := parseRequiredLimit(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		events, err := repo.ListEvents(r.Context(), principal.AgencyID, limit)
		if err != nil {
			http.Error(w, "list telemetry events", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, events)
	})
	adminRead := admin.Require(auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	mux.Handle("/v1/events", adminRead(eventsHandler))
	mux.Handle("/admin/debug/telemetry/events", adminRead(eventsHandler))

	return mux
}

type acceptingDeviceStore struct{}

func (acceptingDeviceStore) Verify(_ context.Context, input devices.VerifyInput) (devices.Credential, error) {
	return devices.Credential{AgencyID: input.AgencyID, DeviceID: input.DeviceID, VehicleID: input.VehicleID, Status: "active"}, nil
}

func (acceptingDeviceStore) Rebind(context.Context, devices.RebindInput) (devices.RebindResult, error) {
	return devices.RebindResult{}, nil
}

func (acceptingDeviceStore) ListBindings(context.Context, string) ([]devices.Binding, error) {
	return nil, nil
}

func decodeTelemetryPayload(w http.ResponseWriter, r *http.Request) (telemetry.Event, json.RawMessage, error) {
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxTelemetryPayloadBytes))
	if err != nil {
		return telemetry.Event{}, nil, fmt.Errorf("invalid json")
	}

	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return telemetry.Event{}, nil, fmt.Errorf("invalid json")
	}

	var evt telemetry.Event
	if err := json.Unmarshal(raw, &evt); err != nil {
		return telemetry.Event{}, nil, fmt.Errorf("invalid telemetry payload")
	}
	if !evt.Valid() {
		return telemetry.Event{}, nil, fmt.Errorf("invalid telemetry payload")
	}
	return evt, append(json.RawMessage(nil), raw...), nil
}

func bearerToken(r *http.Request) string {
	fields := strings.Fields(strings.TrimSpace(r.Header.Get("Authorization")))
	if len(fields) == 2 && strings.EqualFold(fields[0], "Bearer") {
		return fields[1]
	}
	return ""
}

func parseRequiredLimit(raw string) (int, error) {
	if raw == "" {
		return 0, fmt.Errorf("limit is required")
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 1 || limit > 500 {
		return 0, fmt.Errorf("limit must be between 1 and 500")
	}
	return limit, nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
