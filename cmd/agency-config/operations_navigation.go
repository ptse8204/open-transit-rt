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
			Label: "Start",
			Items: []operationsNavItem{
				{Label: "Dashboard", Href: "/admin/operations", Section: "dashboard"},
				{Label: "Launchpad", Href: "/admin/operations/launchpad", Section: "launchpad"},
				{Label: "Setup Wizard", Href: "/admin/operations/setup-wizard", Section: "setup-wizard"},
				{Label: "Setup", Href: "/admin/operations/setup", Section: "setup"},
			},
		},
		{
			Label: "GTFS and feeds",
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
			Label: "Realtime operations",
			Items: []operationsNavItem{
				{Label: "Telemetry", Href: "/admin/operations/telemetry", Section: "telemetry"},
				{Label: "Devices", Href: "/admin/operations/devices", Section: "devices"},
				{Label: "Simulator", Href: "/admin/operations/telemetry-simulator", Section: "telemetry-simulator"},
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
			Label: "Readiness and diagnostics",
			Items: []operationsNavItem{
				{Label: "Readiness", Href: "/admin/operations/readiness", Section: "readiness"},
				{Label: "Checklist", Href: "/admin/operations/checklist", Section: "checklist"},
				{Label: "Reliability", Href: "/admin/operations/reliability", Section: "reliability"},
			},
		},
		{
			Label: "Records and boundaries",
			Items: []operationsNavItem{
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
