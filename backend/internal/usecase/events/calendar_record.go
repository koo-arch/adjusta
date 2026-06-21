package events

import (
	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
)

type CalendarRecord struct {
	ID                uuid.UUID
	GoogleCalendarID  string
	Summary           string
	Description       *string
	Timezone          *string
	SyncProposedDates bool
}

func toCalendarRecord(calendar *repoCalendar.Calendar) *CalendarRecord {
	return toCalendarRecordWithSync(calendar, false)
}

func toCalendarRecordWithSync(calendar *repoCalendar.Calendar, syncProposedDates bool) *CalendarRecord {
	if calendar == nil {
		return nil
	}

	return &CalendarRecord{
		ID:                calendar.ID,
		GoogleCalendarID:  calendar.GoogleCalendarID,
		Summary:           calendar.Summary,
		Description:       calendar.Description,
		Timezone:          calendar.Timezone,
		SyncProposedDates: syncProposedDates,
	}
}

func toCalendarRecords(calendars []*repoCalendar.Calendar) []*CalendarRecord {
	records := make([]*CalendarRecord, 0, len(calendars))
	for _, calendar := range calendars {
		records = append(records, toCalendarRecord(calendar))
	}
	return records
}
