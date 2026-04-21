package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"open-transit-rt/internal/compliance"
	appdb "open-transit-rt/internal/db"
	"open-transit-rt/internal/feed/schedule"
	"open-transit-rt/internal/server"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type scheduleBuilder interface {
	Snapshot(ctx context.Context, generatedAt time.Time) (schedule.Snapshot, error)
}

type publicationStore interface {
	BootstrapPublication(ctx context.Context, input compliance.BootstrapInput) error
	FeedDiscovery(ctx context.Context, agencyID string, generatedAt time.Time) (compliance.FeedDiscovery, error)
	UpsertConsumer(ctx context.Context, input compliance.ConsumerInput) (compliance.ConsumerRecord, error)
	ListConsumers(ctx context.Context, agencyID string) ([]compliance.ConsumerRecord, error)
	BuildAndStoreScorecard(ctx context.Context, agencyID string, at time.Time) (compliance.Scorecard, error)
	StoreValidationResult(ctx context.Context, result compliance.ValidationResult) error
}

type handler struct {
	agencyID string
	schedule scheduleBuilder
	store    publicationStore
	ready    pinger
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := appdb.Connect(ctx, appdb.LoadConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	agencyID := os.Getenv("AGENCY_ID")
	scheduleBuilder, err := schedule.NewBuilder(pool, agencyID)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run("agency-config", newHandler(agencyID, scheduleBuilder, compliance.NewPostgresRepository(pool), pool)); err != nil {
		log.Fatal(err)
	}
}

func newHandler(agencyID string, scheduleBuilder scheduleBuilder, store publicationStore, ready pinger) http.Handler {
	h := &handler{agencyID: agencyID, schedule: scheduleBuilder, store: store, ready: ready}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/readyz", h.readyz)
	mux.HandleFunc("/public/gtfs/schedule.zip", h.publicScheduleZIP)
	mux.HandleFunc("/public/feeds.json", h.publicFeedsJSON)
	mux.HandleFunc("/admin/publication/bootstrap", h.bootstrapPublication)
	mux.HandleFunc("/admin/compliance/scorecard", h.scorecard)
	mux.HandleFunc("/admin/consumer-ingestion", h.consumerIngestion)
	mux.HandleFunc("/admin/validation/run", h.runValidation)
	return mux
}

func (h *handler) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"service": "agency-config",
		"status":  "ok",
		"modes":   []string{"publication", "gtfs-schedule", "compliance"},
	})
}

func (h *handler) readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := h.ready.Ping(ctx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"service": "agency-config", "status": "unavailable", "error": "database unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"service": "agency-config", "status": "ready"})
}

func (h *handler) publicScheduleZIP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	snapshot, err := h.schedule.Snapshot(r.Context(), time.Now().UTC().Truncate(time.Second))
	if err != nil {
		http.Error(w, "build schedule zip", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Last-Modified", snapshot.RevisionTime.Format(http.TimeFormat))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(snapshot.Payload)
}

func (h *handler) publicFeedsJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	agencyID := r.URL.Query().Get("agency_id")
	if agencyID == "" {
		agencyID = h.agencyID
	}
	discovery, err := h.store.FeedDiscovery(r.Context(), agencyID, time.Now().UTC().Truncate(time.Second))
	if err != nil {
		http.Error(w, "load feed discovery metadata", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, discovery)
}

func (h *handler) bootstrapPublication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input compliance.BootstrapInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	fillBootstrapDefaults(&input)
	if err := h.store.BootstrapPublication(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"stored": true})
}

func (h *handler) scorecard(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		agencyID := r.URL.Query().Get("agency_id")
		if agencyID == "" && r.Method == http.MethodPost {
			var input struct {
				AgencyID string `json:"agency_id"`
			}
			_ = json.NewDecoder(r.Body).Decode(&input)
			agencyID = input.AgencyID
		}
		if agencyID == "" {
			agencyID = h.agencyID
		}
		scorecard, err := h.store.BuildAndStoreScorecard(r.Context(), agencyID, time.Now().UTC().Truncate(time.Second))
		if err != nil {
			http.Error(w, "build scorecard", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, scorecard)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) consumerIngestion(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		agencyID := r.URL.Query().Get("agency_id")
		if agencyID == "" {
			agencyID = h.agencyID
		}
		records, err := h.store.ListConsumers(r.Context(), agencyID)
		if err != nil {
			http.Error(w, "list consumer ingestion", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"consumers": records})
	case http.MethodPost:
		var input compliance.ConsumerInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		record, err := h.store.UpsertConsumer(r.Context(), input)
		if err != nil {
			http.Error(w, "upsert consumer ingestion", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, record)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) runValidation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input compliance.ValidationRunInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if input.FeedType == "schedule" && input.Command == "" {
		input.Command = compliance.StaticValidatorCommand(os.Getenv("GTFS_VALIDATOR_PATH"))
	}
	if input.FeedType == "schedule" && input.ScheduleZIPPath == "" {
		snapshot, err := h.schedule.Snapshot(r.Context(), time.Now().UTC().Truncate(time.Second))
		if err != nil {
			http.Error(w, "build schedule zip for validation", http.StatusInternalServerError)
			return
		}
		temp, err := os.CreateTemp("", "open-transit-rt-schedule-*.zip")
		if err != nil {
			http.Error(w, "create validation schedule temp file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(temp.Name())
		if _, err := temp.Write(snapshot.Payload); err != nil {
			_ = temp.Close()
			http.Error(w, "write validation schedule temp file", http.StatusInternalServerError)
			return
		}
		if err := temp.Close(); err != nil {
			http.Error(w, "close validation schedule temp file", http.StatusInternalServerError)
			return
		}
		input.ScheduleZIPPath = temp.Name()
		if input.FeedVersionID == "" {
			input.FeedVersionID = snapshot.FeedVersionID
		}
	}
	if input.FeedType != "schedule" && input.Command == "" {
		input.Command = os.Getenv("GTFS_RT_VALIDATOR_COMMAND")
	}
	if input.ValidatorName == "" {
		input.ValidatorName = "canonical-" + input.FeedType
	}
	result, err := compliance.RunValidation(r.Context(), h.store, input)
	if err != nil {
		http.Error(w, "run validation", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func fillBootstrapDefaults(input *compliance.BootstrapInput) {
	if input.AgencyID == "" {
		input.AgencyID = os.Getenv("AGENCY_ID")
	}
	if input.PublicBaseURL == "" {
		input.PublicBaseURL = os.Getenv("PUBLIC_BASE_URL")
	}
	if input.FeedBaseURL == "" {
		input.FeedBaseURL = os.Getenv("FEED_BASE_URL")
	}
	if input.TechnicalContactEmail == "" {
		input.TechnicalContactEmail = os.Getenv("TECHNICAL_CONTACT_EMAIL")
	}
	if input.LicenseName == "" {
		input.LicenseName = os.Getenv("FEED_LICENSE_NAME")
	}
	if input.LicenseURL == "" {
		input.LicenseURL = os.Getenv("FEED_LICENSE_URL")
	}
	if input.PublicationEnvironment == "" {
		input.PublicationEnvironment = os.Getenv("PUBLICATION_ENVIRONMENT")
	}
	if input.PublicationEnvironment == "" {
		input.PublicationEnvironment = compliance.EnvironmentDev
	}
	if input.ActorID == "" {
		input.ActorID = "system"
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
