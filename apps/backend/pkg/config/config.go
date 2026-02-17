package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerConfig struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	RedisConfig struct {
		URL      string `yaml:"url"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`

	MinIOConfig struct {
		URL      string `yaml:"url"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		UseSSL   bool   `yaml:"use_ssl"`
	} `yaml:"minio"`

	RabbitMQConfig struct {
		URL      string `yaml:"url"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"rabbitmq"`

	PostgresConfig struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Port     string `yaml:"port"`
		SSLMode  string `yaml:"sslmode"`
		TimeZone string `yaml:"timezone"`
	} `yaml:"postgres"`
}

func New(config string) (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	if err := cleanenv.ReadConfig(config, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
