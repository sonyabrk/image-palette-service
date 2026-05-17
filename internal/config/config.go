package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	Workers     int
	CacheTTL    int
	MaxClusters int
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Workers:     getEnvInt("WORKERS", 4),
		CacheTTL:    getEnvInt("CACHE_TTL_SECONDS", 300),
		MaxClusters: getEnvInt("MAX_CLUSTERS", 8),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}
