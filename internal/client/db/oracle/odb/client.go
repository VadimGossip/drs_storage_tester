package odb

import (
	"database/sql"

	db "github.com/VadimGossip/drs_storage_tester/internal/client/db/oracle"

	_ "github.com/godror/godror"
)

type odbClient struct {
	masterDBC db.DB
}

func New(dsn string) (db.Client, error) {
	dbc, err := sql.Open("godror", dsn)

	if err != nil {
		return nil, err
	}

	return &odbClient{
		masterDBC: NewDB(dbc),
	}, nil
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
