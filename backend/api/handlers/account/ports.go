package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/dto"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
	"github.com/koo-arch/adjusta-backend/internal/usecase/account/calendarsetting"
)

type ProfileUsecase interface {
	FetchGoogleProfile(ctx context.Context, userID uuid.UUID) (*usecaseAccount.GoogleProfile, error)
}

type CalendarSettingsUsecase interface {
	ListCalendarSettings(ctx context.Context, userID uuid.UUID, email string) ([]calendarsetting.CalendarSettingOutput, error)
	UpdateCalendarSetting(ctx context.Context, userID, userCalendarID uuid.UUID, email string, req calendarsetting.CalendarSettingUpdateRequest) (*calendarsetting.CalendarSettingOutput, error)
	GetCandidateSyncSetting(ctx context.Context, userID uuid.UUID) (*calendarsetting.CandidateSyncSettingOutput, error)
	SetCandidateSyncSetting(ctx context.Context, userID uuid.UUID, email string, enabled bool) (*calendarsetting.CandidateSyncSettingOutput, error)
}

func toCalendarSettingUpdateRequest(req *dto.CalendarSettingUpdate) calendarsetting.CalendarSettingUpdateRequest {
	if req == nil {
		return calendarsetting.CalendarSettingUpdateRequest{}
	}

	return calendarsetting.CalendarSettingUpdateRequest{
		Role:              req.Role,
		IsVisible:         req.IsVisible,
		SyncProposedDates: req.SyncProposedDates,
	}
}
