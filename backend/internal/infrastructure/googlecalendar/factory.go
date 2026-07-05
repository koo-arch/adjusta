package googlecalendar

import (
	"context"

	"github.com/koo-arch/adjusta-backend/internal/google"
	infraGoogleOAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/googleoauth"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
)

type calendarServiceFactory struct {
	oauthClient *infraGoogleOAuth.Client
}

type calendarService struct {
	service *Client
}

func NewCalendarServiceFactory(oauthClient *infraGoogleOAuth.Client) usecaseCalendar.CalendarServiceFactory {
	return &calendarServiceFactory{oauthClient: oauthClient}
}

func (f *calendarServiceFactory) New(ctx context.Context, token *google.AuthToken) (usecaseCalendar.CalendarService, error) {
	service, err := NewClient(ctx, f.oauthClient, toOAuth2Token(token))
	if err != nil {
		return nil, normalizeGoogleAPIError(err)
	}

	return &calendarService{service: service}, nil
}

func (s *calendarService) FetchCalendarList() ([]*usecaseCalendar.ExternalCalendar, error) {
	calendars, err := s.service.FetchCalendarList()
	if err != nil {
		return nil, normalizeGoogleAPIError(err)
	}

	result := make([]*usecaseCalendar.ExternalCalendar, 0, len(calendars))
	for _, calendar := range calendars {
		if calendar == nil {
			continue
		}
		result = append(result, &usecaseCalendar.ExternalCalendar{
			CalendarID: calendar.CalendarID,
			Summary:    calendar.Summary,
			Primary:    calendar.Primary,
		})
	}

	return result, nil
}

func (s *calendarService) CreateCalendar(summary string) (*usecaseCalendar.ExternalCalendar, error) {
	calendar, err := s.service.CreateCalendar(summary)
	if err != nil {
		return nil, normalizeGoogleAPIError(err)
	}

	return &usecaseCalendar.ExternalCalendar{
		CalendarID: calendar.CalendarID,
		Summary:    calendar.Summary,
		Primary:    calendar.Primary,
	}, nil
}
