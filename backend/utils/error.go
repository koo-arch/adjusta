package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"google.golang.org/api/googleapi"
)

func GetAPIError(err error, fallbackMessage string) error {
	if apiErr, ok := err.(*internalErrors.APIError); ok {
		return apiErr
	}
	return internalErrors.NewAPIError(http.StatusInternalServerError, fallbackMessage)
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
	var gErr *googleapi.Error
	if errors.As(err, &gErr) {
		switch gErr.Code {
		case http.StatusUnauthorized:
			return internalErrors.NewAPIError(http.StatusUnauthorized, "認証エラーが発生しました")
		case http.StatusForbidden:
			return internalErrors.NewAPIError(http.StatusForbidden, "アクセス権限がありません")
		case http.StatusNotFound:
			return internalErrors.NewAPIError(http.StatusNotFound, "リソースが見つかりません")
		default:
			return internalErrors.NewAPIError(http.StatusInternalServerError, "Google APIエラーが発生しました")
		}
	}

	return internalErrors.NewAPIError(http.StatusInternalServerError, "予期せぬエラーが発生しました")
}
