package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	entAccount "github.com/koo-arch/adjusta-backend/ent/account"
)

type AccountRepositoryImpl struct {
	client *ent.Client
}

func NewAccountRepository(client *ent.Client) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{
		client: client,
	}
}

func (r *AccountRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.Account, error) {
	if tx != nil {
		return tx.Account.Get(ctx, id)
	}
	return r.client.Account.Get(ctx, id)
}

func (r *AccountRepositoryImpl) FindByUserID(ctx context.Context, tx *ent.Tx, userID uuid.UUID) (*ent.Account, error) {
	query := r.client.Account.Query()
	if tx != nil {
		query = tx.Account.Query()
	}

	return query.
		Where(entAccount.UserIDEQ(userID)).
		Only(ctx)
}

func (r *AccountRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, userID uuid.UUID, opt AccountMutationOptions) (*ent.Account, error) {
	create := r.client.Account.Create()
	if tx != nil {
		create = tx.Account.Create()
	}

	create.SetUserID(userID)
	applyAccountCreateOptions(create, opt)

	return create.Save(ctx)
}

func (r *AccountRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt AccountMutationOptions) (*ent.Account, error) {
	update := r.client.Account.UpdateOneID(id)
	if tx != nil {
		update = tx.Account.UpdateOneID(id)
	}

	applyAccountUpdateOptions(update, opt)

	return update.Save(ctx)
}

func applyAccountCreateOptions(create *ent.AccountCreate, opt AccountMutationOptions) {
	if opt.GoogleUserID != nil {
		create.SetGoogleUserID(*opt.GoogleUserID)
	}
	create.SetNillableAccessToken(opt.AccessToken)
	create.SetNillableRefreshToken(opt.RefreshToken)
	create.SetNillableExpiresAt(opt.ExpiresAt)
	create.SetNillableScope(opt.Scope)
}

func applyAccountUpdateOptions(update *ent.AccountUpdateOne, opt AccountMutationOptions) {
	if opt.GoogleUserID != nil {
		update.SetGoogleUserID(*opt.GoogleUserID)
	}
	update.SetNillableAccessToken(opt.AccessToken)
	update.SetNillableRefreshToken(opt.RefreshToken)
	update.SetNillableExpiresAt(opt.ExpiresAt)
	update.SetNillableScope(opt.Scope)
}
