package tenant

import (
	"fmt"
	"strings"
)

const publicAgencyPrefix = "/public/agencies/"
const maxAgencyIDLength = 64

// ValidateAgencyID accepts only a conservative path-segment-safe agency ID.
func ValidateAgencyID(agencyID string) error {
	if agencyID == "" {
		return fmt.Errorf("agency_id is required")
	}
	if len(agencyID) > maxAgencyIDLength {
		return fmt.Errorf("agency_id is too long")
	}
	if agencyID == "." || agencyID == ".." {
		return fmt.Errorf("agency_id must not be a traversal segment")
	}
	if strings.HasPrefix(agencyID, ".") {
		return fmt.Errorf("agency_id must not start with a dot")
	}
	if strings.ContainsAny(agencyID, `/\?#`) {
		return fmt.Errorf("agency_id must be a single path segment")
	}
	for _, r := range agencyID {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '-' || r == '_':
		default:
			return fmt.Errorf("agency_id contains unsafe characters")
		}
	}
	return nil
}

// PublicAgencyPath parses paths shaped like /public/agencies/{agency_id}{suffix}.
func PublicAgencyPath(path string, suffix string) (string, bool, error) {
	agencyID, gotSuffix, matched, err := PublicAgencyRoute(path)
	if !matched || err != nil {
		return agencyID, matched, err
	}
	if gotSuffix != suffix {
		return "", true, fmt.Errorf("agency route does not match %s", suffix)
	}
	return agencyID, true, nil
}

// PublicAgencyRoute parses a public agency feed route into agency ID and suffix.
func PublicAgencyRoute(path string) (string, string, bool, error) {
	if !strings.HasPrefix(path, publicAgencyPrefix) {
		return "", "", false, nil
	}
	if strings.ContainsAny(path, "?#") {
		return "", "", true, fmt.Errorf("agency route must be parsed from path only")
	}
	rest := strings.TrimPrefix(path, publicAgencyPrefix)
	agencyID, remaining, ok := strings.Cut(rest, "/")
	if !ok {
		return "", "", true, fmt.Errorf("agency route is missing a feed path")
	}
	if hasEncodedSlashOrBackslash(agencyID) {
		return "", "", true, fmt.Errorf("agency_id must not contain encoded slashes")
	}
	if err := ValidateAgencyID(agencyID); err != nil {
		return "", "", true, err
	}
	return agencyID, "/" + remaining, true, nil
}

func hasEncodedSlashOrBackslash(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "%2f") || strings.Contains(lower, "%5c")
}
