package devices

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"open-transit-rt/internal/appconfig"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	TokenPepper string
}

func ConfigFromEnv() (Config, error) {
	pepper := strings.TrimSpace(os.Getenv("DEVICE_TOKEN_PEPPER"))
	if pepper == "" {
		return Config{}, fmt.Errorf("DEVICE_TOKEN_PEPPER is required")
	}
	if appconfig.IsProduction() && appconfig.IsDevPlaceholder(pepper) {
		return Config{}, fmt.Errorf("DEVICE_TOKEN_PEPPER must not use a development placeholder in production")
	}
	return Config{TokenPepper: pepper}, nil
}

type Credential struct {
	AgencyID  string
	DeviceID  string
	VehicleID string
	Status    string
}

type VerifyInput struct {
	Token     string
	AgencyID  string
	DeviceID  string
	VehicleID string
}

type RebindInput struct {
	AgencyID  string
	DeviceID  string
	VehicleID string
	ActorID   string
	Reason    string
	Now       time.Time
}

type RebindResult struct {
	AgencyID  string `json:"agency_id"`
	DeviceID  string `json:"device_id"`
	VehicleID string `json:"vehicle_id"`
	Token     string `json:"token"`
	RotatedAt string `json:"rotated_at"`
}

type Store interface {
	Verify(ctx context.Context, input VerifyInput) (Credential, error)
	Rebind(ctx context.Context, input RebindInput) (RebindResult, error)
}

type PostgresStore struct {
	pool   *pgxpool.Pool
	config Config
}

func NewPostgresStore(pool *pgxpool.Pool, config Config) *PostgresStore {
	return &PostgresStore{pool: pool, config: config}
}

func (s *PostgresStore) Verify(ctx context.Context, input VerifyInput) (Credential, error) {
	if input.Token == "" || input.AgencyID == "" || input.DeviceID == "" || input.VehicleID == "" {
		return Credential{}, fmt.Errorf("missing device credential input")
	}
	var cred Credential
	err := s.pool.QueryRow(ctx, `
		SELECT agency_id, device_id, vehicle_id, status
		FROM device_credential
		WHERE agency_id = $1
		  AND device_id = $2
		  AND token_hash = $3
		  AND status = 'active'
		  AND revoked_at IS NULL
		  AND rotated_at IS NULL
	`, input.AgencyID, input.DeviceID, s.HashToken(input.Token)).
		Scan(&cred.AgencyID, &cred.DeviceID, &cred.VehicleID, &cred.Status)
	if err != nil {
		return Credential{}, fmt.Errorf("invalid device credential")
	}
	if cred.VehicleID != input.VehicleID {
		return Credential{}, fmt.Errorf("device is not bound to vehicle")
	}
	_, _ = s.pool.Exec(ctx, `
		UPDATE device_credential
		SET last_used_at = now()
		WHERE agency_id = $1 AND device_id = $2
	`, input.AgencyID, input.DeviceID)
	return cred, nil
}

func (s *PostgresStore) Rebind(ctx context.Context, input RebindInput) (RebindResult, error) {
	if input.AgencyID == "" || input.DeviceID == "" || input.VehicleID == "" || input.ActorID == "" {
		return RebindResult{}, fmt.Errorf("agency_id, device_id, vehicle_id, and actor_id are required")
	}
	if input.Now.IsZero() {
		input.Now = time.Now().UTC()
	}
	token, err := GenerateToken()
	if err != nil {
		return RebindResult{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return RebindResult{}, fmt.Errorf("begin device rebind: %w", err)
	}
	defer tx.Rollback(ctx)

	var oldVehicleID, oldStatus string
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(vehicle_id, ''), status
		FROM device_credential
		WHERE agency_id = $1 AND device_id = $2
		FOR UPDATE
	`, input.AgencyID, input.DeviceID).Scan(&oldVehicleID, &oldStatus)
	if err != nil && err != pgx.ErrNoRows {
		return RebindResult{}, fmt.Errorf("query device credential: %w", err)
	}
	if err == pgx.ErrNoRows {
		_, err = tx.Exec(ctx, `
			INSERT INTO device_credential (agency_id, device_id, vehicle_id, token_hash, status, valid_from, created_at)
			VALUES ($1, $2, $3, $4, 'active', $5, $5)
		`, input.AgencyID, input.DeviceID, input.VehicleID, s.HashToken(token), input.Now)
	} else {
		_, err = tx.Exec(ctx, `
			UPDATE device_credential
			SET vehicle_id = $3,
			    token_hash = $4,
			    status = 'active',
			    valid_from = $5,
			    rotated_at = NULL,
			    revoked_at = NULL
			WHERE agency_id = $1 AND device_id = $2
		`, input.AgencyID, input.DeviceID, input.VehicleID, s.HashToken(token), input.Now)
	}
	if err != nil {
		return RebindResult{}, fmt.Errorf("store rebound device credential: %w", err)
	}
	oldPayload, _ := json.Marshal(map[string]any{"device_id": input.DeviceID, "vehicle_id": oldVehicleID, "status": oldStatus})
	newPayload, _ := json.Marshal(map[string]any{"device_id": input.DeviceID, "vehicle_id": input.VehicleID, "status": "active", "token": "redacted"})
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (agency_id, actor_id, action, entity_type, entity_id, old_value_json, new_value_json, reason)
		VALUES ($1, $2, 'device.rebind', 'device_credential', $3, $4::jsonb, $5::jsonb, $6)
	`, input.AgencyID, input.ActorID, input.DeviceID, string(oldPayload), string(newPayload), nullIfEmpty(input.Reason)); err != nil {
		return RebindResult{}, fmt.Errorf("insert device rebind audit log: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return RebindResult{}, fmt.Errorf("commit device rebind: %w", err)
	}
	return RebindResult{
		AgencyID:  input.AgencyID,
		DeviceID:  input.DeviceID,
		VehicleID: input.VehicleID,
		Token:     token,
		RotatedAt: input.Now.Format(time.RFC3339),
	}, nil
}

func (s *PostgresStore) HashToken(token string) string {
	mac := hmac.New(sha256.New, []byte(s.config.TokenPepper))
	_, _ = mac.Write([]byte(token))
	return "hmac-sha256:" + hex.EncodeToString(mac.Sum(nil))
}

func GenerateToken() (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate device token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
