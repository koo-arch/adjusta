package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type fakeSignInReader struct {
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
	store  SignInStore
	called bool
}

func (f *fakeSignInTransaction) Do(ctx context.Context, fn func(store SignInStore) error) error {
	f.called = true
	return fn(f.store)
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

func TestAuthServiceProcessUserSignInCreatesUserWhenUserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	scope := "openid email profile"
	expiry := time.Now().Add(time.Hour).UTC()
	profile := &appmodel.GoogleUserProfile{
		GoogleID: "google-user-id",
		Email:    "user@example.com",
		Name:     "Adjusta User",
		Picture:  "https://example.com/avatar.png",
	}
	token := &appmodel.GoogleAuthToken{
		AccessToken:  "access-token",
		TokenType:    "Bearer",
		RefreshToken: "refresh-token",
		Expiry:       expiry,
		Scope:        &scope,
	}

	var createUserCalled bool
	var createAccountCalled bool

	tx := &fakeSignInTransaction{
		store: &fakeSignInStore{
			createUserFn: func(ctx context.Context, email string, opt UserMutation) (*repoUser.User, error) {
				createUserCalled = true
				if email != profile.Email {
					t.Fatalf("unexpected email: %s", email)
				}
				if opt.Name == nil || *opt.Name != profile.Name {
					t.Fatalf("unexpected user name mutation: %#v", opt.Name)
				}
				if opt.AvatarURL == nil || *opt.AvatarURL != profile.Picture {
					t.Fatalf("unexpected avatar mutation: %#v", opt.AvatarURL)
				}
				return &repoUser.User{
					ID:    userID,
					Email: email,
					Name:  opt.Name,
				}, nil
			},
			updateUserFn: func(ctx context.Context, userID uuid.UUID, opt UserMutation) (*repoUser.User, error) {
				t.Fatalf("update user should not be called")
				return nil, nil
			},
			createAccountFn: func(ctx context.Context, gotUserID uuid.UUID, opt AccountMutation) (*repoAccount.Account, error) {
				createAccountCalled = true
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				if opt.GoogleUserID == nil || *opt.GoogleUserID != profile.GoogleID {
					t.Fatalf("unexpected google user id mutation: %#v", opt.GoogleUserID)
				}
				if opt.AccessToken == nil || *opt.AccessToken != token.AccessToken {
					t.Fatalf("unexpected access token mutation: %#v", opt.AccessToken)
				}
				if opt.RefreshToken == nil || *opt.RefreshToken != token.RefreshToken {
					t.Fatalf("unexpected refresh token mutation: %#v", opt.RefreshToken)
				}
				if opt.ExpiresAt == nil || !opt.ExpiresAt.Equal(expiry) {
					t.Fatalf("unexpected expiry mutation: %#v", opt.ExpiresAt)
				}
				if opt.Scope == nil || *opt.Scope != scope {
					t.Fatalf("unexpected scope mutation: %#v", opt.Scope)
				}
				return &repoAccount.Account{ID: uuid.New(), UserID: gotUserID}, nil
			},
			updateAccountFn: func(ctx context.Context, accountID uuid.UUID, opt AccountMutation) (*repoAccount.Account, error) {
				t.Fatalf("update account should not be called")
				return nil, nil
			},
			createSessionFn: func(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error) {
				t.Fatalf("create session should not be called")
				return nil, nil
			},
		},
	}

	service := NewAuthService(
		&fakeSignInReader{
			findUserByEmailFn: func(ctx context.Context, email string) (*repoUser.User, error) {
				return nil, repoerr.ErrNotFound
			},
			findAccountByUserIDFn: func(ctx context.Context, userID uuid.UUID) (*repoAccount.Account, error) {
				t.Fatalf("find account should not be called when user is missing")
				return nil, nil
			},
		},
		tx,
		&fakeSessionStore{},
		time.Hour,
	)

	gotUser, err := service.ProcessUserSignIn(ctx, profile, token)
	if err != nil {
		t.Fatalf("ProcessUserSignIn returned error: %v", err)
	}
	if gotUser == nil || gotUser.ID != userID {
		t.Fatalf("unexpected user: %#v", gotUser)
	}
	if !tx.called || !createUserCalled || !createAccountCalled {
		t.Fatalf("expected transaction and create operations to be called")
	}
}

