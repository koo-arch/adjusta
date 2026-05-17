package user

import (
	"fmt"
)

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	RefreshToken string `json:"refresh_token"`
}

func NewUser (id int, email, refreshToken string) *User {
	return &User{
		ID: id,
		Email: email,
		RefreshToken: refreshToken,
	}
}

func (u *User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}