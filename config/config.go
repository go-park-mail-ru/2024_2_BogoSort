package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	JWT struct {
		SecretKey      string        `yaml:"secret_key"`
		ExpirationTime time.Duration `yaml:"expiration_time"`
		Issuer         string        `yaml:"issuer"`
	} `yaml:"jwt"`
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

func GetServerAddress() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return ":" + port
}

func GetJWTSecretKey() string {
	if envKey := os.Getenv("JWT_SECRET_KEY"); envKey != "" {
		return envKey
	}

	return cfg.JWT.SecretKey
}

func GetJWTExpirationTime() time.Duration {
	if envTime := os.Getenv("JWT_EXPIRATION_TIME"); envTime != "" {
		duration, err := time.ParseDuration(envTime)
		if err == nil {
			return duration
		}
	}

	return cfg.JWT.ExpirationTime
}

func GetJWTIssuer() string {
	if envIssuer := os.Getenv("JWT_ISSUER"); envIssuer != "" {
		return envIssuer
	}

	return cfg.JWT.Issuer
}
