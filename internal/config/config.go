package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         int
	DatabasePath string
	Password     string
	Dev          bool
}

func Load() Config {
	return Config{
		Port:         envInt("PORT", 8080),
		DatabasePath: envStr("DATABASE_PATH", "data/parcel-tracker.db"),
		Password:     envStr("PARCEL_TRACKER_PASSWORD", ""),
		Dev:          envStr("DEV", "") != "",
	}
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
