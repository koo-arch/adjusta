package events

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	internalRepo "github.com/koo-arch/adjusta-backend/internal/repo"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
	"github.com/koo-arch/adjusta-backend/utils"
)

func (uc *Usecase) FinalizeProposedDate(ctx context.Context, userID uuid.UUID, slug, email string, eventReq *models.ConfirmEvent) error {
	err := uc.uow.Do(ctx, func(repos internalRepo.Repositories) error {
		entEvent, err := repos.Event.FindBySlugAndUser(ctx, userID, slug, event.EventQueryOptions{})
		if err != nil {
			log.Printf("failed to get event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
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

		if err := uc.confirmEventDate(ctx, repos, googleEventID, eventReq, entEvent); err != nil {
			log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running finalize proposed date transaction: %v", err)
		return normalizeUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return nil
}

func (uc *Usecase) handleGoogleEvent(calendarService *customCalendar.Calendar, storedEvent *models.StoredEvent, eventReq *models.ConfirmEvent) (*string, error) {
	var googleEventID *string
	if eventReq.ConfirmDate.ID == nil || storedEvent.GoogleEventID == "" {
		eventDraftCreate := models.EventDraftCreation{
			Title:       storedEvent.Summary,
			Location:    storedEvent.Location,
			Description: storedEvent.Description,
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
		convertGoogleEvent := uc.calendarApp.ConvertToCalendarEvent(&storedEvent.GoogleEventID, storedEvent.Summary, storedEvent.Location, storedEvent.Description, *eventReq.ConfirmDate.Start, *eventReq.ConfirmDate.End)
		googleEvent, err := uc.calendarApp.UpdateOrCreateGoogleEvent(calendarService, convertGoogleEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to update events, error: %w", err)
		}

		googleEventID = &googleEvent.Id
	}

	return googleEventID, nil
}

func (uc *Usecase) confirmEventDate(ctx context.Context, repos internalRepo.Repositories, googleEventID *string, eventReq *models.ConfirmEvent, storedEvent *models.StoredEvent) error {
	priority := 1
	dateOptions := proposeddate.ProposedDateQueryOptions{
		Priority: &priority,
	}

	confirmDateID := eventReq.ConfirmDate.ID
	if eventReq.ConfirmDate.ID == nil {
		dateOptions.StartTime = eventReq.ConfirmDate.Start
		dateOptions.EndTime = eventReq.ConfirmDate.End

		storedDate, err := repos.ProposedDate.Create(ctx, dateOptions, storedEvent.ID)
		if err != nil {
			return fmt.Errorf("failed to create proposed date error: %w", err)
		}
		confirmDateID = &storedDate.ID

		if err := repos.ProposedDate.DecrementPriorityExceptID(ctx, storedEvent.ID, storedDate.ID); err != nil {
			return fmt.Errorf("failed to decrement priority error: %w", err)
		}
	}

	if eventReq.ConfirmDate.ID != nil {
		zero := 0
		dateOptions.Priority = &zero

		if _, err := repos.ProposedDate.Update(ctx, *eventReq.ConfirmDate.ID, dateOptions); err != nil {
			return fmt.Errorf("failed to update proposed date error: %w", err)
		}

		if err := repos.ProposedDate.ReorderPriority(ctx, storedEvent.ID); err != nil {
			return fmt.Errorf("failed to reorder priority error: %w", err)
		}
	}

	status := models.StatusConfirmed
	eventOptions := event.EventQueryOptions{
		Status:          &status,
		ConfirmedDateID: confirmDateID,
		GoogleEventID:   googleEventID,
	}
	if _, err := repos.Event.Update(ctx, storedEvent.ID, eventOptions); err != nil {
		return fmt.Errorf("failed to update event status error: %w", err)
	}

	return nil
}
