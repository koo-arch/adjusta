package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type ProposedDateOutput struct {
	ID            *uuid.UUID
	GoogleEventID *string
	Start         *time.Time
	End           *time.Time
	Priority      int
	Status        value.ProposedDateStatus
	SyncStatus    value.SyncStatus
	LastSyncedAt  *time.Time
	LastSyncError *string
}

type EventDraftDetailOutput struct {
	ID                     uuid.UUID
	Title                  string
	Location               string
	Description            string
	Status                 value.EventStatus
	SyncStatus             value.SyncStatus
	ConfirmedDateID        *uuid.UUID
	GoogleEventID          string
	ConfirmedGoogleEventID *string
	LastSyncedAt           *time.Time
	LastSyncError          *string
	ProposedDates          []ProposedDateOutput
}

type UpcomingEventOutput struct {
	ID                     uuid.UUID
	Title                  string
	Location               string
	Description            string
	Status                 value.EventStatus
	SyncStatus             value.SyncStatus
	ConfirmedDateID        uuid.UUID
	GoogleEventID          string
	ConfirmedGoogleEventID *string
	LastSyncedAt           *time.Time
	LastSyncError          *string
	Start                  time.Time
	End                    time.Time
}

type NeedsActionDraftOutput struct {
	ID             uuid.UUID
	Title          string
	Location       string
	Description    string
	Status         value.EventStatus
	Start          time.Time
	End            time.Time
	NeedsAttention bool
}

type FetchedGoogleEvent struct {
	ID          string
	Summary     string
	Description string
	Location    string
	ColorID     string
	Start       string
	End         string
}

type GoogleEventFetchResult struct {
	Events          []*FetchedGoogleEvent
	FailedCalendars []string
}
