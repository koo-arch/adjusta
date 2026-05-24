package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	opCookie "github.com/koo-arch/adjusta-backend/cookie"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	internalModels "github.com/koo-arch/adjusta-backend/internal/models"
	internalRepo "github.com/koo-arch/adjusta-backend/internal/repo"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/repo/account"
	"github.com/koo-arch/adjusta-backend/internal/repo/session"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
	"golang.org/x/oauth2"
)

type AuthManager struct {
	repos internalRepo.Repositories
	uow   internalRepo.UnitOfWork
}

func NewAuthManager(repos internalRepo.Repositories, uow internalRepo.UnitOfWork) *AuthManager {
	return &AuthManager{
		repos: repos,
		uow:   uow,
	}
}

func (am *AuthManager) ProcessUserSignIn(ctx context.Context, userInfo *userinfo.UserInfo, oauthToken *oauth2.Token) (*internalModels.User, error) {
	u, err := am.repos.User.FindByEmail(ctx, userInfo.Email, user.UserQueryOptions{})
	if err != nil {
		log.Printf("failed to get user by email: %v", err)
		if repoerr.IsNotFound(err) {
			// ユーザーが存在しない場合は新規作成
			return am.CreateUser(ctx, userInfo, oauthToken)
		}
		// エラーが発生した場合
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	accountInfo, err := am.repos.Account.FindByUserID(ctx, u.ID)
	if err != nil {
		if repoerr.IsNotFound(err) {
			return am.UpdateUserAndCreateAccount(ctx, u.ID, userInfo, oauthToken)
		}
		log.Printf("failed to get account by user id: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	// ユーザーが存在する場合はログイン情報を更新
	return am.UpdateUserAndAccount(ctx, u.ID, accountInfo.ID, userInfo, oauthToken)
}

func (am *AuthManager) CreateUser(ctx context.Context, userInfo *userinfo.UserInfo, oauthToken *oauth2.Token) (*internalModels.User, error) {
	var u *internalModels.User

	err := am.uow.Do(ctx, func(repos internalRepo.Repositories) error {
		var err error
		u, err = repos.User.Create(ctx, userInfo.Email, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to create user: %v", err)
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		_, err = repos.Account.Create(ctx, u.ID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to create account: %v", err)
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (am *AuthManager) UpdateUserAndCreateAccount(ctx context.Context, userID uuid.UUID, userInfo *userinfo.UserInfo, oauthToken *oauth2.Token) (*internalModels.User, error) {
	var u *internalModels.User

	err := am.uow.Do(ctx, func(repos internalRepo.Repositories) error {
		var err error
		u, err = repos.User.Update(ctx, userID, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to update user during account creation: %v", err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewAPIError(http.StatusNotFound, "ユーザーが見つかりませんでした")
			}
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		_, err = repos.Account.Create(ctx, userID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to create account for existing user: %v", err)
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (am *AuthManager) UpdateUserAndAccount(ctx context.Context, userID, accountID uuid.UUID, userInfo *userinfo.UserInfo, oauthToken *oauth2.Token) (*internalModels.User, error) {
	var u *internalModels.User

	err := am.uow.Do(ctx, func(repos internalRepo.Repositories) error {
		var err error
		u, err = repos.User.Update(ctx, userID, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to update user login state: %v", err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewAPIError(http.StatusNotFound, "ユーザーが見つかりませんでした")
			}
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		_, err = repos.Account.Update(ctx, accountID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to update account token: %v", err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewAPIError(http.StatusNotFound, "アカウント情報が見つかりませんでした")
			}
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u, nil

}

func (am *AuthManager) CreateSession(ctx context.Context, userID uuid.UUID) (*internalModels.Session, error) {
	sessionToken := uuid.NewString()
	expiresAt := buildSessionExpiry()

	entSession, err := am.repos.Session.Create(ctx, userID, sessionToken, expiresAt)
	if err != nil {
		log.Printf("failed to create session: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return entSession, nil
}

func (am *AuthManager) AuthenticateSession(ctx context.Context, sessionToken string) (*internalModels.User, error) {
	entSession, err := am.repos.Session.FindByToken(ctx, sessionToken, session.SessionQueryOptions{WithUser: true})
	if err != nil {
		log.Printf("failed to find session by token: %v", err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusUnauthorized, "認証情報がありません")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	if entSession.ExpiresAt.Before(time.Now()) {
		if deleteErr := am.repos.Session.DeleteByToken(ctx, sessionToken); deleteErr != nil {
			log.Printf("failed to delete expired session: %v", deleteErr)
		}
		return nil, internalErrors.NewAPIError(http.StatusUnauthorized, "セッションの有効期限が切れています")
	}

	if _, err := am.repos.Session.UpdateExpiry(ctx, entSession.ID, buildSessionExpiry()); err != nil {
		log.Printf("failed to update session expiry: %v", err)
	}

	if entSession.User == nil {
		log.Printf("failed to load user from session: user edge is nil")
		return nil, internalErrors.NewAPIError(http.StatusUnauthorized, "ユーザー情報が取得できませんでした")
	}

	return entSession.User, nil
}

func (am *AuthManager) DeleteSession(ctx context.Context, sessionToken string) error {
	if sessionToken == "" {
		return nil
	}

	if err := am.repos.Session.DeleteByToken(ctx, sessionToken); err != nil {
		log.Printf("failed to delete session: %v", err)
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return nil
}

func buildUserMutationOptions(userInfo *userinfo.UserInfo) user.UserMutationOptions {
	return user.UserMutationOptions{
		Name:      nullableString(userInfo.Name),
		AvatarURL: nullableString(userInfo.Picture),
	}
}

func buildSessionExpiry() time.Time {
	return time.Now().Add(time.Duration(opCookie.DefaultCookieOptions().MaxAge) * time.Second)
}

func buildAccountMutationOptions(userInfo *userinfo.UserInfo, oauthToken *oauth2.Token) repoAccount.AccountMutationOptions {
	return repoAccount.AccountMutationOptions{
		GoogleUserID: &userInfo.GoogleID,
		AccessToken:  &oauthToken.AccessToken,
		RefreshToken: nullableString(oauthToken.RefreshToken),
		ExpiresAt:    &oauthToken.Expiry,
		Scope:        tokenScope(oauthToken),
	}
}

func tokenScope(oauthToken *oauth2.Token) *string {
	rawScope := oauthToken.Extra("scope")
	scope, ok := rawScope.(string)
	if !ok || scope == "" {
		return nil
	}
	return &scope
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
