package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	"open-transit-rt/internal/compliance"
	"open-transit-rt/internal/devices"
	"open-transit-rt/internal/state"
)

const (
	defaultTelemetryLimit = 200
	defaultStaleSeconds   = 90
)

type operationsPage struct {
	Title              string
	AgencyID           string
	GeneratedAt        time.Time
	EnvironmentLabel   string
	CSRFToken          string
	Discovery          compliance.FeedDiscovery
	DiscoveryError     string
	ActiveFeedVersion  string
	FeedsUpdatedAt     *time.Time
	TelemetryUpdatedAt *time.Time
	ScorecardUpdatedAt *time.Time
	ConsumersUpdatedAt *time.Time
	EvidenceUpdatedAt  string
	Scorecard          *compliance.Scorecard
	ScorecardError     string
	Consumers          []consumerStatusView
	ConsumerError      string
	Telemetry          []telemetryView
	TelemetryError     string
	StaleCount         int
	Devices            []devices.Binding
	DeviceError        string
	DeviceToken        string
	DeviceTokenMeta    devices.RebindResult
	Links              []evidenceLink
	Section            string
	StaleThreshold     time.Duration
}

type consumerStatusView struct {
	Name      string
	Status    string
	UpdatedAt *time.Time
	Source    string
	Notes     string
}

type telemetryView struct {
	VehicleID        string
	DeviceID         string
	ObservedAt       time.Time
	ReceivedAt       time.Time
	AgeSeconds       int64
	Stale            bool
	AssignmentState  string
	RouteID          string
	TripID           string
	Confidence       string
	ReasonCodes      []string
	DegradedState    string
	AssignmentSource string
	AssignmentAt     *time.Time
}

type evidenceLink struct {
	Label     string
	URL       string
	UpdatedAt string
}

func (h *handler) operationsRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/admin/operations" {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderOperations(w, r, "dashboard")
		return
	}
	trimmed := strings.Trim(strings.TrimPrefix(r.URL.Path, "/admin/operations/"), "/")
	switch trimmed {
	case "feeds", "telemetry", "devices", "consumers", "evidence", "setup":
		if trimmed == "devices" && r.Method == http.MethodPost {
			h.operationsDeviceRebind(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.renderOperations(w, r, trimmed)
	default:
		http.NotFound(w, r)
	}
}

func (h *handler) renderOperations(w http.ResponseWriter, r *http.Request, section string) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, section)
	renderOperationsTemplate(w, section, page)
}

func (h *handler) operationsDeviceRebind(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	if auth.RejectAgencyConflict(w, r.FormValue("agency_id"), principal) {
		return
	}
	result, err := h.devices.Rebind(r.Context(), devices.RebindInput{
		AgencyID:  principal.AgencyID,
		DeviceID:  r.FormValue("device_id"),
		VehicleID: r.FormValue("vehicle_id"),
		ActorID:   principal.Subject,
		Reason:    r.FormValue("reason"),
	})
	page := h.buildOperationsPage(r, principal, "devices")
	if err != nil {
		page.DeviceError = err.Error()
		renderOperationsTemplate(w, "devices", page)
		return
	}
	page.DeviceToken = result.Token
	result.Token = ""
	page.DeviceTokenMeta = result
	renderOperationsTemplate(w, "devices", page)
}

