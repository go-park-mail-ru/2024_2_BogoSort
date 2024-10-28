package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port            int           `yaml:"port" default:"8080"`
		Host            string        `yaml:"host" default:"localhost"`
		ReadTimeout     time.Duration `yaml:"read_timeout" default:"10s"`
		WriteTimeout    time.Duration `yaml:"write_timeout" default:"10s"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout" default:"10s"`
	} `yaml:"server"`
	Session struct {
		ExpirationTime time.Duration `yaml:"expiration_time" default:"12h"`
	} `yaml:"session"`
	Postgres struct {
		IP   string `yaml:"ip"   default:"postgres"`
		Port int    `yaml:"port" default:"5432"`
		User string `yaml:"user" default:"postgres"`
		Pass string `yaml:"password" default:"postgres"`
		DB   string `yaml:"db" default:"emporiumdb"`
	} `yaml:"postgres"`
}

var cfg Config

func Init() (Config, error) {
	file, err := os.Open("./config/config.yaml")
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to open config file")
	}

	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to decode config file")
	}

	return cfg, nil
}

func InitFromEnv() Config {
    var cfg Config

    cfg.Postgres.User = os.Getenv("POSTGRES_USER")
    cfg.Postgres.Pass = os.Getenv("POSTGRES_PASSWORD")
    cfg.Postgres.IP = os.Getenv("POSTGRES_HOST")

    port := os.Getenv("POSTGRES_PORT")
    if port == "" {
        port = "5432" 
    }
    cfg.Postgres.Port, _ = strconv.Atoi(port)

    expirationTime, _ := time.ParseDuration(os.Getenv("SESSION_EXPIRATION_TIME"))
    cfg.Session.ExpirationTime = expirationTime

    return cfg
}

func GetServerAddress() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return ":" + port
}

func (cfg *Config) GetConnectURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Pass, cfg.Postgres.IP, cfg.Postgres.Port, cfg.Postgres.DB)
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
