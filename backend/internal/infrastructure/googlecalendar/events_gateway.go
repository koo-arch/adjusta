package googlecalendar

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
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

func (g *eventGateway) FetchEvents(ctx context.Context, userID uuid.UUID, calendars []*models.GoogleCalendarInfo, startTime, endTime time.Time) (*usecaseEvents.GoogleEventFetchResult, error) {
	calendarService, err := g.newCalendarService(ctx, userID)
	if err != nil {
		return nil, err
	}

	result, err := g.calendarManager.FetchEventsFromCalendars(calendarService, calendars, startTime, endTime)
	if result == nil {
		return nil, err
	}

	return &usecaseEvents.GoogleEventFetchResult{
		Events:          result.Events,
		FailedCalendars: result.FailedCalendars,
	}, err
}

func (g *eventGateway) UpsertEvent(ctx context.Context, userID uuid.UUID, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
	calendarService, err := g.newCalendarService(ctx, userID)
	if err != nil {
		return "", err
	}

	if existingGoogleEventID == nil || *existingGoogleEventID == "" {
		eventReq := &models.EventDraftCreation{
			Title:       title,
			Location:    location,
			Description: description,
			SelectedDates: []models.SelectedDate{
				{
					Start: start,
					End:   end,
				},
			},
		}

		insertedEvents, err := g.calendarManager.CreateGoogleEvents(calendarService, eventReq)
		if err != nil {
			return "", fmt.Errorf("failed to create google event: %w", err)
		}
		if len(insertedEvents) == 0 || insertedEvents[0] == nil {
			return "", fmt.Errorf("failed to create google event: empty response")
		}
		return insertedEvents[0].Id, nil
	}

	googleEvent := g.calendarManager.ConvertToCalendarEvent(existingGoogleEventID, title, location, description, start, end)
	upsertedEvent, err := g.calendarManager.UpdateOrCreateGoogleEvent(calendarService, googleEvent)
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

	return customCalendar.NewCalendar(ctx, token)
}
