package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteAccessDeniedHTMLIsNoStoreAndEscapesReason(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin/operations/audit?agency_id=%3Cscript%3E", nil)
	req.Header.Set("Accept", "text/html")
	rr := httptest.NewRecorder()

	writeAccessDenied(rr, req, `bad <script>alert("x")</script> & agency`)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", got)
	}
	body := rr.Body.String()
	for _, want := range []string{"Access denied", "&lt;script&gt;", "requested agency scope matches your signed-in agency"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{`<script>alert("x")</script>`, "bad <script>", "other agency data exists"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("body leaked raw text %q: %s", forbidden, body)
		}
	}
}

func TestWriteAccessDeniedNonHTMLDoesNotExposeReason(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/admin/api/device-rebind", nil)
	req.Header.Set("Accept", "application/json")
	rr := httptest.NewRecorder()

	writeAccessDenied(rr, req, "requested agency scope does not match other-agency")

	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	body := rr.Body.String()
	if strings.TrimSpace(body) != "forbidden" {
		t.Fatalf("body = %q, want bounded forbidden response", body)
	}
	if strings.Contains(body, "other-agency") || strings.Contains(body, "requested agency") {
		t.Fatalf("non-HTML body leaked denial reason: %s", body)
	}
}

func TestRequireAgencyQueryMatchRejectsUnsafeAgencyID(t *testing.T) {
	principal := Principal{Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []Role{RoleReadOnly}, Method: MethodBearer}
	for _, target := range []string{
		"/admin/operations?agency_id=demo-agency",
		"/admin/operations",
	} {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		rr := httptest.NewRecorder()
		if !RequireAgencyQueryMatch(rr, req, principal) {
			t.Fatalf("RequireAgencyQueryMatch(%q) rejected valid scope: %d %s", target, rr.Code, rr.Body.String())
		}
	}

	for _, target := range []string{
		"/admin/operations?agency_id=demo-agency%2Fbad",
		"/admin/operations?agency_id=demo-agency%5Cbad",
		"/admin/operations?agency_id=.hidden",
		"/admin/operations?agency_id=agency.bad",
		"/admin/operations?agency_id=other-agency",
	} {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		req.Header.Set("Accept", "application/json")
		rr := httptest.NewRecorder()
		if RequireAgencyQueryMatch(rr, req, principal) {
			t.Fatalf("RequireAgencyQueryMatch(%q) succeeded, want rejection", target)
		}
		if rr.Code != http.StatusForbidden || rr.Header().Get("Cache-Control") != "no-store" {
			t.Fatalf("unsafe/conflicting scope response = %d cache=%q body=%s", rr.Code, rr.Header().Get("Cache-Control"), rr.Body.String())
		}
		forbidden := []string{"other-agency", "demo-agency/bad", `demo-agency\bad`, ".hidden", "agency.bad"}
		for _, value := range forbidden {
			if strings.Contains(rr.Body.String(), value) {
				t.Fatalf("forbidden response leaked agency value %q: %s", value, rr.Body.String())
			}
		}
	}
}

func TestRejectAgencyConflictRejectsUnsafeAgencyID(t *testing.T) {
	principal := Principal{Subject: "operator@example.com", AgencyID: "demo-agency", Roles: []Role{RoleAdmin}, Method: MethodBearer}
	for _, agencyID := range []string{"", "demo-agency"} {
		rr := httptest.NewRecorder()
		if RejectAgencyConflict(rr, agencyID, principal) {
			t.Fatalf("RejectAgencyConflict(%q) rejected valid scope", agencyID)
		}
	}
	for _, agencyID := range []string{"other-agency", "demo-agency/bad", `demo-agency\bad`, ".hidden"} {
		rr := httptest.NewRecorder()
		if !RejectAgencyConflict(rr, agencyID, principal) {
			t.Fatalf("RejectAgencyConflict(%q) succeeded, want rejection", agencyID)
		}
		if rr.Code != http.StatusForbidden || rr.Header().Get("Cache-Control") != "no-store" {
			t.Fatalf("unsafe/conflicting form scope response = %d cache=%q body=%s", rr.Code, rr.Header().Get("Cache-Control"), rr.Body.String())
		}
	}
}
