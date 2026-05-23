package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
)

type SessionQueryOptions struct {
	WithUser bool
}

type SessionRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt SessionQueryOptions) (*ent.Session, error)
	FindByToken(ctx context.Context, tx *ent.Tx, sessionToken string, opt SessionQueryOptions) (*ent.Session, error)
	Create(ctx context.Context, tx *ent.Tx, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*ent.Session, error)
	UpdateExpiry(ctx context.Context, tx *ent.Tx, id uuid.UUID, expiresAt time.Time) (*ent.Session, error)
	DeleteByToken(ctx context.Context, tx *ent.Tx, sessionToken string) error
}
