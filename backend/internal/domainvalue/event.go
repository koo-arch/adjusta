package domainvalue

type EventStatus string

const (
	StatusDraft     EventStatus = "draft"
	StatusActive    EventStatus = "active"
	StatusConfirmed EventStatus = "confirmed"
	StatusCancelled EventStatus = "cancelled"
)

type UserCalendarRole string

const (
	UserCalendarRolePrimary          UserCalendarRole = "primary"
	UserCalendarRoleAdjustaCandidate UserCalendarRole = "adjusta_candidate"
	UserCalendarRoleReference        UserCalendarRole = "reference"
)

type ProposedDateStatus string

const (
	ProposedDateStatusActive      ProposedDateStatus = "active"
	ProposedDateStatusConfirmed   ProposedDateStatus = "confirmed"
	ProposedDateStatusNotSelected ProposedDateStatus = "not_selected"
	ProposedDateStatusCancelled   ProposedDateStatus = "cancelled"
)

type SyncStatus string

const (
	SyncStatusNotSynced SyncStatus = "not_synced"
	SyncStatusPending   SyncStatus = "pending_sync"
	SyncStatusSynced    SyncStatus = "synced"
	SyncStatusFailed    SyncStatus = "sync_failed"
)
