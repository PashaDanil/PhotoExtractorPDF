package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	Redis  RedisConfig
	MinIO  MinIOConfig
}

type ServerConfig struct {
	Port string
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type MinIOConfig struct {
	URL      string
	User     string
	Password string
	UseSSL   bool
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		MinIO: MinIOConfig{
			URL:      getEnv("MINIO_URL", "localhost:9000"),
			User:     getEnv("MINIO_ROOT_USER", "minioadmin"),
			Password: getEnv("MINIO_ROOT_PASSWORD", "minioadmin"),
			UseSSL:   getEnvAsBool("MINIO_USE_SSL", false),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Redis.URL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	if c.MinIO.URL == "" {
		return fmt.Errorf("MINIO_URL is required")
	}
	if c.MinIO.User == "" {
		return fmt.Errorf("MINIO_ROOT_USER is required")
	}
	if c.MinIO.Password == "" {
		return fmt.Errorf("MINIO_ROOT_PASSWORD is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
