package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/dto"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
)

type ProfileUsecase interface {
	FetchGoogleProfile(ctx context.Context, userID uuid.UUID) (*usecaseAccount.GoogleProfile, error)
}

type CalendarSettingsUsecase interface {
	ListCalendarSettings(ctx context.Context, userID uuid.UUID, email string) ([]usecaseAccount.CalendarSettingOutput, error)
	UpdateCalendarSetting(ctx context.Context, userID, userCalendarID uuid.UUID, email string, req usecaseAccount.CalendarSettingUpdateRequest) (*usecaseAccount.CalendarSettingOutput, error)
}

func toCalendarSettingUpdateRequest(req *dto.CalendarSettingUpdate) usecaseAccount.CalendarSettingUpdateRequest {
	if req == nil {
		return usecaseAccount.CalendarSettingUpdateRequest{}
	}

	return usecaseAccount.CalendarSettingUpdateRequest{
		Role:              req.Role,
		IsVisible:         req.IsVisible,
		SyncProposedDates: req.SyncProposedDates,
	}
}
