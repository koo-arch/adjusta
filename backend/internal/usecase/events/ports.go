package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

type GoogleTokenProvider interface {
	GetToken(ctx context.Context, userID uuid.UUID) (*appmodel.GoogleAuthToken, error)
}

type CalendarRecord struct {
	ID                uuid.UUID
	GoogleCalendarID  string
	Summary           string
	Description       *string
	Timezone          *string
	SyncProposedDates bool
}

type ProposedDateRecord struct {
	ID            uuid.UUID
	EventID       *uuid.UUID
	GoogleEventID *string
	StartTime     time.Time
	EndTime       time.Time
	Priority      int
	Status        domainvalue.ProposedDateStatus
	SyncStatus    domainvalue.SyncStatus
	LastSyncedAt  *time.Time
	LastSyncError *string
}

type EventRecord struct {
	ID                     uuid.UUID
	PrimaryCalendarID      uuid.UUID
	Title                  string
	Location               string
	Description            string
	Status                 domainvalue.EventStatus
	ConfirmedDateID        uuid.UUID
	ConfirmedGoogleEventID *string
	SyncStatus             domainvalue.SyncStatus
	LastSyncedAt           *time.Time
	LastSyncError          *string
	Slug                   string
	ProposedDates          []*ProposedDateRecord
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
	FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*EventRecord, error)
}

type EventTxStore interface {
	PrimaryCalendarFinder
	AdjustaCandidateCalendarFinder
	FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*EventRecord, error)
	ReadCalendar(ctx context.Context, calendarID uuid.UUID) (*CalendarRecord, error)
	CreateEvent(ctx context.Context, userID, primaryCalendarID uuid.UUID, title, location, description string, start, end time.Time) (*EventRecord, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, opt EventMutation) (*EventRecord, error)
	SoftDeleteEvent(ctx context.Context, id uuid.UUID) error
	ListProposedDatesByEvent(ctx context.Context, eventID uuid.UUID) ([]*ProposedDateRecord, error)
	CreateProposedDates(ctx context.Context, selectedDates []appmodel.SelectedDate, eventID uuid.UUID) ([]*ProposedDateRecord, error)
	UpdateProposedDate(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*ProposedDateRecord, error)
	DeleteProposedDate(ctx context.Context, id uuid.UUID) error
	CreateProposedDate(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*ProposedDateRecord, error)
}

type EventTransaction interface {
	Do(ctx context.Context, fn func(store EventTxStore) error) error
}

type EventSearchOptions struct {
	WithProposedDates bool
	Title             *string
	Location          *string
	Description       *string
	Status            *domainvalue.EventStatus
	StartTimeGTE      *time.Time
	StartTimeLTE      *time.Time
	EndTimeGTE        *time.Time
	EndTimeLTE        *time.Time
	SortBy            string
	SortOrder         string
}

type EventMutation struct {
	Title                  *string
	Location               *string
	Description            *string
	Status                 *domainvalue.EventStatus
	SyncStatus             *domainvalue.SyncStatus
	ConfirmedDateID        *uuid.UUID
	ConfirmedGoogleEventID *string
	LastSyncedAt           *time.Time
	ClearLastSyncedAt      bool
	LastSyncError          *string
	ClearLastSyncError     bool
}

type ProposedDateMutation struct {
	GoogleEventID      *string
	Start              *time.Time
	End                *time.Time
	Priority           *int
	Status             *domainvalue.ProposedDateStatus
	SyncStatus         *domainvalue.SyncStatus
	LastSyncedAt       *time.Time
	ClearLastSyncedAt  bool
	LastSyncError      *string
	ClearLastSyncError bool
}

type GoogleEventFetchResult struct {
	Events          []*appmodel.GoogleEvent
	FailedCalendars []string
}

type GoogleCalendarGateway interface {
	FetchEvents(ctx context.Context, userID uuid.UUID, calendars []*CalendarRecord, startTime, endTime time.Time) (*GoogleEventFetchResult, error)
	UpsertEvent(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error)
}
