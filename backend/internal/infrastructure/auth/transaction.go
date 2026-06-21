package auth

import (
	"context"

	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
)

type authTransaction struct {
	uow infraRepository.UnitOfWork
}

func NewAuthTransaction(uow infraRepository.UnitOfWork) usecaseAuth.SignInTransaction {
	return &authTransaction{uow: uow}
}

func (t *authTransaction) Do(ctx context.Context, fn func(repos usecaseAuth.AuthTxRepositories) error) error {
	return t.uow.Do(ctx, func(repos infraRepository.Repositories) error {
		return fn(usecaseAuth.AuthTxRepositories{
			User:    repos.User,
			Account: repos.Account,
			Session: repos.Session,
		})
	})
}
