package gtfs

import (
	"testing"
	"time"
)

func TestParseGTFSTimeAfterMidnight(t *testing.T) {
	seconds, err := ParseGTFSTime("26:10:00")
	if err != nil {
		t.Fatalf("parse gtfs time: %v", err)
	}
	if seconds != 26*3600+10*60 {
		t.Fatalf("seconds = %d, want %d", seconds, 26*3600+10*60)
	}
	if got := FormatGTFSTime(seconds); got != "26:10:00" {
		t.Fatalf("formatted = %s, want 26:10:00", got)
	}
}

func TestResolveServiceDaysIncludesPreviousLocalDate(t *testing.T) {
	observed := time.Date(2026, 4, 21, 7, 30, 0, 0, time.UTC)
	days, err := ResolveServiceDays(observed, "America/Vancouver")
	if err != nil {
		t.Fatalf("resolve service days: %v", err)
	}
	if len(days) != 2 {
		t.Fatalf("days = %d, want 2", len(days))
	}
	if days[0].Date != "20260421" || days[0].ObservedLocalSeconds != 30*60 {
		t.Fatalf("local day = %+v, want 20260421 at 00:30", days[0])
	}
	if days[1].Date != "20260420" || days[1].ObservedLocalSeconds != 24*3600+30*60 {
		t.Fatalf("previous day = %+v, want 20260420 at 24:30", days[1])
	}
}
