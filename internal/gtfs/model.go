package gtfs

import (
	"context"
	"fmt"
	"time"
)

type Agency struct {
	ID       string
	Timezone string
}

type FeedVersion struct {
	ID       string
	AgencyID string
}

type StopTime struct {
	TripID            string
	StopID            string
	ArrivalSeconds    int
	DepartureSeconds  int
	StopSequence      int
	ShapeDistTraveled float64
}

type ShapePoint struct {
	ShapeID      string
	Lat          float64
	Lon          float64
	Sequence     int
	DistTraveled float64
	HasDistance  bool
}

type Frequency struct {
	TripID       string
	StartSeconds int
	EndSeconds   int
	HeadwaySecs  int
	ExactTimes   int
	StartTime    string
	EndTime      string
}

type TripCandidate struct {
	AgencyID      string
	FeedVersionID string
	ServiceDate   string
	RouteID       string
	ServiceID     string
	TripID        string
	BlockID       string
	ShapeID       string
	DirectionID   *int
	StopTimes     []StopTime
	ShapePoints   []ShapePoint
	Frequencies   []Frequency
}

type Repository interface {
	Agency(ctx context.Context, agencyID string) (Agency, error)
	ActiveFeedVersion(ctx context.Context, agencyID string) (FeedVersion, error)
	ListTripCandidates(ctx context.Context, agencyID string, feedVersionID string, serviceDate string) ([]TripCandidate, error)
}

func ParseGTFSTime(raw string) (int, error) {
	var hour, minute, second int
	if _, err := fmt.Sscanf(raw, "%d:%d:%d", &hour, &minute, &second); err != nil {
		return 0, fmt.Errorf("parse gtfs time %q: %w", raw, err)
	}
	if hour < 0 || minute < 0 || minute > 59 || second < 0 || second > 59 {
		return 0, fmt.Errorf("invalid gtfs time %q", raw)
	}
	return hour*3600 + minute*60 + second, nil
}

func FormatGTFSTime(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}
	return fmt.Sprintf("%02d:%02d:%02d", seconds/3600, (seconds%3600)/60, seconds%60)
}

func ServiceDateString(t time.Time) string {
	return t.Format("20060102")
}

func ResolveServiceDays(observedAt time.Time, timezone string) ([]ServiceDay, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("load agency timezone %q: %w", timezone, err)
	}
	local := observedAt.In(loc)
	localMidnight := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)

	days := make([]ServiceDay, 0, 2)
	for _, midnight := range []time.Time{localMidnight, localMidnight.AddDate(0, 0, -1)} {
		seconds := int(local.Sub(midnight).Seconds())
		days = append(days, ServiceDay{
			Date:                 ServiceDateString(midnight),
			ObservedLocalSeconds: seconds,
			LocalTime:            local,
		})
	}
	return days, nil
}

type ServiceDay struct {
	Date                 string
	ObservedLocalSeconds int
	LocalTime            time.Time
}
