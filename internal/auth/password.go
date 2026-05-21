package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/tenant"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/argon2"
)

const (
	passwordArgonMemory      = 64 * 1024
	passwordArgonIterations  = 3
	passwordArgonParallelism = 1
	passwordSaltBytes        = 16
	passwordKeyBytes         = 32

	defaultBootstrapTokenTTL = 30 * time.Minute
	maxBootstrapTokenTTL     = 24 * time.Hour
	passwordLockoutAttempts  = 5
	passwordLockoutDuration  = 15 * time.Minute
)

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrCredentialLocked      = errors.New("credential locked")
	ErrBootstrapTokenInvalid = errors.New("invalid bootstrap token")
	ErrBootstrapTokenExpired = errors.New("expired bootstrap token")
)

type FirstAdminBootstrapInput struct {
	AgencyID    string
	Email       string
	DisplayName string
	AuthSubject string
	CreatedBy   string
	TTL         time.Duration
	Now         time.Time
}

type BootstrapResult struct {
	AgencyID  string
	Email     string
	Subject   string
	UserID    int64
	TokenID   int64
	ExpiresAt time.Time
}

type BootstrapTokenRecord struct {
	ID        int64
	AgencyID  string
	UserID    int64
	Email     string
	Subject   string
	Purpose   string
	ExpiresAt time.Time
}

type PasswordCredentialInput struct {
	AgencyID      string
	UserID        int64
	PlainPassword string
	Now           time.Time
}

type PostgresAdminStore struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

func NewPostgresAdminStore(pool *pgxpool.Pool) *PostgresAdminStore {
	return &PostgresAdminStore{pool: pool, now: time.Now}
}

func GenerateBootstrapToken() (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate bootstrap token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func HashBootstrapToken(token string) string {
	sum := sha256.Sum256([]byte("admin-bootstrap-token-v1\x00" + strings.TrimSpace(token)))
	return hex.EncodeToString(sum[:])
}

func ValidateAdminPassword(password string) error {
	if len(password) < 12 {
		return fmt.Errorf("password must be at least 12 characters")
	}
	if len(password) > 256 {
		return fmt.Errorf("password must be at most 256 characters")
	}
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("password must not be blank")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	if err := ValidateAdminPassword(password); err != nil {
		return "", err
	}
	salt := make([]byte, passwordSaltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}
	key := argon2.IDKey([]byte(password), salt, passwordArgonIterations, passwordArgonMemory, passwordArgonParallelism, passwordKeyBytes)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		passwordArgonMemory,
		passwordArgonIterations,
		passwordArgonParallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func VerifyPassword(encoded string, password string) (bool, error) {
	params, salt, expected, err := parsePasswordHash(encoded)
	if err != nil {
		return false, err
	}
	actual := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, uint32(len(expected)))
	return subtle.ConstantTimeCompare(actual, expected) == 1, nil
}

type passwordHashParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
}

func parsePasswordHash(encoded string) (passwordHashParams, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return passwordHashParams{}, nil, nil, fmt.Errorf("unsupported password hash format")
	}
	params := passwordHashParams{}
	for _, part := range strings.Split(parts[3], ",") {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			return passwordHashParams{}, nil, nil, fmt.Errorf("invalid password hash parameter")
		}
		parsed, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return passwordHashParams{}, nil, nil, fmt.Errorf("invalid password hash parameter: %w", err)
		}
		switch key {
		case "m":
			params.memory = uint32(parsed)
		case "t":
			params.iterations = uint32(parsed)
		case "p":
			params.parallelism = uint8(parsed)
		default:
			return passwordHashParams{}, nil, nil, fmt.Errorf("unknown password hash parameter")
		}
	}
	if params.memory == 0 || params.iterations == 0 || params.parallelism == 0 {
		return passwordHashParams{}, nil, nil, fmt.Errorf("missing password hash parameter")
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return passwordHashParams{}, nil, nil, fmt.Errorf("decode password salt: %w", err)
	}
	key, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return passwordHashParams{}, nil, nil, fmt.Errorf("decode password key: %w", err)
	}
	return params, salt, key, nil
}

func NormalizeAdminEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if len(email) > 320 || strings.ContainsAny(email, "\r\n\t ") {
		return "", fmt.Errorf("email must be a single address")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return "", fmt.Errorf("email must be valid")
	}
	return email, nil
}

func NormalizeBootstrapTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return defaultBootstrapTokenTTL
	}
	if ttl > maxBootstrapTokenTTL {
		return maxBootstrapTokenTTL
	}
	return ttl
}

