package httpserver

import "github.com/labstack/echo/v4"

type Option func(*Server)

func WithMiddleware(m echo.MiddlewareFunc) Option {
	return func(s *Server) {
		s.srv.Use(m)
	}
}

func WithRouter(httpMethod string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) Option {
	return func(s *Server) {
		s.srv.Add(httpMethod, path, handler, middleware...)
	}
}

func WithCustomErrorHandler(handler func(err error, c echo.Context)) Option {
	return func(s *Server) {
		s.srv.HTTPErrorHandler = handler
	}
}

func WithValidator(validator echo.Validator) Option {
	return func(s *Server) {
		s.srv.Validator = validator
	}
}
