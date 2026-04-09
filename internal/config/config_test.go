package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	// Clear any env vars that might be set
	for _, key := range []string{"PORT", "DATABASE_PATH", "PARCEL_TRACKER_PASSWORD", "DEV", "REFRESH_INTERVAL", "LAPOSTE_API_KEY"} {
		t.Setenv(key, "")
	}

	cfg := Load()

	if cfg.Port != 8080 {
		t.Errorf("expected default Port=8080, got %d", cfg.Port)
	}
	if cfg.DatabasePath != "data/parcel-tracker.db" {
		t.Errorf("expected default DatabasePath, got %q", cfg.DatabasePath)
	}
	if cfg.Password != "" {
		t.Errorf("expected empty default Password, got %q", cfg.Password)
	}
	if cfg.Dev {
		t.Error("expected Dev=false by default")
	}
	if cfg.RefreshInterval != 30*time.Minute {
		t.Errorf("expected default RefreshInterval=30m, got %v", cfg.RefreshInterval)
	}
	if cfg.LaPosteAPIKey != "" {
		t.Errorf("expected empty default LaPosteAPIKey, got %q", cfg.LaPosteAPIKey)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_PATH", "/tmp/test.db")
	t.Setenv("PARCEL_TRACKER_PASSWORD", "secret")
	t.Setenv("DEV", "1")
	t.Setenv("REFRESH_INTERVAL", "5m")
	t.Setenv("LAPOSTE_API_KEY", "my-key")

	cfg := Load()

	if cfg.Port != 9090 {
		t.Errorf("expected Port=9090, got %d", cfg.Port)
	}
	if cfg.DatabasePath != "/tmp/test.db" {
		t.Errorf("expected DatabasePath=/tmp/test.db, got %q", cfg.DatabasePath)
	}
	if cfg.Password != "secret" {
		t.Errorf("expected Password=secret, got %q", cfg.Password)
	}
	if !cfg.Dev {
		t.Error("expected Dev=true")
	}
	if cfg.RefreshInterval != 5*time.Minute {
		t.Errorf("expected RefreshInterval=5m, got %v", cfg.RefreshInterval)
	}
	if cfg.LaPosteAPIKey != "my-key" {
		t.Errorf("expected LaPosteAPIKey=my-key, got %q", cfg.LaPosteAPIKey)
	}
}

func TestLoadInvalidPortFallsBack(t *testing.T) {
	t.Setenv("PORT", "not-a-number")

	cfg := Load()
	if cfg.Port != 8080 {
		t.Errorf("expected fallback Port=8080, got %d", cfg.Port)
	}
}

func TestLoadInvalidDurationFallsBack(t *testing.T) {
	t.Setenv("REFRESH_INTERVAL", "invalid")

	cfg := Load()
	if cfg.RefreshInterval != 30*time.Minute {
		t.Errorf("expected fallback RefreshInterval=30m, got %v", cfg.RefreshInterval)
	}
}
