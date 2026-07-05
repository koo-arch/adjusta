package account

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AccountMutationOptions struct {
	GoogleUserID *string
	AccessToken  *string
	RefreshToken *string
	ExpiresAt    *time.Time
	Scope        *string
}

type AccountRepository interface {
	Read(ctx context.Context, id uuid.UUID) (*Account, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*Account, error)
	Create(ctx context.Context, userID uuid.UUID, opt AccountMutationOptions) (*Account, error)
	Update(ctx context.Context, id uuid.UUID, opt AccountMutationOptions) (*Account, error)
}
