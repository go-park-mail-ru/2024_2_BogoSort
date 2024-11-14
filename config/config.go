package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port            int           `yaml:"port" default:"8080"`
	Host            string        `yaml:"host" default:"localhost"`
	ReadTimeout     time.Duration `yaml:"read_timeout" default:"10s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" default:"10s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" default:"10s"`
}

type SessionConfig struct {
	ExpirationTime time.Duration `yaml:"expiration_time" default:"12h"`
	SecureCookie   bool          `yaml:"secure_cookie" default:"false"`
}

type Config struct {
	Server     ServerConfig  `yaml:"server"`
	Session    SessionConfig `yaml:"session"`
	PGIP       string        `yaml:"pg_ip"`
	PGPort     int           `yaml:"pg_port"`
	PGUser     string        `yaml:"pg_user"`
	PGPass     string        `yaml:"pg_password"`
	PGTimeout  time.Duration `yaml:"pg_timeout" default:"5s"`
	PGDB       string        `yaml:"pg_db"`
	RdAddr     string        `yaml:"rd_addr"`
	RdPass     string        `yaml:"rd_password"`
	RdDB       int           `yaml:"rd_db"`
	Static     StaticConfig  `yaml:"static"`
	CSRFSecret string        `yaml:"csrf_secret"`
	AuthAddr   string        `yaml:"auth_addr"`
}

type StaticConfig struct {
	Path    string `yaml:"path"`
	MaxSize int    `yaml:"max_size"`
}

var cfg Config

func Init() (Config, error) {
	file, err := os.Open("./config/config.yaml")
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to open config file")
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, errors.Wrap(err, "failed to decode config file")
	}

	if user := os.Getenv("PG_USER"); user != "" {
		cfg.PGUser = user
	}
	if pass := os.Getenv("PG_PASSWORD"); pass != "" {
		cfg.PGPass = pass
	}
	if host := os.Getenv("PG_IP"); host != "" {
		cfg.PGIP = host
	}
	if port := os.Getenv("PG_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.PGPort = p
		}
	}
	if expiration := os.Getenv("SESSION_EXPIRATION_TIME"); expiration != "" {
		if dur, err := time.ParseDuration(expiration); err == nil {
			cfg.Session.ExpirationTime = dur
		}
	}
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port, _ = strconv.Atoi(port)
	}

	return cfg, nil
}

func GetServerAddress() string {
	return fmt.Sprintf(":%d", cfg.Server.Port)
}

func (cfg *Config) GetConnectURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PGUser, cfg.PGPass, cfg.PGIP, cfg.PGPort, cfg.PGDB)
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

func GetStaticConfig() StaticConfig {
	return cfg.Static
}

func GetCSRFSecret() string {
	return cfg.CSRFSecret
}
