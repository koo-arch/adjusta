package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

type EventDraftCreation struct {
	Title         string         `json:"title" binding:"required"`
	Location      string         `json:"location"`
	Description   string         `json:"description"`
	SelectedDates []SelectedDate `json:"selected_dates" binding:"required"`
}

type SelectedDate struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Priority int       `json:"priority"`
}

type EventDraftUpdate struct {
	Title           string                  `json:"title"`
	Location        string                  `json:"location"`
	Description     string                  `json:"description"`
	Status          domainvalue.EventStatus `json:"status"`
	ConfirmedDateID *uuid.UUID              `json:"confirmed_date_id"`
	GoogleEventID   string                  `json:"google_event_id"`
	ProposedDates   []ProposedDate          `json:"proposed_dates"`
}

type ProposedDate struct {
	ID            *uuid.UUID                     `json:"id"`
	GoogleEventID *string                        `json:"google_event_id,omitempty"`
	Start         *time.Time                     `json:"start"`
	End           *time.Time                     `json:"end"`
	Priority      int                            `json:"priority"`
	Status        domainvalue.ProposedDateStatus `json:"status"`
	SyncStatus    domainvalue.SyncStatus         `json:"sync_status"`
	LastSyncedAt  *time.Time                     `json:"last_synced_at,omitempty"`
	LastSyncError *string                        `json:"last_sync_error,omitempty"`
}

type ConfirmEvent struct {
	ConfirmDate ConfirmDate `json:"confirm_date" binding:"required"`
}

type ConfirmDate struct {
	ID            *uuid.UUID `json:"id"`
	GoogleEventID string     `json:"google_event_id"`
	Start         *time.Time `json:"start"`
	End           *time.Time `json:"end"`
	Priority      int        `json:"priority"`
}

type EventDraftDetail struct {
	ID                     uuid.UUID               `json:"id" binding:"required"`
	Title                  string                  `json:"title"`
	Location               string                  `json:"location"`
	Description            string                  `json:"description"`
	Status                 domainvalue.EventStatus `json:"status"`
	SyncStatus             domainvalue.SyncStatus  `json:"sync_status"`
	ConfirmedDateID        *uuid.UUID              `json:"confirmed_date_id"`
	GoogleEventID          string                  `json:"google_event_id"`
	ConfirmedGoogleEventID *string                 `json:"confirmed_google_event_id,omitempty"`
	LastSyncedAt           *time.Time              `json:"last_synced_at,omitempty"`
	LastSyncError          *string                 `json:"last_sync_error,omitempty"`
	ProposedDates          []ProposedDate          `json:"proposed_dates"`
}

type UpcomingEvent struct {
	ID                     uuid.UUID               `json:"id" binding:"required"`
	Title                  string                  `json:"title"`
	Location               string                  `json:"location"`
	Description            string                  `json:"description"`
	Status                 domainvalue.EventStatus `json:"status"`
	SyncStatus             domainvalue.SyncStatus  `json:"sync_status"`
	ConfirmedDateID        uuid.UUID               `json:"confirmed_date_id"`
	GoogleEventID          string                  `json:"google_event_id"`
	ConfirmedGoogleEventID *string                 `json:"confirmed_google_event_id,omitempty"`
	LastSyncedAt           *time.Time              `json:"last_synced_at,omitempty"`
	LastSyncError          *string                 `json:"last_sync_error,omitempty"`
	Start                  time.Time               `json:"start"`
	End                    time.Time               `json:"end"`
}

type NeedsActionDraft struct {
	ID             uuid.UUID               `json:"id" binding:"required"`
	Title          string                  `json:"title"`
	Location       string                  `json:"location"`
	Description    string                  `json:"description"`
	Status         domainvalue.EventStatus `json:"status"`
	Start          time.Time               `json:"start"`
	End            time.Time               `json:"end"`
	NeedsAttention bool                    `json:"needs_attention"`
}
