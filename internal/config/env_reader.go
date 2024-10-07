package config

import "github.com/ilyakaznacheev/cleanenv"

type EnvReader struct{}

func NewEnvReader() *EnvReader {
	return &EnvReader{}
}

func (r *EnvReader) Read() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
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
