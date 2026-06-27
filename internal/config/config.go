package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Kubeconfig      string        `yaml:"kubeconfig"`
	RefreshInterval time.Duration `yaml:"-"`
	Namespaces      []string      `yaml:"namespaces"`

	RefreshIntervalRaw string `yaml:"refreshInterval"`
}

func Default() Config {
	return Config{RefreshInterval: 30 * time.Second}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	if cfg.RefreshIntervalRaw != "" {
		duration, err := time.ParseDuration(cfg.RefreshIntervalRaw)
		if err != nil {
			return Config{}, err
		}
		cfg.RefreshInterval = duration
	}
	//TODO: Evaluate if RefreshInterval should have limits or constraints, if so implement validation logic
	if cfg.RefreshInterval == 0 {
		cfg.RefreshInterval = Default().RefreshInterval
	}
	if strings.HasPrefix(cfg.Kubeconfig, "~/") || strings.HasPrefix(cfg.Kubeconfig, `~\`) {
		home, err := os.UserHomeDir()
		if err != nil {
			return Config{}, err
		}
		cfg.Kubeconfig = filepath.Join(home, cfg.Kubeconfig[2:])
	}
	return cfg, nil
}