func (s *PostgresAdminStore) CreateFirstAdminBootstrap(ctx context.Context, input FirstAdminBootstrapInput, tokenHash string) (BootstrapResult, error) {
	if s == nil || s.pool == nil {
		return BootstrapResult{}, fmt.Errorf("admin auth store is unavailable")
	}
	if err := tenant.ValidateAgencyID(strings.TrimSpace(input.AgencyID)); err != nil {
		return BootstrapResult{}, fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	email, err := NormalizeAdminEmail(input.Email)
	if err != nil {
		return BootstrapResult{}, err
	}
	if strings.TrimSpace(tokenHash) == "" || tokenHash == HashBootstrapToken("") {
		return BootstrapResult{}, fmt.Errorf("token hash is required")
	}
	now := input.Now.UTC()
	if now.IsZero() {
		now = s.now().UTC()
	}
	ttl := NormalizeBootstrapTTL(input.TTL)
	expiresAt := now.Add(ttl)
	subject := firstNonEmptyString(input.AuthSubject, email)
	createdBy := firstNonEmptyString(input.CreatedBy, "operator_console")

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("begin bootstrap transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID int64
	var storedSubject string
	err = tx.QueryRow(ctx, `
		INSERT INTO agency_user (agency_id, email, display_name, auth_subject)
		VALUES ($1, $2, NULLIF($3, ''), $4)
		ON CONFLICT (agency_id, email) DO UPDATE
		SET display_name = COALESCE(NULLIF(EXCLUDED.display_name, ''), agency_user.display_name),
		    auth_subject = COALESCE(NULLIF(agency_user.auth_subject, ''), EXCLUDED.auth_subject),
		    disabled_at = NULL
		RETURNING id, COALESCE(auth_subject, email)
	`, input.AgencyID, email, strings.TrimSpace(input.DisplayName), subject).Scan(&userID, &storedSubject)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("upsert agency user: %w", err)
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO role_binding (agency_id, agency_user_id, role)
		VALUES ($1, $2, 'admin')
		ON CONFLICT (agency_id, agency_user_id, role) DO NOTHING
	`, input.AgencyID, userID)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("bind admin role: %w", err)
	}
	var tokenID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO admin_bootstrap_token (agency_id, agency_user_id, purpose, token_hash, expires_at, created_by, created_at)
		VALUES ($1, $2, 'first_admin', $3, $4, $5, $6)
		RETURNING id
	`, input.AgencyID, userID, tokenHash, expiresAt, createdBy, now).Scan(&tokenID)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("insert bootstrap token: %w", err)
	}
	_, _ = tx.Exec(ctx, `
		INSERT INTO audit_log (agency_id, actor_id, action, entity_type, entity_id, new_value_json, created_at)
		VALUES ($1, $2, 'admin.bootstrap_token.create', 'agency_user', $3, jsonb_build_object('purpose', 'first_admin', 'expires_at', $4::text), $5)
	`, input.AgencyID, createdBy, strconv.FormatInt(userID, 10), expiresAt.Format(time.RFC3339), now)
	if err := tx.Commit(ctx); err != nil {
		return BootstrapResult{}, fmt.Errorf("commit bootstrap transaction: %w", err)
	}
	return BootstrapResult{AgencyID: input.AgencyID, Email: email, Subject: storedSubject, UserID: userID, TokenID: tokenID, ExpiresAt: expiresAt}, nil
}

