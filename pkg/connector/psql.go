package connector

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPostgresConnector(dsn string) (*pgxpool.Pool, error) {
	var dbPool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err = dbPool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return dbPool, nil
}
