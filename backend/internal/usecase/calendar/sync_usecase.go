package calendar

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	domainUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type SyncUsecase struct {
	userRepo               repoUser.UserRepository
	googleTokenProvider    GoogleTokenProvider
	calendarServiceFactory CalendarServiceFactory
	tx                     SyncTransaction
}

func NewSyncUsecase(
	userRepo repoUser.UserRepository,
	googleTokenProvider GoogleTokenProvider,
	calendarServiceFactory CalendarServiceFactory,
	tx SyncTransaction,
) *SyncUsecase {
	return &SyncUsecase{
		userRepo:               userRepo,
		googleTokenProvider:    googleTokenProvider,
		calendarServiceFactory: calendarServiceFactory,
		tx:                     tx,
	}
}

func (uc *SyncUsecase) SyncGoogleCalendars(ctx context.Context, userID uuid.UUID, email string) ([]*CalendarRecord, error) {
	entUser, err := uc.userRepo.Read(ctx, userID)
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

func (uc *SyncUsecase) syncCalendar(ctx context.Context, calendarService CalendarService, calendars []*CalendarRecord, entUser *repoUser.User) ([]*CalendarRecord, error) {
	syncedCalendars := calendars

	err := uc.tx.Do(ctx, func(repos SyncTxRepositories) error {
		relations, err := listUserCalendarRelations(ctx, repos, entUser.ID)
		if err != nil {
			return fmt.Errorf("failed to list user calendar relations: %w", err)
		}

		adjustaCandidate, err := uc.ensureAdjustaCandidateCalendar(ctx, calendarService, repos, entUser.ID, calendars, relations)
		if err != nil {
			return err
		}

		adjustaCandidateID := ""
		if adjustaCandidate != nil {
			adjustaCandidateID = adjustaCandidate.CalendarID
			if findIncomingCalendarByID(calendars, adjustaCandidateID) == nil {
				syncedCalendars = append(append([]*CalendarRecord{}, calendars...), adjustaCandidate)
			}
		}

		incoming := make(map[string]struct{}, len(calendars))
		for _, cal := range calendars {
			if cal.CalendarID == adjustaCandidateID {
				continue
			}

			incoming[cal.CalendarID] = struct{}{}

			storedCalendar, err := uc.ensureStoredCalendar(ctx, repos, entUser.ID, cal.CalendarID, cal.Summary)
			if err != nil {
				return err
			}

			role := domainUserCalendar.ExternalSyncRole(cal.Primary)
			if _, err := ensureUserCalendarRelation(ctx, repos, entUser.ID, storedCalendar.ID, role, nil); err != nil {
				return fmt.Errorf("failed to ensure user calendar relation: %w", err)
			}
		}

		relations, err = listUserCalendarRelations(ctx, repos, entUser.ID)
		if err != nil {
			return fmt.Errorf("failed to list user calendar relations: %w", err)
		}

		for _, relation := range relations {
			if !domainUserCalendar.IsExternalSyncRole(relation.Role) {
				continue
			}
			if _, ok := incoming[relation.GoogleCalendarID]; !ok {
				if err := repos.UserCalendar.SoftDeleteByUserAndCalendar(ctx, entUser.ID, relation.CalendarID); err != nil {
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
	repos SyncTxRepositories,
	userID uuid.UUID,
	calendars []*CalendarRecord,
	relations []*UserCalendarRelationRecord,
) (*CalendarRecord, error) {
	existingRelation := findRelationByRole(relations, value.UserCalendarRoleAdjustaCandidate)
	syncProposedDates := resolveAdjustaCandidateSyncProposedDates(existingRelation)

	if existingRelation != nil {
		current := findIncomingCalendarByID(calendars, existingRelation.GoogleCalendarID)
		if current != nil {
			if _, err := updateCalendar(ctx, repos, existingRelation.CalendarID, current.CalendarID, current.Summary); err != nil {
				return nil, fmt.Errorf("failed to update adjusta candidate calendar: %w", err)
			}
			if _, err := ensureUserCalendarRelation(ctx, repos, userID, existingRelation.CalendarID, value.UserCalendarRoleAdjustaCandidate, syncProposedDates); err != nil {
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

	storedCalendar, err := uc.ensureStoredCalendar(ctx, repos, userID, desired.CalendarID, desired.Summary)
	if err != nil {
		return nil, err
	}

	if existingRelation != nil && existingRelation.CalendarID != storedCalendar.ID {
		if err := repos.UserCalendar.SoftDeleteByUserAndCalendar(ctx, userID, existingRelation.CalendarID); err != nil {
			return nil, fmt.Errorf("failed to replace adjusta candidate relation: %w", err)
		}
	}

	if _, err := ensureUserCalendarRelation(ctx, repos, userID, storedCalendar.ID, value.UserCalendarRoleAdjustaCandidate, syncProposedDates); err != nil {
		return nil, fmt.Errorf("failed to ensure adjusta candidate relation: %w", err)
	}

	return desired, nil
}

func (uc *SyncUsecase) ensureStoredCalendar(
	ctx context.Context,
	repos SyncTxRepositories,
	userID uuid.UUID,
	googleCalendarID, summary string,
) (*repoCalendar.Calendar, error) {
	storedCalendar, err := repos.Calendar.FindByFields(ctx, userID, repoCalendar.CalendarQueryOptions{
		GoogleCalendarID: &googleCalendarID,
	})
	if err != nil {
		if !repoerr.IsNotFound(err) {
			return nil, fmt.Errorf("failed to find calendar: %w", err)
		}

		storedCalendar, err = repos.Calendar.FindByGoogleCalendarID(ctx, googleCalendarID)
		if err != nil {
			if !repoerr.IsNotFound(err) {
				return nil, fmt.Errorf("failed to find global calendar: %w", err)
			}

			storedCalendar, err = createCalendar(ctx, repos, googleCalendarID, summary)
			if err != nil {
				return nil, fmt.Errorf("failed to create calendar: %w", err)
			}
		}
	}

	storedCalendar, err = updateCalendar(ctx, repos, storedCalendar.ID, googleCalendarID, summary)
	if err != nil {
		return nil, fmt.Errorf("failed to update calendar: %w", err)
	}

	return storedCalendar, nil
}

func createCalendar(ctx context.Context, repos SyncTxRepositories, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
	return repos.Calendar.Create(ctx, repoCalendar.CalendarMutationOptions{
		GoogleCalendarID: &googleCalendarID,
		Summary:          &summary,
	})
}

func updateCalendar(ctx context.Context, repos SyncTxRepositories, id uuid.UUID, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
	return repos.Calendar.Update(ctx, id, repoCalendar.CalendarMutationOptions{
		GoogleCalendarID: &googleCalendarID,
		Summary:          &summary,
	})
}

func ensureUserCalendarRelation(ctx context.Context, repos SyncTxRepositories, userID, calendarID uuid.UUID, role value.UserCalendarRole, syncProposedDates *bool) (*domainUserCalendar.UserCalendar, error) {
	isVisible := true
	return repos.UserCalendar.Ensure(ctx, userID, calendarID, domainUserCalendar.UserCalendarQueryOptions{
		Role:              &role,
		IsVisible:         &isVisible,
		SyncProposedDates: syncProposedDates,
	})
}

func listUserCalendarRelations(ctx context.Context, repos SyncTxRepositories, userID uuid.UUID) ([]*UserCalendarRelationRecord, error) {
	userCalendars, err := repos.UserCalendar.FilterByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	relations := make([]*UserCalendarRelationRecord, 0, len(userCalendars))
	for _, userCalendar := range userCalendars {
		calendar, err := repos.Calendar.Read(ctx, userCalendar.CalendarID)
		if err != nil {
			return nil, err
		}

		relations = append(relations, &UserCalendarRelationRecord{
			CalendarID:        userCalendar.CalendarID,
			GoogleCalendarID:  calendar.GoogleCalendarID,
			Role:              userCalendar.Role,
			SyncProposedDates: userCalendar.SyncProposedDates,
		})
	}

	return relations, nil
}

func findRelationByRole(relations []*UserCalendarRelationRecord, role value.UserCalendarRole) *UserCalendarRelationRecord {
	for _, relation := range relations {
		if relation.Role == role {
			return relation
		}
	}
	return nil
}

func findIncomingCalendarByID(calendars []*CalendarRecord, calendarID string) *CalendarRecord {
	for _, cal := range calendars {
		if cal.CalendarID == calendarID {
			return cal
		}
	}
	return nil
}

func findAdjustaCandidateCalendar(calendars []*CalendarRecord) *CalendarRecord {
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
