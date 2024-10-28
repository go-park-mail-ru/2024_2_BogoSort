package config

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type PostgresDatabase struct {
	IP   string `yaml:"ip"   default:"postgres"`
	Port int    `yaml:"port" default:"5432"`
	User string `yaml:"-" default:"postgres"`
	Pass string `yaml:"-" default:"postgres"`
}

type RedisDatabase struct {
	Addr     string `yaml:"addr" default:"redis:6379"`
	Password string `yaml:"-"`
	DB       int    `yaml:"db"   default:"0"`
}

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
	Postgres PostgresDatabase `yaml:"postgres"`
	Redis    RedisDatabase    `yaml:"redis"`
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

func (cfg *PostgresDatabase) GetPostgresConnectURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/emporiumDB?sslmode=disable",
		cfg.User, cfg.Pass, cfg.IP, cfg.Port)
}

func (cfg *RedisDatabase) GetRedisAddress() string {
	return cfg.Addr
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
