package grpcserver

import (
	"fmt"
	"net"

	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	srv  *grpc.Server
	log  *logrus.Logger
	addr string
}

func NewGRPCServer(
	cfg *Config,
	log *logrus.Logger,
) *Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcctxtags.UnaryServerInterceptor(grpcctxtags.WithFieldExtractor(grpcctxtags.CodeGenRequestFieldExtractor)),
			grpclogrus.UnaryServerInterceptor(logrus.NewEntry(log)),
		),
		grpc.ChainStreamInterceptor(
			grpcctxtags.StreamServerInterceptor(grpcctxtags.WithFieldExtractor(grpcctxtags.CodeGenRequestFieldExtractor)),
			grpclogrus.StreamServerInterceptor(logrus.NewEntry(log)),
		),
	)
	return &Server{
		srv:  grpcServer,
		log:  log,
		addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	}
}

func (s *Server) Run() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.log.Errorf("failed to listen on %s: %v", s.addr, err)
	}
	return errors.WithStack(s.srv.Serve(listen))
}

func (s *Server) Shutdown() {
	s.srv.GracefulStop()
}

func (s *Server) RegisterService(service grpc.ServiceDesc, handler interface{}) {
	s.srv.RegisterService(&service, handler)
}
