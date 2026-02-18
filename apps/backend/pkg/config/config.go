package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
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

	LoggerConfig struct {
		Service   string `yaml:"service"`
		Env       string `yaml:"env"`
		Version   string `yaml:"version"`
		Level     string `yaml:"level"`
		AddSource bool   `yaml:"add_source"`
	} `yaml:"logger"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
