package usercalendar

import "github.com/koo-arch/adjusta-backend/internal/domain/value"

const AdjustaCandidateCalendarSummary = "Adjusta 候補日程"

func ExternalSyncRole(isPrimary bool) value.UserCalendarRole {
	if isPrimary {
		return value.UserCalendarRolePrimary
	}
	return value.UserCalendarRoleReference
}

func IsExternalSyncRole(role value.UserCalendarRole) bool {
	return role == value.UserCalendarRolePrimary || role == value.UserCalendarRoleReference
}

func IsAdjustaCandidateCalendarSummary(summary string) bool {
	return summary == AdjustaCandidateCalendarSummary
}