func (h *handler) buildOperationsPage(r *http.Request, principal auth.Principal, section string) operationsPage {
	now := time.Now().UTC().Truncate(time.Second)
	page := operationsPage{
		Title:            "Operations Console",
		AgencyID:         principal.AgencyID,
		GeneratedAt:      now,
		EnvironmentLabel: firstNonEmpty(os.Getenv("PUBLICATION_ENVIRONMENT"), "unknown"),
		CSRFToken:        csrfToken(h.csrfSecret, principal),
		Section:          section,
		StaleThreshold:   staleThreshold(),
		Links: []evidenceLink{
			{Label: "OCI hosted evidence packet", URL: "/docs/evidence/captured/oci-pilot/2026-04-24/README.md", UpdatedAt: "2026-04-24"},
			{Label: "Consumer submission tracker", URL: "/docs/evidence/consumer-submissions/README.md", UpdatedAt: "2026-04-26"},
			{Label: "Compliance evidence checklist", URL: "/docs/compliance-evidence-checklist.md", UpdatedAt: "repo docs"},
			{Label: "Small-agency pilot operations runbook", URL: "/docs/runbooks/small-agency-pilot-operations.md", UpdatedAt: "Phase 17"},
			{Label: "Evidence redaction policy", URL: "/docs/evidence/redaction-policy.md", UpdatedAt: "Phase 15"},
		},
		EvidenceUpdatedAt: "2026-04-26",
	}

	discovery, err := h.store.FeedDiscovery(r.Context(), principal.AgencyID, now)
	if err != nil {
		page.DiscoveryError = "publication metadata is not configured yet"
	} else {
		page.Discovery = discovery
		page.EnvironmentLabel = firstNonEmpty(discovery.PublicationEnvironment, page.EnvironmentLabel)
		page.ActiveFeedVersion = activeFeedVersion(discovery.Feeds)
		page.FeedsUpdatedAt = latestFeedTime(discovery)
	}

	scorecard, err := h.store.LatestScorecard(r.Context(), principal.AgencyID)
	if err != nil {
		page.ScorecardError = "no scorecard has been stored yet"
	} else {
		page.Scorecard = &scorecard
		t := scorecard.SnapshotAt.UTC()
		page.ScorecardUpdatedAt = &t
	}

	consumers, err := h.store.ListConsumers(r.Context(), principal.AgencyID)
	if err != nil {
		page.ConsumerError = "consumer status records are not available"
	} else {
		page.Consumers = consumerStatuses(consumers)
		page.ConsumersUpdatedAt = latestConsumerTime(consumers)
	}

	page.Telemetry, page.TelemetryUpdatedAt, page.StaleCount, page.TelemetryError = h.telemetryViews(r, principal.AgencyID, now)

	bindings, err := h.devices.ListBindings(r.Context(), principal.AgencyID)
	if err != nil {
		page.DeviceError = "device bindings are not available"
	} else {
		page.Devices = bindings
	}
	return page
}

func (h *handler) telemetryViews(r *http.Request, agencyID string, now time.Time) ([]telemetryView, *time.Time, int, string) {
	if h.telemetry == nil {
		return nil, nil, 0, "telemetry repository is not available in this runtime"
	}
	latest, err := h.telemetry.ListLatestByAgency(r.Context(), agencyID, defaultTelemetryLimit)
	if err != nil {
		return nil, nil, 0, "latest telemetry is not available"
	}
	vehicleIDs := make([]string, 0, len(latest))
	for _, event := range latest {
		vehicleIDs = append(vehicleIDs, event.VehicleID)
	}
	assignments := map[string]state.Assignment{}
	if h.state != nil {
		if rows, err := h.state.ListCurrentAssignments(r.Context(), agencyID, vehicleIDs); err == nil {
			assignments = rows
		}
	}
	threshold := staleThreshold()
	var newest *time.Time
	var staleCount int
	views := make([]telemetryView, 0, len(latest))
	for _, event := range latest {
		age := now.Sub(event.Timestamp)
		stale := age > threshold
		if stale {
			staleCount++
		}
		observed := event.Timestamp.UTC()
		if newest == nil || observed.After(*newest) {
			t := observed
			newest = &t
		}
		view := telemetryView{
			VehicleID:  event.VehicleID,
			DeviceID:   event.DeviceID,
			ObservedAt: observed,
			ReceivedAt: event.ReceivedAt.UTC(),
			AgeSeconds: int64(age.Seconds()),
			Stale:      stale,
		}
		if assignment, ok := assignments[event.VehicleID]; ok {
			view.AssignmentState = string(assignment.State)
			view.RouteID = assignment.RouteID
			view.TripID = assignment.TripID
			view.Confidence = fmt.Sprintf("%.2f", assignment.Confidence)
			view.ReasonCodes = append([]string(nil), assignment.ReasonCodes...)
			view.DegradedState = string(assignment.DegradedState)
			view.AssignmentSource = string(assignment.AssignmentSource)
			if !assignment.ActiveFrom.IsZero() {
				t := assignment.ActiveFrom.UTC()
				view.AssignmentAt = &t
			}
		}
		views = append(views, view)
	}
	return views, newest, staleCount, ""
}

func activeFeedVersion(feeds []compliance.FeedMetadata) string {
	for _, feed := range feeds {
		if feed.FeedType == "schedule" && feed.ActiveFeedVersionID != "" {
			return feed.ActiveFeedVersionID
		}
	}
	for _, feed := range feeds {
		if feed.ActiveFeedVersionID != "" {
			return feed.ActiveFeedVersionID
		}
	}
	return ""
}

