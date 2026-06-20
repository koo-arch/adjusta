package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

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
	EventID       uuid.UUID
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
	ProposedDates          []*ProposedDateRecord
}
