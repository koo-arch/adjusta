package auth

import (
	"context"
	"log"
	"net/http"

	internalAuth "github.com/koo-arch/adjusta-backend/internal/auth"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	googleOAuth "github.com/koo-arch/adjusta-backend/internal/google/oauth"
	googleUserInfo "github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	"github.com/koo-arch/adjusta-backend/utils"
)

type GoogleSignInResult struct {
	SessionToken string
	UserEmail    string
}

type SessionUsecase struct {
	authManager *internalAuth.AuthManager
}

func NewSessionUsecase(authManager *internalAuth.AuthManager) *SessionUsecase {
	return &SessionUsecase{
		authManager: authManager,
	}
}

func (uc *SessionUsecase) CompleteGoogleSignIn(ctx context.Context, code string) (*GoogleSignInResult, error) {
	oauthToken, err := googleOAuth.GetGoogleAuthConfig().Exchange(ctx, code)
	if err != nil {
		log.Printf("failed to exchange oauth token: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, "OAuthトークンの取得に失敗しました")
	}

	userInfo, err := googleUserInfo.FetchGoogleUserInfo(ctx, oauthToken)
	if err != nil {
		log.Printf("failed to fetch user info: %v", err)
		return nil, utils.GetAPIError(err, "ユーザー情報の取得に失敗しました")
	}

	u, err := uc.authManager.ProcessUserSignIn(ctx, userInfo, oauthToken)
	if err != nil {
		log.Printf("failed to create or update user: %v", err)
		return nil, utils.GetAPIError(err, "ユーザーの作成または更新に失敗しました")
	}

	entSession, err := uc.authManager.CreateSession(ctx, u.ID)
	if err != nil {
		log.Printf("failed to create app session: %v", err)
		return nil, utils.GetAPIError(err, "ログインセッションの作成に失敗しました")
	}

	return &GoogleSignInResult{
		SessionToken: entSession.SessionToken,
		UserEmail:    u.Email,
	}, nil
}

func (uc *SessionUsecase) Logout(ctx context.Context, sessionToken string) error {
	if err := uc.authManager.DeleteSession(ctx, sessionToken); err != nil {
		log.Printf("failed to delete session: %v", err)
		return utils.GetAPIError(err, "セッションの削除に失敗しました")
	}

	return nil
}
