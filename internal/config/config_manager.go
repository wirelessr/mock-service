package config

import (
	"encoding/json"
	"fmt"
	"os"

	"mock-service/internal/models"
)

// ConfigManagerImpl implements the ConfigManager interface
// It handles loading and managing JSON configuration files
type ConfigManagerImpl struct {
	config models.Config
}

// NewConfigManager creates a new instance of ConfigManager
func NewConfigManager() *ConfigManagerImpl {
	return &ConfigManagerImpl{
		config: models.Config{Rules: []models.MockRule{}},
	}
}

// LoadConfig loads configuration from the specified file path
// Returns error if file cannot be read or JSON is invalid
func (cm *ConfigManagerImpl) LoadConfig(filePath string) error {
	// Read the configuration file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	// Parse JSON configuration
	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse JSON config file %s: %w", filePath, err)
	}

	// Store the loaded configuration
	cm.config = config
	return nil
}

// GetConfig returns the current list of mock rules
func (cm *ConfigManagerImpl) GetConfig() []models.MockRule {
	return cm.config.Rules
}
