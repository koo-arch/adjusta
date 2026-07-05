package repository

import (
	"context"

	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent"
	infraTransaction "github.com/koo-arch/adjusta-backend/internal/infrastructure/transaction"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(repos Repositories) error) error
}

type EntUnitOfWork struct {
	client *ent.Client
}

func NewUnitOfWork(client *ent.Client) *EntUnitOfWork {
	return &EntUnitOfWork{
		client: client,
	}
}

func (u *EntUnitOfWork) Do(ctx context.Context, fn func(repos Repositories) error) error {
	entTx, err := u.client.Tx(ctx)
	if err != nil {
		return err
	}

	txErr := fn(NewRepositories(entTx.Client()))
	infraTransaction.Handle(entTx, &txErr)

	return txErr
}
