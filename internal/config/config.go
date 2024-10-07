package config

import (
	"fmt"
	"time"
)

type Reader interface {
	Read() (*Config, error)
	MustRead() *Config
}

type Config struct {
	Env  string `yaml:"env" env-required:"true"`
	DB   DB     `yaml:"db"`
	Auth Auth   `yaml:"auth"`
}

type DB struct {
	Postgres Postgres `yaml:"postgres"`
}

type Postgres struct {
	Host    string `yaml:"host" env-required:"true"`
	Port    string `yaml:"port" env-default:"5432"`
	User    string `env:"DB_POSTGRES_USER" env-required:"true"`
	Pass    string `env:"DB_POSTGRES_PASS" env-required:"true"`
	Name    string `yaml:"name" env-required:"true"`
	SslMode string `yaml:"ssl_mode" env-default:"disable"`
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User, p.Pass, p.Host, p.Port, p.Name, p.SslMode)
}

type Auth struct {
	Salt       string `env:"AUTH_SALT" env-required:"true"`
	Secret     string `env:"AUTH_SECRET" env-required:"true"`
	AccessTok  Token  `yaml:"access_tok" env-required:"true"`
	RefreshTok Token  `yaml:"refresh_tok" env-required:"true"`
}

type Token struct {
	ExpiresAt time.Duration `yaml:"expires_at" env-required:"true"`
}
