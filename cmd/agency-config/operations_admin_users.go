package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
)

const adminUsersPostMaxBytes = 16 << 10

type operationsAdminUsersView struct {
	AgencyID        string                    `json:"agency_id"`
	GeneratedAt     time.Time                 `json:"generated_at"`
	Boundary        string                    `json:"boundary"`
	Status          string                    `json:"status"`
	Notice          string                    `json:"notice,omitempty"`
	Error           string                    `json:"error,omitempty"`
	ResetLink       string                    `json:"reset_link,omitempty"`
	Users           []operationsAdminUserRow  `json:"users"`
	Roles           []operationsAdminRoleRow  `json:"roles"`
	NextAction      string                    `json:"next_action"`
	DoesNotProve    string                    `json:"does_not_prove"`
	PasswordResets  string                    `json:"password_resets"`
	FutureSSO       string                    `json:"future_sso"`
	EmailAllowlist  string                    `json:"email_allowlist"`
	ClaimBoundaries []operationsAdminClaimRow `json:"claim_boundaries"`
}

type operationsAdminUserRow struct {
	ID             int64    `json:"id"`
	Email          string   `json:"email"`
	DisplayName    string   `json:"display_name"`
	Subject        string   `json:"subject"`
	Roles          []string `json:"roles"`
	Status         string   `json:"status"`
	PasswordStatus string   `json:"password_status"`
	CreatedAt      string   `json:"created_at"`
	DisabledAt     string   `json:"disabled_at,omitempty"`
	LastLoginAt    string   `json:"last_login_at,omitempty"`
	CanDisable     bool     `json:"can_disable"`
}

type operationsAdminRoleRow struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type operationsAdminClaimRow struct {
	Label  string `json:"label"`
	Status string `json:"status"`
}

func (h *handler) renderAdminUsers(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	switch r.Method {
	case http.MethodGet:
		page := h.buildOperationsPage(r, principal, "admin-users")
		renderOperationsTemplate(w, "admin-users", page)
	case http.MethodPost:
		h.adminUsersPost(w, r, principal)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) renderAdminUsersJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	page := h.buildOperationsPage(r, principal, "admin-users")
	page.AdminUsers.ResetLink = ""
	writeJSON(w, http.StatusOK, page.AdminUsers)
}

func (h *handler) adminUsersPost(w http.ResponseWriter, r *http.Request, principal auth.Principal) {
	notice := ""
	errText := ""
	resetLink := ""
	if h.users == nil {
		errText = "admin user store is not available in this runtime"
	} else if err := r.ParseForm(); err != nil {
		errText = "user-management form could not be read"
	} else if auth.RejectAgencyConflict(w, r.FormValue("agency_id"), principal) {
		return
	} else {
		switch strings.TrimSpace(r.FormValue("action")) {
		case "create_user":
			created, err := h.createAdminUserFromForm(r.Context(), principal, r)
			if err != nil {
				errText = safeAdminUsersError(err)
			} else {
				notice = fmt.Sprintf("Created or updated %s for agency %s. Generate a password reset link before sharing browser access.", created.Email, principal.AgencyID)
			}
		case "disable_user":
			if err := h.disableAdminUserFromForm(r.Context(), principal, r); err != nil {
				errText = safeAdminUsersError(err)
			} else {
				notice = "User disabled for this agency. Existing cookie sessions expire at their normal TTL; rotate the JWT secret if immediate global invalidation is required."
			}
		case "password_reset":
			link, result, err := h.passwordResetLinkFromForm(r.Context(), principal, r)
			if err != nil {
				errText = safeAdminUsersError(err)
			} else {
				resetLink = link
				notice = fmt.Sprintf("Password setup link generated for %s. It is shown once and expires at %s.", result.Email, result.ExpiresAt.UTC().Format(time.RFC3339))
			}
		default:
			errText = "unknown admin user action"
		}
	}
	page := h.buildOperationsPage(r, principal, "admin-users")
	page.AdminUsers.Notice = notice
	page.AdminUsers.Error = errText
	page.AdminUsers.ResetLink = resetLink
	renderOperationsTemplate(w, "admin-users", page)
}