func TestAuthServiceSignInWithGoogleCreatesSessionInTransaction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	profile := &appmodel.GoogleUserProfile{
		GoogleID: "google-user-id",
		Email:    "user@example.com",
		Name:     "Adjusta User",
		Picture:  "https://example.com/avatar.png",
	}
	token := &appmodel.GoogleAuthToken{
		AccessToken:  "access-token",
		TokenType:    "Bearer",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(time.Hour).UTC(),
	}

	var createdSessionUserID uuid.UUID

	tx := &fakeSignInTransaction{
		store: &fakeSignInStore{
			createUserFn: func(ctx context.Context, email string, opt UserMutation) (*repoUser.User, error) {
				return &repoUser.User{
					ID:    userID,
					Email: email,
				}, nil
			},
			updateUserFn: func(ctx context.Context, userID uuid.UUID, opt UserMutation) (*repoUser.User, error) {
				t.Fatalf("update user should not be called")
				return nil, nil
			},
			createAccountFn: func(ctx context.Context, gotUserID uuid.UUID, opt AccountMutation) (*repoAccount.Account, error) {
				return &repoAccount.Account{ID: uuid.New(), UserID: gotUserID}, nil
			},
			updateAccountFn: func(ctx context.Context, accountID uuid.UUID, opt AccountMutation) (*repoAccount.Account, error) {
				t.Fatalf("update account should not be called")
				return nil, nil
			},
			createSessionFn: func(ctx context.Context, gotUserID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error) {
				createdSessionUserID = gotUserID
				if sessionToken == "" {
					t.Fatalf("expected session token to be set")
				}
				return &repoSession.Session{
					ID:           uuid.New(),
					UserID:       gotUserID,
					SessionToken: sessionToken,
					ExpiresAt:    expiresAt,
				}, nil
			},
		},
	}

	service := NewAuthService(
		&fakeSignInReader{
			findUserByEmailFn: func(ctx context.Context, email string) (*repoUser.User, error) {
				return nil, repoerr.ErrNotFound
			},
			findAccountByUserIDFn: func(ctx context.Context, userID uuid.UUID) (*repoAccount.Account, error) {
				t.Fatalf("find account should not be called when user is missing")
				return nil, nil
			},
		},
		tx,
		&fakeSessionStore{},
		time.Hour,
	)

	session, gotUser, err := service.SignInWithGoogle(ctx, profile, token)
	if err != nil {
		t.Fatalf("SignInWithGoogle returned error: %v", err)
	}
	if !tx.called {
		t.Fatalf("expected transaction to be called")
	}
	if gotUser == nil || gotUser.ID != userID {
		t.Fatalf("unexpected user: %#v", gotUser)
	}
	if session == nil || session.UserID != userID {
		t.Fatalf("unexpected session: %#v", session)
	}
	if createdSessionUserID != userID {
		t.Fatalf("unexpected session user id: %s", createdSessionUserID)
	}
}

func TestAuthServiceAuthenticateSessionDeletesExpiredSession(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expiredAt := time.Now().Add(-time.Minute)
	var deletedToken string

	service := NewAuthService(
		nil,
		nil,
		&fakeSessionStore{
			createSessionFn: func(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error) {
				t.Fatalf("create session should not be called")
				return nil, nil
			},
			findSessionByTokenFn: func(ctx context.Context, sessionToken string, withUser bool) (*repoSession.Session, error) {
				return &repoSession.Session{
					ID:           uuid.New(),
					UserID:       uuid.New(),
					SessionToken: sessionToken,
					ExpiresAt:    expiredAt,
					User: &repoUser.User{
						ID:    uuid.New(),
						Email: "user@example.com",
					},
				}, nil
			},
			updateSessionExpiryFn: func(ctx context.Context, sessionID uuid.UUID, expiresAt time.Time) (*repoSession.Session, error) {
				t.Fatalf("update session expiry should not be called for expired session")
				return nil, nil
			},
			deleteSessionByTokenFn: func(ctx context.Context, sessionToken string) error {
				deletedToken = sessionToken
				return nil
			},
		},
		time.Hour,
	)

	gotUser, err := service.AuthenticateSession(ctx, "expired-session")
	if gotUser != nil {
		t.Fatalf("expected nil user, got %#v", gotUser)
	}
	apiErr, ok := err.(*internalErrors.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Kind != internalErrors.KindUnauthorized {
		t.Fatalf("unexpected error kind: %s", apiErr.Kind)
	}
	if deletedToken != "expired-session" {
		t.Fatalf("expected expired session to be deleted, got token %q", deletedToken)
	}
}
