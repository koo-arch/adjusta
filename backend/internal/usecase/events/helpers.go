package events

import (
	"context"
	"log"
	"sort"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
)

func (uc *Usecase) findPrimaryCalendar(ctx context.Context, finder PrimaryCalendarFinder, userID uuid.UUID, email string) (*repositorymodel.StoredCalendar, error) {
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

func buildProposedDates(storedDates []*repositorymodel.StoredProposedDate) []appmodel.ProposedDate {
	proposedDates := make([]appmodel.ProposedDate, 0, len(storedDates))
	for _, storedDate := range storedDates {
		proposedDates = append(proposedDates, appmodel.ProposedDate{
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

func buildEventDraftDetail(storedEvent *repositorymodel.StoredEvent) (*appmodel.EventDraftDetail, error) {
	if storedEvent.ProposedDates == nil {
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return &appmodel.EventDraftDetail{
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
