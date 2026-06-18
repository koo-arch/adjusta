package calendar

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	domainUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type SyncUsecase struct {
	userReader             UserReader
	googleTokenProvider    GoogleTokenProvider
	calendarServiceFactory CalendarServiceFactory
	tx                     SyncTransaction
}

func NewSyncUsecase(
	userReader UserReader,
	googleTokenProvider GoogleTokenProvider,
	calendarServiceFactory CalendarServiceFactory,
	tx SyncTransaction,
) *SyncUsecase {
	return &SyncUsecase{
		userReader:             userReader,
		googleTokenProvider:    googleTokenProvider,
		calendarServiceFactory: calendarServiceFactory,
		tx:                     tx,
	}
}

func (uc *SyncUsecase) SyncGoogleCalendars(ctx context.Context, userID uuid.UUID, email string) ([]*appmodel.GoogleCalendarList, error) {
	entUser, err := uc.userReader.GetByID(ctx, userID)
	if err != nil {
		log.Printf("failed to get user info for account: %s, %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("ユーザー情報が見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError("ユーザー情報取得時にエラーが発生しました")
	}

	token, err := uc.googleTokenProvider.GetToken(ctx, entUser.ID)
	if err != nil {
		log.Printf("failed to verify token for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.NormalizeAPIError(err, "OAuthトークンの認証に失敗しました")
	}

	calendarService, err := uc.calendarServiceFactory.New(ctx, token)
	if err != nil {
		log.Printf("failed to create calendar service for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.NewInternalError("Googleカレンダー接続に失敗しました")
	}

	calendars, err := calendarService.FetchCalendarList()
	if err != nil {
		log.Printf("failed to fetch calendars for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.NormalizeAPIError(err, "Googleカレンダー情報の取得に失敗しました")
	}

	calendars, err = uc.syncCalendar(ctx, calendarService, calendars, entUser)
	if err != nil {
		log.Printf("failed to sync calendars for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return calendars, nil
}

func (uc *SyncUsecase) syncCalendar(ctx context.Context, calendarService CalendarService, calendars []*appmodel.GoogleCalendarList, entUser *repoUser.User) ([]*appmodel.GoogleCalendarList, error) {
	syncedCalendars := calendars

	err := uc.tx.Do(ctx, func(store SyncStore) error {
		relations, err := store.ListUserCalendarRelations(ctx, entUser.ID)
		if err != nil {
			return fmt.Errorf("failed to list user calendar relations: %w", err)
		}

		adjustaCandidate, err := uc.ensureAdjustaCandidateCalendar(ctx, calendarService, store, entUser.ID, calendars, relations)
		if err != nil {
			return err
		}

		adjustaCandidateID := ""
		if adjustaCandidate != nil {
			adjustaCandidateID = adjustaCandidate.CalendarID
			if findIncomingCalendarByID(calendars, adjustaCandidateID) == nil {
				syncedCalendars = append(append([]*appmodel.GoogleCalendarList{}, calendars...), adjustaCandidate)
			}
		}

		incoming := make(map[string]struct{}, len(calendars))
		for _, cal := range calendars {
			if cal.CalendarID == adjustaCandidateID {
				continue
			}

			incoming[cal.CalendarID] = struct{}{}

			storedCalendar, err := uc.ensureStoredCalendar(ctx, store, entUser.ID, cal.CalendarID, cal.Summary)
			if err != nil {
				return err
			}

			role := domainUserCalendar.ExternalSyncRole(cal.Primary)
			if _, err := store.EnsureUserCalendarRelation(ctx, entUser.ID, storedCalendar.ID, role, nil); err != nil {
				return fmt.Errorf("failed to ensure user calendar relation: %w", err)
			}
		}

		relations, err = store.ListUserCalendarRelations(ctx, entUser.ID)
		if err != nil {
			return fmt.Errorf("failed to list user calendar relations: %w", err)
		}

		for _, relation := range relations {
			if !domainUserCalendar.IsExternalSyncRole(relation.Role) {
				continue
			}
			if _, ok := incoming[relation.GoogleCalendarID]; !ok {
				if err := store.SoftDeleteUserCalendarRelation(ctx, entUser.ID, relation.CalendarID); err != nil {
					return fmt.Errorf("failed to soft delete user calendar relation: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return syncedCalendars, nil
}

func (uc *SyncUsecase) ensureAdjustaCandidateCalendar(
	ctx context.Context,
	calendarService CalendarService,
	store SyncStore,
	userID uuid.UUID,
	calendars []*appmodel.GoogleCalendarList,
	relations []*UserCalendarRelationRecord,
) (*appmodel.GoogleCalendarList, error) {
	existingRelation := findRelationByRole(relations, domainvalue.UserCalendarRoleAdjustaCandidate)
	syncProposedDates := resolveAdjustaCandidateSyncProposedDates(existingRelation)

	if existingRelation != nil {
		current := findIncomingCalendarByID(calendars, existingRelation.GoogleCalendarID)
		if current != nil {
			if _, err := store.UpdateCalendar(ctx, existingRelation.CalendarID, current.CalendarID, current.Summary); err != nil {
				return nil, fmt.Errorf("failed to update adjusta candidate calendar: %w", err)
			}
			if _, err := store.EnsureUserCalendarRelation(ctx, userID, existingRelation.CalendarID, domainvalue.UserCalendarRoleAdjustaCandidate, syncProposedDates); err != nil {
				return nil, fmt.Errorf("failed to ensure adjusta candidate relation: %w", err)
			}
			return current, nil
		}
	}

	desired := findAdjustaCandidateCalendar(calendars)
	if desired == nil {
		if !shouldCreateAdjustaCandidateCalendar(existingRelation) {
			return nil, nil
		}

		var err error
		desired, err = calendarService.CreateCalendar(domainUserCalendar.AdjustaCandidateCalendarSummary)
		if err != nil {
			return nil, fmt.Errorf("failed to create adjusta candidate calendar: %w", err)
		}
	}

	storedCalendar, err := uc.ensureStoredCalendar(ctx, store, userID, desired.CalendarID, desired.Summary)
	if err != nil {
		return nil, err
	}

	if existingRelation != nil && existingRelation.CalendarID != storedCalendar.ID {
		if err := store.SoftDeleteUserCalendarRelation(ctx, userID, existingRelation.CalendarID); err != nil {
			return nil, fmt.Errorf("failed to replace adjusta candidate relation: %w", err)
		}
	}

	if _, err := store.EnsureUserCalendarRelation(ctx, userID, storedCalendar.ID, domainvalue.UserCalendarRoleAdjustaCandidate, syncProposedDates); err != nil {
		return nil, fmt.Errorf("failed to ensure adjusta candidate relation: %w", err)
	}

	return desired, nil
}

func (uc *SyncUsecase) ensureStoredCalendar(
	ctx context.Context,
	store SyncStore,
	userID uuid.UUID,
	googleCalendarID, summary string,
) (*repoCalendar.Calendar, error) {
	storedCalendar, err := store.FindCalendarByGoogleCalendarID(ctx, userID, googleCalendarID)
	if err != nil {
		if !repoerr.IsNotFound(err) {
			return nil, fmt.Errorf("failed to find calendar: %w", err)
		}

		storedCalendar, err = store.FindAnyCalendarByGoogleCalendarID(ctx, googleCalendarID)
		if err != nil {
			if !repoerr.IsNotFound(err) {
				return nil, fmt.Errorf("failed to find global calendar: %w", err)
			}

			storedCalendar, err = store.CreateCalendar(ctx, googleCalendarID, summary)
			if err != nil {
				return nil, fmt.Errorf("failed to create calendar: %w", err)
			}
		}
	}

	storedCalendar, err = store.UpdateCalendar(ctx, storedCalendar.ID, googleCalendarID, summary)
	if err != nil {
		return nil, fmt.Errorf("failed to update calendar: %w", err)
	}

	return storedCalendar, nil
}

func findRelationByRole(relations []*UserCalendarRelationRecord, role domainvalue.UserCalendarRole) *UserCalendarRelationRecord {
	for _, relation := range relations {
		if relation.Role == role {
			return relation
		}
	}
	return nil
}

func findIncomingCalendarByID(calendars []*appmodel.GoogleCalendarList, calendarID string) *appmodel.GoogleCalendarList {
	for _, cal := range calendars {
		if cal.CalendarID == calendarID {
			return cal
		}
	}
	return nil
}

func findAdjustaCandidateCalendar(calendars []*appmodel.GoogleCalendarList) *appmodel.GoogleCalendarList {
	for _, cal := range calendars {
		if cal.Primary {
			continue
		}
		if domainUserCalendar.IsAdjustaCandidateCalendarSummary(cal.Summary) {
			return cal
		}
	}
	return nil
}

func resolveAdjustaCandidateSyncProposedDates(relation *UserCalendarRelationRecord) *bool {
	if relation == nil {
		return nil
	}

	syncProposedDates := relation.SyncProposedDates
	return &syncProposedDates
}

func shouldCreateAdjustaCandidateCalendar(relation *UserCalendarRelationRecord) bool {
	return relation != nil && relation.SyncProposedDates
}
