package main

type operationsNavGroup struct {
	Label string
	Items []operationsNavItem
}

type operationsNavItem struct {
	Label                string
	Href                 string
	Section              string
	ExternalAdminSurface bool
	Current              bool
}

func operationsNavGroups(currentSection string) []operationsNavGroup {
	groups := []operationsNavGroup{
		{
			Label: "Start Here",
			Items: []operationsNavItem{
				{Label: "Start Here", Href: "/admin/operations", Section: "dashboard"},
				{Label: "Private Launchpad", Href: "/admin/operations/launchpad", Section: "launchpad"},
				{Label: "Agency Setup", Href: "/admin/operations/setup-wizard", Section: "setup-wizard"},
				{Label: "Advanced Setup Details", Href: "/admin/operations/setup", Section: "setup"},
			},
		},
		{
			Label: "Schedule",
			Items: []operationsNavItem{
				{Label: "GTFS Import", Href: "/admin/operations/gtfs-import", Section: "gtfs-import"},
				{Label: "GTFS Studio", Href: "/admin/gtfs-studio", Section: "gtfs-studio", ExternalAdminSurface: true},
				{Label: "Feeds", Href: "/admin/operations/feeds", Section: "feeds"},
				{Label: "Feed Health", Href: "/admin/operations/feed-health", Section: "feed-health"},
				{Label: "GTFS Quality", Href: "/admin/operations/gtfs-quality", Section: "gtfs-quality"},
				{Label: "Validator Health", Href: "/admin/operations/validation-health", Section: "validation-health"},
			},
		},
		{
			Label: "Realtime",
			Items: []operationsNavItem{
				{Label: "Telemetry", Href: "/admin/operations/telemetry", Section: "telemetry"},
				{Label: "Device Credentials", Href: "/admin/operations/devices", Section: "devices"},
				{Label: "Telemetry Simulator", Href: "/admin/operations/telemetry-simulator", Section: "telemetry-simulator"},
				{Label: "Alerts", Href: "/admin/alerts/console", Section: "alerts", ExternalAdminSurface: true},
			},
		},
		{
			Label: "Connectors",
			Items: []operationsNavItem{
				{Label: "Connector Hub", Href: "/admin/operations/connectors", Section: "connectors"},
				{Label: "Connector Tests", Href: "/admin/operations/connectors/tests", Section: "connector-tests"},
			},
		},
		{
			Label: "Health",
			Items: []operationsNavItem{
				{Label: "Readiness", Href: "/admin/operations/readiness", Section: "readiness"},
				{Label: "Checklist", Href: "/admin/operations/checklist", Section: "checklist"},
				{Label: "Reliability", Href: "/admin/operations/reliability", Section: "reliability"},
			},
		},
		{
			Label: "Maintain",
			Items: []operationsNavItem{
				{Label: "Maintenance", Href: "/admin/operations/maintenance", Section: "maintenance"},
			},
		},
		{
			Label: "Learn",
			Items: []operationsNavItem{
				{Label: "Help", Href: "/admin/operations/help", Section: "help"},
				{Label: "Consumers", Href: "/admin/operations/consumers", Section: "consumers"},
				{Label: "Evidence", Href: "/admin/operations/evidence", Section: "evidence"},
			},
		},
	}

	current := normalizeOperationsNavSection(currentSection)
	for groupIndex := range groups {
		for itemIndex := range groups[groupIndex].Items {
			groups[groupIndex].Items[itemIndex].Current = groups[groupIndex].Items[itemIndex].Section == current
		}
	}
	return groups
}

func normalizeOperationsNavSection(section string) string {
	if section == "" {
		return "dashboard"
	}
	return section
}

func operationsPageTitle(section string) string {
	switch normalizeOperationsNavSection(section) {
	case "dashboard":
		return "Agency Operations Cockpit / Start Here"
	case "launchpad":
		return "Private Agency Launchpad"
	case "setup-wizard":
		return "Agency Setup"
	case "setup":
		return "Advanced Setup Details"
	case "gtfs-import":
		return "Browser GTFS Import"
	case "gtfs-quality":
		return "GTFS Quality Triage"
	case "feeds":
		return "Feed URLs And Validation"
	case "feed-health":
		return "Feed Health Dashboard"
	case "validation-health":
		return "Validator Health"
	case "readiness":
		return "Readiness Checklist V2"
	case "checklist":
		return "Private Operator Checklist"
	case "reliability":
		return "Operations Reliability"
	case "maintenance":
		return "Maintenance Center"
	case "telemetry":
		return "Telemetry Freshness"
	case "telemetry-simulator":
		return "Telemetry Simulator Guide"
	case "devices":
		return "Device Credentials"
	case "connectors":
		return "Connector Hub"
	case "connector-tests":
		return "Connector Test Instructions"
	case "consumers":
		return "Consumer Preparation Tracker"
	case "evidence":
		return "Evidence Links And Runbooks"
	case "help":
		return "Operations Console Help"
	default:
		return "Operations Console"
	}
}
