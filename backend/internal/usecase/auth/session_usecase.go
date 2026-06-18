package auth

import (
	"context"
	"log"

	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

type GoogleSignInResult struct {
	SessionToken string
	UserEmail    string
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

type SessionUsecase struct {
	authenticator SessionAuthenticator
	oauth         OAuthGateway
	userInfo      UserInfoFetcher
}

func NewSessionUsecase(authenticator SessionAuthenticator, oauth OAuthGateway, userInfo UserInfoFetcher) *SessionUsecase {
	return &SessionUsecase{
		authenticator: authenticator,
		oauth:         oauth,
		userInfo:      userInfo,
	}
}

func (uc *SessionUsecase) GoogleLoginURL(state string) string {
	return uc.oauth.AuthCodeURL(state)
}

func (uc *SessionUsecase) CompleteGoogleSignIn(ctx context.Context, code string) (*GoogleSignInResult, error) {
	oauthToken, err := uc.oauth.Exchange(ctx, code)
	if err != nil {
		log.Printf("failed to exchange oauth token: %v", err)
		return nil, internalErrors.NewInternalError("OAuthトークンの取得に失敗しました")
	}

	userInfo, err := uc.userInfo.FetchGoogleUserInfo(ctx, oauthToken)
	if err != nil {
		log.Printf("failed to fetch user info: %v", err)
		return nil, internalErrors.NormalizeAPIError(err, "ユーザー情報の取得に失敗しました")
	}

	entSession, u, err := uc.authenticator.SignInWithGoogle(ctx, userInfo, oauthToken)
	if err != nil {
		log.Printf("failed to create or update user: %v", err)
		return nil, internalErrors.NormalizeAPIError(err, "ユーザーの作成または更新に失敗しました")
	}

	return &GoogleSignInResult{
		SessionToken: entSession.SessionToken,
		UserEmail:    u.Email,
	}, nil
}

func (uc *SessionUsecase) Logout(ctx context.Context, sessionToken string) error {
	if err := uc.authenticator.DeleteSession(ctx, sessionToken); err != nil {
		log.Printf("failed to delete session: %v", err)
		return internalErrors.NormalizeAPIError(err, "セッションの削除に失敗しました")
	}

	return nil
}
