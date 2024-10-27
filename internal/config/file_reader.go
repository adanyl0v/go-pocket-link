package config

import "github.com/ilyakaznacheev/cleanenv"

type FileReader struct {
	path string
}

func NewFileReader(path string) *FileReader {
	return &FileReader{path: path}
}

func (r *FileReader) Read() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(r.path, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
