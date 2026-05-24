package events

import (
	"context"
	"log"
	"sort"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) findPrimaryCalendar(ctx context.Context, finder PrimaryCalendarFinder, userID uuid.UUID, email string) (*models.StoredCalendar, error) {
	storedCalendar, err := finder.FindPrimaryCalendar(ctx, userID)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return storedCalendar, nil
}

func buildProposedDates(storedDates []*models.StoredProposedDate) []models.ProposedDate {
	proposedDates := make([]models.ProposedDate, 0, len(storedDates))
	for _, storedDate := range storedDates {
		proposedDates = append(proposedDates, models.ProposedDate{
			ID:       &storedDate.ID,
			Start:    &storedDate.StartTime,
			End:      &storedDate.EndTime,
			Priority: storedDate.Priority,
		})
	}

	sort.Slice(proposedDates, func(i, j int) bool {
		return proposedDates[i].Priority < proposedDates[j].Priority
	})

	return proposedDates
}

func buildEventDraftDetail(storedEvent *models.StoredEvent) (*models.EventDraftDetail, error) {
	if storedEvent.ProposedDates == nil {
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return &models.EventDraftDetail{
		ID:              storedEvent.ID,
		Title:           storedEvent.Summary,
		Location:        storedEvent.Location,
		Description:     storedEvent.Description,
		Status:          storedEvent.Status,
		ConfirmedDateID: &storedEvent.ConfirmedDateID,
		GoogleEventID:   storedEvent.GoogleEventID,
		Slug:            storedEvent.Slug,
		ProposedDates:   buildProposedDates(storedEvent.ProposedDates),
	}, nil
}

func normalizeUsecaseError(err error, fallbackMessage string) error {
	return internalErrors.NormalizeAPIError(err, fallbackMessage)
}
