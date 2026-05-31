package events

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func (uc *Usecase) CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *appmodel.EventDraftCreation) (*appmodel.EventDraftDetail, error) {
	var response *appmodel.EventDraftDetail

	err := uc.tx.Do(ctx, func(store EventTxStore) error {
		storedCalendar, err := uc.findPrimaryCalendar(ctx, store, userID, email)
		if err != nil {
			return err
		}

		storedEvent, err := store.CreateEvent(ctx, storedCalendar.ID, eventReq.Title, eventReq.Location, eventReq.Description, eventReq.SelectedDates[0].Start, eventReq.SelectedDates[0].End)
		if err != nil {
			log.Printf("failed to create event for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		storedDates, err := store.CreateProposedDates(ctx, eventReq.SelectedDates, storedEvent.ID)
		if err != nil {
			log.Printf("failed to create proposed dates for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		response = &appmodel.EventDraftDetail{
			ID:            storedEvent.ID,
			Title:         storedEvent.Summary,
			Location:      storedEvent.Location,
			Description:   storedEvent.Description,
			Status:        storedEvent.Status,
			ProposedDates: buildProposedDates(storedDates),
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running create drafted event transaction: %v", err)
		return nil, normalizeUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return response, nil
}
