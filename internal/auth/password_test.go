package auth

import (
	"strings"
	"testing"
	"time"
)

func TestPasswordHashUsesArgon2IDAndDoesNotStorePlaintext(t *testing.T) {
	password := "correct horse battery staple"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !strings.HasPrefix(hash, "$argon2id$v=19$") {
		t.Fatalf("hash = %q, want argon2id envelope", hash)
	}
	if strings.Contains(hash, password) || strings.Contains(hash, "correct") {
		t.Fatalf("hash leaked password text: %q", hash)
	}
	ok, err := VerifyPassword(hash, password)
	if err != nil || !ok {
		t.Fatalf("VerifyPassword correct = %v, %v; want true nil", ok, err)
	}
	ok, err = VerifyPassword(hash, "wrong password value")
	if err != nil || ok {
		t.Fatalf("VerifyPassword wrong = %v, %v; want false nil", ok, err)
	}
}

func TestPasswordPolicyRejectsUnsafeValues(t *testing.T) {
	for _, password := range []string{"short", "            ", strings.Repeat("x", 257)} {
		if err := ValidateAdminPassword(password); err == nil {
			t.Fatalf("ValidateAdminPassword(%q) succeeded, want error", password)
		}
	}
	if err := ValidateAdminPassword("long enough password"); err != nil {
		t.Fatalf("valid password rejected: %v", err)
	}
}

func TestBootstrapTokenHashIsHashOnlyAndStable(t *testing.T) {
	token := "setup-token-value"
	hash := HashBootstrapToken(token)
	if hash == "" || hash == token || strings.Contains(hash, token) {
		t.Fatalf("hash leaked token: token=%q hash=%q", token, hash)
	}
	if got := HashBootstrapToken(token); got != hash {
		t.Fatalf("HashBootstrapToken not stable: %q vs %q", got, hash)
	}
	if got := HashBootstrapToken(token + "-other"); got == hash {
		t.Fatalf("different tokens produced same hash %q", hash)
	}
}

func TestNormalizeBootstrapTTLCapsAndDefaults(t *testing.T) {
	if got := NormalizeBootstrapTTL(0); got != defaultBootstrapTokenTTL {
		t.Fatalf("zero ttl = %s, want %s", got, defaultBootstrapTokenTTL)
	}
	if got := NormalizeBootstrapTTL(48 * time.Hour); got != maxBootstrapTokenTTL {
		t.Fatalf("large ttl = %s, want %s", got, maxBootstrapTokenTTL)
	}
	if got := NormalizeBootstrapTTL(10 * time.Minute); got != 10*time.Minute {
		t.Fatalf("normal ttl = %s, want 10m", got)
	}
}

func TestNormalizeAdminEmail(t *testing.T) {
	email, err := NormalizeAdminEmail(" Admin@Example.ORG ")
	if err != nil {
		t.Fatalf("NormalizeAdminEmail valid: %v", err)
	}
	if email != "admin@example.org" {
		t.Fatalf("email = %q, want normalized lowercase", email)
	}
	for _, raw := range []string{"", "not-an-email", "admin @example.org", "admin@example.org\nbcc@example.org"} {
		if _, err := NormalizeAdminEmail(raw); err == nil {
			t.Fatalf("NormalizeAdminEmail(%q) succeeded, want error", raw)
		}
	}
}
