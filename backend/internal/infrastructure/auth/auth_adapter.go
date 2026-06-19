package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
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

func (r *authReader) FindUserByEmail(ctx context.Context, email string) (*repoUser.User, error) {
	return r.userRepo.FindByEmail(ctx, email)
}

func (r *authReader) FindAccountByUserID(ctx context.Context, userID uuid.UUID) (*repoAccount.Account, error) {
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
			sessionRepo: repos.Session,
		})
	})
}

type authStore struct {
	userRepo    repoUser.UserRepository
	accountRepo repoAccount.AccountRepository
	sessionRepo repoSession.SessionRepository
}

func (s *authStore) CreateUser(ctx context.Context, email string, opt usecaseAuth.UserMutation) (*repoUser.User, error) {
	return s.userRepo.Create(ctx, email, repoUser.UserMutationOptions{
		Name:      opt.Name,
		AvatarURL: opt.AvatarURL,
	})
}

func (s *authStore) UpdateUser(ctx context.Context, userID uuid.UUID, opt usecaseAuth.UserMutation) (*repoUser.User, error) {
	return s.userRepo.Update(ctx, userID, repoUser.UserMutationOptions{
		Name:      opt.Name,
		AvatarURL: opt.AvatarURL,
	})
}

func (s *authStore) CreateAccount(ctx context.Context, userID uuid.UUID, opt usecaseAuth.AccountMutation) (*repoAccount.Account, error) {
	return s.accountRepo.Create(ctx, userID, repoAccount.AccountMutationOptions{
		GoogleUserID: opt.GoogleUserID,
		AccessToken:  opt.AccessToken,
		RefreshToken: opt.RefreshToken,
		ExpiresAt:    opt.ExpiresAt,
		Scope:        opt.Scope,
	})
}

func (s *authStore) UpdateAccount(ctx context.Context, accountID uuid.UUID, opt usecaseAuth.AccountMutation) (*repoAccount.Account, error) {
	return s.accountRepo.Update(ctx, accountID, repoAccount.AccountMutationOptions{
		GoogleUserID: opt.GoogleUserID,
		AccessToken:  opt.AccessToken,
		RefreshToken: opt.RefreshToken,
		ExpiresAt:    opt.ExpiresAt,
		Scope:        opt.Scope,
	})
}

func (s *authStore) CreateSession(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error) {
	return s.sessionRepo.Create(ctx, userID, sessionToken, expiresAt)
}

type authSessionStore struct {
	sessionRepo repoSession.SessionRepository
}

func NewAuthSessionStore(sessionRepo repoSession.SessionRepository) usecaseAuth.SessionStore {
	return &authSessionStore{sessionRepo: sessionRepo}
}

func (s *authSessionStore) CreateSession(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error) {
	return s.sessionRepo.Create(ctx, userID, sessionToken, expiresAt)
}

func (s *authSessionStore) FindSessionByToken(ctx context.Context, sessionToken string, withUser bool) (*repoSession.Session, error) {
	return s.sessionRepo.FindByToken(ctx, sessionToken, repoSession.SessionQueryOptions{
		WithUser: withUser,
	})
}

func (s *authSessionStore) UpdateSessionExpiry(ctx context.Context, sessionID uuid.UUID, expiresAt time.Time) (*repoSession.Session, error) {
	return s.sessionRepo.UpdateExpiry(ctx, sessionID, expiresAt)
}

func (s *authSessionStore) DeleteSessionByToken(ctx context.Context, sessionToken string) error {
	return s.sessionRepo.DeleteByToken(ctx, sessionToken)
}
