package googleoauth

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Client struct {
	config *oauth2.Config
}

func NewClient(clientID, clientSecret, redirectURL string) *Client {
	return &Client{config: &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/calendar",
			"openid",
		},
		Endpoint: google.Endpoint,
	}}
}

func (c *Client) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return c.config.AuthCodeURL(state, opts...)
}

func (c *Client) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return c.config.Exchange(ctx, code, opts...)
}

func (c *Client) HTTPClient(ctx context.Context, token *oauth2.Token) *http.Client {
	return c.config.Client(ctx, token)
}

func (c *Client) TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	return c.config.TokenSource(ctx, token)
}
