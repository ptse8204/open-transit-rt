package devices

import "testing"

func TestHashTokenUsesPepper(t *testing.T) {
	left := NewPostgresStore(nil, Config{TokenPepper: "pepper-a"}).HashToken("device-token")
	right := NewPostgresStore(nil, Config{TokenPepper: "pepper-b"}).HashToken("device-token")
	if left == right {
		t.Fatalf("hash did not change when pepper changed")
	}
	if len(left) <= len("hmac-sha256:") || left[:12] != "hmac-sha256:" {
		t.Fatalf("hash = %q, want hmac-sha256 prefix", left)
	}
}

func TestGenerateTokenReturnsOpaqueValue(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	if len(token) < 32 {
		t.Fatalf("token length = %d, want opaque random token", len(token))
	}
}