func latestFeedTime(discovery compliance.FeedDiscovery) *time.Time {
	var latest *time.Time
	candidates := []time.Time{discovery.GeneratedAt}
	for _, feed := range discovery.Feeds {
		if feed.RevisionTimestamp != nil {
			candidates = append(candidates, *feed.RevisionTimestamp)
		}
		if feed.LastValidationAt != nil {
			candidates = append(candidates, *feed.LastValidationAt)
		}
		if feed.LastHealthAt != nil {
			candidates = append(candidates, *feed.LastHealthAt)
		}
	}
	for _, candidate := range candidates {
		t := candidate.UTC()
		if latest == nil || t.After(*latest) {
			latest = &t
		}
	}
	return latest
}

func latestConsumerTime(consumers []compliance.ConsumerRecord) *time.Time {
	var latest *time.Time
	for _, consumer := range consumers {
		t := consumer.UpdatedAt.UTC()
		if latest == nil || t.After(*latest) {
			latest = &t
		}
	}
	return latest
}

func consumerStatuses(records []compliance.ConsumerRecord) []consumerStatusView {
	byName := map[string]compliance.ConsumerRecord{}
	for _, record := range records {
		byName[record.ConsumerName] = record
	}
	names := []string{"Google Maps", "Apple Maps", "Transit App", "Bing Maps", "Moovit", "Mobility Database", "transit.land"}
	statuses := make([]consumerStatusView, 0, len(names))
	for _, name := range names {
		if record, ok := byName[name]; ok {
			updated := record.UpdatedAt.UTC()
			statuses = append(statuses, consumerStatusView{Name: name, Status: record.Status, UpdatedAt: &updated, Source: "database record", Notes: record.Notes})
			continue
		}
		statuses = append(statuses, consumerStatusView{Name: name, Status: "not recorded in DB", Source: "docs tracker"})
	}
	return statuses
}

func staleThreshold() time.Duration {
	raw := strings.TrimSpace(os.Getenv("STALE_TELEMETRY_TTL_SECONDS"))
	if raw == "" {
		return defaultStaleSeconds * time.Second
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return defaultStaleSeconds * time.Second
	}
	return time.Duration(seconds) * time.Second
}

func feedOrder(feedType string) int {
	switch feedType {
	case "schedule":
		return 0
	case "vehicle_positions":
		return 1
	case "trip_updates":
		return 2
	case "alerts":
		return 3
	default:
		return 99
	}
}

func csrfToken(secret string, principal auth.Principal) string {
	if strings.TrimSpace(secret) == "" {
		return ""
	}
	return auth.CSRFToken(secret, principal)
}

