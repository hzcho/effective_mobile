package server

import (
	"context"
	"net/http"
	"song_lib/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Server, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTime,
			WriteTimeout: cfg.WriteTime,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
