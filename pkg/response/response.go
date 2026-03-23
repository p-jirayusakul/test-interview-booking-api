package response

type Response[T any] struct {
	RequestID string      `json:"requestId,omitempty"`
	Success   bool        `json:"success"`
	Data      *T          `json:"data,omitempty"`
	Error     *ErrorBlock `json:"error,omitempty"`
}
