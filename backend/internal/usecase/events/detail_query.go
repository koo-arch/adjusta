package events

import (
	"context"
	"log"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func (uc *Usecase) FetchDraftedEventDetail(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) (*EventDraftDetailOutput, error) {
	if uc.tx == nil {
		storedEvent, err := uc.loadDraftedEventDetailRecord(ctx, uc.repos, userID, email, eventID)
		if err != nil {
			return nil, err
		}
		return buildEventDraftDetailOutput(storedEvent)
	}

	var response *EventDraftDetailOutput

	err := uc.tx.DoEvent(ctx, func(repos EventTxRepositories) error {
		storedEvent, err := uc.loadDraftedEventDetailWithSync(ctx, repos, userID, email, eventID)
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
