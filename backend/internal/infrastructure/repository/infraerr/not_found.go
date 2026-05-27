package infraerr

import (
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func MapNotFound(err error) error {
	if err == nil {
		return nil
	}

	if ent.IsNotFound(err) {
		return repoerr.ErrNotFound
	}

	return err
}
