package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	DBUrl          string `json:"db_url"`
	CurentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home dir: %w", err)
	}
	return fmt.Sprintf("%v/%v", home_dir, configFileName), nil
}

func Read() (config Config, err error) {
	configPath, err := getConfigFilePath()
	jsonFile, err := os.Open(configPath)
	if err != nil {
		// Return a wrapped error that clearly explains what failed
		return Config{}, fmt.Errorf("failed to open config file %s: %w", configPath, err)
	}
	defer jsonFile.Close()

	// 2. Read the file contents
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}
	if err := json.Unmarshal(byteValue, &config); err != nil {
		// Return a specific error for unmarshaling failure
		return Config{}, fmt.Errorf("failed to unmarshal JSON from %s: %w", configPath, err)
	}
	return
}
func (c Config) SetUser(username string) (err error) {
	c.CurentUserName = username
	config, _ := json.Marshal(c)
	save_path, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("could not get config file: %w", err)
	}
	err = os.WriteFile(save_path, config, 0644)
	return nil
}
