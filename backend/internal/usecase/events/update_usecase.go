package events

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) UpdateDraftedEvents(ctx context.Context, userID uuid.UUID, slug, email string, eventReq *models.EventDraftUpdate) error {
	err := uc.tx.Do(ctx, func(store EventTxStore) error {
		if _, err := uc.findPrimaryCalendar(ctx, store, userID, email); err != nil {
			return err
		}

		storedEvent, err := store.FindEventBySlug(ctx, userID, slug, false)
		if err != nil {
			log.Printf("failed to get event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("イベントが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		eventOptions := EventMutation{
			Title:       &eventReq.Title,
			Location:    &eventReq.Location,
			Description: &eventReq.Description,
			Status:      &eventReq.Status,
		}
		storedEvent, err = store.UpdateEvent(ctx, storedEvent.ID, eventOptions)
		if err != nil {
			log.Printf("failed to update event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("イベントが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		existingDates, err := store.ListProposedDatesByEvent(ctx, storedEvent.ID)
		if err != nil {
			log.Printf("failed to get proposed dates for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		if eventReq.Status == models.StatusConfirmed {
			confirmDate := models.ConfirmDate{
				ID:       eventReq.ProposedDates[0].ID,
				Start:    eventReq.ProposedDates[0].Start,
				End:      eventReq.ProposedDates[0].End,
				Priority: eventReq.ProposedDates[0].Priority,
			}
			confirmEvent := models.ConfirmEvent{
				ConfirmDate: confirmDate,
			}

			googleEventID, err := uc.handleGoogleEvent(ctx, userID, storedEvent, &confirmEvent)
			if err != nil {
				log.Printf("failed to handle google event for account: %s, error: %v", email, err)
				return internalErrors.NormalizeAPIError(err, "Googleカレンダー更新時にエラーが発生しました")
			}

			if err := uc.confirmEventDate(ctx, store, googleEventID, &confirmEvent, storedEvent); err != nil {
				log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
				return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
			}
		}

		if err := uc.updateProposedDates(ctx, store, eventReq, storedEvent, existingDates); err != nil {
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running update drafted event transaction: %v", err)
		return normalizeUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return nil
}

func (uc *Usecase) updateProposedDates(ctx context.Context, store EventTxStore, eventReq *models.EventDraftUpdate, storedEvent *models.StoredEvent, existingDates []*models.StoredProposedDate) error {
	updateDateMap := make(map[uuid.UUID]models.ProposedDate)
	for _, date := range eventReq.ProposedDates {
		if date.ID != nil {
			updateDateMap[*date.ID] = date
		} else {
			updateDateMap[uuid.New()] = date
		}
	}

	for _, date := range existingDates {
		if updateDate, ok := updateDateMap[date.ID]; ok {
			dateOptions := ProposedDateMutation{
				Start:    updateDate.Start,
				End:      updateDate.End,
				Priority: &updateDate.Priority,
			}
			if _, err := store.UpdateProposedDate(ctx, date.ID, dateOptions); err != nil {
				return fmt.Errorf("failed to update proposed date for account: %s, error: %w", updateDate.ID, err)
			}
			delete(updateDateMap, date.ID)
		} else {
			if err := store.DeleteProposedDate(ctx, date.ID); err != nil {
				return fmt.Errorf("failed to delete proposed date for account: %s, error: %w", date.ID, err)
			}
		}
	}

	for _, date := range updateDateMap {
		dateOptions := ProposedDateMutation{
			Start:    date.Start,
			End:      date.End,
			Priority: &date.Priority,
		}
		if _, err := store.CreateProposedDate(ctx, dateOptions, storedEvent.ID); err != nil {
			return fmt.Errorf("failed to create proposed date for account: %s, error: %w", date.ID, err)
		}
	}

	return nil
}
