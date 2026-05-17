package validation

import (
	"github.com/koo-arch/adjusta-backend/internal/models"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func FinalizeValidation(confirmEvent *models.ConfirmEvent) *internalErrors.ValidationError {
	validationErrors := NewValidationErrors()

	// confirm_dateのバリデーション
	confirmDate := confirmEvent.ConfirmDate
	if confirmDate.Start == nil {
		validationErrors.AddWithCode("confirm_date.start", "date_required")
	}

	if confirmDate.End == nil {
		validationErrors.AddWithCode("confirm_date.end", "date_required")
	}

	if confirmDate.Start.After(*confirmDate.End) || confirmDate.Start.Equal(*confirmDate.End) {
		validationErrors.AddWithCode("confirm_date", "dates_invalid")
	}

	// エラーがあればエラーを返す
	if validationErrors.HasErrors() {
		return validationErrors.ToAPIErrors()
	}

	return nil
}