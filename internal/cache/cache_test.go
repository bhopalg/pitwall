package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileCache(t *testing.T) {
	tmpDir := t.TempDir()
	fc := &FileCache{Dir: tmpDir}

	t.Run("Set and Get successful", func(t *testing.T) {
		key := "test_key"
		value := "some-f1-data"
		ttl := 1 * time.Hour

		err := fc.Set(key, value, ttl)
		if err != nil {
			t.Fatalf("failed to set cache: %v", err)
		}

		var result string
		found, isStale, err := fc.Get(key, &result)

		if err != nil {
			t.Fatalf("error during get: %v", err)
		}
		if !found {
			t.Error("expected to find cached item, but found was false")
		}
		if isStale {
			t.Error("expected item to be fresh, but isStale was true")
		}
		if result != value {
			t.Errorf("expected %s, got %s", value, result)
		}
	})

	t.Run("Get reports stale correctly", func(t *testing.T) {
		key := "stale_key"
		err := fc.Set(key, "old-data", -1*time.Minute)
		if err != nil {
			t.Fatal(err)
		}

		var result string
		found, isStale, _ := fc.Get(key, &result)

		if !found {
			t.Error("expected to find item")
		}
		if !isStale {
			t.Error("expected item to be stale")
		}
	})

	t.Run("Clear removes entries and leaves directory clean", func(t *testing.T) {
		_ = fc.Set("to_clear_1", "data", 1*time.Hour)
		_ = fc.Set("to_clear_2", "data", 1*time.Hour)

		// Add a non-json file to ensure the filter works
		nonJson := filepath.Join(tmpDir, "keep_me.txt")
		_ = os.WriteFile(nonJson, []byte("ignore"), 0644)

		count, err := fc.Clear()
		if err != nil {
			t.Fatalf("clear failed: %v", err)
		}

		if count != 2 {
			t.Errorf("expected to clear 2 files, cleared %d", count)
		}

		var res string
		found, _, _ := fc.Get("to_clear_1", &res)
		if found {
			t.Error("item should have been deleted")
		}

		if _, err := os.Stat(nonJson); os.IsNotExist(err) {
			t.Error("non-json file should not have been deleted")
		}
	})

	t.Run("Info reports correct metadata", func(t *testing.T) {
		_, _ = fc.Clear()

		key := "info_test"
		_ = fc.Set(key, "data", 1*time.Hour)

		infos, absPath, err := fc.Info()
		if err != nil {
			t.Fatal(err)
		}

		if absPath == "" {
			t.Error("expected absolute path to be returned")
		}

		if len(infos) != 1 {
			t.Fatalf("expected 1 info entry, got %d", len(infos))
		}

		if infos[0].Key != key+".json" {
			t.Errorf("expected key %s.json, got %s", key, infos[0].Key)
		}

		if infos[0].Size <= 0 {
			t.Errorf("expected non-zero file size, got %d", infos[0].Size)
		}
	})
}
