package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"golang.org/x/oauth2"
)

type GoogleTokenProvider interface {
	GetToken(ctx context.Context, userID uuid.UUID) (*oauth2.Token, error)
}

type PrimaryCalendarFinder interface {
	FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*models.StoredCalendar, error)
}

type EventReader interface {
	PrimaryCalendarFinder
	ListGoogleCalendarInfosByUser(ctx context.Context, userID uuid.UUID) ([]*models.GoogleCalendarInfo, error)
	SearchEvents(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*models.StoredEvent, error)
	FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*models.StoredEvent, error)
}

type EventTxStore interface {
	PrimaryCalendarFinder
	FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*models.StoredEvent, error)
	CreateEvent(ctx context.Context, calendarID uuid.UUID, title, location, description string, start, end time.Time) (*models.StoredEvent, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, opt EventMutation) (*models.StoredEvent, error)
	SoftDeleteEvent(ctx context.Context, id uuid.UUID) error
	ListProposedDatesByEvent(ctx context.Context, eventID uuid.UUID) ([]*models.StoredProposedDate, error)
	CreateProposedDates(ctx context.Context, selectedDates []models.SelectedDate, eventID uuid.UUID) ([]*models.StoredProposedDate, error)
	UpdateProposedDate(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*models.StoredProposedDate, error)
	DeleteProposedDate(ctx context.Context, id uuid.UUID) error
	CreateProposedDate(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*models.StoredProposedDate, error)
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
	Status            *models.EventStatus
	StartTimeGTE      *time.Time
	StartTimeLTE      *time.Time
	EndTimeGTE        *time.Time
	EndTimeLTE        *time.Time
	SortBy            string
	SortOrder         string
}

type EventMutation struct {
	Title           *string
	Location        *string
	Description     *string
	Status          *models.EventStatus
	ConfirmedDateID *uuid.UUID
	GoogleEventID   *string
}

type ProposedDateMutation struct {
	Start    *time.Time
	End      *time.Time
	Priority *int
}

type GoogleEventFetchResult struct {
	Events          []*models.GoogleEvent
	FailedCalendars []string
}

type GoogleCalendarGateway interface {
	FetchEvents(ctx context.Context, userID uuid.UUID, calendars []*models.GoogleCalendarInfo, startTime, endTime time.Time) (*GoogleEventFetchResult, error)
	UpsertEvent(ctx context.Context, userID uuid.UUID, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error)
}
