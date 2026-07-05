package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func TestAuthenticatorAuthenticateSessionDeletesExpiredSession(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expiredAt := time.Now().Add(-time.Minute)
	var deletedToken string
	sessions := &fakeSessionStore{
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
	}

	service := NewAuthenticator(
		fakeAuthRepositories(&fakeSignInReader{}, sessions),
		nil,
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
