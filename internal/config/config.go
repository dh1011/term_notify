package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration.
type Config struct {
	Server   string `yaml:"server"`
	Topic    string `yaml:"topic"`
	Priority string `yaml:"priority"`
	Token    string `yaml:"token"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Server:   "ntfy.sh",
		Priority: "default",
	}
}

// ConfigDir returns the platform-appropriate config directory.
func ConfigDir() (string, error) {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		return filepath.Join(appData, "term_notify"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return filepath.Join(home, ".config", "term_notify"), nil
}

// ConfigPath returns the full path to the config file.
func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// Load reads the config file, applies env var overrides, and returns the config.
// If the config file doesn't exist, returns defaults with env overrides.
func Load() (*Config, error) {
	cfg := DefaultConfig()

	path, err := ConfigPath()
	if err != nil {
		applyEnvOverrides(cfg)
		return cfg, nil
	}

	data, err := os.ReadFile(path) // #nosec G304 — config path is trusted
	if err != nil {
		if os.IsNotExist(err) {
			applyEnvOverrides(cfg)
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	applyEnvOverrides(cfg)
	return cfg, nil
}

// Save writes the config to disk, creating the directory if needed.
func Save(cfg *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// applyEnvOverrides overrides config values with environment variables if set.
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("TN_SERVER"); v != "" {
		cfg.Server = v
	}
	if v := os.Getenv("TN_TOPIC"); v != "" {
		cfg.Topic = v
	}
	if v := os.Getenv("TN_TOKEN"); v != "" {
		cfg.Token = v
	}
	if v := os.Getenv("TN_PRIORITY"); v != "" {
		cfg.Priority = v
	}
}
