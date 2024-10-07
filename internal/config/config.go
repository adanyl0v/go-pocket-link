package config

import "fmt"

type Reader interface {
	Read() (*Config, error)
	MustRead() *Config
}

type Config struct {
	Env string `env:"ENV" env-required:"true"`
	DB  DB
}

type DB struct {
	Postgres Postgres
}

type Postgres struct {
	Host    string `env:"DB_POSTGRES_HOST" env-required:"true"`
	Port    string `env:"DB_POSTGRES_PORT" env-default:"5432"`
	User    string `env:"DB_POSTGRES_USER" env-required:"true"`
	Pass    string `env:"DB_POSTGRES_PASS" env-required:"true"`
	Name    string `env:"DB_POSTGRES_NAME" env-required:"true"`
	SslMode string `env:"DB_POSTGRES_SSL_MODE" env-default:"disable"`
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User, p.Pass, p.Host, p.Port, p.Name, p.SslMode)
}
