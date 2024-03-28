package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	RedmineURL   string `json:"redmine_url"`
	RedmineToken string `json:"redmine_token"`
	QueryId      int    `json:"query_id"`
}

var config *Config

func getConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	usr, err := user.Current()
	if err != nil {
		return config, err
	}

	configPath := filepath.Join(usr.HomeDir, ".config", "rqw", "config.json")

	file, err := os.Open(configPath)
	if err != nil {
		return config, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("error decoding JSON: %v", err)
	}

	return config, nil
}
