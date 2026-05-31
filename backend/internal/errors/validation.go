package errors

type ValidationError struct {
	Message string
	Details map[string]string
}

func NewValidationError(details map[string]string) *ValidationError {
	return &ValidationError{
		Message: "送信に失敗しました",
		Details: details,
	}
}

func (e *ValidationError) Error() string {
	return e.Message
}
