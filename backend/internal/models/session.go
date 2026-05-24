package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	SessionToken string
	ExpiresAt    time.Time
	User         *User
}
