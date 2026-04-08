package i2ptui

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the persistent configuration file.
type Config struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Path     string `json:"path,omitempty"`
	Password string `json:"password,omitempty"`
	Cert     string `json:"cert,omitempty"`
	Interval string `json:"interval,omitempty"`
	Theme    string `json:"theme,omitempty"`
}

// DefaultConfigPath returns the default path for the config file.
func DefaultConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(dir, "i2ptui", "config.json")
}

// LoadConfig reads and parses the config file. Returns zero Config on error.
func LoadConfig(path string) Config {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, &cfg)
	return cfg
}

// SaveConfig writes the config to disk, creating directories as needed.
func SaveConfig(path string, cfg Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
