package routes

type HTTPError struct {
	Code    int    `json:"code"`    // HTTP status code
	Message string `json:"message"` // Error message
}
