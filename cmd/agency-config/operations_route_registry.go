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
	{ID: "start", Label: "Start"},
	{ID: "setup", Label: "Agency"},
	{ID: "config", Label: "Config"},
	{ID: "gtfs", Label: "Schedule"},
	{ID: "feeds", Label: "Feeds"},
	{ID: "realtime", Label: "Realtime"},
	{ID: "vehicles", Label: "Vehicles"},
	{ID: "connectors", Label: "Connectors"},
	{ID: "readiness", Label: "Readiness"},
	{ID: "maintenance", Label: "Maintain"},
	{ID: "admin", Label: "Admin"},
	{ID: "help", Label: "Support"},
}

var operationsRouteRegistry = []operationsRouteMeta{
	{Section: "dashboard", Path: "/admin/operations", JSONPath: "/admin/operations.json", NavLabel: "Start", PageTitle: "Start", GroupID: "start", Methods: []string{"GET"}, NoStore: true},
	{Section: "launchpad", Path: "/admin/operations/launchpad", JSONPath: "/admin/operations/launchpad.json", NavLabel: "Find Work", PageTitle: "Find Work", GroupID: "start", Methods: []string{"GET"}, NoStore: true},
	{Section: "setup-wizard", Path: "/admin/operations/setup-wizard", JSONPath: "/admin/operations/setup-wizard.json", NavLabel: "Agency Setup", PageTitle: "Agency Setup", GroupID: "setup", Methods: []string{"GET"}, NoStore: true},
	{Section: "setup", Path: "/admin/operations/setup", NavLabel: "Setup Diagnostics", PageTitle: "Setup Diagnostics", GroupID: "setup", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "access", Path: "/admin/operations/access", JSONPath: "/admin/operations/access.json", NavLabel: "Access & Roles", PageTitle: "Access & Roles", GroupID: "setup", Methods: []string{"GET"}, NoStore: true},
	{Section: "config", Path: "/admin/operations/config", NavLabel: "Config Overview", PageTitle: "Focused Config", GroupID: "config", Methods: []string{"GET"}, NoStore: true},
	{Section: "config-agency", Path: "/admin/operations/config/agency", NavLabel: "Agency Profile", PageTitle: "Agency Profile", GroupID: "config", Methods: []string{"GET"}, NoStore: true},
	{Section: "config-feeds", Path: "/admin/operations/config/feeds", NavLabel: "Public Feed URLs", PageTitle: "Public Feed URLs", GroupID: "config", Methods: []string{"GET"}, NoStore: true},
	{Section: "config-auth", Path: "/admin/operations/config/auth", NavLabel: "Login Settings", PageTitle: "Login Settings", GroupID: "config", Methods: []string{"GET"}, NoStore: true},
	{Section: "config-deployment", Path: "/admin/operations/config/deployment", NavLabel: "Deployment Settings", PageTitle: "Deployment Settings", GroupID: "config", Methods: []string{"GET"}, NoStore: true},
	{Section: "config-advanced", Path: "/admin/operations/config/advanced", NavLabel: "Advanced Settings", PageTitle: "Advanced Settings", GroupID: "config", Methods: []string{"GET"}, NoStore: true},
	{Section: "gtfs-workbench", Path: "/admin/operations/gtfs-workbench", JSONPath: "/admin/operations/gtfs-workbench.json", NavLabel: "Schedule Review", PageTitle: "Schedule Review", GroupID: "gtfs", Methods: []string{"GET"}, NoStore: true},
	{Section: "gtfs-import", Path: "/admin/operations/gtfs-import", NavLabel: "Import Schedule", PageTitle: "Import Schedule", GroupID: "gtfs", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "gtfs-studio", Path: "/admin/gtfs-studio", NavLabel: "Draft Schedule Editor", GroupID: "gtfs", Methods: []string{"GET"}, ExternalAdminSurface: true},
	{Section: "gtfs-quality", Path: "/admin/operations/gtfs-quality", NavLabel: "Schedule Quality", PageTitle: "Schedule Quality", GroupID: "gtfs", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "feeds", Path: "/admin/operations/feeds", NavLabel: "Feed URLs", PageTitle: "Feed URLs", GroupID: "feeds", Methods: []string{"GET"}, NoStore: true},
	{Section: "feed-health", Path: "/admin/operations/feed-health", JSONPath: "/admin/operations/feed-health.json", NavLabel: "Check Feeds", PageTitle: "Check Feeds", GroupID: "feeds", Methods: []string{"GET"}, NoStore: true},
	{Section: "validation-center", Path: "/admin/operations/validation-center", JSONPath: "/admin/operations/validation-center.json", NavLabel: "Validation", PageTitle: "Validation", GroupID: "feeds", Methods: []string{"GET"}, NoStore: true},
	{Section: "validation-health", Path: "/admin/operations/validation-health", JSONPath: "/admin/operations/validation-health.json", NavLabel: "Validators", PageTitle: "Validators", GroupID: "feeds", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "realtime", Path: "/admin/operations/realtime", JSONPath: "/admin/operations/realtime.json", NavLabel: "Realtime Output", PageTitle: "Realtime Output", GroupID: "realtime", Methods: []string{"GET"}, NoStore: true},
	{Section: "prediction-lab", Path: "/admin/operations/prediction-lab", JSONPath: "/admin/operations/prediction-lab.json", NavLabel: "Trip Updates", PageTitle: "Trip Updates", GroupID: "realtime", Methods: []string{"GET"}, NoStore: true},
	{Section: "alerts", Path: "/admin/alerts/console", NavLabel: "Alerts", GroupID: "realtime", Methods: []string{"GET"}, ExternalAdminSurface: true},
	{Section: "telemetry", Path: "/admin/operations/telemetry", NavLabel: "Telemetry", PageTitle: "Telemetry", GroupID: "vehicles", Methods: []string{"GET"}, NoStore: true},
	{Section: "devices", Path: "/admin/operations/devices", NavLabel: "Devices", PageTitle: "Devices", GroupID: "vehicles", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "telemetry-simulator", Path: "/admin/operations/telemetry-simulator", JSONPath: "/admin/operations/telemetry-simulator.json", NavLabel: "Simulator", PageTitle: "Simulator", GroupID: "vehicles", Methods: []string{"GET"}, NoStore: true},
	{Section: "connectors", Path: "/admin/operations/connectors", JSONPath: "/admin/operations/connectors.json", NavLabel: "Connectors", PageTitle: "Connectors", GroupID: "connectors", Methods: []string{"GET"}, NoStore: true},
	{Section: "vehicle-avl-setup", Path: "/admin/operations/connectors/vehicle-avl", JSONPath: "/admin/operations/connectors/vehicle-avl.json", NavLabel: "Vehicle / AVL Setup", PageTitle: "Vehicle / AVL Setup", GroupID: "connectors", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "prediction-setup", Path: "/admin/operations/connectors/prediction", JSONPath: "/admin/operations/connectors/prediction.json", NavLabel: "Prediction Setup", PageTitle: "Prediction Setup", GroupID: "connectors", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "validator-setup", Path: "/admin/operations/connectors/validators", JSONPath: "/admin/operations/connectors/validators.json", NavLabel: "Validator Setup", PageTitle: "Validator Setup", GroupID: "connectors", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "connector-workbench", Path: "/admin/operations/connectors/workbench", JSONPath: "/admin/operations/connectors/workbench.json", NavLabel: "Connector Workbench", PageTitle: "Connector Workbench", GroupID: "connectors", Methods: []string{"GET"}, NoStore: true},
	{Section: "connector-tests", Path: "/admin/operations/connectors/tests", JSONPath: "/admin/operations/connectors/tests.json", NavLabel: "Connector Checks", PageTitle: "Connector Checks", GroupID: "connectors", Methods: []string{"GET"}, NoStore: true},
	{Section: "admin-sessions", Path: "/admin/operations/admin/sessions", JSONPath: "/admin/operations/admin/sessions.json", NavLabel: "Login & Sessions", PageTitle: "Login & Sessions", GroupID: "admin", Methods: []string{"GET"}, NoStore: true},
	{Section: "admin-users", Path: "/admin/operations/admin/users", JSONPath: "/admin/operations/admin/users.json", NavLabel: "Users & Roles", PageTitle: "Users & Roles", GroupID: "admin", Methods: []string{"GET", "POST"}, NoStore: true},
	{Section: "readiness", Path: "/admin/operations/readiness", JSONPath: "/admin/operations/readiness.json", NavLabel: "Readiness", PageTitle: "Readiness", GroupID: "readiness", Methods: []string{"GET"}, NoStore: true},
	{Section: "checklist", Path: "/admin/operations/checklist", JSONPath: "/admin/operations/checklist.json", NavLabel: "Checklist", PageTitle: "Checklist", GroupID: "readiness", Methods: []string{"GET"}, NoStore: true},
	{Section: "consumers", Path: "/admin/operations/consumers", NavLabel: "External Sharing Prep", PageTitle: "External Sharing Prep", GroupID: "readiness", Methods: []string{"GET"}, NoStore: true},
	{Section: "maintenance", Path: "/admin/operations/maintenance", JSONPath: "/admin/operations/maintenance.json", NavLabel: "Maintenance", PageTitle: "Maintenance", GroupID: "maintenance", Methods: []string{"GET"}, NoStore: true},
	{Section: "reliability", Path: "/admin/operations/reliability", JSONPath: "/admin/operations/reliability.json", NavLabel: "Reliability", PageTitle: "Reliability", GroupID: "maintenance", Methods: []string{"GET"}, NoStore: true},
	{Section: "help", Path: "/admin/operations/help", JSONPath: "/admin/operations/help.json", NavLabel: "Help", PageTitle: "Help", GroupID: "help", Methods: []string{"GET"}, NoStore: true},
	{Section: "audit", Path: "/admin/operations/audit", JSONPath: "/admin/operations/audit.json", NavLabel: "Audit Log", PageTitle: "Audit Log", GroupID: "help", Methods: []string{"GET"}, NoStore: true},
	{Section: "evidence", Path: "/admin/operations/evidence", NavLabel: "Sharing Limits", PageTitle: "Sharing Limits", GroupID: "help", Methods: []string{"GET"}, NoStore: true},
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
