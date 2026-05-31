package errors

func NewBadRequestError(message string) *APIError {
	return NewAPIError(KindBadRequest, message)
}

func NewUnauthorizedError(message string) *APIError {
	return NewAPIError(KindUnauthorized, message)
}

func NewForbiddenError(message string) *APIError {
	return NewAPIError(KindForbidden, message)
}

func NewNotFoundError(message string) *APIError {
	return NewAPIError(KindNotFound, message)
}

func NewInternalError(message string) *APIError {
	return NewAPIError(KindInternal, message)
}

func NewBadGatewayError(message string) *APIError {
	return NewAPIError(KindBadGateway, message)
}

func NewPartialContentError(message string, details map[string][]string) *APIError {
	return NewAPIErrorWithDetails(KindPartial, message, details)
}

func NormalizeAPIError(err error, fallbackMessage string) error {
	if err == nil {
		return nil
	}

	if apiErr, ok := err.(*APIError); ok {
		return apiErr
	}

	return NewInternalError(fallbackMessage)
}
