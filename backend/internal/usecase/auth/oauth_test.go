package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	"github.com/koo-arch/adjusta-backend/internal/google"
)

type fakeSessionAuthenticator struct {
	signInWithGoogleFn func(ctx context.Context, userInfo *google.UserProfile, oauthToken *google.AuthToken) (*repoSession.Session, *repoUser.User, error)
	deleteSessionFn    func(ctx context.Context, sessionToken string) error
}

func (f *fakeSessionAuthenticator) SignInWithGoogle(ctx context.Context, userInfo *google.UserProfile, oauthToken *google.AuthToken) (*repoSession.Session, *repoUser.User, error) {
	return f.signInWithGoogleFn(ctx, userInfo, oauthToken)
}

func (f *fakeSessionAuthenticator) DeleteSession(ctx context.Context, sessionToken string) error {
	return f.deleteSessionFn(ctx, sessionToken)
}

func TestOAuthUsecaseCompleteGoogleSignInSuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	profile := &google.UserProfile{
		GoogleID: "google-id",
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

	var exchangedCode string
	var fetchedToken *google.AuthToken
	var signedInProfile *google.UserProfile

	usecase := NewOAuthUsecase(
		&fakeSessionAuthenticator{
			signInWithGoogleFn: func(ctx context.Context, userInfo *google.UserProfile, oauthToken *google.AuthToken) (*repoSession.Session, *repoUser.User, error) {
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
		OAuthGatewayFuncs{
			AuthCodeURLFn: func(state string) string {
				t.Fatalf("auth code url should not be called")
				return ""
			},
			ExchangeFn: func(ctx context.Context, code string) (*google.AuthToken, error) {
				exchangedCode = code
				return token, nil
			},
		},
		UserInfoFetcherFunc(func(ctx context.Context, gotToken *google.AuthToken) (*google.UserProfile, error) {
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

func TestOAuthUsecaseGoogleLoginURL(t *testing.T) {
	t.Parallel()

	usecase := NewOAuthUsecase(
		&fakeSessionAuthenticator{
			signInWithGoogleFn: func(ctx context.Context, userInfo *google.UserProfile, oauthToken *google.AuthToken) (*repoSession.Session, *repoUser.User, error) {
				t.Fatalf("sign in should not be called")
				return nil, nil, nil
			},
			deleteSessionFn: func(ctx context.Context, sessionToken string) error {
				t.Fatalf("delete session should not be called")
				return nil
			},
		},
		OAuthGatewayFuncs{
			AuthCodeURLFn: func(state string) string {
				return "https://example.com/oauth?state=" + state
			},
			ExchangeFn: func(ctx context.Context, code string) (*google.AuthToken, error) {
				t.Fatalf("exchange should not be called")
				return nil, nil
			},
		},
		UserInfoFetcherFunc(func(ctx context.Context, token *google.AuthToken) (*google.UserProfile, error) {
			t.Fatalf("fetch user info should not be called")
			return nil, nil
		}),
	)

	got := usecase.GoogleLoginURL("oauth-state")
	if got != "https://example.com/oauth?state=oauth-state" {
		t.Fatalf("unexpected login url: %q", got)
	}
}
