package errors

import "net/http"

type ValidationError struct {
	StatusCode int
	Message    string
	Details    map[string]string
}

func NewValidationError(details map[string]string) *ValidationError {
	return &ValidationError{
		StatusCode: http.StatusBadRequest,
		Message:    "送信に失敗しました",
		Details:    details,
	}
}

func (e *ValidationError) Error() string {
	return e.Message
}
