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

type InfoEntry struct {
	Key       string
	CreatedAt time.Time
	ExpiresAt time.Time
	IsStale   bool
	Size      int64
}

type Cache interface {
	Get(key string, target interface{}) (found bool, isStale bool, err error)
	Set(key string, value interface{}, ttl time.Duration) error
	Clear() (int, error)
	Info() ([]InfoEntry, string, error)
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

func (f *FileCache) Clear() (int, error) {
	files, err := os.ReadDir(f.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			err := os.Remove(filepath.Join(f.Dir, file.Name()))
			if err == nil {
				count++
			}
		}
	}
	return count, nil
}

func (f *FileCache) Info() ([]InfoEntry, string, error) {
	absPath, _ := filepath.Abs(f.Dir)
	files, err := os.ReadDir(f.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, absPath, nil
		}
		return nil, absPath, err
	}

	var infos []InfoEntry
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		path := filepath.Join(f.Dir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var entry CacheEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}

		fInfo, _ := file.Info()
		infos = append(infos, InfoEntry{
			Key:       file.Name(),
			CreatedAt: entry.CreatedAt,
			ExpiresAt: entry.ExpiresAt,
			IsStale:   time.Now().After(entry.ExpiresAt),
			Size:      fInfo.Size(),
		})
	}

	return infos, absPath, nil
}
