package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Config structure to hold configuration values
type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	Database struct {
		Driver string `yaml:"driver"`
		URL    string `yaml:"url"`
	} `yaml:"database"`
}

// LoadConfig loads the configuration from config.yaml
func LoadConfig() Config {
	var cfg Config

	// Read config file
	// Read configuration from the relative path so the application can run
	// both locally and inside a Docker container.
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// Unmarshal YAML data into config struct
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}

	return cfg
}
