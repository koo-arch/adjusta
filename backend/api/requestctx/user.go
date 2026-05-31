package requestctx

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func UserIDAndEmail(c *gin.Context) (uuid.UUID, string, error) {
	userIDValue, ok := c.Get("user_id")
	if !ok {
		return uuid.Nil, "", internalErrors.NewUnauthorizedError("ユーザー情報が取得できませんでした")
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		return uuid.Nil, "", internalErrors.NewBadRequestError("ユーザーIDの形式が正しくありません")
	}

	emailValue, ok := c.Get("email")
	if !ok {
		return uuid.Nil, "", internalErrors.NewUnauthorizedError("ユーザー情報が取得できませんでした")
	}

	email, ok := emailValue.(string)
	if !ok {
		return uuid.Nil, "", internalErrors.NewBadRequestError("ユーザー情報の形式が正しくありません")
	}

	return userID, email, nil
}
