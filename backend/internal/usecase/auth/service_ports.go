package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
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
	FindUserByEmail(ctx context.Context, email string) (*repoUser.User, error)
	FindAccountByUserID(ctx context.Context, userID uuid.UUID) (*repoAccount.Account, error)
}

type SignInStore interface {
	CreateUser(ctx context.Context, email string, opt UserMutation) (*repoUser.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, opt UserMutation) (*repoUser.User, error)
	CreateAccount(ctx context.Context, userID uuid.UUID, opt AccountMutation) (*repoAccount.Account, error)
	UpdateAccount(ctx context.Context, accountID uuid.UUID, opt AccountMutation) (*repoAccount.Account, error)
}

type SignInTransaction interface {
	Do(ctx context.Context, fn func(store SignInStore) error) error
}

type SessionStore interface {
	CreateSession(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error)
	FindSessionByToken(ctx context.Context, sessionToken string, withUser bool) (*repoSession.Session, error)
	UpdateSessionExpiry(ctx context.Context, sessionID uuid.UUID, expiresAt time.Time) (*repoSession.Session, error)
	DeleteSessionByToken(ctx context.Context, sessionToken string) error
}
