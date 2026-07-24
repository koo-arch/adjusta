package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/google"
)

func TestAuthenticatorReauthorizeGoogleUpdatesTokensWithoutCreatingSession(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()
	profile := &google.UserProfile{GoogleID: "google-id", Email: "user@example.com", Name: "Updated User"}
	token := &google.AuthToken{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		Expiry:       time.Now().Add(time.Hour).UTC(),
	}

	reader := &fakeSignInReader{
		readUserByIDFn: func(ctx context.Context, gotUserID uuid.UUID) (*repoUser.User, error) {
			return &repoUser.User{ID: gotUserID, Email: profile.Email}, nil
		},
		findAccountByUserIDFn: func(ctx context.Context, gotUserID uuid.UUID) (*repoAccount.Account, error) {
			return &repoAccount.Account{ID: accountID, UserID: gotUserID, GoogleUserID: profile.GoogleID}, nil
		},
	}
	var accountUpdated bool
	tx := &fakeSignInTransaction{store: &fakeSignInStore{
		updateUserFn: func(ctx context.Context, gotUserID uuid.UUID, opt UserMutation) (*repoUser.User, error) {
			return &repoUser.User{ID: gotUserID, Email: profile.Email, Name: opt.Name}, nil
		},
		updateAccountFn: func(ctx context.Context, gotAccountID uuid.UUID, opt AccountMutation) (*repoAccount.Account, error) {
			accountUpdated = true
			if gotAccountID != accountID || opt.RefreshToken == nil || *opt.RefreshToken != token.RefreshToken {
				t.Fatalf("unexpected account update: %#v", opt)
			}
			return &repoAccount.Account{ID: gotAccountID, UserID: userID}, nil
		},
		createSessionFn: func(ctx context.Context, userID uuid.UUID, sessionToken string, expiresAt time.Time) (*repoSession.Session, error) {
			t.Fatal("session should not be created during reauthorization")
			return nil, nil
		},
	}}

	authenticator := NewAuthenticator(fakeAuthRepositories(reader, &fakeSessionStore{}), tx, time.Hour)
	if err := authenticator.ReauthorizeGoogle(ctx, userID, profile, token); err != nil {
		t.Fatalf("ReauthorizeGoogle returned error: %v", err)
	}
	if !accountUpdated {
		t.Fatal("expected account tokens to be updated")
	}
}

func TestAuthenticatorReauthorizeGoogleRejectsDifferentGoogleAccount(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	reader := &fakeSignInReader{
		readUserByIDFn: func(ctx context.Context, gotUserID uuid.UUID) (*repoUser.User, error) {
			return &repoUser.User{ID: gotUserID, Email: "user@example.com"}, nil
		},
		findAccountByUserIDFn: func(ctx context.Context, gotUserID uuid.UUID) (*repoAccount.Account, error) {
			return &repoAccount.Account{ID: uuid.New(), UserID: gotUserID, GoogleUserID: "expected-google-id"}, nil
		},
	}
	tx := &fakeSignInTransaction{store: &fakeSignInStore{}}
	authenticator := NewAuthenticator(fakeAuthRepositories(reader, &fakeSessionStore{}), tx, time.Hour)

	err := authenticator.ReauthorizeGoogle(
		context.Background(),
		userID,
		&google.UserProfile{GoogleID: "different-google-id", Email: "other@example.com"},
		&google.AuthToken{RefreshToken: "refresh-token"},
	)
	if !internalErrors.IsKind(err, internalErrors.KindForbidden) {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.called {
		t.Fatal("transaction should not run for a different Google account")
	}
}
