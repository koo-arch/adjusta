package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	internalModels "github.com/koo-arch/adjusta-backend/internal/models"
)

type UserMutation struct {
	Name      *string
	AvatarURL *string
}

type AccountMutation struct {
	GoogleUserID *string
	AccessToken  *string
	RefreshToken *string
	ExpiresAt    *time.Time
	Scope        *string
}

type SignInReader interface {
	FindUserByEmail(ctx context.Context, email string) (*internalModels.User, error)
	FindAccountByUserID(ctx context.Context, userID uuid.UUID) (*internalModels.Account, error)
}

type SignInStore interface {
	CreateUser(ctx context.Context, email string, opt UserMutation) (*internalModels.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, opt UserMutation) (*internalModels.User, error)
	CreateAccount(ctx context.Context, userID uuid.UUID, opt AccountMutation) (*internalModels.Account, error)
	UpdateAccount(ctx context.Context, accountID uuid.UUID, opt AccountMutation) (*internalModels.Account, error)
}

type SignInTransaction interface {
	Do(ctx context.Context, fn func(store SignInStore) error) error
}

type SessionStore interface {
	CreateSession(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*internalModels.Session, error)
	FindSessionByToken(ctx context.Context, sessionToken string, withUser bool) (*internalModels.Session, error)
	UpdateSessionExpiry(ctx context.Context, sessionID uuid.UUID, expiresAt time.Time) (*internalModels.Session, error)
	DeleteSessionByToken(ctx context.Context, sessionToken string) error
}
