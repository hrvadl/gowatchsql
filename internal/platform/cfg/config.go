package cfg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	filename = "config.yaml"
	dir      = "gowatchsql"
	filemode = 0o644
)

func NewFromFile(base string) (*Config, error) {
	dirPath := filepath.Join(base, dir)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := createDir(base); err != nil {
			return nil, fmt.Errorf("create dir: %w", err)
		}
	}

	cfgPath := filepath.Join(dirPath, filename)

	f, err := os.OpenFile(cfgPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filemode)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}

	cfg := Config{file: f, Connections: make(map[string]string)}
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	return &cfg, nil
}

func createDir(base string) error {
	dirPath := filepath.Join(base, dir)
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	return nil
}

type Config struct {
	file        io.ReadWriteCloser
	Connections map[string]string `yaml:"connections"`
}

func (c *Config) Save() error {
	if err := yaml.NewEncoder(c.file).Encode(c); err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}

func (c *Config) Close() error {
	return c.file.Close()
}
