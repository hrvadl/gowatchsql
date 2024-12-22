package cfg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"

	"golang.org/x/exp/maps"
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
	f, err := os.OpenFile(cfgPath, os.O_RDWR|os.O_CREATE, filemode)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}

	cfg := Config{file: f, Connections: make(map[string]Connection)}
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
	file        *os.File
	Connections map[string]Connection `yaml:"connections"`
}

type Connection struct {
	Name       string    `yaml:"name"`
	LastUsedAt time.Time `yaml:"last_used_at"`
	DSN        string    `yaml:"dsn"`
}

func (c *Config) AddConnection(ctx context.Context, name, dsn string) error {
	c.Connections[dsn] = Connection{Name: name, DSN: dsn, LastUsedAt: time.Now()}
	return c.Save()
}

func (c *Config) GetConnections(context.Context) []Connection {
	conns := maps.Values(c.Connections)
	slices.SortStableFunc(conns, func(a, b Connection) int {
		return b.LastUsedAt.Compare(a.LastUsedAt)
	})

	return conns
}

func (c *Config) Save() error {
	if err := c.file.Truncate(0); err != nil {
		return fmt.Errorf("truncate config: %w", err)
	}

	if _, err := c.file.Seek(0, 0); err != nil {
		return fmt.Errorf("seek config: %w", err)
	}

	if err := yaml.NewEncoder(c.file).Encode(c); err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	return nil
}

func (c *Config) Close() error {
	return c.file.Close()
}
