package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DB_URL      string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func Read() (Config, error) {
	var config Config
	filePath, err := configFilePath()

	if err != nil {
		return config, err
	}

	jsonData, err := os.ReadFile(filePath)

	if err := json.Unmarshal(jsonData, &config); err != nil {
		return config, err
	}

	return config, nil
}

func (c Config) SetUser(username string) error {

	c.CurrentUser = username
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	configPath, err := configFilePath()

	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

func configFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(homeDir, configFileName)

	return filePath, nil
}
