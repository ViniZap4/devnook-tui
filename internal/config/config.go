package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ServerURL string `json:"server_url"`
	Token     string `json:"token"`
}

func Load() *Config {
	cfg := &Config{
		ServerURL: "http://localhost:8080",
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg
	}

	configPath := filepath.Join(home, ".config", "devnook", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg
	}

	json.Unmarshal(data, cfg)

	if v := os.Getenv("DEVNOOK_SERVER_URL"); v != "" {
		cfg.ServerURL = v
	}
	if v := os.Getenv("DEVNOOK_TOKEN"); v != "" {
		cfg.Token = v
	}

	return cfg
}

func (c *Config) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".config", "devnook")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0o644)
}
