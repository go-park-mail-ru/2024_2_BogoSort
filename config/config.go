package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port         int           `yaml:"port"`
		Host         string        `yaml:"host"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	} `yaml:"server"`
	Session struct {
		ExpirationTime time.Duration `yaml:"expiration_time"`
	} `yaml:"session"`
}

var cfg Config

func Init() error {
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
