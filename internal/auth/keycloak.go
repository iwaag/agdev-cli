package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"agdev/internal/app"
)

const keycloakRequestTimeout = 15 * time.Second

type Refresher interface {
	Refresh(ctx context.Context, refreshToken string) (Session, error)
}

type KeycloakClient struct {
	baseURL     string
	realm       string
	clientID    string
	defaultUser string
	httpClient  *http.Client
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type jwtClaims struct {
	Subject string `json:"sub"`
}

func NewKeycloakClientFromEnv() (*KeycloakClient, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("KEYCLOAK_URL")), "/")
	realm := strings.TrimSpace(os.Getenv("KEYCLOAK_REALM"))
	clientID := strings.TrimSpace(os.Getenv("KEYCLOAK_CLIENT_ID"))
	defaultUser := strings.TrimSpace(os.Getenv("KEYCLOAK_USER_NAME"))

	if baseURL == "" || realm == "" || clientID == "" {
		return nil, app.WithExitCode(app.ExitUsage, fmt.Errorf("KEYCLOAK_URL, KEYCLOAK_REALM, and KEYCLOAK_CLIENT_ID must be set"))
	}

	return &KeycloakClient{
		baseURL:     baseURL,
		realm:       realm,
		clientID:    clientID,
		defaultUser: defaultUser,
		httpClient: &http.Client{
			Timeout: keycloakRequestTimeout,
		},
	}, nil
}

func (c *KeycloakClient) DefaultUser() string {
	return c.defaultUser
}

func (c *KeycloakClient) LoginPassword(ctx context.Context, username, password string) (Session, error) {
	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", c.clientID)
	form.Set("username", strings.TrimSpace(username))
	form.Set("password", password)
	form.Set("scope", "openid")

	return c.exchange(ctx, form)
}

func (c *KeycloakClient) Refresh(ctx context.Context, refreshToken string) (Session, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", c.clientID)
	form.Set("refresh_token", strings.TrimSpace(refreshToken))

	return c.exchange(ctx, form)
}

func (c *KeycloakClient) exchange(ctx context.Context, form url.Values) (Session, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenEndpoint(), strings.NewReader(form.Encode()))
	if err != nil {
		return Session{}, app.WithExitCode(app.ExitInternal, fmt.Errorf("build login request: %w", err))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Session{}, classifyAuthRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Session{}, classifyAuthStatusError(resp)
	}

	var payload tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return Session{}, app.WithExitCode(app.ExitBackend, fmt.Errorf("decode login response: %w", err))
	}

	session := Session{
		AccessToken:  strings.TrimSpace(payload.AccessToken),
		RefreshToken: strings.TrimSpace(payload.RefreshToken),
		ExpiresAt:    time.Now().Add(time.Duration(payload.ExpiresIn) * time.Second),
	}
	session.UserID = extractUserID(payload.IDToken, payload.AccessToken)

	if session.AccessToken == "" {
		return Session{}, app.WithExitCode(app.ExitBackend, fmt.Errorf("login response did not include an access token"))
	}

	return session, nil
}

func (c *KeycloakClient) tokenEndpoint() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.baseURL, url.PathEscape(c.realm))
}

func extractUserID(tokens ...string) string {
	for _, token := range tokens {
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			continue
		}

		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			continue
		}

		var claims jwtClaims
		if err := json.Unmarshal(payload, &claims); err != nil {
			continue
		}
		if claims.Subject != "" {
			return claims.Subject
		}
	}

	return ""
}

func classifyAuthRequestError(err error) error {
	return app.WithExitCode(app.ExitNetwork, fmt.Errorf("authentication request failed: %w", err))
}

func classifyAuthStatusError(resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8192))
	if err != nil {
		return app.WithExitCode(app.ExitBackend, fmt.Errorf("authentication failed with status %d", resp.StatusCode))
	}

	message := strings.TrimSpace(string(body))
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}

	switch resp.StatusCode {
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden:
		return app.WithExitCode(app.ExitAuth, fmt.Errorf("authentication failed: %s", message))
	default:
		return app.WithExitCode(app.ExitBackend, fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, message))
	}
}
