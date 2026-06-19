package googlecalendar

import (
	"context"

	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
)

type calendarServiceFactory struct{}

type calendarService struct {
	service *Client
}

func NewCalendarServiceFactory() usecaseCalendar.CalendarServiceFactory {
	return &calendarServiceFactory{}
}

func (f *calendarServiceFactory) New(ctx context.Context, token *appmodel.GoogleAuthToken) (usecaseCalendar.CalendarService, error) {
	service, err := NewClient(ctx, toOAuth2Token(token))
	if err != nil {
		return nil, normalizeGoogleAPIError(err)
	}

	return &calendarService{service: service}, nil
}

func (s *calendarService) FetchCalendarList() ([]*usecaseCalendar.CalendarRecord, error) {
	calendars, err := s.service.FetchCalendarList()
	if err != nil {
		return nil, normalizeGoogleAPIError(err)
	}

	result := make([]*usecaseCalendar.CalendarRecord, 0, len(calendars))
	for _, calendar := range calendars {
		if calendar == nil {
			continue
		}
		result = append(result, &usecaseCalendar.CalendarRecord{
			CalendarID: calendar.CalendarID,
			Summary:    calendar.Summary,
			Primary:    calendar.Primary,
		})
	}

	return result, nil
}

func (s *calendarService) CreateCalendar(summary string) (*usecaseCalendar.CalendarRecord, error) {
	calendar, err := s.service.CreateCalendar(summary)
	if err != nil {
		return nil, normalizeGoogleAPIError(err)
	}

	return &usecaseCalendar.CalendarRecord{
		CalendarID: calendar.CalendarID,
		Summary:    calendar.Summary,
		Primary:    calendar.Primary,
	}, nil
}
