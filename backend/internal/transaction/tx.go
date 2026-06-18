package transaction

import "github.com/koo-arch/adjusta-backend/ent"

type Tx interface {
	Client() *ent.Client
	Commit() error
	Rollback() error
}
