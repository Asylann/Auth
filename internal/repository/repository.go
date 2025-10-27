package repository

import (
	"github.com/Asylann/Auth/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
	"time"
)

type Repository struct {
	Pool *pgxpool.Pool
}

func NewRepository(cfg config.Config) (Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.DataConnection)
	if err != nil {
		return Repository{}, err
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 10 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	dbPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return Repository{}, err
	}

	if err = dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return Repository{}, err
	}

	return Repository{Pool: dbPool}, err
}
