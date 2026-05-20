package compliance

import (
	"time"

	"open-transit-rt/internal/prediction"
)

const (
	EnvironmentDev        = "dev"
	EnvironmentProduction = "production"

	StatusRed    = "red"
	StatusYellow = "yellow"
	StatusGreen  = "green"

	CanonicalStaticValidatorName    = "mobilitydata-gtfs-validator"
	InternalGTFSImportValidatorName = "open-transit-rt-internal-gtfs-import"
)

var RequiredFeedTypes = []string{"schedule", "vehicle_positions", "trip_updates", "alerts"}
var DefaultConsumers = []string{"Google Maps", "Apple Maps", "Transit App", "Bing Maps", "Moovit"}
var DefaultMarketplaceGaps = []string{
	"hardware_strategy",
	"journey_planner_integrations",
	"sla_kpi_reporting",
	"onboarding_templates",
	"support_runbooks",
	"procurement_documentation",
}

type BootstrapInput struct {
	AgencyID               string `json:"agency_id"`
	PublicBaseURL          string `json:"public_base_url"`
	FeedBaseURL            string `json:"feed_base_url"`
	TechnicalContactEmail  string `json:"technical_contact_email"`
	LicenseName            string `json:"license_name"`
	LicenseURL             string `json:"license_url"`
	PublicationEnvironment string `json:"publication_environment"`
	ActorID                string `json:"actor_id"`
}

type PublicationConfig struct {
	AgencyID               string `json:"agency_id"`
	PublicBaseURL          string `json:"public_base_url"`
	FeedBaseURL            string `json:"feed_base_url"`
	TechnicalContactEmail  string `json:"technical_contact_email"`
	LicenseName            string `json:"license_name"`
	LicenseURL             string `json:"license_url"`
	PublicationEnvironment string `json:"publication_environment"`
}

type FeedDiscovery struct {
	AgencyID               string         `json:"agency_id"`
	AgencyName             string         `json:"agency_name"`
	GeneratedAt            time.Time      `json:"generated_at"`
	PublicationEnvironment string         `json:"publication_environment"`
	PublicBaseURL          string         `json:"public_base_url"`
	TechnicalContactEmail  string         `json:"technical_contact_email"`
	License                License        `json:"license"`
	Feeds                  []FeedMetadata `json:"feeds"`
	Readiness              Readiness      `json:"readiness"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type FeedMetadata struct {
	FeedType             string     `json:"feed_type"`
	CanonicalPublicURL   string     `json:"canonical_public_url"`
	ActivationStatus     string     `json:"activation_status"`
	ActiveFeedVersionID  string     `json:"active_feed_version_id"`
	RevisionTimestamp    *time.Time `json:"revision_timestamp"`
	LicenseName          string     `json:"license_name"`
	LicenseURL           string     `json:"license_url"`
	ContactEmail         string     `json:"contact_email"`
	LastValidationStatus string     `json:"last_validation_status"`
	LastValidationAt     *time.Time `json:"last_validation_at"`
	LastHealthStatus     string     `json:"last_health_status"`
	LastHealthAt         *time.Time `json:"last_health_at"`
}

type Readiness struct {
	Discoverable                     bool `json:"discoverable"`
	HTTPSURLs                        bool `json:"https_urls"`
	LicenseComplete                  bool `json:"license_complete"`
	ContactComplete                  bool `json:"contact_complete"`
	AllRequiredFeedsListed           bool `json:"all_required_feeds_listed"`
	CanonicalValidationComplete      bool `json:"canonical_validation_complete"`
	StablePublicBaseURL              bool `json:"stable_public_base_url"`
	PublicationEnvironmentConfigured bool `json:"publication_environment_configured"`
	ActiveScheduleListed             bool `json:"active_schedule_listed"`
	RealtimeFeedsListed              bool `json:"realtime_feeds_listed"`
}

type ConsumerInput struct {
	AgencyID     string         `json:"agency_id"`
	ConsumerName string         `json:"consumer_name"`
	Status       string         `json:"status"`
	Notes        string         `json:"notes"`
	Packet       map[string]any `json:"packet"`
}

type ConsumerRecord struct {
	ConsumerName string         `json:"consumer_name"`
	Status       string         `json:"status"`
	SubmittedAt  *time.Time     `json:"submitted_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Notes        string         `json:"notes"`
	Packet       map[string]any `json:"packet"`
}

