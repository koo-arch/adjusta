package auth

import (
	"time"
)

type Authenticator struct {
	repos           AuthRepositories
	tx              SignInTransaction
	sessionLifetime time.Duration
}

func NewAuthenticator(repos AuthRepositories, tx SignInTransaction, sessionLifetime time.Duration) *Authenticator {
	return &Authenticator{
		repos:           repos,
		tx:              tx,
		sessionLifetime: sessionLifetime,
	}
}
