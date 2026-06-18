package repository

import (
	"context"

	"github.com/koo-arch/adjusta-backend/ent"
	infraTransaction "github.com/koo-arch/adjusta-backend/internal/infrastructure/transaction"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(repos Repositories) error) error
}

type EntUnitOfWork struct {
	client *ent.Client
	repos  Repositories
}

func NewUnitOfWork(client *ent.Client, repos Repositories) *EntUnitOfWork {
	return &EntUnitOfWork{
		client: client,
		repos:  repos,
	}
}

func (u *EntUnitOfWork) Do(ctx context.Context, fn func(repos Repositories) error) error {
	entTx, err := u.client.Tx(ctx)
	if err != nil {
		return err
	}

	tx := infraTransaction.Wrap(entTx)
	txErr := fn(u.repos.WithTx(tx))
	infraTransaction.Handle(tx, &txErr)

	return txErr
}
