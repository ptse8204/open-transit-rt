package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
	"open-transit-rt/internal/devices"
	"open-transit-rt/internal/feed/schedule"
	"open-transit-rt/internal/prediction"
	"open-transit-rt/internal/state"
	"open-transit-rt/internal/telemetry"
)

func TestOperationsConsoleReferenceDeploymentProductSmoke(t *testing.T) {
	t.Setenv("PUBLICATION_ENVIRONMENT", "reference")
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	discovery := validationHealthTestDiscovery(now)
	discovery.AgencyName = "Reference Agency"
	discovery.GeneratedAt = now
	discovery.PublicBaseURL = "https://feeds.example.org"
	discovery.TechnicalContactEmail = "ops@example.org"
	discovery.PublicationEnvironment = "reference"
	discovery.License = compliance.License{Name: "CC BY 4.0", URL: "https://example.org/license"}
	discovery.Readiness = compliance.Readiness{AllRequiredFeedsListed: true, LicenseComplete: true, ContactComplete: true, HTTPSURLs: true, Discoverable: true, CanonicalValidationComplete: true}
	for i := range discovery.Feeds {
		discovery.Feeds[i].LastValidationStatus = "passed"
		discovery.Feeds[i].LastValidationAt = &now
		discovery.Feeds[i].LastHealthStatus = "ok"
		discovery.Feeds[i].LastHealthAt = &now
	}
	store := &fakePublicationStore{
		discovery: discovery,
		publicationConfig: compliance.PublicationConfig{
			AgencyID:               "demo-agency",
			PublicBaseURL:          "https://feeds.example.org",
			FeedBaseURL:            "https://feeds.example.org/public",
			TechnicalContactEmail:  "ops@example.org",
			LicenseName:            "CC BY 4.0",
			LicenseURL:             "https://example.org/license",
			PublicationEnvironment: "reference",
		},
		scorecard: compliance.Scorecard{AgencyID: "demo-agency", SnapshotAt: now, OverallStatus: compliance.StatusYellow},
		consumers: []compliance.ConsumerRecord{
			{ConsumerName: "Google Maps", Status: "prepared", UpdatedAt: now},
		},
		validationRecords: []compliance.ValidationReportRecord{
			validationHealthRecord(1, "schedule", "feed-v1", "passed", now),
			validationHealthRecord(2, "vehicle_positions", "feed-v1", "passed", now),
			validationHealthRecord(3, "trip_updates", "feed-v1", "passed", now),
			validationHealthRecord(4, "alerts", "feed-v1", "passed", now),
		},
		tripDiagnostics: compliance.TripUpdatesDiagnosticsSummary{
			Recorded:            true,
			SnapshotAt:          now,
			AdapterName:         "deterministic",
			DiagnosticsStatus:   "recorded",
			DiagnosticsReason:   "reference_smoke",
			ActiveFeedVersionID: "feed-v1",
			Metrics: prediction.Metrics{
				EligiblePredictionCandidates: 1,
				TripUpdatesEmitted:           1,
			},
		},
		reliabilityHealth: []compliance.ReliabilityFeedHealthRecord{
			{FeedType: "schedule", SnapshotAt: now},
			{FeedType: "vehicle_positions", SnapshotAt: now},
			{FeedType: "trip_updates", SnapshotAt: now},
			{FeedType: "alerts", SnapshotAt: now},
		},
		reliabilityIncidents: compliance.NormalizeReliabilityIncidentRollup(now, 0, nil, nil, nil, nil, nil, 10),
	}
	srv := newOperationsTestHandler(&handler{
		store: store,
		schedule: fakeScheduleBuilder{snapshot: schedule.Snapshot{
			AgencyID:      "demo-agency",
			FeedVersionID: "feed-v1",
			RevisionTime:  now,
			Payload:       []byte("schedule zip bytes"),
		}},
		devices: fakeDeviceStoreWithBindings{bindings: []devices.Binding{{
			AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Status: "active", ValidFrom: now, CreatedAt: now,
		}}},
		telemetry: fakeTelemetryRepository{latest: []telemetry.StoredEvent{{
			Event:      telemetry.Event{AgencyID: "demo-agency", DeviceID: "device-1", VehicleID: "bus-1", Timestamp: now.Add(-30 * time.Second), Lat: 37.1, Lon: -122.1},
			ReceivedAt: now.Add(-29 * time.Second), IngestStatus: telemetry.IngestStatusAccepted,
		}}},
		state: fakeStateRepository{assignments: map[string]state.Assignment{"bus-1": {
			VehicleID: "bus-1", State: state.StateInService, RouteID: "route-1", TripID: "trip-1", Confidence: 0.95, ActiveFrom: now.Add(-25 * time.Second),
		}}},
		realtime: &fakeRealtimeArtifacts{payloads: map[string][]byte{
			"vehicle_positions": []byte("vehicle positions pb"),
			"trip_updates":      []byte("trip updates pb"),
			"alerts":            []byte("alerts pb"),
		}},
		csrfSecret: "test-csrf",
	}, auth.TestAuthenticator{Principal: auth.Principal{
		Subject: "admin@example.com", AgencyID: "demo-agency", Roles: []auth.Role{auth.RoleAdmin}, Method: auth.MethodBearer,
	}})

	for _, route := range operationsCanonicalHTMLRoutes() {
		req := httptest.NewRequest(http.MethodGet, route.Path, nil)
		rr := httptest.NewRecorder()
		srv.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200: %s", route.Path, rr.Code, rr.Body.String())
		}
		if got := rr.Header().Get("Cache-Control"); route.NoStore && got != "no-store" {
			t.Fatalf("%s Cache-Control = %q, want no-store", route.Path, got)
		}
		body := rr.Body.String()
		for _, want := range []string{"Operations Console", "reference", "demo-agency"} {
			if !strings.Contains(body, want) {
				t.Fatalf("%s missing %q: %s", route.Path, want, body)
			}
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	body := rr.Body.String()
	for _, want := range []string{
		"Start setup",
		"Check feeds",
		"Connect vehicles",
		"Maintain system",
		"https://feeds.example.org",
		"Private Operations Console for local/self-hosted evaluation",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("dashboard missing %q: %s", want, body)
		}
	}

	unauthenticated := newOperationsTestHandler(&handler{store: store, devices: fakeDeviceStore{}}, authRejectAll{})
	req = httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr = httptest.NewRecorder()
	unauthenticated.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK {
		t.Fatalf("unauthenticated /admin/operations returned 200")
	}
	req = httptest.NewRequest(http.MethodGet, "/public/operations", nil)
	rr = httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK {
		t.Fatalf("/public/operations returned 200; private product routes must stay under /admin/operations")
	}
}
