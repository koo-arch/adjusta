package auth

import (
	"context"

	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	googleOAuth "github.com/koo-arch/adjusta-backend/internal/google/oauth"
	googleUserInfo "github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	"golang.org/x/oauth2"
)

type GoogleOAuthGateway struct{}

func NewGoogleOAuthGateway() *GoogleOAuthGateway {
	return &GoogleOAuthGateway{}
}

func (g *GoogleOAuthGateway) AuthCodeURL(state string) string {
	return googleOAuth.GetGoogleAuthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (g *GoogleOAuthGateway) Exchange(ctx context.Context, code string) (*appmodel.GoogleAuthToken, error) {
	token, err := googleOAuth.GetGoogleAuthConfig().Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return &appmodel.GoogleAuthToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		Scope:        tokenScope(token),
	}, nil
}

type GoogleUserInfoFetcher struct{}

func NewGoogleUserInfoFetcher() *GoogleUserInfoFetcher {
	return &GoogleUserInfoFetcher{}
}

func (f *GoogleUserInfoFetcher) FetchGoogleUserInfo(ctx context.Context, token *appmodel.GoogleAuthToken) (*appmodel.GoogleUserProfile, error) {
	userInfo, err := googleUserInfo.FetchGoogleUserInfo(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})
	if err != nil {
		return nil, err
	}

	return &appmodel.GoogleUserProfile{
		GoogleID: userInfo.GoogleID,
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		Picture:  userInfo.Picture,
	}, nil
}

func tokenScope(token *oauth2.Token) *string {
	rawScope := token.Extra("scope")
	scope, ok := rawScope.(string)
	if !ok || scope == "" {
		return nil
	}
	return &scope
}
