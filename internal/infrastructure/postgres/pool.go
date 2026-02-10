package postgres

import (
	"context"
	"fmt"
	"restapi/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgPool(ctx context.Context, cfg *config.Postgres) (*pgxpool.Pool, error) {
	pc, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse pg config: %w", err)
	}

	pc.MaxConnIdleTime = cfg.MaxConnIdleTime
	pc.MaxConnLifetime = cfg.MaxConnLifetime
	pc.MaxConns = cfg.MaxConns
	pc.MinConns = cfg.MinConns

	pgPool, err := pgxpool.NewWithConfig(ctx, pc)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.HealthTimeout)
	defer cancel()

	if err := pgPool.Ping(pingCtx); err != nil {
		pgPool.Close()
		return nil, fmt.Errorf("ping pg pool: %w", err)
	}

	return pgPool, nil
}
