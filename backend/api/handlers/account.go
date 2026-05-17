package handlers

import (
	"net/http"
	"log"
	
	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	"github.com/koo-arch/adjusta-backend/utils"
)

type AccountHandler struct {
	handler *Handler
}

func NewAccountHandler(handler *Handler) *AccountHandler {
	return &AccountHandler{handler: handler}
}

func (ah *AccountHandler) FetchAccountsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			utils.HandleAPIError(c, err, "ユーザー情報確認時にエラーが発生しました")
			return
		}

		authManager := ah.handler.Server.AuthManager

		token, err := authManager.VerifyOAuthToken(ctx, userid)
		if err != nil {
			log.Printf("failed to verify token for account: %s, %v", email, err)
			utils.HandleAPIError(c, err, "OAuthトークン認証に失敗しました")
			return
		}

		userInfo, err := userinfo.FetchGoogleUserInfo(ctx, token)
		if err != nil {
			log.Printf("failed to fetch user info for account: %s, %v", email, err)
			utils.HandleAPIError(c, err, "ユーザー情報取得に失敗しました")
			return
		}

		c.JSON(http.StatusOK, userInfo)
	}
}
