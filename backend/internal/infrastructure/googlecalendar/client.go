package googlecalendar

import (
	"context"
	"time"

	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	infraGoogleOAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/googleoauth"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarList struct {
	CalendarID string `json:"calendar_id"`
	Summary    string `json:"summary"`
	Primary    bool   `json:"primary"`
}

type Client struct {
	Service *calendar.Service
}

func NewClient(ctx context.Context, token *oauth2.Token) (*Client, error) {
	service, err := calendar.NewService(ctx, option.WithTokenSource(infraGoogleOAuth.GetConfig().TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	return &Client{Service: service}, nil
}

func (c *Client) FetchCalendarList() ([]*CalendarList, error) {
	calendarList, err := c.Service.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}

	var calendars []*CalendarList
	for _, item := range calendarList.Items {
		calendar := &CalendarList{
			CalendarID: item.Id,
			Summary:    item.Summary,
			Primary:    item.Primary,
		}
		calendars = append(calendars, calendar)
	}

	return calendars, nil
}

func (c *Client) CreateCalendar(summary string) (*CalendarList, error) {
	created, err := c.Service.Calendars.Insert(&calendar.Calendar{
		Summary: summary,
	}).Do()
	if err != nil {
		return nil, err
	}

	return &CalendarList{
		CalendarID: created.Id,
		Summary:    created.Summary,
		Primary:    false,
	}, nil
}

func (c *Client) FetchEvents(calendarID string, startTime, endTime time.Time) ([]*appmodel.GoogleEvent, error) {
	events, err := c.Service.Events.List(calendarID).
		TimeMin(startTime.Format(time.RFC3339)).
		TimeMax(endTime.Format(time.RFC3339)).
		SingleEvents(true).
		Do()
	if err != nil {
		return nil, err
	}

	var eventsList []*appmodel.GoogleEvent

	for _, item := range events.Items {
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
		event := &appmodel.GoogleEvent{
			ID:          item.Id,
			Summary:     item.Summary,
			Description: item.Description,
			Location:    item.Location,
			ColorID:     item.ColorId,
			Start:       start,
			End:         end,
		}
		eventsList = append(eventsList, event)
	}

	return eventsList, nil
}

func (c *Client) FetchEvent(calendarID, eventID string) (*calendar.Event, error) {
	event, err := c.Service.Events.Get(calendarID, eventID).Do()
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Client) InsertEvent(calendarID string, event *calendar.Event) (*calendar.Event, error) {
	event, err := c.Service.Events.Insert(calendarID, event).Do()
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Client) UpdateEvent(calendarID, eventID string, event *calendar.Event) (*calendar.Event, error) {
	event, err := c.Service.Events.Update(calendarID, eventID, event).Do()
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Client) DeleteEvent(calendarID, eventID string) error {
	err := c.Service.Events.Delete(calendarID, eventID).Do()
	if err != nil {
		return err
	}

	return nil
}
