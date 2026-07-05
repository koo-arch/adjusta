package account

import (
	"context"

	"github.com/google/uuid"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent"
	entAccount "github.com/koo-arch/adjusta-backend/internal/infrastructure/ent/account"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
)

type AccountRepository = repoAccount.AccountRepository
type AccountMutationOptions = repoAccount.AccountMutationOptions

type AccountRepositoryImpl struct {
	client *ent.Client
}

func NewAccountRepository(client *ent.Client) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{
		client: client,
	}
}

func (r *AccountRepositoryImpl) Read(ctx context.Context, id uuid.UUID) (*repoAccount.Account, error) {
	accountEntity, err := r.client.Account.Get(ctx, id)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModel(accountEntity), nil
}

func (r *AccountRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID) (*repoAccount.Account, error) {
	accountEntity, err := r.client.Account.Query().
		Where(entAccount.UserIDEQ(userID)).
		Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModel(accountEntity), nil
}

func (r *AccountRepositoryImpl) Create(ctx context.Context, userID uuid.UUID, opt AccountMutationOptions) (*repoAccount.Account, error) {
	create := r.client.Account.Create()
	create.SetUserID(userID)
	applyAccountCreateOptions(create, opt)
	accountEntity, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}
	return toModel(accountEntity), nil
}

func (r *AccountRepositoryImpl) Update(ctx context.Context, id uuid.UUID, opt AccountMutationOptions) (*repoAccount.Account, error) {
	update := r.client.Account.UpdateOneID(id)
	applyAccountUpdateOptions(update, opt)
	accountEntity, err := update.Save(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModel(accountEntity), nil
}

func toModel(accountEntity *ent.Account) *repoAccount.Account {
	if accountEntity == nil {
		return nil
	}

	return &repoAccount.Account{
		ID:           accountEntity.ID,
		UserID:       accountEntity.UserID,
		GoogleUserID: accountEntity.GoogleUserID,
		AccessToken:  accountEntity.AccessToken,
		RefreshToken: accountEntity.RefreshToken,
		ExpiresAt:    accountEntity.ExpiresAt,
		Scope:        accountEntity.Scope,
	}
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
