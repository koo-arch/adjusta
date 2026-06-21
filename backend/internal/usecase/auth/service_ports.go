package auth

import (
	"context"

	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
)

type UserMutation = repoUser.UserMutationOptions
type AccountMutation = repoAccount.AccountMutationOptions

type AuthRepositories struct {
	User    repoUser.UserRepository
	Account repoAccount.AccountRepository
	Session repoSession.SessionRepository
}

type AuthTxRepositories struct {
	User    repoUser.UserRepository
	Account repoAccount.AccountRepository
	Session repoSession.SessionRepository
}

type SignInTransaction interface {
	Do(ctx context.Context, fn func(repos AuthTxRepositories) error) error
}
