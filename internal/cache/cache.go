package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type CacheEntry struct {
	CreatedAt time.Time   `json:"created_at"`
	ExpiresAt time.Time   `json:"expires_at"`
	Data      interface{} `json:"data"`
}

type Cache interface {
	Get(key string, target interface{}) (found bool, isStale bool, err error)
	Set(key string, value interface{}, ttl time.Duration) error
}

type FileCache struct {
	Dir string
}

func (f *FileCache) Get(key string, target interface{}) (bool, bool, error) {
	path := filepath.Join(f.Dir, key+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return false, false, nil
	}

	var entry CacheEntry
	entry.Data = target
	if err := json.Unmarshal(data, &entry); err != nil {
		return false, false, nil
	}

	isStale := time.Now().After(entry.ExpiresAt)
	return true, isStale, nil
}

func (f *FileCache) Set(key string, value interface{}, ttl time.Duration) error {
	_ = os.MkdirAll(f.Dir, 0755)
	entry := CacheEntry{
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ttl),
		Data:      value,
	}
	data, _ := json.Marshal(entry)
	return os.WriteFile(filepath.Join(f.Dir, key+".json"), data, 0644)
}
