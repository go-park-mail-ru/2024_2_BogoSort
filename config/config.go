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
		SecureCookie   bool          `yaml:"secure_cookie" default:"false"`
	} `yaml:"session"`
	PGIP   string `yaml:"pg_ip"   default:"postgres"`
	PGPort int    `yaml:"pg_port" default:"5432"`
	PGUser string `yaml:"pg_user" default:"postgres"`
	PGPass string `yaml:"pg_password" default:"postgres"`
	PGDB   string `yaml:"pg_db" default:"emporiumdb"`
	RdAddr string `yaml:"rd_addr" default:"redis:6379"`
	RdPass string `yaml:"rd_password" default:""`
	RdDB   int    `yaml:"rd_db" default:"0"`
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

	cfg.PGUser = os.Getenv("POSTGRES_USER")
	cfg.PGPass = os.Getenv("POSTGRES_PASSWORD")
	cfg.PGIP = os.Getenv("POSTGRES_HOST")

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	cfg.PGPort, _ = strconv.Atoi(port)

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

func GetSessionExpirationTime() time.Duration {
	return cfg.Session.ExpirationTime
}
