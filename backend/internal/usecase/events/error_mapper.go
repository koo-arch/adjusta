package events

import (
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func mapUsecaseError(err error, fallbackMessage string) error {
	return internalErrors.NormalizeAPIError(err, fallbackMessage)
}
