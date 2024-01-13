package configuration

import (
	"encoding/json"
	"os"
)

const CONFIG_FILE_DIR = "./ffcu.json"

type Config struct {
	ProfileDir      string
	UserJsUrl       string
	ZippedChromeUrl string
}

func ReadConfig() (Config, error) {
	config := Config{}

	bytes, err := os.ReadFile(CONFIG_FILE_DIR)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return config, err
	}

	return config, nil
}

func CreateConfig() (Config, error) {
	config := Config{}

	bytes, err := json.Marshal(config)
	if err != nil {
		return config, err
	}

	if err := os.WriteFile(CONFIG_FILE_DIR, bytes, os.ModePerm); err != nil {
		return config, err
	}

	return config, nil
}

func (config *Config) SaveConfig() error {
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err := os.WriteFile(CONFIG_FILE_DIR, bytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (config *Config) IsFilledOut() bool {
	return config.ProfileDir != "" && config.UserJsUrl != "" && config.ZippedChromeUrl != ""
}
