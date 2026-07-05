package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	"github.com/koo-arch/adjusta-backend/internal/google"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func TestAuthenticatorProcessUserSignInCreatesUserWhenUserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	scope := "openid email profile"
	expiry := time.Now().Add(time.Hour).UTC()
	profile := &google.UserProfile{
		GoogleID: "google-user-id",
		Email:    "user@example.com",
		Name:     "Adjusta User",
		Picture:  "https://example.com/avatar.png",
	}
	token := &google.AuthToken{
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

	reader := &fakeSignInReader{
		findUserByEmailFn: func(ctx context.Context, email string) (*repoUser.User, error) {
			return nil, repoerr.ErrNotFound
		},
		findAccountByUserIDFn: func(ctx context.Context, userID uuid.UUID) (*repoAccount.Account, error) {
			t.Fatalf("find account should not be called when user is missing")
			return nil, nil
		},
	}

	service := NewAuthenticator(
		fakeAuthRepositories(reader, &fakeSessionStore{}),
		tx,
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

func TestAuthenticatorSignInWithGoogleCreatesSessionInTransaction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	profile := &google.UserProfile{
		GoogleID: "google-user-id",
		Email:    "user@example.com",
		Name:     "Adjusta User",
		Picture:  "https://example.com/avatar.png",
	}
	token := &google.AuthToken{
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

	reader := &fakeSignInReader{
		findUserByEmailFn: func(ctx context.Context, email string) (*repoUser.User, error) {
			return nil, repoerr.ErrNotFound
		},
		findAccountByUserIDFn: func(ctx context.Context, userID uuid.UUID) (*repoAccount.Account, error) {
			t.Fatalf("find account should not be called when user is missing")
			return nil, nil
		},
	}

	service := NewAuthenticator(
		fakeAuthRepositories(reader, &fakeSessionStore{}),
		tx,
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
