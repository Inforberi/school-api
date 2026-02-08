package app

import (
	"context"
	"net/http"
	"time"

	"restapi/internal/config"
	"restapi/internal/log"
	httptransport "restapi/internal/transport/http"
)

type App struct {
	server *httptransport.Server
}

func NewApp(cfg *config.Config) *App {
	return &App{
		server: httptransport.NewServer(cfg),
	}
}

// Run запускает сервер (блокирующий вызов). Для graceful shutdown используй RunUntil.
func (a *App) Run() error {
	return a.server.Run()
}

// Shutdown останавливает сервер (graceful). Передай context с таймаутом.
func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

// RunUntil запускает сервер и блокируется до отмены ctx (например по SIGTERM/SIGINT).
// При отмене ctx выполняет graceful shutdown с заданным shutdownTimeout.
func (a *App) RunUntil(ctx context.Context, shutdownTimeout time.Duration) error {
	runErr := make(chan error, 1)
	go func() { runErr <- a.Run() }()

	select {
	case <-ctx.Done():
		log.Info("shutting down", "reason", ctx.Err().Error())
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := a.Shutdown(shutdownCtx); err != nil {
			log.Error("shutdown error", "err", err)
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
