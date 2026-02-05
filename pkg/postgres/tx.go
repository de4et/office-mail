package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type PostgresqlTransactor struct {
	client *sqlx.DB
}

func NewPostgresqlTransactor(client *sqlx.DB) *PostgresqlTransactor {
	return &PostgresqlTransactor{
		client: client,
	}
}

func (tr *PostgresqlTransactor) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tr.client.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(injectTx(ctx, tx)); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

type txKey struct{}

func injectTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}
