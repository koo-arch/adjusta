package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	GoogleUserID string
	AccessToken  *string
	RefreshToken *string
	ExpiresAt    *time.Time
	Scope        *string
}
