package events

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"github.com/koo-arch/adjusta-backend/utils"
)

func (uc *Usecase) UpdateDraftedEvents(ctx context.Context, userID uuid.UUID, slug, email string, eventReq *models.EventDraftUpdate) error {
	tx, err := uc.client.Tx(ctx)
	if err != nil {
		log.Printf("failed starting transaction: %v", err)
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	defer transaction.HandleTransaction(tx, &err)

	if _, err := uc.findPrimaryCalendar(ctx, tx, userID, email); err != nil {
		return err
	}

	entEvent, err := uc.eventRepo.FindBySlugAndUser(ctx, tx, userID, slug, event.EventQueryOptions{})
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	eventOptions := event.EventQueryOptions{
		Summary:     &eventReq.Title,
		Location:    &eventReq.Location,
		Description: &eventReq.Description,
		Status:      &eventReq.Status,
	}
	entEvent, err = uc.eventRepo.Update(ctx, tx, entEvent.ID, eventOptions)
	if err != nil {
		log.Printf("failed to update event for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	existingDates, err := uc.dateRepo.FilterByEventID(ctx, tx, entEvent.ID)
	if err != nil {
		log.Printf("failed to get proposed dates for account: %s, error: %v", email, err)
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
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

		calendarService, err := uc.getGoogleCalendarService(ctx, userID, email)
		if err != nil {
			return err
		}

		googleEventID, err := uc.handleGoogleEvent(calendarService, entEvent, &confirmEvent)
		if err != nil {
			log.Printf("failed to handle google event for account: %s, error: %v", email, err)
			return utils.HandleGoogleAPIError(err)
		}

		if err := uc.confirmEventDate(ctx, tx, googleEventID, &confirmEvent, entEvent); err != nil {
			log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}
	}

	if err := uc.updateProposedDates(ctx, tx, eventReq, entEvent, existingDates); err != nil {
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return nil
}

func (uc *Usecase) updateProposedDates(ctx context.Context, tx *ent.Tx, eventReq *models.EventDraftUpdate, entEvent *ent.Event, existingDates []*ent.ProposedDate) error {
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
			dateOptions := proposeddate.ProposedDateQueryOptions{
				StartTime: updateDate.Start,
				EndTime:   updateDate.End,
				Priority:  &updateDate.Priority,
			}
			if _, err := uc.dateRepo.Update(ctx, tx, date.ID, dateOptions); err != nil {
				return fmt.Errorf("failed to update proposed date for account: %s, error: %w", updateDate.ID, err)
			}
			delete(updateDateMap, date.ID)
		} else {
			if err := uc.dateRepo.Delete(ctx, tx, date.ID); err != nil {
				return fmt.Errorf("failed to delete proposed date for account: %s, error: %w", date.ID, err)
			}
		}
	}

	for _, date := range updateDateMap {
		dateOptions := proposeddate.ProposedDateQueryOptions{
			StartTime: date.Start,
			EndTime:   date.End,
			Priority:  &date.Priority,
		}
		if _, err := uc.dateRepo.Create(ctx, tx, dateOptions, entEvent); err != nil {
			return fmt.Errorf("failed to create proposed date for account: %s, error: %w", date.ID, err)
		}
	}

	return nil
}
