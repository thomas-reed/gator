package config

import (
	"os"
	"encoding/json"
	"path/filepath"
)

const configFilename = ".gatorconfig.json"

type Config struct {
	DbURL string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUsername = username
	return write(*c)
}

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err = json.Unmarshal(configData, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func write(cfg Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	configData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, configData, 0600)
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homedir, configFilename), nil
}