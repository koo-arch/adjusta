package account

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
)

type Handler struct {
	accountProfileUsecase ProfileUsecase
}

func NewHandler(accountProfileUsecase ProfileUsecase) *Handler {
	return &Handler{accountProfileUsecase: accountProfileUsecase}
}

func (ah *Handler) FetchAccountsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, "ユーザー情報確認時にエラーが発生しました")
			return
		}

		userInfo, err := ah.accountProfileUsecase.FetchGoogleProfile(ctx, userid)
		if err != nil {
			log.Printf("failed to fetch user info for account: %s, %v", email, err)
			respond.Error(c, err, "ユーザー情報取得に失敗しました")
			return
		}

		respond.OK(c, userInfo)
	}
}
