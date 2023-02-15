package server

import (
	"context"
	"net/http"
	"timeline/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *config.Config, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        handler,
		MaxHeaderBytes: cfg.MaxHeaderBytes << 20, // 1 MB
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
