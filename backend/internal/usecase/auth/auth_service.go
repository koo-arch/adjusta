package auth

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	opCookie "github.com/koo-arch/adjusta-backend/cookie"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
	"github.com/koo-arch/adjusta-backend/internal/repositorymodel"
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

func (am *AuthService) ProcessUserSignIn(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repositorymodel.User, error) {
	u, err := am.reader.FindUserByEmail(ctx, userInfo.Email)
	if err != nil {
		log.Printf("failed to get user by email: %v", err)
		if repoerr.IsNotFound(err) {
			// ユーザーが存在しない場合は新規作成
			return am.CreateUser(ctx, userInfo, oauthToken)
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	accountInfo, err := am.reader.FindAccountByUserID(ctx, u.ID)
	if err != nil {
		if repoerr.IsNotFound(err) {
			return am.UpdateUserAndCreateAccount(ctx, u.ID, userInfo, oauthToken)
		}
		log.Printf("failed to get account by user id: %v", err)
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	// ユーザーが存在する場合はログイン情報を更新
	return am.UpdateUserAndAccount(ctx, u.ID, accountInfo.ID, userInfo, oauthToken)
}

func (am *AuthService) CreateUser(ctx context.Context, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repositorymodel.User, error) {
	var u *repositorymodel.User

	err := am.tx.Do(ctx, func(store SignInStore) error {
		var err error
		u, err = store.CreateUser(ctx, userInfo.Email, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to create user: %v", err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		_, err = store.CreateAccount(ctx, u.ID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to create account: %v", err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (am *AuthService) UpdateUserAndCreateAccount(ctx context.Context, userID uuid.UUID, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repositorymodel.User, error) {
	var u *repositorymodel.User

	err := am.tx.Do(ctx, func(store SignInStore) error {
		var err error
		u, err = store.UpdateUser(ctx, userID, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to update user during account creation: %v", err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("ユーザーが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		_, err = store.CreateAccount(ctx, userID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to create account for existing user: %v", err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (am *AuthService) UpdateUserAndAccount(ctx context.Context, userID, accountID uuid.UUID, userInfo *appmodel.GoogleUserProfile, oauthToken *appmodel.GoogleAuthToken) (*repositorymodel.User, error) {
	var u *repositorymodel.User

	err := am.tx.Do(ctx, func(store SignInStore) error {
		var err error
		u, err = store.UpdateUser(ctx, userID, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to update user login state: %v", err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("ユーザーが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		_, err = store.UpdateAccount(ctx, accountID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to update account token: %v", err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("アカウント情報が見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u, nil

}

func (am *AuthService) CreateSession(ctx context.Context, userID uuid.UUID) (*repositorymodel.Session, error) {
	sessionToken := uuid.NewString()
	expiresAt := buildSessionExpiry()

	entSession, err := am.sessions.CreateSession(ctx, userID, sessionToken, expiresAt)
	if err != nil {
		log.Printf("failed to create session: %v", err)
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return entSession, nil
}

func (am *AuthService) AuthenticateSession(ctx context.Context, sessionToken string) (*repositorymodel.User, error) {
	entSession, err := am.sessions.FindSessionByToken(ctx, sessionToken, true)
	if err != nil {
		log.Printf("failed to find session by token: %v", err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewUnauthorizedError("認証情報がありません")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	if entSession.ExpiresAt.Before(time.Now()) {
		if deleteErr := am.sessions.DeleteSessionByToken(ctx, sessionToken); deleteErr != nil {
			log.Printf("failed to delete expired session: %v", deleteErr)
		}
		return nil, internalErrors.NewUnauthorizedError("セッションの有効期限が切れています")
	}

	if _, err := am.sessions.UpdateSessionExpiry(ctx, entSession.ID, buildSessionExpiry()); err != nil {
		log.Printf("failed to update session expiry: %v", err)
	}

	if entSession.User == nil {
		log.Printf("failed to load user from session: user edge is nil")
		return nil, internalErrors.NewUnauthorizedError("ユーザー情報が取得できませんでした")
	}

	return entSession.User, nil
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

func buildSessionExpiry() time.Time {
	return time.Now().Add(time.Duration(opCookie.DefaultCookieOptions().MaxAge) * time.Second)
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
