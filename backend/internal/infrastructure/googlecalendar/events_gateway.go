package googlecalendar

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
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

func (g *eventGateway) FetchEvents(ctx context.Context, userID uuid.UUID, calendars []*usecaseEvents.EventCalendar, startTime, endTime time.Time) (*usecaseEvents.GoogleEventFetchResult, error) {
	calendarService, err := g.newCalendarService(ctx, userID)
	if err != nil {
		return nil, err
	}

	result, err := g.calendarManager.FetchEventsFromCalendars(calendarService, toCalendars(calendars), startTime, endTime)
	if result == nil {
		return nil, normalizeGoogleAPIError(err)
	}

	return &usecaseEvents.GoogleEventFetchResult{
		Events:          toFetchedGoogleEvents(result.Events),
		FailedCalendars: result.FailedCalendars,
	}, normalizeGoogleAPIError(err)
}

func toFetchedGoogleEvents(events []*FetchedEvent) []*usecaseEvents.FetchedGoogleEvent {
	outputs := make([]*usecaseEvents.FetchedGoogleEvent, 0, len(events))
	for _, event := range events {
		if event == nil {
			continue
		}
		outputs = append(outputs, &usecaseEvents.FetchedGoogleEvent{
			ID:          event.ID,
			Summary:     event.Summary,
			Description: event.Description,
			Location:    event.Location,
			ColorID:     event.ColorID,
			Start:       event.Start,
			End:         event.End,
		})
	}
	return outputs
}

func toCalendars(calendars []*usecaseEvents.EventCalendar) []*repoCalendar.Calendar {
	domainCalendars := make([]*repoCalendar.Calendar, 0, len(calendars))
	for _, calendar := range calendars {
		if calendar == nil {
			continue
		}
		domainCalendars = append(domainCalendars, &repoCalendar.Calendar{
			ID:               calendar.ID,
			GoogleCalendarID: calendar.GoogleCalendarID,
			Summary:          calendar.Summary,
			Description:      calendar.Description,
			Timezone:         calendar.Timezone,
		})
	}
	return domainCalendars
}

func (g *eventGateway) UpsertEvent(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
	calendarService, err := g.newCalendarService(ctx, userID)
	if err != nil {
		return "", err
	}

	if existingGoogleEventID == nil || *existingGoogleEventID == "" {
		eventReq := EventDraft{
			Title:       title,
			Location:    location,
			Description: description,
			Dates: []EventDate{
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

func (g *eventGateway) newCalendarService(ctx context.Context, userID uuid.UUID) (*Client, error) {
	token, err := g.googleTokenProvider.GetToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	calendarService, err := NewClient(ctx, toOAuth2Token(token))
	if err != nil {
		return nil, normalizeGoogleAPIError(err)
	}

	return calendarService, nil
}
