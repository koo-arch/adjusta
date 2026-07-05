package googlecalendar

import (
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/google"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/googleapierror"
	"golang.org/x/oauth2"
)

func toOAuth2Token(token *google.AuthToken) *oauth2.Token {
	if token == nil {
		return nil
	}

	return &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
}

func normalizeGoogleAPIError(err error) error {
	if err == nil {
		return nil
	}

	return internalErrors.NormalizeAPIError(googleapierror.Normalize(err), "Google APIエラーが発生しました")
}
