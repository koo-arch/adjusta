package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	entSession "github.com/koo-arch/adjusta-backend/ent/session"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
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

func (r *SessionRepositoryImpl) Read(ctx context.Context, id uuid.UUID, opt SessionQueryOptions) (*repoSession.Session, error) {
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

func (r *SessionRepositoryImpl) FindByToken(ctx context.Context, sessionToken string, opt SessionQueryOptions) (*repoSession.Session, error) {
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

func (r *SessionRepositoryImpl) Create(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error) {
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

func (r *SessionRepositoryImpl) UpdateExpiry(ctx context.Context, id uuid.UUID, expiresAt time.Time) (*repoSession.Session, error) {
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

func toModelSession(sessionEntity *ent.Session) *repoSession.Session {
	if sessionEntity == nil {
		return nil
	}

	model := &repoSession.Session{
		ID:           sessionEntity.ID,
		UserID:       sessionEntity.UserID,
		SessionToken: sessionEntity.SessionToken,
		ExpiresAt:    sessionEntity.ExpiresAt,
	}
	if sessionEntity.Edges.User != nil {
		model.User = &repoUser.User{
			ID:        sessionEntity.Edges.User.ID,
			Email:     sessionEntity.Edges.User.Email,
			Name:      sessionEntity.Edges.User.Name,
			AvatarURL: sessionEntity.Edges.User.AvatarURL,
		}
	}

	return model
}
