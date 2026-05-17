package models

import (
	"time"
)

type JWTToken struct {
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	AccessExpiration   time.Time `json:"access_expiration"`
	RefreshExpiration  time.Time `json:"refresh_expiration"`
}

type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}