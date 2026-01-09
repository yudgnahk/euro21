package sofascore

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	tests := []struct {
		name    string
		config  CacheConfig
		wantErr bool
	}{
		{
			name: "disabled cache",
			config: CacheConfig{
				Enabled:  false,
				CacheDir: "",
				TTL:      0,
			},
			wantErr: false,
		},
		{
			name: "valid cache config",
			config: CacheConfig{
				Enabled:  true,
				CacheDir: filepath.Join(os.TempDir(), "euro21-test-cache"),
				TTL:      5 * time.Minute,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := NewCache(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if cache == nil {
				t.Error("NewCache() returned nil cache")
			}

			// Cleanup
			if tt.config.Enabled && tt.config.CacheDir != "" {
				os.RemoveAll(tt.config.CacheDir)
			}
		})
	}
}

func TestCache_SetAndGet(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "euro21-test-cache-set-get")
	defer os.RemoveAll(cacheDir)

	cache, err := NewCache(CacheConfig{
		Enabled:  true,
		CacheDir: cacheDir,
		TTL:      5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Test Set and Get
	key := "/test/endpoint"
	data := []byte(`{"test": "data"}`)

	err = cache.Set(key, data)
	if err != nil {
		t.Errorf("Cache.Set() error = %v", err)
	}

	retrieved, found := cache.Get(key)
	if !found {
		t.Error("Cache.Get() should have found the cached data")
	}

	if string(retrieved) != string(data) {
		t.Errorf("Cache.Get() = %v, want %v", string(retrieved), string(data))
	}
}

func TestCache_GetMiss(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "euro21-test-cache-miss")
	defer os.RemoveAll(cacheDir)

	cache, err := NewCache(CacheConfig{
		Enabled:  true,
		CacheDir: cacheDir,
		TTL:      5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Test cache miss
	_, found := cache.Get("/nonexistent/endpoint")
	if found {
		t.Error("Cache.Get() should not have found nonexistent key")
	}
}

func TestCache_TTLExpiration(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "euro21-test-cache-ttl")
	defer os.RemoveAll(cacheDir)

	cache, err := NewCache(CacheConfig{
		Enabled:  true,
		CacheDir: cacheDir,
		TTL:      100 * time.Millisecond, // Very short TTL for testing
	})
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	key := "/test/ttl"
	data := []byte(`{"test": "ttl"}`)

	err = cache.Set(key, data)
	if err != nil {
		t.Errorf("Cache.Set() error = %v", err)
	}

	// Should be found immediately
	_, found := cache.Get(key)
	if !found {
		t.Error("Cache.Get() should have found the cached data immediately")
	}

	// Wait for TTL to expire
	time.Sleep(150 * time.Millisecond)

	// Should not be found after TTL expiration
	_, found = cache.Get(key)
	if found {
		t.Error("Cache.Get() should not have found expired data")
	}
}

func TestCache_Clear(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "euro21-test-cache-clear")
	defer os.RemoveAll(cacheDir)

	cache, err := NewCache(CacheConfig{
		Enabled:  true,
		CacheDir: cacheDir,
		TTL:      5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Add some data
	cache.Set("/test/1", []byte("data1"))
	cache.Set("/test/2", []byte("data2"))

	// Verify data exists
	_, found := cache.Get("/test/1")
	if !found {
		t.Error("Cache should have data before Clear()")
	}

	// Clear cache
	err = cache.Clear()
	if err != nil {
		t.Errorf("Cache.Clear() error = %v", err)
	}

	// Verify data is gone
	_, found = cache.Get("/test/1")
	if found {
		t.Error("Cache should not have data after Clear()")
	}
}

func TestCache_DisabledCache(t *testing.T) {
	cache, err := NewCache(CacheConfig{
		Enabled: false,
	})
	if err != nil {
		t.Fatalf("Failed to create disabled cache: %v", err)
	}

	// Operations on disabled cache should not error but should not cache
	err = cache.Set("/test", []byte("data"))
	if err != nil {
		t.Errorf("Disabled cache Set() should not error, got: %v", err)
	}

	_, found := cache.Get("/test")
	if found {
		t.Error("Disabled cache Get() should never find data")
	}
}
