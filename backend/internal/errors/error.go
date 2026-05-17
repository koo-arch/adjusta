package errors

type APIError struct {
	StatusCode int 			  `json:"-"`	// "-" means that this field will not be marshaled
	Message string 			  `json:"message"`
	Details map[string][]string `json:"details,omitempty"`
}

func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message: message,
	}
}

func NewAPIErrorWithDetails(statusCode int, message string, details map[string][]string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message: message,
		Details: details,
	}
}

func (e *APIError) Error() string {
	return e.Message
}

var (
	InternalErrorMessage = "サーバーでエラーが発生しました"
)