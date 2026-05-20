package alerts

import (
	"strings"
	"testing"
	"time"
)

func TestValidateUpsertInputRejectsUnsafeWindowsAndTripScope(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	before := now.Add(-time.Hour)
	tests := []struct {
		name  string
		input UpsertInput
		want  string
	}{
		{
			name:  "reversed active window",
			input: UpsertInput{ActiveStart: &now, ActiveEnd: &before},
			want:  "active_end must be after active_start",
		},
		{
			name: "start date without trip",
			input: UpsertInput{Entities: []InformedEntity{{
				RouteID:   "route-10",
				StartDate: "20260520",
			}}},
			want: "requires trip_id",
		},
		{
			name: "start time without date",
			input: UpsertInput{Entities: []InformedEntity{{
				TripID:    "trip-10",
				StartTime: "08:00:00",
			}}},
			want: "requires start_date",
		},
		{
			name: "bad service date",
			input: UpsertInput{Entities: []InformedEntity{{
				TripID:    "trip-10",
				StartDate: "2026-05-20",
			}}},
			want: "YYYYMMDD",
		},
		{
			name: "bad start time",
			input: UpsertInput{Entities: []InformedEntity{{
				TripID:    "trip-10",
				StartDate: "20260520",
				StartTime: "8am",
			}}},
			want: "HH:MM:SS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpsertInput(tt.input)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("ValidateUpsertInput error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestValidateUpsertInputAcceptsScopedTripEntity(t *testing.T) {
	start := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	err := ValidateUpsertInput(UpsertInput{
		ActiveStart: &start,
		ActiveEnd:   &end,
		Entities: []InformedEntity{{
			RouteID:   "route-10",
			TripID:    "trip-10",
			StartDate: "20260520",
			StartTime: "25:15:00",
		}},
	})
	if err != nil {
		t.Fatalf("ValidateUpsertInput valid scoped trip: %v", err)
	}
}
