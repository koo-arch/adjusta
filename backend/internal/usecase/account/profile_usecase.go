package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/google"
)

type GoogleProfile = google.UserProfile

type GoogleTokenProvider interface {
	GetToken(ctx context.Context, userID uuid.UUID) (*google.AuthToken, error)
}

type UserInfoFetcher interface {
	FetchGoogleUserInfo(ctx context.Context, token *google.AuthToken) (*google.UserProfile, error)
}

type UserInfoFetcherFunc func(ctx context.Context, token *google.AuthToken) (*google.UserProfile, error)

func (f UserInfoFetcherFunc) FetchGoogleUserInfo(ctx context.Context, token *google.AuthToken) (*google.UserProfile, error) {
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

	return userInfo, nil
}
