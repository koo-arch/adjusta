package validation

import (
	"github.com/koo-arch/adjusta-backend/api/dto"
)

func FinalizeValidation(confirmEvent *dto.ConfirmEvent) error {
	validationErrors := NewValidationErrors()

	// confirm_dateのバリデーション
	confirmDate := confirmEvent.ConfirmDate
	if confirmDate.Start == nil {
		validationErrors.AddWithCode("confirm_date.start", "date_required")
	}

	if confirmDate.End == nil {
		validationErrors.AddWithCode("confirm_date.end", "date_required")
	}

	if confirmDate.Start != nil && confirmDate.End != nil && !confirmDate.Start.Before(*confirmDate.End) {
		validationErrors.AddWithCode("confirm_date", "dates_invalid")
	}

	// エラーがあればエラーを返す
	if validationErrors.HasErrors() {
		return validationErrors.ToAPIError()
	}

	return nil
}
