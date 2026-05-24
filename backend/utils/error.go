package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func GetAPIError(err error, fallbackMessage string) error {
	return internalErrors.NormalizeAPIError(err, fallbackMessage)
}

func HandleAPIError(c *gin.Context, err error, fallbackMessage string) {
	if apiErr, ok := err.(*internalErrors.APIError); ok {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fallbackMessage})
	}
	c.Abort()
}

func HandleGoogleAPIError(err error) error {
	return internalErrors.FromGoogleAPIError(err)
}
