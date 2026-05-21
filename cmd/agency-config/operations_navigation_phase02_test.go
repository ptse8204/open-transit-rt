package main

import (
	"encoding/json"
	"html"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"open-transit-rt/internal/auth"
)

func TestOperationsPhase02RegisteredRoutesArePrivateReachableAndNoStore(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02ReadOnlyPrincipal()})
	adminHandler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	unauthenticated := phase02OperationsHandler(t, authRejectAll{})

	for _, route := range operationsCanonicalHTMLRoutes() {
		if !phase02RouteAllows(route.Methods, http.MethodGet) {
			continue
		}
		t.Run(route.Path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, route.Path, nil)
			rr := httptest.NewRecorder()
			phase02HandlerForRoute(route, handler, adminHandler).ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("%s status = %d, want 200: %s", route.Path, rr.Code, rr.Body.String())
			}
			if route.NoStore {
				if got := rr.Header().Get("Cache-Control"); got != "no-store" {
					t.Fatalf("%s Cache-Control = %q, want no-store", route.Path, got)
				}
			}
			body := rr.Body.String()
			for _, want := range []string{
				`<header class="operations-header" role="banner">`,
				`class="app-kicker">Private`,
				`<div class="operations-frame">`,
				`<nav id="operations-nav" class="operations-nav" aria-label="Operations Console sections">`,
				`<main id="operations-main" tabindex="-1" aria-labelledby="operations-page-title">`,
				`<script src="/admin/operations/assets/operations.js" defer></script>`,
			} {
				if !strings.Contains(body, want) {
					t.Fatalf("%s rendered body missing shared private shell marker %q: %s", route.Path, want, body)
				}
			}

			req = httptest.NewRequest(http.MethodGet, route.Path, nil)
			rr = httptest.NewRecorder()
			unauthenticated.ServeHTTP(rr, req)
			if rr.Code != http.StatusUnauthorized {
				t.Fatalf("unauthenticated %s status = %d, want 401", route.Path, rr.Code)
			}

			publicPath := "/public" + strings.TrimPrefix(route.Path, "/admin")
			req = httptest.NewRequest(http.MethodGet, publicPath, nil)
			rr = httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusNotFound {
				t.Fatalf("%s public mirror status = %d, want 404", publicPath, rr.Code)
			}
		})
	}

	for _, path := range operationsCanonicalJSONRoutes() {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()
			route, _ := operationsRouteByPath(path)
			phase02HandlerForRoute(route, handler, adminHandler).ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
			}
			if got := rr.Header().Get("Cache-Control"); got != "no-store" {
				t.Fatalf("%s Cache-Control = %q, want no-store", path, got)
			}
			var decoded map[string]any
			if err := json.Unmarshal(rr.Body.Bytes(), &decoded); err != nil {
				t.Fatalf("%s did not render valid JSON: %v: %s", path, err, rr.Body.String())
			}
			if len(decoded) == 0 {
				t.Fatalf("%s rendered empty JSON object", path)
			}

			req = httptest.NewRequest(http.MethodGet, path, nil)
			rr = httptest.NewRecorder()
			unauthenticated.ServeHTTP(rr, req)
			if rr.Code != http.StatusUnauthorized {
				t.Fatalf("unauthenticated %s status = %d, want 401", path, rr.Code)
			}

			publicPath := "/public" + strings.TrimPrefix(path, "/admin")
			req = httptest.NewRequest(http.MethodGet, publicPath, nil)
			rr = httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusNotFound {
				t.Fatalf("%s public mirror status = %d, want 404", publicPath, rr.Code)
			}
		})
	}
}

