package connector

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPostgresConnector(dsn string, maxPoolSize int32) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.MaxConns = maxPoolSize

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err = dbPool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return dbPool, nil
}
