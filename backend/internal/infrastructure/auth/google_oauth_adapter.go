package auth

import (
	"context"

	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	infraGoogleOAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/googleoauth"
	"golang.org/x/oauth2"
)

type GoogleOAuthGateway struct{}

func NewGoogleOAuthGateway() *GoogleOAuthGateway {
	return &GoogleOAuthGateway{}
}

func (g *GoogleOAuthGateway) AuthCodeURL(state string) string {
	return infraGoogleOAuth.GetConfig().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (g *GoogleOAuthGateway) Exchange(ctx context.Context, code string) (*appmodel.GoogleAuthToken, error) {
	token, err := infraGoogleOAuth.GetConfig().Exchange(ctx, code)
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

func tokenScope(token *oauth2.Token) *string {
	rawScope := token.Extra("scope")
	scope, ok := rawScope.(string)
	if !ok || scope == "" {
		return nil
	}
	return &scope
}
