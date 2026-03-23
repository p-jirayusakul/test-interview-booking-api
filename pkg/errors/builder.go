package errors

func New(code Code, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

func Wrap(code Code, msg string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Cause:   cause,
	}
}
