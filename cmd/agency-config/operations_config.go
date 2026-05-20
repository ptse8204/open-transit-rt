package main

import (
	"net/http"
	"time"

	"open-transit-rt/internal/auth"
)

type operationsConfigView struct {
	GeneratedAt   time.Time                 `json:"generated_at"`
	AgencyID      string                    `json:"agency_id"`
	ActiveSection string                    `json:"active_section"`
	Title         string                    `json:"title"`
	Boundary      string                    `json:"boundary"`
	AdminCanEdit  bool                      `json:"admin_can_edit"`
	Sections      []operationsConfigSection `json:"sections"`
	Rows          []operationsConfigRow     `json:"rows"`
	ClaimBoundary string                    `json:"claim_boundary"`
}

type operationsConfigSection struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Path       string `json:"path"`
	Status     string `json:"status"`
	Summary    string `json:"summary"`
	NextAction string `json:"next_action"`
}

type operationsConfigRow struct {
	Label      string `json:"label"`
	Value      string `json:"value"`
	Status     string `json:"status"`
	NextAction string `json:"next_action"`
}

func (h *handler) renderOperationsConfig(w http.ResponseWriter, r *http.Request, section string) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, section)
	renderOperationsTemplate(w, "config", page)
}

func buildOperationsConfig(page operationsPage) operationsConfigView {
	active := operationsConfigActiveSection(page.Section)
	view := operationsConfigView{
		GeneratedAt:   page.GeneratedAt,
		AgencyID:      page.AgencyID,
		ActiveSection: active,
		Title:         operationsConfigTitle(active),
		Boundary:      "Focused private configuration pages only. These pages separate settings from diagnostics, keep secrets out of HTML, and do not prove compliance, consumer acceptance, production readiness, SSO, hosted service availability, AVL reliability, or ETA quality.",
		AdminCanEdit:  page.IsAdmin,
		Sections:      operationsConfigSections(page),
		ClaimBoundary: "Configuration review is an internal operator workflow. Advanced state remains behind private JSON exports or collapsed review panels.",
	}
	view.Rows = operationsConfigRows(page, active)
	return view
}

func operationsConfigActiveSection(section string) string {
	switch section {
	case "config-agency":
		return "agency"
	case "config-feeds":
		return "feeds"
	case "config-auth":
		return "auth"
	case "config-deployment":
		return "deployment"
	case "config-advanced":
		return "advanced"
	default:
		return "overview"
	}
}

func operationsConfigSectionFromPath(path string) string {
	switch path {
	case "config/agency":
		return "config-agency"
	case "config/feeds":
		return "config-feeds"
	case "config/auth":
		return "config-auth"
	case "config/deployment":
		return "config-deployment"
	case "config/advanced":
		return "config-advanced"
	default:
		return "config"
	}
}

func operationsConfigTitle(active string) string {
	switch active {
	case "agency":
		return "Agency Profile"
	case "feeds":
		return "Public Feed URLs"
	case "auth":
		return "Login Settings"
	case "deployment":
		return "Deployment Settings"
	case "advanced":
		return "Advanced Settings"
	default:
		return "Focused Config"
	}
}

func operationsConfigSections(page operationsPage) []operationsConfigSection {
	return []operationsConfigSection{
		{
			ID:         "agency",
			Label:      "Agency Profile",
			Path:       "/admin/operations/config/agency",
			Status:     configStatus(firstNonEmpty(page.Discovery.AgencyName, page.AgencyID)),
			Summary:    "Agency identity and operator-facing profile details.",
			NextAction: "Review the agency identity used by feed discovery and setup wizard rows.",
		},
		{
			ID:         "feeds",
			Label:      "Public Feed URLs",
			Path:       "/admin/operations/config/feeds",
			Status:     configStatus(firstNonEmpty(page.PublicationConfig.PublicBaseURL, page.Discovery.PublicBaseURL)),
			Summary:    "Public base URL, feed base URL, license, contact, and publication environment.",
			NextAction: "Admins can update publication metadata through the existing setup bootstrap action.",
		},
		{
			ID:         "auth",
			Label:      "Login Settings",
			Path:       "/admin/operations/config/auth",
			Status:     configStatus(page.AuthStatus.PasswordLogin),
			Summary:    "Local demo login, password login, sessions, CSRF, and future SSO boundary.",
			NextAction: "Use Login & Sessions or Users & Roles for auth review and password reset links.",
		},
		{
			ID:         "deployment",
			Label:      "Deployment Settings",
			Path:       "/admin/operations/config/deployment",
			Status:     configStatus(page.EnvironmentLabel),
			Summary:    "Non-secret deployment signals needed for operation and support.",
			NextAction: "Configure deployment-owned environment values outside the browser when rows are missing.",
		},
		{
			ID:         "advanced",
			Label:      "Advanced Settings",
			Path:       "/admin/operations/config/advanced",
			Status:     "review",
			Summary:    "Collapsed links to private diagnostics and JSON exports.",
			NextAction: "Open advanced details only when triaging a specific support or safety question.",
		},
	}
}

