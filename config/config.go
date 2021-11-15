package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Config stores all configuration of the application.
type AppConfig struct {
	User          string `yaml:"User"`
	APIKey        string `yaml:"APIKey"`
	CurrentBuilds string `yaml:"currentBuilds"`
	CreateDir     string `yaml:"createDir"`
	Port          string
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(filePath string) (*AppConfig, error) {
	var config AppConfig
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", filePath, err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config data %s: %w", data, err)
	}

	config.Port = getEnv("PORT", "8080")

	return &config, nil
}

// Helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
