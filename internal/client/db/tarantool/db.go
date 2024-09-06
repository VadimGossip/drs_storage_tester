package tarantool

import (
	"github.com/tarantool/go-tarantool/v2"
)

type Client interface {
	DB() DB
	Close() error
}

type Execer interface {
	Do(req tarantool.Request) *tarantool.Future
}

type DB interface {
	Execer
	Close() error
}
