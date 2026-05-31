package repositorymodel

import "github.com/google/uuid"

type GoogleCalendarInfo struct {
	ID               uuid.UUID
	GoogleCalendarID string
	Summary          string
	IsPrimary        bool
}
