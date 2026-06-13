package errors

func NewValidationError(details map[string]string) *APIError {
	validationDetails := make(map[string][]string, len(details))
	for field, message := range details {
		if message == "" {
			continue
		}
		validationDetails[field] = []string{message}
	}

	return NewAPIErrorWithDetails(KindValidation, "送信に失敗しました", validationDetails)
}
