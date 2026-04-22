package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
	appdb "open-transit-rt/internal/db"
	"open-transit-rt/internal/devices"
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
	LatestScorecard(ctx context.Context, agencyID string) (compliance.Scorecard, error)
	BuildAndStoreScorecard(ctx context.Context, agencyID string, at time.Time) (compliance.Scorecard, error)
	StoreValidationResult(ctx context.Context, result compliance.ValidationResult) error
}

type adminAuth interface {
	Require(...auth.Role) func(http.Handler) http.Handler
}

type realtimeArtifactSource interface {
	RealtimePB(ctx context.Context, feedType string) ([]byte, error)
}

type handler struct {
	agencyID string
	schedule scheduleBuilder
	store    publicationStore
	devices  devices.Store
	ready    pinger
	admin    adminAuth
	cache    *scheduleZIPCache
	realtime realtimeArtifactSource
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
	adminAuth, err := auth.MiddlewareFromEnv(pool)
	if err != nil {
		log.Fatal(err)
	}
	deviceConfig, err := devices.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run("agency-config", newHandler(agencyID, scheduleBuilder, compliance.NewPostgresRepository(pool), devices.NewPostgresStore(pool, deviceConfig), pool, adminAuth)); err != nil {
		log.Fatal(err)
	}
}

func newHandler(agencyID string, scheduleBuilder scheduleBuilder, store publicationStore, deviceStore devices.Store, ready pinger, admin adminAuth) http.Handler {
	return newHandlerWithRealtime(agencyID, scheduleBuilder, store, deviceStore, ready, admin, realtimeArtifactSourceFromEnv())
}

