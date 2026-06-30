package events

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) FetchNeedsActionDrafts(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]NeedsActionDraftOutput, error) {
	storedCalendar, err := uc.loadPrimaryCalendar(ctx, uc.repos, userID, email)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	startTime := currentTime.AddDate(0, 0, daysBefore)
	active := value.StatusActive
	eventOptions := EventSearchOptions{
		WithProposedDates: true,
		Status:            &active,
		StartTimeLTE:      &startTime,
		SortBy:            "ProposedDatePriority",
		SortOrder:         "desc",
	}

	storedEvents, err := uc.repos.Event.SearchEvents(ctx, userID, storedCalendar.ID, toEventSearchOptions(eventOptions))
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	needsActionDrafts := make([]NeedsActionDraftOutput, 0)
	for _, storedEvent := range storedEvents {
		needsActionDraft, err := buildNeedsActionDraftOutput(storedEvent, currentTime)
		if err != nil {
			log.Printf("No association found between calendar and event")
			return nil, err
		}
		if needsActionDraft != nil {
			needsActionDrafts = append(needsActionDrafts, *needsActionDraft)
		}
	}

	return needsActionDrafts, nil
}
