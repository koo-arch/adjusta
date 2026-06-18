package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
)

type fakeSessionAuthenticator struct {
	signInWithGoogleFn func(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repoSession.Session, *repoUser.User, error)
	deleteSessionFn    func(ctx context.Context, sessionToken string) error
}

func (f *fakeSessionAuthenticator) SignInWithGoogle(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repoSession.Session, *repoUser.User, error) {
	return f.signInWithGoogleFn(ctx, userInfo, oauthToken)
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

	usecase := NewSessionUsecase(
		&fakeSessionAuthenticator{
			signInWithGoogleFn: func(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repoSession.Session, *repoUser.User, error) {
				signedInProfile = userInfo
				if oauthToken != token {
					t.Fatalf("unexpected oauth token: %#v", oauthToken)
				}
				session := &repoSession.Session{
					ID:           uuid.New(),
					UserID:       userID,
					SessionToken: "session-token",
					ExpiresAt:    time.Now().Add(time.Hour),
				}
				return session, &repoUser.User{
					ID:    userID,
					Email: userInfo.Email,
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
	if result.SessionToken != "session-token" || result.UserEmail != profile.Email {
		t.Fatalf("unexpected result: %#v", result)
	}
}
