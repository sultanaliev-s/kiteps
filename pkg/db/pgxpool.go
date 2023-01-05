package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sultanaliev-s/kiteps/pkg/logging"
)

// NewPGXPool creates a new pgxpool.Pool.
// url is expected to be a valid postgres connection string. See [docs].
// This function will ping the database to ensure the connection is valid.
// If the connection is not valid it will retry 10 times with 2 seconds delay.
//
// [docs]: https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
func NewPGXPool(url string, logger *logging.Logger, ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.Tracer = logger

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		for i := 0; i < 10; i++ {
			pool, err = pgxpool.NewWithConfig(ctx, config)
			if err == nil && pool.Ping(ctx) == nil {
				goto SUCCESS
			}
			logger.Error(
				"failed to connect to database",
				logging.Error("connErr", err),
				logging.Error("pingErr", pool.Ping(ctx)),
				logging.Int("attempt", i+1),
			)
			time.Sleep(2 * time.Second)
		}
		return nil, errors.New("failed to connect to database")
	}

SUCCESS:

	logger.Info("connected to database")
	return pool, nil
}
