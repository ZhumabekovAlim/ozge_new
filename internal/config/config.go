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
	Mobizon struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"mobizon"`
}

// LoadConfig loads the configuration from config.yaml
func LoadConfig() Config {
	var cfg Config

	// Determine config file path
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/config.yaml"
	}

	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// Unmarshal YAML data into config struct
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}

	// Override database URL if provided via environment
	if envURL := os.Getenv("DATABASE_URL"); envURL != "" {
		cfg.Database.URL = envURL
	}

	return cfg
}
