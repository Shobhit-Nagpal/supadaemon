package errors

type baseError struct {
	Message string
}

func newBaseError(msg string) baseError {
	return baseError{
		Message: msg,
	}
}

func (e baseError) Error() string {
	return e.Message
}
