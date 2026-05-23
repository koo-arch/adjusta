package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	opCookie "github.com/koo-arch/adjusta-backend/cookie"
	"github.com/koo-arch/adjusta-backend/ent"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/repo/account"
	repoSession "github.com/koo-arch/adjusta-backend/internal/repo/session"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"golang.org/x/oauth2"
)

type AuthManager struct {
	client      *ent.Client
	userRepo    user.UserRepository
	accountRepo repoAccount.AccountRepository
	sessionRepo repoSession.SessionRepository
}

func NewAuthManager(client *ent.Client, userRepo user.UserRepository, accountRepo repoAccount.AccountRepository, sessionRepo repoSession.SessionRepository) *AuthManager {
	return &AuthManager{
		client:      client,
		userRepo:    userRepo,
		accountRepo: accountRepo,
		sessionRepo: sessionRepo,
	}
}

func (am *AuthManager) ProcessUserSignIn(ctx context.Context, userInfo *userinfo.UserInfo, oauthToken *oauth2.Token) (*ent.User, error) {
	u, err := am.userRepo.FindByEmail(ctx, nil, userInfo.Email, user.UserQueryOptions{})
	if err != nil {
		log.Printf("failed to get user by email: %v", err)
		if ent.IsNotFound(err) {
			// ユーザーが存在しない場合は新規作成
			return am.CreateUser(ctx, userInfo, oauthToken)
		}
		// エラーが発生した場合
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	accountInfo, err := am.accountRepo.FindByUserID(ctx, nil, u.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return am.UpdateUserAndCreateAccount(ctx, u.ID, userInfo, oauthToken)
		}
		log.Printf("failed to get account by user id: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	// ユーザーが存在する場合はログイン情報を更新
	return am.UpdateUserAndAccount(ctx, u.ID, accountInfo.ID, userInfo, oauthToken)
}

func (am *AuthManager) CreateUser(ctx context.Context, userInfo *userinfo.UserInfo, oauthToken *oauth2.Token) (*ent.User, error) {
	// トランザクションを開始
	tx, err := am.client.Tx(ctx)
	if err != nil {
		log.Printf("failed starting transaction: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	// トランザクションエラー用の変数を別に定義
	defer transaction.HandleTransaction(tx, &err)

	u, err := am.userRepo.Create(ctx, tx, userInfo.Email, buildUserMutationOptions(userInfo))
	if err != nil {
		log.Printf("failed to create user: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	_, err = am.accountRepo.Create(ctx, tx, u.ID, buildAccountMutationOptions(userInfo, oauthToken))
	if err != nil {
		log.Printf("failed to create account: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return u, nil
}

func (am *AuthManager) UpdateUserAndCreateAccount(ctx context.Context, userID uuid.UUID, userInfo *userinfo.UserInfo, oauthToken *oauth2.Token) (*ent.User, error) {
	tx, err := am.client.Tx(ctx)
	if err != nil {
		log.Printf("failed starting transaction: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	defer transaction.HandleTransaction(tx, &err)

	u, err := am.userRepo.Update(ctx, tx, userID, buildUserMutationOptions(userInfo))
	if err != nil {
		log.Printf("failed to update user during account creation: %v", err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "ユーザーが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	_, err = am.accountRepo.Create(ctx, tx, userID, buildAccountMutationOptions(userInfo, oauthToken))
	if err != nil {
		log.Printf("failed to create account for existing user: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return u, nil
}

func (am *AuthManager) UpdateUserAndAccount(ctx context.Context, userID, accountID uuid.UUID, userInfo *userinfo.UserInfo, oauthToken *oauth2.Token) (*ent.User, error) {
	tx, err := am.client.Tx(ctx)
	if err != nil {
		log.Printf("failed starting transaction: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	defer transaction.HandleTransaction(tx, &err)

	u, err := am.userRepo.Update(ctx, tx, userID, buildUserMutationOptions(userInfo))
	if err != nil {
		log.Printf("failed to update user login state: %v", err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "ユーザーが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	_, err = am.accountRepo.Update(ctx, tx, accountID, buildAccountMutationOptions(userInfo, oauthToken))
	if err != nil {
		log.Printf("failed to update account token: %v", err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "アカウント情報が見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return u, nil

}

func (am *AuthManager) CreateSession(ctx context.Context, userID uuid.UUID) (*ent.Session, error) {
	sessionToken := uuid.NewString()
	expiresAt := buildSessionExpiry()

	entSession, err := am.sessionRepo.Create(ctx, nil, userID, sessionToken, expiresAt)
	if err != nil {
		log.Printf("failed to create session: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return entSession, nil
}

func (am *AuthManager) AuthenticateSession(ctx context.Context, sessionToken string) (*ent.User, error) {
	entSession, err := am.sessionRepo.FindByToken(ctx, nil, sessionToken, repoSession.SessionQueryOptions{WithUser: true})
	if err != nil {
		log.Printf("failed to find session by token: %v", err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusUnauthorized, "認証情報がありません")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	if entSession.ExpiresAt.Before(time.Now()) {
		if deleteErr := am.sessionRepo.DeleteByToken(ctx, nil, sessionToken); deleteErr != nil {
			log.Printf("failed to delete expired session: %v", deleteErr)
		}
		return nil, internalErrors.NewAPIError(http.StatusUnauthorized, "セッションの有効期限が切れています")
	}

	if _, err := am.sessionRepo.UpdateExpiry(ctx, nil, entSession.ID, buildSessionExpiry()); err != nil {
		log.Printf("failed to update session expiry: %v", err)
	}

	userEntity, err := entSession.Edges.UserOrErr()
	if err != nil {
		log.Printf("failed to load user from session: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusUnauthorized, "ユーザー情報が取得できませんでした")
	}

	return userEntity, nil
}

func (am *AuthManager) DeleteSession(ctx context.Context, sessionToken string) error {
	if sessionToken == "" {
		return nil
	}

	if err := am.sessionRepo.DeleteByToken(ctx, nil, sessionToken); err != nil {
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
