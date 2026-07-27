package taskauth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/api/idtoken"
)

type OIDCVerifier struct {
	audience      string
	expectedEmail string
}

func NewOIDCVerifier(audience, expectedEmail string) *OIDCVerifier {
	return &OIDCVerifier{audience: audience, expectedEmail: expectedEmail}
}

func (v *OIDCVerifier) VerifyAuthorization(ctx context.Context, authorization string) error {
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return errors.New("bearer token is required")
	}
	payload, err := idtoken.Validate(ctx, strings.TrimSpace(strings.TrimPrefix(authorization, prefix)), v.audience)
	if err != nil {
		return fmt.Errorf("invalid oidc token: %w", err)
	}
	if payload.Issuer != "https://accounts.google.com" && payload.Issuer != "accounts.google.com" {
		return errors.New("unexpected oidc token issuer")
	}
	email, _ := payload.Claims["email"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)
	if email != v.expectedEmail || !emailVerified {
		return errors.New("unexpected oidc token subject")
	}
	return nil
}
