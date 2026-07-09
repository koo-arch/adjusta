package dto

import (
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type CalendarSetting struct {
	ID                uuid.UUID              `json:"id"`
	CalendarID        uuid.UUID              `json:"calendar_id"`
	GoogleCalendarID  string                 `json:"google_calendar_id"`
	Summary           string                 `json:"summary"`
	Description       *string                `json:"description,omitempty"`
	Timezone          *string                `json:"timezone,omitempty"`
	Role              value.UserCalendarRole `json:"role"`
	IsVisible         bool                   `json:"is_visible"`
	SyncProposedDates bool                   `json:"sync_proposed_dates"`
}

type CalendarSettingUpdate struct {
	Role              *value.UserCalendarRole `json:"role"`
	IsVisible         *bool                   `json:"is_visible"`
	SyncProposedDates *bool                   `json:"sync_proposed_dates"`
}
