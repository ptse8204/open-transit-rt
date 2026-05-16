package gtfsrtconformance

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const FixtureSuiteSchemaVersion = "open-transit-rt.gtfsrt_fixture_suite.v1"

type FixtureSuite struct {
	SchemaVersion string        `json:"schema_version"`
	SyntheticOnly bool          `json:"synthetic_only"`
	Boundary      string        `json:"boundary"`
	ClaimFlags    ClaimFlags    `json:"claim_flags"`
	Cases         []FixtureCase `json:"cases"`
}

type FixtureCase struct {
	ID             string   `json:"id"`
	FeedType       string   `json:"feed_type"`
	Family         string   `json:"family"`
	Scenario       string   `json:"scenario"`
	ExpectedStatus string   `json:"expected_status"`
	Assertions     []string `json:"assertions"`
	References     []string `json:"references"`
	DoesNotProve   string   `json:"does_not_prove"`
}

var requiredFixtureCases = map[string][]string{
	FeedVehiclePositions: {
		"after_midnight_trip_descriptor_identity",
		"frequency_non_exact_trip_descriptor_review",
		"suppressed_vehicle_after_stale_threshold",
		"vehicle_position_without_assignment_trip_descriptor",
		"invalid_latitude_or_longitude",
	},
	FeedTripUpdates: {
		"canceled_trip_without_stop_time_updates",
		"missing_start_date_repeated_trip_ambiguity",
		"low_confidence_or_stale_prediction_withheld",
		"missing_trip_identity",
	},
	FeedAlerts: {
		"canceled_service_alert_informed_entity",
		"missing_active_period",
		"missing_header_text",
	},
}

var requiredFixtureFamilies = []string{
	"midnight_rollover",
	"frequency_service",
	"canceled_trip",
	"stale_telemetry",
	"unknown_vehicle",
	"malformed_realtime",
}

func LoadFixtureSuite(path string) (FixtureSuite, error) {
	if evidenceLikePath(path) {
		return FixtureSuite{}, fmt.Errorf("refusing to read evidence-like fixture suite path %q", path)
	}
	payload, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return FixtureSuite{}, err
	}
	var suite FixtureSuite
	if err := json.Unmarshal(payload, &suite); err != nil {
		return FixtureSuite{}, err
	}
	if err := ValidateFixtureSuite(suite); err != nil {
		return FixtureSuite{}, err
	}
	return suite, nil
}

func ValidateFixtureSuite(suite FixtureSuite) error {
	var problems []string
	if suite.SchemaVersion != FixtureSuiteSchemaVersion {
		problems = append(problems, "schema_version must be "+FixtureSuiteSchemaVersion)
	}
	if !suite.SyntheticOnly {
		problems = append(problems, "synthetic_only must be true")
	}
	if strings.TrimSpace(suite.Boundary) == "" || !strings.Contains(strings.ToLower(suite.Boundary), "does not") {
		problems = append(problems, "boundary must be present and claim-bounded")
	}
	if suite.ClaimFlags != (ClaimFlags{}) {
		problems = append(problems, "claim_flags must all be false")
	}
	seenCaseIDs := map[string]bool{}
	seenScenarios := map[string]map[string]bool{}
	seenFamilies := map[string]bool{}
	for i, fixture := range suite.Cases {
		prefix := fmt.Sprintf("cases[%d]", i)
		if strings.TrimSpace(fixture.ID) == "" {
			problems = append(problems, prefix+".id is required")
		}
		if seenCaseIDs[fixture.ID] {
			problems = append(problems, prefix+".id is duplicated")
		}
		seenCaseIDs[fixture.ID] = true
		if !supportedFeed(fixture.FeedType) {
			problems = append(problems, prefix+".feed_type is unsupported")
		}
		if strings.TrimSpace(fixture.Family) == "" {
			problems = append(problems, prefix+".family is required")
		}
		seenFamilies[fixture.Family] = true
		if strings.TrimSpace(fixture.Scenario) == "" {
			problems = append(problems, prefix+".scenario is required")
		}
		if !validExpectedStatus(fixture.ExpectedStatus) {
			problems = append(problems, prefix+".expected_status is invalid")
		}
		if len(fixture.Assertions) == 0 {
			problems = append(problems, prefix+".assertions must not be empty")
		}
		if !strings.Contains(strings.ToLower(fixture.DoesNotProve), "does not prove") {
			problems = append(problems, prefix+".does_not_prove must keep claim boundary")
		}
		for _, ref := range fixture.References {
			if evidenceLikePath(ref) {
				problems = append(problems, prefix+".references must not point at evidence paths")
			}
		}
		if !seenScenarios[fixture.FeedType][fixture.Scenario] {
			if seenScenarios[fixture.FeedType] == nil {
				seenScenarios[fixture.FeedType] = map[string]bool{}
			}
			seenScenarios[fixture.FeedType][fixture.Scenario] = true
		}
	}
	for _, family := range requiredFixtureFamilies {
		if !seenFamilies[family] {
			problems = append(problems, "missing fixture family "+family)
		}
	}
	for feedType, scenarios := range requiredFixtureCases {
		for _, scenario := range scenarios {
			if !seenScenarios[feedType][scenario] {
				problems = append(problems, fmt.Sprintf("missing required %s scenario %s", feedType, scenario))
			}
		}
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

func validExpectedStatus(status string) bool {
	switch status {
	case StatusOK, StatusNeedsReview, StatusFailed:
		return true
	default:
		return false
	}
}

func evidenceLikePath(path string) bool {
	normalized := strings.ToLower(filepath.ToSlash(filepath.Clean(path)))
	for _, marker := range []string{"docs/evidence", "/evidence/", "evidence/", "/submission", "submission/", "/proof", "proof/"} {
		if strings.Contains(normalized, marker) {
			return true
		}
	}
	return false
}
