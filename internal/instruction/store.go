package instruction

import (
	"embed"
	"fmt"
	"path"
	"strings"

	"agdev/internal/app"
)

//go:embed files/**
var files embed.FS

type Document struct {
	Scope           string
	Name            string
	Version         string
	ResolvedVersion string
	Body            string
}

func Get(scope string, name string, version string) (Document, error) {
	base := path.Join("files", scope, name)

	resolvedVersion := strings.TrimSpace(version)
	if resolvedVersion == "" {
		resolvedVersion = "latest"
	}

	if resolvedVersion == "latest" {
		aliasPath := path.Join(base, "latest")
		alias, err := files.ReadFile(aliasPath)
		if err != nil {
			return Document{}, app.WithExitCode(app.ExitInternal, fmt.Errorf("read instruction alias %s: %w", aliasPath, err))
		}

		resolvedVersion = strings.TrimSpace(string(alias))
		if resolvedVersion == "" {
			return Document{}, app.WithExitCode(app.ExitInternal, fmt.Errorf("instruction alias %s is empty", aliasPath))
		}
	}

	bodyPath := path.Join(base, resolvedVersion+".md")
	body, err := files.ReadFile(bodyPath)
	if err != nil {
		return Document{}, app.WithExitCode(app.ExitUsage, fmt.Errorf("instruction not found: %s/%s@%s", scope, name, version))
	}

	return Document{
		Scope:           scope,
		Name:            name,
		Version:         versionOrLatest(version),
		ResolvedVersion: resolvedVersion,
		Body:            strings.TrimSpace(string(body)),
	}, nil
}

func versionOrLatest(version string) string {
	trimmed := strings.TrimSpace(version)
	if trimmed == "" {
		return "latest"
	}

	return trimmed
}
