package auth

import (
	"context"

	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
)

type SignInTransaction interface {
	Do(ctx context.Context, fn func(repos AuthTxRepositories) error) error
}

type SessionAuthenticator interface {
	SignInWithGoogle(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repoSession.Session, *repoUser.User, error)
	DeleteSession(ctx context.Context, sessionToken string) error
}

type OAuthGateway interface {
	AuthCodeURL(state string) string
	Exchange(ctx context.Context, code string) (*appmodel.GoogleAuthToken, error)
}

type UserInfoFetcher interface {
	FetchGoogleUserInfo(ctx context.Context, token *appmodel.GoogleAuthToken) (*appmodel.GoogleUserProfile, error)
}

type OAuthGatewayFuncs struct {
	AuthCodeURLFn func(state string) string
	ExchangeFn    func(ctx context.Context, code string) (*appmodel.GoogleAuthToken, error)
}

func (f OAuthGatewayFuncs) AuthCodeURL(state string) string {
	return f.AuthCodeURLFn(state)
}

func (f OAuthGatewayFuncs) Exchange(ctx context.Context, code string) (*appmodel.GoogleAuthToken, error) {
	return f.ExchangeFn(ctx, code)
}

type UserInfoFetcherFunc func(ctx context.Context, token *appmodel.GoogleAuthToken) (*appmodel.GoogleUserProfile, error)

func (f UserInfoFetcherFunc) FetchGoogleUserInfo(ctx context.Context, token *appmodel.GoogleAuthToken) (*appmodel.GoogleUserProfile, error) {
	return f(ctx, token)
}
