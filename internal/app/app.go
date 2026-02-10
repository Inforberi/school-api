package app

import (
	"context"
	"net/http"
	"time"

	"restapi/internal/config"
	"restapi/internal/infrastructure/postgres"
	log "restapi/internal/logger"
	httptransport "restapi/internal/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	server *httptransport.Server
	pgPool *pgxpool.Pool
}

func NewApp(cfg *config.Config, ctx context.Context) (*App, error) {
	pgPool, err := postgres.NewPgPool(ctx, &cfg.Postgres)
	if err != nil {
		log.Error("failed to create pg pool", "err", err)
		return nil, err
	}

	return &App{
		server: httptransport.NewServer(cfg),
		pgPool: pgPool,
	}, nil
}

// Run запускает сервер (блокирующий вызов). Для graceful shutdown используй RunWithContext.
func (a *App) Run() error {
	return a.server.Run()
}

// Shutdown останавливает сервер и закрывает ресурсы (graceful). Передай context с таймаутом.
func (a *App) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}
	if a.pgPool != nil {
		a.pgPool.Close()
	}
	return nil
}

// RunWithContext запускает сервер и блокируется до отмены ctx (SIGINT/SIGTERM или отмена).
// При отмене ctx выполняет graceful shutdown с заданным shutdownTimeout.
func (a *App) RunWithContext(ctx context.Context, shutdownTimeout time.Duration) error {
	runErr := make(chan error, 1)
	go func() { runErr <- a.Run() }()

	select {
	case <-ctx.Done():
		log.Info("shutting down", "reason", ctx.Err().Error())
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := a.Shutdown(shutdownCtx); err != nil {
			log.Error("shutdown error", "err", err)
			<-runErr // освобождаем горутину
			return err
		}
		if err := <-runErr; err != nil && err != http.ErrServerClosed {
			log.Error("server exit", "err", err)
			return err
		}
		return nil
	case err := <-runErr:
		if err != nil && err != http.ErrServerClosed {
			log.Error("server failed", "err", err)
			return err
		}
		return nil
	}
}
