package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
)

type PgxMockAdapter struct {
	mock pgxmock.PgxPoolIface
}

func NewPgxMockAdapter(mock pgxmock.PgxPoolIface) *PgxMockAdapter {
	return &PgxMockAdapter{mock: mock}
}

func (p *PgxMockAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return p.mock.QueryRow(ctx, sql, args...)
}

func (p *PgxMockAdapter) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return p.mock.Query(ctx, sql, args...)
}

func (p *PgxMockAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	ct, err := p.mock.Exec(ctx, sql, args...)
	return pgconn.CommandTag(ct), err
}

func (p *PgxMockAdapter) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.mock.Begin(ctx)
}