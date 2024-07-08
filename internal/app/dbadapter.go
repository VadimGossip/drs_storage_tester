package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

type DBAdapter struct {
	odb *sql.DB
	rdb *redis.Client
}

func NewDBAdapter() *DBAdapter {
	dba := &DBAdapter{}
	return dba
}

func (d *DBAdapter) connectKeyDb() error {
	d.rdb = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		ReadTimeout:  time.Minute * 10,
		WriteTimeout: time.Minute * 10,
	})

	ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	defer cancel()
	err := d.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to ping keyDb host = %s, port = %d:, error = %s", "host", "port", err)
	}
	return nil
}

func (d *DBAdapter) Connect() error {
	return d.connectKeyDb()
}

func (d *DBAdapter) Close() {
	if err := d.odb.Close(); err != nil {
		logrus.Errorf("Error occurred on oracle db connection close: %s", err.Error())
	}
}
