package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/user"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type UserRepository = repoUser.UserRepository
type UserQueryOptions = repoUser.UserQueryOptions
type UserMutationOptions = repoUser.UserMutationOptions

type UserRepositoryImpl struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		client: client,
	}
}

func (r *UserRepositoryImpl) WithTx(tx transaction.Tx) UserRepository {
	return &UserRepositoryImpl{client: tx.Client()}
}

func (r *UserRepositoryImpl) Read(ctx context.Context, id uuid.UUID, opt UserQueryOptions) (*repoUser.User, error) {
	userEntity, err := r.client.User.Query().
		Where(user.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModelUser(userEntity), nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string, opt UserQueryOptions) (*repoUser.User, error) {
	userEntity, err := r.client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModelUser(userEntity), nil
}

func (r *UserRepositoryImpl) Create(ctx context.Context, email string, opt UserMutationOptions) (*repoUser.User, error) {
	userCreate := r.client.User.Create()
	userCreate.SetEmail(email)
	applyUserCreateOptions(userCreate, opt)
	userEntity, err := userCreate.Save(ctx)
	if err != nil {
		return nil, err
	}
	return toModelUser(userEntity), nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, id uuid.UUID, opt UserMutationOptions) (*repoUser.User, error) {
	userUpdate := r.client.User.UpdateOneID(id)
	applyUserUpdateOptions(userUpdate, opt)
	userEntity, err := userUpdate.Save(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModelUser(userEntity), nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *UserRepositoryImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.client.User.UpdateOneID(id).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *UserRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	err := r.client.User.UpdateOneID(id).
		SetNillableDeletedAt(nil).
		Exec(ctx)
	return infraerr.MapNotFound(err)
}

func applyUserCreateOptions(create *ent.UserCreate, opt UserMutationOptions) {
	create.SetNillableName(opt.Name)
	create.SetNillableAvatarURL(opt.AvatarURL)
}

func applyUserUpdateOptions(update *ent.UserUpdateOne, opt UserMutationOptions) {
	update.SetNillableName(opt.Name)
	update.SetNillableAvatarURL(opt.AvatarURL)
}

func toModelUser(userEntity *ent.User) *repoUser.User {
	if userEntity == nil {
		return nil
	}

	return &repoUser.User{
		ID:        userEntity.ID,
		Email:     userEntity.Email,
		Name:      userEntity.Name,
		AvatarURL: userEntity.AvatarURL,
	}
}
