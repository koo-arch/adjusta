package oauth

import (
	"context"

	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
)

type OAuthUsecase interface {
	GoogleLoginURL(state string) string
	CompleteGoogleSignIn(ctx context.Context, code string) (*usecaseAuth.GoogleSignInResult, error)
	Logout(ctx context.Context, sessionToken string) error
}
