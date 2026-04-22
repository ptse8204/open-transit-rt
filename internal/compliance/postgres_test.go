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

func TestValidationScoreUsesPassedWarningAndFailedStatuses(t *testing.T) {
	if got := validationScore(EnvironmentProduction, []FeedMetadata{{FeedType: "schedule", LastValidationStatus: "passed"}}); got != StatusGreen {
		t.Fatalf("passed validator score = %s, want green", got)
	}
	if got := validationScore(EnvironmentProduction, []FeedMetadata{{FeedType: "schedule", LastValidationStatus: "warning"}}); got != StatusYellow {
		t.Fatalf("warning validator score = %s, want yellow", got)
	}
	if got := validationScore(EnvironmentProduction, []FeedMetadata{{FeedType: "schedule", LastValidationStatus: "failed"}}); got != StatusRed {
		t.Fatalf("failed validator score = %s, want red", got)
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

func TestRunValidationNormalizesPassedJSONReport(t *testing.T) {
	store := &fakeValidationStore{}
	result, err := RunValidation(context.Background(), store, ValidationRunInput{
		AgencyID: "demo-agency",
		FeedType: "schedule",
		Command:  `printf '%s' '{"status":"passed","error_count":0,"warning_count":0,"info_count":3}'`,
	})
	if err != nil {
		t.Fatalf("run validation: %v", err)
	}
	if result.Status != "passed" || result.ErrorCount != 0 || result.WarningCount != 0 || result.InfoCount != 3 {
		t.Fatalf("result = %+v, want passed with info count", result)
	}
	if store.result.Report["raw_report"] == nil {
		t.Fatalf("report = %+v, want raw_report", store.result.Report)
	}
}

func TestRunValidationNormalizesWarningJSONReport(t *testing.T) {
	store := &fakeValidationStore{}
	result, err := RunValidation(context.Background(), store, ValidationRunInput{
		AgencyID: "demo-agency",
		FeedType: "alerts",
		Command:  `printf '%s' '{"notices":[{"severity":"WARNING"},{"severity":"INFO"}]}'`,
	})
	if err != nil {
		t.Fatalf("run validation: %v", err)
	}
	if result.Status != "warning" || result.ErrorCount != 0 || result.WarningCount != 1 || result.InfoCount != 1 {
		t.Fatalf("result = %+v, want warning with notice counts", result)
	}
}

func TestRunValidationNormalizesFailedJSONReport(t *testing.T) {
	store := &fakeValidationStore{}
	result, err := RunValidation(context.Background(), store, ValidationRunInput{
		AgencyID: "demo-agency",
		FeedType: "trip_updates",
		Command:  `printf '%s' '{"summary":{"errors":2,"warnings":1,"infos":4}}'; exit 1`,
	})
	if err != nil {
		t.Fatalf("run validation: %v", err)
	}
	if result.Status != "failed" || result.ErrorCount != 2 || result.WarningCount != 1 || result.InfoCount != 4 {
		t.Fatalf("result = %+v, want failed with parsed counts", result)
	}
	if store.result.Report["error"] == "" {
		t.Fatalf("report = %+v, want command error retained", store.result.Report)
	}
}

type fakeValidationStore struct {
	result ValidationResult
}

func (f *fakeValidationStore) StoreValidationResult(_ context.Context, result ValidationResult) error {
	f.result = result
	return nil
}
