package auth

import (
	"context"

	"agdev/internal/app"
)

type contextKey string

const tokenContextKey contextKey = "auth-token"

func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenContextKey, token)
}

func TokenFromContext(ctx context.Context) (string, error) {
	token, ok := ctx.Value(tokenContextKey).(string)
	if !ok || token == "" {
		return "", app.WithExitCode(app.ExitInternal, ErrTokenNotResolved)
	}

	return token, nil
}
