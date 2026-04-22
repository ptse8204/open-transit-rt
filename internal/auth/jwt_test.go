package auth

import (
	"testing"
	"time"
)

func TestJWTRequiresCoreClaimsAndValidatesAudienceIssuer(t *testing.T) {
	cfg := JWTConfig{Secrets: []string{"test-secret"}, Issuer: "open-transit-rt", Audience: "admin-api", ClockSkew: time.Minute, TTL: time.Hour}
	signer, err := NewSigner(cfg)
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}
	verifier, err := NewVerifier(cfg)
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}
	token, claims, err := signer.Sign("operator@example.com", "demo-agency", time.Hour)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if claims.JTI == "" {
		t.Fatalf("jti is empty")
	}
	verified, err := verifier.Verify(token)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if verified.Subject != "operator@example.com" || verified.AgencyID != "demo-agency" || verified.Issuer != cfg.Issuer || verified.Audience != cfg.Audience {
		t.Fatalf("claims = %+v, want scoped admin claims", verified)
	}

	wrongAudience, _ := NewVerifier(JWTConfig{Secrets: []string{"test-secret"}, Issuer: "open-transit-rt", Audience: "other", ClockSkew: time.Minute})
	if _, err := wrongAudience.Verify(token); err == nil {
		t.Fatalf("verify succeeded with wrong audience")
	}
}

func TestJWTAcceptsOldSecretForRotation(t *testing.T) {
	oldCfg := JWTConfig{Secrets: []string{"old-secret"}, Issuer: "open-transit-rt", Audience: "admin-api"}
	signer, err := NewSigner(oldCfg)
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}
	token, _, err := signer.Sign("operator@example.com", "demo-agency", time.Hour)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	verifier, err := NewVerifier(JWTConfig{Secrets: []string{"new-secret", "old-secret"}, Issuer: "open-transit-rt", Audience: "admin-api"})
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}
	if _, err := verifier.Verify(token); err != nil {
		t.Fatalf("verify with old secret during rotation: %v", err)
	}
}

func TestCSRFTokenBindsToPrincipal(t *testing.T) {
	principal := Principal{Subject: "operator@example.com", AgencyID: "demo-agency"}
	token := CSRFToken("csrf-secret", principal)
	if token == "" {
		t.Fatalf("csrf token is empty")
	}
	other := CSRFToken("csrf-secret", Principal{Subject: principal.Subject, AgencyID: "other-agency"})
	if token == other {
		t.Fatalf("csrf token did not bind to agency")
	}
}
