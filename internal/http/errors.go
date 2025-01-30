package http

import "google.golang.org/grpc/codes"

type InternalApiError struct {
	Message  string
	HttpCode int
	GrpcCode codes.Code
}

func (e InternalApiError) Error() string {
	return e.Message
}

type ExternalApiError struct {
	Code     int
	Message  string
	HttpCode int
}

func (e ExternalApiError) Error() string {
	return e.Message
}
