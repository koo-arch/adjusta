package events

import (
	"context"
	"log"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) FetchAllDraftedEvents(ctx context.Context, userID uuid.UUID, email string) ([]*EventDraftDetailOutput, error) {
	storedCalendar, err := uc.loadPrimaryCalendar(ctx, uc.repos, userID, email)
	if err != nil {
		return nil, err
	}

	eventOptions := EventSearchOptions{
		WithProposedDates: true,
	}
	storedEvents, err := uc.repos.Event.SearchEvents(ctx, userID, storedCalendar.ID, toEventSearchOptions(eventOptions))
	if err != nil {
		log.Printf("failed to get events for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	draftedEvents := make([]*EventDraftDetailOutput, 0, len(storedEvents))
	for _, storedEvent := range storedEvents {
		draft, err := buildEventDraftDetailOutput(storedEvent)
		if err != nil {
			return nil, err
		}
		draftedEvents = append(draftedEvents, draft)
	}

	return draftedEvents, nil
}

func (uc *Usecase) SearchDraftedEvents(ctx context.Context, userID uuid.UUID, email string, query SearchDraftQuery) ([]*EventDraftDetailOutput, error) {
	storedCalendar, err := uc.loadPrimaryCalendar(ctx, uc.repos, userID, email)
	if err != nil {
		return nil, err
	}

	eventOptions := EventSearchOptions{
		WithProposedDates: true,
		Title:             query.Title,
		Location:          query.Location,
		Description:       query.Description,
		Status:            query.Status,
		StartTimeGTE:      query.StartTimeGTE,
		StartTimeLTE:      query.StartTimeLTE,
		EndTimeGTE:        query.EndTimeGTE,
		EndTimeLTE:        query.EndTimeLTE,
	}
	storedEvents, err := uc.repos.Event.SearchEvents(ctx, userID, storedCalendar.ID, toEventSearchOptions(eventOptions))
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	searchResult := make([]*EventDraftDetailOutput, 0, len(storedEvents))
	for _, storedEvent := range storedEvents {
		draft, err := buildEventDraftDetailOutput(storedEvent)
		if err != nil {
			return nil, err
		}
		searchResult = append(searchResult, draft)
	}

	return searchResult, nil
}
