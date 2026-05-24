package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type SessionQueryOptions struct {
	WithUser bool
}

type SessionRepository interface {
	WithTx(tx transaction.Tx) SessionRepository
	Read(ctx context.Context, id uuid.UUID, opt SessionQueryOptions) (*models.Session, error)
	FindByToken(ctx context.Context, sessionToken string, opt SessionQueryOptions) (*models.Session, error)
	Create(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*models.Session, error)
	UpdateExpiry(ctx context.Context, id uuid.UUID, expiresAt time.Time) (*models.Session, error)
	DeleteByToken(ctx context.Context, sessionToken string) error
}
