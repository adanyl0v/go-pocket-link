package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go-pocket-link/pkg/errb"
)

type FileReader struct {
	path string
}

func NewFileReader(path string) *FileReader {
	return &FileReader{path: path}
}

func (r *FileReader) Path() string {
	return r.path
}

func (r *FileReader) Read() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(r.path, cfg)
	if err != nil {
		return nil, errb.Errorf("failed to read %s: %s", r.path, err.Error())
	}
	return cfg, nil
}

func (r *FileReader) MustRead() *Config {
	cfg, err := r.Read()
	if err != nil {
		panic(err)
	}
	return cfg
}
