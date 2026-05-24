package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type UserQueryOptions struct{}

type UserMutationOptions struct {
	Name      *string
	AvatarURL *string
}

type UserRepository interface {
	WithTx(tx transaction.Tx) UserRepository
	Read(ctx context.Context, id uuid.UUID, opt UserQueryOptions) (*models.User, error)
	FindByEmail(ctx context.Context, email string, opt UserQueryOptions) (*models.User, error)
	Create(ctx context.Context, email string, opt UserMutationOptions) (*models.User, error)
	Update(ctx context.Context, id uuid.UUID, opt UserMutationOptions) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
}
