package auth

import (
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
)

type AuthRepositories struct {
	User    repoUser.UserRepository
	Account repoAccount.AccountRepository
	Session repoSession.SessionRepository
}
