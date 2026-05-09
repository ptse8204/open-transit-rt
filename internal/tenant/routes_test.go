package tenant

import "testing"

func TestValidateAgencyID(t *testing.T) {
	for _, agencyID := range []string{"agency-a", "agency_A", "A123", "agency-123"} {
		if err := ValidateAgencyID(agencyID); err != nil {
			t.Fatalf("ValidateAgencyID(%q) = %v, want nil", agencyID, err)
		}
	}

	for _, agencyID := range []string{"", ".", "..", ".hidden", "agency/a", `agency\a`, "agency%2Fbad", "agency bad", "agency.bad", "agency?x", "agency#x", "../agency", "agency-" + "012345678901234567890123456789012345678901234567890123456789"} {
		if err := ValidateAgencyID(agencyID); err == nil {
			t.Fatalf("ValidateAgencyID(%q) succeeded, want error", agencyID)
		}
	}
}

func TestPublicAgencyPath(t *testing.T) {
	agencyID, matched, err := PublicAgencyPath("/public/agencies/agency-a/gtfsrt/vehicle_positions.pb", "/gtfsrt/vehicle_positions.pb")
	if err != nil || !matched || agencyID != "agency-a" {
		t.Fatalf("PublicAgencyPath valid = (%q,%v,%v), want agency-a true nil", agencyID, matched, err)
	}

	_, matched, err = PublicAgencyPath("/public/gtfsrt/vehicle_positions.pb", "/gtfsrt/vehicle_positions.pb")
	if err != nil || matched {
		t.Fatalf("PublicAgencyPath non-match = matched %v err %v, want false nil", matched, err)
	}

	for _, path := range []string{
		"/public/agencies//gtfsrt/vehicle_positions.pb",
		"/public/agencies/./gtfsrt/vehicle_positions.pb",
		"/public/agencies/../gtfsrt/vehicle_positions.pb",
		"/public/agencies/.hidden/gtfsrt/vehicle_positions.pb",
		"/public/agencies/agency%2Fbad/gtfsrt/vehicle_positions.pb",
		"/public/agencies/agency%5Cbad/gtfsrt/vehicle_positions.pb",
		"/public/agencies/agency-a/gtfsrt/vehicle_positions.json",
		"/public/agencies/agency-a/gtfsrt/vehicle_positions.pb?agency_id=agency-b",
		"/public/agencies/agency-a/gtfsrt/vehicle_positions.pb#frag",
	} {
		if _, matched, err := PublicAgencyPath(path, "/gtfsrt/vehicle_positions.pb"); !matched || err == nil {
			t.Fatalf("PublicAgencyPath(%q) = matched %v err %v, want matched error", path, matched, err)
		}
	}
}
