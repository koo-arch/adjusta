package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	entSession "github.com/koo-arch/adjusta-backend/ent/session"
)

type SessionRepositoryImpl struct {
	client *ent.Client
}

func NewSessionRepository(client *ent.Client) *SessionRepositoryImpl {
	return &SessionRepositoryImpl{
		client: client,
	}
}

func (r *SessionRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt SessionQueryOptions) (*ent.Session, error) {
	query := r.client.Session.Query()
	if tx != nil {
		query = tx.Session.Query()
	}

	if opt.WithUser {
		query = query.WithUser()
	}

	return query.
		Where(entSession.IDEQ(id)).
		Only(ctx)
}

func (r *SessionRepositoryImpl) FindByToken(ctx context.Context, tx *ent.Tx, sessionToken string, opt SessionQueryOptions) (*ent.Session, error) {
	query := r.client.Session.Query()
	if tx != nil {
		query = tx.Session.Query()
	}

	if opt.WithUser {
		query = query.WithUser()
	}

	return query.
		Where(entSession.SessionTokenEQ(sessionToken)).
		Only(ctx)
}

func (r *SessionRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*ent.Session, error) {
	create := r.client.Session.Create()
	if tx != nil {
		create = tx.Session.Create()
	}

	return create.
		SetUserID(userID).
		SetSessionToken(sessionToken).
		SetExpiresAt(expiresAt).
		Save(ctx)
}

func (r *SessionRepositoryImpl) UpdateExpiry(ctx context.Context, tx *ent.Tx, id uuid.UUID, expiresAt time.Time) (*ent.Session, error) {
	update := r.client.Session.UpdateOneID(id)
	if tx != nil {
		update = tx.Session.UpdateOneID(id)
	}

	return update.
		SetExpiresAt(expiresAt).
		Save(ctx)
}

func (r *SessionRepositoryImpl) DeleteByToken(ctx context.Context, tx *ent.Tx, sessionToken string) error {
	deleteQuery := r.client.Session.Delete()
	if tx != nil {
		deleteQuery = tx.Session.Delete()
	}

	_, err := deleteQuery.
		Where(entSession.SessionTokenEQ(sessionToken)).
		Exec(ctx)
	return err
}
