package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	infraGoogleOAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/googleoauth"
	"golang.org/x/oauth2"
)

type googleUserInfo struct {
	GoogleID string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}

type GoogleUserInfoFetcher struct{}

func NewGoogleUserInfoFetcher() *GoogleUserInfoFetcher {
	return &GoogleUserInfoFetcher{}
}

func (f *GoogleUserInfoFetcher) FetchGoogleUserInfo(ctx context.Context, token *appmodel.GoogleAuthToken) (*appmodel.GoogleUserProfile, error) {
	client := infraGoogleOAuth.GetConfig().Client(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
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
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, internalErrors.NewUnauthorizedError("ユーザー情報の取得に失敗しました")
		case http.StatusForbidden:
			return nil, internalErrors.NewForbiddenError("ユーザー情報の取得に失敗しました")
		case http.StatusNotFound:
			return nil, internalErrors.NewNotFoundError("ユーザー情報の取得に失敗しました")
		default:
			return nil, internalErrors.NewInternalError("ユーザー情報の取得に失敗しました")
		}
	}

	var userInfo googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return &appmodel.GoogleUserProfile{
		GoogleID: userInfo.GoogleID,
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		Picture:  userInfo.Picture,
	}, nil
}