func (s *PostgresAdminStore) ConsumeBootstrapToken(ctx context.Context, token string, now time.Time) (BootstrapTokenRecord, error) {
	if s == nil || s.pool == nil {
		return BootstrapTokenRecord{}, fmt.Errorf("admin auth store is unavailable")
	}
	tokenHash := HashBootstrapToken(token)
	if tokenHash == HashBootstrapToken("") {
		return BootstrapTokenRecord{}, ErrBootstrapTokenInvalid
	}
	if now.IsZero() {
		now = s.now().UTC()
	} else {
		now = now.UTC()
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return BootstrapTokenRecord{}, fmt.Errorf("begin token transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	var record BootstrapTokenRecord
	var usedAt *time.Time
	err = tx.QueryRow(ctx, `
		SELECT bt.id, bt.agency_id, bt.agency_user_id, au.email, COALESCE(au.auth_subject, au.email), bt.purpose, bt.expires_at, bt.used_at
		FROM admin_bootstrap_token bt
		JOIN agency_user au
		  ON au.agency_id = bt.agency_id
		 AND au.id = bt.agency_user_id
		WHERE bt.token_hash = $1
		FOR UPDATE
	`, tokenHash).Scan(&record.ID, &record.AgencyID, &record.UserID, &record.Email, &record.Subject, &record.Purpose, &record.ExpiresAt, &usedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return BootstrapTokenRecord{}, ErrBootstrapTokenInvalid
	}
	if err != nil {
		return BootstrapTokenRecord{}, fmt.Errorf("load bootstrap token: %w", err)
	}
	if usedAt != nil {
		return BootstrapTokenRecord{}, ErrBootstrapTokenInvalid
	}
	if now.After(record.ExpiresAt) {
		return BootstrapTokenRecord{}, ErrBootstrapTokenExpired
	}
	if _, err := tx.Exec(ctx, `UPDATE admin_bootstrap_token SET used_at = $1 WHERE id = $2`, now, record.ID); err != nil {
		return BootstrapTokenRecord{}, fmt.Errorf("mark bootstrap token used: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return BootstrapTokenRecord{}, fmt.Errorf("commit token transaction: %w", err)
	}
	return record, nil
}

func (s *PostgresAdminStore) SetPasswordCredential(ctx context.Context, input PasswordCredentialInput) error {
	if s == nil || s.pool == nil {
		return fmt.Errorf("admin auth store is unavailable")
	}
	if err := tenant.ValidateAgencyID(strings.TrimSpace(input.AgencyID)); err != nil {
		return fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	if input.UserID <= 0 {
		return fmt.Errorf("agency user id is required")
	}
	hash, err := HashPassword(input.PlainPassword)
	if err != nil {
		return err
	}
	now := input.Now.UTC()
	if now.IsZero() {
		now = s.now().UTC()
	}
	tag, err := s.pool.Exec(ctx, `
		INSERT INTO admin_password_credential (
			agency_id, agency_user_id, password_hash, password_hash_params, status,
			failed_attempts, locked_until, password_changed_at, updated_at
		)
		VALUES ($1, $2, $3, '{"algorithm":"argon2id","version":1}'::jsonb, 'active', 0, NULL, $4, $4)
		ON CONFLICT (agency_id, agency_user_id) DO UPDATE
		SET password_hash = EXCLUDED.password_hash,
		    password_hash_params = EXCLUDED.password_hash_params,
		    status = 'active',
		    failed_attempts = 0,
		    locked_until = NULL,
		    password_changed_at = EXCLUDED.password_changed_at,
		    updated_at = EXCLUDED.updated_at
	`, input.AgencyID, input.UserID, hash, now)
	if err != nil {
		return fmt.Errorf("store password credential: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("password credential was not stored")
	}
	return nil
}

func (s *PostgresAdminStore) CompleteBootstrapPassword(ctx context.Context, token string, password string, now time.Time) (Principal, error) {
	if s == nil || s.pool == nil {
		return Principal{}, fmt.Errorf("admin auth store is unavailable")
	}
	hash, err := HashPassword(password)
	if err != nil {
		return Principal{}, err
	}
	tokenHash := HashBootstrapToken(token)
	if tokenHash == HashBootstrapToken("") {
		return Principal{}, ErrBootstrapTokenInvalid
	}
	if now.IsZero() {
		now = s.now().UTC()
	} else {
		now = now.UTC()
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Principal{}, fmt.Errorf("begin first admin setup transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var record BootstrapTokenRecord
	var usedAt *time.Time
	err = tx.QueryRow(ctx, `
		SELECT bt.id, bt.agency_id, bt.agency_user_id, au.email, COALESCE(au.auth_subject, au.email), bt.purpose, bt.expires_at, bt.used_at
		FROM admin_bootstrap_token bt
		JOIN agency_user au
		  ON au.agency_id = bt.agency_id
		 AND au.id = bt.agency_user_id
		WHERE bt.token_hash = $1
		FOR UPDATE
	`, tokenHash).Scan(&record.ID, &record.AgencyID, &record.UserID, &record.Email, &record.Subject, &record.Purpose, &record.ExpiresAt, &usedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Principal{}, ErrBootstrapTokenInvalid
	}
	if err != nil {
		return Principal{}, fmt.Errorf("load bootstrap token: %w", err)
	}
	if usedAt != nil {
		return Principal{}, ErrBootstrapTokenInvalid
	}
	if now.After(record.ExpiresAt) {
		return Principal{}, ErrBootstrapTokenExpired
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO admin_password_credential (
			agency_id, agency_user_id, password_hash, password_hash_params, status,
			failed_attempts, locked_until, password_changed_at, updated_at
		)
		VALUES ($1, $2, $3, '{"algorithm":"argon2id","version":1}'::jsonb, 'active', 0, NULL, $4, $4)
		ON CONFLICT (agency_id, agency_user_id) DO UPDATE
		SET password_hash = EXCLUDED.password_hash,
		    password_hash_params = EXCLUDED.password_hash_params,
		    status = 'active',
		    failed_attempts = 0,
		    locked_until = NULL,
		    password_changed_at = EXCLUDED.password_changed_at,
		    updated_at = EXCLUDED.updated_at
	`, record.AgencyID, record.UserID, hash, now); err != nil {
		return Principal{}, fmt.Errorf("store password credential: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE admin_bootstrap_token SET used_at = $1 WHERE id = $2`, now, record.ID); err != nil {
		return Principal{}, fmt.Errorf("mark bootstrap token used: %w", err)
	}
	rows, err := tx.Query(ctx, `
		SELECT role
		FROM role_binding
		WHERE agency_id = $1
		  AND agency_user_id = $2
		ORDER BY role
	`, record.AgencyID, record.UserID)
	if err != nil {
		return Principal{}, fmt.Errorf("query role bindings: %w", err)
	}
	defer rows.Close()
	roles, err := RowsToRoles(rows)
	if err != nil {
		return Principal{}, err
	}
	if len(roles) == 0 {
		return Principal{}, ErrInvalidCredentials
	}
	if err := tx.Commit(ctx); err != nil {
		return Principal{}, fmt.Errorf("commit first admin setup transaction: %w", err)
	}
	return Principal{Subject: record.Subject, AgencyID: record.AgencyID, Roles: roles, Method: MethodCookie}, nil
}

func (s *PostgresAdminStore) AuthenticatePassword(ctx context.Context, agencyID string, email string, password string, now time.Time) (Principal, error) {
	if s == nil || s.pool == nil {
		return Principal{}, fmt.Errorf("admin auth store is unavailable")
	}
	if err := tenant.ValidateAgencyID(strings.TrimSpace(agencyID)); err != nil {
		return Principal{}, fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	normalizedEmail, err := NormalizeAdminEmail(email)
	if err != nil {
		return Principal{}, ErrInvalidCredentials
	}
	if now.IsZero() {
		now = s.now().UTC()
	} else {
		now = now.UTC()
	}
	var userID int64
	var subject string
	var hash string
	var status string
	var failedAttempts int
	var lockedUntil *time.Time
	err = s.pool.QueryRow(ctx, `
		SELECT au.id, COALESCE(au.auth_subject, au.email), pc.password_hash, pc.status, pc.failed_attempts, pc.locked_until
		FROM agency_user au
		JOIN admin_password_credential pc
		  ON pc.agency_id = au.agency_id
		 AND pc.agency_user_id = au.id
		WHERE au.agency_id = $1
		  AND lower(au.email) = $2
		  AND au.disabled_at IS NULL
	`, agencyID, normalizedEmail).Scan(&userID, &subject, &hash, &status, &failedAttempts, &lockedUntil)
	if errors.Is(err, pgx.ErrNoRows) {
		return Principal{}, ErrInvalidCredentials
	}
	if err != nil {
		return Principal{}, fmt.Errorf("load password credential: %w", err)
	}
	if status == "disabled" {
		return Principal{}, ErrInvalidCredentials
	}
	if status == "locked" || lockedUntil != nil && now.Before(lockedUntil.UTC()) {
		return Principal{}, ErrCredentialLocked
	}
	ok, err := VerifyPassword(hash, password)
	if err != nil || !ok {
		nextAttempts := failedAttempts + 1
		var nextLockedUntil *time.Time
		nextStatus := status
		if nextAttempts >= passwordLockoutAttempts {
			lock := now.Add(passwordLockoutDuration)
			nextLockedUntil = &lock
			nextStatus = "locked"
		}
		_, _ = s.pool.Exec(ctx, `
			UPDATE admin_password_credential
			SET failed_attempts = $1, locked_until = $2, status = $3, updated_at = $4
			WHERE agency_id = $5 AND agency_user_id = $6
		`, nextAttempts, nextLockedUntil, nextStatus, now, agencyID, userID)
		return Principal{}, ErrInvalidCredentials
	}
	roles, err := NewPostgresRoleStore(s.pool).RolesForSubject(ctx, agencyID, subject)
	if err != nil {
		return Principal{}, err
	}
	if len(roles) == 0 {
		return Principal{}, ErrInvalidCredentials
	}
	_, _ = s.pool.Exec(ctx, `
		UPDATE admin_password_credential
		SET failed_attempts = 0, locked_until = NULL, status = 'active', last_login_at = $1, updated_at = $1
		WHERE agency_id = $2 AND agency_user_id = $3
	`, now, agencyID, userID)
	return Principal{Subject: subject, AgencyID: agencyID, Roles: roles, Method: MethodCookie}, nil
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
