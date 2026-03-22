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
	ReadSession(ctx context.Context) (Session, error)
	WriteSession(ctx context.Context, session Session) error
}

type FileStore struct {
	path string
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

func (s *FileStore) ReadSession(ctx context.Context) (Session, error) {
	select {
	case <-ctx.Done():
		return Session{}, ctx.Err()
	default:
	}

	raw, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Session{}, nil
		}

		return Session{}, fmt.Errorf("read saved login token: %w", err)
	}

	var session Session
	if err := json.Unmarshal(raw, &session); err != nil {
		return Session{}, fmt.Errorf("parse saved login token: %w", err)
	}

	session.AccessToken = strings.TrimSpace(session.AccessToken)
	session.RefreshToken = strings.TrimSpace(session.RefreshToken)
	session.UserID = strings.TrimSpace(session.UserID)

	return session, nil
}

func (s *FileStore) WriteSession(ctx context.Context, session Session) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	session.AccessToken = strings.TrimSpace(session.AccessToken)
	session.RefreshToken = strings.TrimSpace(session.RefreshToken)
	session.UserID = strings.TrimSpace(session.UserID)
	if session.AccessToken == "" {
		return errors.New("access token is required")
	}

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create auth directory: %w", err)
	}

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("encode auth session: %w", err)
	}
	data = append(data, '\n')

	tempFile, err := os.CreateTemp(dir, "auth-*.json")
	if err != nil {
		return fmt.Errorf("create auth temp file: %w", err)
	}

	tempPath := tempFile.Name()
	cleanup := func() {
		_ = os.Remove(tempPath)
	}

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		cleanup()
		return fmt.Errorf("write auth temp file: %w", err)
	}
	if err := tempFile.Chmod(0o600); err != nil {
		tempFile.Close()
		cleanup()
		return fmt.Errorf("set auth temp file permissions: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		cleanup()
		return fmt.Errorf("close auth temp file: %w", err)
	}
	if err := os.Rename(tempPath, s.path); err != nil {
		cleanup()
		return fmt.Errorf("replace auth file: %w", err)
	}

	return nil
}
