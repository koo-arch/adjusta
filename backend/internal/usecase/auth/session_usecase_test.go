package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	"github.com/koo-arch/adjusta-backend/internal/repositorymodel"
)

type fakeSessionAuthenticator struct {
	processUserSignInFn func(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repositorymodel.User, error)
	createSessionFn     func(ctx context.Context, userID uuid.UUID) (*repositorymodel.Session, error)
	deleteSessionFn     func(ctx context.Context, sessionToken string) error
}

func (f *fakeSessionAuthenticator) ProcessUserSignIn(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repositorymodel.User, error) {
	return f.processUserSignInFn(ctx, userInfo, oauthToken)
}

func (f *fakeSessionAuthenticator) CreateSession(ctx context.Context, userID uuid.UUID) (*repositorymodel.Session, error) {
	return f.createSessionFn(ctx, userID)
}

func (f *fakeSessionAuthenticator) DeleteSession(ctx context.Context, sessionToken string) error {
	return f.deleteSessionFn(ctx, sessionToken)
}

func TestSessionUsecaseCompleteGoogleSignInSuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	profile := &appmodel.GoogleUserProfile{
		GoogleID: "google-id",
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

	var exchangedCode string
	var fetchedToken *appmodel.GoogleAuthToken
	var signedInProfile *appmodel.GoogleUserProfile
	var createdSessionUserID uuid.UUID

	usecase := NewSessionUsecase(
		&fakeSessionAuthenticator{
			processUserSignInFn: func(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repositorymodel.User, error) {
				signedInProfile = userInfo
				if oauthToken != token {
					t.Fatalf("unexpected oauth token: %#v", oauthToken)
				}
				return &repositorymodel.User{
					ID:    userID,
					Email: userInfo.Email,
				}, nil
			},
			createSessionFn: func(ctx context.Context, gotUserID uuid.UUID) (*repositorymodel.Session, error) {
				createdSessionUserID = gotUserID
				return &repositorymodel.Session{
					ID:           uuid.New(),
					UserID:       gotUserID,
					SessionToken: "session-token",
					ExpiresAt:    time.Now().Add(time.Hour),
				}, nil
			},
			deleteSessionFn: func(ctx context.Context, sessionToken string) error {
				t.Fatalf("delete session should not be called")
				return nil
			},
		},
		OAuthExchangerFunc(func(ctx context.Context, code string) (*appmodel.GoogleAuthToken, error) {
			exchangedCode = code
			return token, nil
		}),
		UserInfoFetcherFunc(func(ctx context.Context, gotToken *appmodel.GoogleAuthToken) (*appmodel.GoogleUserProfile, error) {
			fetchedToken = gotToken
			return profile, nil
		}),
	)

	result, err := usecase.CompleteGoogleSignIn(ctx, "oauth-code")
	if err != nil {
		t.Fatalf("CompleteGoogleSignIn returned error: %v", err)
	}
	if exchangedCode != "oauth-code" {
		t.Fatalf("unexpected exchanged code: %q", exchangedCode)
	}
	if fetchedToken != token {
		t.Fatalf("unexpected fetched token: %#v", fetchedToken)
	}
	if signedInProfile != profile {
		t.Fatalf("unexpected signed in profile: %#v", signedInProfile)
	}
	if createdSessionUserID != userID {
		t.Fatalf("unexpected session user id: %s", createdSessionUserID)
	}
	if result.SessionToken != "session-token" || result.UserEmail != profile.Email {
		t.Fatalf("unexpected result: %#v", result)
	}
}