func (h *handler) createAdminUserFromForm(ctx context.Context, principal auth.Principal, r *http.Request) (auth.AdminUser, error) {
	roles := adminUserRolesFromForm(r)
	if len(roles) == 0 {
		return auth.AdminUser{}, fmt.Errorf("choose at least one role")
	}
	return h.users.CreateAdminUser(ctx, auth.AdminUserCreateInput{
		AgencyID:    principal.AgencyID,
		Email:       r.FormValue("email"),
		DisplayName: r.FormValue("display_name"),
		AuthSubject: r.FormValue("auth_subject"),
		Roles:       roles,
		ActorID:     principal.Subject,
		Now:         time.Now().UTC(),
	})
}

func (h *handler) disableAdminUserFromForm(ctx context.Context, principal auth.Principal, r *http.Request) error {
	userID, err := strconv.ParseInt(strings.TrimSpace(r.FormValue("user_id")), 10, 64)
	if err != nil || userID <= 0 {
		return fmt.Errorf("valid user id is required")
	}
	if h.adminUsersTargetIsCurrentSubject(ctx, principal, userID) {
		return fmt.Errorf("current signed-in user cannot be disabled from this page")
	}
	return h.users.DisableAdminUser(ctx, auth.AdminUserDisableInput{
		AgencyID: principal.AgencyID,
		UserID:   userID,
		ActorID:  principal.Subject,
		Reason:   r.FormValue("reason"),
		Now:      time.Now().UTC(),
	})
}

func (h *handler) adminUsersTargetIsCurrentSubject(ctx context.Context, principal auth.Principal, userID int64) bool {
	if h.users == nil {
		return false
	}
	users, err := h.users.ListAdminUsers(ctx, principal.AgencyID)
	if err != nil {
		return false
	}
	for _, user := range users {
		if user.ID == userID && strings.TrimSpace(user.Subject) == principal.Subject {
			return true
		}
	}
	return false
}

func (h *handler) passwordResetLinkFromForm(ctx context.Context, principal auth.Principal, r *http.Request) (string, auth.BootstrapResult, error) {
	userID, err := strconv.ParseInt(strings.TrimSpace(r.FormValue("user_id")), 10, 64)
	if err != nil || userID <= 0 {
		return "", auth.BootstrapResult{}, fmt.Errorf("valid user id is required")
	}
	token, err := auth.GenerateBootstrapToken()
	if err != nil {
		return "", auth.BootstrapResult{}, err
	}
	result, err := h.users.CreatePasswordResetToken(ctx, auth.PasswordResetTokenInput{
		AgencyID: principal.AgencyID,
		UserID:   userID,
		ActorID:  principal.Subject,
		TTL:      30 * time.Minute,
		Now:      time.Now().UTC(),
	}, auth.HashBootstrapToken(token))
	if err != nil {
		return "", auth.BootstrapResult{}, err
	}
	link, err := firstAdminSetupLink(adminUsersBaseURL(r), token)
	if err != nil {
		return "", auth.BootstrapResult{}, err
	}
	return link, result, nil
}

func (h *handler) buildOperationsAdminUsersView(r *http.Request, principal auth.Principal, resetLink string) operationsAdminUsersView {
	view := operationsAdminUsersView{
		AgencyID:       principal.AgencyID,
		GeneratedAt:    time.Now().UTC().Truncate(time.Second),
		Boundary:       "Private admin-only user management. This page scopes every action to the signed-in agency, never accepts browser-supplied agency switching, and shows reset tokens only once.",
		Status:         "available",
		ResetLink:      resetLink,
		Roles:          adminUserRoleRows(),
		NextAction:     "Create or update staff users, assign the smallest needed role, then generate a one-time password setup link only when the recipient is ready.",
		DoesNotProve:   "User management does not prove SSO support, compliance, agency adoption, consumer acceptance, hosted service availability, uptime, or production readiness.",
		PasswordResets: "Password setup and reset links are short lived, single use, and stored only as hashes.",
		FutureSSO:      "SSO/OIDC is not implemented. A future provider must verify identity externally, map claims to this same internal agency user and role model, then issue the existing admin_session cookie.",
		EmailAllowlist: "Username/password invitations are limited to explicitly created users in this agency. Broad self-registration is not enabled.",
		ClaimBoundaries: []operationsAdminClaimRow{
			{Label: "SSO/OIDC", Status: "not implemented"},
			{Label: "Cross-agency administration", Status: "blocked"},
			{Label: "Raw credential display", Status: "blocked"},
			{Label: "Password reset token storage", Status: "hash only"},
		},
	}
	if h.users == nil {
		view.Status = "unavailable"
		view.Error = "admin user store is not available in this runtime"
		return view
	}
	users, err := h.users.ListAdminUsers(r.Context(), principal.AgencyID)
	if err != nil {
		view.Status = "error"
		view.Error = "admin users could not be loaded"
		return view
	}
	for _, user := range users {
		view.Users = append(view.Users, adminUserRow(user, principal))
	}
	if len(view.Users) == 0 {
		view.Status = "empty"
		view.Notice = "No admin users are visible for this agency yet."
	}
	return view
}