func TestOperationsPhase02NavigationMatchesRouteRegistryAndRenderedTemplate(t *testing.T) {
	groups := operationsNavGroups("dashboard")
	registryGroups := operationsRouteGroups()
	if len(groups) != len(registryGroups) {
		t.Fatalf("navigation groups = %d, want %d", len(groups), len(registryGroups))
	}

	expectedItems := map[string][]operationsNavItem{}
	for _, route := range operationsRoutes() {
		if route.NavLabel == "" {
			continue
		}
		expectedItems[route.GroupID] = append(expectedItems[route.GroupID], operationsNavItem{
			Label:                route.NavLabel,
			Href:                 route.Path,
			Section:              route.Section,
			ExternalAdminSurface: route.ExternalAdminSurface,
			Current:              route.Section == "dashboard",
		})
	}

	for i, group := range groups {
		wantGroup := registryGroups[i]
		if group.ID != wantGroup.ID || group.Label != wantGroup.Label {
			t.Fatalf("group[%d] = %+v, want id=%q label=%q", i, group, wantGroup.ID, wantGroup.Label)
		}
		wantItems := expectedItems[group.ID]
		if len(group.Items) != len(wantItems) {
			t.Fatalf("group %s item count = %d, want %d: %+v", group.ID, len(group.Items), len(wantItems), group.Items)
		}
		for itemIndex, item := range group.Items {
			wantItem := wantItems[itemIndex]
			if item != wantItem {
				t.Fatalf("group %s item[%d] = %+v, want %+v", group.ID, itemIndex, item, wantItem)
			}
		}
	}

	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02ReadOnlyPrincipal()})
	req := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("dashboard status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()

	previousGroupIndex := -1
	for _, group := range registryGroups {
		groupMarker := `id="nav-group-` + group.ID + `" class="nav-group-label">` + html.EscapeString(group.Label)
		groupIndex := strings.Index(body, groupMarker)
		if groupIndex < 0 {
			t.Fatalf("rendered navigation missing group %q marker %q: %s", group.ID, groupMarker, body)
		}
		if groupIndex <= previousGroupIndex {
			t.Fatalf("rendered navigation group %q appeared out of registry order", group.ID)
		}
		previousGroupIndex = groupIndex
	}

	for _, route := range operationsRoutes() {
		if route.NavLabel == "" {
			continue
		}
		escapedLabel := html.EscapeString(route.NavLabel)
		want := `href="` + route.Path + `"`
		if !strings.Contains(body, want) {
			t.Fatalf("rendered navigation missing href for %s: %s", route.Path, body)
		}
		if !strings.Contains(body, escapedLabel) {
			t.Fatalf("rendered navigation missing label %q for %s: %s", escapedLabel, route.Path, body)
		}
		if route.ExternalAdminSurface {
			wantSurface := `href="` + route.Path + `">` + escapedLabel + ` <span class="nav-surface">`
			if !strings.Contains(body, wantSurface) {
				t.Fatalf("rendered navigation missing external surface marker for %s: %s", route.Path, body)
			}
		}
	}
	if got := strings.Count(body, `aria-current="page"`); got != 1 {
		t.Fatalf("dashboard aria-current count = %d, want 1: %s", got, body)
	}
}

func TestOperationsPhase02UserFacingLabelsAndActiveStateStayRegistryDriven(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02ReadOnlyPrincipal()})
	adminHandler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02AdminPrincipal()})
	for _, route := range operationsCanonicalHTMLRoutes() {
		if route.PageTitle == "" || route.NavLabel == "" {
			t.Fatalf("route must keep user-facing title and nav label: %+v", route)
		}
		if operationsPageTitle(route.Section) != route.PageTitle {
			t.Fatalf("route %s page title helper = %q, want %q", route.Section, operationsPageTitle(route.Section), route.PageTitle)
		}
		phase02AssertNoUnsupportedRouteLabelClaim(t, route.Path, route.PageTitle)
		phase02AssertNoUnsupportedRouteLabelClaim(t, route.Path, route.NavLabel)

		t.Run(route.Section, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, route.Path, nil)
			rr := httptest.NewRecorder()
			phase02HandlerForRoute(route, handler, adminHandler).ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("%s status = %d, want 200: %s", route.Path, rr.Code, rr.Body.String())
			}
			body := rr.Body.String()
			escapedTitle := html.EscapeString(route.PageTitle)
			for _, want := range []string{
				`<title>` + escapedTitle + `</title>`,
				`<h1 id="operations-page-title">` + escapedTitle + `</h1>`,
			} {
				if !strings.Contains(body, want) {
					t.Fatalf("%s missing title marker %q: %s", route.Path, want, body)
				}
			}
			currentMarker := `href="` + route.Path + `" aria-current="page">` + html.EscapeString(route.NavLabel)
			if !strings.Contains(body, currentMarker) {
				t.Fatalf("%s missing current nav marker %q: %s", route.Path, currentMarker, body)
			}
			if got := strings.Count(body, `aria-current="page"`); got != 1 {
				t.Fatalf("%s aria-current count = %d, want 1: %s", route.Path, got, body)
			}
		})
	}

	dashboardReq := httptest.NewRequest(http.MethodGet, "/admin/operations", nil)
	dashboardRR := httptest.NewRecorder()
	handler.ServeHTTP(dashboardRR, dashboardReq)
	if dashboardRR.Code != http.StatusOK {
		t.Fatalf("dashboard status = %d, want 200: %s", dashboardRR.Code, dashboardRR.Body.String())
	}
	body := dashboardRR.Body.String()
	for _, staleLabel := range []string{">Dashboard</a>", ">Start Here</a>", ">Devices &amp; Tokens</a>", ">Telemetry Simulator</a>", ">Validation Center</a>", ">Setup Details</a>"} {
		if strings.Contains(body, staleLabel) {
			t.Fatalf("dashboard navigation still contains stale label %q: %s", staleLabel, body)
		}
	}
	for _, path := range []string{
		"/admin/operations",
		"/admin/operations/devices",
		"/admin/operations/telemetry-simulator",
		"/admin/operations/validation-center",
		"/admin/operations/setup-wizard",
	} {
		route, ok := operationsRouteByPath(path)
		if !ok {
			t.Fatalf("route registry missing %s", path)
		}
		want := `href="` + path + `"`
		if path == "/admin/operations" {
			want += ` aria-current="page"`
		}
		want += `>` + html.EscapeString(route.NavLabel)
		if !strings.Contains(body, want) {
			t.Fatalf("dashboard navigation missing expected user-facing label %q: %s", want, body)
		}
	}
}

