package auth

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repositorymodel"
)

type GoogleSignInResult struct {
	SessionToken string
	UserEmail    string
}

type SessionAuthenticator interface {
	ProcessUserSignIn(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repositorymodel.User, error)
	CreateSession(ctx context.Context, userID uuid.UUID) (*repositorymodel.Session, error)
	DeleteSession(ctx context.Context, sessionToken string) error
}

type OAuthExchanger interface {
	Exchange(ctx context.Context, code string) (*appmodel.GoogleAuthToken, error)
}

type UserInfoFetcher interface {
	FetchGoogleUserInfo(ctx context.Context, token *appmodel.GoogleAuthToken) (*appmodel.GoogleUserProfile, error)
}

type OAuthExchangerFunc func(ctx context.Context, code string) (*appmodel.GoogleAuthToken, error)

func (f OAuthExchangerFunc) Exchange(ctx context.Context, code string) (*appmodel.GoogleAuthToken, error) {
	return f(ctx, code)
}

type UserInfoFetcherFunc func(ctx context.Context, token *appmodel.GoogleAuthToken) (*appmodel.GoogleUserProfile, error)

func (f UserInfoFetcherFunc) FetchGoogleUserInfo(ctx context.Context, token *appmodel.GoogleAuthToken) (*appmodel.GoogleUserProfile, error) {
	return f(ctx, token)
}

type SessionUsecase struct {
	authenticator SessionAuthenticator
	oauthConfig   OAuthExchanger
	userInfo      UserInfoFetcher
}

func NewSessionUsecase(authenticator SessionAuthenticator, oauthConfig OAuthExchanger, userInfo UserInfoFetcher) *SessionUsecase {
	return &SessionUsecase{
		authenticator: authenticator,
		oauthConfig:   oauthConfig,
		userInfo:      userInfo,
	}
}

func (uc *SessionUsecase) CompleteGoogleSignIn(ctx context.Context, code string) (*GoogleSignInResult, error) {
	oauthToken, err := uc.oauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Printf("failed to exchange oauth token: %v", err)
		return nil, internalErrors.NewInternalError("OAuthトークンの取得に失敗しました")
	}

	userInfo, err := uc.userInfo.FetchGoogleUserInfo(ctx, oauthToken)
	if err != nil {
		log.Printf("failed to fetch user info: %v", err)
		return nil, internalErrors.NormalizeAPIError(err, "ユーザー情報の取得に失敗しました")
	}

	u, err := uc.authenticator.ProcessUserSignIn(ctx, userInfo, oauthToken)
	if err != nil {
		log.Printf("failed to create or update user: %v", err)
		return nil, internalErrors.NormalizeAPIError(err, "ユーザーの作成または更新に失敗しました")
	}

	entSession, err := uc.authenticator.CreateSession(ctx, u.ID)
	if err != nil {
		log.Printf("failed to create app session: %v", err)
		return nil, internalErrors.NormalizeAPIError(err, "ログインセッションの作成に失敗しました")
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
