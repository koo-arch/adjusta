package user

import (
	"context"

	"github.com/google/uuid"
)

type UserMutationOptions struct {
	Name      *string
	AvatarURL *string
}

type UserRepository interface {
	Read(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, email string, opt UserMutationOptions) (*User, error)
	Update(ctx context.Context, id uuid.UUID, opt UserMutationOptions) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
}
