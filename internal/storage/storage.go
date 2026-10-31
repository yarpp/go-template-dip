package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var ErrNotFound = errors.New("record not found")

// atomicWrite записывает данные в файл атомарно (через временный файл)
func atomicWrite(path string, v interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// loadJSON загружает данные из файла, если файл отсутствует — оставляет пустым
func loadJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

// Хранилища
type UserStore struct {
	mu    sync.RWMutex
	path  string
	Users []*UserRecord
	next  int
}

type PostStore struct {
	mu    sync.RWMutex
	path  string
	Posts []*PostRecord
	next  int
}

type CommentStore struct {
	mu       sync.RWMutex
	path     string
	Comments []*CommentRecord
	next     int
}
