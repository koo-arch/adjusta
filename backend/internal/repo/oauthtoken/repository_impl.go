package oauthtoken

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/user"
	"github.com/koo-arch/adjusta-backend/ent/oauthtoken"
)

type OAuthTokenRepositoryImpl struct {
	client *ent.Client
}

func NewOAuthTokenRepository(client *ent.Client) *OAuthTokenRepositoryImpl {
	return &OAuthTokenRepositoryImpl{
		client: client,
	}
}

func (r *OAuthTokenRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.OAuthToken, error) {
	if tx != nil {
		return tx.OAuthToken.Get(ctx, id)
	}
	return r.client.OAuthToken.Get(ctx, id)
}

func (r *OAuthTokenRepositoryImpl) FindByUserID(ctx context.Context, tx *ent.Tx, userID uuid.UUID) (*ent.OAuthToken, error) {
	findOAuthToken := r.client.OAuthToken.Query()
	if tx != nil {
		findOAuthToken = tx.OAuthToken.Query()
	}
	return findOAuthToken.
		Where(oauthtoken.HasUserWith(user.ID(userID))).
		Only(ctx)
}

func (r *OAuthTokenRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, token *oauth2.Token, entUser *ent.User) (*ent.OAuthToken, error) {
	oauthTokenCreate := r.client.OAuthToken.Create()
	if tx != nil {
		oauthTokenCreate = tx.OAuthToken.Create()
	}
	return oauthTokenCreate.
		SetAccessToken(token.AccessToken).
		SetRefreshToken(token.RefreshToken).
		SetExpiry(token.Expiry).
		SetUser(entUser).
		Save(ctx)
}

func (r *OAuthTokenRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt OAuthTokenQuertOptions) (*ent.OAuthToken, error) {
	oauthTokenUpdate := r.client.OAuthToken.UpdateOneID(id)
	if tx != nil {
		oauthTokenUpdate = tx.OAuthToken.UpdateOneID(id)
	}

	if opt.AccessToken != nil {
		oauthTokenUpdate.SetAccessToken(*opt.AccessToken)
	}

	if opt.RefreshToken != nil {
		oauthTokenUpdate.SetRefreshToken(*opt.RefreshToken)
	}

	if opt.Expiry != nil {
		oauthTokenUpdate.SetExpiry(*opt.Expiry)
	}

	return oauthTokenUpdate.Save(ctx)

}

func (r *OAuthTokenRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.OAuthToken.DeleteOneID(id).Exec(ctx)
	}
	return r.client.OAuthToken.DeleteOneID(id).Exec(ctx)
}

func (r *OAuthTokenRepositoryImpl) SoftDelete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	softDeleteOAuthToken := r.client.OAuthToken.UpdateOneID(id)
	if tx != nil {
		softDeleteOAuthToken = tx.OAuthToken.UpdateOneID(id)
	}
	return softDeleteOAuthToken.
		SetDeletedAt(time.Now()).
		Exec(ctx)
}

func (r *OAuthTokenRepositoryImpl) Restore(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	restoreOAuthToken := r.client.OAuthToken.UpdateOneID(id)
	if tx != nil {
		restoreOAuthToken = tx.OAuthToken.UpdateOneID(id)
	}
	return restoreOAuthToken.
		SetNillableDeletedAt(nil).
		Exec(ctx)
}