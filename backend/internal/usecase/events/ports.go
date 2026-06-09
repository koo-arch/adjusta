package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
)

type GoogleTokenProvider interface {
	GetToken(ctx context.Context, userID uuid.UUID) (*appmodel.GoogleAuthToken, error)
}

type PrimaryCalendarFinder interface {
	FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*repositorymodel.StoredCalendar, error)
}

type EventReader interface {
	PrimaryCalendarFinder
	ListCalendarsByUser(ctx context.Context, userID uuid.UUID) ([]*repositorymodel.StoredCalendar, error)
	SearchEvents(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*repositorymodel.StoredEvent, error)
	FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error)
}

type EventTxStore interface {
	PrimaryCalendarFinder
	FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error)
	ReadCalendar(ctx context.Context, calendarID uuid.UUID) (*repositorymodel.StoredCalendar, error)
	CreateEvent(ctx context.Context, userID, primaryCalendarID uuid.UUID, title, location, description string, start, end time.Time) (*repositorymodel.StoredEvent, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, opt EventMutation) (*repositorymodel.StoredEvent, error)
	SoftDeleteEvent(ctx context.Context, id uuid.UUID) error
	ListProposedDatesByEvent(ctx context.Context, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error)
	CreateProposedDates(ctx context.Context, selectedDates []appmodel.SelectedDate, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error)
	UpdateProposedDate(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*repositorymodel.StoredProposedDate, error)
	DeleteProposedDate(ctx context.Context, id uuid.UUID) error
	CreateProposedDate(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*repositorymodel.StoredProposedDate, error)
	DecrementPriorityExceptID(ctx context.Context, eventID, excludeID uuid.UUID) error
	ReorderPriority(ctx context.Context, eventID uuid.UUID) error
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
	GoogleEventID          *string
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
	FetchEvents(ctx context.Context, userID uuid.UUID, calendars []*repositorymodel.StoredCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error)
	UpsertEvent(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error)
}
