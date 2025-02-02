package httpserver

import (
	"context"
	"fmt"
	"github.com/go-logrusutil/logrusutil/logctx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	srv  *echo.Echo
	log  *logrus.Logger
	addr string
}

func NewHTTPServer(cfg *Config, log *logrus.Logger, opts ...Option) *Server {
	server := &Server{
		srv:  echo.New(),
		log:  log,
		addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	}

	server.srv.HideBanner = true
	server.srv.HidePort = true

	for _, opt := range opts {
		opt(server)
	}

	server.srv.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))

	server.srv.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().WithContext(logctx.New(c.Request().Context(), logrus.NewEntry(log)))
			return next(c)
		}
	})

	server.srv.Use(middleware.Recover())

	return server
}

func (s *Server) Run() error {
	return errors.WithStack(s.srv.Start(s.addr))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return errors.WithStack(s.srv.Shutdown(ctx))
}
