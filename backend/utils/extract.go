package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func ExtractUserIDAndEmail(c *gin.Context) (uuid.UUID, string, error) {
	session := sessions.Default(c)
	useridStr, ok := session.Get("userid").(string)
	if !ok {
		return uuid.Nil, "", internalErrors.NewAPIError(http.StatusUnauthorized, "ユーザー情報が取得できませんでした")
	}

	userid, err := uuid.Parse(useridStr)
	if err != nil {
		return uuid.Nil, "", internalErrors.NewAPIError(http.StatusBadRequest, "ユーザーIDの形式が正しくありません")
	}

	email, ok := c.Get("email")
	if !ok {
		return uuid.Nil, "", internalErrors.NewAPIError(http.StatusUnauthorized, "ユーザー情報が取得できませんでした")
	}

	emailStr, ok := email.(string)
	if !ok {
		return uuid.Nil, "", internalErrors.NewAPIError(http.StatusBadRequest, "ユーザー情報の形式が正しくありません")
	}

	return userid, emailStr, nil
}