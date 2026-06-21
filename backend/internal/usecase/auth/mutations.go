package auth

import (
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
)

type UserMutation = repoUser.UserMutationOptions
type AccountMutation = repoAccount.AccountMutationOptions

func buildUserMutationOptions(userInfo *appmodel.GoogleUserProfile) UserMutation {
	return UserMutation{
		Name:      nullableString(userInfo.Name),
		AvatarURL: nullableString(userInfo.Picture),
	}
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
