package oauth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/repo/account"
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

func (tm *TokenManager) GetToken(ctx context.Context, userID uuid.UUID) (*oauth2.Token, error) {
	entAccount, err := tm.accountRepo.FindByUserID(ctx, nil, userID)
	if err != nil {
		log.Printf("failed to get account by user id: %v", err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "アカウント情報が見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	if entAccount.AccessToken == nil || entAccount.RefreshToken == nil || entAccount.ExpiresAt == nil {
		return nil, internalErrors.NewAPIError(http.StatusUnauthorized, "Googleアカウントの連携情報が不足しています。再認証してください")
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
			return nil, internalErrors.NewAPIError(http.StatusUnauthorized, "トークンが無効です。再認証してください")
		}
		if strings.Contains(err.Error(), "insufficient_scope") {
			return nil, internalErrors.NewAPIError(http.StatusForbidden, "トークンのスコープが不足しています。再認証してください")
		}
		if strings.Contains(err.Error(), "network error") {
			return nil, internalErrors.NewAPIError(http.StatusBadGateway, "トークン取得サーバーに接続できませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, "トークンの再取得中にエラーが発生しました")
	}

	if shouldPersistRefreshedToken(entAccount, refreshedToken) {
		_, err = tm.accountRepo.Update(ctx, nil, entAccount.ID, buildAccountTokenRefreshOptions(entAccount, refreshedToken))
		if err != nil {
			log.Printf("failed to update token: %v", err)
			if ent.IsNotFound(err) {
				return nil, internalErrors.NewAPIError(http.StatusNotFound, "アカウント情報が見つかりませんでした")
			}
			return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}
	}

	return refreshedToken, nil
}

func buildAccountTokenRefreshOptions(account *ent.Account, oauthToken *oauth2.Token) repoAccount.AccountMutationOptions {
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

func shouldPersistRefreshedToken(account *ent.Account, oauthToken *oauth2.Token) bool {
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
