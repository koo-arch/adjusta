package auth

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type AuthService struct {
	repos           AuthRepositories
	tx              SignInTransaction
	sessionLifetime time.Duration
}

func NewAuthService(repos AuthRepositories, tx SignInTransaction, sessionLifetime time.Duration) *AuthService {
	return &AuthService{
		repos:           repos,
		tx:              tx,
		sessionLifetime: sessionLifetime,
	}
}

func (am *AuthService) SignInWithGoogle(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repoSession.Session, *repoUser.User, error) {
	plan, err := am.resolveSignInPlan(ctx, userInfo)
	if err != nil {
		return nil, nil, err
	}

	var (
		u             *repoUser.User
		storedSession *repoSession.Session
	)

	sessionToken := uuid.NewString()
	expiresAt := repoSession.ExpiresAtFrom(time.Now(), am.sessionLifetime)

	err = am.tx.Do(ctx, func(repos AuthTxRepositories) error {
		var err error
		u, err = am.persistSignIn(ctx, repos, plan, userInfo, oauthToken)
		if err != nil {
			return err
		}

		storedSession, err = repos.Session.Create(ctx, u.ID, sessionToken, expiresAt)
		if err != nil {
			log.Printf("failed to create session during google sign in: %v", err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return storedSession, u, nil
}

func (am *AuthService) ProcessUserSignIn(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repoUser.User, error) {
	plan, err := am.resolveSignInPlan(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	var u *repoUser.User
	err = am.tx.Do(ctx, func(repos AuthTxRepositories) error {
		var err error
		u, err = am.persistSignIn(ctx, repos, plan, userInfo, oauthToken)
		return err
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (am *AuthService) AuthenticateSession(ctx context.Context, sessionToken string) (*repoUser.User, error) {
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

func (am *AuthService) DeleteSession(ctx context.Context, sessionToken string) error {
	if sessionToken == "" {
		return nil
	}

	if err := am.repos.Session.DeleteByToken(ctx, sessionToken); err != nil {
		log.Printf("failed to delete session: %v", err)
		return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return nil
}
