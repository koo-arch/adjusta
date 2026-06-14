package oauth

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
	"golang.org/x/oauth2"
)

type TokenManager struct {
	accountRepo repoAccount.AccountRepository
}

func NewTokenManager(accountRepo repoAccount.AccountRepository) *TokenManager {
	return &TokenManager{
		accountRepo: accountRepo,
	}
}

func (tm *TokenManager) GetToken(ctx context.Context, userID uuid.UUID) (*appmodel.GoogleAuthToken, error) {
	entAccount, err := tm.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		log.Printf("failed to get account by user id: %v", err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("アカウント情報が見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	if entAccount.AccessToken == nil || entAccount.RefreshToken == nil || entAccount.ExpiresAt == nil {
		return nil, internalErrors.NewUnauthorizedError("Googleアカウントの連携情報が不足しています。再認証してください")
	}

	token := &oauth2.Token{
		AccessToken:  *entAccount.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: *entAccount.RefreshToken,
		Expiry:       *entAccount.ExpiresAt,
	}

	refreshedToken, err := oauth2.ReuseTokenSource(token, GetGoogleAuthConfig().TokenSource(ctx, token)).Token()
	if err != nil {
		log.Printf("failed to get new token: %v", err)
		if strings.Contains(err.Error(), "invalid_token") {
			return nil, internalErrors.NewUnauthorizedError("トークンが無効です。再認証してください")
		}
		if strings.Contains(err.Error(), "insufficient_scope") {
			return nil, internalErrors.NewForbiddenError("トークンのスコープが不足しています。再認証してください")
		}
		if strings.Contains(err.Error(), "network error") {
			return nil, internalErrors.NewBadGatewayError("トークン取得サーバーに接続できませんでした")
		}
		return nil, internalErrors.NewInternalError("トークンの再取得中にエラーが発生しました")
	}

	if shouldPersistRefreshedToken(entAccount, refreshedToken) {
		_, err = tm.accountRepo.Update(ctx, entAccount.ID, buildAccountTokenRefreshOptions(entAccount, refreshedToken))
		if err != nil {
			log.Printf("failed to update token: %v", err)
			if repoerr.IsNotFound(err) {
				return nil, internalErrors.NewNotFoundError("アカウント情報が見つかりませんでした")
			}
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}
	}

	return buildGoogleAuthToken(refreshedToken), nil
}

func buildGoogleAuthToken(token *oauth2.Token) *appmodel.GoogleAuthToken {
	return &appmodel.GoogleAuthToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		Scope:        tokenScope(token),
	}
}

func buildAccountTokenRefreshOptions(account *repoAccount.Account, oauthToken *oauth2.Token) repoAccount.AccountMutationOptions {
	refreshToken := account.RefreshToken
	if oauthToken.RefreshToken != "" {
		refreshToken = &oauthToken.RefreshToken
	}

	opt := repoAccount.AccountMutationOptions{
		GoogleUserID: &account.GoogleUserID,
		AccessToken:  &oauthToken.AccessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    &oauthToken.Expiry,
		Scope:        account.Scope,
	}

	if scope := tokenScope(oauthToken); scope != nil {
		opt.Scope = scope
	}

	return opt
}

func shouldPersistRefreshedToken(account *repoAccount.Account, oauthToken *oauth2.Token) bool {
	if account.AccessToken == nil || *account.AccessToken != oauthToken.AccessToken {
		return true
	}

	if oauthToken.RefreshToken != "" && (account.RefreshToken == nil || *account.RefreshToken != oauthToken.RefreshToken) {
		return true
	}

	if account.ExpiresAt == nil || !account.ExpiresAt.Equal(oauthToken.Expiry) {
		return true
	}

	if scope := tokenScope(oauthToken); scope != nil {
		if account.Scope == nil || *account.Scope != *scope {
			return true
		}
	}

	return false
}

func tokenScope(oauthToken *oauth2.Token) *string {
	rawScope := oauthToken.Extra("scope")
	scope, ok := rawScope.(string)
	if !ok || scope == "" {
		return nil
	}
	return &scope
}
