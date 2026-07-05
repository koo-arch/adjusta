package calendar

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	domainUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type CalendarRelation struct {
	CalendarID        uuid.UUID
	GoogleCalendarID  string
	Role              value.UserCalendarRole
	SyncProposedDates bool
}

func (uc *SyncUsecase) syncCalendar(ctx context.Context, calendarService CalendarService, calendars []*ExternalCalendar, entUser *repoUser.User) ([]*ExternalCalendar, error) {
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
				syncedCalendars = append(append([]*ExternalCalendar{}, calendars...), adjustaCandidate)
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
	calendars []*ExternalCalendar,
	relations []*CalendarRelation,
) (*ExternalCalendar, error) {
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
