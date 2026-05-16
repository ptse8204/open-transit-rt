package gtfsrtconformance

import (
	"strings"
	"testing"
)

func TestLoadFixtureSuiteCoversRequiredEdgeCases(t *testing.T) {
	suite, err := LoadFixtureSuite("../../testdata/gtfsrt-conformance/suite.json")
	if err != nil {
		t.Fatalf("load fixture suite: %v", err)
	}
	if len(suite.Cases) < 12 {
		t.Fatalf("case count = %d, want at least 12 edge cases", len(suite.Cases))
	}
	if suite.ClaimFlags != (ClaimFlags{}) {
		t.Fatalf("claim flags = %+v, want all false", suite.ClaimFlags)
	}
}

func TestValidateFixtureSuiteRejectsMissingRequiredCase(t *testing.T) {
	suite, err := LoadFixtureSuite("../../testdata/gtfsrt-conformance/suite.json")
	if err != nil {
		t.Fatalf("load fixture suite: %v", err)
	}
	suite.Cases = suite.Cases[1:]
	err = ValidateFixtureSuite(suite)
	if err == nil {
		t.Fatalf("ValidateFixtureSuite succeeded after removing required case")
	}
	if !strings.Contains(err.Error(), "after_midnight_trip_descriptor_identity") {
		t.Fatalf("error = %v, want missing required scenario", err)
	}
}

func TestLoadFixtureSuiteRejectsEvidencePath(t *testing.T) {
	_, err := LoadFixtureSuite("docs/evidence/captured/gtfsrt-suite.json")
	if err == nil {
		t.Fatalf("LoadFixtureSuite accepted evidence-like path")
	}
	if !strings.Contains(err.Error(), "refusing to read evidence-like fixture suite path") {
		t.Fatalf("error = %v", err)
	}
}
