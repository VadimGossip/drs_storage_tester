package app

import (
	"context"
	"fmt"
	"github.com/VadimGossip/tj-drs-storage/internal/domain"
	"github.com/VadimGossip/tj-drs-storage/internal/rate"
	"github.com/go-redis/redis/v8"
	"time"
)

type DBAdapter struct {
	cfg *domain.Config

	kdb      *redis.Client
	rateRepo rate.Repository
}

func NewDBAdapter(cfg *domain.Config) *DBAdapter {
	dba := &DBAdapter{cfg: cfg}
	return dba
}

// fmt.Sprintf(`user=%s password=%s connectString=%s:%d/%s`

func (d *DBAdapter) connectKdb(ctx context.Context) error {
	kdbCfg := d.cfg.TargetDb

	d.kdb = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", kdbCfg.Host, kdbCfg.Port),
		Username:     kdbCfg.Username,
		Password:     kdbCfg.Password,
		DB:           kdbCfg.Db,
		ReadTimeout:  time.Minute * 10,
		WriteTimeout: time.Minute * 10,
	})

	ctxInner, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	err := d.kdb.Ping(ctxInner).Err()
	if err != nil {
		return fmt.Errorf("failed to ping keyDb host = %s, port = %d:, error = %s", kdbCfg.Host, kdbCfg.Port, err)
	}

	d.rateRepo = rate.NewRepository(d.kdb)
	return nil
}

func (d *DBAdapter) Connect(ctx context.Context) error {
	return d.connectKdb(ctx)
}

func (d *DBAdapter) Disconnect() error {
	if err := d.kdb.Close(); err != nil {
		return fmt.Errorf("error occurred on kdb connection close: %s", err.Error())
	}
	return nil
}
