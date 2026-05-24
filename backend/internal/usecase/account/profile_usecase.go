package account

import (
	"context"

	"github.com/google/uuid"
	googleOAuth "github.com/koo-arch/adjusta-backend/internal/google/oauth"
	googleUserInfo "github.com/koo-arch/adjusta-backend/internal/google/userinfo"
)

type GoogleProfile struct {
	GoogleID string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}

type ProfileUsecase struct {
	googleTokenManager *googleOAuth.TokenManager
}

func NewProfileUsecase(googleTokenManager *googleOAuth.TokenManager) *ProfileUsecase {
	return &ProfileUsecase{
		googleTokenManager: googleTokenManager,
	}
}

func (uc *ProfileUsecase) FetchGoogleProfile(ctx context.Context, userID uuid.UUID) (*GoogleProfile, error) {
	token, err := uc.googleTokenManager.GetToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	userInfo, err := googleUserInfo.FetchGoogleUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}

	return &GoogleProfile{
		GoogleID: userInfo.GoogleID,
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		Picture:  userInfo.Picture,
	}, nil
}
