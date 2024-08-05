package oracle

import (
	"context"
	"database/sql"
)

type Handler func(ctx context.Context) error

type Client interface {
	DB() DB
	Close() error
}

type TxManager interface {
	ReadSerializable(ctx context.Context, f Handler) error
}

type SQLExecer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type Transactor interface {
	BeginTx(ctx context.Context, txOptions *sql.TxOptions) (*sql.Tx, error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close() error
}
