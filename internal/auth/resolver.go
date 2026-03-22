package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"agdev/internal/app"
)

var ErrTokenNotResolved = errors.New("authentication required: specify --token or run agdev login")

type Resolver struct {
	store     Store
	refresher Refresher
}

func NewResolver(store Store, refresher Refresher) Resolver {
	return Resolver{store: store, refresher: refresher}
}

func DefaultResolver() (Resolver, error) {
	store, err := NewFileStore()
	if err != nil {
		return Resolver{}, app.WithExitCode(app.ExitInternal, err)
	}

	client, err := NewKeycloakClientFromEnv()
	if err != nil {
		return NewResolver(store, nil), nil
	}

	return NewResolver(store, client), nil
}

func (r Resolver) Resolve(ctx context.Context, explicitToken string) (string, error) {
	if token := strings.TrimSpace(explicitToken); token != "" {
		return token, nil
	}

	if r.store != nil {
		session, err := r.store.ReadSession(ctx)
		if err != nil {
			return "", app.WithExitCode(app.ExitInternal, fmt.Errorf("load authentication token: %w", err))
		}
		if !session.HasAccessToken() {
			return "", app.WithExitCode(app.ExitAuth, ErrTokenNotResolved)
		}
		if !session.AccessTokenExpired(time.Now()) {
			return session.AccessToken, nil
		}
		if r.refresher != nil && session.RefreshToken != "" {
			refreshed, err := r.refresher.Refresh(ctx, session.RefreshToken)
			if err == nil {
				if refreshed.UserID == "" {
					refreshed.UserID = session.UserID
				}
				if err := r.store.WriteSession(ctx, refreshed); err != nil {
					return "", app.WithExitCode(app.ExitInternal, fmt.Errorf("save refreshed authentication token: %w", err))
				}
				return refreshed.AccessToken, nil
			}
		}
	}

	return "", app.WithExitCode(app.ExitAuth, ErrTokenNotResolved)
}
