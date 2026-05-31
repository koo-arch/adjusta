package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/repo/account"
	repoSession "github.com/koo-arch/adjusta-backend/internal/repo/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/repo/user"
	"github.com/koo-arch/adjusta-backend/internal/repositorymodel"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
)

type authReader struct {
	userRepo    repoUser.UserRepository
	accountRepo repoAccount.AccountRepository
}

func NewAuthReader(userRepo repoUser.UserRepository, accountRepo repoAccount.AccountRepository) usecaseAuth.SignInReader {
	return &authReader{
		userRepo:    userRepo,
		accountRepo: accountRepo,
	}
}

func (r *authReader) FindUserByEmail(ctx context.Context, email string) (*repositorymodel.User, error) {
	return r.userRepo.FindByEmail(ctx, email, repoUser.UserQueryOptions{})
}

func (r *authReader) FindAccountByUserID(ctx context.Context, userID uuid.UUID) (*repositorymodel.Account, error) {
	return r.accountRepo.FindByUserID(ctx, userID)
}

type authTransaction struct {
	uow infraRepository.UnitOfWork
}

func NewAuthTransaction(uow infraRepository.UnitOfWork) usecaseAuth.SignInTransaction {
	return &authTransaction{uow: uow}
}

func (t *authTransaction) Do(ctx context.Context, fn func(store usecaseAuth.SignInStore) error) error {
	return t.uow.Do(ctx, func(repos infraRepository.Repositories) error {
		return fn(&authStore{
			userRepo:    repos.User,
			accountRepo: repos.Account,
		})
	})
}

type authStore struct {
	userRepo    repoUser.UserRepository
	accountRepo repoAccount.AccountRepository
}

func (s *authStore) CreateUser(ctx context.Context, email string, opt usecaseAuth.UserMutation) (*repositorymodel.User, error) {
	return s.userRepo.Create(ctx, email, repoUser.UserMutationOptions{
		Name:      opt.Name,
		AvatarURL: opt.AvatarURL,
	})
}

func (s *authStore) UpdateUser(ctx context.Context, userID uuid.UUID, opt usecaseAuth.UserMutation) (*repositorymodel.User, error) {
	return s.userRepo.Update(ctx, userID, repoUser.UserMutationOptions{
		Name:      opt.Name,
		AvatarURL: opt.AvatarURL,
	})
}

func (s *authStore) CreateAccount(ctx context.Context, userID uuid.UUID, opt usecaseAuth.AccountMutation) (*repositorymodel.Account, error) {
	return s.accountRepo.Create(ctx, userID, repoAccount.AccountMutationOptions{
		GoogleUserID: opt.GoogleUserID,
		AccessToken:  opt.AccessToken,
		RefreshToken: opt.RefreshToken,
		ExpiresAt:    opt.ExpiresAt,
		Scope:        opt.Scope,
	})
}

func (s *authStore) UpdateAccount(ctx context.Context, accountID uuid.UUID, opt usecaseAuth.AccountMutation) (*repositorymodel.Account, error) {
	return s.accountRepo.Update(ctx, accountID, repoAccount.AccountMutationOptions{
		GoogleUserID: opt.GoogleUserID,
		AccessToken:  opt.AccessToken,
		RefreshToken: opt.RefreshToken,
		ExpiresAt:    opt.ExpiresAt,
		Scope:        opt.Scope,
	})
}

type authSessionStore struct {
	sessionRepo repoSession.SessionRepository
}

func NewAuthSessionStore(sessionRepo repoSession.SessionRepository) usecaseAuth.SessionStore {
	return &authSessionStore{sessionRepo: sessionRepo}
}

func (s *authSessionStore) CreateSession(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repositorymodel.Session, error) {
	return s.sessionRepo.Create(ctx, userID, sessionToken, expiresAt)
}

func (s *authSessionStore) FindSessionByToken(ctx context.Context, sessionToken string, withUser bool) (*repositorymodel.Session, error) {
	return s.sessionRepo.FindByToken(ctx, sessionToken, repoSession.SessionQueryOptions{
		WithUser: withUser,
	})
}

func (s *authSessionStore) UpdateSessionExpiry(ctx context.Context, sessionID uuid.UUID, expiresAt time.Time) (*repositorymodel.Session, error) {
	return s.sessionRepo.UpdateExpiry(ctx, sessionID, expiresAt)
}

func (s *authSessionStore) DeleteSessionByToken(ctx context.Context, sessionToken string) error {
	return s.sessionRepo.DeleteByToken(ctx, sessionToken)
}
