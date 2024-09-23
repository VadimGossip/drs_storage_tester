package transaction

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	db "github.com/VadimGossip/drs_storage_tester/internal/client/db/oracle"
	"github.com/VadimGossip/drs_storage_tester/internal/client/db/oracle/odb"
)

type manager struct {
	db db.Transactor
}

func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, txOptions *sql.TxOptions, fn db.Handler) (err error) {
	tx, ok := ctx.Value(odb.TxKey).(*sql.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, txOptions)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}
	ctx = odb.MakeContextTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}

			return
		}

		if nil == err {
			err = tx.Commit()
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) ReadSerializable(ctx context.Context, f db.Handler) error {
	txOpts := &sql.TxOptions{Isolation: sql.LevelSerializable}
	return m.transaction(ctx, txOpts, f)
}
