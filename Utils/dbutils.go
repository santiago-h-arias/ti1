package utils

import (
	"encoding/json"
	"os"
	config "tinc1/Config"
)

func GetConfiguration() (config.DBConfig, error) {
	config := config.DBConfig{}
	file, err := os.Open("./configuration.json")
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
