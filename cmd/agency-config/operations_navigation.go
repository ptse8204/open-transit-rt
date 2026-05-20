package main

type operationsNavGroup struct {
	ID      string
	Label   string
	Current bool
	Items   []operationsNavItem
}

type operationsNavItem struct {
	Label                string
	Href                 string
	Section              string
	ExternalAdminSurface bool
	Current              bool
}

func operationsNavGroups(currentSection string) []operationsNavGroup {
	groups := make([]operationsNavGroup, 0, len(operationsRouteGroupRegistry))
	groupIndexes := make(map[string]int, len(operationsRouteGroupRegistry))
	for _, group := range operationsRouteGroups() {
		groupIndexes[group.ID] = len(groups)
		groups = append(groups, operationsNavGroup{ID: group.ID, Label: group.Label})
	}

	current := normalizeOperationsNavSection(currentSection)
	for _, route := range operationsRoutes() {
		if route.NavLabel == "" {
			continue
		}
		index, ok := groupIndexes[route.GroupID]
		if !ok {
			continue
		}
		groups[index].Items = append(groups[index].Items, operationsNavItem{
			Label:                route.NavLabel,
			Href:                 route.Path,
			Section:              route.Section,
			ExternalAdminSurface: route.ExternalAdminSurface,
			Current:              route.Section == current,
		})
	}
	for groupIndex := range groups {
		for itemIndex := range groups[groupIndex].Items {
			groups[groupIndex].Items[itemIndex].Current = groups[groupIndex].Items[itemIndex].Section == current
			if groups[groupIndex].Items[itemIndex].Current {
				groups[groupIndex].Current = true
			}
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
	if route, ok := operationsRouteBySection(section); ok && route.PageTitle != "" && !route.ExternalAdminSurface {
		return route.PageTitle
	}
	return "Operations Console"
}

func operationsPageNextAction(section string) string {
	switch normalizeOperationsNavSection(section) {
	case "dashboard":
		return "Follow the workflow in order. Open the first missing, blocked, or needs-review step."
	case "launchpad":
		return "Choose the workflow section that matches the current blocker and follow its linked console page."
	case "setup-wizard", "setup":
		return "Confirm agency metadata, feed URL metadata, license, contact, and role visibility before importing or sharing feed links."
	case "gtfs-workbench", "gtfs-import", "gtfs-quality":
		return "Review the active schedule, import or update GTFS when needed, then check validation and feed health."
	case "feeds", "feed-health":
		return "Review all five configured feed paths, then fix missing URLs, stale data, or validation gaps."
	case "validation-center", "validation-health":
		return "Review validation blockers first, then run only server-owned allowlisted validation actions when configured."
	case "realtime":
		return "Review stale, unmatched, suppressed, and low-confidence vehicles before relying on realtime feed output."
	case "telemetry", "devices", "telemetry-simulator":
		return "Review device bindings and freshness, then use safe simulator or deployment-owned telemetry guidance."
	case "prediction-lab":
		return "Review withheld Trip Updates reasons and keep ETA-like output behind the prediction adapter boundary."
	case "connectors", "connector-workbench", "connector-tests":
		return "Choose a connector category, review redaction and fail-closed behavior, then have an administrator or integrator run synthetic conformance checks."
	case "admin-users":
		return "Review users scoped to this agency, assign the smallest needed role, and generate password reset links only when a recipient is ready."
	case "readiness", "checklist", "consumers":
		return "Use readiness rows to prepare for future review without changing consumer status or claiming outside approval."
	case "maintenance", "reliability":
		return "Review routine maintenance and reliability rows, then decide whether a deployment owner needs to run local diagnostics."
	case "access", "audit", "evidence":
		return "Use these support views to understand roles, recent audit metadata, and evidence boundaries without exposing raw private data."
	case "help":
		return "Pick the role-based guide or recovery topic that matches the user in front of the console."
	default:
		return "Return to Start and choose the next visible action."
	}
}
