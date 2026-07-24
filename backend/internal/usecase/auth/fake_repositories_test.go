package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
)

type fakeSignInReader struct {
	readUserByIDFn        func(ctx context.Context, userID uuid.UUID) (*repoUser.User, error)
	findUserByEmailFn     func(ctx context.Context, email string) (*repoUser.User, error)
	findAccountByUserIDFn func(ctx context.Context, userID uuid.UUID) (*repoAccount.Account, error)
}

func (f *fakeSignInReader) FindUserByEmail(ctx context.Context, email string) (*repoUser.User, error) {
	return f.findUserByEmailFn(ctx, email)
}

func (f *fakeSignInReader) FindAccountByUserID(ctx context.Context, userID uuid.UUID) (*repoAccount.Account, error) {
	return f.findAccountByUserIDFn(ctx, userID)
}

type fakeSignInStore struct {
	createUserFn    func(ctx context.Context, email string, opt UserMutation) (*repoUser.User, error)
	updateUserFn    func(ctx context.Context, userID uuid.UUID, opt UserMutation) (*repoUser.User, error)
	createAccountFn func(ctx context.Context, userID uuid.UUID, opt AccountMutation) (*repoAccount.Account, error)
	updateAccountFn func(ctx context.Context, accountID uuid.UUID, opt AccountMutation) (*repoAccount.Account, error)
	createSessionFn func(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error)
}

func (f *fakeSignInStore) CreateUser(ctx context.Context, email string, opt UserMutation) (*repoUser.User, error) {
	return f.createUserFn(ctx, email, opt)
}

func (f *fakeSignInStore) UpdateUser(ctx context.Context, userID uuid.UUID, opt UserMutation) (*repoUser.User, error) {
	return f.updateUserFn(ctx, userID, opt)
}

func (f *fakeSignInStore) CreateAccount(ctx context.Context, userID uuid.UUID, opt AccountMutation) (*repoAccount.Account, error) {
	return f.createAccountFn(ctx, userID, opt)
}

func (f *fakeSignInStore) UpdateAccount(ctx context.Context, accountID uuid.UUID, opt AccountMutation) (*repoAccount.Account, error) {
	return f.updateAccountFn(ctx, accountID, opt)
}

func (f *fakeSignInStore) CreateSession(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error) {
	return f.createSessionFn(ctx, userID, sessionToken, expiresAt)
}

type fakeSignInTransaction struct {
	store  *fakeSignInStore
	called bool
}

func (f *fakeSignInTransaction) Do(ctx context.Context, fn func(repos AuthTxRepositories) error) error {
	f.called = true
	return fn(fakeAuthTxRepositories(f.store))
}

type fakeSessionStore struct {
	createSessionFn        func(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error)
	findSessionByTokenFn   func(ctx context.Context, sessionToken string, withUser bool) (*repoSession.Session, error)
	updateSessionExpiryFn  func(ctx context.Context, sessionID uuid.UUID, expiresAt time.Time) (*repoSession.Session, error)
	deleteSessionByTokenFn func(ctx context.Context, sessionToken string) error
}

func (f *fakeSessionStore) CreateSession(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error) {
	return f.createSessionFn(ctx, userID, sessionToken, expiresAt)
}

func (f *fakeSessionStore) FindSessionByToken(ctx context.Context, sessionToken string, withUser bool) (*repoSession.Session, error) {
	return f.findSessionByTokenFn(ctx, sessionToken, withUser)
}

func (f *fakeSessionStore) UpdateSessionExpiry(ctx context.Context, sessionID uuid.UUID, expiresAt time.Time) (*repoSession.Session, error) {
	return f.updateSessionExpiryFn(ctx, sessionID, expiresAt)
}

func (f *fakeSessionStore) DeleteSessionByToken(ctx context.Context, sessionToken string) error {
	return f.deleteSessionByTokenFn(ctx, sessionToken)
}

func fakeAuthRepositories(reader *fakeSignInReader, sessions *fakeSessionStore) AuthRepositories {
	return AuthRepositories{
		User:    &fakeAuthUserRepository{reader: reader},
		Account: &fakeAuthAccountRepository{reader: reader},
		Session: &fakeAuthSessionRepository{sessions: sessions},
	}
}

func fakeAuthTxRepositories(store *fakeSignInStore) AuthTxRepositories {
	return AuthTxRepositories{
		User:    &fakeAuthUserRepository{store: store},
		Account: &fakeAuthAccountRepository{store: store},
		Session: &fakeAuthSessionRepository{store: store},
	}
}

type fakeAuthUserRepository struct {
	reader *fakeSignInReader
	store  *fakeSignInStore
}

func (r *fakeAuthUserRepository) Read(ctx context.Context, id uuid.UUID) (*repoUser.User, error) {
	return r.reader.readUserByIDFn(ctx, id)
}

func (r *fakeAuthUserRepository) FindByEmail(ctx context.Context, email string) (*repoUser.User, error) {
	return r.reader.findUserByEmailFn(ctx, email)
}

func (r *fakeAuthUserRepository) Create(ctx context.Context, email string, opt repoUser.UserMutationOptions) (*repoUser.User, error) {
	return r.store.createUserFn(ctx, email, opt)
}

func (r *fakeAuthUserRepository) Update(ctx context.Context, id uuid.UUID, opt repoUser.UserMutationOptions) (*repoUser.User, error) {
	return r.store.updateUserFn(ctx, id, opt)
}

func (r *fakeAuthUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *fakeAuthUserRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *fakeAuthUserRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return nil
}

type fakeAuthAccountRepository struct {
	reader *fakeSignInReader
	store  *fakeSignInStore
}

func (r *fakeAuthAccountRepository) Read(ctx context.Context, id uuid.UUID) (*repoAccount.Account, error) {
	return nil, nil
}

func (r *fakeAuthAccountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*repoAccount.Account, error) {
	return r.reader.findAccountByUserIDFn(ctx, userID)
}

func (r *fakeAuthAccountRepository) Create(ctx context.Context, userID uuid.UUID, opt repoAccount.AccountMutationOptions) (*repoAccount.Account, error) {
	return r.store.createAccountFn(ctx, userID, opt)
}

func (r *fakeAuthAccountRepository) Update(ctx context.Context, id uuid.UUID, opt repoAccount.AccountMutationOptions) (*repoAccount.Account, error) {
	return r.store.updateAccountFn(ctx, id, opt)
}

type fakeAuthSessionRepository struct {
	store    *fakeSignInStore
	sessions *fakeSessionStore
}

func (r *fakeAuthSessionRepository) Read(ctx context.Context, id uuid.UUID, opt repoSession.SessionQueryOptions) (*repoSession.Session, error) {
	return nil, nil
}

func (r *fakeAuthSessionRepository) FindByToken(ctx context.Context, sessionToken string, opt repoSession.SessionQueryOptions) (*repoSession.Session, error) {
	return r.sessions.findSessionByTokenFn(ctx, sessionToken, opt.WithUser)
}

func (r *fakeAuthSessionRepository) Create(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error) {
	return r.store.createSessionFn(ctx, userID, sessionToken, expiresAt)
}

func (r *fakeAuthSessionRepository) UpdateExpiry(ctx context.Context, id uuid.UUID, expiresAt time.Time) (*repoSession.Session, error) {
	return r.sessions.updateSessionExpiryFn(ctx, id, expiresAt)
}

func (r *fakeAuthSessionRepository) DeleteByToken(ctx context.Context, sessionToken string) error {
	return r.sessions.deleteSessionByTokenFn(ctx, sessionToken)
}
