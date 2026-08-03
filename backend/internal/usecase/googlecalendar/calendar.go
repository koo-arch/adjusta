package googlecalendar

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type candidateCalendar struct {
	googleCalendarID  string
	syncProposedDates bool
}

func (uc *SyncUsecase) loadCandidateCalendar(ctx context.Context, userID uuid.UUID) (*candidateCalendar, error) {
	userCalendars, err := uc.repos.UserCalendar.FilterByUserID(ctx, userID)
	if err != nil {
		if repoerr.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, userCalendar := range userCalendars {
		if userCalendar == nil || userCalendar.Role != value.UserCalendarRoleAdjustaCandidate {
			continue
		}
		calendar, err := uc.repos.Calendar.Read(ctx, userCalendar.CalendarID)
		if err != nil {
			if repoerr.IsNotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		return &candidateCalendar{
			googleCalendarID:  calendar.GoogleCalendarID,
			syncProposedDates: userCalendar.SyncProposedDates,
		}, nil
	}
	return nil, nil
}
