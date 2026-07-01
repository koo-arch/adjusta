package calendar

import (
	"context"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	domainUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

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

func listUserCalendarRelations(ctx context.Context, repos SyncTxRepositories, userID uuid.UUID) ([]*CalendarRelation, error) {
	userCalendars, err := repos.UserCalendar.FilterByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	relations := make([]*CalendarRelation, 0, len(userCalendars))
	for _, userCalendar := range userCalendars {
		calendar, err := repos.Calendar.Read(ctx, userCalendar.CalendarID)
		if err != nil {
			return nil, err
		}

		relations = append(relations, &CalendarRelation{
			CalendarID:        userCalendar.CalendarID,
			GoogleCalendarID:  calendar.GoogleCalendarID,
			Role:              userCalendar.Role,
			SyncProposedDates: userCalendar.SyncProposedDates,
		})
	}

	return relations, nil
}

func findRelationByRole(relations []*CalendarRelation, role value.UserCalendarRole) *CalendarRelation {
	for _, relation := range relations {
		if relation.Role == role {
			return relation
		}
	}
	return nil
}

func findIncomingCalendarByID(calendars []*ExternalCalendar, calendarID string) *ExternalCalendar {
	for _, cal := range calendars {
		if cal.CalendarID == calendarID {
			return cal
		}
	}
	return nil
}

func findAdjustaCandidateCalendar(calendars []*ExternalCalendar) *ExternalCalendar {
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

func resolveAdjustaCandidateSyncProposedDates(relation *CalendarRelation) *bool {
	if relation == nil {
		return nil
	}

	syncProposedDates := relation.SyncProposedDates
	return &syncProposedDates
}

func shouldCreateAdjustaCandidateCalendar(relation *CalendarRelation) bool {
	return relation != nil && relation.SyncProposedDates
}
