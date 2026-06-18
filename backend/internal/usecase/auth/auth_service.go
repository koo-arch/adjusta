package auth

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	opCookie "github.com/koo-arch/adjusta-backend/cookie"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type AuthService struct {
	reader   SignInReader
	tx       SignInTransaction
	sessions SessionStore
}

func NewAuthService(reader SignInReader, tx SignInTransaction, sessions SessionStore) *AuthService {
	return &AuthService{
		reader:   reader,
		tx:       tx,
		sessions: sessions,
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
	expiresAt := repoSession.ExpiresAtFrom(time.Now(), sessionLifetime())

	err = am.tx.Do(ctx, func(store SignInStore) error {
		var err error
		u, err = am.persistSignIn(ctx, store, plan, userInfo, oauthToken)
		if err != nil {
			return err
		}

		storedSession, err = store.CreateSession(ctx, u.ID, sessionToken, expiresAt)
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
	err = am.tx.Do(ctx, func(store SignInStore) error {
		var err error
		u, err = am.persistSignIn(ctx, store, plan, userInfo, oauthToken)
		return err
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (am *AuthService) AuthenticateSession(ctx context.Context, sessionToken string) (*repoUser.User, error) {
	storedSession, err := am.sessions.FindSessionByToken(ctx, sessionToken, true)
	if err != nil {
		log.Printf("failed to find session by token: %v", err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewUnauthorizedError("認証情報がありません")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	now := time.Now()
	if storedSession.IsExpiredAt(now) {
		if deleteErr := am.sessions.DeleteSessionByToken(ctx, sessionToken); deleteErr != nil {
			log.Printf("failed to delete expired session: %v", deleteErr)
		}
		return nil, internalErrors.NewUnauthorizedError("セッションの有効期限が切れています")
	}

	if _, err := am.sessions.UpdateSessionExpiry(ctx, storedSession.ID, repoSession.ExpiresAtFrom(now, sessionLifetime())); err != nil {
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

	if err := am.sessions.DeleteSessionByToken(ctx, sessionToken); err != nil {
		log.Printf("failed to delete session: %v", err)
		return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return nil
}

func buildUserMutationOptions(userInfo *appmodel.GoogleUserProfile) UserMutation {
	return UserMutation{
		Name:      nullableString(userInfo.Name),
		AvatarURL: nullableString(userInfo.Picture),
	}
}

func sessionLifetime() time.Duration {
	return time.Duration(opCookie.DefaultCookieOptions().MaxAge) * time.Second
}

func buildAccountMutationOptions(userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) AccountMutation {
	return AccountMutation{
		GoogleUserID: &userInfo.GoogleID,
		AccessToken:  &oauthToken.AccessToken,
		RefreshToken: nullableString(oauthToken.RefreshToken),
		ExpiresAt:    &oauthToken.Expiry,
		Scope:        oauthToken.Scope,
	}
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

type signInPlan struct {
	userID    uuid.UUID
	accountID uuid.UUID
	mode      signInMode
}

type signInMode int

const (
	signInModeCreateUserAndAccount signInMode = iota
	signInModeUpdateUserAndCreateAccount
	signInModeUpdateUserAndAccount
)

func (am *AuthService) resolveSignInPlan(ctx context.Context, userInfo *appmodel.GoogleUserProfile) (*signInPlan, error) {
	u, err := am.reader.FindUserByEmail(ctx, userInfo.Email)
	if err != nil {
		log.Printf("failed to get user by email: %v", err)
		if repoerr.IsNotFound(err) {
			return &signInPlan{mode: signInModeCreateUserAndAccount}, nil
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	accountInfo, err := am.reader.FindAccountByUserID(ctx, u.ID)
	if err != nil {
		if repoerr.IsNotFound(err) {
			return &signInPlan{
				userID: u.ID,
				mode:   signInModeUpdateUserAndCreateAccount,
			}, nil
		}
		log.Printf("failed to get account by user id: %v", err)
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return &signInPlan{
		userID:    u.ID,
		accountID: accountInfo.ID,
		mode:      signInModeUpdateUserAndAccount,
	}, nil
}

func (am *AuthService) persistSignIn(ctx context.Context, store SignInStore, plan *signInPlan, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repoUser.User, error) {
	switch plan.mode {
	case signInModeCreateUserAndAccount:
		u, err := store.CreateUser(ctx, userInfo.Email, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to create user: %v", err)
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		_, err = store.CreateAccount(ctx, u.ID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to create account: %v", err)
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return u, nil
	case signInModeUpdateUserAndCreateAccount:
		u, err := store.UpdateUser(ctx, plan.userID, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to update user during account creation: %v", err)
			if repoerr.IsNotFound(err) {
				return nil, internalErrors.NewNotFoundError("ユーザーが見つかりませんでした")
			}
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		_, err = store.CreateAccount(ctx, plan.userID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to create account for existing user: %v", err)
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return u, nil
	case signInModeUpdateUserAndAccount:
		u, err := store.UpdateUser(ctx, plan.userID, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to update user login state: %v", err)
			if repoerr.IsNotFound(err) {
				return nil, internalErrors.NewNotFoundError("ユーザーが見つかりませんでした")
			}
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		_, err = store.UpdateAccount(ctx, plan.accountID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to update account token: %v", err)
			if repoerr.IsNotFound(err) {
				return nil, internalErrors.NewNotFoundError("アカウント情報が見つかりませんでした")
			}
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return u, nil
	default:
		log.Printf("unexpected sign in mode: %d", plan.mode)
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}
}
