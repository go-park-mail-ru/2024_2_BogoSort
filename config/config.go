package config

import (
	"fmt"
	"os"
	"time"

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
	f, err := os.Open("./internal/config/config.yaml")
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	return nil
}

func GetServerAddress() string {
	return fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
}

func GetJWTSecretKey() string {
	if envKey := os.Getenv("JWT_SECRET_KEY"); envKey != "" {
		return envKey
	}
	return cfg.JWT.SecretKey
}

func GetJWTExpirationTime() time.Duration {
	return cfg.JWT.ExpirationTime
}

func GetJWTIssuer() string {
	return cfg.JWT.Issuer
}