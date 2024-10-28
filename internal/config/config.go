package config

import "time"

const (
	EnvLocal = "local"
	EnvProd  = "prod"
	EnvDev   = "dev"
)

type Reader interface {
	Read() (*Config, error)
}

type Config struct {
	Env    string `yaml:"env" env-required:"true"`
	Server struct {
		Host         string        `yaml:"host" env-required:"true"`
		Port         int           `yaml:"port" env-required:"true"`
		ReadTimeout  time.Duration `yaml:"read_timeout" env-required:"true"`
		WriteTimeout time.Duration `yaml:"write_timeout" env-required:"true"`
		IdleTimeout  time.Duration `yaml:"idle_timeout" env-required:"true"`
	} `yaml:"server" env-required:"true"`
	Storage struct {
		Postgres struct {
			Host            string        `env:"POSTGRES_HOST" env-required:"true"`
			Port            int           `env:"POSTGRES_PORT" env-required:"true"`
			User            string        `env:"POSTGRES_USER" env-required:"true"`
			Password        string        `env:"POSTGRES_PASSWORD" env-required:"true"`
			Name            string        `env:"POSTGRES_DATABASE" env-required:"true"`
			SslMode         string        `env:"POSTGRES_SSL_MODE" env-required:"true"`
			MaxOpenConns    int           `yaml:"max_open_conns" env-required:"true"`
			MaxIdleConns    int           `yaml:"max_idle_conns" env-required:"true"`
			ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-required:"true"`
			ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" env-required:"true"`
		} `yaml:"postgres" env-required:"true"`
	} `yaml:"storage" env-required:"true"`
}
