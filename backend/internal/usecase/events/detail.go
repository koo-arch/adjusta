package events

import (
	"context"
	"log"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) FetchDraftedEventDetail(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) (*EventDraftDetailOutput, error) {
	if uc.tx != nil {
		var response *EventDraftDetailOutput
		err := uc.tx.DoEvent(ctx, func(repos EventTxRepositories) error {
			storedEvent, err := uc.loadDraftedEventDetailRecord(ctx, repos, userID, email, eventID)
			if err != nil {
				return err
			}
			response, err = buildEventDraftDetailOutput(storedEvent)
			return err
		})
		if err != nil {
			log.Printf("failed running fetch drafted event detail transaction: %v", err)
			return nil, mapUsecaseError(err, internalErrors.InternalErrorMessage)
		}
		return response, nil
	}

	storedEvent, err := uc.loadDraftedEventDetailRecord(ctx, uc.repos, userID, email, eventID)
	if err != nil {
		return nil, err
	}
	return buildEventDraftDetailOutput(storedEvent)
}

func (uc *Usecase) loadDraftedEventDetailRecord(ctx context.Context, repos EventTxRepositories, userID uuid.UUID, email string, eventID uuid.UUID) (*domainEvent.Event, error) {
	storedEvent, err := repos.Event.FindByIDAndUser(ctx, userID, eventID, domainEvent.EventReadOptions{
		WithProposedDates: true,
	})
	if err != nil {
		log.Printf("failed to get event detail for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return storedEvent, nil
}
