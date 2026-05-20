package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"open-transit-rt/internal/auth"
)

func TestBootstrapAdminLinkConfigRejectsUnsafeAgencyIDBeforeDB(t *testing.T) {
	env := func(key string) string {
		switch key {
		case "AGENCY_ID":
			return "demo-agency"
		case "PUBLIC_BASE_URL":
			return "https://feeds.example.org"
		default:
			return ""
		}
	}
	for _, agencyID := range []string{"agency/bad", `agency\bad`, ".hidden", "agency bad"} {
		_, err := parseBootstrapAdminLinkConfig([]string{"--agency-id", agencyID, "--email", "admin@example.org"}, env)
		if err == nil {
			t.Fatalf("unsafe agency id %q succeeded, want error", agencyID)
		}
	}
}

func TestBootstrapAdminLinkConfigNormalizesTTLAndBaseURL(t *testing.T) {
	env := func(key string) string {
		switch key {
		case "AGENCY_ID":
			return "demo-agency"
		case "ADMIN_BASE_URL":
			return "https://admin.example.org/control"
		default:
			return ""
		}
	}
	cfg, err := parseBootstrapAdminLinkConfig([]string{"--email", "Admin@Example.ORG", "--ttl", "48h"}, env)
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}
	if cfg.Email != "admin@example.org" || cfg.TTL != 24*time.Hour {
		t.Fatalf("config = %+v, want normalized email and capped ttl", cfg)
	}
	link, err := firstAdminSetupLink(cfg.BaseURL, "token-value")
	if err != nil {
		t.Fatalf("setup link: %v", err)
	}
	if link != "https://admin.example.org/control/admin/setup/first-admin?token=token-value" {
		t.Fatalf("link = %q", link)
	}
}

func TestBootstrapAdminLinkOutputShowsTokenOnceAndNotHash(t *testing.T) {
	token := "one-time-token-value"
	link, err := firstAdminSetupLink("https://admin.example.org", token)
	if err != nil {
		t.Fatalf("setup link: %v", err)
	}
	var out bytes.Buffer
	writeBootstrapAdminLink(&out, auth.BootstrapResult{
		AgencyID:  "demo-agency",
		Email:     "admin@example.org",
		Subject:   "admin@example.org",
		UserID:    10,
		TokenID:   20,
		ExpiresAt: time.Date(2026, 5, 20, 12, 30, 0, 0, time.UTC),
	}, link)
	body := out.String()
	if strings.Count(body, token) != 1 {
		t.Fatalf("token count = %d, want one output occurrence: %s", strings.Count(body, token), body)
	}
	if strings.Contains(body, auth.HashBootstrapToken(token)) {
		t.Fatalf("output leaked token hash: %s", body)
	}
	for _, forbidden := range []string{"password", "database_url", "token_hash", "Bearer "} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(forbidden)) {
			t.Fatalf("output leaked forbidden text %q: %s", forbidden, body)
		}
	}
}
