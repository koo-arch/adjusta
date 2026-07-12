package calendarsetting

import (
	"context"

	"github.com/google/uuid"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) GetCandidateSyncSetting(ctx context.Context, userID uuid.UUID) (*CandidateSyncSettingOutput, error) {
	setting, err := uc.findCandidateSetting(ctx, userID)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return &CandidateSyncSettingOutput{Enabled: false}, nil
	}
	return &CandidateSyncSettingOutput{Enabled: setting.SyncProposedDates, Calendar: setting}, nil
}

func (uc *Usecase) SetCandidateSyncSetting(ctx context.Context, userID uuid.UUID, email string, enabled bool) (*CandidateSyncSettingOutput, error) {
	if enabled {
		if uc.candidateEnabler == nil {
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}
		if err := uc.candidateEnabler.EnableAdjustaCandidateCalendar(ctx, userID, email); err != nil {
			return nil, err
		}
	} else {
		candidate, err := uc.repos.UserCalendar.FindByRole(ctx, userID, value.UserCalendarRoleAdjustaCandidate)
		if err != nil {
			if repoerr.IsNotFound(err) {
				return &CandidateSyncSettingOutput{Enabled: false}, nil
			}
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}
		if _, err := uc.repos.UserCalendar.Update(ctx, userID, candidate.ID, repoUserCalendar.UserCalendarQueryOptions{SyncProposedDates: &enabled}); err != nil {
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}
	}
	return uc.GetCandidateSyncSetting(ctx, userID)
}

func (uc *Usecase) findCandidateSetting(ctx context.Context, userID uuid.UUID) (*CalendarSettingOutput, error) {
	relation, err := uc.repos.UserCalendar.FindByRole(ctx, userID, value.UserCalendarRoleAdjustaCandidate)
	if err != nil {
		if repoerr.IsNotFound(err) {
			return nil, nil
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}
	calendar, err := uc.repos.Calendar.Read(ctx, relation.CalendarID)
	if err != nil {
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}
	setting := toCalendarSettingOutput(relation, calendar)
	return &setting, nil
}
