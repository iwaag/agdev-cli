package socket

import "context"

type Client interface {
	Connect(ctx context.Context) error
	Close() error
}
