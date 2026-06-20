package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
)

type GoogleTokenProvider interface {
	GetToken(ctx context.Context, userID uuid.UUID) (*appmodel.GoogleAuthToken, error)
}

type PrimaryCalendarFinder interface {
	FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*CalendarRecord, error)
}

type AdjustaCandidateCalendarFinder interface {
	FindAdjustaCandidateCalendar(ctx context.Context, userID uuid.UUID) (*CalendarRecord, error)
}

type EventReader interface {
	PrimaryCalendarFinder
	AdjustaCandidateCalendarFinder
	ListCalendarsByUser(ctx context.Context, userID uuid.UUID) ([]*CalendarRecord, error)
	SearchEvents(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*EventRecord, error)
	FindEventByID(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*EventRecord, error)
}

type EventTxStore interface {
	PrimaryCalendarFinder
	AdjustaCandidateCalendarFinder
	FindEventByID(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*EventRecord, error)
	ReadCalendar(ctx context.Context, calendarID uuid.UUID) (*CalendarRecord, error)
	CreateEvent(ctx context.Context, userID, primaryCalendarID uuid.UUID, title, location, description string, start, end time.Time) (*EventRecord, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, opt EventMutation) (*EventRecord, error)
	SoftDeleteEvent(ctx context.Context, id uuid.UUID) error
	ListProposedDatesByEvent(ctx context.Context, eventID uuid.UUID) ([]*ProposedDateRecord, error)
	CreateProposedDates(ctx context.Context, selectedDates []SelectedDate, eventID uuid.UUID) ([]*ProposedDateRecord, error)
	UpdateProposedDate(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*ProposedDateRecord, error)
	DeleteProposedDate(ctx context.Context, id uuid.UUID) error
	CreateProposedDate(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*ProposedDateRecord, error)
}

type EventTransaction interface {
	Do(ctx context.Context, fn func(store EventTxStore) error) error
}

type GoogleCalendarGateway interface {
	FetchEvents(ctx context.Context, userID uuid.UUID, calendars []*CalendarRecord, startTime, endTime time.Time) (*GoogleEventFetchResult, error)
	UpsertEvent(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error)
}
