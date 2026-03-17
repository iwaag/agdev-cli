package api

import (
	"agdev/internal/config"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	http *resty.Client
}

func NewClient(cfg config.Config) *Client {
	client := resty.New().
		SetTimeout(cfg.RequestTimeout)

	if cfg.APIBaseURL != "" {
		client.SetBaseURL(cfg.APIBaseURL)
	}

	if cfg.AuthToken != "" {
		client.SetAuthToken(cfg.AuthToken)
	}

	// The request/response layer will be added when backend integration starts.
	return &Client{http: client}
}
