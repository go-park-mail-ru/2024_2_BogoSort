package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port            int           `yaml:"port"`
		Host            string        `yaml:"host"`
		ReadTimeout     time.Duration `yaml:"read_timeout"`
		WriteTimeout    time.Duration `yaml:"write_timeout"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	} `yaml:"server"`
	Session struct {
		ExpirationTime time.Duration `yaml:"expiration_time"`
	} `yaml:"session"`
	DB struct {
		DBUser     string `yaml:"db_user"`
		DBPassword string `yaml:"db_password"`
		DBHost     string `yaml:"db_host"`
		DBPort     int    `yaml:"db_port"`
		DBName     string `yaml:"db_name"`
	} `yaml:"db"`
}

var cfg Config

func ServerInit() error {
	file, err := os.Open("./config/config.yaml")
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}

	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)

	if err != nil {
		return errors.Wrap(err, "failed to decode config file")
	}

	return nil
}

func DBInit() (*Config, error) {
	cfg := &Config{
		DB: struct {
			DBUser     string `yaml:"db_user"`
			DBPassword string `yaml:"db_password"`
			DBHost     string `yaml:"db_host"`
			DBPort     int    `yaml:"db_port"`
			DBName     string `yaml:"db_name"`
		}{},
	}

	return cfg, nil
}

func InitFromEnv() {
	expirationTime, _ := time.ParseDuration(os.Getenv("SESSION_EXPIRATION_TIME"))
	cfg.Session.ExpirationTime = expirationTime
}

func GetServerAddress() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return ":" + port
}

func GetReadTimeout() time.Duration {
	return cfg.Server.ReadTimeout
}

func GetWriteTimeout() time.Duration {
	return cfg.Server.WriteTimeout
}

func GetShutdownTimeout() time.Duration {
	return cfg.Server.ShutdownTimeout
}

func GetSessionExpirationTime() time.Duration {
	return cfg.Session.ExpirationTime
}
