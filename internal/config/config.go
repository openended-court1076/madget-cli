package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type LocalConfig struct {
	RegistryURL string `json:"registry_url"`
	Token       string `json:"token,omitempty"`
}

func Load(path string) (LocalConfig, error) {
	var cfg LocalConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func Save(path string, cfg LocalConfig) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
