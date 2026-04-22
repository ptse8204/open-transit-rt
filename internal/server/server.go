package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const requestIDHeader = "X-Request-ID"

var requestMetrics = newMetrics()

func Run(name string, handler http.Handler) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := os.Getenv("BIND_ADDR")
	if addr == "" {
		addr = "127.0.0.1"
	}

	srv := &http.Server{
		Addr:              addr + ":" + port,
		Handler:           wrap(name, handler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("%s listening on %s", name, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%s failed: %v", name, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown %s: %w", name, err)
	}
	return nil
}

func wrap(service string, next http.Handler) http.Handler {
	metricsEnabled := envBool("METRICS_ENABLED", false)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if metricsEnabled && r.URL.Path == "/metrics" {
			requestMetrics.writePrometheus(w)
			return
		}
		requestID := cleanRequestID(r.Header.Get(requestIDHeader))
		if requestID == "" {
			requestID = newRequestID()
		}
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		rec.Header().Set(requestIDHeader, requestID)
		start := time.Now()
		next.ServeHTTP(rec, r)
		duration := time.Since(start)
		if metricsEnabled {
			requestMetrics.observe(service, r.Method, r.URL.Path, rec.status, duration)
		}
		log.Printf(
			`{"service":%q,"request_id":%q,"method":%q,"path":%q,"status":%d,"duration_ms":%d,"remote_addr":%q}`,
			service,
			requestID,
			r.Method,
			redactedPath(r),
			rec.status,
			duration.Milliseconds(),
			redactedRemoteAddr(r.RemoteAddr),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func envBool(key string, fallback bool) bool {
	raw := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if raw == "" {
		return fallback
	}
	return raw == "1" || raw == "true" || raw == "yes" || raw == "on"
}

func cleanRequestID(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) > 128 {
		return ""
	}
	for _, ch := range raw {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' || ch == '.' {
			continue
		}
		return ""
	}
	return raw
}

func newRequestID() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(raw[:])
}

func redactedPath(r *http.Request) string {
	if r == nil || r.URL == nil {
		return ""
	}
	if r.URL.RawQuery == "" {
		return r.URL.Path
	}
	query := r.URL.Query()
	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		value := "present"
		if secretLike(key) {
			value = "redacted"
		}
		parts = append(parts, key+"="+value)
	}
	return r.URL.Path + "?" + strings.Join(parts, "&")
}

func secretLike(key string) bool {
	normalized := strings.ToLower(key)
	for _, marker := range []string{"authorization", "token", "secret", "password", "csrf", "cookie", "key"} {
		if strings.Contains(normalized, marker) {
			return true
		}
	}
	return false
}

func redactedRemoteAddr(addr string) string {
	if addr == "" {
		return ""
	}
	host := addr
	if idx := strings.LastIndex(addr, ":"); idx > 0 {
		host = addr[:idx]
	}
	if host == "" {
		return "redacted"
	}
	return host
}

type metrics struct {
	mu       sync.Mutex
	requests map[string]int64
	latency  map[string]time.Duration
}

func newMetrics() *metrics {
	return &metrics{requests: map[string]int64{}, latency: map[string]time.Duration{}}
}

func (m *metrics) observe(service string, method string, path string, status int, duration time.Duration) {
	key := service + "\xff" + method + "\xff" + path + "\xff" + strconv.Itoa(status)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requests[key]++
	m.latency[key] += duration
}

func (m *metrics) writePrometheus(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	m.mu.Lock()
	defer m.mu.Unlock()
	keys := make([]string, 0, len(m.requests))
	for key := range m.requests {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	_, _ = fmt.Fprintln(w, "# TYPE open_transit_rt_http_requests_total counter")
	for _, key := range keys {
		parts := strings.Split(key, "\xff")
		_, _ = fmt.Fprintf(w, "open_transit_rt_http_requests_total{service=%q,method=%q,path=%q,status=%q} %d\n", parts[0], parts[1], parts[2], parts[3], m.requests[key])
	}
	_, _ = fmt.Fprintln(w, "# TYPE open_transit_rt_http_request_duration_seconds_sum counter")
	for _, key := range keys {
		parts := strings.Split(key, "\xff")
		_, _ = fmt.Fprintf(w, "open_transit_rt_http_request_duration_seconds_sum{service=%q,method=%q,path=%q,status=%q} %.6f\n", parts[0], parts[1], parts[2], parts[3], m.latency[key].Seconds())
	}
}
