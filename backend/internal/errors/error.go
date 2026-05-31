package errors

type Kind string

const (
	KindBadRequest   Kind = "bad_request"
	KindUnauthorized Kind = "unauthorized"
	KindForbidden    Kind = "forbidden"
	KindNotFound     Kind = "not_found"
	KindInternal     Kind = "internal"
	KindBadGateway   Kind = "bad_gateway"
	KindPartial      Kind = "partial"
	KindValidation   Kind = "validation"
)

type APIError struct {
	Kind    Kind                `json:"-"`
	Message string              `json:"message"`
	Details map[string][]string `json:"details,omitempty"`
}

func NewAPIError(kind Kind, message string) *APIError {
	return &APIError{
		Kind:    kind,
		Message: message,
	}
}

func NewAPIErrorWithDetails(kind Kind, message string, details map[string][]string) *APIError {
	return &APIError{
		Kind:    kind,
		Message: message,
		Details: details,
	}
}

func (e *APIError) Error() string {
	return e.Message
}

func IsKind(err error, kind Kind) bool {
	apiErr, ok := err.(*APIError)
	return ok && apiErr.Kind == kind
}

var (
	InternalErrorMessage = "サーバーでエラーが発生しました"
)
