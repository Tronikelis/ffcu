package ffcu

import (
	"encoding/json"
	"os"
)

type Config struct {
	ProfileDir      string
	UserJsUrl       string
	ZippedChromeUrl string
}

func ReadConfig(dir string) (Config, error) {
	config := Config{}

	bytes, err := os.ReadFile(dir)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return config, err
	}

	return config, nil
}

func CreateConfig(dir string) (Config, error) {
	config := Config{}

	bytes, err := json.Marshal(config)
	if err != nil {
		return config, err
	}

	if err := os.WriteFile(dir, bytes, os.ModePerm); err != nil {
		return config, err
	}

	return config, nil
}

func (config *Config) SaveConfig(dir string) error {
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dir, bytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (config *Config) IsFilledOut() bool {
	// profile dir is required
	if config.ProfileDir == "" {
		return false
	}

	return config.UserJsUrl != "" || config.ZippedChromeUrl != ""
}
