package events

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"github.com/koo-arch/adjusta-backend/utils"
)

func (uc *Usecase) FinalizeProposedDate(ctx context.Context, userID uuid.UUID, slug, email string, eventReq *models.ConfirmEvent) error {
	tx, err := uc.client.Tx(ctx)
	if err != nil {
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	defer transaction.HandleTransaction(tx, &err)

	entEvent, err := uc.eventRepo.FindBySlugAndUser(ctx, tx, userID, slug, event.EventQueryOptions{})
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	calendarService, err := uc.getGoogleCalendarService(ctx, userID, email)
	if err != nil {
		return utils.GetAPIError(err, "サーバーでエラーが発生しました")
	}

	googleEventID, err := uc.handleGoogleEvent(calendarService, entEvent, eventReq)
	if err != nil {
		log.Printf("failed to handle google event for account: %s, error: %v", email, err)
		return utils.HandleGoogleAPIError(err)
	}

	if err := uc.confirmEventDate(ctx, tx, googleEventID, eventReq, entEvent); err != nil {
		log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return nil
}

func (uc *Usecase) handleGoogleEvent(calendarService *customCalendar.Calendar, entEvent *ent.Event, eventReq *models.ConfirmEvent) (*string, error) {
	var googleEventID *string
	if eventReq.ConfirmDate.ID == nil || entEvent.GoogleEventID == "" {
		eventDraftCreate := models.EventDraftCreation{
			Title:       entEvent.Summary,
			Location:    entEvent.Location,
			Description: entEvent.Description,
			SelectedDates: []models.SelectedDate{
				{
					Start: *eventReq.ConfirmDate.Start,
					End:   *eventReq.ConfirmDate.End,
				},
			},
		}

		googleEvents, err := uc.calendarApp.CreateGoogleEvents(calendarService, &eventDraftCreate)
		if err != nil {
			return nil, fmt.Errorf("failed to insert events, error: %w", err)
		}
		googleEventID = &googleEvents[0].Id
	} else {
		convertGoogleEvent := uc.calendarApp.ConvertToCalendarEvent(&entEvent.GoogleEventID, entEvent.Summary, entEvent.Location, entEvent.Description, *eventReq.ConfirmDate.Start, *eventReq.ConfirmDate.End)
		googleEvent, err := uc.calendarApp.UpdateOrCreateGoogleEvent(calendarService, convertGoogleEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to update events, error: %w", err)
		}

		googleEventID = &googleEvent.Id
	}

	return googleEventID, nil
}

func (uc *Usecase) confirmEventDate(ctx context.Context, tx *ent.Tx, googleEventID *string, eventReq *models.ConfirmEvent, entEvent *ent.Event) error {
	priority := 1
	dateOptions := proposeddate.ProposedDateQueryOptions{
		Priority: &priority,
	}

	confirmDateID := eventReq.ConfirmDate.ID
	if eventReq.ConfirmDate.ID == nil {
		dateOptions.StartTime = eventReq.ConfirmDate.Start
		dateOptions.EndTime = eventReq.ConfirmDate.End

		entDate, err := uc.dateRepo.Create(ctx, tx, dateOptions, entEvent)
		if err != nil {
			return fmt.Errorf("failed to create proposed date error: %w", err)
		}
		confirmDateID = &entDate.ID

		if err := uc.dateRepo.DecrementPriorityExceptID(ctx, tx, entEvent.ID, entDate.ID); err != nil {
			return fmt.Errorf("failed to decrement priority error: %w", err)
		}
	}

	if eventReq.ConfirmDate.ID != nil {
		zero := 0
		dateOptions.Priority = &zero

		entDate, err := uc.dateRepo.Update(ctx, tx, *eventReq.ConfirmDate.ID, dateOptions)
		if err != nil {
			return fmt.Errorf("failed to update proposed date error: %w", err)
		}

		if err := uc.dateRepo.ReorderPriority(ctx, tx, entDate.ID); err != nil {
			return fmt.Errorf("failed to reorder priority error: %w", err)
		}
	}

	status := models.StatusConfirmed
	eventOptions := event.EventQueryOptions{
		Status:          &status,
		ConfirmedDateID: confirmDateID,
		GoogleEventID:   googleEventID,
	}
	if _, err := uc.eventRepo.Update(ctx, tx, entEvent.ID, eventOptions); err != nil {
		return fmt.Errorf("failed to update event status error: %w", err)
	}

	return nil
}
