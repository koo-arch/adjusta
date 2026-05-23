package account

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
)

type AccountMutationOptions struct {
	GoogleUserID *string
	AccessToken  *string
	RefreshToken *string
	ExpiresAt    *time.Time
	Scope        *string
}

type AccountRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.Account, error)
	FindByUserID(ctx context.Context, tx *ent.Tx, userID uuid.UUID) (*ent.Account, error)
	Create(ctx context.Context, tx *ent.Tx, userID uuid.UUID, opt AccountMutationOptions) (*ent.Account, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt AccountMutationOptions) (*ent.Account, error)
}
