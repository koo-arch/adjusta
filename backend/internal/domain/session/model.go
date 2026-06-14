package session

import (
	"time"

	"github.com/google/uuid"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
)

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	SessionToken string
	ExpiresAt    time.Time
	User         *repoUser.User
}
