package appconfig

import (
	"fmt"
	"os"
	"strings"
)

const (
	EnvDev        = "dev"
	EnvProduction = "production"
)

func AppEnv() string {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if env == "" {
		return EnvDev
	}
	return env
}

func IsProduction() bool {
	return AppEnv() == EnvProduction
}

func RequireProductionSecret(name string) error {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fmt.Errorf("%s is required in production", name)
	}
	if IsDevPlaceholder(value) {
		return fmt.Errorf("%s must not use a development placeholder in production", name)
	}
	return nil
}

func IsDevPlaceholder(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return true
	}
	return strings.Contains(normalized, "change-me") ||
		strings.Contains(normalized, "dev-") ||
		strings.Contains(normalized, "example") ||
		strings.Contains(normalized, "placeholder")
}
