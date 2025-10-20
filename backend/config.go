package main

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

type config struct {
	aws struct {
		accessKeyID     string `yaml:"access_key_id"`
		secretAccessKey string `yaml:"secret_access_key"`
		region          string `yaml:"region"`
	}
}

func loadConfig() (*config, error) {
	var cfg config

	file, err := os.Open("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("mở file config thất bại: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("đọc file config thất bại: %w", err)
	}

	return &cfg, nil
}