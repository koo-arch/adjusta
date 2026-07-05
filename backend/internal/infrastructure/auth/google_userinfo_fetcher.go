package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/google"
	infraGoogleOAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/googleoauth"
	"golang.org/x/oauth2"
)

type googleUserInfo struct {
	GoogleID string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}

const googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"

type GoogleUserInfoFetcher struct {
	client *infraGoogleOAuth.Client
}

func NewGoogleUserInfoFetcher(client *infraGoogleOAuth.Client) *GoogleUserInfoFetcher {
	return &GoogleUserInfoFetcher{client: client}
}

func (f *GoogleUserInfoFetcher) FetchGoogleUserInfo(ctx context.Context, token *google.AuthToken) (*google.UserProfile, error) {
	client := f.client.HTTPClient(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})

	resp, err := client.Get(googleUserInfoEndpoint)
	if err != nil {
		log.Printf("failed to fetch user info: %v", err)
		return nil, internalErrors.NewInternalError("GoogleAPIからのユーザー情報取得に失敗しました")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to fetch user info: %s", resp.Status)
		return nil, userInfoStatusError(resp.StatusCode)
	}

	var userInfo googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return &google.UserProfile{
		GoogleID: userInfo.GoogleID,
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		Picture:  userInfo.Picture,
	}, nil
}

func userInfoStatusError(statusCode int) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return internalErrors.NewUnauthorizedError("ユーザー情報の取得に失敗しました")
	case http.StatusForbidden:
		return internalErrors.NewForbiddenError("ユーザー情報の取得に失敗しました")
	case http.StatusNotFound:
		return internalErrors.NewNotFoundError("ユーザー情報の取得に失敗しました")
	default:
		return internalErrors.NewInternalError("ユーザー情報の取得に失敗しました")
	}
}
