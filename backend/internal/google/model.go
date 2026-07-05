package google

import "time"

type AuthToken struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
	Scope        *string
}

type UserProfile struct {
	GoogleID string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}
