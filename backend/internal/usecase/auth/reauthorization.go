package auth

import (
	"context"
	"log"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/google"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (am *Authenticator) ReauthorizeGoogle(ctx context.Context, userID uuid.UUID, userInfo *google.UserProfile, oauthToken *google.AuthToken) error {
	if oauthToken.RefreshToken == "" {
		return internalErrors.NewGoogleReauthorizationRequiredError("Googleからリフレッシュトークンを取得できませんでした")
	}

	user, err := am.repos.User.Read(ctx, userID)
	if err != nil {
		log.Printf("failed to get user during google reauthorization: %v", err)
		if repoerr.IsNotFound(err) {
			return internalErrors.NewUnauthorizedError("ユーザー情報が取得できませんでした")
		}
		return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	account, err := am.repos.Account.FindByUserID(ctx, userID)
	if err != nil {
		log.Printf("failed to get account during google reauthorization: %v", err)
		if repoerr.IsNotFound(err) {
			return internalErrors.NewNotFoundError("アカウント情報が見つかりませんでした")
		}
		return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	if user.Email != userInfo.Email || account.GoogleUserID != userInfo.GoogleID {
		return internalErrors.NewForbiddenError("現在連携中のGoogleアカウントで再認可してください")
	}

	return am.tx.Do(ctx, func(repos AuthTxRepositories) error {
		if _, err := repos.User.Update(ctx, userID, buildUserMutationOptions(userInfo)); err != nil {
			log.Printf("failed to update user during google reauthorization: %v", err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}
		if _, err := repos.Account.Update(ctx, account.ID, buildAccountMutationOptions(userInfo, oauthToken)); err != nil {
			log.Printf("failed to update account during google reauthorization: %v", err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}
		return nil
	})
}
