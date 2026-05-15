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
