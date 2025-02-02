package customerrors

import "google.golang.org/grpc/codes"

type InternalApiError struct {
	Message  string
	HttpCode int
	GrpcCode codes.Code
}

func (e InternalApiError) Error() string {
	return e.Message
}
