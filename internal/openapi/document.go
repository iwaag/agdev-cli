package openapi

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
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"agdev/internal/app"
)

const requestTimeout = 15 * time.Second

var operationKeys = []string{
	"get",
	"post",
	"put",
	"patch",
	"delete",
	"head",
	"options",
	"trace",
}

var fileNameUnsafePattern = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

func ResolveDocumentURL(baseURL string) (string, error) {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed == "" {
		return "", app.WithExitCode(app.ExitUsage, fmt.Errorf("base_url is required"))
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", app.WithExitCode(app.ExitUsage, fmt.Errorf("invalid base_url: %w", err))
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", app.WithExitCode(app.ExitUsage, fmt.Errorf("base_url must include scheme and host"))
	}

	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/openapi.json"
	parsed.RawQuery = ""
	parsed.Fragment = ""

	return parsed.String(), nil
}

func FetchDocument(ctx context.Context, documentURL string) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, documentURL, nil)
	if err != nil {
		return nil, app.WithExitCode(app.ExitInternal, fmt.Errorf("build openapi request: %w", err))
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: requestTimeout}

	resp, err := client.Do(req)
	if err != nil {
		return nil, classifyRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, classifyStatusError(resp)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, app.WithExitCode(app.ExitBackend, fmt.Errorf("decode openapi document: %w", err))
	}

	return payload, nil
}

func FilterOperationsByTags(payload map[string]any, rawTags []string) {
	allowedTags := make(map[string]struct{}, len(rawTags))
	for _, tag := range rawTags {
		if trimmed := strings.TrimSpace(tag); trimmed != "" {
			allowedTags[trimmed] = struct{}{}
		}
	}
	if len(allowedTags) == 0 {
		return
	}

	paths, ok := payload["paths"].(map[string]any)
	if !ok {
		return
	}

	for pathKey, pathValue := range paths {
		pathItem, ok := pathValue.(map[string]any)
		if !ok {
			continue
		}

		hasOperation := false
		for _, operationKey := range operationKeys {
			operationValue, exists := pathItem[operationKey]
			if !exists {
				continue
			}

			operation, ok := operationValue.(map[string]any)
			if !ok {
				delete(pathItem, operationKey)
				continue
			}

			if !operationHasAnyTag(operation, allowedTags) {
				delete(pathItem, operationKey)
				continue
			}

			hasOperation = true
		}

		if !hasOperation {
			delete(paths, pathKey)
		}
	}

	filterTagMetadata(payload, allowedTags)
}

func ResolveOutputPath(rawOut string, payload map[string]any) (string, error) {
	if trimmed := strings.TrimSpace(rawOut); trimmed != "" {
		return trimmed, nil
	}

	title := documentTitle(payload)
	return filepath.Join(".hints", sanitizeFileName(title)+".json"), nil
}

func WriteDocument(outputPath string, payload map[string]any) error {
	dir := filepath.Dir(outputPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return app.WithExitCode(app.ExitInternal, fmt.Errorf("create output directory: %w", err))
		}
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return app.WithExitCode(app.ExitInternal, fmt.Errorf("encode openapi document: %w", err))
	}
	data = append(data, '\n')

	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		return app.WithExitCode(app.ExitInternal, fmt.Errorf("write openapi document: %w", err))
	}

	return nil
}

func classifyRequestError(err error) error {
	var netErr net.Error
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return app.WithExitCode(app.ExitTimeout, fmt.Errorf("openapi request timed out: %w", err))
	case errors.As(err, &netErr) && netErr.Timeout():
		return app.WithExitCode(app.ExitTimeout, fmt.Errorf("openapi request timed out: %w", err))
	default:
		return app.WithExitCode(app.ExitNetwork, fmt.Errorf("openapi request failed: %w", err))
	}
}

func classifyStatusError(resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8192))
	if err != nil {
		return app.WithExitCode(app.ExitBackend, fmt.Errorf("openapi request failed with status %d", resp.StatusCode))
	}

	message := strings.TrimSpace(string(body))
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}

	switch resp.StatusCode {
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		return app.WithExitCode(app.ExitTimeout, fmt.Errorf("openapi request failed with status %d: %s", resp.StatusCode, message))
	default:
		return app.WithExitCode(app.ExitBackend, fmt.Errorf("openapi request failed with status %d: %s", resp.StatusCode, message))
	}
}

func operationHasAnyTag(operation map[string]any, allowedTags map[string]struct{}) bool {
	rawTags, ok := operation["tags"].([]any)
	if !ok || len(rawTags) == 0 {
		return false
	}

	for _, rawTag := range rawTags {
		tag, ok := rawTag.(string)
		if !ok {
			continue
		}
		if _, exists := allowedTags[tag]; exists {
			return true
		}
	}

	return false
}

func filterTagMetadata(payload map[string]any, allowedTags map[string]struct{}) {
	rawTags, ok := payload["tags"].([]any)
	if !ok {
		return
	}

	filtered := make([]any, 0, len(rawTags))
	for _, rawTag := range rawTags {
		tagEntry, ok := rawTag.(map[string]any)
		if !ok {
			continue
		}

		name, _ := tagEntry["name"].(string)
		if _, exists := allowedTags[name]; exists {
			filtered = append(filtered, tagEntry)
		}
	}

	payload["tags"] = filtered
}

func documentTitle(payload map[string]any) string {
	info, ok := payload["info"].(map[string]any)
	if !ok {
		return "openapi"
	}

	title, _ := info["title"].(string)
	title = strings.TrimSpace(title)
	if title == "" {
		return "openapi"
	}

	return title
}

func sanitizeFileName(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return "openapi"
	}

	title = fileNameUnsafePattern.ReplaceAllString(title, "_")
	title = strings.Trim(title, "._-")
	if title == "" {
		return "openapi"
	}

	return title
}