func adminUserRolesFromForm(r *http.Request) []auth.Role {
	var roles []auth.Role
	for _, raw := range r.Form["role"] {
		switch auth.Role(strings.TrimSpace(raw)) {
		case auth.RoleAdmin:
			roles = append(roles, auth.RoleAdmin)
		case auth.RoleEditor:
			roles = append(roles, auth.RoleEditor)
		case auth.RoleOperator:
			roles = append(roles, auth.RoleOperator)
		case auth.RoleReadOnly:
			roles = append(roles, auth.RoleReadOnly)
		}
	}
	return roles
}

func adminUserRow(user auth.AdminUser, principal auth.Principal) operationsAdminUserRow {
	row := operationsAdminUserRow{
		ID:             user.ID,
		Email:          user.Email,
		DisplayName:    firstNonEmpty(strings.TrimSpace(user.DisplayName), "not set"),
		Subject:        firstNonEmpty(strings.TrimSpace(user.Subject), user.Email),
		Roles:          safePrincipalRoles(user.Roles),
		Status:         firstNonEmpty(strings.TrimSpace(user.Status), "active"),
		PasswordStatus: firstNonEmpty(strings.TrimSpace(user.PasswordStatus), "not_configured"),
		CreatedAt:      formatAdminUsersTime(user.CreatedAt),
		DisabledAt:     formatAdminUsersTimePtr(user.DisabledAt),
		LastLoginAt:    formatAdminUsersTimePtr(user.LastLoginAt),
		CanDisable:     user.DisabledAt == nil && user.Subject != principal.Subject,
	}
	return row
}

func adminUserRoleRows() []operationsAdminRoleRow {
	return []operationsAdminRoleRow{
		{ID: string(auth.RoleAdmin), Label: "Admin", Description: "Can manage users, roles, credentials, and all private operations pages."},
		{ID: string(auth.RoleEditor), Label: "Editor", Description: "Can use edit-capable operations routes where explicitly allowed; cannot create admins."},
		{ID: string(auth.RoleOperator), Label: "Operator", Description: "Can review operations and run operator-safe review workflows; cannot change user access."},
		{ID: string(auth.RoleReadOnly), Label: "Read only", Description: "Can review private status pages only."},
	}
}

func adminUsersBaseURL(r *http.Request) string {
	if configured := firstNonEmpty(os.Getenv("ADMIN_BASE_URL"), os.Getenv("PUBLIC_BASE_URL")); configured != "" {
		return configured
	}
	scheme := "http"
	if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	host := strings.TrimSpace(r.Host)
	if host == "" {
		host = "localhost:8080"
	}
	base := url.URL{Scheme: scheme, Host: host}
	return base.String()
}

func formatAdminUsersTime(t time.Time) string {
	if t.IsZero() {
		return "not available"
	}
	return t.UTC().Format(time.RFC3339)
}

func formatAdminUsersTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return "not available"
	}
	return t.UTC().Format(time.RFC3339)
}

func safeAdminUsersError(err error) string {
	if err == nil {
		return ""
	}
	text := strings.TrimSpace(err.Error())
	if text == "" {
		return "admin user action failed"
	}
	if len(text) > 160 {
		text = text[:160]
	}
	return text
}
