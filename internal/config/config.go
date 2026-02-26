package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	CacheTTL         int    `json:"cache_ttl"`
	MaxJobs          int    `json:"max_jobs"`
	DefaultOrg       string `json:"default_org"`
	DefaultClonePath string `json:"default_clone_path"`
}

func DefaultConfig() Config {
	return Config{
		CacheTTL: 600,
		MaxJobs:  5,
	}
}

func configDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "ghtools")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "ghtools")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

func oldConfigPath() string {
	return filepath.Join(configDir(), "config")
}

func Load() Config {
	cfg := DefaultConfig()

	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg
	}

	_ = json.Unmarshal(data, &cfg)

	if cfg.CacheTTL <= 0 {
		cfg.CacheTTL = 600
	}
	if cfg.MaxJobs <= 0 {
		cfg.MaxJobs = 5
	}

	return cfg
}

func Init() (string, error) {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	path := configPath()
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	cfg := DefaultConfig()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}

	return path, nil
}

func Path() string {
	return configPath()
}

func CheckMigration() string {
	old := oldConfigPath()
	if _, err := os.Stat(old); err == nil {
		if _, err := os.Stat(configPath()); os.IsNotExist(err) {
			return fmt.Sprintf("Found old bash-style config at %s\nPlease migrate to JSON format at %s\nRun 'ghtools config' to create the new config.", old, configPath())
		}
	}
	return ""
}
