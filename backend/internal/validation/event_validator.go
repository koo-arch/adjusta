package validation

import (
	"github.com/koo-arch/adjusta-backend/internal/models"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func CreateEventValidation(eventDraft *models.EventDraftCreation) *internalErrors.ValidationError {
	validationErrors := NewValidationErrors()

	// 基本情報のバリデーション
	validateBasicFields(eventDraft.Title, eventDraft.Description, eventDraft.Location, validationErrors)

	// 選択日時一つ以上必須
	if len(eventDraft.SelectedDates) == 0 {
		validationErrors.AddWithCode("selected_dates", "dates_required")
	} else if len(eventDraft.SelectedDates) > 10 {
		validationErrors.AddWithCode("selected_dates", "dates_max_length")
	}

	// 選択日時の開始時刻と終了時刻
	for _, date := range eventDraft.SelectedDates {
		if date.Start.After(date.End) {
			validationErrors.AddWithCode("selected_dates", "dates_invalid")
		}
	}

	// エラーが存在する場合はエラーを返す
	if validationErrors.HasErrors() {
		return validationErrors.ToAPIErrors()
	}

	return nil
}

func UpdateEventValidation(eventDraft *models.EventDraftUpdate) error {
	validationErrors := NewValidationErrors()

	// 基本情報のバリデーション
	validateBasicFields(eventDraft.Title, eventDraft.Description, eventDraft.Location, validationErrors)

	// 選択日時一つ以上必須
	if len(eventDraft.ProposedDates) == 0 {
		validationErrors.AddWithCode("proposed_dates", "dates_required")
	}

	// 選択日時の開始時刻と終了時刻
	for _, date := range eventDraft.ProposedDates {
		if date.Start.After(*date.End) {
			validationErrors.AddWithCode("proposed_dates", "dates_invalid")
		}
	}

	// エラーが存在する場合はエラーを返す
	if validationErrors.HasErrors() {
		return validationErrors.ToAPIErrors()
	}

	return nil
}

func validateBasicFields(title, description, location string, validationErrors *ValidationErrors) {
	if title == "" {
		validationErrors.AddWithCode("title", "title_required")
	} else if len(title) > 100 {
		validationErrors.AddWithCode("title", "title_max_length")
	}

	if len(description) > 500 {
		validationErrors.AddWithCode("description", "description_max_length")
	}

	if len(location) > 100 {
		validationErrors.AddWithCode("location", "location_max_length")
	}
}