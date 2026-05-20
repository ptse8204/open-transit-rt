package alerts

import (
	"fmt"
	"regexp"
)

var (
	gtfsDatePattern = regexp.MustCompile(`^\d{8}$`)
	gtfsTimePattern = regexp.MustCompile(`^\d{2,}:[0-5]\d:[0-5]\d$`)
)

func ValidateUpsertInput(input UpsertInput) error {
	if input.ActiveStart != nil && input.ActiveEnd != nil && input.ActiveEnd.Before(*input.ActiveStart) {
		return fmt.Errorf("active_end must be after active_start")
	}
	for i, entity := range input.Entities {
		if entity.TripID == "" && (entity.StartDate != "" || entity.StartTime != "") {
			return fmt.Errorf("entity %d requires trip_id when start_date or start_time is set", i+1)
		}
		if entity.StartTime != "" && entity.StartDate == "" {
			return fmt.Errorf("entity %d requires start_date when start_time is set", i+1)
		}
		if entity.StartDate != "" && !gtfsDatePattern.MatchString(entity.StartDate) {
			return fmt.Errorf("entity %d start_date must use YYYYMMDD", i+1)
		}
		if entity.StartTime != "" && !gtfsTimePattern.MatchString(entity.StartTime) {
			return fmt.Errorf("entity %d start_time must use HH:MM:SS", i+1)
		}
	}
	return nil
}
