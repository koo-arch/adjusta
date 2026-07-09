package account

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/dto"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
)

type Handler struct {
	accountProfileUsecase   ProfileUsecase
	calendarSettingsUsecase CalendarSettingsUsecase
}

func NewHandler(accountProfileUsecase ProfileUsecase, calendarSettingsUsecase CalendarSettingsUsecase) *Handler {
	return &Handler{
		accountProfileUsecase:   accountProfileUsecase,
		calendarSettingsUsecase: calendarSettingsUsecase,
	}
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

func (ah *Handler) ListCalendarSettingsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, "カレンダー設定取得時にエラーが発生しました")
			return
		}

		settings, err := ah.calendarSettingsUsecase.ListCalendarSettings(ctx, userid, email)
		if err != nil {
			log.Printf("failed to fetch calendar settings for account: %s, %v", email, err)
			respond.Error(c, err, "カレンダー設定の取得に失敗しました")
			return
		}

		respond.OK(c, toCalendarSettingResponses(settings))
	}
}

func (ah *Handler) UpdateCalendarSettingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, "カレンダー設定更新時にエラーが発生しました")
			return
		}

		userCalendarID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			respond.BadRequest(c, "カレンダー設定IDが不正です")
			return
		}

		var req *dto.CalendarSettingUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			respond.BadRequest(c, "リクエストのデータ形式が不正です")
			return
		}
		if !isValidUserCalendarRole(req.Role) {
			respond.BadRequest(c, "カレンダー用途が不正です")
			return
		}

		setting, err := ah.calendarSettingsUsecase.UpdateCalendarSetting(ctx, userid, userCalendarID, email, toCalendarSettingUpdateRequest(req))
		if err != nil {
			log.Printf("failed to update calendar setting for account: %s, %v", email, err)
			respond.Error(c, err, "カレンダー設定の更新に失敗しました")
			return
		}

		respond.OK(c, toCalendarSettingResponse(*setting))
	}
}

func toCalendarSettingResponses(settings []usecaseAccount.CalendarSettingOutput) []dto.CalendarSetting {
	responses := make([]dto.CalendarSetting, 0, len(settings))
	for _, setting := range settings {
		responses = append(responses, toCalendarSettingResponse(setting))
	}
	return responses
}

func toCalendarSettingResponse(setting usecaseAccount.CalendarSettingOutput) dto.CalendarSetting {
	return dto.CalendarSetting{
		ID:                setting.ID,
		CalendarID:        setting.CalendarID,
		GoogleCalendarID:  setting.GoogleCalendarID,
		Summary:           setting.Summary,
		Description:       setting.Description,
		Timezone:          setting.Timezone,
		Role:              setting.Role,
		IsVisible:         setting.IsVisible,
		SyncProposedDates: setting.SyncProposedDates,
	}
}

func isValidUserCalendarRole(role *value.UserCalendarRole) bool {
	if role == nil {
		return true
	}

	switch *role {
	case value.UserCalendarRolePrimary, value.UserCalendarRoleAdjustaCandidate, value.UserCalendarRoleReference:
		return true
	default:
		return false
	}
}
