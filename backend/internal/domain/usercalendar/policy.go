package usercalendar

import "github.com/koo-arch/adjusta-backend/internal/domainvalue"

const AdjustaCandidateCalendarSummary = "Adjusta 候補日程"

func ExternalSyncRole(isPrimary bool) domainvalue.UserCalendarRole {
	if isPrimary {
		return domainvalue.UserCalendarRolePrimary
	}
	return domainvalue.UserCalendarRoleReference
}

func IsExternalSyncRole(role domainvalue.UserCalendarRole) bool {
	return role == domainvalue.UserCalendarRolePrimary || role == domainvalue.UserCalendarRoleReference
}

func IsAdjustaCandidateCalendarSummary(summary string) bool {
	return summary == AdjustaCandidateCalendarSummary
}
