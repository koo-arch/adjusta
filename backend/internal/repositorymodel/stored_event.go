package repositorymodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

type StoredProposedDate struct {
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

type StoredEvent struct {
	UserID                 uuid.UUID
	ID                     uuid.UUID
	PrimaryCalendarID      uuid.UUID
	Title                  string
	Location               string
	Description            string
	Status                 domainvalue.EventStatus
	ConfirmedDateID        uuid.UUID
	GoogleEventID          string
	ConfirmedGoogleEventID *string
	SyncStatus             domainvalue.SyncStatus
	LastSyncedAt           *time.Time
	LastSyncError          *string
	Slug                   string
	ProposedDates          []*StoredProposedDate
}
