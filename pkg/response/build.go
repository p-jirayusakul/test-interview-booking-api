package response

import "errors"

import orgerror "github.com/p-jirayusakul/test-interview-booking-api/pkg/errors"

func ErrorResponse(err error, reqID string) *Response[any] {
	if e, ok := errors.AsType[*orgerror.Error](err); ok {
		resp := Response[any]{
			RequestID: reqID,
			Success:   false,
			Error: &ErrorBlock{
				Code:    e.GetCode(),
				Message: e.Message,
			},
		}

		return &resp
	}

	errUnknown := orgerror.New(orgerror.CodeUnknown, err.Error())
	return &Response[any]{
		RequestID: reqID,
		Success:   false,
		Error: &ErrorBlock{
			Code:    errUnknown.GetCode(),
			Message: errUnknown.Message,
		},
	}
}
