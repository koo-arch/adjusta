package account

import (
	"context"

	"github.com/google/uuid"
	googleUserInfo "github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	"golang.org/x/oauth2"
)

type GoogleProfile struct {
	GoogleID string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}

type GoogleTokenProvider interface {
	GetToken(ctx context.Context, userID uuid.UUID) (*oauth2.Token, error)
}

type UserInfoFetcher interface {
	FetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*googleUserInfo.UserInfo, error)
}

type UserInfoFetcherFunc func(ctx context.Context, token *oauth2.Token) (*googleUserInfo.UserInfo, error)

func (f UserInfoFetcherFunc) FetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*googleUserInfo.UserInfo, error) {
	return f(ctx, token)
}

type ProfileUsecase struct {
	googleTokenProvider GoogleTokenProvider
	userInfoFetcher     UserInfoFetcher
}

func NewProfileUsecase(googleTokenProvider GoogleTokenProvider, userInfoFetcher UserInfoFetcher) *ProfileUsecase {
	return &ProfileUsecase{
		googleTokenProvider: googleTokenProvider,
		userInfoFetcher:     userInfoFetcher,
	}
}

func (uc *ProfileUsecase) FetchGoogleProfile(ctx context.Context, userID uuid.UUID) (*GoogleProfile, error) {
	token, err := uc.googleTokenProvider.GetToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	userInfo, err := uc.userInfoFetcher.FetchGoogleUserInfo(ctx, token)
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
