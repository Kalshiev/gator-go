package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Db_URL            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	c.Current_user_name = user
	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return configDir + "/" + configFileName, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY

	permissions := os.FileMode(0644)

	file, err := os.OpenFile(path, flags, permissions)
	if err != nil {
		return err
	}

	defer file.Close()

	enc := json.NewEncoder(file)
	err = enc.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
