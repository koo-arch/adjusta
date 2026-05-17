package userinfo

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/koo-arch/adjusta-backend/internal/google/oauth"
	"golang.org/x/oauth2"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

type UserInfo struct {
	GoogleID string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}

func FetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	client := oauth.GoogleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Printf("failed to fetch user info: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, "GoogleAPIからのユーザー情報取得に失敗しました")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed closing response body: %v", err)
		}
	}()

	// レスポンスのステータスコードが200以外の場合はエラー
	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to fetch user info: %s", resp.Status)
		return nil, internalErrors.NewAPIError(resp.StatusCode, "ユーザー情報の取得に失敗しました")
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return &userInfo, nil
}