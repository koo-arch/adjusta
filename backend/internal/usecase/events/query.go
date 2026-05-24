package events

import (
	"time"

	"github.com/koo-arch/adjusta-backend/internal/models"
)

type SearchDraftQuery struct {
	Title        *string
	Location     *string
	Description  *string
	Status       *models.EventStatus
	StartTimeGTE *time.Time
	StartTimeLTE *time.Time
	EndTimeGTE   *time.Time
	EndTimeLTE   *time.Time
}