func operationsConfigRows(page operationsPage, active string) []operationsConfigRow {
	switch active {
	case "agency":
		return []operationsConfigRow{
			configRow("Agency ID", page.AgencyID, "configured", "Agency ID comes from the authenticated principal and is read-only here."),
			configRow("Agency name", firstNonEmpty(page.Discovery.AgencyName, "missing"), configStatus(page.Discovery.AgencyName), "Set agency profile metadata through the deployment-owned publication workflow."),
			configRow("Setup progress", page.SetupWizard.Summary.Status, page.SetupWizard.Summary.Status, page.SetupWizard.Summary.NextAction),
			configRow("Role visibility", configRoleVisibilitySignal(page), "review", "Open Access & Roles when a user cannot see a needed action."),
		}
	case "feeds":
		return []operationsConfigRow{
			configRow("Public base URL", firstNonEmpty(page.PublicationConfig.PublicBaseURL, page.Discovery.PublicBaseURL, "missing"), configStatus(firstNonEmpty(page.PublicationConfig.PublicBaseURL, page.Discovery.PublicBaseURL)), "Use a stable HTTPS URL before sharing feeds."),
			configRow("Feed base URL", firstNonEmpty(page.PublicationConfig.FeedBaseURL, "missing"), configStatus(page.PublicationConfig.FeedBaseURL), "Store the deployment-owned feed base URL without embedding credentials."),
			configRow("Technical contact", firstNonEmpty(page.PublicationConfig.TechnicalContactEmail, page.Discovery.TechnicalContactEmail, "missing"), configStatus(firstNonEmpty(page.PublicationConfig.TechnicalContactEmail, page.Discovery.TechnicalContactEmail)), "Provide an agency or operator contact before discovery review."),
			configRow("License", firstNonEmpty(page.PublicationConfig.LicenseName, page.Discovery.License.Name, "missing"), configStatus(firstNonEmpty(page.PublicationConfig.LicenseName, page.Discovery.License.Name)), "Publish license metadata before sharing feed URLs."),
			configRow("Publication environment", firstNonEmpty(page.PublicationConfig.PublicationEnvironment, page.Discovery.PublicationEnvironment, page.EnvironmentLabel, "missing"), configStatus(firstNonEmpty(page.PublicationConfig.PublicationEnvironment, page.Discovery.PublicationEnvironment, page.EnvironmentLabel)), "Keep dev, staging, and production labels explicit."),
		}
	case "auth":
		return []operationsConfigRow{
			configRow("Current auth mode", page.AuthStatus.ActiveAuthMode, "review", "Review the current browser or bearer session."),
			configRow("Password login", page.AuthStatus.PasswordLogin, configStatus(page.AuthStatus.PasswordLogin), "Keep username/password login available before any future SSO provider is configured."),
			configRow("Local demo login", page.AuthStatus.LocalDemoLogin, "review", "Local demo sign-in must remain disabled in production."),
			configRow("SSO/OIDC", page.AuthStatus.SSOStatus, "future", "Do not add OIDC endpoints until provider discovery, callback, validation, and role mapping are implemented."),
			configRow("CSRF policy", page.AuthStatus.CSRFPolicy, "configured", "Unsafe cookie-authenticated requests must include the form token."),
		}
	case "deployment":
		return []operationsConfigRow{
			configRow("Environment", page.EnvironmentLabel, configStatus(page.EnvironmentLabel), "Set a non-secret publication environment label."),
			configRow("Active feed version", firstNonEmpty(page.ActiveFeedVersion, "missing"), configStatus(page.ActiveFeedVersion), "Import and activate a schedule before realtime review."),
			configRow("Telemetry stale threshold", page.StaleThreshold.String(), "configured", "Keep stale and future timestamp rules explicit for vehicle data."),
			configRow("Maintenance status", page.Maintenance.OverallStatus, page.Maintenance.OverallStatus, "Use Maintenance for backup, restore, support bundle, and host readiness details."),
			configRow("Reliability status", firstNonEmpty(page.Reliability.OverallStatus, page.ReliabilityError, "not available"), configStatus(firstNonEmpty(page.Reliability.OverallStatus, page.ReliabilityError)), "Review private reliability snapshots without making SLA or uptime claims."),
		}
	case "advanced":
		return []operationsConfigRow{
			configRow("Operations JSON", "/admin/operations.json", "private", "Export only for internal review; do not publish private diagnostics."),
			configRow("Checklist JSON", "/admin/operations/checklist.json", "private", "Use checklist JSON for operator review without changing consumer tracker state."),
			configRow("Route inventory", "make audit-operations-route-inventory", "local", "Run the route audit locally before release."),
			configRow("Product language", "make audit-product-language", "local", "Keep compliance, SSO, hosted service, SLA, AVL, and ETA claims out of product text."),
		}
	default:
		return []operationsConfigRow{
			configRow("Settings moved", "publication metadata lives under Public Feed URLs", "configured", "Use focused config pages instead of the setup diagnostics page for settings."),
			configRow("Dashboard separation", "dashboard shows top issues; config pages hold settings", "configured", "Keep diagnostic tables out of the first screen."),
			configRow("Advanced details", "collapsed or JSON-only", "review", "Open advanced details only for support and release-gate review."),
		}
	}
}

func configRow(label, value, status, nextAction string) operationsConfigRow {
	return operationsConfigRow{
		Label:      label,
		Value:      firstNonEmpty(value, "missing"),
		Status:     firstNonEmpty(status, "unknown"),
		NextAction: firstNonEmpty(nextAction, "Review this setting before release."),
	}
}

func configStatus(value string) string {
	switch firstNonEmpty(value) {
	case "", "missing", "unknown", "not available", "unavailable":
		return "missing"
	default:
		return "configured"
	}
}

func configRoleVisibilitySignal(page operationsPage) string {
	if len(page.SetupWizard.RoleVisibility) == 0 {
		return "no role visibility rows available"
	}
	return firstNonEmpty(page.SetupWizard.RoleVisibility[0].CurrentSignal, page.SetupWizard.RoleVisibility[0].Status)
}
