package transaction

import "github.com/koo-arch/adjusta-backend/ent"

type Tx interface {
	Client() *ent.Client
	Commit() error
	Rollback() error
}

type wrappedTx struct {
	tx *ent.Tx
}

func Wrap(tx *ent.Tx) Tx {
	if tx == nil {
		return nil
	}

	return &wrappedTx{tx: tx}
}

func (tx *wrappedTx) Client() *ent.Client {
	return tx.tx.Client()
}

func (tx *wrappedTx) Commit() error {
	return tx.tx.Commit()
}

func (tx *wrappedTx) Rollback() error {
	return tx.tx.Rollback()
}
