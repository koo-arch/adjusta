package validation

import internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"

type ValidationErrors struct {
	Details map[string]string `json:"details,omitempty"`
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Details: make(map[string]string),
	}
}

func (v *ValidationErrors) Add(field, message string) {
	v.Details[field] = message
}

func (v *ValidationErrors) AddWithCode(field, code string) {
	if message, ok := ErrorMessages[code]; ok {
		v.Details[field] = message
	}
}

func (v *ValidationErrors) ToAPIErrors() *internalErrors.ValidationError {
	return internalErrors.NewValidationError(v.Details)
}

func (v *ValidationErrors) HasErrors() bool {
	return len(v.Details) > 0
}

var ErrorMessages = map[string]string{
	"title_required": 	"タイトルを入力してください",
	"title_max_length": "タイトルは100文字以内で入力してください",
	"description_max_length": "説明文は500文字以内で入力してください",
	"location_max_length": 	"場所は100文字以内で入力してください",
	"dates_required": "日程を一つ以上選択してください",
	"dates_max_length": "日程は10個以内で選択してください",
	"dates_invalid": "開始時刻は終了時刻より前に設定してください",
	"date_required": "日程を選択してください",
}