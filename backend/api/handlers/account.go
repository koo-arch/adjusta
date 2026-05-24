package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

		accountProfileUsecase := ah.handler.Server.AccountProfileUsecase
		userInfo, err := accountProfileUsecase.FetchGoogleProfile(ctx, userid)
		if err != nil {
			log.Printf("failed to fetch user info for account: %s, %v", email, err)
			utils.HandleAPIError(c, err, "ユーザー情報取得に失敗しました")
			return
		}

		c.JSON(http.StatusOK, userInfo)
	}
}
