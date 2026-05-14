package main

type operationsRouteGroupMeta struct {
	ID    string
	Label string
}

type operationsRouteMeta struct {
	Section              string
	Path                 string
	JSONPath             string
	NavLabel             string
	PageTitle            string
	GroupID              string
	Methods              []string
	NoStore              bool
	ExternalAdminSurface bool
}

type operationsCommandRouteMeta struct {
	Section string
	Path    string
	Method  string
	NoStore bool
}

var operationsRouteGroupRegistry = []operationsRouteGroupMeta{
	{ID: "start", Label: "Start Here"},
	{ID: "schedule", Label: "Schedule"},
	{ID: "realtime", Label: "Realtime"},
	{ID: "connectors", Label: "Connectors"},
	{ID: "health", Label: "Health"},
	{ID: "maintain", Label: "Maintain"},
	{ID: "learn", Label: "Learn"},
}

var operationsRouteRegistry = []operationsRouteMeta{
	{Section: "dashboard", Path: "/admin/operations", JSONPath: "/admin/operations.json", NavLabel: "Start Here", PageTitle: "Agency Operations Cockpit / Start Here", GroupID: "start", Methods: []string{"GET"}, NoStore: true},
	{Section: "launchpad", Path: "/admin/operations/launchpad", JSONPath: "/admin/operations/launchpad.json", NavLabel: "Private Launchpad", PageTitle: "Private Agency Launchpad", GroupID: "start", Methods: []string{"GET"}, NoStore: true},
	{Section: "setup-wizard", Path: "/admin/operations/setup-wizard", JSONPath: "/admin/operations/setup-wizard.json", NavLabel: "Agency Setup", PageTitle: "Agency Setup", GroupID: "start", Methods: []string{"GET"}, NoStore: true},
	{Section: "setup", Path: "/admin/operations/setup", NavLabel: "Advanced Setup Details", PageTitle: "Advanced Setup Details", GroupID: "start", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "gtfs-workbench", Path: "/admin/operations/gtfs-workbench", JSONPath: "/admin/operations/gtfs-workbench.json", NavLabel: "GTFS Workbench", PageTitle: "GTFS Workbench", GroupID: "schedule", Methods: []string{"GET"}, NoStore: true},
	{Section: "gtfs-import", Path: "/admin/operations/gtfs-import", NavLabel: "Import Schedule ZIP", PageTitle: "Browser GTFS Import", GroupID: "schedule", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "gtfs-studio", Path: "/admin/gtfs-studio", NavLabel: "Draft Schedule Editor", GroupID: "schedule", Methods: []string{"GET"}, ExternalAdminSurface: true},
	{Section: "feeds", Path: "/admin/operations/feeds", NavLabel: "Feed Links & Health", PageTitle: "Feed URLs And Validation", GroupID: "schedule", Methods: []string{"GET"}, NoStore: true},
	{Section: "feed-health", Path: "/admin/operations/feed-health", JSONPath: "/admin/operations/feed-health.json", NavLabel: "Schedule Health", PageTitle: "Feed Health Dashboard", GroupID: "schedule", Methods: []string{"GET"}, NoStore: true},
	{Section: "gtfs-quality", Path: "/admin/operations/gtfs-quality", NavLabel: "Schedule Quality", PageTitle: "GTFS Quality Triage", GroupID: "schedule", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "validation-health", Path: "/admin/operations/validation-health", JSONPath: "/admin/operations/validation-health.json", NavLabel: "Schedule Validation", PageTitle: "Validator Health", GroupID: "schedule", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "realtime", Path: "/admin/operations/realtime", JSONPath: "/admin/operations/realtime.json", NavLabel: "Realtime Center", PageTitle: "Realtime Operations Center", GroupID: "realtime", Methods: []string{"GET"}, NoStore: true},
	{Section: "prediction-lab", Path: "/admin/operations/prediction-lab", JSONPath: "/admin/operations/prediction-lab.json", NavLabel: "Prediction & ETA Lab", PageTitle: "Prediction & ETA Lab", GroupID: "realtime", Methods: []string{"GET"}, NoStore: true},
	{Section: "telemetry", Path: "/admin/operations/telemetry", NavLabel: "Telemetry", PageTitle: "Telemetry Freshness", GroupID: "realtime", Methods: []string{"GET"}, NoStore: true},
	{Section: "devices", Path: "/admin/operations/devices", NavLabel: "Device Credentials", PageTitle: "Device Credentials", GroupID: "realtime", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "telemetry-simulator", Path: "/admin/operations/telemetry-simulator", JSONPath: "/admin/operations/telemetry-simulator.json", NavLabel: "Telemetry Simulator", PageTitle: "Telemetry Simulator Guide", GroupID: "realtime", Methods: []string{"GET"}, NoStore: true},
	{Section: "alerts", Path: "/admin/alerts/console", NavLabel: "Alerts", GroupID: "realtime", Methods: []string{"GET"}, ExternalAdminSurface: true},
	{Section: "connectors", Path: "/admin/operations/connectors", JSONPath: "/admin/operations/connectors.json", NavLabel: "Connector Hub", PageTitle: "Connector Hub", GroupID: "connectors", Methods: []string{"GET"}, NoStore: true},
	{Section: "connector-workbench", Path: "/admin/operations/connectors/workbench", JSONPath: "/admin/operations/connectors/workbench.json", NavLabel: "Connector Workbench", PageTitle: "Connector Workbench", GroupID: "connectors", Methods: []string{"GET"}, NoStore: true},
	{Section: "connector-tests", Path: "/admin/operations/connectors/tests", JSONPath: "/admin/operations/connectors/tests.json", NavLabel: "Connector Tests", PageTitle: "Connector Test Instructions", GroupID: "connectors", Methods: []string{"GET"}, NoStore: true},
	{Section: "validation-center", Path: "/admin/operations/validation-center", JSONPath: "/admin/operations/validation-center.json", NavLabel: "Validation Center", PageTitle: "Feed Health And Validation Center", GroupID: "health", Methods: []string{"GET"}, NoStore: true},
	{Section: "readiness", Path: "/admin/operations/readiness", JSONPath: "/admin/operations/readiness.json", NavLabel: "Readiness", PageTitle: "Readiness Checklist V2", GroupID: "health", Methods: []string{"GET"}, NoStore: true},
	{Section: "checklist", Path: "/admin/operations/checklist", JSONPath: "/admin/operations/checklist.json", NavLabel: "Checklist", PageTitle: "Private Operator Checklist", GroupID: "health", Methods: []string{"GET"}, NoStore: true},
	{Section: "reliability", Path: "/admin/operations/reliability", JSONPath: "/admin/operations/reliability.json", NavLabel: "Reliability", PageTitle: "Operations Reliability", GroupID: "health", Methods: []string{"GET"}, NoStore: true},
	{Section: "maintenance", Path: "/admin/operations/maintenance", JSONPath: "/admin/operations/maintenance.json", NavLabel: "Maintenance", PageTitle: "Maintenance Center", GroupID: "maintain", Methods: []string{"GET"}, NoStore: true},
	{Section: "access", Path: "/admin/operations/access", JSONPath: "/admin/operations/access.json", NavLabel: "Access & Roles", PageTitle: "Access & Roles", GroupID: "maintain", Methods: []string{"GET"}, NoStore: true},
	{Section: "audit", Path: "/admin/operations/audit", JSONPath: "/admin/operations/audit.json", NavLabel: "Audit Log", PageTitle: "Audit Log", GroupID: "maintain", Methods: []string{"GET"}, NoStore: true},
	{Section: "help", Path: "/admin/operations/help", JSONPath: "/admin/operations/help.json", NavLabel: "Help", PageTitle: "Operations Console Help", GroupID: "learn", Methods: []string{"GET"}, NoStore: true},
	{Section: "consumers", Path: "/admin/operations/consumers", NavLabel: "Consumers", PageTitle: "Consumer Preparation Tracker", GroupID: "learn", Methods: []string{"GET"}, NoStore: true},
	{Section: "evidence", Path: "/admin/operations/evidence", NavLabel: "Evidence", PageTitle: "Evidence Links And Runbooks", GroupID: "learn", Methods: []string{"GET"}, NoStore: true},
}

