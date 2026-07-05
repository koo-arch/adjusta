package events

import (
	"context"
	"log"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type EventCalendar struct {
	ID                uuid.UUID
	GoogleCalendarID  string
	Summary           string
	Description       *string
	Timezone          *string
	SyncProposedDates bool
}

func (uc *Usecase) loadPrimaryCalendar(ctx context.Context, repos EventTxRepositories, userID uuid.UUID, email string) (*EventCalendar, error) {
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

func (uc *Usecase) loadAdjustaCandidateCalendar(ctx context.Context, repos EventTxRepositories, userID uuid.UUID, email string) (*EventCalendar, error) {
	storedCalendar, err := findAdjustaCandidateCalendar(ctx, repos, userID)
	if err != nil {
		log.Printf("failed to get adjusta candidate calendar for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return &EventCalendar{SyncProposedDates: false}, nil
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return storedCalendar, nil
}

func findPrimaryCalendar(ctx context.Context, repos EventTxRepositories, userID uuid.UUID) (*EventCalendar, error) {
	role := value.UserCalendarRolePrimary
	calendar, err := repos.Calendar.FindByFields(ctx, userID, repoCalendar.CalendarQueryOptions{
		Role: &role,
	})
	if err != nil {
		return nil, err
	}
	return toEventCalendar(calendar), nil
}

func findAdjustaCandidateCalendar(ctx context.Context, repos EventTxRepositories, userID uuid.UUID) (*EventCalendar, error) {
	userCalendars, err := repos.UserCalendar.FilterByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, userCalendar := range userCalendars {
		if userCalendar.Role != value.UserCalendarRoleAdjustaCandidate {
			continue
		}

		calendar, err := repos.Calendar.Read(ctx, userCalendar.CalendarID)
		if err != nil {
			return nil, err
		}

		return toEventCalendarWithSync(calendar, userCalendar.SyncProposedDates), nil
	}

	return nil, repoerr.ErrNotFound
}

func toEventCalendar(calendar *repoCalendar.Calendar) *EventCalendar {
	return toEventCalendarWithSync(calendar, false)
}

func toEventCalendarWithSync(calendar *repoCalendar.Calendar, syncProposedDates bool) *EventCalendar {
	if calendar == nil {
		return nil
	}

	return &EventCalendar{
		ID:                calendar.ID,
		GoogleCalendarID:  calendar.GoogleCalendarID,
		Summary:           calendar.Summary,
		Description:       calendar.Description,
		Timezone:          calendar.Timezone,
		SyncProposedDates: syncProposedDates,
	}
}

func toEventCalendars(calendars []*repoCalendar.Calendar) []*EventCalendar {
	eventCalendars := make([]*EventCalendar, 0, len(calendars))
	for _, calendar := range calendars {
		eventCalendars = append(eventCalendars, toEventCalendar(calendar))
	}
	return eventCalendars
}
