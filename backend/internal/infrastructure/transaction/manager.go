package transaction

import (
	"log"

	"github.com/koo-arch/adjusta-backend/ent"
	sharedtx "github.com/koo-arch/adjusta-backend/internal/transaction"
)

type wrappedTx struct {
	tx *ent.Tx
}

func Wrap(tx *ent.Tx) sharedtx.Tx {
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

func Handle(tx sharedtx.Tx, txErr *error) {
	if p := recover(); p != nil {
		log.Printf("Rollback transaction: %v", p)
		if err := tx.Rollback(); err != nil {
			log.Printf("Failed rolling back transaction: %v", err)
		}
		panic(p)
	} else if *txErr != nil {
		log.Printf("Rollback transaction: %v", *txErr)
		if err := tx.Rollback(); err != nil {
			log.Printf("Failed rolling back transaction: %v", err)
		}
	} else {
		if err := tx.Commit(); err != nil {
			log.Printf("Failed committing transaction: %v", err)
		}
	}
}
