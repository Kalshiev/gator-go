package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read() Config {
	path, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("path error: %v", err)
		return Config{}
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("dir error: %v", err)
		return Config{}
	}

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Printf("json decode error: %v", err)
		return Config{}
	}

	return config
}
