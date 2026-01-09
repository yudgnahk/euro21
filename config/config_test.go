package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Check proxy defaults
	if cfg.Proxy.Enabled {
		t.Error("Default proxy should be disabled")
	}
	if cfg.Proxy.Type != "socks5" {
		t.Errorf("Default proxy type = %v, want socks5", cfg.Proxy.Type)
	}

	// Check cache defaults
	if !cfg.Cache.Enabled {
		t.Error("Default cache should be enabled")
	}
	if cfg.Cache.TTL != 300 {
		t.Errorf("Default cache TTL = %v, want 300", cfg.Cache.TTL)
	}

	// Check tournament defaults
	if cfg.Tournament.ID != 1 {
		t.Errorf("Default tournament ID = %v, want 1", cfg.Tournament.ID)
	}
	if cfg.Tournament.SeasonID != 61644 {
		t.Errorf("Default season ID = %v, want 61644", cfg.Tournament.SeasonID)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	// Change to a temp directory where no config exists
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	cfg, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() should not error when file not found, got: %v", err)
	}

	// Should return default config
	if cfg == nil {
		t.Fatal("LoadConfig() returned nil")
	}
	if cfg.Tournament.ID != 1 {
		t.Error("LoadConfig() should return default config when file not found")
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "euro21.config.json")

	// Create a test config
	testConfig := Config{
		Proxy: ProxyConfig{
			Enabled: true,
			Type:    "tor",
			Address: "localhost:9999",
		},
		Cache: CacheConfig{
			Enabled:  false,
			CacheDir: "/tmp/test",
			TTL:      600,
		},
		Tournament: TournamentConfig{
			ID:       999,
			SeasonID: 888,
			Name:     "Test Tournament",
		},
	}

	// Write config to file
	data, _ := json.MarshalIndent(testConfig, "", "  ")
	err := os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Change to temp directory
	oldDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Load config
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Verify loaded values
	if cfg.Proxy.Enabled != true {
		t.Error("Proxy.Enabled not loaded correctly")
	}
	if cfg.Proxy.Type != "tor" {
		t.Errorf("Proxy.Type = %v, want tor", cfg.Proxy.Type)
	}
	if cfg.Cache.Enabled != false {
		t.Error("Cache.Enabled not loaded correctly")
	}
	if cfg.Tournament.ID != 999 {
		t.Errorf("Tournament.ID = %v, want 999", cfg.Tournament.ID)
	}
	if cfg.Tournament.SeasonID != 888 {
		t.Errorf("Tournament.SeasonID = %v, want 888", cfg.Tournament.SeasonID)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "euro21.config.json")

	// Write invalid JSON
	err := os.WriteFile(configPath, []byte("invalid json{{{"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	oldDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Should fall back to default config on invalid JSON
	cfg, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() should not error on invalid JSON, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig() should return default config on invalid JSON")
	}
}

func TestConfig_Save(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "saved_config.json")

	cfg := &Config{
		Proxy: ProxyConfig{
			Enabled: true,
			Type:    "socks5",
			Address: "localhost:1080",
		},
		Cache: CacheConfig{
			Enabled:  true,
			CacheDir: "/tmp/cache",
			TTL:      300,
		},
		Tournament: TournamentConfig{
			ID:       123,
			SeasonID: 456,
			Name:     "Test",
		},
	}

	err := cfg.Save(configPath)
	if err != nil {
		t.Fatalf("Config.Save() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Load and verify content
	data, _ := os.ReadFile(configPath)
	var loaded Config
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved config: %v", err)
	}

	if loaded.Tournament.ID != 123 {
		t.Errorf("Saved config Tournament.ID = %v, want 123", loaded.Tournament.ID)
	}
}

func TestConfig_ToProxyConfig(t *testing.T) {
	cfg := &Config{
		Proxy: ProxyConfig{
			Enabled:        true,
			Type:           "tor",
			Address:        "localhost:9050",
			TorControlAddr: "localhost:9051",
		},
	}

	proxyConfig := cfg.ToProxyConfig()

	if proxyConfig.Enabled != true {
		t.Error("ToProxyConfig() Enabled not converted correctly")
	}
	if proxyConfig.Type != "tor" {
		t.Errorf("ToProxyConfig() Type = %v, want tor", proxyConfig.Type)
	}
	if proxyConfig.Address != "localhost:9050" {
		t.Errorf("ToProxyConfig() Address = %v, want localhost:9050", proxyConfig.Address)
	}
	if proxyConfig.TorControlAddr != "localhost:9051" {
		t.Errorf("ToProxyConfig() TorControlAddr = %v, want localhost:9051", proxyConfig.TorControlAddr)
	}
}

func TestConfig_ToCacheConfig(t *testing.T) {
	cfg := &Config{
		Cache: CacheConfig{
			Enabled:  true,
			CacheDir: "/tmp/test-cache",
			TTL:      600,
		},
	}

	cacheConfig := cfg.ToCacheConfig()

	if cacheConfig.Enabled != true {
		t.Error("ToCacheConfig() Enabled not converted correctly")
	}
	if cacheConfig.CacheDir != "/tmp/test-cache" {
		t.Errorf("ToCacheConfig() CacheDir = %v, want /tmp/test-cache", cacheConfig.CacheDir)
	}
	if cacheConfig.TTL != 600*time.Second {
		t.Errorf("ToCacheConfig() TTL = %v, want 600s", cacheConfig.TTL)
	}
}
