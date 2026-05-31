package googlecalendar

import (
	"context"

	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
)

type calendarServiceFactory struct{}

type calendarService struct {
	service *customCalendar.Calendar
}

func NewCalendarServiceFactory() usecaseCalendar.CalendarServiceFactory {
	return &calendarServiceFactory{}
}

func (f *calendarServiceFactory) New(ctx context.Context, token *appmodel.GoogleAuthToken) (usecaseCalendar.CalendarService, error) {
	service, err := customCalendar.NewCalendar(ctx, toOAuth2Token(token))
	if err != nil {
		return nil, normalizeGoogleAPIError(err)
	}

	return &calendarService{service: service}, nil
}

func (s *calendarService) FetchCalendarList() ([]*customCalendar.CalendarList, error) {
	calendars, err := s.service.FetchCalendarList()
	if err != nil {
		return nil, normalizeGoogleAPIError(err)
	}

	return calendars, nil
}
