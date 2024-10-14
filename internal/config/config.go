package config

import (
	"fmt"
	"time"
)

type Reader interface {
	Read() (*Config, error)
	MustRead() *Config
}

const (
	EnvLocal = "local"
	EnvProd  = "prod"
	EnvDev   = "dev"
)

type Config struct {
	Env    string `yaml:"env" env-required:"true"`
	Server Server `yaml:"server" env-required:"true"`
	DB     DB     `yaml:"db" env-required:"true"`
	Auth   Auth   `yaml:"auth" env-required:"true"`
	Email  Email  `yaml:"email" env-required:"true"`
}

type Server struct {
	Host         string        `yaml:"host" env-required:"true"`
	Port         int           `yaml:"port" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type DB struct {
	Host            string        `yaml:"host" env-required:"true"`
	Port            string        `yaml:"port" env-default:"5432"`
	User            string        `env:"POSTGRES_USER" env-required:"true"`
	Pass            string        `env:"POSTGRES_PASS" env-required:"true"`
	Name            string        `yaml:"name" env-required:"true"`
	SslMode         string        `yaml:"ssl_mode" env-default:"disable"`
	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"10"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"10"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"10s"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" env-default:"10s"`
}

func (p *DB) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User, p.Pass, p.Host, p.Port, p.Name, p.SslMode)
}

type Auth struct {
	Salt            string        `env:"AUTH_SALT" env-required:"true"`
	Secret          string        `env:"AUTH_SECRET" env-required:"true"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env-required:"true"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-required:"true"`
}

type Email struct {
	Username string `env:"EMAIL_USERNAME" env-required:"true"`
	Password string `env:"EMAIL_PASSWORD" env-required:"true"`
}
