package httpserver

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Server struct {
	srv  *echo.Echo
	addr string
}

func NewHTTPServer(cfg *Config, opts ...Option) *Server {
	server := &Server{
		srv:  echo.New(),
		addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	}

	server.srv.HideBanner = true
	server.srv.HidePort = true

	for _, opt := range opts {
		opt(server)
	}
	return server
}

func (s *Server) Run() error {
	return errors.WithStack(s.srv.Start(s.addr))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return errors.WithStack(s.srv.Shutdown(ctx))
}