var operationsCommandRouteRegistry = []operationsCommandRouteMeta{
	{Section: "validation-health", Path: "/admin/operations/validation-health/refresh.json", Method: "POST", NoStore: true},
}

func operationsRouteGroups() []operationsRouteGroupMeta {
	groups := make([]operationsRouteGroupMeta, len(operationsRouteGroupRegistry))
	copy(groups, operationsRouteGroupRegistry)
	return groups
}

func operationsRoutes() []operationsRouteMeta {
	routes := make([]operationsRouteMeta, len(operationsRouteRegistry))
	copy(routes, operationsRouteRegistry)
	return routes
}

func operationsCommandRoutes() []operationsCommandRouteMeta {
	routes := make([]operationsCommandRouteMeta, len(operationsCommandRouteRegistry))
	copy(routes, operationsCommandRouteRegistry)
	return routes
}

func operationsRouteBySection(section string) (operationsRouteMeta, bool) {
	section = normalizeOperationsNavSection(section)
	for _, route := range operationsRouteRegistry {
		if route.Section == section {
			return route, true
		}
	}
	return operationsRouteMeta{}, false
}

func operationsRouteByPath(path string) (operationsRouteMeta, bool) {
	for _, route := range operationsRouteRegistry {
		if route.Path == path || route.JSONPath == path {
			return route, true
		}
	}
	return operationsRouteMeta{}, false
}

func operationsCanonicalHTMLRoutes() []operationsRouteMeta {
	var routes []operationsRouteMeta
	for _, route := range operationsRouteRegistry {
		if route.ExternalAdminSurface || route.Path == "" {
			continue
		}
		routes = append(routes, route)
	}
	return routes
}

func operationsCanonicalJSONRoutes() []string {
	var routes []string
	for _, route := range operationsRouteRegistry {
		if route.JSONPath != "" {
			routes = append(routes, route.JSONPath)
		}
	}
	return routes
}

func operationsExternalAdminSurfaceRoutes() []operationsRouteMeta {
	var routes []operationsRouteMeta
	for _, route := range operationsRouteRegistry {
		if route.ExternalAdminSurface {
			routes = append(routes, route)
		}
	}
	return routes
}
