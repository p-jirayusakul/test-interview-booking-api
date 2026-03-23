package errors

import (
	"errors"
	"net/http"
)

func HTTPStatus(err error) int {
	var e *Error
	if !errors.As(err, &e) {
		return http.StatusInternalServerError
	}

	switch e.Code {
	case CodeInvalidInput, CodeBookingClosed:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound, CodeEventNotFound:
		return http.StatusNotFound
	case CodeConflict, CodeAlreadyBooked, CodeEventFull:
		return http.StatusConflict
	case CodeBusiness:
		return http.StatusBadRequest
	case CodeDependencyUnavailable:
		return http.StatusServiceUnavailable
	case CodeUnknown, CodeSystem:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
