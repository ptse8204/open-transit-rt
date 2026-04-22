package compliance

import "time"

const (
	EnvironmentDev        = "dev"
	EnvironmentProduction = "production"

	StatusRed    = "red"
	StatusYellow = "yellow"
	StatusGreen  = "green"
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
	Discoverable                bool `json:"discoverable"`
	HTTPSURLs                   bool `json:"https_urls"`
	LicenseComplete             bool `json:"license_complete"`
	ContactComplete             bool `json:"contact_complete"`
	AllRequiredFeedsListed      bool `json:"all_required_feeds_listed"`
	CanonicalValidationComplete bool `json:"canonical_validation_complete"`
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
