package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
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
	Title           string            `json:"title"`
	Location        string            `json:"location"`
	Description     string            `json:"description"`
	Status          value.EventStatus `json:"status"`
	ConfirmedDateID *uuid.UUID        `json:"confirmed_date_id"`
	GoogleEventID   string            `json:"google_event_id"`
	ProposedDates   []ProposedDate    `json:"proposed_dates"`
}

type ProposedDate struct {
	ID            *uuid.UUID               `json:"id"`
	GoogleEventID *string                  `json:"google_event_id,omitempty"`
	Start         *time.Time               `json:"start"`
	End           *time.Time               `json:"end"`
	Priority      int                      `json:"priority"`
	Status        value.ProposedDateStatus `json:"status"`
	SyncStatus    value.SyncStatus         `json:"sync_status"`
	LastSyncedAt  *time.Time               `json:"last_synced_at,omitempty"`
	LastSyncError *string                  `json:"last_sync_error,omitempty"`
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
	ID                     uuid.UUID         `json:"id" binding:"required"`
	Title                  string            `json:"title"`
	Location               string            `json:"location"`
	Description            string            `json:"description"`
	Status                 value.EventStatus `json:"status"`
	SyncStatus             value.SyncStatus  `json:"sync_status"`
	ConfirmedDateID        *uuid.UUID        `json:"confirmed_date_id"`
	GoogleEventID          string            `json:"google_event_id"`
	ConfirmedGoogleEventID *string           `json:"confirmed_google_event_id,omitempty"`
	LastSyncedAt           *time.Time        `json:"last_synced_at,omitempty"`
	LastSyncError          *string           `json:"last_sync_error,omitempty"`
	ProposedDates          []ProposedDate    `json:"proposed_dates"`
}

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type EventDraftList struct {
	Items      []*EventDraftDetail `json:"items"`
	Pagination Pagination          `json:"pagination"`
}

type UpcomingEvent struct {
	ID                     uuid.UUID         `json:"id" binding:"required"`
	Title                  string            `json:"title"`
	Location               string            `json:"location"`
	Description            string            `json:"description"`
	Status                 value.EventStatus `json:"status"`
	SyncStatus             value.SyncStatus  `json:"sync_status"`
	ConfirmedDateID        uuid.UUID         `json:"confirmed_date_id"`
	GoogleEventID          string            `json:"google_event_id"`
	ConfirmedGoogleEventID *string           `json:"confirmed_google_event_id,omitempty"`
	LastSyncedAt           *time.Time        `json:"last_synced_at,omitempty"`
	LastSyncError          *string           `json:"last_sync_error,omitempty"`
	Start                  time.Time         `json:"start"`
	End                    time.Time         `json:"end"`
}

type NeedsActionDraft struct {
	ID             uuid.UUID         `json:"id" binding:"required"`
	Title          string            `json:"title"`
	Location       string            `json:"location"`
	Description    string            `json:"description"`
	Status         value.EventStatus `json:"status"`
	Start          time.Time         `json:"start"`
	End            time.Time         `json:"end"`
	NeedsAttention bool              `json:"needs_attention"`
}

type GoogleEvent struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Location    string `json:"location"`
	ColorID     string `json:"color"`
	Start       string `json:"start"`
	End         string `json:"end"`
}
