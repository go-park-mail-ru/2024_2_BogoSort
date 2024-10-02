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

func InitFromEnv() {
	cfg.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")
	expirationTime, _ := time.ParseDuration(os.Getenv("JWT_EXPIRATION_TIME"))
	cfg.JWT.ExpirationTime = expirationTime
	cfg.JWT.Issuer = os.Getenv("JWT_ISSUER")
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
