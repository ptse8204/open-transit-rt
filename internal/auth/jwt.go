package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	DefaultAdminTokenTTL = 8 * time.Hour
	DefaultClockSkew     = 2 * time.Minute
)

var (
	ErrInvalidToken = errors.New("invalid admin token")
	ErrExpiredToken = errors.New("expired admin token")
)

type JWTConfig struct {
	Secrets   []string
	Issuer    string
	Audience  string
	ClockSkew time.Duration
	TTL       time.Duration
}

func (c JWTConfig) validated() (JWTConfig, error) {
	if len(c.Secrets) == 0 || strings.TrimSpace(c.Secrets[0]) == "" {
		return JWTConfig{}, fmt.Errorf("ADMIN_JWT_SECRET is required")
	}
	if c.Issuer == "" {
		return JWTConfig{}, fmt.Errorf("ADMIN_JWT_ISSUER is required")
	}
	if c.Audience == "" {
		return JWTConfig{}, fmt.Errorf("ADMIN_JWT_AUDIENCE is required")
	}
	if c.ClockSkew <= 0 {
		c.ClockSkew = DefaultClockSkew
	}
	if c.TTL <= 0 {
		c.TTL = DefaultAdminTokenTTL
	}
	return c, nil
}

type Claims struct {
	Subject  string `json:"sub"`
	AgencyID string `json:"agency_id"`
	IssuedAt int64  `json:"iat"`
	Expires  int64  `json:"exp"`
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	JTI      string `json:"jti,omitempty"`
	Email    string `json:"email,omitempty"`
}

type Verifier struct {
	config JWTConfig
	now    func() time.Time
}

func NewVerifier(config JWTConfig) (*Verifier, error) {
	validated, err := config.validated()
	if err != nil {
		return nil, err
	}
	return &Verifier{config: validated, now: time.Now}, nil
}

func (v *Verifier) Verify(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}
	headerBytes, err := decodeSegment(parts[0])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return Claims{}, ErrInvalidToken
	}
	if header.Alg != "HS256" || header.Typ != "JWT" {
		return Claims{}, ErrInvalidToken
	}
	signingInput := parts[0] + "." + parts[1]
	signature, err := decodeSegment(parts[2])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	if !v.validSignature(signingInput, signature) {
		return Claims{}, ErrInvalidToken
	}
	payload, err := decodeSegment(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}
	if err := v.validateClaims(claims); err != nil {
		return Claims{}, err
	}
	return claims, nil
}

func (v *Verifier) validSignature(signingInput string, signature []byte) bool {
	for _, secret := range v.config.Secrets {
		secret = strings.TrimSpace(secret)
		if secret == "" {
			continue
		}
		mac := hmac.New(sha256.New, []byte(secret))
		_, _ = mac.Write([]byte(signingInput))
		if hmac.Equal(signature, mac.Sum(nil)) {
			return true
		}
	}
	return false
}

func (v *Verifier) validateClaims(claims Claims) error {
	if claims.Subject == "" || claims.AgencyID == "" || claims.IssuedAt == 0 || claims.Expires == 0 || claims.Issuer == "" || claims.Audience == "" {
		return ErrInvalidToken
	}
	if claims.Issuer != v.config.Issuer || claims.Audience != v.config.Audience {
		return ErrInvalidToken
	}
	now := v.now().UTC()
	skew := v.config.ClockSkew
	if now.After(time.Unix(claims.Expires, 0).Add(skew)) {
		return ErrExpiredToken
	}
	if time.Unix(claims.IssuedAt, 0).After(now.Add(skew)) {
		return ErrInvalidToken
	}
	return nil
}

type Signer struct {
	config JWTConfig
	now    func() time.Time
}

func NewSigner(config JWTConfig) (*Signer, error) {
	validated, err := config.validated()
	if err != nil {
		return nil, err
	}
	return &Signer{config: validated, now: time.Now}, nil
}

func (s *Signer) Sign(subject string, agencyID string, ttl time.Duration) (string, Claims, error) {
	if subject == "" || agencyID == "" {
		return "", Claims{}, fmt.Errorf("subject and agency_id are required")
	}
	if ttl <= 0 {
		ttl = s.config.TTL
	}
	now := s.now().UTC()
	jti, err := randomID()
	if err != nil {
		return "", Claims{}, err
	}
	claims := Claims{
		Subject:  subject,
		AgencyID: agencyID,
		IssuedAt: now.Unix(),
		Expires:  now.Add(ttl).Unix(),
		Issuer:   s.config.Issuer,
		Audience: s.config.Audience,
		JTI:      jti,
	}
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(claims)
	signingInput := encodeSegment(headerJSON) + "." + encodeSegment(claimsJSON)
	mac := hmac.New(sha256.New, []byte(s.config.Secrets[0]))
	_, _ = mac.Write([]byte(signingInput))
	return signingInput + "." + encodeSegment(mac.Sum(nil)), claims, nil
}

func randomID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate jti: %w", err)
	}
	return encodeSegment(raw[:]), nil
}

func encodeSegment(payload []byte) string {
	return base64.RawURLEncoding.EncodeToString(payload)
}

func decodeSegment(segment string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(segment)
}
