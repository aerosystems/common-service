package grpcclient

import (
	"crypto/tls"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

const (
	defaultKeepaliveTime       = 30 * time.Second
	defaultKeepaliveTimeout    = 10 * time.Second
	defaultMaxRetries          = 3
	defaultRetryBackoffInitial = 500 * time.Millisecond
)

func NewGRPCConn(address string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                defaultKeepaliveTime,
			Timeout:             defaultKeepaliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithMax(defaultMaxRetries),
			grpc_retry.WithBackoff(grpc_retry.BackoffExponential(defaultRetryBackoffInitial)),
			grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted, codes.DeadlineExceeded),
		)),
	}

	if address[len(address)-4:] == ":443" {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	return grpc.NewClient(address, opts...)
}
