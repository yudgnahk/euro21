package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yudgnahk/euro21/adapters/sofascore"
)

// Config holds application configuration
type Config struct {
	// Proxy settings
	Proxy ProxyConfig `json:"proxy"`

	// Cache settings
	Cache CacheConfig `json:"cache"`

	// Tournament settings
	Tournament TournamentConfig `json:"tournament"`
}

// ProxyConfig holds proxy configuration
type ProxyConfig struct {
	Enabled        bool   `json:"enabled"`
	Type           string `json:"type"`             // "socks5" or "tor"
	Address        string `json:"address"`          // e.g., "localhost:9050"
	TorControlAddr string `json:"tor_control_addr"` // e.g., "localhost:9051"
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled  bool   `json:"enabled"`
	CacheDir string `json:"cache_dir"`
	TTL      int    `json:"ttl"` // seconds
}

// TournamentConfig holds tournament configuration
type TournamentConfig struct {
	ID       int    `json:"id"`        // SofaScore tournament ID
	SeasonID int    `json:"season_id"` // SofaScore season ID
	Name     string `json:"name"`      // Display name
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".cache", "euro21")

	return &Config{
		Proxy: ProxyConfig{
			Enabled:        false,
			Type:           "socks5",
			Address:        "localhost:9050",
			TorControlAddr: "localhost:9051",
		},
		Cache: CacheConfig{
			Enabled:  true,
			CacheDir: cacheDir,
			TTL:      300, // 5 minutes
		},
		Tournament: TournamentConfig{
			ID:       1,     // Euro 2024 (example, needs verification)
			SeasonID: 61644, // Euro 2024 season (example, needs verification)
			Name:     "UEFA Euro 2024",
		},
	}
}

// LoadConfig loads configuration from file or returns default
func LoadConfig() (*Config, error) {
	// Try to load from config file
	configPaths := []string{
		"./euro21.config.json",
		filepath.Join(os.Getenv("HOME"), ".config", "euro21", "config.json"),
	}

	for _, path := range configPaths {
		if config, err := loadFromFile(path); err == nil {
			return config, nil
		}
	}

	// Return default config if no file found
	return DefaultConfig(), nil
}

// loadFromFile loads configuration from a JSON file
func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// Save saves the configuration to a file
func (c *Config) Save(path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ToProxyConfig converts to SofaScore proxy config
func (c *Config) ToProxyConfig() sofascore.ProxyConfig {
	return sofascore.ProxyConfig{
		Enabled:        c.Proxy.Enabled,
		Type:           c.Proxy.Type,
		Address:        c.Proxy.Address,
		TorControlAddr: c.Proxy.TorControlAddr,
	}
}

// ToCacheConfig converts to SofaScore cache config
func (c *Config) ToCacheConfig() sofascore.CacheConfig {
	return sofascore.CacheConfig{
		Enabled:  c.Cache.Enabled,
		CacheDir: c.Cache.CacheDir,
		TTL:      time.Duration(c.Cache.TTL) * time.Second,
	}
}
