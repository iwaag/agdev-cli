package agcode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"agdev/internal/app"
)

const requestTimeout = 15 * time.Second

type Client struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL:   strings.TrimRight(strings.TrimSpace(os.Getenv("AGCODE_API_URL")), "/"),
		authToken: strings.TrimSpace(os.Getenv("AUTH_TOKEN")),
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (c *Client) GetMission(ctx context.Context, missionID string) (map[string]any, error) {
	if c.baseURL == "" {
		return nil, app.WithExitCode(app.ExitUsage, fmt.Errorf("AGCODE_API_URL is not set"))
	}
	if c.authToken == "" {
		return nil, app.WithExitCode(app.ExitAuth, fmt.Errorf("AUTH_TOKEN is not set"))
	}
	if strings.TrimSpace(missionID) == "" {
		return nil, app.WithExitCode(app.ExitUsage, fmt.Errorf("mission_id is required"))
	}

	endpoint, err := url.Parse(c.baseURL + "/mission/get")
	if err != nil {
		return nil, app.WithExitCode(app.ExitUsage, fmt.Errorf("invalid AGCODE_API_URL: %w", err))
	}

	query := endpoint.Query()
	query.Set("mission_id", missionID)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, app.WithExitCode(app.ExitInternal, fmt.Errorf("build mission request: %w", err))
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, classifyRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, classifyStatusError(resp)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, app.WithExitCode(app.ExitBackend, fmt.Errorf("decode mission response: %w", err))
	}

	return payload, nil
}

func classifyRequestError(err error) error {
	var netErr net.Error
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return app.WithExitCode(app.ExitTimeout, fmt.Errorf("mission request timed out: %w", err))
	case errors.As(err, &netErr) && netErr.Timeout():
		return app.WithExitCode(app.ExitTimeout, fmt.Errorf("mission request timed out: %w", err))
	default:
		return app.WithExitCode(app.ExitNetwork, fmt.Errorf("mission request failed: %w", err))
	}
}

func classifyStatusError(resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8192))
	if err != nil {
		return app.WithExitCode(app.ExitBackend, fmt.Errorf("mission request failed with status %d", resp.StatusCode))
	}

	message := strings.TrimSpace(string(body))
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return app.WithExitCode(app.ExitAuth, fmt.Errorf("mission request failed with status %d: %s", resp.StatusCode, message))
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		return app.WithExitCode(app.ExitTimeout, fmt.Errorf("mission request failed with status %d: %s", resp.StatusCode, message))
	default:
		return app.WithExitCode(app.ExitBackend, fmt.Errorf("mission request failed with status %d: %s", resp.StatusCode, message))
	}
}
