package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"agdev/internal/app"
)

var ErrTokenNotResolved = errors.New("authentication required: specify --token or run agdev login")

type Resolver struct {
	store Store
}

func NewResolver(store Store) Resolver {
	return Resolver{store: store}
}

func DefaultResolver() (Resolver, error) {
	store, err := NewFileStore()
	if err != nil {
		return Resolver{}, app.WithExitCode(app.ExitInternal, err)
	}

	return NewResolver(store), nil
}

func (r Resolver) Resolve(ctx context.Context, explicitToken string) (string, error) {
	if token := strings.TrimSpace(explicitToken); token != "" {
		return token, nil
	}

	if r.store != nil {
		token, err := r.store.ReadToken(ctx)
		if err != nil {
			return "", app.WithExitCode(app.ExitInternal, fmt.Errorf("load authentication token: %w", err))
		}
		if token != "" {
			return token, nil
		}
	}

	return "", app.WithExitCode(app.ExitAuth, ErrTokenNotResolved)
}
