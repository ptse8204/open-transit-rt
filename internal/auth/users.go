package auth

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/tenant"
)

type AdminUser struct {
	ID             int64
	AgencyID       string
	Email          string
	DisplayName    string
	Subject        string
	Roles          []Role
	Status         string
	PasswordStatus string
	CreatedAt      time.Time
	DisabledAt     *time.Time
	LastLoginAt    *time.Time
}

type AdminUserCreateInput struct {
	AgencyID    string
	Email       string
	DisplayName string
	AuthSubject string
	Roles       []Role
	ActorID     string
	Now         time.Time
}

type AdminUserDisableInput struct {
	AgencyID string
	UserID   int64
	ActorID  string
	Reason   string
	Now      time.Time
}

type PasswordResetTokenInput struct {
	AgencyID string
	UserID   int64
	ActorID  string
	TTL      time.Duration
	Now      time.Time
}

func (s *PostgresAdminStore) ListAdminUsers(ctx context.Context, agencyID string) ([]AdminUser, error) {
	if s == nil || s.pool == nil {
		return nil, fmt.Errorf("admin auth store is unavailable")
	}
	if err := tenant.ValidateAgencyID(strings.TrimSpace(agencyID)); err != nil {
		return nil, fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	rows, err := s.pool.Query(ctx, `
		SELECT au.id, au.agency_id, au.email, COALESCE(au.display_name, ''), COALESCE(au.auth_subject, au.email),
		       au.created_at, au.disabled_at, COALESCE(pc.status, 'not_configured'), pc.last_login_at
		FROM agency_user au
		LEFT JOIN admin_password_credential pc
		  ON pc.agency_id = au.agency_id
		 AND pc.agency_user_id = au.id
		WHERE au.agency_id = $1
		ORDER BY au.email
	`, agencyID)
	if err != nil {
		return nil, fmt.Errorf("list admin users: %w", err)
	}
	defer rows.Close()
	var users []AdminUser
	for rows.Next() {
		var user AdminUser
		if err := rows.Scan(&user.ID, &user.AgencyID, &user.Email, &user.DisplayName, &user.Subject, &user.CreatedAt, &user.DisabledAt, &user.PasswordStatus, &user.LastLoginAt); err != nil {
			return nil, fmt.Errorf("scan admin user: %w", err)
		}
		if user.DisabledAt != nil {
			user.Status = "disabled"
		} else {
			user.Status = "active"
		}
		roles, err := s.rolesForUserID(ctx, user.AgencyID, user.ID)
		if err != nil {
			return nil, err
		}
		user.Roles = roles
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate admin users: %w", err)
	}
	return users, nil
}

func (s *PostgresAdminStore) CreateAdminUser(ctx context.Context, input AdminUserCreateInput) (AdminUser, error) {
	if s == nil || s.pool == nil {
		return AdminUser{}, fmt.Errorf("admin auth store is unavailable")
	}
	if err := tenant.ValidateAgencyID(strings.TrimSpace(input.AgencyID)); err != nil {
		return AdminUser{}, fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	email, err := NormalizeAdminEmail(input.Email)
	if err != nil {
		return AdminUser{}, err
	}
	roles := normalizeRoles(input.Roles)
	if len(roles) == 0 {
		return AdminUser{}, fmt.Errorf("at least one role is required")
	}
	now := input.Now.UTC()
	if now.IsZero() {
		now = s.now().UTC()
	}
	subject := firstNonEmptyString(input.AuthSubject, email)
	actor := firstNonEmptyString(input.ActorID, "operator_console")

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return AdminUser{}, fmt.Errorf("begin create user transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	var user AdminUser
	err = tx.QueryRow(ctx, `
		INSERT INTO agency_user (agency_id, email, display_name, auth_subject, created_at, disabled_at)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5, NULL)
		ON CONFLICT (agency_id, email) DO UPDATE
		SET display_name = COALESCE(NULLIF(EXCLUDED.display_name, ''), agency_user.display_name),
		    auth_subject = COALESCE(NULLIF(agency_user.auth_subject, ''), EXCLUDED.auth_subject),
		    disabled_at = NULL
		RETURNING id, agency_id, email, COALESCE(display_name, ''), COALESCE(auth_subject, email), created_at, disabled_at
	`, input.AgencyID, email, strings.TrimSpace(input.DisplayName), subject, now).Scan(&user.ID, &user.AgencyID, &user.Email, &user.DisplayName, &user.Subject, &user.CreatedAt, &user.DisabledAt)
	if err != nil {
		return AdminUser{}, fmt.Errorf("upsert admin user: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		DELETE FROM role_binding
		WHERE agency_id = $1
		  AND agency_user_id = $2
		  AND NOT (role = ANY($3::text[]))
	`, input.AgencyID, user.ID, roleStrings(roles)); err != nil {
		return AdminUser{}, fmt.Errorf("replace role bindings: %w", err)
	}
	for _, role := range roles {
		if _, err := tx.Exec(ctx, `
			INSERT INTO role_binding (agency_id, agency_user_id, role)
			VALUES ($1, $2, $3)
			ON CONFLICT (agency_id, agency_user_id, role) DO NOTHING
		`, input.AgencyID, user.ID, role); err != nil {
			return AdminUser{}, fmt.Errorf("bind role %s: %w", role, err)
		}
	}
	_, _ = tx.Exec(ctx, `
		INSERT INTO audit_log (agency_id, actor_id, action, entity_type, entity_id, new_value_json, created_at)
		VALUES ($1, $2, 'admin.user.create', 'agency_user', $3, jsonb_build_object('email', $4, 'roles', $5::text[]), $6)
	`, input.AgencyID, actor, strconv.FormatInt(user.ID, 10), email, roleStrings(roles), now)
	if err := tx.Commit(ctx); err != nil {
		return AdminUser{}, fmt.Errorf("commit create user transaction: %w", err)
	}
	user.Roles = roles
	user.Status = "active"
	user.PasswordStatus = "not_configured"
	return user, nil
}

func (s *PostgresAdminStore) DisableAdminUser(ctx context.Context, input AdminUserDisableInput) error {
	if s == nil || s.pool == nil {
		return fmt.Errorf("admin auth store is unavailable")
	}
	if err := tenant.ValidateAgencyID(strings.TrimSpace(input.AgencyID)); err != nil {
		return fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	if input.UserID <= 0 {
		return fmt.Errorf("user id is required")
	}
	now := input.Now.UTC()
	if now.IsZero() {
		now = s.now().UTC()
	}
	actor := firstNonEmptyString(input.ActorID, "operator_console")
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin disable user transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	tag, err := tx.Exec(ctx, `
		UPDATE agency_user
		SET disabled_at = $1
		WHERE agency_id = $2 AND id = $3
	`, now, input.AgencyID, input.UserID)
	if err != nil {
		return fmt.Errorf("disable admin user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("admin user not found")
	}
	_, _ = tx.Exec(ctx, `
		UPDATE admin_password_credential
		SET status = 'disabled', updated_at = $1
		WHERE agency_id = $2 AND agency_user_id = $3
	`, now, input.AgencyID, input.UserID)
	_, _ = tx.Exec(ctx, `
		INSERT INTO audit_log (agency_id, actor_id, action, entity_type, entity_id, reason, created_at)
		VALUES ($1, $2, 'admin.user.disable', 'agency_user', $3, NULLIF($4, ''), $5)
	`, input.AgencyID, actor, strconv.FormatInt(input.UserID, 10), strings.TrimSpace(input.Reason), now)
	return tx.Commit(ctx)
}

func (s *PostgresAdminStore) CreatePasswordResetToken(ctx context.Context, input PasswordResetTokenInput, tokenHash string) (BootstrapResult, error) {
	if s == nil || s.pool == nil {
		return BootstrapResult{}, fmt.Errorf("admin auth store is unavailable")
	}
	if err := tenant.ValidateAgencyID(strings.TrimSpace(input.AgencyID)); err != nil {
		return BootstrapResult{}, fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	if input.UserID <= 0 {
		return BootstrapResult{}, fmt.Errorf("user id is required")
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
	actor := firstNonEmptyString(input.ActorID, "operator_console")
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("begin password reset transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	var email string
	var subject string
	err = tx.QueryRow(ctx, `
		SELECT email, COALESCE(auth_subject, email)
		FROM agency_user
		WHERE agency_id = $1 AND id = $2 AND disabled_at IS NULL
	`, input.AgencyID, input.UserID).Scan(&email, &subject)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("load admin user for reset: %w", err)
	}
	var tokenID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO admin_bootstrap_token (agency_id, agency_user_id, purpose, token_hash, expires_at, created_by, created_at)
		VALUES ($1, $2, 'password_reset', $3, $4, $5, $6)
		RETURNING id
	`, input.AgencyID, input.UserID, tokenHash, expiresAt, actor, now).Scan(&tokenID)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("insert password reset token: %w", err)
	}
	_, _ = tx.Exec(ctx, `
		INSERT INTO audit_log (agency_id, actor_id, action, entity_type, entity_id, new_value_json, created_at)
		VALUES ($1, $2, 'admin.password_reset_token.create', 'agency_user', $3, jsonb_build_object('expires_at', $4::text), $5)
	`, input.AgencyID, actor, strconv.FormatInt(input.UserID, 10), expiresAt.Format(time.RFC3339), now)
	if err := tx.Commit(ctx); err != nil {
		return BootstrapResult{}, fmt.Errorf("commit password reset transaction: %w", err)
	}
	return BootstrapResult{AgencyID: input.AgencyID, Email: email, Subject: subject, UserID: input.UserID, TokenID: tokenID, ExpiresAt: expiresAt}, nil
}

func (s *PostgresAdminStore) rolesForUserID(ctx context.Context, agencyID string, userID int64) ([]Role, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT role
		FROM role_binding
		WHERE agency_id = $1 AND agency_user_id = $2
		ORDER BY role
	`, agencyID, userID)
	if err != nil {
		return nil, fmt.Errorf("query role bindings: %w", err)
	}
	defer rows.Close()
	return RowsToRoles(rows)
}

func normalizeRoles(roles []Role) []Role {
	seen := map[Role]bool{}
	var out []Role
	for _, role := range roles {
		switch role {
		case RoleAdmin, RoleEditor, RoleOperator, RoleReadOnly:
			if !seen[role] {
				out = append(out, role)
				seen[role] = true
			}
		}
	}
	return out
}

func roleStrings(roles []Role) []string {
	values := make([]string, 0, len(roles))
	for _, role := range roles {
		values = append(values, string(role))
	}
	return values
}
