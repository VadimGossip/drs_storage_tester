package keydb

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Client interface {
	DB() DB
	Close() error
}

type NamedExecer interface {
	GetWithDur(ctx context.Context, keys string) (string, time.Duration, error)
	MGetWithDur(ctx context.Context, keys ...string) ([]any, time.Duration, error)
}

type Execer interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
}

type RequestExecer interface {
	Execer
	NamedExecer
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	RequestExecer
	Pinger
	Close() error
}
