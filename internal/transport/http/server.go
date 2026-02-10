package http

import (
	"context"
	"net/http"
	"restapi/internal/config"
	log "restapi/internal/logger"
	"restapi/internal/transport/http/router"
)

type Server struct {
	srv *http.Server
}

func NewServer(cfg *config.Config) *Server {
	handler := router.NewRouter()
	h := cfg.App.HTTP

	return &Server{
		srv: &http.Server{
			Addr:              h.Addr,
			Handler:           handler,
			ReadHeaderTimeout: h.ReadHeaderTimeout,
			ReadTimeout:       h.ReadTimeout,
			WriteTimeout:      h.WriteTimeout,
			IdleTimeout:       h.IdleTimeout,
		},
	}
}

// Run запускает HTTP-сервер (блокирующий вызов). При Shutdown возвращает http.ErrServerClosed.
func (s *Server) Run() error {
	log.Info("Server started", "addr", s.srv.Addr)
	return s.srv.ListenAndServe()
}

// Shutdown останавливает сервер с учётом таймаута. Дожидается завершения активных запросов.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
