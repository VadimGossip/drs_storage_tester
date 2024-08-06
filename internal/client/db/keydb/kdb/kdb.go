package kdb

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	db "github.com/VadimGossip/tj-drs-storage/internal/client/db/keydb"
)

type kdb struct {
	dbc *redis.Client
}

func NewDB(dbc *redis.Client) db.DB {
	return &kdb{
		dbc: dbc,
	}
}

func (db *kdb) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	return db.dbc.MGet(ctx, keys...)
}

func (db *kdb) Get(ctx context.Context, key string) *redis.StringCmd {
	return db.dbc.Get(ctx, key)
}

func (db *kdb) MGetWithDur(ctx context.Context, keys ...string) ([]any, time.Duration, error) {
	ts := time.Now()
	values, err := db.dbc.MGet(ctx, keys...).Result()
	return values, time.Since(ts), err
}

func (db *kdb) GetWithDur(ctx context.Context, keys string) (string, time.Duration, error) {
	ts := time.Now()
	value, err := db.dbc.Get(ctx, keys).Result()
	return value, time.Since(ts), err
}

func (db *kdb) Ping(ctx context.Context) error {
	return db.dbc.Ping(ctx).Err()
}

func (db *kdb) Close() error {
	return db.dbc.Close()
}
