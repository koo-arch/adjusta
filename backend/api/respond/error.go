package respond

import (
	"net/http"

	"github.com/gin-gonic/gin"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

// Error converts an application or infrastructure error into an HTTP response.
func Error(c *gin.Context, err error, fallbackMessage string) {
	if apiErr, ok := err.(*internalErrors.APIError); ok {
		payload := gin.H{"error": apiErr.Error()}
		if len(apiErr.Details) > 0 {
			payload["details"] = apiErr.Details
		}
		AbortJSON(c, statusCodeForKind(apiErr.Kind), payload)
	} else {
		AbortJSON(c, http.StatusInternalServerError, gin.H{"error": fallbackMessage})
	}
}

func statusCodeForKind(kind internalErrors.Kind) int {
	switch kind {
	case internalErrors.KindBadRequest, internalErrors.KindValidation:
		return http.StatusBadRequest
	case internalErrors.KindUnauthorized:
		return http.StatusUnauthorized
	case internalErrors.KindForbidden:
		return http.StatusForbidden
	case internalErrors.KindNotFound:
		return http.StatusNotFound
	case internalErrors.KindBadGateway:
		return http.StatusBadGateway
	case internalErrors.KindPartial:
		return http.StatusPartialContent
	default:
		return http.StatusInternalServerError
	}
}
