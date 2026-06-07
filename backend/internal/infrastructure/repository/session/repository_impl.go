package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	entSession "github.com/koo-arch/adjusta-backend/ent/session"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
	"github.com/koo-arch/adjusta-backend/internal/repositorymodel"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type SessionRepository = repoSession.SessionRepository
type SessionQueryOptions = repoSession.SessionQueryOptions

type SessionRepositoryImpl struct {
	client *ent.Client
}

func NewSessionRepository(client *ent.Client) *SessionRepositoryImpl {
	return &SessionRepositoryImpl{
		client: client,
	}
}

func (r *SessionRepositoryImpl) WithTx(tx transaction.Tx) SessionRepository {
	return &SessionRepositoryImpl{client: tx.Client()}
}

func (r *SessionRepositoryImpl) Read(ctx context.Context, id uuid.UUID, opt SessionQueryOptions) (*repositorymodel.Session, error) {
	query := r.client.Session.Query()
	if opt.WithUser {
		query = query.WithUser()
	}

	sessionEntity, err := query.
		Where(entSession.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModelSession(sessionEntity), nil
}

func (r *SessionRepositoryImpl) FindByToken(ctx context.Context, sessionToken string, opt SessionQueryOptions) (*repositorymodel.Session, error) {
	query := r.client.Session.Query()
	if opt.WithUser {
		query = query.WithUser()
	}

	sessionEntity, err := query.
		Where(entSession.SessionTokenEQ(sessionToken)).
		Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModelSession(sessionEntity), nil
}

func (r *SessionRepositoryImpl) Create(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repositorymodel.Session, error) {
	create := r.client.Session.Create()
	sessionEntity, err := create.
		SetUserID(userID).
		SetSessionToken(sessionToken).
		SetExpiresAt(expiresAt).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toModelSession(sessionEntity), nil
}

func (r *SessionRepositoryImpl) UpdateExpiry(ctx context.Context, id uuid.UUID, expiresAt time.Time) (*repositorymodel.Session, error) {
	update := r.client.Session.UpdateOneID(id)
	sessionEntity, err := update.
		SetExpiresAt(expiresAt).
		Save(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModelSession(sessionEntity), nil
}

func (r *SessionRepositoryImpl) DeleteByToken(ctx context.Context, sessionToken string) error {
	_, err := r.client.Session.Delete().
		Where(entSession.SessionTokenEQ(sessionToken)).
		Exec(ctx)
	return err
}

func toModelSession(sessionEntity *ent.Session) *repositorymodel.Session {
	if sessionEntity == nil {
		return nil
	}

	model := &repositorymodel.Session{
		ID:           sessionEntity.ID,
		UserID:       sessionEntity.UserID,
		SessionToken: sessionEntity.SessionToken,
		ExpiresAt:    sessionEntity.ExpiresAt,
	}
	if sessionEntity.Edges.User != nil {
		model.User = &repositorymodel.User{
			ID:        sessionEntity.Edges.User.ID,
			Email:     sessionEntity.Edges.User.Email,
			Name:      sessionEntity.Edges.User.Name,
			AvatarURL: sessionEntity.Edges.User.AvatarURL,
		}
	}

	return model
}
