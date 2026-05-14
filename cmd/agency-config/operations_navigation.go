package main

type operationsNavGroup struct {
	ID    string
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
