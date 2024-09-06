package tdb

import (
	db "github.com/VadimGossip/tj-drs-storage/internal/client/db/tarantool"

	"github.com/tarantool/go-tarantool/v2"
)

type tdb struct {
	dbc *tarantool.Connection
}

func NewDB(dbc *tarantool.Connection) db.DB {
	return &tdb{
		dbc: dbc,
	}
}

func (t *tdb) Do(req tarantool.Request) *tarantool.Future {
	return t.dbc.Do(req)
}

func (t *tdb) Close() error {
	return t.dbc.Close()
}
