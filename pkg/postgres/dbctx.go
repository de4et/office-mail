package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TxClient struct {
	db *sqlx.DB
}

func NewTxClient(db *sqlx.DB) *TxClient {
	return &TxClient{db: db}
}

func (c *TxClient) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return c.db.ExecContext(ctx, query, args...)
}

func (c *TxClient) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	if tx := extractTx(ctx); tx != nil {
		return tx.GetContext(ctx, dest, query, args...)
	}
	return c.db.GetContext(ctx, dest, query, args...)
}

func (c *TxClient) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	if tx := extractTx(ctx); tx != nil {
		return tx.SelectContext(ctx, dest, query, args...)
	}
	return c.db.SelectContext(ctx, dest, query, args...)
}
