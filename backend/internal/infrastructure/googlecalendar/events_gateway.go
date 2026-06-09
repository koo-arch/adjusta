package googlecalendar

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type eventGateway struct {
	googleTokenProvider usecaseEvents.GoogleTokenProvider
	calendarManager     *GoogleCalendarManager
}

func NewEventGateway(googleTokenProvider usecaseEvents.GoogleTokenProvider, calendarManager *GoogleCalendarManager) usecaseEvents.GoogleCalendarGateway {
	return &eventGateway{
		googleTokenProvider: googleTokenProvider,
		calendarManager:     calendarManager,
	}
}

func (g *eventGateway) FetchEvents(ctx context.Context, userID uuid.UUID, calendars []*repositorymodel.StoredCalendar, startTime, endTime time.Time) (*usecaseEvents.GoogleEventFetchResult, error) {
	calendarService, err := g.newCalendarService(ctx, userID)
	if err != nil {
		return nil, err
	}

	result, err := g.calendarManager.FetchEventsFromCalendars(calendarService, calendars, startTime, endTime)
	if result == nil {
		return nil, normalizeGoogleAPIError(err)
	}

	return &usecaseEvents.GoogleEventFetchResult{
		Events:          result.Events,
		FailedCalendars: result.FailedCalendars,
	}, normalizeGoogleAPIError(err)
}

func (g *eventGateway) UpsertEvent(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
	calendarService, err := g.newCalendarService(ctx, userID)
	if err != nil {
		return "", err
	}

	if existingGoogleEventID == nil || *existingGoogleEventID == "" {
		eventReq := &appmodel.EventDraftCreation{
			Title:       title,
			Location:    location,
			Description: description,
			SelectedDates: []appmodel.SelectedDate{
				{
					Start: start,
					End:   end,
				},
			},
		}

		insertedEvents, err := g.calendarManager.CreateGoogleEvents(calendarService, calendarID, eventReq)
		if err != nil {
			return "", fmt.Errorf("failed to create google event: %w", err)
		}
		if len(insertedEvents) == 0 || insertedEvents[0] == nil {
			return "", fmt.Errorf("failed to create google event: empty response")
		}
		return insertedEvents[0].Id, nil
	}

	googleEvent := g.calendarManager.ConvertToCalendarEvent(existingGoogleEventID, title, location, description, start, end)
	upsertedEvent, err := g.calendarManager.UpdateOrCreateGoogleEvent(calendarService, calendarID, googleEvent)
	if err != nil {
		return "", fmt.Errorf("failed to upsert google event: %w", err)
	}

	return upsertedEvent.Id, nil
}

func (g *eventGateway) newCalendarService(ctx context.Context, userID uuid.UUID) (*customCalendar.Calendar, error) {
	token, err := g.googleTokenProvider.GetToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	calendarService, err := customCalendar.NewCalendar(ctx, toOAuth2Token(token))
	if err != nil {
		return nil, normalizeGoogleAPIError(err)
	}

	return calendarService, nil
}
