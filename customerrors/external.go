package customerrors

type ExternalError struct {
	Code     int
	Message  string
	HttpCode int
}

func (e ExternalError) Error() string {
	return e.Message
}
