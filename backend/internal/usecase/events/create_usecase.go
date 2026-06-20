package events

import (
	"context"
	"log"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func (uc *Usecase) CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq DraftCreationRequest) (*EventDraftDetailOutput, error) {
	var response *EventDraftDetailOutput

	err := uc.tx.Do(ctx, func(store EventTxStore) error {
		storedCalendar, err := uc.loadPrimaryCalendar(ctx, store, userID, email)
		if err != nil {
			return err
		}
		candidateCalendar, err := uc.loadAdjustaCandidateCalendar(ctx, store, userID, email)
		if err != nil {
			return err
		}

		storedEvent, err := store.CreateEvent(ctx, userID, storedCalendar.ID, eventReq.Title, eventReq.Location, eventReq.Description, eventReq.SelectedDates[0].Start, eventReq.SelectedDates[0].End)
		if err != nil {
			log.Printf("failed to create event for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		storedDates, err := store.CreateProposedDates(ctx, assignSelectedDatePriorities(eventReq.SelectedDates), storedEvent.ID)
		if err != nil {
			log.Printf("failed to create proposed dates for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		if candidateCalendar.SyncProposedDates {
			storedEvent, err = store.UpdateEvent(ctx, storedEvent.ID, mergeEventChange(EventMutation{}, domainEvent.NewPendingEventChange(nil)))
			if err != nil {
				log.Printf("failed to mark event sync pending for account: %s, error: %v", email, err)
				return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
			}

			for i, storedDate := range storedDates {
				storedDates[i], err = store.UpdateProposedDate(ctx, storedDate.ID, buildProposedDateMutation(domainEvent.NewPendingProposedDateChange(nil, nil, nil, nil)))
				if err != nil {
					log.Printf("failed to mark proposed date sync pending for account: %s, error: %v", email, err)
					return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
				}
			}
		}

		response = &EventDraftDetailOutput{
			ID:            storedEvent.ID,
			Title:         storedEvent.Title,
			Location:      storedEvent.Location,
			Description:   storedEvent.Description,
			Status:        storedEvent.Status,
			SyncStatus:    storedEvent.SyncStatus,
			GoogleEventID: domainEvent.ResolveGoogleEventID(storedEvent.ConfirmedGoogleEventID),
			ProposedDates: buildProposedDateOutputs(storedDates),
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running create drafted event transaction: %v", err)
		return nil, mapUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return response, nil
}
