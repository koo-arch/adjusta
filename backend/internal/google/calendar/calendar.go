package calendar

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"github.com/koo-arch/adjusta-backend/internal/google/oauth"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

type CalendarList struct {
	CalendarID string `json:"calendar_id"`
	Summary    string `json:"summary"`
	Primary   bool   `json:"primary"`
}

type Calendar struct {
	Service *calendar.Service
}

func NewCalendar(ctx context.Context, token *oauth2.Token) (*Calendar, error) {
	service, err := calendar.NewService(ctx, option.WithTokenSource(oauth.GoogleOAuthConfig.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	return &Calendar{Service: service}, nil
}

func (c *Calendar) FetchCalendarList() ([]*CalendarList, error) {
	calendarList, err := c.Service.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}

	var calendars []*CalendarList
	for _, item := range calendarList.Items {
		calendar := &CalendarList{
			CalendarID: item.Id,
			Summary:    item.Summary,
			Primary:   item.Primary,
		}
		calendars = append(calendars, calendar)
	}

	return calendars, nil
}

func (c *Calendar) FetchEvents(calendarID string, startTime, endTime time.Time) ([]*models.GoogleEvent, error) {
	events, err := c.Service.Events.List(calendarID).
		TimeMin(startTime.Format(time.RFC3339)).
		TimeMax(endTime.Format(time.RFC3339)).
		SingleEvents(true).
		Do()
	if err != nil {
		return nil, err
	}

	var eventsList []*models.GoogleEvent

	for _, item := range events.Items {
		// nilチェックを追加
		var start, end string
		if item.Start != nil {
			start = item.Start.DateTime
			if start == "" {
				start = item.Start.Date
			}
		}
		if item.End != nil {
			end = item.End.DateTime
			if end == "" {
				end = item.End.Date
			}
		}
		event := &models.GoogleEvent{
			ID:      item.Id,
			Summary: item.Summary,
			Description: item.Description,
			Location: item.Location,
			ColorID: item.ColorId,
			Start:   start,
			End:     end,
		}
		eventsList = append(eventsList, event)
	}

	return eventsList, nil
}

func (c *Calendar) FetchEvent(eventID string) (*calendar.Event, error) {
	event, err := c.Service.Events.Get("primary", eventID).Do()
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Calendar) InsertEvent(event *calendar.Event) (*calendar.Event, error) {
	event, err := c.Service.Events.Insert("primary", event).Do()
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Calendar) UpdateEvent(eventID string, event *calendar.Event) (*calendar.Event, error) {
	event, err := c.Service.Events.Update("primary", eventID, event).Do()
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Calendar) DeleteEvent(eventID string) error {
	err := c.Service.Events.Delete("primary", eventID).Do()
	if err != nil {
		return err
	}

	return nil
}