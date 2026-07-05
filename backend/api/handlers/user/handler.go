package user

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
)

type Handler struct {
	profileUsecase ProfileUsecase
}

func NewHandler(profileUsecase ProfileUsecase) *Handler {
	return &Handler{profileUsecase: profileUsecase}
}

func (uh *Handler) GetCurrentUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			log.Printf("failed to extract user info for account: %s, %v", email, err)
			respond.Error(c, err, "ユーザー情報確認時にエラーが発生しました")
			return
		}

		ctx := c.Request.Context()

		userInfo, err := uh.profileUsecase.FetchGoogleProfile(ctx, userid)
		if err != nil {
			log.Printf("failed to fetch user info for account: %s, %v", email, err)
			respond.Error(c, err, "ユーザー情報取得に失敗しました")
			return
		}

		respond.OK(c, userInfo)
	}
}
