package auth

import (
	"context"
	"log"

	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

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