func renderOperationsTemplate(w http.ResponseWriter, name string, data operationsPage) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := operationsTemplates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var operationsTemplates = template.Must(template.New("operations").Funcs(template.FuncMap{
	"formatTime": func(t time.Time) string {
		if t.IsZero() {
			return "not available"
		}
		return t.UTC().Format(time.RFC3339)
	},
	"formatTimePtr": func(t *time.Time) string {
		if t == nil || t.IsZero() {
			return "not available"
		}
		return t.UTC().Format(time.RFC3339)
	},
	"join": strings.Join,
	"lower": func(value string) string {
		return strings.ToLower(strings.ReplaceAll(value, "_", "-"))
	},
	"feedURL": func(discovery compliance.FeedDiscovery, feedType string) string {
		for _, feed := range discovery.Feeds {
			if feed.FeedType == feedType {
				return feed.CanonicalPublicURL
			}
		}
		return ""
	},
	"sortedFeeds": func(feeds []compliance.FeedMetadata) []compliance.FeedMetadata {
		out := append([]compliance.FeedMetadata(nil), feeds...)
		sort.SliceStable(out, func(i, j int) bool { return feedOrder(out[i].FeedType) < feedOrder(out[j].FeedType) })
		return out
	},
}).Parse(`
{{define "layoutStart"}}
<!doctype html><html><head><meta charset="utf-8"><title>{{.Title}}</title>
<style>
body{font-family:system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;margin:2rem;line-height:1.4;color:#1f2933}
nav a{margin-right:1rem} table{border-collapse:collapse;width:100%;margin:1rem 0} th,td{border:1px solid #d8dee4;padding:.45rem;text-align:left;vertical-align:top}
th{background:#f6f8fa}.pill{display:inline-block;border:1px solid #c8d1dc;border-radius:3px;padding:.1rem .35rem;background:#f6f8fa}
.warning{background:#fff8c5}.ok{background:#dafbe1}.bad{background:#ffebe9}.muted{color:#59636e}.token{border:1px solid #f0c36d;background:#fff8c5;padding:1rem}
form{margin:1rem 0} label{display:block;margin:.35rem 0} input,select,textarea{min-width:22rem;max-width:100%;padding:.35rem} button{padding:.4rem .7rem}
</style></head><body>
<h1>{{.Title}}</h1>
<p>Agency: <strong>{{.AgencyID}}</strong> · environment: <span class="pill">{{.EnvironmentLabel}}</span> · generated: {{formatTime .GeneratedAt}}</p>
<nav>
<a href="/admin/operations">Dashboard</a>
<a href="/admin/operations/feeds">Feeds</a>
<a href="/admin/operations/telemetry">Telemetry</a>
<a href="/admin/operations/devices">Devices</a>
<a href="/admin/alerts/console">Alerts</a>
<a href="/admin/operations/consumers">Consumers</a>
<a href="/admin/operations/evidence">Evidence</a>
<a href="/admin/operations/setup">Setup</a>
<a href="/admin/gtfs-studio">GTFS Studio</a>
</nav>
{{end}}
{{define "layoutEnd"}}</body></html>{{end}}

{{define "dashboard"}}
{{template "layoutStart" .}}
<h2>Readiness</h2>
{{if .DiscoveryError}}<p class="warning">{{.DiscoveryError}}. Next action: bootstrap publication metadata after a feed is available.</p>{{else}}
<p>Active GTFS feed version: {{if .ActiveFeedVersion}}<strong>{{.ActiveFeedVersion}}</strong>{{else}}not available{{end}}</p>
<table><tbody>
<tr><th>Public URLs</th><td>{{if .Discovery.Readiness.AllRequiredFeedsListed}}listed{{else}}missing or incomplete{{end}}</td></tr>
<tr><th>License</th><td>{{if .Discovery.Readiness.LicenseComplete}}complete{{else}}missing{{end}}</td></tr>
<tr><th>Contact</th><td>{{if .Discovery.Readiness.ContactComplete}}complete{{else}}missing{{end}}</td></tr>
<tr><th>HTTPS URLs</th><td>{{if .Discovery.Readiness.HTTPSURLs}}yes{{else}}not all HTTPS; local/dev URLs may be HTTP{{end}}</td></tr>
<tr><th>Canonical validation</th><td>{{if .Discovery.Readiness.CanonicalValidationComplete}}current passed/warning records exist{{else}}not complete{{end}}</td></tr>
</tbody></table>{{end}}

<h2>Dashboard Sections</h2>
<table><thead><tr><th>Section</th><th>Status</th><th>Last updated</th><th>Next action</th></tr></thead><tbody>
<tr><td>Feeds / validation</td><td>{{if .DiscoveryError}}not configured{{else}}{{len .Discovery.Feeds}} feed records{{end}}</td><td>{{formatTimePtr .FeedsUpdatedAt}}</td><td><a href="/admin/operations/feeds">review feed URLs and validation</a></td></tr>
<tr><td>Telemetry freshness</td><td>{{if .TelemetryError}}{{.TelemetryError}}{{else}}{{len .Telemetry}} vehicles; {{.StaleCount}} stale{{end}}</td><td>{{formatTimePtr .TelemetryUpdatedAt}}</td><td><a href="/admin/operations/telemetry">inspect vehicle freshness</a></td></tr>
<tr><td>Scorecard</td><td>{{if .Scorecard}}{{.Scorecard.OverallStatus}}{{else}}{{.ScorecardError}}{{end}}</td><td>{{formatTimePtr .ScorecardUpdatedAt}}</td><td><a href="/admin/operations/evidence">find scorecard evidence</a></td></tr>
<tr><td>Consumer status</td><td>{{if .ConsumerError}}{{.ConsumerError}}{{else}}{{len .Consumers}} targets shown{{end}}</td><td>{{formatTimePtr .ConsumersUpdatedAt}}</td><td><a href="/admin/operations/consumers">review evidence-only statuses</a></td></tr>
<tr><td>Evidence links</td><td>repo documentation links</td><td>{{.EvidenceUpdatedAt}}</td><td><a href="/admin/operations/evidence">open evidence index</a></td></tr>
</tbody></table>

<h2>Public Feed URLs</h2>
{{if .DiscoveryError}}<p>No public feed metadata is available yet.</p>{{else}}{{template "feedTable" .}}{{end}}
<p class="muted">Validation and public fetch records are supporting evidence only. They are not consumer acceptance or CAL-ITP/Caltrans compliance by themselves.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "feeds"}}
{{template "layoutStart" .}}
<h2>Feed URLs And Validation</h2>
{{if .DiscoveryError}}<p class="warning">No feed metadata is available. Next action: publish or import a GTFS feed, then bootstrap publication metadata.</p>{{else}}
{{template "feedTable" .}}
<h3>Feed discovery document</h3>
<table><thead><tr><th>Item</th><th>URL</th><th>Validation</th><th>Last checked</th></tr></thead><tbody>
<tr><td>feeds.json</td><td>{{.Discovery.PublicBaseURL}}/public/feeds.json</td><td>not a validator result</td><td>{{formatTime .Discovery.GeneratedAt}}</td></tr>
</tbody></table>
{{end}}
<p class="muted">This view shows repo/deployment evidence only. Third-party consumer acceptance requires retained confirmation from the named consumer.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "feedTable"}}
<table><thead><tr><th>Feed</th><th>URL</th><th>Validation</th><th>Validation time</th><th>Health</th><th>Health time</th><th>Active feed version</th></tr></thead><tbody>
{{range sortedFeeds .Discovery.Feeds}}<tr><td>{{.FeedType}}</td><td>{{if .CanonicalPublicURL}}<code>{{.CanonicalPublicURL}}</code>{{else}}missing{{end}}</td><td>{{.LastValidationStatus}}</td><td>{{formatTimePtr .LastValidationAt}}</td><td>{{.LastHealthStatus}}</td><td>{{formatTimePtr .LastHealthAt}}</td><td>{{.ActiveFeedVersionID}}</td></tr>{{end}}
</tbody></table>
{{end}}

{{define "telemetry"}}
{{template "layoutStart" .}}
<h2>Telemetry Freshness</h2>
<p>Stale threshold: {{.StaleThreshold}}</p>
{{if .TelemetryError}}<p class="warning">{{.TelemetryError}}. Next action: confirm the telemetry service and database are running.</p>{{else if not .Telemetry}}<p class="warning">No telemetry has been accepted yet. Next action: create or rotate a device token, configure the device, then send a sample telemetry event.</p>{{else}}
<table><thead><tr><th>Vehicle</th><th>Device</th><th>Observed</th><th>Age seconds</th><th>Freshness</th><th>Assignment</th><th>Trip</th><th>Route</th><th>Confidence</th><th>Reasons</th><th>Assignment time</th></tr></thead><tbody>
{{range .Telemetry}}<tr><td>{{.VehicleID}}</td><td>{{.DeviceID}}</td><td>{{formatTime .ObservedAt}}</td><td>{{.AgeSeconds}}</td><td>{{if .Stale}}stale{{else}}fresh{{end}}</td><td>{{if .AssignmentState}}{{.AssignmentState}}{{else}}not available{{end}}{{if .DegradedState}} / {{.DegradedState}}{{end}}</td><td>{{.TripID}}</td><td>{{.RouteID}}</td><td>{{.Confidence}}</td><td>{{join .ReasonCodes ", "}}</td><td>{{formatTimePtr .AssignmentAt}}</td></tr>{{end}}
</tbody></table>{{end}}
<p class="muted">Safe diagnostics omit raw telemetry payloads, full score details, token fields, and private debug blobs.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "devices"}}
{{template "layoutStart" .}}
<h2>Device Credentials</h2>
<p class="warning">Device tokens are secrets. Store a one-time token immediately; it will not be shown again by this console.</p>
<p>The current supported browser flow is rotate/rebind. If a device has no credential yet, this uses the existing rebind API path; Phase 18 does not add a separate first-time creation API.</p>
{{if .DeviceToken}}<div class="token"><h3>One-time token</h3><p>Device: {{.DeviceTokenMeta.DeviceID}} · Vehicle: {{.DeviceTokenMeta.VehicleID}} · Rotated: {{.DeviceTokenMeta.RotatedAt}}</p><p><code>{{.DeviceToken}}</code></p></div>{{end}}
{{if .DeviceError}}<p class="warning">{{.DeviceError}}</p>{{end}}
<form method="post" action="/admin/operations/devices">
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
<input type="hidden" name="agency_id" value="{{.AgencyID}}">
<label>Device ID <input name="device_id" required></label>
<label>Vehicle ID <input name="vehicle_id" required></label>
<label>Reason <input name="reason" placeholder="rotation or rebind reason"></label>
<button>Rotate / rebind token</button>
</form>
{{if not .Devices}}<p class="warning">No device bindings are recorded. Next action: use the rotate/rebind form for the first device and store the returned token securely.</p>{{else}}
<table><thead><tr><th>Device</th><th>Vehicle</th><th>Status</th><th>Valid from</th><th>Last used</th><th>Rotated</th><th>Revoked</th></tr></thead><tbody>
{{range .Devices}}<tr><td>{{.DeviceID}}</td><td>{{.VehicleID}}</td><td>{{.Status}}</td><td>{{formatTime .ValidFrom}}</td><td>{{formatTimePtr .LastUsedAt}}</td><td>{{formatTimePtr .RotatedAt}}</td><td>{{formatTimePtr .RevokedAt}}</td></tr>{{end}}
</tbody></table>{{end}}
{{template "layoutEnd" .}}
{{end}}

{{define "consumers"}}
{{template "layoutStart" .}}
<h2>Consumer Submission Evidence</h2>
<p class="muted">Statuses are evidence records only. They are not consumer acceptance unless retained third-party evidence exists for the named consumer.</p>
{{if .ConsumerError}}<p class="warning">{{.ConsumerError}}. Next action: review <code>docs/evidence/consumer-submissions/README.md</code>.</p>{{else if not .Consumers}}<p class="warning">No consumer evidence records are available. Next action: review the docs tracker before making any submission claims.</p>{{else}}
<table><thead><tr><th>Target</th><th>Status</th><th>Source</th><th>Updated</th><th>Notes</th></tr></thead><tbody>
{{range .Consumers}}<tr><td>{{.Name}}</td><td>{{.Status}}</td><td>{{.Source}}</td><td>{{formatTimePtr .UpdatedAt}}</td><td>{{.Notes}}</td></tr>{{end}}
</tbody></table>{{end}}
<p>Docs tracker: <code>docs/evidence/consumer-submissions/README.md</code></p>
{{template "layoutEnd" .}}
{{end}}

{{define "evidence"}}
{{template "layoutStart" .}}
<h2>Evidence And Runbook Links</h2>
<table><thead><tr><th>Record</th><th>Location</th><th>Last updated</th></tr></thead><tbody>
{{range .Links}}<tr><td>{{.Label}}</td><td><code>{{.URL}}</code></td><td>{{.UpdatedAt}}</td></tr>{{end}}
</tbody></table>
<p class="muted">These links help operators find repo/deployment evidence. They do not assert consumer acceptance, hosted SaaS availability, agency endorsement, or universal production readiness.</p>
{{template "layoutEnd" .}}
{{end}}

{{define "setup"}}
{{template "layoutStart" .}}
<h2>Guided Setup Checklist</h2>
<table><thead><tr><th>Step</th><th>Current signal</th><th>Next action</th></tr></thead><tbody>
<tr><td>Agency metadata</td><td>{{if .DiscoveryError}}not configured{{else}}{{.Discovery.AgencyName}}{{end}}</td><td>Use publication bootstrap and GTFS Studio agency metadata.</td></tr>
<tr><td>GTFS schedule</td><td>{{if .ActiveFeedVersion}}{{.ActiveFeedVersion}}{{else}}no active feed shown{{end}}</td><td><a href="/admin/gtfs-studio">Open GTFS Studio</a> or import a GTFS ZIP with the existing import command.</td></tr>
<tr><td>Publication URLs</td><td>{{if .Discovery.Readiness.AllRequiredFeedsListed}}listed{{else}}missing{{end}}</td><td><a href="/admin/operations/feeds">Review feed URLs</a>.</td></tr>
<tr><td>Device token</td><td>{{len .Devices}} binding records</td><td><a href="/admin/operations/devices">Rotate/rebind a token</a> and store it securely.</td></tr>
<tr><td>Validation</td><td>{{if .Discovery.Readiness.CanonicalValidationComplete}}supporting records exist{{else}}not complete{{end}}</td><td>Run the existing admin validation workflow, then return to the feed view.</td></tr>
<tr><td>Consumer evidence</td><td>{{len .Consumers}} targets displayed</td><td><a href="/admin/operations/consumers">Review evidence-only statuses</a>.</td></tr>
</tbody></table>
{{template "layoutEnd" .}}
{{end}}
`))
