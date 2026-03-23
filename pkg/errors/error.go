package errors

type Error struct {
	Code    Code
	Message string
	Cause   error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) GetCode() string {
	return string(e.Code)
}
