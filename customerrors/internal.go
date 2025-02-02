package customerrors

import "google.golang.org/grpc/codes"

type InternalError struct {
	Message  string
	HttpCode int
	GrpcCode codes.Code
}

func (e InternalError) Error() string {
	return e.Message
}
