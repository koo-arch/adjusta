package events

import (
	"context"
	"log"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) loadPrimaryCalendar(ctx context.Context, repos EventTxRepositories, userID uuid.UUID, email string) (*CalendarRecord, error) {
	storedCalendar, err := findPrimaryCalendar(ctx, repos, userID)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return storedCalendar, nil
}

func (uc *Usecase) loadAdjustaCandidateCalendar(ctx context.Context, repos EventTxRepositories, userID uuid.UUID, email string) (*CalendarRecord, error) {
	storedCalendar, err := findAdjustaCandidateCalendar(ctx, repos, userID)
	if err != nil {
		log.Printf("failed to get adjusta candidate calendar for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return &CalendarRecord{SyncProposedDates: false}, nil
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return storedCalendar, nil
}

func findPrimaryCalendar(ctx context.Context, repos EventTxRepositories, userID uuid.UUID) (*CalendarRecord, error) {
	role := domainvalue.UserCalendarRolePrimary
	calendar, err := repos.Calendar.FindByFields(ctx, userID, repoCalendar.CalendarQueryOptions{
		Role: &role,
	})
	if err != nil {
		return nil, err
	}
	return toCalendarRecord(calendar), nil
}

func findAdjustaCandidateCalendar(ctx context.Context, repos EventTxRepositories, userID uuid.UUID) (*CalendarRecord, error) {
	userCalendars, err := repos.UserCalendar.FilterByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, userCalendar := range userCalendars {
		if userCalendar.Role != domainvalue.UserCalendarRoleAdjustaCandidate {
			continue
		}

		calendar, err := repos.Calendar.Read(ctx, userCalendar.CalendarID)
		if err != nil {
			return nil, err
		}

		return toCalendarRecordWithSync(calendar, userCalendar.SyncProposedDates), nil
	}

	return nil, repoerr.ErrNotFound
}
