package appmodel

import "time"

type GoogleAuthToken struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
	Scope        *string
}

type GoogleUserProfile struct {
	GoogleID string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}

type GoogleCalendarList struct {
	CalendarID string `json:"calendar_id"`
	Summary    string `json:"summary"`
	Primary    bool   `json:"primary"`
}
