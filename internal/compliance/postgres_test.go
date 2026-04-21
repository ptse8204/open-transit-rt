package compliance

import (
	"context"
	"testing"
)

func TestReadinessRequiresHTTPSOnlyInProductionDiscoverability(t *testing.T) {
	cfg := feedConfig{
		PublicBaseURL:          "http://localhost:8080",
		FeedBaseURL:            "http://localhost:8080/public",
		TechnicalContactEmail:  "ops@example.org",
		LicenseName:            "CC BY 4.0",
		LicenseURL:             "https://creativecommons.org/licenses/by/4.0/",
		PublicationEnvironment: EnvironmentDev,
	}
	feeds := []FeedMetadata{
		{FeedType: "schedule", CanonicalPublicURL: "http://localhost/public/gtfs/schedule.zip", LicenseName: cfg.LicenseName, LicenseURL: cfg.LicenseURL, ContactEmail: cfg.TechnicalContactEmail, LastValidationStatus: "passed"},
		{FeedType: "vehicle_positions", CanonicalPublicURL: "http://localhost/public/gtfsrt/vehicle_positions.pb", LicenseName: cfg.LicenseName, LicenseURL: cfg.LicenseURL, ContactEmail: cfg.TechnicalContactEmail, LastValidationStatus: "passed"},
		{FeedType: "trip_updates", CanonicalPublicURL: "http://localhost/public/gtfsrt/trip_updates.pb", LicenseName: cfg.LicenseName, LicenseURL: cfg.LicenseURL, ContactEmail: cfg.TechnicalContactEmail, LastValidationStatus: "passed"},
		{FeedType: "alerts", CanonicalPublicURL: "http://localhost/public/gtfsrt/alerts.pb", LicenseName: cfg.LicenseName, LicenseURL: cfg.LicenseURL, ContactEmail: cfg.TechnicalContactEmail, LastValidationStatus: "passed"},
	}
	dev := evaluateReadiness(cfg, feeds)
	if !dev.Discoverable || dev.HTTPSURLs {
		t.Fatalf("dev readiness = %+v, want discoverable with non-HTTPS flagged", dev)
	}
	cfg.PublicationEnvironment = EnvironmentProduction
	production := evaluateReadiness(cfg, feeds)
	if production.Discoverable || production.HTTPSURLs {
		t.Fatalf("production readiness = %+v, want non-HTTPS to block discoverability", production)
	}
}

func TestValidationScoreTreatsMissingValidatorsAsProductionRed(t *testing.T) {
	feeds := []FeedMetadata{{FeedType: "schedule", LastValidationStatus: "not_run"}}
	if got := validationScore(EnvironmentDev, feeds); got != StatusYellow {
		t.Fatalf("dev missing validator score = %s, want yellow", got)
	}
	if got := validationScore(EnvironmentProduction, feeds); got != StatusRed {
		t.Fatalf("production missing validator score = %s, want red", got)
	}
}

func TestRunValidationStoresNotRunWhenCommandMissing(t *testing.T) {
	store := &fakeValidationStore{}
	result, err := RunValidation(context.Background(), store, ValidationRunInput{
		AgencyID: "demo-agency",
		FeedType: "schedule",
		Command:  "",
	})
	if err != nil {
		t.Fatalf("run validation: %v", err)
	}
	if result.Status != "not_run" || store.result.Status != "not_run" {
		t.Fatalf("result = %+v stored = %+v, want not_run", result, store.result)
	}
	if store.result.Report["reason"] != "validator_command_missing" {
		t.Fatalf("report = %+v, want missing command reason", store.result.Report)
	}
}

type fakeValidationStore struct {
	result ValidationResult
}

func (f *fakeValidationStore) StoreValidationResult(_ context.Context, result ValidationResult) error {
	f.result = result
	return nil
}
