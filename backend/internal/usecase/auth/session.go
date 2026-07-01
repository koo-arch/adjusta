package auth

import (
	"context"
	"log"
	"time"

	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (am *Authenticator) AuthenticateSession(ctx context.Context, sessionToken string) (*repoUser.User, error) {
	storedSession, err := am.repos.Session.FindByToken(ctx, sessionToken, repoSession.SessionQueryOptions{
		WithUser: true,
	})
	if err != nil {
		log.Printf("failed to find session by token: %v", err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewUnauthorizedError("認証情報がありません")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	now := time.Now()
	if storedSession.IsExpiredAt(now) {
		if deleteErr := am.repos.Session.DeleteByToken(ctx, sessionToken); deleteErr != nil {
			log.Printf("failed to delete expired session: %v", deleteErr)
		}
		return nil, internalErrors.NewUnauthorizedError("セッションの有効期限が切れています")
	}

	if _, err := am.repos.Session.UpdateExpiry(ctx, storedSession.ID, repoSession.ExpiresAtFrom(now, am.sessionLifetime)); err != nil {
		log.Printf("failed to update session expiry: %v", err)
	}

	if storedSession.User == nil {
		log.Printf("failed to load user from session: user edge is nil")
		return nil, internalErrors.NewUnauthorizedError("ユーザー情報が取得できませんでした")
	}

	return storedSession.User, nil
}

func (am *Authenticator) DeleteSession(ctx context.Context, sessionToken string) error {
	if sessionToken == "" {
		return nil
	}

	if err := am.repos.Session.DeleteByToken(ctx, sessionToken); err != nil {
		log.Printf("failed to delete session: %v", err)
		return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return nil
}