func newHandlerWithRealtime(agencyID string, scheduleBuilder scheduleBuilder, store publicationStore, deviceStore devices.Store, ready pinger, admin adminAuth, realtime realtimeArtifactSource) http.Handler {
	h := &handler{agencyID: agencyID, schedule: scheduleBuilder, store: store, devices: deviceStore, ready: ready, admin: admin, cache: newScheduleZIPCache(), realtime: realtime}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/readyz", h.readyz)
	mux.HandleFunc("/public/gtfs/schedule.zip", h.publicScheduleZIP)
	mux.HandleFunc("/public/feeds.json", h.publicFeedsJSON)
	adminRead := admin.Require(auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	mux.Handle("/admin/publication/bootstrap", adminRead(http.HandlerFunc(h.bootstrapPublication)))
	mux.Handle("/admin/compliance/scorecard", adminRead(http.HandlerFunc(h.scorecard)))
	mux.Handle("/admin/consumer-ingestion", adminRead(http.HandlerFunc(h.consumerIngestion)))
	mux.Handle("/admin/validation/run", adminRead(http.HandlerFunc(h.runValidation)))
	mux.Handle("/admin/devices/rebind", adminRead(http.HandlerFunc(h.rebindDevice)))
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
	if err := h.cache.store(snapshot); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	etag := scheduleETag(snapshot)
	if r.Header.Get("If-None-Match") == etag {
		w.Header().Set("ETag", etag)
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Last-Modified", snapshot.RevisionTime.Format(http.TimeFormat))
	w.Header().Set("ETag", etag)
	w.Header().Set("X-Checksum-SHA256", payloadChecksum(snapshot.Payload))
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
	principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
	if !ok {
		return
	}
	var input compliance.BootstrapInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if auth.RejectAgencyConflict(w, input.AgencyID, principal) {
		return
	}
	input.AgencyID = principal.AgencyID
	input.ActorID = principal.Subject
	fillBootstrapDefaults(&input)
	if err := h.store.BootstrapPublication(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"stored": true})
}

func (h *handler) scorecard(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
		if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
			return
		}
		scorecard, err := h.store.LatestScorecard(r.Context(), principal.AgencyID)
		if err != nil {
			http.Error(w, "load latest scorecard", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, scorecard)
	case http.MethodPost:
		principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
		if !ok {
			return
		}
		var input struct {
			AgencyID string `json:"agency_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil && err.Error() != "EOF" {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if auth.RejectAgencyConflict(w, input.AgencyID, principal) {
			return
		}
		scorecard, err := h.store.BuildAndStoreScorecard(r.Context(), principal.AgencyID, time.Now().UTC().Truncate(time.Second))
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
		principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
		if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
			return
		}
		records, err := h.store.ListConsumers(r.Context(), principal.AgencyID)
		if err != nil {
			http.Error(w, "list consumer ingestion", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"consumers": records})
	case http.MethodPost:
		principal, ok := auth.RequireRole(w, r, auth.RoleEditor, auth.RoleAdmin)
		if !ok {
			return
		}
		var input compliance.ConsumerInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if auth.RejectAgencyConflict(w, input.AgencyID, principal) {
			return
		}
		input.AgencyID = principal.AgencyID
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
	principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
	if !ok {
		return
	}
	var input compliance.ValidationRunInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if auth.RejectAgencyConflict(w, input.AgencyID, principal) {
		return
	}
	input.AgencyID = principal.AgencyID
	registry := compliance.ValidatorRegistryFromEnv()
	if spec, ok := registry[input.ValidatorID]; ok && spec.RequiresSchedule {
		snapshot, err := h.schedule.Snapshot(r.Context(), time.Now().UTC().Truncate(time.Second))
		if err != nil {
			http.Error(w, "build schedule zip for validation", http.StatusInternalServerError)
			return
		}
		input.ScheduleZIPPayload = snapshot.Payload
		if input.FeedVersionID == "" {
			input.FeedVersionID = snapshot.FeedVersionID
		}
	}
	if spec, ok := registry[input.ValidatorID]; ok && spec.RequiresRealtime {
		payload, err := h.realtime.RealtimePB(r.Context(), input.FeedType)
		if err != nil {
			http.Error(w, "load realtime protobuf for validation", http.StatusInternalServerError)
			return
		}
		input.RealtimePBPayload = payload
	}
	result, err := compliance.RunValidation(r.Context(), h.store, registry, input)
	if err != nil {
		http.Error(w, "run validation", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *handler) rebindDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
	if !ok {
		return
	}
	var input struct {
		AgencyID  string `json:"agency_id"`
		DeviceID  string `json:"device_id"`
		VehicleID string `json:"vehicle_id"`
		Reason    string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if auth.RejectAgencyConflict(w, input.AgencyID, principal) {
		return
	}
	result, err := h.devices.Rebind(r.Context(), devices.RebindInput{
		AgencyID:  principal.AgencyID,
		DeviceID:  input.DeviceID,
		VehicleID: input.VehicleID,
		ActorID:   principal.Subject,
		Reason:    input.Reason,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type httpRealtimeArtifactSource struct {
	client   *http.Client
	urls     map[string]string
	maxBytes int64
}

func realtimeArtifactSourceFromEnv() realtimeArtifactSource {
	maxBytes := int64(10 * 1024 * 1024)
	if raw := os.Getenv("REALTIME_VALIDATION_MAX_BYTES"); raw != "" {
		if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil && parsed > 0 {
			maxBytes = parsed
		}
	}
	baseURL := firstNonEmpty(os.Getenv("REALTIME_VALIDATION_BASE_URL"), os.Getenv("FEED_BASE_URL"))
	urls := map[string]string{
		"vehicle_positions": firstNonEmpty(os.Getenv("VEHICLE_POSITIONS_FEED_URL"), realtimeFeedURL(baseURL, "vehicle_positions")),
		"trip_updates":      firstNonEmpty(os.Getenv("TRIP_UPDATES_FEED_URL"), realtimeFeedURL(baseURL, "trip_updates")),
		"alerts":            firstNonEmpty(os.Getenv("ALERTS_FEED_URL"), realtimeFeedURL(baseURL, "alerts")),
	}
	return &httpRealtimeArtifactSource{
		client:   &http.Client{Timeout: 15 * time.Second},
		urls:     urls,
		maxBytes: maxBytes,
	}
}

func (s *httpRealtimeArtifactSource) RealtimePB(ctx context.Context, feedType string) ([]byte, error) {
	rawURL := strings.TrimSpace(s.urls[feedType])
	if rawURL == "" {
		return nil, fmt.Errorf("no server-owned realtime validation URL configured for %s", feedType)
	}
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid server-owned realtime validation URL for %s", feedType)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build realtime validation request: %w", err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch realtime validation artifact: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("fetch realtime validation artifact status %d", resp.StatusCode)
	}
	payload, err := io.ReadAll(io.LimitReader(resp.Body, s.maxBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read realtime validation artifact: %w", err)
	}
	if int64(len(payload)) > s.maxBytes {
		return nil, fmt.Errorf("realtime validation artifact exceeds REALTIME_VALIDATION_MAX_BYTES")
	}
	return payload, nil
}

func realtimeFeedURL(baseURL string, feedType string) string {
	if strings.TrimSpace(baseURL) == "" {
		return ""
	}
	filename := map[string]string{
		"vehicle_positions": "vehicle_positions.pb",
		"trip_updates":      "trip_updates.pb",
		"alerts":            "alerts.pb",
	}[feedType]
	if filename == "" {
		return ""
	}
	return strings.TrimRight(baseURL, "/") + "/gtfsrt/" + filename
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

type scheduleZIPCache struct {
	mu       sync.Mutex
	maxBytes int
	entries  map[string]schedule.Snapshot
}

func newScheduleZIPCache() *scheduleZIPCache {
	maxBytes := 50 * 1024 * 1024
	if raw := os.Getenv("SCHEDULE_ZIP_MAX_BYTES"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			maxBytes = parsed
		}
	}
	return &scheduleZIPCache{maxBytes: maxBytes, entries: map[string]schedule.Snapshot{}}
}

func (c *scheduleZIPCache) store(snapshot schedule.Snapshot) error {
	if len(snapshot.Payload) > c.maxBytes {
		return fmt.Errorf("schedule zip exceeds SCHEDULE_ZIP_MAX_BYTES")
	}
	key := snapshot.FeedVersionID + ":" + snapshot.RevisionTime.UTC().Format(time.RFC3339Nano)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = map[string]schedule.Snapshot{key: snapshot}
	return nil
}

func scheduleETag(snapshot schedule.Snapshot) string {
	return `"` + snapshot.FeedVersionID + "-" + payloadChecksum(snapshot.Payload) + `"`
}

func payloadChecksum(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}
