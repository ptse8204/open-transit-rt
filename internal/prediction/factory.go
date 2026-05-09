package prediction

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/gtfs"
	"open-transit-rt/internal/state"
)

const (
	AdapterNameDeterministic       = "deterministic"
	AdapterNameNoop                = "noop"
	AdapterNameExternalHTTP        = "external_http"
	AdapterNameExternalHTTPShadow  = "external_http_shadow"
	defaultExternalHTTPTimeout     = 2 * time.Second
	defaultExternalHTTPRequestCap  = 512 * 1024
	defaultExternalHTTPResponseCap = 1024 * 1024
)

type EnvLookup func(string) (string, bool)

type AdapterFactoryConfig struct {
	AdapterName         string
	DeterministicConfig DeterministicConfig
	ExternalHTTP        ExternalHTTPConfig
}

func AdapterFactoryConfigFromEnv(lookup EnvLookup) (AdapterFactoryConfig, error) {
	if lookup == nil {
		return AdapterFactoryConfig{}, fmt.Errorf("environment lookup is required")
	}
	name := strings.ToLower(strings.TrimSpace(envString(lookup, "TRIP_UPDATES_ADAPTER", AdapterNameDeterministic)))
	config := AdapterFactoryConfig{
		AdapterName: name,
		DeterministicConfig: DeterministicConfig{
			StaleTelemetryTTL:       time.Duration(envInt(lookup, "TRIP_UPDATES_STALE_TELEMETRY_TTL_SECONDS", 90)) * time.Second,
			AssignmentConfidenceMin: envFloat(lookup, "TRIP_UPDATES_ASSIGNMENT_CONFIDENCE_THRESHOLD", state.DefaultConfig().MinConfidence),
			MaxScheduleDeviation:    time.Duration(envInt(lookup, "TRIP_UPDATES_MAX_SCHEDULE_DEVIATION_SECONDS", 2700)) * time.Second,
			DuplicateConfidenceGap:  envFloat(lookup, "TRIP_UPDATES_DUPLICATE_CONFIDENCE_GAP", 0.05),
		},
	}
	switch name {
	case AdapterNameDeterministic, AdapterNameNoop:
		return config, nil
	case AdapterNameExternalHTTP, AdapterNameExternalHTTPShadow:
		external, err := externalHTTPConfigFromEnv(lookup)
		if err != nil {
			return AdapterFactoryConfig{}, err
		}
		config.ExternalHTTP = external
		return config, nil
	default:
		return AdapterFactoryConfig{}, fmt.Errorf("TRIP_UPDATES_ADAPTER must be deterministic, noop, external_http, or external_http_shadow")
	}
}

func NewAdapterFromEnv(lookup EnvLookup, scheduleRepo gtfs.Repository, operationsRepo OperationsRepository) (Adapter, error) {
	config, err := AdapterFactoryConfigFromEnv(lookup)
	if err != nil {
		return nil, err
	}
	return NewAdapter(config, scheduleRepo, operationsRepo)
}

func NewAdapter(config AdapterFactoryConfig, scheduleRepo gtfs.Repository, operationsRepo OperationsRepository) (Adapter, error) {
	switch strings.ToLower(strings.TrimSpace(config.AdapterName)) {
	case "", AdapterNameDeterministic:
		return NewDeterministicAdapter(scheduleRepo, operationsRepo, config.DeterministicConfig)
	case AdapterNameNoop:
		return NewNoopAdapter(), nil
	case AdapterNameExternalHTTP:
		return NewExternalHTTPAdapter(config.ExternalHTTP)
	case AdapterNameExternalHTTPShadow:
		deterministic, err := NewDeterministicAdapter(scheduleRepo, operationsRepo, config.DeterministicConfig)
		if err != nil {
			return nil, err
		}
		external, err := NewExternalHTTPAdapter(config.ExternalHTTP)
		if err != nil {
			return nil, err
		}
		return NewExternalHTTPShadowAdapter(deterministic, external)
	default:
		return nil, fmt.Errorf("TRIP_UPDATES_ADAPTER must be deterministic, noop, external_http, or external_http_shadow")
	}
}

func externalHTTPConfigFromEnv(lookup EnvLookup) (ExternalHTTPConfig, error) {
	config := ExternalHTTPConfig{
		URL:              envString(lookup, "TRIP_UPDATES_EXTERNAL_HTTP_URL", ""),
		AllowedHosts:     splitCSV(envString(lookup, "TRIP_UPDATES_EXTERNAL_HTTP_ALLOWED_HOSTS", "")),
		Timeout:          time.Duration(envInt(lookup, "TRIP_UPDATES_EXTERNAL_HTTP_TIMEOUT_SECONDS", int(defaultExternalHTTPTimeout/time.Second))) * time.Second,
		MaxRequestBytes:  int64(envInt(lookup, "TRIP_UPDATES_EXTERNAL_HTTP_MAX_REQUEST_BYTES", defaultExternalHTTPRequestCap)),
		MaxResponseBytes: int64(envInt(lookup, "TRIP_UPDATES_EXTERNAL_HTTP_MAX_RESPONSE_BYTES", defaultExternalHTTPResponseCap)),
		TokenEnv:         strings.TrimSpace(envString(lookup, "TRIP_UPDATES_EXTERNAL_HTTP_TOKEN_ENV", "")),
	}
	if config.TokenEnv != "" {
		if err := validateExternalHTTPTokenEnvName(config.TokenEnv); err != nil {
			return ExternalHTTPConfig{}, err
		}
		token, ok := lookup(config.TokenEnv)
		if !ok || strings.TrimSpace(token) == "" {
			return ExternalHTTPConfig{}, fmt.Errorf("referenced TRIP_UPDATES_EXTERNAL_HTTP_TOKEN_ENV is not set")
		}
		config.TokenValue = token
	}
	return config.validated()
}

func envString(lookup EnvLookup, key string, fallback string) string {
	if value, ok := lookup(key); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func envInt(lookup EnvLookup, key string, fallback int) int {
	raw, ok := lookup(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback
	}
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0
	}
	return value
}

func envFloat(lookup EnvLookup, key string, fallback float64) float64 {
	raw, ok := lookup(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return -1
	}
	return value
}

func splitCSV(raw string) []string {
	var result []string
	for _, part := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func externalHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
