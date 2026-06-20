package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

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
