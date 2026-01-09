package sofascore

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Cache provides file-based caching for API responses
type Cache struct {
	enabled  bool
	cacheDir string
	ttl      time.Duration
	mu       sync.RWMutex
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled  bool
	CacheDir string
	TTL      time.Duration
}

// cacheEntry represents a cached response
type cacheEntry struct {
	Data      []byte    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// NewCache creates a new cache instance
func NewCache(config CacheConfig) (*Cache, error) {
	if !config.Enabled {
		return &Cache{enabled: false}, nil
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(config.CacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Cache{
		enabled:  true,
		cacheDir: config.CacheDir,
		ttl:      config.TTL,
	}, nil
}

// getCacheKey generates a cache key from the endpoint
func (c *Cache) getCacheKey(endpoint string) string {
	hash := md5.Sum([]byte(endpoint))
	return hex.EncodeToString(hash[:])
}

// getCacheFile returns the cache file path for an endpoint
func (c *Cache) getCacheFile(endpoint string) string {
	key := c.getCacheKey(endpoint)
	return filepath.Join(c.cacheDir, key+".json")
}

// Get retrieves data from cache if available and not expired
func (c *Cache) Get(endpoint string) ([]byte, bool) {
	if !c.enabled {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	cacheFile := c.getCacheFile(endpoint)

	// Read cache file
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if !os.IsNotExist(err) {
			logrus.Debugf("Failed to read cache file: %v", err)
		}
		return nil, false
	}

	// Unmarshal cache entry
	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		logrus.Debugf("Failed to unmarshal cache entry: %v", err)
		return nil, false
	}

	// Check if expired
	if time.Since(entry.Timestamp) > c.ttl {
		logrus.Debugf("Cache expired for %s", endpoint)
		return nil, false
	}

	logrus.Debugf("Cache hit for %s", endpoint)
	return entry.Data, true
}

// Set stores data in cache
func (c *Cache) Set(endpoint string, data []byte) error {
	if !c.enabled {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{
		Data:      data,
		Timestamp: time.Now(),
	}

	// Marshal cache entry
	entryData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	// Write to cache file
	cacheFile := c.getCacheFile(endpoint)
	if err := os.WriteFile(cacheFile, entryData, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	logrus.Debugf("Cached response for %s", endpoint)
	return nil
}

// Clear removes all cached entries
func (c *Cache) Clear() error {
	if !c.enabled {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove all cache files
	files, err := filepath.Glob(filepath.Join(c.cacheDir, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to list cache files: %w", err)
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			logrus.Warnf("Failed to remove cache file %s: %v", file, err)
		}
	}

	logrus.Info("Cache cleared")
	return nil
}

// ClearExpired removes expired cache entries
func (c *Cache) ClearExpired() error {
	if !c.enabled {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	files, err := filepath.Glob(filepath.Join(c.cacheDir, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to list cache files: %w", err)
	}

	removed := 0
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var entry cacheEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}

		if time.Since(entry.Timestamp) > c.ttl {
			if err := os.Remove(file); err != nil {
				logrus.Warnf("Failed to remove expired cache file %s: %v", file, err)
			} else {
				removed++
			}
		}
	}

	if removed > 0 {
		logrus.Infof("Removed %d expired cache entries", removed)
	}
	return nil
}
