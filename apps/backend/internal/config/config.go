package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Redis    RedisConfig
	MinIO    MinIOConfig
	RabbitMQ RabbitMQConfig
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

type RabbitMQConfig struct {
	URL      string
	User     string
	Password string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", "redis"),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		MinIO: MinIOConfig{
			URL:      getEnv("MINIO_URL", "localhost:9000"),
			User:     getEnv("MINIO_ROOT_USER", "minio"),
			Password: getEnv("MINIO_ROOT_PASSWORD", "minio12345"),
			UseSSL:   getEnvAsBool("MINIO_USE_SSL", false),
		},
		RabbitMQ: RabbitMQConfig{
			URL:      getEnv("RABBITMQ_URL", "localhost:5672"),
			User:     getEnv("RABBITMQ_USER", "rabbit"),
			Password: getEnv("RABBITMQ_PASSWORD", "rabbit123"),
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
	if c.RabbitMQ.URL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}
	if c.RabbitMQ.User == "" {
		return fmt.Errorf("RABBITMQ_USER is required")
	}
	if c.RabbitMQ.Password == "" {
		return fmt.Errorf("RABBITMQ_PASSWORD is required")
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
