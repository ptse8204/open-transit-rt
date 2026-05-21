package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	appdb "open-transit-rt/internal/db"
	"open-transit-rt/internal/tenant"
)

type envLookup func(string) string

type bootstrapAdminLinkConfig struct {
	AgencyID    string
	Email       string
	DisplayName string
	AuthSubject string
	CreatedBy   string
	BaseURL     string
	TTL         time.Duration
	Now         time.Time
}

func runBootstrapAdminLink(args []string, stdout io.Writer, stderr io.Writer, getenv envLookup) int {
	cfg, err := parseBootstrapAdminLinkConfig(args, getenv)
	if err != nil {
		fmt.Fprintf(stderr, "bootstrap-admin-link: %v\n", err)
		return 2
	}
	token, err := auth.GenerateBootstrapToken()
	if err != nil {
		fmt.Fprintf(stderr, "bootstrap-admin-link: %v\n", err)
		return 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := appdb.Connect(ctx, appdb.LoadConfigFromEnv())
	if err != nil {
		fmt.Fprintf(stderr, "bootstrap-admin-link: connect database: %v\n", err)
		return 1
	}
	defer pool.Close()
	store := auth.NewPostgresAdminStore(pool)
	result, err := store.CreateFirstAdminBootstrap(ctx, auth.FirstAdminBootstrapInput{
		AgencyID:    cfg.AgencyID,
		Email:       cfg.Email,
		DisplayName: cfg.DisplayName,
		AuthSubject: cfg.AuthSubject,
		CreatedBy:   cfg.CreatedBy,
		TTL:         cfg.TTL,
		Now:         cfg.Now,
	}, auth.HashBootstrapToken(token))
	if err != nil {
		fmt.Fprintf(stderr, "bootstrap-admin-link: %v\n", err)
		return 1
	}
	link, err := firstAdminSetupLink(cfg.BaseURL, token)
	if err != nil {
		fmt.Fprintf(stderr, "bootstrap-admin-link: %v\n", err)
		return 1
	}
	writeBootstrapAdminLink(stdout, result, link)
	return 0
}

func parseBootstrapAdminLinkConfig(args []string, getenv envLookup) (bootstrapAdminLinkConfig, error) {
	if getenv == nil {
		getenv = func(string) string { return "" }
	}
	cfg := bootstrapAdminLinkConfig{
		AgencyID:  strings.TrimSpace(getenv("AGENCY_ID")),
		BaseURL:   firstNonEmpty(getenv("ADMIN_BASE_URL"), getenv("PUBLIC_BASE_URL"), "http://localhost:8080"),
		TTL:       30 * time.Minute,
		CreatedBy: firstNonEmpty(getenv("ADMIN_BOOTSTRAP_CREATED_BY"), "operator_console"),
		Now:       time.Now().UTC(),
	}
	fs := flag.NewFlagSet("bootstrap-admin-link", flag.ContinueOnError)
	fs.StringVar(&cfg.AgencyID, "agency-id", cfg.AgencyID, "agency id for the admin user")
	fs.StringVar(&cfg.Email, "email", "", "first admin email")
	fs.StringVar(&cfg.DisplayName, "display-name", "", "first admin display name")
	fs.StringVar(&cfg.AuthSubject, "subject", "", "internal admin subject; defaults to email")
	fs.StringVar(&cfg.CreatedBy, "created-by", cfg.CreatedBy, "audit-safe actor label for the operator command")
	fs.StringVar(&cfg.BaseURL, "base-url", cfg.BaseURL, "admin base URL used to print the one-time setup link")
	fs.DurationVar(&cfg.TTL, "ttl", cfg.TTL, "bootstrap link TTL, capped at 24h")
	fs.SetOutput(io.Discard)
	if err := fs.Parse(args); err != nil {
		return bootstrapAdminLinkConfig{}, err
	}
	cfg.AgencyID = strings.TrimSpace(cfg.AgencyID)
	if err := tenant.ValidateAgencyID(cfg.AgencyID); err != nil {
		return bootstrapAdminLinkConfig{}, fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	email, err := auth.NormalizeAdminEmail(cfg.Email)
	if err != nil {
		return bootstrapAdminLinkConfig{}, err
	}
	cfg.Email = email
	cfg.AuthSubject = strings.TrimSpace(cfg.AuthSubject)
	cfg.DisplayName = strings.TrimSpace(cfg.DisplayName)
	cfg.CreatedBy = strings.TrimSpace(cfg.CreatedBy)
	if cfg.CreatedBy == "" || strings.ContainsAny(cfg.CreatedBy, "\r\n\t") {
		return bootstrapAdminLinkConfig{}, fmt.Errorf("created-by must be a single audit-safe label")
	}
	cfg.TTL = auth.NormalizeBootstrapTTL(cfg.TTL)
	if _, err := firstAdminSetupLink(cfg.BaseURL, "token-preview"); err != nil {
		return bootstrapAdminLinkConfig{}, err
	}
	return cfg, nil
}

func firstAdminSetupLink(baseURL string, token string) (string, error) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return "", fmt.Errorf("base URL is required")
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("base URL must be absolute")
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/admin/setup/first-admin"
	values := parsed.Query()
	values.Set("token", token)
	parsed.RawQuery = values.Encode()
	parsed.Fragment = ""
	return parsed.String(), nil
}

func writeBootstrapAdminLink(w io.Writer, result auth.BootstrapResult, link string) {
	fmt.Fprintf(w, "First admin setup link for %s (%s)\n", result.Email, result.AgencyID)
	fmt.Fprintf(w, "%s\n", link)
	fmt.Fprintf(w, "Expires at: %s\n", result.ExpiresAt.UTC().Format(time.RFC3339))
	fmt.Fprintln(w, "This token is shown once. Do not paste this output into docs/evidence, screenshots, issue comments, or public logs.")
}
