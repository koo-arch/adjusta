package auth

import (
	"context"

	"github.com/google/uuid"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	"github.com/koo-arch/adjusta-backend/internal/google"
)

type SignInTransaction interface {
	Do(ctx context.Context, fn func(repos AuthTxRepositories) error) error
}

type SessionAuthenticator interface {
	AuthenticateSession(ctx context.Context, sessionToken string) (*repoUser.User, error)
	SignInWithGoogle(ctx context.Context, userInfo *google.UserProfile, oauthToken *google.AuthToken) (*repoSession.Session, *repoUser.User, error)
	ReauthorizeGoogle(ctx context.Context, userID uuid.UUID, userInfo *google.UserProfile, oauthToken *google.AuthToken) error
	DeleteSession(ctx context.Context, sessionToken string) error
}

type OAuthGateway interface {
	AuthCodeURL(state string) string
	Exchange(ctx context.Context, code string) (*google.AuthToken, error)
}

type UserInfoFetcher interface {
	FetchGoogleUserInfo(ctx context.Context, token *google.AuthToken) (*google.UserProfile, error)
}

type OAuthGatewayFuncs struct {
	AuthCodeURLFn func(state string) string
	ExchangeFn    func(ctx context.Context, code string) (*google.AuthToken, error)
}

func (f OAuthGatewayFuncs) AuthCodeURL(state string) string {
	return f.AuthCodeURLFn(state)
}

func (f OAuthGatewayFuncs) Exchange(ctx context.Context, code string) (*google.AuthToken, error) {
	return f.ExchangeFn(ctx, code)
}

type UserInfoFetcherFunc func(ctx context.Context, token *google.AuthToken) (*google.UserProfile, error)

func (f UserInfoFetcherFunc) FetchGoogleUserInfo(ctx context.Context, token *google.AuthToken) (*google.UserProfile, error) {
	return f(ctx, token)
}
