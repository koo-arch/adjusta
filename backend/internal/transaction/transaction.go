package transaction

import (
	"log"

	"github.com/koo-arch/adjusta-backend/ent"
)

func HandleTransaction(tx *ent.Tx, txErr *error) {
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