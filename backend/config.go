package main

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

type config struct {
	AWS struct {
		Bucket          string `yaml:"bucket"`
		Folder          string `yaml:"folder"`
		AccessKeyID     string `yaml:"access_key_id"`
		SecretAccessKey string `yaml:"secret_access_key"`
		Region          string `yaml:"region"`
	} `yaml:"aws"`
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
