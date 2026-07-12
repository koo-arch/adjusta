package calendarsetting

import (
	"context"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
)

type CalendarSettingsRepositories struct {
	Calendar     repoCalendar.CalendarRepository
	UserCalendar repoUserCalendar.UserCalendarRepository
}

type CalendarSettingsTransaction interface {
	DoCalendarSettings(ctx context.Context, fn func(repos CalendarSettingsRepositories) error) error
}

type CandidateCalendarEnabler interface {
	EnableAdjustaCandidateCalendar(ctx context.Context, userID uuid.UUID, email string) error
}

type CandidateCalendarEnablerFunc func(ctx context.Context, userID uuid.UUID, email string) error

func (f CandidateCalendarEnablerFunc) EnableAdjustaCandidateCalendar(ctx context.Context, userID uuid.UUID, email string) error {
	return f(ctx, userID, email)
}
