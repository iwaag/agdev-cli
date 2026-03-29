package openapi

import "testing"

func TestResolveOutputPathUsesSanitizedTitle(t *testing.T) {
	t.Parallel()

	payload := map[string]any{
		"info": map[string]any{
			"title": "Example Backend API",
		},
	}

	got, err := ResolveOutputPath("", payload)
	if err != nil {
		t.Fatalf("ResolveOutputPath returned error: %v", err)
	}

	want := ".hints/Example_Backend_API.json"
	if got != want {
		t.Fatalf("ResolveOutputPath = %q, want %q", got, want)
	}
}

func TestFilterOperationsByTags(t *testing.T) {
	t.Parallel()

	payload := map[string]any{
		"paths": map[string]any{
			"/missions": map[string]any{
				"get": map[string]any{
					"tags": []any{"mission"},
				},
				"post": map[string]any{
					"tags": []any{"admin"},
				},
			},
			"/auth/login": map[string]any{
				"post": map[string]any{
					"tags": []any{"auth"},
				},
			},
		},
		"tags": []any{
			map[string]any{"name": "mission"},
			map[string]any{"name": "auth"},
			map[string]any{"name": "admin"},
		},
	}

	FilterOperationsByTags(payload, []string{"mission"})

	paths := payload["paths"].(map[string]any)
	missions := paths["/missions"].(map[string]any)
	if _, ok := missions["get"]; !ok {
		t.Fatalf("expected GET /missions to remain")
	}
	if _, ok := missions["post"]; ok {
		t.Fatalf("expected POST /missions to be removed")
	}
	if _, ok := paths["/auth/login"]; ok {
		t.Fatalf("expected /auth/login path to be removed")
	}

	tags := payload["tags"].([]any)
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag entry, got %d", len(tags))
	}
	if tags[0].(map[string]any)["name"] != "mission" {
		t.Fatalf("expected remaining tag metadata to be mission, got %#v", tags[0])
	}
}
