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
	IP              string        `yaml:"ip" default:"0.0.0.0"`
	Port            int           `yaml:"port" default:"8080"`
	ReadTimeout     time.Duration `yaml:"read_timeout" default:"10s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" default:"10s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" default:"10s"`
}

type SessionConfig struct {
	ExpirationTime time.Duration `yaml:"expiration_time" default:"12h"`
	SecureCookie   bool          `yaml:"secure_cookie" default:"false"`
}

type DBConfig struct {
	MaxConns          int           `yaml:"max_conns" default:"10"`
	MinConns          int           `yaml:"min_conns" default:"2"`
	MaxConnLifetime   time.Duration `yaml:"max_conn_lifetime" default:"1h"`
	MaxConnIdleTime   time.Duration `yaml:"max_conn_idle_time" default:"30m"`
	HealthCheckPeriod time.Duration `yaml:"health_check_period" default:"1m"`
}

type Config struct {
	Server           ServerConfig  `yaml:"server"`
	Session          SessionConfig `yaml:"session"`
	DB               DBConfig      `yaml:"db"`
	PGIP             string        `yaml:"pg_ip"`
	PGPort           int           `yaml:"pg_port"`
	PGUser           string        `yaml:"pg_user"`
	PGPass           string        `yaml:"pg_password"`
	PGTimeout        time.Duration `yaml:"pg_timeout" default:"5s"`
	PGDB             string        `yaml:"pg_db"`
	RdAddr           string        `yaml:"rd_addr"`
	RdPass           string        `yaml:"rd_password"`
	RdDB             int           `yaml:"rd_db"`
	Static           StaticConfig  `yaml:"static"`
	CSRFSecret       string        `yaml:"csrf_secret"`
	AuthPort         int           `yaml:"auth_port"`
	AuthHost         string        `yaml:"auth_host"`
	CartPurchaseHost string        `yaml:"cart_purchase_host"`
	CartPurchasePort int           `yaml:"cart_purchase_port"`
	StaticHost       string        `yaml:"static_host"`
	StaticPort       int           `yaml:"static_port"`
	SearchBatchSize  int           `yaml:"search_batch_size"`
}

type StaticConfig struct {
	IP      string        `yaml:"ip"            default:"0.0.0.0"`
	Port    int           `yaml:"port"`
	Path    string        `yaml:"path"`
	MaxSize int           `yaml:"max_size"`
	Timeout time.Duration `yaml:"timeout"`
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

	if port := os.Getenv("AUTH_PORT"); port != "" {
		cfg.AuthPort, _ = strconv.Atoi(port)
	}
	if host := os.Getenv("AUTH_HOST"); host != "" {
		cfg.AuthHost = host
	}

	if batchSize := os.Getenv("SEARCH_BATCH_SIZE"); batchSize != "" {
		cfg.SearchBatchSize, _ = strconv.Atoi(batchSize)
	}

	if maxSize := os.Getenv("STATIC_MAX_SIZE"); maxSize != "" {
		cfg.Static.MaxSize, _ = strconv.Atoi(maxSize)
	}

	return cfg, nil
}

func GetStaticMaxSize() int {
	return cfg.Static.MaxSize
}

func GetSearchBatchSize() int {
	return cfg.SearchBatchSize
}

func GetServerAddress() string {
	return fmt.Sprintf(":%d", cfg.Server.Port)
}

func GetAuthAddress() string {
	return fmt.Sprintf("%s:%d", cfg.AuthHost, cfg.AuthPort)
}

func GetCartPurchaseAddress() string {
	return fmt.Sprintf("%s:%d", cfg.CartPurchaseHost, cfg.CartPurchasePort)
}

func GetStaticAddress() string {
	return fmt.Sprintf("%s:%d", cfg.StaticHost, cfg.StaticPort)
}

func (cfg *Config) GetConnectURL() string {
	user := cfg.PGUser
	pass := cfg.PGPass

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=verify-ca&sslrootcert=/etc/postgresql/certs/root.crt&pool_max_conns=%d&pool_min_conns=%d&pool_max_conn_lifetime=%s&pool_max_conn_idle_time=%s&pool_health_check_period=%s",
		user,
		pass,
		cfg.PGIP,
		cfg.PGPort,
		cfg.PGDB,
		cfg.DB.MaxConns,
		cfg.DB.MinConns,
		cfg.DB.MaxConnLifetime,
		cfg.DB.MaxConnIdleTime,
		cfg.DB.HealthCheckPeriod,
	)
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

func (cfg *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", cfg.Server.IP, cfg.Server.Port)
}

func GetStaticConfig() StaticConfig {
	return cfg.Static
}

func GetCSRFSecret() string {
	return cfg.CSRFSecret
}
