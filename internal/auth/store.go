package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Store interface {
	ReadToken(ctx context.Context) (string, error)
}

type FileStore struct {
	path string
}

type fileToken struct {
	Token string `json:"token"`
}

func NewFileStore() (*FileStore, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolve config directory: %w", err)
	}

	return &FileStore{
		path: filepath.Join(configDir, "agdev", "auth.json"),
	}, nil
}

func (s *FileStore) ReadToken(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	raw, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		return "", fmt.Errorf("read saved login token: %w", err)
	}

	var record fileToken
	if err := json.Unmarshal(raw, &record); err != nil {
		return "", fmt.Errorf("parse saved login token: %w", err)
	}

	return strings.TrimSpace(record.Token), nil
}
