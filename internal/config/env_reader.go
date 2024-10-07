package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go-pocket-link/pkg/errb"
)

type EnvReader struct{}

func NewEnvReader() *EnvReader {
	return &EnvReader{}
}

func (r *EnvReader) Read() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, errb.Errorf("failed to read .env: %s", err.Error())
	}
	return cfg, nil
}

func (r *EnvReader) MustRead() *Config {
	cfg, err := r.Read()
	if err != nil {
		panic(err)
	}
	return cfg
}
