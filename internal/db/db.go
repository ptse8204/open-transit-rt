package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"open-transit-rt/internal/appconfig"

	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultDatabaseURL = "postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable"

type Config struct {
	DatabaseURL string
	MaxConns    int32
}

func LoadConfigFromEnv() Config {
	databaseURL := getenv("DATABASE_URL", defaultDatabaseURL)
	if appconfig.IsProduction() && os.Getenv("DATABASE_URL") == "" {
		databaseURL = ""
	}
	cfg := Config{
		DatabaseURL: databaseURL,
		MaxConns:    10,
	}
	if raw := os.Getenv("DB_MAX_CONNS"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			cfg.MaxConns = int32(parsed)
		}
	}
	return cfg
}

func Connect(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	poolConfig.MaxConns = cfg.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
