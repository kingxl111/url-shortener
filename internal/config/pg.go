package config

import (
	"os"
)

const (
	dsnEnvName = "PG_DSN"
	user       = "DB_USER"
	password   = "DB_PASSWORD"
	name       = "DB_NAME"
	host       = "DB_HOST"
	port       = "DB_PORT"
	sslmode    = "disable"
)

type PGConfig struct {
	dsn      string
	DBName   string
	Username string
	Password string
	Host     string
	Port     string
	SSLMode  string
}

func NewPGConfig() (*PGConfig, error) {
	var cfg PGConfig
	cfg.DBName = os.Getenv(name)
	cfg.Username = os.Getenv(user)
	cfg.Password = os.Getenv(password)
	cfg.Host = os.Getenv(host)
	cfg.Port = os.Getenv(port)
	cfg.SSLMode = sslmode
	cfg.dsn = os.Getenv(dsnEnvName)

	return &cfg, nil
}
