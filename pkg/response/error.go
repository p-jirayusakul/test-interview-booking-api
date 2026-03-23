package response

type ErrorBlock struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
