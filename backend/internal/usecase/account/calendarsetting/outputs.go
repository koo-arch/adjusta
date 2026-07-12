package calendarsetting

import (
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type CalendarSettingOutput struct {
	ID                uuid.UUID
	CalendarID        uuid.UUID
	GoogleCalendarID  string
	Summary           string
	Description       *string
	Timezone          *string
	Role              value.UserCalendarRole
	IsVisible         bool
	SyncProposedDates bool
}

type CalendarSettingUpdateRequest struct {
	Role              *value.UserCalendarRole
	IsVisible         *bool
	SyncProposedDates *bool
}

type CandidateSyncSettingOutput struct {
	Enabled  bool
	Calendar *CalendarSettingOutput
}
