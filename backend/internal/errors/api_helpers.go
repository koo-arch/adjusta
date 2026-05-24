package errors

import (
	stderrors "errors"
	"net/http"

	"google.golang.org/api/googleapi"
)

func NewBadRequestError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, message)
}

func NewUnauthorizedError(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, message)
}

func NewForbiddenError(message string) *APIError {
	return NewAPIError(http.StatusForbidden, message)
}

func NewNotFoundError(message string) *APIError {
	return NewAPIError(http.StatusNotFound, message)
}

func NewInternalError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, message)
}

func NewBadGatewayError(message string) *APIError {
	return NewAPIError(http.StatusBadGateway, message)
}

func NewPartialContentError(message string, details map[string][]string) *APIError {
	return NewAPIErrorWithDetails(http.StatusPartialContent, message, details)
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

func FromGoogleAPIError(err error) error {
	var gErr *googleapi.Error
	if stderrors.As(err, &gErr) {
		switch gErr.Code {
		case http.StatusUnauthorized:
			return NewUnauthorizedError("認証エラーが発生しました")
		case http.StatusForbidden:
			return NewForbiddenError("アクセス権限がありません")
		case http.StatusNotFound:
			return NewNotFoundError("リソースが見つかりません")
		default:
			return NewInternalError("Google APIエラーが発生しました")
		}
	}

	return NewInternalError("予期せぬエラーが発生しました")
}
