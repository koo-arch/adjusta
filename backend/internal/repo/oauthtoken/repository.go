package oauthtoken

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"github.com/koo-arch/adjusta-backend/ent"
)

type OAuthTokenQuertOptions struct {
	AccessToken *string
	RefreshToken *string
	Expiry *time.Time
}

type OAuthTokenRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.OAuthToken, error)
	FindByUserID(ctx context.Context, tx *ent.Tx, userID uuid.UUID) (*ent.OAuthToken, error)
	Create(ctx context.Context, tx *ent.Tx, token *oauth2.Token, entUser *ent.User) (*ent.OAuthToken, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt OAuthTokenQuertOptions) (*ent.OAuthToken, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	SoftDelete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	Restore(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
}