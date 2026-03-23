package errors

type Code string

const (
	// Common
	CodeUnknown      Code = "UNKNOWN"
	CodeInvalidInput Code = "INVALID_INPUT"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeForbidden    Code = "FORBIDDEN"
	CodeNotFound     Code = "NOT_FOUND"
	CodeConflict     Code = "CONFLICT"
	CodeSystem       Code = "SYSTEM"
	CodeBusiness     Code = "BUSINESS"

	// Business
	CodeAlreadyBooked Code = "ALREADY_BOOKED"
	CodeEventFull     Code = "EVENT_FULL"
	CodeEventNotFound Code = "EVENT_NOT_FOUND"
	CodeBookingClosed Code = "BOOKING_CLOSED"

	// external
	CodeDependencyUnavailable Code = "DEPENDENCY_UNAVAILABLE"
)
