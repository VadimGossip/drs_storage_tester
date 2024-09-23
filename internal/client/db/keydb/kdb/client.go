package kdb

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	db "github.com/VadimGossip/drs_storage_tester/internal/client/db/keydb"
	"github.com/VadimGossip/drs_storage_tester/internal/domain"
)

type odbClient struct {
	masterDBC db.DB
}

func New(cfg domain.KeyDbConfig) db.Client {
	dbc := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.Db,
		ReadTimeout:  time.Minute * 10,
		WriteTimeout: time.Minute * 10,
	})
	return &odbClient{
		masterDBC: NewDB(dbc),
	}
}

func (c *odbClient) DB() db.DB {
	return c.masterDBC
}

func (c *odbClient) Close() error {
	if c.masterDBC != nil {
		return c.masterDBC.Close()
	}

	return nil
}
