package auth

import (
	"context"
	"log"

	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

type OAuthUsecase struct {
	authenticator SessionAuthenticator
	gateway       OAuthGateway
	userInfo      UserInfoFetcher
}

func NewOAuthUsecase(authenticator SessionAuthenticator, gateway OAuthGateway, userInfo UserInfoFetcher) *OAuthUsecase {
	return &OAuthUsecase{
		authenticator: authenticator,
		gateway:       gateway,
		userInfo:      userInfo,
	}
}

func (uc *OAuthUsecase) GoogleLoginURL(state string) string {
	return uc.gateway.AuthCodeURL(state)
}

func (uc *OAuthUsecase) CompleteGoogleSignIn(ctx context.Context, code string) (*GoogleSignInResult, error) {
	oauthToken, err := uc.gateway.Exchange(ctx, code)
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

func (uc *OAuthUsecase) CompleteGoogleReauthorization(ctx context.Context, code, sessionToken string) error {
	user, err := uc.authenticator.AuthenticateSession(ctx, sessionToken)
	if err != nil {
		return err
	}

	oauthToken, err := uc.gateway.Exchange(ctx, code)
	if err != nil {
		log.Printf("failed to exchange oauth token during reauthorization: %v", err)
		return internalErrors.NewBadGatewayError("Google再認可のトークン取得に失敗しました")
	}

	userInfo, err := uc.userInfo.FetchGoogleUserInfo(ctx, oauthToken)
	if err != nil {
		log.Printf("failed to fetch user info during reauthorization: %v", err)
		return internalErrors.NormalizeAPIError(err, "Google再認可のユーザー情報取得に失敗しました")
	}

	if err := uc.authenticator.ReauthorizeGoogle(ctx, user.ID, userInfo, oauthToken); err != nil {
		log.Printf("failed to persist google reauthorization: %v", err)
		return internalErrors.NormalizeAPIError(err, "Google再認可情報の保存に失敗しました")
	}

	return nil
}

func (uc *OAuthUsecase) Logout(ctx context.Context, sessionToken string) error {
	if err := uc.authenticator.DeleteSession(ctx, sessionToken); err != nil {
		log.Printf("failed to delete session: %v", err)
		return internalErrors.NormalizeAPIError(err, "セッションの削除に失敗しました")
	}

	return nil
}