type AuditLogRecord struct {
	ID               int64     `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	Action           string    `json:"action"`
	EntityType       string    `json:"entity_type"`
	EntityID         string    `json:"entity_id"`
	ActorRecorded    bool      `json:"actor_recorded"`
	ReasonRecorded   bool      `json:"reason_recorded"`
	OldValueRecorded bool      `json:"old_value_recorded"`
	NewValueRecorded bool      `json:"new_value_recorded"`
}

type Scorecard struct {
	AgencyID                string         `json:"agency_id"`
	SnapshotAt              time.Time      `json:"snapshot_at"`
	PublicationEnvironment  string         `json:"publication_environment"`
	OverallStatus           string         `json:"overall_status"`
	ScheduleStatus          string         `json:"schedule_status"`
	VehiclePositionsStatus  string         `json:"vehicle_positions_status"`
	TripUpdatesStatus       string         `json:"trip_updates_status"`
	AlertsStatus            string         `json:"alerts_status"`
	ValidationStatus        string         `json:"validation_status"`
	DiscoverabilityStatus   string         `json:"discoverability_status"`
	ConsumerIngestionStatus string         `json:"consumer_ingestion_status"`
	Details                 map[string]any `json:"details"`
}

type ValidationRunInput struct {
	AgencyID               string `json:"agency_id"`
	FeedType               string `json:"feed_type"`
	FeedVersionID          string `json:"feed_version_id"`
	ValidatorID            string `json:"validator_id"`
	ScheduleZIPPayload     []byte `json:"-"`
	RealtimePBPayload      []byte `json:"-"`
	RealtimeArtifactSource string `json:"-"`
}

type ValidationResult struct {
	AgencyID         string         `json:"agency_id"`
	FeedType         string         `json:"feed_type"`
	FeedVersionID    string         `json:"feed_version_id"`
	ValidatorName    string         `json:"validator_name"`
	ValidatorVersion string         `json:"validator_version"`
	Status           string         `json:"status"`
	ErrorCount       int            `json:"error_count"`
	WarningCount     int            `json:"warning_count"`
	InfoCount        int            `json:"info_count"`
	Report           map[string]any `json:"report"`
}

type ValidationReportRecord struct {
	ID        int64            `json:"id"`
	Result    ValidationResult `json:"result"`
	CreatedAt time.Time        `json:"created_at"`
}

type GTFSImportRecord struct {
	ID             int64      `json:"id"`
	AgencyID       string     `json:"agency_id"`
	FeedVersionID  string     `json:"feed_version_id"`
	SourceFilename string     `json:"source_filename"`
	SourceSHA256   string     `json:"source_sha256"`
	SourceByteSize int64      `json:"source_byte_size"`
	Status         string     `json:"status"`
	ErrorCount     int        `json:"error_count"`
	WarningCount   int        `json:"warning_count"`
	InfoCount      int        `json:"info_count"`
	ActorID        string     `json:"actor_id"`
	StartedAt      time.Time  `json:"started_at"`
	CompletedAt    *time.Time `json:"completed_at"`
}

type GTFSSchedulePreview struct {
	AgencyID      string                         `json:"agency_id"`
	FeedVersionID string                         `json:"feed_version_id"`
	RowLimit      int                            `json:"row_limit"`
	Counts        GTFSSchedulePreviewCounts      `json:"counts"`
	Agency        GTFSScheduleAgencyPreview      `json:"agency"`
	Routes        []GTFSScheduleRoutePreview     `json:"routes"`
	Stops         []GTFSScheduleStopPreview      `json:"stops"`
	Trips         []GTFSScheduleTripPreview      `json:"trips"`
	Calendar      []GTFSScheduleCalendarPreview  `json:"calendar"`
	Frequencies   []GTFSScheduleFrequencyPreview `json:"frequencies"`
}

type GTFSSchedulePreviewCounts struct {
	Routes        int `json:"routes"`
	Stops         int `json:"stops"`
	Trips         int `json:"trips"`
	StopTimes     int `json:"stop_times"`
	Calendar      int `json:"calendar"`
	CalendarDates int `json:"calendar_dates"`
	ShapePoints   int `json:"shape_points"`
	Frequencies   int `json:"frequencies"`
}

type GTFSScheduleAgencyPreview struct {
	AgencyID string `json:"agency_id"`
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

type GTFSScheduleRoutePreview struct {
	ID        string `json:"id"`
	ShortName string `json:"short_name"`
	LongName  string `json:"long_name"`
	RouteType string `json:"route_type"`
}

type GTFSScheduleStopPreview struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type GTFSScheduleTripPreview struct {
	ID          string `json:"id"`
	RouteID     string `json:"route_id"`
	ServiceID   string `json:"service_id"`
	BlockID     string `json:"block_id"`
	ShapeID     string `json:"shape_id"`
	DirectionID string `json:"direction_id"`
}

type GTFSScheduleCalendarPreview struct {
	ServiceID string `json:"service_id"`
	Days      string `json:"days"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type GTFSScheduleFrequencyPreview struct {
	TripID      string `json:"trip_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	HeadwaySecs int    `json:"headway_secs"`
	ExactTimes  int    `json:"exact_times"`
}

type GTFSDraftRecord struct {
	ID                         string    `json:"id"`
	AgencyID                   string    `json:"agency_id"`
	Name                       string    `json:"name"`
	Status                     string    `json:"status"`
	BaseFeedVersionID          string    `json:"base_feed_version_id"`
	LastPublishedFeedVersionID string    `json:"last_published_feed_version_id"`
	LastPublishAttemptID       int64     `json:"last_publish_attempt_id"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

type GTFSDraftPublishRecord struct {
	ID            int64      `json:"id"`
	DraftID       string     `json:"draft_id"`
	FeedVersionID string     `json:"feed_version_id"`
	Status        string     `json:"status"`
	ErrorCount    int        `json:"error_count"`
	WarningCount  int        `json:"warning_count"`
	InfoCount     int        `json:"info_count"`
	ActorID       string     `json:"actor_id"`
	StartedAt     time.Time  `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}

type FeedVersionRecord struct {
	ID               string     `json:"id"`
	AgencyID         string     `json:"agency_id"`
	SourceType       string     `json:"source_type"`
	LifecycleState   string     `json:"lifecycle_state"`
	IsActive         bool       `json:"is_active"`
	ValidationStatus string     `json:"validation_status"`
	PublishedAt      *time.Time `json:"published_at"`
	ActivatedAt      *time.Time `json:"activated_at"`
	RetiredAt        *time.Time `json:"retired_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

type TripUpdatesDiagnosticsSummary struct {
	Recorded                      bool               `json:"recorded"`
	SnapshotAt                    time.Time          `json:"snapshot_at,omitempty"`
	AdapterName                   string             `json:"adapter_name,omitempty"`
	DiagnosticsStatus             string             `json:"diagnostics_status,omitempty"`
	DiagnosticsReason             string             `json:"diagnostics_reason,omitempty"`
	ActiveFeedVersionID           string             `json:"active_feed_version_id,omitempty"`
	VehiclePositionsURL           string             `json:"vehicle_positions_url,omitempty"`
	DiagnosticsPersistenceOutcome string             `json:"diagnostics_persistence_outcome,omitempty"`
	AdapterDetails                map[string]any     `json:"adapter_details,omitempty"`
	Metrics                       prediction.Metrics `json:"prediction_metrics"`
}
