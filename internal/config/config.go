package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	Port            int
	DatabasePath    string
	Password        string
	Dev             bool
	RefreshInterval time.Duration
	LaPosteAPIKey   string
}

func Load() Config {
	return Config{
		Port:            envInt("PORT", 8080),
		DatabasePath:    envStr("DATABASE_PATH", "data/parcel-tracker.db"),
		Password:        envStr("PARCEL_TRACKER_PASSWORD", ""),
		Dev:             envStr("DEV", "") != "",
		RefreshInterval: envDuration("REFRESH_INTERVAL", 30*time.Minute),
		LaPosteAPIKey:   envStr("LAPOSTE_API_KEY", ""),
	}
}

// Validate checks that configuration values are within acceptable ranges.
func (c Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("PORT must be between 1 and 65535, got %d", c.Port)
	}
	if c.DatabasePath == "" {
		return fmt.Errorf("DATABASE_PATH must not be empty")
	}
	dir := filepath.Dir(c.DatabasePath)
	if info, err := os.Stat(dir); err == nil && !info.IsDir() {
		return fmt.Errorf("DATABASE_PATH parent %q is not a directory", dir)
	}
	if c.RefreshInterval > 0 && c.RefreshInterval < time.Minute {
		return fmt.Errorf("REFRESH_INTERVAL must be at least 1m or 0 to disable, got %s", c.RefreshInterval)
	}
	return nil
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
