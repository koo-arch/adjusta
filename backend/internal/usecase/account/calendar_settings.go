package account

import (
	"context"
	"log"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type CalendarSettingsRepositories struct {
	Calendar     repoCalendar.CalendarRepository
	UserCalendar repoUserCalendar.UserCalendarRepository
}

type CalendarSettingsTransaction interface {
	DoCalendarSettings(ctx context.Context, fn func(repos CalendarSettingsRepositories) error) error
}

type CalendarResyncer interface {
	ResyncGoogleCalendars(ctx context.Context, userID uuid.UUID, email string) error
}

type CalendarResyncerFunc func(ctx context.Context, userID uuid.UUID, email string) error

func (f CalendarResyncerFunc) ResyncGoogleCalendars(ctx context.Context, userID uuid.UUID, email string) error {
	return f(ctx, userID, email)
}

type CalendarSettingsUsecase struct {
	repos    CalendarSettingsRepositories
	tx       CalendarSettingsTransaction
	resyncer CalendarResyncer
}

type CalendarSettingOutput struct {
	ID                uuid.UUID
	CalendarID        uuid.UUID
	GoogleCalendarID  string
	Summary           string
	Description       *string
	Timezone          *string
	Role              value.UserCalendarRole
	IsVisible         bool
	SyncProposedDates bool
}

type CalendarSettingUpdateRequest struct {
	Role              *value.UserCalendarRole
	IsVisible         *bool
	SyncProposedDates *bool
}

func NewCalendarSettingsUsecase(repos CalendarSettingsRepositories, tx CalendarSettingsTransaction, resyncer CalendarResyncer) *CalendarSettingsUsecase {
	return &CalendarSettingsUsecase{
		repos:    repos,
		tx:       tx,
		resyncer: resyncer,
	}
}

func (uc *CalendarSettingsUsecase) ListCalendarSettings(ctx context.Context, userID uuid.UUID, email string) ([]CalendarSettingOutput, error) {
	settings, err := listCalendarSettings(ctx, uc.repos, userID)
	if err != nil {
		log.Printf("failed to list calendar settings for account: %s, error: %v", email, err)
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}
	return settings, nil
}

func (uc *CalendarSettingsUsecase) UpdateCalendarSetting(ctx context.Context, userID, userCalendarID uuid.UUID, email string, req CalendarSettingUpdateRequest) (*CalendarSettingOutput, error) {
	var updated *CalendarSettingOutput
	var needsResync bool

	err := uc.tx.DoCalendarSettings(ctx, func(repos CalendarSettingsRepositories) error {
		current, err := findUserCalendarByID(ctx, repos, userID, userCalendarID)
		if err != nil {
			return err
		}

		role := current.Role
		if req.Role != nil {
			role = *req.Role
		}

		syncProposedDates := current.SyncProposedDates
		if req.SyncProposedDates != nil {
			syncProposedDates = *req.SyncProposedDates
		}
		if syncProposedDates && role != value.UserCalendarRoleAdjustaCandidate {
			return internalErrors.NewValidationError(map[string]string{
				"sync_proposed_dates": "候補予定同期は Adjusta 専用カレンダーでのみ有効にできます",
			})
		}
		needsResync = !current.SyncProposedDates && syncProposedDates

		if req.Role != nil && *req.Role == value.UserCalendarRolePrimary {
			if err := demoteExistingPrimary(ctx, repos, userID, current.ID); err != nil {
				return err
			}
		}

		opt := repoUserCalendar.UserCalendarQueryOptions{
			Role:              req.Role,
			IsVisible:         req.IsVisible,
			SyncProposedDates: req.SyncProposedDates,
		}
		saved, err := repos.UserCalendar.Update(ctx, userID, userCalendarID, opt)
		if err != nil {
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("カレンダー設定が見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		calendar, err := repos.Calendar.Read(ctx, saved.CalendarID)
		if err != nil {
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		result := toCalendarSettingOutput(saved, calendar)
		updated = &result
		return nil
	})
	if err != nil {
		log.Printf("failed to update calendar setting for account: %s, error: %v", email, err)
		return nil, err
	}

	// 候補予定同期を OFF→ON にした時点で Adjusta 専用カレンダーを作成/再作成する(requirements 5.7.2)。
	// 同期は更新後の設定値を読むため、トランザクション確定後に実行する。
	// 失敗しても次回 /api/calendar アクセス時の同期ミドルウェアで再作成されるため、設定更新自体は成功として返す。
	if needsResync && uc.resyncer != nil {
		if err := uc.resyncer.ResyncGoogleCalendars(ctx, userID, email); err != nil {
			log.Printf("failed to resync google calendars for account: %s, error: %v", email, err)
		}
	}

	return updated, nil
}

func listCalendarSettings(ctx context.Context, repos CalendarSettingsRepositories, userID uuid.UUID) ([]CalendarSettingOutput, error) {
	userCalendars, err := repos.UserCalendar.FilterByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	calendarIDs := make([]uuid.UUID, 0, len(userCalendars))
	for _, userCalendar := range userCalendars {
		calendarIDs = append(calendarIDs, userCalendar.CalendarID)
	}

	calendars, err := repos.Calendar.FilterByIDs(ctx, calendarIDs)
	if err != nil {
		return nil, err
	}
	calendarByID := make(map[uuid.UUID]*repoCalendar.Calendar, len(calendars))
	for _, calendar := range calendars {
		calendarByID[calendar.ID] = calendar
	}

	settings := make([]CalendarSettingOutput, 0, len(userCalendars))
	for _, userCalendar := range userCalendars {
		calendar, ok := calendarByID[userCalendar.CalendarID]
		if !ok {
			return nil, repoerr.ErrNotFound
		}
		settings = append(settings, toCalendarSettingOutput(userCalendar, calendar))
	}
	return settings, nil
}

func findUserCalendarByID(ctx context.Context, repos CalendarSettingsRepositories, userID, userCalendarID uuid.UUID) (*repoUserCalendar.UserCalendar, error) {
	userCalendar, err := repos.UserCalendar.FindByIDAndUser(ctx, userID, userCalendarID)
	if err != nil {
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("カレンダー設定が見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}
	return userCalendar, nil
}

func demoteExistingPrimary(ctx context.Context, repos CalendarSettingsRepositories, userID, exceptID uuid.UUID) error {
	primary, err := repos.UserCalendar.FindByRole(ctx, userID, value.UserCalendarRolePrimary)
	if err != nil {
		if repoerr.IsNotFound(err) {
			return nil
		}
		return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	if primary.ID == exceptID {
		return nil
	}

	referenceRole := value.UserCalendarRoleReference
	if _, err := repos.UserCalendar.Update(ctx, userID, primary.ID, repoUserCalendar.UserCalendarQueryOptions{Role: &referenceRole}); err != nil {
		return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}
	return nil
}

func toCalendarSettingOutput(userCalendar *repoUserCalendar.UserCalendar, calendar *repoCalendar.Calendar) CalendarSettingOutput {
	return CalendarSettingOutput{
		ID:                userCalendar.ID,
		CalendarID:        calendar.ID,
		GoogleCalendarID:  calendar.GoogleCalendarID,
		Summary:           calendar.Summary,
		Description:       calendar.Description,
		Timezone:          calendar.Timezone,
		Role:              userCalendar.Role,
		IsVisible:         userCalendar.IsVisible,
		SyncProposedDates: userCalendar.SyncProposedDates,
	}
}
