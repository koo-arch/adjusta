package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
)

type UserQueryOptions struct {
	WithOAuthToken bool
}

type UserMutationOptions struct {
	Name               *string
	AvatarURL          *string
	RefreshToken       *string
	RefreshTokenExpiry *time.Time
}

type UserRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt UserQueryOptions) (*ent.User, error)
	FindByEmail(ctx context.Context, tx *ent.Tx, email string, opt UserQueryOptions) (*ent.User, error)
	Create(ctx context.Context, tx *ent.Tx, email string, opt UserMutationOptions) (*ent.User, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt UserMutationOptions) (*ent.User, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	SoftDelete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	Restore(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	SoftDeleteWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	RestoreWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
}