func TestOperationsPhase136PrimaryPagesUseActionFirstProductLanguage(t *testing.T) {
	handler := phase02OperationsHandler(t, auth.TestAuthenticator{Principal: phase02ReadOnlyPrincipal()})
	corePaths := []string{
		"/admin/operations",
		"/admin/operations/setup-wizard",
		"/admin/operations/setup",
		"/admin/operations/gtfs-import",
		"/admin/operations/gtfs-workbench",
		"/admin/operations/gtfs-quality",
		"/admin/operations/feeds",
		"/admin/operations/feed-health",
		"/admin/operations/validation-center",
		"/admin/operations/validation-health",
		"/admin/operations/realtime",
		"/admin/operations/prediction-lab",
		"/admin/operations/devices",
		"/admin/operations/telemetry",
		"/admin/operations/telemetry-simulator",
		"/admin/operations/connectors",
		"/admin/operations/connectors/workbench",
		"/admin/operations/connectors/tests",
		"/admin/operations/readiness",
		"/admin/operations/maintenance",
		"/admin/operations/help",
	}
	bannedVisibleCopy := []string{
		"technical helper",
		"technical-helper",
		"common next action",
		"common next actions",
		"what this does not prove",
	}
	rawClaimFlags := []string{
		"external_evidence_created",
		"consumer_statuses_changed",
		"production_grade_eta_claimed",
		"hosted_saas_claimed",
		"dynamic_backend_plugin_loading_enabled",
	}

	for _, path := range corePaths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("%s status = %d, want 200: %s", path, rr.Code, rr.Body.String())
			}
			body := rr.Body.String()
			lower := strings.ToLower(body)
			for _, banned := range bannedVisibleCopy {
				if strings.Contains(lower, banned) {
					t.Fatalf("%s contains banned primary-page copy %q: %s", path, banned, body)
				}
			}
			if path != "/admin/operations" {
				actionIndex := strings.Index(body, `class="page-next-action"`)
				tableIndex := strings.Index(body, "<table")
				if actionIndex < 0 {
					t.Fatalf("%s is missing the shared page next action: %s", path, body)
				}
				if tableIndex >= 0 && actionIndex > tableIndex {
					t.Fatalf("%s shows diagnostics before the next action: action=%d table=%d body=%s", path, actionIndex, tableIndex, body)
				}
			}
			if flagIndex := firstStringIndex(lower, rawClaimFlags); flagIndex >= 0 {
				advancedIndex := strings.LastIndex(lower[:flagIndex], "advanced safety details")
				if advancedIndex < 0 {
					t.Fatalf("%s exposes raw safety flag names before advanced details: %s", path, body)
				}
			}
		})
	}
}

func firstStringIndex(s string, needles []string) int {
	first := -1
	for _, needle := range needles {
		index := strings.Index(s, needle)
		if index >= 0 && (first < 0 || index < first) {
			first = index
		}
	}
	return first
}

func phase02OperationsHandler(t testing.TB, admin adminAuth) http.Handler {
	t.Helper()
	return newOperationsTestHandler(&handler{store: feedHealthTestStore(t), devices: fakeDeviceStore{}}, admin)
}

func phase02ReadOnlyPrincipal() auth.Principal {
	return auth.Principal{
		Subject:  "reader@example.com",
		AgencyID: "demo-agency",
		Roles:    []auth.Role{auth.RoleReadOnly},
		Method:   auth.MethodBearer,
	}
}

func phase02AdminPrincipal() auth.Principal {
	return auth.Principal{
		Subject:  "admin@example.com",
		AgencyID: "demo-agency",
		Roles:    []auth.Role{auth.RoleAdmin},
		Method:   auth.MethodBearer,
	}
}

func phase02HandlerForRoute(route operationsRouteMeta, readOnly http.Handler, admin http.Handler) http.Handler {
	if route.Section == "admin-users" {
		return admin
	}
	return readOnly
}

func phase02RouteAllows(methods []string, method string) bool {
	for _, candidate := range methods {
		if candidate == method {
			return true
		}
	}
	return false
}

func phase02AssertNoUnsupportedRouteLabelClaim(t testing.TB, context string, value string) {
	t.Helper()
	lower := strings.ToLower(value)
	for _, forbidden := range []string{
		"agency approved",
		"agency adopted",
		"consumer accepted",
		"accepted by",
		"compliance achieved",
		"cal-itp/caltrans compliant",
		"production ready",
		"public launch complete",
		"launch complete",
		"vendor compatible",
		"certified hardware",
		"validated eta performance",
		"consumer-ready",
	} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("%s contains unsupported claim phrase %q in %q", context, forbidden, value)
		}
	}
}
